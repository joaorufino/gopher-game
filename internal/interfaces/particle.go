package interfaces

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// ParticleSystem defines the methods for handling particle effects.
type ParticleSystem interface {
	// AddParticle adds a particle to the system.
	AddParticle(position Vector2D, velocity Vector2D, lifespan float64, size int, color color.Color)
	// Update updates the state of the particle system.
	Update(deltaTime float64)
	// Draw renders the particle system on the screen.
	Draw(screen *ebiten.Image, camera Camera)
}

// Vector2D represents a 2D point.
type Vector2D struct {
	X, Y float64
}

// Rect represents a rectangle.
type Rect struct {
	Position Vector2D `json:"position"`
	Size     Vector2D `json:"size"`
}
