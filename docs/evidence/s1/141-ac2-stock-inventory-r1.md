# 141 AC#2 —— ban #8 缺口段存量盘点（只读程 r1）

**锚点**：`99263cc`（本程开工自量：`git rev-parse --short HEAD` → `99263cc`，分支 `dev`）
**agent**：`auditor-ticket141-inventory-r1`
**性质**：**只读盘点**。本程**没有运行任何仪器**：零次 `go build` / `go test` / `go vet` / `gofmt` / `gofumpt` / `docker` / `wisp` / `wisp slo` / `scripts/slo-check.ps1`，零次 `gh`。全部读数是 `git grep` / `git show` / `git blame` / `sed -n` / 自写的**文本分类脚本**（在 `D:\tmp\wisp141inv\`，不在仓内）。同机 `acceptor-ticket136-ac14-r1` 正在取真读数，本文件不与任何读数争时段。
**只做 AC#2**（纯盘点，不依赖 `Q-46`）。`AC#1` 的选支**不由本程决定，本程也不建议**——见文末「我不替 owner 选支」。

> **本节末尾回执的约定**（写在最前面，免得后一位以为缺了）：共树不得 `--amend`，
> 所以**每一节正文末尾贴的是「上一枚 commit」的回执**；本节自己的回执由下一节带来，
> 最后一节的回执在文末「commit 流水」汇总。所有回执都是 `git log --oneline -1` ＋
> `git show --name-only HEAD` 的原样输出。

---

## 0. 复核：简报给的四个数，与两张要补的账

### 0.1 复核实到（票面用 `git grep -P '[\x{2190}-\x{25FF}]' -- 'internal/**.go' 'cmd/**.go'`，我逐条重跑）

| 简报的数 | 我的复量 | 判定 |
|---|---|---|
| 41 枚 `.go` 文件 | **41 枚**（`git grep -lP` 于 `99263cc`，实测见 §5 命令） | ✅ 对得上 |
| 116 行 | **116 行** | ✅ |
| 121 处（逐字符） | **121 处** | ✅ |
| 注释 78 行 / 非注释 38 行 | **78 / 38**，但**判据被换掉了**（见 0.2） | ✅ 数值对，方法不采信 |
| 非注释散在 11 枚文件 | **11 枚** | ✅ |
| **12 枚在禁改清单** | **41 枚里只有 11 枚**（`internal/risk/` 7 ＋ `internal/winsec/` 4）。第 12 枚 `tools/d22scan/scan_test.go` **根本不在 41 枚的取数 pathspec 里**，而且**不在 ban #8 的 walk scope 里** ⇒ 三支选项里**没有任何一支会因它变红** | ❌ **推翻** |
| 今日新增 0 | `git diff -U0 182daed..HEAD -- '*.go'` 加号行 gap 命中 **0**；`-- 'frontend/'` 亦 **0**（`182daed` 存在，09-24 15:04） | ✅ |

**票面字形分布少列了两枚**：票面「`②`32／`①`20／`⇒`18／`→`16／`③`15／`⑤`7／`⑥`6／`⑦`2／`≥`2／`≤`1」加起来是 **119**，不是它自己写的 121。实测 12 枚码点、总和 121，漏的是 `U+2467` ×1（在 `internal/risk/syncdirs_other_test.go`）与 `U+2229` ×1（在 `internal/risk/pathresolver_anchor_spelling_windows_test.go`）。**总数没错，分布表缺两行。**

### 0.2 推翻的方法：注释／非注释不是「前面有没有 `//`」

票面写「我按"命中位置前有没有 `//`"切的」。那个切法在两处会错，本程不沿用：

1. **它认不出块注释**，也会把**串内的 `//`**（如 `"https://…"`）当成注释起点；
2. **它认不出 raw string（反引号串）里的 `--` SQL 注释**——那在 Go 语法里是字符串，按票面切法会被算成"非注释"，但它在 SQL 语义里**正是注释**。

本程改用一枚**Go 迷你词法器**（`D:\tmp\wisp141inv\classify.pl`，状态机：`code`／`c-line`／`c-block`／`s-dq`／`s-raw`／`s-rune`，处理转义与跨行 raw string），逐 rune 记录命中时刻的状态。取数方式：`git cat-file blob 99263cc:<file>` 把 41 枚文件抽到 `D:\tmp\wisp141inv\tree\`（**不在仓内**）后分类，所以读的是**锚点那版**、不是工作树。

词法器输出 116 行／121 处，与票面**数值全等**，状态分布为：

| 词法状态 | 行数 | 含义 |
|---|---|---|
| `c-line` | **78** | `//` 行注释 |
| `s-dq` | **36** | 双引号字符串字面量内 |
| `s-raw` | **2** | 反引号 raw string 内（**都是 `internal/memory/schema.go` 的 `ddlV1` DDL 文本里的 SQL `--` 注释**） |
| `c-block` / `code`（标识符） | **0** | —— |

