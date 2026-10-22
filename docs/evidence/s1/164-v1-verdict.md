# 票 164 终裁腿 `164-v1` 裁定表 —— AC#2b：那枚用例是"守卫"还是"恒真装饰品"，以及改署名有没有偷动断言

- 裁腿 `164-v1`＝**非实现者**（裁决者≠实现者）。被裁 `164-r1` 四笔：`62a1521c`（唯一 `.go`）／`ec8786e1`／`0786b5c6`／`1146f803`（后三笔 docs-only）。
- 读数时刻 2026-10-22 22:2x–22:4x +08。**HEAD 在我这一程中由 `02b3882c` 前移到 `b4102f1e`**（别的腿提交）；两枚被裁 blob 在两个 HEAD 上逐字相同（`bridge.go` `ea7a739a…`、用例 `f409304b…`）⇒ 全部读数仍钉在同一对象上，不是跨树拼的。
- 写面：本件（+ 仓外 `D:\tmp\wisp\20261022\164v1\**`：9 支可重跑脚本与全部原始日志）。⛔ 未落任何 `.go`、⛔ 翻框、⛔ 改票面／台账、⛔ push。

## 0 总裁定

★**`AC#2b` 能翻（可勾）**。两条判据我都自己重造、不采信腿自报：**署名**已从"未修码读数"改成"未知工具拒绝形状＋D37 分类的守卫"（§2 引文与文件 `:11-49` **逐字相同**，2403 字节／39 行／`diff` 0 行）；**可证伪性**两向都成立（M1 红在 `:73`、M2 红在 `:76`、原样绿）。**"断言一字未动"独立核为真**（AST 语句级 diff 0 行，带正控）。三条残余见 §9，**都不触 AC#2b 的判据**。

## 1 项① 票面前提新鲜度（尺：票面自己的命令，对 HEAD blob 复量）

★它报的三条**逐条复核，三条全部成立**：

- **① 票面 `:4`"`task.output` 全仓零命中"已不成立** ⇒ **成立**。尺＝票面原命令 `grep -rlE '"task\.output"' --include=*.go internal/ cmd/` ⇒ **17 枚**（不是 0）。⚠**这 17 枚是"引号字面量 × internal/+cmd/"这一把尺的数**：同一字面量换成"任意出现"⇒ **32 枚**（internal/+cmd/ 与全树去 `.scratch` 都是 32）；**全树含 `.scratch` ⇒ 162 枚**。它报"17 枚"数字对，但**必须带尺**，否则换把尺就变成另一个数。
- **② 票面 `:56`（与 `:39` 定案⑤同句）"`RunAsync` 生产零调用点"已不成立** ⇒ **成立**。非测试真调用点＝**2 枚**：`cmd/wisp/run.go:1176`、`internal/tools/subagent_197.go:372`（其余非测试命中是注释：`internal/agent/loop.go:327`、`internal/tools/bridge.go:880`、`task_backfill.go:12,14,27`）。
- **③ 票面 `:52`／`:59` 引的 `bridge.go:251`／`:247-253` 已漂** ⇒ **成立**。现读：`:247-253` 是 `Execute` 的文档块（`:251` 逐字＝`// The returned error is reserved for host-internal faults, because the loop`），拒绝支在 **`:272-278`**、文案本体 **`:276`**、分类 **`:277`**。
- **我另量到的两处腐坏（它没报）**：**④** 票面 `:5`／`:10`"`task.list`／`task.cancel` 作为工具在生产注册表里并不存在"**半假**——`task.cancel` 已实现并注册（`task.go:625` `Name() "task.cancel"`；`BuiltinTaskEntries` `task.go:600`；生产 `run.go:557`），只有 `task.list` 仍缺（仅注释／DEFERRED）。**⑤** 票面 `:18`"生产侧实现＝0 枚"已不成立（两枚 `task.*` 有实现）。
- 仍然新鲜的：票面 `:17` 的 awk 原命令 ⇒ `task.cancel`／`task.list`／`task.output` **各 1 次**；`:59`"3 勾／3 未勾"⇒ `grep -c '^- \[x\]'`＝**3**、`'^- \[ \]'`＝**3**、缩进不敏感总数＝**6**；票面引的 6 个锚点全部存在（`e4d8e047`／`b5adc34f`／`58c400b7`／`475cc885`／`b1d9cfda`／`f206e9f`）；AC#3 三枚用例名各 1 处定义。

## 2 项② 核心判点：我自己造的三发（尺：`go test -count=1 -v -run '^TestTaskOutputAC2BeforeLegIsUnreachable$' ./internal/tools/`）

