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

- [x] **AC#0 只读普查（⛔ 不许只答"没找到"）**：ⓐ 托盘**另一枚**勾 `trayPauseWake`（`ball_windows.go:674` 同一次调用送的第二个实参）今天有没有驱动者？给出与上面同形的尺。ⓑ 契约对"托盘勾号"的许诺逐字搜并摘录（`docs/PLAN.md` 的 D16/D43 那几节＋`docs/specs/SPEC-08-ui-ball-panel.md`；字样名册：`静音`/`mute`/`勾`/`check`/`checked`/`托盘`/`tray`/`指示`），并判"这枚勾该镜像门的哪一态"契约写没写。ⓒ 全仓搜有没有**第二处**已经在替托盘维护静音状态的东西（防"我另造一枚"），尺＝按语义变体名扫（`trayMuted`/`menuMute`/`mfCheckd`/`SetTrayChecks`），⛔ 只在 `scripts/` 与票池找过就等于说"没人管"。
- [x] **AC#1 落点形状（本格不改产码，只裁"谁跟着谁"）**：写清"门侧 `muted` 为主、托盘勾为从"这一跳放在哪一侧、要动哪几枚文件。⚠ 硬约束现读：`internal/ball` 的依赖名册里**不含** `internal/audio`（票 290 的 `290-a1` 量过，本票若要动它必须自己复跑 `GOFLAGS= go list -deps ./internal/ball`）⇒ 那一跳**必须由 `cmd/wisp` 注入**（照 `resident_windows.go:303` 后置 setter 的既有形状），⛔ 不许为让勾动起来就新增 `internal/ball → internal/audio` 的包级依赖边。
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

## 编排者裁定（2026-10-09 20:2x，来路＝只读普查腿 `293-a1`，件 `.scratch/wisp/probes/293/a1/` 六枚（`00` 7,403／`10` 8,931／`20` 10,450／`30` 12,661／`40` 13,038／`90` 9,253 字节），五笔 commit `97411044`→`7285ef83`→`7c001a84`→`c70e70ea`→`e371f1b8`；★**五笔我逐枚 `git log -1` 现跑核过存在、六枚件逐枚 `ls` 对上字节、非 `.md` 落件＝0 枚**）

