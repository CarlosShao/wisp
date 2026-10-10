# 票 293 · 写腿 `293-r1` · `20` `AC#4` 门禁四把尺（改前基线自取／改后两发／逐名作差）

本腿 `date` 起手＝`2026-10-10 09:43:32 +0800`；HEAD 起点＝`6ef14788`；产码笔＝`26289b9e`＋`9dade451`。
每把尺都带**逐字命令**与**自己的 rc**；原始 stdout 全在 `logs/`（非 0 字节，见 §6）。

## 0. 进程护栏（⛔ 有活体就别读整包）

`tasklist | grep -ci 'wisp.exe\|balldebug.exe'` ⇒ **0**（起手 09:43 一发、`sh scripts/d22scan.sh` 之前 10:1x 再一发，都是 0）。
跑动期间出现过的只有本腿自己起的 `wisp.test.exe`（`go test` 的测试二进制），⛔ 没有杀过任何不属于我的进程。

## 1. `GOFLAGS= go build ./...` ⇒ **rc=0**

- 改后第一发（重排 `trayMuteState` 之后）：`build_rc=0`（`logs/targeted-after-relocate.txt` 同批）。
- ⚠ 票 290 那一批的 `models.go` CRLF 见 §3 的 gofumpt 段，与本尺无关。

## 2. 整包名册：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ ./internal/audio/ ./internal/ball/ -count=1 -v`

尺＝**`grep -c -- '--- PASS'`（含子测试，本仓沿用的那把）**，另附 `grep -c -- '^--- PASS'`（顶层）与 `grep -c -- '--- FAIL'` 三个数。
⛔ 单发名册不当可靠尺：改前两发、改后两发，逐名作差再取交集。

| 发 | 落件 | rc | `--- PASS` | `^--- PASS` | `--- FAIL` |
|---|---|---|---|---|---|
| 改前 #1（动产码**之前**自取） | `logs/pre-baseline.txt` | rc=1 | **499** | 381 | 7 |
| 改前 #2 | `logs/pre-baseline-2.txt` | rc=1 | **501** | 383 | 5 |
| 改后 #3 | `logs/post-run-3.txt` | rc=1 | **507** | 389 | 6 |
| 改后 #4 | `logs/post-run-4.txt` | rc=1 | **507** | 389 | 6 |

★**改前基线不是绿的**（票面裁定段早就声明不许引旧读数）：改前两发的红名并集＝**7 枚**，逐字在
`logs/pre-fail-union.txt`，全在**面板／WebView2 与 tokens.css** 那一族，与本票射程无关：

```
--- FAIL: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe        （页侧 15s 无回报）
--- FAIL: TestAC14AwaitedBindingReplyReachesThePage
--- FAIL: TestAC14GoSideEvalPushReachesThePage
--- FAIL: TestAC4FocusReturnToPriorWindowGap33r5                    （改前#1 红、#2 绿＝前台时序抖动）
--- FAIL: TestAC4PriorFocusSurvivesARefusedPanelSample              （改前#1 红、#2 绿＝同上）
--- FAIL: TestPanelHostRealWindowHopAndLifecycle                    （cold bring-up -1.000 ms＝本机拉不起 WebView2）
--- FAIL: TestC21TableColourRowsMatchTokensCSS                      （design/assets/tokens.css 在这棵树里不存在）
```

**逐名作差**：
- 改后 #3 的红 6 枚 ＝ 上面并集的子集（**新增红 0 枚**）。
- 改后 #4 的红 6 枚 ＝ 并集子集 ∪ `Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands` ＼ `TestAC4Focus…`。
- **两发红名的交集＝5 枚，全部 ⊂ 改前并集 ⇒ 按票面那条取交集的尺，新增红 0 枚。**
- #4 里那枚**只出现一次**的红具名申报：`Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands`
  是 11.5s 级、要拉真 `wisp run` 子进程的形状，整包并发下最容易抖；本腿把它**单发重跑 3 次全绿**
  （`logs/flake-check-run4-name.txt`，三发 rc 均＝0，逐字 `--- PASS: Test197… (11.56s)/(11.31s)/(11.34s)`），
  且改后 #3 全包那一发里它也是绿的 ⇒ 判定＝**整包负载抖动，非本腿引入**；⛔ 我据此放宽任何断言，也没重跑掉那把尺。
