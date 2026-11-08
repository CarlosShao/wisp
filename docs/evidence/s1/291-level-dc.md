# 票 291 裁决表（读程腿 `291-mic`）— 电平尺不剥直流：AC#2 四形真麦读数 ＋ AC#3 落点表 ＋ AC#4 门禁

**性质**：只读程交付（票面 AC#3 明写「⛔ 本票不许落地」⇒ 零产码改动）。**实现者＝本腿；裁决权在编排者与非实现者验收腿**。
**钟点**：2026-11-08 16:49–17:02 +08。**树锚**：主仓 `8b568eccd`。**读数全文**：`.scratch/wisp/probes/291/mic/00-four-shape-readings.md`（台件源 `mic/resident_audio_291_probe_windows_test.go`，只存于共享克隆树 `/d/tmp/wisp/20261108/291-mic-tree`，从未提交）。

## AC#2 四形真麦读数（同一枚 `麦克风 (Realtek(R) Audio)`，probe＝247 live 台件家族同一 harness，16:57:34–16:58:13）

| 形 | mean | max | 判读 |
|---|---|---|---|
| ① 完全安静 3s | **0.003067** | **0.099908** | 见事实 2 |
| ③ 已知幅度音调 6s（源理论电平 0.258950） | **0.001321** | **0.002530** | 见事实 3 |
| ② 真人说话窗 | **未采** | — | 见事实 1 |
| （无人窗 25s，⛔ 不是形②） | 0.001490 | 0.012336 | harness 活体证据 |
| 对照安静 3s | 0.001305 | 0.001696 | 稳定 |
| ④ 合成 DC(12345)＋正弦(6000) 帧 | **0.398368**（闭式 0.398367，delta 4.369e-7） | — | 尺的算术钉死 |
| ④' 合成纯 DC 12345 帧 | **0.376740** | — | ＝票面夹具值逐位一致 |

冒烟前置：`TestLiveWasapiSmoke` PASS（25 帧/834ms，0 丢帧）；两把 tasklist 尺现读 0。

### 三条如实报的事实

1. **形②未采**：子代理无法触发人工交互（`ask_user` 被运行时拒绝）。真人窗欠着，复采命令已备好（见读数台账文末），编排者摆机主在场那一发即可复跑。
2. **今日安静底噪 0.003067 ≠ 票面 2026-10-09 快照 0.368462**：DC 底噪在今日设备状态不存在（增益/驱动状态已变），票面核心前提「真麦安静读 0.37」今日不可复现；`0.376740` 今天只在合成 DC 帧上成立。照票面禁令，「设备增益」与「尺不减 DC」分开写：前者今日读数已翻面，后者是定义事实（形④钉死：DC+正弦帧读数精确含 DC 分量）。
3. **形③分开不了（具名）**：麦读 max 0.002530 < 安静窗①的 max 0.099908，与对照窗分不开 ⇒ 扬声器→麦通路今日无效。按票面判据具名「分开不了」，⛔ 不挑门槛凑绿。

### AC#2 判据

「说话」与「安静」之间既有的、不动用的判据能否分开——**形②缺席，判不了**。侧证两条：今日安静窗 mean 0.003067 < `SilenceLevelGate 0.06` ⇒ 球侧静默门**今日会正确关上**（10-09「永远开着」形状今日不复现）；尺不剥 DC 的定义事实独立于设备状态成立（形④）。

## AC#3 落点表（⛔ 未落码；全表见读数台账 §AC#3，此处摘要）

最小改法＝**两枚文件**：`internal/audio/level.go:96-107`（求和段 `:100-104` 改两遍遍历：先均值后 `(s-mean)²`）＋同步重写 `level.go:17-39` 尺的注释块（纯 -32768 帧从 1.0 变 0.0，端点语义搬家）。**必红的头号既有钉**＝`level_test.go:72-87` `TestLevelFullScaleSquareEndpoints`（逐字钉着 all -32768 → 恰好 MaxLevel，LevelOfSamples 与 FrameLevel 两处断言）；其余 9 枚测试输入全是零均值信号，读码预判不红，「改前后逐名红名册作差」按编排者 2026-10-09 14:2x 裁定欠在落地腿。`capturelevel_windows_test.go` DC 夹具是派生式期望不会红、但失去见证力；`SilenceLevelGate=0.06`（C21 冻结 token）与票 247 的 2× 门槛**一字不许顺手动**。

## AC#4 门禁读数（全部现跑）

- 主仓产码零改动：`git status --porcelain | grep -v '^??'` 只出现本腿写面（票 291 Progress log 一行）；新增未收件只在 `.scratch/wisp/probes/291/mic/` 与 `docs/evidence/s1/291-level-dc.md`。
- 两把 tasklist 尺 0；克隆树 `GOFLAGS= go build ./...` **rc0**；`sh scripts/d22scan.sh` **clean rc0**；`gofumpt -l` 对台件空。
- 写面清点：`.scratch/wisp/probes/291/mic/{00-four-shape-readings.md, resident_audio_291_probe_windows_test.go}` ＋ `docs/evidence/s1/291-level-dc.md` ＋ 票 291 Progress log 一行。`frontend/**`／`design/**` 零读零写；`PLAN.md`/specs/thresholds/golden/D43/C 表零碰；零 push；commit 显式 pathspec。
