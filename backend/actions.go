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
	ID        uuid.UUID
	Direction Direction
	Timestamp int64 // number of miliseconds since start of the current round
}

func (action MoveAction) Do(game *Game) {
	entity, ok := game.Entities[action.ID]
	if !ok {
		return
	}

	moverEntity, ok := entity.(Mover)
	if !ok {
		return
	}

	// if Mover, it has to be also Positioner, no need to check it
	position := entity.(Positioner).Position()

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

	foodMap := game.getEntityMap(FoodEntity)
	food, ok := foodMap[position]
	if ok {
		game.Score[action.ID] += food.(Fooder).FoodValue()
		// remove food from the map
		delete(game.Entities, food.ID())
		// update map in clients - send action to remove food from map

		moverEntity.Move(position)

		return
	}

	obstaclesMap := game.getEntityMap(ObstacleEntity)
	_, ok = obstaclesMap[position]
	if ok {
		return
	}

	moverEntity.Move(position)
}
