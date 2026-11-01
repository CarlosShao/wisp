# 票 382 残余 R382-1 落地腿 `382-w-r2` —— L2 接线（写腿）

日期 2026-11-01 上午 · 腿 `382-w-r2`（写、`cmd/wisp/**` 独占＋本件；重派——382-w-r1 过夜静默死亡判死）·
起手 HEAD `20392608e`（`dev`）。**开工三核过**：①HEAD 在 dev＝`20392608e`（≥ 派单锚点）；②索引无 staged 脏
（工作树脏区全 `.scratch/**` 与历史 untracked 临时件）；③`git status --porcelain -- cmd/wisp`＝**空**。
环境：staging 全套（no-tag＋`-tags webrtcap` 两枚 staged 测试二进制＋三 DLL 同目录＋五 `WISP_*` env＋
`-test.timeout 30m`＋cwd＝真仓 `cmd/wisp`，runner `D:\tmp\wisp\20261101\382wr2\run.sh`；Git Bash 跑，pwsh 拆
`-test.v` 的坑未踩）；gofumpt 绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`（v0.12.0）；rc 一律重定向现取；
每跑 `unset GOTMPDIR`；突变全走仓外 `-overlay`＋还原核；⛔ push。
★**共享树在飞申明**：本程进行中并发腿 **393-r2** 在 `internal/audio/aec/**` 落了未提交改动（4 改＋
`aec393r2_reset_windows_test.go`＋`393r2-flush-reset.md`，mtime 08:46–09:09，本程 08:28 开工后出现）——
**非本腿写面，全程零触碰**；本程两枚 staged 二进制按共享树现状编译（含其在飞面），零回归读数为两树合并态
的读数；本腿提交 pathspec 逐枚排除其 6 件。

## 0. 一句话

**AC#1①②（382-v1 修法逐字）落地**：`residentVoice` 加 `ttsGate *risk.TTSGate` 字段＋`attachProvenance()`
初始化＋speak() 的 HOP 1 由 `prov.CheckText` 换 `ttsGate.CheckTextWithInheritance`（refused ⇒ 拒播，三态
source 进问题行）——AC#2 三读数全绿（**形 B 污句拒播 refused=1 spoken=0 slot_loads=0＋问题行带
"from closed task scope"；干净句照播 spoken=1 samples=75431；直查作用域仍拒**）；AC#3 反向两发全红
（M1 摘闸调用 ⇒ 污句照播 56482 样本；M2 摘初始化 ⇒ 闸对象消失〔nil panic 形〕＋直查污句照播 42635
样本）＋还原核 3/3 绿逐位复现；AC#4 六族 staged 合跑 **34/1/0 逐位复现基线**（唯一红＝
`Test384ResidentWatchdogEndToEnd` 合跑时窗既有形状）＋**`262/640/902` 逐位同**＋26p 套件 (a)(b)（(a) 单跑
26/83/50/50/34ms 全 ≤400；(b) episodes=1 triggers=0 finals_delta=0＝392 基线逐字）；d22scan clean rc=0＋
正控 selftest 35/0/0（RUN=77）；**AC#1③（任务级积累钩）＝停手条款②具名触发**（§5：生产挂点不存在，
产出=证据链＋残余，⛔ 造假挂点）。

## 1. AC#1①② 接线（382-v1 修法逐字，26p-r1 后的最新形状）

26p-r1 改过 speak() 后的现核：修法四项全部对得上——旧 `prov.CheckText` 仍在 HOP 1（行号 864→898 漂移，
形状逐字同 382-v1 判语 §6 所引），停手条款①**未触发**。落点：

1. `cmd/wisp/resident_voice_windows.go:361`＝`ttsGate *risk.TTSGate` 字段（挨 `prov`，注释钉死唯一写者＝
   `attachProvenance`，与 prov 同生同死）；
2. `cmd/wisp/resident_voice_windows.go:615`＝`rv.ttsGate = risk.NewTTSGate(p, risk.DefaultTTSGateConfig())`
   ——借**同一个**引擎（bridge 往里 mark 的那枚，`resident_windows.go:579` 交来的 `src.run.c25`），私有
   引擎＝"check against marks production never made"那形，注释拒之；`DefaultTTSGateConfig`＝票姿势（形 B
   开、L2 确认开）；`attachProvenance` 的早退分支（作用域打不开）prov 与 gate 一齐留 nil ⇒ speak() 跳检
   ＝与今日 prov==nil 同形；
3. `cmd/wisp/resident_voice_engines_windows.go:904`＝HOP 1 换
   `_, refused, source := rv.ttsGate.CheckTextWithInheritance(req.scope, risk.ChTTS, text)`；**refused ⇒
   ttsGated+1＋noteProblem（带三态 source）＋return，绝不合成**（票 26 AC7 "before playback" 不变）；
   注释钉死：直查保旧 verdict、继承扫描加"已闭任务作用域"拒播、无 provenance 的 boot 闸与引擎一齐缺席。
   闸 nil 肢＝本 boot 没接 provenance（与旧 prov==nil 肢逐字节同义）。

⛔ 未动：`internal/{risk,speech,tools}/**`、D43 表、statemachine、L2 弹卡 UI（NativePanelBridge 交互仍归
票 26 集成面残余——本腿只接机制调用＋读数）、thresholds/golden/testdata 仓根/frontend/design/scripts/
.github/票面/台账/specs/PLAN。

## 2. AC#2 读数（`-test.run Test382WR2`，no-tag staged 二进制，log 05/13）

| 读数 | 判据 | 实测 |
|---|---|---|
| 形 B 污句拒播（`Test382WR2TTSInheritanceRefusesClosedTaskTaint`） | 开任务作用域→Mark→**RecordTaintInheritance**→scope.Close→播含毒句 ⇒ 拒 | `INHERIT refused=1 passed=0 spoken=0 samples=0 slot_loads=0`＋`risk/C25: recorded taint inheritance link for task task-382wr2-closed (1 taints)`＋问题行含 ChTTS 与 **"from closed task scope"**；**合成零次**（slot Loads=0＝闸在合成前跑） |
| 干净句照播（`Test382WR2TTSCleanSentencePlaysThroughTheGate`） | 闸带一条不相关继承链（更狠的正控）⇒ 干净句照过 | `CLEAN passed=1 refused=0 spoken=1 samples=75431`（真 matcha 合成 5.07s） |
| 直查仍拒（`Test382WR2TTSOpenScopeDirectCheckStillRefuses`） | 活作用域毒经新闸仍拒＝换接口不丢直查半边 | `DIRECT refused=1 passed=0 spoken=0 samples=0` |

先例对齐：R26-5 形状（调用点收、TestWireR1AnnounceGateRefusesTaintedText 同款 rig）。本件三测试即
RecordTaintInheritance 的**调用方读数证明**（与 382-v1 仓外 e2e 同形，落在本腿写面内常跑）。

## 3. AC#3 反向判据（仓外 `-overlay` 两发＋还原核；log 09-13）

| 突变 | 摘什么 | 红相（读数） |
|---|---|---|
| **M1** | speak() 的闸调用（`CheckTextWithInheritance` 调用行 ⇒ 恒 no-refusal） | `rc=1`：`timed out waiting for the gate to refuse the closed-task sentence`＋**`announcement spoken runes=26 chunks=1 samples=56482`＝污句照播** |
| **M2** | `attachProvenance` 的 `rv.ttsGate = risk.NewTTSGate(...)` 初始化 | ①形 B 测试 `rc=2`：测试自身 record 调用 **nil receiver panic＝闸对象整个不存在**（挂点被摘的形状）；②直查测试 `rc=1`：**`announcement spoken samples=42635`＝污句照播**＋等拒超时 |
| **还原核** | 无 overlay 的 pristine 二进制（变异前构建、全程未动）复跑三读数 | `rc=0` **3/3 PASS 逐位复现**：INHERIT refused=1 / CLEAN spoken=1 samples=75363 / DIRECT refused=1 |

突变件全在仓外（`D:\tmp\wisp\20261101\382wr2\mut\**`＋overlay JSON 两枚），仓内 `git status` 对 internal/
零新增。

## 4. AC#4 零回归（六族＋262/640/902＋26p (a)(b)＋门禁四数两尺）

- **六族 staged 合跑**（no-tag＋三 DLL＋五 env＋cwd=包目录，`-test.run
  "Test43|Test15e|Test41a2r1|Test381r2|Test390r1|Test26R3|Test384"`）：**34 PASS / 1 FAIL / 0 SKIP**——对
  392-r1 基线**逐位复现**；唯一红＝`Test384ResidentWatchdogEndToEnd`（合跑负载时窗，A1329 既有形状）。
- **`262/640/902` 逐位同**：`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms`（Test15e 六族内）。
- **26p 套件 (a)(b)**（tagged `wisp26p.test.exe`，aec3r1 处方 CGO env 构建 rc=0）：
  - (b) `Test26PEchoNeverTriggersAndNeverTranscribes` PASS（22.75s；`ECHO-ZERO episodes=1
    all_below_cut=true triggers=0 finals_delta=0 handed_delta=0 windows_untrusted=6`＝392 基线逐字）；
  - (a) `Test26PBargeinStopsWithin400ms`：与 (b) 同二进制合跑首遍第 3 轮 **671.76ms＞400ms 红**；单跑复测
    **PASS**，五轮 **26/83/50/50/34ms 全 ≤400**。归因＝(b) 的采捕 rig 刚释放设备后的负载形状（A1329
    「合跑时窗既有形状、单跑 PASS」同族；392-r1 的 (c) 两跑一过一死同族），⛔ 计入本票；两发读数都在案
    （log 06/07）。
- **门禁四数**：d22scan（Git Bash `scripts/d22scan.sh`，一发含正控＋实扫）＝**clean rc=0**；分母现量
  **bans#1-5 internal/=287、cmd/=52（对派单起值零变）**；**ban#8 internal/=730、cmd/=182**（派单起
  729/181：cmd +1＝**本腿新测试件**；internal 730＝HEAD 既有〔392-r1 证据同值 730，派单 729 一枚过期〕，
  本腿 internal 零新增）。正控 selftest＝**PASS=35 FAIL=0 SKIP=0（=== RUN=77）**＝392 基线同值。
- **两尺**：gofumpt `-l`（本腿三枚件逐名）＝**0 行 rc=0**；`go vet ./cmd/wisp`（无 tag）rc=0＋
  `go vet -tags webrtcap ./cmd/wisp` rc=0（⚠ 第一发 tagged vet rc=1＝我只给 `-I` 漏 define 组的**尺架设错**，
  带全 define 复跑 rc=0，log 15 两发都留盘）；`go build ./cmd/wisp` 两形状 rc=0。
- 本腿三测试（05/13 两遍）全绿，且在 M1/M2 两发反向里按预期翻红（牙是活的）。

## 5. ★停手条款②具名：AC#1③ 任务级积累钩——生产挂点不存在

派单要求"任务完成 ⇒ RecordTaintInheritance（对照形 B 调用方语义）"。**结构性挂点找得到**（任务完成边界＝
`cmd/wisp/run.go` `admitTask` 的收尾闭包 ：1185-1190，那发 revoke 在 TEXT loop 的 defer 上跑）；**但
RecordTaintInheritance 的调用方语义（`internal/risk/tts_gate.go:245` 逐字："requires the caller to
accumulate normalized content at mark time"）在这条边界上无从成立**：

1. 生产唯一 mark 点＝`internal/tools/bridge.go:687` `MarkWithHostPath(dec.TaskID, dec.Tool, origin, res.Text,
   hostPath)`——毒内容（`res.Text`）只活在 internal/tools 内部，从不跨界；
2. `tools.Options`（`internal/tools/bridge.go:32`）公开面上**没有任何 mark 观察钩子**（OnUpdate＝流式进度
   增量、OnDecision＝执行前裁决，都不带 mark 内容）；
3. `risk.Provenance` 公开面无内容读出口——`ScopeTaints`（`internal/risk/provenance.go:654`）按安全设计只回
   元数据（Tool/Origin/RuneLen/MarkedAt，注释逐字 "never the content itself"）；
4. ⇒ 在 cmd/wisp 侧能发出的 RecordTaintInheritance 只有两种形：空内容（闸守卫直接吞掉＝**假接线**）或拿
   Origin 冒充内容（改动闸语义＝**按自己的判断填**）。两种都⛔。

**判**：本腿写面（`cmd/wisp/**`）内该钩不可落地；真解是 internal/tools 加 mark 观察缝（如 Options 新钩把
mark 时归一化内容喂给闸）或闸/调用方下沉 internal/tools——两者都越本腿写面且属新契约面，**摆编排者开票**。
残余如实记账，不造假挂点。本程交付的机制面（闸调用＋三态拒播）已全接，生产毒源今日照旧走直查半边
（voice-announce 活作用域），形 B 扫描随链接数增长自动生效。

## 6. 写面申报

仓内（commit pathspec 逐枚）：
`cmd/wisp/resident_voice_windows.go`（M：字段＋attachProvenance 初始化）·
`cmd/wisp/resident_voice_engines_windows.go`（M：speak() HOP 1 换 CheckTextWithInheritance＋头注一行）·
`cmd/wisp/resident_voice_ttsinherit_382wr2_windows_test.go`（新：三读数测试）· 本件。
仓外 `D:\tmp\wisp\20261101\382wr2\**`（bin-notag/bin-tag/mut/logs/cases/run.sh，只建不删）。

⛔ 越界核：`internal/{risk,speech,tools}` 零触碰；D43 表/statemachine/thresholds/golden/testdata 仓根/
frontend/design/scripts/.github/票面/台账/specs/PLAN 零触碰；L2 弹卡 UI（NativePanelBridge）零触碰；
26p-r1 的 bargein/Path-C 产码一行未动；⛔ push（推送归编排者）。共享树在飞 393-r2 的 6 件（§0 申明）
逐枚排除在本腿 pathspec 外。

## 7. 残余

1. **R382-1③（本票残余的残余）**：生产 RecordTaintInheritance 调用方＝internal/tools 新缝或闸下沉（§5，
   摆编排者）；测试面调用方已常跑（§2）。
2. 26p (a) 合跑负载形状一红（§4）＝26p-r1/392-r1 已在账的时窗敏感族，处置归验收腿。
3. 本腿读数取于共享树＋393-r2 在飞态（§0）——验收腿复跑时若 393-r2 已落库，六族/26p 读数应在同带。
