// +build js

package main

import (
	"github.com/relvacode/universe"
	"github.com/relvacode/universe/controller"
	"github.com/relvacode/universe/draw"
	"github.com/relvacode/universe/internal"
	"syscall/js"
)

const defaultColor = "#00b3ff"
const defaultFont = "12px sans-serif"

func onResizeEventHandler(r *universe.Renderer) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		el := args[0].Get("target")
		r.Update(internal.BoundingBox{
			X: 0,
			Y: 0,
			W: el.Get("innerWidth").Float(),
			H: el.Get("innerHeight").Float(),
		})

		return nil
	})
}

func createDefaultContext(el js.Value) *draw.Context {
	ctx := draw.GetContext(el)
	ctx.Push(draw.FillStyle, defaultColor)
	ctx.Push(draw.StrokeStyle, defaultColor)
	ctx.Push(draw.LineWidth, 1)
	ctx.Push(draw.Font, defaultFont)
	ctx.Push(draw.TextBaseline, "top")
	ctx.Push(draw.GlobalAlpha, 1)
	return ctx
}

func main() {
	window := js.Global().Get("window")
	document := js.Global().Get("document")

	renderGroup := document.Get("body")

	width := window.Get("innerWidth").Float()
	height := window.Get("innerHeight").Float()

	createContextElement := func() js.Value {
		element := document.Call("createElement", "canvas")
		element.Set("width", width)
		element.Set("height", height)

		renderGroup.Call("appendChild", element)

		return element
	}

	box := internal.BoundingBox{
		X: 0, Y: 0, W: width, H: height,
	}

	vc := controller.New(box, universe.AbsorbCollisionResolver)
	vc.Bind(document)

	r := universe.NewRenderer(
		createDefaultContext(createContextElement()),
		createDefaultContext(createContextElement()),
		createDefaultContext(createContextElement()),
		box,
		vc,
	)
	//
	//w := &universe.Renderer{
	//	EntityContext:     ,
	//	ViewContext:       ,
	//	DebugContext:      ,
	//	View:              vc,
	//	Box:               box,
	//	CollisionResolver: universe.AbsorbCollisionResolver,
	//}

	window.Call("addEventListener", "resize", onResizeEventHandler(r))
	r.Loop()
}
