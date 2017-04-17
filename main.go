package main

import (
	"strconv"

	tl "github.com/JoelOtter/termloop"
)

// Maze

func generateMaze(w, h int) [][]rune {
	maze := make([][]rune, w)
	return maze
}

// Block (based on the termloop pyramid sample)

type Block struct {
	*tl.Rectangle
	px        int // previous x position
	py        int // previous y position
	move      bool
	g         *tl.Game
	w         int // width of maze
	h         int // height of maze
	score     int
	scoretext *tl.Text
}

func NewBlock(x, y int, color tl.Attr, g *tl.Game, w, h, score int, scoretext *tl.Text) *Block {
	b := &Block{
		g:         g,
		w:         w,
		h:         h,
		score:     score,
		scoretext: scoretext,
	}
	b.Rectangle = tl.NewRectangle(x, y, 1, 1, color)
	return b
}

func (b *Block) Draw(s *tl.Screen) {
	if l, ok := b.g.Screen().Level().(*tl.BaseLevel); ok {
		// Set the level offset so the player is always in the
		// center of the screen. This simulates moving the camera
		sw, sh := s.Size()
		x, y := b.Position()
		l.SetOffset(sw/2-x, sh/2-y)
	}
	b.Rectangle.Draw(s)
}

func (b *Block) Tick(ev tl.Event) {
	if ev.Type == tl.EventKey {
		b.px, b.py = b.Position()
		switch ev.Key {
		case tl.KeyArrowRight:
			b.SetPosition(b.px+1, b.py)
		case tl.KeyArrowLeft:
			b.SetPosition(b.px-1, b.py)
		case tl.KeyArrowUp:
			b.SetPosition(b.px, b.py-1)
		case tl.KeyArrowDown:
			b.SetPosition(b.px, b.py+1)
		}
	}
}

func (b *Block) Collide(c tl.Physical) {
	if r, ok := c.(*tl.Rectangle); ok {
		if r.Color() == tl.ColorWhite {
			// Collision with walls
			b.SetPosition(b.px, b.py)
		} else if r.Color() == tl.ColorBlue {
			b.w += 1
			b.h += 1
			buildLevel(b.g, b.w, b.h, b.score)
		}
	}
}

func buildLevel(g *tl.Game, w, h, score int) {
	maze := generateMaze(w, h)
	l := tl.NewBaseLevel(tl.Cell{})
	g.Screen().SetLevel(l)
	scoretext := tl.NewText(0, 1, "Levels explored: "+strconv.Itoa(score), tl.ColorBlue, tl.ColorBlack)
	g.Screen().AddEntity(tl.NewText(0, 0, "Maze-Box (based on Pyramid)!", tl.ColorBlue, tl.ColorBlack))
	g.Screen().AddEntity(scoretext)
	for i, row := range maze {
		for j, path := range row {
			if path == '*' {
				l.AddEntity(tl.NewRectangle(i, j, 1, 1, tl.ColorWhite))
			} else if path == 'S' {
				col := tl.RgbTo256Color(0xff, 0, 0)
				l.AddEntity(NewBlock(i, j, col, g, w, h, score, scoretext))
			} else if path == 'L' {
				l.AddEntity(tl.NewRectangle(i, j, 1, 1, tl.ColorBlue))
			}
		}
	}
}

func main() {
	game := tl.NewGame()

	screen := game.Screen()

	screen.AddEntity(tl.NewText(0, 0, "Maze-Box!", tl.ColorRed, tl.ColorBlack))

	cell := tl.Cell{
		Bg: tl.ColorBlack,
		Fg: tl.ColorBlack,
		Ch: ' ',
	}

	level := tl.NewBaseLevel(cell)
	screen.SetLevel(level)

	game.Start()
}
