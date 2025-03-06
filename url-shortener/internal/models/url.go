package models

type URL struct {
	Id      int64  `json:"id"`
	URL     string `json:"url"`
	Address string `json:"address"`
}
