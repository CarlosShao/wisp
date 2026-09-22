# 票 118 独立对抗验收表（acceptor-ticket118）

会话：`acceptor-ticket118`（2026-09-22）。工作树 `dev`，被验区间 `d1056c0..7321c76`（两任实现方九格），
复核对象另含 `33c8acd`/`13c972e`/`034080c`（`agent-ticket119b` 在同一条带上交的件）。

**接续记录**：上一枚会话写到 17:42（本文件 9183 字节处，停在"读数一"的 AC#7 之后）被连接中断掐死，
`acceptor-ticket118b` 于 19:2x 接手。**"读数一"那一节（AC#1/AC#2/AC#3/AC#7）是前任验收方亲手打的变异，本方一字不改、
不重写**，只在"读数五"里对它做抽验（快照 md5、归档日志逐发重点数、变异字节落地）并补 AC#3 的恒真/恒假两发；
AC#4/AC#5/AC#6/AC#8/AC#9 与文末收表、`R-118-*` 登记由本方独立打完（"读数二/三/四"）。
本方另用新快照 `/d/tmp/ac118b-c1`（`git archive 7321c76`、`git archive 33c8acd`），与前任的 `/d/tmp/ac118r` 零共用，
POSIX 腿自己起容器重跑，未沿用前任任何一份日志当作自己的读数。

**本表全部读数出自仓外纯净快照**：`git archive <sha> | tar -x -C /d/tmp/ac118r/<名>`，仓库树内未建 worktree、未 checkout、
未对任何生产码/测试码写一行（唯一写件＝本文件）。快照基线两枚：`base_7321c76`（118 交件末态）、`base_fbbecaa`（当前 HEAD，
AC#9 的争议发要用）。

## 标签口径

- 〔独立复现〕＝本验收方自己在纯净快照里造变异、先证落地再读红名；
- 〔日志＋归档，我抽验〕＝读了实现方日志/commit 与归档文件，抽验了其中若干数字但未亲手重造那一发；
- 〔仅自述，不背书〕＝无法在本机/本轮内复现，或实现方自己就没量。

## 裁决表（骨架，逐格落盘）

| 格 | 判据（票面一句话） | 裁决 | 标签 |
|---|---|---|---|
| AC#1 | `kind=` 三取值各一枚；删 switch 必须红 | 通过（附 `R-118-1` 口径） | 〔独立复现〕＋本方抽验 |
| AC#2 | 两桶互换 ≥2 枚红（归属面 + 渲染面） | 通过（实测 9 枚，`R-118-3`/`R-118-4`） | 〔独立复现〕＋抽验归档 |
| AC#3 | `namesEveryone` 不按 2 字节子串 + 反向对照 | 通过（本方补恒真/恒假两发） | 〔独立复现〕 |
| AC#4 | 生产码一行不动（diff 证明） | 通过 | 〔独立复现〕（读数三） |
| AC#5 | 门禁四数点名 + gofmt/gofumpt/vet 双 GOOS + d22scan | 通过 | 〔独立复现〕（读数三） |
| AC#6 | POSIX 叶子方向两枚；MUT-B 从红 1 → ≥3 | 通过（MUT-B 红 3；本方补 MUT-LEAF 红 6） | 〔独立复现〕（读数二） |
| AC#7 | 两枚装饰腿在恒真/恒假下发红名 | 通过（`:294/:401`、`:299/:404` 逐行对上） | 〔独立复现〕＋抽验归档 |
| AC#8 | 跨卷归属：先量再判（票面留 `[ ]`） | **通过（停手移交票 126）**；缝守腿不背书 | 链＋判据＋跨卷读数〔独立复现〕／缝守腿〔仅自述，不背书〕 |
| AC#9 | envfork 那枚断言守语义不守拼写 | 通过；119b 争议判为"两枚 vintage 各得其所" | 〔独立复现〕（读数二） |

## 读数一：AC#1 / AC#2 / AC#3 / AC#7（windows 腿，〔独立复现〕）

快照：`git archive 7321c76 | tar -x -C /d/tmp/ac118r/base_7321c76`。基线（我自己的树，未变异）：
`go test -count=1 -v ./internal/winsec/` rc=0，**`=== RUN`=86 / `--- PASS`=46 / `--- FAIL`=0 / `--- SKIP`=0**
（`ok … 16.468s`）——与 118b 自述的基线逐数一致。

每一发变异都先在快照里 grep 出改后的字节（下附落地证据）、再 `go build ./internal/winsec/` rc=0、
`go vet ./internal/winsec/` rc=0，才读红名。九发全部 build/vet 双 0。

| 发 | 改的字节（快照内 grep 原文） | 我量到的红 | 实现方自述 |
|---|---|---|---|
| MUT-M5 | `kind := "explicit"` + `_ = n // MUT-M5: the whole kind switch is gone`，`grep -c 'case len(n.Principals)'`=0 | **3**：`TestAC1DefaultLogSaysInherited`、`TestAC118KindFieldSaysInheritedForAGrantTheParentHandedDown`、`TestAC118KindFieldSaysBothWhenTheObjectCarriedItsOwnGrantToo` | 红 3 枚，同一对名字 ✔ |
| MUT-M5-INV（**我补的**） | switch 删掉且常数换成 `"inherited"` | **2**：`...SaysExplicit...`、`...SaysBoth...`（`TestAC1DefaultLogSaysInherited` 反而绿） | 未做 |
| MUT-M5-BOTH（**我补的**） | switch 删掉且常数换成 `"explicit+inherited"` | **3**：`TestAC1DefaultLogSaysInherited`、`...SaysExplicit...`、`...SaysInherited...` | 未做 |
| MUT-M4 | `inherited = append(...)` ↔ `explicit = append(...)` 两处互换（`MUT-M4 swapped` 计数 2） | **9**（见下） | 自述"红 5 枚"，但同一行里它自己列了 6 个名字 |
| MUT-M4B（**我补的**） | 生产判定不动，只把渲染的两枚 `strings.Join` 目标互换（kind 仍算对） | **3**：恰是三枚 `TestAC118KindField*`；`TestAC1SealFileReports...`、`TestSealReportsThePrincipalsItCleared` 两枚归属面用例**仍绿** | 未做 |
| MUT-NOINH（**我补的**） | 删掉 `if ace.inherited {...continue}` ⇒ 不看继承位 | **4**：归属面 1 + 渲染面 1 + 三枚 kind 中的 2 | 未做 |
| MUT-TRUE | `_, _ = n.Path, resolved.String()` + `return true // MUT-true` | **4**，含 `TestAC2InheritedNoticeHasANoiseBound`（`inherited_narrow_notice_104_windows_test.go:294`）与 `TestAC3OwnGrantsStaySilentWhicheverWayTheOSNamesThem`（`:401`） | 自述同一对行号 ✔ |
| MUT-FALSE | 同上 `return false` | **9**，含同一对用例的 `:299` 与 `:404` | 自述同一对行号 ✔ |
| MUT-WD | `namesEveryone` 退回 `Contains("WD")\|\|Contains("S-1-1-0")\|\|Contains("Everyone")`（测试码） | **1**：`TestAC118NamesEveryoneIsJudgedBySIDNotByTwoBytes` | 红 1 枚，同名 ✔ |
| MUT-WD2（**我补的**） | 只留 `Contains("WD")` 一枚针 | **1**：同一枚用例，但红在**正向那一腿**（`:560`，`Everyone:(RX)` 的 icacls 渲染里没有 `WD` 两字节） | 未做 |

