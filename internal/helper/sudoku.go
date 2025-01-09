package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/deej-tsn/compSudoku/internal/models"
)

func StringToPosition(positionString string) []int {
	positionArray := strings.Split(positionString, ",")

	position := []int{}
	for i := 0; i < len(positionArray); i++ {
		number, err := strconv.Atoi(positionArray[i])
		CheckError(err)
		position = append(position, number)
	}
	return position
}

func GetBoardAPI() (*models.SudokuResponse, error) {
	difficulty := 2
	sudoku_key := os.Getenv("SUDOKU_API_KEY")
	fmt.Println(sudoku_key)

	url := fmt.Sprintf("https://sudoku-board.p.rapidapi.com/new-board?diff=%d&stype=list&solu=true", difficulty)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-key", sudoku_key)
	req.Header.Add("x-rapidapi-host", "sudoku-board.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var response models.SudokuResponse
	err := json.Unmarshal(body, &response)
	return &response, err

}
