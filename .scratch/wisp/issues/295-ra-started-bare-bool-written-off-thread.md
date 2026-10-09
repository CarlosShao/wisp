# 票 295 — `ra.started` 是一枚**裸 `bool`**：手势线程在写它，boot 那句在它被读的同时可能正被改（邻居 `levels` 用的是 `atomic.Uint64`）

**立票**：2026-10-09 16:5x 编排者（来路＝票 290 非实现者验收腿 `290-v1` 的 N2／`ra.started` 一问，件 `.scratch/wisp/probes/290/v1/40-ac-verdicts.md` "需编排者裁的点 #5"；行号我 16:4x 自己复跑过）
**性质**：⚠ **这是 `290-r1` 新造的一条并发面，不是既有缺陷**——`ra.started` 这枚字段早就存在（票 247 那批），但**今天多了一枚在别的线程上写它的执行者**（静音手势）。**严重性＝低**：唯一的生产读者只在 boot 那几毫秒里读一次，撞上需要用户恰好在那一瞬按键；而 `-race` 不在本仓门禁名册 ⇒ **今天没有任何仪器看得见它**（尺见现量末条）。⛔ 不许因此不修，也⛔ 不许写成"已经出错"。
**这一格不影响票 290 任何一格成立**（腿原话），⛔ 不由票 290 顺手背。

## 现量（每条都带尺；⛔ 引用前先重跑，行号是快照）

- **字段**：`cmd/wisp/resident_audio_windows.go:105` 逐字 `	started bool`（`:102-104` 的注释逐字开头 `// started records that the gate handed the device to the inner source and` / `// capture is actually running. posture() reads it so a leg that built a`）。
- **★邻居同一族用的是原子**：`:110` 逐字 `	levels atomic.Uint64` ⇒ 同一枚 struct 里，"电平计数"有原子保护、"设备到底开没开"这枚**更贴近隐私事实**的布尔**没有**。这是本票最省事的判据形状（照邻居做，不新造规矩）。
- **写点两枚**：ⓐ boot 路径 `:323` 逐字 `		ra.started = true`（在主线程的启动 switch 里）；ⓑ **新写点** `:168` 逐字 `	ra.started = gate.Open()`，在 `toggleMute`（`:156`）体内，注释逐字 `// started mirrors what the gate now reports. It is derived state, not a` `// second truth source,…`。⇒ `toggleMute` 的执行线程＝球的手势回调所在线程，仓里逐字写着（`cmd/wisp/resident_ball_windows.go:245-246` 逐字 `// ball's global hotkeys on the ui-sta thread internal/ball owns (that thread is` / `// the ball's own Registry.Spawn with the frozen D38b resident name "ui-sta",`；`:501` 逐字 `// Threading, because it is not free: callbacks run ON the ui-sta thread and`、`:542` 逐字 `// Callbacks run ON the ui-sta thread (Ball.fire calls them inline, and its`；另 `:482` 逐字 `// ui-sta thread can never see a half-written handoff.`）。
- **唯一生产读者**：`:379` 逐字 `	if !ra.started {`，在 `posture()`（`:375`）体内；`posture()` 的生产调用者只有**一处**＝`cmd/wisp/resident_windows.go:320` 逐字 `	fmt.Printf("wisp: 采集腿：%s\n", raudio.posture())`（我 16:4x 现跑 `grep -rn "posture()" cmd/wisp/*.go`＝其余命中全在 `*_test.go`）。⇒ 时序＝`:303` 挂载**早于** `:320` 读，中间那几毫秒若有人按了 `Ctrl+Alt+M`，就是**未同步的并发读写**。
- ⚠ **已有的那枚锁不管这件事**：`resident_ball_windows.go:173` 一带逐字 `// muteMux is what makes that late attach visible to the ui-sta thread: the` ⇒ `muteMux` 保护的是"**挂载**这件事对 ui-sta 可见"（配 `resident_windows.go:303` 那次后置 setter），⛔ 它**不覆盖** `ra.started` 的写读对。
- **测试面为什么看不见**：读这枚字段的用例全在同线程直调（`cmd/wisp/resident_mute_290_windows_test.go:124`/`:141`/`:174`、`resident_audio_247_windows_test.go:113`/`:155`/`:256`、`resident_audio_247_live_windows_test.go:148`），没有一枚把它放到两个线程上对拉。
- ★**`-race` 不在门禁名册**（尺＝`grep -rn --include=*.sh --include=*.ps1 --include=*.yml -e '-race' scripts/ .github/`＝**0 命中**，我 16:4x 现跑；`grep -rn '-race' scripts/` 不加 `--include` 时只命中 `scripts/spike/bin/*.exe` 那七枚**二进制**里的串 ⇒ 不构成"有人在跑 race"）。

## 要建什么

