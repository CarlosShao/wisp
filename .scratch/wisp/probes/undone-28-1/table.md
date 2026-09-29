# undone-28-1 —— 结案票残余 28 枚未勾判据逐枚判性表

- 盘点腿代号：`undone-28-1`（只读盘点，不产码、不改票面、不翻勾）
- 锚点：本腿现跑 `git rev-parse --short HEAD` = `fac60ad4`，分支 `dev`
- 尺：`grep -c -- '^- \[ \]' <file>`（⛔ 不用 `grep '^- \['`，进度戳 `- [2026-09-29 17:5x +08] ...` 会造成约 1,600 枚假阳性）
- 分母复核（本腿现跑，`.scratch/wisp/issues/`）：78 张 `*-done.md`，未勾合计 **28 枚／13 张票**，
  逐票计数＝`92(4) 07(4) 80(3) 115(3) 97(2) 113(2) 110(2) 105(2) 104(2) 89(1) 63(1) 118(1) 11(1)`
  ⇒ **与编排者 09-29 17:5x 的 28 枚／13 张票完全一致**（逐票分布也一致）
- 四分类定义（本表用字）：甲＝残余真活 / 乙＝措辞残留（不是判据） / 丙＝已归口别的票 / 丁＝今天其实已满足、只是没人翻勾

## 0. 总账表

> 结论列字面：**甲**＝残余真活 / **乙**＝措辞残留（不是判据）/ **丙**＝已归口别的票 / **丁**＝今天已满足没人翻勾 /
> **判不了**＝凭据只能在禁读地界（`frontend/**`、`design/**`）里证。
> 合计：**甲 13／丙 12／乙 1／丁 1／判不了 1 ＝ 28**（与分母对齐，无缺口）。

| # | 票号 | 行号 | 判据名（原文行首片段） | 结论 | 凭据出处（本腿现读） |
|---|---|---|---|---|---|
| 1 | 07 | `:51` | Visual: 20 states rendered in a debug cycle page | **甲** | `docs/evidence/s1/07-adversarial-acceptance.md:10`＋`:14`；台账 `docs/reports/pending-and-issues.md:7`（H1 未闭） |
| 2 | 07 | `:54` | Interactive: click ball → Sleeping→Listening | **丙** | 票 64 `:72`（`- [x]`，64 开放 un=1）；`07-…md:27` FAIL |
| 3 | 07 | `:59` | Hotkeys registered/re-registered on config change | **丙** | 票 64 `:44`/`:59`/`:60`；"mute→Muted"半枚只有 64 `:152` 自述＋`internal/ball/hotkey_live_test.go:225` |
| 4 | 07 | `:61` | Multi-monitor: drag to second monitor | **丙** | 票 64 `:83`（`- [ ]`，硬件缺席）＋`:234`「本票不能置 done」 |
| 5 | 11 | `:53` | Probe suite: mockllm control endpoints | **丙** | 移交 commit `9293518d`（"AC#6 transferred as A11"）；票 12 `:120`/`:123`/`:128`；`internal/llm/probe_health.go:14`、`cmd/wisp/providers.go:189/:195`、`providers_test.go:143` |
| 6 | 63 | `:46` | 端到端：`api_key_ref = "secret:<id>"` 被 provider 解析 | **丙** | 移交 commit `25d5c3ca`（"registered as A8"）；票 12 `:109`/`:112`/`:117-119`；`cmd/wisp/run_test.go:420`/`:446` |
| 7 | 80 | `:43` | AC#3 接线落地 + 真跑在 `risk.Gate` 上的用例 | **丙** | 票 230 `:22`（AC#3 交接要么落地要么作废）＋`:23`（(a) 半枚"本票不接"，待人工）；`internal/tools/mode.go:98` 有生产调用点、`internal/config/unwired.go:81` 说 map 无人填 |
| 8 | 80 | `:50` | AC#4 变异双向 | **丙** | 票 230 `:22` 同一格（"把票 80 那两格转进本票"） |
| 9 | 80 | `:54` | AC#5 门禁（只跑自己碰到的包） | **丙** | 票 83 `:57`（`- [x]` 同族门禁）；⚠ 83 那格包名只写 `internal/config` |
| 10 | 89 | `:348` | 待办：AC#4／AC#5／三条落盘路径／AC#6 POSIX | **乙** | 位置在 Progress log；四条 item 分别在 89 `:363`/`:352`/`:382`/`:392` 已勾；残余同族清单已归口票 132 `:58` 起"编排者增量" |
| 11 | 92 | `:52` | AC#1 面板只读显示档位 | **丙** | 票 114 `:50`（AC#5 不越界，开放 un=9）＋`:47`/`:67`；钉子本体已在 `internal/panel/composer_test.go:93`；`92b-…md:134` 判不通过 |
| 12 | 92 | `:68` | AC#5 变异三向 | **甲** | `docs/evidence/s1/92b-adversarial-acceptance.md:137`：(iii) 未重做、〔仅自述，不背书〕；(i) 红+绿各半 |
| 13 | 92 | `:73` | AC#6 台账与门禁 | **甲** | `92-…md:147` FAIL → `92b-…md:138` 通过（绑 `91b5fc4`/`a8f9459`）；`R-92b-4`（`:290`）未销；当前 HEAD 无读数 |
| 14 | 92 | `:79` | AC#7 负判据（不做 git 切换） | **判不了** | 整格需 `frontend/` 一面（禁读地界）；已验半边：`internal/panel/composer_test.go:269-270`＋本腿负向 grep 0 命中＋`92b-…md:139` 通过 |
| 15 | 97 | `:53` | AC#4 注释与代码一致 | **丙** | 票 230 `:21`（AC#2 那句注释改成实话）；本腿现读 `internal/agent/approval/queue.go:247` 仍写"only reader"，`q.alias` 命中 `:197/:200/:210/:216/:257`；`97-…md:34` R-97-1 |
| 16 | 97 | `:55` | AC#5 门禁（按包） | **甲** | `97-…md:17`「通过」绑 `f140079` 快照；`R-97-3`（`:36`）未销；当前 HEAD 无读数 |
| 17 | 104 | `:52` | AC#3 双向变异 | **甲** | `104-…md:115` 三腿变异〔独立复现〕绑 `/tmp` 快照；104 票面 `:68` 编排者注记同族用例曾在 CI step4 红 |
| 18 | 104 | `:57` | AC#4 回归（四包四数＋三门） | **甲** | `104-…md:162`「每个数字都是我在 `/tmp/ac104-fk`（`4d43447` 纯净树）本机重跑的」；当前 HEAD 无读数 |
| 19 | 105 | `:47` | AC#3 祖先重解析那条腿的行为用例 | **丁** | 票 116 已结案（un=0／chk=6）＋表 `docs/evidence/s1/116-adversarial-acceptance.md:108` 通过并逐字「**票 105 能否因此结案：能**」；用例在 `internal/risk/syncdirs_ancestor_actable_leg_116_test.go:153/:239/:270/:307`；⚠ 读数绑 `80e248c` |
| 20 | 105 | `:54` | AC#5 门禁（按包 scope） | **甲** | `105-…md:151` 口径自陈跑在 `f6818f2` 纯净快照；当前 HEAD 无读数 |
| 21 | 110 | `:39` | AC#3 新步自己会红（两发变异） | **甲** | `110-…md:16`「通过（三发全部独立复现）」绑 `/tmp/wisp-ac110-snap`；仪器在位 `.github/workflows/ci.yml:398`＋`scripts/winsec-tests.sh` |
| 22 | 110 | `:47` | AC#5 门禁 | **甲** | `110-…md:18`「通过（附一条后果登记）」；后果开放中：票 111 un=10／112 un=5／140 un=4；commit `f161c690` 标题自陈"把后面的步骤全吃掉了" |
| 23 | 113 | `:54` | AC#4 变异（关腿＋半修） | **甲** | `113-…md:16`「四发全部我自己下刀、自己复量」绑容器快照；当前 HEAD 无读数；116 接的是另一族腿（`internal/risk`） |
| 24 | 113 | `:56` | AC#5 门禁（容器三包四数＋三门） | **甲** | `113-…md:17`「绿（四数与三门），但 ban#8 那格的 sha↔数字配对错了一格」⇒ 当年就不干净 |
| 25 | 115 | `:44` | AC#1 独立复算四条红＋对质旧 run | **丙** | `docs/evidence/s1` 目录里 `115-*.md` **零命中**（115 名下无表）；票 230 `:19` AC#1 具名接走"115 要一张真裁决表"，230 开放 un=5 |
| 26 | 115 | `:58` | AC#4 变异（退回原状＋只 `EqualFold`） | **甲** | 承载用例在位 `internal/winsec/notice_attribution_115_windows_test.go:191/:240`；无任何变异读数、无非实现者表 |
| 27 | 115 | `:64` | AC#6 门禁（五发） | **甲** | 数只在票面 Progress log `:251`/`:360`/`:444`（实现方自述）；当前 HEAD 无读数 |
| 28 | 118 | `:59` | AC#8 跨卷归属：先量再判 | **丙** | 票面 `:63` 自写"停手回报＋另立票"；票 126 `-done`（un=0）＋表 `docs/evidence/s1/126-adversarial-acceptance.md:40-44` 全通过；生产码今天已修 `internal/winsec/resolve.go:455` `sameVolume` |

## 1. 票 07 —— `07-ball-state-machine-core-done.md`（4 枚）

> 本腿现量：`grep -c -- '^- \[ \]'` = 4，行号 `:51 :54 :59 :61`（未勾）；已勾两枚在 `:49`、`:57`
> （`grep -n -- '^- \[x\]'` 现读）。
> 票面里 `done-fix-1` 已给每枚留了追加段；下面四枚的**判性由本腿独立现读**，只把那些段落当线索不当凭据。

### 1.1 `:51` Visual：20 态渲染 + 人工截图对照 —— **甲（残余真活·阻在 owner 一次眼）**

- 原文行首逐字：`- [ ] Visual: 20 states rendered in a debug cycle page/window; human screenshot review vs`
- 缺什么，说得出口：**owner 本人的一次主观签收**（20 态好不好看，测试绿不代替）。凭据（全部现读、非禁读地界）：
  - 裁决表 `docs/evidence/s1/07-adversarial-acceptance.md:10` 逐字「| 3 | 视觉证据 | **PASS（人工签收挂起）** | docs/evidence/s1/ball-states/ 截图 + c21-native-tokens.md 对照表 |」
  - 同文件 `:14` 逐字「**待人工项**：20 态视觉的主观签收（**用户本人**）→ 已登记 docs/reports/pending-human-review.md」
  - 台账 `docs/reports/pending-and-issues.md:7` 逐字「- [H1] **悬浮球 20 态视觉签收** — **2026-09-20 已改为实况签收**」，
    且 `:9-10` 写「满意则关闭；不满意提修改意见转新工单」⇒ **H1 至今没有关闭行**（本腿 `grep -n 'H1'` 全量扫过台账，
    没有任何一处把这条待人项记成已闭）⇒ 这一格等的不是读数，是一次人眼。
