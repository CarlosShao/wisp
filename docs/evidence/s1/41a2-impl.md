# 票 41 引擎腿 `41-a2`（写腿裁决表）— KWS 引擎落 `internal/speech`，挂上已冻结的 `WakeWordEngine` 缝

> **尺**：开工 `HEAD=d36059b44`（`dev`，索引空、`internal/{speech,config,risk}`／`cmd/wisp` 无未提交脏）→ 收工 HEAD 已被兄弟腿推进到 `a2bfb4ae8`。现跑 2026-10-27 17:41–19:0x +08。★读数逐条现跑，⛔ 转述。
> **写面（三件，逐名）**：`internal/speech/kws.go`（新，12,909 B）＋ `internal/speech/kws_test.go`（新，22,245 B）＋ `internal/speech/slot.go`（**+41 行/−0 行**，纯增量）。`docs/evidence/s1/41a2-impl.md`（本件）。临时件只在 `D:\tmp\wisp\20261026\41a2\**`（仓外，只建不删）。
> **⛔ 未碰**：`cmd/**`、`internal/risk/**`、`internal/config/**`、`internal/session/**`、`internal/speech/cer/**`、`testdata/**`、`frontend/**`、`design/**`、`scripts/**`、`.github/**`、`engines.go`（**C9 冻结零字节**——编译期断言 `var _ WakeWordEngine = (*KwsEngine)(nil)` 放在 `kws.go:105`）。

## 1. 票面 6 格逐条（AC 原文逐字）＋ 本腿覆盖判

| # | 票面原文（逐字） | 判 |
|---|---|---|
| AC1 | `Armed RSS/CPU gate measured (sampler) within D32 row; wake latency ≤300ms.` | 归 A2 接线腿（需 sampler＋真唤醒路径）⇒ 本腿不覆盖 |
| AC2 | `Keywords: custom keyword list change → reload without restart; hit on fixture audio.` | ★**"hit on fixture audio"半格本腿覆盖**（§3 七发夹具读数）；"reload without restart"的接线半格归 A2 |
| AC3 | `Veto: KWS running + Confirming + veto-word wav → L1 call cancelled; KWS absent → UI message asserted.` | 归 A2＋B 腿；本腿已把 `Veto` 机制落好（`kwsVetoWords` 表驱动，现空表——出货词表无 veto 词，`t.Logf` 级零 veto 断言过） |
| AC4 | `P14 false-trigger report (thresholds table + recommended values) committed.` | 归 A2/B；本腿交了它的原料：常量 `DefaultKwsScore=3.0`／`DefaultKwsThreshold=0.1` 的实测依据（§2） |
| AC5 | `Mute cycle Armed→Muted→Armed with inference stop/start asserted (no model reload).` | 归 A2；本腿对应件＝`Slot.LoadKws/Kws()`（同形）＋引擎 `Close` 幂等 |
| AC6 | `Opt-in off → zero mic opens (device-handle probe) + Sleeping budget.` | 归 A2（设备开关裁定未写清，派单明令本腿不碰 `cmd/wisp`） |

## 2. 构造成功（staging；模型路径与配置逐字段）

- **模型**：`WISP_SPIKE_MODELS=D:\work\workspace\projects plans\Wisp\third_party\spike-models` ⇒ `kwsModelDir()` 解析出 `…\sherpa-onnx-kws-zipformer-wenetspeech-3.3M-2024-01-01`。**五件 sha256 逐位 match=True／size_ok=True**（vs `models/manifest.json` kws 条目 archive.files）：encoder/decoder/joiner（int8）＋tokens＋keywords 五发布尔比对全 True，`third_party\model-fetch\kws.tar.bz2` 顶层 sha256 同 True。manifest：`status:"ok"`、`license Apache-2.0`、`quant int8`（README front-matter 逐字 `license: Apache License 2.0`）。
- **sherpa 收到的 KeywordSpotterConfig 逐字段**（探针 `TestKwsProbeSherpaDebug` 的 Debug 打印逐字抓得）：`FeatConfig{16000, 80}`；`Transducer{encoder/decoder/joiner=三枚 epoch-99-avg-1-chunk-16-left-64.int8.onnx}`；`Tokens=tokens.txt`；`NumThreads=1`（IntraOpThreads 钉）；`Provider="cpu"`；`MaxActivePaths=4`；`KeywordsFile=keywords.txt 或 test_keywords.txt`；`KeywordsScore`/`KeywordsThreshold` 见下。**`modeling_unit="cjkchar"`（sherpa 默认）照单全收**——引擎能跑通并命中（§3），改建模单元不在本腿权限内，如实记。
- **搜索两常量＝实测定的，⛔ 不是拍的**：先按 API 默认 `1.0/0.25` 跑 → 对 checkpoint 自带 test_wavs **全 0 命中**（28 发参数全扫探针：`1.0/0.25`、`3.0/0.1`、demo 表、`3.0/0.05-b0` 四组 × 7 wav）；**sherpa c-api.h:1663-1691 对同一族 checkpoint 的官方示例就是 `keywords_score=3.0f`、`keywords_threshold=0.1f`**，按它跑 → 正样本全部命中、负样本保持静默（§3/§4）。⇒ `DefaultKwsScore=3.0`、`DefaultKwsThreshold=0.1` 落为包内常量（⛔ 未加配置键）。
- **构造失败＝具名错误、⛔ 崩**：`NewKeywordSpotter` 返 nil ⇒ `NewKwsEngine` 转具名 error（`15-b checkOnnxruntime()` 前置拒绝同款，老 ORT 红句见 §6）；半配置（缺任一路径）5 发全拒（`TestKwsEngineRejectsIncompleteConfig`）。

