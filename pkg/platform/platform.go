package platform

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joaorufino/gopher-game/internal/interfaces"
)

type Platform struct {
	X, Y   float64
	Width  float64
	Height float64
	Type   string // e.g., "static", "moving"
}

type PlatformManager struct {
	Platforms    []Platform
	ScreenWidth  int
	ScreenHeight int
}

func NewPlatformManager(screenWidth, screenHeight int) *PlatformManager {
	rand.Seed(time.Now().UnixNano())
	return &PlatformManager{
		Platforms:    []Platform{},
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
	}
}

func (pm *PlatformManager) GenerateInitialPlatforms() {
	// Generate initial platforms at the bottom of the screen
	for i := 0; i < 10; i++ {
		pm.Platforms = append(pm.Platforms, Platform{
			X:      rand.Float64() * float64(pm.ScreenWidth),
			Y:      float64(pm.ScreenHeight - (i * 60)),
			Width:  100,
			Height: 20,
			Type:   "static",
		})
	}
}

func (pm *PlatformManager) Update(deltaTime float64) {
	// Update platform positions or types if necessary
	for i := range pm.Platforms {
		// Example: Move platforms down as player climbs
		pm.Platforms[i].Y += 100 * deltaTime
		// Remove platforms that are out of view
		if pm.Platforms[i].Y > float64(pm.ScreenHeight) {
			pm.Platforms = append(pm.Platforms[:i], pm.Platforms[i+1:]...)
			i--
		}
	}
	// Generate new platforms above the view
	if len(pm.Platforms) > 0 && pm.Platforms[len(pm.Platforms)-1].Y < float64(pm.ScreenHeight)-200 {
		pm.Platforms = append(pm.Platforms, Platform{
			X:      rand.Float64() * float64(pm.ScreenWidth),
			Y:      pm.Platforms[len(pm.Platforms)-1].Y - 60,
			Width:  100,
			Height: 20,
			Type:   "static",
		})
	}
}

func (pm *PlatformManager) Draw(screen *ebiten.Image, camera interfaces.Camera) {
	// Draw platforms
	for _, platform := range pm.Platforms {
		vector.DrawFilledRect(screen,
			float32(platform.X),
			float32(platform.Y),
			float32(platform.Width),
			float32(platform.Height),
			color.RGBA{0, 255, 0, 255},
			true)
	}
}
