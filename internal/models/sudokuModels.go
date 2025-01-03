package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/deej-tsn/compSudoku/internal/helper"
)

type (
	Square struct {
		RowIndex          int
		ColumnIndex       int
		Value             int
		ActualValue       int
		Confirmed         bool
		Potential         []int
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
	}

	SudokuResponse struct {
		Response struct {
			Difficulty     string  `json:"difficulty"`
			Solution       [][]int `json:"solution"`
			UnsolvedSudoku [][]int `json:"unsolved-sudoku"`
		} `json:"response"`
	}
)

func (grid Grid) Print() {
	for i := 0; i < len(grid); i++ {
		fmt.Println(grid[i])
	}
}

func (game Game) changeSquareValue(value int) Game {
	if !game.ActiveSquare.Confirmed {
		game.ActiveSquare.Value = value
	}
	return game

}

func (game Game) ChangeActiveSquare(square *Square) Game {
	if game.ActiveSquare != nil {
		game.ActiveSquare.Active = false
		game.activeSquareRelated(false)
		game.activeSquareSameValue(false)
	}

	square.Active = true
	game.ActiveSquare = square
	game.activeSquareRelated(true)
	game.activeSquareSameValue(true)
	return game
}

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

	// Square

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

func (game Game) SetActiveSquare(square *Square) *Game {
	game = game.ChangeActiveSquare(square)
	return &game
}
func (game Game) SetActiveSquareValue(value int) *Game {
	game.activeSquareSameValue(false)
	game = game.changeSquareValue(value)
	game.activeSquareSameValue(true)

	return &game
}

func NewGame(filename string) *Game {
	grid := make([]Row, 9)
	rowsString := strings.Split(filename, "\n")
	for i := 0; i < len(rowsString); i++ {
		grid[i] = stringToRow(i, rowsString[i])
	}
	game := Game{
		Grid: grid,
	}
	return &game
}

func createSquare(rowIndex int, columnIndex int, value int, actualValue int) *Square {
	square := Square{
		RowIndex:          rowIndex,
		ColumnIndex:       columnIndex,
		Value:             value,
		ActualValue:       actualValue,
		Confirmed:         true,
		Potential:         []int{},
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

func createGivenSquare(rowIndex int, columnIndex int, valueOfString int) *Square {
	square := Square{
		RowIndex:          rowIndex,
		ColumnIndex:       columnIndex,
		Value:             valueOfString,
		Confirmed:         true,
		Potential:         []int{},
		Active:            false,
		Valid:             true,
		RelatedToActive:   false,
		SameValueToActive: false,
	}
	return &square
}

func createUnknownSquare(rowIndex int, columnIndex int) *Square {
	square := Square{
		RowIndex:          rowIndex,
		ColumnIndex:       columnIndex,
		Value:             0,
		Confirmed:         false,
		Potential:         []int{},
		Active:            false,
		Valid:             true,
		RelatedToActive:   false,
		SameValueToActive: false,
	}
	return &square
}

func stringToRow(rowIndex int, rowString string) Row {
	row := make([]*Square, 9)
	removeSpaces := strings.ReplaceAll(rowString, " ", "")
	elements := strings.Split(removeSpaces, ",")
	for i := 0; i < len(elements); i++ {
		var number *Square

		if elements[i] != "_" {
			conver, err := strconv.Atoi(elements[i])
			helper.CheckError(err)
			number = createGivenSquare(rowIndex, i, conver)
		} else {
			number = createUnknownSquare(rowIndex, i)
		}
		row[i] = number
	}
	return row
}

func ResponseToGame(response *SudokuResponse) *Game {
	grid := make([]Row, 9)
	for i := 0; i < len(response.Response.Solution); i++ {
		grid[i] = IntegerRowToSudokuRow(response.Response.UnsolvedSudoku[i], response.Response.Solution[i], i)
	}
	game := Game{
		Grid:       grid,
		Difficulty: response.Response.Difficulty,
		Mistakes:   0,
	}
	return &game
}

func IntegerRowToSudokuRow(unsolvedRow []int, solvedRow []int, rowIndex int) Row {
	row := make([]*Square, 9)
	for i := 0; i < 9; i++ {
		row[i] = createSquare(rowIndex, i, unsolvedRow[i], solvedRow[i])
	}
	return row
}
