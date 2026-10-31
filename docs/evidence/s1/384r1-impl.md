# 票 384 落地腿 `384-r1` —— watchdog 宿主接线（写腿）

> **尺**：开工 `HEAD=5c3f44953`（`dev`），现跑 2026-10-31 13:52–15:2x +08。★读数逐条现跑，⛔ 转述。
> **写面（4 枚，逐名）**：新 2 枚（`cmd/wisp/resident_watchdog_windows.go`＝装配四缝＋panel 组件、
> `cmd/wisp/resident_watchdog_windows_test.go`＝3 案）＋插挂 1 枚（`cmd/wisp/resident_windows.go`
> **+14 行纯插入**于 43-r1 挂接行之后，`git diff` 可核 43-r1 行零触碰）＋本文件。
> **零字节**：`internal/watchdog/**`、`internal/session/**`、`internal/risk/**`、`internal/speech/**`、
> `internal/statemachine/**`（D43 表一字未动）、`internal/observe/**`（caps/thresholds 零字节）、
> `internal/config/**`（零配置键）、`internal/panel/**`、`internal/agent/**`、`internal/proc/**`、
> `testdata/**`、`frontend/**`、`design/**`、`scripts/**`、`.github/**`、票面、台账、`docs/PLAN.md`、`docs/specs/**`。
> **临时件（只建不删）**：`D:\tmp\wisp\20261031\384r1\{bin,logs,mut,cases15e,export,export-head,head.tar,d22scan.exe}`。
> 突变在仓外 `export` 导出树（`git archive HEAD`＋覆盖本腿 3 件后变异），工作树零突变窗口。
> **staging**：`go test -c -o bin\wisp.test.exe ./cmd/wisp`＋三 DLL 同目录
> （`third_party\sherpa-onnx\{onnxruntime,sherpa-onnx-c-api,sherpa-onnx-cxx-api}.dll`）＋
> `WISP_SILERO_MODEL/WISP_SPIKE_MODELS/WISP_PUNC_MODEL/WISP_TTS_DIR/WISP_KWS_DIR/WISP_15E_TMP`
> ＋manifest 经 cwd=仓根解析（381-v3 清单同套）。

## 0. 一句话

`internal/watchdog`（票 42 落的库，生产调用者 0）第一次有了生产调用者：常驻 boot 在 43-r1 lifecycle
挂接行旁把四缝一次接齐——**真 TreeReader**（`proc.NewTreeSampler(rt.Job)`，与 `wisp slo` 同一 C30
Job 读数面）、**StateFn 接常驻 D43 机器**（`rv.machine`，四同名直映＋其余态 fail-closed 空 SLOState）、
**OnAlert 按 D43 表行映射**（Armed `AlertOverrun`→`EvWatchdogOverrun`/#9 留态、`AlertFatal`→
`EvWatchdogFatal`/#35 `WatchdogAlert`、Leak→日志/UI 材料）、**panel-webview 真组件注册**（真宿主
`PanelManager`；其余三枚白名单名无生产对象、不注册、具名申报）＋ **Start**（D38b 名册名 `watchdog`，
`rt.Registry.Spawn`）。端到端真机读数：真 tick、over-cap 真 unload 趟、#9 留 Armed、3× settle 败落
`WatchdogAlert`；反向判据两发摘映射 ⇒ D43 事件零派发 ⇒ 红。

## 1. ★★★ 票面 5 格逐条（原文逐字＋判定＋读数）

