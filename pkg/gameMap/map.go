package gameMap

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/joaorufino/cv-game/internal/interfaces"
	"github.com/joaorufino/cv-game/internal/utils"
	"github.com/joaorufino/cv-game/pkg/physics"
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
	Item      interfaces.Item    `json:-`
}

// Map represents the game map with platforms, obstacles, and items.
type Map struct {
	eventManager    interfaces.EventManager
	resourceManager interfaces.ResourceManager
	Platforms       []Platform  `json:"platforms"`
	Obstacles       []Obstacle  `json:"obstacles"`
	Items           []ItemOnMap `json:"items"`
	Background      string      `json:"background"`
	BgImage         *ebiten.Image
}

// NewMap creates a new map instance.
func NewMap(filepath string, eventManager interfaces.EventManager, resourceManager interfaces.ResourceManager, physicsEngine interfaces.PhysicsEngine) (interfaces.Map, error) {
	newMap := &Map{
		resourceManager: resourceManager,
		eventManager:    eventManager,
	}
	newMap.eventManager.RegisterHandler(interfaces.EventItemEquipped, newMap.handleItemPicked)
	return newMap, newMap.LoadMap(filepath, physicsEngine)
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

// LoadMap loads the map from a JSON file.
func (m *Map) LoadMap(filepath string, physicsEngine interfaces.PhysicsEngine) error {
	err := utils.LoadData(filepath, func(data []byte) error {
		err := json.Unmarshal(data, &m)
		if err != nil {
			return fmt.Errorf("failed to unmarshal map JSON: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	bgImage, _, err := ebitenutil.NewImageFromFile(m.Background)
	if err != nil {
		return fmt.Errorf("failed to load background image: %w", err)
	}
	m.BgImage = bgImage

	// Initialize platforms with rigid bodies
	for i := range m.Platforms {
		platform := &m.Platforms[i]
		platform.RigidBody = physics.NewRigidBody(platform.RigidBody.Position, platform.RigidBody.Size, 1, true, fmt.Sprintf("platform%d", i))
		physicsEngine.AddRigidBody(platform.RigidBody)
	}

	// Initialize obstacles with rigid bodies and initial positions
	for i := range m.Obstacles {
		obstacle := &m.Obstacles[i]
		obstacle.RigidBody = physics.NewRigidBody(obstacle.RigidBody.Position, obstacle.RigidBody.Size, 1, true, fmt.Sprintf("obstacle%d", i))
		obstacle.Movement.InitialPosX = obstacle.RigidBody.Position.X
		obstacle.Movement.InitialPosY = obstacle.RigidBody.Position.Y
		physicsEngine.AddRigidBody(obstacle.RigidBody)
	}

	// Initialize items with rigid bodies
	for i, itemData := range m.Items {
		item, err := m.resourceManager.GetItem(itemData.Name)
		if err != nil {
			return fmt.Errorf("failed to get item: %w", err)
		}
		m.Items[i].Item = item
		m.Items[i].RigidBody = physics.NewRigidBody(itemData.RigidBody.Position, itemData.RigidBody.Size, 1, true, itemData.Name)
		m.Items[i].RigidBody.SetPickable(true)
		physicsEngine.AddRigidBody(m.Items[i].RigidBody)
	}

	log.Println("Map loaded successfully")
	return nil
}

func (m *Map) Update(deltaTime float64) {
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
	log.Println("Drawing map...")

	// Get the offset from the camera
	offsetX, offsetY := camera.GetOffset()

	// Draw the background
	if m.BgImage != nil {
		bgOpts := &ebiten.DrawImageOptions{}
		bgOpts.GeoM.Translate(-offsetX, -offsetY)
		screen.DrawImage(m.BgImage, bgOpts)
	} else {
		log.Println("Background image is nil") // Debug log
	}

	// Draw platforms
	for _, platform := range m.Platforms {
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

func (m *Map) SetBackground(imagePath string) {
	m.Background = imagePath
}

// GetPlatforms returns the platforms from the map as a slice of interface{}.
func (m *Map) GetPlatforms() []interface{} {
	platforms := make([]interface{}, len(m.Platforms))
	for i, platform := range m.Platforms {
		platforms[i] = platform
	}
	return platforms
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

// GetObstacles returns the obstacles from the map as a slice of interface{}.
func (m *Map) GetObstacles() []interface{} {
	obstacles := make([]interface{}, len(m.Obstacles))
	for i, obstacle := range m.Obstacles {
		obstacles[i] = obstacle
	}
	return obstacles
}

// SetPlatforms sets the obstacles in the map.
func (m *Map) SetPlatforms(platforms []interface{}) {
	m.Platforms = make([]Platform, len(platforms))
	for i, platform := range platforms {
		if pla, ok := platform.(Platform); ok {
			m.Platforms[i] = pla
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
