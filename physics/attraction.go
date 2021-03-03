package physics

import (
	"github.com/relvacode/universe/internal"
	"math"
)

const (
	G = 8.8e-1
)

func AttractionForceVector(p1, p2 internal.Vector, m1, m2 float64, force float64) internal.Vector {
	direction := internal.Vector{
		X: p2.X - p1.X,
		Y: p2.Y - p1.Y,
	}

	distance := math.Sqrt(direction.Dot())
	if distance == 0 {
		return internal.Vector{}
	}

	normal := internal.Vector{
		X: direction.X / distance,
		Y: direction.Y / distance,
	}

	attraction := force * ((m1 * m2) / distance)

	return internal.Vector{
		X: attraction * normal.X,
		Y: attraction * normal.Y,
	}
}

func AttractionForceImpulse(p1, p2 internal.Vector, m1, m2 float64, timestep float64) internal.Vector {
	force := AttractionForceVector(p1, p2, m1, m2, G)
	return internal.Vector{
		X: force.X * timestep,
		Y: force.Y * timestep,
	}
}

//func Attract(o1, o2 *Object, timestep float64) {
//	impulse := AttractionForceImpulse(o1.P, o2.P, o1.M, o2.M, timestep)
//
//	o1.V.X += impulse.X / o1.M
//	o1.V.Y += impulse.Y / o1.M
//
//	o2.V.X += -impulse.X / o2.M
//	o2.V.Y += -impulse.Y / o2.M
//}
