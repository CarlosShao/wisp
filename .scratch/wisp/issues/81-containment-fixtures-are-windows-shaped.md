# 81 — 两份 containment 测试的**阳性对照是 Windows 形状**，ubuntu 上必红（票 76/79 的判据仪器在 Linux 上不自证）

**Status:** in progress — AC#1 已落地（Windows 本机两侧绿），AC#2 的 Linux 侧/docker 实测、AC#3 全仓扫描、AC#4 门禁待交
**Type:** 测试夹具的平台可移植性缺陷（不是生产洞，但**它让"防逃逸的那条判据"在 Linux 上等于没跑**）
**Blocks:** 票 70 的 AC#2/AC#6（`test-core` 现在只剩 10 条红，其中 2 条就是本票）· **Blocked by:** nothing
**Packages:** `internal/agent/spill_path_invariant_test.go`、`internal/memory/artifacts_path_invariant_test.go`
（必要时它们的兄弟测试文件）。**禁改**：`internal/agent/spill.go`、`internal/memory/artifacts.go` 的**生产行为**、
`internal/risk/**`（票 70/82 的地界）、`tools/d22scan/**`、`allowlist.txt`、任何 `docs/PLAN.md`/`docs/specs/*.md`。

## 实测证据（CI run **35558750456**，commit `17efc2c`，job `test-core` 步骤 `Portable package tests`，ubuntu）

**① `TestArtifactsContainmentByDirectoryListing`（`artifacts_path_invariant_test.go:464`）**
```
purge removed [data/artifacts/nested\canary-separator.txt data/artifacts/real-1.txt data/artifacts/real-2.txt],
want exactly [data/artifacts/nested data/artifacts/nested/canary-separator.txt
              data/artifacts/real-1.txt data/artifacts/real-2.txt]
```
根因：夹具用**字面反斜杠**造"子目录里一个文件"。Windows 上 `\` 是分隔符 ⇒ 真的是 `nested/` 目录；
Linux 上 `\` 是合法文件名字符 ⇒ 落成一个**顶层文件**，目录集合自然对不上。

**② `TestSpillContainmentByDirectoryListing`（`spill_path_invariant_test.go:280`、`:286`）**
```
control: escape id "..\\..\\..\\..\\CONTROL-escape" resolves to ".../data/artifacts/tool-output-..\\..\\..\\..\\CONTROL-escape.txt"
the listing never saw the deliberate escape outside the data dir (…); AC#2 would be vacuous
```
根因：**阳性对照**要证明"如果真有文件写到 data dir 之外，列举一定看得见"。Linux 上那串 `..\` 不构成逃逸，
文件老老实实待在 `artifacts/` 里 ⇒ 对照本身失效。**注意这不代表生产有洞**：票 79 的编码器已经把所有
`/ \ : .` 百分号转义，Linux 那个长文件名恰恰是**正确**行为。

## 判据（1:1，裁决表 `docs/evidence/s1/81-*.md`）

- [ ] **AC#1** 两条用例都改成**平台无关地表达同一个意思**：逃逸/嵌套的"物理路径"用 `filepath.Join` 或
  `filepath.Separator` 构造，同时**保留**字面反斜杠拼写作为额外一例（它在 Linux 上是合法文件名、
  在 Windows 上是分隔符，这个不对称本身就是要钉的东西）。**不许**给这两份文件加 `//go:build windows`
  把它们变成"Linux 上静默不跑"——那是 A49② 里票 78 代理拒绝过的那类修法（整包被排除 = 静默跳过）。
- [ ] **AC#2** 两个方向各自**制造一次红**：(i) 在 Linux（docker `golang:1.27`）证明改后这两条**真跑且绿**，
  且把生产侧的转义临时拆掉时它们**变红**（说明判据仍咬得住）；(ii) 在 Windows 本机同样跑一遍。
  ⚠ 变异锚点选**真正承载行为的那一行**，同一条 `&&` 链里先 grep 证明落地，还原后 `git diff --quiet` 证干净。
- [ ] **AC#3** **全仓扫同族**：`grep -rn` 出所有在没有 build tag 的测试文件里用字面 `\` 当目录分隔符的地方
  （含 `"a\\b"`、`filepath.ToSlash` 反用、`\r?\n` 之类合法的除外），逐条列进本票 log：
  要么本票一并修掉，要么写清"它在两侧语义相同、不需要修"的理由。**不许只报"扫了没问题"**——要给出命中清单。
- [ ] **AC#4** 门禁（只跑自己碰的包）：`gofmt -l` 空、`go vet` rc=0、`go test -count=2` rc=0，
  并 `GOOS=linux go vet` rc=0；逐跑点名 `--- SKIP`/`--- FAIL` 行数与名字。

## Rules（本仓固定）

