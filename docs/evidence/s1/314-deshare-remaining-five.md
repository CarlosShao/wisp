# 票 314 — 补完固定名 WebView2 user-data 夹 de-share（剩 5 处；写腿证据件）

> 角色：实现者（⛔ 审计者；⛔ 未翻任何 AC 框）。时间 `2026-11-14 12:0x–12:2x +08`。
> 树 = `dev`，起点 HEAD `48457eafe`（编排者 12:0x 那条裁定提交；本腿开工时 `git log` 首行已不是上一腿的 `edccb3449`）。
> 写面（逐字照派单）：5 个 `cmd/wisp/*_test.go`、`.scratch/wisp/probes/314-deshare2/**`、本件、票 314 Progress log。
> 读数根：`.scratch/wisp/probes/314-deshare2/{results-*.jsonl,logs/**,overlay/**}`。
> 测试二进制（2×约 63MB，**不入仓**）在 `D:\tmp\wisp\20261114\314-deshare2\bin\`。

---

## 0. 一句话结论

**5 处固定名夹全部换成既有 `panelProfileLeaf("<原固定名>")`（唯一叶＝pid＋纳秒），零新 helper、零产品码。**
**默认跑那处（`panel_unbound_253r6`，tag 只有 `windows`）的正控拿到**：两发背靠背同时起（|skew| 53.2／77.3 ms，两轮），**两发都 `rc=0`、各 1 PASS**；同一 tick 同一 UDF 上**最多 1 个顶层 host、共存 0 tick**，两发 UDF 路径**不同**且都含 pid＋纳秒叶。
**反证拿到**：同形用 `-overlay` 把该处临时改回固定名重编（编译 rc=0），**3 轮里 2 轮复现同一 UDF 两个顶层 host**（各 1 tick，parents 恰是两发 pid）；第 3 轮 skew 225.6 ms 更大、没抓到共存那一瞬，如实记。
**门禁**：`go vet -tags winlive ./cmd/wisp` rc=0；`gofumpt -l cmd/wisp` 空；`sh scripts/d22scan.sh` **clean rc=0**；五枚定向用例最终同轮 **RUN=1/2 PASS=全部 FAIL=0 SKIP=0**（读数见 §2/§3）。
没有任何一处是"故意要固定名"——逐处核过注释，理由见 §1.3。

---

## 1. 改了什么

### 1.1 五处调用点（改后行号）

| # | 文件:行（改后） | before | after | build tag |
|---|---|---|---|---|
| 1 | `cmd/wisp/panel_unbound_253r6_windows_test.go:71` | `filepath.Join(os.TempDir(), "wisp-253r6-panel-profile")` | `panelProfileLeaf("wisp-253r6-panel-profile")` | `windows`（**默认跑**） |
| 2 | `cmd/wisp/panel_host_windows_budget_winlive_test.go:55` | `filepath.Join(os.TempDir(), "wisp-33r1-panel-profile")` | `panelProfileLeaf("wisp-33r1-panel-profile")` | `windows && winlive` |
| 3 | `cmd/wisp/panel_host_windows_live_test.go:36` | 同上 | 同上 | `windows && winlive` |
| 4 | `cmd/wisp/panel_reshow_255r8_winlive_test.go:125` | `filepath.Join(os.TempDir(), "wisp-255r8-panel-profile")` | `panelProfileLeaf("wisp-255r8-panel-profile")` | `windows && winlive` |
| 5 | `cmd/wisp/panel_transport_live_35v2_windows_test.go:185` | `filepath.Join(os.TempDir(), "wisp-35v2-panel-profile")` | `panelProfileLeaf("wisp-35v2-panel-profile")` | `windows && winlive` |

`git diff --stat`（5 件，共 +5/−11）：

```
 cmd/wisp/panel_host_windows_budget_winlive_test.go | 4 +---
 cmd/wisp/panel_host_windows_live_test.go           | 3 +--
 cmd/wisp/panel_reshow_255r8_winlive_test.go        | 2 +-
 cmd/wisp/panel_transport_live_35v2_windows_test.go | 3 +--
 cmd/wisp/panel_unbound_253r6_windows_test.go       | 4 +---
