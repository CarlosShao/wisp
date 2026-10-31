# 票 42 裁决腿 `42-v1` —— watchdog（非实现者裁决表）

> 裁决腿＝`42-v1`（非实现者；实现＝`42-r1` commit `9d99af5c4`，现跑时 HEAD 即该 commit、分支 `dev`）。
> 现跑 2026-10-31 11:01–11:15 +08。写面只有本文件＋仓外 `D:\tmp\wisp\20261028\42v1\{logs,mut}`。
> 读数全部自己现跑；实现件读数（`docs/evidence/s1/42r1-impl.md`）只当被验对象、不预认。
> 工作树在飞脏＝41a2-r1 的 `cmd/wisp` 三枚，本腿零接触；突变全部走仓外 `-overlay`，工作树零突变窗口。

## 0. 一句话结论

**五格全翻：AC#2/#3/#4/#5 净翻，AC#1 翻＋带注（「panel hidden」宿主半边具名交回接线腿）⇒ 可带残余结案**；
R-42r1-1 定性成立（真缺口、修法方向正确、归 internal/config、该现在开新小票）；生产装配缺口具名归宿主接线腿、
排 cmd/wisp 队列（41a2-r1 落地后、与 273 同波或其后）。

## 1. 门禁四数（本腿自跑，件在 `D:\tmp\wisp\20261028\42v1\logs\`）

| 面 | 尺 | rc | 读数 |
|---|---|---|---|
| 主面 | `go test ./internal/watchdog/ ./internal/observe/ -count=1 -v`（件 01） | **0** | watchdog **12 PASS／0 FAIL／0 SKIP**（0.38s）＋observe **82 PASS**（`--- PASS` 共 94−12）＝与实现件读数同 |
| 邻接 | `go test ./internal/config/ ./internal/statemachine/ ./internal/session/ -count=1`（件 04） | **0** | 三包全 ok（1.9/0.2/2.7s）＝零回归 |
| 零字节 | `git show 9d99af5c4 -- internal/observe/thresholds.go` 行数 | —— | **0 行**；commit 名册恰 8 枚新文件、无任何 `M` |
| 突变 | 三发 `-overlay`（§2） | A=1／B=1／还原=0 | 还原态＝同一棵未动树两次 `-count=1` 全绿 |

## 2. ★★★ 突变复刻＋反例探针（仓外 `mut\`，`-overlay`）

- **发 A 复刻（摘 sustain）**：`if l.overTicks < l.sustainTicks {` → `if true {`（`mut/watchdog.mutA.go`＋`overlayA.json`）
  ⇒ **rc=1、恰 7 枚红**（件 03）：`TestConfigMtimePoll…`（`[42] the resource duty never ran on the loop`）、
  `TestSleepingOverCapDropsUnloadablesThenWithinCap`（`panel-webview unload calls = 0, want 1 (one pass over the sustained breach)`）、
  `TestSleepingOverCapStillOverAfterPassAlertsOnce`、`TestArmedOverCapRetainsKWSAndWarns`、`TestArmedWarningMapsToD43Row9`、
  `TestLogAuditNoUnloadOfNonWhitelistedEver`（`no unload action lines found`）、`TestUnloadPanicIsolated`。
  **三枚红句逐字与实现件 §3 发 A 对上**；settle 两枚不红（不依赖卸载趟）＝外科形，计数对上（7）。
- **自加发 B＋反例探针**：发 B＝卸载集改走注册表（`for _, name := range whitelist {` → `range l.Registered()`）；
  探针＝注册表**只装非白名单**（kws／audio-capture／sqlite-db-writer／session-engines／agent-context＋一枚自造名），
  SustainTicks=1 跑两 tick。**pristine ⇒ 探针绿**：Unload 全 0 调、`unload action` 行 0、`unload pass complete` 行 2
  （趟真跑了却什么都没够到）⇒ 白名单∩注册表的结构保证成立；**同探针×发 B ⇒ rc=1、恰 4 枚红**
  （探针＋审计＋Sleeping＋Armed）⇒ 探针有牙、审计不是恒真。探针文件只存仓外、未进仓。

## 3. ★ 票面 5 格逐条（断言形全部现读现跑）

