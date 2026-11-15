# 46 — S6 acceptance: idle/KWS gates, settle, power pre-mortems, cost verification

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 41-kws-wake-word, 42-watchdog, 43-power-lifecycle-events, 44-costmeter-c23,
45-diagnostics-guards
**Parallel slots:** 1
**Spec refs:** S6 done (D44), D32 gates, §10 #7, SPEC-10 §3/§7

## What to build
The S6 slice gate: measured evidence that the resident app truly sleeps — Sleeping and Armed
budgets on real hardware, settle-return ≤10s across workload types, full power-event pre-mortem
matrix, and CostMeter correctness sign-off.

## Key constraints (evidence, not assertion)
- SLO runs (sampler, real desktop): Sleeping ≤25/40MB (per spike verdict) CPU ≤0.5%, handles
  <300, GDI <200, goroutines ≤6, zero periodic writes, zero long connections; Armed ≤90/110MB
  CPU ≤2% with KWS live; Warm ≤350/400MB.
- Settle: after task types {pure-text, tool-heavy, voice, panel-used, model-download} → tree
  RSS back under state cap ≤10s each (FreeOSMemory counter >0 every time).
- Pre-mortems real-triggered: suspend/resume ×2 cycles with mid-task; battery-saver mode
  (simulated) → Warm 30s/Conversation 15s + KWS off; clock jump +5h mid-session; GDI-leak
  fixture crossing 5000; WebView2 masked (no-panel mode + native L2); mic exclusive-hold;
  disk-full (45) — each with correct D37 classes.
- Cost: 7-day simulated ledger fixture reconciles; budget pause e2e.
- Watchdog: Armed-no-KWS-unload rule exercised; settle-failure → WatchdogAlert drill.
- Adversarial acceptance (different agent): stub scan, SLO-integrity (tree metric), contract
  conformance, DEFERRED cross-check, seven-ban scans.

## Out of scope
- Concurrency (S7); plugin runtime; signing/dist.

## Acceptance criteria
- [ ] All SLO rows measured and archived in docs/SLO.md appendix with machine + date.
- [ ] Five settle scenarios each ≤10s with counters.
- [ ] Seven pre-mortem scenarios executed with evidence + correct error classes.
- [ ] Cost ledger reconciliation green; pause e2e green.
- [ ] Adversarial acceptance report committed; gap audit run (per-slice rule).

## Progress log (append-only, newest last)
---

## S6 验收结果（2026-11-15 12:1x，`694cf184d`＋`f18d22ef6`）：**满足 0／总 5**，五格一格不勾

**尺1**：两条提交在，`694cf184d`（主件）＋`f18d22ef6`（精度修正），都只改证据件。
**尺2**：`docs/evidence/s1/46-s6-acceptance.md` **27255 B／85 行**；原始读数 `.scratch/wisp/probes/46/` 8 份。
**尺3**：腿自跑 `./internal/observe`／`./internal/agent`／`./internal/watchdog` 全包＋5 枚点名／`./internal/session`／`./cmd/wisp` 5 枚 —— 全绿。★**它那三条横切发现我逐条自己核过，全部成立**（见下）。

### 逐格判语（**一格不勾**，理由在最后）
| 格 | 判语 | 依据 |
|---|---|---|
| AC#1 SLO 全行归档 | 部分 | 附录只有 A–D（末条票 16），**没有 S6 附录**；盘上唯一六态表六行全 `skeleton`，而 `SLO.md:536-539` 自己写着「当六态已过引用就是假绿」 |
| AC#2 五类 settle | 部分 | 门在盘（`sampler.go:532/537` 要 `FreeOSMemoryCount>0`），但只有合成夹具 261ms；五类任务场景全仓**零枚**；`VerifySettle` **无生产调用者** |
| AC#3 七枚 pre-mortem | 部分 | **真触发 0／7**；类有牙的只 2 枚（GDI→`resource`、独占麦→`audio_device`）；1/3/6 要机主，2/4/5 卡产码，7 卡装配 |
| AC#4 对账＋pause | 部分 | 对账与 pause 腿自跑全绿（真 store＋真 loop、`requests()==0`、`task_log.error_class=budget`）；**7 天台账夹具零枚** |
| AC#5 对抗验收＋缺口审计 | **不满足** | 盘上零 S6 报告、零缺口审计 |

### ★三条横切发现（腿报，**我逐条现核成立**）——这三条票面看不出来，比五格判语更值钱
1. **settle 门在产品里根本不跑**：`cmd/wisp/resident_windows.go:845-848` 现读逐字 —— `wd, wdRoot := startResidentWatchdog(rt, rv, rp)` 之后写「that dispose assembly rides the session-scope wave (**ticket 28**)」然后 `_ = wd`。而**票 28 今天已 `-done`** ⇒ **这条注释立论的前提过期了**，所以 AC#2 连计数都不存在，不只是"缺真机"。**我核**：该段现读逐字如此。
2. **磁盘守卫没装配、跳闸也不分类**：`NewDiskGatedWriter` **产码调用者 = 2 处，全在 `internal/observe/guards_45.go:305/:308` 包内**（**我核**：`internal`／`cmd` 非测试件全扫，包外 0）；`:218-219` 跳闸拒绝是裸 `ErrDiskSpaceLow`，`ClassResource` 只在旁路。
3. **`LogConfig.Disk is set` 是一个不存在的字段**：`guards_45.go:307` 注释逐字这么写，而 `type LogConfig struct`（`internal/observe/logging.go:64`）正文里 **`Disk` 命中 = 0**（**我逐行读过整个结构体**：`Dir`／`Level`／`RollSizeMB`／`RollDays`／`RedactPaths`，没有 `Disk`）。⇒ 注释在描述一个从未存在的字段。
4. **S6 还挂着一行 DEFERRED**：`SPEC-12:80`「完整错误文案体系（含 D37 全 17 类）」归 S6，本票五格**不覆盖它**，而票 216 只落了面板侧 ⇒ 记给 **票 216**（它正卡在 Q-441）。

