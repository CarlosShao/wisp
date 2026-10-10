# 301-v3 — 非实现者验收腿判语表（票 301 `AC#5` 七条形状 ＋ `AC#4` 残留必答 ＋ CI 新料判读）

腿＝`301-v3`（验收腿，非实现者；`AC#5` 的实现＝`301-r2`／`b761b584`，`AC#1` 的实现＝`301-r1`／`0a0f62ef`）。
本件只交**判语**：⛔ 翻任何 `AC` 框（翻勾权在编排者）、⛔ 写产品码、⛔ 改票面／台账／`HANDOVER`、⛔ push。
凡下面出现的数都带**五件事**：哪把尺＋射程目录＋含⛔ 含子测试＋blob 还是工作树＋跑出来的件名。
凡引编排者派单里的读数，一律标〔⛔ 我复跑〕或另给**我自己那把**尺的读数；本件的判语⛔ 有一枚以他的结论为凭据。

---

## 0. 双锚（起手／交回，同一把尺两发）

| 格 | 尺 | 起手读数 | 交回读数 |
|---|---|---|---|
| 时刻 | `date '+%F %T %z'` | `2026-10-10 16:07:23 +0800` | `2026-10-10 16:26:42 +0800` |
| HEAD | `git log -1 --format='%h %ad %s' --date=iso-strict` | `a0e55a51` ＋ `2026-10-10T15:51:05+08:00` | `6da0e688` ＋ `2026-10-10T16:22:20+08:00`（编排者 A818 笔） |
| porcelain 本票射程 | `git status --porcelain -- scripts internal .github docs frontend \| wc -l` | **0 行** | **0 行** |
| tasklist `wisp.exe` | `tasklist //FI "IMAGENAME eq wisp.exe"` | **0 枚**（`INFO: No tasks are running which match the specified criteria.`） | **0 枚**（同形逐字） |
| tasklist `balldebug.exe` | `tasklist //FI "IMAGENAME eq balldebug.exe"` | **0 枚**（同形） | **0 枚**（同形） |
| 件退码 | `echo "anchor_rc=$?"` | `anchor_rc=0` | `close_rc=0` |
| 判过的面在我这发窗口里动过没 | `git diff --name-only a0e55a51..HEAD -- scripts .github internal frontend design \| wc -l` | — | **0 行**（⇒ 我引的每枚 blob 在此窗口内⛔ 漂移） |

件＝`logs/anchor-start.txt`（起手，16:07:23）／`logs/anchor-close.txt`（交回，落笔前现量）。
★**起手锚跑完立刻落第 1 笔**＝`aea74556`（16:07:50，`git show --name-only --format=` 名册＝我自己那 2 枚件），
在任何长跑命令之前。⚠ **本腿在飞期间 `dev` 动过**：`a0e55a51` → `aea74556`（我）→ `6414a4bb`（编排者收 CI 色）→
`36efbd80`（`302-a1` 锚）→ `24ae75b7`（`298-r1` 锚）→ `48e7b0a4`。⇒ 我**所有承重读数都锚在具名 sha 上**
（`0a0f62ef`／`b761b584`／`b761b584~1`／`4be5168c`／`bcd0a543`／`cc315261`／`772ff880`），
只用工作树量的地方＝四枚门禁文件，且我逐枚证过它们与 `HEAD` blob 逐字节相同（V3-R7 那节）。

环境＝本机 `GOOS=windows / GOARCH=amd64`。射程目录＝仓根 `D:/work/workspace/projects plans/Wisp`。

---

## 1. 我这发的尺（命名 V3-R1..R7；与 `301-v1` 的 V-R* **⛔ 同一枚命令**，别混引）

- **V3-R1 逐笔名册尺**＝`git show --name-only --format= <sha> | grep -v '^$'`（⛔ 区间尺）。
- **V3-R2 注释面尺**＝`git diff -U0 <a> <b> -- .github/workflows/ci.yml`，再
  `grep -E '^[+-]' | grep -vE '^(---|\+\+\+)' | grep -vcE '^[+-][[:space:]]*#'`（＝非注释变更行数）。
- **V3-R3 可执行键尺**＝对 V3-R2 的新增行 `grep -cE '^[+][[:space:]]*(name|runs-on|if|run|steps):'`。
- **V3-R4 blob 文本尺**＝`git show <ref>:<path>` ＋ `sed -n` / `grep -n`，一律 **blob**、⛔ 工作树。
- **V3-R5 CI 名册尺**＝`gh run view <id> --json jobs --jq '.jobs[] \| [.name,.conclusion,.databaseId] \| @tsv'`；
  日志名册＝`gh run view --log --job=<id>` 先落**文件**，再 `grep -aE -- '--- FAIL: ' | sed -E 's/.*--- FAIL: ([^ ]+).*/\1/' | sort -u`
  （⚠ `gh run view --log` 的行**带 step 名前缀＋制表符**，⛔ 用 `^TIMESTAMP` 锚；编排者那两份件是另一条形，我两形都按形打）。
- **V3-R6 四数字尺**＝`grep -a 'four numbers'`（脚本自己那行，逐字标 `(all from -v output)`＝**含子测试**）。
- **V3-R7 门禁尺**＝`sh scripts/d22scan.sh`／`bash -n scripts/portable-tests.sh`／`bash scripts/portable-tests.sh --scope=census`，
  每把自己落一行 `rc=N`（`logs/*.rc.txt`，全非 0 字节）。

---

## 2. ① `AC#5` 逐条判语（七条形状 ①..⑦ ＋ 整格）

