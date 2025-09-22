package model

// Bill represents the scraped lawmaking bill
type Bill struct {
	Id                  int    `json:"id"`
	BillId              int    `json:"의안번호"`
	Name                string `json:"의안명"`
	Proposers           string `json:"발의정보"`
	Department          string `json:"상임위원회"`
	ParliamentaryStatus string `json:"국회현황"`
	ResolutionStatus    string `json:"의결현황"`
	MainText            string `json:"주요내용"`
	Summary             string `json:"요약,omitempty"` // new field from GPT
	Categories          string `json:"분류,omitempty"` // example extra processing
	DetailURL           string `json:"detail_url"`
}