## 3. test_wavs 逐个 KeywordHit（★标题读数；对照 test_keywords.txt）

引擎＝`NewKwsEngine`（keywords 指向 **demo 表** `test_wavs\test_keywords.txt`——**配对事实**：test_wavs 是 demo 表的录音，c-api.h 官方示例对同族 checkpoint 就把 `keywords_file` 指到它；生产表 `keywords.txt` 的 8 个唤醒词没有任何一条夹具录音说过 ⇒ 对照面见 §4）。喂法＝生产形状 512 样本窗逐窗 `Accept`。

| wav | 样本数 | KeywordHit.Word（逐字序） | At（样本钟） | test_keywords.txt 对照 |
|---|---|---|---|---|
| 0.wav | 89,784 | **[]（0 hit）** | — | 负样本（官方集合自身行为） |
| 1.wav | 82,449 | **[]（0 hit）** | — | 负样本 |
| 2.wav | 72,389 | **[]（0 hit）** | — | 负样本 |
| 3.wav | 128,481 | `[文森特卡索 法国]` | 373,103（第二词同流内再命中） | ✔ 文森特卡索（+法国） |
| 4.wav | 73,586 | `[蒋友伯 女儿]` | 261,163 起第二遍流 | ✔ 蒋友伯 ✔ 女儿 |
| 5.wav | 66,447 | `[周望军 落实]` | 513,136 流尾 | ✔ 周望军 ✔ 落实 |
| 6.wav | 56,737 | `[朱丽楠 见面会]` | 569,873 流尾 | ✔ 朱丽楠 ✔ 见面会 |

- demo 表 8 词 × 全部命中；`Veto=false` 逐词断言（demo 表无 veto 词）；`At` 全在 `(0, fed]`（样本钟，D42#9）。命中序被**逐字钉死**（`reflect.DeepEqual`）：搜索确定（dither=0／单线程／固定模型音频）。
- **0-2 是官方负样本这一点是量出来的**：参数全扫 28 发里 0/1/2.wav 在任何参数下都 0 hit。
- 停手⑤复核：test_wavs **跑得出 hit**（3-6.wav 共 8 词），⛔ 未硬凑；0 hit 的三枚是集合自己的负样本。

## 4. 正控（给样本数）

| 控 | 夹具 | 读数 |
|---|---|---|
| 纯静音 | 3s 数字静音＝**48,000 样本** | **0 hit**（`TestKwsEngineOnSilenceStaysSilent`） |
| 非关键词语音 | 合成语音替代物 5s＝**80,000 样本** | **0 hit**（`TestKwsEngineOnNonKeywordSpeechStaysSilent`） |
| 生产表对照面 | keywords.txt（8 唤醒词）× 7 夹具（614,884 样本总量） | **7 发全 0 hit**（`TestKwsEngineProductionTableRejectsTheDemoWavs`）——没人说过这些词，引擎正确拒绝 |

## 5. 反向判据：摘"命中后必 Reset" ⇒ 必须红（三读数＋还原核）

变异＝`k.go` 的 `k.spotter.Reset(k.stream)` 一行替换为注释，**只存在于仓外** `D:\tmp\wisp\20261026\41a2\mut\kws-noreset.go`，经 `go test -c -overlay mut\overlay.json` 编译（工作树零突变窗口）。

