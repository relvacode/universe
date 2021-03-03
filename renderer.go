package universe

import (
	"fmt"
	"github.com/relvacode/universe/draw"
	"github.com/relvacode/universe/internal"
	"syscall/js"
	"time"
)

const PhysicsConstantTimestep = 0.0165 // 60FPS

func NewRenderer(
	entityContext, viewContext, debugContext *draw.Context,
	worldBoundary internal.BoundingBox, view View) *Renderer {
	return &Renderer{
		contextEntities: entityContext,
		contextView:     viewContext,
		contextDebug:    debugContext,
		worldBoundary:   worldBoundary,

		simulation: NewSimulation(worldBoundary),
		view:       view,
	}
}

type Renderer struct {
	contextEntities *draw.Context
	contextView     *draw.Context
	contextDebug    *draw.Context
	worldBoundary   internal.BoundingBox
	view            View

	simulation *Simulation

	timeStepRemaining float64
	frame             int

	perfPhysicsIterations int
	perfInteractions      time.Duration
	perfCollisions        time.Duration
	perfPhysics           time.Duration
	perfDraw              time.Duration
}

func (r *Renderer) reset(ctx *draw.Context) {
	ctx.ClearRect(r.worldBoundary.X, r.worldBoundary.Y, r.worldBoundary.W, r.worldBoundary.H)
}

func (r *Renderer) fps(timestep float64, ctx *draw.Context) {
	r.frame++
	var y float64 = 20
	h := ctx.MeasureFontHeight() + 2

	ctx.FillText(fmt.Sprintf("%.3fms, %.2ffps, frame %d, %dit", timestep, 1000/timestep, r.frame, r.perfPhysicsIterations), 10, y)
	y += h

	ctx.FillText(fmt.Sprintf("physics %s", r.perfPhysics), 10, y)
	y += h

	ctx.FillText(fmt.Sprintf("collisions %s", r.perfCollisions), 20, y)
	y += h

	ctx.FillText(fmt.Sprintf("interactions %s", r.perfInteractions), 20, y)
	y += h

	ctx.FillText(fmt.Sprintf("draw %s", r.perfDraw), 10, y)
	y += h
}

func (r *Renderer) runConstantTimeStep(timestep float64, entities EntityList) {
	if len(entities) == 0 {
		return
	}

	perf := time.Now()

	visible, invisible := r.simulation.CompileTree(entities)

	perf = time.Now()
	r.simulation.Interact(timestep, invisible)
	r.perfInteractions = time.Now().Sub(perf)

	for _, o := range entities {
		o.Step(timestep)
	}

	perf = time.Now()
	r.simulation.ResolveCollisions(visible, r.view.CollisionResolver())
	r.perfCollisions = time.Now().Sub(perf)

	return
}

func (r *Renderer) Render(timestep float64) {
	if timestep == 0 {
		return
	}

	timestepSecs := timestep / 1000

	if timestepSecs > 1 {
		return
	}

	timeScale := ((timestep / 1000) * r.view.TimeScale()) + r.timeStepRemaining
	entities := r.view.Entities()

	var perf time.Time

	var physicsIterations int
	if timeScale > 0 {
		perf = time.Now()

		// Run physics for as long as at least one full iteration can run
		for ; timeScale >= PhysicsConstantTimestep; timeScale -= PhysicsConstantTimestep {
			r.runConstantTimeStep(PhysicsConstantTimestep, *entities)
			physicsIterations++
		}

		// Any remaining time is deferred to the next frame
		r.timeStepRemaining = timeScale
		r.perfPhysics = time.Now().Sub(perf)
	}

	r.perfPhysicsIterations = physicsIterations

	entities.DeleteSweep(func(e *Entity) bool {
		return e.Disabled
	})

	perf = time.Now()
	r.reset(r.contextEntities)

	camera := r.view.Camera()
	cameraBounds := camera.Crop(r.worldBoundary)

	for _, o := range *entities {
		// Do not draw entities that are not within the bounds of the current camera
		if !cameraBounds.Intersects(o.BoundingBox()) {
			continue
		}

		o.Draw(r.contextEntities, *camera)
	}

	r.view.Draw(r.contextView)

	r.reset(r.contextDebug)
	r.perfDraw = time.Now().Sub(perf)

	r.fps(timestep, r.contextDebug)
}

func (r *Renderer) Loop() {
	var renderFrame js.Func
	var lastTimestamp float64

	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		timestamp := args[0].Float()
		timestep := timestamp - lastTimestamp
		lastTimestamp = timestamp

		r.Render(timestep)

		js.Global().Call("requestAnimationFrame", renderFrame)
		//c.reqID =  // Captures the requestID to be used in Close / Cancel
		return nil
	})
	defer renderFrame.Release()
	js.Global().Call("requestAnimationFrame", renderFrame)
	<-make(chan struct{})
}

func (r *Renderer) Update(box internal.BoundingBox) {
	r.worldBoundary = box

	r.contextEntities.Resize(box.W, box.H)
	r.contextView.Resize(box.W, box.H)
	r.contextDebug.Resize(box.W, box.H)

	r.view.Update(box)
}
