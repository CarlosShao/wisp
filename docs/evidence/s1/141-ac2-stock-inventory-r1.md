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

## 2. 小结② —— 按码点段：owner 要选「补宽到哪一段」，这是各段的价

### 2.1 先说一件会改变这张表形状的事：**四段不够，缺第五段**

简报要的四段是「箭头 `U+2190–U+21FF`／带圈数字 `U+2460–U+24FF`／方框绘制 `U+2500–U+25FF`／`U+27C0–U+2AFF`」，
票面也写「整段 `U+2190–U+25FF`（箭头、方框绘制、带圈数字）」。
**但 `U+2190–U+25FF` 中间夹着一整段没人点名的 Unicode 区块：`U+2200–U+22FF`（数学算符）。**
本仓存量里有 **10 枚落在这一档**（`≥` `≤` `−` `∩`），**而且外露在 UI 上的那几枚全在这里**（§2.4）。
所以四段之和 **≠** 缺口段全量，下表把第五段一并列出。

### 2.2 实测分布（锚点 `99263cc`；命令＝`git grep -oP '[\x{lo}-\x{hi}]' 99263cc -- <树>`，逐段跑）

**A. `internal/` + `cmd/` 的 `.go`（票面的射程）—— 合计 121 处／116 行／41 枚**

| 码点段 | 处 | 行 | 文件 | 本仓这段的实际字形（处数） |
|---|---|---|---|---|
| `U+2190–U+21FF` 箭头 | **34** | 32 | 15 | `U+21D2` 18 · `U+2192` 16 |
| `U+2200–U+22FF` 数学算符（**第四段之外的漏项**） | **4** | 4 | 4 | `U+2265` 2 · `U+2264` 1 · `U+2229` 1 |
| `U+2300–U+245F` 杂记/技术符号 | 0 | 0 | 0 | —— |
| `U+2460–U+24FF` 带圈数字 | **83** | 80 | 26 | `U+2461` 32 · `U+2460` 20 · `U+2462` 15 · `U+2464` 7 · `U+2465` 6 · `U+2466` 2 · `U+2467` 1 |
| `U+2500–U+25FF` 方框绘制＋块元素＋几何形状 | **0** | 0 | 0 | —— |
| `U+27C0–U+2AFF` 补充箭头/数学 | **0** | 0 | 0 | —— |

**B. `frontend/`（票面漏量、但 ban #8 真走的一棵，`everyFile`）—— 合计 921 处／25 行／10 枚**

| 码点段 | 处 | 行 | 文件 | 实际字形 |
|---|---|---|---|---|
| `U+2190–U+21FF` 箭头 | **3** | 3 | 3 | `U+2191` 1 · `U+2192` 2（**全在 `/**` 注释里**） |
| `U+2200–U+22FF` 数学算符 | **6** | 6 | 4 | `U+2264` 4 · `U+2212` 2（**6 枚全是渲染文本**，见 §2.4） |
| `U+2460–U+24FF` 带圈数字 | 0 | 0 | 0 | —— |
| `U+2500–U+25FF` 方框绘制 | **912** | 16 | 8 | `U+2500` 912（8 枚 `.tsx` 的 `/* ── … ── */` 分节线，16 行 × 约 57 枚） |

**C. `design/`：0**（`git grep` 四段全 0）。PLAN.md:3574-3575（`5d777f7` 版，09-19 09:11）**逐字点名的就是这棵树**，它今天真干净。

### 2.3 候选射程的爆炸半径（**这就是三支的代价对比**，同一锚点实测）

`run '<pattern>'` ＝ `git grep -{l,n,o}P` 于 `99263cc`，pathspec **＝ ban #8 的全部四棵**（`internal/` `cmd/` `frontend/` `design/`）：

| 若把 `emojiRe` 补宽到…… | 点亮文件 | 行 | 处 | 备注 |
|---|---|---|---|---|
| 只补 `U+2190–U+21FF` | 18 | 35 | 37 | 箭头＝台账/注释习惯，几乎不动语义 |
| 只补 `U+2200–U+22FF` | **8** | **10** | **10** | **最小的一段，却是唯一含"已经外露"的那一段** |
| 只补 `U+2460–U+24FF` | 26 | 80 | 83 | 存量最大（占 A 段 69%），但**零枚外露** |
| 补 `2190–21FF + 2460–24FF`（＝四段里非零的两段） | 41 | 115 | 120 | **漏 10 枚**（全在数学算符段），且**漏掉外露的那几枚** |
| 补 `2190–21FF + 2200–22FF + 2460–24FF` | 47 | 125 | 130 | 漏 `U+2500` 那 912 枚（`frontend/` 分节线） |
| 补 `U+2500–U+25FF` | 8 | 16 | 912 | **处数最大的一段**，全在 8 枚 `.tsx` 的注释分节线里 |
| **补 PLAN 的字面全射程**（缺口＝`2190–25FF` ∪ `27C0–2AFF`） | **51** | **141** | **1042** | 口径最硬；`tools/` 那一枚仍点不亮（§3.2） |

