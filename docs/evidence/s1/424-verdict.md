# 424 判词（424-review · 非实现者验收）

- 时刻：2026-11-12 11:2x +08
- 验收腿：`424-review`（只验收翻框，不改任何代码）
- 对象：票 424（Linux 编译断裂：`memory_retention_180k.go` 引用 Windows-only `residentVoice` ＋ 第二批 5 处同族 undefined）
- 结论：**7 / 7 AC 全部成立，票面 7 框已翻 `[x]`**

## 一、CI 真容器读数（本腿现跑，`gh api repos/CarlosShao/wisp/actions/jobs/130945374154`）

run `43146479546` / lint job `130945374154`，step 4–16 结论：

| step | 名称 | conclusion |
|---|---|---|
| 4 | D22 scanner positive control | success |
| 5 | D22 scanner self-test (ticket 161 AC#2) | success |
| 6 | D22 seven-ban + emoji scan (tools/d22scan) | success |
| 7 | Tracked path-length budget (ticket 262) | success |
| 8 | SLO nail census (ticket 266) | **failure**（既存红，非本修） |
| 9 | gofmt (gofumpt) | success |
| 10 | gofmt (gofumpt) tracked set (ticket 161 AC#7 form A) | success |
| **11** | **go vet (module)** | **success**（原 failure，判据核心） |
| 12 | go vet (tools/d22scan module) | success |
| 13 | staticcheck | **failure**（既存红 54 条，票 122，非本修） |
| 14 | mockllm module vet | success |
| 15 | Portable tests carrier self-test (ticket 111 AC#3) | success |
| 16 | Attribution instrument self-test (ticket 171 AC#3) | success |

step 11 = success ⇒ `go vet ./...` 在真 ubuntu 容器 rc=0 ⇒ **Linux 编译面恢复**（原错 `vet: cmd/wisp/memory_retention_180k.go:73:38: undefined: residentVoice` 不再出现）。

## 二、三把尺读数（本腿现跑，非抄派单）

**尺1 提交真在**（`git log --oneline -6` + `git show --stat`）：
- `e98276632` 424 第一批：`memory_retention_180k.go` −16 / 新 `memory_retention_180k_windows.go` +29 / `docs/evidence/s1/424-fix-linux-build.md` +113
- `3f15ede4d` 424-fix-b 第二批：新 `ball_config_appliers.go` +63、新 `config_err_text.go` +20、`ball_idle_appearance_windows.go` −31、`ball_style_65_windows.go` −16、`resident_audio_windows.go` −12、`docs/evidence/s1/424-fixb-linux-build.md` +59
- `aa23f834b` 424 第三批：新 `mode_request_114r2_helper_test.go` +17、`panel_mode_confirm_114_windows_test.go` −8、`docs/evidence/s1/424-fixc-r9-test-helper.md` +30

**尺2 件真在**（Test-Path 全 True，字节数）：`memory_retention_180k_windows.go`(1135)、`config_err_text.go`(695)、`ball_config_appliers.go`(2858)、`mode_request_114r2_helper_test.go`(902)、`424-fix-linux-build.md`(6263)、`424-fixb-linux-build.md`(3679)、`424-fixc-r9-test-helper.md`(2217)。

**尺3 自己复跑**：
- Windows：`go build ./cmd/wisp` **rc=0**；`go vet ./cmd/wisp` **rc=0**；`go test ./cmd/wisp -run 'Test180K|Test423|TestTicket255|Test65|Test255' -count=1` **rc=0**（`ok github.com/CarlosShao/wisp/cmd/wisp 27.010s`）
- Linux overlay type-check（GOOS=linux CGO_ENABLED=0，overlay 把 go.mod→sherpa 纯 Go 桩 `D:\tmp\wisp\20261112\424fix\sherpa-stub`）：
  - `go vet -overlay … ./cmd/wisp` **rc=0，`undefined:` 行数 = 0**
  - `go test -overlay -gcflags=all=-e -c`（全量错误清单，破 go vet 每包首错遮蔽）**rc=0，`undefined:` 行数 = 0**

## 三、逐 AC 判

| AC | 判 | 依据（本腿现跑） |
|---|---|---|
| #1 迁移函数 | ✅ | `memory_retention_180k_windows.go` 首行 `//go:build windows`、含 `attachMemoryL1IfEnabled180K`；原文件 grep `residentVoice\|attachMemoryL1IfEnabled180K` **0 命中**；3 个跨平台函数（`retentionWindowDays180K`:40 / `memoryL1Enabled180K`:56 / `retentionBootLine180K`:70）仍在原文件 |
| #2 Linux 编过（判据核心） | ✅ | CI step 11 `go vet (module)` = **success**；本腿自跑 linux vet rc=0 / test-c rc=0，undefined=0 |
| #3 Windows 不回归 | ✅ | win build rc=0、win vet rc=0、点名测试 rc=0（27.010s ok） |
| #4 越界 | ✅ | 三枚 commit `--stat` 全部落在票面既定文件＋对应 evidence 文档；`git status --porcelain` 无任何已跟踪文件被改动（全部 `??` 为历史遗留未跟踪临时件），本腿写面只有票 424 与本判词两件 |
| #5 门禁 | ✅ | CI step 9/10 gofumpt success；step 11 vet(module) success、step 12 vet(tools/d22scan) success |
| #6 第二批 5 处清零 | ✅ | `configErrText` 在跨平台 `config_err_text.go`（无 build tag）；3 个 `applyBall*` 包装 + 3 个 applier 变量声明在跨平台 `ball_config_appliers.go`；`installBallStyleApplier65` / `installBallIdleAppearanceAppliers` 留在 windows 文件（赋值处 `ball_style_65_windows.go:59`、`ball_idle_appearance_windows.go:65-66`）；linux 全量 undefined=0 |
| #7 两平台不回归 | ✅ | 尺3 双侧全 rc=0（见上） |

## 四、越界自证

`git status --porcelain`（验收动作前后各查一次）：输出全为 `??` 未跟踪项（历史遗留的探针/临时件），**无任何 ` M` 已跟踪改动**；本腿新增/修改仅两处显式 pathspec：`.scratch/wisp/issues/424-linux-build-break-memory-retention-180k-uses-windows-only-residentVoice.md`（翻框＋Progress）与本判词。零代码、零 ci.yml、零其它票改动。

## 五、附注

- step 8（SLO nail census 266）与 step 13（staticcheck 54 条，票 122）为**既存红**，与本修无关，不阻塞票 424 的 AC。
- 本腿未 push（铁律），commit 只含上述两处 pathspec。
