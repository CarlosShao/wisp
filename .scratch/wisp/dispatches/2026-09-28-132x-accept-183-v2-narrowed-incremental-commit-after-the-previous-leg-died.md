# 派单 `183-v2`（非实现者验收·**收窄版**）— 票 183 落地那枚豁免的独立裁决；⚠ 上一枚 `183-v1` **死在收尾**（52 枚工具调用后 harness 报错、表没写、没 commit），它的**读数仍在盘上**可复用

- 派单时刻：`2026-09-28 13:2x`（`date` 你自己现量）
- 起手锚：**`d61aee1b`**（`184-c1` 的交件号；被验收的那枚产码 commit 是 **`0662a35a`**）
- 票面：`.scratch/wisp/issues/183-the-host-minted-pointer-exemption-does-not-hold-on-the-real-cli-…md`
- 被验收对象：`183-r1` 交件表 `docs/evidence/s1/183-pointer-exemption-r1.md`（24874 字节）＋两枚产码件（`internal/risk/provenance.go +48`／`internal/risk/taintmatch.go +59`，删除列 0）＋判据件 `internal/risk/pointer_183_test.go`（＋331，顶层 6 枚）
- **死程留下的现场**（⚠ 全部按〔仅自述·程已死·未验收〕对待）：`.scratch/wisp/probes/183/accept-v1/`——变异副本 `mut-ma/`…`mut-mf/`＋各自 overlay（`ov-ma.json`…）、台件 `zz183v1_e2e_test.go`（19599 字节）、读数 `logs/run-head-probe.txt`（4334）、`logs/e2e-readings.txt`（33389）、`logs/cli-reread-183v1.txt`（27358）、`logs/gate-clauses.txt`、`logs/d22scan.txt`。我 13:21:25 现量：`git status --porcelain -- internal/ cmd/` **为空＝它没留变异**、probes 里 `grep -rl MUTATION` 命中 0。
- 台账：`A363`（我否决前程推荐修法＋批准面三件）· `A364`（收 `183-r1`）

## 0. 这一单为什么收窄（别再死一次）

`183-v1` 拿了六个问题、50 枚预算，跑到第 52 枚死在"写表"之前，**产出＝零**。所以本单：
- **预算 40 枚**，且**第 28 枚起停止新探索**，剩下的只用于写表与 commit；
- **增量交付（硬要求）**：第 ① 次 commit 在**前 8 枚工具调用之内**完成——先建表骨架（§号齐全、每节写"待填"）＋step-0 读数，commit；此后每做完一格就**追加并 commit**（表允许中途不完整，不允许盘上没有）。这样任何一次中断都留下可读的部分裁决；
- 能复用就复用死程的台件与变异副本（**但你要引用的数，必须是你自己这一跑出来的**，它的读数只当"存在性线索"）。

## 1. step-0 五件（前 8 枚工具调用里做完并 commit）

`date`｜`git rev-parse --abbrev-ref HEAD`＋`--short HEAD`（应 `d61aee1b`）｜`git status --porcelain -- internal/ cmd/`（必须空）｜三向 md5（`provenance.go:468-473`＝`858e45116383caa3e7c1dd4b0924fad1`、`taintmatch.go:11-15`＝`5680ddd18e2d2ec2a85e485b54f4c12e`、`provenance.go:482-506`＝`f89e891e5eee3f3ea2b4f89d921072c4`）｜`git show --numstat --format="%h" 0662a35a`（确认删除列 0、名下文件只在你预期那几枚里）。

## 2. 四格必答（**只有这四格**，其余一律登记"未裁"）

