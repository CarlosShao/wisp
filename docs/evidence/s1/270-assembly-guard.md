# 票 270 — 非实现者终裁（腿 `270-v1`）：`task.spawn` 装配 guard 与调用点自保护（形＝甲）

- 腿＝**裁决者·只读**；实现者＝`270-r1`（`03da4941`，形＝甲／台账 `A990` 裁③；写面两枚：`internal/tools/subagent_197.go`＋新件 `internal/tools/nilbase_270_test.go`）
- 裁的落点＝HEAD `232f309e`；本件生成时刻 `2026-10-16 11:1x +0800`（`date` 现读）
- 票框尺（现跑 `grep -c '^[[:space:]]*- \[ \]'`／`'- \[x\]'` 对 `.scratch/wisp/issues/270-task-spawn-assembly-guard-is-also-the-only-nil-protection-for-line-278.md`）＝**6 枚未勾／0 已勾**（本腿一枚未翻）
- 台面＝`git archive HEAD | tar -x` 的**仓外导出树** `D:\work\workspace\wisp-270v1-adjudication-bench\repo`（`find … -name '*.go' | wc -l`＝**1098**，六件套在场）；突变全部打在**物理字节**上，⛔ `-overlay`；逐发还原对 `git cat-file -p HEAD:internal/tools/subagent_197.go` 比 sha256
- 读数件＝`.scratch/wisp/probes/270/v1/logs/`（`read-*.txt` 逐发四数＋红句原文＋还原哈希、`gates.txt`、`gofmt-and-surface.txt`）；台件＝`mutate270v1.py`（按条件文本定位，非恰 1 命中即 rc=3 拒绝突变，⛔ 行号锚）＋`run-bench.sh`

## 一、现量行号（号只当读数；尺＝`git grep -n '<条件文本>' 03da4941`＋Read 工作树同一枚文件；`git log --oneline 03da4941..HEAD -- internal/tools/subagent_197.go internal/tools/nilbase_270_test.go`＝0 行⇒落地后零改动）

| 票面的号 | 现量号 | 那一行的原文（判据＝该行真含这句话／这个条件） |
|---|---|---|
| `:256` 的 guard | **`:265`** | `if t.d.BaseOptions == nil \|\| t.d.ParentTools == nil {`（拒绝文案在 `:266`） |
| `:256` 上方那句 AC#5 注释 | **`:256`–`:264`** | `:256` = `// ONE branch, TWO duties (ticket 270): it refuses an assembly the host never`；`:261` = `// Form 甲 (ledger A990, decision 3) put a nil check of its own at that call,` |
| 甲形新加的那支（派单写的 `:290`） | 注释 **`:290`–`:296`**、判断 **`:297`**、还位 **`:298`**、拒绝 **`:299`** | `:297` = `if t.d.BaseOptions == nil {`；`:298` = `release()`；`:299` = 与 `:266` 逐字节相同的那句装配拒绝 |
| `:278` 那次调用 | **`:301`** | `base, wired := t.d.BaseOptions()`（票面"就在 guard 之后两行"已过期：`:266`→`:301`＝35 行） |
| 生产侧唯一装配点 | **`cmd/wisp/run.go:796`** | `BaseOptions: func() (agent.Options, bool) { return rt.loopOpt, rt.loopOptSet },`（注册在 `:787` `tools.BuiltinSubagentEntries(...)`；票面写的 `:789` 漂了 7 行） |

## 二、五格判语（只用三种；⛔"附条件成立"）

