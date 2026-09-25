# observe AC#15 — 非实现者终裁（final r1）

agent: `acceptor-observe-ac15-final-r1`（**非实现者**；本程不是 `poll` 程也不是 `linux` 程）
ticket face: `.scratch/wisp/issues/136-the-zero-sample-fail-closed-guard-...md` `AC#15`（现量格在 `:348`，`[ ]` 未勾）
被裁的两枚交付：
1. `docs/evidence/s1/observe-ac15-poll-r1.md`（修法本体，498 行，commit `d8390aa`＋`333dfe3`）
2. `docs/evidence/s1/observe-ac15-linux-r1.md`（Linux 分母程，454 行，commit `22be55e`＋`dc44c24`）

> **开头先把归因钉住（台账 `A226⑤` 要求的那句）**
> `observe-ac15-linux-r1.md` 的 **§5 那一节，那个程一个字都没写**（它撞在 150 轮上限上，止于 §4.6）。
> **本程没有、也不会去替它补 §5** —— 本文件是**另一枚程的另一份文件**，只回答它 §3.5 承诺要答的那一问（本文 §4），
> 并按派单出 `AC#15` 的终裁（本文 §6）。两程的笔不混在一份文件里。

本程的角色与写面：
- **唯一可写路径**＝本文件。**没有**改过 `internal/observe/**`、**没有**碰票面（`.scratch/**`）、**没有**碰台账
  （`docs/reports/pending-and-issues.md`）、**没有**做任何 `A##` 编号。
- **没有跑 `go test ./...`**：本程开工时同仓有两枚兄弟程在飞（一枚写 `internal/panel/l2_grant_boundary_test.go`、
  一枚写 `docs/evidence/s1/d22scan-close-r1-accept-r1.md`），派单明令只跑 `./internal/observe/`。
- 所有需要注入的测量（饿窗、正控）只发生在**仓外副本** `D:\tmp\observe-ac15-final\`（只建不删）。

---

## 0. 锚点、挂载证明、两枚门的重跑

### 0.1 锚点（本程每次读数各自贴时刻；HEAD 在共享树里一直在漂）

```
$ date "+%Y-%m-%d %H:%M %z" ; $ git rev-parse HEAD
2026-09-25 11:04 +0800   6b195cf06b2950bd8b45526efb741e716620651c   ← 开工
2026-09-25 11:09 +0800   598620efe585653eadfdddd5e225ee4efcffad12   ← 第一发 go test 之后
2026-09-25 11:11 +0800   bd9861ced5442170acfe6f965cc98a4f9b23fb27
2026-09-25 11:12 +0800   c803840bbf131104c04aac4d99844460643590c7   ← 本节写于其后
```

⚠ **漂的是别人的地界，不是被验的那枚包**（这一条承重，所以三锚各跑一次）：

```
$ for a in 6b195cf 598620e bd9861c c803840; do git diff --name-only 333dfe3..$a -- internal/observe/; done
（四枚全部无输出）          ⇒ 自 AC#15 落地的 333dfe3 起，internal/observe/ 在本程读数的每一个 HEAD 上一字节没动
```

被验六枚文件的字节（工作树 md5 vs `git show 333dfe3:` md5，**逐枚同值**）：

```
19de4896ea4b5122fbbc31418a85c59c  window_wait_136_test.go
87896a7721ba546e3d2c10aa9eb21388  sampler_settle_coverage_136_test.go
0ebfc98d84f808db32c66d606ae678eb  sampler_settle_gate_136_test.go
781206946df861b713ea33b25f163466  sampler_test.go
6de900e183dba3e111cb5cce85832296  sampler_settle_zerosample_136_test.go
bce93602fe280166b97390bcd0b35cc3  sampler_zerosample_136_test.go
```

⇒ **顺带一枚三方对上的字节账**：前两枚 md5 与 `observe-ac15-linux-r1.md` §0.4 贴出的
`19de4896…` / `87896a77…` **逐字同值** ⇒ 本程 Windows 树、本程容器里那棵树、与归档 Linux 程挂载的那棵树，
是同一版码。（本程自己量的，不是引它那句。）

### 0.2 容器复用与**挂载非空证明**（派单点名的静默陷阱）

复用已存在的停止态容器（台账 `A226⑥`：上一程只 `stop` 没 `rm`）：

```
$ docker ps -a --filter name=wisp-obs --format '{{.Names}} | {{.Image}} | {{.Status}}'
wisp-obs-ac15-linux-r1 | golang:1.27 | Exited (137) 4 minutes ago
$ docker inspect ... --format '{{range .Mounts}}...'
bind  D:/work/workspace/projects plans/Wisp -> /repo   rw=false      ← 仓库是**只读**挂载
bind  D:/tmp/observe-ac15-linux-r1-s1      -> /scratch rw=true
volume …/wisp-obs-linux-gomod   -> /go/pkg/mod ；volume …/wisp-obs-linux-gocache -> /root/.cache/go-build

$ docker start wisp-obs-ac15-linux-r1
$ docker exec wisp-obs-ac15-linux-r1 sh -c 'ls /repo | wc -l ; ls /repo | head -25 ; cd /repo && go list ./internal/observe ; go version ; nproc ; uname -r'
20
AGENTS.md  README.md  balldebug.exe  build  cmd  deps.toml  design  docker  docs  frontend
go.mod  go.sum  internal  models  scripts  signmodels.exe  third_party  tmp-none  tools  wisp.exe
github.com/CarlosShao/wisp/internal/observe          ← 空目录产不出这一行
go version go1.27.1 linux/amd64
12
6.6.114.1-microsoft-standard-WSL2
```

⇒ **挂载证明成立在先，之后才有本程任何一枚 Linux 读数**（派单要求）。
容器 `/repo` 是 `:ro` ⇒ 本程结构上写不进仓库。
台件实话（与 Linux 程同一味，本程自量）：`nproc`＝**12 枚 vCPU**、Debian/WSL2 内核，
**不等于 GitHub `ubuntu-latest`（公开规格 4 枚）**；`ls -ld /tmp` ＝ `drwxrwxrwt`（真目录、非软链，
排除台账记过的"容器软链 TMPDIR 自己缩小分母"那一形）；`go env TMPDIR/GOTMPDIR` 皆空＝走默认 `/tmp`。
`GOPROXY` 本程未复核（未拉任何依赖，离线跑通即无需）。

### 0.3 门一：`sh scripts/d22scan.sh` —— **本程现跑 rc=0**，且它自己带正控

```
$ date "+%H:%M:%S" ; sh scripts/d22scan.sh ; echo rc=$?
11:08:56
（前段＝扫描器自身测试套件）runtests.sh: OK - packages=[./...] top-level: PASS=29 FAIL=0 SKIP=0, === RUN=69, '[no tests to run]'=0
（后段＝真扫）
d22scan: examined 225 production Go files under internal/ and cmd/
d22scan: scope bans #1-5 internal/      examined 203 production Go files
d22scan: scope bans #1-5 cmd/           examined  22 production Go files
d22scan: scope ban #6 frontend/          examined  41 text files
d22scan: scope ban #7 internal/tools/    examined  18 production Go files
d22scan: scope ban #8 design/            examined  32 text files
d22scan: scope ban #8 frontend/          examined  41 text files
d22scan: scope ban #8 internal/         examined 407 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/               examined  39 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations
d22scan.sh rc=0            ← 锚点 598620e（11:09 前后）
```

**这一格有两个要点，本程都不按"rc=0 所以干净"收下**：

1. **这一发的正控在同一份输出里**：脚本第一步是扫描器自己的端到端套件，其中
   `TestBuiltBinaryGoesRedEndToEnd/seeded_violation_exits_1`、`…/fully_live_fixture_exits_0`、
   `…/empty_live_scope_exits_2`、`…/armed_ban_6_goes_red_on_the_panel_violation_exits_1` 全部 `--- PASS`
   ⇒ **"造出来的违规会让二进制红"这一味是这一步今天当场验过的**，不是本程的推理。
   （`--- PASS` 那 7 行原文在 `logs/`，本程贴了名字。）
2. ⚠ **poll 程 §0 报的"`scripts/d22scan.sh` 跑不动（`scan_test.go:1197` build failed）"在本程手里已经不成立**
   ——本程 rc=0。台账 `A219⑤` 早就把那句判成"别家写到一半的瞬时态被撞见，不是坏物"，本程是第三次撞上同一枚事实。
   ⇒ **对 `AC#15` 的裁没有任何一支依赖"那枚门跑不动"**；本程交的是**跑通那一发**。
3. ⚠ 分母 `internal/=407` **含兄弟程在飞的 `internal/panel/l2_grant_boundary_test.go`**（工作树里它是 `M` 态）。
   按台账 `A220`/`A222⑤` 的口径：**基线不当常量追，引用时现跑并带 HEAD** ⇒ 本程这句只属于 `598620e` 这一版，
   **不追它后面变成几**。

