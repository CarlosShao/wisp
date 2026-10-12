# 编队起手／终态闸门名册（票 304 `AC#4`＋`AC#4b` 的落点，2026-10-11 17:3x 由编排者定死）

**这份文件是什么**：从本刻起，每一份派单的"起手闸门"那一节**必须逐字抄下面 §1 那一段**，"终态锚"那一节**必须逐字抄 §2**。
**为什么在这里**：票 304 `AC#4` 要的是"该数的进程名定下来"，而编队的闸门文本一直是我每次现写、散在 `probes/orch/**` 的派单拷贝里
（`2026-10-09-live-window-runsheet.md:45`、`296-ac3-runsheet.md:7`、`a793-append.md:7`、`handover-40as.md:4`、`handover-40at.md:4` 那五枚是**历史快照，⛔ 回头改**——它们是当时说过什么的记录，改了就是造假）。
⇒ 定死一份正典，今后只抄它。

---

## 1. 起手闸门（跑真机／整包之前当场跑，缺一不跑）

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
★**这条尺的正控必须同行落**（10-12 实测事故：`20261011/logs/orch-A856-gate-family.txt` 那发的 `-join` 被外壳转义打断，PowerShell 报 `ParserError/ExpectedValueExpression`，而 `count=0` 照写 ⇒ 一枚**假 0**；正形＝把计数放在 PS Inside 里算好再打印，见 `probes/orch/20261012/gate-recheck-1012.ps1`）。
判"这把尺有牙"＝同管线跑一枚**必非零**的名字（10-12 现量：`svchost.exe`＝85／`explorer.exe`＝1／`*.exe`＝349，件 `logs/orch-gate-recheck-1012.txt` 里 `canonical_count=0` 那一发才是本条尺的正形读数）。

**临时目录那一族（票 310 `AC#4`，2026-10-12 11:1x 由编排者落政；⛔ 这一行⛔ 授权任何人删目录）**

```
powershell -NoProfile -Command '(Get-ChildItem -Path $env:TEMP -Directory | Where-Object { $_.Name -like "Test*" }).Count'   # 枚数＝%TEMP% 下测试临时目录残留数
```

- 现量序列（同一把尺）：10-10 立案 **21** → 10-11 19:5x **44** → 10-12 10:2x **50** → 普查腿交件 **53** → 编排者复跑 (b) 两发之后 **55** → `310-r1` 交件 **58** → 编排者复跑 winsec 两发之后**仍是 58**（件 `probes/orch/20261012/logs/orch-gate-20261012.txt`＋`orch-310-recheck-b.txt`＋`orch-310r1-precheck.txt`＋`orch-temp-after-my-shots.txt`）→ `310-v1` 交件 **59**（多那一枚＝它 `AC#2` (a) 那发的具名残渣 `TestRoundTripLoadMarshalLoad4222770164`，`CreationTime`＝`12:48:51.544`，编排者 13:0x 复跑到盘上仍在＝`logs/orch-310v1-recheck.txt`；归它、⛔ 删）。
- 这一行买的是"**看见它在长**"，⛔ 买的是"有人收尾"——收尾那半归票 310 的落地腿（`AC#1`）。
- ★**读数件⛔ 写进任何专用根里**（10-12 由 `310-v1` 现量：往 `D:/tmp/wisp-test-gotmp` 里落了两枚 `.txt` ⇒ 同一格两把尺失配 `Get-ChildItem -Directory`＝6 ⇄ `ls -1`＝8，"根里条目数"从此分不清残渣与文件）。定式＝**专用根只放 `t.TempDir()` 的产物**，读数件一律落 `probes/<票>/<腿>/logs/`；数根里条目时**具名写你用哪一把**（`-Directory` ⇄ `ls -1`）。
- ⚠ **名册件里含非 ASCII（本机实测＝一枚带中文的测试临时目录名经 PowerShell→bash 重定向落成 GBK 字节流）⇒ 数行的那把尺是 locale 敏感的**：10-12 编排者同一格三把尺并排＝PS `.Count`**59** ⇄ `grep -a -E '^Test'`**59** ⇄ `LC_ALL=C grep -cE '^Test[^|]*\|2026'`**59**，而**默认 locale 下那把严格尺报 58**（整行漏掉那一枚；同一件还出现过一次 `Binary file … matches`）。件内 **NUL＝0**、`utf-8` 解码在 offset 1192 失败（`invalid continuation byte 0xc9`）。定式＝**数行一律 `LC_ALL=C` 或 `grep -a`**；名册尺的枚数要与计数尺逐字对得上，对不上⛔ 报数、先怀疑尺；★**⛔ 从"哪把尺少报了一枚"倒推成因**（我第一版把成因写成"长行触发误判"，重跑才发现是 GBK 字节流——同族错见 `A858` §4）。件 `probes/orch/20261012/roster-1012.md`。
- ⚠ **枚数⛔ 是证据，名字才是；而名字只有"作差"那种来路算数**（10-12 实测：两枚腿把突变台件自报的 `dir=Test…<数字>` 当成漏的那枚抄进了交付件——同一秒里 `t.TempDir()` 会建了又删，台件盯到的那枚往往早已自我收尾；盘上真正活着的那枚是**另一个名字**）。⇒ 定式＝**起进程前的名册 ⇄ 杀完之后的名册作差（`comm -13`）＋逐枚带 `CreationTime`**；⛔ 引用任何"某某目录是我这一发漏的"之前，先用 `Get-ChildItem … | Select Name,CreationTime` 证它**今天还在盘上**。裁决腿⛔ 按交件里的名册盖章（本仓 10-12 两枚腿各栽过一次，见 `A857`／`A858`）。

