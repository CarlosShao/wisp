# `00-anchor.md` — 普查腿 `310-a1` 起手锚（票 310 `AC#0`，只读）

leg=`310-a1`  ticket=310  ac=`AC#0`  role=只读普查腿（非实现者／非验收腿）
写面＝仅 `.scratch/wisp/probes/310/a1/**` ＋ 票 310 尾一条 append-only Progress bullet
本件第一笔单独 commit；add 与 commit 都带显式 pathspec。

---

## 1. 锚点现量（本腿自己跑，非引用派单）

| 项 | 尺 | 读数 | rc |
|---|---|---|---|
| HEAD | `git rev-parse --short HEAD` | `49153fa8` | 0 |
| 分支 | `git rev-parse --abbrev-ref HEAD` | `dev` | 0 |
| 产码面脏度 | `git status --porcelain -- internal cmd tools scripts` 行数 | `0`（empty above＝0 枚 ⇒ 产码面干净） | 0 |
| 全索引暂存枚数 | `git -c core.quotePath=false diff --cached --name-only \| grep -c ''` | `0` ⇒ **没有别家在 index 里的路径，我没动任何东西** | 0（计数尺带 `\|\| true`） |
| 全工作树脏行（旁证，非我的面） | `git status --porcelain \| grep -c ''` | `841` ⇒ 共享工作树处于在飞态，其中 `.scratch/wisp/probes/**` 下大量 `M` 是别家／机主的；我⛔ 动其一 | 0 |

派单里**不写锚号**（§0.1），这一节的号＝我盘上现量，与编排者转述冲突时以本节为准。

## 2. 起手闸门正典逐字抄（＝`.scratch/wisp/probes/orch/gate-roster-304.md` §1，provenance＝HEAD blob，`git show HEAD:.scratch/wisp/probes/orch/gate-roster-304.md`）

```
tasklist //FI "IMAGENAME eq wisp.exe"          # 命中行数＝枚数
tasklist //FI "IMAGENAME eq balldebug.exe"     # 命中行数＝枚数
tasklist //FI "IMAGENAME eq mockllm.exe"       # 命中行数＝枚数（票 304 起必数：测试起的"假大脑"会挂着一枚回环套接字不退）
tasklist //FI "IMAGENAME eq wisp.test.exe"     # 命中行数＝枚数（`go test -c` 的产物，同理会留）
```

⚠ **`//FI` 只认精确名**，所以上面那行**只覆盖 `wisp.test.exe` 一枚**，漏掉兄弟名 `<包>.test.exe`（如 `tools.test.exe`／`llm.test.exe`）。
⇒ **家族另用一条宽松尺**（票 304 `AC#4` 的条件，2026-10-11 20:1x 由 `304-v1b` 判出、编排者落进本条；⛔ 拿精确名冒充家族）：

```
powershell -NoProfile -Command 'Get-CimInstance Win32_Process | Where-Object { $_.Name -like "*.test.exe" } | Select-Object ProcessId,ParentProcessId,Name,CommandLine | Format-List'
```

⚠ 这条尺**命中 0 时也要落盘**（把命令的 stderr 与 RC 一起写进名册，⛔ 用"空输出"当"0 枚"；在 bash 里请用**单引号**包住整条 PS 命令，否则 `$_` 会被外壳吃掉）。

**数到 ≠ 归我**。任何一枚非零计数，动手之前先答三问，答不全就**停下来上报**，⛔ 不许停：

1. **谁起的**——它的命令行里那个具名路径是什么（`Get-CimInstance Win32_Process -Filter "ProcessId=<pid>" | Select CommandLine`；
   `mockllm.exe` 的命令行自带它的 `t.TempDir()` 全路径，路径首段就是测试名；webview 系看 `--webview-exe-name=`）。
2. **父 PID 还在不在**——`ParentProcessId` 现量查得到吗？⚠ **这一问只能当旁证，不能当归属**：
   Windows 会**复用 PID**（票 304 现量：曾被判"已不存在"的 `11724` 后来查得到＝`RuntimeBroker.exe`，只是同号换了人）。
   ⇒ 归属只认**命令行里的具名路径**；路径认不出来的＝**归属不明**，⛔ 归本腿。
3. **本腿认下哪几枚**——逐枚写 PID＋来路；写不出来的那几枚在报告里标"别家的，⛔ 动"。

## 3. 闸门现量（本腿当场跑的一发，件＝`logs/00-gate-310a1.txt`，`wc -c`=1089 非 0）