靶的锚（内容锚，⛔ 行号）＝被删那两行的**原文**；我这发在 `b761b584~1` 的 blob 上逐字读到（V3-R4）：

```
  # Windows-only packages (ball GUI, cgo speech) are out of this job's scope
  # by platform, not skipped: they run in test-windows / slo jobs.
```

⇒ 与票面 `AC#5` 的引文**逐字相同**；它们在 `test-core:` 那块的头上（V3-R4＝HEAD 上 `test-core:` 那一行现在在 `:419`，
被新增的 17 行推下去；票面"现量"节写的 `ci.yml:400-401` 是**改动前**的位置，`@@ -400,2` 那一枚 hunk 头就是它俩，⛔ 漂）。

### ①「单独一笔」＝**成立**

- 尺＝V3-R1＋`git log -1 --format='%h parents=[%p] %s' b761b584`。读数＝`b761b584 parents=[4be5168c]`（单亲、⛔ merge）；
  V3-R1 名册 7 枚＝`.github/workflows/ci.yml` ×1 ＋ `.scratch/wisp/probes/301/r2/**` ×6，
  **`scripts/portable-tests.sh` 命中＝0** ⇒ 与 `AC#1` 那笔（`0a0f62ef`，名册 28 枚、`ci.yml` 命中＝0）确实是两笔。
- 尺＝`git log --oneline 4be5168c..HEAD -- .github/workflows/ci.yml` ⇒ 命中**只有** `b761b584` 一枚（这一列里⛔ 第二笔碰 yml）。
- 尺＝`git show --name-only --format= 4be5168c`＝`r2/00-anchor.md`＋`r2/msgs/00-anchor.txt` 两枚（锚件⛔ 带 yml）⇒ 拆形＝锚笔＋注释笔。
- 次序硬约束（票面：`AC#5` 必须按在 `AC#1` **落地并核过**之后）我这发自己打：
  `git merge-base --is-ancestor 0a0f62ef b761b584`＝yes、`--is-ancestor bcd0a543 b761b584`＝yes
  （`bcd0a543`＝`301-v1` 的裁决笔）⇒ **落地笔在 `AC#1` 与其验收件之后**，约束满足。

### ②「pathspec＝`ci.yml` ＋ 落地腿自家探针目录，且写在 `$( … )` 之外」＝**成立（可观测面）／半枚 ⛔ 判不了**

- 可观测面（V3-R1，blob 取集）＝`b761b584` 名册 7 枚，剥掉 authorized 面
  （`^(\.scratch/wisp/probes/301/r2/|\.github/workflows/ci\.yml$)`）后**外人件＝0 枚**；`4be5168c` 同尺＝**0 枚**。
  ⇒ 这一发**没有吞任何别人的 staged 件**，那枚子句要买的后果测得＝0。探针面＝`r2/`（腿自家目录），⛔ 票级 glob `probes/301/**`
  ⇒ `301-v1` §3-⑦ 那条具名改良**这一腿照做了**（我实测名册里 `probes/301/orch/`、`probes/301/r1/` 各 0 枚）。
- ⛔ 判不了的那半枚＝命令行形本身（`git commit` 带没带 pathspec、`$( … )` 内外）。盘上无残留：这一笔**取集＝authorized 面**，
  所以"缺省收整索引"与"显式 pathspec"两形产**同一份名册**，判别尺打不出来（我试过的那把＝`git show --format=%B`＋名册对拉，
  只能证 message 没被 `cat` 吞：`b761b584` 的 subject 是一行正常说明文，⛔ 是文件正文 ⇒ `42617eba` 那枚坑这一腿没踩）。
- ⇒ 判语＝**形状②成立**（后果面 0、面⛔ 越界、message 形正常），命令行那一半记〔⛔ 盘上判不了〕并归进 §3(a) 的同一枚射程。

### ③「纯注释面、零行为」＝**成立（三把尺 ＋ CI 面复跑）**

- V3-R2 尺：`git diff -U0 b761b584~1 b761b584 -- .github/workflows/ci.yml` ＝ **1 枚 hunk `@@ -400,2 +400,19 @@ jobs:`**；
  变更行 **21** 枚（删 2 ＋ 增 19）；其中**非注释行＝0**（件 `logs/ac5-diff-u0.diff`，1,810 B）。
- V3-R3 尺：新增行里 `^\s*(name|runs-on|if|run|steps):` 命中＝**0**。
- 我自己的第三把（⛔ 用 python，比"数 hunk 行"更全局）＝
  `diff <(git show b761b584~1:…ci.yml | grep -vE '^[[:space:]]*#') <(git show b761b584:…ci.yml | grep -vE '^[[:space:]]*#')`
  ＝**0 行**（件 `logs/ac5-noncomment-normalized.diff`，0 B＝空即证据）⇒ 剥掉注释行后的整文件逐字节相同。
- 编排者那把 `yaml.safe_load` 正规化尺我也复跑了（`python -c`，同一枚判语）＝`safe_load_equal= True`、
  `jobs_before= 6`、`jobs_after= 6`。⇒ 他那四把里这一把我打中了（同一判语、两条形）。
- **CI 面（新料，我这发自己下的日志）**＝第一枚**带着这 17 行注释**的 run＝`38035842314`（`headSha a0e55a51`，
  event push，conclusion failure）。它相对改后基线 `38034689386`（`bcd0a543`，⛔ 带注释）的颜色：
  V3-R5 job 层六枚**逐枚同色**、V3-R6 四数字**逐字相同**（`101/58/0/0`·`101/58/0/0`·`398/278/8/2`·`624/427/12/1`）、
  `census totals` 行逐字相同、windows 档红名集合 `diff`＝**0 行** ⇒ **"零行为"这半句在托管 runner 上也被证了一次**，⛔ 只在文本层。