```

### 1.2 import 收窄（只为去掉不再使用的包；未动其他）

| 文件 | 去掉 |
|---|---|
| `panel_unbound_253r6_windows_test.go` | `os`、`path/filepath`（该文件只在被替换那行用过） |
| `panel_host_windows_budget_winlive_test.go` | `os`、`path/filepath`（同上） |
| `panel_host_windows_live_test.go` | `path/filepath`（`os` 仍用于 `os.Getpid()`） |
| `panel_transport_live_35v2_windows_test.go` | `path/filepath`（`os` 仍用于 `os.Getpid()`/`os.Stderr`） |
| `panel_reshow_255r8_winlive_test.go` | 无（`os`/`filepath` 仍用于 `:90` 的 seed 写入） |

### 1.3 「有没有哪一处是故意要固定名？」——逐处核过，结论：**没有**

判据＝每条调用点的上下文注释里是否写了"固定名的必要性"（例如"跨发复用同一夹以复现某缺陷"）。逐处读数：

| 处 | 注释里关于夹/名的内容 | 判定 |
|---|---|---|
| `panel_unbound_253r6_windows_test.go:68-71` | `startPanelFor253r6` 的注释只讲"为什么要把 dispatch 传进来"，**没有一个字讲夹名** | 无理由，照改 |
| `panel_host_windows_budget_winlive_test.go:50-56` | 只讲 D32 预算为何搬进 winlive 档 | 无理由，照改 |
| `panel_host_windows_live_test.go:5-24` | 只讲 2s 退出条款为何搬进 winlive 档 | 无理由，照改 |
| `panel_reshow_255r8_winlive_test.go:118-123` | 唯一一条相关注释：**夹要放在系统 temp、不要 `t.TempDir()`**（TempDir 的 RemoveAll 会与 Edge profile 落盘赛跑，"EBWebView\\Default: The directory is not empty"）。它论证的是**根目录**，不是**固定名**；`panelProfileLeaf` 恰好保留 `os.TempDir()` 根 ⇒ 该理由**不构成**保留固定名的理由 | 无固定名理由，照改 |
| `panel_transport_live_35v2_windows_test.go:3-56` | 只讲 winlive 档归属与"为什么这不是被禁的 N6 形状"，未提夹名 | 无理由，照改 |

⇒ **5 处全部照改，没有一处按"故意固定名"豁免。** 上表即是派单要求的"逐处理由"（结论为"无理由"）。

---

## 2. 门禁读数（全部本会话现跑）

| 门 | 命令 | 读数 |
|---|---|---|
| vet（winlive） | `go vet -tags winlive ./cmd/wisp` | **rc=0**，无输出 |
| gofumpt | `D:\work\base\gopath\bin\gofumpt.exe -l cmd/wisp`（v0.12.0，裸名不在 PATH） | **输出空**、rc=0 |
| d22scan | `sh scripts/d22scan.sh` | **clean rc=0**、97.8s（`logs/d22scan-r2.out.txt` 末行 `d22scan: clean - no D22 ban violations`） |
| 编译（真跑用） | `go test -c -tags winlive -o … ./cmd/wisp` | rc=0，62,980,567 B |
| 编译（反证用） | `go test -c -tags winlive -overlay … ./cmd/wisp` | rc=0（"突变必须能编译"） |

**d22scan 一次假红已记**：第一次跑（`logs/d22scan.out.txt`，12:12）红在 `TestScannerSelfScanOfRealRepoIsGreen`：
`internal/observe/quota_44.go:223: [phantom-citation] comment cites repo path "docs/evidence/s1/44-quota-ledger-race.md" which does not exist on disk`。
现量该文件在 12:15 已由**并行的票 44 腿**落盘（`git status` 显示 `M internal/observe/quota_44.go` ＋ `?? docs/evidence/s1/44-quota-ledger-race.md`，不在本腿写面），随即重跑 **clean rc=0**。⇒ 该红是并行腿的瞬时状态，与本腿 5 处改动无关（本腿未碰 `internal/**`）。

---

## 3. 定向用例四数（最终同轮；⛔ 没跑整包）

同一条命令形：`go test -count=1 [-tags winlive] -run '<regex>' -v ./cmd/wisp`；日志 `logs/gate-final-*.out.txt`。

| 用例 | RUN | PASS | FAIL | SKIP | rc | 关键读数 |
|---|---|---|---|---|---|---|
| `Test253r6TheSilentDropBranchIsLoud`（默认档） | 1 | 1 | 0 | 0 | **0** | 三读全绿（REJECTED / RESOLVED / roster 拒绝＋INBOUND-DISPATCH 审计行） |
| `TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive` | 1 | 1 | 0 | 0 | **0** | tree 0 → peak → 0 |
| `TestPanelHostRealWindowHopAndLifecycle_WinLive` ＋ `TestPanelHostLatencyPercentilesAC2_WinLive` | 2 | 2 | 0 | 0 | **0** | hot 84.127 ms；cold 1157.165 ms（预算 1500/200） |
| `TestTicket255r8RealWindowFollowsTheConfigAcrossAReshow` | 1 | 1 | 0 | 0 | **0** | outer 640x400→916x439，client 900x400，同 HWND |
| `TestLive35v2PagePostMessageEnvelopeReachesTheGoDoor` | 1 | 1 | 0 | 0 | **0** | teardown 后 tree=0 |

### 3.1 负载敏感红——如实记（不是本改动造成的）

这台机器本会话 CPU 长期 91–100%（3×DeepSeek Harness＋微信＋其它 WebView2 常驻，机器级 `msedgewebview2` 基线 31 枚）。同一批用例在负载尖峰下红过，**且用改前（固定名）二进制同样红**，故这些红不归因于夹名；我没有放宽任何断言：

- **2s 退出条款**（`logs/exit2s-*.out.txt`，3 轮 A/B 交替）：deshare 红 2 轮（`now 7/8`）绿 1 轮；before 红 1 轮（`now 7`）绿 2 轮。红的那几轮树内 webview 在 Destroy 后 2s 仍为 7–8。
- **D32 冷/热预算**（`logs/gate-winlive-budget.out.txt`）：负载尖峰时 hot 207.7 ms > 200 红；同一构建在 `gate-winlive-budget-retry1/2.out.txt` 与最终轮分别读 142.1／44.0／84.1 ms 绿。A/B（`logs/budget-*.out.txt`）里 before 构建也读到过 cold 1352.7 ms（接近上限）、deshare 1483.7 ms，两个构建都会在尖峰下红。
- **35v2 teardown**（`logs/transport-*.out.txt`，2 轮 A/B）：round1 两个构建都 green（tree=0）；round2 两个构建都红（tree=8）。
- 这正是 `panel_host_windows_live_test.go:5-20` 自己写过的"负载效应，不是分母假象"；`docs/evidence/s1/314-deshare-fixed-profiles.md` 上一腿也记过并发下必红的那条 D32 冷启预算。**⛔ 我没有改阈值、没有加 retry、没有删断言。**

---

## 4. 正控：默认跑那处两发同时起（★派单要求）

探针 `.scratch/wisp/probes/314-deshare2/concurrent.ps1`（照 `probes/314-deshare/concurrent.ps1` 改：默认选择换成 `^Test253r6TheSilentDropBranchIsLoud$`，竞品计数用本腿二进制名）。
两发 `Start-Process` **背靠背同时起**；采样器每 250 ms 读所有 `msedgewebview2.exe` 的 `--user-data-dir`，按"父不在同一 tick 的浏览器 pid 集合里"判顶层 host，按 UDF 分组。原始行＝`results-concurrent.jsonl`。

### 4.1 de-share 构建（工作树现编）

| tag | skew | A.rc／PASS | B.rc／PASS | 同一 UDF 最大顶层 host | 共存 tick | 两发 UDF 不同？ |
|---|---|---|---|---|---|---|
| `deshare-253r6-r1` | **−53.2 ms** | **0／1** | **0／1** | **1** | **0** | **是** |
| `deshare-253r6-r2` | **−77.3 ms** | **0／1** | **0／1** | **1** | **0** | **是** |

逐发 UDF 路径（r1，逐字）：

```
A(pid 53740): C:\Users\swq\AppData\Local\Temp\wisp-253r6-panel-profile-53740-1794629282092589700\EBWebView
B(pid 55088): C:\Users\swq\AppData\Local\Temp\wisp-253r6-panel-profile-55088-1794629282043751000\EBWebView
```

⇒ **路径含 pid＋纳秒叶、两发不同**；同一 UDF 上没有出现第二个顶层 host（`sameUdfCoexistence="none"`、`maxTopLevelHostsOnOneUdf=1`）。两发 `panic=0`、`bootFailed=0`、`embedFailed=0`、`accessViolation=0`。

### 4.2 反证：`-overlay` 把该处临时改回固定名（同形重跑）

**做法**：跟踪文件**一字未动**；`overlay/` 下 5 份副本逐字节等于改前提交 `48457eafe:cmd/wisp/<file>`（`git show` 经 cmd 重定向；字节数＝`git cat-file -s`，5/5 相等：13101／6364／3292／11435／12668），其中 `panel_unbound_253r6_windows_test.go:73` 逐字回到：

```
	dataPath := filepath.Join(os.TempDir(), "wisp-253r6-panel-profile")
```

`go test -c -tags winlive -overlay overlay/overlay-fixed-names.json …` → **rc=0**（能编译的突变），产物 `D:\tmp\…\wisp-before.test.exe`。

| tag | skew | A.rc／PASS | B.rc／PASS | 最大顶层 host | 共存 tick | 两发 UDF 相同？ |
|---|---|---|---|---|---|---|
| `fixed-253r6-r1` | −86.0 ms | 0／1 | 0／1 | **2** | **1** | **是** |
| `fixed-253r6-r2` | −51.1 ms | 0／1 | 0／1 | **2** | **1** | **是** |
| `fixed-253r6-r3` | −225.6 ms | 0／1 | 0／1 | 1 | 0 | 是（本轮没抓到共存那一瞬） |

**共用 UDF 的 tick 逐字抄**（r2；parents 恰是两发 pid）：

```
tick=0 udf=C:\Users\swq\AppData\Local\Temp\wisp-253r6-panel-profile\EBWebView topLevelHosts=2 parents=53088,18276 browsers=65520,41256
```

r1 同形：

```
tick=0 udf=C:\Users\swq\AppData\Local\Temp\wisp-253r6-panel-profile\EBWebView topLevelHosts=2 parents=34636,51000 browsers=45520,45332
```

⇒ 固定名侧 **3 轮里 2 轮**读到"同一 UDF 同一 tick 两个顶层 host"（各 1 tick）；第 3 轮 skew 更大（225.6 ms）、两窗重叠窗口没被 250 ms 采样抓到，`maxTopLevelHostsOnOneUdf=1`——**这一轮我造不出共存，如实记，不把 2/3 说成 3/3**。
另：反证两发**都 `rc=0`**，即固定名下"共用 UDF"当场可读，但**没有造出致命红**（与上一腿同结论：固定名不等同于必死）。

---

## 5. 还原核（SHA256 同串 ＋ 无残留）

反证走 `-overlay`，5 个跟踪文件**从未被写回固定名**，故"还原"是恒等式而非动作；仍给读数（全部探针跑完后取样）：

| 文件 | SHA256 |
|---|---|
| `cmd/wisp/panel_unbound_253r6_windows_test.go` | `C9A14528D6E64CE98E13EDEF5F1064585B975234DF22BD384E78ADCAFB240551` |
| `cmd/wisp/panel_host_windows_budget_winlive_test.go` | `8174EC64F1A75B84EFE471E26F20A93978A3DC49013DE14D98457D9B2F58C21D` |
| `cmd/wisp/panel_host_windows_live_test.go` | `71C9F7888F2266895BA452AF3D2D735DA4DAEFC67CF211193B336CE4A9894A6D` |
| `cmd/wisp/panel_reshow_255r8_winlive_test.go` | `96A38B69AB60DAE1D9F0C22CAC303CE3620DD1D5D9D2BD79879062275A733086` |
| `cmd/wisp/panel_transport_live_35v2_windows_test.go` | `CA47814127197560D72156F6037EE5E7D0A1C6CE63C217B9DBF2DC046014D7D8` |

- 跟踪测试件里 `filepath.Join(os.TempDir(), "wisp-*panel-profile")` 现量 **0 条**（`Select-String -Path cmd\wisp\*_test.go` 计数 0）。
- `git status --porcelain` 对这 5 件只有 5 行 ` M`；`git diff` 只有 §1.1 那 5 件（+5/−11），**没有第 6 件测试件**。
- "临时改回固定名"的产物全部在 `probes/314-deshare2/overlay/`（5 份字节级副本＋overlay JSON），随本笔入仓备核。

---

## 6. 边界与诚实记录

1. **没有一处按"故意固定名"豁免**；逐处理由见 §1.3。唯一相关注释（`panel_reshow_255r8_winlive_test.go:118-123`）论证的是"根目录别用 `t.TempDir()`"，`panelProfileLeaf` 保留 `os.TempDir()` 根，故不构成保留固定名的理由。
2. **正控两轮都 rc=0**（A/B 各 1 PASS），但这是"每发一窗"的形状；本腿**没有**去复跑上一腿那个"两发各两窗必死一发"的形状（那是上一腿的读数与本票主题，不属本腿派单）。
3. **反证只复现 2/3 轮**；第 3 轮如实记"没抓到共存那一瞬"。没有把 2/3 写成 3/3。
4. **负载敏感红全部留档**（§3.1），且用改前二进制同形对照证明与本改动无关；**没有放宽任何断言、没有加 retry、没有删断言**。
5. 开工时 HEAD 已是 `48457eafe`（编排者 12:0x 的裁定提交，不是上一腿的 `edccb3449`）；本腿代码/证据提交＝`b1e1a738f`（只 commit 不 push）。该提交落地后数秒，并行的票 44 腿在它之上提交了 `6f23ce632`（`internal/observe/quota_44*`＋其证据件，**不是本腿的**）；此后 `HEAD` 不再是 `b1e1a738f`，但本件引用的"改前"一律指 `48457eafe`，不受影响。
6. ⛔ 未跑整包（只 `-test.run` 定向，最长一轮 2 枚用例）；⛔ 零产品码／`scripts/**`／`internal/**`／`frontend/**`／阈值／golden／契约改动；⛔ 未 push。
7. 并行腿（票 44）在 `internal/observe/quota_44.go` 与 `docs/evidence/s1/44-quota-ledger-race.md` 的未提交改动**不是本腿的**，未纳入本腿 commit（显式 pathspec）。

---

## 7. 复跑（逐字）

```powershell
# 0) 两个构建（de-share 用工作树；固定名用 overlay；均 -tags winlive 以便同时含两种 tag）
go test -c -count=1 -tags winlive -o D:\tmp\wisp\20261114\314-deshare2\bin\wisp-deshare.test.exe ./cmd/wisp
go test -c -count=1 -tags winlive -overlay .scratch\wisp\probes\314-deshare2\overlay\overlay-fixed-names.json `
  -o D:\tmp\wisp\20261114\314-deshare2\bin\wisp-before.test.exe ./cmd/wisp

# 1) 正控（两发同时起；期望 maxTopLevelHostsOnOneUdf=1、sameUdfCoexistenceTicks=0、udfPathsDistinctBetweenRuns=true）
pwsh -NoProfile -File .scratch\wisp\probes\314-deshare2\concurrent.ps1 `
  -Bin D:\tmp\wisp\20261114\314-deshare2\bin\wisp-deshare.test.exe -Tag deshare-253r6-r3 `
  -TestRun '^Test253r6TheSilentDropBranchIsLoud$'

# 2) 反证（同形、固定名；期望 maxTopLevelHostsOnOneUdf=2、sameUdfCoexistenceTicks>=1）
pwsh -NoProfile -File .scratch\wisp\probes\314-deshare2\concurrent.ps1 `
  -Bin D:\tmp\wisp\20261114\314-deshare2\bin\wisp-before.test.exe -Tag fixed-253r6-r4 `
  -TestRun '^Test253r6TheSilentDropBranchIsLoud$'

# 3) 门禁
go vet -tags winlive ./cmd/wisp
D:\work\base\gopath\bin\gofumpt.exe -l cmd/wisp
sh scripts/d22scan.sh

# 4) 五枚定向（四数）
go test -count=1 -run '^Test253r6TheSilentDropBranchIsLoud$' -v ./cmd/wisp
go test -count=1 -tags winlive -run '^TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive$' -v ./cmd/wisp
go test -count=1 -tags winlive -run '^(TestPanelHostRealWindowHopAndLifecycle_WinLive|TestPanelHostLatencyPercentilesAC2_WinLive)$' -v ./cmd/wisp
go test -count=1 -tags winlive -run '^TestTicket255r8RealWindowFollowsTheConfigAcrossAReshow$' -v ./cmd/wisp
go test -count=1 -tags winlive -run '^TestLive35v2PagePostMessageEnvelopeReachesTheGoDoor$' -v ./cmd/wisp
```

（`third_party/sherpa-onnx` 需在 PATH，否则该包测试二进制装载即 `0xc0000135` —— 票 98 的既有前提。）
