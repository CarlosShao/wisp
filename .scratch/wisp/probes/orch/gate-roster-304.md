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
- ⛔ **由裸枚数生成"停法"类欠账**。实测事故链：把别家的 `msedgewebview2.exe=18` 报成"18 枚孤儿"⇒ 台账里就多了一笔要"停"的假账；
  那 18 枚**数字对、标签错**（逐树复量＝`SearchHost.exe`／`wetype_update.exe`／`clipsync-desktop.exe` 各 6 枚、三枚父进程全在、wisp 系 0）。
- ⚠ **在飞期间"本机枚数"里必然混着腿自己的窗**（同一台面上 62 秒内从 28 枚含 2 棵 `wisp.test.exe` 树变成 18 枚）⇒ ⛔ 拿总枚数当残留数。

## 3. 现量（这一发改完当场用它跑的那一发）

尺＝上面 §1 那四条 `tasklist` 各一发的命中行数＋`Get-CimInstance` 逐枚命令行，见 `20261011/logs/orch-A85x-gate-live.txt`（`gate_rc=` 行自落）。
