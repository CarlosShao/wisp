# 票 314 — cli 终点锚（commit `832e04489`）独立审计（非实现者）

- 审计腿：314 终点锚审计腿（只读产品面；⛔ 非实现者、⛔ 未翻任何 AC 框）
- 审计对象：`832e04489e35e95dd5112d2c64f095377b61bf33`（2026-11-14 09:21:00 +08）
- 审计起手 HEAD：`7bdf09042482b582a98071143f5cf6ec75512407`；落盘时 HEAD：`ebe98e6f77c70f4529da9fd0988228733c1210d4`（另一腿在此期间提交）
- 脚本同一性：`git diff --stat 832e04489 HEAD -- scripts/portable-tests.sh` 为空（`diff_rc=0`）⇒ 本审计跑的就是该 commit 的脚本字节，不是后来改过的版本
- 读数根：`.scratch/wisp/probes/314-anchor-audit/**`（只建不删；`.out` 命中 `.gitignore:8:*.out`，未入本 commit）
- 结论：**两条声称均通过**（甲 通过；乙 通过，边界见 §7）

---

## 0. 结论逐条（通过 / 不通过）

| 复核项 | 结论 | 一句话读数 |
|---|---|---|
| **甲**：早停时响亮变红并点名（跑了 / 分母 / 缺的锚 / 最后一行），而 carrier 仍打 `ok`、既有 GUARD 也放行 | **通过** | `rc=1`；`=== RUN=1` / `roster=619` / 锚名 / `last line: ok ... 3.00s` 四样齐；同一发拆掉断言 `rc=0` |
| **乙**：跑到锚时不误报（零锚告警） | **通过（小全集＋合成；本审计未跑整包）** | 最后两枚 `RUN=2 PASS=2`、锚 `--- PASS` 命中、`rc=0`、无 `END ANCHOR ... did NOT reach` |
| 锚是「最后一枚枚举到的用例的 `--- PASS/FAIL/SKIP` 裁决行」，不是 `=== RUN` | **通过** | 正则见 §1.2；合成 `=== RUN` 单行 `NOMATCH`；stub 只打锚的 `=== RUN` 仍 `rc=139` 红 |
| 没有顺手放宽 / 吞掉 | **通过** | 脚本 diff **+131 / 0 删**；运行、rc、锚判定上没有 `\|\| true`；唯一的 rc 写是 `0 → 1` 抬高 |
| rc 分类：被外部砍（124/137/128+n、`panic: test timed out`）与「测试失败」分开说 | **通过** | 124/137/139/143 与 panic 均 `KILLED, NOT FAILED`；rc=1 无签名 → `EARLY STOP / HANG`，四数保留 `FAIL=1` |
| 临时改动逐字节还原 | **通过** | SHA256 回到 `56E5504A...`；`git diff` 空；六条 grep 全 0；`bash -n` rc=0 |

---

## 1. 它加了什么（逐行核对）

`git show 832e04489 --numstat`（逐字）：

```text
2	0	.scratch/wisp/issues/314-cmdwisp-suite-process-death-in-window-zone.md
216	0	docs/evidence/s1/314-end-anchor.md
131	0	scripts/portable-tests.sh
```

脚本只有 **+131 / 0 删**，新增两段：

### 1.1 跑前：分母与终点锚（`scripts/portable-tests.sh:997-1066`）

```sh
grep -E '^Test[A-Za-z0-9_]+$' "$cli_roster" >"${cli_roster}.names" || true
cli_roster_count=$(grep -c . "${cli_roster}.names" || true)
cli_anchor_name=$(grep . "${cli_roster}.names" | tail -n 1 || true)
```

- 分母来自 `go test -list '.*' "${scope[@]}"`，只留裸 `Test` 名，滤掉 Go 自己的 `PASS` / `ok` 尾行；锚＝**最后枚举到的那一枚**。
- 本机重跑 `go test -list '.*' ./cmd/wisp/`（`PATH` 前置 `third_party/sherpa-onnx`）读数：输出 **621 行**，其中裸 `Test` 名 **619 枚**；最后五枚逐字：

```text
Test16r1PremortemMicPermissionDeniedLiveProbe
Test16r1PremortemModelSha256CorruptRefusesTheHandOff
TestZZ315A1RootLoopFirstBridgeCall
TestZZ315A1SiblingChildFirstBridgeCall
TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask
ok  	github.com/CarlosShao/wisp/cmd/wisp	0.673s
```

- 锚＝`TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`，与实现者声称一致；分母 619 与实现者声称一致。

### 1.2 跑后：只认 verdict 行（`scripts/portable-tests.sh:1218-1276`）

```sh
cli_anchor_re=$(printf '%s' "$cli_anchor_name" | sed 's/[][\\^$.*+?(){}|]/\\&/g')
if ! grep -Eq "^--- (PASS|FAIL|SKIP): ${cli_anchor_re}( |\$)" "$capture"; then
```

