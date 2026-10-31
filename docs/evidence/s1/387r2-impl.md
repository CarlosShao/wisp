# 票 387 实现腿 `387-r2` —— 落地件证据表（`NewMatchaTts` 补 `checkOnnxruntime` 守卫＋裸跑零崩溃）

> **身份**：写腿（实现者）自报读数。**判定归验收/重裁腿**——本件不翻任何一格，票面五格只给「读数＋实现腿自评」。
> **尺**：开工 `HEAD=46a80b579`（`dev`，`merge-base --is-ancestor 46a80b579 HEAD` rc 0；索引 0 staged；`git status --porcelain -- internal/speech/` 开工前**空**）。共享树他腿在飞件＝`tools/mockllm/**`（386-r2 写面）＋`.scratch/**` 散件＋`.scratch/wisp/issues/30-result-routing-d10-done.md` 他人追加——**全程零碰**。现跑 2026-10-31 21:18–21:5x +08。
> **环境**：Git Bash runner（`D:\work\soft\Git\bin\bash.exe`）；`unset GOTMPDIR` 每跑；rc 一律重定向现取；gofumpt 绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`。
> **staging 形状**（票面 AC#4 集合）：`go test -c -o <tmp>\bin\speech.test.exe ./internal/speech` ＋ **三 DLL 同目录**（`third_party/sherpa-onnx/`：onnxruntime 1.28.2／sherpa-onnx-c-api／sherpa-onnx-cxx-api）＋ `WISP_SILERO_MODEL`/`WISP_SPIKE_MODELS` 两 env（**无** `WISP_PUNC_MODEL`——不在票面集合，punc 族走既有具名 skip）＋ `-test.v -test.timeout 30m`。
> **裸跑形状**（59-r1「System32 1.17.1 遮蔽 PATH 的 1.28.2」的复现形态）：**PATH 注入 `third_party/sherpa-onnx`**（供给 c-api/cxx-api 与 1.28.2 ORT——Windows 装载序 System32 先于 PATH ⇒ ORT 实取 System32 的 1.17.260311，本机现验）、**无 DLL 贴身、无 `WISP_*`**。⚠ 纯净 PATH（连 sherpa DLL 都不在）时进程起不来＝`exit status 0xc0000135`（DLL 未找到，`logs/01-…log`）——那是「链接面缺失」，不是本票机理；本票机理需要 DLL 可载而 ORT 被遮蔽。
> **临时件（只建不删）**：`D:\tmp\wisp\20261031\387r2\{logs,bin-base,bin-fixed,export}`。**⛔ push；⛔ 票面/台账/specs/PLAN/thresholds/golden/testdata/frontend/design 零字节；⛔ cmd/wisp、internal/{risk,session,config,memory,panel,audio} 零字节；引擎语义零碰（守卫只在不满足时早退具名错误）。**

## 0 一句话

守卫家族漏网的 `NewMatchaTts` 挂上同形 `checkOnnxruntime`（AC#1）；裸跑从「`0xc0000005` 崩进程」变成「具名 skip／具名拒绝」（AC#2）；仓外导出树**真摘守卫**⇒崩溃逐字复现（AC#3）；staging 61/0/6 名册逐位同基线（AC#4）。**顺带校正 59-r1 归因一处**：全包 `./internal/speech` 裸跑的**首发**崩点不是 `NewMatchaTts`，是票 41 老测试 `TestKwsProbeSherpaDebug`（同样无守卫直调 C API）——已一并补 test-side 同族 skip，见 §2。

## 1 ★ 两件结构（逐件：接在哪／为什么）

| 件 | 接在哪 | 为什么是这个形状 |
|---|---|---|
| **①AC#1 主件：`NewMatchaTts` 挂守卫** | `internal/speech/tts.go` `NewMatchaTts`：config 校验 switch 之后、`sherpa.NewOfflineTts` 之前，插 `if err := checkOnnxruntime(); err != nil { return nil, err }`；文档句按家族同形补「It refuses, rather than crashes, …: see checkOnnxruntime.」 | 家族先例四点同形全核过：`NewDetector`(vad.go:72)、`NewStreamingAsr`(asr.go:71)、`NewKwsEngine`(kws.go:213)、`NewCTTransformerPunctuator`(punctuate.go:185)——**全部是「config 校验后、sherpa 构造前」同一位置同一句形**。守卫本体 `checkOnnxruntime`(vad.go:102) 只读 `GetOnnxruntimeVersion()` 比对 `requiredORTMinor=28`，零引擎语义触碰；staging 路径下 1.28.2 通过、行为逐位不变（AC#4 实证） |
| **②AC#2 配套：`TestKwsProbeSherpaDebug` 补 test-side skip** | `internal/speech/kws_test.go` 探针测试第一行：`if err := checkOnnxruntime(); err != nil { t.Skipf("raw sherpa probe cannot create sessions on this host: %v", err) }` | 该测试**绕过一切 Go 构造器直调 `sherpa.NewKeywordSpotter`**(kws_test.go:502)，低于所有生产守卫——裸跑全包时它**先于 TTS** 把进程打死（§2），AC#2 的「零崩溃」字面靠它才成立。test-side 同族 skip 在本包有先例：`cgo_boundary_windows_test.go:157/225`（`rawVAD` 也只在该守卫之后被调）。staging 下探针照跑（`--- PASS: TestKwsProbeSherpaDebug` 在基线与改后名册都在）＝零行为改变 |

## 2 ★★★ 归因校正（给编排者/验收的实质新信息）

**59-r1 的「`go test` 直跑会在 `NewMatchaTts` 崩」在它自己的 aecspike 包视角为真**（`59r1-aec-spike.md` §3 原文＋该包无 kws 探针）；**全包 `go test ./internal/speech` 裸跑实测首发崩点在 `TestKwsProbeSherpaDebug`**（票 41 `d8dcf4c31` 就在的老测试，kws_test.go 按文件序先于 tts_test.go 执行）：

- **run 01b**（HEAD 改前全包裸跑）：`rc=1`，栈＝`_Cfunc_SherpaOnnxCreateKeywordSpotter` ← `TestKwsProbeSherpaDebug`（`logs/01b-bare-prefix-path.log`）——进程死在 kws 探针，**TTS 测试根本没轮到跑**。
- **run 02**（HEAD 改前**定向** `-run 'Test26'` 裸跑）：`Exception 0xc0000005 0x0 0x18 0x7ffd…`，栈逐字＝`_Cfunc_SherpaOnnxCreateOfflineTts` ← `NewOfflineTts(sherpa_onnx.go:1359)` ← **`NewMatchaTts` ← `tts.go:158`**（`logs/02-bare-prefix-targeted.log`）——**票面点名的崩点原样坐实**，且 ORT 自报机理行：`The requested API version [28] is not available, only API versions [1, 17] are supported in this build. Current ORT Version is: 1.17.1`。
- **停手条件②核销**：条件要求「裸跑崩点不在 NewMatchaTts ⇒ 停手具名真崩点」。实测＝**两处都真**：`NewMatchaTts` 崩点存在且可定向复现（票面机理不误）；另有**更早**的同族无守卫直调点（kws 探针）把全包首发崩点挪了位。据「最有价值」铁律与写面明示「可能一个测试」，处置＝主件照修＋探针补同族 skip＋**两处具名上报**，不扩到其余测试断言形状（§6 残余①）。

## 3 ★ 五格逐条（实现腿自评；**翻框归验收**）

### AC#1 守卫补齐 —— 读数齐
`tts.go` diff 恰 +9 行（文档 3＋守卫 3＋空行界 3 中实插 6 行代码/注释）；`git diff` 全文与四先例逐形对照无差（§1①）。家族盘点（`sherpa.New|Get` 全包 11 处）：生产面唯 `tts.go:158` 无守卫，其余全在守卫后；测试面唯 kws 探针无守卫（`rawVAD` 已被 157/225 两守卫罩住）——**补完后全包直调点 100% 被守卫覆盖**。

### AC#2 裸跑不崩 —— 读数齐（具名 skip＋零崩溃）
**run 04**（改后全包裸跑，run 01b 同形态）：**零崩溃**，跑完全程 `rc=1`（失败非崩溃），**25 PASS / 6 FAIL / 36 SKIP**（`logs/04-bare-postfix.log`）：
- 全部能崩调用点均具名 skip：matcha 三测（`Test26TheTtsEngineEmitsRealAudio`/`Test26TheClosedEngineRefusesToSpeak`/`Test26TheRealSwitchDoesNotAddTheTwoFamilies`）skip 句＝家族原句形 `matcha tts unavailable on this host: speech: onnxruntime 1.17.1 is older than the 1.28 the linked sherpa-onnx C API requires…`；kws 探针 skip 句＝本腿新句（§1②）；asr/vad/punc/cgo_boundary 全族照旧具名 skip。
- **6 个 KWS 引擎测试 FAIL（非 skip）＝HEAD 既有裸跑形状、非本腿引入**：run 01b（改前 HEAD）同批 5 枚已在崩溃前 FAIL（第 6 枚 `TestKwsEngineCloseIsIdempotent` 在探针之后、当时没轮到跑）。失败句＝测试**直调 `NewKwsEngine` 后 `t.Fatalf("NewKwsEngine: %v", err)`**（kws_test.go:399/416/430/451），错误文本正是守卫具名拒绝——「守卫拒绝以具名失败呈现」，不崩、不哑。归 AC#2 字面「具名 skip、零崩溃」：**零崩溃达、能崩点全具名 skip；这 6 处的 skip 化超出 AC#1 字面（只点名 NewMatchaTts），本腿不扩**（残余①，后续票一句线的事）。

### AC#3 反向判据 —— 读数齐（真改文件＋三读数＋还原核）
**A380 WalkDir 同款方法：仓外导出树真改文件（⛔ -overlay——守卫的拒绝/崩溃都发生在运行时，必须真改编译输入以外的盘上件）**。导出树＝`D:\tmp\wisp\20261031\387r2\export`（robocopy 工作树，剔除 `.git/.scratch/logs/design/frontend/node_modules/build/docs`；`third_party` **实拷**自持，模型/DLL 全在）。
- **读数 A（改前）**：`export\internal\speech\tts.go` SHA256＝`3F4C7DA8…A354` ＝ repo 工作树**逐位同**（导出树裸跑对照 `run 05a`：`-run Test26` ⇒ **rc=0 全具名 skip 零崩溃**——与 repo run 04 同形）。
- **真改**：`edit` 摘除 `if err := checkOnnxruntime(); err != nil { return nil, err }` 三行（守卫块；文档句留作注释，无函数引用）。
- **读数 B（改后）**：SHA256＝`F9FEECC3…F678`（已变）；`if err := checkOnnxruntime` 计数 0；**repo 侧 SHA256 原值未动**。
- **崩溃复现（run 05b）**：同导出树裸跑 `-run Test26` ⇒ **`Exception 0xc0000005 0x0 0x18 0x7ffd…`**，栈逐字 `_Cfunc_SherpaOnnxCreateOfflineTts` ← `NewOfflineTts(…@v1.13.8/sherpa_onnx.go:1359)` ← **`NewMatchaTts` ← `export/internal/speech/tts.go:161`** ← `Test26TheRealSwitchDoesNotAddTheTwoFamilies`（`logs/05b-export-mutant.log`）——**摘守卫即崩、崩点即 AC#1 修的构造器，单变量归因闭合**（161＝158+3，恰为文档句行移， mutant 只少守卫块）。
- **读数 C（还原核）**：repo `tts.go` 拷回导出树 ⇒ SHA256 复等（`3F4C7DA8…`）；守卫块计数回 1；全程 `git status --porcelain -- internal/speech/` 恒＝本腿两枚（`M tts.go`/`M kws_test.go`），**repo 零突变窗口**。

### AC#4 零回归 —— 读数齐（逐位同）
基线（run 03，**改前 HEAD 原样**）与改后（run 06）同一 staging 形状各跑全量：

| 尺 | 基线（改前） | 改后 | 判 |
|---|---|---|---|
| rc | 0 | 0 | 同 |
| 计数 | **61 PASS / 0 FAIL / 6 SKIP** | **61 / 0 / 6** | 同 |
| 名册 | 67 行（roster.txt） | `diff` 基线 vs 改后＝**空（rc 0）** | **逐位同** |
| 探针 | `--- PASS: TestKwsProbeSherpaDebug (2.69s)` | `--- PASS: TestKwsProbeSherpaDebug`（名册行同） | staging 下行为不变 |

6 SKIP 全为 punc 模型缺席族（票面 staging 集不含 `WISP_PUNC_MODEL`，与 41a2r1 记录的 61/0/6 同源）。守卫在 staging 路径（1.28.2 同目录）下放行，matcha 五测真合成照跑（`Test26TheTtsEngineEmitsRealAudio (5.34s)` 等）——**守卫零改 staging 行为**。

### AC#5 纪律格
写面恰 3 枚（§5 申报）；`cmd/wisp`／`internal/{risk,session,config,memory,panel,audio}`／thresholds/golden/testdata 仓根/票面/台账/specs/PLAN/frontend/design 零字节（现跑 `git status` 证）；⛔ push；票面 Status 未动；台账未动。

## 4 ★ 门禁

| 尺 | 读数 |
|---|---|
| `gofumpt -l tts.go kws_test.go`（绝对路径） | **空**（rc 0） |
| `go vet ./internal/speech` | **rc 0**（`logs/07-vet.log`） |
| `d22scan`（`tools/d22scan/d22scan.exe`，仓根） | **rc 0 clean**——`internal/=729 Go files` 分母含本腿两件（`logs/07-d22scan.log`） |

## 5 写面申报（逐枚）

1. `M internal/speech/tts.go` —— AC#1：`NewMatchaTts` 守卫＋家族同形文档句（+6 行）
2. `M internal/speech/kws_test.go` —— AC#2 配套：`TestKwsProbeSherpaDebug` test-side `checkOnnxruntime` skip＋注释三句（+8 行）
3. `A docs/evidence/s1/387r2-impl.md` —— 本件
**提交＝本腿唯一一笔（三文件显式 pathspec，主题逐字 `票 387 落地腿 387-r2`）。⛔ push。**

## 6 ★ 残余具名（不阻塞、不在 AC#1 字面）

① **6 个 KWS 引擎测试裸跑 `Fatalf` 非 skip**（§3 AC#2）——测试侧对守卫错误的断言形状，改之超出本票写面授权（AC#1 只点名 `NewMatchaTts`，「可能一个测试」已用于探针）；要裸跑全绿 rc=0，须把这 4 处 `t.Fatalf("NewKwsEngine: …")` 改为守卫错误具名 skip（kws_test.go:399/416/430/451＋同形两处），**一句线后续票**。② 59-r1 归因校正（§2）——台账口径如需修订（`pending-and-issues.md:28956` 的「`NewMatchaTts` 崩」表述对全包裸跑不完整），归编排者。③ 纯净 PATH 形态 `0xc0000135`（DLL 缺失）≠ 本票机理，已在头部注明，防后人把「起不来」当「守卫没效」。

## 7 ★ 未定义即停两条核销

①守卫家族先例形状——`NewDetector` 的守卫**正是** `checkOnnxruntime` 形（vad.go:72-74），四先例同位同句形 ⇒ 未触发；②裸跑崩点——**NewMatchaTts 崩点坐实**（run 02 定向逐字复现于 tts.go:158），另有更早的 kws 探针崩点（§2 具名上报并同族封堵）⇒ 未触发停手（前提「修的就是 NewMatchaTts」成立）。
