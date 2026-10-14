# 236-r6 — 票 236：r4 钉内两条注释级小瑕疵的清扫（写腿／实现者）

> 判语归后续非实现者腿；⛔ 本腿未翻任何框、未动 `Status:`／票名。
> 落点＝**只改注释行**：`internal/tools/taskid_source_236r4_teeth_test.go`（其余零文件重叠）。
> 起手锚 `6605cb6f`（13:21:56）；执行窗口内 HEAD 移到 `a4590ac2`（13:27 编排者 A922 文档笔：stat＝5 件
> `.scratch/wisp/issues/*`＋`pending-and-issues.md`＋orch txt，**零 Go 文件**，与本腿零重叠）。

## §0 结论先行

- **瑕疵①已按现读改准**：熵支行号 `1132-1137`→`1155-1161`、nibbles `1139-1140`→`1162-1163`
  （loop.go blob `0d53277d`＝`236-v2` 记录同串 ⇒ **未再漂**，与 v2 现读逐字一致）。
- **瑕疵②注释已改成真来源**：12 枚里 **10 枚**在两 golden 目录（注释原引的那把 grep 现读 **33 枚唯一 id**）；
  另 **2 枚**在内联 wire 载荷：`call_w1`＝`internal/tools/loop_approval_test.go:288`、
  `chatcmpl-fake`＝`cmd/wisp/secret_test.go:1133`。
- diff＝`+9/−6` 全注释行，**非注释改动行＝0**；断言／锁逻辑／id 集合**一字未动**。
- 门禁全绿：`Test236` rc=0（前/后各一发）／`-race` rc=0／gofumpt 空（前/后）／build rc=0（前/后）／
  d22scan clean rc=0／长度门 GREEN（over=59／roster=59／not-in-roster=0）。

## §1 瑕疵① 现读（熵支与 nibbles）

| 位置（pin 行） | 原文（旧） | 改为（现读真身） | 真身逐字 |
|---|---|---|---|
| pin:53（文件头注释） | `(loop.go:1132-1137)` | `(loop.go:1155-1161)` | `:1155 if _, err := rand.Read(b[:]); err != nil {` … `:1161 }` |
| pin:107-108（`r236r4MintedShape` 文档注释） | `loop.go:1139-1140` | `loop.go:1162-1163` | `:1162 b[6] = (b[6] & 0x0f) | 0x40` / `:1163 b[8] = (b[8] & 0x3f) | 0x80` |

判据＝本腿 13:27 现读 `awk NR>=1153 && NR<=1165 loop.go` 与 `git hash-object` 双向核对（logs/01-reprojection.txt）。

## §2 瑕疵② 现读（12 枚 id 的真来源）

原注释称 12 枚均从 `internal/agent/testdata/golden`＋`internal/llm/testdata/golden`
（`grep -rhoE '"id":"[^"]+"'`）现读。**现跑该 grep＝33 枚唯一 id**，其中 12 枚的命中：

| id | 在两 golden 目录命中 | 真身（不在目录的 2 枚） |
|---|---|---|
| call_e1 / call_b2 / call_a1 / call_x1 / fc_1 / fc_2 / rs_1 / msg_mock1 / msg_out1 / resp_mock1 | 1/3/2/2/3/2/2/11/12/20 | — |
| call_w1 | **0** | `internal/tools/loop_approval_test.go:288`（内联 SSE 载荷） |
| chatcmpl-fake | **0** | `cmd/wisp/secret_test.go:1133`（httptest 内联 JSON） |

改后注释（pin:79-85）＝"Ten of the twelve were re-read …（33 unique ids）；the other two sit in inline wire
payloads elsewhere in this repo, not in those directories：`call_w1`(…:288)、`chatcmpl-fake`(…:1133)"。
⛔ id 列表（pin:86-89 字面量）与断言本体零接触；两枚仍过不了锁（7 与 13 字符 ≠ 36）、也不等于真铸值。

## §3 判据逐条

