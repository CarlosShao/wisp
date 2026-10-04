# 票 263 r1 — `scripts/slo-check.ps1` 取未赋值 `$LASTEXITCODE` 致 D32 门当场死：成因、修法、正控与钉子

落地腿：`263-r1`。写面＝`scripts/slo-check.ps1` ＋ 本目录台件。全程中文。
本节骨架于开工第 9 轮先提交（派单「节奏」条）；交件发零占位标记＝原派单点名的那三枚词在本文与本目录所有交证件里一枚都不在，
尺与命中数见 `placeholder-check.txt`（pattern 本身不在这里逐字写出来，否则"尺自我命中"会把命中数变成 1——这个坑与本文件 §⑤ 的词面钉同一个形状）。

- 起手锚点：分支 `dev`，HEAD `2ae018be`，取数与落笔同发（2026-10-04 11:11:58 +0800）
- 起手写面真空：`git status --porcelain -- scripts .github` ＝ 空（同发取，`262-r1` 已收档）
- 交件时我的两枚产码 commit：`b92037bd`（主修＋台件）、`cb8dcecc`（词面钉扩形＋反形台件）
- 台件生成的变异副本（`mutations/slo-check-*.ps1`）与所有采样输出目录（`out-*`、`mutations/out-*`）**留在盘上不入库**：
  临时件只建不删，副本可由 `run-nail-mutations.ps1` 从 tracked 字节逐字节重生成（锚不命中即 throw），
  入库只会造成"两份 slo-check.ps1"的第二真相源——这正是本仓禁的形状。逐枚原始控制台读数都以日志形式入库（`*.txt`/`*.log`）。
- 本文写作钟：2026-10-04 11:48 +0800；下方所有读数都是本腿 11:1x–11:4x 现跑，⛔ 别当常量引用
- 在飞写腿：`260-r3`（`cmd/wisp`＋`internal/agent/approval`）⇒ 本腿⛔ 未跑 `-Subset full` 全量、⛔ 未跑整包 `go test`、⛔ 未碰 `build/` 里他人产物（exe 全部来自本目录 `bin/`，输出全部进本目录 `out-*`/`mutations/out-*`）

---

## ① 现量复认（AC#0）

**1. 票面 §现量 4 的尺，我自己重跑（11:11:58 +0800）**

```
wc -l scripts/slo-check.ps1                      -> 397 行
grep -n 'LASTEXITCODE' scripts/slo-check.ps1     -> 4 处：:106 :326 :345 :366
grep -n 'Set-StrictMode' scripts/slo-check.ps1   -> 1 处：:65
git status --porcelain -- scripts .github        -> 空
```
逐字读数（改前，`git cat-file blob HEAD:scripts/slo-check.ps1` 抽出的 397 行副本，同发 `grep -c` ＝ 4 处）：
- `:106` `    if ($LASTEXITCODE -ne 0) { Fail 'build.ps1 failed' }`
- `:326` `    $code = $LASTEXITCODE`
- `:345` `$settleCode = $LASTEXITCODE`
- `:366` `$leakCode = $LASTEXITCODE`
- `:65`  `Set-StrictMode -Version 2.0`

⚠ **派单里两处尺与盘上不符，具名更正（不改派单、不改 ci.yml）**：
1. 派单说 `ci.yml:585` 命令行是 `-Subset smoke -SecondsPerState 6`；盘上逐字（`sed -n '585p'`）是
   `run: powershell -NoProfile -ExecutionPolicy Bypass -File scripts/slo-check.ps1 -Subset smoke -SecondsPerState 4`
   ⇒ smoke 那档是 **4 秒**，不是 6。`:648` 逐字＝`... -Subset full -SecondsPerState 6`（这条派单读对了）。
2. 派单说 `grep -n Set-StrictMode` 只有 `:65` 一处并另有 `:373` 过期指认 ⇒ 两条都复认成立，
   `:373` 逐字：`# Force the filter into an array with @(): under Set-StrictMode 2.0 (line 48)`（StrictMode 真身在 `:65`，指认偏 17 行）。

**2. `gh run view 37166458550` 现跑（复认 §现量 1）**

```
databaseId 37166458550  headSha fd269de12fa1e6e2c1c3ea7e61bb4338a5f8c847
conclusion failure      createdAt 2026-10-04T00:56:01Z  workflow ci
slo-full  failure:  Set up job success / checkout@v4 success / setup-go@v5 success /
                    Build wisp.exe (deps cached on the runner) success /
                    SLO full gate (six states + settle + leak) failure /
                    Upload SLO report skipped
slo-smoke failure:  Set up job success / checkout@v4 success / setup-go@v5 success /
                    Cache third_party success / Build wisp.exe success /
                    SLO smoke gate failure / Upload SLO report skipped
```
⇒ §现量 1 成立：checkout 那层墙今天过了，红在门这一步；`Upload SLO report` 仍 `skipped` ⇒ 那次 run 没有任何 SLO 数字。

**3. 日志逐字（`.scratch/ci-logs/run-37166458550-failed.log`，行号是本腿现跑的 grep 命中行）**