- ⚠ 本腿能判到这儿为止：本格**第二个分句**要对照 `design/screens/ball.html`，那是派单写死的**禁读地界** ⇒
  本腿**永远无法**把这枚升成丁（无法证明对照已发生）。裁决者若要判丁，得由能读 `design/**` 的人补一次核。
- ⚠ 顺带一枚引用缺陷（不是本票的事，但影响后来人找凭据）：裁决表 `:14` 指向的
  **`docs/reports/pending-human-review.md` 在树里不存在**（`git ls-files | grep -c 'pending-human-review'` = 0，
  `find docs -name 'pending-human-review*'` 空）。内容实际在 `pending-and-issues.md` 的
  「## 待人工审核（pending-human-review）」这一节（`:5` 起）。
- 归口核查：票 65（`65-ball-glass-quality-rework.md`，本腿现量 un=6、无 `-done` 后缀＝**仍开放**）
  接的是"玻璃质感返工"，其判据 `:57`-`:63` 里没有一枚写"20 态 debug cycle 页 / 对 `ball.html` 的 20 态签收"；
  票 68（`68-ball-default-visuals-parity.md`，un=2、Status: review）接的是默认值翻转与差分复测，也不同事。
  ⇒ **不是丙**（没被整枚接走），保持甲。

### 1.2 `:54` Interactive：单击唤起 / Esc 取消 / 不抢焦点 / 透明区穿透 —— **丙（整枚归口票 64）**

- 原文行首逐字：`- [ ] Interactive: click ball → Sleeping→Listening; Esc/click cancels from Confirming; focus`
- 归口去处**点名并现读**：`.scratch/wisp/issues/64-ball-defects-hotkey-interactive.md`（本腿现量 un=1／chk=8，
  文件名**不带** `-done` ⇒ 仍开放）`:72` 是一枚 **`- [x]`**，逐字
  「交互四项各有真机测试且**不被 skip**：单击唤起、Confirming 取消、不抢焦点、透明区穿透」——四个分句与本格一一对上。
- 07 当年确实没证据（不是本腿臆断）：`docs/evidence/s1/07-adversarial-acceptance.md:27` 裁决
  「AC#3 交互…| **FAIL** | 唯一被点名的候选 `internal/ball/tokens_test.go:239 TestHitTestAndDPIInjection`
  只覆盖**第 4 个分句的一半**」。
- ⚠ 但 64 的这枚勾本身带条件：64 `:72` 之下的票面文字自陈「今天的判定是"结构成立（无静默 skip）、逐跑记录缺证"」
  并要求「交互四项任一走 SKIP（…）就把本框退回未勾」（两句本腿现读＝64 `:78` 与 `:77`）。本腿**禁跑** `go test -tags winlive`，
  因此既不能确认 64 那枚勾今天仍可复算，也不据此判 07 这枚为丁。⇒ **丙：账在 64 名下，07 这格不该再占本票分母。**

### 1.3 `:59` Hotkeys 注册/重注册 + 静音切 Muted —— **丙（大半归口票 64；"静音"半枚只有实现方记录）**

- 原文行首逐字：`- [ ] Hotkeys registered/re-registered on config change; mute toggles Muted state.`
- 归口去处**现读**（票 64 文件同上）：
  - `:44` `- [x] 改 [hotkey] 后热键真的重注册（端到端：改配置→按新键→球响应；旧键不再响应）。`
  - `:59` `- [x] 注册失败两类语义分开：被占用 vs 未尝试，各自给出用户可见提示（不再只有一行 Warn）。`
  - `:60` `- [x] 默认唤起键为 Ctrl+Alt+Q，且 Ctrl+Alt+Space 作为可配置备选写进文档`
  - 07 侧当年同样是 FAIL：`docs/evidence/s1/07-adversarial-acceptance.md:62`
    「AC#5 热键注册 / 配置变更后重注册；静音切换 Muted 态 | **FAIL** | 被点名候选 … 与 AC#5 语义无关」
- ⚠ **"mute toggles Muted"那半枚没有被 64 的任何判据框逐字接住**：本腿现跑 `grep -n 'Muted' 64-*.md` 只命中
  一处，是 Progress log `:152`「`WM_HOTKEY(mute)` → `OnMuteHotkey` → `EvMuteKey` → 机器 `Muted` 且球渲染 `Muted`；再按一次」
  ＝**有逐跑记录、没有判据框**。盘上另有码侧对应物（`internal/ball/hotkey_live_test.go:225`
  「TestLiveMuteHotkeyEndToEnd is A1d: the registered mute key -> OnMuteHotkey -> EvMuteKey -> the machine and
  the ball both in Muted, and back out.」，本腿现读；该文件是 `winlive` 门控腿，本腿不跑）。
  ⇒ 该半枚的"证据"目前**只有实现方/编排者自己的逐跑叙述**，属〔仅自述〕，**不能**据以翻 07 的勾。
  结论仍记丙，但**请裁决者把"64 要不要为静音半枚单立一格"当成一个待决问题**。

### 1.4 `:61` Multi-monitor：拖到第二屏/持久化/恢复/模拟拔出 —— **丙（归口票 64，且 64 那一格至今未勾）**

- 原文行首逐字：`- [ ] Multi-monitor: drag to second monitor, persist, restore; simulated detach → primary.`
- 归口去处**现读**：`64-ball-defects-hotkey-interactive.md:83` 是一枚 **`- [ ]`**，逐字
  「多显示器：真拖到第二屏验证并留证；本机无第二屏则**保持未勾**并写明所需硬件」，其下一行写明
  「本机物理单屏（`\\.\DISPLAY4` 3440x1440），**硬件缺席**」。64 `:234` 附近另写「本票**不能**置 done」。
- 07 侧裁决：`docs/evidence/s1/07-adversarial-acceptance.md:76`「AC#6 多显示器…**PARTIAL**｜已真正验证的两段：
  ①"模拟拔出 → 主屏"是纯函数级真断言且在默认套件内」。
- ⇒ **丙**：活还在（缺一次副屏硬件跑），但账**已明确挂在开放票 64 名下**。
  ⚠ 本腿补一句给裁决者：这枚是「丙」而**不等于"已交付"**——翻 07 这框＝把别人名下还欠着的一次硬件跑算成完成，
  所以判性给丙（分母归 64）之外，**不能**同时当丁用。

## 2. 票 11 —— `11-llm-adapters-rest-done.md`（1 枚）

> 本腿现量：un = 1（`:53`）；票面 AC 框共 7 枚（`^- \[x\]`=6 ＋ `^- \[ \]`=1，两条命令现跑，非心算）。

### 2.1 `:53` Probe suite（fc-capable / fc-broken / vision-capable → provider_health + mismatch 事件）—— **丙（正式移交票 12，代号 A11，且接收方已闭环）**

- 原文行首逐字：`- [ ] Probe suite: mockllm control endpoints emulate fc-capable / fc-broken / vision-capable`
- **移交是有 git 记录的，不是转述**：本腿现跑 `git show -s --format=%B 9293518d`，标题逐字
  「docs(acceptance): ticket 11 DONE with a 1:1 verdict table; **AC#6 transferred as A11, ruling R12**」，
  正文末段逐字「AC#6 is PARTIAL, not failed: the measuring half is proven, but nothing in cmd/wisp …」。
  票面自身 `:66`（09-20 10:20Z，`agent-ticket11-probe-2`）也写「LEFT AC#6 UNTICKED on purpose」＋三条原因
  （装配根没人调 / thinking 沿用票 09 的检查 / audio 按契约推票 61）。
- 归口去处**本腿现读**：`.scratch/wisp/issues/12-cli-text-path-s1-gate.md`
  （⚠ 该票文件名**不带** `-done`＝仍开放，本腿现量 un=3：`:48` `:148` `:149`）
  - `:120` `- [x] The capability probe is actually called from the composition root: `cmd/wisp`'s provider`
  - `:123` 逐字「event can fire on a real machine. Handed here from ticket 11 AC#6, registered as **A11**.」
  - `:128` 逐字「—— **A11 就此闭环**」，点名 `cmd/wisp/providers_test.go` 的
    `TestProvidersProbeRecordsMeasuredThinkingTrue` / `...False`（正反两向）、
    `TestProvidersDiscoverListsWhatTheServerServes`、`TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless`
- 盘上凭据（本腿现读，全在允许地界）：
  - `internal/llm/probe_health.go:14` 逐字「// Ticket 11 AC#6, consumer side: turn probe.go's per-capability PRIMITIVE」
  - `cmd/wisp/providers.go:189` `var mismatches []llm.ProbeMismatch` ／ `:195` `Mismatch: func(m llm.ProbeMismatch) {`
  - `cmd/wisp/providers_test.go:143` 逐字「t.Errorf("the declared-vs-measured event was not announced:\n%s", errb)」
- 凭据出处是**非实现者**：`9293518d` 正文自陈「Three claims I refused to take on the implementer's word」＋
  自己下的变异（`probe_health.go:255` config-echo 变异弄红三枚探针用例），执行者是编排者验收方，不是票 11 的实现代理。
- ⇒ 判**丙**（分母应归票 12 名下）。⚠ 与丁的分界本腿说清：这枚**事实上**已具备翻勾条件，
  但它在 git 与票面上被**正式移交**走了；若记成丁，等于同一份活在 11/12 两张票里各算一次交付。
  建议处置＝票 11 这格改指向 A11（或直接搬出勾框），**而不是翻勾**。
- ⚠ 一枚本腿不替它认账的边角：原格还含 "vision-capable" 一档与 `✓/✗` 呈现，票 12 的 A11 块写的是
  fc/thinking 两向；`grep -n 'vision' 12-…md` 的射程本腿没逐格核（超出"这枚未勾判据判性"所需），
  若裁决者要按丁翻勾，得先补这一眼。

## 3. 票 63 —— `63-credential-entry-cli-done.md`（1 枚）

> 本腿现量：un = 1（`:46`）；AC 框共 7 枚（chk=6 ＋ un=1，两条 `grep -c` 现跑）。

### 3.1 `:46` 端到端：`api_key_ref = "secret:<id>"` 被 provider 解析出发请求 —— **丙（正式移交票 12，代号 A8，接收方已闭环且出自非实现者）**

- 原文行首逐字：`- [ ] 端到端：`wisp secret set` 存好后，config 里 `api_key_ref = "secret:<id>"` 能被 provider 解析出发请求（用 mockllm，测试里用假 key）。`
- 当年为什么不勾，票面自己写了（`:47-49` 逐字）：`TestSecretEndToEndConfigRefResolvesAtRequestTime` 证明了
  ref→`config.LoadFile`→`ProviderKeys`→`Authorization` 头真上线，但「发请求的是测试自己的 http client，
  不是 `internal/llm` 的 provider；provider 那一段是调用方自证的」。⇒ 缺的东西很具体：**由 provider 自己发出去的那一发**。
