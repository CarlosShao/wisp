# 写腿 `33r11-r2` 申报表 —— 33r11 census 两枚红（384-v1 归因的 381 族债）补名册

> 写腿＝`33r11-r2`（实现者；非实现者终裁归 `33r11-v1`）。基线＝dev `46a80b579`（开工三核：HEAD 在 dev＝派单基线／索引 0 staged／`cmd/wisp` 0 脏）。
> 现跑 2026-10-31 21:1x–21:4x +08。日志件＝仓外 `D:\tmp\wisp\20261031\33r11r2\logs\*`（只建不删）。

## 0. 一句话结论

**两枚 census 红→绿，反向三发全复红且定向，零回归（红名集合零新增），381-r2 本体零字节。**
修＝33r11 census 判据自己要求的 reviewed act：把 `waveOutSpeaker.releaseLocked`／`unprepareLocked` 两条登记进 `lockedNamingRoster33r11`（census 判据逻辑行、`voice_out_windows.go`（381-r2 本体）零字节触碰）。

## 1. 变更本体（写面恰一枚）

`git diff --stat`＝`cmd/wisp/panel_locked_naming_33r11_windows_test.go | 10 ++++++++++`（+10/−0，纯插入）：

1. 名册 `lockedNamingRoster33r11` ＋2 条目（key/value 形状逐字照既有条目 `PanelManager.setPriorFocusLocked` 的现形）：
   - `"waveOutSpeaker.releaseLocked": "requires the caller to hold s.mu; body takes no lock"`
   - `"waveOutSpeaker.unprepareLocked": "requires the caller to hold s.mu; body takes no lock"`
2. 每条目一段注释（照 `setPriorFocusLocked` 的注释段形状）：381-r2 的 winmm 分片等待语境（两个调用方 `Play`/`Drain` 各自在 `s.mu.Lock()/defer Unlock()` 段内调；`releaseLocked` 体内零锁调用、`abortReq` 是 `atomic.Bool` 非互斥锁；`unprepareLocked` 唯一调用方＝锁内的 `releaseLocked`）。

**零触碰清单（现跑为证）**：census 判据逻辑行（`lockedFuncsIn33r11`/`TestLockedSuffixNeverTakesTheLockItself33r11` 的两个 `t.Errorf`/名册反向检查/行为案/四枚正控）零字节；`cmd/wisp/voice_out_windows.go`（381-r2 本体，`releaseLocked`/`unprepareLocked` 实现）零字节；`internal/**`、thresholds/golden、testdata/frontend/design/scripts/.github/票面/台账/specs/PLAN 零字节。

**登记的判据依据**（为什么这是"判据要求的登记形状"、不是改判据）：
- census 判据文件头 (a)：`every function or method … whose name ends in "Locked" must (a) be named in the roster below … A new Locked-suffixed func is red until it is reviewed and named here, or renamed`——加名册条目就是判据写明的唯一两条出路之一。
- `384-v1-verdict.md` §6/§9.4：两红＝两方法带 `Locked` 后缀未入 33r11 名册；处置＝「381 族补名册或改名」。
- 派单点名修法：「按 233 §8 只加读者形状补名册登记（**或按 census 判据该登记的形状**）」——本腿按括号内兜底口径执行：census 判据该登记的形状＝`lockedNamingRoster33r11` 条目（key=`RecvType.method`、value=一句诚实描述）。「233 §8」在盘上无可核出处（票 233 文件无 §8、「读者形状」全仓 grep 仅派单一处）——未触发停手①，因为 census 判据自身的登记形状与名册现形完全对得上（下节现读为证）。

## 2. census 两枚红的现读（改前，`head-census-before.txt` rc=1，红句逐字）

census 头两行：`33-r11 census: 3 Locked-suffixed func(s) across 48 production source file(s) of package main`；两 hit 体内均 `0 lock operation(s)`（判据 (b) 干净、纯 (a) 缺登记）。两枚红句逐字：

```
waveOutSpeaker.releaseLocked (voice_out_windows.go line 338) carries the "Locked" suffix but no roster entry: a new name that claims the caller holds a mutex must be reviewed and named in lockedNamingRoster33r11, or renamed
waveOutSpeaker.unprepareLocked (voice_out_windows.go line 376) carries the "Locked" suffix but no roster entry: a new name that claims the caller holds a mutex must be reviewed and named in lockedNamingRoster33r11, or renamed
```

与 `384-v1` §6 双树红形逐字相同（行号 338/376 同）＝HEAD 既有 381 族债现读确认。

## 3. 正向判据（登记 ⇒ census 红→绿）

`census-after.txt` **rc=0**：census 族 5/5 PASS（`TestLockedSuffixNeverTakesTheLockItself33r11`／行为案／三枚正控全绿）。census 头两行变＝`3 Locked-suffixed func(s) across 48 production source file(s)`，三 hit 全部名册内、体内均 0 lock ops。
尺＝A1326 staging 清单全套：工作树 `go test -c` 出 `wisp.test.exe`（DLL 三枚贴 binary）＋WISP_* 五 env＋cwd＝主仓 `cmd/wisp` 包目录＋`-test.timeout 30m`。

## 4. 反向判据（摘登记 ⇒ census 复红；三发 overlay ＋还原核）

仓外 `-overlay` 三发（零树突变；每发恰摘名册行、`mut.sh` 在盘），每发只跑 census 主案：

