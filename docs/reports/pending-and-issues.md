# Reports — 问题项 / 阻塞项 / 待人项登记

> 持续更新。条目解决后移入文末"已解决"。格式：[级别] 标题 — 状态 — 阻塞什么 — 下一步。

## 待人工审核（pending-human-review）

- [H1] **悬浮球 20 态视觉签收** — **2026-09-20 已改为实况签收**：编排者用 `build\balldebug.exe -stay`
  把真悬浮球启到用户桌面（热键/托盘/点击全活），用户直接看与摸，不必再对 PNG。
  截图对照仍是备选：`docs/evidence/s1/ball-states/*.png` 与 `design/screens/ball.html`。
  满意则关闭；不满意提修改意见转新工单。
- [H2] **LLM API Key** — 等用户自助录入（**不得经对话框传 key，见 R7**）— 阻塞票 09 的黄金 SSE 录制
  （录制需要真 provider 响应；mock 与回放框架不阻塞，已先行）与票 12 的真链路验收、票 60 门控条件②、
  票 61 云级联。临时可用路径：用户在自己终端执行 `setx WISP_LLM_KEY "<key>"`（编排者只引用变量名）。
  **残缺现状**：全仓尚无录入入口（`wisp secret` 子命令未做、配置 GUI 属票 39/40 即 S5），
  已提议把最小录入入口提前，待用户批准。
- [H3] **P10 命名残余核查**（npm/PyPI/域名/商标）— **按 R6 降级：发布前才需要，个人自用阶段不做** —
  仍阻塞票 56（票 56 本身维持 DEFERRED）。用户已在商标网登录，填法见本条附注：
  国际分类填 `9`（软件/已录制程序）与 `42`（软件服务），查询方式按文字类型选（汉字→`一缕`；
  英文→`WISP`），类似群可留空（留空=全类近似），检索要素 1–20 字符；
  结果里同类目/同类似群出现「已注册 / 申请中」同名或近似名即为冲突，冲突时备选名
  **只能由用户从 Mote/唤/聆 中选**（票 56 明文：不得发明新名）。
- ~~[H4] `web.search` 实现路径~~ — **已裁定 R4：搜索 API（anysearch）**，移入文末"已解决"。
- [H6] **Sleeping 态微点在浅色桌布上不可见**（2026-09-20 实况签收发现）— 等用户选方案 —
  **不是渲染缺陷**：差分实拍证明球一直在正常绘制（同区域有球/无球仅 87/5184 像素变化，
  证据 `docs/evidence/s1/ball-contrast/`）。根因是 SPEC-08 §2 第 29/35 行规定
  Sleeping = 直径 12px + opacity 0.35 的浅色微点，叠在白色桌布上肉眼等于零。
  **改它属契约变更（D22），需用户批准**。三个候选：
  ①微点加 1px 描边或投影（任意底色都可见，改动最小，推荐）；
  ②opacity 0.35→0.6 + 直径 12→16px（仍"低调"但有下限保证）；
  ③保持规格，新增"浅色/深色桌布"配色档（改动最大，票 39 配置界面正好承接）。
  现状影响：用户看不见常驻球 → 无法用单击唤起（只能靠热键，而热键另有 A 系列缺陷）。
  **→ 2026-09-20 用户裁定：①②③全部不采用，视觉整体重做（见 R8），本条由票 62 承接关闭。**

## 平台事件记录（platform-incidents）

- [P1] **子代理验证码/配额故障（2026-09-19 持续）**：当日 9 次子代理派发被
  "Captcha instance timed out"/"exceed quota limit" 打断（时长 38s–2.9h 不等）。
  应对：票据进度日志 + 断点续传协议（零工作丢失实证：票 05 五棒接力、票 07 三棒接力均无返工）；
  验收类工作由编排者亲自执行（独立性满足）；并发按约定 3→2 回落。
- [P2] **GitHub 直连 TLS 间歇失败**：push 常规重试 ≤5 次可消化；未造成丢失。
- [P3] **`git add -A` 两次吞并行 WIP**（T03/T02 期间）：已按 `git rm --cached` 先例修复；
  此后全仓强制显式路径提交，未再发生。

## 裁定记录（orchestrator rulings）

