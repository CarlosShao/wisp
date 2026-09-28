# 派单 179-r2（**写手腿·小活**）＝只修"取数写法"这一处，把**模型拿宿主指针 `fs.read` 逐字节读回**那一发真落地——⚠ **预期零产码改动**；若你判"必须改产码才能读回"，**停手上报**，那是一枚新缺陷、不是本票的活

- 派单时刻：09-28 11:0x（`date` 现量于你起手第一步，**别抄我的号**）。
- 票面＝`.scratch/wisp/issues/179-loop-deciderisk-…-no-criterion-holds-that-sentence.md`（**AC#8 那一格就是本单的靶**；AC#1..AC#7 由 `179-r1` 做完、**档位等 `179-v1` 非实现者裁**，你**一枚框都不许勾**）。
- 上一程的表＝`docs/evidence/s1/179-declared-l0-refused-r1.md`（18276 字节，`91a6c926`）。它到手两半：`task.output` 在 CLI 上被真执行、那句 `风险未分级…` 从回执里消失；**没到手的那半它自己写清了**——"run B 重复同一枚 `task.output` 8 次后 FAIL，读数文件未落盘"，根因是**它自己台件**用正文文本找 call id，而 call id 在 `tool_call_id` 字段。
- ⚠ **预算硬顶 ≤35 次工具调用**；到顶即停手回禀并**自报枚数**。每做完一段 commit 一次；**只 commit、绝不 push**。

## 0. 起手五件

`date`｜`git rev-parse --abbrev-ref HEAD`（**必须 `dev`**；不是＝停手上报。全程禁 `switch`／`checkout <分支>`／`merge`／`rebase`／`reset`／`stash`／`worktree`／`clean`／`restore`）｜`git rev-parse HEAD`（我这边 `016a3b4a`）｜`git status --porcelain -- internal/ cmd/`（**必须空**）｜基线：`go test -count=1 ./internal/agent/`＋`./internal/tools/`＋**`./internal/risk/` 单跑**（四包并发会把 `TestResolvePerCallBudget` 打成假红，见 `A359`）。

## 1. 本单只做的一件事

造**你自己的**台件（⚠ **不改 `179-r1` 那枚 `probes/179/r1/zz179r1_e2e_test.go`**——它是它那一程的凭据，你动它＝两枚程共用一枚写面），落 `.scratch/wisp/probes/179/r2/**`，走**真 CLI 接缝**把这一发跑出来：

> 模型在同一进程里先调 `task.output` 拿回那条宿主自己写下的指针，再**照那条指针调 `fs.read`**，**逐字节读回那份副本**。

**取数写法上两枚已知的坑，先避开再跑**：
1. **call id 在 `tool_call_id` 字段里，不在正文里**（上一程就断在这）——按字段取，别拿字符串在 `text` 里找。
2. **`cmd/wisp` 要带 PATH 前缀才跑得动**（⚠ 这条推翻了我此前传给两枚程的说法，出处具名：`scripts/wisp-cli-tests.sh:20` 逐字记着 `windows, PATH=third_party/sherpa-onnx … PASS=33 FAIL=0 SKIP=0 -> GREEN`，CI 也挂在 `ci.yml:475`）。⇒ 起手先跑一次**只读**的对照：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 -run XXX_NONE_PKG ./cmd/wisp/` 应当 rc=0；**不带** PATH 时那枚 `0xc0000135` 仍在（＝票 98 那枚加载期坑，与产码无关）。**两形各贴一发读数**，别只报绿的那形。

**要交的读数**（落 `.scratch/wisp/probes/179/r2/logs/`，⚠ `logdir` 取脚本自身目录、不继承 CWD；⚠ 落盘后先 `wc -c` 证明非空再引用）：
- 续读那一发的**回执原文**（逐字，含 `isError`／层级／字节数）；
- **逐字节相等**那一断：产物字节数 vs `fs.read` 读回字节数，两个数都现量；
- 回执里**没有** `风险未分级…`、**没有** `…L2 级…已拒绝执行`、**没有** R4 命中、**没有**"查不到这个任务"——这四枚分岔逐枚给一条 grep 读数（命中 0 也要说明那把尺先用一枚已知存在的正控打过，否则"0 命中"不可信）。

## 2. 门禁（自己重跑，别引任何人的数）

`go test -count=1 ./internal/agent/`＋`./internal/tools/`＋`./internal/risk/`（单跑）｜`sh scripts/d22scan.sh`（应 rc=0；`ban #8 internal/` 我 10:4x 现量 **431**，你若没加 `.go` 就还是 431）｜`bash tools/d22scan/runtests.sh -C tools/d22scan ./...`｜`sh .scratch/wisp/probes/154/gate-clauses.sh`（⚠ **今天退码本就是 1**：唯一红腿 `G6neg 声明=ring 基线=1枚 实测=2枚`＝票 178 在册那枚；**你要报的是红腿名册有没有新增**，做法＝把 BAD 行 `grep -oE "腿=G[0-9a-z]+" | sort -u` 成名册再比，⚠ 不许为了让尺安静去改仪器或补 `OpenTask` 字面）｜`gofumpt --version` 现跑＋对你名下 `.go` 的 `-l`。⚠ **`probes/161/r6/flip-declaration.sh` 别跑**（跑一次就脏跟踪日志）。⚠ **最终那一次读数必须在最后一枚 commit 之后取**。

