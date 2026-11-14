# 票 314 — `scripts/portable-tests.sh` cli 分支**终点锚断言**（证据件）

- 腿：314c 写腿（脚本面；⛔ 非实现者验收、⛔ 未翻任何 AC 框）
- 派单：`.scratch/wisp/probes/orch/20261113/txt/314c-end-anchor-dispatch.txt`
- 时刻：2026-11-13 22:4x–23:0x +08
- 起手 HEAD：`843ad88aa14d00ef7912bba59a29520848845bb3`（`git rev-parse HEAD` 现读）
- 写面（最终 diff）：`scripts/portable-tests.sh` **+131 行 / 0 删**；本件（新）；票 314 Progress log 一行。
- ⛔ 未动 `cmd/**`、`internal/**`、`.github/**`、`models/**`、任何 Go 测试、任何阈值/golden/契约；⛔ 未 push。
- 临时件根：`D:\tmp\wisp\20261113\314c\`（只建不删）。

---

## 1. 落点与机制（回话第 1 问）

### 1.1 分母从哪来

`scripts/portable-tests.sh` 的 cli 分支在**跑之前**新增一段（落点＝原 `:997` 前，现 `:997-1067`）：

```sh
go test -list '.*' "${scope[@]}"          # scope == (./cmd/wisp/)
grep -E '^Test[A-Za-z0-9_]+$'              # 只留裸 Test 名，滤掉 go 自己的 PASS/ok 尾行
cli_roster_count=$(grep -c . names)        # 分母
cli_anchor_name=$(grep . names | tail -n 1) # 终点锚 = 最后枚举到的那枚
```

- **为什么用 `go test -list`**：脚本里已有的同形做法（`:851` 的 ledger universe、`:886` 的逐名 `-list` 复核），复用同一台仪器，不新造第二把尺。
- **为什么滤掉尾行**：`go test -list` 除用例名外还在 stdout 打 `PASS` 与 `ok  <pkg> <dur>`（本机实测 621 行里含这两行）；不过滤时第一版曾把锚读成 `ok  github.com/CarlosShao/wisp/cmd/wisp	0.184s`，被名字形状检查挡下（那条拒绝读数留在 `posA-first.out`，只建不删）。加 `^Test[A-Za-z0-9_]+$` 后锚稳定＝`TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`。
- **为什么"最后一个枚举名"就是终点**：Go 按文件顺序注册顶层用例、按同一顺序执行；`-list` 按同一顺序打印。本机 `-test.list '.*'` 逐字尾部＝`Test16r1Premortem*` → `TestZZ315A1RootLoopFirstBridgeCall` → `TestZZ315A1SiblingChildFirstBridgeCall` → `TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`（619 枚，`grep -n '^TestZZ'` 逐字 617/618/619）。
- **不要求 ran == 分母**：子例程会多加 `=== RUN` 行、carrier 的名册回退（ticket 111 wave 8b 的 `-test.skip`）与台账 skip 合法减枚。两个数**只报不等**；断言的是锚的**verdict 行**。

### 1.2 跑完断言什么

跑完在 GUARD B 之后（落点＝原 `:1145` 后，现 `:1218-1264`）：

```sh
grep -Eq "^--- (PASS|FAIL|SKIP): <anchor>( |$)" "$capture"
```

- 认 **verdict 行**而不是 `=== RUN` 行：进程死在最后一枚里也会打 `=== RUN`，只有 verdict 能分开"跑到头"与"死在中途"。
- 见不到 ⇒ 响亮失败，点名四样：**跑了多少**（四数）、**分母多少**、**缺的锚**、**最后一行输出**；并区分**被外部砍了**（rc=124/137、128+n 信号含 139、或 `panic: test timed out`）与**测试失败/早停**。
- rc==0 但缺锚 ⇒ 强制 `rc=1`（carrier 会在早停时打自己的 `ok` 行并退出 0，这条是关键）。

---

## 2. 正控（回话第 2 问）★真造早停，两发读数

**造法**：临时在 cli 分支的 carrier 调用后追加 `-test.run '^TestAC1AlwaysBranchDoesNotRevertAHandEditedKey$'`（只跑 1 枚 ⇒ 必然早停；该临时改动在控制结束后逐字节还原）。命令形（两发相同，只有断言开关不同）：

```
bash scripts/portable-tests.sh --scope=cli      # PATH 前置 third_party/sherpa-onnx
```

### 2.1 红那一发（断言在）— `posA.out`，rc=1

```
portable-tests.sh: END ANCHOR (ticket 314): roster=619 top-level case(s) from `go test -list '.*' ./cmd/wisp/`; the run must reach the last one and print its verdict: TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask
...
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

★ 注意倒数第 6 行：carrier 自己打了 `ok` 行、GUARD B 也认了它——**两条既有判据都放行**，只有新断言抓住了这次早停（`rc=0 → rc=1`）。