- **[R1] matcha-zh-baker P3（非商用）**（2026-09-20）：个人使用期可继续用（许可允许）；C29 manifest 分发被挡（blocked-p3 硬拒已实现）。**默认 TTS 替代选型归票 26**（与 P7 音质门禁一并评测；候选： sherpa vits-zh 许可干净系列 / 云 TTS）。
- **[R2] 截止计划**（2026-09-20 07:2x，配额 09:00 到期）：在途 14✅/08（deadline-pause 规则已发）；**票 10（2–4h）与 12 留给新会话**——交接文档 docs/reports/HANDOVER.md 已含新会话行动清单。
- **[R3] 车队并发上限 2→3（用户 2026-09-20 批准）**：旧值「并发 2 为底线、3 为试验位」是 **ZCode 免费配额**校准的产物；用户已切换 Qoder 付费。新上限按实测定：**写码代理稳态 3（极限 4）、编队 3 实现+1 验收、需安静测量的票独占 1、只读检索 10+**。瓶颈已从平台配额转移到本机资源（6C12T / 32GB 空闲仅 ~9.4GB / C 盘剩 30.5GB 且 GOCACHE 在此）。>1 个写码代理必须 worktree 隔离 + `third_party` junction 复用。已同步改 `.scratch/wisp/issues/README.md` 规则 2。
- **[R4] H4 web.search 路径已定：搜索 API**（用户已有 **anysearch** key，2026-09-20）。票 22 的 `SearchProvider` 接口按 anysearch 兼容实现为默认 provider，抓结果页方案降为备用/离线兜底；票 22 的"OPEN decision"验收框就此可勾。
- **[R5] 默认 TTS 走云 TTS，供应商 StepFun**（用户已有 StepFun 云 TTS 模型，2026-09-20）。R1 的替代选型不再走"换开源中文 TTS 权重"路线：票 26 的 TTS 输出层以 StepFun 云 TTS 为默认实现，与票 61（云级联 ASR/TTS C9）同源，减少一家供应商。⚠ 本地 matcha 条目**保持 `blocked-p3` 不改**（改回 ok 属 D36 级安全决策）；离线/无网场景的本地 TTS 缺口另立待决项，不在票 26 里偷偷塞。
- **[R6] 发布类事项全部推迟（用户取向：先自己用，用好了再说别的）**：SignPath 代码签名资格（P8/票 56）、生产 minisign 密钥仪式、P10 命名/商标核查、winget/Scoop 分发 —— 个人自用阶段一律**不做**，票 56 维持 DEFERRED。dev minisign 密钥（`E:\work\base\wisp-minisign\`）继续用于本地更新/模型链验证。H3/H5 由"待人项"降级为"发布前才需要"。
- **[R7] 凭录入方式约束（用户 2026-09-20 明确要求）**：**任何 API key 不得经对话框传递**，必须由用户自己在界面/本地隐藏输入里录入。现状核查：`internal/secret`（票 06，DPAPI）已有 `Store.Store(ref, secret)`，但**全仓没有录入入口**（无 `wisp secret` 子命令、配置 GUI 属票 39/40 即 S5）。因此 H2/anysearch/StepFun 三项凭据当前只能走「用户自己 `setx` 环境变量 + 我只引用变量名」的临时路径；已提议把最小录入入口（`wisp secret set` 隐藏输入）提前做，待用户批准。→ **用户已批准，落成票 63**。
- **[R8] 悬浮球视觉整体重做，否掉"描边/投影"修补路线**（用户 2026-09-20，原话：「我也不知道你说的什么劳什子描边投影是什么东西，估计也是上不了台面的样式，按照我说的来」）。H6 的 ①②③ 三个候选**全部不采用**。新方向（参考 vivo 蓝心小V 语音球，用户附图两张）：玻璃质感外壳 + 内部蓝紫青液态色团；**靠边自动收缩半隐（"就像迅雷那种悬浮球"）**；唤起时液态流动动画；**说话时液体随音频节奏旋转**；不说话时通过过渡动画丝滑引入边框。→ 落成**票 62**。
  ⚠ 我在票里写死的取舍：D32 规定 Sleeping 态 CPU ≤0.5% 且**零动画零定时器**，所以流动动画只允许出现在会话态，Sleeping 必须是**静态玻璃球本体**（不再是 12px 微点）——"看得见"靠本体可见性达成，不靠给 Sleeping 加动画。票 62 的 AC 已钉住这一点。
  契约变更顺序：**先可运行原型 → 用户实机签收满意 → 才回填 SPEC-08 §2 并请签字**；签字前不得改契约文本。
- **[R9] 外部服务端点登记（非机密，用户 2026-09-20 提供）**：
  - **AnySearch**（票 22 `web.search`/`web.fetch` 的默认 provider）：`baseUrl = https://api.anysearch.com/v1`，
    文档 `https://www.anysearch.com/docs/api-endpoints/v1-search`。接口面：
    `POST /v1/search`（`query` 必选；`max_results` 1–10 默认 10；`tag` 形如 `{domain}.{sub_domain}`；
    `zone` = `cn|intl`；`language` = `zh-CN|en`；`params` 透传 AnyMix；`format` = `json|markdown`）、
    `GET /v1/sub-domains`、`POST /v1/extract`。认证 `Authorization: Bearer <key>`，
    **可选**——匿名按 IP 限流并消耗每日免费额度（→ 无 key 时仍可降级工作，这条要写进票 22 的降级链）。
  - **StepFun**（票 26 云 TTS 与票 61 云级联、H2 真 LLM 黄金流的候选）：`baseUrl = https://api.stepfun.com/step_plan/v1`。
  - key 本身一律走 R7/票 63 录入，**不得出现在对话、argv、日志或 config 明文里**。
- **[R11] 票 19 收口 → 14/61 done**（2026-09-20，编排者亲自验收，独立性=实现者≠验收者）。
  验收链：首轮对抗 **FAIL（2 BLOCKER：同步根拼写绕过并实测写盘进 OneDrive；具名通道按参数名逃逸）**
  → 修复（祖先锚定，未照抄复核者那行会误伤正常写盘的修法）→ 复验 **PASS WITH CONDITIONS**
  → C-1 修复 → **编排者读码追加 C-3**（`writeGate` 只判第一个 path 形键，诱饵 `path` + 真 `dest` 逃逸，
  探针实跑 3/4 形状漏）→ C-3 修复 → 编排者复验（自跑 vet/`-count=2`/`-race` 全绿 + 四条 writeGate 测试 PASS
  + 变异证据：还原旧码后新测试 6/7 红而其余 72 个顶层测试仍全绿）。
  冻结面零 diff，`TaintHit(params)` 缝未加宽。
  **移交残余（三条，不得无声消失）**：①N-11 `rules_gateway.go` 的 R4 注释仍写「Dormant until ticket 19
  wires the TaintDetector」，缺的其实是循环接线 → **票 20/21 接线时顺手改准**；②N-8 具名 `fs.write` 上契约外键
  在非同步目录仍被扫（键名依赖的**偏严**不对称，已入 PRECHECK §P12）→ **票 20/21 别当 bug 报**；
  ③C-2 env 级同步根证据可被同用户进程伪造（信任边界即此）+ TOCTOU 不属本票威胁模型 + POSIX 分支不可达
  → **归票 55** macOS 端口收口。
  → 票 19 done 后**票 20（host bridge + fs 工具）解锁**，它是污染门的第一个真实消费者。
