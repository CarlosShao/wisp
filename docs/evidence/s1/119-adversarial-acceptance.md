# 票 119 — 独立对抗验收（`acceptor-ticket119`）

被验：`.scratch/wisp/issues/119-posix-link-leg-refuses-legitimate-symlinked-data-roots.md`
票 113 那条 POSIX 链接腿的**副作用面**（合法的软链数据根被拒）。票面六格 AC 由实现方**全部自勾 `[x]`**，
语义裁成 **②**（调用方把 OS 给的答案解析干净，底线一字不改）。

- 验收人：`acceptor-ticket119`（**只读**；唯一写件＝本文件；一行生产码/测试码都没改，见 §末自证）
- 锚定 commit：`8c8aad3`（票面入库）；实现件：`189cb1e`；**我的控制组 = `ce666ea` = `189cb1e^`**
  （`git diff --numstat ce666ea..8c8aad3` 全树只有 5 枚文件：`cmd/wisp/doctor.go`、`internal/proc/envfork.go`、
  `internal/winsec/winsec_other.go`、`internal/winsec/dataroot_symlink_119_other_test.go` + 票面本身
  ⇒ 除这枚实现件外**同基座一字不差**，比票面自用的 `823d457`（隔着票 115/117 三笔）更干净）
- 纯净快照（全部在仓外；仓库内**没有**建 worktree/checkout，A38④）：
  - `D:\tmp\ac119\snap8c8aad3` = `git archive 8c8aad3 | tar -x`（被验侧，429 枚 `.go`）
  - `D:\tmp\ac119\snap-ce666ea` = `git archive ce666ea | tar -x`（控制组；`SealableRoot` 出现 **0** 次，已核）
  - 容器：`golang:1.27`（go1.27.1 linux/amd64，uid=0），挂载 `-v D:\tmp\ac119\snap…:/src`
  - **每一把都在容器内 `ls -l /src/go.mod` + `md5sum` 三枚被验文件**（假绿坑，见 §〇）

## 总判（先给结论，逐格证据在下面）

**退回。** ② 这条纪律在生产里**没有被覆盖完**：同一枚容器、同一枚形状（`$HOME` 经软链）、
**同一枚修后二进制**——`wisp doctor` 的落点行已从 `/varlink/…` 变成 `/realpriv/…`、`wisp providers` 也不再报
"凭据存储不可用"（改报"配置未就绪"），**但 `wisp secret` 仍 rc=1**，
错误原文照旧是 `winsec: refusing to seal …/varlink/home/.config/wisp-dev/secrets`。
本票立项要防的结局（**合法数据根被误伤 ⇒ 产品跑不起来**）在票面自己点名的第二条生产路（`~/.config` 被软链）
上**今天仍然被真实造出来**，且它正是本票那句关键读数所用的那枚命令（`wisp secret list`）。
按票 103/107 的口径：**不许盖"附条件"章。**
其余五格：**AC#1 通过附条件、AC#3 通过、AC#4 通过、AC#5 通过、AC#6 通过**。
AC#3（反半边）**一枚都没被放绿**：票 113 那 16 枚两形 outcome 逐字不差，`MUT-119-D` 拆腿照样 11 枚红；
更要紧的是我亲手把**放行侧降成"存在性/信任调用方"**试了一发（`MUT-119-E`）——它清掉了那 83 条 harness 红里的 62 条，
**同时让本票自己那两枚反半边用例变红** ⇒ **"放宽放行侧换绿"这条路今天是有牙齿挡着的**，票 107b 那个退回理由不成立。
**81 条 harness 红不在 AC#1 的清除射程内**：AC#1 的原文自己写的是"断**今天的真实结局**（红就是红）"，
而且 79 枚去重被拒路径**全部**是测试自造的 `t.TempDir()`（0 枚是 OS 给的数据根）——但它确实还没人清，见 `R-119-8`。

| 格 | 结论 | 标签 |
|---|---|---|
| AC#1 误伤面做成可重跑用例 | **通过附条件** | 〔独立复现〕 |
| AC#2 裁 ①/②/③ + "winsec 一字未动" | **退回** | 〔独立复现〕 |
| AC#3 反半边仍拒 | **通过** | 〔独立复现〕 |
| AC#4 `winsec.go` 边界话不自行改 | **通过** | 〔独立复现〕 |
| AC#5 变异（退回旧语义 + 两发半修） | **通过** | 〔独立复现〕 |
| AC#6 门禁四数 / 格式 / vet / d22scan | **通过** | 〔独立复现〕 |

---

## 〇、开工前的仪器自检：我自己也假红/假绿过一发

**坑（我自己的，方向是假绿）**：第一版 `asym.sh` 写的是 `rm -rf /realpriv /varlink` 之后
`mkdir -p /realpriv /varlink/home` —— `/varlink` 那一刻还不存在，`mkdir -p` 于是把它**建成了真目录**，
链接根本不存在，`wisp secret list` 在 dev 形状下**读到 rc=0**。`ls -ld /varlink` 打出来是
`drwxr-xr-x`（不是 `lrwxrwxrwx … /varlink -> /realpriv`）才暴露。
⇒ 后面所有形状脚本都加了硬断言：`case "$(ls -ld /varlink)" in l*) ;; *) exit 97;; esac`
与 `readlink -f /varlink/home = /realpriv/home`。断言不过就当停。**这条与容器挂载假绿是同一族**：
形状必须自证存在，否则"绿"是仪器的。

**派单点名的三条坑我都按口径走了**：
1. 挂载假绿 ⇒ 每把都 `ls -l /src/go.mod`（`-rwxrwxrwx 1 root root 883`）+ `md5sum` 三枚被验文件，见每次日志头；
2. `GOOS=linux go vet` 只编译不执行 ⇒ 本票产品结论**全部来自 `CGO_ENABLED=1` 真建真跑的二进制**（见 §六），
   vet 只当"能不能编译"用；我自己还踩了一发 `rc=127 error while loading shared libraries: libsherpa-onnx-c-api.so`
   （漏了 `LD_LIBRARY_PATH`）⇒ 当场判为**仪器缺件**、没有当成产品红；
3. CR 读数逐文件量 ⇒ 见 §六末段（不全局断言）。

另两条**仪器口径**先说清，免得下面被读错：
- 票面那组 `RUN/PASS/FAIL/SKIP` 是**含子测试**的口径（`^\s*--- FAIL`）。我只按顶层 `^--- FAIL` 数第一遍会读到 **68**，
  换成含子测试才是 **81** ⇒ 下面每一次四数都标了口径，两口径我都有原始行可查。
- 控制组我选 `ce666ea`（= `189cb1e^`）而不是票面自用的 `823d457`：同基座、只差本票四文件，
  这样 `15/0/0/7/26/33 = 81` 与锚点 `15/2/0/7/26/33 = 83` 的差**只能**来自这枚实现件。

---

## 一、AC#1 —— 误伤面做成可重跑的用例（容器真跑，两形，自证挂载）：**通过附条件**〔独立复现〕

**AC#1 的字面要求我先钉住**：它要的是"断**今天的真实结局**（红就是红）"，即**把这副面做成可重跑的用例**，
不是"把这副面清绿"。这个区别决定了下面三件事怎么判。

### 我亲手量的（不是读它的表）

