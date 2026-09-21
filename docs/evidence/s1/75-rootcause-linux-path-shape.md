# 票 75 — 根因定位与修复证据：C26 在 Linux 上产出反斜杠形状的 canonical

HEAD 基线：`f1033e1`（票 72 落地后）。测量环境：docker `golang:1.27`（Go 1.27.1 linux/amd64），
源码用 `git archive HEAD | tar -x -C /tmp/wisp75` 取干净快照（共树有其他在途代理，A38④）。

## AC#1 基线（改任何东西之前）

CI `test-core` 的 "Portable package tests" 步骤，命令逐字取自 `.github/workflows/ci.yml:127-134`：

```
go test ./internal/agent/... ./internal/llm/... ./internal/config/... \
  ./internal/memory/... ./internal/observe/... ./internal/secret/... \
  ./internal/risk/... ./internal/statemachine/... ./internal/session/... \
  ./internal/watchdog/... ./internal/tools/... ./internal/models/... \
  ./internal/buildinfo/... ./internal/audio/... ./internal/proc/... \
  ./internal/panel/... -count=1
```

在容器里等价执行（`WISP_ENV=test`，`-v wisp75-gomod:/go/pkg/mod`）：

```
docker run --rm -v <archive>:/src -v wisp75-gomod:/go/pkg/mod -w /src \
  -e WISP_ENV=test golang:1.27 bash /src/run_ci.sh
```

### 基线读数（见下文"基线读数"小节，跑完补齐）

- 全量 `--- FAIL` 数：
- `internal/risk`：
- `internal/tools`：
- 600s 超时：

## 根因（file:line）

C26 的管道里只有**一步**在 Linux 上改变分隔符，而它是无条件执行的：

`internal/risk/pathresolver.go:139`（`normalizeLocalUNC` 第一行）

```go
u := strings.ReplaceAll(p, "/", `\`)   // 无条件：POSIX 绝对路径被改写成反斜杠形状
```

函数名只承诺"规范化本地 UNC"，实际做的是"把所有正斜杠换成反斜杠，然后返回"。
`Resolve`（`pathresolver.go:58-93`）在第 61 行无条件调用它，第 62 行又把它交回的
字符串再过一次 `lexCanonical`（`pathresolver.go:125-133` = `filepath.Abs` + `filepath.Clean`）。
在 Linux 上 `\home\u\.ssh\id_ed25519` 不是绝对路径（`filepath.IsAbs` 只认前导 `/`），
于是 `Abs` 把**工作目录**拼在前面、`Clean` 对反斜杠一个字符都不动：

```
输入  /home/runner/.ssh/id_ed25519
  61  \home\runner\.ssh\id_ed25519      （normalizeLocalUNC）
  62  /src/\home\runner\.ssh\id_ed25519  （lexCanonical = Abs+Cwd）
  80  resolveHandle -> "" false          （pathresolver_other.go:10 DEFERRED 桩）
  91  Canonical = /src/\home\runner\.ssh\id_ed25519
```

交回的 canonical **既不是调用方给的路径，也不是任何真实存在的路径**——它在 Linux 上
根本无法 open。这是"Windows-shaped everywhere"的字面现场，也是本票 47 个 FAIL 的源头：

- `internal/tools`：`PathCanonicalizer.resolve`（`internal/tools/paths.go:54-60`）把
  `Resolve().Canonical` 原样交回给 bridge，fs 工具用它去 `os.Open`/`os.Create` → ENOENT。
