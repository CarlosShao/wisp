# 票 160 · 落地腿 `160-r1` · AC#4 交件 —— 它治得到什么／治不到什么

**这一格是什么**：`AC#4` 要求**明写这道机器治不到什么**（⛔ 含糊的"部分覆盖"）。它不是再补一道守卫，
而是把**边界**写成可核的句子，且每条"治不到"**各带一发读数**。
**它不是什么**：⛔ 不是票 158"今天响不响"那一格的重判（那格是 G5/G6/G7 的事），**两格不许并成一格结掉**。
**锚**：写前 `ce84d53d`／写后 `c4d65e4e`。**写面**：`internal/risk/provenance.go`＋其测试＋本件，共两枚 `.go`、一枚 `.md`。

---

## §1 前提新鲜度（票面／派单里那两处行号，今天逐条现跑）

| 票面／派单写的 | 今天现读 | 判 |
|---|---|---|
| `provenance.go:603-612` ＝ `Detector` 返回绑定过的借用对象 | `Detector` 今天在 **`:868-870`**，`boundDetector` 在 **`:886-889`**（`TaintHit` `:891-897`） | ★**行号过期**（同尺可核：`34c0b4e4` 之前 `boundDetector` 在 `:609`，今天 `:886`——票 160 自己那次改动推后的）。**句子本身成立**，读数见 §4 |
| `allowdirs.go:157-163` ＝ `statOwnWrite` 只认领自己写出来的内容 | `internal/config/allowdirs.go:157` func／`:162` 只写 `seenMtime, seenSize, seenValid`（全仓唯一一枚 `allowdirs.go`） | ★**行号对**，射程也对（它是**类比**：只认自己造的，不是本票的守卫） |
| 票面 `[ ]` 只剩 `AC#4`／`AC#5` | `AC#1/#2/#3` 已勾、`#4/#5` 未勾 | 对 |
| `AC#2` 引的 `assessor.go:228`（`Assess`）与 `bridge.go:610`（宿主侧盖章） | `func (a *RiskAssessor) Assess` 今天在 **`assessor.go:230`**；真正的判定与盖章在 **`bridge.go:297-299`**（`verdict := b.assessorFor(req.TaskID).Assess(...)`／`dec.Level = verdict.Level`）。`:610` 今天落在 D31 取消那段里 | ★**两处都过期**；"判定语义不许变"这句话本身成立（本笔一字节未动） |

★**尺的射程**（落笔前四问）：①量的是**主检出的 `internal/risk` 与三处生产调用点**，⛔ 导出树；
②输出是 `file:line`＋测试名＋原始行；③本件里"代码"与"注释"分列（`§2` 两栏）；④每条"治不到"都能点到**唯一一枚**新测试或既有测试。

## §2 边界表（每行一个 `file:line`；"治得到"＝有一枚可读的返回值/登记，⛔ 只有日志）

| # | 它**治得到**什么 | 出处 |
|---|---|---|
| 1 | **关别人的作用域**：句柄不再拥有该 id 时 `Close` 返回 `ErrScopeNotOwner`，且**一个字节都不动**（污点不丢、门照响） | `provenance.go:444-446` |
| 2 | **连关两次**：第二次返回 `ErrScopeAlreadyClosed`（幂等、不 double-drop） | `provenance.go:436-438` |
| 3 | **空登记上的误关**：登记在、污点没有，也照样拒绝并留着登记（fail-closed） | `provenance.go:439-446` |
| 4 | **"读 id ≠ 有权关"**：`ID()` 只读；**按 id 关的函数在 `internal/`、`cmd/` 零命中**（`CloseScope` 字样 5 处**全是注释/测试注释**：`provenance.go:365`／`:689`／本笔测试 `:114`／`provenance_test.go:312`／`:339`） | `provenance.go:408-416` |
| 5 | **"一本账"**：登记与它唯一的关者同处一个值、同一把锁（票 158 那本并行账塌掉了） | `provenance.go:258-270` |
| 6 | **生产关点**：任务边界用**发出去的那枚**关，并把拒绝写进审计行 | `bridge.go:864` 开／`bridge.go:940` 关／`:942-943` 审计 |

