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

- 现量序列（同一把尺）：10-10 立案 **21** → 10-11 19:5x **44** → 10-12 10:2x **50** → 普查腿交件 **53** → 编排者复跑 (b) 两发之后 **55**（件 `probes/orch/20261012/logs/orch-gate-20261012.txt`＋`orch-310-recheck-b.txt`）。
- 这一行买的是"**看见它在长**"，⛔ 买的是"有人收尾"——收尾那半归票 310 的落地腿（`AC#1`）。

**数到 ≠ 归我**。任何一枚非零计数，动手之前先答三问，答不全就**停下来上报**，⛔ 不许停：

1. **谁起的**——它的命令行里那个具名路径是什么（`Get-CimInstance Win32_Process -Filter "ProcessId=<pid>" | Select CommandLine`；
   `mockllm.exe` 的命令行自带它的 `t.TempDir()` 全路径，路径首段就是测试名；webview 系看 `--webview-exe-name=`）。
2. **父 PID 还在不在**——`ParentProcessId` 现量查得到吗？⚠ **这一问只能当旁证，不能当归属**：
   Windows 会**复用 PID**（票 304 现量：曾被判"已不存在"的 `11724` 后来查得到＝`RuntimeBroker.exe`，只是同号换了人）。
   ⇒ 归属只认**命令行里的具名路径**；路径认不出来的＝**归属不明**，⛔ 归本腿。
3. **本腿认下哪几枚**——逐枚写 PID＋来路；写不出来的那几枚在报告里标"别家的，⛔ 动"。

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