容器 `golang:1.27`（go1.27.1 linux/amd64，uid=0），形状 `ln -s /realpriv /varlink` + `TMPDIR=/varlink/w119tmp`；
每次运行前 `ls -l /src/go.mod`（`-rwxrwxrwx 1 root root 883`）+ 三枚被验文件 md5，形状还有硬断言（`ls -ld /varlink`
必须是 `lrwxrwxrwx … /varlink -> /realpriv`，否则 `exit 97`）。

| 读数 | 我的复算 | 票面自述 |
|---|---|---|
| 控制组 `ce666ea` 六包 link 形 | `RUN=263 PASS=181 FAIL=81 SKIP=1` rc=1 | 同 |
| 锚点 `8c8aad3` 六包 link 形 | `RUN=268 PASS=184 FAIL=83 SKIP=1` rc=1 | 同 |
| 控制组 六包 plain 形 | `RUN=294 PASS=293 FAIL=0 SKIP=1` rc=0 | 同 |
| 锚点 六包 plain 形 | `RUN=299 PASS=298 FAIL=0 SKIP=1` rc=0 | 同 |
| 红名逐名比对 | 只多一枚 `TestLayoutForTestEnv`（父 + `falls_back_to_TEMP_per_pid` = 2 行） | "只多出这一枚" |
| 那枚红的原文 | `envfork_test.go:108: test DataDir = "/realpriv/w119tmp/wisp-test-310", want "/varlink/w119tmp/wisp-test-310"` | 同 |
| 本票 5 枚用例 | 两形下 **10/10 全 PASS，0 SKIP** | 同 |

⇒ **票面那组数我一名一号全部复现，没有一条对不上。**（含它自己更正过的"四数按含子测试口径"。）

**"那些红全部是 harness 形状"这句，我独立量了，成立**：控制组 link 形里被拒的**去重路径 79 枚，79/79 全在
`/varlink/w119tmp/Test*` 之下**（测试自己 `t.TempDir()` 长出来的根），**0 枚**是 OS 给的数据根。
逐包分布（link 形，`FAIL=` 含子测试）：控制组 winsec=15 / proc=0 / secret=0 / config=7 / agent=26 / memory=33 = **81**；
锚点 winsec=15 / **proc=2** / secret=0 / config=7 / agent=26 / memory=33 = **83**。

### 为什么不是"通过"而是"通过附条件"

**条件 1（这条是硬的，直接喂给 AC#2 的退回）**：`~/.config` 那一枚用例钉的是 `proc.Summarize` →
`proc.DefaultLayout` 这条路；而生产里还有**第二枚读同一份 OS 答案的人** —— `cmd/wisp/secret.go:118
resolveSecretLayout` 直接把自己读的 `os.UserConfigDir()` 交给 `proc.LayoutFor`，**不走** `DefaultLayout`。
这条路上"合法数据根被误伤 ⇒ 跑不起来"的结局**今天仍在被真实造出来**，而**用例对它全盲**：

```
/r/wisp-after  dev + HOME=/varlink/home      secret list  rc=1   refusing to seal /varlink/home/.config/wisp-dev/secrets
/r/wisp-sec    同一形状（只多一行 SealableRoot） secret list  rc=0   dir=/realpriv/home/.config/wisp-dev/secrets
锚点 vs SEC 在 ./internal/winsec ./internal/proc ./internal/secret 的 60 条 outcome：diff rc=0（一字不差）
```

⇒ **产品行为差一个 rc，判据仪器一条都不差。** `grep os.Symlink cmd/wisp/*_test.go` 命中 **0** 次，
`grep resolveDataDir cmd/wisp/*_test.go` 命中 **0** 次：这条形状在本仓**没有任何一枚用例**。
AC#1 说"先把误伤面做成可重跑的用例"，误伤面枚举漏了一整条生产路。

**条件 2（这条最硬：② 自己的分界线被抹平，仪器一声不响）**：5 枚里有 2 枚的 **fixture 用被测函数算**
（`root := filepath.Join(proc.SealableRoot(base), "data")`、`injected := filepath.Join(proc.SealableRoot(base), "harness", "picked")`），
正撞本仓自己写在 `internal/proc/envfork_test.go:117` 的那句规矩（"an expectation computed by the function under test cannot fail"）。
我把这一撞量成了两发可重跑的变异：

```
MUT-119-R  落地 internal/proc/envfork.go:104  return SealableRoot(dir) /* 把"声明的树原样返回"这条纪律直接抹掉 */
           go vet rc=0 → ./internal/winsec ./internal/proc plain 形：RUN=56 PASS=56 FAIL=0  ← 全绿
           连 TestAC2POSIXInjectedTestDataDirStandsAsDeclared119（唯一以"声明树不动"命名的那枚）也 PASS
           原因可指名：SealableRoot 幂等，而该用例的 declared 值本身就是 SealableRoot(base)+…
                       ⇒ "原样返回"与"也解析"两种实现下断言同样成立 ⇒ 这一腿恒真

MUT-119-A@HEAD 落地 envfork.go:162 real, err := cur, error(nil)（把解析变成恒等），go vet rc=0
           link 形 FAIL 17→23；新红 4 枚里 TestAC2…StandsAsDeclared119 与
           TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119 的红因是"fixture 不再干净"，
           不是它们各自要钉的语义 ⇒ 读红名会读错对象
```

⇒ **AC#2 要答的那句"该不该由调用方解析成实路径"，票面答对了（该，且只解析 OS 的答案），但钉住它的用例不可伪。**
另两枚（`...UnresolvedSymlinkedRootStillRefused119`、`TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`）
我用 `MUT-119-D`/`MUT-119-E` 各量过，**它们有牙齿**（照样红）。

**条件 3**：AC#1 末句"并自证挂载非空"是**跑者的性质**，不是 Go 用例的性质；实现方那一份挂载自证只活在票面转述里
（日志未入库）⇒ 那半句属〔仅自述，不背书〕。我自己的每把都带 `ls -l /src/go.mod` + md5 + 链接硬断言，
结论建立在我的复现上，不建立在它的转述上。

## 二、AC#2 —— 裁 ①/②/③ 的理由是否切在同一件东西上 + "winsec 一字未动"：**退回**〔独立复现〕

这一格分三层：**(a) 那一刀切的是不是同一件东西**（派单的首要攻击点）、**(b) "winsec 一字未动"这句**、
**(c) ② 这条纪律在生产里是否做完了**。**(a) 我判实现方对，(b) 半对半错（它自己如实登记了），(c) 没做完 ⇒ 退回。**

### (a) ①/② 确实切的不是同一件东西 —— 我用变异把它量出来了，不是靠读它的理由

它的三条硬理由（拿不到"声明的树"、会让 POSIX 比 Windows 更松、这条腿零 CI 覆盖）我不复述，
我只补一句**可重跑的判据**：**"清掉那 81 条"与"让合法产品形状跑起来"不是同一个对象**，
因为清掉那 81 条的捷径恰恰是**放宽放行侧**——票 107b 被判退回的那个动作。我亲手做了那一发（`MUT-119-E`，见 §五）：

```
把 platformVerifyPlacement 的拒绝条件降成"这前缀 os.Stat 得到就放行"（存在性 / 信任调用方）
→ 六包 link 形：RUN=299 PASS=277 FAIL=21     ← 那 83 条 harness 红被清掉了 62 条
→ 同形状 ./internal/winsec ./internal/proc：FAIL=11，其中
     TestAC1POSIXUnresolvedSymlinkedRootStillRefused119        ← 本票自己的反半边，红
     TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119    ← 本票自己的 AC#3 那枚，红
     + 票 113 的 9 枚 AC#1 拒绝腿（含 4 枚子测试）全红
```