### AC#1 判定：通过〔独立复现〕

票面判据"删掉 switch 必须让新用例红"成立（MUT-M5 红 2 枚 kind + 票 104 的渲染用例）。
118b 自己如实补的那一格**我单独验了并且成立**：`TestAC118KindFieldSaysExplicitForAGrantThatStoodOnTheObjectItself`
在 M5 下**必然仍绿**——M5 的输出就是 `explicit`，任何"把自己那格删成自己"的变异都红不了它。
它算不算"三取值各一枚"的例外：**不算**。判据应当是"每一枚取值都被至少一发 switch 删除形红过"，
我补的两发把这格补齐了：常数= `inherited` 时红 `SaysExplicit`+`SaysBoth`，常数= `explicit+inherited` 时红
`SaysExplicit`+`SaysInherited`，常数= `explicit`（原 M5）时红 `SaysInherited`+`SaysBoth`。三枚用例各红过、
且各自由**不同**的常数代入红，没有一枚是恒绿的装饰。

### AC#2 判定：通过（附读数更正），但实现方那一格**数错了**〔独立复现〕

MUT-M4 我量到 **9 枚**顶层红，不是自述的 5 枚（他们自己列的名字已经是 6 个）：
`TestAC1SealFileReportsTheInheritedGrantItCleared`、`TestAC1DefaultLogSaysInherited`、
`TestSealReportsThePrincipalsItCleared`、`TestNoticeAttributionSurvivesAn83AliasOfItsOwnTree`、
三枚 `TestAC118KindField*`、`TestSealNarrowsAndNamesThePrincipalItRemovedBySID`、
`TestGateRefusesADescriptorThatLeavesARealGrantToAnotherAccount`。
票面判据（≥2 枚、一枚归属一枚渲染）成立，"桶换了而全绿"那一形**造不出来**。

残留形状的追查（实现方没做）：MUT-M4B 单独把**渲染面**两桶互换、生产 kind 计算保持正确 ⇒
仍红 3 枚（三枚 kind 用例），所以"字段没跟着换"不可能全绿；但同一发下**归属面**两枚用例是绿的
——即归属面与渲染面各自只被对方一侧的哨兵守住。这不是本票 AC#2 的判据缺口（AC#2 只要求两侧各≥1 枚红、
且这一发两侧红的是渲染面、M4 那一发两侧都红），登记为 `R-118-3` 的建议。

### AC#3 判定：通过〔独立复现〕，且"误认针"那一发复算成立

MUT-WD 红 1 枚，红名点名它自己，明细行：
```
notice_kind_and_everyone_118_windows_test.go:579: AC#3 RED: namesEveryone admitted a SERVICE grant as Everyone:
  C:\Users\swq\AppData\Local\Temp\TestAC118NamesEveryoneIsJudgedBySIDNotByTwoBytes2294106694\002\svc\service-stands-here.txt NT AUTHORITY\SERVICE:(RX)
```
118b 说旧规则的**真正误认针不是 `WD` 而是 `Everyone` 这个词被 `t.TempDir()` 拼进路径**——复算成立：
上面那串文本里 `Everyone` 出现（来自测试函数名 `TestAC118NamesEveryoneIsJudgedBySIDNotByTwoBytes`），
而 `WD` 两字节**一个都不在**。（`...I sJudgedBy...` 无 W、无 D 相邻。）
行号 `:579` 也对得上：它报的是**变异后**的快照行号——我的同一发也是 `:579`，
仓库树里那句 `t.Fatalf` 在 `:585`，差 6 行＝我在同一处删掉的 7 行换 1 行。**这不是他们的读数错，是变异树口径**，
登记一条口径备忘即可。

我补的 MUT-WD2（只留 `WD` 一枚针）红在同一枚用例的**正向腿** `:560`：
`icacls` 把 Everyone 的主体渲染成 `Everyone:(RX)`，其中没有 `WD` ⇒ 旧规则连"真的 Everyone 站在那儿"都认不出。
两发合起来：旧规则既过（误认 SERVICE）又不及（认不出 icacls 文本里的 Everyone），
新判定按 trustee 位→SID 比，在**同一批 fixture** 上 86 枚全绿（基线 FAIL=0）。

### AC#7 判定：通过〔独立复现〕——两枚曾经的装饰腿现在都有牙齿

恒真发：`inherited_narrow_notice_104_windows_test.go:294`（六枚子文件各报一条 `emitted 1 WARN(s) ... bound is 0`）
与 `:401`（`sealing a child that only ever held our own trustees ... emitted 1 WARN(s) attributable to it, want 0`）；
恒假发：`:299`（`the parent that was sealed holds 0 notice(s) attributable to it, want exactly 1`）
与 `:404`（`the witness file held one genuinely foreign grant and 0 notice(s) are attributable to it`）。
红名点到它们自己、行号与自述逐一对齐，**不需按票面规则降级**。票面 AC#7 那句"造不出红就是装饰"不成立。
（顺带：子体那一腿的红名是 `TestAC2InheritedNoticeHasANoiseBound` 顶层，`-v` 明细里能看见
`parent_policy_change_propagates_without_a_per_child_storm` 的 fixture 路径，与自述一致。）



