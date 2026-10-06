# 236 AC#2 的交件（判语版，腿 `236-r3c`）

射程＝票 236 的 **AC#2 一格**：把那把 DEFERRED 尺的 **(b) 支**从"数词面"改成"看能力"——看 `BuiltinTaskEntries`
的注册名册里有没有 `task.list` 这一行——并**先证它今天无牙**（`teeth-m13` 那一形：摘掉标记行仍然 PASS＝未修码读数）。
AC#1／AC#1b／AC#3／AC#4／AC#5／AC#6 一枚不做、不评、不顺手；票面 7 枚未勾框（含 AC#2 那一枚）**一枚未翻**。

写法上两条口径：

1. **本腿只跑四发突变**＝派单的硬约束（前两枚腿 `236-r3`／`236-r3b` 都死在 75–80 次调用附近，本腿把发数砍小）。
   四发＝正控／`m-2b`／`m-2e`／`teeth-m13`，另加**三发仪器配对**（票面点名要的"配正控"）。
   ⛔ **未重跑整包 `go test ./internal/...`**——那是 AC#1 那格（`236-r2` 已交件）的判据，不是 AC#2 的。
2. **归属（派单要求逐枚指名）**：本文件原有 36 行是 `236-r3` 的骨架（每条以"打算答："开头，是实话句、不是判语），
   **逐字保留**在本件文末附录（临时件只建不删），其"本节以下目前是骨架"那句以本件为准。
   `logs/mut/**`（29 枚）、`logs/overlay/**`（11 枚）、`logs/gates/**`（11 枚）、`logs/baseline/**`（3 枚）＝**`236-r3` 跑的**，
   本腿只当**形状参考**（overlay json 的写法、合成副本放 `D:/tmp/` 的做法），⛔ 未引用其中任何 rc／红名册／判语当自己的凭据。
   **本腿自己跑的读数全部在新目录 `.scratch/wisp/probes/236/r3/logs/r3c/**`**（现量枚数：`mut/` 23、`overlay/` 7、`gates/` 10）。

用例载体（`internal/tools/tasklist_deferred_236r3_teeth_test.go`，HEAD 版 200 行）＝`236-r3` 落库（提交 `d9aff5f`，
编排者标〔未验证半成品〕代提）。**本腿对它的改动＝注释三处**（新增"读数归属"一节；把 `M-2a`／`M-2c` 两行按本腿实测改准，
如实写明 `M-2c` 未实测）。零断言改动、零用例增删、零导出名、零产码改动。
文件字节：`8f536366cfe481e446f640d08e460673`（HEAD 原版）→ `21b53d52ac87b893b74fdd9828f2f3aa`（本腿改后）。
★因此本件引用**测试文件行号**按轮次分开看：pass1 四发跑在 `8f536366`（红句里是 `:132`／`:151`／`:189`／`:195`），
三发配对与 pass2 跑在 `21b53d52`（同一句落在 `:144`／`:163`／`:201`／`:207`）。
**两遍红名册逐字相同**（`logs/r3c/mut/pass2-readings.txt`）⇒ 名册不是某次文件字节或缓存的产物（全程带 `-count=1`）。

---

## §0 起手锚（四把尺原文＋本腿 HEAD）

| 尺 | 命令 | 原文读数（本腿现跑，10-06 16:0x） |
|---|---|---|
| 撤票检查 | `sed -n '3,5p' .scratch/wisp/issues/236-six-cells-that-only-surface-at-the-reading-layer.md` | 第 3 行＝「**Status**：**待派**（编排者 09-29 20:2x 立，起手锚点 `c5d88a7f`＝本票现读 `git rev-parse --short HEAD`）。」；第 4 行＝来源两枚（`docs/evidence/s1/221-task-cancel-v1.md`／`.scratch/wisp/probes/ci-red/ci-red-1.md`，台账 `A451`）；第 5 行＝「**本票的射程不是"功能没做"，是"没人能证明它没坏"**」。**零 `WITHDRAWN`、零"撤"、零"作废"** ⇒ 票活着，可开工 |
| HEAD／时间 | `git log -1 --format='%h %ad %s' --date=format:'%m-%d %H:%M'` ＋ `date '+%m-%d %H:%M'` | `d9aff5f 10-06 15:25 probes(236-r3 死腿半成品代提)：一枚能编译、自己会红，但证件只有骨架——编排者只代提、不代判，标〔未验证半成品〕` ＋ `10-06 16:09`。**本腿 HEAD＝`d9aff5f`** |
| 工作树（产码面） | `git status --porcelain -- internal cmd` | 起手 **0 行**（与派单 16:08 现量一致）；四发突变之后仍 **0 行**；本腿注释编辑之后＝唯一一行 ` M internal/tools/tasklist_deferred_236r3_teeth_test.go`；commit 之后回 0 行（`logs/r3c/gates/status.txt`／`final-gates.txt`／`final-gates-second-pass.txt`） |
| 框数 | `grep -n '^- \[ \]' .scratch/wisp/issues/236-*.md` ＋ `grep -c '^- \[x\]'` | 未勾 **7 枚**＝`:26` AC#1／`:27` AC#2／`:28` AC#3／`:29` AC#4／`:30` AC#5／`:31` AC#6／`:85` AC#1b；已勾 **0 枚**（`grep -c` 无匹配）。⛔ 七枚一枚不归本腿，本件不翻任何一枚 |
| 基线（单包） | `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./internal/tools/ -count=1` | `ok  	github.com/CarlosShao/wisp/internal/tools	15.318s`，**rc=0** ⇒ 未突变态是绿的，§3 每一枚红都只能来自那一行突变 |
| 基线名册（同包 `-v`） | 同上追加 `-v`，输出落 `logs/r3c/mut/baseline-full-package-v.log` | 顶层 **PASS=206／FAIL=0／SKIP=0**、`=== RUN` 281、缩进 `--- FAIL` **0 枚**、`ok 13.924s`。★口径：206 是**顶层枚数**（票面 §六 逐字警告过行首锚定尺看不见缩进子测试；本腿判据不依赖名级计数）。交叉核对：`236-r2` 交件报 `PASS=204`，`236-r3` 之后多两枚＝`Test236R3TaskListDeferralIsPinnedByTheRosterNotByAComment`＋`Test236R3InstrumentFactReadDiskIsBlindToOverlayCarriesIt` ⇒ 204＋2＝206 对得上 |

