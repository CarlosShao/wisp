# 88 — 把 `ban #6` 从"豁免"翻成"真覆盖"（`frontend/` 已经进树了，豁免就变成谎话）

**Status:** ready-for-review
**next:** 编排者验收——AC#1~AC#6 六框已勾（数字逐条挂在 log 第 2/3/5/6 条），代码沿用 `84e4161`，
          本票只补票面证据。**push 的 blocker 已从本票转移到票 89**（见 log 第 7 条的红账）。
**Type:** 门禁完整性（A26 / A44① / A54② / A56④ 的同族，**这一次是轮到 ban #6 自己**）
**Blocks:** **编排者的 push**——票 77 的 `frontend/` 已落在 `63ef895`（未推），
             现在推上去会让 `lint` 的 D22 两步红，从而**把刚拿到的 `go vet` 绿证又挡成 skipped**。
             本票闭掉之前 push 压着不走。
**Blocked by:** nothing
**Packages:** `tools/d22scan/main.go`（`declaredScopes()` 那一条 + 必要的 `walkText` 过滤器）、
              `tools/d22scan/scan_test.go`（受影响的那几条台账用例）。
              **禁改**：`allowlist.txt`（**必须仍是 5 行非注释**）、任何 ban 的**文本**（D22：不是代理能缩的）、
              "空作用域致命"这条规则本身、`frontend/**`（票 77 的地界）、`internal/**`。

## 背景（编排者已经试过一次，把账留在这是给你少走弯路）

票 67 当年把 `ban #6 frontend/` 登记成 **exempt-with-a-note**（`live:false, absentOK:true`），
理由写得很好："目录不存在 ⇒ 这条只能走 0 个文件，留着就是假装在扫"。
并配了 **drift guard**：目录一旦冒出来而没人翻牌，扫描器**致命报错**（fail-closed）。
票 71 的 AC#4 又把"0 覆盖"从沉默变成明打字的 `NOT COVERED`。

**现在票 77 把 `frontend/` 建出来了** ⇒ 豁免的前提已经不成立。编排者在 `main.go:306-314`
试把 `live: false, absentOK: true` 改成 `live: true`，`cd tools/d22scan && go test ./...` **当场 5 条红**：

```
--- FAIL: TestVerdictGreenOnFullyLiveFixture
--- FAIL: TestVerdictRedOnEmptyBan7Scope
--- FAIL: TestExemptScopeCannotOutliveItsAbsentTree
--- FAIL: TestRealRepoLedgerIsHonest
--- FAIL: TestBuiltBinaryGoesRedEndToEnd
      ├─ /seeded_violation_exits_1
      ├─ /exempt_scope_whose_tree_appeared_exits_2
      └─ /fully_live_fixture_exits_0
报错原文：d22scan: scope ban #6 frontend/ examined 0 files but is declared live
          in declaredScopes() - an empty instrument is not a verdict (ticket 71 AC#4)
```

⇒ **这套台账是在"ban #6 永远豁免"的假设下写的**，翻牌要把它们**逐条重新表述**，不是删掉。
⚠ 那条"live 作用域若 examine 0 个文件 ⇒ 致命"**是对的，是你最不该碰的东西**：
它存在的意义就是"空仪器不算判据"（票 78 的代理当年也是按这条拒绝我的省事修法）。

## 先回答一个问题再动手（AC#1）

**真实仓里 `ban #6` 翻牌之后到底能扫到几个文件？** 编排者的实验里 fixture 是 0，
真树里 `frontend/` 有 38 个非 `node_modules` 文件（`.tsx`/`.ts`/`.css`/`.mjs`/`.md`）。
⇒ 如果 `walkText` 的后缀过滤器**不收 `.tsx`/`.ts`**，那真树也会量到 0 ⇒
**那才是本票真正要修的洞**：一条"看起来已覆盖、实际永远 0"的门。
先给出真实读数（`examined N text files`）再决定改哪里，**不许**为了让 N>0 而把 ban 的语义改宽或改窄。

## AC（1:1，裁决表 `docs/evidence/s1/88-*.md`）

- [x] **AC#1**（数字见 log 第 2 条建立、第 6 条复测）真实读数：在 `git archive HEAD` 的**仓外纯净快照**（目录带你的会话后缀）里跑
      `sh scripts/d22scan.sh`，报出 `ban #6 frontend/ examined N text files` 的 **N**，
      并说明过滤器现在收哪些后缀。若 N=0 ⇒ 明写"这条门翻牌后仍是空仪器"，**这就是本票的缺陷本体**。
