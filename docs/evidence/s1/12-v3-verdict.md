# 12-v3 终裁腿 —— 票 12 AC#9（registry A25）「回复后 3s 回落」非实现者裁决表（2026-11-01）

> 裁腿 `12-v3`（12-v1/12-v2 两代静默死亡判死后重派）。裁决者≠实现者（D22 双角色）。
> **读数对象＝`1b611f5a6`（12-r2 落地腿，6 枚 +547/−3），现核于 HEAD=`20392608e`**：五枚接缝件在
> HEAD 与该 commit 的差异只有 `resident_voice_engines_windows.go` +13 行（`rv.speak()` 的票 26p
> barge-in 两段，不在 openScope 钩子上）——T0/T1/T2 三戳行与 HEAD 逐字节同。**读数全部现跑**
> （2026-11-01 08:38–08:49 +08），实验全在仓外导出树 `D:\tmp\wisp\20261101\12v3\export`
> （`git archive HEAD` 解树、五件 `cmp` 对 HEAD 一致后构建）；主仓工作树零突变。跑法＝381-v3
> 同套 staging（三 DLL 与二进制同目录＋`WISP_*` 五 env＋cwd=包目录）。⚠ 在飞 `382-w-r2`/`393-r2`
> 在我读数期间落树：**读数对象是 HEAD 导出树，与脏件零关联；两腿写面零触碰**。

## 0. 一句话裁定

AC#9 三条判据在 `1b611f5a6` 上全部可判、全绿，接缝最小侵入成立，反向判据两发全闭合 ⇒
**AC#9 可翻勾**（带两条具名残余，见 §3/§7）。**「3s 回落」语义裁定＝D43#32 的 3s 静置窗**
（Settling→dispose 级联，实测 3000–3001ms）；AC#9① 的 ≤3000ms 是另一条判据（进窗腿），
在具名 staging 暖窗 1.2s 下 5/5 全过。

## 1. 接缝最小侵入核（必跑①）—— PASS

- **T0**＝`run.go` `consoleSink.Publish` 的 `EvTextDelta` 分支 1 行 `c.settle.markLastToken()`
  （票 12 原生 sink，非 30-w-r1 行）＋`agentRuntime`/`consoleSink` 各 +1 字段＋`Sink:` 字面量 1 处——diff 现读坐实。
- **T1**＝engines `openScope` 的 **Logf 既有闭包**内 `strings.HasPrefix(format, "session: SCOPE-SETTLING")`
  → `markSettling()`；锚行＝`internal/session/scope.go:596` 逐字对上（s.log 原始 format 串行首，漂移只会停戳）。
- **T2**＝同函数 **FreeOSMemory 既有闭包**内 `markFellBack()` 置于 `freeOSMemoryNow()` 之前；
  级联调用点 `scope.go:681-682`（C11 step 4，onSettleExpired＝D43#32 触发器）。
- **statemachine/session 零字节**：commit 文件清单 6 枚无一在 `internal/{statemachine,session}`（stat 现核）。
- **机器 Sink 面**：产码 `statemachine.New(` 构造点实测 **3 枚**（`models.go:358`＝票 121 `7fe5e7375` 引入、
  `resident_voice_windows.go:534/1778`），**0 枚带 Sink**——cmd/wisp 产码 `Sink:` 仅 `providers.go:194`
  （storeHealthSink，probe 面）与 `run.go:1295`（consoleSink，agent 面），均非机器 Sink ⇒ **273 投递面未碰成立**。
  ⚠ 实现件写「构造点 2 枚」，计数照录为 3（models.go:358 在其开工前已在），实质主张不受影响。
- **30-w-r1 挂接行零触碰**：`resident_windows.go` diff 纯 +8/−0（4 行注释＋4 行兄弟挂接段），
  且该文件 HEAD 与 commit 逐字节同 ⇒ 后续腿也未碰。
- **活体证据**：基线样本值非零即三戳全真落（任何一枚缺失，§4 的反向判据已证实测案必红）。

## 2. 分布复跑 ≥5 样本（必跑②，staged 导出树、产线全链、单调钟）—— PASS