⇒ **三条结论**：① "① 能清掉这 81 条，② 清不掉"是真的（我用放宽底线这一发复现了"清得掉"）；
② 但清掉的代价是**放行侧被降成存在性判断**，那正是票 107b 的退回理由；
③ 而**本票交回来的用例对这一发是有牙齿的**（两枚自家具名红）。
所以：**没有拿放宽放行侧换绿**。派单第一问的第二个分支（"是不是又放宽了放行侧"）——**否，实测排除**。

再看"81 条在不在票 AC 的射程内"。AC#1 的原文是"断**今天的真实结局**（红就是红）"——**票面自己把这一格定义成
'把结局钉住'，不是'把它清绿'**。加上 79/79 的失败对象都是 harness 自造的 `t.TempDir()`（§一），
我判：**那 81 条不在本票 AC 的清除射程内**，实现方把它们交回票 118/111 是对的走法，不算藏。
（但"不在射程内"不等于"没人管"——见 §八 `R-119-8`，HEAD 实测它们仍是 83 条。）

### (b) "winsec 一字未动"：判定**未动**成立，"一字"这个字面**不成立**（它自己登记过）

换三把我自己的尺（不用它的 `grep -v '^[+-]//'` 那一把）：

| 尺 | 读数 |
|---|---|
| 注释与空行全剥后两版逐行 diff | before **59** 行 / after **59** 行，**diff hunks = 0** |
| `git diff --numstat ce666ea..8c8aad3 -- internal/winsec/` | `winsec_other.go 21 6`、`dataroot_symlink_119_other_test.go 271 0` |
| 它的尺重跑（去注释后的 `+/-` 行数） | **0**（我复算一致） |

⇒ `if ancestorIsLink(prefix)` 那三行判定**确实一字未动**，这我用"剥尽注释再比"这把更狠的尺也验了。
**但 AC#2 的字面是"winsec 是否一字未动"**，真实读数是 `internal/winsec/winsec_other.go` 里**动了 21 增 / 6 删注释**，
`winsec.go` / `resolve.go` / `winsec_windows.go` 是 0（见 §四）。实现方在票面明写"本票动了那段注释，没动那三行判定"
⇒ **不是隐瞒，是措辞**。裁决表里这句只能记成"**判定未动，注释动了 27 行**"，不能记"一字未动"。
而那 27 行注释里有一句**被证伪的全称**（下一小节）。

### (c) ② 只做了一半：同一份 OS 答案有第二个生产读者没被覆盖 ⇒ **这就是退回**

```
internal/proc/envfork.go:106  TestDataDir      -> SealableRoot(os.TempDir())      [已覆盖]
internal/proc/envfork.go:226  DefaultLayout    -> SealableRoot(os.UserConfigDir)  [已覆盖]
cmd/wisp/doctor.go:249        resolveDataDir   -> SealableRoot(os.UserConfigDir)  [已覆盖]
cmd/wisp/secret.go:119        resolveSecretLayout-> os.UserConfigDir() 原样交给 LayoutFor  [未覆盖]
```

真二进制读数（同一容器、同一形状、同一枚锚点构建 `md5 23392aa0…`）：

| 形状 | 修前 `wisp secret list` | **修后（锚点）** |
|---|---|---|
| `WISP_ENV=test` + `TMPDIR=/varlink/w119tmp` | rc=1 `refusing to seal /varlink/w119tmp/wisp-test-13/secrets` | **rc=0** `dir=/realpriv/w119tmp/wisp-test-21/secrets` |
| `WISP_ENV=dev` + `HOME=/varlink/home` | rc=1 | **rc=1，错误原文一字未变** |
| `WISP_ENV=dev` + `/realhome/.config` 本身是链接 | rc=1 `…the link at /realhome/.config…` | **rc=1，一字未变** |
| `WISP_ENV=dev` + `XDG_CONFIG_HOME=/varlink/xdgcfg` | rc=1 | **rc=1，一字未变** |
| 同形状 `wisp doctor` / `wisp providers`（走 resolveDataDir） | doctor 打印 `/varlink/…`；providers rc=2 报"凭据存储不可用" | doctor 打印 **`/realpriv/…`**；providers 仍 rc=2 但红因变成"配置未就绪"（**密封拒绝这一步已经过去了**）|

⇒ **同一枚修后二进制上，落点在 `doctor`/`providers`/`run` 这条路上已经被解析干净，`wisp secret` 这条路仍被底线拒**，
且拒的是**本票关键读数所用的那枚命令**。这不是"登记不藏的残余"，是 AC 声称要防的结局在点名 routes 上仍成立。

**能力类那一问先答掉（派单第 6 条）**：`SealableRoot` **不是**"做出来了没人调"那一族——
`go list -deps ./cmd/wisp` 里 `github.com/CarlosShao/wisp/internal/proc` 命中 1 次，
生产调用点 **3 处**（`cmd/wisp/doctor.go:249`、`internal/proc/envfork.go:106`、`internal/proc/envfork.go:226`），
测试调用点 2 处。而且它的效果在**真二进制**上看得见（`wisp doctor` 那行从 `/varlink/…` 变 `/realpriv/…`），
不是只在测试里被调。⇒ 这一问**通过**；缺陷不在"有没有生产调用者"，在"**生产调用者少了一个**"。

修法我也替它验证过了（`MUT-119-SEC`）：只加 `sessionLayout(env, proc.SealableRoot(root), exeDir)` 一行 ⇒ rc=0。

**"那条文件不在我的地界（票 117 的核心区）"这句我核了，站不住**：票面 `Packages` 枚的是
`winsec_other.go + cmd/wisp/doctor.go + internal/proc/envfork.go + internal/memory/open.go`，
`禁改` 列里**没有** `cmd/wisp/secret.go`——它既没被授权也没被禁；而
`git log --oneline -- cmd/wisp/secret.go` 最近三笔全是**票 63**，工作树里 `cmd/wisp/` 现在**干净**，
票 117 也已入库结案（`9740dd2`）⇒ 没有"在飞所以不能碰"这个约束。

### (d) 顺带量到：② 的"只解析 OS 答案，不动声明树"这条原则句与实现不自洽

`WISP_TEST_DATA_DIR` 原样返回（有用例钉：`TestAC2POSIXInjectedTestDataDirStandsAsDeclared119`，我实测 PASS）。
但 `os.UserConfigDir()` 的答案**本身就可能是声明**（`XDG_CONFIG_HOME` / `HOME` 是注入方写的），而它**会被改写**：

```
声明值 = /varlink/xdgcfg/wisp-dev
修后 wisp doctor 打印 = /realpriv/xdgcfg/wisp-dev   ← 声明过的树被换成解析后的拼写
修前 wisp doctor 打印 = /varlink/xdgcfg/wisp-dev
```

同一枚注入变量 `WISP_TEST_DATA_DIR=/varlink/home/declared` 则**原样**交给底线并**被拒**（rc≠0，两处都不落盘）。
⇒ 原则句要成立，判据得改写成"**是否由本进程之外注入**"或者两处都解析；现在这是一句**只对一半的自律话**。
危害：`Store.Dir()` 与调用方声明的拼写不再相同（票面自己列为"不值"的代价，但那半句在 `memory/open.go`，
这半句在生产默认路上没人登记）。登记为 `R-119-3`。

### (e) 它做对的部分（记进裁决表，别只记退回）