## 读数二：AC#6 / AC#9（POSIX 容器腿，〔独立复现〕，acceptor-ticket118b 接续）

接续口径：下面这一节由 `acceptor-ticket118b` 在断线后自己重打，**不复用上表任何一发读数**；
快照我自己建：`git archive 7321c76 | tar -x -C /d/tmp/ac118b-c1/base`、
`git archive 33c8acd | tar -x -C /d/tmp/ac118b-c1/t33`（与前任验收方的 `/d/tmp/ac118r` 零共用，AC#4 已核这两枚
`internal/proc/envfork.go` 逐字节相同、`internal/winsec/dataroot_symlink_119_other_test.go` 在 `33c8acd..HEAD` 未再动）。
容器命令原文：

```
MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -e WISP_ENV=test \
  -v /d/tmp/ac118b-c1:/work -v wisp118mod:/go/pkg/mod -v wisp118build:/root/.cache/go-build \
  -w /work golang:1.27 bash /work/posix.sh
```

挂载先证明（假绿那一坑）：容器内 `ls -l /work/base/go.mod` = **883 bytes**、`ls /work/base/internal/winsec | wc -l` = **26**、
`go version` = `go1.27.1 linux/amd64`。每一发变异先 grep 出改后的字节（`posix_out.txt` 第 9-17 行＝落地原文），
再 `go build` / `go vet`，才读 test——**参与测量的 14 发 build=vet=0 全过**，无一枚是靠编译错误"红"的。

### AC#6 判定：通过〔独立复现〕，MUT-B 的红数从 1 涨到 3 成立

| 发 | 我打进去的字节 | 容器读数（`-count=1 -v ./internal/winsec/`） |
|---|---|---|
| 基线 `base` | 无 | rc=0，**RUN=42 PASS=27 FAIL=0 SKIP=0** |
| **MUT-B**（`winsec_other.go:129` 走链改成 `pieces[:len(pieces)-1]`＝不查叶子） | `pieces := pathPieces(path) /* MUT-B: the leaf is not looked at */` + `for _, prefix := range pieces[:len(pieces)-1] {` | rc=1，**FAIL=3**：`TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed`、`TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed`、`TestAC1POSIXSealDirThroughASymlinkRefuses` |
| **MUT-LEAF**（**我补的**反向半修：`pieces[len(pieces)-1:]`＝只看叶子、不查祖先） | 同上位置，切片取反 | rc=1，**FAIL=6**：`TestAC1POSIXUnresolvedSymlinkedRootStillRefused119`、`TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`、`TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone`、`TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth`、`TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing`、`TestAC1POSIXSealFileThroughABackslashNamedLink` |

红名是行为的不是拼写的，逐字（我的这一发，行号与 118b 自述逐一对齐）：

```
placement_leaf_118_other_test.go:84: AC#6 RED: SealFile(".../002/root/artifact.txt") left the victim at mode=600 uid=0 gid=0, it was mode=666 uid=0 gid=0 - the mode change followed the leaf
placement_leaf_118_other_test.go:114: AC#6 RED: PrivateFile("...") replaced the victim's bytes ("not this tree's data" -> "top secret"); the open followed the leaf
placement_leaf_118_other_test.go:118: AC#6 RED: PrivateFile("...") left the victim at mode=600 uid=0 gid=0, it was mode=666 uid=0 gid=0
```

票面判据"半修 MUT-B 现在至少红 3 枚"= 正好 3 枚，达成；实现方自述的基线四数（42/27/0/0）与三枚红名**逐字复现**。
我补的 MUT-LEAF 那一发是对称性的加验：叶子腿与祖先腿各自都被至少一枚用例钉住，
"只修一半"在这一族里没有全绿的走法。（登记为 `R-118-5` 的正面项，不是缺口。）

### AC#9 判定：通过〔独立复现〕——那枚断言今天守语义，且不恒真也不恒假

六发（`-run TestLayoutForTestEnv ./internal/proc/`，形状＝容器内 `ln -sfn /tmp/realtree /tmp/linktree`，TMPDIR 取 `/tmp` 或 `/tmp/linktree`）：

| 树 | TMPDIR | rc | 读数 |
|---|---|---|---|
| `base`（现码＋AC#9 新断言） | `/tmp` 实目录 | 0 | 绿 |
| `base` | `/tmp/linktree` 软链 | 0 | 绿，`:139` 记 `test DataDir = "/tmp/realtree/wisp-test-1070" (TMPDIR = "/tmp/linktree", which reaches "/tmp/realtree")` |
| `mut-oldassert`（把 AC#9 那一腿退回 118 之前的字符串比较） | 软链 | **1** | `envfork_test.go:121: test DataDir = "/tmp/realtree/wisp-test-1267", want "/tmp/linktree/wisp-test-1267"` |
| `mut-oldassert`（**我补的对照**） | `/tmp` 实目录 | 0 | 绿——旧断言不是坏在一般形状上，只坏在软链形状 ⇒ 它钉的确实是拼写 |
| `mut-pre119`（生产码 `TestDataDir` 不走 `SealableRoot`） | 软链 | **1** | `envfork_test.go:137: the test data root "/tmp/linktree" still reaches itself through the link at "/tmp/linktree" …` |
| `mut-pre119` | `/tmp` 实目录 | 0 | 绿——同一枚断言在没有链的形状上不制造红 ⇒ 它也不是恒假 |

⇒ 票面 AC#9 那两句（"断该树已被调用方解析的语义而不是某个拼写"、"按 119 的读数以容器真跑复现修前红/修后绿"）**独立复现成立**，
实现方自述的五发读数逐行对上，我补的第 6 发（oldassert 落在实目录上仍绿）把"旧断言只坏在软链形"这一句钉死。
新断言的期望值出自 `filepath.EvalSymlinks(os.TempDir())` + `os.SameFile` + 自己走链的 `firstLinkInPath`，
**不是拿被测函数算期望**——这正是 `R-119-9` 那一族恒真的形状，这里没有。

### AC#9 附：`33c8acd` 那一发争议判完了——**两任读的是同一枚变异、两棵不同的树，都不是假绿**

