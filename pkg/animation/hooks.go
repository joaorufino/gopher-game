package animation

// OnStart sets the function to be called when the animation starts.
func (b *AnimationBuilder) OnStart(onStart func()) *AnimationBuilder {
	b.config.OnStart = onStart
	return b
}

// OnEnd sets the function to be called when the animation ends.
func (b *AnimationBuilder) OnEnd(onEnd func()) *AnimationBuilder {
	b.config.OnEnd = onEnd
	return b
}

// OnFrameChange sets the function to be called when the frame changes.
func (b *AnimationBuilder) OnFrameChange(onFrameChange func(int)) *AnimationBuilder {
	b.config.OnFrameChange = onFrameChange
	return b
}

// Build creates an Animation with the specified configuration.
func (b *AnimationBuilder) Build() *Animation {
	return NewAnimation(b.config)
}
