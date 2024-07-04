package interfaces

// Component defines an interface for entities' components.
type Component interface {
	// Update updates the component with the given delta time.
	Update(deltaTime float64)
	// Serialize serializes the component to a map.
	Serialize() map[string]interface{}
	// Deserialize deserializes the component from a map.
	Deserialize(data map[string]interface{})
}
