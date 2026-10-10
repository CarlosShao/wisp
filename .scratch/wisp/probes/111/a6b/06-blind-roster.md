# 件1＋件3 — 票 111 `AC#12` 盲区名册 与 "加一枚路径的连带代价"（腿 111-a6b，HEAD `f76091cd`，2026-10-10 17:0x）

⛔ 裁形、⛔ 推荐。本件只给现量与尺。所有读数走对象层（`git archive HEAD` → 仓外临时树，`00b-archive-check.txt` 里有逐枚 cksum 对拉），
**本腿 `go build`/`go test`/`go vet`/`go run`/`go list` 跑了 0 次**。

## 0. 尺（本件用的三枚，可逐字重跑）

| 尺 | 命令形状（重跑＝`bash .scratch/wisp/probes/111/a6b/01-rulers.sh`） | 射程 | 整族/抽样 |
|---|---|---|---|
| **R1** 内容 tag 尺 | `head -3 <blob>` 里找 `^//go:build` ⇒ 类别 `WIN / WIN_LIVE / POSIX / POSIX_OTHER / TAG_OTHER / NONE` | `git archive HEAD -- cmd internal` 全部 376 枚 tracked `_test.go`（⛔ `.scratch/**`、⛔ 自带 `go.mod` 的 `scripts/spike`／`tools/d22scan`／`tools/mockllm`） | **整族** |
| **R2** 顶层用例尺 | `grep -c -e '^func (Test\|Benchmark\|Fuzz\|Example)'` | 同 R1 | 整族（**顶层口径，⛔ 含子测试**） |
| **R5** 认领尺 | 四枚 pin（`core_pin` 27 行／`win_pin` 11 行／`cli_pin` 1／`winsec_pin` 1，从 blob 抽）× `ci.yml` 各 job 的 `runs-on` ⇒ 包→档→runner-OS | `scripts/portable-tests.sh` blob ＋ `.github/workflows/ci.yml` blob | 整族 |
| **R4** 前置条件尺 | 逐枚取**函数体**，`index()` 问 20 个门闸字样（`WISP_LIVE_MIC`／`t.Skip(`／`t.Setenv`／`net.Listen`／`net.Dial`／`http.Get`／`exec.Command`／`registry.`／`HWND`／`webview`／`frontend/dist`／`18080`／`CreateWindow`／…） | 4 枚盲包的 450 枚用例全量 | 整族（名册全列在 `04-case-preconditions.tsv`） |

**⛔ 没有用文件名分类那把尺**（票面点名的陷阱本腿现量到了，见 §3）。

## 1. ★盲区名册（定义＝包被某档认领，但那档跑的 OS 上永远不编译它那半边）

**判据**：`R1` 给出包里 `WIN`／`WIN_LIVE` 用例枚数，`R5` 给出认领它的档与 runner OS；认领档里没有一枚跑 windows ⇒ `BLIND_WIN=1`。

| 包 | 档（runner OS） | windows-tagged 文件／用例 | 该包顶层用例总数 | 无 tag 文件／用例 |
|---|---|---|---|---|
| `internal/agent` | core（ubuntu） | 1／1 | 85 | 20／84 |
| `internal/memory` | core（ubuntu） | 1／1 | 37 | 10／36 |
| `internal/models` | core（ubuntu） | 3／3 | 45 | 13／39 |
| `internal/tools` | core（ubuntu） | 10／27 | 209 | 43／182 |
| **合计** | — | **15／32** | **376** | **86／341** |

⇒ **4 枚包／15 枚文件／32 枚用例**。这枚数与票面 `AC#12` 现量②"票 301 `AC#1` 落地后＝4／15／32"**逐枚同**（本腿是独立尺：R1 问内容、R5 问 blob pin）。
交叉核对全模块：带 windows-tagged 用例的模块包＝**12 枚**、tagged 文件＝93＋12(winlive)＝**105 枚**、tagged 顶层用例＝349＋29＝**378 枚**
＝ 票面现量①的第二发（@`41329475`＝105／378）**逐字复现**（含票 300 那枚 `TestParseWaveFormatSubFormatOffset300`）。

其余 8 枚带 windows 用例的包（`cmd/wisp`／`internal/audio`／`ball`／`config`／`proc`／`risk`／`secret`／`winsec`）**不盲**：
它们分别被 `cli`（windows）、`windows`（windows＝`win_pin` 11 枚）或 `winsec` 档（`runs-on: windows-latest` 的 ACL 那一步）认领。

### 1b. 反形（票面没写、同一把尺看见的）＝`BLIND_POSIX`

| 包 | 档（runner OS） | `!windows` 文件／用例 | 后果 |
|---|---|---|---|
| `cmd/wisp` | cli（windows，唯一认领档） | 1／3 — `cmd/wisp/secret_dataroot_119b_test.go`（`//go:build !windows`） | 那 3 枚 POSIX 半边**没有任何档会跑**：`cli` 只在 `test-windows` 被调用（`ci.yml` 里 `--scope=cli` 零枚调用点，调用者是 `scripts/wisp-cli-tests.sh`），`cmd/wisp` ⛔ 在 `core_pin` 里 |

