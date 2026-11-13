# 证据件：票 111 · cmd/wisp CLI 测试的 voice 夹具（test-windows 步 7）

- **腿**：111-cmdwisp-cli-fixtures（写腿，2026-11-13）
- **派单**：`.scratch/wisp/probes/orch/20261113/txt/111-cmdwisp-cli-fixtures-dispatch.txt`
- **现象（权威件引用）**：CI run `43389487057`／job `test-windows`／步 7：hosted runner 上没有
  `third_party/spike-models` ⇒ carrier `wisp-cli-staged-tests.sh` 的 voice_env 判 setup failure
  直接 `exit 1` ⇒ `portable-tests.sh` 四数 `RUN=0 PASS=0 FAIL=0 SKIP=0` ⇒ GUARD B
  `1 of 1 packages ... <NO TOP-LEVEL RESULT LINE>` ⇒ 步红＝「包消失＝红」。
- **根因**：`third_party/` 整体被 `.gitignore:12` 忽略；runner 上它只经 actions/cache +
  build.ps1 恢复，而那两处只补 **sherpa DLL**，从不补 **voice 模型**。carrier 假定「夹具
  已在盘上」（其注释原话 "the bytes are already on disk under third_party/"）——在 hosted
  runner 上这个前提第一次不成立。
- **写面**：`scripts/wisp-cli-staged-tests.sh`、`.github/workflows/ci.yml`（仅 test-windows
  相关步）、新增 `scripts/fetch-voice-fixtures.ps1`、票 111 Progress log、本证据件。
  ⛔ 未动 `cmd/wisp/**`、`internal/**`、`portable-tests.sh`、`tools/d22scan/**`、thresholds、
  golden、allowlist。

## 0. 路线裁决：**A 为主干 ＋ B 作 carrier 内具名回退**（两半都落地）

**为什么 A 是主干（最有价值）**：
1. **主流 harness 形状**：CI 里「数据依赖在测试前 fetch + 哈希验证 + actions/cache」是
   GitHub Actions 的标准形状（本仓自己的 `deps.toml`→`fetch-deps.ps1`→cache 就是这个
   形状；voice 模型此前只是从没被接进同一根管子）。voice 测试在 CI **真跑**，不是被
   记账缺席——cmd/wisp 的 41 枚夹具依赖用例从「CI 上 0 次」变「每次 push 都被执行」。
2. **单一源不新增**：URL 与 SHA256 全部**现读** `build/models/manifest.json`（`wisp models
   ensure` 的同一张 pin 账本；夹具测试自己也经 `WISP_MODELS_MANIFEST` 读它），fetch 脚本
   零重复 pin——manifest 改 pin ⇒ fetch 自动跟随，无第二真相源。
3. **cache 现实可行**：targets 约 650 MB 原始字节（KWS 压缩包 32.7 MB、matcha 压缩包
   75.5 MB、VAD 0.64 MB、ASR 三件 237.2 MB、vocoder 53.9 MB、punc 压缩包 64.7 MB），
   actions/cache 命中后 fetch 步退化为纯哈希校验（秒级）。首次冷跑约多付
   （体积/带宽）的下载时间，上限由步级 `timeout-minutes` 兜住。
4. **fetch 可靠性地基现量过**：本机六目标全部与 manifest pin **逐字同哈希**
   （2026-11-13 现算；silero_vad 的 pin 在 manifest :56，同样吻合）——pin 账本与盘上
   夹具一致，fetch 后测试拿到的字节就是 pin 的字节。

**为什么仍要 B（而不是只做 A）**：A 修的是「runner 没夹具」的**常态**；但当步真的
失效（manifest 某 URL 全死、runner 断网、开发者本机从未 fetch）时，旧 carrier 的答案
是「整包 0 跑」——那正是 GUARD B 当天咬中的形状。B 让失效日的答案变成
「**具名缩小 + 全额记账**」：夹具依赖用例按**全名**跳过（41 枚＝29 直接依赖＋4 间接依赖，
2026-11-13 现量名单在 carrier 内并对着编译二进制自验；首轮探明的 file 级方案已被证伪，
见 §3.4），包内其余用例照跑（r1 证 796 枚全走完），carrier 打出**真实的包级结果行**。
这不是把「包消失」判绿——GUARD B 依旧行使「跑了 0 次＝红」的职权；B 只是让
「跑了一个具名的子集」成为**不同于 0 跑**的第三种、诚实的读数。回退也响亮：
每次回退打 NAMED FALLBACK 段（缺哪些变量、跳过多少枚、怎么恢复）。
`WISP_CLI_VOICE_FORCE_SETUP_FAILURE=1` 可一键还原旧行为（撤销口令在票面）。

