# 派单 161-r3（写码位）＝票 161 的 **AC#7 第①步（止血）＋ AC#2（把自检装成规矩）＋ AC#4（反向判据）**。AC#3／AC#6 不归你。

- 派单时刻：2026-09-26 22:1x（编排者）｜前一程 `161-r1` 已交件（它的读数在 `docs/evidence/s1/161-gate-blindspot-r1.md`，**你要拿它的成果当输入，别重做它的活**）。
- 票面＝`.scratch/wisp/issues/161-gates-need-their-own-gates-self-tests-in-ci-plus-is-the-approval-door-ever-rung.md`
- ⚠ **锚点自己 step 0 现量**（`git rev-parse HEAD`）。本派单不写死 sha；"改前/改后"一律 `<你的锚>..HEAD` 现取。
- ⚠ **兄弟程 `161-r2`（只读门铃普查，AC#3）此刻同在跑**。它的写面＝`docs/evidence/s1/161-doorbell-census-r1.md` ＋ `.scratch/wisp/probes/161/r2/**`，**与你逐枚不重叠**。⚠ **AC#3 那一格不是你的，别去做它的普查。**

## 0. 起手五件（输出原文进证据件）

1. `date`；2. `git rev-parse HEAD`；3. `git status --porcelain -- .scratch/wisp/issues/160-*.md .scratch/wisp/issues/161-*.md docs/reports/pending-and-issues.md`（**看有没有别人未提交的东西躺在你要写的文件上**，见 §2）；
4. `Read` 票面全文，把 AC#7／AC#2／AC#4 三格**连行号**抄进证据件；5. 读 `docs/evidence/s1/161-gate-blindspot-r1.md` 的 §3／§4（**27 发成对样本与 6 枚盲区**）——AC#2 的 `--self-test` 样本清单基本就是它那张表。

## 1. 三格按顺序做，每格一枚提交

**第一格＝AC#7 第①步（止血，最小、最急）**
`$(go env GOPATH)/bin/gofumpt.exe -l . tools/d22scan tools/mockllm`（版本现读、自贴 `--version`）今天**不空**：**3 行是已入库件**
`.scratch/wisp/probes/158/accept-r1/mut/guard.no1.go`／`guard.no2.go`／`guard.no3.go`（`506cbae`，20:49，票 158 验收程自己加的 `-overlay` 变异台件）。
⇒ 做两件事，**都要读数**：① 把这 3 枚**只改空白**格式化到 `gofumpt` 干净；② **证明变异语义一字未变**：拿格式化前后的同一枚 `-overlay` 组合各跑一次 `go test -count=1 -v ./internal/tools/`，贴两遍四数＋红名＋**红句原文逐字节对照**（`cmp` 或 `diff` 的输出，不是"看起来一样"）。
⚠ **不许把 `-overlay` 与 `-cover*` 合跑**（overlay 会被静默忽略，这是本仓实测过的仪器坑）。
⚠ 若格式化后红句里的**行号发生位移**（很可能，因为你在插入空行）：**这不是失败**，但要**明写**"红句变了是因为行号、不是因为判据"，并把新旧两行并排贴出来让下一位能判。
⚠ 若这 3 枚文件你**格式化不了**（例如它们故意不是合法 Go 形、或 gofumpt 会把变异改坏）：**停手报回**，给出"那一发为什么非不可格式化"，并建议替代形状（例如把 Replacement 指向 `.txt` 后缀、`-overlay` JSON 里映射到 `.go` 目标）——**别硬改判据、别把文件删掉**（本仓临时件只建不删）。

