# 票 442 · ③-c 运行时判据落地 ＋ ① 形状尺落地冻结（裁定「甲＋落①」）

> 角色：**写腿（实现者）**，⛔ 审计者，不自审。上游＝票 442 尾部
> `## 编排者裁定（2026-11-14 15:5x）：走 甲＋落①`、派单
> `.scratch/wisp/probes/orch/20261114/txt/442-jiac-r1-dispatch.txt`、前腿
> `.scratch/wisp/probes/442/{option1-readings.md,option3-feasibility.md}`。
> 时间 2026-11-14 16:0x–16:1x +08。树＝仓工作树（HEAD `5209a38e7`）。
> 边界：本波唯一一枚可动 `cmd/wisp` 的写腿；写面＝`cmd/wisp/**`（③-c 接缝 diff ＋6/−0 ＋ 新测试件 ＋
> ① 落地后追加的冻结段）、本件、`.scratch/wisp/probes/442/**`、票 442 Progress log。
> ⛔ 未碰阈值／golden／契约／规格；⛔ 未放宽任何断言；⛔ 未加豁免；⛔ 未 push。

**第一行结论：① 落地（patch 57796 B → 工作树，post-image 同串 `9F976532…A4E2`，采后冻结）；
③-c 落地＝`NewPanelManager` 3 行可执行接缝（diff ＋6/−0，另 3 行是注释） ＋ 收口断言 `cmd/wisp/panel_data_path_distinct_442_windows_test.go:254`
（「同一构造点两次构造拿到同一个 dataPath」即红）；两侧正控＝把 `panel_unbound_253r6_windows_test.go:71`
改回固定名 ⇒ **运行时判据真红并逐字报出同一个值**、还原 ⇒ SHA 同串复绿；36 子例（35 形＋1 汇总）＋① 的
8 条控制全 PASS；四门禁 GREEN。**

---

## 0 交付物与 file:line（现跑行号）

| 件 | 是什么 | 读数 |
|---|---|---|
| `cmd/wisp/panel_host_windows.go:332-337` | **③-c 的接缝**：`NewPanelManager` 里 3 行可执行（`if`／调用／`}`）＋3 行注释＝diff ＋6/−0，把 `dataPath` 交给观察点 | SHA256 `0A63A827…839F`，66891 B |
| `cmd/wisp/panel_data_path_seam_windows.go:32` | 接缝的包级 `var panelDataPathSeam func(dataPath string)`（产码里恒 nil；只有 `_test.go` 赋值） | `A4A4102F…E4C37B3`，1851 B |
| `cmd/wisp/panel_data_path_distinct_442_windows_test.go:139` | ③-c 的收口测试 `Test442EveryDrivenConstructionPointHandsItsOwnDataPath` | `687ED474…783E`，11536 B |
| **`cmd/wisp/panel_data_path_distinct_442_windows_test.go:254`** | **③-c 的断言本身**：同一构造点两发拿到同一 dataPath ⇒ `t.Errorf` 逐字点名该点与值 | 同上 |
| `cmd/wisp/panel_profile_leaf_shape_windows_test.go:108-122` | ① 落地后追加的「冻结为次要启发式＋分工」头段（本波唯一往该件加的字） | `060C998C…90E6`，86017 B |
| `.scratch/wisp/probes/442/option1-ruler.patch` | ① 的落地物（前腿留） | `851AD655…F4F9D`，**57796 B**（与票面数字同串） |

> ⚠ 尺子 SHA 三个数分开记，别混：**HEAD 那枚** `6E54A0B3…23F2` → **① post-image**（patch 落完）`9F976532…A4E2`／84829 B
> → **本波追加冻结段后**（现盘上）`060C998C…90E6`／86017 B。三者的差只有「A1b/A1c 跨过程分析」与「冻结段 16 行」两块。

---

## 1 执行顺序（照派单，先①后③-c）

1. **先落 ①**（`git apply` 的等价物已在盘上，见 §2）→ 跑定向套件与四门禁；
2. **再做 ③-c**（接缝 diff ＋6/−0 ＋ 新测试件）；
3. **两侧正控**（突变 ⇒ 真红；还原 ⇒ SHA 同串复绿）；
4. 证据件（本件）＋ 票 Progress log ＋ 一笔显式 pathspec 的 commit。

