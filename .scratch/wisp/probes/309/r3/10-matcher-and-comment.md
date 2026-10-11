# `309-r3` (A)(B)：改了什么，逐枚 文件:行

判语（哪格"成立"）⛔ 本腿的活，归 `309-v3`。这里只交"改了什么＋用了哪几条 AST 形状＋仍漏哪几形"。

## 1. 改动名册（尺＝`git show --name-only` 逐笔；两笔代码笔＝`df8d517d`(A) ＋ `a1f84d39`(B)）

| # | 文件:行（改后现量） | 改了什麼 | 属 (A) 还是 (B) |
|---|---|---|---|
| 1 | `internal/tools/toolrequest_nail_309_test.go:17` | 头注 `(2) matcher` 那条：`309-r2 认三形` → 加一行"309-r3 再接三形，现在共认六形"，并指回 `20-teeth.md` | (A) 注释面 |
| 2 | 同件 `:55-80` | 射程登记块重排：新增 `:57-63`"309-r3 接进射程的三形"（逐枚带 shape 名与台件名）；`:64-79` 保留仍漏的形逐枚具名；`:80` 那句"判据是恒真形必须响、⛔ 类型级完美（那要引 `go/types`）"**原句未动** | (A) 注释面 |
| 3 | 同件 `:127-135`（`toolRequestNamesIn309` 内注释） | 旧句"定义型那形算残漏"改成"定义型本身⛔ 进名集，够着 `agent.ToolRequest` 靠那枚转换，转换**已**被计（`callShape309`）" | (A) 注释面 |
| 4 | 同件 `:187-224` | **新函数 `callShape309`**（含 doc）——三形判定全在这一枚 | (A) 代码面 |
| 5 | 同件 `:264-274` | `ast.Inspect` 的 switch 里新增 `case *ast.CallExpr:` 一支，把 `callShape309` 的结果按同一枚 `nail309Point` 结构计进分子 | (A) 代码面 |
| 6 | 同件 `:246-262` | `findToolRequestPoints309` 的 doc：旧句"Two shapes count"换成"Six shapes count（三枚 309-r2、三枚 309-r3）"，并指回名册 | (A) 注释面 |
| 7 | 同件 `:426-465` | `controls` 名册 6 枚 → **18 枚**（新增 12 枚；其中两枚是"漏形今天真的计 0"的反向控制）；控制体⛔ `t.Run` ⇒ 整包名册零新增（读数见 `20-teeth.md` 末节） | (A) 测试面 |
| 8 | 同件 `:480` | 件长 373 → **480** 行（尺＝`wc -l internal/tools/toolrequest_nail_309_test.go`，工作树面） | (A) |
| 9 | `internal/tools/bridge.go:959-972` | 那句过头的 "counts **every** production construction point … so a second one rings **wherever** it is written" 换成实测射程段（认哪七枚 shape 标签／明说哪几形⛔ 响／指回钉的登记／收尾"那是债、这句是射程、⛔ 是全覆盖"） | (B) 注释面 |

`git diff --numstat 3974ab4d a1f84d39`（本腿两笔代码笔，尺＝`git diff --numstat`，射程＝这两件）＝
`14 4 internal/tools/bridge.go` ＋ `120 13 internal/tools/toolrequest_nail_309_test.go`。

## 2. (B) 那一处是不是"只改注释"——三把尺（读数＝`logs/bridge-comment-scope.log`、`logs/bridge-comment-scope-2.log`）

| 尺 | 读数 |
|---|---|
| `git diff -U0 -- internal/tools/bridge.go` 里非 `//` 开头的增删行（`grep -vE '^[+-][[:space:]]*//' \| grep -c`） | **`noncomment_changed=0`**（名册为空＝0 枚） |
| 两版（HEAD blob ⇄ 工作树）各 `sed -E '/^[[:space:]]*\/\//d'` 删掉整行注释后 `cmp` | **IDENTICAL**，`854 ⇄ 854` 行，`sha256=496ed6d683e8e070…`（⇄ `309-v2` 报的同一枚 sha） |
| 语义字节逐枚 `grep -c` | `if b.injected \|\| b.prov == nil \|\| taskID == "" {`＝1⇄1、`WithTaintDetector(b.prov.Detector(taskID))`＝1、CR 字节数＝**0**（⛔ 把 `bridge.go` 编辑成 CRLF 那一族病本腿没犯） |

⇒ `assessorFor` 的行为、`internal/risk/**`、`lit236*`、`agent.ToolRequest` 字段：本腿一枚⛔ 动。

## 3. 新射程用了哪几条 AST 形状（只 `go/ast`＋`go/parser`，⛔ 引 `go/types`）

`callShape309`（件 `:195`）三条判定，按顺序：

1. **转换形** `conversion-to-tool-request`：`*ast.CallExpr` ∧ `len(Args)==1` ∧ `denotesToolRequest309(call.Fun, names)`。
   `Fun` 可以是 `agent.ToolRequest`（`*ast.SelectorExpr`）、同文件别名 `T`、或 `package agent` 里的裸名 `ToolRequest`。
   ⇒ 这就是定义型 `type X agent.ToolRequest` 那形今天唯一真造出 `agent.ToolRequest` 值的写法（件 `:196-199`）。
2. **`new` 形** `new-of-tool-request`：`Fun` 是 `*ast.Ident` 且 `Name=="new"` ∧ `denotesToolRequest309(Args[0], names)`。
   ⛔ 读参数名以外的一切（件 `:209-213`）。
