# 票 293 — 托盘"静音"那一项现在**真会拧开麦克风**，而它自己的**勾永远不动**（`SetTrayChecks` 生产调用者＝0）

**立票**：2026-10-09 16:5x 编排者（来路＝票 290 的非实现者验收腿 `290-v1` 的 N5 一问，件 `.scratch/wisp/probes/290/v1/20-questions-n5-n6.md:3-40`；承重读数我 16:4x 全部自己复跑过）
**性质**：★**这一枚是票 290 那一发新造的"有后果的错状态"，不是既有缺陷的重新发现**。两件事必须分开写，⛔ 不许混成一件事：
- **既有（不归本票）**：`trayMuted` 自始没被人驱动过 ⇒ 默认档 boot 时门关着、菜单也不打勾，那句"静音"从第一天起就是零信息。
- **新造（归本票）**：`290-r1` 让那一项**第一次能改变设备状态**（点了真开门），而那枚勾仍永不动 ⇒ 它现在与一件**已经发生的事**并排显示。⚠ 危害面是**隐私**：用户看不出自己刚把麦克风开了。
**优先级**：三枚残余票（293/294/295）里**排第一**，理由＝唯一一枚用户会看见错状态、且方向是"少报警"（不是多报警）。

## 现量（每条都带尺；⛔ 引用前先重跑，行号与枚数都是快照）

- **那枚勾是什么**：`internal/ball/ball_windows.go:128` 逐字 `trayMuted     bool`；被读进菜单的那行 `:674` 逐字 `sel := showMenu(b.hwnd, b.trayMuted, b.trayPauseWake)` → `internal/ball/tray_windows.go:88` 逐字 `appendItem(menuMute, "静音", muted)` → `:81-82` 逐字 `if checked {` / `flags |= mfCheckd`。⇒ 这枚 bool **就是菜单上那个勾**，且它是**用户看得见静音状态的唯一一面**。
- ★**驱动它的那枚 setter 生产调用者＝0 枚**（尺＝`git grep -n "SetTrayChecks" HEAD -- internal cmd`，我 16:4x 现跑＝只命中两行：注释 `internal/ball/ball_windows.go:952`＋定义 `:953`；含测试的全仓尺同样两命中 ⇒ **调用者 0 枚，改前 0、改后仍 0**）⇒ `trayMuted` 恒为零值 `false` ⇒ **"静音"这一项永远不带勾**。
- **这一发改变了"这一项会不会动手"**：`cmd/wisp/resident_ball_windows.go:325` 逐字 `OnTrayMute:      func() { rb.muteGesture("tray-mute") },`（`290-r1` 之前是 `recordBallGesture("tray-mute")` 空落地）→ `:527` `outcome, executed := fn()` → `resident_audio_windows.go:162` `gate.SetMuted(!gate.Muted())`。
- ⚠ **console 那句不是兜底**：`cmd/wisp/resident_windows.go:55-57` 逐字开头 `// Ticket 117: the resident process is the leg owner actually uses - double` / `// click the icon, no terminal attached, stderr going nowhere.`；而 `muteGesture` 那句 `resident_ball_windows.go:536` 逐字 `fmt.Printf("wisp: ball %s: %s\n", name, outcome)` ⇒ **双击启动的形状里没有终端读者**（它进的是日志文件，不是屏幕）。
- **同罪的另一面**：热键那支 `:321` 逐字 `OnMuteHotkey:    func() { rb.muteGesture("mute-hotkey") },` 走的是同一枚 `muteGesture`，同样只留 console/log 那一句 ⇒ **两面（菜单勾／日志句）给出的确定度不一样，勾是零信息的那一面**。
- **默认值这一头不许动**：`internal/config/schema.go:277` 逐字 `MicMutedDefault bool `toml:"mic_muted_default" default:"true"``（我 16:4x 现读）。

## 要建什么（⛔ 先普查再落地）