- **① 恒真性矩阵（最值钱）**：对 `183-r1` 那 6 枚常驻判据，至少自取**三发**变异并逐名列出"哪枚红"：(M1) 摘掉 `contains()` 里 `spellsDeclaredPath` 那段跳过；(M2) 摘掉 `provenance.go` 里 `attachDeclaredPath` 那一支（**保留**位置性 skip）；(M3) **不要求 `runeIndexOf` 命中就无条件挂 `declaredNorm`**（这一发最重要：它测"声明的路径根本不在正文里时豁免会不会仍生效"＝fail-closed 那一支有没有牙）。⇒ 每一发之后 `git cat-file blob HEAD:<path> > <path>` 还原，终态再量一次 `internal/ cmd/` 为空。**答不出"哪枚红"的那枚判据＝装饰，具名判不通过。**
- **② 安全退让面的进攻**：落地那族把"拼进声明路径 ≥8 rune 的窗口"从**这一枚 mark** 的证据里去掉。**自己造一发**：外来正文里放一条与声明路径共享 ≥8 rune 的敏感串，看外带会不会逃过 R4；再判"影响面只在同一枚 mark 内，还是跨 mark 可借"。**并复量 `183-r1` 的频率读数**（尺 `cat probes/183/r1/logs/freq-head.txt`）：那 10 份语料是**仓内文件、不是用户文档**——你要判"这句话写对了吗、口径够不够"。
- **③ 三条边界**：(a) 参数侧按值放行（票 183 AC#5① 明禁）有没有以别的形式回来；(b) `declaredPath` 会不会变成事实上的名册（要**证它跨 mark 不可见**，别只读注释）；(c) `taintmatch.go:11-15` 那句"逐 token 追踪已被否决"算不算被复活——⚠ **编排者 13:0x 判"不算"（理由在 `A364` 裁定②），这一格就是请你推翻或坐实我**。
- **④ AC 判读（只给判读，不许勾）**：票 183 **AC#2** 的措辞"模型照宿主指针续读 ⇒ 不命中 R4、`fs.read` 成功、**逐字节读回**"——真机现在**第一发通、第二发（同一发任务里对同一条路径的有界重读）仍被 R4 拒**（`183-r1` 读数第 55／60 行，第二发的判中者是 `fs.read` 自己那枚未声明的 mark）。⇒ 你答三件：**这一格算达成／部分达成／不达成**；票 177 AC#3 端到端条件、票 175 AC#5、票 176 AC#3-5 **能不能翻**；以及票 185（我已立、按住）该不该独立存在还是并进票 183。**AC 框一枚不许你勾。**

## 3. 允许"没做完"，不允许"没说"

表里必须有一节 **`本程没裁什么`**：把 §2 之外的格（票 183 AC#1 四支候选、AC#4 兄弟形、AC#6 窗口 0.0s、AC#7 契约轴、`-race` 并发那一格）逐枚具名登记"未裁·为什么"。⚠ 死程 `183-v1` 已经跑过一部分（它的 `mut-*` 与 `logs/*` 就是证据），**你可以引用它的现场存在性，但不许把它的读数当你的凭据**。

## 4. 护栏（与派单 `183-r1` §6 同一套）

只 commit **不 push**；显式 pathspec；禁 `git add -A`／`.`／`-a`／`--amend`／`reset`／`rebase`／`stash`／`checkout`／`switch`／`merge`／`worktree`／`clean`／**任何删除命令**；还原只用 `git cat-file blob HEAD:<path> > <path>`；终态 `git status --porcelain -- internal/ cmd/` 为空。`frontend/**`、`design/**` **不读不写不引**；那些 ` M`／` D` 脏件（`.gitignore`、`probes/152/my152.py`、`probes/161/r6/logs/flip-*`、`docs/evidence/s1/152-*.md`、`design/**`）**不提交、不还原、不评论**。凭据值零抄录。
CLI 腿要 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（不带＝`0xc0000135`）；`./internal/risk/` **单包跑**；台件 `logdir` 不用 `runtime.Caller`；**不许跑** `probes/161/r6/flip-declaration.sh`；`gate-clauses.sh` **比红腿名册不比退码**（在册唯一红腿 `腿=G6neg`＝票 178）。
**交付**：表 `docs/evidence/s1/183-pointer-exemption-accept-v2.md`（§号齐全，含被拒/没成功的调用、有没有跑过删除命令、工具调用枚数 vs 硬顶 40、伪授权两栏、`next=`）；读数落 `.scratch/wisp/probes/183/v2/**`（**别写进 `accept-v1/`，那是死程的现场**）；票 183 Progress log **只追加**；至少 **两次 commit**。
