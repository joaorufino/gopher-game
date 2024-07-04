package interfaces

type ActionManager interface {
	GetAction(name string) (Action, bool)
	RegisterAction(name string, action Action)
	GetActions() map[string]Action
}

// Action represents a game action.
type Action func(user Entity)
