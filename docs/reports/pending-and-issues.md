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
- **[A6] 句柄增长 → 已定案结案（2026-09-20，票 64 接续代理 + 编排者独立复跑）** —
  原记"同进程 `-count=2` 两次观测到基线由 104 升至 376/374，`Ball.Close()` 后约 270 个句柄不回落，
  疑为真泄漏，本次未能定案"。**定案：不是每球泄漏。**
  票 64 的代理把句柄数改为**每轮打点**（`live_windows_test.go:79`），我自己在安静桌面上
  独立跑 `go test -count=3 -tags winlive -run TestBallLiveLifecycle ./internal/ball/`，三轮数据：
  **base=114 → afterNew=363 → afterClose=369（净 +255，USER 对象 1→6）；
  base=369 → afterNew=372 → afterClose=370（净 +1，USER 6→6）；
  base=370 → afterNew=372 → afterClose=372（净 +2，USER 6→6）** ⇒
  那 255 是**进程内第一次建球的 D2D/COM 初始化地板，只发生一次**；此后每轮建/销净增 1–2 个句柄、
  USER 对象稳定持平。**不开泄漏票。**
  ⚠ 两点如实保留：①每轮净增是 **+1/+2 而非 0**，量级远低于当初担心的 270/球，
  但长驻进程反复建/销球时是否为 COM 缓存还是慢泄漏，**本轮数据不足以下结论**（要几百轮才看得出斜率），
  若日后票 08 的 SLO 长跑出现句柄单调爬升再回指本条；
  ②`Ball.Close()`（`ball_windows.go:762`）里 `pDestroyWindow.Call` 的**返回值仍未检查**——
  这次没被证明有害，但它正是"销毁失败会静默"的那类代码，属可选加固，不单独开票。
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
  ⚠⚠ **本条目末尾这句"该命中不存在"是错的，撤回（2026-09-20 21:38，见 A22）**：
  我当时在**仓根**跑 `go run ./tools/d22scan -root .`，它打印的是
  `main module does not contain package .../tools/d22scan`（**扫描器从未执行**），我把"没有命中输出"读成了"clean"。
  正确调用（`cd tools/d22scan && go run . -root ../..`，**与 CI 那一步逐字同形**）实跑结果：
  **该命中真实存在**，`[bare-goroutine] bare go func( is banned (D22/D38b)`，由 `78b1466`（票 11）引入、不在 allowlist。
  ⇒ **那位代理的报告是对的，我当众否定了一次如实报告**；根因与纪律修正记在 A22
  （判据仪器"没报错"≠"跑过了"，采信前必须制造一次已知会红的阳性）。**票 20 的 AC#1-4 结论不受影响。**

- **[A14] `parseSystemProcesses` 丢快照末项 ⇒ D32 的 state 口径 SLO 门从未产出过样本窗，且 CI 两处门坏着** —
  `internal/proc/treemetrics_windows.go:240-243` 先 `if next == 0 { break }` 再解析当前项（:247-260），
  于是**系统快照链上最后一个进程永远进不了 map**；新建进程恰在 `ActiveProcessLinks` 末尾 ⇒
  「自己量自己」的 `wisp slo -state X` 常态 fail-closed `exit=2`。
  **我自己独立复现（2026-09-20 21:44，同一二进制）**：无 keeper → `exit=2` 报
  `system snapshot does not contain the sampling process`；加 keeper → 采样跑通。兄弟实现
  `cmd/balldebug/diff_windows.go::privateWorkingSetFor` 顺序是对的（先比 pid 再 break）。
  **影响面（三处，都在已 done 的东西上）**：①`build/slo/slo-report.json`（09-19T23:35Z，票 08 归档跑）
  `Sleeping exit=2 / Warm exit=2 / all_pass=false` ⇒ **票 08 的 state 口径九个样本窗里一个都没有**；
  ②`.github/workflows/ci.yml:166`（slo-smoke）与 `:196`（slo-full）**无 `continue-on-error`** ⇒ 这两处今天按进程创建顺序随机红；
  ③票 12 AC#2 的六份 A.1 样本是靠 keeper 绕法才拿到的（口径本身带偏，见 A15）。
  **完成判据 + 当前残缺表现登记在票 66**；修好前 `slo-check.ps1` **不得**被当成已通过的门引用。
  ✅ **已解决（票 66，agent-ticket66b 复测 2026-09-20T15:06Z）**：解析顺序修复 = commit **`86e868d`**
  （全仓只留一份走链实现 `WalkSystemProcesses`；回归用例经**第二人独立变异复测**咬住：退回旧序 ⇒ 6 红/24 绿，
  原始输出 `docs/evidence/s1/66-mutation-parse-order.md`）。**零 keeper** 连跑 5 次 `-state Sleeping -seconds 10
  -interval-ms 250` ⇒ `exit=0/0/0/0/0`、九门各 9 行（`docs/SLO.md` 附录 C.1，JSON `docs/evidence/s1/66/66b-nokeeper-{1..5}.json`）。
  A.6#1 的 keeper 绕法已从文档删除。⚠ 本条修复**不足以**让 CI 变绿：它反而揭出了下面 C.4 那起聚合崩溃
  （修好 A14 之后 `slo-check.ps1` 是「必红且无报告」，不是 B.3 预告的「随机红」）。**finding 原文照录不删。**

- **[A15] D32 `Sleeping` CPU 行在树内口径下**没有定义**：观测者单次读的绝对成本就等于门限本身** —
  采样器跑在被测进程内，每次 `ReadTree` 的 CPU 被计进被测数字。**我复跑 8 个样本（票 12 那次 6 个，合计 14 个）的结论**：
  ①同一绝对成本在两种间隔下都量到 —— 250ms 间隔跳动样本 0.5062–0.5276%（≈**1.28–1.32ms/次读**）、
  2000ms 间隔 0.0649–0.0652%（≈**1.298–1.304ms/次读**）；
  ②⇒ **在 250ms 间隔下「采样一次」本身就是 0.52% > 0.5% 门**，门与仪器最小可分辨量撞在同一数量级；
  ③同配置重复跑 PASS/FAIL 随机翻转（我 5 个 10s/250ms 样本：0.0387% / 0.5838% / 0.5837% / 0.5053% / 0.2076%）
  ⇒ **只修 A14 会把 CI 从「必红」变成「约一半概率红」**，比确定性地坏更糟；
  ④产品侧清白：`iv250` 40 样本里 **37 个 CPU 恰为 0.0000**，与 A.2 树外实测球体 `Sleeping` **0.000%** 互证。
  ⚠ 同时更正票 12 记录里的**解释形式**：A.1/A.4 的「CPU ∝ 读次数（118/60s=0.629%…每次读折 16–38ms）」
  **不能作为已证事实引用**——我这次 40 次读反而比 5 次读更省（0.0387% vs 0.1032%），keeper 5→20 也无单向变化。
  原始表数字不改写（那是当次真测），但解释以 **`docs/SLO.md` 附录 B** 为准。
  **裁定**：票 12 AC#2 保持未勾，CPU 行判「仪器未定义」（既不 PASS 也不 FAIL），产品证据以 A.2 树外 0.000% 为准；
  **D32 的 0.5% 阈值一字未动**。修法（把观测者移出树 / 或书面裁定「从门集里移除 CPU 行」）归**票 66**。
  ✅ **已解决（票 66，commit `00bbb76`）**：CPU 行改成**树外口径**（父读子 pid，`proc.ExternalSampler`；
  阈值与 limit 字符串两侧仍由同一个 `stateCPULimit` 生成 ⇒ 「改的是量法不是门」机器可核，
  钉住它的是 `internal/observe/observer_cost_test.go`）。完成判据的双条件同时成立：同一 30s 窗内
  **树外 0.0173% ≤0.5%** 而 **树内 0.7629% >0.5%**（`docs/SLO.md` 附录 C.3，JSON
  `docs/evidence/s1/66/66-ac3-30s.json`）；五个 10s 零 keeper 跑里树内 0.5963–0.9318% 对树外 0.0000–0.0130%。
  树内那条仍写进 JSON 并标 `observer_cost: true` + `gate: false`，pass/fail 照判。
  **AC#4 的 5/5 exit=0**（`docs/SLO.md` C.4）依赖另一处修复：`scripts/slo-check.ps1:129` 在
  `Set-StrictMode 2.0` 的**全过路径**上抛 `.Count` 异常，commit `cf0c050`。**finding 原文照录不删。**