| # | 票面原文（逐字） | 判定 | 读数 |
|---|---|---|---|
| 1 | `AC#1 生产装配：★真 TreeReader／StateFn／OnAlert 映射（#9/#35 逐字）/四组件注册＋Start（★对照 42r1-impl.md 的缝形状）` | **①可做（已做）** | `Test384StartResidentWatchdogRegistersRealLoop` rc=0：生产装配函数真调用，读数 `384-READ production assembly registered=[panel-webview], D38b CountByName(watchdog)=1 (273 open: production machine Sink=nil)`；真 TreeReader=proc.NewTreeSampler(rt.Job)（`wisp slo` 同形）；StateFn=sloStateFromD43Machine（observe/sampler.go 自钉的同名对应：Sleeping/Armed/Warm/Conversation 直映；其余 16 态无 D32 采样行 ⇒ 零 SLOState ⇒ 库的 enforceResources fail-closed「not a D32 sampling state; tick skipped」，⛔ 不给会话活跃态套小 cap——D32③ 有害形状的反面）；OnAlert=dispatchWatchdogAlert（#9/#35 逐字见 §3）；Start=`observe.NewRoot("watchdog-host")`＋`loop.Start(root.Ctx)`（内部 `Registry.Spawn("watchdog",…)`） |
| 2 | `AC#2 端到端：★over-cap ⇒ 真 unload 趟＋D43#9 留态；3× settle 败 ⇒ #35 WatchdogAlert＋one-click-restart（★真机读数）` | **①可做（真机读数已落）** | `Test384ResidentWatchdogEndToEnd` rc=0（1.25s，staging，真 Boot/真 Job/真 TreeSampler/真机器，log 01）：Phase0 真 tick 行=3（Debug 级逐 tick 落盘＝goroutine 真跑）＋名册 `CountByName(watchdog)=1`；Phase1 Sleeping（cap 25MB）真持 64MB 页触内存 sustain 2 tick ⇒ **真 unload 趟=1**（`"msg":"watchdog: unload pass complete"`），白名单无注册 ⇒ 趟空走、树仍超 ⇒ AlertOverrun(Sleeping) 走「无 D43 行」日志臂（no-row receipts=1）且 **D43 事件零派发**（sink #9+#35 计数=0）、机器留 Sleeping；Phase2 经 #5 真拨 Armed（EvKwsEnabled）、持 128MB > 90MB ⇒ **D43#9 effects(kws.keep-alive-alert)=1、机器 STAY Armed**（留态读数）＋映射回执行在；Phase3 `VerifySettle`×3（released=false 按契约必败，`BackWithinCapMS=51`）⇒ 第 3 败 AlertFatal ⇒ **#35 effects(ui.one-click-restart)=1、机器落 `StateWatchdogAlert`**。41a2/381/43 各族零回归见 §5 |
| 3 | `AC#3 反向判据：★摘 OnAlert 映射 ⇒ D43 事件零派发 ⇒ 红（三读数＋还原核）` | **①可做（两发定向突变）** | 仓外导出树两发（log 02/03）：**M1** 摘 `Dispatch(EvWatchdogOverrun)`（行换 `rv.machine.State(), error(nil)`）⇒ rc=1、红句逐字 `384: never observed D43 #9 kws.keep-alive-alert effect (Armed over-cap, mapped) within 5s`，且前三读数照常（tick=3／名册=1／unload 趟=1 ⇒ 红的只有映射面）；**M2** 摘 `Dispatch(EvWatchdogFatal)` ⇒ rc=1、红句 `384: never observed D43 #35 ui.one-click-restart effect within 3s`，且读数自证定向：`D43#9 effects=1`（M2 没碰 #9，照常派发）、`#35 effects=0, machine lands in Armed`。**还原核**（时机=M2 后）：导出树复原 ⇒ 三案全绿（log 04）；主仓 `git status -- cmd internal`＝恰本腿 3 件、本腿文件 MUTANT grep=0、生产文件两枚 Dispatch 行在位（:133/:146） |
| 4 | `AC#4 零回归：★15e/26R3/381 打断案/41a2 KWS 案全绿；drained 抖动只报` | **①可做（案级全绿＋串跑抖动如实申报）** | 五族合跑（log 08/13/14）：15e 6/6＋`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms` **逐位同值**；41a2r1 5 案、381r2 5 案、26R3 3 案、Test43 9/9（log 16）；233/223/255 定向 6/6（log 17）；`Test384` 三案、watchdog/observe/session/config/statemachine 全 ok（rc=0）。**抖动申报（非零回归差）**：四族**全量串跑**时 ASR 冷装载案（`Test41a2r1Wake/WarmGate`）出现 1 枚/发的负载敏感红——**纯 HEAD binary 同分布**（18 PASS/1 FAIL，log 13，同族同形状、案随机），本腿 binary 同分布（log 14，18/1）；四红案**单发复跑全绿**（log 15，本腿 binary）⇒ 案级绿、族串跑的装载预算被前案残留负载挤压 ⇒ 与本腿写面无关（这些测试不经 `runResident`/装配函数）。`drained` 抖动只报：本腿零触碰 voice 面（零字节），无 drained 断言被引入 |
| 5 | `AC#5 纪律格：⛔ 翻框/改票面/台账/push；⛔ 改 internal/watchdog/internal/session 语义；⛔ D43 表；⛔ 配置键` | **①遵守** | `git status --porcelain -- cmd internal`＝`M resident_windows.go`（+14 纯插入）＋`?? resident_watchdog_windows.go`＋`?? resident_watchdog_windows_test.go`（§6 四数）；`internal/{watchdog,session,risk,speech,statemachine,observe,config}` 全部 0 字节（git status 无它们）；D43 表（`internal/statemachine/table.go`）零字节——映射走表内 #9/#35 **现行行**（`table.go:69` `From: StateArmed, Event: EvWatchdogOverrun, To: StateArmed, SideEffects:[kws.keep-alive-alert]`、`table.go:256` `From: AnyState, Event: EvWatchdogFatal, To: StateWatchdogAlert`）；配置键 0 新增（caps 全走 `observe` 只读 accessor：`MemCapBytes/CPULimitPercent/GDILimit/HandleLimit/GoroutineLimit`）；票面/台账/PLAN/specs 零字节；未 push（commit 由本腿落、push 归编排者） |