- [x] **AC#2**（落地见 log 第 3 条，`live: true` / 无 `absentOK` 由我在第 6 条复核） 翻牌落地：`declaredScopes()` 里 ban #6 变 `live: true`（去掉 `absentOK`），
      `note` 改写成"谁在哪一批翻的、为什么"（保留 D22 那句"ban 文本不是代理能缩的"）。
      **ban 的判定文本一字不改**。
- [x] **AC#3**（重述见 log 第 3 条，真绿由我在第 6 条以 `=== RUN`=28 / PASS=28 / SKIP=0 复测）上面 5 条红**逐条重新表述并真绿**：
      涉及豁免语义的（`TestExemptScopeCannotOutliveItsAbsentTree`、
      `/exempt_scope_whose_tree_appeared_exits_2`）要**保留 drift guard 的可测性**——
      办法是**用一个合成的 exempt 作用域来测这条规则**，而不是靠"ban #6 恰好豁免"这个巧合；
      **规则本身必须仍然有牙**（不许把它测成永真）。
- [x] **AC#4**（数字与报错原文见 log 第 5 条）阳性对照：`/seeded_violation_exits_1` 用的就是
      `frontend/src/app.js` 里那句 `approval.decide(...)`（`scan_test.go:80`）⇒
      **翻牌后必须还能因这条种子违规红**，并且红的时候打印出 `ban #6` 命中。
      这是本票"门真的有牙"的证明，**不许用"改种子文件名/后缀"绕过 AC#1 的过滤器问题**。
- [x] **AC#5**（我自己那次纯净快照的逐作用域原文见 log 第 6 条；⚠ 当前 HEAD 另有一笔**非本票**的红账在第 7 条）全套绿之后跑 `sh scripts/d22scan.sh`（**逐字同 CI**）于纯净树：rc=0，
      且台账里 `ban #6` 那行的 N>0；同时确认 `allowlist.txt` 仍 5 行非注释、
      `git diff` 里**没有任何 ban 文本被改动**。
- [x] **AC#6**（真树 0 命中见 log 第 3 条建立、第 6 条复测：纯净快照 0 命中 / 35 文件）如果 `frontend/` 里的**真实代码**（票 77 vendored 的 Approval Card）命中 `ban #6`：
      **不要动扫描器也不要加豁免**——把命中点登记给编排者，由**票 77 那侧改名**。
      本票只负责让门有牙，不负责替别人消音。

## Rules（本仓固定）

