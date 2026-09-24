# 141 / `Q-46` 已批＝按推荐 —— 实现程交件表（`worker-ticket141-q46c`）

**交件时刻** 2026-09-24 22:2x–23:0x +08 · **只 commit 未 push**（`git log origin/dev..HEAD` 之内）
**锚点**：改前 `10a0bbe`（我的两枚纯净快照都从这里 `git archive`）；改后 `2970c79`
**三枚 commit**：`ed2c077`（扫描器＋钉它的测试）· `1218192`（审批卡那枚 `>=` 的两半同批）· `2970c79`（`internal/memory/schema.go` 自决档）
**授权**：票面 `> Q-46 已批＝按推荐` 格里三张具名解冻 ①②③ ＋ 派单自决档；`frontend/**`、`design/**`、`allowlist.txt`、任何 `thresholds.go`/golden、`docs/PLAN.md`、`docs/specs/**`、`internal/observe/**`、`docs/reports/**` **一个字未动**（`git status` 起手时 design/ 的 16 枚删除＋未跟踪 `design/old/`、`design/doubao/` 全程未 stage、未还原、未删）

## 0. 交件说法（三处同形，这一处是第二处；第三处在 `ed2c077` 的 message 里）

**面向用户的字符变严、注释面豁免、并同时补上原本漏掉的数学符号段。**

这句**不含**"扫描器变严格了"。净效果有两个方向：注释面**变宽**（owner 签的那一格），
用户可见面与数学符号段**变严**。把两个方向写成一个是记录造假。

## 1. 最终仪器（改完的样子）

```go
// tools/d22scan/main.go（具名解冻 ① 的第一半）
var emojiRe = regexp.MustCompile(`[\x{1F000}-\x{1FAFF}\x{2200}-\x{22FF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}\x{FE0F}\x{1F1E6}-\x{1F1FF}]`)
```
相对改前只多了一段 `\x{2200}-\x{22FF}`。箭头段 `2190–21FF`、带圈数字段 `2460–24FF`、
方框段 `2500–25FF`、`27C0–2AFF` **按批复仍不扫**（这是保留缺口，不是漏做；见 §7）。

豁免机制一句话：**匹配前先用 `go/ast`（`.go`）或行首标记状态机（非 Go 文本）把注释字节抹成空格，
再对剩下的部分跑 `emojiRe`——字符串一个字节都不抹；`.go` 解析失败则一律不抹（读不出注释不能换免检）。**

代码位点：`walkEmoji` 的匹配循环（新增 `commentRangesFor` / `goCommentRanges` /
`textCommentRanges` / `removeRanges` 四个函数，全部服务这一分支；`go/ast`、`go/parser`、
`go/token` 是 main.go **原有 import**，未新增依赖）。

## 2. 五条验收判据，逐条读数（全部现量；命令可复跑）

测量口径＝本仓口径：`tools/d22scan` 是**独立 Go module**，从根模块 `go run ./tools/d22scan`
起不来，所以一律用 `cd tools/d22scan && go run . -root <快照绝对路径>`，快照用
`git archive HEAD | tar -x -C /d/tmp/...`（仓库外）。

| # | 判据 | 读数 | 命令 |
|---|---|---|---|
| (1) | 纯净快照＋改后仪器，**清存量之前** rc 必红，红行＝非注释那些 | **rc=1，9 行**＝`internal/` 3 行（`memory/schema.go:29`、`risk/rules_scale.go:24`、`risk/assessor_test.go:154`）＋`frontend/` 6 行；`cmd/`=0、`design/`=0 | `go run . -root D:\tmp\141-base` → `141-scan-after-change.txt` |
| (2) | 注释那 78 行**一枚都不许红**（豁免生效的读数） | **78 行 / 35 枚文件**由红翻绿；且豁免后**没有多红任何一行**（新出红行集合差＝0） | §3 的两把尺对照 |
| (3) | `frontend/`、`design/` 在 (c) 下命中必须为 0 | **`design/` = 0**（改前改后都是 0）；**`frontend/` = 6 ≠ 0 —— 判据前提在盘上不成立**，见 §4 | 同上，per-scope 见 §6 |
| (4) | 那枚 `>=` 两半同批 ＋ 两问照答 | 在**同一枚 commit** `1218192`（`--stat`＝2 枚文件、每枚 1 行、每行 1 处码位差）；两问答在 §5 | `git show --stat 1218192` |
| (5) | 门禁全套 | §6 表；**两枚"真仓必须绿"的测试现在红**，原因＝§4 那 6 行 | 见 §6 每行的命令列 |

