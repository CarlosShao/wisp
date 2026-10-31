# 票 381 验收腿 `381-v3` —— dev `77c2cabfb`（381-r2）上终裁（非实现者）

> **尺**：被裁实现＝`77c2cabfb`（5 文件 +756/−54；`merge-base --is-ancestor` rc 0＝dev 祖先，开工 tip=`c73bde4b7`）。★**全部实验出自仓外导出树** `D:\tmp\wisp\20261031\381v3\export`（`git archive 77c2cabfb` 2264 文件，381-r2 五文件 `git hash-object` 与提交 blob 逐位同；另补 `.scratch` 子树供 d22scan 引用解析，分母不受影响）。★共享树纪律：102r2-r1（internal/panel）、59-r1（internal/audio）在飞，本腿零碰零归因；共享树只读用 `third_party`（DLL/模型）、`build/models/manifest.json(.minisig)`（未入库生成物，拷入导出树）、`tools/d22scan`（`go run`）。⛔ 未改任何仓库文件（除本件）；⛔ push。

## 0 一句话判语

**381-r2 三件结构全部真、≤400ms 真 MET**：真机 5 轮 unstick/reset 35–48ms（三跑 15 轮全部 ≤59ms，对照轮照播完 2.114s——381-v2 的 6,653ms 形状死了）；mutation 三发复刻全红＋末发后还原核 3/3 绿；veto 删除后 381-v2 的 28 句语料误伤 0/12。**六格：AC#3/#4 可翻 ⇒ 票 381 全勾可结；R26-1 可收口（编排者记账）**。实现件读数本腿逐项核真；两处口径差（"27 PASS"实为 26、staging 清单缺两项）§8 具名，不动摇实质。

## 1 ★★★ ≤400ms 实测复跑（staging 全套；一发定论）

- staging＝导出树 `go test -c -o bin\wisp381v3.test.exe ./cmd/wisp`＋三 DLL 同目录＋`WISP_SILERO_MODEL/SPIKE_MODELS/PUNC_MODEL/TTS_DIR/15E_TMP`＋**`WISP_KWS_DIR`**（41a2r1 的 `voice41a2r1KwsDir` env 优先，实现件 staging 清单没写）＋manifest 拷入＋CWD＝树根（`voice15eRepoFile` 从 CWD 上爬）。
- 终轮（rc 0，真 waveOut＋真 matcha＋生产 `TriggerInterrupt`，`Test381r2RealDevicePlayingPhaseAbortWithin400ms`）：

| 读数 | 本腿终轮 | 实现件对照 |
|---|---|---|
| 5 轮 unstick/reset（trigger→#27／→waveOutReset 落地） | **46/48/46/35/46 ms** | 41–46 ms（同带） |
| parked_to_stop（含 300ms 驻留） | 335–348 ms | 341–347 ms |
| 对照轮（无打断，照播完） | parked_to_done=2.114s | 2.095s |
| abort_returned | 0×4、2×1（两形都停） | 0–2 波动（具名非缺陷） |
| 终轮收尾 | spoken=1 aborted=5 gate 6/6 problems=0 | 同形 |

