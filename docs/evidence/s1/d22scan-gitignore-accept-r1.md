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
3. **`-root` 指到不是一枚 wisp 仓的地方**：响亮且 `rc=2`。我的迷你件（只有 `.gitignore` + `frontend/`，
   没有 `internal/`、`cmd/`、`go.mod`、`allowlist.txt`）实测被 `checkRoot` 拒绝：
   `go run . -root forced/ci` → 无任何 scope 行、`rc=2`（`main.go:1124-1129` 的 required 列表）。
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
