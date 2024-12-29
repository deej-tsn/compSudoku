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
	connections = make([]*WebSocketConnection, 0)
	upgrader    = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func UserConnected(currentConn WebSocketConnection, connections []*WebSocketConnection, msg string) error {
	connectedBypeResp, _ := json.Marshal(SocketResponse{
		From:    "",
		Message: msg,
	})

	for _, v := range connections {
		if currentConn.Conn != v.Conn {
			fmt.Println("sending")
			if err := v.Conn.WriteMessage(1, connectedBypeResp); err != nil {
				return err
			}
		}
	}

	fmt.Println(msg)
	fmt.Println("Current Connection: ", len(connections))
	return nil
}

func (chatH ChatHandler) InitWs(c echo.Context) error {
	//username := c.FormValue("message")
	username := "dempsey"
	conn, _ := upgrader.Upgrade(c.Response(), c.Request(), nil)

	currentConn := WebSocketConnection{Conn: conn, Username: username}
	connections = append(connections, &currentConn)

	connected := username + " connected....."
	err := UserConnected(currentConn, connections, connected)
	if err != nil {
		c.Logger().Error(err)
	}
	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		var thing models.JsonMessage
		err = json.Unmarshal(msg, &thing)
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Println(thing)
		// Print the message to the console
		fmt.Printf("%s %s: %s:%d\n", conn.RemoteAddr(), username, string(msg), msgType)
		resp := SocketResponse{
			From:    currentConn.Username,
			Message: string(msg),
		}
		//byteResp, _ := json.Marshal(resp)
		fmt.Println(resp)
		message := models.Message{
			Text:   thing.Text,
			Author: "dempsey",
		}
		messageComp := new(bytes.Buffer)
		err = components.Message(message).Render(context.Background(), messageComp)
		if err != nil {
			c.Logger().Error(err)
		}
		for _, v := range connections {

			if err = v.Conn.WriteMessage(websocket.TextMessage, messageComp.Bytes()); err != nil {
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
