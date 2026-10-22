# 票 360 终裁腿 `360-v1b` 判语 —— `360-r2` 的补强堵住了 `360-v1` 判的那个洞

- 身份：**非实现者终裁腿**（裁决者≠实现者）。⛔ 落任何 `.go`／⛔ 翻任何框／⛔ 改票名·`Status:`·`-done`·台账／⛔ push。
- 时间：2026-10-22 21:4x–22:2x +08。判据树 HEAD `02b3882c`（开跑 `5b6338f6`；其间他腿 9 笔，`internal/config/` 逐字节未变，见 ⑨）。
- 被测件：`internal/config/writeguard_adopt_360_test.go` **107 行**、blob `7077bc2b353f0009e593b5cebd3aa3a35d13bb4c`；★**HEAD＝工作树＝仓外导出树三处同串**。
- 变异台：仓外 `D:\tmp\wisp\20261022\360v1b\tree`（`git archive` 直出 tar；**17582/17582 件、1156/1156 `.go`**，零缺件、零乱码名）。
- 结论：**`AC#2` 能翻**。那条前置断**不可被同臂断言吞掉** ⇒ `M1＋M3b` 与"全弄哑"两发都**转红在它**（逐字见 ③④）。

## ① 票面前提新鲜度（HEAD blob 逐条复量）

- **① `statOwnWrite` 分层（raw grep 5 行）**：定义 `allowdirs.go:157`（体 `:157-163`，只写 `seenMtime`/`seenSize`/`seenValid`）；代码调用 `writeguard.go:179`（`mergeWrite` 的 `len(diverged)==0` 支内）；注释 3 行 `allowdirs.go:150`／`writeguard.go:105`／`writeguard_226_test.go:186`。`permmode.go` **0** 命中。HEAD 与工作树逐行同。
- **② 两条路**：写回入口 3 枚 `allowdirs.go:62→:139→:142`／`permmode.go:70→:79`／`settings.go:259→:303`，汇聚 `writeguard.go:111 mergeWrite`（`:118` 读盘当 base、`:163 SaveFile` 从不回读）；Adopt 全仓**仅一条** `manager.go:142 CheckAndReload→:160 LoadFile→:193 plan.commit` 写 `cur`，生产驱动 `cmd/wisp/config_reload.go:157`。⇒ **两者互不相连，前提成立**。
- **③ 复现**：`STATE 1 after writeback: file.ball.size=60 mem.ball.size=56`（主检出**连跑 5 次不抖**；导出树同）。
- **④ 正控**：`STATE 2 after one Adopt pass: mem.ball.size=60 file.ball.size=60 report.Hot=[ball]`（同上）。
- **⑤ `M-ADOPT` 前提（票面"红 8 枚"）**：我在导出树 `writeguard.go:169` 后插 `m.cur = deepCopyConfig(base)` ⇒ `go build` **rc=0**、`internal/config` 红 **9 枚**：`TestAC11ControlSnapshotWriteIsWhatRevertsTheHandEdit`／`TestAC1AllowedDirsWriteKeepsAHandEditedKey`／`TestAC3AdoptionClaimsOnlyWhatThisWriteProduced`／`TestAC3HandAddedEntryIsPreservedButNotAppliedToMemory`／`TestAC4PermissionModeWriteKeepsAHandEditedKey`／**`TestWritebackDoesNotAdoptMemoryTheFileKeepsTheHandEdit`**／`TestHandDeletedScalarLineComesBackAtSchemaDefault`／`TestHandDeletedWholeSectionComesBackFromDefaultsNotMemory`／`TestHandDeletedProviderBaseUrlComesBackAsThePresetValue`。
  ★**9 ＝ 票面那 8 枚 ＋ 本票新用例 1 枚**（新用例当年还不存在、量不到）⇒ 票面"8"的**数字今天已不可逐字复现，但结论（那 8 枚全是他票格的顺带断言）成立**；且新用例在 `:60`／`:65`／`:99` 红 ⇒ **它对产码侧那一跳有牙**。还原 `sha256 efe12661…47bd` ＝ HEAD blob。
- **⑥ "间接已钉"那一发：★我没复量**。原因＝**环境限制**：`cmd/wisp` 测试二进制在本机导出树起不来（`exit status 0xc0000135`＝`STATUS_DLL_NOT_FOUND`，**未变异基线同码**）⇒ 非突变所致；`TestAC6AHandEditIsStillSeenUnderTheRealPoll` 全仓只在 `cmd/wisp/always_write_no_clobber_227_test.go:125`，`internal/config` 无同题用例。★⛔ 当"我复现了"。

