# 票 392 落地腿 `392-r1` —— 28-PathC 会话集成件（写腿）

日期 2026-10-31 · 腿 `392-r1`（写、cmd/wisp＋internal/audio＋本件）·
起手 HEAD `2e134f638`（dev，开工三核过：①锚点在 ②索引净（脏区全 `.scratch/**`）③写面无他人脏）
· 进程表核过无 rig 在飞（音频错峰）· 派单存档
`.scratch/wisp/probes/orch/20261031/txt/392-r1-dispatch.txt` 在案。
**26p-v1 判语不在盘**（glob `docs/evidence/s1/*26p*.md` 仅 26pr1-impl）⇒ 前提＝26p-r1
读数＋停手条款，按派单字面执行。

## 0. 一句话

三个设计输入全部消费并出数：**AC#2 调参面落地**（绑定层
`EchoControlFactory` 注入 `EchoCanceller3Config`，AEC3 源树零改）——双讲损伤摆动
K 发实测 **baseline 13.16dB**（-41.28..-28.12，比 26p-r1 两发的 0-8dB 更宽＝设计输入
精化），三参数各一发过链：`trigger_threshold 12→6` 微收 11.69dB、`enr_threshold
0.25→0.45` 反宽 15.03dB、`nearend_max_dec_lf 0.25→0.5` **摆动砍半 6.44dB 但最好一
发 -28.1→-35.0（全跌破收紧档）＝引擎级裕量交换**，还原核回基线带；**AC#1 保温按
铁律裁定＝试过、测死、默认关**——暖姿势实测不可用（第二个暖窗回声抵消失效，复现
3 次，即使全新 APM＋全新实测），生产默认回到逐字节 26p-r1 原形，发现钉成读数＋工
件；冷姿势多窗 (a)/(b) 全健康；dispose 级联收口（EnginesHot 翻转）。零回归：26p 套
件 (a) 五轮 46-195ms ≤400ms、(b)＋M2 逐字、等价 6 窗零漂移、(c) 两跑摆动带内；
六族 staged **34/1/0 逐位复现基线**（唯一红＝384 合跑时窗既有形状）；**262/640/902
逐位同值**。

## 1. AC#1 跨公告保温＝28 的集成深度：实测裁定「默认关」

保温姿势（28-r1 warm-pool 形状投到 monitor）按派单要求做了**真实现＋真测量**，三个
候选形状依次实测：

1. **暖种子（上一窗实测 lag 作下一窗判决时延）**：seed 臂结构上成立（judge 在远端
   头出现即武装，ready 512 样本 vs 冷测 9728），但 fresh 验证与种子发散（-59 vs
   +474，远端 tap 与采集起点的交错每公告抖动），且种子 APM 零抵消。verify-then-seed
   就等价冷测 ⇒ **种子不可传递，实测证伪**。
2. **APM 跨窗保温**：回声路状态是渲染/采配对的**流内资源**，窗末截断留下无配对渲染
   帧——测量：跨窗携带的 APM 对下一窗彻底失准（残余≈近端电平）。改为每窗重建后…
3. **暖姿势整体（silero＋resampler 保温、APM 每窗新）**：第二个暖窗仍不判——
   episodes=0、triggers=0、残余剖面骑在近端电平（echo-only 段 -22~-25dBFS）——
   **复现 3 次**，即使该窗全新 APM＋全新实测。冷姿势同一序列多窗全健康。毒源逐件
   隔离（resampler vs det）未做＝具名残余；AEC3 无公开 flush/reset 面，真解属绑定层
   新 API（后续票）。

**裁定（铁律①④：最有价值，不因动作大而降级，也不为复杂而复杂）**：生产默认
`keepalive=false`＝逐字节 26p-r1 原形；姿势保留（`setKeepalive`）供测试驱动；发现
钉成 `Test392KeepalivePostureFindingsAndColdDefault` 的 finding 断言（第二暖窗不判
⇒ 红名「FINDING MOVED」，未来修好必须连默认一起翻）＋剖面工件
`392-keepalive.json`。**读数**（staged 全链）：

