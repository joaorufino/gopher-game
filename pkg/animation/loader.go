package animation

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joaorufino/cv-game/internal/interfaces"
)

func LoadAnimations(resourceManager interfaces.ResourceManager, basePath string, frameCounts map[string]int) (map[string]*Animation, interfaces.Point) {
	animations := make(map[string]*Animation)
	var frames []*ebiten.Image

	for k, v := range frameCounts {
		frames = loadFrames(resourceManager, basePath+"/"+k, v)
		animations[k] = NewAnimation(AnimationConfig{
			Frames:        frames,
			FrameDuration: 0.1,
			Loop:          true,
		})

	}
	size := interfaces.Point{}
	if len(frames) != 0 {
		// TODO: we are hoping that all images have the same size
		width, height := frames[0].Size()
		size.X, size.Y = float64(width), float64(height)
	}

	return animations, size
}

// loadFrames loads a series of frames from files.
func loadFrames(resourceManager interfaces.ResourceManager, basePath string, frameCount int) []*ebiten.Image {
	frames := []*ebiten.Image{}
	for i := 0; i < frameCount; i++ {
		imgPath := fmt.Sprintf("%s%d.png", basePath, i)
		img, err := resourceManager.LoadImage(imgPath)
		if err != nil {
			log.Fatalf("failed to load frame image: %v", err)
		}
		frames = append(frames, img)
	}
	return frames
}