- **翻两格：`AC#0`／`AC#1` ⇒ `[x]`**（这两格要的是"读数齐不齐＋落点谁跟着谁"，属我裁；`AC#2`..`AC#5` 属落地腿与非实现者，⛔ 一枚不碰）。它⛔ 零产码、⛔ 零翻框、⛔ 没跑 `go build`/`vet`/`test`（按派单让给在飞的写腿），并把**没跑的尺具名归口**（`90-unrun-rulers.md`）——★其中一条我要转述给落地腿：**它一枚"改前基线"都没取，所以 `293-r1` ⛔ 不许引"今天 `cmd/wisp` 全绿"这类旧读数当基线**。
- ★**承重读数我逐把自己复跑（三把对上、一把纠号）**：`git grep -n "SetTrayChecks" HEAD -- internal cmd` ⇒ **只命中注释 `ball_windows.go:952` ＋定义 `:953`＝调用者 0 枚**（与票面 16:4x 那把同形、与它一致）；`grep -rn "mutedAtBoot" cmd/wisp/` ⇒ **两命中＝`:101` 声明＋`:256` 一处写，零读**⇒ **"写读比 1:0 的死字段"成立**；`GOFLAGS= go list -deps ./internal/ball` ⇒ rc=0、**145 枚依赖、`internal/*` 仅 5 枚、`internal/audio` 命中 0**，正控同尺 `./cmd/wisp`＝**287 枚、含 `internal/audio` 命中 1** ⇒ 那把全称否定尺**有牙、不是假否定**（这条是本格"必须 `cmd/wisp` 注入"的根据）；`docs/specs/SPEC-08-ui-ball-panel.md:229` 逐字＝`- 托盘：右键菜单（打开面板/静音/暂停唤醒/退出）；左键 = 打开面板。` ⇒ **只许四项存在、从未许勾**。
- ★**它给的那枚"直接冲突"是真的，但它报的行号两处都不对**：`resident_windows.go` 里那段注释逐字为 `// booked), and no mirror of the gate into the tray's own checkmark - the` / `// ball-side trayMuted display flag stays undriven here rather than becoming a` / `// second authority nobody reconciles.`；**内容锚复跑真身在 `:322-324`**（它一篇写 `:293-303`、另一篇写 `:300-302`，两篇两个号）。⇒ 本格按**内容锚**裁（认 `no mirror of the gate into the tray's own checkmark` 这三行），⛔ 不认行号；这条同时印证台账 `A789` 那条定式：**行号锚会腐烂，三方三号是常态**。
- **裁：`AC#1` 走甲形（推）**——`cmd/wisp` 侧注入、**`internal/ball` 0 改动、新包级依赖边 0、新协程 0**，形状照现成的后置 setter（`resident_windows.go` 里 `rb.attachMuteGate(...)` 那一族：定义在 `resident_ball_windows.go`、带自己的锁、读回门态）。⛔ 乙形（把状态位搬进 `internal/ball` 的 tray 结构体）不落：菜单是**每次右键现建**（`ball_windows.go` 的 `sel := showMenu(...)` 那一行）、且乙要扩 `internal/ball` 导出面，而 `SPEC-01` 自陈该包无业务逻辑 ⇒ 射程大于本票欠账。
- ★**落地腿必须同批做、票面没写的一件**：上面那段"**deliberately does NOT mirror**"的注释，它的前提**在 `290-r1` 之后已经变了**——那时托盘那一项点了不动设备，勾自然无事可镜像；今天它**真拧门**（`OnTrayMute → muteGesture → gate.SetMuted`，本次真机两形已证）。⇒ `293-r1` **必须把这段注释改写成新理由**，理由要写成主从关系：`gate.Muted()` 是**唯一真相源**、球侧那枚 bool 只是**投影**（⛔ 不是第二座权威）。⛔ 不许旧句与新代码并存——留着的后果是下一位读到"这是故意的"就又把它关掉。
- **三条禁区我在此钉死（写进派单）**：① ⛔ **不许用 `ra.mutedAtBoot` 当真相源**（死字段、`toggleMute` 不更新它 ⇒ 接上去＝把"永不动"换成"永远错"）；② ⛔ **`SetTrayChecks` 的第二实参（暂停唤醒）这发不许顺手接**——它今天**无真相源可喂**（`internal/speech` 只有 `doc.go`、`[voice] enabled` 出厂 `false`、而那一项今天只 `recordBallGesture`），交付＝**显式传 `false` 并在注释里写明为什么**（⚠ 这条也解释了 ⓐ 的不对称：**静音不打勾＝谎，暂停唤醒不打勾＝今天无辜**，两枚不能一并"修"）；③ ⛔ 新代码必须先过 `rb.b` 的 nil 护栏（那条早退分支在 `attachMuteGate` 之前）。
- ★**`AC#0 ⓑ` 那问"契约写没写"＝没写（契约空白），但我裁：⛔ 不摆机主**。理由照第 14 款射程：这一格**零新增能力**、只是把已经发生了的事（门开了）显示成它本来的样子、完全可逆、有撤销口令（「**293 撤镜像**」＝摘掉注入那一跳即回到今天）；而它**卡的另一面**（球变灰球＋斜杠）早就是 `Q-81`，机主已裁"先不要"，那需要人工批准才能动 D43 表 ⇒ **本格不碰状态机、不碰 `Q-81`**，只把"契约从没许诺托盘勾号"这一事实具名记进台账（`A790`），⛔ 不改 `PLAN.md`／`SPEC-08` 一字。
- **排程**：`293-r1`（甲形＋同批改注释＋三条禁区）按在 **`296-r1` 交回并裁完之后**（同写 `cmd/wisp`，⛔ 同包两枚写腿不并发）⇒ `293-v1` 非实现者裁 `AC#2`..`AC#4`；`AC#5`（肉眼那枚勾）**归下一次合并真机窗口**，与票 296 `AC#3`、票 247/290 欠的"托盘退出补 `levels_delivered=`"、"设备被占"自造形**同一次**跑完，⛔ 不占机主多次。⛔ 零 push。

**Status:** **`AC#0`／`AC#1` 已由只读腿交付并经编排者裁完（20:2x，见上节：甲形、同批改写那段"故意不镜像"的注释、三条禁区）；`AC#2`..`AC#5` 一枚未翻，等 `293-r1` 与非实现者 `293-v1`；真机那形归合并窗口。**⛔ 零 push。
⚠ 原 Status 逐字留档（已过期、不删）：「**未开工**（本票 16:5x 立，六格全未勾；⛔ 未派腿）。」——"未开工"对 `AC#0`/`AC#1` 已作废，对 `AC#2`..`AC#5` 仍成立。

## Progress log

