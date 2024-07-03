package main

import (
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const BoardWidth = 100
const BoardHeight = 100

var Rules string = "0:0,0\n1:0,0\n2:3,1\n3:1,1\n4:5,0\n5:6,0\n6:7,0\n7:8,0\n8:0,0\n"

var zFactor float32 = 1.0

var cellSize float32 = 8

var padding float32 = (10*zFactor - cellSize) / 2

var PosX, PosY int

var Neighbors [][]int = [][]int{
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{0, -1},
	{0, 1},
	{1, -1},
	{1, 0},
	{1, 1},
}

var Pause bool = false

var StateToColor = map[int]rl.Color{
	0: rl.NewColor(100, 100, 110, 255),
	1: rl.NewColor(110, 110, 130, 255),
	2: rl.NewColor(120, 120, 140, 255),
	3: rl.NewColor(130, 130, 150, 255),
	4: rl.NewColor(140, 140, 160, 255),
	5: rl.NewColor(150, 150, 170, 255),
	6: rl.NewColor(160, 160, 180, 255),
	7: rl.NewColor(170, 170, 190, 255),
	8: rl.NewColor(180, 180, 200, 255),
}

type Cell struct {
	X, Y  int
	State int
}

func (c *Cell) Draw(zoomFactor float32) {
	// Adjust cell size based on zoom factor

	cellSize = 8 * zoomFactor
	padding = (10*zoomFactor - cellSize) / 2 // Adjust padding to keep cell centered

	xPosition := float32(c.X)*10*zoomFactor + padding
	yPosition := float32(c.Y)*10*zoomFactor + padding

	if c.State >= 1 {
		rl.DrawRectangle(int32(xPosition), int32(yPosition), int32(cellSize), int32(cellSize), StateToColor[c.State])
	} else {
		rl.DrawRectangle(int32(xPosition), int32(yPosition), int32(cellSize), int32(cellSize), rl.DarkGray)
	}

}
func (c *Cell) CheckNeighbours(b *Board) int {
	count := 0
	for _, n := range Neighbors {
		x := c.X + n[0]
		y := c.Y + n[1]
		if x >= 0 && x < BoardWidth && y >= 0 && y < BoardHeight {
			if b.Cells[x][y].State == 1 {
				count++
			}
		}
	}
	return count
}

type Board struct {
	Cells        [BoardWidth][BoardHeight]*Cell
	UpdateString string
	UpdateMap    map[int][2]int
}

func NewBoard() *Board {
	b := Board{}
	for x := 0; x < BoardWidth; x++ {
		for y := 0; y < BoardHeight; y++ {
			b.Cells[x][y] = &Cell{x, y, 0}
		}
	}
	b.UpdateMap = make(map[int][2]int)
	b.UpdateString = Rules
	b.ParseUpdateString(b.UpdateString)
	return &b
}

func (b *Board) ParseUpdateString(s string) {
	// Parse the update string
	// Separate by new line
	// Split by ":"
	// Split by ","
	// Part 0 is the state
	// Part 1 is the output dead
	// Part 2 is the output alive
	// Convert to int
	// Store in a map

	// Split by new line
	lines := strings.Split(s, "\n")
	for _, line := range lines[:len(lines)-1] {
		// Split by ":"
		parts := strings.Split(line, ":")
		// Convert to int
		state := parts[0]
		output := parts[1]

		stateI, err := strconv.Atoi(state)
		if err != nil {
			panic(err)
		}

		// Separate the dead and alive states
		outputParts := strings.Split(output, ",")
		dead := outputParts[0]
		alive := outputParts[1]

		deadI, err := strconv.Atoi(dead)
		if err != nil {
			panic(err)
		}

		aliveI, err := strconv.Atoi(alive)
		if err != nil {
			panic(err)
		}

		b.UpdateMap[stateI] = [2]int{deadI, aliveI}

	}

}

func (b *Board) Update() {
	/*
		Any live cell with fewer than two live neighbours dies, as if by underpopulation.
		Any live cell with two or three live neighbours lives on to the next generation.
		Any live cell with more than three live neighbours dies, as if by overpopulation.
	*/
	tempBoard := NewBoard()

	for x := 0; x < BoardWidth; x++ {
		for y := 0; y < BoardHeight; y++ {
			count := b.Cells[x][y].CheckNeighbours(b)
			if b.Cells[x][y].State >= 1 {
				// Alive case
				tempBoard.Cells[x][y].State = b.UpdateMap[count][1]
			} else {
				// Dead case
				tempBoard.Cells[x][y].State = b.UpdateMap[count][0]
			}
		}
	}

	// Apply the new states to the original board
	for x := 0; x < BoardWidth; x++ {
		for y := 0; y < BoardHeight; y++ {
			b.Cells[x][y].State = tempBoard.Cells[x][y].State
		}
	}
}

func main() {
	b := NewBoard()

	rl.InitWindow(1000, 1000, "Conway's Game of Life")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {

		if rl.IsKeyPressed(rl.KeySpace) {
			Pause = !Pause
		}

		mWheel := rl.GetMouseWheelMove()

		if mWheel > 0 {
			zFactor += 0.1
		} else if mWheel < 0 {
			if zFactor > 1 {
				zFactor -= 0.1
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		for x := 0; x < BoardWidth; x++ {
			for y := 0; y < BoardHeight; y++ {

				b.Cells[x][y].Draw(zFactor)
			}
		}

		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			if Pause {
				// Assuming the base cell size is 10, adjust this if your base cell size is different
				cellSize := 10 * zFactor
				x := int((float32(rl.GetMouseX()) / cellSize))
				y := int((float32(rl.GetMouseY()) / cellSize))

				// Ensure x and y are within the bounds of the board
				if x >= 0 && x < BoardWidth && y >= 0 && y < BoardHeight {
					b.Cells[x][y].State = 1
				}
			}
		}

		if rl.IsKeyDown(rl.KeyW) {
			PosY++
			for x := 0; x < BoardWidth; x++ {
				for y := 0; y < BoardHeight-1; y++ {
					b.Cells[x][y].State = b.Cells[x][y+1].State
				}
			}
		}

		if rl.IsKeyDown(rl.KeyS) {
			PosY--
			for x := 0; x < BoardWidth; x++ {
				for y := BoardHeight - 1; y > 0; y-- {
					b.Cells[x][y].State = b.Cells[x][y-1].State
				}
			}
		}

		if rl.IsKeyDown(rl.KeyA) {
			PosX++
			for x := 0; x < BoardWidth-1; x++ {
				for y := 0; y < BoardHeight; y++ {
					b.Cells[x][y].State = b.Cells[x+1][y].State
				}
			}
		}

		if rl.IsKeyDown(rl.KeyD) {
			PosX--
			for x := BoardWidth - 1; x > 0; x-- {
				for y := 0; y < BoardHeight; y++ {
					b.Cells[x][y].State = b.Cells[x-1][y].State
				}
			}
		}

		if rl.IsKeyPressed(rl.KeyR) {
			b = NewBoard()
		}

		if !Pause {
			b.Update()
		}
		rl.DrawFPS(10, 10)
		rl.DrawText("Press SPACE to pause", 10, 30, 20, rl.White)
		rl.EndDrawing()
	}
}