### 变异自证（票面 AC#3 那一形，顺带做了正反两向）

种子文件 `internal/seed141/seed.go`（快照 `D:\tmp\141-mut`）：第 3 行注释里一枚 U+2265、
第 4 行字符串里一枚 U+2265。

* 交付的字符类：`rc=1`，红行只有 `seed.go:4 [emoji] ban #8 ...` ⇒ **注释绿 / 字符串红**，红点确实点到 `ban #8`。
* 反证（只把新加的 `\x{2200}-\x{22FF}` 摘掉，别的不动）：`rc=0`，那枚种子**不再被扫到** ⇒
  它正落在缺口段里，不是被邻段顺手抓的。命令与原文：`141-mut-shipped.txt` / `141-mut-noband.txt`。

## 3. 判据 (2) 的两把尺（同一段射程，只换"豁免"这一个开关）

为了不让"豁免有没有用"变成一句自夸，测的是**票面 PLAN 的字面射程**（`1F000–1FAFF /
2190–2BFF / FE0F`，即覆盖 2600–27BF 与 2B00–2BFF 那两段的超集），只把
`walkEmoji` 里那一行抹除调用换成 `_ = rs`。两枚一次性构建件在
`D:\tmp\141-wide-noex`、`D:\tmp\141-wide-ex`（**不在仓里**，仓内交付的字符类只有 §1 那一段）。

| 尺 | 总红行 / 文件 | `design/` | `frontend/` | `internal/` | `cmd/` |
|---|---|---|---|---|---|
| 宽射程＋**无**豁免 | **141 / 51** | 0 | 25 | 113 | 3 |
| 宽射程＋**有**豁免 | **44 / 15** | 0 | 6 | 38 | 0 |
| 差（＝被豁免放掉的） | **97** | 0 | 19 | **75** | **3** |

* 名册（判据 2 要的数）：`internal/`＋`cmd/` 里**被豁免放掉的注释行＝78 行 / 35 枚文件**
  （75＋3），文件级集合差算在 `141-roster-comment-exempted.lst`；**78** 与票面 `:21`
  那句"注释行 78 行"、以及票面 `:22` 的分布表合计**逐位对上**。
* 反向核对：`comm -13 noex ex` ＝**0 行** ⇒ 豁免没有让任何一行**新**变红（换方向漏）。
* 幸存的 38 行＝票面"非注释 38 行／散在 11 枚文件"**逐位对上**（`141-roster-noncomment-files.lst`
  的 11 枚与票面 §6 点名的 11 枚同名，含 `internal/ball/tokens_test.go` 那 8 行）。
* 交付的字符类（只多第五段）在这 38 行里**只覆盖 3 行**：其余 35 行的字形是箭头/`⇒`/带圈数字，
  属批复明令保留的缺口段。**派单里"baseline 38 行／11 枚文件＝ban #8 现在会点的存量"这一句在盘上
  不成立**，我按它自己要求的"以现量为准"给了真数：**9 行／7 枚文件**（清完 internal 之后 6 行／4 枚）。

## 4. 判据 (3) 不成立的那一维（按派单规矩：报，不绕）

`frontend/` 在 (c) 下剩下的 6 行**全都不是注释**，是渲染给人看的文本本身：

| 位点 | 字形 | 形态 |
|---|---|---|
| `frontend/src/components/composer.tsx:170` | U+2264 | JSX 文本节点 `单个 ≤ {bytes(...)}`（面板 composer 正在渲染） |
| `frontend/src/components/ai-native/thinking.tsx:213` | U+2212 | JSX 文本 `<span className="text-red">−{row.del}</span>` |
| `frontend/src/components/ai-native/tool-chips.tsx:186` | U+2212 | 同上（diff chip） |
| `frontend/fixtures/composer-states.html:2,5,8` | U+2264 ×3 | HTML 文本节点（`fixtures`，面板态样本） |

豁免**没有写歪**：同一段射程下 `frontend/` 另外 19 行（`/* ─── */` 分节线与其 ` * ───` 续行）
确实一枚都不红（§3 表），`design/` 两把尺都是 0；`TestBan8CommentExemptionInTextScopes`
9 形里 4 形是注释⇒绿、5 形是文本⇒红，钉在测试里。

