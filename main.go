package main

import (
	"fmt"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1200
	screenHeight = 800
	gridSize     = 20
	uiHeight     = 100
	gameHeight   = screenHeight - uiHeight
	cellsX       = screenWidth / gridSize
	cellsY       = gameHeight / gridSize
)

var (
	grid             [cellsY][cellsX]bool
	nextGrid         [cellsY][cellsX]bool
	paused           = true
	speedLevel       = 0
	speedMultipliers = []float64{1, 2, 4, 8, 16, 32, 64}
	lastUpdate       time.Time
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Conway's Game of Life - Lab")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	randomizeGrid(0.2)
	lastUpdate = time.Now()

	for !rl.WindowShouldClose() {
		handleInput()

		currentSpeed := speedMultipliers[speedLevel]
		updateInterval := time.Duration(float64(1000/6)/currentSpeed) * time.Millisecond

		if !paused && time.Since(lastUpdate) >= updateInterval {
			updateGame()
			lastUpdate = time.Now()
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		drawGrid()
		drawUI()
		rl.EndDrawing()
	}
}

// Обработка клавиатуры и мыши
func handleInput() {
	if rl.IsKeyPressed(rl.KeySpace) {
		paused = !paused
	}
	if rl.IsKeyPressed(rl.KeyR) {
		randomizeGrid(0.2)
	}
	if rl.IsKeyPressed(rl.KeyC) {
		clearGrid()
	}
	if rl.IsKeyPressed(rl.KeyEscape) {
		rl.CloseWindow()
	}

	// Выбор скорости цифрами 1-7
	if rl.IsKeyPressed(rl.KeyOne) {
		speedLevel = 0
	}
	if rl.IsKeyPressed(rl.KeyTwo) {
		speedLevel = 1
	}
	if rl.IsKeyPressed(rl.KeyThree) {
		speedLevel = 2
	}
	if rl.IsKeyPressed(rl.KeyFour) {
		speedLevel = 3
	}
	if rl.IsKeyPressed(rl.KeyFive) {
		speedLevel = 4
	}
	if rl.IsKeyPressed(rl.KeySix) {
		speedLevel = 5
	}
	if rl.IsKeyPressed(rl.KeySeven) {
		speedLevel = 6
	}

	// Взаимодействие с сеткой через мышь
	mousePos := rl.GetMousePosition()
	cellX := int(mousePos.X) / gridSize
	cellY := int(mousePos.Y) / gridSize

	if cellY >= 0 && cellY < cellsY && cellX >= 0 && cellX < cellsX {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			grid[cellY][cellX] = true
		}
		if rl.IsMouseButtonPressed(rl.MouseRightButton) {
			grid[cellY][cellX] = false
		}
	}
}

// Расчет следующего поколения клеток
func updateGame() {
	for y := 0; y < cellsY; y++ {
		for x := 0; x < cellsX; x++ {
			neighbors := countNeighbors(x, y)
			alive := grid[y][x]

			if alive && (neighbors == 2 || neighbors == 3) {
				nextGrid[y][x] = true
			} else if !alive && neighbors == 3 {
				nextGrid[y][x] = true
			} else {
				nextGrid[y][x] = false
			}
		}
	}

	for y := 0; y < cellsY; y++ {
		for x := 0; x < cellsX; x++ {
			grid[y][x] = nextGrid[y][x]
		}
	}
}

// Подсчет живых соседей
func countNeighbors(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < cellsX && ny >= 0 && ny < cellsY {
				if grid[ny][nx] {
					count++
				}
			}
		}
	}
	return count
}

// Отрисовка игрового поля
func drawGrid() {
	for y := 0; y < cellsY; y++ {
		for x := 0; x < cellsX; x++ {
			if grid[y][x] {
				rl.DrawRectangle(
					int32(x*gridSize), int32(y*gridSize),
					int32(gridSize-1), int32(gridSize-1),
					rl.Green,
				)
			} else {
				rl.DrawRectangleLines(
					int32(x*gridSize), int32(y*gridSize),
					int32(gridSize), int32(gridSize),
					rl.DarkGray,
				)
			}
		}
	}
	rl.DrawLine(0, gameHeight, screenWidth, gameHeight, rl.Red)
}

// Отрисовка нижней панели управления
func drawUI() {
	rl.DrawRectangle(0, gameHeight, screenWidth, uiHeight, rl.NewColor(45, 45, 45, 255))
	rl.DrawRectangleLines(0, gameHeight, screenWidth, uiHeight, rl.Gray)

	liveCount := 0
	for y := 0; y < cellsY; y++ {
		for x := 0; x < cellsX; x++ {
			if grid[y][x] {
				liveCount++
			}
		}
	}

	// ===== КОЛОНКА 1: Статус и скорость =====
	statusText := "PAUSED"
	statusColor := rl.Red
	if !paused {
		statusText = "RUNNING"
		statusColor = rl.Green
	}
	rl.DrawText(statusText, 20, gameHeight+15, 22, statusColor)

	speedText := fmt.Sprintf("SPEED: %.0fx", speedMultipliers[speedLevel])
	rl.DrawText(speedText, 20, gameHeight+50, 22, rl.White)

	// ===== КОЛОНКА 2: Счетчик клеток =====
	rl.DrawText("CELLS:", 220, gameHeight+15, 18, rl.LightGray)
	liveText := fmt.Sprintf("%d / %d", liveCount, cellsX*cellsY)
	rl.DrawText(liveText, 220, gameHeight+45, 28, rl.Yellow)

	// ===== КОЛОНКА 3: Управление скоростью =====
	rl.DrawText("SPEED (1-7):", 420, gameHeight+15, 18, rl.LightGray)
	rl.DrawText("1=1x  2=2x  3=4x", 420, gameHeight+40, 18, rl.White)
	rl.DrawText("4=8x  5=16x 6=32x", 420, gameHeight+60, 18, rl.White)
	rl.DrawText("7=64x", 420, gameHeight+80, 18, rl.White)

	// ===== КОЛОНКА 4: Действия =====
	rl.DrawText("CONTROLS:", 680, gameHeight+15, 18, rl.LightGray)
	rl.DrawText("[SPACE] Play/Pause", 680, gameHeight+40, 18, rl.White)
	rl.DrawText("[R] Randomize", 680, gameHeight+60, 18, rl.White)
	rl.DrawText("[C] Clear", 680, gameHeight+80, 18, rl.White)

	// ===== КОЛОНКА 5: Управление мышью =====
	rl.DrawText("MOUSE:", 900, gameHeight+15, 18, rl.LightGray)
	rl.DrawText("LMB: Create", 900, gameHeight+40, 18, rl.White)
	rl.DrawText("RMB: Delete", 900, gameHeight+60, 18, rl.White)
}

// Случайное заполнение поля
func randomizeGrid(density float32) {
	for y := 0; y < cellsY; y++ {
		for x := 0; x < cellsX; x++ {
			grid[y][x] = rand.Float32() < density
		}
	}
}

// Очистка поля
func clearGrid() {
	for y := 0; y < cellsY; y++ {
		for x := 0; x < cellsX; x++ {
			grid[y][x] = false
		}
	}
}
