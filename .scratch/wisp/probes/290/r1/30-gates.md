# 票 290 / 290-r1 — 30 门禁终态（AC#4／AC#5，逐名照抄终态，每枚门禁件自落一行 `rc=N`）

树＝本腿两笔产码之后：`fdee536b`（开门产码＋测试）→（其间别腿两笔 `beaac644`/`75620005` 只动
`.scratch/wisp/probes/orch/**`，⛔ 零 `.go`）→ `69d8c9ef`（cite 行号修复）＝起手 HEAD。
每个 `rc` 都由**被量命令自己的** `$?` 取（⛔ 不是管道尾）：需要过管道的地方写成
`… | grep … ; echo "rc=$?"` 时，那一行的 `rc` 是 `grep` 的**命中有无**（`rc=1`＝零命中＝我要的终态），
不是被量命令的失败；纯命令一律直接 `$?`。派单点名的这一刀我在自己那份基线里也踩过，逐处标明读的是谁的码。

## 1. AC#5 门禁四数（＋两枚派单点名的形状）

| 门禁件 | 命令（逐字） | 终态读数 | rc 取自 |
|---|---|---|---|
| 建 | `GOFLAGS= go build ./...` | 无输出（全仓可编） | `rc_build=0` |
| D22 闸 | `sh scripts/d22scan.sh`（⛔ `go run ./tools/d22scan` 必失败：独立模块，`scripts/d22scan.sh` 头部注释逐字写明） | `d22scan: clean - no D22 ban violations`；射程逐字：`bans #1-5 internal/=229, bans #1-5 cmd/=39, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=524, ban #8 cmd/=117`（`cmd/` 由 289 时的 116→**117**＝本腿新增那枚 `_test.go`） | `rc_d22scan=0`，原始件 `raw-d22scan-final.md` |
| gofmt | `gofmt -l cmd/wisp/resident_audio_windows.go cmd/wisp/resident_ball_windows.go cmd/wisp/resident_windows.go cmd/wisp/config_readers_255.go cmd/wisp/resident_mute_290_windows_test.go` | 空（未列任何文件） | `rc_gofmt=0` |
| gofumpt（裸跑，派单点名的形状） | `gofumpt -l cmd/wisp/resident_ball_windows.go` | **不在 PATH**，命令未执行 | `rc_bare_gofumpt=127`（复现派单说的这一处，⛔ 不掩盖） |
| gofumpt（用的那一支） | `$(go env GOPATH)/bin/gofumpt.exe -l <上面五枚件>` | 空 | `rc_gofumpt=0` |
| 测试（改前） | `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -v ./cmd/wisp/ ./internal/audio/ -count=1` ×2 | 见 §2 表 | 每发 `rc=1`（5 枚在飞红，见 §2 逐名） |
| 测试（改后） | 同上 ×2 | 见 §2 表 | 每发 `rc=1`（**同一组** 5 枚） |
| 靶向（本票六枚＋被推动的那枚） | `go test ./cmd/wisp/ -run 'TestTicket255Roster\|TestAC290' -count=1 -v` | `--- PASS` ×8：六枚 `TestAC290*` ＋ `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` ＋ `TestTicket255RosterStillMatchesTheActualReadSites`；`ok cmd/wisp 0.571s` | `rc_targeted=0` |
| AC#3 第①形默认档四枚尺 | `go test ./cmd/wisp/ -run 'TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice\|TestAC247ShippedDefaultsAreTheOnesThisLegReads\|TestAC247HandingTheLevelToTheSeamIsNotVisibility\|TestAC247UnmuteReachesTheBallSeam' -count=1 -v` | `--- PASS` ×4、`ok 0.100s`（凭据叙述见 `20-ac-readings.md`） | `rc_form1=0` |
| vet（本腿自加，非票面要求） | `GOFLAGS= go vet ./cmd/wisp/` | 无输出（含 copylocks：`residentBall` 新增 `sync.Mutex` 后全仓仍零值拷贝） | `rc_vet=0` |

## 2. 改前／改后整包名册：三数并排 ＋ 逐名作差（每包两发）

