# 324 — 255 那把名册尺按**字段名文本**匹配 ⇒ 别种结构的同名段字段被算成"配置读者"（今天一枚现量假命中）；同时 `llm` 行**确实缺一枚真读者**（dev 上今天的红＝这两件事叠在一起）

- **Status**：已立，**派 `324-r1`**（面＝`cmd/wisp`，起手时该面无写腿）
- **Type**：仪器缺陷（尺的射程写歪）＋名册欠一行；⛔ 不是"把红的去掉"
- **来源**：`A991` 第三节（编排者接手时现抓到，引入者 `272405f5`＝票 320 的取件端口）＋编排者 2026-10-15 21:1x 逐枚读码定性；票 255 已 `-done`，残余按本仓先例**另立一票**（先例＝票 161 结案→残余立 171），⛔ 重开结案票加射程。

## 现象（我现跑，红句逐字）

`go test ./cmd/wisp -run TestTicket255RosterStillMatchesTheActualReadSites` ⇒ rc=1，两条 `t.Errorf`（`cmd/wisp/config_receipt_255_test.go:271`）：
1. `section "llm" (field LLM) is read by [cmd/wisp/panel_config_store.go cmd/wisp/panel_models_320.go cmd/wisp/providers.go cmd/wisp/run.go internal/llm/resolver.go] in production, the roster names […缺 panel_models_320…]`
2. `section "models" (field Models) is read by [cmd/wisp/models.go internal/panel/composer_dispatch.go] in production, the roster names [cmd/wisp/models.go]`

**为什么 CI 没抓到**：`cmd/wisp` 那一步整包**早停**在开窗段（票 314 死亡区），跑不到这枚用例 ⇒ "CI 早停"会把真红**盖成没红**（口径出处 `A976`，今天补一条后果：**名册/计数类尺也在盲区里**）。

## 现量（编排者 21:1x 现跑；⚠ 引用前重跑，行号会腐烂）

1. 尺的判定是**文本级、不解析类型**：`cmd/wisp/config_receipt_255_test.go:306 func scanReadFiles` 走 `cmd internal tools scripts` 四目录，排除 `_test.go`／`internal/config/` 自己／名册件自己（`rosterScanSkipFiles`）／`testdata`／以 `//` 或 `*` 开头的行，然后对**每一行**问 `readsFieldAccess(该行, field)`。
2. ★**假命中一枚**：`internal/panel/composer_dispatch.go:231` = `if d == nil || d.Models == nil {`——这里 `d` 是 `*ComposerDispatch`，`.Models` 是 `ModelsListRequestHandler`（票 320 AC#4 的派发字段），⛔ 配置读者。它只因**字段名与 `models` 那节同名**被算进去。
3. ★**真读者缺一行**：`cmd/wisp/panel_models_320.go:85` = `p, ok := cfg.LLM.Providers[provider]`——这是 `llm` 那节的真读者，名册里没有它（`cmd/wisp/config_readers_255.go:110` 那行只列了四枚）。
4. 名册两行现量：`config_readers_255.go:110 "llm"`／`:125 "models"`；该尺射程＝**14 枚 hot section**（`:259` 那道 `len(fields) < 14` 的门）。

## 要建什么（⛔ 形不写死，先量代价）

- **甲**：把"读者"判定从**文本**升到 **AST＋类型**（命中要求接收者根是配置类型/`Config()` 返回值，而非任意同名结构体字段）。⇒ 语义最干净，代价＝这枚尺自己变成 AST 尺（注意：AST 尺的突变台件**只能用物理改写＋sha256 还原**，`-overlay` 对它失明——`A991` 第一节那条规矩）。
- **乙**：保留文本匹配，但要求同一行同时出现"配置根"证据（`cfg.`／`.Config()`／`LoadFile` 一类），把纯 `d.Models` 这类挡掉。⇒ 改动小，代价＝词面尺的固有脆弱（可被拼法绕过），要在注释里说实话。
- **丙（⛔ 不许选）**：把 `internal/panel/composer_dispatch.go` 直接写进名册 ⇒ **把假命中固化成"名册宣称读者"**，语义变脏；或把该文件加进 `rosterScanSkipFiles` ⇒ **用点名豁免换绿**，等于给这把尺装一个"谁都能加黑名单"的口子。两条都在本票禁区内。

## 判据（AC 框由编排者翻，腿一枚都不许碰）

- [ ] **AC#0 现量复跑**：上面四枚逐枚自己跑（⛔ 抄本票的号），并把"这枚尺今天为什么红"拆成**两笔**（缺真读者／同名假命中）分别具名。
- [ ] **AC#1 正反两形都有牙**（本票核心）：① **正控**——把 `readsFieldAccess` 那一支退化（例：让接收者判定恒真／把配置根证据要求摘掉）⇒ 必须有一枚**具名**用例红（证明新牙不是装饰）；② **反控**——在测试夹具里造一枚**同名段的非配置字段**（`struct{ Models … }` 之类，落仓外快照或既有测试件，⛔ 往产码里塞假字段）⇒ 尺**不许**把它算成读者。两缺一不作凭据。
- [ ] **AC#2 补 `llm` 那行的真读者**：`cmd/wisp/panel_models_320.go` 进名册（判定词、`hotClaim*` 那一档选哪档**要写理由**，⛔ 随手挑一档）；补完复跑该 section 并引绿句。
- [ ] **AC#3 14 枚 section 逐枚复跑**：⛔ 只修这两行就交差——把整把尺跑到底，红名册**与起手作差＝新增 0**（起手那两枚已知红要归位成"具名已修"或"具名移交"，不许默默消失）。
- [ ] **AC#4 名册件不许被写脏**：`config_readers_255.go` 里那些"verdict strings quote code"的句子（该件被尺特意跳过）⛔ 借本次改动扩写；⛔ 往 `rosterScanSkipFiles` 加文件。
- [ ] **AC#5 门禁＋面核**：`go build ./...`／gofumpt（绝对路径）／`go vet ./cmd/wisp/`／`sh scripts/d22scan.sh`／长度门；`go test ./cmd/wisp` **只 `-run` 定向**（⛔ 整包＝票 314 死亡区，会拿假分母）；`git status --porcelain -- internal cmd` 终态＝起手。

## 禁区

`docs/PLAN.md`／`docs/specs/**`／阈值／golden／`tools/d22scan/allowlist.txt`／`.github/**`／`frontend/**`／`design/**`；冻结三件；⛔ 为变绿放宽 255 那批既有钉（它们是**已翻勾的凭据**）；⛔ 顺手改 14 枚 section 的**热档判定**（`hotClaim*` 语义属票 255 的裁定面，动它要在回执里具名）。

**撤销口令**：「撤 324」（⇒ 尺退回文本匹配＋`llm` 行缺读者那一格重新挂回〔待派〕）。

## Progress log (append-only, newest last)

- [2026-10-15 21:1x +08] 编排者立票：来源＝`A991` 第三节；两笔定性（真读者缺一枚／同名假命中一枚）由我逐枚读码取得（`composer_dispatch.go:231` 的 `d` 类型、`panel_models_320.go:85` 的 `cfg.LLM`）。丙形写死为禁区（那是"把缺陷藏进名册"）。
