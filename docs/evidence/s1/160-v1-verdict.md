# 票 160 · 终裁腿 `160-v1` —— `AC#4`／`AC#5` 两格判语

- 身份：**非实现者裁腿**（裁决者≠实现者）。本件**不落 `.go`、不翻任何框、不改票面/台账、不 push**。
- 被裁：`160-r1` 三笔 —— `c4d65e4e`（`internal/risk/provenance.go` **+25/−1**；新增 `provenance_160_boundary_test.go` **164 行**）／`f895fe5c`＋`02b3882c`（`docs/evidence/s1/160-r1-impl.md`，145+4−2＝**147 行**）。
- 锚：HEAD `d0194273`；父提交 `ec8786e1`（`c4d65e4e^`）。★`git log c4d65e4e..HEAD -- internal/risk` **零命中** ⇒ HEAD 的 `internal/risk` 与 `c4d65e4e` 同字节。
- 写面：本件＋`D:\tmp\wisp\20261022\160v1\**`（导出树/日志，仓外）。每节的尺与射程写在节内。
- **一句话**：`AC#4` **能翻**（两条"治不到"都写成可核句子、各有读数、各有一枚承重用例）；`AC#5` **按字面 ⛔ 翻**（票级 6 笔名册被 09-27 那轮越过 3 枚）、**按本腿 3 笔名册干净**，缺的只是一条追加改判。

## ① 十项逐项判语＋读数

**1) 票面前提新鲜度**（尺＝`git grep -n … HEAD -- <file>`，量 HEAD blob，⛔ 工作树）——它报的 5 处**逐条复核全对**：

| 票面/派单写的 | 今天 HEAD 现读 | 判 |
|---|---|---|
| `provenance.go:603-612` | `Detector` `:868-870`、`boundDetector` `:886-889`、`TaintHit` `:891-897` | ★对（`34c0b4e4^` 时 `boundDetector` 在 `:609`，本票自己推后） |
| `allowdirs.go:157-163` | `func statOwnWrite` `:157`、三赋值 `:162`、闭 `:163` | ★对（唯一一枚 `allowdirs.go`，`git ls-files` 现量 1） |
| `assessor.go:228`／`bridge.go:610` | `Assess` `:230`；判定/盖章 `bridge.go:297`＋`:299` | ★对（`:610` 今天是 `out.ErrorClass = ClassCancelled`，不是判定） |
| `loop.go:360-366` | `:373` `revokeAdmission := func(){}`、`:379` `defer` | ★对（`:360-366` 今天是 journal 注释） |
| `run.go:557-565` | `admitTask` `:1025-1039`，`:1037 rt.bridge.CloseTask` | ★对（`:557-565` 今天是工具注册循环） |

★**它还漏报了 4 处腐坏**（不阻塞，但票面引用会继续误导）：Status 行的 `bridge.go:642`（今天 `:864` 开／`:940` 关）、`panel_assets.go:232`（今天 `:243`）、`provenance.go:55-61`（今天 `:58-65`）、`gate.go:150`（今天 `:184`，它在 §2-E 引对了但没进腐坏表）。票面 `AC#1/#2/#3` 已勾、`#4/#5` 未勾——复核成立。

**2) `AC#4` 判点一（"它治不到什么"可核吗）**——三条分开判，全成立：
- **从不动的泄漏**：我自造探针（`package risk_test`，只走导出 API，生产形状 `ProvOptions{NoProbe:true}`＋`_ = OpenScope`）⇒ `taints=1 hit_after_gc=true logged_total=1 logged_about_scope=0`（那 1 行是构造期的 P12 同步目录告警，与泄漏无关）。★口径补正：它那枚用例的 `logged_lines=0` 是在 `NewProvenance` **之后**才装 sink（`provenance_160_boundary_test.go:50-51`）⇒ 量的是"构造后没人说话"；**泄漏专属行数＝0 成立**。且 `risk.Logf` 在**生产码零赋值**（尺：`git grep -n 'risk\.Logf\|Logf = ' -- cmd internal | grep -v _test.go` ⇒ **0 行**；别的包各自的 `Options.Logf` 字段不是它）⇒ sink 恒 nil，引擎的话今天根本送不出去——比它写的更强。
- **钉子承重**：我造 `M-LEAK-SPEAKS`（`OpenScope` 里加一行 `logf`）⇒ `logged_lines=1`、`FAIL …UntouchedLeakIsSilent` @ `:72`。⇒ "泄漏在、没人说话"**被一枚用例钉住**，⛔ 恒真装饰品。
- **借用那一形**：我造 `M-BORROW-CLOSE`（`TaintHit` 里补 `_ = d.p.(*Provenance).OpenScope(d.scope).Close()`）⇒ `hit#1 true taints=0`／`hit#2 false taints=0`／`:94` Fatal，**并连带** `TestSyncWriteNegative` @ `provenance_test.go:277`。`TestSyncWriteNegative` 是**票 19（`73b60571`）就有的用例**、在 HEAD 与父树上基线**全绿**，只在变异下红 ⇒ **是"既有用例"、⛔ 是"既有红"，也不算新增红**。

