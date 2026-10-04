# 票 263 验收表（263-v1，非实现者）— 攻 `263-r1` 那五枚 commit

被验对象：`82af691d`（骨架 11:20）→ `b92037bd`（落地 11:32）→ `cb8dcecc`（词面钉改 scope 拼法 11:43）→ `befb779c`（交件 11:55）→ `c9411a1e`（收尾 12:02）。
本腿＝`263-v1`（非实现者验收腿）。⛔ 未改任何产码、⛔ 未碰票 263 的任何 AC 框、⛔ 未 push。
终尺钟点 `2026-10-04 13:05:02+0800`，起手 HEAD＝`67f04685`（＝本腿骨架枚），未推枚数 `79`（`git rev-list --count origin/dev..HEAD` 现跑）。
`r1` 证据件（`.scratch/wisp/probes/263/r1/fix-and-readings.md`，473 行）里每一句都按待验断言处理：**下面每一格给的都是本腿自己取的数**，凡引用它的数都写明〔仅它的读数，本腿未复跑〕。
本腿台件全在 `.scratch/wisp/probes/263/v1/**`（假件源码 `src/`、两枚 exe `bin/`、改前字节 `slo-check-prefix-head.ps1`、突变副本与日志 `mutations/`、两门读数 `gate-*.txt`）。

七格总判（逐格正文见同名节号）：

| 格 | 判语 |
|---|---|
| ① AC#1 根因链 | **成立**（三验齐；②为同旗标等价件读数，具名〔非今天 CI 那枚产物的直接读数〕） |
| ② AC#2 四条偷懒 | **成立**（四条逐条有尺；另抓到它一枚枚数错） |
| ③ AC#3 正控三发 | **成立**（三发本腿全部重跑，未引用它的读数文件） |
| ④ 词面钉有没有牙 | **带条件成立**（两枚正控有牙；三形绕过本腿实测全过） |
| ⑤ AC#4 是否同义反复 | **带条件成立**（不是"扫自己的常量"，但是同文件自锁；另抓到一枚误报方向） |
| ⑥ AC#6 卫生两门 | **成立**（两门本腿现跑 rc=0；+8 归因独立验成立；502/99→503/100 归因给 260-r4） |
| ⑦ AC#5 CI 销账 | **不成立＝没验证**（本腿不 push；⛔ 任何"D32 已被 CI 保护"都不写） |

---

## ① 现量复认与名册独立拉取

**五枚 commit 逐枚名册**（`git show --name-only --format= <c>`，probes 之外的全部文件）：

```
82af691d total=2  probes263=2  非 probes：（零）
b92037bd total=12 probes263=11 非 probes：scripts/slo-check.ps1
cb8dcecc total=12 probes263=11 非 probes：scripts/slo-check.ps1
befb779c total=9  probes263=9  非 probes：（零）
c9411a1e total=4  probes263=4  非 probes：（零）
```
⇒ **"probes 之外只碰了 `scripts/slo-check.ps1` 一枚"成立**，且是**我拉的**名册不是它的。禁区逐名（同一段 diff 名册上 grep）：
`docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**`／三枚冻结件／`.github/workflows/ci.yml` —— **零命中**（grep rc=1＝无匹配）。
`internal/observe/thresholds.go` 另加一把 md5 尺：工作树＝`b8fabe32710b3ccc725b930237b3ffe9`＝HEAD＝`git cat-file blob 82af691d^`＝`git cat-file blob c9411a1e` 四个锚**逐字节同** ⇒ 这枚 D32 契约面在本票全程没动过。
`git status --porcelain -- scripts .github`＝**空**（13:0x 现跑）⇒ 写面此刻真空，没留疤。

**编排者要我复认的三件事，我自己现跑**：
```
wc -l scripts/slo-check.ps1                      -> 625（自述"397→625"成立）
Invoke-ExternalProgram（function 定义行）        -> :143
INSTRUMENT BROKEN（具名失败的屏幕上第一句）      -> :104（Fail-InstrumentBroken 内）
INSTRUMENT BROKEN（两枚钉调用它的那两行）        -> 词面钉 :256、能力钉 :225/:228
```
⇒ 派单说的"`:256`"是**词面钉那一发的调用点**（`Fail-InstrumentBroken ('shape nail: ...')`），不是具名句的真身；两种数法都成立，我把两处分开写清，免得下一任以为只有一枚具名出口。派单里的三件事**全部复认**。

**改前字节我自己抽**（`git cat-file blob b92037bd^:scripts/slo-check.ps1` → `slo-check-prefix-head.ps1`）：397 行，`LASTEXITCODE` 命中 `:106`/`:326`/`:345`/`:366`，`Set-StrictMode -Version 2.0` 在 `:65` ⇒ 与票面 §现量 4 逐字同、与 r1 §① 逐字同（两把尺都独立复跑到）。
**四个调用运算符 `&` 的位置**（同一份改前字节）：`:105`（缺件自建那支）、`:325`（逐档采样）、`:344`（settle）、`:365`（leak）⇒ 与四个读取点一一成对；改后文件里 `^\s*& ` **零命中**（尺：`grep -n "^\s*& "` 两版对拉）。

