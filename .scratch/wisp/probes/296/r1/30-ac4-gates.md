# 票 296 AC#4 — 门禁读数（每把自落一行 rc=N）＋ 名册三数（带尺名）＋ 越界自查

leg = 296-r1 · 全部尺都在 20:39–21:14 这一段现跑 · 工作树＝`fc4aedaf`（`git diff --stat -- cmd/wisp/` 空）

## 1. 编译与仪器

```
GOFLAGS= go build ./...                                  rc=0     （最终态 21:1x 再跑一遍，同 rc=0）
sh scripts/d22scan.sh                                    rc=0     （d22scan: clean - no D22 ban violations）
```

d22scan 对本腿射程的读数（逐字尾两行，全文在 `/tmp/d22final.md`，另 `logs/build-final.md` 自落 rc 行）：

```
d22scan: scope ban #8 cmd/              examined 118 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=229, bans #1-5 cmd/=39, ... ban #8 cmd/=118
```

⇒ 新增那枚测试文件在 ban #8 的射程里（cmd/ 118 枚 Go 文件），零 emoji、零裸 `go func(`、零 `filepath.Clean|Abs`、零墙钟超时（夹具的驱动用的是既有 seam `rb.hotkeyBridgeCheck258()`，逐名 bounded loop，不睡墙钟判超时）。

## 2. 格式那两把（并排交，射程写清）

**尺①＝工作树那一把，射程 `cmd/wisp`**

```
$ gofmt -l cmd/wisp
cmd\wisp\models.go
cmd\wisp\panel_inbound_guards_35r3_test.go
cmd\wisp\panel_transport_35r2_test.go
rc=0        （3 枚）

$ "$(go env GOPATH)/bin/gofumpt.exe" -l cmd/wisp
cmd\wisp\models.go
cmd\wisp\panel_inbound_guards_35r3_test.go
cmd\wisp\panel_transport_35r2_test.go
rc=0        （3 枚，逐名同 gofmt 那把）
```

- `models.go` 那一枚＝**CRLF 假枚**（尺＝`tr -cd '\r' < cmd/wisp/models.go | wc -c` = **334**，而 `git show HEAD:cmd/wisp/models.go` 的 blob = **0**；`gofmt -l` 对 blob 拷贝＝空）⇒ checkout 换行造成的，不是谁写脏的。
- 另两枚＝票 298 名下那两枚已知脏件（本腿⛔ 没顺手洗：它们不在我的 pathspec 里，`git diff` 也零字节）。
- **本腿碰过的两枚文件单测两把尺都空**：
  ```
  gofmt -l cmd/wisp/resident_windows.go cmd/wisp/resident_hotkey_296_windows_test.go          rc=0（空）
  $(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp/resident_windows.go cmd/wisp/resident_hotkey_296_windows_test.go   rc=0（空）
  ```
- 具名报：**裸 `gofumpt` 不在 PATH** ⇒ `/usr/bin/bash: line 1: gofumpt: command not found`，`rc=127`。这不是"跳过"，是那一把尺必须用绝对路径。

**尺②＝判据「相对 HEAD 新增 0 枚」的改前／改后两把读数**（同一把尺：`git archive <commit> cmd/wisp | tar -x` 到仓外临时目录，再 `gofmt -l cmd/wisp`；两把之间只有 CRLF 那一枚差，所以逐名作差才作数）

| 读数点 | gofmt | gofumpt | 名册 |
|---|---|---|---|
| 改前 `28a2ff4e~1`（＝本腿锚点 `7285ef83` 的 `cmd/wisp`） | 2 枚 | 2 枚 | `panel_inbound_guards_35r3_test.go`／`panel_transport_35r2_test.go` |
| 改后 `HEAD=fc4aedaf` | 2 枚 | 2 枚 | 同上，逐名相同 |

⇒ **新增 0 枚**。工作树尺那一把（3 枚）改后没直接取"改前"，因为工作树已被本腿推进过；但可代数证成：本腿在 `cmd/wisp` 只碰两枚文件，而这两枚在改后名册里**都不在** ⇒ 名册增量只能来自这两枚 ⇒ 0。

## 3. 整包 `go test ./cmd/wisp/ -count=1`（改前 2 发／改后 2 发取交集）

尺（四发同一把）＝`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1 -v`
每发起手都先跑 `tasklist //FI "IMAGENAME eq wisp.exe"` 与 `... balldebug.exe` ⇒ **四发都 0 枚**（逐字 `INFO: No tasks are running which match the specified criteria.`），并 `date` 记锚。

| 发 | 起手 date | rc | `^--- FAIL`（顶层红） | `^--- PASS`（只数顶层） | `--- PASS`（含子测试） | `^--- SKIP`（顶层） | 钟点 |
|---|---|---|---|---|---|---|---|
| 改前 1（无 `-v`，`logs/prefix-full-1.md`） | 20:40:02 | 1 | **6** | 0（这把尺不产 PASS 行） | 0 | 未取 | 472.191s |
| 改前 2（`-v`，`logs/prefix-full-2.md`） | 20:49:11 | 1 | **6** | **273** | **383** | 2 | 474s |
| 改后 1（`-v`，`logs/postfix-full-1.md`） | 20:58:23 | 1 | **4** | **275** | **385** | 2 | 460s |
| 改后 2（`-v`，`logs/postfix-full-2.md`） | 21:06:16 | 1 | **4** | **275** | **385** | 2 | 466s |

