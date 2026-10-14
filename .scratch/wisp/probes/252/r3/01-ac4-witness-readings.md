# 252-r3 AC#4：三枚 witness 同机对照读数（改前必红 / 改后必绿）

口径：改前＝仓外树 `C:\tmp\wisp252r2\tree-before`（`62ca6c34`＝`2b1a3071^`）；改后＝工作树交付态（载具 commit `d3cccfda`）。
两造同机、同壳（Git Bash）、同 PATH（`third_party/sherpa-onnx` 前置）、同 TMP 短名根（见 `00-anchor-and-method.md`）。
每发＝`go test -count=1 -v -run '^<name>$' ./cmd/wisp/`；⛔ 未跑整包。日志（原始件）在 `C:\tmp\wisp252r2\logs\`，字节数逐枚点在各节。

## W1 TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking

- **改前红**（12:02:57 → 12:03:48；rc=1；`--- FAIL: TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking (43.52s)`；`w1-before-short.log` 4,619 B）逐字：
  - `ticket224_assembly_test.go:269: the control call errored: the L1 window running out means execute`
  - `ticket224_assembly_test.go:293: the granted call errored: true`
  - `ticket224_assembly_test.go:311: a live session grant left 1 cards on screen for the covered call: 「命中授权就不再弹卡」 is not happening in the assembled run (ledger id=sess_7a16…, pattern="C:\\tmp\\wisp252r2\\SHORT-~1\\…\\granted-in-this-session.txt", row=1)`
  - `:316 covered call booked decision="timeout", want "allow_session_grant"`；`:319 grant_id=<nil>`；`:327 the covered call should have written …: The system cannot find the file specified.`
- **改后绿**（12:07:43 → 12:07:59；rc=0；`--- PASS: … (5.55s)`；`TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking-after-short.log` 3,336 B）
- 归因佐证（AC#5 同形红发里逐字）：`reason="… ; R2: 目标路径在授权目录之外: C:\tmp\wisp252r2\SHORT-~1\…never-granted.txt"`

## W2 TestComposedGateBlocksAWriteForTwoSeconds（`run_test.go:378`）

- **改前红**（12:02:52 → 12:08:02；rc=1；`--- FAIL: TestComposedGateBlocksAWriteForTwoSeconds (303.23s)`；`w2-before-short.log` 2,844 B）逐字：
  `run_test.go:378: an unvetoed L1 window means EXECUTE, got: 审批超时（300 秒未确认），C18 一律判拒绝，已自动拒绝`
- **改后绿**（12:07:59 → 12:08:12；rc=0；`--- PASS: … (5.78s)`；`TestComposedGateBlocksAWriteForTwoSeconds-after-short.log` 3,231 B）
- **账目归口（票面要求的原句，逐字落在这里）**：这 301 秒的账（本机 303.23s ＝ C18 default 300s ＋ 开机；CI 两发 301.70s／301.12s 同形）**从"C18 超时既有一笔"改归本票（票 252）**，⛔ 别留两笔——它等的不是 C18 本身，而是被 R2 判成 L2 后无人应答的那张卡。

## W3 TestTicket101SessionGrantDoesNotCrossRestart（`run_mode101_test.go:506`）

- **改前红**（12:02:26 → 12:02:39；rc=1；`--- FAIL: TestTicket101SessionGrantDoesNotCrossRestart (6.20s)`；`w3-before-short.log` 3,457 B）逐字：
  `run_mode101_test.go:506: control write should NOT error under auto_approve (审批超时（1 秒未确认），C18 一律判拒绝，已自动拒绝); the control half is broken, so the refusal above proves nothing`
- **改后绿**（12:08:12 → 12:09:42；rc=0；`--- PASS: … (83.47s)`；`TestTicket101SessionGrantDoesNotCrossRestart-after-short.log` 3,626 B）。
  83.47s 是该载具自己的两枚 40s C18 窗口（267-r2 抬过的 `confirm_timeout_sec = 40`；日志内 boot1 12:08:22 → boot2 12:09:02），rc=0 且 `PASS`，不是挂死。

## 形状对照（同机、同树、只换 TMP 拼法）

| 树 | TMP | 结果 |
|---|---|---|
| pre-fix `62ca6c34` | 长名（GOTMPDIR 透传形） | W2 PASS 5.80s ／ W3 PASS 5.47s（`w2/w3-before-longtmp.log`） |
| pre-fix `62ca6c34` | 短名 `SHORT-~1` | **W1 FAIL 43.52s ／ W2 FAIL 303.23s ／ W3 FAIL 6.20s** |
| 交付态（`d3cccfda`） | 短名 `SHORT-~1` | **W1 PASS 5.55s ／ W2 PASS 5.78s ／ W3 PASS 83.47s** |

⇒ 三枚红都只在 runner 的短名拼法下出现（本机长名天生绿，与 `224-c2` census 记的本机 3.17/3.20/3.23s 同向）；短名下 pre-fix 三枚全红、修后三枚全绿。

## 没做（具名）

- ⛔ 未跑整包 `./cmd/wisp/`；⛔ 未 push；⛔ 产码语义零动（本轮唯一 commit 只动 `cmd/wisp/ticket224_assembly_test.go` 的 dump 滤网）。
- 改前红的树是 `62ca6c34` **全树**（含当时的载具与常数）；唯一把"现载具塞进旧码"的发＝AC#5 的 `-overlay` 发，具名在 `02-ac5-level-dump.md`。
- W1 只有短名红发、没有长名对照发（W2/W3 有）：按 census 先例本机长名天生绿，未补跑；点名候选留给异议者。
