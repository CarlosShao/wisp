# 309-r1 · `AC#2` 枚数钉（`internal/tools/toolrequest_nail_309_test.go`）

## 定形（编排者裁语：枚数钉，⛔ 依赖任何未来接线；⛔ 新造生产可达路径）
断言＝**`agent.ToolRequest{` 在 `internal/agent` 之外的生产构造点枚数 == 0**，且**分母非零**（自带分母是本仓既有规矩，抄 `task_state_188_test.go` 那句 *"a scan whose denominator is zero reports clean for a tree it never read"*）。

## 实现（只用范本已有的两把形状，⛔ 自造遍历、⛔ 去 `cmd/wisp` 找样板）
- 读源码：`filepath.WalkDir` over roots＝`filepath.Join("..","..","cmd")` 与 `filepath.Join("..","..","internal")`（`go test` 工作目录＝包目录 `internal/tools`，故模块根＝`../..`；只读、⛔ import、⛔ 调用）。
- 逐文件 `parser.ParseFile` + `ast.Inspect`，按**限定类型选择器** `agent.ToolRequest` 的 `CompositeLit` 计数（⛔ 靠变量名/helper 名匹配、⛔ 字符串 grep，改名不骗过尺、真新构造点藏不住）。
- 排除：`_test.go`（含 `reflect.TypeOf(agent.ToolRequest{})` 那两枚）；`internal/agent/**`（那枚可信的裸 `ToolRequest{` 构造点在 `loop.go`，同包不写限定名、且落排除项）。
- 返回 `(constructions, examined)`；`examined<10 ⇒ t.Fatalf`（分母不可信＝假绿）。名册/排除项/射程目录全写进失败句，下一位可逐字重跑。

## 现量基线（两把独立尺，写进用例红句）
- 枚数尺：`grep -rn 'ToolRequest{' internal cmd --include=*.go` 全仓 **60**；非测试仅 **1**＝`internal/agent/loop.go:675`；`reflect.TypeOf(agent.ToolRequest{})` **2**（都在 `_test.go`）。⇒ 排除测试与 `internal/agent` 后＝**0**。
- 门尺：票 154 `gate-clauses.sh` 的 **G1** 登记 `want_n 0`，今日对得上（⛔ 在 309 重登记 want_n）。

## 实测（原始读数＝`logs/nail-rig.log`）
| 台件 | constructions | 分母 examined | 期望 | 实测 |
|---|---|---|---|---|
| BEFORE（无额外构造点） | 0 | 245 | 绿 | `--- PASS` rc=0 |
| Rig（在 `internal/tools/` 内我写面里临时加一枚非测试 `agent.ToolRequest{TaskID:""}` 生产文件） | 1 | 246 | **红** | `--- FAIL` rc=1（枚数钉响句命中，点名"要么接非空 id，要么进 `AC#1` 覆盖入口"） |
| 还原 | — | — | 逐字节 | `SCRATCH_REMOVED=yes`；`git status --porcelain -- internal/tools` 只剩两枚 `??`（我的正式测试文件） |

- 内置正控（不碰盘、⛔ 生产路径）：`parser.ParseFile` 读内存里一段 `agent.ToolRequest{TaskID:""}`，`countToolRequestLits==1` 否则 `t.Fatalf("the ruler is blind")` ⇒ 证明绿不是"扫了 0 个然后报绿"。

## 越面判定（重要）
枚数钉**读** `cmd/**` 与 `internal/**`（含跨包源码），但**写**只在 `internal/tools/*_test.go`：跨包读源＝只读、范本 `task_state_188_test.go` 已确立（它 `WalkDir ../panel`）⇒ ⛔ 越写面、⛔ 为覆盖 `cmd/wisp` 而把测试写到别处。**未触发停手上报**：`cmd/wisp` 现量零构造点，且它在扫描 roots 内被覆盖到。`AC#2` 闭合。
