package controller

import (
	"fmt"
	"github.com/relvacode/universe"
	"github.com/relvacode/universe/internal"
	"github.com/relvacode/universe/physics"
	"math/rand"
)

type EntityPlacer func(box internal.BoundingBox) []*universe.Entity

func generateRandomField(n int) EntityPlacer {
	return func(box internal.BoundingBox) (objects []*universe.Entity) {
	place:
		for i := 0; i < n; {
			en := universe.NewEntity(
				internal.Vector{
					X: box.X + (box.W * rand.Float64()),
					Y: box.Y + (box.H * rand.Float64()),
				},
				internal.Vector{},
				2,
			)

			for _, o := range objects {
				if _, ok := physics.Colliding(en.Object, o.Object); ok {
					fmt.Println("place collision detected")
					continue place
				}
			}

			objects = append(objects, en)
			i++
		}
		return
	}
}