### 为什么五格一格不勾
AC#1–#4 都是"部分满足"，AC#5 是"不满足"。按 `issues/README` 与铁律，**部分满足＝不勾**。⛔ 我不拿"腿自跑全绿"当满足 —— 那是**下层模块绿**，不是**切片验收通过**。

### ⇒ 本票的真前置不是「41–45 已完成」，是下面这 4 件 + 1 次机主会话
**产码／装配（4 件，本编队可派）**
1. settle 门接进 dispose 路径（`resident_windows.go:845-848` 的 `_ = wd` 那条注释立论已过期）
2. `NewDiskGatedWriter` 真装配到日志链上 ＋ 跳闸按 `ClassResource` 分类
3. `LogConfig` 加 `Disk` 字段（**或**把 `guards_45.go:307` 那句注释改成事实句 —— 二选一，⛔ 留一条描述不存在字段的注释）
4. 五类 settle 任务场景夹具（今天零枚）＋ 7 天台账夹具（今天零枚）

**机主一次会话（5 步，见证据件末节）**：起 wisp 静置 3 分钟量 Sleeping／连说三次唤醒词＋冷载一次量 Armed+Warm／手动睡眠唤醒两次中间留一张待确认卡／系统时间后拨 5 小时再拨回／先开会议软件占麦再启动 wisp 看提示。

⛔ 本票在 4 件产码与那次会话完成前**不许改名结案**。

---

## S6 两处空洞收件（2026-11-15 18:4x，`ce993cd1e`，13 文件 +1373/−31，件 `docs/evidence/s1/46-s6-seams-impl.md`）

**尺1**：提交在。**尺2**：13 件＋证据件。**尺3 我自己复跑**：`Test46DiskGatedWriterHoldsAndReplaysThroughInstallLogSink` **PASS 0.04s**／`Test46SettleGateVerifiesTheDisposeCascade` **PASS 0.13s**；三处赋值点逐字现读核对过（`resident_voice_engines_windows.go:550 gate.request(settleRequest{peak: peak, released: released})`／`cmd/wisp/logsink.go:210 obs.Disk = guard`／`resident_watchdog_windows.go:81 rv.attachSettleGate(gate)`／`internal/observe/logging.go:257 func (p *LogPipeline) DiskHeldStats() (held, lost, replayed int, ok bool)`）。

### 三处空洞的收口（三处都实拆过一次，红了，装回）
| 拆掉什么 | 红句逐字 | 耗时 |
|---|---|---|
| `logsink.go:210`（`obs.Disk` 置 nil） | `46: this pipeline has no gated writer at all (LogConfig.Disk never reached InitLog): nothing is holding bytes, so the log class is ungated` | 0.03 s |
| `resident_watchdog_windows.go:81`（`rv.attachSettleGate(gate)`） | `46: the residency holds settle gate 0x0, want the one this assembly built (0x362153830210): the dispose cascade would verify nothing` | 0.00 s |
| `resident_voice_engines_windows.go:550`（`gate.request(...)`） | `46: the settle verification the cascade should have started never happened within 10s` | 10.00 s |

### 闸门的四个决定（它拍的，我认）
1. **触发点**＝dispose 级联第 3 步（**全进程里唯一「引擎已卸、堆已交回」同时为真的时刻**）
2. **测量**走 `Registry.Spawn("watchdog","settle-gate")`（`VerifySettle` 阻塞整个窗口，⛔ 不能压级联线程；D38b 六枚冻结名册一字节没动）
3. **`released` 用 `FreeOSMemoryCount` 的增量**（谁改回 `debug.FreeOSMemory` 就按契约红）
4. **目标态读 D43 机现态**，非采样态才回落 Sleeping（Armed 90MB vs Sleeping 25MB —— 硬编码会把常驻 KWS 的机器判成每发都不过）；唤醒 `abort` **不计失败**

真校验那发读数：`final_bytes=419430400 cap_bytes=26214400 back_within_cap_ms=-1 free_os_memory_count=1 consecutive_failures=1`。

### ★那句假话的处理：**没改词，改成了真**
旧句 `LogConfig.Disk is set`（当时 `Disk` 在该结构体里 **0 命中**）⇒ 现在**字段真造出来了**（`internal/observe/logging.go:88`）＋`guards_45.go:338-353` 四行可核的组装说明（谁填、谁装、唯一生产点、探针在哪，并明写「旧句描述的是一个不存在的字段」）。⇒ ⛔ 不是把人话删掉，是让那句话变成真的。

### ★★它在装配里撞出一个真缺陷（已记台账 A1887）
磁盘门控写入器**第一次真接上就自锁**：`Write` 被 slog 处理器持锁调用，旧 `Check` 拒绝时 `slog.Warn` ⇒ 记录流回同一枚处理器 ⇒ **同一 goroutine 对不可重入 mutex 自锁** ⇒ **0% CPU 静默挂死**（那枚挂 18.6 分钟的 `wisp.test.exe` 就是它）。
**真机后果**：磁盘 <200MB ⇒ **下一次启动挂在「装日志」这一行、不留任何日志**。修法：`Check` 拆出 quiet `admit`（`guards_45.go:201`）＋警告改由 flusher 在跨越时发（`logging.go:169`／`:189`）＋`DiskHeldStats()`（`logging.go:257`，「被扣住」与「还没写」在文件系统上长得一样，这是唯一分辨读数）。

### 本票仍未结（剩余＝第 4 件夹具 ＋ AC#3 三枚要机主的预演）

