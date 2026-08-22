# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
ok  	ruralfolk/api	0.007s
ok  	ruralfolk/cmd/server	0.003s
ok  	ruralfolk/config	0.001s
ok  	ruralfolk/domain	0.001s
ok  	ruralfolk/media	0.001s
--- FAIL: TestBusinessChain50 (0.00s)
    business_chain_test.go:28: publication returned mismatched record: domain.Exhibit{ID:"N-24", Title:"50 results", Story:"一段完整的村史", MediaURL:"", PublishedAt:"", Status:"published"}
FAIL
FAIL	ruralfolk/service	0.005s
ok  	ruralfolk/store	0.017s
ok  	ruralfolk/web	0.002s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/server): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/server): exit `0`
- Frontend build (web): exit `0`