> 读法提示：**「处数」与「要改的文件数」不成比例**——补 `2500–25FF` 一段就多 912 处但只 8 枚文件，
> 补 `2200–22FF` 一段只 10 处却牵动 8 枚文件里的**生产文案**。按"处"比代价会误判。

### 2.4 段内**外露**明细（与 §1 的两档交叉，owner 判"哪一段真有后果"只看这张）

| 段 | 已外露于面板/CLI 的存量 | 出处与**那一版** |
|---|---|---|
| `U+2200–U+22FF` | **7 处在真外露面上**：① `frontend/src/components/composer.tsx:170` 的 **JSX 文本节点** `支持 {state.acceptedAttachmentMimes.join(" / ")}，单个 ≤ {bytes(state.maxAttachmentBytes)}` —— 面板正在渲染这枚 `U+2264`；② 同形的 fixture `frontend/fixtures/composer-states.html:2,5,8` 的 `<span class="text-ink-3">…单个 ≤ …B</span>`（3 处）；③ `thinking.tsx:213`、`tool-chips.tsx:186` 的 JSX 文本 `<span className="text-red">−{row.del}</span>`（2 处 `U+2212`）；④ `internal/risk/rules_scale.go:24` 的 `U+2265` → 审批卡（§1.3，1 处）。**合计 `U+2264`×4 ＋ `U+2212`×2 ＋ `U+2265`×1 ＝ 7 处** | `composer.tsx`／`composer-states.html` ＝ **`8e10095`（09-21 19:31）**；`thinking.tsx`／`tool-chips.tsx` ＝ **`9318bb8`（09-21 13:05）**；`rules_scale.go` ＝ **`67ffbd8`（09-20 08:22）** |
| `U+2190–U+21FF` | 生产文案 **4 行**：`fs_write.go:346,480`（步骤账本→回执「已执行」）＋ `:485,567`（`Result.Text`→回灌上下文）。`frontend/` 那 3 处箭头**全在注释里，不外露** | `fs_write.go` ＝ `ed74595`（09-21 10:58）；`approval-card.tsx`／`task-rows.tsx` ＝ `9318bb8`（09-21 13:05） |
| `U+2460–U+24FF` | **0 枚外露**：83 处全部落在注释与测试诊断文案里（表 A 可逐行核） | —— |
| `U+2500–U+25FF` | **0 枚外露**：912 处全在 8 枚 `.tsx` 的 `/* */` 分节线注释 | —— |
| `U+27C0–U+2AFF` | 0（本仓无存量） | —— |

**⇒ 一张数（不是建议）：`Q-46` 三支里任何一支若射程不覆盖 `U+2200–U+22FF`，**已经在面板/审批卡上外露的这 7 处**继续扫不到；
若要覆盖它，牵动的**生产文件只有 4 枚**（`frontend/src/components/composer.tsx`、`ai-native/thinking.tsx`、`ai-native/tool-chips.tsx`、`internal/risk/rules_scale.go`）＋ 1 枚 fixture（`frontend/fixtures/composer-states.html`），
**其中 `internal/risk/rules_scale.go` 在禁改清单里**，另 4 枚在 `frontend/`（**不在**简报给的禁改清单上，但 `frontend/` 是 owner 自己未提交改动的邻域——本程没动、也没量它是否还有别的未提交态，见 §5）。

---

<!-- 回执 5dfb8ba -->
```
$ git log --oneline -1
5dfb8ba docs(141 AC#2 盘点 r1 第2节): 小结① 非注释 38 行两档 = (i)31 行/7 枚 + (ii)7 行/4 枚; 判据 = 沿拼接续行爬回语句首行的 sink pattern(可复跑), 并给改判路径(i)35/(ii)3; 生产外露 5 行逐条点名到 fs_write.go:70-72 / tool.go:40-43 / panel/approval.go:47,83。
$ git show --name-only HEAD

docs/evidence/s1/141-ac2-stock-inventory-r1.md
```

---

## 3. 小结③ —— 冻结件逐枚，与**每枚各需一张什么样的具名解冻**

