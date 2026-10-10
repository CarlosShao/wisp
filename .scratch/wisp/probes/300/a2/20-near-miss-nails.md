# 300-a2 · 20 — 擦边钉普查（⛔ 只找音频包；⛔ 替编排者裁）

**名册级尺（全部 blob 层、内容锚，⛔ 靠行号引产码）：**
- 尺①＝`git grep -n -e 'convertPacket' -e 'FloatToPCM16' -e 'MonoDownmix' HEAD -- internal cmd`（剔 `.scratch`）
- 尺②＝`git grep -n -e 'SineLevelTolerance' -e 'SilenceLevelGate' HEAD -- internal cmd`（剔 `.scratch`）
- 尺③＝`git grep -n -i -e 'thresholds' -e 'golden' -e 'observe\.' HEAD -- internal/audio`（剔 `.scratch`）
- 尺④＝`git grep -n -e 'wasapi_windows\.go:' -e 'parse_wave_format_300' -e 'hotplug_test\.go:' -e 'gate\.go:' -e 'device\.go:' HEAD -- cmd internal scripts tools docs .github`（剔 `.scratch`）→ 逐字存 `logs/anchor-rot.txt`
- 尺⑤＝`git ls-tree -r --name-only HEAD -- cmd/wisp | grep -e 255` ＋逐枚 `git show HEAD:<file> | grep -n '^func ' / 'internal/audio'` → 存 `logs/near-miss.txt`
- 尺⑥＝`git grep -n -e 'TestGoFiles' -e 'XTestGoFiles' -e 'NO-SCOPE' HEAD -- scripts cmd internal tools`（剔 `.scratch`）
- 尺⑦＝`git grep -n -e 'floating' HEAD -- cmd/wisp`（20 枚命中逐枚读）

**逐枚具名（它断的是什么／撞⛔ 撞）——共 7 枚（N1–N7）**

