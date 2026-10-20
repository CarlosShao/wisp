# 票 353 `AC#2` 基线侧读数（`353-r1`＝★写腿·乙形·修前＝乙未落时的现形）

尺＝`sh 02-run-pair.sh host base353 <三枚定向 -run> nametopup_353_windows_test.go`（本件同目录 `sh/`）；射程＝**仓外导出树物理字节**（锚 `a59f5630`，三枚承重 blob 与 census/断腿锚逐枚同串＝`winsec.go f3c1b68e`⇄`pinnedleaf_342_windows.go db077e30`⇄`winsec_windows.go 63714a83`〔量的，尺＝`git rev-parse <锚>:<path>`〕）；★量的（跑出来的分子），⛔ 算的。⛔ 共享工作树跑过任何 `go test`。

台面对齐＝go1.27.1 windows/amd64＋`GOTMPDIR`⇄`TMPDIR`＝反斜杠短形落 `%LOCALAPPDATA%` 同卷且先 `mkdir -p`；PATH 用 `cygpath -w` 只加仓外那两枚 DLL 段（★⛔ `cygpath -wp` 转整条 PATH）。两侧树在**这一趟字节相同**（＝乙未落），故两发是同形复量、⛔ 成对判据。

## 一、四数（`^=== RUN`／`^--- PASS`／`^--- FAIL`／`^--- SKIP`，⛔ 混用⛔ 相减）

| 发 | 树 | RC | RUN | PASS_TOP | FAIL_TOP | SKIP | SHAPE | 名册 |
|---|---|---|---|---|---|---|---|---|
| s1 | treeB（＝修前现形） | 1 | 3 | 2 | **1** | 0 | OK | 1 枚＝`TestAC2WindowsNameTopUpBeforeTheReCheckLeavesTheTopUpObjectNamed353r1` |
| s2 | treeL（同字节，复量） | 1 | 3 | 2 | **1** | 0 | OK | 同枚 |

名册尺＝`grep -E '^--- FAIL'`（逐字见 `logs/base353-s1-treeB.rednames.txt`）；原始逐发＝`logs/base353-s1-treeB.raw.txt`⇄`...s2-treeL.raw.txt`。

## 二、两枚分子（★到达数，⛔ 拿"总数大"当到达）

- **`entered=3/3` ⇄ 同发 `realSealRanFirst=true`**（liveness 钥匙＝逐发同发在场；`plantErr=<none>`，⛔ 借 `truncated`／`deleted` 任何一枚分子）＝**世界进得去**＝本票要买的那枚世界在这台机器确定性可复现（342-r2 台架行 S4：落地句柄带 `FILE_SHARE_DELETE`）。
- **`damage`（顶进来的对象被删到⛔ 剩一枚名）＝`3/3`** ＝★**修前红**＝票 353 那一跳的伤，逐发都是。
- 对照那枚既有钉（同发跑进来了）＝`ancestorwindow_342r2_windows_test.go:379` 的 `AC#1-one over 3 shots: worlds entered=3, ... unnamed objects left with no name=3` ＋ `:382` 的 RESIDUAL OBSERVED 分支照旧打印＝census 第 2 节说的"乙之后这枚钉要重述"的世界**今天真在**。

逐字分子行（s1，`nametopup_353_windows_test.go:146`）＝

```
353r1 AC#2 over 3 shots: entered=3/3, worlds the platform froze=0, unnamed objects left with no name (DAMAGE)=3, refusals that left the caller's name standing=0
```

s2 同形（逐字同串＝三发都 entered、都 damage，⛔ 竞态行⛔ 要求逐字同号这条本趟用不上，因为它是确定性仪器）。

## 三、★逐字红句（修前＝乙未落）

`nametopup_353_windows_test.go:150`（s1 ⇄ s2 逐字同串）＝

```
353r1 AC#2 RED: 3 of 3 entered worlds left the swapped-in object with NO name, so the refusal-path removal ran on a name that carried an object the caller never named - ticket 353's damage, on the platform where the route's own handle makes that name vacate-able (bench row S4)
```

逐发 liveness 行（三发同形，s1 `:109`）＝

```
353r1 AC#2 shot 1 LIVENESS+REFUSAL: leg=seal refusal on this route (case one injects it through applyDescriptor and names it as injected; case two never reaches this leg) | realSealRanFirst=true windowEntered=true plantErr=<none> | isNotSealable=true | verbatim=winsec: cannot apply a private security descriptor: C:\Users\swq\AppData\Local\wisp-test-gotmp\353r1-20261020-200542-742039\TestAC2WindowsNameTopUpBeforeTheReCheckLeavesTheTopUpObjectNamed693340766\001\leaf-342r2w.bin: injected through this package's declared AC#5 seam: the descriptor could not be applied
```

