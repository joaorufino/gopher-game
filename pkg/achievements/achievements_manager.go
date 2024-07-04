package achievements

import (
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/joaorufino/cv-game/internal/interfaces"
)

type AchievementManager struct {
	achievements map[string]Achievement
	displayQueue []achievementDisplay
	mu           sync.Mutex
	config       Config
	eventManager interfaces.EventManager
}

type Config struct {
	ScreenWidth            int
	ScreenHeight           int
	DisplayDuration        time.Duration
	MaxAchievementsDisplay int
	PaddingX               float64
	PaddingY               float64
	AchievementOffsetY     float64
	TextOffsetX            float64
	TextOffsetY            float64
}

func NewAchievementManager(config Config, eventManager interfaces.EventManager) *AchievementManager {
	am := &AchievementManager{
		achievements: make(map[string]Achievement),
		displayQueue: []achievementDisplay{},
		config:       config,
		eventManager: eventManager,
	}

	am.registerEventHandlers()
	return am
}

func (am *AchievementManager) registerEventHandlers() {
	am.eventManager.RegisterHandler(interfaces.EventTypeAchievementUnlocked, am.handleAchievementUnlocked)
}

func (am *AchievementManager) handleAchievementUnlocked(event interfaces.Event) {
	achievementName, ok := event.Payload.(string)
	if !ok {
		return
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	if achievement, exists := am.achievements[achievementName]; exists {
		am.displayQueue = append(am.displayQueue, achievementDisplay{
			achievement: achievement,
			startTime:   time.Now(),
		})
	}
}

func (am *AchievementManager) AddAchievement(name, description string, icon *ebiten.Image) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.achievements[name] = Achievement{
		Name:        name,
		Description: description,
		Icon:        icon,
	}
}

func (am *AchievementManager) Update() {
	am.mu.Lock()
	defer am.mu.Unlock()

	now := time.Now()
	filteredQueue := []achievementDisplay{}

	for _, display := range am.displayQueue {
		if now.Sub(display.startTime) < am.config.DisplayDuration {
			filteredQueue = append(filteredQueue, display)
		}
	}

	// Keep only the most recent achievements within the limit
	if len(filteredQueue) > am.config.MaxAchievementsDisplay {
		filteredQueue = filteredQueue[len(filteredQueue)-am.config.MaxAchievementsDisplay:]
	}

	am.displayQueue = filteredQueue
}

func (am *AchievementManager) Draw(screen *ebiten.Image) {
	am.mu.Lock()
	defer am.mu.Unlock()

	for i, display := range am.displayQueue {
		x := float64(am.config.ScreenWidth) - am.config.PaddingX
		y := float64(am.config.ScreenHeight) - am.config.PaddingY - float64(i)*am.config.AchievementOffsetY
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)
		screen.DrawImage(display.achievement.Icon, op)

		textX := x + am.config.TextOffsetX
		textY := y + am.config.TextOffsetY
		ebitenutil.DebugPrintAt(screen, display.achievement.Name, int(textX), int(textY))
		ebitenutil.DebugPrintAt(screen, display.achievement.Description, int(textX), int(textY+am.config.TextOffsetY))
	}
}