**3) `AC#4` 判点二（收窄是否只改形状）**：`git show c4d65e4e --numstat` 现读 **25/1**（`provenance.go`）＋**164/0**（新测试），`git diff --stat c4d65e4e^ c4d65e4e` ＝ **1 file, +25/−1**，⛔ 1115 行重写。全 diff 只有一段：`Detector` 注释＋`detectorSeam` 接口＋`var _ detectorSeam = (*Provenance)(nil)`＋`boundDetector.p` 静态类型。`Inspect`/`Assess`/`Mark`/`TaintHit` 函数体**一字未进 diff**；`assessor.go`、`bridge.go` **不在提交里**。⇒ **只改形状、零判定语义**（判法见 ⑤）。

**4) 反向判据**：我造 `M-HOP-WIDE`（`p detectorSeam` → `p *Provenance`）⇒ **恰一枚红**：`provenance_160_boundary_test.go:156 READING borrow.face: type=*risk.Provenance kind=ptr methods=[CheckText Detector Inspect IsSyncPath Mark MarkWithHostPath OpenScope ScopeTaints SyncDetectionComplete SyncRoots]`／`:158` Fatal。**"未修码上也响"复核成立**：父树＋本笔测试文件 ⇒ `:158` 同一行红，另三枚同读数绿（`taints=1`／`hit#1/#2=true`／`ownerless survives`）。

**5) `AC#5` 契约轴**：本腿 3 笔名册＝`internal/risk/provenance.go`＋`provenance_160_boundary_test.go`＋`docs/evidence/s1/160-r1-impl.md`（逐枚 `--no-walk` 现量；`.go`＝2）；负向名单（`docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`／`tools/d22scan/**`／`frontend/**`／`design/**`）**零命中**。票级 6 笔名册另见 ③。

**6) "新增红 0"独立核**：父树 `-v`＝PASS **140**／FAIL **1**／SKIP **1**／RUN（顶层）**142**（含子测试 245/1/1/247）；HEAD `-v`＝PASS **144**／FAIL **1**／SKIP **1**／RUN（顶层）**146**（含子测试 249/1/1/251）。红名两树都只有 `TestC26…`；`TestResolvePerCallBudget` 在**父树**也红过（`1.779 ms/op`），HEAD 红过（`1.715`／`1.311 ms/op`）⇒ 计时噪声，**⛔ 新红**。★它报的"PASS=144／RUN=251"是**两把尺混用**（144 是顶层 PASS、251 是含子测试 RUN；单尺应为 144/1/1/146 或 249/1/1/251）——结论不变，读数须改口径。

**7) 门禁复跑**：`go build ./...` **rc=0**（零输出）／`go vet ./internal/risk/` **rc=0**（零输出）／`go test -count=1 ./internal/risk/` **rc=1**（`--- FAIL: TestResolvePerCallBudget 1.715ms`＋`--- FAIL: TestC26…`，末行 `FAIL`）／`sh scripts/d22scan.sh` **rc=0**／`sh scripts/check-path-length-budget.sh` **rc=0**。格式甲乙见 ⑦。

**8) 写面纪律**：三笔逐笔 `git show --no-walk --format='' --name-only` ⇒ `c4d65e4e` 2 枚（`.go`=2，越界 0）／`f895fe5c` 1 枚 md／`02b3882c` 1 枚 md；**越界 0**。**正控**同尺打 `394cc212`（票 223，非本票）⇒ 1 枚、`.go`=1 ⇒ 这把尺**看得见 `.go`**。

