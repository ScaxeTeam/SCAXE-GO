package network

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/sandertv/go-raknet"
	"github.com/scaxe/scaxe-go/pkg/logger"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

type Server struct {
	listener *raknet.Listener
	sessions map[string]*Session
	mu       sync.RWMutex

	Address    string
	ServerName string
	MaxPlayers int
	MOTD       string

	OnSessionConnect    func(*Session)
	OnSessionDisconnect func(*Session)
	OnPacket            func(*Session, protocol.DataPacket)

	ctx    context.Context
	cancel context.CancelFunc
}

func NewServer(address, serverName string, maxPlayers int) *Server {
	logger.Debug("NewServer", "address", address, "serverName", serverName, "maxPlayers", maxPlayers)

	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		sessions:   make(map[string]*Session),
		Address:    address,
		ServerName: serverName,
		MaxPlayers: maxPlayers,
		MOTD:       serverName,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (s *Server) Start() error {
	logger.Debug("Server.Start", "address", s.Address)

	var err error
	s.listener, err = raknet.Listen(s.Address)
	if err != nil {
		logger.Error("Server.Start", "error", "failed to start listener", "err", err)
		return fmt.Errorf("failed to listen on %s: %w", s.Address, err)
	}

	s.listener.PongData([]byte(s.buildMOTD()))

	logger.Info("Server.Start", "status", "listening", "address", s.Address)

	go s.acceptLoop()

	return nil
}

func (s *Server) Stop() {
	logger.Debug("Server.Stop", "action", "stopping server")

	s.cancel()

	if s.listener != nil {
		s.listener.Close()
	}

	s.mu.Lock()
	for addr, session := range s.sessions {
		session.Close("Server shutdown")
		delete(s.sessions, addr)
	}
	s.mu.Unlock()

	logger.Info("Server.Stop", "status", "stopped")
}

func (s *Server) acceptLoop() {
	logger.Debug("Server.acceptLoop", "status", "started")

	for {
		select {
		case <-s.ctx.Done():
			logger.Debug("Server.acceptLoop", "status", "context cancelled, exiting")
			return
		default:
		}

		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.ctx.Done():
				return
			default:
				logger.Error("Server.acceptLoop", "error", "accept failed", "err", err)
				continue
			}
		}

		session := NewSession(conn, s)
		addr := conn.RemoteAddr().String()

		s.mu.Lock()
		s.sessions[addr] = session
		s.mu.Unlock()

		logger.Info("Server.acceptLoop", "event", "new connection", "address", addr)

		if s.OnSessionConnect != nil {
			s.OnSessionConnect(session)
		}

		go s.handleSession(session)
	}
}

func (s *Server) handleSession(session *Session) {
	addr := session.conn.RemoteAddr().String()
	logger.Debug("Server.handleSession", "address", addr, "status", "started")

	defer func() {
		s.mu.Lock()
		delete(s.sessions, addr)
		s.mu.Unlock()

		if s.OnSessionDisconnect != nil {
			s.OnSessionDisconnect(session)
		}

		logger.Info("Server.handleSession", "address", addr, "status", "disconnected")
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-session.ctx.Done():
			return
		default:
		}

		data, err := session.ReadPacket()
		if err != nil {
			if err.Error() != "EOF" {
				logger.Error("Server.handleSession", "address", addr, "error", err)
			}
			return
		}

		if len(data) == 0 {
			continue
		}

		packets, err := session.ProcessRawPacket(data)
		if err != nil {
			logger.Error("Server.handleSession", "error", "failed to process packet", "err", err)
			continue
		}

		for _, pkt := range packets {

			if pkt.ID() == protocol.IDLogin {
				loginPkt := pkt.(*protocol.LoginPacket)
				logger.Info("Server.handleSession", "event", "login_request", "username", loginPkt.Username, "protocol", loginPkt.Protocol)

				statusPx := protocol.NewPlayStatusPacket()
				statusPx.Status = protocol.PlayStatusLoginSuccess
				if err := session.SendPacket(statusPx); err != nil {
					logger.Error("HandleLogin", "stage", "login_success", "error", err)
					return
				}

				startGame := protocol.NewStartGamePacket()
				startGame.EntityID = 1
				startGame.RuntimeID = 1
				startGame.SpawnX = 128
				startGame.SpawnY = 70
				startGame.SpawnZ = 128
				startGame.Seed = 12345
				startGame.Dimension = 0
				startGame.Generator = 1
				startGame.Gamemode = 1
				startGame.LevelID = "world"
				startGame.X = 128.0
				startGame.Y = 70.0
				startGame.Z = 128.0
				if err := session.SendPacket(startGame); err != nil {
					logger.Error("HandleLogin", "stage", "start_game", "error", err)
					return
				}

				spawnStatus := protocol.NewPlayStatusPacket()
				spawnStatus.Status = protocol.PlayStatusPlayerSpawn
				if err := session.SendPacket(spawnStatus); err != nil {
					logger.Error("HandleLogin", "stage", "player_spawn", "error", err)
					return
				}

				logger.Info("Server.handleSession", "event", "login_complete", "username", loginPkt.Username)
			}

			logger.Debug("Server.handleSession", "packetID", pkt.ID(), "packetName", pkt.Name())
			if s.OnPacket != nil {
				s.OnPacket(session, pkt)
			}
		}
	}
}

