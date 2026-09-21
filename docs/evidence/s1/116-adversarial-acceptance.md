# 票 116 对抗验收 —— `acceptor-ticket116`

**被验 sha**：`80e248c`（用例）＋ `63fdbd2`（票面 append）。锚定复算/变异一律 `git archive 80e248c | tar -x`。
**我的快照**：`/tmp/ac116-acc116`（= `C:/Users/swq/AppData/Local/Temp/ac116-acc116`，A38④：仓库目录内不建 worktree/checkout）。
**纪律**：未改一行生产码、未 commit、未 push、未改票面、未动台账；写权限只用于本证据文件与 `/tmp`。
**脏工作树隔离**：全程只读快照；`internal/winsec/`、`cmd/wisp/` 此刻有别的票在写，一律未碰、未用作被验版本。
**快照纯净证明**：`git show 80e248c:internal/risk/syncdirs.go | diff - <snapshot>/internal/risk/syncdirs.go` 空；
`pathresolver.go` 同样逐字节相同；变异跑完还原后 `grep -c MUTATION-116` = 0（两文件）。

## 环境读数
- 起笔 `date`：2026-09-21 21:52 CST；门禁段 22:06；POSIX 段 22:08 CST。
- `go version go1.27.1 windows/amd64`；`nproc`=12；门禁时 `% Processor Time`≈66%（**有并发负载**，与 `R-116-1` 同根，见 R 表）。
- docker server 29.6.2；`golang:1.27` 在；`GOMODCACHE=D:\work\base\gopath\pkg\mod`（POSIX 挂载复用，容器内先 `ls -l /wisp/go.mod` 证非空挂载）。

## 三档证据口径
〔独立复现〕＝我亲手在这台机/这台机 docker 里跑过、命令与读数在下方。
〔日志＋归档，抽验〕＝实现方日志我抽验过关键读数，未逐行重跑。
〔仅自述，不背书〕＝我没跑、也不背书的格子。

## AC#1 —— 这枚用例是"只钉符号"的第二版吗？红得有没有意义〔独立复现〕
**判：不是钉符号。红在行为（`Sync` 判定翻转 / `Why` 归属），且票 105 那两条在 M1 下仍 2/2 PASS。**

新用例只此一枚文件 `internal/risk/syncdirs_ancestor_actable_leg_116_test.go`（4 枚顶层用例，无 build tag）。
`git diff --name-status 88d8956..80e248c -- internal/risk/` = 单一 `A`（我复算过）；`syncdirs.go` 在 80e248c 与快照逐字节相同。

**锚点先证原位（动手前）**：`awk 'NR==226' | cat -A` → `^IanceCanon, err := ares.Actable()$`（tab 前导、行尾无杂字节），与自述一致。

**基线（纯净 80e248c，Windows 原生，`-count=2 -v` 六枚：4×新＋票105的2枚）**：`RUN=12 PASS=12 FAIL=0 SKIP=0`（6 名×2，票105两枚因 `mklink /J` 免管理员**真跑**、非 skip）。

**三发变异（每发：改 → 同一条命令里 grep 证明那几字节真落进文件 → `go build ./internal/risk/` rc=0 → 才读红名）：**

| 刀 | 落地证明 | build | 红的测试名（-count=2 全 2/2） | 红形 |
|----|----------|-------|------------------------------|------|
| **M1** `:226-229 → anceCanon := ares.Canonical` | `grep -n MUTATION-116-M1`→226；原 `anceCanon, err := ares.Actable()` count=0 | rc=0，vet rc=0 | **唯一红** `TestSyncAncestorActableLegFailsClosedOnMovedAncestor116` | **行为**：`verdict: Sync=false Root={} Why="write target is not under any sync root"`（第220行断言红）＋ 第229行 `Why` 丢 `expansion moved` |
| **M2** `:226-229 → anceCanon, _ := ares.Actable()` | `grep -n MUTATION-116-M2`→226 | rc=0 | 同一枚唯一红 | 红在**第三条断言**（第229行）：`Sync=true suspect-fallback` 仍在，但 `Why` 退化为通用 unverified-ancestor 文案、丢 `expansion moved` |
| **M1 on POSIX**（docker golang:1.27，容器内先 `ls -l /wisp/go.mod` 证非空挂载，`grep` 再证 MUTATION 在容器侧） | 容器内 grep→226 | 过 | 同一枚唯一红（其余三枚 PASS），`-run 116 -count=1` RC=1 | **行为**：`Sync=false`，raw 形 `/…/001/%WISP116ANC//B%/new.md`（靠 `//` 折叠出同一配对差） |

