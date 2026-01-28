package raknet

import (
	"bytes"
	"encoding/binary"
	"net"
	"sync"
	"time"

	"github.com/scaxe/scaxe-go/pkg/logger"
)

type splitPacketData struct {
	splitCount uint32
	fragments  map[uint32][]byte
}

type Session struct {
	server    *Server
	addr      *net.UDPAddr
	mtu       uint16
	clientID  uint64
	connected bool

	sendSeqNum    uint32
	receiveSeqNum uint32

	splitID      uint16
	messageIndex uint32

	receivedPackets map[uint32]bool
	ackQueue        []uint32

	splitPackets map[uint16]*splitPacketData

	mu sync.Mutex

	lastActivity time.Time
}

func NewSession(server *Server, addr *net.UDPAddr, mtu uint16, clientID uint64) *Session {
	logger.Debug("raknet.NewSession", "address", addr.String(), "mtu", mtu, "clientID", clientID)
	return &Session{
		server:          server,
		addr:            addr,
		mtu:             mtu,
		clientID:        clientID,
		receivedPackets: make(map[uint32]bool),
		splitPackets:    make(map[uint16]*splitPacketData),
		lastActivity:    time.Now(),
	}
}

func (s *Session) Address() string {
	return s.addr.String()
}

func (s *Session) handleDataPacket(data []byte) {
	if len(data) < 4 {
		return
	}

	s.lastActivity = time.Now()

	seqNum := uint32(data[1]) | uint32(data[2])<<8 | uint32(data[3])<<16
	logger.Debug("raknet.Session.handleDataPacket", "seqNum", seqNum, "size", len(data))

	s.mu.Lock()
	s.ackQueue = append(s.ackQueue, seqNum)
	s.receivedPackets[seqNum] = true
	s.mu.Unlock()

	s.sendACK()

	offset := 4
	for offset < len(data) {
		encapsulated, newOffset := s.decodeEncapsulated(data, offset)
		if encapsulated == nil {
			break
		}
		offset = newOffset

		s.handleEncapsulatedPacket(encapsulated)
	}
}

type encapsulatedPacket struct {
	reliability  byte
	hasSplit     bool
	splitCount   uint32
	splitID      uint16
	splitIndex   uint32
	messageIndex uint32
	orderIndex   uint32
	orderChannel byte
	payload      []byte
}

func (s *Session) decodeEncapsulated(data []byte, offset int) (*encapsulatedPacket, int) {
	if offset >= len(data) {
		return nil, offset
	}

	pkt := &encapsulatedPacket{}

	flags := data[offset]
	offset++

	pkt.reliability = (flags & 0xe0) >> 5
	pkt.hasSplit = (flags & 0x10) != 0

	if offset+2 > len(data) {
		return nil, offset
	}
	lengthBits := binary.BigEndian.Uint16(data[offset:])
	offset += 2
	length := int((lengthBits + 7) / 8)

	if pkt.reliability >= 2 && pkt.reliability <= 4 {

		if offset+3 > len(data) {
			return nil, offset
		}
		pkt.messageIndex = uint32(data[offset]) | uint32(data[offset+1])<<8 | uint32(data[offset+2])<<16
		offset += 3
	}

	if pkt.reliability == 1 || pkt.reliability == 3 || pkt.reliability == 4 {

		if offset+4 > len(data) {
			return nil, offset
		}
		pkt.orderIndex = uint32(data[offset]) | uint32(data[offset+1])<<8 | uint32(data[offset+2])<<16
		offset += 3
		pkt.orderChannel = data[offset]
		offset++
	}

	if pkt.hasSplit {
		if offset+10 > len(data) {
			return nil, offset
		}
		pkt.splitCount = binary.BigEndian.Uint32(data[offset:])
		offset += 4
		pkt.splitID = binary.BigEndian.Uint16(data[offset:])
		offset += 2
		pkt.splitIndex = binary.BigEndian.Uint32(data[offset:])
		offset += 4
	}

	if offset+length > len(data) {
		length = len(data) - offset
	}
	pkt.payload = make([]byte, length)
	copy(pkt.payload, data[offset:offset+length])
	offset += length

	logger.Debug("raknet.decodeEncapsulated",
		"reliability", pkt.reliability,
		"hasSplit", pkt.hasSplit,
		"payloadSize", len(pkt.payload))

	return pkt, offset
}

