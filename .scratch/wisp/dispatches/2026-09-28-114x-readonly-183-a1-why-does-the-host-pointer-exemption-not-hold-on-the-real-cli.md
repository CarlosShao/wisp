# 派单 183-a1（**只读设计核·零跟踪改动**）＝裁"票 177 那枚精确豁免为什么在真机 CLI 上没护住宿主自己写下的指针"——四支候选各要**现量**，判"是哪一支"必须拿变异兑现（⚠ 变异只许走 `-overlay`，不许脏跟踪文件）

- 派单时刻：09-28 11:4x（`date` 现量于你起手第一步，**别抄我的号**）。
- 票面＝`.scratch/wisp/issues/183-…-cannot-be-read-back.md`（**AC#1 那一格就是本单**；AC#2..AC#8 是落地腿的，别替它做）。
- 现场证据（先读这三件再动）：`docs/evidence/s1/179-cli-reread-r2.md`（20739 字节）＋读数 `.scratch/wisp/probes/179/r2/logs/e2e-readings.txt`（28496 字节，**阻断那一行在第 51 行**）＋票 177 的八枚常驻判据 `internal/risk/shape_a_exemption_test.go`。
- ⚠ **预算硬顶 ≤45 次工具调用**；到顶即停手回禀并**自报枚数**。每裁完一段 commit 一次；**只 commit、绝不 push**。

## 0. 起手五件

`date`｜`git rev-parse --abbrev-ref HEAD`（**必须 `dev`**；不是＝停手上报。全程禁 `switch`／`checkout <分支>`／`merge`／`rebase`／`reset`／`stash`／`worktree`／`clean`／`restore`）｜`git rev-parse HEAD`（我这边 `95218c4c`）｜`git status --porcelain -- internal/ cmd/`（**必须空**）｜基线**只跑一枚**：`go test -count=1 -run 'ShapeA|177' ./internal/risk/`（⚠ **另一枚验收程 `179-v1` 此刻在飞、它会临时改 `internal/agent/**`** ⇒ 你**不许跑整包门禁、不许跑 `gate-clauses.sh`、不许跑 d22scan**，那些数不可归因；你要的凭据是**读码＋overlay 变异＋现成读数文件**）。

## 1. 四支候选，逐支给〔成立／不成立／无法判定〕＋可复制命令