**M1 之下票 105 那两条（本票存在理由的直接证据）——我独立复算＝仍 PASS/PASS：**
`--- PASS: TestSyncAncestorLegReResolvesThroughJunctionToRealTree`×2、`--- PASS: TestSyncAncestorLegNotExemptedStaysFailClosed`×2，
且三枚 AC#2 用例（env/home/同树）全绿 ⇒ **红名唯一、只点名祖先那枚，不是"顺带红"、不是"包被改坏全红"**。
⇒ 复现 `acceptor-ticket105b` 那格读数（`Sync=false`、`Why="write target is not under any sync root"`）**逐字同形**，
说明 M1 真把"那条腿被删→写被放行不疑"的结局造了出来并被行为用例接住。

## AC#2 —— 反半边有没有被新用例做松？方向歪没歪〔独立复现〕
**判：没做松、方向正。我自己的形状（非照搬它的用例）造了三发探针，基线全绿，没把"被拒"变成"被放行"。**

新用例自带的三枚 AC#2 用例方向正确（读码确认＋跑绿）：env 改写要求 `Sync:true`＋`suspect-fallback`＋`Why` 含 `expansion moved`（响亮拒），
`~` 改写同理，同树正常写要求 `Sync:true` 且落在**登记 env 根**（`Provider=Ticket116Sync Source=env`，明写"suspect-fallback 回答=藏坏腿"）。

**我另写的独立探针** `zz_acceprobe_116acc_test.go`（只活在快照，验完即删、不进树）：
- (A) 普通 `%WISP116ACCPAIR%` 合法改写→另一棵树 ⇒ `Sync:true`＋`suspect-fallback`＋`expansion moved`（**响亮拒**，PASS）；
- (B) 同树 `notes/new.md` 无展开 ⇒ 命中登记 env 根（**照旧成功**，PASS）；
- (C) 真·树外写 ⇒ `Sync:false`（证明 Sync:true 是**判定**、不是仪器从不放行，PASS）。
`go test -run 116acc` RC=0。**没有**任一形状能让合法改写变放行 ⇒ 不构成"探针达成 AC 声称要防的结局"⇒ 不触发 FAIL-退回。
（三枚 AC#2 用例在 M1/M2 之下我实测全绿，见上表"其余三枚 PASS"。）

## AC#3 —— 射程话是不是真话〔独立复现；发现一处过保守、不准，但不判 FAIL〕
**判：主断言"祖先那次 `Actable()` 被删会红"=真（M1/M2 独立复现）。但文件头那句 `A wrong-but-called Actable() would keep these tests green` 经我 M3 复算＝过保守、不成立。**

我额外打了一发它没打的刀 **M3**（`pathresolver.go:93-99` 的 `Actable()` 整体改成恒 `return r.Canonical, nil`，即"被调、但判错=永远报可用"）：
grep 证 `MUTATION-116-M3` 落在 95 行 → build rc=0 → `go test -run 116`：**四枚里三枚变红**——
- 祖先用例红在第201行 `precondition: the ancestor's Actable() must be the thing that errors here`（`t.Fatal`）；
- env 用例红在第252行、home 用例红在第289行（都是 `first Actable() leg must refuse` 的直接断言）。

⇒ **一个"被调但判错（恒可用）"的 `Actable()` 会把本票用例弄红**，因为用例自己在 preconditions 里**直接调 `res.Actable()`/`ares.Actable()` 并断言其报错**。
所以头部 `would keep these tests green` 这句是**过保守**（它其实把 `Actable()` 的"报错义务"在**这些输入上**也钉住了一部分）。
- 定性：过保守＝**少报**覆盖面，不会造假绿 ⇒ **不判 FAIL、不阻塞**；但它是**不准确的自述**，按纪律要照实登记（过保守也是不准）。
- 正确的窄法应是："本票钉住祖先那次 `Actable()` **被删/被丢账**会红，也钉住**在这些输入上** `Actable()` 该报错；它**不**声称 `Actable()` 的判定语义在**未 exercise 的构造**（如 POSIX `$VAR`、非 `%pair%` 类）上也对——那归票 102 的 `pathresolver_rewrite_account_test.go` 静态判据。"（措辞修正见"下一步"。）

