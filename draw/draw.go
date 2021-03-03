package draw

import (
	"github.com/relvacode/universe/internal"
)

type Renderer interface {
	Draw(ctx *Context, c *Camera)
}

type Camera struct {
	Zoom   float64
	Offset internal.Vector
}

func (c Camera) Crop(global internal.BoundingBox) internal.BoundingBox {
	return internal.BoundingBox{
		X: global.X + c.Offset.X,
		Y: global.Y + c.Offset.Y,
		W: global.W / c.Zoom,
		H: global.H / c.Zoom,
	}
}

func (c *Camera) FitCenter(global internal.BoundingBox) {
	c.Offset.X = (global.W - (global.W / c.Zoom)) / 2
	c.Offset.Y = (global.H - (global.H / c.Zoom)) / 2
}