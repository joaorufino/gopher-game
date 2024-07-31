package interfaces

// PhysicsEngine defines the methods for handling physics and collisions.
type PhysicsEngine interface {
	// AddRigidBody adds a rigid body to the physics engine.
	AddRigidBody(rb RigidBody)
	// RemoveRigidBody removes a rigid body from the physics engine.
	RemoveRigidBody(rb RigidBody)
	// DetectCollision detects a collision between two rigid bodies.
	DetectCollision(rb1, rb2 RigidBody) bool
	// ResolveCollision resolves a collision between two rigid bodies.
	ResolveCollision(rb1, rb2 RigidBody)
	// Update updates the state of the physics engine.
	Update(deltaTime float64)
	GetRigidBodies() []RigidBody
}

// RigidBody represents a physical object in the game.
type RigidBody interface {
	// GetPosition returns the position of the rigid body.
	GetPosition() Vector2D
	// SetPosition sets the position of the rigid body.
	SetPosition(position Vector2D)
	// GetVelocity returns the velocity of the rigid body.
	GetVelocity() Vector2D
	// SetVelocity sets the velocity of the rigid body.
	SetVelocity(velocity Vector2D)
	// GetSize returns the size of the rigid body.
	GetSize() Vector2D
	SetSize(size Vector2D)
	GetIsStatic() bool
	SetPushable(bool)
	GetPushable() bool
	SetCollidable(bool)
	GetCollidable() bool
	GetIdentifier() string
	GetCanPick() bool
	SetCanPick(bool)
	GetPickable() bool
	SetPickable(bool)

	Update(deltaTime float64)
}
