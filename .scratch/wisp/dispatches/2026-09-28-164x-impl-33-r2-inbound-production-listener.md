# 派单 33-r2（写腿）— 片 B：给入向那一跳接上**一个具名的生产调用者**（此前一程死在我自己派单的矛盾上）

- 时间：2026-09-28 16:4x +08　编排者锚（**不等于你的锚**）：`ba65b29f`
- 工单：`.scratch/wisp/issues/33-panel-host-c27.md`（票 33）· 上游设计核：`docs/evidence/s1/33-inbound-hop-design-a1.md`（`33-a1`，台账 `A374`）· 片 A 交付：`internal/panel/composer_dispatch.go`（199 行）＋`composer_dispatch_test.go`（481 行／12 用例），台账 `A380`
- ⚠⚠ **这程的来路是一枚我欠的错**：上一程（片 A）我一面要求"必须有真生产调用者、不许拿测试假宿主充数"，一面在禁区里写死"不许动 `cmd/wisp/**`"——**两条同时成立＝无解**，那程诚实停手上报了。⇒ **本轮放开 `cmd/wisp/**`**，矛盾作废。账记在 `A384` 裁定③。

## 你要落的那一格（一句话）

`ComposerDispatch.Handle`（`internal/panel/composer_dispatch.go:120`）**今天生产调用者零枚**（我 16:4x 复跑：`grep -rn "ComposerDispatch" internal/ cmd/ tools/ | grep -v _test.go` 只命中它自己的定义与注释）。片 B＝**从"用户动作那一端"起一条真能读的通道，把原始字符串送进 `Handle`**，并让这条通道在**生产装配路径**上有一个具名听众。

## 三种形状（自己选，选完在表里具名说为什么）

- **甲（推荐起点）**：在常驻/GUI 那条回路里挂一枚具名入向组件（如 `ComposerInbound`），它的 raw 源＝**今天已存在的、本机可端到端跑的通道**（命名管道／stdin／已有事件循环里的一格），宿主未接线那一段**显式命名成没接**，不许用注释冒充。
- **乙**：给 `wisp` CLI 加一枚**具名子命令或 flag**当入向缝（走 `C17 PanelBridge` 之外的那条"CLI `wisp run` 注入面"，`AGENTS.md` §1.3 明文许可），并在 `cmd/wisp` 装配根里真调用。⚠ 若选这一支，**必须同时说明它离"界面点一下"还差哪几跳**（`33-a1` 那张 H1–H10 表里 H2／H3／H10 要真宿主）。
- **丙**：**落不了地就报回来**——具名写出"缺的是哪一层、需要放开哪一枚文件/依赖"，不许为了变绿自创形状。

## 硬约束（一条破即判失败）

1. **不许引入任何新依赖**：`go.mod`／`go.sum` 一字节不动（**不引 WebView2**，那属票 33 余下的 H2/H3/H10，另程处理）。
2. **不许拿"只有测试能构造的假宿主"当生产听众**（`A205`／`[[verification-blind-spots]]` 第 29ⓑ 那一族：注释里声称的钉、尺本身是空的）。判据：至少一枚**非测试**文件里的调用方，且端到端能在本机跑一次并留读数。
3. **写面**：`cmd/wisp/**`（**除 `cmd/wisp/run.go`**，那枚文件在别人的串行队列上）＋ `internal/panel/composer_dispatch.go` ＋ 你新建的文件（`internal/panel/` 或 `cmd/wisp/` 下，命名自取并具名登记）。
4. `internal/panel/bridge.go:42-45` 那**四枚白名单常量的顺序与命名一字不动**（`composer_test.go:502` 有钉）；**不新增任何 `panel.*` 或 `config.*` 常量**——名册补齐属票 194 堆1，排在你的活之后。
5. **禁区不动**：`frontend/**`／`design/**` **零写面**（读可以，写不行，也不许计入任何"零命中"宣称）；`internal/panel/tokens_fourway_test.go`、`l2_grant_boundary_test.go`、`frontend_hygiene_test.go` 一字不修；`allowlist.txt`／`thresholds.go`／golden／C18 审批超时常量／`docs/PLAN.md`／`docs/specs/**` 一字节不动；票面 AC 里点名的既有测试用例不许改（要新判据就**新建文件**）。
6. **不要碰在飞那程的文件**：`internal/risk/provenance.go`、`internal/risk/taintmatch.go`、`internal/risk/hostpath_185.go`、`internal/risk/pointer_185_test.go`、`internal/tools/pointer_185_cli_seam_test.go`（`185-r1` 的活）。它们出现在 `git status` 里是**别人的**。
7. **别把"面板可以撤销授权"这类语义顺手实现**（那是票 194 堆1 与票 195 的射程）。

## Git 纪律（逐字遵守）

