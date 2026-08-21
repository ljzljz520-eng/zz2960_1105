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
?   	inventoryseal/cmd/inventoryctl	[no test files]
ok  	inventoryseal/internal/archive	0.009s
ok  	inventoryseal/internal/config	0.013s
ok  	inventoryseal/internal/crypto	0.006s
ok  	inventoryseal/internal/domain	0.008s
ok  	inventoryseal/internal/httpapi	0.013s
ok  	inventoryseal/internal/importer	0.013s
ok  	inventoryseal/internal/ledger	0.003s
ok  	inventoryseal/internal/policy	0.004s
ok  	inventoryseal/internal/query	0.004s
ok  	inventoryseal/internal/reconcile	0.004s
ok  	inventoryseal/internal/report	0.004s
--- FAIL: TestBusiness11Regression (0.02s)
    workflow_test.go:112: expected cancellation, got <nil>
FAIL
FAIL	inventoryseal/internal/service	0.087s
ok  	inventoryseal/internal/store	0.012s
ok  	inventoryseal/internal/telemetry	0.003s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/inventoryctl): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/inventoryctl): exit `0`