### 2.2 回绿那一发（拆掉断言，同一早停）— `posB.out`，rc=0

把 `:1218` 的 `if [ "$mode" = cli ]; then` 临时改成 `if false ; then`（只拆断言，早停照旧）：

```
portable-tests.sh: END ANCHOR (ticket 314): roster=619 top-level case(s) from `go test -list '.*' ./cmd/wisp/`; the run must reach the last one and print its verdict: TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask
portable-tests.sh: four numbers (all from -v output): === RUN=1  --- PASS=1  --- FAIL=0  --- SKIP=0
portable-tests.sh:   ok (own line)  github.com/CarlosShao/wisp/cmd/wisp
```

**rc=0，全绿，零锚告警** ⇒ 同一发早停，没有这条断言就真的是绿的。这就是 CI 上"四数不作分母"却无人过问的形状。

两发之后的还原自证：

```
grep -c 'TEMPORARY POSITIVE-CONTROL\|-test\.run' scripts/portable-tests.sh  → 0
git diff --stat -- scripts/portable-tests.sh                                → 131 insertions, 0 deletions
```

---

## 3. 负控（回话第 3 问）本机正常跑

### 3.1 第一发全量：**仍红，但红得明白**（`neg-full.out`）

命令形：`timeout 1800 bash scripts/portable-tests.sh --scope=cli`（PATH 前置 `third_party/sherpa-onnx`）。

| 读数 | 值 |
|---|---|
| rc | **139**（＝128+11 SIGSEGV；ticket 314 的"整进程倒下"签名） |
| 耗时 | **370.00s**（carrier 实测墙钟） |
| 四数 | `=== RUN=418  --- PASS=329  --- FAIL=0  --- SKIP=0` |
| 分母 | 619（`go test -list '.*' ./cmd/wisp/`） |
| 锚 | `TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask` —— **未出现** |
| 死点 | `=== RUN   Test145ArgsJSONComesFromTheBookedForensicsRow` 之后**无 verdict**，紧接 `FAIL  github.com/CarlosShao/wisp/cmd/wisp	370.00s` |
| 最后一行 | `FAIL  github.com/CarlosShao/wisp/cmd/wisp	370.00s` |

逐字（这一发跑在分类器扩展**之前**，所以当时印的是 EARLY STOP；同一发在新码下会印 KILLED，见 §3.4 的受控复验）：

```
portable-tests.sh: four numbers (all from -v output): === RUN=418  --- PASS=329  --- FAIL=0  --- SKIP=0
portable-tests.sh:   FAIL (own line) github.com/CarlosShao/wisp/cmd/wisp
portable-tests.sh: END ANCHOR (ticket 314) - the cli run did NOT reach its last enumerated case.
portable-tests.sh:   anchor    : a '--- PASS|FAIL|SKIP: TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask' verdict line
portable-tests.sh:   missing   : the anchor above never appeared in the -v output
portable-tests.sh:   ran       : === RUN=418  --- PASS=329  --- FAIL=0  --- SKIP=0
portable-tests.sh:   roster    : 619 top-level case(s) enumerated before the run
portable-tests.sh:   last line : FAIL  github.com/CarlosShao/wisp/cmd/wisp	370.00s
portable-tests.sh:   verdict   : EARLY STOP / HANG - ...
portable-tests.sh: strict runner exited 139 for scope=[./cmd/wisp/]
NEG_CTRL_RC=139
```

★ **诚实结论**：本机这一刻**给不出"全量绿"的读数**——套件自己死在 `Test145ArgsJSON...`（与 2026-11-11 那发 ac3 读数同一死点）、rc=139。这不是本改动引入的红：rc 本来就是非零、四数行本来就在、GUARD B 本来也认了 carrier 的 `FAIL` 行；**新断言买到的是**"这 418 枚不是分母、缺的是哪枚锚、最后停在哪一行"。按派单"负控必须仍绿"的字面，此格**未达成**，原因具名如上（产品面 ticket 314 本人）。

### 3.2 第二发全量：**跑到 795 枚仍被自己的 90m `-test.timeout` 砍了**（`neg-full2.out`）

命令形同上。这一发用的是**最终码**（分类器扩展之后），所以 panic 分支是真实读数、不是 rig：

| 读数 | 值 |
|---|---|
| rc | **2**（Go `-test.timeout` panic） |
| 耗时 | **36149.00s**（carrier 实测墙钟；含本机一段约 9h41m 的挂起/睡眠，`date +%s` 把挂起算进去了） |
| 四数 | `=== RUN=795  --- PASS=601  --- FAIL=2  --- SKIP=12` |
| 分母 | 619（顶层）；795 含子例程 |
| 锚 | `TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask` —— **未出现** |
| 最后一行 | `FAIL  github.com/CarlosShao/wisp/cmd/wisp	36149.00s` |

逐字：

