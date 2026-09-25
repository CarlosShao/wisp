# d22scan gitignore 修 · index-aware 第二程 r1（A214/F1+F2+F3）

**我是谁**：实现程（写码这一侧）。本件是取证件，不是裁决——**裁决必须出自非实现者**
（`SPEC-12 §4.3` #1、`AGENTS.md` §0.3）。本件里任何"我证过了"都不等于"通过"。
**票面**：编排者派单（A214 ② 点名的正解＝index-aware，**不采纳**验收者 §1.4 的 `.git` 存在性闸门）。
**被改交付**：`3bb99aa`（r1，已被 `987a79a`/`d22scan-gitignore-accept-r1.md` 整批退回，缺陷 F1/F2/F3）。

| 用途 | 锚点 | 说明 |
|---|---|---|
| 开工锚点 | `ee5a25e` | `git rev-parse HEAD` 现量（09:0x）。前一次读到 `13673e9`，其间 `ee5a25e`（台账 A213+A214）落地 |
| 交件时 HEAD | `d8390aa` | §1-§6 写作时刻（09:15）。它是 `internal/observe/**` 那一程的 commit，**没碰我这三枚文件**（`git show --numstat` 现量只带 `internal/observe/*` + 它自己那份 evidence） |
| 被验代码锚 | `3bb99aa` | 其父 `467f8a4`；`git log --oneline -3 -- tools/d22scan` 现量 = `3bb99aa` / `ed2c077` / `5e8f87b` ⇒ **交付后无人再动这枚仪器**，本批是唯一写者 |

**地界**：只写 `tools/d22scan/**`（三枚已有文件，无新增文件）＋本件。
`git status --porcelain` 现量另有 16 枚 ` D design/**` 与两枚 `?? design/`（owner/前端会话地界）——
**我一枚都没 stage、没 commit、没恢复**；每次提交前 `git diff --cached --name-only` 逐枚核过（见 §12）。