- 认的是 `--- PASS|FAIL|SKIP: <anchor>`，**不是** `=== RUN`；`^---` 同时排掉缩进的子测试 verdict。
- rc 分类在 `:1241-1253`；`rc=0` 但缺锚强制抬到 1 在 `:1274`。
- 这一段整体在 `if [ "$mode" = cli ]` 内，非 cli 档一字未动。

---

## 2. 正控（甲）：真造早停，红并点名

**造法**：只临时改 cli 分支的 carrier 调用，加 `-test.run '^TestAC1AlwaysBranchDoesNotRevertAHandEditedKey$'`（只跑 1 枚 ⇒ 必然早停）；脚本其余字节不动。命令形：

```text
PATH=<repo>/third_party/sherpa-onnx:$PATH bash scripts/portable-tests.sh --scope=cli
```

### 2.1 断言在：`rc=1`，四样齐（`posA.out`，逐字）

```text
--- PASS: TestAC1AlwaysBranchDoesNotRevertAHandEditedKey (2.84s)
PASS
ok  github.com/CarlosShao/wisp/cmd/wisp	3.00s
portable-tests.sh: four numbers (all from -v output): === RUN=1  --- PASS=1  --- FAIL=0  --- SKIP=0
portable-tests.sh:   ok (own line)  github.com/CarlosShao/wisp/cmd/wisp
portable-tests.sh: END ANCHOR (ticket 314) - the cli run did NOT reach its last enumerated case.
portable-tests.sh:   anchor    : a '--- PASS|FAIL|SKIP: TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask' verdict line
portable-tests.sh:   missing   : the anchor above never appeared in the -v output
portable-tests.sh:   ran       : === RUN=1  --- PASS=1  --- FAIL=0  --- SKIP=0
portable-tests.sh:   roster    : 619 top-level case(s) enumerated before the run
portable-tests.sh:   last line : ok  github.com/CarlosShao/wisp/cmd/wisp	3.00s
portable-tests.sh:   verdict   : EARLY STOP / HANG - the run ended before its last case and no external
portable-tests.sh:               kill signature was found. The four numbers above are NOT a coverage
portable-tests.sh:               denominator (ticket 314 / ledger A976).
portable-tests.sh:   Do NOT read this run's four numbers as a denominator, and do NOT report a killed run
portable-tests.sh:               as a test failure. Re-run under the same conditions; if it stops again, the stopping
portable-tests.sh:               case is the last '=== RUN' name printed above.
portable-tests.sh: strict runner exited 1 for scope=[./cmd/wisp/]
```

**逐条对齐派单**：

| 派单要求 | 这一发的逐字读数 |
|---|---|
| 跑了多少 | `=== RUN=1`（`ran` 行） |
| 分母多少 | `roster=619` |
| 缺的锚 | `TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask` |
| 最后一行 | `ok  github.com/CarlosShao/wisp/cmd/wisp	3.00s` |
| carrier 自己仍打 `ok` | `ok  github.com/CarlosShao/wisp/cmd/wisp	3.00s` |
| 既有 GUARD 也放行 | `portable-tests.sh:   ok (own line)  github.com/CarlosShao/wisp/cmd/wisp` |
| 红且 rc=1 | `POSA_RC=1`；末行 `strict runner exited 1` |

### 2.2 拆掉断言：同一发 `rc=0`（`posB.out`）

只把 `:1228` 的 `if [ "$mode" = cli ]; then` 临时改为 `if false; then`，carrier 的唯一一枚早停照旧。读数：

```text
--- PASS: TestAC1AlwaysBranchDoesNotRevertAHandEditedKey (4.18s)
PASS
ok  github.com/CarlosShao/wisp/cmd/wisp	5.00s
portable-tests.sh: four numbers (all from -v output): === RUN=1  --- PASS=1  --- FAIL=0  --- SKIP=0
portable-tests.sh:   ok (own line)  github.com/CarlosShao/wisp/cmd/wisp
POSB_RC=0
```

`posB.out` 里 `END ANCHOR` 只出现 **1 次**，且是跑前的信息行（`roster=619 ... the run must reach ...`），没有 `END ANCHOR ... did NOT reach` 的判决段。⇒ 同一发早停，没有这条断言就真的是绿的；红确实来自这条断言，不是别的门。

---

## 3. 负控（乙）：跑到锚时不误报

本审计**没有**跑整包 `go test ./cmd/wisp`（派单也要求不需要跑整包）。用仓里既有的一枚真二进制做「小全集」：只跑最后两枚，让最后一枚的 verdict 真的出现。

**命令形**（只临时改 carrier 调用，其它字节不动）：

```text
-test.run '^(TestZZ315A1SiblingChildFirstBridgeCall|TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask)$'
```

