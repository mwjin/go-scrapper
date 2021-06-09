package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var baseUrl string = "https://kr.indeed.com/jobs?q=python&l="

func main() {
	totalPages := getPages()
	fmt.Println(totalPages)
}

// How many pages are there?
func getPages() int {
	pageCnt := 0
	res, err := http.Get(baseUrl)
	checkErr(err)
	checkStatus(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pageCnt = s.Find("a").Length()
	})
	return pageCnt
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkStatus(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with status: ", res.StatusCode)
	}
}