---

## 2 ①：`option1-ruler.patch` 落地

### 2.1 落地读数

- **接手时盘上状态**：`cmd/wisp/panel_profile_leaf_shape_windows_test.go` SHA `9F976532…A4E2`，
  `git apply -R --check .scratch/wisp/probes/442/option1-ruler.patch` **rc=0**
  ⇒ 盘上正是该 patch 的 post-image（反向应用可前进＝正向已应用），不是"看起来像"。
  前腿 Progress log 写"已还原成 HEAD 那枚"，与本读数不符；**以现跑读数为准**（本件记这一笔，不改前腿的话）。
- patch 本体 `851AD655…F4F9D`，57796 B —— 与票面 57796 B 同数。
- 落地后既有子例：**35 形＋1 汇总＝36 条逐条 PASS**，加 ① 自带的 8 条控制，本文件 `--- PASS: .../` 实测
  **44 条**（36＋8），一条没少、一条没钝化（`run-06-restored-green.log` 数出 `SUBCASE_PASS=44`）。
- **豁免表 2 条逐字节未动**：把 `panelProfileLeafExemptions = []panelProfileLeafExemption{…}` 整块从
  `git show HEAD:` 版本与现盘版本各抽一遍，23 行、SHA256 **`1F6424E8…3FE6388` 两边同串**（`SAME=True`）。
- 尺子耗时：`Test314FixedPanelProfileLeafShapeIsRed` 现跑 **0.83 s**（扫 **276** 个源文件；HEAD 那轮是 274 个，
  ＋2 ＝ 本波新增的 `panel_data_path_seam_windows.go` 与 `panel_data_path_distinct_442_windows_test.go`）。

### 2.2 落完即冻结

① 的头注释追加段（`panel_profile_leaf_shape_windows_test.go:108-122`）逐条写明：主判据已换成运行时观测、
本尺**只读不扩张**、新形状归运行时判据或另立票、两边各盖各的洞且互不许冒领对方的覆盖。
**本波此后未再往尺子里加一条规则。**

---

## 3 ③-c：运行时判据

### 3.1 接缝（3 行可执行＋3 行注释，diff ＋6/−0）

```go
// Ticket 442's runtime observation seam (nil in production; see
// panel_data_path_seam_windows.go): the one send point a test reads the
// value the host was actually handed from, whatever spelled it.
if panelDataPathSeam != nil {
	panelDataPathSeam(dataPath)
}
```

- 位置：`cmd/wisp/panel_host_windows.go:332-337`（`NewPanelManager` 函数体第一段）。
- 值：`panelDataPathSeam` 是包级 func 值，**产码里恒 nil**，只有 `_test.go` 装记录器，`t.Cleanup` 还原；
  没有锁、没有分配、不留记录 ⇒ 空观察点时与补前代码不可区分（代价＝每次构造一次 nil 判断）。
- 为什么用 **func 值**而不是切片：观察顺序不是判据，且包级切片会把每条路径留到整发结束。

### 3.2 收口断言（**`panel_data_path_distinct_442_windows_test.go:254`**）

判据一句话：**这一发里，被驱动的每个构造点各拿两次值，两次必须不同、且不得为空**。

| 断言 | 行 | 红的样子 |
|---|---|---|
| 记录数与「点数×轮数」不等 ⇒ 测量缺口 | `:199-206` | 构造点不再走 `NewPanelManager`（判据对它变瞎）或接缝被拆 |
| 记录来自不同点数 ≠ 期望点数 | `:216-221` | 某个夹具开始从第二处构造，或两个点塌成一个 |
| 期望文件里必须**恰好一个**构造点 | `:224-232` | 判据读到的值不来自该文件自己的构造调用（标签是假的） |
| 每个点必须恰好 `rounds` 个值 | `:235-241` | 某轮没构造 ⇒ 测量缺口 |
| 空值 ⇒ 红 | `:248-252` | `""` ＝ 文档化的"交给 go-webview2 默认"，那是**一个共享目录** |
| **同点两次同值 ⇒ 红（判据本体）** | **`:254-263`** | 固定名回退：两发并发会落在同一个 WebView2 user-data 目录 |

