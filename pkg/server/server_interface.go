package server

import (
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

func (s *Server) BroadcastMessage(message string) {
	textPk := protocol.NewTextPacket()
	textPk.TextType = protocol.TextTypeRaw
	textPk.Message = message
	s.BroadcastPacket(textPk)
}

func (s *Server) GetTime() int {
	if s.Level != nil {
		return int(s.Level.GetTime())
	}
	return 0
}

func (s *Server) SetTime(t int) {
	if s.Level != nil {
		s.Level.SetTime(int64(t))

		timePk := protocol.NewSetTimePacket()
		timePk.Time = int32(t)
		timePk.Started = true
		s.BroadcastPacket(timePk)
	}
}

func (s *Server) GetAverageTPS() float64 {
	return s.GetTPS()
}

func (s *Server) GetDifficulty() int {
	return s.Config.Difficulty
}

func (s *Server) SetDifficulty(difficulty int) {
	s.Config.Difficulty = difficulty

	diffPk := protocol.NewSetDifficultyPacket()
	diffPk.Difficulty = int32(difficulty)
	s.BroadcastPacket(diffPk)
}

func (s *Server) GetSeed() int64 {
	if s.Level != nil {
		return s.Level.GetSeed()
	}
	return 0
}
