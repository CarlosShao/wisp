# 191-audit-r1 — 票 191 三波的对抗验收（裁决者 ≠ 实现者，2026-11-20 23:0x）

**总判：3 成立 / 1 夸大 / 0 不成立。放宽断言 0 条。** 变异序列（1/2/5 三组）`rc=1 · rc=1 · rc=1` 全红。

## ① "nil 回调逐字节一致" —— **成立**（口径另记）
- **overlay 对照**（`-overlay` 把 `6512d5a73~1` 的 `shell_session.go` 换进去，仓内文件零改动）：两次 `AUD191PROBE text=` **同一串** ⇒ 增量在 nil 回调下**零字节差**。
- 变异 M1c（sink 无条件安装并影响 Result）⇒ `Test191NilCallbackByteIdentical` FAIL，逐字 `nil "ok: echo hello 2>&1\r\n\n退出码：0…"` vs `cb "ok: echo hello 2>&1\n…"`。
- ⚠ **夸大的一处（属编排者派单口径错）**：我派单要求"少一次 flush ⇒ 定向用例应红（逐字）"。**实际不红**——分片器把 marker 行拆成三段跨读，丢一次 flush 只让它**晚一轮**并入（连跑 10 次全绿）；只有"把首轮字节直接丢弃"才红，而红形是 `panic: test timed out`（等不到边界），不是逐字断言。⇒ 该条要求**用例没兑现**，牙只兑现了"nil 路径也调 sink"那一侧。此口径错误记在本件，不掩饰。
- 口径挑明：`internal/tools/bridge.go:732-736` 的 onUpdate 闭包**恒非 nil**（生产唯一调用路径 `:738`）⇒ "nil 回调"只在直接调用/测试面可达（191 证据 `191-stream-tools-r1.md:55` 已自认）。

## ② "上限真拦得住" —— **成立**（工具侧数值未钉死）
- 面板 8192：M2c（8192→16384）⇒ `terminal_carrier_191_test.go:83` FAIL 逐字 `the terminal bounds moved: window=16384 input=4096, want 8192/4096`。
- 工具 64MiB：**M2a（64→128MiB）rc=0 全绿**（`go test ./internal/tools` ok 30.611s）⇒ **上限数值本身无用例钉死**（测试期望由同一常量派生，自证型）。M2b（删 chunk 侧帽守卫）⇒ `:206` FAIL `streamed chunks must stay inside the buffer cap: delivered=135266304 cap=134217728`；M2c（删搜索循环拒绝块）⇒ `:203` FAIL（60s 到点，帽路径消失）。
- 还原：两文件 SHA256 回原值，`git diff` 空。

## ③ "只读不执行" —— **成立**（派单 grep 口径拆开）
- 派单原文 `git grep -n "exec.Command|CommandContext" -- internal/panel` 是**基本正则**，`|` 为字面竖线 ⇒ 命中 0 条，**无信息量**（口径错，记此）。
- E-regex 数：ALL 13→19（+6 全在 `terminal_carrier_191_test.go` 的注释/测试名字面），**生产面 3→3→HEAD 3**。
- **AST 独立普查**（不靠 grep、不靠仓内尺）：`internal/panel` 生产面 `exec.*` 调用点 **2**（`git_readonly_exec.go`、`tree_reader_190.go`，均 `CommandContext`），`os/exec` import 2；**改前改后同为 2 ⇒ 面板侧零新增执行点**。全仓生产面 `exec.*` 38 处，191 三波未增删其一。
- 载体文件 import 名册 = `{sync}`（`:3` 唯一 import，逐字 + AST 双证）；`TerminalFeed` 方法集 = `{Input, OutputChunk, Section}`（`go doc` 恰好三枚）；正控在跑（植入 `exec.Command` 到副本必须被点名）。

## ④ 账本收口与界面三态
- `composer_test.go` 两笔：**删除/改写断言逐条列出后判"收口/收紧"，放宽 0 条**（见审计正文第 4 节）。
- 界面三态：缺席显「终端未连接」＋空态「尚无输入/尚无输出」两态**不共用一张脸**；变异"缺席态改显示空态文案"⇒ rc=1 红。

## ⑤ 零 emoji / 零 dist 字节 / 未碰他人件：三笔提交无 `frontend/dist`；`frontend/fixtures/composer-states.html` 未入任一提交。

## 遗留（已派小活补齐）
工具侧 `sessionPendingCap` 的**数值**未被钉死（M2a 绿）⇒ 追加一枚数值钉 `sessionPendingCap == 64<<20` 并证明"改大即红"。
