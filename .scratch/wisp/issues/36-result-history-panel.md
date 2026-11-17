# 36 — Result & history panel: chat stream, agent work states, sanitizer, transcript view

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 30-result-routing-d10, 35-panel-bridge-c17
**Parallel slots:** ≤2 sub-agents (A: chat stream + 14 agent-state visuals; B: history/transcript
+ purge/export + virtualization)
**Spec refs:** SPEC-08 §5.4, §17.5, F2 sanitizer layer, D10 panel surface, §14.4 transcript rules

## What to build
The main result surface: streaming chat view with the 14 agent-work-state visuals (thinking
band, SSE cursor, reasoning collapse, tool chips with risk badges, approval-wait, cancel/error/
stuck cards), plus the history & transcript browser over SQLite with privacy controls.

## Key constraints
- 14 states per SPEC-08 §5.4 table — all rendered (design/screens/chat.html is the visual
  reference; human sign-off required): thinking light-band + elapsed timer; streaming cursor +
  token counter; reasoning collapsed-by-default; tool chip = category icon + name (mono) +
  truncated args + status arc/check/x/ban + L0/L1/L2 outline pill; expanded chip = colored JSON
  + duration + correlationId; approval-wait amber pulse-once + depth; error card human-copy +
  expandable error_class; stuck card names the repetition; cost line mono tabular-nums.
- Markdown: react-markdown + remark-gfm; **rehype-sanitize whitelist mode** (tightened schema);
  `web.fetch`/`doc.read` outputs NEVER take the HTML path (plain text or sanitized markdown
  only); no `dangerouslySetInnerHTML` anywhere (CI grep); code blocks → Shiki static highlight.