- **[A16] 球热键的配置值与注册 id 靠**下标 zip** 对齐，无一用例看管（票 64 验收时读码挖出，与 C-3/M-7 同族）** —
  `internal/ball/hotkey_windows.go:382-386`：
  `bindings := []string{cfg.Summon, cfg.Mute, cfg.Cancel, cfg.Panel}` 与 `hkNames`（`{hkSummon,"summon"}`…）
  在 `for i, p := range hkNames` 里靠 `bindings[i]` 配对。`grep hkNames --include=*_test.go` **零命中**
  ⇒ 没有任何断言钉住这个对应。两个方向后果不同：给 `hkNames` 加第五项而忘扩字面量 → 越界 **panic（响的）**；
  但**调换 `hkNames` 顺序** → `cfg.Summon` 的字符串以 `hkMute` 的 id 注册，`Problems()` 依旧报 "summon live"，
  按下唤起键去切了静音——**静默错位**。**根因同族**：M-7 是"按键名决定扫不扫"，C-3 是"只看第一个命中的路径键"，
  A16 是"判定建立在两个平行数组的下标恰好对齐上"。
  **修法方向已定**：合成**一张具名单表**（`{id, name, cfgField}`）消掉 zip；
  ⚠ **不接受**"保留 zip + 加一条顺序断言"这种修法——断言挡不住下一个人只改一侧。
  **完成判据**：① mutation 检验先行——把 `hkNames` 前两项对调，若默认套件与 winlive **全绿**，即证明今天无人看管，
  该记录留在 `docs/evidence/s1/64-*`；②改成具名单表后**重跑同一变异**，必须有测试转红；
  ③四项 `Problems()` 句子与 `RegisteredHotkeys()` 集合逐 id 对得上（正向断言，不只断"不 panic"）。
  - **✅ 判据① 已于 2026-09-21 08:38 执行，结论是"有人看管"那一支**（编排者亲自跑，桌面与测试窗当时空着：
    `tasklist` 无 `go.exe`）。变异**先证落地**（`grep -n MUTATION-A16` 出 `:35`/`:36` 两行）再跑
    `go test -tags winlive -count=1 -v ./internal/ball/` ⇒ **`MUTATED_EXIT=1`，6 条转红**：
    `TestLiveHotkeyRebindEndToEnd`、`TestLiveHotkeyOccupiedVsNotAttempted`、`TestLiveMuteHotkeyEndToEnd`、
    `TestRegisterAllLiveSet`、`TestRegisterAllSplitsFailureFamilies`、`TestHotkeyReloaderRebindsOnConfigChange`
    （64 条 `=== RUN` / 51 PASS / 6 FAIL；原始日志 `docs/evidence/s1/64-winline-retest-2026-09-21/a16-mutation-red.log`）。
    还原后复跑 `ok 12.189s`、`grep -c MUTATION-A16`=0、`git status --short internal/ball/` 空。
    ⇒ **我 08:35 那条自我更正（"保护是副作用而非意图"）成立**：`hkNames` 顺序被动确实会被咬，
    咬点是"四槽状态各不相同 ⇒ 换序会把状态搬到错误的 id"这一族断言。**静默错位今天不存在**。
  - **但这条不该被当作"已解决"**：挡住换序的是**巧合的形状**（四槽状态恰好互不相同），不是**意图**。
    真会变静默的时刻是**加第五个热键**那天（第 5 项与第 5 个配置字段谁先忘都不一定有测试红）
    ⇒ **判据②③ 与修法（具名单表）绑到票 39 的"第五热键"那一批**，本条保持 open 但**优先级下调**，
    且**禁止**用"保留 zip + 补一条顺序断言"交差（原登记已写明理由：断言挡不住下一个人只改一侧）。

  **当前残缺表现**：今天顺序恰好对，**所以没有任何用户可见故障**；这是纯潜在缺陷，
  但在加第五个热键（票 39 配置 GUI 必然要加）之前不修，就会变成"改一处顺序、四个键悄悄串位"。
  **归属**：票 64 已 `review` 且本票声明不碰球码 ⇒ **归票 21 段 2 的同批球面改动，或独立小票**；
  编排者在桌面/测试窗空出后先做 ① 的变异检验再定。
  ⚠⚠ **登记后 25 分钟内自我更正（2026-09-20 21:25，把上面"无一用例钉住"收回）**：
  我按 `grep hkNames --include=*_test.go` 零命中就下结论"没人看管"，**这条推断是错的**——
  字面命中为零属实，但保护**藏在别的断言里**。逐条读测试体后重算 6 种两两对调，**全部至少被一条现有断言抓住**：
  - `TestRegisterAllSplitsFailureFamilies`（`hotkey_status_test.go:89-126`）四槽状态**各不相同**
    （summon=Live / mute=Taken(真 1409) / cancel=Live / panel=Unparsable），
    于是凡涉及 **mute 或 panel** 的对调都会把状态搬错 id ⇒ :98 `IsLive` 断言或 :103/:111 状态断言转红；
    且 :120-126 只配 Summon 却断言 `hkMute/hkCancel/hkPanel` **全不 live**，这条正好抓住 summon↔mute。
  - **唯一两个同为 Live 的槽（summon↔cancel）**由另一条抓住：
    `hotkey_status_test.go:225` 与 `:241` 断言 `rep.Binding(hkSummon).Acc.VK == 'R' / 'Q'`，
    即**配置里的 summon 串必须落在 `hkSummon` 这个 id 上**——这正是配对断言本身。
  **所以本条的严重度从"静默错位的潜在缺陷"降为「保护是副作用，不是意图」**：
  配对靠三条与配对无关的断言**偶然**兜住，改测试的人不会意识到自己在拆配对保护；
  而"加第五个热键"（票 39 配置 GUI 必然要做）若五槽状态相似、又没写新的 `.Acc` 断言，
  错位就会**真的变静默**。修法方向不变（具名单表消掉 zip，不接受"保留 zip + 加断言"），
  但**优先级从"尽快"改为"与票 39 第五键同批"**——因为那才是它变成真缺陷的时刻。
  **变异检验的预期结果随之反转**：对调 `hkNames` 前两项**应当转红**（而不是我原先预言的"全绿"）。
  跑红 ⇒ 证明保护是真的，本条按"文档性 MINOR"收尾；跑绿 ⇒ 我这次更正又错了一层，升级回原严重度。
  **教训同族**：与 [[adversarial-review-orchestration]] 第 5 条同形——**"grep 符号名零命中"≠"没有保护"**，
  断言可能通过别的变量名咬住同一件事；下"无人看管"的结论前必须读测试体，不能只读测试标题与 grep 命中表。

- **[A17] 票 20：`Needs` 的逐条断言与它判的代码同 commit 被删（MAJOR）→ 已修** —
  `0986d63` 把 `TestFSRegistrationIsTheL0Pair` 改名成 `TestFSRegistrationIsTheD34Roster` 时
  删掉 `len(e.Decl.Needs)!=1 || Needs[0]!=CapFSRead`，只留下 `fs.delete` 一条。
  **为何有牙**：`bridge.go:276-296` 的 C3 判定 `declared.missing(need)` / `authz.missing(need)`
  取的就是 `Needs` ⇒ 少报一个能力，就等于**只被授权 `fs.read` 的机器放行一个真写盘的调用**，
  且 `Decision.Capabilities` 记的也是错值；`registry.go:140` 的 `Capabilities ⊇ Needs` 只挡多报不挡少报。
  **公道话**：老断言写死"五槽全 `[fs.read]`"，新增 write/trash/move 后**必须**改；
  错的是**换成"什么都没有"**而不是换成逐条表（删守卫须与替代同批，这次替代缺席）。
  **已修**：`fs_test.go` 补 `wantNeeds` 逐条表；**变异检验**（`fs_write.go:675` 改 `CapFSRead`）→
  仅新断言转红、**其余全套件仍 `ok 20.600s` 全绿** ⇒ 量化了"旧套件对这个形状完全盲"。
  证据 `docs/evidence/s1/20-needs-assertion-restored.md`。**当前无可利用缺陷**（今天六条声明都正确）。
