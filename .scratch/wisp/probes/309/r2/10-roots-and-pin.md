# 309-r2 · 10 枚数钉的射程与落点（job 1＋job 2）

> ⚠ 本件与先落盘的 `10-nail-roots-and-pin.md`（随 `9d3f9da7`）**同题、只追加⛔ 取代**：那一枚是交回时的正文，本枚补齐它⛔ 展开的两半
> ——`scripts` 到底进⛔ 进的具名答复（两把尺的读数摆在这）＋每 root 分母与"空名册即红"那两支牙。临时件⛔ 删（工单池 `issues/README` 规则 8），
> 所以两件都在；`90-final.md` §1 那行的"件"名按这一对读。

件＝`internal/tools/toolrequest_nail_309_test.go`（⛔ 一枚名册新增：测试函数名⛔ 动，见 `50-gates-behaviour.md`）。

## 1. roots 的两把尺 ⇒ tools 进、scripts **也进**

尺＝`logs/roots-rulers.log`（R359-1/R359-3，工作树面，HEAD `202ed4a5` 之后、本腿代码笔之前）：

```
go.mod                      :: module github.com/CarlosShao/wisp
scripts/spike/go.mod        :: module github.com/CarlosShao/wisp/scripts/spike
tools/d22scan/go.mod        :: module github.com/CarlosShao/wisp/tools/d22scan
tools/mockllm/go.mod        :: module github.com/CarlosShao/wisp/tools/mockllm
-- 尺二: grep -rn "CarlosShao/wisp/internal" tools scripts --include=*.go
tools/signmodels/main.go:34:	"github.com/CarlosShao/wisp/internal/models"
```

**判据①（前缀）**：三枚 nested 模块的 module 路径逐字落在 `github.com/CarlosShao/wisp/` 之内。
Go 的 `internal/` 可见性是按 **import 路径前缀**判定的（`internal` 父亲那一截之内的包都可 import），
⇒ "它是独立模块所以进不来 `internal/agent`"这个说法⛔ 成立（编排者 §2 里 `nail-roots-miss-tools-dir` 那条自纠同理）。

**判据②（活例）**：`tools/signmodels/main.go:34` 逐字 import `github.com/CarlosShao/wisp/internal/models`
⇒ nested 模块伸手进 `internal/` 在本仓**已是事实**，⛔ 是我脑补的语言规则。

**`scripts` 到底要不要进 ⇒ 进。** 三条读数支撑：
1. 前缀尺：`scripts/spike` ＝ `github.com/CarlosShao/wisp/scripts/spike` ⇒ 在 `wisp/` 那一截之内 ⇒ 可见性**放行**它 import `internal/agent`。
2. 今天它自己的 import 名册＝0 命中（尺＝上面那把 grep 打在 `scripts`，⛔ 一行）；`find scripts -name '*.go' -not -name '*_test.go'`＝**11** 枚源。
3. ⛔ 用"此刻还没伸手"当排除理由＝把**接线事实**当**结构保证**，那正是本票在防的那族病（票面 `AC#2` 第一句就是这句）。
⇒ 判据与前两条同源，`scripts` 与 `tools` 同判，**两枚都进 roots**。今天各 root 的非测试源枚数：cmd 40／internal 229／tools 13／scripts 11 ⇒ 分母 **293**。

## 2. 排除支换成"具名文件＋具名枚数"

旧：`strings.Contains(rel, "internal/agent/")` ⇒ `internal/agent` 的 **24 枚**非测试源一枚⛔ 解析，
而重试包装器最自然的落点恰在那包里；且 `loop.go:675` 是**同包裸名** `ToolRequest{`，旧 matcher（只认限定名 SelectorExpr）本来就认不出它。
新：⛔ 排任何目录，断言**全 roots 生产构造点合计 == 1 枚**，且**具名**那枚落在 `internal/agent/loop.go`；名册⛔ 空也照打印（文件:行(形状) 逐枚）。

读数（`logs/r359-7-nail-roster.log`，钉自打印⇄第三把 `find` 尺对得上）：
```
AC#2 枚数钉：roots=[..cmd ..internal ..tools ..scripts] 里 agent.ToolRequest 生产构造点=1（==1，那枚＝internal/agent/loop.go），
分母=293 枚非测试源（每 root=cmd=40 internal=229 scripts=11 tools=13）
```
第三把尺（⛔ 用钉自己的数自证）＝`find <root> -name '*.go' -not -name '*_test.go' | wc -l` 逐 root＝40／229／13／11，和＝293＝钉自打印 ⇒ 分母是真算出来的。
旧钉那枚 245 的来历也留在钉头（＝269−24），两把尺的口径差别⛔ 隐蔽。

新增两枚牙（⛔ 一枚空名册当"扫过了"）：
- 每 root 解析数⛔ 零：为 0 直接 `t.Fatalf("root %q 一枚非测试源也没解析：空名册等于扫过了，这一格没交")`。
- 全局地板从 `<10` 抬到 `<200`，注释写明它**是地板⛔ 是常数**（今天 293；将来有人把某枚 root 从名册里划掉就会响）。

## 3. 与票 154 那扇门的口径关系（⛔ 重记、⛔ 改它）

`gate-clauses.sh` 的 **G1** 仍登记 `want_n 0`（尺＝`git grep -nE 'ToolRequest\{' ... ':!internal/agent/*'`＝"internal/agent **之外**几枚"）。
本钉现在问的是另一问（"全 roots 共几枚、那枚在哪个文件"）。两问今天对得上：G1 侧＝0，本钉侧＝1 且那枚在 agent 里。
⛔ 动 `gate-clauses.sh` 一字（G1b 那格归票 154 那一族，编排者 §U-4 已裁）。
