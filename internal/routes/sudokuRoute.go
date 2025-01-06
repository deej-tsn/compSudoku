package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

func (sudoku *SudokuRoute) SetActiveSquare(c echo.Context) error {
	position := stringToPosition(c.FormValue("position"))
	sudoku.game = sudoku.game.SetActiveSquare(sudoku.game.Grid[position[0]][position[1]])
	return helper.Render(c, http.StatusAccepted, components.Grid(sudoku.game))
}

func (sudoku *SudokuRoute) SetActiveSquareValue(c echo.Context) error {
	value, err := strconv.Atoi(c.FormValue("value"))
	helper.CheckError(err)
	sudoku.game = sudoku.game.SetActiveSquareValue(value)
	return helper.Render(c, http.StatusAccepted, components.Grid(sudoku.game))
}

func (sudoku *SudokuRoute) GetBoard(c echo.Context) error {
	return helper.Render(c, http.StatusAccepted, components.Grid(sudoku.game))
}

func (sudoku *SudokuRoute) PostFlipEditMode(c echo.Context) error {
	sudoku.game.EditState = !sudoku.game.EditState
	return helper.Render(c, http.StatusAccepted, components.InputType(sudoku.game.EditState))
}

func GetBoardAPI() (*models.SudokuResponse, error) {
	difficulty := 2
	sudoku_key := os.Getenv("SUDOKU_API_KEY")

	url := fmt.Sprintf("https://sudoku-board.p.rapidapi.com/new-board?diff=%d&stype=list&solu=true", difficulty)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-key", sudoku_key)
	req.Header.Add("x-rapidapi-host", "sudoku-board.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))
	var response models.SudokuResponse
	err := json.Unmarshal(body, &response)
	return &response, err

}