`doctor.go` 的 test 分支不再复制一份 `filepath.Join(os.TempDir(), …)` 而是改走 `proc.TestDataDir`——
**重复实现被消掉一处**；票面点名的两处（`doctor.go:235`、`envfork.go:98`）我用控制组原文核过，
行号与形状**属实**（两处都是 `return filepath.Join(os.TempDir(), fmt.Sprintf("wisp-test-%d", os.Getpid()))`）。
`SealableRoot` 的失败方向是"把原拼写交回去"而不是吞成 `""`，这一点与它的注释一致（判定尺见 §五 MUT 表）。

## 三、AC#3 —— 反半边：穿过链接把密封带到另一棵树，今天还拒不拒：**通过**〔独立复现〕

派单口径："这条拒不了就是退回，不是附条件。" 我按三把尺量，**结论是仍拒**。

**(1) 票 113 那 16 枚的 outcome 集合，控制组 vs 锚点，逐名逐状态比对：**

```
link  形：control=16 行  anchor(去掉本票 5 枚)=16 行   diff rc=0   ← 一字不差
plain 形：control=16 行  anchor(去掉本票 5 枚)=16 行   diff rc=0
link 形下仍红的 5 枚（修前修后同名同因，都是 harness 未解析根）：
  TestAC2POSIXDoesNotFoldABackslashIntoASeparator
  TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator
  TestAC3POSIXSealFileStillNarrowsAPlainFileInsideTheNamedTree
  TestAC3POSIXSealStillWorksNextToAndThroughRealDirectoriesAndLinks
  TestAC4POSIXFloorAnswersInsideTheNamedTree
plain 形：16/16 绿（修前修后都是）
SKIP：两形下 113/119 相关 0 枚；六包里唯一 SKIP = 既有的 TestSubprocessCrashWriter
```

⇒ **一枚都没被"变绿方式"处理**：没 Skip、没删、没改断言、没把"拒"换成"静默不封"。

**(2) 拆腿反证 `MUT-119-D`（`if false && ancestorIsLink(prefix)`）—— 我自己重做，锚点快照：**

| 形状 | 读数 | 红的名字 |
|---|---|---|
| plain | `RUN=21 PASS=10 FAIL=11 SKIP=0` | 票 113 的 AC#1 那 9 行（含 4 枚子测试）+ **本票 2 枚**：`TestAC1POSIXUnresolvedSymlinkedRootStillRefused119`、`TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119` |
| link | `RUN=21 PASS=9 FAIL=12 SKIP=0` | 同上 + `TestAC2POSIXDoesNotFoldABackslashIntoASeparator` |

⇒ **FAIL=11 复现**，且腿一拆本票自己的两枚反半边用例同红 ⇒ 它们不是装饰腿。（它自述"113 的 6 枚 AC#1 全红"，
我数到 9 行——那 6 枚里有 1 枚带 4 个子测试，口径差在子测试行，红名集合与它列的一致。）

**(3) 真二进制层，"穿过链接跑到别人的树"今天拒不拒：**

```
在②解析过的数据根里种一个链接 <root>/secrets2 -> /realpriv/elsewhere，然后 winsec.PrivateDirAll(<root>/secrets2)：
  -> refusing to seal /realpriv/poison/wisp-test-24/secrets2: winsec: path is not provably resolved,
     refusing to seal: … reaches it through the link at /realpriv/poison/wisp-test-24/secrets2,
     which is not the tree this call names
被指到的 /realpriv/elsewhere 事后仍是 drwxr-xr-x（没被 chmod 收窄）⇒ 别人的树一个字节没动
```

**同时量到 ② 换出来的那半面（如实登记，不是我推的）**：修后 `TMPDIR=/varlink/poison` 下真二进制 rc=0，
并在**链接指向的那棵树里**建出并收窄了 `wisp-test-24`(0700) / `secrets`(0700)；修前同形状 rc=1 且**什么都不建**。
⇒ 底线仍然**不穿过**链接密封（上面那条 refusal 就是它），但**"哪棵树是数据根"这个选择权，从"响亮拒绝"变成了
"跟着链接走一次"**。这是 ② 的固有代价，票面只写了"harness 那半边不解决"，**没写这半面** ⇒ 登记 `R-119-4`。
危害评估我按住不裁：要走到这一步得先控制目标进程的 `TMPDIR`/`HOME`/`XDG_CONFIG_HOME`，
而捏住这三个变量本身就已经决定了写入位置——是否算独立缺陷请编排者裁。

## 四、AC#4 —— `winsec.go` 那段边界话（票 113b 交的字）没自行改：**通过**〔独立复现〕

派单要我换一把尺复算，我用了三把（都不是它那把 `grep -v 注释`）：

```
(1) md5 逐版比对（最硬：改了注释也会变）
  git show ce666ea:internal/winsec/winsec.go  | md5sum = a6144c880de80e43bb1393f3624e7221
  git show 8c8aad3:internal/winsec/winsec.go  | md5sum = a6144c880de80e43bb1393f3624e7221
  git show HEAD:internal/winsec/winsec.go     | md5sum = a6144c880de80e43bb1393f3624e7221
  resolve.go: control 与 anchor 同为 7eb8a754eb3db1c8cf42a9dfeaa33074
(2) 1499efe..HEAD 的 log 尺（"票 113b 交的字之后有没有人碰"）
  git log --oneline 1499efe..HEAD -- winsec.go resolve.go winsec_windows.go placement_windows.go
    → 只有 391878a 一笔；git show --numstat 391878a 显示它动的是 **winsec_windows.go**（票 115，
      票面禁改列里明写"winsec_windows.go（票 115 在飞）"）⇒ **winsec.go 与 resolve.go 自 1499efe 起 0 笔**
(3) numstat 尺
  git diff --numstat 1499efe..HEAD -- internal/winsec/winsec.go internal/winsec/resolve.go  → 空（0 行）
  git diff --numstat ce666ea..8c8aad3 -- internal/winsec/  → 只有 dataroot_symlink_119_other_test.go 271/0
                                                             与 winsec_other.go 21/6
```

⇒ **`winsec.go` 包文档（113b 交的字）一字未动，`resolve.go` 一字未动，`winsec_windows.go` 一字未动。**
AC#4 的字面（"若需要改 `winsec.go` 那段边界话就先登记，不要自行改"）成立：**它没改，也确实用不着改**——
判定语义未变 ⇒ 包文档那句"none of the placement checks asks whose tree it is"仍然真。

**但要把一件事实同时记在这里**：它**动了** `internal/winsec/winsec_other.go` 里"两条成本"那段注释 27 行，
而那段新注释里含有一句**全称断言被我的 `wisp-sec` 读数证伪**：

> "the layer that asks the OS resolves what the OS answered (proc.SealableRoot, used by internal/proc's
> data-root resolution and cmd/wisp's resolveDataDir), **so a root handed to this floor names a real tree**"

真实读数：`cmd/wisp/secret.go` 这条路**既不经过 `DefaultLayout` 也不经过 `resolveDataDir`**，
交给底线的仍是一个经软链的拼写（三形 rc=1，§二(c)）。⇒ 这句"so a root …"在**生产默认路上过头**，
它读的正好是 AC#4 不许自行改的那本边界账的邻居文件。**下一格（AC#2）的退回里含这一条**，
但若编排者要把"改注释"单独记账，`R-119-2` 是给它的。

## 五、AC#5 —— 变异：**通过**（它四发我全部复现；我自己另加三发）〔独立复现〕

