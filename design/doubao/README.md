# Wisp 高保真交互原型 · demo

> 纯静态 HTML/CSS/JS，双击 `demo/index.html` 即可在浏览器中打开（file:// 可用，无需构建、无需联网）。
> 视觉语言：**Wisp Minimal · 磨砂极简**——Linear/Notion/shadcn 文档风，克制留白，轻 Acrylic/Mica 磨砂。
> 当前版本：v4 全面对齐 BeautifulUI 21 组件 + ReactBits 动画体系（Sidebar 滑翔/Spotlight/磁吸/屏切换）。

## 打开方式

```
design/doubao/demo/index.html
```

直接双击或拖入浏览器。所有依赖（Tailwind play CDN、Lucide 图标）已下载到 `lib/` 本地，离线可用。

## 目录结构

```
demo/
├── index.html          # 入口：桌面背景 + WebView2 窗口 + 悬浮球 + Modal 容器
├── styles.css          # shadcn CSS 变量(token) + 磨砂材质 + 组件工具类
├── app.js              # 路由、导航、Toast、Ctrl+K 命令面板、Modal、主题切换
├── lib/
│   ├── tailwind.js     # Tailwind play CDN（本地）
│   └── lucide.min.js   # Lucide 图标 UMD（本地）
├── screens/            # 每屏一个 JS 文件，模板字符串注册到全局 Screens
│   ├── chat.js         # 对话（流式/工具chip四态/结果分流/附件/模式切换/历史分组）
│   ├── approval.js     # 审批（队列选中联动/L1L2区分/批量拒绝/队列满）
│   ├── palette.js      # 命令面板（5组/搜索过滤/空态/快捷键）
│   ├── tasks.js        # 任务（六态/路径锁/cron+电源条件/新建modal）
│   ├── ball.js         # 球状态（20态交互列表/CSS球联动/转移图/热键）
│   ├── config.js       # 设置（7节/三档生效/校验错误/搜索/未保存提示）
│   ├── security.js     # 安全（权限模式/授权撤销/黑名单/taint时间线/注入拦截）
│   ├── privacy.js      # 隐私（5tab真实字段/保留期/麦克风开关/存储位置）
│   ├── cost.js         # 成本（超限态/日月视图/Top5/柱状图/token分解）
│   └── firstrun.js     # 首次引导（真实模型清单/校验/目录授权/DPAPI）
└── screenshots/        # 逐屏验证截图（v1/ v2/ v3/ v4/）
```

## 交互说明

| 操作 | 效果 |
|---|---|
| 点击左侧导航图标 | 切换屏幕（9 屏） |
| 点击右缘悬浮球 | 显示/隐藏主窗口 |
| Ctrl+K | 唤起交互式命令面板（搜索/上下导航/Enter执行/Esc关闭） |
| 点击标题栏太阳/月亮 | 浅色/深色主题切换（localStorage 持久化） |
| 点击导航底部救生圈 | 重新打开首次引导 Modal |
| 球状态页点击任意态行 | 上方 CSS 预览球随态变色/变尺寸/切动画 |
| 审批页点击队列项 | 右侧详情联动，可展开完整参数 |
| 设置页点击左侧节导航 | 切换配置节，搜索框跨节过滤 |
| 隐私页点击 tab | 切换 L1画像/L2记忆/L3日志/tool_call取证/artifacts |
| 成本页点击日/月视图 | 切换 token 分解表 / Top5 任务榜 |
| 任意无真实功能的按钮 | 弹出 Toast「演示原型，仅作展示」 |

## 屏幕清单与对应工单

