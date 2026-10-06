# 271-a1 — 票 271 AC#0 三张代价表的**只读现量普查件**

- 腿：`271-a1`（只读代价普查腿）／时刻 2026-10-06 19:3x +08／分支 `dev`
- 射程：只为 AC#0 三张代价表出**现量**。⛔ 本腿不选形、⛔ 一枚判语不给、⛔ 票面 `- [ ]`／`- [x]` 一字未动。
- 硬约束兑现：**零 `go` 命令**（§0.3 自证）／零产码／零 `docs/**` 改动。
- 长输出先落 `.scratch/wisp/probes/271/a1/logs/*.txt`，本件只抄读数。

---

## §0 起手锚

| 尺 | 命令原文 | 读数 |
|---|---|---|
| HEAD 短号 | `git rev-parse --short HEAD` | `bd39f17f` |
| 分支 | `git branch --show-current` | `dev` |
| HEAD 标题 | `git log -1 --format='%h %ad %s'` | `bd39f17f Tue Oct 6 19:21:34 2026 +0800 probes(票 111 追加编排者记 19:2x)：ci-census-read-1 死腿的 logs 我自己复跑后读回三问` |
| 被审两棵树脏度 | `git status --porcelain -- cmd internal` \| `wc -l` | **0 行**（＝派单"此刻应为 0"成立） |
| 全树脏度（背景，非本腿所造） | `git status --porcelain` \| `wc -l` | **750 行**（别人的活；本腿未碰其中任何一枚） |

### §0-补 被审文件行数与字节（`wc -l -c` 原文）

```
 289  63098 .scratch/wisp/probes/236/v1/verdict.md
 889  37391 internal/tools/task.go
 223  12328 internal/tools/tasklist_deferred_236r3_teeth_test.go
 535  21579 internal/tools/task_cancel_221_legs_test.go
1427  67155 cmd/wisp/run.go
```

### §0.3 零 `go` 命令的自证（本腿跑过的命令全集，逐枚）

`git rev-parse` · `git branch --show-current` · `git status --porcelain` · `git log --oneline --diff-filter=A -- <path>` ·
`git log --oneline -- <path>` · `git log -1 --format` · `grep -rn / -rl / -rnE / -c`（含 `-A`／`-B`） · `sed -n` · `awk` · `cut` ·
`sort` · `uniq` · `wc -l / -c` · `ls` · `mkdir -p` · `head`／`tail`（管道末端） · `cat`（只读本腿自己的 logs） ·
`md5sum` ＋ `git show HEAD:<file> | md5sum`（§1.5(b) 那把对拉，⛔ 只读不改盘） ＋ 工具侧 `Read`／`Write`。

⇒ **清单里没有任何一枚 `go`**：无 `go test`／`go build`／`go vet`／`go list`／`go env`／`go run`，
也没有任何间接起 `go` 的包装件（`scripts/*.sh`、`tools/d22scan/runtests.sh` 一枚未跑）。
⇒ 与那枚正在 `cmd/wisp` 跑突变的并发腿**无进程争用**，它的读数没被本腿洗过。

### §0.4 地界自证（⛔ 未读、未引用）

`.scratch/wisp/probes/232/**` · `pool-validity/5g2/**` · `expired-premises/7b2/**` · `probes/270/**` · `frontend/**` · `design/**`
＝**一枚未打开**。`.scratch/wisp/probes/236/**` 只读派单点名的 `236/v1/verdict.md`（§2 第 8／9 行、§5 第 3 条），
其中读数在本件里一律**标注为「`236-v1` 造的、本腿未重跑」**，⛔ 未把它的判语当本腿结论。

---

## §1 尺与名册逐字（全部本腿现读，⛔ 不抄票面也不抄派单的号）

### 1.1 尺本体＝`internal/tools/task_cancel_221_legs_test.go:183-185`（逐字，含行号）

```
183:	if len(names) != 2 {
184:		t.Errorf("task 家族注册了 %d 枚：%v（task.list 仍须是 DEFERRED，本票不许顺手注册它）", len(names), names)
185:	}
```

⇒ **派单给的两枚锚点未漂**（`:183` 与 `:184` 逐字一致，含"本票不许顺手注册它"那句）。
同一枚用例的上下文（本腿现读）：

```
170:	func Test221TaskCancelIsRegisteredAtItsFrozenLevel(t *testing.T) {
171:		entries := BuiltinTaskEntries(TaskDeps{Roster: NewTaskRoster()})
172:		names := make([]string, 0, len(entries))
174:		for i, e := range entries {
175:			names = append(names, e.Tool.Name())
180:		if cancelEntry == nil {
181:			t.Fatalf("BuiltinTaskEntries 没注册 task.cancel（现名册：%v）—— 说明书又在许诺一枚不存在的工具", names)
187:		if decl.Declared != risk.L1 {
190:		if len(decl.Capabilities) != 0 || len(decl.Needs) != 0 {
193:		if len(decl.PathParams) != 0 {
```

★结构事实（选形要用它，本腿不据它下判语）：**同一枚用例里 `:181` 已经是"按名字判"的尺**（`task.cancel` 缺席即 `Fatalf`），
`:183` 是**同一具名分母上唯一的枚数形断言**；`:187`／`:190`／`:193` 数的是被点名那一枚 entry **自己的字段**，
合法新增第三枚不会让它们红。⇒ 乙形的写面**只可能落在 `:183-185` 这一块**。

同一文件另一枚相关尺（词面共现，票 221 AC#4 改写形的 (b) 支；本腿现读）：

```
223: func Test221DeferredMarkerForCancelLiftedButListStillMarked(t *testing.T) {
224: 	src, err := os.ReadFile("task.go")
     （循环数 cancelMarked / listMarked）
240: 	if cancelMarked != 0 { … }
244: 	if listMarked == 0 {
245: 		t.Errorf("task.list 的 DEFERRED 标记不见了：本票不许摘它（那半支仍然没有实现）")
```

### 1.2 名册本体＝`internal/tools/task.go`（逐字）

```
597: // BuiltinTaskEntries returns the task family this configuration registers:
598: // task.output and (since ticket 221 甲形, ruling A434) task.cancel. See the
599: // header comment for why task.list is not here.
600: func BuiltinTaskEntries(d TaskDeps) []Entry {
601: 	return []Entry{
602: 		{Tool: taskOutput{d: d}, Decl: TaskOutputDecl()},
603: 		{Tool: taskCancel{d: d}, Decl: taskCancelDecl()},
604: 	}
605: }
```