- ⛔ 有人事后动过这枚文案：`git diff --stat b761b584 HEAD -- .github/workflows/ci.yml`＝空；工作树 `ci.yml` 与 `HEAD` blob `cmp` 相同。

### ④「文案必须把不对称写进去」＝**成立（而且逐句我都能验真）**

新增 19 行的原文见 `logs/ac5-diff-u0.diff`（我从 blob 逐字读的）。不对称是**分开两句**写的，我逐句打尺：

1. 「`./internal/audio/` is named in this tier's own array (`core)`)」→ V3-R4 on `HEAD:scripts/portable-tests.sh`
   ＝`:242` 那一行 `./internal/buildinfo/... ./internal/audio/... ./internal/proc/...`（`core)` 数组）✓，
   另有 `:168`(`core_pin`)/`:195`(`win_pin`)/`:253`(`windows)` 数组)/`:594`(ledger) 四处命中，枚枚具名。
2. 「a linux target still resolves sources plus untagged test files there」＋它自带的尺
   `CGO_ENABLED=0 GOOS=linux go list -f '{{len .GoFiles}}/{{len .TestGoFiles}}' ./internal/audio/`
   → 我照它跑＝**`9/5`**（`TestGoFiles`＝5 ≠ 0），同尺 `GOOS=windows`＝**`11/8`** ⇒ 那句"包级 out of scope"确实是旧的谎，
   新句"粒度只到文件"是真的。**这一枚⛔ 是引编排者：我自己跑的。**
3. 「test-windows has taken `./internal/audio/` since ticket 301 `AC#1` landed」→ V3-R4：`win_pin` 与 `windows)` 数组各含 audio（上面那两处），
   且 CI 字节直接证：census 里 audio 那行**基线 `5/0 core` → 改后 `8/0 corewindows`**（V3-R5 打在两份已落盘日志上）。
4. 「while on that same date the slo tiers reached no package list …`slo-check.ps1` builds `wisp.exe` and samples D32 states instead of running package tests」
   → 我这发现量：`ci.yml` 里 `slo-smoke:`/`slo-full:` 两块**只有** checkout/setup-go/cache/`Build wisp.exe`（`scripts/build.ps1 -Env dev`）/
   `SLO … gate`（`scripts/slo-check.ps1 -Subset smoke|full`）/`upload-artifact`；`slo-check.ps1` blob 里 `portable-tests` 命中＝**0**、
   `go test` 命中＝**1** 且那一处（`:336`）在注释里、`:294-316` 确实在 `wisp.exe` 缺失时拉 `scripts/build.ps1 -Env dev` 自建
   ⇒ "builds"⛔ 假、"reached no package list"⛔ 假。
5. 「`(memory/handle subset, no audio)`」→ 逐字命中 `slo-check.ps1:15`。**这一枚在 `301-v1` 是〔⛔ 我复跑〕的引文，我这发升成实测。**
6. 文案自己还带了一句防再造谎的：`this comment is not the evidence` ＋ `re-read them` ⇒ 它⛔ 长成第二座真相源。

⇒ 判语＝④**成立**：一句变真、一句仍假，两半各带尺，且⛔ 有一句我打不中。

### ⑤「⛔ 写会被自己扫的现在时计数」＝**成立（没写）**

- 尺＝V3-R2 的新增行取数字 token：`grep -E '^\+' | grep -oE '[0-9]+[0-9a-zA-Z.#/-]*' | sort | uniq -c`
  ⇒ 出现的数字**只在 5 行里**，token 全集＝`2026-10-10`（日期，×2）、`301`（票号，×2）、`32`（＝`D32` 决策标号）、`1`（＝`AC#1` 标号，×3）、`0`（＝尺串里的 `CGO_ENABLED=0`）。
- ⇒ **零枚**"11 枚 pin／42 枚用例／0 命中／只有 3 枚"这类现在时全称计数；枚数一律**留给尺**（文中两处尺串逐字可跑）。
- 经验面：`sh scripts/d22scan.sh` rc=0（V3-R7，我这发）⇒ 这 17 行没有一枚被自家仪器扫到；
  `--scope=census` totals 行逐字未变 ⇒ 文案⛔ 与任何一枚在跑的计数打架。

### ⑥「门禁」＝**成立（三把我都自己跑）**

| 发 | 命令（V3-R7） | rc | 件（stdout/stderr） | 承重读数 |
|---|---|---|---|---|
| 1 | `sh scripts/d22scan.sh` | **0** | `logs/d22scan.txt` 23,086 B／`logs/d22scan.err` 0 B／`logs/d22scan.rc.txt` | 逐字 `d22scan_rc=0` |
| 2 | `bash -n scripts/portable-tests.sh` | **0** | `logs/bashn.txt` 0 B／`logs/bashn.err` 0 B／`logs/bashn.rc.txt` | `bashn_rc=0` |
| 3 | `bash scripts/portable-tests.sh --scope=census` | **0** | `logs/census.txt`／`logs/census.err` 0 B／`logs/census.rc.txt` | totals＝`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`（逐字**未变**）；audio 行＝`8/0 corewindows` |
| 4 | 同尺 CI 面（两份 run 的 test-windows 日志） | — | `logs/fournums-*.txt`、`logs/a0e55a51-test-windows.log` 1,419,603 B | `census totals` 两枚 run 逐字相同 |

⚠ 我这发**⛔ 跑** `--scope=windows`／`--scope=core` 的本机颜色（理由具名见 §5）；那两档的颜色我用 **CI 字节**（V3-R5/R6）替代，
且那些字节证的恰是这两档在**托管 runner** 上的形——比本机一发更贴 `AC#2` 那句"CI 现在会跑它"。

