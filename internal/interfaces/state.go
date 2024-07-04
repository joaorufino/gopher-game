package interfaces

import "github.com/hajimehoshi/ebiten/v2"

// State defines the interface for a game state.
type State interface {
	// Initialize initializes the state. This is called once when the state is created.
	Initialize() error
	// Update updates the state with the given delta time.
	// deltaTime: The time elapsed since the last update in seconds.
	Update(deltaTime float64) error
	// Draw renders the state on the screen.
	// screen: The screen to draw the state on.
	Draw(screen *ebiten.Image) error
	// OnEnter is called when the state is entered.
	OnEnter() error
	// OnExit is called when the state is exited.
	OnExit() error
}

// StateMachine defines the interface for managing game states.
type StateMachine interface {
	// AddState adds a new state to the manager.
	// name: The name of the state.
	// state: The state to be added.
	AddState(name string, state State)
	// RemoveState removes a state from the manager.
	// name: The name of the state to be removed.
	RemoveState(name string)

	// Update updates the current state.
	Update(deltaTime float64) error
	// Draw renders the current state on the screen.
	Draw(screen *ebiten.Image) error
	// GetCurrentState returns the current state.
	GetCurrentState() State
	// GetState returns the state with the specified name.
	GetState(name string) (State, bool)

	ChangeState(newState State)
}
