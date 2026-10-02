# e2e-recensus-1 — H1–H14 端到端跳表复核（只读普查腿）

> 本件唯一问题：台账 `A531`（2026-10-02 13:1x，锚 `cf87f7a2`）那轮端到端跳表 H1–H14，**今天（锚 `f17b1165` 起）哪些跳已经闭上**。
> 口径沿用 `e2e-panel-1`（`.scratch/wisp/probes/e2e-panel-1/readiness.md`）：「通」＝非测试调用者 ≥1 枚（调用点尺）＋码内分支阅读；⛔ 本腿零 Go 命令。
> 复核基准＝A531 当时的三态；今天只对**承重处变化**的那些跳改判，其余照录并写明"照录"。

## §0 起手锚

- 首发锚（同发取）：`date`＝`2026-10-02 18:59:26 +0800`；`git rev-parse HEAD`＝`19b634255788f1e04fd5e9d2aad0ba08971583ee`；分支 `dev`。
- 骨架 commit 时刻复量：`2026-10-02 19:08:57 +0800`，HEAD 已漂到 `f17b1165`（`245-c1` 普查骨架，`--stat` 只动 `.scratch/wisp/probes/245/c1/census.md`，零产码）；`git status --porcelain -- cmd internal`＝**0 行**。
- 自 A531 锚 `21a15d2` 以来产码漂移全量：`git diff --name-only 21a15d2..HEAD -- cmd internal tools scripts`＝**恰好 5 枚**：
  - `cmd/wisp/firstrun.go`（`fa5593c0`，198-r2，回执三段）
  - `cmd/wisp/firstrun_198r2_test.go`（同发，473 行五枚用例）
  - `internal/config/manager.go`（`dd92bb92`，255-r1c，hot 表同源守卫）
  - `internal/config/tiers.go`（同发，TierRegistry 登记表）
  - `internal/config/tiers_255_test.go`（同发，四枚用例）
- 票面状态：198 已改名 `-done`（`f8c45ac8`，五格全勾）；255 AC#0/ⓑ 选形已裁、AC#2-ⓑ＋AC#3 产码已落（`dd92bb92`）但 **AC 框一枚未碰**（翻勾归 255-v1）；257 AC#0 翻勾＋选形＝ⓒ（`5c8fa195`/`ed3fd270`），落地腿 257-r1 未派；248 仍未勾 5 格（AC#2/#4/#8/#9/#10 剩余面）。
- 共享树警示：今天已有六程写腿断流（A541/A542/A544），HEAD 在本腿写作期间继续漂（`245-c1`、`f17b1165`）；本腿全部行号在 `cmd internal` porcelain＝0 的窗口里取，HEAD 漂移若只动 `.scratch/**` 则行号不受影响。

## §1 H1–H14 逐跳三态表

> 三态定义：**已闭**＝A531 判「不通/判不了」的那一环，今天的产码里补上了承重调用点；**半闭**＝补了一半（写明缺的那一环 file:line）；**没动**＝与 A531 时逐字节同形。
> 「已闭」的格子仍受 A531 同一边界约束：调用点尺判「通」**不含运行时证据**。