- **[A18] 票 20 AC#3：真·外部 kill 下的暂存文件残留未证** —
  `TestAtomicWriteKillsMidWrite` 用的是**进程内** `Hooks.Kill`（`fs_write.go:47-53` 返回 error，**会**跑清理），
  而真 `taskkill` 不跑 Go 的清理 ⇒ `.wisp-tmp-*`（`fs_write.go:38`）会留在用户目录，
  测试里"不留暂存文件"那条断言（`fs_write_test.go:120-128`）在真实故障下**不可观测**。
  ⚠ 这条是审计代理的判断，**我没复现 ⇒ 票 20 的 AC#3 框不动**（纪律：不复现不改勾），改为登记。
  **完成判据**：①一条用真子进程 + `taskkill /F` 的用例，断言**目标文件要么完整要么不存在**（这条今天已结构性成立：
  `fs_write.go:276` 暂存建在目标目录、`:322` 单次 `os.Rename` 落地）；②再断言遗留 `.wisp-tmp-*` 的**下次启动自愈清理**
  或明确申报为可接受残留并写进 SPEC。**归属**：票 20 收尾批；清理器若归票 39 需书面转办。
- **[A19] 票 21 BLOCKER：审批层有判定、没有输入设备（今天 L1 必无否决执行、L2 必 300s 自动拒）** —
  ⚠ **定级已按 R14 更正为「S1 可接受、S3 必做」**（PLAN `:1424` 把「门控 UI」排除在 S1 之外）；
  **发现内容本身不变**，安全约束是"S1 不得依赖 L2 达成出口判据"。**下面的 grep 与后果照录。**
  我自己的 grep：`DecideFromNative`/`DecideFromPanel`/`.Native()`/`.Veto(` 在非测试码里**只有定义与一句注释**；
  `cmd/wisp/run.go:257` 用 `approval.NewChannels()`（无参⇒零通道已加载），
  而 `approval.go:163` 那个载入 Ball+Esc 的构造函数**没人调**，`SetLoaded`（`:178`）零非测试调用者。
  **当前残缺表现**：真机上 L1 阻止窗一定到时执行（人点不了、Esc 不接、语音不接），L2 一定自动拒；
  `consoleApprovalUI`（`run.go:500-538`）从不应答。
  **完成判据（票 21 段 2）**：①`cmd/` 或段 2 的 winlive 里存在一条调用链真正抵达
  `gate.Native().Allow`/`DecideFromNative`，且 `Prompt.Grant` 取自**被显示的那个面**；
  ②一条 cmd 级用例断言 L2 请求以 **allow 结束**而非 timeout；③通道加载改走带参构造，
  且"未加载 ⇒ 提示「语音取消不可用」"与"已加载 ⇒ 否决生效"两向都有测。
  **⚠ 段 2 的验收重点不是美观，是「让一个真人能改变一个判定」。**
- **[A20] 票 21 BLOCKER：D45-1 批量聚合站错了轴，且对全部在产工具恒为死码** —
  ⚠ **归属按 R14 更正**：PLAN `:3105` 把「批量聚合确认（D45-1）」落在 **S3 本片** ⇒
  真修法归 S3；**票 21 段 2 只负责把这段死码删除或显式隔离**。
  下面"恒返回 nil / 测试喂的是产不出的形状"两条事实**不变且必须留在票面**，
  否则 S3 会以为它已经能用。
  `batch.go:35-38` 的 `Aggregate(d tools.Decision)` 要求 `len(d.Paths) >= 3`，判的是**单次调用内的路径数**；
  契约原文（`docs/PLAN.md:1590`「批量聚合（500 个 L1 → 一次确认）」、`:1795`「批量场景 = 500 次 L1 确认」）
  说的是**跨多次调用**。可达性我也核了：`fs.write`/`fs.trash`/`fs.delete` 各 1 个 `path`
  （`fs_write.go:223-225/349-350/549`，且 `additionalProperties:false`），`fs.move` 2（`from`/`to`，`:415-417`）
  ⇒ **任何生产工具单次最多 2 条路径，`Aggregate` 恒返回 nil**。
  **当前残缺表现**：确认疲劳（B2）这一整类风险今天**没有任何缓解**，而包内 5 条 `batch_test.go` 全绿——
  它们喂的是生产里产不出来的工具形状。
  **完成判据**：①按契约语义改成**跨调用**聚合（窗口内合并同源 L1 确认），或书面申报 D45-1 延后并由我改 PLAN 的
  S3 done 判据引用（改契约需 owner，D22）；②若保留 per-call 聚合，必须有一个**真在产**的多路径工具作载体，
  并加一条桥级用例证明"该形状真会被模型产出"；③`batch_test.go` 改名或加注，说明它测的是哪种形状。
  **归属**：PLAN `:3105` 把 D45-1 落在 **S3** ⇒ 票 21 段 2 **只登记不顺手改**。
- **[A21] 票 21 两条 MINOR：`corr` 可猜 + 污点归属取 `Paths[0]`** —
  `bridge.go:219` 的 corr 来自 `loop.go:640`（= taskID），`revokeGrants`/`reject`/`Veto` 直接吃它
  ⇒ 猜中者可使诚实卡片的活 nonce 作废（**fail-closed DoS，非提权**：伪造放行仍被 bind 摘要挡住）。
  `bridge.go:501-503` `origin = dec.Paths[0]` 是"取第一个"选择器，今天确定只因 `pathArgs` 先排序
  ——**同族第四次**（M-7 / C-3 / A16）。**完成判据**：前者断言 corr 对模型不可得或按条目哈希；
  后者在污点归属处**遍历全集**（任一路径带污点即整体带污点），并加一条"乱序输入 ⇒ 同一结论"的用例。
- **[A22] ⚠ 我（编排者）的一条假绿：`d22scan` 是独立 Go module，我从仓根调用它**根本没跑起来**，我把"无输出"读成"clean"** —
  **事实链（我今天自己复跑双方向确认）**：
  ① 我在仓根跑 `go run ./tools/d22scan -root .` → 输出是
  `main module (github.com/CarlosShao/wisp) does not contain package .../tools/d22scan`
  ——**扫描器从未执行**，我却据此写下"HEAD 上 d22scan clean"。
  ② 正确调用 `cd tools/d22scan && go run . -root ../..` → **1 条真命中**：
  `internal/llm/adaptertest/mockllm.go:68: [bare-goroutine] bare \`go func(\` is banned (D22/D38b)`，
  由 **`78b1466`（票 11）** 引入，`tools/d22scan/allowlist.txt` 里**没有**它。
  ③ CI 的调用方式恰好是正确的那种（`.github/workflows/ci.yml` "D22 seven-ban + emoji scan" 步：
  `cd tools/d22scan; go run . -root "${{ github.workspace }}"`，**无 `continue-on-error`**）
  ⇒ **lint job 自 `78b1466`（约 05:59Z）起就是红的**。
  **我造成的二次伤害**：票 21 的代理早前**如实报告过同一处命中**，我用自己的假绿**当众否定了一次真实报告**
  （话留在本文件 A13 条目末尾与票 12 的 12:35Z log 行里）。**已在两处原文下追加更正，不覆盖。**
  **元教训（进偏好，与 [[adversarial-review-orchestration]] 第 2/7 条同族）**：
  **判据仪器"没报错"≠"跑过了"**；凡是外部扫描器，采信前必须至少制造**一次已知会红的阳性**（seeded violation）
  或核对它的非零退出码与"我用的调用方式与 CI 里那一行逐字相同"。d22scan 自带 seeded-violation 测试
  （CI 里那步的上一行 `go test ./...` 就是它），**我一次都没跑**。
  **完成判据**：①`mockllm.go:68` 改走 `observe.Registry.Spawn`（命名 + owner + recover）
  或**由我书面**加入 allowlist 并写明理由——**不许静默放宽**；②跑 `cd tools/d22scan && go test ./... && go run . -root ../..`
  两条都过；③本条与 A13 末尾、票 12 log 的更正三处互相引用可追溯。
  **归属**：票 11 的包（`internal/llm/adaptertest`）⇒ 开**票 67** 承接，连带票 12 的三处 UI 字形（下条）。
