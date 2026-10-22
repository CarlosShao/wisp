# 164-r2 落地腿交付件 —— 件一：修活归档仪器那半死掉的 Text 读数｜件二：AC#4 取消尾巴的判据

- 腿：**票 164 续跑腿 `164-r2`（写腿／实现者）**；工作树 `D:\work\workspace\projects plans\Wisp`；提交 `821302c9`（本件随后单独一笔）。
- 写面实动 **2 枚**：`.scratch/wisp/probes/164/v1/ac2-anchor-reading.sh`（M，件一专用）＋ `internal/tools/task_164_ac4_cancel_tail_test.go`（新，件二用例）。`internal/tools` **产码零改动**（§2.5 有逐文件 sha256）。
- 票面**一枚框没翻**：`grep -cE '^[[:space:]]*- \[x\]'`＝**4**、`- \[ \]`＝**2**（规整缩进后现跑；未勾仍是 AC#4／AC#5）。票名／`Status:`／`-done`／台账／push 全部未动。

## 0 开工前（⛔ 跳步）

**0.1 写面脏检**：起手 `git status --porcelain -- .scratch/wisp/probes/164/ internal/tools/ docs/evidence/s1/` **为空**。期间冒出**别人的**未跟踪件 `docs/evidence/s1/166-v1-verdict.md` ⇒ 未动、未 add（本次两笔提交都是逐文件名点名的 pathspec）。

