package AI

import (
	"context"
	"fmt"
	"log"
	"time"
	"web-scraper/backend/model"

	"github.com/sashabaranov/go-openai"
)

type AIStage struct {
	client *openai.Client
}

// NewPipeline initializes GPT client
func NewAgent(apiKey string) *AIStage {
	return &AIStage{
		client: openai.NewClient(apiKey),
	}
}

func (a AIStage) Transform(in <-chan model.Bill, out chan<- model.Bill) {
	defer close(out)
	for bill := range in {
		// Enrich each bill with GPT
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

		prompt := fmt.Sprintf("다음 법안을 2줄로 요약하고 카테고리(예: 환경, 의료, 노동)를 지정해 주세요:\n\n의안명: %s\n발의정보: %s\n주요내용: %s",
			bill.Name, bill.Proposers, bill.MainText)

		resp, err := a.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{Role: "system", Content: "당신은 전문 입법 분석가입니다."},
				{Role: "user", Content: prompt},
			},
		})

		cancel()

		if err != nil {
			log.Printf("[AI] GPT error for bill %s: %v\n", bill.Name, err)
			out <- bill // send original bill if GPT fails
			continue
		}

		bill.Summary = resp.Choices[0].Message.Content
		out <- bill
	}
}
