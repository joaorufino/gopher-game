package gameAudio

import (
	"log"
	"os"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// AudioManager manages game audio, including sound effects and background music.
type AudioManager struct {
	mu         sync.Mutex
	context    *audio.Context
	sounds     map[string]*audio.Player
	bgm        *audio.Player
	sampleRate int
	onLoad     func(string)
	onError    func(error)
}

// AudioManagerConfig holds configuration options for creating an AudioManager.
type AudioManagerConfig struct {
	Context    *audio.Context
	SampleRate int
	OnLoad     func(string)
	OnError    func(error)
}

// NewAudioManager initializes a new AudioManager with the given configuration.
func NewAudioManager(config AudioManagerConfig) *AudioManager {
	return &AudioManager{
		context:    config.Context,
		sounds:     make(map[string]*audio.Player),
		sampleRate: config.SampleRate,
		onLoad:     config.OnLoad,
		onError:    config.OnError,
	}
}

// LoadSound loads a sound from the specified file and stores it in the manager.
func (am *AudioManager) LoadSound(name, path string) (*audio.Player, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	f, err := os.Open(path)
	if err != nil {
		am.handleError(err)
		return nil, err
	}
	defer f.Close()

	d, err := wav.DecodeWithSampleRate(am.sampleRate, f)
	if err != nil {
		am.handleError(err)
		return nil, err
	}

	p, err := am.context.NewPlayer(d)
	if err != nil {
		am.handleError(err)
		return nil, err
	}

	am.sounds[name] = p
	am.handleLoad(name)
	return p, nil
}

// PlaySound plays the sound identified by the given name.
func (am *AudioManager) PlaySound(name string) {
	am.mu.Lock()
	defer am.mu.Unlock()

	if sound, ok := am.sounds[name]; ok {
		if err := sound.Rewind(); err != nil {
			am.handleError(err)
			return
		}
		sound.Play()
	} else {
		log.Printf("sound %s not found", name)
	}
}

// LoadBGM loads background music from the specified file.
func (am *AudioManager) LoadBGM(path string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	f, err := os.Open(path)
	if err != nil {
		am.handleError(err)
		return err
	}
	defer f.Close()

	d, err := wav.DecodeWithSampleRate(am.sampleRate, f)
	if err != nil {
		am.handleError(err)
		return err
	}

	p, err := am.context.NewPlayer(d)
	if err != nil {
		am.handleError(err)
		return err
	}

	am.bgm = p
	am.handleLoad("BGM")
	return nil
}

// PlayBGM plays the loaded background music.
func (am *AudioManager) PlayBGM() {
	am.mu.Lock()
	defer am.mu.Unlock()

	if am.bgm != nil {
		am.bgm.Play()
	}
}

// StopBGM stops the currently playing background music.
func (am *AudioManager) StopBGM() {
	am.mu.Lock()
	defer am.mu.Unlock()

	if am.bgm != nil {
		am.bgm.Pause()
	}
}

// handleLoad is a helper method to call the onLoad callback if set.
func (am *AudioManager) handleLoad(name string) {
	if am.onLoad != nil {
		am.onLoad(name)
	}
}

// handleError is a helper method to call the onError callback if set.
func (am *AudioManager) handleError(err error) {
	if am.onError != nil {
		am.onError(err)
	} else {
		log.Printf("error: %v", err)
	}
}
