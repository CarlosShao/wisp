# `274-v3` 裁决表 —— 票 274 AC#11「本趟构建出处」闸门的**含空格名假绿洞**是否被真修掉（腿＝非实现者，攻 `274-r3` 落进仓里的那道数组对拉）

被审件＝`scripts/build.ps1`（盘上 **313 行**、md5 `43fb7e20c597192104728befb231e8ea`，与 `git show 674c1920:scripts/build.ps1` 逐字节相同：`90-final-gates.txt:13-15` 两串全等）。
判据原文＝`.scratch/wisp/issues/274-no-nail-requires-shipped-exe-to-carry-page-build-ps-step-2-comment-is-false-in-both-halves.md:59`（AC#11，整行逐字读完，>1000 字符），核心两句：**「在仓外镜像树的 `frontend/dist` 里预放一枚名字含空格的页面文件，其余照 v2 那发的形状（桩 npm 一字不写）⇒ 必须红」**＋**「修形＝跑前跑后都按数组对拉（不要过字符串），⛔ 不许用"忽略带空格的名字"换绿」**。⛔ 本腿 AC 框一枚没翻。⛔ `probes/274/r3/logs/10·11·12·13` 一字节没当凭据引，四发全部本腿自跑。

## §0 起手锚与台件（`logs/00-anchor.txt`、`logs/01-setup.txt`；第 1 笔 commit `83badb67` 早于任何长跑命令）

- `date`＝`Wed Oct  7 17:31:21 CST 2026`；分支 `dev`；HEAD＝`10899b77`（A678 收 274-r3）；被审修形 commit＝`674c1920`。
- 六族 porcelain（`-- cmd internal scripts .github docs frontend build`）＝**0 行**（`:15-16`；终态仍 0 行，§4）。全仓 porcelain＝753 行＝别人在飞的数（`design/**` 等，`00-anchor.txt:17-18`），本腿零触碰。
- 起手仓内 `frontend/dist` 四枚 md5＝`d41d8cd9…`／`f98bfc4b…`／`db4db7a2…`／`70128a3d…`，`build/wisp.exe`＝`e6c8e52bb14f05…` 见 `:20-25`（逐字：`e6c8e52bb14f15a7983bbc6b093e6058`）。
- 简报行号逐条对盘上复认（`00-anchor.txt:26-45`＋全文现读）：**313 行**✓、md5 ✓、`function Fail`＝`:47`✓、全文唯一 `exit 1`＝`:49`✓（计数 1）、名册快照数组＝`:158`✓、判定＝`:209-210`（`-notcontains` 数组对拉）✓、`-split` 全文命中 **0**✓、出处红因 `Fail`＝`:212`✓、`git rev-parse`＝`:240`✓。**编排者这串行号无一枚过期。**
- 台件＝仓外 `D:/tmp/wisp274v3/tree-a·tree-b·tree-c`：`cp -al` 逐条顶层（⛔ 不含仓内 `.git`、⛔ `build/` 只建空目录防截断写盖仓内 exe、⛔ 不含 `frontend/node_modules`），`frontend/dist` 用**常规 `cp`** 拷入仓内那 4 枚陈旧真页面的字节（先 `rm` 镜像里的硬链再拷 ⇒ 仓内 inode 全程只读，`01-setup.txt` 末段 md5 复跑全等）；`tree-a`/`tree-b` 预放含空格陈旧页 `stale spaced page.html`（31 字节），`tree-c` 不放。桩 npm 两枚（`stubs/npm-stub-nowrite.cmd`＝任何调用退 0 一字不写；`stubs/npm-stub-write.cmd`＝仅在 `run build` 时往指定树写那枚含空格新名）——**全部用 Write 工具落盘，没用 bash printf**；PATH 前置桩目录后 `(Get-Command npm.cmd).Source` 逐字进每发日志。调用形态＝`powershell -NoProfile -ExecutionPolicy Bypass -File <树>\scripts\build.ps1 -Env dev`（PS 5.1，与 CI 同形）。
- ⚠ **台件偏差一枚（具名，见 §3-2）**：简报的「把 `git.exe` 从 PATH 摘掉」本机做不到也不该做；本腿改用**镜像树内 `git init` 新建空仓（`--allow-empty` commit）**中和 `:240`。仓内 `.git` 内容零拷入、仓外树各自独立。`node_modules` 全程硬链没用、真 npm 没跑（两枚桩），构建走 `node_modules ABSENT ⇒ npm ci` 支由桩 `ci` 退 0 满足。

