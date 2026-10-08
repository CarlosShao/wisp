# 35-v5 — 非实现者对抗验收判决（票 35 新框 `AC#8`，写腿 `35-r5` 四笔）

锚：`HEAD=d2a4597d58c47a1613e9c93281d53038694a4d7b`（**不是**派单写的 `396c39b6`；那枚差是编排者自己的台账 A694 一笔，`git diff --name-only 396c39b6..HEAD` ＝只有 `docs/reports/pending-and-issues.md` ⇒ 被验树面未变，本判决锚在现量 HEAD）。
时间窗 `2026-10-08 09:44 → 09:58 +08`。⛔ 本腿没改任何产码／测试夹具／`ci.yml`／台账，⛔ 没跑任何 `-tags winlive`，⛔ 没翻任何一枚框（追加前后各量：**8 未勾／2 勾**）。

## 结论一览

| 问 | 判 | 一句 |
|---|---|---|
| Q1 (i) 新用例有没有牙 | **成立** | 我自造的三发新坏法里两发只红那一枚新用例、一发被新用例的**非抛那五支**抓住；把断言削成"不抛就算过"的弱形在静默丢件的坏法下**转绿** ⇒ 那五支是承重的，用例不是装饰 |
| Q2 "唯一红" | **成立** | 红名册恰 1 枚＝`TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsNotCallable` |
| Q3 (ii) 拆前缀拆没拆掉牙 | **成立（判据没变松）** ＋★抓到我自己的题外新面：**枚数尺与行为尺都过，但那枚可复用的 v4 反形驱动被这次行位移打死了 6/7 枚锚** |
| Q4 (iii) 免责声明还在不在射程 | **成立** | 那块"这格不买什么"逐字还在，三处都改成夹具自称，运行时红句里 `NOT a WebView2 behaviour record` 逐字没删；文案更谦虚与判据变松两格分开写、两格都是好消息 |
| Q5 (iv) 〔仅本机可量〕四要素 | **文字成立／载体缺一枚**（今天只能这么判，且缺一枚那件我具名了） |
| Q6 名册 13 PASS 与两发门 | **成立** | 13/0/0 复跑出来，与 35-v4 那份 12 枚名册作差＝**恰且仅**新用例那一枚；`go vet` rc=0、`d22scan` rc=0 clean |

---

## Q1｜(i) 那枚新用例是不是真的有牙 — **成立**

尺件：`scripts/b1-countershape-v5.sh`、`logs/b2-mutant-landing.txt`、`logs/b6-overlay-selfcheck.txt`、`logs/run-r{0..11}*.txt`、`logs/b3-run-summary.txt`、`logs/e1-r11-and-roster.txt`。
⛔ 一份写腿的 overlay 件都没复用：所有突变源由本腿从工作树现场重造（`diff` 逐枚落 `changed-line-count=2`，见 `b2`）。

**正控两发（缺一不发凭据；本仓被 `{mutant:mutant}` 空操作与键路径带前缀咬过两次）**
- `r0-pristine`（同字节产码进 overlay）＝ `buildrc=0 runrc=0`、**PASS=10 FAIL=0 SKIP=0**、`exit status 0xc0000135` 命中 **0** ⇒ 载体活着、overlay 不是空操作。
- `r1-zsyntax`（产码塞故意语法错）＝ **`buildrc=1`** ⇒ overlay 的**键路径确实被 Go open 了**。
- ★本腿自伤登记（**这是"正控抓到空操作"那一类的真实样本，不是我抓到被测物是假的**）：第 1 遍我把 overlay 的**替换侧**写成 MSYS 路径 `/d/tmp/…`，Go 按当前盘解析成 `D:\d\tmp\…` ⇒ **8 发全 `buildrc=1`、PASS=0 FAIL=0**。是 `r0-pristine` 本该 `rc=0` 却给了 1 这枚反常暴露的（不是某发红）。第 2 遍 `runrc=127`＝`sherpa-onnx-c-api.dll: cannot open shared object file`：只把 DLL 目录塞进 `PATH` 不够（那条 PATH 项含空格），定式＝**三枚 DLL 拷到测试 exe 旁边 ＋ CWD＝包目录**，两样都做了才取到读数。builder 现在拒绝任何非盘符绝对路径的键/值并逐枚查存在（`b6`）。