---

## ② ★ AC#1 根因链的三验

**① `scripts/build.ps1:115` 的 `-H=windowsgui` 真在**
逐字 `    "-H=windowsgui"`（在 `$ldflags` 数组 :103-116 内，:103 注释写"SPEC-11 §2.2: no-arg launch = GUI"）。
`git log -S'-H=windowsgui' -- scripts/build.ps1` 现跑只回一枚 commit＝`cc6eaa65 2026-10-02 19:39 244-r1: build.ps1 追加 -H=windowsgui，双击不再带黑控制台` ⇒ **肇因枚号与票号都对得上，且没有第二枚动过它**。

**② 产物 PE 头真子系统＝GUI(2)**（⛔ 不读注释不读文件名）
本仓 `build/` 目录此刻**不存在真产物**（`ls build/` 空），而本腿⛔ 不许跑 cgo 全量构建、⛔ 不许动别人的产物 ⇒ **今天 CI 那枚 `wisp.exe` 的直接 PE 读数我量不到，具名说没量到**。
我改用**同旗标同工具链等价件**：同一份 Go 源（本腿自己写的 `src/main.go`）两发构建——
```
go build -o bin/probe263v1-cui.exe .
go build -ldflags "-H=windowsgui" -o bin/probe263v1-gui.exe .
objdump -p bin/probe263v1-cui.exe | grep -E "^\s*Subsystem"  -> Subsystem 00000003 (Windows CUI)
objdump -p bin/probe263v1-gui.exe | grep -E "^\s*Subsystem"  -> Subsystem 00000002 (Windows GUI)
```
⇒ **`-H=windowsgui` ⇒ Subsystem=2** 这一环在本机、本工具链（go1.27.1 windows/amd64）上成立；
CI 那一步确实跑的就是这枚旗标所在脚本：`.github/workflows/ci.yml:645`（`... -File scripts/build.ps1 -Env dev`，步名 `Build wisp.exe (deps cached on the runner)`）与 `:648`（`... slo-check.ps1 -Subset full -SecondsPerState 6`）——两处**我自己 grep 现读**，⛔ 没碰 `ci.yml` 一字。
顺带复认派单一处尺错：`ci.yml:585` 逐字是 `-Subset smoke -SecondsPerState 4`（**4 秒**）⇒ r1 §① 对我派单的更正**成立**，我独立读到了同一行。

**③ "PowerShell 的 `&` 对 GUI 子系统进程不等待"——我亲手复现过**
最小件 `repro-wait.ps1`＝把改前 `:324/:325/:326` 三条语句逐字搬出（`Write-Host` → `& $Exe ...` → `$code = $LASTEXITCODE`），同一宿主（`powershell 5.1.26100.8457`，`-NoProfile -File`，`$ErrorActionPreference='Stop'`＋`Set-StrictMode -Version 2.0`）两发差分：
```
CUI：call returned after 2027ms -> state exit=0 -> rc=0
GUI：call returned after   14ms -> The variable '$LASTEXITCODE' cannot be retrieved because it has not been set.
        + FullyQualifiedErrorId : VariableIsUndefined
        子进程 probe263v1: started / exiting 两行出现在异常之【后】
```
⇒ "不等"不是推论：**14ms vs 2027ms**，且子进程的输出在取值失败之后才落到控制台（进程还在跑，调用方已经走了）。
再用**改前门本体**（397 行字节件）× 两枚假件各一发（`-Subset smoke -SecondsPerState 1`，`WISP_ENV` 由门自己设 `test`，⛔ 未跑 `-Subset full`）：
```
门 × GUI 件：precheck ok ... cpu max 26% -> sampling state Sleeping for 1s -> 红句（主体逐字同 CI）-> rc=1
             盘上只有 state-Sleeping.json，【没有 slo-report.json】（那份 JSON 是死后孤儿进程补写的）
门 × CUI 件：同字节、同假件、只换子系统 -> 两档 + settle + leak 全跑完 -> all_pass=True -> rc=0
```
⇒ **差分把肇因钉在"子系统"这一枚变量上**，不是"忘了初始化"、不是争用（两发 cpu 26%/未拒）、不是假件行为（同一枚 Go 源）。票面 §现量 3 那句"肇事的读取点不是那三行"复认成立，真点＝改前 `:326`。
**为什么这一路径今天才走到**（票面 AC#1 那半格）：我这一发 `-WispExe` 指了现成 exe ⇒ 改前 `:105`（缺件自建那支的 `&`）被 `Test-Path` 跳过；CI 同形——`:645` 已先建好 `build\wisp.exe`，门进去就跳过 `:105`。于是 precheck 通过后**本会话第一条外部程序就是 `:325`**，取值在 `:326` ⇒ 前面 `:102-:311` 全是 cmdlet，自动变量从未被写过。这条链条我在改前字节上逐段读过，并由上面那发的顺序（precheck ok → 第一档即红）复现。

