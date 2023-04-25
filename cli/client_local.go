package main

import (
	"log"
	"time"

	"github.com/atamocius/gameloop"
	"github.com/dubravaj/go-game/backend"
	"github.com/dubravaj/go-game/frontend"
	"github.com/google/uuid"
)

type Client struct {
	ID     uuid.UUID
	game   *backend.Game
	UIView *frontend.UIView
}

func (client *Client) Init() {
	game, err := backend.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	client.UIView = &frontend.UIView{}
	game.Init()
	client.game = game
}

func (client *Client) Run(player *backend.Player) {

	game := client.game
	client.UIView.Init(game, player)
	client.UIView.Run()
	config := gameloop.Config{

		TargetFPS: 60,

		IdleThreshold: 5,

		CurrentTimeFunc: func() float64 {
			return float64(time.Now().UnixNano()) * 1e-9
		},

		ProcessInputFunc: func() bool {
			return false
		},

		UpdateFunc: func(dt float64) {

			go func() {

				command := <-game.CommandsChan
				command.Execute(game)

			}()

		},

		RenderFunc: func() {
			client.UIView.Render()
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

func main() {

	client := Client{}
	client.Init()

	player := &backend.Player{UUID: uuid.New(), Icon: "\u25CF"}
	client.game.AddPlayer(player)

	client.Run(player)

	close(client.game.CommandsChan)

}
