package models

import (
	"fmt"
)

type (
	Square struct {
		RowIndex        int
		ColumnIndex     int
		Value           int
		Confirmed       bool
		Potential       []int
		Active          bool
		RelatedToActive bool
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

func (game Game) changeSquareValue(square *Square) {
	if !square.Confirmed {
		square.Value = 9
	}

}

func (game Game) changeActiveSquare(square *Square) Game {
	if game.ActiveSquare != nil {
		game.ActiveSquare.Active = false
		game.activeSquareRelated(false)
	}

	square.Active = true
	game.ActiveSquare = square
	game.activeSquareRelated(true)
	return game
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

	if game.ActiveSquare.Value != 0 {
		for i := 0; i < len(game.Grid); i++ {
			for j := 0; j < len(game.Grid[0]); j++ {
				if game.Grid[i][j].Value == activeSquare.Value {
					game.Grid[i][j].RelatedToActive = state
				}
			}
		}
	}

	activeSquare.RelatedToActive = false
}

func (game Game) SetSquare(square *Square, value int) *Game {
	fmt.Println(game.ActiveSquare)
	if game.ActiveSquare != nil && game.ActiveSquare == square {
		game.changeSquareValue(square)
	} else {
		if game.ActiveSquare != nil {
			fmt.Printf("change Active from (%d,%d) to (%d,%d)", game.ActiveSquare.RowIndex, game.ActiveSquare.ColumnIndex, square.RowIndex, square.ColumnIndex)
		}

		game = game.changeActiveSquare(square)
		fmt.Println(game.ActiveSquare)

	}
	return &game
}
