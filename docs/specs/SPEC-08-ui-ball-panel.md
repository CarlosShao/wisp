# SPEC-08 · UI：悬浮球、状态机与面板

> 追溯：D2/D29/§14.6/§17、C12/C17/C21/C27/D43；设计 token 参考实现 = [../../design/assets/tokens.css](../../design/assets/tokens.css)。
> UI 质量是与 SLO 同级的验收维度（D29 自我纠错），视觉验收必须人工（用户本人）。

## 1. 背景与问题

- 悬浮球常驻且必须**原生**（Direct2D/DirectWrite），不能是 WebView（D2）；「好看」不可机器判定，
  但 20 态的颜色/透明度/动效节奏必须可执行地定死。
- 面板按需拉起（C27 单窗口复用）；确认交互按视觉要求分层而非轻重分层（D29）。
- 用户硬要求：**零 emoji、内联 SVG、简约但高级**（毛玻璃/磨砂，v2.3 冷调）。

## 2. 悬浮球原生实现（§14.6）

- **窗口样式**：`WS_EX_LAYERED | WS_EX_NOACTIVATE | WS_EX_TOOLWINDOW | WS_EX_TOPMOST`；
  **永不抢焦点**（否则每次唤起打断正在打字）。
- **绘制**：`UpdateLayeredWindow` + 每帧 32bpp premultiplied ARGB；图形 Direct2D、中文文本
  DirectWrite（GDI 中文渲染质量不可接受）；与面板共用 `ui-sta` STA 线程（D38a）。
- **动画纪律**：`Sleeping` 态完全静态（不启动定时器，CPU≈0）；仅
  `Listening/Thinking/Acting/Speaking/Confirming` 启动动画定时器，帧率上限 **30fps**，
  CPU 计入工作态预算。**所有循环动画只允许出现在非空闲态**（`Warm` 静态或 ≥2s 周期 opacity
  呼吸走合成器；面板关闭后不得有任何动画在跑）。
- **点击穿透**：透明像素区域 `WS_EX_TRANSPARENT` 或 hit-test 返回 `HTTRANSPARENT`；球体本体可点。
- **多显示器**：Per-Monitor V2 DPI 感知；位置按显示器分别保存；拔插后落到不可见区域 → 自动回
  主屏可见位置；DPI 变化按新 DPI 重建 D2D 设备相关资源（否则模糊）。
- **视觉一致性**：所有颜色/圆角/阴影/缓动**必须取自 C21 DesignTokens**（`design/assets/tokens.css`
  为唯一样式真相源），绘制代码硬编码色值 = 违规。
