package source

import (
	"math"
	"math/rand"
)

// PseudocodePool contains the 48 pseudocode text snippets used for fragments.
var PseudocodePool = []string{
	"if (signal.ready())",
	"await decrypt(key)",
	"for i in range(n)",
	"buffer.flush()",
	"return hash(data)",
	"while !complete",
	"socket.connect()",
	"yield next()",
	"mutex.lock()",
	"chan <- payload",
	"defer close(conn)",
	"select { case <-",
	"go func() {",
	"async fn process()",
	"match result {",
	"loop { recv() }",
	"try! decode(stream)",
	"let cipher = new()",
	"pub fn verify(sig)",
	"impl Trait for Node",
	"spawn(|| task())",
	"pipe(filter, map)",
	"reduce(acc, val)",
	"emit(event, data)",
	"subscribe(topic)",
	"publish(channel)",
	"query.execute()",
	"index.rebuild()",
	"cache.invalidate()",
	"merkle.verify(root)",
	"consensus.propose()",
	"validate(block)",
	"broadcast(peers)",
	"replicate(shard)",
	"compress(frame)",
	"encrypt(payload)",
	"authenticate(token)",
	"authorize(scope)",
	"serialize(struct)",
	"deserialize(bytes)",
	"allocate(pages)",
	"gc.collect()",
	"interrupt(handler)",
	"syscall(trap)",
	"mmap(addr, len)",
	"futex.wait()",
	"epoll.add(fd)",
	"io_uring.submit()",
}

// ComputeFragmentOpacity returns the current opacity for a fragment at the given time.
// Returns -1 when the fragment has expired (elapsed >= total duration).
func ComputeFragmentOpacity(f ActiveFragment, time float64) float64 {
	elapsed := time - f.StartTime
	if elapsed < 0 {
		return 0
	}

	totalDuration := f.FadeInDuration + f.HoldDuration + f.FadeOutDuration
	if elapsed >= totalDuration {
		return -1
	}

	// Fade-in phase.
	if elapsed < f.FadeInDuration {
		return f.PeakOpacity * (elapsed / f.FadeInDuration)
	}

	// Hold phase.
	if elapsed < f.FadeInDuration+f.HoldDuration {
		return f.PeakOpacity
	}

	// Fade-out phase.
	fadeOutElapsed := elapsed - f.FadeInDuration - f.HoldDuration
	return f.PeakOpacity * (1 - fadeOutElapsed/f.FadeOutDuration)
}

// SpawnFragment attempts to spawn a new pseudocode text fragment near an active cluster.
// Returns true if a fragment was successfully spawned.
func SpawnFragment(state *FragmentState, time float64, activeClusters []ClusterConfig, panelW, panelH int, rng *rand.Rand, scaleFactor float64) bool {
	// Constraint (a): max active fragments.
	if len(state.ActiveFragments) >= MaxActiveFragments {
		return false
	}

	// Constraint (b): minimum spawn interval.
	if time-state.LastSpawnTime < MinSpawnInterval {
		return false
	}

	// Constraint (c): at least one active cluster.
	if len(activeClusters) == 0 {
		return false
	}

	// Select text from pool excluding lastSpawnedText.
	text := selectFragmentText(rng, state.LastSpawnedText)

	// Select target cluster uniformly at random.
	cluster := activeClusters[rng.Intn(len(activeClusters))]

	// Position: random angle + distance from cluster center.
	angle := rng.Float64() * 2 * math.Pi
	radius := rng.Float64() * scaledProximityRadius(scaleFactor)
	cx := (cluster.CenterXPct / 100.0) * float64(panelW)
	cy := (cluster.CenterYPct / 100.0) * float64(panelH)
	x := cx + radius*math.Cos(angle)
	y := cy + radius*math.Sin(angle)

	// Clamp to panel bounds.
	x = clamp(x, 0, float64(panelW))
	y = clamp(y, 0, float64(panelH))

	// Generate parameters.
	fadeIn := randRange(rng, 1, 2)
	hold := randRange(rng, 2, 5)
	fadeOut := randRange(rng, 1, 2)
	fontSize := randRange(rng, scaledFontSizeMin(scaleFactor), scaledFontSizeMax(scaleFactor))
	peakOpacity := randRange(rng, 0.3, 0.8)

	// Color: 50% green, 50% cyan.
	var color HSLColor
	if rng.Float64() < 0.5 {
		color = HSLColor{
			H: randRange(rng, 100, 140),
			S: randRange(rng, 70, 100),
			L: randRange(rng, 40, 60),
		}
	} else {
		color = HSLColor{
			H: randRange(rng, 170, 190),
			S: randRange(rng, 80, 100),
			L: randRange(rng, 45, 65),
		}
	}

	fragment := ActiveFragment{
		Text:            text,
		X:               x,
		Y:               y,
		StartTime:       time,
		FadeInDuration:  fadeIn,
		HoldDuration:    hold,
		FadeOutDuration: fadeOut,
		FontSize:        fontSize,
		Color:           color,
		PeakOpacity:     peakOpacity,
	}

	state.ActiveFragments = append(state.ActiveFragments, fragment)
	state.LastSpawnTime = time
	state.LastSpawnedText = text
	return true
}

// selectFragmentText selects a random entry from PseudocodePool that differs from lastText.
func selectFragmentText(rng *rand.Rand, lastText string) string {
	for {
		text := PseudocodePool[rng.Intn(len(PseudocodePool))]
		if text != lastText {
			return text
		}
	}
}