## 2. ★★★ 开工第一核：OnAlert→D43 的 sink 依赖票 273（现跑申报）

- **273 状态现跑**（2026-10-31 14:0x）：`.scratch/wisp/issues/273-shipping-process-builds-state-machine-without-a-sink-so-every-d43-effect-falls-into-no-op.md` **在、未 `-done`**；7 枚框 **2 勾**（AC#0/#1，10-07 翻）／**AC#2–AC#6 未勾**——AC#2「三形代价表 ⛔ 本票不选形」至今未选形 ⇒ **真收件人（sink）缺**。
- **本腿处置＝做非 sink 半边＋映射写成具名交回**（按派单，⛔ 发明 sink——那是契约面）：
  1. **映射与 Dispatch 是真的**：`dispatchWatchdogAlert` 把两枚有表行的告警**按名派发**进机器；转移真发生、`rv.machine.State()` 可读（`WatchdogAlert` 落态有真机读数）。
  2. **Effect 投递面是 273 的**：生产机器由 `buildResidentVoice` 建构、`Sink=nil` ⇒ `statemachine.New` 的包内 no-op 缺省（`machine.go:64-68`）收下 `kws.keep-alive-alert`／`ui.one-click-restart`。**本腿零发明**：测试侧的 sink（`effectLog384`）是测试自建观察面（273 普查承认的合法测试形状），生产装配一字未带 sink。
  3. **273 收口时本腿零改动**：映射已经 `rv.machine.Dispatch(...)` 走机器自己的 Sink 字段——273 落「装配根注真 Sink」（甲形）后，效果投递自动可见。
- **测试文件的机器带测试 sink 与生产不矛盾的钉**：`Test384StartResidentWatchdogRegistersRealLoop` 用 `Sink=nil` 机器走生产装配函数，读数里具名 `(273 open: production machine Sink=nil)`。

## 3. ★ 装配缝逐项（42-v1 §7 的四件＋Start）

