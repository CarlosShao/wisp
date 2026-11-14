# 23 — System/window/input tool family: system.*, window.*, clipboard, media, screen, input.type

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 20-host-bridge-fs-tools
**Parallel slots:** ≤2 sub-agents (A: system.*/media/clipboard; B: window.*/screen.capture/
input.type + C7 Image wiring)
**Spec refs:** SPEC-07 §3 (D34 rows), C7 Image, D45 input binding, §16.11#2 UIPI, D34 note②

## What to build
The machine-control builtin family per the D34 authoritative table, including the C7 Image
content part that makes `screen.capture` actually reach the LLM, and `input.type` with its
per-target-process session-binding groundwork.

## Key constraints (RiskLevels fixed by D34 — verify each)
- `system.get` L0 (focused-window title taint-marked), `system.set` L1 (volume/brightness/mute/
  media keys/IME), `system.power` L2 (lock/sleep/restart/shutdown/logoff — shutdown+restart
  require typed confirmation word).
- `window.list` L0; `window.focus` L1; `window.manage` L1 (min/max/move/resize/snap) but
  close = L2 (unsaved-content risk).
- `clipboard.read` L0 (taint) / `clipboard.write` L1 (exfil channel per F4).
- `media.control` L1 (play/pause/next/prev).
- `screen.capture` L1 → C7 ImagePart{mime, bytes_ref, alt}; bytes_ref memory-only — never on
  disk/DB/logs/diagnostics unless user explicitly saves; DRM-window black-frame documented.
- `input.type` L2 first-use per target process → session-grant downgrade hook (binding table
  now; grant lifecycle at 49): confirm card carries "仅本次 / 本会话内允许向此应用 / 所有应用
  (danger-styled, expanded click)" options; UIPI: injecting into elevated windows fails with
  explicit "Wisp 无法向管理员窗口输入" (never attempt bypass); types text AND key sequences at
  the focused window.
- All system calls through the host bridge; risk decisions via 17 (R5-R8 rows exercised).
- Fixtures: a test target app (notepad/scripted window) for window/input acceptance on the
  Windows runner.

## Out of scope
- Grant persistence/lifecycle (49); voice dictation flow (32); panel screens (37/40).

## Acceptance criteria
- [x] RiskLevel matrix test: every action above asserted (esp. window.close=L2, power=L2,
      set=L1, get=L0+window-title taint).
- [x] screen.capture: screenshot reaches mock LLM as ImagePart (C7 wire format asserted);
      bytes_ref absent from disk/db/logs; DRM window → black frame + doc note.
- [ ] input.type: types 300-char text into fixture app focused window; elevated-process
      injection → explicit UIPI error; confirm-card three options render with danger styling
      on "all apps".
- [x] clipboard round-trip + taint propagation into R4 (write of tainted content → L2).
- [x] system.power shutdown path requires typed confirmation word (unit test with fake executor).

## Progress log (append-only, newest last)

- 2026-11-03 21:4x 编排者：后端+工具接线层落地（119c3fa91 后端6文件；a2f6f4a95 system.go BuiltinSystemEntries 11工具+run.go注册+system_23测试 1450行；df67e002e 修229普查登记52的OS进程spawn站点）。internal/tools 整包 ok、d22 clean。**已满足**: AC1 风险矩阵、AC5 关机确认词、AC4 taint→R4。**残余**: AC2 C7 ImagePart→LLM 接缝(需 internal/agent loop+llm.ImagePart 集成)+bytes_ref内存only/DRM黑帧；AC3 input.type 真fixture app+UIPI 真机+确认卡三选项UI(前端)；AC4 真剪贴板round-trip。**卡点**: AC3 确认卡UI=frontend(机主外部agent)+真机集成，非纯后端可闭合。票仍 open。
- 2026-11-09 14:3x 写腿（票23复核腿）：**上格进度已过期，实况如下（本腿零写面，只追加）**。AC2 的 C7 ImagePart→LLM 接缝**已闭合**：a67e1d335（loop+llm 半）+ 111bfdb04（tools 半：`Result.Images`+`bridgeImages`+system.go 文案改写）+ e3abe0ebb（遗产腿代提 5 测试腿+AC4 真剪贴板往返门闩+SrcClipboardRead→R4+AC3 UIPI 谓词/D45 绑定键）。**非实现者裁决表 `docs/evidence/s1/23-verdict.md` 终裁：AC#1/#2/#4/#5 PASS、AC#3 PARTIAL**——残余 R1 真fixture app 300字符注入／R2 真提权窗口 UIPI 实机／R3 确认卡三选项 UI（`frontend/**`，机主外部 agent 写面）／R4 真屏截图→真 LLM 演示（owner 真机签收类）。**名册差集（D34 PLAN.md:2560-2567 ＋ SPEC-07 §3:56-67 ↔ `BuiltinSystemEntries` 11 枚，run.go:773 注册）＝空集**（多 0 缺 0），且 `Test23RosterCoversElevenD34Rows`＋`Test23RosterDeclarationsMatchTheFrozenTable` 双钉。**门禁在共享克隆（`git clone --shared`，锚 cc5b297e2）fresh**：gofumpt 空 rc0 / vet rc0 / `go test ./internal/tools -count=2` ok 64.2s / d22scan clean；点名 `Test23` 族 10 枚 = 9 PASS + 1 SKIP（`Test23LiveClipboardRoundTripThroughRealOS` 门闩默认不碰操作者剪贴板，形状正确）。勾不翻（裁决表 §8 明文把翻勾留给编排者/票面流程）；票保持 open，残余归属见裁决表 §8。