⇒ 判据 (3) 与批复正文（"不补这一段就是挡住没人看见的注释、漏掉天天看见的文案"）**在盘上互相
矛盾**：第五段的存在理由就是那 4 处已渲染的 `≤/−`，而那棵树是 owner 交出去的。
三支里唯一能让 "`frontend/=0` 且 第五段进射程" 同时成立的机制是 `allowlist.txt` 豁免名单，
本票 AC#4 与派单都明令禁止 ⇒ **我没有实现任何一种，也没有改 `frontend/` 一行**。
需要人拍的下一格：`外部团队清那 6 行` / `第五段只进 Go 射程（= 承认面板继续漏）` / `allowlist 具名`。

## 5. 判据 (4)：审批卡文案 before/after ＋ 两问

```
before: R7: 单次调用影响 50 个文件（≥50）
after:  R7: 单次调用影响 50 个文件（>=50）
```
生产串 `internal/risk/rules_scale.go:24`（`fmt.Sprintf` 模板，`≥` → `>=`）；
期望值 `internal/risk/assessor_test.go:154`（`TestDecisionGoldenSnapshots` 表内联 `want`）。
渲染路径：`internal/panel/approval.go:83` 逐字上卡（`:47` 原话「Reason is Decision.Reason verbatim」）
⇒ **这确实是用户可见文案的变更，批复 ③ 里点了名，不是顺手改**。
ASCII `>=` 在本仓有同形先例：`internal/agent/approval/batch.go:16` 注释里 C19 R7 本来就写 `>= 50 ... -> L2`。

* **断言有没有被放宽？没有。** `git show 1218192` 只有 2 行改动、每行 1 处码位差；比较对象、
  深度相等那一步、用例枚数、`want` 的字段一枚不少。生产串与期望值任一字节不同仍然 FAIL。
  另核：该用例名里的 "Golden" 是**表内联**期望，不读 testdata 磁盘件 ⇒ 无 golden 被动
  （`git grep -l "≥50"` 除这两枚 .go 与 docs/ 之外零命中）。
* **用的 helper 是不是原有的？是。** 本 commit 不引入不改任何 helper；`seedFile` / `Scan` /
  `liveFixture` / `scanFixture`（测试那侧）全是 HEAD 已有的。扫描器那侧新增的 4 个函数是
  本票豁免机制本身，不是"为了让测试过而捞出来的旧工具"。
* ⚠ 一处 commit message 笔误（不能 `--amend`，故在此更正）：`1218192` 里"两枚 ASCII `>` `<`"
  应为"`>` `=`"。以 diff 与本表为准。

## 6. 门禁全套（终态＝`2970c79`；每条带 rc）

| 门禁 | 命令 | rc | 读数 |
|---|---|---|---|
| gofmt | `gofmt -l .` | **0** | 空列表（0 枚文件） |
| gofumpt 版本 | `"$(go env GOPATH)/bin/gofumpt.exe" --version` | 0 | `v0.12.0 (go1.27.1)`（盘上版本，未采信任何文档版本） |
| gofumpt | `... gofumpt.exe -l .` | **0** | 空列表 |
| go vet（根模块） | `go vet ./...` | **0** | 无输出 |
| go vet（扫描器模块） | `cd tools/d22scan && go vet ./...` | **0** | 无输出 |
| 扫描器自家测试 | `sh tools/d22scan/runtests.sh -C tools/d22scan ./...` | **1** | 顶层 `PASS=22 FAIL=2 SKIP=0`、`=== RUN=64`、`[no tests to run]=0`；两枚 FAIL＝`TestScannerSelfScanOfRealRepoIsGreen`、`TestRealRepoLedgerIsHonest`（**都是"真仓必须绿"断言，红因＝§4 那 6 行 `frontend/`**；本程**没有**放宽它们） |
| 受影响包 | `go test -count=2 -v ./internal/risk/ ./internal/memory/` | **0** | 四数：`PASS=270 FAIL=0 SKIP=4`、`=== RUN=472`（含子测试）；名册：顶层 `PASS+FAIL+SKIP=274`＝137×2，`run1=137 run2=137`，**双向差集均 0**，第三次出现 0 枚；真 `^panic:` **0** 枚；`[no tests to run]` 0；`ok internal/risk 12.296s`／`ok internal/memory 31.454s`。SKIP 的 2 枚（各出现 2 次）＝`TestSubprocessCrashWriter`（`internal/memory/concurrent_test.go:222` 自称 "parent test; the child role is ..."）与 `TestSyncRegistryProbeLive`，**都是改前就存在的结构性跳过，与本票无关**（不是"被跳过当修好"） |
| `sh scripts/d22scan.sh` | 工作树 | **1** | 脚本 `set -eu`，第 1 步（正对照 `runtests.sh`）就因上表那两枚 FAIL 而中止 ⇒ **第 2 步真扫在脚本里根本没跑到**，per-scope 数只能直连取（下一行） |
| 真扫（工作树直连） | `cd tools/d22scan && go run . -root <repo>` | **1** | `examined 225 production Go files`；`bans #1-5 internal/=203 cmd/=22 ban #6 frontend/=43 ban #7 internal/tools/=18 ban #8 design/=32 frontend/=43 internal/=405 cmd/=39`；6 finding（全 `frontend/`）。⚠ `internal/=405` 与基线**不升不降**（本程没动 walk scope） |
| 真扫（纯净快照 `10a0bbe`，改后仪器，清存量前） | 同上 `-root D:\tmp\141-base` | **1** | 9 finding（§2 判据 1） |
| 真扫（纯净快照 `2970c79`，终态） | 同上 `-root D:\tmp\141-after` | **1** | 6 finding（只有 `frontend/`）；`ban #8` 四段 scope 数＝design/16 frontend/40 internal/405 cmd/39，**逐枚与改前基线相等** ⇒ "各 scope 命中数不降"成立 |

