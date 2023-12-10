# GitHub Action: `go-test-summary`

> Concept inspired heavily by `gotestsum` and GitHub blog on step summaries
> * https://github.com/gotestyourself/gotestsum
> * https://github.blog/2022-05-09-supercharging-github-actions-with-job-summaries/

This is a GitHub action to run `go test` using the json formatted output from the test operation, then parsing
and rendering the output to provide a nice summary. Injecting that into the $GITHUB_STEP_SUMMARY file (which exists as
an environment variable in GitHub runners).
If the test errors, the action will return the same exit code (and the output)

## Inputs

| Input | Default | Description |
| -     | -       | -           |
| hideUntestedPackages | `false` | Whether to minimise output for packages with no tests |

## Usage

### In GitHub actions
You should be able to hot-swap your existing `go test` step with usage of this GHA like so:
```diff
name: on-push-test-and-build-app

on: push

jobs:
  test-and-build-app:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        name: Checkout Repo

      - name: Setup Go 1.21.0
        uses: actions/setup-go@v4
        with:
          go-version: '^1.21.0'

-      - name: Run Go Test
-        run: go test ./...
+      - name: Go Test -> Summary
+        uses: dmatryx/go-test-summary@v1

      - name: Build
        run: CGO_ENABLED=0 GOOS=linux go build -o my-app-name .
```


### If you want to run this locally for some reason
Not really designed to be run locally, but you can grab and compile the binary itself.
The program only requires one environment variable, and that's `$GITHUB_STEP_SUMMARY` - This is where the markdown file
will be output to.

It's also worth noting that the Markdown generated is GitHub Flavour Markdown, which may or may not be supported by
whatever you are trying to view it through.

## Development

### Contributions
I'm all ears, but don't be a douche and declare into the void that the project is abandoned if you don't get a response
to your issue immediately.

### Future ideas

* Should move releases off to a feature branch which is really cut down
* Probably add some tests (oh the irony)
* Provide the coverage output as a file in the repo's wiki
* Maybe re-use the feature from `go-coverage-report` below for coverage buttons


### Random other references
> * https://full-stack.blend.com/how-we-write-github-actions-in-go.html#fn:8
>   * This was an interesting blog, although used javascript shims to launch their go app.
> * https://docs.github.com/en/actions/creating-actions/creating-a-composite-action
>   * GitHub docs on creating composite actions
> * https://github.com/vakenbolt/go-test-report
>   * Renders go output as html
> * https://github.com/marketplace/actions/go-coverage-report
>   * Go coverage report as an action