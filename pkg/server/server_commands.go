package server

import (
	"time"
)

func (s *Server) GetMaxPlayers() int {
	return s.Config.MaxPlayers
}

func (s *Server) GetStartTime() time.Time {
	return s.StartTime
}

func (s *Server) GetRakNetSessionCount() int {
	if s.RakNet != nil {

		return len(s.Players)
	}
	return 0
}
