package gameMap

import (
	"fmt"
	"image/color"
	"log"
	"math"
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
// In soccer game context, this generates player positions
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
	// For soccer game, place players in formation
	
	// Generate blue team players
	// Goalkeeper
	pg.addPlayerPlatform(50, 300, "goalkeeper")
	
	// Defenders
	pg.addPlayerPlatform(150, 150, "defender")
	pg.addPlayerPlatform(150, 300, "defender")
	pg.addPlayerPlatform(150, 450, "defender")
	
	// Midfielders
	pg.addPlayerPlatform(300, 150, "midfielder")
	pg.addPlayerPlatform(300, 300, "midfielder")
	pg.addPlayerPlatform(300, 450, "midfielder")
	
	// Strikers
	pg.addPlayerPlatform(450, 200, "striker")
	pg.addPlayerPlatform(450, 400, "striker")
	
	// Generate red team players (as obstacles) - done in the Map class
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

// For soccer game - position player at specific x,y with role
func (pg *PlatformGenerator) addPlayerPlatform(x float64, y float64, role string) {
	platform := Platform{
		RigidBody: physics.NewRigidBody(interfaces.Vector2D{X: x, Y: y}, interfaces.Vector2D{X: pg.config.PlatformWidth, Y: pg.config.PlatformHeight}, 1, true, role),
	}
	pg.platforms = append(pg.platforms, platform)
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
	
	// Add soccer ball as an item
	ballSize := interfaces.Vector2D{X: 20, Y: 20}
	ballRigidBody := physics.NewRigidBody(interfaces.Vector2D{X: 400, Y: 300}, ballSize, 0.5, false, "soccer_ball")
	ballRigidBody.SetCanPick(true)
	// Make the ball move more naturally - lower mass and friction
	ballRigidBody.Mass = 0.2
	ballRigidBody.Friction = 0.98
	
	// Add the ball to the physics engine
	physicsEngine.AddRigidBody(ballRigidBody)
	
	// Create a dummy item for the ball
	dummyItem := &SoccerBall{Name: "soccer_ball"}
	
	// Add the ball to the map's items
	newMap.Items = append(newMap.Items, ItemOnMap{
		Name:      "soccer_ball",
		RigidBody: ballRigidBody,
		Item:      dummyItem,
	})
	
	// Add red team players as obstacles
	newMap.addRedTeamPlayers(physicsEngine)
	
	// Add boundary walls to keep everything inside the field
	newMap.addBoundaryWalls(physicsEngine)
	
	return newMap
}

// Add boundary walls to keep players and ball inside the field
func (m *Map) addBoundaryWalls(physicsEngine interfaces.PhysicsEngine) {
	// Field dimensions
	fieldWidth := 800.0
	fieldHeight := 600.0
	wallThickness := 20.0
	
	// Top wall
	topWall := physics.NewRigidBody(
		interfaces.Vector2D{X: 0, Y: -wallThickness},
		interfaces.Vector2D{X: fieldWidth, Y: wallThickness},
		100, true, "wall_top",
	)
	physicsEngine.AddRigidBody(topWall)
	
	// Bottom wall
	bottomWall := physics.NewRigidBody(
		interfaces.Vector2D{X: 0, Y: fieldHeight},
		interfaces.Vector2D{X: fieldWidth, Y: wallThickness},
		100, true, "wall_bottom",
	)
	physicsEngine.AddRigidBody(bottomWall)
	
	// Left wall (except goal area)
	leftWallTop := physics.NewRigidBody(
		interfaces.Vector2D{X: -wallThickness, Y: 0},
		interfaces.Vector2D{X: wallThickness, Y: fieldHeight/2 - 75},
		100, true, "wall_left_top",
	)
	physicsEngine.AddRigidBody(leftWallTop)
	
	leftWallBottom := physics.NewRigidBody(
		interfaces.Vector2D{X: -wallThickness, Y: fieldHeight/2 + 75},
		interfaces.Vector2D{X: wallThickness, Y: fieldHeight/2 - 75},
		100, true, "wall_left_bottom",
	)
	physicsEngine.AddRigidBody(leftWallBottom)
	
	// Right wall (except goal area)
	rightWallTop := physics.NewRigidBody(
		interfaces.Vector2D{X: fieldWidth, Y: 0},
		interfaces.Vector2D{X: wallThickness, Y: fieldHeight/2 - 75},
		100, true, "wall_right_top",
	)
	physicsEngine.AddRigidBody(rightWallTop)
	
	rightWallBottom := physics.NewRigidBody(
		interfaces.Vector2D{X: fieldWidth, Y: fieldHeight/2 + 75},
		interfaces.Vector2D{X: wallThickness, Y: fieldHeight/2 - 75},
		100, true, "wall_right_bottom",
	)
	physicsEngine.AddRigidBody(rightWallBottom)
}

// SoccerBall implements the interfaces.Item interface
type SoccerBall struct {
	Name string
}

func (sb *SoccerBall) GetName() string {
	return sb.Name
}

func (sb *SoccerBall) GetDescription() string {
	return "A soccer ball"
}

func (sb *SoccerBall) GetIconPath() string {
	return "/images/icons/ball.png" // This path might not exist, handled in Draw method
}

func (sb *SoccerBall) GetImagePath() string {
	return "/images/icons/ball.png" // This path might not exist, handled in Draw method
}

func (sb *SoccerBall) GetAbilities() []string {
	return []string{}
}

func (sb *SoccerBall) GetAppearance() interfaces.Appearance {
	return interfaces.Appearance{
		Type:     "ball",
		Color:    "black and white",
		Material: "leather",
	}
}

func (sb *SoccerBall) GetVersion() int {
	return 1
}

// Add red team players as obstacles
func (m *Map) addRedTeamPlayers(physicsEngine interfaces.PhysicsEngine) {
	// Goalkeeper
	m.addOpposingPlayer(750, 300, 30, 30, "goalkeeper", physicsEngine)
	
	// Defenders
	m.addOpposingPlayer(650, 150, 30, 30, "defender", physicsEngine)
	m.addOpposingPlayer(650, 300, 30, 30, "defender", physicsEngine)
	m.addOpposingPlayer(650, 450, 30, 30, "defender", physicsEngine)
	
	// Midfielders
	m.addOpposingPlayer(500, 150, 30, 30, "midfielder", physicsEngine)
	m.addOpposingPlayer(500, 300, 30, 30, "midfielder", physicsEngine)
	m.addOpposingPlayer(500, 450, 30, 30, "midfielder", physicsEngine)
	
	// Strikers
	m.addOpposingPlayer(350, 200, 30, 30, "striker", physicsEngine)
	m.addOpposingPlayer(350, 400, 30, 30, "striker", physicsEngine)
}

// Helper to add opposing team players
func (m *Map) addOpposingPlayer(x, y, width, height float64, role string, physicsEngine interfaces.PhysicsEngine) {
	obstacle := Obstacle{
		Type: role,
		RigidBody: physics.NewRigidBody(
			interfaces.Vector2D{X: x, Y: y},
			interfaces.Vector2D{X: width, Y: height},
			1, true, "red_"+role,
		),
	}
	
	m.Obstacles = append(m.Obstacles, obstacle)
	physicsEngine.AddRigidBody(obstacle.RigidBody)
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
	
	// Update obstacles (red team players)
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
	
	// Make the ball move toward the center of the field when not being pushed
	for _, item := range m.Items {
		if item.Name == "soccer_ball" {
			// Get center coordinates
			centerX := 400.0
			centerY := 300.0
			
			// Calculate vector toward center
			dirX := centerX - item.RigidBody.Position.X
			dirY := centerY - item.RigidBody.Position.Y
			
			// Calculate distance to center
			distance := math.Sqrt(dirX*dirX + dirY*dirY)
			
			// Only apply force if the ball is not at the center and moving slowly
			if distance > 5.0 && math.Abs(item.RigidBody.Velocity.X) < 50 && math.Abs(item.RigidBody.Velocity.Y) < 50 {
				// Normalize direction vector
				if distance > 0 {
					dirX /= distance
					dirY /= distance
				}
				
				// Apply a gentle force toward the center
				forceStrength := 10.0 * deltaTime
				item.RigidBody.ApplyForce(interfaces.Vector2D{
					X: dirX * forceStrength,
					Y: dirY * forceStrength,
				})
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

	// Draw soccer field
	// Field background (green)
	fieldWidth := 800.0
	fieldHeight := 600.0
	fieldX := -offsetX
	fieldY := -offsetY
	
	// Draw the green field
	vector.DrawFilledRect(screen,
		float32(fieldX),
		float32(fieldY),
		float32(fieldWidth),
		float32(fieldHeight),
		color.RGBA{34, 139, 34, 255}, // Forest Green
		true)
	
	// Draw field lines (white)
	// Center line
	vector.DrawFilledRect(screen,
		float32(fieldX+fieldWidth/2-2),
		float32(fieldY),
		4,
		float32(fieldHeight),
		color.RGBA{255, 255, 255, 255},
		true)
	
	// Center circle
	centerX := float32(fieldX + fieldWidth/2)
	centerY := float32(fieldY + fieldHeight/2)
	radius := float32(50)
	segments := 30
	for i := 0; i < segments; i++ {
		angle1 := float32(i) * 2 * 3.14159 / float32(segments)
		angle2 := float32(i+1) * 2 * 3.14159 / float32(segments)
		x1 := centerX + radius*float32(math.Cos(float64(angle1)))
		y1 := centerY + radius*float32(math.Sin(float64(angle1)))
		x2 := centerX + radius*float32(math.Cos(float64(angle2)))
		y2 := centerY + radius*float32(math.Sin(float64(angle2)))
		vector.StrokeLine(screen, x1, y1, x2, y2, 2, color.RGBA{255, 255, 255, 255}, true)
	}
	
	// Goal areas
	// Left goal
	vector.DrawFilledRect(screen,
		float32(fieldX),
		float32(fieldY+fieldHeight/2-75),
		20,
		150,
		color.RGBA{200, 200, 200, 255},
		true)
	
	// Right goal
	vector.DrawFilledRect(screen,
		float32(fieldX+fieldWidth-20),
		float32(fieldY+fieldHeight/2-75),
		20,
		150,
		color.RGBA{200, 200, 200, 255},
		true)
	
	// Draw platforms (as players or obstacles)
	for _, platform := range m.platformGenerator.GetPlatforms() {
		vector.DrawFilledRect(screen,
			float32(platform.RigidBody.Position.X-offsetX),
			float32(platform.RigidBody.Position.Y-offsetY),
			float32(platform.RigidBody.Size.X),
			float32(platform.RigidBody.Size.Y),
			color.RGBA{0, 0, 255, 255}, // Blue for players
			true)
	}
	
	// Draw obstacles (as opposing team players)
	for _, obstacle := range m.Obstacles {
		var cl color.RGBA
		switch obstacle.Type {
		case "goalkeeper":
			cl = color.RGBA{255, 0, 0, 255} // Red for goalkeepers
		case "defender":
			cl = color.RGBA{255, 100, 100, 255} // Light red for defenders
		case "striker":
			cl = color.RGBA{200, 0, 0, 255} // Dark red for strikers
		default:
			cl = color.RGBA{255, 0, 0, 255} // Red for default players
		}
		vector.DrawFilledRect(screen,
			float32(obstacle.RigidBody.Position.X-offsetX),
			float32(obstacle.RigidBody.Position.Y-offsetY),
			float32(obstacle.RigidBody.Size.X),
			float32(obstacle.RigidBody.Size.Y),
			cl,
			true)
	}
	
	// Draw items (as soccer ball)
	for _, itemOnMap := range m.Items {
		itemOpts := &ebiten.DrawImageOptions{}
		itemOpts.GeoM.Translate(itemOnMap.RigidBody.Position.X-offsetX, itemOnMap.RigidBody.Position.Y-offsetY)
		
		// Try to load the image first
		iconImage, err := m.resourceManager.LoadImage(itemOnMap.Item.GetIconPath())
		if err != nil {
			log.Println(err)
			// If image loading fails, draw a simple soccer ball
			ballRadius := float32(10)
			ballCenterX := float32(itemOnMap.RigidBody.Position.X - offsetX + itemOnMap.RigidBody.Size.X/2)
			ballCenterY := float32(itemOnMap.RigidBody.Position.Y - offsetY + itemOnMap.RigidBody.Size.Y/2)
			
			// Draw white circle
			segments := 20
			for i := 0; i < segments; i++ {
				angle1 := float32(i) * 2 * 3.14159 / float32(segments)
				angle2 := float32(i+1) * 2 * 3.14159 / float32(segments)
				x1 := ballCenterX + ballRadius*float32(math.Cos(float64(angle1)))
				y1 := ballCenterY + ballRadius*float32(math.Sin(float64(angle1)))
				x2 := ballCenterX + ballRadius*float32(math.Cos(float64(angle2)))
				y2 := ballCenterY + ballRadius*float32(math.Sin(float64(angle2)))
				vector.StrokeLine(screen, x1, y1, x2, y2, 2, color.White, true)
			}
			
			// Draw black pentagon pattern (simplified)
			vector.StrokeLine(screen,
				ballCenterX-ballRadius/2, ballCenterY-ballRadius/2,
				ballCenterX+ballRadius/2, ballCenterY-ballRadius/2,
				1, color.Black, true)
			vector.StrokeLine(screen,
				ballCenterX+ballRadius/2, ballCenterY-ballRadius/2,
				ballCenterX+ballRadius/2, ballCenterY+ballRadius/2,
				1, color.Black, true)
		} else {
			screen.DrawImage(iconImage, itemOpts)
		}
	}
	
	// Draw boundaries (optional - makes them visible for debugging)
	/*
	// Top and bottom walls
	vector.DrawFilledRect(screen, float32(-offsetX), float32(-offsetY-20), float32(fieldWidth), 2, color.RGBA{255, 0, 0, 255}, true)
	vector.DrawFilledRect(screen, float32(-offsetX), float32(-offsetY+fieldHeight), float32(fieldWidth), 2, color.RGBA{255, 0, 0, 255}, true)
	
	// Left and right walls (except goal areas)
	vector.DrawFilledRect(screen, float32(-offsetX-20), float32(-offsetY), 2, float32(fieldHeight/2-75), color.RGBA{255, 0, 0, 255}, true)
	vector.DrawFilledRect(screen, float32(-offsetX-20), float32(-offsetY+fieldHeight/2+75), 2, float32(fieldHeight/2-75), color.RGBA{255, 0, 0, 255}, true)
	vector.DrawFilledRect(screen, float32(-offsetX+fieldWidth), float32(-offsetY), 2, float32(fieldHeight/2-75), color.RGBA{255, 0, 0, 255}, true)
	vector.DrawFilledRect(screen, float32(-offsetX+fieldWidth), float32(-offsetY+fieldHeight/2+75), 2, float32(fieldHeight/2-75), color.RGBA{255, 0, 0, 255}, true)
	*/
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
func (m *Map) GetItems() []interfaces.ItemOnMap {
	result := make([]interfaces.ItemOnMap, len(m.Items))
	for i, item := range m.Items {
		result[i] = interfaces.ItemOnMap{
			Name:      item.Name,
			RigidBody: item.RigidBody,
			Item:      item.Item,
		}
	}
	return result
}