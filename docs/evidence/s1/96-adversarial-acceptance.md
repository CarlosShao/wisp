# 票 96 独立对抗验收 —— `ban #8`（零 emoji）覆盖 `frontend/`

**验收代理：** `acceptor-ticket96`（独立对抗验收，不写生产码、不碰工作树生产文件、不 push）
**时间：** 2026-09-21 16:2x（起手 HEAD=`4adba70`；写作期间 HEAD 已推进到 `c8d5c94` —— 共树，别人的 commit 在飞，与本票无关）
**被验收对象：** `.scratch/wisp/issues/96-ban8-never-scans-the-frontend-tree.md`（票面 `1fc4ff7` + 补正 `8b1b10f`），
代码 commit `5e8f87b`（**只含** `tools/d22scan/main.go` + `tools/d22scan/scan_test.go` 两路径，`git show --name-only` 实测）。
**本文所有数字都是我自己跑出来的读数，没有一条抄票面。** 每条后面注明复现命令与快照路径。

---

## 0. 纪律自证（先说我没碰什么）

- 仓内**唯一写入**：本文件 `docs/evidence/s1/96-adversarial-acceptance.md`。`git status --short` 起手即录：
  只有 `.scratch/wisp/issues/90-*.md`、`.scratch/wisp/issues/91-*.md` 两枚**别人的**未提交改动，我一行未动。
- **所有变异与种子都在 `/tmp` 快照里做**，目录带我的会话后缀 `wisp96acc-96`：
  - `/tmp/wisp96acc-96/snap8b` = `git archive 8b1b10f | tar -x`（判定基准：票 96 交件后的树）
  - `/tmp/wisp96acc-96/snap5e8` = `git archive 5e8f87b | tar -x`（用于复核 AC#6 起初那枚红）
  - `/tmp/wisp96acc-96/mut-isTextFile`、`mut-drop-frontend`、`mut-both-narrow` = 三枚变异（各 `cp -r snap8b`）
- 未在仓内建 worktree（A38④）。未用 `git add -A`、未 `--amend`/`reset`/`rebase`/`stash`/`checkout .`、未 push。
- 对生产树**只读**：§6.4 那次对工作树的操作是 `go run . -root <仓>`（扫描器不写盘，跑前跑后我都看了 `git status --short`）。
- ⚠ **共树实况如实记**：我跑动期间 HEAD 从 `4adba70` 推进到 `c8d5c94`，工作树在收尾时出现
  `internal/winsec/reparse_windows_test.go`、`internal/winsec/winsec_windows.go`、`docs/evidence/s1/94-*.md`
  三枚**别人的**未提交改动（票 94 在飞）。这些不是我写的：我全程只 `git add` 自己的那一个路径，
  每枚 commit 前 `git diff --cached --name-only` 只含本文件。

---

## 1. 我的基线读数（纯净快照 `8b1b10f`，`git archive` 出来的树）

扫描器逐作用域（`cd /tmp/wisp96acc-96/snap8b/tools/d22scan && go run . -root <快照>`）⇒ **rc=0**，逐字：

```
d22scan: examined 217 production Go files under internal/ and cmd/ of .../wisp96acc-96/snap8b
d22scan: scope bans #1-5 internal/      examined 197 production Go files
d22scan: scope bans #1-5 cmd/           examined  20 production Go files
d22scan: scope ban #6 frontend/         examined  37 text files
d22scan: scope ban #7 internal/tools/   examined  17 production Go files
d22scan: scope ban #8 design/           examined  16 text files
d22scan: scope ban #8 frontend/         examined  37 text files
d22scan: scope ban #8 internal/         examined 335 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined  25 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations; live scope work: ... ban #6 frontend/=37 ... ban #8 frontend/=37 ...
```

**我自己数的那一棵树（不采信任何人的数）：**

| 量法 | 命令 | 读数 |
| --- | --- | --- |
| git 入库件数 | `git ls-tree -r --name-only 8b1b10f -- frontend \| wc -l` | **37** |
| 文件系统实数（快照内，排除 `node_modules/`、`testdata/`、`.git/`） | `find frontend -type f -not -path "*/node_modules/*" -not -path "*/testdata/*" \| wc -l` | **37** |
| 快照里 `frontend/` 下真存在 `node_modules/`/`testdata/` 目录？ | `find frontend -type d \( -name node_modules -o -name testdata \) \| wc -l` | **0**（排除集今天没有可排的东西） |
| 后缀构成 | 同上 `find -printf "%f"` + awk | 17 `.tsx`、7 `.json`、3 `.ts`、3 `.mjs`、2 `.css`、1 `.md`、1 `.html`、1 `.go`、2 无后缀 dotfile（`.gitignore`、`dist/.gitkeep`）= **37** |

