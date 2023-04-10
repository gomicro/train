package bar

import (
	"errors"
	"fmt"
)

var (
	DefaultTheme = &BarTheme{
		Empty:    '-',
		Fill:     '=',
		Head:     '>',
		LeftEnd:  '[',
		RightEnd: ']',
		Width:    70,
	}

	ErrEmtpyCharNotSet = errors.New("empty character not set")
	ErrFillCharNotSet  = errors.New("fill character not set")
	ErrHeadCharNotSet  = errors.New("head character not set")
	ErrWidthTooSmall   = errors.New("width set too small")
)

type BarTheme struct {
	Empty    byte
	Fill     byte
	Head     byte
	LeftEnd  byte
	RightEnd byte
	Width    int

	appendFuncs  []DecoratorFunc
	prependFuncs []DecoratorFunc
}

type DecoratorFunc func(b *Bar) string

func NewTheme(empty, fill, head byte, width int) (*BarTheme, error) {
	switch {
	case empty == 0:
		return nil, fmt.Errorf("bar: theme: %w", ErrEmtpyCharNotSet)
	case fill == 0:
		return nil, fmt.Errorf("bar: theme: %w", ErrFillCharNotSet)
	case head == 0:
		return nil, fmt.Errorf("bar: theme: %w", ErrHeadCharNotSet)
	case width == 0:
		return nil, fmt.Errorf("bar: theme: %w", ErrWidthTooSmall)
	}

	return &BarTheme{
		Empty: empty,
		Fill:  fill,
		Head:  head,
		Width: width,
	}, nil
}

func NewThemeFromTheme(t *BarTheme) *BarTheme {
	c := *t
	return &c
}

func (t *BarTheme) Append(f DecoratorFunc) {
	t.appendFuncs = append(t.appendFuncs, f)
}

func (t *BarTheme) Prepend(f DecoratorFunc) {
	t.prependFuncs = append(t.prependFuncs, f)
}
