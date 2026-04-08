package report

import (
	"testing"

	"github.com/ubnetdef/ScoreTrak/pkg/report"
)

func TestServiceEnabled(t *testing.T) {
	s := &report.SimpleService{Pause: false, Points: 123}
	if TotalServicePoints(s) != 123 {
		t.Fail()
	}
}

func TestServiceDisabled(t *testing.T) {
	s := &report.SimpleService{Pause: true, Points: 123}
	if TotalServicePoints(s) != 0 {
		t.Fail()
	}
}
