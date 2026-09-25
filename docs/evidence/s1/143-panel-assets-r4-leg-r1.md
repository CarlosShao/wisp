# 票 143 — `wisp panel-assets` 那条产字节的路现在能打出由 assessor 真判出来的 R4 卡（r1）

裁决者：本件作者＝实现者本人（票 143 的落地程）。**对抗验收与缺口审计必须另派一枚 agent**，
本件是交付说明，不是验收表（`SPEC-12 §4.3` #1/#3、D22 双角色）。

- 锚点自量：进场 `git rev-parse --short HEAD` = **`6103120`**（与简报 12:4x 给的同一枚）。
- 本程每一格读数都带它自己那次读数的 HEAD（下面逐格标）。共享树里 HEAD 一直在走：
  本文引用过的锚点 `038ed2b 81ad6fd 7d6a019 4712ea6 b293784 ae9f27f 051134b 6e5c1d2 9839b98 2dbb6af dd007d8 5ef1632 fe4a301`
  逐枚点名核过，全部落在 `git log --format=%h 6103120..HEAD`（锚点 `fe4a301` 现量，那一段共 23 枚，含别程的 commit，此处不逐枚点名）；
  `6f9a703` 与 `a2a7489` 两枚是读完那条链之后新落上去的锚点，`git rev-parse --short HEAD` 现量。
  成对读数（改前/改后）的两半取自不同 HEAD，所以每半单独标；这不是笔误，是并发。
- commit 落点（各枚 `git show --name-only` 现量都只带本程的路径，别家路径 0 枚）：
  - **代码两枚，逐枚点名**：`9be3288` `feat(票143 AC#1+AC#2)` = `cmd/wisp/panel_assets.go` + `cmd/wisp/panel_assets_143_test.go`（387 insertions / 4 deletions）；
    `4037539` `fix(票143 AC#2 追正)` = `cmd/wisp/panel_assets.go` 一枚（旗标说明按现量改写，见 §2.3）。这两枚不会再动。
  - **本件是追加式的若干枚，只改这一枚 md，所以不在这里点名枚数**（写枚数的那一行自己就会被下一枚改错，
    `8a5c691` 就是干过这件事的现量）。要名单就跑
    `git log --format="%h %s" -- docs/evidence/s1/143-panel-assets-r4-leg-r1.md`。

---

## §0 前提复算：七条 —— 四条成立、两条要说准、一条（票面自己的引用）不成立

| # | 简报那句前提 | 复算结果 | 凭据 |
|---|---|---|---|
| P1 | "能在 cmd 侧接住（AC#1 ⓐ）" | **成立**，`internal/risk/**` 一字未动 | §1；零命中读数在 §4.3 |
| P2 | "这条 CLI 路只有 `panel_assets.go` 一处" | **成立**：全仓 `NewApprovalCardView`/`CardViewFromDecision` 的生产调用点只有 `cmd/wisp/panel_assets.go` 那一处（改前 `:55`，本程改后同一句在 `:68`），其余三处调用加一处注释全在 `internal/panel/approval_test.go`（:8/:53/:92/:141） | `grep -rn "NewApprovalCardView\|CardViewFromDecision" --include=*.go .`，两枚锚点各读一次：`81ad6fd`（进场时的形状）与 `6f9a703`（本程落完之后） |
| P3 | "`ruleTaint` 只在 `ctx.taint == nil` 时休眠" | **成立**：`internal/risk/rules_gateway.go:102` 第一句就是 `if ctx.taint == nil {`（103 行 `return nil`），注释里那句 "Dormant until ticket 19 wires the TaintDetector" 逐字在，但**它跨 99/100 两行**（`... ticket 19 wires the` / `TaintDetector.`），按行号 grep 单行的人会漏 | 直读该文件，锚点 `6f9a703` 复量 |
| P4 | "`bridge.go:670` 那是生产唯一接线处" | **改前成立**（`grep -rn "WithTaintDetector(" --include=*.go . \| grep -v _test.go` 在生产侧只有三处：`internal/risk/assessor.go:207` 的定义、`internal/risk/provenance.go:51` 的一句注释、以及 `internal/tools/bridge.go:670` 的接线）；**本程落完之后它多了第二枚接线处，就是票 143 自己那枚**（`cmd/wisp/panel_assets.go:66`，锚点 `6f9a703` 现量）。**但"生产是接了的"这半句还要补一条：生产里另有第二枚 assessor 实例是没接的** | 见下面 P4-补 |
| P5 | "`go test ./cmd/wisp/` 改前改后各一次" | **这条口令在本机按字面跑不通**：不带 DLL 路径时它连加载都过不去，改前那一次就是红的，红因不是代码 | 见 §4.1 |
| P6 | （简报没提，我自己撞上的）旗标位置 | 新旗标与 `-l2` 的先后**不限**，但**必须在第一个工具参数之前**；否则整串被吞进 argv | 见 §2.3 现量 |
| P7 | **票面自己那句引用不成立**：「`PLAN.md:3481` 那张差距表里"聊天/SSE 之外的 L2 卡"一格举了 `R4：包含来自 <source> 的内容`」 | 三处对不上：① `PLAN.md:3481` 是 D29 那张**状态-视觉-动效**表里的 `L2 确认卡（面板内）` 那一行，不叫"差距表"；② 那句 R4 文案在同一张表的 **`:3488`（`注入检出`（C25/R4）那一行）**，写作「本次操作包含来自 **`<url>`** 的内容」（占位符是 `<url>`，不是 `<source>`），另两处在 `:2476` 与 `:3026`；③ 全 `docs/` 与 `.scratch/` 里 "聊天/SSE 之外的 L2 卡" 这个短语**只存在于这张票自己的措辞**（`grep -rn` 一枚命中，就是票面） | 锚点 `6f9a703` 现量：`grep -n "包含来自" docs/PLAN.md`、`sed -n '3481p;3488p'`、`grep -rn "聊天/SSE" docs/ .scratch/`；`docs/PLAN.md` 最后一次改动是 `45623e4`（11:01，早于本票开），**不是行号腐坏，是那格本来就指错了** |

