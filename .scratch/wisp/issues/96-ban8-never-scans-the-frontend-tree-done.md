# 96 — `ban #8`（零 emoji）**根本没扫 `frontend/`**：面板那 40 个文本文件今天是"门开着但没人看"

**Status:** **accepted-done**（2026-09-21 17:4x 独立对抗验收判 **AC#1–AC#6 六格全 PASS**，裁决表 `docs/evidence/s1/96-adversarial-acceptance.md`。要点：相等性**不是自证**——真正承重的是 `TestLedgerCountsMatchAnIndependentWalk` 的独立 walk（两门同时收窄时，相等用例仍绿、独立 walk 红 33 vs 37）；五类种子各自 rc=1 点名到行；三枚变异按预期。⚠ 验收带回两条残留：脚本第 1 步可被缓存回放归**票 99**；`ci.yml:25` 与 `main.go:30` 两处**抄来的清单已腐烂**归票 85 的台账段（见 A73④）。）
**AC#1..AC#6 六格全达成**（AC#6 起初因票 94 的 `winsec.go:126` 记 PARTIAL，该红源随 `7910bcd` 落地而消失，
读数与两次的差别写在 Progress log 末两条）（2026-09-21 15:5x 编排者建；来源=票 77 接续代理交件时**点名为它 AC#4 的硬缺口**，不是它要偷工）
**Type:** 门禁完整性（票 71 AC#4 / A44① / A54② / 票 88 的同族：**覆盖面自己会烂，而输出长得和"检查过"一模一样**）
**Blocks:** 票 **77 的 AC#4**（那一框字面就要"`ban #8` 对 `frontend/` 的自报文件数"，今天**打不出这个数**）
· **Blocked by:** nothing（`tools/d22scan/` 此刻无人写；票 88 已交件）
**Packages:** 只 `tools/d22scan/main.go`（`emojiScopes()` 那一处 + 必要的 walk 过滤器）与 `tools/d22scan/scan_test.go`。
              **禁改**：`allowlist.txt`（**5 行非注释，只许变短或持平**）、任何 ban 的**文本/正则/字符类**（D22：不是代理能缩的）、
              "live 作用域 examine 0 个文件 ⇒ 致命"这条规则本身、`frontend/**`（票 77 的地界）、`internal/**`、`cmd/**`。

## 实测（编排者自己读到源码，不是引用别人的话）

`tools/d22scan/main.go:451-457` 的 `emojiScopes()` 逐字是三条：

```go
{dir: filepath.Join(root, "design"),   label: "design/"},
{dir: filepath.Join(root, "internal"), label: "internal/", goOnly: true},
{dir: filepath.Join(root, "cmd"),      label: "cmd/",      goOnly: true},
```

而 `ban8Scopes()`（`:439-451`）就是把它原样搬进台账 ⇒ CI 自报的三行是
`ban #8 design/ 16 text files`、`ban #8 internal/ 324 Go files`、`ban #8 cmd/ 25 Go files`。
**`frontend/` 不在清单上** ⇒ 票 88 刚武装的 `ban #6` 覆盖 `frontend/` 的 40 个文本文件（票 77d 的实测数），
而**零 emoji 那条门对这些文件一个字节都没读过**。
⚠ 这不是"CI 没跑"：它跑了、绿了、自报了三个数，**只是那三个数里没有我们新长出来的那棵树**——
和票 67 当年把 `ban #6` 登记成豁免是同一族，只不过这次的形状是"**清单忘了同步**"，因此**连豁免都没写过**。

## 为什么这条值得单独一张票而不是"顺手加一行"

加一行 `frontend/` 会**当场把台账用例弄红**（票 88 刚演示过一遍：翻牌代价是 5 条重新表述）。
更要紧的是**判据本身**：`frontend/` 里最可能出现 emoji 的地方是**注释与 UI 文案**（`.tsx`/`.css`/`.md`），
而 `design/` 那条走的是"text files"模式（不要求 `.go`）⇒ **要用哪条路、`node_modules/` 与 `testdata/` 怎么排、
`VENDORED.md` 这种"上游原文里可能就有 emoji"的 vendored 文件怎么办**——这三个决定都会改变 CI 的颜色，
必须**写在 commit 正文里**，不许静默选一个。⚠ 尤其**不许**为了变绿把 vendored 文件整目录排除掉：
那是把 R18 的"三库都是 copy-paste 进仓"这条暴露面重新藏起来。

