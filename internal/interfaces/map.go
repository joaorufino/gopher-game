package interfaces

import "github.com/hajimehoshi/ebiten/v2"

// Map defines the methods for managing the game map.
type Map interface {
	Update(deltaTime float64)
	Draw(screen *ebiten.Image, camera Camera)
	LoadMap(filepath string, physicsEngine PhysicsEngine) error
	SetObstacles(obstacles []interface{})
	SetPlatforms(platforms []interface{})

	GetObstacles() []interface{}
	GetPlatforms() []interface{}
	SetBackground(imagePath string)
}

// Tile represents a tile on the game map.
type Tile interface {
	GetPosition() Point
	SetPosition(position Point)
	Draw(screen *ebiten.Image, camera Camera)
}
