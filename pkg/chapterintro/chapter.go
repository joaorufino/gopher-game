// chapterintro/chapterintro.go
package chapterintro

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/joaorufino/cv-game/internal/interfaces"
	"github.com/joaorufino/cv-game/pkg/physics"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

type Letter struct {
	Char      rune
	RigidBody *physics.RigidBody
}

type ChapterIntro struct {
	letters       []Letter
	text          string
	position      interfaces.Point
	font          font.Face
	physicsEngine interfaces.PhysicsEngine
	startTime     time.Time
}

func NewChapterIntro(text string, startPosition interfaces.Point, physicsEngine interfaces.PhysicsEngine) *ChapterIntro {
	ci := &ChapterIntro{
		text:          text,
		position:      startPosition,
		font:          basicfont.Face7x13,
		physicsEngine: physicsEngine,
		startTime:     time.Now(),
	}

	ci.createLetters()
	return ci
}

func (ci *ChapterIntro) createLetters() {
	x, y := ci.position.X, ci.position.Y
	for i, char := range ci.text {
		rb := physics.NewRigidBody(interfaces.Point{X: x, Y: y}, interfaces.Point{X: 10, Y: 10}, 1.0, false, fmt.Sprintf("char%d", i))
		letter := Letter{Char: char, RigidBody: rb}
		ci.letters = append(ci.letters, letter)
		ci.physicsEngine.AddRigidBody(rb)
		x += 10
		if char == '\n' {
			x = ci.position.X
			y += 20
		}
	}
}

func (ci *ChapterIntro) Update(deltaTime float64) {
	for _, letter := range ci.letters {
		letter.RigidBody.Update(deltaTime)
	}
	ci.physicsEngine.Update(deltaTime)
}

func (ci *ChapterIntro) Draw(screen *ebiten.Image, camera interfaces.Camera) {

	// Get the offset from the camera
	offsetX, offsetY := camera.GetOffset()
	for _, letter := range ci.letters {
		text.Draw(screen, string(letter.Char), ci.font, int(letter.RigidBody.Position.X-offsetX), int(letter.RigidBody.Position.Y-offsetY), color.White)
	}
}
