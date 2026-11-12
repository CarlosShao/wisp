# Round 88 汇总 — 环境三层全部打通，暴露 8 枚真红（2026-11-12）

## 0. 一句话

上一轮我报"52 枚红、根因是模型未下载要花钱"是**错的**。本轮把环境三层逐个打通：
**①DLL 遮蔽（39 处→0）②夹具环境变量（20 枚→真跑）③shell 选错（WSL relay）**，
整包从 **52 红 / 几乎零 PASS** 变成 **593 PASS / 14 FAIL**，且四类环境错误签名**精确归零**。

## 1. 三层环境，逐层归零（可复跑）

| 层 | 签名 | 修前 | 修后 | 治法 |
|---|---|---|---|---|
| ①DLL 遮蔽 | `onnxruntime 1.17` | 39 处 | **0** | `scripts/wisp-cli-staged-tests.sh` 把 deps.toml 钉的 DLL 摆到测试二进制同目录（Windows 加载序 exe 目录 > System32 > PATH） |
| ②夹具寻址 | `WISP_SPIKE_MODELS is not set` | 多枚 | **0** | carrier 内 `voice_env()` 五变量逐个默认＋逐目录验存在 |
| ②夹具形状 | `checkpoint is incomplete` | 多枚 | **0** | `WISP_TTS_DIR` 指向 `spike-models/matcha-icefall-zh-baker` 而非 `model-fetch/x-matcha` |
| ②夹具形状 | `标点级未接上` | 每发必现 | **0** | `WISP_PUNC_MODEL` 指向**文件** `.../model.int8.onnx` 而非目录 |
| ③shell | WSL relay 静默零测试 | 未知 | **exit 2 响亮拒绝** | carrier 开头查 `OSTYPE` 必须 msys/cygwin |

## 2. 整包读数对比

| | 上一轮（裸 go test，15:0x） | 本轮（carrier，19:0x） |
|---|---|---|
| `cmd/wisp` | FAIL 749s，**52 枚红** | **PASS=593 / FAIL=14 / SKIP=9** |
| 日志 | `D:\tmp\wisp\20261112\230ac5-self\full.log` | `D:\tmp\wisp\20261112\fullclean\cmdwisp.log` |

## 3. 剩余 14 枚，逐名归因（⛔ 不混为一谈）

### 3a. 台账 `:32512` 已具名存量真红（3 枚本轮出现）
`TestTicket198FirstRunCallerIsTheRunEntryOnly` · `Test313WritableUnionEqualsTheTwoRosters` · `TestLockedSuffixNeverTakesTheLockItself33r11` · `TestAC247HandingTheLevelToTheSeamIsNotVisibility`

### 3b. 已知偶发（1 枚，复跑转绿）
`TestTicket223HandEditedFsLooseningCostsAnL2Card`——台账 `:28515` 记 60 发 A/B（base 30 发 3 红）；我单独复跑 **PASS 3.77s**。

### 3c. 环境修好后暴露的真红（8 枚，⛔ 不是环境问题）

| 测试 | 逐字断言 | 定性 |
|---|---|---|
| `TestAC1ResidentLegInstallsItsLogListenerOnDisk` | `:498: the child never reached the event loop, so the disk read above proves little` | 预算不够（10s 窗口 vs 本机 3.4–5.6s 启动，大头 `speech.NewDetector` 1.83s） |
| `TestAC1ResidentLegOutlivesItsOwnLogFailure` | `:539: the refusal stopped the app: the child never reached the event loop` | **偶发**（腿跑 5 发：3 PASS / 2 FAIL），`pollUntil127(20,...)` 1s 预算撞 2.9–5.6s 启动 |
| `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink` | `:635: the last record in the sink is "wisp: D38(e) step 10 (exit): ok", want a shutdown-step record` | **★断言落后于产码**（见 §4） |
| `TestWireR1WarmIdleDisposesOnTheOwner` | `:745: mic close counted 1 times, want at least 2` | 产品缺陷（Warm 行＋级联防御应各关一次麦） |
| `Test26PWithoutAECTheEchoTranscribesVerbatim` | `:216: MUTANT 26p-M2 SHAPE: the raw echo opened zero episodes` | 待深查 |
| `Test384ResidentWatchdogEndToEnd` | `:230: 384: never observed the tree back under the Armed cap (a real within-cap reading) within 3s` | 待深查 |
| `Test12r2SeamCallSitesAreTheMeasurement` | `:193: 12r2: openScope no longer stamps T1 on the scope's SCOPE-SETTLING announcement` | 待深查 |
| `Test365RunRowsAStuckTaskAsStuck` | （日志未抓到 `_test.go:` 断言行，需细上下文） | 待深查 |

