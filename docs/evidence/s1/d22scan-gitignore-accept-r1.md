# d22scan gitignore 修 · 非实现者对抗验收 r1

**我是谁**：裁决者（非实现者）。本件不修任何东西，只判、只把最小修法点名交回。
**被验交付**：`3bb99aa`（代码）+ `364b9eb` / `00f9f65` / `7b4e212`（实现程自己的取证件 `docs/evidence/s1/d22scan-gitignore-r1.md`）。
**本件写作用锚点**：`42d141d`（`docs(台账,A212)`）。逐节标注该节读数时的锚点——HEAD 在我脚下动。

| 用途 | 锚点 | 说明 |
|---|---|---|
| 交付锚点（被验代码最后一次落地） | `3bb99aa` = `3bb99aad02358b76f2cfd5a155f3ca017c2f54d8` | 其父 = `467f8a4` |
| 我的读数锚点 ① | `ab785bf` | 五发读数 / 工作树对净快照 / 两枚变异 |
| 我的读数锚点 ② | `42d141d` | 行号重抽 / `git log 7b4e212..42d141d -- tools/d22scan` **现量空** ⇒ 交付后无人再动这枚仪器，`ab785bf` 的读数对 `42d141d` 仍成立 |

**仓库外工作区**（全部临时件在此，真树未写一字节）：
`D:\tmp\d22scan-gitignore-accept-r1-snapshots\`
- `head\` = `git archive ab785bf | tar -x`（净快照，无 `.git`）
- `pre\` = `git archive 3bb99aa^` = `467f8a4`（修前代码，无 `.git`）
- `headbuild\` = `git archive` + 真实 `frontend/dist/` 构建产物（无违规）
- `forcedfull\` = `git archive` 后**在仓库外 `git init`** 并 `git add -f frontend/dist/tracked-forced.tsx`（§1 的眼罩变异）
- `forcedfull-ci\` = 对 `forcedfull` 再跑一次 `git archive HEAD`（＝CI 那一侧的形状：没有 `.git`，但那枚受追踪的违规件**在场**）
- `forced\repo` + `forced\ci` = 迷你版同形实验（被 `checkRoot` 响亮拒绝，见 §3）
- `probe\` = `head/tools/d22scan` 的副本 + 我新增的 `zz_accept_probe_test.go`（七发语义探针，§2）

**日志留档**：`head-d22scan.log` / `pre-d22scan.log` / `head-seeded-runtests.log` / `headbuild-runtests.log` / `run-pre-seeded.log` / `run-head-seeded.log` / `run-forced-pre.log` / `run-forced-post.log` / `run-worktree-pre.log` / `run-worktree-post.log` / `head-roster.txt` / `pre-roster.txt`（均在上面那枚 `D:\tmp` 目录里）。

---

## 0. 交付物核对（锚点 `42d141d` 现跑）

`git show --numstat 3bb99aa` 与票面一致，**只有三枚文件**：

```
420	0	tools/d22scan/gitignore.go
82	9	tools/d22scan/main.go
231	0	tools/d22scan/scan_test.go
```

`git diff --name-status 3bb99aa^..3bb99aa` = `A gitignore.go` / `M main.go` / `M scan_test.go`。
三枚取证件 commit（`364b9eb` `00f9f65` `7b4e212`）各自只碰 `docs/evidence/s1/d22scan-gitignore-r1.md` 一枚路径，
`273/0`、`5/3`、`6/2` —— 没有一枚顺手改了别的文件。
`git log --oneline 7b4e212..42d141d` = 4 枚（`ab785bf` `827b797` `5540087` `42d141d`），其中
`git log -- tools/d22scan` **现量空输出** ⇒ 交付之后无人再动仪器，本件的行号不会因为别人又改这枚文件而漂。

票面背景核对：`A207`（缺陷本体）/`A209`②/`A210`③/`A211`② 均在 `docs/reports/pending-and-issues.md:5842,5860,5870,5878` 现量存在，
`A211` 已把台账基线重标做完（票面说"归编排者那一批"——那一批已交，`ab785bf`）。**这一条不构成缺陷。**

---

## 1. 主攻击：这枚 matcher 判"被忽略"用的是模式，还是受追踪状态

### 1.1 代码事实：纯模式，从不查受追踪状态

- 入口 `func (g *gitIgnore) skip(path string, isDir bool) bool` = `tools/d22scan/gitignore.go:126`；
  它唯一的输入是 `g.rel(path)` + `g.decide()`。
- `decide` → `match`（`:206`）→ `ruleMatch`（`:226`）→ `rulesFor`（`:263`）。
- `rulesFor` 读的是 `path := filepath.Join(g.root, filepath.FromSlash(dir), ".gitignore")`（`tools/d22scan/gitignore.go:271`）。
  整枚文件的全部数据来源 = **被扫树里的 `.gitignore` 文件本身**。
- 全文没有一处读 index、没有一处 `os/exec`、没有任何"这枚在不在 `git ls-files` 里"的判断。
  `import` 块（`tools/d22scan/gitignore.go:67-74`）只有 `fmt os path/filepath sort strings unicode/utf8`。

⇒ **答案：只按模式判，根本不查受追踪状态。** git 的真实判据是"受追踪的文件永远不会被忽略"
（`git ls-files -i -c --exclude-standard` 就是这条的官方仪器），这枚实现把这条丢了。

### 1.2 仓里**现在**有没有这种件：没有（两枚独立口径）

```
$ git ls-files -i -c --exclude-standard          -> 空（0 行）
$ for f in $(git ls-files); git check-ignore --no-index -q "$f" && echo MATCH;  -> 0 行 MATCH
```

两枚口径都是 0 ⇒ 今天没有一枚"受追踪且模式命中 ignore"的文件，工作树读数因此没有说谎。
**造成它的动作只有两种**：①`git add -f <被 ignore 的路径>`；②先提交、后补 `.gitignore` 规则（更常见、更 accident-prone）。
本仓 `frontend/.gitignore:12` = `dist/*`、`:13` = `!dist/.gitkeep`；根 `.gitignore:22` = `frontend/dist/*`、`:23` = `!frontend/dist/.gitkeep`。
`!` 只救 `.gitkeep` 一枚，`frontend/dist/` 下**任何其它名字**都在模式射程里。

### 1.3 我把那枚形状建出来了（仓库外），它把门变成了眼罩

在 `forcedfull\`（`git archive ab785bf` 解出的完整树，仓库外 `git init -b main`，`git add -A` 后**再** `git add -f frontend/dist/tracked-forced.tsx`）：

```
$ git ls-files -i -c --exclude-standard
frontend/dist/tracked-forced.tsx          <- git 自己承认：受追踪 且 模式命中 ignore
$ git archive HEAD | tar -x -C forcedfull-ci   # 造 CI 那一侧的形状
$ cat forcedfull-ci/frontend/dist/tracked-forced.tsx
export const doc = "ready ≤";
```

同一棵 `forcedfull-ci`，两版扫描器（`-root` 指过去，其余一字不差）：

**修前（`pre\tools\d22scan`，代码 = `467f8a4`）**

```
rc=1
d22scan: scope ban #6 frontend/         examined  41 text files
d22scan: scope ban #8 frontend/         examined  41 text files
frontend/dist/tracked-forced.tsx:1: [emoji] ban #8 glyph in scope frontend/ is banned (D23): ...
d22scan: 1 finding(s); D22 bans are not negotiable (...)
```

**修后（`head\tools\d22scan`，代码 = `ab785bf`）**

```
rc=0
d22scan: skipped as git-ignored: 1 file(s) under 0 ignored director(ies) [], decided by frontend/.gitignore (1 path(s))
d22scan: scope ban #6 frontend/         examined  40 text files
d22scan: scope ban #8 frontend/         examined  40 text files
d22scan: clean - no D22 ban violations; live scope work: ... ban #6 frontend/=40, ban #8 frontend/=40, ...
```

⇒ **一枚 git 承认受追踪、`git archive` 真的搬进 CI 的交付文件里的 `≤`，修前红、修后绿。**
这是本修**新引入**的检测面回退：修前 CI 检出的每一个字节都被扫（CI 树里天然没有"未追踪"件）；
修后 CI 树里凡模式命中的路径都不进分母、不读、且 `rc=0`。stdout 上有一行 note，**exit code 不受它影响**。

### 1.4 交付文本里那句保证是**假**的，而且它点名的那枚测试钉不住它

`tools/d22scan/main.go:61-64`（HEAD 现号）逐字：

> `// NOT: it does not suppress a finding on a TRACKED path (an exclusion that could`
> `// hide a violation would be the D22 run-away this tool exists to catch, so`
> `// scan_test.go's TestWalksSkipGitIgnoredPaths seeds the ignored and the tracked`
> `// copy of one byte and demands the first stay silent and the second go red), and`

`TestWalksSkipGitIgnoredPaths`（`tools/d22scan/scan_test.go:1263`）里"tracked 那一侧"种的是
`frontend/dist/.gitkeep`（`scan_test.go:1275`）与 `frontend/src/live.tsx`。
`.gitkeep` 之所以被扫，是**因为 `!frontend/dist/.gitkeep` 那枚否定把它捞回来了**，不是因为任何"受追踪"事实——
那枚夹具里根本没有 index，也没有 `.git`，它没有能力表达"受追踪且模式命中"这一支。
测试断言 `emoji` 名单恰为 `{frontend/dist/.gitkeep, frontend/src/live.tsx, internal/pkg/keep.go}`
（`scan_test.go:1315`）—— 三支全是"模式不命中或命中后被否定"，**没有一支是 §1.3 那个形状**。

而且第二份手抄副本 `ignoredLikeGit`（`tools/d22scan/scan_test.go:1458`）写的是：

```go
if strings.HasPrefix(rel, "frontend/dist/") {
    return filepath.Base(rel) != ".gitkeep"
}
```

⇒ 它对 `frontend/dist/tracked-forced.tsx` 返回 **true**（跳过）。两份副本共用同一条**纯模式**前提，
所以 `TestLedgerCountsMatchAnIndependentWalk` 在这种树上两边一起跳、一起绿。
**本批没有任何一枚测试能发现 §1.3 那发变异**——这恰好是 `main.go` 那段注释自我担保的地方。

**最小修法（不在本件射程内动，交回实现程/编排者拍）**：
`skip()` 前先问"这棵树是不是一次导出"——`root` 下没有 `.git` 目录时**不应用任何 ignore 规则**。
这条对本修的目标是**无损**的：净快照今天的 note 是空的（§6.4 现量），
即"CI 树里被模式命中的路径 = 0 枚"，所以按 `.git` 存在与否加闸，CI 读数一字不动（还是 40），
工作树读数也一字不动（`.git` 在，照旧排除构建产物，还是 40），而 §1.3 那发重新变红
（`git archive` 出来的 CI 树里那枚 `tracked-forced.tsx` 会被读）。
代价必须一起说：`scan_test.go:1458` 的 `ignoredLikeGit` 得在**同批**接上同一条 `.git` 闸门，
否则净快照上两份副本会当场打架——那正是这枚副本设计 wanted 的摩擦，不是 bug。
更便宜的底线（无行为变更也必须做）：`main.go:61-64` 那句"does not suppress a finding on a TRACKED path"
要么兑现，要么删；现在它是这枚仪器上**未被任何断言钉住的担保**。

---

## 2. 语义清单：逐条"处理了吗 / 有测试钉吗"

"钉"= `tools/d22scan/scan_test.go` 里存在断言这条语义的用例（现号 = `42d141d`）。
探针 = 我在 `probe\zz_accept_probe_test.go` 里问给**已交付的 matcher** 自己的七发（真跑，`go test -run TestAcceptProbe -v -count=1`，全 PASS 因为我只 t.Logf 不硬判）。

| 语义 | 处理？ | 钉？ | 证据 |
|---|---|---|---|
| 否定 `!pattern` | 是（`gitignore.go:301-304` 置 `negate`；`:248` `ignored = !r.negate`） | **是** | `scan_test.go:1358` `TestGitIgnoreRuleSemantics` 的 `!dist/.gitkeep` / `notes/blocked.md` 两支 |
| 只配目录 `dir/` | 是（`:305-308` 剥尾斜杠置 `dirOnly`；`:237-239` `r.dirOnly && !isDir` 跳过） | **是** | `TestGitIgnoreRuleSemantics` 的 `src/build`（文件）vs `web/build`（目录） |
| 锚定 `/foo` | 是（`:309-311` 前导 `/` → `anchored`） | **否** | `TestGitIgnoreRuleSemantics` 的 9 条规则里**没有一枚以 `/` 开头**；`/build` 一支未测 |
| 不锚定 `foo/`（任意深度） | 是（不含 `/` 则按 base name 匹配） | **是** | `web/build`、`web/logs`、`sub/y.log` 三支 |
| 模式内含内部 `/` ⇒ 锚到该规则文件目录 | 是（`:312-314` `strings.Contains(line,"/")` → `anchored`；`:241-245` `segs[r.dirSegs:]`） | **是** | `a/b/mid.txt` 与 `c/a/b/mid.txt` 正反两支 |
| `**` 跨路径段 | 部分（`:330-337` 段级 `**`；`a/**/b` 与 `**/gen/*.go` 对） | **是**（跨段这支） | `sub/gen/a.go`、`gen/a.go`（零段）两支 |
| **`foo/**` 会不会把 `foo` 目录本身剪掉** | **否，剪错了** | **否** | 探针 P2 现量：规则 `dist/**` + `!dist/.gitkeep` ⇒ `decide("dist",isDir=true)=ignored=true`，`.gitkeep` 跟着被吞。git 的 `dist/**` **不**匹配 `dist` 本身，否定因此还能救。⇒ 一旦有人写下这种规则，CI 的 40 会变 39（与本文件 `:50` 自己写的顾虑同一枚，只是没防住） |
| 尾随 `*` | 是（`:353-363` `segMatch` 的 `*` 分支；连续 `*` 折叠成一枚 = git 文档那条） | **是** | `dist/*`、`dist/assets`、`*.log` 三支 |
| `#` 注释行 | 是（`:297` TrimSpace 后判前缀） | **是** | `# a comment line` + `{"comment", false}` 两支 |
| 空行 | 是（`:297`） | **是** | 夹具里那枚空行 |
| 转义 `\#` | 是（`:385-389` `case '\\\\'`） | **否** | 探针 P3：`decide("#frag/a.txt") ignored=true`，与 git 同；无测试 |
| **行首空格** | **否——被 TrimSpace 掉了** | **否** | `gitignore.go:296` `strings.TrimRight(strings.TrimSpace(raw), " \t")`；git **不**剥行首空白。探针 P1 现量：规则 `  dist/*`（git 意为名叫 `"  dist"` 的路径）⇒ `decide("dist/index.html")=ignored=true`。**方向 = 多吞 = 眼罩向**。今天两枚 `.gitignore` 里 `grep -P '^\s+\S'` 现量 0 行，所以还不成害 |
| **逐目录 `.gitignore`（含 `frontend/.gitignore`）** | 是 | **是** | `rulesFor(dir)` 对每个祖先目录各解析一次（`:263-282`，读文件在 `:271`）；"深枚文件优先"由 `ruleMatch` 从深到浅（`:229-231`）实现。测试 `notes/.gitignore` 的 `!blocked.md` 压过根文件那一支（`scan_test.go:1358` 体内）。真树两侧都读到了：note 现量同时点名 `frontend/.gitignore` 与 `.gitignore`（§6.3） |
| `.git/info/exclude` | 否（明示不支持，`:61-64`） | n/a | 探针 P6：`ignored=false` ⇒ 方向是**多扫**。但代价要说清：用 info/exclude 挡住的构建产物会继续进工作树分母而 CI 里没有该件 ⇒ **A207 那枚"两形状不同数"在这条支上原样活着** |
| `core.excludesFile` | 否（同上） | n/a | 同上；无测试 |
| 被忽略目录连内容一起剪 | 是（`match` 先逐级问父目录，`:209-218`） | **是** | `dist/assets/b.js`、`web/build/x.go` 两支 |

`gitignore.go:61-64` 那句 "Every unsupported shape fails in the loud direction: the path stays scanned,
so the worst this matcher can do is examine too much, never less" ——**这句对自己文件不成立**：
行首空白（P1）与 `foo/**`（P2）两条未支持形状都是**少扫**方向，而且都写在它自己的射程声明里
（`:44-60` 列了它「实现的语义」，没列这两条的例外）。这是本件第 2 枚实质缺陷。

---

## 3. 树不是 git 仓库时会怎样（CI 复制件 / `-root` 指到别处）

**没有"是不是仓库"这个概念**：`gitignore.go` 从不看 `.git`（全文对 `.git` 的引用只有注释 `:38-39` 与
`main.go` 里 walk 的 `d.Name() == ".git"`）。分三种实际形状，都真跑过：

1. **没有 `.git` 但 `.gitignore` 在**（`git archive` 形状，也就是本仓 CI 与所有取证件用的形状）：
   照常套用规则。探针 P5 + §6.4 现量：净快照里 note 为空 ⇒ 树内无模式命中件，读数与修前逐字相同。
2. **树里根本没有 `.gitignore`**：探针 P5 `decide(...)=false`、`note()=""` ⇒ **什么都不排除、什么都不说**。
   方向 = 全扫 = 响亮（多扫不误判，不会把违规藏掉）。这是安全的那一支。
3. **`-root` 指到不是一枚 wisp 仓的地方**：响亮（仪器自己的码是 2；
   我的迷你件（只有 `.gitignore` + `frontend/`，没有 `internal/`、`cmd/`、`go.mod`、`allowlist.txt`）
   实测被 `checkRoot` 拒绝：`go run . -root forced/ci` → 无任何 scope 行 +
   `is not a wisp repository root: missing go.mod`（`main.go:1124-1129` 的 required 列表）。
   ⚠ 我初稿在这里写的是"shell `rc=2`"，那是**我的引用口径错**，`§11.1` 已更正：`go run` 把子进程的 2 拍平成 1，
   仪器自己的码 2 要用仓库外二进制才看得见。
   补出 10 枚以下 production `.go` 也一样拒绝（`minProductionGoFiles = 10`，`main.go:1117`）。
   ⇒ **"静默全放行"这个危险答案不是本实现的行为**。

---

## 4. 是 shell out 还是重写？重写算不算重复造轮子？

- **重写**：`tools/d22scan/gitignore.go:67-74` 的 import 里没有 `os/exec`，全仓 `tools/d22scan/*.go` 里
  `exec.` 现量 0 命中。⇒ 不调 `git`，因此"`git` 不在 PATH 会怎样"这一问**不适用**——
  这条选择顺带把"仪器是否可用"和"这台机器有没有 git"解耦了，票面给的理由（净快照没有 `.git`，
  用它就把读数重新绑回机器状态）我在 §3 形状 1 上独立复现成立。
- **仓里是否已有一枚做同一件事的依赖**：`tools/d22scan/go.mod` 全文只有两行
  （`module github.com/CarlosShao/wisp/tools/d22scan` + `go 1.27`），**零 require**；
  根 `go.mod` 的 require 与 `go.sum` 里 `grep -i 'gitignore|go-git'` 现量 0 命中。
  ⇒ 420 行**没有**重复任何已 vendored 的东西。反过来说：这枚独立 module 想引第三方也必须新拉依赖，
  而 `vendor/` 在本仓是被忽略的（`.gitignore:9`），所以"用库"在本仓并不可得。
- `TestGitIgnoreRuleSemantics`（`scan_test.go:1358`）直接调 `newGitIgnore` / `g.decide`，
  是包内白盒测试，不需要真 `git` ⇒ 语义可证伪性不依赖外部程序。这一条我认可。

---

## 5. 另一根轴：有没有碰它不许碰的契约

### 5.1 ban #8 字符类：**逐字节未动**

```
$ for a in 3bb99aa^ 3bb99aa HEAD; do git show $a:tools/d22scan/main.go \
      | grep 'emojiRe = regexp.MustCompile' | sha256sum; done
3bb99aa^  len=135  sha256[:16]=a6150ef689cc20c9
3bb99aa   len=135  sha256[:16]=a6150ef689cc20c9
HEAD      len=135  sha256[:16]=a6150ef689cc20c9
```

三枚锚点上同一串，且与票面期望值逐字相同（HEAD 现号 `tools/d22scan/main.go:149`）：

```go
var emojiRe = regexp.MustCompile(`[\x{1F000}-\x{1FAFF}\x{2200}-\x{22FF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}\x{FE0F}\x{1F1E6}-\x{1F1FF}]`)
```

`git diff 3bb99aa^..3bb99aa -- tools/d22scan/main.go` 里 **没有任何一行 `+`/`-` 触及 `emojiRe` 或 `1F000`**
（`grep -E '^[+-].*(emojiRe|1F000)'` 现量 rc=1）。⇒ 这一条不构成退回项。

### 5.2 `allowlist.txt`：`git diff 3bb99aa^..3bb99aa -- tools/d22scan/allowlist.txt` 输出 **0 字节**，
且 `git diff --name-status 3bb99aa^..3bb99aa` 根本不列它。⇒ **零豁免新增**，成立。

### 5.3 `main.go` 那 9 行删除：逐枚点名（旧号来自 `git diff -U0 3bb99aa^..3bb99aa`）

| # | 旧号 | 删除的原文 | 有没有"带分母的行为"被摘掉 |
|---|---|---|---|
| 1 | 30 | `//	                      The walk now examines 35 text files in this repo` | 注释，无分母 |
| 2 | 31 | `//	                      (measured 2026-09-21 on a \`git archive HEAD\` snapshot)` | 注释，无分母 |
| 3 | 349 | `"frontend/ that is not node_modules/ or testdata/ - .tsx, .ts, .mjs, .css, " +` | ban #6 的**说明文字**，替换成 `"… or testdata/ or a gitignored path (A207) - …"`（`main.go:375`），射程只增不减 |
| 4 | 351 | `"the ban into (35 files as measured 2026-09-21).",` | 同一串说明，替换为 `"(35 files as measured 2026-09-21; 40 as measured 2026-09-25 at 467f8a4, …)"`，旧数**保留在文里**（不是抹掉） |
| 5 | 594 | `if d.Name() == "testdata" \|\| d.Name() == ".git" {` | walkGo 目录判定 → `:633` 变为 `… \|\| s.ign.skip(path, true)`：**三支全在，追加第四支** |
| 6 | 778 | `if d.Name() == "testdata" \|\| d.Name() == "node_modules" \|\| d.Name() == ".git" {` | walkText → `:823` 同上，三支原样保留 |
| 7 | 849 | `if d.Name() == "node_modules" \|\| d.Name() == ".git" {` | walkEmoji → `:901` 同上，两支原样保留 |
| 8 | 1081 | `if d.Name() == "testdata" \|\| d.Name() == ".git" {` | checkRoot → `:1146` 同上 |
| 9 | 1086 | `if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {` | checkRoot 计数条件 → `:1151` 追加 `&& !ign.skip(path, false)`，前两支原样 |

⇒ **没有一枚删除摘掉过"带分母的行为"**：5–9 是 `\|\|`/`&&` 单向加严（剪掉更多，不是更少），1–4 是测量文字。
"摘掉任意一味是否存在一发变异从此打不红"这一问，答案落在**新增的那一味**上（§1.3 那发），不是删除上——
即：本修的红/绿差异不是"少了一条排除"造成的，而是"多了一条会吞掉交付件的排除"造成的。

### 5.4 断言有没有放水（本仓口径只有两问：断言被改过吗 / helper 是否为通过而新造）

- **删除行数**：`scan_test.go` `+231 / -0` ⇒ 旧断言一枚未删、一枚未改号。
- **新增 `t.Skip` / 阈值 / golden**：`git show 3bb99aa | grep -E '^\+.*(t\.Skip|SKIP|threshold|golden)'` 现量 **rc=1（零命中）**。
  新增两枚测试用的是 `t.Errorf`/`t.Fatal`（更强），不是放宽。
- **唯一一处"断言逻辑被改"的地点** = `TestLedgerCountsMatchAnIndependentWalk` 的对照走查里插进了
  `if ignoredLikeGit(relToRepo(root, p)) { return nil }`（`tools/d22scan/scan_test.go:1197`）。
  这一支是**给 oracle 加了排除**，性质上属于"两边一起改才能维持一致"，所以它是否有牙必须实测，不能靠自述。
  ⇒ 见 §6.8 那发"牙齿对照"：我在快照里往 `.gitignore` 追加一条落进被扫 scope 的规则，看这枚 oracle 是否当场打脸。
- **helper 是否为通过而新造**：`relToRepo`(`:1433`) 与 `ignoredLikeGit`(`:1458`) 确为本批新造，
  但它们的用途是给 oracle **加**约束（不是给被测方加豁免），且实现程在 `scan_test.go:1164-1173` 那段
  注释里明写了这是"故意的第二份手抄副本 +  disagreement 即打脸"。判定：
  **不是为通过而造，但它与实现共用同一条"纯模式"前提**（§1.4），所以它证伪不了 §1.3 那一支——
  这枚局限是实现程自述里没有点破的。

---

## 6. _required_ 读数（全部现跑，日志在 `D:\tmp\d22scan-gitignore-accept-r1-snapshots\`）

### 6.1 仪器自己那套的四数 + 名单枚数，与修前对比

`scripts/d22scan.sh`（`set -eu`）两步：**step 1 = 正向对照**（`sh tools/d22scan/runtests.sh -C tools/d22scan ./...`），
**step 2 = 真扫**（`go run . -root`）。四数**只由 step 1 产出**，八枚 per-scope 读数**只由 step 2 产出**。

| 锚点 / 树 | 命令 | PASS | === RUN | FAIL | SKIP | 顶层名单枚数 |
|---|---|---|---|---|---|---|
| 修前 `467f8a4`，净快照 `pre\` | `sh scripts/d22scan.sh` 全流程 | **24** | **64** | 0 | 0 | **24** |
| 修后 `ab785bf`，净快照 `head\` | `sh scripts/d22scan.sh` 全流程 | **26** | **66** | 0 | 0 | **26** |

两枚 `rc=0`。差值 **+2/+2/0/0**，与自述（24→26 / 64→66）一致。
名单差集（`diff pre-roster.txt head-roster.txt`）现量恰为两行、**只有新增没有消失**：

```
12a13
> TestGitIgnoreRuleSemantics
24a26
> TestWalksSkipGitIgnoredPaths
```

⚠ 诚实记一枚我自己量错形状的地方：我用 `grep -c '\[no tests to run\]'` 复算 step 1 日志得 **1**，
那枚 1 是 `runtests.sh` 自己的 OK 行里那句 `'[no tests to run]'=0` 的字面，**不是真实事件**；
`runtests.sh` 报的 0 才是对的。step 1 里"哪一步产出哪个数"因此必须按上面的归属读，不能按我的复算读。

⚠ 另一枚流程性事实（票面已写，本件复现）：step 1 红 ⇒ `set -eu` 让 step 2 **根本不跑**。
`head-seeded-runtests.log` 就是这一形状的真例：我种了违规探针后 step 1 rc=1（两枚 FAIL：
`TestScannerSelfScanOfRealRepoIsGreen` 报 `repo HEAD violates: frontend/src/still-bites.tsx:1`，
`TestRealRepoLedgerIsHonest` 报 `HEAD must be green, rc=1`），**"step 2 没输出"绝不能读成"step 2 过了"**。

### 6.2 still-bites（非忽略件仍必须打红）

树 = `head\`（净快照 + 真实 `frontend/dist/` 构建产物 3 枚 + `frontend/dist/probe.txt`(≤) + `frontend/src/still-bites.tsx`(≤)）。
扫描器 = `head\tools\d22scan`（修后代码）。

```
rc=1
d22scan: skipped as git-ignored: 2 file(s) under 1 ignored director(ies) [frontend/dist/assets/], decided by frontend/.gitignore (3 path(s))
d22scan: scope ban #6 frontend/         examined  41 text files
d22scan: scope ban #8 frontend/         examined  41 text files
frontend/src/still-bites.tsx:1: [emoji] ban #8 glyph in scope frontend/ is banned (D23): non-comment text, string literals included; comments are exempt per Q-46(c)
d22scan: 1 finding(s); D22 bans are not negotiable (...)
```

⇒ **成立**：不被忽略的 `≤` 照样 rc=1，点名的 scope 是 `frontend/`，分母随种子涨到 41。

### 6.3 no-longer-counted（同一棵树、修前报 / 修后不报）

同一串字节、同一份种子的 `pre\`（修前代码）：

```
rc=1
d22scan: scope ban #6 frontend/         examined  45 text files
d22scan: scope ban #8 frontend/         examined  45 text files
frontend/dist/probe.txt:1: [emoji] ban #8 glyph in scope frontend/ is banned (D23): ...
frontend/src/still-bites.tsx:1: [emoji] ban #8 glyph in scope frontend/ is banned (D23): ...
d22scan: 2 finding(s); ...
```

修前 45 / 2 条 finding，修后 41 / 1 条。**finding 名单差集恰为 `frontend/dist/probe.txt` 一枚**
⇒ "被 git 忽略的路径不再算进分母、不再被读"这条**行为变更是真的**，而 §6.2 那枚仍红说明它**目前**不是眼罩。
分母差 45−41 = 4 = `index.html` + 2 枚 bundle + `probe.txt`（`.gitkeep` 因 `!` 否定留在 40 里），
与"CI 的 40 含 `.gitkeep`"这一实现约束对得上。

### 6.4 工作树 vs 净快照：八数并排（修后代码，读数锚点 `ab785bf`）

| scope | 真工作树（只读 `-root`） | 净快照 `head\` | 同数？ |
|---|---|---|---|
| bans #1-5 internal/ | 203 | 203 | ✅ |
| bans #1-5 cmd/ | 22 | 22 | ✅ |
| ban #6 frontend/ | **40**（修前 43） | 40 | ✅（**本修要合流的就是这一枚**） |
| ban #7 internal/tools/ | 18 | 18 | ✅ |
| ban #8 design/ | **32** | **16** | ❌ **残留** |
| ban #8 frontend/ | **40**（修前 43） | 40 | ✅ |
| ban #8 internal/ | 405 | 405 | ✅ |
| ban #8 cmd/ | 39 | 39 | ✅ |
| examined production Go | 225 | 225 | ✅ |
| rc | 0 | 0 | ✅ |

修前对照（真工作树，`pre\` 代码）：`ban #6 frontend/=43`、`ban #8 frontend/=43`、`ban #8 design/=32`、其余同值、rc=0。
⇒ **A207 的成因复现成立**（差 3 枚、正是那 3 枚被忽略的构建产物）。
净快照修后 **不打印排除行**（`head-d22scan.log` step 2 全文无 `skipped as git-ignored`）⇒ 该沉默时沉默，成立。

### 6.5 design/ 残留：票面假设**证实**；交付文本与台账**都没有过度声明**（我草稿里那句判错了，已在正文收回）

现量（`42d141d`，只读）：

```
find design -type f                                        = 68
git ls-files design | <isTextFile 的 15 类后缀>            = 16
find design -type f | <同后缀>                            = 32
git status --porcelain --ignored design  的 !! 行          = 0
comm -13 <tracked 文本> <on-disk 文本>  按二级目录分组     = 16 design/doubao + 16 design/old
```

⇒ 磁盘上那 32 枚**全部是"未追踪且未被忽略"**，正正落在票面猜的 `design/old/`(16)、`design/doubao/`(16)。
   ⚠ 我初稿把这一行写成"32 = 16 受追踪 + 16 未追踪"，是 prose 错、与上面那行 `comm` 读数自相矛盾；
   重跑发现 HEAD 里那 16 枚受追踪 `design/**` 此刻**已被另一程从磁盘删除且尚未 commit**（16 枚 ` D`），
   与磁盘 32 枚的交集实测为 **0**。`§11.2` 已更正并追加 F9。
`!!` 行为 0 ⇒ 没有任何一枚是 ignore 过滤器能吞的；**ignore 过滤器不能也不该吞它**。
票面假设 **成立**，"worktree 恒等于 CI"这句话作为一般命题是**假的**（八枚里 7 枚同，`ban #8 design/` 不同）。

- **交付文本自己**：`docs/evidence/s1/d22scan-gitignore-r1.md:186-189` 明写
  "`ban #8 design/` 在这棵工作树上是 **32**、净快照是 **16**，本修**不动它，也不该动**"，并给出 `!!` 行=0 的现量。
  ⇒ **实现程没有过度声明**，这一条不记在它头上。
- **代码注释**：`main.go:529` "and 40 is the reading of both shapes" 所在段通篇只谈 `frontend/`
  两枚数，语境里成立；`main.go:30-36` 同理。**但** `gitignore.go:33-34`
  "`.gitignore` IS the policy; the scanner follows it, so the instrument and the repository's ignore rules cannot disagree"
  是**一枚过度声明**：git 的 policy = 模式 **+ index**，实现只跟了前者（§1.3），所以"cannot disagree"当场就不成立。
- **台账 A211**：**我这条草稿原本要判它过度声明，现跑之后收回**。
  `docs/reports/pending-and-issues.md:5878` ② 末句 "**这道门的分母现在在工作树与干净检出上是同一个数。**"
  单看是歧义的，但同一格 ③（`:5879`）逐字写着
  "**一枚残留要说清，别让下一位以为'从此两边永远相等'**：工作树 `ban #8 design/=32` vs 净快照 `16`……
  '跳 ignore' 治不了它，这是另一个因"。⇒ **台账没有过度声明，它自己就把射程划完了**；
  ④（`:5880`）还预先登记了本件这发主攻击（"被跟踪的文件永远不会被 git 忽略……`git add -f` 的
  `frontend/dist/*.txt` 会对门永久隐身＝把量尺换成蒙眼布"）。
  ⇒ 本件的结论是：**编排者猜中了形状，但交付里没有任何一枚测试或一行代码防住它**（§1.3/§1.4）。

### 6.6 note() 的自述数字与移动的分母不是同一个量纲（现量两例）

`gitignore.go:173` 的模板：`skipped as git-ignored: %d file(s) under %d ignored director(ies) [%s] …`。

| 树 | note 说的 file(s) | 分母实际移动 | 差额来源 |
||---|---|---|
| §6.2 seeded `head\` | 2 | 45 → 41（**4** 枚） | `assets/` 里 2 枚 bundle 被 SkipDir 剪掉，**从未被 visit**，不进 `files` |
| 真工作树（§6.4） | 1 | 43 → 40（**3** 枚） | 同上 2 枚 bundle |
| §1.3 forcedfull-ci | 1 | 41 → 40（1 枚） | 无被剪目录，恰好对得上 |

⇒ `files` 记的是**walk 访问次数**，`pruned` 记的是目录名，模板那句 "N file(s) **under** M ignored director(ies)"
把这二者读成了包含关系。实现程在 `d22scan-gitignore-r1.md:127-130` 自捉过**另一种**错（两份 walk 重复计数，
`counted` 去重后修好），但这一枚"剪掉的目录里不再计数"没修也没登记。
最小修法：note 里对每个被剪目录补一句它里面**还有多少枚文件**（walk 已 visit 过该目录本身，`os.ReadDir` 一次即得），
或者把措辞改成 `N path(s) skipped individually, M ignored director(y|ies) pruned`——
这枚 note 存在的唯一理由就是"移动了的数必须能被解释"（`main.go:1302-1309` 那段自述），
所以它现在**只解释了三枚里的第一枚**。

### 6.7 构建产物存在时，正向对照还绿不绿（修后代码）

`headbuild\` = `git archive` + 真 `frontend/dist/{index.html,assets/×2}`，**不种违规**：

```
$ sh tools/d22scan/runtests.sh -C tools/d22scan ./...      rc=0
runtests.sh: OK - packages=[./...] top-level: PASS=26 FAIL=0 SKIP=0, === RUN=66, '[no tests to run]'=0
```

⇒ 有构建产物的树上 step 1 仍绿。附带一枚对我有利的否证：`head-seeded-runtests.log:116-122,130`
现量 `ban #6 frontend/=41`、`ban #8 frontend/=41`、独立对照 `verified ban #6 frontend/: 41 files`
⇒ 三方（两枚 walk + oracle）在**带构建产物**的树上给出同一个 41，等式那两枚守卫没被绕过。

### 6.8 那枚"故意留成第二份手抄"的 oracle 到底有没有牙（实测两发）

> **行号口径**（本件唯一一处跨版本引用，先说清免得被当成漂移）：
> §6.8 里 `scan_test.go:269`/`:1273`/`:1277` 三枚是**修前那份文件**（`3bb99aa^` = `467f8a4`）的现号，
> 来自 `blind-suite-pre.log` 的真回显；`:1239`/`:1504` 是**修后**（`ab785bf`/`42d141d`）的现号，来自 `blind-suite.log`。
> 两版 scan_test.go 差 `+231` 行，所以同一句断言在两边不可能是同一个号——这里贴的是日志，不是我重数出来的号。

不自述，直接改树跑 step 1 的三枚守卫。

**牙齿正向对照** — `headtooth\` = `git archive` + 往根 `.gitignore` 追加一行 `frontend/src/`
（一条**落进被扫 scope** 的新 ignore 规则）：

```
rc=1
--- PASS: TestRealRepoBan8CoversFrontendTreeAtBan6sCount (0.44s)
--- FAIL: TestLedgerCountsMatchAnIndependentWalk       (0.44s)
--- PASS: TestRealRepoLedgerIsHonest                  (0.41s)
    scan_test.go:1239: scope ban #6 frontend/ reported examining 19 files, independent walk says 40 - the number in the self-report is wrong, not just small
    scan_test.go:1239: scope ban #8 frontend/ reported examining 19 files, independent walk says 40 - the number in the self-report is wrong, not just small
runtests.sh: go test exited 1 - top-level: PASS=2 FAIL=1 SKIP=0
```

**牙齿负向对照** — `headtoothout\` = 同样的 archive，但追加 `/docs/` 与 `/models/`（scope 之外的规则）：

```
rc=0
--- PASS: TestRealRepoBan8CoversFrontendTreeAtBan6sCount
--- PASS: TestLedgerCountsMatchAnIndependentWalk
--- PASS: TestRealRepoLedgerIsHonest
runtests.sh: OK - top-level: PASS=3 FAIL=0 SKIP=0, === RUN=3
```

⇒ **oracle 真有牙，而且是精确的牙**：规则吃进被扫 scope ⇒ 当场红并说出"这是基线变更"；
规则在 scope 之外 ⇒ 不冤枉。`d22scan-gitignore-r1.md:57-62` 那段设计说明**兑现了**。

**但同一枚 oracle 对 §1.3 那一支完全无力**（这才是关键，实测不是推理）——
`forcedfull\` 是一枚**真的 git 仓库**，`git ls-files -i -c --exclude-standard` 现量点名
`frontend/dist/tracked-forced.tsx`（受追踪 且 模式命中），该文件里就写着 `≤`：

```
rc=0
--- PASS: TestScannerSelfScanOfRealRepoIsGreen        (0.46s)
--- PASS: TestRealRepoBan8CoversFrontendTreeAtBan6sCount (0.43s)
--- PASS: TestLedgerCountsMatchAnIndependentWalk      (0.47s)
--- PASS: TestRealRepoLedgerIsHonest                  (0.46s)
    scan_test.go:1504: ban #6 examined 40 frontend/ text files in the real repo
runtests.sh: OK - top-level: PASS=4 FAIL=0 SKIP=0
```

**同一棵树在修前的代码上**（`pre2\` = `git archive 3bb99aa^`，只加这一枚文件，别的都没动）：

```
rc=1
    scan_test.go:269: repo HEAD violates: frontend/dist/tracked-forced.tsx:1: [emoji] ban #8 glyph in scope frontend/ is banned (D23): ...
--- FAIL: TestScannerSelfScanOfRealRepoIsGreen (0.44s)
--- PASS: TestRealRepoBan8CoversFrontendTreeAtBan6sCount (0.39s)
--- PASS: TestLedgerCountsMatchAnIndependentWalk        (0.41s)
    scan_test.go:1273: ban #6 examined 41 frontend/ text files in the real repo
--- FAIL: TestRealRepoLedgerIsHonest (0.42s)
    scan_test.go:1277: HEAD must be green, rc=1
runtests.sh: go test exited 1 - top-level: PASS=2 FAIL=2 SKIP=0
```

⇒ **回退发生在 `scripts/d22scan.sh` step 1 自己身上**：修前两枚守卫红、step 2 也会红；
修后 step 1 的四枚守卫全绿（`ban #6 examined 40`），CI 的两步门**没有任何一层看得见这一发违规**。
所以 §1 那一条不是"读数变了"级别的问题，是"门从两层看变成零层看"。

### 6.9 同一枚过配打到 Go scope（比 `frontend/dist/` 更要命的一发）

根 `.gitignore:14` 是不锚定的目录规则 `build/`，射程直达 `internal/build/`。
两棵 fresh archive（`headmut\` 修后代码 / `premut\` 修前代码），只加同一枚文件：

```
internal/build/leak.go:  package build / func worker() {} / func leak() { go worker() }
$ git ls-files -i -c --exclude-standard   ->  internal/build/leak.go   (受追踪 且 模式命中)
```

**修前**：

```
d22scan: scope bans #1-5 internal/      examined 204 production Go files
internal/build/leak.go:5: [bare-goroutine] bare `go worker(...)` is banned (D22/D38b, R16: named calls count too): use observe.Registry.Spawn (...)
d22scan: 1 finding(s); ...
```

**修后**：

```
rc=0
d22scan: skipped as git-ignored: 0 file(s) under 1 ignored director(ies) [internal/build/], decided by .gitignore (1 path(s))
d22scan: scope bans #1-5 internal/      examined 203 production Go files
d22scan: scope ban #8 internal/         examined 405 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations; ...
```

⇒ 一枚 **AGENTS.md §1.2 明列的禁止形状**（裸 `go func`/`go worker`）在修后**对 bans #1-5 与 ban #8 双双不可见**，
`rc=0`，且 note 说的是 **"0 file(s)"** —— 整包被剪，note 一个文件都没报。
这条同时是 §6.6 那枚量纲缺陷的最坏实例：**被剪目录里可以藏着 N 枚文件而 note 报 0**。
实现程并非不知道这枚规则会打到 `internal/`：它的夹具 `scan_test.go:1265` 就往 `.gitignore` 里写了 `build/`、
并在 `:1272` 种 `internal/build/probe.go`——只不过**当作"该被排除的正面例子"来种的**，
没把"这枚件同时是受追踪交付件"那一支写进任何断言。

---

## 7. 同源副本普查：这修有没有造出第 7 份

票面给的命令跑完是 30 行命中，但**其中 12 行是注释、5 行是夹具种子字符串**，
真正的**可执行跳过逻辑**我按位置重新点一遍（HEAD 现号）：

| # | 位置 | 抄的是哪条规则 | 归属 | 本批是否动 |
|---|---|---|---|---|
| ① | `tools/d22scan/main.go:633` `walkGo` 目录名 | `testdata`/`.git`（+ 新接的 ignore） | 实现程 | **动了** |
| ② | `tools/d22scan/main.go:823` `walkText` 目录名 | `testdata`/`node_modules`/`.git`（+ ignore） | 实现程 | **动了** |
| ③ | `tools/d22scan/main.go:901` `walkEmoji` 目录名 | `node_modules`/`.git`（+ ignore） | 实现程 | **动了** |
| ④ | `tools/d22scan/main.go:1146` + `:1151` `checkRoot` | `testdata`/`.git` + `.go` 后缀判（+ ignore） | 实现程 | **动了** |
| ⑤ | `tools/d22scan/scan_test.go:1192` 独立对照走查 | `testdata`/`.git`/`node_modules` 目录名（+ `:1197` 调 `ignoredLikeGit`） | 实现程 | **动了** |
| ⑥ | **`tools/d22scan/scan_test.go:1458` `ignoredLikeGit`** | **`frontend/dist/` 前缀 + `.gitkeep` 例外——手抄的 gitignore 政策** | 实现程（本批**新造**） | **新建** |
| ⑦ | `internal/panel/frontend_hygiene_test.go:289` | 子串 `/node_modules/`、`/dist/` | **编排者/前端会话地界** | 未动（正确） |
| ⑧ | `internal/panel/composer_test.go:286`/`:434`/`:606` | case 表列 `node_modules`/`.git`/`dist`/`fixtures` | 同上 | 未动（正确） |
| ⑨ | `internal/risk/pathresolver_rewrite_account_test.go:223` | case 表列 `.git`/`.scratch`/`node_modules`/`dist`/`build`/`testdata`/`design`/`docs` | 同上 | 未动（正确） |

**判：有没有造出第 7 份？按"手抄的跳过规则"这个口径，答案是有一份新的 = ⑥ `ignoredLikeGit`。**
实现程报的"6 份"不是假账，但它的口径有两处会让人少算一枚：③ 那行把 `main.go` **两处**（594/1081）合成一枚 ③，
而 ⑥ `ignoredLikeGit` 被折进它自己的 ④ 里（`d22scan-gitignore-r1.md:268` 有写，表里没列）。
⇒ 台账口径建议按位置数（我这枚表是 9 行、其中新增 1 行）。
**这份新副本是披露的、不是偷藏的**（`scan_test.go:1164-1173` 整段论证 + r1 §1 脚注），而且 §6.8 证明它**真有牙**。
它的缺陷不在"多抄了一份"，在**抄的是同一条纯模式前提**（§1.4），所以它对 §1.3/§6.8 那一支天生看不见。

**归属边界是否被尊重：是。** `git diff --name-status 3bb99aa^..3bb99aa` 只有 `tools/d22scan/` 三枚路径，
`internal/**` 一字节未动；`grep -rn 'tools/d22scan' --include='*.go' internal cmd` 现量只有注释命中，
**没有任何 Go 文件 import 这枚 module** ⇒ ⑦⑧⑨ 不可能因本修变红（与 r1 §7 的自述一致）。

**⑦ 需不需要同批改？判：不需要，也不应该在本批改。** 三条理由：
(a) 它跳的是**子串** `/dist/`，本就比 `.gitignore` 宽，本修没让它变得更不一致；
(b) 它是 root module 的测试，与 tools/d22scan 零依赖；
(c) 它归编排者/前端会话，按票面"不许动别人地界"。
但必须登记一条**它自己那一族的既有暴露**（不是本批引入）：`internal/panel/frontend_hygiene_test.go:289`
同样会把一枚受追踪的 `frontend/dist/x.tsx` 从它自己的 ban #6/#8 对照里跳掉 ——
即 §1.3 那发在**第二台仪器**上本来就成立，且早于本修。这一条应记给⑦的主人，别记在本批头上。

顺带一枚与本批无关但普查路上撞见的**好消息**：`internal/panel/frontend_hygiene_test.go:71` 的
`emojiRangesRe` 与 `tools/d22scan/main.go:149` 的 `emojiRe` 我把两串正则抽出来做逐字比较，
**IDENTICAL**（`U1` 那一族同源字符类今天没有漂）。本批没碰它，也没有把它碰漂。

---

## 8. 缺陷清单 + 每枚的最小修法（本件一律不动手）

| # | 严重度 | 缺陷 | 证据 | 最小修法（点名，不代做） |
|---|---|---|---|---|
| **F1** | **阻断** | matcher 只按**模式**判忽略，从不查**受追踪状态**；于是 `git add -f`（或"先提交后补规则"）进来的交付件对门**永久隐身**，且 `scripts/d22scan.sh` **step 1 的正向对照也看不见** | §1.3 二进制两读（修前 rc=1/41、修后 rc=0/40）；§6.8 套件级四枚守卫全绿 vs 修前两枚 FAIL；§6.9 同一形状打到 bans #1-5（`internal/build/leak.go`，note 还报 "0 file(s)"） | `newGitIgnore`/`skip()` 前加一道**"这棵树是不是一次导出"**：`root` 下无 `.git` 目录时**不套用任何 ignore 规则**。对本修目标无损——净快照今天 note 为空（§6.4），即"CI 树里模式命中件 = 0"，所以加闸后 CI 仍 40、工作树仍 40，而 §1.3 那发重新变红。代价必须同批：`scan_test.go:1458` `ignoredLikeGit` 接上同一条 `.git` 闸门，否则两份副本在净快照上当场打架（那是设计要的摩擦，不是 bug）；再补一枚**能表达"受追踪且模式命中"**的断言（可用 `TestBuiltBinaryGoesRedEndToEnd` 那台真跑二进制的形状，在仓库外 `git init` 一棵 fixture 树）。 |
| **F2** | **阻断（文本面）** | 交付里两处**未被任何断言钉住的担保**：`main.go:61-64` "it does not suppress a finding on a TRACKED path"（它点名的测试实际是 `!` 否定捞回的 `.gitkeep`，夹具无 index，表达不了这一支）、`gitignore.go:33-34` "the instrument and the repository's ignore rules cannot disagree"（git 的 policy 是 模式 **+ index**，只跟了前半） | §1.4 代码事实；§1.3 实测直接反证 | 二选一，不许并存：把 F1 实现掉（担保随之为真），或把这两句改成可被证伪的真话（"按模式判；受追踪且模式命中的形状今天为 0 枚，`git ls-files -i -c` 为这条的口径"）。 |
| **F3** | 中 | 两条**少扫方向**的语义偏差，且 `gitignore.go:61-64` 自称"every unsupported shape fails in the loud direction: … never less"——**对它自己的文件不成立** | 探针 P1：规则 `  dist/*`（git 意为名叫 `"  dist"` 的路径）⇒ `decide("dist/index.html")=ignored true`；根因 `gitignore.go:296` 的 `TrimSpace`。探针 P2：规则 `dist/**` ⇒ `dist` **目录本身**被剪、`!dist/.gitkeep` 失效 ⇒ 写下这种规则的当天 CI 40→39（正是 `gitignore.go:50` 自己担心的那件事，没防住） | `:296` 改 `strings.TrimRight(raw, " \t")`（git 只剥尾随未转义空白）；`pathMatch` 的 `**` 分支加一条"`**` 不得吃掉整个待匹配路径的父目录"（即 `foo/**` 不匹配 `foo` 自身）。两条各配一枚 `TestGitIgnoreRuleSemantics` 的用例。 |
| **F4** | 低-中 | `note()` 的数与移动的分母不是同一个量纲：`gitignore.go:173` 那句 "N file(s) **under** M ignored director(ies)" 把 walk 访问数和目录名数读成了包含关系 | §6.6 三例：真工作树 note 说 1、分母动 3；seeded 树说 2、动 4；§6.9 说 **0 file(s)** 而整包被剪 | 对被剪目录做一次 `os.ReadDir` 把里面的文件数并进 note，或改措辞为 `N path(s) skipped individually, M ignored director(y|ies) pruned`。这枚 note 的唯一职责就是解释动过的数（`main.go:1302-1309` 自述），现在只解释了三枚里的第一枚。 |
| **F5** | 低（登记性） | `.git/info/exclude` / `core.excludesFile` 不支持（`:61-64` 已声明），但**没写它会让 A207 的病在这一支原样活着** | 探针 P6：`.git/info/exclude` 写 `dist/*` ⇒ `ignored=false` ⇒ 工作树继续数、CI 树里根本没有该件 ⇒ 两形状又不同数 | 文档补一句射程限定即可（不改行为）；或者在 `note()` 里报一声"检测到 `.git` 存在但本仪器不读 `info/exclude`"。 |
| **F6** | 低（信息） | `TestRealRepoBan8CoversFrontendTreeAtBan6sCount` 只断言两枚 walk **相等**，两边一起降时它不报 | §6.8 正向对照真回显：`19 vs 19` ⇒ 该枚 **PASS**，而基线已从 40 掉到 19 | 不是错（等式守卫本来的职责就窄），但它意味着**全仓没有一枚断言钉住 `frontend/` 的绝对基线**——所以台账重标（`A211`① 已做）就是唯一的锚，别把它省掉。 |
| **F7** | 信息 | 普查口径：本批确实**新增了一枚手抄副本** `ignoredLikeGit`，实现程的"6 份"是按 owner 数、不是按位置数（我按位置数是 9 处） | §7 表 | 台账里给 `ignoredLikeGit` 单列一行（它能独立于 ⑤ 移动）。 |
| **F8** | 不属于本批 | `internal/panel/frontend_hygiene_test.go:289` 用**子串** `/dist/` 跳，本来就会把受追踪的 `frontend/dist/x.tsx` 从它自己的对照里跳掉，且比 `.gitignore` 更宽（任何深度的 `/dist/`） | §7 (a)(b)(c) + 该行现号 | **记给⑦的主人（编排者/前端会话），早于本修存在**。本批不许动它，本件也不动。 |

### 8.1 本批**确实做到**的事（免得沉默被读成不必说）

- A207 的主目标**达成且被我独立复现**：工作树 `ban #6/#8 frontend/` 43→40，与三枚锚点上的 CI 读数 40 合流；CI 形状八数**逐字未变**、`rc=0`、净快照上不打印排除行。
- 契约面**零越界**：`emojiRe` 三枚锚上同一串（sha 前缀 `a6150ef689cc20c9`）、`allowlist.txt` 0 字节差异、`internal/**`/`frontend/**`/`design/**`/`docs/reports/**` 零字节、9 枚删除行无一枚摘掉带分母的行为、测试零删除零 `t.Skip` 零阈值。
- 规则是**读出来的不是抄来的**（`gitignore.go:271`），并且**故意留了第二份手抄副本当证伪器**，这枚证伪器被我用正/负两发实测证明**真有牙**（§6.8）。
- 不依赖 `git` 可执行文件，且 `tools/d22scan` 零依赖 module ⇒ 420 行**没有**重复任何已 vendored 之物（§4）。
- 取证件诚实：design/ 残留是它**自己先说出来的**（`d22scan-gitignore-r1.md:186-189`），note 重复计数那半枚缺陷也是**它自捉的**（`:127-130`），连自己 commit 标题写成 `A2207` 都在正文点名更正（`:225-228`）。

---

## 9. 总裁（一格一词）

| 格 | 判 | 一理由 |
|---|---|---|
| §1 主攻击：按模式还是按受追踪状态 | **退回** | 纯模式、无 index，`git add -f` 的交付件对门与对 step 1 **双双永久隐身**（§1.3/§6.8/§6.9 三发实测，非推理）。 |
| §2 语义覆盖 | **附条件成立** | 15 支里 11 支有实现且其中 9 支被测试钉住；**行首空白与 `foo/**` 剪父目录**这两支是少扫方向、无测试、且与 `:61-64` 的自述矛盾（F3）。 |
| §3 非 git 树 / `-root` 指错 | **成立** | 无 `.gitignore` 时不排除也不说（响亮），指错树由 `checkRoot` `rc=2` 硬拒；"静默全放行"不是本实现的行为。 |
| §4 shell out 还是重写 / 是否重复 | **成立** | 重写（import 无 `os/exec`），`tools/d22scan/go.mod` 零 require ⇒ 没有重复任何已 vendored 库；不用 `git check-ignore` 的理由被我独立复现成立。 |
| §5 契约轴（emojiRe/allowlist/删除行/断言） | **成立** | 字符类三锚同串、allowlist 0 字节、9 枚删除行全是加严或注释、测试零删零 Skip 零阈值。 |
| §6 读数（四数/两对照/工作树 vs 净快照） | **成立** | step 1 `24/64→26/66`、名单 +2 无消失；still-bites 仍 rc=1 点名 `frontend/`；no-longer-counted 差集恰为那枚 ignore 件；八数 7 同。 |
| §6.5 design/ 残留算不算缺陷 | **成立**（不算） | 现量 32 = 16 受追踪 + 16 未追踪且 `!!` 行 0，正落 `design/doubao`+`design/old` ⇒ ignore 过滤器不能也不该吞；交付文本与台账 A211③ **都没过度声明**。 |
| §6.6 note 自述精度 | **附条件成立** | 数出的是 walk 访问数不是文件数，最坏实例 "0 file(s)" 而整包被剪（F4）；不影响 rc，只影响可解释性。 |
| §7 同源副本与地界 | **成立** | 只碰自己三枚文件、root module 零 import 关系；新增那枚手抄副本是披露的且被证明有牙（口径应按位置数成 9 枚，F7）。 |
| **本批交付整批** | **退回** | F1+F2 一起才是要害：**门把自己的"我没看"重新命名成了"这里没有东西"**，而它为此写的两句担保恰好在它唯一没能表达的那一支上失效。 |

---

## 10. 我**没有**测什么（绝不让裁决者的沉默被读成批准）

1. **没有把 F1 的修法实现或验证**——我只证明它可被构造、且现网形状为 0 枚（`git ls-files -i -c --exclude-standard` 与 `check-ignore --no-index` 两枚口径都空）。修法在 `.git` 存在性加闸之后是否让八数仍然逐枚不动，**未量**。
2. **没有跑 root module 的任何测试**（`internal/**`、`cmd/**` 全 suite）：本修代码在独立 module、`grep` 现量无 import 关系，但"因此不可能影响 root module"是我的**推理**，不是我的读数。特别是 `internal/panel/frontend_hygiene_test.go`、`composer_test.go` 今天是否仍绿，**未跑**。
3. **没有跑完整 CI**（GitHub Actions 那套步骤），只跑了 `scripts/d22scan.sh` 这一个门的两步。
4. **没有测 Linux/macOS 行为**：全部在 win32 + `git 2.52.0.windows.1` + `go1.27.1 windows/amd64`。探针 P7（mode-000 的 `.gitignore`）在 Windows 上**仍能读**（`ignored=true`），所以"不可读规则文件 ⇒ 无规则 ⇒ 路径继续被扫"这条**在 Linux 上的响亮与否未证**。
5. **没有测并发/共享工作树场景**（本修是读文件的代码，我没查它在别人正在 commit 时读树的一致性）。
6. **没有测大小写/路径分隔符分支**：Windows 上 `filepath.Rel` 与 `ToSlash` 我依赖实现，未针对"盘符大小写不一致"、UNC、长路径（>260）构造用例。
7. **没有穷举 gitignore(5) 的其它细分支**：`!` 否定一个**目录**（`!logs/`）、`[!a-z]` 与 `[a-]` 这类畸形类、`**` 出现在段**内部**（`a**b`，代码按 git 文档折叠成单 `*`，我**没实测**）、以 `\` 结尾的转义尾随空白（`:61-64` 已声明不支持，未证方向）。
8. **没有验证 `design/` 那 16 枚未追踪件该不该进交付**——那是 owner 与前端会话的范围决定，本件只判"ignore 过滤器有没有权利吞它"（没有）。
9. **没有复算 37 这个历史读数**来自哪一枚锚点（票面只说 37/40/43 三个读数打架；我只量到 40 与 43，**37 未出现在我任何一发读数里**）。
10. **没有把 §6.9 那发做到"真实提交"级别**：我只在仓库外快照里 `git add -f` 造出形状并跑门；真仓 `dev` 分支上今天**不存在**这种件（第 1 条的两枚口径均为空），所以我不能说"CI 现在正在漏一个真违规"，只能说"门现在有这个洞、且 step 1 不再兜它"。
11. **没有测 `core.excludesFile` 被真的设过之后**的读数（只测了 `.git/info/exclude` 一支的 `ignored=false`）。
12. **权限**：本程**没有任一枚工具调用被权限系统拒绝**，因此也没有"因被拒而绕行"的事可报。
13. **本件没动过任何已存在的文件**、没 push、没 `git add -A`/`.`、没 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；每次提交前 `git diff --cached --name-only` 均只列 `docs/evidence/s1/d22scan-gitignore-accept-r1.md` 一枚路径（四次全中，无第二枚混入）。

---

## 11. 追加更正与追加发现（脚下 HEAD 与工作树状态都动过，这一节是现跑重导的）

写下 §1–§10 之后 HEAD 从 `42d141d` 走到 `686d7e7`，`design/` 的磁盘状态也被另一程改了。
按"每个 `file:line`/每个读数都必须从我真正读过的那一版重导"的规矩，这里改两处我写错的地方，并追加一枚发现。
**我未对别人那 16 枚未提交的删除做任何动作**（不 stage、不 commit、不恢复）；
本件四枚 commit 每次 `git diff --cached --name-only` 均只有 `docs/evidence/s1/d22scan-gitignore-accept-r1.md` 一枚路径。

### 11.1 更正 §3 第 3 支：我把 `go run` 的壳层退出码当成了仪器自己的码

`scripts/d22scan.sh` step 2 用的是 `go run . -root`，而 **`go run` 把子进程的 2 拍平成 shell 的 1**。现跑两遍：

```
$ go run . -root <非 wisp 树>
d22scan: "...\forced\ci" is not a wisp repository root: missing go.mod
d22scan: this tool is its own Go module; run `scripts/d22scan.sh` or ...
exit status 2
shell rc=1            <-- 壳层

$ ./bin/d22scan-acc.exe -root <同一棵树>       (仓库外自建二进制)
d22scan: "...\forced\ci" is not a wisp repository root: missing go.mod
binary rc=2           <-- 仪器自己给的码
```

⇒ §3 那句"响亮且 `rc=2`"里，**"响亮"成立**（两行错误文字 + 一行"该怎么跑"的提示，不是静默），
但**经 `go run` 时 shell 看到的数是 1**。`set -eu` 之下 1 与 2 都致命，门不会因此变绿，
所以这不是新缺陷，是**我引用读数时的口径错**；顺带一枚可登记的既有事实：
**CI 日志里 step 2 的 rc 永远区分不出"发现问题(1)"和"仪器没跑成(2)"**——这是 `scripts/d22scan.sh` 的形状，不是本修引入的。
本件其余各发凡是 `rc=0` 的读数不受影响（`go run` 不拍平 0）。

### 11.2 更正 §6.5 那句"32 = 16 受追踪 + 16 未追踪"

`42d141d`/`686d7e7` 现跑（`EXT` = `isTextFile` 那 15 类后缀）：

```
git ls-files design | <EXT>            = 16      （HEAD 里受追踪的文本件）
find design -type f | <EXT>            = 32      （磁盘上的文本件）
comm -12 tracked disk                  =  0      <-- 交集为空！
comm -23 tracked disk                  = 16      <-- 16 枚受追踪的**今天不在磁盘上**
comm -13 tracked disk                  = 32      <-- 磁盘上 32 枚**全部未追踪**
   按二级目录分组: 16 design/doubao + 16 design/old
git status --porcelain design          = 16 枚 " D"（未提交的删除） + 2 枚 "??"
git status --porcelain --ignored design 的 !! 行 = 0
```

⇒ 正确表述是：**磁盘上那 32 枚 100% 是"未追踪且未被忽略"**（`design/doubao/` + `design/old/`），
而 HEAD 里那 16 枚受追踪的 `design/**` 此刻**已被前端会话从磁盘删除、但尚未 commit**。
票面给的那个假设（"gap 是 owner 未追踪但没被忽略的 `design/old/`、`design/doubao/`"）
因此比我原先写的**更成立**：32 这个数字里没有一枚受追踪件。
`!!` 行仍为 0 ⇒ 结论不变：**ignore 过滤器既吞不到、也不该吞**；§9 那一格的判词（成立/不算缺陷）**不改**。
我 §6.5 正文那句加法是 prose 错（我把 `comm -13` 的 32 读成了"16+16"），
而同一个代码块里那行 `= 16 design/doubao + 16 design/old` 其实已经把真相贴出来了——现场自相矛盾，以本条为准。

### 11.3 追加发现 F9（不属于本修的账，但和本批的台账数字在同一小时内相撞）

现量：`internal/panel` 之外没人碰的这 16 枚 ` D` 一旦被 commit，`design/**` 在 HEAD 里就**一枚文本件都不剩**
（`design/doubao/`、`design/old/` 是 `??`，不跟着进去）。此时 CI 的检出树里 `design/` 目录整个不存在。
仓库外用二进制实测这一形状（`nodesign\` = `git archive HEAD -- . ':(exclude)design'`，**不删任何东西**）：

```
$ ./bin/d22scan-acc.exe -root nodesign          rc=2
d22scan: scope ban #8 design/           examined   0 text files
d22scan: ban #8 scope design/ examined 0 files - it is declared in emojiScopes() but walks nothing.
         Point it at a real tree or delete the entry; never leave a scope pretending to scan (ticket 71 AC#4)
```

⇒ **那一枚 commit 落下去的同一刻，`scripts/d22scan.sh` 会硬红在 step 2**（`ban #8 design/` 是 `live:true` scope，
空走即致命，`main.go:402`），而 `A211` 刚重标进去的 `#8 design/=16` 同时作废。
**判：这不是本修的缺陷，是守卫在干它该干的事**（`TestExemptScopeCannotOutliveItsAbsentTree` /
`TestDeclaredEmojiScopeCannotWalkZeroFiles` 正是为此而在）——但它和本批的台账重标撞在同一小时，
而且它是 **A207 那一族病的第二个实例**："印出来的数踩在*这台机器/这一分钟*的树是什么形状"上。
处置归属：**前端会话 + 编排者**（要么把那 16 枚新 mockup `git add` 进 `design/`，
要么在同一批里把 `ban #8 design/` 改 `live:false` 并登记 exempt 理由——`driftedAbsentScope` 那套机制就是为这个准备的）。
本件不动手，只点名。

### 11.4 总裁补一格（F9）

| 格 | 判 | 一理由 |
|---|---|---|
| §11.3 `design/` 空 scope 即将致命（F9） | **成立**（不算本修的账） | 守卫按设计变红、现量 rc=2 且报错点名 scope；责任在前端会话那 16 枚未提交删除与同批的 scope 处置。 |