| 样本 | token→Settling | Settling→回落 | token→回落 |
|---|---|---|---|
| 1 | 1229ms | 3000ms | 4229ms |
| 2 | 1226ms | 3000ms | 4226ms |
| 3 | 1227ms | 3000ms | 4228ms |
| 4 | 1229ms | 3000ms | 4230ms |
| 5 | 1219ms | 3000ms | 4219ms |

- 还原核兼复现跑（独立重建重跑）：1221/1228/1223/1235/1221 · 3000/3000/3001/3000/3000 · 4222/4228/4225/4235/4222。
- **判据① 10/10 ≤3000ms 全过**；**「3s 回落」本体 3000–3001ms＝冻结 `SettleStaticDefault` 真值**
  （计时器不早火、经 owner 队列到戳）。对照窗 1201–1206/3000–3003/4201–4206：**同形**，进窗腿高 13–29ms
  （机器负载差，非分布漂移；两轮内部自洽 ±8ms）。
- 走位行逐行在案（log）：`SCOPE-WARM reason=announce-end warm_timeout=1.2s` →
  `SCOPE-SETTLING reason=warm.idle static=3s` → `SCOPE-DISPOSED reason=settling.expired
  cascade=models-unload,panel-hide,free-os-memory,grants-expire`；D43#15/#26 产线派发同 log 可引。

## 3. ★关键裁量：AC#9 的「3s」指哪段（必跑③，12-v1 死前未裁）—— 裁定

1. **「3s 回落」＝D43#32 的 3s 静置窗（Settling→dispose 级联）**。三处互证：`PLAN.md:1430` 逐字
   （票面引 `:1424`＝第 20 处漂移，本腿现跑 `sed -n` 坐实 S1 行在 1430）＋票 12 What-to-build `:26`
   「Settle: result presentation done → 3s → DisposalScope session teardown (incl. FreeOSMemory) →
   RSS back to idle cap ≤10s」＋D43#32 行。这一段与配置无关（冻结默认），6/6 样本 3000–3001ms——
   它就是票 12 自己写的「3s Settling→unload」的实测记录，**是 AC#9 标题的本体**。
2. **AC#9① 的 ≤3000ms 判的是进窗腿**（最后 token→进 Settling）＝管道时延上限，不是 3s 本体。
   它的可判形状依赖暖窗配置：staging 暖窗 1.2s（`voicePlan.WarmTimeout`＝`[session].warm_timeout_sec`
   的读数位）下 5/5 过；**生产默认 90s**（`internal/config/schema.go:193 default:"90"`、
   `loader_test.go:72` 钉死）下按字面必超（≈90s＋管道）。**诚实边界照录，且裁定它不挡翻勾**：
   90s 保热是票 28/D43#26/#31 的引擎热窗设计、不是 S1 语义（S1 文本路无语音引擎、无保热窗，
   What-to-build 的「3s Settling→unload」正是无暖窗形状）；该腿在**任何配置下都可实测**
   （缝不依赖配置），残余具名即可。⇒ 非未定义，不停手。

## 4. 反向两发＋AST 常驻牙（必跑④）—— 全闭合

- **M1 摘 T1/T2**（导出树删两戳行，编译 rc=0）⇒ 实测案 **rc=1**，红句逐字
  `timed out waiting for the reply's fall-back to land in the trace (T2)`；**log 里 `SCOPE-DISPOSED`
  仍在**＝行为在、记录无处落（正是 AC#9③ 要防的形状）；AST 牙同变异 **rc=1（T1 句，0 秒红）**。
- **M2 摘 T0**（`markLastToken` 行换 `_ = c`，编译 rc=0）⇒ **rc=1 红句同形逐字**。
- **还原核×2**：按件从 archive 还原后五件 `cmp` 对 HEAD 全 IDENTICAL（两轮共 10/10），重建 rc=0
  重跑 **5/5 同分布**＋AST rc=0。
- **AST 常驻牙随套**：`Test12r2SeamCallSitesAreTheMeasurement` 与五样本案同一二进制同套跑
  （staged/restored 两轮 PASS）；对两枚变异体 0 秒红。

## 5. 零回归（必跑⑤）—— PASS

- **六族合跑 rc=0，28/28 PASS**（Test43 9＋15e 6＋41a2r1 5＋381r2 5＋26R3 3；含 36.4s 真设备案），
  **`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms` 逐位同值**（连
  seamWaitToFirstText=5.51s finalLatency=4.096s speechMS=3392 closedBy=end 逐字同）。
