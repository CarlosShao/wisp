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

---

## 7. 契约轴（我自己重导的，不抄上一程）

| 判据 | 现量 | 判 |
|---|---|---|
| `emojiRe` 逐字节 | `main.go:164` 现串＝ ``var emojiRe = regexp.MustCompile(`[\x{1F000}-\x{1FAFF}\x{2200}-\x{22FF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}\x{FE0F}\x{1F1E6}-\x{1F1FF}]`)``；**五枚锚点（`3bb99aa^`/`3bb99aa`/`ca84b75^`/`ca84b75`/HEAD）抽该行做 sha 全同**（`c8fb5593539afd07`，len 136）；`git diff ca84b75^..ca84b75 -- main.go \| grep -E '^[+-].*(emojiRe\|1F000)'` ⇒ **rc=1 零命中** | ✅（上一程 §5.1 的 sha 前缀与我不同＝它那侧管道不同，**"五锚同值"这一条两把尺独立同结论**） |
| `allowlist.txt` | `git diff ca84b75^..HEAD -- tools/d22scan/allowlist.txt` ＝ **0 字节** | ✅ 零新增豁免 |
| 断言有没有被删 | `scan_test.go` 里 `t.Error*/t.Fatal*` 计数 `ca84b75^=139 → ca84b75=164 → HEAD=164`；该文件删除行**逐字只有两枚**（旧 `ignoredLikeGit` 调用行＋旧签名行），**没有一枚断言、没有一枚 `func Test` 消失**（`grep -E '^-func Test'` 空） | ✅ 只加不减 |
| 新增 `t.Skip` | `t.Skip(f)(` 计数 **7 → 7 → 7** 三锚同值；`git show ca84b75 \| grep '^\+.*t\.Skip'` 只命中一枚**注释行**（`:1381` 那句"No t.Skip on this path, on purpose"） | ✅ |
| 阈值 / golden | `internal/observe/thresholds.go` 与 `*slo*` 路径 `git diff ca84b75^..HEAD` **0 字节**；`internal/agent/testdata/golden/*.sse` 根本不在被碰的 4 枚路径里 | ✅ |
| `main.go` 的 −22 | **22 枚全是注释行**（12 枚旧 DESCEND 段、6 枚 `ign` 字段注释、3 枚 `walkGo` 注释、1 枚 "…machine)." 被并入长句）——**没有一枚代码/分母行为被摘掉**；`walkGo/walkText/walkEmoji/checkRoot` 四枚 `skip` 接入点（`:648/662/841/853/919/952/1168/1173`）一行未动 | ✅ 逐枚点名 |
| `gitignore.go` 的 −33 | 24 枚注释（含上一程 F2 点名的两句：删除行 #2 逐字就是 `the instrument and the repository's ignore rules cannot disagree`，#14 就是 `never less`）＋ 9 枚 `note()` 函数体行（`if g == nil \|\| (files==0 && pruned==0)` 的守卫挪进函数内、`dirs`/`src` 组装与 `Sprintf` 原样保留在新 `if g.files > 0 \|\| len(g.pruned) > 0` 块里）⇒ **note 的旧那一行文案一字未改**（`:421-422`），新增的是第 1、3 两行 | ✅ |
| 函数面 | `gitignore.go` 顶层 `func` 13 → 20（＋7：`holds`/`runGitIndex`/`runGitAt`/`firstLine`/`parseIndexPaths`/`parseIndexPathsList`/`index`）；**全部只读**（三枚命令我逐字核过：`rev-parse`／`ls-files`／`ls-files -i -c`，无任何写操作、无 `--force` 类参数） | ✅ |

---

## 8. 新环境依赖：门的正控现在要 `git`

**编排者已判"可接受"，我复核后同意结论、更正它给的理由。**

- 现量（真工作树 `2206cd4`，同一命令 `sh tools/d22scan/runtests.sh -C tools/d22scan ./...`）：

  | 环境 | 四数 | rc |
  |---|---|---|
  | PATH 有 git（`/mingw64/bin/git`，`git version 2.52.0.windows.1`） | **PASS=28 FAIL=0 SKIP=0 ＝ 28 枚 `^--- PASS` 我自己数过（正控）**，RUN=68 | 0 |
  | **PATH 无 git**（正控：`env PATH=… git --version` ⇒ `No such file or directory`；`go version` 仍在） | **PASS=25 FAIL=3 SKIP=0**，RUN=68 | **1** |

  三枚红：`TestWalksSkipGitIgnoredPaths`、`TestTrackedPathsAreNeverSkippedByTheIgnoreFilter`、
  `TestIgnoreRulesAreOffWhenTheIndexCannotBeAsked`。红句点名要什么：
  `scan_test.go:1450: git init -q -b main in …: exec: "git": executable file not found in %PATH% …
  - the tracked-path guarantees of A214 cannot be built without a git binary`。
