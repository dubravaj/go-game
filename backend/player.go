package backend

import "github.com/google/uuid"

// player object
type Player struct {
	Identifier
	Mover
	Diplayer
	UUID            uuid.UUID
	CurrentPosition Coordinates
	Icon            string
}

func (p *Player) ID() uuid.UUID {
	return p.UUID
}

func (p *Player) Move(c Coordinates) {
	p.CurrentPosition = c
}

func (p *Player) Position() Coordinates {
	return p.CurrentPosition
}

func (p *Player) Display() string {
	return p.Icon
}