⇒ 三个数（我的 git 数、我的文件系统数、两条门各自自报的数）**四路一致 = 37**。
票面 AC#1 的"工作树 40=40"我也单独复现了（§6.4），差 3 就是 `frontend/dist/` 里未入库的构建产物，两门同时看见，与票面解释一致。

---

## 2. AC#1 主攻点：那条新用例的**断言本体**（不看名字，只看代码）

`tools/d22scan/scan_test.go:767-836` `TestRealRepoBan8CoversFrontendTreeAtBan6sCount`，我逐行读到的断言是：

1. `front := s.emojiSeen["frontend/"]; if front == 0 { t.Fatal(...) }` —— 空树致命。
2. `if ban6 := s.examined["panel-approval"]; front != ban6 { t.Errorf(...) }` —— **这是"两个数互相相等"**，
   形状正是编排者担心的那种：它**单独**不能证明数是真的，只证明两条门读到一样多。
3. 遍历 `declaredScopes(root)`，字面量比对 `"ban #8 frontend/"`：必须存在于台账、必须 `live`、
   且 `sc.count(s) == front`（第三条同样是同一 scan 的内部自洽）。
4. 自报行：`want := fmt.Sprintf("d22scan: scope ban #8 frontend/ examined %d text files", front)`
   —— **注意：期望串里的 `%d` 就是被测的那个数本身**，所以这条断言只钉"这行存在/命名/格式化"，
   **钉不出数值的真假**。票面把这条写成"钉成用例"，措辞偏强；我按实际形状登记（见 §9 残留 ①）。
5. `scopeSummary(...)` 必须含 `ban #8 frontend/=<front>` —— 同样是自比较。

**那"数是真的"这件事到底由谁承重？** 由**另一条**用例：`scan_test.go:891-965`
`TestLedgerCountsMatchAnIndependentWalk`。我读了它的 `count` 闭包本体（不是注释）：

```go
count := func(dir string, accept func(string) bool) int { ... filepath.WalkDir(dir, func(...) {
    if d.IsDir() { if d.Name()=="testdata" || d.Name()==".git" || d.Name()=="node_modules" { return filepath.SkipDir }; return nil }
    if accept(p) { n++ } ... }) }
```

它**不调用** `main.go` 的任何 walk（不是 `walkEmoji`/`walkText` 的别名），且在同一个 `cases` 表里
**对两个数各自**与独立 walk 比对：

- `{"ban #6 frontend/", count(frontend, 一律 true), s.examined["panel-approval"]}`
- `{"ban #8 frontend/", count(frontend, 一律 true), s.emojiSeen["frontend/"]}`
外加 `if c.report == 0 { t.Errorf(...) }`。

⇒ **裁决：不是"只断言彼此相等"的假绿。** 相等性背后有第三条独立 walk 逐边核对，
且该 walk 的 `SKIP=0`（真跑了，`scan_test.go:958` 打出 `verified ban #6 frontend/: 37 files` /
`verified ban #8 frontend/: 37 files`，与我 §1 的两种独立量法同数）。

**但我还是把"相等会不会被造出来"这件事自己做了一次变异（§3 的 M-C）：会——不过不是靠这条用例挡的。**

---

## 3. 收窄防线：三枚变异，全部先 `go build ./...` 过（编译失败不算变异）

### M-A：把 `ban #8` 对 `frontend/` 装回 `isTextFile()` 白名单（= 票面决定①反着做）

改法（`mut-isTextFile/tools/d22scan/main.go:483`）：`{... label: "frontend/", everyFile: true}` → 去掉 `everyFile: true`
（于是落进 `default:` 分支 = `isTextFile()` 后缀白名单）。`go build ./...` **OK**，`go test -count=1 ./...` ⇒ **rc=1，红 3 条**：