**它没试过的两道守卫写法（新坏法 a）**
| 发 | 产码 `:680` 换成 | 红名册 | 判 |
|---|---|---|---|
| `r2-a1` | `typeof window.wispDispatch === "undefined"` | **FAIL=1**，只有新用例，句带 `M-B3 RED` | **红**＝有牙 |
| `r3-a2` | `window.wispDispatch == null` | **FAIL=1**，只有新用例，句带 `M-B3 RED` | **红**＝有牙 |

两发都与真值判断一样"只在门存在但非函数那个世界分岔"（`undefined`/`null` 世界里旧那枚 `:2032` 照旧走原生出口⇒绿，函数世界两者都不分岔⇒全绿），所以它们不是"再证一遍同一件事"：它们说明**这枚用例抓的是一整类把 `typeof` 半支削弱成取值判断的写法**，不是只对 `!x` 那一形敏感。

**反向弱形（b，只在 overlay 上做，跟踪件零改动）**
- `r5-weak-shape`（只削夹具、产码不动）＝ **PASS=10 FAIL=0** ⇒ 弱形照样编译照样跑，不是 build 假绿。
- `r6-weak-truthy`（弱形 ＋ 真值判断突变）＝ **FAIL=1** 仍只有新用例 ⇒ "不抛就算过"那半支**单独**就已经能抓住真值判断那一发。
- `r11-weak-silentdrop`（弱形 ＋ 静默丢件突变 `return;`）＝ **新用例 `--- PASS`**，名册 4 PASS／6 FAIL。
- `r4-a3`（**原强形** ＋ 同一发静默丢件突变）＝ 新用例 **红**，红句 `M-B3 READING WRONG: the fallback must reach the native exit exactly once…`。
⇒ **这一对才是"用例不是装饰"的决定性证据**：同一发坏法，强形红、弱形绿 ⇒ 承重的是那五支非抛断言（原生出口次数／receiver／逐字转发／死在 unbound 槽）。派单里那句"确认它〔弱形〕会红"我按字面跑会得**相反**的颜色，所以我把两发配对跑并这样读：**弱形红＝只有抛那一支有用；弱形绿而强形红＝断言整族有用**。今天量到的是后者，更强。
- 还原自证：`logs/b5-restore.txt` ＝ 跑前跑后 **14 枚**（夹具＋产码＋**12 枚 winlive**，`grep -rlE '//go:build.*winlive' --include=*.go`＝cmd/wisp 7＋internal/ball 5）逐枚 `git hash-object` **名册全等**、`git status` 对两枚被引件为空、产码 blob 仍 `26b5de83b93a9a141f1546dc19a9b03ffa45ede0`。

## Q2｜"唯一红"站不站得住 — **成立**

复跑写腿的 `b`（`r7-truthy-new`，我自己的突变源）：`buildrc=0`、**PASS=9 FAIL=1 SKIP=0**，红名册逐枚抄回**只有 1 枚**：
`--- FAIL: TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsNotCallable (0.00s)`，红句 `M-B3 RED (guard reduced to a truthiness test)…TypeError: string is not a function`。
⇒ **不是级联红**，与上一轮 `M-A` 那种要靠成对摘除读数才能归因的形不同：这枚读数的价钱是"专指新面"，可以直接引。
对照 `r10-typeofgone`（摘掉整个 `typeof` 半支）＝ FAIL=2（`:2054` 门不在 ＋ `:2123` 门不可调用），红名册与 `r7` 的 1 枚不同 ⇒ 两台探测器不是同一台，写腿 (i) 那句"名册与 `b` 不同"复现成立。

## Q3｜(ii) 拆前缀有没有把牙拆掉 — **成立（判据没变松）**，另抓到一枚我自己的题外新面