绿名册里与本格直接相关的四枚（`-v` 原文行，逐字）：

```
--- PASS: Test221TaskCancelIsRegisteredAtItsFrozenLevel (0.00s)
--- PASS: Test221DeferredMarkerForCancelLiftedButListStillMarked (0.00s)
--- PASS: Test236R3TaskListDeferralIsPinnedByTheRosterNotByAComment (0.00s)
--- PASS: Test236R3InstrumentFactReadDiskIsBlindToOverlayCarriesIt (0.00s)
```

⛔ 本腿与在飞的腿无交叠：`probes/evidence-close/5/**` **未读**、结论**未引用**；`probes/pool-validity/4e/**` **零读零写**。
起手后 HEAD 被推进过 `f84b455 10-06 16:12 ledger(A639)`（只动 `docs/**`／台账，非本腿所为），本腿未与它抢任何文件。

---

## §1 现量

### 1.1 被审的那把尺在哪枚文件哪几行、它到底数什么

尺本体＝`internal/tools/task_cancel_221_legs_test.go:223` 的
`Test221DeferredMarkerForCancelLiftedButListStillMarked`（`grep -n 'func Test221DeferredMarkerForCancelLiftedButListStillMarked'`
现命中行号 **223**）。逐行读它数什么（原文抄自 `logs/r3c/gates/s1-ruler-location.txt` 第 1 段）：

| 行 | 原文 | 在数什么 |
|---|---|---|
| `:224` | `src, err := os.ReadFile("task.go")` | **读物理盘**——§1.5 那条仪器事实的主体 |
| `:230` | `if !strings.Contains(line, "DEFERRED") { continue }` | 先把不含 `DEFERRED` 的行整行丢掉 |
| `:233-235` | `if strings.Contains(line, "task.cancel") { cancelMarked++ }` | **(a) 支**：同时含两枚词面的行数 |
| `:236-238` | `if strings.Contains(line, "task.list") { listMarked++ }` | **(b) 支**：同上，纯词面共现 |
| `:240`／`:244` | `cancelMarked != 0` 报红／`listMarked == 0` 报红 | (b) 支的断言只问「行数是否掉到 0」 |

⇒ **(b) 支数的不是能力、是"同时含 `DEFERRED` 与 `task.list` 的行数 >0"**。同文件另一枚相关尺（本腿**原样保留**）：
`:170 Test221TaskCancelIsRegisteredAtItsFrozenLevel` 里的 `:183 if len(names) != 2` ⇒ **计数形**，票面 §四 预言它会误判
"合法新增第三枚 task 工具"（假红）。本腿没动它，并在 §3 发② 实测到它的红句。

### 1.2 词面被叙述句撑住的机制（现跑枚数）

```
$ grep -c "DEFERRED" internal/tools/task.go                        -> 3
$ grep "DEFERRED" internal/tools/task.go | grep -c "task\.list"    -> 2
$ grep "DEFERRED" internal/tools/task.go | grep -c "task\.cancel"  -> 0
```

三行逐字（现读行号）：`:23` `//	task.list    -   DEFERRED with five fields, PLAN.md §7 :1531`（**就是那一行标记**）、
`:32` `// and the DEFERRED line above is the only one left standing.`（讲这张名册的叙述句，含 `DEFERRED`、不含 `task.list`）、
`:279` `// diagnostics, not as a tool surface: task.list is DEFERRED (§7 :1531), and`（`Count()` 的文档注释，**同时含两枚词面**）。

⇒ **摘掉 `:23` 之后 (b) 支数到的仍是 1**（同一把尺打在合成副本上现量：`2` 掉到 `1`，见
`logs/r3c/mut/copies-and-landing-proofs.txt` 末段），`listMarked == 0` 永不成立。
票面"现量"表第 2 行＋§四"词面被叙述句撑住"的机制**在本腿现量下复认**，撑住它的那一枚＝`:279`（叙述句，不是名册行）。

### 1.3 `BuiltinTaskEntries` 名册枚数与 `task.list` 在不在（现跑）

产码逐字（`internal/tools/task.go:600-605`，见 `logs/r3c/gates/s1-ruler-location.txt` 第 3 段）：

```go
func BuiltinTaskEntries(d TaskDeps) []Entry {
	return []Entry{
		{Tool: taskOutput{d: d}, Decl: TaskOutputDecl()},
		{Tool: taskCancel{d: d}, Decl: taskCancelDecl()},
	}
}
```

- 名册枚数＝**2**；两名实测＝`task.output` / `task.cancel`（本腿用例 `t.Logf` 原文：
  `现名册：task.output / task.cancel（DEFERRED 的判据读的是这一行，不是注释里的词面）`）。
- **`task.list` 不在这两名里** ⇒「task.list 未注册」今天**为真**。这是事实陈述；⛔ 本腿不裁它该不该注册
  （D34／`PLAN.md §7 :1531` 那一格与产码腿的地界）。
- ★**这枚钉的分母不是测试专用件**（能力类断言必问生产调用者）：`cmd/wisp/run.go:544`
  `for _, e := range tools.BuiltinTaskEntries(tools.TaskDeps{Roster: rt.tasks, Paths: paths})` ⇒
  本腿读的是**装配根注册的那一枚函数**，不是测试自造的列表（全仓 `BuiltinTaskEntries` 命中名册见
  `logs/r3c/gates/s1b-existing-capability-rulers.txt`）。
- 分母口径＝`BuiltinTaskEntries(TaskDeps{Roster: NewTaskRoster()})` ⇒ 读的是**编译期构造出的名册**，
  所以 `go test -overlay` 量得到它（与 §1.5 那把读盘尺正好相反，这一对差异就是本格的落点）。

### 1.4 与票 225 的分界（两行逐字引；⛔ 不做 225 的活）

- 票 236 AC#2 原文（`.scratch/wisp/issues/236-six-cells-that-only-surface-at-the-reading-layer.md:27`）里那一句逐字：
  「⚠ 与票 225 分开算账：225 管"标记与 `SPEC-12 §5` 双向对账"，本格只管**这把尺本身有没有牙**」。
