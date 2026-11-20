# 191-cap-pin-r1 — 票 191 补牙：把 `sessionPendingCap` 的**数值**钉死（实现腿自记，2026-11-20 23:2x）

**身份**：写腿（实现者，agent-168）。**起手锚**：`ee2475111`（该笔只含派单存档 `191-cap-pin-dispatch.txt`，`internal/tools` 无未提交改动）。裁决者≠实现者：本件是实现腿的证据件，对抗复验由编排者另派。

## 1. 审计抓到的缺口与本腿的补法
`191-audit-r1.md` 第②节：`internal/tools/shell_session.go:64` 的 `sessionPendingCap`（64MiB）只被"上限存在"钉住——既有尺 `Test191StreamChunkCapIsTheBufferCap` 的剂量（`sessionPendingCap + (1<<20)`）与上限读数（`col.bytes() > sessionPendingCap`）都**由同一常量派生**（自证型），M2a（64→128MiB）rc=0 全绿。本腿在 `internal/tools/shell_session_cap_pin_191_test.go` 落两枚钉：

| # | 用例 | 钉的是什么 | 期望从哪来 |
|---|---|---|---|
| 1 | `Test191SessionPendingCapNumberPinnedTo64MiB` | 常量**恰等于 `64 << 20`**（67108864 字节） | **字面量**；红句携带消费方清单与 64MiB 的预算同源理由（C22，PLAN.md:1378） |
| 2 | `Test191FloodRefusedPastLiteral68157440Bytes` | 灌 **67108864 + 1048576 = 68157440 字面字节**（marker 永不到）：断言①走「单发上限」拒绝路径（非 deadline）②拒因文本含字面 `67108864` ③消费端累计 ≤ 67108864 | **全部字面数字**，零派生；常量改大必红 |

**为什么是 64MiB（钉在用例 1 的注释与红句里）**：该帽回答"一发命令流得比它结束快"（`yes` 型），与 C22 预算同源（每工具 timeoutMs、单任务 token 预算同族）——失控命令至多吃掉一个 64MiB 缓冲就由帽路径回收会话并拒绝，比同发调用可花费的任务级预算低一个量级；也是面板侧引用的数（`internal/panel/terminal_carrier_191.go:55` 称它为 tool side's own 64MiB transport cap）。**改它必须连消费方一起改**：`shell_session.go` 的 sink 帽（:324）、搜索循环帽（:358）、打印该数的拒因文本（:362）＋面板措辞＋票面理由，然后更新两枚钉与本件。

## 2. 牙必真红（M2a 重演，`-run 'Test191SessionPendingCapNumberPinnedTo64MiB|Test191FloodRefusedPastLiteral68157440Bytes' -count=1 -v`）
变异：`shell_session.go:64` 改 `const sessionPendingCap = 128 << 20`（只此一处）⇒ **两枚都红**（比"至少一枚"更强），rc=1，逐字红句：

```
=== RUN   Test191SessionPendingCapNumberPinnedTo64MiB
    shell_session_cap_pin_191_test.go:47: sessionPendingCap moved off 64MiB: got 134217728 bytes, want 67108864 (64<<20). The number is C22-budget-sourced; if it moved, move its consumers with it - shell_session.go's sink guard and search-loop cap and the refusal text that prints the number, internal/panel/terminal_carrier_191.go:55's wording, the ticket rationale - then update this nail and docs/evidence/s1/191-cap-pin-r1.md
--- FAIL: Test191SessionPendingCapNumberPinnedTo64MiB (0.00s)
=== RUN   Test191FloodRefusedPastLiteral68157440Bytes
    shell_session_cap_pin_191_test.go:76: 68157440 literal bytes must take the cap refusal (the constant was 67108864): 这条命令没在 C22 预算内回来，而且会话流里没有出现本发的结束标记——这是失去同步，不是「命令失败」：整条会话进程树已回收，本发不报退出码（拿不到真值），本发已积的输出也不归属给任何结果。下一发 shell.session 会开一条全新会话（工作目录与环境从头开始）。
--- FAIL: Test191FloodRefusedPastLiteral68157440Bytes (60.00s)
FAIL
FAIL	github.com/CarlosShao/wisp/internal/tools	60.162s
```

用例 2 的红因形状＝帽路径随常量一起变宽、68157440 字节灌不满 134217728 的帽、60s C22 到点走 timeout——正是审计 M2a"自证型"病的直接镜像。日志：`.scratch/wisp/probes/orch/20261120/191-cap-pin-mut128-run1.log`。
**还原**：`const sessionPendingCap = 64 << 20`，SHA256 改前 `C92EADFE512B0B1E469C60B239DA5C8FB504976B574C50B846FD85FA633049F6`＝改后同值 ⇒ 逐字节还原；`git status --porcelain -- internal/tools` 只余新钉文件。
**还原后两枚绿**（同口径 `-count=1 -v`）：`--- PASS: Test191SessionPendingCapNumberPinnedTo64MiB (0.00s)`、`--- PASS: Test191FloodRefusedPastLiteral68157440Bytes (0.30s)`，rc=0（日志 `191-cap-pin-teeth-run1.log`）。

## 3. 门禁四数与三把静态尺（最终态现跑）
- **`go test ./internal/tools -count=1 -v` rc=0**（`ok ... 35.448s`）：`=== RUN`=**603**，`--- PASS`=**446**（顶层）＋156（子测）＝**602**，`--- FAIL`=**0**，`--- SKIP`=**1**（`Test23LiveClipboardRoundTripThroughRealOS`，既有操作者点名门 `WISP_T23_LIVE_CLIPBOARD!=1`，非本腿引入）。日志：`191-cap-pin-tools-final-v.log`。
- **gofmt -l internal/tools = 空**；**gofumpt v0.12.0 -l internal/tools = 空**（CI 同款，`go install mvdan.cc/gofumpt@v0.12.0`，`v0.12.0 (go1.27.1)`）。⚠ 如实记一笔**顺手格式化**：起手尺点名两枚**既有**漂移文件（`git status` 无 M ⇒ 入库即有、非本腿引入）：`internal/tools/gate.go`（struct 字段对齐）、`internal/tools/ticket49_limits_test.go`（字符串续行 `+` 位置）——两处均为纯格式、零语义，且都在本腿写面 `internal/tools/**` 内，已 `gofmt -w` 补齐并在本笔 commit 具名；格式化后四数与首跑逐数相同（RUN=603/PASS=602/FAIL=0/SKIP=1）。
- **go vet ./internal/tools rc=0**。
- **`sh scripts/d22scan.sh` rc=0，0 findings**：`clean - no D22 ban violations`；末行逐字：`d22scan.sh: OK - control rc=0, scan rc=0`（日志 `191-cap-pin-d22scan.log`）。

## 4. 写面与硬约束对账
- 动了：`internal/tools/shell_session_cap_pin_191_test.go`（新，两枚钉）＋ `docs/evidence/s1/191-cap-pin-r1.md`（本件）＋ `internal/tools/gate.go`、`internal/tools/ticket49_limits_test.go`（仅 gofmt，见 §3）。
- `sessionPendingCap` **只钉不改**：变异只活在变异窗口内，交件树 = 原值 64MiB（SHA256 逐字节同）。
- 既有断言零放宽、零删除；`internal/panel/**`、`cmd/wisp/**`、`internal/ball/**` 零字节（`terminal_carrier_191.go:55` 只是**引用**该帽的注释，属他腿地界，本腿只读）。
- ⛔ push：本笔只 commit。