15 次工具调用内交第一枚 checkpoint commit；每完成一组**就把结论追加进票面**（别攒，本仓代理死在轮数上限上过）；
`git commit -q -F - -- <显式路径>` + 带引号 heredoc；禁 `git add -A`/`.`；commit 前核对
`git diff --cached --name-only`（**别把活着代理的票面 add 进去**，编排者今天栽过一次）；
禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`；**不 push**；不在仓内建 worktree（A38④）；
四种假绿逐跑点名；数红/绿用全量输出仪器；票面 append-only，**要改的那行先读再替换**。

## Progress log（append-only）

- 2026-09-21（编排者）：建票。我做过一次 `live:true` 的试验并**已 `git restore` 还原**（工作树干净），
  5 条红的名单与报错原文在上面。我的判断：本票的价值不是"翻一个布尔"，
  而是**回答 AC#1 那个 N**——如果过滤器不收 `.tsx`，那我们今天所有关于"ban #6 已覆盖面板"的说法都是空的。
- 2026-09-21（票 88 代理）：**AC#1 答：N = 35，不是 0。假设里的那个洞不存在，但差一点就存在。**
  测量方法：`rm -rf /tmp/wisp88-base && mkdir -p /tmp/wisp88-base && git archive HEAD | tar -x -C /tmp/wisp88-base`
  然后 `cd /tmp/wisp88-base && sh scripts/d22scan.sh`。翻牌**前**的输出已经把 N 量出来了
  （walk 一直在跑，只是账上记成 exempt）：
  `d22scan: scope ban #6 frontend/         examined  35 text files  [NOT COVERED]`，
  并且**同一跑当场 rc=2**：`scope ban #6 frontend/ is registered as absent-but-exempt while its directory EXISTS`
  ⇒ 编排者"push 会被挡"的判断成立，纯净树 HEAD 上 `sh scripts/d22scan.sh` 现在就是红的（rc=1，因为 step1 的 go test 先 FAIL）。
  过滤器原文依据（`tools/d22scan/main.go`）：`walkText(dir, ban, check, goOnly)` 的后缀判断只有一句
  `if goOnly && (!strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go")) { return nil }`，
  而 ban #6 的调用是 `walkText(filepath.Join(root,"frontend"), "panel-approval", s.panelCheck, **false**)`
  ⇒ `goOnly=false` 时**根本没有后缀过滤器**：frontend/ 下除 `node_modules/`、`testdata/`、`.git/` 之外的
  每个文件都被读、被逐行扫（`.tsx`/`.ts`/`.mjs`/`.css`/`.md`/`.json`/无后缀全收）。
  `isTextFile()`（那份后缀清单，含 .ts/.tsx/.jsx/.css/.md）**只服务 ban #8 的 design/**，与 ban #6 无关。
  35 个文件的构成：16 tsx / 6 json / 3 ts / 3 mjs / 2 css / 1 md / 1 html / 1 go / 1 .gitkeep / 1 .gitignore。
  ⇒ **本票的缺陷本体不是"空仪器"，是"豁免本身变成谎话"**（drift guard 已经把它抓住了）。
  ⚠ 反手一条给未来的自己：**不许**为了"整齐"给 walkText 的非 goOnly 分支加 `isTextFile` 白名单——
  那是把 35 缩成 33 并漏掉无后缀文件的**收窄**（R16#4），且会被 `TestLedgerCountsMatchAnIndependentWalk` 的 ban #6 行抓住。
  为防止"armed but blind"另加了一条正向钉子 `TestBan6ScopeIsNotNarrowedByAnExtensionFilter`。
- 2026-09-21（票 88 代理）：**AC#2 落地 + AC#3 五条红全部重新表述完毕**（同枚 commit）。
  `declaredScopes()` 里 ban #6 = `live: true`、`absentOK` 去掉、`note` 改写为"谁在哪一批翻的、为什么"，
  D22 那句"ban 文本不是代理能缩的"保留。**判定文本/正则一个字未改**（见 AC#5 的 diff 证明）。
  "live 作用域 examine 0 个文件 ⇒ 致命"这条规则**未动**，翻牌后 ban #6 自己也开始受它管。
  AC#6：**真树 0 命中**——`grep -rn 'approval\.decide' frontend/`（排除 node_modules）无匹配，
  `cd tools/d22scan && go test ./...` rc=0，其中扫真仓的 `TestScannerSelfScanOfRealRepoIsGreen` /
  `TestRealRepoLedgerIsHonest` 都在跑 ban #6 ⇒ 票 77 vendored 的 Approval Card 不需要改名。
- 2026-09-21（票 88b 代理接续；Status→ready-for-review）：**接续说明**（我断在哪、前任欠什么）。
  前任（票 88 代理）在 `84e4161 feat(88,AC#1,AC#2,AC#3): arm ban #6` 里做完了 **AC#1 的测量（N=35）、
  AC#2 的翻牌、AC#3 的 5 条台账重述**，但**六个 AC 框一个都没勾、AC#4 的阳性对照没跑、AC#5 没跑、
  Status 还停在 `open`**——票面上没有任何一处断言"这套绿过了"。我从**断点第 1 项 AC#4** 接上（第 5 条），
  再跑 AC#5（第 6 条），然后回头用**我自己重测的数字**勾框。**代码我一行未改**：
  `git log --oneline 84e4161..HEAD -- tools/d22scan` 输出**空**（HEAD 在 `de15a6b`，其间只有票 89/90/91 的东西），
  所以本票的 commit 只动 `.scratch/wisp/issues/88-*.md` 这一个文件。
  会话后缀：快照目录用 `/tmp/wisp88b-flip`（翻牌 SHA）与 `/tmp/wisp88b-head`（当前 HEAD），
  与前任的 `/tmp/wisp88-base` 不同名；**仓内没建 worktree、没 checkout**。
