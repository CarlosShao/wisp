# 票 294 — "拧门那只手是装配根挂上去的"这句话**今天只由 grep 尺与一句注释保证**：`runResident` 从未被任何测试调用，球侧两支闭包也没被走

**立票**：2026-10-09 16:5x 编排者（来路＝票 290 非实现者验收腿 `290-v1` 的 N3／M2 两问，件 `.scratch/wisp/probes/290/v1/10-mutations-and-cites.md:48-105`；承重读数我 16:4x 自己复跑过，并**更正了腿那把尺的字面读数**，见下面现量第二条）
**性质**：★**不是"仪器里没有造出那个世界"那一类，是"仪器在仓里躺着、没被指向这一跳"**。票 290 的六枚用例判的是"门被拧了没"（读的是 `gate.Muted()`/`gate.Open()`，M1 变异证明它们有牙）；而**"谁把手挂到门上"**这一跳——装配根那一次 `attachMuteGate` 调用、以及球侧 `Events` 里那两支闭包真被路由到 `muteGesture`——**在生产形状上今天没人验过**。
⚠ 腿一度把这条写重了（写成"仓里没有能造生产形状的仪器"），现跑后它自己更正：仓里**有两件现成仪器**，290 只是没用（本票的台件写法因此⛔ 不需要新造 seam，见 AC#1/AC#2）。

## 现量（每条都带尺；⛔ 引用前先重跑，行号与枚数都是快照）

- **生产挂载点只有一处**：`cmd/wisp/resident_windows.go:303` 逐字 `rb.attachMuteGate(raudio.toggleMute)`（起因是同一函数体里 `:217` 建球**早于** `:278` 建门 ⇒ 只能是后置 setter，这条**时序**就是整条链成立的前提）。
- ★**`runResident` 在任何测试里都没被调用过**（尺＝`grep -rn "runResident(" --include=*_test.go cmd/ | grep -v '// '`＝**0 行**，我 16:4x 现跑）。⚠ **腿原文那把尺的字面读数是假的、实质结论我复跑后收下**：它写 `grep -rn "runResident\b" --include=*_test.go cmd/` → "零命中，rc=1"，而该尺现跑＝**27～29 命中**（全是注释、AST 里 `fn.Name.Name == "runResident"` 那类字符串比较、以及 `t.Fatalf` 文案）⇒ ⛔ 本票引用时以**调用形状**那把为准。"为什么只能走查不能直调"仓里早写了：`cmd/wisp/resident_ball_228_test.go:17` 逐字 `// Why a source walk and not a runtime call: runResident() never returns (it`。
- **台件是手抄了那句话、没走装配根**：`cmd/wisp/resident_mute_290_windows_test.go:93` 与 `:206` 逐字各是 `rb.attachMuteGate(ra.toggleMute)`（`rb := &residentBall{}` 手搓的球），`:235` 是 `nilBall.attachMuteGate(func() (string, bool) { return "", false })` ⇒ 尺＝`grep -rn "attachMuteGate" --include=*_test.go cmd/wisp/`＝**3 命中、全直调、零走查**。
- **球侧两支闭包没被走**：台件直接调 `rb.muteGesture("mute-hotkey")`（`resident_mute_290_windows_test.go:117`），跳过 `resident_ball_windows.go:321` `OnMuteHotkey:    func() { rb.muteGesture("mute-hotkey") },` 与 `:325` `OnTrayMute:      func() { rb.muteGesture("tray-mute") },`；`internal/ball/ball_windows.go:660` `case hkMute:` 与 `:678` `case menuMute:` 更要**真窗口＋真消息泵**才走得到，台件里一枚都没走。
- ★**两件现成仪器（本票⛔ 不许新造 seam，直接用）**：
  ① 真窗口形状＝`cmd/wisp/resident_hotkey_258_windows_test.go:266` `Test258BridgeRebindsLiveKeysFromConfigEdit`，`:278` 逐字 `rb := startResidentBall(observe.NewRegistry(), nil, src, src)` 起**真球窗**，`:289-291` 读 `rb.b.HotkeyReport()` 的注册表（`:290` 逐字 `if !boot.IsLive(hkSummon258) {`），且它的 `writeConfig258` 段里 `:273` **逐字已经写了** `mute = "Ctrl+Alt+M"` ⇒ "mute 那枚热键到底注册成功了没"这一跳**今天有现成的尺可读，只是没人为 `hkMute` 读它**。
  ② 装配根形状＝对 `runResident` 函数体做 **AST 走查**的先例三枚：`cmd/wisp/resident_hotkey_258_test.go:25`（`func runResidentBody258(t *testing.T) (*ast.BlockStmt, string) {`，`:33` 抽 `runResident` 的体，`:45-68` 逐条断言）、`resident_ball_228_test.go:166-247`、`resident_approval_risk_256_windows_test.go:527-569`。
