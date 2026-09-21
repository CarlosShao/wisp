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
- 对生产树**只读**：§7 那次对工作树的操作是 `go run . -root <仓>`（扫描器不写盘），跑完 `git status --short`
  仍只有起手那两枚别人的改动。

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
票面 AC#1 的"工作树 40=40"我也单独复现了（§7），差 3 就是 `frontend/dist/` 里未入库的构建产物，两门同时看见，与票面解释一致。

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
不能说成"两条门数相等"（见 §10）。