- 2026-09-21（票 88b 代理接续）：**AC#4 阳性对照做了，门真的有牙，红的原文如下。**
  命令：`cd tools/d22scan && go test -count=1 -v -run 'TestBuiltBinaryGoesRedEndToEnd/armed_ban_6' ./...` ⇒ 测试本身 **rc=0**，
  因为它断言的是**子进程** rc=1（这正是 AC#4 要的东西）。⚠ 用了 `-run` 就按本仓规矩点名匹配证据：
  该跑的 `=== RUN` 计数 = **2**（父 `TestBuiltBinaryGoesRedEndToEnd` + 子 `armed_ban_6_goes_red_on_the_panel_violation_exits_1`），**> 0**，
  不是"pattern 打空造成的假绿"。子进程原文（`scan_test.go:952` 的 t.Logf + 断言）：
  ```
  rc=1 stdout=d22scan: examined 15 production Go files under internal/ and cmd/ of
      C:/Users/swq/AppData/Local/Temp/TestBuiltBinaryGoesRedEndToEndarmed_ban_6_goes_red_on_the_panel3079960430/001
      d22scan: scope ban #6 frontend/         examined   2 text files
      frontend/src/app.js:1: [panel-approval] `approval.decide` in frontend/ is banned (D33/F2: allow decisions are native-side only)
  stderr=d22scan: 1 finding(s); D22 bans are not negotiable (see PLAN.md D22, tools/d22scan/allowlist.txt)
  ```
  ⇒ **红的时候打印出了 `ban #6` 命中**（`scope ban #6 frontend/` 那行 + `panel-approval` finding + 文件名 `frontend/src/app.js`）。
  种子就是票 71 那条老种子（`frontend/src/app.js` 里的 `approval.decide(...)`），**没改名、没改后缀**，
  所以不存在"用改文件名绕开过滤器问题"。同一枚测试的**全量** 5 个子用例（`-count=1 -v -run 'TestBuiltBinaryGoesRedEndToEnd'`）：
  `seeded_violation_exits_1` rc=1 / `empty_live_scope_exits_2` rc=2 / `armed_ban_6_goes_red_on_the_panel_violation_exits_1` rc=1 /
  `ban_6_tree_gone_while_declared_live_exits_2` rc=2 / `fully_live_fixture_exits_0` rc=0，`--- PASS: TestBuiltBinaryGoesRedEndToEnd (1.88s)`，
  `=== RUN` 共 6 行（父 1 + 子 5）、`--- SKIP` **0** 条（四种假绿逐条点名：无 SKIP 冒充 ok、`-run` 已证明匹配、未用 `-count=N>1` 故无 RUN/N 核对问题、无步骤被静默跳过）。