`grep -n "slo-check.ps1\|precheck ok\|cannot be retrieved" <那份日志>` 命中 `:8734-8748`（slo-full）与 `:8719-8734`（slo-smoke）。
slo-full 那段逐字（ANSI 已剥、时间戳保留）：
```
00:56:57.8757289Z ##[group]Run powershell ... slo-check.ps1 -Subset full -SecondsPerState 6
00:56:58.9540673Z slo-check.ps1: sampling validity precheck (ticket 134 AC#4)
00:57:02.2618570Z slo-check.ps1: precheck ok - no foreign toolchain/runner process, machine-wide cpu max 17%
00:57:02.2640129Z slo-check.ps1: sampling state Sleeping for 6s
00:57:02.2943006Z time=2026-10-04T08:57:02.293+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
00:57:02.4597483Z E:\work\base\actions-runner\_work\wisp\wisp\scripts\slo-check.ps1 : The variable '$LASTEXITCODE' cannot be retrieved because it has not been set.
00:57:02.4599920Z     + CategoryInfo          : InvalidOperation: (LASTEXITCODE:String) [slo-check.ps1], RuntimeException
00:57:02.4601599Z     + FullyQualifiedErrorId : VariableIsUndefined,slo-check.ps1
00:57:07.7518268Z ##[error]Process completed with exit code 1.
```
⇒ §现量 2 成立（precheck 真在拦：`cpu max 17%`、无 foreign 进程 ⇒ 这枚红不是抢 CPU）。
⇒ §现量 3 成立并**加一条本腿新读数**：从 `sampling state Sleeping for 6s`（`02.264`）到红句（`02.459`）＝**195 ms**，
而那一档要求 6 秒；slo-smoke 那段同形是 **70 ms**（`04.426`→`04.496`）。
⇒ **PowerShell 根本没等那枚进程**——这一条是后面 §② 的第一手证据，不是推论。
⛔ 我没有拿今天 CI 那发当"我的正控读数"；§④ 的改前红是我自己造的。

**4. `slo-smoke` 那一发（票面说未归因、本票不许顺手提它）**

同一次 run 的 `slo-smoke` 步骤日志（`:8719-8731`）里有一发红句，**形状与本枚逐字同**（只有路径与文件标签不同）：
```
00:57:04.4260305Z slo-check.ps1: sampling state Sleeping for 4s
00:57:04.4964178Z D:\a\wisp\wisp\scripts\slo-check.ps1 : The variable '$LASTEXITCODE' cannot be retrieved because it has not been set.
```
⇒ 本腿**不据它下归因结论、不为它单独改任何东西**（票面禁）。只具名一句事实：两步跑的是同一枚文件，
所以本票的改动客观上也覆盖 smoke 那一步的取值点；下一步它什么颜色，要等下一次真 run 才算读数（AC#5 射程，归编排者）。

---

## ② 肇事实点与到达路径（AC#0 真实读取点／AC#1 成因）

**真实读取点＝改前 `scripts/slo-check.ps1:326`（逐字 `    $code = $LASTEXITCODE`），不是票面 §现量 4 那张地图里的 :106/:345/:366。**

排除法（不照抄票面四行）：
- `:106` 在 precheck **之前**（原文件 `:102-107`），而日志已打出 `precheck ok` ⇒ 已经走过 :108，:106 今天没被执行到。
- `:345`（settle）与 `:366`（leak）都在逐档循环**之后**；红句落在第一档 `Sleeping` 与 `Warm` 之间 ⇒ 未到达。
- `:326` 恰在 `:324 Write-Host 'sampling state ...'` 与 `:325 & $WispExe slo ...` 之后一行 ⇒ 与日志顺序唯一吻合。
- 复认：票面 §现量 3 自己写的"肇事读取点不是那三行"这条待验断言，实测成立。

**为什么这一取会炸（机制链，每一环都有自跑读数）**

1. `scripts/build.ps1:115` 在 `$ldflags` 数组里带 `"-H=windowsgui"`（逐字；台账 A548／commit `cc6eaa65`＝票 244-r1 落的地）。
   `.github/workflows/ci.yml:644-645` 的 `Build wisp.exe (deps cached on the runner)` 那一步就是 `powershell ... -File scripts/build.ps1 -Env dev`
   ⇒ **CI 用的 wisp.exe 是 Windows GUI 子系统二进制**。
2. PowerShell 对 GUI 子系统的外部程序**不等待**（`&` 立即返回），因此 `$LASTEXITCODE` 这一自动变量**从未在该作用域被写过**；
   `Set-StrictMode -Version 2.0` 下"取从未赋值的变量"＝运行时异常 ⇒ 不是"忘了赋值"，是**等都没等过，就没有码可取**。
3. 台件对照（同一份 Go 源码，两种子系统，`objdump -p` 读 PE 头，不读注释）：
```
bin/probe-cui.exe  Subsystem 00000003 (Windows CUI)
bin/probe-gui.exe  Subsystem 00000002 (Windows GUI)
```
   `repro-wait.ps1`（把改前 `:324/:325/:326` 三条语句逐条搬出来的最小件）：
