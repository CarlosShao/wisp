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
