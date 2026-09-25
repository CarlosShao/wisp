# d22scan-gitignore-r1 — 扫描器跳过 git-ignore 路径（台账 `A207` 那一修）

**时刻**：2026-09-25 07:34–08:05 `+0800`（每节另标该节读数的时刻与锚点）
**锚点**：开工 `7771409` → 全部读数复算于 `467f8a4`（07:45，前端会话在那 30 分钟里推了 4 枚）→
交付三枚 commit：代码 `3bb99aa`（只带 `tools/d22scan/{main.go,gitignore.go,scan_test.go}`）·
本取证件 `364b9eb` · 这段抬头自己的更正 `00f9f65`（见 §7 末那条 ⚠，标题里的账号笔误按 `A207` 读）。
§8 那批复算用最终代码重跑，并附 `git archive 3bb99aa` 的 rc=0 读数。
**文中每个 `file:line` 都是在交付版上现 grep 重导的**，锚点变了就重导，不抄上一轮。
**作者身份**：实现者自证。**本文件不是裁决表**——`SPEC-12 §4.3` #1/#3 要求缺口审计与对抗
验收出自**另一个** agent，这一格请按 `A207` 原文验收，尤其验 §4 那发（它是这修的唯一"不是眼罩"凭据）。
**改动面**：`tools/d22scan/main.go` + 新文件 `tools/d22scan/gitignore.go` + `tools/d22scan/scan_test.go`。
**没碰**：`emojiRe` 字符类 · `allowlist.txt` · `frontend/**` · `design/**` · `internal/**` · 任何阈值/golden。
**快照纪律**：所有种子种在 `git archive | tar -x` 解到 `D:\tmp\qoder-d22scan-ignore-scratch\` 的仓库外副本里，
真树里没写过 probe 文件；临时件只建不删。

---

## 0. 进场测量①：扫描器在哪里决定"下不下潜"（锚点 `7771409`，改前行号）

`tools/d22scan/main.go` 四处独立写了同一套"跳过"判断，全部按 **目录名字面量** 硬编码，
**没有一处读过 `.gitignore`**（改前 `grep -n gitignore tools/d22scan/*.go` = 0 命中 ⇒ 票面前提"扫描器已在忽略某些东西"
只对了一半：它忽略名字，不忽略 git 的判断）：

| 位置（改前） | 判断 | 忽略了什么 |
|---|---|---|
| `walkGo` `main.go:594` | `d.Name() == "testdata" \|\| ".git"` → SkipDir | Go 生产码 walks（bans #1-5 / #7） |
| `walkText` `main.go:778` | `"testdata" \|\| "node_modules" \|\| ".git"` | ban #6（`frontend/`）、ban #7 |
| `walkEmoji` `main.go:849` + `:857` | `"node_modules" \|\| ".git"`，另 `goOnly\|\|everyFile` 时跳 `testdata` | ban #8 四个 scope |
| `checkRoot` `main.go:1081` | `"testdata" \|\| ".git"` | 头部那枚 `examined N production Go files` |

⇒ `dist` 不在任何一张名字表里，`frontend/dist/` 的构建产物因此被算进分母。

**票面前提核对（全部现跑，锚点 `7771409`）**：
- `.gitignore:22-23` = `frontend/dist/*` + `!frontend/dist/.gitkeep`；`frontend/.gitignore:12` = `dist/*`（`git check-ignore -v` 报的是这枚，深路径优先）。
- 工作树里 `frontend/dist/` 实存 4 枚文件：`.gitkeep`（受追踪）+ `index.html` + `assets/index-CEH-Pz8P.css` + `assets/index-UL9kYvYl.js`（后 3 枚被 ignore）⇒ **差的正是 3 枚**，票面成立。
- 净快照里 `find frontend/dist -type f` 只有 `.gitkeep` ⇒ CI 的 40 含 `.gitkeep`。**这条否决了"整目录剪掉"的省事实现**：剪掉 `dist/` 会把 CI 也改成 39。
- 结构性保证：`git ls-files -i -c --exclude-standard` **输出为空** ⇒ 本仓没有任何受追踪文件落在 ignore 规则里，
  所以"跳过 ignore 路径"在净快照上**不可能**让任何计数下降（§2 再实测一次）。
- ⚠ 一条与票面不同的事实（不改判据，只补一句）：`git archive HEAD` 的快照**没有 `.git` 目录**
  （实测 `ls snap/.git` → No such file or directory）。⇒ **不能**用 `git check-ignore` 实现这一修，
  否则读数会取决于"这台机器上这是不是一次 checkout"，正是本修的病。实现改为读 `.gitignore` 本身。

## 1. 进场测量②：同源副本清点（`A207②`/`U1`/`U2` 那一族病）

`grep -rn 'node_modules' --include='*.go' internal tools cmd | grep -v '_test.go:.*//'` 的命中 +
手工过 `dist` 一支，**"谁跳 `node_modules`/`dist`"这条规则在仓里有 6 份**：

