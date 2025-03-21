package hud

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/joaorufino/gopher-game/pkg/score"
	"golang.org/x/image/font"
)

type HUD struct {
	ScoreManager *score.ScoreManager
	Font         font.Face
	ScreenWidth  int
	ScreenHeight int
}

func NewHUD(scoreManager *score.ScoreManager, font font.Face, screenWidth, screenHeight int) *HUD {
	return &HUD{
		ScoreManager: scoreManager,
		Font:         font,
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
	}
}

func (h *HUD) Draw(screen *ebiten.Image) {
	// Draw team scores at the top center
	homeTeam := h.ScoreManager.GetTeamName(0)
	awayTeam := h.ScoreManager.GetTeamName(1)
	homeScore := h.ScoreManager.GetScore(0)
	awayScore := h.ScoreManager.GetScore(1)

	scoreText := fmt.Sprintf("%s %d - %d %s", homeTeam, homeScore, awayScore, awayTeam)

	// If we have a font, use it for nicer rendering
	if h.Font != nil {
		bounds := text.BoundString(h.Font, scoreText)
		x := (h.ScreenWidth - bounds.Dx()) / 2
		text.Draw(screen, scoreText, h.Font, x, 30, color.White)
	} else {
		// Fallback to debug print
		ebitenutil.DebugPrintAt(screen, scoreText, (h.ScreenWidth-len(scoreText)*6)/2, 20)
	}

	// Draw match time if match is active
	if h.ScoreManager.IsMatchActive() {
		timeInSeconds := h.ScoreManager.GetMatchTime()
		minutes := timeInSeconds / 60
		seconds := timeInSeconds % 60
		timeText := fmt.Sprintf("%02d:%02d", minutes, seconds)

		if h.Font != nil {
			bounds := text.BoundString(h.Font, timeText)
			x := (h.ScreenWidth - bounds.Dx()) / 2
			text.Draw(screen, timeText, h.Font, x, 60, color.White)
		} else {
			ebitenutil.DebugPrintAt(screen, timeText, (h.ScreenWidth-len(timeText)*6)/2, 40)
		}
	} else {
		// Show "Match Ended" or "Press Space to Start" if no match is active
		statusText := "Press Space to Start Match"
		if h.Font != nil {
			bounds := text.BoundString(h.Font, statusText)
			x := (h.ScreenWidth - bounds.Dx()) / 2
			text.Draw(screen, statusText, h.Font, x, 60, color.RGBA{255, 255, 0, 255})
		} else {
			ebitenutil.DebugPrintAt(screen, statusText, (h.ScreenWidth-len(statusText)*6)/2, 40)
		}
	}
}
