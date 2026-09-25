# 票 143 对抗验收 r1 — `wisp panel-assets -taint-source` 那条腿：R4 到底是不是 assessor 判出来的

裁决者：**非实现者**（本件作者，对抗验收程）。被验物＝票 143 的交付（`cmd/wisp/panel_assets.go` + `cmd/wisp/panel_assets_143_test.go`）。
实现方自证表＝`docs/evidence/s1/143-panel-assets-r4-leg-r1.md`（**本件不复用它任何一行读数当自己的凭据**，逐格现量）。

## §0 进场自量与工作树前置

- 进场 `git rev-parse HEAD` = **`d24216809eed57482c1bbf71f4174ce1a3020916`**（短 `d242168`；与简报 14:0x 那枚同值）。
- `git status --porcelain -- cmd/wisp/` 进场读数 **空**（0 行）⇒ 被验目录干净，未读脏工作树。
- 本程全程共享树并发漂移：`d242168` →（读数时）`5b4352f` → `0e17914` → `3855dc6`。
  每次漂移后都复算过 `git diff --name-only d242168 HEAD -- cmd/wisp/` = **0 行**（现量两次：`0e17914`、`3855dc6`）
  ⇒ 被验版本自始至终是那一份，别程没往 `cmd/wisp/` 写过一字节。
- 交付两枚 commit 的路径（`git show --pretty=format: --name-only` 现量，非引用自证表）：
  - `9be3288` → `cmd/wisp/panel_assets.go`、`cmd/wisp/panel_assets_143_test.go`
  - `4037539` → `cmd/wisp/panel_assets.go`
  **别家路径 0 枚**。
- 工作树里与本程无关的现场（本程既不动也不计入任何零命中证据）：
  `M docs/reports/frontend-watchdog.md`（别程）、`design/**` 的 20 行未提交状态（owner 的东西，单独量到就单独写，见 §8）。