⇒ **两个结论**（都影响 owner 的代价估算）：
- **38 枚非注释命中全部在字符串字面量里，一枚都不在标识符里**。简报说「测试名 `func Test…` 里的箭头」这一档**在本仓实测为 0**：Go 标识符不能含这些字形，票面上那 20 枚 `→`/18 枚 `⇒` 没有一枚在函数名位。
- 那 2 枚 raw string 命中是**「Go 眼里是字符串、SQL 眼里是注释」**。若 `Q-46` 走 (c)「注释豁免」，这 2 行**落在哪一边需要 owner 单独认一句**——词法上讲它们在字符串内（不豁免），语义上讲它们是 SQL 注释（该豁免）。本程只把这个岔口标出来。

### 0.3 补的第二张账：票面的 pathspec 只圈了 ban #8 射程的一半

`emojiScopes()`（`tools/d22scan/main.go:480-487`，`git blame` 该函数体所在版：`5e8f87b` 09-21 15:36 是这枚文件最近一次改动）声明 ban #8 **实际走四棵树**：`design/`、`frontend/`（`everyFile`）、`internal/`、`cmd/`（后两棵 `goOnly`）。票面的取数 pathspec 只有 `internal/**.go` ＋ `cmd/**.go`，**漏了 `frontend/`**——而 `frontend/` 恰是这条禁令**动机上真正要管的那棵树**（面板就是 UI）。实测（同一 `99263cc`）：

| 树 | 文件 | 行 | 处 | 是否 ban #8 walk scope |
|---|---|---|---|---|
| `internal/` + `cmd/`（`.go`） | 41 | 116 | 121 | 是（`goOnly`） |
| `frontend/` | **10** | **25** | **921** | **是**（`everyFile`） |
| `design/` | 0 | 0 | 0 | 是——PLAN.md:3574 真点名的那棵树，**今天确实是 0** |
| `tools/`（含 `tools/d22scan/`） | 1 | 1 | 1 | **否**（`emojiScopes()` 里没有 `tools/`） |

⇒ **补宽后的真实爆炸半径是 51 枚文件／141 行／1042 处**，不是 41／116／121。差额里 `frontend/` 一枚就贡献 921 处（96% 的字符量），且 `frontend/` 的 25 行里有 **6 行是渲染文本**（`≤`、`−`），不是注释——详见 §1 与表 B。

---

## 1. 小结① —— 非注释那 38 行的两档

### 1.1 判据（**不是**"看起来会不会显示"）

一条命令、一张 pattern 表，可复跑：

```
perl /d/tmp/wisp141inv/sink.pl <文件> <行号>
```

它先沿**字符串拼接续行**（上一行以 `+` 结尾）爬回**打开这枚字面量的那条语句**，
再按下面的 pattern 顺序给它贴 sink 标签。**分档规则 = 这枚字面量是不是「要被写出去的东西」**：

| 标签 | pattern（语句首行） | 档 |
|---|---|---|
| `test-print` | `t\.(Fatalf\|Errorf\|Logf\|Fatal\|Error\|Log\|Skipf\|Skip)\(` | **(i)** 进 `go test -v` 输出＝CI 可见 |
| `render-field` | `(\w*(Text\|Reason\|reason))\s*:\s*fmt\.(Sprintf\|Sprint)` | **(i)** 进面板/回执 |
| `audit-ledger` | `se\.record\(` | **(i)** 进步骤账本（渲染证据见 1.3） |
| `struct-field` | 仅 `(\w*(Text\|Reason\|reason))\s*:`（右边不是 `fmt.`） | **(ii)** 输入 fixture 或表驱动期望值 |
| `parse-template` | `Sscanf\(` | **(ii)** 解析模板 |
| `NONE` | 上面都不沾（raw string 里的 SQL DDL 文本） | **(ii)** |

### 1.2 结果：**(i) 31 行／7 枚文件**，**(ii) 7 行／4 枚文件**（合 38／11，与 §0 全等）

**(i) 会进可见输出 —— 31 行**

| 文件 | 行 | sink | 谁把它写出去 |
|---|---|---|---|
| `internal/ball/tokens_test.go` | 203,206,216,220,224,228,234,237,247,278（**10 行**） | `t.Errorf` | `go test -v` |
| `internal/tools/fs_staging_windows_test.go` | 66,73,157,235,246,337,365,371（**8 行**） | `t.Fatalf` | 同上 |
| `internal/winsec/reparse_windows_test.go` | 121,126,236,245（**4 行**） | `t.Logf` | 同上（`-v` 下**通过时也打印**，比 Fatalf 更外露） |
| `internal/tools/fs_write.go` | 346,480（`se.record`）、485,567（`Text: fmt.Sprintf`）（**4 行**） | 生产渲染 | 见 1.3 |
| `internal/agent/approval/ticket97_alias_direction_test.go` | 119,132（**2 行**） | `t.Fatalf` | `go test -v` |
| `internal/tools/bridge_a18_kill_windows_test.go` | 220,222（**2 行**，同一枚 `t.Fatalf` 的拼接续行，语句首行在 218） | `t.Fatalf` | `go test -v` |
| `internal/risk/rules_scale.go` | 24（**1 行**） | 生产渲染 | 见 1.3 |

**(ii) 只在源码字符串里 —— 7 行**