func (s *Session) handleEncapsulatedPacket(pkt *encapsulatedPacket) {
	if len(pkt.payload) == 0 {
		return
	}

	if pkt.hasSplit {
		reassembled := s.handleSplitPacket(pkt)
		if reassembled == nil {

			return
		}

		pkt.payload = reassembled
	}

	packetID := pkt.payload[0]
	logger.Debug("raknet.handleEncapsulatedPacket", "innerPacketID", packetID, "payloadSize", len(pkt.payload))

	switch packetID {
	case IDConnectionRequest:
		s.handleConnectionRequest(pkt.payload)
	case IDNewIncomingConnection:
		s.handleNewIncomingConnection(pkt.payload)
	case IDDisconnectNotification:
		s.server.removeSession(s.addr.String())
	case 0xfe:

		if s.server.OnPacket != nil {
			s.server.OnPacket(s, pkt.payload[1:])
		}
	case 0x8e:

		if s.server.OnPacket != nil {
			s.server.OnPacket(s, pkt.payload[1:])
		}
	default:

		if packetID >= 0x8f && packetID <= 0xcb {
			if s.server.OnPacket != nil {
				s.server.OnPacket(s, pkt.payload)
			}
		} else {
			logger.Debug("raknet.handleEncapsulatedPacket", "unhandled", packetID)
		}
	}
}

func (s *Session) handleSplitPacket(pkt *encapsulatedPacket) []byte {
	s.mu.Lock()
	defer s.mu.Unlock()

	splitID := pkt.splitID
	splitIndex := pkt.splitIndex
	splitCount := pkt.splitCount

	logger.Debug("raknet.handleSplitPacket",
		"splitID", splitID,
		"splitIndex", splitIndex,
		"splitCount", splitCount,
		"fragmentSize", len(pkt.payload))

	data, exists := s.splitPackets[splitID]
	if !exists {
		data = &splitPacketData{
			splitCount: splitCount,
			fragments:  make(map[uint32][]byte),
		}
		s.splitPackets[splitID] = data
	}

	data.fragments[splitIndex] = pkt.payload

	if uint32(len(data.fragments)) != data.splitCount {
		logger.Debug("raknet.handleSplitPacket",
			"status", "waiting for more fragments",
			"received", len(data.fragments),
			"total", data.splitCount)
		return nil
	}

	logger.Info("raknet.handleSplitPacket",
		"status", "all fragments received, reassembling",
		"splitID", splitID,
		"fragmentCount", data.splitCount)

	totalSize := 0
	for i := uint32(0); i < data.splitCount; i++ {
		frag, ok := data.fragments[i]
		if !ok {
			logger.Error("raknet.handleSplitPacket", "error", "missing fragment", "index", i)
			delete(s.splitPackets, splitID)
			return nil
		}
		totalSize += len(frag)
	}

	result := make([]byte, 0, totalSize)
	for i := uint32(0); i < data.splitCount; i++ {
		result = append(result, data.fragments[i]...)
	}

	delete(s.splitPackets, splitID)

	logger.Debug("raknet.handleSplitPacket",
		"status", "reassembly complete",
		"totalSize", len(result))

	return result
}

func (s *Session) handleConnectionRequest(data []byte) {
	if len(data) < 17 {
		return
	}

	clientID := binary.BigEndian.Uint64(data[1:9])
	sendPingTime := binary.BigEndian.Uint64(data[9:17])
	logger.Debug("raknet.handleConnectionRequest", "clientID", clientID, "pingTime", sendPingTime)

	buf := new(bytes.Buffer)
	buf.WriteByte(IDConnectionRequestAccepted)

	writeAddress(buf, s.addr)

	binary.Write(buf, binary.BigEndian, uint16(0))

	for i := 0; i < 10; i++ {
		buf.WriteByte(4)
		buf.Write([]byte{127, 0, 0, 1})
		binary.Write(buf, binary.BigEndian, uint16(19132))
	}

	binary.Write(buf, binary.BigEndian, sendPingTime)

	binary.Write(buf, binary.BigEndian, uint64(time.Now().UnixMilli()))

	s.sendReliable(buf.Bytes())
	logger.Debug("raknet.handleConnectionRequest", "sent", "ConnectionRequestAccepted")
}

func (s *Session) handleNewIncomingConnection(data []byte) {
	logger.Debug("raknet.handleNewIncomingConnection", "address", s.addr.String())
	s.connected = true

	if s.server.OnConnect != nil {
		s.server.OnConnect(s)
	}
}