### P4-补：生产里那枚"没接的"实例（不改，只登记）

`internal/tools/bridge.go:663-671` 的 `assessorFor(taskID)` 开头是
`if b.injected || b.prov == nil || taskID == "" { return b.assess }`，
而 `b.assess` 在 `bridge.go:194-196` 构造时只接了 `WithCanonicalizer` + `WithSensitiveClassifier`，
**没有** `WithTaintDetector`。调用点是 `bridge.go:272`（`verdict := b.assessorFor(req.TaskID).Assess(...)`）。
⇒ 净效果：**一条 `TaskID` 为空的调用在生产里同样判不出 R4**（配套的另一半：`bridge.go:552` 在
`dec.TaskID == ""` 时连 `Mark` 都不做，所以那枚 scope 本来就是空的）。
今天不是活洞——`internal/agent/loop.go` 每条工具调用都带 taskID（`grep "TaskID:"` 生产侧命中 `loop.go:644` 等，非空）——
但"生产是接了的"这句要说准：**接的是 per-task 那枚实例，不是所有实例**。
归口建议：与票 143 无关，动的是 `internal/tools/**`，要开就另开，本程一字节没碰。

### 另一条要看的人需要知道的（工作树状态）

- `git status --porcelain` 里那 16 枚 ` D design/**` 与 4 枚 `?? design/...` **不是本程造的**：
  我在 `6103120` 的第一枚读数里它们就已经在了（原文贴在 §4.3）。

---

## §1 AC#1：最小真修法落在 ⓐ，而且接的是真引擎不是小 cmd 件

**问题**："要让 `wisp panel-assets -l2` 打出由 assessor 真判出来的 R4 卡，最少需要注入什么？"

答案：三样，全在 `internal/risk` 的**导出口**上，本程一行都没碰 `internal/risk/**`：

1. `risk.NewProvenance(risk.ProvOptions{NoProbe: true})` — C25 引擎本体；
2. `OpenScope(scopeID)` + 逐条 `Mark(scopeID, tool, origin, content)` — 声明"这次任务读过哪些敏感来源"；
3. `assessor.WithTaintDetector(prov.Detector(scopeID))` — 冻结的 C19 接缝，与 `bridge.go:670` **同一个构造函数**。

**为什么 ⓑ（要解冻 `internal/risk`）根本不需要**：`ruleTaint` 缺的只是 `ctx.taint` 非空，
而 `TaintDetector` 是消费侧接口（`assessor.go:165-167`）、`WithTaintDetector` 是导出的装配口（`assessor.go:207`）。
休眠是**装配缺失**，不是**能力缺失**。

