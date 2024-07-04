package physics

import (
	"github.com/joaorufino/cv-game/internal/interfaces"
)

const GRAVITY = 9.8

// RigidBody represents the physical properties of an entity.
type RigidBody struct {
	Identifier      string
	Position        interfaces.Point `json:"position"`
	Velocity        interfaces.Point `json:"velocity"`
	Acceleration    interfaces.Point `json:"acceleration"`
	Mass            float64          `json:"mass"`
	Size            interfaces.Point `json:"size"`
	IsStatic        bool
	OnGround        bool
	IsCollidable    bool
	IsPushable      bool
	IsPickable      bool
	CanPick         bool
	CollidingBodies []*RigidBody
}

// NewRigidBody creates a new RigidBody.
func NewRigidBody(position, size interfaces.Point, mass float64, isStatic bool, identifier string) *RigidBody {
	return &RigidBody{
		Identifier:      identifier,
		Position:        position,
		Size:            size,
		Mass:            mass,
		IsStatic:        isStatic,
		IsCollidable:    true,
		IsPushable:      false,
		OnGround:        false,
		CollidingBodies: []*RigidBody{},
	}
}

// IsStatic
func (rb *RigidBody) GetIsStatic() bool {
	return rb.IsStatic
}

// IsPushable
func (rb *RigidBody) GetPushable() bool {
	return rb.IsPushable
}

// IsCollidable
func (rb *RigidBody) GetCollidable() bool {
	return rb.IsCollidable
}

func (rb *RigidBody) GetPickable() bool {
	return rb.IsPickable
}

func (rb *RigidBody) SetPickable(value bool) {
	rb.IsPickable = value
}

func (rb *RigidBody) GetIdentifier() string {
	return rb.Identifier
}

func (rb *RigidBody) GetCanPick() bool {
	return rb.CanPick
}

func (rb *RigidBody) SetCanPick(value bool) {
	rb.CanPick = value
}

func (rb *RigidBody) SetCollidable(b bool) {
	rb.IsCollidable = b
}

func (rb *RigidBody) SetPushable(b bool) {
	rb.IsPushable = b
}

// ApplyForce applies a force to the rigid body.
func (rb *RigidBody) ApplyForce(force interfaces.Point) {
	if rb.IsStatic {
		return
	}
	rb.Acceleration.X += force.X / rb.Mass
	rb.Acceleration.Y += force.Y / rb.Mass
}

// Update updates the position of the RigidBody based on its velocity and the elapsed time.
func (rb *RigidBody) Update(deltaTime float64) {
	if rb.IsStatic {
		return
	}

	// Apply gravity if the rigid body is not on the ground
	if !rb.OnGround {
		rb.Velocity.Y += GRAVITY * deltaTime
	}

	// Update the position based on the velocity
	rb.Position.X += rb.Velocity.X * deltaTime
	rb.Position.Y += rb.Velocity.Y * deltaTime
}

// GetPosition returns the current position of the rigid body.
func (rb *RigidBody) GetPosition() interfaces.Point {
	return rb.Position
}

// SetPosition sets the position of the rigid body.
func (rb *RigidBody) SetPosition(position interfaces.Point) {
	rb.Position = position
}

// GetVelocity returns the current velocity of the rigid body.
func (rb *RigidBody) GetVelocity() interfaces.Point {
	return rb.Velocity
}

// SetVelocity sets the velocity of the rigid body.
func (rb *RigidBody) SetVelocity(velocity interfaces.Point) {
	rb.Velocity = velocity
}

// GetSize returns the size of the rigid body.
func (rb *RigidBody) GetSize() interfaces.Point {
	return rb.Size
}

// SetSize sets the size of the rigid body.
func (rb *RigidBody) SetSize(size interfaces.Point) {
	rb.Size = size
}

// Serialize converts the RigidBody into a map representation.
func (rb *RigidBody) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"position": map[string]float64{
			"x": rb.Position.X,
			"y": rb.Position.Y,
		},
		"size": map[string]float64{
			"x": rb.Size.X,
			"y": rb.Size.Y,
		},
		"velocity": map[string]float64{
			"x": rb.Velocity.X,
			"y": rb.Velocity.Y,
		},
		"mass":     rb.Mass,
		"isStatic": rb.IsStatic,
	}
}

// Deserialize populates the RigidBody with data from a map.
func (rb *RigidBody) Deserialize(data map[string]interface{}) {
	if position, ok := data["position"].(map[string]interface{}); ok {
		if x, ok := position["x"].(float64); ok {
			rb.Position.X = x
		}
		if y, ok := position["y"].(float64); ok {
			rb.Position.Y = y
		}
	}

	if size, ok := data["size"].(map[string]interface{}); ok {
		if x, ok := size["x"].(float64); ok {
			rb.Size.X = x
		}
		if y, ok := size["y"].(float64); ok {
			rb.Size.Y = y
		}
	}

	if velocity, ok := data["velocity"].(map[string]interface{}); ok {
		if x, ok := velocity["x"].(float64); ok {
			rb.Velocity.X = x
		}
		if y, ok := velocity["y"].(float64); ok {
			rb.Velocity.Y = y
		}
	}

	if mass, ok := data["mass"].(float64); ok {
		rb.Mass = mass
	}

	if isStatic, ok := data["isStatic"].(bool); ok {
		rb.IsStatic = isStatic
	}
}
