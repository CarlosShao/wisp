# 450 — 整包「约 400 秒死亡」：可判事实与二分（本机探针，2026-11-15）

**状态**：探针完成。AC#1 / AC#2(a) / AC#3 有读数；AC#2(b) 按票面**具名留给机主**（做它必须停球）。
**锚点**：`caccccfc6`。**代码面与死亡那两发相同**——`git diff --stat cf3050401 HEAD -- cmd internal` **空**（cf3050401＝昨 22:xx 那两发的未提交树在 23:41 落成的提交），所以下面的读数不是"代码变了"造成的。
**球**：`wisp.exe` PID **49348**，`2026-11-14 16:23:02` 起，**本波全程在跑、一次未碰**（未 kill / 未停 / 未重启）。
**载具**：`scripts/wisp-cli-staged-tests.sh`（一次一发；本波**五发串行**：两发整包 + AC#3 的三臂单独跑，每发单独 stage 目录）。

---

## §0 一句话结论（先看这四行）

1. **死没复现：没复现。** 球在跑、同一条命令、同一代码面，**两发整包都跑到底**：`RUN=861 / PASS=619 / FAIL=2 / SKIP=17`，**exit 1**，最后一行 `FAIL`，载具判行 `FAIL github.com/CarlosShao/wisp/cmd/wisp 1207.00s` 与 `1308.00s`。两枚 FAIL 都是本机具名已知红（§1.1）。
2. **昨日那两发不是"崩溃"，是静默消失**：Windows 事件日志里**0 条**提到 `wisp.test.exe`（Application 日志覆盖到 2026-07-30、共 20,715 条；对照：同机近 30 天**有 26 条 `Application Error`**）⇒ 没有 WER 记录 ⇒ **不是**未处理异常／访问冲突／fail-fast。这一条同时**削弱"WebView2 不稳"那一支**：WebView2 host 的 fail-fast 会写 WER 1000。
3. **球争用：否**（作为"整包在这台机上跑不完"的原因）。证据：球在跑时整包**今天两次跑完**，且**昨晚 21:17–21:48 还有一发跑完**（§2.2）。若问的是"球是不是昨日两次死亡的**触发因素**" ⇒ **证据不足**（AC#2(b) 未做，票面禁止本腿 kill 球）。
4. **留给机主的那条命令**（球一停就跑这一条；球没停它会自己拒绝，**不 kill 任何东西**）：

```
powershell -NoProfile -ExecutionPolicy Bypass -File ".scratch\wisp\probes\450\run-balloff.ps1"
```

---

## §1 AC#1 把「死」变成可判事实

### §1.1 今天两发整包（球在跑）——逐项读数

| 发 | 起–止 | 耗时 | exit | RUN | PASS | FAIL | SKIP | 载具判行 | 最后一行 | 死? |
|---|---|---|---|---|---|---|---|---|---|---|
| 1 | 10:01:28–10:21:49 | 1221 s（测试段 1207 s） | **1** | 861 | 619 | 2 | 17 | `FAIL github.com/CarlosShao/wisp/cmd/wisp 1207.00s` | `FAIL` | **否** |
| 2 | 10:22:14–10:44:12 | 1318 s（测试段 1308 s） | **1** | 861 | 619 | 2 | 17 | `FAIL …/cmd/wisp 1308.00s` | `FAIL` | **否** |

两发命令**逐字同一条**：

```
cd "D:\work\workspace\projects plans\Wisp"
$env:WISP_CLI_STAGE_DIR=".scratch\wisp\probes\450\stage-full1"   # 第 2 发是 stage-full2
& "D:\work\soft\Git\bin\bash.exe" scripts/wisp-cli-staged-tests.sh
```

原始读数：`.scratch/wisp/probes/450/run1.out.txt`（1,266,835 B / 7,567 行）、`run2.out.txt`（1,265,678 B / 7,556 行）；
每发的 `-test.v` 全流另存 `stage-full1/wisp-cli-tests.log`（1,258,559 B）、`stage-full2/wisp-cli-tests.log`（1,257,413 B）。