⚠ 尺名口径照 10-09 新规矩写清：`--- PASS`（含子测试）与 `^--- PASS`（只数顶层）在两把读数里都各就各位（383↔273、385↔275），差 110 枚是本仓既有形状，不是名册在漂。
改前 1 那发**没带 `-v`** ⇒ 它天生零 PASS 行，只用来作红名册的交集，不许拿去和 `-v` 那三发的 PASS 数并排。

**逐名作差（改前 2 发交集 ∩ 改后 2 发交集）**：

- 改前红名册（6 枚，逐字）：
  ```
  --- FAIL: TestPanelHostRealWindowHopAndLifecycle (5.11s)
  --- FAIL: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (20.01s)
  --- FAIL: TestAC14AwaitedBindingReplyReachesThePage (20.01s)
  --- FAIL: TestAC14GoSideEvalPushReachesThePage (20.01s)
  --- FAIL: Test296BridgeKeepsThreeLiveKeysWithOwnerEmptySlots (0.04s)
  --- FAIL: Test296BridgeKeepsThreeLiveKeysWhenConfigUnreadable (0.04s)
  ```
  （改前 1 那发逐名相同，只钟点不同 ⇒ 两发改前交集＝这 6 枚）
- 改后红名册（4 枚，两发逐字相同 ⇒ 交集＝这 4 枚）：
  ```
  --- FAIL: TestPanelHostRealWindowHopAndLifecycle (5.05s)
  --- FAIL: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (20.02s)
  --- FAIL: TestAC14AwaitedBindingReplyReachesThePage (20.01s)
  --- FAIL: TestAC14GoSideEvalPushReachesThePage (20.00s)
  ```
- **新增红＝0 枚**；作差只有两个方向：本票那两枚夹具由红转绿（`--- PASS: Test296BridgeKeepsThreeLiveKeysWithOwnerEmptySlots (0.04s)`／`... WhenConfigUnreadable (0.04s)`，改后两发都在），顶层 PASS 由 273→275（含子 383→385）。
- 留下的 4 枚红**不属本票射程**（逐字因由，尺＝`grep -B 6 -- '^--- FAIL: <名>' logs/postfix-full-2.md`）：
  三枚 `panel_resident_windows_test.go` 的 AC#13/AC#14 全是 `no report "ac13-probe"/"ac14r-0"/"ac14-push" from the page within 15s (what DID arrive at the door: nothing at all)`，
  一枚 `panel_host_windows_test.go:662` 的 `cold bring-up did not produce a browser round trip (got -1.000)` ⇒ 都是 WebView2／页面投递那一族（`frontend/dist` 空那笔账的邻域），改前改后**都在**、与本腿那枚热键闭包无关。本腿⛔ 未动、未"顺手修"、也没为了让整包绿去放宽任何断言。

## 4. 越界自查（`git show --stat` 逐笔数过的文件枚数）

| commit | 内容 | 文件枚数 | 名册 |
|---|---|---|---|
| `ff5c193e` | 起手锚 | 2 | `probes/296/r1/00-anchor.md`／`probes/296/r1/logs/msg-01.md` |
| `28a2ff4e` | AC#1 夹具（两形，改前红）＋闭包提取 | 5 | `cmd/wisp/resident_hotkey_296_windows_test.go`／`cmd/wisp/resident_windows.go`／`probes/296/r1/10-ac1-red-before-fix.md`／`logs/red-before-fix.md`／`logs/msg-02.md` |
| `fc4aedaf` | AC#2 甲-全 ＋ 注释更正 | 4 | `cmd/wisp/resident_windows.go`／`probes/296/r1/20-ac2-fix-and-comment.md`／`logs/green-after-fix.md`／`logs/msg-03.md` |

尺＝`git diff --name-only 7285ef83..HEAD | grep -E "internal/|frontend/|design/|golden|thresholds\.go|allowlist\.txt|\.gitignore|config\.toml|tools/d22scan"` ⇒ **NONE（0 命中）**。
逐格确认零字节：`internal/config/schema.go`／`internal/ball/**`／三枚冻结件（`internal/panel/tokens_fourway_test.go`、`internal/panel/l2_grant_boundary_test.go`、`internal/perm/ticket90_persist_test.go`）／golden／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**`／`.gitignore`／D43 转移表／`config.toml`（机主那份未读未写）。
零文案改动：`muteGestureWhy296` 那一族注释体一字未动（裁定节归口"已裁可接受的注释体"）。
零 push；每笔 commit 都带显式 pathspec；无 `add -A`／`amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；未删任何文件（含四发整包原始日志留在 `logs/`，见下）。

**未入库的原始日志（具名，按规矩只建不删，尺＝`wc -c .scratch/wisp/probes/296/r1/logs/*.md`）**：
`logs/prefix-full-1.md` 196,715／`logs/prefix-full-2.md` 316,731／`logs/postfix-full-1.md` 313,893／`logs/postfix-full-2.md` 313,932（四发 `-v` 全文，合计约 1.1 MB）
／`logs/d22scan.md` 23,086／`logs/build-final.md` 523／`logs/build-tail.md` 11／`logs/build.md` **0 字节**。
⇒ 那 0 字节的 `logs/build.md` 是 `GOFLAGS= go build ./...` **成功时的正常零输出**，不是"那一格没交"：
build 这一把的 rc 与名册落在 `logs/build-final.md`（自落 `build rc=0`／`d22scan rc=0` 两行）与本件 §1，
判据用的正是"件自己带一行 rc=N"那一把尺。

