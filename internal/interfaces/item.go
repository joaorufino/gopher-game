// interfaces/item.go
package interfaces

// Appearance represents the appearance details of an item.
type Appearance struct {
	Type           string `json:"type"`
	Color          string `json:"color"`
	Material       string `json:"material,omitempty"`
	SpecialEffects string `json:"specialEffects,omitempty"`
}

// Item represents an item in the game.
type Item interface {
	GetName() string
	GetImagePath() string
	GetIconPath() string // New method for icon path
	GetDescription() string
	GetAppearance() Appearance
	GetAbilities() []string
	GetVersion() int
}

// ItemManager defines the methods for managing items in the game.
type ItemManager interface {
	LoadItems(path string) error
	GetItem(name string) (Item, error)
	AddItem(item Item)
	RemoveItem(name string)
	GetAllItems() []Item
}