**ⓐ 那一支要求我说清"算不算假件"。判据是 owner 给的那条死的**：那句解释文字与那条 hit 是不是 assessor 自己产出的。

- CLI 给的：`<source-tool>|<origin>|<content>` — 只有来源事实，没有等级、没有规则号、没有那句话；
- assessor 产的：`rulesHit` 里的 `R4`、`level = L2`、`reason` 里那句 `R4: 包含来自 <来源> 的内容`、
  以及 `sessionOverrideBlocked = true`（这个 bool 只有 `rules_gateway.go:113` 会置）；
- 渲染那行没动：`panel.NewApprovalCardView(...)` 一字未改（`internal/panel/**` 零字节，见 §4.3）。
- **并且**：来源名字要出现在卡上，唯一的路径是它被写进了 C25 的 provenance 记录，
  再由 `Hit.Source()`（`provenance.go:257-262`）交给 `ruleTaint` 的 `fmt.Sprintf`。
  这条通路是**用例钉住的**（§3 那枚 `TestAC1AC2...` 断言 reason 含声明的 origin），
  不是本件口头承诺的。

**没有降级成 cmd 内自定义小 detector**：`cmd/wisp` 里没有任何实现 `TaintHit` 的类型，
`detector()` 返回的就是 `prov.Detector(...)` 那个 `risk.boundDetector`。这一点是 AC#1 ⓐ
"用一枚 cmd 内定义的小 detector"的**加强版**——比简报允许的形状更真，所以顺手记一句：
如果哪天有人在这条路上塞一个返回常量的 detector，§3 那两枚用例（R4 要含声明的来源、
不过度匹配）会一起红。

**`NoProbe` 那一处偏差，写明并给出不改变判据的理由**：生产桥的 `Provenance` 是全配的（会探 P12 同步目录），
本 CLI 探针用 `NoProbe: true`，为了不去碰注册表/环境变量探针。它**改变不了这条路的结论**，理由是形状而不是推导：
面板把 argv 变成 params 的函数只产出 `command` 与 `argv` 两枚 key（`internal/panel/approval.go:107-114`），
两枚都不是 `pathKeys`（`provenance.go:143`）里的写目标，所以 `writeGate("", params)` 在任何 sync 目录集合下都返回
`hasPath == false` → `open == true`（`provenance.go:660-668`）→ 走通用扫描。同步目录那条唯一带落点条件的通道在这条路上根本不参与。
现量：§2.4 那张卡命中的是通用扫描，reason 里点名的是声明的来源。

---

## §2 AC#2：落地与"没渗进默认路径"

### 2.1 旗标

```
-taint-source <source-tool>|<origin>|<content>      （可重复）
```

`-h` 里可见（锚点 `5ef1632` 现量的**最终**版本，即 `4037539` 那枚落下去的文字；两行只是节选）：

```
usage: wisp panel-assets [-manifest] [-check] [-render <path>] [-taint-source <source-tool>|<origin>|<content>]... [-irreversible <ops>] -l2 <tool> <args...>
  -taint-source value
    	with -l2: repeatable INPUT FACT of the shape <source-tool>|<origin>|<content>, meaning "this task read <content> from <source-tool> at <origin>". It feeds the C25 provenance engine the taint rule judges; it declares no verdict, no rule id and no level, and the taint rule stays dormant without it. Every flag has to come before the tool's own arguments: parsing stops at the first argument that does not start with a dash, and everything after it is argv.
```

顺手把 `fs.Usage` 加了 `fs.PrintDefaults()`：原来那枚 Usage 只印一行，**任何旗标的说明文字都到不了 `-h`**，
"要在 `-h` 里看得见、说明它是输入事实"这句在没有它的情况下做不到。这是 `-h` 那**一路**输出的变化，
不是 `-l2` 那一路（§2.2 的 cmp 只管 stdout 的卡字节）。

形状校验在 `Set` 里：四枚畸形值一律 rc=2 且不印卡（§3 的 `TestAC1MalformedTaintSourceIsRefused` 钉住）。

### 2.2 不带新旗标时输出逐字节不变（AC#2 的判据）

改前二进制：`go build -o /d/tmp/ticket-143/wisp-before.exe ./cmd/wisp/`，锚点 `7d6a019`，工作树 `cmd/wisp/` 干净。
改后二进制：`wisp-after.exe`，同一棵树之后（锚点 `4712ea6`，`git status --porcelain -- cmd/wisp/` 只有本程那两枚文件）。

