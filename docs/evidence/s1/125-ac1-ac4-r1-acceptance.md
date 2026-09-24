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

### 1.4 §0＋§1 正文所在 commit 的回执（这枚回执由 §2 那一枚 commit 携带）

```
$ git log --oneline -1
6c0ad3a evidence(125,AC#1 r1): 第一份裁决表 §0 锚点/闸门/被验版本身份 + §1 AC#1 成立
$ git show --name-only HEAD
commit 6c0ad3af63b7882963c036d5143d7662dccb3e2e
Author: CarlosShao <1933942520@qq.com>
Date:   Thu Sep 24 15:39:34 2026 +0800

    evidence(125,AC#1 r1): 第一份裁决表 §0 锚点/闸门/被验版本身份 + §1 AC#1 成立

    非实现者 acceptor-ticket125-r1。AC#1 用自造的 MUTATION-R1A（让安装被跳过）在容器里
    独立复现红名 TestAC1POSIXSeamHoldsC26Pipeline125，RUN=153 两侧同、panic=0。
    零生产码改动；全部变异在仓外快照 D:/tmp/wisp125r1-*。

docs/evidence/s1/125-ac1-ac4-r1-acceptance.md
```

---

## §2 AC#2 —— "守门人拿未解析的 `os.TempDir()` 造探针"该不该改、改得对不对

**判**：**成立附条件**。技术面（裁定方向＋判据可红＋拒绝侧没放宽）**成立**；
授权面**未闭合**，条件见 2.4。证据档：**〔独立复现〕**（三台对照与两发变异都是本程自己跑的）。

### 2.1 三台对照（本程自己跑，容器 `golang:1.27` 真跑，每把先 `ls -l go.mod` = 883 B）

构造法：`pre` ＝ `git archive ff3faf9`（**未修的生产码**）＋ 从 `4824bb8` 拷来的**两枚交件用例**；
`post` ＝ `git archive 4824bb8`。命令都是 `go vet ./internal/winsec/ ./internal/config/` 先证落地、再 `go test -count=1 -v` 同两包。

| 台 | VET_RC | TEST_RC | 读数（顶层/子用例） | panic | 红名 |
|---|---|---|---|---|---|
| `pre`（未修码＋新 leg） | 0 | **1** | `RUN=153 顶层PASS=82 顶层FAIL=3 子PASS=65 子FAIL=3 SKIP=0` | 0 | `TestAC2POSIXSeamProbeShapesAreBuiltOnAResolvedRoot125/measured_symlink_spelled_temp`、`TestAC2POSIXSeamAcceptsTheHonestPOSIXAnswer125/measured_symlink_spelled_temp`、`TestAC2POSIXSeamInstallSurvivesASymlinkSpelledTemp125/measured_symlink_spelled_temp` |
| `post`（本改动） | 0 | **0** | `RUN=153 顶层PASS=85 顶层FAIL=0 子PASS=68 子FAIL=0 SKIP=0` | 0 | 无 |

三枚红名**全部只红在 `measured_symlink_spelled_temp` 这一子用例**，四枚 `control_*` 子用例两台全绿 ⇒
红的是"这一形"，不是仪器自己坏了。

**AC#1 那枚钉子在同一形上的四枚读数（一票两形，本程现量）**：

| 形 | `pre` | `post` |
|---|---|---|
| plain temp | `--- PASS` | `--- PASS` |
| 软链拼写的 temp | **`--- FAIL`** | `--- PASS` |
| `PathResolverInstalled()` | 软链形 `<nil>` ⇒ 修后 `risk.c26Pipeline`；plain 形两台都 `risk.c26Pipeline`（红因/现场行见下） | |

现场行（逐字，路径里的 `<BASE>` 是本程替掉容器随机临时前缀，形状未改）：
```
pre  : --- FAIL: TestAC2POSIXSeamInstallSurvivesASymlinkSpelledTemp125/measured_symlink_spelled_temp
       child: c26_seam_posix_125_test.go:58: AC#1 RED: winsec.PathResolverInstalled() = <nil> …
post : --- PASS: …/measured_symlink_spelled_temp
       child: INFO winsec: sealing path resolver installed resolver=risk.c26Pipeline probes_passed=1
```

