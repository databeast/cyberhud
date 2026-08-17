package svgrender

import "time"

// Animator is a time-driven frame sequencer that advances through a sequence
// of SVG frames based on elapsed time. It implements the widgets.Animated
// interface via its Tick method.
type Animator struct {
	frames         []Frame
	loop           bool
	currentIndex   int
	elapsedInFrame time.Duration
	done           bool
}

// NewAnimator creates a new Animator with the given frame sequence and loop flag.
// Returns nil if frames is empty or nil.
func NewAnimator(frames []Frame, loop bool) *Animator {
	if len(frames) == 0 {
		return nil
	}
	return &Animator{
		frames: frames,
		loop:   loop,
	}
}

// Tick advances the animator's internal clock by elapsed duration,
// potentially advancing the frame index when the current frame's duration
// is exceeded. This method satisfies the widgets.Animated interface.
func (a *Animator) Tick(elapsed time.Duration) {
	if a.done {
		return
	}

	a.elapsedInFrame += elapsed

	for a.elapsedInFrame >= a.frames[a.currentIndex].Duration {
		a.elapsedInFrame -= a.frames[a.currentIndex].Duration

		if a.loop {
			a.currentIndex = (a.currentIndex + 1) % len(a.frames)
		} else {
			if a.currentIndex < len(a.frames)-1 {
				a.currentIndex++
			} else {
				// Past the final frame duration for non-looping: hold on last frame.
				a.done = true
				return
			}
		}
	}
}

// CurrentFrame returns the SVG content and index of the current frame.
func (a *Animator) CurrentFrame() (svg string, index int) {
	return a.frames[a.currentIndex].SVG, a.currentIndex
}

// Done returns true if the animator is non-looping and has completed
// its full frame sequence (elapsed time exceeded the final frame's duration).
func (a *Animator) Done() bool {
	return a.done
}
