# 297-a1 · 40 现有仪器名册 ＋ 这批改动会不会顶红 CI

锚＝`e96da9f4`。⚠ 派单点名的一条我照做：**⛔ 不写 `ci.yml:605`**（那一枚锚在本仓被证过错）。下面凡是 CI 位置的引用，一律**逐字那一行＋它跑的脚本名**做内容锚，行号只作旁证，并标明尺与射程。

---

## 1. build tag 名册（尺 R15，现量）

尺＝`for f in $(git ls-tree --name-only -r e96da9f4 internal/audio cmd/wisp | grep '_test\.go$'); do t=$(git show e96da9f4:$f | grep -m1 '^//go:build' || echo NO-TAG); ...`
射程＝`internal/audio/` ＋ `cmd/wisp/` 的**全部 `_test.go`**（对象层，非工作树），**共 86 枚件**（尺＝R15 输出逐行数点名，现量：`cmd/wisp` **79** 枚＋`internal/audio` **7** 枚；⚠ 这一枚数我第一次落笔写成"87＝80＋7"，是把 R15 输出的**表头行**也数进去了，已就地改正——落笔时手数，未跑 `wc -l`，所以这枚 86/79 是**算的**，要现量请跑 `git ls-tree --name-only -r e96da9f4 internal/audio cmd/wisp | grep -c '_test\.go$'`）。

### 1.1 与**电平／音频链**直接相关的那几枚（本票射程）

| 件 | 测试／内容锚（函数名，逐字） | build tag（R15 现量） | 现量级 |
|---|---|---|---|
| `internal/audio/level_test.go` | 票面 `:14` 记的"全零帧读 0"那节（票面写作 `:46` 起的 `// AC#1 first of the three: …`＋`:52-55` 断言） | **无 tag** | 票面引文级（⛔ 我未逐字复对该行号） |
| `internal/audio/resample_test.go` | `TestMonoDownmixAndFloatConvert`（票面 `:16` 记的名），内含 `MonoDownmix(stereo, 2)` 与 `FloatToPCM16(f)` | **无 tag** | 本腿现量＝尺 R13 `git grep -c "MonoDownmix\|FloatToPCM16" e96da9f4 -- '*.go'` 给出该件 **5 枚命中**；⚠ 函数名本身是票面转述，我未逐字复对 |
| `internal/audio/capturelevel_windows_test.go` | `TestAC247RealCaptureLoopEmitsLevels`（我**全文读过**，逐字块见 20 件 §ⓑ） | `//go:build windows` | **本腿现量** |
| `internal/audio/hotplug_test.go` | `TestAC247…` 无关；本票关心的是 `fakeStream` 的 `func (s *fakeStream) Drain() ([]int16, error)`（R12 现量，`internal/audio/hotplug_test.go:108`）＋ `TestLiveWasapiSmoke`＋`TestPinnedThreadStable10s`（R16 现量上下文） | `//go:build windows` | **本腿现量**。⚠ 陷阱具名：**文件名没有 `_windows` 后缀却是 windows-tagged** ⇒ 任何靠 `*_windows_test.go` 猜射程的尺都会漏它 |
| `internal/audio/wavinjector_test.go` | 未逐字读；R13 现量它含 **1 枚** `MonoDownmix|FloatToPCM16` 命中 | **无 tag** | 现量（计数级） |
| `cmd/wisp/resident_audio_247_windows_test.go` | 未逐字读 | `//go:build windows` | **本腿现量**（R15） |
| `cmd/wisp/resident_audio_247_live_windows_test.go` | `TestAC247LiveMicrophoneLevelsReachTheBallSeam`＋`playAlarmThroughSpeakers`＋`requireNoCompetingAudioProcesses`＋`bootAudioRuntime`（R16 现量上下文里逐字可见后两枚调用） | `//go:build windows` | **本腿现量** |
| `cmd/wisp/resident_mute_290_windows_test.go` | 未读内容 | `//go:build windows` | 现量（R15） |
| `cmd/wisp/resident_tray_mute_293_windows_test.go` | ★ **`293-r1` 那枚在飞写腿刚落的件**（`e96da9f4` 提交标题自证："AC#4 门禁件…"） | `//go:build windows` | 现量（R15）。⇒ **`cmd/wisp` 的音频族名册在我读数期间仍在呼吸**，落地腿必自己复跑 |
| `cmd/wisp/resident_ball_228_windows_test.go`／`resident_ball_live_228_windows_test.go` | 球侧吃电平的那两枚（票面"为什么要紧"节说球的液态吃的就是这枚数） | 前者 `//go:build windows`、后者 **`//go:build windows && winlive`** | 现量（R15） |

