# 票 263 — `scripts/slo-check.ps1` 在 strict mode 下取**从未赋值的 `$LASTEXITCODE`** ⇒ 那枚"D32 唯一求值路径"在 checkout 被修好的**同一个清晨**、开采 5 秒后即死（这枚缺陷此前被文件名墙整层掩盖，盘上从未有人见过它的红句）

**立票时刻**：2026-10-04 09:3x +08，锚点 HEAD `ada5ed79`（`dev`）
**来路**：推送 191 枚触发的 run `37166458550`（`sha=fd269de1`，2026-10-04T00:56Z）里 `slo-full / SLO full gate (six states + settle + leak)` 那一步的**逐字红句**；日志已落 `.scratch/ci-logs/run-37166458550-failed.log`（2,053,538 字节）；台账 `A586`。母案＝票 262（那一枚治的是"到不了这一步"，本枚治的是"到了这一步却死在半路"）。

## 现量（编排者 2026-10-04 09:2-09:3x 现跑，⚠ 引用前先重跑，别把这几行当常量）

1. **checkout 已通**：同一 job 里 `Run actions/checkout@v4` / `Run actions/setup-go@v5` / `Build wisp.exe (deps cached on the runner)` 三步全 `success` ⇒ 票 262 的改名把那层墙拆掉了（**这条只说明"今天这一枚 commit 过得了"，不说明有仪器拦得住下一枚**，那格仍归票 262）。
2. **争用被脚本自己的 precheck 排除**：日志逐字 `precheck ok - no foreign toolchain/runner process, machine-wide cpu max 17%`（票 134 AC#4 那枚 validity precheck 真的在拦）⇒ 下面那发红**不是**编排者编队抢 CPU 造成的，别拿"机器忙"当解释。
3. **红句逐字**（开采第一档 `Sleeping for 6s` 之后约 0.2 秒）：
   `E:\work\base\actions-runner\_work\wisp\wisp\scripts\slo-check.ps1 : The variable '$LASTEXITCODE' cannot be retrieved because it has not been set.`
   `+ CategoryInfo : InvalidOperation: (LASTEXITCODE:String) [slo-check.ps1], RuntimeException`
   `+ FullyQualifiedErrorId : VariableIsUndefined,slo-check.ps1` ⇒ 5 秒后 `##[error]Process completed with exit code 1.`
   ⚠ **注意它死在哪**：死在**逐档采样**里，早于 `:326/:345/:366` 那三处显式 `$settleCode/$leakCode` 取值 ⇒ 肇事的读取点**不是**那三行，落点由落地腿自己现读（本条属待验断言）。
4. 脚本形状（`scripts/slo-check.ps1`，共 **397 行**）：出现 `LASTEXITCODE` 的行＝`:106`、`:326`、`:345`、`:366`；命令面逐字 `powershell -NoProfile -ExecutionPolicy Bypass -File scripts/slo-check.ps1 -Subset full -SecondsPerState 6`，env 里带 `WISP_ENV: test`、`MINGW64_ROOT: E:\work\base\msys64\mingw64\bin`。
5. 该 job 的 `Upload SLO report` 因上一步红而 `skipped` ⇒ **报告整份没落**，也就是说今天这条路上**没有任何 SLO 数字存在**（既不是达标也不是超标，是**没测**）。
6. 同一次 run 的对照读数：`slo-smoke` 那一步同样 `failure`（⚠ 是否同因**未归因**，本票不许顺手提它）。全 run `--- FAIL` 名册＝**24 枚**，vs 旧基线 `941805d0` 的 29 枚（9 枚消失、4 枚新增）；新增里两枚属 `tools/d22scan`（另一张票的射程，见 `A586` §4）。

## 为什么会发生（一句机制，⛔ 别让下一枚程把它当"PowerShell 的怪癖"绕过去）

`$LASTEXITCODE` 是 PowerShell 的**自动变量**：只有当该作用域里**真的执行过一条外部命令**之后它才存在。在 `Set-StrictMode -Version Latest`（或等价严格档）下，**取一个从未赋值的变量＝运行时异常**，而脚本恰好在"这一档还没跑任何外部命令"的路径上取了它 ⇒ 红句落在采样档而不是落在三处显式取值行。
⇒ 净结果：只要某一档的采样顺序里那次外部命令没跑／跑前就取，**整道 D32 门当场死掉**，且**它在过去几天根本没被执行到过**（被票 262 那层墙挡在前面），所以今天才第一次露出来。**这是"两层故障互相掩盖"的形状，不是编码手滑。**

## 要建什么

