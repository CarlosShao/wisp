# 票 111 — 五格收口裁决表（编排者出，非实现者；2026-10-15 13:0x）

- **被裁对象**：腿 `111-r7`（落地腿；**平台 stream 超时死于 41 min**，但**盘上活完整**：5 笔 commit `a45f6776`／`9655ed02`／`3998be73`／`d5190e0e`／`33e429c1`＋`verdict.md`＋`probes/111/r7/**`；**零产码改动**——`git status --porcelain -- scripts .github` 空、五笔逐枚只带 `probes/111/r7/**`＋票尾 bullet）。
- **裁定人**：编排者（非实现者，独立性满足：本轮五格没有一枚由我实现；我**另取**读数复核，未引用它日志里的数）。
- **锚**：`33e429c1`（收笔）／起手 `457e7513`；裁定时 HEAD 含 319-r2 `2ad0aa86` 之后各笔。
- **裁定方式**：逐格判 + **编排者亲跑**（下面每格的"编排者亲跑"栏都是我自己在 2026-10-15 13:0x 现跑的读数，命令逐条写明）。

## 逐格

### AC#3（空分母守卫）＝**成立**
- 守卫本体在盘：`scripts/portable-tests.sh` 的 GUARD A（`:651` 区）／GUARD C（`:629` 区）／GUARD D（`:474`/`:487`/`:490` 区）——**编排者亲跑**：见下两条变异。
- **编排者亲跑 M1（抽包⇒红）**：`git archive HEAD` 快照 `D:/tmp/wisp111orch/m1`，把 core scope 行里的 `./internal/plugin/` 摘掉（pin 不动，1 行）⇒ `bash scripts/portable-tests.sh --scope=core`：
  `GUARD C - scope mode=core resolved to a DIFFERENT package set than the one pinned next to it. Pinned: 28, resolved: 27.` ＋ 逐字点名 `< github.com/CarlosShao/wisp/internal/plugin`，**rc=1**。
- **编排者亲跑 M2（声明空分母⇒红）**：同快照另切 `m2`，core scope 行尾追加 ` ./internal/watchdog/` ＋ core_pin 插同名（2 行）⇒ 同一命令：
  `GUARD A - these packages are declared in scope mode=core but compile NO test file at all for GOOS=windows:` ＋ `EMPTY github.com/CarlosShao/wisp/internal/watchdog` ＋ "a scope entry with no denominator cannot fail…（ticket 111 AC#3）"，**rc=1**。
- ⇒ 两个方向都有牙、报红句点名包与平台。**判：成立**（它 §2/§3 与我这两发同形；我的读数独立取得）。

### AC#5（门禁＋步序）＝**成立**
- `bash -n`：**编排者亲跑** 7 枚 shell 脚本（`portable-tests.sh`／`portable-tests-selftest.sh`／`winsec-tests.sh`／`wisp-cli-tests.sh`／`slo-freshness.sh`／`check-path-length-budget.sh`／`d22scan.sh`）逐枚 **rc=0**。
- `sh scripts/d22scan.sh`：**编排者亲跑 rc=0 clean**；各 scope `internal/=237/40/115/24/30/115/558/131`（r5 基线 228/38/85/23/30/85/514/104）⇒ **逐数不降**。
- 步序（新步在红步之后仍答话）：**编排者现读同一枚 run**（`38676514264`）的步表——`lint` job 内 step8/9/12=failure 而 step13（mockllm vet）/step14（本票 carrier self-test）=success；`test-windows` job 内 step7=failure 而 step8（winlive gate）/step9（census）=success ⇒ `if: !cancelled()` 的形状成立。
- YAML 可解析：**编排者亲跑** `python -c yaml.safe_load` 成功（jobs＝lint/test-core/test-windows/slo-smoke/slo-full/lint-frontend）。