**探针材料的形状读数（这一行就是"改了什么"的最短证明）**：
```
pre  TMPDIR="<BASE>/varlink125/tmproot125" shapes=["<BASE>/varlink125/tmproot125/../wisp-103-conformance-probe"]
post TMPDIR="<BASE>/varlink125/tmproot125" shapes=["<BASE>/real125/tmproot125/../wisp-103-conformance-probe"]
```
⇒ **敌意 leg 一个字没少**（`/../wisp-103-conformance-probe` 还在），换掉的只有它站的那枚根；
plain 形两台的 `shapes` 逐字相同（`real125/…`）⇒ 这台机器上解析是恒等操作，改动没碰它。

### 2.2 本程另造的两发变异（各自红在哪一行）

**MUTATION-R1B（把"解析"这一刀摘掉，生产码其余一字未动）**，快照 `D:/tmp/wisp125r1-mutRoot`：
```diff
@@ -271,7 +271,7 @@
 		missing = append(missing, filepath.Base(cur))
 		cur = parent
 	}
-	real, err := filepath.EvalSymlinks(cur)
+	real, err := cur, error(nil) // MUTATION-R1B the probe root is no longer resolved
 	if err != nil {
 		return path
 	}
```
读数：`VET_RC=0`、`TEST_RC=1`、`RUN=153 顶层FAIL=3 子FAIL=3 PANIC=0`，红名与 `pre` 台**逐字同名**。
红因原文：
```
seam_probe_root_125_other_test.go:134: AC#2 RED: the seam guard built probe "<BASE>/varlink125/tmproot125/../wisp-103-conformance-probe"
                                       on a root that reaches itself through the link at /tmp…
seam_probe_root_125_other_test.go:295: AC#2 RED: with the process's temp dir spelled through a symlink the seam guard
                                       still refuses a candidate whose only sin is answering honestly about a root nobody resolved
```
⇒ **AC#2 的判据不是"修完必须红/绿"的装饰**：把这一刀摘掉，今天它就响，且响的名字与未修版一模一样。

**MUTATION-R1C（删掉守门人对候选答案的底线复算）**，快照 `D:/tmp/wisp125r1-mutFloor`：
```diff
@@ -297,7 +297,7 @@
 		if out == probe {
 			return fmt.Sprintf("it answered the hostile shape %q by passing it through unchanged, …", probe)
 		}
-		if _, fErr := (builtinVerifier{}).Resolve(out); fErr != nil {
+		if _, fErr := (builtinVerifier{}).Resolve(out); fErr != nil && false { // MUTATION-R1C the floor re-run leg is removed
 			return fmt.Sprintf("it answered %q with %q, a spelling the built-in floor itself refuses: %v", probe, out, fErr)
 		}
```
读数：`VET_RC=0`、`TEST_RC=1`、`RUN=153 顶层PASS=84 顶层FAIL=1 子PASS=66 子FAIL=2 PANIC=0`：
```
--- FAIL: TestAC2POSIXSeamGuardStillRefusesEveryHostileShape125
    --- FAIL: …/control_plain_temp
    --- FAIL: …/measured_symlink_spelled_temp
seam_probe_root_125_other_test.go:253: AC#2 RED (the refusal side moved): the seam guard accepted
                                       answer_through_the_link …
```
⇒ 这格的用途要说准：**leg 3 在 `pre` 与 `post` 两台都是绿的**（它按设计不是"修复判据"，是"边界钉"），
它的可红性**只由这发变异供给**。本程自己造出来并证实了：底线复算那一支一旦拿掉，一枚"答案穿过链接"的候选
在**两形**下都被放行。⇒ 放行侧的牙齿是真的，而票 125 没动它。

### 2.3 三层裁定（按派单要求，逐层分开、不混）

#### ① 事实层：那 82 行到底改了什么

- **改了什么**：两枚探针根的来源（`resolverProbeShapes()` 里那枚绝对拼写、`resolverTreeOwnershipFailure()` 里的 `parent`），
  从"直接拿 `os.TempDir()`"换成"拿 `resolverProbeRoot()`"；外加两枚新函数（`resolverProbeRoot` / `resolveProbeRoot`）
  与 121 行注释。**§0.5 已把全部非注释增删行贴出**，可逐字节重放。
- **`sameTree` / `sameVolume` / `foldSegment` / `answerInsideTree` 有没有被放宽**：**没有**，且这是**读 diff 原文**得出的，
  不是"看它没动就算清白"——把 82 行的 `+/-` 行喂 `grep -iE 'sameTree|sameVolume|foldSegment|answerInsideTree|resolverConformanceFailure|treeOwnershipFailureForPair'`
  现量**零命中**。`resolverConformanceFailure` 的四枚拒绝 leg（`no_account` / `pass_through` / `constant`（同一枚拼写）/
  **底线复算**）**函数体一字未改**，且最后那一支经 MUT-R1C 证明它**活着且有牙齿**（不是死代码）。