**① 枚数尺**（`logs/a03-prefix-census.txt`）：`M-B RED`=**0**（`daf1f5a0` 上是 6）、`M-B1`=**9**、`M-B2`=**5**、`M-B3`=**9**、`M-A RED`=**2**。⚠ 尺写法必须带**尾空格**：`grep -c 'M-B1 RED'` 只给 **4**，另 5 枚是 `READING WRONG`／`CONTROL RED`／`face has teeth` 一族（派单与 A694 写的 9 只有在前缀尺下才对，引这组数要连尺写法一起引＝本仓第 129/130 条那个形）。
`t.Fatalf` **47→56**（增量 9）逐枚归因过：新用例自己就是 9 枚 `t.Fatalf`＝1 枚 `installPanelTransport failed` ＋ 3 枚 `M-B3 CALIBRATION RED` ＋ 5 枚断言 ⇒ **既有用例一枚没增一枚没减**；`t.Errorf` 9→9、`func Test` 8→9。

**② 行为尺（挑两枚被改过文案的既有用例，跑它们**原有**那发该红的突变）**
| 被改的那枚用例 | 它原有的该红突变 | 色 | 红句认不认得出用例 |
|---|---|---|---|
| `TestForwardingHookIsSilentlyUnarmedByANonWritableNativeExit`（4 条句子前缀 `M-B RED`→`M-B1 RED`） | 夹具 `jsObject.set` 不再拒写（v4 的 `v-set-gone`） | **FAIL=1**，红的正是它自己 | `1 M-B1 RED` ＋ `--- FAIL` 名 ⇒ 认得出 |
| `TestForwardingHookMustNotLoseTheNativeExitReceiver`（红句 :1960 被 (iii) 就地改写） | 产码 `native.call(cw,…)`→`native(message)`（v4 的 `v-ma-hook`） | **FAIL=7**（级联，与写腿/`35-v4` 同读数），红名册含它 | `1 M-A RED` ＋ `--- FAIL` 名 ⇒ 认得出 |

⇒ 两枚都**仍红且红句指得认对应用例**，**没有任何一枚的判据被改松**；(ii) 越格自报的那 14 枚整族拆分在行为面上成立。
⚠ 上面 `v-set-gone` 我**没按行号**取（按行号会取错，见下），改成按内容 `s/exists && o\.unwritable\[name\] {/exists && false {/`，落地尺 `changed-line-count=2`。

**③ 顺手查那件事：有没有"活的仪器"去 grep `M-B RED`** — 有，而且比派单问的更严重一层，具名报：
- **CI／tools／scripts 三棵树零命中**（`a06` 第 (3) 支）⇒ 没有任何门今天起静默找不到东西。这一格**好消息**。
- 字面 `M-B RED` 的活仪器在 `.scratch` 里，三枚：`probes/35/v4/scripts/countershape-v4.sh`、`probes/35/r5/scripts/anchor-rulers.sh`、`probes/35/r5/scripts/after-rulers.sh`。
- ★**新面（不是派单问的、也不是写腿自报的）**：`countershape-v4.sh` 是**按行号**给夹具开刀的（`mkmut = awk NR==n`），这次 (ii)/(iii) 的行位移把它 7 枚 fixture 突变锚里的 **6 枚打死了**（尺件 `a08`，逐枚给了 before/after 两 blob 的同一行＋真身现在的行号）：
  `1288→注释行`、`1290→注释行`、`1947→注释行`、`1977→注释行`、`1954→空行`、`2048→误落在 if bf.window.has(...)`；真身已移到 1293/1295/1953/1983/1956。只有 `164` 那枚侥幸没动。`daf1f5a0` 上 6 枚**全 LANDS** ⇒ 位移由 35-r5 造成，因果有对拉。
  今天重跑它多半是 **`buildrc=1`（响）**而不是静默恒绿，但 `v-label-ma`/`v-label-mb` 那两枚"按文案归因"的反形**从此再也造不出它们声称的那个世界**＝一枚躺在仓里的、名字与用途只活在注释里的仪器。**产码那侧没事**：`675/680/682` 三锚仍逐字对得上（`a09`）。
