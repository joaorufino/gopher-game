// abilities/manager.go
package abilities

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/joaorufino/gopher-game/internal/interfaces"
	"github.com/joaorufino/gopher-game/internal/utils"
)

// AbilitiesManagerImpl implements the AbilitiesManager interface.
type AbilitiesManagerImpl struct {
	abilities    map[string]*Ability
	actionMap    map[string]interfaces.Action
	eventManager interfaces.EventManager
	mu           sync.RWMutex
}

// NewAbilitiesManager creates a new AbilitiesManagerImpl.
func NewAbilitiesManager(actionMap map[string]interfaces.Action, eventManager interfaces.EventManager) interfaces.AbilitiesManager {
	am := &AbilitiesManagerImpl{
		abilities:    make(map[string]*Ability),
		actionMap:    actionMap,
		eventManager: eventManager,
	}

	am.registerEventHandlers()
	return am
}

// LoadAbilities loads abilities from a JSON file.
func (am *AbilitiesManagerImpl) LoadAbilities(path string) error {
	return utils.LoadData(path, func(data []byte) error {
		var abilitiesData []Ability
		if err := json.Unmarshal(data, &abilitiesData); err != nil {
			return fmt.Errorf("failed to unmarshal abilities JSON: %w", err)
		}
		for _, abilityData := range abilitiesData {
			am.abilities[abilityData.Name] = NewAbility(
				abilityData.Name,
				abilityData.Image,
				abilityData.Icon,
				abilityData.Description,
				abilityData.ActionName,
				abilityData.Cooldown,
				nil,
				am.eventManager,
			)
		}
		return nil
	})
}

// GetAbility retrieves an ability by name.
func (am *AbilitiesManagerImpl) GetAbility(name string) (interfaces.Ability, bool) {
	am.mu.RLock()
	defer am.mu.RUnlock()
	ability, exists := am.abilities[name]
	return ability, exists
}

func (am *AbilitiesManagerImpl) registerEventHandlers() {
	am.eventManager.RegisterHandler(interfaces.EventTypeAbilityUsed, am.handleAbilityUsed)
}

func (am *AbilitiesManagerImpl) handleAbilityUsed(event interfaces.Event) {
	payload, ok := event.Payload.(map[string]interface{})
	if !ok {
		return
	}

	abilityName, ok := payload["abilityName"].(string)
	if !ok {
		return
	}

	user, ok := payload["user"].(interfaces.Entity)
	if !ok {
		return
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	if ability, exists := am.abilities[abilityName]; exists {
		ability.Activate(user, am.actionMap)
	}
}

// Update maintains the state of all abilities over time.
func (am *AbilitiesManagerImpl) Update(deltaTime float64) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	for _, ability := range am.abilities {
		ability.Update(deltaTime)
	}
}
