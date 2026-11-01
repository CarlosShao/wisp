# 票 392 终裁腿 `392-v2` —— 非实现者裁决表（五格逐条判）

日期 2026-11-01 08:28–09:20 +08 · 腿 `392-v2`（终裁、非实现者）· 被裁件 **`a87ffe9b3`**（9 枚
+938/-49：cmd/wisp 3＋internal/audio/aec 5＋件 `392r1-session-integration.md`）·
派单存档 `.scratch/wisp/probes/orch/20261031/txt/392-v1-dispatch.txt`（v1 过夜静默死亡判死，本腿重派）。
**方法**：`git archive a87ffe9b3` → 仓外导出树 `D:\tmp\wisp\20261101\392v2\tree`（工作树零碰；
在飞 `382-w-r2`/`393-r2` 未触未归因）；staged 全链（A1326 三 DLL 同目录＋third_party/build 拷贝＋
cwd=包目录＋WISP_* env＋GOTMPDIR 只 unset）；音频错峰（每跑前进程表现核）；工具链 go1.27.1＋
mingw gcc 16.2.0＋staged `libwebrtc_audio_processing.a`（aec3r1 成品，零重建零源改）。
★**全部读数本腿现跑**；对照实验含**父提交 `2e134f638` 对照件**与本腿重构件、**实现件原二进制**
（`D:\tmp\wisp\20261031\392r1\bin\wisp392.test.exe`）三方互证。

## 0. 一句话

**五格全过判，可带具名残余收口**：交付为真——保温三形测死复现成立（本腿 ×3）、keepalive=false
逐字节 26p-r1 原形核成立、调参面四设置 K=4 复跑且「三旋钮无收窄且不亏解」双侧成立、反向三对
闭合、零回归结构面全持、写面恰 9 枚禁区零触碰。**本腿主新发现＝(a) 面是摆动脆弱面**（三方二进制
互证：pre-392 同红 ⇒ 非 392 回归，⇒ 「五轮全≤400ms」是时段敏感读数而非可复现保证）。

## 1. 五格逐条判

| # | 派单要求 | 判 | 本腿独立读数（全部现跑） |
|---|---|---|---|
| ① | ★★★保温三形测死复现＋keepalive=false 逐字节 26p-r1 原形核 | **可翻** | 见 §2/§3 |
| ② | ★★★调参面 K=4×四设置＋三旋钮结论核 | **可翻** | 见 §4 |
| ③ | ★★反向三对复刻 | **可翻** | 见 §5 |
| ④ | ★★零回归（26p 套件＋六族＋262/640/902＋d22scan） | **可翻（带 §6 具名残余）** | 见 §6 |
| ⑤ | ★写面恰 9 枚＋禁区零触碰 | **可翻** | 见 §7 |

## 2. ★★★保温 finding 复现（三形之一 ×3，即使全新 APM＋全新实测）

`Test392KeepalivePostureFindingsAndColdDefault` 本腿三跑全 PASS（47.65/47.02/53.84s）：

- **WARM-POSTURE round=2 triggers=0 episodes=0 三连**（第二暖窗不判）。
- 残余剖面骑**近端电平**（run1 res_dbfs=[-26.2 -21.8 -26.6 -22.4 -24.6 -29.3 -38.4] vs
  near=[-20.4 -23.6 -21.5 …]；同跑冷窗 echo-only res=[-50.3 -39.2 -43.6 …]）＝该窗回声抵消灭绝。
- ★**全新 APM 在场是码面事实**：keepalive 姿势下 release 无条件 `apm.Close()`＋arm 每窗重建
  APM，而 warm2 仍不判 ⇒ 毒源在跨窗携带的 silero/resampler/对齐态，不在 APM 实例。
- 暖种子不可传递／APM 跨窗携带失准两形＝码面钉死（monitor 头注＋`armJudge` 形状＋finding 字符串
  与实现件逐字节同）；本腿复现的三形取其最强形（姿势整体死、即使全新 APM）。

## 3. ★★★keepalive=false ＝ 逐字节 26p-r1 原形核（git diff 形状＋tripwire）

