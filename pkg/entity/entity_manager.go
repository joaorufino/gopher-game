package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joaorufino/cv-game/internal/interfaces"
	"github.com/joaorufino/cv-game/pkg/entities"
)

type EntityManager struct {
	entities map[string]*entities.Entity
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		entities: make(map[string]*entities.Entity),
	}
}

func (em *EntityManager) AddEntity(entity *entities.Entity) {
	em.entities[entity.ID] = entity
}

func (em *EntityManager) RemoveEntity(id string) {
	delete(em.entities, id)
}

func (em *EntityManager) Update(deltaTime float64) {
	for _, entity := range em.entities {
		entity.Update(deltaTime)
	}
}

func (em *EntityManager) Draw(screen *ebiten.Image, cam interfaces.Camera) {
	for _, entity := range em.entities {
		entity.Draw(screen, cam)
	}
}
