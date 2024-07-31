package particle

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joaorufino/gopher-game/internal/interfaces"
)

// UpdatePosition updates the position based on velocity and delta time.
// position: The current position.
// velocity: The velocity vector.
// deltaTime: The time elapsed since the last update in seconds.
// Returns the updated position.
func UpdatePosition(position, velocity interfaces.Vector2D, deltaTime float64) interfaces.Vector2D {
	return interfaces.Vector2D{
		X: position.X + velocity.X*deltaTime,
		Y: position.Y + velocity.Y*deltaTime,
	}
}

// Particle represents a single particle in a particle system.
type Particle struct {
	Position interfaces.Vector2D
	Velocity interfaces.Vector2D
	LifeTime float64
	Image    *ebiten.Image
}

// NewParticle creates a new particle with the specified properties.
func NewParticle(position, velocity interfaces.Vector2D, lifeTime float64, size int, col color.Color) *Particle {
	img := ebiten.NewImage(size, size)
	img.Fill(col)
	return &Particle{
		Position: position,
		Velocity: velocity,
		LifeTime: lifeTime,
		Image:    img,
	}
}

// ParticleSystem manages a collection of particles.
type ParticleSystem struct {
	particles    []*Particle
	maxParticles int
}

// NewParticleSystem creates a new ParticleSystem with a maximum number of particles.
// maxParticles: The maximum number of particles.
// Returns a pointer to the created ParticleSystem.
func NewParticleSystem(maxParticles int) *ParticleSystem {
	return &ParticleSystem{
		particles:    make([]*Particle, 0, maxParticles),
		maxParticles: maxParticles,
	}
}

// AddParticle adds a new particle to the system.
// position: The initial position of the particle.
// velocity: The initial velocity of the particle.
// lifeTime: The lifetime of the particle in seconds.
// size: The size of the particle.
// col: The color of the particle.
func (ps *ParticleSystem) AddParticle(position, velocity interfaces.Vector2D, lifeTime float64, size int, col color.Color) {
	if len(ps.particles) < ps.maxParticles {
		ps.particles = append(ps.particles, NewParticle(position, velocity, lifeTime, size, col))
	}
}

// Update updates all particles in the system based on the elapsed time.
// deltaTime: The time elapsed since the last update in seconds.
func (ps *ParticleSystem) Update(deltaTime float64) {
	for i := len(ps.particles) - 1; i >= 0; i-- {
		p := ps.particles[i]
		p.Position = UpdatePosition(p.Position, p.Velocity, deltaTime)
		p.LifeTime -= deltaTime
		if p.LifeTime <= 0 {
			ps.particles = append(ps.particles[:i], ps.particles[i+1:]...)
		}
	}
}

// Draw renders all particles in the system on the screen.
// screen: The screen to draw the particles on.
func (ps *ParticleSystem) Draw(screen *ebiten.Image, camera interfaces.Camera) {
	offsetX, offsetY := camera.GetOffset()
	for _, particle := range ps.particles {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(particle.Position.X-offsetX, particle.Position.Y-offsetY)
		screen.DrawImage(particle.Image, opts)
	}
}
