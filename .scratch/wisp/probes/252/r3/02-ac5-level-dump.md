# 252-r3 AC#5：审计现场那一行不再被载具自己滤掉

点名现验：`cmd/wisp/ticket224_assembly_test.go:206-217 grantLinesOf`＝真。旧滤网针＝`GRANT`／`SESSION-MINT`／`REPLY `
（大小写敏感、逐字）；`tools: call …` 那行只带**小写** `grant_id=`、**不含任何一枚旧针** ⇒ 红名册天生看不见
`risk=`／`in_allowlist_scope=` 这一维（与 `224-c2` 的判一致）。

## 改前（旧滤网，pre-fix 树的 W1 红发现场）

- 发：`go test -count=1 -v -run '^TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking$' ./cmd/wisp/`（`62ca6c34`，短名 TMP；12:02:57→12:03:48；rc=1）
- `ticket224_assembly_test.go:324: audit is missing the GRANT-USE line naming row 1; got:` 之后的 dump 逐字**只有两行**：
  - `[audit] wisp run: SESSION-MINT id=sess_7a16170a273a3eef0dd0f70a417094cd (结束点＝本进程退出，A435 第 2 条)`
  - `[audit] session: GRANT-RECORD id=1 tool=fs.write pattern="C:\\tmp\\wisp252r2\\SHORT-~1\\…\\granted-in-this-session.txt" session=sess_7a16… expires_at=1794542608`
- 全日志 `grep -c "tools: call"` ＝ **0**；`grep -c in_allowlist_scope` ＝ **0**（`w1-before-short.log` 4,619 B）。

## 改后（交付载具，同形红复刻）

- 复刻法：在 pre-fix 树发上 `-overlay` 把 `cmd/wisp/ticket224_assembly_test.go` 换成交付态那一枚
  （overlay JSON＝`C:\tmp\wisp252r2\overlay-ac5.json`；⛔ 工作树零写、产码零动）。
- 发：`go test -count=1 -v -overlay='C:\tmp\wisp252r2\overlay-ac5.json' -run '^TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking$' ./cmd/wisp/`
  （12:05:35→12:07:06；rc=1；`--- FAIL: TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking (83.16s)`；`ac5-after-red.log` 5,415 B）
- 同一处 dump（现载具行号 `:333`）逐字新增两行，样行（另一枚 `corr=t224r2-hit` 同形）：
  `[audit] tools: call kind=approval-timeout task=t224r2-read corr=t224r2-ask tool=fs.write risk=L2 decision=timeout outcome=error rules_hit=[R1 R2] in_allowlist_scope=false grant_id=0 reason="R1: 工具声明为下界（L1）; R2: 目标路径在授权目录之外: C:\tmp\wisp252r2\SHORT-~1\TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking2885818523\002\never-granted.txt"`
- ⇒ 判据成立：等级维（`risk=L2`）、范围维（`in_allowlist_scope=false`）、命中规则（`R1 R2`）、原因句（`R2: 目标路径在授权目录之外`）都在红名册可读的 dump 里，且点名到短名目标路径本身。

## 边界（没做的）

- 只追加一枚针 `tools: call `（审计行自己的字面前缀，出处 `internal/tools/bridge.go` 的 `book()`），
  ⛔ 未放宽任何断言（只加 dump 行；diff＝+10/−1 全在注释与这一针）；⛔ 未动 `grantLinesIn`（AC#5 未点名）；⛔ 未覆盖/删除任何日志。

## 门禁（10-14 现跑，逐发时刻）

| 门 | 命令 | 读数 |
|---|---|---|
| build | `go build ./...` | 12:11:45→12:11:50，rc=0（`gate-build.log` 0 B） |
| gofumpt | `gofumpt -l cmd/wisp/ticket224_assembly_test.go` | 0 行（`gate-gofumpt.txt` 0 B） |
| d22scan | `sh scripts/d22scan.sh` | 12:12:16→12:12:42，rc=0；`d22scan: clean - no D22 ban violations`（`gate-d22scan.log` 23,758 B） |
| 长度门 | `sh scripts/check-path-length-budget.sh` | 12:12:50→12:12:52，**VERDICT GREEN**；`over-budget=59 covered by roster=59 not in roster=0`（`gate-length.log` 1,457 B） |

临时件（只建不删，全部在 `C:\tmp\wisp252r2\`）：`tree-before\`（62ca6c34 全树）· `SHORT-~1\`（短名临时根）·
`overlay-ac5.json` · `logs\` 13 枚＝`w1-before-short.log` 4,619 B／`w2-before-longtmp.log` 2,810 B／`w2-before-short.log` 2,844 B／
`w3-before-longtmp.log` 3,197 B／`w3-before-short.log` 3,457 B／`TestTicket224…-after-short.log` 3,336 B／
`TestComposedGate…-after-short.log` 3,231 B／`TestTicket101…-after-short.log` 3,626 B／`ac5-after-red.log` 5,415 B／
`gate-d22scan.log` 23,758 B／`gate-length.log` 1,457 B／`gate-build.log` 0 B／`gate-gofumpt.txt` 0 B。
