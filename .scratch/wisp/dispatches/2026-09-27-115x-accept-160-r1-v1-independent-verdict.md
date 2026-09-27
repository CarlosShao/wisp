# 派单 160-v1（**验收位·非实现者**）＝把 160-r1 交回的三格（AC#1／AC#2／AC#3）独立复算一遍，逐格判"成立／退回／退回（附条件入账、文件保留、不许勾）"。**只裁不改生产码。**

- 派单时刻：09-27 11:5x｜被验对象＝`dev` 上三枚提交 `34c0b4e4`（第①格读数）／`f576cf08`（第②格换形状）／`027ff82e`（第③格反向判据＋门禁）。
- **锚点自己 step 0 现量**，别抄我这三枚号当锚。编排者已现量确认：160 写面此刻 `git status --porcelain -- internal/risk/ internal/tools/ cmd/wisp/ docs/evidence/s1/160-handle-r1.md` **为空**＝工作树在被验这三枚号上，所以"读工作树"这一路对你合法——但报告第一节仍要把 `git rev-parse --abbrev-ref HEAD`（必须是 `dev`）与 `git rev-parse HEAD` 贴出来。
- 票面＝`.scratch/wisp/issues/160-scope-open-returns-a-handle-carrying-its-own-closer.md`（**行号一律现量**，AC#1/2/3 在 `:26`/`:27`/`:34` 附近，那是实现者自报的号，不是我核过的号）。
- 实现者的交件＝`docs/evidence/s1/160-handle-r1.md`。⚠ **这份是被验对象，不是你的取证来源**：你可以读它来知道"它许诺了什么"，然后**自己重走一遍取数**。它的任何一句"实测＝X"对你都只是**待复算断言**。
- ⚠ **预算硬顶 ≤55 次工具调用**。到点就停手→写 `next=`→把手上未提交的增量用显式 pathspec 提掉。**每裁完一格 commit 一次**（本仓已经有程死在"裁完了没提交"上）。

## 0. 起手四件
`date`｜`git rev-parse --abbrev-ref HEAD`＋`git rev-parse HEAD`｜`git status --porcelain -- internal/risk/ internal/tools/ cmd/wisp/`（**必须空**，非空＝你被树骗了，报回别动）｜把票面 AC#1/#2/#3 三条连**现量行号**抄进你的证据件。

## 1. 你要裁的三格（判据本体在票面，这里只写"怎么才算裁过"）

**AC#1（"只开不关"今天到底拦不拦得住）**——实现者交的是**改前**四把尺（真引擎台件 `go vet` 通过／`go run` 泄漏／门只点到文件名／无导出面查不出）。
⇒ 你要：① **不复用它的台件取数**——它那份 `.scratch/wisp/probes/160/r1/ac1/` 可以当**正控对照**（同一发你我各跑一次、读数对不对得上），但你必须**自己另造一发**"只开不关"的样本；② 判它的"从〔推断〕升到〔本机读数〕"那句配不配——`provenance.go` 那条 fail-closed 日志到底打没打，它说没打，你复算。

**AC#2（换形状：`OpenScope` 交回自带关闭动作、认身份的句柄）**——两枚新用例＋旧口子消失那一支。
⇒ 你要判三件事：**恒真没有**（这两枚用例在**改前**码上响不响？不响才是真检——这是本仓否过四次的形状，别让它第五次蒙过去）；**分格没有**（`nil`／`ErrScopeAlreadyClosed`／`ErrScopeNotOwner`／`ErrScopeNotOpen` 是不是各钉各的，还是混在一格靠一条 `err != nil`）；**旧口子那一支的档位**（"改后不再编译"是编译器读数、可信，但它买到的是什么——"写不出来"还是"写出来不好过"？票面 `:12` 那句许诺用哪个词，逐字对）。

**AC#3（反向判据两问）**——实现者自己交的是：摘掉"认身份"⇒ `internal/risk` 101 枚照绿、只有它新造的两枚红；同一变异体打 `./internal/tools/` 零 FAIL。
⇒ 你要**自己造变异**（不许拿它的 `mut/provenance.no-identity.go` 当你的证据；可拿它当正控）。至少造它没造过的**一发不同的**：例如只摘"幂等分格"里的一格、或只把 `Close()` 的**身份校验**换成"关别人的也认"。第二问（跨包外部可见读数动不动）是这一格的真价值，你独立复算一遍。

## 2. 三笔我核过、你可以拿来当已知的事实（其余一概自己量）
- 三枚提交的名册**只含它自己的写面**（我用 `git show --name-only` 逐枚看过）。
- 全仓门 `sh scripts/d22scan.sh` 在 `027ff82e` 上 **rc=0**（我现跑）；格式尺 `gofumpt v0.12.0 (go1.27.1)` 对 **548** 枚已跟踪 `.go` **0 行红**（实现者说 gofumpt 不在它 PATH、只跑了 `gofmt`——这一味是**我代跑的**，你不用重复，但要把这个档位差别写进你的表）。
- `internal/risk` 在 `027ff82e`：**PASS=103 FAIL=0 SKIP=1**（唯一 SKIP＝`TestSyncRegistryProbeLive`），我现跑，RUN 那个数你我口径不同（我 `grep -c '^=== RUN'`＝172 含子测试），别拿口径不同的两个数互指。