## §1 逐格裁决

### ① 复现「假绿被修掉」＝**成立**

`tree-a`：陈旧真页面 4 枚＋含空格陈旧页 1 枚预放，桩 npm（nowrite）报成功、一字不写。

命令（逐字）＝`powershell -NoProfile -ExecutionPolicy Bypass -File "D:\tmp\wisp274v3\tree-a\scripts\build.ps1" -Env dev`（PATH 前置 `D:\tmp\wisp274v3\stub-nowrite`）。

- **rc＝1**（`logs/10-case1-spaced-stub-nowrite.txt:22`：`CASE_01-spaced-stub-nowrite_PROCESS_EXITCODE=1`）。
- 前置快照逐字（`:18`）：`build.ps1: frontend: pre-build dist snapshot = 5 file(s) [.gitkeep assets\index-BRKj5OIJ.css assets\index-BVKlegVD.js index.html stale spaced page.html] (this run's npm must produce an artifact name outside it).` ⇒ **含空格整名进了快照、没被拆碎**（这正是修形后的应有形状）。
- 红因逐字（`:21`）：`build.ps1: FATAL: frontend step: npm run build reported success but frontend/dist holds no artifact name outside the snapshot taken before this run's npm (…snapshot was 5 file(s) [… stale spaced page.html], dist now holds [.gitkeep assets\index-BRKj5OIJ.css assets\index-BVKlegVD.js index.html stale spaced page.html]). Named cause: unprovenanced page bytes …` ⇒ **`Named cause: unprovenanced page bytes` 在**，跑在 `:212` 那支（红因文本与盘上 `:212` 逐字吻合）。
- 跑后 `tree-a` dist 名册逐枚未变（`:23`），仓内 4 枚 md5＋exe 逐枚未变（日志尾 `Repo md5 AFTER`）。

### ② ★恒真检查（牙是否恰在 `:209-210`）＝**成立（含空格假绿原形复现）**

`tree-b`：与 ① 同形，唯一突变＝把 `:209-210` 的数组对拉**换回旧的 `-join ' '`／`-split ' '` 两趟**（`logs/02-revert-scope-manifest.txt:4-20` 的 diff 逐字：仓内那份的两行 `$preRoster`／`$postRoster` 直比，被换成 `@(($preRoster -join ' ') -split ' ')` 与 `@(($postRoster -join ' ') -split ' ')`) 再对拉；除此一处外与盘上 313 行版逐字节相同。⛔ 仓内那份一字节没动（终态 md5＝起手，§4）。

命令＝`powershell -NoProfile -ExecutionPolicy Bypass -File "D:\tmp\wisp274v3\tree-b\scripts\build.ps1" -Env dev`（PATH 前置同一枚 nowrite 桩）。

- **rc＝0**（`logs/11-case2-reverted-join-split.txt:47`：`CASE_02-reverted-join-split_PROCESS_EXITCODE=0`）⇒ **① 的红确实由 `:209-210` 这一处给出**，不是别的支顺路拦的。
- 假绿签名逐字（`:21`）：`build.ps1: frontend ok: 5 file(s) … || provenance: this run's npm added new artifact name(s) [stale spaced page.html] over the pre-build snapshot 5 file(s) [.gitkeep assets\index-BRKj5OIJ.css assets\index-BVKlegVD.js index.html stale spaced page.html]; gone since snapshot: [stale spaced page.html]` ⇒ 桩一字不写、名册一枚未变，同一枚**整名**却被同时报成「本趟新增」和「快照后消失」——碎片对不上整名的双向指纹。
- 整条通路走完：`:22` `go build ok (cgo linked against sherpa-onnx C API)`、`:46` `wisp doctor: PASS`。⇒ 旧那一版的假绿不只是闸门放行，是**真出货一枚带陈旧页面的 exe**（这一形正是 AC#11 立案要拦的）。

