package events

import (
	"encoding/json"
	"strings"
	"time"
)

type TestEvent struct {
	Time    time.Time
	Action  string
	Package string
	Test    string
	Elapsed float64
	Output  string

	Subtest          bool
	PackageLevel     bool
	TestStatusResult bool
	Cached           bool
}

var TestStatusResults = []string{
	"pass",
	"fail",
	"skip"}

func ParseTestOutput(output string) ([]TestEvent, string) {
	var events []TestEvent
	var nonTestOutput []string
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		var event TestEvent
		if len(line) > 0 {
			err := json.Unmarshal([]byte(line), &event)
			if err != nil {
				nonTestOutput = append(nonTestOutput, line)
				println(line)
			}
			if event.Test == "" {
				event.PackageLevel = true
			}
			if strings.Contains(event.Output, "\t(cached)") {
				event.Cached = true
			}
			if strings.Contains(event.Test, "/") {
				event.Subtest = true
			}
			for _, result := range TestStatusResults {
				if strings.Contains(event.Action, result) {
					event.TestStatusResult = true
				}
			}
			events = append(events, event)
		}
	}
	return events, strings.Join(nonTestOutput, "\n")
}
