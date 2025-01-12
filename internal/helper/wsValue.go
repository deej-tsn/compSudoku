package helper

import (
	"fmt"

	"github.com/deej-tsn/compSudoku/internal/models"
)

func webSocketFormat(typeOfMessage string, data string) string {
	ws := fmt.Sprintf("{\"type\" : \"%s\", \"data\" :  %s }", typeOfMessage, data)
	return ws
}

func ActiveSquareFormat(square models.Square) string {
	ws := fmt.Sprintf("{\"position\" : \"%d,%d\"}", square.RowIndex, square.ColumnIndex)

	return webSocketFormat("set_active_square", ws)
}

func SquareValueFormat() string {
	return webSocketFormat("set_square_value", "{}")
}

func FlipEditMode() string {
	return webSocketFormat("flip_edit_mode", "{}")
}

func chatMessage() string {
	return ""
}

func NumberValueFormat(number int) string {
	ws := fmt.Sprintf("{\"value\" : \"%d\"} ", number)
	return webSocketFormat("set_square_value", ws)
}

func NewBoardFormat() string {
	return webSocketFormat("new_board", "{}")
}
