package interfaces

// Ability represents a special power or skill that can be used in the game.
type Ability interface {
	// Activate triggers the ability's effects.
	Activate(user Entity, actionMap map[string]Action) error
	// Deactivate ends the ability's effects if applicable.
	Deactivate() error
	// Update maintains the ability's state over time.
	Update(deltaTime float64) error
	// GetName returns the name of the ability.
	GetName() string
	// GetCooldown returns the current cooldown duration.
	GetCooldown() float64
	// SetCooldown sets the cooldown duration for the ability.
	SetCooldown(duration float64)
	// IsOnCooldown checks if the ability is currently on cooldown.
	IsOnCooldown() bool
}

type AbilitiesManager interface {
	LoadAbilities(path string) error
	GetAbility(name string) (Ability, bool)
	Update(deltatime float64)
}
