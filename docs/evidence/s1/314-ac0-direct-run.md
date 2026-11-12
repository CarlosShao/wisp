# 票 314 AC#0 — 直跑测试二进制、取进程退出码、数浏览器进程增量（编排者现跑）

- 锚：`git rev-parse --short HEAD` = **`c38898db4`**（`dev`）
- 时刻：**2026-11-12 09:38–09:40 +08**
- 环境：`PATH` 含 `third_party/sherpa-onnx` 与 `build`（缺则 `0xC0000135`）
- 腿性质：**编排者亲自现跑**（不走子代理 —— 本轮两条 `group/group-review-and-audit` 腿均返回**空交件**，见 §5）

---

## 1. ★先纠正 `A1746` 的一条猜测（本条最重要）

`A1746` 据"78 个 `msedgewebview2` 存活 + 39 个 `wisp-33r5-*` profile 夹 39/39 非空"推断：

> **"`t.Cleanup(os.RemoveAll)` 每次都失败 ⇒ 夹具不收尾浏览器子进程"**

**现跑实测推翻了"夹具不收尾子进程"这半句**：

| 测试 | 是否开窗 | rc | 退出码 | 耗时 | **浏览器进程增量** |
|---|---|---|---|---|---|
| `TestPanelThreadIsSTAAndExitsCleanly` | 否 | 0 | `0x00000000` | 25s | **0**（66 → 66） |
| `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` | **是** | 0 | `0x00000000` | **10s**（测试本身 3.13s） | **0**（66 → 66） |

- **开窗那一发**的测试日志逐字：
  ```
  panel thread ending  why="the library pump returned"  window_opened=true
  panel thread exited cleanly  shows=1 toggles=0 disposals=0
  --- PASS: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (3.13s)
  ```
- ⇒ **面板线程干净退出、窗口正常开、浏览器子进程数不增**。
- ⇒ **`A1746` 的"夹具不收尾子进程"猜测是错的**。真正的原因是别的（见 §3 的剩余候选）。

## 2. 两条单发读数（可复跑）

```powershell
$env:PATH = 'D:\work\workspace\projects plans\Wisp\third_party\sherpa-onnx;D:\work\workspace\projects plans\Wisp\build;' + $env:PATH
$before = (Get-Process -Name msedgewebview2 -ErrorAction SilentlyContinue | Measure-Object).Count
go test ./cmd/wisp -run '^TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe$' -count=1 -v
Start-Sleep -Seconds 8
$after = (Get-Process -Name msedgewebview2 -ErrorAction SilentlyContinue | Measure-Object).Count
"DELTA = $($after - $before)"
```

**读数**：`before=66`、`after=66`、`DELTA=0`、`rc=0`。

## 3. 那 78 个泄漏进程从哪来？（剩余候选，⛔ 未定论）

既然"正常路径"不泄漏，候选只剩：

| # | 候选 | 支持 | 反对 |
|---|---|---|---|
| 甲 | **崩溃/被杀路径**（进程被 `TerminateProcess` 或测试超时 panic 时 `t.Cleanup` 不跑） | `scripts/portable-tests.sh:104-114` 逐字写了这个形状；`A1746` 有 39 个非空 profile 夹 | 需在崩溃那一发上量 |
| 乙 | **17.5 天前那 6 个**（2026-10-25 09:35 启动）来自更早的、已被修掉的代码 | 时间对得上 | 需 git 历史确认 |
| 丙 | **不是本仓测试产生的**（机主自己的面板/浏览器） | 66 个进程里可能混有他正在用的 | 未区分 |
| 丁 | **`-timeout` 触发的那一发**（`A1746` 的 ac3 在 392.7s 处死） | 时间对得上 | **未复现** |

⇒ **`A1746`/`A1747` 里"单发数增量即可定论"的简化，只在"正常路径"上成立；正常路径实测不泄漏。** 要抓真凶必须在**崩溃路径**上量。

## 4. 对票 314 的含义

- **AC#0 的"直跑二进制 + 取进程退出码"这一半**：本件用 `Start-Process -PassThru` + `$p.ExitCode` 做了两发，**均 `0x00000000`**（对照 `A1746` 记录的 `-1073741819` = `0xC0000005`）。
- **AC#0 的"四次死法分类"**：**两发都没死** ⇒ 死亡是**条件性**的，不是每次必死。这与 `A1746` 的"330 是 616 的精确位置前缀 ⇒ 早停位置确定、但是否发生有条件"一致。
- ⇒ **票 314 的下一步不该是"修夹具"**（正常路径无泄漏），而应是**在崩溃路径上取证**（候选甲/丁）。

## 5. ★如实登记：两条子代理腿都是空交件

| 腿 | 路由 | 返回 | 盘上核（三把尺） |
|---|---|---|---|
| `346-a2` | `group/group-review-and-audit` max | `"I'll start with the environment setup (Step 0)."` | **无 commit、无件** |
| `314-rep` | `group/group-review-and-audit` max | **空字符串** | **无 commit、无件** |

- **三把尺**：① `git log` 最新仍是 `c38898db4`（我自己的台账提交）② 两件均 `MISSING` ③ `go/wisp` 进程 0。
- ⇒ **这是本会话第 5 次证实**：`group/group-review-and-audit` **作为执行者返回空**（`323-a1` 曾 4 次空返回）。
- ⇒ **处置**：**测量/执行类腿改走 `auto/glm-5.3-flash` high**；**group 组只用于"读已存在材料出判语"的纯评审**。
- ⚠ 本节是对 `A1754` 派单的**事后更正**：那两条腿派错了路由。

## 6. 撤销口令

删本件即复原（本件零产码，未改任何 `.go`）。