**⛔ 一条本程不overclaim 的话**：ban #1–5 的分母是 **production Go files only**（本程读到的就是"203 production Go files"这句），
`_test.go` **不在它的射程里**；只有 ban #8 明写"comments and `_test.go` included"。
⇒ **不许说"如果 AC#15 那批改出了裸 `go func(` 或墙钟超时，门会抓到"**——按盘上射程，**门抓不到测试文件里的那些形状**。
墙钟那一支（ban #4）只认 `.Sub(time.Now())` 一条拼写。本程**没有**把任何一枚门打成红来证明这条射程（见 §7 第 5 条），
所以这句话的强度只到"扫描器自己印出来的分母行"，不到"实验"。

### 0.4 门二：`gofmt` / `go vet`（本包）

```
$ gofmt -l internal/observe/      （空输出）     rc=0
$ go vet   ./internal/observe/                  rc=0
$ date "+%H:%M:%S"   11:09:22      $ git rev-parse HEAD   598620e…
```

⚠ 这两枚的"空输出"本程**不做零发现用**——它们只证形状与编译，不证判据；判据面的正控在 §1.4 与 §4。

---

## 1. 四数＋名册：两平台各自现跑（不引归档那两枚读数）

命令与派单/票面逐字同形，未改动的树、未摘名册成员、无 `-run`、无 `-cpu`：

```
Windows（工作树，只读地跑）  $ go test ./internal/observe/ -count=2 -v   11:06:39→11:06:46  rc=0
                             日志 D:\tmp\observe-ac15-final\logs\win-count2.v.log（包时 ok … 6.904s）
Linux（容器 /repo 只读挂载） $ go test ./internal/observe/ -count=2 -v   03:08:14→03:08:27 UTC  rc=0
                             日志 /scratch/final-logs/linux-count2.v.log（＝D:\tmp\observe-ac15-linux-r1-s1\final-logs\，包时 ok … 9.849s）
```

### 1.1 四数与两枚"生死"计数（两平台各一发，本程自己跑）

| 读数（grep 式子逐字同形） | Windows | Linux |
|---|---|---|
| `^=== RUN` | **142** | **142** |
| `^--- PASS`（顶层） | **142** | **142** |
| `^--- FAIL` | **0** | **0** |
| `^--- SKIP` | **0** | **0** |
| `^panic:` | **0** | **0** |
| 包级 rc | 0 | 0 |
| 不同名枚数 | **71** | **71** |
| 声明数（`grep -h '^func Test' *_test.go \| wc -l`） | **71** | **71** |
| `=== RUN` 去重名 | **71** | **71** |

⇒ **两个 OS 上 142＝71×2 闭合、`RUN` 去重＝`PASS` 去重＝声明数＝名册三者同值**（后两条是"有没有读数被吞"的正向证据，
只看四数看不到，见 §1.3）。

### 1.2 名册**逐名**两向差集（不是只比枚数）

```
$ comm -23 roster-win.txt roster-linux.txt   →（空）
$ comm -13 roster-win.txt roster-linux.txt   →（空）
$ diff -q roster-win.txt roster-linux.txt    → 无输出（rc=0）
$ wc -l 两份                                  71 / 71
```

⇒ **Windows 与 Linux 的名册是同一枚 71 枚，逐名同值**（不是"两枚碰巧相等的枚数"）。

### 1.3 "`*_other_test.go` / `//go:build !windows` 那一族在 Linux 有分母"这一问的现量答复

派单的预期是"Linux 名册**可能合法地更大**，而这是那些腿唯一有分母的地方"。**本包现量：不成立，且理由是结构性的。**

```
$ grep -rn "go:build" internal/observe/*.go
internal/observe/treemetrics_other.go:1://go:build !windows        ← 全包唯一一枚平台约束文件
$ ls internal/observe/ | grep -i other
treemetrics_other.go                                               ← 只有 .go，没有任何 *_other_test.go
$ grep -c '^func Test' 于该文件                                    0（它里面零枚用例）
```

⇒ `internal/observe` 在 Linux 上**没有一枚 Windows 不跑的用例**；那枚 `!windows` 文件只让本包在 ubuntu 上**编译得过**。
⇒ **"两向差集"这一味本程给的是"按定义不适用"而不是"核对后为空"**：两 OS 名册同一枚 71，
Linux 侧不存在额外分母。**这一格不能写成"差集已核对＝守恒"**（那会把"没有可比对象"说成"核过没问题"）。
⇒ 派单说"这是那些腿唯一有分母的地方"——**在本包是真命题但空集**：全仓那些腿（若有）在别的包里，本包零枚。

### 1.4 ⚠ 三枚"零"的**正控**（本仓这一周被骗过四次的就是这一形）

**"0 红／0 吞"这三枚零读数按字面跑在真码上会毫无意义**，除非同一把尺在已知正例上响过。
本程在同一台机、同一条 grep 式子上做了三发正控：

| 零读数 | 正控怎么做 | 读数 |
|---|---|---|
| `^--- FAIL` ＝ 0 | 仓外最小模块 `D:\tmp\observe-ac15-final\panicctl\`（一枚故意 `xs[0]` 越界的用例＋一枚常绿兄弟），同一条 `go test -v`、同一条 `grep -c '^--- FAIL'` | **FAIL＝1**（尺响） |
| `^panic:` ＝ 0 | 同一发正控里同一条 `grep -c '^panic:'` | **1**，逐字 `panic: runtime error: index out of range [0] with length 0 [recovered, repanicked]` |
| "名册会缩小＝读数被吞"这一味 | 同一发正控：声明数 **2**、`=== RUN` 去重 **2**、`^--- ` 判定行去重 **1** | **2→1＝吞**，机器逮到了 |

正控原文（同一条计数式子，未改一字）：

```
$ cd /d/tmp/observe-ac15-final/panicctl && go test ./... -v -count=1 > logs/panic-control.log ; echo rc=$?
rc=1
RUN=2  PASS=0  FAIL=1  SKIP=0  ^panic:=1
control: declared=2  rosterPASS=1  RUNnames=2
```

⇒ **同一把尺在正例上四处全响** ⇒ 本程上面 §1.1 那三枚 0 与"声明＝RUN＝PASS 三账闭合"是有牙的读数。
⇒ 顺带把派单那句**机制**验到实物：一包里一枚 panic 之后，**同包其余用例连 `=== RUN` 都不再有**
（正控里那枚常绿兄弟 `TestStillGreenSibling` 在日志里零痕），所以包级 rc 只说"这包失败"、不说谁被吞——
**本程因此比的是名册，不只是计数**（§1.2）。

### 1.5 `t.Skip` 与 `SKIP=0` 的形状（Linux 侧现量，Windows 侧同码同值）

```
$ grep -rn "t\.Skip" internal/observe/
internal/observe/sampler_test.go:36            // …t.Skip and not a silent TreeMetrics{} return…   ← 注释
internal/observe/sampler_test.go:357           // …never t.Skip, never a silent return.            ← 注释
internal/observe/window_wait_136_test.go:94    // …no time.Sleep papering over a window, and no t.Skip: ← 注释
```

⇒ **可执行位 0 枚**（三枚全在注释行）⇒ `SKIP=0` 不是"跳过了没看见"，本包**结构上产不出 SKIP**。
⇒ ⚠ 这条 grep 的注释命中枚数在**两 OS 都是 3**（同一版码，§0.1 的字节账）。

---

## 2. `AC#15` 的判据**按票面现在怎么写**（逐行引，不引派单摘要）

票面＝`.scratch/wisp/issues/136-the-zero-sample-fail-closed-guard-underneath-the-freshness-nail-has-no-nail-54-tests-stay-green-without-it.md`。
本程现量：**该格在 `:348`，仍是 `- [ ]`**（`grep -n 'AC#15'` ⇒ `:348` 为那一行）。

### 2.1 派单点名要我读的三枚 `>` 块（`:274` / `:286` / `:288`，逐字）

`:274`（**地界句的归属更正**——它改变"这一格能不能动生产码"的答案）：

> **先认一枚**：`:250-251` 那句"派单时由编排者按'这条判据要落盘必须碰哪几段'现划"——我上一版派单简报里把 **AC#15 的 `:278` 地界句（"只许动 `internal/observe/**_test.go`"）当成了 AC#14 的**。

⇒ 现量：AC#15 的地界句**今天不在 `:278`**，在 **`:361`**（原格内）与 **`:379`**（重划块 ⑥，"地界不变"），
票面自己那句"AC#15 只许动测试"这一**约束属于 AC#15**、**不属于 AC#14**（AC#14 另有具名解冻）。
⇒ 派单简报里"AC#15 的射程句被我当成 AC#14 的"这一支**在盘上确认成立**。

`:286`（**串行约束**，两格共用一枚文件）：

> **③ 与 AC#15 串行，不并行**：两格共用 `internal/observe/sampler_settle_coverage_136_test.go`——AC#14 判据② 要的正是**推翻**那里 `:208-210` 那句 `if !rep.Pass { t.Fatalf("disclosure leg, not a verdict leg: this window still passes, ...") }`（它现在逐字要求"丢一半读数那一形**仍然 pass**"，与 AC#14 方向相反）。

`:288`（**行号必须现量重划**）：

> AC#15（那里 `:189` 的 `precondition broken: only %d reads taken` 偶发红）**排在 AC#14 之后**，且 `AC#14` 落地后它的靶子形状可能已经变了 ⇒ 到时要现量重划，不许照抄今天的行号。