- 构造点字面 `keepalive: false`（`newVoiceBargeinMonitor`）；release 的冷分支＝**原
  `m.closeEngines()` 一行原位**；APM 构造 `NewApmTuned(零值)→NewApm` 原路（`aec.go`
  `Overrides()` 三字段全 >0 判定，零值＝引擎默认）；stop() 增 7 行 nil-safe 幂等 dispose
  （冷姿势下引擎已收，纯 no-op）。
- **26p 面零字节**：`git diff 2e134f638 a87ffe9b3 --stat` 对
  `resident_voice_bargein_26p_{windows,echo_windows}_test.go`、`gate.go`、`align.go`、
  `bargein.go` = **空**；红轮行号两形状同在 `26p_windows_test.go:354`（互证）。
- **FINDING MOVED tripwire 在**：committed 测试 :293 逐字节在盘（改好必须连默认一起翻）。
- 运行时默认断言三跑全过（`readings().Keepalive==false`）＋冷多窗 (a)(b) 6/6 窗过＋还原核回冷。

## 4. ★★★调参面复核（全四设置 K=4，超过「至少两设置」）

`Test392AecTuningDoubleTalkSwing` rc=0（101.58s），control-26 (-26.70)×matched echo 0.2：

| 设置 | 本腿 4 发 min..max | 摆动 | 实现件（过夜） |
|---|---|---|---|
| baseline 默认 | -28.73..-27.41 | **1.33dB** | 13.16dB |
| enr_threshold 0.45 | -44.98..-30.62 | **14.36dB（反宽）** | 15.03dB |
| trigger_threshold 6 | -40.04..-29.15 | 10.89dB | 11.69dB |
| max_dec_factor_lf 0.5 | -44.57..-33.79 | 10.78dB | 6.44dB |
| 还原核（关） | -29.21，tuned=false 回基线带 ✓ | —— | -28.44 ✓ |

**结论核**：「三旋钮无收窄且不亏解」双侧成立——enr 本腿再证反宽；trigger 同带；max_dec 最好
一发 -33.79 仍**全数跌破 -30 收紧档**（「有时过」→「从不过」，与昨同向）。★**读数精化**：摆动
宽度是**会话量**不是旋钮常数（baseline 昨 13.16 今 1.33；今晨最宽反而是 enr）——「真解不在这三
个旋钮里、指向 suppressor 状态机/流内一致性」被**强化**。验收腿引用宽度时必须带会话上下文。

## 5. ★★反向三对（全带还原核）

1. **冷默认多窗**：Part A 三跑 6/6 窗（-20 说话者 triggers=1、ready=9728 全新实测）＋(b) 收紧档
   吃 ✓、还原核回冷 ✓。
2. **调参关**：`AecTuned=false`＋level -29.21 回基线带 ✓。
3. **dispose**：暖姿势窗后 `EnginesHot=true` → `stop()` 后 `false`；
   `cascade=models-unload,panel-hide,free-os-memory,grants-expire` 零 incomplete ✓。

## 6. ★★零回归复核

- **(b)** PASS：episodes=2 全被收紧档吃、triggers=0、finals/handed Δ=0，残余 **-51.63 与实现件
  逐位同**；**M2** 逐字 PASS（本腿转写「我已经完成明天下午三点的会议改到四点请准时参加谢谢」，
  实现件「(不)已经完成…」——同一 distinctive 断言两形皆过＝ASR 抖动级）＋还原核 episodes=0 ✓；
  **等价 6 窗零漂移** PASS（lag=509 strength=0.70）。
- **(c)** 两发 **rc=0 双过**（实现件过夜一过一死）；摆动在 (c) 面独立复现：L1 **-28.94↔-33.60**、
  L2 **-36.94↔-43.41**；M1 回声自触发复现（-45 轴 -36.14/-32.90）＋还原核 triggers=0 ✓。
  「L2 断言摆动敏感＋keyed-episode 依赖＋未改 26p-r1 一字」核 ✓。
