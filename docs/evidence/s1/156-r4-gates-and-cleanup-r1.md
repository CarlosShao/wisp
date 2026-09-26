# 票 156 — r4 程证据件（AC#4② ＋ AC#7 改后那一发 ＋ AC#6 名册 ＋ 交件报告）

- 程：**r4**（本票第四个实现程）。接手 r3：它 16:18:34 **死于模型服务连接中断**（`status=failed`，不是撞轮次上限），
  交出了 `3fbb1ce`（gofumpt 修形）与那发"改前"门禁（编排者代提为 `1137781`），**没来得及做** AC#4②／AC#7 改后／AC#6／证据件。
- 派单：`.scratch/wisp/dispatches/2026-09-26-162x-impl-156-r4-after-r3-died.md`
- 本件由 **r4 建**（`Write` 工具），r1/r2 的 `156-exited-asks-os-r1.md` 与 r3 的 `gate-r3-pre/**` **一字未改**。

---

## 第 0 格　锚点·前提·写面

### 0.1 step 0 现量（**照本程自己的量，不照派单印的号**）

```
$ date                                     Sat Sep 26 16:25:35 CST 2026
$ git rev-parse --short HEAD               1133977
$ git status --porcelain -- cmd/wisp       (空)
$ git rev-parse --abbrev-ref HEAD          dev
```

⇒ **派单给的锚点 `1133977` 与本程 step 0 现量一致**（不是"我信了它"：它是现量出来的相符）。

⚠ **但本程干活期间 HEAD 又动了**，这是编排者预告过的事（"我一直在往自己的写面提交"），本程全程现量、不拿 step 0 那个号当凭据：

| 时刻 | HEAD 现量 | 那几枚动了什么 |
|---|---|---|
| step 0 16:25:35 | `1133977` | — |
| AC#4 变异测量 | `b8ebe7d` | `b8ebe7d`＝台账 `A301`＋票 156 进度（**只动** `docs/reports/pending-and-issues.md` 与票面） |
| AC#4 改后控制发 | `88050b1` | `88050b1`＝停车点更正（**只动** `docs/reports/HANDOVER.md`） |
| AC#7 改后门禁 | `1d53526` | `1d53526`＝**本程**那枚 AC#4② commit |

**跨这些号本程的射程字节没变**（本程现量，不是推）：

```
$ for f in cmd/wisp/slo_windows.go cmd/wisp/slo_report_144_windows_test.go \
      cmd/wisp/slo_exit_os_156_windows_test.go; do
      git rev-parse 1133977:$f | cut -c1-16 ; git rev-parse b8ebe7d:$f | cut -c1-16 ; done
slo_windows.go                  d1f1b881a006a569  == d1f1b881a006a569   同一枚
slo_report_144_windows_test.go  33dc0a98baa1a776  == 33dc0a98baa1a776   同一枚
slo_exit_os_156_windows_test.go cabd5c6d926f6ae8  == cabd5c6d926f6ae8   同一枚
```

⇒ 本程所有读数钉在 `1133977`≡`b8ebe7d`≡`88050b1` 的同一批 cmd/wisp 字节上；AC#7 那一发跑在 `1d53526`（＝上面那批字节 ＋ 本程的注释改动）。

### 0.2 派单的五条硬前提，逐条现量

| 前提 | 派单写的 | 本程现量 | 判 |
|---|---|---|---|
| **P1** 票面 7 枚未勾／0 枚已勾 | `^- \[ \] ` 锚定 | `grep -c '^- \[ \] '` = **7**；`grep -c '^- \[x\] '` = **0** | **成立** |
| **P2** `gofumpt -l cmd/wisp/` 必空 | 不许再动那枚文件 | v0.12.0 (go1.27.1)，`-l cmd/wisp/` **输出空**、rc=0（改前改后各量一次，两次都空） | **成立** |
| **P3** AC#4③ 已由编排者闭 | `7ad2d98` 入库 9 枚 | 见 §2.4 现量 | **成立**（本程未重做、未降级） |
| **P4** 四支"本票不解决"不许越 | 尤其不引入 reaper goroutine | 本程**零枚**生产码改动、**零枚**新 goroutine；只改了两处注释文字 | **成立** |
| **P5** 不许把抄来的数当凭据 | — | 本程所有四数都是自己的命令打出来的（§1、§2.1、§3.1 各带原文） | **成立** |

