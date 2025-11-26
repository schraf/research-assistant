package models

import "fmt"

type ResourceMode int

const (
	ResourceModeMinimal ResourceMode = 0
	ResourceModeBasic   ResourceMode = 1
	ResourceModePro     ResourceMode = 2
)

// Validate ensures the ResourceMode is one of the defined values.
func (m ResourceMode) Validate() error {
	switch m {
	case ResourceModeMinimal, ResourceModeBasic, ResourceModePro:
		return nil
	default:
		return fmt.Errorf("invalid resource mode: %d", m)
	}
}

type ResearchDepth int

const (
	ResearchDepthShort  ResearchDepth = 0
	ResearchDepthMedium ResearchDepth = 1
	ResearchDepthLong   ResearchDepth = 2
)

// Validate ensures the ResearchDepth is one of the defined values.
func (d ResearchDepth) Validate() error {
	switch d {
	case ResearchDepthShort, ResearchDepthMedium, ResearchDepthLong:
		return nil
	default:
		return fmt.Errorf("invalid research depth: %d", d)
	}
}

type Schema map[string]any

// ResearchRequest represents the payload passed to the Cloud Run Job
type ResearchRequest struct {
	RequestId string        `json:"request_id"`
	Mode      ResourceMode  `json:"resource_mode"`
	Depth     ResearchDepth `json:"research_depth"`
	Topic     string        `json:"topic"`
}
