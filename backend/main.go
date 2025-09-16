package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

type Bill struct {
    Name       string `json:"의안명"`       // 의안명
    Proposers  string `json:"발의정보"`     // 발의정보
    MainText   string `json:"주요내용"`     // 주요내용
}

func main() {
	// Create a new collector
	c := colly.NewCollector(
		colly.Async(true), // enable async requests
		colly.AllowedDomains("opinion.lawmaking.go.kr"), // stay in domain
	)

	// Set concurrency rules
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
	})


	// On every <a href="..."> found
	c.OnHTML("a[name=outDetail]", func(e *colly.HTMLElement) {
		onclick := e.Attr("onclick")

		re := regexp.MustCompile(`outDetailR\((\d+)\)`)
		match := re.FindStringSubmatch(onclick)

		if len(match) > 1 {
			id := match[1]
			detailURL := fmt.Sprintf("https://opinion.lawmaking.go.kr/gcom/nsmLmSts/out/%s/detailRP", id)
			fmt.Println("Following detail URL:", detailURL)

			err := e.Request.Visit(detailURL)
			if err != nil {
				log.Println("Visit error:", err)
			}
		}
	})

	// Print policy name, proposer info, and main content from detail page
	var bill Bill

	c.OnHTML("tr", func(e *colly.HTMLElement) {
		thText := strings.TrimSpace(e.ChildText("th"))

		switch {
		case strings.Contains(thText, "의안명"):
			bill.Name = strings.TrimSpace(e.ChildText("td span"))
		case strings.Contains(thText, "발의정보"):
			bill.Proposers = strings.TrimSpace(e.ChildText("td"))
		case strings.Contains(thText, "주요내용"):
			bill.MainText = strings.TrimSpace(e.ChildText("td"))
		}
	})


	// Debug: when any page is visited
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
	})

	// Start scraping from page 1
	startURL := "https://opinion.lawmaking.go.kr/gcom/nsmLmSts/out?pageIndex=1"
	c.Visit(startURL)

	// Wait until everything is done
	c.Wait()

	fmt.Println("Scraping finished!")
}