⚠ **一处派单文字与盘的现状不符（不影响前提成立，但会骗下一个按派单行号做事的人）**：派单 §2 说那三行注释在"现 `:754`–`:756`"、r2 的更正块在"`:741`–`:756`"。本程量到更正块是 **`:742`–`:752`**、前瞻句是 **`:753`–`:758`**（各 6 行）。票面 `:23` 写的 `:737-740` 才是准的。**票面为准**，派单那两处行号已过期——派单自己 `:4` 就声明过"那文件在漂"，所以这不是本程换前提，是把过期坐标换成现量坐标。

### 0.3 写面自证（每枚 commit 之前都跑过 `git diff --cached --name-only`）

本程只碰过：`cmd/wisp/slo_report_144_windows_test.go`（派单 §2 那两处）、
新建 `docs/evidence/s1/156-r4-gates-and-cleanup-r1.md`（本件）、
新建 `.scratch/wisp/probes/156/{gate-r4-post/**,mut-156-r4/**,mut-156-r4-postedit/**,my156-r4.py,gate156-r4.sh}`。

派单 §5 点名的"别人的未提交增量"（`design/**` 25 枚、`152-…-accept-r1.md` 3 行自校、`probes/152/my152.py`、
`?? .zcodeignore`、`probes/156/__pycache__/`、`mut-156-r2/asis.log`）：**一枚未碰、未还原、未提交**，
且**一枚都没进本程任何零命中宣称**（§4 的判据只数本程自己那几枚 commit）。

---

## 第 1 格　r3 那发"改前"读数的**复算**（不重跑整包）

**档位先写明**：〔**前程现跑、读数已入库 `1137781`、本程复算**〕——**这不是 r4 跑的一发**，
本程只是把已入库的日志重新数过，用来确认自己拿对了对照件。

- 用了哪枚文件：`.scratch/wisp/probes/156/gate-r3-pre/{cmdwisp,internaltools}.log`（`meta.txt` 记它跑在 `3293192`）
- 哪条命令（与 r3 脚本同一套锚定形）：

```
$ for n in cmdwisp internaltools; do f=.scratch/wisp/probes/156/gate-r3-pre/$n.log;
    echo "RUN=$(grep -c '^=== RUN' $f) PASS=$(grep -c '^--- PASS' $f) \
          FAIL=$(grep -c '^--- FAIL' $f) SKIP=$(grep -c '^--- SKIP' $f)"; done
cmdwisp        RUN=144 PASS=84 FAIL=0 SKIP=0     verdict: ok github.com/CarlosShao/wisp/cmd/wisp 84.960s
internaltools  RUN=115 PASS=79 FAIL=0 SKIP=0     verdict: ok github.com/CarlosShao/wisp/internal/tools 13.279s
```

**名册也复算了**（不只是数）：本程从 `.log` 现抽两向名册，与 r3 已入库的 `-runs.txt`/`-verdicts.txt` 逐字节比：

```
$ diff <(本程从 cmdwisp.log 抽的 runs)  gate-r3-pre/cmdwisp-runs.txt      → 0 行差
$ diff <(本程从 cmdwisp.log 抽的 verdicts) gate-r3-pre/cmdwisp-verdicts.txt → 0 行差
   cmdwisp runs 144/144  verdicts 84/84 ； internaltools runs 115/115  verdicts 79/79
```

⇒ **r3 那发的四数与四枚名册，本程独立复算全部对得上**，可以当 AC#7 改后的对照件用。

**本程自己的一条仪器细节（比派单写得更细，不改派单的结论）**：`cmdwisp.log` 里含字符串 `FAIL` 的行
`grep -c 'FAIL'` = **4**，派单 §3 说"有 3 行印出来的 `[FAIL]`"。本程逐行现量：**3 枚 `[FAIL]` 标签属实**
（`sherpa-onnx C API`、`onnxruntime.dll version`、`data dir resolvable (dev)`），**第 4 枚是同一次 doctor 报告的汇总行
`wisp doctor: FAIL`**——同一枚用例、同一个成因、同一支形状。按 `^--- FAIL` 数＝**0**。
⇒ 差别方向不变（用锚定形），但**报数时别把 4 说成 3**。

---

## 第 2 格　AC#4② — 本票唯一一枚被批准的删除

### 2.1 先按**本程自己的锚点**现量（改前，`1133977`≡`b8ebe7d` 字节）

