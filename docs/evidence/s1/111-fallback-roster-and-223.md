# 证据件：票 111 波 8b · 回退名册三路并集＋FIXTURE-GAP 防呆 ＋ 票 223 具名跳过

- **腿**：111-wave8b-roster-223（写腿，2026-11-13 晚）
- **派单**：`.scratch/wisp/probes/orch/20261113/txt/111-wave8b-roster-223-dispatch.txt`
- **对象 CI 读数**：run `43413718734`／`test-windows`／`cmd/wisp CLI tests`：
  `RUN=743 --- PASS=566 --- FAIL=2 --- SKIP=8`，两枚红：
  ① `Test41a2r1WakeEndToEndArmedToListening`
  （`resident_voice_41a2r1_windows_test.go:225: cannot find third_party\spike-models\matcha-icefall-zh-baker above the working directory`）；
  ② `TestTicket223PermissionDeniedSitsInItsOwnSentence`
  （`config_reload_perm_223_windows_test.go:95: attempt 1: icacls denied nothing, so this case cannot show the permission sentence`）。
- **写面**：`scripts/wisp-cli-staged-tests.sh`、`scripts/portable-tests.sh`（只加 ledger 行）、
  `cmd/wisp/config_reload_perm_223_windows_test.go`、本证据件、票 111 Progress log。
  ⛔ 未动 GUARD B 判据、thresholds、golden、allowlist、契约、`.github/**`、`build/models/manifest.json`、
  `internal/speech|tools|ball`；⛔ 未删测试、未放宽任何断言（223 的 Fatalf→Skip 只发生在
  「前提造不出」那一条支路上，能造出时全部断言照旧）。

## 1. 红一根因与修法：名册从手打 41 名正则改为三路并集现算

**根因（复算，非转述）**：41 枚名册是"逐**测试函数体** grep helper 名"建起来的，而
`Test41a2r1WakeEndToEndArmedToListening` 的函数体只调 `voice41a2r1Start`——
夹具触碰发生在那个**二次封装**构造器的体内（`voice15eRepoFile`／`voice41a2r1KwsDir`／
`voice15eInstallAsr`，`resident_voice_41a2r1_windows_test.go:115-147`）。函数级 grep 看不见
隔着封装的依赖 ⇒ 名册漏它 ⇒ 回退态它照跑 ⇒ 在无夹具的 runner 上红。
同型漏洞还有一枚**反向**的：`Test41a2r1GuardRowsWalkWithTheProductionFact` 在名册上但
**不依赖夹具**（它用 `voice41a2r1StubKws` 桩引擎走 D43 表；前腿无夹具读数 r1 里它
**PASS**，`after-nofixtures.log:4595`）——过跳同样在删覆盖，只是不红而已。

**修法（三路并集，全部长在 carrier 脚本里、各一处具名清单＋说理注释）**：

| 路 | 清单 | 具体成员 | 抓什么 |
|---|---|---|---|
| (a) | fixture helper/wrapper **函数名** | 7 枚 leaf：`voice15eVadModel` `voice15eSpikeModels` `voice15eRepoFile` `voice15eInstallAsr` `voice41a2r1KwsDir` `voiceWirePuncModel` `voiceWireTtsFiles`；8 枚二次封装：`voice41a2r1Start` `voice15eStartLeg` `voice15eArm` `voiceWireStart` `ttsInheritLeg` `newVoice26pRig` `voice26pRenderTtsWav` `voice26pSpikeAsr` | 函数体引用任一名字的用例（CI 那枚红走 wrapper 半路被抓住） |
| (b) | 夹具**路径字面量** | `spike-models` `model-fetch` `matcha-icefall` `silero_vad` `sherpa-onnx-punct` `sherpa-onnx-kws` `vocos` `asr-encoder` `test_wavs` | 夹具树改名/新增时只改这一处；按路径摸夹具的用例无 helper 也被抓住 |
| (c) | 人工残余 4 枚（不变） | `Test394R1Seam×2`、`TestSLO159…`、`TestAC4Focus…` | 源码 grep 原则上看不见的间接依赖（理由仍是前腿证据件 §4） |

