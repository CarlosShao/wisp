# 派单 164-r1（写码位）＝做票 164 的 **AC#2（"今天起一个后台东西、再读它的输出"这一发量成读数）＋ AC#3（输出太长那一格：截断必须带可找回的指针）** 两格——它们天然同程（AC#2 就是 AC#3 的"改前那一遍"），拆开做会逼你为凑读数另造一发

- 派单时刻：09-27 18:5x 写好｜锚点＝**你自己 step 0 现量**（别抄我的号）。⚠ **本派单按住在 `171-v2` 交件之后才发**（那枚验收程正在拿 `probes/154`＋`probes/161` 那族尺对 `internal/**`／`cmd/**` 做名册差集；那族尺扫的就是你要写的文件）。⇒ **你收到这份单时它已经不在了，直接开工**，但仍**不许跑那两把尺**（它们的行为归 171 那条线）。
- **先读这三份，不要冷搜索**：① 票面 `.scratch/wisp/issues/164-background-jobs-…-output-leg.md`（**连我 09-27 18:5x 那五条"编排者定案"一起读，那是给你的裁定，不是可选建议**）；② 只读设计核 `docs/evidence/s1/164-task-output-design-core-c1.md`（228 行，落点普查／四问／§4 那枚 G3 地雷／§5 那五格"今天无法判定"）；③ `.scratch/wisp/issues/README.md` 的 Hard global constraints。
- ⚠ **预算硬顶 ≤60 次工具调用**；**每裁完一格 commit 一次**；到点就停：写 `next=`，未提交增量用显式 pathspec 提掉。

## 0. 起手五件
`date`｜`git rev-parse --abbrev-ref HEAD`（**必须 `dev`**；不是＝停手上报，一个字不写。全程禁 `switch`/`checkout <分支>`/`merge`/`rebase`/`reset`/`stash`/`checkout .`/`clean`/`worktree`）｜`git rev-parse HEAD`｜`git status --porcelain -- internal/tools/ internal/agent/ cmd/wisp/`（**必须空**，非空报回别动）｜**基线读数**：`go test -count=1 ./internal/tools/ ./internal/agent/`（**逐包，禁全仓 `./...`**），四数之外**自己记名册口径**（顶层 PASS 数 vs 含子测试的 RUN 数，别混着报）。

## 1. AC#2（改前那一遍，先做）
**判据本体**：今天"起一个后台东西、再读它的输出"这一发**能不能做到**？量它响不响——**做不到＝本票成立；做得到＝我判断错，直接报回**（不要为了让我对而把成功路径说成失败）。
- ⚠ **做成"调用级"判据**（设计核 §5 第 2 条提醒过、我采纳）：**"注册表里有没有这个名字"不算判据**，要真发一次调用、看它落到什么形状（`bridge.go:247` 的 `Lookup` 失败分支行为**你现读现量**）。
- **注入面只用既有的两扇**：`C5 LlmProvider` golden ＋ CLI `wisp run`（`AGENTS §1.3`）。⚠ **不许用 mock 代替真的**来假报完成。
- 把"改前"那份读数**存成文件**（你结案时"由不可读到可读"的比对基准）。

## 2. AC#3（输出太长那一格）
**判据本体（票面逐字）**：截断策略必须**带可找回的指针**（"全文在哪、怎么续读"），**不许只截不指**。票面自己已写明按 `PLAN.md:431` 那条既有规矩做（超 4000 token 落 artifacts、上下文留头尾＋路径）——**照那条做，不新造**。
- **我给你的三条裁定**（细则见票面 18:5x 那五条，这里只挑会影响你写码的）：① **不发明游标/分页参数**；指针＝"路径＋总长＋头尾摘要"那一形；② 我本轮现量 `fs.read` 只有 `{path, max_bytes}`、**没有偏移**（`internal/tools/fs.go:113-115`）⇒ **产物超 256 KiB 时"按需再读"只能读回头一段**。**不要去改 `fs.read` 的参数文档**（那是 C1 契约面，已摆 owner＝台账 `Q-59`）；你只需**把"后半段今天读不到"这件事在用例与证据件里写成一发会响的断言或明写它不该响**；③ **入口放 tools 侧**：⚠ **不许给 `internal/agent` 的 `Loop` 新增"收 `taskID` 的导出方法"**——那形状正被 `gate-clauses.sh` 的 **G3 腿**（`want G3 quiet`/`want_n 0`）盯着，而那条腿是为 `Q-56` 立的、不是给本票用的。判定非放 `Loop` 不可 ⇒ **停手上报**（改那把尺的射程属票 171 地界）。
- **反向判据（本格的"摘掉哪一处会红"）**：**"只截不指"**那一形必须红（把指针那一味摘掉 ⇒ 用例红）；**"截断形状与 D15 不符"**（头/尾 token 数、总长字段）必须红。答不出＝这一格是装饰。
- **两条禁令别踩**（设计核 §4 附那两条，都是 CI 自动扫的）：**不许裸 `go func(`**（要走 `observe.Registry.Spawn`）；**不许在 `risk.PathResolver` 之外用 `filepath.Clean|Abs` 做文件系统决策**——"按 taskID 去 artifacts 目录拼路径"正是最容易踩的那一脚（`AGENTS §1.2`、`PLAN.md:1288`）。