**★ 一枚我要写清、r1 没写的机制细节**：最小件那一发 **rc=0**（异常之后脚本继续跑下一条语句， powershell 以 0 收）。门本体那一发是 rc=1，因为失败之后没有别的语句把它掰回 0。⇒ "**取未赋值变量＝当场死**"在**门本体的形状**下成立（rc=1、无 report，与 CI 同），但机制上它是**语句级作废**而不是脚本级终止；这意味着"同一形状落在别的语句位置上"存在 **rc=0 而什么都没测**的可能。这一条我给编排者当**判据射程信息**（它正好压在 AC#2④ 与"绿＝什么都没测"那枚既有的坑上），⛔ 不构成对本票修法的否定——改后的形状压根不再走这条语句。

**副作用 vs 缺陷（分开写，这条决定票 244 与票 263 的关系）**
- **票 244 裁的"双击不带黑控制台"这一裁成立**：它要的 GUI 子系统出厂件，我在同旗标等价件上量到 Subsystem=2；我没跑真机双击，那格归票 244 自己的 AC，本票不判。
- 票 244 引入的是**一次跨票副作用**，不是缺陷：它改了"出厂件子系统"这枚全局变量，而**下游消费者名册里少了一列"谁在等这个进程"**（本仓的求值器 `scripts/slo-check.ps1` 是唯一的 `&` 消费者）。副作用的**受害面**是求值器，⛔ 不是票 244 的实现错了。
- **票 263 治的是求值器自身的 latent 形状**：`& + 自动变量`这一族在 CUI 时代也只是"恰好每档都等到"才安全——改前 `:106` 那枚潜伏点（r1 说它是同形状第二枚，实测它今天不炸）就是证据。
- 建议（**登记，不属本票射程**）：任何"改出厂件子系统／启动形态"的票，落点普查里加一格"谁在等这个进程"（`grep -n "^\s*& "` 全仓脚本），否则同类副作用还会第三次出现。归口＝编排者是否立新票；本腿不动任何产码。

---

## ③ AC#2 四条 ⛔ 偷懒形状：每条一把尺＋判语

**第 1 条「没用 try/catch 吞」——尺：两版 catch 逐枚对拉（`grep -n -A3 catch`）**
```
改前 4 枚：:141（枚举进程失败 -> Fail，exit 1）、:217（perf 类读不到 -> cpuMax=-1; break）、
           :333/:357（JSON 解不开 -> WARNING，pass 保持 $false）
改后 7 枚：上述四枚【逐字保留】在 :349/:425/:546/:575；新增三枚：
           :122（INSTRUMENT BROKEN 的记录写不出去 -> 只补一句报告，随后 :128 仍 exit 1）
           :183（起不动 -> return ok=$false + why）
           :195（等不到/取不到码 -> return ok=$false + why）
```
新增三枚里没有一支"之后当作测过了继续跑"：:183/:195 的 `ok=$false` 在**四个调用点一律**走 `Fail-InstrumentBroken ⇒ exit 1`（`:311`、`:537`、`:561`、`:591`），:122 之后落 `:128 exit 1`。⇒ **成立**。
附一把邻尺：`SilentlyContinue` 枚数 3→5，两枚新增在 `Get-SelfHostExecutable`（`:205`、`:208`），其失败出口是 `:305 Fail-InstrumentBroken` ⇒ 不是把异常染色成绿。

**第 2 条「没放宽 `Set-StrictMode`」——尺：两版逐字对拉**
`Set-StrictMode -Version 2.0`：改前 `:65` → 改后 `:82`，**字符串逐字相同**（只是被上面新增的注释块整体推下 17 行）；`$ErrorActionPreference = 'Stop'`：`:64` → `:81` 逐字同。全文件 grep `-Version 1`／`Off`／`-Version 3`：**零命中**。⇒ **成立**（这条我把整枚文件扫了，不是只看那一行）。

**第 3 条「阈值/golden/`thresholds.go` 一字节未动」——尺：名册＋四锚 md5**
见 §①：五枚 commit 里 probes 之外只有 `scripts/slo-check.ps1`；`thresholds.go` 在工作树/HEAD/`82af691d^`/`c9411a1e` 四处 md5 全同。
门里出现的 `Sleeping CPU <=0.5%`、`RSS <=25MB` 我逐处看了落点：`:29`/`:343`/`:449`/`:464-465` 全是**注释与 step-summary 文案里的引用文字**，求值仍在 `wisp slo` 一侧 ⇒ 判定线没被搬进脚本、也没被搬走。⇒ **成立**。
⛔ 我没有对全仓 golden 逐枚取 md5（不在本票写面、也不在派单要求内），这一格的射程＝"本腿五枚名册里零 golden"。