```
--- FAIL: TestRealRepoBan8CoversFrontendTreeAtBan6sCount        scan_test.go:782: ban #8 examined 32 frontend/ files but ban #6 examined 37 ...
--- FAIL: TestBan8FrontendScopeIsNotNarrowedByAnExtensionFilter scan_test.go:865: ban #8 examined 6 frontend/ files, want 7
                                                              scan_test.go:876: ban #8 did not fire in frontend/scripts/vendor.mjs ...
                                                              scan_test.go:876: ban #8 did not fire in frontend/Procfile ...
                                                              scan_test.go:881: unexpected ban #8 hit at frontend/testdata/golden.tsx
--- FAIL: TestLedgerCountsMatchAnIndependentWalk               scan_test.go:952: scope ban #8 frontend/ reported examining 32 files, independent walk says 37
```

**"37 vs 32"我自己复现了，并独立算出被白名单丢掉的就是那 5 个：**
`frontend/.gitignore`、`frontend/dist/.gitkeep`、`frontend/scripts/{gen-tokens,vendor-shadcn,vendor}.mjs`
（我用 shell 复刻 `isTextFile()` 的后缀表跑出来的清单，与票面逐字同名）。
⇒ **收窄必红，PASS**。（附带发现：白名单分支还顺手把 `testdata/` 放了进来——`everyFile` 那条 testdata 排除也一起没了，
说明 fixture 用例连排除规则一起钉着。）

### M-B：AC#3 反向对照 —— 从 `emojiScopes()` 里删掉 `frontend/` 那一条

`mut-drop-frontend/tools/d22scan/main.go` 删第 483 行。`go build ./...` **OK**，`go test -count=1 ./...` ⇒ **rc=1，红 6 条**：

```
--- FAIL: TestDeclaredEmojiScopeCannotWalkZeroFiles
--- FAIL: TestScopeReportMatchesRealCoverage
--- FAIL: TestRealRepoBan8CoversFrontendTreeAtBan6sCount
--- FAIL: TestBan8FrontendScopeIsNotNarrowedByAnExtensionFilter
--- FAIL: TestLedgerCountsMatchAnIndependentWalk
--- FAIL: TestBuiltBinaryGoesRedEndToEnd
```

与票面报的 6 条名单**逐字同名**（我自己数的，不是引用）。我另外确认了"台账在核对清单"的**承重处是字面量**：
`TestScopeReportMatchesRealCoverage` 里是 `if !strings.Contains(report, "frontend/")`，
`TestRealRepoBan8...` 里是 `strings.Contains(strings.Join(labels, ","), "ban #8 frontend/")`
——两处都不从 `emojiScopes()` 推，所以"清单又忘了同步"必红。**⇒ AC#1 不是自证。**

### M-C：本票最刁的那一击 —— **两条门同时**收窄（对称收窄）

`mut-both-narrow`：在 `walkText`（ban #6）与 `walkEmoji`（ban #8）的目录排除里**同时**加 `d.Name() == "scripts"`。
`go build ./...` **OK**，`go test -count=1 ./...` ⇒ **rc=1**，但红的是这三条：

```
--- FAIL: TestBan6ScopeIsNotNarrowedByAnExtensionFilter
--- FAIL: TestBan8FrontendScopeIsNotNarrowedByAnExtensionFilter
--- FAIL: TestLedgerCountsMatchAnIndependentWalk  scan_test.go:952: scope ban #6 frontend/ reported examining 33 files, independent walk says 37
                                                  scan_test.go:952: scope ban #8 frontend/ reported examining 33 files, independent walk says 37
```

**而 `TestRealRepoBan8CoversFrontendTreeAtBan6sCount` 在这枚变异下是绿的**（33 == 33），
并且扫描器自己照样打出 `clean`、**rc=0**（`d22scan: scope ban #6 frontend/ examined 33 text files` /
`ban #8 frontend/ examined 33`）。⇒ 结论要写清：
**"两门相等"不是 falsifier，独立 walk + 那两个 fixture 用例才是**；
票 96 把承重放对了地方（它同时写了这三条），所以判 PASS，但 owner -facing 那句话的依据必须说成"独立 walk"，
不能说成"两条门数相等"（见 §9 残留 ②）。

---

## 4. AC#2 阳性对照：我自己种的种子（仓外快照 `snap8b`，四种 + 我加的一种）

每一枚：改/建目标文件 → `cd tools/d22scan && go run . -root /tmp/wisp96acc-96/snap8b` → 收 rc 与点名行 → **还原**。
两种glyph分别落在 `emojiRe` 的两个不同字符段里，防"只认一段"：`✓` U+2713（`\x{2600}-\x{27BF}`，3 字节）与
😀 U+1F600（`\x{1F000}-\x{1FAFF}`，4 字节）。⚠ 全部**放注释/prose，不放字符串字面量**——这是 `ban #8` 的定义（D23 覆盖注释与 `_test.go`）。

