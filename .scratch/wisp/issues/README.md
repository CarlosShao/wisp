# Wisp Ticket Index

Source: `docs/specs/SPEC-*.md` (authoritative for all numbers/rules cited here).
Slice order: S0→S1→{S2∥S3}→S4→S5→S6→S7→S8 (SPEC-12 §2).
All tickets are **vertical tracer bullets**; work the **frontier** (all blockers done).

## Status lifecycle (anti-interruption rules — MANDATORY)

| Status | Meaning |
|---|---|
| `ready-for-agent` | Frontier ticket, unclaimed |
| `in-progress` | Claimed; must have ≥1 Progress-log entry |
| `blocked` | Waiting on external decision (ticket notes which) |
| `review` | Work done, awaiting adversarial acceptance (SPEC-10 §8) |
| `done` | All boxes checked; **rename file with `-done` suffix** + title `(DONE ✅)` |

**Rules (prevent abandoned/in-progress tech debt):**
1. Before writing any code: set `Status: in-progress`, fill `Claimed by`, `Last update`,
   append a Progress-log line → **commit + push (both remotes)** in the same step.
2. Append a Progress-log line after every meaningful unit (`- [UTC ts] agent=<id> did=… next=…`),
   commit+push. Sub-agent lines tagged with their id. **Max 2 sub-agents per ticket**（每张票仍按票内
   `Parallel slots` 走，同票多写手会在同一工作树里互相踩）；**整体车队宽度 3**（R3，2026-09-20 用户批准，
   旧值 2 是 ZCode 免费配额校准的、已失效）。附加口径：需安静测量的票（SLO 私有工作集/延迟分段/CER）
   **独占 1 个代理**、其余排队；只读检索/审计类可放到 10+。>1 个写码代理必须 worktree 隔离
   （`third_party` 用 junction 复用，避免每个工作树重下 1.3GB）。
3. An `in-progress` ticket with no log update for >24h: next agent must verify on-disk state
   against the log before continuing; never reset checked boxes; append `resumed` line.
4. Completion: check all acceptance boxes → `Status: done` → rename file `NN-slug.md` →
   `NN-slug-done.md` → update this index → commit+push.
5. If a decision is missing, set `Status: blocked` with the open question in the log
   (D22 闸门③: undefined = stop and ask, never assume).
6. **置 done 前，验收报告必须含一张与 AC 编号 1:1 的裁决表**（每行 = 一条 AC 的 PASS/FAIL/PARTIAL
   + 判定依据的 `file:line` 或实跑输出）。**缺行即 FAIL**，无论其余部分多好。
   （2026-09-20 owner 批准新增。起因：票 07 标 DONE 时 6 条 AC 只有 1 条被书面裁决过，
   其中"热键随配置重注册"的标的 `RebindHotkeys` 实为全仓零调用者的死代码；
   补裁另有 2 框 FAIL、2 框 PARTIAL。旧报告不追溯改写，但 13 张已 done 票的未决框由票 64 等消化。）

## Dependency graph (blockers in parentheses)

- **S0**: 01 build-chain ✅done · 02 s0-spike ✅done (01)
- **S1**: 03 skeleton ✅done (01) · 04 sqlite-core ✅done (03) · 05 config-model ✅done (03) · 06 secretstore-envs ✅done (03) ·
  07 ball-state-machine ✅done (03) · 08 observability-slo ✅done (03,07) · 09 llm-provider+mockllm ✅done (03) ·
  10 agent-loop-core (05,09) · 11 llm-adapters-rest (09) · 12 cli-text-path=S1 gate (04,07,08,10)
- **S2**: 13 audio-capture ✅done (03) · 14 model-distribution ✅done (03) · 15 speech-engines+CER (02,13,14) ·
  16 s2-acceptance (12,15)
- **S3**: 17 risk-assessor C19 ✅done (03) · 18 path-resolver C26 ✅done (03) · 19 provenance C25 ✅done (17) ·
  20 host-bridge+fs-tools (17,18) · 21 approval-gates-minimal (17) · 22 web-tools+D30 (19,20) ·
  23 system/window/input-tools (20) · 24 doc/search-tools (20) · 25 s3-acceptance (21,22,23,24)
- **S4**: 26 tts-output (15) · 27 punctuation (15) · 28 session-scope Warm (07,15) ·
  29 memory-l1/l2 (04,10,28) · 30 result-routing D10 (10,26) · 31 reminders (04,30) ·
  32 s4-acceptance scenarios ③④ (27,29,30,31) — **26/28/32 的 Path C 部分（D47 全双工/barge-in）
  另被 59 阻塞**
