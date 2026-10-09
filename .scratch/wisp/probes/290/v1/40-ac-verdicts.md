# 290-v1 · 40 裁决表（与 AC 编号 1:1，⛔ 本腿零枚框翻动）

腿＝`290-v1`（非实现者验收）。被裁＝`290-r1` 三笔 `fdee536b`／`69d8c9ef`／`75128ae5`。
票面＝`.scratch/wisp/issues/290-gate-mounted-nobody-can-open-it.md`（未勾 4 枚＝`AC#2`..`AC#5`，已勾 2 枚＝`AC#0`/`AC#1`，本腿现跑 `grep -c`＝4／2）。
`AC#0`/`AC#1` 已归编排者裁（`:54-66` 那一节），本腿不重裁、只列名以免缺行。

| AC | 判语 | 凭据（`文件:行` 逐字 或 本腿实跑读数） |
|---|---|---|
| **AC#0**（契约普查，只读） | 〔已由编排者裁定并翻勾，本腿不重裁〕 | 票面 `:29`＝`- [x]`；裁定段 `:54-66`；件 `.scratch/wisp/probes/290/a1/00-findings.md`（38,744 字节，本腿现跑 `stat` 已核存在） |
| **AC#1**（三形代价表） | 〔同上，`[x]` 在 `:30`；本腿只复核了甲-1 的落地形状与其一致〕 | 落地形状＝甲-1：本腿现读三枚文件全在 `cmd/wisp`（`00-anchor…` §0.7 逐笔名册）＋`WithGateEvents` 产码调用者 0 枚复现（尺见 `20-questions-n5-n6.md` 末表） |
| **AC#2 落地那一发** | **【成立】** | ① 票面尺（对象层，⛔ 工作树）：`git grep -n "SetMuted(\|SetSpeaking(" 6c976aa6 -- internal cmd \| grep -v _test` ＝**3**（`gate.go:97` 注释＋`:130`/`:168` 两枚定义）→ 同尺在 `HEAD` ＝**4**，多出的那枚逐字 `HEAD:cmd/wisp/resident_audio_windows.go:162:	gate.SetMuted(!gate.Muted())`；**调用者枚数 0→1**。② 那一枚"在不生产链路上"：本腿逐跳现读到 `internal/ball/ball_windows.go:660 case hkMute:` → `:661 b.fire(b.opts.Events.OnMuteHotkey)` → `cmd/wisp/resident_ball_windows.go:321` → `:515 muteGesture` → `:527 fn()` → `resident_windows.go:303 rb.attachMuteGate(raudio.toggleMute)` → `resident_audio_windows.go:156/:162`（全表＝`00-anchor…` §0.4）；门本体是 `:273 NewHalfDuplexGate(...)`、`:277 ra.mic, ra.gate, ra.cancel = ...` 装配出来的那枚。③ ⛔ 不是又一枚调试 `cmd`：命中文件＝`cmd/wisp`（常驻本体），`cmd/balldebug` 零命中（`EvMuteKey` 尺反证：唯一产码发射者仍是 `cmd/balldebug/main.go:603`）。④ **牙**＝M1 变异：把那行改成空拧（`gate.SetMuted(gate.Muted())`）⇒ `TestAC290BothMuteGesturesTurnTheGateAndBack`／`…OutcomeIsReadOffTheGateNotOffTheRequest`／`…MutedDefaultStillDecidesTheBoot` 三枚红（`rc_m1=1`，断言行 `:119`/`:178`/`:265`），同一次那把枚数尺**仍报 1** ⇒ 见下面"需编排者裁 #1/#2" |
| **AC#3 默认档一字不改＋两形真机读数** | **【不可判(欠机主在场那一发)】** ⛔ 不据此判失败 | 形① 的本腿独立现跑（装配级，四枚出厂表尺）：`go test ./cmd/wisp/ -run 'TestAC247…' -count=1 -v` → `rc_form1=0`、`--- PASS` x4；另六枚 `TestAC290*` 基线 `rc_targeted_baseline=0`。三枚默认值现读一字未动（`internal/config/schema.go:197 default:"false"`／`:245 default:"true"`／`:277 MicMutedDefault … default:"true"`），boot 那句仍逐字在 `cmd/wisp/resident_audio_windows.go:298`（`采集腿已装配、门处于静音：设备未打开（[audio] mic_muted_default=true，来源 %s）`）。⛔ **形② 不可判**：判据是"用户主动开门之后 `SetAudioLevel` 真的收到随声音变化的数"，那需要机主在场对麦说话；本腿⛔ 未跑 `WISP_LIVE_MIC=1`（`resident_audio_247_live_windows_test.go:129-130` 的 `t.Skip` 在两发整包里都记为 `--- SKIP`，本腿没把它读成读数），⛔ 未用外放/合成源/临时配置冒充。欠账归编排者（那一发的命令与逐字期望句在 `20-ac-readings.md:99-125`，本腿核其引用的行号属实）。⚠ 另报一条现场前置：票面 `AC#3` 要求"跑真机前 `wisp.exe` 必须为 0"，本腿现跑＝**1 枚在跑**（PID 32888＝`build\wisp.exe`，16:10:05 起），见 `30-gates-and-rosters.md` §3 |
| **AC#4 越界检查** | **【成立】** | ① 逐笔尺（⛔ 未用区间 diff 作判据）：`git diff-tree -r --name-only --no-renames` 对三笔＝4＋1＋18 枚，过票面十枚禁列（含 `frontend/`/`design/`/`docs/PLAN.md`/`docs/specs/`/`thresholds.go`/`gate.go`/`internal/ball/`/三枚冻结件/testdata/golden）＝**hits 0、rc=1** 三笔各自。② 区间尺作旁证：`git diff --no-renames --name-only 6c976aa6..HEAD -- internal cmd` ＝**恰好那 5 枚文件**（别腿在区间内零 `.go`），且禁列区间尺＝**空（零字节）**。③ 票面动作尺：`git diff-tree --numstat -r 75128ae5 -- <票面旧长名>` ＝ **`1 0`**（＋1 行 Progress log／−0 行），`^- [ ]`＝4、`^- [x]`＝2 ⇒ **零 AC 框被翻、原句零改**。④ 状态机面：三枚产码件新增行内 `statemachine.` 计数＝0/0/0、`EvMuteKey`＝0（表边仍只有 `internal/statemachine/table.go:65/:75/:79`）⇒ D43 表零触碰。⑤ 无新配置键、无持久化改动（工作树 `internal cmd` 干净、`schema.go` 零字节）。⑥ `frontend/**`/`design/**` 零读零写（名册零命中；那 16 枚 ` D design/**` 是别人的在飞件，本腿与实现者都未碰） |
| **AC#5 门禁四数** | **【成立】**（附两处具名口径欠账，⛔ 不构成本格失败） | 本腿独立现跑：`GOFLAGS= go build ./...` **`rc=0`**；`sh scripts/d22scan.sh` **`rc=0`**、`clean`、射程逐字与自报一致（`internal/=524`、`cmd/=117`）；`$(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp` **空、`rc=0`**；裸 `gofumpt` **`rc=127`（具名不掩盖）**；`go vet ./cmd/wisp/` **`rc=0`**；`gofmt -l cmd/wisp` **非空 3 枚**＝`models.go`/`panel_inbound_guards_35r3_test.go`/`panel_transport_35r2_test.go`（三枚**都不在本票名册**，`git log -1` 分别 `5e8748b3 10-03`/`7f9d6e40 10-07`/`3a343bc7 10-08`，`gofmt -d` 现读是 CRLF 整件重写形；实现者用的口径是"五枚自己动过的件"，本腿逐字复跑＝**空、`rc=0`**）⇒ **票面那条尺原文写的是 `gofumpt -l <目录>`，目录口径下本票过**。整包两发（票面**两枚包**）：本腿 1 `423/6/3 rc_run1=1`、本腿 2 `423/6/3 rc_run2=1`、**同树两发红名册 diff＝零行（`rc=0`）**、`0xc0000135`＝0、`0xc000013a`＝0；与实现者 `rosters/raw-post-final-{1,2}.red.txt` 作差＝**多 1 枚** `TestAC246DevLegIgnoresTheTestTaskInjection`，其自身逐字判语 `resident_task_source_246_windows_test.go:432-435`＝"AC#7 RED (**desktop state, not our code**): a dev Wisp is already running…"，本腿 `tasklist` 现跑证实会话里有活的 `build\wisp.exe` ⇒ **本票新增红＝0 枚**（判据是逐名作差，⛔ 不是"全绿"）。名册侧另核：实现者四发 raw 的三数（418/5/3、423/6/3、424/5/3 x2）本腿在其 raw 件上逐字复算**口径一致**（`grep -cE "^[[:space:]]*--- PASS"` 含子测试），且 `raw-post-1` 多出的那枚正是 `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`（与 §7 自报的漂移归因同一条）。**具名口径欠账两处**：① 三数口径（含子测试 vs 顶层）未在票面钉死；② `raw-post-2.md`＝**0 字节**（见下） |

