# 票 289 / 腿 289-r1 — 门禁（AC#4）与越界自证

> 每枚门禁件自己落一行 `rc=N`。零字节的件＝那一格没交。

- `GOFLAGS= go build ./...`：待填 rc=
- `gofmt -l`（4 枚件）：待填 rc=
- `gofumpt -l`（4 枚件）：待填 rc=
- `go test ./internal/panel/ ./cmd/wisp/ -count=1`（改后）：待填 rc=
- `sh scripts/d22scan.sh`（纯净树）：待填 rc=
- `git show --stat` 名册（应只含这 4 枚件）：待填
- 越界自证：`frontend/**`／`design/**`／三枚冻结件／golden／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt`／`docs/PLAN.md`／`docs/specs/**` 零字节：待填