| 尺 | 读数 | 自落 rc |
|---|---|---|
| `tasklist //FI "IMAGENAME eq wisp.exe"` | `hit_lines=0`（empty above＝0 枚） | `tasklist_rc=0` |
| `tasklist //FI "IMAGENAME eq balldebug.exe"` | `hit_lines=0`（empty above＝0 枚） | `tasklist_rc=0` |
| `tasklist //FI "IMAGENAME eq mockllm.exe"` | `hit_lines=0`（empty above＝0 枚） | `tasklist_rc=0` |
| `tasklist //FI "IMAGENAME eq wisp.test.exe"` | `hit_lines=0`（empty above＝0 枚） | `tasklist_rc=0` |
| 家族宽松尺 `Get-CimInstance … *.test.exe` | 空输出，`family_named_lines=0`（empty above＝0 枚），stderr 与命中数同批落盘 | `ps_family_rc=0` |
| **本票新增的尺**：`%TEMP%` 下目录枚数／其中 `Test*` 枚数 | `total=2754`／`test=50`／`temp=C:\Users\swq\AppData\Local\Temp` | `temp_count_rc=0` |

- ⚠ 编排者 10-12 10:21:55 那一发的快照（`probes/orch/20261012/logs/orch-gate-20261012.txt`：四把＝0／家族尺空／`Test*`＝50）**我只当参照不当读数**；上面这一列是本腿现跑。
- 本腿**一枚目录都没删**、一个进程都没杀（起手时非零计数＝0 枚，无三问可答）。

## 4. 派单 §0 第 3～4 步（逐字抄，provenance＝工作树未跟踪件 `.scratch/wisp/probes/orch/20261012/txt/310a1-dispatch.txt`，`git ls-files` 判＝未入库）

> 3. **起手闸门逐字跑一遍**（正典＝`.scratch/wisp/probes/orch/gate-roster-304.md` §1，把那段**逐字抄进**你的 `00-anchor.md` 再跑）：四条 `tasklist //FI "IMAGENAME eq <名>"`（`wisp.exe`／`balldebug.exe`／`mockllm.exe`／`wisp.test.exe`）＋那条家族宽松尺（`Get-CimInstance Win32_Process | Where-Object { $_.Name -like "*.test.exe" }`，**bash 里整条 PS 命令用单引号包住**，否则 `$_` 被外壳吃掉）。⚠ 命中 0 也要落盘（stderr＋rc 并排），⛔ 空输出当"0 枚"。编排者 10-12 10:21:55 那一发的读数在 `probes/orch/20261012/logs/orch-gate-20261012.txt`（四把＝0／家族尺空／`%TEMP%` 下 `Test*`＝**50**），**那是快照，⛔ 当你的读数**。
> 4. **第 1 笔 commit 立刻做**：只交 `00-anchor.md`，形式＝`git add -- <显式路径>` 之后 `git commit -F <你自家 logs 里的消息件> -- <同样的显式路径>`。**⛔ 裸 `git commit`**（它提交的是整个 index；add 带 pathspec、commit 也带 pathspec）。⛔ `git add -A`／`git add .`／`git add <目录>` 一把扫。

## 5. 派单 §3 硬约束（逐字抄，全六条）