| # | 命令（都走 `-l2` 那一路） | 读数 |
|---|---|---|
| A | `panel-assets -l2 shell.run rm -rf D:/tmp/x` | `cmp outA-before.json outA-after.json` -> **identical**（rc 都是 0） |
| B | `panel-assets -irreversible delete,overwrite -l2 fs.write D:/notes/a.txt` | **identical**，卡是 `L2 / [R1 R8]` |
| C | `panel-assets -l2 notify hello world` | **identical**，卡是 `L1 / [R1]` |

现量命令与读数原文（改前二进制 `wisp-before.exe` 建自锚点 `7d6a019`；改后跑了两遍，因为 `4037539` 追正了说明文字，两遍都要有读数）：

```
$ for x in A B C; do cmp out$x-before.json out$x-after.json && echo "cmp $x: identical"; done     # wisp-after.exe，锚点 4712ea6
cmp A: identical
cmp B: identical
cmp C: identical
$ for x in A B C; do cmp out$x-before.json out$x-after2.json && echo "cmp2 $x: identical"; done   # wisp-after2.exe，锚点 5ef1632（= 4037539 的内容）
cmp2 A: identical
cmp2 B: identical
cmp2 C: identical
```

⇒ 三发默认路径 stdout 与改前**逐字节**相同，两枚二进制都过；`rc` 三发全 0。

### 2.3 一条旗标位置的实测（简报没写，写下来省得下一个人心算）

Go 的 flag 包在**第一个不以 `-` 开头的参数**处停止解析，`-l2` 的值不算 positional：

- `-l2 fs.write -irreversible delete D:/tmp/x` -> R8 照样命中（值被 flag 吃掉，argv 只剩 `D:/tmp/x`）；
- `-l2 notify 第一句 -taint-source "web.fetch|u|..."` -> **整串进了 argv**（`args` 三行含 `-taint-source`），卡是 `L1 / [R1]`。

⇒ 结论：**旗标可以写在 `-l2` 后面，只要写在任何工具参数前面**；usage 行给的"全部在前"是最不会错的那一版。

### 2.4 新能力真的产出了一张由 assessor 判出来的卡（锚点 `5ef1632`，stdout 逐字节）

```json
{
  "correlationId": "panel-assets-l2",
  "tool": "notify",
  "args": [
    "把 合同编号 HT-2026-0731-KX 发到远端"
  ],
  "level": "L2",
  "rulesHit": [
    "R1",
    "R4"
  ],
  "reason": "R1: 工具声明为下界（L1）; R4: 包含来自 web.fetch https://files.example.com/q3-notes.txt 的内容",
  "reasonKnown": true,
  "sessionOverrideBlocked": true,
  "callChain": [
    "cli",
    "panel-assets",
    "notify"
  ],
  "decidedBy": "native"
}
```

同一枚卡在 `4712ea6`（`9be3288` 的内容）上也读过一次，形状一模一样；`4037539` 只动了说明文字。
`sessionOverrideBlocked: true` 只可能来自 `rules_gateway.go:113`——这条是"这句话不是 CLI 抄的"最硬的一枚旁证。

---

## §3 AC#3：承重自证，两发变异都落在仓外副本

副本：`/d/tmp/ticket-143/mut`（`go.mod`/`go.sum` + `cmd/` + `internal/` + `frontend/{embed.go,dist}`，5.8 M）。
**仓库树全程没进过变异态**：变异只 `sed`/`perl` 副本里的文件，每发之后复算。两枚内容各跑一整轮
（第一遍 = `9be3288` 的落件，第二遍 = `4037539` 追正说明文字之后的落件），四份读数都留在 `mut1.txt / mut2.txt / mut1b.txt / mut2b.txt`：

