package interfaces

// AIManager defines the methods for managing AI behaviors.
type AIManager interface {
	// Initialize initializes the AI manager with necessary parameters.
	Initialize(world World)

	// PlaceEnemies places enemies in the specified rectangles.
	PlaceEnemies(rects []Rect)

	// Update updates the AI manager and all managed AI entities.
	Update(deltaTime float64)

	// AddEnemy adds an enemy to the manager with a specific behavior.
	AddEnemy(x, y float64, behavior string)
}

// AIAgent defines the interface for an AI agent in the game.
type AIAgent interface {
	// Initialize initializes the AI agent. This is called once when the agent is created.
	Initialize() error
	// Update updates the AI agent with the given delta time.
	// deltaTime: The time elapsed since the last update in seconds.
	Update(deltaTime float64) error
	// SetTarget sets the target for the AI agent.
	// target: The target point the AI agent should move towards or interact with.
	SetTarget(target Vector2D)
	// GetTarget returns the current target of the AI agent.
	GetTarget() Vector2D
	// GetPosition returns the current position of the AI agent.
	GetPosition() Vector2D
	// SetPosition sets the position of the AI agent.
	// position: The new position for the AI agent.
	SetPosition(position Vector2D)
	// OnEnter is called when the AI agent is activated or enters a new state.
	OnEnter() error
	// OnExit is called when the AI agent is deactivated or exits its current state.
	OnExit() error
}