只 commit、**绝不 push**；commit 必带**显式 pathspec**（`git commit -q -F - -- <你的路径>`）；禁 `git add -A`／`git add .`／`-a`；禁 `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`／`switch`／`worktree`；**仓内零删除命令**（临时件只建不删；要恢复跟踪件用 `git cat-file blob HEAD:<path> > <path>`）；**绝不在仓库目录内建 worktree 或 checkout**。每完成一格 commit 一次（这把损失从"整票重跑"压到"半张表"）。

## 起手必做（三件，顺序不能反）

1. `date "+%Y-%m-%d %H:%M %z"` ＋ `git log --oneline -1` ＋ `git status --porcelain`：**把起手名册逐枚抄进证据件第一行**（含上面那批 185 的脏件）。
2. 现量并写进表：`grep -cE '=\s*"panel\.' internal/panel/bridge.go`（我 16:4x 现量＝**4**；⚠ 别用 `grep -c 'panel\.'`，那把尺会把 `bridge.go:35` 的**注释**也数进去、得出 5——上一程我就是这样写错派单的）、`wc -l internal/panel/composer_dispatch.go`、`grep -rn "ComposerDispatch" internal/ cmd/ tools/ | grep -v _test.go` 的**非 test 命中枚数**（起手应＝**0**，这就是你要改的那枚数）。
3. 读 `33-a1` 那张 H1–H10 表（`docs/evidence/s1/33-inbound-hop-design-a1.md`）再决定甲／乙／丙。**前提与之不符就报回来，别硬改。**

## 你要交的判据与自证

- **新常驻判据**（进新建的 `*_test.go`）至少三条：① 从你那条通道投一枚具名形状的请求 ⇒ `Handle` 真被走到（断言可观测后果：审计行／回执串），② 白名单外的方法名 ⇒ 拒答**且写审计**（复用片 A 的 `rosterMismatch` 那支，别新造），③ **反例**：把生产调用者那一跳摘掉（改名或注释掉真不可，就用变异构造）⇒ **判据必须红**（"函数存在但零调用方"这一形要被看得见）。
- **变异自证**：至少两发，逐发具名"改了哪一行 ⇒ 哪条用例红"，跑完 `internal/`＋`cmd/` 树必须回到干净态。⚠ 别用 `go test -overlay` 造载体：**扫磁盘真实文件的尺看不见 overlay**（`A382` 的教训）；要副本就复制到**仓外**，且**仓内零删除**（在仓外临时目录里 `rm` 也请避免，用带序号的新目录）。
- **门禁四数（口径按 `A384` 的新读法，别再自造矛盾）**：`sh scripts/d22scan.sh`（只看 **clean／红＋红点逐名**；⚠ 它输出的 `ban #8 internal/ = N` 那格是**扫到的文件枚数**、不是违规枚数，**在飞的 185 会让它涨，涨不算漂移**）；`bash .scratch/wisp/probes/154/gate-clauses.sh`（只比 **BAD 腿名册**，今天在册＝**只 `G6neg`**；⚠ 别看退码；**绝不要跑 `probes/161/r6/flip-declaration.sh`**）；`go test -count=1 ./internal/panel/ ./cmd/wisp/`（⚠ `internal/panel/` 今天有**三枚在册红**：`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`——**逐名照实记、不修、不当绿**；`./internal/risk/` **别单跑也别顺带跑**，那是 185 的写面）；`"$(go env GOPATH)/bin/gofumpt" -l` 空。**若要跑 CLI 测试**：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（否则 `0xc0000135`，票 98）。
- **终态闸门**：`git status --porcelain` **必须等于起手名册（逐枚具名差集为空）**——不是"必须为空"（共享树里做不到，`A374` 已入账）。新增未跟踪文件**只建不删**并具名。

## 落盘位置与报告格式

- 证据件：`docs/evidence/s1/33-inbound-listener-r2.md`（骨架先建、每节 commit 一次；⚠ 半份表比没有更危险，写不出结论就写"未裁完"）。
- 票面：给票 33 **追加一格** `AC#9（片 B：入向有具名生产听众）` 并把判据、尺、锚点写进去；**别勾**（勾要非实现者裁）。若你判"票 33 现有某格才是归宿"，在票面 Progress log 具名说理由，别改既有 AC 原文。
- 回报（≤400 字）：① 你选甲／乙／丙与一句理由；② 起手与终态那三枚现量数（非 test 调用方枚数：0→?）；③ 变异逐发"哪行→哪条红"；④ 门禁四数原样读数；⑤ 你被拒过的每次工具调用（一栏，别省）；⑥ 没测到什么＋"这一程之后离'界面点一下后端真收到'还差哪几跳"。
- 工具调用预算：硬顶 **35**，第 5 枚前给骨架、第 15 枚前给第一段结论。