- 冷默认多窗：round 1/2 各 `triggers=1`（-20 说话者过收紧档）、ready=9728（全新实
  测）、lag 536/490；(b) echo-only `triggers=0` episodes 全被收紧档吃；还原核同形。
- 暖姿势：round 1 判到（triggers=1）；round 2＝finding（triggers=0 episodes=0）。
- dispose 级联（28-r1 对齐）：暖姿势窗后 `EnginesHot=true`（模型引擎在窗间持有）→
  `stop()` 的 C11 release 步收走 → `EnginesHot=false`（`Test392DisposeReleasesTheEngines`）。

## 2. AC#2 AEC 调参读数（绑定层面，⛔ AEC3 源树改未触）

**参数面**：`internal/audio/aec` 新 `Tuning`（三字段，零值＝引擎默认＝`NewApm` 原路）
＋`NewApmTuned`（走公开 `AudioProcessingBuilder.SetEchoControlFactory` 注入
`EchoCanceller3Config`——shim 加 `wap_new_tuned`，`EchoCanceller3Config::Validate` 先
行；静态库 RTC_EXPORT 展开为空，链接零改）。**停手条款③未触发**（面在绑定层暴露）。

**K 发摆动读数**（L1 形状：control-26 -26.70 over matched echo 0.2，talker@1.2s，
每设置 4 发，全部冷姿势＝26p-r1 可比；工件 `392-aec-tuning.json`）：

| 设置 | post-AEC episode RMS（4 发） | 摆动 |
|---|---|---|
| baseline（引擎默认） | -41.28 .. -28.12 | **13.16dB** |
| dominant_nearend enr_threshold 0.25→0.45 | -43.26 .. -28.23 | 15.03dB（不收反宽） |
| dominant_nearend trigger_threshold 12→6 | -40.71 .. -29.01 | 11.69dB（微收，带内） |
| suppressor.nearend_tuning max_dec_factor_lf 0.25→0.5 | -41.40 .. -34.96 | **6.44dB（砍半）** |
| 还原核（调参面关） | -28.44 | 回 baseline 带 ✓ |

**如实结论（设计输入精化）**：摆动主导项确证且比两发读数更宽（13dB）；三个可进绑
定层的参数里，**没有「拉平且不亏」的解**——max_dec 砍半摆动的代价是最好一发被拉下
6.8dB（透明的 nearend 态增益下潜变慢＝引擎级裕量交换，全部 4 发跌破 -30 收紧档，
「有时过」变「从不过」）；enr/trigger 两态切换阈值在 ±1.5dB 内＝抖动级。**真解不在
这三个旋钮里**，方向指向 suppressor 状态机深处/流内一致性与 AEC3 调参面扩展（后续
票的输入）。

## 3. AC#3 反向判据（三读数＋还原核，全在 staged 全链）

1. **冷默认＝生产姿势**：六窗序列 (a)(b) 形状全持（§1）；反向＝暖姿势第二窗不判
   （finding 断言钉住）；还原核＝切回冷，(b) held ✓。
2. **调参面关**：`AecTuned=false` 读数＋level 回 baseline 带（-28.44）✓。
3. **dispose**：`EnginesHot true→false`（stop() 的 release 步）✓；scope 级联日志
   `cascade=models-unload,panel-hide,free-os-memory,grants-expire` 零 incomplete。

## 4. AC#4 零回归（Path-C 三真数不退化＋六族）

- **26p 套件**（392 二进制，新 monitor，冷默认＝原姿势）：(a)
  `Test26PBargeinStopsWithin400ms` PASS——五轮 **58/76/195/57/46ms 全 ≤400ms**；
  (b) `Test26PEchoNeverTriggersAndNeverTranscribes` PASS（episodes=1 被收紧档吃、
  triggers=0、finals/handed Δ=0）＋M2 `Test26PWithoutAECTheEchoTranscribesVerbatim`
  PASS（逐字「(不)已经完成明天下午三点的会议改到四点请准时参加谢谢」＋还原核）；
  等价网格 `Test26PFeedOnlineMatchesWindowReplay` PASS（6/6 逐位同、零漂移）；
  (c) `Test26PMarginRetestAndGateReverse` **两跑一过一死**：run1 L1=-28.69 过／L2
  key 到 -45.79 死；run2 全过（L1=-41.61 被吃＝如实代价、L2=-41.46 标称可达）。同
  一 control-26 post-AEC **-28.69↔-41.61（13dB）**＝§2 摆动在 (c) 面独立复现；(c)
  的 L2 断言对摆动与 keyed-episode 顺序敏感＝26p 测试形状在实测摆动下的脆弱点
  （⛔ 未改 26p-r1 一字，读数在案供验收腿裁量）。