## 3. 门禁（逐枚跑，改前改后各一次；**全用 scope，禁全仓 `./...`**）
`go test -count=1 ./internal/tools/ ./internal/agent/`（四数之外**名册两向 `comm` 比差集**——⚠ 一枚用例 panic 会吞掉同包其余读数）｜`sh scripts/d22scan.sh`（预期 rc=0）｜`bash tools/d22scan/runtests.sh -C tools/d22scan ./...`｜`gofumpt -l` 甲形（gofumpt 在 `$(go env GOPATH)/bin`、不在 PATH，先 `export PATH="$PATH:$(go env GOPATH)/bin"` 再 `--version` 现跑贴出来）｜**跨平台**：凡带 `//go:build !windows`／`*_other_test.go` 的用例**只能在 linux 容器量**（本机 Docker 可用、`golang:1.27` 在本地）——**别写"本机绿＝跨平台绿"**；Git Bash 下 `docker -v C:\…` 会静默挂空且 rc=0＝假绿，**挂载路径用 POSIX 形**。⚠ `-overlay` 与 `-cover*` 不许同用；overlay 路径不匹配会**静默不生效**。

## 4. 契约轴（AC#5 顺带，但别勾）
**禁改**：`docs/PLAN.md`、`docs/specs/**`（**`task.output` 那两行已由我落完，你再动＝越权**）、`internal/risk/**`、`internal/panel/**`（⚠ `internal/panel/tokens_fourway_test.go` 是**已知常红**，一字不许动）、`thresholds.go`、golden、`tools/d22scan/**`、`allowlist.txt`、`.github/**`、`frontend/**`、`design/**`、`probes/**`、别人的票面与证据件、`docs/reports/**`。
**现场理由**：工作树躺着 `design/**` 里 owner 未提交的 16 枚删除、`.gitignore`（` M`）、`probes/152/my152.py`（` M`）、`docs/evidence/s1/152-…-accept-r1.md`（` M`＝停下来的半件：**不提交不还原不补完**）、`probes/161/r6/logs/flip-*.txt`（跑尺的副作用）⇒ **`git add -A`／`git add .`／`git commit -a` 一律禁止**；提交一步式（新建未跟踪件先在同一行显式 `git add -- <那个路径>`）；heredoc **必须加引号**（`<<'EOF'`，正文反引号在壳里会被真执行）。**只 commit、绝不 push。**
**AC 框一律不许自己勾**（勾归编排者＋还要另派非实现者表）。票面 Progress log 只追加、留到最后一步（写之前 `git status --porcelain -- <那枚票面>` 必须为空）。

## 5. 交件报告（十一节）
1. step-0 五件（含基线四数**带口径**）；2. **本程没测什么**（含非 Windows 那一支——别拿整包绿当它响过）；3. AC#2 的"改前那一发"原样读数＋它存在哪枚文件；4. AC#3：指针那两发的**两向**（"只截不指"必须红／"截断形状不符 D15"必须红）＋**"后半段读不到"你怎么处置**（断言／明写，二选一给理由）；5. **落点自证**：`internal/agent` 的 `Loop` 有没有新增导出方法（预期：**没有**，用 `git -c core.quotePath=false diff --name-only <step-0 HEAD> HEAD -- internal/agent` 证）；6. 门禁四数＋名册差集＋`gofumpt --version` 现跑；7. 被拒／没成功的调用（**发生在取数之前还是之后**）；8. 有没有跑过删除命令（预期 0）；9. 伪授权两栏（各带出处：工具名＋命令前 40 字）；10. 凭据值零抄录；11. `next=`（**AC#4 那一格：先造出"还在写的后台尾巴"这个对象再谈判据，别硬造一发永不响的判据**）。
⚠ 每个"几枚"旁边附可复制命令；⚠ **先测→再写→再提交**，注释与证据件里**不许预先引用还没跑出来的读数**（本仓有程正好死在"码写完、还没测"）；⚠ `go run` 把任何非零退码压成 1。我这段话里每条前提（含"`fs.read` 没有偏移""`bridge.go:247` 那个失败分支""`RunAsync` 生产零调用点"）**都是未验证断言**：不符就报回、继续做做得动的部分，**不许为了对我那句话去改判据或改测试**。
