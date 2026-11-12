# cmd/wisp 整包终态 — 从 52 红到 14 红（2026-11-12 round 88，编排者亲跑）

## 0. 命令与台面

- 锚点：`HEAD=4d6189dc2`（分支 dev）
- 命令：`bash -c "scripts/wisp-cli-staged-tests.sh"`（carrier 自解析全部夹具环境）
- 日志：`D:\tmp\wisp\20261112\fullclean\cmdwisp.log`
- **读数：rc=1，`--- PASS`=593 / `--- FAIL`=14 / `--- SKIP`=9**

## 1. 与上一轮整包（2026-11-12 15:0x，`230ac5-self/full.log`）对比

| 指标 | 上一轮（裸 go test） | 本轮（carrier） |
|---|---|---|
| 逐包 | `internal/*` 35 包全 ok；`cmd/wisp` **FAIL 749s** | `cmd/wisp` **PASS=593 / FAIL=14** |
| `cmd/wisp` FAIL 数 | **52** | **14** |
| `onnxruntime 1.17` 遮蔽 | **39 处** | **0** |
| `WISP_SPIKE_MODELS is not set` | 多枚 | **0** |
| `checkpoint is incomplete` | 多枚 | **0** |
| `标点级未接上` | 有 | **0** |
| PASS 数 | 极少（多数 0.00s 即死） | **593** |

## 2. 四类环境错误精确归零（逐字 grep 计数）

| 签名 | 计数 |
|---|---|
| `WISP_SPIKE_MODELS is not set` | **0** |
| `onnxruntime 1.17` | **0** |
| `标点级未接上` | **0** |
| `checkpoint is incomplete` | **0** |

⚠ **计数陷阱（留档）**：日志里 `incomplete` 字样出现 150 次，全部是
`msg="scope disposed" ... disposal_incomplete=false` 这个**正常字段**，不是错误。
按 `incomplete` 裸词 grep 会得出"150 处夹具缺失"的错误结论——必须用
`checkpoint is incomplete` 这个完整短语。

## 3. 剩余 14 枚 FAIL 逐名归因

### 3a. 台账 `:32512` 已具名的存量真红（5 枚）

| 测试 | 断言/性质 |
|---|---|
| `TestTicket198FirstRunCallerIsTheRunEntryOnly` | `firstrun_198_test.go:281: AC J1 RED: ensureFirstRunConfig is called by [runTextTask runConfigSet65], want exactly [runTextTask]` |
| `Test313WritableUnionEqualsTheTwoRosters` | `panel_config_313_test.go:252: the door reports the union [...]` 与两名册不等 |
| `TestAC4FocusReturnToPriorWindowGap33r5`（本轮未现） | 需真桌面焦点 |
| `TestLockedSuffixNeverTakesTheLockItself33r11` | `panel_locked_naming_33r11_windows_test.go:254: panelHolder.stampPublishLocked carries the "Locked" suffix but no roster entry` |
| `TestAC247HandingTheLevelToTheSeamIsNotVisibility` | `resident_audio_247_windows_test.go:355: prototypeVisuals is ON in this assembly: the shipped default is off` |

### 3b. 已知偶发（1 枚，有台账 A/B 记录）

`TestTicket223HandEditedFsLooseningCostsAnL2Card`
- 台账 `:28515` 逐字：`★TestTicket223HandEditedFsLooseningCostsAnL2Card 60 发 A/B：★base 30 发 3 红 / mine 60 发 2 红、同红行 :436/:439`
- 我本轮**单独复跑：PASS 3.77s**
- ⇒ **偶发，⛔ 不算真红**

### 3c. 环境层修好后暴露的真红（8 枚）

| 测试 | 已知断言 | 状态 |
|---|---|---|
| `TestAC1ResidentLegInstallsItsLogListenerOnDisk` | 待取证（腿 `resident-ac1-verdict` 在飞） | 新暴露 |
| `TestAC1ResidentLegOutlivesItsOwnLogFailure` | 同上 | 新暴露 |
| `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink` | 同上 | 新暴露 |
| `TestWireR1WarmIdleDisposesOnTheOwner` | `resident_voice_wire_r1_windows_test.go:745: mic close counted 1 times, want at least 2 (the Warm row and the cascade's defence in depth)` | 已具名（`voice-fixture-env.md` §8） |
| `Test26PWithoutAECTheEchoTranscribesVerbatim` | 待取证 | 新暴露 |
| `Test384ResidentWatchdogEndToEnd` | 待取证 | 新暴露 |
| `Test12r2SeamCallSitesAreTheMeasurement` | 待取证 | 新暴露 |
| `TestTaskOutputPointerNotice174r8SpeaksTheRealReason` | 待取证 | 新暴露 |
| `Test365RunRowsAStuckTaskAsStuck` | 待取证 | 新暴露 |

⚠ 这些**不是环境问题**（四类签名已全零），是真实断言。⛔ 把它们当环境红一笔带过。

## 4. 结论（可核）

1. **本轮的治法有效**：`cmd/wisp` 从 **52 红 / 几乎零 PASS** 变成 **593 PASS / 14 FAIL**，四类环境错误全部归零。
2. **剩下 14 枚里**：5 枚台账已具名存量真红、1 枚已知偶发（复跑 PASS）、**8 枚是环境修好后暴露的真红**。
3. **⛔ 不宣称"整包全绿"**。距全绿还差那 8 枚新暴露的真红，需逐枚取证（`resident-ac1-verdict` 腿正在做 3 枚）。
4. **零 Go 产码改动**：本轮全部改动在 `scripts/**` 与 `docs/**`。