尺＝本程新建的 `.scratch/wisp/probes/156/my156-r4.py`：**零枚新变异字面**，
全部 `spec`/`case_off` import 自 r1 那把已入库的尺 `probes/156/my156.py`
（`probes/152/my152.py` 是别程的半件：未执行、未还原、未提交）。
选择器＝`TestSLO144|TestSLO147|TestSLO149|TestSLO152`（票 152 那一族原选择器，**不含**本票新增的 `TestSLO156`）。

```
$ python .scratch/wisp/probes/156/my156-r4.py .scratch/wisp/probes/156/mut-156-r4 \
      r4-asis-152-surface r4-152g1-case14off r4-152g1-case14live
P156R4|ANCHOR|head=b8ebe7d|branch=dev|status_cmd_wisp=[]
P156R4|CONTENTION|runner_listener_procs=1|runner_worker_procs=0
r4-asis-152-surface    rc=0  RUN=21  PASS=14  FAIL=0  SKIP=0  LOADED=YES
r4-152g1-case14off     rc=0  RUN=20  PASS=13  FAIL=0  SKIP=0  LOADED=YES   <- 写进注释的那一发
r4-152g1-case14live    rc=1  RUN=21  PASS=13  FAIL=1  SKIP=0  LOADED=YES
                                red=TestSLO152CorruptLegWithNoNamedPositionRefusesToBorrowOne
P156R4|LANDING|cell=r4-152g1-case14off|compiled_file=cmd/wisp/slo_windows.go|...|new1=yes|old1=gone
P156R4|LANDING|cell=r4-152g1-case14off|compiled_file=cmd/wisp/slo_report_144_windows_test.go|...|new1=yes|old1=gone
P156R4|AFTER|cell=...每一发|slo_windows.go=d2bc41ecfef06bd4|...=daf3455916f6bfd3|...  (落地后树哈希一字未动)
```

⇒ **票面 ② 给的 `20/13/0` 在本程的锚点上成立**（不是"预期值撞对了"：本程量到 `20/13/0`，且配对的两发说明它是哪一发）。
本程量到**三发**而不是两发，因为派单 §2 要"写明你跑了哪两发"，而第三发（`asis`）是错因的证据不是猜测：
**`21/14/0` 正是这个选择器的零变异基线**（`r4-asis-152-surface`），被贴进了一句承诺变异读数的话。

### 2.2 一处本程现量到的**加强**：那枚日志自己就说 20/13/0

`:739` 那一句是"**点名引用**某枚日志的读数"的句式（`probes/152/mut-post/g1-restored-borrowed-value-case14off.log
reports === RUN=…`）。本程不推断那枚日志印了什么，直接数它：

```
$ git ls-files --error-unmatch .scratch/wisp/probes/152/mut-post/g1-restored-borrowed-value-case14off.log   tracked: yes
$ grep -c '^=== RUN' … = 20   '^--- PASS' = 13   '^--- FAIL' = 0   '^--- SKIP' = 0
$ head -2 … | grep -o "\-run '[^']*'"      -run 'TestSLO144|TestSLO147|TestSLO149|TestSLO152'   ← 与本程同选择器
```

⇒ 把 `:739` 改成 `20/13/0` 同时改真**两件事**：那一形的读数、以及它引的那枚 tracked 日志真的印的数。
本程的 `r4-152g1-case14off` 与那枚 tracked 日志**两发同形同数**，所以这一枚不依赖任何前程的表。

### 2.3 改了两处，一枚 commit（`1d53526`）

1. `:739` 三个数 `RUN=21, PASS=14, FAIL=0` → `RUN=20, PASS=13, FAIL=0`（票面 ② 要的那一枚）。
2. `:753`–`:758` 那句前瞻改过去式，**保留它讲的因果**（零变异基线被贴进承诺变异读数的句子→读到 20/13 的人
   会去找"失踪的用例"而不是一句过期注释），并写明 **r4 已在本枚 commit 就地改真**。

**两处都做成"同行数替换"（7 加／7 删，全文 938 行→938 行）**，理由是本程现量到的：
`:875` 那句自指本文件的 `:835`（case 13 裁定头），**任何净增删都会把那条行号引用弄假**。

