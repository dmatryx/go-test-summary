package renderer

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/dmatryx/go-test-summary/internal/results"
)

type Renderer struct {
	TestResults          results.TestingResults
	HideUntestedPackages bool
}

func (t *Renderer) testIcon(icon string) string {
	if icon == "pass" {
		return "🟢"
	} else if icon == "fail" {
		return "🔴"
	} else if icon == "skip" {
		return "🟡"
	} else if icon == "test" {
		return "🔬"
	} else if icon == "output" {
		return "🖨"
	} else if icon == "results" {
		return "📝"
	} else if icon == "package" {
		return "📦"
	} else if icon == "duration" {
		return "⏳"
	} else {
		return ""
	}
}

func (t *Renderer) getPackageDetails() string {
	tableBody := "<tr>" +
		"<th>" + t.testIcon("package") + " Package</th>" +
		"<th>" + t.testIcon("pass") + " Pass</th>" +
		"<th>" + t.testIcon("fail") + " Fail</th>" +
		"<th>" + t.testIcon("skip") + " Skip</th>" +
		"<th>" + t.testIcon("duration") + " Dur</th>" +
		"</tr>"
	for _, result := range t.TestResults.PackageResults {
		if result.HasTests() || !t.HideUntestedPackages {
			testSummaries := ""
			if result.HasTests() {
				testSummaries += "<details><summary>" + t.testIcon("test") + " Tests</summary><ul>"
				for s, testResult := range result.Tests {
					testSummaries += "<li>" + t.testIcon(testResult.TestStatusResult) + "<code>" + s + "</code></li>"
					if len(testResult.Subtests) > 0 {
						keys := []string{}
						for key := range testResult.Subtests {
							keys = append(keys, key)
						}
						sort.Strings(keys)
						testSummaries += "<ul>"
						for i := range keys {
							key := keys[i]
							testSummaries += "<li>" + t.testIcon(testResult.Subtests[key].TestStatusResult) + "<code>" + key + "</code></li>"
						}
						testSummaries += "</ul>"
					}
				}
				testSummaries += "</ul></details>"
			}
			testOutput := "\n\n```\n"
			outputLines := 0
			for _, event := range result.PackageLevelEvents {
				if len(event.Output) > 0 {
					testOutput += event.Output
					outputLines += 1
				}
			}
			for _, event := range result.Events {
				if len(event.Output) > 0 {
					testOutput += event.Output
					outputLines += 1
				}
			}
			if outputLines > 1 {
				testOutput = "<details><summary>" + t.testIcon("output") + " Output</summary>" + testOutput
			}
			testOutput += "```\n\n"
			if outputLines > 1 {
				testOutput = testOutput + "</details>"
			}
			tableBody += fmt.Sprintf(
				"<tr><td>%s<br><sub><i>%s</i></sub></td><td>%d</td><td>%d</td><td>%d</td><td>%.1fms</td></tr>",
				t.testIcon(result.PackageEvent.Action)+" <code>"+result.PackageEvent.Package+"</code>",
				result.Coverage,
				result.TestStatusResults["pass"],
				result.TestStatusResults["fail"],
				result.TestStatusResults["skip"],
				result.PackageEvent.Elapsed*1000)
			tableBody += "<tr><td colspan=\"5\">" + testSummaries + testOutput + "</td></tr>"
		}
	}
	return "<table>" + tableBody + "</table>\n"
}

func (t *Renderer) getSummaryText() string {
	total := 0
	statusTotals := make(map[string]int)
	for _, result := range t.TestResults.PackageResults {
		for statusType, i2 := range result.TestStatusResults {
			total += i2
			statusTotals[statusType] += i2
		}
	}
	summaryText := fmt.Sprintf("%d test", total)
	if total != 1 {
		summaryText += "s"
	}
	if total != 0 {
		var subSummaryParts []string
		subSummaryText := " ("
		for s, i := range statusTotals {
			subSummaryParts = append(subSummaryParts, fmt.Sprintf("%d %s", i, s))
		}
		subSummaryText += strings.Join(subSummaryParts, ", ") + ")"
		summaryText += subSummaryText
	}
	if t.HideUntestedPackages {
		skippedQuantity := 0
		for _, result := range t.TestResults.PackageResults {
			if !result.HasTests() {
				skippedQuantity += 1
			}
		}
		if skippedQuantity > 1 {
			summaryText += fmt.Sprintf(" -- %d packages with no tests.", skippedQuantity)
		}
	}
	summaryText += "\n"
	return summaryText
}

func (t *Renderer) Header(headerLevel int, headerText string) string {
	output := " " + headerText + "\n\n"
	for i := 0; i < headerLevel; i++ {
		output = "#" + output
	}
	return output
}

func (t *Renderer) Render() {
	output := t.Header(2, t.testIcon("results")+" Test summary")
	output += t.Header(3, "`"+t.TestResults.ModuleName+"`")
	output += t.getSummaryText()
	output += t.getPackageDetails()
	// Output a table of the package results
	outputFile := os.Getenv("GITHUB_STEP_SUMMARY")
	err := os.WriteFile(outputFile, []byte(output), 0666)
	if err != nil {
		log.Fatal(err)
	}
}
