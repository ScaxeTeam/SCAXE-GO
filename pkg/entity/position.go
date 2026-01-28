package entity

import (
	"math"
)

const (
	SideDown  = 0
	SideUp    = 1
	SideNorth = 2
	SideSouth = 3
	SideWest  = 4
	SideEast  = 5
)

type Vector3 struct {
	X, Y, Z float64
}

func NewVector3(x, y, z float64) *Vector3 {
	return &Vector3{X: x, Y: y, Z: z}
}

func (v *Vector3) Zero() *Vector3 {
	return &Vector3{0, 0, 0}
}

func (v *Vector3) Add(x, y, z float64) *Vector3 {
	return &Vector3{v.X + x, v.Y + y, v.Z + z}
}

func (v *Vector3) AddVector(other *Vector3) *Vector3 {
	return v.Add(other.X, other.Y, other.Z)
}

func (v *Vector3) Subtract(x, y, z float64) *Vector3 {
	return &Vector3{v.X - x, v.Y - y, v.Z - z}
}

func (v *Vector3) SubtractVector(other *Vector3) *Vector3 {
	return v.Subtract(other.X, other.Y, other.Z)
}

func (v *Vector3) Multiply(number float64) *Vector3 {
	return &Vector3{v.X * number, v.Y * number, v.Z * number}
}

func (v *Vector3) Divide(number float64) *Vector3 {
	return &Vector3{v.X / number, v.Y / number, v.Z / number}
}

func (v *Vector3) FloorX() int {
	return int(math.Floor(v.X))
}

func (v *Vector3) FloorY() int {
	return int(math.Floor(v.Y))
}

func (v *Vector3) FloorZ() int {
	return int(math.Floor(v.Z))
}

func (v *Vector3) Floor() *Vector3 {
	return &Vector3{
		X: float64(v.FloorX()),
		Y: float64(v.FloorY()),
		Z: float64(v.FloorZ()),
	}
}

func (v *Vector3) Ceil() *Vector3 {
	return &Vector3{
		X: math.Ceil(v.X),
		Y: math.Ceil(v.Y),
		Z: math.Ceil(v.Z),
	}
}

func (v *Vector3) Round() *Vector3 {
	return &Vector3{
		X: math.Round(v.X),
		Y: math.Round(v.Y),
		Z: math.Round(v.Z),
	}
}

func (v *Vector3) Abs() *Vector3 {
	return &Vector3{
		X: math.Abs(v.X),
		Y: math.Abs(v.Y),
		Z: math.Abs(v.Z),
	}
}

func (v *Vector3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

func (v *Vector3) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v *Vector3) Normalize() *Vector3 {
	len := v.LengthSquared()
	if len > 0 {
		return v.Divide(math.Sqrt(len))
	}
	return &Vector3{0, 0, 0}
}

func (v *Vector3) Dot(other *Vector3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (v *Vector3) Cross(other *Vector3) *Vector3 {
	return &Vector3{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

func (v *Vector3) Distance(pos *Vector3) float64 {
	return math.Sqrt(v.DistanceSquared(pos))
}

func (v *Vector3) DistanceSquared(pos *Vector3) float64 {
	dx := v.X - pos.X
	dy := v.Y - pos.Y
	dz := v.Z - pos.Z
	return dx*dx + dy*dy + dz*dz
}

func (v *Vector3) Equals(other *Vector3) bool {
	return v.X == other.X && v.Y == other.Y && v.Z == other.Z
}

func (v *Vector3) GetSide(side int, step int) *Vector3 {
	switch side {
	case SideDown:
		return &Vector3{v.X, v.Y - float64(step), v.Z}
	case SideUp:
		return &Vector3{v.X, v.Y + float64(step), v.Z}
	case SideNorth:
		return &Vector3{v.X, v.Y, v.Z - float64(step)}
	case SideSouth:
		return &Vector3{v.X, v.Y, v.Z + float64(step)}
	case SideWest:
		return &Vector3{v.X - float64(step), v.Y, v.Z}
	case SideEast:
		return &Vector3{v.X + float64(step), v.Y, v.Z}
	default:
		return v
	}
}

func (v *Vector3) SetComponents(x, y, z float64) *Vector3 {
	v.X = x
	v.Y = y
	v.Z = z
	return v
}

func (v *Vector3) Clone() *Vector3 {
	return &Vector3{v.X, v.Y, v.Z}
}

type Location struct {
	*Vector3
	Yaw   float32
	Pitch float32
}

func NewLocation(x, y, z float64, yaw, pitch float32) *Location {
	return &Location{
		Vector3: NewVector3(x, y, z),
		Yaw:     yaw,
		Pitch:   pitch,
	}
}
