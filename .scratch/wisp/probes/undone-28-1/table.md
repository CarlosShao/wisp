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

## 10. 票 110 —— `110-no-ci-step-runs-internal-winsec-done.md`（2 枚）

## 11. 票 113 —— `113-posix-platformverifypplacement-has-no-link-leg-done.md`（2 枚）

## 12. 票 115 —— `115-seal-notices-carry-the-resolvers-answer-while-the-cases-compare-caller-spelling-done.md`（3 枚）

## 13. 票 118 —— `118-winsec-test-hardening-from-ticket104-acceptance-done.md`（1 枚）

## 附. 与本腿分母不一致的证据 / 判不了的格

（待填）
