package interfaces

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Renderer is an interface for components that can render themselves.
type Renderer interface {
	Render(screen *ebiten.Image)
}
