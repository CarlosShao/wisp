# 票 398 终裁腿 `398-v2` —— silero RNN hangover 的 session reset 面（非实现者裁决表）

日期 2026-11-01 · 裁决对象 `ced345c3c`（3 枚 +289：`internal/speech/vad.go` Reset 一枚 +13/−0、
`vad_reset398_test.go` 116 行、证据件 160 行；父 `d26a29083`）· 材料现核：票面 5 格＋
`398r1-hangover-reset.md`＋`392r2-keepalive-fix.md`（探针读数序）＋`392-v3-verdict.md`（采纳路径②）
· 本腿七项全部自跑；全部实验在导出树 `D:\tmp\wisp\20261101\398v2\**`（`git archive` 三树：
base=d26a29083、impl=ced345c3c、tree=HEAD）；`git diff ced345c3c HEAD -- internal/speech/`＝空＝
实现自落地零再动；仓内只写本件；⛔ push。

## 0. 一句话

**五格全过、可收口。** 形状核（Reset 只清状态核、Endpointer/Finalize 零字节、RNN 状态三锚点
具名）、探针复跑 `(silence=1 talker=0)⇒(silence=0 talker=1)` 逐字（rc=0）、口径差裁定成立、
等价案按段计数＋突变实测咬合、摘 `d.vad.Reset()` 双红＋同突变旁案绿＋还原核四件、零回归
staged 69/63/0/6＋基线 67 枚名册逐位同序＋裸跑 25/0/44＋cmd/wisp KWS 41a2r1 两案 2/2＋
门禁 gofumpt/vet/d22scan clean（selftest 35/0，RUN=77）。

## 1. ★AC#1 reset 面形状核 ✓

- **diff 现读**：+13/−0 全在 `Finalize` 之后的新方法块；`Accept`/`Finalize`/`Consume`/`Close`/
  `acceptWindow`/Endpointer **零字节**——⛔ 动 silero 推理语义与 Endpointer/Finalize 语义核销。
- **形状对照**：`Finalize`=零垫尾窗→分类→`ep.Finalize`（关开口话态）；`Reset`=`d.rest[:0]`＋
  `d.vad.Reset()`（vad.go:185-188）。两面正交：开口话态仍归 Finalize、session 态归 Reset，
  组合＝owner 完整窗间清。只清状态三件：silero RNN 语境（hangover 载体）＋C 侧段队列（Go 侧
  `acceptWindow` 本就每窗 drain，行为恒等）＋`d.rest` 半窗缓冲（防旧样本拼进新会话）。
- **RNN 状态位置三锚点现核**（v1.13.8，go.mod pin 现读）：①`sherpa-onnx-go-windows@v1.13.8`
  `sherpa_onnx.go:1703` `impl *C.struct_SherpaOnnxVoiceActivityDetector`（状态在 C 侧会话对象内
  持有）；②同件 `:1787-1788` `func (vad *VoiceActivityDetector) Reset()` →
  `C.SherpaOnnxVoiceActivityDetectorReset(vad.impl)`；③同包 `c-api.h:2235` 该 API 声明＋注释
  「Reset a voice activity detector so it can process a new stream」逐字。面已存在、直接可达。

## 2. ★AC#2 hangover 清除读数 ✓＋口径差裁定（成立）

本腿 staged 复跑（自建 exe＋三 DLL 贴身＋两 WISP env；log `probe-staged.log`，rc=0）：

| 臂 | 本腿读数 | 与申报 |
|---|---|---|
| PRE-RESET（392-r2 控制形） | `(silence=1 talker=0)` | 逐字同 |
| POST-RESET | `(silence=0 talker=1)` | 逐字同 |
| RESET-EQUALS-FRESH | `(gapA=0 talker=1 gapB=0)` 两臂全等 | 逐字同 |

