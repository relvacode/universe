package universe

import (
	"github.com/relvacode/universe/draw"
	"github.com/relvacode/universe/internal"
)

type View interface {
	TimeScale() float64
	CollisionResolver() CollisionResolver
	Camera() *draw.Camera
	Entities() *EntityList

	Update(box internal.BoundingBox)
	Draw(ctx *draw.Context)
}