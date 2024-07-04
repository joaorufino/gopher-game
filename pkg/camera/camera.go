package camera

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Camera represents the game camera.
type Camera struct {
	x, y         float64
	screenWidth  int
	screenHeight int
}

// NewCamera initializes and returns a new Camera instance.
func NewCamera(screenWidth, screenHeight int) *Camera {
	return &Camera{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

// Follow updates the camera position to follow the target.
func (c *Camera) Follow(targetX, targetY float64) {
	c.x = targetX - float64(c.screenWidth)/2
	c.y = targetY - float64(c.screenHeight)/2
}

// Apply applies the camera's translation to the given DrawImageOptions.
func (c *Camera) Apply(opts *ebiten.DrawImageOptions) {
	opts.GeoM.Translate(-c.x, -c.y)
}

func (c *Camera) GetOffset() (float64, float64) {
	return c.x, c.y
}
