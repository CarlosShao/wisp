# 180 — `[panel] width` 是一枚**生产码零读者**的配置项：用户在 `config.toml` 里改它、或我们改 `default:"640"`，**界面上什么都不会变**（"提示存在但结果永不受影响"那一族在配置面上的复发）

- Status: **ready-for-agent（只读普查＋一支小落地）**。
- 来路：09-28 10:2x，前端那侧递给 owner 的"要拍板三件事"里第 ① 件（"`[panel] width` 默认 640 — 你要多少？改 `config.toml` 还是让 Go 改默认值？"）。我按盘上现量裁定：**这一件不该摆给他**——因为**两条路今天都没有可见效果**，先修效果再谈数值。台账 `A360`。
- 关联：票 92（composer 的 mode/附件/工作区那一格，同族"面板要真值不要自造"）· 票 145（快照扩成真载体）· 票 147／139／97（"字段存在、无人读"那一族）· `AGENTS §1.2` 与台账里那条"一条提示存在、但结果永远不受它影响的检查，比没有检查更坏"

## 现量（锚 `039efb47`，09-28 10:2x 本程现跑，每把尺可复制）

1. 字段在场：`internal/config/schema.go:527-528` 逐字
   `// Width is the panel width in px.` ＋ `Width int \`toml:"width" default:"640"\``（尺：`grep -rn "Width" --include=*.go internal/config/ | grep -v _test`）。
2. **读者＝零枚**：尺 `grep -rn "\.Width" --include=*.go internal/ cmd/ | grep -v _test.go` ⇒ **空输出**。
   ⇒ 今天**没有任何一条路径**把这枚数送到窗口尺寸上：改 `config.toml` 无效、改 `default` 也无效。
3. 它也不是"待接线的公开缺口"被别处登记过：尺 `grep -rn "width" internal/config/unwired_test.go` ⇒ 只有 `:217` 那行**测试用的 TOML 样本**（`width = 800`），**不在"未接线名册"的断言里**。
   ⚠ 这一条要写手复核：`unwired_test.go` 的语义到底是"登记未接线字段"还是别的；若它本应登记，那本票的第二格是"补登记"而不是"新发现"。

## 为什么值得做（不做会怎样）

最坏后果不是"宽度不对"，是**它假装可配**：用户（这里是 owner 自己）照文档/照配置树去设一个值，界面不动，而他**没有任何一处能读到"这项今天不生效"**——他会归因成"我设错了"或"这产品坏了"。同一族我们已经在提示音、审批窗口、SLO 步上各收过一次。

## AC（每格都要答"这一发在未修码上响不响"）

- [ ] **AC#1 先判归属，再谈修法**：逐枚答三问并落表——① 面板窗口尺寸今天由**谁**决定（现读 `internal/panel/` 里创建窗口/WebView 的那处，给出 file:line；尺：`grep -rn "NewWindow\|SetBounds\|Rect{\|Width:" --include=*.go internal/panel/ | grep -v _test`）；② `cfg.Panel.Width` 到它之间**缺哪一跳**；③ D36 配置树里这一项的**规格出处**是哪一行（尺：`grep -n "\[panel\]" -A 6 docs/PLAN.md` 与 `grep -rn "width" docs/specs/SPEC-03*.md docs/specs/SPEC-08*.md`）——**规格有没有要求它生效**，决定本票是"补实现"还是"改规格文字（＝人工批准）"。
- [ ] **AC#2 正向落地（若 AC#1 判"规格要求生效"）**：把那一跳接上，并给一发**常驻判据**：设成两个不同值 ⇒ 面板侧读到的尺寸**随之变**（不许只断"字段非零"）。⚠ 判据要能区分"接上了"与"又抄了一遍默认值"。
- [ ] **AC#3 反向判据（与 AC#2 成对，缺一形＝装饰）**：**不许**用"把默认值改掉"来交差——把 `default:"640"` 改成别的数而**不接读者**，本票判不通过（那只是把一枚装饰换成另一枚）。若 AC#1 判"规格今天不要求生效"，则正解是**响亮地登记**：在配置项注释＋`docs/reports/pending-and-issues.md` 里写明"此项 RESERVED，设了不生效，缺口在票 N"，**不许静默留着**。
- [ ] **AC#4 顺带查同族**：`internal/config/schema.go` 里**还有几枚**"有 `default` 但生产零读者"的字段（尺：对每一枚字段名跑 `grep -rn "\.<字段名>" --include=*.go internal/ cmd/ | grep -v _test.go`，逐枚记 0/非 0）。⚠ 枚数只许现量、不许目测；本票**不修**它们，只出名单并逐枚归口（新票或既有票），⚠ 不许写"以后加固"。
- [ ] **AC#5 契约轴**：`docs/PLAN.md`（含 D36 那一节与 `:1531`／`:1532` 两行 DEFERRED）、`docs/specs/**`、`thresholds.go`／golden／`allowlist.txt` 一字节不许动；`internal/config/schema.go` 的**默认数值本身今天不许改**（要改数值是 owner 的一句话，见 Q-64 之外的另一枚待答项——若 AC#1 判"规格要求生效"，数值仍保持 640，**本票只接效果、不改数**）。
- [ ] **AC#6 门禁**：逐包 `go test -count=1 ./internal/config/ ./internal/panel/`（⚠ `./cmd/wisp/` 在本机测不到任何东西＝票 98，别把它算进你的红；CLI 那一面只能走 `-overlay` 台件）＋`sh scripts/d22scan.sh`＋名册两向 `comm` 差集；⚠ **最终读数取在最后一枚 commit 之后**。

