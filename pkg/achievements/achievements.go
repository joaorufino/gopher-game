package achievements

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Achievement struct {
	Name        string
	Description string
	Icon        *ebiten.Image
}

type achievementDisplay struct {
	achievement Achievement
	startTime   time.Time
}