**并集是现算的**：carrier 每次回退用一段 POSIX awk 逐函数扫 `cmd/wisp/*_test.go`
（纯 `index()` 字符串匹配，无动态正则），加 route-(c) 四枚拼成 `-test.skip` 全名锚定
表达式；`go test -list` 对编译二进制自验（匹配 0 枚即拒跑，前腿的绿 no-op 防线原样保留）。
新用例只要引用清单上任一名字，**第一次编译就自动进名册**。

**名册枚数（前/后）**：
- 前：**41** 枚（手打，`voice_case_pattern` 静态正则）。
- 后：**46** 枚（(a)∪(b) 现算 44 ＋ (c) 独有 2）。
- 净差（逐枚，机算 `comm` 对比，非手抄）：
  - 新增 6：`Test41a2r1WakeEndToEndArmedToListening`（本次 CI 红，route (a) wrapper 半）、
    `Test26PBargeinStopsWithin400ms`／`Test26PFeedOnlineMatchesWindowReplay`
    （经 `newVoice26pRig` → `voice15eVadModel`＋`voiceWireTtsFiles` 摸夹具；前腿 r1 里它们
    是设备门 SKIP，故 CI 从未红——依赖在 r1 读数旁证），
    `Test392AecTuningDoubleTalkSwing`／`Test392DisposeReleasesTheEngines`／
    `Test392KeepalivePostureFindingsAndColdDefault`（经 `newVoice26pRig`＋`voice26pRenderTtsWav` → 真	matcha 渲染）；
  - 落下 1：`Test41a2r1GuardRowsWalkWithTheProductionFact`（stub 路径，无夹具 PASS 实证）。

**刻意不收进 (a) 的名字与理由**（过跳＝静默删覆盖，清单必须精确）：
`voice15eWaitFor`／`voice15eWriteWav`／`voice15eDataDir`／`voice15eWriteConfig`
（纯测试工具，无夹具触碰；只用它们的用例在无夹具读数里 PASS）、
`voice41a2r1Walked`／`voice41a2r1Feed`／`voice41a2r1PushSilence`／`voice41a2r1Wav*`
（操作用例已有数据；feed/wav 只经 Start 构造器到达）。
曾试过两个更宽的口径都被证据否掉：前缀家族扫描把 `Test392JSONSentinel…`／
`Test381r2Veto…` 两枚**纯单元**（r1 无夹具 PASS）误捕；函数级不动点闭包在真实语料上
爆炸（645 枚，方法名子串级联污染）——两条路都弃，取逐名清单。

## 2. FIXTURE-GAP 防呆（红一要求的"不让同类再漏"机制）

回退态跑完后，carrier 重读自己的日志：凡 `--- FAIL` 的失败消息含夹具缺失签名
（`cannot find third_party`＝voice15eRepoFile 的 walk-up 失败；`WISP_SPIKE_MODELS is not set`
＝voice15eInstallAsrFiles 的守卫），单列 **FIXTURE-GAP** 判读：逐枚打印
`依赖夹具但不在名册上：<名>`，carrier 以非零退出（1），**⛔ 不混进普通 FAIL 计数**
（四数如实——那枚确实跑了确实失败；GAP 标签加在四数之上，不替换它们）。
归因形状＝消息行在其 `=== RUN` 与 `--- FAIL` 之间（go -v 的输出形状），归不上号的 FAIL
保持普通 FAIL（防呆的失败模式是**少标签**，绝不会是**假绿**）。
**GUARD B 职权未动**：包级结果行照旧先打（`ok/FAIL <import path> <时长>`），每次运行都有。

## 3. 门禁（逐条如实）

| # | 门 | 读数 |
|---|---|---|
| 1 | `bash -n scripts/wisp-cli-staged-tests.sh` | **rc=0**（2026-11-13 20:2x 现跑） |
| 2 | `bash -n scripts/portable-tests.sh` | **rc=0**（2026-11-13 20:2x 现跑） |
| 3 | `go build ./cmd/wisp/`（CGO_ENABLED=1＋mingw64 gcc 在 PATH） | **rc=0** |
| 4 | `go vet ./cmd/wisp/` | **rc=0**（零输出） |
| 5 | `sh scripts/d22scan.sh` | **rc=0 clean**（`no D22 ban violations`；ban#8 cmd/=276、internal/=1005，与前腿基线同量级） |

## 4. 核心验证（junction 沙盒，复刻 hosted runner＝有 DLL 无 voice 模型）

