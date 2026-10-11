
### 编排者复跑＋裁定（2026-10-11 14:1x +08，收 **`309-v1`**＝一笔 `47d3f6da`，件 `.scratch/wisp/probes/309/v1/**`，裁决表 `docs/evidence/s1/309-v1-verdict.md`；台账 `A847`；⛔ 上面各节原句一字不改，本节追加）

★**本轮翻两格＝`AC#4`＋`AC#5`**（尺后＝**3 勾／3 未勾**，尺＝`grep -cE '^[[:space:]]*- \[x\]'` ⇄ `- [ ]`）。`AC#1`／`AC#2`／`AC#3` **不翻**，条件全部是**真残余**且修法便宜，落成一枚小落地腿 **`309-r2`**（范围见末节）。

**1. 我自己复跑的（⛔ 一枚⛔ 抄腿的数；件＝`probes/orch/20261011/logs/orch-A847-*.txt`）**

| 承重读数 | 腿报 | 我这一发 | 判 |
|---|---|---|---|
| 枚数钉分母 | 245（自打印） | 第三方尺＝tracked 非测试 `.go`（`internal`＋`cmd`）**269** − `internal/agent/` **24** ＝ **245** | ✅ 对得上 |
| `bridge.go` 换行归一那次自抓自纠 | `CR=0`／`i/lf w/lf attr/text eol=lf` | `tr -dc '\r' \| wc -c`＝**0**；`git ls-files --eol`＝`i/lf w/lf attr/text eol=lf` | ✅ |
| 三门禁 rc | 各自落行、件⛔ 0 字节 | `GOFLAGS= go build ./...`＝**rc=0**、`go vet ./internal/tools/`＝**rc=0**（两件输出皆 0 字节＝无诊断）；`sh scripts/d22scan.sh` 我这发 13:31＝**rc=0**；腿那枚 `gate-d22scan.log`（23078 字节）里**真带 live 门禁那句** `clean - no D22 ban violations`（我先怀疑它只装了自测输出，量完才收） | ✅ |
| ③ 拒绝文案的出处 | 票 197 正文 0 命中＝在产码 | `grep -c '没有宿主给的任务 id' .scratch/wisp/issues/197-*.md`＝**0**；`internal/tools/subagent_197.go`＝**1** | ✅ 收下那枚措辞更正 |
| 形状洞的先例 | `probes/154-accept/ac1-g1-blindspot-synthetic.txt` 已量过五形 | `ls` 该件在场（另有同族 `ac1-gate-clauses-rerun-by-accept.txt` 等） | ✅ ⇒ **⛔ 另立新账**，引用即可 |
| `AC#3` 那把自指尺 | 注释自称 1、实测 3 | 逐字跑注释自己写的那把（`grep -rn 'ToolRequest{' internal cmd --include=*.go` 减 `_test.go`）＝**3 行**：`internal/agent/loop.go:675`（真构造点）＋**`bridge.go:879` 与 `:940` 两行注释自己**；再剥行首 `//`＝**1** | ✅ 条件成立＝本格⛔ 翻 |

**2. 十枚欠账逐枚签收（收／驳都写理由，⛔ 整批收下）**

- **收**：`Q-309-v1-selfref-ruler`（上表那发，阻 `AC#3` 的"带尺重述"那半句）／`Q-309-v1-nail-misses-decl-and-alias-shapes`（声明＋赋值、别名两形我读 `countToolRequestLits` 的 AST 匹配码自证＝只认 `*ast.CompositeLit` 且 `Type` 是 `agent.ToolRequest` 那个 `SelectorExpr`，`var` 声明与 `X{}`（别名＝`*ast.Ident`）两形**结构上必漏**，⛔ 跑腿的 Rig 也认）／`Q-309-v1-nail-excludes-whole-agent-package`（排除支是 `strings.Contains(rel,"internal/agent/")`＝**整包 24 枚源不解析**，而重试包装器最自然的落点就在那包里；且 `loop.go:675` 是**裸名** `ToolRequest{`，那把 matcher 本来就⛔ 认——这条比我裁语里那句"外部落点会响"更窄，我改口）／`Q-309-v1-stale-count-comment`（钉头复述的 60／59／1，同一把尺今日＝**69／66／3**，增量＝本票那两枚件自己；我这一发逐枚数过）／`Q-309-v1-prov-nil-limb-blind`（两肢⛔ 仪器）／`Q-309-v1-other-gates-uninstrumented`／`Q-309-v1-nail-control-not-exercisable`（内置正控的红⛔ 可测＝它⛔ 可能盲，腿改用扫描面三形代证区分度，我签收这个替代）／`Q-309-v1-anchor-debt-154`＝**⛔ 本票账**，归票 154/151 那一族（那 5 处绝对 `loop.go:N` ⛔ 在 `AC#3` 点名段里，腿具名停手＝对）。
- **驳一枚＝`Q-309-v1-who-rings-overclaim`**。腿指 `bridge.go:947` 那句"Who rings if this limb's shape changes"⛔ 全称。我读了原句：它写的是 **this limb**（单数＝`taskID == ""` 那一肢），并且同段末句自己写了 "The compiler guards nothing here - that test is the only instrument, and this comment states the debt, not a guarantee"。⇒ **⛔ 算 overclaim**（射程与自陈都⛔ 越）。⚠ 但它把"这肢⛔ 是唯一一枚"这句逼出来了＝`prov == nil` 那肢的欠账仍按上一行收下，两件事⛔ 混一格。
- **★我自己差点判错的一枚＝`Q-309-v1-nail-roots-miss-tools-dir`**。我第一反应是"可以关掉：`tools/` 是独立模块，Go 的 `internal/` 可见性规则⛔ 让它 import `internal/agent`"。**落笔前核了两件事**：①`tools/d22scan/go.mod` 的 module 路径＝`github.com/CarlosShao/wisp/tools/d22scan`＝**在 `internal` 父亲那一截之内**⇒ 可见性规则**放行**；②本仓**已有活例**＝`tools/signmodels/main.go:34` 逐字 import `github.com/CarlosShao/wisp/internal/models`。⇒ **这枚⛔ 关、收进 `309-r2`**（把 `tools` 加进 roots，`scripts/spike` 那枚独立模块同问）。★记我：这正是 `⛔ 拿"我记得的语言规则"当尺`——与 `A846` 那枚假清洁同一族（**先证前提，再判死活**）。