- **修**：把取值改成**先赋再取**（进入采样循环前显式 `$LASTEXITCODE = 0` 一形，或改用"跑完外部命令才读"的结构），⛔ **不许**用 `try/catch` 吞掉、⛔ **不许**放宽 `Set-StrictMode`、⛔ **不许**为了变绿改任何**阈值/golden/`internal/observe/thresholds.go` 一字节**、⛔ 不许把任何一步从"必须出结论"改成"跳过"。
- **正控**（按既有定式：负向尺必配"种 X 必响"）：交件要含两发——① 改前：那枚采样档取未赋值变量的红句**逐字重现**（可用最小复现件在 `WISP_ENV=test` 下跑，不许拿今天 CI 那发当自己的读数）；② 改后：同一发**不再红**，且**必须证明它仍会因真超标而红**（把某一档的判据临时改成必红形状 ⇒ 门红 ⇒ 还原，两发逐字）。
- **诚实档**：脚本如果本来打算"取不到结论就报 `NO CONCLUSION`"，那今天这发**连 NO CONCLUSION 都没走到**＝异常直接终止。要把"未赋值即死"这一形变成**显式具名的一种失败**（而不是运行时异常），否则下一任看到的还是一个 exit 1 而不知道死因。⚠ 既有读数：step 层的绿可以等于"什么都没测"（`NO CONCLUSION (machine-contended)`＋`exit 0`），本票不许把这一形做得更糟。
- **销账**：本票交件后**必须真出现一次 `slo-full` 出结论的 run**（`Upload SLO report` 不再 skipped），并在那之前任何地方都不许写"D32 已被 CI 保护"。

## 判据（AC 框由编排者翻，产码腿一枚都不许碰）

- [x] **AC#0**：落地腿自己现跑复认 §现量 1/2/3（`gh run view 37166458550`＋那份日志），并现读肇事的**真实读取点**行号（不许照抄 §现量 4 的四行）。
- [x] **AC#1**：红句成因定位到"哪一档/哪一路径在外部命令之前取了 `$LASTEXITCODE`"，具名 `file:line`＋该路径为什么今天走到（票 134 AC#4 的 precheck 通过后第一件事是什么）。
- [x] **AC#2**：修法落地且**不含**上面四条 ⛔ 里的任何一条（逐条自查并在交件里点名"我没做哪一种偷懒"）。
- [x] **AC#3**：正控两发读数齐（改前必红逐字／改后不再红**且**真超标仍红），并附还原证明（`git diff --stat scripts/slo-check.ps1` 只剩正式改动）。
- [x] **AC#4（口径由验收腿 263-v1 §⑥ 收窄后翻，⛔ 不许读成"已保证不再回来"）**：★这枚钉买到的**只有本文件两枚自锁钉**（词面钉＋能力钉），**射程必须随框一起读**——三形实测绕过（运行期 `Get-Variable` 取名／把两枚钉的调用删掉⇒**rc=0 全绿、本机零外牙**／读取挪去别的文件⇒**CI 原形红句逐字回来而两枚钉全打 ok**）；外牙只有 `slo-freshness.sh` 的 P3 aging。凭据更正：交付字节里 `LASTEXITCODE` **字面 2 行（`:241` 注释、`:252` 模式串）、带 sigil 0 处**（实现腿原写"1 次"是枚数错）。⚠ 另记一枚**反方向**读数：只在注释里写这三个词（零读取）也会把整道 D32 门打死（rc=1）——它比语义宽，故"在本文件写解释性注释"从今天起是有代价的动作。
- [ ] **AC#5**：一次真 run 里 `slo-full` 走到出结论（`Upload SLO report` 不再 skipped）；拿不到就**具名停在"没验证"**，不许用本机一发冒充 CI 读数（既有定式：CI 与本机不同口径不可互比）。⇒ **停在未勾**：本机 12 发全是 `-Subset smoke -SecondsPerState 1`×假件，证的是**路的形状**、不证 D32 那两个数；推送窗口在我手里（取数时刻 `2026-10-04 13:13` 现跑 `git rev-list --count origin/dev..HEAD`＝**81 枚未推**）。
- [x] **AC#6（两门适用＋两门不适用，具名）**：`tools/d22scan`（rc=0，`ban #8 internal/=503`、`cmd/=100`，比 12:0x 各 +1 ⇒ 归因＝**260-r4 新增的两枚测试件**，与本票无关）＋票面长度的门（rc=0，`VERDICT GREEN`）＝**本腿现跑的两门**；`gofumpt`/`go vet` 因"五枚 commit 零 tracked Go 产码"（验收腿自己拉的 `git show --name-only` 名册）**不可能因本票变红**，⛔ 下一任不许以为这两门真跑过。ban #9：本机 d22scan clean ⇒ 未扩大。`tools/d22scan` 那两枚"CI 红／本机绿"本腿没碰一字、按票面具名上报。

## 编排者更正（2026-10-04 16:5x）：本票 `:53` 那句 scope 名册**写宽了一枚**（⛔ 原句不抹，按这一条读）

