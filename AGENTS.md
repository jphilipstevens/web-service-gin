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

## Commit message style

Craft commit messages using the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) format:

```
<type>(<scope>): <subject>

<optional body>
```

Use one of the following types to clearly convey the intent of the change:

- `feat`: a new feature
- `fix`: a bug fix
- `chore`: routine maintenance tasks
- `docs`: documentation changes
- `style`: formatting or stylistic updates that do not affect behavior
- `refactor`: code refactoring that neither fixes a bug nor adds a feature
- `test`: adding or updating tests
- `perf`: performance improvements

Keep the subject concise and place any additional details in the body. This convention ensures a clear history that is easy to scan and maintain.

