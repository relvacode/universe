package controller

import (
	"fmt"
	"github.com/relvacode/universe"
	"github.com/relvacode/universe/draw"
	"github.com/relvacode/universe/internal"
	"math"
	"syscall/js"
)

var _ universe.View = (*View)(nil)

func New(initial internal.BoundingBox, collisionResolver universe.CollisionResolver) *View {
	var el universe.EntityList
	return &View{
		box:               initial,
		entities:          &el,
		collisionResolver: collisionResolver,
		timescale:         1,
		inputs: NewInputController(initial,
			new(PlayStateInput),
			new(DeleteStateInput),
		),
		camera: &draw.Camera{
			Zoom: 1,
		},
	}
}

type View struct {
	box               internal.BoundingBox
	entities          *universe.EntityList
	collisionResolver universe.CollisionResolver

	timescale float64
	paused    bool

	camera *draw.Camera

	inputs       *InputController
	targetCursor string
	mouseHandler MouseHandler
}

func (v *View) TimeScale() float64 {
	if v.paused {
		return 0
	}
	return v.timescale
}

func (v *View) CollisionResolver() universe.CollisionResolver {
	return v.collisionResolver
}

func (v *View) Camera() *draw.Camera {
	return v.camera
}

func (v *View) Entities() *universe.EntityList {
	return v.entities
}

func (v *View) Update(world internal.BoundingBox) {
	v.box = world
	v.inputs.Update(world)
}

func (v *View) SetCursor(cursor string) {
	if v.targetCursor == cursor {
		return
	}

	v.targetCursor = cursor
	if cursor == "" {
		cursor = "initial"
	}
	fmt.Println("set cursor to", cursor)
	js.Global().Get("document").Get("body").Get("style").Set("cursor", cursor)
}

func (v *View) Draw(ctx *draw.Context) {
	ctx.Clear(v.box)

	if v.paused && len(*v.entities) < 256 {
		v.drawForwardProjections(ctx)
	}

	if v.mouseHandler != nil {
		v.mouseHandler.Draw(v, ctx)
	}

	ctx.Push(draw.FillStyle, "#FFFFFF")
	v.inputs.Draw(v, ctx)
	ctx.Pop(draw.FillStyle)
}

func (v *View) modifyZoomLevel(in bool) {
	if !in {
		v.camera.Zoom /= 2
		if v.camera.Zoom < 1 {
			v.camera.Zoom = 1
		}
	} else {
		v.camera.Zoom *= 2
	}

	v.camera.FitCenter(v.box)
}

func (v *View) findEntityAtTarget(xy internal.Vector) *universe.Entity {
	for _, e := range *v.entities {
		if !e.BoundingBox().ContainsPoint(xy) {
			continue
		}

		delta := internal.Vector{
			X: e.P.X - xy.X,
			Y: e.P.Y - xy.Y,
		}

		distance := math.Sqrt(delta.Dot())
		if distance < e.R {
			return e
		}
	}

	return nil
}