```
$ git -C "/d/work/workspace/projects plans/Wisp" status --porcelain -- cmd/wisp/
 M cmd/wisp/panel_assets.go                      <- 第一遍时：本程的改动（当时未提交）
?? cmd/wisp/panel_assets_143_test.go             <- 就这两枚，没有第三枚
（第二遍同一命令的读数： M cmd/wisp/panel_assets.go            —— 测试文件那时已提交，只剩这一枚）
$ md5sum cmd/wisp/panel_assets.go /d/tmp/ticket-143/pristine/<备份> /d/tmp/ticket-143/mut/cmd/wisp/panel_assets.go
4789e18f6c0ed42e68eeb10e8c77ad9c   第一遍三处同一个值（9be3288 内容）
7f200f275ed8ea854e84a104ea0fc58f   第二遍三处同一个值（4037539 内容）
$ grep -c MUTANT cmd/wisp/panel_assets.go
0                                   仓库树里从没出现过 MUTANT 字样
```

⇒ 三处 md5 相同 = 仓里那枚、变异前的备份、变异后又还原的副本是同一份内容。

| 发 | 内容 | 结果（两遍同一形，读数分别 `mut1/mut1b`、`mut2/mut2b`） |
|---|---|---|
| 控制 | 副本、未变异 | `go test -count=1 -run 'Taint\|CmdSide' ./cmd/wisp/` -> `ok`（rc=0；两遍都过） |
| **M1 摘接线** | `s/assessor = assessor.WithTaintDetector(det)/_ = det \/\/ MUTANTI/` | rc=1，**唯一红**：`--- FAIL: TestAC1AC2TaintSourceLegProducesAJudgedR4`；卡退回 `level` 为 `L1`、`rulesHit` 只剩 `R1`、`reason` 只剩 `R1: 工具声明为下界（L1）`（`mut1.txt` 第 12-16 行是那段 stdout 原文）；其余 5 枚（默认路径、过度匹配、usage、畸形值、字面量扫描）全绿 |
| **M2 换永真 detector** | `detector()` 末尾改返回一枚恒 `true` 的 `alwaysHitTaint143`（`prov` 变 `_ = prov`） | rc=1，**红两枚**：`--- FAIL: TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit`（**两个子用例都红**：`declared_but_absent_from_the_outgoing_call` 与 `fragment_below_the_contract_floor_cannot_match`）**以及** `TestAC1AC2TaintSourceLegProducesAJudgedR4` 一起红（永真 detector 报的来源不是声明的那枚，reason 里点不出 origin）；默认路径那枚仍绿（没接线时就该绿，这是对的） |
| 还原 | 用备份 `cp` 回去 | `ok`（rc=0），md5 三处相同 |

**答得出"是哪条用例红"**：M1 = `TestAC1AC2TaintSourceLegProducesAJudgedR4`；M2 = `TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit`。
临时件全部留下（`mut1.txt`、`mut2.txt`、`mut1b.txt`、`mut2b.txt`、`pristine/`、`mut/`），只建不删。

新加的 6 枚顶层用例（`cmd/wisp/panel_assets_143_test.go`）：

1. `TestAC1AC2TaintSourceLegProducesAJudgedR4` — R4 由 assessor 判出，reason 点名声明的来源，`sessionOverrideBlocked` 为真；
2. `TestAC2TaintFlagDoesNotLeakIntoTheDefaultCard` — 不给旗标时 stdout 里**不含** `R4`、`rulesHit` 恰为 `[R1]`（仓内进程内的那半，字节级那半在 §2.2）；
3. `TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit` — 过度匹配两形（来源声明了但调用没带 / 来源短于 8 字符契约地板）；
4. `TestAC2TaintSourceIsVisibleInTheUsageBlock` — `-h` 里看得见旗标，说明文字含 "INPUT FACT" 与 "no verdict"；
5. `TestAC1MalformedTaintSourceIsRefused` — 四枚畸形值 rc=2、不印卡；
6. `TestAC1CmdSideEmitsNoVerdictTokens` — **P9 红线落成仪器**：`go/ast` 扫 `panel_assets.go` 的**每一枚字符串字面量**，出现 `R4` / `R9` / `包含来自` / `rulesHit` / `sessionOverrideBlocked` / `L2:` 即红（注释豁免，字面量不豁免——和 ban #8 同一口径）。

---

## §4 AC#4：门禁、名册、契约轴

### 4.1 `go test ./cmd/wisp/` 改前改后各一次（先报一条前提不成立）

**按字面跑不通**（锚点 `038ed2b`，同一枚在 `81ad6fd` 复现）：

```
$ go test ./cmd/wisp/
exit status 0xc0000135
FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.035s
```