| 格 | 判语 | 承重读数（哪把尺＋射程＋量的还是算的） |
|---|---|---|
| **AC#1** | **成立**（残余两处具名见下） | `read-strip-guard.txt`／`read-strip-both.txt`／`read-guard-off+jia-panic.txt`／`read-guard-off+jia-silent.txt`／`read-guard-off+jia-no-release.txt`／`read-strip-jia.txt`／`k0-run-single.txt` |
| **AC#2** | **不适用**（票面自写"只在选乙时适用"，本票走甲⇒`:265` 未拆，这一格无人钉；实测见下） | `read-guard-keeps-baseoptions-half.txt`／`read-guard-keeps-parenttools-half.txt` |
| **AC#3** | **成立** | `read-strip-guard.txt`（构造那一发）＋`read-guard-off+jia-borrowed-sentence.txt`（替身形必响）＋`read-strip-both.txt` |
| **AC#4** | **成立**（含两枚具名代用尺／台面归因） | `gates.txt`＋逐发 `MATCH=yes`＋`read-base.txt`↔`read-restore-green.txt` 作差 0 |
| **AC#5** | **成立**（含一处过期文字缺陷，具名见 §四） | `:256`–`:264` 与 `:290`–`:296` 原文对同文件 `:125`–`:127`、`failclosed_236_teeth_test.go:201`–`:207` 逐字读 |

### AC#1 的三发（AC 原文＝"`:278`〔现量 `:301`〕不再无条件调用；`BaseOptions == nil` 走到具名拒绝，⛔ panic、⛔ 静默退化"）

1. **复跑该用例**（量）：`go test ./internal/tools/ -run Test270R1NilBaseOptionsRefusesWithTheAssemblySentence -count=1 -v` ⇒ `--- PASS`、rc=0。
2. **A990 裁③ 那一发＝定向判据**（突变 `strip-guard`＝物理删掉 `:265`–`:267` 整支，突变后 sha256=`a9fdc98c…`，与实现者报的甲形突变串**独立复现同串**）：
   `RUN=325／PASS_TOP=241／FAIL_TOP=1／SKIP=0／^panic=0`、rc=1 ⇒ 名册**全量**（＝起手 325，⛔ 25/2 那种失明形）；
   `Test236R2TaskSpawnRefusesWhenBaseOptionsIsUnwired` 在这一发里 **`--- PASS`**（＝要求的那一发"红→绿"），`Test270R1NilBaseOptions…` 亦 PASS ⇒ `:301` 之前那支真接住了；唯一红＝`Test236R2TaskSpawnRefusesWhenParentToolsIsUnwired`（`||` 的另一半，射程外）。
3. **去牙＝物理摘掉那句 nil 判断**（`strip-both`＝同时摘 `:265` 整支与 `:297`–`:300` 那一支；sha256=`0767820b…`）⇒ **必红**，四枚具名：
   `failclosed_236_teeth_test.go:219`（got `"panic（guard 摘掉后不是拒绝而是崩）: runtime error: invalid memory address or nil pointer dereference"`）、
   **`nilbase_270_test.go:60`** = `调用点没自己判空（摘掉上游 guard 后的旧形）："panic（guard 摘掉后不是拒绝而是崩）: runtime error: invalid memory address or nil pointer dereference"`、
   **`nilbase_270_test.go:132`** = `装配未接线 那一支说的不是自己那句（got "panic…"）`、`failclosed_236_teeth_test.go:196`（ParentTools 半边，射程外）。
   这一发的 `RUN=325／PASS_TOP=238／FAIL_TOP=4／SKIP=0／^panic=0` ⇒ 摘掉保护后**名册仍全量**，即票面 :26 那条"拆掉任一半 guard 之后不许出现零读数"**成立**。
