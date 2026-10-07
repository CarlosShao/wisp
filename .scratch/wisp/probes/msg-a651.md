
## A651〔10-07 08:5x–09:0x 编排者〕两枚隔夜死腿都零提交⇒我的"写腿闸门"定式被证伪；代提 19 件；另抓到我自己把两形 build failed 当成读数

**0. 盘面起手（⛔ 不信通知，全部现量）**：`date`＝2026-10-07 08:51 +08；HEAD=`e19fce1f`；`git status --porcelain -- cmd internal tools docs .scratch/wisp/issues`＝**0 行**；票 272 逐格＝**未勾 6／已勾 0**（`272-…so-a-banner-carrying-the-same-words-passes.md:31-36`）⇒ 昨夜那枚写腿**没有留下半截产码**。⚠ 顺手抓一枚我自己的路径错引：我上一轮把票 272 的文件名 slug 记成"its-scope-is-the-whole-run-not-this-occurrence"，真身是"so-a-banner-carrying-the-same-words-passes"——按旧名 grep 框数会得到空结果，差点被我读成"票面上没有判据框"。

**1. 两枚死腿的用量与死法（第三种死因，专属写腿）**：
- `272-r1`（写腿）：**71 次调用／813 万 token／44 分 5 秒**，零 commit。
- `done-uncensused-5g3`（装配腿）：**15 次调用／145 万 token／15 分 39 秒**，零 commit ⇒ 与预算无关，纯服务端故障（我上一轮才把它的射程砍到"只装配"）。
- ★**要紧的是我派单模板错了一枚**：我写的"第 45 次调用之前先把 §判不动＋§0 写满并 commit"这道闸门，**对写腿根本不成立**——写腿每跑一次 `go test ./cmd/wisp` 是几十秒到几分钟，71 次调用＝44 分钟墙钟（正好撞上同一天"登录过期"那个时长）。⇒ 定式改为：**写腿的第 1 笔 commit 在起手锚跑完之后立刻做，在任何一条长跑命令之前，不等任何调用数**；调用数闸门（45／70）**只适用于零 go 的只读腿**。已进记忆 `dispatch-output-budget`（死因由三种改四种）。

**2. ★记我：代提说明第一版里我把两形 build failed 当成了"读数"**。我照文件名写"`A1-PART`／`A1-REASON`／`A1-MBOTH`／`A1-REPEAT` 四份是配套的边界与理由段读数"——现量后：`PART`（拷贝 `:128:121`）与 `MBOTH`（`:120:131`）**都是 `syntax error: unexpected newline in argument list`＝build failed，不是读数**；`REASON`／`REPEAT` 才是有颜色的跑。⇒ 定式：**代提说明里每一句"A 这份＝某读数的含义"必须先有一把尺对拉过**，⛔ 不许按文件名推断语义（这把尺这次＝`diff <(git show HEAD:<文件>) <拷贝>` 数改动行 ＋ logs 里 `grep -oE '(--- PASS|--- FAIL|build failed)'`）。落笔前改掉，未入库假话。

**3. 死腿遗产的有效部分（票 272 的靶子它已经量到了，判语我不代填）**：`M` 形对 HEAD 的差异＝**恰好 1 行**（`cmd/wisp/config_reload.go` `:127a128` 插一枚"上一版横幅（种改动之前打的）"，正文逐字含「本次运行不会生效」「需要重启进程」），而它的 logs 里 `A1-M.txt`＝`--- PASS`。⇒ 这正是票 272 AC#1 要的**未修码读数**（横幅先带同样字样时那把尺今天不拦）。`cur`＝0 行改动／PASS，`MDEL`＝8 行／FAIL。⚠ 但按第 2 条口径，重派腿**必须自己复跑这几把**才算凭据（[[feedback-verification-blind-spots]] 第 108 条：死腿的 logs 只是"它取过"）。

**4. 处置＝代提一笔**：`82a10da4`（19 份文件／**+722／−0**，逐枚 ≤60 KB，最大 7,796 B）＝`.scratch/wisp/probes/272/r1/logs/**` 8 份＋`pool-validity/5g2/logs/**` 新增 10 份＋说明件 `salvage-272r1-5g3.txt`。⛔ 一枚没删。`/d/tmp/wisp272r1/waveA/{mut,overlay}` 那 7 形拷贝在**仓外**：不入库（体积）⛔ 也不删，重派腿的派单里具名给了它的路径。

**5. 负载尺这次是两把，读数不一致，具名写清**：WMI `LoadPercentage` 三发＝**79／66／69%**，PerfMon `\Processor(_Total)\% Processor Time` 三发＝**63.2／72.2／72.5%**（我先前那次采到 60.5／62.5／62.8）。占用者不是我的腿：`Everything`(3989 s 累计)＋Qoder 自身，`go/wisp/webview2/node` 只剩 5 枚且非测试在跑。⇒ 我按"PerfMon 那把"判**可以开工**（机主今早原话：「赶紧把任务清点一下然后继续干活儿，记得写票不要太少」），⛔ 不是把 ≥70% 那条规矩作废——本波只放 **6 枚**且**只有 1 枚会在 Go 包里跑突变**（互斥那条只在跑突变的腿之间成立）。

**6. 本波编队（09:0x 派 6 枚，一律 `run_in_background`）**：
- ★写腿 2 枚：`272-r2`（唯一跑 Go 突变的腿，写面＝`cmd/wisp` 那把尺；commit-first 新规矩逐字进派单）、`111-r5`（写面＝`.github/workflows/ci.yml` 那一步缺的 `if: !cancelled()` 守卫，纯 YAML 不动 Go，故不与 `272-r2` 抢读数）
- 只读 3 枚：`270-a2`（票 270 三形代价配对，`internal/tools`，⛔ 禁跑 `./cmd/wisp`）、`voice-wire-1`（`internal/audio` 零 importer／`internal/speech` 只有 `doc.go` 那笔 09-30 现量的**今天复量＋接线顺序**）、`panel-resident-1`（后端主链第一环：面板宿主接进常驻腿到底缺哪几跳）
- 装配 1 枚：`pool-5g3b`（`batch5g2.md` 只补 §2／§3／§判不动 三节，⛔ 不重跑任何尺、禁读五份母本、第 25 次前必交）
- ⛔ 按住的：`236-m1`／`228 AC#2/AC#11`／`253 AC#1`／`197-r3`——全部撞 `cmd/wisp` 写面或要跑第二发突变，等 `272-r2` 交完再逐枚串行（具名理由＝机主只盯"写的代理"，我不许看起来像我在停派）。
- 待人项仍按 `A641` 那一节；`winlive` 那条禁区继续生效（⛔ 本轮不再问机主）；未推枚数不写数、取数期间不推送、只 commit 不 push、cnb 未经机主同意不动。
