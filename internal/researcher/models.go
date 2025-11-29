package researcher

type ResearchDepth int

const (
	ResearchDepthShort  ResearchDepth = 0
	ResearchDepthMedium ResearchDepth = 1
	ResearchDepthLong   ResearchDepth = 2
)

func (d ResearchDepth) Validate() bool {
	switch d {
	case ResearchDepthShort, ResearchDepthMedium, ResearchDepthLong:
		return true
	default:
		return false
	}
}
