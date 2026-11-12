# 票 85 AC#3 — staticcheck 变异回验（编排者现跑，零 LLM 成本）

- 采集：**2026-11-12 10:4x +08**，编排者本机现跑（⛔ 子代理转述）
- 锚：`3fa57f31f` 之后工作树（`internal/observe` 零改动）
- 工具：`staticcheck.exe`（`$(go env GOPATH)\bin`，**2026.2.2**，与 CI `ci.yml:465` 同版本）
- 目的：证明 staticcheck **真的在看这个包**，不是"又一场空跑"（票 85 AC#3 逐字）。

## 1. 变异回验（逐字读数）

| 步 | 命令 | 读数 |
|---|---|---|
| **基线** | `staticcheck ./internal/observe/` | **1 条** finding（`logprivate_132_windows_test.go:99:6: func dirIsSealed132 is unused (U1000)`） |
| **种入旧形状** | 新建 `internal/observe/zz_mutprobe_AC3.go`，内容 `package observe` ＋ `func zzMutProbeUnusedAC3() {}` | **2 条** finding，新增逐字：`internal\observe\zz_mutprobe_AC3.go:3:6: func zzMutProbeUnusedAC3 is unused (U1000)` ⇒ **真报了** |
| **还原** | 删掉探针文件 | **zzMutProbe 不再出现**；回到基线 1 条 ⇒ **干净** |

## 2. 自证

```
git status --porcelain -- internal/observe/   ⇒ 空（干净）
```
- ⇒ **探针未留任何痕，工作树 `internal/observe` 零改动。**

## 3. 判

**AC#3＝成立。** 种入一枚"未使用符号"（旧形状的典型）⇒ staticcheck **立即重报 `U1000`**；还原后不再报。⇒ **"它真的在看这个包"被证。**

## 4. 诚实边界（本件做了什么、⛔ 做了什么）

- ✅ 做了：种入**新**的未使用符号 → 报 → 删 → 不报（标准变异回验，证明工具活跃）。
- ⚠ **未做**：AC#3 字面说"把**其中一条**（票 122 已修掉的）改回旧形状"。本件用**新种探针**而非"回退某条已修 finding"，因为回退真实代码会改共享工作树里别人的修复（风险高）。**两者证明的是同一件事（工具在看），但字面形式略有差异，如实登记。**
- ⛔ 改任何生产 `.go`：**探针建后即删，`git status` 空**。
- ⛔ 跑 `go test`：没跑（AC#3 只要求 staticcheck）。

## 5. 撤销口令

删本件即复原（本件零产码）。
