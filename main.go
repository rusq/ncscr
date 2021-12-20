// Command ncscr - Starry Night screensaver from Norton Commander from 1980s.
// Serves no real purpose other than nostalgic feelings
package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
)

const (
	starsPercent = 0.018                  // % of screen taken by the stars
	explodeSpeed = 150 * time.Millisecond // explosions speed

	// max number of simultaneous explosions
	explosions = 2
)

type skies struct {
	stars []Star

	width, height int
}

// maxStars returns the maximum number of stars for the given screen size
func maxStars(width, height int) int {
	return int(float64(width*height) * float64(starsPercent))
}

func NewSkies(width, height, numStars int) *skies {
	if numStars == 0 {
		numStars = maxStars(width, height)
	}

	stars := make([]Star, numStars)

	return &skies{
		stars:  stars,
		width:  width,
		height: height,
	}
}

func (sky *skies) Play(quit <-chan struct{}, s tcell.Screen) {

	targets := make(chan int, explosions)
	defer close(targets)

	// demolition goroutines
	for i := 0; i <= explosions; i++ {
		go func() {
			for {
				select {
				case <-quit:
					return
				case target := <-targets:
					newX, newY := sky.getRandXY()

					sky.stars[target].Explode(s)
					sky.stars[target].ResetAt(s, newX, newY)
				}
			}
		}()
	}

	timer := time.NewTicker(explodeSpeed)
	defer timer.Stop()

	for i := 0; i < len(sky.stars); i++ {
		select {
		case <-quit:
			return
		case <-timer.C:
			x, y := sky.getRandXY()

			sky.stars[i] = *NewStar(x, y)
			sky.stars[i].Shine(s)
		}
	}

	for {
		select {
		case <-quit:
			return
		case <-timer.C:
			targets <- rand.Int() % len(sky.stars)
		}
	}
}

func (sky *skies) getRandXY() (int, int) {
	return rand.Int() % sky.width, rand.Int() % sky.height
}

func main() {

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Init(); err != nil {
		log.Fatal(err)
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))
	s.Clear()

	quit := make(chan struct{})
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter, tcell.KeyCtrlC:
					close(quit)
					return
				case tcell.KeyCtrlL, tcell.KeyF10:
					s.Sync()
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	x, y := s.Size()
	skies := NewSkies(x, y, maxStars(x, y))
	go skies.Play(quit, s)

	<-quit
	// loop:
	// 	for {
	// 		select {
	// 		case <-quit:
	// 			break loop
	// 		}
	// 	}

	s.Fini()
}

// func makebox(s tcell.Screen) {
// 	w, h := s.Size()

// 	if w == 0 || h == 0 {
// 		return
// 	}

// 	glyphs := []rune{'@', '#', '&', '*', '=', '%', 'Z', 'A'}

// 	lx := rand.Int() % w
// 	ly := rand.Int() % h
// 	lw := rand.Int() % (w - lx)
// 	lh := rand.Int() % (h - ly)
// 	st := tcell.StyleDefault
// 	gl := ' '
// 	if s.Colors() > 256 {
// 		rgb := tcell.NewHexColor(int32(rand.Int() & 0xffffff))
// 		st = st.Background(rgb)
// 	} else if s.Colors() > 1 {
// 		st = st.Background(tcell.Color(rand.Int() % s.Colors()))
// 	} else {
// 		st = st.Reverse(rand.Int()%2 == 0)
// 		gl = glyphs[rand.Int()%len(glyphs)]
// 	}

// 	for row := 0; row < lh; row++ {
// 		for col := 0; col < lw; col++ {
// 			s.SetCell(lx+col, ly+row, st, gl)
// 		}
// 	}
// 	s.Show()
// }
