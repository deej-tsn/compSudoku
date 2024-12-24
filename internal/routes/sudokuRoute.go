package routes

import (
	"fmt"
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
		game: &models.Game{},
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

func stringToActive(activeString string) bool {
	active, err := strconv.ParseBool(activeString)
	helper.CheckError(err)
	return active
}

func (sudoku *SudokuRoute) UpdateGrid(c echo.Context) error {
	position := stringToPosition(c.FormValue("position"))
	active := stringToActive(c.FormValue("active"))
	fmt.Printf("row:%d column:%d active:%t\n", position[0], position[1], active)
	newSquare := sudoku.game.SetSquare(position[0], position[1], active, 9)
	return helper.Render(c, http.StatusAccepted, components.Square(*newSquare))
}