- **[R10] 默认唤起热键 `Ctrl+Alt+W` → `Ctrl+Alt+Q`**（用户 2026-09-20 答"换"）。实测依据：本机 84 个候选组合中
  **只有 `Ctrl+Alt+W` 被第三方程序占用**，其余 83 个空闲；`Ctrl+Alt+Space` 也空闲，但因部分中文输入法会抢
  该组合，只作可配置备选、不作默认。→ 接线与错误语义拆分（"被占用" vs "未尝试"）归**票 64**。
  ⚠ 同时修正一条误判：补裁阶段报告称"四个默认热键全部注册失败"，实为**当时我的演示球进程还在跑**、
  占了其中三个键（`TestBallLiveLifecycle` 也因 `FindWindowW` 按类名匹配到演示球而变红）。
  教训已入票 64 的前置断言要求：测量类测试发现外部同窗口必须 **fail-fast 并提示，禁止静默 skip**。
- **[R12] 票 12 的定位改为"接线票"，不是"集成收尾票"**（2026-09-20，验收票 11 时定）。
  连续两张票以**同一形状**留下未勾的 AC：票 63 AC#6（A8，真实 provider 从 DPAPI 取密钥）与
  票 11 AC#6（A11，能力探针没有生产调用者）——都是"测试里已经证明、装配根 `cmd/wisp` 从未调用"。
  这不是两张偶然的残票，而是**当前拆票方式的系统性盲区**：把"做一个能力"和"把它接上"拆给不同票时，
  后者没人认领，能力就以"完成度很高但线上永不生效"的状态沉底。
  ⇒ 裁定：票 12 开工时**先消化 A8/A11 两张外来框再写自己的新功能**；今后凡"新增一个能力"的票，
  票面必须显式回答"谁调用它、在哪个装配根"，答不出就把接线写进同一张票的 AC，不许转嫁给"以后"。

## 完成度审计（2026-09-20 08:58，audit-A + audit-B）- **结论：13 张 done 票无 BLOCKER；实质工作全部有证据支撑**（audit-B 深查 4 张的引用抽验全属实）。
- **[A1] Audit-A MAJOR→已修**：票 02 归档漏改 Status 字段（已修，commit 08f9ff8）。CGO=0 全仓限制为已登记基线（HANDOVER §6），非违规。
- **[B1] Audit-B MAJOR → 转新会话第一组簿记任务**：13 张 done 票的 AC checkbox 未逐条勾选（71 条 `[ ]`）——**实质裁决都在 docs/evidence/ 各验收报告里**，勾选动作纯簿记；票 18 的"perf ≤1ms"AC 为实现自认 DEFERRED（bench 未写），已在票内登记，新会话补 bench 或正式降级。新会话开工前先花 15 分钟做此勾选 + 逐条裁决行。
- **[B2] MINOR×6 全部接受为簿记**：日志时序（02/17/18）、报告路径笔误（17/18）、票 04 双时间戳、票 18 缺 handoff 行、H5/H6 编号、票 07 报告引用文件名笔误（内容在 pending-and-issues.md 未丢失）——新会话顺手清理，不阻塞任何事。
- **[B3] H2/H5 登记补强**：票 16 真机三场景（原 H6 漏登 registry）已由本条补登。

## 阻塞项（blocked）

（当前无硬阻塞。票 09 的黄金录制部分依赖 H2；其余在途/排队票均可推进。）
- ~~[H5] TTS 模型选型（P3 BLOCKED：matcha-zh-baker 非商用）~~ — **已裁定 R5：云 TTS（StepFun）**，
  移入文末"已解决"。⚠ 本地 matcha manifest 条目**保持 `status=blocked-p3` 不变**（改回 ok 属 D36 级
  安全决策）；`Manager.Ensure` 对它的拒绝与 `TestP3BlockedModelRefused` 继续生效。
  遗留缺口：**离线/无网场景没有可用中文 TTS**——该缺口按 R6"先自用"暂不排票，但不得无声消失，
  登记于此（完成判据：票 26 的降级链里出现"无网时静默降级为文字+通知"的明文行为，或选定本地开源模型）。
- ~~[H5-旧编号/HANDOVER 的 minisign 项] 生产 minisign 密钥仪式~~ — **按 R6 降级：发布前才需要**。
  dev 密钥（`E:\work\base\wisp-minisign\`，公钥硬编码 `internal/buildinfo`）继续用于本地更新/模型链验证。

## AC 补裁遗留（2026-09-20）

> 来源：对 13 个从未被书面裁决的 AC 框做补裁（票 05×2、06×1、07×4、08×2、14×4；票 07 AC#2 为
> 待人项 H1，不在补裁范围）。**9 框判 PASS 并改勾，2 框 FAIL、2 框 PARTIAL 保持未勾**。
> 逐条证据与实跑输出见各票 `docs/evidence/{s0,s1,s2}/0X-adversarial-acceptance.md`
> §"Addendum 裁决（2026-09-20，AC 补裁）"。下列为**必须由后续票消化**的遗留。

- **[A1] 票 07 AC#5 热键：配置变更→重注册这条链根本没接线（FAIL，最高优先）** —
  `Ball.RebindHotkeys`（`internal/ball/ball_windows.go:693`）**全仓零调用者**；`[hotkey]` 热重载
  不驱动重注册。且 `registerAll` 失败只 `slog.Warn` 并丢弃该 id（`internal/ball/hotkey_windows.go:189`），
  **没有任何测试断言注册集合**——本机 winlive 实跑四项默认热键全注册失败、测试仍报 PASS。
  阻塞什么：票 07 不能算真完；票 39（配置 GUI）改热键后将静默无效。
  下一步：新票做 ①`config.Manager.OnReload`→`RebindHotkeys` 接线 ②注册集合四项齐的断言
  ③区分"被他人占用"与"未尝试"的错误语义 ④`OnMuteHotkey`→`EvMuteKey`→`Muted` 端到端用例。
  完成判据：票 07 AC#5 框可勾。
- **[A2] 票 07 AC#3 交互四项从未验证（FAIL）** — 现仅 `HitTest()` 纯函数几何分区有测试；
  真实点击→`EvSummon`→`Listening`、Confirming 下 Esc/点击取消＋会话归还（B1）、
  **"焦点不被抢（需用焦点编辑器验证）"** 三项零证据（NOACTIVATE 只经代码审读）。
  下一步：与 H1 实况签收并做一次人工交互签收，另加一条按 hwnd（非类名）查找的 hermetic live 用例。
  完成判据：票 07 AC#3 框可勾。
- **[A3] 票 07 AC#4 动画定时器句柄断言不在默认套件内（PARTIAL）** —
  `TestAnimationPolicyZeroTimerInSleeping` 断的是策略表返回值；句柄级断言仅存于
  `//go:build windows && winlive` 门控的 `live_windows_test.go:95`，**`go test ./...` 与 CI 均不覆盖**，
  且 `Ball.TimersAlive()`（`ball_windows.go:387`）返回的是镜像策略的 `animTimerActive` 布尔。
  另：AC 首句"Sleeping CPU≈0"仍未在真实 runner 量得（票 08 自记首次 CI 运行为最终验证点）。
  下一步：句柄级断言升进默认套件或把 winlive 接入 nightly CI 并转绿；slo-full 上量出 Sleeping CPU≈0。