每发先证落地（`grep -n`＋`go build ./internal/tools/` rc=0）再读数，读完立刻单独还原（`cp` 回备份，⛔ 不用 `git checkout`）。

```
[原样] rc=0  --- PASS: TestTaskOutputAC2BeforeLegIsUnreachable (0.01s)
       :66 unknown-tool refusal: IsError=true ErrorClass="tool" Truncated=false Text="未知工具 task.output，可用工具见 list_tools"
[M1]   bridge.go:276 "未知工具 "→"无法识别工具 "；落地 grep 276；sha256=89210b00…c09c7a；build rc=0；test rc=1
       :73 expected the unknown-tool refusal (Bridge.Execute, Lookup-failure branch), got: "无法识别工具 task.output，可用工具见 list_tools"
       --- FAIL: TestTaskOutputAC2BeforeLegIsUnreachable (0.01s)
[M2]   bridge.go:277 observe.ClassTool→observe.ClassInternal；落地 grep 277；sha256=589b537e…27e51c；build rc=0；test rc=1
       :76 error_class = "internal", want tool (an unknown tool is self-correctable, D37)
       --- FAIL: TestTaskOutputAC2BeforeLegIsUnreachable (0.01s)
[还原×2] sha256 回到 a0c0736c…ccfdd；git status --porcelain -- bridge.go 为空；git diff --stat 为空；再跑原样 rc=0 / PASS
```

★两发变异的 sha256 与它报的**逐字相同**（我独立用 `sed` 行定位造，落点字节一致）⇒ 它确实做了它说的那两发，不是抄的。
★**判语：那枚用例是"守卫"，⛔ 恒真装饰品**——它钉的是桥 `Bridge.Execute` 的**活行为**（文案形状＋D37 分类），不是树状态；桥下面挂的注册表是 `helpers_test.go:45` 的 `BuiltinFSEntries`（fs 家族），故它对"后台任务今天能不能读"这一问仍**不可证伪**——这正好是它该有的射程，不是缺陷。

## 3 项③ "断言一字未动"我的独立核法（⛔ 只数 `Fatalf`）

- **甲法（复现它的去注释 diff）**：去注释后逐行 `diff` ⇒ **4 行**，全是消息串：`t.Logf` 标签 1 行＋`t.Fatalf` 措辞 3 行；**条件、实参、语句顺序零变化**。
- **乙法（我自造的 AST 语句级探针，仓外 `D:\tmp\…\astcmp\main.go`）**：用 `go/parser` 读两版，**注释在 AST 里结构性不存在**；只把 `t.Logf/t.Fatalf/t.Errorf` 的**第一枚格式串**遮成 `<format>`（其余字面量逐字保留，`strings.Contains(out.Text,"未知工具")` 这种**看得见**）。读数：**before 29 条语句／after 29 条／`logic_diff_lines=0`**。
- **正控**：同一探针把用例里 `out.ErrorClass != "tool"` 改成 `!= "internal"` ⇒ `DIFF @18 IF (out.ErrorClass != "internal")`，rc=1 ⇒ 探针能看见断言逻辑，0 行不是空转。
- **结论**：`62a1521c` 对断言逻辑**一字节未动**；动的只有 4 处消息串＋注释。它自报的计数（`Fatalf`4／`Errorf`2／`if `7／`ErrorClass`4）我按**行**复量为同值（`strings.Contains` 按行 2／按出现 3，它报 2＝按行；不构成假话，但"计数"这把尺本身验不出逻辑，所以我不拿它当结论）。

## 4 项④ 署名是否真改到位（尺：md §2 引文 vs 文件 `:11-49`）

`sed -n '11,49p'` 与 md `:51-89` 两块 ⇒ **IDENTICAL，diff 0 行**，两边都是 **39 行／2403 字节**（它说第一次引文末三行换行不同、已改 ⇒ **改后确已逐字相同**）。§2 引文不是"大意"，是可逐字复跑的引文。

## 5 项⑤ 零产码复核（尺：blob SHA-1 vs 工作树；⚠ 它报的是 SHA-256）

- `internal/tools` 全部 `.go`＝**94 枚**（产码 25＋测试 69）；逐枚 `git rev-parse HEAD:<p>` 对 `git hash-object <p>` ⇒ **SAME 94／DIFF 0**；且 94 枚 **CR 全为 0**，故裸 `hash-object` 不受 CRLF 干扰（若有 CRLF 我会改报 `--path` 过滤后的值）。
- `bridge.go`：blob `ea7a739a…` ＝ 工作树 `hash-object` `ea7a739a…`；它件里报的 `a0c0736c…ccfdd` 是 **64 位＝SHA-256**，我用 `sha256sum` 复量＝**同值** ⇒ **它报的是 SHA-256，⛔ SHA-1**（编排者本轮踩的那一脚，这里没再踩）。
- `git status --porcelain -- internal/ cmd/ scripts/ .github/ docs/` ＝ **空**。

