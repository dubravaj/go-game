package backend

import "github.com/google/uuid"

type Food struct {
	Identifier
	Diplayer
	UUID            uuid.UUID
	CurrentPosition Coordinates
	Value           int
	Icon            string
}

func (f *Food) ID() uuid.UUID {
	return f.UUID
}

func (f *Food) Position() Coordinates {
	return f.CurrentPosition
}

func (f *Food) Display() string {
	return f.Icon
}
