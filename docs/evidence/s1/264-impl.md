# 票 264 实现腿证据表（形甲：咽喉一处出站脱敏，2026-11-07）

**腿**：票 264 实现写腿（264-impl）。**裁定来源**：编排者铁律令——甲形（模型可见副本出站掩路径），取形已由机主/编排者拍（HEAD `763963b87` 的 commit 标题即载「264甲」）。**基线 HEAD**：`763963b873b2d18cc704f636064a55be7ec6e09e`（派单锚 `4656ec715` 为其祖先）。**本腿零 push**；`docs/reports/pending-and-issues.md`、`frontend/**`、`design/**` 一字未碰。

## 0. 落点（咽喉一处，非 103 处逐点）

| 件 | 内容 |
|---|---|
| `internal/agent/pathmask.go`（新增） | `maskOutboundPaths(text)`：Windows 盘符绝对路径 → `<path>`；含 174 指针豁免、URL 防误伤（`https://` 前缀字节判别）、续段吸收规则 |
| `internal/agent/loop.go` `toolResultParts`（唯一改动点） | 该函数是**每一条模型可见工具答案的必经咽喉**：正文 TextPart 与两支 imageDropNotice 均先过 `maskOutboundPaths` 再入史。用户侧（`EvToolEnd` 事件 :846-848、`res.ToolLog`）在此调用**之前**已从 `log.Text` 取值，天然不掩 |
| `internal/agent/pathmask_264_test.go`（新增） | 规则表 + 咽喉单元钉 + e2e（harness 真 golden 流）+ 174 指针保命钉 |

未动的 103 枚装配点：一处咽喉覆盖全部，逐点改等于给每个新工具重开泄漏口（264-ac0 §5 本腿技术裁定的原话）。

## 1. 掩码规则（一句话可引）

**从盘符绝对路径（`X:\` 或 `X:/`，且不是长词尾巴如 `https://`）开始、到文件名非法字符（空白、`:`、`：`、`，`、`。`、`；`、`、`、引号、`<`、`>`、`|`、`*`、`?`）为止的整段，整体替换为稳定占位 `<path>`；段内以空白接续的分段仅当其自身含分隔符（`\` 或 `/`）才吸收——故「段中间带空格的路径」整段掩掉，而路径后的错误正文（`…: 系统找不到指定的文件。`）逐字保留。**

- 占位符：`<path>`（稳定、可判读；上线正文以 JSON HTML 转义形 `\u003cpath\u003e` 出现）。
- **唯一豁免＝174 指针**：紧跟 `全文见 ` 之后的路径不掩。那是宿主自己写的回执（`internal/agent/spill.go:189` 与 `internal/tools/task.go:629` 同形），是模型回到全文的唯一路——掩掉＝只截不指（PLAN.md:2564），174r8 钉会红。
- 已知边界（有意接受，写死在 pathmask.go 头注）：末段含空格的路径（`C:\a\new file.txt`）尾段不掩；POSIX `/usr/…` 不掩（D7 Windows-only，掩了会吃 URL 与时刻比例）。掩码只会替换，不会发明路径。

## 2. 三枚判据 → 读数（全部绿）

| 判据 | 测试 | 读数 |
|---|---|---|
| (a) 模型可见副本无机器目录结构 | `TestToolResultPathMaskedForModelButNotForUser264`（e2e：真 golden SSE 流 + 真 openai-chat adapter 上线正文）+ `TestToolResultPartsMasksModelCopyOnly264`（咽喉单元） | 历史 RoleTool 正文与第二发请求体均无 `C:\Users`/`operator`，含 `<path>`；错误正文 `系统找不到指定的文件。` 活着 |
| (b) 用户侧仍说得出是哪个文件 | 同 e2e 的 (b) 支 | `EvToolEnd` 事件 Text 与 `res.ToolLog[0].Text` 均逐字含 `C:\Users\operator\笔记\draft.txt` |
| (c) 174r8 指针钉原样通过 | `TestMaskKeeps174TaskOutputPointer264`（task.go:629/spill.go:189 同形回执逐字节不动）+ 实跑 `go test ./cmd/wisp -run TestTaskOutputPointerNotice174r8` | 单元钉 PASS；174r8 两枚端到端（SpeaksTheRealReason / StaysQuiet）PASS，回执逐字日志里 `全文见 C:\Users\swq\…\tool-output-call-174r8.txt…` 原样在 |

AC#2 两向：`TestToolResultPartsLeavesCleanText264`（无路径文本逐字节不动＝去掉必不响）＋规则表 `TestMaskOutboundPathsRule264`（10 形：patherror／C26 expands／rename 双路径／段中空格／路径后正文保留／全角冒号前缀／URL 不误伤／时刻不误伤／裸盘符不动／幂等）。

## 3. 门禁读数（本腿现跑，2026-11-07 09:2x–09:3x）

| 门 | 命令 | 读数 |
|---|---|---|
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l internal/agent` | 空（首跑检出测试表缩进 1 件，`gofumpt -w` 后复空） |
| gofmt | `gofmt -l internal/agent` | 空 |
| vet | `go vet ./internal/agent` | 干净 |
| 包测 | `go test ./internal/agent -count=2 -timeout 25m` | `ok … 11.760s`（两轮全绿） |
| D22 | `sh scripts/d22scan.sh` | **clean**（首跑检出本腿注释幻影引用 `tools/task.go` 一枚，改为 `internal/tools/task.go` 后复跑 clean；positive control 全 PASS） |
| 174r8 | `go test ./cmd/wisp -run TestTaskOutputPointerNotice174r8 -count=1` | `ok`（指针钉原样通过） |

环境：Git Bash；`PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`；`unset GOTMPDIR TMPDIR TEMP TMP`；`CGO_ENABLED=1`（全程未置 0）。未开 GUI、未跑 wisp.exe、未跑任何 *winlive* 测试。

## 4. 中途红过的两次（诚实账）

1. 首版 `pathSpanEnd` 从盘符字母起扫，把 `C:` 里的冒号当止点 → 掩成 `<path>:\Users…`（规则表 6 形红）。修：跳过 `X:/` 三字节前缀再扫。
2. 首版续段吸收在续段自身又是盘符路径时（`rename D:\a D:\b`）被其冒号截断。修：续段允许起手 `X:\` 前缀（`isDrivePrefix`）。
另：d22scan 首跑红一枚幻影引用（本腿注释把 `internal/tools/task.go` 写成 `tools/task.go`），改引用后 clean。

## 5. 写面与提交

- 写面：`internal/agent/pathmask.go`、`internal/agent/loop.go`、`internal/agent/pathmask_264_test.go`、`docs/evidence/s1/264-impl.md`（本件）。仅此四件入 commit，显式 pathspec，不 push。
- 共享树里他腿在飞件（`cmd/wisp/run.go`、`internal/ball/*`）未触碰、未入 commit。