1. **AC#1 ①可做，翻＋带注**。`TestSleepingOverCapDropsUnloadablesThenWithinCap`：30MB>25MB×3 tick ⇒
   `panel-webview`＋`bm25-index-cache` 各 dropped=1、before 30MB/after 20MB、0 告警；不回落臂
   `…StillOverAfterPassAlertsOnce`：全 skip ⇒ Overrun 恰 1/episode；Armed `…RetainsKWSAndWarns`：100MB>90MB、
   kws calls=0、bm25=2、Detail 具名 "KWS"+"disabling"。★**「panel hidden」半格归口成立**：真隐藏动作在面板宿主
   （cmd/wisp 地界，派单明令本腿 ⛔），本腿交付白名单机制＋`panel-webview` 组件语义（用户在看 ⇒ skip 分支有测）＋
   `Register` 装配缝 ⇒ **该半格随宿主接线腿收口**（残余名见 §7）。
2. **AC#2 净翻**。`TestThreeConsecutiveSettleFailuresReachWatchdogAlertWithRestart`：3 连败 ⇒ AlertFatal 恰 1、
   计数重置；宿主缝真 Dispatch `EvWatchdogFatal` ⇒ `sm.State()==StateWatchdogAlert`＋sink 收
   `Effect{D43:35, Name:"ui.one-click-restart"}`；连续性（败败**成**→第 6 败才第 2 发 Fatal）有测；
   `released=false` 按契约必败；未知目标 config 类拒绝。
3. **AC#3 净翻**。`TestLeakIndicatorsAlertWithSnapshot`：handles=601/GDI=201/goroutines=7（Sleeping cap 6）⇒
   3 条 Leak 各带非空 JSON 快照（metric/state/穿越树本体）、同 episode 去重、handles 回落再越重臂（共 2）、
   `panel.unloadCalls()==0`、日志自述 `no auto-unload`。USER record-only（不发明 cap）＝正确裁量。
4. **AC#4 净翻**。`TestLogAuditNoUnloadOfNonWhitelistedEver`：白名单 4＋NEVER 5 全注册、四态各一趟；审计三层
   （动作行 component∈白名单、outcome∈{dropped,skipped}、NEVER 名 0 行＋Unload 0 调）＋趟行带
   before/after/cap/dropped。＋本腿独立反例（§2 探针）。`watchdog.go:467-499` 读码复认：卸载趟只走冻结白名单、
   按名查注册表，`Register` 对非白名单名只做透明登记。
5. **AC#5 净翻**。`TestConfigMtimePollFiresReloadEventsIdenticalTo05Contract`：**真 shared body**——同一
   `driveReloadScenario` 双驱动（05＝裸 `CheckAndReload` 直调、42＝同一 Manager 挂 watchdog tick）、同一
   `assertReloadContract`、末尾 `reflect.DeepEqual(got05, got42)`；链式三改（[voice] reload／[ball] hot／[app]
   restart）两驱动事件序列全同；42 侧另证一条循环双职责（`CountByName("watchdog")==1`＋同循环日志同时含
   `config poll applied changes` 与 `unload action`）。

## 4. ★★ D43#9/#35 核表（`internal/statemachine/table.go` 现跑行号）

