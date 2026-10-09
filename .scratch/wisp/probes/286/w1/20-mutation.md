# 286-w1 / 20 反形自证四件套（票 286 AC#4 ③）

目标：证明搬过的钉**仍然咬得住** —— 把产码 `callCorr`（`internal/agent/loop.go:603`）在盘上种成
"恒返回 taskID"那一形，新期望必须红。⛔ 不回退产码、种完即还原。

## 件一：种前 hash

```
cb7f9da7a6d277f341d3b40caee59eecec3bea4e4faee31777b1ad47c3763189 *internal/agent/loop.go
```
（`logs/mutant-pre-hash.md`；同一枚串在起手量里也出现过一次，两处一致。）

## 件二：`sed -n` 复量（种前 = 还原后，逐字同文）

`sed -n '600,612p' internal/agent/loop.go`：

```
// still goes through taskID - and appends the provider's call id, the id the
// wire, the history and the tool_call rows already key each call by. An
// id-less call falls back to its position inside the turn.
func callCorr(taskID, callID string, index int) string {
	if callID == "" {
		return fmt.Sprintf("%s#call-%d", taskID, index)
	}
	return taskID + "#" + callID
}

// executeCalls runs a turn's tool calls: risk policy first (pass-through until
// ticket 21), then D38d-capped concurrent execution with the C22 per-tool
// timeout, then D15(3) spill, then the results back into the history. Tool
```

种形（两枚 `return` 都换成 `return taskID`，即"恒返回 taskID"整形，`perl -pi -e` 两条替换）：

```
func callCorr(taskID, callID string, index int) string {
	if callID == "" {
		return taskID
	}
	return taskID
}
```

⛔ 唯一调用点 `loop.go:676 TaskID: taskID, CorrelationID: callCorr(taskID, p.call.ID, i),` 未动 ⇒
种的是**函数返回形**，不是把调用点改回 `taskID`（那才是回退产码）。
`grep -c 'fmt\.' internal/agent/loop.go` = 9 ⇒ 种形不会让 `fmt` 变成未用导入（不会假绿成编译失败）。

## 件三：红句逐字（`文件:行` 为新期望所在行）

命令（输出**先落文件** `logs/mutant-run.md` 再读，命令自身的 rc 单独取）：

```
GOFLAGS= go test ./internal/agent/ -count=1 -run TestGoldenSingleToolCall
```

```
rc=1
```

红句逐字：

```
--- FAIL: TestGoldenSingleToolCall (0.01s)
    loop_golden_test.go:81: call identity = {TaskID:fbb95865-367d-4b21-ae9c-b3e4c7e2a539 CorrelationID:fbb95865-367d-4b21-ae9c-b3e4c7e2a539 CallID:call_e1 Name:echo Args:{"text":"22 摄氏度，晴"} Timeout:2s}, want task fbb95865-367d-4b21-ae9c-b3e4c7e2a539 and correlation id fbb95865-367d-4b21-ae9c-b3e4c7e2a539#call_e1 (task id prefix + this call's own id, never the bare task id)
FAIL
FAIL	github.com/CarlosShao/wisp/internal/agent	0.040s
FAIL
```

★ 行号说明（⚠ 别按搬前的 `:69`/`:70` 找）：搬完之后 `if` 落在 **`:77`**，`t.Errorf` 落在 **`:81`**，
红句报的是 **`Errorf` 的行号 `loop_golden_test.go:81`** —— 与搬前红句报 `:70`（那正是当时的 `Errorf` 行）同形。
（`fbb95865-…` 三处 UUID（TaskID／CorrelationID／want）逐字同值，实测值为 `fbb95865-367d-4b21-ae9c-b3e4c7e2a539`，
逐字以 `logs/mutant-run.md` 为准；该 UUID 每次运行重新铸造，不承载判据。）

红之所以红（三钉逐条对上）：`HasPrefix(corr, taskID)` 仍真 ⇒ 第一钉不咬；
`TrimPrefix(corr, taskID+"#")` 得 `fbb95865-367d-4b21-ae9c-b3e4c7e2a539`（前缀不匹配时原样返回），不等于 `call_e1` ⇒ **第二钉咬**；
`corr == res.TaskID` 成立 ⇒ **第三钉咬**。第二行条件与 `TaskID` 等值那一半仍真，未一起搬。

## 件四：还原凭据（hash 等值 + porcelain 空 + 复绿）

种形由 `trap 'cp -p /tmp/loop286w1.go.bak internal/agent/loop.go' EXIT` 保证还原（中途炸也还原）。

```
pre  : cb7f9da7a6d277f341d3b40caee59eecec3bea4e4faee31777b1ad47c3763189 *internal/agent/loop.go
post : cb7f9da7a6d277f341d3b40caee59eecec3bea4e4faee31777b1ad47c3763189 *internal/agent/loop.go
HASH_EQUAL=yes
```

`git status --porcelain -- internal cmd` 原文：**空**（`PORCELAIN_LINES=0`）⇒ 突变体没有留在跟踪文件里。

还原后复跑同一条命令（`logs/restored-run.md`）：

```
ok  	github.com/CarlosShao/wisp/internal/agent	0.040s
RESTORED_TEST_RC=0
```

`sed -n '600,612p'` 复量与件二种前**逐字同文**（见上，同一次运行内两处输出一致）。

## ⛔ 我没有碰的东西（本步自证）

`internal/agent/loop.go:369`（`newTaskJournal(l.opt.Journal, taskID, taskID)`）、
`internal/agent/journal.go:59`／`:81`、同文件 `:341`（journal 侧反形，票 282 `AC#4`(a) 已裁"留"、台账 `A758`）
⇒ 全部一字节未动（终局尺＝`git diff --stat <锚>..HEAD -- internal/agent/testdata internal/agent/loop.go internal/agent/journal.go` 输出空，见 `30-gates.md`）。

★ 搬完后 `:341` 那一枚**仍然绿**（改后 FAIL=0 已含它）⇒ 没有触发派单里"若 `:341` 变红立刻停手回报"那一条：
本腿只动了请求侧，账本侧形状未被我顺带改掉。