- 票 225 标题（`.scratch/wisp/issues/225-deferred-markers-do-not-match-the-spec12-registry-and-nothing-checks-it.md:1`）逐字：
  「# 225 — **`DEFERRED(...)` 代码标记与 `SPEC-12 §5` 登记表今天对不上，而且没有任何仪器在查这一条**：`AGENTS.md` §1.1 那条"1:1 双向"目前是**只靠人**的规矩」。
- **分界落在实处、不只是嘴上**：本件的判据（§3 四发＋配对、§1.6 三把尺）**零枚读 `docs/specs/SPEC-12.md §5`**，
  零枚断"标记与登记表对不对得上"。本腿唯一碰"标记"的地方＝§1.2 那三行的**枚数现量**（那是 AC#2 要的"这把尺数什么"，不是对账）。
  同类形状只登记、不核对：`grep -rn 'DEFERRED(' internal/tools --include=*.go`（排 `_test.go`）＝**2 行**
  （`internal/tools/tool.go:147`、`internal/tools/open.go:101`）⇒ 那两枚属 225 的射程，⛔ 本腿一枚不碰。

### 1.5 票面点名要写进判据注释的那条仪器事实（本腿把它变成可量的）

- 事实原文（票 236 §四 `:79`，编排者现读复认）逐字：「尺本体 `task_cancel_221_legs_test.go:223-247`：`:224` 是
  **`os.ReadFile("task.go")`**——**读盘型尺对 `go test -overlay` 结构性不可见**（overlay 只替换编译器看到的字节，
  测试二进制运行时打的是物理盘）。⇒ 今天**连"测这把尺的牙"都做不到**，不是没测、是测不出。这条必须写进判据注释。」
  同一句在票面 AC#2 那一格里也写着（`:27` 末）：「测它的牙必须 overlay 替换测试文件里的读路径并配正控」。
- 注释里在位（`236-r3` 落的字，本腿复认逐字在位）：`internal/tools/tasklist_deferred_236r3_teeth_test.go:33-45`。
- **本腿把这条事实变成可量的两半**：
  1. **编译期那一半**＝同文件 `:86` `//go:embed task.go`（`//go:embed` 走**构建期**，overlay **够得到**）。
     用例 `Test236R3InstrumentFactReadDiskIsBlindToOverlayCarriesIt` 把两把读数并排比（盘 vs 编译）：
     HEAD 上两数相等、只有 `-overlay` 能让它们分开 ⇒ 这一枚红＝"这一发的字节真进了编译器"的**落地证明**。
  2. **正控那一半**＝overlay 替换**测试文件里的读路径**：只改 `:224` 那一行 `os.ReadFile` 的字符串参数，
     指进 `D:/tmp/wisp236r3c/` 的合成副本（副本除那一行外与盘上逐字节相同）。⇒ 这一发证明 overlay 打得进 `_test.go`，
     于是"读盘尺一直绿"就**排除了"overlay 没落地"这一解释**。读数在 §3 配对①②③。

### 1.6 每把负向尺先证命中得了真名（今天栽过两次，逐枚自证）

本件只有一处负向结论：**「钉『`task.list` 这一行不许出现在注册名册里』那一维，今天零载具」**。尺＋命中证明
（原始输出＝`logs/r3c/gates/n1-n5-negative-rulers.txt`）：

| 尺（本腿现跑） | 读数 | 命中真名的证明（同形换真名必须非零） |
|---|---|---|
| `grep -rn '"task\.list"' internal/tools --include=*_test.go` | **4 行**；逐行读：4 行全落在注释或"数词面"的谓词里（`task_cancel_221_legs_test.go:236`、本腿文件 `:108`／`:144`），**零枚比的是注册名** | —（这一把是"命中了，但没有一枚是能力尺"，不是零命中） |
| `grep -rn '== "task\.list"' internal --include=*_test.go` | **1 行**＝本腿自己那枚钉子（`tasklist_deferred_236r3_teeth_test.go:144`，`21b53d52` 版行号） | 同族形状换真名：`grep -rn 'Name() == "task\.cancel"' internal --include=*_test.go` ⇒ **命中 `task_cancel_221_legs_test.go:176`（非零）** ⇒ 这把尺抓得到英文注册名，**不是中文词面尺** |
| `grep -rn 'registered\["task\.' . --include=*.go` | **0 命中** | ★**这一把的零命中本身是伪尺**：名集在 `internal/tools/task_cancel_221_test.go:94` 以**计算键**写入（`registered[e.Tool.Name()] = true`），任何"字面量键"的 grep 结构上都打不到它。⇒ 本腿**不拿这一把下结论**，只拿前两把（各有命中证明） |

⇒ 结论的精确形状：**不是"仓里没有能力尺"**——票面 §四 指名的两枚本腿复认其形状
（`internal/tools/task_cancel_221_test.go:42 allBuiltinEntriesHere` ＋ `:89 Test221EveryPromisedTaskNameIsRegistered`；
`internal/tools/ticket175r2_stamp_live_test.go:349` 的分类表）。
**是"没有任何一枚断言问成名册里不许有 `task.list`"**（`221EveryPromisedTaskNameIsRegistered` 钉的是「说明文字许诺过的必须注册」＝正向；
`175r2` 钉的是「注册了的必须有盖戳分类」＝另一维）。⇒ AC#2 的落点＝在同一枚**编译期分母**上**新增**一枚具名负向断言；
⛔ 不删、不放宽任何既有断言（词面尺两半都留着、计数尺照数它的 2）。

---

## §2 判据表（命令｜原始读数｜判语）

读数逐字来自 `logs/r3c/**`。基线口径见 §0；每一发的完整红名册在 §3。

