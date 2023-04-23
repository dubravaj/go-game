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

type Command interface {
	Execute(game *Game)
}

type AddPlayerCommand struct {
	ID   uuid.UUID
	Name string
	Icon string
}

func (command AddPlayerCommand) Execute(game *Game) {
	player := Player{UUID: command.ID, Icon: command.Icon, Name: command.Name}
	player.CurrentPosition = game.generateRandomPosition()
	game.Entities[player.ID()] = &player
	game.Score[player.ID()] = 0
}

type AddFoodCommand struct {
	ID    uuid.UUID
	Value int
	Icon  string
}

func (command AddFoodCommand) Execute(game *Game) {
	foodInitialPosition := game.generateRandomPosition()
	food := Food{UUID: command.ID, CurrentPosition: foodInitialPosition, Icon: command.Icon, Value: command.Value}
	game.Entities[command.ID] = &food
}

type RemoveEntityCommand struct {
	ID uuid.UUID
}

type MoveCommand struct {
	ID        uuid.UUID
	Direction Direction
	Timestamp int64 // number of miliseconds since start of the current round
}

func (command MoveCommand) Execute(game *Game) {

	//width, height := game.Map.Size()

	entity, ok := game.Entities[command.ID]
	if !ok {
		return
	}

	moverEntity, ok := entity.(Mover)
	if !ok {
		return
	}

	// if Mover, it has to be also Positioner, no need to check it
	position := entity.(Positioner).Position()

	switch command.Direction {
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
		game.Score[command.ID] += food.(Fooder).FoodValue()
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

	//if position.Y-1 > MapHeightOffset || position.Y-1 < height-MoveOffet || position.X-1 < width-MoveOffet || position.X-1 >= MoveOffet {
	moverEntity.Move(position)
	//}
}
