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

- [ ] **AC#1 先判"最小真修法"到底落哪一层，落不了就停手上报，不许自行假设**（D22 闸门③）。
      读 `internal/risk/provenance.go`＋`internal/tools/bridge.go:670`＋`assessor.go:205-212`，
      回答："要让 `wisp panel-assets -l2` 打出**由 assessor 真判出来**的 R4 卡，最少需要注入什么？
      它能在 `cmd/wisp` 侧构造（不碰 `internal/risk/**` 一字）吗？"
      ⓐ 若能在 cmd 侧用一个 cmd 内定义的小 detector 接上 ⇒ 走这一支，**且必须说清它算不算"假件"**：
        判据＝**R4 那句解释文字与那条 hit 是不是 assessor 自己产出的**；CLI 只允许提供"哪条通道携了哪个来源"这一**输入事实**，
        不允许把 `R4:` 字样或 `rulesHit` 数组写进输出。
      ⓑ 若必须动 `internal/risk/**` ⇒ **停手**，把要动的文件＋行＋为什么不够写进证据件并报回编排者（解冻＝人工批准）。
- [ ] **AC#2 落地**：`-l2` 那条路新增一个旗标（名字自取，要在 `-h`/`usage` 里看得见，说明它是"输入事实"不是"结论"），
      用它注入 AC#1 选定的 detector；**输出仍走同一枚 `panel.NewApprovalCardView`，一行渲染代码都不许加**。
      判据：不带新旗标时输出**逐字节不变**（与改前同一命令的 stdout 做 `cmp`，这是"没渗到默认路径"那一条）。
- [ ] **AC#3 承重自证**（不许只加断言）：在**仓外副本**里造两发变异——
      (i) 摘掉新加的接线 ⇒ 那枚"能打出 R4"的用例必须**变不响或变红**（答不出是哪条用例＝装饰，直接作废本票并登记为什么）；
      (ii) 把 detector 换成"永远报命中"⇒ 必须有用例红（防过度匹配；这条是本仓既有的一族，见 `A218` 那批）。
- [ ] **AC#4 门禁与名册**：`go test ./cmd/wisp/` 改前改后各一次（**四数之外点名册差集**：新增几枚 `Test`、有无原有用例消失）、
      `gofmt -l` 空、`go vet ./cmd/wisp/` 空、`sh scripts/d22scan.sh` rc=0（⚠ 它是独立 module，别在根目录 `go vet ./tools/d22scan/`；
      bans #1-5 **不含** `_test.go`，只有 ban #8 含）。
      ⚠ `internal/risk/**`、`rules_gateway.go`、`thresholds.go`、golden、`docs/PLAN.md`、`docs/specs/**`、
      `frontend/**`、`design/**`、`tools/d22scan/**` **一律零字节**（契约轴要先证明这些东西存在，再说零命中，同一把尺跑一枚正控）。
      ⚠ `internal/panel/**` 此刻有另一程在写（门钉 r3）；`bridge.go` 与 `l2_grant_boundary_test.go` **都不要碰**，
      需要读它的判据就 `git show`，别改。
- [ ] **AC#5 交回给前端那一侧的一句话**：新命令逐字长什么样（一条可直接粘贴的命令行＋它的 stdout 前两行），
      写进证据件末节；**不写 `frontend/**`、不替它建 fixture**。

## 规矩

只写 `cmd/wisp/panel_assets.go`、`cmd/wisp/*_test.go` 里新增的那一枚测试文件、以及 `docs/evidence/s1/143-panel-assets-r4-leg-*.md`。
逐格 commit（每裁完一格提交一枚），显式 pathspec（`git commit -q -F - -- <路径>`），commit 前现核 `git diff --cached --name-only`，
**出现别家路径就停手报回**；禁 `add -A`/`--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`/`push`；
变异落仓外副本、临时件只建不删；台账与票面的勾归编排者。
简报里任何一句前提（含"`internal/risk` 够用"、"这条 CLI 路只有 `panel_assets.go` 一处"）复算发现不符 ⇒ 写在证据件里并报回，**不要按我认为的样子改**。