| 号 | 件（blob 名） | 它断的是什么 | 撞⛔ 撞 |
|---|---|---|---|
| **N1** | `internal/audio/parse_wave_format_300_windows_test.go` | 四形夹具喂 `parseWaveFormat`，断 `got.tag == tc.wantTag`（字面 3/1/7/0）＋header 三枚＋:152-153 的 `floating` 恒等式 | **最可能"被迫放宽"的一枚**：它是 `AC#1`/`AC#2` 凭据本体（硬约束②射程）。落点 C1 只要**⛔ 改那三行既有断言**就⛔ 撞；⚠ 但"共用同一枚文件"这件事本身会让 `300-v1`/`300-v2`/`300-v3` 引用过的 `:152-153` 逐字凭据与那对成对基线（`432/320`→`437/321`）**一起漂**——那是 docs 面锚腐烂（本仓老账＝行号锚的突变台件会打死锚位），⛔ 仪器红，⛔ 因此算"没撞" |
| **N2** | `internal/audio/resample_test.go` → `TestMonoDownmixAndFloatConvert` | 断 `MonoDownmix`/`FloatToPCM16` 的字面样本输出（**linux 也编得进**，该文件⛔ tagged） | **⛔ 撞**（⛔ 引用该常量）；但它是 C2/C4 走"输出字节面"判据时**唯一既有邻居**：新用例⛔ 得复用它的样本期望去反过来改它 |
| **N3** | `internal/audio/hotplug_test.go` → `TestLiveWasapiSmoke`（＋该文件其余 7 枚顶层用例） | 8 枚 windows-tagged 用例；`TestLiveWasapiSmoke` 靠 `WISP_LIVE_MIC != 1` 走 skip | **⛔ 撞**，**除非**新判据被塞进这枚文件：那会挪动 `:527`，而 `scripts/portable-tests.sh` 夹具 ledger 里逐字写着行号锚 `hotplug_test.go:527`（尺＝④ ⇒ scripts 侧唯一命中）⇒ **锚腐烂＝顺手提防的一枚** |
| **N4** | `internal/audio/level.go`（`const SineLevelTolerance = 5e-5`）＋ `internal/ball/liquid.go:30`（`SilenceLevelGate = 0.06`）＋ `internal/ball/liquid_test.go:60` ＋ `internal/ball/tokens_table_test.go:552`（表里以字符串键 `"SilenceLevelGate"` 钉着那枚常量）＋ `internal/audio/level_test.go`（10 枚顶层用例逐枚拿它做容差） | 电平尺定义（票面 `AC#2` 与 `AGENTS §1.1` 都把它列进⛔ 动面） | **会撞的那一族**：新判据若走"样本幅度"路线（C2/C4 的输出面），**最省事的写法就是复用 `SineLevelTolerance`**——复用⛔ 撞，**改值／加容差＝打红 `level_test.go` 那 10 枚 ＋ `liquid_test.go` ＋ `tokens_table_test.go`**（尺②⇒ 21 枚命中，逐枚是这三枚文件）。⇒ 落地派单里要把这句写死：**只许引用、不许动值** |
| **N5** | `cmd/wisp/config_receipt_255_test.go` → `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`、`TestTicket255RosterStillMatchesTheActualReadSites`（载具＝`cmd/wisp/config_readers_255.go:141` 那句逐字 `internal/audio/gate.go:57`；另 `cmd/wisp/resident_audio_247_windows_test.go:221` 引 `internal/audio/device.go:86/:89/:103`） | cmd/wisp 侧的"出处名册仍说实话"那族钉 | **⛔ 撞——而且我是量出来⛔ 推出来的**：尺④ 名册级结果＝`wasapi_windows.go:<行号>` 与 `parse_wave_format_300*` 被行号引的命中数 **0 枚**（非 `.scratch` 射程内），被 audio 侧按行号引的**只有 `gate.go`／`device.go` 两枚文件** ⇒ C1/C2/C3/C4/C5 五个落点⛔ 碰这两枚 ⇒ ⛔ 打红这族。**⚠ 两点限定（⛔ 说过头）**：① 尺④ 只证"⛔ 人按行号引那两枚文件"，⛔ 证"那族仪器⛔ 会去看 audio 的字节面"（它读的是 roster 自己的字符串与 read sites 扫描结果）；② `cmd/wisp` 台面此刻由 `303-r1` 独占 ⇒ 本腿⛔ 跑任何编译面去验它，那一发归**编排者或 `cmd/wisp` 的验收腿** |
| **N6** | `scripts/portable-tests.sh`（GUARD A/C/D ＋ 夹具 ledger ＋ `core_pin`/`win_pin`/`cli_pin`/`winsec_pin`） | 包级名册：每枚 scope 条目⛔ 是空分母（GUARD A）、pin 与清单同笔（GUARD C）、⛔ 有人认领（GUARD D） | **⛔ 撞（就"枚数"而言）**：尺 K4 名册＝`internal/audio` 在该脚本只有 5 处命中（`core_pin` :168、`win_pin` :195、`core)` 清单 :242、`windows)` 清单 :253、ledger :594），**逐枚都是包级或单枚用例名、全脚本⛔ 一枚"该包应有几枚用例"的计数钉** ⇒ 新增用例⛔ 改 pin。**⚠ 一枚真会撞的形**：若新判据在某些机器上 `t.Skip`（C3 的 SDK 缺失形）⇒ 未登记 SKIP 是红档，登记＝改 scripts ⇒ **越出票面 `AC#4` 名册（⛔ 只许 `internal/audio/**`＋`probes/300/**`）** ⇒ 那一格是"被迫扩射程"⛔"被迫放宽断言"，具名交编排者 |
| **N7** | `cmd/wisp/resident_audio_247_windows_test.go`、`cmd/wisp/resident_mute_290_windows_test.go`、`cmd/wisp/resident_tray_mute_293_windows_test.go`、`internal/audio/capturelevel_windows_test.go`（`TestAC247RealCaptureLoopEmitsLevels`） | 采集链路下游（电平/静音/常驻）——**常量错了真正会疼的那一族** | **⛔ 撞**（⛔ 一枚能引用未导出常量），⚠ 更要紧的是反方向：**⛔ 指望它们兜住 `AC#6`**——尺① ⇒ `convertPacket` 的测试侧调用者 **0 枚**、尺⑦ ⇒ `cmd/wisp` 里 20 枚 `floating` 命中**逐枚都是"浮动球窗口"的散文，与音频浮点⛔ 干**（我把 20 枚逐枚读了）⇒ 全仓**⛔ 任何一枚既有用例观察过"生产怎么用这枚常量"**。这正是 `AC#6` 存在的原因，也⛔ 是"落点只能长在 `audio` 包里"的理由 |

**⛔ 动任何既有断言（本腿一句纪律复述，出处＝`AGENTS §1.1` ＋ 票面 `AC#6` 硬约束② ＋ 票面"⛔ 为了让它绿而同时改产码"）**：上面 N1–N7 里，唯一"会被顺手打红"的形是 **N4（改容差/token 值）** 与 **N6（新增 skip 而⛔ 登记）**；两者都靠**⛔ 那么写**来避免，⛔ 靠放宽谁。

**本腿自报的一处形状**：本目录 `.md/.txt` 落笔时 git 报 `LF will be replaced by CRLF`（core.autocrlf=true、`.gitattributes` 未覆盖 `.scratch/**/*.md`），⇒ 我的件的 **blob 与工作树 EOL ⛔ 同值**，⛔ 人拿我这堆件做 `cmp` 逐字节凭据；本腿⛔ 动 `.gitattributes`（⛔ 属我射程）。

rc=0
