# 96 — `ban #8`（零 emoji）**根本没扫 `frontend/`**：面板那 40 个文本文件今天是"门开着但没人看"

**Status:** open（2026-09-21 15:5x 编排者建；来源=票 77 接续代理交件时**点名为它 AC#4 的硬缺口**，不是它要偷工）
**Type:** 门禁完整性（票 71 AC#4 / A44① / A54② / 票 88 的同族：**覆盖面自己会烂，而输出长得和"检查过"一模一样**）
**Blocks:** 票 **77 的 AC#4**（那一框字面就要"`ban #8` 对 `frontend/` 的自报文件数"，今天**打不出这个数**）
· **Blocked by:** nothing（`tools/d22scan/` 此刻无人写；票 88 已交件）
**Packages:** 只 `tools/d22scan/main.go`（`emojiScopes()` 那一处 + 必要的 walk 过滤器）与 `tools/d22scan/scan_test.go`。
              **禁改**：`allowlist.txt`（**5 行非注释，只许变短或持平**）、任何 ban 的**文本/正则/字符类**（D22：不是代理能缩的）、
              "live 作用域 examine 0 个文件 ⇒ 致命"这条规则本身、`frontend/**`（票 77 的地界）、`internal/**`、`cmd/**`。

## 实测（编排者自己读到源码，不是引用别人的话）

`tools/d22scan/main.go:451-457` 的 `emojiScopes()` 逐字是三条：

```go
{dir: filepath.Join(root, "design"),   label: "design/"},
{dir: filepath.Join(root, "internal"), label: "internal/", goOnly: true},
{dir: filepath.Join(root, "cmd"),      label: "cmd/",      goOnly: true},
```

而 `ban8Scopes()`（`:439-451`）就是把它原样搬进台账 ⇒ CI 自报的三行是
`ban #8 design/ 16 text files`、`ban #8 internal/ 324 Go files`、`ban #8 cmd/ 25 Go files`。
**`frontend/` 不在清单上** ⇒ 票 88 刚武装的 `ban #6` 覆盖 `frontend/` 的 40 个文本文件（票 77d 的实测数），
而**零 emoji 那条门对这些文件一个字节都没读过**。
⚠ 这不是"CI 没跑"：它跑了、绿了、自报了三个数，**只是那三个数里没有我们新长出来的那棵树**——
和票 67 当年把 `ban #6` 登记成豁免是同一族，只不过这次的形状是"**清单忘了同步**"，因此**连豁免都没写过**。

## 为什么这条值得单独一张票而不是"顺手加一行"

加一行 `frontend/` 会**当场把台账用例弄红**（票 88 刚演示过一遍：翻牌代价是 5 条重新表述）。
更要紧的是**判据本身**：`frontend/` 里最可能出现 emoji 的地方是**注释与 UI 文案**（`.tsx`/`.css`/`.md`），
而 `design/` 那条走的是"text files"模式（不要求 `.go`）⇒ **要用哪条路、`node_modules/` 与 `testdata/` 怎么排、
`VENDORED.md` 这种"上游原文里可能就有 emoji"的 vendored 文件怎么办**——这三个决定都会改变 CI 的颜色，
必须**写在 commit 正文里**，不许静默选一个。⚠ 尤其**不许**为了变绿把 vendored 文件整目录排除掉：
那是把 R18 的"三库都是 copy-paste 进仓"这条暴露面重新藏起来。

## AC（1:1，裁决表 `docs/evidence/s1/96-*.md` 由验收方出，不是自裁）

- [ ] **AC#1** 台账里出现第四行 **`ban #8 frontend/ examined N text files`，N>0**，且**逐文件后缀构成**要贴出来
      （多少 `.tsx`/`.ts`/`.css`/`.md`/`.json`/无后缀）。**判据是"N 与票 88 量的 `ban #6` 那个数对得上或解释得清差在哪"**
      ——两条门扫同一棵树，数不同就必须给原因（排除规则不同？），**不许含糊**。