| # | 它**治不到**什么（⛔ 不许含糊） | 出处 |
|---|---|---|
| A | **从不动的泄漏照样漏**：句柄被丢掉、没有 `Close` ⇒ 登记＋污点留到进程结束，**一行日志都没有**（没有 goroutine／定时器／finalizer；`D38`＋`ban #1`） | `provenance.go:460-470` 发句柄／`:429-451` 唯一的销毁路径；生产实例 `cmd/wisp/panel_assets.go:243`（`_ = prov.OpenScope`）；测试里也是这个形状，`provenance_test.go:24` |
| B | **借用那一形不该由它治**：`Detector` 交回的是**借用**。在那儿补关 ⇒ 借用方把任务自己的污点删掉，**第二次调用就静默不响**（fail-open），且所有者的 `Close` 不再是所有者的 | `provenance.go:868-870`／`:886-889`／`:891-897`；生产借用面 `cmd/wisp/panel_assets.go:247` → `:66` 进 assessor |
| C | **无主登记改后谁也关不掉**：`Mark` 打到没人开过的 id 上会落一枚 `owner == nil` 的登记；`Close` 只删自己拥有的，按 id 关的口子又被票 160 拿掉 ⇒ **它没有关的路径**，唯一出路是再 `OpenScope` 同一个 id 接下所有权再关 | `provenance.go:626-629`（`:628`）／`:264-266`（owner nil 的说明）／`:444-446` |
| D | **审计行里的 `close_err=%v` 无人钉**：`internal/`、`cmd/` 的 `_test.go` 对它**零命中**；本票写面到不了 `internal/tools`，所以只登记不修 | `bridge.go:942-943` |
| E | **`DisposalScope` 那一支仍未做**（`AC#2` 的"(2) 不承诺"）：全仓 `DisposalScope` 走到 tools／agent／cmd 只剩**一句注释**（`git grep -nw` 恰 1 命中）；关腿真实持有者今天在 `loop.go:373-379` 的 `revokeAdmission` → `run.go:1025-1039` 的 `admitTask` 闭包（`:1037 rt.bridge.CloseTask`） | `provenance.go:58-65`（`DEFERRED(C25-loop-wiring)`）／`:426-428`；`internal/agent/approval/gate.go:184` 是那句注释 |

## §3 读数一：**从不动的泄漏**（生产形状；未修码与修后**逐字相同**）

未修码（`-overlay` 换成 `ce84d53d` 的 `provenance.go`，工作树零写入；`sha256 092590aa…`）与修后各跑一遍：

```
=== RUN   TestScopeBoundary160UntouchedLeakIsSilent
    provenance_160_boundary_test.go:61: READING leak.silent: registered=true owner_was_issued=true taints=1 logged_lines=0
--- PASS: TestScopeBoundary160UntouchedLeakIsSilent (0.02s)      ← 未修码；修后同一行只差耗时
```

★读法：`taints=1` ＋ `logged_lines=0` ＝ 泄漏**在**、**没人说话**。这不是缺陷，是**边界**：它治的是
"被不拥有的人关掉"，⛔ "从来没人关"。同一枚测试还钉住"漏出来的作用域**照样拦**"（`Detector` 照样 `hit=true`），
⛔ 把泄漏说成"反正没影响"。

## §4 读数二：**借用那一形**——`Detector(scopeID)` 是借用，补关会弄坏它

**(a) 未修码上，借用是有生命的、且不花掉作用域**（同一发在修后逐字相同）：

```
    provenance_160_boundary_test.go:91: READING borrow.hit#1: hit=true src="fs.read /legacy/secret" taints=1
    provenance_160_boundary_test.go:91: READING borrow.hit#2: hit=true src="fs.read /legacy/secret" taints=1
--- PASS: TestScopeBoundary160BorrowedDetectorOwnsNothing (0.01s)
```

两次同参命中都响，`taints` 仍是 1，所有者的 `Close()` 返回 `nil`（借用没花掉它）。