```
CUI 那发：repro: sampling state Sleeping for 2s → probe263: started → probe263: exiting → repro: state exit=3 elapsed-ms=2049   rc=0
GUI 那发：repro: sampling state Sleeping for 2s → probe263: started → The variable '$LASTEXITCODE' cannot be retrieved because it has not been set.   rc=1
pwsh 7.6.0 跑 GUI 那发：同形红句（不是 5.1 的怪癖，两个宿主同语义）
```
4. **拿改前的门本体复现**（`.scratch/wisp/probes/263/r1/slo-check-before-fix.ps1`＝`git cat-file blob HEAD:` 抽的 397 行，4 处读取未动）
   对着 GUI 子系统的假 wisp 跑 ⇒ 逐字重现 CI 那句红（§④ 附全文）。
   ⇒ 因果链闭合：**票 244 把出厂件换成 GUI 子系统 ⇒ 票 263 这枚取值从"能取到"变成"从来没被写过" ⇒ 门开采 0.2 秒即死。**
   这也解释了"两层故障互相掩盖"：262 的改名拆掉 checkout 墙之前，这条路径一次都没被踩到。

**这条路径今天为什么走到（precheck 通过后脚本干的第一件事）**

precheck 打出 `precheck ok` 之后（改前 `:314`），脚本按序只做：`:316` 组 `$states` → `:322` 进 foreach → `:324` 打印
`sampling state Sleeping for 6s` → `:325` **启动本会话第一条外部程序**（前面 `:102-:311` 全是 cmdlet：`Test-Path`／`New-Item`／
`Get-ChildItem`／`Get-CimInstance`／`ConvertTo-Json`／`Set-Content`，都不写 `$LASTEXITCODE`）→ `:326` 取值 ⇒ 死。
⇒ 净结果：D32 那两个数字的唯一求值路径，在"第一次采样还没被等完"的那一瞬间就取了一次不存在的数。

**两个附带后果（都量了／标清了）**

- 孤儿进程：`repro-leftover.ps1`（GUI 件、寿命 6s、改前形状）逐字：
  `leftover-probe: THROW - ...` 之后 `probe-gui processes still alive right after the death = 1`（pid=8592），等它退完才归 0。
  ⇒ 改前那形**死后子进程还在跑**；改前脚本 `:153` 的 `$loadNames` 里含 `'wisp.exe'`，所以这台机器上下一次采样的 precheck
  会把它当 foreign contention 而拒采。CI 上那发 wisp.exe 的存活时长我没量到 〔未实测，只有台件同形〕。
  ⚠ 我这把尺第一次跑过滤名写成 `probe263`（Go 程序里的名字）得到 `0 alive`＝**坏尺**，改成 exe 基名 `probe-gui` 才读到 1；已就地修并留此记录。
- 改前 `:106`（缺件自建那支）今天**不炸**，实测理由：`repro-child-parent.ps1` 对两种子脚本形状——
  `child-exit.ps1`（只 `exit 3`）得 `repro-parent: read exit=3`，`child-native.ps1`（只跑 `cmd.exe /c exit 5`、不写 exit）得 `read exit=5`
  ⇒ 被调 `.ps1` 会把码留给调用方。第三形（build.ps1 自己异常终止）会沿 `$ErrorActionPreference='Stop'` 直接终止本脚本、根本走不到取值〔读码判断，未实测〕。
  ⇒ 所以 `:106` 是**同一形状的第二枚潜伏点**，本票一并换形（票面 §现量 4 把它列在四枚候选里，不是新扩射程）。

---

## ③ 修法与四条禁自我点名（AC#2）

**形状＝票面许可的第二支"跑完外部命令才读"，不是第一支"先赋再取"。**

新增（行号为交付版 `scripts/slo-check.ps1`，共 625 行）：
- `:143 function Invoke-ExternalProgram` —— `ProcessStartInfo` 起进程、`$proc.WaitForExit()`（`:193`）等它、
  返回 `@{ ok; code; why }`（成功形 `:194`）。跑不起来／拿不到码 ⇒ `ok=$false` ＋具名 `why`，**绝不返回一个没来路的数**。
- `:131 Get-NativeArgumentLine` —— `ProcessStartInfo` 只吃一条命令行，本仓路径含空格（`projects plans`），带空格/引号的参数逐枚加引。
- `:202 Get-SelfHostExecutable` —— 供"缺件自建"那支找到当前宿主 exe，把 `scripts/build.ps1` 也当**外部程序**起＋等。
- `:214 Test-ExitCodeInstrument`（能力钉）／`:233 Test-ScriptShapeNail`（词面钉），`:291`/`:292` 在任何采样之前调用。
- `:94 Fail-InstrumentBroken`（诚实档，见 §⑥）。
四个取值点全部改走它：`:307`/`:312` 缺件自建、`:533`/`:539` 逐档采样、`:557`/`:563` settle 行、`:583`/`:593` leak 自检。
交付文件里 `$LASTEXITCODE` 字面出现 **1 次**（`:252` 那枚模式串本身，不带 sigil ⇒ 钉不匹配自己），读取 **0 次**。

**为什么不走"先赋再取"（实测理由，不是偏好）**
`repro-probe-slow.ps1` 两发逐字：不等而直读 `ExitCode` 时，.NET 不抛异常，**读回 0**（`NO-WAIT read succeeded exit=0 (hasExited=False)`，
`cmd-fast` 与 `cmd-slow` 两形都是 0）。⇒ 若只补一句初始值，门会拿到"每档都是 0"这种没来路的码，
而 `-out` 文件那时还没写出来 ⇒ 每档 `pass=False` ⇒ **门常红、永远不测**——正好是票面 §现状 要避免的那件事的反面（把"测不到"洗成"看起来测过"）。
等待才是缺陷的本体；顺带这也是 §⑤ 那枚能力钉测得到"没等"的原因（没等 ⇒ 0 ⇒ 不等于 7 ⇒ 当场具名红）。