> 1. ⛔ 删任何目录或文件（包括你自己的残渣、包括"看着明显是垃圾"的那些）；⛔ 移动、改名、清空 `%TEMP%` 下任何东西。这一程**只读＋只新建你自己的 `probes/310/a1/**`**。
> 2. ⛔ 动任何产码／`*_test.go`；⛔ 改 `ci.yml`／`.github/**`；⛔ 动 `docs/reports/**`／`docs/evidence/s1/**`（裁决表归非实现者）／`AGENTS.md`／`docs/PLAN.md`／`docs/specs/**`。
> 3. ⛔ 放宽任何断言、⛔ 造 `--- SKIP` 换绿；⛔ 动 SLO 阈值／golden／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt`／D43 转移表／C1–C32／D1–D47／三枚冻结件（`internal/panel/tokens_fourway_test.go`、`internal/panel/l2_grant_boundary_test.go`、`internal/perm/ticket90_persist_test.go`）。
> 4. ⛔ 碰 `msedgewebview2.exe` 那几棵树；⛔ 杀任何不是你起的进程。要停进程只能按三条守卫（名字必须是 `*.test.exe`／`mockllm.exe` ＋ 命令行含你那棵导出树的具名路径 ＋ 父 PID 现量不存在），且逐枚写进件里。
> 5. ★**Go 编译面**：编排者 10-12 上午把这一面交给你独占。⛔ 跑整包 `go test ./...`（会把 `cmd/wisp` 那批真窗/音频用例卷进来，读数必假）；只跑你为表② 选的那几条窄用例，并在件里具名"我只跑了这几条"。
> 6. ⛔ 在仓内建 worktree／checkout；临时件**只建不删**，只建在你自家 `probes/310/a1/**` 里（导出树在 `mktemp -d`，那是仓外，允许留）。

## 6. 派单 §4 带 ★ 的规矩（逐字抄）

> ★**证据件规矩**：⛔ 叫 `.out`（根 `.gitignore` 第 8 行全仓 `*.out`，会被静默跳过）⇒ 一律 `.md`／`.txt`／`.tsv`；⛔ 在 `probes/**` 下新建 `.go`（会挪别家的 gofumpt／d22scan 分母）；件名⛔ 带 `§`、⛔ 长名（相对 ≤121、末段 ≤100、全路径 ≤165；新建前先跑一次 `sh scripts/check-path-length-budget.sh --with-self-test` 那把长度门）。
>
> ★**每把读数件自己落一行 `rc=N`**；⛔ 0 字节件（＝那格没交）。若你要引用一枚同伴件确认"它确实是空的、不是没交"，**先 `wc -c` 再写**。
>
> ★**尺规矩**：`grep -c` 命中 0 时 rc＝1，会把 `&&` 后面的尺整段跳掉 ⇒ 凡"某枚计数应当为 0"的链一律 `N=$(… || true)`，件里写一行"empty above＝0 枚"；任何"0 命中"的尺，**stderr 要与命中数并排落盘**（一把坏尺＋`2>/dev/null` 会造出像读数的假 0）；引用外部文档/票面原文时标 provenance（HEAD blob／工作树／仓外）。
>
> ★**门禁**（本格只要求这三把，逐件自落 rc）：`sh scripts/check-path-length-budget.sh --with-self-test` rc=0；`sh scripts/d22scan.sh` rc=0（纯净树；⛔ 让你的临时件进名册——它只扫跟踪件，所以 commit 前先 `git status --porcelain` 确认你新增件只在你自家目录）；`git show --name-only --format=` 逐笔量你自己的每笔 commit（越出上面三件套＝**直接停手上报**）。⛔ `GOFLAGS= go build ./...`／⛔ 整包 `go test`（那是写腿与编排者的面）。

派单 §2 表① 的两枚 ★（名册逐枚理由列／分母同表）与 §2 表② 的 ★（跑完当场自报自己新增的目录）随件交出，见 `10-roster.md`、`20-shapes.md`。

## 7. 长度门现量（本腿自己跑）

`sh scripts/check-path-length-budget.sh --with-self-test` → `rc=0`，末行 VERDICT GREEN（件＝`logs/00-pathlen.txt`）。
本腿新建件全部落 `.scratch/wisp/probes/310/a1/` 下：该相对目录前缀长 32 字符，件名 ≤24 字符 ⇒ 相对路径 ≤56，远在 121 帽内。

## 8. 起手对票面／派单的顶回（先记，后续发现续加）

- **无过期冲突**：派单 §0.3 写的编排者快照（四把＝0／家族空／`Test*`＝50）与本腿现量方向一致；票 310 标题那句「从立案的 21 枚长到 44 枚」是 10-11 19:5x 的读数，今天现量＝**50**，票面 §现量 自己写了「⛔ 引用前先重跑，枚数一律当快照」⇒ 不构成矛盾，只按快照处理。
- **尺口径缺口（具名报回，本腿已自作补正并写明判据）**：派单表① 的反查尺「去掉尾部随机数字后的名字拿去 `git grep -n "func <名>"`」对**两类名字会漏**：① Go 的 `t.TempDir()` **第一次调用**的目录名＝测试名本身、**不带任何尾部数字**（`testing.tRunner` 只在 `counter > 0` 时才追加数字），"去尾数"这条尺会把 `TestX1ng` 之类**合法以数字结尾的用例名**削错；② **子测试**里的 `t.TempDir()` 名字形如 `TestFoo_subname1`（`/` 与空格被换成 `_`），"去尾数"后拿原名 `git grep` 必零命中。⇒ 本腿名册对每枚**同时**试三种形态（原名／去尾数原名／首段 `_` 前主名＋`_` 还原成 `/` 的子测试形），逐枚记命中用了哪一形，不靠单一把尺下结论。

## 9. 锚位更正（第一笔 commit 之前现量，⛔ 让 §1 那行停留在过期读数上）

- §1 表里那行 `HEAD=49153fa8` ＝本腿**第 1 次调用**的现量，量没错，但**在它落盘之前同伴腿 `done-check-2` 交了自己的锚件**，仓尖已移动。
- 复量（第一笔 commit 之前当场跑）：`HEAD_now=34184dd7`（`done-check-2: anchor at 49153fa8 + verbatim re-copy of dispatch §0.3-4 and §2 star rules`）、`parent_now=49153fa8`。
- ⇒ 本腿的**起始锚＝`49153fa8`**（产码面 §1 的读数就是在这枚尖上量的），**本腿第一笔 commit 的父＝`34184dd7`**。两枚都记着，⛔ 用一枚冒充另一枚。
- 复量全索引暂存枚数＝`staged_now=0`（同伴腿已把它的件提交掉，index 空，我没动过它）。
- ⛔ 碰 `done-check-2` 的写面（`probes/orch/20261012/done-check-2/**`）；本程与它的唯一交集＝共享 index 与工作树，两次现量都证明 index 空、我没清过别人的东西。