**先纠一笔账**（§0 已提，这里定死）：**41 枚里在禁改清单的是 11 枚，不是 12 枚。**
第 12 枚 `tools/d22scan/scan_test.go:274`（1 行／1 处 `U+2460`，注释）
① 不在票面取数的 pathspec（`internal/**.go` ＋ `cmd/**.go`）里，所以**它不在 41 之内**；
② 更关键：`emojiScopes()`（`tools/d22scan/main.go:480-487`，该文件最近一版 `5e8f87b` 09-21 15:36）只声明
`design/` `frontend/` `internal/` `cmd/` 四棵，**没有 `tools/`**
⇒ **三支选项里没有任何一支会因它变红；给它开解冻是白开。** 若 owner 要把仪器自家测试也纳入射程，那是**第四件事**（扩 scope），票面没写。

### 3.1 十一枚的逐枚表（**解冻要具名，所以每行是一张独立的、范围封闭的授权描述**）

| # | 文件（最近一版，`git log -1` 实取） | 命中行 | 行/处 | 所在段 | 位置 | 需要的具名解冻（**只此范围，不外扩**） |
|---|---|---|---|---|---|---|
| 1 | `internal/risk/provenance.go`（`f342413` 09-21 08:25） | 13,41,42,146,469,499,622,628 | 8/9 | `U+2460/2464/2465` 带圈 | **8 行全注释** | 「仅限该文件上述 8 行**注释内**的字形替换；不得触碰 `//` 之外任何字节、不得改可执行语句」 |
| 2 | `internal/risk/provenance_test.go`（`42ed13f` 09-21 13:18） | 165,483,527,564 | 4/4 | 带圈 | 全注释 | 同 1 形状，**4 行**；虽在 `_test.go` 但只动注释 ⇒ 与断言无关 |
| 3 | `internal/risk/syncdirs_other_test.go`（`42ed13f` 09-21 13:18） | 27 | 1/1 | `U+2467`（**本仓唯一一枚**，票面分布表漏的两枚之一） | 注释 | 同 1 形状，**1 行** |
| 4 | `internal/risk/pathresolver_expansion_test.go`（`0117459` 09-21 17:16） | 17 | 1/3 | `U+2192` 箭头 | 注释 | 同 1 形状，**1 行**（该行 3 处） |
| 5 | `internal/risk/pathresolver_anchor_spelling_windows_test.go`（`6a6c85e` 09-21 10:32） | 200 | 1/1 | `U+2229`（**数学算符段**，票面分布表漏的另一枚） | 注释 | 同 1 形状，**1 行**；⚠ 若射程不含 `U+2200–U+22FF`，**这一枚压根不需要解冻** |
| 6 | `internal/risk/rules_scale.go`（`67ffbd8` 09-20 08:22） | 24 | 1/1 | `U+2265` | **非注释·生产文案** | ⚠ **最重的一张**：解冻须写明「允许改 `reason: fmt.Sprintf(…)` 的**格式串**」，而**这会改变审批卡上显示的文案**（§1.3 的 verbatim 链）；且须**与第 7 枚同批**（期望值耦合，§1.5） |
| 7 | `internal/risk/assessor_test.go`（`67ffbd8` 09-20 08:22） | 154 | 1/1 | `U+2265` | **非注释·表驱动期望值** | 「允许改 `cases[]` 内 `want: Decision{… Reason:"…"}` 一枚字面量」；**改它＝改断言的输入**，AC#4 的"放水"两判据（断言有没有被动／helper 是不是原有的）在这里**必须逐条回答** |
| 8 | `internal/winsec/reparse_windows_test.go`（`01e7007` 09-21 16:34） | 34,47,113,138,146,232,261（注释）＋ **121,126,236,245（`t.Logf`）** | 11/11 | `U+2460/2461` | **7 注释 ＋ 4 非注释** | 一张里两栏：「7 行**注释内**替换」**＋**「4 枚 `t.Logf` 文案里的 `A51②` 台账标号替换」。后者不参与判定（`t.Logf` 不断言），但**会改 `-v` 输出文本**，AC#4 的名册差集要看 |
| 9 | `internal/winsec/seam_guard_windows_test.go`（`c02c609` 09-21 19:31） | 284 | 1/1 | 带圈 | 注释 | 同 1 形状，**1 行** |
| 10 | `internal/winsec/winsec.go`（`1499efe` 09-21 20:57） | 4,74 | 2/2 | 带圈 | 注释 | 同 1 形状，**2 行** |
| 11 | `internal/winsec/winsec_windows.go`（`391878a` 09-21 21:17） | 186,639 | 2/2 | 带圈 | 注释 | 同 1 形状，**2 行** |
| — | `tools/d22scan/scan_test.go`（`5e8f87b` 09-21 15:36） | 274 | 1/1 | `U+2460` | 注释 | **不需要解冻**（不在 ban #8 射程，见本节开头）。要它变红得先扩 `emojiScopes()`——**那是另一张票，本表不替它开口** |