4. **反形检查**（票面问的两形，各配一发，均在 `strip-guard` 世界上才可达）：
   - 把拒绝换成 `panic`（`guard-off+jia-panic`）⇒ **4 枚红**（`nilbase_270_test.go:60`／`:132`、`failclosed…:219`、ParentTools 那枚）⇒ ⛔ panic 形**不放行**。
   - 换成"静默当成没孩子"（`guard-off+jia-silent`，`IsError:false`＋一句自造话）⇒ **4 枚红**（`nilbase_270_test.go:63`/`:66`/`:132`、`failclosed…:216`/`:219`）⇒ 静默形**不放行**。
   - ⇒ 这一族尺对"甲写歪成 panic／静默"两形**都敏感**；**具名不敏感的那一形＝"只摘甲那一支、guard 不动"**（`strip-jia`，sha=`a92f8161…`）：`325／242／0／0`、rc=0 ⇒ **零具名红**。这把尺对那一发不敏感。成因＝甲形本身（`A990` 裁③ 明写"`:256` 的 guard **保持现状**"），不是实现者偷工：`guard` 在前，`:297` 今天**一枚调用者都没有**（`read-strip-jia.txt`＝它的死支证明）；实现者在票尾与证据 §6 第 4 条**自报过同一处**。⇒ 计入残余，⛔ 推翻 AC#1（AC#1 问的是"`:301` 还无条件吗＋摘掉守卫会不会崩/失明"，两问都答了"是/不会"）。
5. **`release()` 真在被摘的那条路上执行了吗**（⛔ 只看注释）：`guard-off+jia-no-release`（摘掉 `:297`–`:300` 里那一行 `release()`，其余不动）⇒ `FAIL_TOP=2`，其中
   **`nilbase_270_test.go:74`** = `被拒的这一枚占着池位：in-flight = 1, want 0（判空那支必须先把位还回去）`＝**计数差**读数 ⇒ 还位这一步是被钉住的（它在 guard 尚存的世界上不可达，与第 4 条同一残余）。

### AC#2 — 不适用（本格永远无人钉，说实话）

- 票面 AC#2 自写"（只在选乙时适用）：半边摘 ⇒ **恰 1 枚**具名红，两半各一形"。本票落**甲**＝`:265` 那枚 `if` 一字未拆（`git show --numstat`＝`20 0`，只有新增行），所以**这一格没有交付物可裁**⇒ 判**不适用**（⛔ 翻框、⛔ 判成立）。
- 但"半边摘"这一形今天**是可测的**，我测了（各一发物理突变，⛔ 行号锚）：
  - `guard-keeps-baseoptions-half`＝摘掉 `||` 的 **ParentTools 半边** ⇒ `325／241／1／0`＝**恰 1 枚**具名红（`failclosed_236_teeth_test.go:192`/`:196`，got 是"子代理 …（父工具面没接线）已结束，状态 Settling"那一形）。
  - `guard-keeps-parenttools-half`＝摘掉 `||` 的 **BaseOptions 半边** ⇒ `325／242／0／0`、rc=0、`^panic=0`＝**0 枚红**。⇒ ★这一发就是"甲形下 AC#2 那一格永远无人钉"的**实测形状**：`BaseOptions` 半边被摘掉后**没人喊**，因为 `:297` 接住了（安全上正是甲要的效果；名册全量、走具名拒绝）。
- ⇒ 具名给编排者：如果机主要的是"摘任一半必恰响一枚"，那**只能靠乙形**（拆语句＋各自文案），甲形在定义上把 BaseOptions 半边换成了"无声但安全"。这一格今后也不会有人来钉，除非另批。

### AC#3 — 三句互不充当替身

- 名册尺（现跑 `git grep -n`，射程＝`internal/tools/subagent_197.go`）指到**四个各自的产出点**：
  `:254` 名册未接线／`:266` 装配未接线（guard 那一支）／`:273`–`:274` 宿主没给任务 id／`:285` 池满（`fmt.Sprintf` 的"同时在跑的子代理已达上限 %d 枚"）。装配那句＝`:299` 复用 `:266` 的**同一理由**（同一批字段），⛔ 是别的门那句的替身。
  测试侧字面量名册：`failclosed_236_teeth_test.go:93`/`:95`/`:98`＋`nilbase_270_test.go:34`（等值钉 `Test270R1AssemblySentenceIsTheOne236Pinned`，base 一发 PASS）。
