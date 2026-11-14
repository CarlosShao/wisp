# 票 448 — 446 的「牙」只在 `probes/` 里：把判定改回英文句子比较，**仓里没有任何测试会红**

**状态**：待派（小活，但这是「回归闸 vs 一次性证明」的分界）
**来路**：独立审计腿 `fc580501b`（`docs/evidence/s1/audit-20261114-wave2.md` §四）现跑查出：`internal/tools` 里**零测试**提到 `errnoPresent`、**零测试**调那两个调用点。票 446 的 AC#3 字面（「一发能红的检查」）**是满足的**，缺的是**牙留下来**。

## 为什么值得单独立票
「我证明过它能红」与「它以后不能悄悄变绿」是两件事。446 的证明用的是 `-overlay` 挂进包的**临时**探针，跑完就没了 ⇒ 下一次有人（或某条腿）把这行改回 `== "The operation completed successfully."`，`go test ./internal/tools` **照样全绿**，而那个 bug 在非 en-US 机器上是真的。⇒ 这类洞正是「绿得像门」的典型，必须补一枚**常驻**的检查。

## 完成判据（AC）
- [x] **AC#1 常驻回归闸**：在仓里（`internal/tools/*_test.go` 或 `tools/d22scan` 的静态闸，二选一，**你定并写理由**）落一枚**常驻**检查，使「判定退回英文句子比较」或「调用点漏用 `errnoPresent`」**都会让仓里已接线的那套跑红**。
  - 判据：**把这行改回去 ⇒ 该检查红并点名**；还原 ⇒ 绿。⛔ `-overlay`、⛔ 只放 `probes/`。
- [x] **AC#2 两个方向都要红**：①判定函数本体退回句子比较 ⇒ 红；②调用点绕过 `errnoPresent` 直接用旧比较 ⇒ 红（静态闸可，但要在证据里说明它为什么能抓到）。
- [x] **AC#3 不引入 env 依赖**：检查要在**任何语言**的 Windows 上都成立（⛔ 依赖 `GetUserDefaultUILanguage` 的返回值），可用 `syscall.Errno(0)` 的字面值／本地化样例做输入。
- [x] **AC#4 门禁**：`go vet ./internal/tools` rc=0；`go test ./internal/tools` rc=0 且**四数给出**；`gofumpt -l` 空（⚠ `gofumpt.exe` 在盘上：`D:\work\base\gopath\bin\gofumpt.exe`）；`d22scan` clean；path-length GREEN。
- [ ] **AC#5 CI**：`test-core` 与 `test-windows` 各一发**含本修复且完成**的运行（run 号＋`merge-base`＋完成时刻三项齐全）。

## 边界
写面：`internal/tools/*_test.go`（新件）或 `tools/d22scan/**`；`docs/evidence/s1/448-*`；`.scratch/wisp/probes/448/**`；本票 Progress log（只追加）。
⛔ 改 `internal/tools/system_backend_windows.go` 的判定逻辑本身（446 已修好，只补牙）；⛔ `cmd/wisp/**`、`internal/panel/**`、`internal/llm/**`；⛔ 阈值/golden/契约；⛔ add -A/--amend/reset/rebase/stash/checkout .；⛔push。

## 撤销口令
「票 448 撤」＝删本票＋回退其写腿 commit（新 commit 反改，不 reset/amend）。

---

## Progress log（写腿，2026-11-14 18:5x–19:1x +08）

- **选点**：`internal/tools/*_test.go`（**不是** d22scan 静态闸）。四条理由：这条不变量是包内的（`errnoPresent` 未导出、出包无调用者）；往 d22scan 加 ban 是动 `PLAN.md` D22 的冻结清单（需业主批准），不该为一个小活背；`internal/tools` 在 `scripts/portable-tests.sh` 的 core scope 里 ⇒ 无 tag 的那一半在两个 CI job（test-core/test-windows）都跑，d22scan 只在 lint；新 ban 还要配 `selftestsamples.go`＋`scan_test.go` 的钉子，而红本来就能从 M1/M2 拿到。全文＝新件头注释 ＋ 证据 `docs/evidence/s1/448-errno-teeth-resident.md` §1。
- **落点**：`internal/tools/ticket448_errno_gate_test.go`（无 build tag，跨 GOOS；闸1 判定体、闸2 调用点、闸3 文本兜底 ＋ 三段种子反证）、`internal/tools/ticket448_errno_gate_windows_test.go`（`//go:build windows`；行为闸 10 行，其中 4 行是结构字段标出的「与语言无关的牙」）。
- **两个方向都红过**（日志 `.scratch/wisp/probes/448/`）：
  - **M1 判定体退回句子比较** → 闸1 点 `system_backend_windows.go:269`（两句）＋`:268`，闸3 点 `:269`，行为闸 **4/10 行红**（localized / wrapped / wrapped-twice / lookalike）。同一变异下 `errno-zero-raw` 在本机**仍是绿**（en-US 消息表恰好对得上）⇒ 那两行不能当牙。
  - **M2 两处调用点退回内联句子比较** → 闸2 点 `:297`（realEnumWindows）与 `:439`（sendInputs），闸3 红，**行为闸 10/10 全绿**。这一格就是静态牙不可替代的理由：行为闸测的是「判定函数对不对」，而 M2 让判定函数正确地没人用。
  - **M2b 调用点内联 `errors.Is` 绕过（无文本）** → 只有闸2 红（`:439`），闸3 绿 ⇒ 闸2 抓的是结构，不是一句英文的黑名单。