- **S5**: 33 panel-host C27 (07,12) · 34 frontend-scaffold (01) · 35 panel-bridge C17 (33,34) ·
  36 result-history-panel (30,35) · 37 approval-ui (21,35) · 38 palette+tasks (35) ·
  39 config-editor-gui (05,35) · 40 security/privacy/cost-pages (29,35,44)
- **S6**: 41 kws-wake-word (15,28) · 42 watchdog (08,28) · 43 power-lifecycle (28) ·
  44 costmeter C23 (04,10) · 45 diagnostics+guards (42) · 46 s6-acceptance (41–45)
- **S7**: 47 task-scheduler+pathlock (10,21) · 48 approval-queue-full (37,47) ·
  49 session-grants D45-2 (04,48) · 50 tier1-manifest-plugins (20) · 51 tier2-goja (17,20) ·
  52 d46-command-plugins (50) · 53 larkcli-plugin-e2e (48,51,52) · 54 s7-acceptance (47–53)
- **S8 (deferred, not ready-for-agent)**: 55 macos · 56 signing+dist+naming · 57 plugin-sdk+registry · 58 i18n
- **D47 追加（2026-09-19，编号晚于 58 但入图拓扑早于被阻塞票）**：
  **59 p15-aec-spike (13) → 阻塞 26/28/32 的 Path C 部分与 60**；
  **60 c32-realtime-engine (16,28,44,59；门控双条件：59 通过 + 用户有 Key，缺一自动推迟登记)**。
  ⚠ 编号惯例说明：后补票用更大编号以保持文件名稳定，忽略编号与拓扑序的差异，以 Blocked by 为准。
- **LLM 接入补充（2026-09-19 用户批准）**：
  **61 cloud-voice-providers-C9 (16,11,05)** —— 云级联 ASR/TTS（语音模型下拉的落地层，
  候选 StepAudio 2.5 ASR / StepAudio 3 TTS）；配套修订：票 05/09/11/40/44/60 均已注入
  目录/能力位/计费模式/配额三层/兜底链/角色化默认/probe/限流要求；
  SPEC-02（provider_health v2 表）、SPEC-03 §3.1、SPEC-05 §3.3a、SPEC-12 登记表（云端语音
  DEFERRED→已采纳·提前）同步。

- **后补票（2026-09-20 owner 裁定 R8/R9/R10，编号续 61 之后以保持文件名稳定）**：
  **62 liquid-glass-ball-visuals (07,13)** —— 液态玻璃球视觉重做 + 靠边吸附收缩 + 音频驱动液面，
  含 SPEC-08 §2 视觉表变更（**先原型、owner 实机签收后才回填契约文本并签字**）；
  **63 credential-entry-cli (06)** —— `wisp secret set/get/list/unset` 隐藏输入→DPAPI，
  兑现 R7「key 绝不进对话框」；
  **64 ball-defects-hotkey-interactive (07)** —— 消化 registry A1–A7（热键接线、交互四项、
  Sleeping 零定时器实测、多显示器实拖），默认唤起键 `Ctrl+Alt+W`→`Ctrl+Alt+Q`。
  ⚠ 62 与 64 的边界：64 只管"球能被唤起、配置改动能生效"，62 只管材质/动效/吸附；两票都不得
  为对方让路而放宽 D32 的 Sleeping 零定时器与 CPU ≤0.5% 纪律。
