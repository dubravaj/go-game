package frontend

import (
	"fmt"
	"os"

	"github.com/dubravaj/go-game/backend"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UIAction struct {
}

type UIView struct {
	Game    *backend.Game
	App     *tview.Application
	Player  *backend.Player
	Views   *tview.Pages
	Actions chan UIAction
}

func (view *UIView) Render() {
}

func (view *UIView) Init() {

	view.App = tview.NewApplication()
	view.Views = tview.NewPages()

	view.Views.AddPage("Login", loginView(view), true, true)
	view.Views.AddPage("Game", gameView(view), true, false)

	if err := view.App.SetRoot(view.Views, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func (view *UIView) Run() {
	if err := view.App.Run(); err != nil {
		panic(err)
	}
}

func loginView(view *UIView) *tview.Form {
	form := tview.NewForm().
		AddInputField("Player name: ", "", 20, nil, nil).
		AddButton("Start", func() {
			// send event instead
			view.Views.SwitchToPage("Game")
		})
	form.SetBorder(true).SetTitle("Welcome to the Game").SetTitleAlign(tview.AlignCenter)
	return form
}

func gameView(view *UIView) *tview.Grid {

	gameView := tview.NewGrid().SetRows(2).SetColumns(200)

	game := tview.NewBox().SetBorder(true).SetBackgroundColor(tcell.ColorDarkBlue)
	infoBar := tview.NewBox().SetBorder(true).SetBackgroundColor(tcell.ColorDarkBlue).SetTitle("KOKOT")
	game.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		defStyle := tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault)
		screen.SetContent(20, 20, 'X', nil, defStyle)

		return 0, 0, 0, 0
	})
	gameView.AddItem(infoBar, 0, 0, 1, 1, 0, 0, false)
	gameView.AddItem(game, 1, 0, 1, 1, 0, 0, false)
	return gameView
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range text {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func renderBorders(gameMap tcell.Screen, style tcell.Style) {
	width, height := gameMap.Size()
	gameMap.SetContent(0, backend.MapHeightOffset-1, tcell.RuneULCorner, nil, style)
	gameMap.SetContent(0, height-1, tcell.RuneLLCorner, nil, style)
	gameMap.SetContent(width-1, backend.MapHeightOffset-1, tcell.RuneURCorner, nil, style)
	gameMap.SetContent(width-1, height-1, tcell.RuneLRCorner, nil, style)
	for i := 0; i < width; i++ {
		if i == 0 || i == width-1 {
			for j := backend.MapHeightOffset; j < height-1; j++ {
				gameMap.SetContent(i, j, tcell.RuneVLine, nil, style)
			}
		} else {

			gameMap.SetContent(i, backend.MapHeightOffset-1, tcell.RuneHLine, nil, style)
			gameMap.SetContent(i, height-1, tcell.RuneHLine, nil, style)
		}
	}
}

func renderGameEntities(game *backend.Game, style tcell.Style) {
	for _, entity := range game.Entities {
		displayerEntity, ok := entity.(backend.Diplayer)
		if !ok {
			continue
		}
		for _, r := range displayerEntity.Display() {
			currentPosition := displayerEntity.Position()
			game.Map.SetContent(currentPosition.X, currentPosition.Y, r, nil, style)
		}

	}
}

func Render(game *backend.Game) {
	game.Map.Clear()
	defStyle := tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault)
	game.Map.SetStyle(defStyle)

	text := "Welcome to the game"
	drawText(game.Map, 5, 3, len(text)+5, 3, defStyle, text)

	i := 0
	for _, score := range game.Score {
		drawText(game.Map, 5, 5, len(text)+5, 5, defStyle, fmt.Sprintf("Player %d score: ", i))
		drawText(game.Map, 21, 5, 23, 5, defStyle, fmt.Sprintf("%d", score))
		i++
	}

	renderBorders(game.Map, defStyle)
	renderGameEntities(game, defStyle)

	game.Map.Show()
}

func HandleInput(game *backend.Game, player *backend.Player) {

	width, height := game.Map.Size()
	currentPosition := player.Position()
	var moveCommand backend.Command

	if game.Map.HasPendingEvent() {

		switch event := game.Map.PollEvent().(type) {
		case *tcell.EventResize:
			game.Map.Sync()
		case *tcell.EventKey:
			switch event.Key() {
			case tcell.KeyEscape:
			case tcell.KeyCtrlC:
				game.Map.Fini()
				os.Exit(0)
			case tcell.KeyUp:
				if currentPosition.Y > backend.MapHeightOffset {
					moveCommand = backend.MoveCommand{ID: player.UUID, Direction: backend.Up}
					game.CommandsChan <- moveCommand
				}
			case tcell.KeyDown:
				if currentPosition.Y < height-backend.MoveOffet {
					moveCommand = backend.MoveCommand{ID: player.UUID, Direction: backend.Down}
					game.CommandsChan <- moveCommand
				}
			case tcell.KeyRight:
				if currentPosition.X < width-backend.MoveOffet {
					moveCommand = backend.MoveCommand{ID: player.UUID, Direction: backend.Right}
					game.CommandsChan <- moveCommand
				}
			case tcell.KeyLeft:
				if currentPosition.X >= backend.MoveOffet {
					moveCommand = backend.MoveCommand{ID: player.UUID, Direction: backend.Left}
					game.CommandsChan <- moveCommand
				}
			}
		}
	}
}
