# 票 384 裁决腿 `384-v1` —— watchdog 宿主接线（非实现者裁决表）

> 裁决腿＝`384-v1`（非实现者；实现＝`384-r1` commit `51899e7c0`，4 枚 681 行）。现跑 2026-10-31 15:21–15:5x +08。
> 写面只有本文件＋仓外 `D:\tmp\wisp\20261031\384v1\{logs,mut,bin,tree,tree-base}`；实验一律导出树（`git archive 51899e7c0` 主树＋`5056a365c` 基线树），主仓工作树零突变窗口。
> 读数全部自己现跑（件 01–16）；实现件读数（`384r1-impl.md`）只当被验对象、不预认。在飞脏＝`30-w-r1` 写 `cmd/wisp`（15:4x 起 `M resident_windows.go`/`M run.go`/`?? resident_present_windows.go`）、`387-r1`/`386-r1`——本腿零接触。

## 0. 一句话结论

**五格全翻、可带残余结案**（15/28/41/42/43 先例）：必跑清单 1–7 全做——端到端 staging 复刻全读数对上、四缝逐缝读码为真、273 sink 处置合规（判语：**不构成伪完成**）、mutation 两发复刻红句逐字＋还原核绿、零回归（Test43 9/9＋43-r1 行零触碰＋255 名册 20/0/0＋LATENCY 262/640/902 逐位同＋15e/41a2r1/381r2/26R3 计 19/19）、33r11 两红归因坐实 HEAD 既有、d22scan clean＋正控 rc=1 逐字点名＋删后复净。

## 1. 本腿自跑读数（件 `logs\01–16`，rc 全部现取）

| 面 | 尺 | rc | 读数 |
|---|---|---|---|
| ★E2E staging | `wisp.test.exe -test.run ^Test384 -test.v`（导出树 build＋三 DLL 贴 binary＋A1307补 PATH，件 01） | **0** | 真 tick=3／名册=1／Sleeping 64MB sustain⇒**unload 趟=1**＋D43 零派发（sink #9+#35=0、留 Sleeping）／Armed 128MB⇒**#9 effect=1 STAY Armed**／3× settle 败⇒**#35 effect=1 落 WatchdogAlert**；生产装配 registered=[panel-webview]＋Sink=nil 具名 |
| mutation | M1 摘 #9 Dispatch（件 02）／M2 摘 #35（件 03） | **1/1** | 红句逐字见 §4；还原核 pristine 复跑 **0**（件 04）＋树内变异串 grep=0＋主仓两枚 Dispatch 行在位 |
| Test43 | `-test.run ^Test43`（件 05） | **0** | **9/9 PASS** |
| 255 名册 | `-test.run TestTicket255`（件 06） | **0** | **20 PASS/0 FAIL/0 SKIP**（含 HotRowRoster／RosterEvidenceLines／RosterStillMatches 三枚名册尺） |
| 15e＋LATENCY | `-test.run ^Test15e`（件 07+08） | **0** | **6/6**；`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms` **逐位同值** |
| 三族 | 41a2r1 5＋381r2 5＋26R3 3（件 14+16） | **0** | **13/13**（RealDevice 真跑 36.04s 非 skip；ASR 冷装载两案首跑 setup 红＝仓外缺未追踪模型件，junction 后 18.7s/20.9s 绿——见 §9⑦） |
| 33r11 | 两红案×双树（件 09/10） | **1/1** | 两树红形逐字相同＝HEAD 既有，见 §6 |
| d22scan | 导出树净扫＋正控（件 11/12/13） | **0/1/0** | clean；正控 `zz_d22_positive_control_384v1.go:4: [emoji] ban #8 …(Q-46(c)` 逐字点名；删后复净 |

## 2. ★★★ 四缝真伪逐缝核（读 `resident_watchdog_windows.go` 现码＝51899e7c0 树）