**第 4 条「没把任何'必须出结论'的一步改成跳过」——尺：三把**
- (a) `exit` 名册对拉：改前真退出 4 枚（`:74` exit1 / `:311` exit0 / `:396` exit1 / `:397` exit0）→ 改后 5 枚（`:91` exit1 / **`:128` exit1（新增）** / `:519` exit0 / `:624` exit1 / `:625` exit0）。⇒ **`exit 0` 一枚没多**，新增那一枚是 `exit 1`；既有的 `NO CONCLUSION (machine-contended)` 那两行（`:310-311` → `:518-519`）**逐字未动**（它自述一致，我核了字）。
- (b) **四个读取点逐处对读**（判"有没有哪一档从此不再向操作系统求值"）：
```
改前 :105 & build.ps1 + :106 读自动变量   -> 改后 :307 Invoke-ExternalProgram + :311 ok 检查 + :314 code≠0 -> Fail
改前 :325 & wisp slo -state + :326       -> 改后 :533 + :536 ok 检查 + :539 $code=$run.code
改前 :344 & wisp slo -settle + :345      -> 改后 :557 + :560 + :563
改前 :365 & wisp slo -leak + :366        -> 改后 :583 + :587 + :593
```
四档**每一档仍真起进程、仍 `WaitForExit()`（`:193`）、仍取 OS 给的 `ExitCode`（`:194`）**；起不到/等不到时不是"跳过"而是 `Fail-InstrumentBroken` ⇒ `exit 1`＋盘上记录。`Invoke-ExternalProgram` 调用点现跑＝5 枚（`:223` 能力钉 + 上述四档），`^\s*& ` 现跑＝0 枚。⇒ **没有哪一档不再向操作系统求值**。
- (c) 参数面：`[ValidateSet('smoke','full')]` 仍在 `:74`，`-Subset` 无第三值、无 `-Skip`/`-Force`/`continue-on-error`（grep 名册零命中）。⚠ 我**没有**跑"往 `ci.yml` 塞一个 `-Subset` 值看门响不响"的正控（那是 262/134 的写面），这半格〔未实测〕，与 r1 §⑧ 8 同格同归口。
⇒ **成立**。

**抓到 r1 一枚枚数错（不影响结论，具名记）**：它 §③ 写"交付文件里 `$LASTEXITCODE` 字面出现 **1 次**（`:252` 那枚模式串本身）"。我现跑 `grep -n "LASTEXITCODE" scripts/slo-check.ps1`＝**2 行**：`:241`（注释里的 `Variable: LASTEXITCODE`）与 `:252`（模式串）。两行都**不带 sigil**（尺：`grep -c '\$LASTEXITCODE\|\${LASTEXITCODE'`＝**0**），所以钉不匹配自己这句结论仍然对；**枚数是错的**。⇒ 若翻勾，凭据文字里这句应改成"字面 2 行、带 sigil 0 处"。

---

## ④ AC#3 正控三发：本腿自己重跑（⛔ 没引用它的任何一个读数文件）

三发全部用**本腿自己写的 Go 假件**（`src/main.go`＋`bin/probe263v1-{cui,gui}.exe`，PE 头自证子系统）与**本腿自己从 git 抽的字节**：

```
① 改前必红（门本体 × GUI 件）rc=1，逐字：
slo-check.ps1: sampling validity precheck (ticket 134 AC#4)
slo-check.ps1: precheck ok - no foreign toolchain/runner process, machine-wide cpu max 26%
slo-check.ps1: sampling state Sleeping for 1s
probe263v1: started pid=31780 args=[slo -state Sleeping -seconds 1 -interval-ms 250 -out ...out-pre-gui\state-Sleeping.json]
D:\...\263\v1\slo-check-prefix-head.ps1 : The variable '$LASTEXITCO
DE' cannot be retrieved because it has not been set.
    + CategoryInfo          : InvalidOperation: (LASTEXITCODE:String) [slo-check-prefix-head.ps1], RuntimeException
    + FullyQualifiedErrorId : VariableIsUndefined,slo-check-prefix-head.ps1
probe263v1: exiting            <- 死后子进程还在跑（孤儿，写在红句之后）
盘上：out-pre-gui/ 只有 state-Sleeping.json，【无 slo-report.json】
```
与 CI 那句只差路径与文件标签，消息主体逐字同。〔它给的 195ms/70ms 我不引用；我的对应数是 `14ms`（最小件）与"第一档即红"。〕

