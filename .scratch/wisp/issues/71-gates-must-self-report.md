# 71 — 每道门必须自报工作量：把"没报问题"和"没看"变成可区分的两件事

**Status:** in progress（我所有路径内的部分已落地：台账 + runner + CI 三处接入 + A43；AC#1/#2 的
`cmd/balldebug` 那半张票**不在本代理所有路径内**，AC#5 的"CI 真红"因**禁止 push** 而物理不可证；均见 Progress log）
**Claimed by:** ticket-71 implementer（2026-09-21，负责 `tools/d22scan/**` 与 `.github/workflows/ci.yml`）
**Last update:** 2026-09-21 10:5x（AC#3 已勾；AC#4 差一个 ban #6 处置的裁决；AC#1/#2 的 balldebug 那半张票移交；A43 两条调查已落 registry）
**Blocked by:** 70-ci-actually-green（只因为两者都要改 `.go`，不是语义依赖；票 70 的 gofumpt 已落地）
**Parallel slots:** ≤1 sub-agent（本票改的是**门禁自身**，改错会让后面所有票的"绿"都失去意义）
**Spec refs:** README 完成规则 6（1:1 裁决表）、D22（契约变更需人批）、C31（可观测性）

## What to build
今天（2026-09-20）连续三次踩到同一个形状：**一道门跑完了，但它的输出无法证明它看过任何东西。**
最贵的一次是我自己：在仓库根跑 `go run ./tools/d22scan -root .`，它只打印一条**模块错误**就退出，
我把"没有 findings"读成"clean"，并据此**否证过一个代理的正确报告**（A22），
又把这句假绿写进了三张 done 票的验收报告（A30②）。
**"外部扫描器没报错" 不等于 "跑过了"**——这条已经进了长期记忆，但记忆只对我有效；
**能挡住下一个会话的东西是让门自己把工作量打印出来并在 0 时致命退出**。

票 67 已经给 `d22scan` 装了这台仪器（`examined N production Go files` + `N==0` 致命退出，
R16 裁定 3 **禁止回退**）。本票把同样的纪律推到**其余三类仪器**上。

## 本轮挖出的空仪器（逐条都有复现路径，不要重新发现）
1. **`cmd/balldebug -diff` 无论像素计数为何都返回 nil**：差分表的 `px_delta_ge8`、`timers`、`cpu`
   全部只是**打印**，主函数结束时的 exit code 与这些数字**无关**（相关代码见
   `cmd/balldebug/diff_windows.go:365` 的打印行与 `:495` 的 `fileHas(..., "timers=true")`）。
   ⇒ 一次"球物理上不成像"的跑（票 62 的 `px_delta_ge8=0`）在 CI 眼里和理想跑**完全一样**。
2. **`go test -run <不匹配的模式>` 打印 `ok`**：Go 自身语义。今天票 64 的 AC#4 就是被这个形状
   坑过一次（"14 项全 PASS" 里 SKIP 也算 ok）。⇒ 我们不能改 Go，但可以**加一层包装断言**。
3. **`d22scan` 的 `frontend/` 作用域实际走 0 个文件**：票 67 之后它虽然自报总数，但**分作用域**的
   覆盖面仍无人看；一个恒为 0 的作用域就是"看着在扫、其实没扫"。

## Acceptance criteria
- [ ] **AC#1**：`cmd/balldebug -diff` 有一个**阈值门**模式——当任一行 `px_delta_ge8` 低于给定的
      最小可见性阈值、或缺行、或 `timers` 非 `no`（在要求零定时器的态上）时**非零退出**；
      阈值**写死在票面/配置里**并带来源（D43/SPEC-08），**不许为了让某次跑变绿而现场调低**（D22 精神）。
- [ ] **AC#2**：阳性对照（**这条不做本票不算完**）：构造一次**故意不成像**的跑（把球隐藏或把阈值调到
      不可能达到），证明 AC#1 的门**会红**。今天的教训正是"没红"被当成"过"。
- [x] **AC#3**：一个可复用的测试运行器（脚本或 `cmd/` 下的小工具），跑 `go test -tags <X> -v <pkg>`
      后**断言 `--- PASS/FAIL/SKIP` 行数 > 0 且 `-run` 模式实际匹配到测试**；
      匹配 0 条时**非零退出**并打印它匹配到了什么。**不得**把 SKIP 计入"通过"。
      ⇒ **落地物 `tools/d22scan/runtests.sh`**（票面简报只给了我 `tools/d22scan/**` 与 `ci.yml` 两处，
      所以它没去 `scripts/`；调用点在 CI 里只有一行，编排者若要迁进 `scripts/` 是一句 `git mv`）。
      实测见 Progress log 10:4x 条：漂移名 `-run` 裸跑 **rc=0** / 经脚本 **rc=1**，SKIP 一律 rc=1，
      真失败原样传递 rc=1，`-tags` 透传已验。CI 的三处接入：`lint` 的阳性对照步、`test-core` 的
      `TestLayoutForTestEnv`、`test-windows` 的 `TestPathResolverJunctionWindows`。
