# 票 398 落地腿 `398-r1` —— silero RNN hangover 的 session reset 面（写腿）

日期 2026-11-01 · 腿 `398-r1`（写、`internal/speech/**` 独占＋本件）·
起手 HEAD `d26a29083`（dev，开工三核过：①锚点在（＝派单门槛本身）②索引 0
staged ③`git status --porcelain -- internal/speech/ docs/evidence/s1/`＝空）·
派单存档 `.scratch/wisp/probes/orch/20261101/txt/398-r1-dispatch.txt` 在案 ·
现跑 2026-11-01 15:34–16:0x +08。

> **环境**：Git Bash runner（`D:\work\soft\Git\bin\bash.exe`）；`unset GOTMPDIR` 每跑；
> rc 一律重定向/`PIPESTATUS` 现取；gofumpt 绝对路径 `D:\work\base\gopath\bin\gofumpt.exe`。
> **staging 形状**（387-r3 §头部同款、A1326 卡）：`go test -c -o <tmp>\staging\speech398.test.exe
> ./internal/speech` ＋ 三 DLL 同目录（onnxruntime 1.28.2／sherpa-onnx-c-api／sherpa-onnx-cxx-api）＋
> `WISP_SILERO_MODEL`/`WISP_SPIKE_MODELS` 两 env ＋ `-test.v -test.timeout 30m`。
> **裸跑形状**（387-r3 §头部同款）：PATH 注入 `third_party/sherpa-onnx`、无 DLL 贴身、无
> `WISP_*`。⚠ 纯净 PATH（连 sherpa DLL 都不在）＝进程起不来 `exit status 0xc0000135`
> （387-r3 已注明非读数形状，本腿第一跑即撞、照方换 PATH 注入形状）。
> **临时件（只建不删）**：`D:\tmp\wisp\20261101\398r1\{staging,mut,logs}\**`。
> **⛔ push；⛔ 票面/台账/specs/PLAN/thresholds/golden/testdata 仓根/frontend/design 零字节；
> ⛔ cmd/wisp、internal/{risk,session,audio} 零字节（cmd/wisp 只跑读数不写字）。**

## 0. 一句话

**silero session reset 面＝`Detector.Reset()` 落地**（对照 `Finalize` 形状、只清状态：
RNN 语境＋sherpa 段队列＋rest 半窗缓冲；Endpointer 与 Finalize 语义零字节）——
探针把 392-r2 的具名 hangover 读数钉进本包：**reset 前 `(silence=1 talker=0)`
（假 start 落 post-speech 静音段）、reset 后 `(silence=0 talker=1)`（gap 干净、
talker 正常开火）**；等价案 `398-RESET-EQUALS-FRESH (gapA=0 talker=1 gapB=0)`
两臂全等（落点敏感读数）。反向：仓外 `-overlay` 摘 `d.vad.Reset()` ⇒ 主案红
`(silence=1 talker=0)`＋等价案红 `(gapA=1 talker=0 gapB=0) vs (0 1 0)`，同突变
二进制旁案三枚照常绿＝定向自证；还原核（突变件单行 delta 核＋正树复跑绿）。
零回归：staged 全套 **rc=0 / 69 RUN / 63 PASS / 0 FAIL / 6 SKIP**（基线 67 枚
61/0/6＋本腿两枚新增绿、SKIP 集仍＝punc 缺席族同形）、裸跑 **rc=0 / 25 PASS /
0 FAIL / 44 SKIP**（387-r3 基线 25/0/42＋本腿两枚具名 skip——rc0 全绿未破）、
cmd/wisp KWS 族 41a2r1 两案 staged **2 PASS / 0 FAIL**（16.27s/14.75s，与
387-r3 run 05 同带）；gofumpt/vet/d22scan 全 clean。

## 1. AC#1 reset 面设计（对照 detector Finalize 形状）

**新 API 一枚**（`internal/speech/vad.go`，置于 `Finalize` 之后）：

```go
func (d *Detector) Reset() {
	d.rest = d.rest[:0]
	d.vad.Reset()
}
```

- **清什么**：①silero RNN 语境（hangover 的载体——392-r2/392-v3 机理：无输入
  不衰减、能在 post-speech 静音上开火）；②sherpa 自己的段队列（C API Reset
  语义「so it can process a new stream」的一部分；Go 侧 `acceptWindow` 本就每窗
  drain，行为恒等）；③`d.rest` 半窗缓冲（session reset 语义＝丢弃上一会话的
  输入尾巴，防「Accept 半窗后 Reset」把旧样本拼进新会话）。
- **不碰什么**：`Endpointer` 一字未动——开口话态的清面仍是 `Finalize`（392-v2/
  393-r2 窗间清在用、语义已定）；两面正交组合＝owner 完整窗间清
  （`Finalize(nil)` ＋ `Reset()`），392-v3 §6.2 的「冷段」从此有更硬的替代。
  `Finalize` 本体零字节。classifier 零字节（cfg/threshold/模型/窗口逻辑原样）。
