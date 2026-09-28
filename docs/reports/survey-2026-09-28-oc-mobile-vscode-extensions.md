# survey-2026-09-28 腿1 — OpenChamber 的 mobile／vscode／extensions 三包读穿

> 只读调研件。参照仓 `D:\work\AI\open source\openchamber`（**无 .git，tarball 解包**，全部引证＝`文件:行`，零提交历史引用）。
> 本仓对照＝2026-09-28 现读 grep/ls 计数。产出人＝survey-oc-3。

## 0. 先数枚数（口径与推导式）

| 包 | 口径 | 读数 | 命令（在本机复跑过） |
|---|---|---|---|
| `packages/mobile` | 全部文件 | **110** | `find mobile -type f \| wc -l` |
| `packages/mobile` | `.ts/.tsx`（排 node_modules/dist） | **1**（只有 `capacitor.config.ts`） | 同上 `\( -name '*.ts' -o -name '*.tsx' \)` |
| `packages/mobile` | 真正要读的代码件 | **15 枚**＝1 ts＋6 mjs（scripts/）＋2 java＋5 swift＋1 androidmanifest | 按扩展名清点 `sed 's/.*\.//' \| sort \| uniq -c` |
| `packages/vscode` | 全部文件 | **151** | `find vscode -type f \| wc -l` |
| `packages/vscode` | `.ts`（src 91＋webview 21＋vite.config 1） | **114 枚／src 总行数 28,964** | `find src -name '*.ts' -exec cat {} + \| wc -l` |
| `packages/extensions` | 全部文件 | **3**（README.md／DOCUMENTATION.md／registry.json） | `find extensions -type f` |
| 溢出面（移动端逻辑真身） | `packages/ui/src` 里 `*mobile*/*Mobile*` 命名件 | **61 枚** | `find src -iname '*mobile*' -o -iname '*widget*' … \| wc -l` |

**结论先说**：`mobile/` 目录本身**装的不是应用，是"壳"**——Capacitor 原生壳（`README.md:3`），真 UI 复用 web 构建；所以"逐目录读完 mobile/"物理上=读那 15 枚壳文件（我全读了）+ 追踪 61 枚 ui 侧移动件的**关键 13 枚**（其余点名见 §8）。`vscode/` 的 114 枚我读了全部名册＋按主题读穿约 25 枚核心件，`gitService.ts`(4073 行)与 `opencodeConfig.ts`(3192 行)两枚巨件**只登记未通读**（§8）。`extensions/` 三包里最小：**只有 3 枚、全读**。

---

## 1. 核心问题人话版回答（owner 五问）

### Q1 同一 harness 换到手机／IDE，多了什么桌面没有的？

**手机多出来的（桌面／web 都没有）**：
- **锁屏与桌面小部件 5 枚＋Control Center 一键开新会话**：中等件=最近会话列表(带未读橙点)+四宫格快捷(新建/状态/实例/设置)；大件=带项目名的 6 行会话列表+需关注计数；锁屏圆件=两枚（logo→新建会话 / **铃铛数字=待处理会话数**）。`OpenChamberWidgets.swift:4-321`、`OpenChamberControl.swift:11-38`（iOS 18+ ControlWidget+AppIntent，**OS API 层**，跑在独立 widget 进程）。
- **离线也能刷新的推送件**：通知服务扩展（NSE）收到推送时**不联网**，只读推送里带的 badge 数和 sessionId，改共享区快照、重载 widget（`NotificationService.swift:4-9,36-81`，OS API 层）。
- **扫码配对**：手机摄像头扫服务器二维码完成连接配对（MLKit/CameraX，Android 侧模型打进 APK **离线可用、不依赖 Google Play**，`HANDOFF.md:84-87`；用途写在 `Info.plist` NSCameraUsageDescription="scan a server's pairing QR code"）。
- **多服务器"实例"管理**：手机同时存多台家里/公司的 OpenChamber 服务器，连接信息＋客户端令牌存**iOS Keychain/Android Keystore**（`mobileConnections.ts:8-15`，"native 上令牌绝不进 localStorage"）。桌面版没有这层"连到哪台"的入口。
- **边缘滑动切会话、安卓返回键语义、状态栏配色跟随主题、iOS 外接键盘检测**（`HANDOFF.md:96-98`；`AppDelegate.swift:94-161`：用 GameController 的 `GCKeyboard.coalesced` 权威回答"有没有接实体键盘"，在 document-start 注入 `window.__OPENCHAMBER_HARDWARE_KEYBOARD__` 并在 0.3/1/2.5 秒三次重发压竞态——**OS API 层**）。

