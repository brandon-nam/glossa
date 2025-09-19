package scraper

import (
	"log"
	"strings"
	"web-scraper/backend/pipeline/model"

	"github.com/gocolly/colly"
)

type Scraper struct{}

func (s Scraper) Run(out chan<- model.Bill) {
	log.Println("[Scraper] Starting Run")
	startURL := "https://opinion.lawmaking.go.kr/gcom/nsmLmSts/out"
	s.ScrapeBills(startURL, out)
	log.Println("[Scraper] Run finished")
}

// ScrapeBills returns a channel of Bills as they are found
func (s Scraper) ScrapeBills(startURL string, out chan<- model.Bill) {

	log.Println("[Scraper] Initializing Colly collector")
	c := colly.NewCollector(
		colly.Async(true),
		colly.AllowedDomains("opinion.lawmaking.go.kr"),
	)

	// Follow detail pages
	c.OnHTML("a.mxW100", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if href != "" {
			detailURL := e.Request.AbsoluteURL(href) // makes it absolute
			log.Printf("[Scraper] Visiting detail page: %s\n", detailURL)
			if err := e.Request.Visit(detailURL); err != nil {
				log.Printf("[Scraper] Error visiting detail page %s: %v\n", detailURL, err)
			}
		}
	})

	// Extract bill info
	c.OnHTML("tr", func(e *colly.HTMLElement) {
		thText := strings.TrimSpace(e.ChildText("th"))
		var bill model.Bill

		switch {
		case strings.Contains(thText, "의안명"):
			bill.Name = strings.TrimSpace(e.ChildText("td span"))
		case strings.Contains(thText, "발의정보"):
			bill.Proposers = strings.TrimSpace(e.ChildText("td"))
		case strings.Contains(thText, "주요내용"):
			bill.MainText = strings.TrimSpace(e.ChildText("td"))
		}

		if bill.Name != "" || bill.Proposers != "" || bill.MainText != "" {
			log.Printf("[Scraper] Extracted bill: %+v\n", bill)
			out <- bill
		}
	})

	log.Printf("[Scraper] Visiting start URL: %s", startURL)
	if err := c.Visit(startURL); err != nil {
		log.Printf("[Scraper] Error visiting start URL %s: %v", startURL, err)
	}
	c.Wait() // wait for all async requests to finish
	log.Println("[Scraper] Finished scraping")
	close(out)
}