- **状态位置具名**（停手条件①核销）：RNN 状态在 C 侧会话
  `*sherpa.VoiceActivityDetector`（Go 绑定 1.13.8，`impl *C.struct_SherpaOnnxVoiceActivityDetector`）
  内持有；C API 提供 `SherpaOnnxVoiceActivityDetectorReset`
  （windows 包 `c-api.h:2235`，注释「Reset a voice activity detector so it can
  process a new stream」），Go 绑定已暴露 `vad.Reset()`
  （`sherpa-onnx-go-windows@v1.13.8/sherpa_onnx.go:1787`）——**面已存在、直接可达**。
  ⚠ 如实记：sherpa-onnx C++ 侧 `VoiceActivityDetector::Reset()` 的源码级逐行核
  未做成（本环境网络策略拦 GitHub/raw.githubusercontent/docs.rs）；以上游 C API
  头文件注释（在案）＋行为级读数（§2 探针＋§4 全套回归＋§3 突变反向）为证。

**停手三条核销**：①「无状态可清」——未触发（状态具名且绑定可达，见上）；
②「必须动 Finalize 语义」——未触发（diff 只增 Reset 方法块；Endpointer 零字节）；
③「reset 后判定行为改变」——未触发（Reset 不碰判定参数；判定通路活性由探针
内建正控钉住：POST-RESET `talker=1`、FRESH-EQUALS `wantTalk=1`——若读数 0/0
反而红，与 392-r2 探针的 `nSil+nTalk==0 ⇒ MECHANISM MOVED` 同判；既有 67 枚
staged 名册逐位同＝判定语义未动的最强旁证）。

## 2. AC#2 hangover 清除读数（探针＝`internal/speech/vad_reset398_test.go`）

读数形状沿用 392-r2 机制探针的具名序 `(silence=N talker=M)`。材料＝本包实测
代理语声：talker 1500ms（校准史：900ms 过不了 onset 爬升＋MinSpeech 300ms，
首轮探针 `the probe talker did not open an utterance` 即此——按既有 1500ms
body 校准放宽）；gap 400ms（< 700ms 关门阈、全程落在 ~720ms hangover 期内，
`vad_test.go` 实测校准注记）。staged 跑（log `probe-tree-final.log`）：

| 臂 | 形状 | 读数 | rc |
|---|---|---|---|
| PRE-RESET（392-r2 对照形） | talker→`Finalize(nil)`→gap400→talker | **`(silence=1 talker=0)`** 假 start 落静音段 | 0 |
| POST-RESET | 同上＋`Finalize` 后 `Reset()` | **`(silence=0 talker=1)`** gap 干净、talker 正常开火 | 0 |
| FRESH-EQUALS | dirty(warmUp→Finalize→Reset) vs 全新 detector，三段计数 | **`(gapA=0 talker=1 gapB=0)`** 两臂全等 | 0 |

**口径差如实记**：派单写「reset 后 0/0」；实测 reset 后＝`(silence=0 talker=1)`。
第一个 0＝静音段假 start 清零（AC#2 主判据「post-speech 静音不再开火」）；
392-r2 读数形里的 `talker=0` 是 hangover 吃掉 endpointer 的形状，reset 后同样
清掉；而字面 0/0 若成立反而触发停手③（detector 被清死＝判定通路断）——
392-r2 探针自己的 Fatalf 就把 `nSil+nTalk==0` 定为 MECHANISM MOVED。本腿读数
以「silence=0（hangover 清）＋ talker=1（判定活着）」为准。

**FRESH-EQUALS 案的读数是落点敏感的按段计数**（gapA/talker/gapB 三段各数
start），不是 kind 序列——§3 的教训：kind 序列对落点不敏感，咬不住「start
从 talker 挪到 gap」的突变。旁证注记：FRESH-EQUALS 的 fresh 臂 gapA/gapB=0
本身是 `TestVadPureSilenceProducesNoEndpoint` 的窗级复刻正控。

## 3. AC#3 反向判据（三读数＋还原核）

突变＝仓外 `-overlay`（仓库零突变；突变件 `D:\tmp\wisp\20261101\398r1\mut\vad.mut.go`
由 `sed` 单行置换生成）：`d.vad.Reset()` → `_ = d.vad // MUTANT 398r1-M1`。

1. **摘 reset 调用 ⇒ hangover 复现**（突变 exe，log `mutant-probe2.log`）：
   POST-RESET 臂读数退回 **`(silence=1 talker=0)`**，rc=1，MUTANT 具名句
   「the hangover survived the reset」——**`d.vad.Reset()` 承重坐实**。
2. **等价案同突变复红**（log `mutant-fresh.log`）：`(gapA=1 talker=0 gapB=0)
   vs fresh (0 1 0)`，rc=1，具名句「the reset did not return the detector to
   its fresh state」——hangover start 挪落 gap 段被落点敏感读数咬住。
