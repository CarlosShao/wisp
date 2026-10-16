票 270 归一微腿 `270-r2` 放行前更正（2026-10-16 12:4x 编排者追加；原文一字未改，⛔ 拿原文的行号/台面清单当现行）
================================================================================================

**释放条件已成立**：12:4x 实测 CPU 连发三发＝**44 / 42 / 43（%）**，机主的尺是"只要 cpu 不高于 70% 就可以继续补位"；同批进程名册＝台上只有 `wisp.exe` PID 21892 一枚、`go.exe`/`*.test.exe`/`gofumpt*` 命中 0（件 `.scratch/wisp/probes/orch/20261016/txt/roster-1016.ps1`，`-File` 跑）。⇒ 本腿**现在可放**，放出去就当他从头读一遍这枚派单＋下面三段更正。

## ① 台面清单改**八件套**（原文第 43 行的"七件套"作废，⛔ 照它建）

`git archive <具名 sha> go.mod go.sum internal cmd tools scripts frontend **models**`
少 `models/` ⇒ `internal/models` 会报三枚**假红** `committed models/manifest.json not from the test working directory`（成因：`internal/models/manifest.go:317 ResolveManifestPath` 运行时读 `models/manifest.json`，那两枚件是跟踪的）。另需手工拷进导出树（非跟踪或只跟踪一枚）：`build/*.dll`、`third_party/sherpa-onnx/*.dll`、`frontend/dist`，否则会拿回 `0xc0000135`／"no native DLLs"／白屏。

## ② ★门禁三数⛔ 在那棵**部分树**上跑（原文第 44 行没说，本条补上）

`tools/d22scan` 的 ban #8 射程**声明了 `design/` 与 `frontend/`**（`tools/d22scan/main.go:581-582`），而一枚声明了的射程扫到 0 枚文件是**硬退出**（同文件 `emptyEmojiScope`，`exit 2`）⇒ 在八件套部分树上跑会拿回 `scope design/ examined 0 files` 这种**GUARD 射程拒绝**，长得像红、⛔ 计入判据。
⇒ 门禁四数（`gofumpt`／`d22scan`／长度门／`go vet`）**只在工作树上跑**（本票写面只有 `internal/tools` 注释，工作树跑得起）；长度门要 `.git`，部分树根本没有。〔读码＋本波 `165-r3` 实测各报过一次同类缺件：它缺 `design/`、`220-v1` 缺 `models/`〕

## ③ 在飞名册换人（原文第 50 行那批已过时）

**现行**：`326-v1`（裁决腿，`internal/memory`／`internal/session`）∥ 本波同时放出的 `327-a1`／`328-a1`（只读普查）∥ `165-v2`（裁决 `internal/agent` 那枚新用例）∥ `216-a2`（只读普查）。
⇒ `internal/tools/subagent_197.go` 这枚面**没人抢**（`internal/tools` 本波零写腿）；`go test ./internal/tools/` 撞上了等一等再跑并具名写你等了。