**专用根落哪一枚卷＝会改闸门颜色的变量（10-12 编排者现量，票 310 `AC#1` 落地腿撞出、我复跑坐实；细节见票 311）**

- ⛔ 把 `WISP_TEST_GOTMP_ROOT`（或直接 `GOTMPDIR`）指到**仓所在卷 `D:`** 再跑 winsec／cli 档：`bash scripts/winsec-tests.sh ./internal/winsec/` 在 `D:` 根＝**PASS=56 FAIL=2**（`TestGateJudgesThePrivateSetByResolvedSID` 两枚子例＋`TestSealNarrowsAndNamesThePrincipalItRemovedBySID`），在默认 `%LOCALAPPDATA%` 兄弟根＝**PASS=58 FAIL=0**；同台面同码，只差那一枚卷。⇒ 改前／改后两发**必须用同一枚根形状**（默认那枚），⛔ 为了"仓外"好看把根挪卷。
- ⛔ 这一档的命令拼写＝`bash scripts/winsec-tests.sh ./internal/winsec/`。`--scope=winsec` 是 `portable-tests.sh` 那一层的拼法，喂给 `winsec-tests.sh` 会被它自己的 GUARD 1 打回 **rc=2**（10-12 编排者踩过，件 `logs/orch-winsec-dtmp-err.txt` 逐字两行）。
- 新增一把同行尺（⛔ 单独当判据）：`grep -c 'TempDir RemoveAll cleanup'`——`t.TempDir()` 的**自我收尾本身会失败**（10-12 现量：`D:` 根那发 4 行逐字 `… \la: Access is denied.`，`%LOCALAPPDATA%` 根那发 **0** 行，件 `logs/orch-winsec-A-dtmp.txt`／`-B-lappdata.txt`）⇒ 残渣⛔ 只来自"进程被杀"这一种死法。
- ★**这一族换卷改颜色的成因今天量到了，措辞要照它写**（10-12 `311-a1` 取证＋编排者六把复跑，件 `probes/311/a1/logs/**`＋`probes/orch/20261012/logs/orch-311a1-recheck-{A,B}.txt`）：差的⛔ 是判据、是**清理器能不能删**。默认根的父链给本账号 `swq:(I)(OI)(CI)(F)`（含 delete-child）⇒ 封存后的叶删得掉；`D:` 卷的父链只给 `Authenticated Users:(I)(M)`＋`Users:(I)(RX)`（Modify ⛔ 含 delete-child）⇒ 夹具封出的 protected DACL 叶（逐字 `D:PAI(A;;FA;;;LA)`，`LA`＝`…-500`＝**已禁用的内置管理员**）谁都不给访问 ⇒ 自我收尾必失败、**并把用例判红**。尺＝`icacls` 逐层＋`Get-Acl`；⚠ 三句判据（`gate accepted a descriptor`／`gate refused a descriptor`／`seal left a foreign grant`）四发各 **0 命中** ⇒ **今后任何人把这一族写成"判据不稳／判据依赖夹具位置"＝过期说法**（我 12:2x 就写过一次，更正见 `A860` §3）。
- ⚠ **令牌／SID 那类读数⛔ 用 `whoami` 走 bash**：10-12 编排者连坏三把（coreutils `whoami` 把 `/user` 当多余参数；MSYS 把 `/groups` 改写成 `C:/Users/…/git/groups`；`cmd //c "whoami /groups …"` 报"too many arguments"）⇒ 正形＝PowerShell 走 `-File`，里面 `[System.Security.Principal.WindowsIdentity]::GetCurrent()` 取令牌、`Get-LocalUser` 取账号状态、`Get-Acl` 取叶的 SDDL。⚠⚠ 同一条上我今天还犯了一次**"坏尺报的 0 当读数"**：`whoami_user_sid=` 那行是空的＝尺坏、⛔ 是"这台机没那枚 SID"。

