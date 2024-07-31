package player

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joaorufino/gopher-game/internal/interfaces"
	"github.com/joaorufino/gopher-game/pkg/animation"
	"github.com/joaorufino/gopher-game/pkg/particle"
	"github.com/joaorufino/gopher-game/pkg/physics"
)

// Player represents the player character.
type Player struct {
	Position            interfaces.Vector2D
	velocity            interfaces.Vector2D
	animations          map[string]*animation.Animation
	currentAnimation    string
	lastAnimationUpdate time.Time
	onGround            bool
	particleSystem      *particle.ParticleSystem
	items               []interfaces.Item
	itemIcons           []*ebiten.Image // Field for item icons
	resourceManager     interfaces.ResourceManager
	config              *Configuration
	canFly              bool
	RigidBody           *physics.RigidBody
	EventManager        interfaces.EventManager
	rotation            float64
	gameWidth           float64 // Add gameWidth to constrain movement
}

// Configuration holds the configurable settings for the Player.
type Configuration struct {
	ScreenWidth  int
	ScreenHeight int
	Gravity      float64
	JumpVelocity float64
	RunVelocity  float64
	ImageScale   float64
}

// NewPlayer initializes a new player instance.
func NewPlayer(startX, startY float64, resourceManager interfaces.ResourceManager, config *Configuration, physicsEngine interfaces.PhysicsEngine, event interfaces.EventManager, gameWidth float64) *Player {
	frameCounts := map[string]int{
		"idle": 1,
		"run":  1,
		"jump": 1,
		"pick": 1,
	}
	pAnimations, size := animation.LoadAnimations(resourceManager, "/images/player/normal", frameCounts)

	size.X = size.X * config.ImageScale
	size.Y = size.Y * config.ImageScale
	player := &Player{
		Position:            interfaces.Vector2D{X: startX, Y: startY},
		velocity:            interfaces.Vector2D{X: 0, Y: 0},
		animations:          pAnimations,
		currentAnimation:    "idle",
		lastAnimationUpdate: time.Now(),
		onGround:            true,
		particleSystem:      particle.NewParticleSystem(100),
		items:               []interfaces.Item{},
		itemIcons:           []*ebiten.Image{}, // Initialize itemIcons
		resourceManager:     resourceManager,
		config:              config,
		canFly:              false, // Initialize without the cloud item
		RigidBody:           physics.NewRigidBody(interfaces.Vector2D{X: startX, Y: startY}, size, 1000, false, "player"),
		EventManager:        event,
		gameWidth:           gameWidth, // Set gameWidth
	}
	player.RigidBody.SetCanPick(true)
	// Add the player's rigid body to the physics engine
	physicsEngine.AddRigidBody(player.RigidBody)

	// Register input handlers for the player
	player.registerInputHandlers()

	return player
}

// Update updates the player's state.
func (p *Player) Update(deltaTime float64) error {
	if err := p.animations[p.currentAnimation].Update(deltaTime); err != nil {
		log.Printf("animation update error: %v", err)
	}

	p.updatePosition(deltaTime)
	p.Position = p.RigidBody.GetPosition() // Sync player position with rigid body position
	p.particleSystem.Update(deltaTime)
	return nil
}

// updatePosition updates the player's position.
func (p *Player) updatePosition(deltaTime float64) {
	// Update the rigid body position
	p.RigidBody.Update(deltaTime)

	// Constrain player within game boundaries
	if p.RigidBody.Position.X < 0 {
		p.RigidBody.Position.X = 0
	} else if p.RigidBody.Position.X > p.gameWidth-p.RigidBody.Size.X {
		p.RigidBody.Position.X = p.gameWidth - p.RigidBody.Size.X
	}
}

// Draw renders the player and particle system on the screen.
func (p *Player) Draw(screen *ebiten.Image, cam interfaces.Camera) error {
	// Get the offset from the camera
	offsetX, offsetY := cam.GetOffset()

	// Calculate the center of the frame
	centerX, centerY := float64(165)/2, float64(205)/2

	// Draw the player at the correct position with the offset
	playerOpts := &ebiten.DrawImageOptions{}
	playerOpts.GeoM.Translate(-centerX, -centerY) // Move to the center of the image
	playerOpts.GeoM.Rotate(p.rotation)            // Apply rotation
	playerOpts.GeoM.Translate(centerX, centerY)   // Move back to the original position

	playerOpts.GeoM.Scale(p.config.ImageScale, p.config.ImageScale)
	playerOpts.GeoM.Translate(p.Position.X-offsetX, p.Position.Y-offsetY)

	p.animations[p.currentAnimation].Draw(screen, playerOpts)

	// Draw the particle system
	p.particleSystem.Draw(screen, cam)

	// Draw item icons
	iconX, iconY := 10.0, 10.0 // Starting position for icons
	iconSpacing := 5.0         // Space between icons
	for _, icon := range p.itemIcons {
		iconOpts := &ebiten.DrawImageOptions{}
		iconOpts.GeoM.Scale(0.5, 0.5) // Adjust scale as necessary
		iconOpts.GeoM.Translate(iconX, iconY)
		screen.DrawImage(icon, iconOpts)
		iconX += float64(icon.Bounds().Dx())*0.5 + iconSpacing
	}

	return nil
}