### ③ 牙没钝（不带空格名）＝**成立**

`tree-c`：只放陈旧真页面 4 枚（无空格名），桩 npm（nowrite）一字不写。

命令＝`powershell -NoProfile -ExecutionPolicy Bypass -File "D:\tmp\wisp274v3\tree-c\scripts\build.ps1" -Env dev`。

- 快照逐字（`logs/12-case3-nospace-stub-nowrite.txt:18`）：`pre-build dist snapshot = 4 file(s) [.gitkeep assets\index-BRKj5OIJ.css assets\index-BVKlegVD.js index.html]`。
- 红因逐字（`:21`）＝同一句 `Named cause: unprovenanced page bytes …`；**rc＝1**（`:22`）。⇒ 修形后 ①(c)／陈旧页那一形照旧红，`274-v2` 量的 AC#10 牙没被这次改动钝化。

### ④ 绿路没断（桩真写一枚含空格新名）＝**成立**

`tree-c`（③ 跑后名册未变）＋桩 `stub-write`：`run build` 时真写 `stale spaced page.html` 一枚（31→43 字节的新页），其余一字不动。

命令＝`powershell -NoProfile -ExecutionPolicy Bypass -File "D:\tmp\wisp274v3\tree-c\scripts\build.ps1" -Env dev`（PATH 前置 `D:\tmp\wisp274v3\stub-write`）。

- 快照 4 枚（`logs/13-case4-spaced-written-green.txt:18`）、桩干活行（`:20`）`274-v3 stub-write npm ran build, wrote the spaced artifact`。
- 成功行逐字（`:21`）：`build.ps1: frontend ok: 5 file(s) in frontend/dist (entry index.html is 1044 bytes): .gitkeep=0 assets\index-BRKj5OIJ.css=49943 assets\index-BVKlegVD.js=553469 index.html=1044 stale spaced page.html=43 || provenance: this run's npm added new artifact name(s) [stale spaced page.html] over the pre-build snapshot 4 file(s) […]; gone since snapshot: []` ⇒ **整名原样出现在新增名列表里，没有被拆成 `stale`／`spaced`／`page.html` 两段**；`gone` 空集＝名册只增没错删。
- **rc＝0**（`:47`），整条通路走完：`:22` `go build ok`、`:46` `wisp doctor: PASS`。

### ⑤ 带空格名会不会让别的通道串位（只读判读＋一枚实测）