### 1.2 一枚与直觉相反的名字（具名，防落地腿踩）

`cmd/wisp/leg_sink_nail_131_windows_test.go` ⇒ R15 现量 tag ＝ **`NO-TAG`**（文件名带 `_windows` 却没有 build 约束）。同目录 `resident_sink_nail_127_windows_test.go` 才是 `//go:build windows`。⇒ **"电平/采集输出钉子"这一族里有一枚今天会在 ubuntu 上跑**。若 `AC#2` 动到电平语义，这枚是最先可能被顶红的一枚（⚠ **算的**：我没读它内容、没跑它，只做了"tag 档位＋文件名族"两级现量 ⇒ 落地腿引之前先读它断言）。

---

## 2. CI 真链路（逐枚读脚本，⛔ 不靠 API 名字推）

尺 R17＝`grep -nE "go test|scripts/|windows|tags|winlive" .github/workflows/ci.yml`（射程＝该 workflow 全文）；尺 R18＝`grep -nE "scope|--tags|go test|SKIP|skip|-run" scripts/portable-tests.sh`（767 行全文）；尺 R19＝`grep -rn "wisp-cli-tests" .github/workflows/ scripts/`；尺 R24＝`awk 'NR>=230 && NR<=275' scripts/portable-tests.sh`（scope 清单原文）；尺 R25＝`grep -n "audio" scripts/portable-tests.sh scripts/wisp-cli-tests.sh`；尺 R26/R28＝SKIP ledger 与 247 的 grep。全部现量。

执行面逐档（**只有这五档真跑测试**）：

| 档 | 逐字锚（那一行） | 实际跑什么 | 与本票的关系 |
|---|---|---|---|
| core（ubuntu job） | `        run: bash scripts/portable-tests.sh --scope=core`（R17 现量行号 470） | scope 数组逐字含 `./internal/buildinfo/... ./internal/audio/... ./internal/proc/...`（R24 现量） | ⇒ **`internal/audio` 的非 windows 件（`level_test.go`／`resample_test.go`／`wavinjector_test.go`／`gate_test.go`／`captureopt_test.go`）在 CI 真跑**；⚠ **windows-tagged 那三枚（capturelevel／hotplug）在 ubuntu 上不编译＝不在执行面** |
| winsec | `        run: bash scripts/winsec-tests.sh`（568） | `./internal/winsec/...` | 无关 |
| **cli（windows job）** | `        run: bash scripts/wisp-cli-tests.sh`（**604**；⛔ 不是 605） | `scope=(./cmd/wisp/)`（R24 逐字），且 `wisp-cli-tests.sh` 里有 GUARD 逐字 `if [ "$(go env GOOS)" != windows ]; then` → `exit 2` | ⇒ **`cmd/wisp` 的 `//go:build windows` 电平族（含 247 live）今天在这档真跑** |
| windows | `        run: bash scripts/portable-tests.sh --scope=windows`（780） | scope 逐字＝`./internal/proc/ ./internal/secret/ ./internal/config/ ./internal/risk/ ./internal/ball/ ./internal/perm/ ./internal/plugin/ ./cmd/llmrecord/ ./internal/session/ ./internal/projctx/` | ★ **`./internal/audio/` 不在 windows scope**（R24＋R25 双尺现量：全文里 `audio` 只出现在 `:168` 的 ledger 与 `:241` 的 core scope）⇒ **`internal/audio` 的 windows-tagged 用例在整条 CI 上没有任何一档真跑**（core 档不编译、windows 档不点名） |
| winlive 编译面 | `        run: go vet -tags winlive ./cmd/wisp/ ./internal/ball/`（655） | 只 vet；ci.yml 自己的注释逐字（R17 现量）：`and are NOT wired here; do not read this step as "the winlive cases now` | ⇒ `winlive` 那 12 枚件（尺 R32＝`grep -rlE '//go:build.*winlive' --include='*.go' . \| wc -l` 现量 **12**，与 ci.yml 注释自报的 `7 under cmd/wisp, 5 under internal/ball` 对得上）**只编译、不执行** |
| census | `        run: bash scripts/portable-tests.sh --scope=census`（744） | 审计，不跑用例（`go test -list`） | 见 §4 |