- **构造那一发**（⛔ 靠正跑的绿冒充）：`strip-guard`＝guard 摘掉后，装配案唯一能说话的产出点只剩 `:299`；那一发里 `Test270R1NilBaseOptions…`（等值 want＝`lit236SpawnNoAssembly`）与 `Test236R2…BaseOptionsIsUnwired` 双双 **PASS** ⇒ 屏幕上出现的就是装配那句。
- **替身必须响**：`guard-off+jia-borrowed-sentence`＝把 `:299` 那句换成**名册门**那句 ⇒ 3 枚响，红句逐字
  `nilbase_270_test.go:66: 空函数值没有走到那具名拒绝（got "任务名册未接线（fail-closed：拒绝派生子代理，派生了也没有地方登记它）", want "宿主没有给出派生用的装配（…）"）`、
  `nilbase_270_test.go:132`（同 got）、`failclosed_236_teeth_test.go:219`。⇒ "装配案被别的门回答"这一形**会红**。
- 236-v1 实测那一形（nil 名册被池满门回答）由 `Test270R1ThreeRefusalSentencesDoNotSubstitute` 钉成具名断言（`nilbase_270_test.go:160`/`:163`），且它在 `strip-both` 那一发**真红**过 ⇒ 有牙。⇒ **成立**。

### AC#4 — 门禁我自己现跑（五发逐枚 rc＋还原自证）

`go build ./...` **rc=0** ／ `go vet ./internal/tools/` **rc=0** ／ `sh scripts/d22scan.sh` **rc=0**（clean 且射程非空：ban #7 `internal/tools/`=25 枚产码文件、ban #8 `internal/`=573、`cmd/`=139、`frontend/`=116、`design/`=30）／
`sh scripts/check-path-length-budget.sh --with-self-test` **rc=0**（`tracked paths=14113  over-budget=60  covered by roster=60  not in roster=0` ⇒ VERDICT GREEN）／
`go test ./internal/tools/ -count=1` **rc=0**（同台件 `-v` 四数＝`RUN=325／PASS_TOP=242／FAIL_TOP=0／SKIP=0`）。
**还原自证**＝九发突变逐发末尾 `restored_tree_sha256＝repo_HEAD_blob_sha256＝86de44c4b79b35172c5ca4bd42862f3659e1d9268ab992f975c282eaa6b172e3`、`MATCH=yes`；被种文件只有那一枚（`internal/tools/subagent_197.go`），仓内工作树 `git status --porcelain -- internal/tools/subagent_197.go internal/tools/nilbase_270_test.go`＝**0 行**（本腿零动产码/测试件）。
起手红名册 0 枚↔终态（还原后再跑一发 `read-restore-green.txt`）0 枚 ⇒ **作差 0**。
两处具名（⛔ 藏）：① 长度门那一发**在导出树里回 rc=2**＝`not inside a git work tree - refusing to guess a denominator`（脚本 `:137`）⇒ **台面故障、⛔ 计入判据**，改在仓根跑（先自证 `git diff --cached --name-only | wc -l`＝0⇒分母＝HEAD 集合）；② `command -v gofumpt`＝**NOT_FOUND**⇒ `gofmt -l` 是**代用尺**（比 gofumpt 松），"两件零输出"这一半⛔ 能推广到 gofumpt，本腿⛔ 拿它冒充。另：台面第一枚树放在 `C:\Users\swq\AppData\Local\Temp\…` 时 `TestTicket252R1AlignmentAddsNoAuthorization` 红于"短名回显与长名同串·前置条件缺失"＝台面归因，迁到带 8.3 别名的父目录后转 PASS（`gates.txt` 逐字留了那两行）。⇒ **成立**。

### AC#5 — 注释只写现状、⛔ 许诺未来

