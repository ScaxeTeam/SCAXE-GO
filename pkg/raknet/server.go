package raknet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/scaxe/scaxe-go/pkg/logger"
)

const (
	IDUnconnectedPing         byte = 0x01
	IDUnconnectedPingOpenConn byte = 0x02
	IDUnconnectedPong         byte = 0x1c
	IDOpenConnectionRequest1  byte = 0x05
	IDOpenConnectionReply1    byte = 0x06
	IDOpenConnectionRequest2  byte = 0x07
	IDOpenConnectionReply2    byte = 0x08

	IDConnectionRequest         byte = 0x09
	IDConnectionRequestAccepted byte = 0x10
	IDNewIncomingConnection     byte = 0x13
	IDDisconnectNotification    byte = 0x15
	IDIncompatibleProtocol      byte = 0x19
	IDAcknowledge               byte = 0xc0
	IDNAcknowledge              byte = 0xa0

	IDCustom0 byte = 0x80
	IDCustomF byte = 0x8f
)

var RakNetMagic = []byte{
	0x00, 0xff, 0xff, 0x00,
	0xfe, 0xfe, 0xfe, 0xfe,
	0xfd, 0xfd, 0xfd, 0xfd,
	0x12, 0x34, 0x56, 0x78,
}

var SupportedProtocols = []byte{7, 8}

type Server struct {
	conn     *net.UDPConn
	address  string
	serverID int64
	pongData []byte

	sessions   map[string]*Session
	sessionsMu sync.RWMutex

	OnConnect    func(*Session)
	OnDisconnect func(*Session)
	OnPacket     func(*Session, []byte)

	running bool
	stopCh  chan struct{}
}

func NewServer(address string) *Server {
	logger.Debug("raknet.NewServer", "address", address)
	return &Server{
		address:  address,
		serverID: rand.Int63(),
		sessions: make(map[string]*Session),
		stopCh:   make(chan struct{}),
	}
}

func (s *Server) ServerID() int64 {
	return s.serverID
}

func (s *Server) SetPongData(data []byte) {
	s.pongData = data
	logger.Debug("raknet.SetPongData", "size", len(data))
}

func (s *Server) Start() error {
	logger.Debug("raknet.Server.Start", "address", s.address)

	addr, err := net.ResolveUDPAddr("udp", s.address)
	if err != nil {
		logger.Error("raknet.Server.Start", "error", "failed to resolve address", "err", err)
		return err
	}

	s.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		logger.Error("raknet.Server.Start", "error", "failed to listen", "err", err)
		return err
	}

	s.running = true

	go s.readLoop()

	logger.Info("raknet.Server.Start", "status", "listening", "address", s.address)
	return nil
}

func (s *Server) Stop() {
	logger.Debug("raknet.Server.Stop", "action", "stopping")
	s.running = false
	close(s.stopCh)
	if s.conn != nil {
		s.conn.Close()
	}
	logger.Info("raknet.Server.Stop", "status", "stopped")
}

func (s *Server) readLoop() {
	buf := make([]byte, 2048)

	for s.running {
		s.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			if !s.running {
				return
			}
			logger.Error("raknet.readLoop", "error", err)
			continue
		}

		if n == 0 {
			continue
		}

		data := make([]byte, n)
		copy(data, buf[:n])

		s.handlePacket(addr, data)
	}
}

func (s *Server) handlePacket(addr *net.UDPAddr, data []byte) {
	if len(data) == 0 {
		return
	}

	packetID := data[0]
	addrStr := addr.String()

	logger.Debug("raknet.handlePacket", "from", addrStr, "packetID", fmt.Sprintf("0x%02x", packetID), "size", len(data))

	switch packetID {
	case IDUnconnectedPing, IDUnconnectedPingOpenConn:
		s.handleUnconnectedPing(addr, data)
		return

	case IDOpenConnectionRequest1:
		s.handleOpenConnectionRequest1(addr, data)
		return

	case IDOpenConnectionRequest2:
		s.handleOpenConnectionRequest2(addr, data)
		return
	}

	s.sessionsMu.RLock()
	session, exists := s.sessions[addrStr]
	s.sessionsMu.RUnlock()

	if !exists {
		logger.Warn("raknet.handlePacket", "warning", "packet from unknown session", "address", addrStr)
		return
	}

	if packetID >= IDCustom0 && packetID <= IDCustomF {
		session.handleDataPacket(data)
		return
	}

	if packetID == IDAcknowledge || (packetID&0xf0) == 0xc0 {
		session.handleACK(data)
		return
	}

	if packetID == IDNAcknowledge || (packetID&0xf0) == 0xa0 {
		session.handleNAK(data)
		return
	}

	if packetID == IDConnectionRequest {
		session.handleConnectionRequest(data)
		return
	}

	if packetID == IDNewIncomingConnection {
		session.handleNewIncomingConnection(data)
		return
	}

	if packetID == IDDisconnectNotification {
		s.removeSession(addrStr)
		return
	}

	logger.Warn("raknet.handlePacket", "warning", "unhandled packet", "id", fmt.Sprintf("0x%02x", packetID))
}

