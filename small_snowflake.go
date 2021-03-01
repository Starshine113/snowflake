package snowflake

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// SmallSnowflake is a uint32-sized ID
type SmallSnowflake uint32

// SmallID is an alias for SmallSnowflake
type SmallID = SmallSnowflake

func (s SmallSnowflake) String() string {
	return strconv.FormatUint(uint64(s), 16)
}

// ParseSmallSnowflake parses a string into a small snowflake, if possible
func ParseSmallSnowflake(s string) (SmallSnowflake, error) {
	i, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return 0, err
	}
	return SmallSnowflake(i), nil
}

// DefaultSmallGenerator is a SmallGenerator with the epoch set to January 1, 2021 at 00:00 UTC
var DefaultSmallGenerator = NewSmallGen(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))

// SmallGenerator holds info needed for generating custom snowflakes
type SmallGenerator struct {
	Epoch     time.Time
	increment uint8
	mu        sync.Mutex
}

// NewSmallGen returns a new small generator
func NewSmallGen(epoch time.Time) *SmallGenerator {
	return &SmallGenerator{
		Epoch:     epoch,
		increment: uint8(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100)),
	}
}

// Get gets a new snowflake
func (g *SmallGenerator) Get() SmallSnowflake {
	var id uint32

	t := time.Now().Sub(g.Epoch).Minutes()
	id |= (uint32(t) << 7)
	id |= uint32(g.getIncrement())

	return SmallSnowflake(id)
}

// GetFromTime gets a snowflake, with only the time field set, that parses as the given time
func (g *SmallGenerator) GetFromTime(t time.Time) SmallSnowflake {
	if g.Epoch.After(t) {
		return 0
	}

	d := t.Sub(g.Epoch).Minutes()
	return SmallSnowflake(uint32(d) << 7)
}

// Parse parses the given snowflake into a timestamp
func (g *SmallGenerator) Parse(s SmallSnowflake) (t time.Time, err error) {
	i := uint32(s)

	timestamp := (i >> 7) + uint32(g.Epoch.UnixNano()/int64(time.Minute))
	t = time.Unix(0, int64(time.Duration(timestamp)*time.Minute)).UTC()
	return
}

func (g *SmallGenerator) getIncrement() uint8 {
	g.mu.Lock()
	g.mu.Unlock()
	inc := g.increment
	g.increment++
	if g.increment >= 127 {
		g.increment = 0
	}
	return inc
}