**IDE 多出来的**：
- **编辑器内的行评论线程**：右键选区或行号槽 `+`，起一个**钉在那几行上的** VS Code 原生评论线程，"评论=待发上下文卡"，随下一条消息发出；**还能开在 diff 视图上（original/modified 两侧都行）**（`InlineCommentThreads.ts:1-50`、`InlineCommentDraftPayload.source: 'diff'|'file'` :30；`vscode.comments` API，宿主进程层）。
- **右键菜单动作族**：Add to Context / Explain / Improve Code（编辑器选区右键子菜单，`package.json` contributes.menus `editor/context`）。
- **会话开进编辑器标签页**：`openchamber://` 之外还有 `SessionEditorPanelProvider`（708 行，未通读），命令"Open Session in Editor"。
- **订阅配额仪表盘**：约 **20 家**编码订阅的用量窗口（日/周/月/session 五类窗口，含重置倒计时、余额、消费上限）——claude/codex/github-copilot(+addon)/gemini/antigravity/cursor/cline-pass/deepseek/kimi/minimax-coding-plan(含-cn)/openrouter/ollama-cloud/opencode-go/exe-dev/xai/nano-gpt/neuralwatt/wafer/hyper（`quotaProviders.ts:1-60`＋名册 grep，宿主进程层直接读各家 API/auth 文件）。
- **宿主代理全家桶（bridge，32 种消息型）**：webview 没权限干的都在扩展宿主进程代办——fs 十种 CRUD/exec/reveal/pick、git **24 种**（含 stage/unstage/**apply-hunk**/commit/push/pull/fetch/worktree 五种/**PR 描述生成**/冲突详情）、编辑器开文件开 diff（虚拟文档 provider，`bridge-system-runtime.ts:127-145`）、`workspace:addFolder`（聊天下指令往工作区加文件夹）、把 webview 的 SSE 流**在宿主里代接**（3 次重连/1s 起指数退避/20s 停滞超时，`sseProxy.ts:23-25`）。
- **诊断命令**"Show OpenCode Status"：一次列出 server URL、模式、启动/重启/就绪次数、检测端口、API 前缀、上次退出码、连接时长（`extension.ts:755-795`）。
- **扩展自带的 opencode CLI 自升级链**：`api:opencode/install-v2`、`upgrade-status`、`upgrade`、兼容性检查（`bridge-system-runtime.ts:229-246`）。

### Q2 手机上"要人批准"那一跳怎么答？

