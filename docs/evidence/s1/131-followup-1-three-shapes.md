# 票 131 续单 第 1 件：三形（X4／X8／X12）打进枚举门，另加两形自证与一枚邻票新形的答复

会话：`worker-ticket131-followup-r2`。接续撞轮次上限的前任 `worker-ticket131-followup`
（它留下的三枚未提交件本轮全部保留并在其基础上继续，未 `reset`／未 `stash`／未 `--amend`）。

- 锚点：开工 `git rev-parse HEAD` = **`b723978c0f17612c5e3210670fa6dda9db8e94ff`**（13:3x 实测；简报说的"约 b723978"成立）。
  跑到 14:3x 共享树被兄弟代理推进（`7b4c36a` → `c66d683`），新落的 5 枚里**零枚碰 `cmd/wisp/**`**（逐枚 `git show --name-only` 查过），
  所以本档的分母仍是 `b723978`。
- 判据原文：`docs/evidence/s1/131-adversarial-acceptance.md` §4（X1..X12）、§6（`R-131-1..4`）、§8.3（结案三件事）。
- 简报给的行号/枚数是断言：`main.go:54-93` 与"基线约 100 枚"实测成立；`slo_windows.go:179` 那一格实为 `:171`（130 在同段追加过句子）；
  "P1 ⇒ helper 仍须红"**不成立**（读数在 `131-followup-2-ratification-and-linux.md` §2）。

## 0. 次序：先测、后写、再提交

前任的 `cmd/wisp/leg_sink_gate_131_test.go` 文件头写着三形"each was re-measured on a planted leg
(readings in docs/evidence/s1/131-followup-1-three-shapes.md)"，而这份文件当时**不存在**
（`ls docs/evidence/s1 | grep 131-followup` 只有 `131-followup-0-crossplatform.md`）。
本档就是补齐的那一步：下面每一发都是**先量完**才写进注释，码与本档同批入库。

仪器（每一发固定五步）：
`git archive b723978 | tar -x -C /tmp/wisp131-r2-<发>` ⇒ 拷 `third_party/sherpa-onnx/*.dll` 三枚进快照
（少了它们测试二进制在装载期 `0xc0000135`、四数全零）⇒ 变异 `diff`/`grep` 证落地 ⇒ `go build ./cmd/wisp/` rc=0 ⇒
`go test -v` 取红名原文 ⇒ 还原 + `diff -r` 证快照回到未变异态 ⇒ 复绿。全程不在仓内建 worktree。

**仪器同一性**（免得读数对不上码）：最终码三枚件的 sha256 前 16 位在 `/tmp/wisp131-r2-fin`（§5..§7 用的快照）
与工作树逐枚相同：`leg_sink_gate_131_test.go 259cdc4802c7accb`、
`leg_sink_nail_131_windows_test.go 69106117257dba2d`、`slo_windows.go bd5c02ca11fd745d`。
本文 §2..§5 四发的**整包四数**与 §10 全部数出自这份最终码：四发各在 `/tmp/wisp131-r2-w4-<发>.log`
（`go test -count=1 -v ./cmd/wisp/`，四发同为 `RUN 100 / 52 / 1 / 0`、唯一顶层红都是门），
门内红名另在 `/tmp/wisp131-r2-fin-<发>-gate.log` 复测一遍同形。§8 六发的四数取自 13:5x 首轮
（那一版还没有 §6 的第二把尺），它们的**红名文本**在最终码上逐发复测过、逐字相同，
只有门自己的 `leg_sink_gate_131_test.go:<行号>` 前缀因文件加长下移（正文两批都给了）。

⚠ 一枚仪器坑（登记，不硬改）：`/tmp/wisp131-r2-fix` 对照树在 14:4x 之前存的是**旧一份**门文件
（sha `abbff201b26daf57`，还没有 `rulingCrossChecks131`），所以那一版 `wisp131-r2-final.out` 里每发末尾的
`reset: DRIFTED` 是**对照树落后**、不是"变异没还原干净"。把对照树刷成工作树同一份之后，
`diff -rq` 对 `fin` 与 `w4` 两棵树的还原态报**零差异**（15:0x 复测）。


