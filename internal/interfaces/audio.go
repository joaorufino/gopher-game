package interfaces

import "github.com/hajimehoshi/ebiten/v2/audio"

// AudioManager defines the interface for managing game audio.
type AudioManager interface {
	// PlaySound plays a sound by its identifier.
	PlaySound(id string)

	// LoadSound loads a sound from the specified file and stores it in the manager.
	LoadSound(name, path string) (*audio.Player, error)

	// LoadBGM loads background music from the specified file.
	LoadBGM(path string) error

	// PlayBGM plays the loaded background music.
	PlayBGM()

	// StopBGM stops the currently playing background music.
	StopBGM()
}
