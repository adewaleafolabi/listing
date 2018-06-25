package model

type Property struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Sold   bool    `json:"sold"`
	Body   string  `json:"body"`
	Amount float64 `json:"amount"`
}