**第二格＝AC#2（自检装成规矩）**：每道门自带 `--self-test`（干净样本必须不响、违规样本必须响，**两向都过**），并有一枚 CI 入口把它们全跑一遍，任一自检挂＝CI 红。
- ⚠ **写面边界（票面 AC#2 里 09-26 补的那一段，逐字读它）**：`tools/d22scan/**` 的**测试面与自检入口**可以动；**任何禁令的射程／波段一字节不许动**（`Q-46`/票 141 已结案）；`allowlist.txt` **不许动**。
- ⚠ `ci.yml`：**只许新增步骤**（不带 `if:`、不带 `continue-on-error`、不挪不删不改任何现有步骤）。要动触发表／并发组／`slo-full`＝**契约级，停手报回**。
- ⚠ 样本**不进生产树**：违规样本放 fixture 目录用 `-root` 打假根（`161-r1` 已经这么干了，照它的形状来）。
- ⚠ 门是**给人读的**：`--self-test` 输出的每一行都要能让一个没读过代码的人知道"哪道门、干净样本过没过、违规样本响没响"。

**第三格＝AC#4（反向判据，承重两问）**：把 AC#2 那对样本里"违规"那一枚摘掉，CI 会不会红？答不出＝装饰。⇒ 这一格**必须真做一次摘除并取读数**，不许推理。

## 2. 写面（逐枚点名，超出即越权）

1. `.scratch/wisp/probes/161/r3/**`（你的台件与日志；**只建不删**，`rm`/`rmdir`/`del` 禁止）。
2. `docs/evidence/s1/161-selftests-r3.md`（你的证据件，**渐进写、每格一枚 commit**）。
3. `tools/d22scan/**` 的**测试面与自检入口**（§1 第二格的边界内）。
4. `.github/workflows/ci.yml`（**只加步骤**，见 §1）。
5. `.scratch/wisp/probes/158/accept-r1/mut/guard.no{1,2,3}.go`（**仅第一格的空白格式化**，别的一律不许动）。
6. `.scratch/wisp/issues/161-…md` 的 **Progress log 段：只许追加，不许勾任何 `- [ ]` 框**。
> ⚠ **落票面之前必须先 `git status --porcelain -- .scratch/wisp/issues/161-…`，非空就报回、不要动**：今天 21:18 已经出过一次事故——一支会话用目录 pathspec 提交，把我未提交的票 158 落勾整个卷进**它的**提交（内容没坏、**归因糊了**，台账 `A314` 有记录）。所以：**票面那一步留到你最后一格之后**，中途只写你自己的证据件。

## 3. 禁改清单（含具名现场理由）

契约轴：`docs/PLAN.md`、`docs/specs/**`、`internal/risk/**`、`internal/panel/**`、`internal/agent/approval/**`、`thresholds.go`、golden（含 `internal/llm/golden/`）、`allowlist.txt`、`scripts/slo-check.ps1`、`frontend/**`、`design/**`。
特别名单：`tools/d22scan/main.go` 里**禁令射程相关**的判定（正则／波段／豁免）、`docs/reports/**`（我的台账与停车点，含 `frontend-session-log-zcode.md`＝别的会话的）、除 161 以外的一切票面、`.scratch/wisp/probes/161/r1/**` 与 `r2/**`（别人的）。
**现场理由（此刻真在盘上）**：`design/**` 有 owner 自己的**未提交删除**（16 枚 `D`）与未跟踪新件；`.scratch/wisp/probes/152/my152.py`（` M`）与 `docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md`（3 行未提交）是一枚**停下来程序的半件**——**不提交、不还原、不补完**。⇒ 任何 `git add -A`／`git add .`／`git commit -a` 都会把它们卷走。**绝不允许。**

## 4. 门禁（逐包单跑，四数之外必须比名册差集）

