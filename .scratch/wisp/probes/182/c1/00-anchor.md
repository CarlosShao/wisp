# 182-c1（二轮只读普查）起手锚 — 2026-10-08

- 派单口径：起手 HEAD **不早于** `601c2b18`（2026-10-08 10:09:48 +0800，本人现量）；起手闸门落笔时**复量** HEAD＝`17a54338a8b3cd98085456c0cd2af2ce0118ed69`（并行腿正在推进 dev，锚件只认这次现量，不认"必须等于"）。
- `git status --porcelain` 全仓行数：起手现量 **756**（第一次读为 753，两读之间 HEAD 与在飞件都在动）——**工作树本来就有别人的在飞改动（`design/**`、`.gitignore` 等），本程一个字都不动它们**。
- 分根起手读数（尺：`git status --porcelain -- cmd internal docs .scratch/wisp/issues tools scripts | wc -l` ⇒ **2**；`git status --porcelain -- design | head -5` ⇒ 有 ` D design/assets/*` 等在他程名下）。⚠ 后果写进各尺读数口径：产码根（cmd/internal/tools）近干净，本程读数≈HEAD；`design/**` 脏，凡读 demo 一律具名"工作树值"。
- ⛔ 零 go 命令；⛔ 不开窗口/不起 GUI；⛔ 不 push；⛔ 不碰 `docs/reports/pending-and-issues.md`。

## 本程三处写点（除此之外一律不写）
1. 工单面 `.scratch/wisp/issues/182-*.md` 末节**追加**（⛔ 不改写既有行、⛔ 不勾任何 `- [ ]` 框）。
2. 证据件目录 `.scratch/wisp/probes/182/c1/`（只用 `.md/.txt/.tsv`，⛔ `.out`；每把尺自带 `rc=N` 行）。
3. 以上两处的 commit（显式 pathspec，只 add 本程路径）。

## 查重尺 A（已读，凡它给过读数的堆本程只引用不重跑）
`.scratch/wisp/probes/167/c2/census.md`（锚 `0ce6cb91`，10-02 22:13）已量三堆，本程引用其尺与读数：
- **占用条（上下文 used÷total）**：`Snapshot`（`internal/panel/composer.go:57-92`）6 枚字段、`ComposerState`（`:235-270`）14 枚 JSON 字段里**无 used/total 任何一格**；`PumpSources`（`internal/panel/pump.go:121-188`）12 槽**无占用 reader**；分母候选 `loop.Budgets()`（`internal/agent/loop.go:253`）cmd/wisp 产码调用者 3 枚（`run.go:1056/:1062/:1109`）无一取 `.ContextWindow`；⚠ `loop` 是 `execute()` 局部变量，runtime 不握读口（其 §6.1）。
- **排队序号**：真源 `LiveApproval.Position`（`internal/agent/approval/queue.go:268-275`、`pending_read.go:104-122`）；断点＝`NativeVerdict` 无 Position 字段（`pump.go:77-91`）＋`liveVerdicts()` 未抄（`cmd/wisp/panel_pump.go:58-76`）；`Queue.Depth()`（`queue.go:133-137`）导出、pump 装配未接。
- **停止**：入向名册 6 枚＝`internal/panel/bridge.go:42-45` 四枚 `panel.*`＋`:66-67` 两枚 config，**无一与停止有关**；四枚"写好零生产调用方"名册（`Steer`/`RunningTask.Cancel`/`RequestWorkspaceSwitch`/`Options.Control`）见其 §3.1。
- 同件另一结论（引用）：**常驻面板与 resident runtime 快照 pump 之间今天没有装配连线**（其 §5.3，常驻腿 disp 链只有入向六门、不接快照 pump）。
- ⚠ 冲突复核义务：若本程任何复量与上列冲突 ⇒ 具名报回，不默默取其一。

## 查重尺 B（票面框已现读，AC#3 指认前不凭印象）
- 票 `145-*`：`- [x]` AC#1／AC#3／AC#4／AC#5；**未勾**＝AC#2（落地有源集）、AC#6（反向判据）、AC#2b（模型清单/档位，09-28 追加）。其 AC#1 行来源 `PLAN.md:3473-3488` 十四行——本程要复核行号是否漂。
- 票 `181-*`：`- [x]` AC#1–AC#5；**未勾**＝AC#6（门禁，欠 `TestComposerContractTypesMatchFrontend` 一枚红＋`./internal/risk/` 那族 C26 读数）、AC#7（`WorkspaceView.Rewritten` 生产者恒 false）。⇒ 票 181 的 git **状态显示**已落地（`internal/panel/git.go`＋`panel_pump.go` `gitView()`）。
- 票 `163-*`：`- [x]` AC#1；**未勾**＝AC#1b（名册差集：`shell.exec`/`shell.session` 生产注册表都不存在）、AC#2–AC#5。⇒ 终端地基＝163 未开实现格，本票只登记依赖。
- 票 `182-*` 本体：AC#1–AC#7 **七框全未勾**（本程交件也不翻勾，AC 勾选归编排者）。