- **#9＝table.go:69-72**：`From: StateArmed, Event: EvWatchdogOverrun, To: StateArmed, Guard:"watchdog-overrun",
  SideEffects:["kws.keep-alive-alert"]`。实现侧：KWS 结构性不可卸（不在白名单）＋`warnArmedKWS`（watchdog.go:447-461）
  在 sustain 窗发 AlertOverrun（Detail 具名 KWS retained＋建议停用、每 episode 去重、并压掉通用 Overrun 防双响）
  ⇒ 宿主映射后「留态＋keep-alive 告警」两味与实现一致。Guard 的宿主实现归接线腿。
- **#35＝table.go:256-257**：`From: AnyState, Event: EvWatchdogFatal, To: StateWatchdogAlert,
  SideEffects:["ui.one-click-restart"]`。实现侧：3 连败 ⇒ AlertFatal（class=resource）⇒ settle_test 真机
  Dispatch ⇒ WatchdogAlert＋`ui.one-click-restart` 落 sink。一致。

## 5. ★ observe caps 冻结

`caps.go` 48 行全为只读 accessor（MemCapBytes/CPULimitPercent/GDILimit/HandleLimit/GoroutineLimit），
零 loosening／零再推导；`GoroutineLimit` 仅 Sleeping/Armed gated、其余 record-only。`TestFrozenCapAccessors`
钉 6 态内存行（25/90/350/350/600/700-target MB）、CPU 六值、GDI 200、handles 600、goroutines 6/7＋四态非 gated。
thresholds.go 零字节（§1）⇒ 一张表两个消费者、无第二张表（188 定式）。

## 6. ★★ R-42r1-1 定性（新缺口，真）

- **病灶**：`internal/config/manager.go:160-168`——`LoadFile` 失败支**无条件收养 stat**（`seenValid=true`）
  ⇒ 瞬时 IO 读失败（temp+rename 竞态的 sharing violation 一类）把该次写**永久隐藏**（下一 tick mtime+size
  无变化 ⇒ 幂等早退）。我读码复认该支无条件收养形状；42 腿并发轮询实证且其测试件自陈 workaround
  （config_poll_test.go:100-104／186-198 的「真实写入者会再写」重种建模）。
- **修法判断：「解析类收养／IO 类不收养」方向正确**。解析/校验类失败＝文件不会自己变好，收养防每 tick
  重报错（05 契约原意）成立；IO 类失败＝盘上内容通常有效，不收养 ⇒ 下一 tick 重试，成功即生效、失败每 tick
  一条 Warn（有界、自愈、cmd/wisp 223 已有 describeReloadFailure 出线），严格优于永久隐藏。可行性：错误链
  今天已可分——`readConfigFile` 的 IO 支（loader.go:81 `config.toml read`，原始 `os.ReadFile` err 在链上、
  可 errors.Is/As）与 parse/校验支（loader.go:102/:118 等）是两条出线；落法建议失败支按「读 IO 错 ⇒ 不收养」
  一刀切，并为该分界装一枚定向突变牙。
- **归属与开票**：归 **internal/config 所有者**。**该现在开新小票**：226 已终裁且射程是写路径 statOwnWrite
  收窄、不是这条 reload 失败支，并入会搅浑两格 AC；reload 轮询今天已由 223 在 cmd/wisp 生产接线、42 宿主
  接线再加第二个生产轮询者 ⇒ 暴露是活的。⛔ 本腿不开票不改票面 ⇒ 编排者开票，票面可直接引用本节＋
  42r1-impl.md §6。

## 7. ★ 生产装配缺口具名（何时派）

生产调用者 0（`cmd/wisp` 无 `watchdog.New/Start`）＝具名真缺口。归宿主接线腿：真 TreeReader（proc 侧）、
StateFn 接状态机、OnAlert 映射（Armed 的 AlertOverrun→EvWatchdogOverrun/D43#9、AlertFatal→EvWatchdogFatal/
D43#35、Leak→日志/UI）、四组件注册（panel-webview 真隐藏、long-output-buffers、goja-vm if enabled、
bm25-index-cache）＋ Start。**时机**：cmd/wisp 现被 41a2-r1 独占 ⇒ **41a2-r1 落地后派**；OnAlert 映射依赖
273（状态机无 sink、AC#2 甲＝装配根注真 Sink，未结）的真收件人 ⇒ **与 273 同波或其后**（两票两缺口、
一波装配；同波不许两腿同写 cmd/wisp）。

## 8. 结案意见（先例 15/26/347/382-done 同形）

- 五格：AC#2/#3/#4/#5 **净翻**；AC#1 **翻＋带注**（「panel hidden」宿主半边具名交回接线腿；白名单机制＋
  skip 分支＋装配缝已交付并有牙）。
- 残余（全部具名、不阻结案）：①宿主接线（§7，新票）②R-42r1-1（§6，新小票归 config、现在开）③
  `internal/watchdog/doc.go:18` 的 `DEFERRED(watchdog loop/thresholds)` 摘标记＝225 票写面（同 28 先例）
  ④USER record-only／CPU 无卸载语义＝票面外如实登记，非缺口。
- 翻勾与票面状态行改动归编排者（本腿 ⛔）。
