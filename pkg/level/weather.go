package level

import (
	"math/rand/v2"

	"github.com/scaxe/scaxe-go/pkg/protocol"
)

const (
	WeatherClear   = 0
	WeatherRain    = 1
	WeatherThunder = 2
)

func (l *Level) IsRaining() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.Raining
}

func (l *Level) IsThundering() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.Thundering
}

func (l *Level) GetWeather() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.Thundering {
		return WeatherThunder
	}
	if l.Raining {
		return WeatherRain
	}
	return WeatherClear
}

func (l *Level) SetWeather(w int) {
	l.mu.Lock()
	switch w {
	case WeatherClear:
		l.Raining = false
		l.Thundering = false
		l.RainTime = 12000 + rand.IntN(168000)
	case WeatherRain:
		l.Raining = true
		l.Thundering = false
		l.RainTime = 12000 + rand.IntN(12000)
		l.ThunderTime = rand.IntN(12000)
	case WeatherThunder:
		l.Raining = true
		l.Thundering = true
		l.RainTime = 12000 + rand.IntN(12000)
		l.ThunderTime = 3600 + rand.IntN(12000)
	}
	l.mu.Unlock()
}

func (l *Level) TickWeather() {
	l.mu.Lock()
	if l.RainTime > 0 {
		l.RainTime--
		if l.RainTime <= 0 {
			l.Raining = !l.Raining
			if l.Raining {
				l.RainTime = 12000 + rand.IntN(12000)
			} else {
				l.RainTime = 12000 + rand.IntN(168000)
				l.Thundering = false
			}
		}
	}
	if l.Raining && l.ThunderTime > 0 {
		l.ThunderTime--
		if l.ThunderTime <= 0 {
			l.Thundering = !l.Thundering
			if l.Thundering {
				l.ThunderTime = 3600 + rand.IntN(12000)
			} else {
				l.ThunderTime = rand.IntN(12000)
			}
		}
	}
	l.mu.Unlock()
}

func (l *Level) MakeWeatherPackets() []protocol.DataPacket {
	l.mu.RLock()
	raining := l.Raining
	thundering := l.Thundering
	l.mu.RUnlock()

	var pkts []protocol.DataPacket

	if raining {
		pk := protocol.NewLevelEventPacket()
		pk.EventID = uint16(protocol.EventStartRain)
		pk.Data = 65535
		pkts = append(pkts, pk)
	} else {
		pk := protocol.NewLevelEventPacket()
		pk.EventID = uint16(protocol.EventStopRain)
		pkts = append(pkts, pk)
	}

	if thundering {
		pk := protocol.NewLevelEventPacket()
		pk.EventID = uint16(protocol.EventStartThunder)
		pk.Data = 65535
		pkts = append(pkts, pk)
	} else {
		pk := protocol.NewLevelEventPacket()
		pk.EventID = uint16(protocol.EventStopThunder)
		pkts = append(pkts, pk)
	}

	return pkts
}