### ⑦「⛔ 零 push」＝**成立（就落地腿而言）＋一枚具名状态更正**

- 尺＝`git branch -r --contains b761b584`＝**2** 枚远端分支含它、`git rev-list --count origin/dev..HEAD`＝0、`git log -1 --format='%h' origin/dev`＝`48e7b0a4`。
  ⇒ `b761b584` **已经被推了**——这一句必须写清归属，否则下一个读者会把它记成落地腿的违反。
- 归属尺（我这发想的，⛔ 引编排者自述）＝`gh run view 38035842314 --json headSha,createdAt`＝`a0e55a51` / `2026-10-10T07:51:05Z`，
  与 `a0e55a51`（编排者自己的落账笔"收 301-r2"）的 commit 时刻 `2026-10-10T15:51:05+08:00` **同秒**；
  而 `b761b584`（15:47:05）之后⛔ 再产生过任何 run ⇒ **推出的是编排者那一枚 push（把 15:47 与 15:51 两笔一起带走）**，
  ⛔ 落地腿自己推的。⇒ 判语＝**⑦ 成立**（腿名面零 push 证据），⛔ 别把"已推送"记成这一格的违反。
- 我这发：⛔ push（三枚远端读尺全是只读 API）。

### `AC#5` 整格

七条形状 **全部成立**（其中②的命令行那一半记〔⛔ 盘上判不了〕、⛔ 影响判语），
靶文两句的谎都在我这发的尺下**改成了真的**，且新增文案的每一枚承重子句我都自己打中过（§2-④ 的 1-5）。
⇒ **本格可以翻**（翻勾权在编排者）。⚠ 一句给下一个读的人：本格买到的是**文案的诚实性**，
它⛔ 改变 CI 行为（V3-R2/R3 ＋ §2-③ 的 CI 面都＝0），⛔ 任何人别把它读成"coverage 有进展"——coverage 那一笔是 `AC#1`。

---

## 3. ② 必答：`AC#4` 那枚残留（"commit 必带显式 pathspec"）

票面 `AC#4` 末句逐字（V3-R4 on 票面）＝「每把门禁件自落一行 `rc=N`（⛔ 0 字节＝那格没交）；⛔ 零 push、commit 必带显式 pathspec。」

**(a) 那枚子句的字面射程管不管"已经发生的一笔"＝**部分管**：管**判**、⛔ 管**满足**。**

- 尺（射程）＝那枚子句的**出处原文**（`AGENTS.md` §1.4，逐字「`commit` 必须带**显式 pathspec**；**禁 `git add -A` / `git add .`**」）
  ＋票面 `AC#4` 那一格的主语。两处的语法主语都是**一次动作的命令形**（"commit 必带"），⛔ 是"树/历史的状态"。
  ⇒ 已经发生的那一笔（`0a0f62ef`）**落在它的射程内**，因为它要评的正是那一枚 act；违反是一枚**永久事实**。
- 尺（事实）＝V3-R1 on `0a0f62ef`：名册 **28** 枚，剥掉写腿自己的 authorized 面（`probes/301/r1/**`＋`scripts/portable-tests.sh`）
  ＝**27** 枚，**外人件 1 枚**，具名＝`.scratch/wisp/probes/301/orch/2026-10-10-baseline-and-delta-prediction.md`（编排者的件）。
  我这发还量到：该件在 `0a0f62ef~1`（`74eb032c`）的名册里**⛔ 存在**、在 `HEAD` 树里 `git cat-file -e` **存在**
  ⇒ 它**第一次进仓就是进在写腿那一笔下**＝署名与内容不同主，那枚子句要防的事确实发生了（⛔ 只是"名册宽"）。
- ⛔ 管"满足"的那半＝没有任何后发动作能让"那一枚 commit 带了显式 pathspec"变真；
  而 `--amend`/`reset`/改写历史是本仓硬禁（`AGENTS.md` §1.4 逐字「已推送的历史不改写，要更正就**追加新 commit**」）。
  ⇒ "追加闭"**⛔ 满足原句**；它满足的是**另一件事**（把违规记到能防下一次）。
- ⚠ 机制那一半我这发仍判不了（同 `301-v1`，且我另找了一把新判别尺，落空）：
  `42617eba` 的名册**⛔ 含外人件** ⇒ 它证⛔ 了"缺省收整索引"那一形；两形产同一份 28 枚取集。
  ⇒ 记〔机制⛔ 定，**归腿**这件事⛔ 依赖机制〕。

**(b) 若管，`AC#4` 整格能不能翻＝**⛔ 能**（至少⛔ 能翻成"全格成立"；我这发⛔ 授权翻）。**

- 理由具名：那一格是一串**合取**的子句，其中一枚被**射程内、不可追溯**地违反了；
  而且它的**意图**（一笔 commit 只装一枚腿自己的活）被**可测地**打碎过一次（§3(a) 的那 1 枚外人件）。
  "内容无损"我独立复跑过（`git show 0a0f62ef:<orch 件> | cmp - <工作树>` 同 `301-v1` §3-①(c) 的形；我这发用的是
  `cat-file -e` ＋ 该件此后**没被任何别的笔动过**＝`git log -1 -- <path>` 回 `0a0f62ef`）⇒ 丢的是**署名**，⛔ 是字节。