计数尺＝`grep -c -- '--- PASS' / '--- FAIL' / '--- SKIP'`（含子测试；同一支 harness，只带 `-v`）。
⚠ 两枚 `0xc0000135`／`0xc000013a` 计数每发都单独量（`grep -c`），⛔ 未把环境红算进被验物：
四发全 0（用了 sherpa DLL 前缀那一支 harness）。

| 发次 | 原始件 | PASS | FAIL | SKIP | rc | `cmd/wisp` 包级行 | `internal/audio` 包级行 | 135 | 13a |
|---|---|---|---|---|---|---|---|---|---|
| 改前 1 | `raw-pre-1.md` | 418 | 5 | 3 | 1 | `FAIL … 464.624s` | `ok … 15.945s` | 0 | 0 |
| 改前 2 | `raw-pre-2.md` | 418 | 5 | 3 | 1 | `FAIL … 465.462s` | `ok … 16.226s` | 0 | 0 |
| 改后（中间树＝`fdee536b`，一发） | `raw-post-1.md` | 423 | **6** | 3 | 1 | `FAIL … 449.140s` | `ok … 16.314s` | 0 | 0 |
| 改后·中间树第二发（被我主动中止，故 0 字节） | `raw-post-2.md` | — | — | — | 中止 | — | — | — | — |
| 改后终态 1（树＝`69d8c9ef`） | `raw-post-final-1.md` | **424** | 5 | 3 | 1 | `FAIL … 459.789s` | `ok … 15.914s` | 0 | 0 |
| 改后终态 2 | `raw-post-final-2.md` | **424** | 5 | 3 | 1 | `FAIL … 461.112s` | `ok … 15.961s` | 0 | 0 |

逐名作差尺＝`grep -- "^--- FAIL" | sed -E 's/^--- FAIL: ([^ ]+).*/\1/' | sort` 落 `rosters/*.red.txt` 后 `diff`：

| 作差 | 命令 | 读数 |
|---|---|---|
| 改前 1 vs 改后终态 1 | `diff rosters/raw-pre-1.red.txt rosters/raw-post-final-1.red.txt` | **零行**，`rc_diff_pre1_post1=0` ⇒ 新增红 0 枚、转绿 0 枚 |
| 改前 2 vs 改后终态 2 | `diff rosters/raw-pre-2.red.txt rosters/raw-post-final-2.red.txt` | **零行**，`rc_diff_pre2_post2=0` |
| 改后终态 1 vs 改后终态 2（同树两发自身噪声尺） | `diff rosters/raw-post-final-1.red.txt rosters/raw-post-final-2.red.txt` | **零行**，`rc_diff_post1_post2=0` ⇒ 本票这两发没有 289-v1 记过的那枚 ±1 抖动 |
| PASS 三数 | 418 → 424 | ＋6＝本腿新增的六枚 `TestAC290*` 全绿；⛔ 无任何旧测试转绿（不可能由注释/文案产生） |
| 那 5 枚在飞红逐名（改前后同名同时长量级） | `rosters/raw-post-final-1.red.txt` | `TestPanelHostRealWindowHopAndLifecycle`、`TestAC4FocusReturnToPriorWindowGap33r5`、`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`、`TestAC14AwaitedBindingReplyReachesThePage`、`TestAC14GoSideEvalPushReachesThePage`（全 5 枚＝面板/WebView2 族，改前两发即红 ⇒ ⛔ 不是本票账） |
| SKIP 逐名（两态同三枚） | `grep -- '--- SKIP'` | `TestPanelHostLatencyPercentilesAC2`、`TestAC247LiveMicrophoneLevelsReachTheBallSeam`（`WISP_LIVE_MIC=1` 才跑）、`TestLiveWasapiSmoke` ⇒ ⛔ 本腿没把任何 SKIP 当读数、⛔ 未加 `-skip` 变绿、⛔ 未放宽任何断言 |

### 2.1 中间树那一发的 6 枚红＝本腿自己的账（具名＋修法）