- Approval interplay: L2 card inside panel shows full params + rules_hit + reject button +
  "点击悬浮球以批准" persistent guide with pulsing ball glyph — NO allow button (F2, §15#6).
- History: task list (30d) → task detail (transcript + tool calls) from SQLite via bridge;
  virtualized list (react-virtual) for 1000+ rows; per-item delete + purge-all + export JSON
  (privacy APIs from 29/04); transcript display honors §14.4 masking marker.
- Copy buttons (code/paths) go through clipboard tool path for taint bookkeeping.
- All loops respect the CPU rule: no animations when panel hidden (visibilitychange pause).

## Out of scope
- Approval queue multi-task view (38); palette (38); config editor (39).

## Acceptance criteria
- [ ] 14-state visual checklist against chat.html — user sign-off (screenshots committed).
- [x] XSS suite: `<img onerror>`, `<script>`, `javascript:` href, event-handler attrs in fake
      fetch/markdown content → neutralized; zero script execution (CSP report-only listener
      asserts zero violations).
- [x] dangerouslySetInnerHTML grep = 0; sanitizer whitelist config reviewed + committed.
- [ ] History: 1000-row fixture scrolls at 60fps; delete-one/purge/export verified against DB.
- [ ] Long-task stream: 200-chunk golden replay renders with merge (no jank), token counter
      accurate.
  〔**10-08 12:0x 编排者就地裁冲突，来源＝只读普查腿 `197-a1`（件 `.scratch/wisp/probes/197/a1/sweep.md`）；⛔ 本框保持未勾、⛔ 不改判据要什么**：那句 **"renders with merge" 与盘上已重裁的规则方向相反**——产码 `internal/panel/pump.go:397-409` 逐字 `Overflow TRUNCATES and never merges`，是**票 197 leg B 重裁**（判据 票 197 `:32`＋未勾 `AC#3 :39`；裁处 `docs/evidence/s1/197-subagent-stream-r2.md:41-54` §1「乙」；批准＝台账 **`A406`**；代码＝`94ea50f7` 删掉 `mergeOverflowLocked`，现 0 命中）。⇒ 落地时这一格只能按**截断**那形写判据；⛔ **任何程都不许为了对上"merge"这个词去改 pump**（那是把一枚已裁的规则倒回去，要改必须先经人工批准）。另两处同名不同链具名排除：`cmd/wisp/panel_pump.go:397` 与 `internal/panel/pump.go:397` **行号相同、文件不同、说的相反**，引号必须先报文件；`internal/subagent_197.go:274` 那句「把任务并成一枚子代理」是**用户文案**、与本链无关。〕
- [ ] Hidden-panel animation pause asserted (rAF counters).

## Progress log (append-only, newest last)

> **2026-09-26 进度追记（owner 把前端交 ZCode 直管后的可视壳进度，只追加不勾框）**：
> 本票的可视面已按 beautiful-ui 主题收编第四代重画，演示形态在 `?harness=1`
> 应用全观（左栏三段式/会话更改卡/右栏审查·终端·浏览器）与 `?harness=2` 陈列室。
> **本票 AC 的真数据 e2e 格仍等票 145/35/33 的字段与宿主，一格未勾、不由本程
> 勾**；可视壳的 commit 链见票 77 进度追记（`424ac84`…`5e23d99`）与
> `docs/reports/frontend-session-log-zcode.md` §1–§7。

> **2026-11-16 22:1x 腿 1（agent-95，GLM 舰队五票顺序跑第 1 票；铁律 4 前端归本编队）**：
> 工作状态面真接上了：`snapshot.run`（票 145 生产读面，cmd/wisp/run.go:1184/1281 装配、
> :2729-2792 喂值）第一次有 UI 读者（`frontend/src/components/result-history.tsx`）；
> 结果流完成态过 F2 白名单消毒（react-markdown+remark-gfm+rehype-sanitize，schema 收紧：
> 减 img/div/input、href 协议收窄、class 全剥）；git chip 从喂死空对象改真读 `composer.git`。
> 合同尺注册的 6 枚 Go-only 键（`internal/panel/composer_test.go:231-233`）**零造值**，
> 摘注册→panel.ts 补声明的同 commit 流程写在
> `docs/evidence/s1/36-result-history-panel-leg1.md` §4，归 Go 腿。
> 新牙 `npm run render:history36`（四族 XSS 过产码消毒、真数据回显、注册键零造值、
> 缺 reader 具名、git chip 走线）OK 0 失败；build/typecheck/lint/tokens:check rc=0，
> render 仪器 8 枚全绿，d22scan clean。AC 勾框一律留给非实现者；留白格与具名缺引擎
> 清单见证据件 §1/§3。

---

## 编排者收件·第一批（2026-11-17 09:2x，`9d2b873df`，件 `docs/evidence/s1/`＋仪器 `render:history36`）

**尺1** `9d2b873df` 在。**尺3 我自己复跑**：`npm run render:history36` **rc=0**（`4 XSS families neutralized through production MarkdownText, run-feed values round-trip`）；`build/typecheck/lint` 全 **rc=0**。
- **勾 AC#1**：四族 XSS 载荷（`<img onerror>`／`<script>`／`javascript:`／事件属性）走**生产 MarkdownText**（`result-history.tsx:114 rehypePlugins={[[rehypeSanitize, TIGHT_SCHEMA]]}`）全灭，仪器现跑 rc=0。
- **勾 AC#2**：`dangerouslySetInnerHTML` **用法 0**（全仓 grep 命中 1 处，是 `result-history.tsx:35` 注释里的自述，非用法）；白名单 `rehype-sanitize ^6.0.0`（`package.json:41`）＋`TIGHT_SCHEMA` 已入库。
- **留白（具名）**：AC#3 1000 行历史树（缺 SQLite 历史 reader＋`panel.view.request` 被 Q-50 封）、AC#4 token 计数（`ResultChunkView` 无字段）、AC#5 rAF 暂停（需真宿主窗）。

## 编排者追加·历史读面已落地（2026-11-17 13:3x，`3a62bed97`，件 `docs/evidence/s1/readface-wave2-20261117-r1.md`）

我复跑（载具）：`TestHistoryPageDoorAnswersTheRealTables` 0.14s／`TestHistoryPageScrollsOneThousandRowsWithoutOverlap` 0.56s（1005 行三窗滚动不重不漏）／`TestHistoryPageRefusesAnUnservableWindowByName` 0.13s ⇒ **全 PASS**。
- 面：`internal/memory/dao_history_page_36.go:104 HistoryPage`（真表 `task_log`+`tool_call`）→ `internal/panel/bridge.go:175 history.page` → `composer_dispatch.go:344` → `cmd/wisp/panel_history_wave2.go:78`。**Q-50 不影响**：走独立名册方法，`panel.view.request` 不复活的裁定（A238）照守。
⇒ 前端 AC#3 只剩"把它画成 1000 行滚动列表＋delete-one/purge/export 三支走既有 privacy 面"。
