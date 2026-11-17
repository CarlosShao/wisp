# 39 — Config editor GUI: sections ↔ structs round-trip, three-tier badges, security locks

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 05-config-model, 35-panel-bridge-c17
**Parallel slots:** ≤2 sub-agents (A: schema-driven form rendering; B: security-section flow +
restart/reload UX)
**Spec refs:** SPEC-08 §5.2, D6, D36, §14.5, design/screens/config.html

## What to build
The config editor as a pure VIEW over config.toml: schema-driven forms for every section,
three-tier effect badges, security-section lock affordances with L2 re-confirm flow, and
strict single-source-of-truth discipline (GUI holds nothing TOML doesn't).

## Key constraints
- Forms generated from the SAME Go schema structs (bridge `config.get/set`); GUI never owns
  defaults or extra state (D36 rule 4; Zod for input feedback only — Go stays truth).
- Three-tier badges: hot (applies now) / reload (subsystem reload note) / restart (banner
  "需重启生效"); after save, applied-tier feedback surfaced from the reload events.
- 🔒 sections (risk/fs/net/plugins): lock icon; loosening edits trigger the L2 re-confirm
  (native card), reject → keep old value; tightening applies hot; all directions logged.
  Danger-styled affordance, never a silent toggle.
- api_key fields: write-only inputs (never render stored values); set → creates/updates
  `api_key_ref` via SecretStore (06); placeholder shows last-4 only.
- Unknown-key handling: file edited externally with unknown keys → validation errors listed
  with line numbers (05 errors surfaced verbatim).
- Migration/mismatch: `schema_version` mismatch → guided flow, no silent overwrite.
- Round-trip guarantee: open → save unchanged → file diff empty (formatting excepted).

## Out of scope
- First-run wizard (firstrun screen lands with 41/14 polish; basic flow here); models page
  heavy UX (40 covers listing).

## Acceptance criteria
- [ ] Every D36 section renders from schema (golden snapshot per section); edit → save →
      hot/reload/restart behavior matches tier (events observed).
- [ ] Security loosening → native L2 confirm → accept applies + logs; reject keeps old
      (e2e with 21's card).
- [ ] api_key write-only + last-4 display + ref creation verified; plaintext never rendered
      or stored by GUI.
- [ ] Round-trip no-op save produces no diff; external-edit + unknown key shows line-numbered
      errors.
- [ ] Visual sign-off vs design/screens/config.html (lock icons on security sections).

## Progress log (append-only, newest last)

> **2026-09-26 进度追记（owner 把前端交 ZCode 直管后的可视壳进度，只追加不勾框）**：
> 本票的可视面已按 beautiful-ui 主题收编第四代重画，演示形态在 `?harness=1`
> 应用全观（左栏三段式/会话更改卡/右栏审查·终端·浏览器）与 `?harness=2` 陈列室。
> **本票 AC 的真数据 e2e 格仍等票 145/35/33 的字段与宿主，一格未勾、不由本程
> 勾**；可视壳的 commit 链见票 77 进度追记（`424ac84`…`5e23d99`）与
> `docs/reports/frontend-session-log-zcode.md` §1–§7。

---

## 编排者追加事实（2026-11-16 20:3x，来自只读审计腿 `census-unwired`，件 `docs/evidence/s1/census-unwired-20261116-2017.md`）

**这一段是"现读坐标"，比票面正文新；按它做，⛔ 按票面旧行号。**
- **`app.theme` 今天零读者**：`App.Theme` 全仓 22 命中全在 `internal/config` 自转圈＋测试；`internal/ball/renderer_windows.go:106 SetTheme` **零调用者**；`ball.Options.Theme`（`ball_windows.go:93`）装配根**不填**（`resident_ball_windows.go:488-500`）；写侧也不在名册（`config_write_313.go:201-209` 六册并集无它）。⇒ 本票据 GUI 时**要么接上球侧消费者（`SetTheme` 有人调），要么这一行明确标"未接引擎"**，⛔ 显示"即时"。
- **`panel.font_size` 四跳已通（写/读/快照/装配）但零消费引擎**：`config_write_wave.go:40`＋`panel_config_store.go:431-440`＋`settings_appearance_wave.go:179`＋`run.go:1272-1274`；`Panel.FontSize` 全仓 6 命中全在 `internal/config`＋测试，`panel_host_windows.go` **0 命中**。⇒ 要么接到页面字号，要么标明缺引擎。
- **`panel.opacity`／`panel.animations` 不是 config 键**（页面本地 CSS 变量／`html[data-motion]`）⇒ 本页若展示它们，**标明"仅本次会话、不落盘"**，⛔ 让用户以为存了。
- **★ 徽章判据是假的**：`frontend/src/components/config-screen.tsx:411` `const live = row.blocked === ""`，`blocked` 是 `APPEARANCE_ROWS`（`:152-162`）**硬编码字符串** ⇒ 徽章与 Go 到底有没有读者**毫无关系**。本票必须把它改成**读真状态**（有生产读者/引擎才"即时"），并给一枚行为牙：某行读者被摘 ⇒ 徽章翻"未接线"且测试红。

---

## Progress log（续，2026-11-17 08:4x，本票前端写腿 agent-100 · 编排者 20261116 五票批次的第 4 张）

> **AC 勾格状态：5 格 0 勾、5 格留白**（原因逐格在下列表，均因缺 Go 侧读面/写面，前端⛔造数据）。
> 本轮落地的是**编排者追加事实一节的四条 census 指令**（徽章读真状态＋三行诚实标注）＋渲染仪器重钉。

### 本轮落地（frontend/**，三件）

1. **★徽章判据换成读真状态**（census ★）：`frontend/src/components/config-screen.tsx`
   `Row.blocked`/`control` 硬编码字符串删除，改 `RowSource` 三态模型
   （`document`／`config{field,section,tier,consumerMissing?}`／`dormant{reason}`）；
   徽章由 `rowBadge(source, sections)` 从**快照节是否真的到达本泵**推导
   （Go 侧 nil 指针＝键不上线，节缺席＝"这台泵没有读者"）：
   - 未喂到的 config 行 ⇒ 徽章翻 **未接线**＋控件区 等快照（census 牙：读者被摘⇒徽章翻）；
   - `panel.opacity`／`panel.animations` ⇒ 徽章 **仅本次会话**（census：非 config 键，⛔让用户以为存了），desc 加"仅本页样式，不落盘"；
   - `app.theme` ⇒ 保持 **未接线** 并具名缺引擎："SetTheme 暂无调用者，球侧装配根也不填 Theme"（census 原话点名）；
   - `panel.font_size` ⇒ 即使节已到泵，徽章 **缺消费引擎**（橙警示），desc 写明"写入会落盘，但还没有引擎消费这个键，暂不生效"（census：四跳已通但零消费引擎）；
   - 已喂到的行徽章保持原 tier 词（能改／立即/重启／关窗/重启生效）。
2. **渲染仪器重钉**：`frontend/scripts/render-config-screen.tsx`——census 牙**双向**钉死：
   无读者渲染＝6 枚 config 行徽章全翻 未接线＋零 tier 词＋页本地行 仅本次会话×3（含页脚）＋主题行具名 SetTheme；
   全喂渲染＝tier 词到位（能改×2／关窗重启×2／立即重启×1）而 font_size 仍 缺消费引擎；
   外加源级钉：行表六对 `field,section` 配对不许漂、旧 `blocked === ""` 判据不许回来。
   **仪器 exit 0 全绿**（约 24 枚 check）。
3. `frontend/src/components/harness/settings-page.tsx:62` 注释同步（旧行数声明过期）。

### 门禁读数（2026-11-17 08:3x，本机 windows，node24）

`typecheck` rc=0｜`lint` rc=0（38 警告 0 错误，全为存量）｜`build` rc=0｜`render:config-screen` 仪器 exit 0｜
`render:nav` OK（9 屏/70 冻结名）｜`render:composer` 7 态全画 rc=0｜`probe-ballsize-render` OK｜
`d22scan` rc=1 但 **frontend/** 零发现**（4 条全在 `cmd/wisp/approval_reply.go:464/:379`、`internal/session/ticket49_lifecycle_test.go:172`、`internal/tools/ticket49_limits_test.go:139`——本轮 cmd/wisp 另一写腿的在飞文件，未触碰）；`frontend/dist` 仍只有 `.gitkeep`（dist 产物 git-ignored，未入库字节）。

### AC 逐格留白（⛔前端造数据，需 Go 侧字段/file:line）

- **[ ] AC#1**（D36 全 section 按 schema 渲染＋golden snapshot＋tier 事件观测）：**留白**。缺 Go：①config.get 的结构化 payload 现量只有 票313 的 模型/语音两节＋快照推送的 ball/panel 六节（`internal/panel/config_handlers.go`、`composer.go`），**没有** D36 全 section 树的 schema+值 reader；②没有"reload/restart 事件"推送字段。要加的：config.get 增 `sections[]{key,value,tier}` 全节读面＋reload 事件推送（file:line 落点由 cmd/wisp 腿定）。
- **[ ] AC#2**（安全节 loosening → 原生 L2 卡 → accept/reject）：**留白**。缺 Go：risk/fs/net/plugins 四节既无读面也无写面（快照零推送），票21 的原生卡流程无 config 侧入口。
- **[ ] AC#3**（api_key write-only＋last-4＋ref 创建）：**部分现量**：模型页已是 write-only（`frontend/src/components/model-settings.tsx:606-625` type=password、值不进 state、发出即清）＋`apiKeyRef` 展示（`:597-598`）＋凭据写走 SecretStore；**留白部分＝last-4 占位显示**：`internal/panel` 刻意"红线尾不上线"（`frontend/src/lib/config-view.ts:44-45` 注释），与票面"placeholder shows last-4"冲突——要显示 last-4 需 Go 改 payload（改契约文字，须人工批准），未动。
- **[ ] AC#4**（round-trip 无 diff＋外部编辑未知键带行号报错）：**留白**。缺 Go：无"校验错误＋行号"的读面（config.get 只答冻结句＋节 payload，无 `errors[]{line,message}` 字段）。
- **[ ] AC#5**（对照 design/screens/config.html 视觉签收＋安全节锁图标）：**留白**。安全节没有读面就画不出锁行；且 design/** 与视觉签收属验收者（裁决者≠实现者），不在本腿射程。

### 回话摘要（给编排者）

> commit hash 见 git log；勾格 0/5（全留白，原因＋要加的 Go 字段如上）＋census 四条全落地；新牙＝`render-config-screen.tsx` 的 census 双向徽章牙＋行表 field/section 配对钉；npm build/typecheck/lint rc=0。


---

## 编排者收件·第一批（2026-11-17 09:2x，`021bb9832`，仪器 `render-config-screen`）

**尺1** `021bb9832` 在。**尺3 我自己复跑**：`vite build --ssr scripts/render-config-screen.tsx` ＋ `node …` 均 **rc=0**（`render-config: OK, badge reads the pump (census tooth both directions: readerless rows badge 未接线 + 等快照, fed rows badge their tier …)`）；build/typecheck/lint 全 rc=0。
**census 四件已落地（我 grep 复核）**：①`config-screen.tsx` 里 `row.blocked === ""` **判据已删**（全仓零命中；只有别的文件的 `blockedOnApproval`）；②`app.theme` 标「未接线」并具名（`SetTheme` 无调用者）；③`panel.font_size` 标「缺消费引擎」；④`panel.opacity/animations` 改「仅本次会话」。
**0/5 格勾（诚实）**：五格全卡在 Go 读面 —— AC#1 缺 `sections[]{key,value,tier}` 全节树＋reload 事件；AC#2 risk/fs/net/plugins 四节零读零写＋无 L2 入口；AC#3 `api_key` last-4 与 `config-view.ts:44-45`「尾永不上线」契约冲突（**须机主批契约**）；AC#4 缺 `errors[]{line,message}`；AC#5 视觉签收（机主）。

---

## 编排者追加·读面已落地（2026-11-17 10:3x，`b1d2162dd`，件 `docs/evidence/s1/readface-wave1-20261117-r1.md`）

我复跑：`internal/config` 新牙 ok 0.133s；载具 `Test39ConfigStoreReadFaceSectionsErrors` **PASS 0.01s**。**Go 读面已就位，前端可以接**：
- **AC#1**：`internal/config/sections_456.go:45 SectionRows`（全 D36 树反射走查，含 `[plugins]` 动态子表）→ `internal/panel/config_handlers.go` payload `sections[]{key,value,tier}` → `cmd/wisp/panel_config_store.go:147`；**reload 事件** `config_handlers.go:579`（`wisp-config-reload:` marker；accepted 写才发、被拒零事件）。tier 词取自 `TierRegistry`（hot/reload/restart/locked）。
- **AC#4**：`internal/config/loaderrors_456.go:43/:54/:108`（`errors[]{line,message}`，未知键/类型错带 1-based 行号；语义错 Line=0）。
⇒ 前端只剩"把 sections 画成每节一张表＋tier 徽章、把 errors 按行号列出来"。

---

## Progress log（续，2026-11-17 11:4x，本票前端消费腿 agent-109 · 读面消费波）

**本轮交付（frontend/**＋internal/panel/composer_test.go 一处；证据件 `docs/evidence/s1/readface-consume-20261117-r1.md`）：**

1. **AC#1 页半落地**：`frontend/src/lib/config-view.ts` 解析 `sections[]{key,value,tier}`（omitempty 缺席＝腿未答，不画成空文件）＋`wisp-config-reload:` 事件解析；`frontend/src/lib/config-actions.ts` 加 reload 事件总线，五个 config.set 写助手都在 reply 上检测（被拒写零通知）；新件 `frontend/src/components/config-editor.tsx`：每节一张表＋tier 徽章逐字（hot/reload/restart/locked，"" ⇒ 未登记档位，locked 红 tint）＋`ConfigEditorSection`（自读 config.get＋订阅事件整页重读）；挂外观类脚底（config-screen.tsx 尾）。
2. **AC#4 页半落地**：`ConfigErrorList` 按行号列（line=5 ⇒ 「第 5 行」；Line=0 ⇒ 「（无行号）」，不折第 1 行）；errors 缺席 ⇒ 「清洁或未检查分不出」句，不编错误。
3. **牙**：`render-config-screen.tsx` 新增消费牙——18 节⇒18 表、摘一节⇒表消失且旧值零残留、tier 词逐字计数（hot×12/reload/restart/locked×4）、未登记档位、errors 两态、reload 事件逐字＋被拒零通知＋五助手 notify 计数＝5。
4. **[ ] AC#1 / AC#4 勾格留白（裁决者≠实现者）**：页半＋牙已齐，等对抗验收；AC#2（安全节 L2 流）/AC#3（last-4 契约冲突）/AC#5（视觉签收）不动。

**门禁读数**：build/typecheck/lint rc=0；`render:config-screen` rc=0（vite --ssr＋node 手动跑，无 npm script 历史如此）；`go test ./internal/panel` rc=0；d22scan rc=0（frontend 零发现）；⛔ 未 push；dist 未入库（.gitkeep 唯一跟踪件）。
**回归牙记录**：fixture 里 [panel] 行被渲染门门禁（composer_test.go routeLiteralRe 禁 renderer 源里任何 "panel.*" 字面量）抓红一次，已换 demo 占位行并在仪器注释里写明原因——门禁工作正常，未放宽。
