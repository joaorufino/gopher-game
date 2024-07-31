package pet

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joaorufino/gopher-game/internal/interfaces"
	"github.com/joaorufino/gopher-game/pkg/animation"
	"github.com/joaorufino/gopher-game/pkg/particle"
	"github.com/joaorufino/gopher-game/pkg/physics"
)

// Pet represents the pet character that follows the player.
type Pet struct {
	Position            interfaces.Vector2D
	velocity            interfaces.Vector2D
	animations          map[string]*animation.Animation
	currentAnimation    string
	lastAnimationUpdate time.Time
	onGround            bool
	particleSystem      *particle.ParticleSystem
	resourceManager     interfaces.ResourceManager
	config              *Configuration
	player              interfaces.Player
	RigidBody           *physics.RigidBody
}

// Configuration holds the configurable settings for the Pet.
type Configuration struct {
	ImageScale   float64
	RunVelocity  float64
	JumpVelocity float64
}

// NewPet initializes a new pet instance.
func NewPet(startX, startY float64, resourceManager interfaces.ResourceManager, config *Configuration, physicsEngine interfaces.PhysicsEngine, player interfaces.Player) *Pet {
	frameCounts := map[string]int{
		"idle": 1,
		"run":  1,
		"jump": 1,
	}
	animations, size := animation.LoadAnimations(resourceManager, "/images/player/normal", frameCounts)
	size.X = size.X * config.ImageScale
	size.Y = size.Y * config.ImageScale

	pet := &Pet{
		Position:            interfaces.Vector2D{X: startX, Y: startY},
		velocity:            interfaces.Vector2D{X: 0, Y: 0},
		animations:          animations,
		currentAnimation:    "idle",
		lastAnimationUpdate: time.Now(),
		onGround:            true,
		particleSystem:      particle.NewParticleSystem(100),
		resourceManager:     resourceManager,
		config:              config,
		player:              player,
		RigidBody:           physics.NewRigidBody(interfaces.Vector2D{X: startX, Y: startY}, size, 500, false, "pet"),
	}
	physicsEngine.AddRigidBody(pet.RigidBody)
	return pet
}

func (p *Pet) Update(deltaTime float64) error {
	p.followPlayer()
	if err := p.animations[p.currentAnimation].Update(deltaTime); err != nil {
		log.Printf("animation update error: %v", err)
	}

	p.updatePosition(deltaTime)
	p.Position = p.RigidBody.GetPosition()
	p.updateAnimationState()
	p.particleSystem.Update(deltaTime)
	return nil
}

func (p *Pet) followPlayer() {
	playerPos := p.player.GetPosition()
	pPos := p.Position

	if playerPos.X > pPos.X+30 {
		p.RigidBody.Velocity.X = p.config.RunVelocity
	} else if playerPos.X+30 < pPos.X {
		p.RigidBody.Velocity.X = -p.config.RunVelocity
	} else {
		p.RigidBody.Velocity.X = 0
	}

	if playerPos.Y+30 > pPos.Y {
		p.RigidBody.Velocity.Y = p.config.JumpVelocity
	} else if playerPos.Y+30 < pPos.Y {
		p.RigidBody.Velocity.Y = -p.config.JumpVelocity
	} else {
		p.RigidBody.Velocity.Y = 0
	}
}

func (p *Pet) updatePosition(deltaTime float64) {
	p.RigidBody.Update(deltaTime)
}

func (p *Pet) updateAnimationState() {
	if p.RigidBody.Velocity.X != 0 {
		p.currentAnimation = "run"
	} else if !p.RigidBody.OnGround {
		p.currentAnimation = "jump"
	} else {
		p.currentAnimation = "idle"
	}
}

func (p *Pet) Draw(screen *ebiten.Image, cam interfaces.Camera) error {
	offsetX, offsetY := cam.GetOffset()

	petOpts := &ebiten.DrawImageOptions{}
	petOpts.GeoM.Scale(p.config.ImageScale, p.config.ImageScale)
	petOpts.GeoM.Translate(p.Position.X-offsetX, p.Position.Y-offsetY)
	p.animations[p.currentAnimation].Draw(screen, petOpts)

	p.particleSystem.Draw(screen, cam)
	return nil
}

func (p *Pet) GetPosition() interfaces.Vector2D {
	return p.RigidBody.GetPosition()
}

func (p *Pet) SetPosition(po interfaces.Vector2D) {
	p.Position = po
}

func (p *Pet) GetCurrentAnimation() string {
	return p.currentAnimation
}

func (p *Pet) SetCurrentAnimation(animation string) {
	p.currentAnimation = animation
}
