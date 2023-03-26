package backend

import "github.com/google/uuid"

type EntityType int

const (
	PlayerEntity = iota
	ObstacleEntity
	FoodEntity
)

// obstance object
type Obstacle struct {
	Identifier
	Positioner
	Diplayer
	UUID            uuid.UUID
	CurrentPosition Coordinates
	Icon            string
}

func (o *Obstacle) ID() uuid.UUID {
	return o.UUID
}

func (o *Obstacle) Position() Coordinates {
	return o.CurrentPosition
}

func (o *Obstacle) Display() string {
	return o.Icon
}

type Coordinates struct {
	X int
	Y int
}
