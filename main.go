package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/dubravaj/go-game/backend"
	"github.com/gdamore/tcell"
	"github.com/google/uuid"
)

// type Player struct {
// 	X      int
// 	Y      int
// 	Xspeed int
// 	Yspeed int
// }

// type Coords struct {
// 	X int
// 	Y int
// }

// func (p *Player) Display() string {
// 	return "\u25CF"
// }

// func (p *Player) RightMove() {
// 	p.X += p.Xspeed
// }

// func (p *Player) LeftMove() {
// 	p.X -= p.Xspeed
// }

// func (p *Player) UpMove() {
// 	p.Y -= p.Yspeed
// }

// func (p *Player) DownMove() {
// 	p.Y += p.Yspeed
// }

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range text {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func main() {

	game := backend.NewGame()

	screen, err := tcell.NewScreen()
	defStyle := tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	player := backend.Player{UUID: uuid.New(), CurrentPosition: backend.Coordinates{X: 27, Y: 14}, Icon: "\u25CF"}
	game.AddPlayer(&player)

	genChan := make(chan backend.Coordinates)
	timer := time.NewTimer(5 * time.Second)
	go func(genChan chan backend.Coordinates, screen tcell.Screen, timer *time.Timer) {
		for {
			<-timer.C
			width, height := screen.Size()
			y := rand.Intn(height)
			if y < 10 {
				y += 10
			}
			if y == height-1 {
				y -= 10
			}
			x := rand.Intn(width)
			if x == 0 {
				x += 10
			}
			if x == width-1 {
				x -= 10
			}
			coords := backend.Coordinates{X: x, Y: y}
			genChan <- coords
			timer.Reset(5 * time.Second)
			//time.Sleep(1000 * time.Millisecond)
		}

	}(genChan, screen, timer)

	go func(screen tcell.Screen, player *backend.Player, genChan chan backend.Coordinates) {
		var items []backend.Coordinates
		for {

			screen.Clear()
			defStyle := tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault)
			screen.SetStyle(defStyle)
			width, height := screen.Size()

			text := "Welcome to the game"
			drawText(screen, 5, 3, len(text)+5, 3, defStyle, text)

			i := 0
			for _, score := range game.Score {
				drawText(screen, 5, 5, len(text)+5, 5, defStyle, fmt.Sprintf("Player %d score: ", i))
				drawText(screen, 18, 5, len(text)+1, 5, defStyle, fmt.Sprintf("%d", score))
				i++
			}
			//drawText(screen, 13, 5, len(text)+1, 5, defStyle, string())

			screen.SetContent(0, 9, tcell.RuneULCorner, nil, defStyle)
			screen.SetContent(0, height-1, tcell.RuneLLCorner, nil, defStyle)
			screen.SetContent(width-1, 9, tcell.RuneURCorner, nil, defStyle)
			screen.SetContent(width-1, height-1, tcell.RuneLRCorner, nil, defStyle)
			for i := 0; i < width; i++ {
				if i == 0 || i == width-1 {
					for j := 10; j < height-1; j++ {
						screen.SetContent(i, j, tcell.RuneVLine, nil, defStyle)
					}
				} else {

					screen.SetContent(i, 9, tcell.RuneHLine, nil, defStyle)
					screen.SetContent(i, height-1, tcell.RuneHLine, nil, defStyle)
				}
			}

			screen.SetContent(20, 10, tcell.RuneHLine, nil, defStyle)
			screen.SetContent(21, 10, tcell.RuneHLine, nil, defStyle)
			screen.SetContent(22, 10, tcell.RuneHLine, nil, defStyle)
			screen.SetContent(23, 10, tcell.RuneURCorner, nil, defStyle)
			screen.SetContent(23, 11, tcell.RuneVLine, nil, defStyle)

			for _, r := range player.Icon {
				currentPosition := player.CurrentPosition
				screen.SetContent(currentPosition.X, currentPosition.Y, r, nil, defStyle)
			}

			for _, items := range items {
				screen.SetContent(items.X, items.Y, 'X', nil, defStyle)
			}

			go func(items *[]backend.Coordinates, genChan chan backend.Coordinates) {
				coord := <-genChan
				*items = append(*items, coord)
			}(&items, genChan)

			screen.Show()

			time.Sleep(40 * time.Millisecond)
		}

	}(screen, &player, genChan)

	ox, oy := -1, -1
	for {

		switch event := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			switch event.Key() {
			case tcell.KeyEscape:
			case tcell.KeyCtrlC:
				screen.Fini()
				os.Exit(0)
			case tcell.KeyUp:
				currentPosition := player.CurrentPosition
				player.Move(backend.Coordinates{X: currentPosition.X, Y: currentPosition.Y - 1})
			case tcell.KeyDown:
				currentPosition := player.CurrentPosition
				player.Move(backend.Coordinates{X: currentPosition.X, Y: currentPosition.Y + 1})
			case tcell.KeyRight:
				currentPosition := player.CurrentPosition
				player.Move(backend.Coordinates{X: currentPosition.X + 1, Y: currentPosition.Y})
			case tcell.KeyLeft:
				currentPosition := player.CurrentPosition
				player.Move(backend.Coordinates{X: currentPosition.X - 1, Y: currentPosition.Y})
			case tcell.KeyCtrlD:
				screen.SetContent(15, 15, 'W', nil, defStyle)
			}
		case *tcell.EventMouse:
			x, y := event.Position()
			switch event.Buttons() {

			case tcell.Button1, tcell.Button2:
				fmt.Println("HERE")
				if ox < 0 {
					ox, oy = x, y
				}
			case tcell.ButtonNone:
				if ox >= 0 {
					screen.SetContent(ox, oy, 'W', nil, tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault))
					fmt.Println("Click")
					ox, oy = -1, -1
				}

			}

		}

	}

}