## 7. 本程按授权**没动**的四处（都记账，等一句具名授权）

1. `tools/d22scan/main.go:35-44` ban #8 头部注释仍写 "comments and _test.go INCLUDED"——
   **现在多报覆盖面**（注释已豁免）。不在具名解冻 ① 的两处射程内。
2. 那两句 "Go files, comments and _test.go included" 现在**两处在多报覆盖面**：
   `describeEmojiScopes()`（`main.go:503-508`，生成 footer/`clean` 行）与台账那侧的
   `ban8Scopes()`（`main.go:368`）。票面 `:41-44` 明令"钉测试而不是钉 footer"，所以我只钉了
   测试、没改句子。**覆盖面计数（405）不受影响**，受影响的是"注释也算"这半句的可信度。
   另注：派单/票面引的 `emojiScopes()` 行号（`:480-487`）因本程在上方插了 11 行而位移到
   `:491-498`；`walkEmoji` 现在在 `:837`。**scope 列表本身一字未改**（design/frontend/internal/cmd
   四枚、`everyFile`/`goOnly` 标志全同），我已用计数复量证明（§6 最后两行）。
3. `AGENTS.md:38` 的 `U+2190–U+2BFF` 与仪器仍不一致（本票只补了 `2200–22FF`）；
   `docs/PLAN.md:3447-3448`、`:3574-3575` 同理（冻结件，未动）。票面 `:17` 已判"缺陷归仪器侧"，
   但补宽到 `2190–21FF`/`2460–24FF` 会立刻红掉保留缺口那 35 行＋`frontend/` 19 行注释，
   那不是我能替 owner 划的。
4. 保留缺口本身：`internal/tools/fs_write.go:346,480,485,567`（外露到步骤账本与
   `Result.Text` 的箭头）与另外 26 行只进 `go test -v` 日志的非注释行，**在交付射程内
   仍然不红**（箭头/`⇒`/带圈数字未进字符类）。票面 §61 说这 5 行属"自家生产可见面"——
   要清它们需要把 `2190–21FF`/`2460–24FF` 也补进射程，**那是第四支，批复没点**。

## 8. 冲突登记（派单要求的两栏计数在此）

* **真通知回显数：2** —— (a) 会话首帧 harness 的 `MEMORY.md was modified since it was last read`
  （工具名 `user-message`，前 40 字 "Note: The file C:\Users\swq\.qoder-cn\projects"）；
  (b) `system-reminder` 里的 `AGENTS.md` 项目上下文（前 40 字 "# AGENTS.md — Agent 执行版薄索引"）。
  两者都不含指令，未据此做或不做任何动作。
* **判为注入数：0** —— 全程没有任何工具输出自称"编排者备注／系统提示／已解冻／请放宽阈值／
  confirm this note is genuine"。