## `raw-post-2.md`（0 字节）那一格的单独判语

**【成立】（＝那一格没欠）**：它是**中间树**（`fdee536b` 态、cite 未修）那**一发被自己中止**的残件，
不是任何 AC 格的唯一凭据；同格（中间树改后首发）由 `raw-post-1.md`（321,015 字节）＋`rosters/raw-post-1.red.txt`（6 枚红）交付，
而票面 `AC#5` 要求的"改前改后各一次"由 `raw-pre-1/2.md`＋`raw-post-final-1/2.md`＋五枚名册交付（本腿逐件 `stat`＋三数复算）。
⛔ 本腿未删该件。⚠ 形状欠账具名：0 字节件**自带零自证**（"中止"只活在 `30-gates.md:35` 的叙述里）——
建议台件纪律一条：中止那一发也应留一枚**非 0 字节、写明中止原因**的件（不阻塞本票）。

## 我没收下、需要编排者裁的点（编号沿用派单）

1. **`AC#2` 那条尺的字面分母有缺陷（票面文字面就通不过 M1）**：票面 `:35` 写"由 **2** 枚变 ≥3 枚"。
   本腿两把都跑：合尺在起手锚＝**3**（不是 2）、单尺（只 `SetMuted(`）＝**2**（"2"只在单尺上成立）。
   实现者已具名报回（`20-ac-readings.md:37-40`）。**需裁**：结案时把 `AC#2` 的判据文字订正为"调用者 0→1 且新命中在生产链路上"（＋行为断言），
   ⛔ 票面文字改动＝人工批准，本腿与实现者都未动一字。