| 文件:行 | sink | 为什么落 (ii) |
|---|---|---|
| `internal/llm/anthropic/cache_test.go:21,22,23` | `struct-field:Text` | `llm.TextPart{Text: "①identity…"}` 是**喂给 mock 的请求输入**，不是消息文案 |
| `internal/risk/assessor_test.go:154` | `struct-field:Reason` | 表驱动 `want:` **期望值**，与生产串比对用 |
| `internal/tools/fs_write_test.go:521` | `parse-template` | `fmt.Sscanf` 的**解析模板**，且整个 `if` 被 `; false` 守卫短路 |
| `internal/memory/schema.go:29,118` | `NONE`（raw string） | `ddlV1` 的 **SQL DDL 文本**，送进 SQLite，不是任何输出流 |

### 1.3 生产那 5 行的外露证据（逐条点名，不用推测）

- `internal/tools/fs_write.go:346`／`:480` 的 `se.record(...)` → `sideEffect.steps`；`fs_write.go` **最近一版 `ed74595`（09-21 10:58）** 的 `record` 函数自带注释逐字写着「**the report renders these verbatim under 「已执行」**」（`internal/tools/fs_write.go:70-72`，同上一版）。
- `internal/tools/fs_write.go:485`／`:567` 的 `Text: fmt.Sprintf(...)` → `tools.Result.Text`，定义处 `internal/tools/tool.go:40-43` **最近一版 `0986d63`（09-20 19:16）** 注释逐字：**「Text is the payload the host puts back into the context (the loop caps and spills it, D15(3) …)」** ⇒ 回灌上下文＝会进对话呈现。
- `internal/risk/rules_scale.go:24`（最近一版 `67ffbd8` 09-20 08:22）的 `reason: fmt.Sprintf("…（`U+2265`%d）")` → `risk.Decision.Reason` → `internal/panel/approval.go:83` `Reason: decision.Reason`；该文件**最近一版 `63ef895`（09-21 13:26）**，`:47` 注释逐字：**「Reason is Decision.Reason verbatim」** ⇒ **审批卡原样渲染这枚 `U+2265`**。这是 38 行里唯一一枚**已经在真面板上外露**的缺口段字形。

### 1.4 两档的**边界声明**（owner 若要换判据，数会怎么变）

`struct-field` 那 4 行**并非绝对不外露**：
- `internal/risk/assessor_test.go:154` 的期望值，在断言不等时由 `:178-181` 的 `t.Fatalf("%s:\n  want %+v\n  got  %+v", tc.name, tc.want, got)` **整块打印**（`internal/risk/assessor_test.go` 最近一版 `67ffbd8` 09-20 08:22）；
- `internal/llm/anthropic/cache_test.go:21-23` 的输入 fixture，在解析失败时被 `:54`／`:65` 的 `t.Fatalf("…\n%s", body)` / `("system field: %v\n%s", …)` **连带请求体整块打印**（该文件最近一版 `f342413` 09-21 08:25）。

⇒ 若 owner 把 (i) 定义为「**任何可能到达终端的字面量**」，这 4 行改判 (i)，分布变成 **(i) 35 ／ (ii) 3**（只剩 Sscanf 模板 1 ＋ SQL raw 2）。
本程**按 sink 位取 (i) 31／(ii) 7**，并把改判路径写清楚——**不替 owner 决定用哪条判据**。

### 1.5 顺手量的两枚**耦合**（清存量时会互相拽一下）

- `internal/tools/fs_write.go:485`／`:567` 与 `internal/tools/fs_write_test.go:521` **看着像镜像、其实不是**：生产侧回收站文案（`fs_write.go:417`，`ed74595` 09-21 10:58）写的是「已放入回收站：%s（还原记录 %s，可在回收站还原）」——**不含缺口段字形**；测试侧那枚模板却写着「…（条目数 %d `U+2192` %d」。⇒ 这是一枚**陈旧模板**，与生产串不同步。清理它不需要与生产同批改，但**动它就在动一枚断言的输入**（AC#4 的"放水"判据会看这个）。
- `internal/risk/rules_scale.go:24`（生产）与 `internal/risk/assessor_test.go:154`（期望值）是**同一枚 `U+2265` 的两半**：改生产不改期望＝断言红。**必须同批**。两枚文件都在禁改清单（§3）。

---

<!-- 回执 ade897c：本节开头贴上一节那枚 commit 的原样输出 -->
```
$ git log --oneline -1
ade897c docs(141 AC#2 盘点 r1 第1节): 复核实到 41/116/121 全等, 推翻两处—— "12 枚在冻结路径"实为 11(第12枚 tools/d22scan/scan_test.go 不在取数 pathspec 也不在 ban#8 walk scope); 票面 pathspec 漏了 frontend/(ban#8 真射程), 补宽爆炸半径是 51 枚/141 行/1042 处。注释切法改用 Go 词法器, 38 枚非注释 全在字符串字面量内、0 枚在标识符。
$ git show --name-only HEAD
    全在字符串字面量内、0 枚在标识符。

docs/evidence/s1/141-ac2-stock-inventory-r1.md
```

---

<!-- 下一节 -->

