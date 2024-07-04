package interfaces

// EventType represents the type of an event.
type EventType string

// Event represents a generic event with a type, priority, and payload.
type Event struct {
	Type     EventType
	Priority int
	Payload  interface{}
}

// EventHandler is a function that handles an event.
type EventHandler func(Event)

// EventManager defines the methods for managing event registration and dispatching.
type EventManager interface {
	// RegisterHandler registers a handler for the specified event type.
	RegisterHandler(eventType EventType, handler EventHandler)
	// Dispatch dispatches an event to the registered handlers.
	Dispatch(event Event)
	// Wait waits for all events to be processed.
	Wait()
}

// Define your event types as needed.
const (
	EventPlayerJump   EventType = "PlayerJump"
	EventPlayerMove   EventType = "PlayerMove"
	EventItemEquipped EventType = "ItemEquipped"
	EventSceneSwitch  EventType = "SceneSwitch"
	EventTypeInput    EventType = "Input"

	// Input Events
	EventKeyPressed             EventType = "KeyPressed"
	EventKeyJustPressed         EventType = "KeyJustPressed"
	EventMouseButtonPressed     EventType = "MouseButtonPressed"
	EventMouseButtonJustPressed EventType = "MouseButtonJustPressed"

	EventTypeAbilityUsed         EventType = "AbilityUsed"
	EventTypeAchievementUnlocked EventType = "AchievementUnlocked"
)