- **再后补票（2026-09-20 票 62 签收轮与票 12 桌面跑之后）**：
  **65 ball-glass-quality-rework** —— owner「赝品」判决在此消化，`blocked-on-owner`（缺参考图，
  第一步是把图落到 `design/refs/`），SPEC-08 §2 的 INTERIM 标记只在本票被 owner 签收后方可移除；
  **66 slo-instrument-parse-cpu-observer** —— SLO **判据仪器本身**返工（registry A14/A15）：
  `parseSystemProcesses` 丢快照末项（state 口径从未出数、CI 两处 SLO 步骤坏着）+ `Sleeping` CPU 行
  在树内口径下**无定义**（观测者单次读 ≈1.3ms，250ms 间隔下本身就超 0.5% 门）。
  ⚠ **两起必须一起修**：只修前者会把 `slo-check.ps1` 从"必红"变成"约一半概率随机红"。
  **测量类独占**，不与任何其他跑测的代理并发；**D32 阈值不得因本票改动**。
  **67 d22scan-trust-mockllm-goroutine** —— 让 **D22 静态门重新可信**（registry A22/A23）：
  `tools/d22scan` 是**独立 Go module**，编排者从仓根调用它等于**从未跑过**却读成 clean
  ⇒ 真实命中 `internal/llm/adaptertest/mockllm.go:68` 裸 `go func(`（`78b1466`/票 11 引入、不在 allowlist）
  ⇒ **CI lint job 自约 05:59Z 起就是红的**；且 emoji 门的 `emojiRe` 只看 `design/` 与**不存在**的 `frontend/`，
  对 `internal/`+`cmd/` 全盲（票 12 自己的 `cd011b8` 带进 3 处 U+2713/2717）。**不许用 allowlist 豁免凑绿。**
  **68 ball-default-visuals-parity** —— 让**默认构建画的就是 owner 签收的那个球**（争议 D1/D2）：
  `prototypeVisuals` 默认 **false**，唯一开启者是 `cmd/balldebug/main.go:104` ⇒ 今天默认库仍画 12px 微点，
  而我改过的 SPEC-08 §2「44px 静态玻璃体」描述的是**非默认配置**（这笔账是编排者欠的）。
  另含 Sleeping 尺寸三方不一致（44 vs `stateSize`=34.72 vs 旧行 12，而 A.2 实测像框 46×46）。
  ⚠ **SPEC-08 冻结：不许改文本凑数**；票 66/67 未收尾前**不得跑 winlive**（桌面被占，先交付静态半）。
  **69 c21-token-table-machine-check** —— C21 表新补的 **61 行（25 几何动效 + 36 look 色）无任何机器检查**
  （`TestTokenGoldenValues` 只管 20 条配色，`TestNoHardcodedColorsInBallPackage` 管的是另一件事）⇒
  要一条**双向**表↔码一致性断言 + 变异检验。⚠ **被票 68 阻塞：两票同动 `internal/ball`，必须串行**。
  **70 ci-actually-green** —— **本仓的 CI 门禁从未生效过**（registry A26/A27）：
  实测 `gh run view 35517463335 --json jobs` ⇒ **`lint`/`test-core`/`test-windows`/`slo-smoke`/`slo-full`
  5/5 全红**（gofumpt 标 **69 文件** · observe 等 4 条测试红 · junction **placeholder** 步骤 ·
  A14 解析缺陷 · 自托管 runner 构建失败）。⚠ 硬规矩：**不加 `continue-on-error`、不删步骤、不下调阈值凑绿**
  （D22 "no job skippable"）；格式化 sweep 与逻辑改动**永不同 commit**；判据 = **逐 job 全 pass**，
  不是"我本地某一步过了"。**blocked-on-tree**：AC#1 全仓重写，须等票 66/68 的包空出来。
- **再再后补票（2026-09-20 A30 全量证据自查之后）**：
  **71 gates-must-self-report (70)** —— 把"没报问题"与"没看"变成机器可区分的两件事（registry A22/A26/A30⑤）：
  `cmd/balldebug -diff` **无论像素计数为何都 return nil**（一次球完全不成像的跑在 CI 眼里与理想跑同形）、
  `go test -run <不匹配>` **打印 ok**（票 64 AC#4 的"14 项全 PASS"含 SKIP 即为此形状）、
  `d22scan` 的 `frontend/` 作用域**实际走 0 个文件**。票 67 已给 d22scan 装了 `examined N` + `N==0` 致命退出，
  本票把同一纪律推到其余仪器，且 **AC#2 要求一次真实的红**（阳性对照），否则不算完。
  ⚠ 只许从严：不许为了让某次跑变绿现场调低可见性阈值（D22 精神）；`allowlist.txt` 只许变短或不变。
  **blocked-on-70**：本票要改 `.go`，必须等全仓 gofumpt 那一次重写落地，否则每个文件都撞车。
- **票 70 死前发现（2026-09-21 09:11，编排者建票，优先级高于普通票）**：
  **72 atble-classification-runner (18)** —— **A 表（禁区不可放行）在 GitHub `windows-latest` 上退化成
  B 表（一次 L2 确认可放行）**，本机同命令 PASS、`pathresolver_junction_windows_test.go:104` 在 runner 上 FAIL。
  ⇒ **纵深防御真破一格**，不是断言写松、也不是占位步骤（那位的 `a04d3e2` 已逐条排除
  `os.UserHomeDir`/`t.Setenv`/`normPath`/`isUnder`/A 表命中，且**一字未改断言**只加诊断输出）。
  ⚠ 本票要动 `internal/risk/`＝**冻结区 ⇒ 改法先报编排者判**；判据要求"归一化两侧、不是加特例"、
  双向变异、**本机 + runner 两侧都过**。**为什么不留在票 70 里**：安全分类失效留在"让 CI 变绿"的票内，
  极易被下一个代理用"改断言/加 skip"的最短路径解决掉。
