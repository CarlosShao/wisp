# 450 — 落修：真窗两枚用例 profile 隔离 ＋ 载具终点锚（实现腿 r1，2026-11-21）

> **本件是实现件，不是验收件。** 判据全部可复跑（命令逐条在文内）。
> **锚点**：实现基线＝`2ed129e36`（诊断件 `450-death-localization-r1.md` 在 `8333de280`，其后两枚均为
> 台账/派单存档件、零代码）。⛔ 未 push。
> **写面**：`cmd/wisp/panel_host_windows_test.go` · `panel_host_windows_live_test.go` ·
> `panel_host_windows_budget_winlive_test.go`（只 `_test.go`）、`scripts/wisp-cli-staged-tests.sh`、
> 本件。共享树里**别的腿的未提交改动**（`.scratch/wisp/issues/34-*.md` 等）一字未碰、未入本腿 commit。

---

## §0 一句话结论

1. **profile 隔离落了**：4 处调用点各自不同 leaf（`-hop`／`-ac4`／`-live`／`-budget`），
   两枚真窗用例不再共用 `wisp-33r1-panel-profile` 家族。
2. **隔离自证＝3 发全活**（走载具、一次一发，`-test.run` 定向两枚）：3 发 0 死，逐发 rc=0。
3. **载具终点锚落了，且在真死亡上自己红过一次**：整包复跑仍是 **RC=139 静默死亡**，但这次载具
   打印逐字红句 `CARRIER_INCOMPLETE: died before the last test` 并以 139 判红——
   修之前的那发（诊断件 §1）同样的死法却**没有任何"没跑完"的信号**。
4. **如实报（⛔ 不许埋）**：整包死亡**换位置了**。先前致死的那对真窗用例在全包里**双双 PASS**；
   新死点是 `TestResyncOnShowDeliversFreshSnapshotFieldByField35g6`
   （`cmd/wisp/panel_resync_stateless_35g6_windows_test.go:133`，`--- PASS` 缺失，403 s，
   四数 RUN=589／PASS=451／FAIL=0／SKIP=1）。**本腿不改产码、不重开诊断**——这个新位置连同一条
   嫌疑线索（§5）留给下一条诊断腿。
5. **钉函数行为的尺一字未动**：`panel_profile_leaf_shape_windows_test.go` 零改动；定向复跑该尺
   rc=0 全绿（§4.4）。派单的 ⚠（"改调用点参数不该动它"）成立——不需要动它。

---

## §1 改动清单（全部可核 `git show <本 commit>`）

### §1.1 profile 隔离（4 处调用点，同函数 `panelProfileLeaf` 传不同前缀）

| 文件:行（改后） | 用例 | 改后 leaf |
|---|---|---|
| `cmd/wisp/panel_host_windows_test.go:720` | `TestPanelHostRealWindowHopAndLifecycle` | `panelProfileLeaf("wisp-33r1-panel-profile-hop")` |
| `cmd/wisp/panel_host_windows_test.go:1065` | 真窗焦点尺（诊断件 §3 最小集合第二枚；改前 func 声明 :1042，改后 :1051） | `panelProfileLeaf("wisp-33r1-panel-profile-ac4")` |
| `cmd/wisp/panel_host_windows_live_test.go:38` | `TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive`（winlive 层） | `panelProfileLeaf("wisp-33r1-panel-profile-live")` |
| `cmd/wisp/panel_host_windows_budget_winlive_test.go:64` | `TestPanelHostRealWindowHopAndLifecycle_WinLive`（winlive 层） | `panelProfileLeaf("wisp-33r1-panel-profile-budget")` |

- 只改 leaf 前缀参数与注释，**断言零改动**；lifetime 语义不变（仍然永不删除，靠既有 stale sweep）。
- sweep 按 prefix 记账（`sweptPanelProfilePrefixes`），新前缀各自每进程扫一次，行为不变；
  其文件头注释里过期的"three fixtures share two prefixes"更正为现状（注释，非断言）。

### §1.2 载具终点锚（`scripts/wisp-cli-staged-tests.sh`，+83 行含注释）