**口径差裁定：成立，实现侧读数是正确口径**——①AC#2 主判据＝silence=0（post-speech 静音不再
开火＝hangover 清）；②talker=1 是判定通路活性的**内建正控**（Reset 只准清状态、不准碰分类器）；
③字面 0/0 若成立恰是「reset 清死判定通路」＝票面停手③——392-r2 探针自己的 Fatalf
（vad_reset398_test.go:56-59 把 `nSil+nTalk==0` 判 MECHANISM MOVED）与 post 臂断言（:72-76 硬钉
`(0,1)`）都把 0/0 定为红。派单「0/0」只能读作「silence=0＋无 hangover 假 start」；(0,1) 是该
语义的可观测正确形。判据差不是放宽——是「0/0＝死」与「(0,1)＝清＋活」的判然之分。

## 3. ★AC#3 等价案核＋反向复刻＋还原核 ✓

- **等价案核**：`Test398ResetMatchesAFreshDetector` 用 per-segment 三段计数（gapA/talker/gapB，
  vad_reset398_test.go:91-93，落点敏感）；fresh 臂 gapA/gapB=0＝`TestVadPureSilenceProducesNoEndpoint`
  的窗级复刻正控。**开发史核**：kind 序列初版对落点不敏感（dirty 的 start 落 gap、fresh 的落
  talker，`[start end-of-stream]` 两臂同形）而误绿——本腿突变复跑实证改法咬合：摘 `d.vad.Reset()`
  后等价案红 `(gapA=1 talker=0 gapB=0) vs (0 1 0)`（`mutant-test398.log`）——落点挪移被咬住，
  单变量归因闭合。
- **反向复刻（本腿导出树自做突变、自建突变二进制）**：sed 单行 `d.vad.Reset()` →
  `_ = d.vad // MUTANT 398r1-M1`；①主案红：POST-RESET 退回 `(silence=1 talker=0)`＋具名句
  「the hangover survived the reset」rc=1（PRE 控制臂 (1,0) 不动）；②等价案同红 rc=1；③同突变
  二进制旁案 `TestVadPauseDetail`/`TestVadPureSilenceProducesNoEndpoint`/`TestStateMachineIgnoresChunking`
  全绿 rc=0＝突变只打 hangover 形、不打既有判定通路。
- **还原核四件**：①同 sed 再生突变与实现者突变件 **diff 零差**（rc=0/0 字节；注记：用短标记
  `_ = d.vad // MUTANT 398r1-M1` 再生与实现件差 1 行注释后缀、功能置换同一单行——实现件完整
  标记文本下零差）；②pristine ≡ `HEAD:vad.go`（diff rc=0）；③仓内
  `git status --porcelain -- internal/speech/`＝0 行；④正树探针复跑 rc=0 绿（`probe-tree-final2.log`）。

## 4. ★AC#4 零回归（本腿自跑）✓

| 跑 | rc | 读数 |
|---|---|---|
| internal/speech 全套 staged | 0 | **RUN=69 / PASS=63 / FAIL=0 / SKIP=6**；名册＝基线 67 枚**逐位同序**＋恰两枚 Test398 新增；SKIP 集逐名同 |
| 基线 `d26a29083` staged（导出树） | 0 | RUN=67 / PASS=61 / FAIL=0 / SKIP=6（在案基线 61/0/6 同） |
| internal/speech 裸跑（PATH 注入、无 WISP_*） | 0 | **25 PASS / 0 FAIL / 44 SKIP**；25 PASS 逐位同基线裸跑（基线 25/0/42 本腿复跑同）；SKIP 仅＋2 具名 Test398（守卫句＝onnxruntime 1.17.1 System32 影子，在案） |
| cmd/wisp KWS 族 41a2r1 两案 staged | 0 | **2 PASS / 0 FAIL**（18.86s/15.87s 同带）；build 前后 cmd/wisp 恰一枚他腿脏件（397-r4 `resident_voice_bargein_392_windows_test.go`，不在本族、本腿零归因零触碰） |