- **446 的判定逻辑零字节改动**：开工前／每次还原后 sha256 都是 `8CAD968D67AE4F3D4795CDBDB136A6A28E48AB93B44A76FEF716D61026738845`，`git diff --name-only` 空。变异前的原件 Copy-Item 到 `probes/448/system_backend_windows.go.pristine`（不用 `git checkout`／`stash`）。
- **门禁四数**：`go test ./internal/tools` rc=0 **PASS=431 FAIL=0 SKIP=1**（那 1 枚 SKIP 是本票之前就有的 `Test23LiveClipboardRoundTripThroughRealOS`，要 `WISP_T23_LIVE_CLIPBOARD=1`；本票新增 7 个条目 0 SKIP）＋ 包行 `ok … 34.556s`；`go vet ./internal/tools` rc=0；`GOOS=linux go vet ./internal/tools` rc=0；`gofumpt -l` 空；`d22scan` rc=0 clean（ban #8 internal/ 1016 个 Go 文件含新件）；path-length **VERDICT GREEN** rc=0。
- **AC#5 留白**：⛔ 本腿不 push ⇒ `test-core`／`test-windows` 两发含本 commit 且完成的运行，由编排者推送后取（run 号＋`merge-base`＋完成时刻）。
---

## 编排者收件（2026-11-14 19:1x，非实现者）：AC#1–#4 翻，AC#5 等 CI

- **尺1**：`016426b3c` 在；`git show --name-only` 只有 2 枚新测试件（`internal/tools/ticket448_errno_gate_test.go`、`ticket448_errno_gate_windows_test.go`）＋ 证据 ＋ probes ＋ 本票 ⇒ **产码零改动**。
- **我自己核「446 的产码没被动」**：`git hash-object internal/tools/system_backend_windows.go` = **`da42976e6…`**，与 `be0f599f3:` 那版、与 `HEAD:` 那版**三者同串**。
- **尺2**：`docs/evidence/s1/448-errno-teeth-resident.md` **14671 B**；`probes/448` **12 枚**，含 `10-M1…RED`、`20-M2…RED`、`21-M2b…RED`、`30/31-post-restore-GREEN`、`60-package-tests`、`61-linux-compile`，以及原件 `.pristine`。
- **尺3 我自己复跑**：`go test ./internal/tools -count=1 -v` ⇒ **RC=0**，**`RUN=584 PASS=583 FAIL=0 SKIP=1`**，`ok … 30.806s`；5 枚 `Test448*` 全 PASS；唯一 SKIP ＝ `Test23LiveClipboardRoundTripThroughRealOS`（本票之前就有，要 `WISP_T23_LIVE_CLIPBOARD=1`）。⚠ 我的 PASS=583 与腿报的 431 **口径不同**：我数到了子测试，腿数的是顶层 —— 两个都对，但引用时必须说清是哪个口径。
- **我认可它的选型**：闸落包内测试、**不往 `d22scan` 加 ban** —— 因为 d22scan 的 ban 清单是 `PLAN.md` D22 冻结清单的转写，**加 ban #11 等于动契约**（要机主批），为一个小符号把一次人工批准拖进小活不值。它还给了覆盖率理由：`internal/tools` 在 `scripts/portable-tests.sh` 的 core scope（`:385`），无 tag 那半个在 **test-core(ubuntu) 与 test-windows 都跑**；d22scan 只在 lint。
- **它最值钱的一条读数**：M1（判定退回句子比较）下，`errno-zero-raw` 在**本机仍是绿**（en-US 表恰好对得上）⇒ 那两行**不能当牙**。这正是本票 AC#3 要求的「不依赖本机语言」的机器判据（结构字段 `localeIndependentTooth`，门里断言 ≥3，现 4）。M2（调用点回退）下**行为闸 10/10 全绿、静态闸红** ⇒ 这就是静态牙必须存在的证明。
- **AC#5 留白**：等我推后 CI 的 `test-core`/`test-windows` 各一发**含 `016426b3c` 且完成**的运行（run 号＋`merge-base`＋完成时刻）。
