package interfaces

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Transition defines the methods for handling state transitions.
type Transition interface {
	Initialize() error
	Execute(fromState, toState State) error
	Finalize() error
	UpdateProgress()
	IsComplete() bool
	HandleStateChange(stateMachine StateMachine, currentScene UIElement)
	Draw(screen *ebiten.Image)
}
