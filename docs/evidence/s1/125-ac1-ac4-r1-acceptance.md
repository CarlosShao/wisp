# 票 125 — AC#1..AC#4 第一份独立裁决表（r1，验收方 `acceptor-ticket125-r1`）

**本票此前状态**：`ls docs/evidence/s1/ | grep -i 125` 现量 **无命中**（该目录 138 枚文件，一枚都不属于票 125）
⇒ 本表是这张票的**第一份**裁决表。票面四格 AC 的 `[x]` **全部是实现方自勾**，本表**不翻任何一枚勾**（翻勾权在编排者终裁之后）。

**角色**：本程**非实现者**，只做裁决。全程**零枚生产码/测试码改动**（所有变异都造在仓外快照）。
唯一写件＝本文件。

**三档证据标记**在本表每一格逐节打：〔独立复现〕／〔日志＋归档，抽验〕／〔仅自述，不背书〕。

---

## §0 锚点自量 · 争用闸门 · 被验版本的盘上身份

### 0.1 锚点（现量，`date` 每次重取）

```
$ date
Thu Sep 24 15:12:28 CST 2026
$ git rev-parse HEAD
182daed377feeff20f73ea47675c121df5607ca5
$ git cat-file -t 182daed377feeff20f73ea47675c121df5607ca5
commit
$ git rev-parse --abbrev-ref HEAD
dev
```
本节写完时的第二个时刻：`date '+%H:%M +08'` → **`15:38 +08`**。

### 0.2 争用闸门（先做再取数）

- 15:12 现量 `gh run list --branch dev --limit 3`：`35967768017` = **`in_progress`**（15:05 编排者那次 push 自启）。
  ⇒ **该时刻不开测**。§1 之后的每一发读数都在它变成 `completed` 之后取。
- 15:2x 复量：同一枚 run 变 `completed failure 8m38s`；`docker ps` 里 **无 `golang`/`wisp` 容器**
  （只有 union-proxy / clipsync×4 / minio / postgres / redis 常驻件，非 CPU 突发型）。
  `worker-ticket137-ac4` 那枚容器**全程没出现在 `docker ps`** ⇒ 没有撞窗。
- ⚠ 那枚 `ci` run 判 `failure` 与本表无关，本表**不引用任何 CI 结论当读数**（编排者记忆：引用门禁前先查有几枚 green）。

### 0.3 仪器（现量，不采信任何"这台机器没有 Docker"的说法）

```
$ docker version --format '{{.Server.Version}}'
29.6.2
$ docker run ... golang:1.27 sh -c 'uname -s; go version'
Linux
go version go1.27.1 linux/amd64
```
**未执行 `docker pull`**（镜像本机已有）。快照一律 `git archive <sha> | tar -x -C /d/tmp/wisp125r1-<名>`，
挂载一律 `-v /d/tmp/...:/wisp`，**每把先 `ls -l /<mount>/go.mod` 自证非空再开测**（全部读到 `883` 字节）。

### 0.4 被验版本的盘上身份（禁读脏工作树）

工作树**不干净**（`design/**` 16 枚被别人挪走、`internal/winsec/*_other_test.go` 有 `worker-ticket137-ac4` 在飞的改动）。
**本表一处都不读工作树**：所有读数取自 `git archive` 出的纯净快照。

| 用途 | sha | 盘上自证 |
|---|---|---|
| **被验的主版本（AC#2 生产码那一笔）** | `4824bb8`（09-22 20:58 `fix(winsec,125,AC#2)`） | `git cat-file -t` = commit，`git show --name-only` 列 4 枚路径 |
| 改前对照 | `ff3faf9`（= `4824bb8^`，09-22 20:43，AC#1 那笔） | 快照 `wisp125r1-pre`，`grep -c resolverProbeRoot resolve.go` = **0**（确认是未修版） |
| 票面末格（AC#3/AC#4/补记/补记二） | `bd50c63`（09-22 21:20，现量是该票面文件的**最后一枚** commit） | 快照 `wisp125r1-gate` |
| 实现方给的控制组 | `81b4d5f`（09-22 20:22，票 126 收格） | 快照 `wisp125r1-ctrl` |

**两条必须说清的版本事实**（否则读者会把本表当成"HEAD 的读数"）：