**四数怎么数的**：与 `scripts/portable-tests.sh:1238-1244` 同一把尺（`^=== RUN` / `^--- PASS` / `^--- FAIL` / `^--- SKIP` 行数，子测试一并计）。

**两枚 FAIL 具名（都不是本波、也不是本票的红）**：

| 用例 | 第 1 发 | 第 2 发 | 归谁 |
|---|---|---|---|
| `TestAC4FocusReturnToPriorWindowGap33r5` | FAIL 5.20 s | FAIL 5.06 s | 本机前台被占的**已知环境红**（票面已具名） |
| `TestAC246DevLegIgnoresTheTestTaskInjection` | FAIL 37.63 s | FAIL 35.74 s | 驻留球在跑的**已知红**，红句逐字：`AC#7 RED (desktop state, not our code): a dev Wisp is already running in this session…`（票面已具名） |

**这两发覆盖到昨日两发从没跑到的那一片**：`Test145ArgsJSONComesFromTheBookedForensicsRow`（PASS 0.17 s）与 `TestAC4PriorFocusSurvivesARefusedPanelSample`（PASS 1.35 s）都**跑过并变绿**，`r`–`v` 区（`resident_*` / `ticket*` / `voice_*`）也跑到了底（最后一枚 `TestZZ315A1…` PASS 0.35 s）。

### §1.2 事件日志臂——昨日那两发的可判事实

查法（可逐字重跑）：

```
Get-WinEvent -FilterHashtable @{LogName='Application'; StartTime=(Get-Date '2026-11-14 21:50'); EndTime=(Get-Date '2026-11-14 23:45')} |
  Sort-Object TimeCreated | Select-Object TimeCreated,ProviderName,Id,LevelDisplayName
Get-WinEvent -FilterHashtable @{LogName='Application'; StartTime=(Get-Date).AddDays(-7)} |
  Where-Object { $_.Message -match 'wisp\.test' }
```

| 读数 | 值 |
|---|---|
| 昨日死亡窗（21:50–23:45）Application 日志条目 | **7 条，全部无关**：6 条 `Security-SPP 16384/16394`（服务重排）＋1 条 `igccservice PowerEvent` |
| 其中 `Application Error` / `Windows Error Reporting` / `Application Hang` | **0 条** |
| 近 7 天全库检索消息含 `wisp.test` | **0 条** |
| 对照（证明记录机制没关） | 近 30 天 `Application Error` **26 条**，如 `WebHostComponent.exe`（5 条）、`wetype_server.exe`（ntdll.dll） |
| Application 日志覆盖 | 最早 `2026-07-30 08:36:18`，共 **20,715** 条 ⇒ 昨日窗口**在覆盖范围内** |

⇒ **可判事实（不是描述）**：昨日两发**没有留下任何 Windows 错误报告**。没有 WER 1000 就排除了"未处理异常／访问冲突／fail-fast"这一整类（它们必进 WER）。剩下的形状是**进程被正常或强制收回**（`ExitProcess` 或外部 `TerminateProcess`），而这两者**只能靠退出码区分**，昨日两发的退出码**没有留存**（见 §6.1）。**本票据此不改任何产码。**

### §1.3 「WebView2 窗那一区」今天量到的实际形状（解释"为什么越跑越重"，**不写成死因**）

整包期间 msedgewebview2 进程数（全机）与合计工作集，15 s 一个采样点
（`.scratch/wisp/probes/450/mem.csv`）：