沙盒＝`D:\tmp\wisp\20261113\111-wave8b\sandbox\`（只建不删）：`cmd/internal/tools/frontend/
build/forks` 与 `third_party/sherpa-onnx`、`third_party/model-fetch` 用 junction 指回真仓，
**`third_party/spike-models` 不挂**（＝无 voice 模型），`go.mod/go.sum/deps.toml` 拷贝，
`scripts/` 四件拷贝；全链 `wisp-cli-tests.sh`→`portable-tests.sh --scope=cli`→
`wisp-cli-staged-tests.sh`（与 CI 完全同链）。

### 4.1 验证①：回退态全链（`core1-nofixtures.log`，ship carrier，20:30–20:36）

- **名册打印**：`NAMED FALLBACK - fixtures absent, so 46 fixture-dependent cases
  (computed roster 44 from routes (a)+(b) ... go test -list = 46)`——三路并集现算、
  对编译二进制自验通过；缺的四个 env 逐枚点名。
- **`Test41a2r1WakeEndToEndArmedToListening` 不再执行**：`=== RUN` 计数 **0**（grep 实证）。
- **四数（portable-tests.sh 逐字）**：`=== RUN=408  --- PASS=318  --- FAIL=1  --- SKIP=0`。
- **退出码＝139**（测试二进制中途段错误）：真窗族机器耦合偶发——与 r1（同树全量 796 枚
  跑完）互证非本腿改动面；crash 点 `TestAC4PriorFocusSurvivesARefusedPanelSample` 与
  前腿 r3 同位。**GUARD B 职权在场**：carrier 打了包级 own-line
  `FAIL  github.com/CarlosShao/wisp/cmd/wisp	365.00s`，portable 判 `FAIL (own line)`。
- **唯一 FAIL＝`TestTicket255RosterStillMatchesTheActualReadSites`**：沙盒非 git 仓伪影
  （其扫描器走 git 面），**与前腿 r3 逐字同款**（r3 日志同报 `config_receipt_255_test.go:309
  ... scan=[] match=false`），非本腿产物。
- **FIXTURE-GAP 防呆正确沉默**：那枚 FAIL 无夹具签名 ⇒ 防呆块不打印、不误标；
  用船上 awk 原样抽件复扫整份日志＝空输出。
- 覆盖面注：408 枚是段错误截断前的分母（r1 全量 796 枚跑完证明包可全跑；r3 同截断形状）。

### 4.2 验证②：FIXTURE-GAP 正控（`gap-pos-control.log`，20:40–20:42）

**变异**：沙盒 carrier 副本把 `voice41a2r1Start` 从 route-(a) 清单删掉（模拟"名册再漏一枚"），
`-test.run '^Test41a2r1WakeEndToEndArmedToListening$'` 单跑。逐字读数：

```
wisp-cli-staged-tests.sh: NAMED FALLBACK - fixtures absent, so 44 fixture-dependent cases
  (computed roster 42 from routes (a)+(b), ... go test -list = 44) are
=== RUN   Test41a2r1WakeEndToEndArmedToListening
    resident_voice_41a2r1_windows_test.go:225: install the streaming ASR into
      C:\Users\swq\AppData\Local\Temp\...\001: WISP_SPIKE_MODELS is not set
--- FAIL: Test41a2r1WakeEndToEndArmedToListening (0.02s)
FAIL
FAIL  github.com/CarlosShao/wisp/cmd/wisp	0.00s
wisp-cli-staged-tests.sh: FIXTURE-GAP - 1 failed case(s) failed on the fixture-absence
wisp-cli-staged-tests.sh:   signature, not on an assertion: this is a ROSTER GAP, not a package defect.
wisp-cli-staged-tests.sh:   Add these to the route-(a)/(b)/(c) lists above in this same script:
	依赖夹具但不在名册上：Test41a2r1WakeEndToEndArmedToListening