| 屏 | 深化后的真实功能 | 对应工单/规格 |
|---|---|---|
| 对话 | 流式回复+停止按钮、工具chip四态（等待/执行中/完成/被拒）可展开参数结果、结果分流标识（面板/TTS/通知）、附件缩略图可删、工作区路径、干活/陪聊/自动模式切换（陪聊出红环）、历史按日期分组、推理折叠 | 92 composer、30 结果路由D10、26 TTS、28 会话预热、36 结果历史、D10/D11/D47 |
| 审批 | 5条队列选中联动详情、每条引用R规则+风险等级、L1轻确认(2-3s倒计时)与L2强确认(300s)视觉区分、完整参数展开、面板仅拒绝+查看无允许键、批量拒绝、队列9/10接近容量上限+满额fail-closed说明 | 21 审批最小集、37 L2卡、48 队列满、49 会话授权D45-2、SPEC-06、D45/D31 |
| 命令 | 5组（快捷指令/工具带L0L1L2 chip/设置/跳转/最近任务）、每项kbd快捷键、搜索过滤+空态（回车提交为自由文本）、最近任务按时间倒序、可点击跳转 | 38 命令面板、SPEC-08 palette、D12 |
| 任务 | 六态色标（运行中/等待审批/路径阻塞/排队/失败/完成）、路径阻塞标注「在等task-6释放Desktop」、cron中文+电源条件（空闲时/解锁时/接通电源）、运行中可取消、新建任务modal（类型/触发/电源/目录）、定时提醒开关 | 47 定时任务路径锁、31 提醒、43 电源事件、D31/D12 |
| 球状态 | 20态完整列表（英文名+中文名+触发条件+视觉动效+点球行为）、CSS预览球随选中态变色/变尺寸/切动画（pulse/breathe/spin/flow/ring/fade）、关键路径转移图（含热路径Warm→Listening零模型加载）、全局热键卡（Ctrl+Alt+Q/M/K/,/Esc）、球体交互卡（单击/拖拽/双击/右键/穿透） | 62 液态玻璃、64 热键交互、68 缺陷、SPEC-08第2节、D43 |
| 设置 | 7节导航（通用/语音/音频/大模型/工具/外观/隐私）、三档生效徽标（即时绿/重载黄/重启灰）、大模型9键+api_key只读、api_base校验错误红框、搜索框跨节过滤、未保存改动提示条、语音含采样率/VAD阈值、工具有L0L1L2默认级 | 39 配置编辑器、SPEC-03全文、D36 |
| 安全 | 三档权限模式（每步都问默认/只问高危/全自动需确认）切换+差异说明、Job Object隔离状态、生效授权（工具/路径通配/L0L1L2/剩余时间/可撤销）、A/B级黑名单网格（System32/.ssh/config.toml/.env*/.pem）、taint污染时间线、D30间接提示注入拦截记录、导出日志 | 17 风险评估C19、90 权限模式、40 安全页、SPEC-06全文、D33/D30 |
| 隐私 | 5tab真实字段：L1画像(slot枚举+人话+来源+时间+上限20条LRU)、L2记忆(全选批量删)、L3任务日志(状态/耗时/费用)、tool_call取证(工具+L级+参数截断+审批决定+含L2被拒rm -rf)、artifacts(文件/大小/路径)、保留期30/90/永久、麦克风开关+状态说明(D16默认关/缓冲不落盘)、存储位置(wisp.db 12.4MB/artifacts 86/500MB) | 29 记忆L1L2、40 隐私页、SPEC-02、D20/D16/D15 |
| 成本 | 今日/本月/累计三卡、本月超预算红色超限态(58.60/50,117%,新任务暂停banner)、预算80%告警线、7日纯CSS柱状图(内联色值+周三高亮+数值标注)、日视图(按模型IN/OUT/缓存分解+任务明细)、月视图(花费Top5任务榜+月度模型分解)、C23微单位对账说明 | 44 成本计量C23、SPEC-05 C23、D15 |
| 首次引导 | 三步进度轨、模型下载含真实清单（asr-streaming-paraformer/asr-offline-sensevoice/punc-ct-transformer/kws-zipformer，各带大小+minisign校验状态）、68%进度环、目录授权可增删、API Key掩码+DPAPI说明、每步可跳过 | 14 模型分发、63 凭据录入、SPEC-03、D26 |

## 视觉规格（Wisp Minimal）

- **材质**：窗口表面 rgba(255,255,255,0.72)（dark: rgba(22,24,29,0.72)），backdrop-filter: blur(28px) saturate(1.25)，1px 发丝边框，柔和阴影 0 12px 40px rgba(0,0,0,.10)
- **配色**：浅色默认；品牌青雾 #86C2B9 仅用于 primary 按钮/选中态/focus ring；其余中性灰阶
- **字体**：Inter 到 PingFang SC 到 Microsoft YaHei UI；页标题 20/600，分区标题 12/600 大写淡色，正文 14/400，辅助 12 muted
- **组件**：全部 shadcn 化——Button h-9 rounded-lg、Input h-9、Card rounded-xl、Badge h-6、Tabs 下划线、Switch、Progress
- **圆角**：6/8/10/12；间距 4 的倍数；区块间 24/32
- **图标**：仅 Lucide，stroke-width 1.75；零 emoji（CI 硬约束，含 U+2190-2BFF 箭头符号区间）
- **动效**：150ms ease；hover 浅底；点击 scale(.98)；屏切换 120ms 淡入


## BeautifulUI 组件基座（仿真移植）

agent 工作流部分以 BeautifulUI ai-native 组件为基座（工单 77 AC#7，vendored 在 `frontend/src/components/ai-native/`）。本 demo 用纯 CSS/JS 仿真其视觉语法，不引 React：