| # | 命令 | 原始读数 | 判语 |
|---|---|---|---|
| A | `PATH=… go test ./internal/tools/ -count=1` | `ok  	github.com/CarlosShao/wisp/internal/tools	15.318s`，rc=0 | 未突变态全绿 ⇒ §3 的红只能归因于那一行突变（本包没有别人的红可混，见 §0 的 206／0／0） |
| B | 同上 `-v` | 顶层 **PASS=206／FAIL=0／SKIP=0**、`=== RUN` 281、`ok 13.924s` | 票面 §九 那句"并发＝互相洗读数"这一轮**没发生**：红名册可直接用 |
| C | `go vet ./internal/tools/` | **rc=0**（起手一轮、最后一次编辑之后再一轮） | 静态门干净，载体编译得动（`d9aff5f` 已报过一次，本腿复跑复认） |
| D | `sh scripts/d22scan.sh` | 正控先绿 `runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0` ⇒ 真扫末行 `d22scan: clean - no D22 ban violations`，**rc=0**；`ban #8 internal/ = 514` | 门的射程覆盖本腿那枚文件（`514` ＝ `236-r2` 报的 `513` 恰多这一枚 ⇒ 在射程内且没撞禁令）。详读数在 §4 |
| **E** | **正控**：`go test ./internal/tools/ -count=1 -overlay D:/tmp/wisp236r3c/overlay/ctrl.json`（副本与盘上**逐字节相同**，md5 同为 `4138177e…`） | rc=**0**、`ok  	github.com/CarlosShao/wisp/internal/tools	15.116s`（pass2 `14.378s`）、红名册**空**；`-v` 隔离三枚＝**PASS／PASS／PASS**；成对读数原文 `成对读数：盘上 2 行 / 编译期 2 行（HEAD 上两数相等＝这把尺没在空转）` | `-overlay` 这个动作本身不产红 ⇒ 后面每一发的红都是"那一行文字差"造成的。**同时证明这把尺不空转**：它交出两名名册＋两把相等的标记读数 |
| F | `go test ./internal/tools/ -count=1 -run 'Test221DeferredMarkerForCancelLiftedButListStillMarked\|Test236R3' -v`（未突变盘上态） | 三枚 **PASS**、`ok  	github.com/CarlosShao/wisp/internal/tools	0.044s`（`gates/final-gates-second-pass.txt`） | 交件态**全绿**：本腿没留下任何"永久红"，也没靠放宽断言换绿 |
| G | 三把负向尺（§1.6） | 见该表 | "这一维今天零载具"**成立**，且尺命中得了真名 ⇒ 不是"读不到所以以为没有" |
| H | 名册枚数与标记枚数（§1.2／§1.3 的命令） | 名册 2 枚（`task.output`／`task.cancel`，无 `task.list`）；`DEFERRED` 行 3、其中共现 `task.list` 的 2 | 本格判据的**现量基线**：「能力今天为真」与「词面尺数的是共现行数」两件事各自都有数 |
| **I** | **未修码读数**：`-overlay teeth-m13.json`（只摘 `task.go:23` 那一行标记） | rc=1，红名册**恰 1 枚**＝`--- FAIL: Test236R3InstrumentFactReadDiskIsBlindToOverlayCarriesIt (0.00s)`；归档词面尺 **PASS**、能力钉 **PASS**；副本上共现行数 **2 掉到 1** | ★**票面要的"它今天无牙"就是这一行读数**：标记被摘到只剩叙述句撑着的那一行，数词面那把尺照样绿，且这是 overlay 已落地之后的绿（同一发里编译器读到 1 行）。⇒ 无牙＝**机制性**（§1.2 的算术），不是"没测到" |
| **J** | **装牙后的读数**：`-overlay m2b.json`（把 `task.list` 注册进名册） | rc=1，红名册 4 枚，含 `--- FAIL: Test236R3TaskListDeferralIsPinnedByTheRosterNotByAComment (0.00s)`，红句自报名册三名 `[task.output task.cancel task.list]` | ★**同一发上旧尺仍全绿、新尺具名红** ⇒ 这就是"改扫能力"买到的那颗牙（对比 I 行：两发的旧尺都是绿的，差别只在**能力有没有变**） |

### 本格判据一句话（判语靠 §3 的发数撑）

**「`task.list` 是 DEFERRED」这一格从此由编译期名册负责，不再由注释负责**：名册里出现 `task.list` 那一行 ⇒ 具名红（发②／I·J 对比）；
名册少掉已接线的 `task.cancel` ⇒ 具名红（发③）；只摘注释里那行标记、能力不动 ⇒ 这枚钉**不红**（发④，这是它对的行为），
而"数词面"那把旧尺在**同一发上照样绿**、并且**在它真读得到字节的那一发上也照样绿**（配对①）⇒
票面要的"先证它今天无牙"与"改扫能力"两颗牙，两边都有原文读数。

---

## §3 突变名册（★四发＋三发配对全部由 `236-r3c` 跑；⛔ 未引用 `236-r3` 的 54 枚台件当凭据）

统一做法：合成副本放 `D:/tmp/wisp236r3c/<发名>/task.go`，替换 json 放 `D:/tmp/wisp236r3c/overlay/<发名>.json`，
命令＝`go test ./internal/tools/ -count=1 [-v -run '<三枚>'] -overlay <json>`。七份 json 已并档进 `logs/r3c/overlay/`。
⛔ **全程未在共享工作树里原地编辑过任何产码**：七发跑完再验一次
`md5sum internal/tools/task.go`＝`4138177e29ffff427776a73eb0d61a9b`＝`git show HEAD:internal/tools/task.go | md5sum`，
`internal/tools/task_cancel_221_legs_test.go` 同样盘上＝HEAD＝`13a0ba39e30fb95ad4f2be88beed4fc7`
（逐字在 `logs/r3c/mut/rosters-and-restore.txt`）。⇒ 本腿的"还原"**不是一条命令**，是 overlay 的结构性事实（它从不写盘）；
`git status --porcelain -- internal cmd` 四发之后仍 **0 行**。

### 发① 正控 `ctrl`

- 改什么：什么都不改（`cp` 一份，逐字节相同）。overlay＝`D:/tmp/wisp236r3c/overlay/ctrl.json`。
- 落地证明：`diff` 零输出；副本 md5 `4138177e29ffff427776a73eb0d61a9b`＝盘上＝HEAD blob；跑完再量盘上仍是这一串。
- 读数：全量包 rc=**0**、红名册**空**（`logs/r3c/mut/full-ctrl.log`，`ok 15.116s`／pass2 `ok 14.378s`）；
  `-v` 隔离三枚逐字 **PASS／PASS／PASS**（`logs/r3c/mut/isolated-v-lines.txt` `ctrl` 段）。
- 判语：**该绿的绿**。这一发是其余各发的分母——没有它，后面任何"绿"都不能归因。