- 新增的绿：本腿 7 枚 `TestAC293*` 在 #3 与 #4 各全绿（`grep -E '^--- PASS: TestAC293' | wc -l` ⇒ **7 / 7**）。
- `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` 在 #3／#4 均 **PASS**（它一度被我的插行打红，见 §5）。

## 3. `$(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp internal/ball` ⇒ 与票面预期不符，具名报回

裸 `gofumpt` 不在 PATH＝`rc=127`（本腿⛔ 用裸名，一律走 `"$(go env GOPATH)/bin/gofumpt.exe"`；下面读数 rc 都是 0）。

- **工作树口径**（尺＝`gofumpt.exe -l cmd/wisp internal/ball`）⇒ **非空 3 枚**：
  `cmd/wisp/models.go`、`cmd/wisp/panel_inbound_guards_35r3_test.go`、`cmd/wisp/panel_transport_35r2_test.go`
  （原样 stdout＝`logs/gofumpt-post.txt`）。
- **HEAD blob 口径**（⛔ 落在仓外：`git show HEAD:<path>` 逐枚抽到 `$TEMP/293_headblobs/`，149 枚 `.go`）⇒ **非空 2 枚**：
  后两枚（`panel_inbound_guards_35r3_test.go`、`panel_transport_35r2_test.go`），`models.go` **不在其中**。
- ⇒ 两条结论：① `models.go` 那一枚是**检出行尾差异**（工作树 CRLF vs blob LF：`file` 报 `with CRLF line terminators`、
  `cmp` 与 blob 不等而 `git status` 干净）＝**签出形态**，不是谁写的格式债；② 另两枚是 **HEAD 就存在的真未格式化文件**，
  都属票 35 那一批面板传输测试，⛔ 本腿没动、也⛔ 顺手格式化（那会挪别人的分母）。
- **本腿名册（射程＝我动过的 5 枚 `.go`）**：尺＝
  `gofumpt.exe -l cmd/wisp/config_readers_255.go cmd/wisp/resident_ball_windows.go cmd/wisp/resident_audio_windows.go cmd/wisp/resident_windows.go cmd/wisp/resident_tray_mute_293_windows_test.go`
  ⇒ **空，rc=0**（10:1x 现跑）。
- ⚠ 与派单「⇒ 空」的冲突＝**票面／转述过期**：整目录那把尺在今天不可能为空，除非动不属于本票的文件。本腿按
  「我动的文件为空 ＋ 与 HEAD 逐名对齐」交，两条读数都在上面。

## 4. `sh scripts/d22scan.sh`（纯净树）⇒ **rc=0**

`logs/d22scan-post.txt`（41,989 字节）末行逐字：
```
d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=229, bans #1-5 cmd/=39, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=524 Go files, comments and _test.go included, ban #8 cmd/=119 Go files, comments and _test.go included
```
⇒ 裸 `go func(`（ban #1）与零 emoji 那把（ban #8，含注释与 `_test.go`）**同时过了我这两笔**。

## 5. ★中途那一发真红（我自己引入、当场修回，全程留痕）

改后 #1（`logs/post-run-1.txt`，rc=1、7 枚红）里出现一枚**改前没有**的红：
`TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`。根因⛔ 与被测行为无关，而是**票 255 的名册按行号引用产码行**：

- 我把 `trayMuteState` 插在 `resident_audio_windows.go` 中部 ⇒ 它下面 18 行整体下移 ⇒ 名册里
  `resident_audio_windows.go:256 [c.Audio.MicMutedDefault]`／`:257 [c.Voice.Enabled]` 那 5 行引用全部漂走。
- 我在 `residentBall` 的 struct 里加了字段 ⇒ `resident_ball_windows.go:316 [Hotkeys:  cfg,]` 漂到 `:322`。

