# 票 130 · `AC#1`「会掉的记录」清单 —— `130-r2` 刷新（HEAD `dbe861a02`）· 全 129 枚逐枚

**这是什么**：`docs/evidence/s1/130-ac1-inventory-and-ruling.md` §2 那份清单的**刷新版**。那份锚在 `ac6f31c`、分母写 53 行 / 22 包；本表锚在 `dbe861a02`，分母**现跑**。
**为什么重写**：`130-v1` §2 判 `AC#1` 差四件 ——（a）分母翻倍、（b）"其余 24 枚"只列了 19、（c）16 处行号漂移、（d）漏一个 B 窗口站点 `internal/winsec/winsec_windows.go:143`。本表把四件逐条补齐，且**一条行号都不抄**（每条自带 `grep -nF` 的逐字锚）。
**性质**：读数表，不是判据。零生产码。

---

## §0 三把尺（先声明它量的是哪个对象）

```bash
# 尺① 主式子＝本表的分母（包级 slog 记录，非测试）
grep -rn "slog\.\(Error\|Warn\|Info\)(" --include="*.go" internal/ cmd/ | grep -v _test     # 129 行
# 尺② 闭包＝本表的射程
go list -deps ./cmd/wisp/ | grep -c 'CarlosShao'                                            # 27（含 cmd/wisp 自己）
# 尺③ 闭包外的仓内包
comm -23 <(ls -d internal/*/ | sed 's#/$##') <(go list -deps ./cmd/wisp/ | grep -o 'wisp/internal/[a-z0-9]*$' | sed 's#wisp/##' | sort -u)
```

**尺①量不到的三形**（本节明写，别把本表当全集）：
1. `slog.Debug(` —— 今天 **1** 枚（`internal/winsec/resolve.go:176`，本表行 56 之外单列 A6）；Debug 被原生默认的 Info 门挡 ⇒ 今天**连 stderr 都不到**。
2. `slog.{Info,Warn,Error,Debug}Context(` —— 今天 **0** 枚（式子要含它，但今天空）。
3. **捕获式**：构造期 `slog.Default()` 存进字段（`internal/memory/open.go:118`、`internal/agent/loop.go:215`），之后经 `s.logger.X` 出声 —— 今天 **25** 枚站点（memory 16 / agent 9，`grep -rnE "\b[a-z]+\.log(ger|\(\))\.(Info|Warn|Error|Debug)\(" --include="*.go" internal/ cmd/ | grep -v _test`）。它们**不在 129 里**，缺口形状见 `130-ac1-inventory-and-ruling.md` §2 D 组。

**闭包今天的两个事实**（判语 §2(a) 复核成立）：`internal/audio`（7 枚站点）与 `internal/ball`（23 枚站点）**都在 `cmd/wisp` 闭包内**（resident 腿已接球与采集）；闭包外只剩 `internal/speech`、`internal/watchdog` 两包，**主式子 0 行** ⇒ 129 行**全部在闭包内**。

---

## §1 一条命令重验（本表的全部分母 + 每条锚）