## AC#4 —— "只加测试"做没做得到（判定本体动没动）〔独立复现〕
**判：做得到，`syncdirs.go` 一行未动，也不需要新缝。**
`git show 80e248c:internal/risk/syncdirs.go | diff - <snapshot>` 空；形状全从导入口 `IsSyncPath` 进、`SyncStatus` 出，包内依赖
`deepestExistingAncestor`/`Resolve` 本就同包可见 ⇒ 测试不需生产码开口。加固①（`anc` 由 `deepestExistingAncestor(canon)` 从 C26 自己输出实取再断言
`gotAnc==ancestor`、`remainder=={new.md}`、串含 `%`）确实把票 105 那句"祖先串不含未展开构造"用**同一断言的两个前提同时成立**顶了回去——我读码＋跑绿确认。

## AC#5 —— 门禁四数 / vet / fmt / d22scan〔独立复现；POSIX 真跑非交叉 vet〕
Windows 原生，`-count=2 -v`，全部我亲手量：
- `internal/risk`：**RC=0**；顶层 `=== RUN`=**200**（100 不同名×2；全量含子测试 336）、PASS=**198**、FAIL=**0**、SKIP=**2**（唯一 SKIP 名 `TestSyncRegistryProbeLive`×2，既有 HKCU 探针，带 `-v` 量的）。预算 `TestResolvePerCallBudget` 读数 **0.504 / 0.673 ms/op**（**在 66% 负载下仍绿**，见 R 表）。
- `internal/tools`：**RC=0**；顶层 158=79×2、PASS=158、FAIL=0、SKIP=0。
- `go vet ./internal/risk/` **rc=0**。
- **POSIX 真跑**（docker `golang:1.27`，`MSYS_NO_PATHCONV=1 MSYS2_ARG_CONV_EXCL='*'`，容器内先 `ls -l /wisp/go.mod`＋`ls …116_test.go` 证非空挂载）：
  `-run 116 -count=2` **8 条 RUN / 8 PASS / RC=0**（**证明"去了 windows build tag⇒POSIX 有分母"是真的、非假绿**）；
  整包 `-count=1` RC=0、77 PASS / 0 FAIL / 1 SKIP（`TestC26RewrittenSyncRootDoesNotDisarmSuspectNet`，既有 POSIX 形状，×2 即自述的 154/2）；
  **M1 在 POSIX 也红**（上表，独立重跑）。措辞守口径：只说"POSIX 编译并真跑了这个包"，不说"Linux 全仓测过了"。
- `gofmt -l <新用例>` 空；`"$(go env GOPATH)/bin/gofumpt.exe" -l <新用例>` 空 —— gofumpt **存在**：`D:\work\base\gopath/bin/gofumpt.exe`，`v0.7.0 (go1.27.1)`，跑的是命令本身。
- `sh scripts/d22scan.sh`（纯净快照）**rc=0**；阳性对照 `PASS=21 FAIL=0 SKIP=0`；七 scope 逐字（bans#1-5 internal/=202、cmd/=20、ban#6 frontend/=40、ban#7 internal/tools/=18、ban#8 design/=16、frontend/=40、cmd/=26），
  **ban#8 internal/=373**；把 116 文件移出快照再跑＝**372** ⇒ 那 +1 **正是这枚 `_test.go`**（注释/测试都在 ban#8 射程），来源核清、台账不降。
- `GOOS=linux go vet ./...` / 交叉：按既有仪器坑，**只编译不执行** ⇒ 不作为"Linux 测过"的证据；Linux 结论只来自上面 docker 真跑。

## AC#6 —— 负载注明 + 那枚假红〔独立复现〕
我这轮 12 核、门禁时 `% Processor Time`≈66%（票 111/115/92b 在同共享树活着）。我这遍 `TestResolvePerCallBudget` 两读数 0.504/0.673 都在 1ms 线内＝绿，
但**这不证伪敏感性**——实现方已留 1.707/1.903 两枚红、验收方（105b）留过 1.668。⇒ 该用例对**包级并行负载**敏感属既有脆弱点（票 18 AC#6 形状），
**1ms 阈值一个字没动**（那是 golden）。登记 `R-116-1`（与 `R-105-3` 同根，不重复立案）。