`diff rosters/raw-post-1.red.txt rosters/raw-post-final-1.red.txt` ⇒
`6d5 < TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`（`rc_diff_intermediate=1`＝差集非空，读的是 `diff` 的码）：
本腿在 `resident_audio_windows.go`/`resident_ball_windows.go` 的插入把票 255 roster 的
`file.go:LINE [token]` 证据行推走，那枚"cite 必须仍说真话"的尺翻红。修法＝笔 `69d8c9ef`
**只更新三枚行号**（`:196→:256`、`:197→:257`、`:276→:316`），token/verdict 文字未动、
⛔ 未动任何断言、⛔ 未改 255 的句子措辞。叙述与先例见 `10-wiring.md` §7。
⇒ 终态两发（树＝`69d8c9ef`）该枚已绿，红名册与改前**逐名同集**。

## 3. AC#4 越界尺（逐笔，⛔ 用区间 diff——289-v1 已具名那条方法论坑：区间会把别腿的台账笔算到我头上）

| 尺 | 命令 | 读数 |
|---|---|---|
| 逐笔名册 | `git show --name-only --format= fdee536b` | 4 枚：`cmd/wisp/resident_audio_windows.go`、`cmd/wisp/resident_ball_windows.go`、`cmd/wisp/resident_mute_290_windows_test.go`、`cmd/wisp/resident_windows.go` |
| 逐笔名册 | `git show --name-only --format= 69d8c9ef` | 1 枚：`cmd/wisp/config_readers_255.go` |
| 禁区命中 | 上面两串各自管道给 `grep -cE "internal/speech\|scripts/spike\|frontend/\|design/\|docs/\|thresholds\.go\|golden\|allowlist\.txt\|tokens_fourway_test\|l2_grant_boundary_test\|ticket90_persist_test\|statemachine/table\.go\|config/schema\.go"` | 两笔都是 **0 枚**（`rc=1` 读的是 `grep` 无命中，⛔ 不是命令失败）⇒ 票面 AC#4 那十枚路径零命中；`internal/ball`、`internal/statemachine`、`internal/config` **一枚未动** |
| numstat 逐枚 | `git show --numstat --format= fdee536b` / `… 69d8c9ef` | `60/0`、`142/14`、`287/0`、`25/0` ＋ `7/7` ⇒ 合计 **514 插入 / 14 删除**（产码笔）＋ 7/7（cite 修复），全在 `cmd/wisp` |
| 三枚默认值相对 HEAD 零字节 | `git diff fdee536b~1..HEAD -- internal/config/schema.go internal/statemachine/table.go internal/observe/thresholds.go internal/audio/gate.go` | 空（本腿没碰这四枚文件一字；门的读口/写口也未扩，`gate.go` 保持原样） |
| 工作树 Go 面是否只含我的笔 | `git status --short -- internal cmd` | 空 ⇒ 我这两笔之后工作树的 `internal/**`＋`cmd/**` 无未提交改动（别腿在飞的 `design/**`、`.gitignore`、`.scratch/**` 一律未还原、未提交、未据以判绿） |
| 文件是否存在只认对象层 | `git cat-file -e HEAD:design/assets/tokens.css` | `rc=1`（`design/assets/tokens.css` 已被别腿在工作树/后续笔删除）⇒ `TestC21DesignTokensFourWayAgree` 若红是那枚文件的账，⛔ 不是本票；本票两发整包红名册里没有它（它在 `internal/panel`，本腿没进那包） |

## 4. 硬约束逐条自量（违一条即退回，故逐名给尺）

