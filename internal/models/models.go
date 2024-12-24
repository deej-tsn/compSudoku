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
)

func (grid Grid) Print() {
	for i := 0; i < len(grid); i++ {
		fmt.Println(grid[i])
	}
}

func (grid Grid) changeSquareValue(row int, column int) *Square {
	grid[row][column].Value = 9
	return grid[row][column]
}

func (grid Grid) flipSquareActiveState(row int, column int) *Square {
	grid[row][column].Active = !grid[row][column].Active
	return grid[row][column]
}

func (grid Grid) SetSquare(row int, column int, active bool, value int) *Square {
	var squarePointer *Square
	if active {
		squarePointer = grid.changeSquareValue(row, column)
	} else {
		squarePointer = grid.flipSquareActiveState(row, column)
	}
	return squarePointer
}