## 1. 新基线（`b723978` 纯净快照，未加本轮三枚件）

| 仪器 | 读数 |
| --- | --- |
| `go test -count=2 -v ./cmd/wisp/` | rc=0，`=== RUN` **200** / 顶层 **PASS 106 / FAIL 0 / SKIP 0**（子用例 94/0/0），`ok` 104.830s |
| `sh scripts/wisp-cli-tests.sh`（CI 形状） | rc=0，`=== RUN` **100** / 顶层 **53 / 0 / 0**（子用例 47），`ok` 49.824s |
| `sh scripts/d22scan.sh` | rc=0，八 scope `203 / 22 / 40 / 18 / 16 / 40 / 390 / 37` |

与验收方 §2 的 `82 / 46` 差多少、差在谁身上：验收方锚 `56d8026`，其后 `ea05cf5`(129) 与 `9b5d64d`(130)
往 `cmd/wisp` 追加了用例；前任在 `3029415` 量到 `100 / 53`，我在 `b723978` 复测逐数相同
⇒ **差 +18 RUN / +7 顶层 PASS 全归 129/130，本件零贡献**；`717d822..b723978` 对本包分母零影响。
⇒ 本档全部"四数不降"都对 **100/53/0/0**（CI 形状）与 **200/106/0/0**（`-count=2`）比。

## 2. 形 a ＝ X4：早退 `if` 分发一条装了听众的腿（红名必须点到 `--diag`）

造法（照验收方 X4）：`main()` 里 `switch args[0] {` 之前插

```go
	if len(args) > 0 && args[0] == "--diag" {
		os.Exit(cmdDiag131(args[1:]))
	}
```

+ 新生产文件 `cmd/wisp/diag131.go`（`installLogSink(root)` + `defer sink.close()`），不给它钉、不碰测试文件。
落地：`diff -u` 只有那 4 行（快照 `main.go:60`）；`grep -c installLogSink cmd/wisp/diag131.go` = 1。
`go build` rc=0、`go vet` rc=0。先证这条腿是活的（`go build -o /tmp/wisp131-r2-x4.exe ./cmd/wisp`，真二进制）：

```
$ /tmp/wisp131-r2-x4.exe --diag <根>      # rc=0
level=INFO msg="wisp: persistent log sink installed" dir=...\wisp131-x4data\root\logs min_level=info early_records=1 early_dropped=0
<根>\logs\wisp-20260923-001.jsonl:{"time":"...","level":"INFO","msg":"wisp: persistent log sink installed","dir":"...\\root\\logs",...}
```

**修复后读数**（最终码，`-count=1 -v -run TestAC4EveryLegIsNailedOrRuled`）：rc=**1**，红名原文两枚：

```
    leg_sink_gate_131_test.go:348: AC#4 RED: main.go:60 reads args (the local os.Args was bound to) in func main outside every branch this gate classifies (the argv binding, the `if len(args) == 0` leg, and the `switch args[...]`). Literals in that statement: "--diag"
        Statement: if len(args) > 0 && args[0] == "--diag" { os.Exit(cmdDiag131(args[1:])) }
        A dispatch that lives beside the switch is invisible to the leg list, so the leg it selects can install a listener with no nail and no row. Fix: move it inside the switch, or teach this gate the new shape and make it prove the reading on a planted leg.
    leg_sink_gate_131_test.go:348: AC#4 RED: main.go:61 reads args (the local os.Args was bound to) in func main outside every branch this gate classifies (...). The statement names no literal command, so this reader cannot say which leg it dispatches.
        Statement: os.Exit(cmdDiag131(args[1:]))
```

