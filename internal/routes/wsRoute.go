package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/deej-tsn/compSudoku/internal/helper"
	"github.com/deej-tsn/compSudoku/internal/models"
	componentsChat "github.com/deej-tsn/compSudoku/web/components/chat"
	componentsSudoku "github.com/deej-tsn/compSudoku/web/components/sudokuBoard"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketHandler struct {
	Game    *models.Game
	ChatLog *models.ChatLog
}

type WebSocketMessage struct {
	Type   string                 `json:"type"`
	Data   json.RawMessage        `json:"data"`
	Header map[string]interface{} `json:"HEADERS"`
}

type WebSocketConnection struct {
	Conn     *websocket.Conn
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
		Game:    game,
		ChatLog: chatLog,
	}
}

func (wsHandler *WebSocketHandler) InitWs(c echo.Context) error {
	conn, err := webSocketUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	currentConn := &WebSocketConnection{
		Conn:     conn,
		SignedIn: false,
	}

	webSocketConnections[currentConn] = true
	defer func() {
		delete(webSocketConnections, currentConn)
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
			fmt.Println(wsMsg)
			data, err := wsMsg.Data.MarshalJSON()
			if err != nil {
				c.Logger().Error(err)
				continue
			}
			fmt.Println(string(wsMsg.Data))
			switch wsMsg.Type {
			case "chat_message":

				wsHandler.handleChatMessage(c, currentConn, data)

			case "new_board":
				wsHandler.handleNewBoard(c, conn, data)

			case "set_active_square":
				if wsHandler.Game.State == models.GAME_STATE_IN_PROGRESS {
					wsHandler.handleSetActiveSquare(c, conn, data)
				}

			case "set_square_value":
				if wsHandler.Game.State == models.GAME_STATE_IN_PROGRESS {
					wsHandler.handleSetSquareValue(c, conn, data)
				}

			case "flip_edit_mode":
				wsHandler.handleFlipEditMode(c, conn)
			}
		} else {
			wsHandler.handleSignIn(c, conn, currentConn, msg)
		}

	}

	return nil
}

func (wsHandle *WebSocketHandler) handleSignIn(c echo.Context, conn *websocket.Conn, currentConn *WebSocketConnection, data []byte) {
	var SignedInReq models.JSONSignIn
	err := json.Unmarshal(data, &SignedInReq)
	if err != nil {
		c.Logger().Error(err)
		return
	}
	currentConn.Username = SignedInReq.Username
	currentConn.SignedIn = true
	currentConn.Color = SignedInReq.Color
	wsHandle.UserConnected(c, *currentConn)
}

func (wsHandler *WebSocketHandler) handleChatMessage(c echo.Context, currentConn *WebSocketConnection, data []byte) {
	var messageReq models.JSONMessage
	if err := json.Unmarshal([]byte(data), &messageReq); err != nil {
		c.Logger().Error(err)
		return
	}

	if messageReq.Text == "" {
		return
	}

	message := *models.NewMessage(messageReq.Text, currentConn.Username, currentConn.Color)
	wsHandler.ChatLog.Messages = append(wsHandler.ChatLog.Messages, &message)
	c.Response().Header().Set("HX-Trigger", "newMessage")
	bytes := wsHandler.renderChatMessage(c, message, true)
	if err := currentConn.Conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
		c.Logger().Error(err)
		return
	}
	bytes = wsHandler.renderChatMessage(c, message, false)
	broadcastToOthers(c, currentConn, bytes)

}

func (wsHandler *WebSocketHandler) UserConnected(c echo.Context, user WebSocketConnection) {
	systemMessage := fmt.Sprintf("%s has connected", user.Username)
	message := *models.NewMessage(
		systemMessage,
		"system",
		"purple",
	)

	bytes := wsHandler.renderChatMessage(c, message, false)
	go broadcastToAll(c, bytes)
}

func (wsHandler *WebSocketHandler) handleNewBoard(c echo.Context, conn *websocket.Conn, data []byte) {
	var difficultyJson models.DifficultyJSON
	if err := json.Unmarshal(data, &difficultyJson); err != nil {
		c.Logger().Error(err)
		return
	}
	newGrid, err := helper.GetBoardAPI(difficultyJson.Difficulty)
	if err != nil {
		conn.WriteJSON(map[string]string{"error": "Failed to fetch new board"})
		return
	}
	wsHandler.Game = models.ResponseToGame(newGrid)
	gameBuffer := new(bytes.Buffer)
	err = componentsSudoku.Game(wsHandler.Game).Render(context.Background(), gameBuffer)
	if err != nil {
		c.Logger().Error(err)
		return
	}
	broadcastToAll(c, gameBuffer.Bytes())
}

