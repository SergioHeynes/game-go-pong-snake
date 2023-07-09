package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const PaddleSymbol = 0x2588
const BallSymbol = 0x25CF
const GameFrameWith = 30
const GameFrameHeight = 15
const GameFrameSymbol = '|'

type GameObject struct {
	row, col, width, height int
	velRow, velCol          int
	symbol                  rune
}

var screen tcell.Screen

var isGamePaused bool
var debugLog string

var gameObjects []*GameObject

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {

	InitScreen()
	InitGameState()
	inputChan := InitUserInput()

	for {
		HandleUserInput(ReadInput(inputChan))
		updateState()
		drawState()

		time.Sleep(90 * time.Millisecond)
	}

	screen.Fini()
}

func InitScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err = screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
}

func InitUserInput() chan string {
	inputChan := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				inputChan <- ev.Name()
			}
		}
	}()
	return inputChan
}

func InitGameState() {
	gameObjects = []*GameObject{}
}

func drawState() {
	if isGamePaused {
		return
	}
	screen.Clear()
	PrintString(0, 0, debugLog)
	for _, obj := range gameObjects {
		PrintFilledRect(obj.row, obj.col, obj.width, obj.height, obj.symbol)
	}
	screen.Show()
}

func CollidesWithWall(obj *GameObject) bool {
	_, screenHeight := screen.Size()
	return obj.row+obj.velRow < 0 || obj.row+obj.velRow >= screenHeight
}

func CollidesWithPaddle(ball *GameObject, paddle *GameObject) bool {
	var collidesOnColumn bool

	if ball.col < paddle.col {
		collidesOnColumn = ball.col+ball.velCol >= paddle.col
	} else {
		collidesOnColumn = ball.col+ball.velCol <= paddle.col
	}

	return collidesOnColumn &&
		ball.row >= paddle.row &&
		ball.row < paddle.row+paddle.height

}

func updateState() {
	if isGamePaused {
		return
	}
	for i := range gameObjects {
		gameObjects[i].row += gameObjects[i].velRow
		gameObjects[i].col += gameObjects[i].velCol
	}
}

func ReadInput(inputChan chan string) string {
	var key string
	select {
	case key = <-inputChan:
	default:
		key = ""
	}

	return key
}

func HandleUserInput(key string) {
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	}
}

func PrintStringCenter(row, col int, str string) {
	col = col - len(str)/2
	PrintString(row, col, str)
}

func PrintString(row, col int, str string) {
	for _, c := range str {
		screen.SetContent(col, row, c, nil, tcell.StyleDefault)
		col += 1
	}
}

func PrintFilledRect(row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {

		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}

	}
}
