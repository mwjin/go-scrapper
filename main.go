package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	var allJobs []job
	jobsC := make(chan []job)
	totalPages := GetPageCnt()
	fmt.Println(totalPages)

	for i := 0; i < totalPages; i++ {
		go getPage(i, jobsC)
	}

	for i := 0; i < totalPages; i++ {
		jobs := <-jobsC
		allJobs = append(allJobs, jobs...) // append all the contents of jobs
	}

	writeJobs(allJobs)
	fmt.Printf("Extract %d jobs.\n", len(allJobs))
}

func writeJobs(jobs []job) {
	file, err := os.Create("jobs.csv")
	checkErr(err)
	w := csv.NewWriter(file)
	defer w.Flush()

	header := []string{"URL", "Title", "Location", "Salary", "Summary"}
	err = w.Write(header)
	checkErr(err)

	c := make(chan []string)
	for _, oneJob := range jobs {
		go makeJobSlice(oneJob, c)
	}

	var jobSlices [][]string

	for i := 0; i < len(jobs); i++ {
		jobSlice := <-c
		jobSlices = append(jobSlices, jobSlice)
	}

	if err := w.WriteAll(jobSlices); err != nil {
		checkErr(err)
		log.Fatalln("Error write a job to csv: ", err)
	}
}

func makeJobSlice(oneJob job, c chan<- []string) {
	jobBaseUrl := "https://kr.indeed.com/jobs?q=python&l&vjk="
	c <- []string{
		jobBaseUrl + oneJob.id,
		oneJob.title,
		oneJob.location,
		oneJob.salary,
		oneJob.summary}
}

func getPage(page int, jobsC chan<- []job) {
	var jobs []job
	jobC := make(chan job)
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
		go extractJob(card, jobC)
	})

	for i := 0; i < searchCards.Length(); i++ {
		oneJob := <-jobC
		jobs = append(jobs, oneJob)
	}
	jobsC <- jobs
}

func extractJob(card *goquery.Selection, c chan<- job) {
	id, _ := card.Attr("data-jk")
	title := cleanString(card.Find(".title>a").Text())
	location := cleanString(card.Find(".sjcl").Text())
	salary := cleanString(card.Find(".salaryText").Text())
	summary := cleanString(card.Find(".summary").Text())
	c <- job{id, title, location, salary, summary}
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

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
