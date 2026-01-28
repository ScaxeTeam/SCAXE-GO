package world

import (
	"fmt"
	"math"
)

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

func NewVector3(x, y, z float64) *Vector3 {
	return &Vector3{X: x, Y: y, Z: z}
}

func (v *Vector3) Floor() *Vector3 {
	return &Vector3{
		X: math.Floor(v.X),
		Y: math.Floor(v.Y),
		Z: math.Floor(v.Z),
	}
}

func (v *Vector3) String() string {
	return fmt.Sprintf("Vector3(x=%g, y=%g, z=%g)", v.X, v.Y, v.Z)
}
