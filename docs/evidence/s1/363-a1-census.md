# 票 363 普查腿 `363-a1` — AC#0 只读现读＋归因：`MODE-REDLINE … still asks` 在每一发 L0 上都是假话

- **腿**：`363-a1`（**只读**）。锚 `46e61076`（取读数那一发的 HEAD）→ 结件时 HEAD `e9caddff`；
  `git log --oneline 46e61076..HEAD -- internal/risk internal/tools` ⇒ **0 笔** ⇒ 四条链的行号与字节在**两端同一**。
- **写面**：本件一枚（`docs/evidence/s1/363-a1-census.md`）。⛔ 未落任何 `.go` 字节、⛔ 改票面／台账／`Status:`／`.github/**`、⛔ push、⛔ 翻框。
- **一句话读数**：**L0 → `mode.Screen` 落进"无法识别的风险级"支 → `Kept` 非空 → 日志打 "still asks" → 同一发里 `route` 走 `case risk.L0` 直接放行 → 确认卡 0 张 → 工具已执行完。**
  ⇒ ★**"会问"与"0 张卡"出现在同一发、同一进程、同一行决策链上**，⛔ 两发对比。

## §1 四条链逐条复核（HEAD `e9caddff` 现跑；尺＝`git show HEAD:<p> | nl -ba`）

| 票面 | 现跑实测 | 漂 | 逐字（HEAD blob） |
|---|---|---|---|
| ①`internal/tools/bridge.go:311` | **311** | 无 | `sil := mode.Screen(verdict)` |
| ②`internal/risk/mode.go:163-165` | **163-165** | 无 | `if d.Level != L1 && d.Level != L2 {` ／ `return Silenced{Level: d.Level, Kept: "无法识别的风险级，fail-closed 保留询问", Mode: m}` ／ `}` |
| ③`internal/tools/bridge.go:324` | **323**（条件）／**324-325**（日志call） | 无 | `case sil.Kept != "":` ／ `b.log("tools: MODE-REDLINE mode=%s still asks tool=%s rules=%v because=%s",` ／ `mode, req.Name, verdict.RulesHit, sil.Kept)` |
| ④`internal/tools/bridge.go:352` | **352**（`b.route`） | ★**注释在 351，⛔ 352** | `// --- routing: L0 pass, L1 window, L2 approval, Deny never reaches a gate -`（351）／ `ok2, why := b.route(ctx, &dec, sil, grantID)`（352） |

- ④的落点是 `Bridge.route` 的 switch：`bridge.go:432` `switch sil.Level {` ／ **`433-435`** `case risk.L0:` `dec.DecisionColumn = agent.DecisionAllow` `return true, ""` ⇒ **L0 不碰任何 gate**（⛔ 是"闸门放行"，是**没有闸门**）。
- **行号报告**：①②③ 与票面**逐字同号**；④ 票面把那条注释并进了 `352`，实测注释在 **351**、`b.route` 在 **352** ⇒ **唯一处 1 行偏差**，不影响定性。
- **归因闭合**：链①把 L0 判定交给 `Screen`；链②在**任何模式判断之前**把 `L0`（≠L1、≠L2）判成"无法识别"；链③只看 `Kept != ""` 就宣称 "still asks"、**⛔ 看这一发到底有没有卡**；链④同一发里 0 卡放行。⇒ 四环都在，且**环③的判据与环④的事实相反**。

## §2 ★造一发真读数：那行 `still asks` ＋ 同一发确认卡 ＝ 0（本格核心）

工具＝**真的** `fs.read`（D34 真声明 `FSReadDecl()`、`Declared: risk.L0`，`fs.go:296-305`），路径＝临时工作根里**真实存在**的文件（在授权目录内 ⇒ 判定走真 R1–R9、无规则命中）。
读数为**逐字**（`LOGVERBATIM|` 后即 `Options.Logf` 收到的原文，未加一字；三个档各一发）：

```
----- L0 fs.read / ask_every_step -----
ASSESSED level=L0 rules=[] reason="无规则命中（L0 直接执行）"
EMITTED level=L0 mode=ask_every_step silenced=false kept="无法识别的风险级，fail-closed 保留询问" decisionColumn=""
CARDS L1_WINDOWS=0 L2_APPROVALS=0 TOTAL=0
LOGVERBATIM|tools: MODE-REDLINE mode=ask_every_step still asks tool=fs.read rules=[] because=无法识别的风险级，fail-closed 保留询问
LOGVERBATIM|tools: call kind=success task=task-1 corr=corr-1 tool=fs.read risk=L0 decision=allow outcome=success rules_hit=[] in_allowlist_scope=true grant_id=0 refusal=none reason="无规则命中（L0 直接执行）"
EXEC err=<nil> isError=false text="hello"          ← 工具真跑了，返回了文件内容

----- L0 fs.read / ask_high_risk -----  同上，只有 mode= 变：LOGVERBATIM|tools: MODE-REDLINE mode=ask_high_risk still asks tool=fs.read rules=[] because=无法识别的风险级，fail-closed 保留询问   CARDS TOTAL=0   EXEC text="hello"
----- L0 fs.read / auto_approve -----   同上，只有 mode= 变：LOGVERBATIM|tools: MODE-REDLINE mode=auto_approve still asks tool=fs.read rules=[] because=无法识别的风险级，fail-closed 保留询问   CARDS TOTAL=0   EXEC text="hello"
```