- **没有一枚走 `t.Skip`**（SKIP=0，且 `gitCommand` `:1384-1402` 只有 `t.Fatalf`）⇒ 缺 git 不会被 `runtests.sh` 读成通过
  （票 71 AC#3 那条"SKIP 不是 pass"在这里真的用上了）。
- **实现程那句自述我逐字复核为真**："扫描器本体不要求 git，变的是门自己那一步"——
  §3.2 现量：把 git 从 PATH 摘掉，**二进制照跑、rc 只朝红、分母只涨不落**。
- ⚠ **CI 那一侧的更正**：编排者 `A218③` 的理由是"本机 runner 是 self-hosted Windows、本来就靠 git 克隆"——
  **跑 `d22scan` 的那一枚 job 不在 self-hosted Windows 上**：`.github/workflows/ci.yml:65-66` `lint: runs-on: ubuntu-latest`，
  step `:81`＝正控、step `:109`＝`sh scripts/d22scan.sh`；self-hosted 那枚是 `ci.yml:538 [self-hosted, wisp-slo]`（`slo-full`）。
  ⇒ **结论仍成立但依据要换**：GitHub 托管的 `ubuntu-latest` 镜像必带 git，且 step `:68` 用的就是 `actions/checkout@v4`
  （它自己就得调 git 才能建出带 `.git` 的检出）——**这一步同时保证了 CI 侧扫的是一枚真仓库**，
  于是 §3.1 那枚"规则全关"的形状**在 CI 上不发生**（这是 `A214②` 否决 `.git` 闸门时说的同一件事，我在读数层复核到）。
  我**没有**跑过 GitHub CI（见 §12）。
- 方向判断：**红而不静＝对的方向**。三枚红的代价（缺 git 时门跑不动）换来的是"受追踪"这一判据可被证伪；
  反过来（缺 git 时静默少扫）正是这仓抓过多次的那枚假绿。**判：成立（含上面那句 runner 归属更正）。**
- 顺带一枚脆性（低，登记不升级为缺陷）：`TestIgnoreRulesAreOffWhenTheIndexCannotBeAsked` 把 git 的错误文案
  `not a git repository`（`scan_test.go:1588` 那一支）写进了断言。⇒ 若某台机器的 `TMPDIR` 落在一枚真仓库里面，
  `liveFixture` 的临时树会走"子目录"那一支、文案换成 `is the subdirectory … not its top` ⇒ **行为仍正确、测试却红**
  （FAIL 而非静默，方向安全）。本机与 ubuntu 都不属这形（本机现量 `C:/Users/swq/AppData/Local/Temp/…`，非仓库内）。

---

## 9. 读数：用**当下**基线，逐枚归因

**我的两种形状、两个二进制（同树同参数）**：

| scope | 真工作树（`77e89f9` 与 `2206cd4` 两时刻各一次，`sh scripts/d22scan.sh` step 2） | 净快照 `git archive HEAD`（`5d463bb` 一次、`2206cd4` 一次） |
|---|---|---|
| bans #1-5 internal/ | 203 | 203 |
| bans #1-5 cmd/ | 22 | 22 |
| ban #6 frontend/ | 40 | 40 |
| ban #7 internal/tools/ | 18 | 18 |
| ban #8 design/ | **32** | **30** |
| ban #8 frontend/ | 40 | 40 |
| ban #8 internal/ | **407** | **406**（@`5d463bb`）／**407**（@`2206cd4`） |
| ban #8 cmd/ | 39 | 39 |
| rc | 0（整道门含 step 1 也 0） | 0（step 2 与整道门都是 0） |

