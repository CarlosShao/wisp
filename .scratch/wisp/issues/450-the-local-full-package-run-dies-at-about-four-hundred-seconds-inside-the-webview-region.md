# 票 450 — 本机整包跑到 **约 400 秒**就在真 WebView2 窗那一区**整进程死亡**，两发两次都不是断言红

**状态**：待派（能否复现、是不是我自己的驻留球造成的，都要读数）

## 为什么值得单独一张票
本编队每个写腿都要跑一次 `go test ./cmd/wisp` 整包来证明"没把别人弄红"。这条路**今天在本机走不完** —— 它在 `resident_*`／`ticket*`／`voice_*` 那一区**还没开始跑**之前就整进程死了。后果不是"少几枚绿"，是**那一片测试在本机等于没人跑过**，而大家会拿"两发都 FAIL 0"当覆盖证明。这正是"看起来绿、实际没跑"的形状。

## 现量（2026-11-14 22:xx，票 228 切片 1 的实现腿，两发）
| 发 | 读到 | 死在哪 | 用时 |
|---|---|---|---|
| 1 | `PASS 339  FAIL 0` | `Test145ArgsJSON…` 半途 | 416 s |
| 2 | `PASS 334  FAIL 0` | `TestAC4PriorFocusSurvivesARefusedPanelSample` 半途 | 403 s |

两发**都不是断言失败**，是**进程消失**（无 `--- FAIL` 行）。`228-r3b` 记过同一族。

## 关键背景（可能是真因，⭐别先入为主）
两次死亡发生时，**编排者自己的驻留球在跑**（PID 49348，2026-11-14 16:23 起）。已知的**本机专属红**里就有 `TestAC246DevLegIgnoresTheTestTaskInjection`（"a dev Wisp is already running in this session"）——那说明**驻留球会与整包里的用例争用同一份会话/窗口资源**。⇒ 第一嫌疑是"球在跑"，⛔ 一上来就判"WebView2 不稳"。

## AC

- [ ] **AC#1 先把"死"变成可判事实**：复现一次，抓**退出码／最后一行输出／Windows 事件日志里的 application error 条目（进程名＋模块）**。⛔ 用"卡住了"描述。
- [ ] **AC#2 二分：球在跑 vs 球不在跑**：同一命令各跑 1 发 —— (a) 球在跑；(b) 球停（**由机主停，⛔ 腿自己 kill**；若机主不在，用"从未起过球的干净用户会话"替代并写明）。给两发的 `RUN/PASS/FAIL/SKIP` 与是否死亡。
- [ ] **AC#3 二分：WebView2 窗那一区 vs 其余**：把死亡点附近那几枚用 `-test.run` 单独跑（含 `Test145ArgsJSON*`、`TestAC4PriorFocus*`），看在**没有整包负载**时是否也死。
- [ ] **AC#4 判归属并落成规矩**：若 (a) 死 (b) 不死 ⇒ 真因＝驻留球争用，**把"跑整包前必须停球"写进本仓的测试前置说明**（`scripts/wisp-cli-staged-tests.sh` 的注释或 `docs/evidence/` 一份具名说明，⛔ 改 `docs/specs/**`）。若不是 ⇒ 交回可判事实与最小复现，**本票不许自己猜着改产码**。
- [ ] **AC#5 门禁**：本票若只改脚本注释/证据件，跑 `d22scan` 与 path-length 即够；若动产码，按整包 + 本机两枚已知红具名口径执行。

## 边界
⛔ `thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／D43 转移表；⛔ 放宽或跳过任何既有断言（若真需要跳过，必须**具名＋给数字**并写进证据）；⛔ 自己 kill 机主的球；⛔ `frontend/**`／`design/**`；⛔ add -A/--amend/reset/rebase/stash/checkout .；⛔push。
写面：`scripts/wisp-cli-staged-tests.sh`（**仅注释**）、`docs/evidence/s1/450-*`、`.scratch/wisp/probes/450/**`、本票 Progress log。

## 撤销口令
「450 撤」＝删本票与其证据件（不 reset/amend）。