| BeautifulUI 组件 | demo 落地位置 | 仿真要点 |
|---|---|---|
| `thinking.tsx` | chat.js 推理块 | thinking-dots 三点呼吸 + shimmer-text「思考中」标签 + 可折叠推理步骤（animated-list-item 120ms stagger） |
| `streaming-text.tsx` | chat.js Agent 回复 | stream-word 逐词 blur 渐显（45ms/词）+ stream-cursor 光标 + 重播按钮 |
| `tool-chips.tsx` | chat.js 工具调用行 | icon+label+mono chip + tool-detail 左边框展开详情 + 四态（spinner-ring/shimmer-text/check/x） |
| `approval-card.tsx` + `l2-approval-card.tsx` | approval.js L2 卡 | LevelBadge 风险等级 + RuleList outline badge + fail-closed 原因文案 + mono 参数 pre 块 + 调用链 |
| `task-rows.tsx` | tasks.js 任务行 | SVG spinner-ring（dasharray 0.28/0.72）+ 80ms stagger 入场 + 状态转换 |
| `shimmer.tsx` + `loading-state.tsx` | 全局加载骨架 | shimmer-line/shimmer-block 渐变扫动 + 各屏切换/首载过渡 |
| `result-stream.tsx` | chat.js 消息流 | 消息分组 + 工具 chip 行 + 结果分流徽章 + animated-list-item 入场 |


## BeautifulUI 21 组件全集映射（v4）

| # | BeautifulUI 组件 | demo 落地屏 | 实现要点 |
|---|---|---|---|
| 1 | Loading State | 全局 shimmer 骨架 | 各屏切换/首载 shimmer-line/block 渐变扫动 |
| 2 | Thinking | chat | thinking-dots + shimmer-text + 可展开 traces（Steps/Reasoning/Search/Coding 子页签） |
| 3 | Streaming Text | chat | stream-word 逐词 blur 渐显（45ms/词）+ cursor + 重播 + 内联来源 + Follow-ups |
| 4 | Approval Card | approval | 头部 LevelBadge + pagination-dots + 原因/RuleList/参数分区 + 倒计时 + 无允许按钮 |
| 5 | Tool Chips | chat | icon+label+mono chip + tool-detail 左边框展开 + 四态（等待/执行中/完成/被拒） |
| 6 | Task Rows | tasks | SVG spinner-ring + 子任务进度条 + 六态色标 + 80ms stagger 入场 |
| 7 | Chat | chat | 会话页签 + 用户气泡/Agent 无气泡 + reasoning 回复 + 结果分流 badge |
| 8 | Prompt Bar | chat | .prompt-bar Rounded 变体：@来源 //命令/textarea/模型下拉/mic/发送 |
| 9 | Recommendation Card | approval | 置信度条 confidence-bar 78% + 备选方案 + 接受建议 |
| 10 | Context Cards | privacy | .context-card + .ctx-source 知识块（L1画像/L2记忆，2列网格+card-spotlight） |
| 11 | Diff Table | （预留，tool-detail 中用 add/del 色） | chat 工具详情 diff-add 行 |
| 12 | Records Table | security/privacy/cost/tasks | .records-table 授权/取证/日志/明细/定时任务，th 可点排序 |
| 13 | Filter Table | tasks | .filter-chips 状态过滤（全部/运行中/等待/阻塞/排队/完成/失败）实时重组 |
| 14 | Sidebar Nav | 全局 app.js | .nav-glide 滑翔指示条（280ms cubic-bezier 平滑滑动）+ 底部引导位 |
| 15 | Search | palette + Ctrl+K | 大字号搜索框 + 五组实时过滤 + 空态 + kbd 快捷键 + 底部键位提示 |
| 16 | Flowchart | security | .flow-canvas 虚线点格画布 + trigger/gate/tool .flow-node + arrow-right |
| 17 | Insight Cards | cost | .insight-card + .insight-num 大数字 + count-up + pagination-dots 洞察分页 |
| 18 | Code Block | config | .code-block 行号 + tok-key/str/num/com 语法色 + config.toml 预览随节联动 |
| 19 | Fine-tune Card | config | inspector 式键值控件（mono 键名+描述+input/select/switch/slider+生效徽标） |
| 20 | Selection Actions | （预留，chat 划句交互） | Prompt Bar 选中文本时的 @ 来源机制 |
| 21 | Agent Screen | chat 整屏 | 会话页签+消息流+thinking+streaming+tool chips+prompt bar 完整 agent 界面 |

## ReactBits 动画清单（R19「留5」）

工单 77 第 225-239 行 R19 表定案的 5 个动画，全部在本 demo 仿真落地：