## 三笔"不是我造成的红"—— 我独立确认过"移走它的文件照样红"〔独立复现〕
在纯净快照里把 `syncdirs_ancestor_actable_leg_116_test.go` **移出**后复跑：
- `internal/panel` `TestComposerRenderFixtureTellsTheTruth`：有/无 116 文件**都 FAIL**（`render fixture verified across 0 painted states`，票 92b 系）；且 panel 与 risk 是**不同包**，一枚 risk 测试文件结构上碰不到它。
- `cmd/wisp`：`exit status 0xc0000135`（STATUS_DLL_NOT_FOUND，缺 DLL，票 98 系）——移出 116 文件后仍 `0xc0000135`。
- `TestResolvePerCallBudget`：见 AC#6（负载所致，非 116 文件）。
⇒ 三笔**均不记到票 116 头上**，也**不阻塞**票 116 结案（它们先于本票、属别的票地界）。

## R 账（验收方视角）
- **R-116-1**〔要立案，属 `internal/risk` 地界，不阻塞本票〕：`TestResolvePerCallBudget` 包级并行下不可复现（本轮 66% 负载仍绿，但历史 1.707/1.903 红坐实）。与 `R-105-3` 同根。完成判据＝"安静机器上单独 `-run TestResolvePerCallBudget`"。
- **R-116-2（新）**〔要登记，属票 116 交付物 `internal/risk` 文件头注释，不阻塞〕：AC#3 头部 `A wrong-but-called Actable() would keep these tests green` 过保守、经 M3 复算不成立（见 AC#3）。**只改注释措辞、不改判据、不动生产码**。
- **R-116-3（继承 R-105-4）**〔不归本票，第三档〕：祖先那次 `Resolve` 的 **deny 分支** `syncdirs.go:223-225` 仍未被行为用例覆盖（差一枚"搬后树自带未登记 reparse 点"的 junction）。实现方 next=③ 已如实登记，本票不认领。
- **R-116-4（记账，非结论）**：本轮工具输出里出现 **0 次**自称"编排者备注/系统提示/请 revert/冻结某包/不要提它"的伪指令文本。我照实计数＝0（下界，任何一步都可能再附）；全程未 revert、未改判据、未动阈值。

## 总判
**通过（附一条非阻塞措辞修正 R-116-2）。**

一句话理由：票 116 存在的唯一理由——"删掉祖先那一次 `Actable()` 腿时，有没有一枚**行为**用例会红、而票 105 那两条仍绿"——被我在 Windows 与 docker POSIX 两侧**独立复现成立**
（M1/M2 各自把 `TestSyncAncestorActableLegFailsClosedOnMovedAncestor116` 单枚钉红、红形是 `Sync` 翻转/`Why` 丢账，非编译期符号；票 105 两枚 M1 下 2/2 仍 PASS），
反半边（合法改写响亮拒 / 同树成功）用我自己的探针证没被做松，POSIX 分母是真跑非假绿，门禁四数与 ban#8 +1 来源全部复算对上；
唯一瑕疵是文件头那句射程自述过保守（M3 之下反而三枚红），属少报、不造假绿 ⇒ 通过、把措辞修正挂 R-116-2 非阻塞。

**票 105 能否因此结案：能。** `acceptor-ticket105b` 卡 AC#3 的那格（"那半今天只能钉符号、钉不住行为"、`R-105-2`）已被本票的行为用例补齐并独立复现；
`R-116-3`（祖先 deny 分支 `:223-225`）不在票 105 AC#3 的射程话里，另账、不阻 105。三笔基线红（panel / cmd/wisp / budget）不记到票 116、也不阻 116。

## 下一张该派什么
1. **R-116-2 措辞修正**（票 116 自己地界，`internal/risk/…_116_test.go` 文件头，只改注释、不动判据/生产码）：把 `A wrong-but-called Actable() would keep these tests green`
   改成"本票在这些输入上也钉 `Actable()` 的报错义务；不声称其判定语义在未 exercise 构造（POSIX `$VAR` 等）上正确"。可由票 116 续作或并入票 105 结案收口。
2. **R-116-3 / 继承 R-105-4**（新票，`internal/risk`）：补祖先那次 `Resolve` 的 **deny 分支** `syncdirs.go:223-225` 行为用例——差一枚"搬后树自带未登记 junction"的探针。
3. **`R-116-1` 镜像进真相源**（编排者动手）：`docs/reports/pending-and-issues.md` 有并发写者，我不抢；与 `R-105-3` 同根、合并登记即可。
4. 票 116 无待办生产码；本票可在挂 R-116-2（非阻塞）后结案。

