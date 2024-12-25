package models

import (
	"fmt"
)

type (
	Square struct {
		RowIndex          int
		ColumnIndex       int
		Value             int
		Confirmed         bool
		Potential         []int
		Active            bool
		RelatedToActive   bool
		SameValueToActive bool
	}

	Row []*Square

	Grid []Row

	Game struct {
		Grid         Grid
		ActiveSquare *Square
	}
)

func (grid Grid) Print() {
	for i := 0; i < len(grid); i++ {
		fmt.Println(grid[i])
	}
}

func (game Game) changeSquareValue(square *Square, value int) Game {
	fmt.Println(square)
	if !square.Confirmed {
		fmt.Println("changing")
		square.Value = value
	}
	fmt.Println(square)
	return game

}

func (game Game) changeActiveSquare(square *Square) Game {
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

func (game Game) SetSquare(square *Square, value int) *Game {
	fmt.Println(value)
	if game.ActiveSquare != nil && game.ActiveSquare == square {
		game = game.changeSquareValue(square, value)
	} else {
		game = game.changeActiveSquare(square)

	}
	fmt.Println(game.ActiveSquare)
	return &game
}