- **[A4] 票 07 AC#6 "拖到第二屏"从未执行（PARTIAL）** — 本机物理单显示器；
  `enumMonitors()`（`internal/ball/monitors_windows.go:40`）**无注入接缝**，双屏夹具只喂纯
  `ResolvePosition`，从不喂进 `Ball`。下一步：给 `Ball` 可注入的显示器快照做端到端 DISPLAY2 恢复，
  并在真双屏机/D42 显示器拓扑预演上补真实跨屏拖拽 + DPI 资源重建（SPEC-10 §7 明文"必须真实触发"）。
- **[A5] live 球测试非 hermetic（补裁期间新发现，影响 A2/A3/A4 可裁定性）** —
  `live_windows_test.go:30-34 findBallWindow()` 用 `FindWindowW("WispBallWindow", 0)`，
  **跨进程按类名匹配全桌面**：用户 H1 签收常驻的 `balldebug.exe`（PID 29844）同持该类名窗口，
  使 `TestBallLiveLifecycle` 在 `:65` 稳定假红（"ball window still alive after Close"）。
  ⇒ 该红**不可**用作泄漏证据。下一步：测试改按自身 hwnd 断言，并保证运行期桌面上无其他 Wisp 球。
- **[A6] 未定案的句柄增长（需安静环境复跑归类）** — 同一进程内 `-count=2` 两次观测到
  `handlesOfProcess()` 起点基线由 **104 升至 376/374**（该计数走
  `GetProcessHandleCount(CurrentProcess())`，是**进程内**值，外部进程无法抬高，不受 A5 混淆影响），
  ⇒ `Ball.Close()`（`ball_windows.go:762`，其中 `pDestroyWindow.Call` 返回值未检查）后本进程约 270 个
  句柄未回落，与票面约束"release path must bound non-Sleeping handle growth"相悖。
  **可能是真泄漏，也可能是 D2D/USER 对象待 finalizer；本次未能定案**（本机有他球常驻且不可中止用户签收）。
  下一步：在桌面**无其他 Wisp 球**的环境跑 `go test -tags winlive ./internal/ball/ -run TestBallLiveLifecycle -count=3`
  并记录基线是否单调增长，再决定是否立泄漏修复票（另参 D42 "GDI 泄漏"预演项）。
- **[A7] 票 08 诊断包两处残余（不阻塞，转票 45）** — `BuildDiagnosticsBundle` 全仓**无生产调用者**
  （"采样器数据可挂"目前成立在 API 契约层）；`TestDiagnosticsBundleCollectsAndRedacts` 只断
  `slo-snapshot.json` **存在于包内**，未断其字节等于传入值（等值靠 `diagnostics.go:161` 直传保证）。
  票 45 做 UX 时一并补内容回读断言。
