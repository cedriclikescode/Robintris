package main

import (
	"embed"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
	"time"

	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed BACK.png BACK2.png BLOCK.png
var f embed.FS

var (
	blocksIMG *ebiten.Image
	backIMG *ebiten.Image
	back2IMG *ebiten.Image
	blockDirection int
	blockType int = rand.Intn(7)
	nextBlockType int = rand.Intn(7)
	blockX, blockY int = 3, 0
	blocksStruct [25][4][4]int
	blocks [10 * 23]bool
	keyToDAS int
	keyHoldDuration float64
	keyHoldProgress float64
	oldKeys = []ebiten.Key{}
	newKeys = []ebiten.Key{}
	currentTime time.Time
	forceMoveDown bool
	forceInvalid bool
	numberOfLinesCleared int
)

func init() {
	blocksStruct = [25][4][4]int{
	// T_BLOCK
	{
		{0, 1, 0, 0},
		{1, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 1, 0, 0},
		{0, 1, 1, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 0, 0, 0},
		{1, 1, 1, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 1, 0, 0},
		{1, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0},
	},
	// L1_BLOCK
	{
		{0, 0, 1, 0},
		{1, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 0, 0, 0},
		{1, 1, 1, 0},
		{1, 0, 0, 0},
		{0, 0, 0, 0},
	},
	{
		{1, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0},
	},
	// L2_BLOCK
	{
		{1, 0, 0, 0},
		{1, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 1, 1, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 0, 0, 0},
		{1, 1, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{1, 1, 0, 0},
		{0, 0, 0, 0},
	},
	// S1_BLOCK
	{
		{0, 1, 1, 0},
		{1, 1, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 1, 0, 0},
		{0, 1, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 0, 0, 0},
		{0, 1, 1, 0},
		{1, 1, 0, 0},
		{0, 0, 0, 0},
	},
	{
		{1, 0, 0, 0},
		{1, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0},
	},
	// S2_BLOCK
	{
		{1, 1, 0, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 0, 1, 0},
		{0, 1, 1, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 0, 0, 0},
		{1, 1, 0, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 1, 0, 0},
		{1, 1, 0, 0},
		{1, 0, 0, 0},
		{0, 0, 0, 0},
	},
	// BAR_BLOCK
	{
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
	{
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0},
	},
	{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{0, 0, 0, 0},
	},
	{
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0},
	},
	// CUBE_BLOCK
	{
		{0, 1, 1, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
	}
	var err error
	blocksIMG, _, err = ebitenutil.NewImageFromFileSystem(f, "BLOCK.png")
	if err != nil {
		log.Fatal(err)
	}
	backIMG, _, err = ebitenutil.NewImageFromFileSystem(f, "BACK.png")

	if err != nil {
		log.Fatal(err)
	}
	back2IMG, _, err = ebitenutil.NewImageFromFileSystem(f, "BACK2.png")
	if err != nil {
		log.Fatal(err)
	}
}

func isValid() bool {
	if forceInvalid {
		return false
	}

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if blocksStruct[blockType * 4 + blockDirection][y][x] == 1 && 
			(x + blockX < 0 || x + blockX > 9 || y + blockY > 22 || blocks[x + blockX + (blockY + y) * 10]) {
				return false
			}
		}
	}
	return true
}

type Game struct {}

func placeBlock(){
	//place block
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if blocksStruct[blockType * 4 + blockDirection][y][x] == 1 {
				blocks[x + blockX + (blockY + y) * 10] = true
			}
		}
	}
	//new block
	blockType = nextBlockType
	nextBlockType = rand.Intn(7)
	blockX, blockY = 3, 0
	blockDirection = 0
	//check for lines
	var x int
	for y := 22; y >= 0; y-- {
		for x = 0; x < 10; x++ {
			if !blocks[x + y * 10] {
				break //check failed
			}
		}
		if x == 10 {
			numberOfLinesCleared++
			for x = 0; x < 10; x++ {
				//blocks[x + y * 10] = false
				for i := y; i >= 1; i-- {
					blocks[x + i * 10] = blocks[x + i * 10 - 10]
					blocks[x + i * 10 - 10] = false
				}
			}
			y++
		}
	}
	forceInvalid = false
}

func setDASkey(i int) {
	keyToDAS = i
	keyHoldDuration = 0
	keyHoldProgress = 0
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
}

func MoveLeft(){
	blockX--
	if !isValid() {
		blockX++
	}
}

func MoveRight(){
	blockX++
	if !isValid() {
		blockX--
	}
}