整包同发：rc=**1**，`RUN 100 / 顶层 PASS 52 / FAIL 1 / SKIP 0`（基线 53/0）——唯一一枚红在门上。
账本仍 15 行、`--diag` 零命中 ⇒ 拦下这一形的不是"多了一条腿"，是"越权读 argv 就红"。
**⇒ 结案判据 2 的"红名点到 `--diag`"逐字成立**（第一枚的第一行）。
还原：删 `diag131.go` + 还原 `main.go` ⇒ `diff -r` 与未变异快照只差一枚 `go build` 产物 ⇒ 整包 `100 / 53 / 0 / 0` 复绿。
零误伤：未变异树上 `strayArgvReads131` 零命中（§10 那趟 `AC#4 RED` 计数 = 0）。

## 3. 形 b ＝ X8：`case "slo":` 的标签写成命名常量

造法：`case "slo":` ⇒ `case sloCmdName131:`，并加 `const sloCmdName131 = "slo"`（`diff -u` 两处）。build/vet rc=0。
修复前（验收方 X8）：门 `PASS`、`leg slo` 静默出账、换成两行 `default`。
**修复后**（最终码，门那一发）：

```
    leg_sink_gate_131_test.go:337: AC#4 RED: the case label at main.go:84 is sloCmdName131, which is not a string literal. It is package-level const sloCmdName131 = "slo", so the leg this clause dispatches is "slo" - a name this reader would only know by walking somebody else's declaration, which is the assumption that just got a leg dropped.
        This gate enumerates legs by their labels, so a label it cannot read is a leg it cannot name: the row above carries it as "unparsed-label@main.go:84:sloCmdName131" instead of dropping it or calling it default.
        Fix: dispatch on the literal, or teach caseLabels131 the form and make it prove the reading on a planted leg.
```

整包同发：rc=1，`100 / 52 / 1 / 0`。账本里那条腿以 `unparsed-label@main.go:84:sloCmdName131` 自己的行在账，
`default` 那一行不受影响 ⇒ 验收方点名的两个后果（静默出账 / 与别人撞键后静默）各自都有了对面的红。
还原 ⇒ `diff -r` 零差异、门 `--- PASS`。

## 4. 形 c ＝ X12：`var 别名 = installLogSink`（外加 X11 那一侧）

造法与验收方 **逐字同形**：把它的 `cmd/wisp/fake131.go` 原样拷进快照
（`diff -q /tmp/wisp131-acc131r2-x12-aliasnewleg/cmd/wisp/fake131.go` 同字节；内含
`var sinkAlias131 = installLogSink` + `sinkAlias131(root)` + `defer sink.close()`），
`switch args[0] {` 后插 `case "fake131": os.Exit(cmdFake131(args[1:]))`。build/vet rc=0。
活证（真二进制）：

```
$ /tmp/wisp131-r2-x12.exe fake131 <根>     # rc=0
level=INFO msg="wisp: persistent log sink installed" dir=...\wisp131-x12data\root\logs
<根>\logs\wisp-20260923-001.jsonl: {"msg":"wisp: persistent log sink installed", ...}
```

**修复后**：门红，且与 X6（直呼 install 的新腿）**同一枚红名**：

```
    leg_sink_gate_131_test.go:275: AC#4 RED: leg "fake131" (main.go:61) reaches installLogSink on 1 line(s) of this package and no registered nail names it.
        Delete that install block and this package stays green - the reading ticket 127 measured for models, and ticket 117 for resident.
        Fix: write the nail, and claim it with registerLegNail131("fake131", TestYourCase) in the file that owns it.
    leg_sink_gate_131_test.go:371: call edges added through package-level function values (R-131-1 形 c): cmdFake131:sinkAlias131->installLogSink
```

整包同发：rc=1，`100 / 52 / 1 / 0`。还原 ⇒ 零差异 ⇒ `100 / 53 / 0 / 0`。

X11（同一形打在**已有钉**的 models 腿上）：验收方量到"门确实红、但红名是假话"
（`which no longer installs the persistent sink`，而它装得好好的）。本轮改的是"看得懂别名"而不是"红得更响"：
`models.go` 加 `var sinkViaVar131 = installLogSink`、调用点改走别名 ⇒
门 `-count=1 -v -run 'TestAC4...|TestAC2Models...'` rc=**0**、两枚 `--- PASS`，并在日志里报出自己新加的边：