| 动画 | 落地位置 | 参数 |
|---|---|---|
| Thinking Dots | chat 推理块、firstrun 下载状态、全局 Preloader | 三点呼吸 1.2s，delay 0/0.15/0.3s |
| Staggered Text | chat Agent 回复正文 | 逐词 45ms，buStreamIn 420ms cubic-bezier(0.22,0.61,0.25,1) |
| Animated List | 全部屏幕的列表/消息/任务行 | buFadeUp 320ms，stagger 40-80ms/行，maxItems 兼渲染上限 |
| Blur Highlight | chat 工具 chip 展开 | 选中项 bh-focused 清晰，其余 blur(1.5px)+opacity 0.45，250ms 过渡 |
| Preloader | 全局首载 + firstrun 模型下载步 | 全屏覆盖 1.1s 后淡出，品牌球呼吸 + 「一缕 · 正在唤醒」 |
| Screen Transition | 全局屏切换 | .screen-enter 淡入+上滑 220ms cubic-bezier(0.22,1,0.36,1) |
| Card Spotlight | 全局卡片悬停 | .card-spotlight 径向高光跟随鼠标（--mx/--my），200ms 淡入 |
| Button Micro-magnetic | 全局按钮 | .btn-magnetic 极轻微 transform，active scale(0.97) |
| Sidebar Glide | 左侧导航 | .nav-glide 指示条 280ms 平滑滑动到 active 项 |
| Count-up | cost 三卡 | requestAnimationFrame 600ms easeOutCubic 0->目标值 |
| Tooltip | 全局 | .tooltip 弹入 150ms scale+opacity |

**氛围层 Fog Sphere**：桌面背景两个极淡呼吸色块（青雾 420px + 蓝紫 360px，blur 80px，12s 呼吸），放在窗口背后不抢内容。

### 许可口径与生产门

- ReactBits = MIT + Commons Clause，**生产引码属二期、需 owner 复核**；本 demo 仅做效果仿真，不引其代码。
- Fog Sphere 生产门（SPEC 约束）：同屏最多 1 个、面板隐藏即销毁、CPU 空闲回落；demo 里做了 2 个极淡版并在 prefers-reduced-motion 下降级为静态低透明度。
- 全部动画遵守 `@media (prefers-reduced-motion: reduce)` 降级（时长 0.01ms、迭代 1 次）。

## 与真实前端的对应关系

本 demo 的 token 体系与组件形态可直接移植到工单 77 定案的 React + shadcn 技术栈：

| demo 实现 | React 前端对应 |
|---|---|
| styles.css CSS 变量 | tokens.css（C21 契约，四方对账） |
| screens/*.js 模板字符串 | React 组件（result-stream、l2-approval-card、composer 等） |
| app.js 路由/状态 | Go 侧 PanelSnapshot 推送 + 前端受控渲染 |
| Tailwind play CDN | Tailwind v4 + @tailwindcss/vite |
| Lucide UMD | lucide-react 1.47 |

## 验证记录（v4 全面对齐版）

- v4 全局：Sidebar Nav 滑翔指示条平滑滑动、屏切换淡入上滑、card-spotlight 悬停高光、窗口分层投影
- 10 屏全部在 Chrome 中逐屏渲染验证，布局无溢出
- 全局 Preloader 首载 1.1s 后正常淡出，切屏不闪
- Fog Sphere 桌面氛围正常呼吸，深色/浅色主题下均可见且不抢内容
- chat 屏：流式逐词播放正常（45ms/词+光标）、重播按钮可重放、Thinking Dots 呼吸、工具 chip Blur Highlight 焦点切换正常、工具 chip 四态视觉正确
- approval 屏：L2 卡按 BeautifulUI 族重画（LevelBadge/RuleList/fail-closed/mono 参数）、队列选中联动、无允许按钮
- tasks 屏：spinner-ring SVG 加载环旋转、shimmer 骨架 0.8s 后落真实内容、六态色标正确
- firstrun 屏：进度环 0%->68% 模拟推进、模型行 shimmer->check 错峰落定、「在装模型/在起引擎」状态标签
- config/security/privacy/cost 屏：节/tab/视图切换均有 shimmer 过渡，成本三卡数字 0->目标值滚动
- 浅色/深色主题切换正常；Ctrl+K 命令面板正常；Toast 正常
- 控制台零报错（仅 Tailwind CDN 生产提示 + 2 个无效 Lucide 图标名已修复为 quote/audio-lines）
- 全 10 文件 emoji/禁扫符号扫描：0 匹配
- 全部 10 文件 + app.js node --check 语法通过
- 成本页柱状图用内联 style（height 像素 + background-color 色值）

## 归档说明

- `01-ball-states.jpg`：用户满意的球状态定稿图，保留在本目录（ball.js 引用）
- `design/old/rejected-image-drafts/`：前一版 8 张面板生图（02–09），因视觉方向不符已作废归档（只移不删）
- `design/old/`：旧 HTML 原型（assets/screens/index.html），参考保留
