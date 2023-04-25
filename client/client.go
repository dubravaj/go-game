package client

import (
	"github.com/dubravaj/go-game/backend"
	"github.com/gdamore/tcell/v2"
)

type Client struct {
	Player   *backend.Player
	Entities backend.Diplayer
	Map      tcell.Screen
}

func (c *Client) watchCollisions() {

}

func (c *Client) MovePlayer() {

}
