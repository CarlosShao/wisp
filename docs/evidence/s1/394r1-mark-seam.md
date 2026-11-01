# 票 394 落地腿 `394-r1` —— risk 层 mark 观察缝（第三形，机主 A1362 批）

日期 2026-11-01 下午 · 腿 `394-r1`（写腿；写面＝`internal/risk/**` 新缝＋`cmd/wisp/**` 一枚新测试件＋本件）·
起手 HEAD `608693adc`（`dev`）。**开工三核过**：①HEAD 在 dev＝`608693adc`（≥ 派单锚点）；②索引无 staged 脏
（`git diff --cached` 空）；③`git status --porcelain -- internal/risk cmd/wisp docs/evidence/s1`＝**空**。
环境：本地全套（staged 语音案按 382-w-r2 runner 配方：no-tag 测试二进制＋三 DLL 同目录＋五 `WISP_*` env＋
`-test.timeout 30m`＋cwd＝真仓 `cmd/wisp`，runner `D:\tmp\wisp\20261101\394r1\run.sh`；**Git Bash 跑测试**
——pwsh 拆 `-test.v` 的坑本程也踩了半步即按卡改道）；gofumpt 绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`
（`-l` 六枚件 0 行 rc=0）；rc 一律重定向现取；每跑 `unset GOTMPDIR`；突变全走仓外 `-overlay`＋还原核；
⛔ push。

★**共享树在飞申明**：本程进行中工作树有他腿未提交改动（`internal/panel/{pump.go,subagent_roster_197.go,
backpressure_bound_35r8_test.go}`、`internal/winsec/**` 五件、`frontend/src/lib/panel.ts`）——**非本腿写面，
全程零触碰**；本程 staged 二进制按共享树现状编译（含其在飞面），零回归读数为两树合并态的读数；本腿提交
pathspec 逐枚排除其件。

## 0. 一句话

**第三形落地**：`internal/risk/markseam_394.go` 新缝（未导出观察者接口＋`Provenance` 注册/通知路由＋
`TTSGate.observeTaintMark` upsert）＋`MarkWithHostPath` 落账后一处通知＋`NewTTSGate` 借引擎自注册——
**RecordTaintInheritance 的生产调用方落地**：生产 mark 内容（非手选常量）经缝进继承链，`normalizeTaint`
义务收回 risk 内履行（382-w-r2 停手②就此闭合）；AC#2 生产读数 **`SEAM-PRODUCTION refused=1 passed=0
spoken=0 samples=0 slot_loads=0 links=1`**（纯 `prov.Mark` 喂链、零手动 record、全角原文 vs 半角播报——
只有缝内归一化副本才可能匹配）；AC#3 两发突变（M1 摘汇流口通知／M2 摘自注册）risk 层 5 红＋staged 层
**喂数零链＋污句照播（runes=27 samples=81534）**＋382WR2 三案在突变下保持绿（缝独立性）＋还原核 5/5
绿；AC#4 risk 全套 rc=0＋六族 **34/1/0 逐位复现基线**（唯一红＝Test384 合跑时窗既有形状）＋
`LATENCY 262/640/902` 逐位同＋race rc=0＋d22scan clean rc=0＋正控 35/0/0（RUN=77）；三条停手条款
逐一核过未触发。

## 1. 选型与落地（AC#1）

**两亚形拍板＝TTSGate 借引擎自注册（亚形 B），不取 ProvOptions 钩（亚形 A）**。理由：
1. **接线顺序**：`ProvOptions` 在 `NewProvenance` 时定格，而生产闸在 `attachProvenance`
   （resident_voice_windows.go:615）才存在——亚形 A 必须在组合根加晚绑定管道（cmd/wisp 写面扩获）；
   亚形 B 的 `NewTTSGate(prov, cfg)` 天然持有引擎，382-w-r2 已落的接线**零改动即生效**。
2. **契约面最小**：B 的全部新面（接口、注册、通知）都可包内未导出——契约面只在 risk 的要求逐字达成；
   A 要动公开配置面。
3. **注册范围天然正确**：闸只观察它借的那台引擎（"same engine the bridge marks into" 规则的延续），
   私有引擎的 mark 永远不会串进别家闸。

**落点四件**（全部新增面，既有契约逐字未动）：
1. `internal/risk/markseam_394.go`（新）：`taintMarkObserver` 未导出接口；`Provenance.addMarkObserver/
   notifyTaintMark`（观察者切片锁内拷贝、锁外调用——观察者拿自己的锁永不与引擎锁倒序；零观察者时一次
   len() 读零分配）；`TTSGate.observeTaintMark`（upsert）。
2. `internal/risk/provenance.go`：`Provenance` 加 `markObs` 字段；`MarkWithHostPath` 落账（`p.mu.Unlock()`
   之后）通知 `string(rp)`——**已归一化＋已截断**的索引坐标内容，早退分支（归一化后为空）不通知。
   mark 语义（返回值/索引/豁免/预算）一字未动。
3. `internal/risk/tts_gate.go`：`NewTTSGate` 在 `EnableTaintInheritance` 开时自注册（关时链永不被扫，
   注册纯属浪费）；`CheckText*`/`Scope.Close`/`RecordTaintInheritance` 契约逐字未动。
4. `internal/risk/markseam_394_test.go`（新）＋ `cmd/wisp/resident_voice_markseam_394_windows_test.go`（新）。

**留存姿势（已批代价的落地形）＝每作用域一条 link 的 upsert**，不逐 mark 追加：已落扫描契约
`checkAncestorScopes` 只扫**前** `MaxAncestorDepth`(16) 条链，逐 mark 追加会让第 17 个 mark 起的生产毒
**永久**落在扫描窗外（保护性失聪）；upsert 让链增长＝O(不同已 mark 作用域数)，一个任务的多次 mark 折进
同一条链（`TaintContent`/`ToolSources` 平行追加、`MarkedAt` 更新为末次）——这恰是 382-w-r1 文档里
"accumulatedContents" 的复数形状。**382-w-v2「只收任务作用域」建议的处理**：risk 内无任务作用域判别器
（生产 id 混形：环路自铸、`host:*`、`run:*`），按前缀过滤就是在新缝里造脆弱契约；upsert 已达成该建议
追逐的留存界，且不过滤是 fail-closed 姿势（永不闭合的作用域如 `run:project-instructions` 的 fs.read
毒，链是它们到达播报裁决的**唯一**通路——直查只看 voice-announce 作用域）。

**导出的 `RecordTaintInheritance` 保持逐字节不动**（append 语义由 `Test394SeamUpsertsOneLinkPerScope`
尾段钉住：显式记录永远追加、绝不并进缝的 link）；缝的生产路径走自己的 upsert 面并共享同一
`taintInheritanceLink` 模型——「生产调用方落地」的兑现形在证据 §2/§3。

**4 处 mark 点覆盖核**（AC#2 的现场判定）：

| mark 点 | 引擎 | 经缝喂闸 |
|---|---|---|
| `internal/tools/bridge.go:687`（MarkWithHostPath，任务作用域） | 共享 `rt.c25` | ✅ 直喂（闸借同一引擎） |
| `internal/tools/subagent_197.go:512`（Mark，父任务作用域 task.output） | 共享（run.go:896/917） | ✅ 直喂 |
| `internal/projctx/projctx.go:349`（Mark，`run:project-instructions`） | 共享（run.go:1374） | ✅ 直喂 |
| `cmd/wisp/panel_assets.go:245`（Mark，`panel-assets-l2`） | **私有引擎**（`detector()` 现建，活死在本进程） | 缝照常触发（汇流口在 MarkWithHostPath），但该进程**无闸可喂** |

第 4 处处置具名：`wisp panel-assets -l2` 是一次性 CLI 诊断命令，其引擎「lives and dies inside this
process」（该文件注释原文），进程里没有播报路径也没有闸——硬造一个被消费方为零的闸就是 382-w-r2 停手②
拒绝的同款**假接线**。**判：停手条款①未触发**——两条缝的观察面（汇流口）对第 4 处照常够得着（同一次
`prov.Mark` 调用、同样的归一化通知），缺的是那个进程里的消费方而那个进程本无播报；本腿不对一次性诊断
命令造消费方。

**归一化义务收回的钉子**：单元侧 `Test394SeamRecordsNormalizedContentAtMarkTime` 白盒断言链内内容
== `normalizeTaint(raw)`；生产侧 AC#2 测试用**全角原文 mark、半角文句播报**——raw 对 raw 零共享窗，只有
缝存了归一化副本才可能拒。382-w-r2 测试被迫手选已归一化常量的痛点，两根钉子共同宣告闭合。

## 2. AC#2 生产真挂读数（`-test.run "Test394R1|Test382WR2"`，log 03/05/09）

| 读数 | 判据 | 实测 |
|---|---|---|
| 喂数（`Test394R1SeamRecordsProductionMark`） | 生产 mark 一发、零手动 record ⇒ 链上恰 1 条 | PASS（links=1） |
| 拒播（`Test394R1SeamFedRefusalRefusesClosedTaskPoison`） | mark→Close→播含毒句 ⇒ 合成前拒 | `SEAM-PRODUCTION refused=1 passed=0 spoken=0 samples=0 slot_loads=0 links=1`＋问题行含 ChTTS 与 "from closed task scope" |
| 382WR2 三案不退化（对照 382-w-r2 §2） | 机制换喂源后原三读数逐位保持 | `INHERIT refused=1`／`CLEAN passed=1 spoken=1 samples=75326~75473`（真 matcha 合成）／`DIRECT refused=1` |

对照 382-w-r2 的 INHERIT refused=1 形：同形状、同拒播点（合成前 slot_loads=0），**喂的是生产 mark 内容
（`prov.Mark` 原文）而非手选常量**，且其显式 record 调用在缝在场时变成冗余而非依赖（§3 突变下独立绿）。

## 3. AC#3 反向判据（仓外 `-overlay` 两发＋还原核；log 06-08/10-12）

| 突变 | 摘什么 | 红相 |
|---|---|---|
| **M1** | `MarkWithHostPath` 的缝通知（汇流口聋） | risk 层 rc=1：4 个缝测试＋`TestTTSGate_AncestorInheritance`（编码新姿势的计数）红，独立案（NotRecorded/Disabled/其余 TTSGate）绿；staged 层 rc=1：喂数案 `the seam recorded 0 links`＋拒播案 `timed out waiting for the seam-fed gate to refuse`＋**`announcement spoken runes=27 samples=81534`＝污句真合成上了设备** |
| **M2** | `NewTTSGate` 的自注册（闸不订阅） | risk 层红相与 M1 同形；staged 层红相与 M1 同形（0 links＋拒播超时） |
| **缝独立性对照** | 同两发突变下跑 382WR2 | `INHERIT refused=1`／`CLEAN`／`DIRECT` 三案**全绿**——显式 record 通路与缝互不依赖，382-w-r2 的机制在缝被摘时零退化 |
| **还原核** | 无 overlay 的 pristine 二进制（突变后复跑） | rc=0 **5/5 PASS**：喂数 links=1＋`SEAM-PRODUCTION refused=1 … slot_loads=0`＋INHERIT/CLEAN/DIRECT 逐位复现 |

突变件全在仓外（`D:\tmp\wisp\20261101\394r1\mut\**` 两份副本＋overlay JSON 两枚），仓内 `git status`
对 internal/ 零新增脏（仅本腿写面六件，§5）。

## 4. AC#4 零回归＋门禁四数两尺

- **risk 全套**：`go test ./internal/risk` **rc=0**（log 02；含 382-w-r1 的闸读数全族 tts_gate_test.go——
  唯一既有断言改动＝`TestTTSGate_AncestorInheritance` 链数 1→2，即已批留存姿势变更的可见面，注释钉明）。
- **六族 staged 合跑**（`Test43|Test15e|Test41a2r1|Test381r2|Test390r1|Test26R3|Test384`）：**34 PASS /
  1 FAIL / 0 SKIP**——对 382-w-r2 基线**逐位复现**；唯一红＝`Test384ResidentWatchdogEndToEnd`（合跑负载
  时窗，A1329 既有形状）。
- **`262/640/902` 逐位同**：`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms`（Test15e 六族内）。
- **race 尺**：`go test -race ./internal/risk -run Test394` **rc=0**（缝并发面干净；upsert 与通知的锁纪律
  在并发测试＋race 双覆盖）。
- **两尺**：gofumpt `-l`（本腿六枚件逐名）＝**0 行 rc=0**；`go vet ./internal/risk` rc=0＋`go vet
  ./cmd/wisp` rc=0＋`go build ./cmd/wisp` rc=0。tagged（webrtcap）形状未跑——本腿零触碰音频面
  （`internal/{session,speech,audio}`、L2 弹卡 UI、26p 测试面），tagged 尺架设归 26p-k1 在飞件，不重。
- **门禁四数**：d22scan（Git Bash `scripts/d22scan.sh`，一发含正控＋实扫）＝**clean rc=0**；正控
  selftest＝**PASS=35 FAIL=0 SKIP=0（=== RUN=77）**＝基线同值。分母现量与归因：
  **bans#1-5 internal/=288**（382-w-r2 起值 287：+1＝本腿 `markseam_394.go` 唯一新增生产件）；
  **bans#1-5 cmd/=52**（零变）；**ban#8 internal/=734**（382-w-r2 起点树上 729 枚 tracked＋两腿之间
  他人落库 3 件〔393-r2/392-r2/398-r1 测试件〕＋本腿 2 件＝734，`git ls-tree` 逐点核过）；
  **ban#8 cmd/=184**（382-w-r2 记 730/182 均含其自件；树上 182＋他人 16r1 一件＋本腿一枚测试件＝184）。

## 5. 写面申报

仓内（commit pathspec 逐枚）：
`internal/risk/markseam_394.go`（新：缝本体）·
`internal/risk/markseam_394_test.go`（新：六读数单元）·
`internal/risk/provenance.go`（M：markObs 字段＋MarkWithHostPath 落账后通知）·
`internal/risk/tts_gate.go`（M：NewTTSGate 借引擎自注册＋头注）·
`internal/risk/tts_gate_test.go`（M：AncestorInheritance 链数 1→2＋姿势注释）·
`cmd/wisp/resident_voice_markseam_394_windows_test.go`（新：生产接线两读数，复用 382wr2 腿，零碰其文件）·
本件。仓外 `D:\tmp\wisp\20261101\394r1\**`（bin-notag/mut/logs/cases/run.sh，只建不删）。

⛔ 越界核：`internal/tools/**` 零触碰（第三形的全部意义）；`internal/{session,speech,audio}/**`、
D43 表、statemachine、thresholds/golden、testdata 仓根、frontend/design、票面、台账、specs、PLAN、
docs/SLO.md 零触碰；26p-k1 的测试面（(a) 案）零触碰；既有 `CheckText*`/`Scope.Close`/
`RecordTaintInheritance` 契约零字节改动；⛔ push（推送归编排者）。

## 6. 三条停手条款核

1. **两亚形都无法覆盖 4 处 mark 点** ⇒ **未触发**：汇流口（MarkWithHostPath）是 4 处的公共 funnel，
   观察面对 4 处全部够得着；第 4 处（panel-assets CLI）缺的是其进程内的消费方（本无播报、本无闸），
   处置＝照常触发不硬造消费方（§1 具名），非够不着。
2. **缝必须动 tools 契约** ⇒ **未触发**：4 处 mark 点 0 字节触碰，缝完全在 risk 包内自洽。
3. **normalizeTaint 无法在缝内履行** ⇒ **未触发**：通知点在 `MarkWithHostPath` 内部、归一化之后，
   义务在 risk 内闭合；white-box 断言＋全角/半角两钉共同证明。

## 7. 残余

1. **注册前 mark 不入链**：缝从 `NewTTSGate` 注册起观察；boot 早期（闸建立前）的 mark（如启动期
   projctx 盖戳）不追溯——其毒仍在各自作用域内活着，与今日保护面相同，无回退；追溯需要内容读出口
   （`ScopeTaints` 按安全设计只回元数据），不动。
2. **缝日志走 `logf`**（宿主注入汇，nil=drop）：与 mark 路径既有日志惯例一致；staged 测试里缝日志
   不可见是路由使然，链数断言是缝触发的判据。
3. **`inherited` 链无淘汰**：upsert 后链长＝不同已 mark 作用域数，超长会话（>16 个不同作用域）触及
   382-w-r1 既有 `MaxAncestorDepth` 首-16 扫描帽——既有防御参数的既有性质，本腿不动 CheckText 语义；
   如需演化归 382 契约面另开票。
4. 26p (a)(b) 不在本腿 AC#4 清单（派单六族＋risk 全套口径）；音频面零触碰。
