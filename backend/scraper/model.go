package scraper

// Bill represents the scraped lawmaking bill
type Bill struct {
	Name       string `json:"의안명"`
	Proposers  string `json:"발의정보"`
	MainText   string `json:"주요내용"`
	Summary    string `json:"요약,omitempty"`   // new field from GPT
	Categories string `json:"분류,omitempty"`   // example extra processing
}
