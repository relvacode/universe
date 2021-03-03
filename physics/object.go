package physics

import (
	"github.com/relvacode/universe/internal"
	"math"
)

type Object struct {
	P internal.Vector
	V internal.Vector
	R float64
	M float64
}

func (o *Object) KineticEnergy() float64 {
	return 0.5 * o.M * math.Pow(math.Abs(o.V.X)+math.Abs(o.V.Y), 2)
}

func (o *Object) ReflectBounds(bb internal.BoundingBox) {
	var (
		top    = bb.Y
		left   = bb.X
		right  = bb.X + bb.W
		bottom = bb.Y + bb.H
	)
	switch {
	case o.P.Y < top+o.R:
		o.P.Y = top + o.R
		o.V.Y *= -1
	case o.P.X > right-o.R:
		o.P.X = right - o.R
		o.V.X *= -1
	case o.P.Y > bottom-o.R:
		o.P.Y = bottom - o.R
		o.V.Y *= -1
	case o.P.X < left+o.R:
		o.P.X = left + o.R
		o.V.X *= -1
	}
}

func (o *Object) BoundingBox() internal.BoundingBox {
	return internal.BoundingBox{
		X: o.P.X - o.R,
		Y: o.P.Y - o.R,
		W: o.R * 2,
		H: o.R * 2,
	}
}

func (o *Object) Step(timestep float64) {
	o.P.X = o.P.X + (o.V.X * timestep)
	o.P.Y = o.P.Y + (o.V.Y * timestep)
}