`agent-ticket119b` 主张恒真那一发现在红名点名 `TestAC2POSIXInjectedTestDataDirStandsAsDeclared119`，
而票 119 的验收方在同一发上读到 56/56 全绿（票 119 验收表 `:343` 的 `119-R`，票面 `:75` 引同一读数）。
我自己打同一枚变异（`envfork.go:104` 的 `return dir` → `return SealableRoot(dir)`）在**两枚 vintage** 上：

| 树 | fixture 形状 | 同一发变异下的 `./internal/winsec/` 读数 |
|---|---|---|
| `base`＝`7321c76`（`33c8acd` 之前） | `injected := filepath.Join(proc.SealableRoot(base), "harness", "picked")` | rc=0，**RUN=42 PASS=27 FAIL=0**——那枚用例**绿**，验收方的全绿读数复现 |
| `t33`＝`33c8acd`（返修后） | leg 1 经链接声明：`declared != asTheKernelSpellsIt`，且有 `t.Fatalf("premise broke: …")` 自证前提 | rc=1，**FAIL=1**，唯一红名正是 `TestAC2POSIXInjectedTestDataDirStandsAsDeclared119`——119b 的红名读数复现 |
| `t33` 未变异 | — | rc=0，RUN=42 PASS=27 FAIL=0 |

机制逐字可查：`SealableRoot` 幂等，`33c8acd` 之前那枚 fixture 自己先把根解析了一遍再声明，
于是"原样返回"与"也解析"给出同一枚字符串，这一发变异在这棵树上**没有任何可观察后果**（用例不是恒真，
是把注入值换成别的仍会红；准确说法是**对它本该守的那一发不敏感**）。
⇒ **判：这一发是 119b 的形状，且它的说法成立；验收方的 56/56 也成立，读的是修前那棵树。两边都不退回。**
分界点就是 `33c8acd`（17:03）与票 119 验收表入库的 `bbcb965`（16:33）谁先谁后，不是读数谁在撒谎。
账**记在票 119 的 `R-119-9`**，与本票 AC#9 无涉：AC#9 那一格动的是 `internal/proc/envfork_test.go`，
`git show --name-only 3b03f00` 只有这一枚文件（我已核，见读数三 AC#4）。
顺带一条口径：同一发变异在 `./internal/proc/`（linktree 形）下两棵树都红 2 枚
（`TestLayoutForTestEnv` 的 `:84` 与 `TestPortableOverride` 的 `:282`，都是既有的注入值字符串比较，
**不是** AC#9 那一腿），plain `/tmp` 形下 0 红——这解释了为什么票 119 验收方那句"没人守这条线"
只在 plain 形成立，登记为口径备忘 `R-118-6`。

## 读数三：AC#4 / AC#5（仓库树内只读命令＋纯净快照门禁，〔独立复现〕）

### AC#4 判定：通过〔独立复现〕——118 那一串 commit 一行生产码都没动；两任的账也分开了

我自己打的命令与读数（全部 `git ... | wc -l` 计数，不靠目测）：

| 命令原文 | 读数 |
|---|---|
| `git diff --stat d1056c0..7321c76 -- internal/winsec/winsec_windows.go internal/winsec/winsec.go internal/winsec/resolve.go internal/proc/envfork.go` | **空输出**；`--numstat` 同 **0 行**，rc=0 |
| 同上再加 `internal/winsec/winsec_other.go` | 仍 **0 行**（票面 AC#4 的禁改名单里也有它） |
| 九枚 118 commit 逐枚 `git show --name-only` | 只出现 5 枚测试文件 + 本票票面：`69c7236`＝winsec 三枚 `_test.go`、`712d048`＝1 枚、`0b1fd06`＝1 枚新 `placement_leaf_118_other_test.go`、`3b03f00`＝`internal/proc/envfork_test.go`、`f9dcd3b`/`db469bc`/`12a6804`/`34d606c`/`7321c76` **只有票面** |
| `git diff --name-status d1056c0..7321c76`（全量 17 路径） | 除上述 5 枚测试文件与本票票面外，另有票 121 的 6 枚 commit（`7fe5e73`…`708221d`）带进来的 `cmd/wisp/main.go`、`cmd/wisp/models.go`、`internal/models/**` 九枚——**那是票 121 的账，不在本票禁改名单上，也不动 AC#4 那四枚文件**（上一行已量：四枚文件 0 行） |

两任分账（不许糊过去的那一笔）：

- `7321c76`（10:45）之后才有 119 返修。`agent-ticket119b` 的 **`034080c`（17:11）** 才是动过
  `internal/proc/envfork.go`（+26/−7）与 `internal/winsec/winsec_other.go`（+41/−14）的那一发，
  我按字节机器验了它是**纯注释**：对 `034080c` 涉及的三枚文件（另含 `cmd/wisp/secret.go` +4/−0）逐行取 `^[+-]`、
  滤掉 `+++`/`---`、再滤掉以 `//`、`*`、`/*` 开头的行与空行，**剩下 0 行**——即判定分支 0 hunks，119b 自述站得住。
- `36294c2`（119b，17:03）动的是 `cmd/wisp/secret.go` 的生产一行 + 一枚新用例，**没碰** `envfork.go`/`winsec_other.go`
  （`git show --name-only` 只有那两枚路径）——票面上"两任都提过 36294c2 动过 envfork.go"的说法要勘正一笔，登记为 `R-118-2`。
- 结论：本票 AC#4 的"一行不改"**成立**，`d1056c0..7321c76` 这个区间口径也成立（119 返修全在它之后）。

### AC#5 判定：通过〔独立复现〕——门禁我在 `git archive 7321c76` 的快照里逐数重跑，与自述逐数相同

快照 `/d/tmp/ac118b-c1/base`（`git archive 7321c76`，树内无 `.git`），命令原文与读数：