| # | 判据 | 读数 | 出处 |
|---|---|---|---|
| ① | `git diff` 只动注释行 | 非注释改动行＝**0**；numstat `9 6`；sha256 `47fb441d…`→`7e23ec32…`；275→278 行 | logs/02 |
| ② | `-run 'Test236' ./internal/tools/ -count=1` | 前 rc=0（1.989s）／后 rc=0（3.361s；`Test236R4…` PASS） | logs/03 |
| ③ | `-race` 一发 | rc=0（4.650s） | logs/03 |
| ④ | gofumpt 该件空 | `gofumpt -l` 前/后皆空 rc=0 | logs/03 |
| ⑤ | 断言/锁逻辑/id 集合一字不动 | diff 面全为 `//` 行；断言体、12 枚 id 字面量零接触 | logs/02 |
| 补 | build／d22scan／长度门 | build rc=0（前/后）／d22scan clean rc=0（正面控制 PASS=35 FAIL=0）／GREEN 59=59 | logs/03 |

## §4 窗口逐报（他腿在飞）

- 13:21:56 起手锚 HEAD `6605cb6f`；pin sha256 `47fb441d…`（＝r4 声明＝v2 复量）；loop.go blob `0d53277d`；
  run.go blob `154279f7`；ticket224 blob `ebe020dd`。
- 13:22-13:25 门禁两轮（前/后）；13:27 reprojection 时 HEAD 已在编排者笔 `a4590ac2`（见件头）。
- 全程 `internal/tools/task_output_pointer_notice_test.go` 挂 `M`＝**174-r6 在飞件**（同包不同件）：
  本腿两发测试均带它一起编译，读数已具名；未读未碰。
- CI `38508311376` 在跑——本件未引任何 CI 读数。

## §5 具名顶回（登记外漂移，**未动**，留裁定）

1. **run.go 消费行号族**（pin 注释 `:17-19`／`:136`，断言字符串内 `:205-206`／`:244`）：数字写 r4 时代
   `1105/1110/1124/1135/1145`，在 `f88d98ae` 上逐字 ✓（236-v2 已记"引用者须重跑"）。HEAD 现读已漂为
   `1137`(MarkRoot)／`1139-1142`(Backfill)／`1156-1160`(StartTaskLog)／`1167`(FinishTaskLog)／`1177`(状态行)
   ——315-r1 `+29` 与 279-r1 `+3` 所致。⛔ 未动：同一组号在 `t.Errorf/Fatalf` 字符串里各有一份（非注释行，
   判据①⑤禁碰），只改注释半边会让文件自相矛盾。⇒ 需**整族一并处置**的裁定。
2. **`ticket224_assembly_test.go:433`**（pin:28 注释行）：r4 blob 上＝`Valid()` 断言处；252-r3 `d3cccfda` `+9` 行后
   现读 `:441-443`（`:442` 是 `t.Fatalf` 行）。⛔ 未动（登记外）。
3. **loop.go:676 那句引用**（pin:210-211）：自 r4 起 676 行本就含 `CorrelationID`、`CallID` 在 677 ⇒ 是 r4 的
   压缩引用、**非漂移**；实质（两格并排）仍真。⛔ 未动。

## §6 具名没做

- 未动 run.go／ticket224／任何断言字符串（理由见 §5）。
- 未复跑 internal/tools、internal/agent **整包**基线（面＝单件注释；`-run 'Test236'` 已覆盖本件族；
  整包 220/0/0 读数在 v2 窗口）。
- 未跑 CI 复读、未 push、未碰他腿脏件。

## §7 件清单（committed）与 raw（留盘未提交）

committed：`00-anchor-and-gates.md`(本件)／`logs/01-reprojection.txt`／`logs/02-diff-comment-only.txt`／
`logs/03-gates.txt`（＝pre/post 合订）。
raw（只建不删，未提交）：`logs/01-anchor.txt`、`logs/02-pre-Test236.txt`、`logs/03-pre-gofumpt-build.txt`、
`logs/04-diff-comment-only.txt`、`logs/05-post-Test236-race.txt`、`logs/06-post-gofumpt-build-length.txt`、
`logs/07-d22scan.txt`（均为上列 committed 的来源快照）。
