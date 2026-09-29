# undone-28-1 —— 结案票残余 28 枚未勾判据逐枚判性表

- 盘点腿代号：`undone-28-1`（只读盘点，不产码、不改票面、不翻勾）
- 锚点：本腿现跑 `git rev-parse --short HEAD` = `fac60ad4`，分支 `dev`
- 尺：`grep -c -- '^- \[ \]' <file>`（⛔ 不用 `grep '^- \['`，进度戳 `- [2026-09-29 17:5x +08] ...` 会造成约 1,600 枚假阳性）
- 分母复核（本腿现跑，`.scratch/wisp/issues/`）：78 张 `*-done.md`，未勾合计 **28 枚／13 张票**，
  逐票计数＝`92(4) 07(4) 80(3) 115(3) 97(2) 113(2) 110(2) 105(2) 104(2) 89(1) 63(1) 118(1) 11(1)`
  ⇒ **与编排者 09-29 17:5x 的 28 枚／13 张票完全一致**（逐票分布也一致）
- 四分类定义（本表用字）：甲＝残余真活 / 乙＝措辞残留（不是判据） / 丙＝已归口别的票 / 丁＝今天其实已满足、只是没人翻勾

## 0. 总账表

| # | 票号 | 行号 | 判据名（逐字行首片段） | 结论 | 凭据出处 |
|---|---|---|---|---|---|

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

> 本腿现量：un = 1（`:53`）；票面 8 枚 AC 里 7 枚已勾。

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

## 4. 票 80 —— `80-blacklist-overrides-never-wired-to-gate-done.md`（3 枚）

## 5. 票 89 —— `89-0600-is-decorative-on-windows-acl-for-private-data-done.md`（1 枚）

## 6. 票 92 —— `92-panel-composer-mode-attachments-workspace-done.md`（4 枚）

## 7. 票 97 —— `97-dead-strict-param-and-the-comment-that-invents-a-caller-done.md`（2 枚）

## 8. 票 104 —— `104-sealfile-silently-drops-inherited-grants-done.md`（2 枚）

## 9. 票 105 —— `105-c26-rewrite-account-has-no-production-reader-done.md`（2 枚）

## 10. 票 110 —— `110-no-ci-step-runs-internal-winsec-done.md`（2 枚）

## 11. 票 113 —— `113-posix-platformverifypplacement-has-no-link-leg-done.md`（2 枚）

## 12. 票 115 —— `115-seal-notices-carry-the-resolvers-answer-while-the-cases-compare-caller-spelling-done.md`（3 枚）

## 13. 票 118 —— `118-winsec-test-hardening-from-ticket104-acceptance-done.md`（1 枚）

## 附. 与本腿分母不一致的证据 / 判不了的格

（待填）
