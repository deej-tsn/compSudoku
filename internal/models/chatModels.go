package models

import (
	"strings"

	"github.com/deej-tsn/compSudoku/internal/helper"
)

type (
	Message struct {
		Text   string
		Author string
		Color  string
	}

	JSONMessage struct {
		Text   string                 `json:"message"`
		Header map[string]interface{} `json:"HEADERS"`
	}

	JSONSignIn struct {
		Username string                 `json:"username"`
		Color    string                 `json:"radioColor"`
		Header   map[string]interface{} `json:"HEADERS"`
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

func NewChatLogFromFile(filePath string) *ChatLog {
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

func NewChatLog() *ChatLog {
	chatlog := ChatLog{
		Messages: []*Message{},
	}
	return &chatlog
}
