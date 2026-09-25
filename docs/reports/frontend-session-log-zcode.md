# 前端会话流水账 · ZCode 接管（2026-09-25 起）

> 本文件是前端角色的 append-only 流水账。前一程（Qoder 侧）的 §1…§62 在
> `docs/reports/frontend-session-log.md`，按交接文档要求新会话另起文件，不改旧节。
> 写出者：ZCode（owner 2026-09-25 当面移交「前端这块完全交给你重构了」）。

## §1 接管与三项 owner 裁决（2026-09-25 晚）

**接管**。两枚子代理独立核实过（票池审计 + 组件库摸底），与交接文档 `frontend-handover-to-zcode.md`
无矛盾，但三处口径修正：145 整体是 Go 侧票（前端只欠乙段 L2 卡文案一格）；92 已 accepted-done；
34 是被 77 覆盖却无人关的僵尸票（已列入报回编排者清单）。

**owner 三项裁决（对话原文意，逐条可撤）**：

1. **demo 与 PLAN 冻结规范撞车时以 demo 为准**。射程：§17.4 图标扩充 7 枚
   （message-square-text / shield-check / list-checks / orbit / settings / chart-column / life-buoy，
   render-nav.tsx 的 APPROVED_EXTENSION 常量是批准记录）＋ §17.3.2 中文 14px 下限让位 demo 的
   11-13px。**不包含** sparkles/brain 等 §17.4 点名的"AI 俗套"禁图标——继续避开；**不包含**
   十四态表（:3475-3488）的动效与"明确不用"列——那张表仍是硬约束（SSE 光标 1px×14px --accent
   呼吸、L2 卡无允许键、思考中不用三点跳动等全部保留）。
   撤销口令：「撤 demo 图标裁决」⇒ panel-views 回到 Q1=甲 替代图标 + render-nav 回旧判据。
2. **分两批交付**。第一批 = token 第三代 + 窗口骨架 + 对话屏 + 审批屏；第二批 = 其余六屏 + firstrun。
3. **R19 动画走 react-bits 免费档近亲 + demo 的 CSS 仿真兜底**（Pro-only 无源码的 10 枚不上 Pro 渠道）。

**token 第三代（owner 授权"推翻重新按照你的来"）**：`scripts/gen-tokens.mjs` 重写。
真相源不变（demo styles.css + index.html），证明机制不变（逐行引用 + --check + 悬空 var 拒绝），
表形状改为**直接讲 demo 自己的名字**：层 1 = demo 26 亮 + 25 暗枚逐字（HSL 三元组原样保留）；
层 2 = 41 枚 demo 内联画过的值（success 绿/warn 琥珀/tok 语法色/fog/ball/三层投影/--panel-alpha 旋钮）；
层 3 = 30 枚尺度。0 interim（每枚都有出处）。**一处刻意改名**：demo 的 `--accent`（浅色调）
发射为 `--accent-soft`，`--accent` 留给品牌青——因为 PLAN.md:3476 冻结行点名光标色 token 是
`--accent`（指品牌青），render-stream.tsx 按 `var(--accent)` 断言。旧第二代语义层
（--bg-base/--fg-primary 等 79 键）退役。

**已知恒红（非前端债，报回编排者）**：`TestC21DesignTokensFourWayAgree`——两枚红因
（design/assets/tokens.css 被 owner 挪走 + 表形状已换）在第三代下同样存在，需 Go 侧同步
internal/ball/tokens.go 与 docs/evidence/s1/c21-native-tokens.md 才能绿。前端不擅自凑绿。

## §2 第一批施工读数（续）