| # | 副本（锚点 `7771409`） | 归属 | 本批处置 |
|---|---|---|---|
| ① | `tools/d22scan/main.go:778`（`walkText`） | 我（141 已结案） | **改为读 `.gitignore`** |
| ② | `tools/d22scan/main.go:849`（`walkEmoji`） | 我 | **同上** |
| ③ | `tools/d22scan/main.go:594` + `:1081`（`walkGo`/`checkRoot`） | 我 | **同上**（原来连 `node_modules` 都不跳，这次一并统一） |
| ④ | `tools/d22scan/scan_test.go:1180`（改前号；修后现号 `:1192`，`TestLedgerCountsMatchAnIndependentWalk` 的**独立对照走查**） | 我 | **改，但故意保留成第二份字面副本**（见下） |
| ⑤ | `internal/panel/frontend_hygiene_test.go:289`（子串跳 `/node_modules/` 与 `/dist/`） | **前端会话/owner 交出的地界** | **不动，只登记** |
| ⑥ | `internal/panel/composer_test.go:286`/`:434`/`:606`、`internal/risk/pathresolver_rewrite_account_test.go:223`（各自 case 表里列 `node_modules`/`dist`/`build`） | 同上（`internal/**` 测试） | **不动，只登记** |

**④ 为什么故意留成两份**：那枚函数是"仪器自己说 N 枚 ⇒ 第二台独立仪器也数出 N 枚"的证伪器。
如果它去调用同一个 `gitignore.go`，那么"一条 ignore 规则吞掉整个 scope 的文件"就会**被仪器数错、
又被自己的对照验为正确**。两份分开才有牙齿，代价是它会主动打脸：今后谁在 `.gitignore` 里新增一条
**能落到 `design/` `frontend/` `internal/` `cmd/` 之内**的规则，`TestLedgerCountsMatchAnIndependentWalk`
就会红，报错文案那句"the number in the self-report is wrong"要读成
"ignore 政策现在吃进被扫 scope 了 ⇒ 这是**基线变更**，按 `A207②` 必须同批重标台账各 scope 基线"。
⑤⑥ 属 owner 已交出的那棵树/别人的 blast radius（票面也点名"那 test 在另一枚 blast radius 里，不要改，登记即可"），
**本批一并不动**：它们跳的是同一族路径，语义上与本修同向（不会因本修变红——它们本来就已经在跳 `dist`）。

---

## 2. 证明 (iii)：CI 的形状零行为变更（锚点 `467f8a4`，07:52）

同一枚 `git archive 467f8a4` 解两棵仓库外快照，一棵用 HEAD 的原扫描器、一棵只把 `tools/d22scan/{main.go,gitignore.go,scan_test.go}` 换成修后的，
两棵各跑一次 CI 逐字同形的 `sh scripts/d22scan.sh`（两步都跑，`set -eu`）：