这是票 98 登记过的那枚坑（`scripts/wisp-cli-tests.sh` 头部逐字写着：测试二进制链 sherpa-onnx 的 cgo
导入库，Windows 在 Go 第一行之前就把进程杀了），不是本程造的伤：`go build`、`go test -c` 都 rc=0，
把编译好的二进制直接跑，报的是
`error while loading shared libraries: sherpa-onnx-c-api.dll`。
⇒ 改前那一枚必须**先按 CI 那一步把 `third_party/sherpa-onnx` 交给 PATH**才拿得到（`go test -c` 之后跑 exe 也一样要它）。
本程两枚读数都用同一把尺（`PATH=<repo>/third_party/sherpa-onnx go test -count=1 -v ./cmd/wisp/`），锚点各自标。

| 读数 | 锚点 | `=== RUN`（含子用例） | 顶层 `--- PASS` | 顶层 `--- FAIL` | `--- SKIP` | panic |
|---|---|---|---|---|---|---|
| 改前 | `81ad6fd` | 101 | 54 | 0 | 0 | 0 |
| 改后（第一遍，`9be3288` 内容） | `b293784` | 113 | 60 | 0 | 0 | 0 |
| 改后（第二遍，`4037539` 内容） | `5ef1632` | 113 | 60 | 0 | 0 | 0 |

（`grep -c panic` 三枚读数都是 0 ⇒ 没踩"一枚 panic 吞掉同包几十条"那枚坑；
`--- SKIP` 三枚都是 0 ⇒ 本程**没有**把任何用例改成 `t.Skip`。）

**点名册差集**（`diff roster-before.txt roster-after2.txt`，只增不减；第一遍 `roster-after.txt` 差集逐字相同）：

```
> TestAC1AC2TaintSourceLegProducesAJudgedR4
> TestAC1CmdSideEmitsNoVerdictTokens
> TestAC1MalformedTaintSourceIsRefused
> TestAC2TaintFlagDoesNotLeakIntoTheDefaultCard
> TestAC2TaintSourceIsVisibleInTheUsageBlock
> TestAC3TaintSourceThatTheCallDoesNotCarryIsNotAHit
```

新增 6 枚、消失 0 枚、子用例 47 -> 53（+6，全在 3/5 两枚里）。
另有一条独立交叉核对：包内声明的 `func Test` 57 枚、跑到的 54 枚，那 3 枚是
`TestAC1POSIXSecretRouteSymlinkedHomeBecomesSealable119`、
`TestAC1POSIXSecretRouteSymlinkedXDGConfigHomeBecomesSealable119`、
`TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119`——`comm -13` 反向为空，
即**没有任何"跑到了但没声明"的用例**，差集是本机的平台约束，不是本程造成的。

### 4.2 另外三道门（现量）

| 门 | 命令 | 锚点 | 读数 |
|---|---|---|---|
| 格式 | `gofmt -l cmd/wisp/` | `b293784` 与 `5ef1632`（两遍） | **空**（改前 `6103120` 也空） |
| vet | `go vet ./cmd/wisp/` | `b293784` 与 `5ef1632`（两遍） | **空**（没在根目录跑 `go vet ./tools/d22scan/`，那是独立 module，按简报是设计不是伤） |
| D22 | `sh scripts/d22scan.sh` | 第一遍 `ae9f27f`、第二遍 `5ef1632`、第三遍 `4350bc2`（本件已入树之后） | 三遍都 **rc=0**；第一步正控先过：`runtests.sh: OK ... top-level: PASS=30 FAIL=0 SKIP=0, === RUN=70`；第二步真扫描 `d22scan: clean - no D22 ban violations`，分母三遍同为 `bans #1-5 internal/=203, cmd/=22, ban #6 frontend/=46, ban #7 internal/tools/=18, ban #8 design/=32, frontend/=46, internal/=407, cmd/=40`（ban #8 那几项逐字写着 "comments and _test.go included" ⇒ 我新增的 `_test.go` 真被扫过；`docs/` 不在 ban #8 的四棵树里，本件那枚 md 不是它的射程）。**下一行这句话本身在 `4350bc2` 之后，没再跑第四遍** |
| CI 那一步 | `sh scripts/wisp-cli-tests.sh` | `6e5c1d2` | **rc=0**，`portable-tests.sh: four numbers: === RUN=113 --- PASS=60 --- FAIL=0 --- SKIP=0`（与 §4.1 那枚改后用另一把尺独立复现；它的 `-skip` 名单里 5 枚名字在 `cmd/wisp` 里 `grep "func <name>("` 全部无命中，所以两把尺的分母是同一套。这一遍跑的是 `9be3288` 的内容，`4037539` 之后没重跑它——重跑的是 §4.1 那把 `-count=1 -v` 尺） |