func (s *Server) handleUnconnectedPing(addr *net.UDPAddr, data []byte) {
	if len(data) < 25 {
		return
	}

	pingTime := binary.BigEndian.Uint64(data[1:9])
	logger.Debug("raknet.handleUnconnectedPing", "from", addr.String(), "pingTime", pingTime)

	buf := new(bytes.Buffer)
	buf.WriteByte(IDUnconnectedPong)
	binary.Write(buf, binary.BigEndian, pingTime)
	binary.Write(buf, binary.BigEndian, s.serverID)
	buf.Write(RakNetMagic)

	binary.Write(buf, binary.BigEndian, uint16(len(s.pongData)))
	buf.Write(s.pongData)

	s.conn.WriteToUDP(buf.Bytes(), addr)
	logger.Debug("raknet.handleUnconnectedPing", "sent", "pong", "size", buf.Len())
}

func (s *Server) handleOpenConnectionRequest1(addr *net.UDPAddr, data []byte) {
	if len(data) < 18 {
		return
	}

	if !bytes.Equal(data[1:17], RakNetMagic) {
		logger.Warn("raknet.handleOpenConnectionRequest1", "warning", "invalid magic")
		return
	}

	protocolVersion := data[17]
	logger.Debug("raknet.handleOpenConnectionRequest1", "from", addr.String(), "protocolVersion", protocolVersion)

	supported := false
	for _, p := range SupportedProtocols {
		if p == protocolVersion {
			supported = true
			break
		}
	}

	if !supported {
		logger.Warn("raknet.handleOpenConnectionRequest1", "warning", "unsupported protocol", "version", protocolVersion)

		buf := new(bytes.Buffer)
		buf.WriteByte(IDIncompatibleProtocol)
		buf.WriteByte(SupportedProtocols[0])
		buf.Write(RakNetMagic)
		binary.Write(buf, binary.BigEndian, s.serverID)
		s.conn.WriteToUDP(buf.Bytes(), addr)
		return
	}

	mtuSize := len(data) + 28

	buf := new(bytes.Buffer)
	buf.WriteByte(IDOpenConnectionReply1)
	buf.Write(RakNetMagic)
	binary.Write(buf, binary.BigEndian, s.serverID)
	buf.WriteByte(0)
	binary.Write(buf, binary.BigEndian, uint16(mtuSize))

	s.conn.WriteToUDP(buf.Bytes(), addr)
	logger.Debug("raknet.handleOpenConnectionRequest1", "sent", "reply1", "mtu", mtuSize)
}

func (s *Server) handleOpenConnectionRequest2(addr *net.UDPAddr, data []byte) {

	if len(data) < 34 {
		logger.Warn("raknet.handleOpenConnectionRequest2", "warning", "packet too small", "size", len(data))
		return
	}

	logger.Debug("raknet.handleOpenConnectionRequest2", "rawDataSize", len(data))

	if !bytes.Equal(data[1:17], RakNetMagic) {
		logger.Warn("raknet.handleOpenConnectionRequest2", "warning", "invalid magic")
		return
	}

	offset := 17
	addrVersion := data[offset]
	offset++

	if addrVersion == 4 {

		offset += 4 + 2
	} else if addrVersion == 6 {

		offset += 16 + 2
	} else {
		logger.Warn("raknet.handleOpenConnectionRequest2", "warning", "unknown address version", "version", addrVersion)

		offset += 6
	}

	if offset+2 > len(data) {
		logger.Warn("raknet.handleOpenConnectionRequest2", "warning", "not enough data for MTU")

		offset = len(data) - 10
	}
	mtu := binary.BigEndian.Uint16(data[offset : offset+2])
	offset += 2

	if offset+8 > len(data) {
		logger.Warn("raknet.handleOpenConnectionRequest2", "warning", "not enough data for GUID")
		offset = len(data) - 8
	}
	clientGuid := binary.BigEndian.Uint64(data[offset : offset+8])

	logger.Debug("raknet.handleOpenConnectionRequest2", "from", addr.String(), "mtu", mtu, "clientGuid", clientGuid)

	session := NewSession(s, addr, mtu, clientGuid)

	s.sessionsMu.Lock()
	s.sessions[addr.String()] = session
	s.sessionsMu.Unlock()

	buf := new(bytes.Buffer)
	buf.WriteByte(IDOpenConnectionReply2)
	buf.Write(RakNetMagic)
	binary.Write(buf, binary.BigEndian, s.serverID)

	writeAddress(buf, addr)

	binary.Write(buf, binary.BigEndian, mtu)
	buf.WriteByte(0)

	s.conn.WriteToUDP(buf.Bytes(), addr)
	logger.Info("raknet.handleOpenConnectionRequest2", "sent", "reply2", "session", addr.String(), "mtu", mtu)
}

func (s *Server) removeSession(addrStr string) {
	s.sessionsMu.Lock()
	session, exists := s.sessions[addrStr]
	if exists {
		delete(s.sessions, addrStr)
	}
	s.sessionsMu.Unlock()

	if exists && s.OnDisconnect != nil {
		s.OnDisconnect(session)
	}
	logger.Debug("raknet.removeSession", "address", addrStr)
}

func (s *Server) SendTo(addr *net.UDPAddr, data []byte) error {
	_, err := s.conn.WriteToUDP(data, addr)
	return err
}

func writeAddress(buf *bytes.Buffer, addr *net.UDPAddr) {
	ip := addr.IP.To4()
	if ip == nil {

		buf.WriteByte(6)

		buf.Write(make([]byte, 16))
		binary.Write(buf, binary.BigEndian, uint16(addr.Port))
	} else {

		buf.WriteByte(4)

		buf.WriteByte(^ip[0])
		buf.WriteByte(^ip[1])
		buf.WriteByte(^ip[2])
		buf.WriteByte(^ip[3])
		binary.Write(buf, binary.BigEndian, uint16(addr.Port))
	}
}