⇒ **派单参考的 `:600-605` 轻漂**：`:600-605` 是**整个函数体**；两枚注册行是 **`:602`／`:603`**。
`236-v1` 记的"`task.go:603` 后插一行"与盘上对上（`:603`＝`task.cancel` 那一行）。

DEFERRED 标记行与那句共现注释（逐字；禁区件，本腿一字未动）：

```
23: //	task.list    -   DEFERRED with five fields, PLAN.md §7 :1531
279: // diagnostics, not as a tool surface: task.list is DEFERRED (§7 :1531), and
280: // giving task.output a "list everything" mode would smuggle that row back in.
```

那把词面尺"共现数撑在 1"的来源（本腿自数，口径＝同时含 `DEFERRED` 与 `task.list` 的行）：

```
grep -c DEFERRED internal/tools/task.go      → 3   （:23 ／ :32 ／ :279）
其中同时含 task.list 的行                     → 2   （:23 标记行 ＋ :279 那句 Count() 文档注释）
```

⇒ 与 `docs/reports/pending-and-issues.md:12619` 记的"3 枚 DEFERRED、2 枚 co-carry"**同口径可对拉**。
★`:279` **不是标记行，是一句注释碰巧同时含两个词**——它就是词面尺被叙述句撑住的形状（票 236 AC#2 的分母，本腿未动）。

### 1.3 生产调用者（★锚在**带括号的调用形状**，⛔ 不用裸符号名）

调用形状尺：`grep -rn "BuiltinTaskEntries(" --include="*.go" cmd internal`

```
cmd/wisp/run.go:544	for _, e := range tools.BuiltinTaskEntries(tools.TaskDeps{Roster: rt.tasks, Paths: rt.paths}) {
internal/tools/task.go:600	func BuiltinTaskEntries(d TaskDeps) []Entry {	← 声明，不是调用
其余 12 枚命中全在 _test.go：pointer_183_cli_seam_test.go:111 ／ pointer_185_cli_seam_test.go:103 ／
subagent_197_test.go:201 ／ tasklist_deferred_236r3_teeth_test.go:130 ／ task_cancel_221_legs_test.go:101 与 :171 ／
task_cancel_221_test.go:50 ／ task_output_leg_test.go:40 与 :144 ／ task_output_pointer_notice_test.go:54 ／
ticket175r2_stamp_live_test.go:79 与 :353
```

⇒ **产码（非 `_test`）里的调用点＝1 枚：`cmd/wisp/run.go:544`**（派单锚点**未漂**）。
★裸符号名的对照读数（这一发就是派单点名的那把尺差）：`grep -rn "BuiltinTaskEntries" --include="*.go" .` ＝ **114 枚**
（含 `.scratch/wisp/probes/**` 里别人突变副本的声明与构造），`grep -rn "TaskDeps{"` ＝ **51 枚** ⇒ 114 枚命中里能算"被产码调用"的只有 1 枚。
另：`grep -rn "BuiltinTaskEntries(.*)\[" --include="*_test.go" internal cmd` → **rc=1 零命中**（＝好消息：
**没有一把尺按下标读名册**，乙形不必担心"顺序变了"那一维）。

**双向追到 `main`**（逐层现读，⛔ 不引任何腿的结论）：

| 层 | 证据（file:line＋逐字） |
|---|---|
| 被调者 | `internal/tools/task.go:600` |
| 调用者 | `cmd/wisp/run.go:544` |
| enclosing | `cmd/wisp/run.go:382` `func assembleRuntime(s runSpec) (*agentRuntime, int)`（awk 现取） |
| assembleRuntime 的**非测试**调用者 | `grep -rn "assembleRuntime(" --include="*.go" cmd internal` 去掉 `_test` ⇒ **2 枚**：`cmd/wisp/run.go:249`（enclosing `func runTextTask(s runSpec) int`＝`:181`）＋ `cmd/wisp/resident_task_source_windows.go:278`（enclosing `func startResidentTaskSource(...)`＝`:228`）；另 4 枚在 `_test.go` |
| runTextTask | `cmd/wisp/main.go:159` `return runTextTask(runSpec{`（在 `func cmdRun(args []string) int`＝`:151` 内） |
| 入口 | `cmd/wisp/main.go:58 func main()` → `:91` `case "run":` → `os.Exit(cmdRun(args[1:]))` |

⇒ **现量结论（不是判语）**：名册在**两条**产码路径上都被装配（文本 CLI `wisp run` ＋ 常驻任务来源那一支），
所以这把尺量的**不是测试专用件**。`run.go:537` 那段注释逐字写明谁填名册（`Who FILLS the roster is ticket 197's spawner`，在 `:537-541` 那一段里）。

**装配根一共往同一枚 registry 里放三家**（本腿现读 `grep -n "\.Register(" cmd/wisp/run.go` 去掉 `_test` ⇒ 三处循环）：

```
cmd/wisp/run.go:517-519	for _, e := range tools.BuiltinFSEntries(tools.FSDeps{Paths: rt.paths, DeleteEnabled: cfg.FS.DeleteEnabled}) { reg.Register(e) }
cmd/wisp/run.go:544-545	tools.BuiltinTaskEntries(...)  → reg.Register(e)
cmd/wisp/run.go:780-791	tools.BuiltinSubagentEntries(...) → reg.Register(e)
```

内置构造器全集（`grep -rnE "^func Builtin[A-Za-z]*Entries" --include="*.go" internal cmd`）＝**4 枚**：
`fs.go:325`／`fs_write.go:772`／`subagent_197.go:233`／`task.go:600`（`fs_write` 那一枚装配根没用，`task_cancel_221_test.go:41-43` 的并集里用了）。
⇒ ★**没有任何一把尺数"装配根最终注册出去的全集枚数"**：`grep -rn -A3 "\.Tools(" --include="*_test.go" internal cmd | grep -E "len\("` ＝ **rc=1 零命中**（`logs/union-counts.txt`）。
★**这把零命中本腿按"过滤后空 ≠ 没跑"自证过两段**（`logs/zero-hit-proof.txt`）：第一段 `grep -rn -A3 "\.Tools(" --include="*_test.go" internal cmd` ＝ **24 行**（真命中 5 枚调用点：
`bridge_test.go:729`／`fs_test.go:259`／`fs_write_test.go:569`／`subagent_197_test.go:641`／`task_output_ac2_before_test.go:71`），第二段那道 `len(` 过滤才是 rc=1
⇒ **"零命中"是"这些调用点上没有计数断言"，不是"我 grep 错了"**（同一把尺的第二个用法：`grep -rn "BuiltinTaskEntries(" --include="*_test.go" internal cmd`＝**12 行**，
再 `grep "\["`＝**0 行** ⇒ 真的没有一把尺按下标读名册）。
⇒ 这条现量对选形要紧：**乙形把 `:183` 换成名集判之后，"整个装配的枚数"这一维仍然零尺**（本来也是零尺），
不存在"改完乙就没人管枚数了"这一说——今天管枚数的只有**分家族**那 5 枚。