- 编排者在 `A218⑥` 给的基线 `bans #1-5 internal/=203 cmd/=22 / #6 frontend/=40 / #7 internal/tools/=18 /
  #8 design/=30 frontend/=40 internal/=406 cmd/=39`：**在 `5d463bb` 的净快照上我逐枚复现同值**（`gate-snapshot.log`），
  并复现到"输出的第一行＝NOT APPLIED 自陈"（`A218⑦` 那一发我独立量到同一行）。
- **`internal/=406 → 407` 是 `d88c356`（panel 会话新加 `internal/panel/l2_grant_boundary_test.go`）造成的**，
  不是本批：`find internal -name '*.go'` 两树差集现量**恰那一枚文件**；工作树 `407` 与净快照 `407` 现已同值。
- **`design/` 30（＝HEAD 里 16 枚旧件 ＋ `5d463bb` 入库的 14 枚 doubao）**、工作树 32（盘上另有 18 枚未追踪的
  `design/old/**` 文本 − 16 枚已删除未提交）⇒ `git ls-tree` 口径两枚现量：**HEAD=30 / disk=32**。
- **`ca84b75` 有没有动任何一枚读数？没有——这条我用两把尺各量一次**：
  ①**真工作树**上 `d22scan-pre.exe` 与 `d22scan-post.exe` 输出 **`diff` 现量 0 行**（八数、账本行、rc 全同，均 rc=0）；
  ②**净快照**（`2206cd4`）上两版 `diff` **唯一差异就是多出的第 1 行自陈**（`0a1`），八数一字不动。
  原因也钉住了：真仓 `git ls-files -i -c --exclude-standard` 现量**空**（今天没有一枚"受追踪且撞规则"的件），
  净快照里没有未追踪件可被规则吞 ⇒ 两版判据在这些树上必然同值。
  ⇒ **§1 那两发的分母只涨在"有人真去 `git add -f`"的树上，今天这仓没有这种树**——这既是好消息也是本批价值所在（洞堵上了，读数没动）。
- `git ls-files -i -c --exclude-standard` 在两枚锚点（我自己的会话开始与写本节时）都现量空；
  上一程 §1.2 用过的第二枚口径（逐件 `git check-ignore --no-index`）我**没重跑**（见 §12 第 6 条）。
- 正控与"没输出≠过了"：`scripts/d22scan.sh` 是 `set -eu`、step 1 就是正控——
  本节所有 step 2 读数我都是**分开取的**（step 2 单独用二进制、step 1 单独用 `runtests.sh`），
  两次整门运行（`gate-worktree.log`/`gate-snapshot.log`）都是 **step 1 绿（28/0/0/68）之后 step 2 才有输出**。
  另外缺 git 那一发恰好是一枚天然的反证：`step1 rc=1` 时 step 2 根本不该被期望有输出。
  **格判：成立。**

---

## 10. 往前账：哪些是本批的缺陷、哪些是登记过的兄弟

