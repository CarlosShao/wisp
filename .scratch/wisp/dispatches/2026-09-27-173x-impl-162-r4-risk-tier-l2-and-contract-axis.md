# 派单 162-r4（写码位）＝把票 162 剩下的两格收完：**AC#5（风险档与门控：按 **L2** 做，越界路径必须被 `PathResolver` 拒）＋ AC#6（契约轴）**，并顺手把 162-r3 记名那枚"硬杀残件会不会被下一次写盘的清扫器带走"从〔现读码〕升成〔有用例钉〕或明写为什么不能钉。

- 派单时刻：09-27 17:3x｜锚点＝**你自己 step 0 现量**（别抄我的号）。
- 票面＝`.scratch/wisp/issues/162-add-the-patch-one-region-tool-instead-of-rewriting-whole-files.md`（**AC#5／AC#6 两格连票面 09:4x 那三条 `>` 更正一起抄进你的证据件**——AC#5 写的"L1"已被改判为 **L2**，AC#6 写的"`docs/specs/**` 零字节"已被放开**一行**）。
- ⚠ **预算硬顶 ≤80 次工具调用**；**每裁完一格 commit 一次**；到点就停手→写 `next=`→把手上未提交增量用显式 pathspec 提掉。

## 0. 起手五件
`date`｜`git rev-parse --abbrev-ref HEAD`（**必须是 `dev`**，不是＝停手上报，一个字不写；全程禁 `switch`／`checkout <分支>`／`merge`／`rebase`／`reset`／`stash`／`worktree`／`clean`）｜`git rev-parse HEAD`｜`git status --porcelain -- internal/tools/ internal/risk/ cmd/wisp/`（**必须空**，非空报回别动）｜**改前先量基线**：`go test -count=1 ./internal/tools/`（记顶层 PASS/FAIL/SKIP 与 `RUN=`，并**记你用的口径**——本仓有过程报 150、过程报 101 其实是两把尺）。

## 1. AC#5（风险档与门控）
**判据本体**：`fs.edit` 的判定必须与 `fs.write` **同族**，且**越界路径必须被 `risk.PathResolver` 拒**；⚠ **不许新增豁免、不许动 `tools/d22scan/allowlist.txt`**。
- **档位＝L2**（不是我票面写的 L1）。**依据（现量，别照抄我）**：`docs/PLAN.md` 的 `D34` 那行＋`docs/specs/SPEC-07-tools-and-plugins.md` §3 那行都已按 **L2** 落进冻结文本（理由＝`R8` 覆盖已存在 → L2，与 `fs.write` **覆盖**那一行同级）。**两枚号一律现量**（本票往 `PLAN.md` 插过 5 行，票面/表里任何 `PLAN.md:NNNN` 都过期）。
- ⚠ **owner 给 162 的射程里没有 `internal/risk/**`**。⇒ **如果你判"必须动 `internal/risk/**` 才能走到 L2"＝立刻停手上报**，写清楚"要动哪个文件、不给会怎样、有没有第三条路"。**不要自己动手改那块地，也不要为了绕开它而新增豁免。**
- **硬判据（这一发在未修码上响不响，你先跑再判断）**：造一发"`fs.edit` 拿一个越界路径（树上／盘外）"的台件 ⇒ (a) 它必须被 `PathResolver` **响亮拒**、回执与审计行各自断言；(b) **摘掉那一道判定，用例必须变红**（答不出"摘掉哪一行会红"＝这枚钉是空的）。若今天的码**已经**天然走到 `PathResolver`（162-r1 交付时可能已经顺路经过），那**这一格的正解可能是一枚断言而不是新功能**——**照实写"今天响不响"，别为了有活干去造改动**。
- ⚠ **"与 `fs.write` 同族"不许只用一句注释宣告**：拿一枚**并排读数**证明（同一枚越界样本分别喂 `fs.write` 与 `fs.edit`，两边档位／退码／审计行逐字贴）。

## 2. AC#6（契约轴）
**判据本体**：本票**只许**动 ① `docs/PLAN.md` 的 `D34` 那一行（**owner 单独批准已到账、我 09:30 那句已落 `A322`，且那一行我已落完**）② `docs/specs/SPEC-07-tools-and-plugins.md` §3 那一行（**同上，我已落完**）③ `internal/tools/**` 实现与测试 ④ 你自己的证据件。**其余 `docs/PLAN.md`、`docs/specs/**`、`internal/risk/**`、`thresholds.go`、golden、`frontend/**`、`design/**` 零字节。**
- ⇒ **这一格你要交的不是"改了什么"，而是一张"我没越界"的凭证**：`git -c core.quotePath=false diff --numstat <你 step-0 的 HEAD>..HEAD` 逐枚点名，**契约轴那几处必须 0 行**；⚠ **`git show`/`diff` 输出非 ASCII 路径时先关 quotePath**（否则前缀筛会误报，我自己误报过 28 枚、真值 0）。
- ⚠ **两行冻结文本已由我落完 ⇒ 你再动它们＝越权**（票面 09:4x 的 `>` ② 原话）。