**仓库外工作区**（全部临时件在此，真树未写一字节；**只建不删**）：`D:\tmp\d22scan-index-fix-r1\`
- `snap-prebase\` = `git archive ee5a25e | tar -x`（净快照，无 `.git`）
- `bin\d22scan-pre.exe` = 在 `snap-prebase\tools\d22scan` 里 `go build` 出来的**修前**二进制（代码 = `3bb99aa` 那一版）
- `bin\d22scan-post.exe` = 本批**修后**代码编译出的二进制（从工作树源码编，产物落仓库外）
- `shapes\` = 净快照副本 + 仓库外 `git init` + `git add -A` + **`git add -f`** 两枚件（§2 的 F1 形状）＋两枚**未追踪且被忽略**的构建产物（§3 的反向对照）
- `ci-shapes\` = 对 `shapes` 的 `git init` 之后 `git archive HEAD | tar -x`（＝CI 那一侧：没有 `.git`，但那两枚受追踪的交付件在场）
- `stillbites\` = 净快照副本 + `git init` + 一枚受追踪非忽略的 `frontend/src/still-bites.tsx` + 一枚未追踪被忽略的 `frontend/dist/untracked-only.html`
- `mut-indexoff\` / `mut-semantics\` = 修后代码的两发**变异副本**（§6 / §5），仓库外跑
- `logs\` = `pre-shapes.log` `pre-ci.log` `post-shapes.log` `post-ci.log` `post-stillbites.log` `pre-clean-archive.log` `post-clean-archive.log` `pre-suite.log` `post-suite.log` `post-gate-worktree.log` `mut-M1.log` `mut-M2.log` `roster-pre.txt` `roster-post.txt` `census.txt`

**读数环境**：win32 · `git version 2.52.0.windows.1` · `go version go1.27.1 windows/amd64` ·
`date` 现量 `2026-09-25 09:15 +0800`。测试**串行**跑（`internal/observe/**` 那一程正在用 CPU）。

---

## 0. 派单前提核过一遍：有一条**我读了盘上事实之后同意、一条要补一句**

1. **"验收者点名的 `.git` 闸门不够"——我独立复核，成立，且盘上有硬证据。**
   验收者自己那套形状（`accept §1.3`）用的是 `forcedfull\` ＝ **带 `.git` 的真仓库**，其 `git ls-files -i -c`
   现量点名 `frontend/dist/tracked-forced.tsx`；那一棵上"有没有 `.git`"的答案是**有**，
   所以 `.git` 闸门在那棵上**根本不会关掉过滤** ⇒ 真实检出（CI clone 与工作树）里
   "受追踪且模式命中"那一支仍然被吞。票面给的推理方向对。
   **但要补一句**（免得下一位把 §4 读成"`.git` 闸门完全没用"）：那闸门确实能治
   `git archive` 那一支（我这轮的 `ci-shapes` 读数就是证据），它治不了的是**带 index 的那一支**——
   而带 index 的那一支才是 CI 的真实形状。所以我实现的是"问 index"，不是"看 `.git` 在不在"。
2. **"用 `git ls-files -i -c --exclude-standard` 当判据"——成立，但我把它当**诊断**用了一半。**
   盘上事实（我这台，真工作树只读）：`git ls-files -i -c --exclude-standard` 现量 **0 行**
   （与 `accept §1.2` 两枚口径均为空一致）。这枚命令给的是"受追踪 **且** 模式命中"的交集，
   而 `skip()` 要的正确判据其实是 git 的真判据"**受追踪就永远不跳**"。交集 ⊂ 全集，
   所以只查交集会漏掉一类：我的**模式匹配器比 git 多跳**（一枚受追踪件、git 说没被忽略、
   但我的解析说被忽略）——那一类**不会出现在 `-i -c` 的名单里**，交集闸门救不了它。
   ⇒ 我实现的是**两枚都问**：`git ls-files -z`（全集＝结构性保证）＋
   `git ls-files -z -i -c --exclude-standard`（票面点名的仪器＝诊断＋冗余并集成员）。
   这是**超出票面点名的一味**，理由与代价都写在这里，交回裁决者判：见 §11.3。
3. **"F3 两枚语义偏差"——成立，而且我用真 git 逐枚复算过**（不是照抄验收者的探针）：

   | 支 | 真 git 现量（仓库外 `probe1/2/3`） | r1 行为 | 本批 |
   |---|---|---|---|
   | 行首空格 | 规则 `  dist/*` ⇒ `git check-ignore -v dist/index.html` **rc=1**（不忽略）；规则 `  lead/*` ⇒ `git check-ignore -v '  lead/x.txt'` **rc=0** 且回显 `.gitignore:1:  lead/*` | `TrimSpace` 把规则读成 `lead/*` ⇒ 吞掉 `dist/index.html`＝**少扫** | 只剥尾随 `" \t\r"`（`:557`） |
   | `foo/**` | 规则 `dist/**` ⇒ `git check-ignore -v dist`（目录本身）**rc=1**；`dist/index.html` rc=0；`!dist/.gitkeep` rc=0 | 匹配器把 `dist` 目录本身判为忽略 ⇒ 父目录被剪、否定失效＝**少扫** | 尾随 `**` 至少吃一段（`:598`）；中段 `**` 仍可吃零段（`a/**/b` 打 `a/b`，`git check-ignore` rc=0 实测） |
4. **一条派单没点名、但我撞见的事实**：`git ls-files` **没有** `--full-index-name` 选项
   （`git 2.52` 现量 `error: unknown option`，正确名是 `--full-name`）。
   我没有用任何一枚，而是**把 `cmd.Dir` 设成被扫根**，输出天然就是根相对路径，
   并用 `git rev-parse --show-prefix` 先证明"这枚根就是一次仓库顶"（见 §11.2）。

---

## 1. 改了什么（三枚文件，无新增文件）

`git diff --numstat`（锚 `d8390aa`，工作树）：

```
304	33	tools/d22scan/gitignore.go
44	22	tools/d22scan/main.go
340	2	tools/d22scan/scan_test.go
```

1. **`gitignore.go` 新增 index 层**（现号）：`gitIndexState` `:157`、`holds()` `:184`、
   `runGitIndex()` `:202`、`runGitAt()`、`gitIgnore.index()` `:382`；
   `skip()` `:338` 的判据从"模式命中即跳"改成"**模式命中 且 git 说 index 不持有它**"，
   并且 `index 问不到 ⇒ 一条规则都不应用`。
   三枚命令各带 `context.WithTimeout` 10s（`indexTimeout` `:190`）——
   **没有一处用墙钟差**（`AGENTS.md §1.2` 明禁）：全文不出现 `time.Now()`，超时只由 context 计时器决定。
2. **`note()` `:400` 变成多行**（原为一行）：
   ⓐ index 问不到时**无条件**打印 `gitignore rules NOT APPLIED - <原因>; every path in every scope is being scanned …`
   （即便本轮一件都没跳、即便 rc=0）；ⓑ 原有的 `skipped as git-ignored: …`；
   ⓒ `N path(s) git reports as TRACKED and matching an ignore rule … none of them was skipped: <逐枚点名>`。
   ⓐ/ⓒ 是 `main.go:1307` 那一处 `stats.ign.note()` 现成打出来的，**没动 main 的调用形状**，
   所以 verdict() 的 output 逐行断言不受影响（`note` 故意在 ledger 之外，`:158` 那段理由原样保留）。
3. **两枚语义修正**（F3）：`parseIgnoreFile` 只剥尾随空白＋一个尾随 CR（`git` 的 CR 处理我实测过：
   CRLF 的 `.gitignore` 规则仍生效）；`pathMatch` 的尾随 `**` 必须吃至少一段。
4. **`main.go` 的注释**：`:55-77` 重写（F2 那两句现在**指向钉它的断言**）、`scanner.ign` 字段注释、
   `walkGo` 里那句"A git-ignored .go file is not on anyone's delivery path"（它对受追踪那支是**假**的，已改）、
   `checkRoot` 的 matcher 多问一次 index（一次运行两枚 matcher ⇒ 三枚 git 命令 ×2）。
   **删除的 22 行逐枚核过全是注释行**（`git diff -U0 | grep '^-'` 现量，逐条贴 §9 表），没有一行是代码。
5. **`scan_test.go`**：新增两枚测试（§11.1）＋三枚 git 夹具 helper（`:1384` `gitCommand`、
   `:1410` `gitIndexFixture`、`:1421` `gitForceAdd`）＋第二副本 `ignoredLikeGit` 接 index（`:1758`）
   ＋`gitTrackedIndex`（`:1775`）＋`TestGitIgnoreRuleSemantics` 追加 3 条规则与 6 个用例。
   **删除的 2 行**＝`ignoredLikeGit` 的函数签名行与它在独立走查里的调用行（同批改签名，非删断言）。

**缓存与稳定性（票面要求"决定并写明"）**：
- index **每枚 matcher 问一次**，惰性（`index()` 里 `g.idx == nil` 才问），**跑完不过期**。
  理由写在 `:376-381`：四枚 walk 共用一枚 matcher，按路径重问＝上千次 subprocess；
  而 mid-walk 重问会让同一枚路径在两个 ban 下给出不同答案（分母就会随遍历顺序动）。
  `main()` 里 `checkRoot` 与 `scanWithStats` 各持一枚 matcher ⇒ 一次运行问 git **两轮**（6 次 spawn，实测约 0.5s），
  两轮的正确答案必然相同，**若不同就说明 index 在跑的过程中变了**，那本身是要暴露的事。
- **分母与顺序无关**：`note()` 里三处名单全部 `sort.Strings`（`pruned` 目录、`origins` 规则文件、
  `ignoreSet` 受追踪-被忽略路径），map 迭代顺序进不了输出。
  `files`/`pruned`/`origins` 的计数靠 `counted` 去重（r1 已有，本批未动）。
- `skip()` 的顺序：**先 `rel == ""`（walk 自己的起点目录，永不跳）→ 再问 index → 再问模式**。
  起点那一支排在前面，所以 `TestGitIgnoreRuleSemantics` 里"starting directory must never be skipped"
  那一支**不会为了证它而 spawn git**。

---

## 2. 格一：F1 去红（票面要求"把验收者的形状自己重造一遍"）

形状（仓库外 `shapes\`）：净快照副本 → `git init -b main` → `git add -A`（**1125 枚受追踪**）
→ 再落两枚"模式命中但被强推进 index"的交付件 → `git add -f`。git 自己的两台仪器先作证：

```
$ git ls-files | wc -l                                        -> 1127
$ git ls-files -i -c --exclude-standard
frontend/dist/tracked-forced.tsx
internal/build/leak.go
$ git check-ignore -v frontend/dist/tracked-forced.tsx internal/build/leak.go
rc=1                                   <- 带 index 的真判据：两枚都**没被忽略**
$ git check-ignore -v --no-index frontend/dist/tracked-forced.tsx internal/build/leak.go
frontend/.gitignore:12:dist/*	frontend/dist/tracked-forced.tsx
.gitignore:14:build/	internal/build/leak.go
rc=0                                   <- 只按模式（＝r1 实现的判据）：两枚都"被忽略"
```

**同一棵 `shapes`，两版二进制，`-root` 指过去，其余一字不差**（锚 `d8390aa`）：

| 读数 | 修前 `bin/d22scan-pre.exe` | 修后 `bin/d22scan-post.exe` |
|---|---|---|
| rc | **0** | **1** |
| ban #6 frontend/ | 40 | **41** |
| ban #8 frontend/ | 40 | **41** |
| bans #1-5 internal/ | 203 | **204** |
| examined production Go | 225 | **226** |
| findings | 无 | `internal/build/leak.go:5: [bare-goroutine] bare \`go worker(...)\` is banned (D22/D38b, R16…)`<br>`frontend/dist/tracked-forced.tsx:1: [emoji] ban #8 glyph in scope frontend/ is banned (D23)…` |
| note | `skipped as git-ignored: 2 file(s) under 2 ignored director(ies) [frontend/dist/assets/, internal/build/], decided by .gitignore (1 path(s)), frontend/.gitignore (3 path(s))` | `skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/], decided by frontend/.gitignore (2 path(s))` ＋ **`2 path(s) git reports as TRACKED and matching an ignore rule (git ls-files -i -c --exclude-standard) - a tracked path is part of the delivery, so none of them was skipped: frontend/dist/tracked-forced.tsx, internal/build/leak.go`** |

⇒ 两发都翻回来了：`≤` 那一发（ban #8）**和**裸 `go worker()` 那一发（`AGENTS.md §1.2` 明禁的安全形状，bans #1-5）
从"对门永久隐身"变成"点名且 rc=1"。修前的 note 还写着 `internal/build/` 被剪掉，修后那一剪**没有了**
（目录里有受追踪件 ⇒ `holds()` 拒剪），这正是 `accept §6.9` 那发的最坏形状被拆掉的样子。

**CI 那一侧的形状**（`ci-shapes\` ＝ 对 `shapes` 再 `git archive HEAD | tar -x`，**没有 `.git`**，
但那两枚受追踪交付件**在场**——`cat` 现量 `export const doc = "ready ≤";`）：

```
$ bin/d22scan-pre.exe  -root ci-shapes   -> rc=0, frontend/=40, internal/=203, clean   （＝accept §1.3 的读数，独立复现）
$ bin/d22scan-post.exe -root ci-shapes   -> rc=1, frontend/=41, internal/=204, 两枚 finding 点名
   第一行输出：d22scan: gitignore rules NOT APPLIED - git cannot be consulted in
   D:/tmp/d22scan-index-fix-r1/ci-shapes: `git rev-parse --show-prefix`: fatal: not a git
   repository (or any of the parent directories): .git; every path in every scope is being
   scanned, so the counts below may include build output (A207's machine-dependent
   denominator). This is the loud direction: a scanner that cannot ask git which paths are
   tracked does not get to skip any.
```

⇒ 这一支**同时**由两条独立机制保：index 那一半在带 `.git` 的树上生效，"问不到就全扫＋响亮"那一半在
`git archive` 副本上生效。验收者 §1.4 的闸门只覆盖后一支；两支都盖住才是"不往少扫方向偏"。

---

## 3. 格二：不蒙眼对照，**两个方向**

| 方向 | 树 | 读数 |
|---|---|---|
| 受追踪、**不被忽略**的违规件仍必须红 | `stillbites\`（`git init`＋`git add frontend/src/still-bites.tsx`） | **rc=1**，`frontend/src/still-bites.tsx:1: [emoji] ban #8 glyph in scope frontend/ is banned`，`ban #6 frontend/=41`、`ban #8 frontend/=41`，`skipped as git-ignored: 1 file(s)` |
| 真·被忽略（未追踪）的件**不得**算进分母（＝A207 本体，不许退回去） | 同一棵里的 `frontend/dist/untracked-only.html`(含 `≤`) ＋ `shapes\` 里的 `frontend/dist/untracked-bundle.html` 与 `frontend/dist/assets/hash-a12.js` | 上面那行 `skipped as git-ignored: 1 file(s)` 点名的就是它；分母 **41 而不是 42**；finding 名单里**没有**它。`shapes` 那棵里三枚未追踪构建产物同理（修前修后都不进分母），而修后 `frontend/=41` 恰＝ 40（CI 基线）＋ 1（受追踪那枚 `.tsx`） |
| 同一条不变量在单元层 | `TestTrackedPathsAreNeverSkippedByTheIgnoreFilter` | 一枚 `seed()` 造两棵**字节相同**的树，唯一差别＝`git add -f` 有没有落下；断言分母 **+1 而非两边一起动**、两枚 finding 点名、未追踪那一半**不得**出现这些 finding（"A207's exclusion was lost" 那句就是防我把修退回去的）、note 两侧各自点名、verdict 一侧 1 一侧 0 |

**测试层同一对形状的真 git 前置断言**（不是推理）：未追踪那一棵 `git ls-files -i -c` 必须为空、
受追踪那一棵必须点名两枚、且 `git check-ignore --no-index` 必须仍说"模式命中"——
三枚都写进断言，所以将来若 `git` 改了这些仪器的语义，这一格会自己红给我看。

---

## 4. 格三：没有 `.git` / 问不到 index 的形状——**必须响亮**

| 形状 | 读数 | 判 |
|---|---|---|
| 净快照（`snap-prebase`，无 `.git`） | rc=0、八数逐枚同 §7 表、**第一行就是** `gitignore rules NOT APPLIED - … not a git repository …` | ✅ 绿是**带着自陈**的绿，不是静默 |
| `ci-shapes`（无 `.git`、场内有受追踪违规件） | rc=1 ＋ 同一行自陈 ＋ 两枚点名 | ✅ 见 §2 |
| `-root` 指到不是一棵 wisp 仓的地方 | `checkRoot` 硬拒（`accept §3` 已量：`missing go.mod`），本批**未改**这一层 | 不适用；未重跑 |
| 工作树（有 `.git`、是仓库顶） | `sh scripts/d22scan.sh` rc=0，**不打印** NOT APPLIED 行，打印 `skipped as git-ignored: 1 file(s) …`（§8） | ✅ 该沉默时沉默 |
| 单元层两枚 | `TestIgnoreRulesAreOffWhenTheIndexCannotBeAsked`：无 index ⇒ 三枚分母**全部包含**被忽略件、两枚 finding 点名、note 单行且含 `not a git repository`、verdict=1；**再一棵干净树**（无违规）⇒ rc=0 但 note 仍含 `gitignore rules NOT APPLIED` | ✅ "安静地绿"被断言本身拒绝 |

---

## 5. 格四：F3 两枚语义（含变异对照＝这一格真有牙）

在 `mut-semantics\`（修后代码，把两处改动逐枚还原）跑：

```
$ go test -count=1 -v -run 'TestGitIgnoreRuleSemantics|TestWalksSkipGitIgnoredPaths|TestTrackedPathsAreNeverSkippedByTheIgnoreFilter|TestIgnoreRulesAreOffWhenTheIndexCannotBeAsked' ./...
M2 rc=1
--- PASS: TestWalksSkipGitIgnoredPaths (0.36s)
--- PASS: TestTrackedPathsAreNeverSkippedByTheIgnoreFilter (0.99s)
--- PASS: TestIgnoreRulesAreOffWhenTheIndexCannotBeAsked (0.11s)
--- FAIL: TestGitIgnoreRuleSemantics (0.00s)
    scan_test.go:1686: decide("deep", isDir=true) = ignored true, want false (a TRAILING `**` never matches its own parent directory …)
    scan_test.go:1686: decide("deep/.gitkeep", isDir=false) = ignored true, want false (the negation still reaches it …)
    scan_test.go:1686: decide("lead/inside.txt", isDir=false) = ignored true, want false (git does NOT strip leading whitespace …)
    scan_test.go:1686: decide("  lead/inside.txt", isDir=false) = ignored false, want true (and it does bind the directory whose name really starts with those two spaces)
FAIL	github.com/CarlosShao/wisp/tools/d22scan	1.496s
```

⇒ 新加的 6 个用例（含 `deep/x`、`deep/x/y.js` 两枚**正向**用例，防"为了修父目录把整条 `**` 语义废掉"）
逐枚咬得住。正/反两向：`lead/inside.txt` 不该被吞（少扫方向），`  lead/inside.txt` 该被吞（多扫方向）。
⚠ 一枚取证口径要写清：第一版 M2 我漏了 `-v`，那份日志里**只有 FAIL 行**（Go 不加 `-v` 不打印 PASS），
所以"另外三枚仍 PASS"当时**没有凭据**；上面这份是补跑 `-v` 的读数（`logs/mut-M2.log` 已被重跑覆盖）。
另外 `accept §2` 那行"`/foo` 锚定 未测"**本批没补**（它不在 F1/F2/F3 射程里，且方向是**多锚定**＝少扫的潜在面；
登记进 §12 没测清单第 6 条，交回裁决者判要不要另开一格）。

---

## 6. 格五：新测试的牙齿＝把 index 那一味摘掉

`mut-indexoff\`＝修后代码，只在 `skip()` 里把两处 index 闸门还原成 r1 的"纯模式"：

```
$ go test -count=1 -v -run 'TestTrackedPathsAreNeverSkippedByTheIgnoreFilter|TestIgnoreRulesAreOffWhenTheIndexCannotBeAsked|TestWalksSkipGitIgnoredPaths|TestGitIgnoreRuleSemantics' ./...
M1 rc=1
--- PASS: TestWalksSkipGitIgnoredPaths (0.35s)
--- FAIL: TestTrackedPathsAreNeverSkippedByTheIgnoreFilter (1.07s)
    scan_test.go:1479: ban #6 counted 2 frontend/ files tracked and 2 untracked - the force-added .tsx must add exactly 1
    scan_test.go:1508: ban "emoji" did NOT fire at the tracked path frontend/dist/tracked-or-not.tsx - a rule that hides a delivered byte is a blindfold, not a fix (all findings [])
    scan_test.go:1508: ban "bare-goroutine" did NOT fire at the tracked path internal/build/leak.go - …
    scan_test.go:1534: verdict on the tracked half must be 1 (findings), got 0
--- FAIL: TestIgnoreRulesAreOffWhenTheIndexCannotBeAsked (0.16s)
    scan_test.go:1564: ban #6 examined 2 frontend/ files, want 3 - with no index to consult nothing may be excluded
    scan_test.go:1581: the index-less tree must go red at frontend/dist/copied-delivery.tsx, findings=[]
    scan_test.go:1597: verdict must be 1 on findings, got 0
--- PASS: TestGitIgnoreRuleSemantics (0.01s)
```

同时**记录一枚我不想让它被读歪的事实**：`TestWalksSkipGitIgnoredPaths` 与
`TestGitIgnoreRuleSemantics` 在 M1 下**仍 PASS**——这正是 `accept §1.4` 说的"本批没有任何一枚测试能发现
这一发变异"的形状：那两枚测的是"该跳的跳了"，它表达不了"受追踪且模式命中"那一支。
现在能表达的是新那两枚，且上表证明它们打得红。

---

## 7. 格六：净快照八数**逐枚不动**（票面 #4）

树 ＝ `snap-prebase`（`git archive ee5a25e`，无 `.git`）；两版二进制跑同一棵（日志
`pre-clean-archive.log` / `post-clean-archive.log`）：

| scope | 修前 | 修后 | 同数？ |
|---|---|---|---|
| bans #1-5 internal/ | 203 | 203 | ✅ |
| bans #1-5 cmd/ | 22 | 22 | ✅ |
| ban #6 frontend/ | 40 | 40 | ✅ |
| ban #7 internal/tools/ | 18 | 18 | ✅ |
| ban #8 design/ | 16 | 16 | ✅ |
| ban #8 frontend/ | 40 | 40 | ✅ |
| ban #8 internal/ | 405 | 405 | ✅ |
| ban #8 cmd/ | 39 | 39 | ✅ |
| examined production Go | 225 | 225 | ✅ |
| rc | 0 | 0 | ✅（修后多打印第一行 NOT APPLIED，见 §4） |

⇒ 与票面期望的八个数逐枚相同，**没有一个数需要解释它的移动**。
工作树那一侧（`sh scripts/d22scan.sh` step 2，锚 `d8390aa` 的工作树，含 `internal/observe/**` 那一程的在飞件）：
`ban #6 frontend/=40`、`ban #8 frontend/=40`（A207 的合流**没退**）、`bans #1-5 internal/=203 cmd/=22`、
`ban #7 internal/tools/=18`、`ban #8 design/=32`、`internal/=406`、`cmd/=39`、rc=0。
`design/=32` vs 净快照 `16` ＝ `accept §6.5`/`A211③` 那枚**已定性的残留**（100% 未追踪且未被忽略），
`ignore` 过滤器吞不到也不该吞 ⇒ 本批**没动**它。`internal/=406` 比 `405` 多一枚＝那一程新落的
`window_wait_136_test.go`（`_test.go` 进 ban #8 不进 bans #1-5，故只有 405→406 这一枚动，203 未动）——
**不是本批改出来的**，`git show --numstat d8390aa` 现量只带 `internal/observe/**`。

⚠ **F9（`accept §11.3`）本批一字未动**：owner 那 16 枚未提交的 `design/**` 删除一旦入库，
`ban #8 design/` 走 0 枚 ⇒ 门硬退出 2。我没碰那 16 枚删除、没碰 `emojiScopes()`、没碰 `driftedAbsentScope`，
也**没有**为它做任何"预防性放宽"。

---

## 8. 格七：`sh scripts/d22scan.sh` 端到端（票面 #5；step 1 与 step 2 四数分开报）

锚 `d8390aa`，真工作树，串行（`logs/post-gate-worktree.log`，rc=0，实测 10.8s）：

- **step 1（正向对照，四数只由它产出）**：`PASS=28 FAIL=0 SKIP=0, === RUN=68, '[no tests to run]'=0`
- **step 2（真扫）**：八数见 §7 工作树行，`rc=0`，打印 `skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/], decided by frontend/.gitignore (2 path(s))`，**不打印** NOT APPLIED（工作树问得到 index）

修前基线（同一套命令跑在 `3bb99aa` 的代码上，净快照 `snap-prebase`，`logs/pre-suite.log`）：
`PASS=26 FAIL=0 SKIP=0, === RUN=66`。⇒ 差值 **+2/+2/0/0**，名单差集（`diff roster-pre.txt roster-post.txt`）
现量**恰两行、只有新增**：

```
13a14
> TestIgnoreRulesAreOffWhenTheIndexCannotBeAsked
21a23
> TestTrackedPathsAreNeverSkippedByTheIgnoreFilter
```

⇒ 名单没缩、没 `--- SKIP`、step 1 绿（`set -eu` 之下 step 1 红会让 step 2 根本不跑；本批两枚 step 都跑了，
读数逐枚在上面）。

---

## 9. 格八：契约轴**未动且已证**（票面 #6）

```
$ for a in 3bb99aa^ 3bb99aa HEAD WORKTREE; do printf '%-10s ' "$a"; \
    ( [ "$a" = WORKTREE ] && cat tools/d22scan/main.go || git show $a:tools/d22scan/main.go ) \
      | grep 'emojiRe = regexp.MustCompile' | tr -d '\n' | sha256sum | cut -c1-16; done
3bb99aa^   a6150ef689cc20c9
3bb99aa    a6150ef689cc20c9
HEAD       a6150ef689cc20c9
WORKTREE   a6150ef689cc20c9
$ <同一取法> | tr -d '\n' | wc -c
135
```

- `emojiRe` 那一行**逐字节同串**，且我自己重导的 sha 前缀与 `accept §5.1` 记录的 `a6150ef689cc20c9` **相同**。
  ⚠ 一枚口径要说清免得下一位以为打架：那一行本身是 **135 字节**（与验收件写的 135 相同），
  而我先前那次单发命令写的 `wc -c` ＝ **136** 是**管道里带了行尾 `\n`** 的数；
  上表统一用 `tr -d '\n'` 之后再算，所以四枚 sha 同值、长度与验收件同值——**不是两处不同串**。
- `git diff HEAD -- tools/d22scan/main.go | grep -E '^[+-].*(emojiRe|1F000)'` 现量 **rc=1（零命中）**
  ⇒ 我这 44/22 行没有一行触及字符类。
- `git diff HEAD -- tools/d22scan/allowlist.txt` = **0 字节**；`git status --porcelain` 不列它 ⇒ 零豁免新增。
- 断言删除：`git diff -U0 -- tools/d22scan/scan_test.go | grep '^-'` 现量**只有两行**
  （`ignoredLikeGit` 签名行＋它唯一的调用行，均被三参数版替换）。**没有一枚 `t.Error/t.Errorf/t.Fatalf` 消失**
  （新测试净增，旧断言未改号：`TestWalksSkipGitIgnoredPaths` 的 5/5/3/3 与 finding 名单断言**逐枚原样**，
  只是夹具多了一步 `git init`＋`git add -A`）。
- `git diff -- tools/d22scan/ | grep -E '^\+.*t\.Skip'` 现量 **rc=1** ⇒ 没新增任何 Skip。
  （`runtests.sh` 把 SKIP 判为不通过，我若加 Skip 这一步会当场红，形状上就堵住了。）
- 阈值/golden：`grep -E '^\+.*(threshold|golden)'` 命中的两行**全是注释**
  （`gitignore.go` 里"as a SOURCES OF A SKIP"那段的散文），无 SLO/golden 变更。
- `gofmt -l tools/d22scan` **空**、`go vet ./...` **空**（锚 `d8390aa` 工作树）。

---

## 10. 同源副本普查（票面点名的命令，`logs/census.txt`）

`grep -rn 'node_modules\|/dist/\|skipIgnore\|ignoredLikeGit\|matchIgnore' --include='*.go' internal tools cmd`
＝ 55 行命中，按文件分组：**tools/d22scan/scan_test.go 36 · tools/d22scan/main.go 6 · tools/d22scan/gitignore.go 6 ·
internal/panel/frontend_hygiene_test.go 3 · internal/panel/composer_test.go 3 ·
internal/risk/pathresolver_rewrite_account_test.go 1**。
其中**可执行的跳过逻辑**逐枚（现号）与处置：

| # | 位置（现号） | 抄的是哪条 | 归属 | 本批处置 |
|---|---|---|---|---|
| ① | `main.go:648`/`:662` `walkGo` | 目录名 + `s.ign.skip` | 我 | **动**（skip 内部改 index-aware；调用形状一字未动） |
| ② | `main.go:841`/`:853` `walkText` | 同上 | 我 | **动**（同上） |
| ③ | `main.go:919`/`:952` `walkEmoji` | 同上 | 我 | **动**（同上） |
| ④ | `main.go:1168`/`:1173` `checkRoot` | 同上 | 我 | **动**（同上；它多问一轮 index，见 §1） |
| ⑤ | `scan_test.go:1192-1205` 独立走查 | `testdata/.git/node_modules` 目录名 + `ignoredLikeGit(rel, tracked, indexOK)` | 我 | **动（同批）**：oracle 接上 index＋fallback 两条参数 |
| ⑥ | `scan_test.go:1758` `ignoredLikeGit` | **手抄的 ignore 政策**（`frontend/dist/` 前缀 − `.gitkeep`） | 我 | **动（同批）**：签名 `+tracked +rulesOn`；`!rulesOn ⇒ 一条都不跳`；`tracked[rel] ⇒ 不跳` |
| ⑦ | `internal/panel/frontend_hygiene_test.go:289` | 子串 `/node_modules/`、`/dist/` | **编排者/前端会话** | **未动（正确）**：票面点名"不同 owner，本批不许改，只登记"。`accept F8` 已把它那一族的既有暴露（受追踪的 `frontend/dist/x.tsx` 会被**这台仪器**从它自己的对照里跳掉，且早于本修）记给它的主人——**本批的 index-aware 修不到那一枚，这条账仍在它那边**，我在 §12.9 明确写"没测它今天是否仍绿" |
| ⑧ | `internal/panel/composer_test.go:286`/`:434`/`:606` | case 表列目录名（夹具数据，不是跳过政策） | 同上 | 未动（正确） |
| ⑨ | `internal/risk/pathresolver_rewrite_account_test.go:223` | case 表列 `.git/.scratch/node_modules/dist/build/testdata/design/docs` | 同上 | 未动（正确） |
| ⑩ | `main.go:390` | ban #6 的**说明文字**里那句 `not node_modules/, testdata/ or a gitignored path` | 我 | **未动**：文字仍成立（射程只增不减），改它只是 cosmetics；登记在此免得被当成漏数 |

**判：新增了几份可执行副本？** 零。①-⑥ 全都收敛到同一枚 `gitIgnore.skip()`；
⑤⑥ 是测试侧**故意留的第二份手抄**（r1 就有，本批按票面**同批**改）。
普查里"手抄的跳过规则"这一族，按位置数还是 ⑥ 一枚，与 `accept §7` 的 9 行表对齐（我这里是 10 行，
多的那一行是 ⑩ 说明文字，它不可执行）。

---

## 11. F2 那两句话的下场：一句兑现、一句改窄（票面要求"要么真要么重写"）

1. **`main.go:61-64`（修前）"it does not suppress a finding on a TRACKED path"**
   ⇒ **兑现**，并且它现在**点名一枚能打死它的断言**。新文本（现号 `main.go:55-79`）写的是机制
   （"`skip()` asks `git ls-files -z` and `git ls-files -i -c --exclude-standard` and refuses to skip
   anything either one reports"）＋测试名（`TestTrackedPathsAreNeverSkippedByTheIgnoreFilter`，
   `scan_test.go:1439`）＋第二枚测试名（`TestIgnoreRulesAreOffWhenTheIndexCannotBeAsked`，`:1551`）。
   证伪方式已在 §6 实测：摘掉 index ⇒ `--- FAIL` 并说出"blindfold, not a fix"。
2. **`gitignore.go:33-34`（修前）"the instrument and the repository's ignore rules cannot disagree"**
   ⇒ **删除并改成可证伪的窄命题**（现号 `gitignore.go:34-48`）："git 的判据是规则 **加 index**；
   本文件的主张是 `for a path git reports as tracked, skip() returns false`"——这一句能被
   §6 那种变异直接打红。原句那种"cannot disagree"的**全称担保**不再存在于这枚文件里
   （`grep -c 'cannot disagree' tools/d22scan/gitignore.go` 现量 0）。
   另外 `:61-65`（修前）"every unsupported shape fails in the loud direction … never less"
   也重写成"列出的**不支持**形状只会让本匹配器**少跳**；唯一可能多跳的形状（问不到 index）
   由上一条规则处理：一条都不应用"——这条现在与 §5/§6 的实测一致。
3. **一处超出票面点名的实现选择，主动交回裁决**（免得下一位当成"实现者自己扩了射程"）：
   票面点名 `-i -c` 一枚命令，我实现了**两枚**（`-i -c` 之外再加全量 `ls-files -z`）。
   理由见 §0.2（交集挡不住"匹配器比 git 多跳"那一类，而 F3 只治了今天已知的两例）。
   代价：一次运行多 3 次 git spawn（实测约 0.5s）、`gitTrackedIndex` 在测试侧多问 2 次。
   **如果裁决者认为应当严格按票面只查交集**，砍掉 §1.1 里 `holds()` 的 `ix.tracked[rel]` 一支即可，
   §2/§3/§6 的读数不会变（两棵树上 `-i -c` 已点名同样两枚），变的只有未来那类未知形状。

---

## 12. 我**没有**测什么（绝不让实现者的沉默被读成批准）

1. **没让裁决者看**：本件是实现者自证。**F1/F2/F3 的合上与否必须另一位非实现者判**
   （`SPEC-12 §4.3` #1/#3），并且它已经具备那发变异的复现装置（`accept §10.1` 那句"别重造"）。
2. **`-root` 指到"别人仓库的子目录"那一支只实现了、没实测**：`runGitIndex` 里
   `rev-parse --show-prefix` 非空 ⇒ 判"问不到"⇒ 全扫＋响亮。我**没有**造出这棵树跑一遍
   （它要求把快照放进某枚真仓库的子目录里，而 `AGENTS.md §1.4` 禁"在仓库目录内建 checkout"）。
   `TestGitIgnoreRuleSemantics` 也表达不了它（`decide()` 不碰 index）。**这一支的响亮与否未证。**
3. **git 二版本身不在 PATH 的形状没实测**：只实测了"有 git 但不是仓库"（§2/§4 两发）。
   "PATH 里没有 git"走的是同一条 `!ix.ok` 分支（`runGitAt` 返回 `exec: "git": …not found`），
   我用代码路径推理它响亮，**没有真跑一遍**（这台机器 git 在 PATH 里）。
4. **超时那一支没实测**：`indexTimeout = 10s` 与 `ctx.Err()` 分支是写了但**没被任何一发读数触发过**
   （git 在这三枚问题上最慢 0.14s）。没有一枚测试能构造它——要构造就得 mock `exec`，
   而 `AGENTS.md §1.3` 只许在 C8/C5/C17/`wisp run` 那四枚接缝注入。**登记为未证，不当已证。**
5. **`.git/info/exclude` 与 `core.excludesFile` 的行为没实测**（`accept §2` 的 P6 那一支）。
   本批它们**仍是"不作为 SKIP 来源"**，所以 A207 的"两形状不同数"在这一支**照旧活着**（F5 未修，不在本批射程）。
   ⚠ 一处**新的耦合**必须让下一位知道：`-i -c --exclude-standard` 这个查询**会把 info/exclude 与
   core.excludesFile 计入"受追踪且被忽略"的诊断名单**（受追踪件的 index 保护与 SKIP 来源无关），
   ⇒ ⓐ 那两枚外部文件现在能影响**note 的第三行**、ⓑ 若有人把 `frontend/dist/*` 从 `.gitignore` 挪进
   `info/exclude`，那枚**受追踪**件仍被扫描（多扫方向）。**两枚都没实测**。
6. **`/foo` 前导斜杠锚定仍无测试**（`accept §2` 那行"钉＝否"未补），F3 之外的语义缺口**一律没补**：
   `\#` 转义（P3）、`[!a-z]`／`[a-]` 畸形类、段内 `a**b`、`\` 结尾的转义尾随空格。
7. **step 1 现在需要 PATH 里有 `git`**——这是本批**新引入的依赖面**，必须写在案上：
   §11 那两枚新测试（以及 `TestWalksSkipGitIgnoredPaths`，它现在要 `git init`）拿不到 git 就
   `t.Fatalf` 当场红（**没有** Skip；`runtests.sh` 也判 SKIP 为不通过）。
   我**没有**在无 git 的机器上验证 step 1 的形状，也**没有**测 `git` 存在但
   `core.autocrlf`/`safe.directory`/只读 index 那类异常配置。**扫描器本身仍可无 git 运行**（§4），
   变的是"门的正向对照"变严了一格：跑不了 git 的机器跑不了这一格 step 1。
8. **没跑 root module 的任何测试**（`internal/**`、`cmd/**` 全 suite）。
   ⚠ 一枚**我自己的操作失误要入账**：为了取"修前基线名单"，我在 `snap-prebase/tools/d22scan` 里误打了
   `sh runtests.sh -C . ./...`（`-C .` 被脚本解析成**快照根**），结果跑的是 **root module 的全套测试**
   （`logs/pre-suite.log` 的第一版，838 PASS / 若干 FAIL，rc=1）。那一批 FAIL **不是本批造成的**、
   也**不是我读的基线**——我发现后按 `sh tools/d22scan/runtests.sh -C tools/d22scan ./...` 重跑，
   正确基线是 `PASS=26 / RUN=66`（本文 §8 用的就是重跑那一版）。
   全程在仓库外快照里发生，真树未写一字节。`internal/panel/frontend_hygiene_test.go`（⑦）今天是否仍绿：**没跑、未测**。
9. **⑦ 那台仪器的既有暴露（`accept F8`）本批修不到**，也没测：它用**子串** `/dist/` 跳，
   与 `.gitignore` 与 index 都无关。归它的主人。
10. **submodule / sparse-checkout / worktree-linked `.git`（`git dir` 是文件不是目录）三支没测**：
    我只测了"目录里有 `.git` 且 root 就是仓库顶"与"根本没有"。`rev-parse --show-prefix`
    在 linked worktree 里同样返回空（我**没实测**这一句——它是推理，登记在此）。
11. **大小写/盘符/长路径/UNC 未测**（`accept §10.6` 同条未补）。`git ls-files` 输出与 `filepath.Rel`
    的一致性我依赖既有实现（r1 已如此），没有新增用例。
12. **F4（note 的"N file(s) under M ignored director(ies)"量纲）本批**没修**：
    票面没点它，我改了 note 的**行数**（一行⇒最多三行）却没改那一行的措辞。
    ⇒ `shapes` 那发是现量证据：修后 note 报 `1 file(s) under 1 ignored director(ies)`，
    而 `frontend/dist/assets/` 里实际藏着 1 枚从未被 visit 的文件——**F4 仍活着，账不变**。
    `A211`/`A214` 的台账登记**我没做**（不在我的 pathspec 里：`docs/reports/**` 归编排者），
    需要落的文字见本件 §1 与 §11，请编排者按 `A##` 追加。
13. **工具调用被拒**：本程**没有任一次**被权限系统拒绝，因此也没有"被拒后绕行"可报。
14. **Git 纪律自查**：只 commit、不 push；每次 `git commit -q -F - -- <显式 pathspec>`；
    提交前 `git diff --cached --name-only` 逐枚现量（只出现我这四枚路径）；
    没有 `add -A`/`add .`/`--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；
    16 枚 ` D design/**` 与两枚 `?? design/` **一枚都没被我带走**；临时件只建不删。

---

## 13. 总判（实现者自判，不是裁决）

**改了什么**：`skip()` 的判据从"模式命中"换成"模式命中 **且** git 的 index 不持有它"，
index 问不到就**一条规则都不应用**并打印一行点名原因；`note()` 增至最多三行（新增"规则未应用"与
"受追踪且模式命中的名单"）；两枚少扫方向的语义偏差（行首空白被 TrimSpace、尾随 `**` 吃掉父目录）
各自改正并配用例；第二副本 `ignoredLikeGit` 与其唯一调用点**同批**接上同样两条参数；
`main.go:61-64` 与 `gitignore.go:33-34` 那两句担保，一句**兑现并指向新断言**、一句**删掉换成窄命题**。
新增两枚测试，删除零枚断言、零枚 Skip。

**F1 两发复现读数**（同一棵 `shapes`，两版二进制）：
`frontend/dist/tracked-forced.tsx`(含 `≤`) 与 `internal/build/leak.go`(裸 `go worker()`) ——
**修前 rc=0 / `ban #6 frontend/=40` / `bans #1-5 internal/=203` / clean**；
**修后 rc=1 / `41` / `204` / 两枚 finding 点名**，note 另加一行
`2 path(s) git reports as TRACKED and matching an ignore rule … none of them was skipped: <两枚路径>`。
CI 那一侧（`ci-shapes`，无 `.git`、场内同样两枚）：修前 rc=0/40/203 ⇒ 修后 rc=1/41/204 且**第一行**就是
`gitignore rules NOT APPLIED - … not a git repository …`。

**两枚不蒙眼对照**：受追踪非忽略的 `frontend/src/still-bites.tsx`(`≤`) ⇒ **rc=1 点名 `scope frontend/`、`frontend/=41`**；
未追踪且被忽略的 `frontend/dist/untracked-only.html`(`≤`) ⇒ **不进分母、不进 finding 名单**（41 而非 42，
note 里 `skipped as git-ignored: 1 file(s)` 就是它）。单元层同一对形状由
`TestTrackedPathsAreNeverSkippedByTheIgnoreFilter` 以**两棵字节相同的树 + 分母 +1** 钉住。

**净快照八数（票面 #4）**：`bans #1-5 internal/=203 cmd/=22 / ban #6 frontend/=40 / ban #7 internal/tools/=18 /
ban #8 design/=16 frontend/=40 internal/=405 cmd/=39`，rc=0 —— **修前修后逐枚相同，没有一枚需要解释**。
门的端到端：step 1 `PASS=28 FAIL=0 SKIP=0 / === RUN=68`（修前 26/0/0 与 66），step 2 rc=0，名单只增两行无消失。
契约轴：`emojiRe` 那一行 sha 前缀 `a6150ef689cc20c9` 在 `3bb99aa^`/`3bb99aa`/HEAD/工作树四处同值（与验收件记录一致），
`allowlist.txt` 0 字节差异，diff 里无一行触及 `emojiRe|1F000`，无断言删除、无新增 Skip、无阈值/golden 改动。

**明确没测**：`-root` 落在别人仓库子目录那一支、PATH 里没有 git 那一支、10s 超时那一支、
`.git/info/exclude` 与 `core.excludesFile` 两支（含它们与新诊断行的新耦合）、
linked worktree / submodule / sparse-checkout、大小写与长路径、`/foo` 锚定仍无测试、
F4（note 量纲）原样活着、⑦ `internal/panel/frontend_hygiene_test.go` 与 root module 全套今天是否仍绿。
**F9（`ban #8 design/` 走 0 枚将硬退出 2）一字未动**，处置权在前端会话与编排者。
**本件的合上与否由非实现者裁**；若裁决者认为 §11.3 那枚超出票面的实现选择不可接受，
按 §11.3 给出的那一行删除路径退回即可。