### 发② `m-2b`＝把 `task.list` 注册进名册

- 改什么：`task.go:603` 之后插一行 `{Tool: taskListRow236r3c{taskCancel{d: d}}, Decl: taskCancelDecl()},`，
  末尾追加 `type taskListRow236r3c struct{ taskCancel }` ＋ `func (taskListRow236r3c) Name() string { return "task.list" }`。
  ⇒ 名册 2 枚变 3 枚、第三名正是被 DEFERRED 那一行；**注释一个字没动**（词面尺的输入不变）。
  ★形状说明：内嵌 `taskCancel` 只为借够 `Tool` 接口另两枚方法，`Name()` 覆盖成 `task.list`；⛔ 未新造工具面、未动产码、未进任何生产装配。
- overlay＝`m2b.json`。落地证明：`diff` 恰三处（`603a604`、`889a891,896`）＋副本 md5 `e51a9a9a59f1e97f8421f8fa368a08c5`（盘上仍 `4138177e…`）；
  **行为层落地证明**＝红句里名册自己报出第三名，不靠 diff 猜。
- 红句（本腿那枚钉，`-v` 隔离原文）：
  `tasklist_deferred_236r3_teeth_test.go:132: DEFERRED 那一支和名册分叉了：BuiltinTaskEntries 的注册名册里出现了 task.list 这一行（现名册：task.output / task.cancel / task.list）。§7 :1531 与 D34 都还把它记成 DEFERRED，「在册＋无实现＋无人认领」正是票 164 AC#1 要杀的那一形；真要解冻得由人工批准、标记与接线同批核销（A434 批准范围 item 7），而不是留下一行注释替一枚已注册的工具说它还不可用。`
  ＋ `:155: 现名册：task.output / task.cancel / task.list（DEFERRED 的判据读的是这一行，不是注释里的词面）`。
- 全量包红名册：rc=**1**、**恰 4 枚**（`logs/r3c/mut/full-m2b.log` 逐字 `^--- FAIL`）＝
  `Test221TaskCancelIsRegisteredAtItsFrozenLevel`／`Test236R3TaskListDeferralIsPinnedByTheRosterNotByAComment`／
  `Test236R3InstrumentFactReadDiskIsBlindToOverlayCarriesIt`／`TestEveryRegisteredToolIsClassifiedForMarking175r2`。
- **判语**：本格要的对比就在这一发——**能力变了、注释没变**时，读注释那把尺
  （`Test221DeferredMarkerForCancelLiftedButListStillMarked`）**逐字不在红名册里**（它看不见），
  **读名册这枚钉具名红** ⇒ (b) 支改扫能力之后，这个方向第一次有牙。
  另两枚红是**别人尺的正当反应**、本腿不裁：计数尺红句
  `task_cancel_221_legs_test.go:184: task 家族注册了 3 枚：[task.output task.cancel task.list]（task.list 仍须是 DEFERRED，本票不许顺手注册它）`；
  `175r2` 红句 `ticket175r2_stamp_live_test.go:367: task.list 既不在盖戳名册里、也没有被分类：它返回的内容是不是外部内容？…`
  （属票 175 的盖戳普查射程，与 AC#2 无关）。

### 发③ `m-2e`＝注销 `cancel`（名册少一枚）

- 改什么：删 `task.go:603` 那一行 `{Tool: taskCancel{d: d}, Decl: taskCancelDecl()},`（`sed '603d'`，只这一处）。
- overlay＝`m2e.json`。落地证明：`diff` 恰 `603d602` 一处＋副本 md5 `f2152339178d911b1d7944748aa9773c`（盘上仍 `4138177e…`）。
- 红句（本腿那枚钉，`-v` 隔离两行）：
  `tasklist_deferred_236r3_teeth_test.go:151: task 家族名册 = task.output, want 含 task.output 与 task.cancel（task.output 是票 164 的，task.cancel 是票 221 甲形接的线）——名册少一枚就是「说明书与能力分叉」的另一形`
  ＋ `:155: 现名册：task.output（…）`。
- 全量包红名册：rc=**1**、**10 枚**（`logs/r3c/mut/full-m2e.log` 逐字）＝
  `Test236R2TaskCancelRefusesWhenHostGaveNoCallerID`／`Test236R2TaskCancelRefusesWhenRosterIsUnwired`／
  `Test221TaskCancelIsRegisteredAtItsFrozenLevel`／`Test221ParentStopsItsOwnChildRowAndStreamSettle`／
  `Test221SubagentCannotStopSiblingOrItself`／`Test221ParentCancellationStillDoesNotCascade`／
  `Test221TaskCancelUnderNoGateStopsAtTheWindow`／`Test221EveryPromisedTaskNameIsRegistered`／
  `Test236R3TaskListDeferralIsPinnedByTheRosterNotByAComment`／`Test236R3InstrumentFactReadDiskIsBlindToOverlayCarriesIt`。
- **判语**：**反向也有牙**（名册少一行同样具名红）⇒ 这不是一枚只盯一个名字的单向尺。
  本发另给出第三条只有读数能给的判语：**注销一支＝10 枚具名红**（`236-r2` 那格装上的五枚拒绝尺全在里面）
  ⇒ "名册这一层动一行的爆炸半径"今天是**实测到的数**。本腿**没有**因为红得多去收窄任何断言（禁区）；
  ⛔ "该红几枚才算对"不是 AC#2 的射程，见 §5 第 3 条。

### 发④ `teeth-m13`＝摘掉 DEFERRED 标记行（票面点名那一形：先证它今天无牙）

- 改什么：**只**删 `task.go:23` 那一行 `//	task.list    -   DEFERRED with five fields, PLAN.md §7 :1531`（`sed '23d'`）。
  名册两行注册、`Count()` 注释全不动 ⇒ 这是"标记与能力分叉"里**只动标记**那一支。
- overlay＝`teeth-m13.json`。落地证明：`diff` 恰 `23d22` 一处＋副本 md5 `bb011a71b1949f23283191420825fb5a`（盘上仍 `4138177e…`）。
- **副本上的词面算术（"无牙"的机制证明，`logs/r3c/mut/copies-and-landing-proofs.txt` 末段）**：
  `grep "DEFERRED" <副本> | grep -c "task\.list"` 由 **2 降到 1** ⇒ (b) 支那句 `listMarked == 0` **永远碰不到 0**。
