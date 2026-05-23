package main

import (
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1200
	screenHeight = 800
	gridSize     = 20
	cellsX       = screenWidth / gridSize  // 60
	cellsY       = screenHeight / gridSize // 40
)

// Игровое состояние
var (
	grid             [cellsY][cellsX]bool
	nextGrid         [cellsY][cellsX]bool
	paused           = true
	speedLevel       = 0 // 0 = 1x, 1 = 2x, 2 = 4x
	speedMultipliers = []float64{1.0, 2.0, 4.0}
	lastUpdate       time.Time
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Conway's Game of Life - Лабораторная работа")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	randomizeGrid(0.2) // Заполняем 20% клеток живыми

	lastUpdate = time.Now()

	for !rl.WindowShouldClose() {
		handleInput()

		currentSpeed := speedMultipliers[speedLevel]
		updateInterval := time.Duration(float64(1000/3)/currentSpeed) * time.Millisecond // 3 FPS обновления симуляции

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

func handleInput() {
	// Управление симуляцией
	if rl.IsKeyPressed(rl.KeySpace) {
		paused = !paused
	}

	if rl.IsKeyPressed(rl.KeyR) {
		randomizeGrid(0.2)
	}

	if rl.IsKeyPressed(rl.KeyC) {
		clearGrid()
	}

	// Изменение скорости: 1, 2, 3 на основной клавиатуре или NumPad
	if rl.IsKeyPressed(rl.KeyOne) || rl.IsKeyPressed(rl.KeyKp1) {
		speedLevel = 0
	}
	if rl.IsKeyPressed(rl.KeyTwo) || rl.IsKeyPressed(rl.KeyKp2) {
		speedLevel = 1
	}
	if rl.IsKeyPressed(rl.KeyFour) || rl.IsKeyPressed(rl.KeyKp4) {
		speedLevel = 2
	}

	// Выход
	if rl.IsKeyPressed(rl.KeyEscape) {
		rl.CloseWindow()
	}

	// Взаимодействие с мышью
	mousePos := rl.GetMousePosition()
	cellX := int(mousePos.X) / gridSize
	cellY := int(mousePos.Y) / gridSize

	if cellX >= 0 && cellX < cellsX && cellY >= 0 && cellY < cellsY {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			grid[cellY][cellX] = true
		}
		if rl.IsMouseButtonPressed(rl.MouseRightButton) {
			grid[cellY][cellX] = false
		}
	}
}

func updateGame() {
	// Расчёт следующего поколения
	for y := 0; y < cellsY; y++ {
		for x := 0; x < cellsX; x++ {
			neighbors := countNeighbors(x, y)
			isAlive := grid[y][x]

			if isAlive && (neighbors == 2 || neighbors == 3) {
				nextGrid[y][x] = true
			} else if !isAlive && neighbors == 3 {
				nextGrid[y][x] = true
			} else {
				nextGrid[y][x] = false
			}
		}
	}

	// Копируем следующее поколение в текущее
	for y := 0; y < cellsY; y++ {
		for x := 0; x < cellsX; x++ {
			grid[y][x] = nextGrid[y][x]
		}
	}
}

func countNeighbors(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			// Тороидальные границы (мир закольцован)
			if nx < 0 {
				nx = cellsX - 1
			} else if nx >= cellsX {
				nx = 0
			}
			if ny < 0 {
				ny = cellsY - 1
			} else if ny >= cellsY {
				ny = 0
			}
			if grid[ny][nx] {
				count++
			}
		}
	}
	return count
}

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
}

func drawUI() {
	// Полупрозрачная панель внизу
	rl.DrawRectangle(0, screenHeight-120, screenWidth, 120, rl.Fade(rl.Black, 0.85))

	// Информация о состоянии
	statusText := "PAUSED"
	statusColor := rl.Red
	if !paused {
		statusText = "RUNNING"
		statusColor = rl.Green
	}

	speedText := "SPEED: 1x"
	if speedLevel == 1 {
		speedText = "SPEED: 2x"
	} else if speedLevel == 2 {
		speedText = "SPEED: 4x"
	}

	rl.DrawText(statusText, 20, screenHeight-90, 30, statusColor)
	rl.DrawText(speedText, 20, screenHeight-55, 30, rl.White)

	// Подсказки по управлению
	rl.DrawText("[SPACE] Play/Pause", screenWidth-350, screenHeight-90, 20, rl.LightGray)
	rl.DrawText("[1/2/4] Speed: 1x/2x/4x", screenWidth-350, screenHeight-65, 20, rl.LightGray)
	rl.DrawText("[R] Randomize", screenWidth-350, screenHeight-40, 20, rl.LightGray)
	rl.DrawText("[C] Clear", screenWidth-350, screenHeight-15, 20, rl.LightGray)

	// Подсказки по мыши
	rl.DrawText("LMB: Create cell | RMB: Delete cell", 20, screenHeight-20, 20, rl.LightGray)
}

func randomizeGrid(density float32) {
	for y := 0; y < cellsY; y++ {
		for x := 0; x < cellsX; x++ {
			grid[y][x] = rand.Float32() < density
		}
	}
}

func clearGrid() {
	for y := 0; y < cellsY; y++ {
		for x := 0; x < cellsX; x++ {
			grid[y][x] = false
		}
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
