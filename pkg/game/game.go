package game

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/joaorufino/gopher-game/internal/interfaces"
	"github.com/joaorufino/gopher-game/pkg/abilities"
	"github.com/joaorufino/gopher-game/pkg/achievements"
	"github.com/joaorufino/gopher-game/pkg/actions"
	"github.com/joaorufino/gopher-game/pkg/chapterintro"
	"github.com/joaorufino/gopher-game/pkg/gameMap"
	"github.com/joaorufino/gopher-game/pkg/hud"
	"github.com/joaorufino/gopher-game/pkg/pet"
	"github.com/joaorufino/gopher-game/pkg/physics"
	"github.com/joaorufino/gopher-game/pkg/score"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

// Params defines the dependencies for the Game struct.
type Params struct {
	fx.In

	Player          interfaces.Player
	Background      *ebiten.Image
	ScreenWidth     int `name:"screenWidth"`
	ScreenHeight    int `name:"screenHeight"`
	GameMap         interfaces.Map
	Settings        interfaces.Settings
	ResourceManager interfaces.ResourceManager
	ItemManager     interfaces.ItemManager
	Camera          interfaces.Camera
	PhysicsEngine   interfaces.PhysicsEngine
	EventManager    interfaces.EventManager
	InputHandler    interfaces.InputHandler
}

// Game represents the main game structure.
type Game struct {
	Player             interfaces.Player
	Background         *ebiten.Image
	ScreenWidth        int
	ScreenHeight       int
	GameMap            interfaces.Map
	Settings           interfaces.Settings
	ResourceManager    interfaces.ResourceManager
	ItemManager        interfaces.ItemManager
	Camera             interfaces.Camera
	PhysicsEngine      interfaces.PhysicsEngine
	EventManager       interfaces.EventManager
	InputHandler       interfaces.InputHandler
	Pet                *pet.Pet
	AbilitiesManager   interfaces.AbilitiesManager
	AchievementManager interfaces.AchievementManager
	chapterIntro       *chapterintro.ChapterIntro
	ScoreManager       *score.ScoreManager
	HUD                *hud.HUD
	matchTimer         float64 // Timer for soccer match in seconds
}

// NewGame creates a new Game instance using dependency injection.
func NewGame(params Params) *Game {
	fullscreen, err := params.Settings.Get("fullscreen")
	if err != nil {
		log.Fatalf("Failed to get fullscreen: %v", err)
	}

	ebiten.SetWindowSize(params.ScreenWidth, params.ScreenHeight)
	if fullscreen.(bool) {
		ebiten.SetFullscreen(true)
	}

	// Initialize the player and pet
	player := params.Player
	petConfig := &pet.Configuration{ImageScale: 0.1, RunVelocity: 50, JumpVelocity: 20}
	petInstance := pet.NewPet(player.GetPosition().X, player.GetPosition().Y, params.ResourceManager, petConfig, params.PhysicsEngine, player)

	achievementConfig := achievements.Config{
		ScreenWidth:            800,
		ScreenHeight:           600,
		DisplayDuration:        5 * time.Second,
		MaxAchievementsDisplay: 3,
		PaddingX:               10,
		PaddingY:               10,
		AchievementOffsetY:     40,
		TextOffsetX:            50,
		TextOffsetY:            10,
	}
	actionManager := actions.NewActionManager()
	actions.RegisterBasicActions(actionManager)
	achievementManager := achievements.NewAchievementManager(achievementConfig, params.EventManager)
	abilitiesManager := abilities.NewAbilitiesManager(actionManager.GetActions(), params.EventManager)
	// Load abilities from a JSON file
	err = abilitiesManager.LoadAbilities("game/abilities.json")
	if err != nil {
		log.Fatalf("Failed to load abilities: %v", err)
	}
	chapterIntro := chapterintro.NewChapterIntro("Soccer Match - Score Goals to Win!", interfaces.Vector2D{X: 100, Y: 400}, params.PhysicsEngine)

	// Configure the PlatformGenerator (now used for soccer players)
	platformGenConfig := gameMap.PlatformGeneratorConfig{
		MinPlatformDistance: 100,
		MaxPlatformDistance: 150,
		PlatformWidth:       30, // Smaller for player representations
		PlatformHeight:      30, // Square for player representations
		ScreenWidth:         float64(params.ScreenWidth),
		ScreenHeight:        float64(params.ScreenHeight),
	}
	platformGenerator := gameMap.NewPlatformGenerator(platformGenConfig, params.PhysicsEngine)

	// Initialize the game map (soccer field) with the PlatformGenerator
	gameMapInstance := gameMap.NewMap(params.EventManager, params.ResourceManager, params.PhysicsEngine, platformGenerator)

	// Create score manager for soccer game
	scoreManager := score.NewScoreManager()
	scoreManager.SetTeamName(0, "Blue Team")
	scoreManager.SetTeamName(1, "Red Team")
	scoreManager.StartMatch(90) // 90 second match

	// Create HUD to display scores
	// Note: This part assumes we have font loading - if not, this can be adjusted
	hud := hud.NewHUD(scoreManager, nil, params.ScreenWidth, params.ScreenHeight)

	game := &Game{
		Player:             player,
		Pet:                petInstance,
		Background:         params.Background,
		ScreenWidth:        params.ScreenWidth,
		ScreenHeight:       params.ScreenHeight,
		GameMap:            gameMapInstance,
		Settings:           params.Settings,
		ResourceManager:    params.ResourceManager,
		ItemManager:        params.ItemManager,
		Camera:             params.Camera,
		PhysicsEngine:      params.PhysicsEngine,
		EventManager:       params.EventManager,
		InputHandler:       params.InputHandler,
		AbilitiesManager:   abilitiesManager,
		AchievementManager: achievementManager,
		chapterIntro:       chapterIntro,
		ScoreManager:       scoreManager,
		HUD:                hud,
	}

	game.registerEventHandlers()

	return game
}