- 要求原文＝"`:256` 处写清'这支同时是 `:278` 的 nil 防护'，⛔ 只写现状，不许在注释里许诺未实现的能力"。
  交付＝`:256`–`:264` 那一块，逐句可验：`:257`–`:260` 写了"一枚分支两桩差事"（拒绝＋nil 防护）；`:261`–`:263` 写了甲形把判断放进调用点。
  **没有一句许诺未实现的东西**（无"将来会拆成两支""ParentTools 也会补上""以后加 ErrorClass"之类）。三句现状级断言我都现场验为真：
  ① "deleting this branch now lands on the same named refusal instead of on that panic"＝`strip-guard` 一发（BaseOptions 那枚**红→绿**、`^panic=0`）；
  ② "The guard above fires first in every wiring this package builds today, which is why this check has no ordinary caller"（`:292`–`:293`）＝`strip-jia` 一发（摘掉甲那一支**零具名红**⇒ 今天确实无人可达）；
  ③ "The Text is the same sentence as that guard's"＝`Test270R1AssemblySentenceIsTheOne236Pinned` PASS＋我的突变尺要求装配句出现次数恰为 1／2 才放行。
- **打脸检查**（引两枚原文，同文件）：`:257`–`:258` "it is also **the only site** that keeps the `t.d.BaseOptions()` call below from being reached with an empty func value"
  对 `:261` "Form 甲 … **put a nil check of its own at that call**" 与 `:294`–`:295` "**this line is the only thing left** that cannot be called empty"——
  同一次落地把两枚"only"都写进了同一枚文件，二者互相顶（`strip-jia` 与 `guard-keeps-parenttools-half` 两发共同证明：两处**都**站在 nil 与调用之间）。⇒ 这一句是**过期文字**，归更正（⛔ 由我改）。  另引一枚**票 236 侧的旧注释**（同包不同文件）：`failclosed_236_teeth_test.go:202`–`:203` "and the half where 'guard gone' is a CRASH, not a fall-through, because `:278` calls `t.d.BaseOptions()` **unconditionally two lines later**"——这句在甲形落地后**同样过期**（现量 `:301`，且它已⛔ 无条件）；它⛔ 本腿写面，登记给编排者。
- 判语＝**成立**。具名一句：AC#5 的现文只要求"写清两桩差事＋只写现状＋⛔ 许诺未实现"，我按这三条验；若机主要的是"注释里任何两句⛔ 互相顶"，那是**改判据**（须人工批准），本腿⛔ 自造尺子。

## 三、能力类必答：`task.spawn` 在生产链上的调用点枚数

- 尺（原样）＝`git grep -n '"task.spawn"' HEAD -- '*.go'` ⇒ **126 命中**。★其中 **104 命中在 `.scratch/wisp/probes/**` 的存档突变副本里**（别的腿的历史 `.go` 快照，⛔ 代码——这就是票 324 那种"名册尺同名假命中"的形状，⛔ 去掉它会得出"126 枚调用点"）。
- 真代码内 **22 命中**＝产码 **2 枚**（都在 `internal/tools/subagent_197.go`：`:187` `func (subagentSpawn) Name() string { return "task.spawn" }`、`:583` `const subagentHiddenTool = "task.spawn"`，后者的两枚用法点 `:595`/`:604`＝"孩子的目录里摘掉这枚工具＋有人指名要它就拒"）＋测试 **20 命中**。
- 生产链＝**1 枚装配点**：`cmd/wisp/run.go:787` `tools.BuiltinSubagentEntries(tools.SubagentDeps{…})`，其中 `:796` 那枚 `BaseOptions` 函数字段是本票那两处 nil 判断的**唯一真数据源**（⛔ `go build` 的 rc=0 当"接上了"；尺＝上面那把 grep＋`git grep -n 'SubagentDeps'` 的非测试面只剩 `cmd/wisp/run.go`／`internal/tools/bridge.go`／`internal/tools/subagent_197.go`）。
- ⛔ 全仓存在第二处产码直接以字符串调用 `task.spawn`；孩子侧那枚 `ParentTools` 半边（`:316` `newSubagentToolProvider`）本票零改动（射程外，仍候裁）。

## 四、交付缺陷具名（⛔ 我动一字，归编排者/机主处置；五格判语都⛔ 依赖这一节任一条）

