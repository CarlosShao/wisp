# 304-a1 — (A) 当前残留名册重量（只量⛔ 动手）

尺族＝`Get-CimInstance Win32_Process`（经 `powershell -NoProfile -Command`，内联⛔ `$_`，全部用简化形 `Where-Object Name -match/-like`），交叉尺＝`tasklist //FI "IMAGENAME eq <名>"` 命中行数。件＝`logs/proc-4fam-list.txt`／`logs/proc-4fam-cmdline.txt`／`logs/proc-parents-and-snapshot.txt`／`logs/proc-loose-scan.txt`。时刻＝2026-10-11 15:1x。

## 1. `mockllm.exe`：现量 **0 枚**

| 尺 | 现跑 | 快照（派单／票面） |
|---|---|---|
| `Name -match "^(mockllm\|wisp\|balldebug\|msedgewebview2)\.exe$"`（`logs/proc-4fam-list.txt`） | `mockllm.exe` **0 行** | 票面 `:9-:10`＝**2 枚**（pid 20904 起 2026-10-08 08:42:21／pid 25660 起 11:35:29） |
| 宽松尺 `Name -match "mockllm\|wisp\|balldebug\|test\.exe"`（`logs/proc-loose-scan.txt`） | **0 行**（含 `*.test.exe` 一枚⛔ 在跑） | — |
| `tasklist` 命中行数（交叉尺） | **0** | — |
| 逐枚 PID 尺 `Win32_Process -Filter "ProcessId=20904 OR ProcessId=25660 OR …"`（`logs/proc-parents-and-snapshot.txt`） | **20904／25660 两枚都⛔ 在名册里**＝都不在了 | 票面 `:17` 已记"处置已发生（立案时按具名 `-Id` 停掉），改后现量＝`mockllm.exe` 0 枚" |

⇒ **重跑值 0 ⇄ 派单快照 2 ⇄ 票面立案后自量 0**：本腿读到的与票面 `:17` 那一句**一致**，与派单 §2 引的 `:9` 那句（2 枚）**不一致**——那半句是**立案时的历史读数**，⛔ 是现值（这正是票面 `:7` 自己写着的"⛔ 引用前先重跑"那一枚雷，本腿按现跑交回）。

**父 PID 那一栏的更正**（票面 `:10` 说两枚 PPID `11724`／`20036` "当时都查不到"）：现量尺（`logs/proc-parents-and-snapshot.txt`）里 **`11724` 现在存在＝`RuntimeBroker.exe`，`CreationDate=2026-10-11 7:22:36`** ⇒ **PID 已被回收另用**，⛔ 那枚还是当初的父进程；`20036` 仍⛔ 在名册里。这一条只影响"父在不在"这把尺的可靠性（**PID 会复用**⇒归属必靠命令行具名路径，⛔ 靠号），⛔ 改本票任何结论。

**占用/套接字那一栏⛔ 复量**：票面 `:10` 的工作集/private/累计 CPU/`127.0.0.1:54887`·`58918` 两枚 Listen 都挂在**那两枚已不存在的进程**上 ⇒ 本腿⛔ 能复量（进程没了）。＝**具名欠读数 `Q-304-a1-F`**。

## 2. 另三族现量（归属三件逐枚：谁起的／命令行具名路径／父链在不在）

| 族 | 枚数（CIM 尺） | 交叉尺 `tasklist` | 归属 |
|---|---|---|---|
| `wisp.exe` | **0** | **0** | — |
| `balldebug.exe` | **0** | **0** | — |
| `msedgewebview2.exe` | **18**（＝3 棵树 × 每棵 6 枚） | **18** | **三枚树根逐枚具名，父链全在，wisp 系 0**（见下表） |

三棵树（尺＝`logs/proc-4fam-cmdline.txt` 逐枚 `CommandLine` 里的 `--webview-exe-name=` ＋ `logs/proc-parents-and-snapshot.txt` 的父进程名册）：