⇒ 收口如果只做"windows 那半边"，`cmd/wisp` 的 POSIX 半边仍是盲区且新门⛔ 抓不到（同一形、方向反）。
全模块 `!windows` 文件＝12 枚／38 枚用例，其中只有 `cmd/wisp` 这 1／3 落在"⛔ ubuntu 档认领"的包里。

## 2. ★件3：连带名册（"往清单里加一枚路径"到底拉进多少枚）

票面 `AC#12` 排程段要求落地腿起手做这件事，本腿先量出**规模**：把 4 枚盲包逐枚拉进 `windows` 档（`scope=(…)` ＋ `win_pin` 各加一行）时，
windows job 的分母**不是** 32 枚 tagged 用例，而是那 4 枚包里**为 GOOS=windows 编译出来的全部顶层用例**：

| 若拉进 | 目标 tagged 用例 | 被连带拉进（无 tag，**今天从未在 windows 上求值过**） | 该包 windows 分母 | 被排除（`!windows`） |
|---|---|---|---|---|
| `internal/agent` | 1 | 84 | 85 | 0 |
| `internal/memory` | 1 | 36 | 37 | 0 |
| `internal/models` | 3 | 39 | 42 | 3 |
| `internal/tools` | 27 | 182 | 209 | 0 |
| **合计** | **32** | **341** | **373** | 3 |

⇒ **买 32 枚可见性＝同时把 341 枚从未在 windows 求值过的用例塞进 windows 腿的分母**（放大倍数 11.7 倍）。
逐枚名册＝`04-case-preconditions.tsv`（450 行，含包／文件／用例名／体内门闸字样），前置条件类别现量：

| 前置条件（按**用例体内容**判，⛔ 按名字猜） | 枚数（4 盲包全量 450 枚里） | 尺注 |
|---|---|---|
| `t.Setenv`（自带 env，⛔ 外部依赖） | 12 | R4 |
| 含子串 `mic` | 11 | ⚠ **R4 这一味是子串尺，会误命中（如 `chemical`）⇒ ⛔ 当"真麦克风"用**；`WISP_LIVE_MIC` 命中 **0 枚** |
| `t.Skip(` | 7 | 其中 **1 枚是 windows-tagged 目标用例**：`internal/tools/recycle_windows_test.go` 的 `TestShellTrashLeavesARecycleBinEntryTheOSCanEnumerate` |
| `time.Sleep` | 6 | ⛔ 等价于"要设备"，只是时序 |
| `exec.Command`（起真子进程／`mklink /J`／`taskkill`） | 4 | 需要宿主能起进程＋能造 reparse point（托管 runner 上成不成＝**欠读数**） |
| `watcher` | 2 | — |
| `HWND` | 1 | 只问结构体布局，⛔ 不需要真窗 |
| `WISP_LIVE_MIC`／`net.Listen`／`net.Dial`／`http.Get`／`18080`／`webview`／`frontend/dist`／`CreateWindow` | **0** | ⇒ 这 373 枚**不吃麦克风、不吃网络、不吃真窗**（网络那味只在别的包里） |

### 2b. 三枚"连带"的具体红法（现量，⛔ 推断）

1. **`SKIP` 变红那一枚**（票面 ⓑ 说的"正面相遇"）：`TestShellTrashLeavesARecycleBinEntryTheOSCanEnumerate`（`recycle_windows_test.go`，`//go:build windows`，体内有 `t.Skip(`）
   一旦被拉进 windows 腿就成了**该腿的分母**；而 `tools/d22scan/runtests.sh`（`portable-tests.sh:15` 自述"Every verdict comes from tools/d22scan/runtests.sh"，规则 3 原文 `:19` "zero top-level PASS *and* zero FAIL is fatal"）
   把 SKIP 记成非 pass ⇒ **这一枚是"拉进清单"当天就可能把 `test-windows` 弄红的具体对象**（它 skip 不 skip ＝ 欠读数 §4-3）。
2. **ledger 会指错方向**：`internal/tools` 有 3 行 ledger 登记的是 `platform=linux`（`TestWorkspaceSwitchRefusesAJunctionToOutside`／`TestD34WriteMatrix`／`TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop`，原文自陈"on windows it … runs and passes"）。
   ⇒ 用例级/档级两形若把 ledger 的 `platform` 列读成"该用例只在那台 OS 有分母"，这 3 行会把**windows 那侧**读成豁免、把 linux 那侧留成断言＝**方向恰好相反**，需要落地腿具名处理。
3. **GUARD A 与"用例级"会给同一件事两种颜色**：`internal/tools` 的 windows 分母＝209 枚（非零）⇒ GUARD A（`:566` 起，判"声明在 scope 里但本平台编译出 0 枚测试文件"）**⛔ 红**；
   与票面 ⓓ 的"谁赢"问题在这里是**不冲突**的（票面那 5 枚盲包的无 tag 文件也都非零）⇒ 冲突真发生在**包整体为 0 分母**的那形，本腿现量：4 枚盲包在 windows 上**全都非零**，
   所以 ⓓ 那枚冲突今天**没有活样本**（这是一枚读数，⛔ 一枚安慰）。

