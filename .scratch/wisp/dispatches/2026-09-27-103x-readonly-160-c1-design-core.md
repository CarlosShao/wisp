# 派单 160-c1（**只读设计核**，不写产品码）＝把票 160（`OpenScope` 交回自带关闭动作的句柄）从"票面上写着推荐"变成"能派代码程"

- 派单时刻：09-27 10:3x｜owner 09-27 09:30 已批本票的**开地**（台账 `A322`：批准单位＝"放开 `internal/risk/provenance.go`（含测试）＋两枚唯一生产调用点的必要改动"；撤销口令「160 别动」）。
- 票面＝`.scratch/wisp/issues/160-scope-open-returns-a-handle-carrying-its-own-closer.md`。**"我的推荐"那一节＋AC#2 的三条硬约束（不新增 goroutine／不许墙钟差／`RiskLevel` 语义一字节不许变）是本单的边界**。
- ⚠ **本程只读**：一字节产品码都不改。⚠ **预算硬顶 ≤55 次工具调用**；到点停手→写 `next=`→把证据件用显式 pathspec 提交（**每核完一节写一节并提交**）。

## 0. 起手四件
`date`｜`git rev-parse HEAD`（锚只认自己现量）｜`git status --porcelain -- internal/ cmd/`（**非空就报回别动**）｜`Read` 票面并把 AC#1–AC#5 连行号抄进证据件。

## 1. 要交回来的六件事（每件都要"命令＋读数"，不许只给推断）
1. **今天的形状现量**：`OpenScope`／`CloseScope` 的签名与**全部调用点**（含 `internal/tools/bridge.go` 的开合两处、`cmd/wisp/panel_assets.go` 那枚只开不合的在册活形状）；逐枚给 `文件:行号`＋那一行原文。
2. **依赖方向这道题（本单最要紧）**：清扫清单（`Defer`／`DeferNamed`／`Dispose`／`Incomplete`）在 `internal/plugin/disposal.go`，而"开作用域"在 `internal/risk/provenance.go`。⇒ **`internal/risk` 现在依不依赖 `internal/plugin`？（`go list -deps` 现跑两向）** 若依赖方向反过来（plugin→risk），那"开的时候就登记怎么关"**在 risk 层根本做不到**，只能由调用方（bridge.go）登记 ⇒ **"忘了关写不出来"这一许诺就只覆盖 bridge 那条腿、不覆盖 `panel_assets.go` 那条腿**。⇒ 把这个结论写成一句话，并说明它让票面 AC#2 的哪半句必须改写（**只登记，不改票面**）。
3. **两条设计支各自要动几枚文件**：（甲）关闭动作是**闭包里的局部量、不出包**（参照票面点名的 Pi `session.ts:1475-1490` 那一路）；（乙）`OpenScope` 返回一枚带 `Close()` 的结构体。逐枚列**要改的文件与行**，并各自答：**"摘掉认身份（只有发出去那枚能关）这一味，是否存在一发跨包误关从此打不红？"**
4. **同源拷贝枚数**：票面与台账写过的那几句复述（`provenance.go` 里契约注释那句、`:453` 一带的日志句、`bridge.go` 的注释复述、`cmd/wisp/panel_assets.go` 的复述）——**逐枚现量行号还在不在**，并数一把"改一把尺有几份自称逐字抄自它的拷贝"（本仓的老坑）。
5. **测试面预演**：现有 `provenance_test.go` 里"开必须配关"那类用例（票面点了 `:310`）有几枚、各自钉的是什么；并预演 **AC#3 那一问**：摘掉"认身份"这一味，**有没有任何外部可见读数会变**（答不出＝装饰）。
6. **DEFERRED(C25-loop-wiring) 的兑现面**：该标记在码里的每一处现量位置＋`SPEC-12 §5` 登记表那一条的五字段，与本票只收第 (1) 条的边界；⚠ **不许顺手把它别条一起"顺便兑现"**。

## 2. 不许做／不许结论化
不改任何文件；不新增 goroutine；不用墙钟差；不动 `thresholds.go`／golden／`allowlist.txt`；**不给"那就选甲"这种结论**——那是我与 owner 的事，你只交**两支各自要动什么＋各自会红在哪**。
`frontend/**`、`design/**` **不查、不引用、不转述**（不是本编队的地界）。

## 3. 写面（超出即越权）
`docs/evidence/s1/160-design-core-c1.md`（**渐进写**）＋ `.scratch/wisp/probes/160/c1/**`（只放读数与命令记录，落盘前 **gofumpt 干净**——它会被算进全仓乙形归因）。
**禁改**：`internal/**`、`cmd/**`、`tools/**`、`.github/**`、`docs/PLAN.md`、`docs/specs/**`、所有票面、`docs/reports/**`、`allowlist.txt`、golden、`frontend/**`、`design/**`。
**现场理由**：`design/**` 有 owner 未提交的删除；`probes/152/my152.py`（` M`）与 `docs/evidence/s1/152-…-accept-r1.md`（3 行未提交）是停下来的半件——不提交、不还原、不补完 ⇒ `git add -A`／`git add .`／`git commit -a` 一律禁止。

## 4. 仪器坑（实测过）
`grep -c` 命中 0 ⇒ rc=1 吃掉 `&&`｜**引别人的行号先按他的锚取**（`git show <锚>:<文件>`），漂移多半是我拿 HEAD 比出来的｜"某文件里有 N 处"必须**逐枚点名**再给数，分组结构天然漏尾项｜`go list -deps` 在两枚不同 module（`tools/d22scan` 自成模块）上会给你不同的图，先确认你在哪个模块里问。

## 5. Git 与噪声
只 commit 不 push；显式 pathspec ＋ **加引号的 heredoc**；提交后 `git show --name-only <自己号>` 出现别人路径＝停手报回；禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；**不许跑删除命令**（`rm` 每次都会在真人机器上弹授权窗）。

## 6. 交件报告（顺序固定）
1. step-0 四件；2. **本程没核什么**；3. 六件事逐件（每件带命令）；4. **对编排者的不服**（我给的"Pi 那一形状／依赖方向猜想／拷贝枚数"都是未验证断言）；5. 被拒／没成功的调用（取数前后）；6. 有没有跑过删除命令；7. 伪授权两栏计数（各带出处）；8. 凭据值零抄录；9. 档位（四档）；10. **`next=`＝一份"代码程派单就绪包"**：文件级写面名单＋每格可复算的会响判据草案＋必须回报的未定义项清单。
⚠ 我这段话里每条前提都是未验证断言；**不成立就报回并继续做做得动的部分，不许为迁就我去改判据**。