- 基线红插记（如实）：基线 staged 第一发 rc=1＝**本腿导出树形状缺口**（`TestDefaultKwsConfigResolution`
  用 runtime.Caller 从源码路径 walk-up 并与 env 解析串比较——导出树不自洽）；junction 补
  third_party＋env 指向导出树后 61/0/6 全绿。非基线红、非 398 红。
- 门禁：gofumpt `-l internal/speech/`＝0 行 rc=0；`go vet ./internal/speech` rc=0；d22scan
  clean rc=0＋selftest PASS=35 FAIL=0（RUN=77）。**分母漂移归属**：HEAD 现 ban#8 internal/=734
  （ls-tree .go 计数 734 同）；impl 时 **732**（ls-tree 732 同＝派单口径）；+2 全来自
  `ced345c3c..HEAD` 的 394-r1（internal/risk 测试件）与 400-r1（panel/winsec）；cmd 183→184
  同因、两侧各 −1 为工具系统口径。ban#8 不扫 docs/（392-v3 在案口径），本件不入分母。

## 5. ★五格逐条判＋可否收口

- **AC#1 ✓**（§1：只清状态核、语义零字节、三锚点具名）。
- **AC#2 ✓**（§2：读数逐字＋口径差裁定成立）。
- **AC#3 ✓**（§3：双红＋旁案绿＋还原核四件＋按段计数改法实证咬合）。
- **AC#4 ✓**（§4：三跑逐位＋门禁 clean＋分母漂移具名归属）。
- **AC#5 纪律 ✓**（398-r1 写面恰 3 枚 pathspec 现核；越界面 `cmd/wisp|internal/audio|risk|panel`
  diff＝0 行；本腿仓内只写本件、实验全在导出树、cmd/wisp 只读；无 push）。
- **可收口：是**（票面 Status/勾选归编排者收口动作，本腿⛔改）。

## 6. 残余与边界（如实记账）

1. **C++ 源码级核仍缺**：本腿同拦（raw.githubusercontent/github/jsdelivr/docs.rs 四路均解析到
   非公共 IP）——行为级三面（探针/全套回归/突变反向）闭合维持 398r1 定性；补核需可通 GitHub
   的环境（上游 `sherpa-onnx/csrc/voice-activity-detector.cc` 的 `VoiceActivityDetector::Reset()`）。
2. **owner 接入＝后续裁定**：`Finalize(nil)＋Reset()` 免冷段属 392-v3 §6 第 3 步（keepalive 翻转
   前置）的实现面；398-r1 与本腿对 owner（internal/audio、cmd/wisp 生产面）均零字节。
3. **hangover 衰减扫参未关闭**（「衰减多少够」，398r1 §6.3 原样）——reset 面使其失去 owner 侧
   紧迫性，但读数仍只有 fixture 级一点。
4. 导出树形状缺口两处如实记（walk-up 自洽、d22scan selftest 的 self-scan 两案需真仓）——实验
   环境形状，非读数判据、非回归。

## 附：本腿自跑尺（全部可重放）

`D:\tmp\wisp\20261101\398v2\`：三导出树（base=d26a29083／impl=ced345c3c／tree=HEAD，均
`git archive`＋junction 补 third_party）；staging＝自建 `speech398v2`/`speech-base398`/
`mut398v2`/`wisp398v2` 四枚测试二进制＋三 DLL（与 third_party 同 SHA256）；logs/＝
{probe-staged, suite-staged, suite-base, suite-bare, suite-bare-base, mutant-test398,
mutant-bystander, probe-tree-final2, kws-41a2r1-staged, gofumpt.out, vet-speech,
d22scan, d22scan-at-impl}；mut/＝{vad.pristine, vad.mut398v2, vad.mut.regenerated}＋
logs/mut-vs-impl.diff。Git Bash runner（`D:\work\soft\Git\bin\bash.exe`）、`unset GOTMPDIR`
每跑、rc 一律重定向现取。