规矩我照走：**每发先证落地**（`grep` 出那几字节原文 → `go vet` rc=0 → 才读红名），
且全部在 `git archive 8c8aad3 | tar -x` 的仓外纯净快照里做（`D:\tmp\ac119\mut\{A,B,C,D,E,R,SEC}`，
`MUT-119-A@HEAD` 打在 `mut/HEADA` = `git archive HEAD` 的副本上；`mut/HEAD` 是同 sha 的未变异对照组），
**仓库内没建 worktree/checkout**（A38④）。每把另带 `ls -l /mut/<x>/go.mod` 证明文件真解出来。

| # | 锚点（落地后 grep 原文） | vet | 读数（我的，含形状） |
|---|---|---|---|
| **119-A** 新语义退回旧行为 | `internal/proc/envfork.go:162: real, err := cur, error(nil) /* MUTATION-119A */` | rc=0 | plain 形 `./internal/winsec ./internal/proc`：`RUN=56 PASS=54 FAIL=2 SKIP=0`，红的正是 `TestAC1POSIXSymlinkedTempDirRouteBecomesSealable119` + `…ConfigDirRoute…119` |
| **119-B** 半修（只解 TMPDIR） | `envfork.go:226: l, err := LayoutFor(env, dir) /* MUTATION-119B */` | rc=0 | `FAIL=1` = config 那一枚 |
| **119-C** 半修（只解 config） | `envfork.go:106: return filepath.Join(os.TempDir(), fmt.Sprintf("wisp-test-%d", os.Getpid())) /* MUTATION-119C */` | rc=0 | `FAIL=1` = TMPDIR 那一枚 |
| **119-D** 拆掉票 113 那条腿 | `winsec_other.go:130: if false && ancestorIsLink(prefix) { /* MUTATION-119D */` | rc=0 | plain `FAIL=11` / link `FAIL=12`（见 §三） |
| **119-E**（我的）**放行侧降成存在性** | `if _, err := os.Stat(prefix); err == nil { continue }` | rc=0 | 六包 link 形 `RUN=299 PASS=277 FAIL=21`（harness 83→21）**同时** `winsec+proc` plain 形 `FAIL=11`，含本票 2 枚反半边 ⇒ "清 81 条"与"守住底线"是互斥动作 |
| **119-R**（我的）**抹掉 ② 自己的分界线** | `envfork.go:104: return SealableRoot(dir) /* MUTATION-119R */`（声明的注入根也解析） | rc=0 | `RUN=56 PASS=56 FAIL=0` **全绿** ⇒ 包括那枚专门钉"声明树原样"的用例在内，**没人守这条线**（§一 条件 2） |
| **119-SEC**（我的）**补上缺的那一行** | `cmd/wisp/secret.go:129: return sessionLayout(env, proc.SealableRoot(root), exeDir)` | rc=0 | 真二进制 dev 三形 rc=1→**rc=0**；而 `winsec+proc+secret` 的 **60 条 outcome diff rc=0** ⇒ 本票的变异仪器看不见这一格 |
| **119-A@HEAD**（我的） | 同 119-A，但打在 `git archive HEAD`（`9740dd2`）快照上 | rc=0 | link 形 `FAIL=17→23`；新红里含 `envfork_test.go:137`（118 改写后的那枚），见 §七① |
| **探针**（我的，非变异） | `probe119/main.go` 直接问 `winsec.PathResolverInstalled()` | build rc=0 | plain 形 → `risk.c26Pipeline`；**TMPDIR 经软链 → `<nil>`**；**HOME 经软链 + TMPDIR plain → `risk.c26Pipeline`** ⇒ ③ 的开关是"OS 的 temp 树经不经软链"，与数据根无关 |

**判**：AC#5 的字面（"新语义退回旧行为 ⇒ AC#1 用例必须红；再试一发半修 ⇒ 也要红；每发先证落地"）**逐条成立**，
读数与票面自述一致（它 119-A 报 `RUN=21 PASS=19 FAIL=2`，包集不同所以我 `RUN` 更大，**FAIL 数与红名一致**）。
**唯一的缺口不在它有没有做变异，而在它能被谁看见**：`MUT-119-SEC` 证明这套仪器对第二条生产 config 路**零敏感**——
这条记在 AC#1/AC#2 的账上，不记在 AC#5 的账上（AC#5 只承诺"做了这些发"）。

## 六、AC#6 —— 门禁四数 / gofmt / gofumpt / vet / d22scan：**通过**〔独立复现〕

全部在容器（`golang:1.27`，uid=0，plain 形状 `TMPDIR=/tmp/p119`）里对**锚点纯净快照**跑，每把带挂载自证。

| 项 | 我的读数 | 票面自述 | 判 |
|---|---|---|---|
| 六包 `-count=2 -v` | `RUN=598 PASS=596 FAIL=0 SKIP=2` rc=0；**不同 `=== RUN` 名 = 299** ⇒ 299×2=598 对上 | 同 | 〔独立复现〕 |
| 同六包非 `-v` | `PASS` 命中 0 行、`SKIP` 命中 0 行、`ok` 6 行 rc=0 | 同 | 〔独立复现〕 |
| `gofmt -l`（四文件） | 空 | 空 | 〔独立复现〕 |
| `gofumpt v0.7.0 -l`（四文件） | `go install mvdan.cc/gofumpt@v0.7.0` rc=0，`-version` = **v0.7.0 (go1.27.1)**，`-l` 输出空 | "本机 v0.7.0 存在" | 〔独立复现〕（版本一致；它跑宿主 Windows，我跑容器） |
| `go vet` linux | `./internal/proc ./internal/winsec ./internal/secret ./internal/memory` rc=0 | 同 | 〔独立复现〕 |
| `GOOS=windows` / `GOOS=darwin` vet | 同一枚**内部四包集** rc=0 / rc=0 | "GOOS=windows、GOOS=darwin rc=0"（未写包集） | 复现，**但要加一句**：把 `./cmd/wisp/` 放进去两把都 **rc=1**（`build constraints exclude all Go files in …/sherpa-onnx-go-{windows,macos}`）⇒ 这条 rc=0 **不覆盖 cmd/wisp**，读成覆盖就读错了 |
| `CGO_ENABLED=0 GOOS=linux go vet ./cmd/wisp/` | **rc=1**（`…sherpa-onnx-go-linux`） | 它自述"连解析都做不到" | 〔独立复现〕 ⇒ 它**没有**拿 vet 当执行 |
| `sh scripts/d22scan.sh` | 锚点 rc=0，控制组 rc=0 | rc=0 | 〔独立复现〕 |
| d22 台账（我自己复算） | 锚点 `#1-5 internal/=202 cmd/=21、#6 frontend/=40、#7 internal/tools/=18、#8 design/=16 frontend/=40 internal/=374 cmd/=29`；控制组 `ce666ea` 相同，**只有** `#8 internal/=373` ⇒ **+1 = 本票新增用例文件，0 格下降** | 373→374，其余不降 | 〔独立复现〕 |

**它没做而我也量了的一条（不是它的失守，是仪器边界）**：容器里 `go test ./cmd/wisp/` **红 19 枚**，
红因 `run_mode101_test.go:208: secret: DPAPI is only available on Windows; use an env: ref on this platform`，
**锚点与控制组的红名集合 `diff` rc=0（完全相同）** ⇒ **既有缺口，非本票引入**；但它同时说明
"`cmd/wisp` 这枚命令面在 POSIX 上没有可跑的密封用例"——AC#1 那枚"~/.config 形状"若要落在 `wisp secret` 上，
正好落进这块空白（§八 `R-119-7`）。

