package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"

	"github.com/deej-tsn/compSudoku/internal/helper"
	"github.com/deej-tsn/compSudoku/internal/models"
	componentsChat "github.com/deej-tsn/compSudoku/web/components/chat"
	componentsSudoku "github.com/deej-tsn/compSudoku/web/components/sudokuBoard"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketHandler struct {
	game    *models.Game
	chatLog *models.ChatLog
}

type WebSocketMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type WebSocketConnection struct {
	*websocket.Conn
	SignedIn bool
	Username string
	Color    string
}

var (
	webSocketConnections = make(map[*WebSocketConnection]bool)
	webSocketUpgrader    = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func NewWebSocketHandler(game *models.Game, chatLog *models.ChatLog) *WebSocketHandler {
	return &WebSocketHandler{
		game:    game,
		chatLog: chatLog,
	}
}

func (wsHandler *WebSocketHandler) InitWs(c echo.Context) error {
	conn, err := webSocketUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	currentConn := WebSocketConnection{
		Conn:     conn,
		SignedIn: false,
	}

	webSocketConnections[&currentConn] = true
	defer func() {
		delete(webSocketConnections, &currentConn)
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
			break
		}

		if currentConn.SignedIn {

			var wsMsg WebSocketMessage
			if err := json.Unmarshal(msg, &wsMsg); err != nil {
				c.Logger().Error(err)
				continue
			}

			switch wsMsg.Type {
			case "chat_message":
				wsHandler.handleChatMessage(c, &currentConn, wsMsg.Data)

			case "new_board":
				wsHandler.handleNewBoard(c, conn)

			case "set_active_square":
				wsHandler.handleSetActiveSquare(c, conn, wsMsg.Data)

			case "set_square_value":
				wsHandler.handleSetSquareValue(c, conn, wsMsg.Data)

			case "flip_edit_mode":
				wsHandler.handleFlipEditMode(c, conn)
			}

		} else {
			wsHandler.handleSignIn(c, conn, currentConn, msg)
		}

	}

	return nil
}

func (wsHandle *WebSocketHandler) handleSignIn(c echo.Context, conn *websocket.Conn, currentConn WebSocketConnection, data []byte) {
	var SignedInReq models.JSONSignIn
	err := json.Unmarshal(data, &SignedInReq)
	if err != nil {
		c.Logger().Error(err)
		return
	}
	delete(connections, currentConn)

	currentConn = WebSocketConnection{Conn: conn, Username: SignedInReq.Username, Color: SignedInReq.Color, SignedIn: true}
	connections[currentConn] = true

	UserConnected(c, currentConn)
}

func (wsHandler *WebSocketHandler) handleChatMessage(c echo.Context, currentConn *WebSocketConnection, data string) {
	var messageReq models.JSONMessage
	if err := json.Unmarshal([]byte(data), &messageReq); err != nil {
		c.Logger().Error(err)
		return
	}

	if messageReq.Text == "" {
		return
	}

	message := *models.NewMessage(messageReq.Text, currentConn.Username, currentConn.Color)
	wsHandler.chatLog.Messages = append(wsHandler.chatLog.Messages, &message)
	bytes := wsHandler.renderChatMessage(c, message, false)
	broadcastToAll(c, bytes)
}

func (wsHandler *WebSocketHandler) handleNewBoard(c echo.Context, conn *websocket.Conn) {
	newGrid, err := helper.GetBoardAPI()
	if err != nil {
		conn.WriteJSON(map[string]string{"error": "Failed to fetch new board"})
		return
	}
	wsHandler.game = models.ResponseToGame(newGrid)
	wsHandler.broadcastSudokuState(c)
}

func (wsHandler *WebSocketHandler) handleSetActiveSquare(c echo.Context, conn *websocket.Conn, data string) {
	position := helper.StringToPosition(data)
	wsHandler.game = wsHandler.game.SetActiveSquare(wsHandler.game.Grid[position[0]][position[1]])
	wsHandler.broadcastSudokuState(c)
}

func (wsHandler *WebSocketHandler) handleSetSquareValue(c echo.Context, conn *websocket.Conn, data string) {
	value, err := strconv.Atoi(data)
	if err != nil {
		conn.WriteJSON(map[string]string{"error": "Invalid value"})
		return
	}
	wsHandler.game = wsHandler.game.SetActiveSquareValue(value, c)
	wsHandler.broadcastSudokuState(c)
}

func (wsHandler *WebSocketHandler) handleFlipEditMode(c echo.Context, conn *websocket.Conn) {
	wsHandler.game.EditState = !wsHandler.game.EditState
	wsHandler.broadcastSudokuState(c)
}

func (wsHandler *WebSocketHandler) renderChatMessage(c echo.Context, message models.Message, selfSent bool) []byte {
	messageComp := new(bytes.Buffer)
	err := componentsChat.Message(message, selfSent).Render(context.Background(), messageComp)
	if err != nil {
		c.Logger().Error(err)
		return []byte{}
	}
	return messageComp.Bytes()
}

func (wsHandler *WebSocketHandler) broadcastSudokuState(c echo.Context) {
	gridBytes := wsHandler.renderSudokuGrid(c)
	broadcastToAll(c, gridBytes)
}

func (wsHandler *WebSocketHandler) renderSudokuGrid(c echo.Context) []byte {
	gridBuffer := new(bytes.Buffer)
	err := componentsSudoku.Grid(wsHandler.game).Render(context.Background(), gridBuffer)
	if err != nil {
		c.Logger().Error(err)
		return []byte{}
	}
	return gridBuffer.Bytes()
}

func broadcastToAll(c echo.Context, message []byte) {
	for conn := range webSocketConnections {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			delete(webSocketConnections, conn)
			conn.Close()
			c.Logger().Error(err)
		}
	}
}
