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
- **[R7] 凭录入方式约束（用户 2026-09-20 明确要求）**：**任何 API key 不得经对话框传递**，必须由用户自己在界面/本地隐藏输入里录入。现状核查：`internal/secret`（票 06，DPAPI）已有 `Store.Store(ref, secret)`，但**全仓没有录入入口**（无 `wisp secret` 子命令、配置 GUI 属票 39/40 即 S5）。因此 H2/anysearch/StepFun 三项凭据当前只能走「用户自己 `setx` 环境变量 + 我只引用变量名」的临时路径；已提议把最小录入入口（`wisp secret set` 隐藏输入）提前做，待用户批准。

## 完成度审计（2026-09-20 08:58，audit-A + audit-B）

- **结论：13 张 done 票无 BLOCKER；实质工作全部有证据支撑**（audit-B 深查 4 张的引用抽验全属实）。
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