## 编排者翻框（2026-11-09 15:1x，非实现者判语）

依据 `docs/evidence/s1/23-verdict.md` §7 逐 AC 终裁表（非实现者裁决）＋该腿 fresh `go test ./internal/tools -count=2` ok 64.2s：
- **AC#1 风险矩阵＝PASS** ⇒ 翻（count=2 绿，无残余）
- **AC#2 C7 ImagePart＝PASS** ⇒ 翻（bytes_ref 内存 only 由类型形状＋两枚 DB/日志腿＋`assertNoPNGOnDisk23` 钉住；无本域残余）
- **AC#4 剪贴板＝PASS** ⇒ 翻（taint→R4 fresh 绿，live 腿门闩正确）
- **AC#5 关机确认词＝PASS** ⇒ 翻（fresh 绿）
- **AC#3 input.type＝PARTIAL，⛔ 不翻**：后端余量已闭（L2 声明/UIPI 拒绝序/绑定键），三枚残余 R1（真 fixture app 注入）/R2（真提权窗实机）/R3（确认卡三选项 UI 在 `frontend/`）全在真机/前端写面。
- ⇒ **票 23 保持 open 不 `-done`**；残余 R1–R3 归属 frontend/owner-open-window，等真机窗口期一次收口。

## 写腿真机读数（2026-11-14，写腿；只追加，未翻勾）

- 2026-11-14 10:3x 写腿（23 真机读数腿）：**AC#3 那格上了真机，结论是「第一分句成立、第二分句产品侧不成立」**。仪器＝`.scratch/wisp/probes/23/`（`fixture23` 脚本窗 fixture 进程＋`probe23` 读数腿＋`diag23`/`uipi23` 只读谓词腿；注入走真名册 `tools.BuiltinSystemEntries(tools.SystemDeps{})`，与 `cmd/wisp/run.go:897` 同形，零产码改写）。**300 字那半量到并逐字对上**：runA ASCII 300 字 → 窗内 `chars=300`、`wm_char=300`、sha256 `1f47881e…fbf` 与期望**相同**，产品回「已向 wisp23fixture.exe 注入 300 个字符」（1632 ms／600 事件）；runB CJK 300 字同形（sha `84056f1d…e0c`，1598 ms）。前置焦点两条独立读数：probe 侧 `GetForegroundWindow=0x82e1476`＋`GetGUIThreadInfo.hwndFocus=0x6f10fba`，产品侧 `system.get` 报 `focused_title="WISP23-FIXTURE input.type target"`／`focused_process="wisp23fixture.exe"`。**提权那一半失败**：对真 High 窗（Taskmgr pid 6904，IL `12288`／`S-1-16-12288`）调 `input.type` → `IsError=false`＋「已向 taskmgr.exe 注入 15 个字符」，**没回冻结句 `ErrElevatedTarget`**；根因只读定位到 `internal/tools/system_backend_windows.go:741`「`integrityOf(windows.Token(h))` 把**进程句柄**当**令牌句柄**」→ `GetTokenInformation` 报 `The handle is invalid.`→ 返 0 → 谓词恒假（`diag23` 并排读数：shipped=0／corrected=12288／self=8192；`uipi23` 复读同形）。门禁之所以绿＝拒绝腿把 `processElevated` 打桩成 true（`ticket23_input_uipi_windows_test.go:66-72`），真谓词那条腿只测"自己窗口放行"（恒假的正确答案就是 false，`:129`）。**第三分句（卡片三选项/danger 样式）本腿不可判**：生产 `frontend/**` 三档文案 0 命中，仅设计稿 `design/old/screens/approval.html:848/:856/:871`＋`SPEC-06:117/:120` ⇒ 仍是 R3 归 frontend。**未翻 AC 框、未改 `internal/**`、未提权未弹 UAC**；证据件 `docs/evidence/s1/23-input-type-real-window.md`（含重跑命令、run2–run8 为何按安全闸退出、以及"再跑一轮需机主把前台让回普通窗口"这条人工前置）。**新欠账（具名归写腿修）**：`system_backend_windows.go:741` 改成先 `OpenProcessToken(h, TOKEN_QUERY, &tok)` 再比 token IL——修完 `probe23` 原样复跑即可复核，无需 UAC。