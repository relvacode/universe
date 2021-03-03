package universe

import (
	"github.com/relvacode/universe/draw"
	"github.com/relvacode/universe/internal"
	"github.com/relvacode/universe/physics"
	"math"
)

type EntityList []*Entity

func (el *EntityList) Clear() {
	*el = (*el)[:0]
}

func (el *EntityList) DeleteSweep(f func(e *Entity) bool) int {
	var i, n int
	for _, e := range *el {
		if !f(e) {
			(*el)[i] = e
			i++
		} else {
			n++
		}
	}

	for j := i; j < len(*el); j++ {
		(*el)[j] = nil
	}
	*el = (*el)[:i]
	return n
}

func NewEntity(p internal.Vector, v internal.Vector, r float64) *Entity {
	return &Entity{
		Object: physics.Object{
			P: p,
			V: v,
			R: r,
			M: math.Pow(r, 3),
		},
	}
}

type Entity struct {
	physics.Object
	Disabled bool
}

func (e *Entity) Draw(ctx *draw.Context, c draw.Camera) {
	ctx.BeginPath()
	ctx.Arc((e.P.X-c.Offset.X)*c.Zoom, (e.P.Y-c.Offset.Y)*c.Zoom, e.R*c.Zoom, 0, 2*math.Pi)
	ctx.Fill()
	ctx.ClosePath()
}
