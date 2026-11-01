# 票 397 终裁腿 `397-v2` —— monitor 复活三处窗间清（非实现者裁决表）

日期 2026-11-01 · 裁决对象 `4c8b8ea60`（3 枚 +345/-77：产品码＋测试件＋证据件；真父 `d8cc6f785` 现核一致；397 实现归属全库唯一此提交）·
重派背景：397-v1 静默死亡（24h 零产物）判死，本腿从零复裁、⛔ 采信实现件字面 ·
本腿全部实验在导出树 `D:\tmp\wisp\20261101\397v2\`（`git archive 4c8b8ea60`；两枚码面 SHA256 **commit=导出树=工作树逐位同** `A07C1929…E59`/`D5EFFA2A…2CB`，pristine webrtcap 二进制本腿自建 rc=0）·
staging＝A1326 全套（三 DLL 贴 binary＋WISP_* 五 env＋cwd=包目录＋`unset GOTMPDIR`＋Git Bash runner；**每跑前进程表核**——在飞只读同行 394-v2 的 rig 反复入窗，本腿等待环错峰，无一发带争用强跑）·
工作树零触碰、⛔ push/amend/reset/rebase/stash/checkout。

## 0. 一句话

**五格全过、可收口（附三条新具名残余）。** ★三处窗间清读码逐块核（非采信）＋`keepalive=false` 默认逐字节不动现核；★AC#2 权威读数 rc=0 复跑（warm2＝episodes=1/triggers=1/engine_warm=true/engine_built=false/**ready=9728 逐位同**/rms=-25.71@-30，warm 剖面与冷 (a) 同带、(b) 残余 -47.2..-58.9 ⊂ -41..-59）；★AC#3 反向四发＋还原核：**M3 判死形两发红＋M4 瞬态形红定向复现**，M1 行为绿×3/M2 rc=0＝本材料无失败形（与实现件同向、跨面差异再证）；★零回归批量：26p 定向 **4/4**、调参族 rc=0、六族补齐 staging 后 **34/1/0 逐位归位**＋**`262/640/902` 逐位**、四数 **98/75/0/11 逐位同基线**、`./...` 红全在 package-main 编译域外（cmd/wisp 唯一红＝导出树无 `.git` 的结构预期红，具名）；★纪律格全过。新残余：M3 一发摆动、M1 仪器面 `-Inf` marshal 偶发、导出树实验须补 gitignored 面（§6）。

## 1. AC#1 三处窗间清一体落地 ✓（读码核）

1. **`release()` keepalive 分支**（新 :556-576）：**APM 持有＋`Reset()`**（拒绝→warn＋`Close()+nil`＝下窗冷重建兜底）＋**resampler `Flush()`**＋**det `Finalize(nil)`**（开口话态关闭、边界事件弃置）＋**1.2s 静音冷段**（`voiceBargeinCoolSamples=19200`、`CoolWindows=38` 整窗、`Accept(voiceBargeinSilence,nil)` 事件弃置）；else 分支＝`closeEngines()` 原样。旧形对照现核：旧分支＝APM **不持有、无条件 `Close()+nil`**、det/resampler 裸保温——申报「替换原 Close/nil 形」属实。
2. **`arm()` 重建点消解**（新 :389-421）：`!keepalive`→`closeEngines()`（冷姿势每窗照拆）；`cold:=det==nil` 只建缺的；**`apmTuning` 守卫**＝同调参持有复用（release 已 Reset＝just-created 态）、异调参 `Close()` 冷重建（one-parameter-one-firing 保持）；**judge measure-then-arm 零动**（drain.`armJudge` 区零 hunk，盲头/全新实测读数 ready=9728 与冷姿势同值）。
3. **`speakFar` rate 守卫**（新 :480-493）：held resampler `Rates()` from-rate ≠ 本窗 engineRate → 放下重建（engineRate≤0 拒建）。
4. **keepalive=false 默认逐字节不动**：旧 :294＝新 :316 `keepalive: false, // the measured-good default...` **逐字同**；`git diff -U0` hunk 全清单＝文件头注释/常量块/struct 注释＋apmTuning 字段/arm 构造块/speakFar 守卫/release 注释与分支——**closeEngines/dispose/readings/drain 零 hunk**；`resident_voice_windows.go` 零触碰。

