package main

import (
	"log"
	"net/http"

	"github.com/deej-tsn/compSudoku/internal/helper"
	"github.com/deej-tsn/compSudoku/internal/models"
	"github.com/deej-tsn/compSudoku/internal/routes"
	layoutComponents "github.com/deej-tsn/compSudoku/web/components/layout"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	sudokuRoutes "github.com/deej-tsn/compSudoku/internal/routes"
)

func main() {

	pathToSudoku := "sudoku.txt"
	pathToWeb := "./web/public"
	godotenv.Load(".env")
	grid, err := routes.GetBoardAPI()
	var game *models.Game
	if err != nil {
		log.Panicln("Cannot encode response to Object")
		game = models.NewGame(helper.ReadFileToString(pathToSudoku))
	} else {
		game = models.ResponseToGame(grid)
	}

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
		return helper.Render(c, http.StatusOK, layoutComponents.Index("Sudoku", game, chatLog))
	})

	//Sudoku
	e.POST("/sudoku", sudokuController.SetActiveSquare)
	e.POST("/sudoku/active", sudokuController.SetActiveSquareValue)
	e.POST("/sudoku/htmx/flipEditMode", sudokuController.PostFlipEditMode)
	e.GET("/sudoku/board", sudokuController.GetBoard)

	// CHATS
	e.GET("/chats", chatController.InitWs)
	e.POST("/messages", chatController.PostMessage)

	e.Logger.Fatal(e.Start(":8080"))

}