| 命令原文 | rc | 读数 |
|---|---|---|
| `go test -count=2 -v ./internal/winsec/` | 0 | `=== RUN`=**172** == 不同名 **46** × 2；顶层 `--- PASS`=**92**、`--- FAIL`=**0**、`--- SKIP`=**0**；`grep -c '(cached)'`=**0** ⇒ `-count=2` 确实不缓存；收尾 `ok … 19.472s` |
| `go test -count=2 ./internal/winsec/`（非 `-v`） | 0 | 全文一行 `ok github.com/CarlosShao/wisp/internal/winsec 19.411s`；`--- PASS` 行数 **0**、含 `SKIP` 行数 **0**、`=== RUN` 行数 **0** |
| `go test -count=2 -v ./internal/proc/` | 0 | `=== RUN`=**72**（36×2）、`--- PASS`=**62**、`--- FAIL`=0、`--- SKIP`=**2**，两枚都是既有 `TestHelperProcess`（CI 的 `-skip` 清单里本就有），非本票新增 |
| `gofmt -l .`（快照根） | 0 | **0 枚候选** |
| `$(go env GOPATH)/bin/gofumpt.exe --version` / `-l .` | 0 / 0 | `v0.7.0 (go1.27.1)`（**存在**，不是"未跑"）；`-l .` **0 枚候选** |
| `go vet ./internal/winsec/ ./internal/proc/` | 0 | 干净（host=windows） |
| `GOOS=linux GOARCH=amd64 go vet …` | 0 | **只编译不执行**，这一格证不到任何 Linux 行为；真正的 Linux 执行读数在上面的容器那一跑 |
| `GOOS=darwin GOARCH=arm64 go vet …`（**我补的**第三 GOOS） | 0 | 同样只编译；`winsec_other.go` 的 `_other` 腿在 darwin 也编得过 |
| `bash scripts/winsec-tests.sh`（CI 同形） | 0 | `four numbers (all from -v output): === RUN=86 --- PASS=46 --- FAIL=0 --- SKIP=0`，`winsec result line: ok … 10.085s` |
| `sh scripts/d22scan.sh` | 0 | `d22scan: clean - no D22 ban violations`；正向对照 `runtests.sh: OK - packages=[./...] PASS=21 FAIL=0 SKIP=0, === RUN=31` ⇒ 门不是瞎的 |

台账八 scope（`7321c76` 这一枚）：`bans #1-5 internal/=202、cmd/=22；#6 frontend/=40；#7 internal/tools/=18；
#8 design/=16、frontend/=40、internal/=382、cmd/=30`——与 118b 自述**逐数相同**。对今天 HEAD（`8ec04f3`）的差是
`cmd/ 30→31`（票 121/127 在 `fbbecaa` 之后各带一枚文件；前任验收方在 `base_fbbecaa` 上量到的就是 31），
**八格无一下降**。AC#5 票面要求的四数点名、`-count=2` 不缓存、非 `-v` 不印 PASS/SKIP、gofumpt 存在性、双 GOOS 口径、
d22scan 纯净快照 rc=0 + 台账不降——**六项全部独立复现**。
一处如实登记的口径分歧（不影响裁决）：118b 的 AC#5 那一跑是在 `34d606c` 上量的，
我在 `7321c76`（交件末态，含它自己那五枚 commit 的票面文字）重跑，四数完全一致，
差别只剩 `cmd/wisp` 那几枚文件不在本票范围内。

## 读数四：AC#8（跨卷归属，责任链与三条判据〔独立复现〕；缝守那一腿**不背书**）

### 责任字节链：我逐字 grep，不抄任何人的行号

`git grep -n "sameTree\|pathComponents\|VolumeName" 7321c76 -- internal/winsec/` 的读数（去掉测试与注释行）：

```
internal/winsec/winsec_windows.go:111:  return sameTree(n.Path, resolved.String())          <- noticeNamesTree 的那一发
internal/winsec/resolve.go:325:         if sameTree(anchorAns, parentAns) || answerInsideTree(anchorAns, parentAns) {
internal/winsec/resolve.go:335:         func sameTree(a, b string) bool {
internal/winsec/resolve.go:336:          compsA, compsB := pathComponents(a), pathComponents(b)
internal/winsec/resolve.go:359/360:     func answerInsideTree(...) { childComps, parentComps := pathComponents(child), pathComponents(parent)
internal/winsec/winsec.go:347:          func pathComponents(path string) []string {
internal/winsec/winsec.go:353:          	i := len(filepath.VolumeName(path))                 <- 卷段从来不进比较，就是这一字节
internal/winsec/winsec.go:323:          	vol := filepath.VolumeName(path)（pathPieces，保留卷，不参与比较）
internal/winsec/resolve.go:514:          i := len(filepath.VolumeName(path))（lexicalTraversal，形状检查，不参与比较）
internal/winsec/placement_windows.go:37/41 卷前缀与 UNC 的落点底线（118b 说 subst/UNC 表达不了这一形，位置就在这两行）
```

⇒ 票 126 票面 `:9` 写的那条链 `winsec_windows.go:111 → resolve.go:335 → winsec.go:353` **逐字成立**（三处行号我自己在
`7321c76` 上验的）。编排者给我的"resolve.go:353?"不是责任字节——那一行是 `sameTree` 的收尾空行，
比较两侧的函数在 `winsec.go`（347/353），`resolve.go` 里只有调用者（336/360）。登记为 `R-118-2` 的同一条勘误。

### 我自己的跨卷探针（快照内文件，不进仓库）

`/d/tmp/ac118b-c1/base/internal/winsec/ac118b_xvol_118b_probe_windows_test.go`（`go vet ./internal/winsec/` rc=0 后才读），
两枚 face，`go test -count=1 -v -run TestAC118B ./internal/winsec/`：**RUN=2、rc=1**（红的是行为面，正是量到了才红）。

```
unit face : pathComponents(A)=[ac118b-xvol store-3131 artifact.txt] pathComponents(B)=[ac118b-xvol store-3131 artifact.txt]
            VolumeName(A)="C:" VolumeName(B)="D:"  =>  sameTree(A,B)=true   answerInsideTree(A, D 的父)=true
            四枚对照全绿：sameTree(A,A)=true、不同叶子=false、不同目录=false、不同深度=false
behavioral: AC#8 shapes: A=C:\ac118b-xvol\store-3131\artifact.txt B=D:\ac118b-xvol\store-3131\artifact.txt (tails identical)
            AC#8 notices after sealing A only: [{Path:C:\...\artifact.txt Principals:[] Inherited:[S-1-1-0(A;ID;0x1200a9;;;WD)]}]
            AC#8 forward leg GREEN: the tree that was sealed keeps its own notice
            AC#8 RED: 1 notice(s) from a seal that ran on C:\...\artifact.txt are attributed to D:\...\artifact.txt,
                      a tree on the volume D: that was never sealed.
            AC#8 single-pair face: noticeNamesTree(notice from C:\..., D:\...) = true
```