- ⛔ 同意把这格读成"永远⛔ 能结"——那⛔ 是本仓要的口径；正确处置是**改判据要具名**，⛔ 由验收腿顺手盖章。
  所以我把"能翻"的前提写死：**要么保持未勾＋台账具名**，**要么由编排者追加一格改写判据**（见 (c) 第 2 件）。
  追加 AC 格这一手在本票**已有先例**＝`AC#5` 自己就是编排者 13:3x 追加的（票面逐字「`AC#5`（2026-10-10 13:3x 编排者追加…）」），
  且 AC 格⛔ 属 `C1–C32`/`D1–D47` 那两枚契约面 ⇒ 落在他名下，⛔ 需要新的批准层级。
- 本格其余子句我这发**逐枚复跑＝全部成立**（写给编排者落账用，⛔ 引我的结论当凭据的人请重跑）：
  `d22scan` rc=0／`census` rc=0（§2-⑥）；`windows`·`core` 改前改后逐名作差**具名新增红 0 枚**——我这发在 **CI 字节**上重证一次：
  `ci-fail-after-windows.txt` ∩ audio 42 枚名册＝**0**，且 `bcd0a543`↔`a0e55a51` 两枚 run 的 windows 红名册 `diff`＝**0 行**；
  `git show --stat` 名册尺（V3-R1，六笔＋r2 两笔）相对 authorized 面 `scripts/portable-tests.sh`＋`probes/301/**` **越界 0 枚**；
  `frontend/**`／`design/**`／三枚冻结件／golden／`thresholds.go`／`allowlist.txt`／D43 表＝**逐笔 0 枚**
  （r2 两笔名册里只出现 `ci.yml`＋`r2/**`，我逐枚比过）；`:594` 那枚 ledger 行**逐字节未动**
  （尺＝`diff <(git show 0a0f62ef~1:… | sed -n '593p') <(git show HEAD:… | sed -n '594p')`＝空）；
  ⓿字节证据件＝**⛔ 一枚 rc 件都是 0 字节**（r1 的 9 枚、r2 的 3 枚、v1 的 12 枚；我这发＝6 枚 `.rc.txt` ＋ 3 枚件内嵌 `rc=` 行，全非 0），
  其余 0 字节件全是结构性空（我这发 **9** 枚逐枚给原因：`bashn.txt`/`bashn.err`＝`bash -n` 本来⛔ 产 stdout/stderr；
  `d22scan.err`/`census.err`/`golist.err`/`ci.err`/`core.err`/`a0e55a51-dl.err`＝全 rc=0 的发次stderr／下载⛔ 报错，应空；
  `ac5-noncomment-normalized.diff`＝"空即证据"那一形）。⛔ 零 push＝写腿与两枚验收腿名面都无证据（§2-⑦ 那枚归属更正除外）。

**(c) 残的到底是哪一件东西、该由谁、用什么形状收口**

- **残体＝一件历史事实，⛔ 一件产物**：`0a0f62ef` 同时背着「命令形违规」＋「外人件署在写腿名下」，两者都不可追溯消除、
  也⛔ 能靠"再交一枚件"抵掉。⇒ 任何"追加闭＝那格绿了"的说法都⛔ 成立；能追加的只有**记录**与**前瞻闸门**。
- 收口三件（各带归属，⛔ 一件归我执行）：
  1. **台账追加一行**（`docs/reports/pending-and-issues.md`，`A##` 只追加⛔ 删既有行；归**编排者**）：
     内容五件＝笔号 `0a0f62ef`、子句"显式 pathspec"、机制两说并存⛔ 定、归 `301-r1`、
     "追加闭⛔ 满足原句"判定＝本件 §3(a)-(b)（凭据尺＝V3-R1 取集 28 ⊃ authorized 27 ＋ 外人件名）。
  2. **票面追加一格**（建议名 `AC#4b`，归**编排者裁**；⛔ 改 `AC#4` 原句一字）：把本格的结案形写死成
     "其余子句全交＋违规具名入台账＋前瞻闸门已立"，并逐字带上"pathspec 那一枚子句本身判为**违反且不可满足**"。
     ⛔ 这一格，`AC#4` 那个框无论勾不勾都会让下一个人读不出它到底结在哪。
  3. **前瞻闸门**（⛔ 我实现、⛔ 我动文本；只交形状给编排者/人工裁）：
     (甲) 派单/工单池规矩文本里加一条"commit 前必跑 `git diff --cached --name-only`，输出须 ⊆ 自己的面"；
     (乙) 根因那一半⛔ 在腿手上＝**编排者的落账件⛔ 在别人可能收整索引的窗口里 stage**；
     形状＝编排者在自己的窗口里 add＋commit（同 `301-v1` §3-①(b) 那句），或把补位件落到⛔ 被任何票级 glob 命中的面
     （`probes/orch/<ticket>/`，`A812` 已把它降级为**条件性**纵深防御——我这发同意那个降级，理由＝
     `0a0f62ef` 那一枚外人件恰好**躺在** `probes/301/**` 里，所以换目录在那一形下挡得住，而在"缺省收整索引"那一形下⛔ 挡得住）。

---

## 4. ③ 新料判读：CI 色回填会不会改动 `AC#2`／`AC#4`／`AC#5` 任何一格

**结论先写：⛔ 任何一格被打红；`AC#2` 的凭据被加强（v1 具名那枚欠账可销）；`AC#5` 的"零行为"多了一层 runner 面证据；
`AC#4` 那枚残留⛔ 因此变好也⛔ 因此变坏。** 下面枚枚是我自己复跑的读数（⛔ 引派单那句当凭据）。

### 我复跑了什么（派单要求至少一把，我打了六把）