| 缝 | 落法 | 依据 |
|---|---|---|
| 真 TreeReader | `proc.NewTreeSampler(rt.Job)` | 与 `wisp slo`（slo_windows.go:283）同一 C30 JobScope 读数面；零第二张指标路径 |
| StateFn 接状态机 | `sloStateFromD43Machine(rv)`：Sleeping/Armed/Warm/Conversation 同名直映；default ⇒ `observe.SLOState("")` | observe/sampler.go 自钉「Sleeping/Armed/Warm/Conversation ARE D43 machine states (names pinned by test)」；其余 16 态 D32 无采样行 ⇒ fail-closed skip（库的 enforceResources 原生 Warn），不给会话活跃态套错行 cap（D32③ 有害形状的反面）。PanelOpen/WorkPeak 是 SLO 姿态、无生产者 ⇒ 不产 |
| OnAlert 映射 | `dispatchWatchdogAlert`：`AlertOverrun && State==Armed ⇒ Dispatch(EvWatchdogOverrun)`（#9）；`AlertFatal ⇒ Dispatch(EvWatchdogFatal)`（#35）；其余（非 Armed 的 Overrun、AlertLeak）⇒ 日志/UI 材料、**不派发** | 42 契约「the loop never touches the state machine ... the host owns the D43 mapping」；未列转移一律非法（D43），发明边＝改契约 ⛔ |
| 四组件注册 | `panel-webview` 真注册（`panelWebviewUnloadable{mgr: rp}`，rp nil/装配失败则不注册）；`long-output-buffers`／`goja-vm`／`bm25-index-cache` **不注册（具名申报）** | 白名单机制语义＝「卸载趟走冻结白名单、按名查注册表」⇒ 注册表里没有的名字天然跳过。三枚无生产对象：long-output-buffers＝spill 落盘文件（`internal/agent/spill.go`，归票 04 的 500MB LRU 配额，无内存缓冲对象可丢）；goja-vm＝票 51 未落（`internal/tools/registry.go` "ticket 51 lands the Tier-2 goja runtime"）；bm25-index-cache＝`bm25Select` 纯函数无缓存对象（`internal/agent/inject.go`）。⛔ 不造假组件占位 |
| panel-webview 语义 | `Unload`＝nil host ⇒ skip；`!IsCreated` ⇒ skip「no panel window exists」；`IsShown` ⇒ skip「the user may be viewing it」（42 的 skip 分支）；否则 `mgr.Destroy()` ⇒ dropped=true | Destroy 的 webview 半边＝go-webview2 `Destroy()`＝**一记 `PostMessageW(hwnd, WM_CLOSE)`**（webview.go:347）⇒ 任意线程可调，面板线程 pump 自跑关闭链；Destroy（非 Hide）才是卸载语义（Hide 只挪窗口、msedgewebview2 子进程仍驻；WM_CLOSE 链让子进程退出、树私有内存真回落）；可恢复：Show 走 recreate path＋panelHolder 按票 319 重建线程。**具名残余**：真窗口的 Destroy 读数（子进程退出、树回落）属 winlive 面（未批），本腿给的是 PostMessage 线程安全依据＋recreate 引文；`Test384PanelWebviewComponentArms` 钉 nil/无窗两臂 |
| Start | `observe.NewRoot("watchdog-host")`＋`loop.Start(root.Ctx)`；`runResident` 内 `defer wdRoot.Cancel()`（注册在 43-r1 行之后 ⇒ LIFO 先于面板/球 teardown 与 D38(e) 序列停 watchdog） | 223 的 config tick 用同名 `"watchdog"`（SPEC-03 §4.3 预留槽位）；根 cancel 是本腿自己的 defer，D38(e) 十步与 hookableRoster 零触碰 |

## 4. ★★★ 端到端真机读数（log 01，staging `-test.v` 逐字）

```
384-READ real tick lines=3 (the loop goroutine is alive and ticking)
384-READ D38b roster CountByName(watchdog)=1
384-READ Sleeping unload passes=1, unload-action lines=0, no-row alert receipts=1
384-READ D43#9 effects(kws.keep-alive-alert)=1, machine stays Armed
384-READ settle failures=3 (BackWithinCapMS=51, released=false fails by contract); D43#35 effects(ui.one-click-restart)=1, machine lands in WatchdogAlert
384-READ production assembly registered=[panel-webview], D38b CountByName(watchdog)=1 (273 open: production machine Sink=nil)
```

## 5. ★ 零回归（ staging 合跑；binary＝本腿写面 build）

| 面 | 结果 | 出处 |
|---|---|---|
| `Test384` 三案 | 3/3 PASS（E2E 1.25s） | log 01 |
| 15e 6 案＋LATENCY | 6/6；`openDelay=262ms chunk=640ms ballLatency=902ms` 逐位同值 | log 08/13/14 |
| 41a2r1 5 案／381r2 5 案／26R3 3 案 | 全绿（四红案单发 4/4，log 15；全量串跑的装载抖动申报见 §1#4） | log 08/13/14/15 |
| `Test43` 9 案（43-r1 lifecycle 面原样） | 9/9 PASS | log 16 |
| 233/223/255 定向 | 6/6（255 名册未破——本腿零 config 读点） | log 17 |
| `internal/{watchdog,observe,session,config,statemachine}` | 全 ok（rc=0） | §1#4 |
| `go vet ./cmd/wisp` | rc=0 | 门禁段 |

## 6. ★ 四数＋两尺＋d22scan