**(b) 补关变异 `M-BORROW-CLOSE`**（给 `TaintHit` 真的补一句关；`-overlay` 打在盘上，工作树零写入）：

```
    provenance_160_boundary_test.go:91: READING borrow.hit#1: hit=true src="fs.read /legacy/secret" taints=0
    provenance_160_boundary_test.go:91: READING borrow.hit#2: hit=false src="" taints=0
    provenance_160_boundary_test.go:94: hit #2: the borrowed detector must stay live across calls (a close at the borrow point is what kills it)
--- FAIL: TestScopeBoundary160BorrowedDetectorOwnsNothing (0.02s)
--- FAIL: TestSyncWriteNegative (0.01s)          ← 本票之前就存在的用例，同一条变异把它也打红
rc=1
```

★两枚**独立**见证 ⇒ "在那儿补关会弄坏它"不是纸上话：一枚是本笔新钉的，一枚是**别人早写的**。
★另注：收窄之后，那句补关**还得先做一次 `d.p.(*Provenance)` 类型断言**才写得出来——这正是 §5 那一跳买到的"⛔ 顺手改"。

## §5 反向判据：把新加的那一跳**摘掉** ⇒ 转红并逐名报行

那一跳＝`boundDetector.p` 的静态类型从 `*Provenance` 收窄成 `detectorSeam`（**只有 `Inspect`**，见 `provenance.go:872-880`／`:884`／`:886-889`）。
为什么只能是**形状**：`AC#2` 冻住了判定语义（`assessor.go:230`／`bridge.go:297-299`），`AC#4` 新加的任何一跳都不许改一个判定。

- **变异红**（`p detectorSeam` → `p *Provenance`，`rc=1`）：

```
    provenance_160_boundary_test.go:156: READING borrow.face: type=*risk.Provenance kind=ptr methods=[CheckText Detector Inspect IsSyncPath Mark MarkWithHostPath OpenScope ScopeTaints SyncDetectionComplete SyncRoots]
    provenance_160_boundary_test.go:158: the borrowed detector holds *risk.Provenance itself, so a close at the borrow point becomes writable again: Detector must hand out detectorSeam
--- FAIL: TestScopeBoundary160BorrowFaceCarriesNoCloser (0.00s)
```

- **还原绿**（`cp` 冻结件回来，`rc=0`）：`--- PASS` ×4（`…LeakIsSilent`／`…BorrowedDetectorOwnsNothing`／`…OwnerlessRegistrationHasNoCloser`／`…BorrowFaceCarriesNoCloser`）
- **原样绿**（变异之前那一跑，`rc=0`）：同四枚 `--- PASS`
- **未修码上它也响**（`-overlay` 打 `ce84d53d` 的 `provenance.go`，`rc=1`）：上面那两行 `READING borrow.face` ＋ `FAIL` 逐字相同 ⇒ 这一跳**今天不在**，不是自证。

**还原核 `sha256`（三处逐字相同）**：`cbd7093382ff893030aa3a656832b71922529f2c3798d69a6d40f64e5ecc5991`

| 处 | 对象 |
|---|---|
| a | `/d/tmp/…/probe/provenance.hop.go`（变异前冻结的那枚） |
| b | 还原后的工作树 `internal/risk/provenance.go` |
| c | `git show c4d65e4e:internal/risk/provenance.go`（提交后的 blob；`CR` 字节 = 0） |

（写前锚 `ce84d53d` 的同一文件是 `092590aae9465518c6511a2dee2e8d1e1819858e3d9f4c227f0ad97c487c6f68`，⛔ 与上面三个混为一谈。）

## §6 门禁（本腿现跑，串行；rc 逐条）

