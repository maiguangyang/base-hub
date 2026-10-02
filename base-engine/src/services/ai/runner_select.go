package ai

import (
	"google.golang.org/adk/v2/tool"
)

func (s *Service) toolsForRun(state *runState) ([]tool.Tool, error) {
	s.mu.Lock()
	registered := append([]tool.Tool(nil), s.tools...)
	s.mu.Unlock()
	selected := make([]tool.Tool, 0, len(registered)+2)
	selector, err := state.selectionTool()
	if err != nil {
		return nil, err
	}
	selected = append(selected, selector)
	for _, item := range registered {
		spec, exists := s.config.Catalog.Lookup(item.Name())
		if exists && specAllowed(spec, state.principal, state.phase, state.approvedIDs) {
			selected = append(selected, item)
		}
	}
	if state.phase == PhasePreview {
		proposal, err := state.proposalTool()
		if err != nil {
			return nil, err
		}
		selected = append(selected, proposal)
	}
	return selected, nil
}