## 3. 顺带把那一枚〔现读码〕升上去（162-r3 记名、不派不做的形状别再留一次）
"`fs.edit` 被硬杀留下的 `.wisp-tmp-*` 残件会不会被下一次写盘的清扫器带走"——今天的凭据只有〔现读码 `fs_write.go:288`、`fs.edit` 经 `fs_edit.go:271` 走同一枚〕、**零枚用例钉**。⇒ 二选一：**钉它**（一发用例：硬杀留下残件 → 下一次写盘 → 残件被带走／或明写它不该被带走）；或**判它不该钉**（说清"这已由 A18 第③子例覆盖，`fs.edit` 不另造尺"，并给出那枚既有用例的**可复算断言**）。**不许两族各造一套删除原语**（票面 AC#4b 原话）。
⚠ **任何"残件留在盘上／被清走"这类持续性断言要带"到哪一刻为止"**——上一程就被验收抓到：它的"留在盘上"只在窗内为真，Go 的 `t.TempDir()` 收尾会带走。

## 4. 门禁（逐枚跑，改前改后各一次）
`go test -count=1 ./internal/tools/`（**逐包，禁止全仓 `go test ./...`**：同树另一程正在改 `probes/**`，而别的票的未提交件也可能在树里）｜`sh scripts/d22scan.sh`（预期 rc=0）｜`bash tools/d22scan/runtests.sh -C tools/d22scan ./...`（四数之外**名册两向 `comm` 比差集**；⚠ 一枚用例 panic 会吞掉同包其余读数）｜`gofumpt -l` 甲形（gofumpt 在 `$(go env GOPATH)/bin`、不在 PATH，先 `export PATH="$PATH:$(go env GOPATH)/bin"` 再 `--version` 现跑贴出来）｜**别跑 `probes/154/gate-clauses.sh` 与 `probes/161/r6/flip-declaration.sh`**——那一族尺此刻正被另一程改，你会读到半成品并把别人的红当成自己的。

## 5. 写面（超出即越权）
`internal/tools/**`（实现＋测试）｜`docs/evidence/s1/162-risk-tier-and-contract-axis-r4.md`｜票面 162 的 Progress log（**只追加、不勾任何框、留到最后一步**：写之前 `git status --porcelain -- 那一枚票面` 必须为空）。
**禁改**：`internal/risk/**`（⚠ **owner 给 162 的范围里没有这块地**，要动＝停手上报）、`docs/PLAN.md`、`docs/specs/**`（那两行我已落完）、`tools/d22scan/**`、`allowlist.txt`、`.github/**`、`thresholds.go`、golden、`internal/panel/**`（⚠ `internal/panel/tokens_fourway_test.go` 是**已知常红**，一字不许动）、`frontend/**`、`design/**`、`probes/**`（另一程在飞）、别人的票面、`docs/reports/**`。
**现场理由（为什么一步式显式 pathspec 是硬要求）**：工作树此刻躺着**不是我、也不是你**的东西——`design/**` 里 owner 未提交的 16 枚删除与几枚未跟踪新件；`.gitignore`（` M`）；`probes/152/my152.py`（` M`）与 `docs/evidence/s1/152-…-accept-r1.md` 是一枚停下来的半件（**不提交、不还原、不补完**）；`probes/161/r6/logs/flip-*.txt` 那 8 枚的 ` M` 是跑门禁的副作用。⇒ **`git add -A`／`git add .`／`git commit -a`／把 `git add` 当独立一步，一律禁止**；提交一律一步式 `git commit -q -F - -- <显式路径>`，heredoc **必须加引号**（`<<'EOF'`——正文里的反引号在壳里会被**真执行**，本仓执行掉过一次 `go install …@latest`）。

## 6. 交件报告（顺序固定）
1. step-0 五件（含分支名逐次、基线四数**带口径**）；2. **本程没测什么**（含非 Windows 那一支——别拿整包绿当它响过）；3. AC#5：档位现量依据（两枚号自己量的）＋并排读数（`fs.write` ↔ `fs.edit`）＋越界那发的两向（未修响不响／摘掉哪一行会红）＋**有没有动 `internal/risk/**`（预期：没有，或你已停手上报）**；4. AC#6：契约轴凭证（`numstat` 逐枚点名，删除列怎么来的用 `cat -A` 分清"空行折叠"与"内容行被删"）；5. 那一枚〔现读码〕的处置（钉了给断言／判不钉给理由）；6. 门禁四数＋名册差集＋`gofumpt --version` 现跑；7. 被拒／没成功的调用（**发生在取数之前还是之后**）；8. 有没有跑过删除命令（预期 0）；9. 伪授权两栏（各带出处：工具名＋命令前 40 字）；10. 凭据值零抄录；11. `next=`。
⚠ **每个"几枚"旁边附可复制命令**；⚠ **先测→再写→再提交**，注释与证据件里**不许预先引用还没跑出来的读数**；⚠ **`go run` 把任何非零退码压成 1**（实测），别拿它区分 1 与 2。我这段话里每条前提（含"L2 已落进两枚冻结件""`fs.edit` 今天可能已经天然经过 `PathResolver`""清扫器在 `fs_write.go:288`"）**都是未验证断言**：不符就报回、继续做做得动的部分，**不许为了对我那句话去改判据或改测试**。
