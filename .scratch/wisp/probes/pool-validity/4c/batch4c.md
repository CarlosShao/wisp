# pool-validity-4c / 批4c —— 零勾开放票有效性普查（号段 15–26，8 枚最老首轮切片票）

> 只读普查腿 `pool-validity-4c`。没有前文，本件自证。
> 起手 HEAD `3d9b8374`（2026-10-06 13:12:18 +0800，分支 `dev`），量尺时刻 10-06 13:1x +08。⚠ 换 HEAD 要重量。
> 本腿**零 go 命令**（并行写腿 `236-r2` 此刻在 `internal/tools` 取突变读数，抢一把就洗掉它）。
> 枚数／行数一律 `git ls-tree`／`git ls-files`／`wc`；看码一律 `git show HEAD:<path>`／`git grep … HEAD`。
> 票面零改动：零翻勾、零改名、零 `Status:` 改写；台账 `docs/reports/pending-and-issues.md` 与 `docs/**` 一字不碰。
> 唯一写入路径＝`.scratch/wisp/probes/pool-validity/4c/**`（新建；起手 `ls` 验过 `4c` 未占用）。
> ⛔ 零读零引 `frontend/**`、`design/**`、`.gitignore`、`.scratch/wisp/probes/236/**`、`.scratch/wisp/probes/pool-validity/4b/**`。
> ⛔ `cmd/wisp/**` 与 `internal/**` 的工作树内容面不读（写腿地界），一律走 `HEAD:` 尺。
> `docs/evidence/s1/**` 只按**文件名**取、不读内容。
> 本段与 §0 之外的所有读数都是**本腿现量**，不采信任何旧腿结论；号段 17–20、27 以上不归本腿。

## §0 四档尺（本腿判法，学 `.scratch/wisp/probes/pool-validity/3a/batch3a.md` §0 的表形，口径自立）

| 档 | 判据 | 复法（必须留命令原文＋读数） |
|---|---|---|
| **仍成立** | 票面点名的缺陷／缺口今天还在树上 | 把票面 AC 里**最硬的一条**自己跑尺：`git --no-pager grep -n '<符号或那句字面>' HEAD -- internal/ cmd/ scripts/` 的命中数＋file:line；或 `git ls-files '<那个包>/*_test.go'` 之类存在性尺 |
| **已失效／已被别人做掉** | 票面要的东西已经在树上 | ⛔ 硬门：一枚 commit 号＋`git log --name-only` 命中，**或**今树读到实现件在场的复算尺；拿不出就不许写这一档 |
| **差翻勾** | 东西已在，但凭据是"票内注记／台账 `A##`／裁决表" | 凭据种类写进复法列；⚠ 不等于本腿裁它可结案，翻勾归编排者 |
| **量不到** | 判据落运行期行为／真麦克风／真 TTS 播放／真开窗／真机双击／CI 侧读数／`frontend/**`·`design/**`（禁读） | 归口写清**缺哪一行读数**；⛔ 不许由 grep 命中外推成"这功能能用" |

★ 这批是 9 月开的首轮切片票，写的时候产品还不存在 ⇒ **"票面点名的东西今天存不存在"只由尺回答，不由名字像不像回答**。
★ 切片级验收票（票 16／票 25）的射程盖的是**别的票的产物**：本腿不替它判"整片通过"，只回答"它点名的那几枚产物今天在不在树上"，其余进 §4。
★ `上次动过` 列＝`git log -1 --format='%h %ad' --date=format:'%m-%d' -- <票文件>`，本腿逐枚实跑。

## §1 名册（8 枚，票面标题原文）

分母尺（本腿 10-06 13:1x 现量，逐枚精确 glob，不含同前缀他号）：

```
for f in .scratch/wisp/issues/15-*.md .scratch/wisp/issues/16-*.md .scratch/wisp/issues/21-*.md .scratch/wisp/issues/22-*.md .scratch/wisp/issues/23-*.md .scratch/wisp/issues/24-*.md .scratch/wisp/issues/25-*.md .scratch/wisp/issues/26-*.md; do
  b=$(grep -c '^- \[ \]' "$f"); x=$(grep -c '^- \[x\]' "$f");
  echo "== $(basename $f) 未勾=$b 已勾=$x"; done
```

读数原文：

```
== 15-speech-engines-cer-harness.md 未勾=6 已勾=0
== 16-s2-acceptance.md 未勾=5 已勾=0
== 21-approval-gates-minimal.md 未勾=5 已勾=0
== 22-web-tools-d30.md 未勾=6 已勾=0
== 23-system-window-input-tools.md 未勾=5 已勾=0
== 24-doc-search-tools.md 未勾=5 已勾=0
== 25-s3-acceptance.md 未勾=5 已勾=0
== 26-tts-output.md 未勾=7 已勾=0
```