`negC.out` 逐字：

```text
=== RUN   TestZZ315A1SiblingChildFirstBridgeCall
--- PASS: TestZZ315A1SiblingChildFirstBridgeCall (0.64s)
=== RUN   TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask
--- PASS: TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask (0.44s)
ok  	github.com/CarlosShao/wisp/cmd/wisp	1.00s
portable-tests.sh: four numbers (all from -v output): === RUN=2  --- PASS=2  --- FAIL=0  --- SKIP=0
portable-tests.sh:   ok (own line)  github.com/CarlosShao/wisp/cmd/wisp
NEGC_RC=0
```

- 锚 verdict 命中，**没有任何** `END ANCHOR ... did NOT reach` 判决段（只有跑前信息行）。
- 跑序＝`-list` 末两枚的顺序（sibling 在前、anchor 在后），与「最后枚举名＝最后跑到」一致。
- **边界**：这是「小全集」负控，不是全量绿；它证明的是断言在锚 verdict 出现时保持静默，代码路径与全量跑完完全相同（该断言只 grep capture 里的锚 verdict）。全量绿仍需要一次真正的整包完成跑，本审计不冒充。

---

## 4. 锚判定宽度：`=== RUN` 不算，只有 verdict 算

用脚本里的同一条正则做合成 capture 矩阵（`logic-*.txt`）：

```text
logic-runonly=NOMATCH
logic-verdict-pass=MATCH
logic-indented-subtest=NOMATCH
logic-prefix-extra=NOMATCH
logic-verdict-skip=MATCH
logic-verdict-fail=MATCH
```

对应：

- 只有 `=== RUN   <anchor>`、没有 verdict ⇒ **NOMATCH**（不会把「死在最后一枚里」读成到达）；
- `--- PASS/FAIL/SKIP: <anchor>` ⇒ MATCH；
- 缩进的子测试 `    --- PASS: <anchor>` ⇒ NOMATCH（子测试不能给顶层锚作证）；
- `--- PASS: <anchor>Extra` ⇒ NOMATCH（前缀不能冒充）；
- 只有 `--- PASS|FAIL|SKIP` 三种裁决行算到达。

再用 stub carrier 做一次端到端「死在最后一枚里」：stub 只打 `=== RUN   <anchor>`、不打 verdict、`exit 139`。真实脚本的读数（`classify-anchor-run-only.out`）：

```text
=== RUN   TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask
portable-tests.sh: END ANCHOR (ticket 314) - the cli run did NOT reach its last enumerated case.
portable-tests.sh:   verdict   : KILLED, NOT FAILED - the run was cut by an external clock or signal
portable-tests.sh: strict runner exited 139 for scope=[./cmd/wisp/]
```

⇒ 锚行出现了，仍然红；这条断言确实不是只查 `=== RUN`。

---

## 5. 没有顺手放宽 / 吞掉 + rc 分类

### 5.1 放宽 / 吞掉

- 脚本 diff **+131 / 0 删**；`cmd/**`、`internal/**`、carrier、CI、阈值、golden、ledger 行零改动（commit 的另外两个文件只是本票 Progress log 与证据件）。
- 新增代码里的 `|| true` 只有 4 处，全部在**只读的名单/末行读取**上：

```text
grep -E '^Test[A-Za-z0-9_]+$' ... || true
cli_roster_count=$(grep -c . ... || true)
cli_anchor_name=$(grep . ... | tail -n 1 || true)
cli_last_line=$(tail -n 1 "$capture" || true)
```

- **没有**在 carrier 运行、rc 捕获、锚判定上写 `|| true`；唯一的 rc 写是 `if [ "$rc" -eq 0 ]; then rc=1; fi`（只抬高，不吞掉任何非零 rc）。
- 断言只在 `mode=cli` 生效；`speech`、`core`、`windows`、`winsec` 等档不经过这一段。

### 5.2 rc 分类（受控 stub，逐个真跑）

stub 只产生一段短 capture（非锚的 `=== RUN`、包级结果行、可选 panic），并返回指定 rc；真实 `portable-tests.sh` 的判决逐字：

| 注入 | 脚本 rc | `verdict` 逐字 |
|---|---|---|
| `rc=139` | 139 | `KILLED, NOT FAILED - ... (rc=139 (SIGSEGV / 0xC0000005 - whole process down)). This is NOT a test failure...` |
| `rc=124` | 124 | `KILLED, NOT FAILED - ... (rc=124 (external timeout(1))). This is NOT a test failure...` |
| `rc=137` | 137 | `KILLED, NOT FAILED - ... (rc=137 (SIGKILL)). This is NOT a test failure...` |
| `rc=143` | 143 | `KILLED, NOT FAILED - ... (rc=143 (128+signal 15 - killed by a signal)). This is NOT a test failure...` |
| `panic: test timed out` + `rc=2` | 2 | `KILLED, NOT FAILED - ... (panic: test timed out (Go -test.timeout)). This is NOT a test failure...` |
| `rc=1` + `--- FAIL`、无 kill 签名 | 1 | `EARLY STOP / HANG - the run ended before its last case and no external kill signature was found.`；四数保留 `=== RUN=1 --- PASS=0 --- FAIL=1 --- SKIP=0` |

