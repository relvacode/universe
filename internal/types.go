package internal

import (
	"math"
)

type Vector struct {
	X, Y float64
}

func (v Vector) Equals(o Vector) bool {
	return v.X == o.X && v.Y == o.Y
}

func (v Vector) Add(x, y float64) Vector {
	return Vector{
		X: v.X + x,
		Y: v.Y + y,
	}
}

func (v Vector) Direction() float64 {
	return math.Atan2(v.Y, v.X)
}

func (v Vector) Dot() float64 {
	return v.X*v.X + v.Y*v.Y
}

type BoundingBox struct {
	X, Y, W, H float64
}

func (bb BoundingBox) Abs() BoundingBox {
	if bb.W >= 0 && bb.H >= 0 {
		return bb
	}

	var abs = BoundingBox{
		X: bb.X,
		Y: bb.Y,
		W: math.Abs(bb.W),
		H: math.Abs(bb.H),
	}

	if bb.W < 0 {
		abs.X = bb.X + bb.W
	}
	if bb.H < 0 {
		abs.Y = bb.Y + bb.H
	}

	return abs
}

func (bb BoundingBox) Center() Vector {
	return Vector{
		X: bb.X + (bb.W / 2),
		Y: bb.Y + (bb.H / 2),
	}
}

func (bb BoundingBox) ContainsPoint(v Vector) bool {
	return bb.X < v.X && bb.X+bb.W > v.X &&
		bb.Y < v.Y && bb.Y+bb.H > v.Y
}

func (bb BoundingBox) Intersects(o BoundingBox) bool {
	return (bb.X <= o.X+o.W && bb.X+bb.W >= o.X) &&
		(bb.Y <= o.Y+o.H && bb.Y+bb.H >= o.Y)
}

func (bb BoundingBox) Contains(o BoundingBox) bool {
	return bb.X < o.X && bb.Y < o.Y && bb.X+bb.W > o.X+o.W && bb.Y+bb.H > o.Y+o.H
}
