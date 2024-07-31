package settings

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"runtime"

	"github.com/joaorufino/gopher-game/internal/utils"
)

// SettingsImpl represents the game settings.
type SettingsImpl struct {
	settings map[string]interface{}
}

// NewSettings creates a new SettingsImpl instance.
func NewSettings() *SettingsImpl {
	return &SettingsImpl{
		settings: map[string]interface{}{
			"screenWidth":  800,
			"screenHeight": 600,
			"fullscreen":   false,
		}}
}

// Get retrieves a setting by key.
func (s *SettingsImpl) Get(key string) (interface{}, error) {
	value, exists := s.settings[key]
	if !exists {
		return nil, errors.New("setting not found: " + key)
	}
	return value, nil
}

// Set sets a value for a specific key.
func (s *SettingsImpl) Set(key string, value interface{}) {
	s.settings[key] = value
}

// Save saves the settings to a file.
// Uses conditional compilation to handle WebAssembly constraints.
// TODO: use a post method for updating
func (s *SettingsImpl) Save(path string) error {
	if runtime.GOARCH == "wasm" {
		return errors.New("saving settings not supported in WebAssembly")
	}

	data, err := json.Marshal(s.settings)
	if err != nil {
		log.Printf("Failed to marshal settings: %v", err)
		return err
	}
	err = os.WriteFile(path, data, 0600)
	if err != nil {
		log.Printf("Failed to write settings to file: %v", err)
		return err
	}
	return nil
}

// Load loads the settings from a file.
func (s *SettingsImpl) Load(path string) error {
	return utils.LoadData(path, func(data []byte) error {
		err := json.Unmarshal(data, &s.settings)
		if err != nil {
			log.Printf("Failed to unmarshal settings: %v", err)
			return err
		}
		log.Println("Settings loaded successfully")
		return nil
	})
}

// Convenience methods for specific settings
func (s *SettingsImpl) GetScreenWidth() int {
	if width, ok := s.settings["screenWidth"].(int); ok {
		return width
	}
	return 800
}

func (s *SettingsImpl) GetScreenHeight() int {
	if height, ok := s.settings["screenHeight"].(int); ok {
		return height
	}
	return 600
}

func (s *SettingsImpl) IsFullscreen() bool {
	if fullscreen, ok := s.settings["fullscreen"].(bool); ok {
		return fullscreen
	}
	return false
}