| 跳 | A531 三态 | 今天三态 | 依据（哪发 commit 哪段码） |
|---|---|---|---|
| H1 | 通 | **已闭（照录＝通）** | 无承重码漂移：`cmd/wisp/main.go:58-67` 无 argv → `runResident()`；今天该链未被动过（diff 名册里没有 main.go／resident_windows.go） |
| H2 | **不通**（建件只挂 `wisp run`，双击那条链第一步因缺文件返错） | **半闭→接近闭但**：`fa5593c0` 只给 `ensureFirstRunConfig` 补了**回执三段**（firstrun.go:103-112：key 走 `wisp secret set` 进 DPAPI／文件里补的是 blob 名不是明文／模型只有改文件一条路），**建件调用点仍是 1 枚**＝`cmd/wisp/run.go:245`，常驻/双击那条链（`panel_inbound.go:230` → `loader.go:64-68`）今天**仍无建件者**；`grep -rn ensureFirstRunConfig cmd internal --include=*.go \| grep -v _test`＝定义 1（firstrun.go:72）＋产码调用点 1（run.go:245）。⇒ **票 198 的四格（AC#2 同源牙／AC#5 失败支／AC#4 去哪儿补／J1 落点）全闭并五格全勾 `-done`（`f8c45ac8`），但那四格判据里从来没有"常驻腿建件"——H2 的缺口（常驻入口建件）不在票 198 射程内，A531 也只写"归票 198"是把优先级提到队首，不是把这一环划给它** ⇒ 本腿今天判：H2 **没动**（常驻入口上仍无人建首份配置），只是"补什么"的指引变好了 | 承重码：`cmd/wisp/run.go:245`（唯一建件调用点）；`cmd/wisp/panel_inbound.go:228-233`（常驻链先读盘）；`internal/config/loader.go:64-68`（缺文件返错）——三处今天与 A531 逐字节同形（diff 名册可证） |
| H3 | 通 | **已闭（照录＝通）** | `NewPanelManager` 定义 `panel_host_windows.go:175`、产码调用点 1（`panel_resident_windows.go:204`）——该两枚文件今天零漂移 |
| H4 | 部分通（热键＋托盘两形通；点球那形不通） | **没动** | `cmd/wisp/resident_ball_windows.go:174` `OnClickBall` 仍只 `recordBallGesture("click")`；`internal/ball/ball_windows.go:50-59` 事件表无双击槽位（`grep -rin "doubleclick\|DBLCLK" cmd internal --include=*.go` 非测试零命中，本腿现跑）。diff 名册里没有任何 `internal/ball/**` 或 `resident_ball_windows.go` |
| H5 | **不通（干净检出必不通）**（embed 目标目录跟踪件 1 枚 .gitkeep） | **没动** | `git ls-files \| grep -c "/dist/"`＝**1**（现量 19:0x，与 A531 同值）；`internal/panel/assets.go:54-60`、`panel_host_windows.go:340-343` 未漂移。⚠ 本机 ignored `dist/` 件 6249 枚照旧（`git ls-files --others --ignored --exclude-standard \| grep -c "/dist/"`＝6249）——"本机那包在不在"仍判不了（A531 §5-2 原样带过） |
| H6 | Go 侧通／页面侧判不了 | **没动（Go 侧照录＝通，页面侧仍禁令判不了）** | `internal/panel/bridge.go:66-67` 两枚名在白名单——bridge.go 今天零漂移；`frontend/**` 禁令不变，本腿不读不引 |
| H7 | 通 | **已闭（照录＝通）** | `composer_dispatch.go:153-165/197-206`、`panel_host_windows.go:319-322` 零漂移 |
| H8 | 通 | **已闭（照录＝通）** | `internal/panel/config_handlers.go:57-69` 名册 7 枚、四道闸门（`:315/:321/:328/:334/:349`）零漂移 |
| H9 | 通（形状是一句人话） | **已闭（照录＝通）** | `panel_config_store.go:88-133 ReadSettings`、`config_handlers.go:292-303` 投递链零漂移 |
| H10 | 通（落盘栈齐） | **已闭（照录＝通）** | `settings.go:199-249 writeOneKey`、`writeguard.go:111/160`、`panel_config_store.go:173-230` 零漂移 |
| H11 | **不通**（干净机器七枚字段全部写不进去） | **半闭**：选形已定（ⓒ＝不建行、回执指路 `[llm.providers.<名>]`，账 `A543`/票 257 §8），且 `fa5593c0` 的回执三段**已经把那条路写出来了**（firstrun.go:103-112 逐字"`[llm.providers.<名>]` 里补 api_key_ref、base_url 与 models.<id>，再在 [llm] 的 text_chain 或 roles.chat 里点名 provider/model"）⇒ "指引"那一环闭了；**但"行可建/可写"那一环没动**：`internal/config/defaults.go` map 仍「leave nil」（`:77-78` 现读同形）、`schema.go` Providers 无 default 标签、写侧 `settings.go:310-314`/`:294-308` 仍要求行已存在、名册（`config_handlers.go:90-96`）仍无 provider 创建字段。落地腿 **257-r1 未派**（A543 §4：按住等 cmd/wisp 写面）。⇒ 判：**半闭**——干净机器上那七枚字段今天仍然一枚都写不进去，只是"指路人"换了张更准的地图 | 依 `fa5593c0` firstrun.go 回执三段＋`ed3fd270` 257-a1 普查（选形 ⓒ）＋写侧三处现读（defaults.go:77-78／settings.go:310-314／config_handlers.go:90-96） |
| H12 | 接线通、运行时判不了 | **没动（照录）** | `panel_host_windows.go:320-321` 回话原样 return、`panel_resident_windows.go:250` Run() 泵——零漂移；"到达"仍需真页面 await，本腿无射程 |
| H13 | **不通**（恒答 restart vs `llm` 在 hot 段表，两套零行码互译；宿主拿不到配置；`[panel]` 五键零读者） | **半闭（A531 的三环里闭了半环）**：`dd92bb92` 落了 **TierRegistry**（`internal/config/tiers.go`：16 枚段级＋app/voice 逐键）＋`manager.go` hot 表改名 `hotSections` 并加**同源守卫**（`:296-303`：循环逐段查 `TierRegistry`，缺失或档位不符当场 panic）＋四枚用例（`tiers_255_test.go`，完备性尺/名册数钉/键级分档钉/voice 走查）⇒ **"两套词表互译"的机器可读底座闭了**。但 H13 判「通」还缺的两环**没动**：① 面板写入侧恒答 restart 仍在——`panel_config_store.go:228`/`:275` 现读仍是 `res.Tier = panel.EffectiveRestart` 硬填，`grep -rn "TierOf(" cmd internal --include=*.go`＝**产码调用点 0 枚**（tiers.go:93 定义＋0 调用；tiers.go:89-91 注释自陈"the panel write path is the intended caller once 255-r2 lands"）；② 宿主仍拿不到配置对象——`NewPanelManager` 签名仍无配置参数（`panel_host_windows.go:175`）、`Width: 420` 仍写死（`:304-305`，本腿现读）、`[panel]` 键仍零读者（`grep -rn "\.Panel\.Width\|cfg\.Panel\."` 非测试＝仅 boundary_test.go:160 一枚测试命中）。⇒ AC#5 那格（回执由登记表同源产出）与 AC#4（真生效）那格的落地腿 255-r2 未派 | 依 `dd92bb92` 三枚文件＋panel_config_store.go:228/275 现读＋TierOf 产码调用点 0 的现量 |
| H14 | **不通（面板侧无泵）** | **没动** | 快照泵非测试构造点仍唯一挂在 `cmd/wisp/run.go:699`；`grep -rn "\.Eval(\|ExecuteScript"` 产码仍零命中（本腿现跑）；diff 名册里没有 pump.go／run.go 泵段 |

