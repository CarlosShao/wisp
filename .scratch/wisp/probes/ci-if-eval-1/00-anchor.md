# 00-anchor｜只读普查腿 `ci-if-eval-1`｜起手读数

现量时刻：`date '+%Y-%m-%d %H:%M %z'` = **2026-10-08 15:47 +0800**（起手另有一发 14:48 +0800，同一程内）

## 1. 锚

| 件 | 现量 |
|---|---|
| 分支 | `dev` |
| 本机 HEAD | `6547fd3046066d0f9e2fde9ba2cad3ee384573e8` |
| 远端 `origin/dev`（GitHub） | `cc31526165734e612de297848bb2080bd459ccba`（`gh api repos/CarlosShao/wisp/branches/dev` 与 `.../git/ref/heads/dev` 两端点同值） |
| `e6dc79ed` 全号 | `e6dc79ed5bd14380cce1507f94886d65a51ee63f`，committer date `2026-10-08 12:06:33 +0800`（= 04:06:33Z） |
| `e6dc79ed` 与本程关系 | `git merge-base --is-ancestor e6dc79ed HEAD` ⇒ **真**；父＝`05992d0520cb2de1abfad0e6ade0c9360f3bf5b4` |
| HEAD 上的 ci.yml 是否等于 `e6dc79ed` 上的 | **等于**：`git log --oneline e6dc79ed..HEAD -- .github/workflows/ci.yml` 回空（rc=0） |
| 工作树的 ci.yml 是否脏 | **不脏**：`git status --porcelain -- .github/workflows/ci.yml` 回空（rc=0）⇒ 本件全部读数取自查 `HEAD` 快照，与工作树一致 |
| 快照尺寸 | `git show HEAD:.github/workflows/ci.yml` = **1015 行 / 62,192 字节**；本腿取到仓外 `D:/tmp/ciifeval/ci-head.yml`（⛔ 不在仓内落临时件、⛔ 不建脚本） |
| 另一枚 workflow | `git ls-tree --name-only HEAD .github/workflows/` = `ci.yml` + `slo-fresh.yml`（后者 3 枚步、**步级 `if:` 0 枚**；不在本格射程，只具名存在） |

## 2. 尺的写法（⛔ 不是词频尺）

两把**独立实现**的结构尺，同数才落笔：

**尺 A（缩进＋键形状）**：一枚 `if:` 计入"步级"当且仅当
ⓐ 该行 `indent == （该 job 内第一枚 `- ` 列表项的 indent）+2`（实测 ci.yml 里＝第 8 列），
ⓑ `stripped` 以 `if:` 起手，
ⓒ 该行**不在 block scalar 内**——先扫一遍把 `key: |` / `key: >`（含 `|-` `+` 等修饰）起、以 indent 大于该键的续行为止的行全部打掩码，所以 `run:` 里的 shell/Python 源码**不可能被当成语法键**，
ⓓ 该行不是注释。
"作业级"同理：`indent == 4` 且位于 `jobs:`（第 53 行）之下、键名为 `if`。

**尺 B（逐步状态机）**：按 `- ` 列表项切步，在每一项内扫到第一个 `if:` 键即记该步的守卫行号；同样带 block-scalar 掩码。

两把尺在 HEAD 上同数（36／52／16）。⚠ 尺 B 自己错过两回、都留在账上：第一回把步键缩进写成"列表项缩进"而非"＋2" ⇒ 报 **0 枚**；修掉后又因步内扫描起点差一行（`range(li+1,…)` 跳过了紧邻 `- name:` 的那一行），把 14 枚"守卫紧贴步名行"的步整枚漏掉 ⇒ 报 **22 枚**；两处都修后回到 **36**，与尺 A 同数。⇒ 本件里凡"枚数"格都由**两把独立实现＋逐 commit 加法**（§3 末行）三处对拉，⛔ 没有任何一格是单把尺说的。

**词频尺的形状（只为对照，⛔ 不用它下结论）**：
- `6547fd30:.github/workflows/ci.yml` 里含 `cancelled`（大小写不敏感）的**行数** = **50**；
- 其中**注释行** = **14**；剩下的 36 行恰＝尺 A 的步级 `if:` 数；
- 所以词频尺与结构尺在**本 ref**差 14 枚，全差在注释上。⚠ 这与 `A710 §5` 记的"25＝含注释的行数／27＝改前行数／49＝改后行数／50＝出现次数"**不是同一把尺也不是同一时刻**，本腿只报本 ref 上的数，⛔ 不去对齐那三四个数。

