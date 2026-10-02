package ai

import (
	"errors"
	"sync"
)

func (s *Service) acquire(accountID string, phase Phase) (func(), error) {
	if s == nil || accountID == "" || phase != PhasePreview && phase != PhaseRun {
		return nil, errors.New("INVALID_AI_RUN")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.activeAccounts[accountID] {
		return nil, errors.New("AI_RUN_ALREADY_ACTIVE")
	}
	hour := s.now().Unix() / 3600
	rate := s.rateByAccount[accountID]
	if rate == nil || rate.hour != hour {
		rate = &accountRate{hour: hour}
		s.rateByAccount[accountID] = rate
	}
	if err := advanceAccountRate(rate, phase); err != nil {
		return nil, err
	}
	s.activeAccounts[accountID] = true
	var once sync.Once
	return func() { once.Do(func() { s.mu.Lock(); delete(s.activeAccounts, accountID); s.mu.Unlock() }) }, nil
}

func advanceAccountRate(rate *accountRate, phase Phase) error {
	if phase == PhasePreview {
		if rate.previews >= 10 {
			return errors.New("AI_PREVIEW_RATE_LIMIT")
		}
		rate.previews++
		return nil
	}
	if rate.runs >= 10 {
		return errors.New("AI_RUN_RATE_LIMIT")
	}
	rate.runs++
	return nil
}
