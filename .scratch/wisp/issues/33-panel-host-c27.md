# 33 — Panel host: C27 PanelManager singleton, WebView2 window, embed.FS serving

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 07-ball-state-machine-core, 12-cli-text-path-s1-gate
**Parallel slots:** 1
**Spec refs:** SPEC-08 §5.1, D29, C27, D32 panel rows, D42#11, S5

## What to build
The `panel` module host side: singleton PanelManager owning at most ONE WebView2 window per
session (hide-don't-destroy), embed.FS resource serving via `AddWebResourceRequestedFilter`
(no localhost server), cold/hot show paths, focus return, and the no-runtime fallback
declaration consumed by 37's native card.

## Key constraints
- `jchv/go-webview2` (pure Go); one window per session; hide/show instead of destroy within a
  session; destroy at session dispose (31 hook). Cold show ≤1500ms / hot show ≤200ms — measured
  and recorded (D32).
- Resources: `//go:embed assets/web/dist` + `AddWebResourceRequestedFilter` mapping to embedded
  FS; NO localhost HTTP listener (D29; D21 anti-pattern).
- WebView2 created on the shared `ui-sta` STA thread (D38a); environment notes from spike (02)
  recorded — single-window strategy sidesteps Environment sharing entirely.
- Focus: only the panel takes focus when shown; on hide/close → focus returns to the recorded
  previous foreground window (D29).
- Multi-task partition hook: render context keyed by correlationId (surface built at 36/38).
- Runtime-missing detection: creation failure → `panel.unavailable` event + no-crash; L2 native
  fallback flag set (37 consumes); guide user to Evergreen install, never auto-download.
- CSP header/meta injection point for the embedded page (actual CSP value enforced at 35/36 —
  provide the injection hook + default strict CSP now: connect-src 'none' etc., SPEC-06 §9).

## Out of scope
- Frontend bundle contents (34); bridge protocol (35); panel pages (36–40).

## Acceptance criteria
- [ ] Show/hide/destroy lifecycle: second show within session reuses window (process-tree
      child-count stable); session dispose destroys + WebView children exit ≤2s.
- [ ] Latency: cold ≤1500ms / hot ≤200ms (10-run P50/P95 into SLO appendix).
- [ ] Embedded resource serving works offline; no listening socket exists (netstat assertion).
- [ ] Focus round-trip: editor → panel → hide → focus back in editor (automated + manual).
- [ ] Runtime-missing fixture (rename/mask loader) → fallback event, app alive, native L2 flag on.
- [ ] CSP present on served documents (response inspection) with connect-src 'none'.
- [ ] **AC#7（从票 35 快照泵 r1 验收转入，2026-09-25，来源 `35-panel-snapshot-pump-r1-accept-r1.md` F-PUMP-4）**：
      面板真的能看到快照之后，`SnapshotPump.Snapshot()` 那份包**跨两个瞬间**的形状要有一条会响的检——
      它不持泵锁地先后读 `Verdicts()` 与 `Results()`（`internal/panel/pump.go`），而今天四个读数点各自有锁、
      `SnapshotPump.mu` 只护计数器，所以 `-race` 两包全跑**是干净的**（实测 `ok 4.809s`／`ok 75.715s`、0 DATA RACE）。
      ⚠ **这一格今天不许在票 35 上开**：没有 hook 可测，硬开只会造出一枚**永远绿的假钉**（本仓那族叫"恒真判据"）。
      它要有 hook 的时刻，正是"有人在生产里读这份包的全量字节"那天——也就是本票把窗口接上的那天。
      判据：本票落地 `Marshal()` 的第一枚**生产**调用者之后，造一发"两个触发点之间状态变了"的用例，
      它必须能红（把 `Snapshot()` 改成持锁一次性快照 ⇒ 该发由红转绿）。
- [ ] **AC#8（从同一份验收转入，来源 §7 第 3 条 ＋ F-PUMP-6）**：本票接上出站通道时，**顺手量一次
      "哪些泵里的分支今天零执行"**并给它们用例或删掉。今天已知的两枚：
      ⓐ `cmd/wisp/panel_pump.go:186-192` 那条"摘要超过 440 字符就丢 ids"的分支，在 13 发变异＋工作树那发里
      **从未被走到**（落盘记录的 `pending=` 字段只出现过一枚值）；
      ⓑ `cmd/wisp/run.go:738` 那枚 `case agent.EvToolStart:`（`run.go:741` 里泵的新触发点 `changed = true`）
      **今天不可达**——`EvToolStart` 全仓**零枚发射者**（编排者 17:4x 独立复量：非测试引用恰 3 处＝定义
      `internal/agent/sink.go:25` ＋ 它的注释 `:24` ＋ 这一处消费；没有任何一处 `emit`）。
      判据：这两枚分支各要么有一发"拿掉它就红"的用例，要么被删；**不许留着让人以为它们在守什么**。

## Progress log (append-only, newest last)
- 2026-09-25 17:5x（编排者）：**加 AC#7／AC#8 两格，都是从票 35 快照泵 r1 的验收表转入的，不是本票新造的活。**
  转入口径写在这里：这两格**在票 35 上今天开不出来**——不是漏做，是**没有可测的 hook**，硬开会得到一枚
  永远绿的假钉（本仓那族缺陷有名字："恒真判据是一类新假绿"）。它们的 hook 恰好在**本票**把
  "Go → 面板"那根管子接上那天出现（AC#7 要 `Marshal()` 的第一枚生产调用者；AC#8 要真有人消费落盘摘要）。
  ⚠ 本票的 `Blocked by` 那行**没变**、`Out of scope` 那节**没动**（bridge protocol 35 仍在本票范围外）——
  这两格买的是"接上那天要顺手量"，不是把 35 的活搬过来。
  出处逐字在 `docs/evidence/s1/35-panel-snapshot-pump-r1-accept-r1.md`（F-PUMP-4／§7 第 3 条／F-PUMP-6）。
  **撤销口令：「撤 33 AC#7 AC#8」**（撤＝只删这两格＋在本节留一行撤销记录，不改写上面任何一行）。
- 2026-09-28 15:1x（只读设计核 `33-a1`，**零产码、AC 框一格未勾**）：交件全文
  `docs/evidence/s1/33-inbound-hop-design-a1.md`。本节只记四句结论，**不改写上面任何一行**：
  ① **"网页点一下 → Go 收到"这一跳被拆成 10 环（H1–H10），其中 H4–H9 六环一环都不引用 WebView2 符号**
  （现量：`grep -c webview go.mod` = 0；`ParseComposerRequest` 的签名就是 `raw string`，`bridge.go:84`）。
  ⇒ **票 186／187 不是"整体等地基"，是能先动 Go 侧那一半**（最小入向那片 = 1 枚新产码文件 + 1 枚同包测试 +
  零依赖 + 有 ubuntu CI 分母 + 零批准）。真宿主那片（H2/H3/H10 的手段）今天**在 windows scope 零用例分母**，
  要它就得先拍"改门禁形状"或"接受该腿分母＝本机人工"。
  ② **一枚规格级缺口，只报不填**：SPEC-08 §5.1 那五条宿主责任（`:143-154`）点名了生命周期／资源服务
  （`AddWebResourceRequestedFilter`，是**出向**）／无状态／Runtime 缺失／correlationId 分区，
  **没有一条写"消息接收"**；本票票面 `grep -nE "WebMessage|postMessage|inbound|receive"` = **0 命中**，
  而票 35 把 `bridge transport + dispatch`（`:7` slot A）划给自己、可 transport 要的那只手长在本票的窗口上，
  本票 `:33` Out of scope 又明写 bridge protocol 归 35。⇒ **入向接收器今天介于两票之间、两票都没逐字认领**；
  责任矩阵少一行是本跳断口的规格根因（补哪一行由编排者定，本程不动票面）。
  ③ **C17 白名单"定到哪一步"的现量比"待定稿"更难看**：规格侧 SPEC-08 `:163-174` 那 12 枚 invoke 方法名
  在**非测试 Go 代码里 0 命中**，代码侧 `bridge.go:42-45` 那 4 枚（`panel.mode.request` /
  `panel.workspace.request` / `panel.attachment.add` / `panel.message.send`）在 **`docs/specs/**` 里 0 命中**
  ⇒ **两条平行线，交集为 0**；节标题自带 `【SPEC 提案，S5 定稿走契约批准】`（`:156`）。`C24 GojaHostAPI`
  初始集一个名字都还没定（产码命中只有 `internal/config/manager.go`、`schema.go` 的配置字段）。
  ④ **P11（冷拉起 >2s 那条待定项）的现量其实已经在仓里**：`docs/evidence/s0/02-spike-report.md` §3.4
  两 run 的 cold P95 = 1041.6 / 1256.4 ms（**均未越过 2s，也均在 D32 冷 ≤1500ms 内**），唯一越过 1500ms 的是
  "进程内首次 create（竞争态）"那枚 1808.9 ms；而入向这一跳自己的传输成本有一枚可算的代理指标：
  **hot 浏览器往返 p50 = 3.5 ms**（量的就是 `w.Bind` JS→Go 那一趟，`scripts/spike/webview2-latency/main.go:118-134`）。
  ⚠ 两 run 同机、2026-09-19 未复跑，`grep -n "P11" docs/reports/pending-and-issues.md` = **0 命中**
  ⇒ 这条待定项在真相源台账里**没有对应 `A##`**，自偿要编排者先入册。
  ⚠ 本程另发现**一枚对本票不利的依赖现量**（利好写腿）：go-webview2 已被仓内那枚独立模块
  `scripts/spike/go.mod:9` 钉过（`v0.0.0-20260205173254-56598839c808`，含 indirect `go-winloader`）且
  `scripts/spike/bin/webview2-latency.exe` 真编译过 ⇒ "这枚依赖在本机能不能解析链接"**S0 已答**，
  本票落地时只剩"并进根 `go.mod` + `go.sum` + `deps.toml` 许可登记"三处动作（今天三处都 0 命中）。
  **撤销口令：「撤 33 progress 33-a1」**（撤＝只删本段＋在本节留一行撤销记录）。
