package main

import (
	"example/user/hello/maps"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Direction int
type Tile int
type GhostTile int

const (
	Right Direction = iota + 1
	Left
	Top
	Down
)

const (
	RedGhost GhostTile = iota
	OrangeGhost
	BlueGhost
	CyanGhost
	PinkGhost
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
	ghostPosX        float32 = 80
	ghostPosY        float32 = 80
	mapSprite        rl.Texture2D
	direction        Direction = Right
	animationSpeed   int32     = 8
	speed            float32   = 2
	playerPosX       float32   = 16
	playerPosY       float32   = 16
	playerSpritePosX float32   = 0
	frameCount       int32     = 0
	numberOfGhosts   int       = 4
	nextDirection    Direction
	gameMap          [][]int
	mapHeight        int
	mapWidth         int
	screenWidth      int32
	screenHeight     int32
)

func drawScene() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawTextureRec(playerSprite, rl.NewRectangle(playerSpritePosX*SPRITE_SIZE, 0, SPRITE_SIZE, SPRITE_SIZE), rl.NewVector2(playerPosX, playerPosY), rl.White)
	// XXX: ghost position is fixed for now
	rl.DrawTextureRec(ghostsSprite, rl.NewRectangle(0*SPRITE_SIZE, 0, SPRITE_SIZE, SPRITE_SIZE), rl.NewVector2(ghostPosX, ghostPosY), rl.White)

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
	gridX := int(playerPosX / SPRITE_SIZE)
	gridY := int(playerPosY / SPRITE_SIZE)

	// Check if the player is walking over a tile with value 15, and change it to 17
	if gameMap[gridY][gridX] == int(Point) {
		gameMap[gridY][gridX] = int(Empty)
	}

	// Check if the player can move in the nextDirection
	switch nextDirection {
	case Right:
		if !isCollision(playerPosX+speed, playerPosY) {
			direction = Right
		}
	case Left:
		if !isCollision(playerPosX-speed, playerPosY) {
			direction = Left
		}
	case Top:
		if !isCollision(playerPosX, playerPosY-speed) {
			direction = Top
		}
	case Down:
		if !isCollision(playerPosX, playerPosY+speed) {
			direction = Down
		}
	}

	// Move in the current direction
	switch direction {
	case Right:
		if !isCollision(playerPosX+speed, playerPosY) {
			playerPosX += speed
			playerSpritePosX = 0 // Set the sprite to the "right" animation frame
		}
	case Left:
		if !isCollision(playerPosX-speed, playerPosY) {
			playerPosX -= speed
			playerSpritePosX = 3 // Set the sprite to the "left" animation frame
		}
	case Top:
		if !isCollision(playerPosX, playerPosY-speed) {
			playerPosY -= speed
			playerSpritePosX = 6 // Set the sprite to the "up" animation frame
		}
	case Down:
		if !isCollision(playerPosX, playerPosY+speed) {
			playerPosY += speed
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

func drawGameOverScene() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	text := "Game Over!"
	var fontSize int32 = 40
	textWidth := rl.MeasureText(text, fontSize)
	textHeight := fontSize
	x := (screenWidth - textWidth) / 2
	y := (screenHeight - textHeight) / 2
	rl.DrawText(text, x, y, fontSize, rl.Red)
	rl.EndDrawing()
}

func main() {
	gameMap, mapHeight, mapWidth = maps.LoadMap("./assets/maps/one.map")
	screenWidth = int32(mapWidth) * int32(SPRITE_SIZE)
	screenHeight = int32(mapHeight) * int32(SPRITE_SIZE)

	rl.InitWindow(screenWidth, screenHeight, "goman")
	rl.InitAudioDevice()
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	rl.SetExitKey(0)

	gameOverSound := rl.LoadSound("./assets/audio/gameOverSound.wav")
	music := rl.LoadMusicStream("./assets/audio/music.wav")
	rl.SetMusicVolume(music, 0.5)
	rl.PlayMusicStream(music)

	playerSprite = rl.LoadTexture("./assets/player.png")
	ghostsSprite = rl.LoadTexture("./assets/ghosts.png")
	mapSprite = rl.LoadTexture("./assets/tile.png")

	for !rl.WindowShouldClose() {
		if isGameOver() {
			return
		}

		if playerPosX == ghostPosX && playerPosY == ghostPosY {
			rl.StopMusicStream(music)
			rl.PlaySound(gameOverSound)
			for rl.IsSoundPlaying(gameOverSound) {
				drawGameOverScene()
			}
			return
		} else {
			drawScene()
			handleMovement()
			input()
			playerAnimation()
		}

		rl.UpdateMusicStream(music)
	}

	rl.UnloadSound(gameOverSound)
	rl.CloseAudioDevice()
	rl.StopMusicStream(music)
	rl.CloseAudioDevice()
	rl.CloseWindow()
}