## 4. ★最硬的一条发现：`TestAC1ResidentLegBooks…` 是「断言落后于产码」

**机制（腿现读，可核）**：
- 断言要求**盘上最后一条记录**是 `shutdown step ...`（`resident_sink_nail_127_windows_test.go:635`）。
- 但 `wisp: D38(e) step N ...` 由 `cmd/wisp/resident_windows.go:246` 在 **defer 尾部把十步全量回放**；
- `shutdown step ...` 由 `internal/proc/shutdown.go:130` 在**序列执行当时**打印；
- ⇒ 十步回放**必然排在**序列当时那些记录**之后** ⇒ **末位永远是 `D38(e) step 10`**，`shutdown step` 不可能在末位。
- **这是票 43（`5056a365c`，2026-10-31）引进的**：它在 defer 里新加了 `:246` 这个回放循环，**断言没跟上**。
- **旁证极硬**：`t.Logf` 自己逐字说 `RESIDENT LEG RECORDED ON DISK: 5 shutdown record(s), steps [1 2 5 6 7]`——**5 条 shutdown 记录确实在盘上**，`len(trail) != 0` 那关（`:629`）**过了**；红的只是"它在不在末位"。
- **稳定性**：现跑 **5/5 稳定红**（7.80s／6.96s／6.79s／5.82s／5.97s），voice 与 novoice 都红，**时长/负载/环境不参与**。

**判据（⛔ 本件未改，只陈述）**：把 `last` 改成"`trail` 的最后一条必须是 `recs` 里最后一条 `shutdown step` 记录"，或改成"`trail` 非空且步序单调"即绿。**⚠ 这是改断言，需单独立票并由非实现者翻框。**

## 5. 本轮我自己的重大失误（已回退）

票 174 的 AC#4：我把 `mem.ArtifactsDir()` 无条件加进 `[fs]` 授权根，**违反票面 `Status:` 行逐字禁区**（"本票不碰 `[fs]` 的默认值"），且**与 AC#2 方向相反**（AC#2 要求回执说"读不到"，我让它"读得到"）。
- **抓住我的测试**：`TestTaskOutputPointerNotice174r8SpeaksTheRealReason`，单独复跑 **FAIL 4.13s**，稳定真红。
- **回退**：`7dcb145ff`；回退后该族 **PASS=2 / FAIL=0**（1.93s／1.55s），`go build` rc=0。
- **台账**：`A1800`；票 174 的 Progress `:120` 已就地更正。
- **教训**：红变绿⛔等于修对；要问"它原来测的那个状态还存在吗"。我那枚三判据测试件**全绿却测的是错方向**。

## 6. 提交清单（本轮，全部已推 origin/dev）

| 提交 | 内容 |
|---|---|
| `32ef54dcd` | carrier 解析五个夹具环境变量 |
| `946355c13` | `WISP_PUNC_MODEL` 指向文件而非目录 |
| `b7697197f` | A1798 更正"模型未下载"的错误结论 |
| `4d6189dc2` | A1799 记录三类环境红归零 |
| `0703cc60f` | 整包终态证据件 |
| `7dcb145ff` | **回退我自己的 174 AC#4** |
| `238835c3e` | A1800 记录我的失误 |
| `dd8286a6b` | 更正票 174 Progress `:120` |
| `192821996` | carrier 拒绝非 msys shell（WSL relay） |
| `829c6da1c` | 腿 `resident-ac1-verdict` 的三枚取证 |

## 7. 边界（诚实）

- ⛔ **不宣称整包全绿**：距全绿还差 **8 枚真红**（§3c），需逐枚立票／并入现有票。
- ⛔ 未下载任何模型、未改 System32、未提交二进制、未动 `PLAN.md`／`specs/**`／`thresholds.go`／golden／`allowlist.txt`／`frontend/**`。
- ⛔ 本轮的 8 枚真红**一枚都没修**——只做了归因。
