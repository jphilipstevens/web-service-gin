# Logging with Trace IDs

- Use `logutils.WithTrace(ctx)` for all logs inside a request.
- Example:

```go
logutils.WithTrace(ctx).Info("Election started")
```

This ensures traceability across the distributed system.
