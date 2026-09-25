# d22scan gitignore 修复批（`ca84b75`）· 非实现者对抗验收 r1

**我是谁**：裁决者（非实现者）。本件不修任何东西——只判、只把最小修法点名交回。
**被验交付**：`ca84b75`（`fix(d22scan,A214/F1): skip() 改问 git index`）
＝ `tools/d22scan/gitignore.go` **＋304/−33** ／ `main.go` **＋44/−22** ／ `scan_test.go` **＋340/−2**
／ `docs/evidence/s1/d22scan-gitignore-fix-r1.md` ＋478/0（`git show --numstat ca84b75` 现量，与票面逐枚同值）。
**上一程的退回来龙去脉**：`docs/evidence/s1/d22scan-gitignore-accept-r1.md`（712 行，判 `3bb99aa` 整批退回，F1＋F2）。
**它点名的最小修法（`.git` 存在性加闸）已被编排者否决**（`A214②`），本批判的是**换掉的那一味**（index-aware），
不是"有没有照上一程的话做"——上一程的修法按盘上事实治不到它自己发现的那一发，这句我复核后**认**（见 §3.4）。

| 用途 | 锚点 | 说明 |
|---|---|---|
| 交付锚点（被验代码） | `ca84b75` | 其父 `885f50b`；**`git log ca84b75..HEAD -- tools/d22scan` 现量空** ⇒ 本件全部 `file:line` 对 HEAD 仍成立 |
| 我的读数锚点 | `5d463bb` → `77e89f9`（写本件时 HEAD 在动，每节标各自时刻） | §0/§1/§2 在 `5d463bb`；§6/§9 在 `77e89f9`；**HEAD 每节都现查**，不假定单调 |
| 实现程自述（**只当声明，不当证据**） | `docs/evidence/s1/d22scan-gitignore-fix-r1.md`（568 行） | 它的 §11.3 就是本件第 4 格要重量的"摘掉这一味"那发 |

