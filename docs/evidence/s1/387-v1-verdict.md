# 票 387 终裁（腿 387-v1，非实现者）——五格复核与收口判定

- 被裁件：`9d9607559`（387-r2，3 枚 +97：internal/speech/tts.go +6 ／ kws_test.go +8 ／ docs/evidence/s1/387r2-impl.md +83）。
- 三把尺：提交真在（git show 现跑）／件真在（tts.go blob 4dc2fe521..52b1aed3f 与 diff 头一致；导出件 SHA256=3F4C7DA8...A354 与开工时 repo 工作树逐位同）／全部读数自己复跑（本件各节均出自本腿在导出树的现跑，未采信实现件任何字面读数）。
- 改前基线：`9d9607559^` = `1b9721519`；与开工基线 `46a80b579` 在 internal/speech 零差异（git diff --stat 空）。
- 方法：全部实验出自仓外导出树 `D:\tmp\wisp\20261031\387v1\`（after=git archive 9d9607559，18202 文件；before=git archive 9d9607559^，18201 文件，差 1 枚=新增证据件，与三文件 diff 吻合；mutant=robocopy 后真改 tts.go）。third_party 用 junction 指主仓（0 枚 .go，不进生产分母）。
- 纪律：GOTMPDIR 每跑置空；rc 一律重定向现取；裸跑=PATH 前置 third_party\sherpa-onnx 且六个 WISP_* 全清（现验初值为空）；staging=-c 二进制＋三 DLL 同目录＋两 env（各指本导出树资产）。⛔ 全程未动工作树在飞件（391-r1=tools/mockllm；387-r3=internal/speech/kws_test.go 在本腿末段出现 M——零碰零归因）；repo 侧 tts.go SHA256 开工/中段/收尾三次同值。

## AC#1 守卫补齐 —— 判：站得住（不能翻）

- 四先例同位同形逐一对核（本腿现读现行代码，非引实现件）：NewDetector vad.go:72-74、NewStreamingAsr asr.go:71-73、NewKwsEngine kws.go:213-215、NewCTTransformerPunctuator punctuate.go:185-187——全部是「config 校验后、sherpa 构造前」插同句形三行。tts.go 新守卫 :142-144 同位（config switch :134-141 之后、sherpa.NewOfflineTts :164 之前），文档句与 asr/kws 家族同形；checkOnnxruntime 本体未动（vad.go:102-114，requiredORTMinor=28）。
- 全包 sherpa.New/Get 直调点盘点（本腿自数，grep 现跑，11 处）：vad.go:79（守卫 :72 之后）；vad.go:103（守卫自身的版本探测，全部裸跑证据显示 1.17.1 可答不崩）；asr.go:74（:71 之后）；asr.go:125 NewOnlineStream（rec 只能出自守卫构造器，Begin :119 nil 早退）；kws.go:216（:213 之后）；kws.go:222（spotter 出自守卫构造器）；punctuate.go:191（:185 之后）；tts.go:164（新守卫 :142 之后）；kws_test.go:510/515 探针（新 test-side skip :490-492 之后）；cgo_boundary_windows_test.go:111 rawVAD（唯一调用点 :238 在 test-side 守卫 :225 之后，TestSingleOnnxruntimeInstance leg2 亦在 :157 守卫之后）。补完后 100% 被守卫覆盖——独立复核成立。

## AC#2 裸跑不崩 —— 判：站得住（不能翻），带残余①

- 本腿复跑（export-after，无 staging）：rc=1（失败非崩溃）、25 PASS／6 FAIL／36 SKIP、全 log「Exception 0x」计数=0——与实现件读数同形；67=25+6+36 与 staging 总数一致。
- 能崩调用点全部具名 skip，skip 句逐字含守卫具名错误（8 处：tts_test.go:78/135/171/206/240/319、slot_real_test.go:28、kws_test.go:491 探针新句），错误文本=「speech: onnxruntime 1.17.1 is older than the 1.28 the linked sherpa-onnx C API requires; a second onnxruntime.dll is shadowing the one on PATH (on Windows the loader searches System32 before PATH)」。
- 6 FAIL 恰为 6 枚 KWS 引擎测试（HitsTheFixtureWavs／ProductionTableRejectsTheDemoWavs／OnSilenceStaysSilent／OnNonKeywordSpeechStaysSilent／ResetAfterHitLetsTheNextUtteranceHitAgain／CloseIsIdempotent）——测试直调 NewKwsEngine 后 t.Fatalf 呈现守卫拒绝：具名失败、不崩不哑。skip 化超出 AC#1 字面（残余①，387-r3 在飞修，不记入本格）。

## AC#3 反向判据 —— 判：站得住（不能翻）

- A380 WalkDir 方法属实：mutant=仓外导出树 robocopy 后**真改盘上 tts.go**（非 -overlay），摘 :142-144 三行守卫。
- 三读数：改前 SHA256=3F4C7DA8...A354（与 repo 逐位同）；改后 SHA256=F9FEECC3...F678、守卫计数 0——**与 387-r2 证据件记录的 mutant SHA 逐字相同**（同一单变量编辑的独立复现，反证实现腿真用真改文件法）；repo 侧全程未动。
- 崩溃复现（mutant 裸跑 -run Test26）：「Exception 0xc0000005 0x0 0x18 ...」，栈逐字 = _Cfunc_SherpaOnnxCreateOfflineTts <- NewOfflineTts(sherpa-onnx-go-windows@v1.13.8/sherpa_onnx.go:1359) <- NewMatchaTts <- tts.go:161 <- slot_real_test.go:27 <- (*Slot).LoadTts <- Test26TheRealSwitchDoesNotAddTheTwoFamilies；ORT 自报机理行逐字（「The requested API version [28] is not available, only API versions [1, 17] ... Current ORT Version is: 1.17.1」）；161=164-3 与「只少守卫三行」算术闭合。
- 还原核：repo tts.go 拷回 => SHA256 复等 3F4C7DA8...A354、守卫计数回 1、复跑 rc=0（4 PASS/7 SKIP/0 FAIL/0 Exception）。

## AC#4 staging 零回归 —— 判：站得住（不能翻）

- 形状：两树各 `go test -c` ＋ 三 DLL（onnxruntime 1.28.2／sherpa-onnx-c-api／sherpa-onnx-cxx-api）同目录 ＋ WISP_SILERO_MODEL/WISP_SPIKE_MODELS ＋ -test.v -test.timeout 30m；涉 sherpa 全部 staged（A1326 卡；本票 AC#4 集合=三 DLL＋两 env）。
- 改前基线（before 树）：rc=0、61 PASS/0 FAIL/6 SKIP；改后（after 树）：rc=0、61/0/6——逐位同。
- 名册：两树各提 67 行（结果+测试名），diff=0 行；TestKwsProbeSherpaDebug 两册都在且 staging PASS；matcha 真合成在跑（Test26TheTtsEngineEmitsRealAudio 计时行在册）；两 log Exception 计数=0。
- 复跑陷阱具名（给后续腿）：TestDefaultKwsConfigResolution 用 runtime.Caller 从编译进二进制的源路径上溯——二进制编自哪棵树，两个 WISP_* env 就必须指同一棵树；跨树指 env 会让该测试 1 枚具名 FAIL（本腿首跑曾中，属 staging 形状错、非回归）。

## AC#5 纪律格 —— 判：站得住（不能翻）

- git show --stat：写面恰 3 枚（tts.go／kws_test.go／证据件）；零触碰 cmd/wisp、internal/{risk,session,config,memory,panel,audio}、thresholds/golden/testdata/票面/台账/specs/PLAN/frontend/design。未 push（推送归编排者）。
- 本腿自己：repo 内只写本件；⛔ push/amend/reset/rebase/stash/checkout 全程未犯；实验全部仓外。

## 归因校正核（清单第 5 条）—— 两处都真，成立

- 改前基线全包裸跑：rc=1，「Exception 0xc0000005」栈 = _Cfunc_SherpaOnnxCreateKeywordSpotter <- TestKwsProbeSherpaDebug（kws_test.go:502 直调）；崩前 5 枚 KWS Fatalf 已 FAIL；全 log TTS 测试零提及 => 全包首发崩点=kws 探针，TTS 没轮到跑。
- 改前基线定向 -run Test26：栈逐字同 AC#3 链但行号=改前 tts.go:158（158+3+3=164，与改后行号算术闭合）；ORT 机理行在同 => NewMatchaTts 崩点独立坐实。
- 结论：59-r1「NewMatchaTts 崩」在其 aecspike 包视角为真＋全包首发=kws 探针——两处都真、非误归因；A1335 台账修订与事实一致。

## d22scan（清单第 6 条）—— 净＋正控过＋两尺

- 尺一（件真在）：tools/d22scan/d22scan.exe 在盘（4,365,312 字节）。
- 导出树 clean：rc=0；bans#1-5 internal/=287 生产 .go、ban#8 internal/=729 .go（与实现件同分母）；导出树无 .git，扫描器自报「不适用 gitignore、全路径扫」=更严方向。
- 尺二（自己复跑）＋正控：向导出树**生产面**塞裸 goroutine（无 owner/recover）=> rc=1、「1 finding(s)」、分母 287 变 288；撤件复扫 rc=0 回 287/729。注意：裸 goroutine 禁令豁免 _test.go（本腿第一发正控件放 _test.go 未被咬、仅分母 +1）——后续腿做正控必须用生产面文件。
- 证据件小瑕疵（不阻塞）：387r2-impl.md 引的 07-d22scan.log 实为一次 go run 失败残留（main module does not contain package tools/d22scan），真读数来自预编译 exe；判语本身成立，记录备查。

## 残余定性与收口判定

- 残余①：6 枚 KWS 引擎测试裸跑 Fatalf 非 skip——测试断言形状、超出 AC#1 字面授权；387-r3 已派在飞（A1335；本腿收尾时工作树已现 M kws_test.go=其件）=> 不阻塞本票。
- 残余②：台账 :28956「NewMatchaTts 崩」表述对全包裸跑不全 => 编排者已办（A1335 append-only 修订为「首发崩点=kws 探针直调＋NewMatchaTts 定向另证、两处都真」）=> 已闭合。
- 残余③（本腿新具名、轻微）：证据件 d22scan 日志为 go-run 失败残留＋正控豁免面未写明 => 不阻塞，供后续腿绕坑。

## 结论

- 五格全站得住（零翻框）；**票 387 可带残余结案**，残余①②③如上（①387-r3 在飞、②已闭合、③轻微备查）。
- 复跑入口：D:\tmp\wisp\20261031\387v1\（logs\10-15=裸跑与 mutant 三读数、21/23+roster-*=staging、30-34=d22scan；export-after/before/mutant 三树在盘）。
