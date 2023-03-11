package backend

import "github.com/google/uuid"

// obstance object
type Obstacle struct {
	UUID     uuid.UUID
	Position Coordinates
	Icon     rune
}

type Coordinates struct {
	X int
	Y int
}