- **[A1–A7 票 64 消化标注（2026-09-20，逐条，不许无声消失）]** —
  - **A1 已修（机制）**：`internal/ball/hotkey_reload.go` 桥 + `cmd/balldebug -config` 宿主；
    winlive `TestLiveHotkeyRebindEndToEnd` 证明"真改 `[hotkey]` → `RegisteredHotkeys()` 集合真变 →
    新键真唤起 → 旧键死"。**生产触发者仍缺 → 移交票 12**：`cmd/wisp` 只有 `config.LoadFile`
    （`run.go:189`/`providers.go:90`）、全仓 `ball.New(` 只在 `cmd/balldebug/main.go:171`，
    故装配根需自建 `config.NewManager`+`NewHotkeyReloader` 并把 `Check()` 挂 watchdog tick
    （`[hotkey]` 属 HOT 档，只挂 `OnReload` 永远不会触发）。票 64 的对应 AC 框因此**保持未勾**。
  - **A1b 已修**：`HotkeyStatus` 五分（`Live/Disabled/Unparsable/Taken/Error`）+ `Attempted()` 把
    "被他人占用（真 1409）"与"未尝试"分开，`Problems()` 出用户可见句子；
    实况 `TestLiveHotkeyOccupiedVsNotAttempted` 以真蹲键取真 1409 为证，默认套件
    `TestRegisterAllSplitsFailureFamilies`（fakeRegistry）保住分类不依赖桌面。
  - **A1c 已修（默认值）/文档半待解冻**：`DefaultHotkeys().Summon="Ctrl+Alt+Q"`（R10）+
    `AltSummonSpace="Ctrl+Alt+Space"` 常量注释说明"可配但不作默认（中文输入法抢 Space）"；
    用户可见文档三处（SPEC-03/SPEC-08/PLAN）对本票冻结，未改一字。
  - **A2 已修**：交互四项成 `interaction_live_test.go` 真机行为测试并实跑全绿（单击唤起/拖动不唤起、
    Confirming Esc+单击取消且归还 B1 Esc、`WS_EX_NOACTIVATE`+`MA_NOACTIVATE`+前台未变、透明角按 OS 自己的
    `WindowFromPoint` 判穿透）；**这一轮跑抓出并修掉隐形边框真缺陷**（`2d063f0`：`dwStyle=0`→
    `WS_POPUP`，此前点击落点比球心偏 8x31）。
  - **A3 半修**：句柄/消息级实测已在 winlive（`TestLiveSleepingZeroTimerHandles`：探针先看得到 Warm 的真
    `WM_TIMER`，再断 Sleeping 600ms 内 0 条 + UI 线程只持 1 个球窗），且 fail-fast 不静默 skip；
    **仍待**：winlive 进 nightly CI（`go test ./...` 不覆盖），以及 slo-full 上量 Sleeping CPU≈0。
    D32 阈值未被放宽、Sleeping 未加定时器。
  - **A4 待硬件**：本机物理单屏（`\\.\DISPLAY4` 3440x1440），第二屏实拖**未做**，也没有拿双屏 mock 冒充；
    `enumMonitors()` 仍无注入接缝（属票 64 范围外的接缝工作）。
  - **A5 已修**：live 断言改按测试自身 `DebugHWND()`（`307e24e`），并有前置守卫在发现外部同窗口类球时
    `t.Fatalf` 指名 pid/title（`live_guard_windows_test.go:165`）——A5 那两次假红从此不可能再被读成结论。
  - **A6 已定案：非逐球泄漏**（安静桌面 `-count=3`：base/afterNew/afterClose =
    120/363/369、369/372/370、370/372/370；USER 对象 1→6 后两轮 6→6）。第一颗球的 ~243 句柄是一次性
    D2D/DirectWrite/COM 进程底座，原先"104→376/374"是把"建球前底座"当成"第二颗球的残留"。
    **不立泄漏修复票**；余留观察项：Close 后那 5 个 USER 对象恒定不增长，若哪天变成逐球累加再开票。
  - **A7 维持原判**：仍属票 45（UX）范围，票 64 未触碰诊断包。
- **[A8] 票 63 AC#6：真实 provider 取密钥这条路还没被证明（转票 12）** — `wisp secret` 本身已证明
  ref→`config.LoadFile`→`ProviderKeys`→`Authorization` 头，但发请求的是**测试自己写的 http client**，
  不是 `internal/llm` 的 provider；按 README 规则「安全判定不得由调用方可控的选择器决定」，这一格实现者
  不能自证，故留着未勾。**当前残缺表现**：把 `api_key_ref` 接错、或 provider 侧偷偷回退读明文 config
  字段，票 63 的套件不会变红。**完成判据**已写进票 12 的新 AC 框（判据含变异检验：ref 指向不存在的
  blob 时必须走 Unconfigured 失败路径，不得静默成功）。票 12 因此从"无待消化 AC"变为"背一张外来的框"。
- **[A9] 票 63 MINOR-1：审计轨迹断言被"收窄"开出一个假绿洞（不阻塞）— 已解决 `ef026a4`（2026-09-20）** — `4477f56` 把 `unset` 的
  "forced=true 只属于强制删除"从**整缓冲 grep** 收窄成只看 `auditRecordFor(needle)` 返回的那一条记录。
  收窄修的是真污染（共享日志缓冲后前一个子测试的合法 forced=true 会误伤），但丢了旧断言的一种检测力：
  把 `forced=true` 记进**另一条不含 needle 的记录**、且强制/非强制两条分支都记，
  `go test -run TestSecretUnset ./cmd/wisp/` **0.064s 全绿**。复现补丁存在仓库外
  `ticket63-mutation-M2.patch`（验收方替被杀的验收代理跑完实验后已 `git checkout` 还原生产文件）。
  **当前残缺表现**：弱的是审计轨迹而非泄漏判据——明文入日志仍会红 5 处（已独立验证）。
  **完成判据**：正半边保留 `auditRecordFor`，另加全缓冲**计数**不变式（含 `forced=true` 的记录条数
  == 本测试实际执行的强制删除次数），且 M2 补丁重新 apply 后必须变红。归属：下次动 `cmd/wisp`
  时顺手做，不单独立票。
  - **修法（判据逐条对齐）**：`cmd/wisp/secret_test.go` 新增 `countRecords(logs, needle)`（整缓冲计数，
    与 `auditRecordFor` 并列成一对：前者答"本测试总共说过几次 X"，后者答"这一次操作记了什么"）；
    `TestSecretUnsetRefusesWhileReferenced` 顶部显式 `trail := captureLogs(t)` 取得本测试自己的缓冲，
    并在跑 `unset --force` 的那一处 `forcedDeletes++`，末了断言
    `countRecords(trail, "forced=true") == forcedDeletes`。分母是本测试实际执行的强制删除次数，
    所以共享缓冲不再构成污染（判据第 1 条满足）。既有断言一条未改、未弱化（含 856-862 行的
    record-scoped 半边）。
  - **变异复验**：重新 apply `ticket63-mutation-M2.patch` → `--- FAIL: TestSecretUnsetRefusesWhileReferenced`
    `audit trail carries 4 records saying forced=true, this test performed 1 forced deletes`；
    失败点只有新增的那一条，旧断言全绿，即这正是原洞。`git checkout -- cmd/wisp/secret.go` 还原后
    `-run TestSecretUnset` 0.063s 全绿，全包 `go test -count=2 ./cmd/wisp/` ok 12.586s，
    `go test -race ./cmd/wisp/` ok 8.319s。工作树无残留变异（`git status --porcelain` 已核）。
