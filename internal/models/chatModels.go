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

	ChatLog struct {
		Messages []*Message
	}
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
	messages := []*Message{}
	chatLog := ChatLog{}
	for i := 0; i < len(fileLines); i++ {
		message := newMessage(fileLines[i])
		messages = append(messages, message)
	}
	chatLog.Messages = messages
	return &chatLog
}