```
$ git archive 467f8a4 | tar -x -C /d/tmp/.../snap-base     # 修前代码
$ cd snap-base && sh scripts/d22scan.sh          -> rc=0
runtests.sh: OK - packages=[./...] top-level: PASS=24 FAIL=0 SKIP=0, === RUN=64, '[no tests to run]'=0
d22scan: examined 225 production Go files under internal/ and cmd/
d22scan: scope bans #1-5 internal/      examined 203 production Go files
d22scan: scope bans #1-5 cmd/           examined  22 production Go files
d22scan: scope ban #6 frontend/         examined  40 text files
d22scan: scope ban #7 internal/tools/   examined  18 production Go files
d22scan: scope ban #8 design/           examined  16 text files
d22scan: scope ban #8 frontend/         examined  40 text files
d22scan: scope ban #8 internal/         examined 405 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined  39 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=203, bans #1-5 cmd/=22,
ban #6 frontend/=40, ban #7 internal/tools/=18, ban #8 design/=16, ban #8 frontend/=40,
ban #8 internal/=405, ban #8 cmd/=39; ...

$ git archive 467f8a4 | tar -x -C /d/tmp/.../snap-final    # 换入修后三枚文件
$ cd snap-final && sh scripts/d22scan.sh          -> rc=0
runtests.sh: OK - packages=[./...] top-level: PASS=26 FAIL=0 SKIP=0, === RUN=66, '[no tests to run]'=0
（step 2 的九行 d22scan: 输出与上面逐字相同；diff 只有一行不同：路径里的快照目录名）

$ diff <(修前 step2 的 scope/examined/clean 行) <(修后 step2 同行) | sed 's|snap-base|SNAP|;s|snap-final|SNAP|'
（无输出）⇒ PROOF(iii): CI 形状逐字同读数
```

- 各 scope **一个数都没动**：`203 / 22 / 40 / 18 / 16 / 40 / 405 / 39`，`examined 225`，rc=0。与票面预期值逐枚相符，
  也与 `A209①`/`A210③` 在 `bb61dc5`/`471af50`/`1755903` 三枚锚上量的 `design/16 frontend/40 internal/405 cmd/39` 同值。
- 修后那一步**没打印** `skipped as git-ignored` 那行：`awk '/scan of/{f=1} f' | grep -c 'skipped as' = 0`
  ⇒ 净 checkout 里没有 ignore 件，仪器在该沉默时沉默（那行的 `grep -c`=1 出现在 step 1 里，是新测试自己的 `t.Logf`，不是真扫描）。
- step 1 四数 `PASS 24→26`、`RUN 64→66`、`FAIL=0`、`SKIP=0` ⇒ 新增两枚测试，没让任何一枚存量测试变 SKIP 或变红。

## 3. 证明 (i)：同一棵树、同一串字节，修前报、修后不报（07:47，`snap-mut`）

```
$ printf 'ready \xe2\x89\xa4\n' > snap-mut/frontend/dist/probe.txt      # U+2264，落在被 ignore 的路径
$ git check-ignore -v frontend/dist/probe.txt
frontend/.gitignore:12:dist/*	frontend/dist/probe.txt

【修前的扫描器（HEAD 代码）跑同一棵树】
$ bin/d22scan-before.exe -root snap-mut       -> rc=1
d22scan: scope ban #6 frontend/         examined  41 text files
d22scan: scope ban #8 frontend/         examined  41 text files
frontend/dist/probe.txt:1: [emoji] ban #8 glyph in scope frontend/ is banned (D23): ...

【修后的扫描器跑同一棵树、同一字节】
$ bin/fixed.exe -root snap-mut                -> rc=0
d22scan: skipped as git-ignored: 1 file(s) under 0 ignored director(ies) [], decided by frontend/.gitignore (1 path(s))
d22scan: scope ban #6 frontend/         examined  40 text files
d22scan: scope ban #8 frontend/         examined  40 text files
d22scan: clean - no D22 ban violations; ...
```

**红-前 / 静-后两读都贴了真输出。** 顺带一枚自捉的仪器缺陷（不贴读数就没人知道它存在过）：
第一次实现把 `note()` 的"跳过 N 枚"按 **walk 访问次数** 计，`frontend/` 被 ban #6 与 ban #8 各数一遍
⇒ 明明只排除 1 枚文件却打印"2 file(s)"，而计数只从 41 动到 40。`gitignore.go:105` 的 `counted` 集合（去重门在 `:141-149`）后
读数自洽（`1 file(s)`）。这条也是新增测试 `TestWalksSkipGitIgnoredPaths` 的第一发红（"want 4, got 5" + note 那句），
它先红后才可信。

## 4. 证明 (ii)：不是眼罩——受追踪形状的种子照样红（07:49，同一棵 `snap-mut` 上累加）

同一棵快照再种 4 枚**不被 ignore**的违规（`frontend/fixtures/`、`frontend/src/`、`design/screens/`、`internal/probe/`）：

