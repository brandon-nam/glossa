package handlers

import (
	"web-scraper/backend/model"
)

type Source interface {
	Run(out chan<- model.Bill)
}

type Transformer interface {
	Transform(in <-chan model.Bill, out chan<- model.Bill)
}

type Sink interface {
	Consume(in <-chan model.Bill)
}
