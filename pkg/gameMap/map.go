package gameMap

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/joaorufino/gopher-game/internal/interfaces"
	"github.com/joaorufino/gopher-game/pkg/physics"
)

// Movement defines the movement properties for an obstacle.
type Movement struct {
	Type        string  `json:"type"`
	Distance    float64 `json:"distance"`
	Speed       float64 `json:"speed"`
	InitialPosX float64 `json:"-"`
	InitialPosY float64 `json:"-"`
}

// Obstacle represents an obstacle with potential movement.
type Obstacle struct {
	Type      string             `json:"type"`
	Movement  Movement           `json:"movement"`
	RigidBody *physics.RigidBody `json:"body"`
}

// Platform represents a static platform in the game.
type Platform struct {
	RigidBody *physics.RigidBody `json:"body"`
}

// ItemOnMap represents an item placed on the map.
type ItemOnMap struct {
	Name      string             `json:"name"`
	RigidBody *physics.RigidBody `json:"body"`
	Item      interfaces.Item    `json:"-"`
}

// PlatformGeneratorConfig holds configuration for generating platforms.
type PlatformGeneratorConfig struct {
	MinPlatformDistance float64
	MaxPlatformDistance float64
	PlatformWidth       float64
	PlatformHeight      float64
	ScreenWidth         float64
	ScreenHeight        float64
}

// PlatformGenerator generates platforms dynamically.
type PlatformGenerator struct {
	config        PlatformGeneratorConfig
	platforms     []Platform
	lastPlatformY float64
	physicsEngine interfaces.PhysicsEngine
}

func NewPlatformGenerator(config PlatformGeneratorConfig, physicsEngine interfaces.PhysicsEngine) *PlatformGenerator {
	rand.Seed(time.Now().UnixNano())
	return &PlatformGenerator{
		config:        config,
		platforms:     []Platform{},
		lastPlatformY: config.ScreenHeight,
		physicsEngine: physicsEngine,
	}
}

func (pg *PlatformGenerator) GenerateInitialPlatforms() {
	for y := pg.config.ScreenHeight; y > 0; y -= pg.randomDistance() {
		pg.addPlatform(y)
	}
}

func (pg *PlatformGenerator) Update(deltaTime float64) {
	// Generate new platforms as the player moves up
	for !(pg.lastPlatformY > 0) {
		pg.addPlatform(pg.lastPlatformY)
	}
}

func (pg *PlatformGenerator) addPlatform(y float64) {
	x := rand.Float64() * (pg.config.ScreenWidth + pg.config.PlatformWidth)
	platform := Platform{
		RigidBody: physics.NewRigidBody(interfaces.Vector2D{X: x, Y: y}, interfaces.Vector2D{X: pg.config.PlatformWidth, Y: pg.config.PlatformHeight}, 1, true, "platform"),
	}
	pg.platforms = append(pg.platforms, platform)
	pg.lastPlatformY += pg.randomDistance()
	pg.physicsEngine.AddRigidBody(platform.RigidBody)
}

func (pg *PlatformGenerator) randomDistance() float64 {
	return pg.config.MinPlatformDistance + rand.Float64()*(pg.config.MaxPlatformDistance-pg.config.MinPlatformDistance)
}

func (pg *PlatformGenerator) GetPlatforms() []Platform {
	return pg.platforms
}

// Map represents the game map with platforms, obstacles, and items.
type Map struct {
	eventManager      interfaces.EventManager
	resourceManager   interfaces.ResourceManager
	platformGenerator *PlatformGenerator
	Obstacles         []Obstacle  `json:"obstacles"`
	Items             []ItemOnMap `json:"items"`
	Background        string      `json:"background"`
	BgImage           *ebiten.Image
}

// NewMap creates a new map instance.
func NewMap(eventManager interfaces.EventManager, resourceManager interfaces.ResourceManager, physicsEngine interfaces.PhysicsEngine, platformGenerator *PlatformGenerator) *Map {
	newMap := &Map{
		resourceManager:   resourceManager,
		eventManager:      eventManager,
		platformGenerator: platformGenerator,
	}
	newMap.eventManager.RegisterHandler(interfaces.EventItemEquipped, newMap.handleItemPicked)
	platformGenerator.GenerateInitialPlatforms()
	return newMap
}

func (m *Map) handleItemPicked(event interfaces.Event) {
	payload, ok := event.Payload.(map[string]interface{})
	if !ok {
		return
	}

	itemName, ok := payload["itemName"].(string)
	if !ok {
		return
	}

	m.removeItem(itemName)
	log.Printf("Item removed: %s", itemName)
}

func (m *Map) removeItem(itemName string) {
	for i, item := range m.Items {
		if item.Name == itemName {
			m.Items = append(m.Items[:i], m.Items[i+1:]...)
			return
		}
	}
}