- **哨兵＝动态计算的整包尾**：跑前用**本载具刚构建的那枚二进制** `wisp.test.exe -test.list '.*'`
  取**最后一枚**——今天实测 790 枚顶层用例，末尾＝
  **`TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`**
  （读数存档：`.scratch/wisp/probes/orch/20261121/450-test-list-tail-check.txt`）。
  ⚠ 派单里举例的 `TestZZ315A1RootLoopFirstBridgeCall` **已经不是整包尾**（现排倒数第三），
  动态取尾使改名/新增尾用例不会把锚搞烂——空名册则载具自拒（GUARD，exit 1）。
- **武装判据＝"哨兵会不会真的跑"**，不是"有没有传 `-test.run`"：
  - 整包（无 `-test.run`/`-test.skip`）→ 武装；
  - 调用方 `-test.run P` → 用 `-test.list P` 查哨兵是否在运行集内，在则武装、不在则关（定向跑本来没有哨兵）；
  - 调用方 `-test.skip` 或 wave-8b 名册 fallback → 同理查哨兵是否被跳过，被跳过则关。
  关闭时打印一行 `end anchor OFF (<原因>)`，武装时打印哨兵名——两种状态都响亮可见。
- **判红**：整包跑完后日志里没有 `^=== RUN   <哨兵>$` ⇒ 打印逐字红句
  **`CARRIER_INCOMPLETE: died before the last test`**（stderr）＋四行解释，exit＝二进制原 rc
  （非零时保留，如 139），rc=0 时 exit 1——**二进制"ok"也照红**（正控实测，§4.2）。
  哨兵在 ⇒ 打印 `end anchor seen: ...` 后走原有 verdict 路径，exit 语义不变。
- 原有 guard（stage 拒绝、onnx 检查、名册空拒、GUARD B 结果行、FIXTURE-GAP）**全部未动**；
  例外一处：fallback 名册从"只存计数"改为"同时存名字"（`listed_names`），计数算法逐字等价，
  供哨兵成员判定用。

---

## §2 自证读数（全部走载具 `scripts/wisp-cli-staged-tests.sh`，一次一发串行；跑前 `Get-Process wisp.test` 均＝空）

### §2.1 隔离 3 发（`-test.run '^(TestPanelHostRealWindowHopAndLifecycle|<诊断件 §3 死亡点用例>)$'`）

| 发 | rc | 结果 | 读数文件 |
|---|---|---|---|
| 1 | **0** | **活**：两枚都 `--- PASS`（1.17s／1.10s） | `.scratch/wisp/probes/orch/20261121/450-shot1.log` |
| 2 | **0** | **活**（0.96s／0.95s） | `450-shot2.log` |
| 3 | **0** | **活**（1.24s／1.09s） | `450-shot3.log` |

**3 发 0 死 ⇒ 按派单判据写"隔离有效"**（口径＝那两枚的定向组合；对照诊断件基线：两枚一起
14 发 6 死）。每发载具均打印 `end anchor OFF (the sentinel is not in the caller's -test.run set)`
——负控一：定向跑不受锚影响，rc=0。

### §2.2 终点锚正控（哨兵名写错的临时副本，⛔ 不进仓）

- 副本：`D:\tmp\wisp\20261121\fix450\carrier-posctl-wrong-sentinel.sh`（仓外），三处外科手术式改动：
  ① repo root 写死为本仓路径（副本不在 `scripts/` 下）；② 哨兵名写成 **`TestNoSuchSentinel450`**；
  ③ 在武装打印前强制 `end_anchor_armed=1`（跑定向短臂用——哨兵名错了会被武装判据合法关闭，
  强制武装才走得到红臂）。其余逐字节同真载具。
- 跑法：`-test.run '^(TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask)$'`（单枚、0.19s、
  二进制 rc=0）⇒ 红臂全路径走通，**逐字红句**：

  ```
  CARRIER_INCOMPLETE: died before the last test
  ```

  载具 exit＝**1**（二进制 rc=0 的情形）。读数：`450-posctl.log`。