**`UseShellExecute` 选 `$false` 的实测理由**：四发对照（`repro-helper.ps1` × CUI/GUI × true/false）——
两种取值都能等到正确退出码（`exit=3`），但 `UseShellExecute=true` 那两发**看不到子进程自己打的行**
（`probe263: started/exiting` 两行缺失）⇒ CI 日志会失去 wisp 的输出。`$false` 那两发退出码对、日志行也在。

**四条 ⛔ 逐条自查并点名"我没做哪一种偷懒"**

1. **没用 `try/catch` 吞**：新代码里两枚 catch 在 `Invoke-ExternalProgram` 内（`:183` 起不动、`:195` 等不到／取不到码），
   全部只做"把运行时异常换成具名失败"，返回 `ok=$false`＋`why`，调用方一律走 `Fail-InstrumentBroken` ⇒ `exit 1`。
   没有任何一支 catch 之后继续往下跑。
   `Fail-InstrumentBroken` 内部那一枚 catch（`:122`）只兜"盘上记录写不出去"，判定与退出码不丢。
2. **没放宽 `Set-StrictMode`**：交付版 `:82` 逐字 `Set-StrictMode -Version 2.0`（改前在 `:65`，只是被上面新增的注释块推下 17 行），
   一字未改；`$ErrorActionPreference='Stop'` 也未改。**没有**为绕开这枚缺陷加 `-Version 1`／关掉严格档／包一层 SilentlyContinue。
3. **阈值/golden/`internal/observe/thresholds.go` 一字节未动**：见 `git show --stat b92037bd cb8dcecc`——
   改动文件只有 `scripts/slo-check.ps1` 与本目录台件；`git status --porcelain -- scripts` 交件为空。
   脚本里出现的 `Sleeping CPU <=0.5%`／`RSS <=25MB` 全是**注释与文案里的引用文字**（`:29`、`:343`、`:449`、`:464-465`），
   求值仍在 `wisp slo` 里，本票没碰判定线。
4. **没把任何"必须出结论"的一步改成跳过**：新增分支的方向一律是 `exit 1`（更坏＝更响）；
   既有的 `NO CONCLUSION (machine-contended)`＋`exit 0`（交付版 `:437-519` 那一段）一字未动；
   没有加 `-Skip`/`-Force`/`if:`/continue-on-error 任何一枚（D22 mode-6）。
   ⛔ 我也**没有**顺手把 smoke/full 之外的第三档、或 `ValidateSet('smoke','full')`（交付版 `:74`）改掉。

**单列一行（票面之外、在我写面里的第二处改动，理由摊开）**
交付版 `:600-601`：把改前 `:373` 那句 `under Set-StrictMode 2.0 (line 48)` 改成 `under Set-StrictMode 2.0 (the one
Set-StrictMode line near the top of this file)`。理由：那句行号本来就是错的（真身 `:65`，偏 17 行），而我这一改又把行号整体推走——
**枚号型指认在我这发之后必然再过期**，所以换成内容锚。这一族（过期指认）归票 212 射程，我只改自己文件里这一句指认，
⛔ 没动 `tools/d22scan` 一字。⚠ 顺带上报一条读数（不修）：`tools/d22scan/main.go` 的 scope 名册只走
`internal/`、`cmd/`、`frontend/`、`design/`、`internal/tools/`——**`scripts/*.ps1` 不在任何一把尺的射程里**，
所以 ps1 里的过期指认今天无名册尺；这条判断是读 `main.go` 的 scope 表得出的，我**没有**跑"往 scripts 里种一枚过期指认看门响不响"的正控〔未实测〕。

---

## ④ 正控读数（AC#3，四发齐＋还原证明）

**① 改前那发红句，我自己造（⛔ 没拿 CI 那发当读数）**

台件：`slo-check-before-fix.ps1`（＝HEAD 抽出的 397 行改前脚本，逐字带 4 处读取）＋ `bin/probe-gui.exe`（GUI 子系统假 wisp）。
命令（`WISP_ENV=test`）：
```
powershell -NoProfile -ExecutionPolicy Bypass -File .../slo-check-before-fix.ps1 `
  -Subset smoke -SecondsPerState 2 -WispExe .../bin/probe-gui.exe -OutDir .../out-before