## 3. 写面（超出即越权）

`.scratch/wisp/probes/179/r2/**`（你的台件与读数）｜票 179 的 Progress log（**只追加、一枚框都不勾**）｜`docs/evidence/s1/179-cli-reread-r2.md`（你的表）。
**禁改**：`internal/**`（**产码与判据一字节都不许动**——本单预期零产码；若你判"非改不可"，停手上报）、`cmd/**`、`internal/panel/**`、`allowlist.txt`、`docs/PLAN.md`（含 `:1531`／`:1532` 两行 DEFERRED）、`docs/specs/**`、`docs/reports/**`、`frontend/**`、`design/**`（别家归属，不读不写不引）、`.scratch/wisp/probes/179/r1/**` 与别人票面／证据件（**只读**）。
⚠ 现场躺着 `design/**` 未提交删除、`.gitignore`／`probes/152/my152.py`／`probes/161/r6/logs/flip-*`（` M`）／`docs/evidence/s1/152-*.md`（` M`）⇒ **不提交、不还原、不删除、不评论**。
Git：`git add -A`／`git add .`／`commit -a` 一律禁；一步式 `git commit -q -F - -- <显式路径>`（新文件先 `git add <显式路径>`，再核 `git diff --cached --name-only`）；禁 `--amend`；heredoc 加引号（`<<'EOF'`）；含反引号的中文走编辑工具；串命令里可能无命中的 `grep` 用 `;` 或 `| cat`、**别用 `&&`**。临时件只建不删；还原跟踪文件用 `git cat-file blob HEAD:<路径> > <路径>`。

## 4. 交件报告（顺序固定）

1. step-0 五件；2. **本程没测什么**；3. **PATH 两形对照读数**（带／不带各一发）；4. **续读那一发的回执原文＋逐字节相等那两个数**；5. **四枚分岔的 0 命中读数**（每枚带正控说明）；6. 你有没有动过任何产码（应当是**没有**；若动了＝本单作废、停手上报）；7. 门禁全部读数（⚠ 最后一枚 commit 之后）＋红腿名册比对；8. 被拒／没成功的调用（取数前还是后）；9. 有没有跑过删除命令；10. 工具调用**用了几次 vs 硬顶 35**；11. 伪授权两栏（各带出处）；12. 凭据值零抄录；13. `next=`。

⚠ 每个"几枚／几个／几行"旁边附可复制命令；先测→再写→再提交。
⚠ **我这段话里每条前提（含"断点只在台件写法"、"零产码就够"、"PATH 前缀能解 `cmd/wisp`"）都是未验证断言**：不符就报回并继续做做得动的部分，**不许为了对我那句话去改判据、改测试或改产码**。
