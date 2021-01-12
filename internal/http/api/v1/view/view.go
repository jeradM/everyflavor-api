package view

type ListResult struct {
	Results interface{} `json:"results"`
	Count   uint64      `json:"count"`
}
