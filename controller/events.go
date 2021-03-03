package controller

import (
	"fmt"
	"github.com/relvacode/universe"
	"github.com/relvacode/universe/draw"
	"github.com/relvacode/universe/internal"
	"math"
	"syscall/js"
)

type MouseHandler interface {
	Draw(v *View, ctx *draw.Context)
	Move(v *View, xy internal.Vector) bool
	Release(v *View, xy internal.Vector) bool
}

func (v *View) onKeypressEventHandler() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		keyCode := event.Get("keyCode").Int()
		fmt.Println(keyCode)

		switch keyCode {
		case ' ':
			v.paused = !v.paused
		case 'r':
			v.entities.Clear()
		case 's':
			*v.entities = append(*v.entities, generateRandomField(1024)(v.box)...)
		case 'v':
			for _, e := range *v.entities {
				e.V = internal.Vector{}
			}
		case '_':
			v.modifyZoomLevel(false)
		case '+':
			v.modifyZoomLevel(true)
		}

		return nil
	})
}

func (v *View) mouseEventOrigin(event js.Value) internal.Vector {
	return internal.Vector{
		X: event.Get("offsetX").Float(),
		Y: event.Get("offsetY").Float(),
	}
}

func (v *View) onMouseUpEventHandler() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if v.mouseHandler == nil {
			return nil
		}

		origin := v.mouseEventOrigin(args[0])
		if v.mouseHandler.Release(v, origin) {
			v.mouseHandler = nil
		}

		return nil
	})
}

func (v *View) onMouseMoveEventHandler() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		origin := v.mouseEventOrigin(args[0])

		if v.mouseHandler != nil {
			if v.mouseHandler.Move(v, origin) {
				v.mouseHandler = nil
			}
			return nil
		}

		if v.inputs.Get(v, origin) != nil {
			v.SetCursor("pointer")
			return nil
		}

		v.SetCursor("")

		return nil
	})
}

func (v *View) onMouseDownEventHandler() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		origin := v.mouseEventOrigin(args[0])

		if i := v.inputs.Get(v, origin); i != nil {
			v.mouseHandler = &inputMouseHandler{
				Input: i,
			}
			return nil
		}

		target := v.findEntityAtTarget(origin)
		if target != nil {
			v.mouseHandler = &VelocityModifier{
				target:  target,
				initial: origin,
				final:   origin,
			}
			return nil
		}

		if args[0].Call("getModifierState", "Alt").Bool() {
			v.mouseHandler = &Deleter{
				initial: origin,
				final:   origin,
			}
			return nil
		}

		v.mouseHandler = &EntitySpawner{
			origin: origin,
		}

		return nil
	})
}

func (v *View) Bind(el js.Value) {
	el.Call("addEventListener", "keypress", v.onKeypressEventHandler())

	el.Call("addEventListener", "mousedown", v.onMouseDownEventHandler())
	el.Call("addEventListener", "mousemove", v.onMouseMoveEventHandler())
	el.Call("addEventListener", "mouseup", v.onMouseUpEventHandler())
}

var _ MouseHandler = (*inputMouseHandler)(nil)

type inputMouseHandler struct {
	Input
}

func (i *inputMouseHandler) Move(_ *View, xy internal.Vector) bool {
	return !i.Input.Contains(xy)
}

func (i *inputMouseHandler) Release(v *View, xy internal.Vector) bool {
	if i.Input.Contains(xy) {
		i.Input.Click(v, xy)
	}
	return true
}

var _ MouseHandler = (*Deleter)(nil)

type Deleter struct {
	initial internal.Vector
	final   internal.Vector
}

func (d *Deleter) box() internal.BoundingBox {
	return internal.BoundingBox{
		X: d.initial.X,
		Y: d.initial.Y,
		W: d.final.X - d.initial.X,
		H: d.final.Y - d.initial.Y,
	}.Abs()
}

func (d *Deleter) Draw(_ *View, ctx *draw.Context) {
	ctx.BeginPath()

	box := d.box()
	ctx.Rect(box.X, box.Y, box.W, box.H)
	ctx.Stroke()
}

func (d *Deleter) Move(_ *View, xy internal.Vector) bool {
	d.final = xy
	return false
}

func (d *Deleter) Release(v *View, xy internal.Vector) bool {
	d.final = xy
	box := d.box()

	v.entities.DeleteSweep(func(e *universe.Entity) bool {
		return box.Contains(e.BoundingBox())
	})

	return true
}

var _ MouseHandler = (*VelocityModifier)(nil)

type VelocityModifier struct {
	target *universe.Entity

	initial internal.Vector
	final   internal.Vector
}

func (vm *VelocityModifier) targetVelocity() internal.Vector {
	return internal.Vector{
		X: (vm.final.X - vm.initial.X) * 4,
		Y: (vm.final.Y - vm.initial.Y) * 4,
	}
}

func (vm *VelocityModifier) Draw(_ *View, ctx *draw.Context) {
	ctx.BeginPath()
	ctx.MoveTo(vm.initial.X, vm.initial.Y)
	ctx.LineTo(vm.final.X, vm.final.Y)
	ctx.Stroke()
	ctx.ClosePath()

	target := vm.targetVelocity()
	ctx.FillText(fmt.Sprintf("%.2f, %.2f", target.X, target.Y), vm.final.X, vm.final.Y)

	//ctx.FillText(fmt.Sprintf("%.0f, %.0f", vm.target.P.X, vm.target.P.Y), vm.target.P.X, vm.target.P.Y+vm.target.R+10)
}

func (vm *VelocityModifier) Move(v *View, xy internal.Vector) bool {
	vm.final = xy

	if v.paused {
		vm.target.V = vm.targetVelocity()
	}

	return false
}

func (vm *VelocityModifier) Release(v *View, xy internal.Vector) bool {
	vm.final = xy
	vm.target.V = vm.targetVelocity()
	return true
}

var _ MouseHandler = (*EntitySpawner)(nil)

type EntitySpawner struct {
	origin internal.Vector
	edge   *internal.Vector
}

func (s *EntitySpawner) radius() float64 {
	if s.edge == nil || (s.origin.X == s.edge.X && s.origin.Y == s.edge.Y) {
		return 1
	}

	delta := internal.Vector{
		X: s.origin.X - s.edge.X,
		Y: s.origin.Y - s.edge.Y,
	}

	radius := math.Abs(math.Sqrt(delta.Dot())) / 4
	if radius < 1 {
		return 1
	}

	return radius
}

func (s *EntitySpawner) drawOuterEdge(ctx *draw.Context) {
	ctx.BeginPath()
	ctx.Arc(s.origin.X, s.origin.Y, s.radius(), 0, math.Pi*2)
	ctx.Stroke()
	ctx.ClosePath()
}

func (s *EntitySpawner) Draw(_ *View, ctx *draw.Context) {
	ctx.BeginPath()
	ctx.Arc(s.origin.X, s.origin.Y, 1, 0, math.Pi*2)
	ctx.Fill()
	ctx.ClosePath()

	ctx.FillText(fmt.Sprintf("%.2f", s.radius()), s.origin.X+10, s.origin.Y)

	if s.edge != nil {
		s.drawOuterEdge(ctx)
	}
}

func (s *EntitySpawner) Move(_ *View, xy internal.Vector) bool {
	s.edge = &xy
	return false
}

func (s *EntitySpawner) Release(v *View, xy internal.Vector) bool {
	s.edge = &xy
	*v.entities = append(*v.entities, universe.NewEntity(
		s.origin,
		internal.Vector{},
		s.radius(),
	))
	return true
}
