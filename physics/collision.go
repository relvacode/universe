package physics

import (
	"github.com/relvacode/universe/internal"
	"math"
)

func Colliding(o1, o2 Object) (float64, bool) {
	delta := internal.Vector{
		X: o2.P.X - o1.P.X,
		Y: o2.P.Y - o1.P.Y,
	}

	dist := math.Sqrt(delta.Dot())

	return dist, dist < o1.R+o2.R
}

// CollisionTime calculates the timestep at which the two objects will collide
func CollisionTime(n1, n2 Object) float64 {
	delta := internal.Vector{
		X: n2.P.X - n1.P.X,
		Y: n2.P.Y - n1.P.Y,
	}

	dist := math.Sqrt(delta.Dot())
	dist = dist - (n1.R + n2.R) // The remaining distance

	// combined velocity
	vec := internal.Vector{
		X: n2.V.X - n1.V.X,
		Y: n2.V.Y - n1.V.Y,
	}

	speed := math.Sqrt(vec.Dot())
	return dist / speed
}

//func (c Collision) Move(o1, o2 *Object) {
//	touchDistFromO1 := c.Distance * (o1.R / (o1.R + o2.R))
//
//	contact := internal.Vector{
//		X: o1.P.X + c.Normal.X*touchDistFromO1,
//		Y: o1.P.Y + c.Normal.Y*touchDistFromO1,
//	}
//
//	o1.P.X = contact.X - c.Normal.X*o1.R
//	o1.P.Y = contact.Y - c.Normal.Y*o1.R
//
//	o2.P.X = contact.X + c.Normal.X*o2.R
//	o2.P.Y = contact.Y + c.Normal.Y*o2.R
//}

func Reflect(o1, o2 *Object) {
	delta := internal.Vector{
		X: o2.P.X - o1.P.X,
		Y: o2.P.Y - o1.P.Y,
	}

	dist := math.Sqrt(delta.Dot())
	normal := internal.Vector{
		X: delta.X / dist,
		Y: delta.Y / dist,
	}

	massTotal := o1.M + o2.M
	massDelta := o1.M - o2.M
	direction := normal.Direction()

	n1 := math.Sqrt(o1.V.Dot())
	n2 := math.Sqrt(o2.V.Dot())

	dir1 := o1.V.Direction()
	dir2 := o2.V.Direction()

	v1s := n1 * math.Sin(dir1-direction)
	cp := math.Cos(direction)
	sp := math.Sin(direction)

	cpd1 := n1 * math.Cos(dir1-direction)
	cpd2 := n2 * math.Cos(dir2-direction)

	cpp := math.Cos(direction + math.Pi/2)
	spp := math.Sin(direction + math.Pi/2)

	t := (cpd1*massDelta + 2*o2.M*cpd2) / massTotal
	o1.V.X = t*cp + v1s*cpp
	o1.V.Y = t*sp + v1s*spp

	direction += math.Pi

	v2s := n2 * math.Sin(dir2-direction)
	cpd1 = n1 * math.Cos(dir1-direction)
	cpd2 = n2 * math.Cos(dir2-direction)

	t = (cpd2*-massDelta + 2*o1.M*cpd1) / massTotal
	o2.V.X = t*-cp + v2s*-cpp
	o2.V.Y = t*-sp + v2s*-spp
}