- **★★(a) 面新发现（本腿主产出）**：今晨同机三方实测 (a)
  `Test26PBargeinStopsWithin400ms`：本腿重构件 **0/3**、实现件原二进制 **0/2**、**父提交
  `2e134f638` 对照件 1/3**。红形一致：某窗 -20 说话者被引擎深压至 **-39.75..-45.16**（过链窗
  -20.12..-26.77），不触发→超时。⇒ 深压 firing 与 AC#2 文档化的双讲抑制摆动同源；**(a) 的红
  非 392 回归（pre-392 同红）**，**(a) 的绿非稳定性质（实现件 5/5＝昨夜机器态读数）**。
  「58/76/195/57/46ms 全≤400」判语在「其过夜读数为真」意义上成立，**不构成可复现保证**。
- **六族 staged 合跑 34/1/0 逐位复现**（唯一红 `Test384ResidentWatchdogEndToEnd`＝合跑时窗，
  A1329 既有归因；本腿首跑 33/2/0 系 staging 缺 `third_party\model-fetch`——补齐后逐位复现，
  非码面）。**`262/640/902` 逐位同**（`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms`）。
- `internal/audio` no-tag **RUN=98 FAIL=0 SKIP=11 逐位同**（两包 ok；aecspike 秒红系测试二进制
  DLL 搜索序 0xc0000135，PATH 前置三 DLL 目录后全绿——staging 事项非码面）；gofumpt `-l`=0；
  vet（无 tag／`-tags webrtcap`）双 rc=0；d22scan **clean rc=0**（分母逐位同：bans#1-5
  internal/=287、cmd/=52；ban#8 internal/=729、cmd/=181）＋**selftest 35/0/0** ✓。

## 7. ★写面恰 9 枚＋禁区零触碰

`git show --stat a87ffe9b3` 恰 9 枚：`resident_voice_bargein_392_windows_test.go`（新 457 行）、
`resident_voice_bargein_windows.go`（M）、`resident_voice_windows.go`（M +7 dispose 挂点）、
`internal/audio/aec/{aec.go,aec_stub.go,aec_webrtc_windows.go,wap_shim.h,wap_shim.cpp}`（M）、
`docs/evidence/s1/392r1-session-integration.md`（新）。禁区逐枚核过：`internal/{risk,session,speech}`
零触碰、D43/statemachine/thresholds/golden 零触碰、`internal/audio/bargein.go` 零触碰、**AEC3 源树
零触碰**（调参全在本仓 shim/Go：`wap_new_tuned` 走公开 `SetEchoControlFactory`＋
`EchoCanceller3Config::Validate` 先行）、deps.toml/testdata/frontend/design/scripts/.github 零触碰。

## 8. 裁决：可否收口＋残余定性

**可以带残余结案**（AC#1–AC#5 全可翻；判语见 §1）。残余三件，逐件定性：

1. **(a) 面时段脆弱（本腿新钉，双侧共有、摆动同源）**：处置＝26p (a) 断言的验收语义按「摆动面
   读数」裁量（加稳/重 key/放宽属验收腿＋编排者；26p-r1 一字未改）。保温毒源逐件隔离
   （resampler vs det）仍未做＝后续票输入。
2. **flush-reset API＝票 393 已开**：AEC3 无公开 flush/reset 面，真解在绑定层新 API＋suppressor
   状态机深处（本腿数据强化此方向）。
3. **工件溯源**：本腿复跑覆盖实现件 hardcode out 工件（`392-keepalive.json`/`392-aec-tuning.json`
   落 `D:\tmp\wisp\20261031\392r1\out\`，现为**本腿**读数）；实现件原工件已备份
   `D:\tmp\wisp\20261101\392v2\out-backup\`。26p 工件同（`26pr1\out\`，备份同处）。

## 9. 写面申报

仓内仅本件 `docs/evidence/s1/392-v2-verdict.md`；仓外 `D:\tmp\wisp\20261101\392v2\**`
（tree＝a87ffe9b3 导出、parent-tree＝2e134f638 对照导出、bin/bin-notag、logs 全量、out-backup），
只建不删。⛔ push（推送归编排者）。