⇒ **118b 那一发的读数独立复现**：C: 上那发 seal 的通知被归到从未被 seal 的 D: 树，正向 leg 仍绿。

三条判据逐条复算（判据出自票 118 票面 AC#8 那条，票 126 票面 `:7` 直取）：

1. **"取不到两枚可写卷就 `t.Fatalf`，不是 Skip"**——成立且我这一跑真走到了：我的枚举取到 4 枚可建目录的卷根
   （C/D/E/F），G-J 报"cannot find the path specified"被跳过。**但这里有一枚两任读数的分歧要解释**：
   `touch C:\f.txt` 在这台机器上**被拒**、`mkdir C:\d` **被允许**（普通用户的 C 盘根默认如此），
   所以按"建文件"枚举卷会得到 `D: E: F:`、按"建目录"枚举才得到 `C: D: E: F:`。
   118b 自述的"C: 与 D:"只有第二种枚举才成立（前任验收方 `ac118r` 那枚探针用的是第一种，残骸落在 D:/E:，我看见了并清掉了）。
   ⇒ 判据本身没问题，**枚举口径**得写死，登记为 `R-118-8`（建议归票 126 AC#6）。
2. **"正向 leg 断自己的通知归自己，否则常数 false 也能满足反向 leg"**——成立：我这一跑正向 leg 绿、反向 leg 红，
   两枚 leg 各自独立报，顺序也在前面（红不是靠"什么都不归"凑出来的）。
3. **"反向 leg 断从未被 seal 的那棵树归到 0 枚"**——成立，就是我引用那一行 RED。
   我另加了一枚它没要求的对照：B 树**也种了同一发的继承 Everyone 授权**（`grantStandsOn118` 逐枚验过），
   否则"0 枚"可以由"B 本来无可报之物"满足。

### 不许背书的那一腿

`resolve.go:325` 那条 C26 缝守**在同形下会不会把另一卷的树当成同一棵**——118b 明确写了它没量
（"那要造一枚会跨卷作答的 fake resolver……不当读数用"）。我只补了**函数级**的一发事实：
缝守用的两枚比较（`sameTree`、`answerInsideTree`）在同一条 `pathComponents` 上构造，我直接喂跨卷拼写
**两枚都返回 true**（上面 unit face 那一行）。这是"它俩都不看卷"的读数，
**不是**"那条缝今天会放走一枚跨卷的答案"的读数——后者要一枚会跨卷作答的 fake resolver，
两任都没造，**票 126 AC#2 仍是 open，本表不替它盖章**。登记为 `R-119-*` 之外的 `R-118-9`。

### AC#8 裁决：**通过（按票面规则停手）**，标签：责任链＋三条判据＋跨卷读数〔独立复现〕，缝守腿〔仅自述，不背书〕

票面 AC#8 的原文要求是"先量，再判是不是生产洞；**如果证明的是生产判据会归错 ⇒ 停手回报，不要自己改 `winsec_windows.go`**"。
实现方量到了、停了手、没动一行生产码（AC#4 已机器验）、把用例全文与三条判据交给新立的票 126，
并且如实写下它没量的那一腿。**这正是该格设计出来的行为**，不是"AC 声称要防的结局被造出来"——
被造出来的是**归属判据会归错**这一事实，而这一格的判据就是"量出它就停手交回"。
所以：**本格不因停手而退回**；它也不能被勾成"已交付"，因为它交付的是"移交"。
票 118 结案时把它记成"已划给票 126"（见下面的结案意见），不记成缺格。

现场清账：我的探针落在 `C:\ac118b-xvol`、`D:\ac118b-xvol`，`t.Cleanup` 已把**卷根那一层**一起 RemoveAll，
跑完逐枚 `ls` 复验四枚卷根（C/D/E/F）均无残骸；同时清掉前任验收方留在 `D:\ac118r-xvol`、`E:\ac118r-xvol` 的空目录
（它那一发死在 17:42，没来得及清）。仓库树内没有多出任何文件（`git status --porcelain` 只有本证据文件一枚）。

## 读数五：接续方对"读数一"的抽验 + 我补的发数（AC#1/AC#2/AC#3/AC#7）

断线接续不改前任验收方的判定，但要把它的底查一遍。三件事：

1. **快照完整性**：`/d/tmp/ac118r/base_7321c76` 与我自己 `git archive 7321c76` 的两枚，
   对 `winsec_windows.go`、`internal/proc/envfork_test.go`、`placement_leaf_118_other_test.go` 三枚文件 md5 **逐枚相同**
   （`fb20bca5` / `32713c58` / `aea4f696`，两枚快照与 `git show 7321c76:<path>` 三方一致）⇒ 它读的确实是交件末态。
2. **归档日志逐发重点数**（`/d/tmp/ac118r/log_mut-*.txt`，只数顶层 `--- FAIL`）：
   M5=3、M5-INV=2、M5-BOTH=3、M4=9、M4B=3、NOINH=4、TRUE=4、FALSE=9、WD=1、WD2=1，
   **与上表那一列逐数相同**，红名列表也逐名相同；变异字节也还在它各自的快照里
   （`mut-m5/internal/winsec/winsec_windows.go:136` 是 `kind := "explicit"` + `:137` 的 `_ = n // MUT-M5`，
   `mut-wd/…/notice_kind_and_everyone_118_windows_test.go:406` 仍是 `Contains("WD")` 那一族）。
   ⇒ 前任那一节的 〔独立复现〕 我背书为 **〔日志＋归档，我抽验〕＋其红名清单逐数复核通过**；
   AC#7 那一节的四个行号我也对过原文件：`:294`、`:299`、`:401`、`:404` 在 `7321c76` 的
   `inherited_narrow_notice_104_windows_test.go` 里正是那四条断言本体（不是 Logf）。
