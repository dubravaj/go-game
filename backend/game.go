package backend

import (
	"log"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
)

type Game struct {
	Entities map[uuid.UUID]Identifier
	Score    map[uuid.UUID]int
	Map      tcell.Screen
	Food     map[int]map[int]Identifier
}

func NewGame() *Game {
	gameMap, _ := tcell.NewScreen()
	game := Game{
		Entities: make(map[uuid.UUID]Identifier),
		Score:    make(map[uuid.UUID]int),
		Food:     make(map[int]map[int]Identifier),
		Map:      gameMap,
	}
	return &game
}

func (g *Game) Init() {
	if err := g.Map.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func (g *Game) AddPlayer(e Identifier) {
	g.Entities[e.ID()] = e
	g.Score[e.ID()] = 0
}

func (g *Game) watchCollisions() {
	//
}

func (g *Game) GenerateFood() {
	timer := time.NewTimer(5 * time.Second)
	for {
		<-timer.C
		width, height := g.Map.Size()
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
		food := Food{UUID: uuid.New(), CurrentPosition: Coordinates{X: x, Y: y}, Value: 1, Icon: "⛄"}
		g.Entities[food.ID()] = &food
		timer.Reset(5 * time.Second)
	}
}

type Mover interface {
	Move(Coordinates)
}

type Identifier interface {
	ID() uuid.UUID
}

type Positioner interface {
	Position() Coordinates
}

type Diplayer interface {
	Display() string
	Position() Coordinates
}