```
    leg_sink_gate_131_test.go:371: call edges added through package-level function values (R-131-1 形 c): modelsEnsure:sinkViaVar131->installLogSink
```

⇒ 账本恢复真话（`leg models install=true ... -> nailed`），假话红没有存在必要。
还原 ⇒ 零差异、门 PASS。

## 5. 第四形 ＝ 别名"放不下"那一半（X13）：红着说看不见，不许当叶子走过

前任注释主张这一半，但它落地的码**举不出可编译的实例**：`resolveAlias131` 只在名字进了 `pkg.declared` 时才报盲，
而 `declared` 原先只记**有初值**的 `var`/`const`/`type` ⇒ `var factory func(...)`（无初值、由 `init()` 填，
正是"按平台分派／注入替身"那形）会被判成"不是本包的名字"而**静默当叶子**。本轮两处改：

1. 包级 `ValueSpec` 的每一枚名字都记进 `declared`（与有无初值无关）；
2. 盲边的红只在**有人 call 它**时报，并把 caller 名单写进读数（没人 call 的仍是叶子，理由写在码旁）。

造法（`cmd/wisp/blind131.go`，可编译、可跑）：

```go
var sinkFactory131 func(string) (*logSink, error)   // 无初值
var callSink131 = sinkFactory131                     // 腿调用这枚别名
func init() { sinkFactory131 = installLogSink; callSink131 = sinkFactory131 }
```

+ `case "blind131": os.Exit(cmdBlind131(args[1:]))`。build/vet rc=0。
活证：`/tmp/wisp131-r2-x13.exe blind131 <根>` rc=0，`<根>\logs\*.jsonl` 里 `persistent log sink installed` 命中 1 次。
红名原文（门自认是仪器问题）：

```
    leg_sink_gate_131_test.go:360: AC#4 RED (the instrument, not the code): var callSink131 = sinkFactory131, called by cmdBlind131: this reader cannot place what the alias holds, and the functions named above call through it.
        This gate walks names; an edge it cannot resolve is an edge a leg can reach installLogSink through without the ledger noticing. Say which it is: point the alias at a function this package declares, or teach loadMainPackage131 the form and prove the reading on a planted leg.
```

整包同发 `100 / 52 / 1 / 0`；还原 ⇒ 零差异 ⇒ 门 PASS。未变异树上这一枚红的命中数 = **0**（零误伤）。

## 6. 第五形＝回答邻票：**Y3 打在这扇门上成不成立**

`129-adversarial-acceptance.md` §4.4 的 Y3 是"**改仪器自己的一行**，让一列去读自己的期望值"，
它同时指出"承重只能是行数读数"这句话自己没有证人。两个问题分开答。

**(1) 先给"那枚标记到底承不承重"的证人**（这是 A0，仪器一字未动）：
只改 `cmd/wisp/slo_windows.go` 一处，把 `WISP-LEG-SINK-RULING:` 改成 `LEG-SINK-DECISION-REMOVED:` ⇒

```
    leg_sink_gate_131_test.go:298: AC#4 RED: leg "slo" (main.go:82) books records (slog.InitLog@slo_other.go:34 | slo_windows.go:185) through the process logger, installs no listener, and carries no WISP-LEG-SINK-RULING: sentence.
--- FAIL: TestAC4EveryLegIsNailedOrRuled (0.03s)
```

⇒ 删掉那一发裁决，红的是 `TestAC4EveryLegIsNailedOrRuled` 本人，且点名 `slo`。
覆盖面主张现在答得出"删掉它哪条用例会红"了。

**(2) 再把 Y3 那一形照搬进门里**，逐枚实测（每一枚都 `go build` rc=0 后才读数）：

