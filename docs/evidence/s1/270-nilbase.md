# 270-nilbase — 票 270 落地腿 `270-r1` 交件证据（形＝甲，`A990` 裁③）

时刻：2026-10-15 21:01 +08（`date` 现读）。锚：起手 `HEAD=6a6b1b18`（＝`origin/dev`，本腿未推）。
写面（逐枚，越界零）：`internal/tools/subagent_197.go`（两处：`:256` 上方注释＝AC#5、`:290` 判空＝AC#1）
＋新件 `internal/tools/nilbase_270_test.go`（191 行）＋本件＋`.scratch/wisp/probes/270/r1/**`＋票 270 末尾 1 条 bullet。
⛔ 未碰：`cmd/wisp/**`、`:296`/`ParentTools` 半边、`ErrorClass` 那一维、`failclosed_236_teeth_test.go` 的字面量钉、任何 `- [ ]` 框／票名／Status。

## 0. 现跑锚（票面旧号已漂 +3，本腿逐枚重跑，⛔ 抄号）

| 锚 | 现量 | 对上 |
|---|---|---|
| `internal/tools/subagent_197.go:135` | `BaseOptions func() (agent.Options, bool)` | 是 |
| `internal/tools/subagent_197.go:256-257` | 一枚 guard 管两枚字段＋具名拒绝整串 | 是 |
| `internal/tools/subagent_197.go:281` | `base, wired := t.d.BaseOptions()` | 是（本腿改后＝`:290`） |
| `internal/tools/subagent_197.go:296` | `opt.Tools = newSubagentToolProvider(t.d.ParentTools)` | 是（只读，零改动） |
| `cmd/wisp/run.go:792` | `BaseOptions: func() (agent.Options, bool) { return rt.loopOpt, rt.loopOptSet },` | 是（只读） |
| `internal/tools/failclosed_236_teeth_test.go:95` | `lit236SpawnNoAssembly` 两行拼接常量 | 是 |
| `internal/tools/failclosed_236_teeth_test.go:208` | `Test236R2TaskSpawnRefusesWhenBaseOptionsIsUnwired` | 是 |

★AC#1 前提（`270-a2` 给的"那两处复用同一句"，本腿**自己再验一次**）：`:257` 与 `:299` 两行**逐字节相同**
（python 抽行比对＝`True`，两枚命中行号 `[266, 299]`），且该串与 `lit236SpawnNoAssembly` 相等
（`Test270R1AssemblySentenceIsTheOne236Pinned` 在盘上钉住，红则具名）。前提成立 ⇒ 甲形成立，未换文案。

## 1. AC#0 引用＋现跑复认（⛔ 抄 a1/a2 的数）

引用两件名（本腿只引，不重做普查）：`.scratch/wisp/probes/270/a1/census.md`（三形代价表：甲 0 枚具名红／
乙 3 枚／丙 0 枚但对 `-overlay` 失明）·`.scratch/wisp/probes/270/a2/pairing.md`（更正一：原 AC#1 判据恒真）。
现跑复认"甲形今天仍 0 枚具名红"：起手（未改产码、未加测试）`go test ./internal/tools/ -count=1 -v` ＝
**rc=0／PASS=235／FAIL=0** ⇒ 顶红枚数 **0**，与 `a1` 的读数同色。改动后正跑＝**PASS=239／FAIL=0**（239＝235＋本腿新 4 枚）。

## 2. AC#1 定向判据（突变台件，形＝物理改写＋sha256 还原）

方法学一条（承 `321-r1` 已入库的规矩，本腿采纳）：**⛔ `go test -overlay`**——overlay 只换编译期输入，
对任何"运行时读源码"的 AST／尺类台件失明；本腿用**物理改写＋sha256 还原**，`.bak` 留底**只建不删**。
台件＝`.scratch/wisp/probes/270/r1/mutate.sh`＋`guard_strip.py`（按条件文本定位、**要求恰 1 枚命中**，
命中数≠1 或紧随两行不是那句拒绝即 rc=3/4 拒绝突变；⛔ 行号锚，本腿自己的注释会移行）。
摘的形＝把 `:256` **整支 3 行删掉**（`if … || … {` ＋ return ＋ `}`），保留调用点。

基线（终态产码，未突变）sha256＝`86de44c4b79b35172c5ca4bd42862f3659e1d9268ab992f975c282eaa6b172e3`
·突变后（甲形＋摘 guard）＝`a9fdc98c29e3cee5af862bf5285b7bb68722e75547e7592d874860b5fd08132f`
·突变后（HEAD 产码＋摘 guard＝旧形）＝`16729a92f5fd32a71ca7f8f6505737bb9f0eb6240ee7d75c7796333c4c915a6b`
·两轮还原均 `restore verified: 86de44c4…`（＝基线同串），`git diff --numstat` 该件＝`20 0`（只本腿两处）。
中间留底：`subagent_197.go.bak`／`subagent_197.go.bak-r1`（注释改动前那份）／`subagent_197.go.head`。

