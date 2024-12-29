package main

import (
	"net/http"

	"github.com/deej-tsn/compSudoku/internal/helper"
	"github.com/deej-tsn/compSudoku/internal/models"
	layoutComponents "github.com/deej-tsn/compSudoku/web/components/layout"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	sudokuRoutes "github.com/deej-tsn/compSudoku/internal/routes"
)

func main() {

	pathToSudoku := "sudoku.txt"
	pathToWeb := "./web/public"
	game := models.NewGame(helper.ReadFileToString(pathToSudoku))
	chatLog := models.NewChatLog()
	e := echo.New()

	//CONTROLLERS

	sudokuController := sudokuRoutes.NewSudokuRoute(game)
	chatController := sudokuRoutes.NewChatHander(*chatLog)

	// MIDDLEWARE

	// logs all http requests
	e.Use(middleware.Logger())

	// static files in public folder
	e.Static("static", pathToWeb)

	// resets server if error
	e.Use(middleware.Recover())

	// CORS - cross origin resource sharing
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// allow any origin - DEVELOPMENT ONLY
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// ROUTES

	///HOME
	e.GET("/", func(c echo.Context) error {
		return helper.Render(c, http.StatusOK, layoutComponents.Index("Grid", game, chatLog))
	})

	//Sudoku
	e.POST("/sudoku", sudokuController.UpdateGrid)

	// CHATS
	e.GET("/chats", chatController.InitWs)
	e.POST("/messages", chatController.PostMessage)

	e.Logger.Fatal(e.Start(":8080"))

}
