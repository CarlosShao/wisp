# 票 387 残余腿 `387-r3` —— 落地件证据表（KWS 六引擎测试补 test-side 守卫 skip＋裸跑 rc0 全绿）

> **身份**：写腿（实现者）自报读数。**判定归验收/重裁腿**——本件不翻任何一格，票面五格只给「读数＋实现腿自评」。
> **尺**：开工 `HEAD=e2ce24d2b`（`dev`，即派单门槛锚本身；索引 0 staged）。共享树他腿在飞件＝`tools/mockllm/**`（391-r1 写面：`messages.go`/`responses.go`/`caps_maxtokens_391r1_test.go`）＋`.scratch/wisp/issues/30-result-routing-d10-done.md` 他人追加＋`.scratch/**` 散件——**全程零碰**。开工核 `git status --porcelain -- internal/speech/`＝**空**。现跑 2026-10-31 21:52–22:1x +08。
> **环境**：Git Bash runner（`D:\work\soft\Git\bin\bash.exe`）；`unset GOTMPDIR` 每跑；rc 一律重定向现取；gofumpt 绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`。
> **staging 形状**（387-r2 §头部同款）：`go test -c -o <tmp>\bin\speech.test.exe ./internal/speech` ＋ **三 DLL 同目录**（`third_party/sherpa-onnx/`：onnxruntime 1.28.2／sherpa-onnx-c-api／sherpa-onnx-cxx-api，实存三枚——390-v1 §7 已具名「四 DLL」计口误）＋ `WISP_SILERO_MODEL`/`WISP_SPIKE_MODELS` 两 env ＋ `-test.v -test.timeout 30m`。
> **裸跑形状**（387-r2 §头部同款、59-r1「System32 1.17.1 遮蔽 PATH 的 1.28.2」复现形态）：**PATH 注入 `third_party/sherpa-onnx`**、无 DLL 贴身、无 `WISP_*`。⚠ 更裸的纯净 PATH（连 sherpa DLL 都不在）＝进程起不来 `exit status 0xc0000135`（链接面缺失，本腿开工第一跑即撞、387-r2 已注明非票机理）——**不是本腿任何读数的形状**。
> **临时件（只建不删）**：`D:\tmp\wisp\20261031\387r3\{bin,export,*.txt,*.rc}`。**⛔ push；⛔ 票面/台账/specs/PLAN/thresholds/golden/testdata 仓根/frontend/design 零字节；⛔ cmd/wisp、internal/{risk,session,config,memory,panel,audio} 零字节（cmd/wisp 只跑读数不写字）；⛔ 守卫本体（tts.go 等）/NewKwsEngine 语义零碰（diff 仅 kws_test.go 一件）。**

## 0 一句话

派单残余①落地：kws_test.go 六处引擎测试（点名 399/416/430/451 四处＋同形两处 336/542）test-top 补同族 `checkOnnxruntime` ⇒ `t.Skipf` 具名 skip（先例＝`cgo_boundary_windows_test.go:157/225`＋387-r2 探针同款）；**裸跑 `go test ./internal/speech` 从 rc1（25/6/36）收敛到 rc0（25/0/42）**；staging 61/0/6 名册 67 枚逐位同基线；仓外导出树摘 skip ⇒ 原六枚逐名复红（同一守卫句），单变量归因闭合。

## 1 ★ 结构（一件：接在哪／为什么）

| 件 | 接在哪 | 为什么是这个形状 |
|---|---|---|
| **六处 test-top 守卫 skip** | `internal/speech/kws_test.go`：`TestKwsEngineHitsTheFixtureWavs`/`TestKwsEngineProductionTableRejectsTheDemoWavs`/`TestKwsEngineOnSilenceStaysSilent`/`TestKwsEngineOnNonKeywordSpeechStaysSilent`/`TestKwsEngineResetAfterHitLetsTheNextUtteranceHitAgain`/`TestKwsEngineCloseIsIdempotent` 各函数体第一行：`if err := checkOnnxruntime(); err != nil { t.Skipf("KWS engine cannot create sessions on this host: %v", err) }`；每测文档句按探针先例补「It skips by name when this host's onnxruntime is shadowed (see checkOnnxruntime): the staging recipe is what makes it run.」 | **只加 test-side 判断、零改帮手**：`newTestKws`/`newTestKwsWithTable`（kws_test.go:157/173）的 `t.Fatalf("NewKwsEngine: …")` **一字未动**——其文档句「on the staging host every ingredient is present, so a failure is a finding, not an environment note」原样成立（staging 路径下 `checkOnnxruntime` 放行 ⇒ skip 块不触发 ⇒ 构造失败仍是红，tripwire 语义零损）。家族先例三点同形全核过：`cgo_boundary_windows_test.go:157`（leg2 前）、`:225`（测顶）、387-r2 `TestKwsProbeSherpaDebug` 探针测顶（387-r2 §1②）。skip 句 `%v` 携守卫具名错误逐字（§3 run 04 句形）。⇒ 未触发停手①：六测在 staging 照跑真断言，裸跑 FAIL 只是环境噪声非验证价值 |

## 2 ★★ 另 2 枚 FAIL 定性（派单②；给编排者/验收的实质新信息）

裸跑基线（run 01，§3）6 枚 FAIL 中派单点名 4 枚（399/416/430/451），另 2 枚＝**同案、非产品问题**：

- **`TestKwsEngineHitsTheFixtureWavs`**（:336 帮手调用）与 **`TestKwsEngineCloseIsIdempotent`**（:542 帮手调用）——**失败句与点名 4 枚逐字同一**：`NewKwsEngine: speech: onnxruntime 1.17.1 is older than the 1.28 the linked sherpa-onnx C API requires; a second onnxruntime.dll is shadowing the one on PATH (on Windows the loader searches System32 before PATH)`（run 01 六段逐字提取，`bare-before2.txt`）。387-r2 §3 早把两枚计入「同形两处」（其 §6① 原文「kws_test.go:399/416/430/451＋同形两处」），派单残余①只点名 4 处行号。
- **fixture 面在裸跑可解析**：`TestKwsEngineHitsTheFixtureWavs` 先走 `kwsFixtureDir`→`kwsModelDir`（kws.go:130，无 env 时从包目录向上八级 walk-up 找 `third_party/spike-models/<kws>`）成功后**才**死在帮手的 `NewKwsEngine`（失败行 336＝帮手调用行）⇒ 唯一阻塞就是守卫，无第二机理。
- **⇒ 处置＝同法 skip**（派单②给出的两枝之一），停手②「真产品问题」未触发；六案在 staging 全部真跑（61/0/6 名册六名全 PASS），真验证价值无损。

## 3 ★ 读数总表（全部 rc 重定向现取）

| # | 跑 | 形状 | rc | PASS/FAIL/SKIP | 判 |
|---|---|---|---|---|---|
| run 01 | **裸跑基线**（改前 HEAD） | PATH 注入 `third_party/sherpa-onnx`，无 WISP_* | **1** | **25/6/36** | 逐位复现 387-r2 run 04 形状；6 枚 FAIL 名单＝§2 六名，失败句六段逐字同（守卫具名错误） |
| run 02 | **裸跑改后** | 同 run 01 | **0** | **25/0/42** | ★派单目标达成：rc0 全绿（36 旧 skip＋6 枚新具名 skip）；六枚 skip 句逐字含守卫具名错误（kws_test.go:336/407/429/448/471/569） |
| run 03 | **staging 基线**（改前 exe） | 三 DLL 贴身＋两 WISP env | **0** | **61/0/6** | 与 387-r2 在册 61/0/6 同 |
| run 04 | **staging 改后**（改后重建 exe） | 同 run 03 | **0** | **61/0/6** | **名册逐位同**：67 枚（名⇒态）Compare-Object＝零差（时长浮动除外）；6 SKIP 仍＝punc 缺席族 |
| run 05 | cmd/wisp KWS 族定向（staging 形状：`wisp.test.exe`＋三 DLL 贴身＋两 WISP env＋cwd＝cmd/wisp 包目录，`-test.run 'Test41a2r1WakeEndToEndArmedToListening\|Test41a2r1WarmGateFeedsKwsNotTheAsr'`） | 同上 | **0** | **2 PASS / 0 FAIL** | 16.26s/13.86s——33r11-r2/390-r1 后该族全绿，零回归 |

## 4 ★★★ 反向判据（仓外导出树摘 skip ⇒ 复红；三读数＋还原核）

导出树＝`D:\tmp\wisp\20261031\387r3\export`（**仓外**；go.mod/go.sum＋speech 依赖五 internal 包＋`third_party/spike-models` 的 KWS 目录与 silero_vad.onnx——首跑对照曾 11 FAIL 全为「silero model not found」修剪伪影，补件后对照即同形；`tools/d22scan` 拷入供扫描）。**全程 repo 零突变**：

- **读数 A（同一性）**：导出树 `kws_test.go` SHA256＝`D4AAD4B24B48F932F4AC658CFEA62EF2AF9ED30FBB10F0F34ACE08D5A0297E5C` ＝ repo 工作树逐位同；对照裸跑（run 06a）⇒ **rc=0，25/0/42** 与 repo run 02 逐位同形。
- **真改**：`edit`（replace_all）摘除六处 `if err := checkOnnxruntime(); …t.Skipf…}` 三行块（`Skipf("KWS engine cannot create sessions` 计数 0；探针 387-r2 的 raw-sherpa skip 句不同、原样保留——突变面恰＝本腿六块）。
- **读数 B（复红）**：突变树裸跑（run 06b）⇒ **rc=1，25/6/36**，复红名单＝§2 六名**逐名同**；失败句逐字＝同一守卫具名错误（首枚 `kws_test.go:338`＝摘块后帮手行位移＋2，单变量归因闭合：只少 skip 块、复红即原 FAIL）。
- **还原核（读数 C）**：repo `kws_test.go` SHA256 恒＝`D4AAD4B2…97E5C`（突变前后两次现取同值）；`git status --porcelain -- internal/speech/` 恒＝本腿一件 `M kws_test.go`。

## 5 ★ 门禁（两尺）

| 尺 | 读数 |
|---|---|
| `gofumpt -l kws_test.go`（绝对路径） | **空**（rc 0） |
| `go vet ./internal/speech` | **rc 0** |
| `d22scan`（`tools/d22scan/d22scan.exe`，仓根） | **rc 0 clean**——四数＝bans #1-5 internal/**287**、bans #1-5 cmd/**52**、ban #6 frontend/**116**、ban #7 internal/tools/**25**；ban #8＝design/39、frontend/116、internal/**729** Go files（与 387-r2 在册 729 同）、cmd/180（390-r1 三枚 cmd/wisp 落地后 +3，逐源可数） |
| d22scan 导出树＋正控 | 导出树扫描 rc 2＝**30 条 `[phantom-citation]` 修剪伪影**（cited 件在真 repo 存在、未拷入导出树；control/突变/摘探针三轮计数恒 30，非本腿引入；repo 正扫 rc 0 clean 已证门真）。**正控两发中双旗**：production 探针件 `[bare-goroutine]`（:8）＋`[emoji]`（:11）rc 2；摘探针 ⇒ zz 命中 0（回到 30 条伪影基线）＝可逆。注：bare-goroutine 旗只扫 production 件——`_test.go` 探针只中 emoji，已换 production 探针补双旗 |

## 6 写面申报（逐枚）

1. `M internal/speech/kws_test.go` —— 唯一代码件：六处 test-top 守卫 skip＋每测文档句（+35/−6，diff 全在六测函数体与注释；`newTestKws`/`newTestKwsWithTable`/`checkOnnxruntime`/引擎文件零字节）
2. `A docs/evidence/s1/387r3-impl.md` —— 本件
**提交＝本腿唯一一笔（两文件显式 pathspec，主题逐字 `票 387 残余腿 387-r3`）。⛔ push。**

## 7 ★ 停手三条核销

①skip 损验证价值？——未触发：六测 staging 全真跑（run 03/04 名册六名全 PASS），帮手 `Fatalf` tripwire 原样（staging 构造失败仍是红）；skip 只在守卫拒绝的裸跑主机触发，skip 句含守卫具名错误。②另 2 枚 FAIL 真产品问题？——未触发：与点名 4 枚**失败句逐字同一**（§2），fixture 面裸跑可解析、唯一阻塞是守卫。③必须动守卫/引擎？——未触发：diff 仅测试件，守卫本体与 `NewKwsEngine` 语义零字节。