## §2 与 A531 结论的差集

A531（via e2e-panel-1）的十四跳里，今天**三态发生变化**的只有两跳，其余十二跳逐字节同形或属"照录＝通/照录＝判不了"：

1. **H11：「不通」→「半闭」**。变化的凭据＝`fa5593c0` 的回执三段（firstrun.go:103-112）把"干净机器缺什么、去哪儿补"从"没说"变成"说且指对了路"（`[llm.providers.<名>]` 这个真路径来自 `257-a1` 普查对票面 `[providers.x]` 的顶回，`5c8fa195` 入账）；但"行不存在则写侧仍拒"那一环原样，七枚字段仍写不进 ⇒ 半闭不是全闭。
2. **H13：「不通」→「半闭」**。变化的凭据＝`dd92bb92` 的 TierRegistry＋manager.go 同源守卫＋四枚用例 ⇒ "两套词表、零行码互译"这一句**在 internal/config 一侧**不再成立（现在是同源、缺行会 panic）；但面板写入侧（`panel_config_store.go:228/275`）与宿主侧（`:175` 签名／`:304-305` 硬编码）三处原样 ⇒ 回执与生效说法仍未接上那张表，`TierOf` 产码调用点 0 枚是硬读数。

其余十二跳（H1–H10、H12、H14）今天的产码承重处与 A531 锚 `21a15d2` 逐字节同形（差集证据＝`git diff --name-only 21a15d2..HEAD -- cmd internal tools scripts` 恰好 5 枚，无一枚落在这些跳的承重件上）。

**一句话**：A531 列的六跳缺口（H2/H11/H5/H6/H4/H13），今天**没有任何一跳全闭**；H11 与 H13 各从「不通」推进到「半闭」，H2 的"去哪儿补"指引变好但常驻建件那一环没动。

## §3 我可能判错的条目