### 1.4 ⚠ 派单／票面锚点漂移核查表（逐枚具名，未漂的也列）

| 派单／票面锚点 | 盘上真值 | 判定 |
|---|---|---|
| `task_cancel_221_legs_test.go:183` ＝ `if len(names) != 2 {` | 同 | **未漂** |
| `:184` ＝ 那句红句 | 逐字同 | **未漂** |
| `task.go:600-605` ＝ `BuiltinTaskEntries` 的**构造行**、现量 2 行 | 函数体 `:600-605`；两枚注册行＝**`:602`／`:603`** | **轻漂**（下一枚腿照抄会把 `:597-599` 注释算进写面） |
| `task.go:23` DEFERRED 标记 | 同 | **未漂** |
| `task.go:279` 带 `task.list` 的文档注释 | 同（`:280` 是同句第二行） | **未漂** |
| `cmd/wisp/run.go:544` 生产调用者 | 同 | **未漂** |
| `tasklist_deferred_236r3_teeth_test.go:174` 红句 | 同（逐字 `task 家族名册 = %s, want 含 task.output 与 task.cancel`） | **未漂** |
| 票 221 老票面 `internal/tools/task.go:595` | 现 `:600` | **修码前的历史读数**，不算漂移（但引用它的人会漂） |

### 1.5 凭据件的引用口径**逐枚验过**＋二次锚点核查（本腿自己拉尺，⛔ 不照抄派单）

**（a）凭据件 §2 那张表的**行序号**本腿自己数过（口径＝`:112` 表头／`:113` 分隔行之后的**数据行**，第 1 枚＝`:114`）：**

| 数据行序号 | verdict.md 位置 | 那一格的内容 |
|---|---|---|
| 1–4 | `:114`–`:117` | AC#1 两支／AC#1b 两支 |
| 5 | `:118` | AC#1b `parentID == ""` ⇒ ★**禁区别处引用的那句就在这一格**，逐字 `⛔ 本腿不据此判"断言太弱"，也不许任何人为了看孤行去**新增名册读口**（＝产码改动＋可能新导出名）` ⇒ **票面 §禁区那句"§2 第 5 行"引用成立** ✔ |
| 6 | `:119` | 一票两形：panic 那一形 |
| 7 | `:120` | AC#2 (b) 支"今天无牙"（`B-1`／`B-6`） |
| 8 | `:121` | AC#2 (b) 支改扫能力＝**`B-4`（注册 `task.list`，FAIL=4）／`B-5`（注销 `task.cancel`，rc=1／PASS=196／FAIL=10 枚具名红）** |
| 9 | `:122` | AC#2 欠量①（(a) 支，`B-2a`） |
| **10** | **`:123`** | **AC#2 欠量② `m-2c`＝`B-3`（注册无关的 `task.note`）⇒ rc=1／PASS=203／**FAIL=3**；★名册钉 **PASS**＝假红预测坐实；判语栏逐字把"枚数形旧尺会假红"作为**附带读数交回编排者** |
| 11 | `:124` | AC#2 的仪器事实 |

⇒ ★**一处序号引用要更正给编排者**：票面 Status 行与派单说的"§2 **第 8 行**的附带读数"，按本腿这把数法，
"枚数形尺今天会假红"那条附带读数的真身＝**第 10 枚数据行（`:123`）**；第 8 枚（`:121`）是 `B-4`／`B-5` 那一格。
⇒ **内容全都对得上、序号差两枚**（本腿不猜是哪一种数法造成的，只具名报）。
★更要紧的一处**数字归属冲突**（已进 §3 第 5 条）：票面 §现量 3 写的是
`B-3`＝注册无关的 `task.note` ⇒ **rc=1／PASS=196／FAIL=10 枚具名红**＋"其中含两枚 AC#1 拒绝钉走真注册路径"，
而凭据 `:123`（`B-3` 那一格）逐字＝**PASS=203／FAIL=3**、`:121`（`B-5` 那一格）才是 **PASS=196／FAIL=10 且含两枚拒绝钉**。
⇒ **票面 §现量 3 把 `B-5` 的三个数（196／10／两枚拒绝钉）挂到了 `B-3` 名下**；本腿 ⛔ 未重跑、原样报回，不裁谁对。
（`236-v1` §5 第 3 条＝`:261` 同口径：`B-5`＝10 枚／`B-4`＝4 枚／`B-3`＝3 枚，且逐字
`⛔ 本腿没有因红得多而判"钉太宽"，也没有建议收窄任何断言（禁区）`。）

**（b）二次锚点核查（本腿写完全件后重拉一把，防共享工作树漂移）**：

```
git rev-parse --short HEAD                                 → 6895b994（本腿第一发交件；起手快照＝bd39f17f）
git log --oneline bd39f17f..HEAD --name-only -- internal cmd → 零枚命中（rc=0，无任何 commit 动过 internal/** 或 cmd/**）
git status --porcelain -- cmd internal | wc -l              → 0
md5 对拉（disk vs `git show HEAD:<f>`）                     → 五枚全部 same=YES：
  internal/tools/task.go                        4138177e29ffff427776a73eb0d61a9b
  internal/tools/task_cancel_221_legs_test.go   13a0ba39e30fb95ad4f2be88beed4fc7
  internal/tools/tasklist_deferred_236r3_teeth_test.go 21b53d52ac87b893b74fdd9828f2f3aa
  cmd/wisp/run.go                               1d95cfaf116d9092514ee0a905a82b23
  internal/tools/ticket175r2_stamp_live_test.go 2a5fde81b78e56c8644c8463ce8738e3
```