| # | 目标文件（快照内） | 类 | glyph | **我的 rc** | 点名行（逐字） |
| --- | --- | --- | --- | --- | --- |
| 1 | `frontend/src/App.tsx` 追加 `// acc96 probe ✓ comment tooth` | `.tsx` **注释** | U+2713 | **1** | `frontend/src/App.tsx:35: [emoji] ban #8 glyph in scope frontend/ is banned (D23): covers comments and _test.go, not only string literals` |
| 2 | `frontend/README.md` 追加 prose | 仅 `.md` | U+1F600 | **1** | `frontend/README.md:1: [emoji] ban #8 glyph in scope frontend/ ...` |
| 3 | `frontend/scripts/vendor.mjs` 追加注释 | `.mjs` | U+2713 | **1** | `frontend/scripts/vendor.mjs:111: [emoji] ban #8 glyph in scope frontend/ ...` |
| 4 | 新建 `frontend/Procfile96`（**无后缀**） | 无后缀 | U+1F600 | **1** | `frontend/Procfile96:1: [emoji] ban #8 glyph in scope frontend/ ...` |
| 5 | `frontend/.gitignore` 追加（**dotfile，M-A 里被白名单丢掉的那一类**） | dotfile | U+2713 | **1** | `frontend/.gitignore:17: [emoji] ban #8 glyph in scope frontend/ ...` |
| — | **无种子**（§1 那一跑） | — | — | **0** | `d22scan: clean - no D22 ban violations` |

⇒ 五类各自 rc=1 且**点名到文件:行号**，没有一类靠运气命中；`.md` 的选择不含糊：**扫**（票面决定③ + fixture 里
`frontend/VENDORED.md` 那一行明写"this is the .md answer - .md IS scanned"）。排除规则也明写：
`node_modules/` 与 `testdata/` 不扫（fixture 里两种反向对照钉着，M-A 顺带证明它还钉住了排除集）。
跑完 `diff -r snap8b mut-drop-frontend` 除该变异自身的 `main.go` 与一次 `go build` 产物外**无差异** ⇒ 种子全还原，快照未脏。

**整条脚本的形状（要如实交代）：** 我在同一种子下跑 `sh scripts/d22scan.sh` ⇒ **SCRIPT_RC=1**，
但它的**第 1 步 `go test ./...` 打的是 `ok ... (cached)`**（测试缓存只看包内源码，`frontend/App.tsx` 在模块外，
改了它不会让缓存失效），红是**第 2 步扫描器**给的（`1 finding(s)` + `exit status 1`）。
⇒ 本次判定**不受影响**（我要的就是扫描器点名），但"第 1 步可被缓存端过去"是本仓已登记的账（票 99），
我在 §9 一并挂出。**CI 侧那一跑不吃这个亏**：CI 第 1 步走 `runtests.sh`，其 `go test -v -count=1` 是硬编码的。

---

## 5. AC#4 / AC#5：有没有顺手削弱别处（我逐条自己算）