- **[A10] 票 10 复验登记两处小项（都不阻塞 done，但不得无声消失）** —
  - **A10-a（转票 04）**：`j.startCall(ctx, …)`（`internal/agent/loop.go:585`）跑在活的 task ctx 上，
    取消若落在 executeCalls 的预记账循环里，`InsertToolCall` 失败 → `startCall` 返回 0 →
    `journal.finish` 找不到行 id 就什么都不写 ⇒ **被取消任务的 tool_call 行集静默比模型请求的调用数少，
    且连一条 WARN 都没有**。票 10 已裁定"不在本票修"（本地唯一修法=给每次正常路径的工具记账加 2s 死线，
    是在热路径上新造失败模式去补微秒级窗口），持久解 = 票 04 的 `writeQueue.exec` 在 `BeginTx` 前做
    `context.WithoutCancel`。**附加条件（本次复验新提）**：票 04 落地时必须同时补一条 `startCall` 失败的
    WARN，**只记事实不记内容**（避免踩 C28 日志不落明文）。
    完成判据：一个"任务在记账循环中被取消"的用例，断言 tool_call 行数 == 模型请求的调用数，且失败时有 WARN。
  - **A10-b（下次动 `internal/agent/prompt.go` 时改）— 已解决 `c193ec7`（2026-09-20）**：`enforceTotal` 上方注释写着
    "Every pass must remove at least one token, so the loop terminates from any starting point"。
    **这句作为机制描述是错的**：复验用变异证明，真正的终止保证有**两条且互为冗余**——`nb > toks-1` 钳子，
    以及 `s.Budget = nb` 回写（`nb = Budget - over()` 让 Budget 以 `over()≥1` 严格递减，数轮后必然低于
    `toks`）。删任一条，`TestEnforceTotalTerminatesAtEveryWindow` **仍然全绿**（我实测：拆钳子 ok 0.041s）；
    两条都删才红。后果不是 bug 而是**误导**：后人若信了这句注释，会去删那行真正在兜底的 Budget 回写。
    完成判据：注释改为陈述两条保证各自的充分性，并补一个"删任一条即红"的变异锚定用例（或显式记录该冗余为有意设计）。
    - **实测更正：终止保证是三条，不是两条**。除钳子与回写外，真正兜住小窗口死循环的是**候选过滤**
      （`t <= 1` 的地板段不参与挑选 + `idx == -1` 退出）：把过滤条件改回修复前的 `t == 0`（钳子与回写都在），
      窗口表在 window 1 挂死 → `FAIL … Sections() never returned` 10.035s。旧注释把功劳记错了一处。
    - **实测更正二：本条登记的前提"两条都删才红"不成立**。把钳子与回写**同时**中和，
      `TestEnforceTotalTerminatesAtEveryWindow` 仍 **ok 0.036s**——它压根不复现原始挂死。
      换句话说这两条冗余保证在旧套件下**一点都没被锚定**，比登记里说的更危险。
    - **编排者独立复现（同日，我自己重跑）**：两条变异确实都落进文件后（`:357` 钳子被 `false &&` 短路、
      `:366` 回写换成 `_ = nb`），`TestEnforceTotalTerminatesOnABelowBudgetSection` **FAIL 5.00s** 并给出
      准确诊断，而 `TestEnforceTotalTerminatesAtEveryWindow` 在同一份双删代码下**不在失败列表里 = 仍绿**。
      更正二成立。⚠ 附带一条方法论教训：我第一次跑这个实验时 sed 的第二条**没匹配上**却得到 `ok 0.035s`，
      差点据此判代理说谎——**变异实验必须先 grep 证明变异真的进了文件，再跑测试**；
      已把这条并入 [[adversarial-review-orchestration]]。
    - **做了什么**：注释重写为三条机制（各自的充分性、2/3 互为冗余、只有 1 被窗口表钉住），并在钳子与
      回写两行就地标注 guarantee 2 / guarantee 3 与"keep both"。选项上取 **(i) 的可达部分 + (ii) 的记账**：
      新增 `TestEnforceTotalTerminatesOnABelowBudgetSection`，直接喂给 `enforceTotal` 一个窗口表到不了的
      状态（suffix 段 token 数低于其 granted budget：`Sections()` 会先把每段裁到自己的预算，所以进不去），
      该状态下只有钳子能让首轮就咬、回写要空转几轮 → **锚定"这一对不能同时删"**。
      变异矩阵（`-run TestEnforceTotal`）：删钳子 ok 0.037s、删回写 ok 0.034s、两条都删
      `--- FAIL: TestEnforceTotalTerminatesOnABelowBudgetSection (5.00s)`（窗口表那一条此时仍绿，正是登记的症状）。
    - **为何不写"删任一条即红"的两条单独锚定用例**：能区分二者的可观测量只有*轮数*（需给生产码加计数器）
      与*回写后的 `s.Budget` 字段值*（一段 10 token 的文本挂着 64 的预算并不违反任何契约，现有断言
      `used <= s.Budget` 两种设计都满足）。钉这两个都等于钉实现细节，任何等价的优雅重构都会误红——
      即协调者所说的"钉一条它自己到不了的分支，负价值"。故按完成判据的后一支，把冗余显式记为
      **deliberate defense-in-depth**，并把"删任一条仍绿、两条都删新用例挂死"写进注释，让下一个动手的人
      拿到的是可核对的事实而不是一句错的口号。
    - 门禁：`gofmt -l` 对四个改动文件无输出；`go vet ./cmd/wisp/ ./internal/agent/` 干净；
      `go test -count=2 ./cmd/wisp/ ./internal/agent/` ok（wisp 12.586s / agent 2.477s）；
      `go test -race -count=1` ok（wisp 8.319s / agent 2.886s）。所有变异均以
      `git checkout -- <该一个文件>` 或备份文件还原，`git status --porcelain` 复核无残留。
