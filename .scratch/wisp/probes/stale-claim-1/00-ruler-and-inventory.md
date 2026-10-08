# stale-claim-1 / 00 — 尺、名册来源、存在性尺

- 腿名：`stale-claim-1`（纯只读普查）
- 取数时刻：`2026-10-08 15:04:03 +0800`（`date` 现量）；派单时刻 `14:48:17 +0800` 也记着
- ref：`dev` @ `6547fd30`（全号 `6547fd3046066d0f9e2fde9ba2cad3ee384573e8`，commit 时间 `2026-10-08T14:41:19+08:00`）
- **本程零写入既有文件**；只新建本目录下 5 枚 `.md`。**未跑任何 `go build`/`vet`/`test`/`go doc`/`go list`**；`go env` 也未用（没这个需要）。用到的取数只有 `git show HEAD:<path>` / `git grep <ref>` / `git ls-tree` / `wc -l` / `sed -n` / `awk` / `mkdir` / `ls` / `date` / `rev-parse`。
- 工作树里别人的在飞改动（`design/**` 的 16 删 4 改、`cmd/wisp`、被重写的 `probes/**`）**一概未读、未 add、未 checkout、未 restore、未 clean、未 stash**。取数一律走 HEAD 快照。

---

## 1. 这把尺长什么样（内容扫，不是符号名扫）

### 1.1 起手两串（按"这句话在说什么"扫）

```
git grep -nE "production caller|production callers|no caller|zero caller|not wired|unwired|had no assignment|never called|DORMANT|no production reader|no production writer|零调用者|生产零|未接线|未接入|尚未接入|没有接入|无人调用" HEAD -- ':(exclude).scratch' ':(exclude)*_test.go' '*.go'
```
中文那一族另跑过一遍宽串（含 `不存在／占位／没有读者／未使用／零引用`），英文另跑过 `no caller / is not wired yet / not connected / nothing calls / not reachable / no consumer`。

**枚数（同一条尺，四个射程，全部现量）**

| 射程 | 命中行数 | 说明 |
|---|---|---|
| 产码·非测试（`HEAD`，剥 `.scratch`、剥 `_test.go`） | **101 行** | 本程的主名册来源 |
| 产码·测试文件（剥 `.scratch`，只留 `_test.go`） | **87 行** | 单列，见 §1.4 |
| `.scratch/wisp/issues` + `.scratch/wisp/dispatches` + `docs/**` | **660 行** | 表二只取其中"派单会害到人"的那族，⛔ 不是逐枚裁决 |
| 全树（含 `.scratch/wisp/probes/**`） | 未取全量（输出 159 KB / 29.5 KB 两枚被截） | **具名声明**：probes 目录里是历次变异／pristine 快照，同一个产码注释会被复制成 3–9 枚，计入枚数就是把一把尺量出九枚 |

⇒ **回答"扫出多少枚"**：内容扫命中 **101 行**（产码非测试）。这 101 行不全是状态级断言：里头一大半是**运行时拒绝文案**（`任务名册未接线（fail-closed：…）`、`票 21 未接入`、`ResolveBytes is not wired`）与**标识符自身**（`unwiredKey`／`unwiredKeys`／`hotClaimNoReader`）。**逐枚读过之后判为"状态级断言"的是 39 枚**，全部进表一，逐枚带尺与分档（已过期 15／判不动 9／仍成立 15）。⛔ 我没把 101 行全裁完，剩下 62 行的分类＝运行时拒绝文案／fail-closed 行为分支／标识符与表名自指，其中**被本程显式剔出名册并具名的三组**见 `03` 丙节 R14。

⛔ 这一族**不能用符号名扫**：`RunProbeSuite`／`askConfirmation` 这种裸串会把定义行、注释行、测试名（`TestRunTextTask…`）、以及同名的 OS 类型全算成"调用者"。本程所有"零调用者"的判定一律走**调用形状**。

### 1.2 调用形状尺（判"死它"和"立它"同一把）

```
git grep -nE "<Receiver>\.<method>\(|<method>\(" HEAD -- ':(exclude).scratch' ':(exclude)*_test.go' '*.go'
```
命中里再**人工剥掉**：`func ` 开头那行＝定义、行首 `//`＝注释。剩下的才是调用点。同一枚判据正反两向都这么量。

实发例（表一里可复核）：
- `SealDir`：尺 `git grep -nE "SealDir\(|func SealDir" HEAD -- ':(exclude).scratch' ':(exclude)*_test.go' '*.go'` -> **1 命中，且它就是 `internal/winsec/winsec.go:222` 的定义行** ⇒ 生产调用点 **0**。
- `TaskRoster.Cancel`：尺 `git grep -nE "\.Cancel\(|func \(r \*TaskRoster\) Cancel" HEAD -- ...` -> 非测试里 `TaskRoster` 那族的调用点**只有** `internal/tools/task.go:740`（其余 15 枚 `.Cancel(` 是 `observe.Root`/`context` 根，同形不同符号，剥掉）。