### 4.3 契约轴：先证明那些东西存在，再说本程零命中，同一把尺跑一枚正控

锚点 `051134b`。分母（`git ls-files`，不是 `ls`）：

```
internal/risk      37 files tracked      （含 rules_gateway.go / assessor.go / provenance.go / taintmatch.go）
internal/panel     18 files tracked
tools/d22scan       6 files tracked
frontend/src       25 files tracked
docs/specs         14 files tracked
rules_gateway.go   internal/risk/rules_gateway.go
thresholds.go      internal/observe/thresholds.go
golden files       58 tracked（git ls-files -- '*golden*'）
PLAN.md            docs/PLAN.md
```

同一把尺、只把 pathspec 换成上面那一串：

```
$ git status --porcelain -- internal/risk internal/panel tools/d22scan frontend docs/PLAN.md docs/specs design internal/observe/thresholds.go
 D design/assets/base.css   ... （共 16 枚 D）
?? design/doubao/01-ball-states.jpg
?? design/doubao/demo/lib/
?? design/doubao/demo/screenshots/
?? design/old/
```

⇒ `internal/risk`、`internal/panel`、`tools/d22scan`、`frontend`、`docs/PLAN.md`、`docs/specs`、
`thresholds.go`：**这一格里零行**。出现的 20 行全是 `design/**`，
**不是本程造的**：进场第一枚读数（锚点 `6103120`）里那 16 枚 ` D design/...` 与那 4 枚 `?? design/...` 就已经在了，
本程从没往 `design/**` 写过一字节。

golden 与 `thresholds.go` 不在上面那串 pathspec 的射程里（它们在 `internal/agent[/testdata]`、
`internal/llm[/openaichat|/testdata]`、`tools/d22scan` 六处目录下，`git ls-files -- '*golden*' | sed 's|/[^/]*$||' | sort -u` 现量），
所以单独补两把同形状的尺（锚点 `7ea0b31`）：

```
$ git status --porcelain -- '*golden*'          # 58 枚 tracked golden
$ git status --porcelain -- internal/observe/   # thresholds.go 所在目录
（两条读数都是空）
```

正控（同一条命令换 pathspec，证明这把尺看得见改动；锚点 `051134b`，那两枚当时未提交）：

```
$ git status --porcelain -- cmd/wisp/
 M cmd/wisp/panel_assets.go
?? cmd/wisp/panel_assets_143_test.go
```

`internal/panel/l2_grant_boundary_test.go`：**没碰过**（那是另一程的门钉 r3）。
`internal/panel/approval.go` 与 `internal/tools/bridge.go`：**只读**——用 Read 与 `sed -n '<区间>p'`
看过判据（§1 那几处行号引用就是从这两条读法来的），**写入 0 字节**；上面 §4.3 的 pathstatus 读数就是这条陈述的证明。
本程落地 commit 的 `--name-only` 只有两枚 `cmd/wisp/` 路径（见本文顶部），别家路径一枚都没有；
`git add` 与 `git commit` 之间另一程挤进来一枚 `9839b98`，我的 commit 仍只带我的 pathspec——
这条就是"我带了 pathspec 不是豁免"那位的现量：`git diff --cached --name-only` 在 commit 前一刻核过，两枚路径。

### 4.4 没做的测量

- **没有任何计时类判据**。§4.1/§4.2 表里那些 `0.035s / 68.593s / 11.966s` 是命令自己印出来的原文，
  当下另有 2-3 程在跑测试，本程**不拿它们作任何判据**；承重判断只用了 rc、颜色、枚数、字节（`cmp`/`md5`）。

---

## §5 AC#5：交回给前端那一侧的一句话（只在本件落这一处）

一条可直接粘贴的命令行（在仓库根；`wisp` 用**从 `4037539` 或之后**建出来的二进制——
仓里那枚 `build/wisp.exe` 是旧的，实测它回 `flag provided but not defined: -taint-source`）：

```
wisp panel-assets -taint-source "web.fetch|https://files.example.com/q3-notes.txt|合同编号 HT-2026-0731-KX" -l2 notify "把 合同编号 HT-2026-0731-KX 发到远端"
```

