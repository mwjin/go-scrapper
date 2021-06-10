package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type job struct {
	id       string
	location string
	title    string
	salary   string
	summary  string
}

var baseUrl string = "https://kr.indeed.com/jobs?q=python"

func main() {
	totalPages := GetPageCnt()
	fmt.Println(totalPages)

	for i := 0; i < totalPages; i++ {
		getPage(i)
	}
}

func getPage(page int) {
	pageUrl := baseUrl + "&start=" + strconv.Itoa(page*10)
	fmt.Println("Requesting " + pageUrl)
	res, err := http.Get(pageUrl)
	checkErr(err)
	checkStatus(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, card *goquery.Selection) {
		id, _ := card.Attr("data-jk")
		title := card.Find(".title>a").Text()
		location := card.Find(".sjcl").Text()
		fmt.Println(id, title, location)
	})
}

// How many pages are there?
func GetPageCnt() int {
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
