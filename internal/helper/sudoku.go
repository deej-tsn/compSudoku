package helper

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/deej-tsn/compSudoku/internal/models"
)

func checkError(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func ReadFileToString(filepath string) string {
	file, err := os.ReadFile(filepath)
	checkError(err)
	return string(file)
}

func NewGrid(filename string) models.Grid {
	grid := make([]models.Row, 9)
	rowsString := strings.Split(filename, "\n")
	for i := 0; i < len(rowsString); i++ {
		grid[i] = stringToRow(rowsString[i])
	}
	return grid
}

func createGivenSquare(valueOfString int) *models.Square {
	square := models.Square{
		Value:     valueOfString,
		Confirmed: true,
		Potential: []int{},
	}
	return &square
}

func createUnknownSquare() *models.Square {
	square := models.Square{
		Value:     0,
		Confirmed: false,
		Potential: []int{},
	}
	return &square
}

func stringToRow(rowString string) models.Row {
	row := make([]*models.Square, 9)
	removeSpaces := strings.ReplaceAll(rowString, " ", "")
	elements := strings.Split(removeSpaces, ",")
	for i := 0; i < len(elements); i++ {
		var number *models.Square

		if elements[i] != "_" {
			conver, err := strconv.Atoi(elements[i])
			checkError(err)
			number = createGivenSquare(conver)
		} else {
			number = createUnknownSquare()
		}
		row[i] = number
	}
	return row
}
