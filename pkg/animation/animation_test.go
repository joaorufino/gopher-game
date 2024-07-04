package animation

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Helper function to create a slice of ebiten.Images
func createFrames(count int) []*ebiten.Image {
	frames := make([]*ebiten.Image, count)
	for i := 0; i < count; i++ {
		frames[i] = ebiten.NewImage(1, 1)
	}
	return frames
}

func TestAnimation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Animation Suite")
}

var _ = Describe("Animation", func() {
	var (
		frames []*ebiten.Image
		anim   *Animation
	)

	BeforeEach(func() {
		frames = createFrames(2)
	})

	Describe("Creating a new animation", func() {
		BeforeEach(func() {
			anim = NewAnimation(AnimationConfig{
				Frames:        frames,
				FrameDuration: 0.2,
				Loop:          true,
				MaxLoops:      3,
			})
		})

		It("should have a frame duration of 0.2", func() {
			Expect(anim.frameDuration).To(Equal(0.2))
		})

		It("should loop", func() {
			Expect(anim.loop).To(BeTrue())
		})

		It("should have a maximum of 3 loops", func() {
			Expect(anim.maxLoops).To(Equal(3))
		})
	})

	Describe("Updating the animation", func() {
		BeforeEach(func() {
			anim = NewAnimation(AnimationConfig{
				Frames:        frames,
				FrameDuration: 0.1,
				Loop:          false,
				MaxLoops:      0,
			})
		})

		Context("when updated with 0.15 seconds", func() {
			BeforeEach(func() {
				if err := anim.Update(0.15); err != nil {

				}
			})

			It("should advance to the next frame", func() {
				Expect(anim.currentFrame).To(Equal(1))
			})
		})
	})

	Describe("Pausing and resuming the animation", func() {
		BeforeEach(func() {
			anim = NewAnimation(AnimationConfig{
				Frames:        frames,
				FrameDuration: 0.1,
				Loop:          false,
				MaxLoops:      0,
			})
			anim.Pause()
		})

		Context("when paused and updated with 0.15 seconds", func() {
			BeforeEach(func() {
				if err := anim.Update(0.15); err != nil {

				}
			})

			It("should not advance the frame", func() {
				Expect(anim.currentFrame).To(Equal(0))
			})
		})

		Context("when resumed and updated with 0.15 seconds", func() {
			BeforeEach(func() {
				anim.Resume()
				if err := anim.Update(0.15); err != nil {

				}
			})

			It("should advance to the next frame", func() {
				Expect(anim.currentFrame).To(Equal(1))
			})
		})
	})

	Describe("Resetting the animation", func() {
		BeforeEach(func() {
			anim = NewAnimation(AnimationConfig{
				Frames:        frames,
				FrameDuration: 0.1,
				Loop:          false,
				MaxLoops:      0,
			})
			if err := anim.Update(0.15); err != nil {

			}
			anim.Reset()
		})

		It("should reset to the first frame", func() {
			Expect(anim.currentFrame).To(Equal(0))
		})

		It("should reset the elapsed time to zero", func() {
			Expect(anim.elapsedTime).To(Equal(0.0))
		})
	})
})