## 本票**不**解决

- 不裁"面板应该多宽"——那是 owner 的一句话，且**在效果接上之前问他没有意义**（他会答一个不生效的数）。
- 不改任何前端文件（`frontend/**` 归别家，不读不写不引）。
- AC#4 只出名单不修。

## Progress log

- 09-28 10:2x 编排者立票：上面三把尺本程现跑（锚 `039efb47`）。未派。
- 09-28 14:4x `180-c1`（只读普查，与票 182 合派一程；表 `docs/evidence/s1/180-182-panel-fields-census-c1.md` §2–§3）：
  **AC 框一枚没勾、产码零字节未动**。现跑复核（起手 HEAD `e9ef94d0`，非票面锚 `039efb47`）——
  ① 字段仍在 `internal/config/schema.go:527-528`（`default:"640"`）；
  ② 尺 `grep -rn "\.Width" --include=*.go internal/ cmd/ | grep -v _test.go` ⇒ **仍然空输出（0 读者）**；
  `PanelSection` 生产里唯一被碰到的是 `internal/config/manager.go:211` 的**整块指针拷贝**，不读值不生效；
  ③ 三档结论＝**「规格要求生效」**：`PLAN.md:2726` 逐字「`hot`（立即生效）」＋`PLAN.md:2743` 把 `width` 列进 `[panel]` 的 `hot` 行
  ＋`SPEC-03-config-secrets-envs.md:39` 同行；SPEC-08 里 `width` 唯一命中是 `:218` 的 SVG `stroke-width`，与面板窗口无关。
  **缺的不是"一跳"是整条链**：这棵树今天**没有面板窗口**（`internal/panel/doc.go:16`／`pump.go:16`／`cmd/wisp/run.go:438`
  三处注释逐字指票 33/35 未开工），⇒ 票面 AC#1 的"谁决定尺寸"答案是**没有人**。
  AC#4 同族现量：带 `default:` 的行 **69**、去重字段名 **61**、**零生产读者 18 枚**（下界，粗尺同名歧义见表 §3），
  逐枚名与 struct 见表；⚠ **这 18 枚全部已在票 83（done）的 AC#1 全量表里登记并归口**
  （`issues/83-...md:155` 逐字把 `panel.{enabled,width,height,keep_alive_in_session,scale}` 归给 📋 票 33/34/36）。
  现量第 3 条的复核：`unwired_test.go` 的"名册"（`TestEveryLockedSectionKeyIsAccountedFor` `:311`）
  **射程只有 `risk`/`fs`/`net`/`plugins` 四枚 locked section** ⇒ `[panel] width` 不在册是设计如此，
  本票第二格既不是"补登记"也不是"新发现"；**落点本程不选，交编排者**。