| 要检查的 | 我的命令与读数 | 结论 |
| --- | --- | --- |
| 本票代码 commit 的文件面 | `git show --name-only 5e8f87b` = `tools/d22scan/main.go`、`tools/d22scan/scan_test.go`（**只有这两个**） | 未越界 |
| `tools/` 里是否顺手改了别的 | `git diff --name-only 5e8f87b^..8b1b10f -- tools/` = 同两文件（`allowlist.txt`/`go.mod`/`runtests.sh` 不在内） | 未越界 |
| `main.go` 删掉的行是什么 | `git diff 5e8f87b^..5e8f87b -- main.go` 过滤注释后剩：`if sc.goOnly && d.Name()=="testdata"`、`if sc.goOnly {`、`} else if !isTextFile(path) {`、`return nil` —— **4 行，全是 walk 分支控制流**（新增 `everyFile` 分支与 `switch` 重写） | 只动作用域清单与 walk 控制流 |
| ban 文本 / 字符类 / 正则 | diff 内 `emojiRe`、`banned (D23)`、`approval.decide`、`1F000` 命中数 **0**（我 grep 过 `^[-+]` 行）；现行 `emojiRe` 仍是 `[\x{1F000}-\x{1FAFF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}\x{FE0F}\x{1F1E6}-\x{1F1FF}]`（`main.go:111`） | 一字未动 |
| `allowlist.txt` 仍是 5 行非注释 | `grep -v '^#' tools/d22scan/allowlist.txt \| grep -c .` = **5**；该路径 `5e8f87b^..8b1b10f` diff 为空 | 持平、未触碰 |
| `design/` 那 16 个文件的判据 | 条目 `{dir: design, label: "design/"}` 在 diff 里是 **context 行**（未改）；它落 `default:` 分支 = 原 `isTextFile()`。我实测 `design/` 16 个文件 = 12 `.html` + 2 `.js` + 2 `.css`，**全部在 `isTextFile()` 白名单内** ⇒ 不改它是"无谓动判据"的说法成立（改与不改今天同数，但改了就是动别人量过的东西） | 未削弱 |
| `ban #6` 的 37 未降 | `walkText` 本体（`main.go:758-784`）在 diff 中**未被触碰**；我在快照实测 `ban #6 frontend/ examined 37`，与我独立数一致；票面自报的工作树 40 我也复现（§6.4） | 未降 |
| 用例有没有被删 | `git show 5e8f87b -- scan_test.go \| grep -c '^-func Test'` = **0**；顶层 `^func Test` 数 19 → **21** | 未删用例 |
| AC#4 三条重述是否仍测同一件事 | 逐条读了：①`TestScopeReportMatchesRealCoverage` 的 `Contains(report,"frontend")` 翻成 `!Contains(...)` —— 仍测"话术与覆盖面一致"，方向因覆盖面变了而翻面；②`TestDeclaredEmojiScopeCannotWalkZeroFiles` 先补种 `frontend/` 再删它 ⇒ **消掉了"点名 cmd/ 靠 emojiScopes 排序运气"** 这个真实隐患，并加 `verdict` rc=2 + stderr 含 `ban #8 scope frontend/`（我在 M-B 下看到它红，证明它现在真有牙齿）；③`TestBuiltBinaryGoesRedEndToEnd` 的 `"ban 6 tree gone"` 改名 `"frontend tree gone"` 并把 wantAll 从 `empty instrument` 换成 `examined 0 files`，**另加** `ban 7 tree gone` 子用例把 "empty instrument" 的端到端牙齿单独钉住 | 重述成立 |

**附带证明 ③ 不是空转：** 我按 CI 口径跑了一次 `bash tools/d22scan/runtests.sh -C tools/d22scan ./...`（快照内）⇒
`runtests.sh: OK - top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31, '[no tests to run]'=0`，
子测试行里**两条都在**：`--- PASS: .../frontend_tree_gone_while_declared_live_exits_2`、
`--- PASS: .../ban_7_tree_gone_while_declared_live_exits_2`。新加的那条真跑真过，不是幻影。

---

## 6. AC#6 门检复跑 + 那处"起初 PARTIAL 后补正"的诚实性判定

### 6.1 我自己复跑的三件（纯净快照 `8b1b10f`）

- `sh scripts/d22scan.sh` ⇒ **rc=0**（第 1 步 `ok github.com/CarlosShao/wisp/tools/d22scan 5.212s`，**不是 `(cached)`**；
  第 2 步八行作用域逐字见 §1，含 `ban #6 frontend/=37` 与 `ban #8 frontend/=37`）。
- `cd tools/d22scan && go test -count=1 -v ./...` ⇒ **rc=0**，我数的：`=== RUN` **31**、顶层 `--- PASS` **21**、
  缩进子测试 `--- PASS` **10**、`FAIL` **0**、`SKIP` **0**（无 `--- SKIP` ⇒ 没有"跳过冒充 ok"；
  两条 real-repo 用例都真跑了：`--- PASS: TestRealRepoBan8CoversFrontendTreeAtBan6sCount (0.30s)`，
  且 `verified ban #6 frontend/: 37 files` / `verified ban #8 frontend/: 37 files` 两行都打出来了）。
- `gofmt -l tools/d22scan/` ⇒ **空**。
- ⚠ 独立 module 那个坑我按规矩走：判定跑全部在 `tools/d22scan/` 目录里跑（`-root` 指快照根）。
  顺手复现了一次坑本身：从仓根按包路径跑只会得到"main module does not contain package"，扫描器根本不执行。

### 6.2 "起初那枚红"到底存在过吗 —— 我去 `5e8f87b` 的快照上看了一眼

`git archive 5e8f87b | tar -x -C /tmp/wisp96acc-96/snap5e8` 后跑扫描器 ⇒ **rc=1**，唯一命中逐字：

