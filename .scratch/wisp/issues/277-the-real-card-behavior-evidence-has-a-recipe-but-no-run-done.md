# 票 277 — 「常驻进程会不会真举出一张审批卡」的行为凭据：配方齐了，那一发没有跑

**立票**：2026-10-08 19:3x 编排者（来路＝台账 `A725` 举卡链量清／`A732` 裁定／`A735` 配方收档）
**性质**：`AC#7`（票 246）的勾**只覆盖到「起一条任务管线」**（台账 `A732` 已裁）；本票承接**另一半**——**真举出一张卡**的**行为凭据**。⚠ 本条**不是契约问题**、⛔ 不需要 owner 拍板。
**Status:** **done**（2026-10-09 10:0x 编排者结案：三格全勾。凭据＝腿 `277-r1` 一发＋**编排者自己第三发独立复跑**（rc=0／2 PASS／红句 0 命中，corr 与腿那发不同＝真跑）。⚠ 结案范围**只到〔接缝注入栏〕**：〔真凭据栏〕（真控制台＋真模型凭据）今天仍空，缺的三件逐字登记在 AC#2 与 `probes/277/r1/`；B 道条件齐之前 ⛔ 不派。台账＝`A753`）

## 现量（引用前先重跑）

- 链的静态形状全接：`cmd/wisp/main.go:66`→`resident_windows.go:260` 起任务源→装配→`run.go` 注门/交 bridge→`internal/tools/bridge.go:424/:443/:458`（有效级 L1 无会话授权 或 L2 才问）→`internal/agent/approval/gate.go:294/:536` 两处 `ui.Prompt`→常驻侧 `ballCardUI.Prompt`（`resident_approval_windows.go:864` 一带）。
- 卡出现的可观测信号真身：`resident_approval_windows.go:885` `fmt.Printf("wisp: 卡片挂起：%s %s（编号 %s）
", p.Level, p.Tool, p.CorrelationID)`＋`:880` 审计行 `approval: 常驻进程显示一张确认卡片`。
- 配方（可照做）：`.scratch/wisp/probes/card-proof-prep-1/01-recipe.md`（A 道＝winlive 台件注入文本，今天跑得动；B 道＝真控制台 `task` 动词，今天跑不动：`%APPDATA%\wisp` 无凭据＋`build/wisp.exe` 旧于 HEAD）。

## 要建什么