| 约束 | 尺 | 读数 |
|---|---|---|
| ⛔ 零新协程／裸 `go func(` | d22scan ban #1 射程 `internal/=229 cmd/=39`，`clean`；另尺＝`grep -n "go func(" cmd/wisp/resident_{ball,audio}_windows.go cmd/wisp/resident_windows.go` | 三枚件内**零** `go func(`；采集线程仍只由既有 `audio.SpawnCapture` 在既有 `audio-capture` 名下起（测试尺 `reg.CountByName("audio-capture")==0` 在静音位钉住，见 `10-wiring.md` §5） |
| ⛔ D38(e) 十步顺序一字不动 | `git diff fdee536b~1..HEAD -- internal/proc` | 空；本腿只在既有 `proc.StepStopAudio`（247 注册的钩子）之外读门，⛔ 未新增步、未新增钩子名 |
| ⛔ 不扩"拒绝启动" | 本腿三枚件内零 `os.Exit`／零新增 refusal 分支（`grep -c "os.Exit" ` 三枚件＝我未加行）；装配仍只有票 128 那一种拒绝 | 成立 |
| ⛔ 跨过 `internal/ball` 边界的只许一枚标量 | 本腿未动 `internal/ball` 一字（§3 逐笔名册）；送进球的仍是 `residentAudio.levelOut`→`rb.setAudioLevel(float32)`；静音那一跳送过去的是一枚**函数值**（`muteGestureFunc`），⛔ 零 PCM／零 samples 过边 | 成立 |
| ⛔ 门侧为真相源、零第二枚布尔 | `grep -n "muted\|isMuted" cmd/wisp/resident_ball_windows.go` 只命中注释里的 `muteGate muteGestureFunc`（一枚函数值）与 `gate.go:20-21` 的引文；`ra.started` 逐字注释为 `gate.Open()` 的派生镜像 | 成立 |
| ⛔ 不碰别人的在飞活 | 我的两笔 pathspec 显式列 5 枚件；`git show --stat` 逐枚可数（§3）；工作树 16 枚 `design/` 删除＋`.gitignore`＋`.scratch/**` 在飞件：⛔ 未还原、⛔ 未提交 | 成立 |
| ⛔ 零读零写 `frontend/**`/`design/**`、本票不涉凭据 | §3 逐笔名册零命中；本件与三份证据件内零凭据字符串 | 成立 |
| Git 纪律 | 四笔 commit（骨架 1＋产码 1＋cite 1＋证据 1），每笔显式 pathspec；⛔ 零 `add -A`／`git add .`、⛔ 零 `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`、⛔ **零 push**（`git log dev..origin/dev` 与本地笔数由编排者核） | 成立 |
| 票面纪律 | `.scratch/wisp/issues/290-…md` 只在 `## Progress log` 末追加一条；⛔ 票面原句零改动；`AC#2`/`AC#3`/`AC#4`/`AC#5` 四框**一枚未勾**（框等非实现者裁） | 成立 |
| 真机前置进程数 | `tasklist //FI "IMAGENAME eq wisp.exe" \| grep -c "wisp.exe"` → `0`；`… balldebug.exe` → `0`（`rc=1` 读的是 grep 无命中） | 两枚 0；⛔ 本腿全程未拉起 `wisp.exe`（没注册全局热键、没开这台机器的麦） |

## 5. 本腿具名欠的（⛔ 不自行翻勾）

1. **AC#3 第②形真机读数**（用户主动开门后 `SetAudioLevel` 收到随声音变化的数）＝要机主在场对着麦说话。
   仪器、命令、期望逐字句、采样出处都在 `20-ac-readings.md`；那一格欠账归编排者（与票 247 的
   `AC#2`/`AC#4`/`AC#6` 三格＋票 291 `AC#2` 四形同一次真机窗口批发；别腿 `livewin-1` 的操作单
   `.scratch/wisp/probes/orch/2026-10-09-live-window-runsheet.md` §3「C 族」按票面等我这一发，
   本腿产码已入库＝`fdee536b`，可按那单执行）。
2. **AC#3 第①形的真机 console 凭据**：本腿证到装配级读数（真源构造＋真门＋出厂表），
   ⛔ 没跑那一次真实 `wisp.exe` 双击形状（理由同上，机主不在场时不开麦、不注册全局热键）。
3. **`internal/panel` 包级名册**：票面 AC#5 与派单第 5 段的尺只有 `./cmd/wisp/ ./internal/audio/` 两包，
   本腿⛔ 未自行扩包集；派单让 `00-anchor` 记 `internal/panel` 的那一列**具名欠**（已在 `00-anchor-and-baseline.md` §6 写明）。
4. **托盘勾号 `trayMuted` 不跟门**（本发按裁定不镜像）＝用户开门后托盘菜单勾号仍不显示，那一格归后续票。
5. 非本腿、非我造成的红，具名带一句：别腿 `livewin-1` 在 `75620005` 记的
   `scripts/check-path-length-budget.sh` VERDICT RED 红在**票 290 的票名 106 字符**（>README 规矩 9 的 100 帽），
   改票名＝改名＝票面动作，⛔ 不由写腿做，交回编排者。