⇒ 这一列同时钉住三件事：ⓐ 注入走的是本包**已声明的** `applyDescriptor` 缝（⛔ 新缝）、ⓑ 真 seal 先跑完并保留其判语、ⓒ 拒的是 seal 那一支（AC#5 哨兵在场），所以 damage 那一步确实是**拒绝路径**的删除、⛔ 成功路径那一支。

## 四、正控（同发跑进来、照绿）——这一枚是 damage=0 那把尺的牙

`s1 :182` 逐字（★`kind=absent`＝删除**真的在跑**，所以 damage 那枚 0 将来⛔ 是"尺⛔ 咬东西"）＝

```
353r1 AC#2 CONTROL (no top-up, name still carries the object the route opened): isNotSealable=true | leg=seal refusal on this route (...) | name=absent (GetFileAttributesEx C:\Users\...\no-top-up-control-353r1.bin: The system cannot find the file specified.) | dir entries []
```

⇒ 判据形状＝**没有顶入的世界，删除照旧执行**（AC#5"⛔ 留半枚文件"三形都⛔ 改这一条）。乙落地后这一枚**必须照绿**＝它和主用例合起来才把"damage 归 0"夹成"复验拦住了删"、⛔ 夹成"删除没了"。

## 五、★写进判据的那味：窗内那一支⛔ 测不到

写进**用例注释 ⇄ 用例自己打印的 GREEN SCOPE 行**两处，⛔ 只在散件里写：

- 文件＝`internal/winsec/nametopup_353_windows_test.go`，★注释块＝第 `29-38` 行（`★WHAT THIS INSTRUMENT CANNOT COVER ...`，逐字写着"NOT refused and NOT measurable here ... a hook INSIDE the platform's `NtOpenFile`"）；
- ★**运行期逐字打印**＝同一枚文件里那条 `AC#2 GREEN SCOPE` 行＝第 `157-163` 行（落在 `entered==3 && damage==0` 那一支＝**只在仪器绿的时候才打印**＝绿行自己带着⛔ 覆盖窗内那一支的声明）；
- 另两处＝同文件 `:117`（`!entered` 那一支的红句：把"⛔ 进得去"写成红、⛔ 冒充绿＝基线侧那 `0/3` 的凭据形状）与 `:166`⇄`:172-195` 正控用例（⛔ 咬东西的尺⛔ 许报绿）。

★行号口径＝上面这几枚是**笔④（纯注释修 ban#9 拼法）之后**的终量行号，尺＝`sh/28-delivery-rulers.sh` F 段现跑；本节上方写着的"逐字写着 (`internal/syscall/windows/at_windows.go:247`)"那形引用已在笔④改成中性同义写法（＝Go 运行时自家的件、⛔ 是仓相对路径；d22scan ban #9 只判仓相对拼法），判据内容一字⛔ 减。交件以 `readings-ac2-landed.md` 第七节的四处落点为准。

⇒ 形状＝票 322 `AC#1` 的"报窗口、⛔ 圆窗口"先例；⛔ 让仪器在它⛔ 覆盖的世界上报绿＝census 必答③ 对乙的判语（"先顶入＝确定性可证伪；窗内那一支⛔ 可证伪＝要钩在 `Deleteat` 的 `NtOpenFile` **之内**＝平台内部"）。

## 六、★基线侧两枚分子的最终复量（笔④后同尺重取＝交件用的那对数在此）

同一条尺换到最终台面（锚 `a59f5630` 重 archive＋HEAD 的测试字节，产码＝基线 blob）再取一遍＝`entered=3/3` ⇄ `DAMAGE=3` ⇄ 正控 `name=absent` 照绿，逐字分子行⇄红句见 `readings-ac2-landed.md` 第二节 finalB 那一列（＝本节那两枚数在新台面上逐字复现＝⛔ 一次读数当定论）。本节上方那两发（`base353-s1/s2`）＝同一枚基线字节在**上一版台面**上的读数，两发同号；保留⛔ 删（临时件只建不删），交件以第二台面的数为准。

## 七、纪律

框终量＝票 353 **未勾 3／已勾 2**（`AC#0`⇄`AC#1` 已由编排者翻；尺＝`git cat-file blob HEAD:<票> | grep -cE '^- \[ \]'`⇄`'^- \[x\]'` 现量，射程＝该票面）＝本腿一枚⛔ 翻。票名／`Status:`／`-done` 一字⛔ 动；票尾那条唯一追加已被 `353-a1` 占额⇒本腿话落本 probes 件、⛔ 补第二条 bullet。⛔ push。临时件只建不删＝`D:\tmp\wisp\20261020\353r1\{treeB,treeL,treeB2,finalB,finalL,mut-final-*,pair*,incoming,base,logs,gocache-win,gocache-col-*}`。写面＝`internal/winsec/nametopup_353_windows_test.go`（新增）＋本腿 probes 件；`cmd/wisp`⇄`internal/models`⇄`internal/risk/**` 字节 0 枚。★本件＝笔①入库的修前记录；交件终量＝同目录 `readings-ac2-landed.md`。
