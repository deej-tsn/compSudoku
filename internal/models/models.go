package models

import "fmt"

type (
	Square struct {
		RowIndex    int
		ColumnIndex int
		Value       int
		Confirmed   bool
		Potential   []int
		Active      bool
	}

	Row []*Square

	Grid []Row

	Game struct {
		Grid         Grid
		ActiveSquare Square
	}
)

func (grid Grid) Print() {
	for i := 0; i < len(grid); i++ {
		fmt.Println(grid[i])
	}
}

func (game Game) changeSquareValue(row int, column int) *Square {
	game.Grid[row][column].Value = 9
	return game.Grid[row][column]
}

func (game Game) flipSquareActiveState(row int, column int) *Square {
	game.Grid[row][column].Active = !game.Grid[row][column].Active
	return game.Grid[row][column]
}

func (game Game) SetSquare(row int, column int, active bool, value int) *Square {
	var squarePointer *Square
	if active {
		squarePointer = game.changeSquareValue(row, column)
	} else {
		squarePointer = game.flipSquareActiveState(row, column)
	}
	return squarePointer
}
