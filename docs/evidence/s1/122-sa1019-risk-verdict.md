# 122 — risk 那条 SA1019 的裁决建议：解冻 vs 挂票 vs 不修

- **出件腿**：`122-sa1-census`（只读取证，零 `.go` 改动）
- **时间**：2026-11-08 15:2x +08 · 起手锚 `31482f271`（批 2 交件后 HEAD）
- **对象**：`internal/risk/proccpu_windows_test.go:20` 那一条 `SA1019`
- **本件性质**：裁决建议（判据与出处齐全），**不是**解冻授权——解冻须编排者具名落派单＋台账补 `A##`。

---

## 1. 用的哪个 deprecated API、替代是什么、改法几行、行为是否逐位不变

**现场（2026-11-08 现读全文，29 行）**：`proccpu_windows_test.go:20` 调 `windows.GetCurrentProcess()`（`golang.org/x/sys/windows`）。

**deprecated 标记原文**（现跑，module cache `golang.org/x/sys@v0.49.0/windows/syscall_windows.go:481-485`）：

```go
// The returned error is always nil.
//
// Deprecated: use CurrentProcess for the same Handle without the nil
// error.
func GetCurrentProcess() (Handle, error) {
```

- **API**：`windows.GetCurrentProcess`（返回伪句柄＋一个恒为 nil 的 error）。
- **推荐替代**：`windows.CurrentProcess` —— upstream 逐字写明 "the same Handle without the nil error"，是 drop-in。
- **改法**：改 1 行＋删 3 行：
  - `:20` `proc, err := windows.GetCurrentProcess()` → `proc := windows.CurrentProcess()`
  - 删 `:21-23` 的 `if err != nil { return 0, err }`（死码分支，见下）。
  - 可选：照 speech 先例（`internal/speech/cgo_boundary_windows_test.go:36-38`，批 1 `4df9964d1` 已落）加 2–3 行注释说明"非弃用伪句柄访问器、不可失败"。
- **行为是否逐位不变**：**是**。两个访问器返回的是同一个 Windows 伪句柄常量（当前进程伪句柄），传给 `:25` 的 `GetProcessTimes` 的句柄值逐位相同；`GetCurrentProcess` 的 error 恒为 nil（upstream 原文 :481），`:21-23` 是从未走过的分支，删除不改变任何可观测行为。函数签名、返回值、全部调用方（仅 `pathresolver_budget_norace_test.go` 的 5 处调用，现跑 grep 坐实）均不动。

## 2. `internal/risk` 为什么冻结——具名出处

