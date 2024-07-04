package interfaces

// Settings defines the methods for managing game settings.
type Settings interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{})
	Save(path string) error
	Load(path string) error
	GetScreenWidth() int
	GetScreenHeight() int
	IsFullscreen() bool
}
