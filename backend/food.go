package backend

import "github.com/google/uuid"

type Food struct {
	Identifier
	UUID     uuid.UUID
	Position Coordinates
	Value    int
}

func (f *Food) ID() uuid.UUID {
	return f.UUID
}