```
逐字输出（rc=1）：
```
slo-check.ps1: sampling validity precheck (ticket 134 AC#4)
slo-check.ps1: precheck ok - no foreign toolchain/runner process, machine-wide cpu max 14%
slo-check.ps1: sampling state Sleeping for 2s
probe263: sampling started pid=3040 args=[2 3]
D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\263\r1\slo-check-before-fix.ps1 : The variable '$LASTEXITCOD
E' cannot be retrieved because it has not been set.
    + CategoryInfo          : InvalidOperation: (LASTEXITCODE:String) [slo-check-before-fix.ps1], RuntimeException
    + FullyQualifiedErrorId : VariableIsUndefined,slo-check-before-fix.ps1
```
与 CI 那句的差别只有路径与文件标签（消息主体逐字相同）。`out-before/` 里只有 `state-Sleeping.json`、**没有 `slo-report.json`**
（那份 JSON 是死后的孤儿进程后来写的——正是 §② 尾注那个形状）。

**② 改后同一发不再红（交付字节）**

```
powershell ... -File D:\...\scripts\slo-check.ps1 -Subset smoke -SecondsPerState 2 -WispExe bin/probe-gui.exe -OutDir out-tracked-green2
```
逐字（`tracked-green.txt`，rc=0）：两枚钉 ok → `precheck ok ... cpu max 45%` → `state Sleeping exit=0 pass=True` →
`state Warm exit=0 pass=True` → `settle exit=0 pass=True` → `leak exit=1 flipped_to_fail=True` →
`report written to ...\out-tracked-green\slo-report.json (all_pass=True)`；`Measure` 那一发 `real 0m18.748s`
（＝2+2+10+2 秒采样都**等满**了，对照改前的 195 ms 落空）⇒ 门不但活着，而且真的在测。

**③ 真超标仍红（交付字节；不是改判据，是改**样本**）**

同一命令加 `PROBE263_VERDICT=fail`（假 wisp 按 `wisp slo` 的契约同时给出 `pass:false` 的 JSON **和** 非 0 退出码）。
逐字（`retry-console-1.txt`，rc=1）：
```
slo-check.ps1: precheck ok - no foreign toolchain/runner process, machine-wide cpu max 49%
slo-check.ps1: sampling state Sleeping for 1s
slo-check.ps1: state Sleeping exit=1 pass=False
slo-check.ps1: sampling state Warm for 1s
slo-check.ps1: state Warm exit=1 pass=False
slo-check.ps1: settle exit=1 pass=False
slo-check.ps1: leak exit=1 flipped_to_fail=True
slo-check.ps1: report written to ...\out-retry-1\slo-report.json (all_pass=False)
```
⇒ 这一发证明"门红"仍然有路可走到（不是只能死、也不是只能绿）。

**④ 判据临时换成必红形状 ⇒ 门红 ⇒ 还原（票面点名的那一形）**

⚠ 偏离票面 wording 具名：票面说"临时把 `scripts/slo-check.ps1` 的某一档判据改成必红形状再还原"。
这是**共享工作树**，我不在跟踪文件上做瞬时破坏（别人随时可能读/提交它），改成由
`run-nail-mutations.ps1` 从 tracked 字节**逐字节生成副本**（副本断言锚命中才落盘，锚不中就 throw 拒绝出读数），
在副本上改判据。判据那一枚是交付版 `:545` 的 `            $pass = [bool]$reportJson.pass` → `            $pass = $false`。
逐字（`nail-mutations-v2.txt`，`force-state-fail` 一发，rc=1；样本本身是达标形状 `exit=0`）：
```
slo-check.ps1: shape nail ok - this file reads no automatic exit-code variable (ticket 263)
slo-check.ps1: instrument self-check ok - ...
slo-check.ps1: precheck ok - ... cpu max 24%
slo-check.ps1: state Sleeping exit=0 pass=False
slo-check.ps1: state Warm exit=0 pass=False
slo-check.ps1: settle exit=0 pass=True
slo-check.ps1: leak exit=1 flipped_to_fail=True
slo-check.ps1: report written to ...\mutations\out-force-state-fail\slo-report.json (all_pass=False)
=== RUN force-state-fail rc=1 ===
```
⇒ 判据那一行是活路，不是装饰。还原证明：tracked 文件从头到尾没带过这枚变异 ⇒
`git status --porcelain -- scripts` ＝ 空（§⑦ 附同发读数），`git show --stat b92037bd` 里
`scripts/slo-check.ps1 | 240 ++++--`（`git diff --numstat` 同发＝`231  9`），第二枚 commit 只再动 4 行（钉的模式与注释）。

**⑤ 两发被拒采的诚实交代（不藏）**

`tracked-fail.txt`（cpu 75%）与 `tracked-fail2.txt`（cpu 54%）两发都被 precheck 拒采 ⇒ 打 `NO CONCLUSION (machine-contended)`、
`slo-report.json NOT written`、**rc=0**。这正是票面 §现状 点名的"绿＝什么都没测"，本票不许把它做得更糟，也说明
**这台机器上任何采样读数都可能被编队洗掉**（`run-tracked-fail-retry.ps1` 就是为等一个窗口而写的）。

**⑥ 我这把台件尺自己也中过一次同一族雷（记我）**

`run-tracked-fail-retry.ps1` 第一版用"输出里有没有 `precheck ok`"判断是否走到采样——但门是用 `Write-Host` 打的，
`Write-Host` 不进管道 ⇒ 该字段恒 `False`，把 rc=1 的**真红读数读成了"没走到采样"**。改成以
`out-retry-N/slo-report.json` 是否存在为 oracle（结论文件只在做完判定时才写），并在文件里留了这段说明。
⇒ 教训与票面 §② 同形：**判据要用产物，不要用看起来像输出的空串**。

---

## ⑤ 钉子与反形（AC#4：正形读数＋反形读数＋绕得过的形状具名）

**钉 A＝能力钉 `Test-ExitCodeInstrument`（`:214`，形状：起一个已知退出码的外部程序并等它）**
- 正形（交付字节）：`instrument self-check ok - external programs are started, waited for, and their real exit code read (ticket 263)`
- 反形 1 `drop-wait`（把 `:193` 的 `$proc.WaitForExit()` 删掉）逐字（rc=1）：
  `slo-check.ps1: INSTRUMENT BROKEN: exit-code instrument: asked the shell to exit 7 and read 0; the start-wait-read path is broken (ticket 263)`
- 反形 2 `hardcode-zero`（把返回码写死成 0）：同一句红（`read 0`），rc=1。
- 为什么这枚钉对"没等"敏感（实测，不是推）：`repro-probe-slow.ps1` 两发逐字 `NO-WAIT read succeeded exit=0 (hasExited=False)`
  ⇒ 不等时读到的是 **0**，不是异常 ⇒ 与钉期望的 7 不符 ⇒ 必红。**这条正是"先赋再取"那一支会静默放过去的形状。**
- 它的射程限制（具名）：它证的是"退出码通道可用"，用的是 `cmd.exe` 而不是 `wisp.exe`；wisp.exe 起不动那一支由调用点的
  `ok=false ⇒ Fail-InstrumentBroken` 兜（同一形状，路径不同）。

**钉 B＝词面钉 `Test-ScriptShapeNail`（`:233`，模式在 `:252`）**
- 模式＝`(?i)` + `\$`（由 char code 36 拼出，故不匹配自己）+ `\{?(?:[A-Za-z]+:)?LASTEXITCODE`。
- 正形：`shape nail ok - this file reads no automatic exit-code variable (ticket 263)`
- 反形（四发逐字，都 rc=1，`nail-mutations-v2.txt`）：
```
plant-bare         -> First hit text: $LASTEXITCODE
plant-braced       -> First hit text: ${LASTEXITCODE
plant-global       -> First hit text: $global:LASTEXITCODE
plant-global-used  -> First hit text: $global:LASTEXITCODE
```
  ⚠ 前两条在 `cb8dcecc` 之前也红，后两条**是第一次读数不红、我改了模式之后才红**（第一版只认 `$`/`${` 两形）。
  这是我自己造出来的反形逼的改动，不是顺手扩射程：见 `nail-mutations-raw.txt`（旧版：plant-global 全绿）与
  `nail-global-fresh-session.txt`（旧版新会话：同形无名死）。
- **绕得过的形状（具名说绕得过，⛔ 不当牙）**——`nail-mutations-v2.txt` 的 `plant-getvariable` 一发：
  种 `$stale = (Get-Variable -Name 'LASTEXITCODE' -ValueOnly)` ⇒ `shape nail ok`、`report ... (all_pass=True)`、**rc=0 全绿**。
  ⇒ 词面钉认文本，运行期取名的这一族它看不见。**同一族更坏的一形我也量了**（`bypass-pair.txt`）：
  同一会话里先跑 `clean-copy`（`exit 0`）再跑 `plant-global-used` ⇒ 那发**rc=0 全绿**，因为它的"档位退出码"是从
  **上一枚程序的退出码借来的**（新会话里同形则是 `VariableIsUndefined` 无名死，逐字在 `nail-mutations-extra.txt`）。
  ⇒ 结论写死：词面钉挡不住（i）运行期取名／拼串取变量、（ii）把钉和等待一起删掉、（iii）把读取挪去别的文件（模式只扫 `$PSCommandPath` 这一个文件）。
  CI 里每一 step 是一个新进程 ⇒ 在 CI 借不到号，那一形会退回 `$null`；`$null -eq 0` 为假 ⇒ 档位判红（fail-closed 方向）——这一支是读码推的〔未实测〕。
- 外层牙（不在本文件、本票不改，只具名指路）：`scripts/slo-freshness.sh` 的 P3 在"最新上传的 slo-report artifact"上 aging
  ⇒ "门再也不出结论"这一族由那枚钉自己变红，不靠我这两枚自锁钉。

**两枚钉的共同局限（一句实话）**：它们都在**被保护的文件内部**，是自锁；本仓 CI 里没有任何一把尺扫 `scripts/*.ps1`
（§③ 末那条上报）。所以"未赋值即死不再回来"这件事，今天的真牙是：**词面/能力两枚自锁钉（挡顺手写回旧形状）＋
下一次真 run 出结论（AC#5）＋ P3 那枚 aging 钉（挡整道门不再执行）**。

---

## ⑥ 诚实档：把"取不到结论"变成具名失败

交付版现在有三色，语义互不覆盖（`ci.yml:650-654` 的 Upload 步只认 `build/slo/slo-report.json`）：

| 形状 | 屏幕上具名句 | 退出码 | 盘上记录 | 是不是结论 |
|---|---|---|---|---|
| 仪器坏了（起不动／等不到／拿不到码／钉红） | `slo-check.ps1: INSTRUMENT BROKEN: ...`（`:94-129`） | **1** | `build/slo/slo-instrument-broken.json`（`verdict='instrument-broken'`、`d32_evaluated=false`、零 state 文件） | 不是结论，是坏 |
| 有效性拒采（既有，机器争用） | `NO CONCLUSION (machine-contended) ...`（`:437-519`，一字未动） | **0**（Q-36 裁定） | `slo-no-conclusion.json` | 不是结论 |
| 采样成功但超标（既有） | `state X exit=N pass=False` / `report written ... (all_pass=False)` | **1** | `slo-report.json` | 是结论 |

要点逐条：
- **今天这发连 `NO CONCLUSION` 都没走到**＝纯运行时异常。交付版把这类"取不到"全部换成第一行的具名形状，
  退出码仍是 1（和现状同色），但下一任看得见名字、原因（`why` 逐字进句）和一份能留存的记录。
- **没有把现状做坏**：既有"绿＝什么都没测"那一支我原样保留（那是 134 AC#6 的裁定射程，本票不许我改颜色规则），
  新增的一律是 exit 1 方向；`instrument-broken` 记录**故意不叫 report、也不叫 no-conclusion**，
  所以 Upload 步不会把它当报告、P3 也不会被它骗老——仪器坏掉的门不会伪装成"有结论"。
- 记录文件写在 `$OutDir`：`Fail-InstrumentBroken` 自己先 `New-Item -Force`（`:107`），因为坏在"缺件自建"那支时
  脚本主体的建目录（`:318`）还没跑到。

---

## ⑦ 门禁读数（AC#6，2026-10-04 11:44-11:48 +0800 现跑）

**1. `sh scripts/d22scan.sh`** — **rc=0，clean**
```
runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0
d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=228, bans #1-5 cmd/=38,
  ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85,
  ban #8 internal/=502, ban #8 cmd/=99
```
⚠ 与派单转的编排者 10:40 终值 `228/38/85/23/39/85/501/98` 相比，本发是 `... /502/99`：`internal/` ban #8 501→502、`cmd/` 98→99。
归因（不是我这发）：`git log --diff-filter=A` 现跑显示 10:20 `256-r1` 新增 `cmd/wisp/resident_approval_risk_256_windows_test.go`、
11:10 `260-r3` 新增 `cmd/wisp/resident_cancel_key_wording_260r3_windows_test.go`。⇒ 分母随他人提交长大，**不是卫生扩大**；
本腿没动 `internal/`、`cmd/` 任何 Go 文件（§③ 的 `git show --stat`）。
`tools/d22scan` 那两枚 CI 新增红（A586 §4，属票 212 射程）：**本机这发看不见**（本机正控 35 PASS/0 FAIL），我不动它、只具名。

**2. `sh scripts/check-path-length-budget.sh --with-self-test`** — **rc=0，VERDICT GREEN**
```
control 1/3 ok / control 2/3 ok（种一枚超长名被逐字点名，相对 125、名 104、全路径 169 vs budget 165）/ control 3/3 ok
denominator: tracked paths=5484  over-budget=57  covered by roster=57  not in roster=0
longest=180 chars relative（.scratch/wisp/issues/252-...md）
worst full path on the self-hosted runner=224 chars (hat budget 165, wall open interval (206,217])
bands: over the hat=57  ...  in the wall interval=0
```
⚠ 按 A596 结清的单位口径读：**165/180/224 与 (206,217] 是两种单位**（前者帽是全路径预算、后者是相对量开区间），本票不把它们放进同一次比较。
本腿新增路径最长一枚相对 54 字符（`git ls-files .scratch/wisp/probes/263 | awk '{print length}'|sort -rn` 现跑，前 5 名 57/52/52/50/50），
远低于 121 的相对帽 ⇒ 不给 262 那扇门添名册。
⚠ 我这发也**没有**新增 `issues/` 长文件名（票名沿用 263 既有票面）。

**3. emoji 自查（ban #8 仪器实际射程，字符串不豁免）**
```
perl -CSD -ne 'print if /[\x{1F000}-\x{1FAFF}\x{2200}-\x{22FF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}\x{FE0F}]/' scripts/slo-check.ps1
-> 0 命中
```
非 ASCII 只在既有行（`§`、`:438` 的中文引语"都按推荐"），我的新增行全 ASCII；箭头我也用 ASCII `->`，
没用 `→`（仪器不扫 `U+2190`，但没必要留灰色地带），更没用 `✓`/`≤`。
`tools/d22scan` 不扫 `scripts/` ⇒ 这一条是自查、不是门（见 §③ 末的上报）。

**4. 其余两门的适用性**
- `gofumpt`/`go vet`：**本票零 tracked Go 产码改动**（`git show --stat` 两枚 commit 只有 ps1＋本目录台件），故 tracked Go 的名册不会因我变红。
  台件里的 `src/main.go` 有**独立 `go.mod`**（`module probe263`）且不进根模块 `./...`（Go 的包模式跳过以 `.` 开头的目录），
  `gofumpt -l src` 本机现跑＝空（用 `tools/gofumpt.exe`，⚠ CI 用 `@latest`，版本线不同 ⇒〔仅本机可量〕）。
- 票 212 ban #9（同名双写那族）：本机 d22scan clean ⇒ 未扩大。
- **红名集合逐名比对**：本腿没有跑 CI（不 push），所以只能给"本机两门 rc=0＋d22scan 分母逐名已归因"这一半；
  CI 侧红名册的变化归 AC#5 那一次真 run（编排者核）。

**5. 占位标记尺（交件判据＝0）——`placeholder-check.txt` 同发跑（11:53:01 +0800）**
```
证据件本体 fix-and-readings.md = 0 命中（单独复尺：evidence-hits=0）
scripts/slo-check.ps1 = 0 命中
本目录 24 枚 tracked 件合计 = 1 命中，唯一那枚在 msg-skeleton.txt
```
⚠ 那 1 枚命中**不是未填的空位**：它是骨架 commit `82af691d` 的消息正文，句子里在**说明**"骨架发允许占位标记、交件发必须清零"这条规矩
（逐字：〔占位标记〕只出现在骨架发，交件发必为 0 枚）。已提交的消息件不改写（AGENTS §1.4 追加不抹），所以具名留在这里而不是抹掉它。
⚠ 另一枚小雷我也避开了并把形状写下来：**占位词的 pattern 本身不许写进被扫的文档**，否则"尺自我命中"会让命中数变成 1——
这与 §⑤ 词面钉"模式串由 char code 36 拼出、不匹配自己"是同一个形状。

---

## ⑧ 判不动的地方与量不到的格子

1. **AC#5 整格＝没验证，具名停在这里。** 我没有推送权限（⛔ 只 commit 不 push），本机一发不能冒充 CI 读数
   （既有定式：CI 与本机不同口径不可互比；且这台机器本身就是 runner，编队会洗读数）。
   销账条件照票面：一次真 run 里 `slo-full` 走到出结论、`Upload SLO report` 不再 `skipped`。在那之前任何地方都不许写"D32 已被 CI 保护"（本文件也没写）。
2. **`slo-smoke` 那一步是否同因**：票面禁止本票归因/处理。我只在 §① 4 留了一条同形事实读数（70 ms 落空＋逐字红句），没为它改一个字。
3. **改后的"缺件自建"那支（`:299-315`）没有对真 `scripts/build.ps1` 实测**：它会跑一次 cgo 全量构建并动 `build/`
   （派单禁：⛔ 别碰 `build/` 里别人的产物、别在编队飞行时抢 CPU）。已实测的是**机制同形**：
   `child-exit.ps1`／`child-native.ps1` 两发证明子脚本的 `exit N` 与子进程退出码都能被父侧读到，
   而新形状用 `Get-SelfHostExecutable`＋`Invoke-ExternalProgram` 起 `-File scripts/build.ps1 -Env dev`
   （参数含空格路径 ⇒ 走 `Get-NativeArgumentLine` 加引，实测在含空格的本仓路径上跑通）。真 build.ps1 那一发的 rc〔未实测〕。
4. **CI 上 wisp.exe 死后的存活时长**：台件同形量到 1 枚孤儿（§②），CI 那一发的进程表当时没人取 ⇒〔未实测〕。
   ⇒ "改前那形会让下一次采样被自己的孤儿判成 machine-contended"这句只是方向成立、无 CI 证据。
5. **Get-Variable 反形在 CI 的颜色**：本机量到"钉不响＋门可全绿"（借号形，同会话）与"新会话无名死"两形；
   CI 每步新进程 ⇒ 我推它是 `$null` ⇒ 判红方向，但**没在 CI 上量过**〔未实测〕。
6. **`tools/d22scan` 对 `scripts/*.ps1` 的盲区**：只读 `main.go` 的 scope 名册判定，没种样本做正控（§③ 末）；
   要修那把尺属票 212 射程，本腿不动。
7. **本机争用**：两发被 precheck 拒采（§④ ⑤）证明我这台的读数可被编队洗掉；所以本文件每条采样读数都带当时的 cpu 数。
8. **"每一档必须出结论"能不能被参数绕过**：`ValidateSet('smoke','full')`（交付版 `:74`）＋两步命令行都写死（`ci.yml:585`/`:648`），
   我读码判定"没有第三档、没有跳过档"；**没有**跑"往 ci.yml 塞一个 `-Subset` 值看门响不响"的正控（那是 262/134 的写面）〔未实测〕。
9. 我不判"票 244 该不该切 GUI"：那是已定案的产品决策（双击不带黑窗），本票只让仪器跟随它，不动那枚旗标一字。

---

## ⑨ 交件自问：改完之后，那两道门有没有哪一道从此永远不再执行

**门一＝D32 那两个数字（Sleeping CPU <=0.5%、private RSS <=25MB，由 `wisp slo` 求值）。**
- 有没有被我做成"永不执行"？没有新增任何跳过：唯一能提前结束脚本的新分支是 `Fail`／`Fail-InstrumentBroken`／
  两枚钉，全部 `exit 1`（响），不是 `exit 0`（绿）。既有那支 `exit 0` 的机器争用拒采我一字未动（那是 134 AC#6 的裁定）。
- 但它今天**仍未被证明求值过**：本腿没有 CI 读数（AC#5 未验证）。我能说的只有"本机用假 wisp 走完了全路径并等到真实退出码、
  写出了 `all_pass`"，那证的是**路的形状**，不是那两个数字的值。⇒ 我不写"D32 已被 CI 保护"。

**门二＝退出码语义（0＝通过或无结论，1＝有数字且没过，或仪器坏了）。**
- 三档语义没被稀释：通过＝0、超标＝1、仪器坏＝1（新增具名）、无结论＝0（既有）。
  新增的 `INSTRUMENT BROKEN` **不占用** 0 那一侧，所以"绿"从此只有两种来源：真通过、或既有的机器争用拒采（后者由 P3 aging 钉兜）。
- 有没有从此不再执行的风险点？有一个我如实留下：两枚钉在被保护文件内部，删掉它们不需要动别的文件；
  外部只剩 P3（跑在 `slo-fresh.yml`，另一步）能间接发现"门再也不出报告"。⇒ 这条不是牙，是**记录**。

**最后一句诚实话**：本票把"测不到"改成了"要么测到、要么大声说坏"，但**这两个 D32 数字有没有真的达标，今天仍然没人知道**——
只有 AC#5 那一次真 run 能回答，而那不归我（不 push、不冒充 CI 读数）。
