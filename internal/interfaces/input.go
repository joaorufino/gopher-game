package interfaces

// InputHandler defines the interface for input handling
type InputHandler interface {
	IsJumpPressed() bool
	IsUpPressed() bool
	IsDownPressed() bool
	IsLeftPressed() bool
	IsRightPressed() bool
	Update() error
}

// Key represents a key on the keyboard.
type Key int

const (
	KeyArrowLeft Key = iota
	KeyArrowRight
	KeyArrowUp
	KeyArrowDown
	KeySpace
)
