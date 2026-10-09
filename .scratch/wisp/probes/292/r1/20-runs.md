# 票 292 腿 292-r1 AC#3 — 四发整包三数并排 + 每发 tasklist 读数 + 逐名作差

harness 逐字（写在 runner `.scratch/wisp/probes/292/r1/run-pass.sh`）：
`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1 -v`
整包原始日志：`logs/pre-1.txt` `logs/pre-2.txt` `logs/post-1.txt` `logs/post-2.txt`（各约 311 KB，未读进上下文，只用 grep 局部取数）。
三数尺（四发同一把）：`^=== RUN` / `^--- PASS` / `^--- FAIL` / `^--- SKIP`（顶层行，缩进的子测试不计数）。
环境红族尺：`grep -o 'exit status 0xc000013[5a]'` ⇒ **四发全空**，即没有任何一发是 DLL 缺失或 Ctrl+C 族，四发都真跑了用例。

| 发 | 起点（date stdout） | tasklist before | tasklist after | RUN | PASS | FAIL | SKIP | 日志字节 |
|---|---|---|---|---|---|---|---|---|
| pre-1 | 17:06:00 | wisp.exe=0 balldebug=0 | wisp.exe=0 balldebug=0 | 389 | 273 | 4 | 2 | 310840 |
| pre-2 | 17:14:00 | wisp.exe=0 balldebug=0 | wisp.exe=0 balldebug=0 | 389 | 272 | 5 | 2 | 311178 |
| post-1 | 17:31:16 | **wisp.exe=1** balldebug=0 | **wisp.exe=1** balldebug=0 | 389 | 272 | 5 | 2 | 311454 |
| post-2 | 17:40:30 | **wisp.exe=1** balldebug=0 | **wisp.exe=1** balldebug=0 | 389 | 270 | 7 | 2 | 312360 |

## 逐名红册（原文抄，含秒）

改前（既有红名册，两发各自全抄）：
- pre-1（4 枚）：`TestPanelHostRealWindowHopAndLifecycle (5.06s)` / `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (20.01s)` / `TestAC14AwaitedBindingReplyReachesThePage (20.06s)` / `TestAC14GoSideEvalPushReachesThePage (20.01s)`
- pre-2（5 枚）：以上 4 枚 + `TestAC4FocusReturnToPriorWindowGap33r5 (5.12s)`
- **改前交集 = 4 枚**；改前并集 = 5 枚（`TestAC4FocusReturnToPriorWindowGap33r5` 只在 pre-2 出现 ⇒ 单发尺不可靠，正是票面 `AC#3` 警告的 ±1 漂移，本腿实测复现）

改后：
- post-1（5 枚）：改前交集那 4 枚 + `TestAC246DevLegIgnoresTheTestTaskInjection (34.30s)`
- post-2（7 枚）：改前交集那 4 枚 + `TestAC246DevLegIgnoresTheTestTaskInjection (33.49s)` + `TestAC4FocusReturnToPriorWindowGap33r5 (5.23s)` + `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (2.38s)`
- **改后交集 = 5 枚** = 改前交集 4 枚 + `TestAC246DevLegIgnoresTheTestTaskInjection`
- 改后并集 − 改前并集 = {`TestAC246DevLegIgnoresTheTestTaskInjection`, `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`}（后者只出现在 post-2 一发，与 pre-2 的 `AC4FocusReturn` 同族＝漂动枚）

## 那 1 枚新增红的具名归因（不是我把树改红的）

两发改后都带 `wisp.exe=1`，两发改前都是 0 —— 而这枚红的自身判语就写着「a dev Wisp is already running」，
所以它与「会话里有一枚活的 wisp.exe」同时出现、同时消失，形状完全对上。那枚进程是谁起的（`who-owns-wisp.ps1` 现跑）：

```
ProcessId       : 9084
ParentProcessId : 32128
CreationDate    : 2026/10/9 17:27:04
ExecutablePath  : D:\work\workspace\projects plans\Wisp\build\wisp.exe
CommandLine     : "D:\work\workspace\projects plans\Wisp\build\wisp.exe"
父进程 32128 = bash.exe，其 CommandLine 内含：
  o=.scratch/wisp/probes/orch/2026-10-09-c1c2-resident.raw.md && ... ./build/wisp.exe >> "$o"
```

⇒ **它是编排者自己那一发常驻 Wisp 探针（17:27:04 起，写 probes/orch/ 的 raw 件），不是我的 `go test` 残留，也不是机主手起的。**
证据：我的两发改前跑完时 after=0（runner 自己数过），该进程 17:27:04 才出生，正落在我两发改前之间。

**因此我没有停它。** 派单那条「跑完若还有残留就按具名 `-Id` 停掉」的前提是「残留＝我这发 `go test` 拉起来的子进程」，这个前提在盘上不成立：
它由编排者的 bash 起、还在往 `probes/orch/2026-10-09-c1c2-resident.raw.md` 追加写，我 `-Id 9084` 停它＝掐断别人正在收的证据。
⛔ 我不按字面执行这条指令，并按派单自己那条「与转述冲突时以现场为准、具名报回」把它报回来。

## 这格的结论怎么写才不越权

- 主证（语义面零）：`+N==-N` 3/3、`git diff -U0` 非 `+++`/`---` 行 total=6 全为注释、noncomment=0、唯一 hunk `@@ -598,3 +598,3 @@`。改动不可能影响任何用例的判定。
- 辅证（名册作差）：按票面要求的**交集尺**是 5−4＝**新增 1 枚**（`TestAC246DevLegIgnoresTheTestTaskInjection`），
  该枚的在场/缺席与编排者那枚 `wisp.exe` 的在场/缺席**逐发一一对齐**；改前后共有的稳定红是同一批 4 枚。
- 所以本格**不自判成立**：判据原文写「新增红 0 枚」，交集尺上现在是 1 枚，虽然它有具名外因、且主证为零，
  **过不过由编排者裁**（见回报「需要裁的点」C1）。要把它坐成 0，最便宜的办法是：编排者那发常驻探针结束后我再补一发改后整包（约 8 分钟）取交集。
