package snowflake

import (
	"math/rand"
	"sync"
	"time"
)

// Generator holds info needed for generating custom snowflakes
type Generator struct {
	Epoch     time.Time
	increment uint16
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

	return Snowflake(id)
}

// GetFromTime gets a snowflake, with only the time field set, that parses as the given time
func (g *Generator) GetFromTime(t time.Time) Snowflake {
	if g.Epoch.After(t) {
		return 0
	}

	d := t.Sub(g.Epoch).Milliseconds()
	return Snowflake(uint64(d) << 22)
}

// Parse parses the given snowflake into a timestamp
func (g *Generator) Parse(s Snowflake) (t time.Time, err error) {
	i := int64(s)

	timestamp := (i >> 22) + int64(g.Epoch.UnixNano()/1000000)
	t = time.Unix(0, timestamp*1000000).UTC()
	return
}

func (g *Generator) getIncrement() uint64 {
	g.mu.Lock()
	inc := uint64(g.increment)
	g.increment++
	g.mu.Unlock()
	return inc
}

func (g *Generator) getRandom() uint64 {
	return uint64(g.rand.Uint32() >> 16)
}
