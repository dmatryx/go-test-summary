package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/dmatryx/go-test-summary/internal/renderer"
	"github.com/dmatryx/go-test-summary/internal/results"
)

func main() {
	// Get hideUntestedPackages from the provided input (GHA converts this to an environment variable)
	hideUntestedPackagesValue := os.Getenv("INPUT_HIDEUNTESTEDPACKAGES")
	hideUntestedPackages, err := strconv.ParseBool(hideUntestedPackagesValue)
	if err != nil {
		// Log a warning and make the value false.
		println("::warning ::Invalid format for hideUntestedPackages input value")
		hideUntestedPackages = false
	}

	testDirectories := strings.Split(os.Getenv("INPUT_TESTDIRECTORIES"), "\n")

	testResults, exitCode := results.GetTestResults(testDirectories)
	// Establish the renderer with test results, and have it output
	markdownRenderer := renderer.Renderer{testResults, hideUntestedPackages}
	markdownRenderer.Render()

	os.Exit(exitCode)
}