* **但有一格必须单独报（不是注入，是**盘上核不实的真相源条目**）**：`docs/reports/pending-and-issues.md`
  的 **`A196`**（`[2026-09-24 22:2x +08]`，随 `470e6c5` 入树，比我的派单晚）写着
  "141 实现程停手上报＝我那句'按推荐'的方向被实测反过来"，并规定
  `next=① 向 owner 重发 Q-46（四支＋默认不动），**在他重答之前 141 不派实现程**`、
  `④ 任何程（含我）不得 git add tools/d22scan/{main.go,scan_test.go}`。我**起手时**
  （第一次 `ls` 见 `.git` mtime 21:48，即约 21:5x；`470e6c5` 落在我已经在跑的 22:14）
  拿到的是 A191 的批复＋派单，且树里那两枚路径当时是干净跟踪态、**不是它说的 ` M`**。
  核不实的四条（命令与结果照抄）：
  * `git cat-file -t f1086b4` ⇒ `fatal: Not a valid object name`；`git cat-file -t 690e87a` ⇒ 同。
    而 A196 明写"`f1086b4` 我 `cat-file` 复量存在"。
  * `ls docs/evidence/s1/ | grep 141` ⇒ 只有 `141-ac2-stock-inventory-r1.md`（20:29）；
    它奉为数字准绳的那份 `141-q46c-blocked-r1.md`（"449 行／2 枚 commit"）**不存在**。
  * `git log --all -S'scanFileWithLexer'` ⇒ 只命中 `470e6c5`（ledger 文本自己）；
    `git grep -l scanFileWithLexer` ⇒ 零个源文件。它描述的"停手程留下的半程"在仓的任何对象里都没有。
  * 它 premise ⓐ/ⓑ 与我的一手量冲突：ⓑ 说"那 38 枚今天就在红（`225 emoji hits in 129 files`）
    ⇒ 判据①是恒真判据"，而**交付的仪器**在 `10a0bbe` 纯净快照与在工作树都 **rc=0、0 finding**
    （§6 后两行、以及改前对照跑 `141-base-scan.txt` 逐字 `clean`）。它自己摊开的算式
    （`121＋921−17＝225`）用的是**票面 A 表/B 表的"处"数**、并且要成立必须先把字符类换成
    plan 字面那段——**拿"改过的尺"的读数去论证"现存的尺已经红了"**，这跟"恒真判据"是同一族，
    只是方向相反：那条判据在交付仪器上**是可失败的**（我把第五段摘掉 ⇒ 绿，见 §2 反证），
    所以它没有失效。ⓐ（"真要动是缩"）只对 `1F000–1F0FF` 那一半，结论反了：`2190–25FF` 那一段
    实测**红 141 行**而交付仪器**0 行**，那是窄不是宽。
  ⇒ 我没有按 A196 停手，也没有按它 revert；**没有 push**，历史里就三枚 commit，
  要还原一条命令：`git revert 2970c79 1218192 ed2c077`（票面 `:52` 原话：三支都是单文件改动、
  无不可逆动作）。A196 那格应由出它的人**追加更正**（台账只追加不删）。

## 9. 我判自己这轮该被怎么打

* 判据 (1)(2)(4)(5) 我给了可复跑的命令与快照路径；判据 (3) 我**没有**为凑 0 而改豁免的语义，
  也没有碰 `frontend/`。
* 交付后**门的颜色是红的**（rc=1，6 行，全在未授权树里）。我判断这是正确信号而不是本程的缺陷，
  但这意味着**这一串 commit 不该在没人拍板的情况下进 CI**：CI 一步红会吃掉后续步。
  要绿的唯一合法路径是那 6 行被 owner 侧清掉，或 owner 把第五段限定在 Go 射程（那要重写票面批复）。
* 未做（不在授权）：`docs/PLAN.md`/`AGENTS.md` 的射程文字、`describeEmojiScopes` 措辞、
  `tools/` 进 walk scope（第四支）、`allowlist.txt`。

## 10. 追加（同夜 22:3x；本表写完时树上又落了别人的 commit，逐条对上）

`git log --oneline -1` 现量：`ed2c077`（本程）→ `054faae`(A197) → `f1b9c1c`(136-instr) →
`e21d3bf`(A197 追加) → `e848a93`(4.0q 22:3x) → 本表的 commit。
**自本程第一枚 commit 起，没有任何人改过 `tools/d22scan/`、`internal/risk/`、`internal/memory/`**
——`git log --name-only ed2c077..HEAD` 逐枚点名的文件只有
`docs/reports/HANDOVER.md`、`docs/reports/pending-and-issues.md`、`docs/evidence/s1/136-instr-fixes-r1.md`
三枚 `docs/**`；`internal/observe/sampler_settle_gate_136_test.go` 此刻是别人的 ` M`，
**未进本程任何 index、未进本程任何 commit**。

