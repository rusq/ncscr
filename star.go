package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

var chars = []rune{'\u00b7', '\u2219', '\u2022', '\u2666', '\u25A0', '\u263C'}

const (
	boomIdx        = 1
	animationSpeed = 100 * time.Millisecond
)

// Star is one star
type Star struct {
	x, y int
	sync.RWMutex
}

// NewStar creates a new star
func NewStar(x, y int) *Star {
	return &Star{x: x, y: y}
}

// Reset resets the exploded star at new coordinates.  If coordinates are
// not defined, it resets it at old coordinates.
func (s *Star) Reset(scr tcell.Screen) {
	s.Lock()
	defer s.Unlock()
	s.ResetAt(scr, -1, -1)
}

// ResetAt resets the star at the new position
func (s *Star) ResetAt(scr tcell.Screen, x, y int) {
	if !(x == -1 || y == -1) {
		s.Lock()
		s.x, s.y = x, y
		s.Unlock()
	}
	s.Shine(scr)
}

// Shine draws the star
func (s *Star) Shine(scr tcell.Screen) {

	s.draw(scr, tcell.Style.Foreground(tcell.StyleDefault, tcell.ColorTeal), chars[0])
}

// Explode explodes the star
func (s *Star) Explode(scr tcell.Screen) {
	defer s.clear(scr)

	for _, c := range chars[boomIdx : boomIdx+rand.Int()%(len(chars))] {
		s.draw(scr, tcell.Style.Foreground(tcell.StyleDefault, tcell.ColorWhite), c)
		<-time.After(animationSpeed)
	}
}

// clear clears the star
func (s *Star) clear(scr tcell.Screen) {
	s.draw(scr, tcell.Style.Foreground(tcell.StyleDefault, tcell.ColorTeal), ' ')
}

// draw is the actual draw
func (s *Star) draw(scr tcell.Screen, style tcell.Style, r rune) {
	s.RLock()
	defer s.RUnlock()
	scr.SetCell(s.x, s.y, style, r)
	scr.Show()
}