- **[A23] 票 12 自己带进了三个用户可见 emoji，而 emoji 门今天根本不看 Go 源码** —
  ①**门是瞎的**：`tools/d22scan/main.go:71` 的 `emojiRe` 只被 `walkEmoji` 用在 `design/` 与 `frontend/`
  （`main.go:137-142`），而 **`frontend/` 在本 HEAD 不存在** ⇒ 该门今天的实际覆盖面只有 `design/`，
  **对 `internal/`+`cmd/` 的 Go 字符串字面量完全不可见**。⇒ 票 12 AC#7 那句"zero emoji scan over new UI strings"
  **前提是半假的**，不是"做完了"。
  ②**真命中 3 处，且是票 12 自己的 commit 带进来的**：`cmd/wisp/providers.go:201` `verdict := U+2717`、
  `:203` `verdict = U+2713`、`:209` 打到 stdout，另 `:41` 出现在 help 文本里；U+2713/U+2717 落在 ban #8 的
  `2600-27BF` 区间内。引入 commit：**`cd011b8`（票 12 的装配 commit）**。
  **我没有删它们**（改用户可见文案属 UI 变更，要签收）。**完成判据**：①把 ban #8 扩到 `internal/`+`cmd/`
  （代码改动；先决条件是 `mockllm.go:68` 清零，否则门一直红）或改写 AC#7 措辞为"design/"——**二者都要我书面裁定**；
  ②`providers.go` 三处字形改成 `PASS`/`FAIL` 文本（推荐，Windows 控制台字体不保证有这两个码位）
  或书面豁免；③AC#7 的两个半边各按上面的真实覆盖面重判。 **归属**：票 67。
- **[A24] 票 12 AC#7 复验带出的 C21 token 表六处争议（D1–D6）：签收尺寸四方不一致 + 默认构建画的不是签收的样子** —
  本条是**登记收敛**（票 68 的成立理由与票 12 的 AC#7 都指回这里），逐条给去向。
  - **D1（=票 68 AC#1）`Sleeping` 尺寸四方对不上**：SPEC-08 §2 INTERIM 写 **44px**；
    `internal/ball/statevisual.go::stateSize` 默认 56 基准 × `SleepRestRatio=0.62` = **34.72px**（下限 `SleepingRestMinPx=30`）；
    token 表旧行写 **12px**；而 `docs/SLO.md` §A.2 实测差分化像框 **46×46 / 2103px≥8/255**
    ⇒ 被量的那个跑法里 Sleeping 确实是 44px 级。**完成判据**：三列真值表（flag off / on 未吸附 / on 已吸附），
    与契约不符的格子如实报不符。**SPEC-08 冻结：不许改文本凑数。**
  - **D2（=票 68 AC#2）owner 签收的球不在默认构建里**：`statevisual.go:80` `var prototypeVisuals bool` 默认 **false**，
    全仓唯一开启者是 `cmd/balldebug/main.go:104`（`EnablePrototypeVisuals(!*frozen)`）；
    被门控的是 `dock_windows.go:90/114/176/206/229/245` + `liquid.go:313` + `liquid_windows.go:52/154`
    ⇒ **默认库画 12px 微点**，我今天改的 SPEC-08 §2 描述的是**非默认配置**（这笔账是编排者欠的）。
    **裁定**：翻转默认值使"默认构建 == 被签收的样子"，`-frozen` 留作对照逃生门；
    ⚠ **迁移测试不是删测试**——钉旧冻结规格的那批要显式 `EnablePrototypeVisuals(false)` **保留**。
  - **D3（已裁定并落文档）**：36 个 look 色**在 `tokens.css` 里无出处**（原生首创），而表头原文写"权威真相源 = `tokens.css`"。
    已在 `docs/evidence/s1/c21-native-tokens.md` 表头加「覆盖面限定」块，把两类行的约束来源分开写死。**不改码**（设计面，等票 65/68）。
  - **D4（新缺口，待认领）**：AC#7 代理补进表里的 **61 行（25 几何动效 + 36 look）无任何机器检查**
    （`TestTokenGoldenValues` 只管 20 条配色）。**完成判据**：纳入某条机器断言（golden 或表↔码反向 grep），
    否则它们只是人写的一篇作文。**归属：票 69**（`69-c21-token-table-machine-check`，
    **被票 68 阻塞**——同包并行是假并行，两票必须串行）。
  - **D5（已修，根因在我）**：票 12 票面出现"截断的重复 log 行"，是我用 `printf` 追加时 `%` 被当格式符、写了一半。
    **教训：往文件写含 `%` 的长行一律用 Edit，不用 printf/echo。**
  - **D6（范围裁定）**：`tokens.css` 129 条声明里 **80 条是面板/屏幕侧、原生侧无对应物** ⇒ 已在表内以「范围」块
    **按名声明为域外**，而非补 80 行空 Go 列；重导命令在 `审计` 段。

- **[A25] S1 四项出口判据里，「回复后 3s 回落」这条**今天完全没有被量过**（建 S1 清单时才发现）** —
  PLAN `:1424` 的 S1 行第一项原文是「打字→得回复→**3s 回落**」。逐项对账的结果：
  文字通路端到端已证（票 12 AC#1，我复跑 `ok 13.794s`）、≤25MB 有实测（`docs/SLO.md` §A.1/§A.2）、
  `wisp run` 可用、C21 落地（检查有洞＝A24-D4）——**唯独"3s 回落"找不到任何一条实测记录**。
  `wisp slo` 量的是**空闲态**，`balldebug -diff` 量的是**逐态驻留**，两者都不量"最后一次回复之后多久进 `Settling`/`Sleeping`"。
  **这不是"大概没问题"，是本阶段唯一一条零证据的验收判据。**
  **完成判据**：①一条用例/实测记录，从"最后一个 token 送达"起算到"进入 `Settling`"，
  **用单调钟**（`observe.Timeout` 那条路，禁墙钟差，D22 禁令 #4），断言 ≤3s；
  ②多样本全报（≥5 次），**不许用平均抹掉坏尾部**，超线就报 FAIL；
  ③若计时起点需要事件钩子（当前 `Dispatch` 收尾没有暴露"最后 token 时间"），**这一步是代码改动**，
  归入票 12 的装配半或票 28（session scope/warm）而不是塞给在跑的代理。
  **归属已定（21:02，我核过票 28 之后）**：**塞不进票 28**——它的 AC 是
  「90s Warm→Settling / 30s Conversation→Warm」，与"回复后 3s 回落"**不是同一件事**；
  ⇒ 已作为**新的一条 AC 直接加在票 12（S1 gate 票）**上，判据写死：单调钟计时（禁墙钟差，D22 #4）、
  ≥5 样本全报、超线录 FAIL、若需事件钩子则明写"那是代码改动，写进本票而不是用肉眼数秒表绕过去"。
  **理由**：它是 S1 的出口判据，出口判据的测量责任必须落在 gate 票上，否则又变成"人人以为别人测了"。

## S1 出口清单（2026-09-20 21:54，编排者按 PLAN 原文逐条对账）

判据来源：`docs/PLAN.md:1424` 的 S1 行**原文四项**，不是我自己发明的清单。
⚠ 同一行还**显式排除**了「语音全部、WebView 面板、命令面板、**门控 UI**、插件抽象、记忆」——
这条排除**改变了我刚给两条发现的定级**，见下面「裁定 R14」。