- ★**同一发里的三件事**：③那行日志说"still asks"；②`CARDS L1_WINDOWS=0 L2_APPROVALS=0 TOTAL=0`（一张卡都没有）；①`kind=success … decision=allow outcome=success` ＋ `text="hello"`（**已经执行完、连结果都回来了**）。
- ★**与档无关**：`ask_every_step`／`ask_high_risk`／`auto_approve` **三个档都假**——因为链②的 `L0` 支排在 `m == ModeAskEveryStep` 判断**之前**（`mode.go:163` 在 `168` 之前）⇒ 不存在"某个档才假"。
- ★假话的**形状**：`Kept` 的文档语义是「**拒绝静默**的那条红线」（`mode.go:143-146` 逐字：`Kept names the red line that refused to silence the verdict ("" when the mode applied as asked)`）。L0 上**没有红线**、也**没有问题被保留** ⇒ 这句 `because=` 是**范畴错误**，它把"我认不出这个级"说成了"我保留了这个询问"。

## §3 ★正控（硬闸）：L1／L2 上这条链**说真话**，⛔ 是"所有档都假"

| 腿 | 判定 | 卡数 | `MODE-REDLINE … still asks` | 判 |
|---|---|---|---|---|
| L1 `fs.write`(新文件) / `ask_every_step` | `L1` `[R1]` | **L1窗口=1** | ★**根本没有这行**（LOG_LINES=2，只有 `kind=refused`＋`PATH-ACCOUNT`） | 该问就真问了一张卡，**⛔ 多话** |
| L1 `fs.write`(新文件) / `ask_high_risk` | `L1` `[R1]` | **0** | ⛔ 这行；改打 `MODE-SILENCE … (a silenced question, not an allow: no user click happened)` | 静默有**另一条**日志，语义正确 |
| L2 `fs.write`(覆盖已有) / `ask_high_risk` | `L2` `[R1 R8]` | **L2审批=1** | `… still asks tool=fs.write rules=[R1 R8] because=ask_high_risk 只静默 L1，本判定是 L2` | ★**说的正是真的**（1 张卡在） |
| L2 `fs.write`(覆盖已有) / `auto_approve` | `L2` `[R1 R8]` | **L2审批=1** | `… still asks tool=fs.write rules=[R1 R8] because=R1 工具自身声明的 L2（不可逆自认，PLAN.md:1629）` | ★**真的**（1 张卡在） |

⇒ ★**正控成立**：`Screen` 认得的两档，日志与卡数**一致**；只有 `L0` 这一档"日志说会问、实际 0 卡"。⇒ 票面定性「**只在 L0 上假**」**成立**，⛔ 需改。

### §3b ⚠ 超出票面四条链的一条实测（**反向**的第三种假话，供 AC#1 裁形用）
- `fs.read` 打 `%APPDATA%\wisp\config.toml` ⇒ `ASSESSED level=Deny rules=[R2 R3]`、`kept="拒绝级 verdict（Deny 级不受任何模式影响）"`、**`CARDS TOTAL=0`**、`text="该目标被绝对禁止访问（R3 A 档…）"` ⇒ ★**日志仍打 `still asks`，而事实是"拒绝、⛔ 问"**。`~/.ssh/id_rsa` 同形（`auto_approve` 下也一样）。
- ⇒ `still asks` 这句在**两个**方向上都假：L0（放行了却说会问）、Deny（拒绝了却说会问）。★这决定 §9 的判读，但**裁形仍归 AC#1**。

## §4 「`L0` 是不是一个认得的级」— 我的判读：**是，`Screen` 这一支的判据错了**

