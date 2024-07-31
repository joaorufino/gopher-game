package resource

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/joaorufino/gopher-game/internal/interfaces"
	"github.com/joaorufino/gopher-game/pkg/gameAudio"
)

// ResourceManagerImpl manages game resources.
type ResourceManagerImpl struct {
	images      map[string]*ebiten.Image
	sounds      map[string]*audio.Player
	cache       map[string]*ebiten.Image
	itemManager interfaces.ItemManager
}

// NewResourceManager initializes a new resource manager.
func NewResourceManager(itemManager interfaces.ItemManager) interfaces.ResourceManager {
	return &ResourceManagerImpl{
		images:      make(map[string]*ebiten.Image),
		sounds:      make(map[string]*audio.Player),
		cache:       make(map[string]*ebiten.Image),
		itemManager: itemManager,
	}
}

// LoadImage loads an image from a file path, using cache if available.
func (rm *ResourceManagerImpl) LoadImage(filePath string) (*ebiten.Image, error) {
	if img, ok := rm.cache[filePath]; ok {
		return img, nil
	}

	img, _, err := ebitenutil.NewImageFromFile(filePath)
	if err != nil {
		return nil, err
	}

	rm.cache[filePath] = img
	return img, nil
}

// GetImage retrieves an image from the manager.
func (rm *ResourceManagerImpl) GetImage(name string) (*ebiten.Image, error) {
	img, ok := rm.images[name]
	if !ok {
		return nil, errors.New("image not found")
	}
	return img, nil
}

// LoadSound loads a sound from file and stores it in the manager.
func (rm *ResourceManagerImpl) LoadSound(name, path string, context *audio.Context) error {
	am := gameAudio.NewAudioManager(gameAudio.AudioManagerConfig{
		Context:    context,
		SampleRate: 44100,
		OnLoad:     nil,
		OnError:    nil,
	})
	p, err := am.LoadSound(name, path)
	if err != nil {
		return err
	}
	rm.sounds[name] = p
	return nil
}

// GetSound retrieves a sound from the manager.
func (rm *ResourceManagerImpl) GetSound(name string) (*audio.Player, error) {
	sound, ok := rm.sounds[name]
	if !ok {
		return nil, errors.New("sound not found")
	}
	return sound, nil
}

// GetItem retrieves an item from the manager.
func (rm *ResourceManagerImpl) GetItem(name string) (interfaces.Item, error) {
	return rm.itemManager.GetItem(name)
}