- `anchor-rulers.sh:29` 那条把 `M-B RED` 打在 **"--- (ii) ruler BEFORE"这个标题下**却读**工作树**（`after-rulers.sh` 读的是 `git show daf1f5a0` 的落盘副本＝blob 锚定，没这毛病）⇒ **今天重跑 anchor-rulers.sh 会打印一行 `M-B-RED-count=0` 挂在 BEFORE 标题下**。这是"旧读数被当成当前读数"的形，具名入账（⛔ 我不改，历史记录只追加不删）。

**④ 与派单转述冲突（具名）**：派单/A694 说"文档/票面/台账里 **19** 处"。我七把尺都取不到 19（`a10`）：跟踪件全仓 **60 行／25 枚文件**、非 `.go` 60 行、`*.md` **13 行／4 枚**、非 `.go` 剔掉 `probes/35/*/logs/**` ＝ 20 行、`.go` **0**（这枚与转述同）。⇒ A694 里那枚 19 是**枚数不明的一把尺**，请后续程以「60 行/25 枚（全跟踪件）」「13 行/4 枚（`*.md`）」这两把带口径的读数作准；A694 那句"⛔ 后续程不许拿『文档还写着 M-B RED』去夹具里找用例"我完全同意，且与枚数无关。

## Q4｜(iii) 文案改完，免责声明还在不在射程内 — **成立**

尺件：`logs/c1-gates.txt` 之外的直接 `sed` 逐字抄（正文在判决里，件＝`impl` 的 :1900-1920 段与 `a04`）。
- **那块"这格不买什么"还在**，现量在 `cmd/wisp/panel_transport_35r2_test.go:1911-1920`：`WHAT THIS BUYS AND WHAT IT DOES NOT. It buys "a hook that drops this, or whose assignment is silently refused, is no longer certified as a delivery". It does NOT buy "WebView2 really binds the receiver / really lets the page replace postMessage" — A684 §1 puts that格 with ticket 35:52's real-window evidence, and no assertion here stands in for it.` 后面 `Object.defineProperty` 未实现、descriptor READS/configurable/delete 具名列为 un-modelled 那半段也逐字在。⛔ 没被顺手删。
- 三处收紧逐处对拉：`:1285-1291`（`Chrome and WebView2 **are reported** to answer…` ＋ `THIS RULER CANNOT PRODUCE THAT CREDENTIAL (ticket 35 AC#8(iii))` ＋ `THIS FIXTURE'S OWN ANSWER, modelled after Chrome/WebView2 and **never offered as a WebView2 behaviour record**` ＋ 指回 winlive 件并写 `〔仅本机可量〕`）；`:1944-1946`（`the rule THIS FIXTURE answers with, modelled after what Chrome and WebView2 are reported to do (this yard never observes a browser; AC#8(iii))`）；`:1960` **在 `t.Fatalf` 运行时文案里**（`THIS FIXTURE answers that with Illegal invocation - modelled after what Chrome and WebView2 are reported to do, **NOT a WebView2 behaviour record** (AC#8(iii); the real-window reading is the 〔仅本机可量〕 winlive rig)`）。
- ⛔ 无残留未限定的浏览器断言：`grep -c 'Chrome and WebView2 answer an\|the way a browser.s host method does\|A real browser answers that with Illegal invocation'` ＝ **0**（`after-iii.txt` 同尺，写腿也报 0，我对拉一致）。
- 红句仍指得认用例：`M-A RED (receiver lost at the native exit)` 前缀仍在（2 枚命中），`%d/%s/%s` 参数序未动，`r8-mahook` 那发红名册里这枚在、红句可读。
- **两格分开写**：文字更谦虚＝**成立**；判据变松＝**没有**（`r8` 那发它照红，`r10` 那发它没被惊动）。

## Q5｜(iv) 那句〔仅本机可量〕是不是四要素齐 — **文字成立／载体缺一枚**（今天只能这么判）