wisp-cli-staged-tests.sh:   FIXTURE-GAP carrier verdict: non-zero exit (1); do not read this run as a plain FAIL count.
```

**carrier rc=1**。⇒ 名册再漏时红的是「名册缺项」（带名可修），不是"测试坏了"；
签名命中的是第二枚（`WISP_SPIKE_MODELS is not set`），第一枚（`cannot find third_party`）
由合成探针覆盖（§4.4）。

### 4.3 验证③：派单字面场景——调用方把 `WISP_SPIKE_MODELS` 指到不存在目录
（`core3-callerset-absent.log`，20:42–20:44，ship carrier 还原后）

```
wisp-cli-staged-tests.sh: WISP_SPIKE_MODELS=D:\tmp\wisp\20261113\111-wave8b\nonexistent-models is SET but not on disk (streaming ASR + VAD tree).
wisp-cli-staged-tests.sh: missing fixture set: WISP_SPIKE_MODELS=D:\...nonexistent-models(caller-set but not on disk: ...)
wisp-cli-staged-tests.sh: NAMED FALLBACK - fixtures absent, so 46 fixture-dependent cases ...
```

- 名册照打（46 枚）、`Test41a2r1WakeEndToEndArmedToListening` **0 次执行**（`=== RUN` 0）；
  二进制应答 `no tests to run`（`-test.run` 选中那枚、`-test.skip` 跳掉它——两个旗子都生效）；
  包级行 `ok  github.com/CarlosShao/wisp/cmd/wisp	0.00s`、**rc=0**（0 跑 0 红，ok 如实；
  GUARD B 的锚行在场）。旧语义下这一发会是无名册的 30 枚 missing-model 红——语义修正
  的行为差有此对照。
- 全链四数读数由 4.1 提供（单枚直调不经 portable-tests.sh，无四数行——如实注明）。

### 4.4 验证④：归因 awk 合成日志单元探针（`synthetic-gap.log`＋`ship-gap.awk`）

从**船上脚本**逐字节抽取归因 awk，喂合成 -v 日志（5 枚：PASS 邻居、带
`cannot find third_party` 签名的 FAIL、无签名普通 FAIL、SKIP、把签名词句印在
PASS 用例里的诱饵）。读数：**恰好点名 `TestRosterMissedVoiceCase`**——
PASS/SKIP/普通 FAIL/诱饵全部不误伤。核心验证 #1 后又复跑一遍（awk 随脚本两轮编辑
漂移后重抽）＝同读数。

## 5. 红二修法：`TestTicket223…` 的环境判读从 Fatalf 改为具名 Skip

判据一字未换（仍是"造得出拒绝就断言"）：

```diff
 			if _, err := os.ReadFile(path); err == nil {
-				t.Fatalf("attempt %d: icacls denied nothing, so this case cannot show the permission sentence", attempt)
+				// TICKET 111 WAVE 8B: this is an environment verdict, not an
+				// assertion one. ... (理由注释，含 run 43413718734 锚)
+				t.Skip("ticket 223: this runner cannot deny the read (icacls is a no-op for an administrator), so the permission sentence cannot be produced here")
 			}
```

- **不是永远跳**：只有"stat 成功＋read 失败"这一正控**失败**（＝deny 没落上）才走到 Skip；
  deny 落上了（普通用户会话、机主本机）正控通过，**后面每一行断言照旧执行**。
- **台账行**（`scripts/portable-tests.sh` ledger，格式 `name|pkg|platform|class|reason`）：
  `TestTicket223PermissionDeniedSitsInItsOwnSentence|./cmd/wisp/|windows|fixture|…`
  （票 71 AC#3 的 SKIP 必须有台账行；理由串含 run 号、机制、恢复条件）。
- 派单的"ledger 加行"与"⛔ 不许永远跳过"同时满足：Skip 分支只在正控失败时到达。

## 6. 边界与残余（如实）

- **CI 首跑读数欠**（零 push）：名册现算、FIXTURE-GAP、223 Skip 的 CI 颜色都要等下一次推送。
- 本件 §4 的沙盒读数为**本机**读数；`go test ./cmd/wisp/` 全量本机跑不了（票 98 豁免），
  走 carrier 链，未跑未绿的事一律如实标注。
- `TestTicket223R2FailureSentenceRouting` 等其它 223 族用例未动（它们的 r1 无夹具读数
  PASS，不依赖 icacls 前提）。

## 7. 撤销口令

回「**撤 111-wave8b**」⇒ carrier 名册回 41 枚静态正则＋`voice_env` 语义回"caller 值无条件赢"＋
FIXTURE-GAP 块删除；portable-tests.sh ledger 撤最后一行；223 测试的 Skip 支路回 t.Fatalf；
本件与票面 Progress 段保留（append-only）。