- 全量包读数：rc=**1**、红名册**恰 1 枚**＝`--- FAIL: Test236R3InstrumentFactReadDiskIsBlindToOverlayCarriesIt (0.00s)`。
  `-v` 隔离三枚逐字＝归档读盘尺 **PASS**／能力钉 **PASS**／仪器尺 **FAIL**。仪器尺红句（原文两行）：
  `:189: 编译期与运行期读到的 task.go 字节不一致（编译 37326 字节 / 盘上 37391 字节，标记行数 1 对 2）：只有 -overlay 能让它们分开。这一发就是那条仪器事实的落地证明：盘上那把尺看不见突变，所以它绿的时候不能当作「标记还在」的证据。…`
  ＋ `:195: 同一枚 DEFERRED＋task.list 标记行数，编译器看到 1 行、os.ReadFile 看到 2 行：两个数只有在 -overlay 下才会分开（交件 §3 M-2a 的配对读数）`。
- **判语（本发是核心，三句都钉死）**：
  1. **票面要的"未修码读数"坐实＝这把尺今天无牙**：摘掉标记行，数词面那把尺**逐字 PASS**（机制见 §1.2、算术见上一段），
     而且这个绿是 overlay **已落地之后**的绿（`编译器看到 1 行` 那一行就是落地证明，⛔ 不是"overlay 没进去"）。
  2. **能力钉在这一发 PASS 是对的**：能力一个字没变，它不该红 ⇒ 它不是"更宽的词面尺"，是**换了维度**（注释归注释、名册归名册）。
  3. 这一发真正的后果留给后程一个明确形状：**"注释说谎"这件事今天仍拦不住**——只有仪器尺因为 overlay 在跑而红，
     没有任何一行为"标记不许被摘"报警。⇒ 见 §5 第 4 条（本腿**故意不造**那枚常驻尺，并写明为什么）。

### 配对①②③＝票面点名要的"overlay 替换测试文件里的读路径并配正控"

| 配对 | 改什么 | 归档读盘尺 | 仪器尺（`//go:embed` vs `os.ReadFile`） | 判语 |
|---|---|---|---|---|
| ① `readpath-teeth-m13` | 同时 overlay 两枚文件：`task_cancel_221_legs_test.go` 的 `:224` 读路径改指 teeth-m13 副本，且 `task.go` 也 overlay 成 teeth-m13 | **PASS** | **FAIL**（`:199: 成对读数：盘上 2 行 / 编译期 1 行`） | 尺的谓词**命中得了真名**：它读的字节里 `:279` 那句还留着，数到 1 就绿。⇒ 它的绿不是 bug，是词面尺的本性；同一发里编译器读到 1 行 ⇒ overlay 确实打得进 `_test.go` |
| ② `readpath-no-markers` | 同上，但读路径与 `task.go` 都指向"两枚标记全摘"的副本（`no-markers/task.go`，md5 `474ca494c266dc345ed786cda81eddd4`，`grep DEFERRED｜grep task.list`＝**0**） | **FAIL**（`:245: task.list 的 DEFERRED 标记不见了：本票不许摘它（那半支仍然没有实现）`） | **FAIL**（`:199: 成对读数：盘上 2 行 / 编译期 0 行`） | ★**这一发把"测这把尺的牙"从做不到变成做得到**：改读路径之后它**有牙**（词面真没了它会红）⇒ 票面说的"无牙"精确含义＝**对 overlay 突变无牙（结构性不可见）**，不是"谓词写坏了"。也要摘到 0 行它才红 ⇒ 叙述句 `:279` 撑着 ⇒ 与 §1.2 的算术互咬 |
| ③ `readpath-nomarkers-testonly` | **只** overlay 那枚测试文件（读路径→无标记副本），`task.go` 一字不 overlay | **FAIL**（同上红句） | **PASS**（两把读数相等） | 反向对照：尺红的时候**盘上标记仍在**（`盘上 2 行`）⇒ 这一对读数把"尺测的是它读到什么"与"盘上有什么"这两件事彻底分开。⛔ 本腿**没有**因此去改读盘尺的断言，也没把 `os.ReadFile` 改成"能看见 overlay"（那条路不存在，注释里也写了这条） |

### 前腿报过的形状，本腿怎么处理

`d9aff5f` 的提交说明写「主用例在 **m-2b 与 m-2e 两发 FAIL**，在 m-2a／m-2c／m-2z／fix-green 等发 PASS」。
★本腿**既不采信也不推翻**，只把自己四发的答案写在上面，并说明为什么形状自洽：r3 的"主用例"＝能力钉，
它只在**名册变了**时红（发②③），在**只有注释变了**时绿（发④）——这正是 AC#2 要的行为（能力尺不替注释负责），不是遗漏。
**"该红却绿"（尺无牙）的例子：四发里一例也没有出现**——除了发④那一枚**故意要它绿**的能力钉，
唯一该红的读盘尺在配对②里确实红了。⛔ 本腿**没有跑** `m-2c`／`m-2z`（第四发帽给的是票面点名的四形），见 §5 第 2 条。

---

## §4 门禁

本腿全部由 `236-r3c` 现跑，逐字落在 `logs/r3c/gates/`（`vet-fmt.txt` 起手一轮、`final-gates.txt` 与
`final-gates-second-pass.txt` 在最后一次注释编辑之后再两轮；两轮读数一致）。

| 步 | 命令 | 读数 |
|---|---|---|
| 扫描门 | `sh scripts/d22scan.sh` | 正控先绿：`runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0`；真扫末行 `d22scan: clean - no D22 ban violations`，**rc=0**。分母逐枚在案：`bans #1-5 internal/=228`、`cmd/=38`、`#6 frontend/=85`、`#7 internal/tools/=23`、`#8 design/=39`、`#8 frontend/=85`、`#8 internal/=514`、`#8 cmd/=104` ⇒ `514`＝`236-r2` 报的 `513` 恰多本腿那枚文件，**它在门的射程里且没撞禁令** |
| 格式 | `gofmt -l internal/tools/tasklist_deferred_236r3_teeth_test.go` | **零输出**，rc=0（改前改后各一轮） |
| 格式 | `"$(go env GOPATH)/bin/gofumpt.exe" -l <同一枚>`（`v0.12.0 (go1.27.1)`） | **零输出**，rc=0 |
| 静态 | `go vet ./internal/tools/` | rc=0（两轮） |
| 未突变态 | `go test ./internal/tools/ -count=1 -run 'Test221DeferredMarker\|Test236R3' -v` | 三枚 **PASS**、`ok 0.044s` ⇒ 交件态全绿 |
| 工作树 | `git status --porcelain -- internal cmd` | 起手 **0 行** → 四发之后 **0 行** → 本腿注释编辑后只有 ` M internal/tools/tasklist_deferred_236r3_teeth_test.go` 一行 → commit 之后回 **0 行** |
| 工作树 | `git status --porcelain -- scripts .github` | **0 行**（本腿零动 CI／脚本，禁区合规） |