## 3. 起手结论（详表在 `report.md`）

| 格 | 现量 |
|---|---|
| 步级 `if:` 枚数（HEAD） | **36**（`${{ !cancelled() }}` **35** ＋ `always()` **1**（`6547fd30:.github/workflows/ci.yml:473` `Stop compose services`）） |
| 作业级 `if:` 枚数（HEAD） | **0** ⇒ **复认编排者此前那句**，不推翻。两把尺同数，另有硬旁证：全档"去空格后以 `if:` 起手"的行＝**36**，与步级 36 枚**一一抵尽**⇒ 作业级不可能还有第 37 枚 |
| 步总数／有守卫／无守卫 | **52 / 36 / 16**（lint 13、test-core 7、test-windows 10、slo-smoke 6、slo-full 5、lint-frontend 11） |
| `e6dc79ed` 到底加了几枚 | **23**（尺：`git show e6dc79ed -- .github/workflows/ci.yml` 中 `^+` 紧跟 8 空格再跟 `if: ` 的行 = **23**；`^-` 同形 = **0**；逐枚内容全为 `if: ${{ !cancelled() }}`）⇒ **复认**台账 `A710 §1/§2` 那句"23 枚、0 删" |
| 按 job 拆那 23 枚 | lint **7** ＋ test-core **4** ＋ slo-smoke **3** ＋ lint-frontend **9** ＋ test-windows **0** ＋ slo-full **0** = **23** ⇒ **复认** `e6dc79ed` 提交标题里那句拆分（本腿不是照抄标题，是用 `cc31526165`／`e6dc79ed^`／HEAD 三版名册做差集推出的） |
| 步级 `if:` 的三代枚数 | `cc31526165`（GitHub 上最新一版 ci.yml）**10** → `e6dc79ed^`（`05992d05`）**13** → HEAD **36** ⇒ `A710 §2` 那句"13→36"复认；**10→13 这一段是编排者没列的第 3 枚**，见 `report.md` §2 |
| ci.yml 在 `cc31526165..HEAD` 的形状 | `git diff --stat` = **160 insertions(+), 0 deletions** ⇒ 零删除＝**没有一枚既有步被改名／挪位／删掉** ⇒ 用步名去对齐历史 run 是**安全**的（这条是本腿最能撑住对齐结论的一枚） |

## 4. 命令与退码（取数全档；⛔ 零 go 命令、零 push、零触发 run）

