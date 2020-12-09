package snowflake

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Generator holds info needed for generating custom snowflakes
type Generator struct {
	Epoch     time.Time
	increment uint64
	rand      *rand.Rand
	mu        sync.Mutex
}

// NewGen returns a new generator
func NewGen(epoch time.Time) *Generator {
	return &Generator{Epoch: epoch, rand: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

// Get gets a new snowflake
func (g *Generator) Get() Snowflake {
	var id uint64

	t := time.Now().Sub(g.Epoch).Milliseconds()
	id |= (uint64(t) << 22)
	id |= (g.getRandom() << 16)
	id |= g.getIncrement()

	return Snowflake(fmt.Sprintf("%v", id))
}

func (g *Generator) getIncrement() uint64 {
	g.mu.Lock()
	inc := g.increment
	g.increment++
	g.mu.Unlock()
	return inc
}

func (g *Generator) getRandom() uint64 {
	return uint64(g.rand.Uint32() >> 16)
}
