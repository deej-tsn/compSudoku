package models

import (
	"strings"

	"github.com/deej-tsn/compSudoku/internal/helper"
)

type (
	Message struct {
		Text   string
		Author string
	}

	ChatLog []*Message
)

func newMessage(s string) *Message {
	message := Message{
		Text:   s,
		Author: "dempsey",
	}
	return &message
}

func NewChatLog(filePath string) *ChatLog {
	file := helper.ReadFileToString(filePath)
	fileLines := strings.Split(file, ";")
	chatLog := []*Message{}
	for i := 0; i < len(fileLines); i++ {
		message := newMessage(fileLines[i])
		chatLog = append(chatLog, message)
	}
	return (*ChatLog)(&chatLog)
}
