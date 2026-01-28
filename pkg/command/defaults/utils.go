package defaults

import (
	"strconv"
	"strings"

	"github.com/scaxe/scaxe-go/pkg/command"
	"github.com/scaxe/scaxe-go/pkg/entity"
)

func ParseRelativeCoordinate(val string, relativeTo float64) (float64, error) {
	if strings.HasPrefix(val, "~") {
		offsetStr := val[1:]
		if offsetStr == "" {
			return relativeTo, nil
		}
		offset, err := strconv.ParseFloat(offsetStr, 64)
		if err != nil {
			return 0, err
		}
		return relativeTo + offset, nil
	}
	return strconv.ParseFloat(val, 64)
}

func ParseBlockArg(arg string) (int, int, bool) {
	parts := strings.Split(arg, ":")
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, false
	}
	meta := 0
	if len(parts) > 1 {
		m, err := strconv.Atoi(parts[1])
		if err == nil {
			meta = m
		}
	}
	return id, meta, true
}

type Positional interface {
	GetPosition() *entity.Vector3
}

func parseLocation(sender command.CommandSender, args []string, startIndex int) (float64, float64, float64, bool) {
	if len(args) < startIndex+3 {
		return 0, 0, 0, false
	}

	var px, py, pz float64
	if posSender, ok := sender.(Positional); ok {
		pos := posSender.GetPosition()
		px, py, pz = pos.X, pos.Y, pos.Z
	} else if p, ok := sender.(command.PlayerSender); ok {

		if posSender, ok := p.(Positional); ok {
			pos := posSender.GetPosition()
			px, py, pz = pos.X, pos.Y, pos.Z
		}
	}

	x, err1 := ParseRelativeCoordinate(args[startIndex], px)
	y, err2 := ParseRelativeCoordinate(args[startIndex+1], py)
	z, err3 := ParseRelativeCoordinate(args[startIndex+2], pz)

	if err1 != nil || err2 != nil || err3 != nil {
		return 0, 0, 0, false
	}
	return x, y, z, true
}