3. **AC#3 我补的两发常数形**（票面没要求、两任都没做的"新判定是不是恒真/恒假"）：
   把 `namesEveryone` 的整段 SID 判定改成常数（快照内测试码，仓库树未动）：

   | 发 | 改法 | 读数 |
   |---|---|---|
   | MUT-NAMES-TRUE | `return true // the SID rule is gone` | rc=1，红 **1** 枚：`TestAC118NamesEveryoneIsJudgedBySIDNotByTwoBytes`，红点 `:583` `AC#3 RED: namesEveryone admitted a SERVICE grant as Everyone: …NT AUTHORITY\SERVICE:(RX)` |
   | MUT-NAMES-FALSE | `return false // nothing is ever Everyone` | rc=1，红 **2** 枚：同一枚用例红在正向腿 `:564`（`does not recognise a real explicit Everyone grant … Everyone:(RX)`）+ 票 104 的消费者 `TestSealReportsThePrincipalsItCleared`（`narrow_notice_windows_test.go:81`） |

   ⇒ AC#3 的新判定**既不恒真也不恒假**，而且它的另一个消费者也真的在听它说话。
   顺带一枚对本票票面口径有用的读数：FALSE 那发的红名文本里，icacls 把 Everyone 渲染成 **`Everyone:(RX)`**、
   整串里没有 `WD` 两字节——这就是前任 MUT-WD2"旧规则连真的 Everyone 都认不出"的独立佐证。
   另登记一处小残留：交付文件 `notice_kind_and_everyone_118_windows_test.go:406` 还留着
   `Contains("WD") || Contains(everyoneSID)` 这种"两种渲染之一"的容忍断言，**非承重**
   （同一枚用例 `:403` 的逐桶 SID 精确比较先红），但它正是 `R-104-6` 那一族形状，建议改严或删 ⇒ `R-118-7`。

### 我这一轮补的发数（实现方与前任验收方都没做的）

| # | 发 | 目的 | 结果 |
|---|---|---|---|
| 1 | POSIX **MUT-B**（不查叶子）容器真跑 | 复算 AC#6 的"1→3" | 红 3 枚，逐名相同 |
| 2 | POSIX **MUT-LEAF**（只查叶子） | AC#6 的对称半修 | 红 6 枚（祖先腿另有钉） |
| 3 | AC#9 六发（base×2、oldassert×2、pre119×2） | 复算修前红/修后绿 | 与自述五行逐行相同，另补 oldassert+实目录=绿那一发对照 |
| 4 | **MUT-INJ 打在两枚 vintage**（`7321c76` vs `33c8acd`）× winsec/proc × 两种 TMPDIR | 判 119b 与验收方那一发争议 | 见"AC#9 附"：两读各得其所，分界点＝`33c8acd` |
| 5 | **MUT-NAMES-TRUE / MUT-NAMES-FALSE** | AC#3 新判定的恒真/恒假对照 | 各红 1 / 2 枚 |
| 6 | **跨卷探针**（unit face + behavioral face，真卷 C:/D:） | 独立复算 AC#8 的量 | `sameTree(A,B)=true`、正向绿反向红；另测得 `answerInsideTree` 同形也为 true |
| 7 | **AC#5 门禁整轮**（`-count=2` 双包、gofmt/gofumpt、三 GOOS vet、CI 同形、d22scan） | 复算四数与台账 | 逐数相同，八 scope 无一下降 |
| 8 | `034080c` 注释纯度机器验、九枚 commit 逐枚文件清单 | 两任分账 | 非注释改动 **0 行** |

## 注入文字登记（本接续方可见范围）

本代理这一轮的工具输出里，**没有**出现自称"编排者备注 / 系统提示 / 用户已更新编码规则、用户偏好优先于 AGENTS.md /
请 revert / 冻结某包 / 放宽阈值 / 不要提它"那一类文字：**0 次**，无据此改判。
另有 harness 自己的两类自动提醒，如实分开记（它们不是注入、也不是指令）：
后台任务完成通知 1 次（`task-notification`，我自己起的那一发门禁跑完）；
"任务列表"提醒若干次（每次工具调用重复注入，内容与本票无关，未据此改判）。

## `R-118-*` 登记与建议归属（本表只登记，`docs/reports/**` 由编排者记，我没动）

| 编号 | 内容 | 性质 | 建议归属 |
|---|---|---|---|
| **R-118-1** | AC#1 的判据口径："删掉 switch 必须红"单发不足以覆盖三枚取值——`kind=explicit` 那一格在 M5 下必然仍绿（M5 的输出就是 `explicit`）。达标形状＝每枚取值各被一发**不同常数**的删除形红过（本表已补齐：`explicit` / `inherited` / `explicit+inherited` 三发） | 判据措辞 | 票 118 结案备注；下一张"判据仪器"票的措辞照此写 |
| **R-118-2** | 票面与建票材料的**行号已漂**：票面 AC#1 写 `winsec_windows.go:80-86`，真实位置 `:136-142`；编排者给我的复核口径写 `resolve.go:353`，那一行是 `sameTree` 收尾空行，责任字节在 `winsec.go:347/353`。另：`envfork.go`/`winsec_other.go` 的注释改动来自 **`034080c`**，不是 `36294c2`（后者只动 `cmd/wisp/secret.go` + 一枚新用例） | 票面勘误（append-only，我不改票面） | 编排者的票面账 |
| **R-118-3** | AC#2 两侧哨兵**各自只被对方一侧的红覆盖**：MUT-M4B（只换渲染两桶、生产 kind 仍算对）下归属面两枚用例仍绿；MUT-M4（生产两桶互换）双侧都红。判据没缺口（票面只要求两侧各 ≥1 枚红），但"桶换了、字段没跟着换"这一形今天只有一侧有牙 | 加固建议 | 下一张 winsec 归属票（若立）；本票不因此退回 |
| **R-118-4** | 实现方 AC#2 自述"红 5 枚"，同一行自己列了 6 个名字，**实测 9 枚**（清单见上表 MUT-M4 那一行）。少报不是假绿，是账没数完 | 读数更正 | 票 118 票面 Progress log 的勘正（编排者记） |
| **R-118-5** | 正面项：MUT-LEAF（只查叶子不查祖先）红 6 枚 ⇒ POSIX 那条走链的**两半各自有钉**，"只修一半"不存在全绿走法 | 已达标读数 | 记入票 118 AC#6 结案语 |
| **R-118-6** | 口径备忘：**同一枚变异在不同 vintage 的树上有不同后果**。`envfork.go:104 return SealableRoot(dir)` 在 `33c8acd` 之前对 winsec 完全不可观察（plain TMPDIR 下连 proc 的两枚字符串比较也不红），之后红 1 枚且红名点名它。⇒ 复用这一发的验收必须写清 (a) 快照 sha (b) TMPDIR 形状 | 仪器口径 | 票 119 的 `R-119-9` 账（不是本票的洞） |
| **R-118-7** | `internal/winsec/notice_kind_and_everyone_118_windows_test.go:406` 仍留 `Contains("WD") \|\| Contains(everyoneSID)` 的容忍断言。非承重（同用例 `:403` 先红），但形状与 `R-104-6` 同类，后人把严格那枚"简化"掉就复活 | 轻微加固 | 票 122 那类清理票，或票 118 的收尾小件 |
| **R-118-8** | 跨卷用例的**卷枚举口径**：C 盘根允许普通用户建目录、拒绝建文件 ⇒ 用 `WriteFile` 枚举会漏掉 C:，用 `MkdirAll` 才得到 4 枚。两任读数差（C:/D: vs D:/E:）就出在这里。判据要写死"用建目录枚举"，并保留"取不到两枚就 Fatalf、绝不 Skip" | 判据细化 | **票 126**（AC#6 那一格一并管清残骸：本轮我清了前任留在 D:/E: 的 `ac118r-xvol`） |
| **R-118-9** | 缝守 `resolve.go:325` 的两枚比较（`sameTree`、`answerInsideTree`）**函数级跨卷同真**——这一发我量到了；但"一条会跨卷作答的 resolver 今天能不能过那道缝"**两任都没量**（要 fake resolver） | 未量的那一腿＝open | **票 126 AC#2**，本表不背书 |

