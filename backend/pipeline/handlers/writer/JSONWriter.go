package writer

import (
	"encoding/json"
	"log"
	"net/http"

	"web-scraper/backend/pipeline/model"
)

type JSONSink struct {
	Writer http.ResponseWriter
}

func (s *JSONSink) Consume(in <-chan model.Bill) {
	w := s.Writer
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	w.Write([]byte("[")) // open JSON array

	first := true
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Println("Warning: ResponseWriter does not support flushing")
	}

	for bill := range in {
		if !first {
			w.Write([]byte(","))
		}
		first = false

		if err := enc.Encode(bill); err != nil {
			log.Println("encode error:", err)
			break
		}

		if ok {
			flusher.Flush() // stream chunk immediately
		}
	}

	w.Write([]byte("]")) // close JSON array
}
