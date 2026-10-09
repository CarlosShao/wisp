# 票 292 腿 292-r1 AC#2 — 改前/改后逐字 + 两把尺

改动范围：`cmd/wisp/subagent_carrier_197_test.go` 一枚件，`:598-600` 三行，全是注释行。

## 1. 改前（HEAD `a3782ff` 逐字）

```
:598 //  1. blockedOnApproval is never asserted on a production packet. The join needs a card
:599 //     whose correlation id equals a CHILD task id, and nothing in this tree can make
:600 //     that happen from the outside: the agent loop never sends tool_choice, so mockllm
:601 //     never answers with a tool call, so a spawned child never reaches the gate. The
:602 //     field is therefore asserted at the pump against a real ApprovalCardView
:603 //     (internal/panel's TestTheRosterReaderPutsSubagentsOnTheWire), and on this path it
:604 //     is measured only as false.
```

## 2. 改后（盘上逐字，本腿交付）

```
:598 //  1. blockedOnApproval is never asserted on a production packet. The join has two
:599 //     producers (subagent_roster_197.go:214, :219-220): a card's corr, never a task id
:600 //     for the loop (loop.go:603), and L1 windows no run feeds; no tool_choice is sent, so mockllm
```

（`:601-604` 一字未动，见本件 §4。）

## 3. 三节是否都说了实话（对票面 AC#2 的三条要求）

| 要求 | 新句里的落点 | 现场尺（写进注释的每个名与行号都现跑过） |
|---|---|---|
| join 有两枚来源 | `The join has two producers (subagent_roster_197.go:214, :219-220)` | `sed -n '214p' internal/panel/subagent_roster_197.go` = `indexWaitingKey(waiting, card.CorrelationID)`；`:219` = `indexWaitingKey(waiting, w.CorrelationID)`；`:220` = `indexWaitingKey(waiting, w.TaskID)` |
| 卡片那一半的形状今天不可能（不再写「只是从外面做不到」） | `a card's corr, never a task id for the loop (loop.go:603)` | `sed -n '603p' internal/agent/loop.go` = `func callCorr(taskID, callID string, index int) string {`，两支（`:605`/`:607`）都追加 `#` 段 ⇒ 恒不等于 taskID。措辞用「never」而非「hard to reach from outside」 |
| 生产侧 L1 那一半还没有人供数据（⛔ 不许写成「已经能点亮」） | `and L1 windows no run feeds` | `grep -c L1Windows cmd/wisp/run.go` = 0；装配点 `cmd/wisp/run.go:699` 的 `PumpSources` 无 `L1Windows` 字段（读数在 00 件 §3 问 c） |

## 4. 两把尺（并排）

- 尺 1 `git diff --numstat -- cmd/wisp/subagent_carrier_197_test.go` = `3	3	cmd/wisp/subagent_carrier_197_test.go` ⇒ `+N == -N` 成立。
- 尺 2 `git diff -U0` 非 `+++`/`---` 行（全量存 `logs/ac2-changed-lines.txt`）：`total=6 comment=6 noncomment=0` ⇒ 纯注释行枚数＝全量枚数，非注释行 0 枚。
- 行宽：`:598=83 / :599=87 / :600=98`。该文件 HEAD 版本本身最长行 = 110 列、≥96 列的行本就有 13 枚，故 98 列未越出本文件既有风格；gofmt/gofumpt 不复排注释，实测 `gofmt -l cmd/wisp/subagent_carrier_197_test.go` 空（rc=0）。
- `gofmt -l cmd/wisp` 整包在 HEAD 上本来非空（3 枚既有件），见 30 件。

## 5. 两处必须向编排者交代的取舍（不是我擅扩范围，是被射程逼的）

1. **那三行不是一个完整句子**：票面「现量」只列 `:598-600`，但这句话的**语法尾巴在 `:601`**（`never answers with a tool call, so a spawned child never reaches the gate.`），`-602` 才是 `field is therefore asserted…`。
   ⇒ `:600` 的结尾**必须是一个能当「never answers」主语的名词**（原文是 `mockllm`）。我在中途一度把 `:600` 写成以 `so the field` 收尾，读出来变成 "so the field never answers with a tool call" —— 那是病句，已改回。
   ⇒ 直接后果：**`no tool_choice is sent, so mockllm` 这一节我必须留在 `:600`**，不能删。本腿起手时按派单给的三节清单，以为可以把这条因果链整个丢掉（派单 §AC#2 只要求那三节，没提 tool_choice）；实际做不到，理由具名报回：删了它 `:601` 就没有主语，而 `:601` 属「本件其它注释」、⛔ 不许动。
   保留它同时也是实话：AC#1 问 ⓓ 裁过 `grep -rn ToolChoice internal/agent/` 全量 0 命中，`tools/mockllm/chat.go:172` 把整个 tool-call 决策挂在 `req.ToolChoice != nil` 上 ⇒ 两环今天都为真。
2. **行宽**：三行装下三节 + 四个引用，`+N==-N` 又锁死了行数，`:600` 只能到 98 列（该文件最长 110 列）。若编排者更看重「注释一律 ≤88 列」，正确处置是**另立一票允许把这段重排成 4 行**（`+N==-N` 仍成立、但会改到 `:601`，本票明令禁止），本腿不自作主张。

## 6. 没有碰的东西（自证）

- `:311-313` 那句 `t.Error` 原文仍在原位、逐字未改（`sed -n '305,320p'` 读数在 00 件之后的现场核）。
- 本件其余注释、任何 `func`/`if`/字段：零改动（尺 2 的 `noncomment=0`）。
- `PumpSources.L1Windows` 的生产供数据：**没接**（票 220 落地腿的活）。
- `git diff --numstat` 全仓（本腿名下）只有这一枚 `.go` 件，见 commit 的 `--name-status`。