1. `:257`–`:258` 那句 "the **only** site" 与 `:261`/`:294`–`:295` 互相顶（见 AC#5）⇒ 过期文字。
2. 票面 4 枚行号锚全部过期：`:256`→**`:265`**、`:278`→**`:301`**、`:135` 字段声明现量在 `:135`（未漂）、`cmd/wisp/run.go:789`→**`:796`**；票面那句"就在 guard 之后两行"由 10-07 更正二登记过，今天实际是 **35 行**。
3. **零具名红的两枚形状**（甲形固有）："只摘 `:297` 那一支"与"摘掉 `||` 的 BaseOptions 半边"两发都 `FAIL_TOP=0` ⇒ 甲形那支**今天没有自己的牙**，只有摘掉 guard 那一发能验它；`A990` 裁③ 就是这样设计的，⛔ 是实现者的漏，但**下一批若要拆语句必须连带补牙**。
4. `nilbase_270_test.go:79`–`:81` 那枚"准入次数＝0"钉（AC#1 的"⛔ 退化成另起一套运行时"半边）**在我全部十发里一次都没红过**：nil 函数字值无法被一行物理突变填成可用装配，而更早的 `t.Fatalf`（`:60`）让那一发跑不到该断言 ⇒ 这一半句的可观察证明**只剩文案等值钉**，具名弱牙（实现者自报的那条是文案不可观察，与本条是两枚不同的牙）。
5. `ParentTools` 半边（`:316`）零改动；`guard-keeps-baseoptions-half` 那一发的红句证明它是**延迟崩/空目录**形（got "子代理 …已结束，状态 Settling"）⛔ 当场 panic ⇒ 是否单立一票仍候裁（票面 10-07 裁定＋`A990` 裁③ 都按下未立）。

## 五、交件自证三把尺

1. **提交号＋逐枚 `git show --stat`**：本腿共 4 笔（台件与读数两笔＋本件一笔＋票尾一笔），每笔只含派单枚举的三处写面：`docs/evidence/s1/270-assembly-guard.md`、`.scratch/wisp/probes/270/v1/**`、票 270 末尾**一条** bullet。⛔ push、⛔ amend/reset/rebase/stash/checkout ./clean；`git add`/`git commit` 两处逐枚显式 pathspec。
2. **证据件绝对路径＋`wc -c`**：见本腿交回口径。
3. **去牙必红／还原必绿**（含还原哈希）：去牙＝`strip-both`，`RUN=325／PASS_TOP=238／FAIL_TOP=4／SKIP=0`，红名 `Test236R2TaskSpawnRefusesWhenBaseOptionsIsUnwired`＋`Test270R1NilBaseOptionsRefusesWithTheAssemblySentence`＋`Test270R1ThreeRefusalSentencesDoNotSubstitute`＋`Test236R2TaskSpawnRefusesWhenParentToolsIsUnwired`，红句原文见 §二 AC#1 第 3 条，报到的红行＝`nilbase_270_test.go:60`/`:132` 与 `failclosed_236_teeth_test.go:196`/`:219`；还原＝同一文件写回后 `sha256=86de44c4b79b35172c5ca4bd42862f3659e1d9268ab992f975c282eaa6b172e3`＝`git cat-file -p HEAD:internal/tools/subagent_197.go` 同串（`MATCH=yes`），还原后台面再跑一发＝`read-restore-green.txt`：`325／242／0／0`、rc=0。

**没裁完的格具名**：`AC#0`（三形代价表）**⛔ 本腿射程**——派单只交我 AC#1–AC#5＋能力类；AC#0 的交付物是 `270-a1/census.md`＋`270-a2/pairing.md` 两枚只读普查件，本腿**未裁**（⛔ 缺读数：它没交"拆完之后剩几枚零读数"的实测，票面 :44 自己具名过这一格欠账）。另：本腿**没跑** `go test ./cmd/wisp`（票 314 死亡区＋写腿 `220-r2` 在那面）与 `internal/tools` **整包以外**的任何面——派单明写 ⛔。