`rounds=2` 是判据的要害：**一次构造分不出"每发唯一"与"固定名"**，而固定名正是全部题目。
判据读的是**构造点交给构造函数的那个值**，不问它怎么拼出来的 ⇒ 跨过程／跨文件／常量折叠／`Sprintf`／
`Replace`／`Join`／`Getenv` 全部自动覆盖 —— 这正是候选③相对源码形状尺的价值，也是票面 ③-a
（直接断言 `panelProfileLeaf` 两次不同＝同义反复）被拒的理由所在。

### 3.3 驱动了哪三个构造点（现跑点名）

```
ticket 442 runtime criterion: 3 construction point(s) x 2 round(s) = 6 data path(s),
all distinct within the run and none empty - points
[panel_unbound_253r6_windows_test.go:72 panel_resident_windows_test.go:74 panel_resident_windows.go:367]
```

1. **253r6 的窗体夹具**（`panel_unbound_253r6_windows_test.go:72`）——审计 A3 用来证"尺子没牙"的那条回退路径就在它上游 `:71`；
2. **33r5 的窗体夹具**（`panel_resident_windows_test.go:74`，`uniqueProfileLeaf` 家族）；
3. **产码唯一的构造点：常驻宿主装配根**（`panel_resident_windows.go:367`，值由 `:362` 的 `dataDir` 推出）——
   用票 323 的装配台件驱动（`newResidentPanelManager`，`panel_resident_assembly_323_test.go`），
   两个不同数据根 ⇒ 两个不同 dataPath；若把 `:362` 改成固定名/常量，这里两发同值即红。

> **⛔ 造窗**：三点都在"只建管理器＋线程、从不向 `bringUp` 要窗口"的层级上构造
> （319 的头注释 `:47-53` 把这一形写清楚了），所以本测试 **0.04 s**、无 WebView2 进程、无 profile 目录，
> 不占 253r6 文件头 `:34-38` 点名的那条"两窗接连 mid-create 死"的 flake 面。

### 3.4 覆盖范围（写清，⛔ 隐瞒）

- **盖**：上述三个构造点，**每次**运行都盖（`rounds=2`），且值来源与写法无关。
- **不盖（有意，逐条给理由）**：
  1. **319 holder 夹具**故意共用一个固定目录（`panel_holder_319_windows_test.go:47-53`），且从不建窗；
     它的正当性由**静态豁免表那一条**（含 `noWindow` 子句在尺子里被断言）持有 —— **本件不重抄一份豁免语义**，
     那正是候选 ③-b 被拒的代价；
  2. 交 `""`（"交给库默认"的文档化契约）或在**不建窗**的测试体里用 `t.TempDir()` 的那些点；
  3. 只能在该测试自己建窗时才走到的构造点（`panel_death_replay_319_winlive_test.go` 等）—— 它们留给静态尺。
- **分工**（两份头注释都写了）：本判据判**值**，尺子判**拼写**。尺子看不见"没跑到的那段码"里的固定名；
  本判据看不见"没被它驱动的那段码"。两边都不许冒领对方。

---

## 4 两侧正控（③-c 的牙）

### 4.1 突变：`panel_unbound_253r6_windows_test.go:71` 改回固定名 ⇒ 真红

突变体（审计 A3 逐字形）：`dataPath := filepath.Join(os.TempDir(), "wisp-253r6-panel-profile")`（并补 `os`／`path/filepath` 两行 import 使其可编译）。
**只跑** `-test.run '^Test442EveryDrivenConstructionPointHandsItsOwnDataPath'`（把静态尺挡在外面，证明红不是尺子给的）：

```
--- FAIL: Test442EveryDrivenConstructionPointHandsItsOwnDataPath (0.05s)
  panel_data_path_distinct_442_windows_test.go:254: construction point panel_unbound_253r6_windows_test.go:74
  handed the host the SAME data path in two constructions of one run: round 1 and round 2 both got
  "C:\\Users\\swq\\AppData\\Local\\Temp\\wisp-253r6-panel-profile". That is the fixed-name shape ticket 442
  exists to catch: ...
```

