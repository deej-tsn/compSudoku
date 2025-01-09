package main

import (
	"log"
	"net/http"
	"os"

	"github.com/deej-tsn/compSudoku/internal/helper"
	"github.com/deej-tsn/compSudoku/internal/models"
	layoutComponents "github.com/deej-tsn/compSudoku/web/components/layout"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	sudokuRoutes "github.com/deej-tsn/compSudoku/internal/routes"
)

func main() {

	pathToWeb := "./web/public"
	godotenv.Load(".env")
	grid, err := helper.GetBoardAPI()
	var game *models.Game
	if err != nil {
		log.Panicln("Cannot encode response to Object")
		game = nil
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
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = ":80"
	}
	e.Server.Addr = PORT

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
	e.POST("/sudoku/htmx/active", sudokuController.POSTActiveSquare)
	e.POST("/sudoku/htmx/active/change", sudokuController.POSTActiveSquareValue)
	e.POST("/sudoku/htmx/flipEditMode", sudokuController.POSTFlipEditMode)
	e.GET("/sudoku/htmx/board", sudokuController.GETBoard)
	e.GET("/sudoku/htmx/board/new", sudokuController.GETNewBoard)

	// CHATS
	e.GET("/chats", chatController.InitWs)
	e.POST("/messages", chatController.PostMessage)

	e.Logger.Fatal(e.Start(":8080"))

}