func (s *Session) handleACK(data []byte) {

	logger.Debug("raknet.Session.handleACK", "size", len(data))
}

func (s *Session) handleNAK(data []byte) {

	logger.Debug("raknet.Session.handleNAK", "size", len(data))
}

func (s *Session) sendACK() {
	s.mu.Lock()
	if len(s.ackQueue) == 0 {
		s.mu.Unlock()
		return
	}
	acks := s.ackQueue
	s.ackQueue = nil
	s.mu.Unlock()

	buf := new(bytes.Buffer)
	buf.WriteByte(IDAcknowledge)

	binary.Write(buf, binary.BigEndian, uint16(len(acks)))

	for _, seqNum := range acks {
		buf.WriteByte(1)

		buf.WriteByte(byte(seqNum))
		buf.WriteByte(byte(seqNum >> 8))
		buf.WriteByte(byte(seqNum >> 16))
	}

	s.server.SendTo(s.addr, buf.Bytes())
}

func (s *Session) sendReliable(payload []byte) {
	const maxPayloadSize = 1400

	if len(payload) <= maxPayloadSize {

		s.sendsingleReliable(payload)
		return
	}

	s.mu.Lock()
	splitID := s.splitID
	s.splitID++
	s.mu.Unlock()

	totalLen := len(payload)
	splitCount := uint32((totalLen + maxPayloadSize - 1) / maxPayloadSize)

	logger.Debug("raknet.sendReliable", "action", "splitting packet", "totalLen", totalLen, "splitID", splitID, "count", splitCount)

	for i := uint32(0); i < splitCount; i++ {
		start := int(i) * maxPayloadSize
		end := start + maxPayloadSize
		if end > totalLen {
			end = totalLen
		}
		chunk := payload[start:end]

		s.sendSplitFragment(chunk, splitID, splitCount, i)
	}
}

func (s *Session) sendsingleReliable(payload []byte) {
	s.mu.Lock()
	seqNum := s.sendSeqNum
	s.sendSeqNum++
	msgIndex := s.messageIndex
	s.messageIndex++
	s.mu.Unlock()

	buf := new(bytes.Buffer)

	buf.WriteByte(0x84)

	buf.WriteByte(byte(seqNum))
	buf.WriteByte(byte(seqNum >> 8))
	buf.WriteByte(byte(seqNum >> 16))

	buf.WriteByte(0x40)

	lengthBits := uint16(len(payload) * 8)
	binary.Write(buf, binary.BigEndian, lengthBits)

	buf.WriteByte(byte(msgIndex))
	buf.WriteByte(byte(msgIndex >> 8))
	buf.WriteByte(byte(msgIndex >> 16))

	buf.Write(payload)

	s.server.SendTo(s.addr, buf.Bytes())

}

func (s *Session) sendSplitFragment(chunk []byte, splitID uint16, splitCount uint32, splitIndex uint32) {
	s.mu.Lock()
	seqNum := s.sendSeqNum
	s.sendSeqNum++
	msgIndex := s.messageIndex
	s.messageIndex++
	s.mu.Unlock()

	buf := new(bytes.Buffer)

	buf.WriteByte(0x84)

	buf.WriteByte(byte(seqNum))
	buf.WriteByte(byte(seqNum >> 8))
	buf.WriteByte(byte(seqNum >> 16))

	buf.WriteByte(0x50)

	lengthBits := uint16(len(chunk) * 8)
	binary.Write(buf, binary.BigEndian, lengthBits)

	buf.WriteByte(byte(msgIndex))
	buf.WriteByte(byte(msgIndex >> 8))
	buf.WriteByte(byte(msgIndex >> 16))

	binary.Write(buf, binary.BigEndian, splitCount)
	binary.Write(buf, binary.BigEndian, splitID)
	binary.Write(buf, binary.BigEndian, splitIndex)

	buf.Write(chunk)

	s.server.SendTo(s.addr, buf.Bytes())
}

func (s *Session) SendPacket(data []byte) {

	wrapped := make([]byte, len(data)+1)
	wrapped[0] = 0x8e
	copy(wrapped[1:], data)
	s.sendReliable(wrapped)
}

func (s *Session) Close() {

	s.sendReliable([]byte{IDDisconnectNotification})
	s.server.removeSession(s.addr.String())
}