⇒ 本程因此**没有**用任何一枚派单/票面给的行号，全部现量（§2.3）。

### 2.2 本格判据按哪一版执行（票面自己声明的优先级）

`:372` 那句是票面自设的替换令：

> **AC#15 的靶形已现量重划（09-24 21:3x 编排者…）**：…**上面原句一字不抹，本格判据按下面这版执行**

⇒ 所以** operative 判据＝`:348-363` 原格 ＋ `:365-370` 更正 (a)-(d) ＋ `:372-380` 重划块（①②③④⑤⑥，其中 `④` 改写 `①`）＋ `:383-388` 的再更正（①／①b／⑤＋那枚仪器缺陷）**。
逐条落到本程手里的尺：

| # | 票面原文（去引号，位置） | 本程怎么量 |
|---|---|---|
| 判据①（被 `:377` 的 `④` **改写**） | "命中率必须**分两口径各给一张表**，永不加总——**(A) 外循环 `go test -count=1`**…**(B) 同进程 `-count=N`**…两口径的 RUN／顶层 PASS／FAIL／SKIP 四数与**名册双向差集**都要落表" | §1（两平台四数＋名册两向差集）＋ §6 逐格裁两程交件 |
| 判据② | "根因归到**取读数的窗口/间隔与断言之间缺的那条判据**，⛔ 不许收在'机器负载高'" | §3.1（缺的那条判据＝"数了但不等"，代码级定位） |
| 判据③ | "修法只许把'等待'变成**有判据的等待**（轮询到条件成立＋单调超时上界）；⛔ `time.Sleep` 糊窗口／`Skip` 这条用例／放宽或调高任何阈值／删断言" | §3.2＋§3.3（机械读数） |
| 判据④ | "改完同一命令复跑 **≥30 发命中必须 0**…并给**名册双向差集为空**（证明 0 命中不是'用例变少／被跳过／被 panic 吞掉'）" | §6 裁两程；本程 §1.1 那一发＝71 名 ×2 |
| 地界（`:361`＋`:379`） | "只许动 `internal/observe/**_test.go`；`sampler.go` 的 `CheckSettle`…语义若必须变 ⇒ 停手报回（属人工批准面）" | §3.4 `git diff --numstat` 逐枚 |
| 方向锁（`:379` 末） | "⚠ 现在多一层：**AC#14 的门已经 true**，所以任何'让那一窗重新 pass'的改法都会被门行的折叠逻辑判成假——**修法方向只许是'多等一会儿并且等到'，不许是'改判定为不算丢'**" | §3.5：现量 `sampler.go:556 const settleCoverageRowGates = true` 确认这一层生效，再逐条对修法形状 |
| 窄口径尺（`:385`） | "**窄口径**＝真正'会因窗口没读满而红'那一族 ⇒ **8 处／3 枚文件**…派单第 0 步的分类尺**用窄口径这 8 处**" | §2.4 现量：交付态已是 **13 个 await 站点／12 枚腿／5 枚文件**（票面那 8 处是 `a1fd5bf` 版、且**单位是守卫不是窗**） |
| 前置一发（`:387`） | "`sampler_settle_gate_136_test.go:262` 是 `buildSettleVerdicts(SettleReport{})[0]` 直取下标、前一行没有长度守卫…（**列入 AC#15 真跑程的第 0 步之前，作为独立一发**）" | §2.5 现量：已由独立 commit `021a549` 落（且**是 `d8390aa` 的祖先**）⇒ 前置成立 |
| 上界算术（`:378` ⑤） | "超时上界要由这 400ms／200ms 的算术推出来，不是拍脑袋" | 本文 §4（派单 Job A 那一问的正面答复） |

### 2.3 靶子的**当前**形状（本程现量；票面 `:189` 那一族今天在哪几行）

```
$ grep -n "^func Test" internal/observe/sampler_settle_coverage_136_test.go internal/observe/sampler_settle_gate_136_test.go
coverage:142 TestCheckSettleSingleTrustworthyReadReportsItsLoss
coverage:236 TestCheckSettleHalfTheReadsFailedReportsItsLoss      ← 票面点名的靶子腿（原句说 `:189`，盘上现量 `:236`）
coverage:320 TestCheckSettleZeroFootprintDropsAreCountedToo
coverage:373 TestCheckSettleFullyMeasuredWindowReportsNoLoss
gate:160    TestSettleCoverageRowExistsAndPassesWhenFullyMeasured
gate:212    TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed
gate:257    TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured   ← §4 那枚"最长腿"
```

靶子腿内部（同一枚文件现量）：await 站点 `:237`（`awaitSettleReads(t, 4, …)`）、
少判据那条守卫 `:243`（`if reads < 4 { t.Fatalf("precondition broken: only %d reads taken, half-and-half needs a window to lose in", reads) }`）、
影子守卫 `:246`（`kept < 2`）。窗口参数（`:112`，逐字）＝`CheckSettle(…, 100*time.Millisecond, 10*time.Millisecond, …)`。
⇒ **票面 `:350` 那句 "`:189` 那一族" 与 `:373` ① 那整张行号表都已腐坏**（AC#14 落地＋本程被裁的这批落地各推了一段），
本程按 `:288` 那句"现量重划"给上面这版。

### 2.4 窄口径家族在交付态的真实枚数（本程自跑）

```
$ grep -rnE 'await(Settle|State)Reads\(t, ' internal/observe/ | grep -v window_wait
coverage:143 want=3   coverage:237 want=4   coverage:321 want=3   coverage:374 want=3
gate:166 want=1       gate:217 want=2       gate:264 want=1       gate:267 want=1
settle_zerosample:43 want=1  settle_zerosample:140 want=1
sampler_test:93 want=5(state)  sampler_test:335 want=3(state)
zerosample:152 want=3(state)
⇒ 13 个 await 站点 / 去重 12 枚腿 / 5 枚文件（gate:264＋:267 同属一枚腿＝"两枚窗"）
```

⇒ 与 poll 程 §7.7 那张**更正后**的表逐枚同值、与 Linux 程 §3.6 末那句"站点行号…与 poll 程 §7.7 逐枚相同"同值；
⇒ 与票面 `:385` 的"窄口径 8 处／3 枚文件"**枚数不同**——那是 `a1fd5bf` 版的**守卫**枚数，本程数的是**窗/站点**，
且中间隔了 AC#14 与本批两次落地。**两个数不打架，单位不同，引哪个都要带单位**（本仓"枚数要带口径"那条口径）。

### 2.5 `:387` 那枚"排在第 0 步之前"的前置：现量已闭合

```
$ sed -n '300,306p' internal/observe/sampler_settle_gate_136_test.go
	neverRows := buildSettleVerdicts(SettleReport{})
	if len(neverRows) != 1 {
		t.Fatalf("precondition broken: buildSettleVerdicts must return exactly 1 self-describing row for a never-measured report, got %d rows: %+v", len(neverRows), neverRows)
	}
	never := neverRows[0]
$ git merge-base --is-ancestor 021a549 d8390aa   → 真（021a549 是 d8390aa 的祖先，09-24 22:46）
$ git show --numstat 021a549
  277 0  docs/evidence/s1/136-instr-fixes-r1.md
   25 1  internal/observe/sampler_settle_gate_136_test.go      ← 纯测试，且早于本批修法
```

⇒ 票面要求的"先修那枚会 panic 的仪器、**只给它加守卫让它自己红**、不许改 Skip"这一前置**已按字面闭合**
（守卫＋`t.Fatalf`＋无 Skip；且它不在被裁的两枚交付里，是另一程 `136-instr-fixes-r1`）。

### 2.6 门行已 `true` 这一层（`:379` 方向锁的前提，本程现量）

```
$ grep -n "settleCoverageRowGates" internal/observe/sampler.go
556:const settleCoverageRowGates = true      580:  Gate: settleCoverageRowGates,
```

⇒ `AC#15` 裁的时候**必须**带着这一层：AC#14 落地后"丢一半读数"那一形**不再 pass**，
所以**任何"把判据改成不算丢"的修法都会被门判成假**——本程据此把判据③的方向锁当成**硬判据**用（§3.5）。

---

## 3. 判放水：只用这两条尺（断言有没有被弱化／helper 是原有的还是为了能过而新造的）

### 3.1 尺之外先钉一句：本程判的是**交付态码**，不是两程的报告

下面每一枚读数都由本程自己在 `ee5a25e`（改前）与交付态（`333dfe3`＝工作树，§0.1 字节账）之间跑出来。

### 3.2 机械面：**条件行**逐行 diff（式子自定，不引 poll 的式子）

```
$ 每枚文件：grep -cE '^[[:space:]]+if '  改前 / 改后
coverage            52 → 51     gate               40 → 40     settle_zerosample 19 → 19
sampler_test        47 → 47     zerosample         25 → 25     window_wait        0 → 3（新文件）
$ t.Fatal* 语句枚数（grep -coE 't\.Fatal'）
coverage 49 → 48   gate 37 → 37   settle_zerosample 20 → 20   sampler_test 37 → 37   zerosample 23 → 23   window_wait 0 → 1
⇒ 全包条件行 183 → 182；t.Fatal* 166 → 167（净 +1 在 window_wait 那枚 bound 红）
```

