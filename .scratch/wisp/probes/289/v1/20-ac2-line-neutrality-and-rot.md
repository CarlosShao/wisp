# 票 289 / 腿 289-v1 — AC#2 判语：两把行中性尺自跑＋票面锚腐烂有没有发生

**判语：成立。**（两把尺都是本腿自己现跑；并且**本票确实没有推动任何一枚票面行号锚**——腐烂检查现读通过，但现读同时暴露出**两处先于本票就已过期**的票面锚，具名上报、不算本票的账。）

## 1. 尺 A — 每枚件 `+N == -N`

命令逐字：`git show --numstat --format="" 73c30dfe`
原样输出：

```
4	4	cmd/wisp/subagent_carrier_197_test.go
2	2	internal/panel/subagent_blocked_220_test.go
3	3	internal/panel/subagent_roster_197.go
5	5	internal/panel/subagent_roster_197_test.go
```

`git show --stat --format="%H %ad" 73c30dfe` 的合计行逐字：` 4 files changed, 14 insertions(+), 14 deletions(-)`（作者时间戳 `Fri Oct 9 14:26:18 2026 +0800`）。
⇒ 四枚逐枚相等（4/4、2/2、3/3、5/5），⛔ 无一枚净增删。**成立。**

## 2. 尺 B — diff 里非 `+++`/`---` 的行**每一行都是纯注释行**

命令逐字（先落盘再量，`rc(dump)=0`）：

```
git show -U0 --format="" 73c30dfe > .scratch/wisp/probes/289/v1/raw-v1-diff-u0.md
total=$(grep -E '^[+-]' .scratch/wisp/probes/289/v1/raw-v1-diff-u0.md | grep -vE '^(\+\+\+|---)' | wc -l)
cmt=$(grep -E '^[+-][[:space:]]*//' .scratch/wisp/probes/289/v1/raw-v1-diff-u0.md | wc -l)
```

读数原样：`total=28 cmt=28 noncomment=0` ⇒ **非注释行＝0**。

## 3. 为什么要行中性（本票的牙）＋第三把尺：语义面为零

`internal/panel/subagent_roster_197.go` 是**产码件**，票 220／票 197 的现量表**按行号**引它 ⇒ 净增删行会让那些写在票面里的锚**静默腐烂**（同形先例＝票 279 `AC#2` 实测"插 +3 行 ⇒ `config_receipt_255_test.go:179` 那把按 `lines[n-1]` 直取的尺真红"）。
本腿补一把更硬的尺（不只"行数相等"，而是"把注释行整枚剥掉后两态逐字节相同"）：

命令逐字（`fd692f7e`＝r1 改前的 HEAD，`73c30dfe`＝改后）：

```
for p in <这 4 枚件>; do
  x=$(git show "fd692f7e:$p" | grep -vE '^[[:space:]]*//' | md5sum)
  y=$(git show "73c30dfe:$p" | grep -vE '^[[:space:]]*//' | md5sum)
done
```

原样输出（四枚 `same=YES`）：

```
internal/panel/subagent_roster_197.go          pre=57b80471868328ff4835dcde24a7c4b1 post=57b80471868328ff4835dcde24a7c4b1 same=YES
internal/panel/subagent_roster_197_test.go     pre=c62a823e19061e7cdcf2f66a38076c34 post=c62a823e19061e7cdcf2f66a38076c34 same=YES
internal/panel/subagent_blocked_220_test.go    pre=38eed8e1e54a9a71edb8052cc9435106 post=38eed8e1e54a9a71edb8052cc9435106 same=YES
cmd/wisp/subagent_carrier_197_test.go          pre=9369cb5f166f20debad5bb033eac3889 post=9369cb5f166f20debad5bb033eac3889 same=YES
```

⚠ 这把尺有一个已知盲点：`grep -vE '^[[:space:]]*//'` 只剥**整行注释**，如果改动把 `//go:` 之类的编译指令塞进/拆出注释，两态仍可能"看起来相同"。所以补测：那 28 行里**指令形注释＝0**（`grep -cE '^[+-][[:space:]]*//[[:space:]]*(go:|line |export|extern|go:build|\+build)' raw-v1-diff-u0.md` → `0`，`rc(directive-grep)=1` 即"零命中"）。
⇒ **语义面＝零**，这一条同时是 AC#3 里"注释-only 能不能改变测试行为"的回答（见 `30-ac3-test-comparison.md` §3）。

