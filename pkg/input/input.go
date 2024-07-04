package input

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/joaorufino/cv-game/internal/interfaces"
)

// https://github.com/hajimehoshi/ebiten/blob/main/examples/touch/main.go
const TOUCH_THRESHOLD = 30

// InputHandler handles input for the game.
type InputHandler struct {
	keyJump       ebiten.Key
	keyUp         ebiten.Key
	keyDown       ebiten.Key
	keyLeft       ebiten.Key
	keyRight      ebiten.Key
	mouseJump     ebiten.MouseButton
	eventManager  interfaces.EventManager
	mouseStartX   int
	mouseStartY   int
	isDragging    bool
	touchStartX   int
	touchStartY   int
	isTouchActive bool
}

// NewInputHandler creates a new InputHandler with default key bindings and an event manager.
func NewInputHandler(eventManager interfaces.EventManager) *InputHandler {
	return &InputHandler{
		keyJump:      ebiten.KeySpace,
		keyUp:        ebiten.KeyW,
		keyDown:      ebiten.KeyS,
		keyLeft:      ebiten.KeyA,
		keyRight:     ebiten.KeyD,
		mouseJump:    ebiten.MouseButtonLeft,
		eventManager: eventManager,
	}
}

// IsJumpPressed checks if the jump key, mouse button, or touch input is pressed.
func (ih *InputHandler) IsJumpPressed() bool {
	return ih.isKeyPressed(ih.keyJump) || ih.isMouseButtonPressed(ih.mouseJump) || ih.isTouchJustPressed()
}

// IsUpPressed checks if the up key or mouse/touch drag up is pressed.
func (ih *InputHandler) IsUpPressed() bool {
	return ih.isKeyPressed(ih.keyUp) || ih.isKeyPressed(ebiten.KeyUp) || ih.isMouseDragUp() || ih.isTouchDragUp()
}

// IsDownPressed checks if the down key or mouse/touch drag down is pressed.
func (ih *InputHandler) IsDownPressed() bool {
	return ih.isKeyPressed(ih.keyDown) || ih.isKeyPressed(ebiten.KeyDown) || ih.isMouseDragDown() || ih.isTouchDragDown()
}

// IsLeftPressed checks if the left movement key or mouse/touch drag left is pressed.
func (ih *InputHandler) IsLeftPressed() bool {
	return ih.isKeyPressed(ih.keyLeft) || ih.isKeyPressed(ebiten.KeyLeft) || ih.isMouseDragLeft() || ih.isTouchDragLeft()
}

// IsRightPressed checks if the right movement key or mouse/touch drag right is pressed.
func (ih *InputHandler) IsRightPressed() bool {
	return ih.isKeyPressed(ih.keyRight) || ih.isKeyPressed(ebiten.KeyRight) || ih.isMouseDragRight() || ih.isTouchDragRight()
}

// Private helper method to check if a key is pressed and dispatch the event.
func (ih *InputHandler) isKeyPressed(key ebiten.Key) bool {
	if ebiten.IsKeyPressed(key) {
		ih.eventManager.Dispatch(interfaces.Event{Type: interfaces.EventType(fmt.Sprintf("KeyPressed_%d", key)), Priority: 1})
		return true
	}
	return false
}

// Private helper method to check if a mouse button is pressed and dispatch the event.
func (ih *InputHandler) isMouseButtonPressed(button ebiten.MouseButton) bool {
	if ebiten.IsMouseButtonPressed(button) {
		ih.eventManager.Dispatch(interfaces.Event{Type: interfaces.EventType(fmt.Sprintf("MouseButtonPressed_%d", button)), Priority: 1})
		return true
	}
	return false
}

// Private helper method to check if a key was just pressed.
func (ih *InputHandler) isKeyJustPressed(key ebiten.Key) bool {
	if inpututil.IsKeyJustPressed(key) {
		ih.eventManager.Dispatch(interfaces.Event{Type: interfaces.EventType(fmt.Sprintf("KeyJustPressed_%d", key)), Priority: 1})
		return true
	}
	return false
}

// Private helper method to check if a mouse button was just pressed.
func (ih *InputHandler) isMouseButtonJustPressed(button ebiten.MouseButton) bool {
	if inpututil.IsMouseButtonJustPressed(button) {
		ih.eventManager.Dispatch(interfaces.Event{Type: interfaces.EventType(fmt.Sprintf("MouseButtonJustPressed_%d", button)), Priority: 1})
		return true
	}
	return false
}