```
$ git show --numstat -1 1d53526                7	7	cmd/wisp/slo_report_144_windows_test.go   ← 唯一一枚文件
$ git show --name-only -1 --format=            cmd/wisp/slo_report_144_windows_test.go
$ awk 'NR==739||NR==835' <改后>                :739 已印 20/13/0 ； :835 仍是 case 13 裁定头（没漂）
$ grep -c 'reports === RUN=21, PASS=14, FAIL=0' <改后>   0    （改前=2：一处过期读数＋一处历史引用）
$ "$(go env GOPATH)/bin/gofumpt.exe" -l cmd/wisp/        空   $ go build ./cmd/wisp/   rc=0
```

**改后控制发**（证明本程这两处注释没移动任何读数）：

```
$ python my156-r4.py .scratch/wisp/probes/156/mut-156-r4-postedit r4-asis-152-surface r4-152g1-case14off
r4-asis-152-surface   RUN=21 PASS=14 FAIL=0 SKIP=0   r4-152g1-case14off  RUN=20 PASS=13 FAIL=0 SKIP=0
   两发与改前一字未变；P156R4|AFTER|slo_windows.go=d2bc41ecfef06bd4（生产码一枚字节没动）
```

### 2.4 本格**没做**的两件事（都不在本程写面，指名归属）

- **AC#4③ 已由编排者闭**（P3）：本程现量
  `git ls-files .scratch/wisp/probes/152/mut-shipped/ | wc -l` = **9**、`git log --oneline -- 该路径` 命中 `7ad2d98`。
  ⇒ 本程**未重做、未再降级**，也没"顺手"核对那 9 枚的内容（不是本程的凭据）。
- **残留缺陷（越界，本程不动，写清给验收程裁）**：`:742`–`:752` 那一整段是 **r2 的更正块**，本程被授权只动两处，
  所以它现在带着三处**因本程这枚 commit 而变得指涉过期**的文字：
  ①`:742` "Correction of that triple" —— 被纠正的那三个数**已经在 `:739` 就地改真**，这段从"更正"变成了"更正的历史"；
  ②`:749` `So "=== RUN=21, PASS=14, FAIL=0" is neither reading` —— 那串字符**今天已不在 `:739`**（它作为
    **被讨论的历史值**留在这一段里，仍然可 grep，但时态是现在时）；
  ③`:743` "At head 0f18652" 与 `:752` "again at this anchor" —— "this anchor" 是 **r2 的**锚点，
    不是读者的、也不是本程的（本程锚点 `1133977`，同一形本程量到 20/13/0 与 21/14/0，与 r2 一字不差）。
  **本程判断**：这三处都属"改 r2 已入库证据文字"，票面 ② 的射程只到 `:737-740` 那三行；
  本程不扩写面，把它们**登记为残留**交非实现者验收程裁。

### 2.5 放水两问自答（本格）

① **断言方向动没动**：没动。本格新增/删除的都是 `//` 注释文字；`git show 1d53526` 里 14 行改动**全部以 `//` 开头**，
零枚 `t.` 调用、零枚判据、零枚阈值、零枚 golden。控制发（§2.3）证明读数一字未变。
② **helper 是不是原有的那枚**：是——尺面是 r1 已入库的 `my156.py`（`spec`/`case_off`/`run`/`apply_overlay` 全 import），
本程的 `my156-r4.py` 只换了三样**打印与命名**：日志名带选择器（不写 `asis.log`）、每发打印落地哈希、跑前查争用。
本程**没有**新造变异字面，也就没有"造一枚打不红的变异"这扇门。

---

## 第 3 格　AC#7 "改后"那一发（逐包单跑 ＋ 对 `gate-r3-pre/` 的名册差集）

### 3.1 跑法与四数（原文可复算）

- 脚本 `.scratch/wisp/probes/156/gate156-r4.sh`：**与 r3 的 `gate156-r3.sh` 同形**（同包列、同 flag 集、
  同锚定模式、同名册抽法），因为两发要做差集，**脚本差一点就是仪器差一点**；落地件 `.scratch/wisp/probes/156/gate-r4-post/`。
- 跑前查了争用（`meta.txt`）：`runner_listener_procs=1`、`runner_worker_procs=0`——本机 self-hosted runner
  **有一个 Listener 常驻**，与 r3 那发同形，所以两发的 CPU 处境可比；不是"无争用"。
- `meta.txt` 现量：`anchor_head=1d53526`、`branch=dev`、`status_cmd_wisp_internal=[]`（**干净树**）、
  `go1.27.1 windows/amd64`、`dll_dir` 走 **shell 形** `/d/work/...`（非 `pwd -W` 盘符形）。

