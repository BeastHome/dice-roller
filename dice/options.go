package dice

import (
	"math/rand"
	"time"
)

// EngineOptions configures an Engine.
//
// NOTE on Seed: a value of zero is treated as "use the time-based
// default" — to seed deterministically with literal 0, use
// NewEngineWithSeed(0) instead. NewEngineWithSeed is also the
// recommended constructor for any deterministic use (replay,
// embedded consumer test harnesses, the --seed CLI flag).
type EngineOptions struct {
	Seed int64
}

// NewEngineWithOptions constructs an Engine using the given options.
// If Seed is zero, a time-based default seed is used.
func NewEngineWithOptions(opt EngineOptions) *Engine {
	seed := opt.Seed
	if seed == 0 {
		seed = defaultSeed()
	}
	return &Engine{rng: newRNG(seed)}
}

// NewEngineWithSeed constructs an Engine with an explicit deterministic
// RNG seed. Prefer this for reproducible output; call NewEngine() for a
// time-based default.
func NewEngineWithSeed(seed int64) *Engine {
	return &Engine{rng: newRNG(seed)}
}

// defaultSeed returns the default RNG seed.
func defaultSeed() int64 {
	return time.Now().UnixNano()
}

// newRNG constructs a new RNG from a seed.
func newRNG(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}