- 移交是 git 记录，不是转述：本腿现跑
  `git log -S'Handed here from ticket 63' -- .scratch/wisp/issues/12-cli-text-path-s1-gate.md`
  → `25d5c3ca|CarlosShao|2026-09-20|docs(tickets): hand ticket 63's unticked AC#6 to ticket 12, registered as A8`
- 归口去处**本腿现读**（`12-cli-text-path-s1-gate.md`，无 `-done` 后缀＝仍开放，un=3）：
  - `:109` 是一枚 **`- [x]`**，逐字「Key comes from the store, not the config: the provider that serves the `wisp run` request is」
  - `:112` 逐字「(not one a test-written client set). **Handed here from ticket 63 AC#6**, whose end-to-end proves」
  - `:117-118` 逐字「—— **A8 就此闭环**：`cmd/wisp/run_test.go::TestRunTextTaskKeyResolvesInTheStore`（provider 自己
    发出请求、key 从 DPAPI 存储解析）+ `TestMissingBlobFailsUnconfiguredNeverSilently`（正是上面那条变异判据…）」
- 盘上凭据（本腿现读，允许地界）：`cmd/wisp/run_test.go:420` `func TestRunTextTaskKeyResolvesInTheStore(t *testing.T)`、
  `:446` `func TestMissingBlobFailsUnconfiguredNeverSilently(t *testing.T)`、`:25` 头注释逐字
  「TestRunTextTaskKeyResolvesInTheStore A8: the Authorization header the」
- 凭据出处＝**非实现者**：票 12 的实现代理 `agent-ticket12-assembly` 跑完 165 次工具调用后
  「票面 8 个 AC 框一个没勾、Progress log 一行没写」（台账式记录在 `12:162`，`agent=orchestrator`），
  A8 那段的闭环文字与「我亲自跑通（含在 ok 13.794s 那一批里）」（`12:119`）都是**编排者对账时补写并亲自跑的**。
- ⇒ 判**丙**。与丁的分界同 11：`-done` 票的这格已被具名移交，记丁会造成 63/12 双重计账；
  建议处置＝这格改成指向 A8 的引用（或搬出勾框），**不是翻勾**。
- ⚠ 附带一条给裁决者：票 63 的 `:46` 与票 12 的 `:109` 说的是**同一次端到端**，
  而票 12 那张 A8 框已经勾了 ⇒ 这类"移交＋接收方已勾"的残余格（本腿在票 11 已核到一枚同形状，
  其余各票本腿按同一规矩逐枚核，见下文对应节）
  **不该再回到分母**，只需在票面上留指针。

## 4. 票 80 —— `80-blacklist-overrides-never-wired-to-gate-done.md`（3 枚）

> 本腿现量：un = 3（`:43` `:50` `:54`）。票面头部逐字（`:4-8`）：「本票的交付物是**裁决**，不是码：**零 Go 改动**…
> 后续工作**全部转移到票 83**；AC#3/AC#4/AC#5 的接手方逐条写在框后」＋「编排者裁决 = 选项 (C)（见 A53②）」。

### 4.1 `:43` AC#3 接线落地 + 一条真跑在 `risk.Gate` 上的用例 —— **丙（账在票 230 AC#3；⚠ 上一轮腿用的那把尺是坏的，本腿重尺后结论不变）**

- 原文行首逐字：`- [ ] **AC#3** 接线落地 + 一条**真跑在 `risk.Gate` 上**的用例：同一份配置里写 override，`
- 框后自带「⇒ **未做（被 AC#2 挡住）**」，并写明 (a) 那半的形状＝`PLAN.md:2399`「必须由人点」禁止的静态预授权。
- 归口去处**本腿现读**（`.scratch/wisp/issues/230-four-cells-left-unfinished-inside-closed-tickets.md`，
  **无** `-done` 后缀＝仍开放，本腿现量 un=5：`:19 :21 :22 :23 :24`）
  - `:22` 逐字「**AC#3 票 80 的交接要么落地要么作废**：二选一——**甲＝把那格真写进票 21 的判据**…
    或**乙＝把票 80 那两格转进本票并具名作废原交接句**。⛔ 不许留着"接手方＝票 21"这句话继续存在而票 21 一字不知。」
  - `:23` 逐字「**AC#4 静态放行开关那一枚＝必须由人点（本票不接）**」⇒ (a) 那半**没有**被 230 接走做，
    只是摆给 owner（＝待人拍板，Q-27 线）。
- 悬空交接本腿自己复算，**成立**：`grep -l 'blacklist_overrides' .scratch/wisp/issues/*.md` 只命中
  **80／83／90／230** 四张；`grep -c -i 'override' 21-approval-gates-minimal.md` = **0**（票 21 一字不知）。
- ⛔ **推翻一枚上一轮腿的尺**：票 80 `:44` 的 `done-fix-1` 追加段引
  「`grep -rn 'NewGate(' --include=*.go internal cmd | grep -v _test` ＝ **0 命中**（真跑在 `risk.Gate` 上那条用例今天仍无承载体）」。
  本腿现读：`risk.Gate` 在本仓**不是构造函数**，`NewGate(` 这个串**全仓一处都没有**（含测试）——
  它的真身是 `internal/risk/blacklist.go:76` 的函数 `func Gate(canonical string, bOverrides map[string]bool) PathDecision`。
  **一把全仓零命中的尺证明不了任何事。**
- ⚠ 重尺后的**新事实**（对裁决者有用，不改变判性）：`risk.Gate` **今天已有生产调用点**——
  `internal/tools/mode.go:98` `d := risk.Gate(c, overrides)`，`overrides := b.confirmations()`（`:92`）。
  但填这张 map 的是**运行期 L2 确认**，不是配置键；这一点仓内守卫自己写着：
  `internal/config/unwired.go:81` 逐字「nothing populates the bOverrides map risk.Gate reads: ticket 90 gave Gate
  a production call site (internal/tools readBlacklist…), but the confirmations that would fill this map are
  minted only by an L2 answer, and that flow is ticket 21's approval queue」。
  ⇒ AC#3 缺的具体东西仍然具体：**"配置里的 override 被解析成 bOverrides 并翻转那一个文件"这条腿今天没接**，
  且 (a) 那半**卡在人工裁决**上。判丙（分母归 230），但**这是"丙而活未死"**：230 AC#3 收的是交接文书，不是这条功能。

### 4.2 `:50` AC#4 变异双向 —— **丙（与 AC#3 同批，账在票 230 AC#3）**

- 原文行首逐字：`- [ ] **AC#4** 变异双向：(i) 把接线断开 ⇒ AC#3 必须红；(ii) 把 (b) 的"其它文件仍 B"断言指向 override`
- 框后逐字「⇒ **未做**：AC#3 没有码可断，变异无承载体。」⇒ 本格**没有独立生命**，AC#3 落地才有它。
- 归口去处＝票 230 `:22` 同一格（其文字逐字写"把票 80 **那两格**转进本票"，那两格＝AC#3／AC#4）。
- ⚠ 记一句给裁决者（与 `pool-2` 同劝告）：**只处理 AC#3 不管本格**会留下"接线有牙没验过"的第二枚洞；
  两枚必须同批处置。

### 4.3 `:54` AC#5 门禁（只跑自己碰的包）—— **丙（交接写得最标准的一枚：账在票 83 `:57`，已勾）**

- 原文行首逐字：`- [ ] **AC#5** 门禁（**只跑自己碰到的包**，共树不跑整仓）：`gofmt -l <pkgs>` 空、`
- 归口去处**本腿现读**：`83-config-keys-that-lie-must-fail-loudly-done.md:57` 是一枚 **`- [x]`**，逐字
  「**AC#5** 门禁（只跑自己碰的包）：`gofmt -l` 空、`gofumpt -l <files>` 空、`go vet ./internal/config/` rc=0、
  `GOOS=linux go vet ./internal/config/` rc=0、`go test -count=2 ./internal/config/` rc=0，逐跑点名 `--- SKIP`/`--- FAIL`」，
  并且 `:60` 承接了票 80 量到的 `TestResolvePerCallBudget` 并跑抖动（「如实登记为外项，不许顺手调那个 1ms 预算」）。
- 接手句在票 80 框后逐字：「⇒ **编排者补（A53）**：接手方 = **票 83 的 AC#5**（同两个包，而票 83 会真改 Go 文件 ⇒ 门禁有承载体）」。
- ⚠ 一枚射程差，本腿如实报：83 那格的**包名只写 `./internal/config/`**，而票 80 的 AC#5 要求
  `internal/risk` + `internal/config` 两个包都过门。票 80 自称「零 Go 改动 ⇒ 门禁无改动面可验」，
  所以这格按丙处置没问题；但若裁决者想拿 83 那枚勾去翻 80 这格（＝丁路线），
  **`internal/risk` 那半边在 83 的判据文字里没有对应物**，得另找票 83 的实跑记录才能翻。
- 另注：框后那句「**本票永远不勾这两框**，它们是**被裁决取消**、不是**没做完**」指的是 AC#3／AC#4；
  AC#5 这框**没人说过不勾**，它只是被具名接走了 ⇒ 若走乙类处置（搬出勾框），**别把它和 AC#3/AC#4 混成一句**。

## 5. 票 89 —— `89-0600-is-decorative-on-windows-acl-for-private-data-done.md`（1 枚）

> 本腿现量：un = 1，但**它不在 AC 区**——「## Acceptance criteria」的六枚框在 `:41 :44 :47 :51 :55 :57`，
> **全部 `- [x]`**（本腿 `grep -c -- '^- \[x\]'` 含进度日志条目共 14 枚，全票唯一一枚 `- [ ]` 就是 `:348`）；
> `:348` 的位置在「Progress log（append-only）」里面。

### 5.1 `:348` 待办：AC#4（链接）、AC#5（失败注入）、三条落盘路径接 `winsec`、AC#6 POSIX 侧点名 —— **乙（勾框形状用在进度日志的待办行上；四条 item 在同票后续条目各自有落点）**

- 原文逐字：`- [ ] 待办：AC#4（链接）、AC#5（失败注入）、把 memory/agent/secret 三条落盘路径接上 `winsec`、AC#6 的 POSIX 侧点名。`
- 为什么不是判据：它是**进度日志里的一行待办**，不是本票「## Acceptance criteria」那一节的东西（本腿现读 AC 区
  `:41`-`:57` 六枚全勾）。它列的四件事，**同票后来的已勾条目各自都写了落点**：
  - AC#4 链接 → `:363` `- [x] **AC#4 A51②：本机可构造，用的是 junction（`mklink /J`，普通权限即可）**`
  - AC#5 失败注入 → `:352` `- [x] **AC#5 失败方向只能收紧（红→绿有名字）**`（点名 `internal/winsec/private_fail_test.go`、
    四条腿、`assertNoBytesOnDisk`）
  - 三条落盘路径 → `:382` `- [x] **落盘路径接线（票面点名的三个家）**`，逐字点名
    `agent/spill.go`／`secret/store.go`／`memory/open.go`／`memory/artifacts.go` 各自的 `PrivateDirAll`/`PrivateFile*` 落点，
    并附生产端到端 icacls 证据（`TestAC3ProductionDataRootIsPrivateEndToEnd`、`TestAC3SpillArtifactLandsPrivate`）
  - AC#6 POSIX 侧 → `:392` `- [x] **AC#6 门禁（只跑本票碰的包）**` 里有 `GOOS=linux go vet` 同四包 rc=0；
    POSIX 语义另在 `:179`「⚠ POSIX 侧不适用：`sealDir` 在 `_other.go` 里没有 walk…」与 `:377`
    `TestPOSIXSymlinkAtArtifactPositionIsNotRecursed` 具名