⇒ 本件所有静态读数**落在同一份产码字节上**，且那五枚 md5 与 `236-v1` §4 第 8 行自报的前四串**逐字相同**
（`4138177e…`／`13a0ba39…`／`21b53d52…` ＋ `run.go` 那一枚它未列）⇒ **两枚腿取的是同一份字节**，本件的静态归因可以直接接在它那些读数后面。
★注意：`cmd/wisp/run.go` 的 md5 本腿量到 `1d95cfaf…`，而 `236-v1` §4 第 8 行那五枚文件名册里那串第五位是 `21b53d52…`＝**它排的是 `236r3` 那枚文件**，
两枚腿的被审文件集**只有 4 枚重叠**（本腿多一枚 `ticket175r2_stamp_live_test.go`）⇒ 不构成矛盾，只记口径差。

**（c）射程外那一枚"数工具枚数"的雷（不在本票两张表里，但派单第 4 条问的是"全仓还有几把"，具名给）**：

- `grep -rlE "len\(…\) (!=|==) [0-9]+" --include="*_test.go" internal cmd tools pkg` 去掉 `internal/tools/`＋`cmd/wisp/` ⇒ **139 枚文件**（口径＝**文件枚数不是断言枚数**，命令与读数在 `logs/marker-ruler.txt` 末段）。
- 其中**数工具名册**的：零枚（逐枚看过 `logs/union-counts.txt` 里的候选：`internal/agent/spill_*` 数的是目录项、`internal/ball/hotkey_status_test.go:398` 数的是注册尝试、`internal/risk/pointer_185_test.go:367` 数的是污点标记）。
- 一枚**同族形状但射程是面板 API**的：`internal/agent/approval/ticket259_panel_capability_rulers_test.go:81` `if typ.NumMethod() != len(panelAPIAllowedMethods) {`
  ⇒ 它是"枚数对名集长度"的**双形**写法（乙形若想照同一形状办，这一枚是仓里现成的**写法先例**，⛔ 不是放行理由）。
- ⇒ **结论性现量**：**枚数形工具尺在整个 `internal/**`＋`cmd/**` 全集里只有 §2.0(a) 那 5 枚**，task 名册那一枚只有 1 枚。

---

## §2 三形代价表

### 2.0 两张分母先钉清（三张表共用）

**（a）枚数形尺全集**——口径＝`len(...) != <数字>`／`== <数字>` 落在 `internal/tools/**_test.go` ＋ `cmd/wisp/**_test.go`。
原始命中 **331 枚**（`logs/count-rulers.txt`）→ 名字相关收窄 **104 枚**（`logs/count-rulers-narrow.txt`）→ 逐枚读上下文后**分类**
（`logs/ctx-hits.txt`／`logs/extra.txt`／`logs/exactN.txt`／`logs/misc-context.txt`）。

**数名册／工具枚数的尺＝5 枚**（逐枚 file:line＋断言逐字＋它数的是什么）：

| # | file:line | 断言逐字 | 数的是什么 | 新增第三枚 **task** 工具会不会红 |
|---|---|---|---|---|
| 1 | `internal/tools/task_cancel_221_legs_test.go:183` | `if len(names) != 2 {` ＋ `:184` 红句（见 §1.1） | **task 家族名册枚数**（分母＝`BuiltinTaskEntries`） | **会**——本票射程内唯一一把 |
| 2 | `internal/tools/fs_test.go:190` | `if len(entries) != 6 {` ＋ `t.Fatalf("the default fs roster has %d tools, want 6 (read/list/write/edit/trash/move)", len(entries))` | **fs 默认名册枚数** | 不会（不同名册） |
| 3 | `internal/tools/fs_test.go:235` | `if len(entries) != 7 {` ＋ `t.Fatalf("with delete_enabled the roster has %d tools, want 7", len(entries))` | **fs 名册（带 `delete_enabled`）枚数** | 不会 |
| 4 | `internal/tools/fs_test.go:263` | `if len(dir) != 6 {` ＋ `t.Fatalf("directory has %d entries: %+v", len(dir), dir)` | **模型可见工具目录枚数**（`b.Tools()`，`:259`；桥＝`:258` 的 `fsBridgeWith`，其构造面＝`internal/tools/helpers_test.go:45` `for _, e := range BuiltinFSEntries(FSDeps{Paths: paths})` ＝**只装 fs、不装 task**（本腿现读，`fs_test.go` 里 `grep -c BuiltinTaskEntries`＝**0**）） | 不会（静态可判：task 那两枚根本不进这把尺的桥；⛔ 仍未跑，见 §4 U1） |
| 5 | `internal/tools/bridge_test.go:292` | `if len(AllCapabilities) != 11 {` ＋ `t.Fatalf("C3 has %d tokens, want 11", len(AllCapabilities))` | **C3 能力 token 枚数**（不是工具枚数） | 不会；若新工具**顺带**加能力则红＝**那是该红的** |

**同族但形状不同**（`== 0` 的"非空分母自证"，合法新增**不会**红、摘掉注册才会）：
`internal/tools/ticket175r2_stamp_live_test.go:356` `if len(names) == 0 {`＋`t.Fatal("一枚内置工具都没注册到，这条普查就是空转")`、
`internal/tools/task_cancel_221_test.go:65`＋`:130` 同形、`internal/tools/tasklist_deferred_236r3_teeth_test.go:139`＋`:140` 同形。
★**这把尺本腿放宽过窗口重扫一次以防漏**（`-A3` 可能盖不住隔两行的断言）：`grep -rn -A8 "\.Tools(" --include="*_test.go" internal cmd | grep -E "len\((dir|list|tools)"`
⇒ 只多出一枚 `internal/tools/bridge_test.go:733` `if len(dir) != 0 {`＋`t.Fatalf("empty bridge reports %d tools", len(dir))`，
它钉的是 `New(Options{})` 那枚**空桥**（`:728`），分母里没有任何家族 ⇒ 合法新增同样不红，**上表那 5 枚是全集**（读数 `logs/tools-window.txt`）。

**枚数形＋名字形并存**（乙形最容易被误读成"也要动它们"，具名列出）：
`internal/tools/fs_test.go:190` 之后**另有一张 `want map[string]risk.Level` 名集表**（红句 `:211` `unexpected fs tool %q`）
＋`:225` 一轮 `fs.delete` 缺席断言 ⇒ **fs 家族那把尺本身就是"枚数＋名集"双形**；
task 家族的名集那一半在**另一枚文件**里（`236r3`，见 2.0(b)）——**这是两家族形状最要紧的差别，选甲／乙的决定料**。

