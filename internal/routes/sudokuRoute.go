package routes

import (
	"net/http"
	"strconv"
	"strings"

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

func stringToPosition(positionString string) []int {
	positionArray := strings.Split(positionString, ",")

	position := []int{}
	for i := 0; i < len(positionArray); i++ {
		number, err := strconv.Atoi(positionArray[i])
		helper.CheckError(err)
		position = append(position, number)
	}
	return position
}

func (sudoku *SudokuRoute) UpdateGrid(c echo.Context) error {
	position := stringToPosition(c.FormValue("position"))
	sudoku.game = sudoku.game.SetSquare(sudoku.game.Grid[position[0]][position[1]], 9)
	return helper.Render(c, http.StatusAccepted, components.Game(*sudoku.game))
}
