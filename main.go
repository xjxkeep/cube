package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"time"

	"github.com/eiannone/keyboard"
)

var (
	cubes = [][3][3]int{
		{
			{1, 1, 1},
			{0, 1, 0},
			{0, 1, 0},
		},
		{
			{1, 1, 1},
			{1, 1, 1},
			{1, 1, 1},
		},
		{
			{1, 1, 0},
			{0, 1, 0},
			{0, 1, 0},
		},
		{
			{0, 1, 0},
			{0, 1, 0},
			{0, 1, 0},
		},
		{
			{0, 1, 0},
			{1, 1, 1},
			{0, 1, 0},
		},
		{
			{0, 1, 0},
			{1, 1, 1},
			{0, 0, 0},
		},
	}
)

type Cube struct {
	Body          [3][3]int
	X, Y          int
	Width, Height int
}

func (c *Cube) Rotate(angle int) {
	temp := [3][3]int{}
	switch angle {
	case 90:
		for i := range c.Body {
			for j := range c.Body[i] {
				temp[j][2-i] = c.Body[i][j]
			}
		}
	case -90:
		for i := range c.Body {
			for j := range c.Body[i] {
				temp[2-j][i] = c.Body[i][j]
			}
		}
	}
	c.Body = temp
}

type Game struct {
	View      [][]int
	Backgroud [][]int
	Cube      *Cube
	Height    int
	Width     int
	Score     int
	NextCube  *Cube
	speed     int
}

func (g *Game) Move(dx, dy int) {
	g.Cube.X += dx
	g.Cube.Y += dy
	g.Clip()
}

func (g *Game) Reset() {
	g.Score = 0
	g.Cube.X = 0
	g.Cube.Y = 0
	g.speed = 500
	for i := range g.Backgroud {
		clear(g.Backgroud[i])
		clear(g.View[i])
	}
	for i := range g.Backgroud[0] {
		g.Backgroud[g.Height-3][i] = 1
		g.Backgroud[g.Height-2][i] = 1
		g.Backgroud[g.Height-1][i] = 1
	}
	for i := range g.Backgroud {
		g.Backgroud[i][0] = 1
		g.Backgroud[i][1] = 1
		g.Backgroud[i][g.Width-1] = 1
		g.Backgroud[i][g.Width-2] = 1
	}
}

func (g *Game) Clip() {
	if g.Cube.X < 0 {
		g.Cube.X = 0
	}
	if g.Cube.Y < 0 {
		g.Cube.Y = 0
	}
	if g.Cube.Y+g.Cube.Height > g.Height {
		g.Cube.Y = g.Height - g.Cube.Height
	}
	if g.Cube.X+g.Cube.Width > g.Width {
		g.Cube.X = g.Width - g.Cube.Width
	}
}
func (g *Game) MoveAbs(x, y int) {
	g.Cube.X = x
	g.Cube.Y = y
	g.Clip()
}

func (g *Game) Next() {
	g.NextCube = &Cube{
		Body:   cubes[rand.Intn(len(cubes))],
		X:      g.Width / 2,
		Y:      0,
		Width:  3,
		Height: 3,
	}
	if v := rand.Intn(3); v == 0 {
		g.NextCube.Rotate(90)
	} else if v == 1 {
		g.NextCube.Rotate(-90)
	}
}

func (g *Game) GenernateCube() {
	g.Cube = g.NextCube
	g.speed = 500
	g.Next()
}

func (g *Game) SettleScores() int {
	row := g.Height - 4
	for {
		if !slices.Contains(g.View[row], 0) {
			for r := row; r > 0; r-- {
				copy(g.View[r], g.View[r-1])
			}
			g.Score++
		} else {
			row--
			if row < 0 {
				return g.Score
			}
		}
	}
}

func (g *Game) Check() bool {
	for i := range g.Cube.Body {
		y := g.Cube.Y + i
		for j := range g.Cube.Body[i] {
			if g.Cube.Body[i][j] == 1 {
				x := g.Cube.X + j
				if g.Backgroud[y][x] == 1 {
					return false
				}
			}
		}
	}
	return true
}

func (g *Game) GameOver() {

}
func (g *Game) Stash() {
	for i := range g.Backgroud {
		copy(g.Backgroud[i], g.View[i])
	}
}

func (g *Game) HandleInput(key rune) {
	switch key {
	case 'a':
		g.Move(-1, 0)
		if !g.Check() {
			g.Move(1, 0)
		}
	case 'd':
		g.Move(1, 0)
		if !g.Check() {
			g.Move(-1, 0)
		}
	case 'w':
		g.Cube.Rotate(90)
		if !g.Check() {
			g.Cube.Rotate(-90)
		}
	case 's':
		g.speed = 100
	}
}

func (g *Game) Refresh() {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("\033[H")
	for i := range g.View {
		copy(g.View[i], g.Backgroud[i])
	}
	for i := range g.Cube.Body {
		for j := range g.Cube.Body[i] {
			g.View[g.Cube.Y+i][g.Cube.X+j] += g.Cube.Body[i][j]
		}
	}
	for i := range g.View[:g.Height] {
		for j := range g.View[i] {
			if g.View[i][j] == 0 {
				buf.WriteString("  ")
			} else if g.View[i][j] == 1 {
				buf.WriteString("██")
			} else {
				buf.WriteString("░░")
			}
		}
		buf.WriteString("\n")
	}

	buf.WriteString("next: \n")
	for i := range g.NextCube.Body {
		for j := range g.NextCube.Body[i] {
			if g.NextCube.Body[i][j] == 0 {
				buf.WriteString("  ")
			} else {
				buf.WriteString("██")
			}
		}
		buf.WriteString("\n")
	}
	buf.WriteString(fmt.Sprintf("\nscore: %d", g.Score))
	os.Stdout.Write(buf.Bytes())
}
func (g *Game) Speed() time.Duration {
	return time.Millisecond * time.Duration(g.speed)
}

func (g *Game) GameStart(ctx context.Context) {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	keyChan := make(chan rune)
	go func() {
		defer keyboard.Close()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				char, key, _ := keyboard.GetKey()
				if char != 0 {
					keyChan <- char
				}
				if key == keyboard.KeyCtrlC {
					return
				}
			}
		}
	}()
	for {
		select {
		case <-time.Tick(g.Speed()):
			g.Move(0, 1)
			if !g.Check() {
				g.SettleScores()
				g.Stash()
				g.GenernateCube()
				if !g.Check() {
					g.Refresh()
					g.GameOver()
					return
				}
			}
			g.Refresh()
		case key := <-keyChan:
			g.HandleInput(key)
			g.Refresh()
		case <-ctx.Done():
			fmt.Println("game exit")
			return
		}
	}
}
func NewGame(Height, Width int) *Game {
	Cube := Cube{
		Width:  3,
		Height: 3,
	}
	Height += 3
	Width += 6
	game := Game{
		Cube:   &Cube,
		Height: Height,
		Width:  Width,
	}
	game.View = make([][]int, game.Height)
	game.Backgroud = make([][]int, game.Height)
	for i := range game.View {
		game.View[i] = make([]int, game.Width)
		game.Backgroud[i] = make([]int, game.Width)
	}
	game.Reset()
	game.Next()
	game.GenernateCube()
	fmt.Print("\033[H\033[2J")
	return &game
}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	game := NewGame(10, 16)
	game.GameStart(ctx)

}
