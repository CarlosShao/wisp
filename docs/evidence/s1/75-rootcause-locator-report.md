# 票 75 根因定位报告（只读代理产出，10:41）—— 由编排者代它落档

**为什么是我在写而不是它**：它开工前被指派写这个文件名，但落地时该文件**已被票 75 的实现代理认领并在改**
（`git status` 显示 +32/−5 未提交），它**两次被写入工具拒绝后选择停下来报告，而不是覆盖**
⇒ 这是共享工作树里正确的次序（A31/A34 那一族的正确应对），**这份报告因此没有丢**。
**⚠ 但它原本要写的文件现在归票 75 的实现者所有 ⇒ 本文件是唯一副本，别再合并、别再删。**

## Q1 反斜杠是在哪儿产生的：是"形状本身错了"，不是"比较错了"

一行无条件替换，藏在规范化步骤里：

| HEAD | 函数所在 | 出问题的那一行 |
|---|---|---|
| `186ad31` | `internal/risk/pathresolver.go:132` | **:133** |
| `131f722` | `internal/risk/pathresolver.go:138` | **:139** |
| `536998b` | `:167` | `:171`（票 75 的实现者已加了守卫） |

```go
u := strings.ReplaceAll(p, "/", `\`)   // 无条件执行；POSIX 绝对路径被折成反斜杠形状
...
return u                                // 即使没匹配到 UNC 前缀也照样返回
```

放大器：`Resolve` 在 `131f722:62` 又跑了一遍 `lexCanonical`（= `filepath.Abs` + `Clean`）。
在 Linux 上 `\tmp\x` **不算绝对路径**，于是 `Abs` 在前面拼上**进程 cwd**，`Clean` 又把 `\` 原样留着。
容器（`golang:1.27`，`HOME=/tmp/fakehome`，cwd `/src`）实测，`186ad31` 与 `131f722` 输出**逐字相同**：

```
"/tmp/wisp-probe/dir/file.txt" -> "/src/\tmp\wisp-probe\dir\file.txt"  resolved=false class=none
"/etc/passwd"                  -> "/src/\etc\passwd"                   resolved=false class=none
"~/probe/.ssh/id_testkey"      -> "/src/\tmp\fakehome\probe\.ssh\id_testkey" resolved=false class=B
"~/.git-credentials"           -> "/src/\tmp\fakehome\.git-credentials" resolved=false class=none
os.Stat("/src/\etc\passwd") -> no such file or directory
os.Stat("/etc/passwd")      -> exists
```

**证明是"产生端"而不是"比较端"**：把同一条路径两种写法分别喂给 `Classify`，**两种都判 A**
⇒ `blacklist.normPath/isUnder/baseName` 与 `tools.foldPath` 在形状上本来就是对称的，不是它们的问题。

## ⚠ 安全结论（在票 72 落地之后测的 `131f722`）：**R17 的不变式在 Linux 上仍然不成立**

`~/.git-credentials` → `none`、`~/.aws/credentials` → `none`、`~/.ssh/id_testkey` → **B**。
票 72 的 `formsOf` 只把**锚点那一侧**展开了，而**候选路径在进比较之前就已经在 `Resolve` 里被弄坏了**。
⇒ **A 表（永不可放行）在 Linux 上静默降级成 B 表（点一次确认就能放行）**，与票 72 修的是同一件事，
只是发生在另一侧。**这条必须在票 75 的验收里被单独问一次**，不能靠"测试变绿"过关。

## Q3 复现与逐条归类（用一次性树做的，补丁只打在仓库外，树已删）

| 测量 | 结果 |
|---|---|
| 逐字 CI 命令 @ `186ad31` | **4 个包红 / 49 条 `--- FAIL`**：risk 17、tools（FAIL + `panic: test timed out after 10m0s`，600.09s）、agent 1、memory 14 |
| `./internal/risk/` @ `131f722` | **12 条**（票 72 落地已经顺手治好了 17 条里的 5 条） |
| `./internal/tools/` @ `186ad31`，`-timeout 2400s` | **25 条**（20 顶层 + 5 子项），包在 **1211.3s** 失败。票面写的"19"在我测的 sha 上**复现不出来**（`24a66b6` 记为未证）；CI 自己的 600s 把它截断了 |
| 600s 超时的真名 | `TestVetoInsideTheWindowWritesNothing (4m48s)`，栈在 `approval/gate.go:473`。**未截断时有四条测试各烧满 C18 的 300s 拒绝超时**（300.01 / 300.02 / 300.02 / 300.03）＝ 1200s 里的 1211s |

**分层修复矩阵**（这才是给实现者的地图）：

| 树 | 打了哪些补丁 | risk | tools |
|---|---|---|---|
| 基线 `131f722` | 无 | **12** | 25 + 600s panic |
| A | P1 `normalizeLocalUNC` 守卫 + P2 `syncdirs` 用本机分隔符 | **10** | — |
| B | A + P3 `resolveHandle`→`EvalSymlinks` | **2** | **8，9.4s —— 超时消失** |
| C | B + P4 `dirOf`/`baseOf` 本机化 | **2** | **`ok`，0 FAIL，9.4s** |

**归类结论（三条不同性质，别混）**：
- **`internal/tools` 25/25 同一个根因，零例外。**
- **`internal/risk` 12 条 = 4 条同根因 + 8 条别的 bug + 1 条被 bug 掩盖的测试假设。**
  那 8 条的真因是 `pathresolver_other.go:10` 的 `DEFERRED(macOS/Linux)` 桩**无条件返回 `("", false)`**，
  而 `syncdirs.go:139 add()` 要求 `res.Resolved` ⇒ `s.complete` 在 Linux 上永远为假。
  **树 A 证明：形状修复治不了这 8 条。** ⇒ 归**票 55**（macOS 移植，realpath+lstat）。
  ⚠ **禁止为了让它们变绿而放宽 `add()` 的 `res.Resolved` 判定**——那是拿字符串启发去糊一个 realpath 缺口，方向是 fail-open。
- **`TestFourChannelExfilSuite/fs.write_into_sync_dir`（`provenance_test.go:135-137`）是"因为 bug 才绿"的用例**：
  它把 `C:\Users\test\OneDrive\Notes\shared.md` 硬编码在 POSIX 的 `t.TempDir()` 家里，
  **今天之所以通过，只因为 bug 把所有路径都判成同步嫌疑**。修好之后它会**变红**
  ⇒ 正确处置：挪进 `_windows_test.go` 层，**并补一条用 `t.TempDir()` 构造的 POSIX 等价用例**，
  否则覆盖面会静默下降。

## Q4 最小正确改动面（提议，不是实现）—— 每一项在 Windows 上都必须是"恒等"

理由共用：`lexCanonical` 里的 `filepath.Clean` 在 Windows 上已经把 `/` 折成 `\`，
所以任何"`/` 开头才走"的分支在 Windows 上根本不会被取到。

1. `internal/risk/pathresolver.go`（**冻结/D22**）—— `normalizeLocalUNC`：只在真是 UNC 形状（`\\` 或 `//`）时才动。
2. `internal/risk/syncdirs.go` —— `splitPathComponents` / `deepestExistingAncestor` / `resolveTarget` 重新拼接：
   POSIX 根用本机分隔符；`\\?\`、`\\server\share`、`X:` 三个分支保持逐字不变。
   ⚠ `//server/share` **既是合法 POSIX 路径又是本仓的 UNC 候选** ⇒ 两个分支的先后顺序必须显式写清并测。