1. **TreeReader＝真 C30 Job**：`startResidentWatchdog` :58 `TreeReader: proc.NewTreeSampler(rt.Job)`——与 `wisp slo`（slo_windows.go:283 `observe.NewSampler(proc.NewTreeSampler(rt.Job),…)`）**同一构造形状、同一 Job 对象**，零第二张指标路径。E2E 的 `currentTree384` 真读 PrivateWorkingSetBytes>0（件 01）。
2. **StateFn 16 态 fail-closed**：`sloStateFromD43Machine` :87-98 只直映四名（statemachine/states.go 恰 20 态 ⇒ 其余 **16 态** default ⇒ `SLOState("")`），库侧 watchdog.go:353 对空值 Warn「not a D32 sampling state; tick skipped」。**D32③ 反面核真**：会话活跃态（Listening，经 #5→#7 真拨）答零 SLOState 而非小行 cap——Test384Start :338-343 有钉；`rv==nil||machine==nil` ⇒ Sleeping：机器不存在则无会话可活，最小 cap 方向安全。
3. **OnAlert 零发明边**：`dispatchWatchdogAlert` :131-162 全文只有两枚 Dispatch——`AlertOverrun && State==Armed ⇒ EvWatchdogOverrun`（表行恰一枚：table.go:69 #9 From Armed To Armed 自环）、`AlertFatal ⇒ EvWatchdogFatal`（table.go:256 #35 AnyState→WatchdogAlert）；Leak／非 Armed Overrun 落 :159 日志臂零派发。Dispatch 被机器拒绝时 Warn+return（fail-closed）。对照 D43 表：`EvWatchdogOverrun`/`EvWatchdogFatal` 各恰一行 ⇒ **映射只落表内行、零发明**。
4. **panel-webview 真注册＋三 skip 臂＋Destroy 可恢复**：:63-65 `panel!=nil ⇒ Register(ComponentPanelWebView, panelWebviewUnloadable{mgr: panel})`＝真 PanelManager 指针；三臂 :188-197（nil⇒"panel host not assembled"／!IsCreated⇒"no panel window exists"／IsShown⇒"user may be viewing it"）读的是 panel_host_windows.go:275/:283 真互斥锁态；`Destroy()`（panel_host_windows.go:822-833 自述「After Destroy a later Show creates a fresh window (the recreate path)」）⇒ 可恢复按构造成立；go-webview2 Destroy＝**一记 PostMessageW WM_CLOSE**（模块缓存 v0.0.0-20260205173254 webview.go:**348**，实现件引 :347 同内容差一行）⇒ 任意线程可调、面板 pump 自走关闭链。**三组件不注册＝白名单查表天然跳过**（watchdog.go:476 `for _, name := range whitelist` 按名查注册表；Test384Start :308 钉 registered 恰=`[panel-webview]`）——三引文现核：long-output-buffers＝spill 落盘文件（internal/agent/spill.go D15(3)）、goja-vm＝票 51 未落（internal/tools/registry.go:54 占位）、bm25-index-cache＝纯函数（internal/agent/inject.go:172）。**⛔ 无造假组件**。

## 3. ★★★ 273 sink 处置合规裁——判语：**不构成伪完成**

- **形状核真**：①273 票现读**未 `-done`、7 框 2 勾、AC#2「⛔ 本票不选形」**——甲（装配根注真 Sink）未定案，sink 是契约面，写腿发明 sink＝改契约须人工批准；②`machine.go:64-67` nil ⇒ **显式 no-op**（doc.go:11-12/machine.go:39 冻结自述「nil = explicit no-op」）、:163 Effect 走**机器自有 sink 字段**；③生产构造点现核＝resident_voice_windows.go:526/:1734 两枚 `statemachine.New` 均无 Sink ⇒ no-op；④测试 sink（effectLog384）＝273 普查承认的合法测试形状（bridge_test.go:16-18 先例），且 Test384Start 用 **Sink=nil 机器**驱动生产装配函数、读数具名 `(273 open: production machine Sink=nil)`。
- **判语与理由**：**转移是状态机的、效果投递是 273 的、两件事**。本票交付的缝＝映射＋Dispatch＋留态，全部为真且有三读数（#9 留 Armed／#35 落 WatchdogAlert／机器态可读）；no-op 缺省是 machine 包自己的冻结契约、非本腿引入或遮掩；映射经 `rv.machine.Dispatch` 走机器自有 Sink 字段 ⇒ 273 落甲形（装配根注真 Sink）时**本缝零改动**（形状核②③已证）。伪完成的判据应是「声称的交付面其实没交付」——本票声称的交付面不含效果投递，且交回是具名的（文件头 39-48 行＋impl §2）。**合规**。