**9) 三项欠账**：判语见 ⑩（①`DisposalScope` 支 ②`close_err=%v` ③无主登记清理路径）。

**10) CRLF 归一**：判语见 ⑩（**不算越界**，但是本机检出状态的一次局部改动，应登记）。

## ② `AC#4`：**能翻**

判据逐字对照：**"从不动的泄漏它照样漏"**——`§3` 读数＋我独立探针＋`M-LEAK-SPEAKS` 承重；**"借用那一形不该由它治、在那儿补关会弄坏它"**——`§4` 两读数＋`M-BORROW-CLOSE` 打红两枚（一枚本笔、一枚票 19 既有）；**"与票 158 两格、不许并成一格"**——它的证据件头部（`160-r1-impl.md:3-5`）逐字写了"⛔ 不是票 158 今天响不响那一格的重判"，且 `§3` 明写"点名的是 G5 普查尺，不是本引擎"。
09-27 验收程点名要在本格写清的四件——`DisposalScope` 支／从不动的泄漏／`close_err=%v` 无人钉／无主登记——在它 `§2` 的 A–E 五行逐条在案（⛔ 省任何一条）。
**缺什么**：⛔ 不缺承重读数。两条口径补正（⛔ 退回理由）：①`logged_lines=0` 是"构造后"读数（sink 晚装），泄漏专属行数 0 仍成立；②`TestSyncWriteNegative` 是既有**用例**、不是既有**红**，别写反。另：三项欠账的"归宿"里，`DisposalScope` 支仍**没有票号**（见 ⑩①），建议编排者在结票时点名，⛔ 让它随 `-done` 蒸发。

## ③ `AC#5`：**按字面 ⛔ 翻；按本腿名册能翻**

- **本腿 3 笔（`--no-walk` 集合，逐枚现量）**：`internal/risk/provenance.go`（`c4d65e4e`）／`internal/risk/provenance_160_boundary_test.go`（`c4d65e4e`）／`docs/evidence/s1/160-r1-impl.md`（`f895fe5c`＋`02b3882c`）⇒ **全在字面白名单内**，负向名单**零字节**。
- **票级 6 笔（`34c0b4e4`／`f576cf08`／`027ff82e`＋本腿 3 笔）**：另带 `cmd/wisp/panel_assets.go`、`internal/tools/bridge.go`、`internal/tools/bridge_scope_open_ticket158_test.go` 三枚**在白名单外**（前两枚＝`A322` 批准的"两枚文件必要改动"，第三枚＝编译器必需，09-27 已追认）；负向名单仍**零命中**（`docs/PLAN.md`／specs／thresholds／golden／allowlist／d22scan／frontend／design 一枚未动）。
- **缺什么**：只缺一条**追加改判**（把 `A322` 批准的两枚生产文件＋编译器必需的 tools 测试写进 `AC#5`，原话不抹）。⛔ 缺任何代码或读数。**在那条落笔前按字面勾＝假勾**。

## ④ 我自造的变异与反向判据（逐字读数）

- **`M-BORROW-CLOSE`**（补关；`d.p` 需先做 `.(*Provenance)` 断言才写得出来——正是收窄买到的那一跳）：
  `READING borrow.hit#1: hit=true src="fs.read /legacy/secret" taints=0`／`READING borrow.hit#2: hit=false src="" taints=0`／`hit #2: the borrowed detector must stay live across calls (a close at the borrow point is what kills it)`／`--- FAIL: TestScopeBoundary160BorrowedDetectorOwnsNothing`／`fs.write without a verifiable path must fail-closed to the scan`／`--- FAIL: TestSyncWriteNegative`／`rc=1`。
- **`M-HOP-WIDE`**（摘掉那一跳）：见 ①-4 的两行逐字＋`--- FAIL: TestScopeBoundary160BorrowFaceCarriesNoCloser`；同包另三枚绿。
- **`M-LEAK-SPEAKS`**（让泄漏说话）：`READING leak.silent: … logged_lines=1`／`the leak must be silent here (the G5 census names it, not this engine)`／`--- FAIL: TestScopeBoundary160UntouchedLeakIsSilent`。
- **"未修码上也响"**（父树 `ec8786e1` ＋本笔测试文件）：`:156` 同读数（`kind=ptr`）、`:158` 同红行；另三枚绿。⇒ 这一跳**今天不在**，⛔ 自证。
- **探针件**（仓外、⛔ 入库）：`D:\tmp\wisp\20261022\160v1\tree-head\internal\risk\zz_160v1_probe_test.go`（`package risk_test`，只走导出 API）；变异树 `tree-mut-wide`／`tree-mut-borrowclose`／`tree-mut-leakspeaks` 与日志同在 `D:\tmp\wisp\20261022\160v1\**`。
- 还原核对：所有变异只打在仓外导出树；仓内 `git status --porcelain -- internal/risk` **零命中**。