3. `internal/tools/fs_write.go` —— `dirOf`、`baseOf`。顺带检查 `recycle_windows.go:226-227`
   与票 73 挂在 `fs_write.go:282` 的 `reclaimStaging(parent,…)`（孤儿清扫必须还能匹配 `.wisp-tmp-*`）。
4. `internal/tools/paths.go` —— `const sep` 改成 `string(filepath.Separator)`、`foldPath` 同步。
   ⚠ **这条不是卫生问题**：POSIX 上 `\` 是合法文件名字符，无条件 `/`→`\` 会让 `a\b` 与 `a/b` 折成同一个串，
   而 `bOverrides` 与白名单根都是**在这个折叠串上做 map/前缀查** ⇒ **跨文件放行泄漏，方向是 fail-open**。
   与票 72 的"放行侧只认已解析形式"是同一条不变式的另一处漏口。
5. **不属于票 75**：`pathresolver_other.go` 的 `resolveHandle`/`reparseComponents`（realpath+lstat）归票 55。

## 未证 / 局限（它自己声明的，我照录）
- **Windows 那一半它没测**（只读代理不跑 Windows 测试）⇒ "在 Windows 上恒等"目前是**推理**，不是证据。
  ⇒ 票 75 验收时这一格必须由实现者补：票 18 的 `pathresolver_junction_windows_test.go` Case 1–10、
  票 20 的 `bridge_junction_windows_test.go`、票 73 的 `fs_staging_windows_test.go` 与
  `bridge_a18_kill_windows_test.go`，全部 `-count=2` 真跑。
- CI 命令的行号：草稿里写的 `ci.yml:127-134` 是错的，实际在 `:107-113`。