1. **`bd50c63`（票面末格）与 `4824bb8` 的代码逐字节相同**：
   `git diff --name-only 4824bb8 bd50c63 -- internal/ cmd/ scripts/ tools/` 现量 **空输出**。
   ⇒ 本表 §1/§4 用 `bd50c63` 取的门禁数，就是 `4824bb8` 那一版的门禁数；`88eba32`／`bd50c63` 两枚补记
   `git show --name-only` 现量各只列票面一枚 `.md`。
   ⚠ 顺带纠正一处口径：实现方在"补记"里说"真正的末格 sha 是 `6d8554e`"，而盘上现量票面文件的最后一枚是 **`bd50c63`（补记二）**——
   这一枚它自己没量（内容纯 `.md`，结论不受影响，但"最终 sha 必须有读数"那句按字面看**差一枚**）。
2. **HEAD 的 `resolve.go` 已经不是票 125 那一版**：`git log -- internal/winsec/resolve.go` 现量
   `a45b2e9`（票 129，09-22 **23:41**，"绝对性进两枚比较(只更严)"）**晚于** `4824bb8`。
   ⇒ 本表 §2 那句"`sameTree`/`answerInsideTree` 一字未动"判的是 **`4824bb8^ → 4824bb8` 这一刀**，
   不覆盖票 129 后来对同两枚函数做的事（那是票 129 的账，不在本票射程）。

### 0.5 本表读的 diff 原文（生产码那一刀，可逐字节重放）

`git diff 4824bb8^ 4824bb8 -- internal/winsec/resolve.go` 现量 **+82/−2**，非注释非空行的**全部**增删行只有这些：

```
+	"io/fs"
-	shapes := []string{os.TempDir() + sep + ".." + sep + "wisp-103-conformance-probe"}
+	shapes := []string{resolverProbeRoot() + sep + ".." + sep + "wisp-103-conformance-probe"}
+func resolverProbeRoot() string {
+	return resolveProbeRoot(os.TempDir())
+}
+func resolveProbeRoot(path string) string {
+	if path == "" {
+		return path
+	}
+	var missing []string
+	cur := path
+	for {
+		_, err := os.Lstat(cur)
+		if err == nil {
+			break
+		}
+		if !errors.Is(err, fs.ErrNotExist) {
+			return path // unreadable, not absent: not this guard's to reinterpret
+		}
+		parent := filepath.Dir(cur)
+		if parent == cur {
+			return path
+		}
+		missing = append(missing, filepath.Base(cur))
+		cur = parent
+	}
+	real, err := filepath.EvalSymlinks(cur)
+	if err != nil {
+		return path
+	}
+	for i := len(missing) - 1; i >= 0; i-- {
+		real = filepath.Join(real, missing[i])
+	}
+	return real
+}
-	parent := os.TempDir()
+	parent := resolverProbeRoot()
```
机器可核的"没动哪几枚"：把上面 diff 的 `+/-` 行（去掉文件头）喂给
`grep -iE 'sameTree|sameVolume|foldSegment|answerInsideTree|resolverConformanceFailure|treeOwnershipFailureForPair'`
⇒ 现量 **零命中**（`rc=1`）。这是"被动没动"的**证据**，不是"实现方说没动"的转述。

---

## §1 AC#1 —— POSIX 侧那枚"能红的钉子"

### 1.1 判据与读数

**判**：**成立**（无附条件）。证据档：**〔独立复现〕**（本程自己在容器里跑的变异与两形读数）。

钉子本体：`internal/config/c26_seam_posix_125_test.go:52`（`//go:build !windows`，`package config`，
`TestAC1POSIXSeamHoldsC26Pipeline125`），断言三件：
`PathResolverInstalled() != nil` → `fmt.Sprintf("%T", r) == "risk.c26Pipeline"`（**逐字面量，不是拿被测函数算的**）
→ 一枚只有真管线做得到的行为 leg（对 `<base>/../wisp-125-probe` 给出不含 `..` 的折叠答案）。

**"钉子能红"这一发是本程自己造的，不是复述实现方的 MUT-125A2。**

**MUTATION-R1A（本程造，让安装被跳过）**，落在仓外快照 `D:/tmp/wisp125r1-mutA`：

