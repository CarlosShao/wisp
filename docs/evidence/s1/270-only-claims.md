# 票 270 残余 · `270-r2`：两枚互顶的 "only" 改成说实话（只动整行注释）

腿＝`270-r2`（落地腿·注释面一枚）｜日期 2026-10-16｜起手锚 `2e706255`｜写面＝派单那四处
本件是**实现者自证**，⛔ 裁决；对抗验收归非实现者腿。

## 1. 清的是什么

同一块注释里两枚 "only" 互相把对方顶成假话（行号＝起手读数，⛔ 当判据）：

- 旧 `subagent_197.go:257`（护 `:265` 那支）：`it is also the only site that keeps the
  t.d.BaseOptions() call below from being reached with an empty func value`；
- 旧 `:295`（护 `:297` 那支）：`this line is the only thing left that cannot be called empty`。

甲形（台账 `A990` 裁定 3）落地后今天盘上是**两支**，所以两枚 "only" 都是过期文字。
按编排者裁定（⛔ 重新发明）改成两支关系：

| 位置 | 改后说什么 |
|---|---|
| `:265` 上方那块 | 它是 **FIRST of the two sites**，并**自限**："only the first one while its left half still names BaseOptions"；`:297` 是第二支 |
| `:297` 上方那块 | 它是 **SECOND／后备**：`that guard carries a BaseOptions half today, so it answers first in every wiring this package builds today`；"唯一"用条件句 —— `Only in the world where that || is split or loosened` ＋ `today it is one of two` |

两枚引用逐枚真件（⛔ 指名不存在的测试）：`internal/tools/failclosed_236_teeth_test.go` 把
panic 读成自己那枚红（`:219`/`:220` 原文）；`internal/tools/nilbase_270_test.go` 头部原话就
写着上游 guard 先响（`:8`–`:15`）。AC#5 照旧：只写今天谁先响／谁兜底／为什么，⛔ 许诺"将来也
永远拦得住"。

## 2. 产码字节零改动（本腿最重要那把尺）

```
$ git show 2e706255:internal/tools/subagent_197.go | grep -vE '^[[:space:]]*//' > /tmp/270r2/a.nocomment
$ grep -vE '^[[:space:]]*//' internal/tools/subagent_197.go > /tmp/270r2/b.nocomment
$ diff -u /tmp/270r2/a.nocomment /tmp/270r2/b.nocomment; echo "rc=$?"
11291 /tmp/270r2/a.nocomment
11291 /tmp/270r2/b.nocomment
rc=0
```

正控（先证这把尺有牙）＝复制 `a.nocomment`、只在**副本的非注释行**把 `MaxSubagentDepth`
改成 `MaxSubagentDepthZZ` ⇒ 同一把 `diff` **非空 1687 字节、rc=1**。两种形状都判对 ⇒ 那
个"空"才算证据。逐行核：改动 31 枚、其中非整行注释 **0** 枚。分母：起手 612／258／354 ⇒
终态 617／263／354（非注释一字未动）。

## 3. 颜色与门禁

`go test ./internal/tools/ -count=1 -v`（尺口径：`^=== RUN` 含子测试 ⇄ `^--- PASS`/`FAIL`/`SKIP`
只数顶层；台面＝`git archive` 的**仓外**导出树，⛔ 共享工作树）：

| 发 | rc | RUN | PASS_TOP | FAIL_TOP | SKIP | ^panic |
|---|---|---|---|---|---|---|
| 改前 `2e706255` | 0 | 325 | 242 | 0 | 0 | 0 |
| 改后 `087936bc` | 0 | 325 | 242 | 0 | 0 | 0 |

终态红名册与起手名册作差＝**0 枚**（`diff` rc=0，两侧皆空）。
门禁四数（只在工作树上跑，⛔ 部分导出树）：`gofumpt.exe -l` **rc=0**（真身
`/d/work/base/gopath/bin/gofumpt.exe`）／`sh scripts/d22scan.sh` **rc=0**（clean，ban #7
internal/tools/=25、ban #8 internal/=577 含注释）／`sh scripts/check-path-length-budget.sh`
**rc=0**（VERDICT GREEN 60==60）／`go vet ./internal/tools/` **rc=0**。红句零枚含
`walk`/`build failed`/`no test files`/`setup failed`/`0xc0000135`。

## 4. 必答题：摘掉 `:265` 的 `BaseOptions == nil` 半边后，"谁先响"还站得住吗

台件＝`mutate_270r2.py`＋`run-mut.sh`，物理字节改**仓外导出树**（⛔ 碰工作树），锚文本非
恰 1 枚命中即拒绝突变。两发：

| op | rc | RUN | PASS_TOP | FAIL_TOP | ^panic | 那三枚点名腿 |
|---|---|---|---|---|---|---|
| `half-off`（只摘半边） | **0** | 325 | 242 | 0 | 0 | 两枚 nil-BaseOptions 腿**全绿**、`ParentTools` 腿绿 |
| `half-off-no-jia`（再摘 `:297` 后备） | 1 | 325 | 239 | **3** | 0 | 两枚 nil-BaseOptions 腿＋`ThreeRefusalSentences` **红**，got 以 `panic（guard 摘掉后不是拒绝而是崩）` 开头；`ParentTools` 腿仍绿 |

**答**：**半站得住 —— 而且注释已经把这个"半"写明了。**
1. 对 **`ParentTools` 半边**："先响"照旧成立（两发里 `Test236R2…ParentTools…` 都绿＝那半支没被摘、
   仍先答）。