- [ ] **AC#4**：`d22scan` **逐作用域**自报 `examined N files`（含 `frontend/` 这类子作用域），
      任何**声明了但恒为 0** 的作用域 ⇒ 致命退出；或者**删掉那个作用域**并在 allowlist 留一行说明。
      现存的 `allowlist.txt` 只允许**变短或不变**，变长需我在票面批。
- [ ] **AC#5**：把上面三处接进 CI（`lint` 或 `test-windows`），并给出**一次真实的红**：
      在分支上故意造一条不合格产物、贴出 CI 红的那一行，然后回滚。**只贴"本地能红"不算**。
- [ ] **AC#6**：本票结束时，`docs/reports/pending-and-issues.md` 的 **A30⑤** 追加一句"已闭于票 71"，
      并把本票新挖到的同类空仪器（如果有）继续登记在 A30 之后，**不要改已完成票的证据文件**。
- [ ] **AC#7**：对抗验收由非实现者执行，验收报告含与上述 AC **1:1 的裁决表**（README 规则 6），
      每行给出可复现命令与真实 exit code。

## 明确不在本票范围
- 不放宽任何既有阈值（D32 `Sleeping` CPU ≤0.5% 一动没动，也不许动）。
- 不改 `internal/risk/**`、`rules_gateway.go`、`docs/PLAN.md`、`docs/specs/**`（冻结契约）。
- 不为了 CI 变绿去加豁免：**只许从严**（R16 裁定 4）。

## log
- 2026-09-20 23:5x 编排者建票。**为什么现在建**：A30 的全量证据自查一次性给出三类空仪器，
  我此前只在报告里"建议新立一张小票"，那等于把推迟项写成口头技术债（已在长期记忆里被批评过一次）。
  完成判据 = AC#1..AC#7；当前残缺表现 = **今天所有依赖 `-diff` 与 `go test -run` 的"绿"，
  其证明力都限于人肉读输出，机器不会替我们拦**。

## 编排者追加（2026-09-21 09:41，从票 67/67b 移交过来，不要重做它们的已交付部分）

本票的主题正是"门必须自报自己走了多少文件"。票 67 交完覆盖面扩展后留下三条**同族**的活，都归本票：

1. **ban #6 仍指向 `frontend/`** —— 那棵树在本 HEAD **不存在**，所以那条禁令是**恒 0 文件的死作用域**。
   ⚠ 不要"顺手缩掉"：禁令文本属 D22 范围，要动先报编排者。可先做的是让 `emptyScope` 守卫也覆盖 ban #6。
   > 2026-09-21 票 71 implementer：**已做，且做的比"覆盖"更多**——禁令文本一字未动；ban #6 现在在
   > 台账里，每次运行打一行 `[NOT COVERED]` + 一条解释，且"树出现但豁免还在"= rc 2（自己过期的豁免）。
   > 待你裁决：这算不算 AC#4 的第三种处置（判据字面是"致命退出 或 删作用域+allowlist 留一行"，
   > 前者会让 CI 为代理无权修的原因永久红，后者要给 allowlist 加行 ⇒ 违反"只许变短或不变"）。
2. **`.github/workflows/ci.yml:23` 的注释仍写 "design/ and frontend/"**，与 `emojiScopes()` 现在的真相源
   （`design/` + `internal/` + `cmd/`）不符。注释归票 70 的文件 ⇒ **等票 70 落地后同批改**，别提前撞车。
3. **新判据（今天用一次红 CI 换来的）**：**放宽/扩大一个门禁的覆盖面，必须在同一批 commit 里改完它新照到的
   存量违规**。实测反例：`bcf44d6` 把 ban #8 扩到 `internal/`+`cmd/` 后，
   5 分钟前入库的 `d63bc49` 里那条注释 `⚠`(U+26A0) 立刻变成真命中 ⇒
   `git archive HEAD` 纯净树 **1 finding / exit 1**，**CI lint 红约 13 分钟**，
   直到 `a8ae9ad` 用 1 行纯注释 ASCII 替换收掉。
   **本票 AC 判据里要加一条**：任何"扩覆盖面"的改动，收尾时必须跑一次**纯净树**
   （`git archive HEAD` 或等价手段）证明它是绿的——在有别的代理并发的共享树里，
   `git status` 干净**不等于** HEAD 干净。