| S1 判据（PLAN 原文） | 今天状态 | 证据 / 缺口 |
|---|---|---|
| 打字→得回复→**3s 回落** | **半** | 文字通路端到端已证（票 12 AC#1 `TestRunTextTaskTextPathEndToEnd`，我复跑 `ok 13.794s`）。⚠ **"3s 回落"我没找到任何一条实测记录**——`wisp slo` 测的是空闲态，不测"回复后多久回 Sleeping"。**这是 S1 四项里今天唯一完全没被量过的一条。** |
| 空闲 private RSS **≤25MB 实测通过** | **达成（口径受限）** | `docs/SLO.md` §A.1 树内 8.43–9.10MB、§A.2 带球树外 11.79MB，均 ≤25MB。⚠ 口径边界照录：**S1 唯一带球宿主是 `balldebug`**，`wisp run` 按 D12 无球。CPU 行按 A15 判「仪器未定义」，票 66 修。 |
| `wisp run "..."` CLI 可用 | **达成** | 票 12 AC#1/#4/#5/#6 已勾（含 key 走 store、探针被装配根调用），A8/A11/A13 三条同形遗留全部闭环。 |
| **C21 设计 token 落地**（原生侧） | **达成，但检查有洞** | 表 227 行、78 配色行与 `tokens.css`/`tokens.go` 逐值零漂移（票 12 AC#7 代理的对账）。⚠ **A24-D4**：新补 61 行无机器检查 ⇒ **票 69**（被票 68 阻塞）。AC#7 框仍不勾，等票 67 AC#3 + 票 68 AC#1。 |

### 裁定 R14（2026-09-20 21:54）：A19/A20 的定级按 PLAN 的阶段归属更正，**不按"能力应当存在"的直觉**

- **A19 从 BLOCKER 降为「S1 可接受、S3 必做」**。理由：PLAN `:1424` 把「**门控 UI**」明确排除在 S1 之外，
  所以"没有生产应答路径"（`DecideFromNative`/`Native()`/`Veto()` 零非测试调用者）
  **在 S1 是符合计划的**，不是违约。**但安全后果必须写死**：无应答路径 ⇒ **每一次 L2 都必然 300s 自动拒绝**，
  因此 **S1 不得依赖任何 L2 工具达成其出口判据**（今天满足：S1 四项都不需要 L2）。
  归 **票 21 段 2**，与它原本的 native 卡片渲染同批——**那才是它的验收阶段**。
- **A20 维持"站错轴 + 恒为死码"的事实认定，但归属从"票 21 段 2"改为「S3 本片」**，
  因为 PLAN `:3105` 把「批量聚合确认（D45-1）」明确落在 **S3**。
  ⇒ 票 21 段 2 **只负责删除或显式隔离这段死码**，真正的跨调用聚合归 S3 的票。
  ⚠ 不撤销的发现内容：`Aggregate` 对全部在产工具**恒返回 nil**、而包内 5 条测试全绿喂的是产不出的形状——
  **这条要留在票面上**，否则 S3 会以为它已经能用。
- **纪律沉淀**：我给发现定级时**必须先查 PLAN 的阶段归属**，再查"能力是否可达"。
  只查后一条会把**计划内的延后**报成缺陷，连着几次之后我会开始低估真缺陷——那比漏报更贵。

### S1 之后按依赖顺序的下一步（不是待办清单，是排程）

1. **票 66**（测量独占，进行中）→ 解锁 `slo-check.ps1` 作为 CI 门。
2. **票 68**（进行中）→ **票 69**（同包，必须串行）。
3. **票 67 AC#3**（`providers.go` 三处字形 + emoji 覆盖面扩展，等 66 让出 `cmd/wisp`）。
4. **票 20 收尾**：AC#5（junction/短名 18 条集成用例只在 `internal/risk` 层有，桥层零命中）+
   **A18**（真 `taskkill` 下的 `.wisp-tmp-*` 残留）+ A17 已修的票面表述同步。
5. **票 64 的 MINOR-1/2**：桌面空出后跑 `-tags winlive -v` 逐条记 PASS/SKIP，
   **交互四项任一走 SKIP 就把 AC#4 退回未勾**。
6. **票 12 AC#8**：owner 视觉签收（R13 只是 INTERIM）。
7. **票 21 段 2**（A19 的应答路径 + 死码隔离）→ 然后才是 S2/S3 的能力票（15/16/22+）。
8. **票 65**：`blocked-on-owner`，等 owner 重发图1/图2 落 `design/refs/`。

- **[A26] 两份**已 done 票的验收证据**里，静态门跑的是"扫描零个文件"的命令，并记成 PASS** —
  我自己核过的原文：
  - `docs/evidence/s1/10-adversarial-acceptance.md:95-96`：`$ cd tools/d22scan && go run .` → 记 `clean`；
  - `docs/evidence/s3/19-adversarial-acceptance.md:20`：**同一条命令**记 PASS。
  **问题**：不带 `-root` 时它**默认扫 `.`**，也就是**扫描器自己那个 module**——
  里面**没有 `internal/`、没有 `cmd/`、allowlist 读不到就被当空**，
  于是它打印 `clean` 且 **exit 0，而实际 examined 0 个文件**。
  ⇒ **票 10 与票 19 的四项安全断言（裸 goroutine / 墙钟差超时 / `filepath.Clean` 越权 / emoji）
  在那两次验收里等于没有门。** 我今天还在多处引用"d22scan clean"作为既有基线，那话的地基就是这两行。
  **已修（票 67，`23ebb59`）**：`checkRoot()` 让空范围**致命退出**，输出里加 **`examined N production Go files`**
  ——**门必须自报工作量**，这正是本条缺的东西；CI 改走 `scripts/d22scan.sh`。
  代理还留了一条**我没让它修**的实话：ban #1 只匹配 `go func(){...}` **闭包字面量**，
  `go probeReader()` 这种**具名函数调用扫不出来**（生产范围内实测 4 处：`internal/observe/goroutine.go:281`
  那一条是受管的 + `cmd/balldebug/main.go:231/259/440`）⇒ 今天的 `clean` 含义是
  "无闭包字面量裸协程"，**不是**"没有绕过 Spawn 的东西"。
  **完成判据（剩余）**：要么把 ban #1 扩到具名调用并**由我书面**给 `internal/observe` 豁免，
  要么在文档里把"clean"的语义写死成当前范围。**不许默默扩语义。** 归**票 70**。
  **元教训（与 A22 同条，但更狠）**：我自己踩的是"调用方式错了所以没跑"，
  这两份证据踩的是"调用方式对扫描器自己成立、但范围是空的"——**同一类假绿的第二种形态**。
  今后凡是引用门禁结论，必须看到它**报出被检查的条目数**；只报"clean"的门一律不算证据。
- **[A27] ⚠ CI 从来没绿过：最新一次完整跑 **5/5 个 job 全红**，五种不同的因** —
  我没有推断，是直接查的（`gh run view 35517463335 --json jobs`，run 时间 2026-09-20T14:45:32Z，commit `0891907`）：
  `test-windows / lint / slo-full / slo-smoke / test-core` 全部 `conclusion=failure`。逐 job 的失败步骤：

  | job | 失败步骤 | 因 | 归属 |
  |---|---|---|---|
  | `lint` | `gofmt (gofumpt)` | 我用 CI 钉的 **v0.7.0** 在本地复跑 `gofumpt -l . tools/d22scan tools/mockllm` → **69 个文件**被标 | **票 70**（一次机械格式化 sweep；`f088ce3` 起就在，所以**这个 job 可能从未通过**） |
  | `test-core` | `Portable package tests` | `internal/observe`：`TestNoBareGoFuncInProductionCode`、`TestSampleStateCPUTotalDrivenMean`；另有 `TestExistsAndDelete`、`TestBlobsListsMetadataOnly` | **票 70** 先复现再分诊（observe 那两条**可能**是票 66/67 改动带出来的，也可能是 ubuntu/Windows 平台差） |
  | `slo-smoke` | `SLO smoke gate` | **就是 A14**（解析丢末项 → 常态 fail-closed） | **票 66**（在跑） |
  | `slo-full` | `Build wisp.exe (deps cached on the runner)` | 自托管 runner 的依赖/构建问题，未查 | **票 70**；`ci.yml` 注释写明该 runner 需 `wisp-slo` 标签 |
  | `test-windows` | `PathResolver junction placeholder (real cases tickets 18/20)` | 步骤名自带 **placeholder**——**这很可能是"设计上先红"的占位**，需要确认是不是该改成显式 TODO 而非失败 | **票 70** 定性；我不在确认前改它 |
  **我这一轮又错了一次，就地更正**：14:59 我在 commit `988c833` 里写"CI lint job 九小时来第一次真绿"——
  **错了**。我验的只是**其中一个步骤**（d22scan），而 `lint` job 里 **gofumpt 步骤排在它前面**且今天仍红。
  **纪律：报"job 绿"之前必须逐步骤看，或直接把 `--json jobs` 的结论贴出来。**
  已被 `gh` 实测否证的推断，我不再用"我记得"顶。
  **本条的意义大于任何单张票**：我们把"CI 门禁"当成本仓的既有护栏写进了 HANDOVER 与各票简报，
  而它**一次都没生效过**——所有"CI 会拦住 X"的论证都缺前提。