- [2026-10-10 10:36:55 +08] agent=293-r1 did=交 AC#2/AC#3/AC#4（甲形，⛔ 翻框）：产码四枚文件全在 cmd/wisp（26289b9e＋9dade451），internal/ball 逐字节零改动（尺＝git diff 6ef14788..HEAD -- internal/ball | wc -c ⇒ 0）、新包级依赖边 0（go list -deps ./internal/ball＝145 枚、internal/audio 命中 0）、新协程 0（diff 里 grep -c 'go func(' ⇒ 0）；形状＝读半 residentAudio.trayMuteState 回 gate.Muted()（⛔ 未接 mutedAtBoot，那枚仍是写读 1:0）＋写半 rb.b.SetTrayChecks 作方法值挂在 residentBall 的 trayMuteRead/trayCheckPush（共用现成 muteMux，⛔ 新锁）＋mirrorTrayMute 两跳触发（muteGesture 的 executed 出口＝起点①点那一项与②热键共用的一条尾巴；resident_windows.go 挂载点下一跳＝起点③出厂静音位），第二实参逐字 push(muted,false) 并把两枚不对称的理由写进注释，nil 护栏 if rb.b != nil 在取方法值之前、票 290 那条早退分支位置一字未动；那段"故意不镜像"的注释（内容锚 no mirror of the gate into the tray's own checkmark，本腿量到 :333-335）已改写为主从实话且旧句 0 命中（grep ⇒ 0），挂载行 rb.attachMuteGate(raudio.toggleMute) 逐字未动，muteGesture 文档同批补一段防下一位当既有政策；AC#2ⓐ＝7 枚用例（三枚起点各一枚＋拒交接读真相源＋无门不写＋nil 安全＋写次对齐），定向尺 rc=0 全绿；AC#2ⓑ＝两枚突变体各有牙（A 恒假→红 4 枚含 boot、B 摘掉手势那跳→红 4 枚含拒交接，两枚互补短板都具名写了杀不到哪几枚），红句逐字抄进 10 件，还原证明＝sha256 前后同值 74a52c42…9c32 ＋ git status 复原为空（⛔ 拿"调用者枚数＝1"当凭据；⛔ 放宽任何既有断言）；AC#4＝GOFLAGS= go build ./... rc=0、d22scan rc=0（clean，cmd 产码 39 枚／ban#8 cmd 119 枚）、整包同尺改前 2 发改后 2 发逐名作差：改前并集 7 枚红（★改前⛔ 不是绿的，我没引任何旧读数当基线）、改后两发红名交集 5 枚全 ⊂ 并集 ⇒ 新增红 0 枚，我自造的那一枚真红（TestTicket255RosterEvidenceLines…，根因＝票 255 名册按行号引产码行、我的插行使 :256/:257/:316 漂移）已当场修回并全程留痕：trayMuteState 改挂文件尾使 5 行引用复位、struct 注释 11→6 行、名册那一行逐字改指 :322（token 未变、断言未放宽）——这一枚越界改动请 293-v1 单独判；改后两发 PASS 计数 507/507（含子测试那把尺，顶层 389），我的 7 枚在 #3/#4 各全绿；logs/post-run-1、post-run-2 两发属"跑动中改了树"的中间态、⛔ 计入门禁，已在 20 件具名作废 next=请编排者派非实现者 293-v1 裁 AC#2..AC#4；三件待人裁或待跑：① 票面「现量」那条"SetTrayChecks 调用者 改后仍 0"与 AC#2 自相矛盾（勾要跟门就必然多出写侧，我按 AC#2 裁形并把矛盾记进 30 件）② gofumpt 那把尺今天不可能为空（HEAD blob 口径 2 枚真未格式化＝票 35 面板传输两枚测试、工作树再多一枚 models.go 是 CRLF 签出形态；我动的 5 枚为空 rc=0，⛔ 顺手格式化别人的文件）③ AC#5 那枚肉眼勾归合并真机窗口，本腿⛔ 声称用户看得见；另 AC#1 那格"internal/ball 依赖名册"我复跑同号（145/0）
（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
- [2026-10-09 20:42:36 +08] agent=编排者 did=收 293-a1 并裁完 AC#0/AC#1（翻两格，AC#2..AC#5 一枚不碰）：五笔 commit 与六枚 .md 件逐枚 git log -1／ls 现核、非 .md 落件=0；承重读数我自跑对上三条（SetTrayChecks 调用者 0＝注释 :952＋定义 :953、mutedAtBoot 写读比 1:0＝:101/:256 零读、go list -deps ./internal/ball=145/5/0 且正控 cmd/wisp=287/含 audio 1⇒尺有牙）、SPEC-08:229 逐字只许四项存在从未许勾；★它报的那枚"注释已决定不镜像"属实但**两篇两个行号都不对**，内容锚现量真身 :322-324 ⇒ 落地腿必须同批把该三行改写成主从理由（gate.Muted() 唯一真相源、球侧 bool 是投影），⛔ 旧句不许与新代码并存；裁甲形（cmd/wisp 注入、internal/ball 0 改动、新包级边 0、新协程 0）；三条禁区钉死＝⛔ 不用 mutedAtBoot 死字段、⛔ 不接第二实参 pause-wake（显式传 false＋注释）、⛔ 先过 rb.b nil 护栏；AC#0 ⓑ 契约空白＝**我裁不摆机主**（零新增能力、可逆、口令「293 撤镜像」；被许诺的那面卡在 D43 无 Sleeping→Muted＝已裁"先不要"的 Q-81，本票不碰）；另把它申报的 AGENTS.md D 表行号问题用更严的尺重扫＝**15 枚过期**（D32/33/34/45 偏+2、D36-D44/46/47 偏+6）已改正并把尺写进该文件（PLAN.md 一字未动），腿只报到 2 枚、根因＝宽松包含判据 ⇒ 定式"比对那一行是不是这个标题、不是有没有这个号" next=293-r1 按在 296-r1 交回并裁完之后（同 cmd/wisp 写面），AC#5 肉眼那形并入合并真机窗口
