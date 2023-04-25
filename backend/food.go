package backend

import "github.com/google/uuid"

type Food struct {
	UUID            uuid.UUID
	CurrentPosition Coordinates
	Value           int
	Icon            string
}

func New(position Coordinates, value int, icon string) Food {
	return Food{UUID: uuid.New(), CurrentPosition: position, Value: value, Icon: icon}
}

func (f *Food) ID() uuid.UUID {
	return f.UUID
}

func (f *Food) FoodValue() int {
	return f.Value
}

func (f *Food) Position() Coordinates {
	return f.CurrentPosition
}

func (f *Food) Display() string {
	return f.Icon
}