| 发 | 摘 | rc | 红句（逐字） | 定向自证 |
|---|---|---|---|---|
| M1 | 摘 `releaseLocked` 条目 | **1** | `waveOutSpeaker.releaseLocked (voice_out_windows.go line 338) carries the "Locked" suffix but no roster entry: …`（恰 1 枚） | `unprepareLocked` 不红（仍在册）⇒ 红的只是被摘那条 |
| M2 | 摘 `unprepareLocked` 条目 | **1** | `waveOutSpeaker.unprepareLocked (voice_out_windows.go line 376) carries the "Locked" suffix but no roster entry: …`（恰 1 枚） | `releaseLocked` 不红 |
| M3 | 两条全摘（＝HEAD 原状形状） | **1** | 两枚红句逐字（§2 同文） | ＝改前 HEAD 既有债逐字复现 |

**还原核三读数**：①overlay 只换编译输入、主树零突变——`git status --porcelain -- cmd/wisp` 恰一枚 `M panel_locked_naming_33r11_windows_test.go`；②还原态复跑＝§3 `census-after.txt` rc=0（同一棵树、无 overlay）；③树内变异串 grep＝0（`git grep d22PositiveControl33r11r2 -- cmd internal` rc=1 无命中）。

## 5. 零回归

**六族合跑**（`^Test(15e|26R3|41a2r1|381r2|43|384)`，staged 全套尺）：**RUN=31 / PASS=30 / FAIL=1 / SKIP=0**，唯一红＝`Test384ResidentWatchdogEndToEnd (3.31s)`。

- **归因**：单跑 `^Test384` ×3 全绿（rc=0、各 3 PASS，`384-solo1..3.txt`）⇒ 合跑负载时序红＝**HEAD 既有形状**（390-r1 A1329 / 26p-r1 A1331 同族照录先例）；红名集合**零新增**。本腿写面仅 `_test.go` 名册条目，机制上接触不到 watchdog 内存时窗。
- ★**381r2 winmm 分片案 5/5 全绿**（派单点名必须仍绿）：`Test381r2RealDevicePlayingPhaseAbortWithin400ms` 真设备 **37.00s PASS**（≤400ms 判据）＋其余四枚 PASS。
- ★**`262/640/902` 逐位同**：`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms`（`resident_voice_15e_windows_test.go:428`）。
- ★**Test256 照录（⛔ 本腿碰）**：`^TestTicket256` 定向 rc=1、红句逐字 `the resident leg passes an unexpected Options field "Answer": the shipped field set is no longer the one ticket 256 adjudicated, so re-read what this gate is now being fed`（`TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims`）＝**375 族债**（384-v1 §6 归因 pickaxe＝`b63d56413`），与 381 族无关，归其所有者。

## 6. 门禁四把尺 ＋ 两尺（共享工作树，rc 现取）

| 尺 | 读数 |
|---|---|
| `go build ./...` | **rc=0** |
| `go vet ./cmd/wisp` | **rc=0** |
| `gofmt -l cmd/wisp` | **空**（rc=0） |
| `gofumpt -l cmd/wisp` | **空**（rc=0，`D:\work\base\gopath\bin\gofumpt.exe` 绝对路径） |

## 7. d22scan（导出树 ＋ 正控）

- **导出树**＝`git archive HEAD(46a80b579)` 至 `D:\tmp\wisp\20261031\33r11r2\export`，`d22scan.exe -root` 净扫 **rc=0 clean**（`d22scan.txt`）。
- **分母四数**＝**`bans#1-5 internal/=287, cmd/=52, ban#8 internal/=729, cmd/=180`**。⚠ 对照派单 286/51/726/177：+1/+1/+3/+3 漂移＝**26p-r1（`4c21ac42c`）的 Go 文件**——与本腿无关且本腿**零新增文件**（改现有 `_test.go` 内容不改枚数），修复后分母仍为该四数、零漂移。
- **正控两发可逆**：种 `zz_d22_positive_control_33r11r2.go`（字符串含 U+2264）⇒ **rc=1** 逐字点名 `cmd/wisp/zz_d22_positive_control_33r11r2.go:4: [emoji] ban #8 glyph in scope cmd/ is banned (D23): non-comment text, string literals included; comments are exempt per Q-46(c)`；删后复净 **rc=0**、分母逐字同（`d22scan-emoji.txt`/`d22scan-revert.txt`）。

## 8. 三条硬边界 ＋ 未定义即停核

1. ⛔ 改 381-r2 本体：**守住**——`voice_out_windows.go` 零字节（`git status` 为证）。
2. ⛔ 改 census 判据／Test256（375 族）／`internal/**` 语义：**守住**——判据文件只有名册条目＋注释段插入，判定逻辑零字节；Test256 只照录；`internal/**` 零字节。
3. ⛔ thresholds/golden/testdata/frontend/design/scripts/.github/票面/台账/specs/PLAN：**守住**——diff 恰一枚 `_test.go`。
- 停手①登记形状对不上：**未触发**（census 判据的登记形状与名册现形逐字对上）。
- 停手②登记后 census 仍红：**未触发**（census 即刻 5/5 绿）。
- 停手③必须动 381-r2 本体：**未触发**。

## 9. 纪律与环境

- commit＝显式 pathspec 两枚（本件＋`panel_locked_naming_33r11_windows_test.go`）、主题逐字 `roster 债修复腿 33r11-r2`、无 `# ` 前缀、`-F` 在 `--` 前、**⛔ push（归编排者）**、临时件全在 `D:\tmp\wisp\20261031\33r11r2\**`、无裸 `go func(`、无墙钟超时、`unset GOTMPDIR`、rc 取法＝先重定向再 `$?`。
- staging 尺＝A1326 补卡清单全套（wisp.test.exe＋DLL 三枚贴 binary＋WISP_* 五 env＋cwd＝包目录＋`-timeout 30m`）；跑批一律 Git Bash（pwsh 拆 `-test.v` 踩坑先例在册）。
- 本腿一发射向主仓工作树的字节只有申报件＋名册文件；census 判据、381-r2 本体、375 族、其余禁区零字节。
