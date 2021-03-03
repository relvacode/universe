package universe

import (
	"github.com/relvacode/universe/draw"
	"github.com/relvacode/universe/internal"
	"math"
)

func NewQuadTree(boundary internal.BoundingBox, maxEntries, maxDepth int) *QuadTree {
	return &QuadTree{
		boundary:       boundary,
		maxEntries:     maxEntries,
		remainingDepth: maxDepth,
	}
}

type QuadTree struct {
	boundary internal.BoundingBox
	objects  []*Entity

	maxEntries     int
	remainingDepth int

	totalMassVector internal.Vector
	totalMass       float64

	nw *QuadTree
	ne *QuadTree
	sw *QuadTree
	se *QuadTree
}

func (qt *QuadTree) Nodes() int {
	var c int
	if len(qt.objects) > 0 {
		c = 1
	}
	if qt.nw == nil {
		return c
	}

	c += qt.nw.Nodes()
	c += qt.ne.Nodes()
	c += qt.sw.Nodes()
	c += qt.se.Nodes()

	return c
}

func (qt *QuadTree) insertChildren(e *Entity) bool {
	return qt.nw.Insert(e) || qt.ne.Insert(e) || qt.sw.Insert(e) || qt.se.Insert(e)
}

func (qt *QuadTree) Subdivide() {
	depth := qt.remainingDepth - 1
	w := qt.boundary.W / 2
	h := qt.boundary.H / 2
	qt.nw = &QuadTree{
		maxEntries:     qt.maxEntries,
		remainingDepth: depth,
		boundary: internal.BoundingBox{
			X: qt.boundary.X,
			Y: qt.boundary.Y,
			W: w,
			H: h,
		},
	}
	qt.ne = &QuadTree{
		maxEntries:     qt.maxEntries,
		remainingDepth: depth,
		boundary: internal.BoundingBox{
			X: qt.boundary.X + w,
			Y: qt.boundary.Y,
			W: w,
			H: h,
		},
	}
	qt.se = &QuadTree{
		maxEntries:     qt.maxEntries,
		remainingDepth: depth,
		boundary: internal.BoundingBox{
			X: qt.boundary.X,
			Y: qt.boundary.Y + h,
			W: w,
			H: h,
		},
	}
	qt.sw = &QuadTree{
		maxEntries:     qt.maxEntries,
		remainingDepth: depth,
		boundary: internal.BoundingBox{
			X: qt.boundary.X + w,
			Y: qt.boundary.Y + h,
			W: w,
			H: h,
		},
	}

	for i := 0; i < len(qt.objects); i++ {
		qt.insertChildren(qt.objects[i])
	}

	qt.objects = nil
}

func (qt *QuadTree) Insert(e *Entity) bool {
	if !qt.boundary.ContainsPoint(e.P) {
		return false
	}

	qt.totalMassVector.X += e.P.X * e.M
	qt.totalMassVector.Y += e.P.Y * e.M
	qt.totalMass += e.M

	if qt.nw == nil && (len(qt.objects) < qt.maxEntries || qt.remainingDepth == 0) {
		qt.objects = append(qt.objects, e)
		return true
	}

	if qt.nw == nil {
		qt.Subdivide()
	}

	return qt.insertChildren(e)
}

func (qt *QuadTree) CenterOfMass() internal.Vector {
	return internal.Vector{
		X: qt.totalMassVector.X / qt.totalMass,
		Y: qt.totalMassVector.Y / qt.totalMass,
	}
}

// Intersections iterates over all entities that intersect each tree node.
// If the callback function returns false, then not more entities are evaluated.
// Returns true if all possible entities were iterated over.
func (qt *QuadTree) Intersections(box internal.BoundingBox, f func(m *Entity) bool) bool {
	if !qt.boundary.Intersects(box) {
		return true
	}

	for i := 0; i < len(qt.objects); i++ {
		if !f(qt.objects[i]) {
			return false
		}
	}

	if qt.nw == nil {
		return true
	}

	return qt.nw.Intersections(box, f) && qt.ne.Intersections(box, f) && qt.sw.Intersections(box, f) && qt.se.Intersections(box, f)
}

func (qt *QuadTree) AppendLeaves(arr []*QuadTree) []*QuadTree {
	if qt.nw == nil {
		if len(qt.objects) > 0 {
			arr = append(arr, qt)
		}
		return arr
	}

	arr = qt.nw.AppendLeaves(arr)
	arr = qt.ne.AppendLeaves(arr)
	arr = qt.sw.AppendLeaves(arr)
	arr = qt.se.AppendLeaves(arr)

	return arr
}

func (qt *QuadTree) Draw(ctx *draw.Context) {
	if len(qt.objects) > 0 {
		ctx.BeginPath()
		ctx.MoveTo(qt.boundary.X, qt.boundary.Y)
		ctx.LineTo(qt.boundary.X+qt.boundary.W, qt.boundary.Y)
		ctx.LineTo(qt.boundary.X+qt.boundary.W, qt.boundary.Y+qt.boundary.H)
		ctx.LineTo(qt.boundary.X, qt.boundary.Y+qt.boundary.H)
		ctx.LineTo(qt.boundary.X, qt.boundary.Y)
		ctx.Stroke()

		ctx.BeginPath()

		c := qt.CenterOfMass()
		ctx.Arc(c.X, c.Y, math.Cbrt(qt.totalMass), 0, math.Pi*2)
		ctx.Stroke()
	}

	if qt.nw != nil {
		qt.nw.Draw(ctx)
		qt.ne.Draw(ctx)
		qt.sw.Draw(ctx)
		qt.se.Draw(ctx)
	}
}