```
② 改后不再红（tracked 字节 × GUI 件，-Subset smoke -SecondsPerState 1）rc=0：
shape nail ok -> instrument self-check ok -> precheck ok cpu max 24%
state Sleeping exit=0 pass=True / state Warm exit=0 pass=True / settle exit=0 pass=True
leak exit=1 flipped_to_fail=True / report written ... (all_pass=True)
盘上：state-Sleeping.json state-Warm.json settle.json leak.json slo-report.json 五枚齐
```
两枚假件行的 `started` 与 `exiting` **夹在门的两个档位行之间**（而不是像改前那样掉在红句后面）＝等待真的发生了。

```
③ 真超标仍红（同字节、同命令，只把样本换成 fail 形状：V1_VERDICT=fail ⇒ pass:false + exit 1）rc=1：
shape nail ok / instrument self-check ok / precheck ok cpu max 18%
state Sleeping exit=1 pass=False / state Warm exit=1 pass=False / settle exit=1 pass=False
leak exit=1 flipped_to_fail=True / report written ... (all_pass=False)
```
⇒ **判红这一条路是活的**：既不是只能死，也不是只能绿。

```
④ 判据那一行的活路（票面点名的"临时把判据改成必红形状"）——见 §⑤ 的 force-state-fail 一发：rc=1，
   样本本身达标（exit=0）却 pass=False、all_pass=False。⚠ 形状偏离票面 wording 具名：
   票面说改 tracked 文件再还原，本腿同样选择【在副本上改】——共享工作树里不做瞬时破坏（与 r1 §④ 的理由一致）。
```

**还原证明**：本腿台件每发跑前跑后各取一次 tracked md5，终值 `8a288306e95fc98c31bdb8386ea6a675`＝起始值，`tracked UNCHANGED = True`；`git status --porcelain -- scripts`＝空。
**争用交代**：本腿 `smoke/1s` 共 12 发，**零发被 precheck 拒采**（每发 cpu max 读数在档：15/18/22/24/26%）。⛔ 我不拿 r1 那两发 refusal 当自己的凭据，也不据它们否定它们。

---

## ⑤ ★ 词面钉有没有牙：本腿实测八发（副本突变，tracked 字节零破坏）

编排器／生成方式：`mutations/base.ps1`＝`git cat-file blob HEAD:scripts/slo-check.ps1`（md5 与跟踪文件同），
`run_mutations_v1.py` 每枚突变断言锚**恰好命中 1 次**否则 raise；编排一开始用 PowerShell，**连撞四把我自己的尺**（见 §⑨ 记我），最终用 Python 起-等-取，退码由 `subprocess` 给（⛔ 不是 `cmd | head; echo $?` 那种取到 `head` 的假退码）。

| 突变 | 含义 | rc | 门打了什么 |
|---|---|---|---|
| `copy-noop` | 副本机制自身对照（零改动） | **0** | 与 tracked 字节同色（全路径绿）⇒ 副本不失真 |
| `drop-wait` | 删 `$proc.WaitForExit()`（`:193`） | **1** | `INSTRUMENT BROKEN: exit-code instrument: asked the shell to exit 7 and read 0` |
| `plant-bare` | 种裸读 `$code = $LASTEXITCODE` | **1** | `INSTRUMENT BROKEN: shape nail: 1 direct read(s) ... First hit text: $LASTEXITCODE` |
| `plant-getvariable` | 改用 `&` ＋ `Get-Variable -ValueOnly` 取名 | 1 | **两枚钉都 ok**；红的是 `Get-Variable : Cannot find a variable with the name 'LASTEXITCODE'` |
| `delete-nails` | 只删两枚钉的调用（`:291-:292`） | **0** | 全绿，屏上连 `shape nail ok` 两行都没有 |
| `read-in-other-file` | 读取挪到 dot-source 的 helper 文件 | 1 | **两枚钉都 ok**；**CI 那句原形逐字回来**（`VariableIsUndefined`） |
| `force-state-fail` | 判据那行改成 `$pass = $false` | 1 | `state X exit=0 pass=False`×2、`all_pass=False` ⇒ 判据是活路 |
| `sample-fail` | 交付字节零改动，只换 fail 形状样本 | 1 | 与 §④③ 同形（互证） |
| `delete-nails-and-wait` | 复合最坏形：删两枚钉**并**删等待 | 1 | 打印 `state Sleeping exit=0 pass=False`、`settle exit=0 pass=False`、`leak exit=0 flipped_to_fail=False` ⇒ `FATAL: forced 100MB leak did NOT flip` |
| `plant-in-comment` | 只在注释里写这三个词（零读取） | **1** | `INSTRUMENT BROKEN: shape nail: 1 direct read(s) ... First hit text: $LASTEXITCODE` |

**逐格读数原文**在 `mutations/log-*.txt`（10 枚，每枚一份）。

