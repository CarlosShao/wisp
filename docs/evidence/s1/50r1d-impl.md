# 50-r1d AC2: host_api 版本门（caret gate）

## AC2 读数（host_api ^1.0 vs host 1.x loads / 2.x refuses）

**测试矩阵与通过率**：`go test ./internal/plugin -count=1 -run TestAC2CaretGate` → **37 例全绿**。其中：

* **装载通过用例**：^1.0/host1.0.0、^1.0/host1.2.3、^1.0/host1.99.99、^1.2/host1.2.0、^1.2/host1.5.0、主版本≥1 且 min 之上的各类边界；以及主流 caret 语义的 `^1`（等同于 `^1.0.0`）。
* **拒绝通过用例**：主版本不匹配（^1.0 vs 2.0.0、^2.1 vs 1.9.9）、低于 min（^1.2 vs 1.1.9、^0.2 vs 0.3.0）——错误信息包含 **"refused"**, **"major-version incompatibility"**, **"below the minimum"**。
* **grammar 不可解析即拒**：`"1.0"`、`">=1"`、`=1.0"`、`"^*"`, `"^"`, `""` 等 C16 未定义语法均以 `"unsupported range syntax"` 或 `"caret needs a version"` 明确拒绝，不以静默方式绕过。
* **fail closed**：dev 默认 `DefaultHostVersion="0.0.0-dev"` 不匹配任何 `^1.x` 范围 → 拒绝所有插件加载。

**反向判据**（改坏匹配⇒点名用例红）已验证：

| 突变 | 描述 | 红案数（命名） |
|---|---|---|
| M1 | 去掉主版本相等检查 → 任意主机接受 | 5 例：`caret1.0/host2.0.0` / `caret1.0/host2.4.1` / `caret1.0/host3.0.0` / `caret2.1/host1.9.9` / `caret_major_only_upper` |
| M2 | 去掉 min 下界检查 | 2 例：`caret1.2/host1.1.9_below_min` / `caret1.2.3/host1.2.2` |

每次突变后均**恢复原核并复测全绿**（见下文 commit diff）。

## 改动文件清单（仅 internal/plugin）

以下 3 个文件构成本次交付，git 显式 pathspec 提交：

```text
internal/plugin/hostapi.go       337 lines (new)  // HostVersion/CaretRange 解析 + CaretRange.Check + CheckHostAPI + Manifest.CompatibleWith + DefaultHostVersion=buildinfo.Version
internal/plugin/hostapi_test.go  156 lines (new)  // 37 例 AC2 矩阵 + 装载入口联动 + grammar 封口 + fail-closed 默认值
internal/plugin/validator.go       8 lines (改)   // ParseManifest 侧把 host_api 的 C16 range 语法封口，不可解析的 range 进不了 gate
```

**commit**：`0e8f2ccea`（3 files changed, 499 insertions(+), 2 deletions(-)）。
**diff 核对**：两次突变（M1/M2）之后 `git status --porcelain -- internal/plugin` 与 `git diff --stat HEAD -- internal/plugin` 均为**空**（还原核＝文件回到提交态，非“看起来差不多”）。

## gates 读数（本机复现，2026-11-02 22:0x +08）

```text
gofumpt -l internal/plugin                 -> 无输出
go vet ./internal/plugin                   -> 无输出
go test ./internal/plugin -count=1         -> ok  0.131s
go test -run TestAC2CaretGate -v -count=1  -> 37 PASS / 0 FAIL（subtest 级）
go test -run TestAC2 -v -count=1           -> 4 个顶层测试全 PASS
go test -v -count=1（整包）                 -> 26 个顶层测试全 PASS（含 50-r1c 的 AC1/AC4 用例，未回归）
tools/d22scan/d22scan.exe internal/plugin  -> clean - no D22 ban violations, exit=0
```

## AC3 / AC6：后续腿登记（not done this round）

* **AC3 install flow e2e incl. authorization UI showing capabilities + source; uninstall teardown leaves zero goroutines/files (DisposalScope test)**  
  → **待后续腿**。本 leg 只完成 C16 gate（load-time check），install/uninstall/authorization UI 不在范围内。

* **AC6 plugin tools visible in list_tools/discovery + execution via bridge records plugin id**  
  → **待后续腿**。本 leg 不涉及 `list_tools` 工具注册和调用链路审计。

## TODO（自动继续时不再重复本文档的“未定义”判断）

* host API 进一步扩展（支持更多 range 运算符）**不是实现细节，而是契约扩大**，需机主人工批准（AGENTS §0: 改契约＝人工批准）。
* 与 cmd/wisp 的 load hook 集成在票 50 后续腿完成。
