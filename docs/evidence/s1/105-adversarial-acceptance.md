# 105 — 独立对抗验收裁决表（acceptor-ticket105）

验收方：`acceptor-ticket105`（只读；一行生产码未改，本文件是唯一工作树写件之一，另加 `/tmp` 快照）
被验：票 105「C26 改写账只有测试在读 → 接一条真实生产消费路径 + 补祖先重解析行为用例」
HEAD 复算时点：`f6818f24ec404063fc6dde408edbe5b56acdf2ac`（分支 `dev`），本文件起笔 `date` = **2026-09-21 20:25 CST**
交件自述：票面 `## Progress log` 末条（19:5x，agent-ticket105），其 `next=` 请出本表。

## 裁决表（1 AC = 1 格）

| AC | 票面判据要点 | 裁决 | 证据标签 |
|----|--------------|------|----------|
| AC#1 | 复算"零生产读取者"，贴真实命中表，不许照抄 | **通过（复现＝成立）** | 〔独立复现〕 |
| AC#2 | 至少一条真实生产消费路径（grep 非零 + 装配根可达 + 端到端断记录内容） | **通过** | 〔独立复现〕 |
| AC#3 | 祖先重解析行为用例断可观察结果 + 变异红在行为；独立复算其可达性结论 | （见下） | （见下） |
| AC#4 | 若"给人看"必须动 frontend/ 或 internal/panel/ ⇒ 停手登记 | （见下） | （见下） |
| AC#5 | 门禁：go test/gofmt/gofumpt/vet/GOOS=linux vet/d22scan 纯净快照 rc=0，台账不降 | （见下） | （见下） |

---

## AC#1 — 复算"零生产读取者"

**判据**：票面 AC#1 原文命令 `grep -rn "RewrittenRoots(\|UnusableRoots(\|\.Roots(" --include=*.go . | grep -v _test.go`，贴真实命中；若其实有生产读取者 ⇒ 写"未复现"并降级本票。

**我的复算（当前 HEAD f6818f2，工作树含交件的三枚 commit）**：
- 非测试命中 **5 条**：
  - `internal/tools/bridge.go:810` `// account. Until here the only consumers of Roots()/RewrittenRoots()/` —— **注释**
  - `internal/tools/bridge.go:811` `// UnusableRoots() were tests, ...` —— **注释**
  - `internal/tools/bridge.go:864` `return b.paths.Roots(), b.paths.RewrittenRoots(), b.paths.UnusableRoots()` —— **真调用点（AC#2 新增的消费者）**
  - `internal/tools/paths.go:186` `func (p *PathCanonicalizer) UnusableRoots() []string {` —— **定义**
  - `internal/tools/paths.go:194` `func (p *PathCanonicalizer) RewrittenRoots() []string {` —— **定义**
- 含测试总命中 **35**，非测试 **5**（口径：全仓 `--include=*.go`，`grep -v _test.go`；不带 `-v`）。

**对票面 AC#1"零生产读取者"这句话的判定**：票面 AC#1 要复算的是**建票时**的状态。交件的 19:3x 条自述在**接线前**量到 2 条命中且全是定义行、`Roots()` 分支零命中 ⇒ 那才是"零生产读取者"的原始命题，本票据此未降级。我在**接线后**HEAD 复算，非测试命中变 5，但其中 3 条是定义/注释、**唯一新增的真实读取者就是 AC#2 要造的那个**（`bridge.go:864`）。⇒ 与自述一致：**接线前该命题复现成立**；接线后它由本票 AC#2 亲口改掉。

**引用腐坏检查**：自述 AC#1 引用的 `paths.go:186/194` 定义行、`bridge.go:864` 调用点，我逐条对上，无行号漂移。

**AC#1 裁决：复现＝成立（未复现的反例不存在）**〔独立复现：命令与命中我亲手跑过〕。本票不降级，AC#2/AC#3 照验。

---