**判（第 4 格）**：**带条件成立**。
- 有牙的两形我亲眼看它响：`drop-wait`、`plant-bare` 都 rc=1 且落 `slo-instrument-broken.json`（不是只打印）。
- 它自陈的三种绕过，本腿**全部实测通过（即：三形都绕得过）**，超出派单"至少试两种"的下界：运行期取名（`plant-getvariable`）、删钉（`delete-nails`）、别的文件（`read-in-other-file`）。
- ⚠ 反形不是"全绿"，所以这枚判据**不是不敏感的凭据**：10 形里 8 红 2 绿（`copy-noop`、`delete-nails`）。绿的这一枚恰恰是自锁的本体（钉可被就地删除且本机零外牙）——这条不否定它的价值，但**必须写清射程**，否则下一任会把"两枚钉 ok"读成"这台机器上等得到退出码"。
- **本腿给它补两枚它没写的读数**：
  (i) `delete-nails-and-wait` 那一形里，门**伪造了 `exit=0`**（没等就读到 0，与 r1 §③"不等直读 ExitCode 得 0"的实测一致）——最后兜住它的是 **leak 自检要求 `exit==1`** 这枚**第三牙**（`:594`/`:596-598`）。这枚牙不在它 §⑤ 的名册里，是**本腿量到的**；它的存在也说明残余风险：哪天 leak 的形状不再要求 1，这一族就只剩可删的自锁钉。
  (ii) `plant-in-comment` 是**误报方向**：注释里出现这三个词就能把整道 D32 门打死（rc=1）。方向是响不是骗，但它使"在本文件写一句解释性注释"变成有代价的动作——具名入账给编排者（属 AC#4 的形状选择，不属违规，⛔ 本腿不改它）。

---

## ⑥ AC#4 那枚"结构/词面钉"是不是同义反复：买到了什么、没买到什么

**不是"扫自己的字符串常量"那一型同义反复**，尺与理由：模式串由 `[char]36` 拼（`:251`），文件里两处 `LASTEXITCODE` 字面都不带 sigil（`:241`、`:252`；带 sigil 命中＝**0**），而钉扫的是 `$PSCommandPath` 的**整份文本**（`Get-Content -Raw`＋`[regex]::Matches`，`:253-254`）。所以它抓的不是"我的常量"，是"这个文件的字节里任何 `$`＋自动变量名"的形状——`plant-in-comment` 一发（纯注释、零读取）也被它抓住，证明它的判据**比语义更宽**（宽到误报），而不是**自我循环**。

**买到的（本腿实测）**：① 下一人在**本文件**写回 `&`＋自动变量（含 `${...}`、`$global:`、`$script:` 各 scope 拼法）⇒ 当场具名红；② 有人改坏 `Invoke-ExternalProgram` 内部的等待或取码 ⇒ 能力钉当场具名红（`drop-wait`）。

**没买到的（本腿实测，且 r1 自己逐条具名承认过，我没抓到它夸大）**：① 运行期取名（`Get-Variable -ValueOnly`）——钉 ok、能力钉 ok，最终红，但**红得无名**（PowerShell 自己的 `ObjectNotFound`），"诚实档"在这一形失守；② 把两枚钉的调用删掉 ⇒ **rc=0 全绿**，本机没有任何外部尺会发现（`tools/d22scan` 的 scope 名册不含 `scripts/*.ps1`，这一条 r1 §③ 末上报过，**本腿复跑 d22scan 时也无法反驳**——它的 scope 行只列 internal/cmd/frontend/design/tools，见 `gate-d22scan.txt`）；③ 读取挪去别的文件 ⇒ **原形红句逐字回来**，两枚钉全打 ok。

**结论（第 5 格）**：**带条件成立**。票面 AC#4 的措辞是"保证'未赋值即死'不再回来"＋"绕得过就具名说绕得过，别当牙"。前半句这枚钉**达不到**（三形实测都能让它视而不见），后半句 r1 **做到了**（§⑤ 逐条写"绕得过"，还给了反形读数）。⛔ 我不同意把 AC#4 读成"已保证不再回来"；我同意按"本文件词面/能力两枚自锁钉，射程具名，外牙只有 `slo-freshness.sh` 的 P3 aging"这个口径成立。
（它 §⑤ 末与 §⑨ 已经把 P3 指为外牙，并明说"这条不是牙，是记录"——**这句是它自己收缩射程的话，我没抓到它把话说满**。）

---

## ⑦ AC#6 卫生两门逐名（本腿现跑，13:0x）