## AC（1:1，裁决表 `docs/evidence/s1/96-*.md` 由验收方出，不是自裁）

- [x] **AC#1** 台账里出现第四行 **`ban #8 frontend/ examined N text files`，N>0**，且**逐文件后缀构成**要贴出来
      （多少 `.tsx`/`.ts`/`.css`/`.md`/`.json`/无后缀）。**判据是"N 与票 88 量的 `ban #6` 那个数对得上或解释得清差在哪"**
      ——两条门扫同一棵树，数不同就必须给原因（排除规则不同？），**不许含糊**。
- [x] **AC#2** 阳性对照：在**仓外纯净快照**（`git archive <SHA> | tar -x -C /tmp/<带你的会话后缀>`）的
      `frontend/` 里种一个 emoji（**注释里放一个**，别只测字符串字面量），`sh scripts/d22scan.sh` 必须 rc=1 并点名它；
      ⚠ 再种一个只在 `.md` 里的 ⇒ 明写你的选择（`.md` 扫不扫）与理由，别让它靠运气命中。
- [x] **AC#3** 反向对照：**把 `frontend/` 从 `emojiScopes()` 里删掉**（模拟"清单又忘了同步"）⇒
      必须有用例红。若全绿，说明台账根本没在核对清单，那 AC#1 就是自证。
      锚点=承载行为那一行，**同链 grep 证落地再跑**，还原后 `git diff --quiet` 证干净，**编译失败不算变异**。
- [x] **AC#4** 因加作用域而红的台账用例**逐条重新表述**（票 88 刚做过一遍，照它的做法，**不是删用例**）：
      逐条列"旧断言 / 新断言 / 为什么新断言仍然在测同一件事"。
- [x] **AC#5** `ban` 文本零改动证明：`git diff` 到本票第一枚 commit 的父，贴出**ban 文本与字符类那几行未动**的证据；
      `allowlist.txt` 仍是 **5 行非注释**（`grep -v '^#' allowlist.txt | grep -c .`）。
- [x] **AC#6** 门禁：`cd tools/d22scan && go test -count=1 ./...` rc=0（**必须 `-count=1`**，
      ⚠ 它是**独立 Go module**，从仓根跑会打印 `main module does not contain package` 而**扫描器根本没执行**——本仓栽过）、
      `gofmt -l tools/d22scan/` 空、纯净树 `sh scripts/d22scan.sh` **rc=0** 且逐作用域贴全。

## Rules（本仓固定）

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc）；提交前 `git diff --cached --name-only`；
禁 `git add -A`/`.`、`--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34，共树）；**不 push**；
不在仓内建 worktree（A38④）；票面 append-only（改行前先读；标题前插段落要把标题重抄，`git diff --numstat` 删除列必须 0）；
四种假绿逐条点名；数字不达标写 FAIL 附数字。**收尾前必跑 `sh scripts/d22scan.sh`**（A64② 立的规矩）。
15 次工具调用内交回第一枚 checkpoint；每次提交同步 Status + 勾框 + `next=`；接近轮数上限主动收尾留断点。
⚠ **共树提醒**：`internal/winsec/`（票 94）、`internal/config/`+`internal/agent/`（票 90）有人在写，别碰；
`.scratch/wisp/issues/{82,86,87,88,89,90,91,92,93,94,95}-*.md` 是别人的票面。

## Progress log（append-only）

- 2026-09-21 15:5x（编排者）：建票。来源：票 77 的接续代理（`agent-ticket77d`）在交件时写明
  **它的 AC#4 只能记 PARTIAL**——"框字面要 `ban #8` 对 `frontend/` 的自报数，而它的声明作用域不含 `frontend/`，
  改 `tools/d22scan` 非本票地界"。**这个判断是对的，我按它说的单独立案**，不去把 77 的框圆过去。
  我自己读了 `emojiScopes()`（file:line 在上面）确认缺口是真的，**没有采信它的叙述**。
  ⚠ 与票 88 的区别要说清：那次是"目录不存在所以登记豁免"（有账、会漂移报错）；
  这次是**清单根本没同步过**，所以**连豁免都不存在** ⇒ 它比豁免更隐蔽：CI 自报三行、一行都不提 `frontend/`。
  next= 可立即派单（只碰 `tools/d22scan/`，与在飞的票 89/90/94 无文件交集）。