- [ ] **AC#0 只读普查（⛔ 不许只答"没找到"）**：把 `residentAudio` / `residentBall` 这两枚 struct 里**所有"手势线程写、主线程读"的裸字段**列成名册（本票已知两枚候选：`started`、`verdict`——`:168` 只写 `started`，而 `verdict` 的写点全在 boot 路径 `:263/:284/:298/:313/:320/:324`，读点在 `:382`/`:389-390`；⇒ **本格必须现跑尺判"手势路径到底写不写 `verdict`"**，⛔ 不许照抄我这句）。判据＝每枚字段给"写点在哪个线程／读点在哪个线程"两行＋尺。
- [ ] **AC#1 修法落地**：**形状＝照邻居 `:110` 的 `atomic` 家族来**——`started` 换 `atomic.Bool`（先例与前提都现量过：`go.mod:3` 逐字 `go 1.27`，而 `atomic.Bool` 全仓已有 **5 处**在用，尺＝`grep -rn "atomic.Bool" --include=*.go internal cmd` 我 16:5x 现跑 ⇒ ⛔ 不许为这枚改动新造规矩，也⛔ 不许自己发明第三种包装）。⛔ **不加锁**（`sync.Mutex` 会把 boot 那条只读路径也拖进锁面）；⛔ **不加协程**（D38(b) 名册六枚零膨胀）；⛔ **不动 D43 状态表、不动默认值、不改门的语义**。若 `verdict` 那一枚在 AC#0 判出"也会被手势路径写"⇒ **同批改**（两枚一起，别留第二票）。
- [ ] **AC#2 判据要有牙**：ⓐ 现有那些读 `started` 的用例必须**全绿照跑**（不许放宽任何一条，含 `resident_mute_290_windows_test.go:124`/`:141`/`:174`）；ⓑ **反形正控**＝把 `:168` 换成"只在 boot 写过、手势路径不再更新"那一形 ⇒ 指名用例必须红（先例＝票 290 的 M1：只改类型不改语义的"重构"必须能被现有仪器看住）；ⓒ 若本票要造"两线程对拉"的用例 ⇒ 必须**自带正控**（同一枚台件在不并发时绿、并发时红），并写明它跑不跑 `-race` 都能判；⛔ 不许把"加了 `atomic`"本身当凭据。
- [ ] **AC#3 越界检查**：`git diff-tree -r --name-only --no-renames` 逐笔过十枚禁列（`frontend/**`／`design/**`／`PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／三枚冻结件 `internal/panel/tokens_fourway_test.go`·`internal/panel/l2_grant_boundary_test.go`·`internal/perm/ticket90_persist_test.go`）任一路径出现 ⇒ 直接退回。
- [ ] **AC#4 门禁四数**：`GOFLAGS= go build ./...` rc=0；`$(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp` 空（裸 `gofumpt` 不在 PATH＝`rc=127`，具名不掩盖）；`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ ./internal/audio/ -count=1 -v` 改前改后**各 ≥2 发取交集**并逐名作差（⛔ 单发名册不是可靠尺；三数必须带计数命令逐字，本仓沿用 `grep -c -- '--- PASS'`（含子测试））；`sh scripts/d22scan.sh` 纯净树 rc=0；每枚门禁件自落一行 `rc=N`（⛔ 0 字节件＝那格没交）；跑读数前 `tasklist` 现跑 `wisp.exe`／`balldebug.exe`＝0。
- [ ] **AC#5 〔可选加分项，排最后〕要不要把 `-race` 请进门禁**：本格**只出代价表不改门禁**——现跑一次 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -race ./cmd/wisp/ -count=1` 的**墙钟与颜色**，报"加它要多花多少、今天会不会因别腿的红而恒红"。⛔ 改 CI 那一步＝人工批准面（本票不许动 `.github/workflows/ci.yml` 一字）；判"要不要请进来"归编排者＋机主一句。

## 禁区

- ⛔ **不新造真相源**：`gate.Muted()`/`gate.Open()` 仍是主，`started` 是派生镜像（`resident_audio_windows.go:164-167` 的注释已经这么写死了）；本票动的是**这枚镜像的存储形状**，⛔ 不许顺手改成"自己维护一份静音状态"。
- ⛔ 不加锁、不加协程、不碰 D38(e) 那十步与 `internal/proc`；⛔ 裸 `go func(`（ban #1）。
- ⛔ 不改任何默认值（`internal/config/schema.go:197`/`:245`/`:277`）、不碰 D43 表、不改门的语义。
- ⛔ 不许顺手把票 290 `AC#3`、票 293、票 294 的任何格勾掉。
- git：只 commit 不 push；commit 必带**显式 pathspec**（`git mv` 的改名笔旧名与新名两个路径都要给）；禁 `add -A`／`amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内不删文件、不建 worktree。

## 排程与串行

- 写面＝`cmd/wisp/resident_audio_windows.go`（＋同包测试）⇒ ⛔ 不与票 294 同批（它可能动同一枚文件的 `executed` 那一支）、⛔ 不与票 293 的落地腿同批。
- 编排者队列：重 build `build/wisp.exe` → 票 292 → 机主真机窗口 → 票 293 `a1` → 票 294 → **本票**。
- ⚠ 本票排在最后是有意的：它是三枚里唯一**用户今天看不见错状态**的一枚，而那两枚（293 的勾、294 的仪器缺口）分别关系"隐私看得见"与"后续所有落地腿的凭据可不可信"。

**Status:** **未开工**（本票 16:5x 立，六格全未勾；⛔ 未派腿）。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