## ② 本格唯一判点：那条前置断是否真的不可被吞掉

**判语：不可被吞掉 ⇒ `AC#2` 能翻。** 四发实测（导出树、逐发 `go vet ./internal/config/` rc=0 证编译、`go build ./...` rc=0）：

| 发 | 变异（逐字 diff 已落 `logs/s03-mutations.log`） | rc | 红在哪 |
|---|---|---|---|
| `M1` | 删 `:75-84` 整条 `if/else` Adopt 臂 | **1** | `:84` 前置断 |
| ★`M1＋M3b` | 同上 ＋ `:98` → `if false && (…)` | **1** | `:84` 前置断 |
| ★"更狠" | 删臂 ＋ `:55`／`:59`／`:64`／`:98` 全弄哑 | **1** | `:84` 前置断 |
| 加固 | 上一发 ＋ `:103` 两态守卫也弄哑 | **1** | `:84` 前置断 |

- **机理（我现读，⛔ 采信自报）**：前置断只读 `ranAdopt`／`adoptedMem`／`adoptedFile` 三个哨兵；摘臂后 `ranAdopt=false`、两哨兵仍 `-1` ⇒ **它先咬**；`:98` 的恒真化与 `:103` 的守卫都不在它的表达式里 ⇒ **吞不掉**。
- **它另有一条牙（我实测）**：产码侧把 `CheckAndReload` 改成 `return nil, nil`（臂"跑过"但不 adopt）⇒ `:98` 活时红在 `:99`；`:98` 弄哑时**仍红在 `:104`**（两态守卫）⇒ "臂跑了却没接管"这一形也跑不掉。

## ③ `M1＋M3b`（我自己造的）逐字红行

变异：删 `:75-84` ＋ `:98` 改 `if false && (…)`（`grep -c 'if false &&'`＝1、`CheckAndReload` 剩 1 处＝`:19` 注释、`ranAdopt = true`＝0；`go vet` rc=0；diff 逐字落日志）。读数：
```text
    writeguard_adopt_360_test.go:54: STATE 1 after writeback: file.ball.size=60 mem.ball.size=56
    writeguard_adopt_360_test.go:75: STATE 2 after one Adopt pass: mem.ball.size=-1 file.ball.size=-1 report.Hot=[]
    writeguard_adopt_360_test.go:84: the Adopt arm did not run: ranAdopt=false mem.ball.size=-1 file.ball.size=-1; the positive control never executed, so the state 1 readings above are unbacked
--- FAIL: TestWritebackDoesNotAdoptMemoryTheFileKeepsTheHandEdit (0.60s)
FAIL	github.com/CarlosShao/wisp/internal/config	0.830s
FAIL
```
（`M1` 单发红行同串、红在 `:84` 同句。★`360-v1` 那发当年是 `--- PASS`，今天是 `--- FAIL`。）

## ④ "更狠那一发"（全弄哑）逐字红行

变异：删臂 ＋ `:55`／`:59`／`:64`／`:98` 四条皆 `if false && …`（`:64` 形为 `; false && !slices.Equal(…)`；diff 落日志）。整条用例**只剩前置断在咬**：
```text
    writeguard_adopt_360_test.go:54: STATE 1 after writeback: file.ball.size=60 mem.ball.size=56
    writeguard_adopt_360_test.go:75: STATE 2 after one Adopt pass: mem.ball.size=-1 file.ball.size=-1 report.Hot=[]
    writeguard_adopt_360_test.go:84: the Adopt arm did not run: ranAdopt=false mem.ball.size=-1 file.ball.size=-1; the positive control never executed, so the state 1 readings above are unbacked
--- FAIL: TestWritebackDoesNotAdoptMemoryTheFileKeepsTheHandEdit (0.19s)
FAIL	github.com/CarlosShao/wisp/internal/config	0.372s
FAIL
```
加固发（`:103` 也哑）逐字同 `:84` 红行、`--- FAIL … (0.11s)`。⇒ **判据"改恒真"那一半补上了**。

## ⑤ 原样必须仍绿（主检出、未变异树）

