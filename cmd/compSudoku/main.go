package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/deej-tsn/compSudoku/internal/helper"
	"github.com/deej-tsn/compSudoku/internal/models"
	layoutComponents "github.com/deej-tsn/compSudoku/web/components/layout"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	routes "github.com/deej-tsn/compSudoku/internal/routes"
)

func main() {

	pathToWeb := "./web/public"
	godotenv.Load(".env")
	grid, err := helper.GetBoardAPI("2")
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

	wsController := routes.NewWebSocketHandler(game, chatLog)

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
		fmt.Println(game.State)
		return helper.Render(c, http.StatusOK, layoutComponents.Index("Sudoku", wsController.Game, wsController.ChatLog))
	})

	//Sudoku

	e.GET("/sudoku", wsController.InitWs)

	e.Logger.Fatal(e.Start(":8080"))

}
