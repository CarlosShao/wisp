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

- [ ] **AC#0**：落地腿自己现跑复认 §现量 1/2/3（`gh run view 37166458550`＋那份日志），并现读肇事的**真实读取点**行号（不许照抄 §现量 4 的四行）。
- [ ] **AC#1**：红句成因定位到"哪一档/哪一路径在外部命令之前取了 `$LASTEXITCODE`"，具名 `file:line`＋该路径为什么今天走到（票 134 AC#4 的 precheck 通过后第一件事是什么）。
- [ ] **AC#2**：修法落地且**不含**上面四条 ⛔ 里的任何一条（逐条自查并在交件里点名"我没做哪一种偷懒"）。
- [ ] **AC#3**：正控两发读数齐（改前必红逐字／改后不再红**且**真超标仍红），并附还原证明（`git diff --stat scripts/slo-check.ps1` 只剩正式改动）。
- [ ] **AC#4**：另有一枚**词面/结构钉**保证"未赋值即死"不再回来（例如把"每个取值点之前必有外部命令或显式初始化"做成可判形状），并说明它能不能被换个反形绕过——绕得过就具名说绕得过，别当牙。
- [ ] **AC#5**：一次真 run 里 `slo-full` 走到出结论（`Upload SLO report` 不再 skipped）；拿不到就**具名停在"没验证"**，不许用本机一发冒充 CI 读数（既有定式：CI 与本机不同口径不可互比）。
- [ ] **AC#6**：卫生四门读数不扩大（`tools/d22scan`、`gofumpt`、`go vet`、票 212 ban #9），红名集合逐名比对；⚠ `tools/d22scan` 今天新增两枚红属票 212 射程，**不许顺手修**，看到就具名上报。

## 排程与禁区

- **按住原因（写面互斥）**：本票写面＝`scripts/**`，而 `262-r1` 正在写 `scripts/check-path-length-budget.sh`＋`.github/workflows/ci.yml`。**同目录一律串行**，不接受"改的是不同文件"这种推理（先例：09-29 我自己破过这条）。等 `262-r1` 交件并由我核过之后再派 `263-r1`。
- ⛔ 不动 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／任何 golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**`／三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）。
- ⛔ **D32 那两个数字属于 D32/D18 契约射程**：本票只修"测不到"，不许动"判定线"。任何"顺手把阈值调一下让它过"＝对抗验收直接判失败。
