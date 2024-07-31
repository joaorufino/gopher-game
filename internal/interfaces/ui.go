package interfaces

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

// UIElement defines the interface that all UI elements must implement.
type UIElement interface {
	Initialize() error
	// Update updates the state of the UI element.
	// deltaTime: The time elapsed since the last update in seconds.
	Update(deltaTime float64) error
	// Draw renders the UI element on the screen.
	// screen: The screen to draw the UI element on.
	Draw(screen *ebiten.Image) error
	// OnEnter is called when the UI element is entered.
	OnEnter() error
	// OnExit is called when the UI element is exited.
	OnExit() error
}

// UI defines the methods for managing the game's user interface.
// TODO: this interface doesnt make much sense
// Ideally we would move the creation of parts to other packages
type UI interface {
	AddElement(element UIElement)
	RemoveElement(element UIElement)
	Draw(screen *ebiten.Image)
	NewMenu(switchScene SwitchSceneFunc, ui UI, resources ButtonResources) UIElement
	NewSettings(switchScene SwitchSceneFunc, ui UI) UIElement
	NewAchievements(switchScene SwitchSceneFunc, ui UI) UIElement
	NewEducation(switchScene SwitchSceneFunc, ui UI) UIElement
	NewPersonalInfo(switchScene SwitchSceneFunc, ui UI) UIElement
	NewSkills(switchScene SwitchSceneFunc, ui UI) UIElement
	Update(deltaTime float64) error
	HandleGlobalInput(switchSceneFunc SwitchSceneFunc)
}

// SwitchSceneFunc is a function type that handles scene switching.
type SwitchSceneFunc func(scene string)

// StatusBar defines the methods for managing the game's status bar.
type StatusBar interface {
	Initialize() error
	Update(deltaTime float64) error
	Draw(screen *ebiten.Image) error
	SetValue(value int)
	GetValue() int
	GetContainer() *widget.Container
}

// ButtonResources defines the methods for managing button resources.
type ButtonResources interface {
	LoadResources(path string) error
	GetButton(name string) (Button, error)
}

// Button represents a button resource.
type Button interface {
	Draw(screen *ebiten.Image, position Vector2D) error
	GetTextPadding() widget.Insets
	GetImage() *widget.ButtonImage
}
