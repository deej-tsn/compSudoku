package helper

import (
	"os"
	"strconv"
	"strings"

	"github.com/deej-tsn/compSudoku/internal/models"
)

func ReadFileToString(filepath string) string {
	file, err := os.ReadFile(filepath)
	CheckError(err)
	return string(file)
}

func NewGame(filename string) models.Game {
	grid := make([]models.Row, 9)
	rowsString := strings.Split(filename, "\n")
	for i := 0; i < len(rowsString); i++ {
		grid[i] = stringToRow(i, rowsString[i])
	}
	game := models.Game{
		Grid: grid,
	}
	return game
}

func createGivenSquare(rowIndex int, columnIndex int, valueOfString int) *models.Square {
	square := models.Square{
		RowIndex:        rowIndex,
		ColumnIndex:     columnIndex,
		Value:           valueOfString,
		Confirmed:       true,
		Potential:       []int{},
		Active:          false,
		RelatedToActive: false,
	}
	return &square
}

func createUnknownSquare(rowIndex int, columnIndex int) *models.Square {
	square := models.Square{
		RowIndex:        rowIndex,
		ColumnIndex:     columnIndex,
		Value:           0,
		Confirmed:       false,
		Potential:       []int{},
		Active:          false,
		RelatedToActive: false,
	}
	return &square
}

func stringToRow(rowIndex int, rowString string) models.Row {
	row := make([]*models.Square, 9)
	removeSpaces := strings.ReplaceAll(rowString, " ", "")
	elements := strings.Split(removeSpaces, ",")
	for i := 0; i < len(elements); i++ {
		var number *models.Square

		if elements[i] != "_" {
			conver, err := strconv.Atoi(elements[i])
			CheckError(err)
			number = createGivenSquare(rowIndex, i, conver)
		} else {
			number = createUnknownSquare(rowIndex, i)
		}
		row[i] = number
	}
	return row
}