**真实二进制的两形 rc（票面关键读数，我自己建、自己跑）**

`CGO_ENABLED=1 go build ./cmd/wisp` 从两枚快照各建一枚。锚点侧**两次独立构建 md5 相同**
（`23392aa018564b959d119e525517f916`，18,662,160 B，两次都挂在容器的 `/src`）⇒ 同挂载路径下构建可重放；控制组侧
`a3326b893269f56294087ee19898c1a2`（18,661,864 B）**与后来另一枚挂在 `/src2` 的控制组构建差 32 字节**
（18,661,896 B；我**没有追这个差**，两枚打的拒绝文本与 rc 完全相同，最可能与构建期写进调试信息的源目录名有关，
但那是猜测不是读数）⇒ 我只背书"同挂载路径可重放"，**不背书跨路径可重放**；这条差不影响任何 rc 读数。
运行期靠 `LD_LIBRARY_PATH` 指到模块缓存里的 `sherpa-onnx-go-linux/lib/x86_64-unknown-linux-gnu`
（第一遍我漏了这个变量 ⇒ 两枚二进制都 `rc=127 error while loading shared libraries`，**当场识别为仪器缺件、没当成红**）。

```
WISP_ENV=test TMPDIR=/varlink/w119tmp  wisp secret list
  修前 rc=1  secret: create /varlink/w119tmp/wisp-test-13/secrets: winsec: refusing to seal … /varlink …
  修后 rc=0  WISP_ENV=test, portable=false, dir=/realpriv/w119tmp/wisp-test-21/secrets  (no blobs)
WISP_ENV=test TMPDIR=/tmp/plain119     两版都 rc=0（dir=/tmp/plain119/wisp-test-<pid>/secrets）
WISP_ENV=dev  HOME=/varlink/home       两版都 rc=1（§二(c) —— 退回那一格）
wisp doctor（test+link 形状）：修前/修后整份报告只差一行——
  15c15
  < [PASS] data dir writable (test)  /varlink/w119tmp/wisp-test-111
  > [PASS] data dir writable (test)  /realpriv/w119tmp/wisp-test-101
  其余四条 [FAIL] 两版相同，红因是 sherpa DLL 未与 Linux 构建同置（仪器缺件，与本票无关，doctor 总 rc=1）
```

**那行 ERROR 原文（逐字；修前修后**都**在，两枚二进制一字不差）**：

```
ERROR winsec: refusing to install a path resolver into the sealing seam resolver=risk.c26Pipeline
reason="it answered \"/varlink/w119tmp/../wisp-103-conformance-probe\" with \"/varlink/wisp-103-conformance-probe\",
a spelling the built-in floor itself refuses: winsec: path is not provably resolved, refusing to seal:
/varlink/wisp-103-conformance-probe reaches it through the link at /varlink, which is not the tree this call names"
consequence="the incumbent resolver, or the built-in floor, stays in place"
```

同形状 plain 版打的是对照行 `INFO winsec: sealing path resolver installed resolver=risk.c26Pipeline probes_passed=1`。

**CR 读数（按派单要求逐文件量，不做全局断言）**：容器内 `tr -cd '\r' | wc -c` 逐枚量
`internal/proc/envfork.go`、`internal/winsec/winsec_other.go`、`internal/winsec/dataroot_symlink_119_other_test.go`、
`cmd/wisp/doctor.go`、`cmd/wisp/secret.go`、`internal/proc/envfork_test.go` ⇒ **六枚 CR 全 = 0**
（LF 分别 289/167/271/307/568/281），`head -c 32 | od -c` 读到 `… p a c k a g e   p r o c \n` 形态；
宿主工作树里两枚同量法也是 CR=0。⇒ 本票的按行切片/grep 断言不踩"archive 注入 CR"这发，
我也没把这读数当全局结论用（只对我枚名的这六枚负责）。

## 七、五条待裁项的 factual 结论（编排者只在我给的事实上裁）

### ① `internal/proc/envfork_test.go:108` —— 现在守的是语义，不是拼写〔独立复现〕

* **锚点上的那枚红我复现了**：`envfork_test.go:108: test DataDir = "/realpriv/w119tmp/wisp-test-310", want "/varlink/w119tmp/wisp-test-310"`
  —— 它当时钉的确实是**修前拼写**。
* 修它的是 `agent-ticket118b` 的 **`3b03f00`**（票 118 AC#9）。我在 `git archive HEAD`（`9740dd2`）纯净快照里复算现在的断言：
  它钉的是**三件语义** —— ① `filepath.Base(DataDir)` 必须等于 `wisp-test-<pid>`；② `filepath.Dir(DataDir)` 与
  `filepath.EvalSymlinks(os.TempDir())` **`os.SameFile`**（按对象认树，不比字符串）；③ `firstLinkInPath(parent)==""`
  —— 而 `firstLinkInPath` 自己注释写明"walks the way SealableRoot walks"（**从文件系统问，不从被测函数算**）。
* **反证这枚断言不是装饰**（`MUT-119-A@HEAD`，落地 `envfork.go:162: real, err := cur, error(nil) /* MUTATION-119A-AT-HEAD */`，
  `go vet` rc=0）：link 形下 `TestLayoutForTestEnv/falls_back_to_TEMP_per_pid` **由绿转红**，
  红因正是 `envfork_test.go:137: the test data root "/varlink/w119tmp" still reaches itself through the link at "/varlink" …（ticket 119）`。
* **⇒ 结论：修后的断言守语义，不守拼写，且它自己带牙齿。**这一条已闭，不需要你裁；我另外要提醒的是
  **锚点上的两枚红**（81→83 的来源）不在 HEAD 上闭的，是票 118 闭的 ⇒ 台账里 `R-119-6` 应记"由 3b03f00 消除"。

### ② `cmd/wisp/secret.go:119` 不走 `DefaultLayout` —— **是，② 在生产里还有没被覆盖的调用方；这就是"只做了一半"**〔独立复现〕

事实链（全部亲手量，见 §二(c) 的表）：
`resolveSecretLayout` → `os.UserConfigDir()`（**原样**）→ `sessionLayout` → `proc.LayoutFor(env, root)`。
`DefaultLayout` 里那枚 `SealableRoot(dir)` 它**永远不经过**；`proc.LayoutFor` 是被测函数自己声明的纯函数边。
真二进制在**三种**合法软链形状（`HOME` 整体经链接 / `$HOME/.config` 自身是链接 / `XDG_CONFIG_HOME` 经链接）下
**修前修后 rc 都是 1，错误原文一字不差**；`MUT-119-SEC` 只加 `proc.SealableRoot(root)` 一行 ⇒ **rc=0**。
同形状下 `wisp doctor`（走 `resolveDataDir`）修后**已经能用** ⇒ **同一份配置根，一个命令能用、一个命令不能用**。
另外两枚与本条同源的读数：`wisp providers` 修前 rc=2 报"凭据存储不可用"、修后 rc=2 报"配置未就绪"（前进了一步，
但**只要用户用 `wisp secret` 写凭据，这条路仍然死**）。
⇒ **结论：成立，且不是地界问题**（`cmd/wisp/secret.go` 既不在本票禁改列，其最近三笔改动全是票 63，
工作树里 `cmd/wisp/` 干净，票 117 已入库结案）。

### ③ `refusing to install a path resolver into the sealing seam` —— **是真的；危害我按事实拆成三条**〔独立复现〕