## 4. ★★ mutation 两发复刻（仓外 mut\＋`-overlay`，零树突变；各恰 1 行差异）

- **M1**（:133 `Dispatch(EvWatchdogOverrun)`→`rv.machine.State(), error(nil)`）⇒ **rc=1**、红句逐字 `384: never observed D43 #9 kws.keep-alive-alert effect (Armed over-cap, mapped) within 5s`；前三读数照常（tick=3／名册=1／unload 趟=1）⇒ 红的只有映射面（件 02）。
- **M2**（:146 `Dispatch(EvWatchdogFatal)`→同形）⇒ **rc=1**、红句逐字 `384: never observed D43 #35 ui.one-click-restart effect within 3s`；**自证定向**：`D43#9 effects=1`（M2 没碰 #9）＋`#35 effects=0, machine lands in Armed`（件 03）。
- **还原核**：pristine binary（overlay 只换编译输入、盘上零改）复跑三案 **rc=0**（件 04）；树内变异串 grep=0；主仓 `resident_watchdog_windows.go` 两枚 Dispatch 行在位。

## 5. ★★ 零回归

- **43-r1 行零触碰**：`git diff 5056a365c 51899e7c0 -- cmd/wisp/resident_windows.go`＝**恰 +14 纯插入**（43-r1 挂接行后一枚 hunk、43-r1 行零触碰）；`git show --stat 51899e7c0`＝恰 4 枚（cmd/wisp 3＋impl 证据表），D43 表/thresholds/票面/台账零字节。
- Test43 **9/9**（件 05）＋255 名册 **20/0/0**（件 06）＋LATENCY **262/640/902 逐位同**（件 08）＋15e/41a2r1/381r2/26R3 **19/19**（件 07/08/14/16）。
- `BackWithinCapMS=66`（本腿）vs 51（实现件）：settle 时序读数、非钉值断言、负载相关——如实记，非差。

## 6. ★★ 33r11 两红归因核——**HEAD 既有说成立**（双树对照同红）

- 两案在 51899e7c0 树（件 09）与基线 5056a365c 树（件 10）**红形逐字相同**：`waveOutSpeaker.releaseLocked`（voice_out_windows.go:338）／`unprepareLocked`（:376）带 Locked 后缀未入 33r11 名册；Test256 红＝`unexpected Options field "Answer"`（7/11 fields, want 6）。
- 归因精化（对实现件 §7.5 的一处更正）：33r11 名册两枚＝**381 族**（voice_out_windows.go）✓；但 Test256 红与 Locked 名无关——`Answer:` 经 pickaxe 引入 commit＝`b63d56413`（**票 375-r1**），测试件最后触＝`a04a095f7`（2026-10-04，265-r1c）⇒ 第二红属 **375 族债**。两红均早于 384、非本腿，归其各自所有者。

## 7. ★ d22scan（导出树）

净扫 **rc=0 clean**（bans#1-5 internal/=283、cmd/=49、ban#6 frontend/=116、ban#7 internal/tools/=25、ban#8 internal/=719、cmd/=171——与实现件同族、±本腿 3 枚新 .go；design/ 分母 14 vs 在册 39＝主仓未追踪 design 件、commit 内全查零违规）。正控种 U+2264 ⇒ **rc=1** 逐字点名（§1 表）⇒ 尺是活的；删后复净 **rc=0**。⛔ 未扫 docs/（ban#8 无 docs/ scope，与在册口径一致）。

