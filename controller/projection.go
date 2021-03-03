package controller

import (
	"github.com/relvacode/universe"
	"github.com/relvacode/universe/draw"
	"github.com/relvacode/universe/internal"
	"math"
)

func (v *View) drawForwardProjections(ctx *draw.Context) {
	if len(*v.entities) == 0 {
		return
	}

	copied := make(universe.EntityList, len(*v.entities))
	path := make(map[*universe.Entity][]internal.Vector, len(*v.entities))

	for i, e := range *v.entities {
		n := &universe.Entity{
			Object: e.Object,
		}
		copied[i] = n
		path[n] = []internal.Vector{n.P}
	}

	iterations := 6 / universe.PhysicsConstantTimestep
	if math.IsInf(iterations, 0) {
		return
	}

	simulation := universe.NewSimulation(v.box)

	var it float64
	for ; it < iterations && len(copied) > 0; it++ {
		visible, invisible := simulation.CompileTree(copied)

		simulation.Interact(universe.PhysicsConstantTimestep, invisible)

		for _, e := range copied {
			e.Step(universe.PhysicsConstantTimestep)
		}

		simulation.ResolveCollisions(visible, v.collisionResolver)

		(&copied).DeleteSweep(func(e *universe.Entity) bool {
			return e.Disabled
		})

		for _, e := range copied {
			if e.V.X == 0 && e.V.Y == 0 {
				continue
			}
			path[e] = append(path[e], e.P)
		}

	}

	for e, paths := range path {
		if len(paths) < 2 {
			continue
		}

		if e.Disabled {
			ctx.Push(draw.StrokeStyle, colorDanger)
		} else {
			ctx.Push(draw.StrokeStyle, colorDefault)
		}

		ctx.BeginPath()

		for i, xy := range paths {
			if i == 0 {
				ctx.MoveTo(xy.X, xy.Y)
				continue
			}

			ctx.LineTo(xy.X, xy.Y)
		}

		ctx.Stroke()

		ctx.Pop(draw.StrokeStyle)
	}
}
