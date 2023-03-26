package backend

import "github.com/google/uuid"

type Direction int

// move directions
const (
	Left Direction = iota
	Right
	Up
	Down
)

type Action interface {
	Do(game *Game)
}

type AddPlayerAction struct {
	ID   uuid.UUID
	Name string
	Icon string
}

func (action AddPlayerAction) Do(game *Game) {
	player := Player{UUID: action.ID, Icon: action.Icon, Name: action.Name}
	player.CurrentPosition = game.generateRandomPosition()
	game.Entities[player.ID()] = &player
	game.Score[player.ID()] = 0
}

type AddFoodAction struct {
	ID    uuid.UUID
	Value int
	Icon  string
}

func (action AddFoodAction) Do(game *Game) {
	foodInitialPosition := game.generateRandomPosition()
	food := Food{UUID: action.ID, CurrentPosition: foodInitialPosition, Icon: action.Icon, Value: action.Value}
	game.Entities[action.ID] = &food
}

type RemoveEntityAction struct {
	ID uuid.UUID
}

type MoveAction struct {
	ID uuid.UUID
	Direction
}

func (action MoveAction) Do(game *Game) {
	entity, ok := game.Entities[action.ID]
	if !ok {
		return
	}

	positionerEntity, ok := entity.(Positioner)
	if !ok {
		return
	}

	moverEntity, ok := entity.(Mover)
	if !ok {
		return
	}

	_ = moverEntity

	position := positionerEntity.Position()

	switch action.Direction {
	case Left:
		position.X--
	case Right:
		position.X++
	case Up:
		position.Y--
	case Down:
		position.Y++
	}

	// obstaclesMap := game.getEntityMap(ObstacleEntity)
	// obstacles := obstaclesMap[position]
	// if len(obstacles) > 0 {
	// 	return
	// }

	moverEntity.Move(position)
}