// Private helper method to check if a mouse drag left is detected.
func (ih *InputHandler) isMouseDragLeft() bool {
	if ih.isDragging {
		x, _ := ih.GetMousePosition()
		return x < ih.mouseStartX-30 // Drag left threshold
	}
	return false
}

// Private helper method to check if a mouse drag right is detected.
func (ih *InputHandler) isMouseDragRight() bool {
	if ih.isDragging {
		x, _ := ih.GetMousePosition()
		return x > ih.mouseStartX+30 // Drag right threshold
	}
	return false
}

// Private helper method to check if a mouse drag up is detected.
func (ih *InputHandler) isMouseDragUp() bool {
	if ih.isDragging {
		_, y := ih.GetMousePosition()
		return y < ih.mouseStartY-30 // Drag up threshold
	}
	return false
}

// Private helper method to check if a mouse drag down is detected.
func (ih *InputHandler) isMouseDragDown() bool {
	if ih.isDragging {
		_, y := ih.GetMousePosition()
		return y > ih.mouseStartY+30 // Drag down threshold
	}
	return false
}

// Private helper method to check if a touch drag left is detected.
func (ih *InputHandler) isTouchDragLeft() bool {
	if ih.isTouchActive {
		touches := ebiten.TouchIDs()
		if len(touches) > 0 {
			x, _ := ebiten.TouchPosition(touches[0])
			return x < ih.touchStartX-30 // Drag left threshold
		}
	}
	return false
}

// Private helper method to check if a touch drag right is detected.
func (ih *InputHandler) isTouchDragRight() bool {
	if ih.isTouchActive {
		touches := ebiten.TouchIDs()
		if len(touches) > 0 {
			x, _ := ebiten.TouchPosition(touches[0])
			return x > ih.touchStartX+30 // Drag right threshold
		}
	}
	return false
}

// Private helper method to check if a touch drag up is detected.
func (ih *InputHandler) isTouchDragUp() bool {
	if ih.isTouchActive {
		touches := ebiten.TouchIDs()
		if len(touches) > 0 {
			_, y := ebiten.TouchPosition(touches[0])
			return y < ih.touchStartY-30 // Drag up threshold
		}
	}
	return false
}

// Private helper method to check if a touch drag down is detected.
func (ih *InputHandler) isTouchDragDown() bool {
	if ih.isTouchActive {
		touches := ebiten.TouchIDs()
		if len(touches) > 0 {
			_, y := ebiten.TouchPosition(touches[0])
			return y > ih.touchStartY+30 // Drag down threshold
		}
	}
	return false
}

// Private helper method to check if a touch just pressed is detected.
func (ih *InputHandler) isTouchJustPressed() bool {
	touches := inpututil.AppendJustPressedTouchIDs(make([]ebiten.TouchID, 0))
	if len(touches) > 0 {
		ih.touchStartX, ih.touchStartY = ebiten.TouchPosition(touches[0])
		ih.isTouchActive = true
		return true
	}
	return false
}

// GetMousePosition returns the current mouse position.
func (ih *InputHandler) GetMousePosition() (int, int) {
	return ebiten.CursorPosition()
}

// Update updates the input handler and dispatches events.
func (ih *InputHandler) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		ih.mouseStartX, ih.mouseStartY = ebiten.CursorPosition()
		ih.isDragging = true
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		ih.isDragging = false
	}

	touches := ebiten.TouchIDs()
	if len(touches) > 0 && inpututil.IsTouchJustReleased(touches[0]) {
		ih.isTouchActive = false
	}

	if ih.IsJumpPressed() {
		ih.eventManager.Dispatch(interfaces.Event{Type: "KeyPressed_32", Priority: 1})
	}
	if ih.IsUpPressed() {
		ih.eventManager.Dispatch(interfaces.Event{Type: "KeyPressed_87", Priority: 1})
	}
	if ih.IsDownPressed() {
		ih.eventManager.Dispatch(interfaces.Event{Type: "KeyPressed_83", Priority: 1})
	}
	if ih.IsLeftPressed() {
		ih.eventManager.Dispatch(interfaces.Event{Type: "KeyPressed_65", Priority: 1})
	}
	if ih.IsRightPressed() {
		ih.eventManager.Dispatch(interfaces.Event{Type: "KeyPressed_68", Priority: 1})
	}
	// Check if no keys are pressed and dispatch the "NoKeyPressed" event
	if !ih.IsJumpPressed() && !ih.IsUpPressed() && !ih.IsDownPressed() && !ih.IsLeftPressed() && !ih.IsRightPressed() {
		ih.eventManager.Dispatch(interfaces.Event{Type: "NoKeyPressed", Priority: 1})
	}
	return nil
}