- **底线判据有没有被放宽**：**没有**。`platformVerifyPlacement`（拒未解析拼写的那位）一字未动；
  2.1 的 `shapes` 读数证明交到底线手里的**材料形状**仍带 `/../` 敌意。
- **它有没有把"未解析那一形"从此变成"永远在跑"**：**反过来了，且方向是对的**。
  修后"未解析那一形"只在**三枚失败出口**出现（`Lstat` 报的不是 `ErrNotExist`／`EvalSymlinks` 失败／走到根仍失败），
  每一枚都 `return path` ＝ **退回今天的形状**。⇒ 探针**不会因为解析失败而被跳过**（这是这格最容易长出一枚假绿的地方，
  代码方向选对了：最坏情况是回到本票要修的缺陷，而不是回到"什么都没测"）。
- **一处必须点名的**事实后果（实现方自己登记了 `R-125-1`，本表确认它是真的、不是过度自罚）：
  解析之后 install-time 的树归属 pair 变成"同一枚根的两枚拼写"，**真管线那一发永远走不到第二证人腿**
  （票 126 加的那条腿只在植桩 fake 下有分母）。⇒ 以后任何"跨根/跨卷"判据不许把"环境探针在跑"读成"这一腿被覆盖"。

**"自己算自己"这一刀（`R-119-9` 同族）——逐枚查**：
- AC#1 钉子：期望是字面量 `"risk.c26Pipeline"` ⇒ **不是**拿被测函数算的 ✓
- leg 1：期望 `planted.real` 由测试自己 `filepath.Join(base,"real125","tmproot125")` 造出来的；
  取实际值用的 `probeRootOf125()` 是对 `resolverProbeShapes()` 输出的**文本反解**，`firstLinkInPath125()` 问的是 `os.Lstat` ⇒
  **不调用 `resolveProbeRoot`** ✓（这一枚正是 `MUT-119-R` 恒真用例的形状，它避开了）
- leg 2 的候选 `lexicalFoldResolver125` 是**手写**的 `filepath.Clean(in)`。
  ⚠ 要点名：**"这就是 `internal/risk` POSIX 管线的答法"这句话，在 `internal/winsec` 这枚测试里没有凭据**
  （那包里进不去 risk，见 §1.2）。真管线那一腿的凭据在**别处**＝AC#1 的子进程 leg（本程已独立证实：软链形 `<nil>`→`risk.c26Pipeline`）。
  ⇒ 结论：不是"自己算自己"，但 leg 2 单独看**只裁守门人的反应**，不裁真管线的反应；本表按这个口径接受它，
  读者不许把它读成后者。文件头的注释（"fakes, and fakes on purpose"）说的是实话。

**技术裁定**：把"守门人拿未解析的 `os.TempDir()` 造探针"**判为误伤**，本表**独立同意**，理由不是它的注释，
是 2.1 那两行 `shapes` 读数：**候选做对的唯一一件事，是诚实回答一枚守门人自己没解析过的根**，
代价是整条 C26 掉到内置底线、而进程照常 rc=0。改法方向与票 119 已批准的纪律同侧（OS 给的答案由问它的那层解析），
且**没有**把底线反过来依赖上面那层（没调 `proc.SealableRoot`，方向对：那会拆掉 `envfork.go` 的边界话）。

#### ② 纪律层：**"别的票刚动过这枚文件"能不能构成"它不再是冻结件"的根据**

**正面判：不能。** 三条理由，从弱到强：

1. **冻结的单位不是文件，是"这张票被批准动哪几行"。**
   `a701138`（票 126，19:59）动的是 `sameTree`/`answerInsideTree`；票 125 动的是同文件里**另外两枚函数**。
   "这文件今天被人碰过"对"你这一刀批不批"提供的信息量是 **0**。
2. **这条理由一旦成立，冻结清单会自己缩短，而且缩短的时刻不由人决定。**
   同一枚文件会在 19:00 是冻结件、20:00 变成非冻结件，中间**没有任何一次人的决定**——
   变的是"谁刚提交过"。禁改列的全部意义在于"不管你现在多方便都不许动"，
   用它自己的例外来解释它，等于把它取消。