**11 枚合计：33 行／36 处。** 其中**非注释 6 行**（第 6、7 各 1 ＋ 第 8 的 4），
⇒ **若走"注释豁免"那一档，11 枚里只有 3 枚需要动断言侧**（`rules_scale.go`、`assessor_test.go`、`reparse_windows_test.go`），
另 **8 枚整枚只需注释级解冻**——这是"解冻会不会膨胀"的关键数：**8 张开成"仅注释内"就封住了，3 张必须明写到行。**

### 3.2 段 × 冻结的交叉（**决定"补哪一段要开几张解冻"**）

| 若射程含 | 牵动的冻结枚数（11 之内） | 需具名解冻的文件 |
|---|---|---|
| 只 `U+2190–U+21FF` | **1** | `pathresolver_expansion_test.go`（注释 1 行） |
| 只 `U+2200–U+22FF` | **3** | `rules_scale.go` ＋ `assessor_test.go`（**必须同批**）＋ `pathresolver_anchor_spelling_windows_test.go` |
| 只 `U+2460–U+24FF` | **7** | `provenance.go`、`provenance_test.go`、`syncdirs_other_test.go`、`reparse_windows_test.go`、`seam_guard_windows_test.go`、`winsec.go`、`winsec_windows.go` |
| 只 `U+2500–U+25FF` | **0** | ——（该段在 `internal/`+`cmd/` 无存量） |
| `U+27C0–U+2AFF` | **0** | —— |
| PLAN 字面全射程 | **11** | 上三行并集 |

> **另注一件与本节同形的账**：`frontend/` 那 10 枚**不在禁改清单**上（简报的三条禁改路径是
> `internal/risk/**`、`internal/winsec/**`、`tools/d22scan/**`），
> 所以补宽后 `frontend/` 的 25 行／921 处**不需要任何解冻**——它是 51 枚里最便宜的一大半，
> 但也是唯一含"已经外露"文案的一棵（§2.4）。

---

---

## 4. 顺手要核的两件（都成立，但都比简报说的更宽）

### 4.1 `AGENTS.md:38` 与 PLAN 一致、与仪器不一致 —— **成立**，外加一处措辞岔口

**`AGENTS.md:38` 原文逐字**（该文件最近一版 **`5866c6f`（09-24 09:58）**，取法 `git show HEAD:AGENTS.md \| sed -n '34,42p'`，第 38 行＝那 5 行里的第 5 行）：

> `- UI 代码里出现任何 emoji（\`U+2190–U+2BFF\`、\`U+1F300–U+1FAFF\`、\`U+FE0F\`）。`

**PLAN.md 两处**（`git blame` 逐行：`3447`／`3448`／`3574`／`3575` 四行都出自 **`5d777f7`（09-19 09:11）**，且 `^` 标记说明它是该文件的边界版；文件最近一版 **`45623e4` 09-24 11:01**，工作树对 `docs/PLAN.md` 干净）：

- `docs/PLAN.md:3446-3448`：「⚠ **文字符号也算 emoji 范畴**：对勾用 SVG `check`，箭头用 SVG `chevron-right`，警告用 SVG `alert-triangle`。**这一条对抗验收可扫描**（源码里出现 U+2190–U+2BFF、U+1F300–U+1FAFF、U+FE0F 即判违规）。」
- `docs/PLAN.md:3574-3575`：「**零 emoji 可机器验证**：`design/` 下所有文件扫描 Unicode 区块 U+1F300–U+1FAFF / U+2190–U+2BFF / U+FE0F / U+2600–U+27BF，**命中数必须为 0**」
- `docs/PLAN.md:3443`（同版）：「**绝对禁止**：任何 emoji（含 `✓ ✔ ✗ ⚠ ★ →` 这类**文字符号**）」——**该行实测码点**：`U+2713 U+2714 U+2717 U+26A0 U+2605 U+2192`。前五枚都在 `2600–27BF`（**被扫到**），最后一枚 `U+2192` 在缺口段（**扫不到**）⇒ 同一张"绝对禁止"清单内部，两半执行强度不同。**票面这句我复现了，成立。**

