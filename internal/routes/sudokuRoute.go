package routes

import (
	"context"
	"net/http"
	"strconv"

	"github.com/deej-tsn/compSudoku/internal/helper"
	"github.com/deej-tsn/compSudoku/internal/models"
	components "github.com/deej-tsn/compSudoku/web/components/sudokuBoard"
	"github.com/labstack/echo/v4"
)

type SudokuRoute struct {
	game *models.Game
}

func NewSudokuRoute(game *models.Game) *SudokuRoute {
	return &SudokuRoute{
		game: game,
	}
}

func (sudoku *SudokuRoute) POSTActiveSquare(c echo.Context) error {
	position := helper.StringToPosition(c.FormValue("position"))
	sudoku.game = sudoku.game.SetActiveSquare(sudoku.game.Grid[position[0]][position[1]])
	return helper.Render(c, http.StatusAccepted, components.Grid(sudoku.game))
}

func (sudoku *SudokuRoute) GETNewBoard(c echo.Context) error {
	newGrid, err := helper.GetBoardAPI()
	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}
	sudoku.game = models.ResponseToGame(newGrid)
	return helper.Render(c, http.StatusAccepted, components.Game(sudoku.game))
}

func (sudoku *SudokuRoute) POSTActiveSquareValue(c echo.Context) error {
	value, err := strconv.Atoi(c.FormValue("value"))
	helper.CheckError(err)
	sudoku.game = sudoku.game.SetActiveSquareValue(value, c)
	for i := range 9 {
		if sudoku.game.NumbersLeft[i] <= 0 {
			components.Number(i+1, true).Render(context.Background(), c.Response().Writer)
		}
	}
	return helper.Render(c, http.StatusAccepted, components.Grid(sudoku.game))
}

func (sudoku *SudokuRoute) GETBoard(c echo.Context) error {
	return helper.Render(c, http.StatusAccepted, components.Grid(sudoku.game))
}

func (sudoku *SudokuRoute) POSTFlipEditMode(c echo.Context) error {
	sudoku.game.EditState = !sudoku.game.EditState
	return helper.Render(c, http.StatusAccepted, components.InputType(sudoku.game.EditState))
}