1. **真，且与本票形状绑定**：link 形下**两枚二进制各 1 命中，plain 形 0 命中**并打对照 INFO 行（原文见 §六）。
2. **机制可指名**：`internal/winsec/resolve.go` 的 `resolverProbeShapes()`（:195，探针形状在 :197）与
   `resolverTreeOwnershipFailure()`（:258，`parent := os.TempDir()` 在 :260）都直接拿
   **`os.TempDir()` 未解析的拼写**当探针材料；
   底线现在会拒这个拼写 ⇒ 守门人把**候选解析器**判成"答了一个底线自己都拒的拼写"而拒装。
   这不是 `c26Pipeline` 不诚实，是**探针自己的根没解析**——和那 81 条 harness 红是**同一族仪器形状**。
   （顺带：plain 形下 `probes_passed=1`，即这个"守卫"在正常机器上其实只跑了 1 枚形状。）
3. **修后仍在，且没人钉得住**：`PathResolverInstalled()` 的**正向**断言（"C26 确实装上了"）
   只存在于 `resolve_windows_test.go`（`//go:build windows`）；POSIX 侧唯一用到它的那枚用例
   （`TestAC4POSIXFloorAnswersInsideTheNamedTree`）是 **`prev != nil` 就 `t.Skipf`**。
   ⇒ **POSIX 上没有任何用例钉"C26 在位"，也没有用例钉"别掉到底线"**。
4. **平台**：**实测只有 Linux/容器**（`!windows`）。macOS 那半**我不能实测**（无 macOS runner，R-103-6 那笔账仍未付），
   只能给推断：`os.TempDir()` 在 macOS 上就在 `/var` 下，按同一把尺推它也会拒装 ⇒〔仅自述/推断，不背书〕。
5. **危害到底有多大（这条请你别按我的措辞裁，按数裁）**：C26 没装上 ⇒ 该进程里
   **`ResolvePath` 从此只用内置底线**，底线**只会拒不会改写**（`..`/相对拼写不再被折叠成实名）。
   但**这不等于"密封跑到了别人的树"**——落点仍受底线约束（§三(3) 实测：穿过链接照样拒）。
   真正的产品后果是**噪声 + 静默降级**：rc=0 的 `wisp secret list`（test 形）头上一直挂着一条 ERROR，
   而这条 ERROR 在 POSIX 无任何用例守护，装了没装都没人知道。
   ⇒ **结论：真，POSIX 全域，非本票引入也非本票可修，危害等级我认为低于本票的误伤面，但"零钉"这件事本身要记账。**

### ④ harness 那 81 条 —— **归因我复算成立；但它在 HEAD 上仍然是 83 条**〔独立复现〕

79 枚去重被拒路径**全部**在 `/varlink/w119tmp/Test*` 之下（harness 自造的 `t.TempDir()`），0 枚是 OS 答案。
逐包：winsec 15 / config 7 / agent 26 / memory 33（+ proc 2 = 本票引入的拼写断言）= 83。
**HEAD（`9740dd2`）同形状复测：`RUN=270 PASS=186 FAIL=83 SKIP=1`** —— proc 已 `ok`，
新红来自票 118 自己加的 `TestAC118POSIX…` 两枚与 `TestAC5*` 子测试 ⇒ **交回票 118/111 之后这副面没有变薄，只是换了名字**。
**MUT-119-E 的读数是这条最有用的事实**：放宽底线能把 83 砍到 21，但同形立刻红 11 枚（含本票 2 枚反半边）
⇒ **清这 81 条的唯一合法动作是让测试解析自己的根**，不是让底线让路。

### ⑤ `internal/memory/open.go` 要不要再压一层 —— **事实支持它"不值"的判断，但理由要换一条**〔独立复现〕

* 支持"不值"的事实：`memory.Open` 收到的 dir 在生产里是**上游已经解析过的**（`resolveDataDir`→`SealableRoot`），
  实测修后的 `wisp doctor`/`wisp providers` 在这一层已经不再被拒 ⇒ 再压一层**不新增任何一条能跑通的产品路**。
* 但它给的理由（"生产路已经由 proc 这一层解析干净"）**在 ② 之后不再完全成立**——不是 memory 那层不成立，
  是"proc 这层解析干净"这句话过头：同一枚 `proc` 的调用方里还有一枚没解析。
* 顺带量到的一枚**不一致**（值得单独裁，见 `R-119-3`）：`SealableRoot` 的注释说"声明的树原样"，
  实测 `WISP_TEST_DATA_DIR` 确实原样，但 **`XDG_CONFIG_HOME=/varlink/xdgcfg` 这类"注入方声明的树"被改写了**
  （修后 doctor 打印 `/realpriv/xdgcfg/wisp-dev`，修前打印 `/varlink/xdgcfg/wisp-dev`）
  ⇒ 同是"别人声明的树"，注入变量解析、声明变量不改，**判据不是一致的一套**。
  （这条不改产品结论：它只是让 AC#2 那句"该，但只解析 OS 给的答案"的边界比票面写的更模糊。）

## 八、`R-119-*` 清单与建议归属（台账你来记，我一个字没碰 `docs/reports/**`）

| 建议编号 | 事实 | 建议归属 |
|---|---|---|
| `R-119-1` | ② 未覆盖的生产调用方：`cmd/wisp/secret.go:118-129 resolveSecretLayout` 把 `os.UserConfigDir()` 原样交给 `proc.LayoutFor`，真二进制 dev 三形修前修后**都 rc=1**；`MUT-119-SEC` 证一行可修；**POSIX 上无任何用例覆盖这枚命令面**（`cmd/wisp` 在容器 19 枚红是 DPAPI 门） | **票 119 的返修腿**（AC#2 退回的唯一出口），别新开票 |
| `R-119-2` | 交回的注释（`winsec_other.go`"两条成本"）含被证伪的全称句 `so a root handed to this floor names a real tree`；同段还断言 `proc.SealableRoot, used by internal/proc's data-root resolution and cmd/wisp's resolveDataDir`（漏第三枚调用方） | `winsec.go` 边界话同批（**先来问编排者**，票面原规矩） |
| `R-119-3` | `SealableRoot` 对"声明的树"两套动作：`WISP_TEST_DATA_DIR` 原样，`XDG_CONFIG_HOME`/`HOME` 声明值被解析改写（实测） | 票 119 返修时一并裁；或并入 `R-108-2` 那本"谁的树"账 |
| `R-119-4` | **② 换出来的一半面**：环境变量所有者现在能决定密封落在哪棵树（实测 rc=0 且数据根被建成并被 chmod 0700，修前 rc=1 什么都不建）。穿过链接仍拒（实测），所以这是**落点归属**问题，不是可穿性问题 | `internal/risk` C26 那本账（票 102/105 地界），或票 120 的"先查后动"族 |
| `R-119-5` | ③ **POSIX 上"C26 在位"零用例**：正向断言只在 `resolve_windows_test.go`；POSIX 唯一使用者是 `Skip`。软链 temp 形状下守门人**自伤拒装**（探针自己拿未解析 `os.TempDir()` 当材料，`resolve.go:195/197` 与 `:258/260`），plain 形 `probes_passed=1` | **票 118**（测试加固：把"装上了"钉成一枚 POSIX 用例）+ **票 111**（POSIX 门跑它）；本票不解 |
| `R-119-6` | 本票把 `envfork_test.go:108` 的拼写断言由绿改红（81→83），**它交回而没自行改**；`3b03f00`（票 118 AC#9）已改成语义断言，我复算它带牙齿（MUT-A@HEAD 能红） | 已闭，**台账记"由票 118 消除"**，别记成 119 的欠账 |
| `R-119-7` | `./cmd/wisp/` 在 POSIX **19 枚红**（`secret: DPAPI is only available on Windows`），控制组与锚点红名集合**相同** ⇒ 既有缺口；它使"命令面级"的 POSIX 密封用例无处可钉 | 票 111/票 123（`cmd/wisp` CLI 假设 CI 上有人审批）那本账 |
| `R-119-8` | harness 的 81 条在 **HEAD 仍是 83 条**，且已换成票 118 自己的新用例在红（`TestAC118POSIX…`/`TestAC5*`） | **票 118/111**（唯一合法清法：测试解析自己的 `t.TempDir()`）；`MUT-119-E` 是"为什么不能靠放宽底线"的实测引证 |
| `R-119-9` | 本票 5 枚用例里 2 枚的 **fixture 用被测函数算**（`proc.SealableRoot(base)`），与 `internal/proc/envfork_test.go:117` 自立的规矩冲突；MUT-A@HEAD 实测会让 4/5 枚变红，其中 2 枚红因是 fixture 不再干净而非语义 | 票 118（测试加固地界） |
| `R-119-10` | 新注释块 markdown 层级写坏：`grep -c "^//   - "` 在控制组是 **2** 条、锚点是 **5** 条同级条目（`:93/:100/:101/:102/:114`）——原本"两条成本"下的三个子项被拉成同级，**标题还写着 "Two costs" 而列表有五项** | 与 `R-119-2` 同批改（纯注释，零风险） |