```
$ bash .scratch/wisp/probes/156/gate156-r4.sh post
pkg              rc   0xc0000135 RUN    PASS   FAIL   SKIP   verdict
cmdwisp          0    0     144    84     0      0      ok github.com/CarlosShao/wisp/cmd/wisp 114.336s
internaltools    0    0     115    79     0      0      ok github.com/CarlosShao/wisp/internal/tools 15.130s
```

| 包 | 改前（r3 跑、**本程复算**） | 改后（**本程跑**） | 四数 |
|---|---|---|---|
| `./cmd/wisp/` | 144／84／0／0 | **144／84／0／0** | **相同** |
| `./internal/tools/` | 115／79／0／0 | **115／79／0／0** | **相同** |

（`0xc0000135`＝**0** 且 `=== RUN` 非零 ⇒ 两包都**真加载真跑了**，不是"看着像没测到"。）

### 3.2 名册两向 `comm`（对**已入库的 `gate-r3-pre/`**）——差集为空，逐枚归属

```
$ pre=gate-r3-pre post=gate-r4-post; for n in cmdwisp internaltools; do for k in runs verdicts; do
      echo "pre-only=$(comm -23 <(sort $pre/$n-$k.txt) <(sort $post/$n-$k.txt) | wc -l) \
            post-only=$(comm -13 <(sort $pre/$n-$k.txt) <(sort $post/$n-$k.txt) | wc -l)"; done; done
cmdwisp      runs     pre-only=0  post-only=0      （144 ↔ 144）
cmdwisp      verdicts pre-only=0  post-only=0      （ 84 ↔  84）
internaltools runs     pre-only=0  post-only=0      （115 ↔ 115）
internaltools verdicts pre-only=0  post-only=0      （ 79 ↔  79）
```

**四个方向差集全空**，且本程**没有**只凭"数相等"下结论——名册是逐枚比过的。
预期就是空，本程这一发跨的两枚改动都**不产生改名**：
① `3fbb1ce`（r3 的 gofumpt 修形）只换 `slo_exit_os_156_windows_test.go:272` 那枚结构体字面量的花括号形；
② `1d53526`（本程 AC#4②）只换注释文字，且行数不变。
⇒ **非空才是发现**（会有人问"是不是有一枚用例被格式化吃掉了"）；空集在此**是可解释的**，不是"没测到"：
`=== RUN` 名册 144/115 与改前**同名同数**、两包 verdict=`ok`、红名册（`-reds.txt`）**改前改后都是 0 行文件**。

### 3.3 派单点名的那枚会骗尺的形状，本程量到了并写清（给验收程）

- **按 `^--- FAIL` 数＝0；按字符串 `FAIL` 数＝4**（改前改后**都是 4**，形状逐行 `diff` 相同）。
  这 4 行全在**PASS 着的** `TestAC2RealProcessRefusesOnEveryLegWithoutAppData128` 内部：
  那枚用例故意在 `APPDATA` 未设的情况下真起子进程跑 `wisp doctor`，doctor 自己的报告就印
  `[FAIL] sherpa-onnx C API`、`[FAIL] onnxruntime.dll version`、`[FAIL] data dir resolvable (dev)`
  ＋汇总行 `wisp doctor: FAIL`。⇒ **这不是真伤**（本票/本仓的判据下 0 红），
  但**也不能读成"什么都没测"**：那枚用例恰恰在断言"没有 `%AppData%` 就该拒绝启动"，它**跑过了**才 PASS。
- 射程外那枚常红 `TestC21DesignTokensFourWayAgree`（`internal/panel/tokens_fourway_test.go`）：
  本程两包日志里 `grep -c` = **0**／**0**，`grep -c 'internal/panel'` 亦 0 ⇒ **本程射程没跑偏**，不必上报。
- `-race`／`-cover`／`-overlay` 在 AC#7 这两发里**一枚都没用**（plain shipped tree）；
  变异那一侧用 `-overlay`（§2.1），**从不与 `-cover*` 同用**（`my156.py` 会 FATAL）。
- 取干净树：本程**没有**做 `git archive | tar -x`，也没有需要——两发都跑在工作树本身上，
  且 `meta.txt` 的 `status_cmd_wisp_internal=[]` 与每发前后的树哈希共同钉住"跑的就是这批字节"。
