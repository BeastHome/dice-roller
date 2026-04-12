package dice

import (
	"math/rand"
	"time"
)

// EngineOptions defines optional configuration for the Engine.
// All fields are optional and default to zero values.
type EngineOptions struct {
	Seed     int64 // if non-zero, overrides RNG seed
	Verbose  bool  // TODO: future feature - reserved for controlling AttachVerbose output
	MaxDepth int   // TODO: future feature - reserved for limiting explode/compound recursion
}

// NewEngineWithOptions constructs an Engine using optional settings.
func NewEngineWithOptions(opt EngineOptions) *Engine {
	seed := opt.Seed
	if seed == 0 {
		seed = defaultSeed()
	}

	return &Engine{
		rng: newRNG(seed),
	}
}

// defaultSeed returns the default RNG seed.
func defaultSeed() int64 {
	return time.Now().UnixNano()
}

// newRNG constructs a new RNG from a seed.
func newRNG(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}