修法（`9dade451`）＝**① 把新函数改挂到文件尾**（那 5 行的号复位、⛔ 动名册），
**② struct 注释从 11 行压到 6 行**（Go 的字段必须在 struct 体内，`Hotkeys: cfg,` 在字段之下 ⇒ 任何插行都会让它挪，
这一枚躲不掉），**③ 把名册那一行引用逐字改指 `:322`**（token `[Hotkeys:  cfg,]` 一字未动、断言形状未松、
该行现在仍必须逐字携带那个 token ⇒ ⛔ 这是放宽，是把引用重新对准仍在原地的读者）。
`config_readers_255.go` 只在 `cmd/wisp/**` 名册内、⛔ 不在十枚禁列里；这一枚越界风险已具名交给 `293-v1` 判。
复位后：定向 `go test ./cmd/wisp/ -count=1 -run 'TestAC293|TestTicket255|TestAC290' -v` ⇒ **rc=0、零 FAIL**
（`logs/targeted-after-repoint.txt`），随后 #3／#4 两发全包也都 PASS 那枚尺。
⇒ 本腿**没有**用「改前也红」遮掩这一发，也**没有**把它算进环境抖动：它是我的改动造成的、已修回并有前后读数。

## 6. 那两发**作废**的改后读数（⛔ 拿来当门禁数）

`logs/post-run-1.txt`（rc=1，506 `--- PASS`，7 枚红）与 `logs/post-run-2.txt`（rc=1，506，7 枚红）是
**在跑动中我把树改了**（10:05／10:12 起跳，10:15-10:18 才做复位那一发）⇒ 它们的二进制属于中间态，
两发都带着 §5 那枚 `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` 红（其余 6 枚＝改前并集里的老面孔），
同时 `^--- PASS: TestAC293` 在 #1／#2 里也都是 7 枚。
**门禁那四把尺的数＝§2 表里的 #3 与 #4 两发**；#1／#2 只作 §5 的过程留痕，⛔ 计入门禁。

## 7. 终态读数（`2026-10-10 10:36:55 +0800` 现跑，HEAD 已是编排者的 `cfacfd3c`）

```
GOFLAGS= go build ./...                                   ⇒ rc=0
gofumpt.exe -l <本腿动过的 5 枚 .go>                        ⇒ 空，rc=0
tasklist | grep -ci 'wisp.exe\|balldebug.exe'              ⇒ 0
```

## 8. 落件清单（⛔ 无 `.out`；`find .scratch/wisp/probes/293/r1 -type f -size 0` ⇒ **零枚**）

件＝`00-anchor.md`（射程＋四把内容锚）·`10-ac2-behaviour-and-mutation.md`（ⓐ＋ⓑ＋还原证明）·
`20-ac4-gates.md`（本文件）·`30-ac3-scope.md`（越界名册＋那枚票面自相矛盾）。
`logs/` 逐枚字节（10:36 现跑）：`pre-baseline.txt` 347,356／`pre-baseline-2.txt` 347,148／
`pre-1-roster.txt` 28,816／`pre-fail-union.txt` 360／`post-run-1.txt` 363,532／`post-run-2.txt` 363,533／
`post-run-3.txt` 361,559／`post-run-4.txt` 371,739／`mutant-a-const-false.txt`／`mutant-b-no-gesture-mirror.txt`／
`targeted-after-relocate.txt` 41,989／`targeted-after-repoint.txt` 41,679／`flake-check-run4-name.txt`／
`d22scan-post.txt` 41,989／`gofumpt-post.txt`。四枚 `.rc` 件各 5 字节（`rc=1`）。
⚠ `logs/` 是派单 §9 指定的落点（大 stdout 先落件再取小窗），形态＝`.txt`＋`.rc`；
两把仪器的分母都⛔ 被它挪动（`gofumpt` 只看 `.go`；`d22scan` 的 emoji／ban 扫描按 `.go`／白名单扩展名走，
且 `tools/d22scan/main.go` 对 `.txt` 只在其被列入口径时才读），§3／§4 那两发的计数即证。

