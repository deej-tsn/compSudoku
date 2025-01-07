package models

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type (
	Square struct {
		RowIndex          int
		ColumnIndex       int
		Value             int
		ActualValue       int
		Confirmed         bool
		Options           []bool
		Active            bool
		Valid             bool
		RelatedToActive   bool
		SameValueToActive bool
	}

	Row []*Square

	Grid []Row

	Game struct {
		Grid         Grid
		Difficulty   string
		Mistakes     int
		ActiveSquare *Square
		EditState    bool
		NumbersLeft  []int
	}

	SudokuResponse struct {
		Response struct {
			Difficulty     string  `json:"difficulty"`
			Solution       [][]int `json:"solution"`
			UnsolvedSudoku [][]int `json:"unsolved-sudoku"`
		} `json:"response"`
	}
)

// Related Squared

func (game Game) activeSquareSameValue(state bool) {
	if game.ActiveSquare.Value != 0 {
		for i := 0; i < len(game.Grid); i++ {
			for j := 0; j < len(game.Grid[0]); j++ {
				if game.Grid[i][j].Value == game.ActiveSquare.Value {
					game.Grid[i][j].SameValueToActive = state
				}
			}
		}
	}

	game.ActiveSquare.SameValueToActive = false

}

func (game Game) activeSquareRelated(state bool) {
	activeSquare := game.ActiveSquare
	// Row
	for i := 0; i < len(game.Grid); i++ {
		game.Grid[activeSquare.RowIndex][i].RelatedToActive = state
	}
	// Column
	for i := 0; i < len(game.Grid); i++ {
		game.Grid[i][activeSquare.ColumnIndex].RelatedToActive = state
	}

	// 3x3 Block
	minRowIndex := 3 * (activeSquare.RowIndex / 3)
	maxRowIndex := 3*(activeSquare.RowIndex/3) + 2

	minColumnIndex := 3 * (activeSquare.ColumnIndex / 3)
	maxColumnIndex := 3*(activeSquare.ColumnIndex/3) + 2

	for rowIndex := minRowIndex; rowIndex <= maxRowIndex; rowIndex++ {
		for columnIndex := minColumnIndex; columnIndex <= maxColumnIndex; columnIndex++ {
			game.Grid[rowIndex][columnIndex].RelatedToActive = state
		}
	}

	activeSquare.RelatedToActive = false
}

// Update Grid

func (game Game) SetActiveSquare(square *Square) *Game {
	if game.ActiveSquare != nil {
		game.ActiveSquare.Active = false
		game.activeSquareRelated(false)
		game.activeSquareSameValue(false)
	}

	square.Active = true
	game.ActiveSquare = square
	game.activeSquareRelated(true)
	game.activeSquareSameValue(true)
	return &game
}
func (game Game) SetActiveSquareValue(value int, c echo.Context) *Game {
	game.activeSquareSameValue(false)

	if !game.ActiveSquare.Confirmed {
		game.ActiveSquare.Value = value
		if value != game.ActiveSquare.ActualValue {
			game.Mistakes += 1
			fmt.Printf("Mistakes : %d\n", game.Mistakes)
		} else {
			game.ActiveSquare.Confirmed = true
			game.NumbersLeft[value-1] -= 1
			fmt.Println(game.NumbersLeft)

		}
	}

	game.activeSquareSameValue(true)

	return &game
}

// Create Sudoku Board

func makeNumberLeftSlice() []int {
	NumbersLeft := make([]int, 9)

	for i := 0; i < len(NumbersLeft); i++ {
		NumbersLeft[i] = 9
	}
	return NumbersLeft
}

func ResponseToGame(response *SudokuResponse) *Game {
	grid := make([]Row, 9)
	NumbersLeft := makeNumberLeftSlice()
	for i := 0; i < len(response.Response.Solution); i++ {
		grid[i] = IntegerRowToSudokuRow(response.Response.UnsolvedSudoku[i], response.Response.Solution[i], i, NumbersLeft)
	}
	game := Game{
		Grid:        grid,
		Difficulty:  response.Response.Difficulty,
		Mistakes:    0,
		EditState:   true,
		NumbersLeft: NumbersLeft,
	}
	return &game
}

func IntegerRowToSudokuRow(unsolvedRow []int, solvedRow []int, rowIndex int, numbersLeft []int) Row {
	row := make([]*Square, 9)
	for i := 0; i < 9; i++ {
		row[i] = createSquare(rowIndex, i, unsolvedRow[i], solvedRow[i])
		if unsolvedRow[i] == solvedRow[i] {
			index := unsolvedRow[i] - 1
			numbersLeft[index] = numbersLeft[index] - 1
		}
	}
	return row
}

func createSquare(rowIndex int, columnIndex int, value int, actualValue int) *Square {
	clearOptions := make([]bool, 9)
	square := Square{
		RowIndex:          rowIndex,
		ColumnIndex:       columnIndex,
		Value:             value,
		ActualValue:       actualValue,
		Confirmed:         true,
		Options:           clearOptions,
		Active:            false,
		Valid:             true,
		RelatedToActive:   false,
		SameValueToActive: false,
	}
	if value == 0 {
		square.Confirmed = false
	}
	return &square
}