尺件：`logs/a05-iv-label.txt`、`logs/d1-q5-label.txt`、`logs/d2-q5-followup.txt`。
- 四要素**齐**，且是带序号的四支：①`WHAT IS MACHINE-LOCAL`（现象＝`:52` 那个"页面自己发的 postMessage 到不到 Go 门"只有本机开真窗能取，夹具"models that edge, it does not measure it, and its green is not this credential"）／②`WHO OWNS IT`（**归编排者**，并逐字禁三条引用口径："not in a ticket, not in an evidence table, not in a commit message"，还加了一句"A leg that has not opened a window has not read this file's result"）／③`RE-RUN CADENCE`（每波、桌面空时、`-tags winlive` 一发**并同波再跑一发控制组**、件名＋`rc` 逐波进 `.scratch/wisp/probes/35/**`）／④`THE PRICE PAID`（CI 永远不会拦下"这条传输在真浏览器里坏掉"那类回归）。
- **撤销口令在**：`:21` `Revocation口令 for this label: give ':52's decisive reading a CI-reachable stand-in (then this block is rewritten to name the stand-in, not deleted).`
- **"一处写两处"齐**：文件头 1 处 ＋判决件 `probes/35/r5/impl.md` 4 处（票面 `:81`、台账另立一格也在）。
- **同形核对**：文件头用的是 `〔仅本机可量、CI 永看不见〕`（顿号形）＝先例 `docs/evidence/s1/33-panel-host-c27-r5.md:78` **逐字符同**。⚠ 一处与我无关的形差供参考：台账里有 1 枚写成 `〔仅本机可量，CI 永看不见〕`（逗号）；票面 11 枚是简形 `〔仅本机可量〕`。写腿件没走形。
- ★**欠的具体是哪一件（具名，⛔ 我不替它判"载体已建"）**：③那句口径要求"同波控制组＝同一台件跑在 `fb2fb802` 那枚 hook 上、`go test -overlay`、仓不动"，而我现量：
  1. 全仓**没有任何可跑的控制组台件**——`control-group`/`同波控制组` 只活在 3 处散文里（该文件头／`impl.md`／台账），`probes/35/**` 下零枚控制组件；
  2. 全仓**没有任何一枚 `-tags winlive` 的 runner 脚本**（`.sh/.ps1/.py` 里只有 `v4/scripts/countershape-v4.sh` 提到过它，且它跑的是默认档）；仓里唯一带 winlive 的自动化是 `ci.yml:595/644` 那步 `go vet -tags winlive ./cmd/wisp/ ./internal/ball/`＝**只补编译、不开窗、取不到 `:52` 的读数**；
  3. 引用的 `fb2fb802` **是真提交**（`git cat-file -t`＝commit，全号 `fb2fb802f75a0e3eeacad488f1adc6f064e29f85`）⇒ 那半句不是指空。
  ⇒ 要补的那一件＝**一枚 named runner**（建议 `.scratch/wisp/probes/35/<wave>/scripts/winline-control.sh`：同一 exe 跑两遍，一遍 shipped 钩子、一遍 `-overlay` 换成 `fb2fb802` 的钩子，两遍的 `-test.v` 件名＋各自 `rc` 落盘），⛔ 不是补注释、⛔ 不是把 ③ 的口径再润色一遍。这一格归**编排者**（A694 已这么记，我复认）。
- ⛔ 本腿**没有**跑任何 `-tags winlive`（派单禁），因此我不声称读过这个文件的任何结果，也不声称 winlive 那 12 枚今天能编译（我没跑那步 `go vet -tags winlive`；那是 CI 那一步的格）。

## Q6｜名册 13 PASS／0 FAIL 能不能复跑 — **成立**

尺件：`logs/run-r0-pristine.txt`、`logs/e2-roster-reconcile.txt`、`logs/c1-gates.txt`。
- 载体照派单定式：`CWD=cmd/wisp` ＋ sherpa 三枚 DLL（`ls third_party/sherpa-onnx/*.dll`＝3 枚，现量）；`exit status 0xc0000135` 哨兵在每一发里命中 **0**（否则就是"根本没跑"的假绿）。
- **13/0/0 复跑出来**：`delivered-baseline.txt` 的名册我逐名抽回＝13 PASS／0 FAIL／0 SKIP；作差 `comm` 对 `35-v4/logs/names-baseline.txt` 的 12 枚：
  - 13 有而 12 无 ＝ **恰 1 枚：`TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsNotCallable`**（正是派单要求的那枚）；
  - 12 有而 13 无 ＝ **∅**。
  ⇒ **作差为空／非空两侧都干净，这格成立。**