- **[A28] 票 68 AC#1 顺手挖出：`SizePx` 的 DPI 处理不对称，且**我们全部证据都是在 96 DPI 下测的**（潜在 MAJOR）** —
  代理读码（`internal/ball/statevisual.go`，`81afe3c` 的注释里记着）：DC render target 以
  `dpiX/dpiY = 96` 创建 ⇒ 其单位**就是物理像素**，`drawFrame` 用 `R = SizePx/2` **不缩放**；
  而**描边宽度与 ring margin 会乘 `dpi/96`、窗口边长走 `WindowEdgePx()`**。
  ⇒ **>96 DPI 时球体保持像素尺寸、窗口却放大**（配 56 的球在 144 DPI 下画 56px 的体、装进 108px 的窗）。
  **为什么今天没人看见**：`docs/SLO.md` §A.2 与 `docs/evidence/s1/62-*` 的差分数字全部来自
  本机 **3440×1440 @ 96 DPI**，那里这个不对称**恰好是 no-op**。
  ⚠ **连带影响点击命中**：`HitTest` 按窗口矩形与 `sizePx/2` 判定，体比窗小 ⇒ 高 DPI 上"看得见但点不着"
  或"点到透明区穿帮"都可能出现，而这正是票 64 刚修过的隐形边框那一族的近亲。
  **完成判据**：①在 >100% 缩放的显示器上实测一次差分化像 + 命中（**本机单屏 96 DPI，做不到 → 需 owner 或第二块屏**）；
  ②或先用一条**纯函数**用例把不对称钉住（喂 `dpi=144`，断言体径==窗口径所要求的比例）；
  ③修不修属渲染行为变更，归**票 68 AC#2 的桌面跑同一批**。
  **另：AC#1 判定我对契约写错了一处数字**——`stateSize` 算出的静止体是
  **56 × `SleepRestRatio`(0.62) = 34.72px**（44/48 配置还会被 `SleepingRestMinPx=30` 兜底），
  **不是我今天在 SPEC-08 §2 INTERIM 里写的"直径 44px"**（46×46 是差分化像框，含 ring 与边缘过渡）。
  SPEC-08 冻结、且这条是我自己写错的 ⇒ **更正文案属"修我自己的录入错误"，不是改契约意图**，
  但我仍把它列进**待 owner 过目**清单，不悄悄改。

- **[A29] 票 68 AC#1 一次性挖出四处"窗口半径 vs 态半径"混用 + 我今天写进契约的数字是错的** —
  先记**方法论**：代理没有新测任何像素（桌面被占），它用**已录数据的可预测性**反推真值——
  模型在非 `Sleeping` 行先验证：`SizePx=56` 预测被窗口裁切的圆盘 **4844** 像素，
  而 `docs/SLO.md` §A.2 的 Thinking/Acting **实测正是 4844**；再用同一模型判 `Sleeping`：
  `34.72 → 2130`（与三次实测 2098/2103/2120 差 0.5–1.5%）、`44 → 3421`（**高出 62%**，且那个 72px 窗裁不掉它）。
  ⇒ **"44px" 被定量否证**，这是本仓少见的"用旧数据证新事实"，我按它改判。
  - **① 我今天在 SPEC-08 §2 INTERIM 里写的「静态玻璃体 44px」是错的**：真体是
    `56 × SleepRestRatio(0.62) = 34.72px`（配置 44/48 时被 `SleepingRestMinPx=30` 兜底）；
    40×38 / 44×44 / 46×46 是**差分化像框**（含 halo），且框列本身在不同跑次浮动 ±6px，**px≥8 才是稳定列**。
    `44` 最可能的来源是 `BallSizeSmallPx`/`size_min`（**最小可配置尺寸**，`tokens.go:369-371`）被误当成体径。
    **更正契约属"修我自己在 6 小时前录入的事实错误"，不是改契约意图**；仍列入**待 owner 过目**，不悄悄改完就算。
  - **② 停靠态与文档相反**：`DockSquash`/`DockOverlapFrac=0.42` 只压到液滴(`:561`)、高光(`:573`)、
    rim 描边(`:577`)，而 shell/halo/caustic/折射核仍是整圆；`DockPos` 的 `orbR` 取**窗口**半径(28)、
    体用**态**半径(17.36) ⇒ 实测露出 **≈82%，不是文档承诺的 42%**（`62-diff-signoff` 的 35×46 行佐证）。
  - **③ 命中半径与可见半径不一致**：prototype 自由态 `HitTest` 的 r = **21px**，而可见 halo 到 **26px**
    ⇒ **外圈看得见但点了穿**；frozen 逃生门更糟：r=**10px** 装在 72px 窗里。
    ⚠ 这与票 64 刚修的"隐形边框把点击挪走"是**同一族**（几何消费者各用各的矩形），**同形排查应扩展到全球包**。
  - **④ 又一条零生产调用者的实现**：`hit.go:58` 的 `SleepWindowEdgePx` **只有测试在用**
    （`tokens_test.go:375`）⇒ **A8/A11/A12/A13/A19/A20 之后这个形状的第 7 次复现**。
  - **⑤ `Settling` 在 prototype 下合成到 ≈0.33**：`ball_windows.go:471/474` 把 ULW `constAlpha`
    在**两种模式**下都驱到 `SleepOpacity=0.35`，而 `applyGlassForm` 让 Settling 落到 `RestSettledOpacity=0.95`
    ⇒ 两个 alpha 相乘。**待我裁定**（这是 `Sleeping` 相邻的契约数字，代理正确地没动）。
  - **⑥ 已登记的地雷（翻转默认值时会炸）**：`live_windows_test.go:609`
    `TestBallLiveAudioLiquidGate` 靠**环境默认值**断言"先测 frozen"，**没有显式 `EnablePrototypeVisuals(false)`**
    ⇒ 票 68 AC#2 一旦翻默认，这条会**自称为 frozen 而实测 prototype**。修 AC#2 时必须同批改它。
  **归属**：①我改契约并告知 owner；②③⑤⑥ 归**票 68 AC#2/AC#3 的桌面跑同批**；④ 与 A24-D4 一起在票 69 或票 39 认领。

