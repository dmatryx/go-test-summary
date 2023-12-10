# GitHub Action: `go-test-summary`

> Concept inspired heavily by `gotestsum` and GitHub blog on step summaries
> * https://github.com/gotestyourself/gotestsum
> * https://github.blog/2022-05-09-supercharging-github-actions-with-job-summaries/

The action is basically a wrapper for `go test` using the json formatted output from the test operation, then parsing
and rendering the output to provide a nice summary. Injecting that into the $GITHUB_STEP_SUMMARY file (which exists as
an environment variable in GitHub runners)

## Usage

### In GitHub actions

### If you want to run this locally for some reason

## Development

## Future ideas

* Provide the coverage output as a file in the repo's wiki
* Maybe re-use the feature from `go-coverage-report` below for coverage buttons


## Random other references
> * https://full-stack.blend.com/how-we-write-github-actions-in-go.html#fn:8
>   * This was an interesting blog, although used javascript shims to launch their go app.
> * https://docs.github.com/en/actions/creating-actions/creating-a-composite-action
>   * GitHub docs on creating composite actions
> * https://github.com/vakenbolt/go-test-report
>   * Renders go output as html
> * https://github.com/marketplace/actions/go-coverage-report
>   * Go coverage report as an action