- ★**M2 实测出的另一半缺口**：把 `resident_audio_windows.go:172`／`:174` 那两枚 `true` 改成 `false`（`:162` 拧门那行保持原样）⇒ **六枚用例全绿**（腿读数 `rc_m2=0`）。⇒ happy path 上**没有一枚用例断言 `executed==true`**；`TestAC290NoGateSaysWhichShapeThisProcessIsIn`（`resident_mute_290_windows_test.go:196-199`）只在**没有门**那一支断言它为假。后果：`muteGesture`（`resident_ball_windows.go:527-537`）拿到 `executed=false` 时改打 `:531` 逐字 `slog.Warn("mute gesture found no gate to turn", "gesture", name, "why", outcome)`，而 `:536` 那句用户可见输出**仍是门读回的真话** ⇒ **严重性＝低（可观测性面，不是用户可见错状态）**，但这是那枚函数值契约**没被钉的另一半**。
- **为什么 M2 并进本票、不另立一枚（理由具名，⛔ 不是省事）**：M2 的修法就是"给那一跳做一枚生产形状用例"里**必须断言的第一件事**（`executed==true`），拆成两票会让第二票重复造同一件夹具。⇒ 本票 AC#3 就是那一格，判据独立成框。

## 要建什么（⛔ 先普查再落地；三格各自成框）

- [x] **AC#0 只读普查（⛔ 不许只答"没找到"）**：把"生产形状"这一族今天**谁在验**列成名册——`grep -rn "startResidentBall" --include=*_test.go cmd/` 命中逐枚判"起的是真窗还是手搓的 `&residentBall{}`"；并答一问：**除 `attachMuteGate` 外，装配根里还有几枚后置 setter/注入挂钩同样是"直调有、走查零"**（本票只治静音这一条，但欠账名册要一次列全，⛔ 不要每波新量一枚）。
- [ ] **AC#1 装配根那一跳的形钉（读 AST、不需要跑 `runResident`）**：照 `resident_hotkey_258_test.go:25-68` 的走查形状加一枚用例，断言**三件事**：ⓐ `runResident` 体内 `callExpr` 里 `id.Name == "attachMuteGate"` 的次数 **== 1**（不是 0、不是 2）；ⓑ 它的接收者是 `startResidentBall` 那枚返回值 `rb`、实参点名是 `raudio.toggleMute`（**不接受"任意函数值"**）；ⓒ 它在体内的位置**晚于** `startResidentAudio` 那一次调用（时序是这条链成立的前提：`:217` 建球 < `:278` 建门 < `:303` 挂载）。⚠ 硬约束：走查尺的锚要落在**形状**上（`callExpr`/selector 链），⛔ 不许写成"某行必须存在"的行号锚（本仓实测：夹具一插就把行号锚打死）。
- [ ] **AC#2 球侧注册那一跳的形钉**：复用 `resident_hotkey_258_windows_test.go:266-294` 的真球窗形状（含它那句 `SKIP-LOUD` 的 headless 退路），在 `:289` 那把 `HotkeyReport` 尺上**对 `mute` 那枚 id 现读 `IsLive`**，红／跳都具名。⚠ 本格若因 headless 而 skip ⇒ **skip 不算读数**，具名归"欠机主在场那一形"（同票 290 `AC#3` 的处置），⛔ 不许把 skip 读成过。
- [ ] **AC#3 `executed` 的另一半（M2 那一格）**：happy path 上必须有用例断言 **`executed == true`**（门真被拧了那一条路）。判据带反形正控：把 `:172`/`:174` 那两枚 `true` 改 `false` ⇒ 指名用例必须红（今天它绿，`rc_m2=0`＝本格的改前读数）。⛔ **不许改成"读返回值里的字样"**——票 290 的 M1 已证：读返回值会被"句子里写了成功"骗过去；判据锚只能是门自己的读数（`gate.Muted()`/`gate.Open()`）＋这枚 `executed`。
- [ ] **AC#4 越界检查**：`git diff-tree -r --name-only --no-renames` 逐笔过十枚禁列（`frontend/**`／`design/**`／`PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／三枚冻结件 `internal/panel/tokens_fourway_test.go`·`internal/panel/l2_grant_boundary_test.go`·`internal/perm/ticket90_persist_test.go`）任一路径出现 ⇒ 直接退回；⛔ 不许为此放宽任何既有断言（含 `loop_approval_test.go:213-214` 那枚唯一等值断言、票 255 那张 roster 表）。
- [ ] **AC#5 门禁四数**：`GOFLAGS= go build ./...` rc=0；`$(go env GOPATH)/bin/gofumpt.exe -l <自己动过的目录>` 空（裸 `gofumpt` 不在 PATH＝`rc=127`，具名不掩盖）；`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1 -v` 改前改后**各 ≥2 发取交集**并逐名作差（⛔ 单发名册不是可靠尺；三数必须带计数命令逐字，本仓沿用 `grep -c -- '--- PASS'`（含子测试））；`sh scripts/d22scan.sh` 纯净树 rc=0；每枚门禁件自落一行 `rc=N`（⛔ 0 字节件＝那格没交）；跑读数前 `tasklist` 现跑 `wisp.exe`／`balldebug.exe`＝0。