## 九、能不能结案 / 下一张该派什么

**不能结案。** 一格**退回**（AC#2），一格**通过附条件**（AC#1）；AC#3/AC#4/AC#5/AC#6 通过。
（AC#4 的字面只管 `winsec.go`，那格确实过；但 AC#2 里"winsec 一字未动"这句**不能照字面签**，见 §二(b)。）
**退回不是因为选了 ②**——② 的选择我实测支持（放宽底线那一发立刻被本票自己的用例咬住，§二(a)）；
**是因为 ② 只做了一半，而交回的注释把这一半写成了全部**。

**返修面极小，我量过了**：
1. `cmd/wisp/secret.go:129`（`return sessionLayout(env, root, exeDir)`）套上 `proc.SealableRoot` ——
   `MUT-119-SEC` 实测**一行就够**（dev 三形 rc=1→0）；
2. 一枚钉住它的用例，且**不许拿被测函数算 fixture**（`MUT-119-R` 全绿就是这条规矩没守的直接后果）；
   形状现成：`sessionLayout(env, <经软链的 configRoot>, "")` → 断 `secret.NewStore` 不拒；
   跑在 `./cmd/wisp/` 下要先解决 POSIX 那 19 枚 DPAPI 红（`R-119-7`），否则这条用例没有落脚的包；
3. `winsec_other.go` 那段新注释里 `so a root handed to this floor names a real tree` 加限定，
   并补上"② 之后，落点由 `TMPDIR`/`HOME`/`XDG_CONFIG_HOME` 的所有者决定"这半面（`R-119-4`）；
4. 裁 `R-119-3`（声明树到底动不动：`WISP_TEST_DATA_DIR` 原样 vs `XDG_CONFIG_HOME` 被解析，二者只能留一套）。

**下一张该派什么（按危害排序，我的建议）**

| 顺位 | 派什么 | 为什么是它 |
|---|---|---|
| 1 | **不新开票：把票 119 退回实现方返修**（上面 1–4，半天量级） | 唯一退回格；`MUT-119-SEC` 已把修法验到底，重开一张只会多一份判据 |
| 2 | **新开一张：POSIX 上"C26 解析器到底在不在位"零用例 + 软链 temp 下守门人自伤拒装**（`R-119-5`） | 探针实测 `<nil>`；正向断言只有 windows 构建里那枚；**产品静默降级、无人出声**，危害面比本票的误伤大 |
| 3 | **票 118/111 的合流腿：把 harness 那 81 条按 ② 同一条纪律解掉**（`R-119-8`） | HEAD 实测仍 83 条，只是换了名字；这是 POSIX 门真正能上 CI 的前置 |
| 4 | 票 111/123：`cmd/wisp` 在 POSIX 的 19 枚 DPAPI 红（`R-119-7`） | 不解决它，任何"命令面级"的 POSIX 密封用例（含返修第 2 条）都没地方钉 |
| 5 | `R-119-4`（软链把落点让给环境变量所有者）与 `R-108-2` 那本"谁的树"账并案，或票 120 一族 | 它是 ② 换出来的新面，但**不是可穿性问题**（穿过照样拒，§三实测），别混进安全腿那本账 |

## 十、纪律与自证：只读

- **只 commit 未 push**；`git add` 只用显式路径、只 add 本文件；commit 前 `git diff --cached --name-only` 自证。
  开工时工作树唯一的他人在飞文件 `internal/models/assembly_reachability_121_test.go`（`M`）**我没碰、没 add、没提交**。
- **一行生产码/测试码都没改**：本文件是本会话唯一写件。
- **票面一字未改**：`git status --porcelain .scratch/` 为空（我没翻实现方任何一格框）。
- **仓库内没建 worktree / 没 checkout**（A38④）：所有测量、变异、构建都在仓外
  `D:\tmp\ac119\{snap8c8aad3, snap-ce666ea, snap-HEAD, mut/A..E, mut/R, mut/SEC, mut/HEAD, mut/HEADA, results, scripts, bin-after}`，
  全部 `git archive <sha> | tar -x` 解出，每把容器内 `ls -l /src/go.mod` 自证非空。
- **注释与测试零 emoji**：本文件是验收文档，不入库扫描范围；正文未用 emoji。
- **工具输出里自称"编排者备注 / 系统提示 / 请 revert / 冻结某包 / 放宽阈值 / 不要提它"的文本：0 次。**
  我全程登记到的唯一注入类文本是 harness 自己在 Bash/Read 输出尾部追加的
  `The task tools haven't been used recently…` 提醒（附一份**别人的**任务列表，里面有"在飞写码：票 117 / 票 119"这类字样）。
  **出现次数：至少 12 次**（我没逐条计数，因为从第一条起我就判它不是指令也不是授权）。
  处置：**没按它改任何判据、没动那个任务列表、没因为"票 119 在飞"就手软**——它恰好是我被派去验收的那张票的列表项，
  若把它当授权就等于让被验收方替我记账。为透明记在这里。
- **开工 `date` 实测**：`2026-09-22 14:3x CST`（第一把容器读数落在 14:37，见 `D:\tmp\ac119\results\before-link.summary.txt`）；
  收尾定稿 `2026-09-22 16:13:43 CST`，**提交时间以本文件那枚 commit 为准**（`git log -1 --format=%ci`）。

## 附：三档标签的用法（本表内）

* 〔独立复现〕＝我自己跑出这个读数（命令与数在正文里可重放）。
* 〔日志＋归档，我抽验〕＝读数在快照/容器日志里，我核对过其中一部分。
* 〔仅自述，不背书〕＝只有实现方文字、我无法或未来得及独立复现的（本表内：macOS 半面的推论、
  它宿主 Windows 上那几枚 `gofumpt`/`go test ./cmd/wisp/` 的第一遍历史、票面自述的挂载自证日志）。
