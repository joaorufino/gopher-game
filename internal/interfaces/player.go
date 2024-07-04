package interfaces

import "github.com/hajimehoshi/ebiten/v2"

// Player defines the methods for the player character.
type Player interface {
	Update(deltaTime float64) error
	Draw(screen *ebiten.Image, camera Camera) error
	GetPosition() Point
	SetPosition(position Point)
	EquipItem(item Item)
}