- 2026-09-21 16:2x（agent-ticket96）：**已交件**，commit `5e8f87b`（只含 `tools/d22scan/main.go`
  + `tools/d22scan/scan_test.go` 两个路径；`git diff --cached --name-only` 提交前为空，
  未用 `git add -A`，未碰任何 `internal/**`）。逐字复核了编排者的前提：`emojiScopes()` 确实只有
  design/ + internal/ + cmd/ 三条，**没有 frontend/，也没有登记过豁免**。

  **三个决定的落地（写进 commit 正文，摘要于此）**
  ① walk 路径：新增 `emojiScope.everyFile`，`frontend/` 走**全文件、无后缀白名单**，
     与 `ban #6` 的 `walkText` 逐字同形。**不选** `goOnly:false`（那条走 `isTextFile()`）：
     实测 `git ls-files frontend` = 37 个文件，`isTextFile()` 会丢 **5** 个
     —— `scripts/{vendor,vendor-shadcn,gen-tokens}.mjs` 全部三个、`.gitignore`、`dist/.gitkeep`，
     即 ban #8 对 `frontend/` 的覆盖面会变成 ban #6 的**真子集**。票 88 给 `ban #6` 加白名单是收窄，
     这边复用白名单是同一种收窄的反方向，所以没重复。`design/` 一行未动
     （实测其 16 个文件 12 .html + 2 .js + 2 .css，**0 个**落在 `isTextFile()` 外，动它是无谓改判据）。
  ② `node_modules/` 与 `testdata/`：照 `ban #6` 的 `walkText` 排除（连 `.git`），
     即把 `walkEmoji` 原有的 `goOnly ⇒ skip testdata` 扩到 `everyFile` 分支；
     两条门在同一棵树上的排除集因此**完全相同**。**没有排除任何子树**：
     `VENDORED.md` 与 vendored 的 `src/components/ui/*.tsx` 全在扫描范围内，
     实测这批上游原文今天 **0 命中**（perl 独立探针扫 40 文件、命中 0），故不需要藏。
  ③ 非 `.go` 文本（`.md`/`.json`/`.mjs`/无后缀）：**全扫**。由 ①② 决定，且 AC#2 四类种子逐个证明（见下）。

  **AC#1 台账（工作树，`go run . -root .` 逐字读数）**
  `d22scan: scope ban #8 frontend/         examined  40 text files` —— N=40>0，逐文件后缀构成（37 入库件）：
  17 `.tsx`、7 `.json`、3 `.ts`、3 `.mjs`、2 `.css`、1 `.md`、1 `.html`、1 `.go`、1 `.gitkeep`、1 `.gitignore`；
  工作树另加 3 个未入库构建产物（`dist/index.html`、`dist/assets/index-*.css`、`dist/assets/index-*.js`）= 40。
  **与 `ban #6` 的对应关系（不许含糊，故直说）**：两数**相等是构造要求而非巧合** ——
  `ban #6 frontend/ examined 40`、`ban #8 frontend/ examined 40`（工作树），
  `git archive HEAD` 纯净快照上两边都是 **37**；差 3 纯粹是 `frontend/dist/` 的构建产物未进 git、
  两条门同时看到。**不存在**排除规则差异。钉成两条用例：
  `TestRealRepoBan8CoversFrontendTreeAtBan6sCount` 直接断言 `emojiSeen["frontend/"] == examined["panel-approval"]`
  且 >0，并逐字核对自报那一行；`TestLedgerCountsMatchAnIndependentWalk` 用第三条独立 walk
  （accept-everything + 同样的三个 skip）同时核对这两个数。

  **AC#2 阳性对照（仓外纯净快照 `/tmp/wisp96-ac2` = `git archive 5e8f87b | tar -x`）**
  `.tsx` **注释**里种 U+2713 ⇒ rc=1，`frontend/src/App.tsx:1: [emoji] ban #8 glyph in scope frontend/`；
  只在 `.md` 里 ⇒ rc=1 点名 `frontend/README.md:3`；`.mjs` ⇒ `frontend/scripts/probe96.mjs:1`；
  无后缀 ⇒ `frontend/Procfile96:1`。四类都是 rc=1 且点名，**没有一类靠运气命中**。
  ⚠ 取证口径要交代：纯净快照上 `sh scripts/d22scan.sh` 整条是 **rc=1**，但**唯一命中是
  `internal/winsec/winsec.go:126`（票 94 的账，其修复此刻只在别人的工作树里、未提交）**，
  且脚本第一步 `go test` 因该命中 FAIL + `set -eu` ⇒ **走不到 scan 步**。为隔离本票的判据，
  我在**快照内**（不是仓内）往 `allowlist.txt` 追加了一条中和该已知命中，
  之后 `go run . -root /tmp/wisp96-ac2` 无种子时 **rc=0（clean）**，四类种子各自 rc=1。
  仓内文件的 allowlist 未被触碰（见 AC#5）。

  **AC#3 反向对照（两向变异，锚点=承载行为那一行，同链 grep 证落地，还原后 `git diff --quiet` 证干净）**
  (a) 从 `emojiScopes()` 删掉 `frontend/` 那一条（模拟"清单又忘了同步"）⇒ `go test -count=1 ./...` **rc=1**，
      红 **6** 条：`TestDeclaredEmojiScopeCannotWalkZeroFiles`、`TestScopeReportMatchesRealCoverage`、
      `TestRealRepoBan8CoversFrontendTreeAtBan6sCount`、`TestBan8FrontendScopeIsNotNarrowedByAnExtensionFilter`、
      `TestLedgerCountsMatchAnIndependentWalk`、`TestBuiltBinaryGoesRedEndToEnd` ⇒ **台账确实在核对清单**，
      AC#1 不是自证（关键在 `TestScopeReportMatchesRealCoverage`/`...AtBan6sCount` 里的 `"frontend/"` 是**字面量**，
      不是从 `emojiScopes()` 推出来的——推出来的那种断言当初全体绿着放行了这个缺口）。
  (b) 把 `everyFile: true` 摘掉（模拟"决定①选错、复用 `isTextFile()`"）⇒ rc=1，红 3 条：
      `TestRealRepoBan8CoversFrontendTreeAtBan6sCount`、`TestBan8FrontendScopeIsNotNarrowedByAnExtensionFilter`、
      `TestLedgerCountsMatchAnIndependentWalk`（独立 walk 37 vs 报告 32，即被白名单丢掉的 5 个）。
  两次还原后 `git diff --quiet -- tools/d22scan/main.go` 均无输出（编译均正常，不是"编译失败冒充变异"）。

  **AC#4 因加作用域而红的旧用例：3 条，逐条重新表述，未删任何用例**
  1. `TestScopeReportMatchesRealCoverage`：旧=`if strings.Contains(report,"frontend") {红}`（票 71 故意删条目时的正确话术）；
     新=`if !strings.Contains(report,"frontend/") {红}`。仍测同一件事：话术与实际覆盖面一致；覆盖面变了话术跟着翻。
  2. `TestDeclaredEmojiScopeCannotWalkZeroFiles`：旧断言"点名 cmd/"在新排序下会由 frontend/ 抢先 ⇒ 红。
     重述=先补种 `frontend/` 让 cmd/ 仍是唯一空作用域（**消掉排序运气**），再加第二相：删 `frontend/`、补 `cmd/`
     ⇒ 必须点名 `frontend/` 且 `verdict` rc=2、stderr 含 `ban #8 scope frontend/`。仍测"声明的作用域扫 0 文件⇒响亮失败"。
  3. `TestBuiltBinaryGoesRedEndToEnd` 的 `ban 6 tree gone while declared live exits 2` 子用例：旧 wantAll
     `{"ban #6 frontend/","empty instrument"}`；新 `{"ban #8 frontend/","examined 0 files"}`，改名
     `frontend tree gone...`。原因：同一棵树现由 ban #8 声明，guard 1 结构性优先于 guard 2 ⇒ 仍测"活树没了⇒停止出结论"，
     只是由哪个守卫出这句变了；为免"empty instrument"端到端牙齿被遮蔽，**另加**
     `ban 7 tree gone while declared live exits 2`（`internal/tools/` 不在 ban #8 清单里）钉住它。

  **AC#5 ban 文本零改动证明**：`git diff 5e8f87b^..5e8f87b --stat` = 只 `main.go` + `scan_test.go` 两文件
  （`allowlist.txt` 不在其中）；`main.go` 删掉的**非注释行总共 4 行**，全是 walk 的分支控制流
  （`if sc.goOnly && d.Name()=="testdata"`、`if sc.goOnly {`、`} else if !isTextFile(path) {`、`return nil`）；
  对 `emojiRe`/`banned (D23)`/`approval.decide` 的 diff 命中数 **0** ⇒ 字符类、正则、ban 措辞一字未动。
  `grep -v '^#' allowlist.txt | grep -c .` = **5**（持平，未改）。

  **AC#6 门禁读数**：`cd tools/d22scan && go test -count=1 -v ./...` ⇒ **rc=0**，
  `=== RUN` 31、`--- PASS`(顶层) 21（其余是子测试）、**FAIL 0、SKIP 0**（无 `--- SKIP`，故无"跳过冒充 ok"）；
  `gofmt -l tools/d22scan/` **空**。收尾必跑 `sh scripts/d22scan.sh`（工作树）⇒ **rc=0**，
  逐作用域八行已贴（bans #1-5 internal/=197、cmd/=20、ban #6 frontend/=40、ban #7=17、
  ban #8 design/=16、**frontend/=40**、internal/=335、cmd/=25）。
  ⚠ 两处如实交代、不圆场：(i) 该脚本第一步是 `go test ./...`（无 `-count=1`），本次显示 `(cached)`，
  所以我另跑了一次带 `-count=1` 的全量（上面那组数）；scan 步每次都是真跑。(ii) **本票框面写的
  "纯净树 rc=0" 在 HEAD 上打不到**：`git archive 5e8f87b` 快照 `sh scripts/d22scan.sh` = **rc=1**，
  唯一命中 `internal/winsec/winsec.go:126 [pathresolver-bypass]`（票 94 在修，我一行未碰），
  其余 7 个作用域与本票工作树读数一致 ⇒ 判据按"除该已知命中外无新命中"达成（快照内中和后 rc=0），
  AC#6 记 **PARTIAL**，勾框留空。