- [ ] **AC#0 只读普查（⛔ 不许只答"没找到"）**：ⓐ 托盘**另一枚**勾 `trayPauseWake`（`ball_windows.go:674` 同一次调用送的第二个实参）今天有没有驱动者？给出与上面同形的尺。ⓑ 契约对"托盘勾号"的许诺逐字搜并摘录（`docs/PLAN.md` 的 D16/D43 那几节＋`docs/specs/SPEC-08-ui-ball-panel.md`；字样名册：`静音`/`mute`/`勾`/`check`/`checked`/`托盘`/`tray`/`指示`），并判"这枚勾该镜像门的哪一态"契约写没写。ⓒ 全仓搜有没有**第二处**已经在替托盘维护静音状态的东西（防"我另造一枚"），尺＝按语义变体名扫（`trayMuted`/`menuMute`/`mfCheckd`/`SetTrayChecks`），⛔ 只在 `scripts/` 与票池找过就等于说"没人管"。
- [ ] **AC#1 落点形状（本格不改产码，只裁"谁跟着谁"）**：写清"门侧 `muted` 为主、托盘勾为从"这一跳放在哪一侧、要动哪几枚文件。⚠ 硬约束现读：`internal/ball` 的依赖名册里**不含** `internal/audio`（票 290 的 `290-a1` 量过，本票若要动它必须自己复跑 `GOFLAGS= go list -deps ./internal/ball`）⇒ 那一跳**必须由 `cmd/wisp` 注入**（照 `resident_windows.go:303` 后置 setter 的既有形状），⛔ 不许为让勾动起来就新增 `internal/ball → internal/audio` 的包级依赖边。
- [ ] **AC#2 落地那一发（判据必须带牙）**：形状＝门态一变，勾立刻镜像（点击拧门／热键拧门／boot 出厂静音位三种起点都要对）。⛔ **不许新造第三枚真相源**（勾必须读回 `gate.Muted()`，不许另存一枚"我以为的静音"）。判据写成两半：ⓐ **行为断言**＝点那一项之后 `gate.Muted()` 与菜单送进 `appendItem` 的那个 bool 同值；ⓑ **反形正控**＝把镜像那一行换成空操作（`SetTrayChecks(false, …)` 恒假）⇒ 指名用例必须红（先例＝票 290 的 M1：光有"调用者枚数"这种尺不算牙）。逐名照抄终态、每枚门禁件自落一行 `rc=N`。
- [ ] **AC#3 越界检查**：`git diff-tree -r --name-only --no-renames` 逐笔过十枚禁列（`frontend/**`／`design/**`／`PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／三枚冻结件 `internal/panel/tokens_fourway_test.go`·`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`／`testdata/golden`）任一路径出现 ⇒ 直接退回。⛔ 若要把球打进 D43 的 `Muted` 态（灰球＋斜杠）＝**那不是本票**，那要先落一枚人工批准的 `A##`（票 290 裁定段已量：表里**没有 `Sleeping → Muted`** 那条边）。
- [ ] **AC#4 门禁四数**：`GOFLAGS= go build ./...` rc=0；`$(go env GOPATH)/bin/gofumpt.exe -l <自己动过的目录>` 空（⚠ 裸 `gofumpt` 不在 PATH＝`rc=127`，具名不掩盖）；`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ ./internal/audio/ ./internal/ball/ -count=1 -v` 改前改后各一次并**逐名作差**（⛔ 单发名册不是可靠尺：本仓 10-09 实测同一棵已入库的树两发之间会 ±1 枚；三数必须带**计数命令逐字**，本仓沿用的那把是 `grep -c -- '--- PASS'`（含子测试））；`sh scripts/d22scan.sh` 纯净树 rc=0；**跑任何真机/整包读数前先 `tasklist` 现跑 `wisp.exe`／`balldebug.exe`＝0**（票 290 这一波就是被一枚活的 `build\wisp.exe` 多顶出一枚 `TestAC246DevLegIgnoresTheTestTaskInjection` 才学到这条）。
- [ ] **AC#5 肉眼那一形（⛔ 归机主在场，不算腿的欠账）**：托盘菜单上勾的真机读数＝点开→勾在、再点→勾没，两形各一次；本仓今天没有能造"托盘菜单可见态"的仪器（尺＝`grep -rn "showMenu\|mfCheckd\|appendItem" --include=*_test.go internal cmd`＝**0 命中**，我 16:5x 现跑），⛔ 腿不许拿"读回那个 bool"冒充"用户看得见勾"。

## 禁区

- ⛔ **不改默认值**（`schema.go:197`/`:245`/`:277` 三枚一字不动）；本票动的是"勾跟不跟门"，不是"门开不开"。
- ⛔ **不许新造第二/第三枚真相源**：门的 `muted` 是主（`internal/audio/gate.go:20-21` 逐字写死的），球态与托盘勾都是从；⛔ 不许把主从倒过来。
- ⛔ 不许新增包级依赖边（尤其 `internal/ball → internal/audio`）；⛔ 裸 `go func(`（ban #1）；⛔ 不为镜像新造协程（D38(b) 名册六枚零膨胀）。
- ⛔ 不许顺手把票 290 `AC#3` 那格（欠机主在场真机）勾掉，也不许顺手改球侧 `Muted` 态或 D43 表。
- git：只 commit 不 push；**commit 必带显式 pathspec**（`git mv` 那类改名，旧名与新名**两个路径都要给**——10-09 实测只给新名会把删除侧留在索引里被下一笔顺手收走）；禁 `add -A`／`amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内不删文件、不建 worktree。

## 排程与串行

- 两拍：`293-a1`（只读普查＝AC#0＋AC#1）→ 编排者裁 → `293-r1`（落地＝AC#2..AC#4）。
- 写面＝`cmd/wisp`＋`internal/ball` ⇒ ⛔ 与票 290 残余真机窗口**不混批**（那一次窗口只验开门两形＋票 247 三格＋票 291 四形）；⛔ 与票 294（同写 `cmd/wisp` 测试面）串行，两枚不能同时在飞。
- 队尾现状（编排者排）：重 build `build/wisp.exe` → 票 292（三行注释小发）→ **机主那一次 7-8 分钟真机窗口** → 本票 `a1` → 票 294 → 票 295。

**Status:** **未开工**（本票 16:5x 立，六格全未勾；⛔ 未派腿）。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
