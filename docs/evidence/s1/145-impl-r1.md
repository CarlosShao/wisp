# 145-impl-r1 — §2.14 终态族（OPTION A＝Event.Stop）＋ ToolRunRow.Level ＋ AC#2b 两半枚：实现程

- 日期：2026-11-09 13:1x–14:1x（Asia/Shanghai）
- 腿：票 145 实现写腿（A 腿＝本波唯一 cmd/wisp 写腿）；写面＝`internal/panel/**`＋`internal/agent/loop.go`（＋sink.go 注释）＋`cmd/wisp/**`＋`internal/config/**`＋`internal/llm/**`
- 起手锚：`71ff0fc87`（HEAD，dev）；码落地 commit＝`4111ada61`（只 commit 未 push）
- 判定程前置＝`docs/evidence/s1/145-ac2-r1.md`（同日 12:4x）：其 §1 表的「不落」三族（终态族／ToolRunRow.Level／AC#2b 两半枚）**就是本程的三件活**，其 §3 表预判「三面同开一轮收齐」——本程照此兑现

## 0. 起手三核（现跑，不信转述）

1. **派单点名的串/行还在不在**：现量 `internal/agent/loop.go`——`EvDone` 发布点 `:1111-1114` 只带 Text ✓（未过期）；`Event.Stop`（`sink.go:65`）唯一填点＝`EvError` `:553` ✓；`memory.ToolCall.RiskLevel`（`internal/memory/models.go:80`）✓；bookedArgs 喂线行 `cmd/wisp/run.go:2364-2379`（旧名 bookedArgs）✓。**判定程预判的全部前提现量成立，无过期项**（与前两腿「票面前提已过期」不同形）。
2. **引文出处核**：census §2.14（`145-snapshot-field-census-r1.md:398-411`）现读原文——「给 EvDone 补 Stop 还是走 Result.Status＝取向交编排者」；编排者已裁定 **OPTION A＝Event.Stop**（派单原文），本程照做不再问。
3. **红/绿归属**：全部门禁一律 `git clone -c core.longpaths=true --shared` 导出树跑（`.scratch/wisp/probes/145/impl-r1-clone/`，锚＝本程 commit `4111ada61`；基线对照树 `impl-r1-base-clone/` 检出父提交 `71ff0fc87`）。主树在飞件（323 两枚未跟踪测试件等）零影响。⚠ 环境注：克隆树无 `build/`、`third_party/sherpa-onnx/`（两目录 gitignored），cmd/wisp 测试二进制首跑 `exit status 0xc0000135`（DLL 缺）——补主树两目录绝对路径进 PATH 后复跑即真读数（两次克隆树均同法）。

## 1. 落地清单（按判定程 §1 表三族顺次，file:line＝commit `4111ada61` 现量）

### 1a. §2.14 终态族（OPTION A＝Event.Stop，判定程 §1 行 13 的锁全开）

| 落点 | 内容 |
|---|---|
| `internal/agent/loop.go:1092` | `finish` 增末参 `stop llm.StopReason`；`EvDone` 发布（`:1128-1134`）带 `Stop: stop`——**全终态一个发布点收口** |
| `loop.go` 八处 finish 调用点 | 完成 `:585`（turn.Stop 逐字）／截断 `:575`（turn.Stop）／取消三处 `:452,:514,:596`（`llm.StopCancelled`）／流失败 `failWith`（turn.Stop，适配器真发 StopError）／**loop 自有终态（C22 刹车 `:1080`、预算拒 `:447`、buildRequest 失败 `:499`）传空串**——没有 provider 词的终态不借词，空＝诚实 |
| `loop.go:1070-1074` | `failWith` 同增 stop 参，其终态 `EvError` 发布带 `Stop: stop`（填点从唯一 `:553` 扩到全终态派单原文兑现） |
| `internal/agent/sink.go` | `Event.Stop` 字段注释补去向说明（EvDone 每终态带、空＝loop 自有终态、⛔ 不从 Status 派生） |
| `internal/panel/run_feed_145.go` | 新 `RunDoneRow{Stop string json:"stop,omitempty"}`（`:288-295`）＋ `RunSection.Done`（`:184`，omitempty）＋ `TaskEnded(stop string)`（`:295`，记行＋退 running chips）；文件头「NOT landable」清单把终态词移出、注明 OPTION A 已落 |
| `cmd/wisp/run.go:2505` | `EvDone` 喂线行 `c.runFeed.TaskEnded(string(e.Stop))`——事件词逐字进 done 行 |