1. `gh run view` 两枚 run 的 `headSha/conclusion/event/createdAt`＝`bcd0a543…`/failure/push/`07:31:01Z`；`cc315261…`/failure/schedule/`23:26:26Z`
   ⇒ 派单给的"同一台托管 runner、改后＝`bcd0a543`、基线＝`cc315261` 的 schedule 发"**逐字对上**。
2. **六枚 job 色逐枚同色**（V3-R5，两枚 run 各 6 枚 job，件 `logs/ci-jobs-{baseline,after}.tsv`）＝
   `lint` failure／`test-core` failure／`test-windows` failure／`lint-frontend` success／`slo-smoke` success／`slo-full` success。
3. **四数字逐字复跑**（V3-R6 打在两份已落盘日志上）＝windows `577/386/12/1` → `624/427/12/1`；
   cli `335/231/6/1` → `398/278/8/2`；winsec `101/58/0/0` 两枚 run 相同。
   ⇒ **A-evaluated 我自己算**＝386+12+1＝**399** → 427+12+1＝**440**（**＋41**）。
4. **红名册按包分档**（我自己的分档尺：全 run 红名 − `cmd/wisp` 304 枚顶层名 − `internal/winsec` 84 枚）
   ＝ windows 档红 **12 枚 ↔ 12 枚**、`diff`＝**0 行**；`∩` audio 名册（42 枚，V3-R4 on `HEAD` 的 8 枚 `_test.go`）**＝0 枚** ✓。
5. **audio 侧逐名求值尺**（42 枚名一轮询"有没有顶层结果行"）＝改后 run 里**唯一没被求值的＝`TestLiveWasapiSmoke`**（1/42），
   它在日志里被点名 9 次（`-skip` 模式串回声）而 `--- SKIP` 枚数两枚 run 都＝1 且⛔ 是它 ⇒ **静默过滤器**这一条我在 **CI 字节**上复认。
6. **CI 侧 census**（`--scope=census` 的调用点其实**长在 test-windows 那枚 job 里**，见 §6 报回第 3 条）＝
   `census totals` 两枚 run **逐字相同**；audio 那行 `5/0 core` → `8/0 corewindows`。

### 三格逐格判

- **`AC#2`＝⛔ 改动判定（判语维持"够格"），且 v1 具名欠的那半句可以销。**
  票面那句「⛔ 拿本机跑绿当凭据」今天被 CI 色**补上⛔ 被它推翻**：本格判据的**定义**就是"同一枚用例两种档各一发＋名册逐名作差"，
  本机档差已经闭环；CI 字节买的是**另一枚东西**（"CI 现在真的会跑它"），而那枚东西现在有直接 witness：
  **票 300 那枚判据 `TestParseWaveFormatSubFormatOffset300` 在基线 run 的结果行＝0 枚、在改后 run ＝1 枚顶层 `--- PASS` ＋ 4 枚子测试 `--- PASS`**
  ⇒ 本票立票理由（"那枚判据交进来以后在 CI 上从未被执行过一次"）**第一次有 CI 字节把它反过来了**。
  v1 §3-② 还留了两枚"红了先看它们"的具名计时风险，我这发在 CI 字节上都读到色：
  `TestPinnedThreadStable10s`＝`--- PASS (10.06s)`（基线里⛔ 存在＝它确实是新求值的，且那 10 秒墙钟真付给了托管 runner）；
  `capturelevel_windows_test.go` 的唯一顶层用例 `TestAC247RealCaptureLoopEmitsLevels`＝`--- PASS (0.28s)`。⇒ 两枚风险**都⛔ 红**，欠账可销。
- **`AC#4`＝⛔ 被这批数字改动（⛔ 加严也⛔ 放松）。**
  最要紧的是**别**让那 3 枚 cli 红挂到本票头上（下面另给尺），也**别**把"windows 档红数 12↔12、`∩audio`＝0"读成新的凭据：
  它证的正是 `AC#4` 那句"具名新增红 0 枚"，而那一枚子句 v1 与我都用本机名册证过 ⇒ 这批数字是**第二次**打同一枚判语（同向），
  ⛔ 它救 `AC#4` 的那枚 pathspec 残留（残留是命令形/署名问题，⛔ 是颜色问题）。
- **`AC#5`＝⛔ 被改动判定，但判语多了一层。** 注释笔被 CI 求值了吗？求了：`38035842314`（`a0e55a51`，含 `b761b584`）
  的 job 色/四数字/红名册/census totals 与 `38034689386` 逐枚同 ⇒ **文本层的"零行为"在 runner 面也是零行为**。⛔ 翻勾理由，⛔ 反对理由。

### 有没有哪格被这批数字**反过来打红**＝**没有**；但有两枚"⛔ 是本票的账、是台账的账"的新料

1. **Δ 的一枚未解释数我这发解释了**（⛔ 需要谁裁，只是补账）：`=== RUN` 577→624＝**＋47**，而被求值的顶层＝**＋41**。
   ⇒ 那 **＋6** ＝子测试粒度（四数字那行逐字"(all from -v output)"含子测试）：
   尺＝改后日志里**缩进的结果行**按父名 `∩` audio 名册＝`TestOpenOccupiedAndPermissionDenied` 2 枚 ＋ `TestParseWaveFormatSubFormatOffset300` 4 枚 ＝ **6**，
   与 `+47 − +41` 逐枚相抵。⇒ 派单那句"与我在腿交件前写死的 Δ 预测同数"（＋41）**⛔ 因此被动摇，反而被补齐**。