| 读数（同发突变＝摘 `:256` 整支，`go test ./internal/tools/ -count=1 -v`） | rc | PASS | FAIL | 名册合计 | `^panic` 行 |
|---|---|---|---|---|---|
| **旧形**（HEAD 产码＋摘 guard）`mu-nofix-full.txt` | 1 | 235 | 4 | **239** | 0 |
| **甲形**（本腿产码＋摘 guard）`mu-fix-full.txt` | 1 | 238 | 1 | **239** | 0 |
| 正跑（未突变）`final-count1.txt` | 0 | 239 | 0 | 239 | 0 |

① 两发突变**都出全量名册 239**（＝起手 235＋新 4 枚），⛔ 没有 25/2 那种失明形；`^panic` 恒 0，
`FAIL github.com/CarlosShao/wisp/internal/tools` 之后**没有**"零读数"尾巴（日志末行仍是逐名 PASS/FAIL）。
② 走到具名拒绝：甲形突变下 `Test236R2TaskSpawnRefusesWhenBaseOptionsIsUnwired` **＝PASS**，
本腿 `Test270R1NilBaseOptionsRefusesWithTheAssemblySentence`／`Test270R1ThreeRefusalSentencesDoNotSubstitute` 亦 PASS。
③ ★红→绿现场验（同一支 guard、同一摘法，只差调用点判空）：

旧形红句（`failclosed_236_teeth_test.go:219`，连文案原文）：
`subagent_197.go:256 那一支的完整字面量没被原样说出（got "panic（guard 摘掉后不是拒绝而是崩）: runtime error: invalid memory address or nil pointer dereference", want "宿主没有给出派生用的装配（BaseOptions/ParentTools 未接线，fail-closed：不起第二套运行时）"）——这一支是 || 的 BaseOptions 半边；got 以「panic」开头就是那道 guard 已经不在了`

甲形同发＝`--- PASS: Test236R2TaskSpawnRefusesWhenBaseOptionsIsUnwired (0.00s)` ⇒ **红→绿成立，且是这一支的功劳**。
本腿两枚尺在旧形下的红句（证明尺有牙、不是恒真）：
`nilbase_270_test.go:60: 调用点没自己判空（摘掉上游 guard 后的旧形）："panic（guard 摘掉后不是拒绝而是崩）: runtime error: invalid memory address or nil pointer dereference"`
`nilbase_270_test.go:132: 装配未接线 那一支说的不是自己那句（got "panic（…nil pointer dereference", want "宿主没有给出派生用的装配（…不起第二套运行时）"）`

甲形突变下**只剩 1 枚具名红**＝`Test236R2TaskSpawnRefusesWhenParentToolsIsUnwired`（`||` 的**另一半**，
`A990` 裁③与票面 10-07 裁定明写射程外）。红句原文：
`failclosed_236_teeth_test.go:192: ParentTools 没接线时 task.spawn 回了非拒绝："子代理 …（父工具面没接线）已结束，状态 Settling。…（guard 不在时子代理照样派生，只是拿到一个空目录）"`
`failclosed_236_teeth_test.go:196: …那一支的完整字面量没被原样说出…这一支是 || 的 ParentTools 半边`
⇒ 与票面 10-07 的静态读法同形（延迟崩／空目录，不是当场 panic），本腿**没有**为它变绿去动那枚钉。

## 3. AC#3 三句互不充当替身（新尺 `Test270R1ThreeRefusalSentencesDoNotSubstitute`）

四枚输入、四句理由，逐枚**等值**钉＋两两**不等**钉＋反控：
① 装配未接线（`d.BaseOptions = nil`）＝`lit236SpawnNoAssembly`；② 池满（`TryAcquireSubagentSlot` 手动占满
`MaxConcurrentSubagents=4`，再派第 5 枚）含 `拒绝派生`／`上限`／`不是排队`；③ 没有父任务 id
（`spawnWithoutHostTaskID`，两个 id 字段皆空）＝`lit236SpawnNoHostTaskID`；④ 名册未接线＝`lit236SpawnRosterUnwired`。
`236-v1` 实测到的那一形（nil 名册被池满门回答"你排满了"）在此钉成**具名断言**：
`if strings.Contains(gotNoRoster, "上限") { … }` ⇒ 名册门与池门对调即红；另钉装配句里不许出现`上限`／`父任务`。
反控另有一枚 `Test270R1WiredBaseOptionsIsNotRefused`：装配齐全时**必须不**拒、且读数里不出现那句拒绝
⇒ 上面几枚不是一枚"逢派生即拒"的恒真尺。
池满那枚还钉 `InFlightSubagents`／名册行数不被被拒调用污染（承 `Test197FullPoolRefusesNextSpawnWithReadableReason` 的口径）。

