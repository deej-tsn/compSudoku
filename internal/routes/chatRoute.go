package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deej-tsn/compSudoku/internal/helper"
	"github.com/deej-tsn/compSudoku/internal/models"
	components "github.com/deej-tsn/compSudoku/web/components/chat"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type (
	ChatHandler struct {
		ChatLog *models.ChatLog
	}

	WebSocketConnection struct {
		*websocket.Conn
		Username string
	}

	SocketResponse struct {
		From    string
		Type    string
		Message string
	}
)

var (
	connections = make(map[WebSocketConnection]bool)
	upgrader    = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func UserConnected(c echo.Context, msg string) {
	message := models.Message{
		Text:   msg,
		Author: "system",
	}

	bytes := messageToComponentByte(c, message)
	go broadcast(c, bytes)
}

func (chatH ChatHandler) InitWs(c echo.Context) error {
	username := "dempsey"
	conn, _ := upgrader.Upgrade(c.Response(), c.Request(), nil)

	currentConn := WebSocketConnection{Conn: conn, Username: username}
	connections[currentConn] = true

	connected := username + " connected....."

	UserConnected(c, connected)
	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil || msgType == websocket.CloseGoingAway {
			c.Logger().Error(err)
			break
		}
		var thing models.JsonMessage
		err = json.Unmarshal(msg, &thing)
		if err != nil {
			c.Logger().Error(err)
		}

		fmt.Printf("%s %s: %s:%d\n", conn.RemoteAddr(), username, string(msg), msgType)

		if thing.Text == "" {
			continue
		}

		message := models.Message{
			Text:   thing.Text,
			Author: username,
		}

		bytes := messageToComponentByte(c, message)
		go broadcast(c, bytes)
	}

	delete(connections, currentConn)

	return currentConn.NetConn().Close()
}

func messageToComponentByte(c echo.Context, message models.Message) []byte {
	messageComp := new(bytes.Buffer)
	err := components.Message(message).Render(context.Background(), messageComp)
	if err != nil {
		c.Logger().Error(err)
		return []byte{}
	}
	return messageComp.Bytes()
}

func broadcast(c echo.Context, message []byte) {
	for v := range connections {
		if err := v.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			delete(connections, v)
			v.Close()
			c.Logger().Error(err)
		}
	}
}

func NewChatHander(chat models.ChatLog) *ChatHandler {
	ch := ChatHandler{
		ChatLog: &chat,
	}
	return &ch
}

func (chatH ChatHandler) GetMessages(c echo.Context) error {

	return helper.Render(c, http.StatusAccepted, components.Chatbox(*chatH.ChatLog))
}

func (chatH ChatHandler) PostMessage(c echo.Context) error {
	messageText := c.QueryParam("message")
	message := models.Message{
		Text:   messageText,
		Author: "dempsey",
	}
	chatH.ChatLog.Messages = append(chatH.ChatLog.Messages, &message)
	return helper.Render(c, http.StatusAccepted, components.Chatlist(chatH.ChatLog.Messages))
}
