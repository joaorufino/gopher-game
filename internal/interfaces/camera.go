package interfaces

import "github.com/hajimehoshi/ebiten/v2"

// Camera defines the methods for managing the game camera.
type Camera interface {
	Follow(targetX, targetY float64)
	Apply(opts *ebiten.DrawImageOptions)
	GetOffset() (float64, float64)
}
