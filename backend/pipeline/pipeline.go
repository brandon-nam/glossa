package pipeline

import (
	"context"
	"fmt"
	"log"

	"web-scraper/backend/scraper"

	openai "github.com/sashabaranov/go-openai"
)

// Pipeline with channels
type Pipeline struct {
	client *openai.Client
}

func (p *Pipeline) ProcessBills(param any, out chan scraper.Bill) any {
	panic("unimplemented")
}

// NewPipeline initializes GPT client
func NewPipeline(apiKey string) *Pipeline {
	return &Pipeline{
		client: openai.NewClient(apiKey),
	}
}

// RunPipeline starts processing: scrapeChan -> GPT -> resultChan
func (p *Pipeline) RunPipeline(scrapeChan <-chan scraper.Bill, resultChan chan<- scraper.Bill) {
	go func() {
		defer close(resultChan)
		for bill := range scrapeChan {
			// Enrich each bill with GPT
			ctx := context.Background()
			prompt := fmt.Sprintf("다음 법안을 2줄로 요약하고 카테고리(예: 환경, 의료, 노동)를 지정해 주세요:\n\n의안명: %s\n발의정보: %s\n주요내용: %s",
				bill.Name, bill.Proposers, bill.MainText)

			resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
				Model: "gpt-4o-mini",
				Messages: []openai.ChatCompletionMessage{
					{Role: "system", Content: "당신은 전문 입법 분석가입니다."},
					{Role: "user", Content: prompt},
				},
			})
			if err != nil {
				log.Printf("GPT error for bill %s: %v\n", bill.Name, err)
				resultChan <- bill // send original bill if GPT fails
				continue
			}

			bill.Summary = resp.Choices[0].Message.Content
			resultChan <- bill
		}
	}()
}