⇒ **判定：AGENTS.md:38 的三组范围与 PLAN 逐字一致，与 `main.go:111` 不一致。按 AGENTS.md 自己的让位条款，缺陷归仪器侧。**

**新增读数（简报没说的一半）**：三处对**"哪一片代码"**的说法也互不相同——

| 出处 | 措辞 | 实际覆盖 |
|---|---|---|
| `AGENTS.md:38` | 「**UI 代码**里」 | 最窄：字面只指 UI |
| `PLAN.md:3447` | 「**源码**里出现…即判违规」 | 全部源码 |
| `PLAN.md:3574` | 「**`design/` 下所有文件**」 | 只管那棵树 |
| `main.go:480-487`＋`:794-799` | `design/` 全文本、`frontend/` 每文件、`internal/`＋`cmd/` 全部 `.go`，**注释与 `_test.go` 都算** | **最宽**，且宽过 PLAN 点名的树 |

⇒ 薄索引在**范围**上跟仪器不一致、在**表面**上比仪器窄，两处都朝"让门禁看着比实际规格松"的方向偏。本程只登记，不改（`AGENTS.md` 在禁改清单）。

### 4.2 仪器比 PLAN **多扫** `1F000–1F0FF` —— **成立，而且多扫的是连续三大段**

`tools/d22scan/main.go:111` **原文逐字**（`git blame -L 111,111` → **`64d083d8`（09-20 07:40）**；文件最近一版 `5e8f87b` 09-21 15:36）：

```
emojiRe = regexp.MustCompile(`[\x{1F000}-\x{1FAFF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}\x{FE0F}\x{1F1E6}-\x{1F1FF}]`)
```

PLAN 侧的**并集**（3447-3448 ∪ 3574-3575）＝ `2190–2BFF` ∪ `1F300–1FAFF` ∪ `FE0F` ∪ `2600–27BF`（最后一项已是 `2190–2BFF` 的子集）。

| 差集方向 | 实测区间 | 里面是什么 |
|---|---|---|
| **仪器 − PLAN**（仪器多扫） | **`1F000–1F2FF`（一整段连续）** | `1F000–1F0FF` 麻将／骨牌 · `1F100–1F1FF` 带圈字母补充＋区域指示符 · `1F200–1F2FF` 带圈表意补充 |
| **PLAN − 仪器**（仪器漏扫） | **`2190–25FF`** 与 **`27C0–2AFF`** | 箭头／数学算符／杂记技术／带圈数字／方框绘制／块元素／几何形状 —— **本仓全部存量都在这两段里** |

**三条推论**（都是数出来的，不是推的）：

1. **这是两段各自手工枚举、不是一处笔误**：仪器自己就把 `1F1E6–1F1FF` 单列了一行，而它是 `1F000–1FAFF` 的**真子集**——冗余项的存在证明写正则的人**拼过两张清单**。⇒ **"补宽"不是单向加一段**，要逐段对齐（对齐后 `1F000–1F2FF` 这段要么保留、要么明写为什么保留）。
2. **多扫那段今天零成本**：`1F000–1F2FF` 在 ban #8 的四棵射程内**实测 0 处**；全仓只有 2 处，都在 `docs/reports/pending-and-issues.md:649,652`（`U+1F195`），而 **`docs/` 根本不在 ban #8 射程**（`emojiScopes()` 无 `docs/`）。⇒ 若 owner 选"双向对齐"（既补又裁），**裁掉多扫段不会让任何现有文件变红**。
3. **漏扫段与"当前已经扫到的段"完全不相交**：`2600–27BF`／`2B00–2BFF`／`FE0F`／`1F300–1FAFF` 四段在四棵射程内**全部 0 处**。⇒ 门是绿的**不是因为存量干净，是因为射程与存量不相交**——票面这句我复现成立，且比我预期更彻底：**仪器今天匹配的字符数是 0。**

---

<!-- 回执 64664d6 -->
```
$ git log --oneline -1
64664d6 docs(141 AC#2 盘点 r1 第4节): 小结③ 冻结件逐枚 11 枚/33 行/36 处 + 每枚一张范围封闭的 具名解冻描述; 纠账——41 之内只有 11 枚冻结, 第 12 枚 tools/d22scan/scan_test.go 不在 ban#8 射程(emojiScopes 无 tools/)故三支都点不亮; 段×解冻交叉表: 只补箭头 1 张、只补数学 3 张、只补带圈 7 张、全射程 11 张; 11 枚里只有 3 枚需要动断言侧。
$ git show --name-only HEAD

docs/evidence/s1/141-ac2-stock-inventory-r1.md
```

---

<!-- 下一节 -->