`go test -count=1 -run '^TestWritebackDoesNotAdoptMemoryTheFileKeepsTheHandEdit$' -v ./internal/config/` ⇒ **rc=0、PASS**，两态读数逐字：
```text
    writeguard_adopt_360_test.go:54: STATE 1 after writeback: file.ball.size=60 mem.ball.size=56
    writeguard_adopt_360_test.go:85: STATE 2 after one Adopt pass: mem.ball.size=60 file.ball.size=60 report.Hot=[ball]
```
★连跑 5 次（0.08/0.18/0.09/0.05/0.05s）**两态逐字同、不抖**；导出树原样同串（0.14s）。

## ⑥ 零产码复核（逐枚 `rev-parse HEAD:<p>` 对 `hash-object <p>`）

| 件 | HEAD ＝ 工作树 ＝ 导出树 |
|---|---|
| `internal/config/allowdirs.go` | `ef54b37c094ffd1349bcc8b27752313046a9b968` SAME |
| `internal/config/writeguard.go` | `0ebc415b05651104fa13d6e71c76e61ba8fdd93c` SAME |
| `internal/config/manager.go` | `f53d8a0749cd5ceda1cc6700c7210a2ff019e063` SAME |
| `internal/config/settings.go` | `c4e3b09b21316a897412d7f04aff9e2406ff4061` SAME |
| `internal/config/permmode.go` | `6f385353ec366ed1576512855bbd80ff97cd5c97` SAME |

★更宽一把：`internal/config` **17 枚非测试 `.go`、工作树 diff 0**。★另见一处**不是我**的树脏：` M internal/risk/provenance.go` 在 22:03 出现、22:12 已自行消失（他腿瞬态；`git diff 095699f2^..7ad0c83e4 -- internal/risk/` ＝ **0 枚**）。

## ⑦ 门禁（我自己现跑，HEAD `02b3882c`）

`GOFLAGS=` 已置空、`GOTMPDIR TMPDIR TEMP TMP` 已 unset、⛔ `CGO_ENABLED=0`。