## 2. AC#2 保温复活读数 ✓（断言改钉＋壳修复＋权威读数复跑）

- **断言改钉现核**：旧钉 `FINDING MOVED`（`warm2.triggers!=0||warm2.episodes!=0`）→ 新钉 **`WARM REVIVAL MOVED`**（`warm2.episodes!=1||warm2.triggers!=1`，:342-345）；翻动归属本落地腿＝A1355 裁定在案。
- **壳修复现核**：`voice392RoundOut` 全导出＋json tag＋转换器（+39 行 :139-176）接入两件工件、artifact 键 `finding`→`revival`；**392-r1 空壳本腿现读实锤**（原件备份 `D:\tmp\wisp\20261101\397r2\out-backup\392-keepalive.json`：`cold_rounds=[,]`、四 round 全 `{}`）——对照本腿 run1 件全字段落盘，**392-r1 空壳缺陷闭环**。
- **权威读数复跑**（pristine，A1326 全套）：**run1 rc=0（38.12s，七窗全序）**：冷 control＝keepalive=false built=true warm=false；冷 (a)×2＝triggers=1/ready=9728/built=true；冷 (b)＝triggers=0；暖 W1＝triggers=1/built=true；**暖 W2＝episodes=1、triggers=1、engine_warm=true、engine_built=false、ready=9728、rms=-25.71@-30**；restore＝(b) triggers=0/built=true。warm2 剖面与冷 (a) 轮同带（near -20.2..-23.6、res talker 带 vs 冷 -29.6..-22.6/-33.3..-24.3）；(b) 残余 **-47.2..-58.9** ⊂ 申报带 -41..-59。rms/lag 与实现件单点数字（-24.13/453）异属真 TTS 材料跨跑微差（392-v3 §7 在册口径：判据＝类别/相对）。

## 3. AC#3 反向判据 ✓（四发＋还原核；两形定向复现）

突变件＝`release()` keepalive 分支逐块摘除（**四件 diff 逐行核：每件恰为目标块**；四 overlay 构建二进制 SHA 互异且均 ≠ pristine；overlay JSON 全正斜杠键＝坑处方照用）。同一 pristine 测试、同一 staging：

| 发 | 摘除 | 本腿读数 | 判 |
|---|---|---|---|
| M3 | det Finalize＋冷段 | 首发绿（§6①）；**重跑两发 rc=1＝`episodes=0 triggers=0`（WARM REVIVAL MOVED 红句逐字）** | ✓ **判死形定向复现**（392-r1 finding 形回归） |
| M4 | 仅冷段 | **rc=1＝`episodes=1 triggers=0`** | ✓ **瞬态形定向复现**（hangover 假 start 吃掉 endpointer） |
| M1 | APM Reset（裸携） | **行为断言全过**（warm2 复活 1/1、restore 净、res 首 -25.4＝talker 带）＝复刻失败形（W2 triggers=0）未迁移；首 rc=1 落**仪器面**（§6②）；**补两发 rc=0 全 PASS** | ✗ 失败形未复现（与实现件同向） |
| M2 | resampler Flush | **rc=0 全 PASS**（warm2 1/1） | ✗ 无失败形（与实现件同向） |

**还原核：四发后 pristine 复读 rc=0（35.95s）**；工作树零触碰（全程导出树＋overlay）。**跨面差异具名（本腿三发再证）**：392-r2 复刻矩阵的 no-reset 全形（W2 triggers=0、残余骑近端 ≥14.8dB）在 staged 链三发（实现件 1＋本腿 1 行为绿＋M2）均未迁移——staged 几何上裸携 APM 仍抵消、判官仍判。

## 4. AC#4 零回归 ✓（导出树批量）