### 1b. `ToolRunRow.Level`（判定程 §1 行 4 尾：真源＋构造点＋消费者三件齐）

| 落点 | 内容 |
|---|---|
| `internal/panel/run_feed_145.go:99` | `ToolRunRow.Level string json:"level,omitempty"` |
| `cmd/wisp/run.go:2401` | `bookedArgs` → `bookedForensics(e) (argsJSON, riskLevel string)`：同一枚 booked 行的 `ArgsJSON` 与 `RiskLevel` 一次 join（「args 与 level 来自同一行」由同一循环保证，不造第二次查表） |
| `cmd/wisp/run.go:2444-2447` | `EvToolEnd` 喂线行把 level 递给 `ToolEnded(..., level)` |
| panel `ToolEnded` | 增 riskLevel 参（`:252`）；无 booked 行的调用两值同缺（omitempty）＝「未执行」「无门决策 booked」双真 |

### 1c. AC#2b 两半枚（判定程 §1 两行「卡 internal/config／internal/llm 禁面」——本程写面已含，开）

**半枚① Display：**
- `internal/config/settings_187.go:130-137`：`ChatModelOption` 增 `Display`；`:215,:220` 同一 walk 填 `Models[].Display` 与新 `CurrentDisplay`（一次遍历，名字与对不可能读自不同时刻）；**未声明的条目保持空串，不造名**
- `internal/panel/model_switch_187.go:51`：`ModelOption.Display json:"display,omitempty"`
- `cmd/wisp/run.go:531`：装配根映射收敛为具名函数 `modelSwitchSection(mgr, prov)`（生产与测试共用一体，替掉原「镜面 helper 靠人眼 diff」的形状）；`Current.Display`/`Models[].Display` 逐字带

**半枚② 适配器档位词表：**
- `internal/llm/thinking_levels.go`（新）：`RegisterThinkingLevels`/`ThinkingLevels(protocol)` 导出面——判据照 `RegisterProtocol` 同款纪律（包内 init 自注册、重复注册 panic、返回副本）
- 三家适配器 init **从自家映射表派生注册**（⛔ 不复述第二份词表）：`openairesponses`（`request.go` effortVocabulary() ← effortLevels）、`anthropic`（thinkingVocabulary() ← thinkingBudgets）各注册 `high|low|medium`（排序后）；`openaichat` 注册空表＝**已知协议、零档**（其 adapter doc：intensity 有意不发），与「未注册协议」读数相同、语义用注册动作本身表达
- `internal/panel/model_switch_187.go:88-94`：`ModelSwitchSection` 增 `Adapter`（live provider 自己的 `Info().Protocol` 词）＋ `AdapterThinkingLevels`
- `cmd/wisp/run.go:514,531`：`chatChainProvider(rt)`（无链＝nil，诚实空形）＋映射 `llm.ThinkingLevels(adapter)`
- **resident 无管线形状**（`cmd/wisp/panel_resident_pump.go:166-186`）：Display 照带（ChatSwitchOptions 有）；**Adapter/词表不带**——该形状无 provider 链，词表无源，照该文件自己那句「无源不冒充」留空，⛔ 不猜协议

### 1d. AC#2b 两维本体

**已由 187 落**（`ModelSwitchSection` 五键＋装配根接线），本程零重做；187 的镜面 helper（`panel_model_switch_187_test.go:54`）改为委托 `modelSwitchSection(m, nil)`——生产与测试同体，字节钉（Test187AssemblyRootCarriesTheModelSwitchReader）原样保留。

## 2. AC#6 反向判据（每枚新字段一句）

