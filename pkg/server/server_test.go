package server

import (
	"net"
	"testing"

	"github.com/google/uuid"

	"github.com/scaxe/scaxe-go/pkg/config"
	"github.com/scaxe/scaxe-go/pkg/player"
)

type MockSession struct {
	addr net.Addr
}

func (m *MockSession) ID() int64           { return 123 }
func (m *MockSession) Address() net.Addr   { return m.addr }
func (m *MockSession) SendPacket(b []byte) {}
func (m *MockSession) Close()              {}

func TestServer_GetOnlineCount(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.MaxPlayers = 10
	s := NewServer(cfg)

	if s.GetOnlineCount() != 0 {
		t.Errorf("Expected 0 players, got %d", s.GetOnlineCount())
	}

	addr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345}

	p := player.NewPlayer(nil, addr.String(), addr.Port)

	p.HandleLogin("Steve", uuid.New().String(), "skin", []byte("data"), 60)

	s.mu.Lock()
	s.PlayersByName["Steve"] = p
	s.mu.Unlock()

	if s.GetOnlineCount() != 1 {
		t.Errorf("Expected 1 player, got %d", s.GetOnlineCount())
	}

	p2 := player.NewPlayer(nil, "127.0.0.2", 12346)
	p2.HandleLogin("Steve", uuid.New().String(), "skin", []byte("data"), 60)

	s.mu.Lock()
	s.PlayersByName["Steve"] = p2
	s.mu.Unlock()

	if s.GetOnlineCount() != 1 {
		t.Errorf("Expected 1 player after overwrite, got %d", s.GetOnlineCount())
	}

	s.mu.Lock()
	delete(s.PlayersByName, "Steve")
	s.mu.Unlock()

	if s.GetOnlineCount() != 0 {
		t.Errorf("Expected 0 players after logout, got %d", s.GetOnlineCount())
	}
}