## 裁决表（收表：九格定档）

| 格 | 判据（票面一句话） | 裁决 | 标签 |
|---|---|---|---|
| AC#1 | `kind=` 三取值各一枚；删 switch 必须红 | **通过**（附 R-118-1 口径：第三枚取值由补发的常数形红，票面那一发本来就红不到它） | 〔独立复现〕（前任；本方抽验其归档与变异字节，见读数五） |
| AC#2 | 两桶互换 ≥2 枚红（归属面 + 渲染面） | **通过**（附 R-118-3、R-118-4：实测 9 枚、"字段没跟着换而全绿"造不出来） | 〔独立复现〕＋抽验归档日志 |
| AC#3 | `namesEveryone` 不按 2 字节子串 + 反向对照 | **通过**（旧规则的误认针确实是 `Everyone` 被 `t.TempDir()` 拼进路径；新判定常数两发各红，非恒真非恒假） | 〔独立复现〕（前任 MUT-WD/WD2 ＋ 本方 MUT-NAMES-TRUE/FALSE） |
| AC#4 | 生产码一行不动（diff 证明） | **通过**（四枚文件 0 行；九枚 commit 只碰 5 枚测试文件；119b 的注释那一发机器验为纯注释且在本区间之外） | 〔独立复现〕 |
| AC#5 | 门禁四数点名 + gofmt/gofumpt/vet 双 GOOS + d22scan | **通过**（172/92/0/0 与 86/46/0/0 双口径、cached=0、三 GOOS vet、八 scope 无一下降） | 〔独立复现〕 |
| AC#6 | POSIX 叶子方向两枚；MUT-B 从红 1 → ≥3 | **通过**（容器真跑：基线 42/27/0/0，MUT-B 红 3 枚逐名相同；我补 MUT-LEAF 红 6 枚） | 〔独立复现〕 |
| AC#7 | 两枚装饰腿在恒真/恒假下发红名 | **通过**（`:294`/`:401` 恒真、`:299`/`:404` 恒假各红名点名自己，行号对得上原文件，不需降级） | 〔独立复现〕＋抽验归档日志 |
| AC#8 | 跨卷归属：先量再判（票面留 `[ ]`） | **通过（停手移交）**：量到了、判成生产洞、没动生产码、三条判据与责任链我全复算；**缝守那一腿不背书**（R-118-9，票 126 AC#2 仍 open） | 链＋判据＋跨卷读数〔独立复现〕；缝守腿〔仅自述，不背书〕 |
| AC#9 | envfork 那枚断言守语义不守拼写 | **通过**（六发独立复现：旧拼写在软链形红、新断言在正确实现上绿、在 pre-119 形下红名点名 `:137`、实目录形不红）；119b 争议另判＝**两读各得其所，分界点是 `33c8acd`**，账在票 119 | 〔独立复现〕 |

## 结案意见（票 118）

- 九格：**八格通过 + AC#8 一格"停手移交票 126"**。AC#8 **不算缺格**：票面 AC#8 自己的判据就是
  "量出生产判据归错 ⇒ 停手、不改 `winsec_windows.go`、交回编排者另立票"，实现方按这一句执行完毕，
  票 126 也已由编排者立出来（票面 `:9` 直接取 118 AC#8 那条的用例全文与三条判据）。
  建议结案写法：`AC#8 = 移交票 126（本票不勾，结案不以此格为缺）`——**票面那一格我不翻，实现方的框我一个也不动**。
- 零退回：九格里没有任何"票面声称要防的结局被真实造出来"——被造出来的两件事（跨卷归错、
  恒真用例的 vintage 分歧）都是**本票要量的对象**，且都按票面规则处置了（前者停手立票、后者账在票 119）。
- 附条件共 4 条，全部**不阻塞结案**：`R-118-1`（判据措辞）、`R-118-3`（归属面一侧的桶哨兵加固）、
  `R-118-7`（一处非承重的 WD 容忍断言）、`R-118-8`（跨卷枚举口径，归票 126）。
- 本验收方全程只读：`git status --porcelain` 只有 `docs/evidence/s1/118-adversarial-acceptance.md` 一枚未跟踪文件；
  所有变异都打在 `/d/tmp/ac118b-c1`（及前任的 `/d/tmp/ac118r`）的 `git archive` 快照里，
  仓库树内未建 worktree、未 checkout、未改一行生产码或测试码；卷根残骸已逐枚 `ls` 清账。