**结构／接口面的枚数尺**（合法新增一枚工具**不会**红，但"多开一条出口／多一个字段"会红；列出来是防止被当枚数尺漏判或误判）：
`internal/tools/subagent_197_test.go:879` `if got := subagentDepsFieldNames(); len(got) != 5 {`、
`cmd/wisp/subagent_selfapproval_197_test.go:569` `if len(gateNames) != 2 || !hasMethod197(gateNames, "PendingWindow") || !hasMethod197(gateNames, "PendingApproval") {`（名集＋枚数**同一行**）、
`internal/tools/ticket224_setter_scope_test.go:53` `if len(grantFields) != 1 {`、
`internal/tools/paths_twocontainments_252_r2_test.go:296` `if len(book) != 2 {`（AST 授权册读取次数；红句自己写着 `the ruler is the thing that needs a human approval, not a green test run`）、
`cmd/wisp/resident_approval_246_windows_test.go:146` `if len(records) != 10 {`＋`cmd/wisp/resident_approval_live_246_windows_test.go:245` 同形、
`internal/tools/fs_edit_ac3b_test.go:118` `if len(original) != 14`（字节数）、
`internal/tools/task_output_leg_test.go:263`／`:266`（D15 头尾字节数）、
`cmd/wisp/firstrun_257_test.go:297`／`cmd/wisp/secret_test.go:985`（计数为 3 的行为面）、`cmd/wisp/providers_test.go:112`（报告行字段数）。

⇒ **2.0 的答案（选甲／乙的决定料）**：
1. `cmd/wisp/**_test.go` 里**没有一枚**数工具名册的枚数尺（8 处命中逐枚读过：全是 `== 0` 自证或行格式／记录数）。
2. **"合法新增一枚 task 工具"能直接顶红的枚数形尺＝1 枚**（`task_cancel_221_legs_test.go:183`）。
3. 同族雷**不属于 task 家族**的有 4 枚（`fs_test.go:190`／`:235`／`:263`＋`bridge_test.go:292`）
   ⇒ **"合法新增一枚工具会让几把枚数形尺同时假红"是分家族的答案：task＝1 枚、fs＝3 枚、能力面＝1 枚（且那一枚该红）**。
   ★票 271 的票面射程只写了 221 那一把；`236-v1` 那 10 枚红里也**不含** fs 那三枚（它们分母不同）。

**（b）已经按名字集合判 task 家族的名册钉**（乙形的"顶红谁"与丙形的"能不能区分"都靠它）：

| file:line | 断言逐字（关键处） | 判据形状／区分到哪一层 |
|---|---|---|
| `internal/tools/tasklist_deferred_236r3_teeth_test.go:130` | `for _, e := range BuiltinTaskEntries(TaskDeps{Roster: NewTaskRoster()}) {` ＋ `:131-136` 两枚 `t.Fatal`（nil Tool／无名）＋ `:139-141` 非空分母自证 | 分母自证，**不是**枚数尺 |
| 同文件 `:153-161` | `if n == "task.list" { t.Errorf("DEFERRED 那一支和名册分叉了：BuiltinTaskEntries 的注册名册里出现了 task.list 这一行（现名册：%s）…") }` | **名集·负向**：`task.list` 在册即红 |
| 同文件 `:164-177` | `if !sawOutput \|\| !sawCancel { t.Errorf("task 家族名册 = %s, want 含 task.output 与 task.cancel（…）") }` | **名集·正向（含即可）**：多一枚**不**红 |
| 同文件 `:178` | `t.Logf("现名册：%s（DEFERRED 的判据读的是这一行，不是注释里的词面）", …)` | 自报口径（非断言） |
| 同文件 `:145-149` 头注释 | `Deliberately NOT a count, so a legitimately approved third task tool cannot false-red (交件 §3 M-2c measures that pair).` | **写下来的意图**：名集形明知故不数枚 |
| 同文件 `:71-74` 头注释 | `⛔ nothing new is invented, nothing existing is deleted or loosened: … the count-shaped assertion at task_cancel_221_legs_test.go:183 keeps counting` | **它对那把枚数尺的处置＝原样保留**——乙形动的正是这一句承诺的范围 |
| `internal/tools/task_cancel_221_test.go:134-141` | `for _, root := range roots { if registered[root] { continue } … t.Errorf("说明书对模型许诺了一枚不存在的工具：%s 的 %s 写着「%s」，而 %q 没有注册进这次装配的并集（并集 %d 枚：%s）") }` | **并集名集**（分子＝说明书词根、分母＝四家 `Builtin*Entries` 构造并集，`:40-47`）：多注册**不**红、少注册**会**红。★红句里那个 `%d 枚` 是**报告**不是**判据** |
| `internal/tools/ticket175r2_stamp_live_test.go:349-377` | 四支 `case marked && !decided:` 等；表＝`:333-348` `classified175r2`（含 `"task.cancel": "unstamped: a receipt naming which roster row this call stopped"`） | **逐名分类普查**：新增第三枚**必须就地答一句**否则红（红句自己写"新工具请就地回答"）＝**真伤维** |

### 2.1 甲形（不动断言，升格成明写的警报线＋一枚正控）