- ⚠ 一枚必须带上的口径差：我自己复跑写腿那套 `-run`（`TestForwarding|TestPagePostMessage|TestShapeA3|TestLegacySubShape|TestUnforwardedPageEnvelope`）得到的是 **10 PASS／0 FAIL**，不是 13——少的正是 3 枚 `TestInbound*GuardRefuses*OnPageEdge`（在那枚 `daf1f5a0` 之后的 `-run` 里没点名）。写腿在 `impl.md:71` 自己就写了"旧名册 9／本腿新名册 10"，⛔ 三个分母（10／12／13）不许互相作数；上面那次作差用的是**同分母**（13 vs 12），所以成立。我另 8 发反形都跑在 10 枚那把尺上，红名册读数的分母＝10，已在各件里写明。
- 两发门（`logs/c1-gates.txt`，各带 `rc`）：`go vet ./cmd/wisp/` **rc=0**（零输出）；`cd tools/d22scan && go run . -root ../../` **rc=0 clean**（266 production Go files；ban #8 覆盖 `cmd/`=108、`internal/`=514，含注释与 `_test.go`）。⛔ 没写成 `go run ./tools/d22scan`。
- ⚠ **与写腿自报冲突（具名）**：它写 `gofmt -l` 列 **1** 枚。我这条 `gofmt -l cmd/wisp/` 列 **3** 枚＝`models.go`＋`panel_inbound_guards_35r3_test.go`＋`panel_transport_35r2_test.go`。拆开后**它的实质结论成立、数不对**（`logs/c2-gofmt-drift.txt`）：把 `daf1f5a0` 的 blob 落到临时目录跑同一把尺＝**before 已经列 2 枚**（夹具＋guards，夹具那枚的 hunk 就是 `jsParser peek/next` 那 4 行对齐，与 35-r5 的改动逐字无关）；多出来的 `models.go` 是 **CRLF 假象**（工作树该文件带 334 枚 CR 字节、`git show` 出来是 LF，`gofmt -d` 要对整文件 334 行重写）。⇒ 引 `gofmt` 枚数必须写"**哪棵树＋过没过 autocrlf**"（本仓第 124 条同一个机制）。写腿"没顺手改无关行"这句我也验了：漂移没被它抹平＝对的。

---

## 我这把尺今天看不见什么

1. **任何真浏览器事实**：⛔ 没开过窗、⛔ 没跑 `-tags winlive`。`r2/r3` 里"字符串被调用会抛 `TypeError: string is not a function`"是**这台手写解释器**的答案，不是 WebView2 的答案——而这正是 (iii) 刚被收紧的那件事，我自己也得受它约束。
2. **winlive 那 12 枚能不能编译**：我没跑 `go vet -tags winlive ./cmd/wisp/ ./internal/ball/`（那是 CI `ci.yml:644` 那一步的格，且 A693/A694 记着它的真实颜色欠一次 push）。
3. **整包颜色**：`cmd/wisp` 346 枚、单程 ~14 分钟、其中 5 枚是窗口依赖既有红 ⇒ 我全部读数都跑在 10/13 枚名册上，**不等于包绿**。
4. **`postMessageHopCap` 的上界**：8→80 那一发写腿没跑、我也没跑 ⇒ "帽数被钉死"这句我同样不认，只认下界（8→2／8→1）已有凭据。
5. **再入标志的释放**（`finally { inside = true; }` 那一发只红 `:1764`）与"去掉门那支的 `return`＝等价形"——票面 `:82` 自己划的两条覆盖面边界，我没越，也没重新验。
6. **弱形/强形那对读数只覆盖 `return;` 这一发静默坏法**：还有别的"不抛但送错"的写法（例如 `native.call(window, message)` 之外的 receiver 变形、把帧改写成半截 JSON）我没逐发配对，所以"非抛那五支各自独立承重"我只证到**其中一支以上**，没证到五支全独立。
7. **(ii) 那 14 枚的句子语义**：我只量了"枚数没变＋原有突变照红"，没有逐枚判"这句新文案是否still在描述它原来描述的那面"。
8. **CI**：⛔ 零 push ⇒ 任何"门在 CI 里有牙"的问题今天都取不到读数（第 109/124 条同一格）。
9. `probes/35/v4/**` 那批历史件我现在**只诊断没修**（它们不是我的写面，且只追加不删）；`a08` 那枚仪器位移我只记了行号对拉，没逐枚重跑 v4 的 11 发。