1. **变异红**：`speech.mut.exe -test.run TestKws` ⇒ **rc=1**，红句逐字：`kws_test.go:369: 3.wav hits = [文森特卡索], want exactly [文森特卡索 法国] over 373103 samples: … a missing second keyword is the sherpa Reset clause broken`——摘 Reset 后 3.wav 丢同流第二词，即绑定 `sherpa_onnx.go:2466`（"You MUST call it right after detecting a keyword"）的可观测后果。
2. **还原绿**：正树 `speech.test.exe` 整包 staging ⇒ **rc=0**（RUN=67／PASS=61／FAIL=0／SKIP=6，§6）；`TestKwsEngineHitsTheFixtureWavs` 在其中，3.wav=`[文森特卡索 法国]`。
3. **原样绿（第二遍）**：`run-full3.log` 再发整包 ⇒ **rc=0** 同四数。
- **还原核（时机＝最后一发变异之后）**：真件 sha256 `0B5BFAFF…F58BAEF` 与变异后同值；文本级反证——树件**有** `k.spotter.Reset(k.stream)`（kws.go:286）、**无**变异标记；mut 件（sha256 `10DB5337…`）**无** Reset、**有**标记。⇒ 工作树从未进突变窗口。
- 附：`TestKwsEngineResetAfterHitLetsTheNextUtteranceHitAgain`（同夹具两遍喂、第二遍仍 `[蒋友伯 女儿]`、`At` 只前进）是这颗钉子的行为读数；裸 sherpa 探针（不 Reset）同流只见 6 词 vs 引擎 8 词，同向佐证。

## 6. staging 跑法完整命令 ＋ "为什么 PATH 不够"（三层口径）

```
# 1) 编译测试二进制，并把三枚 sherpa DLL 放到它自己旁边
go test -c -o D:\tmp\wisp\20261026\41a2\bin\speech.test.exe ./internal/speech
copy third_party\sherpa-onnx\*.dll D:\tmp\wisp\20261026\41a2\bin
# 2) 模型路径（staging 目录里没有 third_party）
set WISP_SPIKE_MODELS=D:\work\workspace\projects plans\Wisp\third_party\spike-models
set WISP_SILERO_MODEL=<同目录>\silero_vad.onnx
# 3) 就地跑（cwd = 二进制所在目录；包内常量路径兜底也指向同一目录）
D:\tmp\wisp\20261026\41a2\bin\speech.test.exe -test.v     ⇒  RUN=67 PASS=61 FAIL=0 SKIP=6，rc=0
```
- **为什么 `PATH` 打不过 `C:\Windows\system32`**：Windows 对依赖 DLL 的搜索序＝「主模块所在目录 → System32 → … → `PATH`」；`go test` 把二进制建在临时目录 ⇒ 加载器走到 System32 拿到 **1.17.1**，`PATH` 排在其后改不了序；1.17 的 `OrtGetApi(28)` 返 NULL ⇒ `0xc0000005` 崩。staging 把 **1.28.2**（third_party 实测 FileVersion）放进"主模块所在目录"——唯一排得比 System32 前的位子。System32 那枚本机现量＝`1.17.260311-1434.1…`。
- **三层口径本腿实测**：①裸跑无 DLL＝`0xc0000135`（DLL_NOT_FOUND，15c-r1 同录）；②**只给 `PATH`（无 staging）＝rc=1，不崩**——`checkOnnxruntime()` 前置拒绝，红句逐字：`NewKwsEngine: speech: onnxruntime 1.17.1 is older than the 1.28 the linked sherpa-onnx C API requires; a second onnxruntime.dll is shadowing the one on PATH (on Windows the loader searches System32 before PATH)`（本腿 KWS 测试直接吃这套保护，无 skip）；③staging＝全绿 rc=0（上表）。

## 7. 四数 ＋ audio 零字节 ＋ cmd/wisp 零字节