- [ ] **AC#2** 阳性对照：在**仓外纯净快照**（`git archive <SHA> | tar -x -C /tmp/<带你的会话后缀>`）的
      `frontend/` 里种一个 emoji（**注释里放一个**，别只测字符串字面量），`sh scripts/d22scan.sh` 必须 rc=1 并点名它；
      ⚠ 再种一个只在 `.md` 里的 ⇒ 明写你的选择（`.md` 扫不扫）与理由，别让它靠运气命中。
- [ ] **AC#3** 反向对照：**把 `frontend/` 从 `emojiScopes()` 里删掉**（模拟"清单又忘了同步"）⇒
      必须有用例红。若全绿，说明台账根本没在核对清单，那 AC#1 就是自证。
      锚点=承载行为那一行，**同链 grep 证落地再跑**，还原后 `git diff --quiet` 证干净，**编译失败不算变异**。
- [ ] **AC#4** 因加作用域而红的台账用例**逐条重新表述**（票 88 刚做过一遍，照它的做法，**不是删用例**）：
      逐条列"旧断言 / 新断言 / 为什么新断言仍然在测同一件事"。
- [ ] **AC#5** `ban` 文本零改动证明：`git diff` 到本票第一枚 commit 的父，贴出**ban 文本与字符类那几行未动**的证据；
      `allowlist.txt` 仍是 **5 行非注释**（`grep -v '^#' allowlist.txt | grep -c .`）。
- [ ] **AC#6** 门禁：`cd tools/d22scan && go test -count=1 ./...` rc=0（**必须 `-count=1`**，
      ⚠ 它是**独立 Go module**，从仓根跑会打印 `main module does not contain package` 而**扫描器根本没执行**——本仓栽过）、
      `gofmt -l tools/d22scan/` 空、纯净树 `sh scripts/d22scan.sh` **rc=0** 且逐作用域贴全。

## Rules（本仓固定）

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc）；提交前 `git diff --cached --name-only`；
禁 `git add -A`/`.`、`--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34，共树）；**不 push**；
不在仓内建 worktree（A38④）；票面 append-only（改行前先读；标题前插段落要把标题重抄，`git diff --numstat` 删除列必须 0）；
四种假绿逐条点名；数字不达标写 FAIL 附数字。**收尾前必跑 `sh scripts/d22scan.sh`**（A64② 立的规矩）。
15 次工具调用内交回第一枚 checkpoint；每次提交同步 Status + 勾框 + `next=`；接近轮数上限主动收尾留断点。
⚠ **共树提醒**：`internal/winsec/`（票 94）、`internal/config/`+`internal/agent/`（票 90）有人在写，别碰；
`.scratch/wisp/issues/{82,86,87,88,89,90,91,92,93,94,95}-*.md` 是别人的票面。

## Progress log（append-only）

- 2026-09-21 15:5x（编排者）：建票。来源：票 77 的接续代理（`agent-ticket77d`）在交件时写明
  **它的 AC#4 只能记 PARTIAL**——"框字面要 `ban #8` 对 `frontend/` 的自报数，而它的声明作用域不含 `frontend/`，
  改 `tools/d22scan` 非本票地界"。**这个判断是对的，我按它说的单独立案**，不去把 77 的框圆过去。
  我自己读了 `emojiScopes()`（file:line 在上面）确认缺口是真的，**没有采信它的叙述**。
  ⚠ 与票 88 的区别要说清：那次是"目录不存在所以登记豁免"（有账、会漂移报错）；
  这次是**清单根本没同步过**，所以**连豁免都不存在** ⇒ 它比豁免更隐蔽：CI 自报三行、一行都不提 `frontend/`。
  next= 可立即派单（只碰 `tools/d22scan/`，与在飞的票 89/90/94 无文件交集）。
