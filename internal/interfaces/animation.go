package interfaces

import "github.com/hajimehoshi/ebiten/v2"

// Animation defines the methods for handling animations.
type Animation interface {
	// Update updates the state of the animation.
	Update()
	// Draw renders the animation on the screen.
	// screen: The screen to draw the animation on.
	// options: The options for drawing the image.
	Draw(screen *ebiten.Image, options *ebiten.DrawImageOptions)
}