**(a) 跨度定位偏**：`MarkWithHostPath` 那条路径窗口的坐标是在**归一化之后**算的吗？产物路径里有 `\`、盘符、`agent-task-<uuid>` 长串，而 `normalizeTaint` 会丢空白与零宽字符。⇒ 现读 `internal/risk/provenance.go` 里那段定位逻辑与 `taintmatch.go` 的 `runeIndexOf`／`windowExcluded`，并**造一发 overlay 变异**：把窗口坐标改成"原始下标"，看八枚常驻判据里哪几枚响。
**(b) 载具没送到**：真 CLI 上 `mark` 那一刻 `hostPathBox` 是不是空的？⚠ **这一支我最怀疑**：名册里那条 `ArtifactPath` 是 **`176-r1` 的回填**（`internal/tools/task_backfill.go`，跑完之后才写）填的，而 `task.go` 里那句 `box.set` 是**当场**执行的——两条路是不是同一枚 box、先后对不对，要拿**真 CLI 读数**证明（`179-r2` 那发就是现成的：`task.output` 成功、回执带路径、`fs.read` 仍命中 R4）。
**(c) 命中的不是路径窗口**：R4 匹到的片段会不会是那句**中文说明**（`全文见 `／`省略 17200 字符`／`注意：路径授权判定者未接线…`）而不是路径本体？⇒ 把回执那句拆成"路径子串"与"非路径子串"，各造一发 overlay 变异看谁响。
**(d) 索引外的第二条路**：`provenance.go:662`（路径参数从不豁免）与本次命中的关系——豁免做在**索引侧**，而命中判在**匹配侧**，两者之间有没有第二枚索引（片段之外的整串比对、`f.src` 的 `strings.Contains` 复核）把路径又捞回来？⚠ 票 177 §4.1 第 1 行写过"`f.src` 复核面不动 ⇒ 不会产生假命中"——**这句在真机形状上成不成立，本次要现量，不许照抄**。

## 2. 还要顺手答的三格（都是"下一程要知道的"，不是本单的判据）

1. **为什么八枚常驻判据全绿却挡不住这一发**——逐枚答"它测的是哪条接缝"（桥级／包级），并给一句"要长在哪一枚接缝上才拦得住"（候选：真 CLI `-overlay` 台件、桥级真具、还是 `internal/risk` 内一条端到端形状）。⚠ 这一格是本票最值钱的产出，**不许写成"测试覆盖不足"**。
2. **`dropped=1` 那枚线索**（尺：`grep -an "dropped=" .scratch/wisp/probes/179/r2/logs/e2e-readings.txt`）：run A `dropped=0`、run B `dropped=1` ⇒ 说明续读那一发自己那枚 `task.output` mark 在场。裁一句：豁免要作用的是**哪一枚 mark**（run B 自己刚产出的那枚，还是 run A 的？跨任务/跨 scope 时名册与 mark 的对应关系是什么）。
3. **行号重锚**：`cmd/wisp/run.go` 的行号被上一枚程的注释改动整体推了 3 行（`TaskDeps{` 组装点 362→**365**）。凡你要在报告里引的 `run.go`／`loop.go`／`task.go` 行号，**一律现跑 `grep -n` 取号**，别从票面抄。

## 3. 写面（超出即越权）

`docs/evidence/s1/183-exemption-why-not-live-a1.md`（你的表，新建）｜`.scratch/wisp/probes/183/a1/**`（overlay 变异件与读数；⚠ `logdir` 取脚本自身目录，**不许用 `runtime.Caller` 在 `-overlay` 下取路径**——它会报虚拟路径、把读数落进 `cmd/wisp/logs/`，`179-r2` 刚踩过并留了证）｜票 183 的 Progress log（**只追加、一枚框都不勾**）。
**禁改**：`internal/**`、`cmd/**`（**跟踪文件一字节都不许改**；变异只走 `-overlay`，跑完 `git status --porcelain -- internal/ cmd/` 必须仍为空并贴出来）、`docs/PLAN.md`、`docs/specs/**`、`docs/reports/**`、`allowlist.txt`、`thresholds.go`／golden／审批超时常量、别人的票面与证据件、`probes/**` 既有台件（只读）。
⚠ 现场躺着 `design/**` 未提交删除、`.gitignore`／`probes/152/my152.py`／`probes/161/r6/logs/flip-*`／`docs/evidence/s1/152-*.md`（` M`）⇒ **不提交、不还原、不删除、不评论**。`frontend/**` 与 `design/**` 别家归属，不读不写不引。
Git：`git add -A`／`git add .`／`commit -a` 一律禁；一步式 `git commit -q -F - -- <显式路径>`；禁 `--amend`；heredoc 加引号；含反引号的中文走编辑工具；串命令里可能无命中的 `grep` 用 `;` 或 `| cat`、**别用 `&&`**。临时件只建不删。

## 4. 交件报告（顺序固定）

1. step-0 五件；2. **本程没测什么**；3. **四支候选逐支判定**（每支：结论＋可复制命令＋变异落地证明＋两向读数）；4. **还原证明**（`git status --porcelain -- internal/ cmd/` 空，逐次贴）；5. 第 2 节三格的答案；6. **推荐落地点**（行，不是码）＋"要动哪一枚冻结件"（若判要动 `provenance.go:468-473` 或 `taintmatch.go:11-15` 那两句冻结文字＝**人工批准，具名报回别自己动**）；7. 被拒／没成功的调用（取数前还是后）；8. 有没有跑过删除命令；9. 工具调用**用了几次 vs 硬顶 45**；10. 伪授权两栏（各带出处）；11. 凭据值零抄录；12. `next=`（含一句：**`183-r1` 落地腿派之前还缺什么**）。

⚠ 每个"几枚／几行"旁边附可复制命令。
⚠ **我这段话里每条前提（含"豁免没生效"、"回填路径与 box.set 可能不同一条"、"命中的可能是中文说明"、"八枚常驻判据全在桥级/包级"）都是未验证断言**：不符就报回并继续做做得动的部分，**不许为了对我那句话去改判据或改测试**。