| 栏 | 现量 |
|---|---|
| 动哪几行 | `internal/tools/task_cancel_221_legs_test.go`：`:183` 上下**加注释**＋（可选）改 `:184` 的**文案**（票面 §甲 要求"在文案里写清响了之后要走哪道工序"）；**新增一枚 overlay 台件**（AC#1 要求）＋"删掉台件行后名册尺回到起手读数"的还原对拉。`internal/tools/task.go`／`cmd/wisp/**` 一字不动。 |
| 顶红谁现有的钉 | **HEAD 态：0 枚**（注释与文案都不改任何断言的形状；2.0(b) 那几枚名集／普查钉照旧绿，`fs_test.go`／`bridge_test.go` 那四枚分母不含 task）。<br>★**跑正控时那 1 枚非零**：甲形正控要 `-overlay` 换进 `task.go`，而同包仪器尺 `internal/tools/tasklist_deferred_236r3_teeth_test.go:202` `func Test236R3InstrumentFactReadDiskIsBlindToOverlayCarriesIt`（`:109` `//go:embed task.go` 与 `:204` `os.ReadFile(taskSrcPath236r3)` 两读对拉，红句在 **`:212`**`编译期与运行期读到的 task.go 字节不一致（编译 %d 字节 / 盘上 %d 字节，标记行数 %d 对 %d）`与 **`:217-218`**`同一枚 DEFERRED＋task.list 标记行数，编译器看到 %d 行、os.ReadFile 看到 %d 行`）在**任何**换入 `task.go` 的**不同字节**的 overlay 上都红——该文件头注释 `:33-45` 与 `:197-199` 逐字自认这是它**存在的读数**（`this leg goes red under ANY -overlay of that file, benign ones included … not a defect and not a new always-red`）；`236-v1` §4 第 8 行同口径。⇒ **甲形每跑一次正控必带 1 枚"红得有道理"的仪器红**，名册口径要提前写清，否则下一位会读成"新增了常红"。 |
| 要不要解冻别人 `-done` 票 | **只加注释＝不动判据**；★但若甲选择**改 `:184` 那句红句文案**，那一行属已 `-done` 票 221 的判据文件 ⇒ 严格讲仍需一枚**只到文案**的具名解冻。**这一支算不算"改判据"＝要编排者裁**（§3 第 4 条）。票 221 的 `-done` 现量见 §2.2 同名栏（同一枚文件、同一发 commit）。 |
| "那一天接手的人会看到几枚红" | **不减**：甲买到的是"红得**有说明书**"，枚数尺照旧红。已测单侧读数（`236-v1` 造，本腿**未重跑**）：多注册无关的 `task.note`＝**3 枚**（§5 第 3 条：枚数尺＋`175r2` 普查＋仪器尺）；多注册 `task.list`＝**4 枚**（`B-4`）；注销 `task.cancel`＝**10 枚**（`B-5`）。★"合法新增那一发共几枚"本件**给不出确数**（§4 U1／U2）。 |

### 2.2 乙形（改成按名字集合判，枚数只作次要信号）

| 栏 | 现量 |
|---|---|
| 动哪几行 | **只有** `internal/tools/task_cancel_221_legs_test.go:183-185` 这一块（`len(names) != 2` 换成期望名集判）。同一用例其余断言（`:180-181`／`:187`／`:190`／`:193`）与同文件其余用例**零必要改动**（§1.1 结构事实）。名册那两行（`task.go:602`／`:603`）**不必动** ⇒ 零产码。 |
| 顶红谁现有的钉 | **代码层 0 枚**。逐枚核过理由：`236r3:153-161`／`:164-177` 是名集判、HEAD 本就绿，乙形同向；`:139-141` 只反空分母；`task_cancel_221_test.go:134-141` 只反"许诺了没注册"；`ticket175r2:349-377` 只反"没答"；`fs_test.go:190/235/263`＋`bridge_test.go:292` 分母不含 task 名册。<br>★**不红在代码、红在承诺**：`236r3:71-74` 那句 `nothing existing is deleted or loosened … keeps counting` 是票 236 交付时对这把尺的**处置承诺**（该文件 AC#2 框已 `[x]`——`236` 票面 `- [x] AC#2`，本腿现读），乙形动的正是它。⇒ **"顶红"这一栏在乙形是零枚代码钉＋一枚写在注释里的承诺**，这是它和甲形代价的真实差别所在。 |
| 要不要解冻别人 `-done` 票 | **要，且这是前置不是理由**。现量：`git log --oneline --diff-filter=A -- internal/tools/task_cancel_221_legs_test.go` ＝ **唯一一枚 `c44b30c4 221 甲形＋乙形同一发：注册 task.cancel 并把说明书那句裸许诺换成当下为真的话`**；`git log --oneline -- <该文件>` 也**只有这一枚** ⇒ 尺的主人＝票 221（`.scratch/wisp/issues/221-task-spawn-description-promises-task-cancel-that-is-not-registered-done.md`，AC#1–AC#5 **五格全 `[x]`**）。★**枚数尺不在票 221 任何一格 AC 原文里**（五格逐字读毕：AC#1 并集能力尺／AC#2 生产调用者／AC#3 两枚正控＋非级联／AC#4 DEFERRED 标记／AC#5 整包读数），裁决表 `docs/evidence/s1/221-task-cancel-v1.md` 里也没有一处把 `:183` 写进任一 AC 的判语 ⇒ **"改它算不算改别人已勾判据"今天没有书面答案**（这条摆给编排者，本腿不裁）。<br>落台账要什么（**形**＝仓里既有定式，⛔ 不是"有先例所以可放行"）：`docs/reports/pending-and-issues.md:11738` `## A602（2026-10-04 15:18:03，**具名解冻**·五样齐：票 256 那枚 AST 钉的 Options 字段集期望集 5 枚 → 6 枚）`，其 `:11743` 边界逐字 `**只解冻那一处期望集与那一条反向断言**。⛔ 不含同包其它任何断言…`、`:11744` 撤销口令逐字 `**「265 撤 AST 钉解冻」**`；同形另见 `:11351`／`A595`（`:11542`）／`A598`（`:11629`）／`A611`（`:11939`）。票面 AC#4 原话照抄在这里：`⛔ 不许按"票 221 刚被别人动过"当先例放行`。 |
| "那一天接手的人会看到几枚红" | 枚数尺那一维**不再红** ⇒ 相对甲**至少 −1 枚**。已测单侧（`236-v1`，未重跑）：`task.note` 那发 **3 枚 → 2 枚**，剩下那两枚是 `ticket175r2:349-377` 逐名普查（**真伤**：新工具没在分类表里答一句）＋ overlay 仪器尺（**必带**）；`236r3` 名集钉在两发里都 PASS ⇒ 乙形不把任何现有钉改红。★确数同样要一次实测才能钉（§4 U1）。 |

### 2.3 丙形（不做，登记为已知常红风险）