- **写面四数**：`git status --porcelain -- cmd internal`＝`M cmd/wisp/resident_windows.go`（+14/−0）＋`?? cmd/wisp/resident_watchdog_windows.go`＋`?? cmd/wisp/resident_watchdog_windows_test.go`；docs 侧＋`?? docs/evidence/s1/384r1-impl.md`。提交后 `git show --stat` 只含本腿 4 枚。
- **甲 `gofmt -l cmd`**：空（rc=0）。
- **乙 `gofumpt -l`**（绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`）：本腿 3 枚单点跑＝空。
- **`go vet ./cmd/wisp`**：rc=0。
- **d22scan**（`tools/d22scan` 独立 build→`D:\tmp\…\384r1\d22scan.exe`，扫工作树）：**rc=0 clean**——bans#1-5 internal/=280 cmd/=49、ban#6 frontend/=116、ban#7 internal/tools/=25、ban#8 design/=39、frontend/=116、internal/=715、cmd/=171（log 05）。**正控**：种真 U+2264 进 `internal/config/zz_d22_positive_control_384r1.go` ⇒ **rc=1 且逐字点名** `internal/config/zz_d22_positive_control_384r1.go:4: [emoji] ban #8 glyph in scope internal/ is banned (D23): non-comment text, string literals included; comments are exempt per Q-46(c)`（log 06）⇒ 尺是活的；种后即删、复扫 rc=0（log 07）。

## 7. 未做／具名残余（⛔ 藏）

1. **VerifySettle 的生产调用者**：库的 settle 门设计给 dispose 路径（「runs on the caller's goroutine (the dispose path)」）；常驻进程的 SessionScope dispose 装配是票 28 波次（28「库＋读数齐、装配等下一波」同形状）⇒ 今天生产无人调 `VerifySettle`，`AlertFatal` 的生产触发点随之在 28 落地时激活（映射缝已就位）。E2E 的 3× settle 败读数由测试以 dispose 形状（调用者 goroutine 直调）驱动。
2. **三枚白名单名的生产对象不存在**（§3 表）：long-output-buffers/goja-vm/bm25-index-cache 不注册＝机制语义的诚实形状（表查不到 ⇒ 跳过）；它们各自的落地票（goja=51、spill 配额=04）落对象后由宿主 `loop.Register` 一行接入。
3. **panel-webview 的真窗口 Destroy 读数**：winlive 面（未批），本腿交付 PostMessage 线程安全依据＋recreate path 引文＋两臂单测。
4. **全包 `go test ./cmd/wisp` 在本机不作读数**：全量直跑撞 33 族真窗口案（`0xC0000005` 于 `MARK299 before CreateWindowExW` 后）——winlive 未批的固有形状（43-r1 同口径只跑定向＋族）；本腿以五族＋定向＋邻接包为尺。另两发全包的经验教训已具名：**binary 的 build 树必须=主仓**（export 树 build 会把 `runtime.Caller` 嵌成仓外路径，golden/manifest 类案连坐红）＋**binary 直跑的 cwd 必须=包目录**（源码扫描案以相对路径读源码）。
5. **HEAD 既有红两枚（非本腿，不冒领）**：`TestLockedSuffixNeverTakesTheLockItself33r11`／`TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims`——`waveOutSpeaker.releaseLocked/unprepareLocked`（voice_out_windows.go，381 族文件）带 `Locked` 后缀未入 33r11 名册；**纯 HEAD binary 对照同红**（cmd\wisp cwd），属 381 族债，归其所有者。
6. **四族全量串跑的装载抖动**（§1#4）：ASR 冷装载案在族串跑下 1 枚/发负载敏感红，HEAD 同分布；单发全绿。若验收腿要 19/0，需空载机器或把两枚装载案拆开发跑（属测试基建票，非本腿写面）。

## 8. 停手条件自核（派单四条「未定义即停」）

①273 sink 缺但有半边可做 ⇒ **未触发**（做了全部非 sink 半边＋具名交回）；
②OnAlert 映射须改 D43 表 ⇒ **未触发**（#9/#35 走表内现行行，零字节）；
③caps/白名单与 observe 现值冲突 ⇒ **未触发**（caps 全走 accessor、白名单原样消费、零第二张表）；
④43-r1 lifecycle 挂接与装配正面冲突 ⇒ **未触发**（+14 纯插入于其行后，`git diff` 可核；`Test43` 9/9 复跑绿）。
