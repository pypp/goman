package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pypp/goman/maps"
	"strconv"
	"time"
)

type Direction int
type Tile int
type GhostTile int

const AnimationLength = 3

const (
	Right Direction = iota + 1
	Left
	Up
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
	RightStartAnimation = 0
	LeftStartAnimation  = 3
	UpStartAnimation    = 6
	DownStartAnimation  = 9
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
	isSoundOn          bool = true
	isAttackMode       bool = false
	gameScore          int  = 0
	isForwardAnimation bool = true
	gameOverSound      rl.Sound
	music              rl.Music
	playerSprite       rl.Texture2D
	ghostsSprite       rl.Texture2D
	ghostPosX          float32 = 80
	ghostPosY          float32 = 80
	mapSprite          rl.Texture2D
	direction          Direction = Right
	animationSpeed     int32     = 8
	speed              float32   = 2
	playerPosX         float32   = 16
	playerPosY         float32   = 16
	playerSpritePosX   float32   = 0
	frameCount         int32     = 0
	nextDirection      Direction
	gameMap            [][]int
	mapHeight          int
	mapWidth           int
	screenWidth        int32
	screenHeight       int32
)

func drawBottomText() {
	var textSound string

	if isSoundOn {
		textSound = "sound: on"
	} else {
		textSound = "sound: off"
	}

	// TODO: calcualte the position instead of using fixed sizes
	rl.DrawText("score: "+strconv.Itoa(gameScore), 0, screenHeight-32, 30, rl.Red)
	rl.DrawText(textSound, screenWidth-160, screenHeight-32, 30, rl.Red)
}

func drawMap() {
	for y, h := range gameMap {
		for x, cell := range h {
			rl.DrawTextureRec(mapSprite, rl.NewRectangle(float32(cell)*SPRITE_SIZE, 0, SPRITE_SIZE, SPRITE_SIZE), rl.NewVector2(float32(x)*SPRITE_SIZE, float32(y)*SPRITE_SIZE), rl.White)
		}
	}
}

func drawPlayer() {
	rl.DrawTextureRec(playerSprite, rl.NewRectangle(playerSpritePosX*SPRITE_SIZE, 0, SPRITE_SIZE, SPRITE_SIZE), rl.NewVector2(playerPosX, playerPosY), rl.White)
}

func drawGhost() {
	var ghostSpritePosX float32

	if isAttackMode {
		ghostSpritePosX = 2
	} else {
		ghostSpritePosX = 0
	}

	rl.DrawTextureRec(ghostsSprite, rl.NewRectangle(ghostSpritePosX*SPRITE_SIZE, 0, SPRITE_SIZE, SPRITE_SIZE), rl.NewVector2(ghostPosX, ghostPosY), rl.White)
}

func drawScene() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	drawPlayer()
	drawGhost()
	drawMap()
	drawBottomText()

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

		if gameMap[gridY][gridX] < int(Point) {
			return true // Wall detected
		}
	}
	return false // No collision
}

func input() {
	keyToDirection := map[int]Direction{
		rl.KeyRight: Right,
		rl.KeyLeft:  Left,
		rl.KeyUp:    Up,
		rl.KeyDown:  Down,
	}

	for key, direction := range keyToDirection {
		if rl.IsKeyDown(int32(key)) {
			nextDirection = direction
			break
		}
	}

	if rl.IsKeyPressed(rl.KeyM) {
		isSoundOn = !isSoundOn
	}
}