- 124/137/139/143 与 panic 全部归「被砍」，且明说不是测试失败；
- rc=1 无签名归「早停 / 挂起」，四数里的 `FAIL=1` 照常打印，不把早停包装成测试失败；
- 锚 verdict 出现的真实完成流（§3）不进入这个判决段，rc=0 零告警。
- **措辞边界（如实）**：这条断言没有第三个「TEST FAILED」标签；真正的测试失败如果跑到锚，靠既有的 `FAIL` 行/四数报，不由这条断言改色。这是分工，不是吞掉。

---

## 6. 临时改动的逐字节还原证明

临时改动共 4 类（只动 `scripts/portable-tests.sh`，动前先备份到 `portable-tests.sh.orig`）：

1. cli carrier 调用加 `-test.run`（正控 A/B）；
2. `:1228` 的 `if [ "$mode" = cli ]; then` 临时改 `if false; then`（正控 B，拆断言）；
3. carrier 调用改成最后两枚的 `-test.run`（负控 C）；
4. carrier 调用改成 `bash "$AUDIT_STUB"`（rc 分类 stub）。

还原用**动前备份整文件覆盖**，不是 `git checkout`。还原后读数：

```text
restored_sha256=56E5504A774E3C3095337326FC65B481EC60D34DB24474AF2880FF422F6224F9
orig_sha256=56E5504A774E3C3095337326FC65B481EC60D34DB24474AF2880FF422F6224F9
git diff --stat HEAD -- scripts/portable-tests.sh  → （空），diff_rc=0
git status --porcelain -- scripts/portable-tests.sh → （空）
```

六个 grep 计数（逐字）：

```text
-test.run=0
AUDIT_STUB=0
TestAC1AlwaysBranchDoesNotRevertAHandEditedKey=0
TestZZ315A1SiblingChildFirstBridgeCall=0
if false; then=0
AUDIT_=0
```

语法门：

```text
bash -n scripts/portable-tests.sh
BASH_N_RC=0
```

⇒ 临时注入的字符串一个不剩，脚本字节回到审计对象的状态（SHA256 相同）。

---

## 7. 边界与如实说明

1. **没有跑整包**：本审计按派单用 `-test.run` 小选集与 stub，未跑 `go test ./cmd/wisp` 全量；因此「乙」是**小全集**负控 + 合成矩阵，不是本审计自己产出的全量绿。全量绿仍应由产品面/编排者的一发完整跑提供（本审计不冒充）。
2. **锚是单边的**：它证明「最后一枚枚举到的用例被到达」，不证明「每一枚前面的用例都跑了」；脚本自己在注释里明确 `ran` 与 `roster` **只报不等**。生产路径（`portable-tests.sh` 不向 carrier 传 `-test.run`）下顺序固定，早停不可能跳过前缀又到达末枚；但调用者若自己注入 `-test.run` 选中锚，确实能让这条断言静默——这是「终点锚」的固有边界，不是本次实现的额外放宽。
3. **名单过滤是 ASCII**：`^Test[A-Za-z0-9_]+$` 与脚本既有 ledger 检查同形；当前 `cmd/wisp` 的 652 条 `Test/Benchmark/Fuzz/Example` 声明里非 ASCII 名字为 **0**，所以今天无影响。将来若出现非 ASCII 顶层测试名，会被这条过滤器排除在名单外（潜在漏锚），实现者未处理；本审计记为残余，不构成当前两条声称的失败。
4. **并发环境干扰（如实）**：10:37 左右另一条腿正在写未跟踪的 `internal/panel/error_copy_216.go`（当时为空文件），导致两次 rc 分类跑在 ledger 预检就以 `expected 'package', found 'EOF'` 退出（`classify-panic2.out`、`classify-rc1fail.out` 保留为失败读数）；等该文件可编译后重跑，成功读数在 `classify-panic2b.out`、`classify-rc1failb.out`。这是环境噪声，不是本实现的问题，也未改动任何 `internal/**` 字节。
5. **没有翻 AC 框、没有 push、没有动产品码/测试/阈值/golden**；本件只审计 `scripts/portable-tests.sh` 的这一笔。

---

## 8. 撤销口令

回「撤 314 终点锚审计」⇒ 删本件 + 删票 314 Progress log 的那一行；脚本一字节未动（SHA256 `56E5504A...`），无需还原。
