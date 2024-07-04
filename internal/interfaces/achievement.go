package interfaces

import "github.com/hajimehoshi/ebiten/v2"

type AchievementManager interface {
	AddAchievement(name, description string, icon *ebiten.Image)
	Update()
	Draw(screen *ebiten.Image)
}