- **[A30] ⚠ 一次全量证据自查：我自己盖的三张"D22 clean"是假的，一张票在 0/8 框的情况下被我改了 `-done`** —
  这条是**我（编排者）的错误登记**，不是代理的。触发点：票 70 排队时我顺手把"done 票的证据是否还成立"
  过了一遍。**逐条我自己复现过才写在这里**（扫描代理给了四条，第四条我否证了，见 ③）。
  - **① 票 62 是 `-done`，但它的 8 个 AC 框**一个都没勾**（`grep -c "^- \[ \]"` = 8，`^- \[x\]` = 0），
    且 AC#8 要求的 `62-adversarial-acceptance.md` **根本不存在**（`docs/evidence/s1/` 下只有
    `62-ball-visual-prototype.md` 与四个 `62-diff-*` 目录）。
    ⇒ 直接违反 README 完成规则第 4 条（"check **all** acceptance boxes → rename `-done`"）与第 6 条
    （"置 done 前，验收报告必须含与 AC 编号 1:1 的裁决表"）。**这个改名是我今天做的**，
    做的是把票 65 的 INTERIM 签收当成票 62 的验收——**owner 原话是"算是赝品…先勉强用吧"，那是降级放行不是验收**。
    ⇒ 处置：**改回 `62-liquid-glass-ball-visuals.md` / `Status: review`**，理由写进票面；
    真正的裁决表要么由后续跑补齐（AC#3/#5 需要桌面），要么由票 65 消化。
  - **② 三张 done 票的验收报告里写着"d22scan 全仓 clean"，而这个命令的真实覆盖面是零**：
    `s1/08:10`「d22scan 实跑全仓 clean」、`s1/18:11`「tools/d22scan 全仓 clean（allowlist 1 行豁免…）」、
    `s1/17:12`「越界/D22 | PASS」——**且 s1/17 那行既没有命令也没有 file:line**，是一张无凭据的 PASS。
    根因与 A22/A26 同一条：在仓库根跑 `go run ./tools/d22scan -root .` 只打印一条**模块错误**就退出，
    我把"没有 findings"读成了"干净"。真调用（`cd tools/d22scan && go run . -root ../..`）当时**exit 1**
    且指向 `internal/llm/adaptertest/mockllm.go:68`。**这三张票的其余 AC 不受影响**（各自有独立证据），
    受影响的只是"越界/D22"这一行——它当时的**证明力为零**。
  - **③ 我否证了扫描代理给的第四条**：它把 `s1/05` 也列成假 D22 章。我读了原文：
    `s1/05:12` 那行写的是「越界 | PASS | 实现提交仅 internal/config + go.mod/go.sum + 票 05」——
    那是一个**提交范围检查**，全文没有提到 d22scan，不构成假章。**故本条只算三张，不是四张。**
    （登记这个否证，是为了不让"证据自查"本身变成一条新的未验证引用。）
  - **④ 已失效的具体引用（今天抽验）**：`s1/09:42` 引的 `TestGoldenCancellationMidStream`
    （连同 `adapter_test.go:389-427`）在树里**已无定义**（我抽了该报告 6 个测试名，5 个仍可解析、这 1 个不能）。
    `s0/01` 的 emoji 扫描面主张已过期（同类命中今天散布 20+ 文件，见 A23）。
    `62-visual-spec-draft.md` 被票 62 的正文(`:51`)与 AC#7(`:60`)引用，**文件从未写出**（`find` 零命中）。
    - **✅ ②④ 已于 2026-09-20 23:59 闭（提交 `a5741a2`，我逐条独立核对）**：
      四份报告各在**文末追加**一节 `## 更正（A30…）`，`--stat` 实测 **527 insertions / 0 deletions、
      `.go` 文件 0 个** ⇒ 原文一字未改，历史证据没有被销毁。并排贴了两条命令的真实输出
      （仓根那条**不打印 examined 且 EXIT=1**、正确那条 **`examined 194` + clean + EXIT=0**），
      没有拿"今天 clean"去追认"当时 clean"。
    - **④ 的两处修正，以核对后的版本为准**（我先前那句话也偏了）：
      (a) `TestGoldenCancellationMidStream` **不是断言丢失**——`5ddedf7` 把它逐字搬进
      `internal/llm/adaptertest/harness.go:579 runCancel`，现名 **`TestSharedGoldenSuite/cancel`**
      （`internal/llm/anthropic/suite_test.go:36` 等三协议腿各跑一遍），我 grep 复现了这两个位置。
      ⇒ 受影响的只是**引用字符串与行号**（`adapter_test.go` 现在只有 291 行，427 已越界），不是覆盖面。
      (b) `s1/18:11` 的"**allowlist 1 行豁免**"在**写下当时就失实**：`git show 3df0218:tools/d22scan/allowlist.txt`
      非注释行 = **4**（其中属票 18 的只有 1 行）——我自己刚复跑这条命令确认 =4。
      ⇒ 这不是"引用过期"而是"原句写错"，性质更重一档。
    - **代理如实声明的未证实项（不要当成已证）**：三处 D22 句的**历史**状态需要 checkout 旧树才能定，
      它按授权没做，故那些行判"无凭据"而非"当时确红"。

  - **⑤ 元教训：门必须自报工作量，否则"没报问题"和"没看"在输出上长得一模一样。** 本轮挖出的同类空仪器：
    `cmd/balldebug -diff` 无论像素计数为何都 `return nil`；`go test -run <不匹配>` 打印 `ok`；
    d22scan 的 `frontend/` 作用域实际走 0 个文件。票 67 已把 `examined N` + `N==0 致命退出` 落进 d22scan
    （R16 裁定 3 禁止回退），另外两条**已建票面收口**：**票 71 `71-gates-must-self-report.md`**
    （AC#1..#7，含"必须给出一次真实的红"作为 AC#2；`blocked-on-70` 只因为两者都要改 `.go`）。
  **归属**：①我今天就地改回（已完成：文件回到 `62-liquid-glass-ball-visuals.md`、Status→review，
  提交 `fc33532`）；②④ 归**证据更正代理**（只追加、不覆写原文，每条带可复现命令）；
  ⑤ 归**票 71**；②的根因修复（正确调用姿势）已在 A22/票 67 落地。
  - **⑥ 同轮补记（自查不能只查别人）**：票 66 的活**全交了但缺 README 规则 6 的 1:1 裁决表**，
    文件也一直没改名——这是与 ① **同一规则的另一个方向**的破口（①是没做完就改名，⑥是做完了没补表）。
    于 23:56 闭：补 `docs/evidence/s1/66-adversarial-acceptance.md`（7/7 PASS，每行标注
    "编排者独立复现／代理日志＋归档我抽验／仅代理自述"三档）后改名 `-done`。
    ⚠ **其中 AC#1 的变异检验我在第三档**（没独立重做，因为改 `.go` 会撞票 70 的全仓格式化）
    ⇒ **登记一条必做后续**：票 70 落地后，把 `internal/proc` 的 `WalkSystemProcesses` 顺序退回旧实现，
    确认 `TestParseSystemProcessesKeepsLastSnapshotEntry`（`systemprocs_windows_test.go:79`）变红。
    **这条没做之前，票 66 AC#1 的证明力 = 归档文档 + 我核对过的代码形状，不等于"护栏被真验过"。**
    已同时写进票 66 票面的"残口四条 ④"。


## 待 owner 拍板（编号清单 R15，2026-09-20 23:14）
规则：每项给「选项 / 我的推荐 / 不答的代价」。**1–6 全部不需要你写代码**，只回答问题或跑一条命令。

| # | 问题 | 选项 | 推荐 | 不答的代价 |
|---|---|---|---|---|
| **1** | 我今天把 SPEC-08 §2 的 `Sleeping` 写成"直径 44px"，**已被定量否证**（真体 34.72px；44 是差分化像框）。我已就地更正并标注"待你过目"。 | 追认 / 回退我的更正 | **追认**（原文与证据链都在 §2 更正块里） | 契约里留一个证明是错的数字，票 65 会照它返工 |
| **2** | 同节 `Armed` 行写"直径 44px"，而代码用配置尺寸 56。这是**既有契约**与代码的分歧，我没自行动。 | 改契约文本为"配置尺寸" / 改码到 44 / 保持现状 | **改契约文本**（代码自 S1 起如此且有实测） | 每次对照规格都要先解释"为什么差 12px" |
| **3** | prototype 下 `Settling` 合成到 **≈0.33**（ULW `constAlpha`=0.35 与 glass 0.95 双乘）。回落中的球该多可见？ | 0.95 / 0.35 / 就要 0.33 | **0.95**（`Settling` 是"正在回落"，看不见会让人以为卡死） | 球在回落后半隐形的观感会被当成 bug 报回来 |
| **4** | 停靠时实测露出 **≈82%**，文档承诺 42%（squash 用窗口半径、体用态半径）。你要的是迅雷那种"半个球藏在边里"吗？ | 修码到 ≈50% / 把文档改成 82% | **修码到 ≈50%**（"贴边收起"的语义就是藏一半） | 票 62 的吸附效果与你签收时看到的不是同一个东西 |
| **5** | 命中半径 21px < 可见外沿 26px ⇒ **最外圈看得见但点了穿**（frozen 逃生门更糟：10px 装在 72px 窗里）。 | 对齐可见外沿 / 维持 | **对齐可见外沿** | 用户会反复点边缘没反应，且这类反馈最难归因 |
| **6** | **>100% 缩放的显示器上球体不放大**（DC 目标固定 96 DPI），窗口和描边放大 ⇒ 56 的体装进 108 的窗。**我们所有证据都是在 96 DPI 上测的**，所以今天没人看得见。 | 你借一块高分屏跑一条我给的命令 / 我先加纯函数断言锁住 / 推到 S3 | **我先加断言锁住**，同时请你在真机跑一次 | 上线到笔记本用户那里才发现"球小得像噪点" |
| **7** | 票 65（玻璃质感返工）还缺**图1/图2**——到今天为止没有任何代理看见过参考图。 | 重发两张图给我落到 `design/refs/` / 明确说"S3 再说" | **重发**（否则票 65 第一步就动不了） | INTERIM 标记一直挂着，"赝品"变成长期默认 |
| **8** | **CI 从来没绿过**（A27：五个 job 全红）。要不要我把票 70 排到功能票前面？ | 先转绿 CI（建议顺序：gofumpt sweep → slo-smoke → test-core 分诊）/ 先推功能 | **先转绿 CI** | 我们所有"CI 会拦住"的论证继续缺前提，且红 job 会掩盖真回归 |
| **9** | `slo-full` 跑在带 `wisp-slo` 标签的**自托管 runner**上，就是你这台机；它的 build 步骤现在是失败的。 | 我远程诊断并配好 / 暂时把该 job 标记为不可用 | **我诊断**（它一红，D32 的全量门就没人跑） | SLO 全量门长期无人执行 |