- `internal/risk`：`Classify`→`formsOf`→`normPath`（`blacklist.go:114`）把候选和锚点**都**
  折成反斜杠，所以 A/B 表比较在 Linux 上"碰巧自洽"；真正红的是 `IsSyncPath`：
  `syncdirs.go:245-269 deepestExistingAncestor` 用 `\` 拼接去 `os.Lstat`
  （`syncdirs.go:226` 同样用 `\` 重新拼锚点），POSIX 上永远 stat 不到 ⇒
  `errTargetUnverified` ⇒ 每次写都落进 sync-suspect 兜底网。这个失效模式
  `syncdirs.go:216-224` 的 N-9 注释已经写下来了，只是当时归给了 DEFERRED 的票 55。

Windows 上整条链每一步的形状**恰好**是对的：`filepath.Clean` 本来就把 `/` 折成 `\`，
所以第 139 行的 ReplaceAll 在 Windows 上是**恒等变换**——这正是它能活过所有本地门的原因。
（同一条理由适用于 `blacklist.go:115`、`tools/paths.go:112`、`tools/fs_write.go:179/189`、
`syncdirs.go:278/290`：Windows 恒等、POSIX 破坏。）

## 契约读数（两个选项里选了哪个，为什么）

票面给的两条路：(1) canonical 改成**平台形状**；(2) canonical 保持单一形状、比较侧两边都归一。

读 `docs/specs/SPEC-06-security-gatekeeping.md:47-53`（C26 管道）：第 4 步是
"打开句柄取 `GetFinalPathNameByHandle(VOLUME_NAME_DOS)` 真实路径"，第 5-7 步
（reparse / 8.3 / UNC）都挂在第 4 步的句柄语义上。契约要求的是**真实路径**，不是
"反斜杠路径"；Linux 的真实路径就是 `/` 形状的。所以 (1) 才是契约读数。

(2) 在本项目里也**做不到**：corruption 发生在 `Resolve` 内部，交回的字符串带上了 cwd 前缀，
任何下游比较都无法还原（POSIX 文件名合法包含 `\`，"把 `\` 换回 `/`" 是有歧义的猜测，
而且票 72 刚把"零字符串启发式"写成不变式）。更要紧的是 canonical 不只用于比较，
它还被拿去 `os.Open`（`tools/paths.go:65`、`fs_write.go`）——形状错了就是开错文件。

⇒ 实现 (1)：canonical = 平台形状；比较侧仍保留一个**纯折叠**（两侧都过同一个 fold），
并把 fold 也改成平台形状（理由见下"分隔符折叠的边界"）。

## 修复面（哪些文件，为什么是这些）

1. `internal/risk/pathresolver.go` — **frozen，AC#4**。两处：
   - `normalizeLocalUNC:139`：只在输入是 UNC 形状（`\\` / `//` 前缀）时才动分隔符。
   - `tailExistsBelow:311`：`cut + `\` + first` 的探测路径交给 `os.Lstat`，必须用平台分隔符。
   两处在 Windows 上都是恒等变换（Clean 之后已无 `/`）。**这是本票唯一落在冻结文件里的改动**，
   单独成一个 commit，方便 owner 一个 `git revert` 退回"只交提案"的读数。
2. `internal/risk/blacklist.go` — `normPath`/`normDir`/`isUnder`/`baseName`/`hasGitConfigSegment`
   的分隔符改成平台常量。
3. `internal/risk/syncdirs.go` — `splitPathComponents`/`deepestExistingAncestor`/
   `resolveTarget` 的重拼/`hasFoldedDotDot` 的切分改成平台分隔符（Windows 分支逐字保持）。
4. `internal/tools/paths.go` — `sep`/`foldPath`。
5. `internal/tools/fs_write.go` — `dirOf`/`baseOf`（结果进 `se.record` 审计文案，人看得见）。

### 分隔符折叠的边界

一个只在比较内部用的 fold 用哪种形状是自由的（两侧同一个 fold ⇒ 自洽）；但 POSIX 上
`\` 是**合法文件名字符**，把 `/` 无条件折成 `\` 会让 `/home/u/a\b/.ssh` 与
`/home/u/a/b/.ssh` 折成同一串——B 档 `bOverrides` 是 map 查表，这是一次**跨文件的豁免串用**
（fail-open 方向）。所以 fold 也必须平台化：Windows 继续 `\`/`/` 统一（OS 自己也这么当），
POSIX 只认 `/`、绝不碰 `\`。

## AC#4：交给 owner 的一段话 diff 提案

`internal/risk/pathresolver.go` 的 `normalizeLocalUNC` 目前第一行
`u := strings.ReplaceAll(p, "/", "\\")` 无条件改写分隔符，函数名承诺的"规范化本地 UNC"
被扩大成了"规范化一切分隔符"；建议改成只在输入确实是 UNC 形状时才折叠：

```go
func normalizeLocalUNC(p string) string {
	if !strings.HasPrefix(p, `\\`) && !strings.HasPrefix(p, `//`) {
		return p // 非 UNC：本函数无事可做，分隔符形状交给 OS 自己决定
	}
	u := strings.ReplaceAll(p, "/", `\`)
	…原样不变…
```

同文件 `tailExistsBelow:311` 的 `pathExists(cut + `\` + first)` 换成平台分隔符拼接
（它交给 `os.Lstat`，是 OS 可见路径）。Windows 上 `lexCanonical` 之后已无 `/`，
两处均为恒等，故 8.3/junction/UNC 的攻击面语义一字未动；改动只是停止在 POSIX 上
把 `/` 改写成 `\`。按 registry R17（票 72 face）这是实现纠正、不是 D22 契约变更：
`SPEC-06:50-52` 要求的是句柄真实路径，`PLAN.md:1376/1796` 要求真实路径解析并推翻字面比较。