## ⑤ "只改形状⛔改语义"：我的判法＋结论

判法四层，缺一不判：**(a) 提交面**——只有 `provenance.go`＋新测试两枚，`assessor.go`/`bridge.go` 不在提交里；**(b) 行面**——`--numstat` 25/1，唯一 hunk 是 `Detector` 注释＋接口声明＋`var _`＋字段静态类型；`Inspect`/`Assess`/`Mark`/`TaintHit` 函数体零改动；**(c) 调用面**——`TaintHit` 仍调 `d.p.Inspect(d.scope, "", params)`，接收者仍是同一枚 `*Provenance`（接口只是收窄了"还能调到哪些方法"），`Detector` 仍返回 `boundDetector{p, scopeID}`；**(d) 行为面**——整包红名与父树**逐名相同**、PASS 增量 **+4 恰是四枚新测试**（顶层 140→144、RUN 142→146），无任何既有用例翻面。
**结论：只改形状，判定语义零碰。** 编译期还有 `var _ detectorSeam = (*Provenance)(nil)` 把 `Inspect` 签名钉住；`Detector` 的返回值类型 `TaintDetector` 未变。

## ⑥ "新增红 0"：我的独立核（另起尺，⛔ 采信编排者）

- 父树（`c4d65e4e^`，仓外导出树）`go test -count=1 -v ./internal/risk/` ⇒ `rc=1`；顶层 **PASS 140／FAIL 1／SKIP 1／RUN 142**；唯一红 `TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（`internal/panel/composer.go:127`、`display_scrub_216.go:26/37/218` 未读 rewrite account）。
- HEAD `-v` 同尺 ⇒ `rc=1`；顶层 **PASS 144／FAIL 1／SKIP 1／RUN 146**；红名同上。
- **计时红不是新的**：父树 `-run TestResolvePerCallBudget` 三次里红一次（`1778983 ns/op = 1.779 ms/op`），HEAD 红过 `1.715`／`1.311 ms/op`、也绿过 ⇒ 同一棵树两次跑可红可绿 ⇒ 与 `c4d65e4e` 无因果（该提交未碰 `pathresolver`/`syncdirs`）。
- ⇒ **新增红＝0**。★口径：它报的 144/1/1/251 是混尺，单尺见上。

## ⑦ 门禁 rc 与末行逐字（本腿现跑、串行；`unset GOTMPDIR TMPDIR TEMP TMP`、`GOFLAGS=`）

| 门 | rc | 末行逐字 |
|---|---|---|
| `go build ./...` | **0** | （零输出） |
| `go vet ./internal/risk/` | **0** | （零输出） |
| `go test -count=1 ./internal/risk/` | **1** | `FAIL`（前一行 `FAIL	github.com/CarlosShao/wisp/internal/risk	15.542s`） |
| `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations; … ban #8 internal/=619 … cmd/=148 …` |
| `sh scripts/check-path-length-budget.sh` | **0** | `…unit cross-check: the byte view and the character view select the SAME 61 paths on this tree…` |

**格式甲乙两尺**（裁定件 `docs/evidence/s1/pipe-format-gate-reads-tracked-blob-not-worktree.md`）：
- **甲（blob 树，判据）**：`git archive c4d65e4e` 到仓外跑 `gofumpt.exe -l .` ⇒ **20 行，全部 `.scratch/`，非 `.scratch` = 0**；`provenance.go` 与新测试**一枚都不在红行里**。HEAD 树同尺 21 行＝20＋我自己的探针文件（仓外产物，⛔ 计入）。
- **乙（工作树，只作诊断）**：**49 行＝46 行 `.scratch` ＋3 行非 `.scratch`**（`cmd/wisp/models.go` CR=**389**／`internal/agent/approval/pending_read.go` CR=**131**／`internal/agent/tools.go` CR=**225**，三枚 `git ls-files --eol` 现读全 `i/lf w/crlf`）⇒ **本机 `core.autocrlf=true` 的假红**（我另做对照：把入库的 `provenance.go` 原样转 CRLF ⇒ `gofumpt -l` 立刻列出它 ⇒ CRLF 单独就足以致红）。它报的 54＝51+3 里 `.scratch` 那半已随别的腿漂移，**非 `.scratch` 三枚逐名相同**。

