package universe

import (
	"github.com/relvacode/universe/internal"
	"github.com/relvacode/universe/physics"
	"math"
)


type entityCollisionMap map[[2]*Entity]struct{}

func (m entityCollisionMap) check(e1, e2 *Entity) bool {
	_, ok := m[[2]*Entity{e1, e2}]
	return ok
}

func (m entityCollisionMap) store(e1, e2 *Entity) {
	m[[2]*Entity{e2, e1}] = struct{}{}
}

func NewSimulation(worldBoundary internal.BoundingBox) *Simulation {
	return &Simulation{
		worldBoundary: worldBoundary,
	}
}

type Simulation struct {
	tree          *QuadTree
	worldBoundary internal.BoundingBox
}

func (s *Simulation) ResolveCollisions(entities []*Entity, resolver CollisionResolver) {
	var collisionMap = make(entityCollisionMap)

	var collisions int
	for _, o := range entities {
		if o.Disabled {
			continue
		}
		s.tree.Intersections(o.BoundingBox(), func(m *Entity) bool {
			if m.Disabled || o == m {
				return true
			}

			distance, colliding := physics.Colliding(o.Object, m.Object)

			if !colliding || distance == math.NaN() || collisionMap.check(o, m) {
				return true
			}

			collisionMap.store(o, m)

			collisions++

			resolver(o, m, distance)
			return !o.Disabled
		})
	}
}

func (s *Simulation) Interact(timestep float64, invisible EntityList) {
	var leaves = s.tree.AppendLeaves(nil)

	for il := 0; il < len(leaves); il++ {
		l1 := leaves[il]

		for _, e1 := range invisible {
			if e1.Disabled {
				continue
			}
			coM := l1.CenterOfMass()

			impulse := physics.AttractionForceImpulse(e1.P, coM, e1.M, l1.totalMass, timestep)
			e1.V.X += impulse.X / e1.M
			e1.V.Y += impulse.Y / e1.M

			impulseLeafShare := internal.Vector{
				X: -impulse.X / l1.totalMass,
				Y: -impulse.Y / l1.totalMass,
			}

			for i := 0; i < len(l1.objects); i++ {
				e2 := l1.objects[i]
				eMassRatio := e2.M / l1.totalMass
				e2.V.X += impulseLeafShare.X * eMassRatio
				e2.V.Y += impulseLeafShare.Y * eMassRatio
			}
		}

		for i := 0; i < len(l1.objects); i++ {
			e1 := l1.objects[i]
			if e1.Disabled {
				continue
			}

			for j := 0; j < i; j++ {
				e2 := l1.objects[j]
				if e2.Disabled {
					continue
				}

				impulse := physics.AttractionForceImpulse(e1.P, e2.P, e1.M, e2.M, timestep)

				e1.V.X += impulse.X / e1.M
				e1.V.Y += impulse.Y / e1.M

				e2.V.X += -impulse.X / e2.M
				e2.V.Y += -impulse.Y / e2.M
			}

			for jl := 0; jl < len(leaves); jl++ {
				if il == jl {
					continue
				}

				l2 := leaves[jl]

				impulse := physics.AttractionForceImpulse(e1.P, l2.CenterOfMass(), e1.M, l2.totalMass, timestep)
				e1.V.X += impulse.X / e1.M
				e1.V.Y += impulse.Y / e1.M
			}
		}
	}
}

func (s *Simulation) CompileTree(entities EntityList) (EntityList, EntityList) {
	s.tree = NewQuadTree(s.worldBoundary, 28, 8)

	// Perform a visibility sort.
	// Entities that are not in the current view are shuffled to the end until all entities are sorted.
	// The resulting visible cursor is the number of entities that are still visible.
	var visible int
	var curE = len(entities)
	for ; visible < curE; {
		e := entities[visible]

		ok := s.tree.boundary.Intersects(e.BoundingBox()) && s.tree.Insert(e)

		// Everything ok, the entity at the current position is inside the view
		// Go to the next entity
		if ok {
			visible++
			continue
		}

		// The entity is not inside the view
		// Swap the current index with an entity at the current end pointer
		// And try that one in the next iteration
		entities[visible], entities[curE-1] = entities[curE-1], entities[visible]
		curE--
	}

	return entities[:visible], entities[visible:]
}