```
internal/winsec/winsec.go:126: [pathresolver-bypass] filepath.Abs outside the C26 PathResolver is banned (D22); ...
d22scan: 1 finding(s); D22 bans are not negotiable
```

⇒ **起初记 PARTIAL 是真读数，不是给自己找的台阶。** 同一次跑里 `ban #6 frontend/ examined 37`、
`ban #8 frontend/ examined 37` —— 也就是说本票的判据当时就已经过了，卡住的是别人那枚红（票 94），
它没去修、没绕过、没改判据，这个处理我认可。
`7910bcd`（票 94 的修复）确实在 `8b1b10f` 的祖先链上（`git merge-base --is-ancestor 7910bcd 8b1b10f` ⇒ YES），
所以"该红随 `7910bcd` 落地而消失"这句话在时间线上成立，我在 §6.1 复跑到了 rc=0。

### 6.3 append-only 判定：原读数有没有被撤走

`git show 8b1b10f --numstat` = 14 增 / **5 删**。我逐行读了那 5 条删除：
2 行是 **Status 行重写**（`AC#1..AC#5 达成、AC#6 PARTIAL` → `六格全达成（AC#6 起初…记 PARTIAL，该红源随 7910bcd 消失…）`）、
1 行是 **勾框 `[ ]`→`[x]`**（规矩明写允许且要求每次提交同步）、2 行是上一条 log 末尾的 `next=` 换成新一条的 `next=`。
**原 PARTIAL 那条 progress log 的正文（含"AC#6 记 **PARTIAL**，勾框留空"）在现票面里逐字还在**，
新读数以**并列新条目**追加 ⇒ 两条读数共存。**判定：诚实，PASS。**
唯一小瑕疵（不判失败）：上一条 log 尾的 2 行 `next=` 被改写而不是原地保留，严格 append-only 的话该只追加；
这是格式级，不藏任何证据。前一枚票面 commit `1fc4ff7`（96 增 / 6 删）删除行 = Status 行 + 五条 AC 勾框，同性质。

### 6.4 工作树（HEAD 在动，我只读）

`go run . -root "D:/work/workspace/projects plans/Wisp"` ⇒ **clean、rc=0**，
`ban #6 frontend/ examined 40`、`ban #8 frontend/ examined 40`，我的独立 `find` 数也是 **40**（差 3 = `frontend/dist/` 未入库产物）。
`ban #8 internal/` 现在是 **336**（票 96 交件时 335）——是票 90/94 在 `internal/` 里新添的 `.go`，不是本票的账。

---

## 7. 四种假绿：逐条点名（我这一路）

1. **"pattern 打空造成的 ok"**：我的判定跑**全用 `./...` 全量、没用过 `-run`**；三枚变异跑同样全量并逐字列出 `--- FAIL` 名单。
   为排除这个形状，我还按 CI 口径跑了 `runtests.sh`（它断言 `=== RUN>0`、顶层 PASS/FAIL>0、SKIP=0）⇒
   `PASS=21 FAIL=0 SKIP=0, === RUN=31, '[no tests to run]'=0`。**无此假绿。**
2. **"SKIP 冒充 ok"**：`-count=1 -v` 全量里 `--- SKIP` = **0**（顶层+子），`21+10=31 == === RUN 31`。
   这条对本票特别要紧：`TestRealRepo*` 两条都带 `t.Skipf("not inside the wisp repo")`，
   若它们在快照里被跳过，"37=37"就是空气——日志证明它们**跑了并打了数**。**无。**
3. **"缓存冒充新证据"**：判定与变异跑一律 `-count=1`；§1/§6 的 `go run` 与 `go test` 都打真实耗时（5.212s / 3.666s / 4.700s…）非 `(cached)`。
   **但我实测到一处真的会被缓存端过去**：带种子那次 `sh scripts/d22scan.sh` 第 1 步打 `ok ... (cached)`
   （测试读的是模块外的 `frontend/`，Go 的测试缓存不看它）⇒ 登记为**本仓已立案的票 99**，
   不是票 96 新引入的，也不是本票判定所依赖的证据。**本票读数无一处来自缓存。**