**3. 承重那句我复述一遍（⛔ 让下一位以为裁完就闭）**

形 Ⓑ 之后**对空 `taskID` 的实际防御＝零**（腿三把尺坐实＋我自己一把：代码行 `cmp` 同一／两枚 236 空 id 载体绿／整包名册纯新增）。盘上只有**两枚跳线**：一枚定向仪器（只认"挂没挂探测器"那一件事的两向）＋一枚枚数钉（只认 `agent.ToolRequest{` 限定名字面量、`internal/agent` 之外、`cmd|internal` 射程之内）。⛔ 一格⛔ 读成"以后有人递空 id 会有人管"——这句是本票的**边界**，⛔ 是缺陷清单。

**4. `309-r2`（落地腿）范围＝把上面三枚"条件"关掉，五件、逐件带尺**

1. 枚数钉**roots** 加 `tools`（＋具名回答 `scripts` 那枚独立模块要不要进，判据＝module 路径前缀＋有没有同类 import 活例）；**排除支**从"整包 `internal/agent/`"改成**具名文件＋具名枚数**＝断言"全模块 `agent.ToolRequest`／同包裸名 `ToolRequest` 生产构造点合计 **1** 枚且那枚在 `internal/agent/loop.go`"⇒ 有人往 `internal/agent` 里再加一条也会响。
2. matcher 补两形：**声明＋字段赋值**（`var req agent.ToolRequest` … `req.TaskID = ""`）与**类型别名**（`type X = agent.ToolRequest` 后 `X{}`）；补不上的形状**具名写在注释里**（⛔ 装成全覆盖）。
3. 注释三处说实话：`bridge.go:879-881` 那把自指尺（自带"剥行首 `//`"那一支，或改用解析器尺）；钉头 `:9-11` 那组过期数（60／59／1 ⇒ 现量 69／66／3，并写明"这把尺会把注释自己算进去"）；`:935` 那句 "Every other gate still runs"（结构真、断言⛔ ⇒ 改成说实话的射程句）。
4. `b.prov == nil` 那一肢：我裁＝**⛔ 造第二枚仪器**，照 `b.injected` 那行的既有格式**在注释里另记一句**（同族同处置：两肢都属"装配坏了"，票面现量早已把 `injected` 归装配族另裁）；欠账 `Q-309-v1-prov-nil-limb-blind` 由那行注释＋台账结掉，⛔ 留成"没人管"。
5. `AC#4` 那套门禁照旧：三门禁 rc 自落行、格式名册两把并排各写射程、显式 pathspec、票尾一条 bullet（1 加／0 删）。**⛔ 动框**＝`AC#1`/`AC#2`/`AC#3` 翻不翻归 `309-v2`（非实现者）。

**Status（本程更新）**：`AC#0`（裁形）＋`AC#4`（门禁与越界）＋`AC#5`（三本账对账）＝**三格闭合**；`AC#1`/`AC#2`/`AC#3`＝**附条件成立、⛔ 翻**，条件＝上面 §2 那几枚，落点＝`309-r2` ⇒ 之后 `309-v2` 终裁。⛔ 零 push。撤销口令＝逐枚「**309-v1 撤 `AC#n`**」（腿的判语）与「**309 我翻的 AC#4 撤**」「**309 我翻的 AC#5 撤**」（我的翻框）；改判形 ⓐ 仍用「309 改形 ⓐ」，届时本票五格判语整批作废。

- [2026-10-11 14:1x] agent=编排者 did=收 `309-v1`（一笔 `47d3f6da`，72 枚件、越界 0、`.go` 0）⇒ 六把承重尺自己复跑（分母 245＝269−24／eol 三把／三门禁 rc／票 197 那句 0 命中／自指尺 3⇄1／154 先例件在场）＋**翻两格＝`AC#4`＋`AC#5`**、`AC#1`/`AC#2`/`AC#3` 判"附条件成立"⛔ 翻；十枚欠账逐枚签收＝**收 9 驳 1**（驳＝`who-rings-overclaim`，原句写的是 this limb 且同段自陈"⛔ guarantee"）；★**我自己差点把 `nail-roots-miss-tools-dir` 当"语言规则⛔ 可达"关掉**——现跑 `tools/d22scan/go.mod` 的 module 前缀＋`tools/signmodels/main.go:34` 那条活例后才改判＝收下；落 `309-r2` 五件（roots 加 tools／排除支改具名文件＋枚数／matcher 补声明＋别名两形／三处注释说实话／prov-nil 那肢按 injected 同族另记）；⛔ 动别家票、⛔ 动产码语义、⛔ push
