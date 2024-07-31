package actions

import "github.com/joaorufino/gopher-game/internal/interfaces"

// ActionManager manages all available actions.
type ActionManager struct {
	actions map[string]interfaces.Action
}

// NewActionManager creates a new ActionManager instance.
func NewActionManager() *ActionManager {
	return &ActionManager{
		actions: make(map[string]interfaces.Action),
	}
}

// RegisterAction registers a new action.
func (am *ActionManager) RegisterAction(name string, action interfaces.Action) {
	am.actions[name] = action
}

// GetAction retrieves an action by name.
func (am *ActionManager) GetAction(name string) (interfaces.Action, bool) {
	action, exists := am.actions[name]
	return action, exists
}

// GetActions retrieves all actions
func (am *ActionManager) GetActions() map[string]interfaces.Action {
	return am.actions
}

// RegisterBasicActions registers the basic actions with the ActionManager
func RegisterBasicActions(am *ActionManager) {
	am.RegisterAction("spawnContainer", SpawnContainer)
	am.RegisterAction("moveContainer", MoveContainer)
	am.RegisterAction("fly", Fly)
}