**如实结论：`docs/PLAN.md` 与 `docs/specs/**` 里没有「internal/risk/** 冻结」这一句字面**（两处现跑 grep `冻结…risk|risk…冻结` 均 0 命中）。冻结是编排者层面的纪律，出处按时间序：

1. **台账 A35**（`docs/reports/pending-and-issues.md:883`，票 72 时期）：「⚠ 它要动 `internal/risk/`（冻结区）⇒ **改法必须先报我判**」——最早把 risk 定为冻结区的裁定。
2. **冻结机制的上位依据**：`docs/specs/README.md:3,13`（「C1–C32 冻结契约……任何与冻结契约的偏离 = 契约变更，必须人工批准」）＋ `SPEC-01 §3`（模块边界冻结）。risk 包承载四枚冻结契约的实现：C18/C26（`rules_gateway.go`，票 123 票面 `:8` 「冻结 C18/C26 语义」）、C19（`assessor.go:168` 冻结缝，台账 `:15992`）、C25（`provenance.go`/`taintmatch.go` 判据文字，台账 `:8114/:8224` 逐字钉 md5）、C26（`pathresolver*.go`，票 102 票面 `:16,42`）。
3. **现行禁改列**：票 122 票面 `:21` 自己写着「禁改：……`internal/risk/**`（冻结）」；票 104/107/113/114/115/117/118/119/120/121/138/140/139 等同款逐票复制。
4. **解冻的既定口径与先例**：
   - 台账 `:5235` 一带（票 125/126 裁定）：冻结身份不随"别的票刚动过"自动消失，禁改列不会「在没人决定的情况下自己缩短」；**解冻必须由编排者具名决定并落 `A##`**，正确措辞是「编排者在派单里解冻了它」（票 126 `:401-403` 教训原文）。
   - 票 130（`:8-10`）先例：曾解冻过一次并登记，教训是「申请解冻前先核依赖方向」。
   - **测试面不在冻结语义承载面的实践先例**：票 86（`4d0e2577d`）本枚文件就是以"测试仪器件"身份进 risk 包的；票 216（台账 `:18993`）在 `internal/risk/**` 里**新建两枚 `_test.go`**，裁定口径＝「`internal/risk/**` 的**非测试字节零动**」——即冻结的执行口径自始区分产码契约面与测试面。

## 3. 测试件还是产码？

**测试件，且是纯测试 helper 件**（现跑坐实）：
- 文件名 `_test.go` 后缀＋`//go:build windows`，只参与 `go test` 构建，不进任何产品二进制；
- 全文件只有一枚函数 `processCPUTime()`（无 `func Test*`），是 C26 每调用预算的测试仪器（`:11-14` 注释自述，票 86 立的）；
- 引用面 grep：只有 `pathresolver_budget_norace_test.go`（同为测试件）5 处调用；`proccpu_unix_test.go`/`proccpu_unsupported_test.go` 是它的平台孪生件（同签名、unix 侧不用此 API）；
- **零产码调用者**。改动它的爆炸半径＝risk 包测试构建，产品字节零动。

## 4. 裁决建议（我拍，理由如下）

**甲：具名解冻这一枚文件这一条，让批 2（或下一批）顺手修。** 撤销口令：「撤 risk 测试件解冻，SA1019 留红」。

理由（按铁律 ⓪⓪⓪，判据＝最有价值＋主流做法，不是风险最小）：

1. **主流做法**：staticcheck `SA1019` 的设计意图就是让 deprecated 调用被换成新写法；upstream 给的是逐字 drop-in 替代（"the same Handle without the nil error"），修法 1 行、行为逐位不变。在测试件里换掉 deprecated API 是标准维护动作，不是语义变更。
2. **冻结防的是什么**：risk 冻结防的是 C18/C19/C25/C26 **契约语义漂移**（assessor/pathresolver/provenance/rules_gateway 那几枚判据文字与冻结缝），不是钉住一行 deprecated 系统调用访问器。本仓自己的执行口径也一直这么划：票 86 立仪器件、票 216 动测试面都判「非测试字节零动」即可。本枚修法**非测试字节零动**，与两枚先例同形。
3. **为什么不是乙（单开票）**：对 1 行 drop-in 替换开新票＝负价值——开票、派单、非实现者验收、CI 读数核验的全套开销买不到任何超出"顺手 1 行"的收益；且 SA1019 清账的归属本来就在票 122 里（台账 `:32092` 登记的正是这个问题），在同一张票内具名解冻是最短的正确路径。
4. **为什么不是丙（留着不修）**：windows 侧 staticcheck 永远背 1 条红，AC#6 的 CI 读数（批 1 期望 109、批 2 期望 98）兑不了现——每次都要向每个读数的人解释一遍"那条是冻结面"，账目噪音永久化；deprecated 调用留着没有任何收益。

**边界（若编排者采纳甲）**：
- 解冻**只具名** `internal/risk/proccpu_windows_test.go` 一枚文件、**只许**本件 §1 的修法（改 `:20`＋删 `:21-23`＋可选注释），不许顺手动 risk 包任何其他字节；
- `assessor.go`/`pathresolver*.go`/`rules_gateway.go`/`provenance.go`/`taintmatch.go` 照旧冻结，一枚不许动；
- 措辞照票 126 教训写「编排者在派单里解冻了它」，台账补 `A##`。

## 5. 现跑凭据（复跑尺）

| 凭据 | 命令 | 读数 |
|---|---|---|
| 对象在场 | `sed -n '20p' internal/risk/proccpu_windows_test.go` | `proc, err := windows.GetCurrentProcess()` |
| deprecated 原文 | `grep -B4 "func GetCurrentProcess" "$(go env GOMODCACHE)/golang.org/x/sys@v0.49.0/windows/syscall_windows.go"` | `:481` error 恒 nil、`:483-484` Deprecated→CurrentProcess |
| 引用面 | `grep -rn "processCPUTime" internal/risk` | 仅 3 枚测试件、5 处调用、零产码 |
| PLAN/specs 无字面 | `grep -rn "冻结.*risk\|risk.*冻结" docs/PLAN.md docs/specs` | 0 命中（如实：冻结出处是台账 A35＋各票禁改列） |
| 批 1 先例 | `git log --oneline -1 -- internal/speech/cgo_boundary_windows_test.go` | `4df9964d1`（票 122 批 1） |
| 本枚文件来历 | `git log --oneline -1 -- internal/risk/proccpu_windows_test.go` | `4d0e2577d`（票 86，测试仪器件） |