| 通道 | 判语 | 凭据 |
|---|---|---|
| `wisp.exe panel-assets -manifest` | **不会串位（打印侧）**。`internal/panel/assets.go:114` 逐字 `out = append(out, fmt.Sprintf("%s %d %s", p, len(data), hex.EncodeToString(sum[:8])))`＋枚枚一行打印 ⇒ 实测：拿 ② 那枚**假绿构建产出的树 exe** 跑（`tree-b/build/wisp.exe`），`panel assets embedded: 5 files, entry=index.html built=true`（rc=0）与 `-manifest` 末行 `stale spaced page.html 31 db88a040f8812b50`（rc=0）——整名原样、枚行分隔（`logs/02-revert-scope-manifest.txt:38-47`）。⚠ 顺带登记（不属 AC#11）：该行格式若被**从左按空格切**的读者消费会歧义（名可含空格、size/hash 不可，右切两 token 无歧义）；本仓现读无这种读者，`build.ps1` 也不消费 `-manifest`。 | `logs/02:36-47` |
| `Get-ChildItem` 的 `FullName.Substring($DistDir.Length + 1)` 起点 | **不会串位**。起点只由**目录**路径长度决定（`:158`/`:208`/`:214` 三处同一表达式），文件名含什么字符与起点无关；实测正反两向都在 ①④ 名册里逐字可见（含空格名与 `assets\…` 整名均原样提取参与数组比对，`10:18`、`13:21`）。 | 盘上 `:158`/`:208`/`:214`＋①④ 读数 |
| `$OFS` 隐式拼接 | **不会串位（判定侧）**。`:212`/`:215` 的 `[$preRoster]` 插值与 `:161`/`:214`/`-join ' '` 全部只是**打印文本**；判定输入永远是数组本身，全文 `-split` 命中 0（`00-anchor.txt` 计数行）⇒ 没有任何代码把这些串再拆回去。实测：① 的红因里 `dist now holds [.gitkeep … stale spaced page.html]` 整名原样（`10:21`）。⚠ 唯一残留＝**人读歧义**：`$distRoster` 把 `name=size` 对再 `-join ' '`，含空格名时字段边界肉眼不可分（`13:21` 逐字 `… index.html=1044 stale spaced page.html=43 …`），机器读者不存在，不构成假绿面。 | 盘上 `:161`/`:212`/`:214`/`:215`＋①④ 读数 |

〔判不动〕无——三通道都拿到了读数或逐字代码＋实测组合。

## §2 对 AC#11 本身的裁（给编排者翻勾用，本腿一枚不翻）

**成立**。四发全中且互锁：① 修后含空格名假绿形**红**（`Named cause: unprovenanced page bytes`，rc=1）；② 只把 `:209-210` 换回旧 `-join`/`-split` 两趟**同一形即复绿**（rc=0、整枚 exe 真出货）＝牙**恰好**这一处；③ 无空格陈旧形照旧红＝牙没钝；④ 含空格真新名照旧绿且**整名原样打印**＝没拿"忽略带空格名"换绿（盘上也查无按名过滤的逻辑，判定两支 `:209-210` 无任何名字模式豁免）。修形与票面「跑前跑后都按数组对拉（不要过字符串）」逐字相符：`:156-158` 数组快照、`:208-210` 数组判定、串只在打印位。残余见 §5，均不否决本格。

## §3 ⓐ 推翻／收窄编排者的哪几句（⛔ 以盘上为准）

1. **无一枚行号或 md5 被推翻**：简报那串行号（`:47`/`:49`/`:158`/`:209-210`/`:212`/`:240`、313 行、`-split` 0、md5 `43fb7e20…`）逐枚对盘上复认（§0），且 `git show 674c1920` blob＝盘上逐字节。票面 `:59` 原文里的 `:152`/`:198`/`:200` 是 91c90aa1 那一版的旧行号，票面 `:60` 已自我改口——本腿复认改口后的那串就是盘上现值。
2. **推翻（台件指令的可执行性）**：简报「把 `git.exe` 从 PATH 摘掉」——`git rev-parse 2>$null` 在无 `.git` 目录下于 `$ErrorActionPreference='Stop'` **确实炸成终止错误**（编排者担心的这一半由实测坐实：`logs/05-git-throw-test.txt:4` 逐字 `THREW git=C:\Users\swq\.qoder-cn\bin\git\mingw64\bin\git.exe EXC=System.Management.Automation.RemoteException MSG=fatal: not a git repository…`）；**但摘 PATH 这枚修法在本机不可取**：`(Get-Command git.exe)` 命中多枚（qoder 自带 mingw64、`D:\work\soft\Git\{cmd,bin}`、以及 build.ps1 自己 `:78-79` 会把 gcc 所在 mingw64 目录再前置回去），摘不干净且连坐 gcc。⇒ 本腿的实际处置＝镜像树内 `git init`＋`--allow-empty` commit（仓外、零仓内 `.git` 内容），四发的 `:240` 全部安静走过。**今后派构建脚本腿请把这句改成「镜像树自带空 .git」**。
3. **收窄（假绿指纹的形状，给台账加一枚可读信号）**：旧 `-join`/`-split` 那版的假绿在成功行里留**双向签名**——同一枚整名同时出现在 `added new artifact name(s) [...]` 和 `gone since snapshot: [...]`（`11:21` 逐字）。今后任何人核这类读数，看到「新增与消失同名单」即可断定名册过了字符串，不必重跑。
4. **登记（不推翻）**：`674c1920` 的 numstat 逐枚现量（`logs/02:24-34`）＝`scripts/build.ps1` **14/6** 与简报一致，另 9 枚全在 `probes/274/r3/**` 自己写面内＝r3 越界枚数 0（尺＝numstat 清单逐行，无 `cmd internal scripts(.github|docs|frontend|build)` 之外的产码）。