func loadBackgroundImage(path string) (*ebiten.Image, error) {
	bg, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load background image: %w", err)
	}
	return bg, nil
}

// registerEventHandlers registers event handlers for the game.
func (g *Game) registerEventHandlers() {
	g.EventManager.RegisterHandler(interfaces.EventPlayerJump, func(event interfaces.Event) {
		logrus.Info("Player jumped:", event.Payload)
	})
	g.EventManager.RegisterHandler(interfaces.EventPlayerMove, func(event interfaces.Event) {
		logrus.Info("Player moved:", event.Payload)
	})
	g.EventManager.RegisterHandler(interfaces.EventItemEquipped, func(event interfaces.Event) {
		logrus.Info("Item equipped:", event.Payload)
	})
}

// Update updates the game state.
func (g *Game) Update() error {
	deltaTime := 1.0 / 60.0 // Fixed time step for consistent physics

	g.chapterIntro.Update(deltaTime)
	// Update the input handler
	if err := g.InputHandler.Update(); err != nil {
		return err
	}

	if err := g.Player.Update(deltaTime); err != nil {
		return err
	}
	if err := g.Pet.Update(deltaTime); err != nil {
		return err
	}
	g.GameMap.Update(deltaTime)

	// Update camera to follow the player
	g.Camera.Follow(g.Player.GetPosition().X, g.Player.GetPosition().Y)

	g.PhysicsEngine.Update(deltaTime)
	g.AchievementManager.Update()
	g.AbilitiesManager.Update(deltaTime)

	// Update match timer and score for soccer game
	if g.ScoreManager.IsMatchActive() {
		g.matchTimer += deltaTime
		// Update every second
		if g.matchTimer >= 1.0 {
			g.ScoreManager.UpdateMatchTime(1)
			g.matchTimer = 0
		}

		// Check for goal scoring - detect if ball enters goal areas
		for _, item := range g.GameMap.GetItems() {
			// Need to cast the RigidBody interface to the concrete type
			if rb, ok := item.RigidBody.(*physics.RigidBody); ok {
				// Check if item is in left goal area
				if rb.Position.X <= 20 &&
					rb.Position.Y >= 225 &&
					rb.Position.Y <= 375 {
					// Score for right team (away/red)
					g.ScoreManager.AddGoal(1)
					// Reset ball position
					rb.Position.X = 400
					rb.Position.Y = 300
					rb.Velocity.X = 0
					rb.Velocity.Y = 0
				}

				// Check if item is in right goal area
				if rb.Position.X >= 780 &&
					rb.Position.Y >= 225 &&
					rb.Position.Y <= 375 {
					// Score for left team (home/blue)
					g.ScoreManager.AddGoal(0)
					// Reset ball position
					rb.Position.X = 400
					rb.Position.Y = 300
					rb.Velocity.X = 0
					rb.Velocity.Y = 0
				}
			}
		}
	}

	return nil
}

// Draw draws the game on the screen.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	options := &ebiten.DrawImageOptions{}
	g.Camera.Apply(options)
	g.GameMap.Draw(screen, g.Camera)
	g.chapterIntro.Draw(screen, g.Camera)
	g.AchievementManager.Draw(screen)

	if err := g.Player.Draw(screen, g.Camera); err != nil {
		log.Printf("could not draw player %v", err)
	}
	if err := g.Pet.Draw(screen, g.Camera); err != nil {
		log.Printf("could not draw pet %v", err)
	}

	// Draw the HUD with score information
	if g.HUD != nil {
		g.HUD.Draw(screen)
	}
}

// Layout sets the screen layout dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.ScreenWidth = outsideWidth
	g.ScreenHeight = outsideHeight
	return g.ScreenWidth, g.ScreenHeight
}