- **26p 套件定向 4/5**：keyed 四枚全 PASS（(a) ≤400ms/(b)＋M2 逐字/等价网格；p26-keyed PASS=4 FAIL=0）；
- **(c) 摆动族双树归因**：本腿窗＝397 树三发全绿、对照树 `git archive d8cc6f785`（码面 SHA `37465123…/d82a0a73…`＝397 前旧码现核）{绿,红}——**红名 `Test26PMarginRetestAndGateReverse` 与实现件对照树红同枚**（其窗：397 树 1 绿 2 红、对照树红）；两树跨窗摆动、红名集合零新增、397 冷路径码面零交集（§1）＝发动机摆动族 HEAD 既有形状归因闭合；
- **392 调参族** `Test392AecTuningDoubleTalkSwing`：本腿复跑 **rc=0 全绿**（对实现件「首发摆动红→复跑绿＋HEAD 对照绿」归因闭合再证一发）；
- **M1 补两发**（§3）：**rc=0 全 PASS ×2**（marshal 干净）——无失败形读数加固；
- **六族 staged 合跑**（no-tag 二进制＋三 DLL＋五 env＋cwd=包目录）：本腿首跑 29/6/0、二跑 33/2/0 秒失败假红＝**本腿 staging 缺口**（`git archive` 导出树缺 gitignored 的 `build/models/manifest.json` 对与 `third_party/model-fetch/` 实体——41a2r1 :224/:297 相对路径走查、A1326 卡「manifest 可达」在导出树实验同样适用，§6⑤）；**补齐后三跑＝34/1/0 归位**（唯一红 `Test384ResidentWatchdogEndToEnd` 合跑负载时窗既有形状，单跑 rc=0 两跑核）；**`262/640/902` 逐位**（`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms`，六族日志在案）；
- **四数**（无 tag、`./internal/audio/... -count=1 -v`、静窗）：**rc=0，RUN=98/PASS=75/FAIL=0/SKIP=11 顶层口径与 392-r1:108/392-r2/392-v3 逐位同**（票面 AC#4 的 94/83/0/11＝A1326 时代旧口径，台账 :29600 已勘误，非本腿缺陷）；SKIP 集＝TestLiveWasapiSmoke＋aecspike 硬件/工件门（含一枚密掩占位门名 `TestSpikeAec3LongWindow`）逐名同形；
- **`go test ./...` 全仓扫**（无 env 形，三跑）：**cmd/wisp 唯一红＝`TestReadGitOnThisRepositoryIsSelfConsistent`（0.00s：要求 cwd 上方有 `.git`——git 自省测试在导出树（刻意无 `.git`）的结构性预期红，工作树有；与 397 diff 零关系）**；其余红包＝panel（35r8 backpressure 界＋composer/approval 契约族）/risk/winsec/speech/tools/audio-aecspike（env 门形）/agent-approval（run1 一发、run2/3 过＝摆动）——**全在 package-main 编译域外**（`package main` 不可被 import，两枚改动文件的编译域只有 cmd/wisp 自身）；实现件申报红名族（panel/observe/risk/agent/winsec）与本腿观察集不同形＝env/负载形状差，且其 `./...` 日志未留档、红名逐名比对不可做，observe 本腿三跑全绿——**如实记，判断只取「无红在写面编译域内」这半句**。

## 5. AC#5 纪律格 ✓

写面恰 3 枚（`git show --name-only` 现核：产品码＋测试件＋证据件）；`resident_voice_windows.go`/`internal/{audio,speech,risk,session}`/AEC3 源树/389 门（`internal/audio/bargein.go`）/thresholds/golden/deps.toml 零触碰；票面/台账/specs 零触碰；本腿零 `.go` 写、工作树零触碰、⛔ push。

## 6. 残余与边界（如实记账）