```
sh scripts/d22scan.sh                                rc=0
  runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0
  d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=228, bans #1-5 cmd/=38,
    ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85,
    ban #8 internal/=503, ban #8 cmd/=100
```
⚠ 与我派单转的 12:0x 终值 `.../502/99` 相比 **+1/+1**。归因（**我独立取的，不是顺着它写的**）：`git log --diff-filter=A --name-only c9411a1e..HEAD -- cmd internal` 现跑只回两枚——`cmd/wisp/resident_cancel_key_label_260r4_windows_test.go`、`internal/agent/approval/cancel_key_label_260r4_test.go`＝**260-r4 的台件**，与 263 无关。263-r1 名册里**零 tracked Go 文件**（§①）⇒ 分母随他人提交长大，卫生射程没扩大。
`tools/d22scan` 那两枚"CI 红／本机绿"（票面 AC#6 点名不许顺手修）：本机这发 **PASS=35 FAIL=0** 看不见它们，本腿也没动 `tools/d22scan` 一字（名册）。
⛔ 我没有跑整包 `go test ./...`：读了 `scripts/d22scan.sh:50-54` 才知道它的正控只走 `runtests.sh -C tools/d22scan ./...`（**仅 d22scan 包**），这才敢跑——这条具名写下，免得下一任以为本腿违了禁。

```
sh scripts/check-path-length-budget.sh --with-self-test   rc=0
  control 1/3 ok / control 2/3 ok（种一枚超长名被逐字点名：relative 125、name 104、
    full path on the self-hosted runner 169 chars, budget 165）/ control 3/3 ok
  positive control PASSED
  denominator: tracked paths=5509  over-budget=57  covered by roster=57  not in roster=0
  longest=180 chars relative（.scratch/wisp/issues/252-...md，仍是他人的票面）
  worst full path on the self-hosted runner=224 chars (hat budget 165, wall open interval (206,217])
  bands: over the hat=57 ... in the wall interval=0
  VERDICT GREEN
```
**单位口径**（守 `A596`）：`budget 165` 与 `169` 是**同口径**（都是全路径字符数），所以控制 2/3 那句比较是合法的；而 `wall open interval (206,217]` 是**相对量**，与 `224`（全路径）**不同口径**——本腿**没有**把这两个数放进一次比较，也没据它们互相推翻。

**那枚 +8 归因，我独立验**（派单点名这条）：
```
git ls-tree -r cb8dcecc --name-only | wc -l = 5484
git ls-tree -r befb779c --name-only | wc -l = 5492      (+8)
git show --name-status befb779c 的 A 行 = 8 枚，逐名全部在 .scratch/wisp/probes/263/r1/
  （msg-evidence.txt / placeholder-check.txt / retry-console-1.txt / retry-wrapper.log /
    retry-wrapper2.log / run-tracked-fail-retry.ps1 / tracked-fail-retry.txt / tracked-fail2.txt）
git ls-tree -r c9411a1e --name-only | wc -l = 5495      （收尾发再 +3，也是自己的台件）
```
⇒ **它那句归因成立**（+8＝本腿自己新增台件，不是他人提交）；我这一发的 `5509` 与它的 `5492` 差 17，来源是 `c9411a1e..HEAD` 之间 260-r4 的 12 枚台件＋两枚 Go 测试件＋本腿的 `verdict.md`（`git diff --name-status c9411a1e..HEAD` 现跑），`over-budget/roster/not-in-roster` 三格逐字未变 ⇒ **谁都没给长名册添新名字**。
`gofumpt`/`go vet` 两门：名册已证"五枚零 tracked Go 产码"（§①）⇒ 这两门**不可能因本票变红**，本腿没有另行整包跑（与上面的禁令同）。票 212 ban #9：本机 d22scan clean ⇒ 未扩大。

---

## ⑧ AC#5 未验证格（不许被任何人判成"已验证"）

**本腿没有真 CI run ⇒ 这一格判"没验证"，停在未勾。** 凭据面：未推枚数 `79`（13:05 现跑，`git rev-list --count origin/dev..HEAD`），推送在编排者手里；`slo-full` 出结论＋`Upload SLO report` 不再 skipped 这两件事**今天盘上不存在**。
本腿本机 12 发全部是 **`-Subset smoke -SecondsPerState 1` × 假件**——它证的是**路的形状**（起-等-取、四档都到 OS 取码、三色谱得对），⛔ **不证 D32 那两个数字的值**，也不得被引用成"CI 读数"（既有定式：CI 与本机不同口径不可互比）。
本件通篇没写"D32 已被 CI 保护"，也没有"票 263 已销账"。销账条件照票面：一次真 run 里 `slo-full` 走到出结论。
另：`slo-smoke` 那一步是否同因——票面禁本票归因，本腿也不归因；我只留一句"同一枚文件、同一处取值点，改后那一步的读取点已一并换成 `Invoke-ExternalProgram`"，**它下一次什么颜色要等真 run**〔未实测〕。

---

## ⑨ 判不动的地方＋票 263 现在够格翻哪几枚框