| # | 项 | 本批动了吗 | 我量到的现状 | 判 |
|---|---|---|---|---|
| R1 | **`live:true` scope 走 0 枚 ⇒ 硬退出 2**（上一程编号 **F9**；编排者 `A218⑨` 把它记成了 "`F3`"——**同一台账里 `F3` 另有所指**，见下面 R2） | 否（也不该由它动） | 现量 `ban #8 design/` 净快照 **30**、`git ls-tree HEAD design` 文本件 30；把那 16 枚未提交删除**假设入库**后仍剩 **14** 枚（`comm -23` 现量）⇒ **雷被 `5d463bb` 拆了，不是被本批拆的**；`!!`（被忽略的 design 件）现量 0 ⇒ ignore 过滤器仍然吞不到、也不该吞它 | **不是本批缺陷**；⚠ 台账那一处 `F3`/`F9` 同号相撞请编排者追一行 `>` 更正（append-only，不改写） |
| R2 | 上一程 **F3**＝两枚"少扫"语义（行首空白被 `TrimSpace`、尾随 `**` 吃掉父目录） | **是，本批正是来修它的** | M4a 红（`scan_test.go:1686` 两支）、M4b 红（`decide("deep",isDir=true)` 被吞）⇒ **两枚都实现且有断言**；且拿真 git 复核过方向（`  weird/*` 只绑真名叫 `  weird` 的目录、`git check-ignore -v dist` 对 `dist/**` 的实测） | **成立（旧账已结）** |
| R3 | `/foo` 前导斜杠锚定**无断言** | 否 | `TestGitIgnoreRuleSemantics` 的规则表（`:1626-1641`）逐枚看过：**没有一枚以 `/` 开头**；实现支在 `gitignore.go:570-572` | **登记过的兄弟账，仍在**（本批没让它变坏，也没顺手补）——最小修一：表里加 `"/root-only.log"` ＋ 两支用例 |
| R4 | `.git/info/exclude` ＋ `core.excludesFile` 不作为跳过来源 | 否（`gitignore.go:97-102` 明示不支持） | 现测：往 `f3repo/.git/info/exclude` 写 `excluded-probe.tsx` ⇒ `git check-ignore -v` 点名 `.git/info/exclude:8`（rc=0，git 认它被忽略），**扫描器仍算进分母（41）并点名为 finding** ⇒ 方向＝**少跳＝响亮**；代价＝这一支上"工作树 41 / CI 40"的 A207 病仍活着 | **登记过的兄弟账（上一程 F5），方向安全** |
| R5 | 上一程 **F4**＝`note()` 的数与移动的分母不同量纲 | 否（模板 `:421-422` 一字未改） | 本批新独立复现一发：`f3ci`(规则关) 44 → `f3repo` 40，note 说 **"2 file(s) under 1 ignored director(ies)"** ⇒ 动掉 4 枚只报了 2 枚 | **登记过的兄弟账，本批未恶化也未修**；票面没点它，不算越界 |
| R6 | 上一程 **F8**＝`internal/panel/frontend_hygiene_test.go:289` 用子串 `/dist/` 跳，早于本修且比 `.gitignore` 宽 | **否，正确地没碰**（`git diff --name-only ca84b75^..ca84b75` 不含 `internal/**`） | 该行现号仍在；归它的主人（`A217④` 已把新仪器记在 panel 名下） | **不是本批的账** |
| **N1（新）** | `gitignore.go:103-105` "The one shape that could make it skip MORE than git" 是本批**新增行**，"唯一"为假、且把治它的那一味归错了地方（真治它的是全量 tracked 支，不是"问不到就不套规则"） | **是，本批写的** | §4 的 M4a＋M1 交叉发＝它的反例；交付名册里没有任何断言能证伪这句 | **本批的文本面缺陷（低-中）**，修法见 §5（一句文字＋§4 那发测试） |
| **N2（新）** | `gitignore.go:166-171` "an independent union member so one command's parsing bug cannot blind the other" 站不住 | **是，本批写的** | `-i -c` 严格 ⊆ `ls-files`（现量 `comm -13` 为空）⇒ `holds()` 文件支里那一味**空转**；两枚清单过同一枚 `parseIndexPathsList`（`:279`→`:301`）⇒ 解析 bug 同打两枚 | **本批的文本面缺陷（低）**；窄清单真正承重的是 note（M6 红），把理由改成这句即可 |
| **N3（新）** | 全量清单的**文件**那一半没有任何断言要求（M1 全绿） | 本批引入的代码，非回归 | §4 表 | **附条件**（§4 的最小修一即解） |
| **N4（新，信息）** | `gitignore.go:346-348` "costs no git call"、`:199-201` "any failure yields ok=false plus one reason"（含 `parseIndexPathsList` 对绝对路径/`..` 的拒读）没有测试喂过畸形 git 流 | 本批写的 | 前者被 §4 末节的 6-spawn 记账**旁证为真**；后者方向安全（拒读⇒不套规则⇒多扫） | **不升级为缺陷，登记** |
| **N5（旧账，点名归属）** | `gitignore.go:11` 引 `.gitignore:24` 指 `frontend/dist/*`，**现量该行在 `.gitignore:22`** | 不是本批写的（`ca84b75^` 逐字同串）⇒ 是 `3bb99aa` 留下的引用 | `sed -n '22,24p' .gitignore` 现量：22＝`frontend/dist/*` | 上一程的引用缺陷，本批重写了这段头的上下邻居却没顺手改；**登记，一句话改** |

---

## 11. 总裁