| 面 | 命令 | 读数 |
|---|---|---|
| `internal/speech` 改前 | `git ls-tree -r -l HEAD -- internal/speech` | **30 文件 ／ 251,952 B**（与派单逐字同） |
| `internal/speech` 改后 | 同上＋盘上 | **32 文件 ／ 288,189 B**（＋2 文件 kws.go/kws_test.go；slot.go +41/−0） |
| staged | `git diff --cached --name-status` | 开工时 **空**（0/0/0）；本腿交付前不 stage 别人面 |
| staging 全包 | §6 | **rc=0，RUN=67／PASS=61／FAIL=0／SKIP=6** |
| `internal/audio` | `git status --porcelain -- internal/audio`＝**空**（零字节）；`go test -count=1 ./internal/audio` | **rc=0**，`ok … 16.183s`（独立跑） |
| `cmd/wisp` | 本腿 pathspec 逐名三件，`cmd/` 命中＝0；`git show --name-only` 复核 | **零字节**（⚠ 开工时 `cmd/wisp` 无脏；证据期出现 A2 接线腿在飞脏件——`resident_voice_*.go` 3 `M`＋1 `??`，触者非本腿，本腿零碰零清零背） |

## 8. 格式甲乙两尺 ＋ d22scan（第 21 处定式）

| 尺 | 树 | 输出 |
|---|---|---|
| 甲 `gofmt -l ./internal ./cmd ./tools` | 导出树（HEAD＋本腿三件） | **0 行**，rc=0 |
| 乙 `gofumpt.exe -l`（同范围） | 同上 | **0 行**，rc=0 |
| 甲＋乙 | 共享树现况 | 恰 **3 条既存假红**＝派单点名的 `internal\agent\approval\pending_read.go`／`internal\agent\tools.go`／`cmd\wisp\models.go`（CRLF 假红，⛔ 本腿责任）；按 `[\\/]` 归属：`internal\speech` 命中 **0** |
| `sh scripts/d22scan.sh` | **`git clone --shared` 导出树**（HEAD `a2bfb4ae8` 原样，status 0 行） | **rc=0 clean**；正控逐字 `runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77` |
| 同上 | 导出树＋本腿三件（唯一差异） | **rc=0 clean**；`bans#1-5 internal/ 263→264`、`ban#8 internal/ 682→684`（增量即本腿两枚新件）、cmd/=158 不变 |
| 同上 | 共享树现况 | rc=1 两红，**逐条归因不在本腿写面**：①`cmd/wisp/resident_voice_windows.go:311 [phantom-citation]`＝在飞接线腿的未提交改；②`design/` 39 vs 导出 14＝机主外部 agent 的未跟踪文件。⛔ 本腿不碰不代改 |

## 9. 边界与停手核对（逐枚）

- 四条硬边界：`cmd/**` 0 字节✔；`engines.go` 零改（`git status` 无它）✔／`internal/risk`、`internal/config`、`internal/session`、`internal/speech/cer` 零碰✔；`testdata/frontend/design/scripts/.github` 零碰✔；零配置键（阈值/词表全吃包内常量）、零 `t.Skip`、零放宽断言（Skip=6 全是**既存**用例的具名环境 skip，非本腿新增）✔。
- 停手五枚：①五件 sha256 match=True ⇒ 不触发；②绑定与普查逐字符（`KeywordSpotterConfig`/`NewKeywordSpotter` 失败返 nil/`:2466` MUST Reset 全在）⇒ 不触发；③`engines.go` 未改 ⇒ 不触发；④未碰 `cmd/wisp` ⇒ 不触发；⑤test_wavs 出 hit（8 词）⇒ 不触发。
- 派单写法遵守：`★` 等符号只在本文档；`.go` 内注释全 ASCII（d22scan rc=0 即仪器读数）；无裸 `go func(`、无墙钟超时（`Accept` 纯同步，`At` 全样本钟）。

## 10. 交回父代理（一句话＋给 A2 的注记）

**引擎腿落地**：`KwsEngine` 挂上已冻结的 `WakeWordEngine` 缝（编译期断言在 kws.go）、`slot.go` 同形增 `FamilyKws/Kws()/LoadKws`、staging 整包 rc=0（67/61/0/6）、夹具读数＝demo 表配对 8 词全中＋生产表对照 0 命中＋静音/非关键词 0 hit、摘 `Reset` ⇒ 红（同流丢第二词）、还原核过、d22scan（导出树净/带我件）双 rc=0 clean。**给 A2**：①阈值常量 `Score=3.0/Threshold=0.1` 有实测依据，换值请带 P14 读数；②`kwsVetoWords` 表驱动已就位，B 腿改表即分类；③`WISP_SPIKE_MODELS` 双口径（spike-models 根或 kws 目录）都被 `DefaultKwsConfig` 接受；④接线时喂 512 窗即可（`Accept` 任意分块）；⑤demo 配对读数在 `TestKwsEngineHitsTheFixtureWavs`，勿改成"必须每夹具命中"——0-2 是官方负样本。
