package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/joaorufino/gopher-game/internal/event"
	"github.com/joaorufino/gopher-game/internal/interfaces"
	"github.com/joaorufino/gopher-game/pkg/abilities"
	"github.com/joaorufino/gopher-game/pkg/actions"
	"github.com/joaorufino/gopher-game/pkg/camera"
	"github.com/joaorufino/gopher-game/pkg/gameAudio"
	"github.com/joaorufino/gopher-game/pkg/gameMap"
	"github.com/joaorufino/gopher-game/pkg/input"
	"github.com/joaorufino/gopher-game/pkg/items"
	"github.com/joaorufino/gopher-game/pkg/particle"
	"github.com/joaorufino/gopher-game/pkg/physics"
	"github.com/joaorufino/gopher-game/pkg/player"
	"github.com/joaorufino/gopher-game/pkg/resource"
	"github.com/joaorufino/gopher-game/pkg/settings"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
}

// Provide the configuration for the game
func provideConfiguration() *player.Configuration {
	return &player.Configuration{
		ScreenWidth:  800,
		ScreenHeight: 18000,
		Gravity:      1000,
		JumpVelocity: 100,
		RunVelocity:  200,
		ImageScale:   0.2,
	}
}

// Provide the InputHandler implementation
func provideInputHandler(em interfaces.EventManager) interfaces.InputHandler {
	return input.NewInputHandler(em)
}

// Provide the EventManager implementation
func provideEventManager() interfaces.EventManager {
	return event.NewEventManager()
}

// Provide the PhysicsEngine implementation
func providePhysicsEngine(eventManager interfaces.EventManager) interfaces.PhysicsEngine {
	return physics.NewPhysicsEngine(eventManager, interfaces.Vector2D{X: 0, Y: 9.8}, 3000)
}

// Provide the GameMap implementation
func provideGameMap(physicsEngine interfaces.PhysicsEngine, eventManager interfaces.EventManager, resourceManager interfaces.ResourceManager, platformGenerator *gameMap.PlatformGenerator) (interfaces.Map, error) {
	return gameMap.NewMap(eventManager, resourceManager, physicsEngine, platformGenerator), nil
}

// Provide the Camera implementation
func provideCamera(settings interfaces.Settings) interfaces.Camera {
	return camera.NewCamera(settings.GetScreenWidth(), settings.GetScreenHeight())
}

// Provide the Settings implementation
func provideSettings() interfaces.Settings {
	return settings.NewSettings()
}

// Provide the AudioManager implementation
func provideAudioManager(config gameAudio.AudioManagerConfig) interfaces.AudioManager {
	return gameAudio.NewAudioManager(config)
}

// Provide the ParticleSystem implementation
func provideParticleSystem() interfaces.ParticleSystem {
	maxParticles := 100
	return particle.NewParticleSystem(maxParticles)
}

// Provide the Background Image
func provideBackgroundImage() (*ebiten.Image, error) {
	return loadBackgroundImage("images/background.png")
}

// Provide the screen dimensions
func provideScreenDimensions(settings interfaces.Settings) (int, int) {
	return settings.GetScreenWidth(), settings.GetScreenHeight()
}

// Provide the Player implementation
func providePlayer(resourceManager interfaces.ResourceManager, inputHandler interfaces.InputHandler, config *player.Configuration, engine interfaces.PhysicsEngine, events interfaces.EventManager) interfaces.Player {
	startX := 200.0
	startY := 200.0
	return player.NewPlayer(startX, startY, resourceManager, config, engine, events, 800)
}

// Provide screen dimensions
func provideScreenWidth(config *player.Configuration) int {
	return config.ScreenWidth
}

func provideScreenHeight(config *player.Configuration) int {
	return config.ScreenHeight
}

func loadBackgroundImage(path string) (*ebiten.Image, error) {
	bg, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load background image: %w", err)
	}
	return bg, nil
}

func provideItemManager(abilitiesMgr interfaces.AbilitiesManager) interfaces.ItemManager {
	im := items.NewItemManager(abilitiesMgr)
	err := im.LoadItems("game/items.json")
	if err != nil {
		log.Fatalf("Failed to load items: %v", err)
	}
	return im
}

func provideAbilitiesManager(eventManager interfaces.EventManager) interfaces.AbilitiesManager {
	acm := actions.NewActionManager()
	actions.RegisterBasicActions(acm)
	am := abilities.NewAbilitiesManager(acm.GetActions(), eventManager)
	err := am.LoadAbilities("game/abilities.json")
	if err != nil {
		log.Fatalf("Failed to load abilities: %v", err)
	}
	return am
}

// Provide the ResourceManager implementation
func provideResourceManager(itemManager interfaces.ItemManager) interfaces.ResourceManager {
	return resource.NewResourceManager(itemManager)
}

// Provide the PlatformGenerator implementation
func providePlatformGenerator(config *player.Configuration, physicsEngine interfaces.PhysicsEngine) *gameMap.PlatformGenerator {
	platformGenConfig := gameMap.PlatformGeneratorConfig{
		MinPlatformDistance: 50,
		MaxPlatformDistance: 150,
		PlatformWidth:       100,
		PlatformHeight:      20,
		ScreenWidth:         float64(config.ScreenWidth),
		ScreenHeight:        float64(config.ScreenHeight),
	}
	return gameMap.NewPlatformGenerator(platformGenConfig, physicsEngine)
}