- **[A11] 票 11 AC#6：探针会测量已证明，但生产装配根还没调它（转票 12，不是转 owner）** —
  `llm.RunProbeSuite`（`internal/llm/probe_health.go`）+ 9 个用例已入库（`43f9a53`/`a355bf6` 同族），
  验收方独立复核结论：**判据仪器没被改过**（`capability_test.go` 在 git 历史里只出现过一次，即它被新增的那次），
  且"探针拿 config 自证"这条我用自己的变异验伪：把 `rep.Flags.Set(..., res.OK)` 改成
  `declaredCapability(o.Declared, c)` → `MeasuresBrokenFC` / `MeasuresBrokenVision` /
  `BrokenOnAllThreeDialects` 三红。**当前残缺表现**：没有任何生产代码调用 `RunProbeSuite`，
  所以"声明 ✓ / 实测 ✗"这条 SPEC-05 §3.1 要求的事件在真机上永远不会发出；
  另外 thinking 那格沿用票 09 的检查，它接受纯文本回答，**检测不出坏思考器**。
  **完成判据（票 12 承接）**：①`cmd/wisp` 的 provider 保存/发现路径调用 `RunProbeSuite`
  （测试里的 `memoryHealthSink` 可直接当适配器抄）；②mockllm 增一个 thinking 能力档，
  thinking 探针要求必须出现 reasoning delta，否则判 broken。
  ⚠ **票 12 现在同时背 A8（票 63 AC#6：真实 provider 取密钥）与 A11（本条）两张外来框**，
  两者都属"装配根没接线"同一形状——票 12 是 S1 端到端门禁，正好一并消化。
  实现者代理把它写成了 `next=user wires…`，**该措辞已作废**：owner 不写代码，此项由票 12 的代理完成。
- **[A12] 票 20 段 1 挖出的系统性事实：C19 风险评审器在**今天之前**没有任何生产调用者（票 17/18 的"done"名不副实）** —
  我用 git 独立证实：在 `c9c3a6f` 的父提交上 `git grep "NewAssessor\|risk\.Assess" -- '*.go' ':!*_test.go'`
  **命中数为 0**；`PathCanonicalizer` / `SensitiveClassifier` 两个冻结缝在票 20 之前**只有接口定义与测试桩**
  （仅出现在 `internal/risk/` 的 4 个文件里，其中 2 个是 `_test.go`）。
  ⇒ **R2（越出 allowlist→L2）与 R3（敏感文件）在生产里从来没有成为过真实判定**，
  而交付它们的票 17、18 早已置 done。代码本身正确、测试也真，缺的是"有人调它"。
  **现状（已改善但未闭环）**：`internal/tools/bridge.go` 现在是全仓第一个 C19 生产调用者
  （`risk.NewRiskAssessor()` @ `:161`/`:529`，真实 `Assess(...)` @ `:239`），
  但**桥自己仍无人调用**——`grep -rn "internal/tools" cmd/` 零命中，`cmd/wisp` 走不到它。
  链路 = 评审器(票17，DONE) → 桥(票20段1，DONE 且零调用者) → **票 12 才通到装配根**。
  **完成判据**：票 12 让 `cmd/wisp` 构造并使用 `tools.New`，并有一条**端到端**用例证明
  "越出 allowlist 的 `fs.read` 真的被判 L2 并走审批支"（今天这条只在包内测试层成立）。
  ⚠ 本条**不重开票 17/18**（它们的 AC 从没声称"已被生产接线"，重开属于事后加要求）；
  它重开的是 **R12 的判断**：这不是两三次偶然，而是本仓 done 的系统性含义缺陷——
  **"done" 只保证包内正确，不保证可达**。今后 done 前必须回答"谁调用我"，
  答不出就在票面与 registry 各留一行，像本条一样显式登记，而不是等下一个代理挖出来。
- **[R13] 票 62 判"临时通过"，质感返工另立票 65；SPEC-08 §2 的改动带 INTERIM 标注**（owner 2026-09-20 10:4x）。
  owner 实机跑 `-tour` 后原话：「**可以说算是赝品吧**，离我发的那种质感还是有不小差距，但是先勉强用吧，
  以后再换样式，就先这样吧。开始阶段不能要求太高，本末倒置了就，先完成核心功能」。
  ⇒ 三点裁定：
  ① **这不是签收。** SPEC-08 §2 只回填了「`Sleeping` 12px/0.35 → 静态玻璃体 44px」这一项，
     且改动理由是可证伪的（旧规格差分取证 = **0 个像素**变化 ≥8/255，物理上不成像），
     **不是**"质感被认可"。批准范围写死为尺寸/可见性/零定时器三项，D32 一条未放宽。
  ② **一条必须承认的事实：图1/图2 从未落到磁盘。** 我逐个核对本会话 11 张 PNG 附件
     （商标查询 4 + 托盘 1 + AnySearch 文档 1 + 代理面板 1 + 本轮 tour 3 + 小裁图 1），
     无一为玻璃参考图 ⇒ 票 62 的**四个代理全部是按文字盲做配色**。
     所以"赝品"这个判决在信息缺失下是**必然的、且对代理公平的**；
     票 65 的第一步因此不是写码，而是**向 owner 取回参考图并入库 `design/refs/`**。
  ③ **优先级重排照 owner 的意思执行**：核心功能优先 ⇒ **票 12（装配）升为当前第一优先**
     （它同时消化 A8/A11/A13 三张同形遗留），票 64（球缺陷）并行，票 65 挂 blocked-on-owner 不占车队名额。
  ⚠ 今后任何读到 SPEC-08 §2 的人：那里的 INTERIM 标注**在 owner 对票 65 明确说"这次对了"之前不得摘除**。
