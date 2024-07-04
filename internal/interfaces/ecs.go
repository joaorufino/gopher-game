package interfaces

import (
	"encoding/json"

	"github.com/hajimehoshi/ebiten/v2"
)

type Entity interface {
	AddComponent(name string, component Component)
	GetComponent(name string) Component
	Update(deltaTime float64)
	Serialize() (string, error)
	Deserialize(jsonStr json.RawMessage) error
	Render(screen *ebiten.Image)
}

type World interface {
	AddEntity(entity Entity)
	RemoveEntity(entity Entity)
	Render(screen *ebiten.Image)
	AddSystem(system System)
}

// System represents a system in the ECS.
type System interface {
	Update(deltaTime float64)
}
