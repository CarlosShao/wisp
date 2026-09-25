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