2. **枚数尺必要不充分（M1 实测）**：一枚 `gate.SetMuted(gate.Muted())` 满足票面 AC#2 的尺却一格不拧门。
   ⇒ 建议 `AC#2` 结案判语显式带上"由 `resident_mute_290_windows_test.go` 的三枚行为断言背书"。
3. **M2 实测出的覆盖缺口（新，非实现者自称的形状）**：门真拧时 `executed` 谎报 `false` ⇒ **六枚全绿**（`rc_m2=0`）。
   happy path 上无人断言 `executed==true`；后果只在日志面（`resident_ball_windows.go:531` 会打 `mute gesture found no gate to turn`），
   用户可见那句仍是门读回的真话 ⇒ 危害低，但这是**那一枚函数值契约的另外一半**，建议落票。
4. **N3 缺口（新）**：装配根那一跳（`resident_windows.go:303`）与球侧两支闭包**没有生产形状用例**，
   而仓里**躺着两件现成仪器没用**——`resident_hotkey_258_windows_test.go:278` 的真球窗＋`:289-291` 的 `HotkeyReport`（且 `:273` 逐字已经写了 `mute = "Ctrl+Alt+M"`），
   以及 `resident_hotkey_258_test.go:25-68` 的 `runResident` AST 走查先例（同类先例共三枚）。
   `attachMuteGate` 的三枚测试命中**全是直调、零走查**（尺见 `10-mutations…` §4）。⇒ 建议落票（含具体台件写法）。
5. **`ra.started` 的新跨线程裸写（新，本腿现读到）**：`resident_audio_windows.go:168` 在 ui-sta 线程写裸 bool `started`（`:105`；邻居 `levels` 是 `atomic.Uint64` `:110`），
   唯一生产读者 `resident_windows.go:320` 在挂载（`:303`）**之后**才读 ⇒ boot 那几毫秒内是未同步并发读写（`-race` 不在门禁名册，今天量不到）。
   ⛔ 不影响任何一格成立；建议落票（改 `atomic.Bool`，⛔ 不加锁、⛔ 不加协程）。
6. **`SetTrayChecks` 生产调用者＝0 枚（本腿全仓尺复现）** ⇒ 托盘"静音"项现在**真拧门而勾永不动**（N5 判语：勾不准＝既有，"有后果"＝本次新造）。
   本腿⛔ 未修。⇒ 建议落票（名册第一枚）。
7. **`gofmt -l cmd/wisp` 在 HEAD 非空（三枚别腿旧件、CRLF）**：不是本票账、也不触票面尺（票面尺是 gofumpt）。
   ⇒ 要不要单立一枚"行尾归一"的小票＝编排者的裁；⛔ 不该由 290 背。
8. **实现者具名欠的三格本腿都收为"非其账"**：`AC#3` 形②（欠机主在场）、`internal/panel` 包级名册（票面尺只有两枚包，本腿也没扩——理由更强：工作树里那 16 枚 ` D design/**` 是在飞件，扩包集会量到别人的红）、托盘镜像（裁定明写本发不做）。
9. **票面"现量"节的行号今天全部漂移**（`:213`→`:273`、`:281`→`:321`、`:238`→`:298`、`internal/ball/ball_windows.go:669-682` 段仍对得上）。
   按 `AC#4` 规矩实现者⛔ 未改原句 ⇒ 要不要在票面追加"行号已漂移"一节＝票面动作，归编排者。
10. **派单本身两处与盘上不符**（本腿以原文为准并具名）：`case hkMuted:` 实为 `case hkMute:`（`ball_windows.go:660`）；
    `AC#5` 的整包尺在票面原文里是**两枚包**（`./cmd/wisp/ ./internal/audio/`），派单只写了 `./cmd/wisp/`——本腿按原文跑了两枚包。

## 新票建议名册（交编排者立，⛔ 本腿一枚未动）

1. `29x-tray-checkmark-never-follows-the-gate-the-mute-item-now-acts-while-its-own-check-stays-unchecked`（N5，最坏形状已给）
2. `29x-executed-flag-on-the-happy-path-is-asserted-by-nobody-m2-green`（M2）
3. `29x-no-production-shape-test-for-the-assembly-root-attach-or-the-two-ball-closures`（N3，含两件现成仪器的名字与行号）
4. `29x-ra-started-is-written-from-the-ui-sta-thread-while-the-boot-report-reads-it`（并发面，`atomic.Bool` 那一发）