| 枚 | 门里被改掉的那一行 | 结果 |
| --- | --- | --- |
| A1 | `if strings.Contains(fi.doc, legSinkRulingMarker) {` ⇒ `if !fi.calls[installsFuncName] {`（不装听众就算裁决过）＋标记删掉 | **不绿**：4 枚新红 `leg "models"/"no-args"/"run"/"secret" carries a WISP-LEG-SINK-RULING: sentence (...) AND reaches installLogSink` |
| A2 | install 谓词 ⇒ `if len(fi.slogSel) > 0 {`（发过记录就算装过听众）＋叠 X1（删 models 那 9 行） | **不绿**：3 枚红，其中 X1 原有的 `which no longer installs the persistent sink` 仍在 |
| A3 | records 谓词 ⇒ `if len(fi.slogSel) >= 0 {`（恒真）＋叠 X2 | **不绿**：每条不发记录的腿各红一枚 `books records () ... carries no WISP-LEG-SINK-RULING: sentence` |
| A1b | 同一枚谓词只加一个**指名豁免**：`... \|\| fi.name == "cmdSLO" {` ＋标记删掉 | **补判据前：静默**（`rc=0`、`--- PASS`、账本仍写 `leg slo ... ruled=true`，而那枚句子已不在文件里）；**补判据后：红**（见下） |

⇒ **答复：成立一半。** 整列拉满那三支在这扇门上撞得出矛盾红（因为门另有两条互斥读数：
`installs && ruled` 必红、`records && !installs && !ruled` 必红），**但"指名豁免一条腿"这一支真的静默**——
A1b 在补判据之前的旧码上实测 `--- PASS`，这就是本票原本藏着的那一格。

补的那一枚判据（`rulingCrossChecks131`：第二把尺，扫**字节**不走 AST doc）把 A1b 变成红的：

```
    leg_sink_gate_131_test.go:379: AC#4 RED: leg "slo" is booked as ruled at [slo_other.go:34 | slo_windows.go:185], but the marker sentence "WISP-LEG-SINK-RULING:" is in none of those files: the reading came from the instrument's own predicate, not from a ruling anybody wrote down.
```

同一把尺反向也咬：`X carries a "WISP-LEG-SINK-RULING:" sentence that no ruled ledger row claims:
a ruling nobody's dispatch reaches is prose, not a decision.`
⇒ 从此"往门里塞一枚指名豁免"必须同时瞒过两把尺才静默；这不叫堵死，叫把一行的成本抬到两把尺。

**两枚误伤实测（新判据自己不是一枚免检主张）**：
① 第一版把标记扫描对所有 `.go` 文件都做 ⇒ 本文件自己那行 `const legSinkRulingMarker = "WISP-LEG-SINK-RULING:"`
被读成"一枚没人认领的裁决" ⇒ 未变异树上门当场红（原文
`slo_windows.go carries a ... no ruled ledger row claims` 与 `leg "slo" is booked as ruled at ... in none of those files` 各 1 枚）；
改成测试文件不进扫描。② `cmdSLO` 在 `slo_windows.go` 与 `slo_other.go` 各有一枚声明 ⇒ 它的 site 是
`slo_other.go:34 | slo_windows.go:185` 这种**合并串**，只取第一段就看不见带标记的那半 ⇒ 未变异树上又红一次；
改成按 `" | "` 逐段查。**两枚都是在未变异树上量出来的红，不是推理**。
最终码在未变异树上：门 `--- PASS`、整包 `-count=2` 里 `AC#4 RED` 命中数 **0**（§10）。

## 7. `R-131-2`（X5）：钉的归属列现在可核

注册带第三枚参数（被驱入口符号），门核三条：入口是本包生产函数、这条腿的闭包真能到它、用例真 **call** 它；
`subprocess:` 前缀只免掉第三条并把那一行标成 name-only claim（127 的 resident 钉就是这形状，日志原文：
`nail "TestAC1ResidentLeg..." -> leg "no-args", entry "runResident": declared with the "subprocess:" prefix, so this row is a name-only claim ... NOT evidence that anybody asserts "no-args"'s semantics in-process.`）。

X5（只换腿名）在最终码下两枚红：

