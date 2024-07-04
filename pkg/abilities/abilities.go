package abilities

import (
	"fmt"
	"time"

	"github.com/joaorufino/cv-game/internal/interfaces"
)

// Ability represents a special power or skill that can be used in the game.
type Ability struct {
	Name         string        `json:"name"`                // Unique identifier
	Image        string        `json:"image"`               // Path to image asset
	Icon         string        `json:"icon"`                // Path to icon asset
	Description  string        `json:"description"`         // Text description
	Cooldown     time.Duration `json:"cooldown"`            // Time between uses
	ActionName   string        `json:"actionName"`          // Name of the action to execute
	animationStr *string       `json:"animation,omitempty"` // Animation associated with the ability (if any)
	Animation    interfaces.Animation
	lastUsed     time.Time
	eventManager interfaces.EventManager
}

// NewAbility creates a new Ability instance.
func NewAbility(name, image, icon, description, actionName string, cooldown time.Duration, animation interfaces.Animation, eventManager interfaces.EventManager) *Ability {
	return &Ability{
		Name:         name,
		Image:        image,
		Icon:         icon,
		Description:  description,
		Cooldown:     cooldown,
		ActionName:   actionName,
		Animation:    animation,
		eventManager: eventManager,
	}
}

// CanActivate checks if the ability is off cooldown and can be used.
func (a *Ability) CanActivate() bool {
	return time.Since(a.lastUsed) >= a.Cooldown
}

// Activate executes the ability's action and updates the last used time.
func (a *Ability) Activate(user interfaces.Entity, actionMap map[string]interfaces.Action) error {
	if a.CanActivate() {
		if action, exists := actionMap[a.ActionName]; exists {
			action(user)
			a.lastUsed = time.Now()
			a.triggerAbilityUsedEvent(user)
			return nil
		}
	}
	return fmt.Errorf("ability is on cooldown or action does not exist")
}

// Deactivate ends the ability's effects if applicable.
func (a *Ability) Deactivate() error {
	// Implement the logic to deactivate the ability, if necessary
	return nil
}

// Update maintains the ability's state over time.
func (a *Ability) Update(deltaTime float64) error {
	// Implement any ongoing effects or state maintenance for the ability
	return nil
}

// GetName returns the name of the ability.
func (a *Ability) GetName() string {
	return a.Name
}

// GetCooldown returns the current cooldown duration.
func (a *Ability) GetCooldown() float64 {
	return a.Cooldown.Seconds()
}

// SetCooldown sets the cooldown duration for the ability.
func (a *Ability) SetCooldown(duration float64) {
	a.Cooldown = time.Duration(duration) * time.Second
}

// IsOnCooldown checks if the ability is currently on cooldown.
func (a *Ability) IsOnCooldown() bool {
	return !a.CanActivate()
}

// triggerAbilityUsedEvent triggers an event indicating that the ability was used.
func (a *Ability) triggerAbilityUsedEvent(user interfaces.Entity) {
	a.eventManager.Dispatch(interfaces.Event{
		Type:     interfaces.EventTypeAbilityUsed,
		Priority: 1,
		Payload: map[string]interface{}{
			"abilityName": a.Name,
			"user":        user,
		},
	})
}