3. **定向自证**：同突变二进制旁案 `TestVadPauseDetail`/`TestVadPureSilenceProducesNoEndpoint`/
   `TestStateMachineIgnoresChunking` 全绿 rc=0——突变只打 hangover 形、不打
   既有判定通路。**开发史如实记**：等价案初版用 kind 序列对比，在突变 exe 下
   误绿（`[start end-of-stream]` 两臂同形——dirty 的 start 落 gap、fresh 的落
   talker，kind 序列不可分）；已改按段计数后同突变复红，单变量归因闭合。
4. **还原核**：①突变件可再生性——对仓内 `vad.go` 现跑同一 sed 再 diff 突变件
   ＝零差（「mutant-source-still-single-delta: yes」，突变面恰一行）；②仓内核
   ——`git status --porcelain -- internal/speech/`＝本腿恰两件（`M vad.go`＋
   `?? vad_reset398_test.go`）；③正树探针复跑 rc=0 全绿（`probe-tree-final.log`）。

## 4. AC#4 零回归

| # | 跑 | 形状 | rc | 读数 |
|---|---|---|---|---|
| 1 | internal/speech 全套 staged | 三 DLL 贴身＋两 WISP env＋`-test.v` | **0** | **RUN=69 / PASS=63 / FAIL=0 / SKIP=6**；6 SKIP＝punc 缺席族（与 387-r3 run 04 同名同形）；名册＝基线 67 枚逐位同＋本腿两枚新增绿（`Test398SileroHangoverReset`/`Test398ResetMatchesAFreshDetector`） |
| 2 | internal/speech 裸跑 | PATH 注入、无 WISP_*、`-test.v` | **0** | **PASS=25 / FAIL=0 / SKIP=44 / RUN=69**——387-r3 收敛的 rc0 全绿未破；新增 2 SKIP＝本腿两枚具名 skip（skip 句＝守卫具名错误「onnxruntime 1.17.1 is older than the 1.28…」） |
| 3 | cmd/wisp KWS 族 41a2r1 两案 staged | `wisp.test.exe`＋三 DLL 贴身＋两 WISP env＋cwd=cmd/wisp，`-test.run 'Test41a2r1WakeEndToEndArmedToListening\|Test41a2r1WarmGateFeedsKwsNotTheAsr'` | **0** | **2 PASS / 0 FAIL**（16.27s／14.75s；387-r3 run 05＝16.26s/13.86s 同带） |

**口径差如实记**：基线「61/0/6」「25/0/42」是 67 枚总体时代的数字；本腿落地后
总体 69 枚——staged 63 PASS（+2 新增绿）、裸跑 44 SKIP（+2 具名 skip）。
「⛔ 破」的判据落点是 **FAIL=0 且 rc=0 且既有名册逐位同**：三条全过（staged
SKIP 集逐名同、裸跑 25 PASS 逐位同基线）。

**门禁**：gofumpt `-l internal/speech/`＝空（rc 0）；`go vet ./internal/speech`
rc 0；d22scan clean rc=0——分母 bans#1-5 internal/=287、cmd/=52、ban#6
frontend/=116、ban#7 internal/tools/=25、ban#8 internal/=732（392-r2 在册 731
＋本腿探针件 1，逐源可数）、cmd/=183（392-r2 同值）。

## 5. 写面申报（commit pathspec 逐枚）

1. `M internal/speech/vad.go` —— 唯一代码件：`Detector.Reset()` 一枚（+13/−0，
   全在 Finalize 之后新增方法块；Accept/Finalize/Consume/Close/Endpointer 零字节）
2. `A internal/speech/vad_reset398_test.go` —— 探针件（AC#2 两案＋AC#3 反向读数，
   裸跑具名 skip 同守卫家族形状＝`newTestDetector`）
3. `A docs/evidence/s1/398r1-hangover-reset.md` —— 本件

## 6. 残余与边界（如实记账）

1. **owner 窗间清的新组合**：`Finalize(nil)` ＋ `Reset()`（零冷段等待）——本腿
   只交付面与读数，未改任何 owner（`internal/audio`/`cmd/wisp` 零字节）；接入
   属后续裁定（392-v3 §6 第 3 步 keepalive 翻转前置里做）。
2. **C++ 实现级语义未逐行核**（网络受限，§1）——行为级三面（探针/全套回归/
   突变反向）闭合；若验收腿要源码级核，需可通 GitHub 的环境。
3. **hangover 读数是 fixture 级**：材料＝1500ms 代理语声＋400ms gap；silero 对
   更长静音的语境衰减形状（「衰减多少够」）仍无扫参读数（392-v3 §7.4 原样）——
   reset 面让该问题失去 owner 侧紧迫性，但未关闭。
4. **测试内反向与突变反向并立**：主案 PRE-RESET 臂本身就是「摘 reset」的测试内
   形状（392-r2 控制臂同款）；仓外 overlay 突变是实现半承重的第二读数。两读数
   在本腿 logs 内各自独立在案。
5. **停手三条核销**见 §1——三条全部未触发，无停手交回。