| 新字段 | 哪一枚用它 | 在哪读 | 读了做什么 |
|---|---|---|---|
| `RunSection.Done.Stop` | 十四行表行 13（已取消/被打断） | 页面读 run 段 done 行 | 区分「写完了」与「被取消/截断」——EvDone.Text 取消时是部分回复正文，表达不了这件事（census §2.14 ⓐᐟ 原话） |
| `ToolRunRow.Level` | 行 4 尾的工具风险档 | 页面读工具行的 level 徽标 | 显示该次调用被门按 L0/L1/L2 哪档放行/拒绝（booked 行的门决策） |
| `ModelOption.Display`／`Current.Display` | AC#2b① 模型清单 | 页面读 picker 行 | 渲染人面向名（配置 display），空则回退 model id |
| `Adapter`＋`AdapterThinkingLevels` | AC#2b② 档位集第二问 | 页面读 modelSwitch 词表 | 把「模型声明了哪些档」∩「本适配器真接受哪些档」两问分开答，只 offer 交集（⛔ 四档原样外发＝AC#2b 明禁） |

四枚全部「真源＋构造点填真值＋消费者」三件齐；每枚都带**反档牙**（突变臂：无 booked 行→level 缺席；无 display 声明→缺席；无 provider 链→adapter 缺席；zero-tier 适配器→词表键缺席；loop 自有终态→stop 词缺席）。**乙组纯无源各枚仍维持既有裁定，零触碰。**

## 3. 门禁读数（导出树 `impl-r1-clone/`，锚＝`4111ada61`；unset GOTMPDIR/TMPDIR/TEMP/TMP＋PATH 加 build、third_party/sherpa-onnx〔主树绝对路径，克隆树无此二 gitignored 目录〕；⛔ 未设 CGO_ENABLED=0）

| 门 | 读数 |
|---|---|
| `/d/work/base/gopath/bin/gofumpt.exe -l`（五包） | **空（rc=0）** |
| `go vet`（panel/agent/config/llm/cmd/wisp） | **rc=0 净** |
| `go build`（同上五包） | **rc=0 净** |
| `go test -count=1 -v ./internal/panel` | **PASS 212 / FAIL 2**；红名册＝`TestC21DesignTokensFourWayAgree`（在册 ball token 族）＋`TestComposerRenderFixtureTellsTheTruth`（环境性：fixture 不入 git）。**零新红**；PASS 212＝判定程轮 210＋本腿新牙 2，对得上 |
| `go test -count=1 ./internal/agent ./internal/config ./internal/llm` | **三包 ok**（9.2s／3.3s／50.1s） |
| `go test -count=1 ./cmd/wisp` | **rc=1＝既有真窗族**：64 FAIL vs **基线对照**（父提交 `71ff0fc87` 同尺同环境）63 FAIL；逐名对拉**唯一差集**＝`TestRunPacketLightsTheRosterRowOfItsLiveL1Window`（活窗 publish 计数时序）；**`-count=3` 复跑 3/3 全绿**（35.4s）⇒ 判 flake 非回归，如实登记不掩盖。本腿新牙（Test145*/Test187* 共 12 枚）**0 FAIL** |
| `sh scripts/d22scan.sh` | **rc=0 clean**（bans #1-5 internal/=388＋cmd/=64、#6 frontend/=121、#7 tools/=48、#10 risk/=21、#8 design/14＋frontend/121＋internal/971＋cmd/243） |
| LF | 本程全部触碰文件 `grep -rlIP "\r"` 零命中（含 gofumpt -w 重排过的 sink.go） |

## 4. Git

只 commit、不 push。commit `4111ada61`（码＋牙，20 文件）与本件＋票面追加（第二枚 commit）均 **git add / commit 两处显式 pathspec**，逐枚点名；`frontend/**`、`forks/**`、`internal/winsec/**`、thresholds/golden、AC 原文零字节。**票面 AC 勾零动**——勾要非实现者表（历波同规矩，本程不自勾）；翻勾材料＝本件 §1/§2/§3。

## 5. 边界与如实

- **resident 无管线形状**的 modelSwitch 不带 adapter 词表（无 provider 链）——这是该形状既有「无源不冒充」posture 的延伸，不是欠账；管线形状（`wisp run`）两枚半枚全带。
- **cmd/wisp 整包 64 红**全部为在册真窗族（63 枚与基线逐名相同）＋1 枚已复跑证伪的活窗 flake；未在本程引入任何新红。
- 临时件只建不删：`impl-r1-clone/`、`impl-r1-base-clone/`、`impl-r1-*.txt` 日志均留在 `.scratch/wisp/probes/145/`。
