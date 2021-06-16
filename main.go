package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.Logger.Fatal(e.Start(":8700"))
}

func handleHome(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
