package main

import (
	"example/user/hello/maps"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Direction int
type Tile int

const (
	Right Direction = iota + 1
	Left
	Top
	Down
)

const (
	HorizontalWall Tile = iota
	VerticalWall
	TopLeftCorner
	TopRightCorner
	BottomRightCorner
	BottomLeftCorner
	LeftSideWall
	RightSideWall
	BottomSideWall
	TopSideWall
	LeftEndWall
	RightEndWall
	BottomEndWall
	TopEndWall
	Apple
	Point
	Strawberry
	Empty
)

const (
	SPRITE_SIZE float32 = 16
)

var (
	playerSprite     rl.Texture2D
	ghostsSprite     rl.Texture2D
	mapSprite        rl.Texture2D
	direction        Direction = Right
	animationSpeed   int32     = 8
	speed            float32   = 2
	posX             float32   = 16
	posY             float32   = 16
	playerSpritePosX float32   = 0
	frameCount       int32     = 0
	numberOfGhosts   int       = 4
	nextDirection    Direction
	gameMap          [][]int
	mapHeight        int
	mapWidth         int
)

func drawScene() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawTextureRec(playerSprite, rl.NewRectangle(playerSpritePosX*SPRITE_SIZE, 0, SPRITE_SIZE, SPRITE_SIZE), rl.NewVector2(posX, posY), rl.White)

	// draw map
	for y, h := range gameMap {
		for x, cell := range h {
			rl.DrawTextureRec(mapSprite, rl.NewRectangle(float32(cell)*SPRITE_SIZE, 0, SPRITE_SIZE, SPRITE_SIZE), rl.NewVector2(float32(x)*SPRITE_SIZE, float32(y)*SPRITE_SIZE), rl.White)
		}
	}

	rl.EndDrawing()
}

func isCollision(x, y float32) bool {
	corners := []struct {
		offsetX, offsetY float32
	}{
		{0, 0},                                 // Top-left
		{SPRITE_SIZE - 0.5, 0},                 // Top-right
		{0, SPRITE_SIZE - 0.5},                 // Bottom-left
		{SPRITE_SIZE - 0.5, SPRITE_SIZE - 0.5}, // Bottom-right
	}

	for _, corner := range corners {
		gridX := int((x + corner.offsetX) / SPRITE_SIZE)
		gridY := int((y + corner.offsetY) / SPRITE_SIZE)

		if gridX < 0 || gridX >= len(gameMap[0]) || gridY < 0 || gridY >= len(gameMap) {
			return true // Out of bounds
		}

		// Check if the tile is either 15 or 17 (both are walkable)
		if gameMap[gridY][gridX] != int(Point) && gameMap[gridY][gridX] != int(Empty) {
			return true // Wall detected
		}
	}
	return false // No collision
}

func input() {
	// Update nextDirection based on key press
	if rl.IsKeyDown(rl.KeyRight) {
		nextDirection = Right
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		nextDirection = Left
	}
	if rl.IsKeyDown(rl.KeyUp) {
		nextDirection = Top
	}
	if rl.IsKeyDown(rl.KeyDown) {
		nextDirection = Down
	}
}

func handleMovement() {
	// Get current grid coordinates
	gridX := int(posX / SPRITE_SIZE)
	gridY := int(posY / SPRITE_SIZE)

	// Check if the player is walking over a tile with value 15, and change it to 17
	if gameMap[gridY][gridX] == int(Point) {
		gameMap[gridY][gridX] = int(Empty)
	}

	// Check if the player can move in the nextDirection
	switch nextDirection {
	case Right:
		if !isCollision(posX+speed, posY) {
			direction = Right
		}
	case Left:
		if !isCollision(posX-speed, posY) {
			direction = Left
		}
	case Top:
		if !isCollision(posX, posY-speed) {
			direction = Top
		}
	case Down:
		if !isCollision(posX, posY+speed) {
			direction = Down
		}
	}

	// Move in the current direction
	switch direction {
	case Right:
		if !isCollision(posX+speed, posY) {
			posX += speed
			playerSpritePosX = 0 // Set the sprite to the "right" animation frame
		}
	case Left:
		if !isCollision(posX-speed, posY) {
			posX -= speed
			playerSpritePosX = 3 // Set the sprite to the "left" animation frame
		}
	case Top:
		if !isCollision(posX, posY-speed) {
			posY -= speed
			playerSpritePosX = 6 // Set the sprite to the "up" animation frame
		}
	case Down:
		if !isCollision(posX, posY+speed) {
			posY += speed
			playerSpritePosX = 9 // Set the sprite to the "down" animation frame
		}
	}
}

func playerAnimation() {
	if frameCount%animationSpeed == 0 {
		nextFrame := playerSpritePosX + 1
		switch nextFrame {
		case 3:
			playerSpritePosX = 0
		case 6:
			playerSpritePosX = 3
		case 9:
			playerSpritePosX = 6
		case 12:
			playerSpritePosX = 9
		default:
			playerSpritePosX++
		}
	}

	frameCount++
}

func isGameOver() bool {
	for _, h := range gameMap {
		for _, cell := range h {
			if cell == int(Point) {
				return false
			}
		}
	}

	return true
}

func main() {
	gameMap, mapHeight, mapWidth = maps.LoadMap("./assets/maps/one.map")
	screenWidth := int32(mapWidth) * int32(SPRITE_SIZE)
	screenHeight := int32(mapHeight) * int32(SPRITE_SIZE)

	rl.InitWindow(screenWidth, screenHeight, "goman")
	rl.InitAudioDevice()
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	rl.SetExitKey(0)

	music := rl.LoadMusicStream("./assets/music.wav")
	rl.SetMusicVolume(music, 0.5)
	rl.PlayMusicStream(music)

	playerSprite = rl.LoadTexture("./assets/player.png")
	ghostsSprite = rl.LoadTexture("./assets/ghosts.png")
	mapSprite = rl.LoadTexture("./assets/tile.png")

	for !rl.WindowShouldClose() {
		if isGameOver() {
			return
		}

		rl.UpdateMusicStream(music)

		drawScene()
		handleMovement()
		input()
		playerAnimation()

	}

	rl.StopMusicStream(music)
	rl.CloseAudioDevice()
	rl.CloseWindow()
}
