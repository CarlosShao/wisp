# 309-v1 · `AC#2` 枚数钉是不是真钉（`internal/tools/toolrequest_nail_309_test.go`）

判据（票面逐字）＝"本格只回答并落地一件**最小事**：**将来任一 caller 手搭 `ToolRequest` 递进空 `TaskID`，要让它在盘上可核地响**"。
`AC#0` 裁语把落点收窄成形 B 枚数钉：`go/parser` 读生产源码，断言 `agent.ToolRequest{` 在 `internal/agent` 之外的生产构造点枚数＝**0**，且**必须自带分母**。

## 1. 分母是真算出来的，⛔"扫了 0 个文件然后报绿"
- 本腿自己跑（件＝`logs/pristine-309.log`／`logs/nail-0-baseline.log`），钉**打印**出来的那行逐字＝
  `AC#2 枚数钉：internal/agent 之外 agent.ToolRequest 生产构造点=0（==0），分母=245 枚非测试源`
- 独立第三方尺对账（⛔ 用钉自己的数自证）：
  `find cmd internal -name '*.go' -not -name '*_test.go'` ＝ **269**；其中 `internal/agent` 下＝**24**；`269 − 24` ＝ **245** ✅
  ⇒ 分母＝"两把尺互相咬住"，⛔ 假清洁。
- 内置正控（`:119-130`，内存里 parse 一枚 `agent.ToolRequest{TaskID: ""}` 断言计数==1，否则 `t.Fatalf "the ruler is blind"`）在场且随包跑绿。
  ⚠ **它自己的红⛔ 可能**（那要把腿写的测试件改掉，⛔ 我的写面）。本腿改用**扫描面**把"这把尺⛔ 瞎"证成可测：见 §3 三形 Rig——同一枚 matcher 对字面量计 1、对另外两形计 0，`==1` 那枚正控因此是**有区分度的 1**，⛔ 恒真。

## 2. 腿的 `nail-rig.log` 声明能不能复现（加一枚构造点⇒钉红），且 Rig 加在**钉真扫的文件**里
- 腿那一发的件＝`.scratch/wisp/probes/309/r1/logs/nail-rig.sh`：它新建 `internal/tools/zz309rig_scratch.go`（非 `_test.go`、⛔ 在 `internal/agent` 下）⇒ 确实落在钉的射程里；读数 `BEFORE=0/分母245` ⇒ `RIG=1/分母246`／`RIG_rc=1`。
- 本腿⛔ 建别人的瞬态文件，改打**自己那枚已扫的文件** `internal/tools/bridge.go`（Rig＝`logs/nail-rig.txt`，读数＝`logs/nail-1-literal.log`）：
  追加逐字 `return agent.ToolRequest{TaskID: "", Name: "fs.read"}`，LANDS 自证＝`tail -5` 打出那三行；钉红逐字＝
  `toolrequest_nail_309_test.go:108: 枚数钉响：internal/agent 之外发现了 1 枚 agent.ToolRequest{ 生产构造点 (分母=245 枚解析源，roots=["..\\..\\cmd" "..\\..\\internal"])`／`rc=1`
  ⇒ 复现成立，且**分母没动（245）而 constructions 从 0→1**＝钉响的键值确实是"枚数"，⛔ 只是"文件数变了"。
- 还原：`RESTORE cmp: byte-identical OK` ＋ `sha256sum cb86eaa7…`（与 pristine 同一枚），收尾 `git status --porcelain -- internal/tools/bridge.go`＝**0 行**。

## 3. 覆盖攻击：这枚钉只认"限定名复合字面量"那一种拼法（三形实测，⛔ 形容词）
桥注释具名了三类未来候选者（自带 id 的面板/球宿主、重试包装器、跨轮次带 id 的定时任务）。本腿把三种形状各打一次到**被扫文件**里：

| 形 | 打进去的代码（逐字形状） | 钉的反应 | 读数 |
|---|---|---|---|
| **1 字面量** | `return agent.ToolRequest{TaskID: "", Name:"fs.read"}` | **红**（constructions=1） | `logs/nail-1-literal.log` rc=1 |
| **2 声明＋逐字段赋值** | `var rig agent.ToolRequest` ＋ `rig.TaskID = ""` | **绿**（constructions=0）＝瞎 | `logs/nail-2-decl.log` rc=0 |
| **3 别名后字面量** | `type T = agent.ToolRequest` ＋ `return T{TaskID: "", …}` | **绿**（constructions=0）＝瞎 | `logs/nail-3-alias.log` rc=0 |

成因逐字在 matcher 里（`:33-51`：只认 `*ast.CompositeLit` 且 `cl.Type` 是 `*ast.SelectorExpr` 且 `pkg.Name=="agent" && sel.Sel.Name=="ToolRequest"`）⇒ 形 2 根本没有复合字面量、形 3 的 `cl.Type` 是 `*ast.Ident`。
另两枚⛔ 打也能定的射程洞（尺＝读数，⛔ 猜测）：
- **`internal/agent` 整包被跳过**：`strings.Contains(rel, "internal/agent/")` ⇒ 今天 `internal/agent` 下 24 枚非测试源**一枚⛔ 解析**。重试包装器最自然的落点就在 `internal/agent` 里 ⇒ 钉⛔ 响。（而且 `loop.go:675` 那枚真构造点是**裸名** `ToolRequest{`，同包⛔ 需要限定名，matcher 本来就⛔ 认它。）
- **射程只有 `cmd`＋`internal`**：同 module 的 `tools/` 下还有 **13** 枚非测试 `.go`（`git ls-files` 尺；今天它们⛔ import `internal/agent`＝`grep -rn "wisp/internal/agent" tools/` 命中 0），但一枚未来放那边的构造点⛔ 钉的射程内。

## 4. `AC#2` 自己的文字要不要得比钉交的更多？——要，故⛔ 判"成立"
本格逐字写的是"**任一** caller 手搭 `ToolRequest` 递进空 `TaskID`…要让它在盘上可核地响"。钉兑现的是"在 `cmd|internal` 射程内、以 `agent.ToolRequest{` 字面量拼法、写在 `internal/agent` 之外"的那一**子集**会响（§3 三形里两形实测不响）。
另：`AC#0` 裁语给的正是这一形（枚数钉、matcher 名字都写死了 `agent.ToolRequest{`），所以腿⛔ 走歪；洞是**裁语与格文之间的缝**，⛔ 是实现者擅自缩水。
钉头注释逐字"要变成'有人新写构造点就会红'"，与注释里那句 "Test309ToolRequestProductionConstructionNail pins `agent.ToolRequest{` construction points outside internal/agent to 0" —— 前一句**overclaim**（见 `60-overclaim.md` 第 2 条），后一句带 backtick 限定拼法，⛔ 越界。

## 5. 残余登记（具名，交编排者）
`Q-309-v1-nail-misses-decl-and-alias-shapes`（§3 形 2/3 实测绿）、`Q-309-v1-nail-excludes-whole-agent-package`（24 枚源整包跳）、`Q-309-v1-nail-roots-miss-tools-dir`（13 枚源在射程外）。三枚都是**假阴性**（该响⛔ 响），⛔ 一枚是假阳性。