1. **H2 我判「没动」而非「半闭」，尺是"常驻入口上有人建件"那一环**。如果编排者的尺是"机主照回执指引能自己走出未配置困境"（即票 198 AC#4 的"去哪儿补"），那一环 `fa5593c0` 确实闭了、H2 该升半闭。本腿取的是 e2e-panel-1 §1 H2 的原判据（"缺「常驻/双击那条入口上有人建首份配置」这一环"），那一环今天仍无人补。
2. **H11 判「半闭」可能偏宽**。回执三段写在 `wisp run` 首建那条 stderr 路上——机主要**先跑过一次 `wisp run`**（且通常失败退码 2 之后）才看得到这两句；"指引到达"的运行时态本腿零跑判不了。若编排者只认"写侧/名册/默认表"三处码形，H11 该回「没动」。
3. **H13 的「半闭」我把范围钉在票 255 的 AC#2-ⓑ＋AC#3 落地上**，但票面 AC 框一枚未勾（A544 明写翻勾归 255-v1），而 255-r1c 是**编排者代笔＝实现者**。本腿只判"产码承重处在不在"，不判"AC 翻不翻勾"——这两层若被读混，是我的格子写宽了。
4. **行号是工作树单口径**。本腿没跑 `git cat-file blob HEAD:` 复核（时间窗内 `cmd internal` porcelain＝0 ⇒ 工作树＝HEAD，但这是推断不是双尺复认）。凡引用的行号在 porcelain 非 0 窗口里可能漂。
5. **"照录＝通"的十二跳我没有逐跳重跑调用点尺**（只对 H4/H13/H14/H5/H11 的承重处现跑过 grep）。依据是 diff 名册恰好 5 枚、不落那些承重件 ⇒ 逐字节同形。若 18:59 之后还有未 commit 的产码变异（porcelain 非 0 窗口），"照录"失效。
6. **H12/H14 的"判不了"格子沿用 A531 的"为什么判不了"**（零跑／需真机时序），本腿没有重新论证这些格子的判不了性。

## §4 量不到的格子

1. **H5 的"本机那包页面产物在不在"**：需要跑 `panel.BuiltinAssets()` 或读 embed 目录文件名——前者是 Go 命令（禁）、后者落禁令树（`frontend/**` 不读不引）。A531 §5-2 原样带过，本腿同样判不了。
2. **H6 的页面侧**：`frontend/**`／`design/**` 两层禁令，设置入口/控件/`wispDispatch` 是否存在，本腿不读不引不判。
3. **H12 的"回执真到达页面"**：需一发真页面 await（真机时序），本腿零跑。
4. **H14 的"Go→页推送运行时态"**：`.Eval(` 产码 0 命中是静态读数；泵有没有可能经第三方库间接推送，需读第三方模块源码，本腿没读。
5. **H13 的"真窗口宽度随配置变"**（票 255 AC#4 判据）：只有本机可量那一族，需要真跑 GUI，本腿零跑。
6. **`d22scan` 今天红不红**：尺的形状可读（runtests.sh:98-102），今天跳几枚只有跑才知道，本腿没跑。

## §5 我推翻 A531 哪一句

**没推翻。**

A531 的两条死结定性（H2 归票 198／H11 归票 257）、行号更正（`res.Tier` 在 `:228` 非 `:227`）、以及"AC#13 前提收窄成干净检出必跳"三处，本腿全部现量复认成立：

- `res.Tier = panel.EffectiveRestart` 本腿现读仍在 `cmd/wisp/panel_config_store.go:228`（且 `:275` 同形）——与 A531 更正一致。
- 干净检出的 embed 口径：`git ls-files | grep -c "/dist/"`＝1，与 A531/票 248 那句"只跟踪一枚 .gitkeep"同值。
- H2 的缺口定性（"缺的是常驻入口建件"）与今天的 diff 名册一致——票 198 结案五格全勾，但它的四格判据（AC#2/AC#4/AC#5/J1）里没有一格覆盖常驻腿建件 ⇒ "归票 198"在 A531 原文里是**优先级提到队首**，不是射程归属，本腿这一读与 A531 正文不冲突。

唯一一处**措辞层面**要补的（不是推翻）：A531 说 H11"归票 257"——今天 `257-a1` 已把选形裁成 ⓒ 并顶回了票面的 `[providers.x]` 拼法（真路径 `[llm.providers.<名>]`），`fa5593c0` 的回执已经按真路径写。这枚进展 A531 落笔时（13:1x）尚未发生，属**新增事实**，不是 A531 的错。

---

### 交件态

- 写点一枚：`.scratch/wisp/probes/e2e-recensus-1/recensus.md`。commit 带显式 pathspec、只含这一枚文件。⛔ 未 push。
- 零跑复述：`go build`／`go vet`／`go test`／`./...`／`go list` 一次都没执行。工具面＝`grep`／`sed`／`ls`／`wc`／`git log/show/diff/status/ls-files/rev-parse`。
- 零污染自证：起手与骨架 commit 两发 `git status --porcelain -- cmd internal` 均 0 行（18:59:26／19:08:57）。`docs/**` 本腿只读台账与票面、一字未改。
- 冻结件（`docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／三枚冻结件／`.github/workflows/ci.yml`）零字节触碰。
- 票面勾选框一枚未碰；凭据值零出现（全文只出现字段名/命令名/引用形状前缀名）。
- 占位词自证：本件不含任何待填类占位词（那三个英文缩写与那枚括号词，全文零命中——本行是唯一提及处，且以"零命中"为内容）。