- **壁纸自适应**：球叠在任意壁纸上 → 跟随系统明暗主题 + 边缘描边/玻璃环保证浅色壁纸上可辨（D29）。
- **尺寸**：配置范围 44–72px（默认 56）；`Sleeping` 态为**静态玻璃体**，
  静止体径 = `配置尺寸 × SleepRestRatio(0.62)`（默认 56 ⇒ **34.72px**，下限 `SleepingRestMinPx=30`）
  ——**不是**"直径 44px"（见下面的更正记录）。INTERIM 状态不变，见下。
  - ⚠ **更正（2026-09-20 23:11，编排者本人的录入错误，非契约意图变更，待 owner 过目）**：
    我在本条 INTERIM 变更里第一次写的是「静态玻璃体 **44px**」。**那个数字是错的**：
    44/46 是**差分化像框**的边长（含 halo 与边缘过渡，跨三次实测为 40×38 / 44×44 / 46×46，
    **框本身浮动 ±6px；稳定的量是 `px_delta_ge8`**），而真实绘制的体径由
    `internal/ball/statevisual.go::stateSize` 算出 = `56 × 0.62 = 34.72px`。
    **否证方式是定量的**：以"最宽可见物是 `1.5R` 的 halo 填充椭圆"建模，先验证非 `Sleeping` 行
    （`SizePx=56` 预测 4844 像素，`docs/SLO.md` §A.2 的 Thinking/Acting 实测**正是 4844**），
    再回判 `Sleeping`：`34.72 ⇒ 2130`（实测 2098/2103/2120，差 0.5–1.5%）、
    `44 ⇒ 3421`（**比实测高 62%，且 72px 的窗裁不掉它**）。
    `44` 最可能的来源是把 `BallSizeSmallPx`（**最小可配置尺寸**，`tokens.go:369-371`）误读成体径。
    证据与推导：`docs/evidence/s1/68-*` + `internal/ball/tokens_test.go::TestSleepingSizeTruthTable`、
    `TestRecordedSleepingDiffBoxIsNotA44pxBody`。registry **A29**。
    ⚠ 本更正**只覆盖我今天引入的那句**；同节 §2.1 表内 `Armed` 行的「直径 44px」是**既有契约文本**，
    它与代码（Armed 用配置尺寸 56）**也不一致**，但那是**待裁定的契约-代码分歧**，我不自行动它。
  - ⚠ **INTERIM 变更（2026-09-20，owner 授权继续推进核心功能）**：原规格「`Sleeping` 缩到 12px 微点
    （opacity 0.35）」已被**实测否证**——差分取证显示该规格在浅色桌布上"进程在/进程不在"**零个像素**
    变化 ≥8/255（最大 4/255），即球物理上不成像、用户无法用单击唤起（票 62 `docs/evidence/s1/62-diff-baseline/`）。
    现改为静态玻璃体（体径见上，**默认 34.72px**），实测 `px_delta_ge8=2120`、`timers=no`、`cpu=0.000`、
    私有工作集 11.76MB（**D32 的 Sleeping 零定时器纪律未放宽**；可见性来自静态球体本身，不来自动画）。
    **owner 本次验收原话**：「可以说算是赝品吧，离我发的那种质感还是有不小差距，但是先勉强用吧，
    以后再换样式，就先这样吧。开始阶段不能要求太高，本末倒置了就，先完成核心功能」
    ⇒ **本节只批准"尺寸/可见性/零定时器"这三项，不批准玻璃质感本身**；
    质感返工另立**票 65**，届时本节视觉描述须再改一次。

### 2.1 各态视觉（§17.7，原型 `design/screens/ball.html` 为参照；生产是 D2D 非 Web）

| 态 | 视觉 |
|---|---|
| `Sleeping` | 静态玻璃体，体径 = `配置尺寸 × 0.62`（默认 **34.72px**，下限 30px；**不是 44px**，理由见 §2 尺寸条的更正记录），**零动画、零定时器**（INTERIM：原「12px 微点 + opacity 0.35」经差分取证否证；质感返工归票 65） |
| `Armed` | opacity 0.6，直径 44px，静态 |
| `Muted` | opacity 0.4 + 1.5px 斜杠 |
| `Listening` | 56px + 外环随音量波动 + accent 描边 |
| `Thinking` | 球体底部 2px accent 光带流动 |
| `Acting` | accent 常亮实心 opacity 1 |
| `Confirming` | danger 脉冲环（2s，仅一次或低频）+ 环下 micro 倒计时数字 |
| `AwaitingApproval` | danger 环 + **右上角标**（12px 圆，白色数字 = 队列深度） |
| `Speaking` | success 呼吸（1.6s） |
| `Settling` | opacity 1→0.35 渐隐，260ms |
| `Warm` | **warm 暖色微呼吸（2.4s，opacity 0.55↔0.7）**——与 Sleeping 必须明显可辨（日常最常见态） |
| `Conversation` | **danger 常亮环，不得渐隐、不得弱于 Confirming**（麦克风开着的唯一交代） |
| `Downloading` | 环形进度 + 百分比 |
| `Error` | danger 底 + x 图标 |
| `NoNetwork` | fg-tertiary 底 + wifi-off |
| `Unconfigured` | warn 底 + key-round |
| `FirstRun` | accent 底 + 引导脉冲 |
| `WatchdogAlert` | warn 底 + alert-triangle |
| `Queued` | 主态 + 左下 6px info 小点叠加 |
| `Stuck` | warn 底 + refresh-cw |