- **六族 staged 合跑**（no-tag 二进制＋三 DLL 同目录＋WISP_SILERO_MODEL/
  WISP_SPIKE_MODELS/WISP_15E_TMP＋cwd=包目录，`-test.run
  "Test43|Test15e|Test41a2r1|Test381r2|Test390r1|Test26R3|Test384"`）：
  **34 PASS / 1 FAIL / 0 SKIP**——对 26p-r1 基线逐位复现；唯一红
  `Test384ResidentWatchdogEndToEnd`＝合跑负载时窗（384-r1/390-r1 A1329 既有归因，
  单跑 PASS 的 HEAD 既有形状）。**`262/640/902` 逐位同值**
  （`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms`）。
- `internal/audio/...` 无 tag `-count=1` 三连 `ok`（FAIL=0；其一早跑曾逢六族收尾负
  载红一次，复跑三绿，归因同 384 形状）；`-v` 口径 RUN=98 FAIL=0 SKIP=11。
- 两尺：gofumpt `-l`（cmd/wisp＋internal/audio）＝0；vet（无 tag／`-tags webrtcap`
  两形状）rc=0；build 两形状 rc=0。
- d22scan（Git Bash `scripts/d22scan.sh`）：实扫 **clean rc=0**（bans#1-5 internal/
  =287、cmd/=52；ban#8 internal/=729、cmd/=**181**＝+1 本腿测试件）；self-test
  **PASS=35 FAIL=0 SKIP=0（RUN=77）**。

## 5. 写面申报

仓内（commit pathspec 逐枚）：
`cmd/wisp/resident_voice_bargein_windows.go`（M：keepalive 姿势＋armJudge/退点＋
调参接面＋读数扩）· `cmd/wisp/resident_voice_windows.go`（M：stop() 的 dispose 收口
3 行）· `internal/audio/aec/aec.go`（M：Tuning＋Overrides）·
`internal/audio/aec/aec_webrtc_windows.go`（M：NewApmTuned）·
`internal/audio/aec/aec_stub.go`（M：stub）· `internal/audio/aec/wap_shim.h/.cpp`
（M：wap_aec3_tuning＋wap_new_tuned 工厂注入）·
`cmd/wisp/resident_voice_bargein_392_windows_test.go`（新：三测试）· 本件。
仓外 `D:\tmp\wisp\20261031\392r1\**`（bin/bin-notag/cases/logs/out，只建不删）。

⛔ 越界核：`internal/{risk,session,speech}` 零触碰（SetConversation 公开面未动）；
D43 表/statemachine/thresholds/golden 零触碰；**AEC3 源树零触碰**（调参全在绑定层
shim/Go）；`internal/audio/bargein.go` 零触碰（389 门本体）；deps.toml/下载/
testdata 仓根/frontend/design/scripts/.github/票面/台账/specs/PLAN 零触碰；⛔ push
（推送归编排者）。

## 6. 残余与边界（如实记账）

1. **暖姿势毒源逐件隔离未做**（resampler vs det 哪个带毒）：姿势已默认关＋finding
   钉住；隔离属后续票（连同 AEC3 绑定层 flush/reset API）。
2. **AC#2 拉不平如实报**：三参数无一「收窄且不亏」；方向具名（§2）。
3. **(c) L2 断言的摆动敏感性**：26p-r1 测试形状在实测 13dB 摆动下 2 跑 1 死；本腿
   未改其一字；处置（放宽/加稳/重 key）属验收腿与编排者。
4. **`-20dBFS` 说话者与 control-26 裕量**：26p-r1 残余②原样成立（收紧档对 -26.7
   的过链不保证，摆动是主导项）。
5. 26p-v1 判语未在盘；若后到且与本件前提冲突 ⇒ 停手重对齐条款仍在。