⛔ 本腿未把任何 CI 步改成 `continue-on-error`、未加 `if:`、未放宽／未删除任何既有断言、未新增导出名
（文件里全部是非导出符号，`git diff --numstat` 的删除列＝8，全是注释段落重排，⛔ 无一枚落在别人文件上）。
★整包 `go test ./internal/...` **未重跑**＝派单的收缩决定（那是 AC#1 那格的判据）。
⚠ 一处过期提醒（本腿不复述、只指名）：`logs/gates/gofumpt-internal-cmd.txt`（r3 留的台件）里那 5 枚 `internal\...` 脏名册
是**全目录分母**的读数，属 AC#6 那格，⛔ 不是本腿的门禁项。

---

## §5 判不动的地方（每条给"为什么"，⛔ 不空白、不写"没有"了事）

1. **(a) 支（`cancelMarked` 那一支）今天有没有牙——本格判不了。**
   `:240` 断的是 `cancelMarked != 0`（"不许复活 `task.cancel` 的标记"）。要证它有牙，得造一发"往文件里**写进**一行同时含
   `DEFERRED` 与 `task.cancel` 的注释"的突变；而**同一枚载体是 `os.ReadFile`**（§1.5），overlay 打不进去 ⇒
   只能像配对①那样**改读路径**才能测，那是**又一发配对**。本腿帽在四发，名额已给票面点名的 teeth-m13／m-2b／m-2e／正控。
   ⇒ (a) 支**留在判不动**，⛔ 不写"顺带也行了"。
2. **`m-2c`（注册一枚与 `task.list` 无关的新 `task.*`）未实测**，票面 §四 的假红预测（计数形 `len(names) != 2` 会误判合法新增第三枚）
   本腿**只有源码级论证**：那一行断言的字面是枚数，名册 2→3 必红（发② 实测到了 3 枚的形态）；能力钉的字面是**名字集合**，
   多一枚别的名字不动它。★前腿的 `m-2c` 台件在盘上（`logs/mut/full-package-m-2c.txt`），那是**别人的读数**，
   本腿不拿它当自己的凭据 ⇒ 这一维**要结就得再派一发**，本腿在测试注释里也如实写了"NOT measured by r3c"。
3. **爆炸半径该定在几枚，不归本格裁。** 发③（注销一行）＝10 枚具名红，发②（多注册一行）＝4 枚。
   "动名册一行＝十个具名用例响"是好事（行为尺密），但 AC#2 只管**这把尺自己有没有牙**。
   本腿因此**没有**因"红得多"收窄任何断言，也**没有**把它写成"钉太宽"的缺陷；若要裁"该红几枚"，请另立一格。
4. **本格拦不住"注释与代码分叉"，只拦得住"名册与能力分叉"——而且这是**故意**的。**
   发④ 的真实后果：产码注释可以说谎而**零枚行为尺报警**（能力钉正确地绿，仪器尺红的是"overlay 在跑"这件事）。
   要在**未突变态**也拦住这一形，唯一形状是再加一枚词面尺（"文件里必须留着 `DEFERRED`＋`task.list` 那一行"）。
   ⛔ 本腿**不造**它，两条理由：①它要求"注释必须含某一行文字"，而这条承诺**没有任何契约凭据**
   （票面 §四 逐字只写「保留计数、另加名集」，没写"再加一条词面断言"）；②它是**新增一条约束文档措辞的断言**＝改契约形状＝要人工批准，
   而派单写死"agent 单方面改契约＝跑歪"。⇒ 风险**登记在案**，⛔ 不许把 AC#2 读成"这一格修好了"。
5. **同一把尺两遍读数不同：没发生，但有一处口径必须说清否则后人会误读。**
   本腿四发在**两种文件字节**上各跑一遍，红名册逐字相同（`mut/pass2-readings.txt`），基线两次也一致。
   ★要说清的口径＝§0 那个 `PASS=206` 是**顶层枚数**（`grep -c '^--- PASS'`），票面 §六 逐字警告过行首锚定尺**看不见缩进子测试**；
   本腿判据全靠"具名红名册"与逐行原文，不依赖名级计数，所以不与此冲突；但**引用 206 必须带"顶层"两字**。
6. **仪器尺那枚红对"良性 overlay"也红**（注释里也写了）：任何一次对 `task.go` 的 `-overlay`（哪怕字节完全相同以外的任何改动）
   都会让 `Test236R3InstrumentFactReadDiskIsBlindToOverlayCarriesIt` 红。⛔ 这不是新常红（HEAD 与 CI 上两把读数相等＝绿，
   且全仓 `scripts/**`／`.github/**` 现查 **0 处**使用 `-overlay`），⛔ 也不许"修"成少比几样——少比正是旧尺量不到的原因。
   谁下一格要用 overlay 改 `task.go`，**会先撞上这一枚红**，这是设计行为，已就地写明。
7. **本格不答的清单**：AC#1／AC#1b（五枚拒绝分支，`236-r2` 已交、另有终裁腿）、AC#3（task id 由谁铸造）、AC#4（`ci.yml` 注释）、
   AC#5（红名册口径）、AC#6（格式门分母）、票 225（标记与 `SPEC-12 §5` 双向对账）、`task.list` 到底该不该注册（D34／`PLAN.md §7 :1531`）。
   ⛔ "AC#2 算不算修好"由编排者凭本件＋终裁腿判，本件不勾任何框、不写"完成"。

---

## §6 复现命令（后程逐条可重打）