- staging 修正途中两轮全量（环境错、零产品判语）：46/49/46/37/37 与 36/47/37/59/45——**三跑 15 轮全部 ≤59ms**。前两轮失败＝pwsh 把未引号 `-test.v=true` 拆成裸 `-test`、CWD 放 `bin\` 找不到 manifest——均为本腿 staging 错，非产品。⇒ **≤400ms：MET**。

## 2 ★★★ mutation 三发复刻（仓外 `-overlay`；变异件＋SHA256 互异在 `381v3\mut\`）

| 发 | 变异（`git diff --no-index` 实测） | rc | 逐字红句（截引） |
|---|---|---|---|
| **M1** | `TriggerInterrupt` 摘 abort 接线（−3 行，token 保留＝381-r1 形） | **1** | `MUTANT 381r2-M1 SHAPE: TriggerInterrupt did not abort the speaker in flight (spoken=1, aborted=0)…the sentence plays out (381-v2's 6,653ms shape)` |
| **M2** | `speak()` 摘入口清 token（−4 行） | **1** | `MUTANT 381r2-M2 SHAPE: the announcement was killed by a stale interrupt token (spoken=0, chunks=0)…` |
| **M3** | `consume()` 复插 381-r1 veto 支路＋`strings` import（+9−1） | **1** | `MUTANT 381r2-M3 SHAPE: …"导航去最近的停车场" (case 0: finals=0, handed=0)…` |
| **还原核**（末发之后） | 原版二进制同三案重跑 | **0** | 3/3 PASS；树×三变异件 SHA256 differs=True×3；树文件＝提交 blob |

（本腿造变异件的两轮脚本抖折＝腿的工具错，与树无关；M1/M2/M3 形状与实现件 §4 逐字同源。）

## 3 ★★ R26-7 三字段 ＋ Abort 出手线程形状核

- **R26-7**：`buildResidentVoice`（resident_voice_windows.go:1735-1737）与测试构造器 `newResidentVoice`（:529-531）的 `provScope: voiceAnnounceScope`／`newSpeaker: openWaveOutSpeaker`／`interruptCh: make(chan struct{}, 1)` **逐字同值**；frames/plan/residency/queue 由函数体 ：1733、:1768-1773 补齐（26r3 修的真队列在）。**维持收口**。
- **Abort 从触发方线程出手＝是**：`TriggerInterrupt`（Esc 在 ball UI 线程）同步 `rv.speakAbort.Load().(func())`→`abort()`（:556-558）；`speak()` 在 `newSpeaker` 成功后 `Store(func(){ speaker.Abort() })`、defer 清（engines :958-966）；`Abort()` **先** `abortReq.Store(true)`（锁外，:424）**再**拿锁 `waveOutReset`（:432）。owner 停在 `Drain` ⇒ 唯一能停出声的线程就是触发方——**与 381-v2 点名的形状完全一致**；M1 红证明这半边承载真停。

## 4 ★★ winmm 分片等待核（D22 墙钟禁令合规证明）

`voice_out_windows.go` `releaseLocked`（:344-368）：预算＝`samples*1000/rate + voiceOutWaitGrace(2000ms)`（**音钟推导**）；按 `voiceOutWaitSlice=100ms` 分片 `WaitForSingleObject`；**片间查 `abortReq`**→`errSpeakerAborted`；`WHDR_DONE` 旗标核在 unprepare 之前（修掉既有"WOM_OPEN 残留信号提前 unprepare"潜在形）。**零墙钟读数**：等待界全部来自样本数→毫秒算术，事件等待不读时间差——D42#9/D22 合规。`Drain` 入口也查 `abortReq`（:403-405）。

## 5 ★★ veto 选删核 ＋ 28 句语料复跑

- `consume()` 内 `strings.Contains` 代码命中 **0**（只剩注释），`strings` import 已移除；final 一律 `finals.Add`→`finish()`→文本路。真 veto＝KWS 词级（onKwsHit，#13/#22）＋Esc（`vetoByEsc`→`approval.voiceInterrupt`→`TriggerInterrupt`，resident_audio_windows.go:395），都不骑转写路。
- **381-v2 的 28 句语料逐句过生产 `consume()`**（fresh leg，本腿探针 `zz_381v3_corpus_windows_test.go`，导出树）：良性 12/12、干净 10/10、真停 6/6 全部 **finals=1、handed=1、problems=0、零 EvInterrupt 边** ⇒ **误伤率 0/12**（381-r1 支路同语料 12/12）。真停句现走文本路＝删支路时具名的取舍；转写路结构性当不了 veto 通道（381-v2 §4 三事实）。

## 6 ★★ 正控（「无事发生」的那一半）

`Test381r2IdleInterruptIsInert` PASS（空闲 Esc：Sleeping 不动、timeline 空、problems 空、aborted/spoken=0、nil/无通道 leg 不 panic）；`SILENCE-26R3 gate_closed=0 gate_reopened=0 finals=0 spoken=0`（**静音门针零触发**）；`TestWireR1SilenceAndEmptyTextNeverSpeak`＋`Test15eSilenceNeverCallsTheHandOff` PASS。

## 7 ★★ 零回归复跑 ＋ d22scan

- **合跑终轮 rc 0：26 PASS／0 FAIL／0 SKIP**（Test15e 6＋WireR1 7＋26R3 3＋41a2r1 5＋381r2 **5**；`-test.run` 正则在此树恰匹配 26 枚测试函数）。**实现件写「27 PASS（381r2 6）」是计数口径差——381r2 文件只有 5 个 Test 函数，总数 26**；缺的不是测试，是 1 的算术。
- `LATENCY openDelay=262ms chunk=640ms ballLatency=902ms` **逐位同**（26r3/381-v2/实现件基线值）。`GATE-26R3 closed=1 reopened=1 events=2 mid_speech_transcripts=0`／`ECHO …capture refused frames_fed=6 transcripts_after=1`／`REOPEN finals=3 warm_wake=2` 全在；`spoken_samples` 逐跑浮动（89165–89466，含本腿 89436/89466）＝读数非断言。`internal/speech` 本腿未复跑（写面 cmd/wisp 专属，不在必跑清单；实现件 61/0/6 记录在案）。
- `d22scan` 导出树 **rc 0 clean**；分母 **bans#1-5 internal/=274 cmd/=45、ban#8 internal/=706 cmd/=162**（另 design 14、frontend 116、tools 25）——**274/45/706/162 全对**。⚠ 首扫 rc 1（23 条 phantom-citation）＝本腿首版导出排除了 `.scratch` 的探针假象，补入后清零，与代码无关。**正控两发＋可逆**：`✓`(U+2713) 串 ⇒ rc 1 `[emoji] zz_381v3_poscontrol.go:3`；裸 `go func(` ⇒ rc 1 `[bare-goroutine] :6`；删探针复跑 rc 0 分母回落。★⛔ 未拿 ban#8 扫 docs/。
- 两尺＋vet（导出树）：`gofmt -l`／`gofumpt -l` 全空 rc 0（381-v2 时代的 `internal/memory/undo.go` 真红已被他人修）；`go vet ./cmd/wisp` rc 0。

## 8 实现件读数可信度（本腿逐项核）

**核真**：三件结构的代码声称、真机五轮同带（41–46 vs 35–48ms）、M1/M2/M3 红句逐字、还原核、`262/640/902`、d22scan 分母、`consume()` 删支、`speakAbort` 接线形状、abort_returned 两形波动具名——全真。**口径差两处**：①"27 PASS（381r2 6）"实为 26（381r2 5），算术差、无缺测；②staging 清单缺 `WISP_KWS_DIR` 与 `build/models/manifest.json`（41a2r1 家族依赖；本腿首次复跑时 41a2r1 4 枚因此环境红，补齐后绿）。净判：**实现实质真、报告两处小口径失真**——比 381-r1 时代「自评过宽」轻一档，读数面本身可信。

## 9 ★★★ 六格逐条判 ＋ 收口判定

| 格 | 判 | 依据 |
|---|---|---|
| AC#1 构造器补全 | ✅（381-v2 已翻，本腿维持） | §3 三字段逐字同 |
| AC#2 nil-panic 消除 | ✅（381-v2 已翻，本腿维持） | `newSpeaker: openWaveOutSpeaker` 在；合跑零 panic |
| **AC#3 打断通道＋≤400ms** | **✅ 可翻** | Esc 真通道端到端（§5）＋出声期 Abort 真停：**真机 5 轮 35–48ms（三跑 15 轮 ≤59ms）≤400ms MET**（§1）；合成期 token 通路 M1/M2 验真（§2）；6,653ms 形状死（对照轮照播完 2.114s vs 打断轮 ~347ms 停） |
| **AC#4 端到端** | **✅ 可翻** | 播报中打断 ⇒ 真停（waveOutReset ≤59ms 落地×15 轮）＋门放开（live 6/6、probe 2/2）＋下一句能进（`ABORT-27 aborted=1 next_spoken=1 probe2_aborts=0`；残留 token 杀句消＝M2 红＋`STALE-TOKEN spoken=1 problems=0`）——逐跳读数齐 |
| AC#5 反向判据 | ✅（维持） | r2 新结构三发全红＋还原核（§2，本腿自跑，非实现件字面采信） |
| AC#6 纪律格 | ✅（维持） | 77c2cabfb 恰 5 枚 cmd/wisp＋证据笔；internal/**、D43、配置键、testdata 零碰；d22scan 其树 clean |

- **票 381 可否收口：可**——AC#1/2/5/6 已翻（381-v2）＋本腿翻 AC#3/4 ⇒ 六格全勾，**带残余结案**（残余全是非 AC 项，§10）。
- **R26-1 收口判定：具备**——「打断」半格的 ≤400ms 实测达标＋真通道端到端全绿，**请编排者在台账记收口**。R26-7 381-v2 已记，维持。
- **未定义即停四条**：全部未触发（Abort 真停实测、≤400ms 修到最优、veto 选删已做、internal/speech 零碰）。

## 10 残余逐条定性（全不阻塞收口）

1. **球单击打断未装配**——381-r1 自认、不在 381-r2 三件内；AC#3 判据是「至少一条真通道」（Esc 已真）。归口：后续独立票。
2. **`Close()`/`Abort()` 抢锁两形**——abort_returned 0 或 n 波动，两形都落 `waveOutReset`、出声都停（15 轮实测无差）；读数非缺陷，实现件已具名。
3. **speaker 层「播完」早于真实出声**——对照轮 parked_to_done=2.114s vs 全句出声 ~6.4s：Drain 仍早回（DONE 旗标分片等待已把 381-v2 时代的 9.84s/15.1s 错位收窄）；账仍归 26 侧，不在本票射程。
4. **转写路真停语义**——口说「停」现在是普通语句进 agent（具名取舍）；KWS veto 的在飞射程＝Listening/Confirming（#13/#22），与出声期打断两通道不混。
5. **本腿两条脚注**：我的 28 句探针初版漏 `Endpoint` facts 导致一次假红（腿 bug 已修，产品无涉）；staging 需 `WISP_KWS_DIR`＋manifest 两物——建议后续派单把这两项写进 staging 清单。

## 11 纪律核

写面＝本件＋`D:\tmp\wisp\20261031\381v3\**`（导出树/mut/bin/logs/cases 全在仓外；探针 2 枚、变异件 3 枚、poscontrol 探针已删并复跑可逆）。共享树零字节改动、零归因；rc 一律 `$LASTEXITCODE` 重定向取值；GOTMPDIR 每跑置空；三发变异全走仓外 `-overlay`；还原核时机＝末发之后；⛔ push/amend/reset/rebase/stash/checkout；⛔ 票面/台账改动（收口动作留编排者）。

---

**交回**：编排者——①票 381 六格全翻、按 §9 收口并在台账记 R26-1；②§8 两处口径差请在收 381-r2 的账面注记（不改票面）；③§10 残余 1/3 的归口（球单击票、26 侧 Drain 早回）请落账。
