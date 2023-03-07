package backend

import "github.com/google/uuid"

type Game struct {
	Entities map[uuid.UUID]Identifier
	Score    map[uuid.UUID]int
}

func NewGame() *Game {
	game := Game{
		Entities: make(map[uuid.UUID]Identifier),
		Score:    make(map[uuid.UUID]int),
	}
	return &game
}

func (g *Game) AddPlayer(e Identifier) {
	g.Entities[e.ID()] = e
	g.Score[e.ID()] = 0
}

type Mover interface {
	Move(Coordinates)
}

type Identifier interface {
	ID() uuid.UUID
}
