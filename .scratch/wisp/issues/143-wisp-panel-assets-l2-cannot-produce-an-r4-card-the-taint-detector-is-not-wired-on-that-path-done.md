# 143 — `wisp panel-assets -l2` 产不出 R4 卡：产出的那条路没接污染检测器，而消费它的一侧已在渲染

- Status: ready-for-agent
- 来源：前端会话 09-25 12:3x 对 `F2` 的顶回（`frontend-session-log.md` §46.1），编排者 12:4x 逐枚复量成立
- 关联：`A233`（§10 队列 `F2` 的前提就是本票）、`A232③(2)`、票 114（`-l2` 那条 CLI 路是它开的）、`Q-50`（面板侧新增的第六枚方法名，另一件事）
- 地界：**只做 `cmd/wisp/panel_assets.go`＋它自己的测试**。`internal/risk/**` 在禁改清单里，本票默认**不碰**（见 AC#1 的停手条件）。

## 为什么这枚票存在（不是"顺手补个 fixture"）

`PLAN.md:3481` 那张差距表里"聊天/SSE 之外的 L2 卡"一格举了 `R4：包含来自 <source> 的内容` 这句解释文字；
前端会话原本按我 11:6x 那句"补一枚 `rulesHit:["R4"]` 的 fixture、走 `wisp panel-assets` 同一条路"去做，
**做不下去**，因为那条路**根本产不出 R4**：

- `cmd/wisp/panel_assets.go:55` 构造的是**裸** assessor：`panel.NewApprovalCardView(risk.NewRiskAssessor(), …)`，**没有** `WithTaintDetector`。
- `internal/risk/rules_gateway.go:100-102` 的 `ruleTaint` 第一句就是 `if ctx.taint == nil { return nil }`，注释逐字写
  "Dormant until ticket 19 wires the TaintDetector"。
- 而生产里**是接了的**：`internal/tools/bridge.go:670` `WithTaintDetector(b.prov.Detector(taskID))`。

⇒ **净效果：一台演示机与一台生产机的判据不一样**——真判得出的规则，在产字节那条 CLI 路上是休眠的。
所以今天想要一张带 R4 的卡**只有一条路＝手写字节**，那正面撞 owner 的 P9 红线（"不得拿假件当真实字段"），
也正是我自己在 `A232③(2)` 里警告过它别做的事。**这一格不是前端的活，是我这边的地基缺了一段。**

## AC

- [x] **AC#1 先判"最小真修法"到底落哪一层，落不了就停手上报，不许自行假设**（D22 闸门③）。
      读 `internal/risk/provenance.go`＋`internal/tools/bridge.go:670`＋`assessor.go:205-212`，
      回答："要让 `wisp panel-assets -l2` 打出**由 assessor 真判出来**的 R4 卡，最少需要注入什么？
      它能在 `cmd/wisp` 侧构造（不碰 `internal/risk/**` 一字）吗？"
      ⓐ 若能在 cmd 侧用一个 cmd 内定义的小 detector 接上 ⇒ 走这一支，**且必须说清它算不算"假件"**：
        判据＝**R4 那句解释文字与那条 hit 是不是 assessor 自己产出的**；CLI 只允许提供"哪条通道携了哪个来源"这一**输入事实**，
        不允许把 `R4:` 字样或 `rulesHit` 数组写进输出。
      ⓑ 若必须动 `internal/risk/**` ⇒ **停手**，把要动的文件＋行＋为什么不够写进证据件并报回编排者（解冻＝人工批准）。
- [x] **AC#2 落地**：`-l2` 那条路新增一个旗标（名字自取，要在 `-h`/`usage` 里看得见，说明它是"输入事实"不是"结论"），
      用它注入 AC#1 选定的 detector；**输出仍走同一枚 `panel.NewApprovalCardView`，一行渲染代码都不许加**。
      判据：不带新旗标时输出**逐字节不变**（与改前同一命令的 stdout 做 `cmp`，这是"没渗到默认路径"那一条）。
- [x] **AC#3 承重自证**（不许只加断言）：在**仓外副本**里造两发变异——
      (i) 摘掉新加的接线 ⇒ 那枚"能打出 R4"的用例必须**变不响或变红**（答不出是哪条用例＝装饰，直接作废本票并登记为什么）；
      (ii) 把 detector 换成"永远报命中"⇒ 必须有用例红（防过度匹配；这条是本仓既有的一族，见 `A218` 那批）。
- [x] **AC#4 门禁与名册**：`go test ./cmd/wisp/` 改前改后各一次（**四数之外点名册差集**：新增几枚 `Test`、有无原有用例消失）、
      `gofmt -l` 空、`go vet ./cmd/wisp/` 空、`sh scripts/d22scan.sh` rc=0（⚠ 它是独立 module，别在根目录 `go vet ./tools/d22scan/`；
      bans #1-5 **不含** `_test.go`，只有 ban #8 含）。
      ⚠ `internal/risk/**`、`rules_gateway.go`、`thresholds.go`、golden、`docs/PLAN.md`、`docs/specs/**`、
      `frontend/**`、`design/**`、`tools/d22scan/**` **一律零字节**（契约轴要先证明这些东西存在，再说零命中，同一把尺跑一枚正控）。
      ⚠ `internal/panel/**` 此刻有另一程在写（门钉 r3）；`bridge.go` 与 `l2_grant_boundary_test.go` **都不要碰**，
      需要读它的判据就 `git show`，别改。