## 3. 分类陷阱现量（⚠ 派单点名"必须问文件内容"那一枚）

`R1`（内容）与文件名尺（路径含 `_windows_test.go`）逐枚对拉：**16 枚文件内容带 `windows` tag 而文件名不带**，全名册：

```
cmd/wisp/panel_geometry_255_test.go (5)      cmd/wisp/panel_inbound_guards_35r3_test.go (3)
cmd/wisp/panel_transport_35r1_test.go (1)    cmd/wisp/panel_transport_35r2_test.go (9)
cmd/wisp/resident_hotkey_v1probe_test.go (2) internal/audio/hotplug_test.go (8)   ← 票面引的先例，本腿坐实
internal/ball/hotkey_borrow_refused_260r5_test.go (3)  internal/ball/hotkey_cancel_borrow_260_test.go (7)
internal/ball/hotkey_cancel_borrow_expect_260r2_test.go (3)  internal/ball/hotkey_status_test.go (11)
internal/ball/hotkey_test.go (2)             internal/secret/delete_test.go (5)
internal/secret/migrate_test.go (5)          internal/secret/store_test.go (7)
internal/tools/paths_shortname_252_probe_test.go (3)   internal/tools/paths_shortname_252_r1_test.go (4)
```

⇒ **盲区包 `internal/tools` 里就有 2 枚文件／7 枚用例命中这一形**（`paths_shortname_252_probe_test.go` 3 枚＋`paths_shortname_252_r1_test.go` 4 枚，体内都是 `//go:build windows`）。
⇒ **按文件名分类那把尺会把 `internal/tools` 的盲区从 27 枚读成 20 枚**（少 7 枚），并把"哪个文件需要登记"整个读错。
其余 14 枚（78 枚用例）落在⛔ 盲包：`cmd/wisp` 5 枚／20 枚、`internal/ball` 5 枚／26 枚、`internal/secret` 3 枚／17 枚、`internal/audio` 1 枚／8 枚（＝票面引的 `hotplug_test.go` 先例本身）。
反方向也有：**1 枚**文件名带 `_windows` 而内容⛔ 没有 tag（`cmd/wisp/leg_sink_nail_131_windows_test.go`，`//go:build` 行缺失＝NONE 类，3 枚用例，两平台都编译）。
`tag` 只在前 3 行找（票面"⛔ 全文 grep"那枚陷阱照本腿口径执行）；`TAG_OTHER` 类全模块只有 **1 枚**：`internal/risk/pathresolver_budget_norace_test.go`（`//go:build !race`）＝⛔ 平台 tag，⛔ 算进盲区。

## 4. 欠读数（具名＋谁能取；本腿⛔ 自造）

| # | 欠哪一枚 | 为什么本腿⛔ 能取 | 谁能取 |
|---|---|---|---|
| 1 | `--scope=census` 今天的真实退码与 `unclaimed-with-tests=` 读数 | 要 `go list ./...`＝编译面，`303-a1` 独占台面（票面现量③本身就标〔仅腿量〕） | 落地腿起手／303-a1 收台后的验收腿／编排者开的专用窗口 |
| 2 | `go list ./...` 的**包全集枚数**（票面写 35，本腿的尺只看有 `_test.go` 的 28 枚目录，⛔ 包⛔ 无测试文件＝本尺看不见） | 同上 | 同上 |
| 3 | 341 枚连带用例在托管 windows runner 上**红不红**（尤其 §2b-1 那枚 skip、§2b-2 那 3 行 junction/双卷形状） | 要跑 `go test` | 落地腿跑 `--scope=windows`（＋路径）前后各一发取具名差集（票面排程段就是这么要求的） |
| 4 | `TestShellTrashLeavesARecycleBinEntryTheOSCanEnumerate` 的 `t.Skip` 条件到底是不是"无交互 shell" | 要读体＋跑 | 落地腿（本腿只量到"体内有 `t.Skip(`"这一枚事实） |
| 5 | 11 枚含子串 `mic` 的用例到底⛔ 是不是麦克风门（子串尺的假阳性未消） | ⛔ 编译面无关，但要把 11 枚体逐枚读完再分类；本腿按预算⛔ 做（见 §5） | 落地腿逐枚读体 |

## 5. 本腿自报：`go env` 跑了 **0** 次，`go list`／`go build`／`go test`／`go vet`／`go run` **各 0 次**；
只用了 `git`（`show`／`ls-files`／`grep`／`archive`）、`tar`、`grep`、`awk`、`wc`、`cksum`、`tasklist`。

rc=0  # 本件由 111-a6b 撰写；数字来路＝`01-file-class.tsv`／`02-pkg-summary.tsv`／`03-blind-roster.tsv`／`04-case-preconditions.tsv`／`05-ruler-meta.txt`（同批，每枚自带 rc 行）