2. **census 首列 `5/0 → 8/0` ⛔ 是本票造成的**：尺＝`git ls-tree --name-only <ref> internal/audio/ | grep -c '_test\.go$'`
   ＝`cc315261` **5** 枚／`772ff880` **7** 枚／`HEAD` **8** 枚。⇒ 两枚 run 之间 audio 的测试文件多了 3 枚（其中 1 枚＝票 300 那枚判据），
   这是**代码世代差**⛔ 是档差。写在这里是为了别让下一个人把 `5/0→8/0` 记成"本票让包长大了"。
3. **3 枚 cli 红与本票的因果，我这发用**时序尺**切开（⛔ 靠名册、⛔ 靠"大概是负载"）：**
   test-windows 那枚 job 内的档次序＝`winsec` → `cli` → `windows`，逐档四数字行的时间戳为证
   （改后 run：cli `07:41:34.576Z`、windows `07:42:24.861Z`；基线 run 同形 `23:36:11`/`23:36:56`）
   ⇒ `AC#1` 新加进这一档的 41 枚 audio 用例**严格跑在那 3 枚面板判据之后**，⛔ 可能是它们红的原因。
   补强尺：三枚名在 `cc315261` **都已存在**（`git grep -l <name> cc315261 -- cmd/wisp` 各 1 枚），
   基线色＝`--- SKIP (0.00s)`×2（AC13）/`--- PASS (0.97s)`＋`(1.25s)`（两枚 AC14），改后＝三枚 `--- FAIL (≈20.04s)`
   配"15 秒内没从页面收到任何回执（door 上到的＝什么都没有）"那句红句 ⇒ **形＝票 33 那族真窗判据超时**，
   加上同族第 4 枚反向漂（`TestPanelHostLatencyPercentilesAC2`：基线 `--- FAIL (0.00s)` → 改后 2× `--- SKIP`，
   红句⛔ 是红而是一枚**具名 skip**："no cold/hot sample recorded in this process: the lifecycle test did not run in this binary"）
   ⇒ 归**票 33 那族／托管 runner 负载与前置件未同批**，交编排者入台账；⛔ 本票射程、⛔ 本票能闭。
   ⚠ 一枚给下一个人的读法警告：这**两枚 run 之间是 741 枚 commit**（编排者 `A813` 自己写的推送枚数），
   所以任何 Δ 都是"代码世代差"⛔ 是"本票差"；本票唯一能在 CI 字节上归因的事＝**audio 那 41 枚进/没进结果行**（第 5 把尺），
   那一枚我归因得动、其余归因我⛔ 试。

---

## 5. ④ 我这发的自报缺口（⛔ 给自己盖章"已闭"）

- **跳过的长跑**＝`--scope=windows` 与 `--scope=core` 的本机实跑（`301-v1` 跑过、我这发⛔ 再跑）。
  理由具名＝编队此刻**在飞**（我这发亲眼见 `36efbd80` `302-a1` 起手锚、`24ae75b7` `298-r1` 锚进仓），
  而 windows 档里躺着 `TestPinnedThreadStable10s` 那种判 `LockOSThread` 独占性的 10 秒用例 ⇒ 我一跑就污染别人的计时读数。
  替代凭据＝CI 字节（§4 那六把）；代价＝我**没有**本机↔CI 的同机对照，若有人问"托管 runner 与这台机有没有分歧"，我这发答⛔ 了。
- **没复跑的转述**＝① 红句层"12 枚里 5 枚点名 `C:\Users\RUNNER~1`、2 枚是同文件的 A/B 表判定"——我只做到**名层＋包层**
  （12 枚名＝`TestAListWinsWhereBothTablesHit`/`TestBListDefaultDenyAndOverride`/`TestCanonicalInputGainsNoSecondForm`/
  `TestClassifyAnchorSpellingIsNotVerdict`/`TestPathResolver{ExtendedLengthPrefix,ShortName,UNC}AListDenied`
  ＋ 5 枚 `TestSync*`（`SyncFallbackNotDisarmableByWeakRoot`/`SyncFixtureFallbackAndMatch`/`SyncSuspectFallbackIsComponentBounded`/
  `SyncSuspectFallbackWhenUndetectable`/`SyncUnverifiedRootKeepsFallback`）＝7＋5；
  那 7 枚各自红句里有没有那串 8.3 别名，我⛔ 逐枚读（`RUNNER~1` 在改后日志里命中 **736 行**，枚数⛔ 能从它反推）。
  ⇒ 〔采编排者转述＋标〔⛔ 我复跑〕〕。② `lint` 那两枚 gofumpt 红与 `staticcheck` 欠账的归因我⛔ 复跑（⛔ 本票射程）。
- **我动过的工作树面＝0**：起手与交回 scoped porcelain 都＝0 行；我⛔ 摘过 pin、⛔ 装过 before-blob（`301-v1` 那两个窗口我这发⛔ 重做，
  它的 before 复现＋反形已经覆盖了那半枚判据，而我这发的新料是 CI 字节）。四枚门禁文件逐枚 `cmp` 与 `HEAD` blob 相同（§0 末）。
  产物＝只新建 `.md`/`.txt`/日志件，全在 `probes/301/v3/**`；⛔ 新建 `.sh`/`.ps1`/`.go`；⛔ 删任何件。
- **判断力够不到的半格**＝① `0a0f62ef`/`b761b584` 两笔的 **`git add`/`git commit` 命令行形**（盘上无残留，见 §2-②、§3(a)）；
  ② "`AC#4` 该⛔ 该以追加判据翻勾"——那是编排者的裁量，我只交形状（§3(c)）；
  ③ 票 33 那族 3 枚超时红的真归因（要真窗＋控制组，本票与派单都⛔ 许我碰那枚资源）⇒ 那半格⛔ 算我的欠账，算台账的一行。
