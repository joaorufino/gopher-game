// items/manager.go
package items

import (
	"encoding/json"
	"fmt"

	"github.com/joaorufino/cv-game/internal/interfaces"
	"github.com/joaorufino/cv-game/internal/utils"
)

// ItemManagerImpl is a concrete implementation of the ItemManager interface.
type ItemManagerImpl struct {
	items        map[string]interfaces.Item
	abilitiesMgr interfaces.AbilitiesManager
}

// NewItemManager creates a new ItemManagerImpl.
func NewItemManager(abilitiesMgr interfaces.AbilitiesManager) *ItemManagerImpl {
	return &ItemManagerImpl{
		items:        make(map[string]interfaces.Item),
		abilitiesMgr: abilitiesMgr,
	}
}

// LoadItems loads items from a JSON file.
func (im *ItemManagerImpl) LoadItems(path string) error {
	return utils.LoadData(path, func(data []byte) error {
		var itemsData []Item
		if err := json.Unmarshal(data, &itemsData); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
		for _, itemData := range itemsData {
			item := &Item{
				Name:        itemData.Name,
				Image:       itemData.Image,
				Icon:        itemData.Icon,
				Description: itemData.Description,
				Appearance:  itemData.Appearance,
				Abilities:   itemData.Abilities,
				Version:     itemData.Version,
			}
			im.items[item.Name] = item
		}
		return nil
	})
}

// GetItem returns an item by its name.
func (im *ItemManagerImpl) GetItem(name string) (interfaces.Item, error) {
	item, exists := im.items[name]
	if !exists {
		return nil, fmt.Errorf("item not found: %s", name)
	}
	return item, nil
}

// AddItem adds a new item to the manager.
func (im *ItemManagerImpl) AddItem(item interfaces.Item) {
	im.items[item.GetName()] = item
}

// RemoveItem removes an item by its name.
func (im *ItemManagerImpl) RemoveItem(name string) {
	delete(im.items, name)
}

// GetAllItems returns all items managed by the ItemManager.
func (im *ItemManagerImpl) GetAllItems() []interfaces.Item {
	allItems := make([]interfaces.Item, 0, len(im.items))
	for _, item := range im.items {
		allItems = append(allItems, item)
	}
	return allItems
}