| 栏 | 现量 |
|---|---|
| 动哪几行 | **0 行代码、0 行票面判据**。只：台账登记＋票面 §AC#0 那句"选丙那要在台账写明今后任何腿接第三枚 task 工具前须先读本票"（归编排者写，本腿未动 `docs/**`）。 |
| 顶红谁现有的钉 | **0 枚**（什么都不碰）。 |
| 要不要解冻别人 `-done` 票 | **不需要**——丙形是三形里唯一一枚**零解冻**的。 |
| "那一天接手的人会看到几枚红、其中几枚是真伤" | 会先看到：① `task_cancel_221_legs_test.go:183-184` 那枚**假红**（它自己那枚 `t.Errorf`，非致命，用例其余部分照跑）；② 同包 `ticket175r2_stamp_live_test.go:349-377` 那枚**真红**（新工具没答分类，红句自己写"新工具请就地回答"）；③ 若那一发同时 overlay 了 `task.go` → `236r3` 仪器尺那枚**必带红**。已测单侧枚数（`236-v1`，⛔ 本腿未重跑、也不在手上）：`task.note`＝**3 枚**、`task.list`＝**4 枚**、摘 `task.cancel`＝**10 枚**。⇒ **"接手的人先看到几枚"这一格在丙形是"三到十枚，取决于那一枚动的是哪一支"，且他手上没有本件这张归因表就分不清哪枚是真伤**（票面 §丙 逐字警告的就是这一条：`B-3 那份读数今天不在他手上`）。 |

★**丙形那一支必须配一句事实判定（本腿现读，⛔ 不空写）**：
"**今天没有任何一把尺能区分'合法新增'与'误注册'**"这一句**不成立**——但**只成立一半**，逐枚具名：

- **能区分的那一半**＝`internal/tools/tasklist_deferred_236r3_teeth_test.go:153-161`（负向名集：在册的**恰好是** `task.list` 才红）
  ＋ `:164-177`（正向名集：`task.output`／`task.cancel` 少任一枚才红）。⇒ "误注册 `task.list`"这一形**今天有具名红、红句带现名册**；
  "合法新增第三枚"那一形**今天有 PASS 证明它不红**（该文件头注释 `:145-149` 写明 Deliberately NOT a count；`236-v1` §2 第 8 行的 PASS 读数是它的实测）。
- **不能区分的那一半**＝**除 `task.list` 之外**的任何新名字：没有任何一把尺能说出"这一枚新名字是**批过的**还是**没批的**"。
  `175r2:349-377` 只问"你答没答它是外来内容吗"，`task_cancel_221_test.go:134-141` 只问"说明书有没有许诺没注册的东西"，
  `236r3` 两枚名集钉只认三枚既定名字。⇒ **批准维（人工有没有批这一枚）今天零尺**——这是三形**共同的**残余，甲／乙／丙都没买到它。
- 半句话之外的事实：**枚数尺那把**确实能报"名册变了"，但它**说不出变的是哪一枚、批没批**（红句里 `%v` 只列名不列归属）。

---

## §3 要编排者裁的清单（★本腿一枚不裁、不预选形状；丙＝不做照实摆着）

1. **AC#0 选形**：三形代价表已交（§2）。本腿只补一句现量差别——**甲／乙在"顶红现有钉"这一栏都是 0 枚代码钉**，
   两形的真差别是：乙要一枚**具名解冻**（票 221 已 `-done`、尺在 `c44b30c4` 那一发里），甲**不要**；
   而买到"合法新增不再假红"的只有乙。⇒ 请裁。
2. **乙形那枚解冻的"算不算改判据"**：★现量摆明一个没人写过的洞——`task_cancel_221_legs_test.go:183` 那把尺**不在票 221 任何一格 AC 原文里**，
   `docs/evidence/s1/221-task-cancel-v1.md` 的逐 AC 判语也没引用它（本腿逐格读毕）。⇒ **"它在法律上是已勾判据的一部分，还是实现腿顺手拉的线"要编排者裁**，
   因为这决定解冻要不要写"改别人判据"那一档、以及 AC#4 的五样怎么落。
3. **甲形的写面边界**：加注释＝零解冻；改 `:184` 那句**文案**算不算动判据 ⇒ 同一枚问题在甲形这里的影子，请一并裁。
4. **甲形正控带来的那 1 枚必带仪器红**：`236r3` 仪器尺在任何换入 `task.go` 的 overlay 上都红（它自己注释逐字承认）。
   ⇒ 请裁"这一枚要不要写进本票的 AC#3 门禁名册口径"，否则下一位会把它当新常红。
5. **★票面 §现量 3／§现量 4 与凭据的**数字归属**冲突（本腿逐格对拉，未重跑，见 §1.5(a)）**：
   票面 §现量 3 写 `B-3`＝注册无关的 `task.note` ⇒ **rc=1／PASS=196／FAIL=10 枚具名红、其中含两枚 AC#1 拒绝钉走真注册路径**；
   凭据 `verdict.md:123`（`B-3` 那一格）逐字＝**rc=1／PASS=203／FAIL=3**，而 **196／10 枚／两枚拒绝钉**在 `:121`＝**`B-5`（注销 `task.cancel`）**那一格
   （`236-v1` §5 第 3 条 `:261` 同口径：`B-5`＝10／`B-4`＝4／`B-3`＝3）。
   ⇒ 请裁：票面 §现量 3 那三个数是**换错发**还是**另有名册**（本腿 ⛔ 不重跑、不裁）。
   另一处同源的序号漂：票面 Status 行与派单说"§2 **第 8 行**的附带读数"，本腿按数据行逐枚数下来，"枚数形尺今天会假红"那条附带读数的真身＝**第 10 枚数据行（`:123`）**，
   第 8 枚是 `B-4`／`B-5` 那格 ⇒ **内容都对得上、序号差两枚**，请一并裁引用口径（`第 N 格`／`第 N 枚数据行`／`第 N 行文件`）。
   ★这一处影响的是"那一发接手的人会看到几枚红"那一栏怎么写（3 还是 10），**不影响三形各自的写面与顶红名册**。
6. **射程要不要扩到 fs 那三枚同族雷**：现量＝`fs_test.go:190`／`:235`／`:263` 是同一形状的枚数尺，
   "合法新增一枚 **fs** 工具"会同时假红 3 枚；票 271 只写了 221 那一把。⇒ 请裁是本票扩射程、还是另立一票、还是就地登记残余。
7. **丙形若要选**：台账那行"今后任何腿接第三枚 task 工具前须先读本票"归编排者写（本腿 ⛔ 未动 `docs/**` 一字）。
8. **AC#1 那枚正控的载体**：票面要求 overlay、⛔ 不许原地编辑共享工作树；本腿零 `go` 命令所以**没验过任何台件的可行性**。
   ⇒ 派单时"哪一枚腿跑 go、它跑时 `cmd/wisp` 那枚并发腿要不要停"是排程问题，请裁。

---