15 次工具调用内交第一枚 checkpoint commit；每次 commit 同步票面 Status + 勾框 + 末行 `next=`；
`git commit -q -F - -- <显式路径>` + **带引号** heredoc `<<'MSGEOF'`；禁 `git add -A`/`.`；commit 前核对
`git diff --cached --name-only`；禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`（共树 A34）；**不 push**；
不在仓内建 worktree（A38④），纯净树用 `git archive HEAD | tar -x -C /tmp/...`；
docker 挂载**快照**而不是工作树（别人的在飞改动会混进你的结论）；
票面 Progress log append-only：**要改的那一行先读再替换**。

## Progress log（append-only）

- **C1（AC#1 的代码侧）** 两份夹具改成平台无关表达，字面 `\` 一例保留并**显式钉住不对称**：
  - `internal/memory/artifacts_path_invariant_test.go`：`inv76Shapes` 的 `separator`/`dotdot` 改用
    `inv76Sep()`（`filepath.Separator`）拼，target 一律 `filepath.Join(被指的目录, 调用者给的名字)`
    = "这台机器把它解析成哪儿"；新增 `separator_backslash_literal`、`dotdot_backslash_literal` 两例保留字面 `\`。
    (e) 段原本硬编码 4 条 `data/artifacts/...` 键（Linux 上必然对不上），改为由夹具真实种下的路径经
    `inv76ReclaimKeys` 派生（自己 + 被 OS 认可的父目录，artifacts 根本身不算）；purge 计数也不再硬编码 2，
    改由 `before` 快照里 artifacts 树下的**文件**条数 -1 推得（`PurgeArtifacts` 的 n 只数文件，实测 Windows 6 键/3 文件）。
    新增 `TestArtifactsLiteralBackslashKeepsItsAsymmetry`：同一串字节在 Windows 是两层、在 Linux 是一个扁平文件名，
    两侧各自断言，并注明"加 build tag 就等于删掉这条判据"。
  - `internal/agent/spill_path_invariant_test.go`：`inv76aEscapeID` 改用 `inv76aEscapeIDWithSep(..., string(filepath.Separator))`
    （Linux 上 `..\` 不构成逃逸 ⇒ 阳性对照失效，这就是票面 ② 的根因）；字面 `\` 拼写保留为**第二条对照** (d)，
    Windows 断言它落在 data dir 之外、Linux 断言它落回 artifacts 里成为一个扁平名；端到端用例的 `nested\canary.txt` /
    `..\canary.txt` 扩成"平台分隔符 + 字面反斜杠"双拼写。
  - Windows 本机：`gofmt -l` 空、`go vet ./internal/agent ./internal/memory` 无输出、
    `-run 'Containment|LiteralBackslash|HostileShapes|StaysUnderDataDir|RejectsTheFour'` 两包 `ok`（0 SKIP / 0 FAIL）。
  - AC#1 勾框**留到 Linux 侧绿了再打**——票面要的是两侧都跑，不在只测了一侧时就记完成。
  next= docker `golang:1.27` 跑 `git archive` 快照，证明这两条在 Linux 上真跑且绿（AC#2(i)），然后做生产侧转义/守卫的变异检验
- **C2（Linux 首轮实测：改对了方向，但暴露我第一版仍带 Windows 预设）** 快照 `git archive HEAD` → `D:\tmp\wisp81snap`，
  `docker run golang:1.27 -count=2 -v` 两包：50 条 `=== RUN`、**4 条 FAIL**（同一对，×2）、0 SKIP；agent 包 `ok`，memory 包红。
  失败原文（`artifacts_path_invariant_test.go:502`）：
  `purge removed [data/artifacts/..\canary-dotdot-literal.txt data/artifacts/nested data/artifacts/nested-backslash\canary-literal.txt
  data/artifacts/nested/canary-separator.txt data/artifacts/real-1.txt data/artifacts/real-2.txt] want exactly [去掉第一条]`。
  根因是**我自己的 C1 版本**：`wantRemoved` 仍按名字点了两例字面反斜杠 canary，而 Linux 上 `..\canary-dotdot-literal.txt`
  这个字面名**不越界、老老实实落在 artifacts 里**，于是它也归 purge 管——这正是票面要的不对称的第二面，
  第一版把它漏掉了（同一原因让 `LiteralBackslash` 那条断言误报"意外反斜杠条目"）。
  修法：期望集合改为**按前缀从 target 派生**（`inv76Rel(root, sh.target)` 落在 `data/artifacts/` 里就算），
  不对称测试 `TestArtifactsLiteralBackslashKeepsItsAsymmetry` 重写成两侧通用的"含 `\` 而本平台不当它是分隔符 ⇒ 必是扁平一条"，
  并加 `len(planted) < 4` / `slices.Equal(flats, want)` 两条反-vacuous 卫兵。
  Windows 本机重跑：`rc=0`、50 条 `=== RUN`、0 FAIL、0 SKIP。
  next= 重新 archive 快照跑 Linux（AC#2(i) 要绿），再做 AC#2 的变异检验