```
$ bin/fixed.exe -root snap-mut     -> rc=1，5 findings
internal/probe/leak.go:7: [bare-goroutine] bare `go worker(...)` is banned (D22/D38b, ...)
frontend/src/probe.tsx:1: [panel-approval] `approval.decide` in frontend/ is banned (D33/F2: ...)
design/screens/probe.html:1: [emoji] ban #8 glyph in scope design/ is banned (D23): ...
frontend/fixtures/probe.txt:1: [emoji] ban #8 glyph in scope frontend/ is banned (D23): ...
internal/probe/leak.go:3: [emoji] ban #8 glyph in scope internal/ is banned (D23): ...
d22scan: skipped as git-ignored: 1 file(s) ...        ← 排除的仍只有 dist/probe.txt 那一枚
d22scan: 5 finding(s); D22 bans are not negotiable ...
分母同步随种子变大：#1-5 internal/ 203→204、#6 frontend/ 40→42、#8 design/ 16→17、#8 internal/ 405→406
```

**同树对照旧扫描器：rc=1、6 findings**，把两台的 finding 名单做差：

```
$ diff <(旧 6 条) <(新 5 条)
2d1
< frontend/dist/probe.txt:1: [emoji]
```

⇒ 这一修**只**让被 ignore 的那一枚闭嘴；三枚 ban（#1/#6/#8）在四个 scope 上全部照常红，
且被 ignore 路径与被追踪路径上的字节内容完全相同（同一串 `≤`），差别只在路径。

仓内同形的两枚单元测试（`tools/d22scan/scan_test.go`，锚点 `467f8a4` 现号）：
- `TestWalksSkipGitIgnoredPaths` `:1263`（fixture 里同时种 ignore 件与追踪件，分母断言为**字面量** 5/5/3/3，
  断言方向两头都钉：追踪件不报=测试红，ignore 件报了=测试红；另断 `ban #6 == ban #8` 同一整数与 `note()` 点名规则文件）。
- `TestGitIgnoreRuleSemantics` `:1358`（20 行逐路径表：`!dist/.gitkeep` 必须**不**被排除、
  `dist` 目录本身必须**不**被排除（否则 CI 40→39）、`dist/assets` 目录必须被剪、dirOnly 规则不得命中同名**文件**、
  `**` 跨段、`?` 不吃空、锚定规则不跨层、深枚 `.gitignore` 的否定胜过浅枚、`#` 行不是规则、
  nil matcher 与 walk 起点目录永不跳过）。

## 5. 证明 (iv)：工作树效应 43 → 40（07:53，真树只读；锚点 `467f8a4`）

```
【修前扫描器，真工作树（这台机器 09-21 跑过 npm run build）】
$ bin/d22scan-before.exe -root "D:/work/workspace/projects plans/Wisp"     rc=0
d22scan: scope ban #6 frontend/  examined  43 text files
d22scan: scope ban #8 frontend/  examined  43 text files

【修后扫描器，同一棵真树，未写入任何字节】
$ bin/fixed.exe -root "D:/work/workspace/projects plans/Wisp"              rc=0
d22scan: skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/], decided by frontend/.gitignore (2 path(s))
d22scan: scope ban #6 frontend/  examined  40 text files     ← = CI 的数
d22scan: scope ban #8 frontend/  examined  40 text files     ← = CI 的数
其余七数不变：203 / 22 / 18 / 32 / 405 / 39（design/ 见下）
```

⇒ `ban #6/#8 frontend/` 的工作树读数与净快照读数**现在是同一个 40**；
被排除的 2 枚路径 = `frontend/dist/index.html` + 被剪掉的 `frontend/dist/assets/`（其内 2 枚 bundle 随目录一起不出现在分母里）。

**一处必须说清的"没变"**：`ban #8 design/` 在这棵工作树上是 **32**、净快照是 **16**，本修**不动它，也不该动**。
现量（07:53）：`find design -type f` = 68、`git status --porcelain --ignored=matching design/` 的 `!!` 行 = **0**
⇒ 这 16 枚差额全是**未追踪但没被 ignore** 的件（前端会话把 16 枚受追踪 mockup 移出线、另建 `design/doubao/` `design/old/`）。
**"还没 `git add`" 不等于 "不是交付物的一部分"**：这类件正该被 ban 看见。本修的边界就是 `A207` 的边界——
只处理 git-ignore，不处理 git-untracked。