```sh
# 基线（单包，非整包）
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./internal/tools/ -count=1
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./internal/tools/ -count=1 -v   # 206/0/0

# 四发（合成副本在 D:/tmp/wisp236r3c/，json 在 D:/tmp/wisp236r3c/overlay/，已并档 logs/r3c/overlay/）
for m in ctrl teeth-m13 m2b m2e; do
  go test ./internal/tools/ -count=1 -overlay "D:/tmp/wisp236r3c/overlay/$m.json"
  go test ./internal/tools/ -count=1 -v -run 'Test221DeferredMarkerForCancelLiftedButListStillMarked|Test236R3' \
      -overlay "D:/tmp/wisp236r3c/overlay/$m.json"
done
# 三发配对（读路径替换）
for m in readpath-teeth-m13 readpath-no-markers readpath-nomarkers-testonly; do
  go test ./internal/tools/ -count=1 -v -run 'Test221DeferredMarkerForCancelLiftedButListStillMarked|Test236R3' \
      -overlay "D:/tmp/wisp236r3c/overlay/$m.json"
done
# 门禁
sh scripts/d22scan.sh; gofmt -l internal/tools/tasklist_deferred_236r3_teeth_test.go
"$(go env GOPATH)/bin/gofumpt.exe" -l internal/tools/tasklist_deferred_236r3_teeth_test.go
git status --porcelain -- internal cmd; git status --porcelain -- scripts .github
# 还原证明（overlay 不写盘，跑完仍是 HEAD 字节）
md5sum internal/tools/task.go; git show HEAD:internal/tools/task.go | md5sum
```

---

## 附录（逐字保留）：`236-r3` 留下的 36 行骨架原文

下面这份是本文件被本腿续写**之前**的全部内容（腿 `236-r3`，`d9aff5f` 代提）。它的每一节都写着"打算答"——
那是**打算**、不是判语，所以本件不引用它任何一句当结论（⛔ 临时件不删，只在文末存档）。
它的正文标题是「236-r3 — AC#2 的交件（六节）」，从下一行起到本节结束都属于 `236-r3`；
本腿另存了一份逐字节拷贝在 `.scratch/wisp/probes/236/r3/evidence-skeleton-236r3.md`（md5 `d2fc7e8ef2694cb82bc42405aaf4eb55`）。

下面把 36 行骨架**逐字内联**（内容与 `evidence-skeleton-236r3.md` 那份留档逐字节相同；
两处若将来漂移，以本文件为准并把漂移具名上报，不许静默择一）。

---

# 236-r3 — AC#2 的交件（六节）

腿＝`236-r3`；射程＝票 236 的 **AC#2 一格**（那把 DEFERRED 尺的 (b) 支改扫能力）。
⛔ 本件不答也不裁 AC#1／AC#1b／AC#3／AC#4／AC#5／AC#6（派单 §0 第 4 把尺），七枚框一枚不碰。

> **本节以下目前是骨架**：每节标题＋本腿打算答什么。读数逐节回填，占位词一个不留（写完才 commit）。

---

## §0 起手锚（四把尺原文）

打算答：`sed -n '3,5p'` 的撤票口令原文与 grep 计数、`git log -1`／`date`、`git status --porcelain -- internal cmd`（★必须证明起手 `internal/tools` 无未提交脏文件）、`grep -n '^- \[ \]'` 的逐条行号与枚数、`git status --porcelain -- scripts .github`。

## §1 现量（AC#2 那把尺的三点量＋分界）

打算答：①尺本体在哪枚文件哪几行、它到底在数什么（逐行读）；②`BuiltinTaskEntries` 注册名册今天有几枚、`task.list` 在不在（现跑）；③与票 225 的分界（逐字引票面＋225 标题，⛔ 不做 225 的活）；④词面被叙述句撑住的机制复认（`grep -n DEFERRED internal/tools/task.go` 行数、同时含 `task.list` 的行数）；⑤"零命中不是证据"——每把负向尺先证它命中得了真名。

## §2 判据表（跑了什么命令｜原始读数｜判语）

打算答：基线绿名册（用例枚数＋rc）、AC#2 判据逐条对应哪一枚用例、每把尺的口径。

## §3 突变名册（teeth 名册）

打算答：
- **M-2a**＝摘掉 `internal/tools/task.go:23` 那行 DEFERRED 标记（`teeth-m13` 那形）：**未修码读数**＝归档 leg `Test221DeferredMarkerForCancelLiftedButListStillMarked` 仍 PASS ⇒ 今天无牙；**装牙后读数**＝本腿新 leg 具名红。
- **M-2c**＝正控：overlay 替换测试文件里的读路径 ⇒ 归档 leg 变红，证明 overlay 打得进测试文件、也证明 M-2a 的绿是"读盘看不见"而不是"overlay 没落地"。
- **M-2b／M-2e**＝能力侧突变：`BuiltinTaskEntries` 名册里加／减 `task.list`，分别证新判据的**两个方向都有牙**、且合法同批核销（摘标记＋接线同时）**不假红**——这一发同时把票面 §四 说的计数形 `len(names) != 2` 的假红风险实测出来。
- 每发记：摘哪一行＋overlay json 路径＋落地证明＋红句原文＋还原后 `md5sum` 与 `git show HEAD:<file> | md5sum` 对拉。

## §4 门禁读数

打算答：`sh scripts/d22scan.sh`（含正控）、`gofmt -l`、`gofumpt -l` 对本腿新文件、`go vet ./internal/tools/`、`git status --porcelain -- internal cmd` 收尾为空。本腿若动 `scripts/**` 则补 `bash -n`；⛔ 本腿不动 `.github/**`。

## §5 判不动的地方

打算答（实质内容，不写"没有"）：AC#2 射程外但今天量到的东西——票 225 的双向对账、`task.list` 到底该不该注册（D34／§7 :1531，归 owner 与产码腿）、读盘尺在产码里的其它同类（若查到）、以及本腿判据"哪一维今天量不到"。

---

★本节以上（§0–§6）没有一句结论来自上面那份骨架：它每一节只写「打算答」，是实话句、不是判语。
本腿唯一提到它的地方＝§3 末段「前腿报过的形状」，那里写死**既不采信也不推翻**，只登记形状、另跑自己的数。