- `:53` 逐字写着 `tools/d22scan` 的 scope 名册"只列 internal/cmd/frontend/design/**tools**"——**`tools/` 那一枚不在名册里**。凭据＝只读腿 `266-a1` 指出后，**我 16:5x 自己现读源码复认**：
  - `tools/d22scan/main.go:55-59` 逐字写着 ban #9 "Judged against … the production Go files of **internal/ and cmd/ ONLY** - **`tools/**` and `_test.go` are outside this ban's range**，…**adding either tree is a scope change for an owner to approve, not a test fix**"；
  - 走查那几枚调用＝`main.go:264`（`internal`）／`:267`（`cmd`）／`:271`（`frontend`，ban #6）／`:274`（`internal/tools`，ban #7）／`:581-582`（`design`／`frontend`，ban #8）⇒ **仓根那枚 `tools/` 目录今天不被任何一条 ban 走查**（`:260` 只是去**读** `tools/d22scan/allowlist.txt` 这枚白名单文件，不是扫那棵树）。
- **这一处写宽不是纯洁癖**：它正是票 266 的命门——"**仪器自己的源码不在任何外牙射程内**"。我把它列进 scope 名册＝把"仪器有外牙"这句说大了。⇒ **`266-a1` 的更正成立**，本票 `:53` 那句按上面读。
- ⚠ **为什么我当初会写宽**（记我）：我是照 `AGENTS.md` 的"薄索引"口气与 d22scan 的**自述文案**拼的，⛔ 没有逐条走 `walkGo/walkText` 的调用点。⇒ 定式：**凡把"某道门扫哪几棵树"写进票面，必须现读那几枚 walk 调用，不读注释**。

## 排程与禁区

- **按住原因（写面互斥）**：本票写面＝`scripts/**`，而 `262-r1` 正在写 `scripts/check-path-length-budget.sh`＋`.github/workflows/ci.yml`。**同目录一律串行**，不接受"改的是不同文件"这种推理（先例：09-29 我自己破过这条）。等 `262-r1` 交件并由我核过之后再派 `263-r1`。
- ⛔ 不动 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／任何 golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**`／三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）。
- ⛔ **D32 那两个数字属于 D32/D18 契约射程**：本票只修"测不到"，不许动"判定线"。任何"顺手把阈值调一下让它过"＝对抗验收直接判失败。

## Progress log

- **2026-10-04 13:1x（编排者翻勾，凭据＝非实现者验收腿 `263-v1` 的 `.scratch/wisp/probes/263/v1/verdict.md`，304 行／33,526 字节／占位 0，件在盘 `515ca5c5`）**：翻 **AC#0/AC#1/AC#2/AC#3/AC#6** ＋ **AC#4 带条件**（口径见框内收窄那句）；**AC#5 停勾**。被验的五枚 commit＝`82af691d`→`b92037bd`→`cb8dcecc`→`befb779c`→`c9411a1e`，实现腿＝`263-r1`；验收腿自己拉的名册证明 **probes 之外只碰了 `scripts/slo-check.ps1` 一枚**，禁区逐名零命中，`internal/observe/thresholds.go` 四锚 md5 逐字节同。
- ★ **本票那条因果链的结论（要单独立档，已入 `A598` §4／`A599`）**：肇因不是"PowerShell 的怪癖"，是**票 244 在 `scripts/build.ps1:115` 加的 `-H=windowsgui`** ⇒ `wisp.exe` 成 GUI 子系统二进制 ⇒ PowerShell 的 `&` 不等它退出 ⇒ `$LASTEXITCODE` 从未被赋值 ⇒ 严格模式在**逐档采样**那一发当场抛。⇒ 这答了本仓一直没答的问题：**"D32 那两个数字为什么从来没出现过"**——不是不达标，是**求值器第一次真跑到它就死**，且此前被票 262 那堵"文件名太长导致 checkout 先失败"整层掩盖（两层故障互相掩盖）。
- ⚠ **一句机制收紧（验收腿自己量的，覆盖票面"当场死"那句的射程）**：最小复现件那一发 **rc=0**——异常是**语句级作废**、不是脚本级终止。⇒ "取未赋值变量＝当场死"只在**门本体的形状**下成立（rc=1、无 report）；同一形状落在别的语句位置上，存在 **rc=0 而什么都没测**的可能。票面 §为什么会发生 那段按此读，不许把"当场死"写满。
- ★ **验收腿量到一枚名册外的第三牙**：`delete-nails-and-wait` 那形里门会**伪造 `exit=0`**（没等就读到 0），最后兜住它的是 leak 自检要求 `exit==1`（`scripts/slo-check.ps1:594`/`:596-598`）。⚠ 残余风险具名：哪天 leak 的形状不再要求 1，这一族就只剩**可被就地删除**的自锁钉。
- **归口（本票不做，另立）**：① 给 `scripts/*.ps1` 找一枚**外牙**（`tools/d22scan` 的 scope 名册只列 internal/cmd/frontend/design/tools，删钉本机零尺会发现）；② 改出厂件子系统（`-H=windowsgui` 这一类）要带**下游消费者名册**的普查尺——⚠ 这**不属 263 也不属 244 任何一方**，由编排者裁要不要立票（记在 `A599`）。