| 时刻 | 距测试开始 | msedgewebview2 枚数 | 合计工作集 | 可用内存 |
|---|---|---|---|---|
| 10:06:31 | 289 s | **37** | 3,369 MB | 8,441 MB |
| 10:08:19 | 397 s | **47** | 3,921 MB | 8,365 MB |
| 10:09:06 | 444 s | 63 | 4,658 MB | 8,085 MB |
| 10:10:39 | 537 s | **112**（峰值） | 7,363 MB | 7,440 MB |
| 10:20:07（测试尾段） | 1,105 s | 82 | 5,904 MB | 7,566 MB |
| 10:22（进程退出后实测） | — | **37** | 3,433 MB | — |

- 37 这个基线不是空的：球自己 6 枚（16:23 起）＋早先别的腿留下的 25 枚（18:53 / 19:20 / 20:45 三批）。
- 整包里创建的浏览器进程**是测试进程的子进程**，实测 `ParentProcessId=28040`、`--webview-exe-name=wisp.test.exe`、每枚独立 `--user-data-dir=%TEMP%\wisp-33r5-panel-profile-28040-…\EBWebView` ⇒ **不跨进程泄漏**：进程一退，82 → 37（上表最后一行）。
- **「400 秒」不是墙钟触发器**：今天 t+397 s 时刚建起第 2–3 个真窗口（37→47），而死亡窗落在这段"第一次大量建窗"里（昨日 403 s / 416 s）。这是**观察**，不是死因判断。

---

## §2 AC#2(a) 球在跑那一支（四数＋死没死）

就是 §1.1 两发（AC#1 的复现发与 AC#2(a) 的对照发是**同一条件**：球 PID 49348 在跑、命令相同、串行）：`RUN=861 / PASS=619 / FAIL=2 / SKIP=17`，**两发都没死**，exit 1（红在两枚已知红上）。

### §2.2 历史对照（本波没跑，读的是盘上日志）

`.scratch/wisp/stage-447b/wisp-cli-tests.log`（2026-11-14 **21:17:48 → 21:48:05**，1,241,866 B）：
`RUN=853 / PASS=611 / FAIL=2 / SKIP=17`，同样以 `FAIL` 收尾，两枚 FAIL 是同一对本机已知红
（`TestAC4FocusReturnToPriorWindowGap33r5` 5.42 s、`TestAC246DevLegIgnoresTheTestTaskInjection` 37.42 s）。
**那一发的球也在跑**（PID 49348 起于 16:23，进程至今未换）⇒ **"球在跑"与"整包能跑完"不矛盾，昨晚 21:48 就同时成立过**。

---

## §3 AC#2(b) 球不在跑那一支——**留给机主的就这一条命令**

**为什么不做**：做它必须停球，票面逐字 `⛔ 自己 kill 机主的球`；票里给的替代（"从未起过球的干净用户会话"）本腿拿不到（不能开新用户会话），也不该冒这个险。**本腿一次都没碰 PID 49348。**

### 3.1 命令（球从托盘退出后，跑这一条）

```
powershell -NoProfile -ExecutionPolicy Bypass -File ".scratch\wisp\probes\450\run-balloff.ps1"
```

脚本：`.scratch/wisp/probes/450/run-balloff.ps1`（本波已自测它的**拒绝支**）。

### 3.2 它会做什么、不会做什么

- **不会**：kill / stop / suspend / 重启球，也不读球的内存。只**看**仓库根有没有活着的 `wisp.exe`。
- 球还在 ⇒ **拒绝，exit 3**，一个字都不跑（本波实测输出：`REFUSED - the resident ball is still running: PID 49348 (started 11/14/2026 16:23:02)`）。
- 有另一发 `wisp.test.exe` 在跑 ⇒ 同样拒绝（exit 3），因为本机规矩是整包**一次一发**。
- 球不在 ⇒ 用**同一个载具、同一条命令**跑整包（stage 目录 `probes/450/stage-balloff`），并自己算四数、退出码、最后一行、判行，追加到 `probes/450/balloff-readings.txt`。

### 3.3 判据（跑完怎么读）