// LoadBackground loads the background image.
func (m *Map) LoadBackground(imagePath string) error {
	bgImage, _, err := ebitenutil.NewImageFromFile(imagePath)
	if err != nil {
		return fmt.Errorf("failed to load background image: %w", err)
	}
	m.BgImage = bgImage
	return nil
}

func (m *Map) Update(deltaTime float64) {
	m.platformGenerator.Update(deltaTime)
	for i := range m.Obstacles {
		obstacle := &m.Obstacles[i]
		switch obstacle.Movement.Type {
		case "horizontal":
			obstacle.RigidBody.Position.X += obstacle.Movement.Speed * deltaTime
			if obstacle.RigidBody.Position.X > obstacle.Movement.InitialPosX+obstacle.Movement.Distance || obstacle.RigidBody.Position.X < obstacle.Movement.InitialPosX-obstacle.Movement.Distance {
				obstacle.Movement.Speed = -obstacle.Movement.Speed
			}
		case "vertical":
			obstacle.RigidBody.Position.Y += obstacle.Movement.Speed * deltaTime
			if obstacle.RigidBody.Position.Y > obstacle.Movement.InitialPosY+obstacle.Movement.Distance || obstacle.RigidBody.Position.Y < obstacle.Movement.InitialPosY-obstacle.Movement.Distance {
				obstacle.Movement.Speed = -obstacle.Movement.Speed
			}
		}
	}
}

func (m *Map) Draw(screen *ebiten.Image, camera interfaces.Camera) {
	// Get the offset from the camera
	offsetX, offsetY := camera.GetOffset()

	// Draw the background
	if m.BgImage != nil {
		bgOpts := &ebiten.DrawImageOptions{}
		bgOpts.GeoM.Translate(-offsetX, -offsetY)
		screen.DrawImage(m.BgImage, bgOpts)
	}

	// Draw platforms
	for _, platform := range m.platformGenerator.GetPlatforms() {
		vector.DrawFilledRect(screen,
			float32(platform.RigidBody.Position.X-offsetX),
			float32(platform.RigidBody.Position.Y-offsetY),
			float32(platform.RigidBody.Size.X),
			float32(platform.RigidBody.Size.Y),
			color.RGBA{0, 255, 0, 255},
			true)
	}

	// Draw obstacles
	for _, obstacle := range m.Obstacles {
		var cl color.RGBA
		switch obstacle.Type {
		case "docker_container":
			cl = color.RGBA{0, 0, 255, 255}
		case "docker_image":
			cl = color.RGBA{255, 0, 0, 255}
		default:
			cl = color.RGBA{255, 255, 0, 255}
		}
		vector.DrawFilledRect(screen,
			float32(obstacle.RigidBody.Position.X-offsetX),
			float32(obstacle.RigidBody.Position.Y-offsetY),
			float32(obstacle.RigidBody.Size.X),
			float32(obstacle.RigidBody.Size.Y),
			cl,
			true)
	}

	// Draw items
	for _, itemOnMap := range m.Items {
		itemOpts := &ebiten.DrawImageOptions{}
		itemOpts.GeoM.Translate(itemOnMap.RigidBody.Position.X-offsetX, itemOnMap.RigidBody.Position.Y-offsetY)
		iconImage, err := m.resourceManager.LoadImage(itemOnMap.Item.GetIconPath())
		if err != nil {
			log.Println(err)
		}
		screen.DrawImage(iconImage, itemOpts)
	}
}

// GetPlatforms returns the platforms from the map as a slice of interface{}.
func (m *Map) GetPlatforms() []interface{} {
	platforms := make([]interface{}, len(m.platformGenerator.GetPlatforms()))
	for i, platform := range m.platformGenerator.GetPlatforms() {
		platforms[i] = platform
	}
	return platforms
}

// SetPlatforms sets the platforms in the map.
func (m *Map) SetPlatforms(platforms []interface{}) {
	m.platformGenerator.platforms = make([]Platform, len(platforms))
	for i, platform := range platforms {
		if pla, ok := platform.(Platform); ok {
			m.platformGenerator.platforms[i] = pla
		}
	}
}

// GetObstacles returns the obstacles from the map as a slice of interface{}.
func (m *Map) GetObstacles() []interface{} {
	obstacles := make([]interface{}, len(m.Obstacles))
	for i, obstacle := range m.Obstacles {
		obstacles[i] = obstacle
	}
	return obstacles
}

// SetObstacles sets the obstacles in the map.
func (m *Map) SetObstacles(obstacles []interface{}) {
	m.Obstacles = make([]Obstacle, len(obstacles))
	for i, obstacle := range obstacles {
		if obs, ok := obstacle.(Obstacle); ok {
			m.Obstacles[i] = obs
		}
	}
}

// SetItems sets the items in the map.
func (m *Map) SetItems(items []ItemOnMap) {
	m.Items = items
}

// GetItems returns the items from the map.
func (m *Map) GetItems() []ItemOnMap {
	return m.Items
}