2. 对 **`BaseOptions` 半边**："先响"**不再成立** —— 摘掉后 `:265` 根本不再看这枚字段，第一发
   0 枚红只可能是 `:297` 接住的（第二发把 `:297` 也摘掉 ⇒ 同两枚腿立刻 panic 红，反向钉死
   "答话的就是后备"）。
3. 因此本腿把 `:265` 那块写成**自限句**（"only the first one **while its left half still names
   BaseOptions**"），把 `:297` 那块的"唯一"关进**条件句**（"Only in the world where that || is
   split or loosened"）。摘半边正是走进那个世界 ⇒ 新文字与行为一致，⛔ 需要第二枚腿再补。
4. 这一发只用来验注释文字与行为一致，⛔ 报成"我装了牙"：摘半边对测试名册**零敏感度**（0 枚红）
   是甲形固有弱牙，`270-v1` 判语②与编排者"不另立票"裁定在前，本腿⛔ 翻案、⛔ 为它动产码。
5. 还原＝两发逐发 `MATCH=yes`，`restored_tree_sha256=86de44c4b79b35172c5ca4bd42862f3659e1d9268ab992f975c282eaa6b172e3`
   ＝`git cat-file -p 2e706255:internal/tools/subagent_197.go` 同串；工作树该路径
   `git status --porcelain` 终态 **0 行**（本腿从未拿它当突变台面）。

## 5. 具名上报（三笔，归编排者处置）

1. **本腿自己的缺陷（同一枚错连犯两次，⛔ 抹、⛔ 改历史）**：`-F` 两次都指错了文件 ——
   `9233fc19`（台件与读数九枚）与 `0176ee19`（本件＋`commit-msg-2.txt`）的 message 都是
   `commit-msg-1.txt` 那段（＝与产码提交 `087936bc` 重复）。**intended 正文**逐枚在盘：
   `commit-msg-2.txt`＝台件那一发、`commit-msg-3.txt`＝本件那一发、`commit-msg-4.txt`＝本处
   更正那一发、`commit-msg-5.txt`＝票尾那一发。按"已推送历史⛔ 改写、要更正就追加"的规矩
   ⛔ `--amend`／⛔ `reset`／⛔ `rebase`，只追加更正提交；这一处缺陷的成因是本腿复用同一条
   命令行、只改了 pathspec 没改 `-F` 参数 ⇒ 已把"提交前逐枚对 `-F` 与 pathspec"记进票尾
   bullet。
2. **台面配方还缺一枚**：更正①的八件套（补 `models/`）**仍不够**——本包 `paths_shortname_252_*`
   那 8 枚腿要仓根 `docs/`，且导出树所在路径**必须有 8.3 短名别名**。读数：`/tmp`（C 卷、名已合
   8.3）那发 `FAIL_TOP=11`（8 枚 252 说"短名回显与长名同串……本探针不许 skip"＋3 枚 172 说
   `no docs/PLAN.md above C:Users`）；补 `docs` ⇒ 8 枚；搬到 `D:/work/workspace/wisp-270r2-comment-bench`
   （目录名 >8 字符 ⇒ 有别名）⇒ **0 枚、rc=0**。建议配方＝八件套＋`docs`＋手工三件＋**容器树放
   `D:` 卷且目录名要产生短名**（与记忆里"容器树必须放 D: 卷"同源、那是另一枚成因）。
3. **分母分歧**：派单记"行尾注释形状 **9 行**"，本腿四把候选尺一律复现不出（尺A
   `grep -vE '^[[:space:]]*//' | grep -c '//'`=0；尺B `//.*//`=0；尺C 含 `//` 总行数 258==整行
   注释 258；尺D `/*`=0）。后果只有一句：本文件今天没有可漏的行尾注释，所以第 2 节那把尺没有
   "被同行注释骗过"这一形；9 那一发⛔ 当判据、请编排者复核其尺口径。

## 6. 没做完／不适用（逐枚具名）

- ⛔ 翻任何 AC 框、⛔ 改票名／`Status:`／加 `-done`（六格归 `270-v1` 判语、结残余归编排者按
  `:74`/`:76`）；本件只追加一枚 bullet 到票尾。
- ⛔ 动产码一字（第 2 节那把尺＝证明）；⛔ 拆 `:265` 成两支（乙形代价归编排者裁定）；⛔ 碰
  `ParentTools` 半边／`internal/tools/**` 测试件／`cmd/wisp`／`internal/panel`／`internal/agent`／
  `internal/memory`／`internal/session`／`frontend`／`design`／`PLAN.md`／`docs/specs`／
  `thresholds.go`／golden／`allowlist.txt`／`go.mod`／`go.sum`／`.gitignore`／`ci.yml`／`.github`。
- ⛔ 跑整包 `./cmd/wisp`、⛔ 跑全树 `./...`；⛔ 为变绿放宽断言／删用例／`t.Skip`；⛔ 杀 PID 21892。
- 并发：本腿只读别人的在飞字节都⛔ 当凭据；`go test ./internal/tools/` 全程未与别腿撞车（本波
  `internal/tools` 零写腿），所以**没有等过**；起手锚现跑 `git merge-base --is-ancestor 44f51fc8 HEAD`
  ＝真，`44f51fc8` 在本腿 HEAD 的祖先链上。
- ⛔ 未 push（推送只由编排者做、只推 `origin/dev`）。