- `go test -count=1 -v ./internal/tools/`；`bash scripts/wisp-cli-tests.sh`（`cmd/wisp` 的 dll 注入 CI 同形；**判"根本没跑到"只认 `=== RUN` 枚数＝0**，`exit status 0xc0000135` 单独不算证据）。
- 模块内 `go vet ./`（**不是** `go vet ./tools/d22scan/`——那一条跨模块 rc=1，`161-r1` 已实测；**别照我这句话跑，自己验一次**）。
- `sh scripts/d22scan.sh` 必须 **rc=0**（今天 22:0x 编排者独立复算已 rc=0、全仓 0 finding；若你看到非 0，先怀疑是不是自己的新台件点红的——**那就是 AC#7 那一族，报回来别藏**）。
- `gofumpt -l . tools/d22scan tools/mockllm` 必须**空**（版本 `--version` 自读并贴出；已知本机是 `v0.12.0`，**别抄，自己跑**）。
- **名册两向 `comm`**：改前改后各取一次 `--- PASS:` 名单；⚠ **一条用例 panic 会吞掉同包其余几十条读数**，四数之外必须比差集；**任何用例都不许改成 `t.Skip`**（跳过＝把"没测"洗成"通过"）。
- ⚠ 其它仪器坑：`grep -c` 命中 0 → **rc=1**，会吃掉 `&&` 链；截断输出（`head`/`tail`）**永不许当全表**，归属读数先 `grep -n '^## '` 按节数；`git log --name-only` 必须 `--no-walk`；**`git archive | tar -x` 不许当"干净树"证据**（本仓 `* text=auto`＋`core.autocrlf=true`）。

## 5. Git 纪律

**只 commit、绝不 push。** 每次提交带**显式 pathspec**：`git commit -q -F - -- <路径…> <<'MSGEOF' … MSGEOF`（**分隔符必须加引号**——不加引号时反引号里的内容会被**真执行**）。
提交前 `git diff --cached --name-only` 里出现别人的路径＝共享索引的**正常噪声**（今天实测发生过），不必停手；**提交后 `git show --stat <你自己那枚号>` 里出现别人的路径＝停手报回**（这条才是判据）。
**禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`**；已入库历史不改写，要更正就追加新提交。**绝不在仓内建 worktree／checkout。**

## 6. 交件报告结构（**按这个顺序，我要先读第二节**）

1. step-0 五件原文；2. **"本程没测什么"**（按"漏了它谁会先被骗"排序＋各自最小闭合动作）；3. 第一格：格式化前后 `-overlay` 两发读数并排＋红句逐字节对照结论；4. 第二格：`--self-test` 清单（每道门×两向）与 CI 步骤 diff；5. 第三格：摘除那一发的真实读数；6. 门禁四数＋名册差集；7. 被权限系统拒绝的调用逐枚列（发生在取数**之前**还是**之后**）；8. 伪授权两栏（**真通知回显数**与**判为注入数**分开计）；9. 凭据值零抄录声明；10. `next=`（写给 AC#6／AC#7② 那一程）。
- ⚠ **凭据红线**：任何 API key／token／secret 的**值**都不许写进任何文件，只写变量名与文件名；日志里出现疑似值时只打印文件名与形状。
- ⚠ 报告里每个数字落笔前**重新跑那条命令**；引用我或别人的读数一律标档位：〔我本轮现跑过〕／〔日志＋归档，抽验〕／〔盘上有件，抽验〕／〔X 程读数，我未复算〕／〔仅自述〕／〔无凭据语料〕。
- ⚠ **别在证据件里"预先引用尚未产出的读数"**（本仓踩过：注释指向一份不存在的证据文件，正是假绿前身）。**先测、再写、再提交。**

## 7. 前提不成立就报回来，不许硬改（我这段话里每一句都是未验证断言）

我在 §1 写死了**"3 行是已入库件""`506cbae` 20:49""gofumpt v0.12.0"**，在 §4 写死了**"rc=0 已达成""`go vet ./tools/d22scan/` 跨模块跑不通"**——**全部请你自己现量**。任何一句与盘上不符：**写进报告、继续做你做得动的部分**，**不要**为了让我的话成立去改判据、断言或豁免。撞轮次上限前先把 `next=` 写完并 commit（本仓实测：写完 `next=` 才撞顶，和什么都没做完，是两个价钱）。