- ⇒ **处置建议＝搬出勾框**（本腿不动手）：这枚之所以占分母，纯粹因为进度日志借用了 `- [ ]` 语法。
  ⚠ 但它**不是零内容**：那行待办里"每条落盘写路径各自过没过封"这一问，今天确实有一处**没接的同族**，
  而且已经**具名归口**到别张票——见下。
- 残余内容的归口去处（**本腿现读**，`.scratch/wisp/issues/132-log-files-are-never-sealed-sealdirdir-has-zero-production-callers.md`，
  **无** `-done` 后缀＝仍开放，un=5：`:31 :34 :38 :40 :43`）：
  - 票 89 `:385-386` 自己写了「⚠ **没接的同族**（本票落点之外，要编排者拍板）：`models/downloader.go:224/:389`（staging…）、
    `config/migrate.go:83`+`config/parse.go:213`（0o600 配置备份）、`observe/logging.go:244`（0o644 日志）、`ball/position.go:77`」
  - 票 132 的「## 编排者增量（09-29 17:4x，起手锚点 `a8f3c020`）」（`:58` 起）逐字
    「**本票 AC#1 那张"现状表"从此多一枚具名行，而且它不是日志**：`internal/memory/open.go:533` 的
    **迁移前 DB 备份副本**＝全仓（本尺射程内）唯一一处"写私有数据而不经封条"的裸创建。**不必新立 AC#6**——
    它 rides 在 AC#1 现成的分母」⇒ "逐条写路径过封"这张清单的**承载体＝票 132 AC#1**（`:31`，仍 `- [ ]`）。
  - ⚠ 本腿只核到"132 具名接了备份副本那一处＋AC#1 现状表"，**没核**票 89 列的另外四处
    （downloader / config 备份 / logging / ball position）是否逐条进了 132 的射程——132 正文 `:20` 的表里有
    `config.toml` 一行（逐字「`config.toml` | **已锁** | `internal/config/parse.go:215` 调 `winsec.SealFile(tmpName)`」）。
    **要交回这张"同族清单是否被 132 全覆盖"的核对，得由能跑尺的人做**，本腿停在这儿并注明射程。
- ⚠ 顺带一枚引用漂移（详见附节）：本行上方 `done-fix-1` 追加段（`:349`）引
  「其余三样本腿在票内现读到已勾：AC#4 `:262`、AC#5 `:343`、AC#6 `:57` 与 `:383`」——本腿现读这四个数：
  AC#4 的 junction 条目在 **`:363`**（不是 262；`:265` 另有一枚同名早期条目）、AC#5 在 **`:352`**（不是 343）、
  落盘接线在 **`:382`**（不是 383）。只有 `:57` 对得上。

## 6. 票 92 —— `92-panel-composer-mode-attachments-workspace-done.md`（4 枚）

> 本腿现量：un = 4（`:52` `:68` `:73` `:79`）；AC 区共 7 枚（AC#1–AC#7），已勾 3 枚在 `:56` `:62` `:65`
> （两条 `grep` 现跑：un=4／chk=3，非心算）。
> 票头 Status 逐字（`:3`）「**accepted-done（附条件，条件已归位）**」＋`:9`「**条件我全部落到票 114**（见其 AC#8–AC#11），
> 不在本票静默结案」⇒ 这张票的"残余"在票面上就已经**具名归口**过一次，本腿的任务是核实那个归口今天还在不在。

### 6.1 `:52` AC#1 面板只读显示档位 + 不存在改档位路径（ban #6 正向钉子）—— **丙（账在票 114；钉子本体已在盘上，欠的是覆盖面）**

- 原文行首逐字：`- [ ] **AC#1** 面板**只读地**显示当前档位（三档之一），并且**面板侧不存在任何能改档位的路径**：`
- 归口去处**本腿现读**：`114-composer-request-has-no-production-caller-and-the-real-gate-must-be-native.md`
  （**无** `-done` 后缀＝仍开放）`:50` 是一枚 `- [ ]`，逐字
  「**AC#5** 不越界：面板侧**只许显示 + 发起请求**（R20 的明写）。工作区/档位两个输入口不得变成授权口；
  `frontend/` 里出现 `approval.decide` 即 ban #6 红；零 emoji（ban #8）含测试文件与注释」；
  覆盖面那半另对准 114 `:47` AC#4（「`R-92-1` 的覆盖面：要么把门二的扩展名集扩到与 d22scan 的文本类一致」）与
  `:67` AC#9（「结构钉的扫描根比 `ban #6` 与"渲染器实际加载的文件集合"都窄」＝`R-92b-2`）。
- 裁决表读数（本腿现读）：`docs/evidence/s1/92b-adversarial-acceptance.md:134` 逐字
  「| AC#1 档位只读 | **不通过**（按派单口径：门的覆盖面残留） | 〔独立复现〕 | 反半边我钉住了（2 处合法调用点照旧绿）；正半边 F5/F6 绿 |」
- ⚠ 本腿补一条 done-fix-1 没说的盘上事实：**AC#1 要的那枚正向钉子今天已经有实现体**——
  `internal/panel/composer_test.go:93` `TestPlantedComposerModeWriteGoesRed`，其注释逐字
  「is AC#1's positive nail. The fixture is a composer component written the WRONG way on purpose,
  and the assertion is that the instruments catch it - not that the real tree happens to be quiet」，
  两枚 planted 形状在 `:97`（`approval.decide` 改自己的 mode）与 `:109`（不含禁词的 `panel.mode.set` 写入），
  禁词正则 `composerModeWriteRe` 在 `:88-90`。⇒ **这格欠的不是钉子本身，是"扫描根/扩展名覆盖面"**，
  而那件事 114 已具名接走 ⇒ 判**丙**，且**丙得很干净**（功能面已落、残面有主）。
- ⛔ 不判丁：`92b:134` 是验收方（非实现者）的**不通过**裁决，且 114 的对应格仍 `- [ ]`。

### 6.2 `:68` AC#5 变异三向 —— **甲（残余真活：第 (iii) 向从未被非实现者复跑）**

- 原文行首逐字：`- [ ] **AC#5** 变异三向：(i) 把"面板只能显示"改成"面板能写 mode" ⇒ 必须有用例红；`
- 缺什么，说得出口：**三向里的 (iii)**（把 C26 reparse 那次的拒绝改成放行 ⇒ **既有**安全用例必须红）
  与 **(i) 的绿半边**。凭据（本腿现读）：`docs/evidence/s1/92b-adversarial-acceptance.md:137` 逐字
  「| AC#5 变异三向 | **(i) 本轮我重做为"红+绿各半"**（F3a 红 / F5 绿）；(ii) MUT-H 红；**(iii) 未重做** |
  (i)(ii) 〔独立复现〕／(iii) 〔仅自述，不背书〕 |」
- ⇒ (iii) 今天**只有实现方自述**，验收方明写"不背书"；(i) 那向连验收方自己也只做到"红＋绿各半"
  （F5 仍绿＝`R-92-1` 未清，见票头 `:7`「`R-92-1` 未清（四形里 F5/F6 两形全绿）」）。
- 归口核查：本腿把 114 的 AC 区 `:34`-`:75` 逐枚读下来，**没有任何一枚**是"重做 AC#5 的 (iii) 变异"
  ⇒ **不是丙**；也不是乙（它是一条要读数的判据，不是承诺句）。
- ⚠ 本腿补不了这把尺：派单禁一切编译/测试，`/tmp` 纯净快照里的变异跑属同一禁面 ⇒
  **这一格的丁类升级只能由一枚能跑的腿做**，本腿交"甲＋那把尺的形状"（票面 `:68-72` 已写死锚点与还原证明要求）。

### 6.3 `:73` AC#6 台账与门禁 —— **甲（残余真活：三发读数只对两棵旧树有效，当前 HEAD 无读数）**

- 原文行首逐字：`- [ ] **AC#6** 台账与门禁（只跑自己碰的范围）：`sh scripts/d22scan.sh` 纯净树 rc=0 且贴出**逐作用域文件数**`
- 缺什么：**`sh scripts/d22scan.sh` rc=0、`gofmt/gofumpt -l` 空、`go test -count=2 ./internal/panel/` 四数**
  这三发在**当前 HEAD** 上的一次非实现者读数。凭据（本腿现读）：
  - `docs/evidence/s1/92-adversarial-acceptance.md:147` 第一轮逐字
    「**结论：FAIL —— 附我本机实测数字。**（AC#6 原文要求 `gofumpt -l` 空，它不空；且 POSIX 读数不可复现。）」
  - `docs/evidence/s1/92b-adversarial-acceptance.md:138` 第二轮「| AC#6 台账与门禁 | **通过（数字全复算）** | 〔独立复现〕 |
    `gofmt`/`gofumpt` 空、`go vet ./internal/panel/ ./internal/tools/ ./internal/memory/` rc=0、`go test -count=2 -v …`」
  - ⇒ "通过"是**对 `91b5fc4`／`a8f9459` 那两棵树**的读数（票头 `:5-6` 的复算式子写的就是这两个号），
    结案票的门禁读数天然会过期；本格还另留一句未销的 `R-92b-4`（`92b-…md:290` 逐字
    「票面 `20:5x` 编排者注要求"AC#6/AC#7 读数必须包含整步 `portable-tests.sh` 并说清是谁的账"⇒ 92b 交件 grep `portable` **0 次命中**，未答」）
- ⚠ 给裁决者的一条结构性话：**门禁类 AC 在 `-done` 票里永远会"过期"**，把它记成甲就等于"每次 HEAD 变化都要重跑一次已结案票的门禁"；
  若编排者不想背这个，处置方式应该是**改判据形状**（例如"结案时点有效＋过期即由后续票的门禁继承"）并搬出勾框，
  而不是留着它当下一次复算的靶子。本腿只指出形状，**不动手、也不替它选**。
- 归口核查：票 83 的 `:57` 那枚门禁勾了，但那是**票 83 自己的包面**（`internal/config`），不是 92 的 `./internal/panel/` ⇒ 不是丙。

### 6.4 `:79` AC#7 负判据（不做 git 切换）—— **判不了＝凭据有一半在禁读地界（`frontend/**`）**

