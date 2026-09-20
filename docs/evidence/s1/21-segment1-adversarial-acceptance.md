# 票 21 段 1 对抗验收（决策层）— 非实现者执行：编排者

**时间：** 2026-09-20 21:31–21:33 · **HEAD：** `5c6e661` · **被验范围：** 仅段 1（L1 窗 + 否决通道 + L2 队列 + 来源校验 + 批量 + applied-steps 类型）
**方法：** 只读审计代理先出判据表，**其两条 BLOCKER 与一条 MAJOR 我逐条用自己的 grep/读码复现后才登记**（不复现不采信，含"好消息"）。
审计代理还跑了四个包：`approval ok 0.094s`、`tools ok 10.9s`、`risk ok 3.9s`、`cmd/wisp ok 24.6s`。

## 裁决表（与票面段 1 的 5 框 1:1）

| # | AC（票面摘要） | 勾选 | 裁决 | 证据类别 | 判据 |
|---|---|---|---|---|---|
| 1 | L1 到时执行；每个有效否决通道取消；KWS 未加载时提示「语音取消不可用」；否决失败转 LLM | `[ ]` | **正确未勾 → PARTIAL** | machine | `window_test.go:20/:55/:84(句子逐字)/:144/:163/:264` + 真桥 `wiring_test.go:112/:151`。**但否决通道在生产里没有输入设备**（A19） |
| 2 | L2 展示全参 + `rules_hit`；原生放行执行；伪造面板放行被拒；300s 拒 + 270s 预警 + 重放 | `[ ]` | **正确未勾 → PARTIAL** | machine（数据/路由）+ live（卡片、人的放行） | `queue_test.go:22/:55/:59/:89/:147/:185/:264/:313`；`gate.go:602-615`。伪造轴我复看：授权靠 **route + 活 nonce + bind 摘要**（`approval.go:253-261` 摘要、`:292-306` `subtle.ConstantTimeCompare`），`Request.Source` 只进日志（`gate.go:591/:604`）不进判定 ⇒ **结构上无法伪造放行**，这条我认可 |
| 3 | 10 个 L1 操作 → 一次确认；L2 永不聚合；单次 ≥50 → L2（R7） | `[ ]` | **正确未勾，且实现站错了轴**（A20） | machine，但**测试自建的工具形状生产里产不出来** | `batch_test.go:19/:65/:87/:119/:140`、`batch.go:35-41`、`gate.go:228-244` |
| 4 | 取消中途执行 → applied-steps 报告**被渲染** | `[ ]` | **正确未勾 → PARTIAL** | machine（类型/措辞/bus）+ live（渲染属段 2） | `report.go:94-153`、`window_test.go:207`、`wiring_test.go:182/:247`、`fs_write_test.go:790` |
| 5 | 状态机 `Confirming`/`AwaitingApproval` = D43 #17,21–24 | `[ ]` | **正确未勾（本票段 2 的活）** | **纯 self-report** | 本包**零仪器**：`approval` 不 import `statemachine`（只有 `doc.go:11` 一句注释），态在 `internal/statemachine/*` 与 `internal/ball/anim.go:46` |

**self-report 清点（用户最在意的一类）：AC#5 整条；AC#1/#2/#4 各自的"真机可见"半。**
**没有一框是靠 Progress log 自述支撑其"决策层"半的**，五框全部**当前不应勾**，本验收**不新增任何勾选**。

## BLOCKER（复现后成立，各归其主）

- **BLOCKER-1（=registry A19）审批层有判定、没有输入设备。** 我自己的 grep：
  `grep -rn "DecideFromNative|\.Native\(\)|\.Veto\(|DecideFromPanel" --include=*.go cmd/ internal/ | grep -v _test`
  ⇒ 命中**只有定义本身与 `doc.go:32` 一句注释**，**零生产调用者**。
  另一半同形：`cmd/wisp/run.go:257` 用的是 `approval.NewChannels()`（**无参 ⇒ 零通道已加载**），
  而 `approval.go:163` 那个"载入 Ball+Esc"的构造函数**没人调**，`SetLoaded`（`:178`）**零非测试调用者**。
  ⇒ **今天真机的实际行为：L1 窗一定在无否决的情况下到时执行；L2 一定在 300s 自动拒绝。**
  门控"看起来在"，但**没有任何一条路能让人回答它**。闭合判据见 A19。归属：**票 21 段 2**。