## 6 项⑥ 门禁（我自己现跑；四数分开报、⛔ 相减）

| 门 | 命令 | rc | 末行／四数逐字 |
|---|---|---|---|
| 构建 | `GOFLAGS= go build ./...` | **0** | 输出 0 行（末行为空） |
| 静态 | `go vet ./internal/tools/` | **0** | 输出 0 行 |
| 单包 | `go test -count=1 -v ./internal/tools/` | **0** | `RUN=329`／顶层 `PASS=246`／`FAIL=0`／`SKIP=0`；末行 `ok  github.com/CarlosShao/wisp/internal/tools  20.560s` |
| 族 | `go test -count=1 -run TaskOutput ./internal/tools/` | **0** | `RUN=4`／`PASS=4`／`FAIL=0`／`SKIP=0` |
| D22 | `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations; …ban #8 cmd/=148…` |
| 路径 | `sh scripts/check-path-length-budget.sh` | **0** | `…VERDICT GREEN - every over-budget tracked path is rostered…`（over=61） |

**格式甲乙两尺**（裁定件＝`docs/evidence/s1/pipe-format-gate-reads-tracked-blob-not-worktree.md`）：
- **甲（判据，blob 树）**：`git archive HEAD`（`b4102f1e`）→ 仓外 `export-b4102f1e/`（无 `.git`）→ `gofumpt -l .` ⇒ **20 行，全部 `.scratch\…`**；`grep -vcE '^\.scratch[\\/]'`＝**0** ⇒ 被裁那枚**不在名单**。**正控**：往导出树丢一枚未格式化 `.go` 再跑 ⇒ 非 `.scratch`＝**1**（`zz_probe_unformatted.go`）⇒ 尺能看见坏 `.go`。
- **乙（只作诊断，工作树）**：`gofumpt -l internal/` ⇒ **2 行**，都不是它的：`internal\agent\approval\pending_read.go`（**CR=131**）、`internal\agent\tools.go`（**CR=225**）。⚠**它 §5 写"各 225 个 CR"是错的**：131／225，不是各 225。它对单文件跑 `gofumpt -l` ⇒ 零输出、rc=0。
- 顺带复核它自报的盲尺：`grep -c "$(printf '\r')"` 对 225 CR 的文件读 **0**、对 131 CR 的文件也读 **0**；`tr -cd '\r' | wc -c` 才读得出 225／131 ⇒ **盲尺属实**。

## 7 项⑦ 红名册（逐名；取主检出）

`go test -count=1 -v ./internal/tools/` 全层级 `--- FAIL` ＝ **0 枚** ⇒ **逐名＝（无）**，新增红＝**0**。**正控**（证明这把 grep 不是瞎子）：导出树里同一 pattern 打出了 `--- FAIL: TestTicket252R1AlignmentAddsNoAuthorization`。
⚠**编排者那张名册提醒里的数也不准**：`internal/tools` 里 `TestTicket252*` 是 **10 枚**（不是 8）；`TestReadGitOnThisRepositoryIsSelfConsistent` 在 **`internal/panel/git_test.go:596`**（不在 `internal/tools`）；导出树那一发我只看到 **1 枚红**，根因是 8.3 短名前提出缺（`paths_shortname_252_r1_test.go:163-166`），**不是"无 `.git` ⇒ 8 枚假红"**。结论不变：名册必须取主检出。

## 8 项⑧ 写面纪律（四笔逐笔 `git show --name-only`；正控＝同尺打含 `.go` 的真提交）

| 提交 | 文件 | `.go` | 写面外 |
|---|---|---|---|
| `62a1521c` | `internal/tools/task_output_ac2_before_test.go` | **1** | 0 |
| `ec8786e1` | `docs/evidence/s1/164-r1-impl.md` | 0 | 0 |
| `0786b5c6` | `docs/evidence/s1/164-r1-impl.md` | 0 | 0 |
| `1146f803` | `docs/evidence/s1/164-r1-impl.md` | 0 | 0 |

**正控**：`git show --name-only --format='' 62a1521cb | grep -c '\.go$'`＝**1** ⇒ 尺能看见 `.go`，"0"才是真读数。越界面（`internal/risk|panel|session|memory|config/**`、`cmd/**`、`scripts/**`、`.github/**`、`thresholds.go`、golden、`allowlist.txt`）＝**零字节**。**未 push**：`git branch -r --contains` 对四笔全部为空（`origin/dev`＝`6f808a83`，本支领先 436 笔，都是别人的在飞活）。