---

## 3. "跳过≠通过"这台机器长什么样（这是本票 `AC#0` 落地时最容易被顶红的一枚）

`scripts/portable-tests.sh` 头部规矩，逐字（R18 现量，整块三行未省）：

```
#   1. a non-zero `go test` exit code propagates unchanged;
#   2. any top-level `--- SKIP` in what actually ran is fatal;
```

执行点逐字（R18 现量）：`skipped=$(count '^--- SKIP')` → `if [ "$skipped" -ne 0 ]; then` 并打印每条 SKIP 的 `file:line` 与理由。
`-skip` 模式**由一枚 ledger 生成**（逐字：`skip_pattern="^($alt)\$"`，空 ledger 时退化为哨兵 `'^portable_tests_ledger_is_empty_on_this_platform$'`），ledger 每行的形状＝`名字|包|平台|类别|理由`。

### 3.1 音频族在 ledger 里的现状（尺 R25/R26/R28，全部现量）

- `scripts/portable-tests.sh:593` 逐字（整行）：
```
    "TestLiveWasapiSmoke|./internal/audio/|windows|fixture|live WASAPI capture needs WISP_LIVE_MIC=1 and a physical microphone (hotplug_test.go:527, //go:build windows); a hosted runner has no audio endpoint, so the case has no subject there. Real-hardware smoke is ticket 16's"
```
- **`grep -n "247" scripts/portable-tests.sh` ＝ 0 命中**（尺 R28 现量）⇒ `TestAC247LiveMicrophoneLevelsReachTheBallSeam` **不在 ledger 里**；`grep -nE "TestAC247|WISP_LIVE" scripts/wisp-cli-tests.sh` 也 0 命中（同尺）。

### 3.2 ★ 由此浮出的一枚**先于本票存在**的隐患（本件最重要的读数，⚠ 一半是算的）

四条现量放一起：① cli scope＝`./cmd/wisp/` 且真跑（R24）；② 那枚 live 用例是 `//go:build windows`、windows runner 上会编进去（R15）；③ 它进门第一句就 `t.Skip`（R16 逐字）；④ 规矩 2 "any top-level `--- SKIP` is fatal"＋它不在 ledger（R18＋R28）。
⇒ **可测的预言（算的）**：windows CLI 那一发今天要么**因为一枚 unaccounted SKIP 而红**，要么存在我没找到的第五枚豁免机制。
⇒ 反证线索（我也现量到了）：`scripts/wisp-cli-tests.sh:20` 与 `ci.yml:588` 都写着同一枚读数 `PASS=33 FAIL=0 SKIP=0` ⇒ 两者**不可能同时为真**，除非那枚 33 读数是**在 247 live 用例进仓之前**测的（注释里没写日期，尺 R27 我未复跑头部全文）。
⇒ **具名：本腿判不了，且判它需要 `go test -list` 或一次 CI 日志读数 ⇒ 前者被派单禁跑、后者不在我预算内 ⇒ 这一条写"未证成"，⛔ 不许被任何一方读成"CI 今天是红的"或"CI 今天是绿的"。**
⇒ 但对**落地腿有直接可执行的后果**：若 `AC#0` 新增第二枚 `WISP_LIVE_MIC` 门控用例并落进 `cmd/wisp/`，它会撞上同一台机器；**要么先给它登记 ledger 行（那是改 `scripts/` 的写面，票面 `AC#4` 的名册要认），要么把它落在 `internal/audio`（那里 windows-tagged 用例根本不在执行面，见 §2 的 ★ 行）**。这是个真取舍，⛔ 不该在实现时才想起。