- 原文行首逐字：`- [ ] **AC#7** **负判据**：把"不做 git 切换"变成可检查的东西——在 `frontend/` 与 `internal/panel/` 里`
- 本格要求**两面都** `grep -rn`。`frontend/**` 是派单对本腿写死的零读零引零转述地界 ⇒ **本腿无法判整格**，
  按纪律停在这儿并具名登记（见附节）。
- 能验的那一半本腿验了，且**形状与读数都在**：
  - 禁词正则与用例在盘：`internal/panel/composer_test.go:269-270`（逐字
    `\bgit\s+checkout\b|\bgit\s+switch\b|\bswitchBranch\b|\bcheckoutBranch\b|\bchangeRepo(?:sitory)?\b|` ＋
    `\brepoPicker\b|\bbranchSelect(or)?\b|\bworktree\b|\bgit\.branch\b|\bgit\.repo\b|\bvcs\.switch\b`）
  - 本腿现跑负向尺：`grep -rniE 'git (checkout|switch)|switchBranch|checkoutBranch|changeRepo|repoPicker|branchSelector' internal/panel --include='*.go' | grep -v '_test.go'` ＝ **0 命中**
  - 裁决表：`docs/evidence/s1/92b-adversarial-acceptance.md:139` 逐字
    「| AC#7 不做 git 切换 | **通过** | 〔独立复现〕 | 包内用例 PASS + 我把 `fixtures/`、`dist/` 也 grep 了一遍 0 命中」
    ⇒ 出处＝非实现者（92b 验收腿），**但它那次 grep 跨 `fixtures/`／`dist/`，本腿不能替它现验**
- ⇒ 本腿**不写丁**（done-fix-1 那段把它记成"丁类·本腿碰不到全格"是自相矛盾的：既碰不到就不能判丁）。
  需要一枚有 `frontend/**` 读权的腿补半格，才谈得上翻勾。

## 7. 票 97 —— `97-dead-strict-param-and-the-comment-that-invents-a-caller-done.md`（2 枚）

> 本腿现量：un = 2（`:53` `:55`）；AC 区共 5 枚，已勾 3 枚在 `:45` `:48` `:51`（un/chk 两条尺现跑：un=2／chk=3）。

### 7.1 `:53` AC#4 注释与代码一致 —— **丙（账在票 230 AC#2；⚠ 本腿重读确认那句注释今天仍然撒谎，所以是"丙而活未死"）**

- 原文行首逐字：`- [ ] **AC#4** 注释与代码一致：贴出你改前后的注释原文，并说明**新注释的每句话在代码里能找到对应物**。`
- 缺什么，说得出口：`internal/agent/approval/queue.go:247` 的注释逐字（本腿现读）
  「// stays an unknown correlation id. This function is the only reader of q.alias」——而 `q.alias` 在同文件里
  有 **5 处命中**：`:197`、`:200`、`:210`、`:216`、`:257`（本腿现跑 `grep -n 'q\.alias' internal/agent/approval/queue.go`）。
  验收方当年就把它记成外项：`docs/evidence/s1/97-adversarial-acceptance.md:34` 逐字
  「**R-97-1** `queue.go:247` 那句 "This function is the only reader of `q.alias`" 字面为假（`:197`/`:210` 也读）」，
  并在 `:16` 的 AC#4 行记「**通过（附 2 条找不到对应物的句子）**」。
- 归口去处**本腿现读**：`230-four-cells-left-unfinished-inside-closed-tickets.md:21` 是一枚 `- [ ]`，逐字
  「**AC#2 票 97 那句注释改成实话**：`queue.go:247` 的"唯一读者"要么改成"另有 `:197`／`:210` 两处读者，
  但 `allow` 侧仍不可达，理由是……"，要么把那两处读者的存在具名解释掉」
  ⇒ 活已被开放票 230 具名接走 ⇒ 判**丙**；⚠ 但**这枚丙的功能面为零**：230 AC#2 要的就是改那一行注释，
  它一天不勾，97 这格一天就没有可翻的凭据。
- ⛔ 不判丁、也不判乙：判据本体（"注释不许撒谎"）今天未兑现；它是 AC 区的正式一格，不是措辞残留。

### 7.2 `:55` AC#5 门禁（按包）—— **甲（残余真活＝当前 HEAD 上的五个数；同一"门禁过期"形状见 6.3）**

- 原文行首逐字：`- [ ] **AC#5** 门禁（按包）：`gofmt -l`/`gofumpt -l` 空、`go vet ./internal/agent/approval/` rc=0、`
- 缺什么：**`gofmt -l`／`gofumpt -l` 空、`go vet` rc=0、`go test -count=2` 四数逐条点名（2 条 SKIP 必须点名）、
  收尾 `sh scripts/d22scan.sh`** 这五发在**当前 HEAD** 上的一次非实现者读数。凭据（本腿现读）：
  `docs/evidence/s1/97-adversarial-acceptance.md:17` 的 AC#5 行逐字开头
  「**通过（PASS 相加≠RUN 的写法要更正，数字自洽）** | 〔独立复现〕 | 全部在 `f140079` 纯净快照…」
  ⇒ 验收方的数**绑在 `f140079` 那棵快照树上**，对本锚点 `fac60ad4` 的树无效；本腿被派单禁跑任何门 ⇒ **补不了这把尺**。
- 同一条里还留了一笔**没销的格式账**：`:36` 逐字「**R-97-3** 票内计数式"58+38+2 = 98 对得上"跨层相加
  （顶层 58 + 子测试 38 + 顶层 SKIP 2），正确式子是"顶层 60（58 PASS + 2 SKIP」」⇒ 记甲时把这笔一起带上。
- 归口核查：本腿把 230 的 5 枚框（`:19 :21 :22 :23 :24`）逐条读下来，**没有一枚**接"97 的门禁重跑" ⇒ 不是丙。
- ⚠ 结构性话与 6.3 同一条：门禁类 AC 在结案票里天生会过期。**票 92 AC#6 与票 97 AC#5 是同一个形状**
  （读数绑旧快照／HEAD 无新读数），建议编排者把它们当**一枚类别**处置（统一改成"结案时点有效"或统一搬出勾框），
  别一张票一个判法。

## 8. 票 104 —— `104-sealfile-silently-drops-inherited-grants-done.md`（2 枚）

> 本腿现量：un = 2（`:52` `:57`）；AC 区共 5 枚（AC#1–AC#5），已勾 3 枚在 `:45` `:49` `:63`（un/chk 两条尺现跑：un=2／chk=3）。

### 8.1 `:52` AC#3 双向变异（三腿＋锚点＋`go build` rc=0 先量）—— **甲（残余真活＝当前 HEAD 上的一次非实现者变异复跑；验收读数只绑旧快照）**

- 原文行首逐字：`- [ ] **AC#3** 双向变异：① 把检测退回"只看显式 ACE" ⇒ AC#1 红；`
- 当年这格被裁过而且裁得很硬（**非实现者**）：本腿现读 `docs/evidence/s1/104-adversarial-acceptance.md:115` 逐字
  「**结论：通过。〔独立复现〕**（票面要求的三腿变异 + 我自加两发全部我自己在 `/tmp` 仓外快照重抽、重打、重量；
  另有一发 HEAD 漂移核对，非变异）」
- 缺什么，说得出口：**"绑旧快照"这四个字**。同表 `:162` 逐字「每个数字都是我在 `/tmp/ac104-fk`（**`4d43447` 纯净树**）
  本机重跑的，未照抄自述」⇒ 三腿变异的红/绿是 `4d43447` 那棵树的读数，对本锚点 `fac60ad4` 的树**无效**。
  本腿禁跑一切编译/门，**补不了这把尺** ⇒ 判甲，不判丁（丁要求盘上凭据对**今天**成立）。
- ⚠ 一枚加重这条判断的票面事实（本腿现读 104 `:68-71` 的编排者注，逐字）：
  「**不要替票 115 背账。** 现在 CI 的 `test-windows` **step 4**（run `35599458439` / job `106331840177`）上有三枚与本票同族的红：
  `TestSealReportsThePrincipalsItCleared`…**上一枚 run 它还是 PASS**、`TestAC1SealFileReportsTheInheritedGrantItCleared`
  （`inherited_narrow_notice_104_windows_test.go:134`…逐字 `reported 0 notice(s), want exactly 1`…）」
  ⇒ **AC#1 的承载用例本身在同族里出现过 CI 红**，说明"变异在旧快照红过"不等于"今天在树里还有牙"。
  那批红的账由票 115 接手（115 也在本 28 枚清单里，见第 12 节）。
- 归口核查：票 118 的 AC#1–AC#9（本腿现读 `:29 :30 :32 :34 :36 :43 :53 :59 :67`）**没有一枚**是"重跑 104 的三腿变异"
  （118 收的是 `R-104-1`/`R-104-6` 那组测试加固）⇒ **不是丙**。

### 8.2 `:57` AC#4 回归（四包 `-count=2 -v` 四数＋逐条点名 SKIP＋门禁）—— **甲（同一"读数过期"形状，第 3 次出现）**

- 原文行首逐字：`- [ ] **AC#4** 回归：`go test -count=2 ./internal/winsec/ ./internal/memory/ ./internal/secret/ ./internal/agent/` rc=0，`
- 非实现者读数（本腿现读）：`docs/evidence/s1/104-adversarial-acceptance.md:162` 逐字
  「**结论：通过。〔独立复现〕**（每个数字都是我在 `/tmp/ac104-fk`（`4d43447` 纯净树）本机重跑的，未照抄自述）」
  ⇒ 四数＋`gofmt`/`gofumpt`＋`GOOS=linux go vet`＋`sh scripts/d22scan.sh` 全是**纯读数**，本腿一律不许跑，一发也补不上。
- ⇒ 判**甲**：缺的东西很具体＝这五发在**当前 HEAD** 上的一次非实现者读数。
- ⚠ **结构性话（第三次出现，请当一枚类别看）**：6.3（票 92 AC#6）、7.2（票 97 AC#5）、本格＋9.2（票 105 AC#5）、
  11.2（票 113 AC#5）、12.3（票 115 AC#6）、10.2（票 110 AC#5）全是**"门禁/回归读数绑结案时点"**这一族。
  把每一枚都当"加回分母的真活"，等于给每次 HEAD 变化欠一发重跑；本腿建议编排者**一次性定一条规则**
  （要么"读数绑快照号，过期由后续票的门禁继承"，要么统一搬出勾框），不要逐枚判。
  ⚠ 上面那串行号里凡本腿尚未走到的小节，以本表后续节的现读为准（本行是写作顺序，不是证据顺序）。

## 9. 票 105 —— `105-c26-rewrite-account-has-no-production-reader-done.md`（2 枚）

> 本腿现量：un = 2（`:47` `:54`）；AC 区共 5 枚（AC#1 `:38`／AC#2 `:42`／AC#3 `:47`／AC#4 `:51`／AC#5 `:54`），
> 已勾 3 枚（un/chk 两条尺现跑：un=2／chk=3）。

### 9.1 `:47` AC#3 祖先重解析那条腿的行为用例 + 变异 —— **丁（今天其实已满足，凭据在票 116 名下且验收方具名放行 105 结案）**

