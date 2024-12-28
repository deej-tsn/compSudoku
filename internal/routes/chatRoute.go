package routes

import (
	"net/http"

	"github.com/deej-tsn/compSudoku/internal/helper"
	"github.com/deej-tsn/compSudoku/internal/models"
	components "github.com/deej-tsn/compSudoku/web/components/chat"
	"github.com/labstack/echo/v4"
)

type (
	ChatHandler struct {
		ChatLog *models.ChatLog
	}
)

func NewChatHander(chat models.ChatLog) *ChatHandler {
	ch := ChatHandler{
		ChatLog: &chat,
	}
	return &ch
}

func (chatH ChatHandler) GetMessages(c echo.Context) error {
	return helper.Render(c, http.StatusAccepted, components.Chatbox(*chatH.ChatLog))
}
