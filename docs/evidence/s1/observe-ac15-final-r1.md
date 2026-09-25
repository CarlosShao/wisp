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
