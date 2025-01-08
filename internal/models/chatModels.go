package models

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

func NewMessage(s string, author string, color string) *Message {
	message := Message{
		Text:   s,
		Author: author,
		Color:  color,
	}
	return &message
}

func NewChatLog() *ChatLog {
	chatlog := ChatLog{
		Messages: []*Message{},
	}
	return &chatlog
}
