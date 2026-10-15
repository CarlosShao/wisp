
## A991｜2026-10-15 20:5x +08（`date` 现读）｜收件 `321-r1`＋`120-r2`（两枚都我自己复跑过）＋★抓到一枚**dev 上的存量红**（上一班没发现）＋立票 322＋**我自己两发台面故障**如实登记

### 一、`321-r1` 收件（票 321，托盘/球/热键那一跳的执行器判据）
- **三把尺**：① `git log -1 58de465a` 在（`321-r1 落地腿交件（票 321，零产码）…`）；② 件真在：`cmd/wisp/exec_holder_toggle_321_test.go`（200 行）／`docs/evidence/s1/321-toggle-hop.md`（102 行）／`probes/321/r1/**`（13 读数＋1 备份件），`git show --stat`＝**18 件、+1207、删除 0**，逐枚都在派单枚举的四类路径里；③ **我自己复跑**（`go test ./cmd/wisp -run 'TestAC321PanelHopDrivesTheHolderToggle|TestAC321NamedExecutorIsTheToggleRoad' -count=1 -v`，GOTMPDIR＝`C:\Users\swq\AppData\Local\wisp-test-gotmp\orchwave1`）＝**两枚全 `--- PASS`**，包级 `FAIL` 只因同批那枚 255 尺（见第三节）。
- **产码零动**：`resident_windows.go` 不在该笔改动清单里（腿自报 sha256 基线 `816017b6…` 已复现两次，备份件 `probes/321/r1/resident_windows.go.bak` 入库留底）。
- **判据有牙（腿交三发，我复跑正控那一发）**：突变甲（`:269` 同行 `holder.toggle`⇒`holder.show`）两枚尺全红、红句含 `this hop is NOT toggle semantics`；突变乙（反形：闭包体 `return false`）同一两枚全红 ⇒ "判据换成反形也红"这一条**它做到了**，我没盖章——**仍要另一枚非实现者**裁"没顺手放宽别的断言"（票 321 的 AC 框一枚未翻）。
- ★**方法学新条款（影响本仓所有 AST/源形状钉的变异台件，腿交回、我复按其逻辑成立）**：`go test -overlay` **只换编译期输入**；凡"运行时用 `parser.ParseFile`/`os.ReadDir` 读盘上文件"的尺，用 overlay 去变异被测量那一边 ⇒ **尺恒绿＝假读数**。⇒ 定式：**先问"被测量那一边是不是盘上文件"**，是 ⇒ 突变只能**物理改写＋sha256 还原＋`.bak` 留底**；319-v2 那类 overlay 成立是因为它驱动的是**编译后的函数**。**今后每枚带 AST 钉的派单都写这一条**（本条同时进我的派单模板与记忆）。
- **腿自抓一枚违纪（如实入账，归腿）**：为记突变甲的 sha/numstat，`cp` 出的临时件 `probes/321/r1/logs/mutant-A-subject.txt` 被它 `rm` 删了 ⇒ 违反"临时件只建不删"（`issues/README` 规则 8）。内容零丢失（同字节仍在 `.bak`、基线 sha 复现两次）。规矩重申写进后续派单。

### 二、`120-r2` 收件（票 120，fd 锚定那一族）
- **三把尺**：① `git log -1 6a6b1b18` 在；② 面核：非 `probes` 改动＝`internal/winsec/winsec_other.go`（+239）／`posix_seal_window_120_other_test.go`（544 改）／`docs/evidence/s1/120-fd-anchor.md`（176 行）／票 120 ＋39 行；**那 39 行是不是碰了框＝行为尺判定**：`git show … | grep -cE '^[-+] *- \['` ＝ **0** ⇒ 框一字未动（那条 bullet 本身 39 行长，属"一格一条"形状，不违规）。③ **我自己进容器复跑关键读数**（`git archive HEAD` 导出到 `/d/tmp/orch120chk/base`，先数 `.go`＝**1072 枚证树非空**，`golang:1.27`、`CGO_ENABLED=0`、`WISP_120_ATTEMPTS=300`、`-run TestPOSIX120`）：三形 `hit-foreign=0` **且 liveness 非零**（`ancestor-swap-absolute-link` total `attempts=1200 refused-at-check=67 of-which-refused-at-act=6 sealed-own-object=70`；`…relative-sibling-link` total `refused-at-check=69 … sealed-own-object=108`），包级 `ok`、`inside_rc=0`。⇒ **"0 命中"旁边站着分母与非零活性**，不是死仪器；我这发是**缩小档（每形 4 样本×300）**，与腿的 2000 发**不可相减**，只作独立复认。
- **腿具名交回的缺口（不当已完成）**：`sealHandle` 自己那半 `SameFile` **零具名单元台件**；形 1/2 去牙后仍 0（`os.Root` 本身拒绝对目标链接，红句已为这条写）；"每形每发 `refused-at-act>0`"今天只 `t.Logf` 不 `t.Errorf`。⇒ 这三条是 `120-v1`（非实现者）要裁的对象，⛔ 我在此追认。
- 票面那句 `154/4000` 仍是**改前口径**（措辞属改票面＝须人工批准）⇒ 本班次**没动它**，登记在案。