**少掉的那一枚条件与那一枚 `t.Fatal` 本程逐枚追到了身份**（"少了＝可能被弱化"，所以必须点名）：

```
改前 if err != nil 站点：coverage:100 / :111 / :209      改后：coverage:113 / :130
改后 :93 settleSUT(...) 整个函数体＝ return settleTreeSUT(t, &coverageScriptTree{steps})
改后 :103-122 settleTreeSUT(...) 内 CheckSettle(...) 之后 :113 查 err
```

⇒ **身份＝两条开窗路径合并成一枚 `settleTreeSUT`**：靶子腿原来自带一段 `CheckSettle` ＋ `if err != nil`，
现在改叫 `settleTreeSUT`，**每条开窗路径仍各查一枚 err**（脚本树路径 `settleSUT → settleTreeSUT`、
状态树路径直接 `settleTreeSUT`），**没有哪条路径失去 err 检查** ⇒ 判**不是弱化**。

**条件文本两向 diff（去缩进后逐行）**：改前后唯一差异两枚，全在 coverage——

```
21,22c21
< if tree.reads < 4 {                 > if reads < 4 {
< if err != nil {                     （并入 settleTreeSUT，见上）
25c24
< if kept+lost != tree.reads {        > if kept+lost != reads {
```

⇒ 两枚都是**同一个量的名字换了**（`tree.reads` → 等待带回的那一窗的 `reads`），**阈值一字未动**。

**阈值面单独再来一把尺**（本程自己的式子，比 poll 程那把宽：把 `reads/kept/lost/samples/len(rep.Samples)` 与
数字的比较全抓出来、排序、取 md5）：

```
sampler_settle_coverage_136_test.go       before=c5de18f3  after=c5de18f3  IDENTICAL
sampler_settle_gate_136_test.go           before=132f8828  after=132f8828  IDENTICAL
sampler_settle_zerosample_136_test.go     before=8f24384e  after=8f24384e  IDENTICAL
sampler_test.go                           before=e8dde583  after=e8dde583  IDENTICAL
sampler_zerosample_136_test.go            before=dc6797a0  after=dc6797a0  IDENTICAL
⇒ 五枚被改文件里"变量 比较 常数"这一族的**全集**改前后逐字节同值：< 4 / < 3 / < 2 / == 0 / != reads 一枚没动。
```

### 3.3 第二问：helper 是原有的还是为了能过而新造的

| 问 | 本程现量 |
|---|---|
| 有没有既有的有界等待可复用？ | **有，且被复用了**：`awaitWindow` 用的 `NewTimeout/Expired/Remaining` 是本包**生产级**既有仪器（`clock.go`，`goroutine.go:139` 生产在用）。本程现量 `grep -n 'tm := NewTimeout' internal/observe/*.go` ⇒ `window_wait_136_test.go:99` 与既有的 `goroutine_test.go:57` 同形 ⇒ **不是新造第三种超时写法** |
| 新造的那枚仪器是不是"为了能过而造"？ | 它只做三件事：再开一整窗、比较**判据**、到期 `t.Fatalf`。**不读产品侧任何字段**：判据是 `settleWindow.reads` / `stateWindow.reads`，而 `reads` 是 **fixture 自己的计数器**，在 `ReadTree()` **入口第一行**自增（`coverage` 侧 `coverageScriptTree`/`alternatingTree`、`gate` 侧 `gateScriptTree.ReadTree:84 t.reads++`、`gateAlternatingTree.ReadTree:99`）——**先计数、后判错误/后判零足迹** ⇒ 产品把读数丢了、把披露打错了、把门行判错了，**都满足不了这枚等待** ⇒ 判**不是为洗绿而造** |
| 等待是谁结束的？ | **只由那一窗自己的枚数判据结束**（`if n >= want { return w }`）。到期那一支不是"睡够了就算过"，是 `t.Fatalf` ⇒ 见 §3.4 的反面情形 |
| ⛔ 墙钟结束？ | **无**：`awaitWindow` 体内零枚 `time.Now()`／`time.Since`／`Sleep`（本程现量：`sed -n '97,122p'` 全文只有 `NewTimeout`/`Expired`）。⇒ 但**这一条门管不到**：ban #4 只认 `.Sub(time.Now())` 一条拼写，本程**没把它打成红**去证明（§7 第 5 条） |
| ⛔ `t.Skip`／`Sleep` 糊窗？ | `t.Skip` 可执行位 0 枚（§1.5）；`time.Sleep` 于本包：`grep -rn 'time.Sleep' internal/observe/*_test.go` ⇒ 见 §3.6 |

**结论（判放水这一格）**：**不放水，是收紧**。三条支撑，全是读数不是话：
① 阈值集合 md5 五枚文件逐枚同值；② 少掉的那一枚条件有明确身份（开窗路径合并）且 err 检查数与路径数仍相等；
③ 判据取**接缝计数器**、且计数器在**产品分支之前**自增 ⇒ "让等待被产品坏掉满足"这条路结构上不存在。
⚠ 本程**没有**替 poll 程复算它那五发变异（§7 第 2 条）；③ 是本程从代码字节读出的**结构**论证，
不是"变异仍在咬"的**实验**论证——那两味分开记。

### 3.4 反面形状：**"以墙钟结束"那一支确实存在，但它红而不是过**

`awaitWindow` 里唯一一处"时间到了"的处理是：

```
if tm.Expired() { t.Fatalf("precondition broken: %d windows opened inside the %v monotonic bound (clock.go Timeout) all fell short of %d; counts seen: %v%s", …) }
```

⇒ 预算耗尽的出口是**红**（印尝试枚数＋每窗实收枚数），**不是 return、不是绿** ⇒ 派单点名要防的
"等一下就当等到了"这一形在代码上不存在。本程在 §4 里**把这枚红真叫出来了**（含最长那一枚腿）。

### 3.5 票面 `:379` 那把方向锁：交付形状落在哪一支

现量：`settleCoverageRowGates = true`（`:556`）＋ 五枚被改文件的阈值集合**一字未动**（§3.2）＋
新增的全是"再开一窗"（13 站点）＋**没有新增任何"这一形不算丢"的分支**
（`grep -cE 'SampleErrors|sample_errors' 于改动行`：本程比对 `git diff ee5a25e..HEAD -- internal/observe/*_test.go` 里
新增的 `+` 行**没有一枚**去掉或绕过披露断言）。
⇒ 方向＝票面允许的那一支：**"多等一会儿并且等到"**；不是被禁的那一支："改判定为不算丢"。

### 3.6 `time.Sleep` 与 `t.Parallel`（两枚"糊窗口"的廉价出口，本包现量）

```
$ grep -rn "time.Sleep" internal/observe/*_test.go
clock_test.go:23:        time.Sleep(90 * time.Millisecond)          ← 不在本批射程（本批未碰这枚文件）
goroutine_test.go:47/:94/:200:  time.Sleep(2~5 * time.Millisecond)  ← 同上，AC#11 那枚文件
window_wait_136_test.go:93:// …no time.Sleep papering over a window…  ← 注释行，非语句
⇒ 全包命中 5 处：4 处语句全在**本批未碰**的两枚文件里，1 处是注释。

$ 逐枚（六枚被裁文件，可执行位；parallel 同式）
window_wait sleep=0(注释 1)  coverage sleep=0  gate sleep=0  settle_zerosample sleep=0  sampler_test sleep=0  zerosample sleep=0
$ 同一批式子打在改前 ee5a25e 的五枚文件上：sleep 全 0、parallel 全 0
$ t.Parallel：六枚被裁文件 0 枚，改前亦 0 枚
```

⇒ **判据③ ⛔ 清单里"`time.Sleep` 糊窗口"那一味在本批交付里零枚**（连"注释之外新写一枚 sleep"都没有）；
⇒ `t.Parallel` 0 枚 ⇒ **名册内不存在并发用例**，负载标签可直接归因（与两程同形）。
⇒ ⚠ 那 4 处既有 sleep 属 `clock_test.go`／`goroutine_test.go`，**不属本批也不在本格射程**——
本程点名是为了下一位不把"包里有 sleep"错记成"这批留的口子"。

---

### 3.7 契约轴（派单点名的"不许"清单，逐枚现量）

| 不许 | 本程读数 | 判定 |
|---|---|---|
| 加 `t.Skip` | 可执行位 0 枚（§1.5，两 OS 同码）；票面 `SKIP=0` | 未犯 |
| 调高／放宽任何阈值 | 五枚文件"变量-比较-常数"全集 md5 逐枚同值（§3.2）；`want` 由 13 个调用点传入，取值 1/2/3/4/5 | 未犯 |
| 删断言 | 条件行 183→182、`t.Fatal*` 166→167；少的那一枚有明确身份且 err 检查数与开窗路径数相等（§3.2） | 未犯 |
| 加豁免名单／allowlist | `git diff --name-only ee5a25e..HEAD` 里**无** `tools/d22scan/allowlist.txt`；本批只碰 6 枚 `internal/observe/*_test.go` ＋它自己的证据件 | 未犯 |
| 动 SLO 阈值／golden | `git diff --numstat ee5a25e..HEAD -- internal/observe/thresholds.go` **空输出**；`git diff --name-only ee5a25e..HEAD \| grep -iE 'golden\|testdata\|\.sse$'` **空输出**（仓里被跟踪的 golden 名册本程现量 58 枚） | 未犯 |
| 动生产码 | 见 §3.8 | 未犯 |

