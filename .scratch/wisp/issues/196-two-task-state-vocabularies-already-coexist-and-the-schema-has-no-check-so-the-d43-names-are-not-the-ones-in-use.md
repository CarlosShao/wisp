# 196 — 任务状态其实有**两套词表**在库里跑：面板/取证读的是 D43 那 40 行，而写进去的四个词一枚都不在那张表上，schema 又不拦

- Status: **未派（先摆 owner 一次：要不要给 schema 加约束＝可能触 `D35` 数据模型）**。⚠ **本票不是"缺字段"**（那是票 188），本票是**已经并存的二致性缺陷**。
（资料/裁决记录型：本票无判据框，⛔ 不参与"全勾才收口"那把尺；要拍板的事项见正文"要拍的"一节。）
- 来路：`188-r1` 停手上报（台账 `A394`），我逐条复跑对格。
- 现量（锚 `8c671afb`，尺可复制）：
  | 事实 | 读数 | 尺 |
  |---|---|---|
  | D43 名字的仓内逐字拷贝 | `internal/statemachine/states.go:11-31`（`:8` 自陈"names exactly as in D43 / SPEC-08"） | `sed -n '8p;11,31p' internal/statemachine/states.go` |
  | D43 出处（冻结） | `docs/PLAN.md:3055` 标题＋`:3061-3100` | `sed -n '3055p' docs/PLAN.md` |
  | 今天真写进库的状态词 | **`done`／`cancelled`／`running`／`succeeded`——一枚都不在 D43** | `sed -n '983,998p' internal/agent/loop.go`；`sed -n '651,670p' cmd/wisp/run.go` |
  | schema 有没有拦 | **没有 CHECK** | `sed -n '57p' internal/memory/models.go` 邻域（字段定义 `:62`） |
- **后果（可观察形状）**：同一枚任务在"取证列表"与"面板任务监控"里可能一个说 `succeeded`、一个说 `finished`；任何新读面默认自己的那一套；第五个词今天加进去**门不会响**。
- **要拍的（不答＝默认甲）**：甲＝**只定映射表**（新增一枚 Go 侧唯一权威映射，两套词各自保留、读写都走它，不动 schema）；乙＝**再加 CHECK 约束**（数据模型变更，触 `D35` 射程，须人工批准＋迁移路径）；丙＝先只做票 188 那枚新字段、本票按住。
- 禁区：`PLAN.md`／`docs/specs/**` 一字不动；`frontend/**`／`design/**` 零写面；**不许顺手把 `TaskLog.State` 改成 D43 名**（那会洗掉现有生产者写入的四个词，是另一程的事）。