3. **`resolve.go` 的"冻结"身份在本仓只有一处来源，而那一处是编排者写给这张票的。**
   现量：`grep -rn 'winsec' docs/specs/*.md | grep -i '冻结|frozen|禁改'` → **零命中**；
   `.scratch/wisp/issues/README.md`「Hard global constraints」的冻结区是 `internal/risk/`（`:141`）与 SPEC-08 文本，**不含 `internal/winsec/**`**。
   ⇒ 票面 `:23` 那句"`internal/winsec/resolve.go` 与 `internal/risk/**` 都在冻结/禁改列"**不是仓级规矩，是这一票的范围条款**，
   范围条款只能由立它的人撤。派单如果真说了那句，那**授权是有的、只是给错了理由**；
   而如果它只是把"126 动过"当成依据，那**授权从来没被给出过**——这两种情况的处置不同，见 ③。

配套事实（本仓先例，不引别的）：
- 票 124（同一族、编排者同期写的 Rules）`:34` 逐字：**「若必须动生产码，停下来报给我，不要自己解冻。」**
- 票 125 自己 AC#2 `:24` 逐字：**「停在「裁完交回」是合格交件，不是失败」**。
  ⇒ 合规的那条路**当时就写着、而且明确不罚款**。它没走。
- 同族先例：**票 130 超授权改票 131 的判据，处置是"追认取证"，不是"默认放行"**。

#### ③ 可核层：**"派单里批过"这句话本身可不可验证**

**不可验证。** 三条盘上事实：

1. `.scratch/wisp/` 现量只有一枚 `issues/` 目录 ⇒ **仓里没有派单归档制度**，那句被引用的派单原文**不在盘上任何地方**。
2. 台账现量 `grep -n 'resolve\.go' docs/reports/pending-and-issues.md` = **9 行，无一枚是解冻授权**
   （9 行分别属于票 113/112/126/125 的机制描述与 A160 本案登记本身）。
3. **时间顺序**（本表新量到的一枚，比"台账里没有"更硬）：
   与代码同一枚 commit `4824bb8` 的 AC#2 日志正文（票面 `:86-161`）里
   `grep -nE '冻结|解冻|授权|禁改'` 现量 **零命中**，commit message 也没有；
   "派单更新了"这句话第一次出现在盘上，是 **`bd50c63`（09-22 21:20，补记二）—— 比代码落地（20:58）晚 22 分钟**。
   ⇒ 这不是"当时报备、事后落账"，是**事后追述**。

**这一格记在谁身上——两栏分开记，不互相抵账：**

| 记在 | 内容 | 轻重 |
|---|---|---|
| **编排者（本表认为更该修的那一枚）** | 具名解冻没落 `A##`。"任何具名解冻必须在台账落一枚 `A##`（文件/行/理由/边界/撤销口令）"这条规矩缺的正是这一枚。**没有它，任何一次越界事后都能补一句"派单里批过"，而且无人能证伪。** | 制度洞，影响后面每一张票 |
| **实现方（轻，但真实）** | (i) 动码那一刻没在自己那格日志里标出"这一步依赖一条与票面 `:23` 不同的授权"；(ii) **它写下来的理由是错的那一枚**（"这文件不是冻结件"），而不是唯一能成立的那一枚（"编排者在派单里解冻了它"）。错理由比没理由更贵，因为它会被下一张票照抄。 | 归因/纪律瑕疵，不是伪造 |

**若要追认，还缺的凭据（一条，形状写死在这里）**：一枚**由编排者署名**的 `A##`，五字段齐全——
① 哪枚文件：`internal/winsec/resolve.go`；② 哪几行：`resolverProbeShapes()` 首枚 shape + `resolverTreeOwnershipFailure()` 的 `parent` + 新增两枚函数；
③ 为什么：守门人把 OS 自己的合法形状当攻击 ⇒ 整条 C26 掉底线；④ 边界：`internal/risk/**` 仍冻结（本票 `git show --name-only` 现量零枚 risk 路径，此点成立）；
⑤ 撤销口令。
**可选加强**：给派单建一个归档位置（现在结构上不可核，这一条不针对本票）。

**本表不下处置结论**：不撤、不 revert、不改任何东西——**处置权在编排者**。
本表只把三层裁清：**技术面成立、纪律面不成立、可核面缺凭据** ⇒ AC#2 记 **成立附条件**，
条件是 ③ 那枚 `A##`（**或**把生产码那一刀退回成"只出裁定与判据"，判据与本表 §2.1/§2.2 三枚 leg 可以原样留在仓里）。

### 2.4 §2 正文所在 commit 的回执（这枚回执由 §3 那一枚 commit 携带）