### 3.8 地界／领土主张**独立复核**（不引 poll 那句）

```
$ git show --numstat d8390aa
  218   0 docs/evidence/s1/observe-ac15-poll-r1.md
   59  24 internal/observe/sampler_settle_coverage_136_test.go
   40   8 internal/observe/sampler_settle_gate_136_test.go
   52  32 internal/observe/sampler_settle_zerosample_136_test.go
   41  19 internal/observe/sampler_test.go
   24   8 internal/observe/sampler_zerosample_136_test.go
  142   0 internal/observe/window_wait_136_test.go
$ git show --numstat 333dfe3
   71  0 docs/evidence/s1/observe-ac15-poll-r1.md
   12  2 internal/observe/window_wait_136_test.go        ← 红句列表截到 8 项，纯消息文本
$ git log --oneline ee5a25e..HEAD -- internal/observe/sampler.go
（无输出）      ⇒ 生产侧 settle 读循环在本批一次没被碰
$ git diff --name-only ee5a25e..HEAD | grep -v '^docs/' | grep -v '_test.go$'
design/doubao/**（前端会话）  tools/d22scan/{gitignore,main,scan_test}.go（扫描器收尾程）
⇒ 范围内**非测试、非别家地界**的 Go 文件＝**零枚**；`internal/panel/l2_grant_boundary_test.go` 是兄弟程在飞的
   **未提交**改动（工作树 `M`），本程**没把它算进任何一发读数**，也没跑 `./internal/panel/`。
```

⇒ **poll 程那句"生产码零字节"成立**，且本程是按**两枚 commit 各自的 numstat** 加的，不是按它那句自述。

---

## 4. 那一问：最长腿（两枚 200 ms 窗＝400 ms）的真实预算——**本程现量，不做推理**

### 4.0 问题原文（从 `observe-ac15-linux-r1.md` §3.5 逐字引；它承诺答而没答的就是这一问）

> **要登记的余量（本程给的算术，用的全是 Linux 现量）**：家族里 nominal 最长的腿
> `TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured` 是**两枚 200ms 窗＝400ms**（§3.6 的净窗家族表），
> 它的预算只有 **2s/400ms ≈ 5 枚窗**；靶子腿是 **2s/100ms＝20 枚**。
> ⇒ "预算≈20 次中位尝试"这句在**家族里不均匀**——那枚 400ms 腿只有 5 次、两枚 200ms 腿只有 10 次。
> …§4 落地后在 §5 给"要不要按最长腿／按 OS 加余量"的答复。

派单把这一问交给本程的开放半边是：**能不能真把那枚 400 ms 腿饿到 5 枚窗以外，饿到的时候失败长什么样。**
本程的答案全部出自**仓外副本**（`git archive` 到 `/scratch/final-r1/` 与 `D:\tmp\observe-ac15-final\`），
**交付码一字未动**（§0.1 的字节账在每次读数前后都成立）。

### 4.1 第一发（结构下界）：**一枚 `CheckSettle` 窗交不回 0 枚读**

`sampler.go:495-525` 的环是 `select{<-ctx.Done(); <-t.C}` → `ReadTree()` → 之后才 `if time.Now().After(deadline) { break }`
⇒ **第一枚 tick 一定会被等、等到了就一定读一次**。本程不去读码定罪，把它**打到盘上**：
把 tick 拉长到**比窗还长**（`within=200ms` 不动、`interval` 10ms→250ms），此时那一枚 `want=1` 的判据**该不该饿**？

```
Linux （t-floor）
--- PASS: TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured (0.50s)
Windows（tw-floor，同一改动，3 发）
--- PASS: … (0.52s)   --- PASS: … (0.60s)   --- PASS: … (0.51s)
```

⇒ **两 OS 都不饿**：腿时＝2×250ms＝0.5s（**两枚窗各只跑了一发，零重开**），判据 `want=1` 在"整个窗期只交付一次读"的极端下**仍成立**。
⇒ 同形的那一枚 `SampleState` 腿（`sampler_test.go:341`，`want=3`＝环外 2 枚＋tick 1 枚）本程照同一招量：
`interval` 10ms→100ms 而 `duration` 留 30ms ⇒ **5 发全 PASS、每发 0.10s（无重开）**。
⇒ **结构结论（现量支持，非推理）**：`CheckSettle` 窗的下界＝1 枚读、`SampleState` 窗的下界＝3 枚读
（`sampler.go:267` 开窗读 ＋ 第一枚 tick 读 ＋ 收窗读）。
**⇒ 派单问的那枚"最长腿"用的是 `want=1` ⇒ 它是家族里最不可能饿的一枚，不是最容易饿的一枚。**

### 4.2 第二发（真实窗长下的预算）：**2 s 放得下 10 枚 200 ms 窗，两 OS 同值——不是 5**

判据既然饿不到，本程就**把判据换成不可达**（`gate:264` 的 `want` 1→9999，窗长仍 200 ms／10 ms 一字未动），
直接量"这一枚 await 在 2 s 里重开了几整窗"：

```
Linux （t-200ms）
    sampler_settle_gate_136_test.go:264: precondition broken: 10 windows opened inside the 2s monotonic
      bound (clock.go Timeout) all fell short of 9999; counts seen: [20 20 20 20 20 20 20 20] (+2 more, all of them short of 9999)
--- FAIL: TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured (2.01s)

Windows（tw-200ms，两发；此刻容器内 12 枚自旋仍在跑）
    sampler_settle_gate_136_test.go:264: precondition broken: 10 windows opened inside the 2s monotonic bound
      … counts seen: [2 11 14 10 13 13 14 16] (+2 more, all of them short of 9999)
--- FAIL: … (2.11s)
    …  counts seen: [18 12 9 20 15 13 20 12] (+2 more …)
