package controller

import (
	"github.com/relvacode/universe/draw"
	"github.com/relvacode/universe/internal"
)

type Input interface {
	Update(world internal.BoundingBox) internal.Vector
	Draw(v *View, ctx *draw.Context)

	Contains(xy internal.Vector) bool
	Enabled(v *View) bool
	Click(v *View, xy internal.Vector)
}

func NewInputController(world internal.BoundingBox, inputs ...Input) *InputController {
	ctrl := &InputController{
		inputs: inputs,
	}
	ctrl.Update(world)
	return ctrl
}

type InputController struct {
	inputs []Input
}

func (i *InputController) Get(v *View, xy internal.Vector) Input {
	for _, input := range i.inputs {
		if input.Contains(xy) {
			if input.Enabled(v) {
				return input
			}
			return nil
		}
	}

	return nil
}

func (i *InputController) Update(world internal.BoundingBox) {
	layout := internal.BoundingBox{
		X: 18,
		Y: (world.Y + world.H) - (32 + 18),
		H: 32,
		W: (world.X + world.W) - 32,
	}
	for _, input := range i.inputs {
		size := input.Update(layout)
		layout.X += size.X + 12
	}
}

func (i *InputController) Draw(v *View, ctx *draw.Context) {
	for _, input := range i.inputs {
		input.Draw(v, ctx)
	}
}

type iconButtonInput struct {
	box internal.BoundingBox
}


func (i *iconButtonInput) Update(world internal.BoundingBox) internal.Vector {
	i.box = internal.BoundingBox{
		X: world.X,
		Y: world.Y,
		W: 32,
		H: 32,
	}
	return internal.Vector{
		X: 32,
		Y: 32,
	}
}

func (i *iconButtonInput) Contains(xy internal.Vector) bool {
	return i.box.ContainsPoint(xy)
}

func (i *iconButtonInput) drawIcon(ctx *draw.Context, icon string) {
	ctx.Push(draw.Font, fontSolidIcons)
	ctx.FillText(icon, i.box.X + 8, i.box.Y + 7)
	ctx.Pop(draw.Font)
	//
	//ctx.Rect(i.box.X, i.box.Y, i.box.W, i.box.H)
	//ctx.Stroke()
}

type PlayStateInput struct {
	iconButtonInput
}

func (PlayStateInput) Enabled(_ *View) bool {
	return true
}

func (PlayStateInput) iconForState(v *View) string {
	if v.paused {
		return iconSolidPlay
	}
	return iconSolidPause
}

func (i *PlayStateInput) Draw(v *View, ctx *draw.Context) {
	i.drawIcon(ctx, i.iconForState(v))
}

func (i *PlayStateInput) Click(v *View, _ internal.Vector) {
	v.paused = !v.paused
}

type DeleteStateInput struct {
	iconButtonInput
}

func (DeleteStateInput) Enabled(v *View) bool {
	return len(*v.entities) > 0
}

func (i *DeleteStateInput) Draw(v *View, ctx *draw.Context) {
	if len(*v.entities) == 0 {
		ctx.Push(draw.GlobalAlpha, .2)
		defer ctx.Pop(draw.GlobalAlpha)
	}

	i.drawIcon(ctx, iconSolidTrash)
}

func (i *DeleteStateInput) Click(v *View, _ internal.Vector) {
	v.entities.Clear()
}