| 格 | 判 | 一理由 |
|---|---|---|
| 1 F1 真关了吗（复现） | **成立** | 仓库外重造 `git add -f` 两枚件：修后 `rc=1`、三枚 finding 点名（**含 `internal/build/leak.go` 那枚 AGENTS.md §1.2 的裸 goroutine**）、分母各 ＋1、note 点名"TRACKED and matching an ignore rule"；修前同树同参数 `rc=0`/40/203 |
| 2 是不是换方向的蒙眼／伪装退回 | **成立** | 未追踪且被忽略那一头：两版输出 `diff` **0 行**（40/203/406/39、`skipped as git-ignored` 行逐字同） |
| 3 "问不到"够响、朝多扫 | **成立** | 四支各测：无 `.git`（首行自陈＋44 vs 40＋rc 1）、**PATH 无 git 而 `.git` 在**（自陈点名 `exec: "git" … not found in %PATH%`、40→41、rc 0→1）、10 s 超时（`real 20.7s`、自陈 `timed out after 10s`）、别人仓库的子目录（自陈 `not its top`）；note 与账本同在 stdout |
| 4 那味偏离（全量 `ls-files`） | **附条件成立** | 摘掉它**今天不打红任何一发**（M1 全绿）＝按本仓定义它未被测试要求；但 M4a 单跑绿／M4a＋M1 探针红证明它挡的是**这一批刚犯过两次的活家族**（匹配器比 git 多跳）⇒ **防御不是装饰，条件＝补那第三枚半**；代价实测 6 枚只读 spawn／＋0.34 s 每次、只付在 step 2 |
| 5 F2 那两句现在可证伪吗 | **附条件成立** | 两句各有 pinner（M2 8 条红含 `:1534 got 0`；M5/M6 各咬一半）；**但本批新写了一句同族未钉 blanket（`:103-105`）＋一枚站不住的 union 理由（`:166-171`）** ⇒ 旧账结、新账两行文字 |
| 6 同批改副本这条规矩 | **成立** | `ignoredLikeGit` 签名同批加 `tracked`/`rulesOn`；在 F1 树上把副本那一半抽掉 ⇒ 正控 41/41 PASS、**副本退回纯模式即 `41 vs 40` 当场分家**；名册普查无第 7 份；`internal/panel/**` 零字节 |
| 7 契约轴 | **成立** | `emojiRe` 五锚同值且 diff 不碰它；allowlist 0 字节；断言 139→164 只增、`t.Skip` 7→7、无 `func Test` 消失；thresholds/golden 零字节；`main.go` −22 全注释、`gitignore.go` −33 逐枚点名（含被删的那两句担保） |
| 8 新环境依赖（缺 git 就红） | **成立（含更正）** | 有 git：28/0/0/68 rc=0；无 git：**25/3/0/68 rc=1**、零 Skip、红句点名"cannot be built without a git binary"；扫描器本体无 git 照跑朝多扫。⚠ 依据要换：CI 那枚 job 是 `ci.yml:66 ubuntu-latest`（不是 self-hosted Windows），但托管镜像必带 git 且 `actions/checkout` 就是 git ⇒ 结论不变 |
| 9 当下基线与归因 | **成立** | `A218⑥` 的八数我在 `5d463bb` 净快照逐枚复现（含首行自陈）；现锚 `2206cd4` 净快照＝`203/22/40/18 · 30/40/407/39`（`406→407` 归 `d88c356`、`16→30` 归 `5d463bb`）；**`ca84b75` 一枚读数都没动**：真工作树两版 `diff` 0 行、净快照仅多首行自陈 |
| 10 往前账 | **成立** | R2（旧 F3 两枚语义）结；R1/R3/R4/R5/R6 是登记过的兄弟账、本批未恶化；**新账只有 N1/N2 两枚文本面 ＋ N3 一枚缺断言**，都不改判据方向 |
| **本批整批** | **附条件成立（可入账，不需再来一程）** | 主攻击格（F1）在我自己重造的变异上真的关了，且没把另一头蒙上；两枚 F2 担保兑现；副本同批接上并实测有牙；契约轴零越界。**条件三件，全在一枚小 commit 里**：①补 §4 那发第三枚半（把全量 `ls-files` 那一味钉住，约 15 行测试、零生产码）；②`gitignore.go:103-105` 的"唯一形状"改成两支并列并归对治理者；③`gitignore.go:166-171` 的"独立并集成员"理由换成"note 点名"（顺带 `:11` 的 `.gitignore:24`→`:22`）。**这三件没有一件让门今天少看见东西**，所以不判退回。 |

---