## 6. 证明 (v)：扫描器自己的测试套件（07:57，真树，仓内跑的正是 CI 那台仪器）

```
$ sh tools/d22scan/runtests.sh -C tools/d22scan ./...     -> rc=0
runtests.sh: OK - packages=[./...] top-level: PASS=26 FAIL=0 SKIP=0, === RUN=66, '[no tests to run]'=0
$ gofmt -l tools/d22scan/   ->（空）
$ go vet -C tools/d22scan ./...  ->（空）
```

四数：`PASS=26 / FAIL=0 / SKIP=0 / === RUN=66`（开工基线 `24/0/0/64`，见 §2 修前列）。
**没放宽任何断言、没加 SKIP、没加 allowlist 行**（`git diff` 里 `allowlist.txt` 不出现，`emojiRe` 一字未动）。

⚠ **`scripts/d22scan.sh` 是 `set -eu`、step 1 是正向对照**：§4 那种"往真树里种违规"的形态下 step 1 会先红
（`TestScannerSelfScanOfRealRepoIsGreen` 读的是运行时树），脚本于是在 step 2 之前退出 ⇒
**"step 2 没输出"绝不能读成"step 2 过了"**。本文件所有 §4/§3 的单步读数都是用二进制直接 `-root` 跑出来的，
就是为了不让这条顺序把两件事混成一件事。

## 7. 故意没动 / 已登记的过期引用

- `internal/panel/frontend_hygiene_test.go:289`、`internal/panel/composer_test.go:286/434/606`、
  `internal/risk/pathresolver_rewrite_account_test.go:223`：见 §1 表⑤⑥，别人地界，**语义同向、本修不需它们改**。
- `tools/d22scan/main.go:245` 与 `:516-518`（`emojiScope.everyFile` 与 `emojiScopes` 的注释里那句
  "isTextFile 丢掉 37 枚受追踪文件里的 5 枚"）：`git ls-files frontend` 现量 **40** ⇒ "37" 是过期测量。
  它论证的是**文件类别**（`.mjs`/无扩展名/dotfile 会被白名单漏掉），结论仍成立，只是总数老了；
  改它要重写一段与"ignore"无关的论证 ⇒ **不在这一修的射程，登记不改**。
- 台账 `docs/reports/pending-and-issues.md`：**未追加、未改**。`A207②` 写明"必须与 `Q-48`/`U1` 同批、
  在那一批里把台账各 scope 基线一次重标"，且 `A210 next=①` 把重标列为编排者那一批的事（`U1=3c80352` 已落地）。
  可直接粘贴的重标读数：**修后基线（锚点 `467f8a4`）= `#1-5 internal/=203、cmd/=22` · `#6 frontend/=40` ·
  `#7 internal/tools/=18` · `#8 design/=16（净快照）/32（本工作树，含未追踪件）` · `#8 frontend/=40` ·
  `internal/=405` · `cmd/=39` · `examined 225 production Go` · `runtests PASS=26 FAIL=0 SKIP=0 RUN=66` · rc=0**；
  旧"脏树 43"读数自本批起不再存在。
- 存量工单里那些 `frontend/=40` 与 `frontend/=43` 的读数（`.scratch/wisp/issues/105:166`、`111:186/287` 记 43；
  `101:83`、`103:149`、`106:153`、`110:86` 记 40）：**append-only，不改写**，本修让这两个数从此合流为 40。
- ⚠ 本取证件自己的第二枚 commit `00f9f65` 标题里把账号写成了 **`A2207`**（正确是 `A207`，仓里没有 `A2207` 这一枚）。
  标题不可改写（`AGENTS.md` §1.4：已提交历史只追加、不重排），所以在这里点名一次：
  **读那枚 commit 时按 `A207` 对号，不要以为另有一枚台账。** 内容本身没有因此改动。

## 8. 交付前的复算（08:03，锚点 `467f8a4` ⇒ 代码 commit `3bb99aa`）

§2 之后我又改了 `main.go` 的三处**注释/测量文字**（`35 files` 那两枚过期读数与一段 A207 说明），
所以本修的五发读数**全部用最终代码重跑一遍**，不拿旧读数交差：