⇒ 8 枚**已勾全部＝0**，与编排者 10-06 13:1x 的现量一致；本腿无"已被别人翻过勾"要具名跳过的票。

| 票号 | 票文件（`.scratch/wisp/issues/`） | 票面标题原文 | 未勾格数 | 已勾 |
|---|---|---|---|---|
| 15 | `15-speech-engines-cer-harness.md` | 15 — Speech engines: VAD + streaming ASR via sherpa, engine slot mutex, CER harness | 6 | 0 |
| 16 | `16-s2-acceptance.md` | 16 — S2 acceptance: voice end-to-end gate, latency segments, hotplug pre-mortems | 5 | 0 |
| 21 | `21-approval-gates-minimal.md` | 21 — Approval gates minimal: L1 block window, L2 native card, queue trivial, batch D45-1 | 5 | 0 |
| 22 | `22-web-tools-d30.md` | 22 — Web & open tools: web.search/fetch/open, file.open, app.launch, D30 five layers | 6 | 0 |
| 23 | `23-system-window-input-tools.md` | 23 — System/window/input tool family: system.*, window.*, clipboard, media, screen, input.type | 5 | 0 |
| 24 | `24-doc-search-tools.md` | 24 — doc.read + local search tools: PDF/docx extraction, search.files/content/apps | 5 | 0 |
| 25 | `25-s3-acceptance.md` | 25 — S3 acceptance: capability scenarios ①② + full security red-team suite | 5 | 0 |
| 26 | `26-tts-output.md` | 26 — TTS output: sherpa matcha engine, serial slot, Speaking pipeline, P7 gate | 7 | 0 |

## §2 判档表（8 行 × 4 列）

