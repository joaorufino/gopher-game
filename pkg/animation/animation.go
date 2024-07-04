package animation

import (
	"errors"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// Animation represents a sequence of images to create an animation.
type Animation struct {
	mu            sync.Mutex
	frames        []*ebiten.Image
	currentFrame  int
	frameDuration float64
	elapsedTime   float64
	paused        bool
	loop          bool
	loopCount     int
	maxLoops      int
	onStart       func()
	onEnd         func()
	onFrameChange func(int)
}

// AnimationConfig holds configuration options for creating an Animation.
type AnimationConfig struct {
	Frames        []*ebiten.Image
	FrameDuration float64
	Loop          bool
	MaxLoops      int
	OnStart       func()
	OnEnd         func()
	OnFrameChange func(int)
}

// NewAnimation creates a new Animation with the given configuration.
func NewAnimation(config AnimationConfig) *Animation {
	if config.FrameDuration <= 0 {
		config.FrameDuration = 0.1 // default frame duration
	}
	return &Animation{
		frames:        config.Frames,
		frameDuration: config.FrameDuration,
		loop:          config.Loop,
		maxLoops:      config.MaxLoops,
		onStart:       config.OnStart,
		onEnd:         config.OnEnd,
		onFrameChange: config.OnFrameChange,
	}
}

// Update progresses the animation based on the elapsed time.
// deltaTime: The time elapsed since the last update in seconds.
func (a *Animation) Update(deltaTime float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.paused {
		return nil
	}
	a.elapsedTime += deltaTime
	if a.elapsedTime >= a.frameDuration {
		a.elapsedTime = 0
		a.currentFrame++
		if a.onFrameChange != nil {
			a.onFrameChange(a.currentFrame)
		}
		if a.currentFrame >= len(a.frames) {
			if a.loop && (a.maxLoops == 0 || a.loopCount < a.maxLoops) {
				a.currentFrame = 0
				a.loopCount++
			} else {
				a.currentFrame = len(a.frames) - 1
				if a.onEnd != nil {
					a.onEnd()
				}
				return errors.New("animation completed")
			}
		}
	}
	return nil
}

// Draw renders the current frame of the animation on the screen.
// screen: The screen to draw the animation on.
// opts: Options for drawing the image.
func (a *Animation) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.frames) == 0 {
		return // No frames to draw
	}
	screen.DrawImage(a.frames[a.currentFrame], opts)
}

// Pause pauses the animation.
func (a *Animation) Pause() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.paused = true
}

// Resume resumes the animation if it was paused.
func (a *Animation) Resume() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.paused = false
	if a.onStart != nil {
		a.onStart()
	}
}

// Reset resets the animation to the first frame and elapsed time to zero.
func (a *Animation) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.currentFrame = 0
	a.elapsedTime = 0
	a.loopCount = 0
	a.paused = false
}

// SetFrameDuration dynamically adjusts the frame duration.
// frameDuration: The new duration for each frame in seconds.
func (a *Animation) SetFrameDuration(frameDuration float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if frameDuration > 0 {
		a.frameDuration = frameDuration
	}
}

// AnimationBuilder helps in constructing an Animation using a fluent interface.
type AnimationBuilder struct {
	config AnimationConfig
}

// NewAnimationBuilder initializes a new AnimationBuilder.
func NewAnimationBuilder() *AnimationBuilder {
	return &AnimationBuilder{
		config: AnimationConfig{},
	}
}

// Frames sets the frames for the animation.
func (b *AnimationBuilder) Frames(frames []*ebiten.Image) *AnimationBuilder {
	b.config.Frames = frames
	return b
}

// FrameDuration sets the duration for each frame.
func (b *AnimationBuilder) FrameDuration(duration float64) *AnimationBuilder {
	b.config.FrameDuration = duration
	return b
}

// Loop sets the loop behavior for the animation.
func (b *AnimationBuilder) Loop(loop bool) *AnimationBuilder {
	b.config.Loop = loop
	return b
}

// MaxLoops sets the maximum number of loops for the animation.
func (b *AnimationBuilder) MaxLoops(maxLoops int) *AnimationBuilder {
	b.config.MaxLoops = maxLoops
	return b
}