**数到 ≠ 归我**。任何一枚非零计数，动手之前先答三问，答不全就**停下来上报**，⛔ 不许停：

1. **谁起的**——它的命令行里那个具名路径是什么（`Get-CimInstance Win32_Process -Filter "ProcessId=<pid>" | Select CommandLine`；
   `mockllm.exe` 的命令行自带它的 `t.TempDir()` 全路径，路径首段就是测试名；webview 系看 `--webview-exe-name=`）。
2. **父 PID 还在不在**——`ParentProcessId` 现量查得到吗？⚠ **这一问只能当旁证，不能当归属**：
   Windows 会**复用 PID**（票 304 现量：曾被判"已不存在"的 `11724` 后来查得到＝`RuntimeBroker.exe`，只是同号换了人）。
   ⇒ 归属只认**命令行里的具名路径**；路径认不出来的＝**归属不明**，⛔ 归本腿。
3. **本腿认下哪几枚**——逐枚写 PID＋来路；写不出来的那几枚在报告里标"别家的，⛔ 动"。
4. ★**派单里只要让一枚腿跑 `go test`，而编队里有另一枚腿正在改 `scripts/**`／`tools/d22scan/runtests.sh`，就要具名写"⛔ 走那三枚入口、直接 `GOTMPDIR=<你的根> go test -v -count=1 ./<包>/`"**（10-12 实测编队：`310-r2` 在改入口、`311-a1` 在取证，两件事撞在同一批文件上；今天没撞坏是靠派单里那句"别用那三枚脚本"，⛔ 靠运气）。
5. ★**名册尺的枚数要与计数尺逐字对得上；对不上先怀疑尺、⛔ 报数**（10-12 编排者自抓：同一格第一发读到 **58**、真值 **59**，成因＝一枚**带中文的测试临时目录名**经 PowerShell→bash 重定向落成 **GBK 字节流**（件内 NUL＝0、`utf-8` 解码在 offset 1192 失败）⇒ UTF-8 locale 下 `grep` 整件判 binary 且严格正则**整行漏掉那一枚**。定式＝**数行一律 `LC_ALL=C` 或 `grep -a`**；⚠ 我第一版把成因写成"长行触发误判"＝⛔，成因也从"哪把尺少报一枚"倒推不出来（同族错＝`A858` §4／`A860` §6）。

## 2. 终态锚那一侧（`99-final-anchor.md`／停车点里那句进程计数）

同一枚尺、同三问，另加两条硬约束：

- 每一枚非零计数都要**指到一棵树**（根进程名＋根 PID 现量在不在）。指不到就写"**归属不明**"，⛔ 把裸枚数写进终态锚。
  ★**"指到一棵树"＝沿父链往上爬到一枚能命名的根为止**（10-12 由编排者补，来路＝`done-check-2` 报本格⛔ 被现形过）。只爬一跳会把**孙辈**误写成归属不明；爬到根又可能一路爬到 `explorer.exe` 那种"什么都是它"的根 ⇒ 定式＝**爬到第一枚"名字或命令行能具名到某个宿主/某个仓"的祖先为止，并逐枚写下爬了几跳**；爬不到＝归属不明。
  ⚠ 本格原先只在**全 0 台面**上验过（`304-v1b` 自报"未被现形"）⇒ 那条句当时是枚**恒真尺**。今天非零台面已在场（10-12 现量 `logs/orch-gate-recheck-1012.txt`：`msedgewebview2.exe`＝**24**，其中**命令行⛔ 带 `--webview-exe-name=`＝4 枚**、其父本身是 webview ⇒ 必须再上一跳才指到树；父已消失＝**0** 枚）⇒ 这条句**现在会被咬到**，判语见台账 `A857`。

- ⛔ **由裸枚数生成"停法"类欠账**。实测事故链：把别家的 `msedgewebview2.exe=18` 报成"18 枚孤儿"⇒ 台账里就多了一笔要"停"的假账；
  那 18 枚**数字对、标签错**（逐树复量＝`SearchHost.exe`／`wetype_update.exe`／`clipsync-desktop.exe` 各 6 枚、三枚父进程全在、wisp 系 0）。
- ⚠ **在飞期间"本机枚数"里必然混着腿自己的窗**（同一台面上 62 秒内从 28 枚含 2 棵 `wisp.test.exe` 树变成 18 枚）⇒ ⛔ 拿总枚数当残留数。

## 3. 现量（这一发改完当场用它跑的那一发）

尺＝上面 §1 那四条 `tasklist` 各一发的命中行数＋`Get-CimInstance` 逐枚命令行，见 `20261011/logs/orch-A85x-gate-live.txt`（`gate_rc=` 行自落）。