--- FAIL: … (2.17s)
```

⇒ **两 OS 各 10 枚整窗／2 s**（Linux 每窗 20 枚读；Windows 每窗 2–20 枚读——那一枚 `n=2` 的窗就是被容器自旋压出来的）。
⇒ **那一问的前提"2 s/400 ms ≈ 5 枚窗"在两 OS 都不成立**：`awaitWindow` 每次调用**新起一枚** `NewTimeout`
（`window_wait_136_test.go:99`），**预算是"每一枚窗 2 s"，不是"每一条腿 2 s"**。
⇒ 而且红**印在第一枚 await 的站点**（`gate:264`）——第二枚窗（`:267`）**根本没轮到**：
一条腿只有在"第一枚窗成功、第二枚窗饿死"这一串里才会花到 2 s＋2 s，本程没造出那一串（§7 第 4 条）。

### 4.3 第三发（真负载能不能饿到它）：**不能，量到的是 500 发零重开**

| 发 | 命令（逐字） | 负载标签 | 读数 |
|---|---|---|---|
| Linux | `go test ./internal/observe/ -count=300 -cpu=1 -v -run '^TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured$'` | 容器内 12 枚纯 CPU 自旋、无 `--cpus` 限额；批前 `3.65`、批后 `12.85`（1min 档）；测试压到 1 枚 P | **RUN 300／PASS 300／FAIL 0／SKIP 0／`^panic:` 0**；腿时 p50 **0.40s** p99 0.42s **max 0.43s**；包时 121.564s ⇒ **重开 0 枚**（重开会多烧一整枚 200 ms 窗＝≥0.60s） |
| Windows | 同命令，同容器自旋在跑（`/proc/loadavg` 13.58 那段），`-cpu=1`，`-count=200` | 本机 6C12T 被容器拿走 12 枚线程 | 未加仪器那一版：p50 0.42／p90 0.46／p99 0.61／max 0.68，**3 发 ≥0.60s**（＝重开的时长签名？） |
| Windows·仪器版 | `t-instr`：在 `awaitWindow` 每次尝试印一行 `AWAITCOUNT want= try= n=`（**只在仓外副本**） | 同上 | **400 行 await ＝ 200 发 × 2 枚窗，try 全 ＝1** ⇒ **那 3 发 ≥0.60s 不是重开，是慢窗**；每窗交付枚数 **min 2**、均值 13、max 20（判据要 1）⇒ FAIL 0 |

⇒ **"能不能真把那枚 400 ms 腿饿到 5 枚窗以外"的正面答复：本程造不出来**（在两台机、两种 OS、
可造的最重负载下 500 发**一枚重开都没有**），唯一造得出来的那一发是**把判据改成不可达**（§4.2），
而那已经不是"负载"而是"假题目"。
⇒ 失败形状（本程真叫出来过 13 发）逐字是：`<站点>: precondition broken: N windows opened inside the 2s monotonic bound
(clock.go Timeout) all fell short of <want>; counts seen: [前 8 枚每窗实收] (+M more, all of them short of <want>)`，
**红在 await 的调用点**（`t.Helper()` 归因），腿时 ~2.0–2.2 s，包 rc=1，**名册不缩**（§4.5）。

### 4.4 第四发（本程最强的一枚）：**同一台机、同一负载、同一命令——改前 16 枚红，改后 0 枚红**

派单要"每个 0 都要有正控"。本程的正控不是另造一发，而是**把被修的旧版本扔进同一副负载**：

| 树 | 命令（逐字同形） | 读数 |
|---|---|---|
| **`t-pre`**＝`git archive ee5a25e`（**改前**、无 `window_wait_136_test.go`） | `go test ./internal/observe/ -count=200 -cpu=1 -v -run '^TestCheckSettleHalfTheReadsFailedReportsItsLoss$'`，容器 12 自旋在跑（批时 `loadavg 13.58`） | **PASS 184／FAIL 16／SKIP 0／`^panic:` 0**，rc=1；红句逐字：`sampler_settle_coverage_136_test.go:214: precondition broken: only 3 reads taken, half-and-half needs a window to lose in`（另有 `only 2 reads taken`） |
| **交付码**（`t-instr`，判据未动＝`want=4`） | 同命令、同负载、同机 | **PASS 200／FAIL 0／SKIP 0／`^panic:` 0**；await 行 **208＝200 枚 try=1 ＋ 8 枚 try=2** ⇒ **8 次重开、8 发全被救回**；每窗枚数直方：**n=1 ×1、n=2 ×4、n=3 ×3**（＝那 8 枚低于 4 的窗）、n=4 ×5、n=5 ×5、n=6 ×15、n=7 ×17、n=8 ×20、n=9 ×37、**n=10 ×101** |
| 交付码·**净窗**（自旋已停，`loadavg` 回落） | 同命令 | await 行 **200＝200 枚 try=1、0 枚重开**，PASS 200／FAIL 0，包时 20.144s |

⇒ 三行合起来把**两程都没能给的**那格交出来了：
1. **这枚 flake 在本机可造**（不需要人造停顿——把容器自旋加上就够：改前 200 发**红 16 发＝8%**）；
   红因逐字就是归档与票面那一句（`only 2/3 reads taken`）。
2. **重开路径今天真被执行过**：poll 程 §6.3 与 Linux 程 §4.4 各自量的"每枚腿 p50＝max ⇒ 一次重开都没发生"
   在**这副负载下不成立**——本程量到 **8/200 发重开**、**每发恰好 1 次**（无 try≥3），且**一枚没红**。
   那一枚 **n=1** 的窗（100 ms 里只交付 1 枚读）在**改前**必红、在**改后**被下一枚窗救回。
3. **净窗仍 0 重开**（200 发）⇒ "改后全绿"不是"等待没执行"，两味本程分开量、分开报。

⚠ **口径与不合并**（票面判据①／`:377` ④）：上面是**口径 B（同进程 `-count=N`）·忙窗·Windows·`-cpu=1`**，
分母 **200 发／枚**，与 poll 程的 40 发口径 A、Linux 程的 4000 发口径 B **各记各的**，**不并成一枚百分比**。
⇒ 本程不写"率从 8% 降到 0"（同一副负载只跑过一次、且负载不是本程控制的实验变量，Linux 程 §4.6 那句限定同样适用于本程）。

### 4.5 第五发（预算的代价上限＝为什么不该按最长腿去缩它）

`t-blast`：把 **13 个 await 站点的判据全改成不可达**（`want`→9999），整包一发，Linux：

```
RUN=71  PASS=59  FAIL=12  SKIP=0  ^panic:=0     bound 红＝12 枚
包时 26.867s   （同一命令净窗基线＝9.849s，§1.1）
12 枚红名册＝家族那 12 枚腿，逐枚红在**自己第一枚 await 的站点**：
  coverage:143 :237 :321 :374   gate:166 :217 :264   settle_zerosample:43 :140
  sampler_test:93  :335          zerosample:152
（gate:267 零枚 ⇒ §4.2 那句"第二枚窗没轮到"在这发里是读数不是推断）
```

⇒ **最坏代价是量出来的**：整包被这枚上界顶到 **26.867s ＝ 净窗的 2.7 倍**，
而红名册**恰好**＝12 枚被包的腿、**一枚不多一枚不少**、`RUN` 仍 71、`SKIP` 仍 0、`panic` 仍 0。
⇒ 这一发同时是**本程计数仪器的第二枚正控**（同一条 grep 式子在真红上响 12 次、名册不缩）。
⇒ **对本问的意义**：上界若"按最长腿推导"（把 2 s 换成按 400 ms 计的更小值），
省下来的是**今天 0 概率那条路径**的代价上限，砍掉的却是**今天 8% 概率那条路径**的重开预算——
方向与判据③（"多等一会儿并且等到"）**相反**。

### 4.6 答复（一句话＋三行账）

> **一句话**：**这枚 2 s 上界不改、也不按最长腿缩、也不按 OS 加余量**——因为"2 s/400 ms≈5 枚窗"这个前提**在两 OS 上都被现量否掉**（同一枚 200 ms 窗实测放得下 **10 枚**，预算是**每枚 await** 而非**每条腿**），而那条"最长腿"用的是 `want=1`、`CheckSettle` 窗**结构上交不回 0 枚读**（本程把 tick 拉到比窗还长，两 OS 仍零重开、仍 PASS），**它恰是家族里最饿不到的一枚**。

要登记的三行账：
1. **不均匀是真命题、但方向反了**：每窗预算 = 2 s ÷ 窗长 ⇒ 100 ms 窗 20 枚、200 ms 窗 10 枚（§4.2 实测），
   而**判据高低与窗长成反比**（want=4/5 的腿都跑在 100/120 ms 窗上，want=1 的腿跑在 200 ms 窗上）
   ⇒ "最长腿预算最少"这句要换成"**最长腿的判据最低、且它有结构性下界 1 枚读**"。
2. **两味该修的措辞（都是文档/账面，不是码）**：
   (a) `window_wait_136_test.go:53-55` 那句 "2s is 10x the p99 of the slowest leg (200ms) … buys roughly 20 median-cost attempts"
   ——**20 只对 100 ms 那一族成立**，200 ms 那一族是 10 枚（本程实测两 OS 同值）；
   (b) Linux 程 §3.5 的 "2s/400 ms ≈ 5 枚窗"——把**每-await** 的预算当成了**每-腿**。
   本程**不改那两枚文件的字**（一枚是交付码、一枚是别程的证据件），只在此具名登记。
3. **OS 余量不加**：两 OS 的"窗数／2 s"逐枚同值（10 对 10），差别在**每窗交付几枚读**
   （Linux 净窗 20 枚、Windows 忙窗 2–20 枚）与**饿穿所需的停顿**（Linux ≥80 ms、Windows 60 ms 就红——归档两程的读数），
   这两味差别都被"重开一整窗"的形状吸收掉了，不需要在**常数**里再加 OS 分支；
   加一枚 OS 分支反而会把 §4.5 那枚"最坏代价 26.9 s"变成**每台机不同**的一个数。
   ⚠ 但**这一支只到"今天两枚 OS 各 10"为凭**，不等于"任意 runner 上够"（§7 第 3 条）。

### 4.7 本程自己复算的那一发变异（不引 poll 的 §5）：等待**不能**被产品坏掉满足

`t-mut1`（仓外副本，**改的是副本里的 `sampler.go`**，仓库那份一字未动）：把 `CheckSettle` 读错误分支里的
`rep.SampleErrors++` 摘掉（＝poll 程的 M1），整包一发，Linux：

```
M1 landed: SampleErrors++ removed from the read-error branch of CheckSettle
rc=1
RUN=71  PASS=68  FAIL=3  SKIP=0  ^panic:=0
--- FAIL: TestCheckSettleSingleTrustworthyReadReportsItsLoss (0.10s)
--- FAIL: TestCheckSettleHalfTheReadsFailedReportsItsLoss (0.10s)
--- FAIL: TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed (0.20s)
    sampler_settle_coverage_136_test.go:156: report counted 0 dropped reads but the seam took 10 reads and kept 1: 9 unaccounted
    sampler_settle_coverage_136_test.go:251: the seam lost 5 of 10 reads but the report says sample_errors=0: a half-covered window must report its losses
    sampler_settle_gate_136_test.go:223: precondition broken: an alternating tree must both keep and lose reads within the window, kept=10 lost=0 reads=20
