# 191-readonly-survey-r1 — 票 191（面板内置终端）只读现量 · 2026-11-20

**角色**：只读设计核（审计路由），实现者不在场；本件是编排者对腿报告的**摘要＋自己复核过的那几条**。
**口径**：凡是标「(我复核)」的行，是编排者现跑 `grep`/`sed` 亲自看过的；未标的来自腿报告 file:line，尚未逐行复核。

## 一、AC#1 现在有哪枚工具面能"跑命令并流式回显"
- 结论：**名册级候选 2 枚（`shell.exec` / `shell.session`），能流式回显的 0 枚**。
- 注册点：`internal/tools/shell.go:235-243` `BuiltinShellEntries` 一次注册两枚；**只在 `[risk] shell_enabled=true` 才注册**（`shell.go:236-238`）。
- 无流式的机制：`shell_exec.go:112-116` stdout 挂 `guardedBytes` 整段缓冲，`:131-161` 只在 Wait/ctx.Done 才读；`shell_session.go:219-236` 有读协程但 `:296-352` 只在见到 marker 时整段取出；`shell_exec.go:121-123` 唯一 onUpdate 只发"已启动"，`shell_session.go:431` 零 onUpdate。

## 二、AC#2 命令输出进面板的管子与外泄面
- 管子＝**只有"工具结果"这一条**（一次性，不是字节流）：`cmd/wisp/run.go:2942` `runFeed.ToolEnded` → `internal/panel/run_feed_145.go:93-100` → `internal/panel/pump.go:691-698` `snap.Run.tools[].resultText`。
- **(我复核) 显示层不过洗两处**：`internal/panel/composer.go:313` `Results: results` 与 `pump.go:697` `snap.Run = p.src.Run()` 都在 `composer.go:319` 唯一的 `scrubComposerForDisplay(composer)` 之外 —— 今天只有 Composer 段被洗（`display_scrub_216.go`），`Results`/`Run` 不洗。
- **(我复核) C25 敏感源名册里没有 shell**：`internal/risk/provenance.go:113-117` `sensitiveSourceTools = {SrcFSRead, SrcSearchContent, SrcClipboardRead, SrcSystemGet, SrcWebFetch, SrcDocRead, SrcScreenCapture, SrcTranscript, SrcTaskOutput, SrcPanelAttachment}`；`grep -n "shell\." internal/risk/*.go | grep -v _test` 只有 `rules_shell.go:14` 的 R6 注释（强制 argv 向量），没有 SrcShell。
- 推论（腿报告）：shell 输出里带的绝对路径/token 不进污染索引 ⇒ 后续出向调用 `Inspect` 扫不到 ⇒ R4 不响。**这与票 175 修掉的洞同型（175 只补了 `task.output` 一枚）**。
- 注意：给名册加源＝改 `SPEC-06 §5` 的清单（该文件里 `SrcPanelAttachment` 一条即因票 453 挂 `Q-456` 待机主）⇒ **属"改契约文字"，须机主批**，本条按 Q 登记，不自行实施。

## 三、AC#4 取消与超时语义
- 任务取消＝有（`internal/tools/cancel.go:21-94` + `bridge.go:720-730` + `cmd/wisp/run.go:1441`）。
- 超时＝真 deadline（`internal/tools/shell.go:145-152` 30m → `bridge.go:836-841` `context.WithTimeout`），非时钟差。
- **关面板不等于取消**：`cmd/wisp/panel_host_windows.go:1287-1329 Hide()` / `:1345-1351` / `:1368-1379 Destroy()` 三处都不碰 task ctx；真取消在进程退出（`resident_windows.go:533-535` → `resident_approval_windows.go:973-1001`）。
- 面板按取消＝有但缺一跳（`MethodRosterCancel="roster.cancel"` `internal/panel/bridge.go:266` → `composer_dispatch.go:493-497` → `cmd/wisp/panel_roster_cancel_wave2.go:119-149`）。

## 四、AC#3 能否钉"面板送进来的命令走同一套 risk 门控"
- 有现成接缝可钉：`internal/tools/bridge.go:371-372` 唯一 `Assess` 点 + `:1227-1237 assessorFor`；名册级工具都过桥。
- 但要钉"面板侧来源不得绕过"，还需一条"来源＝面板"的判据（今天 shell 工具没有来源标记）。

## 五、风险清单（实现时动哪些包）
- `internal/tools`（流式回显：把 `guardedBytes`/`guardedText` 变成增量交付 + onUpdate 真发）＋`internal/panel`（新快照字段/流通道）＋可能 `internal/risk`（若同时补 C25 名册，属契约）＋`cmd/wisp`（装配）。**未越出** tools/panel/agent 范围，但 `internal/panel` 是共享编译面 ⇒ 每波只能一条腿。

## 六、结论与建议
- **AC#1 甲支建议**：先做 `shell.session` 的**增量交付**（已有读协程与 marker 抽取，缺的只是"边读边发"），而不是新造第三枚工具；`shell.exec` 保持一次性输出（短命令语义）。
- **另一件事与票 191 解耦**：上述 C25 名册缺 shell 的洞建议**单独立票**（属 SPEC-06 契约文字，得机主批），不塞进 191。