func RotateClockwise(){
	blockDirection = (blockDirection + 1) % 4
	if !isValid() {
		if blockDirection == 0 {
			blockDirection = 3
		} else {
			blockDirection = (blockDirection - 1)
		}
	}
}

func MoveDown(){
	blockY++
	if !isValid() {
		blockY--
		placeBlock()
	}
	forceMoveDown = false
}

func isDifferentKeyPressed() bool {
	if len(newKeys) != len(oldKeys) {
		return true
	}

	for i := 0; i < len(newKeys); i++ {
		if newKeys[i] != oldKeys[i] {
			return true
		}
	}
	return false
}

func cheatCodeOutput() int {
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		return 0
	} else if inpututil.IsKeyJustPressed(ebiten.KeyF4) {
		return 1
	} else if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		return 2
	} else if inpututil.IsKeyJustPressed(ebiten.KeyF6) {
		return 3
	} else if inpututil.IsKeyJustPressed(ebiten.KeyF7) {
		return 4
	} else if inpututil.IsKeyJustPressed(ebiten.KeyF8) {
		return 5
	} else if inpututil.IsKeyJustPressed(ebiten.KeyF9) {
		return 6
	} else if inpututil.IsKeyJustPressed(ebiten.KeyF12) {
		return 7
	}
	return -1
}

func (g *Game) Update() error {
	timeDelta := time.Since(currentTime).Seconds()
	currentTime = time.Now()
	newKeys = inpututil.AppendPressedKeys([]ebiten.Key{})
	if isDifferentKeyPressed(){
		keyToDAS = 0
		ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMinimum)
	}

	oldKeys = newKeys

	keyHoldDuration += timeDelta

	if keyToDAS != 0 && keyHoldDuration > 0.225 {
		if keyHoldProgress == 0 {
			keyHoldProgress = (keyHoldDuration - 0.225) * 30
		} else {
			keyHoldProgress += timeDelta * 30
		}
	}

	for ; keyHoldProgress > 0; keyHoldProgress-- {
		switch keyToDAS {
			case 1:
				MoveLeft()
			case 2:
				MoveRight()
			case 3:
				RotateClockwise()
			case 4:
				MoveDown()
		} 
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		MoveLeft()
		setDASkey(1)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		MoveRight()
		setDASkey(2)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) && blockType < 6 {
		RotateClockwise()
		setDASkey(3)
	}


	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || forceMoveDown {
		MoveDown()
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			setDASkey(4)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape){
		os.Exit(0)
	}
	//cheat codes
	cheatCode := cheatCodeOutput()
	if cheatCode >= 0 {
		if cheatCode == 7 {
			forceInvalid = true
		} else {
			blockType = cheatCode
			blockX, blockY = 3, 0
			blockDirection = 0
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	//draw sides
	screen.DrawImage(backIMG, op)
	op.GeoM.Translate(420, 0)
	screen.DrawImage(back2IMG, op)

	op.GeoM.Translate(105, 67)
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if blocksStruct[nextBlockType * 4][y][x] == 1 {
				screen.DrawImage(blocksIMG.SubImage(image.Rect(0, 0, 20, 20)).(*ebiten.Image), op)
			}
			op.GeoM.Translate(20, 0)
		}
		op.GeoM.Translate(-80, 20)
	}

	//ebitenutil.DebugPrint(screen, fmt.Sprint(fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()), fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS())))
	ebitenutil.DebugPrintAt(screen, fmt.Sprint(numberOfLinesCleared), 450, 158)
	//draw blocks
	op.GeoM.Translate(-200 - 105, 20 - 80 - 67)
	for i := 0; i < 23; i++ {
		for n := 0; n < 10; n++ {
			if blocks[n + i * 10] {
				screen.DrawImage(blocksIMG.SubImage(image.Rect((n % 7) * 20, 0, (n % 7 + 1) * 20, 20)).(*ebiten.Image), op)
			}
			op.GeoM.Translate(20, 0)
		}
		op.GeoM.Translate(-200, 20)
	}
	op.GeoM.Translate(float64(blockX * 20), float64(20 * blockY) - 460)
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if blocksStruct[blockType * 4 + blockDirection][y][x] == 1 {
				screen.DrawImage(blocksIMG.SubImage(image.Rect(60, 0, 80, 20)).(*ebiten.Image), op)
			}
			op.GeoM.Translate(20, 0)
		}
		op.GeoM.Translate(-80, 20)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Robintris")
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMinimum)
	go func(){
		for {
			time.Sleep(1200 * time.Millisecond)
			forceMoveDown = true
			ebiten.ScheduleFrame()
		}
	}()
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
