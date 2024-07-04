// items/item.go
package items

import (
	"github.com/joaorufino/cv-game/internal/interfaces"
)

// Item represents a game item.
type Item struct {
	Name        string                `json:"name"`
	Image       string                `json:"image"`
	Icon        string                `json:"icon"`
	Description string                `json:"description"`
	Appearance  interfaces.Appearance `json:"appearance"`
	Abilities   []string              `json:"abilities"`
	Version     int                   `json:"version"`
}

// GetName returns the name of the item.
func (i *Item) GetName() string {
	return i.Name
}

// GetImagePath returns the image path of the item.
func (i *Item) GetImagePath() string {
	return i.Image
}

// GetIconPath returns the icon path of the item.
func (i *Item) GetIconPath() string {
	return i.Icon
}

// GetDescription returns the description of the item.
func (i *Item) GetDescription() string {
	return i.Description
}

// GetAppearance returns the appearance of the item.
func (i *Item) GetAppearance() interfaces.Appearance {
	return i.Appearance
}

// GetAbilities returns the abilities of the item.
func (i *Item) GetAbilities() []string {
	return i.Abilities
}

// GetVersion returns the version of the item.
func (i *Item) GetVersion() int {
	return i.Version
}
