# Audit-A 自动化回归审计报告

- 日期：2026-09-19（UTC 时间盒 08:50 硬停；报告按到点已有证据撰写）
- 审计对象：13 张已完成票（01–09、13、14、17、18，`-done` 后缀）
- 方法：纯命令回归 + 全仓扫描，未读长文档。环境：go1.27.1 windows/amd64，branch `dev`，HEAD `d101c7f`。

## 1. 全仓回归（命令结果）

| 检查 | 结果 |
|---|---|
| `go vet ./...` | PASS（exit 0） |
| `go build ./...`（CGO_ENABLED=1） | PASS（exit 0） |
| `go test ./... -count=1`（CGO_ENABLED=1） | PASS（exit 0；20 个包 ok，audio 15.7s / openaichat 17.2s / memory 12.2s 最慢） |
| `go test ./... -count=1`（CGO_ENABLED=0） | **FAIL（exit 1）**：仅 `cmd/wisp` setup 失败——`main.go:20` / `doctor.go:14` 直接 import `sherpa-onnx-go/sherpa_onnx`（经 sherpa-onnx-go-windows，无 CGO 时 build constraints exclude all Go files）；internal/ 全部 ok |
| `go test -race ./internal/memory/... ./internal/secret/... ./internal/risk/...` | PASS（exit 0） |

## 2. Stub / 假完成扫描

`grep -rn "todo!|unimplemented|panic(\"not implemented\")" internal/ cmd/ tools/ --include="*.go"`（排除 `_test.go`）：**0 命中**。无可抽查上下文（要求抽查 3 个，实际无可查）。

## 3. D22 扫描器与并发纪律

- `tools/d22scan go run .`：`d22scan: clean - no D22 ban violations, no emoji in design/ or frontend/`（exit 0）。
- 独立 grep emoji 码点（Go 文件）：2 处命中均为中文注释中的箭头 `→`（U+2192，非 emoji）：`internal/ball/position.go:11`、`internal/memory/schema.go:118`。与 d22scan clean 结论一致，仅记录为 INFO。
- `grep -rn "go func(" internal/ cmd/ tools/`（排除测试）：实际代码 0 处裸 goroutine；仅 `tools/d22scan/main.go:8,239` 注释文字提及该模式。符合"只允许 observe.Spawn"纪律。

## 4. gofmt 漂移

`gofmt -l internal cmd tools`：**无输出（0 漂移）**。

## 5. 票据簿记一致性

- 13 张 `-done` 票：13/13 含 `(DONE` 标记；12/13 含 `**Status:** done`。
- **例外**：`02-s0-spike-done.md` 文件名已带 `-done`、内含 DONE 标记，但 `**Status:**` 仍为 `ready-for-agent`（其 Progress-log 明言 "Status left for orchestrator"，orchestrator 归档时改了文件名未改状态字段）。
- 索引 `.scratch/wisp/issues/README.md` 依赖图：13 张 done 票全部正确标注 `✅done`，无遗漏、无错标。
- 未完成票（10–12、15–16、19–61）：全部 `ready-for-agent`，无 in-progress 滞留，无"done 完成度但未归档"（非 done 文件中含 DONE 标记的仅 README.md 索引自身，属正常）。

## 6. 依赖白名单（go.mod 直接 require）

| 依赖 | 白名单判定 |
|---|---|
| `github.com/k2-fsa/sherpa-onnx-go v1.13.8` | 允许 |
| `github.com/pelletier/go-toml/v2 v2.2.4` | 允许（go-toml/pelletier） |
| `golang.org/x/sys v0.48.0` | 允许 |
| `modernc.org/sqlite v1.59.0` | 允许 |

- `golang.org/x/crypto v0.57.0` 仅 indirect（blake2b 用途，白名单允许）。
- 未出现 webrtc / goja / go-webview2 / goreleaser 直接依赖。**白名单 CLEAN**。

## 7. 密钥/敏感扫描

`grep` 高熵模式 + `sk-[A-Za-z0-9]{10,}` + `-----BEGIN`（覆盖 internal/ cmd/ tools/ build/ scripts/ models/ design/ docker/，排除 third_party/，排除测试）：**0 命中**（minisign 假钥 fixture 亦无 BEGIN 块命中）。

## 发现清单

### BLOCKER
（无）

### MAJOR
（无）

### MINOR
1. `.scratch/wisp/issues/02-s0-spike-done.md` — 票据状态字段不一致：文件名 `-done` + DONE 标记，但 `**Status:** ready-for-agent`，违反索引规则 4（完成须 Status: done + 重命名）。元数据问题，索引与交接记录均按 done 处理。
2. `cmd/wisp/main.go:20`、`cmd/wisp/doctor.go:14` — CGO_ENABLED=0 下整包无法编译（sherpa-onnx 无 CGO 被排除）。internal/ 在 CGO=0 下全绿（票 04 记录的 "CGO_ENABLED=0-ok" 仅覆盖 internal/memory，仍成立）；但若 SPEC-01 存在"全仓 CGO=0 构建目标"，此条须升级为 MAJOR——留给 orchestrator 对照 SPEC-01 裁决。

### INFO
- `internal/ball/position.go:11`、`internal/memory/schema.go:118` 中文注释含箭头 `→`（U+2192），非 emoji，d22scan 判定 clean，不计问题。

### 未覆盖项
- Stub 扫描上下文抽查（0 命中，无可抽查）。
- minisign 测试 fixture 假钥的内容核验（扫描无 `-----BEGIN` 命中，未逐文件人工复核 fixture 格式）。

## 结论

13 张票的自动化回归面：编译、vet、双轮测试、race、D22、gofmt、goroutine 纪律、依赖白名单、密钥扫描全部通过；簿记 12/13 完整；CGO=0 唯一红点集中于 cmd/wisp 的 sherpa 原生依赖，属已知 CGO 边界而非回归。

**AUDIT-A VERDICT: ISSUES: 2**