1. **M3 首发摆动（新）**：M3 首发 rc=0 绿（恰逢 394-v2 rig 入窗重叠），重跑两发判死形稳定红——**staged 链判死形非逐发确定**（真 TTS 材料/实时时序摆动族，与 26p (c)/392 调参在册同族）；多发稳定性读数归默认翻转腿全链复读。
2. **M1 仪器面 `-Inf` marshal 脆弱性（新）**：`voice392Round.rms` 无 episode 时保持 `-Inf` 哨兵（:214），profile 桶纯静音同；**壳修复把 rounds 首次接上 JSON 序列化后**，(b) 轮 `episodes=0`（echo 完全门净＝比常形更干净的健康形）→ `write392JSON` 整件 `json: unsupported value: -Inf` → 测试尾部红（m1 发：行为断言全过、仅工件写盘失败）。非产品、非断言、非失败形；**默认翻转腿全链复读必撞面**——建议 JSON 安全哨兵（如 -999/null）小改测试件。
3. **跨面差异维持具名**：复刻矩阵 no-reset 形未迁移到 staged 链（三发）；**更广材料/几何上的 Reset 承重读数归默认翻转腿全链复读**（实现件申报同向；复刻矩阵 19.24dB 承重读数仍为最敏感证据，Reset 形维持不依赖本发）。
4. **冷段＝owner 侧纪律**：silero session 无 RNN reset 面，398（internal/speech 票）待开非阻塞；「衰减多少够」仍 1.2s 一点读数。
5. **导出树实验须补 gitignored 面（新，本腿踩坑入账）**：`git archive` 导出树缺 gitignored 的 `build/models/manifest.json`（12KB 对）与 `third_party/model-fetch/`（886MB 实体）⇒ 六族首跑 29/6/0、二跑 33/2/0 秒失败假红（41a2r1 :224/:297「cannot find … above the working directory」族）；补齐后三跑归位（§4 补记）。A1326 卡「staged 二进制＋DLL 全家＋WISP_* 五 env＋manifest 可达＋cwd=包目录」在**导出树实验**同样适用且须加 model-fetch 实体。同坑双向在案：**失败发不写 artifact ⇒ 工件目录残留上一过跑件**（本腿 cp 快照与实现件 `392-keepalive-m3.json` 同踩——其内容 1/1 与其 m3 日志 FAIL 矛盾，证据一律以日志 `--- PASS/FAIL` 行为准）。
6. 停止条款三件全未触发（judge 盲头原样；复活读数全愈；keepalive 默认零字节）。

## 7. 五格判＋可否收口

- **AC#1 三处窗间清一体落地 ✓**（§1 读码逐块＋默认逐字节）。
- **AC#2 保温复活读数 ✓**（§2 权威读数 rc=0、ready=9728 逐位、断言改钉＋壳修复闭环）。
- **AC#3 反向判据 ✓**（§3 两形定向复现＋M1/M2 无失败形如实记＋还原核）。
- **AC#4 零回归 ✓**（§4 批量；红名集合零新增）。
- **AC#5 纪律 ✓**（§5）。
- **可收口：是**（票 397 候选；AC 框勾选与票面 Status 归编排者收口动作，本腿⛔改）。

## 附：本腿自跑尺（全部可重放）

`D:\tmp\wisp\20261101\397v2\`：`tree\`＝`git archive 4c8b8ea60` 导出＋补 gitignored 面（`build/models/manifest.json` 对＋`third_party/model-fetch/`，§6⑤；两枚码面 SHA 对上）；`bin\`＝pristine＋M1-M4＋no-tag＋对照树二进制＋三 DLL；`mut\`＝四突变件＋四 overlay JSON（diff 在 `logs/ac-runs.log`）；`logs\ac-*.log`＝AC#2/AC#3 十一发全量；`out\392-keepalive-*.json`＝各发工件快照（**失败发不写盘、快照会残留上一过跑件——证据一律以 log `--- PASS/FAIL` 行与 ROUND-*-PROFILE 为准**，本腿与实现件同踩在案）；`head-control\`＝d8cc6f785 对照树；`logs\{quiet-batch,fixup,fixup2,four-numbers,sixfamilies3,dotdotdot3}.log`＝零回归批量全量。runner/处方＝397r2 run.sh 同形（A1326 全套＋每跑前进程表核＋与在飞 394-v2 rig 错峰等待环）。