**为什么不选「只 B」**：voice 41 枚将在 CI 永久缺席，等于把夹具依赖面从 CI 分母里
永久删掉——那是 AC#2 禁的「把被检对象从门禁里删掉」的缓释版。A 买到的是 41 枚的
CI 执行面，B 买到的是失效日的诚实降级；两半互补，B 不替代 A。

## 1. 改动清单（真落盘）

| 件 | 改动 |
|---|---|
| `scripts/wisp-cli-staged-tests.sh` | voice_env 不再 exit 1：缺目标**收集**进 `voice_missing`；全解析后二选一——全在⇒全量跑；有缺⇒`-test.skip='^(41 枚全名\|...)$'` 具名回退（名单先经 `go test -list` 对编译二进制自验，匹配 0 枚即拒跑），`WISP_CLI_VOICE_FORCE_SETUP_FAILURE=1` 还原旧硬退。新增**包级结果行**（`ok/FAIL <import path> <时长>`，由二进制真实退出码+实测时长推导）。GUARD（staging/onnxruntime/WSL/GOOS）逐字未动 |
| `scripts/fetch-voice-fixtures.ps1`（新） | 读 manifest，fetch 五目标（VAD 文件、ASR 三件、matcha 包+vocoder、KWS 包、punc 包），逐文件 sha256 验证后才落位；幂等（先验后载）；exit 0/1 |
| `.github/workflows/ci.yml` | test-windows 内 **ADDED STEP ONLY** ×2：`Cache voice fixtures (manifest key)`（deps.toml 键旁边的第二 cache 键，deps.toml 变更不再连带丢弃模型缓存）＋ `Fetch voice fixtures`（在 cgo build 之后、CLI 步之前）；CLI 步追加 `timeout-minutes: 60`（hang bound：步内含冷夹具下载上限）与 fixtures-leg 注释段。既有步**零移动/零编辑/零删** |

## 2. 门禁（逐项，如实）

| # | 门 | 读数 |
|---|---|---|
| 1 | `bash -n scripts/wisp-cli-staged-tests.sh` | **rc=0**（2026-11-13 17:1x 现跑） |
| 2 | fetch 脚本 PowerShell Parser.ParseFile | **0 errors**；实跑 rc=0（本机缓存命中形状：六目标全 `ok (cached, sha256 verified)`，零下载） |
| 3 | `ci.yml` yaml 解析（`py -3 -c yaml.safe_load`） | **OK / rc=0**；`git diff --numstat` ci.yml ＝ **79 增 / 0 删**（纯加步，无既有行被动） |
| 4 | `sh scripts/d22scan.sh` | 见 §5（后台跑，收数后补） |

## 3. 改前/改后本机模拟（junction 沙盒，⛔ 未动真树）