## 8. ★★★ 票面五格逐条判

1. **AC#1 生产装配——能翻**。四缝＋Start 全真（§2 逐缝读码）；Test384Start rc=0、registered 恰 `[panel-webview]`、名册=1。
2. **AC#2 端到端——能翻（带 273 注）**。本腿 staging 复刻全读数对上（§1 件 01）：真 tick／真 unload 趟＋D43 零派发／#9 留 Armed／3× settle 败⇒#35 落 WatchdogAlert。`one-click-restart`＝D43#35 SideEffect 名，effect 已派发并落测试 sink；**真 UI 消费者随 273 收口**（§3 判语：非本票缺角、具名交回）。
3. **AC#3 反向判据——能翻**。M1/M2 两发红句逐字＋还原核（§4）；「摘映射 ⇒ D43 零派发 ⇒ 红」成立且外科定向。
4. **AC#4 零回归——能翻（带抖动注）**。案级本腿自跑 19/19＋Test43 9/9＋255＋LATENCY 逐位同（§5）；`drained` 抖动只报＝本腿零触碰 voice 面 ✓。全量串跑的 ASR 冷装载负载敏感抖动（实现件申报 18/1 同分布）属测试基建债、非本票写面（§9③）。
5. **AC#5 纪律格——能翻**。写面四数对上（stat＝恰 4 枚 681 行）；43-r1 行零触碰；`internal/{watchdog,session,speech,statemachine,observe,config}` 零字节；D43 表/thresholds/配置键/票面/台账零触碰；本腿复核 33r11/273 申报均如实。

## 9. ★ 残余定性（全部具名、不阻结案）

1. **三组件无生产对象**：long-output-buffers/goja-vm/bm25-index-cache 不注册＝白名单查表语义的诚实形状（§2.4 三引文）；各自落地票（51/04）落对象后宿主一行接入。
2. **273 sink 未结**（AC#2 未选形）：效果投递面归 273；映射缝零改动交回（§3）。
3. **ASR 冷装载串跑抖动**：全量四族串跑 1 枚/发负载敏感红、HEAD 同分布（实现件 18/1×2）；本腿单发/定向 19/19 绿——测试基建债（装载案拆发或空载尺），归测试基建票。
4. **33r11/256 债**：两红 HEAD 既有（§6）；381 族补名册或改名、375 族对 Test256 重读 Options 场——归各自所有者。
5. **VerifySettle 生产调用者**：随票 28 dispose 装配激活（AlertFatal 生产触发点随之活）；映射缝已就位。
6. **panel 真窗口 Destroy 读数**（子进程退出、树回落）＝winlive 面（未批），本票交付 PostMessage 线程安全依据（webview.go:348）＋recreate path 引文＋两臂单测。
7. **staging 经验（尺核教训，供复跑腿）**：`third_party/`、`build/` 不入 git ⇒ 导出树复跑须自带三 DLL 贴 binary（A1308 正解）＋A1307补 PATH＋`WISP_SILERO_MODEL/WISP_SPIKE_MODELS/WISP_TTS_DIR/WISP_KWS_DIR/WISP_15E_TMP` 指主仓资产＋`build/models/manifest.json` 与 `third_party/model-fetch`、`third_party/spike-models` 须对 cwd 可达（voice15eRepoFile 从 cwd 上溯；本腿用 junction 零拷贝解决）；binary 从导出树 build 对 Test384/43/255/15e/41a2r1/381r2/26R3 全绿无连坐（runtime.Caller 连坐红只及 golden/manifest 全包案）。

## 10. 结案意见

**五格全翻、可带残余结案**（残余九条具名如 §9，无一阻结案）。翻勾与票面状态行改动归编排者（本腿 ⛔）。本腿写面＝本文件一枚，随 commit pathspec 落；push 归编排者。
