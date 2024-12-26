package models

import (
	"strings"

	"github.com/deej-tsn/compSudoku/internal/helper"
)

type (
	Message struct {
		text   string
		author string
	}

	ChatLog []Message
)

func newMessage(s string) {

}

func newChatLog(filePath string) *ChatLog {
	file := helper.ReadFileToString(filePath)
	fileLines := strings.Split(file, ";")
	chatLog := []Message{}
	for i := 0; i < len(fileLines); i++ {
		message := newMessage(fileLines[i])
		chatLog = append(chatLog, fileLines)
	}
}