## §4 判不动的地方（★本腿结构性看不见的，逐枚写满，⛔ 不推测填空）

**U1. 没跑过"真注册第三枚"那一发。** 本腿零 `go` 命令（§0.3）⇒ §2 三形"看到几枚红"那一栏里，**凡是本腿自己给的数只有"0 枚代码钉被顶红"这一类静态归因**，
所有"几枚红"的数都**来自 `236-v1` 的读数并已具名标注为未重跑**。⇒ **假红的确数（3／4／10）本腿不能证实，也不能证伪。**

**U2. 那 10 枚红的**完整名册**本件里没有。** `236-v1` §5 第 3 条只逐名给了三枚代表（`:181`、`:97`／`:100` 两枚拒绝钉、名集钉第二半），
其余是计数。**"其中几枚是真伤"这一问本腿答不了**：没有红名册原文＋每枚的归因，静态读码无法确定另外几枚是别的行为尺的连带还是独立伤。
⇒ 要它必须跑（禁）。

**U3. 一把尺**红过没有**看不见。** 枚数尺 `:183` 自 `c44b30c4` 起在 HEAD 上是绿的（本腿静态推的：名册 2 行、`!= 2` ⇒ 绿）。
但"它历史上有没有红过、被谁为什么改红过"要读 CI 日志或跑测试才知道 ⇒ 本腿零 `go`、也没打开任何 run 日志（`gh` 一枚未用）。
⇒ **"它可能就是故意拉的警报线"这一支持证据（它有没有在 CI 上真响过）本腿一条也没取到。**

**U4. ⛔ 禁止新增导出读口 ⇒ 拿不到"名册孤行"的另一种读法。** `236-v1` §2 第 5 行明确禁止、派单照抄。
所以"名册里那一枚新工具的**其它属性**（Decl／能力／路径参数）会不会连带顶红别的尺"这一维，本腿**结构上无法枚举**——
只能逐枚读现有断言的构造面（已做，见 2.0(b)），不能问"还有没有别的尺从别的角度数 task 家族"。

**U5. 并行工作树里的"别人正在写的尺"看不见。** `git status --porcelain` ＝ **750 行**（§0），
本腿只审了 `cmd`／`internal` 两棵树的**已提交态**（porcelain 0 行）。⇒ 若此刻有另一枚腿正在 `internal/tools/**_test.go` 里**新写一枚枚数尺**，
本普查**没算它**，三形代价表的"顶红谁"分母会漂。这一条只能靠编排者在派单时按包级互斥对齐。

**U6. "合法"这个字尺读不到。** §2.3 已具名：批准维零尺。本腿**无法判断**"哪天加第三枚"算不算合法——
`task.list` 的解冻要人工批（`PLAN.md §7`／D34 那一格、票 164 与票 225 各管一半，票面禁区段逐字），
`task.note` 那形今天**不在任何契约表里**。⇒ "警报线该不该响"的语义判据本腿给不出。

**U7. CI 端会不会因乙形少一枚红而"变松"看不见。** 本腿读了 `.github/workflows/ci.yml` 里那几条 `runtests.sh` 的 `-run` 定向步骤（`:404`／`:703`），
现量＝**ci.yml 里没有任何一处按名字引用 `Test221TaskCancelIsRegisteredAtItsFrozenLevel`**（`grep -rn` 具名尺在 `.go`／`.yml`／`.sh` 之外只命中 `docs/**` 与 `internal/tools/*.go` 注释，见 `logs/namerefs.txt`）。
⇒ 但"整包跑的红名册少一枚会不会让某道门（如'终态整包逐名红名册作差＝0'那条 AC#3 规矩）在下一片失去一个比对项"＝**要读别的片的门禁记录才知道，本腿未读、也没跑。**

**U8. 词面尺与名册尺的"分开发红"这一维本腿只能转述。** `236-v1` §2 第 8／9 行说 `B-3` 是唯一能把两把尺分开发出来的一发；
本腿**静态确认了它的机理**（`:279` 那句注释把 co-carry 数撑在 2 ⇒ `listMarked != 0` 恒成立；`task.go:23` 摘掉也不落到 0，见 §1.2 现数）
⇒ **机理本腿自己核过，读数没有**。这条区分对本票很重要（它是"枚数尺是不是唯一在管误注册的那一把"的依据）。

**U9. 三形的"以后维护成本"没有现量。** 例如甲形那条"新增第三枚那天仍需人工碰一次它"要几次、谁碰、走哪道工序——
仓里没有一枚尺数过这件事（本腿 grep 过 `解冻`／`常红`／`爆炸半径` 的落账形状，见 `logs/discipline.txt`，全是人写的 A## 台账、不是仪器）。
⇒ **人工碰一次的代价无法量化**，选甲／丙那一栏这里是空的，只能按定性摆。

**U10. `236-v1` 判语与本件读数的**同一性**只验到了"字节级"，没验到"读数级"。** 本腿按派单只当背景引用它，⛔ 未重跑那些 overlay 突变。
§1.5(b) 那把 md5 对拉**证到了**两枚腿取的是同一份产码字节（四枚重叠串逐字相同）⇒ 静态归因可以接在它那些读数后面；
但**它那些 rc／PASS／FAIL 的数本身**要跑才能校准。⇒ **本腿能做的是"确认我们对的是同一份文件"，不是"确认它数得对"。**

**U11. 三形各自"要不要解冻"这件事，本腿给的是**现量**不是**定性**。** §2.2 摆明：枚数尺不在票 221 任何一格 AC 原文里、
裁决表也没把它列进任一 AC 判语 ⇒ 本腿**无法判断**它是"已勾判据的组成部分"还是"实现腿自拉的线"（§3 第 2 条）。
这条一旦定错，乙形的解冻要么写成"改契约"（过重）、要么写成"没改任何 AC"（漏账），⛔ 两种都由编排者担、不由本腿选。

**U12. "有没有第三枚 task 工具在路上"这件事本腿没查、也不该查。** 那是路线图地界（票面禁区段：`⛔ 不许把"注册 task.list"顺手做进本票`，
`PLAN.md §7`／D34 那一格、票 164 与票 225 各管一半）。⇒ 三张表里"那一天"全是**假设态**，没有一枚真实排队工单支撑它；
若其实**根本没有**那种工单在路上，丙形的代价就比表里写的更轻（轻到多少＝要读路线图才知道，本腿 ⛔ 未读）。