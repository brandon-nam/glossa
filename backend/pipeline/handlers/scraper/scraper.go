package scraper

import (
	"log"
	"strconv"
	"strings"
	"web-scraper/backend/model"

	"github.com/gocolly/colly/v2"
)

type Scraper struct {
	StopAtId int
}

func (s Scraper) Run(out chan<- model.Bill) {
	log.Println("[Scraper] Starting Run")
	startURL := "https://opinion.lawmaking.go.kr/gcom/nsmLmSts/out"
	s.ScrapeBills(startURL, out)
	log.Println("[Scraper] Run finished")
}

// ScrapeBills returns a channel of Bills as they are found
func (s Scraper) ScrapeBills(startURL string, out chan<- model.Bill) {

	log.Println("[Scraper] Initializing Colly collector")
	detailCollector := colly.NewCollector(
		colly.Async(true),
		colly.AllowedDomains("opinion.lawmaking.go.kr"),
	)

	detailCollector.Limit(&colly.LimitRule{Parallelism: 3})

	// Channel for partial bills from list page
	listCh := make(chan *model.Bill, 100)

	go func() {
		for bill := range listCh {
			if bill.DetailURL != "" {
				// Creating context for this specific bill
				log.Println("[Scraper] Visiting detail url: ", bill.DetailURL)
				ctx := colly.NewContext()
				ctx.Put("bill", bill)
				detailCollector.Request("GET", bill.DetailURL, nil, ctx, nil)
			} else {
				// If no detail URL, send directly to out
				out <- *bill
			}
		}
		detailCollector.Wait()
		close(out) // signal that all full bills are done
	}()

	// Detail page callback
	detailCollector.OnHTML("tr", func(e *colly.HTMLElement) {
		bill := e.Request.Ctx.GetAny("bill").(*model.Bill) // can ensure that this bill is the one that needs to be populated

		// Populate additional fields from detail page
		bill.MainText = strings.TrimSpace(e.ChildText("td"))
	})

	detailCollector.OnScraped(func(r *colly.Response) {
		bill := r.Ctx.GetAny("bill").(*model.Bill)
		out <- *bill
	})

	// Start list page collector
	listCollector := colly.NewCollector(
		colly.Async(true),
		colly.AllowedDomains("opinion.lawmaking.go.kr"),
	)
	listCollector.Limit(&colly.LimitRule{Parallelism: 1})

	// Extract partial bill info and detail URL
	listCollector.OnHTML("tr", func(e *colly.HTMLElement) {
		var bill model.Bill

		// Extract '의안번호'(Bill ID)
		bill.BillId, _ = strconv.Atoi(strings.TrimSpace(e.ChildText("td[data-th='의안번호 (대안번호)']")))

		if bill.BillId <= s.StopAtId {
			log.Printf("[Scraper] Reached stop ID %d, stopping scraper.", s.StopAtId)
			return
		}

		// Extract '의안명' (Bill Name)
		bill.Name = strings.TrimSpace(e.ChildText("td[data-th='의안명'] a"))

		// Extract '제안자' (Proposer Info)
		bill.Proposers = strings.TrimSpace(e.ChildText("td[data-th='제안자(제안일자)']"))

		// Extract '상임위원회' (Department)
		bill.Department = strings.TrimSpace(e.ChildText("td[data-th='상임위원회(소관부처)]"))

		// Extract '국회현황(추진일자)'
		bill.ParliamentaryStatus = strings.TrimSpace(e.ChildText("td[data-th='국회현황(추진일자)']"))

		// Extract '의결현황(의결일자)'
		bill.ResolutionStatus = strings.TrimSpace(e.ChildText("td[data-th='의결현황(의결일자)']"))

		// Extract detail URL
		href := e.ChildAttr("a.mxW100", "href")
		if href != "" {
			bill.DetailURL = e.Request.AbsoluteURL(href)
		}

		// Send partial bill to channel
		if bill.Name != "" || bill.Proposers != "" || bill.MainText != "" {
			listCh <- &bill
		}
	})

	// Visit list page
	if err := listCollector.Visit(startURL); err != nil {
		log.Printf("[ListCollector] Error visiting %s: %v", startURL, err)
	}

	listCollector.Wait()
	close(listCh) // signal no more partial bills

	log.Println("[Scraper] Run finished")
}
