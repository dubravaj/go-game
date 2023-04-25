package backend

import "github.com/google/uuid"

// player object
type Player struct {
	UUID            uuid.UUID
	CurrentPosition Coordinates
	PrevPosition    Coordinates
	Icon            string
	Name            string
}

func (p *Player) ID() uuid.UUID {
	return p.UUID
}

func (p *Player) Move(c Coordinates) {
	p.PrevPosition = p.CurrentPosition
	p.CurrentPosition = c
}

func (p *Player) Position() Coordinates {
	return p.CurrentPosition
}

func (p *Player) PreviousPosition() Coordinates {
	return p.PrevPosition
}

func (p *Player) Display() string {
	return p.Icon
}

func (p *Player) Collide(entity Positioner) bool {
	return p.Position() == entity.Position()
}