- **定义**（`internal/risk/assessor.go`）：`41` 逐字 `// L0 < L1 < L2 < Deny: fusion takes the max, and Deny is terminal.`；`46` `L0 Level = iota` ⇒ ★**`L0` 是这个枚举的零值**；`45` `// L0: read-only, no side effects — execute directly.`
- **它是判定器的"正常结论"而不是异常**：`315` `level := L0`（融合的**起点**）；`347` `reason = "无规则命中（L0 直接执行）"`；`56-59` `Level.String()` **把 `L0` 渲染成 `"L0"`**（有名字、进得了日志）。本格实测那行 `risk=L0` 就是它渲染出来的。
- **`Screen` 自己也把 `L0` 当输出级用**：`mode.go:175`、`188` 逐字 `return Silenced{Level: L0, Silenced: true, Mode: m}` ⇒ ★**同一个函数，出口认 `L0`、入口不认**。⇒ `mode.go:155-156` 的自述 `Deny and any unknown level are handled first` 里那个 "unknown" **被实现成了 "not L1/L2"**：`L0` 已知，却被归进 unknown。
- **`Silenced` 的契约**（`mode.go:138-149`）：`Kept` 只在**有红线拒绝静默**时非空（`"" when the mode applied as asked`）。`L0` 上无红线、无问题 ⇒ **该支应返回 `Kept:""`**（且 `Silenced:false`：没有"问题"被移走，`Silenced` 的自述是 `the caller must NOT ask`，语义上 L0 属于"压根没有询问"，与 `ask_every_step` 对未触碰判定的形状同形）。
- **⇒ 判读**：`Screen` 把 `L0` 判成"无法识别的风险级"**这件事本身是错的**（⛔ 只是日志难看）。假话的载体有两个：**日志行**（`bridge.go:324`）与 **`Decision.ModeKept`**（`bridge.go:314` 写入、`gate.go:48-50` 字段自述"names the red line…"）。
- **行为面今天无损失、契约面已破**：L0 的 `route` 仍放行（`433-435`），所以**没有功能回归**；但 `ModeKept` 的文档语义在**每一发 L0** 上被违反。★实测：全仓 `ModeKept` 的读取方**只有测试**（`grep -rn ModeKept --include=*.go` ⇒ 12 处，`cmd/**` 与 `internal/panel/**` **零处**），所以**今天唯一的操作员可见面就是那行日志**。

## §5 我用的尺（对象／形状／谁算的／射程）

- **仪器**：`go test -overlay`。探针文件在**仓外** `D:\tmp\wisp\20261022\363a1\zz_t363a1_probe_test.go`，用 overlay 映射成 `internal/tools/zz_t363a1_probe_test.go` 的**虚拟路径**参与编译 ⇒ ★**仓内 `.go` 字节 0 枚**，`git status --porcelain -- internal cmd` 全程为空（首／尾各测一次）。
- **组合＝生产组合**：`tools.New(Options{Registry: BuiltinFSEntries(d), Paths: NewPathCanonicalizer(...), Gate, Modes, Provenance, OnDecision, Logf})`，与 `cmd/wisp/run.go:764-788` **同形**；唯一差异＝`Logf` 指**捕获器**（生产指 `rt.auditf`）。★**真** registry／真 D34 声明／真 R1–R9 判定器（⛔ 注入 `Options.Assessor`）。
- **「确认卡数」怎么数（⛔ 数 UI、⛔ 数 tick）**：`Gate` 接口就是两张卡的面（`gate.go:110-117`：`PendingWindow`＝L1 预执行窗口、`PendingApproval`＝C18 审批队列）；探针数**这两个回调被调用几次**，`TOTAL=` 是两者之和。`TOTAL=0` ⇒ **bridge 一张卡都没开**。
- **日志到操作员的那一段＝读码，⛔ 实测**：生产 `Logf: rt.auditf`（`run.go:787`）→ `run.go:979-983` `fmt.Fprint(rt.stderr, "[audit] "+line)` ＋ `sink.logger().Info("audit: "+line)` → 滚动 JSONL `<data>\logs\wisp-<YYYYMMDD>-001.jsonl`（`logsink.go:129` 起）。★本腿**没有**跑整条 `wisp run`（那要 golden SSE 台件），这一段**按码认**、已在此标明；行**格式与调用点**是生产原文（不变量在 `bridge.go:324`，⛔ 在我的探针里）。
- **谁算的／几位**：行号＝`nl -ba` on `git show HEAD:<path>`；文件枚数＝`git show --name-only --format= <sha> | tr -d "\r" | tr "\\" "/" | grep -Ec "\.go$"`；卡数＝探针内 `gateSpy.counts()`（互斥锁保护的两个 int）。
- **环境**：`GOFLAGS=`（置空）、`GOTMPDIR/TMPDIR/TEMP/TMP` 未设、⛔ 未用 `CGO_ENABLED=0`；`go1.27`、Windows。**临时件只建不删**（`D:\tmp\wisp\20261022\363a1\`：探针、`overlay.json`、`run-1.txt`、`d22scan.out.txt`、`pathbudget.out.txt`）。

## §6 票面前提新鲜度（逐条对 HEAD blob 复量）

- ①②③④ 四条引用**全部在票面写的号上**（见 §1 表）；**唯一偏差**＝④ 的那条注释在 **351**、票面并入 `352`。
- 票面对 `mode.go:163-165` 的**逐字**引用与 HEAD blob **逐字相同**（含中文串 `无法识别的风险级，fail-closed 保留询问`）。
- `bridge.go` / `mode.go` 自 `46e61076` 到 `e9caddff` **零提交**（`git log 46e61076..HEAD -- internal/risk internal/tools` 空）⇒ 读数窗口内这两包**没被别人动过**，结论与该窗口的 HEAD 绑定。

## §7 门禁（只读腿自证，各报 rc）

- `sh scripts/d22scan.sh` ⇒ **rc=0**（`clean - no D22 ban violations`；`ban #8 internal/=621`、`cmd/=148`、`design/=39`、`frontend/=116`；★`ban #8` **只扫** `internal/`＋`cmd/` 的 `.go`，本件是 `docs/` ⇒ 不在射程）。
- `sh scripts/check-path-length-budget.sh` ⇒ **rc=0**（`denominator: tracked paths=17601 over-budget=61 covered by roster=61 not in roster=0`；`VERDICT GREEN`）。★本件路径 28 字符，远在帽下。