⚠ **`go run ./cmd/wisp` 这一形在本机必须先交 DLL 路径**，否则进程连启动都过不去（票 98 那枚坑，实测）：

```
$ go run ./cmd/wisp panel-assets -l2 notify hello            -> rc=1, stderr: exit status 0xc0000135
$ PATH="<repo>/third_party/sherpa-onnx:$PATH" go run ./cmd/wisp panel-assets -taint-source "..." -l2 notify "..."
                                                            -> rc=0，下面那两行一模一样
```
（两枚读数都在锚点 `a2a7489` 现量。）

它的 stdout 前两行（锚点 `5ef1632`、`wisp-after2.exe` ＝ `4037539` 那枚内容 `go build` 出来的二进制现量，`rc=0`；
同一枚命令在 `4712ea6`、`2dbb6af`、`a2a7489`（`go run` 那一形）三处也各读过一次，前两行一模一样）：

```
{
  "correlationId": "panel-assets-l2",
```

要点四条：

1. 旗标必须写在**第一个工具参数之前**（§2.3）；
2. 卡上的 `rulesHit` 会含 `R4`、`level` 是 `L2`、`sessionOverrideBlocked` 是 `true`、
   `reason` 里点名的来源就是旗标声明的那枚；不带 `-taint-source` 就退回票 143 之前那枚 `L1 / [R1]`；
3. 对齐 `PLAN.md:3488` 那一格的要求（「本次操作包含来自 `<url>` 的内容（**必须指明来源**）」）：
   `reason` 里的源槽位是 `"<源工具> <origin>"`，上面那条命令里就是
   `web.fetch https://files.example.com/q3-notes.txt`——URL 在里面，所以那一格的判据满足；
   卡上另有一枚 `sessionOverrideBlocked` 布尔，对应 D45 那句"这条不许被会话授权盖掉"；
4. **`Hit.Fragment`（命中的那段敏感内容）不会出现在卡上**，也不该出现——`provenance.go:60-61`
   逐字写着 "Hit.Fragment itself must NEVER be persisted (it is sensitive content; native card only)"，
   字段注释在 `:251-252`；`panel.ApprovalCardView` 也没有承载它的键，所以本程无从泄漏，卡上只有来源名。

**本程没有为它建任何 fixture、没有写 `frontend/**` 一字节**（票面地界）。

---

## §6 本程没测什么（留给对抗验收的清单）

1. **没验证这条 CLI 路的字节契约与前端 `panel.ts` 的键一致**——那是 `internal/panel` 的
   `TestApprovalCardViewJSONKeysMatchFrontendTypes` 的活儿，本程没跑 `internal/panel` 的测试（另一程在写那棵树）。
2. **没测 `-taint-source` 与真实 R2/R3 装配的合体**：这条路仍然没接 `WithCanonicalizer`/`WithSensitiveClassifier`
   （票 143 之前就是这样，本程没改），所以如果还要一张同时含 R2/R3/R4 的卡，
   **今天仍然产不出**，得另开一格。
3. **没测多来源命中顺序**：C25 的"最早 mark 优先"（`provenance.go:463` 注释）由 `internal/risk` 自己的测试钉，
   本程只声明了可重复传，没用两枚来源断过序。
4. **没测 `NoProbe` 之外那条全配 `Provenance`（带真同步目录探针）的路**：§1 结尾那句"不改变判据"是
   **形状推导**（params 只有 `command`/`argv`），不是一引量到的。
5. **没在 linux 侧跑过**：`cmd/wisp` 的 linux 半边票 111 记着 19 枚红，与本程无关，本程也没碰。
6. **没测并发**：同包另有程在跑测试时采的数都是 rc/颜色/枚数/字节类读数，但没做过"这条命令同时被两程跑"的一致性断言。
7. **没跑全仓 `go test ./...`**：另两程此刻在写 `internal/panel/**` 与裁 `tools/d22scan/**`，
   全仓读数会把他们的中间态混进来。本程的门只到 `./cmd/wisp/` 与 `scripts/d22scan.sh`（后者自己带全仓正控）。
8. **对抗验收没做**：本件是实现者自述。§0 那七条前提的复算、§3 两发变异的重跑、§4 三门的重跑，
   请验收者**另起 agent 现量**，不要引用本件的读数当自己的读数。
