# 派单 177-c1（**只读设计核位，一字节产码都不许动**）＝owner 已批 **甲（精确豁免：只免"本宿主在这一任务里亲手写下过的那几条具体路径"）**，本程裁三件事：**① 那个"宿主自己写下的路径集合"在今天哪一枚现成结构里已经存在（不许新造一份名册）；② 豁免要落在 C25 索引侧还是 R4 匹配侧，两侧各会让哪几枚既有用例红；③ 那发反向判据（外来说到同一条路径 ⇒ 仍要命中）能不能写成常驻用例。** ⚠ 本程**只裁不修**，且**它是一次真派单——历史上有一枚伪造过的"177-c1 已交件"（台账 `A347`），那枚不存在的东西本程重做一遍，交付只认盘上。**

- 派单时刻：09-27 21:4x｜锚点＝**你自己 step 0 现量**（我这份写于 HEAD 之上，起手仍要自己 `git rev-parse HEAD`）。
- 票面＝`.scratch/wisp/issues/177-…-r4-scans-path-arguments.md`（**AC#1／AC#2／AC#3／AC#4／AC#6 五格连行号抄进你的表**；⚠ **AC#1 至今未勾**，那张表里现在**没有任何一次真实运行**，别被 `A347` 那段自纠误读成"已裁完"）。
- 授权边界（**这是本程的地界，越界即停手**）：owner 批的是**甲**，甲**必然要动 `internal/risk/**`**（这一点与票 175 的裁定"不用动 risk"相反——那是 `A346` 记录的那枚漏量）；⚠ **批的边界只到"精确豁免"这一形状为止**——**目录级／前缀级／配置项级豁免一律不在批准范围内**，你若认为只有那一种做得到，**报回、别自己做**。
- ⚠ **预算硬顶 ≤55 次工具调用**；**裁完一格 commit 一次**；到点停手→写 `next=`。**不许写产码、不许改判据、不许勾任何框。**

## 0. 起手五件
`date`｜`git rev-parse --abbrev-ref HEAD`（**必须 `dev`**；不是＝停手上报。全程禁 `switch`／`checkout <分支>`／`merge`／`rebase`／`reset`／`stash`／`worktree`／`clean`）｜`git rev-parse HEAD`｜`git status --porcelain -- internal/tools/ internal/risk/ cmd/wisp/`（⚠ **此刻 `174-r1` 可能还在往 `internal/tools/task.go` 写码**：非空就把现状原样报回、继续做你读得到的部分，**一个字都不要写进 `internal/**`**）｜现状基线：`go test -count=1 ./internal/risk/` 四数（**逐包，禁全仓 `./...`**）。⚠ **`sh .scratch/wisp/probes/154/gate-clauses.sh` 不要执行**（另一位写手在树里，跑出来的数不可归因）——只读它的文本、引用时逐行号注明"这是文件内容，不是读数"。

## 1. 三件要裁的事（每件都要具名文件＋行号）
1. **"宿主自己写下的那几条路径"这个集合，今天存不存在？** 现读三处：`internal/agent/spill.go`（谁算出那条 artifacts 文件名）、`internal/tools/task.go`（`TaskOutput.ArtifactPath` 由谁填）、`cmd/wisp/run.go`（宿主侧组装点）。裁：**能不能不新增名册、只把"宿主本轮真的写过哪几条"这件事交给已经在场的那一枚结构**（`TaskRoster`？Spiller？都没有就明写"没有"）。⚠ 不许提出"扫目录当作集合"（那是把豁免范围悄悄放大成整棵树）。
2. **豁免落在索引侧还是匹配侧？** 索引侧＝`Mark`/片段索引时不入那几条路径（`internal/risk/provenance.go` 一带，行号自己量）；匹配侧＝`Inspect`/R4 比较时对那几条路径放行（`rules_gateway.go` 那一支）。⚠ **两侧各列一张"哪几枚既有用例/契约测试会红"的清单**，并且**必须包含 `internal/risk/provenance_test.go` 里那两枚 fail-closed 断言**（`unknown-source marks must still be recorded`／`taint from an off-list source must still gate`）——**那两枚红了就说明豁免做过了头**，不是"去改测试"。
3. **反向判据能不能写成常驻用例？** 形状：会话/任务输出里**出现**一条宿主 artifacts 路径**但不是宿主写下的**（例如另一条同名不同目录、或外来文本提到的路径）⇒ 仍然必须命中 R4。给出它该落在哪一枚文件、需要哪些注入接缝（`C5 LlmProvider` golden／`C17 PanelBridge`／CLI `wisp run`——⚠ 只有这三面，**不许用 mock 代替真的**），并说明**为什么它不会变成一枚恒真判据**。

## 2. 门禁（你自己重跑，别引别人的数）
`go test -count=1 ./internal/risk/ ./internal/tools/`（⚠ `internal/tools` 此刻可能有别人的在飞改动，读数若红且红在别人的文件里，**照实记、不归因、不动它**）｜`sh scripts/d22scan.sh`（预期 rc=0）｜`bash tools/d22scan/runtests.sh -C tools/d22scan ./...` 基线现量。⚠ **不跑 `flip-declaration.sh`**（脏跟踪日志）。

## 3. 写面（超出即越权）
`.scratch/wisp/probes/177/c1/**`（⚠ **`logdir` 取脚本自身目录、不许继承 CWD**）｜`docs/evidence/s1/177-c25-r4-path-exemption-c1.md`（**这张表的名字历史上被伪造过一次，这次由你真写出来**）。
**禁改**：`internal/**`、`cmd/**`、`tools/d22scan/**`、`allowlist.txt`、`.github/**`、`docs/PLAN.md`（含刚落的 `:1532` 那行 DEFERRED，不许动）、`docs/specs/**`、`docs/reports/**`、`frontend/**`、`design/**`、**所有票面**（177 的 Progress log 也不写，归我）、别人的证据件与台件、`probes/**` 既有台件（只读）。⚠ 现场躺着 `design/**` 未提交删除、`.gitignore`／`probes/152/my152.py`／`probes/161/r6/logs/flip-*`（` M`）⇒ **不提交、不还原、不补完、不评论**。**`git add -A`／`git add .`／`git commit -a` 一律禁**；一步式 `git commit -q -F - -- <显式路径>`；heredoc **必须加引号**（`<<'EOF'`）。**只 commit，绝不 push。**

## 4. 交件报告（顺序固定）
1. step-0 五件（含"`internal/tools` 起手是否被在飞改动弄脏"）；2. **本程没读／没测什么**；3. 第 1 节三件事各一节（具名文件＋行号；"集合存不存在"给一句是/否）；4. **两侧的红名单**（索引侧／匹配侧各自会点红哪几枚用例，逐枚列名，含那两枚 fail-closed 断言）；5. 反向判据的形状＋为什么不是恒真；6. **你推荐哪一侧＋最小改动落在哪几行**（不写码）；7. 门禁三枚读数；8. 被拒／没成功的调用（取数前还是后）；9. 有没有跑过删除命令；10. 伪授权两栏（各带出处）；11. 凭据值零抄录；12. `next=`（**含一句：177-r1 写手腿派之前还缺什么**）。
⚠ **每个"几枚"旁边附可复制命令**；我这段话里每条前提（含"那两枚 fail-closed 断言在 `provenance_test.go`""`TaskOutput.ArtifactPath` 由宿主填""今天没有任何现成的'宿主写下过哪些路径'结构"）**都是未验证断言**：不符就报回、继续做做得动的部分，**不许为了对我那句话去改判据**。