## 我推翻编排者哪句

1. **`起手 HEAD 应为 396c39b6`** ⇒ 现量 `d2a4597d…`（你自己 A694 那一笔）。已核 `git diff --name-only 396c39b6..HEAD` 只有 `docs/reports/pending-and-issues.md`（16 增/0 删）⇒ 被验面没变，但派单里那枚锚号是过期读数，后续程照它 `git log -1` 会对不上。
2. **`文档/票面/台账里 19 处`（A694 写死的那枚）** ⇒ 我七把尺取不到 19，最接近的两把是 20（非 `.go` 剔 probes logs）与 13（`*.md`）；诚实的全量是 60 行／25 枚文件。这句是**台账里的一枚数**，不是我推翻了判据——是请你在 A695 用带口径的数替换它的口径。
3. **`M-B1=9／M-B2=5／M-B3=9`（不带尺写法）** ⇒ 数是对的，但那把尺**必须**是"前缀＋尾空格"；`grep -c 'M-B1 RED'`＝**4**。你自己第 129/130 条立过同款规矩，这两枚数在 A694 里没带尺写法。
4. **`gofmt -l 列 1 枚`（写腿自报、派单原样转述）** ⇒ 该条工作树列 3 枚，`models.go` 那枚是 CRLF 假象，另两枚 `daf1f5a0` 已列；结论"既有漂移、未顺手改"我**追认**，那句"列 1 枚"要补口径。
5. 我**没有**推翻 `Q2 唯一红`／`Q4 免责声明仍在`／`Q6 名册作差恰 1`／`产码零改动`／票框 8/2 这五句——它们逐条复跑同色。

## 具名欠账（谁欠、欠哪一件）

| # | 欠什么 | 归谁 | 现量出处 |
|---|---|---|---|
| 1 | AC#8(iv) ③ 那半句"同波控制组"**没有可跑台件**、全仓零枚 `-tags winlive` runner | 编排者（A694 已记，本腿复认） | `d1`/`d2` |
| 2 | `probes/35/v4/scripts/countershape-v4.sh` 的 7 枚**按行号**夹具锚里 6 枚被行位移打死；要么改成按内容锚、要么在件首写明"只在 `daf1f5a0` 有效" | 有 `probes/**` 写面的下一枚腿或编排者；⛔ 本腿不改历史件 | `a08`/`a09` |
| 3 | `probes/35/r5/scripts/anchor-rulers.sh:29` 把读**工作树**的 `M-B RED` 计数挂在 `BEFORE` 标题下 ⇒ 今天重跑给 `=0` 且看起来像改前读数 | 同上（`after-rulers.sh` 无此毛病，可照它改成 blob 锚定） | `a06` |
| 4 | 台账 A694 那枚"19 处"要落一枚带口径的更正（只追加） | 编排者 | `a10` |
| 5 | AC#8 帽数**上界**那发（8→80）仍没人跑 | 下一枚验收腿（本腿没跑） | 票面 `:82` |
| 6 | `:52` 真窗读数＋那 12 枚 winlive 的编译颜色（本腿禁跑） | 编排者那一波 | `d2`、`ci.yml:595/644` |

## 本腿交付面

写过的东西全在 `.scratch/wisp/probes/35/v5/**`；件名不以任何 `.out` 结尾；逐件自落 `rc=`（`a01-a10`、`b0-b6`、`c1-c2`、`d1-d2`、`e1-e2`、`run-r0…r11`、`winline-files.txt`）。
⛔ 票面框枚数追加前后各量：8 未勾／2 勾，一字未动。
