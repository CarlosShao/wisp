# 票 29 落地腿 `29-r2` —— 229 全仓普查对 `internal/memory/l1_warm.go:221` 的唯一红（实现者）

> **HEAD**：`f4ad34ee`（现跑 `git rev-parse --short=8 HEAD`）。
> **写面**：`internal/tools/spawn_roster_census_229_test.go`（+5 行登记）、`internal/memory/l1_warm_roster_29r2_test.go`（新 1 枚）、本文件。
> **零字节**：`frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`cmd/wisp/**`、`internal/observe/**`、票面、台账（`git diff --name-only` 只有上一行那两枚 .go）。

## 0. 一句话

✅ 票 229 的全仓普查尺（`Test229SpawnNamesAreRosterShaped`）在 `f4ad34ee` 上**只剩一枚红**：`internal/memory/l1_warm.go:221` 的 `reg.Spawn(l1ExtractGoroutine, ...)` 首参是**包常量**（`l1ExtractGoroutine = "memory-extract"`），AST 判不了字面量，而 `spawn229ReviewedSites` 里没有这条复审记录。按该测试自己给出的补救路径补登记（运行期名字＋名册类别），**一个断言未动**；并给这条"散文式登记"补上**机器咬合**（`internal/memory` 侧新测试钉住常量值＋`ClassifyGoroutine` 判定），改名即红。改前红、改后绿，有反向钻。

## 1. 红在哪儿（RED-before，逐字）

`PATH=… go test ./internal/tools -run Test229SpawnNamesAreRosterShaped -count=1` ⇒ **FAIL**：

```
spawn_roster_census_229_test.go:181: internal/memory/l1_warm.go:221: Spawn 首参不是字符串字面量（l1ExtractGoroutine）且没有复审记录：
  把运行期名字与它落的名册类别写进 spawn229ReviewedSites
spawn_roster_census_229_test.go:192: 229 census 汇总：生产 .go 427 枚（解析跳过 0 枚）、Spawn 站点 25 个、根标签 17 个；名册外有登记 2 枚、未复审 1 枚
FAIL	github.com/CarlosShao/wisp/internal/tools	1.659s
```

红口**只有这一枚**（未复审 1 枚）；同一次扫描里 `名册外有登记 2 枚`＝`panel-sta` 族与 `mockllm-stdout-reader` 的历史债务，未动。

## 2. 改了什么（为什么是"登记"而不是"改名"）

- 先读现实：`internal/memory/l1_warm.go:55` `l1ExtractGoroutine = "memory-extract"`，`:56` owner＝`memory`；`internal/observe/goroutine.go:57-60` `TemporaryNames` 里逐字含 `memory-extract` ⇒ `ClassifyGoroutine("memory-extract") = CategoryTemporary`。**这枚名字本来就在 D38b 名册内**，红的原因是普查尺看不到常量值，不是名字跑出名册。
- 因此可动的只有两支：①把 `Spawn` 首参改成字面量（改生产代码，且让名字脱离单一常量源）；②按测试给的路补复审记录。**选②**：票 229 AC#2 的判据是"名字落在名册内"，改名不改分类、只增加重复字面量；登记则把"这枚站点为什么安全"写成人类可读的一行，且**不改任何表成员**（往名册加名＝D38b 契约改＝须人工批准，本腿⛔ 未加）。
- 落点（`internal/tools/spawn_roster_census_229_test.go`，map 内 +5 行）：

```go
// ticket 29 AC#1: Fire() hands the batch to the L1 extraction goroutine
// through a package const (l1ExtractGoroutine = "memory-extract"), so the
// census cannot see a literal here. The runtime name IS the frozen D38b
// temporary-roster name (observe.TemporaryNames), not an invented one.
`internal/memory/l1_warm.go|l1ExtractGoroutine`: "memory-extract：temporary（D38b 冻结临时族；L1 抽取只在 Warm 窗口内跑，永不上响应路径）",
```

- 登记键格式逐字照抄该 map 的既有约定：`<仓相对路径>|<源码里写的那段表达式>`（`expr229` 折叠空白 ⇒ 抗重排）；值＝运行期名字＋它落的名册类别。

## 3. 登记的"散文"补上机器咬合（`internal/memory/l1_warm_roster_29r2_test.go`）

散文登记跨包，常量被改名时**普查尺照样绿、登记语句变成谎话、每次 spawn 开始打"疑似泄漏"WARN**。新测试在**拥有该常量的包**里钉三件事：

| 断什么 | 反向钻（现跑） |
|---|---|
| `l1ExtractGoroutine == "memory-extract"` | 把常量改成 `l1-extract-drill` ⇒ `FAIL … l1_warm_roster_29r2_test.go:25: L1 extraction goroutine name = "l1-extract-drill"…` |
| `ClassifyGoroutine(l1ExtractGoroutine) == CategoryTemporary` | 同上（名字一出名册这一格也红） |
| 正控：名册外名字仍是 `CategoryUnknown`（分类器没被改成"什么都放行"） | 三格全绿时该行仍咬（`29r2-not-on-any-roster-table` ⇒ unknown） |

反向钻后按 `cp` 还原，`sha256(l1_warm.go)` 改前/改后同为 `3caea5fa2f5e83a7b3bbcfcb357ddb8997cabac2e5063f0b294e6b4f16cbb04d`（`RESTORED-IDENTICAL`），复跑绿。**生产文件净变动 0 字节。**

## 4. GREEN-after

```
go test ./internal/tools -run Test229SpawnNamesAreRosterShaped -count=1 -v
  → 229 census 汇总：生产 .go 427 枚（解析跳过 0 枚）、Spawn 站点 25 个、根标签 17 个；名册外有登记 2 枚、未复审 0 枚
  --- PASS: Test229SpawnNamesAreRosterShaped (0.49s)
  ok  	github.com/CarlosShao/wisp/internal/tools	0.559s

go test ./internal/memory -count=1                     → ok  github.com/CarlosShao/wisp/internal/memory  73.647s
go test -count=1 -timeout 200s ./internal/tools ./internal/memory
  → ok  internal/tools  27.822s ／ ok  internal/memory  36.846s
```

断言零放宽：`spawn229ReviewedSites` 只增一键，`spawn229OffRosterKnown`、`checked < 50`、`spawnSites < 10` 三道闸与错误文案逐字未动（`git diff` 只有 +5 行）。

## 5. 门禁（现跑读数）

| 闸 | 命令 | 读数 |
|---|---|---|
| build（本腿写面可达面） | `GOFLAGS= go build ./internal/... ./cmd/... ./tools/...` | **rc=0** |
| build（整仓） | `GOFLAGS= go build ./...` | **rc≠0，且非本腿引入**：`build/mut_rvw.go` 的未定义符号（`voiceTtsFiles`/`residentBall`/…）——该文件被 `.gitignore:14 build/` 忽略、mtime `2026-11-04 10:26:45`、91003 字节，是别的腿留下的变异评审副本；`go build` 不编译 `_test.go`，本腿 diff（两枚 `_test.go`）不可能造成它。**不删临时件**（规矩：只建不删），交回编排者 |
| 测试 | `go test -count=1 -timeout 200s ./internal/tools ./internal/memory` | **全绿**（见 §4） |
| 格式 | `gofumpt -l internal/tools internal/memory` | **空输出，rc=0** |
| D22 | `sh scripts/d22scan.sh` | **rc=0**，`d22scan: clean - no D22 ban violations` |

## 6. 残余（具名）

1. `go build ./...` 的 `build/mut_rvw.go` 红是**仓外/被忽略树里的既有脏件**（2026-11-04 10:26 留下），本腿写面外、未删（规矩禁删），需要该文件的作者或编排者处理。
2. 229 普查尺用**散文**记名册类别（本腿已补机器咬合，但止于 `memory-extract` 这一枚）；其余非字面量站点仍只有散文。是否把"复审记录→名册类别"做成结构化表（可机检）＝另一张票的事，未在本腿扩面。
3. 未 push；只 commit。