- **票 20 落地后新建（2026-09-21 09:33，编排者建票）**：**73 orphan-staging-file-sweep** ——
  `761447f` 用**真 `taskkill /F`** 测出：每次打断恰好留 1 个 `.wisp-tmp-*` 且**没有任何东西扫它**
  （D31 本身成立，目标文件逐字完整或不存在）。本票把那条评论式钉住"仍在"的断言**翻成"已清"**，
  判据核心是**归属证明**（不许删别人的同名文件）与**扫描时不跟 junction/symlink**——
  否则"清理"会变成一套能写到授权目录外的删除原语。owner 侧同名条目 **Q-16**。
- **票 69 落地后新建（2026-09-21 09:33，编排者建票）**：**74 c21-table-must-mirror-live-drivers** ——
  机器检查装了以后报出 **29 条漂移且不判红**，其中最危险的一族是**语义漂移**：`tokens.go` 声明
  `SwimLevelGain`/`SpinLevelGain` 而**零消费者**，真正驱动像素的是 `liquid.go` 的 `Spin*RadPerS`
  ⇒ 与 **A33（声明✓/实测✗）同族**。⚠ 删已记录的 token **必须与替代同 commit**；
  表↔`tokens.css` 那格若本票不收，**必须写成带票号的书面接手**，"暂缓"无票号即不合格。

- **状态回写（2026-09-20 23:5x，编排者自我更正）**：**62 liquid-glass-ball-visuals：`-done` → `review`**。
  我今天在**八个 AC 框一个都没勾**（实测 `^- [ ]`=8 / `^- [x]`=0）且 **AC#8 要求的
  `62-adversarial-acceptance.md` 不存在**的情况下给它加了 `-done` 后缀——同时违反规则 4 与规则 6。
  根因：我把 owner 对**票 65** 的降级放行（「算是赝品…先勉强用吧」）当成了票 62 的验收结论。
  ⇒ 票面留了一张八行的"每条 AC 现在归谁"表（#1/#6 有证据缺裁决、#3/#4 归票 65、#5 归票 68、
  #7 的 `62-visual-spec-draft.md` **从未写出**、#8 **从未执行**）。**新会话请勿按 AC#1 重跑实现**，
  按那张表认领。登记为 registry **A30①**。
  同轮：`07` 的五个未勾框补了逐条归属（其中 #4 我写下"断言不存在"后一次 grep 自我推翻，
  真话是"断言存在、但藏在 `-tags winlive` 后面、今天是否真执行未证"）。
- **状态回写（2026-09-20 复验轮）**：**10 agent-loop-core → DONE**（编排者亲自独立复验：把
  `enforceTotal` 整体退回修复前实现，被提交的终止测试 10.00s 变红，护栏为真；两处小项登记 A10-a/b）；
  **63 credential-entry-cli → DONE**，但 **AC#6 未勾、正式转票 12**（登记 A8）、**MINOR-1 未修**（登记 A9）
  —— 别因为看到 `-done` 就以为它零残余。

⚠ **Safety-incomplete period**: tickets 21–32 run with `TaskScheduler` single-task only;
multi-task concurrency unlocks only at 47 (SPEC-12 §2).

## Hard global constraints (apply to EVERY ticket)

- Never modify D1–D46 / C1–C31 / R1–R9 / D43 transition table (contract change = human approval).
- Forbidden patterns (auto-checked in CI): bare `go func(` without owner/recover; `filepath.Clean|Abs`
  for fs decisions outside `risk.PathResolver`; plaintext API keys; wall-clock-difference timeouts;
  hashes fetched from mirrors; panel-sourced L2 "allow"; host-internal artifact writes implemented
  as gated Tools; any emoji in UI code (U+2190–U+2BFF, U+1F300–U+1FAFF, U+FE0F).
- `DEFERRED(D-xx)` code markers must map 1:1 to SPEC-12 §5 registry (bidirectional check).
- Tests inject at seams only: C8 AudioSource (wav), C5 LlmProvider (golden SSE), C17 PanelBridge,
  CLI `wisp run`. No mock-instead-of-real to fake completion; never weaken SLO thresholds.
- Each slice completion requires: gap audit + pre-mortem + adversarial acceptance by a DIFFERENT
  agent (SPEC-10 §7–8).