```
$ git log --oneline -1
20a6397 evidence(125,AC#2 r1 终裁 §2): 三台对照 + 自造 MUT-R1B/R1C 各自红点 + 冻结件三层裁定
$ git show --name-only HEAD
commit 20a63979cc88018b44597ba483246219a72abe20
Author: CarlosShao <1933942520@qq.com>
Date:   Thu Sep 24 15:42:00 2026 +0800

    evidence(125,AC#2 r1 终裁 §2): 三台对照 + 自造 MUT-R1B/R1C 各自红点 + 冻结件三层裁定

    技术面成立（判据可红、拒绝侧未放宽、探针失败方向＝退回未解析而非跳过）；
    纪律面不成立（"别的票刚动过这枚文件"不构成解冻，冻结的单位是这一票批准的行）；
    可核面缺凭据（派单不在盘上、台账零枚 A##，且"授权"一句比代码落地晚 22 分钟才出现在盘上）。
    AC#2 判成立附条件＝补一枚编排者署名的 A##，处置权在编排者。

docs/evidence/s1/125-ac1-ac4-r1-acceptance.md
```

---

## §3 AC#3 —— 可见性：那行 `ERROR` 到不到得了人眼前

**判**：**成立**（这一格要的是"复算＋写清哪一格该红"，实现方两件都做了，本程独立复现同向）。
证据档：**〔独立复现〕**（真二进制是本程自己在容器里 build 的，盘上读数是本程自己 grep 的）。

### 3.1 本程的真二进制复算（容器 `golang:1.27` 真跑，两形两台）

| 台 | 二进制 | `RUN_RC` | stderr 第 1 行 | 盘上 JSONL | `grep -rl winsec <植的树>` |
|---|---|---|---|---|---|
| **改前** `ff3faf9` | 19,092,792 B（`go build ./cmd/wisp` rc=0） | 2 | **`ERROR winsec: refusing to install a path resolver into the sealing seam resolver=risk.c26Pipeline …`** | 存在，2 行 | **`NONE`** |
| **改后** `bd50c63`（码＝`4824bb8`） | 19,097,696 B（rc=0） | 2 | `INFO winsec: sealing path resolver installed resolver=risk.c26Pipeline probes_passed=1` | 存在，2 行 | **`NONE`** |

形状先断言过真是链接（`test -L $b/varlink125`，`readlink` → `$b/real125`）；两把 `TMPDIR` 都是 `$b/varlink125/tmproot125`。

改前那一行的**全文**（逐字，含它拒绝的理由，这是本票的现场读数）：
```
2026-09-24 07:29:55 ERROR winsec: refusing to install a path resolver into the sealing seam
  resolver=risk.c26Pipeline reason="it answered \"<B>/varlink125/tmproot125/../wisp-103-conformance-probe\" with
  \"<B>/varlink125/wisp-103-conformance-probe\", a spelling the built-in floor itself refuses:
  winsec: path is not provably resolved, refusing to seal: <B>/varlink125/wisp-103-conformance-probe"
```
紧随其后才是持久 sink 自己那行（**顺序就是答案**）：
```
time=2026-09-24T07:29:56.030Z level=INFO msg="wisp: persistent log sink installed" dir=<B>/real125/tmproot125/wisp-test-3164/logs min_level=info
```
盘上那枚 JSONL（`…/wisp-test-3164/logs/wisp-20260924-001.jsonl`）**全文 2 行**：一条 `wisp: persistent log sink installed`、
一条 `audit: perm: MODE-READ-FAILED …`；`grep -c winsec` = **0**。改后那台同样 **0**。

⇒ 实现方那枚读数（容器内 `grep -rl winsec <data root>` = `NONE`）**本程独立复现，两形两台都是 `NONE`**。

### 3.2 机制（本程自己量的调用点，不是抄注释）

- 记录发在 `winsec.SetPathResolver` 里；生产侧唯一的调用方是 `internal/risk/winsec_c26.go` 的 **`init()`**（任何 `main()` 之前）。
- 听众（`installLogSink`）在盘上的**全部**生产安装点：`cmd/wisp/run.go:164`、`cmd/wisp/resident_windows.go:57`（`cmd/wisp/models.go` 走同一族），
  现量 `grep -rn 'installLogSink' cmd/wisp/*.go` ⇒ **零枚在 `init()` 里**。
- 顺带补一枚实现方没写的边界事实：`cmd/wisp/main.go` 里那 9 处 `attachParentConsole()`（含无参数＝GUI/常驻入口那一处）**也都发生在 `main()` 之内**
  ⇒ 那条记录写的是 **rebind 之前**的 `os.Stderr`。⇒ "换听众之前没人听"这一族里，本格比票 127 的 `logging.go:204` 更早，**这一点本表确认**。