- 点名到**构造点 file:line**（`:74`＝突变后行号，夹具体内两行 import 使 `:71` 下移两行）与**同一个值逐字**。
- 同突变下静态尺也红（`run-05-mutant-static.log`：扫 276 件、命中 3、豁免 2、**实红 1**）——
  这一形两边都能抓，说明的是"重叠区"，不是"尺子够用"。

### 4.2 还原 ⇒ SHA 同串、复绿

| 读数 | 值 |
|---|---|
| 突变前／还原后 `panel_unbound_253r6_windows_test.go` SHA256 | `C9A14528D6E64CE98E13EDEF5F1064585B975234DF22BD384E78ADCAFB240551`（**两边同串**） |
| 还原后 `git status --porcelain -- 该件` | **0 行**（=与 HEAD 一致） |
| 还原后 `-test.run '^Test314|^Test442'` | `ok github.com/CarlosShao/wisp/cmd/wisp 1.00s`，4 个顶层用例全 PASS |

---

## 5 四门禁读数（现跑，全部在本波最终盘上）

| 门 | 命令 | 读数 |
|---|---|---|
| vet | `go vet -tags winlive ./cmd/wisp` | **rc=0**，零输出（`final2-vet.log` 空） |
| gofumpt | `gofumpt -l cmd/wisp` | **0 行**（`final2-gofumpt` 空） |
| d22scan | `scripts/d22scan.sh` | **rc=0**，`clean - no D22 ban violations`；分母 `cmd/ exam 282 Go files`、`internal/ 1010`（`final2-d22scan.log`） |
| path-length | `scripts/check-path-length-budget.sh` | **GREEN**，本波件入索引后 tracked=**20026**／over=69／roster=69／not-in-roster=0；最长 185、worst full 229、wall interval 0（`final3-path-budget-staged.log`；入索引前那一发是 tracked=20002，同一结论 `final2-path-budget.log`） |
| 定向套件 | `GOFLAGS='-tags=winlive' scripts/wisp-cli-staged-tests.sh -test.run '^Test314|^Test442'` | **rc=0**；`Test314FixedPanelProfileLeafShapeIsRed 0.83s`、`Test314PanelProfileShapeRulerRingsOnTheFixedNameShape`（44 条子例 PASS）、`Test442CrossFileHelperReturnReachesTheSink`、`Test442EveryDrivenConstructionPointHandsItsOwnDataPath 0.04s` |

> 载具：**只用** `scripts/wisp-cli-staged-tests.sh`（⛔ 裸 `go test ./cmd/wisp` 会因 System32 的 onnxruntime 假红
> `0xc0000135`）；⛔ 全程没有同时跑两发。

---

## 6 诚实边界（没做到／没测的，逐条写）

1. **没跑整包 `./cmd/wisp`**（90 分钟量级、且本波只写三处）：本件所有测试读数都来自 `-test.run` 定向发。
2. **③-c 只盖被驱动的三个点**（§3.4 逐条列了不盖的与理由）。它不是"扫全仓"，它不是尺子的替代品；
   票面要求的分工已写进两份头注释。
3. **接缝在产码里编译**（`panel_data_path_seam_windows.go` 无 build tag 之外的开关）：Go 没有"仅测试构建"的标签，
   所以做法是"恒 nil、只有测试赋值、零状态"；本件把这点写明，⛔ 说成"只在该测试构建里存在"。
4. **本波把尺子的扫描分母从 274 推到 276 件**（我加的两个文件），尺子耗时 0.73 s→0.83 s 的差也来自此；
   `d22scan` 的 `cmd/` 分母同步 280→282。
5. **本波未新增/未放宽任何豁免**；`panelProfileLeafExemptions` 23 行与 HEAD 同串（§2.1）。

---

## 7 撤销口令

- 「**票 442 撤**」＝删本票＋回退本波与该 patch 的 commit（新 commit 反改，⛔ reset/amend）。
- 「**442 撤 ③-c**」＝删 `panel_data_path_distinct_442_windows_test.go` ＋ `panel_data_path_seam_windows.go` ＋
  `panel_host_windows.go:332-337` 这 6 行（三处一起删才编得过）。
- 「**442 撤 ①**」＝`git apply -R .scratch/wisp/probes/442/option1-ruler.patch`（解冻时连冻结段一并处理）。