### 三、★dev 上的一枚**存量红**（我接手才抓到；⛔ 追认给腿、⛔ 说成"CI 会抓到"）
- **现象（我现跑）**：`go test ./cmd/wisp -run TestTicket255RosterStillMatchesTheActualReadSites` **红**。红句逐字两条：`section "llm" (field LLM) is read by [cmd/wisp/panel_config_store.go cmd/wisp/panel_models_320.go cmd/wisp/providers.go cmd/wisp/run.go internal/llm/resolver.go] in production, the roster names […没有 panel_models_320…]`；`section "models" (field Models) is read by [cmd/wisp/models.go internal/panel/composer_dispatch.go] in production, the roster names [cmd/wisp/models.go]`。
- **归因（不是腿说的算，我自己读的两枚新增读者）**＝`272405f5`（10-15 16:27 票 **320 的取件端口**，`cmd/wisp/panel_models_320.go` 与 `internal/panel/composer_dispatch.go` 各新增一处 section 读取）**未回写** `cmd/wisp/config_readers_255.go` 的 `sectionReadSites`。
- **为什么没人发现**：CI 的 `cmd/wisp` 那一步**早停**在开窗段（票 314 那个死亡区），**跑不到这枚用例** ⇒ 这一族"CI 早停"会把**真红盖成没红**（口径出处 `A976` 只说了"四数不作分母"，今天多一条后果：**名册类尺也在盲区里**）。
- **处置**：归口＝票 **255**（那把尺的主人）。落地腿 `255-r2`（回写两枚读者位＋复跑名册＋门禁）已排入队列，**按在 `320-v1` 之后**（它要跑 `cmd/wisp`，与在飞腿撞面）。⛔ 不许"为了让它绿"去放宽那把尺的判据；⛔ 也不许把红算进 `321-r1` 的分母（起手＝终态、红名集合 diff 为空，两枚读数都在 `probes/321/r1/logs/`）。
- **记我**：上一班把 320 的 AC 框翻了六格、并自报"cmd/wisp 段 ok"，但那一发没带上这枚名册尺 ⇒ 我接手时**没有**现跑它。教训＝**收件时把"该包所有名册/计数尺"整族跑一遍**，别只跑与本票同前缀的那几枚。

### 四、立票 322（`A990` 裁①当场兑现）＋派只读普查
- 票面 `.scratch/wisp/issues/322-privatefile-open-truncates-through-a-swapped-leaf.md`（六格 AC；`AC#0` 未出三形代价表之前⛔ 任何一形开工；`AC#3` 把 **Windows 那一半单独答**——本仓目标平台＝Windows，⛔ 用 POSIX 绿冒充）。
- 现量（我 20:5x 现跑）：`winsec.go:108` `privateFile`／`:113` `os.OpenFile(resolved, flags, perm)`／调用者两枚 `:87`（`O_CREATE|O_TRUNC`＝**敞着的这一支**）与 `:97`（`O_CREATE|O_EXCL`）；`winsec_other.go:82` 注释已逐字指向本票。
- ⚠ 票面**没有假设** `O_EXCL` 那一支"天然安全"——POSIX 该规定与直觉相反的先例本仓刚有过（`120-r1` 实测 `O_NOFOLLOW` 经 `os.Root` 失效）⇒ 要实测。

### 五、★我自己两发**台面故障**（长得像读数，都不是读数；如实登记⛔ 抹进结论）
1. `GOTMPDIR` 指到**没建过的目录** ⇒ `go: creating work dir: … The system cannot find the file specified.` ＋ `rc=0`。定式＝**先 `mkdir -p` 再跑**；红句里出现 `creating work dir` ⇒ 判台面故障，⛔ 记进用例。
2. `git archive` 的导出树放到了 **Git Bash 的 `/tmp`（＝`C:\Users\swq\AppData\Local\Temp`）** ⇒ Docker Desktop 挂到的是**空目录**，容器里回 `go.mod file not found`。定式＝**容器台件的树放 `D:` 卷**（腿用的 `/d/tmp/probe120/base` 那形），且 `-v` 要配 `MSYS_NO_PATHCONV=1`；量之前先数 `.go` 证非空（我这发＝1072 枚）。
3. 附带一枚：**⛔ 在 `cmd | tail` 之后取 `${PIPESTATUS[0]}` 当被测命令的退码**（我这次 `host_rc` 就是这么来的，不可信）⇒ 定式＝`cmd > /d/tmp/x.log 2>&1; rc=$?`，退码与被测命令**紧邻**。

### 六、编队现状（本条落账时刻）
- **在飞 3 枚**：`270-r1`（写 `internal/tools`）／`320-v1`（验收，读 `cmd/wisp`＋`internal/config/panel`）／`322-a1`（只读普查＋容器探针，⛔ 写面外零改）。
- **按住的与具名理由**：`321-v1`（它要**物理改写** `resident_windows.go` 做突变 ⇒ 撞 `320-v1` 正在读的 `cmd/wisp` 面，会污染它的读数）；`255-r2`（同样要跑 `cmd/wisp`）；`165` 装配半格（`cmd/wisp` 面）；`287`／`171`（`internal/tools` 面，`270-r1` 在飞）；`86`（条件＝无 docker／无并发 go test，今天两条都不满足）。
- 未推笔数：`78ea61e2`＋`58de465a`＋`6a6b1b18`＋本条＝**推送前门禁还没跑**（长度门/d22scan 由腿报，但**我推之前自己再跑一遍**）。