### 1.3 三层结论（"有人调" 不等于 "生产到得了"）

表一每一枚都写满三层：
1. **包内可见** —— 非测试调用点在不在同一个包里；
2. **包外有引用** —— 别的产码包调不调；
3. **到达生产入口** —— `main()` → `cmd/wisp/main.go` 的 `case "run"`（`:91` -> `cmdRun` -> `runTextTask`）、`case "panel-inbound"`（`:115` -> `cmdPanelInbound`）、或常驻腿 `startResidentTaskSource`（`cmd/wisp/resident_windows.go:260`）。

第 3 层是本程最容易翻车的一层：`askConfirmation` 有生产调用点（`resident_approval_windows.go:700`），但那唯一一处住在 `AskOnTaskRoot` 的函数体里，而 `AskOnTaskRoot` 今天生产调用点为 **0**（U2 尺）⇒ 它**编译进产物、生产走不到**。表一按〔建了但没接〕这一档记。

### 1.4 测试文件那一族（单列，不混进产码名册）

87 行命中。本程只在两种时候引它们：**表三的正向断言的对照**，以及**"注释里写着我今天不成立"的那一枚**（`internal/tools/task_output_canonicalize_fail_174_test.go:21` 见表一·补）。⛔ 测试注释不是产码声明，但它会被派单当现状读，正是这一族坑的一半。

### 1.5 来源分档（判据没写"来源"这一栏就不给读数降档）

- **〔已证〕** = 尺与输出在案（本目录留了命令原文）
- **〔建了但没接〕** = 符号在、编译进产物、生产路径不达
- **〔仅注释／仅自述〕** = 只有注释或某腿的话在撑着

---

## 2. 存在性尺（`git show HEAD:<file> | wc -l`，全部现量）

```
247 internal/panel/composer_dispatch.go      38 internal/agent/approval/doc.go
527 internal/agent/approval/replies.go       637 cmd/wisp/approval_reply.go
362 cmd/wisp/config_readers_255.go           437 cmd/wisp/config_reload.go
185 cmd/wisp/firstrun.go                     310 cmd/wisp/panel_inbound.go
420 cmd/wisp/panel_pump.go                   988 cmd/wisp/resident_approval_windows.go
244 cmd/wisp/providers.go                   1427 cmd/wisp/run.go
195 cmd/wisp/approval_always.go              172 cmd/wisp/main.go
569 internal/panel/git.go                   288 internal/panel/subagent_roster_197.go
147 internal/config/unwired.go              279 internal/config/validate.go
353 internal/risk/assessor.go              1280 internal/tools/bridge.go
889 internal/tools/task.go                  389 internal/tools/paths.go
736 internal/winsec/resolve.go              388 internal/winsec/winsec.go
1042 internal/ball/ball_windows.go           211 internal/ball/d2d_windows.go
817 tools/d22scan/gitignore.go              497 internal/llm/anthropic/request.go
182 internal/panel/workspace.go             376 internal/agent/spill.go
1143 internal/agent/loop.go
```

⚠ 顺带一枚**指向错了的引用**：票面/台账里说的那串 `composer_dispatch.go:48`，在 HEAD 上**没有**那句断言——真正那一枚在 **`HEAD:internal/panel/composer_dispatch.go:50`**（`git show HEAD:internal/agent/composer/composer_dispatch.go` -> **空文件，`wc -l` = 0**，包路径不存在；文件在 `internal/panel/`）。行号已经漂了 2 行。⛔ 引这一族行号必须带 ref。

被 `composer_dispatch.go:51` 具名引用的裁决件 **存在**（尺：`git ls-tree -r --name-only HEAD docs/evidence/s1/ | grep '^docs/evidence/s1/33-'`）：`33-minimal-inbound-hop-r1.md` 在册，另有 `33-inbound-listener-r2/r3.md`、`33-inbound-hop-design-a1.md`、`33-panel-host-c27-r1…v2.md` 共 14 枚。

---

## 3. 文件导航

| 文件 | 内容 |
|---|---|
| `00-ruler-and-inventory.md` | 本文件：尺、枚数、存在性尺 |
| `01-prod-comment-roster.md` | **表一** 产码注释名册，30 枚逐尺复跑 |
| `02-ticket-and-ledger-roster.md` | **表二** 票面／派单／台账那一族 + "今天拿去派单会害到谁" |
| `03-reverse-roster.md` | **表三** 正向断言（已接线／已在生产／唯一调用者）名册 |
| `04-unjudgeable-and-blind.md` | **末节** 判不动的与看不见的（逐枚写缺哪一发读数）|