### 3.3 具名落点（派单要求的"落进 `R-125-3`／票 130 那一族的哪一格"）

- 台账现量：`R-125-3` 已在 `docs/reports/pending-and-issues.md:3294-3297`（`A99③`，第八次）落地，
  且 `:3353` 写明 **新建＝票 130（`R-117-2` + `R-125-3` 合一，需 `internal/risk` 解冻）**。
- ⇒ AC#3 的残余缺口**不属于票 125**，它落在 **票 130 的第一格：init-time 安全记录的可达性**。
  票 125 该红的部分（"POSIX 侧没人断言这条降级"）已经被 AC#1 的钉子接住（见 §1：软链形改前红、改后绿）。
- 实现方在本格给的判据形状（`init()` 里走一次拒绝 → 装 sink → 断言盘上看得到那条）本表**认可它是可写的**，
  并且它今天必然红——本程 3.1 两把读数就是它的红因。**本格零码改动这条也成立**：
  `git show --name-only a03f7ef` 现量只列票面一枚 `.md`。

### 3.4 要打折的一处措辞（结论不变）

票面 AC#3 `:191-192`：「它唯一的去处是 stderr，而 `wisp run` 的 stderr 在无终端的入口（resident/GUI 双击）**没人接**」。
本程现量：仓里 `grep -rn 'windowsgui\|-ldflags' scripts/*.sh` **零命中** ⇒ 交付的是 **console-subsystem** 构建，
双击起来的 console 窗口**在进程存活期间是能滚到那行的**。
⇒ **"等于无记录"对"持久、可回溯"这一档成立（3.1 的 `NONE` 就是证据）；对"当场有没有一扇窗口"这一档字面过强**。
按 `A98④` 的口径处理：**结论签、措辞打折**。⚠ Windows 双击那一形**本程没测**（见 §6），这条只是代码级归因，不是读数。

### 3.4 附：实现方留在盘上的归档抽验（这一档是〔日志＋归档，抽验〕）

