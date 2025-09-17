package scraper

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

// ScrapeBills returns a channel of Bills as they are found
func ScrapeBills(startURL string, out chan<- Bill)  {

	c := colly.NewCollector(
		colly.Async(true),
		colly.AllowedDomains("opinion.lawmaking.go.kr"),
	)

	// Follow detail pages
	c.OnHTML("a[name=outDetail]", func(e *colly.HTMLElement) {
		re := regexp.MustCompile(`outDetailR\((\d+)\)`)
		if match := re.FindStringSubmatch(e.Attr("onclick")); len(match) > 1 {
			id := match[1]
			detailURL := "https://opinion.lawmaking.go.kr/gcom/nsmLmSts/out/" + id + "/detailRP"
			_ = e.Request.Visit(detailURL)
		}
	})

	// Extract bill info
	c.OnHTML("tr", func(e *colly.HTMLElement) {
		thText := strings.TrimSpace(e.ChildText("th"))
		var bill Bill

		switch {
		case strings.Contains(thText, "의안명"):
			bill.Name = strings.TrimSpace(e.ChildText("td span"))
		case strings.Contains(thText, "발의정보"):
			bill.Proposers = strings.TrimSpace(e.ChildText("td"))
		case strings.Contains(thText, "주요내용"):
			bill.MainText = strings.TrimSpace(e.ChildText("td"))
		}

		if bill.Name != "" || bill.Proposers != "" || bill.MainText != "" {
			out <- bill
		}
	})

	go func() {
		defer close(out)
		c.Visit(startURL)
		c.Wait()
	}()
}