| 门 | 读数 |
|---|---|
| `GOFLAGS= go build ./...` | **rc=0** |
| `go vet ./internal/risk/` | **rc=0** |
| `go test -count=1 -v ./internal/risk/` | **rc=1**｜PASS=**144**、FAIL=**1**、SKIP=**1**、RUN=**251**（四数分开，⛔ 相减；都按 `^--- (PASS\|FAIL\|SKIP):` 行数，**含子测试行**） |
| 基线（写前同尺） | rc=1｜PASS=**139**、FAIL=**2**、SKIP=**1**、RUN=**247** |
| `sh scripts/d22scan.sh` | **rc=0**（`clean - no D22 ban violations`；ban #8 扫 internal/ **619** 枚 Go） |
| `sh scripts/check-path-length-budget.sh` | **rc=0**（`VERDICT GREEN`） |

**红名册（主检出，逐名）**：写前两枚 ＝ `TestResolvePerCallBudget`（计时红：`1417518 ns/op = 1.418 ms/op`，预算 1.000ms，**阈值一字未动**）＋
`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（读 `internal/panel/**` 的消费者，**归别的腿**）。
写后**只剩后者**；前者这一跑是绿的 ＝ 计时噪声，★**⛔ 当成"改好了"**。**新增红 = 0**。
★**对账**（自洽不等于正确，必须和总数对上）：`RUN 247→251` ＝ **+4，恰是本笔四枚**；`PASS 139→144` ＝ +4（本笔）**＋1**（写前那枚计时红这一跑转绿）；`FAIL 2→1`／`SKIP 1→1` 是同一件事的另一面。四笔能互相解释，⛔ 有一笔悬空。

**格式门禁（甲乙两尺都报）**：甲（判据）＝`git archive c4d65e4e` 到仓外跑
`/d/work/base/gopath/bin/gofumpt.exe -l` ⇒ **20 行红，`internal/`+`cmd/`+`tools/` 下 0 行**，与写前基线**逐行相同**（`comm` 双向零新增、零消失），20 行**全在 `.scratch/`**（那些是探针变异件）。
乙（只作诊断）＝工作树同尺 ⇒ 54 行，其中 `.scratch` 51 行、**非 `.scratch` 3 行**（`cmd\wisp\models.go`／`internal\agent\approval\pending_read.go`／`internal\agent\tools.go`，现读全为 `i/lf w/crlf`）⇒
这三行是**工作树 CRLF 假红**、甲尺里不存在、也⛔ 是本笔碰的（本笔的两枚文件在乙尺里都不出现）。★归判一律用 `[\\/]`。

## §7 这一笔**没有**做的事（登记，不许随 `-done` 蒸发）

1. `AC#2` 的**"挂到任务 `DisposalScope`"**那一支仍未做（§2-E），它要动 `internal/agent/tools/cmd` 三层＝超出 owner 批准单位。
2. `close_err=%v` 仍无人钉（§2-D）；它属 `internal/tools` 写面。
3. 无主登记**没有**新增清理路径（§2-C）——本笔只把它写成句子，⛔ 顺手改语义。
4. 票 154 的 G5 那枚 regress（`BAD 腿=G5pos 声明=quiet 实测=ring`）仍归票 171，**不是**本票射程，本笔未碰 `probes/154/**`。
5. `internal/risk` 的**计时红**未修（阈值不许动、用例不许 Skip）。

## §8 提交与名册

| 提交 | 内容 | `git show --name-only` 判 |
|---|---|---|
| `c4d65e4e` | `internal/risk/provenance.go`（25 增 1 删）＋ `internal/risk/provenance_160_boundary_test.go`（新，164 行） | 2 枚，**`.go` = 2**，**越界 = 0** |
| 本件 | `docs/evidence/s1/160-r1-impl.md` | 1 枚，越界 = 0 |

**正控**（同尺打一笔含 `.go` 的真提交）：`027ff82e` ⇒ 21 枚、`.go` = 1（`provenance.no-identity.go`）⇒ 这把尺**看得见 `.go`**。
★**另记一笔**：`provenance.go` 的工作树原本是 **CRLF**（`i/lf w/crlf`、`git diff` 为空），随手落笔会混出两种行尾 ⇒ 先按 `.gitattributes` 的 `*.go text eol=lf` 归一成 LF；**对 git 零差异**（提交 diff 仍是 25 增 1 删，⛔ 1115 行重写），甲尺本来也看不到它，乙尺那枚假红顺手消失。