```
$ go build -C tools/d22scan -o .../bin/final.exe . && gofmt -l tools/d22scan/  ->（空）
$ go vet   -C tools/d22scan ./...                                             ->（空）

(iii) 把最终三枚文件覆进 snap-final（仍从 467f8a4 archive）后重跑 CI 两拍：
$ cd snap-final && sh scripts/d22scan.sh                          rc=0
runtests.sh: OK - top-level: PASS=26 FAIL=0 SKIP=0, === RUN=66
$ diff <(修前 step2 全部 d22scan: 行) <(修后同批行)               -> 空输出 ⇒ 逐字相同
   （examined 225 / 203 / 22 / 40 / 18 / 16 / 40 / 405 / 39，clean 行同）

(i)+(ii) 同一棵 snap-mut（dist 种子 + 4 枚受追踪种子）用最终二进制重跑：
   旧 binary  -> rc=1，6 条 finding，ban #6 frontend/ = 43
   final      -> rc=1，5 条 finding，ban #6/#8 frontend/ = 42
                 d22scan: skipped as git-ignored: 1 file(s) under 0 ignored director(ies) [], decided by frontend/.gitignore (1 path(s))
   finding 名单差集：2d1  < frontend/dist/probe.txt:1: [emoji]     ← 只少了这一枚

(iv) 真工作树（只读）用最终二进制重跑：
   旧 -> ban #6/#8 frontend/ = 43        修后 -> = 40（另七数不变，design/ 仍 32）
   d22scan: skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/], decided by frontend/.gitignore (2 path(s))

(v) 真树里 CI 那台仪器重跑：
$ sh tools/d22scan/runtests.sh -C tools/d22scan ./...   rc=0
runtests.sh: OK - packages=[./...] top-level: PASS=26 FAIL=0 SKIP=0, === RUN=66, '[no tests to run]'=0
```

**交付物自己的净快照读数**（`git archive 3bb99aa` 解到仓库外再跑同一脚本，08:02）：
`rc=0`、`PASS=26 FAIL=0 SKIP=0 RUN=66`、`203/22/40/18/16/40/405/39`、`examined 225`。

**交付时的现号引用**（`3bb99aa`，非上一轮抄来的号）：
- 四处 walk 的 ignore 判断：`tools/d22scan/main.go:633`+`:644`（`walkGo`）、`:823`+`:835`（`walkText`，ban #6/#7）、
  `:901`+`:934`（`walkEmoji`，ban #8 四个 scope）、`:1139`–`:1151`（`checkRoot` 那枚头部读数）；
  matcher 在 `:209` 随 scanner 建。
- 新文件 `tools/d22scan/gitignore.go`（420 行）：`counted` 去重字段 `:105`，去重门 `:141-149`。
- 新测试：`tools/d22scan/scan_test.go:1263`（`TestWalksSkipGitIgnoredPaths`）、`:1358`（`TestGitIgnoreRuleSemantics`）；
  第二份字面副本 `ignoredLikeGit` `:1458`（`relToRepo` `:1433`），它接进独立对照走查的位置是 `:1192` 那个 `count()`。
- 未动的合同面：`emojiRe` 仍在 `tools/d22scan/main.go:149`，类内容与票面 quoted 的那串逐字相同；
  `git show --name-only 3bb99aa` 只列三枚 `tools/d22scan/*`，`allowlist.txt` 不在其中。

## 9. 一句总判

`tools/d22scan` 的分母不再踩在流动的地面上：四处 walk（`main.go:633/644`、`823/835`、`901/934`、`1139-1151`）
统一问 `gitignore.go`"git 会追踪这枚吗"，答案来自被扫树自己的 `.gitignore`（含 `!` 否定与"深枚优先"），
不是第二份硬编码名单，也不是 `git check-ignore`（快照里没有 `.git`，用它就等于把读数重新绑回机器状态）。
CI 形状逐字不变（`203/22/40/18/16/40/405/39`，rc=0，前后 diff 空），工作树 `frontend/` 43→40 与 CI 合流，
同一棵树同一串字节在 ignore 路径上静默、在四枚受追踪路径上仍 rc=1 报 5 条（差集恰为那枚 ignore 件），
套件 `PASS=26 FAIL=0 SKIP=0 RUN=66`。