- 2026-09-21 16:3x（agent-ticket96）：**AC#6 的纯净树补正 —— 上面那框留空的原因已经消失，AC#6 整框勾上。**
  上一条写"纯净快照 rc=1、唯一命中 `internal/winsec/winsec.go:126`"是 `5e8f87b` 当时的真读数，**不撤不回**；
  但它是有时效的：票 94 的 `7910bcd fix(94,AC#1-AC#3)` 随后落了地，本票交件后我重测了一次
  `git archive HEAD | tar -x -C /tmp/wisp96-head && sh /tmp/wisp96-head/scripts/d22scan.sh`
  ⇒ **rc=0**，且那一步 `go test ./...` 是**真跑不是缓存**（`ok github.com/CarlosShao/wisp/tools/d22scan 8.550s`，
  缓存命中不会花 8.5 秒），零 finding，八行作用域逐字：
  `bans #1-5 internal/ 197`、`bans #1-5 cmd/ 20`、`ban #6 frontend/ 37`、`ban #7 internal/tools/ 17`、
  `ban #8 design/ 16`、**`ban #8 frontend/ 37`**、`ban #8 internal/ 335`、`ban #8 cmd/ 25`。
  ⇒ AC#1 的"两边相等"在**已提交的树**上再次成立（37=37），不依赖我工作树里那些未入库的构建产物。
  全程未碰 `internal/winsec/`（上一条已写明该账不属本票），也未因它绿了就回头改判据。
  next= 只剩验收方出 `docs/evidence/s1/96-*.md` 裁决表；票 77 的 AC#4 自报数可解阻塞。
