# 252-v2 AC#4 独立复读（三枚 witness：改前封存树红 / 交付封存树绿）

口径：改前＝`C:\tmp\wisp252v2\tree-before`（`2b1a3071^`＝`62ca6c34`）；改后＝`C:\tmp\wisp252v2\tree-delivery`（`d3cccfda`；零漂移核在 00 件）。两造同机、同 Git Bash、同 PATH（sherpa 前置）、同短名 TMP 根、`GOTMPDIR` unset。每发＝`go test -count=1 -v -run '^<name>$' ./cmd/wisp/`；⛔ 未跑整包。

| witness | 发向 | 起止（+08） | rc | 结果句（逐字） | 日志 |
|---|---|---|---|---|---|
| W1 `TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking` | tree-before | 12:26:29→12:27:29 | 1 | `--- FAIL: … (43.48s)` | 4,615B |
| W1 | tree-delivery | 12:46:05→12:46:19 | 0 | `--- PASS: … (5.76s)` | 3,331B |
| W2 `TestComposedGateBlocksAWriteForTwoSeconds` | tree-before | 12:28:34→12:33:44 | 1 | `--- FAIL: … (303.07s)` | 2,844B |
| W2 | tree-delivery | 12:40:11→12:40:23 | 0 | `--- PASS: … (5.15s)` | 3,219B |
| W3 `TestTicket101SessionGrantDoesNotCrossRestart` | tree-before | 12:36:00→12:36:14 | 1 | `--- FAIL: … (6.55s)` | 3,448B |
| W3 | tree-delivery | 12:40:37→12:42:07 | 0 | `--- PASS: … (83.11s)` | 3,615B |

红句逐字（键行；全文在 `logs/10-ac4-verbatim.txt`）：
- W1：`ticket224_assembly_test.go:269: the control call errored: the L1 window running out means execute`；`:311 a live session grant left 1 cards on screen for the covered call: 「命中授权就不再弹卡」 is not happening in the assembled run (… pattern="C:\\tmp\\wisp252v2\\SHORT-~1\\…granted-in-this-session.txt", row=1)`；`:316 covered call booked decision="timeout", want "allow_session_grant"`；`:327 …The system cannot find the file specified.`
- W2：`run_test.go:378: an unvetoed L1 window means EXECUTE, got: 审批超时（300 秒未确认），C18 一律判拒绝，已自动拒绝`
- W3：`run_mode101_test.go:506: control write should NOT error under auto_approve (审批超时（1 秒未确认），C18 一律判拒绝，已自动拒绝); the control half is broken, so the refusal above proves nothing`

与 r3 读数同形（r3→v2）：43.52→43.48／303.23→303.07／6.20→6.55；5.55→5.76／5.78→5.15／83.47→83.11。W3 的 83s＝其载具两枚 40s C18 窗口（rc=0 非挂死）。短名形状现证：各发日志 `SHORT-~1` 命中 4~9 次。

## 301 秒归账（票面原句逐字复核）

- 票尾 bullet 在（票面 `:119` 行内）逐字（含排版符）：`★**301 秒的账从"C18 超时既有一笔"改归本票（票 252）**，⛔ 别留两笔：它等的不是 C18 本身，而是被 R2 判成 L2 后无人应答的那张卡（本机 303.23s／CI 301.70s·301.12s 同形）。`
- 01 件在（`probes/252/r3/01-ac4-witness-readings.md:22`）逐字（含排版符）：`**账目归口（票面要求的原句，逐字落在这里）**：这 301 秒的账（本机 303.23s ＝ C18 default 300s ＋ 开机；CI 两发 301.70s／301.12s 同形）**从"C18 超时既有一笔"改归本票（票 252）**，⛔ 别留两笔——它等的不是 C18 本身，而是被 R2 判成 L2 后无人应答的那张卡。`
- 因果本腿自核＝**成立**：①超时句只出自 `internal/agent/approval/queue.go:479`（`审批超时（%d 秒未确认），C18 一律判拒绝，已自动拒绝`）＝审批队列里无人应答的卡到点自动拒；`internal/tools/bridge.go:351` 路由注释逐字 `L0 pass, L1 window, L2 approval` ⇒ W2 红句证明该发走进了队列（＝非 L1）。②"进队"的因＝短/长拼法分岔使 `InAllowlist=false` ⇒ R2 ⇒ L2（`internal/risk/rules_gateway.go:49`；修法提交 `2b1a3071` 自陈 `became R2 and paid for an L2 card`；同形直接读数＝本腿 AC#5 overlay 红 dump 逐字 `risk=L2 … in_allowlist_scope=false … rules_hit=[R1 R2] … kind=approval-timeout`）。③封存交付树下 W2 5.15s 绿（无卡、无超时），两树唯一变量＝折叠修复 ⇒ 300 秒的账归本票。旁证：W3 同机制配 `confirm_timeout_sec=1` 的 1 秒超时句 ⇒ 等的是卡，不是 300 这个常数。