| 读数 | 结论 |
|---|---|
| 球不在 + **跑到底**（`RUN=861 PASS=619 FAIL=2 SKIP=17`，exit 1，两枚已知红） | 与今天球在跑的两发**同形** ⇒ 球争用**否**，两天的死亡都是间歇现象，需按"静默退出"继续二分 |
| 球不在 + **仍死**（无 `--- FAIL` 行、`FAIL` 判行、PASS 停在 ~330–340） | 球争用**否**，且死亡与球无关 ⇒ 直奔"静默退出"的两种机制（§5.2） |
| 球不在 + **死得更早/更频繁** | 说明"球在跑"反而把时序推开了，同样否掉"球争用"这个方向 |

---

## §4 AC#3 二分：把死亡点附近那几枚**单独**跑（无整包负载）

命令逐字（载具 + 全名锚定的 `-test.run`，一次一发；读数 `.scratch/wisp/probes/450/ac3-readings.txt`）：

```
& "D:\work\soft\Git\bin\bash.exe" scripts/wisp-cli-staged-tests.sh "-test.run=^Test145ArgsJSONComesFromTheBookedForensicsRow$"
& "D:\work\soft\Git\bin\bash.exe" scripts/wisp-cli-staged-tests.sh "-test.run=^TestAC4PriorFocusSurvivesARefusedPanelSample$"
& "D:\work\soft\Git\bin\bash.exe" scripts/wisp-cli-staged-tests.sh "-test.run=^(Test145ArgsJSONComesFromTheBookedForensicsRow|TestAC4PriorFocusSurvivesARefusedPanelSample|TestAC4FocusReturnToPriorWindowGap33r5)$"
```

| 臂 | 用例 | 用时（含建二进制） | exit | RUN/PASS/FAIL/SKIP | 逐枚结果 | 死? |
|---|---|---|---|---|---|---|
| A1 | `Test145ArgsJSONComesFromTheBookedForensicsRow` 单独 | 9 s | **0** | 1/1/0/0 | PASS 0.22 s | **否** |
| A2 | `TestAC4PriorFocusSurvivesARefusedPanelSample` 单独 | 12 s | **0** | 1/1/0/0 | PASS 1.27 s | **否** |
| A3 | 死亡区三枚一起（另加 `TestAC4FocusReturnToPriorWindowGap33r5`） | 14 s | **1** | 3/2/1/0 | PASS 0.18 s／PASS 0.79 s／**FAIL 4.98 s（已知环境红）** | **否** |

**A1/A2 的载具判行**：`ok github.com/CarlosShao/wisp/cmd/wisp 1.00s` / `2.00s`（exit 0）；A3：`FAIL … 6.00s`（exit 1，红的只有那枚已知环境红）。

### §4.1 一条必须写下来的判读（否则这张票会被读歪）

`Test145ArgsJSONComesFromTheBookedForensicsRow` 的**函数体里没有窗**：它 `memory.Open` + `panel.NewRunFeed` + `panel.NewSnapshotPump.Marshal`（`cmd/wisp/panel_run_145_test.go:217-320`），**不调** `startPanelForTest`。
⇒ 昨日"死在 WebView2 窗那一区"是**从文件序 `panel_*` 推出来的**，**不能**读成"这枚用例自己开窗、自己崩"。
⇒ 单独跑它**本就不该复现**，A1 的 0.22 s PASS 正是这个意思；今天两发整包在同一枚上也没死。
**结论（AC#3）**：这三枚**不是**死因的载体；没有整包负载时它们全部按预期出结果（唯一红是本机具名环境红）。

---

## §5 归属结论与被排除的方向

### §5.1 结论

