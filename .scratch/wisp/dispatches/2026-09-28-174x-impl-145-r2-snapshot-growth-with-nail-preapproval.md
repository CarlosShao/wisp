# 派单 145-r2（写腿）— 把"有来源"的那几枚字段真填进面板快照；⚠ 我已在派单里预先裁掉两枚会挡你的名册钉

- 时间：2026-09-28 17:4x +08　编排者锚（**不等于你的锚**）：`318d4bb3`
- 工单：票 145（`.scratch/wisp/issues/145-panel-snapshot-has-four-fields-so-eleven-of-the-fourteen-ui-states-have-no-input-to-render-grow-the-go-side-carrier.md`）——**先读它的 `AC#1` 那张表**（已交的普查，哪几维"有源"写在里面）＋三格未勾：`AC#2`／`AC#6`／`AC#2b`
- 解冻面（台账 **`A273`** 那批批准，owner 原话「批了」，撤销口令「145 别动」）：`internal/panel/composer.go`（`PanelSnapshot` 那四字段与生产构造点）／`internal/panel/pump.go`／`cmd/wisp/panel_pump.go`，**批准只到"新增字段"为止**——不改语义、不删字段、不动 `bridge.go:42-45` 既有方法名。
- ⚠ **本轮我另加一枚具名解冻（理由＝`A388`，撞钉预检现量出来的）**：`internal/panel/pump_test.go` 的 **`:124` 与 `:276`** 那两条断言，允许从"键集合恰好等于四枚"改成"**恰好等于 `PanelSnapshot` 当前声明的那一集**（由结构体推导、不写死名单）"。**除此两行之外那枚文件一字不动**；改完必须自带正控：**手工少发一枚键 ⇒ 那两条要红**（没有正控＝把门换成好看）。

## 要落的三格

1. **`AC#2` 落地有源那一集**：按 `AC#1` 表里"有 Go 侧真来源"的那几维给 `PanelSnapshot` 加字段＋在**生产构造点**填真值。⚠ **构造点是谁、几处，要在证据件里点名**（`grep` 出真 `file:line`，不接受"应该在哪"）。
2. **`AC#6` 反向判据**：每一枚新字段答一句"**哪一枚用例断言了它的值来自真来源**"。**答不出的字段＝装饰品，从落地集里拿掉**并在表里写明为什么拿掉——这一条优先于"多补几维"。
3. **`AC#2b` 两维**：① 当前可选模型清单（来源＝配置目录 ∩ `enabled=true`，⚠ **不许把未启用的也列进去**）；② 那一节的第二维照票面原文落。**若票面这两维与 `AC#1` 的"有源"结论矛盾，停手报回**，不要自己改票面。

## 硬禁区（一条破即判失败）

- **`frontend/**`／`design/**` 零写面**。TS 侧 interface 对齐**不在本程**：只在证据件里给出"**应当长这样**"的逐键清单（键名／类型／含义／哪枚 Go 字段供它），**由 owner 自己带给他用的那枚 agent**。**不要提"转给谁"**。
- 冻结件不动：`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`frontend_hygiene_test.go`／`go.mod`／`go.sum`／`thresholds.go`／golden／`allowlist.txt`／C18 超时常量／`docs/PLAN.md`／`docs/specs/**`／`bridge.go` 那四枚常量。
- 今天**在册红的三枚**（`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`）**不修、不当绿、不豁免**；⚠ 你加字段**很可能让第一枚更红**（它拿 Go 结构体去对 `frontend/**` 那份声明，而那份归别人）——**这是预期后果，照实记进证据件就行，不许为了它去动前端，也不许放宽它**。
- 别人的在飞／脏件不碰：`internal/risk/**`、`internal/tools/pointer_185_cli_seam_test.go`、`cmd/wisp/panel_inbound*.go`（刚落的片 B）、`design/**` 那批未提交删除、`probes/152`、`probes/161/r6/logs/flip-*`。
- **不新增任何 `panel.*`／方法名**（名册补齐属票 194 堆1）；**不做审批出口**（面板侧永不发起 allow 那一侧）。

## Git 纪律

只 commit **绝不 push**；每次 commit **显式 pathspec**；禁 `git add -A`／`.`／`-a`；禁 `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`；**仓内零删除命令**（临时件只建不删；恢复跟踪件用 `git cat-file blob HEAD:<path> > <path>`）；**每完成一格 commit 一次**（撞轮次上限时损失只有半张表）。

## 起手必做（含两条"现量再说"）

1. `date`＋`git log --oneline -1`＋`git status --porcelain`：起手名册逐枚抄进证据件第一行。
2. 现量并写进表：`grep -c "" internal/panel/composer.go`、`PanelSnapshot` 当前字段枚数（用结构体数，别抄票面）、`grep -rn "PanelSnapshot{" --include=*.go internal/ cmd/ | grep -v _test.go`（**生产构造点枚数**）、`go test -count=1 ./internal/panel/`（＝恰三枚在册红）。⚠ **任何一条与票面/派单不符就停手具名报回**——不许按自己的判断替我改规格或改判据。
3. **别用 `go test -overlay` 造变异载体**：会扫真实磁盘／`dist` 的尺看不见 overlay（`A382`；本仓 `internal/panel` 还经 `frontend` 的 `all:dist` 嵌入，仓外副本要带 `dist`，`A386` 后记）。载体建在**仓外**，仓内零删除。

## 门禁与终态

`sh scripts/d22scan.sh`（只看 clean／红点逐名；`ban #8 internal/=N` 那格是**文件枚数**、不是违规数，涨≠漂移）；`bash .scratch/wisp/probes/154/gate-clauses.sh`（只比 **BAD 腿名册**＝今天在册**只 `G6neg`**；**绝不跑 `probes/161/r6/flip-declaration.sh`**）；`go test -count=1 ./internal/panel/ ./cmd/wisp/`；`"$(go env GOPATH)/bin/gofumpt" -l internal/panel/ cmd/wisp/` 空。跑 CLI 相关测试带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`。终态 `git status --porcelain` **等于起手名册**（逐枚具名差集为空）。

## 落盘与回报

- 证据件：`docs/evidence/s1/145-snapshot-growth-r2.md`（骨架先建、逐格 commit；含"应当长这样"那张前端清单）。
- 票面：`AC#2`／`AC#6`／`AC#2b` 三格**都不许勾**（勾要非实现者表）；把落了几维、拿掉几维、构造点枚数追加进 Progress log。
- 回报（≤350 字，六栏）：① 落了哪几维／拿掉了哪几维（**每维答"真来源在哪一行"**）；② `pump_test.go:124/:276` 改法与**正控那一发红在哪条**；③ `PanelSnapshot` 字段枚数起手→终态＋生产构造点枚数；④ 变异／正控逐发"哪行→哪条红"；⑤ 门禁读数原样＋你被拒过的每次调用；⑥ **没测到什么**＋"离界面那十四态能真画还剩哪几维没源"。
- 预算：硬顶 **30**；第 5 枚内落第一枚 commit。**超预算不是放宽断言的理由**（今天两程各超 16／17 枚，账都记在案上）。