func (s *Server) GetSession(addr string) *Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions[addr]
}

func (s *Server) GetSessionCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.sessions)
}

func (s *Server) Broadcast(pkt protocol.DataPacket) {
	logger.Debug("Server.Broadcast", "packetID", pkt.ID(), "packetName", pkt.Name())

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, session := range s.sessions {
		if err := session.SendPacket(pkt); err != nil {
			logger.Error("Server.Broadcast", "error", err, "address", session.conn.RemoteAddr().String())
		}
	}
}

func (s *Server) buildMOTD() string {

	motd := fmt.Sprintf("MCPE;%s;%d;%s;%d;%d;0;%s;Survival;0",
		s.MOTD,
		70,
		"0.14.3",
		s.GetSessionCount(),
		s.MaxPlayers,
		"SCAXE-GO",
	)
	logger.Debug("Server.buildMOTD", "motd", motd)
	return motd
}

type Session struct {
	conn   net.Conn
	server *Server
	mu     sync.Mutex

	Username string
	UUID     string
	Protocol int32

	ctx    context.Context
	cancel context.CancelFunc
}

func NewSession(conn net.Conn, server *Server) *Session {
	ctx, cancel := context.WithCancel(server.ctx)

	session := &Session{
		conn:   conn,
		server: server,
		ctx:    ctx,
		cancel: cancel,
	}

	logger.Debug("NewSession", "address", conn.RemoteAddr().String())
	return session
}

func (s *Session) RemoteAddr() string {
	return s.conn.RemoteAddr().String()
}

func (s *Session) ReadPacket() ([]byte, error) {
	buf := make([]byte, 4096)
	n, err := s.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	logger.Debug("Session.ReadPacket", "size", n)
	return buf[:n], nil
}

func (s *Session) ProcessRawPacket(data []byte) ([]protocol.DataPacket, error) {
	logger.Debug("Session.ProcessRawPacket", "size", len(data))

	if len(data) == 0 {
		return nil, nil
	}

	packetID := data[0]
	logger.Debug("Session.ProcessRawPacket", "packetID", packetID)

	if packetID == protocol.IDBatch {
		batch := protocol.NewBatchPacket()
		stream := protocol.NewBinaryStreamFromBytes(data[1:])
		if err := batch.Decode(stream); err != nil {
			return nil, fmt.Errorf("failed to decode batch: %w", err)
		}

		decompressed, err := batch.Decompress()
		if err != nil {
			return nil, fmt.Errorf("failed to decompress batch: %w", err)
		}

		return protocol.DecodePackets(decompressed)
	}

	pkt := protocol.GetPacket(packetID)
	if pkt == nil {
		logger.Warn("Session.ProcessRawPacket", "warning", "unknown packet", "id", packetID)
		return nil, nil
	}

	stream := protocol.NewBinaryStreamFromBytes(data[1:])
	if err := pkt.Decode(stream); err != nil {
		return nil, fmt.Errorf("failed to decode packet %d: %w", packetID, err)
	}

	return []protocol.DataPacket{pkt}, nil
}

func (s *Session) SendPacket(pkt protocol.DataPacket) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	logger.Debug("Session.SendPacket", "packetID", pkt.ID(), "packetName", pkt.Name())

	stream := protocol.NewBinaryStream()
	if err := pkt.Encode(stream); err != nil {
		return fmt.Errorf("failed to encode packet: %w", err)
	}

	packetData := stream.Bytes()

	batchStream := protocol.NewBinaryStream()
	batchStream.WriteInt(int32(len(packetData)))
	batchStream.WriteBytes(packetData)

	batch := protocol.NewBatchPacket()
	if err := batch.Compress(batchStream.Bytes()); err != nil {
		return fmt.Errorf("failed to compress batch: %w", err)
	}

	finalStream := protocol.NewBinaryStream()
	if err := batch.Encode(finalStream); err != nil {
		return fmt.Errorf("failed to encode batch packet: %w", err)
	}

	_, err := s.conn.Write(finalStream.Bytes())
	if err != nil {
		logger.Error("Session.SendPacket", "error", err)
		return err
	}

	logger.Debug("Session.SendPacket", "sent", len(finalStream.Bytes()), "bytes")
	return nil
}

func (s *Session) SendPacketDirect(pkt protocol.DataPacket) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	stream := protocol.NewBinaryStream()
	if err := pkt.Encode(stream); err != nil {
		return err
	}

	_, err := s.conn.Write(stream.Bytes())
	return err
}

func (s *Session) Close(reason string) {
	logger.Debug("Session.Close", "address", s.RemoteAddr(), "reason", reason)

	disconnect := protocol.NewDisconnectPacket()
	disconnect.Message = reason
	s.SendPacket(disconnect)

	s.cancel()
	s.conn.Close()
}

func (s *Session) Kick(message string) {
	logger.Info("Session.Kick", "address", s.RemoteAddr(), "message", message)
	s.Close(message)
}