**0.2 票面逐字（sed 现跑，⛔ 转述）**
```
:58 - [ ] **AC#4 取消要真取消**：取消之后不许有"还在写"的尾巴（拿一发变异证：取消后仍写输出 ⇒ 必须红）。
:59 - [ ] **AC#5 契约轴**：`docs/specs/**`、`internal/risk/**`、`internal/panel/**`、`thresholds.go`、golden、`allowlist.txt`、`frontend/**`、`design/**` 零字节；`D34` 那一行若要动＝owner 单独批准、原文落台账。
```

**0.3 票面前提新鲜度（逐条现跑，★三条已腐坏）**
- ✅ 仍真：`TestTaskOutputAC2BeforeLegIsUnreachable` 在盘；`BEFORE verbatim` 在测试件里 **0 命中**（旧仪器因此只印 PASS 行）。
- ❌ **票面 `:67` 定案⑤"`RunAsync` 生产零调用点 ⇒ 连被测对象都没有"已不成立**：`grep -rn "RunAsync" --include=*.go internal cmd | grep -v _test.go` ⇒ 生产调用点 **1 枚**＝`cmd/wisp/run.go:1176`，写入口 **1 枚**＝`cmd/wisp/run.go:1187`（`.Backfill(bg.Root().Ctx, …)`）⇒ **AC#4 的对象今天真的存在**。
- ❌ **票面 `:84`"名册零写者"已不成立**：非测试 `.Record(` 现量 **5 处**（`subagent_197.go:462`／`task.go:343`／`task.go:356`／`task_backfill.go:141`；另 `approval/gate.go:684` 等是别的 Record）。
- ❌ **票面 `:26`"生产注册表只有 `BuiltinFSEntries`"只在锚点 `f206e9f` 成立**：今天 `run.go` 还注册 `BuiltinTaskEntries`（仪器 [6] 那段的标签也该整，"TaskRoster.Record"实为所有 `.Record(`；**找到但没动**，不属本腿两件）。

**0.4 基线（四数分开，主检出）**：`go test -count=1 ./internal/tools/` **rc=0**（42.007s，起手现跑）；四数用
`go test -count=1 -v -skip '^Test164AC4CancelLeavesNoWritingTail$' ./internal/tools/` ⇒ **RUN=329／PASS=246／FAIL=0／SKIP=0，rc=0**（与票面 `:38` 记的 329／246 独立吻合）。

## 1 件一：归档仪器 `.scratch/wisp/probes/164/v1/ac2-anchor-reading.sh`

### 1.1 改前逐字输出（假绿形状，rc=0）
```
== [5] 同一枚改前用例：改前树 vs 改后树，读数是否逐字相同 ==
-- 改前树（f206e9f）--
--- PASS: TestTaskOutputAC2BeforeLegIsUnreachable (0.04s)
-- 改后树（HEAD）--
--- PASS: TestTaskOutputAC2BeforeLegIsUnreachable (0.02s)
  -> 两行 Text 逐字相同 ⇒ 该用例的读数与被测树无关（它自建注册表、且注册表是手抄的 BuiltinFSEntries）
```
⇒ **两行 Text 读数 0 行**（`:42`/`:45` grep 的 `BEFORE verbatim` 已不存在），而 `:46` 的结论**照样打印** ⇒ 假绿成立，编排者现跑核实的三条我全部复核为真。旧件 sha256＝`15fc049d…f609b8602`（`sha256sum`，SHA-256/64 位）。

### 1.2 改了什么（逐行说清；⛔ 没改回标签）
- `:32` 新增 **`LABEL='unknown-tool refusal:'`**（现行 `Logf` 标签，`task_output_ac2_before_test.go:66` 处）＋ `MISSING=0`。**grep 与 sed 共用这一条** ⇒ 标签一错，两处一起错，结论拿不到读数。
- `[5]`：两发 `go test` 先落文件（`READ_BEFORE`／`READ_HEAD`），**分别记 rc**；打印两段 `grep -E "$LABEL|^--- (PASS|FAIL)"`；结论改由 `TEXT_BEFORE`/`TEXT_HEAD`（`sed -n "s/.*$LABEL.* Text=//p"`）**实测比较**得出，且**把实测 Text 一起印出来**。
- `:46` 原那句无条件 `echo` 被**守卫**取代：读数取不到／两边不同／那一发 `go test` 没跑成 ⇒ 打印「**读数缺失、结论不成立**」（并把取到的值或 `<取不到>` 列出）＋置 `MISSING=1`。
- 结尾新增：`MISSING≠0 ⇒ rc=1`（`== 本件有读数缺失（见 [5]）：rc=1 —— 这份输出⛔当绿的凭据 ==`）。其余段落（[1]-[4]、[6]、[7]）一字未动。
- 条数：`63 → 109` 行；`tr -cd "\r" | wc -c` ⇒ **CR=0**；`sh -n` ⇒ `SYNTAX_OK`。

### 1.3 改后逐字输出（两行 Text 真印出来，rc=0）
```
== [5] 同一枚改前用例：改前树 vs 改后树，读数是否逐字相同 ==
-- 改前树（f206e9f，go test rc=0）--
    task_output_ac2_before_test.go:66: unknown-tool refusal: IsError=true ErrorClass="tool" Truncated=false Text="未知工具 task.output，可用工具见 list_tools"
--- PASS: TestTaskOutputAC2BeforeLegIsUnreachable (0.47s)
-- 改后树（HEAD，go test rc=0）--
    task_output_ac2_before_test.go:66: unknown-tool refusal: IsError=true ErrorClass="tool" Truncated=false Text="未知工具 task.output，可用工具见 list_tools"
--- PASS: TestTaskOutputAC2BeforeLegIsUnreachable (1.42s)
  -> 两行 Text 逐字相同 ⇒ 该用例的读数与被测树无关（它自建注册表、且注册表是手抄的 BuiltinFSEntries）
     Text 实测＝"未知工具 task.output，可用工具见 list_tools"
```
⇒ **`:46` 那句结论重新变成"有证据支撑"**（同一句结论，后面紧跟它比较出来的 Text）。

### 1.4 反向判据（守卫有牙）：把 `LABEL` 改回 `BEFORE verbatim` ⇒ 结论 ⛔ 再被打印
```
== [5] 同一枚改前用例：改前树 vs 改后树，读数是否逐字相同 ==
-- 改前树（f206e9f，go test rc=0）--
--- PASS: TestTaskOutputAC2BeforeLegIsUnreachable (0.05s)
-- 改后树（HEAD，go test rc=0）--
--- PASS: TestTaskOutputAC2BeforeLegIsUnreachable (0.02s)
  -> 读数缺失、结论不成立：Text 改前=<取不到>｜改后=<取不到>（grep 标签＝BEFORE verbatim）
...
== 本件有读数缺失（见 [5]）：rc=1 —— 这份输出⛔当绿的凭据 ==
```
⇒ 整件 **rc=1**，结论那句**没有出现**。改回前先证落地：`grep -n "^LABEL=" ⇒ 32:LABEL='BEFORE verbatim'`。

### 1.5 还原核（先问"谁算的、什么算法、几位"）
- 尺＝**`sha256sum`（SHA-256，64 位十六进制）**；⛔ 没拿 `git hash-object`（SHA-1/40 位）比。
- 改后件 `ec02dade9e31e492d749b8444791f3ab0f7a97ed5e7834dde85de16318b43de6`（跑出 §1.3 的那一份）。
- 反向判据跑完用**跑前留的同一份副本**还原 ⇒ `sha256sum` 回到 **`ec02dade…b43de6`**（逐字符相同）＋ `grep -n "^LABEL="` 回 `unknown-tool refusal:` ＋ CR=0 ⇒ §1.3 的读数就是对交付字节的读数。

## 2 件二：AC#4「取消要真取消，取消之后不许有还在写的尾巴」

### 2.1 今天是什么读数（先量后写，⛔ 假设）
- **对象存在**（§0.3 已证）：`run.go:1176` 起任务、`:1187` 用 `bg.Root().Ctx` 回填；子代理那一支由 `finalize` 落终态。
- **守卫已存在**：`task_backfill.go:123` `if stopped, why := Stopped(ctx); stopped {`（`Stopped`＝`cancel.go:86`，票 176 AC#3 落的）⇒ **取消之后那股输出今天已经写不进去**。
- **新用例原样跑：`--- PASS: Test164AC4CancelLeavesNoWritingTail (0.00s)`** ⇒ **AC#4 的"今天读数"是绿，⛔ 是红**。它是**判据**，不是"历史读数"；变异红（§2.3／2.4）才是它的牙。

### 2.2 判据长什么样（`internal/tools/task_164_ac4_cancel_tail_test.go`，139 行）
拼的是**两条真缝**，不 mock：`TaskRoster.Cancel`（唯一停法，`task.go:447`）＋ `TaskBackfill.Backfill`（输出写入口，`task_backfill.go:105`）。顺序＝取消 → 经**真** `subagentSpawn.finalize` 落终态（票 176-r1 的 L-4 明确把这次收尾写入排除在射程外）→ 再拿**子任务自己的 ctx** 写尾巴。四条断言：① 尾巴必须被拒且理由说"取消"；② 那一行 `Text`/`State` 一个字不许动；③ 成功 `Cancel` ⇒ 句柄真跑过（`ctx.Err()≠nil`）；④ 对照组（没取消）必须写得进去（⛔ 绿在沉默上）。

### 2.3 变异 M1（"取消只说不做"：`task.go:462` `cancel()` → `_ = cancel`）
落地证明：`grep -n "MUTATION 164-r2 M1" ⇒ 462`。逐字红（rc=1）：
```
    task_164_ac4_cancel_tail_test.go:107: 取消之后那股输出还是被写进名册了：Backfill 没有拒绝——取消之后不许有还在写的尾巴
    task_164_ac4_cancel_tail_test.go:114: 取消后名册那一行被改动了：Text "被停掉前的半句正文" -> "取消之后尾巴还在写的正文"，State "Muted" -> ""（取消之后落到这一行的写就是那条尾巴）
    task_164_ac4_cancel_tail_test.go:125: Cancel 返回成功，但它挂上的句柄没有真停掉子任务：ctx.Err() == nil（取消只是把话说出去了）
--- FAIL: Test164AC4CancelLeavesNoWritingTail (0.00s)
```

### 2.4 变异 M2（`task_backfill.go:123` `stopped` → `stopped && false`，即摘掉写入口的取消判据）
落地证明：`grep -n "MUTATION 164-r2 M2" ⇒ 123`。逐字红（rc=1）：
```
    task_164_ac4_cancel_tail_test.go:107: 取消之后那股输出还是被写进名册了：Backfill 没有拒绝——取消之后不许有还在写的尾巴
    task_164_ac4_cancel_tail_test.go:114: 取消后名册那一行被改动了：Text "被停掉前的半句正文" -> "取消之后尾巴还在写的正文"，State "Muted" -> ""（取消之后落到这一行的写就是那条尾巴）
--- FAIL: Test164AC4CancelLeavesNoWritingTail (0.00s)
```
⇒ 两发变异**可区分**：M1 多红 `:125`（取消没真发生），M2 的 `:125` 仍绿（取消真发生了、只是门不读它）。⛔ 同一个洞的两句说法。

### 2.5 三读数 ＋ 还原核
- 变异红：M1（§2.3）／M2（§2.4），各 `rc=1`。
- 还原绿／原样绿：`--- PASS: Test164AC4CancelLeavesNoWritingTail (0.00s)`（还原后与改动前各一次）。
- 还原核（`sha256sum`，SHA-256/64 位，⛔ 混 `git hash-object`）：`task.go` **`055fad4b…356fbc0`**、`task_backfill.go` **`b0b5a3ab…dd41bec`** —— 变异前与还原后**逐字节同值**（`diff` 两行 ⇒ `IDENTICAL`），`git status --porcelain -- internal/tools/task.go internal/tools/task_backfill.go` **空**，`grep -c "MUTATION 164-r2"` ⇒ **0／0**。

## 3 门禁 rc（主检出；变异读数与整包门禁**串行**跑）
| 门 | 命令 | rc | 读数 |
|---|---|---|---|
| build | `go build ./...` | **0** | — |
| vet | `go vet ./internal/tools/` | **0** | — |
| 整包 | `go test -count=1 -v ./internal/tools/` | **0** | **RUN=330／PASS=247／FAIL=0／SKIP=0**（基线 329／246，＋1＝本腿新用例） |
| d22scan | `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`（含 runtests 正控） |
| path-length | `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`，over=61、`in the wall interval=0` |

- **红名册（取主检出）：逐名（无）**；**新增红 0**。正控：两发变异都真的打出 `--- FAIL`（尺不瞎）。
- 导出树那发**不作名册**：`git archive 394cc212` 树里 `FAIL=1 TestTicket252R1AlignmentAddsNoAuthorization`＝**8.3 短名前提**（票面 `:54` 已登记），⛔ 算进本腿红。

## 4 格式甲乙两尺（都报）
- **甲（判据）**：`git archive 821302c9 | tar -x -C <仓外>` + `/d/work/base/gopath/bin/gofumpt.exe -l .` ⇒ raw **20 行、全部 `.scratch/`**（正控：尺能看见坏文件）⇒ **非 `.scratch` ＝ 0**（`tr '\\' '/'` 后用 `[\\/]` 判归属）。含别人两笔的 HEAD `46e61076` 复跑同值（20／0）⇒ 交付面不因谁先提交而变。
- **乙（工作树·只诊断）**：`gofumpt -l .` ⇒ 非 `.scratch` **3 行**＝`cmd/wisp/models.go`／`internal/agent/approval/pending_read.go`／`internal/agent/tools.go`；三者 HEAD blob 的 CR 均为 **0**、工作树 CR 分别 **389／131／225** ⇒ **纯工作树 CRLF 假象**，与本次两枚交付件无关（`gofumpt -l` 也没列我这枚）。乙 `rc=2` 来自 `.scratch` 里故意写坏的探针件（既有）。两枚交付件 `tr -cd "\r"|wc -c` ⇒ **CR=0／CR=0**。

## 5 写面／越界
- 提交 **`821302c9`**：2 files changed，+191/−6（`ac2-anchor-reading.sh` M；`task_164_ac4_cancel_tail_test.go` 新）。**产码 0 字节改动**（§2.5）。
- 越界：`internal/session|memory|risk|config|panel|agent|winsec`、`cmd/**`、`scripts/**`、`docs/specs|PLAN|reports`、`frontend|design`、`thresholds.go|golden|allowlist.txt` 全部 **0**；`add`/`commit` 两处**逐文件名**，⛔ `add -A`/`.`；⛔ push。

## 6 留给编排者拍的三句（不自行假设）
1. **AC#4**：判据已落地且有牙（M1/M2），但**今天原样是绿**——因为票 176 AC#3 的 `Stopped` 守卫已在盘、`run.go:1176/1187` 又把对象接上了。是否翻框由你判；我只报读数。
2. **AC#5 零字节**：本腿没碰契约面（§5 逐目录 0），但**框归你**。
3. **仪器 `[6]` 那枚标签**写"对 `TaskRoster.Record` 的调用"，实际 grep 是所有 `.Record(`（把 approval grant 等也算进来）⇒ 又一处"尺的射程"错；**本腿两件之外，没动**，登记给你。
