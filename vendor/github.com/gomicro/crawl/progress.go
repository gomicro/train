package crawl

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/gomicro/crawl/bar"
	"github.com/gosuri/uilive"
)

const defaultRefreshInterval = time.Millisecond * 10

// Progress represents the container that renders progress bars
type Progress struct {
	ctx             context.Context
	out             io.Writer
	refreshInterval time.Duration
	lw              *uilive.Writer
	mtx             *sync.RWMutex

	// Bars is the collection of progress bars
	Bars []*bar.Bar
}

// New returns a new progress bar with defaults
func New(ctx context.Context, out io.Writer) *Progress {
	lw := uilive.New()
	lw.Out = out

	return &Progress{
		ctx:             ctx,
		out:             out,
		refreshInterval: defaultRefreshInterval,
		lw:              lw,
		mtx:             &sync.RWMutex{},

		Bars: make([]*bar.Bar, 0),
	}
}

func (p *Progress) SetOut(o io.Writer) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.out = o
	p.lw.Out = o
}

func (p *Progress) SetRefreshInterval(interval time.Duration) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.refreshInterval = interval
}

func (p *Progress) AddBar(b *bar.Bar) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.Bars = append(p.Bars, b)
}

// Listen listens for updates and renders the progress bars
func (p *Progress) Listen() {
	for {
		p.mtx.Lock()
		interval := p.refreshInterval
		p.mtx.Unlock()

		select {
		case <-time.After(interval):
			p.print()
		case <-p.ctx.Done():
			p.print()
			return
		}
	}
}

func (p *Progress) print() {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	for _, bar := range p.Bars {
		fmt.Fprintln(p.lw, bar.String())
	}
	p.lw.Flush()
}

// Start starts the rendering the progress of progress bars. It listens for
// updates using `bar.Set(n)` and new bars when added using `AddBar`
func (p *Progress) Start() {
	go p.Listen()
}

// Stop stops listening
func (p *Progress) Stop() {
	_, cancel := context.WithCancel(p.ctx)
	cancel()
}

// Bypass returns a writer which allows non-buffered data to be written to the
// underlying output
func (p *Progress) Bypass() io.Writer {
	return p.lw.Bypass()
}