// EquipItem equips an item to the player and updates the player's animations.
func (p *Player) EquipItem(item interfaces.Item) {
	p.items = append(p.items, item)
	frameCounts := map[string]int{
		"idle": 1,
		"run":  1,
		"jump": 1,
		"pick": 1,
	}
	p.currentAnimation = "pick"

	p.animations, _ = animation.LoadAnimations(p.resourceManager, item.GetImagePath(), frameCounts)

	// Load and add the item icon
	icon, err := p.resourceManager.LoadImage(item.GetIconPath())
	if err != nil {
		log.Printf("failed to load item icon: %v", err)
	} else {
		p.itemIcons = append(p.itemIcons, icon)
	}

	for _, ability := range item.GetAbilities() {
		if ability == "Cloud Ride" {
			p.canFly = true
		}
	}
}

func (p *Player) GetPosition() interfaces.Vector2D {
	return p.RigidBody.GetPosition()
}

func (p *Player) SetPosition(po interfaces.Vector2D) {
	p.Position = po
}

// registerInputHandlers registers input handlers for the player.
func (p *Player) registerInputHandlers() {
	if !p.RigidBody.OnGround && !p.canFly {
		p.currentAnimation = "jump"
	}
	p.EventManager.RegisterHandler("NoKeyPressed", func(event interfaces.Event) {
		p.handleStatic()
	})
	p.EventManager.RegisterHandler("KeyPressed_32", func(event interfaces.Event) {
		p.handleJump()
	})
	p.EventManager.RegisterHandler("KeyPressed_87", func(event interfaces.Event) {
		p.handleMoveUp()
	})
	p.EventManager.RegisterHandler("KeyPressed_65", func(event interfaces.Event) {
		p.handleMoveLeft()
	})
	p.EventManager.RegisterHandler("KeyPressed_83", func(event interfaces.Event) {
		p.handleMoveDown()
	})
	p.EventManager.RegisterHandler("KeyPressed_68", func(event interfaces.Event) {
		p.handleMoveRight()
	})
	p.EventManager.RegisterHandler(interfaces.EventItemEquipped, p.handleEquipItemEvent)
}

func (p *Player) handleEquipItemEvent(event interfaces.Event) {
	payload, ok := event.Payload.(map[string]interface{})
	if !ok {
		return
	}

	itemName, ok := payload["itemName"].(string)
	if !ok {
		return
	}

	item, err := p.resourceManager.GetItem(itemName)
	if err != nil {
		log.Printf("failed to get item: %v", err)
		return
	}
	p.EquipItem(item)
}

func (p *Player) handleJump() {
	if p.RigidBody.OnGround {
		p.currentAnimation = "jump"
		p.RigidBody.Velocity.Y = -p.config.JumpVelocity
		p.RigidBody.OnGround = false
		p.particleSystem.AddParticle(p.Position, interfaces.Vector2D{X: 0, Y: -100}, 1.0, 5, color.RGBA{255, 255, 255, 255})
		p.rotation = 0
	}
}

func (p *Player) handleMoveUp() {
	// Implement behavior for upward movement if necessary
}

func (p *Player) handleMoveLeft() {
	if p.RigidBody.OnGround {
		p.currentAnimation = "run"
		p.rotation = 0
	} else {
		p.rotation -= 0.1
	}
	p.RigidBody.Velocity.X = -p.config.RunVelocity
}

func (p *Player) handleMoveDown() {
	// Implement behavior for downward movement if necessary
}

func (p *Player) handleMoveRight() {
	if p.RigidBody.OnGround {
		p.currentAnimation = "run"
		p.rotation = 0
	} else {
		p.rotation += 0.1
	}
	p.RigidBody.Velocity.X = p.config.RunVelocity
}

func (p *Player) handleStatic() {
	if p.RigidBody.OnGround {
		p.currentAnimation = "idle"
		p.rotation = 0
	}
	p.RigidBody.Velocity.X = 0
}