- **一枚过程欠账，同一把尺打我自己**＝tasklist 我起手量了并**当场落件**（这点比 `301-v1` §5 干净），
  但在我**唯一那发长跑**（census，16:1x）之前我⛔ 再量一次；⇒ 若有人要"每发长跑前各量一次并并排落同一枚 rc 件"那一形，我这发**⛔ 满足**。
- ★**一枚我自己抓到的写面越界（越界体在我这发，⛔ 在写腿）**：为了证"我跑门禁时那四枚文件与 `HEAD` blob 逐字节相同"，
  我用 `git show HEAD:<path> > tmp-<basename>` 造了 **6 枚逐字节拷贝**（`tmp-ci.yml`/`tmp-d22scan.sh`/`tmp-portable-tests.sh`/`tmp-slo-check.ps1`
  ＋ `tmp2-ci.yml`/`tmp2-portable-tests.sh`）——派单写面那条逐字是"只新建 `.md`／`.txt`／日志件，⛔ 新建 `.sh`／`.ps1`／`.go`"，
  我新建了 `.sh`/`.ps1`/`.yml` 后缀的件＝**违反字面**（意图面⛔ 违反：它们⛔ 是产码、是 `cmp` 的底片）。
  处置形＝**⛔ 删**（`issues/README` 规则 8 临时件只建不删）⇒ 我这发把它们**逐字节改名**成 `*.removed-ext-cmp.txt`
  （字节⛔ 变、内容⛔ 丢，改完再 `cmp` 回 `HEAD` blob 三枚全 rc=0），现量尺＝
  `find .scratch/wisp/probes/301/v3 -type f \( -name '*.sh' -o -name '*.ps1' -o -name '*.go' -o -name '*.yml' \) \| wc -l` ＝ **0 枚**。
  影响面我也量了：改名后重跑 `sh scripts/d22scan.sh`＝**rc 0**（`logs/d22scan-final.rc.txt`），
  与起手那发的输出差 **138 行**、差行全是 `-v` 结果行的**耗时字段回声**（`0.06s`→`0.05s` 那一形），
  `grep -ciE 'violation\|banned\|forbid'` 在差行里＝**0** ⇒ 没有一枚分母被我挪动。**定式建议（归编排者裁，⛔ 我动文本）**＝
  "`cmp` 的底片一律落成 `.txt` 名"，别按被比文件的后缀起名。
- ⛔ push、⛔ `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`/`--no-verify`、⛔ `git add -A`/`.`、⛔ 在仓内建 worktree、
  ⛔ 改 git 配置、⛔ 动 SLO 阈值/golden/`internal/observe/thresholds.go`/`tools/d22scan/allowlist.txt`/D43 表/C1–C32/D1–D47、
  ⛔ 为变绿放宽任何断言、⛔ 顺手修任何产码（我看到⛔ 舒服的地方只有派单措辞，全部写进 §6⛔ 动盘上件）。

---

## 6. 具名报回派单（以盘上原文为准；⛔ 一处我改过盘上件）

1. **`cmd/wisp` 那档⛔ 是"一个 `--scope=cli` 的调用点"**。现量＝`ci.yml` 里只有 `--scope=core`／`--scope=census`／`--scope=windows`
   三枚调用点（尺＝`grep -nE 'scope=[a-z]+' .github/workflows/ci.yml`，件 `logs/ci2.rc.txt` 那次输出）；
   cmd/wisp 那发是 **`scripts/wisp-cli-tests.sh`** 在 **test-windows 那枚 job 内**起的
   （`portable-tests.sh:258` 注释逐字 "Only scripts/wisp-cli-tests.sh may call this"、`:261` `scope=(./cmd/wisp/)`）。
   ⛔ 只是措辞——要紧在**它和本票新分母在同一枚 job、同一台 runner**，所以"与本票无关"是一枚**需要凭据的结论**；
   我这发给了它凭据（§4 的时序尺＋名册存在尺），结论同向成立。
2. **"新增 3 枚红"应读成"3 枚 `PASS`/`SKIP` → `FAIL`"**：三枚名在基线 sha `cc315261` 都已存在（各 1 枚文件），
   基线色是 `--- SKIP (0.00s)` 与 `--- PASS (0.97s/1.25s)`。"新增"会让人少算一层（以为是新代码带进来的判据）。
3. **"census totals 两发逐字相同"这句没写它落在哪枚档**：它落在 **test-windows** 的日志里。
   我这发另下的两份 **test-core** 日志 `grep -ci census`＝**0**；顺带一枚票面"现量"节的行号腐烂＝
   `--scope=census` 的调用点现在在 `:761`（票面写 `:744`），差 **＋17**＝`AC#5` 那一笔新增的行数。⇒ 内容锚＝
   `run: bash scripts/portable-tests.sh --scope=census`；又一枚"行号锚⛔ 扛得住任何改动"的实例（本仓第四次？记我口径⛔ 数枚数）。
4. **派单"我已跑的"四把我全打中**（名册 7 枚／非注释 `+/-` 行 0／可执行键新增 0／census totals 逐字未变），
   `yaml.safe_load` 那把我用 `python -c` 复跑到 `True`＋job 数 6↔6——**这一把我打中了**，
   只是我还另加一把更狠的（剥注释行整文件逐字节相同）；⛔ 我的复跑改变他那四把的读法。
5. ⛔ 有一处我用工作树量而该用 blob：本件每一枚名册/文案/数组读数都标了尺形；
   只有四枚门禁件是工作树（`go`／`sh` ⓿法只能跑工作树），且逐枚 `cmp` 证与 `HEAD` blob 相同（§0 末）。