- **BLOCKER-2（=registry A20）D45-1 聚合站错了轴，而且即便站对了也永不触发。** 我读 `batch.go:35-38`：
  `func Aggregate(d tools.Decision)` 判的是 `len(d.Paths) >= 3` —— **单次调用内的路径条数**。
  而契约原文（`docs/PLAN.md:1590` / `:1795`）写的是「**批量聚合（500 个 L1 → 一次确认）**」
  「批量场景 = 500 次 L1 确认」——**跨多次调用**。两个轴不是一回事。
  再核可达性：`fs.write` 1 个 `path`（`fs_write.go:223-225`，`additionalProperties:false`）、
  `fs.trash` 1、`fs.delete` 1、`fs.move` **2**（`from`/`to`，`:415-417`）
  ⇒ **任何在产工具单次最多贡献 2 条路径，`Aggregate` 对全部生产输入恒返回 nil**（死码）。
  闭合判据见 A20。归属：**票 21 段 2 只许登记、不许顺手改**（D45-1 按 PLAN `:3105` 落 **S3** 本片）。

## MAJOR / MINOR

- **MAJOR（A13 的第③条判据未满足）**：`cmd/wisp/run_test.go:272-276` 那条 L2 用例被 **D47 守卫在执行前拒掉**
  （`:317` 断言 L2/reject）⇒ **没有任何 cmd 级用例真的把一张 L2 卡显示出来**。
  闭合：让该用例走到 admit，再断言 `rt.ui.shown()==1`。
- **MINOR-1**：`corr` 是调用方可得的**条目地址**（`bridge.go:219` ← `loop.go:640` 用 taskID），
  `revokeGrants`/`reject`/`Veto` 都直接吃它 ⇒ 猜中 corr 的**面板路由**调用方可以作废诚实卡片的活 nonce。
  **只是 fail-closed 的 DoS，不是提权**（伪造放行仍被 bind 摘要挡住）。闭合：断言 corr 对模型不可得，或按条目哈希。
- **MINOR-2**：`bridge.go:501-503` `origin = dec.Paths[0]` —— 污点归属里的"取第一个"选择器；
  今天**确定**只因 `pathArgs` 先排序。**同族第四次**（M-7 / C-3 / A16），登记为形状而非当期缺陷。

## 判据完整性（我最担心的那类，结论是好的）

逐文件 `git log --oneline --`：**没有**任何测试文件晚于它所管的代码被改。两批"码+测同 commit"都是**加强**：
- `b1535d7`：`DecideFromPanel` 从 `allow(corr, r.Grant+"!burned")` 改成 `revokeGrants(corr)`（`queue.go:256-262` 新增）
  —— 旧码什么都没烧掉，新行为由 `queue_test.go:129` 钉住；
- `e669567`：`gate.go:447-483` 加 `markStarted(incoming)`（L2 放行才记已开工）+ `report.go` 无否决分支，
  修的是**真·假干净取消**缺陷，由 `wiring_test.go:247` 钉住。
⚠ 唯一**没有外部仪器**的一框：AC#3 只在 `approval` 包内自测（`batch_test.go`），无桥级/cmd 级独立用例。

## 必须等桌面空出来才能做的（段 2 排队清单）

真球点击否决；真全局 Esc 否决；原生卡片**可见地**渲染全参 + `rules_hit` + 拒绝钮；
球的 `Confirming` 脉冲与 `AwaitingApproval` 深度徽标；物理条上的「语音取消不可用」；
**一次真人放行让 L2 结束**；applied-steps 报告**渲染**给人看。

## 总体判断

段 1 的**判定内核是我在本仓见过最干净的之一**（route+nonce+digest 三轴，来源字段不进判定，两处同 commit 都是加固），
但它今天**是一台没有接线柱的开关**：两条 BLOCKER 都不是"写错了"，而是**"没人把它接到人和机器之间"**——
与 registry A8/A11/A12/A13/A17 完全同一形状，且这次连"能力面"都缺（`Aggregate` 永不触发）。
**⇒ 段 2 的验收重点不是美观，是"让一个真人能改变一个判定"。**