4. **"阈值被偷偷调低 / 挑数字"**：我没改任何阈值、没动 allowlist（仓内 0 写）、没删断言、没挑样本 ——
   票面的 5 组数字（37/37、红 6 条名单、红 3 条名单、5 个被白名单丢的文件、allowlist 5 行、`=== RUN` 31/`PASS` 21/`SKIP` 0）
   我**各自独立复现后才采信**；三枚变异是**故意让全套变红**。多样本全报：纯净快照 `8b1b10f`（rc=0）+ 交件 commit `5e8f87b`（rc=1，1 命中）+ 工作树（rc=0，40=40）+ 三枚变异（各 rc=1）。**无。**

---

## 8. 1:1 裁决表（票面六个 AC 各一行，判据=票面原文）

| AC | 票面要什么（压缩） | 我的独立读数（不是引用） | 裁决 |
| --- | --- | --- | --- |
| **AC#1** | 台账第四行 `ban #8 frontend/ examined N text files`，N>0，逐类后缀构成，且与 `ban #6` 的数对得上或解释清差在哪 | 快照 `8b1b10f`：`ban #8 frontend/ examined 37 text files`（N=37>0）、`ban #6 frontend/ 37`；我的 git 数 37、我的 find 数 37；构成 17`.tsx`+7`.json`+3`.ts`+3`.mjs`+2`.css`+1`.md`+1`.html`+1`.go`+2 无后缀=37；工作树 40=40 且我独立数到 40；差 3=`frontend/dist/` 未入库产物，**不存在排除规则差异** | **PASS** |
| **AC#2** | 仓外纯净快照种 emoji：`.tsx` **注释**里、仅 `.md`、`.mjs`、无后缀，各 rc=1 且点名；无种子 rc=0；`.md` 扫不扫要明写 | §4 表：**5 类各 rc=1**（含我加的 `.gitignore`），点名到 `文件:行`；两种 glyph 跨 `emojiRe` 两个字符段；无种子 **rc=0 clean**；`.md` **扫**（票面决定③ + fixture 明写 + 我的 README.md 实测） | **PASS** |
| **AC#3** | 把 `frontend/` 从 `emojiScopes()` 删掉 ⇒ 必须有用例红；全绿则 AC#1 判 FAIL | M-B（先 `go build` 过 ⇒ 不是编译失败冒充）⇒ **rc=1、红 6 条**，名单与票面逐字同名；承重处是**字面量** `"frontend/"` / `"ban #8 frontend/"`，不从 `emojiScopes()` 推 ⇒ 台账确实在核对清单 | **PASS** |
| **AC#4** | 因加作用域而红的旧用例逐条重述，不删用例 | 顶层 `^func Test` 19→21、删除 test 函数 **0**；三条重述我逐条读过断言本体，仍测同一件事，且**顺带消掉了"点名叫 cmd/ 靠排序运气"**这个旧隐患；新增 `ban 7 tree gone` 子用例真跑真过（CI 口径 `runtests.sh` 日志为证） | **PASS** |
| **AC#5** | ban 文本零改动证明；`allowlist.txt` 仍 5 行非注释 | `main.go` 非注释删除仅 4 行、全是 walk 分支控制流；`emojiRe`/`banned (D23)`/`approval.decide` 在 diff 中命中 **0**；allowlist 非注释 **5**、该路径 diff 为空；`design/` 判据未动（16 文件全在白名单内）；`ban #6` 的 37 未降；`tools/` 内无越界改动 | **PASS** |
| **AC#6** | `cd tools/d22scan && go test -count=1 ./...` rc=0、`gofmt -l` 空、纯净树 `sh scripts/d22scan.sh` rc=0 且逐作用域贴全 | §6：`-count=1 -v` **rc=0**（RUN 31 / 顶层 PASS 21 / 子 10 / FAIL 0 / **SKIP 0**）；`gofmt -l` **空**；纯净树 `8b1b10f` **rc=0**、八行作用域逐字在 §1；且我复核了"起初那枚红"在 `5e8f87b` 快照上**确为 rc=1、唯一命中 `internal/winsec/winsec.go:126`** ⇒ 补正不是圆场 | **PASS**（补正处理判**诚实**，见 §6.3） |

### 两句直答

**(1) 票 96 能否挂 `-done`？—— 能，六格全 PASS，可挂 `96-...-done`。**
依据：三条独立路径的数（我的 git 数 / 我的 find 数 / 两条门自报）在纯净快照与工作树上都对得上；
三枚编译通过的变异分别把"收窄"“漏清单""对称收窄"打出红；票面六个 AC 的每条硬读我都复现到同一数值，
无一处需要"允许残留"来救判据。挂 `-done` 前**不必**再动代码；§9 那几条残留是**新立的账**，不是本票的返工。