## §8 提交／`.go` 尺／正控尺

- 本件的**读数窗口** HEAD＝`46e61076`，**结件时** HEAD＝`e9caddff`（别的腿在连续提交，两颗都不含 `internal/risk`／`internal/tools` 的改动）。
- 本件**自己的提交号不写进本件**（防自我推翻，票 160 `02b3882c` 先例）⇒ 提交号、逐枚文件、`.go`=0 落在**下一笔**（文末 `§8b`），本节只登记**尺**与**正控**：
  - **尺（逐字）**：`git show --name-only --format= <sha> | tr -d "\r" | tr "\\" "/" | grep -Ec "\.go$"`（★`[\\/]` 归一，⛔ 裸 `grep -c '\.go$'` 在反斜杠输出上会漏）。
  - **正控（本腿现跑，证明这把尺不是瞎尺）**：`ec520750`（票 179 修法那一笔，**非本程**）⇒ **2 枚**（`internal/agent/loop.go`、`internal/agent/declared_l0_risk_179_test.go`）⇒ ★尺在真 `.go` 提交上**报得出来**。
  - **本件的预期读数**：1 枚、`docs/evidence/s1/363-a1-census.md`、**`.go`=0**、`internal/`+`cmd/`=**0**。

### §8b 提交后复量（★读数在提交之后取）

- 本件第 1 笔＝**`7866ad4e`**：`git show --name-only --format= 7866ad4e | grep -c .` ⇒ **1 枚** ＝ `docs/evidence/s1/363-a1-census.md`；
  ★`.go`＝**0**、`internal/`＋`cmd/`＝**0**、全在写面内。★**正控尺**（§8）：同一把尺打 `ec520750` ⇒ **2 枚 `.go`** ⇒ **0 是真 0、⛔ 瞎尺**。
- **提交后**（HEAD `7866ad4e`）两道门重跑：`sh scripts/d22scan.sh` ⇒ **rc=0**；`sh scripts/check-path-length-budget.sh` ⇒ **rc=0**。
- ★本件**自己的**（放本节这一笔的）提交号不写进本件——同 §8 的先例，防自我推翻。

## §9 对 AC#1 的输入（★我判读、⛔ 我裁形）

- ★**乙（`mode.Screen` 认 `L0`）不是"可选的美化"，是正确性**：`L0` 是认得的级（§4）⇒ 现支的判据把"已知"当"未知"。甲单独做，`Decision.ModeKept` 会继续在**每一发 L0** 上写着"无法识别的风险级"，`Screen` 的自述契约继续是假。
- ★**甲也有独立缺陷**（§3b 实测）：`still asks` 在 `Deny` 上同样假（0 张卡、直接拒绝）。乙**修不到**这一条。
- ⇒ ★**指向丙（两者都做）**：乙修 L0 的误判与 `ModeKept`；甲把"有卡才说 still asks／无卡说清是放行还是拒绝"这句文案改成能说实话的形状。**方案与代价的最终裁定归编排者（AC#1）**，本腿只交读数与判读。

## §10 复跑（逐字）

```bash
cd "/d/work/workspace/projects plans/Wisp"
GOFLAGS= go test -overlay 'D:/tmp/wisp/20261022/363a1/overlay.json' -count=1 -run TestT363a1CensusReadings -v ./internal/tools
sh scripts/d22scan.sh; echo rc=$?
sh scripts/check-path-length-budget.sh; echo rc=$?
```
