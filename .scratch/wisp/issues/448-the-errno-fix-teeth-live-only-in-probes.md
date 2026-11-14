# 票 448 — 446 的「牙」只在 `probes/` 里：把判定改回英文句子比较，**仓里没有任何测试会红**

**状态**：待派（小活，但这是「回归闸 vs 一次性证明」的分界）
**来路**：独立审计腿 `fc580501b`（`docs/evidence/s1/audit-20261114-wave2.md` §四）现跑查出：`internal/tools` 里**零测试**提到 `errnoPresent`、**零测试**调那两个调用点。票 446 的 AC#3 字面（「一发能红的检查」）**是满足的**，缺的是**牙留下来**。

## 为什么值得单独立票
「我证明过它能红」与「它以后不能悄悄变绿」是两件事。446 的证明用的是 `-overlay` 挂进包的**临时**探针，跑完就没了 ⇒ 下一次有人（或某条腿）把这行改回 `== "The operation completed successfully."`，`go test ./internal/tools` **照样全绿**，而那个 bug 在非 en-US 机器上是真的。⇒ 这类洞正是「绿得像门」的典型，必须补一枚**常驻**的检查。

## 完成判据（AC）
- [ ] **AC#1 常驻回归闸**：在仓里（`internal/tools/*_test.go` 或 `tools/d22scan` 的静态闸，二选一，**你定并写理由**）落一枚**常驻**检查，使「判定退回英文句子比较」或「调用点漏用 `errnoPresent`」**都会让仓里已接线的那套跑红**。
  - 判据：**把这行改回去 ⇒ 该检查红并点名**；还原 ⇒ 绿。⛔ `-overlay`、⛔ 只放 `probes/`。
- [ ] **AC#2 两个方向都要红**：①判定函数本体退回句子比较 ⇒ 红；②调用点绕过 `errnoPresent` 直接用旧比较 ⇒ 红（静态闸可，但要在证据里说明它为什么能抓到）。
- [ ] **AC#3 不引入 env 依赖**：检查要在**任何语言**的 Windows 上都成立（⛔ 依赖 `GetUserDefaultUILanguage` 的返回值），可用 `syscall.Errno(0)` 的字面值／本地化样例做输入。
- [ ] **AC#4 门禁**：`go vet ./internal/tools` rc=0；`go test ./internal/tools` rc=0 且**四数给出**；`gofumpt -l` 空（⚠ `gofumpt.exe` 在盘上：`D:\work\base\gopath\bin\gofumpt.exe`）；`d22scan` clean；path-length GREEN。
- [ ] **AC#5 CI**：`test-core` 与 `test-windows` 各一发**含本修复且完成**的运行（run 号＋`merge-base`＋完成时刻三项齐全）。

## 边界
写面：`internal/tools/*_test.go`（新件）或 `tools/d22scan/**`；`docs/evidence/s1/448-*`；`.scratch/wisp/probes/448/**`；本票 Progress log（只追加）。
⛔ 改 `internal/tools/system_backend_windows.go` 的判定逻辑本身（446 已修好，只补牙）；⛔ `cmd/wisp/**`、`internal/panel/**`、`internal/llm/**`；⛔ 阈值/golden/契约；⛔ add -A/--amend/reset/rebase/stash/checkout .；⛔push。

## 撤销口令
「票 448 撤」＝删本票＋回退其写腿 commit（新 commit 反改，不 reset/amend）。
