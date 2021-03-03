package draw

import (
	"github.com/relvacode/universe/internal"
	"syscall/js"
)

//

func GetContext(el js.Value) *Context {
	var state styleStateMachine
	return &Context{
		c:     el.Call("getContext", "2d"),
		state: &state,
	}
}

type Context struct {
	c           js.Value
	state       *styleStateMachine
	cursorState string
}

func (ctx *Context) Resize(w, h float64) {
	canvas := ctx.c.Get("canvas")
	canvas.Set("width", w)
	canvas.Set("height", h)
	ctx.c = canvas.Call("getContext", "2d")
	for _, style := range ctx.state {
		if style != nil {
			style.apply(ctx.c)
		}
	}
}

func (ctx *Context) Push(attr attribute, value interface{}) {
	ctx.state.push(attr, value).apply(ctx.c)
}

func (ctx *Context) Pop(attr attribute) {
	n := ctx.state.pop(attr)
	if n == nil {
		ctx.c.Set(attr.String(), nil)
		return
	}

	n.apply(ctx.c)
}

func (ctx *Context) Get(attr attribute) js.Value {
	return ctx.c.Get(attr.String())
}

//func (ctx *Context) Save() {
//	ctx.c.Call("save")
//}
//
//func (ctx *Context) Restore() {
//	ctx.c.Call("restore")
//}

func (ctx *Context) ClearRect(x, y, w, h float64) {
	ctx.c.Call("clearRect", x, y, w, h)
}

func (ctx *Context) Clear(box internal.BoundingBox) {
	ctx.ClearRect(box.X, box.Y, box.X+box.W, box.Y+box.H)
}

func (ctx *Context) BeginPath() {
	ctx.c.Call("beginPath")
}

func (ctx *Context) Rect(x, y, w, h float64) {
	ctx.c.Call("rect", x, y, w, h)
}

func (ctx *Context) MoveTo(x, y float64) {
	ctx.c.Call("moveTo", int(x), int(y))
}

func (ctx *Context) LineTo(x, y float64) {
	ctx.c.Call("lineTo", int(x), int(y))
}

func (ctx *Context) Arc(x, y, radius, startAngle, endAngle float64) {
	ctx.c.Call("arc", x, y, radius, startAngle, endAngle)
}

func (ctx *Context) Fill() {
	ctx.c.Call("fill")
}

func (ctx *Context) Stroke() {
	ctx.c.Call("stroke")
}

func (ctx *Context) ClosePath() {
	ctx.c.Call("closePath")
}

func (ctx *Context) FillText(text string, x, y float64) {
	ctx.c.Call("fillText", text, x, y)
}

func (ctx *Context) MeasureTextWidth(text string) float64 {
	return ctx.c.Call("measureText", text).Get("width").Float()
}

func (ctx *Context) MeasureFontHeight() float64 {
	// The width of `W` is an approximation into the height of the font
	return ctx.MeasureTextWidth("W")
}

func (ctx *Context) DrawImage(src interface{}, x, y float64) {
	ctx.c.Call("drawImage", src, x, y)
}

//func (ctx *Context) Get(attribute string) js.Value {
//	return ctx.c.Get(attribute)
//}
//
//func (ctx *Context) Set(attribute string, value interface{}) {
//	ctx.c.Set(attribute, value)
//}
