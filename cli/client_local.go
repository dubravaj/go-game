package main

import (
	"fmt"
	"os"
	"time"

	"github.com/atamocius/gameloop"
	"github.com/dubravaj/go-game/backend"
	"github.com/dubravaj/go-game/frontend"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
)

func main() {

	game := backend.NewGame()
	game.Init()

	player := backend.Player{UUID: uuid.New(), Icon: "\u25CF"}
	game.AddPlayer(&player)

	for i := 0; i < 10; i++ {
		game.AddEntity(&backend.Obstacle{UUID: uuid.New(), CurrentPosition: backend.Coordinates{X: 5 + i, Y: 15}, Icon: "A"})
	}

	config := gameloop.Config{
		// TargetFPS is used to calculate the seconds per update
		// (1 / TargetFPS).
		TargetFPS: 60,

		// IdleThreshold prevents updating the game if the time
		// elapsed since the previous frame exceeds this number (in seconds).
		IdleThreshold: 1,

		// CurrentTimeFunc accepts a function that returns the current time in
		// seconds. The gameloop library only provides a scaffold, it is up to
		// the user to provide an implementation. In this case, time's UnixNano
		// method was used but had to be multiplied by 0.000000001 to convert
		// to seconds.
		CurrentTimeFunc: func() float64 {
			return float64(time.Now().UnixNano()) * 1e-9
		},

		// ProcessInputFunc accepts a function that processes input logic
		// (ie. keyboard, mouse, gamepad, etc.) and returns a flag to signal the
		// game loop to quit.
		ProcessInputFunc: func() bool {
			width, height := game.Map.Size()
			currentPosition := player.Position()
			var moveAction backend.Action
			switch event := game.Map.PollEvent().(type) {
			case *tcell.EventResize:
				game.Map.Sync()
			case *tcell.EventKey:
				switch event.Key() {
				case tcell.KeyEscape:
				case tcell.KeyCtrlC:
					game.Map.Fini()
					os.Exit(0)
				case tcell.KeyUp:
					if currentPosition.Y > 10 {
						moveAction = backend.MoveAction{ID: player.UUID, Direction: backend.Up}
					}
				case tcell.KeyDown:
					if currentPosition.Y < height-2 {
						moveAction = backend.MoveAction{ID: player.UUID, Direction: backend.Down}

					}
				case tcell.KeyRight:
					if currentPosition.X < width-2 {
						moveAction = backend.MoveAction{ID: player.UUID, Direction: backend.Right}
					}
				case tcell.KeyLeft:
					if currentPosition.X >= 2 {
						moveAction = backend.MoveAction{ID: player.UUID, Direction: backend.Left}
					}
				default:
					return false
				}
				go func(action backend.Action) {
					game.ActionsChan <- moveAction
				}(moveAction)
			}
			return false
		},
		// UpdateFunc accepts a function that updates the game's state.
		// This function will be called based on a fixed interval
		// of 1 / TargetFPS (ie. 1 sec / 60 FPS = 0.01667 secs) and it is passed
		// as a parameter (dt).
		UpdateFunc: func(dt float64) {

			go func() {
				action := <-game.ActionsChan
				action.Do(game)
			}()
			// var removedEntities []backend.Identifier

			// for id, entity := range game.Entities {
			// 	if player.ID() != id {
			// 		positionerEntity, _ := entity.(backend.Positioner)
			// 		collision := player.Collide(positionerEntity)
			// 		if collision {
			// 			switch entity.(type) {
			// 			case backend.Fooder:
			// 				game.Score[player.ID()] += entity.(backend.Fooder).FoodValue()
			// 				removedEntities = append(removedEntities, entity)
			// 			default:
			// 				player.CurrentPosition = player.PreviousPosition()
			// 			}

			// 		}
			// 	}
			// }

			// for _, entity := range removedEntities {
			// 	delete(game.Entities, entity.ID())
			// }
		},

		// RenderFunc accepts a function that contains rendering logic.
		RenderFunc: func() {

			game.Map.Clear()
			defStyle := tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault)
			game.Map.SetStyle(defStyle)
			width, height := game.Map.Size()

			text := "Welcome to the game"
			frontend.DrawText(game.Map, 5, 3, len(text)+5, 3, defStyle, text)

			i := 0
			for _, score := range game.Score {
				frontend.DrawText(game.Map, 5, 5, len(text)+5, 5, defStyle, fmt.Sprintf("Player %d score: ", i))
				frontend.DrawText(game.Map, 21, 5, 23, 5, defStyle, fmt.Sprintf("%d", score))
				i++
			}

			game.Map.SetContent(0, 9, tcell.RuneULCorner, nil, defStyle)
			game.Map.SetContent(0, height-1, tcell.RuneLLCorner, nil, defStyle)
			game.Map.SetContent(width-1, 9, tcell.RuneURCorner, nil, defStyle)
			game.Map.SetContent(width-1, height-1, tcell.RuneLRCorner, nil, defStyle)
			for i := 0; i < width; i++ {
				if i == 0 || i == width-1 {
					for j := 10; j < height-1; j++ {
						game.Map.SetContent(i, j, tcell.RuneVLine, nil, defStyle)
					}
				} else {

					game.Map.SetContent(i, 9, tcell.RuneHLine, nil, defStyle)
					game.Map.SetContent(i, height-1, tcell.RuneHLine, nil, defStyle)
				}
			}

			for _, entity := range game.Entities {
				displayerEntity, ok := entity.(backend.Diplayer)
				if !ok {
					continue
				}
				for _, r := range displayerEntity.Display() {
					currentPosition := displayerEntity.Position()
					game.Map.SetContent(currentPosition.X, currentPosition.Y, r, nil, defStyle)
				}

			}

			game.Map.Show()
		},
	}

	// Call the gameloop.Create() function and pass the config to create
	// a game loop.
	runLoop := gameloop.Create(config)

	// generate food
	go game.GenerateFood()

	// Run the created game loop.
	runLoop()

}