## 4. 腐烂检查：票 220／票 197 按行号引的锚，今天还是不是那一行

尺一（本票有没有推动锚）——同一行号在两态逐字对照：

```
for n in 57 127 190 205 210 236; do
  a=$(git show fd692f7e:internal/panel/subagent_roster_197.go | sed -n "${n}p")
  b=$(git show 73c30dfe:internal/panel/subagent_roster_197.go | sed -n "${n}p")
  [ "$a" = "$b" ] && echo "line $n SAME" || echo "line $n DIFF"
done
```

原样输出：`line 57 DIFF`（＝被改写的那句注释本身，`pre=[	// loop dispatches with CorrelationID == TaskID (internal/agent/loop.go:647),]` → `post=[	// internal/tools/subagent_197.go:262), never its correlation id: since]`）；
`line 127 SAME`／`line 190 SAME`／`line 205 SAME`／`line 210 SAME`／`line 236 SAME`。
另核第 4 枚件里被注释引用的夹具行：`sed -n '173p' cmd/wisp/subagent_carrier_197_test.go` 现读
`			TaskID: rootID, CorrelationID: rootID, CallID: fmt.Sprintf("spawn-197-%02d", i),` ⇒ **仍在 `:173`**。
⇒ **本票没有让任何一枚票面锚腐烂。成立。**

尺二（顺手现读这些锚今天对不对，⚠ 结论：**两处已过期，但都不是本票造成的**）：

| 票面锚 | 票面原文（逐字摘要） | 今天盘上现读 | 判 |
|---|---|---|---|
| 票 220 `:11`＝`subagent_roster_197.go:190-210` | "先 `waiting[card.CorrelationID] = true`（键＝correlation id），再 `BlockedOnApproval: row.TaskID != "" && waiting[row.TaskID]`（查＝task id）"，尺＝`sed -n '185,215p'` | `:190-198` 是 `taskRosterSectionFrom` 的注释＋签名（`:191` 逐字 `// taskRosterSectionFrom is the whole render with one extra source: l1Windows, the`），建表在 `:212-221`、查询在 `:236` | ⚠ **锚已漂**——漂因＝**票 220 自己的落地**（加了 `l1Windows` 形参与两半建表），不是票 289；`line 190/205/210/236 SAME` 已证本票零推动 |
| 票 197 `:66`＝`subagent_197.go:99-104`"那四枚在跑的态"（⚠ 这枚锚指 `internal/tools/subagent_197.go`，不是 panel 那枚同名件；派单把它算进 `subagent_roster_197.go` 的族里，本腿按原文落点核） | 四枚态＝`Thinking`／`Settling`／`Muted`／`Error` | 现读 `:98` 逐字 `//	Thinking  the child loop is running`、`:99` `Settling`、`:100-101` `Muted`、`:102` `Error` ⇒ 四枚今天住在 **`:98-102`**，`99-104` 少了 `Thinking`、多含 `:104` 那句"Reconciling…" | ⚠ **锚差一枚**；`git diff --name-only fd692f7e..HEAD -- internal/tools/subagent_197.go` 输出 0 行 ⇒ 本票没碰这枚件 |
| 票 197 `:66`／票 220 `:16`＝`subagent_roster_197.go:127`"blockedOnApproval 是一枚布尔" | — | `:127` 现读逐字 `	// BlockedOnApproval reports that a card for this very task is on the queue`，字段声明在 `:132` `	BlockedOnApproval bool \`json:"blockedOnApproval"\`` | ✓ 锚落在该字段的注释块首行，语义没错位；行中性尺已证 `:127` 两态 SAME |

⇒ AC#2 本格**成立**；两处票面锚过期**具名上报给编排者**（属票 220／票 197 的票面文字缺陷，⛔ 本腿不改票面一字，也不把它们算进票 289 的账）。
