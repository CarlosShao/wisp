# 派单：`33-a1`（只读设计核·零产码）——"网页点一下 → Go 收到"那一跳要建什么：地基普查与代价表

- 派单时刻：`2026-09-28 14:5x`｜编排者锚点 HEAD＝**`1f22f1aa`**（起手若已漂，按实际 HEAD 做并登记）
- 工单：`.scratch/wisp/issues/33-panel-host-c27.md`（面板宿主层）
- 性质：**只读设计核** ⇒ **AC 框一枚不许勾**、**产码与 `go.mod` 零字节不动**、**不落地**。

## 0. 为什么这一单现在派（一句话）
`181-c1` 刚坐实（编排者复跑同值）：**面板→Go 的入向那一跳今天整层不存在**——`ParseComposerRequest`（`internal/panel/bridge.go:84`）与 `HandleModeRequest`（`composer_handlers.go:111`）**非 test 零生产调用方**，全仓没有一枚 WebView2 消息接收器，`go.mod` 里 webview 依赖 **0 命中**，`cmd/wisp/run.go:229` 自己写着 "Nothing calls it yet"。⇒ **owner 要的那几枚"能点的按钮"（票 186 切工作树／票 187 改模型档位）全都排在这一层后面**。本单不是"该不该做"，是**把它要动什么一次摊清**。

## 1. 写面（只这些）
✅ 新建 `docs/evidence/s1/33-inbound-hop-design-a1.md`｜✅ 新建 `.scratch/wisp/probes/33/a1/**`｜✅ 追加票 33 一段 Progress log（只追加、不改原句、不勾框）。
⛔ `internal/**`、`cmd/**`、`go.mod`／`go.sum`、`docs/PLAN.md`、`docs/specs/**`、`allowlist.txt`、`thresholds.go`／golden **一字不动**。起手与终态 `git status --porcelain -- internal/ cmd/ go.mod` 都必须为空。
⛔ `frontend/**`／`design/**` 不读不写不引不转述。⛔ 不属于你的脏件（`.gitignore`、`probes/152/**`、`probes/161/r6/logs/flip-*`、`design/**` 那批未提交删除、`docs/evidence/s1/152-*.md`）**不动不提交不评论**。
⛔ 零删除命令；不许跑 `probes/161/r6/flip-declaration.sh`；只 commit 不 push、显式 pathspec、禁 `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`／`worktree`。
⚠ 重跑任何既有台件前先 `grep -n "WriteFile\|OpenFile" <那枚台件>` 看写出路径（`A367` 那枚坑）。

## 2. 必答六问（每条现跑，不许推理；引规格要**逐字引原句**不许转述）

**Q1 规格到底要求什么**：`docs/specs/SPEC-08*`（球与面板）与 `docs/PLAN.md` 的 D29／C27／C17 相关段落里，"网页事件进 Go"这一跳的**规格原文**是什么？哪一节把它定成宿主层的责任？⚠ 顺带答一枚在册待定项：`AGENTS §2` 列着 "`C24 GojaHostAPI` 初始集与 `C17` 方法白名单定稿（S7/S5 切片卡批准）"——**本单只报它今天定到哪一步，不替它定稿**。
**Q2 依赖面**：接 WebView2 要引什么（今天 `go.mod` 里 webview **0 命中**，尺 `grep -c webview go.mod`）？仓里有没有已写好的 Win32/COM 地基可复用（`internal/ball/sta_windows.go`、`internal/proc/**`、`internal/winsec/**` 各自今天提供什么，逐枚 `grep -n "^func "` 现量）？⇒ 结论要落到"**新写几枚文件／复用哪几枚**"。
**Q3 最小闭环形状**：从"网页 `postMessage`"到"`HandleModeRequest` 真被调用"，逐环列需要哪几枚新函数/注册点（每环指到现有最近邻：`bridge.go:84` 的解析、`composer_handlers.go:111` 的处理、`pump.go` 的出向管子）。⚠ **出向管子今天真在跑**（`cmd/wisp/panel_pump.go`）⇒ 只补入向那一半，别重画出向。
**Q4 冷拉起代价**：`AGENTS §2` 那条"WebView2 冷拉起 >2s → 重评 L2 卡是否回原生"的待定项，本单能给什么现量（有没有可测的代理指标）？没有就明说"要落地腿才能测"。
**Q5 与 owner 那几枚按钮的关系表**：逐枚列——票 186（切工作树／分支）／票 187（改模型、改思考档位）／票 92＋R20（面板里的档位与工作区是"权限输入口"）／票 114（原生侧门已立但入向没接）——**每一枚答"它等这一跳的哪一部分"**，并标"这一跳建完之后它能不能立刻动"。
**Q6 拆不拆得开**：有没有一条"**先只建最小入向、不引 WebView2**"的路（例：CLI/测试用的假宿主先把 router 与写腿跑通，真宿主后补）？如果有，给出落点与判据形状；如果没有，具名说"必须真宿主"。⚠ 这一问决定票 186/187 是"等地基"还是"能先动一半"。

## 3. 门禁（只读程也取数，按 `A363` 不充当结案凭据）
`sh scripts/d22scan.sh`（基线 rc=0、`ban #8 internal/` **examined=433**）；`bash .scratch/wisp/probes/154/gate-clauses.sh` **比红腿名册**（在册只 `G6neg`）；`go test -count=1 ./internal/panel/ ./internal/config/`。⚠ **`./internal/panel/` 今天有 2 枚已知红**（`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`＋`TestC21DesignTokensFourWayAgree`，同一因＝别家会话把 `design/assets/tokens.css` 删了未 staged）⇒ **照实记、不许当绿、不许去修、不许碰那两枚文件**（`tokens_fourway_test.go` 在册冻结）。⚠ **不要改 `gate-clauses.sh` 本身**。终态读数取在你自己最后一枚 commit 之后。

## 4. 表骨架（第 ≤5 枚调用内落盘并 commit 第一枚）
① 起手锚＋写面闸门 ② Q1 规格原文 ③ Q2 依赖与可复用地基 ④ Q3 最小闭环逐环 ⑤ Q4 冷拉起 ⑥ Q5 与四枚按钮的依赖表 ⑦ Q6 能不能不引依赖先建一半 ⑧ **推荐拆法与代价**（不拍板） ⑨ 本程没测什么（逐名） ⑩ 门禁终态 ⑪ 被拒调用＋零删除自证＋工具调用终值 ⑫ next＝派写腿之前还缺什么、哪几枚要人先批准

## 5. 硬顶与增量交付
工具调用**硬顶 35 枚**，**到第 25 枚停止新探索**；每答完两问 commit 一枚；表没落盘前不许继续取数；判"必须改产码才能答完"＝停手上报。

## 6. 结束消息回我八项
① Q1 规格原句（逐字）；② Q2 要引什么依赖＋可复用清单；③ Q3 逐环表；④ Q6 的结论（能不能不引 WebView2 先建最小入向）；⑤ Q5 那四枚依赖关系；⑥ 门禁四数＋红腿名册＋`internal/panel` 那 2 枚已知红照实记；⑦ 三条 `git status` 为空＋工具调用终值＋被拒调用＋有没有跑删除命令；⑧ next＝写腿还缺什么、哪几枚要人先批准。
