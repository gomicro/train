package bar

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrMaxCurrentReached = errors.New("current value is greater total value")
)

// Bar represents a progress bar
type Bar struct {
	countLock    *sync.RWMutex
	currentCount int
	totalCount   int
	started      time.Time
	elapsed      time.Duration
	theme        *BarTheme
}

// New returns a new progress bar
func New(theme *BarTheme, total int) *Bar {
	return &Bar{
		theme:      theme,
		totalCount: total,
		countLock:  &sync.RWMutex{},
	}
}

func (b *Bar) Set(n int) error {
	b.countLock.Lock()
	defer b.countLock.Unlock()

	if n > b.totalCount {
		return fmt.Errorf("bar: set: %w", ErrMaxCurrentReached)
	}

	b.currentCount = n
	return nil
}

func (b *Bar) Incr() bool {
	b.countLock.Lock()
	defer b.countLock.Unlock()

	n := b.currentCount + 1
	if n > b.totalCount {
		return false
	}

	var t time.Time
	if b.started == t {
		b.started = time.Now()
	}

	b.elapsed = time.Since(b.started)
	b.currentCount = n

	return true
}

// Bytes returns the byte presentation of the progress bar
func (b *Bar) bytes() []byte {
	completedWidth := int(float64(b.theme.Width) * (b.CompletedPercent() / 100.00))

	var buf bytes.Buffer

	for _, prepFunc := range b.theme.prependFuncs {
		buf.WriteString(prepFunc(b))
	}

	if b.theme.LeftEnd != 0 {
		buf.WriteByte(b.theme.LeftEnd)
	}

	for i := 0; i < completedWidth-1; i++ {
		buf.WriteByte(b.theme.Fill)
	}

	if completedWidth > 0 {
		if completedWidth < b.theme.Width {
			buf.WriteByte(b.theme.Head)
		} else {
			buf.WriteByte(b.theme.Fill)
		}
	}

	for i := 0; i < b.theme.Width-completedWidth; i++ {
		buf.WriteByte(b.theme.Empty)
	}

	if b.theme.RightEnd != 0 {
		buf.WriteByte(b.theme.RightEnd)
	}

	for _, appFunc := range b.theme.appendFuncs {
		buf.WriteString(appFunc(b))
	}

	return buf.Bytes()
}

// String returns the string representation of the bar
func (b *Bar) String() string {
	return string(b.bytes())
}

func (b *Bar) CompletedPercent() float64 {
	b.countLock.Lock()
	defer b.countLock.Unlock()

	return (float64(b.currentCount) / float64(b.totalCount)) * 100.00
}

func (b *Bar) Current() int {
	b.countLock.Lock()
	defer b.countLock.Unlock()

	return b.currentCount
}

func (b *Bar) Total() int {
	return b.totalCount
}

func (b *Bar) Elapsed() time.Duration {
	return b.elapsed
}