```
    leg_sink_gate_131_test.go:324: AC#4 RED: nail "TestAC2ModelsLegBooksItsHandOffVerdictOnDisk" claims leg "secret" and entry "cmdModels", but the dispatch main.go makes for "secret" never reaches "cmdModels" (main.go:72). The case drives some other leg's code: the ledger's nailed row for "secret" is a claim about the wrong function.
    leg_sink_gate_131_test.go:324: AC#4 RED: nail "TestAC3SecretLegBooksItsAuditRecordsOnDisk" claims leg "models" and entry "cmdSecret", but the dispatch main.go makes for "models" never reaches "cmdSecret" (main.go:74). ...
```

⇒ 验收方 X5 那发"82/82 全绿、两行各报 nailed"不再可能。

## 8. 六发红名逐字复核（结案判据 1）

最终码、同一台机器、`go test -count=1 -v`；下表"本轮红名"与验收方 §4/§10 引文**逐字相同**，
只有门自己的 `leg_sink_gate_131_test.go:<行号>` 前缀与钉文件的 `<行号>` 因文件加长下移（前→后：`171→275`、`194→298`、`214→318/324`；钉 `307→358`、`455→506`、`278 helper 调用点→366`）。

| 发 | 造法（同验收方） | 本轮整包四数 | 红名 |
| --- | --- | --- | --- |
| X1 | 删 `models.go` 那 9 行 install | `100 / 50 / 3 / 0`（+1 子用例） | `AC#4 RED: nail "TestAC2ModelsLeg..." claims leg "models", which no longer installs the persistent sink (main.go:74): the nail is now proving something else, or nothing` ＋钉两枚，逐字同 |
| X2 | 删 `secret.go` 那 5 行 install | `100 / 48 / 5 / 0`（+1 子用例） | 两条独立读数（要裁决 + 要钉）逐字同；**多两枚红**＝票 130 的 early-record 用例，方向是查得更多 |
| X3 | 删掉 `case "models"` 整支、`models.go` 一字不动 | `100 / 52 / 1 / 0` | `... which main.go does not dispatch to any more` 逐字同 |
| X6 | 新腿直呼 install、不给钉 | `100 / 52 / 1 / 0` | `leg "fake131" (main.go:61) reaches installLogSink on 1 line(s) ... no registered nail names it.` 逐字同 |
| X9 | 两处 install 挪到 `<根>-elsewhere` | `100 / 48 / 5 / 0` | 门 `PASS`（分工：门管可达、钉管落点，与验收方一致）；两枚钉各红一次，红名带 `dir=...-elsewhere\logs` |
| X10 | 注册一枚不读盘的用例当钉 | `100 / 52 / 1 / 0` | `... calls none of the shared sink readers (readLegSink131, readResidentSink, readSink) ...` 逐字同，另**多一枚**入口符号缺失的红（§7） |

X1 的钉侧原文（证明"点名 models 那条腿"没退化）：

```
    leg_sink_nail_131_windows_test.go:358: no wisp-<day>-<seq>.jsonl in ...\TestAC2ModelsLegBooksItsHandOffVerdictOnDisk1466134452\001\logs: the leg ran to completion and its listener wrote nothing there.
    leg_sink_nail_131_windows_test.go:506: a log directory that cannot be opened produced no named refusal ("wisp models: 持久日志未启用" on stderr).
```

## 9. 本轮对前任码/注释改了什么（全在 `cmd/wisp/**`，逐条带读数）

1. 形 a 注释里"一发腿只红一枚"与实测不符（X4 实测两枚，一条语句一枚）⇒ 注释改成实测形状（§2 的两枚红）。
2. `pkg.declared` 只记有初值的名字 ⇒ 会静默放过"无初值包级函数值"那形 ⇒ 改为全记 + `callersOf131`（§5 的 X13 就是这一改的读数）。
3. 盲边红名原文从"…is a package-level function value whose target this reader cannot place, and a leg calls through it"
   改成对读数负责的 "%s: this reader cannot place what the alias holds, and the functions named above call through it"（caller 名单在 `%s` 内）。
