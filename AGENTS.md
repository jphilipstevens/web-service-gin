# Repository Guidelines

This project expects all contributions to be validated with the same checks that run in the CI workflow.

## Required checks

Run the following commands before committing any changes:

```bash
gofmt -s -w $(git ls-files '*.go')

go vet ./...

go test ./...

go build ./...
```

The `gofmt` step rewrites source files in place. Commit any updates after running these commands.

