package routes

import (
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
	connections = make([]*WebSocketConnection, 0)
	upgrader    = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func UserConnected(currentConn WebSocketConnection, connections []*WebSocketConnection, msg string) {
	connectedBypeResp, _ := json.Marshal(SocketResponse{
		From:    "",
		Message: msg,
	})

	for _, v := range connections {
		if currentConn.Conn != v.Conn {
			if err := v.Conn.WriteMessage(1, connectedBypeResp); err != nil {
				return
			}
		}
	}

	fmt.Println(msg)
	fmt.Println("Current Connection: ", len(connections))
}

func (chatH ChatHandler) InitWs(c echo.Context) error {
	username := c.FormValue("message")
	conn, _ := upgrader.Upgrade(c.Response(), c.Request(), nil)
	currentConn := WebSocketConnection{Conn: conn, Username: username}
	connections = append(connections, &currentConn)

	connected := username + " connected....."
	UserConnected(currentConn, connections, connected)

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		}

		// Print the message to the console
		fmt.Printf("%s %s: %s\n", conn.RemoteAddr(), username, string(msg))

		resp := SocketResponse{
			From:    currentConn.Username,
			Message: string(msg),
		}
		byteResp, _ := json.Marshal(resp)

		for _, v := range connections {
			if err = v.Conn.WriteMessage(msgType, byteResp); err != nil {
				return c.String(http.StatusTeapot, err.Error())
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