func handleMovement() {
	// Get current grid coordinates
	gridX := int(playerPosX / SPRITE_SIZE)
	gridY := int(playerPosY / SPRITE_SIZE)

	// Check if the player is walking over a tile with value 15, and change it to 17
	if gameMap[gridY][gridX] == int(Point) || gameMap[gridY][gridX] == int(Strawberry) {

		if gameMap[gridY][gridX] == int(Point) {
			gameScore = gameScore + 10 // if point
		} else {
			gameScore = gameScore + 50 // if strawberry
			isAttackMode = true
			time.AfterFunc(5*time.Second, func() { isAttackMode = false })
		}

		gameMap[gridY][gridX] = int(Empty)
	}

	// Check if the player can move in the nextDirection
	if nextDirection != direction {
		switch nextDirection {
		case Right:
			if !isCollision(playerPosX+speed, playerPosY) {
				direction = Right
				playerSpritePosX = RightStartAnimation
			}
		case Left:
			if !isCollision(playerPosX-speed, playerPosY) {
				direction = Left
				playerSpritePosX = LeftStartAnimation
			}
		case Up:
			if !isCollision(playerPosX, playerPosY-speed) {
				direction = Up
				playerSpritePosX = UpStartAnimation
			}
		case Down:
			if !isCollision(playerPosX, playerPosY+speed) {
				direction = Down
				playerSpritePosX = DownStartAnimation
			}
		}
	}

	// Move in the current direction
	switch direction {
	case Right:
		if !isCollision(playerPosX+speed, playerPosY) {
			playerPosX += speed
		}
	case Left:
		if !isCollision(playerPosX-speed, playerPosY) {
			playerPosX -= speed
		}
	case Up:
		if !isCollision(playerPosX, playerPosY-speed) {
			playerPosY -= speed
		}
	case Down:
		if !isCollision(playerPosX, playerPosY+speed) {
			playerPosY += speed
		}
	}
}

func playerAnimation() {
	if frameCount%animationSpeed == 0 {
		var nextFrame float32

		if isForwardAnimation {
			nextFrame = playerSpritePosX + 1
		} else {
			nextFrame = playerSpritePosX - 1
		}

		switch direction {
		case Right, Left, Up, Down:
			startAnimation := map[Direction]int{
				Right: RightStartAnimation,
				Left:  LeftStartAnimation,
				Up:    UpStartAnimation,
				Down:  DownStartAnimation,
			}[direction]

			if isForwardAnimation && nextFrame == float32(startAnimation)+float32(AnimationLength) ||
				!isForwardAnimation && nextFrame == float32(startAnimation)-1 {
				isForwardAnimation = !isForwardAnimation
			}
		}

		if isForwardAnimation {
			playerSpritePosX++
		} else {
			playerSpritePosX--
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
	rl.StopMusicStream(music)
	rl.PlaySound(gameOverSound)

	for rl.IsSoundPlaying(gameOverSound) {
		var (
			fontSize  int32 = 40
			text            = "Game Over!"
			textWidth       = rl.MeasureText(text, fontSize)
			x               = (screenWidth - textWidth) / 2
			y               = (screenHeight - fontSize) / 2
		)

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawText(text, x, y, fontSize, rl.Red)
		rl.EndDrawing()
	}
}

func main() {
	gameMap, mapHeight, mapWidth = maps.LoadMap("./assets/maps/one.map")
	screenWidth = int32(mapWidth) * int32(SPRITE_SIZE)
	screenHeight = int32(mapHeight)*int32(SPRITE_SIZE) + 30

	rl.InitWindow(screenWidth, screenHeight, "goman")
	rl.InitAudioDevice()
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	rl.SetExitKey(0)

	gameOverSound = rl.LoadSound("./assets/audio/gameOverSound.wav")
	music = rl.LoadMusicStream("./assets/audio/music.wav")
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
			if isAttackMode {
				// TODO: eat the ghost
				ghostPosX = -500
				gameScore = gameScore + 200
			} else {
				drawGameOverScene()
				return
			}
		} else {
			drawScene()
			handleMovement()
			input()
			playerAnimation()
		}

		if isSoundOn {
			rl.UpdateMusicStream(music)
		}
	}

	rl.UnloadSound(gameOverSound)
	rl.CloseAudioDevice()
	rl.StopMusicStream(music)
	rl.CloseAudioDevice()
	rl.CloseWindow()
}
