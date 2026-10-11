# 309-r2 · 20 三形各一发"加一枚即红"（job 3）

脚本＝`logs/nail-rigs.txt`（工作树面，只打 `internal/tools/bridge.go`，每发当场还原）与
`logs/nail-rigs-tree.txt`（两枚仓外导出树）。读数＝`logs/nail-rig-*.log`，汇总＝`logs/nail-rigs.stdout`。
⛔ 一次跑完、⛔ 留突变：每发 `RESTORE cmp: byte-identical OK` ＋ `sha256sum` 与 pristine 同一枚
（`069cb275dc03b8051c46e4c11baad9c4037d62561371defb99aa78a1599e9715`）。

| 形 | 打进去的逐字形状 | 旧钉（309-v1 那发） | 新钉（本腿） | 读数 |
|---|---|---|---|---|
| RIG-0 | ⛔ 打（基线） | 绿 分母 245 | **绿**，构造点 1／分母 293 | `nail-rig-0-baseline.log` |
| RIG-1 限定名字面量 | `return agent.ToolRequest{TaskID: "", Name: "fs.read"}` | 红（v1 `nail-1-literal.log` rc=1） | **红 rc=1**，构造点 2，名册 `internal/agent/loop.go:675(same-package-bare-literal) ｜ internal/tools/bridge.go:1346(qualified-literal)`，分母⛔ 动（293） | `nail-rig-1-literal.log` |
| RIG-2 声明＋逐字段赋值 | `var rig309r2Decl agent.ToolRequest` ＋ `rig309r2Decl.TaskID = ""` | **绿 rc=0＝瞎**（v1 `nail-2-decl.log`） | **红 rc=1**，名册含 `bridge.go:1345(var-decl-no-initializer)` | `nail-rig-2-decl.log` |
| RIG-3 同文件别名 | `type rig309r2Alias = agent.ToolRequest` ＋ `rig309r2Alias{...}` | **绿 rc=0＝瞎**（v1 `nail-3-alias.log`） | **红 rc=1**，名册含 `bridge.go:1348(alias-literal)` | `nail-rig-3-alias.log` |
| RIG-4b `internal/agent` 里第二条（同包裸名） | 树内新文件 `internal/agent/zz309r2second.go`：`func rig309r2Second() ToolRequest { return ToolRequest{TaskID: "", Name: "fs.read"} }` | ⛔ 解析（整包被排除，结构上必绿） | **红 rc=1**，名册 `loop.go:675 ｜ internal/agent/zz309r2second.go:4(同包裸名)`，分母 293→**294** | `nail-rig-4b-second-in-agent.log`，树 `/tmp/wisp309r2-second.5EjiSC` |
| RIG-5 省类型的嵌套元素 | `return []agent.ToolRequest{{TaskID: "", Name: "fs.read"}}` | 绿 | **绿 rc=0＝具名漏形**（⛔ 假装响） | `nail-rig-5-elided.log` |
| RIG-6 跨包别名 | 树内 `internal/tools/zz309r2mid.go`＝`type Mid = agent.ToolRequest`，`cmd/wisp/zz309r2use.go`＝`_ = tools.Mid{TaskID: "", ...}` | 绿 | **绿 rc=0＝具名漏形**，分母 293→**295** | `nail-rig-6-crosspkg-alias.log`，树 `/tmp/wisp309r2-alias.Gljfv5` |

要点三条：
1. **红键在分子⛔ 在分母**：RIG-1/2/3 三发红的时候分母都还是 293（往⛔ 被排除的文件里追加一枚函数而已），
   只有真加文件的 RIG-4b/6 才动分母（294/295）⇒ 钉响的键值是"构造点枚数"，⛔ "文件数变了"。
2. **RIG-4b 是这一格的落点本身**：旧钉结构上⛔ 可能响它（`internal/agent` 整包⛔ 解析），新钉响且具名到文件。
3. **两发绿是故意的**：它们钉的是"本 matcher 的射程边界"，边界内的三形各有一发红；
   钉头那句"★这把尺仍⛔ 是全覆盖"就是靠这两发读数撑住的，⛔ 写成"⛔ 漏"。

## 台面与⛔ 越界

- RIG-4b／RIG-6 打的是**仓外导出树**（`git archive 6f389389` ⇒ 计数 995 枚 `.go`、四 root 非测试源 293 枚先证非空），
  再 `cp` 进本腿那枚⛔ commit 的钉与 `bridge.go`。⛔ 往 `internal/agent` 建任何文件（那⛔ 是本腿写面）。
  ⚠ 这一条越格自报进 `90-final.md`：树里的钉＝工作树副本⇄`6f389389` 的其余部分，**混面**，读数时已具名。
- 工作树面收尾：`git status --porcelain -- internal/tools/bridge.go`＝只 ` M`（本腿注释面），
  `logs/prov-rigs.stdout` 末段与 `logs/nail-rigs.stdout` 末段都留了同一枚名册。