| # | 命令（原样） | rc | 落盘 |
|---|---|---|---|
| 1 | `git rev-parse HEAD`／`git log --oneline -8`／`git status --porcelain \| head -40` | 0 | 直接读 |
| 2 | `git merge-base --is-ancestor e6dc79ed HEAD` | 0（真） | — |
| 3 | `git status --porcelain -- .github/workflows/ci.yml` | 0（回空＝不脏） | — |
| 4 | `git log --oneline e6dc79ed..HEAD -- .github/workflows/ci.yml` | 0（回空） | — |
| 5 | `git show HEAD:.github/workflows/ci.yml` | 0 | `D:/tmp/ciifeval/ci-head.yml` |
| 6 | 尺 A python（stdin heredoc，不落仓内脚本） | 首次 **rc=1**：`/tmp` 在 Git Bash 与原生 Windows python 解析到不同目录 ⇒ `FileNotFoundError`；**改 `D:/tmp/ciifeval` 后 rc=0** | `ruler.txt`／`ruler2.txt` |
| 7 | `git show e6dc79ed -- .github/workflows/ci.yml` | 0 | `diff-e6dc79ed.txt`（161 行） |
| 8 | `git show cc31526165:.github/workflows/ci.yml`／`git show e6dc79ed^:...` | 0 / 0 | `ci-runrev.yml`(855)／`ci-parent.yml`(992) |
| 9 | `gh repo view --json nameWithOwner,defaultBranchRef` | 0 | `CarlosShao/wisp`，默认分支 `dev` |
| 10 | `gh run list --workflow=ci.yml --limit 40 --json databaseId,headSha,...` | 0 | `runlist.json`（12,300 B，40 发全 `completed/failure`） |
| 11 | `gh api repos/CarlosShao/wisp/branches/dev --jq ...` | 0 | `cc31526165… 2026-10-06T02:58:35Z` |
| 12 | `gh api repos/CarlosShao/wisp/commits/e6dc79ed --jq ...` | **1** | `HTTP 422 "No commit found for SHA: e6dc79ed"` |
| 13 | **同一发的重试**（换全号）`gh api repos/CarlosShao/wisp/commits/e6dc79ed5bd14380cce1507f94886d65a51ee63f` | **1** | 同 422 ⇒ 负读数**稳**，⛔ 不是网络抖动 |
| 14 | `gh api "repos/CarlosShao/wisp/actions/runs?per_page=15" --jq ...` | 0 | 全 workflow 最新 15 发 |
| 15 | `gh run view 37703959747 --json number,headSha,workflowName,createdAt,event,status,conclusion,jobs` | 0 | `run-37703959747.json` 14,867 B |
| 16 | `gh run view 37406757402 --json ...` | 0 | 14,822 B |
| 17 | `gh run view 37405698188 / 37396530365 / 37166458550 / 37021179942 --json ...`（循环，逐发验 rc） | 0/0/0/0 | 各 14,325–14,863 B |
| 18 | `gh api "repos/CarlosShao/wisp/actions/workflows/ci.yml/runs?created=%3E=2026-10-08T00:00:00Z&per_page=20"` | 0（**回空**） | 10-08 之后 ci **零发** |
| 19 | **正控**：同端点 `created=%3E=2026-10-07T00:00:00Z` | 0 | `COUNT=1` → `37703959747`；⇒ 证明第 18 发那枚时间过滤**真在生效**，回空不是哑过滤器 |
| 20 | **换写法再重试**：`created=%3E=2026-10-08T04:06:33Z`（＝`e6dc79ed` 的提交时刻） | 0 | `total=0` |
| 21 | `gh api "repos/CarlosShao/wisp/actions/workflows/ci.yml/runs?per_page=1" --jq '.workflow_runs[0]...'` | 0 | ci 的最后一发＝`37703959747` @ `cc31526165734e612de297848bb2080bd459ccba` |
| 22 | `gh api repos/CarlosShao/wisp/git/ref/heads/dev` | 0 | `cc31526165…`（与 #11 同） |
| 23 | `git log --format=... -S 'winlive compile gate' -- .github/workflows/ci.yml` | 0 | 唯一命中 `351e5a5e` |
| 24 | 同上 `-S 'Portable tests carrier self-test'` | 0 | 唯一命中 `1309757b` |
| 25 | `for c in 1309757b f6b79ab0 351e5a5e f8810238 e6dc79ed; do git show $c -- ci.yml \| grep -c '^+        if: '`（字面 8 空格） | 0 | 依次 **1 / 1 / 1 / 0 / 23**；同批 `^-        if: ` 全 **0** |
| 26 | `gh api repos/CarlosShao/wisp/actions/runs/35591482293`（`HEAD:ci.yml:523` 注释里引的一发号，按待验断言处理） | 0 | **真实**：`ci` @ `440dd88765`，2026-09-21，`completed/failure` |
| 27 | `git status --porcelain -- .scratch/wisp/probes/ci-if-eval-1` | 0（落笔前回空） | — |

⛔ 零 `go build`/`vet`/`test`/`list`/`env`；⛔ 零 push；⛔ 未触发任何 run；⛔ 未改任何既有文件。所有 `gh` 调用只用 `run list` / `run view` / `api` 的读端点。

## 5. 卫生

- 本腿写面＝`.scratch/wisp/probes/ci-if-eval-1/` 两枚 `.md`（`00-anchor.md`／`report.md`），⛔ 零 `.sh`／`.ps1`／`.txt`；仓外临时件在 `D:/tmp/ciifeval/`（只建不删）。
- 本文件用的 `⛔ ⇒ ★ ⚠` 与全仓文档同族。按 `A710 §7` 现量过的 scope，`tools/d22scan` 那把尺射程＝`internal/`、`cmd/`、`internal/tools/`、`frontend/`、`design/`，**`.scratch/**` 不在内** ⇒ 这些符号不挪分母；这条是**引那一枚现量**，不是本腿自跑的（本腿⛔ 没跑 d22scan）。
- 大输出全先落文件再读：最大的一发 JSON 14,867 B，⛔ 未把整份 `gh` JSON 或 CI 日志贴进任何回报。