- [x] **AC#5 交回给前端那一侧的一句话**：新命令逐字长什么样（一条可直接粘贴的命令行＋它的 stdout 前两行），
      写进证据件末节；**不写 `frontend/**`、不替它建 fixture**。

## 规矩

只写 `cmd/wisp/panel_assets.go`、`cmd/wisp/*_test.go` 里新增的那一枚测试文件、以及 `docs/evidence/s1/143-panel-assets-r4-leg-*.md`。
逐格 commit（每裁完一格提交一枚），显式 pathspec（`git commit -q -F - -- <路径>`），commit 前现核 `git diff --cached --name-only`，
**出现别家路径就停手报回**；禁 `add -A`/`--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`/`push`；
变异落仓外副本、临时件只建不删；台账与票面的勾归编排者。
简报里任何一句前提（含"`internal/risk` 够用"、"这条 CLI 路只有 `panel_assets.go` 一处"）复算发现不符 ⇒ 写在证据件里并报回，**不要按我认为的样子改**。

---

> **票面三处前提被实现程实测推翻／说重了（09-25 14:0x 编排者落；原句一字不抹，出处＝`docs/evidence/s1/143-panel-assets-r4-leg-r1.md`，逐条我已独立复量）**
> **① 我引错了行，且我说的那句"差距表里的 R4 文案"不在我指的那一格**：票面正文写 `PLAN.md:3481` 那张表"举了 `R4：包含来自 <source> 的内容` 这句解释文字"。现量＝`PLAN.md:3481` 是 **`| **L2 确认卡（面板内）** | …顶部 2px --danger 横条…`**，**不是 R4**；R4 那一格在 **`PLAN.md:3488`**：`| **注入检出**（C25/R4） | --danger 竖条 + shield-alert(16px) + 「本次操作包含来自 \`<url>\` 的内容」（必须指明来源）+ 可展…`——**占位符是 `<url>`，不是我写的 `<source>`**。⇒ 要引 R4 那句界面文案引 `:3488`；`:3481` 引的是 L2 卡本体那一格（我在 `A232③(2)` 里引 `:3481` 指的确实是那一格，那处不欠改；欠改的是本票把两格混成一格）。
> **② "聊天/SSE 之外的 L2 卡"这个短语全仓只存在于这张票的措辞**（实现程 `grep` 现量），**不是规格里的一句话**。⇒ 我据此写的范围没有那个边界可依；实际范围＝票面点名的 `panel-assets -l2` 那一跳。本条按现量收窄。
> **③ 我给的 AC#4 口令"改前跑 `go test ./cmd/wisp/` 取基线"在本机字面跑不通**：不带 sherpa DLL 路径时**改前就是红的**（`0xc0000135`，票 98 那一族）。⇒ 不是被验物的错，是**我那句"这条先例在本机有效"的推论过期了**；今后凡让程取 `cmd/wisp` 基线，要同时给"这台机器要带 DLL"那一行，或让它改跑 `scripts/wisp-cli-tests.sh`。实现程实际拿到的凭据比我给的口令硬：**用两枚分别构建的二进制做 `cmp`**（`4712ea6`／`5ef1632` 各一遍）证"默认路径逐字节不变"。
> **④ "生产里是接了的"这句我说满了**（实现程发现、我 14:0x 复量成立）：`internal/tools/bridge.go:670` 确是**带 taint 的那一枚**，但 **`:193-196` 还有第二枚 assessor 没接** `WithTaintDetector`（只有 `WithCanonicalizer`＋`WithSensitiveClassifier`），而 `assessorFor`（`:663`）在 `injected || prov == nil || taskID == ""` 时返回的正是它；调用点 `:272`，配套 `:552` 在空 `taskID` 时连 `Mark` 都不做。⇒ **今天不是活洞**（agent loop 每条调用都带 `taskID`），但这句话今后要说成**"带 `taskID` 那条路接了"**。登记为独立缺口（归口见 `A241`），**本票一字节没碰它**。

>> **⑤ 上面那条"更正"自己也被更正了一次（09-25 14:3x 编排者；出处＝`143-panel-assets-r4-leg-accept-r1.md` 前提报回第 1 条，我 14:3x `sed -n '3481p'` 复核认可；① 的原句不抹）**
>> 我在 ① 里写"`PLAN.md:3481` 那一格不是 R4，R4 在 `:3488`"——**这句过正了**。现量 `:3481` 逐字含
>> `正文 = 完整参数（不是摘要）+ C19 命中的规则（如「R4：包含来自 \`web.fetch\` 的内容」）` ⇒ **那一格确实举了 R4 的文案**，只是占位符写成 `` `web.fetch` ``。
>> `包含来自` 在 `PLAN.md` 共 **4 枚命中**：`:2476`（`<源>`）／`:3026`（`<url>`）／`:3481`（`web.fetch`）／`:3488`（`<url>`）。
>> ⇒ 本票正文真正错的只有两处、且都不是行号：**把那张表叫成"差距表"**（规格里没这个称呼）、**占位符拼成 `<source>`**（纸上无此拼法）。
>> ⚠ 这条要留在仓里的理由：这是一枚**"更正比原错更错"**的实例——我把"哪一格有 R4"从我票面的一处，改口成了另一处，而两处都有。⇒ 以后落 `>` 更正前，**把被引那一行整行读完**，不要只核"我以为的那一格叫什么"。