## Progress log（票 71 implementer，`tools/d22scan/**` + `.github/workflows/ci.yml`）
- 2026-09-21 09:5x 开工。**已闭**：编排者追加第 2 条（`ci.yml:23` 注释仍写 "design/ and frontend/"）——
  注释改为**指向真相源函数** `emojiScopes()` 而不是复制一份清单，因为复制出来的清单本身就是这次要修的
  病灶（话术与覆盖面脱钩）。**移交出去**：AC#1/AC#2 的 `cmd/balldebug -diff` 阈值门不在本代理所有路径内
  （编排者只给了 `tools/d22scan/**` 与 `ci.yml`），且 `internal/ball` 正被票 74 代理占用 ⇒ 交回编排者派活。
  本代理按同一条**精神**在自己所有物上做等价物：阳性对照必须是**真的红**（AC#2 的形状），不是绿。
  **进行中**：AC#4 的逐作用域自报——实测基线（`cd tools/d22scan && go run . -root <abs>`，
  exit 0）今天只自报了 ban #8 的 3 个作用域，ban #6 `frontend/`、ban #7 `internal/tools/`
  与 `internal/` / `cmd/` 各自的 Go 文件数**一条都没报**。⇒ 下一步：给 d22scan 装**统一作用域台账**
  （每个声明作用域 `examined N`）+ 把 `emptyScope` 守卫推到 ban #6/#7，并保证只严不宽。
- 2026-09-21 10:2x **AC#4 的 d22scan 部分已落地并有真红**（commit `b551fef`）。台账 `declaredScopes()`
  现覆盖 8 个作用域，每个都在 stdout 打 `scope <label> examined N <kind>`；`emptyLiveScope` 把
  致命退出从"只有 ban #8"推到全部 live 作用域；新增 `undeclaredKeys`（扫了却没上报 ⇒ rc=2）与
  `driftedAbsentScope`（登记为"树不存在"而树已出现 ⇒ rc=2，让豁免**自己过期**）。
  **实测（工作树，非纯净树）**：`go test ./...` rc=0，`=== RUN` **26**、`--- PASS` 18 + 子测试 8 = 26、
  `--- FAIL` **0**、`--- SKIP` **0**；`go run . -root <abs>` rc=0 打印
  internal/=184 / cmd/=16 / ban#7=16 / ban#8 design=16 internal=293 cmd=21，ban #6 一行 `[NOT COVERED]`。
  **阳性对照（AC#2 同族，`go test -v` 的 4 条 e2e 子测试）**：编译**真二进制**后跑进程退出码——
  种一条 `go worker()` ⇒ rc=**1**；把 `internal/tools/` 清空 ⇒ rc=**2**；造出 `frontend/` ⇒ rc=**2**；
  全 live fixture ⇒ rc=0。**allowlist.txt 未动**：条目数 5 → **5**（`grep -c '\t'` 实测）。
  **AC#4 框仍不勾的唯一原因**：ban #6 我走了 AC 给的**第三个**处置——既不致命（那会让 CI 为一个
  代理无权修的原因永久红：删禁令文本是 D22 范围）也不删作用域（那会丢掉面板落地当天就会响的匹配器，
  且需在 allowlist 加一行说明 ⇒ 违反"只许变短或不变"），而是**登记为 exempt + 每次运行打印
  `NOT COVERED: ban #6 frontend/` + 树出现即 rc=2**。字面判据是"致命或删掉"，二选一我都没做满，
  所以框留给编排者裁决。⇒ next=纯净树（`git archive HEAD`）复跑证 HEAD 绿，然后把 AC#5 的"CI 真红"
  缺口如实写进票面（本代理被明令**不得 push**，AC#5 的"贴出 CI 红的那一行"在物理上不可证）。