---

## 4. 其余会被这批改动顶红的仪器（按射程列，⛔ 全部**未跑**，档位＝读码级判断）

| 仪器 | 逐字锚 | 会不会被顶红（＋为什么） |
|---|---|---|
| `internal/observe/thresholds.go` 的 SLO 钉 | 本腿 R20 现量：208 行、`grep -nE "Level\|Audio\|Mic\|Silence\|dBFS\|RMS\|Db"` **零命中** | ⛔ 碰不到（现量）。硬禁区仍在票面 `AC#4` 里写着"零字节" |
| `internal/ball/tokens_table_test.go:552` 逐字 `"SilenceLevelGate":  SilenceLevelGate,` | R21 现量 | **会**——只要改动让 `0.06` 的含义搬家（去直流那一支）。它是 C21 冻结 token 表的机器核对钉（票 291 `:34` 逐字指认） |
| `internal/audio/level_test.go` 的 `SineLevelTolerance` 双边钉 | 票 291 `:12` 逐字记的 `TestLevelSineToleranceIsNamedAndBounded` | 单元层：**只要 `level.go` 求和段不动就不红**（算的）；`level_test.go` 无 tag ⇒ **它在 ubuntu core 档真跑** ⇒ 这一枚是"改动真红"最快的信号（现量：core scope 含 `./internal/audio/...`） |
| `internal/ball/liquid_test.go:60` 逐字 `if m.level > SilenceLevelGate {` | R21 现量 | 球侧门限断言；`internal/ball/` 同时在 core 与 windows scope ⇒ **两档真跑** |
| `cmd/wisp/leg_sink_nail_131_windows_test.go` | §1.2 | 见 §1.2（无 tag ⇒ ubuntu 也跑） |
| `tools/d22scan`（`sh scripts/d22scan.sh`，逐字锚见 `ci.yml:137`） | R17 现量 | ⛔ 我没跑、也没读它的扫描射程 ⇒ **未证成**；只提示票面 `AC#4` 要求它 rc=0，且本仓 `.scratch/wisp/probes/**/*.go` 里躺着突变体件（30 件 §2 已具名），新增测试文件的 emoji/裸 `go func(` 形状要按 AGENTS.md §1.2 那两行的**仪器实际射程**自查 |
| `scripts/check-path-length-budget.sh`（`ci.yml:170` 逐字 `run: sh scripts/check-path-length-budget.sh --with-self-test`） | R17 现量 | ⚠ **本票的落点目录名本身就长**：`.scratch/wisp/probes/297/a1/30-known-amplitude-vs-ticket-291.md` ＝ **68 字符**（尺＝落笔后 `wc -c` 那发同批的 `awk` 我未跑，此数＝我按路径手数的**算的**，落地腿引前先现跑）。票 291 曾因文件名 111 字符把这台尺打成 VERDICT RED（票 291 Progress log `:45` 逐字记录了那次改名）。⇒ 新增件的路径长度**在这台尺的射程里**，值得先量 |

---

## 5. 本件没跑成的尺（具名）

`go test -list`（判 §3.2 的生死）、`sh scripts/d22scan.sh`、`scripts/check-path-length-budget.sh`、`go build`/`go vet`、任何 `-tags windows` / `-tags winlive` 的执行、CI 运行日志（`gh` 未查）、`probes/**` 全量名册（只看了 R8/R5 顺带命中）。
⇒ **本件对"今天 CI 什么色"零枚现量，所有档位判断都是"码层＋配置层"读数。**