- 原文行首逐字：`- [ ] **AC#3** 补祖先重解析那条腿的**行为用例**：构造"祖先带 reparse/junction"的输入 ⇒ 断**可观察结果**`
- 本格当年被谁卡住（本腿现读）：`docs/evidence/s1/105-adversarial-acceptance.md:116` 判 **FAIL-退回**（`R-105-2`），
  总判 `:203` 逐字「**FAIL-退回（只退 AC#3 这一格，其余四格通过）**」⇒ 缺的东西当年很具体：**行为用例，不是符号用例**。
- 盘上凭据（**本腿今天在树里现读**）：`internal/risk/syncdirs_ancestor_actable_leg_116_test.go`
  - `:153` `func TestSyncAncestorActableLegFailsClosedOnMovedAncestor116`（正是"祖先腿"那一枚行为用例）
  - `:239` / `:270` `TestSyncFirstActableLegStillRefusesEnvRewriteAtVerdict116`／`...HomeRewriteAtVerdict116`（反半边）
  - `:307` `TestSyncAncestorLegStillMatchesPlainSameTreeWrite116`（同树成功那半）
- 凭据出处＝**非实现者**：`docs/evidence/s1/116-adversarial-acceptance.md:1` 署名「票 116 对抗验收 —— `acceptor-ticket116`」，
  `:108` 总判逐字「**通过（附一条非阻塞措辞修正 R-116-2）**」，理由段逐字
  「票 116 存在的唯一理由——"删掉祖先那一次 `Actable()` 腿时，有没有一枚**行为**用例会红、而票 105 那两条仍绿"——
  被我在 Windows 与 docker POSIX 两侧**独立复现成立**（M1/M2 各自把 `TestSyncAncestorActableLegFailsClosedOnMovedAncestor116`
  单枚钉红、红形是 `Sync` 翻转/`Why` 丢账，非编译期符号…）」，紧接着一句逐字
  「**票 105 能否因此结案：能。** `acceptor-ticket105b` 卡 AC#3 的那格…已被本票的行为用例补齐并独立复现」
- ⚠ 本腿给丁类**如实划出射程边界**（这是 28 枚里唯一一枚丁，别把它读得比证据硬）：
  ① 验收方的红/绿读数绑 `80e248c`（表 `:3` 逐字「**被验 sha**：`80e248c`（用例）＋ `63fdbd2`（票面 append）」），
  本腿**不许跑测试**，因此只能证"用例文件与四枚函数名今天在树里"，**不能**证它今天仍绿；
  ② 交付记在票 116 名下，不在 105 名下 ⇒ 翻 105 这框的人在票面上要留指针。
  ⇒ 判**丁**，但**翻不翻由编排者定**；本腿按纪律**不动任何勾**。

### 9.2 `:54` AC#5 门禁（按包 scope，五发＋各 scope 文件数）—— **甲（"读数绑结案时点"那一族的第 4 次出现，见 8.2 的结构性话）**

- 原文行首逐字：`- [ ] **AC#5** 门禁（按包 scope）：`go test -count=2 -v` 各包 rc=0 并逐条点名 SKIP/FAIL（**报 `=== RUN` 行数 == 不同测试名 × 2**；`
- 当年读数（非实现者，本腿现读）：`docs/evidence/s1/105-adversarial-acceptance.md:151` 那一节声明逐字
  「**口径**：全部跑在 `f6818f2` 的**纯净快照** `/tmp/ac105b-gates`…交件的读数是在**当前工作树**量的
  （它自己标了"树＝当前工作树"）⇒ 两处不同步的数我按快照重算，并点名差异」
- 缺什么很具体：**`go test -count=2 -v` 各包四数 ＋ `gofmt`/`gofumpt -l` 空 ＋ `go vet` 与 `GOOS=linux go vet` 按包 rc=0
  ＋ `sh scripts/d22scan.sh` 纯净快照 rc=0 且各 scope 文件数不降**，这五发在**当前 HEAD `fac60ad4`** 上的一次非实现者读数。
  本腿能静态读到的只有"仪器还在"：`scripts/d22scan.sh` 在场、`tools/d22scan/main.go` 在场（本腿只 `ls`/`grep`，未跑）。
- 归口核查：票 116 的 AC#5（`:52`，已勾）是**它自己包面**（`internal/risk/`）的门禁，不是 105 那四包 ⇒ 不是丙。
- ⚠ 另记一枚本腿**推翻旧话**的读数（不影响判性，但影响后来人引用）：票 105 的标题与 AC#1 钉的是
  "`RewrittenRoots()` 这本账**零生产读取者**"，本腿今天在树上重跑那把尺，**读数已经变了**——
  `internal/tools/bridge.go:1053` 逐字 `return b.paths.Roots(), b.paths.RewrittenRoots(), b.paths.UnusableRoots()`
  （其宿主 `:1049 func (b *Bridge) pathAccount() (roots, rewritten, unusable []string)` 本腿一并现读）。**"零生产读取者"这句今天不为真**，而它不为真正是本票 AC#2 接上的那条线；
  引用这句话当现状的人请把票号/日期带上。那把尺具名：`grep -rn 'RewrittenRoots(' --include=*.go internal cmd | grep -v _test`
  ⇒ 三行命中：`bridge.go:999`（注释）／`bridge.go:1053`（**真读取点**）／`paths.go:194`（定义）。

## 10. 票 110 —— `110-no-ci-step-runs-internal-winsec-done.md`（2 枚）

> 本腿现量：un = 2（`:39` `:47`）；AC 区共 5 枚（AC#1 `:32`／AC#2 `:35`／AC#3 `:39`／AC#4 `:43`／AC#5 `:47`），
> 已勾 3 枚（un/chk 两条尺现跑：un=2／chk=3）。
> ⚠ 本票有一枚勾是 `done-check-1` 翻的（AC#4 `:43`），它与 `done-fix-1` 的判定**公开分歧**——见 10.3，那与分母无关但影响裁决者对这两枚的信任度。

### 10.1 `:39` AC#3 这道新步要自己会红（两发变异）—— **甲（残余真活＝三发变异在当前 HEAD 的纯净快照上的一次非实现者复跑）**

- 原文行首逐字：`- [ ] **AC#3** 这道新步要**自己会红**：在 `/tmp` 快照里把 winsec 某条安全断言人为弄坏（例如私有集改成按名字比）⇒ 新步必须 rc≠0；`
- 当年判过，而且是非实现者独立复现（本腿现读 `docs/evidence/s1/110-adversarial-acceptance.md:16`）逐字
  「| AC#3 | 新步自己会红 + 不是空仪器（同链 `grep -n` 证落地、先 `go build` rc=0） | **通过（三发全部独立复现）** | 〔独立复现〕 |
  全部在 `/tmp/wisp-ac110-snap`；仓库树 winsec 全程未动」⇒ 读数绑那枚 `/tmp` 快照，**不绑今天的树**。
- 缺什么很具体：①弄坏安全断言 ⇒ 新步 rc≠0；②把包清单改成不含 winsec ⇒ "扫描空=红"的守卫必须红；③每发先 `go build` rc=0。
  本腿**禁编译禁跑**，一律补不出；能静态读到的只有仪器本身在位：
  `scripts/winsec-tests.sh` 在场，`.github/workflows/ci.yml:398` 逐字
  `- name: "Windows ACL sealing gate (internal/winsec's own tests, ticket 110)"`。
- 归口核查：111／112／140 三张开放票（本腿现量 **111 un=10／112 un=5／140 un=4**，都不带 `-done`）
  接的是"这一步把后面的步骤吃掉"与"CI 覆盖面"那两族账，**没有一枚**是"重跑 AC#3 的三发变异" ⇒ 不是丙。

### 10.2 `:47` AC#5 门禁（`bash -n`／`d22scan.sh` 各 scope 不降／新步不得排在会失败的步之后）—— **甲（"读数绑结案时点"一族第 5 次；且第三句的后果今天仍挂着）**

- 原文行首逐字：`- [ ] **AC#5** 门禁：`bash -n` 改动脚本 rc=0；`sh scripts/d22scan.sh` 纯净快照 rc=0 且台账各 scope 不降；`
- 当年读数（非实现者，本腿现读同表 `:18`）逐字开头
  「| AC#5 | `bash -n` rc=0；`sh scripts/d22scan.sh` 纯净快照 rc=0 且各 scope 不降；新步不得排在会失败的步之后 | **通过（附一条后果登记）** | 〔独立复现〕…」