## §4 ⓑ 终态自证（`logs/90-final-gates.txt`，`date`＝`Wed Oct  7 17:39:21 CST 2026`）

- 六族 porcelain＝**0 行**（`:5-6`）＝起手。
- 仓内 `frontend/dist` 四枚 md5、`build/wisp.exe`（`e6c8e52b…`）、`scripts/build.ps1`（`43fb7e20…`）终态＝起手逐枚全等（`:7-13`；每发日志尾的 `Repo md5 AFTER` 亦逐发全等）。⛔ 全程没删没盖仓内任何一枚。
- `git show 674c1920:scripts/build.ps1 | md5sum`＝`43fb7e20c597192104728befb231e8ea`＝盘上（`:13-15`）。
- 门禁两把全过：`sh scripts/d22scan.sh` rc＝**0**，且先看见自带正控逐字 `runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0`（`:255`，扫库在 `:256` 之后）；`go vet ./cmd/wisp/ ./internal/panel/` rc＝**0**（`:269`）。⛔ 整包 `go test` 没跑（窗口归产码腿）、⛔ 没格式化任何 `.go`。
- ② 的仓外突变可逆性由 md5 对拉自证：仓内那份全程 `43fb7e20…`，突变只活在 `tree-b` 的常规拷贝里。

## §5 残余（具名，⛔ 不写成通过）

1. **桩 npm 只测了 `run build` 与 `ci` 两支的退出面**；「npm 往 `node_modules` 写文件而 dist 一字不动」那形＝`274-v2` ①(b) 已量过（当时牙同一枚 `Fail`），本腿没重跑（不在 AC#11 射程，且 v2 凭据在案）。
2. **含全角空格／制表符的名字**：`-split ' '` 只管 ASCII 空格；票面 AC#11 字面只要求空格形，本腿按字面打；全角空格那一形在修后版**同样红**（读码可断：数组判定不过字符串，`-join`/`-split` 都不在判定路径上），但没实测。
3. **CI 侧**：⛔ 零推送；CI 三处调用会不会带空格名进 dist 属真实 vite 不产那种名的老口径（`274-v2` §5-3 已具名），本腿没碰。
4. **`$distRoster` 的人读歧义**（⑤ 第三行）＝纯打印面形状问题，不属本票射程，登记交台账。

## 裁决摘要（一句话版）

**AC#11 的牙真在盘上、且只在那两行**：修后含空格名假绿形红在 `:212` 具名因（①），把 `:209-210` 换回旧 `-join`/`-split` 立刻复绿并真出货（②＝恒真检查成立），无空格陈旧形照旧红（③＝牙没钝），真写含空格新名照旧绿且整名原样打印（④＝没拿豁免换绿）；三条相邻通道（`-manifest` 枚行制、`Substring` 目录长起点、`$OFS` 只进打印）逐条判**不会串位**，其中 `-manifest` 通道由假绿构建出的那枚 exe 实测钉死。另实测坐实编排者担心的「镜像树无 `.git` ⇒ `:240` 在 EAP=Stop 下炸终止错误」为真，但「摘 PATH」不可行，正解＝镜像树自带空 `.git`。
