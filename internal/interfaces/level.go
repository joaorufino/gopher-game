package interfaces

// LevelGenerator defines methods for generating game levels.
type LevelGenerator interface {
	GenerateNextLevel() error
}

// Level represents a game level.
type Level interface {
	GetStartVector2D() Vector2D
	GetEndVector2D() Vector2D
	GetEnemies() []AIAgent
}