func (wsHandler *WebSocketHandler) handleSetActiveSquare(c echo.Context, conn *websocket.Conn, data []byte) {

	var dataString models.ActiveSquareJSON
	if err := json.Unmarshal([]byte(data), &dataString); err != nil {
		c.Logger().Error(err)
		return
	}

	position := helper.StringToPosition(dataString.Position)
	wsHandler.Game = wsHandler.Game.SetActiveSquare(wsHandler.Game.Grid[position[0]][position[1]])
	wsHandler.broadcastSudokuState(c)
}

func (wsHandler *WebSocketHandler) handleSetSquareValue(c echo.Context, conn *websocket.Conn, data []byte) {
	var dataString models.SquareValueJSON
	if err := json.Unmarshal([]byte(data), &dataString); err != nil {
		c.Logger().Error(err)
		return
	}
	value, err := strconv.Atoi(dataString.Value)
	if err != nil {
		conn.WriteJSON(map[string]string{"error": "Invalid value"})
		return
	}
	wsHandler.Game = wsHandler.Game.SetActiveSquareValue(value)
	if wsHandler.Game.State == models.GAME_STATE_COMPLETE {
		wsHandler.sendCompleteGame(c)
	} else if wsHandler.Game.State == models.GAME_STATE_FAILED {
		wsHandler.sendFailedGame(c)
	} else {
		for number := range wsHandler.Game.ToDisable {
			disabledComponent := wsHandler.renderDisabledNumber(c, number+1)
			broadcastToAll(c, disabledComponent)
			delete(wsHandler.Game.ToDisable, number)
		}
		wsHandler.broadcastSudokuState(c)
	}

}

func (wsHandler *WebSocketHandler) handleFlipEditMode(c echo.Context, conn *websocket.Conn) {
	wsHandler.Game.EditState = !wsHandler.Game.EditState
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

func (wsHandler *WebSocketHandler) sendFailedGame(c echo.Context) {
	failedGame := new(bytes.Buffer)
	err := componentsSudoku.FailedGame().Render(context.Background(), failedGame)
	if err != nil {
		c.Logger().Error(err)
		return
	}
	broadcastToAll(c, failedGame.Bytes())
}

func (wsHandler *WebSocketHandler) sendCompleteGame(c echo.Context) {
	completeGame := new(bytes.Buffer)
	err := componentsSudoku.CompletedGame().Render(context.Background(), completeGame)
	if err != nil {
		c.Logger().Error(err)
		return
	}
	broadcastToAll(c, completeGame.Bytes())
}

// RENDERERS

func (wsHandler *WebSocketHandler) renderSudokuGrid(c echo.Context) []byte {
	gridBuffer := new(bytes.Buffer)
	err := componentsSudoku.Grid(wsHandler.Game).Render(context.Background(), gridBuffer)
	if err != nil {
		c.Logger().Error(err)
		return []byte{}
	}
	return gridBuffer.Bytes()
}

func (wsHandle *WebSocketHandler) renderDisabledNumber(c echo.Context, number int) []byte {
	disableNumberBuffer := new(bytes.Buffer)
	err := componentsSudoku.Number(number, true).Render(context.Background(), disableNumberBuffer)
	if err != nil {
		c.Logger().Error(err)
		return []byte{}
	}
	return disableNumberBuffer.Bytes()
}

// BROADCAST

func (wsHandler *WebSocketHandler) broadcastSudokuState(c echo.Context) {
	gridBytes := wsHandler.renderSudokuGrid(c)
	broadcastToAll(c, gridBytes)
}

func broadcastToAll(c echo.Context, message []byte) {
	for connection := range webSocketConnections {
		if err := connection.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			delete(webSocketConnections, connection)
			connection.Conn.Close()
			c.Logger().Error(err)
		}
	}
}

func broadcastToOthers(c echo.Context, currentConn *WebSocketConnection, message []byte) {
	for connection := range webSocketConnections {
		if connection != currentConn {
			if err := connection.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				delete(webSocketConnections, connection)
				connection.Conn.Close()
				c.Logger().Error(err)
			}
		}
	}
}
