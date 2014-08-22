package peco

import (
	"testing"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type dummyScreen struct {
	*interceptor
	width  int
	height int
}

func (d dummyScreen) SetCell(x, y int, ch rune, fg, bg termbox.Attribute) {
	d.record("SetCell", interceptorArgs{x, y, ch, fg, bg})
}
func (d dummyScreen) Clear(fg, bg termbox.Attribute) error { return nil }
func (d dummyScreen) Flush() error                         { return nil }
func (d dummyScreen) PollEvent() chan termbox.Event        { return nil }
func (d dummyScreen) Size() (int, int) {
	return d.width, d.height
}

func TestLayoutType(t *testing.T) {
	layouts := []struct {
		value    LayoutType
		expectOK bool
	}{
		{LayoutTypeTopDown, true},
		{LayoutTypeBottomUp, true},
		{"foobar", false},
	}
	for _, l := range layouts {
		valid := IsValidLayoutType(l.value)
		if valid != l.expectOK {
			t.Errorf("LayoutType %s, expected IsValidLayoutType to return %s, but got %s",
				l.value,
				l.expectOK,
				valid,
			)
		}
	}
}

func TestPrintScreen(t *testing.T) {
	i := newInterceptor()
	screen = dummyScreen{
		i,
		100,
		100,
	}

	makeVerifier := func(initX, initY int, fill bool) func(string) {
		return func(msg string) {
			i.reset()
			t.Logf("Checking printScreen(%d, %d, %s, %s)", initX, initY, msg, fill)
			width := utf8.RuneCountInString(msg)
			printScreen(initX, initY, termbox.ColorDefault, termbox.ColorDefault, msg, fill)
			events := i.events["SetCell"]
			if !fill {
				if len(events) != width {
					t.Errorf("Expected %d SetCell events, got %d",
						width,
						len(events),
					)
				}
				return
			}

			// fill == true
			w, _ := screen.Size()
			if rw := runewidth.StringWidth(msg); rw != width {
				w -= rw - width
			}
			if len(events) != w {
				t.Errorf("Expected %d SetCell events, got %d",
					w,
					len(events),
				)
				return
			}
		}
	}

	verify := makeVerifier(0, 0, false)
	verify("Hello, World!")
	verify("日本語")

	verify = makeVerifier(0, 0, true)
	verify("Hello, World!")
	verify("日本語")
}