- **死没复现**：是（两发整包 + 三臂单独跑，**五次载具调用全都没死**）。
- **球争用**：**否**（球在跑时整包今天两次跑完、昨晚一次跑完）。作为昨日两发死亡的**触发因素**：**证据不足**（AC#2(b) 未做）。
- **WebView2 不稳**：**证据不足，且有反向证据**（昨日两发无 WER 记录；WebView2 host 失败会写 WER）。同时今天整包把 WebView2 区（含峰值 112 枚浏览器进程、7.4 GB 工作集）**跑完了**。
- **代码面**：**已排除**（`git diff cf3050401 HEAD -- cmd internal` 空）。
- **WebView2 运行时更新**：**已排除**（`C:\Program Files (x86)\Microsoft\EdgeWebView\Application` 最新 `157.0.4341.46`，2026-11-11 10:00 建；昨日 21:30–00:10 无任何 Edge/Update/MsiInstaller 事件）。

### §5.2 还没排除的（具名，**不猜着改产码**）

1. **静默退出的两种机制**：`ExitProcess`（正常退出码）vs 外部 `TerminateProcess`（异常码）。要判**必须复现并抓退出码**——`scripts/wisp-cli-staged-tests.sh` 已经把二进制退出码打进 `bin_rc` 与判行，下一发复现时把 `run*.out.txt` 尾部三行留档就能分开。
2. **并发负载**：昨日 22:xx 两发前后本机有其它腿在跑，本波期间同机也有 ≥3 条腿。**但本波没有一条读数把并发与"无 WER 记录的静默消失"连起来**：run2 里 `TestTicket223…` 处两枚子 `go build`（`Get-CimInstance Win32_Process` 抓到 `go.exe build -o …\Temp\TestTicket223…`）在 10:25:11–10:25:39 之间跑完（**< 30 s，只是慢，不是挂**），整包仍然跑到底。⇒ 并发是"本机为什么规定一次一发"的理由，**不是**本票能拿去定死因的证据。
3. **桌面会话热/冷**：昨日 9:10 有过一次睡眠唤醒（两发在 22:16 之后）；今天 9:51:49 唤醒、两发在 10:01 起跑。**同一天内两种状态都有整包跑完/死亡的读数**，今天不支持它单独成因。

---

### §5.3 门禁（AC#5：本票只写证据件／票面 Progress log，按"脚本注释或证据件"口径跑两条）

| 门 | 命令 | 读数 |
|---|---|---|
| D22 静态禁令 | `D:\work\soft\Git\bin\bash.exe scripts/d22scan.sh` | **rc=0**，`d22scan: clean - no D22 ban violations`（扫 484 枚产码 Go：internal/ 405、cmd/ 79；`d22scan.sh: OK - control rc=0, scan rc=0`） |
| 路径长度预算 | `… bash scripts/check-path-length-budget.sh` | **rc=0**，`VERDICT GREEN`（tracked=20346、over-budget=68、roster=68、not-in-roster=0） |

原始输出：`probes/450/gate-d22.log`、`probes/450/gate-pathlen.log`。

---

## §6 没做到 / 缺口（具名）

1. **昨日两发死亡的退出码与最后一行没有留存**：`.scratch/wisp/stage-cli/wisp-cli-tests.log` 已被 2026-11-14 23:31 之后的运行覆盖（现在那份是 2026-11-15 09:54 的定向发，26 行）。**没有**从任何残件里"猜"出这两个数。
2. **AC#2(b) 未做**（票面禁止本腿停球），替代的"干净用户会话"本腿也拿不到。
3. **AC#4 的"写进测试前置说明"没做**：它只在"球在跑死、球停不死"成立时才该写；本波读数是"球在跑也能跑完"，照 AC#4 的文字不成立，所以**`scripts/wisp-cli-staged-tests.sh` 一个字未改**（连注释都没动）。
4. **本波不改产码**：`git status --porcelain -- cmd internal` 空（见提交前核对的读数）。
5. 本机两枚已知红**原样保留**：未修、未删、未跳过、未加 skip 名单。

---

## §7 撤销口令

「**450 撤**」＝删本文件、删 `.scratch/wisp/probes/450/**`、把本票 Progress log 那一段留着不动（只追加的账不回收）。**不动** `cmd/`、`internal/`、`scripts/`、`docs/specs/**`（本波一处未改）。
