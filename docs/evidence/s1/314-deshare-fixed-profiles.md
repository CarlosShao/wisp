# 票 314 — 把最后三处固定名 WebView2 user-data 夹改成每发独立（写腿；证据件）

> 角色：实现者（⛔ 审计者）。时间 `2026-11-14 11:33–11:5x +08`。树 = `dev`，HEAD `1ce39a651`。
> 写面：`cmd/wisp/panel_host_windows_test.go`、`cmd/wisp/panel_assets_299_windows_test.go`、
> `.scratch/wisp/probes/314-deshare/**`、本件、票 314 Progress log。
> 读数根：`.scratch/wisp/probes/314-deshare/{results-concurrent.jsonl,results-sweep.jsonl,logs/**}`；
> 测试二进制（2×62MB，**不入仓**）在 `D:\tmp\wisp\20261114\314-deshare\bin\`。

---

## 0. 一句话结论（含我拿不到的那一格）

**三处固定名夹已改成每发独立的叶子（父目录仍 `os.TempDir()`，仍不做 RemoveAll，另加启动时清 >1h 同族旧夹）。并发读数：de-share 树上两发同时跑同一枚真建窗用例，两发 rc=0，盘上两发各自命名**不同** UDF、同一 tick 同一 UDF 从未出现两个顶层 host；把三处用 `-overlay` 临时改回固定名重编同一发，同一形状下**同一 tick 同一 UDF 出现两个顶层 host（连续 2 tick，parents 恰是两发的 pid）**。**

**但派单要求的字面形状 —— 两发各跑两枚真建窗用例、两发都 rc=0 —— 我拿不到**：该形状 4/4 对里必有 1 发以**原生异常**消失（`0xC0000005` ×3、`0xC000041D` ×1；无 `panic:`、无 `--- FAIL`），而且**把三处改回固定名的构建在同形状里也死一发（同一退出码）** ⇒ 这一死法不是夹名造成的，本次改动既没引入它、也没治好它。我拿到 rc=0 的并发形状是**两发各建一窗**（3/3 对、6/6 进程 rc=0）。逐字读数在 §4/§5，不修饰。

---

## 1. 改了什么（三处，形状照 `68275983e`）

### 1.1 新 helper（落在 `panel_host_windows_test.go`，imports 区之后）

| 函数 | 作用 |
|---|---|
| `panelProfileLeaf(prefix string) string` | 先清同族旧夹，再返回 `uniqueProfileLeaf(prefix)`（**复用** `panel_resident_windows_test.go:61` 的先例函数，⛔ 没另造命名法） |
| `sweepStalePanelProfiles(prefix string)` | 每个进程每个 prefix 只跑一次（`sync.Map`），列 `os.TempDir()`，只删名字**恰为** `prefix-<pid>-<nanotime>` 且 mtime 早于 `stalePanelProfileAge` 的目录，错误全忽略 |
| `isPanelProfileLeafName(name, prefix)` | 名字形状判据：`prefix-` 之后两段都得是十进制整数；`wisp-33r1-panel-profile`（旧固定名）、`...-abc-123`、`...-<leaf>-second` 一律不匹配 |
| `stalePanelProfileAge = time.Hour` | 1 小时。本族夹具单发实测 2–4s（§3），并发那发的叶子最老也就分钟级 ⇒ 永不落入删除区 |

**保留的两条边界**（都是仓里已有的话，不是我发明的）：
1. **父目录仍 `os.TempDir()`** —— 不换 `t.TempDir()`（那会把根挪到 GOTMPDIR，牵动票 310/311/317 的读数，`68275983e` 的 helper 注释逐字写着）。
2. **测试结束不删** —— 三处都**没有** `t.Cleanup(os.RemoveAll)`：`panel_host_windows_test.go:701-705` 的原文讲清了 `t.TempDir` 的自动 RemoveAll 会和浏览器 profile-DB 退场赛跑、撞锁文件而红（那是退场，不是产品断言）。"不删"的代价改由**启动时清 >1h** 付，正控见 §6。

### 1.2 三处调用点

| # | 文件:行（改后） | before | after |
|---|---|---|---|
| 1 | `cmd/wisp/panel_host_windows_test.go:706` | `filepath.Join(os.TempDir(), "wisp-33r1-panel-profile")` | `panelProfileLeaf("wisp-33r1-panel-profile")` |
| 2 | `cmd/wisp/panel_host_windows_test.go:1027` | 同上（**同文件第二处**，`TestAC4FocusReturnToPriorWindowGap33r5`） | 同上 |
| 3 | `cmd/wisp/panel_assets_299_windows_test.go:164` | `filepath.Join(os.TempDir(), "wisp-299r1-panel-profile")` | `panelProfileLeaf("wisp-299r1-panel-profile")` |

改后形状（逐字，`grep -n 'panelProfileLeaf(' cmd/wisp/*.go`）：

```
cmd/wisp/panel_host_windows_test.go:706:	dataPath := panelProfileLeaf("wisp-33r1-panel-profile")
cmd/wisp/panel_host_windows_test.go:1027:	dataPath := panelProfileLeaf("wisp-33r1-panel-profile")
cmd/wisp/panel_assets_299_windows_test.go:164:	dataPath := panelProfileLeaf("wisp-299r1-panel-profile")
```

`git diff --stat`：`panel_assets_299_windows_test.go` +9/−5（含 import 收窄：`os`/`path/filepath` 不再用）、
`panel_host_windows_test.go` +109/−5（helper + 两处调用点 + 注释）。**零产品码、零 `scripts/**`、零 `internal/**`、零阈值/golden/契约**。

---

## 2. 门禁读数（全部本会话现跑）

| 门 | 命令 | 读数 |
|---|---|---|
| vet（winlive） | `go vet -tags winlive ./cmd/wisp` | **rc=0**，无输出 |
| gofumpt | `D:\work\base\gopath\bin\gofumpt.exe -l cmd/wisp` | **输出空**、rc=0（裸名不在 PATH，`$(go env GOPATH)/bin/gofumpt.exe` v0.12.0） |
| d22scan | `sh scripts/d22scan.sh` | **clean（rc=0）**，见 §8 |
| 编译（真跑用） | `go test -c -tags winlive -o … ./cmd/wisp` | rc=0，62,981,591 B；固定名变异件 `-overlay` 编译 **rc=0**，62,977,189 B（"突变必须能编译"） |

---

## 3. 串行基线（先证明"改前改后单发一个样"）

| 构建 | 选择 | rc | 墙钟 | PASS/FAIL | cold |
|---|---|---|---|---|---|
| de-share | `^(299|lifecycle)$` 一发 | 0 | 5.78s | 2/0 | 1147.318 ms |
| 固定名（overlay） | 同上 | 0 | 5.58s | 2/0 | 1322.720 ms |
| de-share | `^299$` 单发 | 0 | 3.73s | 1/0 | — |
| de-share | `^lifecycle$` 单发 | 0 | 2.95s | 1/0 | — |
| de-share | `^AC4$` 单发 | 0 | 2.73s | 1/0 | — |

⇒ 单发口径下改动前后一致（两枚用例一枚进程内跑，rc=0、冷启都在 D32 1500 ms 内）。

---

## 4. 并发倒逼（★派单要求的正控；含我拿不到的那一格）

探针：`.scratch/wisp/probes/314-deshare/concurrent.ps1` —— 两发 `Start-Process` **背靠背同时**起（`launchSkewMs` 记两发启动间隔，全部 |skew| ≤ 83 ms），
采样器每 250 ms 用 `Win32_Process` 读所有 `msedgewebview2.exe` 的 `--user-data-dir`，
按"父不在同一 tick 的浏览器 pid 集合里"判**顶层 host**，按 UDF 分组。完整 JSON 行＝`results-concurrent.jsonl`。

### 4.1 两发各建**一**窗（★两发 rc=0 的那一格）

| tag | 两发选择 | skew | A.rc／PASS | B.rc／PASS | 同一 UDF 上最大顶层 host | 共存 tick | 两发 UDF 路径不同？ |
|---|---|---|---|---|---|---|---|
| `deshare-299-r1` | `^299$` | −80.9 ms | **0／1** | **0／1** | **1** | **0** | **是** |
| `deshare-299-r2` | `^299$` | −71.0 ms | **0／1** | **0／1** | **1** | **0** | **是** |
| `deshare-ac4-r1` | `^AC4$` | −57.8 ms | **0／1** | **0／1** | **1** | **0** | **是** |

逐发 UDF 路径（`results-concurrent.jsonl` 的 `a.udfs`/`b.udfs`），例 `deshare-299-r1`：

```
A: C:\Users\swq\AppData\Local\Temp\wisp-299r1-panel-profile-44984-1794627799631995100\EBWebView
B: C:\Users\swq\AppData\Local\Temp\wisp-299r1-panel-profile-29692-1794627799628386200\EBWebView
```

⇒ **同一 tick 同一 UDF 没有两个顶层 host（0 tick），且两发 UDF 路径不同（两条都成立）。**

### 4.2 两发各建**两**窗（★派单字面形状；**两发 rc=0 拿不到**）

| tag | 两发选择 | A.rc／PASS | B.rc／PASS | 死法 | 夹读数 |
|---|---|---|---|---|---|
| `deshare-r1` | `^(lifecycle｜299)$` | 0／2 | **−1073740771（0xC000041D）**／1 | 原生异常，无 `panic:`、无 `--- FAIL:`，死在 `=== RUN TestPanelHostRealWindowHopAndLifecycle` 之后 | maxHosts=1、共存 0、UDF 不同 |
| `deshare-r2` | 同上 | **−1073741819（0xC0000005）**／1 | 1（`cold 1641.1 ms exceeds D32 1500`） | A 原生异常；B 冷启预算红 | maxHosts=1、共存 0、UDF 不同 |
| `deshare-ac4-299-r1` | `^(AC4｜299)$` | 0／2 | **−1073741819**／1 | 原生异常，无 panic 行 | maxHosts=1、共存 0、UDF 不同 |
| `deshare-ac4-299-r2` | 同上 | 0／2 | **−1073741819**／1 | 同上 | maxHosts=1、共存 0、UDF 不同 |
| `deshare-lifecycle-r1` | `^lifecycle$` | 1（`cold 1607.9 ms exceeds 1500`） | 0／1 | 只有冷启预算红 | maxHosts=1、共存 0、UDF 不同 |
| `deshare-split-r1` | A=`^lifecycle$`、B=`^299$` | 1（`cold 1864.3 ms exceeds 1500`） | 0／1 | 同上 | maxHosts=1、共存 0、UDF 不同 |

**4/4 对"两发各两窗"都有一发原生消失**；`0xC000041D`/`0xC0000005` 都是原生异常码，**不是** `os.Exit(n)`、**不是** `-test.timeout` panic（日志里 `panic:` 计 0）。
另：`lifecycle` 这枚用例本身带 D32 冷启 1500 ms 预算，**两发并发时实测 1607.9／1641.1／1864.3 ms**（单发 1147／1322 ms）⇒ 并发下它必然踩预算红；这是单发 SLO 的读数，不是夹具共享的读数。`AC4` 的前台尺也读机器级状态：`deshare-ac4-r2` 两发同时红在同一句 `the panel did not take the foreground on Show: foreground 0x1b8152c`（两发读到同一个别人的前台窗）。

---

## 5. 反证：三处临时改回固定名，同一形状再跑

**做法**：不动跟踪文件，`-overlay` 换成改回固定名的整份副本（`overlay/`），浏览器看到的仍是真 `cmd/wisp` 路径；
副本里三行逐字回到：

```
overlay/panel_host_fixed.go:711:	dataPath := filepath.Join(os.TempDir(), "wisp-33r1-panel-profile")
overlay/panel_host_fixed.go:1027:	dataPath := filepath.Join(os.TempDir(), "wisp-33r1-panel-profile")
overlay/panel_assets_299_fixed.go:166:	dataPath := filepath.Join(os.TempDir(), "wisp-299r1-panel-profile")
```

`go test -c -tags winlive -overlay overlay/overlay-fixed-names.json …` → **rc=0**（能编译的突变），产物 `wisp-fixednames.test.exe`。

### 5.1 同一形状（两发各一窗），固定名

| tag | 选择 | A.rc／PASS | B.rc／PASS | 同一 UDF 上最大顶层 host | 共存 tick | 两发 UDF 相同？ |
|---|---|---|---|---|---|---|
| `fixed-ac4-r1` | `^AC4$` | 0／1 | 0／1 | **2** | **2** | **是** |
| `fixed-ac4-r2` | `^AC4$` | 0／1 | 0／1 | **2** | **4** | **是** |
| `fixed-299-r1` | `^299$` | 0／1 | 0／1 | **2** | **1** | **是** |
| `fixed-299-r2` | `^299$` | 0／1 | 0／1 | 1 | 0 | 是（本次没抓到共存那一瞬） |

逐字（`fixed-ac4-r2`，**入仓的 overlay 副本重编后的那一发**；A.pid=25976、B.pid=5880，与两个 parent 一一对上）：

```
tick=1 udf=C:\Users\swq\AppData\Local\Temp\wisp-33r1-panel-profile\EBWebView topLevelHosts=2 parents=25976,5880 browsers=52452,35892
tick=2 udf=C:\Users\swq\AppData\Local\Temp\wisp-33r1-panel-profile\EBWebView topLevelHosts=2 parents=25976,5880 browsers=52452,35892
tick=3 udf=C:\Users\swq\AppData\Local\Temp\wisp-33r1-panel-profile\EBWebView topLevelHosts=2 parents=25976,5880 browsers=52452,35892
tick=4 udf=C:\Users\swq\AppData\Local\Temp\wisp-33r1-panel-profile\EBWebView topLevelHosts=2 parents=25976,5880 browsers=52452,35892
A.udf: C:\Users\swq\AppData\Local\Temp\wisp-33r1-panel-profile\EBWebView
B.udf: C:\Users\swq\AppData\Local\Temp\wisp-33r1-panel-profile\EBWebView
```

（≈1 s 连续 4 tick；两发自己也把自己命名的 UDF 路径读成**同一条**。）

`fixed-299-r1` 同形：`tick=0 udf=…\wisp-299r1-panel-profile\EBWebView topLevelHosts=2 parents=47844,14500 browsers=48548,41156`（其中 47844 = 该发 B 的 pid；14500 这一枚该发未记到 `testPidsOnBox`，只照录）。

### 5.2 反证结论（照派单要求逐字）

- **固定名 ⇒ 同一 UDF 上两个顶层 host 并存：坐实**（`fixed-ac4-r1` 2 tick、`fixed-ac4-r2` 4 tick ≈1 s、`fixed-299-r1` 1 tick）；**de-share 三对 0 tick**，且两发 UDF 路径不同。
- **但"致命红"我造不出来**：`fixed-ac4-r1`／`fixed-ac4-r2`／`fixed-299-r1`／`fixed-299-r2` 两发**全部 rc=0**、无 `panic:`、无 `boot failed`、无 `WebView2 Embed failed`。**这一形我造不出致命红**。
- 唯一在固定名侧出现的死亡是 **"两发各两窗"那一形**（`fixed-ac4-299-r1`：B `−1073741819`），而**同一形在 de-share 树上也死一发、同一退出码** ⇒ 该死亡不随夹名改变，不能拿它当固定名的罪证，也不能拿 de-share 当它的解药。

---

## 6. 清理正控（"只删 >1h 同族旧夹"）

探针 `sweep-check.ps1`：在 `$env:TEMP` 里种 5 枚死物再跑一发 `^299$`（de-share 构建）：

| 死物 | 名字 | mtime | 期望 | 实测 |
|---|---|---|---|---|
| 旧同族叶 | `wisp-299r1-panel-profile-4242-<ms>` | −2h | 删 | **删** |
| 新同族叶（=并发那发的替身） | `wisp-299r1-panel-profile-4343-<ms>` | now | 留 | **留** |
| 别的 prefix | `wisp-33r1-panel-profile-4444-<ms>` | −2h | 留 | **留** |
| 旧固定名 | `wisp-299r1-panel-profile` | −2h | 留 | **留** |
| 形近 | `wisp-299r1-panel-profile-abc-<ms>` | −2h | 留 | **留** |

`results-sweep.jsonl`：`rc=0`、`elapsed_s=3.57`、`pass=true`，`before` 全 true → `after={"oldLeaf":false,"freshLeaf":true,"otherPrefix":true,"fixedLegacy":true,"lookalike":true}`。
⇒ **不误删并发那发的夹**（新叶留）、**不漏清旧夹**（旧叶删），且不碰别的 prefix／旧固定名／形近名。

---

## 7. 还原核（★派单一句话：SHA256 同串 ＋ `git diff` 0 件）

反证走 `-overlay`，跟踪文件**一个字都没动过**，故"还原"是恒等式而非动作；仍给读数：

| 文件 | SHA256（改动落盘后／全部探针跑完之后） |
|---|---|
| `cmd/wisp/panel_host_windows_test.go` | `C080903DCCF029C47250CA3148D85DA0AA7AF31517578F2461EF96A20A488D7E`（两次取样**同串**） |
| `cmd/wisp/panel_assets_299_windows_test.go` | `42735AA56CFA314ADF9F2C425FA3F74C4606F8AF582B384784EC4D43D102760B`（同上） |

`git diff --stat -- cmd/wisp/…` 恒为上面那两枚（+113/−5），**没有第 3 枚文件**；`git status --porcelain -- cmd/wisp/` 只有这两枚 ` M`。
"临时改回固定名"的产物全部在 `.scratch/wisp/probes/314-deshare/overlay/`（副本 2 枚 + overlay JSON 1 枚），随本笔入仓备核。

---

## 8. 边界与诚实记录

1. **派单字面形状（两发各跑两枚真建窗用例、两发 rc=0）我拿不到**：该形 4/4 对里必有一发原生消失，且固定名构建同形也死一发（§4.2、§5.2）。我拿到 rc=0 的并发形状是"两发各建一窗"（3/3 对、6/6 进程）。**不编、不修饰、不把 299-only 的绿说成"两枚用例都跑过"**。
2. 该死亡我只记读数、**不做归因**（派单：⛔ 不要再归因）。可核事实三条：无 `panic:`、无 `--- FAIL:`、退出码是原生异常（`0xC0000005`／`0xC000041D`）；发生前最后一行为 `=== RUN <第二枚用例>`；同样出现在固定名构建。另注：票面 Progress log `2026-11-14 09:5x` 记过"两发整包并发"下的 `rc=139`／90 分钟挂起，本件是同一族读数（原生异常），但**退出码不同、不是同一对数**。
3. `lifecycle` 这枚用例自带 D32 冷启 1500 ms 预算，**并发下必红**（1607.9/1641.1/1864.3 ms vs 单发 1147/1322 ms）；`AC4` 的前台尺读机器级前台窗，并发下也会红（`deshare-ac4-r2` 两发同时红、读的是同一个别人的窗）。两者都**不是**夹具夹读数，我没有放宽它们任何一条断言（阈值/golden/契约零改动）。
4. 机器负载：全程 CPU 91–100%（`Get-CimInstance Win32_Processor` 读数；Weixin／3×DeepSeek Harness／magpie 常驻），桌面未锁屏、有人在使用 ⇒ 这是"热桌面"读数，⛔ 不是"安静桌面"读数。
5. 收尾清了自己造的 1 枚孤儿浏览器（browser 16996，parent 20524＝`deshare-r2` 崩溃发；已退出进程的孩子），未杀任何产品进程。
6. 未跑整包（只 `-test.run` 定向：单枚/两枚）；未 push；未改 `scripts/**`／`internal/**`／`frontend/**`／阈值／golden／契约。

---

## 9. 复跑（逐字）

```powershell
# 1) 两个构建（de-share 用工作树；固定名用 overlay）
go test -c -tags winlive -o D:\tmp\wisp\20261114\314-deshare\bin\wisp-deshare.test.exe ./cmd/wisp
go test -c -tags winlive -overlay .scratch\wisp\probes\314-deshare\overlay\overlay-fixed-names.json `
  -o D:\tmp\wisp\20261114\314-deshare\bin\wisp-fixednames.test.exe ./cmd/wisp

# 2) 正控（两发各一窗，两发 rc=0；UDF 不共存）
pwsh -NoProfile -File .scratch\wisp\probes\314-deshare\concurrent.ps1 `
  -Bin D:\tmp\wisp\20261114\314-deshare\bin\wisp-deshare.test.exe -Tag deshare-299-r3 -TestRun '^TestPanelPageFetchesItsOwnAssets299$'

# 3) 反证（同形状、固定名）—— 期望 coexist>=1、maxTopLevelHostsOnOneUdf=2、两发 UDF 相同
pwsh -NoProfile -File .scratch\wisp\probes\314-deshare\concurrent.ps1 `
  -Bin D:\tmp\wisp\20261114\314-deshare\bin\wisp-fixednames.test.exe -Tag fixed-299-r4 -TestRun '^TestPanelPageFetchesItsOwnAssets299$'

# 4) 清理正控
pwsh -NoProfile -File .scratch\wisp\probes\314-deshare\sweep-check.ps1 `
  -Bin D:\tmp\wisp\20261114\314-deshare\bin\wisp-deshare.test.exe

# 5) 门禁
go vet -tags winlive ./cmd/wisp
D:\work\base\gopath\bin\gofumpt.exe -l cmd/wisp
sh scripts/d22scan.sh
```
