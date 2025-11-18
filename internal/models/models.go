package models

type Schema map[string]any

// ResearchRequest represents the payload passed to the Cloud Run Job
type ResearchRequest struct {
	RequestId string `json:"request_id"`
	Topic     string `json:"topic"`
}