- 缺什么：**这两发（`bash -n`／`d22scan.sh`）在当前 HEAD `fac60ad4` 上的一次非实现者读数**——本腿连 `bash -n` 都不许跑（属执行）。
- ⚠ **本格第三句不是纸面条件**：票 110 那枚改名 commit `f161c690` 正文自陈这一步上线后"把后面的步骤全吃掉"，
  后果今天仍挂在**三张开放票**上（10.1 现量的 111／112／140）⇒ 即便有人重跑读数翻这格，也**不是零后果结案**。
  本腿现跑 `git show -s --format=%s f161c690`，标题逐字
  「docs(97-done,110-done,111+AC#6-#8): 两张结案；票 110 抓到一条反噬——加一道门把后面的步骤全吃掉了」。
- 归口核查：同 10.1，没有别张票具名接走"重跑 110 的门禁" ⇒ 判**甲**。

### 10.3 附带一条**两腿公开分歧**的记录（本格已勾，不在 28 枚里，但裁决者需要知道）

- 已勾的 `:43` AC#4 上叠着两段互不相同的结论：
  - `done-fix-1` 记「**戊类**·附条件的兑现句在码里、但无表侧复算，⛔ 不翻勾」，并引表 `110-…md:17`
    「**通过（附条件：步级证据至今不存在）**」；
  - `done-check-1` 却把它**勾了**，理由是「本格判据是一个『或』句…第二支不需要读数」，并留下一句
    「⚠ 与 `done-fix-1` 的分歧…分歧交编排者裁，本勾若被撤请连带撤这行的结论」。
- 本腿现读：表 `docs/evidence/s1/110-adversarial-acceptance.md:17` 确为「**通过（附条件：步级证据至今不存在）**」；
  那枚"附条件"的账今天落在**开放票 111 的 `:75`**（本腿现读逐字「- [ ] **AC#7** `R-110-4`：票 110 承诺过…」，111 不带 `-done`）。
- ⇒ 本腿**不裁这枚勾该不该撤**（不在分母内），只登记：同一格上两腿结论相反，且"附条件"部分已被具名记在 111 名下。

## 11. 票 113 —— `113-posix-platformverifypplacement-has-no-link-leg-done.md`（2 枚）

> 本腿现量：un = 2（`:54` `:56`）；AC 区共 6 枚（AC#1 `:42`／AC#2 `:47`／AC#3 `:51`／AC#4 `:54`／AC#5 `:56`／AC#6 `:61`），
> 已勾 4 枚（un/chk 两条尺现跑：un=2／chk=4）。

### 11.1 `:54` AC#4 变异（关掉新腿 ⇒ AC#1 红；"祖先链只查一层"的半修 ⇒ 也红）—— **甲（残余真活＝四发变异在当前 HEAD 的一次非实现者复跑；读数绑容器快照）**

- 原文行首逐字：`- [ ] **AC#4** 变异：把新腿关掉 ⇒ AC#1 那条必须红；再把"祖先链只查一层"这种**半修**形状试一发 ⇒ 也要红（证明它咬的是全集不是某一行）。`
- 当年判过且是非实现者自证的（本腿现读 `docs/evidence/s1/113-adversarial-acceptance.md:16`）逐字
  「| AC#4 变异四发 | MUT-A 9 红 / MUT-B 恰 1 红 / MUT-C 9 红 / MUT-D2 2 红 | **四发全部我自己下刀、自己复量**…|〔独立复现〕| **绿**」
  ⇒ 但这四发红/绿是在**验收当时的容器快照**里量的，对本锚点 `fac60ad4` 的树**没有读数**。
- 缺什么很具体：把 `winsec_other.go` 的祖先链接腿关掉 / 改成只查一层这两发变异，在**当前 HEAD 的纯净快照**里
  ⇒ 看 AC#1 那条与半修检测是否红。本腿禁跑（连 `go build` 都不许），**盘上也没有"专属物证"可代**——
  这一格要的就是一发读数（`done-fix-1` 在票面 `:55` 也是这么记的，本腿复认）。
- 归口核查：票 116 接的是 `internal/risk/syncdirs.go` 的**祖先 `Actable()` 腿**（另一族，见 9.1），
  不是 `internal/winsec/winsec_other.go` 的 POSIX 祖先链接腿 ⇒ **不是丙**，两族别混。

### 11.2 `:56` AC#5 门禁（容器内三包 `-count=2 -v` 四数＋按包 vet＋`gofmt`＋`d22scan.sh`）—— **甲（"读数绑结案时点"一族第 6 次；⚠ 同一条里还有一笔没销的配对错账）**

- 原文行首逐字：`- [ ] **AC#5** 门禁：容器内 `-count=2 -v ./internal/winsec/ ./internal/memory/ ./internal/risk/` rc=0 且四数逐条点名（报 SKIP 要说是不是 `-v`；`
- 当年读数（非实现者，本腿现读同表 `:17`）逐字开头
  「**绿（四数与三门），但 ban#8 那格的 sha↔数字配对错了一格**」——五包 `-count=2 -v` 复算 `540/532/0/8 rc=0`，
  同一行另记 `sh scripts/d22scan.sh` rc=0 但配对有误 ⇒ **这格连"当年"都不是干净的通过**，
  要翻得先补两样：①当前 HEAD 上的这批读数；②把那笔 sha↔数字配对错账销掉。
- ⇒ 判**甲**。与 8.2／9.2／10.2 同族，处置应一并考虑（见 8.2 的结构性话）。
- 归口核查：票 113 自己的 AC#6（`:61`，已勾，编排者追加）不是这五发的接管方；
  本腿在 230／111／112／140 四张票里 `grep -n` 也没有找到"重跑 113 门禁"的具名句 ⇒ 不是丙。

## 12. 票 115 —— `115-…-while-the-cases-compare-caller-spelling-done.md`（3 枚）

> 本腿现量：un = 3（`:44` `:58` `:64`）；AC 区共 7 枚，已勾 3 枚（`:47` `:52` `:68`），
> 另有 **1 枚已被搬出勾框**：`- ⛔ **AC#5**` 在 `:61`（`grep -c '^- ⛔'` = 1）⇒
> **乙类处置在前一轮已经落过一次手**，本轮表里凡判乙的格都可以照那个形状办（本腿自己不动手）。

### 12.1 `:44` AC#1 独立复算四条红 + 对质上一枚 run 的 PASS —— **丙（账在票 230 AC#1；⚠ 115 名下今天仍然一张表都没有）**

- 原文行首逐字：`- [ ] **AC#1** 独立复算上面四条红（**逐字**贴你跑出来的红点），并复算"上一枚 run 里 `TestSealReportsThePrincipalsItCleared` 是 PASS"这句`
- **本腿自己复算的那道硬事实**：`ls docs/evidence/s1 | grep -E '115'` ⇒ **零命中**
  ⇒ 票 115 结案却**没有任何一张裁决表**，而 AC 区开头那行（票面 `:42`）自己写着
  「## AC（1:1，裁决表 `docs/evidence/s1/115-*.md` 由验收方出）」⇒ 这格缺的不是读数，是**连出读数的地方都没建过**。
- 归口去处**本腿现读**：`230-four-cells-left-unfinished-inside-closed-tickets.md:19` 是一枚 `- [ ]`，逐字
  「**AC#1 票 115 要一张真裁决表**：那 7 格**必须由非实现者出一张 `docs/evidence/s1/115-*.md`**（七项检查按 `SPEC-10 §8`），
  表出来之后才谈补勾。⛔ 不许"因为文件名带 `-done` 就默认它通过了"——同形状本仓**撤销过一次**（票 62 的…」
  ⇒ 230 本身**仍开放**（un=5）⇒ 判**丙**：分母该记在 230 名下，而不是让 115 这张结案票继续背着。
- ⚠ 给裁决者的一条提醒：230 AC#1 收的是**"出一张表"**这件治理动作；AC#1 这一格自己还要**一次实跑**
  （复算四条红 + `gh api` 取 run `35595651898` 同一步日志对质）。表出来 ≠ 这格自动满足，
  但**这格能不能勾完全取决于那张表出不出来** ⇒ 两枚必须同批走。

### 12.2 `:58` AC#4 变异（退回原状 ⇒ AC#3 新用例红；只 `EqualFold` 的半修 ⇒ 也红）—— **甲（残余真活＝两发变异的独立复跑；承载用例今天在树里）**

- 原文行首逐字：`- [ ] **AC#4** 变异：把你选的修法退回原状 ⇒ AC#3 新用例必须红；再试一发"只比大小写不敏感"（`EqualFold`）这种**半修** ⇒ 也要红`
- 缺什么很具体：**两发变异的红/绿读数**（至今没有任何非实现者复跑过——115 名下连表都没有，见 12.1）。
- 盘上凭据（本腿今天现读，能证的只有"形状在位"）：
  `internal/winsec/notice_attribution_115_windows_test.go:191` `TestNoticeAttributionSurvivesAn83AliasOfItsOwnTree`
  ＋ `:240` `TestNoticeAttributionKeepsTwoTreesApart`（正是 AC#3 要求的正/反两半枚用例）
  ⇒ **有用例、无变异读数**，所以这格既不能判丁（丁要的是"今天已满足"），也不是丙（115/230 都没具名接走变异这两发）。
- 本腿补不了这把尺：禁编译禁跑，`/tmp` 纯净快照的变异属同一禁面。

### 12.3 `:64` AC#6 门禁（按包 `-count=2 -v` 四数／`gofmt`／`gofumpt.exe -l`／`go vet`／`d22scan.sh`）—— **甲（一族第 7 次；这格还额外欠一句"本机有这个二进制"的原文）**

- 原文行首逐字：`- [ ] **AC#6** 门禁：按包 `-count=2 -v` 四数逐条点名（报 SKIP 要说是不是 `-v` 量的；`-count=2` 不缓存）；`
- 判**甲**的理由与 8.2／9.2／10.2／11.2 同一族：**当前 HEAD 上的一次非实现者门禁读数**（五发）。
  而且 115 的这批数出在**票面自己的 Progress log**里（本腿现读 `:251`／`:360`／`:444` 三处同款
  `=== RUN=164  --- PASS=84  --- FAIL=0  --- SKIP=0`，属实现方自述那批），
  配合 12.1 的事实（**没有裁决表**）⇒ 这格连"当年有没有非实现者复算"都是空的，比其它几枚更硬。
- ⚠ 本格还自带一条**措辞级要求**本腿记一下：判据明写 `gofumpt` 那发「**本机有这个二进制**，写"没有/未跑"必须引命令原文 + 错误原文」
  ⇒ 后来人交这格时不能只贴 rc，得把命令原文贴出来。
- 归口核查：230 的五枚框（`:19 :21 :22 :23 :24`）里没有"重跑 115 门禁"这一枚；
  最接近的是 `:24` AC#5「整包终态（凡动 Go 侧注释都要跑）」，那是 **230 自己动码时的门禁**，不是 115 那五发的接管句 ⇒ 不是丙。

## 13. 票 118 —— `118-winsec-test-hardening-from-ticket104-acceptance-done.md`（1 枚）

> 本腿现量：un = 1（`:59`）；AC 区共 9 枚（AC#1 `:29`／AC#2 `:30`／AC#3 `:32`／AC#4 `:34`／AC#6 `:36`／
> AC#5 `:43`／AC#7 `:53`／AC#8 `:59`／AC#9 `:67`），已勾 8 枚（chk=8／un=1 两条尺现跑）。

### 13.1 `:59` AC#8（`R-115-3`）跨卷归属：先量，再判是不是生产洞 —— **丙（按票面规则本就该留 `[ ]`；账在票 126，且 126 已结案并真改了生产码）**

- 原文行首逐字：`- [ ] **AC#8（`R-115-3`，我新立的）跨卷归属：先量，再判是不是生产洞。**`
- 本格自己写的处置规则就是**交出去**（票面 `:63-64` 逐字）：
  「**如果它证明的是生产判据（不只是测试判据）会归错 ⇒ 停手回报，不要自己改 `winsec_windows.go`**，我会另开立票」
  ⇒ 所以这枚 `- [ ]` **不是"没做完"，是"按规则不该在这张票做完"**。
- 归口去处**本腿现读**（三样都对得上）：
  - 票 126 标题逐字（`126-winsec-tree-attribution-strips-the-volume-so-a-seal-on-one-drive-is-reported-against-another-done.md:1`）
    「`sameTree` 比对前**剥掉 volume 段**：实测 C: 上的密封通知被归到**从未被密封的 D: 树**
    （`R-115-3`／**票 118 AC#8 停手那一格**）」；`:5` 逐字「**Blocks:** 票 118 的 AC#8
    （那格按票面规则留在 `[ ]`，等本票）」
  - **126 已结案**：文件名带 `-done`，本腿现量 `grep -c -- '^- \[ \]'` = **0**
  - **126 有非实现者的表**：`docs/evidence/s1/126-adversarial-acceptance.md`，
    `:34` 逐字「**本文件写作状态：六格已全落盘（渐进写，每格一次），已收总判。**」，
    逐格裁语 `:40`-`:44`（AC#1「**通过**」／AC#2「**通过**（量成"守门错"）」／AC#3「**通过**」／
    AC#4「**通过附条件**」／AC#5「**通过**」），档位全部标〔独立复现〕
- 生产码今天的状态（本腿现读，允许地界）：`internal/winsec/resolve.go:455`
  `func sameTree(a, b string) bool {` 的第一条判断逐字 `if !sameVolume(a, b) || !sameAbsoluteness(a, b) { return false }`
  ⇒ 卷比较**已经装上**，即"是不是生产洞"的答案是**是**，且已被修。
- ⇒ 判**丙**。⚠ 本腿**不建议**翻 118 这枚勾：票面 `:63` 那句"停手回报＋我会另开立票"就是为这种格写的，
  翻它等于把 126 名下的交付记到 118 头上；**处置建议＝把这行的 `- [ ]` 改成指向 126 的指针**（与票 80 AC#3/AC#4 同法）。
- ⚠ 一枚 118 名下的次生事实（不改判性，但表上有人这么记过）：118 的验收表 `docs/evidence/s1/118-adversarial-acceptance.md`
  **对同一格出了两处读数**——`:34` 逐字「| AC#8 | 跨卷归属：先量再判（票面留 `[ ]`） | **通过（停手移交票 126）**；缝守腿不背书 |
  链＋判据＋跨卷读数〔独立复现〕／缝守腿〔仅自述，不背书〕」，收表 `:397` 另有一句
  「**通过（停手移交）**：量到了、判成生产洞、没动生产码、三条判据与责任链我全复算」。
  ⇒ "缝守那一腿"（`R-118-9`，见 126 票面 `:58`）**至今两任都没造出 fake resolver**，那半枚仍开放，
  账在 126 名下不在 118；本腿只登记，不替它认账。

## 附. 分母一致性 / 判不了的格 / 本腿推翻或补强的前提

### A. 28 枚分母：**对得上，而且经得起形状核查**

- 本腿现跑（`.scratch/wisp/issues/`）：`*-done.md` 共 **78 张**，`grep -c -- '^- \[ \]'` 求和 ＝ **28 枚／13 张**，
  逐票计数 `92(4) 07(4) 80(3) 115(3) 97(2) 113(2) 110(2) 105(2) 104(2) 89(1) 63(1) 118(1) 11(1)`
  ⇒ **与编排者 09-29 17:5x 的终数逐票一致**，锚点也是同一个 `fac60ad4`。
- ⚠ 本腿额外做了三道"这 28 会不会少数/多数"的形状核查（都是尺，不是印象）：
  - 缩进版未勾框 `grep -c -- '^[[:space:]]\{1,\}- \[ \]' *-done.md` ＝ **0** ⇒ 没有藏在列表缩进里的第二层勾框。
  - 大写勾 `grep -c -- '^- \[X\]' *-done.md` ＝ **0**；异形框（`^- \[` 且既非 `[ ]` 也非 `[x]` 也非时间戳）
    逐名 `sort -u` 后**零输出** ⇒ 尺的口径没有漏网形状。
  - 已被搬出勾框的行首 `- ⛔` 共 **4 枚**：`115-…md:61`、`84-pendingapproval-bounded-wait-refuted-done.md:46`、
    `89-…md:258`、`89-…md:311` ⇒ 这 4 枚**正确地不在 28 里**，而且 `89:311` 与 `89:348` 是**同一种"待办行"形状**，
    前者已被搬出、后者还留着 ⇒ 本腿 5.1 的"乙＋搬出"建议在同票里就有先例。

### B. 判不了／射程不足的格（具名，供编排者派下一腿）

| 格 | 判不了的原因 | 谁能补 |
|---|---|---|
| 票 92 `:79`（AC#7 负判据） | 整格要求对两面 `grep -rn`，另一面在派单写死的禁读地界 ⇒ **本腿只能判一半**，已记"判不了" | 有该地界读权的一枚腿 |
| 票 07 `:51`（Visual） | 判"甲"用的是允许地界的台账与表；**第二个分句的对照物**在禁读地界 ⇒ 本腿**无法**把它升成丁 | 同上 |
| 票 11 `:53` / 63 `:46` / 92 `:68` / 97 `:55` / 104 `:52` `:57` / 105 `:54` / 110 `:39` `:47` / 113 `:54` `:56` / 115 `:58` `:64` | 一律是**读数类**判据（`winlive` 真机跑 / 变异 / 门禁），本腿禁编译禁跑 ⇒ 只能判"甲/丙"，**判不了丁** | 一枚能跑的腿（且须非实现者） |

### C. 本腿推翻或补强的前提（4 条，都是"当场顶回"，不是事后诸葛）

1. **⛔ 一枚坏尺**：票 80 `:44` 的 `done-fix-1` 追加段用
   `grep -rn 'NewGate(' --include=*.go internal cmd | grep -v _test` ＝ 0 命中，
   推出"真跑在 `risk.Gate` 上那条用例今天仍无承载体"。本腿现读：**`risk.Gate` 不是构造函数**，
   真身是 `internal/risk/blacklist.go:76` 的 `func Gate(canonical string, bOverrides map[string]bool) PathDecision`；
   `NewGate(` 这个串**全仓（含测试）一处都没有** ⇒ 零命中证明不了任何事。
   ⚠ 且重尺后事实相反的一半存在：`internal/tools/mode.go:98` `d := risk.Gate(c, overrides)` **是生产调用点**
   （`overrides := b.confirmations()`，`:92`）。**判性不变（仍丙）**，因为填 map 的是运行期 L2 确认而非配置键
   （`internal/config/unwired.go:81` 自陈），但那句"无承载体"后来人不要再用。
2. **引用漂移实测**（`done-fix-1` 那一轮的票面引用，本腿逐条现读）：
   - 对 `docs/evidence/s1/97-adversarial-acceptance.md` 的五枚行号引用**全体 +1**：
     它引 AC#1 `:14`／AC#2 `:15`／AC#3 `:16`／AC#4 `:17`／AC#5 `:18`，真值依次 **`:13 :14 :15 :16 :17`**。
   - 对 `230-…md` 的三处引用**全体 −1**：80 票面引 `:21`（AC#3）与 `:22`（AC#4）、97 票面引 `:20`（AC#2），
     真值依次 **`:22`、`:23`、`:21`**。
   - 对 `89-…md` 内部的引用是**另一种错法**（不是 ±1）：它写「AC#4 `:262`、AC#5 `:343`、落盘接线 `:383`」，
     本腿现读三枚已勾条目在 **`:363`／`:352`／`:382`**。
   - ✅ 也有全对的：`07-adversarial-acceptance.md:10/:14/:27/:62/:76`、`92b-…md:134/:137/:138/:139`、
     `104-…md:29/:78/:115/:143/:162`、`110-…md:16/:17/:18`、`113-…md:16/:17`、`64-…md:44/:59/:60/:72/:83/:152/:234`
     本腿复算**全部命中**。⇒ 结论：**那一轮的行号不能整体信任，也不能整体作废，逐枚现读是唯一办法**（与本腿派单写法一致）。
3. **一张不存在的文件**：`docs/evidence/s1/07-adversarial-acceptance.md:14` 说待人项 H1「已登记
   **docs/reports/pending-human-review.md**」——该路径**在树里不存在**
   （`git ls-files | grep -c 'pending-human-review'` ＝ 0；`find docs -name 'pending-human-review*'` 空）。
   内容实际在 `docs/reports/pending-and-issues.md:5` 那一节「## 待人工审核（pending-human-review）」，
   H1 在 `:7`。⇒ 后来人按 `:14` 那句话去找文件会以为账丢了。
4. **接收票的计数已经漂了**（影响所有"账在某某票"的结论）：`done-fix-1` 记票 114「un=7／chk=2」，
   本腿现量 **un=9／chk=2**；票 12 现量 **un=3**（`:48 :148 :149`）、票 64 现量 **un=1／chk=8**、
   票 111 **un=10**、票 112 **un=5**、票 140 **un=4**、票 132 与 230 **各 un=5**。
   ⇒ **判性不受影响**（本腿全部现读），但任何"挪分母"的动作前，挪到的那张票自己还得再量一次。
   ⚠ 顺带一枚：票 12 `:148` 还挂着与 07 `:51` **同族**的那枚"人工视觉签收"未勾框 ⇒ 07:51 判甲时，
   它的姊妹格在开放票 12 名下也有账，别只看到 07 一处。

### D. 四分类的形状小结（哪些真该回分母）

- **甲 13** 里只有 **1 枚**是要人的眼睛（07 `:51`，H1 实况签收）；其余 **12 枚全是"要一发读数"**：
  变异类 5（92 `:68`、104 `:52`、110 `:39`、113 `:54`、115 `:58`）、
  门禁/回归类 7（92 `:73`、97 `:55`、104 `:57`、105 `:54`、110 `:47`、113 `:56`、115 `:64`）。
  ⇒ 这 12 枚**互相之间可以合并成一件事**：一枚非实现者腿在当前 HEAD 做一次"变异＋门禁"复跑，
  按票分格登记即可；其中 5 枚还各带一笔未销账（`R-92b-4`、`R-97-3`、113 的 sha↔数字配对、
  110 的反噬后果、115 无表）。
- **丙 12** 的去向要分开看：**10 枚归口到仍开放的票**（→ 64 三枚、→ 12 两枚、→ 230 四枚、→ 114 一枚），
  **2 枚归口到已结案的票**（80 `:54`→票 83 `:57`、118 `:59`→票 126）。
  ⇒ 建议：前者从结案票分母里摘掉、由那些开放票自己的分母承担（它们本来就在 un 计数里）；
  后者是"账已完成、只是名字还挂在旧票"，最省事的处置是**把旧票那两行改成指向句**（118 那格连指向都不用新造，
  票 126 `:5` 已经反向写着"Blocks: 票 118 的 AC#8"）。
- **乙 1**（89 `:348`）：同票 `:311` 已有 `- ⛔` 先例，照抄形状即可。
- **丁 1**（105 `:47`）：唯一一枚"今天其实已满足"，凭据在**已结案的票 116** 名下、
  且 116 的验收表具名放行「票 105 能否因此结案：**能**」；本腿按纪律不翻，交由编排者裁。
- **判不了 1**（92 `:79`）。

### E. 最可能藏着真活的票（按"甲的硬度"排）

1. **票 115**（3 枚）：它名下**连一张裁决表都没有**（`docs/evidence/s1/115-*.md` 零命中，而其 AC 区标题自己要求这张表）；
   2 枚甲（变异／门禁）至今只有实现方 Progress log 的自述数（`:251`/`:360`/`:444`）。
2. **票 92**（4 枚）：`(iii)` 变异被验收方明标〔仅自述，不背书〕（`92b-…md:137`），门禁读数绑旧树，
   还有一枚本腿判不了的负判据。
3. **票 113**（2 枚）：门禁当年那行就写着"ban#8 的 sha↔数字配对错了一格"（`113-…md:17`），不是干净通过。
4. **票 110**（2 枚）：新步上线把后面的步骤吃掉了（`f161c690` 标题自陈），后果今天还挂在 111／112／140 三张开放票上。
5. **票 104**（2 枚）：变异与回归读数都绑 `4d43447` 快照，且票面 `:68` 的编排者注记同族用例曾在 CI step 4 红过。
