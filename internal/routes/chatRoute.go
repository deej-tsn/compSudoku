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
		SignedIn bool
		Username string
		Color    string
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

func UserConnected(c echo.Context, user WebSocketConnection) {
	systemMessage := fmt.Sprintf("%s has connected", user.Username)
	message := *models.NewMessage(
		systemMessage,
		"system",
		"purple",
	)

	bytes := messageToComponentByte(c, message, false)
	go broadcast(c, user, bytes)
}

func (chatH ChatHandler) InitWs(c echo.Context) error {
	conn, _ := upgrader.Upgrade(c.Response(), c.Request(), nil)

	currentConn := WebSocketConnection{Conn: conn, SignedIn: false}
	connections[currentConn] = true

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil || msgType == websocket.CloseGoingAway {
			c.Logger().Error(err)
			break
		}
		if currentConn.SignedIn {
			var messageReq models.JSONMessage
			err = json.Unmarshal(msg, &messageReq)
			if err != nil {
				c.Logger().Error(err)
			}

			fmt.Printf("%s %s: %s:%d\n", conn.RemoteAddr(), currentConn.Username, string(msg), msgType)

			if messageReq.Text == "" {
				continue
			}

			message := *models.NewMessage(messageReq.Text, currentConn.Username, currentConn.Color)

			bytes := messageToComponentByte(c, message, false)
			go broadcast(c, currentConn, bytes)

			bytes = messageToComponentByte(c, message, true)

			if err := currentConn.Conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
				c.Logger().Error(err)
				break
			}

		} else {
			var SignedInReq models.JSONSignIn
			err = json.Unmarshal(msg, &SignedInReq)
			if err != nil {
				c.Logger().Error(err)
			}
			delete(connections, currentConn)

			currentConn = WebSocketConnection{Conn: conn, Username: SignedInReq.Username, Color: SignedInReq.Color, SignedIn: true}
			connections[currentConn] = true

			UserConnected(c, currentConn)
		}

	}

	delete(connections, currentConn)

	return currentConn.NetConn().Close()
}

func messageToComponentByte(c echo.Context, message models.Message, selfSent bool) []byte {
	messageComp := new(bytes.Buffer)
	err := components.Message(message, selfSent).Render(context.Background(), messageComp)
	if err != nil {
		c.Logger().Error(err)
		return []byte{}
	}
	return messageComp.Bytes()
}

func broadcast(c echo.Context, sentFrom WebSocketConnection, message []byte) {
	for v := range connections {
		if v != sentFrom {
			if err := v.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				delete(connections, v)
				v.Close()
				c.Logger().Error(err)
			}
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