**不是推送按钮、不是回短信、也不是只读。是三段式：**
1. 服务端事件触发（ready/error/**question/permission** 四类）→ 给手机发 APNs/FCM 推送，**推送内容刻意空泛**："Agent needs permission" 这个固定标题＋**会话名**，无模型名/项目/消息内容（`web/server/lib/notifications/APNS.md:16-21`；`DOCUMENTATION.md:49-50` permission 触发还会先查"该会话是否已开自动接受"，开了就不推）。
2. 点推送 → **深链** `openchamber://session/<id>` 进 App 对应会话（`APNS.md:27`；`deepLinks.ts:23-30` 全套意图词表；冷启动意图**先 stash 不丢**，`deepLinkNavigation.ts:15-16,133-135`）。
3. App 内**照常弹可交互审批卡**——移动壳复用整个 ChatView（`MobileApp.tsx:11 import { ChatView }`），审批卡组件 `components/chat/PermissionCard.tsx`/`PermissionDock.tsx`（存在已核，正文未读，属腿2口径地界）。

配套的**在场抑制**：任何交互端（桌面/web/vscode）有一个可见心跳在，就不给手机推（"Gated on the desktop's visibility, **never the phone's own**"，`HANDOFF.md:108-110`；服务端 `isAnyUiVisible()` 心跳门，`DOCUMENTATION.md:54,65-66`）。iOS 前台则永远不弹横幅（`capacitor.config.ts:32-37` presentationOptions 空集，注释明说：服务端可见性门在 WKWebView 里有竞态、所以干脆**总是发、由 OS 压前台**，`APNS.md:28-38`）。

**跑在哪层**：推送注册=OS API（APNs/FCM）→ 用户自己的服务器（Node 进程）→ **中央中转 Cloudflare Worker**（每服务器自生成 ECDSA P-256 密钥对绑 token，`APNS.md:13-24`；中转不在本仓，属 `openchamber-website` 仓，`APNS.md:87-90`）。

### Q3 移动端怎么处理语音与键盘？

**语音——比桌面还完整，而且是"服务端权威转写"**：
- 采集：浏览器 `getUserMedia` 单声道 → ScriptProcessorNode 抽 PCM16 → **WebSocket 流式**送 `/api/dictation/ws`（`ui/src/lib/dictation/use-dictation-audio-source.ts:4,175`、`dictation-client.ts:2`；安卓靠 Capacitor 把 WebView 录音权限转成运行时弹窗，`AndroidManifest.xml:50-53` 注释写明）。
- 转写：服务端 **fork 的 worker 子进程**跑 sherpa-onnx **Parakeet TDT**（默认本地；也可指任何 OpenAI 兼容 `/v1/audio/transcriptions`），模型首用自动后台下载，缺模型时报 `model_download_in_progress`（`web/server/lib/dictation/DOCUMENTATION.md:46-57,99-100`——模型在 `~/.config/openchamber/speech-models`）。
- 设计账算得很死：**不做增量转写**（Parakeet 全注意力 conformer，O(n²)：实测 60s=2.1s/+90MB、180s=9.3s/+490MB、300s=21.3s/+1.5GB）→ 60s 处遇静默即段落提交、90s 硬上限；**峰值<300 的静音段直接丢**防 Whisper 幻听；ack 只确认连续 seq、客户端留未确认重传（同文档 :63-98）。
- **TTS 播报有，且本地**：Kokoro 与 Piper/VITS 经 sherpa-onnx OfflineTts，`POST /api/dictation/tts/speak`→WAV；每模型声明语种，`language:'auto'` 用无依赖的"脚本+功能词打分"检测整句语种、自动换模型换默认说话人（同文档 :14-32；`ui/src/hooks/useServerTTS.ts`、`useLocalTTS.ts`；播报入口挂在消息体 `MessageBody.tsx`）。
- 设置页语音区：STT 四模型（Parakeet v2/v3、Whisper base/tiny）**带准确度/速度评分条**、逐个下载/删除；TTS 音色列表（Kokoro zh/en 等）（`ui/src/components/sections/openchamber/VoiceSettings.tsx:42-66,219-220,383`）。
- 细节钉子：**转写稿绑定"起录时那条草稿"**——你在转写回来的窗口里切了会话，文字会塞回原草稿而不是当前屏（`useDictationOrigin.ts:1-10`）。键盘快捷键 `toggle_dictation` 可触录（`ComposerDictation.tsx:127,198-210`）。
- **没有**按住说话（hold-to-talk）证据：入口是点按切换（toggle），文件内"hold"无命中；未找到手机 TTS 自动播报（无 autoSpeak 命中），播报是手动按钮。

**键盘**：
- Capacitor Keyboard 插件 `resize:'none'`，自己用 `keyboardWillShow` 事件驱动 CSS 变量 `--oc-keyboard-inset` 抬界面（内置 native resize 要等键盘动画完≈1.5s 延迟，`capacitor.config.ts:19-27`；`mobileNativeChrome.ts:38,119,213`）。
- iOS 实体键盘用 GCKeyboard 权威判定（见 Q1）；安卓用 `adjustResize`＋返回键监听（`mobileNativeChrome.ts:446`）。

### Q4 IDE 版把"我们打算放在球上的信息"搬过去了哪些？

| 我们球上的打算 | OC 在 IDE 的现量 | 证据 |
|---|---|---|
| 当前工作目录 | 有通道但给 webview 用：`api:opencode/directory`＋活动文件广播 | `bridge-system-runtime.ts:177`；`ChatViewProvider.ts:96` |
| 当前模型 | `api:models/metadata`＋配额面板；会话最后一条回复里锁 model 链 | `bridge-system-runtime.ts:192`；quotaProviders |
| 权限档位 | **每会话"自动接受权限"开关，存 VS Code globalState（带 revision＋跨面板广播＋串行写队列）** | `bridge-permission-auto-accept-runtime.ts:1-103` |
| 成本/token | 订阅配额窗口族（20 家）＋opencode-go/xai 等各家 usage；OpenAI credits/spend_control | `quotaProviders.ts:15-52` |
| **关键反差** | OC 在 IDE **一个状态栏项都不做**（全部 src `createStatusBarItem` grep＝0 命中）；所有状态信息留在 webview 里 | grep 现量：`grep -n createStatusBarItem src/*.ts` 无命中 |

手机端同理：会话头只有标题(兼切换器)+一枚"元数据"钮（上下文填充百分比/-token 数/模型，`MobileHeader.tsx:104-144`、`MobileSessionMetadata.tsx:34-36,386-400`）＋工作区按钮上的**未提交改动小圆点**（`MobileHeader.tsx:166-171`）。

### Q5 有没有"离开桌面任务继续跑、回来接上"那套？

有，而且是**"任务本来就在服务器跑"模型**（这是全套设计的根）：
- 服务器才是 agent 循环的家，手机/桌面/IDE 都只是**连上来的客户端**（连接屏=连到"已有的 OpenChamber server"，`mobile/README.md:10-12`）。
- 手机上**关掉的终端**：PTY 跑在活动服务器上，界面关闭只是渲染器 detach，服务器侧会话**留在原地可重连**（`mobile/README.md:16`）。
- 回来的"接上"三件套：**推送→深链→会话**；**app 图标角标＝"自上次打开以来累积的、不同 banner 数"**（按推送 tag=collapse-id 计 Set 大小，非会话数；任何"用户在用 App"的信号即清零——可见性心跳＋打开会话 POST＋发消息 POST 三路兜底，`APNS.md:40-66`）；**widget 未读橙点**用与会话侧栏同一判据 `needsAttention = unseen>0 && (!isSubtask || notifyOnSubtasks)`（`mobileWidgetSnapshot.ts:18-20,84-98`）。
- 事件流断线重建有"权威态对齐"：IDE 侧 `GET /api/session/active` 是唯一真相，"凡不在册的一律 idle，包括本进程断流前以为 busy 的"（`sessionActivityWatcher.ts:34-52`）。
- 子代理跑着时父会话 `session.idle` **不发完成通知**（等带结果重跑后的下一次 idle 才宣布），检查不了就照发（fail-open 有度），`notifications/DOCUMENTATION.md:51-53`。

---

## 2. 缺口三段式清单（他们／我们现量／后果）

> "我们"现量均为 2026-09-28 对本仓 grep/ls 实读；`〔未量〕`＝量不到。**层**＝他们跑在哪一层。

| # | 他们怎么做（带行） | 我们现在是什么样（现量） | 落到我们身上的后果 | 层 |
|---|---|---|---|---|
| G1 | **远端推送审批一跳**：固定四触发 ready/error/question/permission，payload 只含情景标题+会话名，点击深链回会话交互审批（`APNS.md:16-27`；`notifications/DOCUMENTATION.md:49-54`） | Wisp 全仓 Go `push/relay/apns/fcm` 无实现文件（`ls internal/` 22 包无 notify；`grep -rli "APNs\|FCM\|推送通知" docs/specs/`＝**0 份**）；C17 PanelBridge 的"推送"是**进程内** Go→WebView 事件（`PLAN.md:1367`） | 人话：**你合盖出门，审批就永远悬着**。要出家门，最小形态不是 APNs，而是 Windows 原生 toast＋(锁屏可点?)＋协议唤醒——但 Wisp 是单机 harness，服务器=本机，得先回答"进程还在、屏走了"的边界；D43 状态机没有"已通知未答复外出态" | 服务端＋OS API＋中转 |
| G2 | **在场抑制**：任一交互端可见心跳在，手机不推；且"只盯桌面、从不盯手机"（`HANDOFF.md:108-110`）；iOS 前台改由 OS 压横幅、服务端永远发（`APNS.md:28-38`——注释明说可见性门在 WKWebView 有竞态，他们**主动放弃了有竞态的门**） | 无多客户端概念：`grep -rli "visibility\|isAnyUiVisible" internal/` 〔未量到，无对位包〕 | 这条是 G1 的前提。他们踩过坑给的教训：**别用"问客户端你在不在"做门，用 OS 的前台抑制**。若我们做通知，判据该抄这个 | 服务端（心跳表）＋OS API |
| G3 | **app 图标角标＝待处理数**：badge=Set\<tag\> 绝对值，服务端算、三路信号清、设备前台归零（`APNS.md:40-66`；`AppDelegate.swift:201-208`） | 托盘图标存在（`PLAN.md:86,1227`）但**无角标/数字机制**：`grep -rn badge internal/ --include='*.go'` 仅 ball 测试文件命中 | 球=常驻视觉件，天然比角标强；但"多少个在等你"这枚数我们没有现量。落 Win32：托盘 overlay icon / Toast 计数都可行，数据源就是 perm 队列长度（`internal/agent/approval/queue.go` 在仓，**数是从零现成的**） | OS API＋服务端计数 |
| G4 | **移动端 widget 数据桥**：web 层把会话摘要挂 `window.__OPENCHAMBER_WIDGET_SNAPSHOT__()`，原生壳前后台切换时 evaluateJavaScript 拉走→写 App Group→重载 widget；NSE 收推送**不联网**只改快照（`AppDelegate.swift:231-246`；`mobileWidgetSnapshot.ts:108-127`；`NotificationService.swift:4-9`） | 〔未量——无对位场景〕。最接近的东西是球的状态渲染（`internal/ball/renderer_windows.go`） | 人话对应=**锁屏/Win+1 小组件上的"有 3 件事等你"**。Win32 无 App Group/WidgetKit 等价物，这条**落不了原样**；可迁移的是"快照函数＋事件驱动刷新＋离线推送件"的设计形状 | OS API（跨进程共享存储） |
| G5 | **深链词表**：`openchamber://` 7 种意图，纯函数解析、冷启动 stash、未知返 null 不抛（`deepLinks.ts:23-30,46-53`；`deepLinkNavigation.ts:15-16`） | `grep -rli "deeplink\|protocol handler" internal/`＝**0**；PLAN 无自定义协议登记 | 我们要做"点通知回会话"就需要协议注册（Win32 有 HKEY_CLASSES_ROOT 路径）＋**冷启动意图暂存**这个形状；`internal/session` 有会话 ID 体系（D35），词表可以照搬结构 | OS API（注册表）＋进程内 |
| G6 | **扫码配对**：手机摄像头扫服务器 QR 完成实例配对（`HANDOFF.md:84-87`；Android CameraX 模型打进 APK 离线可用） | 无 QR 面（`grep -rli "qrcode" internal/`＝0）。我们的对位问题是"手机/另一台机器怎么安全连上 Wisp"——D36 三档生效级别里没有配对词表 | 自用期若无手机连接需求可不追；一旦要 G1，配对是前置 | OS API（相机） |
| G7 | **令牌进 Keychain/Keystore**："native 上令牌绝不进 localStorage"、删除实例即复位连接屏（`mobileConnections.ts:4-15`；`package.json` dep `@aparajita/capacitor-secure-storage`） | 我们已有 DPAPI：`internal/secret/dpapi_windows.go` 在仓；SPEC-12 待定案里"便携模式×DPAPI"仍是 §2 未决项 | 我们**不缺这层**；但注意他们的纪律可抄：**每个存储面写明"哪层绝不出现令牌"** | OS API |
| G8 | **多服务器实例管理**：连接列表、密码解锁、client token 随连接存（`mobile/README.md:10-15`） | Wisp=单进程单机 harness（D1/D7），无对位 | 结构差异，不是缺陷；但"C24 白名单定稿"若涉及方法面，这个"连接=身份上下文"模型是参照 | 服务端＋进程内 |
| G9 | **语音全家桶服务端权威**：WS 分片协议(seq/ack/重传/60-90s 分段/静音丢弃/adaptive finalize)＋fork worker 跑 sherpa Parakeet＋Kokoro/Piper 本地 TTS＋语种自选（`dictation/DOCUMENTATION.md` 全文；`use-dictation-audio-source.ts:4`） | 计划同源不同工：`internal/speech/doc.go:1-2` 写明 Wisp 语音引擎就规划用 **sherpa-onnx runtime**（同一家族！），但 ASR=票 15、TTS=票 26、KWS=票 41 全 DEFERRED；`internal/audio/doc.go:27` `DEFERRED(playback)`；转写绑定草稿的"origin 钉回"无对位 | **这是本轮对 D5/D25 最硬的一枚参照**：分段阈值、O(n²) 实测表、静音峰值门、防幻听、重传协议——全是可直接抄的数。语音层将来在 Wisp 走子进程（doc.go 已预留"path X speech subprocess"），与他们"绝不在主进程载 native addon"（invariant :91）同构 | 浏览器 OS API＋服务端 fork 子进程 |
| G10 | **实体键盘权威判定**：GCKeyboard＋document-start 注入＋0.3/1/2.5s 重发压竞态（`AppDelegate.swift:94-161`） | 不适用（PC 恒有键盘）；对应问题是**触摸模式笔记本**（我们 D7 没写过触屏分支） | 低优先；登记"屏=球时键盘/触控两态"未决 | OS API |
| G11 | **IDE 行评论线程（可开在 diff 上）**：评论即上下文卡、webview 草稿库为权威、编辑器线程只是镜像、投递确认超时才弃（`InlineCommentThreads.ts:1-56`） | 无对位：Wisp 面板是全屏 webview，无"钉在代码行上"的宿主件（`grep -rli "comment.thread\|inlineComment" internal/`＝0） | 我们的 harness 用户就在 Windows 桌面上，**VS Code 侧的"指着那行说"没有**——若做 IDE 侧是票级工程；至少 owner 问的"点开设置也不少东西"里，右键四动作是别人第一天就有的手感 | 宿主进程（vscode API） |
| G12 | **IDE 审批梯度独立持久层**：auto-accept per session，globalState＋revision＋广播＋串行化队列（`bridge-permission-auto-accept-runtime.ts:1-103`）；服务端 permission 推送触发先查这层（`notifications/DOCUMENTATION.md:49-50`） | 我们有 `internal/perm`＋`internal/agent/approval/{gate,queue,batch,pending_read,ui}.go` 在仓；R20 三档（D4/D31）；**"通知侧查自动接受"这层联动**无对位（因为无 G1） | 抄点：把"自动接受"当**可广播的快照(revision 单调)**而不是散 flag——这与我们 D43 状态机/快照泵(票 35/145)的哲学同形，可直接对话 | 宿主进程存储 |
| G13 | **订阅配额仪表盘 20 家**：日/周/月/session 窗口+重置倒计时+credits+spend control（`quotaProviders.ts:15-52`＋名册 grep） | `grep -rn "配额\|quota" docs/specs/SPEC-04* docs/specs/SPEC-08*`＝**0**；球/面板规划只有成本 tab（`SPEC-08:139,200`：「12,480 tok · ¥0.31」） | 人话：owner 问"点设置一堆东西"——**这屏就是一堆**。我们有 adapter 族（internal/llm 含 anthropic/openai-compat），窗口型配额的**数据形状**（usedPercent/resetAfterSeconds/windowSeconds）值得进我们面板成本页规格 | 宿主进程直连各家 API＋读本地 auth 文件 |
| G14 | **宿主代理面 32 种 bridge 消息**：fs/git/编辑器/工作区/升级（§Q4 行集） | C17 PanelBridge 是 invoke+事件双工但**方法表白名单未定稿**（AGENTS.md §2：C24/C17 初始集待拍板，**阻塞中**） | 这枚是**给拍板用的现量清单**：32 个方法名就是别人定稿的候选集，git 24 操作(含 apply-hunk)+editor:openDiff 可逐条对照我们的 D34 工具权威表 | 宿主进程 |
| G15 | **GUI 程序拿不到用户 PATH 的 Windows 解法**：pwsh `Get-ChildItem Env:`→powershell.exe→System32 全路径→cmd `set` 四级回退、NUL 分隔解析、windowsHide、10s 超时（`opencode.ts:529-575`）；找不到再 `where opencode`（:447） | 我们的模型分发/插件定位（D26/D46）会在服务/计划任务语境撞上同一坑；`grep -rln "Get-ChildItem Env\|环境快照" internal/`〔未量到对位〕 | **直接抄的算法**。后果：不做这个，"装了 CLI 但 Wisp 说找不到"就是我们工单池第 N 枚幽灵 bug | 子进程（spawnSync） |
| G16 | **断流后权威态对齐**：重连不回放、改拉 `session/active` 全集、不在册一律 idle（`sessionActivityWatcher.ts:34-52`）；SSE 代理 3 连退避＋20s 停滞超时（`sseProxy.ts:23-25`） | 我们有 watchdog/observe 包；"恢复后先对齐再演"在 D37 错误传播有位置但**无这枚具体判据**（`grep -rn "session.active" internal/`＝0） | D37b/票 174"回执说实话"同题：抄"**权威集合替换乐观集合**"这条规则 | 服务端 API＋客户端 |
| G17 | **子代理未完不宣布完成**：父 idle 而有子代理在跑→不发完成通知；查不了才照发（`notifications/DOCUMENTATION.md:51-53`） | 票 197 子代理三层刚批准；`internal/agent` 无对位门（owner 09-28 硬要求的"点击子代理看各自流式页"也还没有） | 这条直接喂票 197：**完成语义=树上全员完成**，不是父会话空转 | 服务端事件门 |
| G18 | **内置扩展信任模型**：空 registry＋保留前缀＋iframe 无宿主上下文＋绑定后"HTTP 请求/清单都改不了 binding"＋用户装不得保留命名空间＋禁用态与数据分离（`extensions/DOCUMENTATION.md:11-19`） | D3=goja 内嵌（无 iframe 先例）；`internal/plugin/doc.go` Tier-1/2 **全 DEFERRED**（票 50/51）；无市场（"plugin distribution is S7"） | 他们 registry 是**空的**（`registry.json:2 []`）——"看着一堆内置其实零内置"；真正可抄的是**授权推导规则**："automatic approval is not unrestricted authority"（DOCUMENTATION.md:15）与我们 R20/Q-49 面板侧批准同题 | 服务端（grant 绑定）＋iframe 客户端 |
| G19 | **扩展面板有 12 语 i18n 且禁跨 iframe 用宿主 i18n**（`extensions/DOCUMENTATION.md:8`；vscode 侧 l10n fr/tr） | 我们零 emoji 规范有、多语言〔未量〕无对位规划 | 不追；登记"若上扩展面板，翻译得自带" | 进程内 |
| G20 | **语音/键盘/审批之外的移动工程账**：native 流式**锁死 SSE 禁 WS**（安卓原生 WebSocket 不稳，`HANDOFF.md:144-146`）、Android WebView 111+ 才认 `color-mix()`（:141-143）、iOS 模拟器构建临时剥 MLKit pod（:138-140）、APNs sandbox/production 按 profile 自判随 token 上报（`AppDelegate.swift:70-80`＋`APNS.md:82`） | 不适用 Win32/WebView2（Edge WebView2 常新） | 价值在反面：**他们各家屏的兼容账全部写在文档里**——我们的 D22"未定义即停"该产的就是这种账 | OS API/构建层 |

## 3. `extensions/` 单列（三包里最容易被高估的一包）

3 枚文件全读。它是**"内置扩展的构建与信任登记"**，不是扩展本体：`registry.json` 现在 `{"extensions": []}`（registry.json:1-3），README 自己承认"The registry is currently empty"（README.md:17）。真机制分布在三处：`packages/sdk`（公开契约：contributes.commands/panel(dock/entry/icon/name/size)/tools＋network/storage/session 能力词表〔src 正文未读〕）、`packages/web/server/lib/guests/`（builtins.js 解析 registry、catalog.js 绑实例持久路径、grant-scope/host-session/auth-store/background 等 20+ js〔只列名未读〕）、构建器 `scripts/build-built-in-extensions.mjs`（DOCUMENTATION.md:7）。**且 DOCUMENTATION.md:19 白纸黑字：扩展面板目前只支持 web/Electron，"VS Code and mobile keep the existing explicit unsupported behavior"** ——手机和 IDE 里扩展面板是明确不支持的。

## 4. 覆盖与深度声明（防"目录级抽样冒充覆盖率"）

- `packages/mobile/` **15 枚代码件＋2 枚 md：100% 逐行读完**（462+ 行 swift/java/mjs 全读）。
- `packages/extensions/` **3/3 全读**。
- `packages/vscode/`：91 枚 src 全部**列名＋行数点名**；按主题通读/读穿约 25 枚（bridge-permission-auto-accept 全文、bridge.ts/extension.ts/sseProxy/project-setup/skillsCatalog/InlineCommentThreads/quotaProviders 头部 40-60 行、bridge 全部消息型名册 grep〔32 种 api:*＋git 24 操作逐 case 行号〕、package.json contributes 全量解析、README/nls 全读）；**未通读**＝gitService.ts(4,073)、opencodeConfig.ts(3,192)、SessionEditorPanelProvider.ts(708)、webviewHtml.ts(410)、bridge-config/fs/git 正文、webview/api 21 枚（其中 `webview/api/permissions.ts` 存在已登记、未读）、全部 15 枚 *.test.ts。
- 移动溢出面 `packages/ui/src/apps/` 61 枚命名件：**全读 6 枚**（useNativePushRegistration/deepLinks/mobileWidgetSnapshot/MobileHeader/mobileConnection 头部/MobileSessionMetadata 关键段 grep+行读）；**MobileApp.tsx(1,394)/mobileConnections.ts(1,744)/mobileNativeChrome.ts(462)/deepLinkNavigation.ts(194) 只读到关键行/关键段**；MobileChangesSurface、MobileFilesSurface、MobileTimelineList、MobileSessionSwipe、mobileQrScan、MobileInstancesSurface、MobileFullscreenSurface 等**只登记未读**。
- 服务端为答四问**定向**读了 `web/server/lib/notifications/`（DOCUMENTATION.md＋APNS.md 两份文档全文）与 `lib/dictation/DOCUMENTATION.md` 全文；这些目录的 js 实现体**未通读**；`lib/relay/` 只知存在。

## 5. 我这腿没读到什么（具名到目录）

1. `packages/ui/`（除上表点名片段）：整个桌面/web 共享 UI、会话侧栏 attention 真身、`PermissionCard/PermissionDock` 正文、`ChatInput.tsx` 语音入口正文。
2. `packages/web/server/` 其余全部（含 `lib/guests/*.js` 实现体、`lib/projects/project-setup.js` 服务端镜像、terminal PTY 服务端实现、relay 目录）。
3. `packages/sdk/` 正文（只读了文件名与 grep 出的能力词表）；`packages/electron/`、`packages/docs/`（不在本腿题）。
4. vscode 大件：`gitService.ts`、`opencodeConfig.ts` 全文、`SessionEditorPanelProvider.ts` 正文、`webviewHtml.ts` 正文、`webview/api/*` 全部 21 枚。
5. 中央中转（Cloudflare Worker）**不在这份克隆里**（在 `openchamber-website` 仓，`APNS.md:87-90` 点名）——我只拿到它的协议与表名（`push_tokens`/`0002_push_tokens.sql`）。
6. `mobile/ios/App/App.xcodeproj/project.pbxproj`、gradle 脚本正文（构建接线只核了 HANDOFF.md:117-120 的描述）。
7. PWA 端 service worker 的 web push 实现体（`push-runtime.js` 读的是文档不是实现）。

## 6. 我认为编排者会误读的地方（可能推翻他哪条账）

1. **"OpenChamber＝桌面＋web 两屏"（上一轮目录级抽样的直接推论）——推翻**。它有**四条一等屏**：Electron、web/PWA、Capacitor 手机壳（含推送/widget/深链）、**VS Code 扩展 2.0.3**。任何"人家只是桌面套壳网页"的判语（第二版差集清单 §13 若写过类似句）需要加限定。
2. **"移动审批靠推送按钮一步批"或"移动只读不给批"——两个都不对**。他们**故意不用推送按钮**：推送内容空泛、唯一作用是深链把人拉回可交互会话（`APNS.md:16-27,110-130`）。若编排者按"推送即审批"去对位我们的 D4/门控，会把"审批上下文展示"这个必须性丢掉。
3. **A409 或后续若有"五家都没有移动端通知回路"类猜测**——对 OC 不成立：完整"APNs/FCM＋中转签名绑定＋badge 语义＋在场抑制"是有文档有路由有测试的成品（`notifications/` 目录 19 枚文件里 6 枚 .test.js）。但注意 G1 末段：**它的"服务器"必须是常驻可达的**，纯单机 Wisp 照搬的代价在"谁来当中转"，别在票里把它写成"抄个托盘 toast 就算对齐"。
4. **"他们扩展生态很丰富"——错读方向**。内置 registry＝**空**（registry.json:2），且扩展面板**明拒 vscode/mobile**（DOCUMENTATION.md:19）。这包的真正交付物是**信任/授权模型**（保留命名空间、grants 由声明推导、禁用与数据分离），对位的是我们 Q-49/D19/D46，不是对位"应用商店"。
5. **语音这条最可能被读歪**：`internal/speech/doc.go:2` 已写明 Wisp 规划用 **sherpa-onnx runtime**（与 OC 同引擎家族），所以这不是"人家有我们没有"的功能差，而是**同一技术路线上人家已交的实测账**（O(n²) 分段数、静音门、sandbox/生产环境自判）。编排者若把 G9 拆成"接入云端 STT"票就浪费了这个同源红利。
6. **APNS.md:132-135 自己说 Android FCM "not implemented"，但 HANDOFF.md:104 与 manifest/committed google-services.json 说已实现**——同仓文档互相矛盾、APNS.md 该节过期。**引用时别拿前者当设计决定**；这是"他们文档也会烂"的现量，正可反过来佐证我们"真相源台账＋只追加"的必要性。
7. `packages/vscode` 的"Settings"命令**不是自建设置页**，是把用户弹到 VS Code 原生设置筛选项（且整个扩展只有 2 枚可配项：apiUrl/opencodeBinary，package.json contributes.configuration）。owner 那句"点开设置也不少东西"的"东西"其实在**webview 内的共享设置屏**（VoiceSettings 那类，属腿2/UI 地界），不在 IDE 壳层——编排者合并差集清单时别把两层设置混成一栏。