```

⇒ 三枚红名/站点/句子与 poll 程 §5-M1 **逐枚同值**（本程独立复算，不是引它那张表）。
⇒ **承重那一味在读数里**：靶子腿红在 **0.10s**＝**等待一次就过了**（判据是接缝计数器，产品摘掉计数器不影响它过），
红的是等待**下面**那枚披露断言，且红句照旧写着 **`the seam lost 5 of 10 reads`**
⇒ **"半丢"这个场景没被等待洗掉**（洗掉了这句就印不出来），也没有一枚腿因为产品坏掉而变绿。
⇒ 这一发同时是 §3.3 那条**结构**论证的**实验**版本：本程不再只说"判据取接缝计数器所以结构上洗不掉"，
本程把它打坏了看红。

---

## 5. 契约轴与门禁读数汇总（派单点名的四问，逐问落点）

| 派单点名的轴 | 本程读数（现跑） | 落点 |
|---|---|---|
| 有没有新增 `t.Skip` | 可执行位 0 枚（两 OS 同码） | §1.5 |
| 有没有调高阈值 | 五枚文件阈值全集 md5 改前后**逐枚同值**；`want` 取值 1/2/3/4/5 全在原守卫处 | §3.2 |
| 有没有删断言 | 条件行 183→182、`t.Fatal*` 166→167，少的那枚身份＝两条开窗路径并入 `settleTreeSUT`（err 检查与路径数仍相等） | §3.2 |
| 有没有加豁免名单 | 新增行里唯一命中 `t.Skip` 字样的是**注释一句**；`git diff ee5a25e..HEAD -- internal/observe/ \| grep -E '^\+.*(nolint\|allowlist\|go:build\|t\.Skip)'` ⇒ 只那一枚注释命中 | §3.7 |
| `thresholds.go`／golden 字节 | `git diff --numstat ee5a25e..HEAD -- internal/observe/thresholds.go` 空；`grep -iE 'golden\|testdata\|\.sse$'` 于 diff 名单为空（仓内被跟踪 golden 名册 58 枚） | §3.7 |
| 生产码字节 | 两枚 commit 的 numstat 全为 `*_test.go`；`git log --oneline ee5a25e..HEAD -- internal/observe/sampler.go` **无输出** | §3.8 |
| 契约面（`D1–D47`／`C1–C32`／`PLAN.md`／`docs/specs/**`） | `git diff --name-only ee5a25e..HEAD \| grep -E 'docs/PLAN.md\|docs/specs/\|docs/contracts'` ⇒ **空** | 本格 |
| 票面有没有被谁偷偷勾 | `git log --oneline ee5a25e..HEAD -- .scratch/wisp/issues/136-…md` ⇒ **无输出**；现量勾数 **7 勾／8 未勾**（与台账 `A226` 那句一致），`AC#15` 在 `:348` 仍 `[ ]`；`git status --porcelain -- .scratch/` 空 | 本格 |
| `sh scripts/d22scan.sh` | **rc=0**，锚 `598620e`，正控（`TestBuiltBinaryGoesRedEndToEnd` 七枚子项）在同一发里 | §0.3 |
| `gofmt`／`go vet`（本包） | 空输出／rc=0，锚 `598620e` | §0.4 |
| `DEFERRED(D-xx)` 标记 | 本包 0 枚（`grep -rn "DEFERRED(D-" internal/observe/` ⇒ 0）⇒ 本批**没有新增推迟项**，不动 `SPEC-12 §5` 那张表 | 本格 |

**⛔ 一句本程不说的话**："如果这批改出了墙钟超时／裸 `go func(`，门会抓到"——**本程没把任何一枚门打成红**，
所以这句在本程手里**没有凭据**；盘上能说的只有 §0.3 那两行分母（ban #1–5 只算 production Go files、ban #8 才算 `_test.go` 与注释）。
派单那句"如果声称门会抓到，就得先把它打红"——本程因此**不声称**。

`A212④` 那四条硬约束的本程回执：**(i) 无墙钟差**（`awaitWindow` 体内只有 `NewTimeout/Expired`，§3.3）／
**(ii) 上界由现量推**（本程复核：2 s 对应"100 ms 窗 20 枚、200 ms 窗 10 枚"两 OS 同值，§4.2）／
**(iii) 不许洗掉场景**（本程自己那一发 M1 复算，§4.7）／
**(iv) 禁 Skip／禁改阈值／禁删断言／禁动 `sampler.go`**（§3.7、§3.8）。**四条全收。**

---

## 6. 逐格裁决与总裁（非实现者口径，档位逐格标）

### 6.1 `AC#15` 的判据逐格

| 判据（票面现量位置） | 被裁交付交了什么 | 本程怎么核 | 档位 | 判 |
|---|---|---|---|---|
| **① 命中率给 n、两口径分开、永不加总**（原①被 `:377` 的④改写） | poll：口径 A 20＋20 发／口径 B 3950 枚，并**自陈**"0 对 0 不构成率降"；Linux：口径 A **30＋30**、口径 B 4000＋枚、逐批负载标签；本程补：口径 B **忙窗 200 发改前 16 红／改后 0 红**＋净窗 200 发 0 红（§4.4） | 本程不复算它们那 12000＋枚，只**自己重跑同名命令**并核形状（§1）＋自己那一发负载对照 | ①〔独立复现〕（四数／名册／负载对照）；两程的大分母〔日志＋归档，抽验〕 | **成立** |
| **② 根因归到"窗口与断言之间缺的那条判据"，不许收在"负载高"** | poll §2 普查＝13 枚腿／15 处守卫"数了但不等" | 本程把它打到盘上：同一枚窗**只交付 1 枚读**时改前必红、改后重开救回（§4.4 直方），且改前树在同一副负载里 200 发红 16 | 〔独立复现〕 | **成立** |
| **③ 修法只许"有判据的等待"（轮询到条件成立＋单调超时上界）；⛔ Sleep 糊窗／Skip／放宽阈值／删断言** | 新增 `window_wait_136_test.go`，13 站点／12 腿改成"重开整窗直到判据成立"，到期 `t.Fatalf` | §3.2 阈值 md5 同值、§3.3 helper 身份、§3.6 sleep/parallel 零枚、§3.4 到期出口是红、§4.7 产品坏掉满足不了等待 | 〔独立复现〕 | **成立，且判"收紧非放水"**（判据见 §3.3 末） |
| **④ 改后同一命令复跑 ≥30 发命中 0＋名册双向差集为空（证明 0 不是用例变少／被跳／被吞）** | Linux 交 30＋30 整包单发命中 0；poll 交 20＋20（**不足 30**，它自己在 §7.5 认了） | 本程自己：`-count=2 -v` 两 OS **142/142/0/0**、名册 71 两向差集空（§1.1、§1.2）；再加忙窗 200 发 0 红；**且 §4.5 那发证明这把尺在真红上响（12 枚红、RUN 仍 71）** | 〔独立复现〕 | **成立**（"≥30"这一味由 Linux 程满足；poll 程那 20＋20 **单看不足**，本格按整格交付合看） |
| **地界：只许动 `internal/observe/**_test.go`；`CheckSettle` 语义要变就停手** | 两枚 commit 只碰 6 枚 `*_test.go` | §3.8 逐枚 numstat ＋ `git log … -- sampler.go` 无输出 | 〔独立复现〕 | **成立（未越界）** |
| **方向锁 `:379`：门已 `true`，不许把"丢"改成"不算丢"** | 阈值集合一字未动；新增全是"再开一窗" | §3.5（现量 `sampler.go:556 const settleCoverageRowGates = true`）＋ §3.2 ＋ §4.7（M1 下红句仍写 `lost 5 of 10`） | 〔独立复现〕 | **成立** |
| **票面 `:288`：AC#14 落地后靶形要现量重划** | 票面 `:350` 的 `:189`、`:373` ① 那张行号表**今天全部腐坏** | §2.3 现量重划（靶腿声明 `coverage:236`／await `:237`／红 `:243`），§2.4 数出 13 站点／12 腿／5 枚文件 | 〔独立复现〕 | **成立**（本程替它重划，票面数字请编排者按需另起 `>` 更正，**本程不碰票面**） |
| **前置一发 `:387`（`gate_:262` 那枚会 panic 的仪器，排在第 0 步之前）** | 由**别程** `021a549`（22:46）落，纯测试＋长度守卫＋`t.Fatalf` | §2.5 现量：守卫在盘上、`021a549` 是 `d8390aa` 的祖先、numstat 只 1 枚 `_test.go` | 〔独立复现〕 | **成立**（不记在本批名下，也不许拿它抵本批的账） |

### 6.2 两枚被裁交付各自一格

- **`observe-ac15-poll-r1.md`（修法本体）**：判**成立**。它自陈的四条"没测"（§7.5）本程**没有一条读成通过**；
  其中两条本程已替它还上半枚（**忙窗读数**＝§4.4 本程自己在 Windows 忙窗量了改前/改后各 200 发；
  **重开路径没被执行**＝§4.4 量到 8 次自然重开并全被救回）。它自己那两处枚数自纠（§7.6／§7.7）本程复核：
  §7.7 那张更正表**与盘上逐枚同值**（§2.4 本程另跑一遍同式）。
- **`observe-ac15-linux-r1.md`（Linux 分母程）**：已入库的 §0–§4.6 判**成立**（本程复算得同值的有：四数 142/142/0/0、
  名册 71、`!windows` 那一族在本包为空集、真饿那发的 16 窗/2.10s 形状、2s 不改的结论方向）。
  ⚠ 它 **§3.5 承诺的 §5 那节它本人没写**，本程在 §4 答了这一问，**没有回写进那份文件**（开头那句声明）；
  它 §3.5 里那枚"2s/400ms ≈ 5 枚窗"的算术**被本程现量否掉**（是 10 枚／每 await，§4.2）——
  **这条否证记在本程名下，不记成它's 缺陷**：它当时如实标着"§4 落地后在 §5 给答复"，那一格本来就是空的。

### 6.3 总裁（`AC#15`）

> **判：成立（PASS，无附条件）。**
> 七枚判据格全部本程自己量过、无一格按两程的自述收下；放水两问（断言有没有弱化／helper 是不是为过而造）
> 答案是"未弱化、判据取接缝自己的计数器且计数器在产品分支之前"。
> **本程没有勾任何东西、也没有写台账**——翻勾与 `A##` 归编排者（派单明令）。

两枚**不阻塞**的账面项（本程无写面，交给下一步落笔的人）：
1. `window_wait_136_test.go:53-55` 那句 "buys roughly 20 median-cost attempts" 要补一句
   **"20 只对 100 ms 窗成立；200 ms 窗实测 10 枚（Windows/Linux 同值，见 `observe-ac15-final-r1.md` §4.2）"**。
2. 票面 `AC#15` 格内 `:350`（"`:189` 那一族"）与 `:373` ① 那张行号表**都已腐坏**，
   按票面自己的规矩"原句不抹、另起 `>` 更正"处理；本程给的现量版在 §2.3／§2.4。

---

## 7. 本程**没有**测什么（一格都不算通过）

1. **没有复算两程的大分母**：poll 程的 3800＋3950 枚、Linux 程的 4000＋枚腿执行、`A212①` 那 4/8200＝0.049%，
   本程**一枚没重跑**（只自己另取了 200＋200＋300＋200＋40 枚）。⇒ 那些枚数在本表里全部是〔日志＋归档，抽验〕。
2. **没有复算 M2/M3/M4/M5 那四发变异**（只自己复算了 M1，§4.7）。⇒ 那四行〔仅自述，不背书〕。
3. **没有在任何 GitHub runner 上量过**：本程 Linux 侧＝本机 WSL2 内核上的 Debian 容器、12 vCPU（`ubuntu-latest` 是 4 枚）；
   承重的 CI 分母（`test-core`/ubuntu）本程只**复读了它的存在性**（引两程的 `ci.yml`／`portable-tests.sh` 行号，未自核）。
   ⇒ "2 s 在那枚 runner 上够不够"＝**〔不可判〕**。
4. **没有造出"第一枚窗成功、第二枚窗饿死"那一串**（§4.2 只证了"第二枚窗在第一条腿饿死时根本没轮到"）
   ⇒ 单腿最坏 4 s 这一支**只有代码结构支持，没有读数**。
5. **没有把任何一枚门打成红**（§5 末那句），因此不声称门能抓到测试文件里的形状。
6. **没有跑 `-race`、没跑 `t.Parallel` 相关**（本包 0 枚，§3.6）；**没跑 `go test ./...`**（派单禁令＋两枚兄弟在飞）；
   **没跑 `wisp slo`**（本程禁跑，且 `AC#14` 那格的前置另有人量）。
7. **没有量 `SampleState` 那三枚 `awaitStateReads` 腿的忙窗暴露**（只量了它的结构下界，§4.1 那一发）。
   ⇒ 想扩这一支的下一位：本程的仪器与树都在 `D:\tmp\observe-ac15-final\t-instr\`（只建不删）。
8. **没有核票面 15 枚格的整体一致性**（本程只裁 `AC#15`）；**没有核 `AC#10`／`AC#2…AC#7` 的旧账**。
9. 那 **4 枚"连前提守卫都没有"的暴露腿**（`A219②` 登记的 `TestCheckSettleVerifiesReleaseCounter` 一族）
   本程**没碰、没量、没判**——它们不属 `AC#15` 判据射程，票面另开一格才管。

---

## 8. 给人看的那一段

**修了什么**：有一批"检查睡眠时机到没到"的自动测试，自己带一个前置条件——"这一百毫秒里我至少得读到几次数"。
以前它们**不等**：窗口一关就去数，读到几算几，读少了就自己报错。这次给它们加了一件事：
**没数够就把整枚窗口重开一次再数，最多再等两秒；两秒还不行就大声报错，并把每一枚窗口实际读到几次数全印出来。**
真正测东西的那段程序（会被打包进你机器上那个 `wisp.exe` 的部分）**一个字节都没改**，界面上的东西一个都没变。

**那枚"偶发红"没了没有**：**我这次把它逼出来了**——在同一台电脑上再开十二个吃满 CPU 的空白进程，
用**旧**版本跑 200 次，红 **16 次**；用**修好**的版本跑同样的 200 次，红 **0 次**，
而且那 200 次里有 8 次确实发生了"第一枚窗口读不够、重开一枚才够"——**新加的等待真的被用上、也真的接住了**。
平时（不加负载）跑 200 次，一次都不需要重开，测试耗时和以前一样（约 7 秒）。

**还没知道的**：① 那台跑正式测试的云机器（Ubuntu runner）上两秒够不够，我在本机容器里量不出来；
② 那一族里有 4 枚测试连"读够了没"这句前置都还没写，它们不在这一格的活里，得另开票；
③ 我裁的只是这一格（`AC#15`），票上另外 8 格没勾的还轮不到我。

**用这台机器的你会察觉到什么**：**什么都察觉不到**——球、面板、日志、语音、审批，全都没有变化，
这次动的是"检查这些功能的测试自己会不会偶尔冤枉人"。**唯一可能被你看到的差别**是：
以前 CI 偶尔会红一条谁也看不出毛病的"precondition broken"，
现在**真出那种事**时报错会写清楚"两秒里重开了 10 枚窗口、每枚各读到 2 次、3 次…"——
也就是说，**红的时候你能一眼看出是机器被压垮了，还是软件真坏了**，而这正是这一格要的收口。

**你需要做什么**：**只需要一件事**——勾 `AC#15`（本程不代勾）与是否把 §6.3 那两枚账面更正派下去。
不需要你拍任何新规矩，两秒那个数我建议不动（理由见 §4.6）。

---

## 9. 纪律回执与临时件

| 项 | 回执 |
|---|---|
| 写面 | 只有本文件。**未**改 `internal/observe/**`、**未**碰票面、**未**写台账、**未**做任何 `A##` |
| commit | 每节一枚，全部显式 pathspec（`git add -- docs/evidence/s1/observe-ac15-final-r1.md` ＋ `git diff --cached --name-only` 自核）；**无** `add -A`／`.`、`--amend`、`reset`、`rebase`、`stash`、`checkout .`、`clean`、`push` |
| 一次并发回显（如实报） | §0-§1 那枚 commit 落地后我的 `git log -1` 打出来的是编排者的台账 commit（`4fab445`，晚我 1 秒）；**我那一枚是 `9450575`，`git show --name-only` 只带本文件**（已现量）。⇒ 这是共享树里渐进 commit 的正常回显，不是我的 pathspec 漏了东西 |
| 兄弟程边界 | 全程只跑 `./internal/observe/`；`internal/panel/l2_grant_boundary_test.go`（工作树 `M`）与 `d22scan-close-r1-accept-r1.md` 两枚在飞件**未参与任何读数**；`ban #8 internal/=407` 含前者，本程按 §0.3 的锚点报、不追 |
| owner 的 16 枚 `design/**` 删除＋未跟踪件 | **未还原、未暂存、未提交**（`git status` 上它们全程在场） |
| 容器 | 复用 `wisp-obs-ac15-linux-r1`（`docker start`），**未新建、未 `rm`**；12 枚自旋**起了也停了**（末次 `loadavg 5.86/10.75/8.05` 且在回落） |
| 临时件（只建不删） | `D:\tmp\observe-ac15-final\`：`logs/`（Windows 六发日志）／`panicctl/`（正控模块）／`assertdiff/`（改前后条件行对照）／`t-instr`（ await 计数版副本）／`tw-floor`／`tw-200ms`／`tw-floorstate`／`t-pre`（`ee5a25e` 副本）／`t-mut1`（M1 副本）；容器侧 `/scratch/final-r1/{t-floor,t-200ms,t-blast,t-mut1,blast.log,mut1.log,busy-400ms-leg.log,starve.sh}` |
| 凭据 | 本程**未**读取、**未**复制任何密钥值（本格不需要） |

### 9.1 交件态复核（在本文件最后一枚 commit 之后的 HEAD 上重跑，不引用 §1 的旧日志）

```
$ date "+%H:%M:%S"                       2026-09-25 11:44 +0800
$ git diff --name-only 333dfe3..HEAD -- internal/observe/      （无输出）
   ⇒ 本程所有读数跑的就是交件态那版码，其后落地的全是别家地界与本文件
$ go test ./internal/observe/ -count=2 -v        rc=0
   RUN=142  PASS=142  FAIL=0  SKIP=0  ^panic:=0    包时 6.865s
$ git status --porcelain -- docs/evidence/s1/observe-ac15-final-r1.md .scratch/     （无输出）
   ⇒ 本文件已入库、票面 136 一格未动
```

⇒ §1.1 那组数在这一发里**同值复现**（本程在 Windows 上共跑了这一族 **三发** `-count=2 -v`：11:06、11:44，
以及 §4.4 那两批忙/净窗单腿批次）。
