package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joaorufino/cv-game/internal/interfaces"
	"github.com/joaorufino/cv-game/pkg/components"
)

type Entity struct {
	ID         string
	Components map[string]interface{}
}

func NewEntity(id string) *Entity {
	return &Entity{
		ID:         id,
		Components: make(map[string]interface{}),
	}
}

func (e *Entity) AddComponent(name string, component interface{}) {
	e.Components[name] = component
}

func (e *Entity) GetComponent(name string) interface{} {
	return e.Components[name]
}

func (e *Entity) Update(deltaTime float64) {
	if physicsComp, ok := e.Components["physics"].(*components.PhysicsComponent); ok {
		physicsComp.Update(deltaTime)
	}
}

func (e *Entity) Draw(screen *ebiten.Image, cam interfaces.Camera) {
	if renderComp, ok := e.Components["render"].(*components.RenderComponent); ok {
		if transformComp, ok := e.Components["transform"].(*components.TransformComponent); ok {
			offsetX, offsetY := cam.GetOffset()
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(transformComp.Position.X-offsetX, transformComp.Position.Y-offsetY)
			opts.GeoM.Scale(transformComp.Scale.X, transformComp.Scale.Y)
			opts.GeoM.Rotate(transformComp.Rotation)
			renderComp.Draw(screen, opts)
		}
	}
}