**仓库外工作区**（全部临时件在此，真树未写一字节；只建不删）：`D:\tmp\d22fix-accept-r1\`
- `pre\` = `git archive ca84b75^`（修前代码）、`head\` = `git archive 5d463bb`（＝HEAD 代码，`tools/d22scan` 与 `ca84b75` 逐字同）
- `bin\d22scan-pre.exe` / `bin\d22scan-post.exe` = 两版代码各建一枚二进制，**同一棵树、同一串命令行**对照
- `f1repo\` = §1 的眼罩形状（仓库外 `git init` ＋ `git add -f` 两枚受追踪违规件）
- `f2repo\` = §2 的对照组（同样两枚件，但**未** force-add ⇒ 未追踪 且 被忽略）
- `f3ci\` / `f3repo\` = §3 的"同字节、只差一枚 `.git`"对照（含 4 枚构建产物）
- `nested\wisp\` = §3.3 的"`-root` 落在别人仓库的子目录"分支
- `shim\git.exe` = §3.3 的 15 s 慢 git（打 10 s 超时那一支）
- `logger\git.exe` = §4.3 的记账 shim（数一次运行 spawn 了几枚 git）
- `mut\` = `head\tools\d22scan` 的副本 ＋ 我的 `zz_accept_probe_test.go`（六发变异，§4/§5）
- 日志留档：`run-f1-pre.log` `run-f1-post.log` `run-f2-pre.log` `run-f2-post.log` `run-f3ci.log` `run-f3repo.log`
  `run-f2-nogit.log` `run-f1-timeout.log` `run-nested.log` `run-f3exclude2.log` `gate-worktree.log` `gate-snapshot.log`
  `step1-nogit.log` `scan-spawns.log` `census.txt`

**两把坏尺的规避**（上一程的教训，本件全程执行）：①台账 `A213④` 那条 `'"20[0-9][0-1][0-9][0-3][0-9]'` 恒不匹配——本件不复用它；
②"0 命中"必先打正控——§6 的名册普查与 §8 的"git 不在 PATH"两发**都先跑了已知存在的正控**，逐处附在正文里。

---

## 0. 交付物核对（`77e89f9` 现跑）

- `git show --numstat ca84b75` ＝ 上面那四枚，**只有这四枚路径**；
  `git diff --name-only ca84b75^..ca84b75` 现量 4 行，**`internal/panel/**` 一字节未动**（那是另一位的地界，§6 再核一次）。
- 修复程自己的三枚后续取证 commit（`ca06fcd`/`b6f596d`/`004c6ec`）只碰它自己那份证据件 ⇒ 无顺手改码。
- 票面背景核过：`A207`（缺陷本体）/`A211`④（点名"被追踪的文件永远不会被 git 忽略"）/`A213`/`A214`（退回＋否决上一程修法）
  /`A216`/`A218`（本批交件的编排者读数）逐枚存在于 `docs/reports/pending-and-issues.md:5842,5876,5896,5906,5920,5940`。

---

## 1. 主攻击格：F1 到底关没关——**复现，不推理**

### 1.1 形状（仓库外重建，两枚种子同时在场）

`f1repo\` = `git archive` 解出的完整树 → 仓库外 `git init -b main` → `git add -A` → 再放两枚文件并 `git add -f`：

| 种子 | 内容 | 命中的规则 | 藏的是哪一枚禁令 |
|---|---|---|---|
| `frontend/dist/tracked-forced.tsx` | `export const doc = "ready ≤";`（`≤`＝U+2264） | `frontend/.gitignore:12` `dist/*` | ban #8（字形） |
| `internal/build/leak.go` | `func leak() { go worker() }` ＋ 同串 `≤` | 根 `.gitignore:14` **不锚定**的 `build/` | **ban #1 裸 goroutine＝AGENTS.md §1.2 的安全禁令** |

git 自己的口径先钉住（现跑）：

```
$ git ls-files -i -c --exclude-standard
frontend/dist/tracked-forced.tsx
internal/build/leak.go
$ git check-ignore -v frontend/dist/tracked-forced.tsx      -> rc=1（受追踪 ⇒ git 说不忽略）
$ git check-ignore --no-index -v ...                        -> frontend/.gitignore:12:dist/* / .gitignore:14:build/  rc=0
```

### 1.2 同一棵树、两版代码（`-root` 指过去，其余一字不差）

**修前 `bin\d22scan-pre.exe`（代码＝`ca84b75^`）**

```
rc=0
d22scan: skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [internal/build/],
         decided by .gitignore (1 path(s)), frontend/.gitignore (1 path(s))
d22scan: scope bans #1-5 internal/      examined 203 production Go files
d22scan: scope ban #6 frontend/         examined  40 text files
d22scan: scope ban #8 internal/         examined 406 Go files ...
d22scan: clean - no D22 ban violations; ...
```

**修后 `bin\d22scan-post.exe`（代码＝`ca84b75`）**

```
rc=1
d22scan: 2 path(s) git reports as TRACKED and matching an ignore rule
         (git ls-files -i -c --exclude-standard) - a tracked path is part of the delivery,
         so none of them was skipped: frontend/dist/tracked-forced.tsx, internal/build/leak.go
d22scan: examined 226 production Go files ...
d22scan: scope bans #1-5 internal/      examined 204 production Go files
d22scan: scope ban #6 frontend/         examined  41 text files
d22scan: scope ban #8 frontend/         examined  41 text files
d22scan: scope ban #8 internal/         examined 407 Go files ...
internal/build/leak.go:7: [bare-goroutine] bare `go worker(...)` is banned (D22/D38b, R16: named calls count too)
frontend/dist/tracked-forced.tsx:1: [emoji] ban #8 glyph in scope frontend/ is banned (D23)
internal/build/leak.go:3: [emoji] ban #8 glyph in scope internal/ is banned (D23)
d22scan: 3 finding(s); D22 bans are not negotiable
```

### 1.3 判

**承重的那一发——`internal/build/leak.go` 的裸 `go worker()`——从"门看不见"变成"门点名"**：
修前它连分母都不进（`internal/`=203 而非 204）、`rc=0`、note 还把整包被剪报成 "1 file(s)"；
修后 bans #1-5 与 ban #8 同时红、分母各 ＋1、note 点名两枚路径。字形一发同样从绿转红。
⇒ 上一程 §6.9 那发"同一形状打到 Go scope"是真关掉了，**且不是靠收窄射程**（`build/` 规则一字未改，
`frontend/.gitignore:12` 一字未改，改的是判据）。
代码落点（现号）：`skip()` 先问 index 再问规则 `gitignore.go:338-376`（`rel==""` 在 `:343-349` 早退、
`!ix.ok` 在 `:350-356` 早退、`holds()` 在 `:357`），`holds()` 本体 `:184-192`，
`runGitIndex()` `:202-239`（`rev-parse --show-prefix` `:206` → `ls-files -z` `:217` → `ls-files -z -i -c --exclude-standard` `:221`）。
**判：成立。**

---

## 2. 不是换方向的眼罩／不是伪装的退回（A207 那一头还在不在）

`f2repo\` ＝ 同两枚种子文件**同样在盘上**，但**没有** `git add -f`（`git ls-files -i -c` 现量空 ⇒ 未追踪 且 被忽略），
`git ls-files` 里查无此二件（现量空）。同一棵树、同一串命令行：

| 读数 | 修前 `467f8a4` 代码 | 修后 `ca84b75` 代码 |
|---|---|---|
| rc | 0 | 0 |
| `skipped as git-ignored` 行 | 有（1 file + `[internal/build/]`） | **逐字同值** |
| `ban #6/#8 frontend/` | 40 / 40 | 40 / 40 |
| `bans #1-5 internal/` | 203 | 203 |
| `ban #8 internal/` / `cmd/` | 406 / 39 | 406 / 39 |

`diff run-f2-pre.log run-f2-post.log` **现量 0 行差异**（两枚文件在盘上、内容含 `≤` 与裸 goroutine，
两版都不算进分母、都不读）。⇒ **"真被忽略且未追踪的件仍然不算"这一头完好＝不是退回**；
而 §1 那一头（受追踪且模式命中）单独变红 ⇒ **也不是往另一个方向蒙眼**。两读并立，判：**成立**。

---

## 3. "问不到"那一支：够不够响、是不是永远朝"多扫"倒

### 3.1 `git archive` 快照（没有 `.git`）

`f3ci\` ＝ 净快照 ＋ 4 枚 `frontend/dist/` 构建产物（`index.html`、`assets/b1.js`、`assets/b2.js`、含 `≤` 的 `probe.txt`）。
`f3repo\` ＝ **同一串字节** ＋ `git init` ＋ `git add -A`（那 4 枚因此未追踪且被忽略）。

| | `f3ci`（无 `.git`） | `f3repo`（有 `.git`） |
|---|---|---|
| 输出的**第 1 行** | `d22scan: gitignore rules NOT APPLIED - git cannot be consulted in …: `git rev-parse --show-prefix`: fatal: not a git repository …` | `d22scan: skipped as git-ignored: 2 file(s) …` |
| `ban #6`／`ban #8 frontend/` | **44 / 44**（40 ＋ 4 枚构建产物全算） | 40 / 40 |
| `skipped as git-ignored` | 现量 **0 次**（`grep -c` ⇒ 0） | 1 次 |
| rc | **1**（`frontend/dist/probe.txt:1 [emoji]`） | 0 |

⇒ **第一行就是自陈、分母随后包含构建产物**（"多扫"不是形容词，是 44 与 40 之差），
note 打在 **stdout**、和账本同一枚流（现测：`2>/dev/null` 后仍在，`1>/dev/null` 后不在）⇒ CI 日志里不会被拆到两个流去。
编排者 `A218⑦` 那发（净快照第一行＝NOT APPLIED、rc=0）我在 `77e89f9` 时刻的 `head\` 上独立复现（§9）。

### 3.2 **PATH 里没有 git 而 `.git` 在**（实现程自陈没测的那一支——本件测了）

```
$ env PATH='D:\tmp\...\shim;C:\Windows\System32;C:\Windows' bin/d22scan-post.exe -root f2repo
rc=1
d22scan: gitignore rules NOT APPLIED - git cannot be consulted in …f2repo:
         `git rev-parse --show-prefix`: exec: "git": executable file not found in %PATH%; …
d22scan: scope ban #6 frontend/  41 text files   （有 git 时是 40）
d22scan: scope ban #8 internal/  407             （有 git 时是 406）
d22scan: 3 finding(s)  -> internal/build/leak.go:7 [bare-goroutine] + 两枚 [emoji]
```

同一棵树**有 git 时**（§2）：`rc=0`、40、无 finding。⇒ 这一支**响亮且朝多扫**，
`exec: "git": …` 那枚 Go 层的错也被原样进了自陈行（`runGitAt` `gitignore.go:244-263` ⇒ `fail(...)` `:203-205`）。

### 3.3 另外两支（超时／子目录）

- **10 s 超时**：`shim\git.exe` 每发睡 15 s。`-root f1repo` ⇒ 总耗时 `real 20.7s`（＝两枚 matcher × 10 s ⇒ 顺带证明
  **每次运行恰好 2 枚 matcher**），自陈行是 `timed out after 10s running `git rev-parse --show-prefix``，
  rc 仍由发现决定（＝1，三枚点名）。超时用 `context.WithTimeout`（`:245`），**不是墙钟差**（AGENTS.md §1.2 那一支合规）。
- **`-root` 落在别人仓库的子目录**：`nested\wisp\`（外层有 `.git`，`git rev-parse --show-prefix` 现量 `wisp/`）
  ⇒ 自陈 `… is the subdirectory "wisp/" of a git repository, not its top, so that index is not this tree's index`、规则全关、不剪任何件。
  这一支是**上一程修法（只看 `.git` 在不在）会判错**的形状（`.git` 不在 `wisp/` 下但在 `nested/` 下），本批用 prefix 判据堵上了。

### 3.4 判（含对上一程修法的一句回执）

四支"问不到"（无 `.git`／无 git 二进制／超时／子目录）**逐支：一行自陈 ＋ 规则全关 ＋ 分母只增不减**，
未发现任何一支能安静变绿。上一程点名的 `.git` 存在性加闸确实治不到 `A214②` 说的那一支
（真检出永远带 `.git`，那扇闸只在 `git archive` 副本上生效）——**编排者否决它、改问 index，判据方向正确**。
**判：成立。**

---

## 4. 那味偏离（编排者已批准保留）：摘掉它，有没有一发变异从此打不红

票面点名的是 `git ls-files -i -c --exclude-standard`（受追踪 ∩ 撞规则）。实现**额外**问了全量 `git ls-files -z`
（`gitignore.go:217`，理由写在 `:166-171`：窄清单挡不住"匹配器比 git 多跳"那一类）。
本仓的承重定义＝**摘掉一味，是否存在一发变异从此打不红**。我用六发变异量（全部在 `mut\`＝`head\tools\d22scan` 的副本里跑，
被测代码逐字等于交付版；`TestLedgerCountsMatchAnIndependentWalk` 等三枚真树守卫在我副本的相对路径上会 `t.Skipf`，
下面凡是涉及它们的读数我改到 `f1repo\tools\d22scan` 里跑，那里 `../..` 就是一枚真仓库）。

| 发 | 摘掉/改坏的东西 | 改动（一行） | 读数 |
|---|---|---|---|
| **M1** | **只摘"全量清单"这一味**：`parseIndexPaths(all)` → `parseIndexPaths(withIndex)` | `gitignore.go:225` | **全绿**：`ok … 1.578s`——交付的六枚相关测试一枚都不红，我的探针也不红 |
| **M2** | `holds()` 整个短路（`if false && ix.holds(...)`） | `:357` | **红**，8 条：`scan_test.go:1479/1482/1485/1488`（四枚分母 ＋1 全灭）、`:1508`×3（三枚点名"did NOT fire at the tracked path"）、`:1534` `verdict on the tracked half must be 1 (findings), got 0` |
| **M3** | `if !ix.ok` 短路（问不到照样套规则＝上一程的形状） | `:351` | **红**，8 条，含 `:1594` 与 `:1597`（note 变成两行、rc 从 1 掉到 0） |
| **M4a** | 行首空白那枚语义修回退（`TrimSpace`） | `:557` | **红**：`scan_test.go:1686` 两支（`lead/inside.txt` 应为不忽略、`"  lead/inside.txt"` 应为忽略） |
| **M4b** | 尾随 `**` 那枚语义修回退（`start = 1` → `0`） | `:600` | **红**：`decide("deep", isDir=true)` 被吞、`deep/.gitkeep` 跟着没 |
| **M5** | 只摘**目录**那一半（`holds()` 的 `isDir` 支 → `false`） | `:189` | **红**：`internal/build/leak.go` 又隐形（分母 ＋1 灭、两枚点名灭） |
| **M6** | 只摘**窄清单**的一半（`ignoreSet` 不再进 note） | `:237-238` | **红**：`scan_test.go:1527` `note "" must name "TRACKED and matching an ignore rule"` |

**M1 与 M4a 的交叉发（这一发才是这一味的真正判据）**——我在 `mut\zz_accept_probe_test.go` 造了
"匹配器比 git 多跳"的形状：规则 `frontend/.gitignore` = `  weird/*`（**行首两空格是模式数据**，
拿真 git 现量钉过：`git add -A` 之后 `weird/a.txt` 受追踪、`git ls-files -i -c` 现量空、
`git check-ignore -v weird/a.txt` ⇒ rc=1；而一枚**真名叫 `  weird` 的目录**才 rc=0 被忽略），
种子 `frontend/weird/inside.tsx` 含 `≤`：

| 配置 | 探针读数 |
|---|---|
| 交付原样 | PASS：`ban #6 examined=3`、finding 点名 `frontend/weird/inside.tsx` |
| **只 M4a**（语义修回退，**全量清单还在**） | **仍 PASS**——`holds()` 的全量 tracked 支把它救回来了（分母 3、finding 在） |
| **M4a ＋ M1**（语义修回退，**只剩窄清单**） | **红**：`ban #6 examined=2`、`findings=[]`、note 说 `skipped as git-ignored: 1 file(s)`——**一枚受追踪的交付字节对门永久隐身，而这正是 git 亲手说"我没忽略它"的那一支** |
| 只 M1（语义修完好） | PASS（探针与交付名册都不红） |

**判：这一味不是装饰，是防御——但它防的那一发今天没有任何一枚断言钉着。**
- 摘掉它**今天**不打红任何一发已交付的变异（M1 全绿），按这仓的承重定义**它是未被测试要求的代码**；
- 可它唯一的作用恰好是"窄清单天生看不见的那一类"（`-i -c` 是 `ls-files` 的**真子集**，现量 `comm -13` 空 ⇒
  它永远举不出"git 不忽略而匹配器忽略"的件），而那一类**这一批刚修过两枚实例**（F3 的行首空白与尾随 `**`）
  ⇒ 它是**活家族的保险**，不是想象中的保险。M4a 那行一旦将来被谁改回去，只有它挡着（第三行读数）。
- **代价核过**：一次运行 spawn **恰好 6 枚** git（`logger\git.exe` 记账现量：3 枚一组、两组，
  cwd 全在 root，命令＝`rev-parse --show-prefix`／`ls-files -z`／`ls-files -z -i -c --exclude-standard`），
  来源是 **2 枚 matcher**（`checkRoot` 自建一枚 `main.go:1160` ＋ scan 一枚），每枚懒加载后缓存（`gitignore.go:382-387`）；
  对 1144 枚受追踪件的树，**同树同机对照实测 修后 0.76s vs 修前 0.42s ⇒ ＋0.34 s/次**（编排者票面的"~0.5 s"同一量级）；
  最坏一支是 git 挂住：实测 `real 20.7s`（10 s × 2 枚 matcher，§3.3）。
  ⇒ **只付在 `scripts/d22scan.sh` step 2 与每次本地跑门**（整道门实测 13.9s，step 1 的 go test 才是大头），
  CI 的 lint job（`.github/workflows/ci.yml:66` `ubuntu-latest`，step 见 `:81`/`:109`）不受影响。
- **它没有让任何既有担保重新变得不可证伪**：窄清单那味仍被 `:1527`（M6 红）钉着，全量清单的**目录**那半被 `:1485/1488`（M5 红）钉着。
  没被钉住的只有全量清单的**文件**那半（M1 全绿）。

**最小修法（点名，不动手）**：把 M4a＋M1 那一发做成第三枚半——在
`TestTrackedPathsAreNeverSkippedByTheIgnoreFilter` 里加一个 seed 分支（`frontend/.gitignore` 写 `  weird/*`、
种子 `frontend/weird/inside.tsx` 含 `≤`、只跑普通 `git add -A`，并**断言 `git ls-files -i -c` 现量为空**
＋**断言 finding 点名该件**）。约 15 行测试、零生产码改动；从此摘掉 `all` 那一味必红。
**格判：附条件成立**——保留是对的（防御、不是装饰），条件是那枚防御今天无断言可及，须按上面一发补上。

---

## 5. F2 那两句担保，现在证伪得动吗

| 句子（现号） | 原文 | 钉它的断言 | 变异证据 |
|---|---|---|---|
| `main.go:64-71` | "it never suppresses a finding on a TRACKED path. skip() asks `git ls-files -z` and `git ls-files -i -c --exclude-standard` and refuses to skip anything either one reports…" | `scan_test.go:1439 TestTrackedPathsAreNeverSkippedByTheIgnoreFilter`：分母增量 `:1477-1492`、findings `:1494-1512`、provenance `:1519-1529`、退出码 `:1531-1537` | **M2 红**（8 条，含 `:1534` 那句 `got 0`）；上一程的 F2 说"它点名的测试是 `!` 否定捞回的 `.gitkeep`、夹具无 index"——本批同一枚夹具真 `git init` ＋ `git add -f`（`:1411 gitIndexFixture`、`:1421 gitForceAdd`），形状表达能力补上了 |
| `gitignore.go:44-47` | "…is the narrow, testable one: **for a path git reports as tracked, skip() returns false.**" | 同一枚测试（`skip()` 是四枚 walk 的唯一入口，`main.go:648/662/841/853/919/952/1168/1173` 现量），祖先目录那一半另由 `TestWalksSkipGitIgnoredPaths`（`:1279`）与 M5 钉 | **M2 红 ＋ M5 红**；句子自陈的 pinner 存在且咬得住 |

⇒ 上一程点名的两枚未钉担保：**一句兑现（M2 可红）、一句从"过度声明"换成窄命题并被同一枚测试钉住**（`cannot disagree` 那句已被删除，
`git diff ca84b75^..ca84b75 -- gitignore.go` 的删除行 #2 逐字就是它）。**这两枚：成立。**

**但本批自己又新写了一句同族的、未被任何断言钉住的 blanket**（这是文本面的新账，不是旧账）：

- `gitignore.go:103-105`：`The one shape that could make it skip MORE than git - not being able to reach the index at all -
  is handled by the rule above it, not by this list`。**这句是 ca84b75 新增行**（`git diff … | grep -n '^\+.*The one shape'` ⇒ 命中）。
  它的"唯一"是**假**的：§4 第四发的 M4a（匹配器把自己的规则读宽）就是**第二种** skip-MORE 形状，
  而它**不是**被"一条规则都不应用"治住的，是被 `holds()` 的全量 tracked 支治住的
  （M4a 单独跑仍绿 ⇒ 今天没有匹配器多跳；M4a＋M1 ⇒ 隐身）。
  **它也不可证伪**：交付名册里没有任何断言能区分这句的真假（M1 全绿就是这件事的形状）。
  最小修法＝把 §4 那发补上，并把这句改成"…is handled by `holds()` 的全量清单 ＋ 问不到时不套规则"两支并列。
- `gitignore.go:166-171`：`an independent union member so one command's parsing bug cannot blind the other`——**窄清单在 `holds()` 的文件支里是空转的**：
  `ls-files -i -c` 严格 ⊆ `ls-files`（现量 `comm -13` 空），且两枚都过同一枚 `parseIndexPathsList`（`:279` 由 `:301` 调用），
  所以"一枚命令的解析 bug 盲不了另一枚"这个理由在代码里不成立（另一枚也瞎）。
  窄清单真正承重的作用是 **note 点名**（M6 红钉着）与测试里的 git 侧前提（`:1465`）。⇒ 同一枚一句话级的修法。
- 顺带两句**属于成本陈述、不属于安全担保**的未钉文字，记下不升级为缺陷：
  `gitignore.go:346-348`（`rel==""` 早退"costs no git call"）、`:199-201`（"any failure yields ok=false plus one reason"，
  其中 `parseIndexPathsList` 对绝对路径/`..` 的拒读没有任何测试喂过畸形流——方向是安全的：拒读⇒不套规则⇒多扫）。

**格判：附条件成立**——F2 点名的两句都从"担保"变成"有 pinner 的窄命题"（M2/M5 各咬一次），
但同一批新写了一句"唯一形状"式 blanket（`:103-105`）与一句站不住的"独立并集成员"理由（`:166-171`），
两枚各一行文字修法、并共用 §4 那发补测试。

---

## 6. 同批改副本这条规矩（本仓栽过两次的那条）

1. **副本确实动了，而且在同一枚 commit 里**：`ignoredLikeGit` 的签名从 `(rel string)` 变成
   `(rel string, tracked map[string]bool, rulesOn bool)`（`scan_test.go:1758`，现量 `git diff ca84b75^..ca84b75 -- scan_test.go`
   删除行只有两枚，逐字是 `if ignoredLikeGit(relToRepo(root, p)) {` 与旧签名行——**没有第三枚手抄新造**）。
   配套新增 `gitTrackedIndex`（`:1775`，**手抄的第二份 git 问答**，不调被测码）。
2. **副本这次真有牙（实测，不是推理）**：把副本那一半的 index 摘掉（只在 `f1repo-mut\` 的测试里改，
   `ignoredLikeGit(relToRepo(root, p), nil, indexOK)` ＋ `_ = tracked`），**同一棵 F1 树**上：

   ```
   --- FAIL: TestLedgerCountsMatchAnIndependentWalk
       scan_test.go:1255: scope ban #6 frontend/ reported examining 41 files,
                          independent walk says 40 - the number in the self-report is wrong, not just small
   ```

   正控（原样）在同一棵树上 PASS 且 oracle 自己报 `verified ban #6 frontend/: 41`、`bans #1-5 internal/: 204`、
   `ban #8 internal/: 407` ⇒ 两边在 F1 形状上**合流到同一个 41**，
   而一旦副本退回"纯模式"就当场分家——**上一程 §1.4 那句"两份副本共用同一条纯模式前提"被这一发改掉了**。
3. **名册普查（票面命令我自己跑的）**：命中 113 行（正控：`grep -rn 'ignoredLikeGit' --include='*.go' tools` 先打出 3 行，
   证明这把尺不瞎）。按**位置**数，真正的"可执行跳过逻辑"仍是

   | # | 位置 | 抄的是什么 | 本批 |
   |---|---|---|---|
   | ①②③④ | `main.go:648/841/919/1168`（目录名）＋ `:662/853/952/1173`（文件） | `testdata`/`node_modules`/`.git` 目录名 ＋ 全部接同一枚 `ign.skip` | ①②③④ 的 ignore 那一支语义变了（判据加 index），**代码只有一枚 `skip()`** |
   | ⑤ | `scan_test.go:1213` 独立对照走查 | 目录名 ＋ `ignoredLikeGit` | **同批改了** |
   | ⑥ | `scan_test.go:1758 ignoredLikeGit` | 手抄的 `frontend/dist/` ＋ `.gitkeep` 政策 ＋ **现在加上了 index/rulesOn** | **同批改了**（不是新造） |
   | ⑦⑧⑨ | `internal/panel/frontend_hygiene_test.go:289`、`internal/panel/composer_test.go:286/434/606`、`internal/risk/pathresolver_rewrite_account_test.go:223` | 子串 `/dist/`、case 表列名 | **一字未动（正确）**，且 `A217④` 已把 `internal/panel` 那枚新仪器记在另一程名下 |

   ⇒ **没有造出新的不受控副本**；`internal/panel/**` 那条既有暴露（上一程 F8）仍归它的主人，本批没碰也没该碰。
4. **地界**：`git diff --name-only ca84b75^..ca84b75` ＝ 4 枚路径，`internal/`、`internal/panel/`、`cmd/`、`frontend/`、`design/`、
   `docs/reports/**` **零字节**（现量）。本件自己全程只写 `docs/evidence/s1/d22scan-gitignore-fix-accept-r1.md` 一枚路径。

**判：成立。**