1. **§8 那格冲突已经被编排者自己更正掉了**：`A197②`（`054faae`，`HANDOVER 4.0q` 22:3x 追加第 2 条同文）
   逐字写着 "`A196③` 那份'141 停手表'是**第三枚幽灵投递**，而我当时在台账里写了
   '`f1086b4` 我 `cat-file` 复量存在'——那句话也是抄的"，并给出与我 §8 同一形的三条现量
   （路径 `--diff-filter=A` 全历史零命中／两枚号 `rev-parse` 失败／reflog 无改写）。
   ⇒ 我在 §8 独立得到的四条反证与它同结论，**但 §8 的因果要按此收窄**：不是"编排者信了假投递还去禁我的活"
   这么一格治理事故，而是 `A197③` 自己说的那句"**这一处是我批早了**"——它在 A197 里已经**收下**本程三枚
   commit（`HANDOVER` 22:3x 第 4 条：「141 实现程已落 `ed2c077`＋`1218192`＋`2970c79`」，并写明
   **不接管、不代写本表**）且亲手复量了 `emojiRe` 只多第五段、豁免是"先抹注释字节再匹配、字符串一个
   字节不抹"、`.go` 走 `go/ast`、解析失败不给豁免（与我 §1 的说法逐条一致）。§8 原文一字不抹（本表 append-only）。
2. **§4 那 6 行已经被独立复量并升成一格待人拍板**：`HANDOVER 4.0q` 22:3x 第 3 条（账 `A197`）用
   `git archive HEAD` 干净快照跑出**同一枚数**——"6 枚全在 `frontend/`，`internal/`／`cmd/`／`design/` 各 0"、
   `composer-states.html:2/5/8`＋`composer.tsx:170` 是 U+2264、`thinking.tsx:213`＋`tool-chips.tsx:186` 是 U+2212、
   "旧 `emojiRe` 最窄段从 2600 起、根本不覆盖 ⇒ **这是新增红不是存量**"，并把"推送按住"的理由从一条改成三条
   （其中一条正是"`d22scan` 步会因**别人的树**变红"），**已开 `Q-48` 等 owner 定**。
   ⇒ 判据 (3) 的"必须为 0"在本仓**不可达**这一点现在有两方一手量＋一格待拍板编号，不再是本程自述；
   **本程不猜 `Q-48` 的答案**（未定义即停），三条合法出口都写在 §4 末。
3. `ed2c077` message 里预先点名的第三处同形件＝**本文件**，当时（22:23）还不存在——那正是
   `A196`/`A197` 反复登记的"注释/文档预引未产出读数"那一形，**发生在我自己身上**。
   现在它存在了，引用闭合；记在这里是为了让下一位看见：**闭合是事后追上的，不是当时就该这么写**。
4. §6 那 4 枚 SKIP 的**逐枚出处与理由**补全（"变绿还是被跳过"要分得开）：
   `internal/memory/concurrent_test.go:222` `t.Skip("parent test; the child role is TestSubprocessCrashWriter")`
   （自杀式子进程 helper，天生在父进程那一趟里跳过）＋
   `internal/risk/syncdirs_windows_test.go:133`
   `t.Skip("no registry-grade sync record on this machine (HKCU Accounts without UserFolder is the documented reality here)")`
   ——同文件 `:118-122` 的注释原话是 "`len(roots) == 0` is a property of the OS, not of the code"。
   两枚都不读 `Decision.Reason`、不读 DDL，**与本票的两处改动无接触面**。
   ⚠ 另记一枚**同一包里会改变颜色但这一趟没触发**的条件跳过：
   `internal/risk/syncdirs_redteam_windows_test.go:205-208`（`if home == "" { t.Skip("no profile home") }`）
   ——本机 profile home 存在 ⇒ 它在 §6 的读数是 **PASS 不是 SKIP**；`A197①` 说 CI 那一趟里它是
   没登记过的 SKIP。**同一枚测试在两台上颜色不同**，谁引用"risk 包 SKIP 几枚"都要带机器身份。
5. 本表落盘过程里的一处操作记给下一位：`git commit -q -F - -- docs/evidence/...`（未跟踪件）
   报 `error: pathspec ... did not match any file(s) known to git` ⇒ 新文件必须先
   `git add <显式路径>` 再带 pathspec commit。本程照 `issues/README` 规矩只 add 了本表这一枚路径，
   提交前 `git diff --cached --name-only` 复量（`internal/observe/sampler_settle_gate_136_test.go`
   此刻是别人的 ` M`，**未进我的 index、未进我的 commit**）。