- [x] **AC#1 真机跑 A 道**：照配方 A 道跑一发，逐字收"卡行"（`:885` 那句）与审计行；**判红绿按配方 §4**（绿＝卡行＋`esc_borrowed=true`/`orb_state=Confirming` 同 corr；红＝入口/管线被拒或超时无卡；判不了＝Esc 被占/无桌面）。⚠ 本票**不要求**面板在场（该腿无面板，观测面＝stdout＋审计＋winlive `t.Logf`）。**（2026-10-09 10:0x 编排者判＝成立（绿），⛔ 行号已就地打旧见本节末"更正"。腿 `277-r1` 交一发（corr `0455ee9f…#call_246r2`，另两发 run1 真跑／run2 是 `go test` 缓存回放不算测量，三发全留档）；**编排者自己第三发独立复跑销账**：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -tags winlive ./cmd/wisp -run 'TestLive246ResidentPipeline|TestLive246ExitCancels' -count=1 -v` ⇒ **rc=0**，卡行逐字 `resident_task_source_live_246_windows_test.go:145: AC#7 LIVE steps 1-2: wisp: 卡片挂起：L1 fs.write（编号 900bb089-cbbd-42d8-9ed8-48dd33ff78df#call_246r2）`（corr 与腿那发不同＝真跑）；`--- PASS` 2 枚（17.48s／25.49s）、`=== RUN` 2 枚、**FAIL/SKIP 0 枚**、配方 §4 那七句红/判不了逐句扫 **0 命中**（`never raised a card`／`任务入口未启用`／`任务管线未装配`／`不被受理`／`卡片无处呈现`／`somebody else owns`／`observer rig is blind`）；读数件＝`probes/orch-277/00-own-rerun.txt`。⚠ **两栏分写**（配方 §4 灰条要求）：〔接缝注入栏〕**成立**＝真进程＋真球窗＋真审批门＋真 `fs.write`→L1→门→`ballCardUI.Prompt` 举卡并被注入的裸 Esc 否决（`ANSWER-VETO channel=esc`）；〔真凭据栏〕**仍空**＝任务文本来自 `WISP_TEST_TASK_TEXT` 注入位、模型那一次是 golden 脚本，不是真凭据真模型。⇒ ⛔ 本格勾的是"常驻进程会真举出一张卡"这件事有行为凭据，⛔ 不等于真凭据那一栏已闭。**）**
- [x] **AC#2 B 道的前置记账**：真控制台那一支要 owner 的凭据＋重出 exe ⇒ 在条件齐之前，⛔ 不派；本格只要求把"缺谁/缺什么"逐字登记（配方 §5 已写，复跑核）。**（三枚承重数编排者自己现跑复认，⛔ 未引腿的读数：① `ls -la "$APPDATA/wisp"`＝只有 `logs/`，**无 `config.toml`**、`secrets` 目录 `No such file or directory`；② `stat -c '%y' build/wisp.exe`＝`2026-10-07 11:57:46`，而 `git log -1 -- cmd/wisp`＝`aa6c1881 2026-10-08 19:20:48` ⇒ **exe 旧 31h23m**（本程未重建、机主未令还原）；③ `git grep -c "runConsoleLoop" HEAD -- '*_test.go'`＝**rc=1 零命中**＝"人坐终端敲 `task`"那一形今天没有自动台件。⇒ B 道今天**判不了**，缺的三件逐字登记在案；条件齐之前 ⛔ 不派。**）**
- [x] **AC#3 与票 246 的边界**：⛔ 本票**不动**票 246 任何框；结论回写进本票。**（腿与本程零触碰票 246：`git status --porcelain -- .scratch/wisp/issues` 在 277 交件时只含本票与票 283；票 246 `AC#7` 框保持原 `- [x]` 一字未动。**结论回写＝本节 AC#1 那两栏分写**——`A732` 裁"AC#7 只覆盖起管线那半"的前半（会不会真举卡）今天有了〔接缝注入栏〕凭据，〔真凭据栏〕仍欠；⛔ 不据此动 246 任何格，也不撤任何勾。**）**

## 禁区

⛔ 不许把"静态链全接"当行为凭据；⛔ 不许为让 A 道变绿放宽任何断言；⛔ 真窗类发帖按 `A718 §3` 的凭据强度限制（老 exe/旧 dist 不可作数）；⛔ 三枚冻结件一字不动；⛔ 零 push。

## 更正（原话逐字留在上面、就地打旧；2026-10-09 编排者自己现量）

- 本节「现量」那两枚行号**在 HEAD 上都漂了 1 行**：`fmt.Printf("wisp: 卡片挂起…` 的真身在 **`cmd/wisp/resident_approval_windows.go:886`**（票面写的 `:885` 现在是 `"channels", channelRosterText(...)` 那枚 attr 行）；审计行 `approval: 常驻进程显示一张确认卡片` 起在 **`:881`**（票面写 `:880`），整块 `:881-885`。尺＝`git show HEAD:cmd/wisp/resident_approval_windows.go | sed -n '878,890p'`（行号与短语同一次取数）。**引句字面未变**，只是行号。来路＝本票由编排者（我）立于 10-08 19:3x，行号是我当时抄 `A725` 的，漂在我这一步。
- 另一枚**由腿顶回、编排者复跑确认**的过期转述：`AGENTS.md`／台账里那句「`winlive` 在 CI 无 tag 档」已被 `351e5a5e`（2026-10-08，票 111 r6 `AC#11`）推翻 ⇒ 精确形状＝**CI 编译 winlive 档但不执行其用例**：`.github/workflows/ci.yml:606` 步名 `winlive compile gate (go vet -tags winlive, ticket 111 AC#11)`、`:655` `run: go vet -tags winlive ./cmd/wisp/ ./internal/ball/`（尺＝`git show HEAD:.github/workflows/ci.yml | grep -nE "winlive"`，编排者自己跑）。⚠ 与 `A721` 那条"撤销口令射程不足"同族，不冲突：那道步在**真实 run** 里有没有求值过，是另一格（`A721` 已具名它受"零求值 26 枚"那一族牵制）。
