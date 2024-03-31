package model

// Job currency
type Job struct {
	ID       string `json:"id"`
	Stage    string `json:"stage"`
	CharCode string `json:"charCode"`
	Date     int64  `json:"date"`
}

type JobRepository interface {
	Find(filter JobFilter) ([]Job, error)
	JobFilter(Job) error
}

type JobFilter struct {
	Date     int64
	CharCode string
}
