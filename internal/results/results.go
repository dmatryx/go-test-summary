package results

import (
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/dmatryx/go-test-summary/internal/events"
)

type TestingResults struct {
	ModuleName     string
	PackageResults []PackageResult
}

type PackageResult struct {
	PackageEvent       events.TestEvent
	PackageLevelEvents []events.TestEvent
	Events             []events.TestEvent
	Coverage           string
	Tests              TestDetails
	TestStatusResults  map[string]int
}

type TestDetails map[string]TestResult

type TestResult struct {
	TestStatusResult string
	Subtests         TestDetails
}

// getModuleName inspects the 'go.mod' file to establish what is being tested via regex match
func (t *TestingResults) getModuleName() {
	file, _ := os.ReadFile("./go.mod")
	re := regexp.MustCompile("module ([-a-zA-Z/.]+)")
	matches := re.FindStringSubmatch(string(file))
	t.ModuleName = matches[1]
}

// HasTests will return a boolean reflecting whether the package has tests
func (t *PackageResult) HasTests() bool {
	return t.testCount() > 0
}

// testCount totalizes all test results across the entire package
func (t *PackageResult) testCount() int {
	var total int
	for _, result := range events.TestStatusResults {
		total += t.TestStatusResults[result]
	}
	return total
}

// enumerateResults isn't a great name - this function takes the package's Events - It rearranges these into Tests with
// subtests to make traversal of the testing structure much easier to work with
func (t *PackageResult) enumerateResults() {
	t.TestStatusResults = make(map[string]int)
	t.Tests = make(map[string]TestResult)

	// Initialize test values to zero - this should probably happen as part of an initialisation thing....
	for _, result := range events.TestStatusResults {
		t.TestStatusResults[result] = 0
	}
	for _, event := range t.Events {
		if event.TestStatusResult {
			t.TestStatusResults[event.Action] += 1
		}
		if event.Subtest {
			parentTest, _, _ := strings.Cut(event.Test, "/")
			currentTest := t.Tests[parentTest]
			currentTest.Subtests[event.Test] = TestResult{TestStatusResult: event.Action}
			t.Tests[parentTest] = currentTest
		} else {
			currentTest := t.Tests[event.Test]
			currentTest.TestStatusResult = event.Action
			if currentTest.Subtests == nil {
				currentTest.Subtests = make(map[string]TestResult)
			}
			t.Tests[event.Test] = currentTest
		}
	}
}

// GetTestResults is currently the function which runs 'go test' and captures the output.  Right now it also then calls
// the parse of that output *and additionally* re-arranges that output into a more useful format.  The enumeration of
// results into a more useful format should really be put in a different function...
func GetTestResults() TestingResults {
	var testingResults TestingResults
	var packageResults []PackageResult
	coverageRegexp, _ := regexp.Compile("^coverage: (.+)\n$")

	output, _ := exec.Command("go", "test", "-json", "-count=1", "./...", "-cover").CombinedOutput()
	// TODO: Support arguments?
	// TODO: Move execution out to a runner, capture the exit code, exit this own process right at the end based on that exit code.

	// It may seem odd to not care about the error output here, but we are trying to capture the test output here - if we fatal error because
	// the process reported an exit code then we may not correctly interpret the test output.

	// if err != nil {
	//	log.Fatal(err)
	// }

	allEvents := events.ParseTestOutput(string(output))

	// TODO: Move this loop to a different function
	for _, event := range allEvents {
		if event.PackageLevel && event.TestStatusResult {
			var thisPackageEvents []events.TestEvent
			var thisPackageLevelEvents []events.TestEvent
			var coverage string
			for _, eventSubIterator := range allEvents {
				// Grab package events
				if eventSubIterator.Package == event.Package && !eventSubIterator.PackageLevel {
					thisPackageEvents = append(thisPackageEvents, eventSubIterator)
				} else if eventSubIterator.Package == event.Package && eventSubIterator.PackageLevel {
					thisPackageLevelEvents = append(thisPackageLevelEvents, eventSubIterator)
				}

				// Try to grab the coverage figure for the package
				if eventSubIterator.Package == event.Package && !eventSubIterator.TestStatusResult {
					if coverageRegexp.MatchString(eventSubIterator.Output) {
						matches := coverageRegexp.FindStringSubmatch(eventSubIterator.Output)
						coverage = matches[1]
					}
				}
			}
			packageResult := PackageResult{event, thisPackageLevelEvents, thisPackageEvents, coverage, nil, nil}
			packageResult.enumerateResults()
			packageResults = append(packageResults, packageResult)
		}

	}
	testingResults.getModuleName() // TODO: This should be called automatically when this object is instantiated.
	testingResults.PackageResults = packageResults
	return testingResults
}