```diff
@@ -154,7 +154,7 @@
 			"consequence", "the incumbent resolver stays in place; nothing was uninstalled")
 		return
 	}
-	if reason != "" {
+	if reason != "" || r != nil { // MUTATION-R1A the install is skipped for every candidate
 		slog.Error("winsec: refusing to install a path resolver into the sealing seam",
```

容器真跑（`ls -l /mutA/go.mod` → 883 B 自证挂载）：

```
mutA VET_RC=0
mutA TEST_RC=1
RUN=153 TOPPASS=83 TOPFAIL=2 SUBPASS=66 SUBFAIL=2 PANIC=0
-- 红名 --
--- FAIL: TestAC1POSIXSeamHoldsC26Pipeline125 (0.00s)
--- FAIL: TestAC2POSIXSeamInstallSurvivesASymlinkSpelledTemp125 (0.02s)
    --- FAIL: .../control_plain_temp (0.01s)
    --- FAIL: .../measured_symlink_spelled_temp (0.01s)
```
红因原文逐字（第一手，不是转述）：
```
c26_seam_posix_125_test.go:58: AC#1 RED: winsec.PathResolverInstalled() = <nil>: this process links internal/risk,
whose init() is supposed to hand C26 to the sealing seam, so either the install never ran or the seam guard refused it
```
⇒ **"让安装被跳过" ⇒ 钉子红、红名逐字点到 `TestAC1POSIXSeamHoldsC26Pipeline125`**。钉子不是装饰。

同一棵快照的**对照组**（未变异 `4824bb8`）：`TEST_RC=0`、`RUN=153 TOPPASS=85 TOPFAIL=0 SUBPASS=68 SUBFAIL=0 PANIC=0`。
两把 `RUN` 同为 153、`PANIC=0` ⇒ 这格的读数**没有被一条 panic 吞掉**（本程每一发都算了 panic 数，全 0，见 §4）。

### 1.2 三件附带核查（都做在钉子上，免得"能红"是空话）

1. **零 Skip**：`awk '/^func TestAC1POSIXSeamHoldsC26Pipeline125/,/^}/' … | grep -c 't.Skip'` → **0**。
   实现方那句"零 Skip 条件"字面成立。
2. **落点理由为真**（这格最容易被读成"钉在别处躲事"）：容器现量
   `go list -deps -test ./internal/winsec/ | grep -c internal/risk` → **0**；
   `go list -deps -test ./internal/config/ | grep -c internal/risk` → **1**。
   ⇒ winsec 的 POSIX 测试二进制**确实不链 risk**，把钉子钉在那儿就只有类型、没有安装事件。
3. **票面 AC#1 的前提为真**（"今天唯一断言在 `resolve_windows_test.go`，POSIX 那枚是 Skip"）：
   从**纯净快照** `wisp125r1-gate` 读 `internal/winsec/ancestor_separator_108_other_test.go`：
   ```
   func TestAC4POSIXFloorAnswersInsideTheNamedTree(t *testing.T) {
       if prev := winsec.PathResolverInstalled(); prev != nil {
           t.Skipf("this platform's binary links a real C26 pipeline (seam holds %T), so the floor leg measures nothing", prev)
       }
   ```
   ⇒ 修前 POSIX 侧唯一碰 `PathResolverInstalled()` 的那枚，正是"缝里有东西就 Skip"，**它不可能红**。
   ⚠ 这一发我特意**不**从工作树读（工作树里同名文件已被 `worker-ticket137-ac4` 改过，读它就等于验错版本）。

### 1.3 AC#1 结论

票面 AC#1 要求的两件事——**POSIX 有一枚 `PathResolverInstalled() == risk.c26Pipeline` 的钉子**、
**并自证它真的会红**——本程用**自己造的变异**独立复现成立。
**判：成立。** 唯一要登记的口径瑕疵（不影响结论）：票面 AC#1 说"今天唯一断言在 windows"，
而"POSIX 有断言"这一格是**这张票自己交的件**（`ff3faf9`）；本表读的是交件之后的树，
所以 §1 的每一枚"红"都是**交件套件上的红**，不是票 125 之前树上的红。这一点实现方在票面写清了，本表确认它没夸大。
