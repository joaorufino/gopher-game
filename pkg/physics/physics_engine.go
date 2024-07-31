package physics

import (
	"github.com/joaorufino/gopher-game/internal/interfaces"
)

type PhysicsEngine struct {
	RigidBodies  []interfaces.RigidBody
	gravity      interfaces.Vector2D
	floorY       float64
	eventManager interfaces.EventManager
}

func NewPhysicsEngine(eventManager interfaces.EventManager, gravity interfaces.Vector2D, floorY float64) *PhysicsEngine {
	return &PhysicsEngine{
		RigidBodies:  make([]interfaces.RigidBody, 0),
		gravity:      gravity,
		floorY:       floorY,
		eventManager: eventManager,
	}
}

func (pe *PhysicsEngine) AddRigidBody(rb interfaces.RigidBody) {
	pe.RigidBodies = append(pe.RigidBodies, rb)
}

func (pe *PhysicsEngine) RemoveRigidBody(rb interfaces.RigidBody) {
	for i, r := range pe.RigidBodies {
		if r == rb {
			pe.RigidBodies = append(pe.RigidBodies[:i], pe.RigidBodies[i+1:]...)
			return
		}
	}
}

func (pe *PhysicsEngine) DetectCollision(rb1, rb2 interfaces.RigidBody) bool {
	return CheckCollisionOnX(rb1.(*RigidBody), rb2.(*RigidBody)) &&
		CheckCollisionOnY(rb1.(*RigidBody), rb2.(*RigidBody))
}

func (pe *PhysicsEngine) ResolveCollision(rb1, rb2 interfaces.RigidBody) {
	if rb1.GetCanPick() && rb2.GetPickable() {
		pe.eventManager.Dispatch(interfaces.Event{
			Type:     interfaces.EventItemEquipped,
			Priority: 1,
			Payload: map[string]interface{}{
				"itemName": rb2.GetIdentifier(),
			},
		})
		pe.RemoveRigidBody(rb2)
	}
	ResolveCollision(rb1.(*RigidBody), rb2.(*RigidBody))
}

func (pe *PhysicsEngine) Update(deltaTime float64) {
	for _, rb := range pe.RigidBodies {
		if !rb.(*RigidBody).IsStatic {
			// Apply gravity
			rb.(*RigidBody).ApplyForce(interfaces.Vector2D{
				X: 0,
				Y: pe.gravity.Y * rb.(*RigidBody).Mass,
			})

			// Update velocity
			rb.(*RigidBody).Velocity.X += rb.(*RigidBody).Acceleration.X * deltaTime
			rb.(*RigidBody).Velocity.Y += rb.(*RigidBody).Acceleration.Y * deltaTime

			// Update position
			rb.Update(deltaTime)

			// Reset acceleration
			rb.(*RigidBody).Acceleration = interfaces.Vector2D{X: 0, Y: 0}
		}
	}

	// Check for collisions and resolve them
	for i := 0; i < len(pe.RigidBodies); i++ {
		rb1 := pe.RigidBodies[i]
		for j := i + 1; j < len(pe.RigidBodies); j++ {
			rb2 := pe.RigidBodies[j]
			if rb2.GetIsStatic() && rb1.GetIsStatic() {
				continue
			}

			if pe.DetectCollision(rb1, rb2) {
				pe.ResolveCollision(rb1, rb2)
			}
		}
	}

	// Check for floor collision and reset OnGround flag if necessary
	for _, rb := range pe.RigidBodies {
		if !rb.(*RigidBody).IsStatic {
			rb.(*RigidBody).OnGround = false // Reset OnGround before checking

			// Check for floor collision
			if rb.(*RigidBody).Position.Y >= pe.floorY {
				rb.(*RigidBody).Position.Y = pe.floorY
				rb.(*RigidBody).Velocity.Y = 0
				rb.(*RigidBody).OnGround = true
			}

			// Check if the rigid body is on top of any platforms or obstacles
			for _, otherRb := range pe.RigidBodies {
				if otherRb != rb && otherRb.GetIsStatic() && CheckIfOnTop(rb.(*RigidBody), otherRb.(*RigidBody)) {
					rb.(*RigidBody).OnGround = true
					break
				}
			}
		}
	}
}

func (pe *PhysicsEngine) GetRigidBodies() []interfaces.RigidBody {
	return pe.RigidBodies
}