| 票号 | 判定 | 复法（命令原文＋读数） | 备注／归口 |
|---|---|---|---|
| **15** | **仍成立**（六格整片零落地：`internal/speech` 今天只有包边界件） | 尺1 存在性：`git ls-tree -r --name-only HEAD internal/speech/ \| wc -l` → `1`，唯一文件 `internal/speech/doc.go`（同尺 `internal/audio/` → `16`）。<br>尺2 硬 AC「CER harness 基线集 COMMITTED 到仓」：`git --no-pager grep -n '\bCER\b' HEAD -- internal/ cmd/ scripts/ tools/ .github/` → **1 命中＝`internal/speech/doc.go:19`**（那行是 DEFERRED 自陈「implemented by ticket 15 (ASR + CER harness)」，不是实现）；`git ls-tree -r --name-only HEAD \| grep -ic 'asr-baseline'` → `0`（`near_clean`／`noisy_far` 同名尺亦 0）。<br>尺3 引擎符号：`git --no-pager grep -n 'AsrEngine' HEAD -- internal/ cmd/ scripts/ tools/` → 1 命中＝`internal/speech/doc.go:1`；`WakeWordEngine` → 1 命中＝`doc.go:2`；`TtsEngine` → 1 命中＝`doc.go:2`（三枚全部只在 doc.go 的注释里被命名）。<br>尺4 词面非零但要读落点：`git --no-pager grep -in 'paraformer' HEAD -- internal/ cmd/ scripts/ tools/` → 7 行，逐行落 `internal/config/manager_test.go:119,127,128`（配置键）与 `internal/models/manifest_real_test.go:44`（模型名 `asr-streaming-paraformer-zh-en`）；`silero` → 15 行同样只落 config/models；`git --no-pager grep -n 'intra_op_num_threads\|IntraOpNumThreads' HEAD -- internal/ cmd/ scripts/ tools/` → 2 命中＝`scripts/spike/speech-baseline/main.go:133`、`scripts/spike/xy-verdict/session_cgo.go:39`（**两枚都在 spike 目录，不在产品码里**）。<br>尺5：`git --no-pager grep -n '没听清' HEAD -- internal/ cmd/` → 0 行（`git grep -ic` 复核仍 0）。<br>尺6 唯一在场件：`git ls-tree -r --name-only HEAD internal/audio/` 含 `wavinjector.go`（C8 wav 注入缝），`internal/models/` 26 枚含 `manifest.go`／`downloader.go`（票 14 分发）。 | 逐格归口：①VAD／ASR 引擎、③端点化、④引擎槽互斥、⑤静音重试＋`asr` 错误类、⑥单核 CPU ⇒ **仍成立**（尺1／2／3／5 直接读数）。②CER 两档基线集那格**加一层量不到**：CI 侧读数本腿拿不到（⛔ 零 go、`.github/workflows/**` 只按行取），且基线集在仓尺已给 0。硬门尺是尺1＋尺2（今树实现件不在场），不是"名字像不像"。 |
| **16** | **仍成立**（切片级验收票；本腿只答"它点名的那几枚产物今天在不在树上"，不判整片） | 尺1 硬 AC「延迟表逐行填实测，提交到 `docs/SLO.md` 附录」：`git --no-pager grep -in '其余行（唤醒' HEAD -- docs/SLO.md` → `docs/SLO.md:121`，读数原文「\| 其余行（唤醒→监听、VAD→ASR、标点、首 token、冷/热会话整链） \| — \| 未测 \| 待票 12/15/26/28（各链路落地时逐行回填本表） \|」。<br>旁尺（表里确有实测行，但不是本票点名的语音分段）：`git --no-pager grep -in '16\.3\.4' HEAD -- docs/SLO.md` → `:4`／`:92`／`:101`（ASR↔TTS 串行切换，**FAIL 超 3 倍**）／`:113`／`:117`（面板冷拉起 PASS）。<br>尺2 六条预演证据件：`git ls-tree -r --name-only HEAD docs/ \| grep -ic 'premortem\|pre-mortem'` → `0`；`git --no-pager grep -ln 'pre-mortem' HEAD -- docs/ internal/ cmd/` → 3 命中＝`docs/PLAN.md`、`docs/evidence/s2/13-adversarial-acceptance.md`、`docs/specs/SPEC-10-testing-acceptance.md`（规格文字，无 S2 执行记录）。<br>尺3 真机 E2E 演示记录：`git ls-tree -r --name-only HEAD docs/evidence/s2/` → **只有 2 枚**＝`13-adversarial-acceptance.md`／`14-adversarial-acceptance.md`（只按名取、⛔ 未读内容，且按名属票 13/14 不属本票）。<br>尺4 热插拔件：`git --no-pager grep -c 'func Test' HEAD -- internal/audio/hotplug_test.go` → `8`（件在场，⛔ 运行期读数不在本腿射程）。 | 它点名的产物：延迟分段回填行（`:121` 自陈未测）／六条预演证据（0 枚）／演示记录（0 枚属本票）⇒ 三样都不在 ⇒ 仍成立。AC#4「15 的 CER 门在 CI 绿」随票 15 落**量不到**；AC#5「另一枚 agent 的对抗验收」＝裁决表凭据种类，本腿**不判**（只按名取表，判"过没过"越界，见 §4 第 4 条）。台账里唯一具名本票的一行不是完成凭据：`git --no-pager grep -n '票 16' HEAD -- docs/reports/pending-and-issues.md` → `:97 [B3] H2/H5 登记补强：票 16 真机三场景（原 H6 漏登 registry）已由本条补登`——它登记的是**待办**，与本档判定同向。 |
| **21** | **仍成立**（AC#2「卡片显示 FULL params（不是摘要）」这一条今天还缺；其余 AC 大面积已在，但票内注记之外还压着 10-05 未销的 BLOCKER） | 尺1（本档的硬门尺＝今树直接读数）：`git --no-pager show HEAD:internal/agent/approval/replies.go \| grep -n 'Params'` → **0 命中**；`ReplyCard` 字段面（`replies.go:70-90`）只有 `CorrelationID/TaskID/Tool/Level/Grant/Paths/Window/Deadline`。`git --no-pager grep -n 'Params' HEAD -- cmd/wisp/approval_reply.go cmd/wisp/resident_approval_windows.go cmd/wisp/panel_pump.go internal/ball/` → **0 命中**。`git --no-pager grep -n '\.Params' HEAD -- internal/ cmd/ \| grep -v '_test\.go'` → 3 命中＝`internal/agent/approval/gate.go:603`、`internal/agent/approval/pending_read.go:68`、`internal/tools/bridge.go:292`，三处全是"往卡／投影里填"，无一处渲染。读数：**装 Params 的动作有、把 Params 印到任何面上的动作零** ⇒ 票面那句 "FULL params (not summary)" 仍在树上。同档旁证（今树自己的自陈，`HEAD:` 尺）：`internal/agent/approval/pending_read.go:44-47` 原文写着「the four the L2 card actually carries (Level/RulesHit/Reason/SessionOverrideBlocked)」——**Params 不在它列的那四枚里**，与本尺读数同向。<br>尺2（其余四格已在的在场尺）：`git ls-tree -r --name-only HEAD internal/agent/approval/ \| wc -l` → `29`（其中 `_test.go` `19`）。commit 凭据：`git --no-pager log --format='%h %ad %s' --date=format:'%m-%d' -1 35200c7` → `35200c75 09-20 feat(21): the L1 block window, its four veto channels and the gate that owns them`；同尺 `--name-only` 命中 `internal/agent/approval/{approval.go,batch.go,gate.go,queue.go,report.go,ui.go,window_test.go}`。<br>尺3 生产可达（票内 09-20 自陈"零可达"那一条，今天已翻）：`git --no-pager grep -n 'rt.gate = approval.New' HEAD -- cmd/wisp/run.go` → `run.go:612`；同文件 `:748 Gate: rt.gate`／`:749 Cancel: rt.gate.ToolsCancelBus()`／`:1008 AdmitTask: rt.admitTask`；`:1344 consoleApprovalUI.Prompt` 印 Level/Tool/Reason/RulesHit/Channels/Paths/Window/CorrelationID；`cmd/wisp/resident_approval_windows.go:864 ballCardUI.Prompt` 真挂球＋`TakeEscForCancel`。<br>尺4 判定层逐 AC 点名：`approval/window_test.go:20/55/84/207/264`、`queue_test.go:89 TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis`、`batch_test.go:19/65/87/120/142`、`gate.go:751 Replay`、`queue.go:122 DefaultApprovalTimeout = 300 * time.Second`、`:143 WarningLead`（last-30s）、`approval.go:227 return "语音取消不可用"`、`statemachine/table.go:118/123/150/154/158/163/168`（D43 #17/21/22/23/24）、`internal/ball/statevisual.go:35 RingPulse`／`:177 case statemachine.StateConfirming`／`:184 BadgeCount = 1`、`cmd/wisp/panel_pump.go:70 RulesHit`／`:256 panelArgs`。 | **凭据种类三样齐全**（⚠ 不等于本腿裁它可结案，翻勾归编排者）：①票内注记（09-20 段 1；`Status:` 行按口径不采信）②台账 `A##`：`git --no-pager grep -n '票 21' HEAD -- docs/reports/pending-and-issues.md` → `:334 [A13]`、`:428`、`:469 [A19] 审批层有判定、没有输入设备`、`:482 [A20] D45-1 批量聚合站错了轴，且对全部在产工具恒为死码`、`:498 [A21] 两条 MINOR`、`:617` 归"票 21 段 2" ③裁决表（只按名取）＝`docs/evidence/s1/21-segment1-adversarial-acceptance.md`。<br>两半互不覆盖：判定／队列／聚合／nonce／D43 表／球视觉今天都在树且有 commit ⇒ 只看这些它像〔差翻勾〕；但尺1 是今树零读者的直接读数、A19/A20 到 10-05 未销 ⇒ 本腿写〔仍成立〕，**不写"整片已落地"也不写"整片还欠"**。<br>**量不到的那一整半截**：原生卡片实际观感、球脉冲／深度徽章渲染、真机 Esc／单击否决（`frontend/**`·`design/**` 禁读；`cmd/wisp/resident_approval_live_246_windows_test.go:72/195/289` 三枚 `TestLive246*` 要真开窗）。缺的那一行读数＝真机上一次 veto 的窗口读数。 |
| 22 | 尚未判，本腿下一步判它（票号升序第 4 枚） | 判档时逐 AC 跑尺 | — |
| 23 | 尚未判，本腿下一步判它（票号升序第 5 枚） | 判档时逐 AC 跑尺 | — |
| 24 | 尚未判，本腿下一步判它（票号升序第 6 枚） | 判档时逐 AC 跑尺 | — |
| 25 | 尚未判，本腿下一步判它（票号升序第 7 枚，切片级验收票） | 判档时只答产物在不在 | — |
| 26 | 尚未判，本腿下一步判它（票号升序第 8 枚） | 判档时逐 AC 跑尺 | — |

## §3 判得心虚的枚数与具名理由

本节当前**尚未写**，本腿判完 §2 之后按票号逐枚写满；⛔ 不留空。

## §4 一处尺有歧义／判不动的地方

本节当前**尚未写**，与 §3 同步补满；已预登记两处，届时展开并补其余。

- 预登记①：票 15/26 的 `internal/speech` 存在性尺——本腿 13:2x 现量 `git ls-tree -r --name-only HEAD internal/speech/ | wc -l` ＝ `1`，唯一文件 `internal/speech/doc.go`；同尺 `internal/audio/` ＝ `16` 枚。旧探针读数一律不采信。
- 预登记②：票 21/23/26 的"界面半边"判据落 `frontend/**`·`design/**`（本腿禁读）＋真机行为，按边界必须落〔量不到〕，不许绕。