球体基础形态 =「自发光玻璃滴」：左上高光斑 + 两档透光玻璃体 + 边缘亮线 + `blur(6px)` 状态色
内发光；**各态只换核心发光色 + 外发光色两个变量；球体禁用 `backdrop-filter`**（D2D 必须能一比一
复刻；面板侧模糊不受此限）。

## 3. 状态机权威转移表（C12/D43——20 态 42 条转移；未列出的一律非法）

| # | From | Event | To | Guard / 副作用 |
|---|---|---|---|---|
| 1 | `FirstRun` | 引导完成（模型+目录授权+Key） | `Sleeping` | 写 config.toml；建 DB schema |
| 2 | `FirstRun` | 缺模型 | `Downloading` | — |
| 3 | `FirstRun` | 缺 Key | `Unconfigured` | — |
| 4 | `Sleeping` | 快捷键/单击球 | `Listening` | **建 SessionScope(C31)**；加载 VAD+ASR（冷会话 1–3s） |
| 5 | `Sleeping` | KWS 开启 | `Armed` | 加载 KWS 模型 |
| 6 | `Sleeping` | 无 Key/配置非法 | `Unconfigured` | — |
| 7 | `Armed` | 唤醒词命中 | `Listening` | 同 #4；KWS 暂停 |
| 8 | `Armed` | 静音键 | `Muted` | 停 KWS 推理（模型可保留） |
| 9 | `Armed` | 看门狗超限 | `Armed`+告警 | **不得卸载 KWS**（D32 修订③）；提示用户关 KWS |
| 10 | `Muted` | 静音键 | `Armed`(KWS 开)/`Sleeping` | — |
| 11 | `Listening` | VAD 判停且语音时长 ≥300ms | `Thinking` | 停采集；标点恢复；输入打 taint 源标记 |
| 12 | `Listening` | 超时（首轮 15s/会话内 90s/Conversation 30s） | `Sleeping`(首轮)/`Warm`(会话内) | — |
| 13 | `Listening` | 否决（否决词/Esc/单击球） | `Warm` | **丢弃音频缓冲，不落盘** |
| 14 | `Listening` | 麦克风被占用/拔出 | `Error(audio_device)` | 明示设备名 |
| 15 | `Thinking` | 首 token | `Acting`(有 tool call)/`Speaking`(纯文本) | 面板若开则流式推送 |
| 16 | `Thinking` | 超时/网络失败 | `Error` | 按 D37 分类与重试 |
| 17 | `Acting` | 需确认 | `Confirming`(L1)/`AwaitingApproval`(L2) | L1 走倒计时；L2 入 C18 队列 |
| 18 | `Acting` | 全部工具完成 | `Speaking`(有播报)/`Settling`(静默完成) | — |
| 19 | `Acting` | 重复调用阈值 [3,5,8] | `Acting`+注入提醒 | 达 8 → `Stuck` |
| 20 | `Acting` | 路径冲突（C20） | `Queued` | 面板可见「在等谁」 |
| 21 | `Confirming`(L1) | 倒计时结束未否决 | `Acting` | 执行；`fs.write` 必须 temp+rename |
| 22 | `Confirming`(L1) | 否决（单击球/Esc/KWS 否决词/面板拒绝） | `Acting`（取消该调用回给 LLM） | KWS 未加载时 UI 明示「语音取消不可用」 |
| 23 | `AwaitingApproval` | 队头被**原生侧**允许（F2） | `Acting` | 出队；角标减一；队列非空则留态 |
| 24 | `AwaitingApproval` | 队头被拒绝/**超时 300s** | `Acting`（取消该调用） | 超时前 30s 醒目提示；可一键重放 |
| 25 | `AwaitingApproval` | 队列空且任务完成 | `Speaking`/`Settling` | — |
| 26 | `Speaking` | 播报完 | `Warm` | **麦克风保持关闭**（默认档）；面板保活 |
| 27 | `Speaking` | 打断（快捷键/单击/否决词） | `Listening` | 停 TTS；释放音频输出 |
| 28 | `Speaking` | 播报完（Conversation 模式） | `Listening` | 麦克风开启 + 红色常亮环 |
| 29 | `Warm` | 单击球/快捷键 | `Listening` | **零模型加载**（B4 的解） |
| 30 | `Warm` | 开启 Conversation | `Conversation`→`Listening` | **首次开启需 L2 级隐私确认** |
| 41 | `Speaking` | 播报中检出用户语音（Path C，AEC） | `Listening` | 停 TTS 并释放输出，≤400ms；播报音频不得进 ASR（D47） |
| 42 | `Conversation` | 任务意图（如「帮我做 X」） | 交回文本循环（Path T） | realtime 大脑会话挂起或结束；上下文经 C7 携带（D47） |
| 31 | `Warm` | 90s 无交互 | `Settling` | **Dispose SessionScope**；卸 ASR/TTS；FreeOSMemory；销毁或隐藏面板 |
| 32 | `Settling` | 3s 计时到 | `Sleeping`(KWS 关)/`Armed`(KWS 开) | **RSS 10s 内达标** |
| 33 | `Settling` | 新唤起 | `Listening` | **取消回落**（用户意图优先） |
| 34 | 任意 | 网络探测失败 | `NoNetwork` | **保留当前任务 ctx**，可重试 |
| 35 | 任意 | 看门狗连续 3 次回落失败 | `WatchdogAlert` | 一键重启 |
| 36 | 任意 | panic 被 recover | `Error(internal)` | 记 stack；给诊断包入口 |
| 37 | `Downloading` | 完成/失败 | `FirstRun`/`Error(model)` | **校验 sha256 + manifest 签名（C29）** |
| 38 | `Error` | 用户确认/10s | 回错误前会话态或 `Warm` | — |
| 39 | `Queued` | 前序任务释放路径锁 | `Acting` | — |
| 40 | `Stuck` | 用户介入（继续/放弃） | `Acting`/`Settling` | — |

## 4. 承载分层（D29 确认交互分层）

| 界面 | 承载 | 理由 |
|---|---|---|
| 悬浮球 | 原生 D2D | 常驻，产品脸面 |
| L1 阻止窗口 | 原生提示条（不抢焦点） | 一行字+计时环，2–3s 窗口经不起 500ms 拉起 |
| L2 强确认卡 | WebView（面板）/ 原生降级卡（面板不可用） | 完整命令/参数/风险说明是 CSS 强项；此时已在工作态；**L2 能力不得因面板缺失消失（C27）** |
| 面板（命令/结果/配置/历史/安全/隐私/成本） | WebView | 低频 |

## 5. 面板（WebView2，S5）

### 5.1 宿主（D29/C27）

- `jchv/go-webview2`（MIT，纯 Go 无 cgo）；**单例 `PanelManager`：一会话至多一个 WebView2 窗口，
  隐藏而非销毁**——既绕开 Environment 共享问题，又把绝大多数交互压到热路径（冷 ≤1500ms /
  热 ≤200ms，D32 修正）。
- 资源加载：`AddWebResourceRequestedFilter` 从 `embed.FS` 喂前端产物——**不起 localhost HTTP
  服务**（避免 dsh-tauri 的 iframe-to-localhost 模式，D21/D29）。
- **前端必须无状态**：WebView 销毁/隐藏后状态全丢；每次 `show` Go 侧推 `panel.resync` 全量状态；
  前端不得缓存任何跨 show 的业务状态（历史读 SQLite、配置读 TOML）。
- WebView2 Runtime 缺失 → **无面板模式**：D10 短结果 + 落文件仍可用 + L2 走原生降级卡；
  引导安装 Evergreen Runtime，**不得自动下载安装器**（D42#11）。
- 多任务共用单窗口，按 correlationId 分区渲染。

### 5.2 C17 PanelBridge 方法白名单【SPEC 提案，S5 定稿走契约批准】

`invoke(method, args) → result` + Go→前端事件推送（流式结果/审批请求/任务状态），回复按
correlationId 路由；每方法标注 capability 与是否需原生侧授权；未列出方法名 → 拒绝并记日志。

| method | 方向 | 需原生侧授权 |
|---|---|---|
| `panel.resync` | Go→前端推送 | — |
| `tasks.list` / `task.detail` | invoke | — |
| `history.query` / `transcript.get` | invoke | — |
| `approval.current` / `approval.queue` | invoke | — |
| `approval.decide` | invoke | **「allow」拒绝一切面板来源（F2）；仅 `reject` 可面板发起** |
| `config.get` / `config.set` | invoke | 安全节放宽走 L2 重新确认（SPEC-03 §4.2） |
| `grants.list` / `grants.revoke` | invoke | — |
| `privacy.purge` / `privacy.export` | invoke | purge 需确认 |
| `cost.summary` | invoke | — |
| `models.list` / `models.delete` | invoke | — |
| `diagnostics.export` | invoke | — |
| 事件推送：`task.delta` `tool.chip` `approval.request` `ball.state` `cost.tick` | Go→前端 | — |

### 5.3 前端技术栈（D29 + §17.6，S5 才应用）

React + TypeScript · Tailwind · shadcn/ui（复制源码进仓库，非 npm 依赖）· Radix UI Primitives
（无头原语，键盘导航/无障碍）· `lucide-react`（ISC，零 emoji 硬要求）· `cmdk`（命令面板）·
CVA + clsx + tailwind-merge · Sonner（Toast，但通知被抑制时必须同时更新悬浮球角标）·
motion（仅面板内，禁无限循环动画）· react-markdown + remark-gfm · **rehype-sanitize（白名单）** ·
Shiki（静态高亮，CSP 友好）· @tanstack/react-virtual（取证列表上千行）· Zod（仅输入即时反馈，
真相源在 Go 结构体 tag）· date-fns（仅显示格式化；超时判断一律单调时钟）。
排除：emoji/Font Awesome/Material Icons/MUI/Ant Design/Lottie/Rive/styled-components/图表库。

### 5.4 agent 工作状态的视觉语言（§17.5，14 种状态全部要在 chat 屏画出来）

| 状态 | 视觉（要点） |
|---|---|
| 思考中 | 2px 高光带左→右循环（1.2s）+ 右侧「思考中 · 已耗时 1.2s」（tabular-nums）；不用跳点/spinner/brain 图标 |
| SSE 流式 | 文字按 SSE 分块追加 + 1px×14px 竖线光标呼吸 + 实时 token 计数 |
| 推理过程 | 默认折叠一行（「推理过程 · 3.1s」），展开为 1px 竖线引导的次级文本块 |
| 工具调用 | chip：类别 SVG + 工具名（mono）+ 参数摘要（渐隐截断）+ 状态指示（旋转弧/check/x/ban）+ L0/L1/L2 描边胶囊徽标；展开 = 着色 JSON + 耗时 + correlationId |
| 审批等待 | chip 变 warn 描边 + 队列深度角标；脉冲只跑一次（不无限循环） |
| L2 确认卡 | 顶部 2px danger 横条；标题=工具名+L2 徽标；正文=**完整参数** + C19 命中规则；底部常驻「点击悬浮球以批准」+ 球图示脉冲引导；**按钮只有「拒绝」与「查看完整参数」——没有「允许」**（F2 定案） |
| L1 阻止窗口 | warn 竖条 + 工具名 + 参数摘要 + 20px 环形倒计时 + 「单击球/Esc 取消」；KWS 未加载 → 「语音取消不可用」 |
| L2 原生降级卡 | 球旁 320×140 原生小卡：工具名+参数摘要+「允许」（danger 实心——允许权在原生侧）+「拒绝」+「打开面板查看完整参数」 |
| 错误 | danger 竖条 + alert-triangle + 人话文案 + 可展开 error_class 与技术细节 + 动作按钮 |
| Stuck | warn 竖条 + refresh-cw + 明示重复了什么 + 继续/换个方式/放弃 |
| 成本 | 纯文字 mono tabular-nums「12,480 tok · ¥0.31」+ IN/OUT/CACHED micro 标注 |
| 已取消 | 分隔线 + 居中「已取消」（取消不是错误，不用红色） |
| 注入检出 | danger 竖条 + shield-alert + 「本次操作包含来自 `<url>` 的内容」+ 可展开命中片段 |

## 6. 设计系统（§17.3，"高级感"的可执行规则）

1. 低饱和功能色（绝不用纯红绿蓝/Tailwind-500 原色）；v2.3 冷调环境光四团（青雾/钢蓝/冷靛/暮紫），
   暖色只作语义色（Warm/L1）。
2. 分层靠透光率 + 内侧高光 + 玻璃环；`backdrop-filter` 模糊（面板侧）；1px 线只承担「分隔」。
3. 磨砂颗粒：单层 feTurbulence 噪点 overlay。
4. 间距 4px 网格：只用 2/4/8/12/16/20/24/32/40/48/64。
5. 字号七档（display 30 / title 18 / sub 15 / body 14 / small 12.5 / micro 11 / mono 12.5）；
   数字 `tabular-nums`；section label = micro 大写 + letter-spacing。
6. 动效三档 120/180/260ms，主缓动 `cubic-bezier(0.32,0.72,0,1)`；禁 bounce/elastic/>300ms。
7. 面板宽 640px；正文行高 1.6；卡片内边距 ≥16px。
8. 暗色默认、亮色对照；每屏明暗切换都可读（正文对比度 ≥4.5:1，micro ≥3:1）。
9. 中文排版硬约束：**正文中文 ≥14px**（micro 11px 只用于英文大写 label 与数字）；**禁用 300 及
   更细字重**（Windows 中文细体发虚）。
10. 零 emoji：全部内联 SVG（Lucide 风格，`viewBox="0 0 24 24"`、stroke-width 1.5、尺寸仅
    14/16/18/20px、`currentColor`）；约 55 个图标清单见 §17.4；禁 AI 俗套图标
    （sparkles/wand/brain/robot）；**文字符号也算 emoji**（U+2190–U+2BFF、U+1F300–U+1FAFF、
    U+FE0F 扫描命中必须为 0——可机器验收）。
11. 所有屏幕无硬编码色值（`grep -E '#[0-9A-Fa-f]{6}' screens/` 命中为 0），只准引用 tokens.css。

## 7. 键盘、焦点与托盘

- 键盘优先：命令面板 cmdk 语义；`Esc` 在 Confirming 期间临时接管为取消键，**会话结束必须归还**
  原绑定（B1）。
- 焦点管理：只有 L2 确认卡与文字输入/面板获取焦点；面板销毁/隐藏后**焦点还给记录的原前台窗口**。
- 托盘：右键菜单（打开面板/静音/暂停唤醒/退出）；左键 = 打开面板。
- 全屏/演示检测：前台全屏应用时自动隐藏悬浮球（`[ball] hide_on_fullscreen`）。

## 8. 测试决策

- 结构性：20 态 × 40 转移表做**穷举转移测试**（非法转移必须被拒）；每态超时值表驱动测试。
- C17：correlationId 路由测试（并发 10 个请求乱序回复）、方法白名单拒绝、resync 后前端状态与
  Go 侧一致（面板「无状态」用重开面板断言）。
- 安全：CSP 头实际注入验证；XSS 样本（`<img onerror>`/`<script>`）在面板渲染不执行；
  `approval.decide({allow:true})` 面板来源被服务端拒绝。
- SLO 联动：面板冷/热拉起计时进 SLO 脚本（SPEC-10 §3）。
- **视觉验收（人工）**：原型五屏（index/ball/chat/palette/approval）先行；S1 悬浮球 Direct2D
  与 S5 面板各做一次用户本人验收（「简约高级不丑」是最终判据）；C21 双侧一致性（原生侧取值
  与 tokens.css 一一对应，列出对照表供核对）。

## 9. 不做什么

- 悬浮球不用 WebView 渲染（D2）；不用 `backdrop-filter` 于球体（D2D 不可复刻）。
- 面板不做响应式（固定 640px）；不做 i18n（中文硬编码，S8 前）；不做真实数据以外的演示态
  混入产品构建。
- 不做系统级皮肤/主题编辑器（明暗两套 + tokens 即全部）。