本程只读地查了实现方声称留下的快照——它们在**宿主 temp**（`/tmp/wisp-t125-*` ＝ `C:\Users\swq\AppData\Local\Temp\`），
不在 `/d/tmp/`（那里现量 0 枚）：
```
wisp-t125-ac1  -blast  -ctrl  -ctrl81b  -final  -final2  -mutA  -mutA2  -mutB  -post  -pre  -s125   （12 枚，全在）
```
抽验两枚关键变异的**落地原文与本程独立造的那发同侧**：
```
wisp-t125-mutA2/internal/winsec/resolve.go:134: reason = resolverConformanceFailure(r) + " MUTATION-125A2 the install is refused on every candidate"
wisp-t125-mutB /internal/winsec/resolve.go:300: _ = out // MUTATION-125B the guard no longer re-runs the floor on the candidate answer
wisp-t125-pre  /internal/winsec/resolve.go: grep -c resolverProbeRoot = 0   （确认是未修版）
```
⇒ 票面 `:115-117` 那张三台表里的**变异原文与行号逐字对得上**（`resolve.go:134`／`:300`），
本程据此把 §2 的实现方那一档从〔仅自述，不背书〕**提升为〔日志＋归档，抽验〕**。
但**它的读数值（`RUN=14` 那组）本程没有复用到**——§2 的三台数全部是本程自己 `-count=1` 现跑的。

### 3.5 §3 正文所在 commit 的回执（这枚回执由 §4 那一枚 commit 携带）

```
$ git log --oneline -1
a088515 evidence(125,AC#3 r1 终裁 §3): 真二进制两形两台独立复现——盘上 winsec 零命中，NONE 成立
$ git show --name-only HEAD
commit a08851552323c6ed747c60f0ae1277f8b48aa8c5
Author: CarlosShao <1933942520@qq.com>
Date:   Thu Sep 24 15:44:30 2026 +0800

    evidence(125,AC#3 r1 终裁 §3): 真二进制两形两台独立复现——盘上 winsec 零命中，NONE 成立

    自建 wisp-pre(19,092,792B) / wisp-gate(19,097,696B) 两枚容器二进制，软链 temp 形下
    改前 stderr 第 1 行就是那枚 ERROR、sink 在其后才装上；JSONL 2 行、grep -c winsec=0、
    grep -rl=NONE（两台都是）。残余缺口具名落到票 130 第一格（init-time 记录可达性）。
    一处措辞打折：交付的是 console-subsystem 构建，"等于无记录"只对可回溯档成立。

docs/evidence/s1/125-ac1-ac4-r1-acceptance.md
```

---

## §4 AC#4 —— 门禁

**判**：**成立**（四数、卫生、双 GOOS、d22scan＋台账逐 scope 都独立复现；两处口径要钉，见 4.2／4.4）。
证据档：**〔独立复现〕**。

### 4.1 受影响包 `-count=2 -v` 四数（容器真跑；数只从 `-v` 取）

树：`bd50c63`（代码＝`4824bb8`，见 §0.4）vs 控制组 `81b4d5f`。每把先 `ls -l /<mnt>/go.mod` = 883 B、`uname -s`=Linux、`go1.27.1 linux/amd64`。

| 包 | `bd50c63` | 控制组 `81b4d5f` | rc |
|---|---|---|---|
| `./internal/winsec/` | `RUN=104 顶层PASS=60 顶层FAIL=0 顶层SKIP=0 / 子PASS=44 子FAIL=0` | `RUN=84 顶层PASS=54 … / 子PASS=30` | 两把 rc=**0** |
| `./internal/config/` | `RUN=202 顶层PASS=110 顶层FAIL=0 顶层SKIP=0 / 子PASS=92 子FAIL=0` | `RUN=194 顶层PASS=106 … / 子PASS=88` | 两把 rc=**0** |
| 两包合并跑 | `RUN=306 顶层PASS=170 顶层FAIL=0 顶层SKIP=0 / 子PASS=136 子FAIL=0` | `RUN=278 顶层PASS=160 / 子PASS=118` | 两把 rc=**0** |
| `panic`／`fatal error` 计数 | **0**（三把全 0） | **0** | — |

⇒ 实现方 AC#4 表那四组数（`104/60`、`202/110`、`306/170`、控制组 `84/54`、`194/106`）本程**逐字复现**。

**⚠ 必须钉的口径（否则下一位会把两种分母加成一个数）**：
实现方这张表的 `RUN` 数的是**所有** `=== RUN` 行（含子用例），`PASS` 只数**顶层** `--- PASS` 行 ⇒
`RUN=104 / PASS=60` **不是**"44 条没过"，那 44 枚是子用例（`子PASS=44`，顶层+子＝104＝RUN）。
本表把两档分开写全。**另：本表 §2 用的是 `-count=1`，本节是 `-count=2`，两节不许互比。**

**名册差集（防止"一条 panic 吞掉几十条读数"）**：
顶层＋子用例名集去重后 `bd50c63` **153 枚** vs `81b4d5f` **139 枚** ⇒ **＋14、丢 0**。
`comm -13` 现量**空**（没有任何一枚名字从名册上消失）。
＋14 逐名可点到本票：`TestAC1POSIXSeamHoldsC26Pipeline125`、
`TestAC2POSIXSeamProbeShapesAreBuiltOnAResolvedRoot125`、`TestAC2POSIXSeamAcceptsTheHonestPOSIXAnswer125`、
`TestAC2POSIXSeamGuardStillRefusesEveryHostileShape125`、`TestAC2POSIXSeamInstallSurvivesASymlinkSpelledTemp125`
五枚顶层 ＋ 九枚子用例。⇒ 增量对得上，不是"分母做小换来的绿"。

### 4.2 `gofmt` / `gofumpt`（写明版本；CI 那步没钉版本）

| 仪器 | 版本读数 | 范围 | 结果 |
|---|---|---|---|
| `gofmt -l internal cmd` | go1.27.1（容器） | 全 `internal/`+`cmd/` | **空**（0 行） |
| `gofmt -l internal cmd` | go1.27.1（宿主 windows/amd64） | 同上 | **空** |
| `gofumpt -l internal cmd` | **`v0.12.0 (go1.27.1)`**（`D:\work\base\gopath\bin\gofumpt.exe`，`--version` 现量） | 同上 | **空** |

⚠ 三点具名：
1. **实现方报的版本是 `v0.7.0`，本程现量这台机器上是 `v0.12.0`**（exe mtime `09-23 22:23`，比票 125 落地晚一天）。
   ⇒ 它那把读数**在今天这台机器上已经不可复现**（工具换了），本表的清白读数**只属于 v0.12.0**。
2. CI 那步逐字是 `.github/workflows/ci.yml:113 go install mvdan.cc/gofumpt@latest` ⇒ **这枚门禁按构造就没有固定版本**。
   两枚版本（v0.7.0／v0.12.0）下都是空的，所以结论稳；但"版本"这一栏**在这张票上不可能长期有效**，
   这是门禁自身的性质，不是票 125 的缺陷（登记在本表 §6 的"没核"边里：没去核 v0.7.0↔v0.12.0 之间是否有规则差异）。
3. 容器内没有 `gofumpt`（未 `docker pull`、未 `go install`，避免动网络）⇒ **gofumpt 只有宿主一把**，与实现方同形，如实记。

### 4.3 `go vet` 双 GOOS ＋ **逐错误行归因**

| 把 | 命令 | rc | 归因 |
|---|---|---|---|
| 宿主原生 | `GOOS=windows go vet ./internal/winsec/ ./internal/config/` | **0** | 无 |
| 宿主交叉 | `CGO_ENABLED=0 GOOS=linux go vet ./internal/winsec/ ./internal/config/` | **0** | 无（**只编译不执行**，正向读数只来自 §2/§4.1 的容器真跑） |
| 宿主交叉·全模块 | `CGO_ENABLED=0 GOOS=linux go vet ./...` | **1** | **1 枚错误、0 枚 `file:line`**：`package github.com/CarlosShao/wisp/cmd/wisp` → `imports …/sherpa_onnx` → `imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in …` ⇒ 指向的是**外部模块目录**，不是本仓任何源文件行 |
| **容器·真平台全模块** | `go vet ./...`（linux/amd64，`CGO_ENABLED=1`） | **0**，输出 0 行 | **这一把把上面那枚 rc=1 判掉了** |
| 容器·关 cgo | `CGO_ENABLED=0 go vet ./cmd/wisp/` | **1** | 与宿主交叉那一把**逐字同文**（只有模块目录前缀不同）⇒ 与"交叉编译"无关，纯粹是 cgo 关掉后那枚包的 build constraint |

**归因结论（本程比"既不算破口也不算清白"多走了一步）**：
- 它**不是破口**：同一枚命令在**控制组 `81b4d5f`（票 125 之前的树）**上现量**同样 rc=1、同样这一行**
  ⇒ 票 125 既没造成也没掩盖它；且真平台 `go vet ./...` rc=0。
- 它**也不需要记成"清白未知"**：本程拿到了真平台（容器 cgo on）那枚 rc=0 的正向读数，
  所以派单里"宿主交叉那发既不算破口也不算清白"这一句，**在本表里被一枚实测定掉了**（不是靠推理）。
- 错误行里**没有**任何本仓路径 ⇒ 按"错误文本点到自家仓库路径才按真伤处理"的尺，这一发**不升真伤**。

### 4.4 `sh scripts/d22scan.sh` ＋ 台账各 scope 不降

纯净快照两棵各一把（容器），**两把 rc=0**，且**正控制那一步真跑了**：
`runtests.sh: OK - packages=[./...] top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31, '[no tests to run]'=0`（`TestBuiltBinaryGoesRedEndToEnd` 六枚子用例逐枚 PASS）。

| scope（`d22scan` 自己打的 examined 数） | `bd50c63` | `81b4d5f` | 判 |
|---|---|---|---|
| bans #1-5 `internal/` | 202 | 202 | 持平 |
| bans #1-5 `cmd/` | 22 | 22 | 持平 |
| ban #6 `frontend/` | 40 | 40 | 持平 |
| ban #7 `internal/tools/` | 18 | 18 | 持平 |
| ban #8 `design/` | 16 | 16 | 持平 |
| ban #8 `frontend/` | 40 | 40 | 持平 |
| **ban #8 `internal/`** | **385** | **383** | **＋2 ＝ 本票新增两枚 `_test.go`** |
| ban #8 `cmd/` | 32 | 32 | 持平 |

⇒ **八枚 scope 零枚下降**，唯一变化那一枚的增量**逐名可归因**（`internal/winsec/seam_probe_root_125_other_test.go` ＋
`internal/config/c26_seam_posix_125_test.go` 都是 `internal/` 下的 `.go`），不需要"逐名解释下降"（因为没有下降）。
实现方给的基线 383 与本程控制组现量**逐字吻合** ⇒ 它没把基线说高。

### 4.5 AC#4 里本表**不签**的一句

票面 AC#4 那格的 `②软链 temp 形状那一把（信息性）`读数（`RUN=144 PASS=62 FAIL=19 SKIP=4`）本程**没复算**
（它是票 124 那本 harness 账，不在本票射程；实现方自己标了"信息性"）。
⇒ 那一组数在本表里的档位是**〔仅自述，不背书〕**，见 §6。
