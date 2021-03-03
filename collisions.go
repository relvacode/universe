package universe

import (
	"github.com/relvacode/universe/internal"
	"github.com/relvacode/universe/physics"
	"math"
)

type CollisionResolver func(e1, e2 *Entity, dist float64)

func MassReflectCollisionResolver(e1, e2 *Entity, _ float64) {
	timeColliding := physics.CollisionTime(e1.Object, e2.Object) * .5
	if timeColliding == 0 {
		// No adjustments need to be made other than reflect the velocities
		physics.Reflect(&e1.Object, &e2.Object)
		return
	}

	// Undo half the time each object was colliding
	e1.Step(timeColliding)
	e2.Step(timeColliding)

	// Reflect the velocities
	physics.Reflect(&e1.Object, &e2.Object)

	// Step each object in the other direction for the other half of the time
	e1.Step(-timeColliding)
	e1.Step(-timeColliding)
}

func NoopCollisionResolver(_, _ *Entity, _ float64) {}

func areaOverlap(r1, r2, distance float64) float64 {
	var rr1 = r1 * r1
	var rr2 = r2 * r2

	var phi = (math.Acos((rr1 + (distance * distance) - rr2) / (2 * r1 * distance))) * 2
	var theta = (math.Acos((rr2 + (distance * distance) - rr1) / (2 * r2 * distance))) * 2
	var area1 = 0.5*theta*rr2 - 0.5*rr2*math.Sin(theta)
	var area2 = 0.5*phi*rr1 - 0.5*rr1*math.Sin(phi)
	return area1 + area2
}

func AbsorbCollisionResolver(e1, e2 *Entity, distance float64) {
	var consumer = e1
	var consumed = e2
	if e2.M > e1.M {
		consumer, consumed = consumed, consumer
	}

	consumedArea := math.Pi * math.Pow(consumed.R, 2)
	overlapArea := areaOverlap(consumer.R, consumed.R, distance)

	if overlapArea/consumedArea < .8 {
		return
	}

	consumed.Disabled = true

	massVelocity := internal.Vector{
		X: (consumer.M * consumer.V.X) + (consumed.M * consumed.V.X),
		Y: (consumer.M * consumer.V.Y) + (consumed.M * consumed.V.Y),
	}

	consumer.V.X = massVelocity.X / consumer.M
	consumer.V.Y = massVelocity.Y / consumer.M

	consumer.M = consumer.M + consumed.M
	consumer.R = math.Cbrt(consumer.M)

}