- `GOFLAGS= go build ./...` **rc=0**。
- `go vet ./internal/config/` **rc=0**。
- `GOFLAGS= go test -count=1 -v ./internal/config/` **rc=0**；四数分开：`RUN=196`／`PASS=115`／`FAIL=0`／`SKIP=0`（另 `^panic:` **0**）；末行逐字 `ok  	github.com/CarlosShao/wisp/internal/config	6.564s`。
- `sh scripts/d22scan.sh` **rc=0**；末行逐字：`d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=247, bans #1-5 cmd/=42, ban #6 frontend/=116, ban #7 internal/tools/=25, ban #8 design/=39, ban #8 frontend/=116, ban #8 internal/=619, ban #8 cmd/=148; ban #8 emoji coverage: design/ 39 text files; frontend/ 116 text files; internal/ 619 Go files, comments and _test.go included; cmd/ 148 Go files, comments and _test.go included`。
- `sh scripts/check-path-length-budget.sh` **rc=0**；末行逐字：`check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`。
- **格式门禁甲乙两尺**（裁定件 `pipe-format-gate-reads-tracked-blob-not-worktree.md`）：
  - **甲（入库 blob 尺，判据）**：HEAD 树 `git archive` 到仓外，**1158 枚** `.go` 跑 `gofumpt -l` ⇒ 列 **20 枚**、**全在 `.scratch/wisp/probes/**`（非 scratch 0）**、★**本票件命中 0**、rc=0、stderr **0**。
  - **乙（工作树尺，只作诊断）**：`-l` 的 **stdout 49 条路径** ＋ **stderr 5 条 parse 错**。★**顺带记我自己的一处尺的射程**：我第一把尺把两股混数成 **54**——那是错的，**49 才是文件数**；`xargs` rc=123 只是 gofumpt 对坏探针 exit 2 的回声。
  - **乙的 49 分解**：**20** ＝ 甲那 20 枚真脏；余 **29** ＝ **14** 枚**纯 CRLF 检出副本**（含 `cmd/wisp/models.go` wt CR=389／blob CR=0、`internal/agent/tools.go` 225/0、`internal/agent/approval/pending_read.go` 131/0，三者 `git status` 全**干净**）＋ **15** 枚**非跟踪** `.scratch` 探针（11 `??` ＋ 4 被 ignore）。★本票件在乙也 **0** 命中。
  - ★归属用 `[\\/]`：gofumpt 输出**反斜杠**（甲日志里 `\` 计数 **116**），我全程 `tr '\\' '/'` 后判归属。

## ⑧ 红名册（逐名，⛔ 只报枚数；取主检出）

`internal/config` 整包 **FAIL=0** ⇒ **名册为空**。**新增红＝0**。

- **结构尺**（同树、非跨发比数）：`go list -f '{{.ImportPath}}|{{join .TestGoFiles ","}}' ./...` ⇒ **只有 1 个包**把该件列进 `TestGoFiles`＝`internal/config`（21 枚测试件）⇒ 它只编进这一个测试二进制。
- **行为尺**：该件在 `M-ADOPT` 突变下**会红**（`:60`／`:65`／`:99`）⇒ 不是被 skip 的死件。

## ⑨ `360-r2` 三笔提交（写面纪律：逐笔 `--name-only`）

- `095699f2`：1 件 `internal/config/writeguard_adopt_360_test.go`，`.go` **1**。
- `16898169a`：2 件 `docs/evidence/s1/360-r2-impl.md`／`-part2.md`，`.go` **0**。
- `7ad0c83e4`：1 件 `docs/evidence/s1/360-r2-impl-part2.md`，`.go` **0**。
- 并集 `095699f2^..7ad0c83e4`：**3 件**、`.go` 1、**越界 0**（白名单＝`internal/config/*_test.go` ＋ `docs/evidence/s1/`）。
- **正控**（同尺打一笔**含 `.go`** 的真提交）：`ecebc7a9` ⇒ `.go` **2** ⇒ 这把尺不空转。
- 票面未被动：票文件 `hash-object` ＝ `rev-parse HEAD:` 同串 `81dd684d…`；`origin/dev..HEAD` ＝ **428**（非空）⇒ 腿笔只在本地、**⛔ push**。

## ★ `AC#2` 判语（含"缺什么"）

**能翻。** 缺的那一半（`360-v1` 判的 `改恒真`）已由 `:93-97` 前置断补上；我三发独立复现"摘臂／摘臂＋恒真／全弄哑 ⇒ 全红在它"。**三条边界须写明（⛔ 读大）**：

1. **`M3b 单独恒真仍绿` ⇒ ⛔ 构成缺口（我判）**。实测：只把 `:98` 弄哑 ⇒ PASS、臂**真跑过**（`:85` 打印 `mem=60 file=60 hot=[ball]`、哨兵 60/60）。理由：此时**还剩两条独立牙**——`:93`（臂必须跑）＋`:103`（两态必须不同）；且我实测"`:98` 哑 ＋ 产码侧 Adopt 路弄死" ⇒ **仍红在 `:104`**。同臂断言的**删除**在行为正确时**任何 in-test 断言都测不到**（要测只能重复同一断言，买不到价值）。
2. **整段 STATE 2 删除（`:68-106`）⇒ 仍 PASS**（我实测 rc=0）。这不在 `AC#2`"摘掉 Adopt 臂"的字面射程内，属"掏空整条用例"的通用边界、in-test 无法防。
3. **判据原文的"或"只有按"对"读才可满足**（摘臂 **＋** 弄哑该臂**唯一**那条牙）。这是 `360-v1` 自己操作化的读法；⛔ 我代改契约文字，也⛔ 把 1／2 两条算成 r2 的缺陷。

⇒ 建议（⛔ 代做）：若要把判据写死，把"或"改成"**且**"，或加一句"以摘臂＋弄哑该臂自己那条断的合发为准"——属**改契约文字**，须机主批。

## ★ 我 ⛔ 做什么／没覆盖什么

- ⛔ 改任何 `.go`、⛔ 翻框、⛔ 改票面／`Status:`／`-done`／台账、⛔ push。本件是唯一仓内写面。
- **未复量**：⑥ 的"间接那一发"（`cmd/wisp` 测试二进制本机 `0xc0000135`，未变异基线同码 ⇒ 环境）；票面 ⑨ 的导出树件数（当刻 17569/17575，今天 17582/17586 ⇒ 只随时间增长、非承重）。
- 变异与整包门禁**串行**：突变只在仓外树、每发后 `sha256` 还原并复核（`writeguard.go`＝`efe12661…47bd`、`manager.go`＝`4bb091d9…`、测试件＝`52a1aaee…`、`DIFF=0`）。

## 附：可复跑入口（仓外，只建不删）

`D:\tmp\wisp\20261022\360v1b\`：`s03-mutations.sh`／`s04-gates.sh`／`s06-premise.sh`／`s09-format-scope.sh`／`s12-dead-adopt.sh`／`s13-m3b-dead.sh`／`s14-gut.sh`／`s16-regates.sh` ＋ `logs/*.log`（逐发原始输出）。
