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

- [ ] **AC#0 只读普查（⛔ 不许只答"没找到"）**：把"生产形状"这一族今天**谁在验**列成名册——`grep -rn "startResidentBall" --include=*_test.go cmd/` 命中逐枚判"起的是真窗还是手搓的 `&residentBall{}`"；并答一问：**除 `attachMuteGate` 外，装配根里还有几枚后置 setter/注入挂钩同样是"直调有、走查零"**（本票只治静音这一条，但欠账名册要一次列全，⛔ 不要每波新量一枚）。
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

**Status:** **未开工**（本票 16:5x 立，六格全未勾；⛔ 未派腿）。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