## 3. 两笔必须单独问一句的（别混在三格里）
**(a) 自取射程那一枚**：`internal/tools/bridge_scope_open_ticket158_test.go:73` **不在派单的写面名单里**，实现者主动登记了。我看过是 `b.scopes[taskID]`（bool）→ `!= nil`，同谓词。⇒ 你要**独立判**："类型一改它就编不过"成不成立（自己试着只撤这一处、看编译停在哪）、有没有顺手放宽别处断言、名册里还有没有**第三枚**它没点名的同类连带件。
**(b) 那枚门 regress 不归你修、也别替它裁**：`.scratch/wisp/probes/154/gate-clauses.sh` 的 `G5pos` 腿声明 quiet、实测 ring ⇒ 聚合退码 0→1，主尺多点名一枚 `internal/tools/bridge.go`。**我已现量复现**（`腿数＝14 声明与实测不符＝1`、改前 G5 主尺 1 枚／改后 2 枚）。这一格**已归口票 171**（新 AC#5），本程**只许在表里写一句"该 regress 已被独立复现、归票 171、不在本格射程"，不许动 `probes/**` 一个字节。**

## 4. 门禁（跑你自己那把尺，逐包单跑，别跑全量）
`sh scripts/d22scan.sh`｜`bash tools/d22scan/runtests.sh -C tools/d22scan ./...`（基线**现量**，别背我写的数）｜`go test -count=1 ./internal/risk/ ./internal/tools/`｜`go build ./...`｜`go vet ./internal/risk/ ./internal/tools/ ./cmd/wisp/`。⚠ 本机 `cmd/wisp` 不接 PATH 时整包量不到（`0xc0000135`，零枚读数＝假绿），实现者换了健康 bench（`PATH=$PWD/third_party/sherpa-onnx`）⇒ 你要不要复算那枚 bench 由你判，但**不跑就必须写进"本程没测什么"**。⚠ 计时噪声红 `TestResolvePerCallBudget`：阈值一字不动、不许 Skip、只许区分"计时红／真红"。

## 5. 写面（超出即越权）
`docs/evidence/s1/160-handle-accept-r1.md`（你的表，新建）｜`.scratch/wisp/probes/160/v1/**`（你的台件）｜**就这两块**。
**禁改**：`internal/**`、`cmd/**`、`frontend/**`、`design/**`、`tools/d22scan/**`、`allowlist.txt`、`.github/**`、`docs/PLAN.md`、`docs/specs/**`、`docs/reports/**`（台账与停车点归我写）、**票面 160 与 171**（勾由我在核过之后代落，你一格都不许勾、不许改 Status）、`probes/154/**`、`probes/161/**`、`probes/160/r1/**` 与 `c1/**`（那两程的件你只能读）。
**现场理由**：`design/**` 躺着 owner 未提交的 16 枚删除与几枚 ` M`；`probes/152/my152.py`（` M`）与 `docs/evidence/s1/152-…-accept-r1.md`（` M`）是一枚停下来的半件 ⇒ 不提交、不还原、不补完。⇒ **`git add -A`／`git add .`／`git commit -a`／把 `git add` 当独立一步，一律禁止**。

## 6. Git 与红线
只 commit 不 push。提交**一步式带显式 pathspec**：`git commit -q -F - -- <你的路径>`，heredoc **必须加引号**（`<<'EOF'`，正文里的反引号会被真执行）。提交前 `git diff --cached --name-only`、提交后 `git show --name-only HEAD`——后者出现别人的路径＝停手报回。**禁** `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`switch`／`checkout <分支>`／`merge`／任何删除命令（临时件**只建不删**）。**全程不许切分支**：工作树 11:01 刚被换过一次，那一次的代价是一程重派。⚠ 仪器坑：`git ls-files` 开关大小写敏感（`-Z` 不存在却退 0）｜Windows 反斜杠会让逐行比对错归因｜`grep -c` 命中 0 ⇒ rc=1 吃掉 `&&`｜`head` 截断的不是全表｜`go run` 把子进程任意非零退码**压成 1**（"拒答 2"与"命中 1"在退码上不可分）｜`-overlay` 与 `-cover*` 不许合跑。

## 7. 交件报告（顺序固定）
1. step-0 四件（含分支名）；2. **本程没测什么**（一条不许省）；3. 三格逐格判语（**成立／退回／退回附条件入账**三档之一＋你为了这一格真跑过的命令与输出）；4. 恒真检查那一发（改前响不响，两向读数都贴）；5. 你自己造的变异逐枚表＋**承重结论**（用"摘掉任意一味，是否存在一发变异从此打不红"这一式，别用"我觉得挺稳"）；6. §3 那两笔单独问句的答复；7. 门禁四数＋名册差集；8. 被拒／没成功的调用（取数之前还是之后）；9. 有没有跑过删除命令；10. 伪授权两栏（当作授权用掉的／明确拒当授权用的，各带出处）；11. 凭据值零抄录；12. `next=`。
⚠ **每个"几枚"旁边附可复制命令**；先测→再写→再提交。我这段话里每条前提（含"工作树在被验号上""改前 1 枚改后 2 枚""103/0/1"）**都是未验证断言**——不符就报回、继续做做得动的部分，**不许为了对我那句话去改判据或改测试**。