4. 新增 `rulingCrossChecks131` + 加载期字节扫描（§6 的第五形），并修掉它自己的两枚误伤（§6 末）。
5. `slo_windows.go` 的裁决段措辞（`R-131-4`）落点与复核在 `131-followup-2-ratification-and-linux.md` §3。

## 10. 结案三件事 + AC#5 八发（同一枚最终码）

| 判据 | 结论 | 依据 |
| --- | --- | --- |
| 1. 六发红名逐字不变（X1/X2/X3/X6/X9/X10） | **成立**（X2/X9/X10 各多红，全部归因到 130 的用例与本轮 R-131-2 的新红） | §8 |
| 2. 三发从绿变红（X4/X8/X12），X4 红名点到 `--diag` | **成立** | §2 §3 §4 |
| 3. §1 那张表四数不降 | **成立**（下表逐数相同） | 下表 |

| 仪器（最终码） | 基线 `b723978` | 修复后 | 判定 |
| --- | --- | --- | --- |
| `go test -count=2 -v ./cmd/wisp/` | 200 / 106 / 0 / 0，rc=0 | **200 / 106 / 0 / 0**，rc=0，`ok` 103.200s | 不降 |
| `sh scripts/wisp-cli-tests.sh` | 100 / 53 / 0 / 0，rc=0 | **100 / 53 / 0 / 0**，rc=0 | 不降 |
| `sh scripts/d22scan.sh` 八 scope | 203/22/40/18/16/40/390/37，rc=0 | **同一枚八数**，rc=0 | 不降 |
| `gofmt -l cmd/wisp/` | 空 | 空（rc=0） | 同 |
| `$(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp/` | 空 | 空（rc=0） | 同 |
| `go vet ./cmd/wisp/`（windows） | rc=0 | rc=0 | 同 |
| `go build ./cmd/wisp/`、`go build ./...` | rc=0 / rc=0 | rc=0 / rc=0 | 同 |
| `GOOS=linux go vet ./cmd/wisp/`（交叉，只编译） | rc=1，唯一错误行是 sherpa 那句 | rc=**1**，逐错误行归因：**仅** `imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in D:\work\base\gopath\pkg\mod\...\sherpa-onnx-go-linux@v1.13.8`（零枚 `cmd/wisp/*.go` 出现在错误里） | 同 |

⚠ 那句交叉编译的 rc=1 **不许**当"整树假象"整体归因：Linux 原生（Docker `golang:1.27`，CGO_ENABLED=1）
`go vet ./cmd/wisp/` 的真读数是 **rc=0**，见 `131-followup-2-ratification-and-linux.md` §1——
上一轮就是这三个代理把它整体归因掉，才让真伤带病 4 小时。

## 11. 未验证项 / `next=`

- `next=` ①：门仍只认 `slog.*` 与 `observe.InitLog` 两族 records 谓词（`recordEmittingSelectors`），
  toast-only／纯 stderr 腿判"无义务"——本件**未动**，归票 133 AC#2（票面已登记）。
- `next=` ②：`rulingFiles` 这把第二把尺只核"文件里有没有那枚标记词"，不核句子内容；
  把裁决句改写成"另一句没有标记词的话"仍然红（A0 已量），**改出同义句而不带标记词**则是红——
  真正没被覆盖的是"标记词写在一个不被任何腿分发的函数 doc 上"这一形，反向那条红认得它（prose, not a decision），
  但"两个文件各写半句"这种拼接仍未测。留 `next=`。
- `next=` ③：`cmd/wisp` 只在 windows leg 有分母（`scripts/wisp-cli-tests.sh`）；Linux 原生那扇门今天**如实红**（followup-2 §1），
  不是恒绿，但也不是覆盖——归票面 `next=` ③，本件未移。
- 未验证：`0xc000013a` 争用早死本轮零命中（测前 `gh run list --limit 5` 只见 completed、`runs?status=in_progress` 计数 0）；
  墙钟只作旁证（本机比验收方慢 10-20% 那一档，本档 `ok` 103.2s/104.8s 与基线 104.8s 同量级）。