```bash
cd '/d/work/workspace/projects plans/Wisp' && git rev-parse --short HEAD && \
grep -rn "slog\.\(Error\|Warn\|Info\)(" --include="*.go" internal/ cmd/ | grep -v _test | cut -d: -f1,2 | sort > /tmp/130r2-now.txt && \
grep -E '^\| [0-9]+ \| `[^`]+\.go:[0-9]+` \|' docs/evidence/s1/130-r2-impl-part2.md | grep -oE '`[^`]+\.go:[0-9]+`' | tr -d '`' | sort > /tmp/130r2-want.txt && \
printf 'census now=%s want=%s\n' "$(wc -l < /tmp/130r2-now.txt)" "$(wc -l < /tmp/130r2-want.txt)" && \
diff /tmp/130r2-now.txt /tmp/130r2-want.txt && echo CENSUS-OK && \
printf 'closure=%s audio=%s ball=%s\n' "$(go list -deps ./cmd/wisp/ | grep -c CarlosShao)" "$(go list -deps ./cmd/wisp/ | grep -c internal/audio)" "$(go list -deps ./cmd/wisp/ | grep -c internal/ball)"
```

期望逐字：`census now=129 want=129` → `CENSUS-OK` → `closure=27 audio=1 ball=1`。**任何一行漂了，diff 会点名 `file:line`**。
逐条锚的复跑形状：`grep -nF '<本行第五列的片段>' <本行第二列的文件>`；片段在同一文件里出现多次时（如 `observe: log retention sweep failed` ×2、`hotkey disabled (unset in config)` ×2），
该片段命中的每一行都在本表里有自己的一行 ⇒ **不靠行号认站点**。

---

## §2 组与恒等式（这就是"总数对得上"那一步）

| 组 | 枚数 | 是什么 | 行号（本表内） |
|---|---|---|---|
| **A** | **5** | `init()` 期：`internal/risk/winsec_c26.go:20` 的 `SetPathResolver` 一路 | 54–58 |
| **B** | **4** | 换 handler 之前的窗口（`InitLog` 已跑、`SetDefault` 未跑）：`observe` 的 sweep 与名册、`goroutine` 的 panic 兜底、**pipeline 自己封 `logs\`** 那条 | 39,40,41,59 |
| **C'** | **1** | 新增：**second launch 腿上根本没有听众**（不是"早"，是"never"） | 52 |
| **F** | **1** | 能力做完、0 生产调用者 | 53 |
| **其余** | **118** | 闭包内、按静态顺序在装听众之后（或只在无听众腿 / 另一枚 main 里可达） | 本表去掉上面四组之后剩下的 118 行 |
| **合计** | **129** | `129 = 5 + 4 + 1 + 1 + 118` ✅ | 与尺①逐字相等（§1 的 diff 就是这一步） |

**A6 不在本表**：`internal/winsec/resolve.go:176` 是 `slog.Debug` 形（尺①不含 Debug），今天**连 stderr 都不到**（原生默认 = Info 门）⇒ 判 `产不可达`（R9）。它让 A 组从 5 变成 6，但**不进 129**。
**件里那句"其余 24 枚"错在哪**：它想说的是 `28 = 6 + 2 + 1 + 19`（判语点名这一式），但 19 只是它逐条列出来的那 19，`internal/config 其余` 那半没有名字；
而今天分母是 129、组也变了（B 从 2 涨到 4、C' 是新的一枚）⇒ **19 与 24 两个数今天都不成立**，正确读数是 **118**，且 118 逐枚有名字（本表）。
**判语 §2(c) 的 16 处行号**：本表逐条重跑后的实际落点见 §3 的对照表（判语给的号**大部分对**，四处要修：B1 的 `379` 不是 B 窗口而是 `flushLoop` 的 tick；agent 构造点 `215` 对；`logsink` 三处 `149/153/184` 对；四腿装点 `109/222/226/339` 对）。

---

## §3 判与据的代码

**判**（"能不能造"三档，⛔ 把"造不出"写成"不会掉"）：

| 判 | 含义 |
|---|---|
| `可造` | 本机能真的把它造出来（给"据"里的造法/读数） |
| `造不出` | 本机条件不满足（给"据"里的具体条件），**不等于不会掉** |
| `产不可达` | 就"早于听众"这件事而言生产不可达（给"据"里的理由：装点行号 / 0 调用者） |
| `非缺口` | 它就是 install 那条记录本身 |

**据**（每条的理由代码；同一码同一条理由，⛔ 逐行重写）：

| 码 | 条数 | 理由 / 造法（可复跑） |
|---|---|---|
| `R1` | 95 | `产不可达`：resident 腿 `cmd/wisp/resident_windows.go:109` 装听众**之后**（事件循环 / 面板 STA / 球 / 热键 / 采集 / 关闭序列）。`cmd/balldebug` 是另一枚 main（只把默认 logger 装成 stderr TextHandler，无 JSONL）⇒ 属票 131 的枚举域。 |
| `R2` | 8 | `造不出`：`internal/config/{manager,writeguard}.go` 的两条可达路径 = `run` 腿 `:222` 之后（落盘）与 `wisp panel-inbound` 腿（**无听众**，`cmd/wisp/panel_inbound.go:50` 自述"installs no persistent log sink"）⇒ 该腿上是**真掉点**；本机没有该腿的 settings/mode 写请求夹具。 |
| `R3` | 3 | `产不可达`：`plugin.DisposalScope.Dispose` 非测试调用者 **0**（`grep -rn '\.Dispose()' --include='*.go' cmd/ internal/ \| grep -v _test` → 空）⇒ 同 F1 一族。 |
| `R4` | 2 | `产不可达`：`internal/proc/boot_windows.go:188` 的 `handleActivation` 只在 `RunEventLoop` 里跑（`resident:109` / `slo:281` 之后）。 |
| `R5` | 4 | `产不可达`：D38(e) 关闭序列，装点之后。 |
| `R6` | 1 | `产不可达`：`flushLoop` 的 `flushInterval` / `retentionSweepPeriod` tick，装点之后（**判语 §2(c) 把 `379` 算进 B1 是错的**）。 |
| `R7` | 1 | `产不可达`：wisp 侧唯一入口 `internal/models/bridge.go:40/46/59/64` ← `cmd/wisp/models.go:358`（`models ensure :339` 之后）。 |
| `R8` | 1 | `非缺口`：install 记录本身（`cmd/wisp/logsink.go:153` `SetDefault` 之后）。 |
| `R9` | 3 | `产不可达`：需 `SetPathResolver(nil)` 或第二次、不同实例的安装；全仓非测试调用者 **1**（`internal/risk/winsec_c26.go:21`）。 |
| `R10` | 1 | `造不出`：本机 conformance probe 全过（`probes_passed=2`，本笔现跑见 §3 A2 行）⇒ 要一台让 probe 失败的主机。 |
| `R11` | 1 | `造不出`：boot sweep 的 `os.ReadDir(<data>\logs)` 要在 `MkdirAll` 刚成功后失败（只读权限 / DACL 形状未造出）。 |
| `R12` | 1 | `造不出`：需窗口内已 spawn 的 goroutine（`logging.go:123` 起的 `log-flusher`）panic。 |
| `R13` | 1 | `可造`：变异——把 `log-flusher` 从 D38 名册摘掉，`goroutine.go:271` 立刻出声（判语 §3 M-2 用的是同形）。 |
| `R14` | 1 | `可造`：A1。**本笔现跑**（发声半）：`go test ./internal/risk/ -run 'TestNothingMatches130Audit' -v` ⇒ stderr 首行 `… msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2`。落盘半由 `130-v1` §4 实测（`slo -settle` post：记录 0 = A1，盘上 2 条 / 622 字节）。 |
| `R14w` | 1 | `可造`：`internal/winsec/winsec_windows.go:143` 两向实测（`130-v1` §4）——`slo -settle` post **记录 1**（盘上，`path=…\logs`，时间戳早于任何 listener 记录）；`providers` 腿 = stderr 且 `find` 零文件。 |
| `R15` | 1 | `产不可达`：`secret.MigratePlaintext` 非测试调用者 **0**。 |
| `R16` | 1 | `造不出`：需 `AllowSetForegroundWindow(ASFW_ANY)` 被拒（前台锁）。**且这条腿没有听众**：second launch 在 `cmd/wisp/resident_windows.go:88` 就 `return`，**永不**走到 `:109` ⇒ 这条 WARN 只到 stderr。 **这是本笔新补的 C' 组。** |
| `R17` | 3 | `产不可达`：`secret` 腿 `cmd/wisp/secret.go:226` 之后。 |

---

## §4 全 129 枚

| # | 站点 | 组 | 级 | 锚（`grep -nF` 的逐字片段） | 判 | 据 |
|---|---|---|---|---|---|---|
| 1 | `internal/audio/audio.go:136` | 其余 | Warn | `audio frames dropped: bounded channel full (slow consumer)` | 产不可达 | R1 |
| 2 | `internal/audio/audio.go:152` | 其余 | Error | `audio source failed` | 产不可达 | R1 |
| 3 | `internal/audio/captureopt.go:107` | 其余 | Error | `audio level: seam frame rejected, no level emitted` | 产不可达 | R1 |
| 4 | `internal/audio/mmdevice_windows.go:60` | 其余 | Error | `audio watcher: enumerator unavailable (hotplug disabled)` | 产不可达 | R1 |
| 5 | `internal/audio/mmdevice_windows.go:65` | 其余 | Error | `audio watcher: notification registration failed (hotplug disabled)` | 产不可达 | R1 |
| 6 | `internal/audio/wasapimic_windows.go:122` | 其余 | Error | `audio capture thread did not exit within 2s (abandoned wait)` | 产不可达 | R1 |
| 7 | `internal/audio/wasapimic_windows.go:295` | 其余 | Info | `audio capture reopened after device change` | 产不可达 | R1 |
| 8 | `internal/ball/ball_windows.go:398` | 其余 | Warn | `ball: slow drawFrame` | 产不可达 | R1 |
| 9 | `internal/ball/ball_windows.go:422` | 其余 | Warn | `ball: slow ULW frame` | 产不可达 | R1 |
| 10 | `internal/ball/ball_windows.go:425` | 其余 | Warn | `ball: UpdateLayeredWindow failed` | 产不可达 | R1 |
| 11 | `internal/ball/ball_windows.go:531` | 其余 | Warn | `ball: UpdateLayeredWindow failed` | 产不可达 | R1 |
| 12 | `internal/ball/ball_windows.go:704` | 其余 | Error | `ball: DPI renderer rebuild failed` | 产不可达 | R1 |
| 13 | `internal/ball/ball_windows.go:768` | 其余 | Warn | `ball: position persist failed` | 产不可达 | R1 |
| 14 | `internal/ball/ball_windows.go:888` | 其余 | Error | `cancel borrow fell back to the default binding: [hotkey] cancel is not a valid key combination` | 产不可达 | R1 |
| 15 | `internal/ball/ball_windows.go:990` | 其余 | Error | `ball: DestroyWindow failed; window handles stay alive in this process` | 产不可达 | R1 |
| 16 | `internal/ball/hotkey_reload.go:84` | 其余 | Warn | `ball: hotkey bridge config refresh failed; keeping current hotkeys` | 产不可达 | R1 |
| 17 | `internal/ball/hotkey_reload.go:99` | 其余 | Error | `ball: hotkey reload` | 产不可达 | R1 |
| 18 | `internal/ball/hotkey_reload.go:101` | 其余 | Info | `ball: hotkeys rebound after config change` | 产不可达 | R1 |
| 19 | `internal/ball/hotkey_reload.go:147` | 其余 | Error | `ball: hotkey reload panicked; hotkeys unchanged` | 产不可达 | R1 |
| 20 | `internal/ball/hotkey_windows.go:456` | 其余 | Info | `hotkey disabled (unset in config)` | 产不可达 | R1 |
| 21 | `internal/ball/hotkey_windows.go:461` | 其余 | Error | `hotkey not attempted (bad binding)` | 产不可达 | R1 |
| 22 | `internal/ball/hotkey_windows.go:471` | 其余 | Info | `cancel hotkey left unbound while idle: the configured key is ` | 产不可达 | R1 |
| 23 | `internal/ball/hotkey_windows.go:496` | 其余 | Info | `hotkey disabled (unset in config)` | 产不可达 | R1 |
| 24 | `internal/ball/hotkey_windows.go:502` | 其余 | Error | `hotkey not attempted (bad binding)` | 产不可达 | R1 |
| 25 | `internal/ball/hotkey_windows.go:508` | 其余 | Error | `hotkey occupied by another program, not registered` | 产不可达 | R1 |
| 26 | `internal/ball/hotkey_windows.go:512` | 其余 | Error | `hotkey registration failed` | 产不可达 | R1 |
| 27 | `internal/ball/hotkey_windows.go:612` | 其余 | Error | `cancel takeover refused by Win32; Confirming will have no cancel key` | 产不可达 | R1 |
| 28 | `internal/ball/renderer_windows.go:407` | 其余 | Error | `ball: EndDraw failed` | 产不可达 | R1 |
| 29 | `internal/ball/sta_windows.go:155` | 其余 | Error | `ball: the ui-sta owner could not destroy its own window before releasing the thread; the handle stays alive in this process` | 产不可达 | R1 |
| 30 | `internal/ball/sta_windows.go:197` | 其余 | Warn | `ball: the ui-sta release pump hit its cap before the thread's queue was empty; the thread goes back to the Go pool non-empty` | 产不可达 | R1 |
| 31 | `internal/config/manager.go:184` | 其余 | Warn | `config: locked loosening approved via L2 re-confirmation (D36 rule 1)` | 造不出 | R2 |
| 32 | `internal/config/manager.go:187` | 其余 | Warn | `config: locked loosening rejected; keeping previous values` | 造不出 | R2 |
| 33 | `internal/config/manager.go:335` | 其余 | Info | `config: locked section tightened, hot-applied` | 造不出 | R2 |
| 34 | `internal/config/manager.go:341` | 其余 | Info | `config: locked section updated (neutral), hot-applied` | 造不出 | R2 |
| 35 | `internal/config/writeguard.go:152` | 其余 | Info | `config: config.toml already carries this key; this write added nothing (the ` | 造不出 | R2 |
| 36 | `internal/config/writeguard.go:172` | 其余 | Info | `config: wrote merged change into config.toml (ticket 226: only the keys listed ` | 造不出 | R2 |
| 37 | `internal/config/writeguard.go:184` | 其余 | Warn | `config: wrote merged change into config.toml and kept hand edits this process ` | 造不出 | R2 |
| 38 | `internal/config/writeguard.go:211` | 其余 | Error | `config: diffKeyPaths called with a nil side; the caller must handle an absent file` | 造不出 | R2 |
| 39 | `internal/observe/goroutine.go:168` | B | Error | `goroutine panic recovered` | 造不出 | R12 |
| 40 | `internal/observe/goroutine.go:271` | B | Warn | `goroutine outside the D38 roster (leak symptom)` | 可造 | R13 |
| 41 | `internal/observe/logging.go:257` | B | Warn | `observe: log retention sweep failed` | 造不出 | R11 |
| 42 | `internal/observe/logging.go:379` | 其余 | Warn | `observe: log retention sweep failed` | 产不可达 | R6 |
| 43 | `internal/plugin/disposal.go:208` | 其余 | Error | `disposal_incomplete: Defer called after Dispose` | 产不可达 | R3 |
| 44 | `internal/plugin/disposal.go:332` | 其余 | Error | `disposal step failed` | 产不可达 | R3 |
| 45 | `internal/plugin/disposal.go:334` | 其余 | Info | `scope disposed` | 产不可达 | R3 |
| 46 | `internal/proc/boot_windows.go:251` | 其余 | Info | `activation requested by second launch` | 产不可达 | R4 |
| 47 | `internal/proc/boot_windows.go:257` | 其余 | Info | `activation requested by second launch` | 产不可达 | R4 |
| 48 | `internal/proc/shutdown.go:130` | 其余 | Info | `shutdown step skipped (module not present)` | 产不可达 | R5 |
| 49 | `internal/proc/shutdown.go:145` | 其余 | Error | `shutdown step abandoned wait (deadline exceeded)` | 产不可达 | R5 |
| 50 | `internal/proc/shutdown.go:147` | 其余 | Error | `shutdown step failed` | 产不可达 | R5 |
| 51 | `internal/proc/shutdown.go:172` | 其余 | Warn | `fast shutdown: skipping non-critical log/db flushes (D38e fast path)` | 产不可达 | R5 |
| 52 | `internal/proc/singleinstance_windows.go:171` | C' | Warn | `proc: AllowSetForegroundWindow(ASFW_ANY) was refused; the resident's SetForegroundWindow may be denied by the foreground lock` | 造不出 | R16 |
| 53 | `internal/secret/migrate.go:198` | F | Warn | `plaintext api_key migrated to DPAPI-protected references (D33)` | 产不可达 | R15 |
| 54 | `internal/statemachine/machine.go:139` | 其余 | Error | `state machine rejected transition` | 产不可达 | R7 |
| 55 | `internal/winsec/resolve.go:147` | A | Info | `winsec: sealing path resolver left at the built-in floor` | 产不可达 | R9 |
| 56 | `internal/winsec/resolve.go:151` | A | Error | `winsec: refusing to release the sealing path resolver` | 产不可达 | R9 |
| 57 | `internal/winsec/resolve.go:158` | A | Error | `winsec: refusing to install a path resolver into the sealing seam` | 造不出 | R10 |
| 58 | `internal/winsec/resolve.go:166` | A | Error | `winsec: refusing to replace the already installed sealing path resolver` | 产不可达 | R9 |
| 59 | `internal/winsec/resolve.go:182` | A | Info | `winsec: sealing path resolver installed` | 可造 | R14 |
| 60 | `internal/winsec/winsec_windows.go:143` | B+C | Warn | `winsec: seal cleared principals that stood on this object` | 可造 | R14w |
| 61 | `cmd/wisp/logsink.go:184` | 其余 | Info | `wisp: persistent log sink installed` | 非缺口 | R8 |
| 62 | `cmd/wisp/panel_control_windows.go:717` | 其余 | Warn | `panel host: inbound message is not an RPC frame - received and dropped, no door and no reply (ticket 253 AC#3)` | 产不可达 | R1 |
| 63 | `cmd/wisp/panel_control_windows.go:734` | 其余 | Warn | `panel host: inbound RPC frame names no bound door - received and refused instead of resolved null (ticket 253 AC#3)` | 产不可达 | R1 |
| 64 | `cmd/wisp/panel_host_windows.go:1070` | 其余 | Warn | `panel host: pre-create quit-drain hit its cap before the queue was clean` | 产不可达 | R1 |
| 65 | `cmd/wisp/panel_resident_windows.go:204` | 其余 | Warn | `panel host: [panel] geometry source unreadable, sizing at the host's own default` | 产不可达 | R1 |
| 66 | `cmd/wisp/panel_resident_windows.go:252` | 其余 | Info | `"panel host: " + line` | 产不可达 | R1 |
| 67 | `cmd/wisp/panel_resident_windows.go:264` | 其余 | Warn | `panel host: embedded bundle unavailable, serving nothing` | 产不可达 | R1 |
| 68 | `cmd/wisp/panel_resident_windows.go:286` | 其余 | Error | `panel thread: refusing to bring up the panel window` | 产不可达 | R1 |
| 69 | `cmd/wisp/panel_resident_windows.go:436` | 其余 | Error | `panel host: show failed on the panel thread` | 产不可达 | R1 |
| 70 | `cmd/wisp/panel_resident_windows.go:440` | 其余 | Error | `panel thread: retiring without a panel window after a named refusal` | 产不可达 | R1 |
| 71 | `cmd/wisp/panel_resident_windows.go:477` | 其余 | Info | `panel host: raise request landed on the panel thread` | 产不可达 | R1 |
| 72 | `cmd/wisp/panel_resident_windows.go:558` | 其余 | Info | `panel thread exited cleanly` | 产不可达 | R1 |
| 73 | `cmd/wisp/panel_resident_windows.go:560` | 其余 | Warn | `panel thread did NOT exit within its budget` | 产不可达 | R1 |
| 74 | `cmd/wisp/panel_resident_windows.go:577` | 其余 | Info | `panel thread ending` | 产不可达 | R1 |
| 75 | `cmd/wisp/resident_approval_windows.go:381` | 其余 | Info | `resident gate: [risk] tier taken at construction` | 产不可达 | R1 |
| 76 | `cmd/wisp/resident_approval_windows.go:478` | 其余 | Warn | `resident gate: [risk] source unreadable at construction; the gate falls back to the compiled approval constants` | 产不可达 | R1 |
| 77 | `cmd/wisp/resident_approval_windows.go:497` | 其余 | Warn | `resident gate: [risk] source present but refused at construction; the gate falls back to the compiled approval constants` | 产不可达 | R1 |
| 78 | `cmd/wisp/resident_approval_windows.go:538` | 其余 | Info | `"audit: " + line` | 产不可达 | R1 |
| 79 | `cmd/wisp/resident_approval_windows.go:868` | 其余 | Error | `approval: 卡片无处呈现（本进程没有悬浮球窗口），已 fail-closed 拒绝` | 产不可达 | R1 |
| 80 | `cmd/wisp/resident_approval_windows.go:881` | 其余 | Info | `approval: 常驻进程显示一张确认卡片` | 产不可达 | R1 |
| 81 | `cmd/wisp/resident_approval_windows.go:919` | 其余 | Warn | `approval: L1 窗口挂起期间取消键未借到，本张卡片无法用该键否决` | 产不可达 | R1 |
| 82 | `cmd/wisp/resident_approval_windows.go:936` | 其余 | Warn | `approval: 常驻进程卡片进入醒目提示` | 产不可达 | R1 |
| 83 | `cmd/wisp/resident_audio_windows.go:251` | 其余 | Warn | `audio: config unreadable at boot; the capture leg is built from the compiled default table` | 产不可达 | R1 |
| 84 | `cmd/wisp/resident_audio_windows.go:264` | 其余 | Info | `audio: capture leg not built` | 产不可达 | R1 |
| 85 | `cmd/wisp/resident_audio_windows.go:285` | 其余 | Error | `audio: capture leg refused to start` | 产不可达 | R1 |
| 86 | `cmd/wisp/resident_audio_windows.go:299` | 其余 | Info | `audio: capture armed but muted at boot; device not opened` | 产不可达 | R1 |
| 87 | `cmd/wisp/resident_audio_windows.go:314` | 其余 | Error | `audio: capture device unavailable` | 产不可达 | R1 |
| 88 | `cmd/wisp/resident_audio_windows.go:321` | 其余 | Warn | `audio: gate reports not open without a device error` | 产不可达 | R1 |
| 89 | `cmd/wisp/resident_audio_windows.go:326` | 其余 | Info | `audio: capture leg running` | 产不可达 | R1 |
| 90 | `cmd/wisp/resident_audio_windows.go:339` | 其余 | Error | `proc: 第 4 步（停采集线程）未能注册为钩子` | 产不可达 | R1 |
| 91 | `cmd/wisp/resident_audio_windows.go:363` | 其余 | Info | `slog.Info(line)` | 产不可达 | R1 |
| 92 | `cmd/wisp/resident_audio_windows.go:366` | 其余 | Error | `audio: capture leg stop returned an error` | 产不可达 | R1 |
| 93 | `cmd/wisp/resident_ball_windows.go:178` | 其余 | Warn | `panel gesture could not be handed over` | 产不可达 | R1 |
| 94 | `cmd/wisp/resident_ball_windows.go:399` | 其余 | Error | `ball: the resident leg could not create the floating ball window` | 产不可达 | R1 |
| 95 | `cmd/wisp/resident_ball_windows.go:435` | 其余 | Info | `ball: hotkey reload bridge armed (polls config, rebinds on change)` | 产不可达 | R1 |
| 96 | `cmd/wisp/resident_ball_windows.go:447` | 其余 | Error | `"ball: " + line` | 产不可达 | R1 |
| 97 | `cmd/wisp/resident_ball_windows.go:453` | 其余 | Info | `ball: the resident leg created the floating ball window` | 产不可达 | R1 |
| 98 | `cmd/wisp/resident_ball_windows.go:503` | 其余 | Info | `ball: resident leg destroyed the ball window, its tray icon and its hotkeys` | 产不可达 | R1 |
| 99 | `cmd/wisp/resident_ball_windows.go:662` | 其余 | Warn | `mute gesture arrived before its executor was attached` | 产不可达 | R1 |
| 100 | `cmd/wisp/resident_ball_windows.go:670` | 其余 | Warn | `mute gesture found no gate to turn` | 产不可达 | R1 |
| 101 | `cmd/wisp/resident_ball_windows.go:679` | 其余 | Info | `mute gesture turned this process's capture gate` | 产不可达 | R1 |
| 102 | `cmd/wisp/resident_ball_windows.go:690` | 其余 | Warn | `ball gesture arrived with no executor` | 产不可达 | R1 |
| 103 | `cmd/wisp/resident_ball_windows.go:704` | 其余 | Warn | `cancel hotkey fired with no executor` | 产不可达 | R1 |
| 104 | `cmd/wisp/resident_ball_windows.go:710` | 其余 | Info | `ball: the cancel key was handled by the injected approval gate` | 产不可达 | R1 |
| 105 | `cmd/wisp/resident_ball_windows.go:720` | 其余 | Info | `ball: drag ended` | 产不可达 | R1 |
| 106 | `cmd/wisp/resident_ball_windows.go:754` | 其余 | Warn | `tray exit requested` | 产不可达 | R1 |
| 107 | `cmd/wisp/resident_ball_windows.go:760` | 其余 | Info | `tray exit requested` | 产不可达 | R1 |
| 108 | `cmd/wisp/resident_ball_windows.go:768` | 其余 | Info | `tray exit requested again` | 产不可达 | R1 |
| 109 | `cmd/wisp/resident_task_source_windows.go:240` | 其余 | Warn | `"task source: " + injectedRefusal` | 产不可达 | R1 |
| 110 | `cmd/wisp/resident_task_source_windows.go:244` | 其余 | Warn | `task source: 任务入口未启用` | 产不可达 | R1 |
| 111 | `cmd/wisp/resident_task_source_windows.go:289` | 其余 | Error | `task source: 任务管线未装配` | 产不可达 | R1 |
| 112 | `cmd/wisp/resident_task_source_windows.go:312` | 其余 | Error | `task source: 面板的工作区切换腿未接上本进程的路径范围` | 产不可达 | R1 |
| 113 | `cmd/wisp/resident_task_source_windows.go:363` | 其余 | Error | `task source: 退出序列步骤未能注册为钩子` | 产不可达 | R1 |
| 114 | `cmd/wisp/resident_task_source_windows.go:374` | 其余 | Error | `task source: 注入的任务没有起来` | 产不可达 | R1 |
| 115 | `cmd/wisp/resident_task_source_windows.go:442` | 其余 | Warn | `task source: 控制台输入流断了` | 产不可达 | R1 |
| 116 | `cmd/wisp/resident_task_source_windows.go:487` | 其余 | Info | `task source: 常驻腿提交一发任务` | 产不可达 | R1 |
| 117 | `cmd/wisp/resident_task_source_windows.go:499` | 其余 | Warn | `task source: 任务没有正常收口` | 产不可达 | R1 |
| 118 | `cmd/wisp/resident_task_source_windows.go:511` | 其余 | Info | `task source: 退出第 1 步完成：任务门已关闭，不再接新任务` | 产不可达 | R1 |
| 119 | `cmd/wisp/resident_task_source_windows.go:540` | 其余 | Warn | `task source: 退出第 7 步到点仍有任务在写，本轮不关闭存储（交给进程退出）` | 产不可达 | R1 |
| 120 | `cmd/wisp/resident_task_source_windows.go:546` | 其余 | Info | `task source: 退出第 7 步完成：日志与存储已冲刷关闭` | 产不可达 | R1 |
| 121 | `cmd/wisp/resident_windows.go:223` | 其余 | Error | `panel host: this process has no panel` | 产不可达 | R1 |
| 122 | `cmd/wisp/resident_windows.go:258` | 其余 | Warn | `ball: [hotkey] source unreadable at construction; falling back to the compiled defaults` | 产不可达 | R1 |
| 123 | `cmd/wisp/resident_windows.go:264` | 其余 | Info | `ball: [hotkey] view read from config.toml` | 产不可达 | R1 |
| 124 | `cmd/wisp/resident_windows.go:334` | 其余 | Error | `proc: 第 3 步（取消任务根）未能注册为钩子，退出时待确认项不会被有序拒绝` | 产不可达 | R1 |
| 125 | `cmd/wisp/resident_windows.go:640` | 其余 | Info | `panel host: starting a fresh panel thread for this gesture` | 产不可达 | R1 |
| 126 | `cmd/wisp/resident_windows.go:699` | 其余 | Info | `panel host: starting a fresh panel thread for this request` | 产不可达 | R1 |
| 127 | `cmd/wisp/secret.go:385` | 其余 | Info | `wisp secret: stored dpapi blob` | 产不可达 | R17 |
| 128 | `cmd/wisp/secret.go:530` | 其余 | Warn | `audit + " fields="` | 产不可达 | R17 |
| 129 | `cmd/wisp/secret.go:532` | 其余 | Info | `slog.Info(audit)` | 产不可达 | R17 |
---

## §5 今天的"掉"面已经**换了形状**（本表最重要的一句话）

`130-r1` 的早期缓冲（`b19c6612e`）落地之后，"A/B 窗口"不再是"记录会掉"的同义词：

1. **装得上听众的腿**（`resident` / `run` / `secret` / `models ensure` / `slo`）——A 组与 B 组的记录**今天落盘**。`130-v1` §4 的 post 发逐条点名：`slo -settle` 盘上 2 条，**记录 0 = A1**（resolver 裁决）、**记录 1 = R14w**（pipeline 自己封 `logs\` 那条 WARN）。⇒ 这两条今天**不是**丢失现场，而是"被缓冲捞回来的那两条"。
2. **无听众的腿**——记录仍然掉，而且掉法与"早不早"无关：`providers` / `doctor` / `panel-assets` / `models list`+`verify` / `version` / `help` / **`panel-inbound`** / **second launch**（`resident_windows.go:88` 就 return） / `cmd/balldebug`。
   这九条腿上 `FlushEarlyLogRecords` 与 `InstallAsDefault` 都**没有**调用者（生产调用者各 1：`cmd/wisp/logsink.go:173`、`cmd/wisp/slo_windows.go:281`）⇒ 缓冲**永不冲刷**，随进程消失，只剩 stderr 镜像。
3. **缓冲溢出**——`earlyLogMaxRecords` 64 条 / 64 KiB（`logging.go:514-515`），满了停止追加并计数；冲刷时补一条 `observe: early log buffer overflow`（`logging.go:625`）带 `dropped`/`first_kept`/`capacity`。**丢的条数会成为盘上的一条事实**，但那些记录本身不回来。

⇒ 所以本表的"判"列不是按时刻判的，是按**腿**判的：同一条站点在 `run` 腿上是"落盘"，在 `panel-inbound` 腿上是"掉"。

### §5.1 与件/判语的行号对照（判语 §2(c) 的 16 处，逐处复跑）

| 判语给的号 | 本笔重跑 | 判 |
|---|---|---|
| B1 `logging.go:204` | **257**（boot sweep，`newRollingWriter`→`sweepLocked`）；**`379` 是 `flushLoop` 的 tick，不是 B 窗口** | 判语对一半：`257` 是 B1，`379` 要移出 B 组（本表行 41 / 42） |
| 四腿装点 57/164/219/276 | **109**（`resident_windows.go`）/ **222**（`run.go`）/ **226**（`secret.go`）/ **339**（`models.go`） | ✅ 全对 |
| `logsink` 144/148/165 | **149**（`InitLog`）/ **153**（`SetDefault`）/ **184**（install 记录）；补发在 **173** | ✅ 全对 |
| agent 七枚 373…933 | **391 / 427 / 650 / 723 / 764 / 804 / 1016** ＋构造点 `loop.go` **215** | ✅ 全对 |
| memory 14 枚 | **16 枚**（`artifacts.go:102/271/321`、`open.go:236/240/343/507/522`、`retention.go:109/120/141/167/174/180`、**`writer.go:148/184` 两枚是新增的**） | ⚠ 判语/件都少两枚 |
| 捕获式合计 21 枚 | **25 枚**（memory 16 ＋ agent 9：`compress.go:487/516` 也是新增的） | ⚠ 同漂 |
| A 组 6 枚 `resolve.go:147/151/158/166/176/182` | **逐字仍在原位** | ✅ |
| B2 `goroutine.go:271`、C1 `winsec_windows.go:143`、F1 `secret/migrate.go:198` | **逐字仍在原位** | ✅ |

### §5.2 与件/判语的三处**结论性**差异（不只是行号）

- **"其余 24 枚"⇒ 118 枚**（§2 的恒等式）；19/24 都不是今天的数。
- **`winsec_windows.go:143` 进 B 组**（判语对）：它在 `observe.InitLog` 内部就发（pipeline 自己封 `logs\` 那次）⇒ 属"`InitLog` 已跑、`SetDefault` 未跑"的窗口；件里只把它记在 C 组（别的腿的封印）⇒ 漏了 B 这一半。
- **新补一枚 C'**：`internal/proc/singleinstance_windows.go:171`（second launch 的 `SignalExistingInstance`）—— 这条腿在 `resident_windows.go:88` 就 `return`，**整条腿没有听众**。件与判语都没有它（件把 `proc` 的 4 枚归进"其余、都在装之后"；`shutdown.go` 那 4 枚确实如此，`singleinstance` 这枚不是）。

---

## §6 未验证（明写，别把本表当运行时读数）

1. **本笔没有重跑 exe 的 A/B**：`R14` 的"落盘半"与 `R14w` 全文都引自 `130-v1` §4（它建了两棵 `git archive` 树各跑 5 发）。本笔只现跑了"发声半"（§3 `R14` 那条命令）。
2. **118 枚"其余"的判是静态顺序证据**（装点行号 + 生产调用者数），不是逐站点运行时读数；改装配顺序就要重走 §1 那条命令并重判 R1/R2 的射程。
3. **`wisp panel-inbound` 腿上的掉**：静态（该文件 `:50-55` 自述不装 sink）＋ 三把调用者式子；**未跑 exe**。
4. **溢出支（64 条 / 64 KiB）没有造**：形状在 `logging.go:509-534,620-630`，读数无。
5. **`cmd/wisp` 整包四数未跑**（既有环境上限，判语 §7 同样未跑）⇒ 本表不声称 `cmd/wisp` 全绿。