### §2.3 负控（整包、真载具、锚武装）

真载具整包一发：锚武装（打印哨兵名）、跑到死点前载具判红（§3），exit＝**139**（二进制原 rc
保留）。真载具定向跑（§2.1 三发＋§4.4 尺）rc=0，锚按判据关闭——两条负控都过。

### §2.4 整包复跑读数（如实，对照诊断件）

| 读数 | 诊断件主发（改前） | 本腿整包（改后） |
|---|---|---|
| rc | 139 | **139** |
| 测试段 | 327.00s | **403.00s** |
| 四数（`^=== RUN`／`^--- PASS`／`^--- FAIL`／`^--- SKIP`） | 412／329／0／1 | **589／451／0／1** |
| 最后一个 `=== RUN` | 真窗焦点尺（诊断件 §1） | **`TestResyncOnShowDeliversFreshSnapshotFieldByField35g6`**（其 `--- PASS` 缺失） |
| 先前致死的两枚 | 第二枚死（`--- PASS` 缺失） | **双双 `--- PASS`**（1.17s／1.66s，在全包内） |

判读：**那对真窗用例的死亡被隔离掉了**（3 发定向 0 死＋全包内双过），但**整包静默死亡仍在**，
位置移动——这不是"修完了"，是"修掉了测量到的那一对，露出了下一个"。读数文件：`450-fullpkg.log`。

---

## §3 门禁（全部现跑）

| 门禁 | 结果 |
|---|---|
| `go build ./...` | rc=**0** |
| `go vet -tags winlive ./cmd/wisp`（两枚 winlive 夹具编译面） | rc=**0** |
| `gofmt -l cmd/wisp` | 空（rc=0） |
| `gofumpt -l cmd/wisp`（经 `go run mvdan.cc/gofumpt@latest`） | 空（rc=0） |
| `bash -n scripts/wisp-cli-staged-tests.sh` | rc=**0** |
| `sh scripts/d22scan.sh` | **0 findings**；末行逐字：`d22scan.sh: OK - control rc=0, scan rc=0`（读数 `450-d22scan.log`） |
| ⛔ `-skip` 掩盖 | 未使用（fallback 名册机制本身未动） |
| ⛔ push | 未 push |

---

## §4 尺与"没做"的边界（诚实清单）

1. **`panel_profile_leaf_shape_windows_test.go` 零改动**（派单 ⚠ 钉函数行为）。定向复跑
   `Test314FixedPanelProfileLeafShapeIsRed`＋`Test314PanelProfileShapeRulerRingsOnTheFixedNameShape`
   （含 44 个子测试）＋`Test442CrossFileHelperReturnReachesTheSink`＋
   `TestAC4Distinctness133NoTwoLegsWearTheSameRow` 全绿（`450-shape-ruler.log`）。
2. **winlive 层只编译验证、未实跑**（`go vet -tags winlive` rc=0）；winlive 无 CI job，本腿不占机实跑。
3. **新死点（35g6）未诊断**——本腿边界是夹具＋载具。给下一条腿的**未验证线索**：
   `TestResyncOnShowDeliversFreshSnapshotFieldByField35g6` 在
   `panel_resync_stateless_35g6_windows_test.go`（文件名自称 stateless），经 `startPanelForTest`
   （`panel_resident_windows_test.go:71`）建宿主，而 `startPanelForTest` 的
   `uniqueProfileLeaf("wisp-33r5-panel-profile")` 前缀被**多个用例共享**——与 450 原死点同形状的
   下一嫌疑点，**仅是线索，未测**。
4. 派单举例的 `TestZZ315A1RootLoopFirstBridgeCall` 非整包尾（倒数第三）；锚因此取动态尾。
   台账/验收判据里若写死旧名，应以 `TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`
   或"动态尾"为准。

## §5 撤销口令

「450 落修撤」＝`git revert <本 commit>`；删 `D:\tmp\wisp\20261121\fix450`（临时件，只建不删；
仓内读数 `.scratch/wisp/probes/orch/20261121/450-*.log` 与 `450-test-list-tail-check.txt` 留档不删）。