**沙盒形状**（`D:\tmp\wisp\20261113\111-cmdwisp\sandbox\`，只建不删）：
`cmd internal tools frontend build forks` 与 `third_party/sherpa-onnx`、`third_party/model-fetch`
用 **junction** 指回真仓（runner 视角=有 DLL、无 voice 模型；`spike-models` 目录**不**挂），
`go.mod/go.sum/deps.toml` 拷贝，`scripts/` 四件拷贝（carrier 按 HEAD/新版两态替换）。
全链走 `wisp-cli-tests.sh`→`portable-tests.sh --scope=cli`→`wisp-cli-staged-tests.sh`
（与 CI 完全同链）。

### 3.1 改前（HEAD 版 carrier）

```
wisp-cli-staged-tests.sh: WISP_SPIKE_MODELS defaults to .../third_party/spike-models, which is not on disk.
wisp-cli-staged-tests.sh: the voice fixtures were never fetched on this host, so the
wisp-cli-staged-tests.sh: voice tests cannot be judged here. Treat this as a setup
portable-tests.sh: four numbers (all from -v output): === RUN=0  --- PASS=0  --- FAIL=0  --- SKIP=0
portable-tests.sh:   <NO TOP-LEVEL RESULT LINE> missing github.com/CarlosShao/wisp/cmd/wisp
portable-tests.sh: GUARD B - 1 of 1 packages in scope mode=cli printed
portable-tests.sh:   NO top-level result line, so this step tested them ZERO times:
portable-tests.sh: strict runner exited 1 for scope=[./cmd/wisp/]
```
**rc=1**；与 CI run 43389487057 步 7 逐字同形（沙盒复现成立）。

### 3.2 改后·无夹具（新 carrier 的回退态）

〔见 §3.4：两发读数，首发的 24 红暴露了 file:: 形 skip 无效——机制被当场证伪并重做〕

### 3.3 改后·有夹具（真树探针 + CI 形状推断）

- **真树单用例探针**（2026-11-13 18:4x，新版 carrier `-test.run ^Test15eSilenceNeverCallsTheHandOff$`，夹具在位）：
  `voice fixture environment resolved (caller values win)` → `--- PASS: Test15eSilenceNeverCallsTheHandOff (1.87s)` →
  carrier 包级行 `ok  github.com/CarlosShao/wisp/cmd/wisp	2.00s` → **rc=0**。零 skip、真跑、GUARD B 锚行成立。
- **全量 566 枚的 CI 形状推断**：夹具齐时与既有 CI 全量跑同形（AC#4 历史读数 328.9–392.9s、
  RUN=335 起；回退态 412 枚跑通证明载体链路无新阻碍）——首发 CI 颜色欠（零 push）。
- **fetch 幂等预检**（本机=缓存命中形状）：六目标全 `ok (cached, sha256 verified)`，rc=0，零下载；
  冷缓存路径＝逐 URL 下载→包级 sha256→OS tar 解包→逐成员 sha256→才落位（机制在
  `fetch-voice-fixtures.ps1`，真冷路径的实测欠 CI 首跑）。

### 3.4 两发回退读数（机制自纠错的全过程，如实）

**r1（file:: 形 skip——被证伪的首次版本）**：`=== RUN=796 PASS=583 FAIL=24 SKIP=0`，
rc=1。24 红里 **21 枚是 voice 族** ⇒ `-test.skip` 的 `file::` 形**静默匹配零个**（探针复证：
单用例加 `file::` skip 后照跑照红；同用例改名字形 skip 后 0 行执行）。判定＝机制坏，
**当场重做**为 41 枚全名锚定形 + `-test.list` 名单防呆。r1 的 24 红成分＝21 voice（机制坏）
＋2 枚 markseam_394（间接依赖，file 级 grep 盲区→进 r2 名单）＋TestSLO159（同类→进名单）
＋TestTicket255Roster…（沙盒非 git 仓伪影，其扫描器走 git 面）
＋TestAC4Focus（真窗族已知红，重试后绿）。

**r2/r3（41 枚全名 skip 形——交付版）**：r2 卡死 1 枚真窗族用例（单测探针 1.17s PASS=
桌面时序偶发，与 r1 同树 796 枚跑完互证）⇒ 重跑；r3（干净一发）＝
`=== RUN=412  --- PASS=323  --- FAIL=1  --- SKIP=0`，rc=139（测试二进制中途段错误=
真窗族机器耦合偶发，r1 全量跑完互证非本腿改动面；411 枚完成枚里 41 枚夹具用例
**0 执行**、唯一 FAIL=上述沙盒伪影）。carrier 输出链逐字：
```
wisp-cli-staged-tests.sh: NAMED FALLBACK - fixtures absent, so 41 fixture-dependent cases
wisp-cli-staged-tests.sh:   (roster counts 41; sources (a)+(b) named in
wisp-cli-staged-tests.sh:   docs/evidence/s1/111-cmdwisp-cli-fixtures.md and ticket 111) are
wisp-cli-staged-tests.sh:   skipped BY NAME this run; every other case in ./cmd/wisp/ still runs.
...
FAIL  github.com/CarlosShao/wisp/cmd/wisp	331.00s
portable-tests.sh: four numbers (all from -v output): === RUN=412  --- PASS=323  --- FAIL=1  --- SKIP=0
portable-tests.sh:   FAIL (own line) github.com/CarlosShao/wisp/cmd/wisp
portable-tests.sh: strict runner exited 139 for scope=[./cmd/wisp/]
```
⇒ **GUARD B 不再判包消失**（包有 own-line，FAIL 状态如实传导），rc=二进制真实退出码。

**r3 的 RUN=412 与 r1 的 RUN=796 差值说明**（如实，防误读）：两发的分母不可直接比——
r3 崩在 331s 处（真窗族段错误）后测试进程终止，剩余约 380 枚未跑；r1 跑完了 796 枚。
r1 的 796 里含 41 枚本应被跳过的夹具用例（机制坏）＋583 PASS；r3 的 412 里夹具用例
0 执行。**权威读数＝r3 的 skip 有效性（41 枚 0 执行）+ r1 的全包可跑性（796 枚全走完）**，
两者合起来覆盖「回退态=具名缩小+其余全跑」的两半。

## 4. 事实清单（fetch 目标与依赖面，全部现量）

- **依赖面**：`go test -list '.*'`（编译二进制）@2026-11-13 树 = **617** 枚顶层用例
  （windows 档）；winlive 档的文件级探查（-list -tags winlive = 632）⇒ 15 个夹具依赖
  测试文件全部是普通 `windows` 档、无一藏在 winlive 后。
- **夹具依赖用例＝41 枚**（skip 名单的构成，两路并集）：
  - **(a) 直接依赖 29 枚**：函数体引用夹具 helper（voice15eVadModel/voice15eSpikeModels/
    voice41a2r1KwsDir/voiceWirePuncModel/voiceWireTtsFiles/voice15eRepoFile/
    voice15eInstallAsr）的逐函数 grep 映射；分布＝15e×5、41a2r1×2、acting_395r1×1、
    bargein_26p_echo×3、interrupt_381r2×3、interrupt_390r1×3、selfecho_26r3×2、
    ttsinherit_382wr2×0（三枚经共享 helper 间接依赖，看 (b)）、wire_r1×6、settle_trace×1、
    e2e_demo_16r1×2、config_wire_180k×2、config_wire_423×2。同文件的其余 28 枚
    （如 180k 的 8 枚纯配置用例）**不进名单、回退态照跑**。
  - **(b) 间接依赖 4 枚**（无夹具读数 r1 的 24 红暴露，grep 盲区）：markseam_394×2
    （驱动产线 voice plan 解析器）、SLO159（断言 395r1 文件的在册读数）、AC4Focus
    （其 hidden-window 前台样本只在 voice run 内产出）。
  - 名单防呆：carrier 每次回退先 `-test.list <同表达式>` 对编译二进制验证，
    匹配 0 枚即拒跑（绿 no-op 防线，票 71 AC#3）。
- **fetch 目标⇒carrier env**：`vad-silero`→WISP_SILERO_MODEL；`asr-streaming-paraformer-zh-en`
  三件（**dev 拼名** asr-encoder/decoder.int8.onnx+asr-tokens.txt）→WISP_SPIKE_MODELS；
  `tts-matcha-zh-baker` 压缩包→WISP_TTS_DIR、`vocos-22khz-univ.onnx`→树根（TTS vocoder
  一级回退点）；`kws-...-2024-01-01` 包→WISP_KWS_DIR；`punc-...-int8` 包的
  `model.int8.onnx`→WISP_PUNC_MODEL。dest 拼名=夹具测试自己的回退路径逐字
  （resident_voice_15e/41a2r1/wire_r1_windows_test.go）。
- **验证**：每文件 sha256 对 manifest pin（archive：包级 pin+解压后逐成员 pin 双层）；
  本机预检六目标**全 MATCH**。

## 5. d22scan

**rc=0 clean**（2026-11-13 17:2x 现跑）：`no D22 ban violations`；ban #8 scope 面
internal/=1005、cmd/=276 Go files（comments+_test.go included）——与 r5 基线同量级。

## 6. 边界与残余（如实）

- **portable-tests.sh 未在写面**：GUARD B 的 mode=cli 表 grep 的是 go test 的 `ok <pkg>` 行，
  而 staged carrier 直跑二进制**从不产**这行——即便夹具齐了，旧 carrier 下 GUARD B 仍会
  红（speech 层当年为此在 portable-tests.sh 里换过 GUARD B 形状，cli 层没换）。本腿的解法
  是在 carrier（写面内）按二进制真实退出码+实测时长打包级结果行（setup refusal 路径
  绝不打 ok 行）。**正解**应是 speech 形的 cli 臂 GUARD B，归 portable-tests.sh 的主人
  （票 409/111 的 GUARD 维护腿），已具名进票 111 Progress log。
- **`timeout-minutes: 60` 与 job 级 `timeout-minutes: 30`**：两者不冲突（步上限≤job 上限）；
  但 job 30 是否该随新步职责上调，属「改 job 预算」——编排者拍板项，本腿未动、已具名。
- **CI 侧首跑读数（夹具真跑、回退真退）欠**：零 push ⇒ 无真实 run；本件全部为本机读数。
- **portable-tests.sh 的 ledger 是按用例名记账的仪器**：41 枚具名跳过的记账落在
  票 111 Progress log + 本件 §4 名单（(a)/(b) 构成），并在 carrier 输出里逐次打印
  （NAMED FALLBACK 段 + `-test.list` 验数）——不是静默跳过；「ledger 行」这一形状
  本腿造不了（会改写面外文件），残余如实具名。

## 7. 撤销口令

- 回「**撤 111 夹具腿**」⇒ carrier 回 HEAD 版（setup-failure 形）、ci.yml 三处还原、
  fetch 脚本删除、票面 Progress 段保留（append-only）。
- carrier 内另有 `WISP_CLI_VOICE_FORCE_SETUP_FAILURE=1`：不改盘上件，单次运行还原
  旧的「夹具缺失=setup failure 硬退」行为。