- 2026-09-21 10:4x **AC#3 落地 + CI 次序修正 + A43 两条调查结论落盘**。commit：`b551fef`（台账）与
  本条所在 commit（runner / ci.yml / registry）。**AC 账实表（1:1，README 规则 6）**：
  | AC | 框 | 命令与真实数字 |
  |----|----|----------------|
  | #1 balldebug 阈值门 | ☐ 未做 | 不在我被给的所有路径内（`cmd/balldebug`），且 `internal/ball` 正被票 74 占用 ⇒ **移交编排者派活** |
  | #2 阳性对照（一次真红） | ☐ 框留给 #1 的红；**同族等价物已做** | 编译真二进制跑进程退出码：种 `go worker()` ⇒ **rc=1**、清空 `internal/tools/` ⇒ **rc=2**、造出 `frontend/` ⇒ **rc=2**、全 live fixture ⇒ **rc=0**（`TestBuiltBinaryGoesRedEndToEnd` 的 4 个子测试全 PASS） |
  | #3 测试运行器 | ☑ | `sh tools/d22scan/runtests.sh ./internal/proc/ -run TestLayoutForTestEnv` ⇒ `PASS=1 FAIL=0 SKIP=0, === RUN=3` **rc=0**；名字漂移版 `-run TestLayoutForTestEnvRenamedAway` ⇒ 裸 `go test` **rc=0**、经脚本 **rc=1**（`PASS=0 FAIL=0 SKIP=0, === RUN=0, '[no tests to run]'=1`）；SKIP fixture ⇒ **rc=1**；真失败 ⇒ **rc=1** 原样传递；`-tags windows` 透传 ⇒ **rc=0**；`-C /tmp`（无 go.mod）⇒ **rc=2** |
  | #4 逐作用域自报 | ☐ **差一个裁决** | 台账 8 条全打数（`bans #1-5 internal/=184 cmd/=16`、`ban #7 internal/tools/=16`、`ban #8 design/=16 internal/=293 cmd/=21`），live 作用域 0 文件 ⇒ rc=2，另有"扫了未登记 ⇒ rc=2""豁免过期 ⇒ rc=2"两条守卫；`allowlist.txt` 条目 **5 → 5** 未增。**不勾的原因**：ban #6 我走了 AC 字面（致命 或 删作用域+allowlist 加一行）之外的第三种处置，两难写在"编排者追加"第 1 条的引用块里 |
  | #5 接进 CI + 一次真实红 | ☐ 物理不可证 | 接入已完成（3 处 + 阳性对照步提到 vet 之前）；**"贴出 CI 红的那一行"需要 push，而我被明令不得 push** ⇒ 这一半不勾，不拿"本地能红"冒充。顺带把 `lint` 次序的问题登记为 **A43②**：自 `fd8f838` 起 D22 扫描步在 CI 上**从未产出过一个结论**（被前面的 `go vet` 失败跳过） |
  | #6 registry A30⑤ | ☑（按"部分闭"写） | A30⑤ 追加逐条三行（a 闭 / b 闭 / c 未闭）+ 新增 **A43**；`git diff --numstat` 实测 **78 insertions / 0 deletions** ⇒ 未销毁任何历史条目（我第一版编辑误删了 A42 的标题行，已当场补回并用 numstat 复验 0 删除） |
  | #7 对抗验收（非实现者） | ☐ 不归我 | 本表是**实现者自账**，按规则 6 必须由编排者或另一代理独立重做变异后才算验收 |
  **纯净树（编排者追加第 3 条）**：`git archive HEAD | tar -x -C /tmp/wisp-head-<pid>`（**树外**，避免 A38 的
  嵌套目录让并发代理多算一份源码）后跑 `go run . -root <abs>` ⇒ **rc=0**、`go test ./...` ⇒ **rc=0**
  （HEAD=`b551fef`）。同一棵纯净树交叉 vet 复现 `undefined: mulA` ⇒ 见 A43②（别人路径，只登记不动）。
  ⇒ next=等编排者对 ban #6 处置的裁决；把 `cmd/balldebug`（AC#1/#2 那半张票）派给能碰它的代理；
  本代理剩余额度只用于文档/票面，不再扩路径。
- 2026-09-21 11:0x **gofumpt 自查（CI lint 第一步真实会拦的东西）**：本地装上 `gofumpt@latest` 后
  `gofumpt -l` 点了**我自己两个文件**（`tools/d22scan/main.go` 的 composite literal 换行、`scan_test.go`
  末尾两行空行）⇒ 已 `-w` 修掉，模块目录 `gofumpt -l .` 现为空，`go test ./...` 仍 **ok / rc=0**。
  **同一次扫描顺手量到 HEAD 上还有一个不属于我的 gofumpt 红**：`internal/memory/artifacts_path_invariant_test.go`
  （由 `d9224af`，票 76 的 checkpoint 引入）⇒ **不代改**（该代理正在飞，同文件并发是假并行），登记给编排者。
  **纯净树复跑（HEAD=`8dd43d0`，含我两个 commit）**：`sh scripts/d22scan.sh`（CI 的原命令）⇒ **rc=0**，
  8 条作用域全打数；`sh tools/d22scan/runtests.sh -C tools/d22scan ./...` ⇒ **rc=0**，
  `PASS=18 FAIL=0 SKIP=0, === RUN=26`。⇒ next=交回编排者：ban #6 处置裁决、balldebug 半张票派活、
  memory 的 gofumpt 红、A43② 的两个符号归属。
