package interfaces

// LevelGenerator defines methods for generating game levels.
type LevelGenerator interface {
	GenerateNextLevel() error
}

// Level represents a game level.
type Level interface {
	GetStartPoint() Point
	GetEndPoint() Point
	GetEnemies() []AIAgent
}
