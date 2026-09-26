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

## §2 第一批施工读数（2026-09-25 晚，commit ab7afca + 62431e9）

- 门禁现量：typecheck rc=0｜tokens:check rc=0（191 键）｜render:nav OK（70 冻结 + 7 批准图标）｜
  render:stream OK（1×14 --accent 1→0.2→1 1s）｜render:composer OK（三态）｜
  render:l2 OK（真夹具 l2-card-fs-delete.json）｜vite build rc=0（gz 94.5KB）｜
  oxlint 0 errors｜d22scan rc=0｜go test ./internal/panel/ = 58 PASS / 1 FAIL
  （TestC21DesignTokensFourWayAgree，已知 Go 侧债）。
- 踩坑两枚（都在前端文件注释里）：①theme.css 注释写了 `--color-*/--radius-*`，
  `*/` 把注释提前闭合，lightningcss 把后面英文当 CSS 解析报「Unterminated string」；
  ②panel-views/render-nav 注释里写了 `∪`（U+222A，在 U+2200-22FF 扫描段）——
  前端文件的注释不豁免 emoji 尺（Go 侧那把尺和 d22scan CLI 的注释豁免口径不同）。
- render-l2 的 needle「点击悬浮球以批准」要求每卡恰好一次：info 行不能复读这句
  （球图示 footer 已有一次），F2 字样已从 info 行撤（原生未接线，A222#7）。
- lucide-react 1.47 命名换代：AlertTriangle→TriangleAlert、AlertCircle→CircleAlert、
  Trash2 不存在（用 Trash）。render:nav 的 PLAN 解析器正控反控都没动。
- 浏览器实检（vite preview 4173）：harness 七屏结构/滑翔条/L2 卡/审批队列/设置滑杆
  全部渲染正确；生产空态 composer 钉底（62431e9）；生产磨砂透浏览器画布底色发灰
  属 --panel-alpha 设计本意，真实 backdrop 归票 33。
- 在挂 vendored 件：Shimmer（骨架切换过渡）＋ StatusPill（标题栏待审批）＋
  Chip（设置屏键名）＝3 枚，挂载尺安全。
- 第二批待做：palette/tasks/ball/config 全节/security/privacy/cost/firstrun 八屏
  ＋L1 阻止窗口条（等 owner 答"谁画"）＋R19 动画 react-bits 免费近亲 vendoring。

## §3 路线转向：demo 废弃，beautiful-ui 主题收编（2026-09-25 晚，owner 第四次裁决）

owner 看过第一批 harness 后原话（意）：「推翻再重来，不按照 demo 来了…按照你的想法
重新以 beautiful ui 组件库…重新设计一下 token…必须是极简风的，就两种风格，暗色/浅色…
所有动画用 reactbits，但是这些动画是在 beautiful ui 没有的前提下再用…越看越觉得这个
harness 有点四不像，起码要有主流 harness 的样子」。

**落地方案**：
1. **token 第四代**：收编 turbo-kach（beautiful-ui 官网那套）自己的 app/globals.css 主题
   逐字为主题——极简纯面、明暗两态、蓝 accent #0285ff、hairline 系投影、6/8/10 圆角。
   demo 主题（青雾磨砂/雾团/环境光）全部废弃；demo 不再是任何意义上的样式来源。
   生成器真身改为「本脚本即设计定案」，来源与上游 commit（05dab2d2）在头注声明。
2. **组件**：turbo 19 构件全部零 props 演示页（grep 证实：TaskRows({variant}) 一族），
   但纯 MIT 允许改写 → **改写成 props 驱动版**（视觉动画照抄蓝本、数据进 props、
   无 props 不渲染），头部 derived-from 登记；slev 11 atoms 全 props 驱动 → 逐字
   vendored（Button/EntityChip/Switch 各 1 处字面量做换 token 微改并登记）。
3. **动画**：beautiful-ui 自带（fade-up/pop-in/stream-in/pixel-on/shimmer-text/spin/
   eq-bounce）先用；react-bits 只补库没有的（缺口清单后补）。
4. **harness 重定义**：?harness=1 渲染 **showcase**（官网式组件展示页：header+章节索引+
   编号章节+框式 demo 卡+明暗切换），main.tsx 已接；生产面板（App）永不吃假数据。
   面板自身的浮窗/磨砂/雾球 chrome 全部删除。
5. **不动的**：render-* 四把尺全部 needle（点击悬浮球以批准/拒绝/查看完整参数/
   l2-ball-ring/bg-stop/档位字符串族/1.0 MB 格式化等）、panel.ts 契约、panel-views、
   App 安全设计（L2 卡全局渲染）、六把尺、git 纪律。
6. 转向口令：「撤 beautiful-ui 转向」⇒ 回 demo 路线（token 第三代 + 第一批提交）。

## §4 react-bits 动画缺口审计（2026-09-26，代理 D 产出）

owner 规矩「动画用 react-bits，但 beautiful-ui 没有的才用」。审计结论：**本轮零 vendored**——
5 枚点名候选全部有依赖：FadeContent→gsap，CountUp/AnimatedList/BlurText/ShinyText→motion/react。
gsap 与 motion 都是「未经 owner 批准不得新增」的 npm 依赖 ⇒ 全部判不可用；ShinyText 另有
props 默认 #b5b5b5/#ffffff 字面量，即使依赖解决也过不了零字面量尺。

R19（票 77:235-249 留5/缓4/砍2）逐枚：
- 已覆盖 7：Thinking Dots（pixel-on 点阵+shimmer）、Staggered Text（stream-in/reveal-text）、
  Animated List（task-rows 80ms fade-up stagger）、Screen Transition（.screen-enter）、
  Sidebar Glide（.nav-glide）、Tooltip（.nav-rail-label+.kbd）、eq-bounce（已收编待票 35 挂载）。
- 缺但现状够用：Blur Highlight（tool-chips 展开即聚焦；要字面版是 3 行 CSS，无需库）。
- 缺被依赖挡：Count-up（react-bits CountUp 只吃 number，SHOWCASE_COST 三值是字符串，
  money/budget 不适配；tokens 一枚可 to=12480——等依赖决定）。
- 缓 4（性能门后）：Glass Flow/Aura Blob/Neural Float/Fog Sphere——react-bits 近亲全依赖
  ogl/three，同须 owner 批；R19 原文的三条性能用例（同屏至多1/隐藏即销毁/CPU 回落）先有才做。
- 砍 2 维持：Glass Cursor 不复活；Agentic Ball 不替代 Direct2D 球（ball-screen 已是 D43 状态表）。
- Preloader：极简转向已弃，不实现。

**范围外发现（待人拍板）**：react-bits 的 Magnet（按钮微磁吸）是全库唯一零依赖、
零字面量的组件——今天就能照抄。不在 R19 清单，要不要抬由 owner 一句话。

给 owner 的问题：①要不要为 react-bits 引入 motion 或 gsap（引入后 CountUp/BlurText 等
5 枚可落，但 bundle 与 D32 CPU 预算要先评）？②Magnet 抬不抬？③背景层四枚继续缓？
