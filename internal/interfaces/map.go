package interfaces

import "github.com/hajimehoshi/ebiten/v2"

// Map defines the methods for managing the game map.
type Map interface {
	Update(deltaTime float64)
	Draw(screen *ebiten.Image, camera Camera)
	SetObstacles(obstacles []interface{})
	SetPlatforms(platforms []interface{})

	GetObstacles() []interface{}
	GetPlatforms() []interface{}
}

// Tile represents a tile on the game map.
type Tile interface {
	GetPosition() Vector2D
	SetPosition(position Vector2D)
	Draw(screen *ebiten.Image, camera Camera)
}