- **[A30] 票 68 的 ①–⑥ 已归 registry A29；本节 R15 是它们中**只有 owner 能答**的那部分** —
  纯索引条目，避免"发现记了但没人拍板"。对应关系：R15#1/#2↔A29①，#3↔A29⑤，#4↔A29②，
  #5↔A29③，#6↔A28/A29 DPI 条；A29④（`SleepWindowEdgePx` 零调用者）与 A29⑥（winlive 地雷）
  **不需要 owner**，分别归票 69 与票 68 AC#2 同批。

### 裁定 R16（23:27）：A26 剩余半——ban #1 覆盖面怎么扩，以及**我更正代理的一处计数**

我读了匹配器本身（`tools/d22scan/main.go:251-252` 只对字面量 `go func(` 报警），
并逐条看它报告的"生产范围内 4 处具名协程"。**其中一处不是违规**：

| 位置 | 语句 | 判定 |
|---|---|---|
| `internal/observe/goroutine.go:281` | `go r.run(h, name, root, owner, fn)` | **这就是受管机制自己的实现**——禁令要禁的是"绕过 Registry 起协程"，而这一行**是** Registry 的内部实现。它不是漏网，是**唯一合法的 spawn 点** |
| `cmd/balldebug/main.go:231` | `go runHotkeyBridge(bridge, pollCtx, time.Second, pollDone)` | 具名绕过（票 64 写的，带 owner+recover+join，**精神合规、形式不过**） |
| `cmd/balldebug/main.go:259`、`:440` | `go feedLevels(...)` ×2 | 同上 |

⇒ **代理说"4 处"应当读作"3 处要处理 + 1 处是机制本身"**。它没让我修是对的（豁免权在我），
但它把 `observe` 那一行与 balldebug 三行并列为"命中"，会误导后来的代理去豁免真正的漏洞、
或反过来把机制改掉。

**裁定（写死，票 70 AC#5 照此执行，不要再问我）**：
1. **豁免按"文件"而不是按"形式"**：新增 ban 匹配 `go <任意>`（含具名调用），
   **唯一按文件路径豁免的是 `internal/observe/goroutine.go`**。
   ⚠ **不得**用"这条是具名调用所以放过"当豁免理由——那正好把漏洞留在门上。
2. **`cmd/balldebug` 的 3 处要改成 `observe.Registry.Spawn`**，不是加豁免。
   理由：`balldebug` 是长期存在的调试宿主，它的协程寿命与球窗口绑定，
   走 Registry 才能被 D38b 的名册与泄漏检查看见；现在这三条**只在人肉 review 下成立**。
3. **输出必须自报工作量**（A26 的元教训落地）：扫描结束打印 `examined N production Go files`，
   且 `N == 0` 时**致命退出**——这条票 67 的 `checkRoot()` 已经做了，扩展覆盖面时**不许回退**。
4. **只许从严**：若扩展后在别处又挖出具名 spawn，**逐个登记进票 70 的 log**，
   由我逐条判定"改码 or 按文件豁免"，**不许默默加豁免凑绿**。
5. 这条与 `mockllm` 的 loose end 挂钩：`Spawn` 对不在名册的名字**无条件 WARN**，
   票 67 因此每条 mockllm 测试多出一行 `WARN goroutine outside the D38 roster`。
   抑制它需要在 `internal/observe/goroutine.go` 加 `TemporaryNames` 条目——
   **那是票 66 的文件**，所以**等 66 收尾再做**；代理"拒绝借用现有名册名来骗过泄漏检测"的判断**是对的**，
   骗过检测器比多一行 WARN 贵得多。

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
- **[A14] `parseSystemProcesses` 丢快照末项 ⇒ state 口径从未出数** → **已解决（票 66，AC#1+AC#2）**。
  关闭 commit **`86e868d`**（唯一走链实现 `internal/proc/systemprocs_windows.go::WalkSystemProcesses`，
  「先解析当前项、再测 `NextEntryOffset==0`」）。判据：回归用例经**第二人独立变异复测**咬住缺陷本身
  （退回旧序 ⇒ `internal/proc` 6 红 / 24 绿，含真机自快照用例与两条树外采样用例）；
  **零 keeper** 连跑 5 次 `wisp slo -state Sleeping -seconds 10 -interval-ms 250` = `exit=0/0/0/0/0`、
  九项门每次各出 9 行。证据：`docs/SLO.md` 附录 C.1/C.2、`docs/evidence/s1/66/66b-nokeeper-{1..5}.json`、
  `docs/evidence/s1/66-mutation-parse-order.md`。A.6#1 的 keeper 绕法已从文档删除。
  ⚠ 上面 A14/A15 两段的 finding 原文**一字未删**（本段是叠加的结案行，不是替换）。
- **[A15] 树内 CPU 门测的是观测者自己 ⇒ D32 `Sleeping` 行无法判定** → **已解决（票 66，AC#3）**。
  关闭 commit **`00bbb76`**：`wisp slo` 的 state 段现在同一 Job 里起两个子进程，**被测主体由父进程树外读取**
  （`internal/proc/externalsampler_windows.go`，复用同一份快照解码器），遗留的树内自采样降级进
  `observer` 段并标 `observer_cost: true` + `gate: false`。**D32 的 0.5% 阈值一字未动**：
  两侧 limit 串由同一个 `stateCPULimit` 生成，用例 `internal/observe/observer_cost_test.go` 钉住「只许 CPU 行动」。
  判据（双条同窗成立）：30s/250ms 内 **树外 0.0173%** 且 **树内 0.7629%**；5 个 10s 零 keeper 跑内
  树外 0.0000–0.0130% 对树内 0.5963–0.9318%。证据：`docs/SLO.md` 附录 C.1/C.3、
  `docs/evidence/s1/66/66-ac3-30s.json`。
- **[票 66 复跑时新挖，第三起仪器缺陷；正式编号请编排者给] `slo-check.ps1` 在「全过」路径上崩溃，
  报告根本不写 ⇒ CI 红且无 artifact** → **已解决（票 66，AC#4）**，关闭 commit **`cf0c050`**。
  `scripts/slo-check.ps1:129` 对 `Where-Object` 的结果取 `.Count`，而 `Set-StrictMode 2.0` 下
  零命中是 `$null`、单命中是标量 ⇒ 恰好在所有 state 都 pass 时抛终止性异常。
  **它是 A14 的下游**：A14 让每个 state `exit=2` 全 fail，过滤器命中 ≥2 项才拿到真数组，掩盖了这条。
  ⇒ 编排者 B.3 的预告「只修①⇒约一半概率随机红」方向要反过来读：**只修① = 必红且没报告**。
  判据：`-Subset smoke` 连跑 **5/5 exit=0**（`docs/SLO.md` C.4，五份日志 `docs/evidence/s1/66/66-smoke-run-{1..5}.log`，
  崩溃原文 `66-smoke-pre-fix-crash.log`）；阳性对照 = 同一表达式在 StrictMode 下喂 0/1/2 项失败得
  `allPass=True/False/False`，即修完仍会红。附带：`-Subset full`（六态）exit=0，但六态 `posture` **全是 skeleton** ⇒
  不得当作 D32 六态表达标（边界写在 C.5）。