**判不动／量不到（具名＋归口）**：
1. **今天 CI 那枚 `wisp.exe` 的直接 PE 头**——量不到：`build/` 无产物，本腿禁跑 cgo 全量构建、禁碰他人产物。已用同旗标等价件补，归口＝若要把这格升成"产物直接读数"，只能由下一次真 run 或编排者授权的受控构建给出。
2. **`-Subset full` 那六档真采样**——本腿没跑（派单禁＋本机就是 runner）。归口＝AC#5 那次真 run。
3. **"往 `ci.yml` 塞一个 `-Subset` 值看门响不响"**——没跑（那是 262/134 的写面，本腿不越）。
4. **`Get-Variable` 那一形在 CI 的颜色**——本机量到"钉不响＋无名红"，CI 每步新进程下是否同形，没量。
5. **改后"缺件自建"那支对真 `scripts/build.ps1` 的一发**——没跑（会起 cgo 全量构建并动 `build/`）。本腿只由 `copy-noop` 证明副本机制不失真，这一格仍是〔未实测〕，与 r1 §⑧ 3 同格。
6. **CI 上 wisp.exe 死后的存活时长**——本机在改前那一发里直接看见"红句之后子进程才输出 `exiting`"（孤儿形状成立），CI 当时没人取进程表 ⇒〔未实测〕。
7. **golden 全名册 md5**——本腿射程只到"五枚 commit 名册零 golden"。

**记我（本腿自己的尺连撞四次，全数留下不当笔误藏）**：
- PS 里 `& child 2>&1` 在 `$ErrorActionPreference='Stop'` 下把子进程 stderr 变成 `NativeCommandError`，把我的台件自己打死（第一发）。
- 改成 `1> file 2> file` 后 `$LASTEXITCODE` **不落值**，读数里出现 `rc=`（空）（第二发）。
- 改 `Start-Process -ArgumentList $args`——`$args` 是**函数内自动变量**，赋值不生效 ⇒ `ArgumentList` 收到 null（第三发）。
- 再改 `[Diagnostics.ProcessStartInfo]` 属性赋值报 `The property 'FileName' cannot be found`（第四发，未查成因，直接弃用 PS 编排）。
⇒ 最终用 Python `subprocess` 起-等-取，退码由 `proc.returncode` 给。**定式建议入库：Windows 上编排 PowerShell 突变的台件用 Python，不要用 PowerShell 自己编。**
另抓到自己一枚好尺：`plant-in-comment` 之前我以为钉只会抓"读取"，实测它也抓注释 ⇒ **它比语义宽**，这条我写进 §⑤ 而不是藏掉。

**票 263 现在够格翻哪几枚框（翻勾由编排者做，本腿一枚都没碰；现跑尺：`grep -c '^- \[ \]'`＝**7**、`'^- \[x\]'`＝**0**，r1 与本腿都零碰）**：
- **够格翻**：`AC#0`、`AC#1`、`AC#2`、`AC#3`、`AC#6`。凭据分别＝本文 §①（名册＋三件事复认）、§②（三验齐＋差分）、§③（四条逐条有尺）、§④（三发全部本腿重跑＋还原证明）、§⑦（两门 rc=0 逐名＋归因独立验）。
- **够格翻但要带条件写清**：`AC#4`。按"本文件两枚自锁钉＋射程具名"这个口径成立；⛔ 不许写成"保证'未赋值即死'不再回来"（§⑥ 三形实测绕过）。翻勾凭据里那句"字面出现 1 次"应更正为"字面 2 行、带 sigil 0 处"（§③ 末）。
- **必须停在未勾**：`AC#5`。它要的是**一次真 CI run 里 `slo-full` 出结论＋`Upload SLO report` 不再 skipped**，推送在编排者手里（未推 79 枚）。本机一发不是凭据。
- **本票之外该新立的格子（我给判断，不代拍）**：票面 AC#6 写"卫生四门"，而本票实际适用的是两门＋两门"因零 Go 改动而不适用"——这不是缺口，但台账里最好具名一句，免得下一任以为 gofumpt/vet 真跑过。另：§② 末尾那格"改出厂件子系统要带下游消费者名册"的普查尺，属**新票射程**（不属 263/244 任何一方），要不要立由编排者裁。

**交件判据（占位符 0）——本腿现跑尺，读数档 `placeholder-check.txt`（13:09:40+0800）**：scope＝本目录 23 枚 tracked 件＋`scripts/slo-check.ps1`，pattern 由 UTF-8 八进制转义在运行时拼（⛔ 不写进被扫文档，否则尺自我命中——与票 263 词面钉同一枚形状）；**非零名册为空、总命中＝0**。
入库范围＝证据件＋假件源码＋两门读数＋十份突变日志＋两份台件脚本（含那枚失败的 PS 编排，作 §⑨ 的读数用）＋改前 397 行字节件；⛔ 未入库＝突变副本（防第二份 `slo-check.ps1`）、`*.exe`、`out-*` 采样目录——三者可由 `run_mutations_v1.py` 从 HEAD 字节逐字节重生成（锚不命中即抛）。
