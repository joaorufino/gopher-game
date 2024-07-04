package interfaces

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

// ResourceManager defines the methods for loading and managing resources.
type ResourceManager interface {
	// LoadImage loads an image from the specified path.
	LoadImage(path string) (*ebiten.Image, error)
	// LoadSound loads a sound from the specified path.
	LoadSound(name, path string, context *audio.Context) error
	// GetImage returns the image with the specified name.
	GetImage(name string) (*ebiten.Image, error)
	// GetSound returns the sound with the specified name.
	GetSound(name string) (*audio.Player, error)
	GetItem(name string) (Item, error)
}
