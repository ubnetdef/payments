package report

import "github.com/ubnetdef/ScoreTrak/pkg/report"

func TotalServicePoints(s *report.SimpleService) uint64 {
	if s.Pause {
		return 0
	}
	return s.Points
}
