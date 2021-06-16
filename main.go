package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/mwjjeong/go-scrapper/scrapper"
)

const fileName string = "jobs.csv"

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrapper)
	e.Logger.Fatal(e.Start(":8700"))
}

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handleScrapper(c echo.Context) error {
	defer os.Remove(fileName)
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	fmt.Println(term)
	scrapper.Scrape(term)
	return c.Attachment("jobs.csv", "jobs.csv")
}