3. **`make` 形** `make-slice-of-tool-request`：`Fun.Name=="make"` ∧ `Args[0]` 是 `*ast.ArrayType` ∧ `arr.Len == nil`（切片**类型**，⛔ 定长数组）∧ `denotesToolRequest309(arr.Elt, names)`。
   ★长度实参⛔ 被读 ⇒ `make(..., 0)` 与 `make(..., n)` 同一形（票面 §1(A) 要的那枚"含 `make(..., 0)` 那种长度写法"，件 `:214-221`）。

枚数不重复计：`*ast.CompositeLit`（旧形）与 `*ast.CallExpr` 是两种节点，`agent.ToolRequest{…}` ⛔ 是 CallExpr；
`var rig = new(agent.ToolRequest)` 只由 `new` 计一枚（`GenDecl` 那一支要求 `len(vs.Values)==0`）——两枚"no double count"控制各钉一条（件 `:424-425`、`:445-447`）。

shape 标签枚数尺（＝(B) 那句"seven named shape labels"的分母）：

```
grep -oE '"(qualified|same-package-bare|alias)-literal"|"var-decl-no-initializer"|"(new-of|make-slice-of|conversion-to)-tool-request"' internal/tools/toolrequest_nail_309_test.go | sort -u | wc -l
=> 7
```
七枚逐字＝`alias-literal` / `conversion-to-tool-request` / `make-slice-of-tool-request` / `new-of-tool-request` / `qualified-literal` / `same-package-bare-literal` / `var-decl-no-initializer`
（读数＝`logs/bridge-comment-scope.log` 末节；⚠ 另有枚 `literalShape309` 的 `default` 分支返回 `"literal"`，结构上到不了——`denotesToolRequest309` 只放行 SelectorExpr/Ident，⛔ 把它数进"七"）。

## 4. 被拒的形：照原样留在名册里（逐枚引 309-r2 交付版原文行号）

原文行号＝**改前**那件（尺＝`git show 3974ab4d:internal/tools/toolrequest_nail_309_test.go` 的行号）；现行号＝工作树改后。

| 原文行 | 原句（逐字，节首） | 改前判定 | 现状态 | 现行 |
|---|---|---|---|---|
| `:53` | `· 别名再套一层落到**别的包**（pkg q 里 type Y = agent.ToolRequest，第三方写 q.Y{...}）＝漏；` | 漏 | **仍漏**，⛔ 装覆盖（实测绿读数＝`309-r2` 的 RIG-6，`logs/nail-rig-6-crosspkg-alias.log` rc=0）；本腿另加一枚同名反向控制（计 0） | `:65` |
| `:54` | `· 函数体内声明的别名（本 matcher 只读文件级 type 声明）＝漏；` | 漏 | **仍漏**（matcher 只遍历 `file.Decls` 的 `token.TYPE`） | `:66` |
| `:55` | `· 改名 import（import ag "…/internal/agent" 后 ag.ToolRequest{...}）与 dot-import ＝漏；` | 漏 | **仍漏**（`denotesToolRequest309` 硬看 `pkg.Name == "agent"`） | `:67` |
| `:56` | `· new(agent.ToolRequest)／make([]agent.ToolRequest, n)／反射或解码构造（json.Unmarshal 进结构体）＝漏；` | 漏（一里含三形） | ★**拆开**：前两形 309-r3 已接进射程（`:58-59`）；第三形"反射或解码构造"**留在名册**并写明"**没有静态尺抓得住**"＝真实残余 | `:58-59`、`:68-70` |
| `:57` | `· 定义型 type X agent.ToolRequest（⛔ 别名）再 agent.ToolRequest(x) 那枚转换＝漏；` | 漏 | **已接进射程**（`conversion-to-tool-request`，RIG-9 两发为证） | `:60` |
| `:58` | `· 省了类型的嵌套元素 []agent.ToolRequest{{TaskID: ""}} 的内层 {...}＝漏（外层那枚是 ArrayType，同样⛔ 计）⇒ 这一形实测两向都计 0，见 …nail-rig-5-elided.log；` | 漏 | **仍漏**，本票⛔ 要它，且加了一枚反向控制钉"今天真计 0"，防注释漂成覆盖 | `:71-72`、`:459-461` |
| `:60` | `· 裸名那一支按 **package 子句**放行⛔ 按目录放行 ⇒ 别处真有一枚叫 agent 的包时，它自己的 ToolRequest{} 会进分子（今天尺＝…读数 R359-2）。` | 漏 | **仍漏**（本腿⛔ 动那条放行） | `:73-75` |
| `:63` | `⇒ 判据是"恒真形必须响"，⛔ 类型级完美（那要引 go/types，本票⛔ 装那个复杂度）。` | 自陈口径 | **原句一字未动**（派单 §1(A) 要求保留） | `:80` |
| —（本腿新添自陈） | `· ★309-r3 自陈新增一枚：转换那一支是**句法**判定…一枚恰好同名、同包、又只吃一枚实参的**函数**调用也会被计进分子` | 新登记的**假阳**方向 | 具名在 `:77-79`，⛔ 装成类型级判定 | `:77-79` |

⇒ 名册枚数：改前 **7 条项／9 形**（台账 `A851` §2 那把尺）→ 改后 **6 条项／7 形仍漏**，其中 `:68-70` 那枚＝反射与解码**两形**（尺＝逐枚数上面那张表"仍漏"那一列：`:53`/`:54`/`:55`/`:58`/`:60` 各一形 ＋ `:68-70` 两形 ＋ 本腿自陈那枚**假阳**⛔ 数进漏形＝**7 形**；枚数⛔ 是常数，引之前先重跑）。
