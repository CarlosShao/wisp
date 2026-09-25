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