- **[A8 / A11 / A13 → 全部闭环（2026-09-20 12:35Z，票 12 装配 + 编排者票面对账）]** —
  R12 所指出的"能力做完了没人接"这一族，三条同形遗留**一次消化完毕**：
  **A13**（审批门 + 桥不可达）由 `cd011b8` 解决——`cmd/wisp/run.go` 现在同时 import
  `internal/agent/approval` 与 `internal/tools`，这是 `cmd/` 下第一次有代码够到这两个包；
  **A8**（真实 provider 从 DPAPI 取密钥）由 `cmd/wisp/run_test.go::TestRunTextTaskKeyResolvesInTheStore`
  闭环，且我写进完成判据的那条变异有独立用例 `TestMissingBlobFailsUnconfiguredNeverSilently`
  （ref 指向缺失 blob → 走 Unconfigured，**不静默成功**）；
  **A11**（能力探针无生产调用者）由 `cmd/wisp/providers_test.go` 闭环，thinking 探测
  **正反两向都钉**（`...MeasuredThinkingTrue` / `...False`），mockllm 侧新增 thinking 能力档。
  守卫删除与替代**确在同一 commit**（`cbdea7c` 只含 `internal/agent/loop.go` +
  `internal/tools/loop_approval_test.go`），且**正反两向都有测试**：无 gate 时 L1 写被拒、记
  `L1/reject`、文件不出现 ⇒ "忘了接审批"的后果是**拒绝执行**，不是放行。
  ⚠ **本条之所以还要单独写，是因为这三件事全部没有出现在票面上**：做它的代理跑完 165 次工具调用、
  落了 3 个 commit，**AC 框一个没勾、Progress log 一行没写**，`Status` 停在 `ready-for-agent`
  ⇒ 下一张读票的代理会把这 165 次调用**整个重做**。根因是我给票 12 的简报漏了"每次 commit 同步票面"这一条
  （票 11、票 64 的简报里都有）。**教训：验收对账要以代码为准，且"票面与 commit 同步"必须写进每一份简报，
  不能依赖代理自觉。**
  另记两处**段 2 的硬障碍**（代理上报、我已读码确认存在）：
  ①`loop.decideRisk` 在 `Execute` **之前**就拒已声明的 L1/L2 ⇒ `fs.write` 到它手里是死的，
  须等票 21 替换该函数；②`tool_call` 表**没有 `rules_hit` 列**（SPEC-02 的 DDL 是契约层的），
  该字段目前只随 `Decision`/审计行走，代理**没有伪造列**，这是正确处置，票 21 需显式定夺。
- **[A13] 票 21 段 1：审批门与桥都造好了，但装配根还没有它们（R12 的第三次同形复现）** —
  `internal/agent/approval/`（19 个用例）与 `internal/tools/`（桥）今天落地，**两者在 `cmd/wisp` 下都不可达**：
  `internal/tools` 至今**零外部 importer**（我自己 `grep -rn "internal/tools" cmd/` 零命中）。
  验收方独立复核过的部分：`Request.Source` 在非测试代码里**只出现在两条日志**（`claimed_source=%q`），
  决策路径一次都不读；`PanelAPI` 没有 `Allow` 方法、`PanelItem` 没有授权字段（"面板来源结构性拒绝"
  是**类型层面**成立，不是 if 里成立）；我的变异把面板路由改成"也接受调用方递来的 grant"，
  `TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis` **立刻变红**，而原生放行/单次使用两条仍绿。
  实现还多做了一步我没要求的：**未受信路由上一出示 grant，就把该项的活 nonce 全部作废**（当泄漏信号）。
  **当前残缺表现**：真机上 L1 阻断窗、L2 原生卡、批量聚合、300s 自动拒绝**一条都不会触发**，
  因为没人构造 Gate；`fs.write` 依旧在 `loop.decideRisk`（`internal/agent/loop.go:719-738`）被执行前被拒。
  **完成判据（票 12 承接）**：①`cmd/wisp` 组装 `Loop` 时注入 `Options.Gate` 并调用 `AdmitTextTask`；
  ②删掉 `decideRisk` 里那段"L1/L2 先拒"的代码（代理确认有了 Gate 之后这一步**只是删除**）；
  ③一条端到端用例证明"越出 allowlist 的 `fs.read` 真的被判 L2 并走审批支"。
  ⚠ **A8 / A11 / A13 三条同形**，全部收敛到票 12 —— 票 12 的实质是**装配**，不是新功能。
  （另：代理报告称 d22scan 在 `internal/llm/adaptertest/mockllm.go:68` 有一处既有命中，
  我在 HEAD 上复跑 `tools/d22scan` 结果为 **clean**，该命中不存在，无需处理。）

## 已解决（resolved）

- **[H4] `web.search` 实现路径** → 2026-09-20 用户裁定：**搜索 API**，供应商 **anysearch**（用户已持有 key）。
  票 22 的 `SearchProvider` 以 anysearch 兼容实现为默认 provider，抓结果页降为离线/被封兜底。
  key 录入走 R7（用户自助，不经对话框）。
- **[H5] 默认 TTS 选型** → 2026-09-20 用户裁定：**云 TTS，供应商 StepFun**（用户已有 StepFun 云 TTS 模型）。
  票 26 的 TTS 输出层以 StepFun 为默认实现，与票 61（C9 云级联 ASR/TTS）同源。
- **[H1 形式变更] 悬浮球签收** → 2026-09-20 由"看截图"改为"实况签收"：`build\balldebug.exe -stay`
  启真球到桌面。首次实况启动记录：`handles=381`（与 SLO §2 实测 407–413 同量级）；
  ⚠ **热键 `Ctrl+Alt+W` 注册失败**（`hotkey registration failed (already taken?)`）——
  本机已有他者占用该全局热键，签收时需一并确认是否换键（属 C-契约外的默认值调整，需用户点头）。
