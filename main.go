package main

import (
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pypp/goman/maps"
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

type Player struct {
	posX float32
	posY float32
	// TODO: make object for player sprite
	spritePosX         float32
	direction          Direction
	nextDirection      Direction
	isForwardAnimation bool
}

type Ghost struct {
	posX float32
	posY float32
	// TODO: will just be color type
	spritePosX float32
}

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
	SPRITE_SIZE     float32 = 16
	ANIMATION_SPEED int32   = 8
	SPEED           float32 = 2
)

var (
	frameCount   int32 = 0
	music        rl.Music
	isMuted      bool = false
	isAttackMode bool = false
	gameScore    int  = 0
	gameMap      [][]int

	ghostPosX float32 = 80
	ghostPosY float32 = 80

	// TODO: move to player struct
	direction          Direction = Right
	playerPosX         float32   = 16
	playerPosY         float32   = 16
	playerSpritePosX   float32   = 0
	nextDirection      Direction
	isForwardAnimation bool = true
)

func drawBottomText() {
	var textSound string

	if isMuted {
		textSound = "sound: off"
	} else {
		textSound = "sound: on"
	}

	winWidth, winHeight := int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight())

	// TODO: calculate the position instead of using fixed sizes
	rl.DrawText("score: "+strconv.Itoa(gameScore), 0, winHeight-32, 30, rl.Red)
	rl.DrawText(textSound, winWidth-160, winHeight-32, 30, rl.Red)
}

func drawMap(mapSprite rl.Texture2D) {
	for y, h := range gameMap {
		for x, cell := range h {
			rl.DrawTextureRec(mapSprite, rl.NewRectangle(float32(cell)*SPRITE_SIZE, 0, SPRITE_SIZE, SPRITE_SIZE), rl.NewVector2(float32(x)*SPRITE_SIZE, float32(y)*SPRITE_SIZE), rl.White)
		}
	}
}

func drawPlayer(playerSprite rl.Texture2D) {
	rl.DrawTextureRec(playerSprite, rl.NewRectangle(playerSpritePosX*SPRITE_SIZE, 0, SPRITE_SIZE, SPRITE_SIZE), rl.NewVector2(playerPosX, playerPosY), rl.White)
}

func drawGhost(posY, poxY float32, ghostColor GhostTile, ghostSprite rl.Texture2D) {
	var tile float32

	if isAttackMode {
		tile = float32(BlueGhost)
	} else {
		tile = float32(ghostColor)
	}

	rl.DrawTextureRec(ghostSprite, rl.NewRectangle(float32(tile)*SPRITE_SIZE, 0, SPRITE_SIZE, SPRITE_SIZE), rl.NewVector2(posY, poxY), rl.White)
}

func drawScene(playerSprite, ghostsSprite, mapSprite rl.Texture2D) {

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	drawPlayer(playerSprite)

	drawMap(mapSprite)
	drawBottomText()

	drawGhost(80, 80, RedGhost, ghostsSprite)
	drawGhost(208, 16, OrangeGhost, ghostsSprite)
	drawGhost(16, 208, CyanGhost, ghostsSprite)
	drawGhost(160, 208, PinkGhost, ghostsSprite)

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
		isMuted = !isMuted
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
			if !isCollision(playerPosX+SPEED, playerPosY) {
				direction = Right
				playerSpritePosX = RightStartAnimation
			}
		case Left:
			if !isCollision(playerPosX-SPEED, playerPosY) {
				direction = Left
				playerSpritePosX = LeftStartAnimation
			}
		case Up:
			if !isCollision(playerPosX, playerPosY-SPEED) {
				direction = Up
				playerSpritePosX = UpStartAnimation
			}
		case Down:
			if !isCollision(playerPosX, playerPosY+SPEED) {
				direction = Down
				playerSpritePosX = DownStartAnimation
			}
		}
	}

	// Move in the current direction
	switch direction {
	case Right:
		if !isCollision(playerPosX+SPEED, playerPosY) {
			playerPosX += SPEED
		}
	case Left:
		if !isCollision(playerPosX-SPEED, playerPosY) {
			playerPosX -= SPEED
		}
	case Up:
		if !isCollision(playerPosX, playerPosY-SPEED) {
			playerPosY -= SPEED
		}
	case Down:
		if !isCollision(playerPosX, playerPosY+SPEED) {
			playerPosY += SPEED
		}
	}
}

func playerAnimation() {
	if frameCount%ANIMATION_SPEED == 0 {
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

func drawGameOverScene(gameOverSound rl.Sound) {
	winWidth, winHeight := int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight())
	rl.StopMusicStream(music)
	rl.PlaySound(gameOverSound)

	for rl.IsSoundPlaying(gameOverSound) {
		var (
			fontSize  int32 = 40
			text            = "Game Over!"
			textWidth       = rl.MeasureText(text, fontSize)
			x               = (winWidth - textWidth) / 2
			y               = (winHeight - fontSize) / 2
		)

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawText(text, x, y, fontSize, rl.Red)
		rl.EndDrawing()
	}
}

func main() {
	var mapHeight, mapWidth int
	gameMap, mapHeight, mapWidth = maps.LoadMap("./assets/maps/one.map")

	screenWidth := int32(mapWidth) * int32(SPRITE_SIZE)
	screenHeight := int32(mapHeight)*int32(SPRITE_SIZE) + 30

	rl.InitWindow(screenWidth, screenHeight, "goman")
	rl.InitAudioDevice()
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	rl.SetExitKey(0)

	gameOverSound := rl.LoadSound("./assets/audio/game_over_sound.wav.wav")
	wakaWakaMusic := rl.LoadMusicStream("./assets/audio/waka_waka.wav")
	mainMusic := rl.LoadMusicStream("./assets/audio/music.wav")
	playerSprite := rl.LoadTexture("./assets/sprites/player.png")
	ghostsSprite := rl.LoadTexture("./assets/sprites/ghosts.png")
	mapSprite := rl.LoadTexture("./assets/sprites/tile.png")

	for !rl.WindowShouldClose() {
		if isAttackMode {
			music = wakaWakaMusic
		} else {
			music = mainMusic
		}

		rl.SetMusicVolume(music, 0.5)
		rl.PlayMusicStream(music)

		if isGameOver() {
			return
		}

		if playerPosX == ghostPosX && playerPosY == ghostPosY {
			if isAttackMode {
				// TODO: eat the ghost
				gameScore = gameScore + 200
			} else {
				drawGameOverScene(gameOverSound)
				return
			}
		} else {
			drawScene(playerSprite, ghostsSprite, mapSprite)
			handleMovement()
			input()
			playerAnimation()
		}

		if !isMuted {
			rl.UpdateMusicStream(music)
		}
	}

	rl.UnloadSound(gameOverSound)
	rl.StopMusicStream(music)
	rl.CloseAudioDevice()
	rl.CloseWindow()
}