### AC#7（`-run TestSyncRegistryProbeLive`）＝**成立**（采票面给的第二支：**当众改口径**）
- **编排者亲跑**：`git grep -c -- "-run TestSyncRegistryProbeLive" HEAD -- scripts tools .github` ⇒ **零命中（rc=1）** ⇒ "① 从未落地"属实。
- 它把"① 会被拒"的两条读数留在 `g3-ac7.md`（`runtests.sh` 下 SKIP 致命＝常红／裸 `go test` 下 skip 印 ok＝假绿，票 110 `:96`＋对抗验收 `:17` 已裁）+ ② 的新硬读数（`R-110-4` 欠的 ledger class+reason 行今天存在：run `38676514264`/job `116231259514`/step 名 `ci.yml:763` 区、`03:40:45Z` 打印）＋票尾 bullet 写死口径。⇒ 票面要求"当众改口径、不许留在原地当已做"**已满足**。
- 残余（不遮）：**"它真跑并 PASS"＝0 次**，归一台有同步记录（HKCU）的主机（self-hosted `wisp-slo` 是候选），属另票/另 runner 面。

### AC#8（包匹配前缀锚定＋阳性自证）＝**成立**
- 锚定在盘：`portable-tests.sh:874/:884/:890` 区＋`winsec-tests.sh:234/:235` 区（它逐字钉；我抽验了 `escape_re`／`result_tail` 的用法）。
- **编排者亲跑阳性自证**（它的可复跑脚本 `probes/111/r7/logs/g4-ac8-proof.sh`，shim＝仓内 `scripts/testdata/portable-tests/go`，`scripts/` 零改）：宽松 `grep -c` 数出 **2**（"名字只出现在字符串里"）、**提交匹配式数出 0**、对照真 `ok` 行命中、`winsec` guard-2 同形 `[]`；Part A 端到端（结果行抽掉＋名字只种进字符串）**rc=1** 且 `GUARD B - 1 of 27 packages in scope mode=core printed…` 逐字点名。**与它自报的数一致**（我独立重跑，非引用）。
- 边界（它 §877-883 自述、我认）：这是**输出行形状**匹配，不是协议性质。

### AC#11（winlive 编译门）＝**成立**
- ⓐ 当前读数：run `38676514264`/job `116231259514` **step 8**（`winlive compile gate (go vet -tags winlive, ticket 111 AC#11)`）**success**（编排者现读步表）＋本机 rc=0。
- ⓑ 反形自证：**编排者亲跑**（`-overlay` 仓外副本，靶 `internal/ball/live_windows_test.go` 种一行语法错）：`go vet -overlay=… -tags winlive ./internal/ball/` ⇒ `internal\ball\live_windows_test.go:5:1: expected declaration, found this`，**rc=1**；**摘 tag 同 overlay** ⇒ **rc=0**（该文件根本不被编译）。
- ⓒ 口径具名：`ci.yml:636-645` 区逐字含 "STATED AS COMPILE-ONLY SO IT IS NOT READ AS A SLIDE (ticket 111 AC#11 c)"＋"explicitly NOT a CI carrier for ticket 35's real-window readings"（编排者现读）。
- 名册：**编排者亲跑** `git grep -lE '//go:build.*winlive' HEAD -- '*.go'` ⇒ **14 枚**（cmd/wisp 9＋internal/ball 5），与它自报的 12→14 生长一致。

## 编排者顺带做的两处（本表之外的收口动作，已入本笔 commit）

1. **`ci.yml` 三处过期注释数字已改**（它具名交回、我裁"现在就改"）：`12 files`→"9 under cmd/wisp, 5 under internal/ball at 2026-10-15 - 14 in total…the count is a READING, not a pin"（另两处改成不带数字的 `those files`／`the winlive-tagged files`）⇒ **纯注释、零行为改动**（`git diff` 无一条非注释行；`yaml.safe_load` 通过）。⚠ 行号因此漂移（原 `:624/:637/:668` 区）。
2. **`ci.yml` 那段"hosted 颜色尚未观测"已闭合**：加 `CLOSED 2026-10-15 (ticket 111-r7): run 38676514264 / job 116231259514 step 8 reads success…` 并写明同日反形复证。

## 未做/不遮

- 本表**没有**把它的日志数当凭据（凡引数字处，都是我另取；同数处已注明"独立重跑一致"）。
- 未跑：`portable-tests.sh`／`winsec-tests.sh` 的**全量**（含 linux 腿）——那是 CI 的活；本表只裁五格。
- AC#7 残余（真跑并 PASS）与 ci.yml 数字曾经过期两件都在票面 log 具名，不随本表消失。
