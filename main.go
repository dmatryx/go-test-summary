package main

import (
	"os"
	"strconv"

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

	// Establish the renderer with test results, and have it output
	markdownRenderer := renderer.Renderer{results.GetTestResults(), hideUntestedPackages}
	markdownRenderer.Render()
}