## 12. 我**没有**测什么（绝不让沉默被读成批准）

1. **没跑 GitHub CI**（一行都没跑）。§8 关于 ubuntu 的两句是**推理**：镜像带 git ＋ `actions/checkout@v4` 依赖 git。
   我没读到任何一枚 run id / step 日志，也没核 `GITHUB_TOKEN` 侧的检出形状。
2. **没测 `core.excludesFile` 真被设过之后的读数**（只测了 `.git/info/exclude` 那一支，R4）。
   也不打算动任何 git 配置（本程禁改 git config）。
3. **没测并发/共享树**：本批代码在别人正在 commit 时读 index 的一致性我没造场景（`index()` 每枚 matcher 只问一次 ⇒ 中途
   有人 commit，同一趟扫描里前一半与后一半可能踩着两枚 index；方向应仍是"多扫或按老 index 少扫"，**未证**）。
4. **没测大小写／长路径／盘符大小写不一致／UNC** 下的 `filepath.Rel` ＋ `ToSlash` ＋ git 输出三者的对账（依赖实现，未构造用例）。
5. **没在 Linux 上跑过任何一发**：全部 win32 ＋ `git 2.52.0.windows.1` ＋ `go1.27.1 windows/amd64`。
   CI 那枚 job 恰好在 ubuntu ⇒ **`exec: "git": executable file not found in %PATH%` 这类文案在 linux 上是另一串**，
   而 `TestIgnoreRulesAreOffWhenTheIndexCannotBeAsked` 只锁了 `not a git repository` 那一支（§8 末的脆性发我只在本机论证）。
6. **没重跑上一程的第二枚"今天有没有这种件"口径**（逐件 `git check-ignore --no-index`）；我只用了 `git ls-files -i -c --exclude-standard`（现量空）。
7. **没测 `git ls-files -z` 输出畸形流**（绝对路径／`..` 那两条 `parseIndexPathsList` 守卫）——我没有办法让真 git 吐出那种流，
   也没造桩（stub）。⇒ N4 那句"any failure yields ok=false"只有代码事实，没有读数。
8. **没测 git 仓库巨大时的形状**（十万级受追踪件的 `ls-files -z` 内存/耗时），也没测 worktree（`.git` 是文件）那一支。
   我只在 1142-1144 枚受追踪件上量过 6-spawn／0.34 s。
9. **没测 step 1 在超时那一支下的表现**（10 s × 若干枚 fixture）——我只量了 step 2 的 20.7 s 最坏形状；
   门的 step 1 若有 fixture 撞上慢 git，总时长会涨，**未量**。
10. **没核 §4 那发探针在生产口径下会不会误伤**：我的探针是"匹配器读宽规则"这一类的**存在性**证明，
    今天交付版没有这种规则（M4a 单独跑才响），所以我**不能**说"这仓现在正漏着什么"，只能说"这一味今天没断言、将来靠它挡"。
11. **没读实现程 `d22scan-gitignore-fix-r1.md` 的每一枚自述去逐条对账**（我只把它当声明用，本件所有格都是现跑）；
    它 §11.3 的"删一行即退回"路径我没有替它验证删除后的**全套**门（我只跑了摘那一味的相关六枚测试＋我的探针）。
12. **没跑 root module 的任何测试**（`internal/**`、`cmd/**` 整包）：本批代码在独立 module、无 import 关系（上一程 §7 已量），
    但"因此不可能影响 root module"仍是**推理**不是我的读数；何况另有两程此刻正在那两个包里写码。
13. **没测 `--no-index` 之外的 git 版本差异**（不同 git 主版本对 `ls-files -i -c` 的 `--exclude-standard` 语义）。
14. **权限**：本程**没有任何一枚工具调用被权限系统拒绝**；也没 push、没 `git add -A`/`.`、没 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；
    每枚 commit 前 `git diff --cached --name-only` 只列 `docs/evidence/s1/d22scan-gitignore-fix-r1.md`（逐次核过，无第二枚混入），
    每次都是**显式 pathspec 提交**；对 owner 那 16 枚未提交的 `design/**` 删除**未做任何动作**（不 stage、不 commit、不还原）。
    ⚠ 一处我自己的手滑如实记：搭 `mut\` 时我跑过一次 `rm -f README*`，**该目录里没有 README ⇒ 实际零删除**；此后一律只建不删。
