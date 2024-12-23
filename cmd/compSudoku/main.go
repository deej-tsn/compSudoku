package main

import (
	"net/http"

	"github.com/deej-tsn/compSudoku/internal/helper"
	layoutComponents "github.com/deej-tsn/compSudoku/web/components/layout"
	sudokuComponents "github.com/deej-tsn/compSudoku/web/components/sudokuBoard"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	pathToSudoku := "sudoku.txt"
	pathToWeb := "./web/public"
	grid := helper.NewGrid(helper.ReadFileToString(pathToSudoku))
	e := echo.New()

	//CONTROLLERS

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
		return helper.Render(c, http.StatusOK, layoutComponents.Index("Grid", sudokuComponents.Grid(grid)))
	})

	e.Logger.Fatal(e.Start(":8080"))

}