- 本程变异全部落仓外副本 `D:\tmp\acc-143\`（`before/`＝`git archive 9be3288^`、`after/`＝`git archive d242168`、
  `mutA/`＝`git archive d242168` 的变异场），**不复用实现方的 `D:\tmp\ticket-143\mut`**。
  每发之后复算 `git status --porcelain -- cmd/wisp/`：末次读数 **空（0 行）**，且 `grep -rc MUTANT cmd/wisp/` 无任何非零命中。
  临时件只建不删（`acc-*.log`、`acc2-*.log`、`roster-*.txt`、`pristine-panel_assets-acc143.go`、`hm1.diff` 全留）。

内容指纹（本程自算，用来钉"被验的就是那枚"）：

```
git show 9be3288^:cmd/wisp/panel_assets.go | md5sum  = 1372bf92d4ca6792d2608b6fe514f1c9
git show 4037539:cmd/wisp/panel_assets.go  | md5sum  = 7f200f275ed8ea854e84a104ea0fc58f
git show d242168:cmd/wisp/panel_assets.go  | md5sum  = 7f200f275ed8ea854e84a104ea0fc58f   （锚点＝交付终版）
md5sum after/cmd/wisp/panel_assets.go      = 7f200f275ed8ea854e84a104ea0fc58f             （归档副本同值）
```

⇒ 简报"锚点必须自量"这一条：**成立**，且简报给的 `d242168` 就是被验版本。

---

## §1 攻点 1（最要害）：ⓐ 那一支有没有降级成"cmd 内自己写一枚 detector"

### 1.1 数据通路逐跳（每跳 `file:line`，全部本程直读源码，锚点 `d242168`／`3855dc6` 两处复核同值）

| # | 跳 | 出处 | 这一跳干了什么（读数） |
|---|---|---|---|
| 1 | CLI 交出的只有来源事实 | `cmd/wisp/panel_assets.go:199-214`（`Set`） | 只切 `<tool>|<origin>|<content>` 三段并校验非空，**没有任何等级/规则号/文案** |
| 2 | 事实交给 C25 引擎 | `cmd/wisp/panel_assets.go:231` `risk.NewProvenance(risk.ProvOptions{NoProbe: true})` | 真引擎构造口（`internal/risk/provenance.go:266`） |
| 3 | 开 scope | `cmd/wisp/panel_assets.go:232` → `internal/risk/provenance.go:339` | `OpenScope("panel-assets-l2")`（常量在 `panel_assets.go:166`） |
| 4 | 逐条登记污点 | `cmd/wisp/panel_assets.go:234` → `internal/risk/provenance.go:362` | `Mark(scope, tool, origin, content)`：归一化＋建 ≥8 字符片段索引（`:380-382`），`tool` 不在 SPEC-06 §5 名单里也**照记**（`:377-379` fail-closed） |
| 5 | 取绑定到 scope 的检测器 | `cmd/wisp/panel_assets.go:236` → `internal/risk/provenance.go:605-607` | `Detector()` 返回 `boundDetector{p, scopeID}`——**cmd 侧没有、也无法构造这个类型**（未导出） |
| 6 | 接进 assessor | `cmd/wisp/panel_assets.go:66` → `internal/risk/assessor.go:207-210` | `WithTaintDetector(t)` 只是 `a.taint = t`（字段 `:181`） |
| 7 | 卡片入口 | `cmd/wisp/panel_assets.go:68` → `internal/panel/approval.go:64-66` | `NewApprovalCardView` 内部 `assessor.Assess(subject.Tool, paramsFromArgs(subject.Args), subject.Facts)`；`paramsFromArgs` 只产 `command`/`argv` 两枚 key（`approval.go:107-114`） |
| 8 | 事实进 judge 上下文 | `internal/risk/assessor.go:228-236` | `ctx.taint = a.taint`（`:235`） |
| 9 | R4 那条规则被跑到 | `internal/risk/assessor.go:267`（`builtinRules` 里 `ruleTaint,`）→ `internal/risk/rules_gateway.go:101-115` | `:102` `if ctx.taint == nil { return nil }`（休眠门，票面引的那句逐字在，跨 `:99-100` 两行）；`:105` `source, hit := ctx.taint.TaintHit(ctx.params)`；**`:110` `rules: []RuleID{R4}`、`:111` `level: L2`、`:112` `reason: fmt.Sprintf("R4: 包含来自 %s 的内容", source)`、`:113` `blockSessionAuth: true`** |
| 10 | 检测器里的真匹配 | `internal/risk/provenance.go:614-620` → `:472` `Inspect` → `:587-599` `matchText` | `matchText` 按 mark 顺序遍历片段索引，命中才 `Hit{SrcTool: m.tool, Origin: m.origin, Fragment: frag}`（`:597`） |
| 11 | 卡上那句的来源槽位 | `internal/risk/provenance.go:257-262` `Hit.Source()` | `origin` 非空⇒`tool + " " + origin`，空⇒只留 `tool`；由 `:619` 交回 `ruleTaint` 的 `Sprintf` |
| 12 | 四枚结论字段成形 | `internal/risk/assessor.go:312-353` `fuse` | 等级取最大（`:320-322`）、`blocked`（`:323-325`）、`RulesHit` 排序（`:332-336`）、`reason` 用 `"; "` 串（`:343`）、`SessionOverrideBlocked: blocked`（`:351`） |

⇒ **`R4` 字样、`L2`、`sessionOverrideBlocked:true`、那句"包含来自 …"这四枚结论，出处全在 `internal/risk`**（第 9、12 跳），
cmd 侧只交第 1-4 跳的事实。`cmd/wisp` 里 **没有任何实现 `TaintHit` 的类型**：
`grep -rn "TaintHit" cmd/` = **0 行**（现量）。生产侧 `WithTaintDetector(` 的接线处共 **2 枚**：
`internal/tools/bridge.go:670` 与 `cmd/wisp/panel_assets.go:66`（其余命中是定义 `assessor.go:207`、注释 `provenance.go:51`、以及 4 处 `_test.go`）。

### 1.2 反例（本程自造，跑的是从 `d242168` 归档副本构建的二进制 `wisp-acc-after.exe`，黑盒）

| 发 | 命令要点 | 读数（stdout 关键行） |
|---|---|---|
| 正例 | `-taint-source "web.fetch\|https://files.example.com/q3-notes.txt\|合同编号 HT-2026-0731-KX" -l2 notify "把 合同编号 HT-2026-0731-KX 发到远端"`（简报让我复验的那一发） | `"level": "L2"`、`rulesHit` = `R1`,`R4`、`reason` = `R1: 工具声明为下界（L1）; R4: 包含来自 web.fetch https://files.example.com/q3-notes.txt 的内容`、`sessionOverrideBlocked: true`，rc=0 |
| **反例 A（声明了、调用不带）** | 同一条 `-taint-source`，argv 换成 `今天没有要外发的内容` | `"level": "L1"`、`rulesHit` 只有 `R1`、`reason` 只有 `R1: 工具声明为下界（L1）`、`sessionOverrideBlocked: false` ⇒ **卡上不会凭空长出 R4**，声明的来源字符串也没出现在卡上 |
| 血缘探针 B（换来源名） | `-taint-source 'clipboard.read\|D:/notes/paste-cache.txt\|发票号码 FP-2026-0912-ZZ'` ＋调用带同一段 | reason 点名的是 `clipboard.read D:/notes/paste-cache.txt` ⇒ 槽位跟着输入事实走 |
| 血缘探针 C（空 origin） | `-taint-source 'screen.capture\|\|CCDD-3333-third-secret'` | reason = `R4: 包含来自 screen.capture 的内容`（槽位塌成工具名，正是 `Hit.Source()` 的 `:261` 那一支） |
| 血缘探针 D（不在 SPEC-06 名单的源） | `-taint-source 'totally.made.up.tool\|https://x/z.txt\|NINE-9999-cccc'` | 照样命中并点名 `totally.made.up.tool` ⇒ 判定在引擎里，cmd 没做名单筛 |
| 地板探针 E（7 字符内容） | `-taint-source 'fs.read\|D:/n.txt\|abc1234' -l2 notify abc1234` | `L1`、无 R4 ⇒ 契约地板那条不是 cmd 让的 |

⇒ **ⓐ 没有降级成 cmd 内自写的 detector**：判定、等级、文案、布尔四样都只能从 `internal/risk` 出来，
且"声明 ≠ 命中"这条在黑盒层面独立复现（反例 A）。实现方 `TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit`
**确实在拦**（§4 用变异 M2 复算：换成永真 detector 时它的两个子用例全红）。

### 1.3 但要给这一格钉一枚精度纠正（简报与自证表都说轻了）

`TestAC1AC2TaintSourceLegProducesAJudgedR4` **单独**不足以证明接的是真引擎。本程自造变异 **M2b**
（仓外副本 `mutA`：`detector()` 返回一枚"永远命中、但把声明的 `tool + " " + origin"` 原样吐回去"的 cmd 内 detector）⇒

```
--- PASS: TestAC1AC2TaintSourceLegProducesAJudgedR4          ← 它照样绿
--- FAIL: TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit （两个子用例都红）
--- PASS: 其余五枚全绿（六枚总数里除 AC3 之外那五枚）
```

⇒ 真正卡住"cmd 自写 detector"的是 **AC3 那一枚**，不是 AC1AC2。
自证表 §1 末句"塞一枚返回常量的 detector ⇒ §3 那两枚用例**会一起红**"对**它那枚**变异（返回一个与声明无关的固定串）成立（§4 的 M2 复算就是两枚红），
但作为一般陈述**说过头了**：返回"声明的来源串＋永远 true"的 detector 只红一枚。判**不构成退回**（用例集合整体有牙，红的那枚本就是票面 AC#3(ii) 要的防过度匹配），
记为**残余精度**：若将来要更强，只需在 AC3 之外补一枚"detector 的类型必须是 `prov.Detector` 的返回值"的白盒断言（本程不修）。

---

## §2 攻点 2："注入事实、不注入结论"——grep 逐枚

尺子：`grep -n -E 'R4|R9|包含来自|rulesHit|sessionOverrideBlocked|L2:' cmd/wisp/panel_assets.go`（锚点 `3855dc6`，`cmd/wisp/` 0 行脏）。

**命中数 = 0 行**——而且这是**连注释一起算**的 0（票面与自证表口径只要求"输出路径 0"）。逐枚点名这三类"结论字样"：

| 结论字样 | 在 `panel_assets.go` 的命中 | 判定 |
|---|---|---|
| `R4` | 0 | 卡上那句 `R4: 包含来自 …` 不可能由这枚文件拼出（拼装点唯一：`internal/risk/rules_gateway.go:112`） |
| `rulesHit` | 0 | JSON 键名来自 `internal/panel/approval.go:46` 那枚结构体 tag，不在这条 CLI 路上 |
| `L2`／`sessionOverrideBlocked`（作值） | 0 | 等级与布尔都由 `fuse`（`internal/risk/assessor.go:320-325,351`）产出 |

**点名区分（注释／既存的、不属于本票的字样，与"输出路径"不是一回事）**：

- `L2` 单独出现 **3 枚**：`:4`（文件头注释）、`:41`（`-l2` 的 help 串 "print the L2 card JSON"）、`:146`（`splitFacts` 的注释）。
  三枚**逐枚对改前文件核过是既存件**：`git show 9be3288^:cmd/wisp/panel_assets.go` 里同一批句子的行号是 `:4`/`:40`/`:133`
  （本票只是把它们往下推了 1 行或原地未动）⇒ **本票没有新增任何 `L2` 字样**，`:41` 那枚 help 串也不是输出路径上的值——
  它是 `-l2` 这枚旗标的名字与说明，卡上等级实际打印的是 `decision.Level.String()`（`internal/panel/approval.go:81`）。
- `R8` 出现 **2 枚**：`:42`（既存 `-irreversible` 的 help 串 "(R8's input)"）、`:145`（`splitFacts` 注释）。
  同样核过是**票 114 时代的既存件**（改前文件 `:41`/`:132` 同句），不是本票造的。

**这条 CLI 路的 stdout 写出点**（`panel_assets.go:79-81`）只有 `enc.Encode(view)` 一枚；
`view` 的九个字段来源逐枚点名：`CorrelationID`/`Tool`/`Args`/`CallChain` ← 调用方给的（`:69-72`，其中 `Args` 就是 argv），
`Level`/`RulesHit`/`Reason`/`SessionOverrideBlocked` ← `risk.Decision`，`ReasonKnown`/`DecidedBy` ← `internal/panel`
（`approval.go:77-88`）。**没有一枚字段由 cmd 侧填判据字符串。**

**仪器本身**：`TestAC1CmdSideEmitsNoVerdictTokens`（`panel_assets_143_test.go:264-293`）用 `go/ast` 扫 `panel_assets.go` 的
**每一枚字符串字面量**（注释豁免、字面量不豁免），禁词表 `:265` = `R4 / R9 / 包含来自 / rulesHit / sessionOverrideBlocked / L2:`。
本程独立复跑该枚：**绿**；并复算 M1（摘接线）与 M2（永真 detector）两发变异下它都**仍然绿**（它只钉字面量形状，钉不出接线缺失，这是设计内的分工）。

⚠ **一枚要报回给编排口的口径差**：那张禁词表只封了 `R1–R9` 里的 **两枚**（`R4`、`R9`），
而 `R8` 恰好以既存 help 串的形状活在同一个文件里（`:42`）。也就是说：**仪器不是"禁所有规则号"，是"禁本票这条路的规则号"**——
今天够用（`R8` 那枚是输入事实的说明、不是判据输出），但它给下一个人的保护比表面上看起来窄一格：
将来若有人在 cmd 侧把 `R2:`／`R3:`／`R8:` 拼进输出，这枚仪器不会红。登记为**残余（不退回）**，归口建议见 §8 末。

## §3 攻点 3：默认路径不许渗——本程重造的 `cmp`（不收它的"三发相同"）

台件（都是 `git archive` 出去的仓外副本，同一台机器同一环境同一把尺）：

```
before/ = git archive 9be3288^   → go build -o wisp-acc-before.exe ./cmd/wisp/   （md5 源文件 1372bf92…）
after/  = git archive d242168    → go build -o wisp-acc-after.exe  ./cmd/wisp/   （md5 源文件 7f200f27…）
两份副本各自 cp 了仓里的 frontend/dist（embed 内容同值），构建 rc 都是 0
运行 PATH 里带 third_party/sherpa-onnx（票 98 那一族：不带则进程连加载都过不去）
```

| 发 | 命令（`panel-assets` 之后） | 改前 stdout 字节 | 改后 stdout 字节 | `cmp` | rc |
|---|---|---|---|---|---|
| A | `-l2 shell.run rm -rf D:/tmp/x` | 365 | 365 | **identical** | 0 / 0 |
| B | `-irreversible delete,overwrite -l2 fs.write D:/notes/a.txt` | 418 | 418 | **identical** | 0 / 0 |
| C | `-l2 notify hello world` | 348 | 348 | **identical** | 0 / 0 |
| **D（本程自加）** | `-l2 fs.write D:/secret/key.txt payload` | 366 | 366 | **identical** | 0 / 0 |
| **P（本程自加，不进 `-l2` 分支）** | `panel-assets`（裸的资产摘要那一支） | 60 | 60 | **identical** | 0 / 0 |

原文读数（`/d/tmp/acc-143/`）：

```
$ for x in A B C D P; do if cmp -s $x-before.json $x-after.json; then echo "$x stdout: identical"; fi; done
A stdout: identical (365 bytes)
B stdout: identical (418 bytes)
C stdout: identical (348 bytes)
D stdout: identical (366 bytes)
P stdout: identical (60 bytes)
```

B 卡现量（两枚二进制同一串字节）：`"level": "L2"`、`rulesHit` = `R1`,`R8`、
`reason` = `R1: 工具声明为下界（L1）; R8: 不可逆操作（永久删除、覆盖已有内容）`、`sessionOverrideBlocked: false`；
C 卡现量：`L1 / [R1]`、`sessionOverrideBlocked: false`。

⇒ **"不带新旗标时 stdout 与改前逐字节相同"这一格：成立**，且比实现方那三发多两发（D 是路径形状、P 是完全不进 `-l2` 分支的那一支），
rc 也逐发对齐。它自报的 `4712ea6`／`5ef1632` 两枚二进制我不引用——本程的两枚是从 `9be3288^` 与 `d242168` 现构的。

---

## §4 攻点 4：承重两发独立重跑（仓外副本 `D:\tmp\acc-143\mutA`，不复用它的 `D:\tmp\ticket-143\mut`）

台件：`git archive d242168` 解出的副本 ＋ 补 `frontend/dist` ＋ 补 `third_party/`（**这两步是必需的**：
不带 `third_party` 时副本里 **8 枚**与本票无关的用例会红，红因逐字是
`secret_argv_windows_test.go:271: no native DLLs in ..\..\third_party\sherpa-onnx - run scripts/fetch-deps.ps1 first`；
补上之后同一副本控制跑＝全绿。⇒ **归档副本天生不是等效环境，必须补这两味**，本程与实现方各自的数都是在补过之后取的。）

| 发 | 变异（只落副本） | 顶层六枚（`-run` 六枚那把尺） | 整包四数（`-count=1 -v ./cmd/wisp/`） |
|---|---|---|---|
| 控制 | 无（md5 `7f200f27…` 与锚点同值） | 12 枚 RUN（6 顶层＋6 子用例）**全绿** | **RUN=114 PASS=60 FAIL=0 SKIP=0**（无 panic） |
| **M1 摘接线** | `:66` `assessor = assessor.WithTaintDetector(det)` → `_ = det // MUTANT-ACC-M1` | **唯一红**：`TestAC1AC2TaintSourceLegProducesAJudgedR4`；其余 **5 枚全绿** | **RUN=114 PASS=59 FAIL=1 SKIP=0**（红的就是那一枚） |
| **M2 换永真 detector** | `detector()` 末尾改成返回 `mutantAlwaysHit{}`（恒 `(src="MUTANT-ACC-M2 always-hit", true)`，`prov` 整块不接） | **红两枚**：`TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit` 的**两个子用例都红**（`:189` 与 `:212` 两条 RED 原文都印出 `rulesHit = [R1 R4]`）**以及** `TestAC1AC2…` 一起红；默认路径那枚仍绿 | **RUN=114 PASS=58 FAIL=2 SKIP=0** |
| **M2b（本程自加，更贴"cmd 自写 detector"那一形）** | `detector()` 返回 `mutantEcho{f.sources}`：**永远命中，但把声明的 `tool + " " + origin` 原样吐回去** | **只红 AC3**（两子用例），`TestAC1AC2…` **仍绿** ⇒ 见 §1.3 | 未跑整包（子集已定性） |
| **M3（本程自加，钉攻点 5 那一格）** | 摘掉 `fs.Usage` 里的 `fs.PrintDefaults()` 一行 | **只红 `TestAC2TaintSourceIsVisibleInTheUsageBlock`**（`:228`/`:231` 两条："help does not say the flag is an input fact"／"does not disclaim that the flag carries a conclusion"） | 未跑整包 |
| 还原 | `cp ../pristine-panel_assets-acc143.go` 回去 | — | 副本 md5 复量 `7f200f27…`＝锚点同值 |

M1 退回的那张卡（红因原文里的 stdout 逐字节，`acc2-M1-full.log`）：
`"level": "L1"`、`rulesHit` 只剩 `R1`、`reason` 只剩 `R1: 工具声明为下界（L1）`、`sessionOverrideBlocked: false`。
⇒ 实现方 §3 表里"卡退回 `L1 / [R1]`、其余 5 枚全绿"这一行：**复算成立**；M2 那行的"两枚红＋连带的另一枚一起红"：**复算成立**。

**补问的那条（编排者点名要判、实现方没给过读数的）：摘掉接线这一发，外部可见读数动过哪些？**

| 外部面 | M1 前后的读数 | 结论 |
|---|---|---|
| `-taint-source … -l2 notify …` 的 stdout | `L2 / [R1 R4] / sessionOverrideBlocked:true` → `L1 / [R1] / false`（本程用 `wisp-acc-M1.exe` 现跑现量） | **变了**——真伤，看得见 |
| `panel-assets -h` 的 stderr | 与未变异二进制 `diff`：只差第一行那条 `winsec: sealing path resolver installed` 的**时间戳**，usage 正文 14 行一字未动 | **没变**⇒ help 仍然承诺 `-taint-source` 能用，而接线已被摘掉 |
| 门禁四数 | `PASS 60→59`、`FAIL 0→1`（RUN/SKIP 不变） | **变了**⇒ 门确实有牙 |
| 不带新旗标的默认路径 stdout | A/B/C/D/P 五发与 M1 二进制仍逐字节相同 | 没变（本来就不该变） |

⇒ 判：**承重两发成立**（红因都能指到具体用例，票面 AC#3 要的"答得出是哪条用例"成立），
但**"help 会撒谎"这一支只有测试在守，外部面一条都不守**——这是接线这件事唯一的仪器，属结构性残余，不退回。

---

## §5 攻点 5：渲染那一行 + 新旗标的射程

### 5.1 `panel.NewApprovalCardView` 一字未加未改——用 git 证

```
git diff --name-only 9be3288^ HEAD -- internal/panel/approval.go      → 0 行
git log --format='%h %s' 9be3288^..HEAD -- internal/panel/            → 空（整段区间没人动过 internal/panel）
```

⇒ 渲染函数本体（`internal/panel/approval.go:64-67` 那三行 + `CardViewFromDecision` `:72-99`）**与本票无关、零改动**。
另加一枚交叉核：`go test -count=1 -run 'TestApprovalCardView' ./internal/panel/` 在实树现量 **ok**（rc=0）
⇒ 卡面键名与字节契约没被这条 CLI 路的改动碰坏（实现方 §6.1 说"没跑 internal/panel 的测试"，本程替它跑的就是这一条线，只这一条线）。

**但"渲染那一行一字未动"这句要说准**：调用行本身确实变了形，`git diff 9be3288^ 4037539 -- cmd/wisp/panel_assets.go` 的 hunk 逐字是

```
-		view := panel.NewApprovalCardView(risk.NewRiskAssessor(), panel.ApprovalSubject{
+		view := panel.NewApprovalCardView(assessor, panel.ApprovalSubject{
```

变的只有**第一个实参**（原来当场 new，现在 new 完可选地接一枚 detector 再传进去），
`panel.ApprovalSubject{…}` 那整块 subject 字面量（`:69-77`）与后面的 `enc.Encode(view)`（`:79-84`）**一字未动**，
`else` 分支（资产那几支）也一字未动。⇒ 票面 AC#2 那句"**一行渲染代码都不许加**"：**成立**（新增 0 行渲染）；
自证表 §1 的"渲染那行没动，仍是 `panel.NewApprovalCardView`"要读成"函数没动"才准确，**调用表达式本身是动过一枚实参的**——这是措辞精度，不是伤。

### 5.2 新旗标渗进 `-h`／`usage` 之外的地方了吗

- stdout 面：§3 的五发 `cmp` 已证**逐字节不变**（含完全不进 `-l2` 分支的 P 发）⇒ 没渗。
- rc 面：`-h` 前后都是 **2**；五发默认路径前后都是 **0** ⇒ 没渗。
- stderr 面：`-h` 的读数**从 2 行变成 14 行**（改前：一行 `winsec` 日志 + 一行 usage；改后：那两行 + **6 枚旗标的说明**）。这是本票唯一的对外可见变化。
  逐字 diff 已存 `/d/tmp/acc-143/hm1.diff` 之外的 `H-before.err`/`H-after.err`（两枚二进制现跑）。

**判"补 `fs.PrintDefaults()` 是不是顺带修了别的东西"**：

1. **必要性成立**：AC#2 要求旗标"在 `-h`/`usage` 里看得见（并且）说明它是'输入事实'不是'结论'"。
   `fs.Var` 的说明串只有 `PrintDefaults()` 才会印出来。自加变异 **M3**（摘掉那一行）现量：
   恰好 `TestAC2TaintSourceIsVisibleInTheUsageBlock` 一枚红，报的就是"help does not say the flag is an input fact"
   /"does not disclaim that the flag carries a conclusion"（`panel_assets_143_test.go:228`、`:231`）
   ⇒ **"那一改有自己的用例"这一条：成立**（它不是没测的顺带手）。
2. **射程确实比"最小"宽**：`PrintDefaults()` 是**一把**印全部旗标，于是 5 枚**既存**旗标的说明一并进了 `-h`，
   其中两枚的措辞本票没核过：`-l2` 那句 "print **the L2** card JSON"（实测不带 `-irreversible`/`-taint-source` 时打的是 `L1` 卡，
   措辞与实际不符，`C-after.json` 现量）与 `-irreversible` 那句 "(**R8**'s input)"。两句都是**票 114 时代的既存文案**
   （`git show 9be3288^` 里逐字同在），本票只是**把它们从注释级提到了界面级**。
   ⇒ 判：**不是"顺手改行为"**（行为面只有 help 变长），但**是"顺手扩大了 exposure"**，且这层 exposure 无断言（M3 只钉新旗标那三枚 token）。
   登记为**残余 R-143-a（不退回）**：要么把 `-l2` 那句 help 的 "the L2 card JSON" 改成不带承诺的说法（属 `cmd/wisp` 地界、一行的事、**要人工拍**因为改的是既存旗标文案），
   要么补一枚"`-h` 全集逐字快照"的用例把界面级说明钉住。本程**一字节没改**。

---

## §6 攻点 6：门禁与名册（自量，用的那一支写在标题里）

### 6.1 先把我用的是哪一把尺说死（编排者点名要的"要么带 DLL、要么改跑脚本"）

**两支都跑了**，两支持平：

| 尺 | 命令（原文） | 读数 |
|---|---|---|
| 甲：带 DLL 的整包 verbose | `PATH="/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx:$PATH" go test -count=1 -v ./cmd/wisp/` | 实树（内容＝锚点 `d242168`）**rc=0**、RUN=113 PASS=60 FAIL=0 SKIP=0；`ok github.com/CarlosShao/wisp/cmd/wisp 74.749s`（这行时长是命令自印，**不作判据**） |
| 乙：CI 那一步 | `sh scripts/wisp-cli-tests.sh` | **rc=0**；末两行逐字：`runtests.sh: OK - packages=[./cmd/wisp/ -count=1 -skip ^(…7 枚…)]: PASS=60 FAIL=0 SKIP=0, === RUN=113, '[no tests to run]'=0` ＋ `portable-tests.sh: four numbers (all from -v output): === RUN=113 --- PASS=60 --- FAIL=0 --- SKIP=0` |

**票面 `>` 更正③（"改前跑 `go test ./cmd/wisp/` 这条口令在本机字面跑不通"）——本程独立复量成立**：

```
$ go test ./cmd/wisp/            （不带任何 PATH 帮助，实树，HEAD=3855dc6）
exit status 0xc0000135
FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.035s
```

⇒ 这一条是**派单口令过期**，不是被验物的伤；本程所有 cmd/wisp 读数都是补过 DLL 路径之后取的。

### 6.2 改前 / 改后各一次（同一把尺、同一环境，两半都从 git 现构）

| 读数 | 版本来源 | RUN | 顶层 PASS | 顶层 FAIL | SKIP | panic |
|---|---|---|---|---|---|---|
| 改前 | `git archive 9be3288^` 副本（补 `frontend/dist` ＋ `third_party/`） | **101** | **54** | 0 | 0 | 0 |
| 改后（副本） | `git archive d242168` 同一副本还原成 md5 `7f200f27…` | **113** | **60** | 0 | 0 | 0 |
| 改后（实树交叉） | 工作树内容＝锚点（`git diff d242168..HEAD -- cmd/wisp/` 在 `0e17914`、`3855dc6` 两处都 0 行） | **113** | **60** | 0 | 0 | 0 |

⇒ 与实现方 §4.1 那三行（101/54、113/60、113/60）**逐枚同值**，但本程是自己重跑的，没引用它。
⚠ **一条流程坑要报回**（我踩过、实现方也绕过的同一形）：`git archive` 出来的副本**不带 `third_party/` 时整包会有 8 枚与本票无关的红**
（`TestAC2RealProcessRefusesOnEveryLegWithoutAppData128`、3 枚 `TestAC*Early*/ResidentLeg*`、`TestSecretArgvCarriesNoSecret`、`TestSecretRealBinaryRefusesValueFlag` 等，
红因逐字 `no native DLLs in ..\..\third_party\sherpa-onnx - run scripts/fetch-deps.ps1 first`）。
把"整包 rc 红"当结论就会把这 8 枚算到被验物头上——**本仓记过这一族**，此处补一枚现量：补过 `third_party` 之后 54/60 全绿。

### 6.3 点名册差集（四数之外的那一栏）

```
$ grep -E '^=== RUN   [A-Za-z]' <改前/改后两本 log> | sed 's/^=== RUN   //' | sort | diff -
只增不减，新增 12 行 = 6 枚顶层 + 6 枚子用例：
> TestAC1AC2TaintSourceLegProducesAJudgedR4
> TestAC1CmdSideEmitsNoVerdictTokens
> TestAC1MalformedTaintSourceIsRefused (+ /empty_content /empty_tool /no_separators /two_parts)
> TestAC2TaintFlagDoesNotLeakIntoTheDefaultCard
> TestAC2TaintSourceIsVisibleInTheUsageBlock
> TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit (+ /declared_but_absent_from_the_outgoing_call /fragment_below_the_contract_floor_cannot_match)
消失：0 枚
```

第三枚独立仪器同向：包内那条覆盖率闸门自己印的披露行，改前印 **`54 startable cases`**、改后印 **`60 startable cases`**
（`leg_dispatch_gate_133_test.go:244` 的 run-roster disclosure，副本与实树两处读数一致）。
包内 `^func Test` 声明数现量 **63** 枚（实现方 §4.1 那行的"声明 57、跑到 54"是改前的形状，57+6=63 对得上）。
`grep -c "t.Skip" cmd/wisp/panel_assets_143_test.go` = **0**；两本 log 的 `--- SKIP` = **0**；`panic` 字样 = **0**
⇒ 没有"改成 `t.Skip` 把读数抹掉"，也没有"一枚 panic 吞掉同包几十条"那种缺口——**本程不需要记任何"未取到"**。

### 6.4 另外三道门（本程现量，HEAD 见每一行）

| 门 | 命令 | HEAD | 读数 |
|---|---|---|---|
| 格式 | `gofmt -l cmd/wisp/` | `3855dc6` | **空** |
| vet | `go vet ./cmd/wisp/` | `3855dc6` | **空，rc=0**（没在根目录碰 `tools/d22scan/`，那是独立 module） |
| D22 | `sh scripts/d22scan.sh` | `3855dc6` | **rc=0**。第一步正控先过：`runtests.sh: OK - packages=[./...] top-level: PASS=30 FAIL=0 SKIP=0, === RUN=70`；第二步真扫描 `d22scan: clean - no D22 ban violations`，分母逐枚 `bans #1-5 internal/=203, bans #1-5 cmd/=22, ban #6 frontend/=46, ban #7 internal/tools/=18, ban #8 design/=32, ban #8 frontend/=46, ban #8 internal/=407, ban #8 cmd/=40`（ban #8 那几项逐字写着 "comments and _test.go included" ⇒ 本票新增的 `_test.go` 真在 40 枚射程里）；另印一行 `skipped as git-ignored: 1 file(s) … frontend/dist/assets/` |

⇒ **攻点 6 判成立**：四数（113/60/0/0）、名册差集（+6/−0）、三道门（空/空/rc=0）全部自量同值，
且用的是"带 DLL"那一支并把另一支（`scripts/wisp-cli-tests.sh`）也对了一遍。