**(2)"零 emoji 这条门现在真的看着面板"这句能不能对 owner 说？—— 能说，但要按下面的说法说，不能说成"两条门数相等所以安全"。**
- **依据（可复现）：** `ban #8` 的作用域清单今天含 `frontend/`，且走的是**无后缀白名单的全文件 walk**
  （`emojiScope.everyFile`，`walkEmoji` 的 `case sc.everyFile:` 分支），计数器 `s.emojiSeen[label]++`
  是在 `os.ReadFile(path)` **成功之后**才加的 ⇒ "examined 37"就是"读过 37 个文件的字节"。
  我自己在纯净快照上往 `.tsx` 注释 / `.md` / `.mjs` / 无后缀 / dotfile 各种一枚 emoji，**每一枚都把 CI 端成了 rc=1 并点名到行**；
  把 walk 退回后缀白名单，独立 walk 那条用例立刻报 `32 vs 37`；把 `frontend/` 从清单删掉，6 条用例红。
- **一句它不覆盖什么：** 它只盯 `design/`、`frontend/`、`internal/` 的 `.go`、`cmd/` 的 `.go` 这四条作用域，
  且**在 `frontend/` 内部排除 `node_modules/` 与 `testdata/`** ⇒
  **仓库里其余路径（根目录文件、`scripts/`、`docker/`、`docs/`、`tools/`、`models/`、`third_party/`）
  以及面板那棵树里的 `node_modules/`（第三方 vendored 依赖）仍然不在零 emoji 的门里**，
  那句话只能说到"我们自己的面板源码这棵树"为止；另外它只在 `sh scripts/d22scan.sh` 真跑时才有牙
  （本地那条脚本第 1 步可被 Go 测试缓存回放 = 票 99，CI 走 `runtests.sh -count=1` 不受影响）。

---

## 9. 允许残留 / 新立的账（都不改上面任何一格的裁决）

1. **`TestRealRepoBan8CoversFrontendTreeAtBan6sCount` 的"自报行"断言是自比较**：
   `want := fmt.Sprintf("... examined %d text files", front)` 的 `front` 就是被测数 ⇒ 该断言只钉行的存在/命名/格式，
   钉不出数值真假。真假由 `TestLedgerCountsMatchAnIndependentWalk` 承重（它确实承住了，M-A/M-C 为证）。
   票面写成"钉成用例"措辞偏强。**建议（可选）**：那条 `want` 改用独立 walk 的数，或在该用例里同时引 `count(...)`。
2. **"两条门数相等"本身不是 falsifier（我的 M-C 实测）**：同时给两条 walk 加一个目录排除 ⇒
   `33 == 33`、相等用例**仍绿**、扫描器仍打 `clean` **rc=0**，只有独立 walk 与两个 fixture 红。
   不是本票的缺陷（承重放对了），但 owner-facing 话术与后续票**必须引用独立 walk** 那一条。
3. **`ban8Scopes` 的 `make([]scanScope, 0, 3)`**（`main.go:353`）容量字面量在清单变 4 条后没跟着改 —— 纯外观，
   `append` 行为正确，台账照样打满 8 行。登记，不返工。
4. **两处"会腐烂的字面清单"在本票地界之外、但它一暴露就该有人认领：**
   (a) `.github/workflows/ci.yml:25` 注释仍写 "emojiScopes() … = design/ + internal/ + cmd/"（今天 4 条，含 `frontend/`）；
   它上一句还自我声明"scope names come from the tool, never from this comment"，紧接着就把清单抄了一遍并抄错了 ——
   正是票 71 当初立这条守卫要消灭的形状。
   (b) `tools/d22scan/main.go:30` 的 ban #6 文档注释仍写 "examines 35 text files"（今天 37）。
   两处本票都**不该改**（`main.go` 那行是票 88 的判据文本、ci.yml 不在它的 Packages 里），
   但**建议编排者开一枚一行的小票**把这两处从"抄来的清单"改成"指向函数"，或直接并进票 99。
5. **`goOnly` 与 `everyFile` 同时为真时 `goOnly` 静默优先**（`switch` 的 case 顺序），注释声明互斥但代码没断言。
   今天清单里没有这种条目，无行为影响；加固位置是 `emojiScopes()` 出口或 `walkEmoji` 入口的显式 panic。
6. **票面上一条 log 尾的 2 行 `next=` 被补正 commit 改写**（§6.3），严格 append-only 的话应原地保留再追加。格式级。