## 禁区

- ⛔ **不新开 seam、不加新包级依赖边**（两件现成仪器够用；`internal/ball` 的依赖里今天不含 `internal/audio`，那是票 247 `AC#0` 当年禁掉的方向）。
- ⛔ 不许为了"走到 `case hkMute:`"去 `PostMessage` 真热键或 press 键盘（AC#2 用 `HotkeyReport` 现读，AC#1 用 AST）；真按键那一形若要做＝归机主在场的真机窗口，⛔ 不在本票。
- ⛔ 不动默认值、不碰 D43 表、不改球侧那两支闭包的语义（本票只给它们**加仪器**）。
- ⛔ 不许顺手把票 290 `AC#3` 或票 293/295 的格勾掉。
- git：只 commit 不 push；commit 必带**显式 pathspec**（`git mv` 的改名笔旧名与新名两个路径都要给）；禁 `add -A`／`amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内不删文件、不建 worktree。

## 排程与串行

- 写面＝`cmd/wisp` 的测试面（＋可能 `cmd/wisp` 产码那一枚 `executed`）⇒ ⛔ 与票 293 串行、⛔ 不与任何"落地腿"混批；与票 295（同写 `cmd/wisp/resident_audio_windows.go`）**不同批**。
- 编排者队列：重 build `build/wisp.exe` → 票 292 → 机主真机窗口 → 票 293 `a1` → **本票** → 票 295。

## 编排者裁定（2026-10-09 21:1x，来路＝只读普查腿 `294-a1`，件 `.scratch/wisp/probes/294/a1/`（`ac0-census.md` 112 行／`ac1-ac3-shapes.md` 82 行＋两枚 commit-msg 件），两笔 commit `b6a41f69`→`23d06056`；⛔ 上面各节原句一字不改，本节追加）

- **合规**：三枚落件**全是 `.md`**（非 `.md`＝0）、票面只 `+2` 行（Progress log）、⛔ 零源码改动、⛔ 零 go 命令（`296-r1` 正独占 Go 编译面）、⛔ 未动 `docs/**`、⛔ 没给自己翻框。⇒ 派单模板"补位只读腿只新建 `.md`＋禁跑 go"这一条**第二次照跑成功**。
- **★我复跑用的是 HEAD 对象层那把尺（`git grep -n … HEAD`），⛔ 不是工作树**：写腿在取"改前名册"那几发时会把同一枚 `resident_windows.go` **在工作树上临时还原**（腿自己撞上：两次 grep 之间同一行从 `:325` 漂到 `:336`；我 20:5x 量到还原态、21:1x 复量 `git diff --numstat HEAD -- cmd/wisp/resident_windows.go` 已空＝它把树还原回去了）。⇒ **结论：树会呼吸，写腿在飞时"盘上现量"自带歧义 ⇒ 锚点一律引 HEAD**。复跑结果：`attachMuteGate` 产码 **1**（`cmd/wisp/resident_windows.go:336` 逐字 `rb.attachMuteGate(raudio.toggleMute)`）／测试直调 **3**（`resident_mute_290_windows_test.go:93`/`:206`/`:235`，全挂手搓 `&residentBall{}`／`nilBall`）／定义 **1**（`resident_ball_windows.go:472`）／**AST 走查 0** ⇒ 腿那句"生产挂载有人调、没人验形状"**成立**。手搓 `&residentBall{` **13** 枚复现一致。
- ★**一枚口径差要钉死（不是谁错）**：`startResidentBall(` 在 `cmd/wisp` 测试里整族命中 **12 行**，其中 **2 行是注释**（`resident_hotkey_296_windows_test.go:19`/`:26` 的说明句），真调用 **10**。腿报 10＝**剥注释那一把**，我这把 12＝**整族那一把**。⇒ 今后名册尺**必须写含不含注释行**（票 292 结的"词面尺带不带冒号"是同一条规矩，本票再兑现一次）。手搓 `&residentBall{` **13** 枚复现一致。
- **`AC#1` 落点裁（本格不翻，缺"跑出来的红/绿"）**：新文件 `cmd/wisp/resident_mute_294_test.go`；模板＝`resident_hotkey_258_test.go:25-68` 那一族（`token.NewFileSet`＋`parser.ParseFile(…, parser.ParseComments)`＋`ast.Inspect`＋`runResident` 的 `*ast.BlockStmt`，我逐字复读到）。⚠ **一处纯照抄会造恒绿假钉**：258 那把尺匹配 `call.Fun.(*ast.Ident)`，而生产那一行是**方法调用** `rb.attachMuteGate(…)` ⇒ `call.Fun` 是 `*ast.SelectorExpr`，**必须加一层 SelectorExpr**，否则尺**永远 0 命中**、用例永远"绿"。（先例＝台账 `A710`：我代笔的仪器被判没牙。）
- **`AC#2` 落点裁**：腿的"尺半存"复现——`HotkeyReport.IsLive` 存在（`resident_hotkey_258_windows_test.go:290` 已在读 summon），`internal/ball/hotkey_live_test.go:404` **已对 `hkMute` 读过一次**⇒ 形状有仓内先例；缺的只是 `cmd/wisp` 侧那枚本地镜像常量（`hkMute` 是 `internal/ball/hotkey_windows.go:37` 的非导出常量，`cmd/wisp` 今天只镜像 `1`／`4`，`mute=2` 零镜像）。⚠ **`SKIP-LOUD` 在本机不算交付**：本机 20:5x/21:0x 两枚腿都在真窗上跑出了改前红／改后绿，退路留着是给 headless 的，⛔ 不许拿它当这一格的读数。
- **`AC#3` 裁＋★排程新增一条硬约束**：全仓唯一 `executed` 断言＝`resident_mute_290_windows_test.go:197` 的 `if executed {`（它断的是**假**那一支：无门的腿不许声称拧过门），happy path 的 `executed == true` **没人断** ⇒ 腿"不成立"复现。⚠ 但它的**反形正控**要翻 `resident_audio_windows.go` 里 `toggleMute` 那两支 `…, true`（内容锚 `case gate.Muted():`／`case gate.Open():`），而**票 295 的 `AC#1` 改的正是同一枚函数**（`ra.started = gate.Open()` 就在那几行里）。⇒ **裁：`295-r1` 必须先于 `294-r1`**，两枚 ⛔ 不并发、⛔ 不混批；`294-r1` 的断言与台件都锚在 295 落地之后的形状上。
- **票面现量三枚行号已过期（记我，⛔ 不改原句）**：`:303` 生产挂载 → HEAD 现量 `:336`（`296-r1` 的 `+11` 位移）；`:217` 建球、`:278` 建门同批过期（现约 `:250`／`:311`）。⇒ 腿判"**次序 球 < 门 < 挂载 仍成立**"我复认 ⇒ `AC#1` 的时序断言不受行号腐烂影响，但派单一律**内容锚**（票 298 `AC#3` 那条"锚点口径"同样适用本票）。
- **队列（覆盖上面「编排者队列」那一行，原句留着）**：`296-r1` 交回并裁＋我亲跑 `AC#3` 真机 → `293-r1` → `298-r1` → **`295-a1` → `295-r1`** → **`294-r1` → `294-v1`** → `297-a1`（独占音频设备面）→ 票 287 `AC#2` → 票 288 → `255-r2`。
- **契约面**：本票零触碰（⛔ 不新开 seam、⛔ 不加包级依赖边、⛔ 不改默认值、⛔ 不碰 D43 表）。⛔ 零 push。

**Status:** **`AC#0` 已由只读腿交付并经编排者 HEAD 复跑翻勾（21:1x，见上节）；`AC#1`..`AC#3` 只裁到"落点形状＋缺哪件仪器"，读数（红/绿、真窗 `IsLive(hkMute)`、`executed==true`＋M2 反形）一枚未翻，等 `294-r1` 与非实现者 `294-v1`；`AC#4`/`AC#5` 归落地腿。★排程新增硬约束＝`295-r1` 先于 `294-r1`（同一段 `toggleMute`）。**⛔ 零 push。
⚠ 原 Status 逐字留档（已过期、不删）：「**未开工**（本票 16:5x 立，六格全未勾；⛔ 未派腿）。」——"未开工"对 `AC#0` 已作废，对 `AC#1`..`AC#5` 仍成立。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）

- [2026-10-09 21:06:56 +08] agent=294-a1（只读普查腿，非落地腿） did=AC#0 普查 + AC#1..3 落点判定两件落 `.scratch/wisp/probes/294/a1/`（`ac0-census.md`/`ac1-ac3-shapes.md`），commit `b6a41f69` + 本笔；名册：`startResidentBall` 测试调用 10 枚全真窗、手搓 `&residentBall{}` 13 枚、`attachMuteGate` 3 枚测试直调全挂手搓球；`attachMuteGate` 锚形状＝产码 1（现随 296-r1 并发从 `:325`→`:336`）/测试直调 3/定义 1/注释 1/AST 走查 0；装配根"直调有走查零"欠账除静音外 ≥4 枚（bindBallHost·RegisterShutdownHook(StepCancelTasks)·startResidentTaskSource 调用点·startResidentAudio）；AC#1 模板 `resident_hotkey_258_test.go` 是 `*ast.Ident` 形状，`rb.attachMuteGate` 是 `*ast.SelectorExpr` 方法调用、不能纯照抄须加一层；AC#2 `HotkeyReport.IsLive` 存在且 internal/ball 已对 hkMute 读，但 hkMute 非导出、cmd/wisp 无镜像常量、从未对 mute 读 IsLive；AC#3 无一人断 `executed==true`（`:197` 只断无门⇒false），happy 世界在 290 夹具里已可造真门、缺的只是那半枚断言，且 M2 反形翻 `:172/:174` 撞票 295 同写面须串行 next=待 294-r1 持 Go 编译面现跑：AC#1 走查红/绿＋AC#2 本机真窗 `IsLive(hkMute 镜像)` 实读（SKIP-LOUD 在本机不是退路）＋AC#3 `executed==true`（与 295 串行）；票面 `:303/:217/:278` 行号已漂一律改锚形状（本腿 ⛔ 未自改票面正文）。本把 ⛔ 零 go 命令、⛔ 零翻框、⛔ 未动 docs。

- [2026-10-09 21:37:41 +08] agent=编排者 did=收 294-a1（两笔 b6a41f69→23d06056，三枚件全 .md、非 .md 落件 0、票面只 +2 行 log、零 go 命令、零源码、未动 docs）并裁完 AC#0（翻勾）＋AC#1..#3（只裁落点，读数一枚不翻）：我复跑用 **HEAD 对象层**那把尺（git grep -n … HEAD），⛔ 不碰工作树：写腿取"改前名册"那几发会把 resident_windows.go 在工作树上临时还原（腿两次 grep 同一行 :325→:336；我 20:5x 量到还原态、21:1x 复量 git diff HEAD 该件已空＝树被还原回去）⇒ 定式"树会呼吸，写腿在飞时盘上现量自带歧义，锚点一律引 HEAD"；attachMuteGate 产码 1（resident_windows.go:336）／测试直调 3（290 件 :93/:206/:235）／定义 1（resident_ball_windows.go:472）／AST 走查 0；★口径钉死：startResidentBall 整族 12 行＝真调用 10＋注释 2（腿报剥注释那一把、我那把整族，两把都对差在射程，今后名册尺必写含不含注释行）；手搓 &residentBall{ 13 复现；hkMute 确是 internal/ball/hotkey_windows.go:37 非导出且 cmd/wisp 只镜像 1/4（AC#2 缺器＝本地镜像常量＋一次现读，先例 internal/ball/hotkey_live_test.go:404）；全仓 executed 唯一断言＝290 件 :197 断 FALSE 那一支、happy 的 true 没人断（AC#3 不成立复现）；⚠AC#1 纯照抄 258 模板会造恒绿假钉——模板匹配 *ast.Ident 而生产那行是 SelectorExpr，必须加一层（先例 A710 我代笔的仪器被判没牙）；★排程新增硬约束＝295-r1 先于 294-r1（两枚都动 resident_audio_windows.go 里同一枚 toggleMute：295 改 ra.started、294 的 M2 反形翻那两支 true），⛔ 不并发不混批；票面现量 :303/:217/:278 三枚行号已过期（记我、原句不改），派单一律内容锚 next=按队列 296-r1 交回并裁＋我亲跑 AC#3 真机 → 293-r1 → 298-r1 → 295-a1 → 295-r1 → 294-r1 → 294-v1