## 4. AC#4 门禁（每条自落 rc，全部终态复跑于还原之后）

| 尺 | rc | 读数 |
|---|---|---|
| `go build ./...` | 0 | 静默 |
| `gofumpt -l <本腿两件>`（`D:\work\base\gopath\bin`） | 0 | 空列表 |
| `go vet ./internal/tools/` | 0 | 静默 |
| `sh scripts/d22scan.sh` | 0 | clean；且射程非空：bans #1-5 internal/=239／cmd/=41／ban #6 frontend/=116／ban #7 internal/tools/=25／ban #8 internal/=564（含注释与 `_test.go`）／cmd/=136 |
| `sh scripts/check-path-length-budget.sh --with-self-test` | 0 | VERDICT GREEN，over-budget=60＝roster=60，not in roster=0 |
| `go test ./internal/tools/ -count=2 -v` | 0 | PASS=478／FAIL=0（＝239×2） |
| `go test ./internal/tools/ -count=1 -v`（两发取交集） | 0 | PASS=239／FAIL=0 |
| 红名册与起手作差 | 0 | 起手＝`0 枚`／终态＝`0 枚`⇒ 差 0（既有红没被本腿算进、也没被顺手弄绿：起手本就全绿） |
| `git status --porcelain -- internal cmd` | 0 | 起手 0 行；提交后终态回到 0 行（本腿只提交自己两枚路径） |

环境：`GOTMPDIR='C:\Users\swq\AppData\Local\wisp-test-gotmp\270r1gotmp'`（C: 卷、反斜杠形）；
`PATH` 加 `third_party/sherpa-onnx`；机主实例 `wisp.exe` PID 21892 全程在跑、未被 kill、未被另拉。

## 5. AC#5 注释落点（`:256` 上方，原文）

`ONE branch, TWO duties (ticket 270): … it is also the only site that keeps the t.d.BaseOptions() call
below from being reached with an empty func value … Form 甲 (ledger A990, decision 3) put a nil check of
its own at that call, so deleting this branch now lands on the same named refusal instead of on that panic;
the sentence below stays the one this branch of the || speaks for both of its fields.`
只写现状（"今天哪一处先说话／摘掉后读数是什么形"），⛔ 未许诺任何未实现能力。
ban #9 自查：注释里出现的唯一路径 `internal/tools/failclosed_236_teeth_test.go` 盘上真存在；
⛔ 引用 `docs/DECISIONS.md`／`DEFERRED.md`／`TOOLS.md`／`STATE_MACHINE.md`／`contracts/`（`270-a1` 实测全 MISSING），
需要点名 SPEC-05 时写作 `docs/specs/SPEC-05-agent-core.md`（真名）。

## 6. 没做完的格子（逐枚具名，⛔ 静默）

1. **AC#2 不适用**：本票选甲、`:256` 未拆 ⇒ "半边摘＝恰 1 枚具名红"那格今天无从验（`ParentTools` 半边那枚红是**射程外**，不是本腿的牙）。
2. **`ParentTools` 半边**（`:296`）零改动：摘 guard 后它走"延迟到 `Execute`／空目录"那一形，甲形**没有**替它装判空；立不立票归编排者（票面 10-07 已按下）。
3. **AC#1 框／AC#3～AC#5 框一枚未翻**，票名／Status／票面既有文本一字未动；本件只是腿的自证，**对抗验收须由非实现者另做**（`SPEC-12 §4.3` #1/#3）。
4. **恒真自查留了一处弱牙并具名**：新尺**无法**在正跑树里直接观察 `:290` 那一句的文案（`:256` 永远先说话），
   所以"两处复用同一句"的可观察证明**只在摘 guard 那一发**成立（本腿已现场取到）；正跑树里本腿只能钉"两串相等"这一半。
   这就是这把尺的真实射程，⛔ 拿正跑的绿冒充 AC#1。
5. **CI 色未取**（本腿不推、不开 worktree）；`-count=3` 那格按票 236 AC#5 的口径留给编排者定夺（本腿交了 `-count=2`＋`-count=1` 两发交集）。

## 7. 台账／禁区核对

`A990` 裁③（选甲）·票面 10-07 裁定（射程收窄为 `BaseOptions` 半边）逐条对齐；⛔ 放宽 236 那批字面量钉（未编辑该文件一字，
`git status --porcelain -- internal/tools/failclosed_236_teeth_test.go` 空）；⛔ 新导出名（新注释＋包内判空而已，测试件符号皆包内小写）；
⛔ `ErrorClass` 那一维；⛔ 拆 `:256`；⛔ 与 AC#2 分批（甲形本就不拆，不存在中间态）；
⛔ 双引号串里的反引号／`$`（进 shell 的正文一律单引号或加引号 heredoc，本腿零事故）；临时件只建不删（`.bak`／`.bak-r1`／`.head`／六份日志全留底）。