## ⑧ 红名册（逐名）

1. `TestC26RewriteAccountIsConsumedAtEverySecurityLeg` —— HEAD 与父树**都红**（既有红，归 `internal/panel` 消费者）。
2. `TestResolvePerCallBudget` —— 计时红，父树 `1.779 ms/op`、HEAD `1.715`／`1.311 ms/op`，**两树都可红可绿**（⛔ 新红、阈值一字未动）。
3. 变异红（⛔ 入库码的红）：`TestScopeBoundary160BorrowFaceCarriesNoCloser` @ `:158`（M-HOP-WIDE）／`TestScopeBoundary160BorrowedDetectorOwnsNothing` @ `:94`＋`TestSyncWriteNegative` @ `provenance_test.go:277`（M-BORROW-CLOSE）／`TestScopeBoundary160UntouchedLeakIsSilent` @ `:72`（M-LEAK-SPEAKS）。

## ⑨ 提交号＋逐枚 `.go`／越界＋正控

| 提交 | 名册（`--no-walk` 现量） | `.go` | 越界 |
|---|---|---|---|
| `c4d65e4e` | `internal/risk/provenance.go`、`internal/risk/provenance_160_boundary_test.go` | **2** | **0** |
| `f895fe5c` | `docs/evidence/s1/160-r1-impl.md` | 0 | 0 |
| `02b3882c` | `docs/evidence/s1/160-r1-impl.md` | 0 | 0 |

**正控**：同尺打 `394cc212`（票 223 落地腿，非本票）⇒ 1 枚、`.go`=1 ⇒ 尺能看见 `.go`（另：本腿三笔的 `--numstat` 现读 25/1＋164/0＋145/0＋4/2，与它自报一致）。

## ⑩ 三项欠账与 CRLF 归一：判语

- **① `AC#2` 的"挂到任务 `DisposalScope`"支未做**：**不算本票欠账**。票面 09-27 10:5x 的 `>` 更正已把它降成"(2) 不承诺"（要动 agent/tools/cmd 三层＝超 `A322` 批准单位），且码里 `DEFERRED(C25-loop-wiring)`（`provenance.go:58-65`）仍逐字点名这半未闭。**但"归宿"没有票号**——建议编排者结票时点名（⛔ 让它蒸发）。附带一笔既有账：`DEFERRED(C25-loop-wiring)` 现量 3 枚而 `SPEC-12 §5` 登记表 0 条（与本笔无关，⛔ 算它头上）。
- **② `bridge.go:942-943` 的 `close_err=%v` 无人钉**：**不算本票欠账**。`internal/tools` 在本腿字面写面外，`AC#4` 只要求"明写"（它写在 `§2-D`）；⛔ 为它越界改测试。真要做，归 `internal/tools` 侧票。
- **③ 无主登记（`provenance.go:628`）没有新增清理路径**：**不算欠账，且"不改"是对的**。`AC#4` 的活是把它写成句子（`§2-C`＋一枚用例钉住"外来句柄关不掉、只能重开同 id 再关"），加清理路径＝改判定语义，撞 `AC#2` 冻结。
- **CRLF 归一（`provenance.go` 工作树 `i/lf w/crlf` → LF）**：**不算越界**。理由三条现量：①对 git **零差异**（`git status --porcelain -- internal/risk` 零命中；提交 diff 仍 25/1，父/子 blob 的 CR 都是 **0**）；②文件在它自己的写面内，且 `.gitattributes` 就是 `*.go text eol=lf`；③裁定件 §4.4 明写"⛔ 修那枚 CRLF 副本当成本票的活"——反过来说它是**本机检出状态**、⛔ 仓库内容。**但**它确实局部改了工作树字节，应登记为"本机 `autocrlf=true` 的一次局部归一"；只治自己写面内那一枚、另三枚没碰，**边界正确**（碰了才越界）。要根治是统一本机 `autocrlf`／全树 renormalize，属另一件事。
