package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"time"
)

type Cube struct {
	Body          [3][3]int
	X, Y          int
	Width, Height int
}

func (c *Cube) Rotate(angle int) {
	// 创建临时矩阵存储旋转结果
	temp := [3][3]int{}
	// 根据角度选择旋转方式
	switch angle {
	case 90: // 顺时针旋转90度
		for i := range c.Body {
			for j := range c.Body[i] {
				temp[j][2-i] = c.Body[i][j]
			}
		}
	case -90: // 逆时针旋转90度
		for i := range c.Body {
			for j := range c.Body[i] {
				temp[2-j][i] = c.Body[i][j]
			}
		}
	}

	// 将旋转后的结果复制回原矩阵
	c.Body = temp
}

type Montior struct {
	View      [][]int
	SnapView  [][]int
	Backgroud [][]int
	Cube      *Cube
	Height    int
	Width     int
}

func (m *Montior) Move(dx, dy int) {
	m.Cube.X += dx
	m.Cube.Y += dy

}
func (m *Montior) Clip() {
	if m.Cube.X < 0 {
		m.Cube.X = 0
	}
	if m.Cube.Y < 0 {
		m.Cube.Y = 0
	}
	if m.Cube.Y+m.Cube.Height > m.Height {
		m.Cube.Y = m.Height - m.Cube.Height
	}
	if m.Cube.X+m.Cube.Width > m.Width {
		m.Cube.X = m.Width - m.Cube.Width
	}
}

func (m *Montior) MoveAbsolute(x, y int) {
	m.Cube.X = x
	m.Cube.Y = y
	m.Clip()
}

func (m *Montior) RemoveFullLine() {
	row := m.Height - 2
	for {
		if !slices.Contains(m.View[row], 0) {
			for r := row; r > 0; r-- {
				copy(m.View[r], m.View[r-1])
			}
		} else {
			row--
			if row < 0 {
				return
			}
		}
	}
}

func (m *Montior) Check() bool {
	for i := range m.SnapView {
		copy(m.SnapView[i], m.Backgroud[i])
	}
	for i := range m.Cube.Body {
		for j := range m.Cube.Body[i] {
			m.SnapView[m.Cube.Y+i][m.Cube.X+j] += m.Cube.Body[i][j]
			if m.SnapView[m.Cube.Y+i][m.Cube.X+j] > 1 {
				return false
			}
		}
	}
	return true
}
func (m *Montior) Stash() {
	for i := range m.Backgroud {
		copy(m.Backgroud[i], m.View[i])
	}
}

func (m *Montior) Refresh() {
	buf := bytes.NewBuffer(nil)
	fmt.Fprint(buf, "\033[H")
	for i := range m.View {
		copy(m.View[i], m.Backgroud[i])
	}
	for i := range m.Cube.Body {
		for j := range m.Cube.Body[i] {
			m.View[m.Cube.Y+i][m.Cube.X+j] = m.Cube.Body[i][j]
		}
	}
	for i := range m.View[:m.Height-1] {
		for j := range m.View[i] {
			fmt.Fprint(buf, m.View[i][j])
		}
		fmt.Fprint(buf, "\n")
	}
	os.Stdout.Write(buf.Bytes())
}

func NewGame(Height, Width int) *Montior {
	Cube := Cube{
		Body: [3][3]int{
			{0, 0, 0},
			{1, 1, 1},
			{1, 1, 1},
		},
		X:      0,
		Y:      0,
		Width:  3,
		Height: 3,
	}
	Monitor := Montior{
		Cube:   &Cube,
		Height: Height,
		Width:  Width,
	}
	Monitor.View = make([][]int, Height)
	Monitor.SnapView = make([][]int, Height)
	Monitor.Backgroud = make([][]int, Height)
	for i := range Monitor.View {
		Monitor.View[i] = make([]int, Width)
		Monitor.SnapView[i] = make([]int, Width)
		Monitor.Backgroud[i] = make([]int, Width)
	}
	for i := range Monitor.Backgroud[0] {
		Monitor.Backgroud[Monitor.Height-1][i] = 1
	}
	return &Monitor
}
func main() {
	// if err := keyboard.Open(); err != nil {
	// 	panic(err)
	// }
	// defer keyboard.Close()
	// keyChan := make(chan keyboard.Key)
	// go func() {
	// 	for {
	// 		char, key, err := keyboard.GetKey()
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		if key != 0 {
	// 			keyChan <- key
	// 		} else {
	// 			keyChan <- keyboard.Key(char)
	// 		}
	// 	}
	// }()

	Monitor := NewGame(10, 3)
	for i := 0; i < 100; i++ {
		select {
		// case key := <-keyChan:
		// 	switch key {
		// 	case keyboard.KeyArrowLeft:
		// 		Monitor.Move(-1, 0)
		// 		if !Monitor.Check() {
		// 			Monitor.Move(1, 0)
		// 		}
		// 	case keyboard.KeyArrowRight:
		// 		Monitor.Move(1, 0)
		// 		if !Monitor.Check() {
		// 			Monitor.Move(-1, 0)
		// 		}
		// 	case keyboard.KeyArrowUp:
		// 		Monitor.Cube.Rotate(90)
		// 		if !Monitor.Check() {
		// 			Monitor.Cube.Rotate(-90)
		// 		}
		// 	}
		default:
			Monitor.Move(0, 1)
			if !Monitor.Check() {
				Monitor.RemoveFullLine()
				Monitor.Stash()

				Monitor.MoveAbsolute(0, 0)
				if !Monitor.Check() {

					return
				}
			}
		}
		Monitor.Refresh()
		time.Sleep(time.Millisecond * 500)
	}
}