## 9 项⑨ 两条残余的判语

**① 第三节那颗牙的射程改写——⛔ 没有变成假话，改写正确。** 逐句核：
- "`BuiltinFSEntries` 不夹带 `task.*`" ⇒ **真**：函数体＝`fsRead`＋`fsList`＋`BuiltinFSWriteEntries`（`fs.go:325` 起），用例断言它并绿。
- "本句只管 fs 注册路径，**不是**关于生产注册表的断言；生产注册表今天确实带 `task.output`" ⇒ **真**：生产在同一 `reg` 上先挂 fs（`run.go:528`）再挂 task 家族（`run.go:557`），`task.go:484` `Name() "task.output"`、`task.go:600` `BuiltinTaskEntries` 注册它。
- 旧句"模型今天看到的目录没有 `task.*`，所以连问都学不会"**对生产确为假**（生产看得见 `task.output`）⇒ 改写是**按射程收紧**，不是换了个说法。附带核：它新增的"(SPEC-07 2 states the C3 rule as a refusal, not as an error)"实质为真（`SPEC-07:25` 逐字"未声明即拒绝调用（不是报错，是拒绝）"），只是定位符"2"是松指，不是行号断言。
- ⚠**但新注释里有一处未证成**："the CI logs address this case by them"——`.github/` 与 `scripts/` 对该用例名／文件名 **0 命中**，CI 日志是仓外物，无法在仓内核。建议下一程要么删这半句、要么给日志出处。

**② 文件名＋函数名保留——决定对，但它给的理由不成立，且这笔提交已经弄坏了那支仪器。**
- 决定对：票面 `AC#2b` **按文件名点它**（票 `:27` 逐字含 `internal/tools/task_output_ac2_before_test.go`）；归档仪器 `ac2-anchor-reading.sh:41,44` **按函数名 `-run`** 它。改名要同时改票面与 `.scratch` 仪器＝**两条都在本腿写面外**。
- 理由错：它说"改名会让归档仪器变成 `no tests to run` 的**假绿**"。现量：脚本第 `:16` 是 **`set -eu`**；`go test -run <旧名>` rc=0 但输出只有 `testing: warning: no tests to run`／`PASS`／`ok … [no tests to run]`，**没有** `--- PASS:` 行 ⇒ 仪器那条 grep rc=1 ⇒ `set -e` 让脚本**响亮中止**（机制正控：`set -eu`＋失败管道 rc=1；管道有匹配时 rc=0）。⇒ 改名是**响**，不是假绿。
- **真正已经发生的损伤（它 §8 没报）**：这笔提交把 `t.Logf` 标签从 `AC#2 BEFORE verbatim:` 改成 `unknown-tool refusal:`，而仪器 `:42,:45` grep 的正是 `BEFORE verbatim` ⇒ 今天在 HEAD 跑那支仪器，它**只印出 1 行** `--- PASS: …`，**Text 读数那半已经没了**，而脚本 `:46` 仍无条件打印结论"两行 Text 逐字相同 ⇒ 该用例的读数与被测树无关"——**这才是假绿形状，而且已经落地**。⇒ 残余记账：要么把标签改回、要么改仪器 grep（`.scratch` 不在本腿写面）；在此之前**不许再把那支仪器当 AC#2 的凭据引用**。

## 10 附：它件里两处可改的小错（不阻塞翻框）

- §5"两枚**各** 225 个 CR" ⇒ 实为 **131／225**。
- §8.1"文件头第 **36–39** 行明写这一取舍" ⇒ 实为 **`:46-49`**（`:36-39` 是"teeth measured"那段）。**这两处都是行号／数字腐坏，正是本票在治的病**，登记给下一程动那枚文件时顺手改。

## 11 复跑台件（仓外，只建不删）

`D:\tmp\wisp\20261022\164v1\`：`m1_premise.sh`／`m2_scopes.sh`／`m3_composition.sh`／`m4_assertions.sh`＋`m4b_ast.sh`（含 `astcmp/main.go`）／`m5_mutants.sh`／`m6_gates.sh`／`m7_commits.sh`／`m8_clause.sh`／`m9_misc.sh`／`m10_instrument.sh`；原始日志 `shot0_original.log`／`shot1_M1.log`／`shot2_M2.log`／`shot3_original_again.log`／`gate_*.log`／`fmt_jia_*.txt`／`hash_compare.txt`／`astcmp.log`。