- **Test384 单发 rc=0，3/3 PASS**（E2E 1.15s）。
- **邻接六包 rc=0 全 ok**：session 6.277s／statemachine 0.326s／observe 6.186s／config 4.587s／
  watchdog 0.423s／proc 5.461s。
- 本腿自证：导出树构建 rc=0×4（基线/mut1/mut2/restored）＝编译面干净；`gofmt -l cmd/wisp/` 空。
- staging 协议备注：六族案的 `build\models\manifest.json` 由 cwd 向上寻仓命中（主仓包目录跑、
  只读 fixture，模型装进 `WISP_15E_TMP` 临时 case）——实现腿 zeroreg log 同源可引（其 manifest
  逐字 `D:\work\...\Wisp\build\models\manifest.json`）。首次误用导出树包目录作 cwd 得 6 枚
  setup 红（找不到 manifest），换 cwd 后全绿——**非回归，协议差**，照录防后人踩。

## 6. 落点判定（必跑⑥）—— PASS

AC#9 字面是「有实测记录」，**不点名** `docs/SLO.md`（同票 AC#2「PASS recorded into docs/SLO.md
appendix」点名＝对比）⇒ 记录落证据表＋票面引用，往 SLO.md 加附录才是自创章程位。
现核：`git status -- docs/SLO.md` 干净、commit 6 枚清单不含它 ⇒ **SLO.md 未动成立**。

## 7. AC#9 可否翻＋票 12 收口路径（必跑⑦）—— 裁定

- **AC#9 可翻勾**（§1–§6 全过；两条具名残余＝§3 边界＋文本路仍无回落走位，见实现件 §8.2、本腿维持）。
- AC#2：对抗验收已在案（票面 `:49` 12-slo-v1 verdict 七项全过、解除条件＝票 66 done 7/7 成立）
  ⇒ 只差编排者引该 verdict 翻勾，无需新测量。
- AC#8＝机主目验：**纯等人、无 agent 替代**（12-r1 §2 判定维持：机审只能备料不能替代签收）。
- **收口路径（15/28 先例判）**：15-done 现核即 **3/6 未勾结案**；28 收口句逐字「可翻 5 格、AC3 归票 59
  ⇒ 带具名残余结案（★先例：15-done 自己就是 3 格未勾结案）」⇒ **票 12 可以带 AC#8 未勾收口**，
  条件＝AC#2/AC#9 先行翻勾＋AC#8 以具名等人残余入账（机主答 R15 三项 → 票 68 翻默认复测 →
  默认构建对 7 态截图一次性签收，同清票 68 目验面）。
- 时序：本表交付 → 编排者翻 AC#9（引本表）＋AC#2（引 12-slo-v1 verdict）→ 票 12 带具名残余转
  done 并改 `-done` → **票 16 解锁**（12-r1 §4：16 的唯一剩余阻塞＝票 12）。

## 8. 未做／边界（照录）

1. 未跑 d22scan（派单明令 ⛔ ban#8 扫 docs/；仪器面是实现腿职责且其 clean 读数在案）。
2. 进窗腿默认 90s 形状未实测（≈90s＋管道由 `schema.go:193` 与冻结面推得；缝可测，今日不跑
   90s×5 样本＝时间成本无增量价值）。
3. 主仓在飞脏件（382-w-r2：`resident_voice_engines_windows.go`/`resident_voice_windows.go` M＋
   一 ??；393-r2：internal/audio/aec 四 M＋一 ??）零触碰、不进读数；382-w-r2 恰动 T1/T2 所在文件，
   **其内容不属本裁**（读数对象＝HEAD 导出树）。
4. 停手条件自核：三处均未触发（PLAN:1430 语义对上；接缝未逼 statemachine/session；与 30-w-r1
   装配无冲突——挂接段是其兄弟段）。

临时件（只建不删）：`D:\tmp\wisp\20261101\12v3\{head.tar,export,bin,logs,cmp,tmp15e}`；
logs＝`12v3-{staged,mut1,mut2,restored,sixfamilies,384-single,adjacent6}.log`。
