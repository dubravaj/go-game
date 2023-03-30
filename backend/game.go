package backend

import (
	"log"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
)

type Game struct {
	Entities    map[uuid.UUID]Identifier
	Score       map[uuid.UUID]int
	Map         tcell.Screen
	ActionsChan chan Action
}

func NewGame() *Game {
	gameMap, _ := tcell.NewScreen()
	game := Game{
		Entities:    make(map[uuid.UUID]Identifier),
		Score:       make(map[uuid.UUID]int),
		Map:         gameMap,
		ActionsChan: make(chan Action),
	}
	return &game
}

func (g *Game) Init() {
	if err := g.Map.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func (g *Game) getCollisionMap() map[Coordinates]Identifier {
	collisionMap := make(map[Coordinates]Identifier)

	for _, entity := range g.Entities {
		positionerEntity, _ := entity.(Positioner)
		collisionMap[positionerEntity.Position()] = entity
	}

	return collisionMap
}

func (g *Game) getEntityMap(entityType EntityType) map[Coordinates]Identifier {
	entityMap := make(map[Coordinates]Identifier)
	var positionerEntity Positioner

	for _, entity := range g.Entities {
		switch entityType {
		case PlayerEntity:
			_, ok := entity.(Mover)
			if !ok {
				continue
			}
		case FoodEntity:
			_, ok := entity.(Fooder)
			if !ok {
				continue
			}

		case ObstacleEntity:
			_, ok := entity.(Fooder)
			// food is also Positioner, but we want to allow move to it
			if ok {
				continue
			}
		}
		positionerEntity, _ = entity.(Positioner)
		entityMap[positionerEntity.Position()] = entity
	}

	return entityMap
}

func (g *Game) generateRandomPosition() Coordinates {
	width, height := g.Map.Size()
	collisionMap := g.getCollisionMap()

	for {

		x := rand.Intn(width-1) + 1
		y := rand.Intn(height-10) + 10

		randomPosition := Coordinates{X: x, Y: y}

		if _, found := collisionMap[randomPosition]; !found {
			return randomPosition
		}
	}
}

func (g *Game) AddPlayer(p *Player) {
	p.CurrentPosition = g.generateRandomPosition()
	g.Entities[p.ID()] = p
	g.Score[p.ID()] = 0
}

func (g *Game) AddEntity(e Identifier) {
	g.Entities[e.ID()] = e
}

func (g *Game) InitMap() {
	g.AddEntity(&Obstacle{UUID: uuid.New(), CurrentPosition: Coordinates{X: 10, Y: 20}, Icon: "A"})
	g.AddEntity(&Obstacle{UUID: uuid.New(), CurrentPosition: Coordinates{X: 10, Y: 21}, Icon: "A"})
	g.AddEntity(&Obstacle{UUID: uuid.New(), CurrentPosition: Coordinates{X: 29, Y: 25}, Icon: "A"})
	g.AddEntity(&Obstacle{UUID: uuid.New(), CurrentPosition: Coordinates{X: 50, Y: 22}, Icon: "A"})
}

func (g *Game) watchActions() {
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
		//"⛄"
		food := Food{UUID: uuid.New(), CurrentPosition: Coordinates{X: x, Y: y}, Value: 1, Icon: "X"}
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

type Fooder interface {
	FoodValue() int
}

type Diplayer interface {
	Display() string
	Position() Coordinates
}