| 树根 PID | 父 PID／父名（现量在不在） | 命令行里的具名归属（逐字取段） | 树内枚数 | `CreationDate` |
|---|---|---|---|---|
| `12812` | `9004`＝`SearchHost.exe`（**在**，起于 2026-10-07 08:36:28） | `--webview-exe-name=SearchHost.exe --user-data-dir="C:\Users\swq\AppData\Local\Packages\MicrosoftWindows.Client.CBS_cw5n1h2txyewy\LocalState\EBWebView"` | 6（`12812`＋`12948`/`12536`/`12524`/`12412`/`13376`） | 2026-10-07 08:36:30 |
| `28208` | `16852`＝`wetype_update.exe`（**在**，起于 2026-10-07 08:36:53） | `--webview-exe-name=wetype_update.exe --user-data-dir="C:\Users\swq\AppData\Local\flutter_webview_windows\wetype_update\EBWebView"` | 6（`28208`＋`30272`/`3460`/`14280`/`28256`/`28516`） | 2026-10-08 15:50:33 |
| `13188` | `8548`＝`clipsync-desktop.exe`（**在**，起于 2026-10-10 18:30:46） | `--webview-exe-name=clipsync-desktop.exe --user-data-dir="C:\Users\swq\AppData\Local\com.clipsync.desktop\EBWebView"` | 6（`13188`＋`34272`/`34208`/`31572`/`32860`/`29940`） | 2026-10-10 18:30:46 |

每棵树的 6 枚子代都带同一枚 `--webview-exe-name=` 与同一 `--user-data-dir`（`--type=crashpad-handler`／`gpu-process`／`utility network`／`utility storage`／`renderer`），⇒ **归属是命令行背书的，⛔ 枚数推出来的**。

**重跑 ⇄ 快照**：票面 `:17`＝24 枚（四棵树，含 `magpie-windows-amd64.exe` 那一棵）；票面 `:30`（`AC#4b` 引的 `forensics-4-trees.txt`，2026-10-10 18:18）＝`total=18 trees=3`。**本腿现量＝18 枚／3 棵**＝与后者同值。差异解释：`magpie-windows-amd64.exe` **本身还活着**（尺＝`tasklist //FI "IMAGENAME eq magpie-windows-amd64.exe"` 现量 **pid 24208**），但它名下**现量 0 枚** `msedgewebview2` ⇒ 24−18＝那一棵的 6 枚已经没了；**为什么没**（机主自己关的／窗退了渲染进程跟着退／别停的）本腿⛔ 能证 ⇒ **`Q-304-a1-G`**。

⇒ **本腿按票面 `:34` 那句边界执行：这 18 枚⛔ 是本票的残留**（父链全活、命令行⛔ 指 wisp），**且本腿一枚⛔ 停**（派单 §2 硬令；处置三形归 `AC#1` 裁）。

## 3. 同一场盘上的第二症状（票面 `:16`，顺手复量）

尺＝`Get-ChildItem -Directory $env:TEMP | Where-Object Name -like "Test*" | Measure-Object`（⛔ `$_`）：现量 **Count=21** ⇄ 票面 `:16` 快照 **21 枚** ⇒ **同值**。⚠ 归属⛔ 做：目录还在、起它的进程已经⛔ 在名册里 ⇒ "哪一枚目录是哪一次杀留下的"这把尺**看不见**（票面 `:16` 自己点过那半把尺落空的日期栏）。⇒ `Q-304-a1-H`。

## 4. 这格"看不见什么"（派单 §4：命中 0 ⛔ 等于没问题）

- CIM 名册尺只回答"**现在**有没有"，⛔ 回答"这一族会不会再长出来"——那要表②的实测（本腿⛔ 做）。
- 名字尺⛔ 能抓的形：任何**改名**的假大脑（例：`go build -o <tmp>/mockllm2.exe`）、**linux 形**的 `mockllm`（无 `.exe`，ubuntu CI 那一支）、以及被别的名字包装的同一家族。本腿用宽松尺（§1 第二行）打了 `mockllm|wisp|balldebug|test\.exe` 四条，**两把尺都是 0 命中** ⇒ 0 命中只说明"这台机器这一时刻⛔ 挂着这一族"，⛔ 说明机制已修。
- 两把尺（CIM 与 `tasklist`）**同时放行**同一读数（18／0／0／0）⇒ 按派单 §4 那句"一枚尺两种形状都放行⛔ 是凭据"，本格的凭据⛔ 是"两把尺一致"本身，而是**逐枚命令行的具名路径**（§1、§2 那些 `--webview-exe-name=`／`golden-dir` 段）。