- 2026-09-21（票 88b 代理接续）：**AC#5 我自己那一次纯净快照跑完，rc=0；AC#1/AC#2/AC#3/AC#6 顺手用我自己的数字复测。**
  快照做法（**翻牌后的 SHA**，不是我改出来的东西）：
  `git rev-parse 84e4161` = `84e4161458b11c5b96718f3ac3e39ff34cf2bbe3` ⇒
  `rm -rf /tmp/wisp88b-flip && mkdir -p /tmp/wisp88b-flip && git archive 84e4161 | tar -x -C /tmp/wisp88b-flip` ⇒
  `cd /tmp/wisp88b-flip && sh scripts/d22scan.sh` ⇒ **rc=0**。逐作用域**原文**（⚠ 这是**我自己的读数**，
  与编排者 A62 那次独立；两者一致这件事是我复跑之后才发现的，不是引用她的数字当我的证据）：
  ```
  d22scan.sh: positive control - go test ./... (tools/d22scan)
  ok      github.com/CarlosShao/wisp/tools/d22scan    (cached)
  d22scan.sh: scan of /tmp/wisp88b-flip
  d22scan: examined 207 production Go files under internal/ and cmd/ of C:/Users/swq/AppData/Local/Temp/wisp88b-flip
  d22scan: scope bans #1-5 internal/      examined 187 production Go files
  d22scan: scope bans #1-5 cmd/           examined  20 production Go files
  d22scan: scope ban #6 frontend/         examined  35 text files
  d22scan: scope ban #7 internal/tools/   examined  16 production Go files
  d22scan: scope ban #8 design/           examined  16 text files
  d22scan: scope ban #8 internal/         examined 314 Go files, comments and _test.go included
  d22scan: scope ban #8 cmd/              examined  25 Go files, comments and _test.go included
  d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=187, bans #1-5 cmd/=20,
             ban #6 frontend/=35, ban #7 internal/tools/=16, ban #8 design/=16, ban #8 internal/=314, ban #8 cmd/=25
  ```
  ⇒ **AC#1 的 N=35 我复跑到同一个数**（>0，不是空仪器），`ban #6 frontend/=35` 明打字在 clean 行里。
  ⚠ 诚实一条：脚本第一步在我这次跑里打的是 `(cached)`（CI 的 `go test ./...` 本身不带 `-count=1`，逐字同 CI 我就照它跑），
  **所以我的门禁不吃缓存**：`cd /tmp/wisp88b-flip/tools/d22scan && go test -count=1 -v ./...` ⇒ **rc=0**（9.436s），
  `=== RUN` **28** / `--- PASS` **28** / `--- SKIP`+`--- FAIL` **0**；`gofmt -l tools/d22scan/` 输出**空**；
  `cd tools/d22scan && go vet ./...` **rc=0**（⚠ 从仓根跑 `go vet ./tools/d22scan/...` 必失败，原文
  `directory prefix tools\d22scan does not contain main module or its selected dependencies`——tools/d22scan 是**自己的 module**，不是代码问题）。
  禁改项核对：`grep -v '^\s*#' tools/d22scan/allowlist.txt | grep -c .` = **5**（仍是 5 行非注释，最近一次改动是 `38b3715` 票 70，本票没碰）；
  `git diff 84e4161^ HEAD -- tools/d22scan/main.go | grep -E '^[-+].*(banned \(|approval\.decide in frontend)'` ⇒ **0 行**，
  判定文本/正则一字未改；"live 作用域 examine 0 个文件 ⇒ 致命"这条规则**未动**，而 ban #6 现在正受它管
  （`main.go:319-328` 的 ban #6 条目：`live: true`、无 `absentOK`、`note` 写明"armed by ticket 88 … The ban TEXT is unchanged"）。
  **AC#6 我自己复测**：同一份纯净快照里
  `grep -rn --exclude-dir=node_modules --exclude-dir=testdata --exclude-dir=.git 'approval\.decide' frontend/` ⇒ **0 命中**（grep rc=1），
  `find frontend -type f -not -path '*/node_modules/*' | wc -l` = **35** ⇒ 真树里翻牌后的 ban #6 是**绿的且真有牙**，
  票 77 vendored 的 Approval Card **不需要改名**，本票没有替任何人消音。
- 2026-09-21（票 88b 代理接续）：⚠ **一笔不属于本票的红账，登记给编排者（AC#5 的"多样本全报"就是这一条）**。
  同一套命令在**当前 HEAD** `de15a6b` 的纯净快照（`/tmp/wisp88b-head`）里 **rc=1**：
  ```
  --- FAIL: TestScannerSelfScanOfRealRepoIsGreen
      scan_test.go:269: repo HEAD violates: internal/winsec/winsec.go:126: [pathresolver-bypass]
      filepath.Abs outside the C26 PathResolver is banned (D22); see tools/d22scan/allowlist.txt ...
  --- FAIL: TestRealRepoLedgerIsHonest
      scan_test.go:815: ban #6 examined 35 frontend/ text files in the real repo
      scan_test.go:819: HEAD must be green, rc=1 ... err=d22scan: 1 finding(s); D22 bans are not negotiable
  ```
  同一跑里 ban #6 仍是 `examined 35 text files`、`#1-5 internal/=190`（比翻牌那批多 3 个 .go）。
  根因：**票 89 的 `de15a6b`** 往 `internal/winsec/winsec.go:126` 放了 `filepath.Abs`，而 `allowlist.txt` 没有对应条目
  （`git show HEAD:internal/winsec/winsec.go | grep -n filepath.Abs` ⇒ 126 行确已在 **HEAD** 里，不是工作树噪声）。
  ⇒ **翻牌（本票）已经不是 push 的 blocker，票 89 才是**：要么票 89 那侧改走 C26 PathResolver，要么由**编排者**决定要不要加豁免。
  `allowlist.txt` 与 `internal/**` 都在本票禁改清单里，我一个字没动，也**不打算**动。
  另：在**工作树**（票 77/89 未提交的在飞改动）里跑包测试，报的是同样两条红 +
  `ban #6 examined 38 frontend/ text files`（38 是票 77 的代理正在往 `frontend/` 加文件，与本票无关）
  ⇒ 这就是为什么 AC#1/AC#5 的证据必须取自 `git archive` 的仓外纯净快照，而不是工作树。