```
portable-tests.sh: END ANCHOR (ticket 314) - the cli run did NOT reach its last enumerated case.
portable-tests.sh:   ran       : === RUN=795  --- PASS=601  --- FAIL=2  --- SKIP=12
portable-tests.sh:   roster    : 619 top-level case(s) enumerated before the run
portable-tests.sh:   last line : FAIL  github.com/CarlosShao/wisp/cmd/wisp	36149.00s
portable-tests.sh:   verdict   : KILLED, NOT FAILED - the run was cut by an external clock or signal
portable-tests.sh:               (panic: test timed out (Go -test.timeout)). This is NOT a test failure and must not be reported as one;
portable-tests.sh: strict runner exited 2 for scope=[./cmd/wisp/]
NEG2_RC=2
```

★ **panic 分类分支的真实读数**（非 rig）：`KILLED, NOT FAILED (panic: test timed out)`。
★ 这一发跑到过全部 voice 族（模型 hand-off 真跑），`PASS=601 / SKIP=12`，比第一发远得多，但 **90m 预算耗尽**仍是 ticket 314 的同一形态（"卡到超时、响亮 panic"）。

### 3.3 负控（断言不误报）＝合成一次"跑到锚"的完成流

本机两发真实全量**都给不出绿**（产品面原因，见 §3.1/§3.2），所以另做一格**合成完成流**证明这条断言在"跑到终点"时保持静默、不误报：临时把 cli 分支的 carrier 调用换成一段含**锚 verdict 行**与包级 `ok` 行的短流（跑完逐字节还原）。

| 读数 | 值 |
|---|---|
| rc | **0（绿）** |
| 四数 | `=== RUN=2  --- PASS=2  --- FAIL=0  --- SKIP=0` |
| 锚 | 命中，**零告警** |
| 包级行 | `ok (own line)  github.com/CarlosShao/wisp/cmd/wisp` |

逐字（除四数与包级行外，无任何 END ANCHOR 段落）：

```
portable-tests.sh: four numbers (all from -v output): === RUN=2  --- PASS=2  --- FAIL=0  --- SKIP=0
portable-tests.sh:   ok (own line)  github.com/CarlosShao/wisp/cmd/wisp
```

**边界写明**：这不是真实套件绿，只是"断言在锚出现时不响"的负控；本机此刻**没有**全量绿的读数可交（原因就是本票）。

### 3.4 分类器的受控复验（`classify-*.out`）

因为 §3.1 那一发跑在分类器扩展之前，扩出来的分支用**受控 rig** 复验（临时把 cli 分支的 carrier 调用换成一段短流＋指定退出码，跑完逐字节还原）：

| rig | rc | 逐字 verdict |
|---|---|---|
| `exit 139` | 139 | `verdict   : KILLED, NOT FAILED - ... (rc=139 (SIGSEGV / 0xC0000005 - whole process down)). This is NOT a test failure...` |
| `panic: test timed out after 10m0s` + `exit 2` | 2 | `verdict   : KILLED, NOT FAILED - ... (panic: test timed out (Go -test.timeout)). This is NOT a test failure...` |
| 无签名早停（§2.1，rc=1） | 1 | `verdict   : EARLY STOP / HANG - ...` |

三类分开：**外部时钟**（124）、**信号/原生崩**（128+n，含 137/139）、**Go 自己的 `-test.timeout` panic** 都归"被砍"，只有无签名的早停才写 EARLY STOP，断言红一律不当成"测试失败"。

---

## 4. 门禁（回话第 4 问）

| 门 | 命令 | 读数 |
|---|---|---|
| 语法 | `bash -n scripts/portable-tests.sh` | **rc=0** |
| D22 | `sh scripts/d22scan.sh` | **rc=0，clean - no D22 ban violations** |
| 路径长度 | `sh scripts/check-path-length-budget.sh` | **rc=0，VERDICT GREEN**（tracked=19545，over-budget=69，roster=69，not in roster=0） |

---

## 5. 边界与禁止（如实）

- 写面仅 `scripts/portable-tests.sh` ＋本件＋票 314 Progress log；⛔ 未加 `|| true`、⛔ 未放宽任何断言/阈值、⛔ 未吞 rc、⛔ 未改 carrier/ledger/CI。
- 断言只作用于 `--scope=cli`；其他档一字未动。
- ⛔ 未跑 winlive 档、⛔ 未动 `.github/**`、⛔ 未 push。
- 欠账：`go test -list` 的前置（DLL 可达性）与脚本既有的 ledger universe 调用同源——该调用今天就在同一分支里，所以本改动**没有引入新的环境依赖**；若二进制列不出名单，脚本早在 `:851` 就已拒绝。

---

## 6. 撤销口令

回「撤 314c 终点锚」⇒ `git checkout` 掉 `scripts/portable-tests.sh` 的这一笔（或 revert 该 commit），删本件；零产码、零测试、零阈值受影响。
