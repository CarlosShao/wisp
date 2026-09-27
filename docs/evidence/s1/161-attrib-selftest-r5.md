# 票 161 · 161-r5（写码位）＝归因尺的合成自测 ＋ CI 甲形那一步 ＋ `TRACKED-DIRTY` 腿的实测

- 派单：`.scratch/wisp/dispatches/2026-09-27-001x-impl-161-r5-attrib-selftest-and-ci-leg.md`
- 锚点（step 0 现量，未抄任何人的号）：`git rev-parse HEAD` = `d49a5bebc88b613507383112c9dc12459f368da7`｜分支 `dev`｜`date` = `Sun Sep 27 00:02:55 CST 2026`
- 本程写面：`probes/161/r5/**` ＋ 本件 ＋ `.github/workflows/ci.yml` 的一步 ＋ 票面 Progress log 一格。`tools/d22scan/**`、`allowlist.txt`、`probes/161/r3/**`、`r4/**`、契约轴全部**零字节**（证据见 §7）。
- 尺：`gofumpt v0.12.0 (go1.27.1)`，命令 `"$(go env GOPATH)/bin/gofumpt.exe" --version`。⚠ 裸 `gofumpt` 不在 PATH（`gofumpt --version` ⇒ `command not found`，rc=127），所以本程每一发都走 GOPATH 全路径或仪器的现跑。

---

## 0. step-0 四件

1. `date`｜`git rev-parse HEAD`｜见抬头。
2. `git status --porcelain -- .scratch/wisp/probes/161/ .github/workflows/ci.yml tools/d22scan/` ⇒ **非空**，两行，都是别人（161-r2）的未跟踪遗留：
   `?? .scratch/wisp/probes/161/r2/__pycache__/`｜`?? .scratch/wisp/probes/161/r2/ctl/`
   ⇒ 按派单"非空报回别动"处理：**未碰、未提交、未还原**（本程写面是 `r5/**`，与那两枚不相交；`git show --stat` 见 §8）。
3. 票面 AC#7 那一格**连行号**逐字抄录（`sed -n '35,43p' .scratch/wisp/issues/161-…-done 无 ⇒ 文件名见下）：

```
35  - [ ] **AC#7（09-26 22:0x 追加，来路＝161-r1 报回＋编排者独立复算）自家取证台件不许把全仓门点红，且"跑门"必须排在"落台件"之后**。
36    **现场（两向都有独立读数）**：`$(go env GOPATH)/bin/gofumpt.exe -l . tools/d22scan tools/mockllm`（版本现读 **`v0.12.0 (go1.27.1)`**）今天**不空**，
37    其中 **3 行是已入库件**＝`.scratch/wisp/probes/158/accept-r1/mut/guard.no1.go`／`no2.go`／`no3.go`（`506cbae`，20:49 由票 158 的**验收程自己**加进去的变异台件）。
38    ⇒ 时序很扎人：那枚验收表在 **20:2x** 现量过 `gofumpt -l . … rc=0 空` 并据此判 AC#5③ **成立**，**不到半小时后它自己加的台件让那条读数过期了**。
39    这不是谁读错，是**"门先跑、台件后落"这个顺序本身没有门**。
40    **本格要做两件事，都要读数**：① **止血**——把那 3 枚入库台件格式化到 `gofumpt` 干净，并**复跑 `-overlay` 那一发证明变异行为一字未变**（红名与红句必须与验收表 §2.1 逐枚对得上；只改空白不改语义这句话**要用读数证明，不许说**）；
41    ② **装顺序**——本仓任何一票的门禁步骤里，**全仓级那三把尺（`gofumpt -l .`／`sh scripts/d22scan.sh`／`go vet`）必须在台件全部入库之后跑最后一次**，并把这一条写成一个可复算的形状（一枚脚本或一票面判据都行，但必须**说得出删掉它哪一发会重新漏**）。
42    ⚠ **判据不许写成"台件永远不许被这些门点红"**——那是恒真检（本仓否过两次）。本格要的是**"门点红了有人能在推送之前知道是自己的台件造的"**这一发可复算。
43    ⚠ 与 169 的区别要说清：169 是**被扫面里的真违规**（界面上的字符），本格是**取证机械自己的产物撞门**（假文件、真红）。两件事，别并账。
```

（票面文件名：`.scratch/wisp/issues/161-gates-need-their-own-gates-self-tests-in-ci-plus-is-the-approval-door-ever-rung.md`；上面那段与 `sed -n '35,43p'` 的输出做过逐字比对，见 §5 末。）

4. 派单给我的三条前提，**两条与盘符、一条不符**（不符那条是本程最大的一个发现）：

| 派单原话 | 盘上 | 命令 |
|---|---|---|
| "stray sample 现在让乙形默认红" | **符**：`lines=10 files=8 attributed=[161] unattributable=1`、rc=1（搬样本之前是 `lines=9 files=7`） | `sh .scratch/wisp/probes/161/r4/attrib.sh` → `logs/a0-r4-baseline.txt` |
| "基线 34 枚通过" | **符**：`PASS=34 FAIL=0 SKIP=0`、`=== RUN=76`、rc=0 | `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` → `logs/g1-runtests-pre.log` |
| "甲形的失败路径从未真的红过" | **符**，两把不同形状的尺：① `grep -rl 'TRACKED-DIRTY' .scratch/wisp/probes/` ⇒ 只命中 `r4/attrib.sh`＋`r5/**`（源码与我本程的日志），**没有任何前程序的日志含这个 token**；② `grep -h 'tracked-dirty=' .scratch/wisp/probes/161/r4/logs/*.txt` ⇒ **3 发全是 `tracked-dirty=0`**，那三发的红一律来自 `unattributable=1`。r3 留下的 `logs/post-r*-m*no*.red.txt` 是**另一把尺**（`go test` 的 `RUN=116 PASS=78 FAIL=2`，摘样本造的红），不是"已跟踪字节不干净" | 两列命令即读数 |
| **"移动后乙形应当全部可归因、rc=0"** | **不符**，见 §2.2 | — |

---

## 1. 本程**没有**测什么（先说这个）

- **没测 CI 真会怎么跑。** 那一步只在本地以同形命令跑过（§3.3），`ubuntu-latest` 上的 `$(go env GOPATH)/bin/gofumpt`（无 `.exe`）路径由 `attrib.sh` 自己那三行 `[ -f/-x ]` 判定，**未在任何 runner 上验证过**。
- **没测 `--self-test` 之外的任何 Go 行为。** 本程零枚 `.go` 改动、零枚 Go 测试用例增删（四数 34→34、名册差集空，§4）。
- **没测"归因＝故意"**：`attrib.sh` 从头到尾只能说"票 NN 拥有这个路径"，说不出这些字节是不是 NN 自己写坏的。（这句是 r4 就登记的边界，本程原样继承，`classify_line` 一字未改那条规则。）
- **没测 `..` 之外的路径归一化**：`classify_line` 沿用 r4 的 `norm()`（只把反斜杠换成正斜杠、剥 `:行:列:` 尾巴），**不解析 `..`**。`.scratch/wisp/probes/161/r5/negative-control/../../../../999/…` 这种文本会被归到 **ticket-161**（形状正确、语义在撒谎）。真实 `gofumpt -l .` 不会产出这种行（它从目录树走上来拼路径），所以今天不是活洞；我一度把它写成自测的一发，跑出来发现期望写反了就**撤掉了那一发而不是改判据**，并把它登记在这里。会响条件：谁哪天让某把尺吐出带 `..` 的行，那一发就会把不该归因的行归掉。
- **没测那枚 `probes/999` 旧样本的去处**：退役它＝删除，删除不在我权限内（§8）。
- **AC#6（聚合 rc ＋ 两族成对普查）一格未碰**，AC#3/AC#4 也未碰（派单明令）。

---

## 2. 第①格：把"会响"从盘上挪进文本

### 2.1 复制（不是移动：本仓只建不删）

```
cp .scratch/wisp/probes/999/bad-sample.go .scratch/wisp/probes/161/r5/negative-control/bad-sample.go
sha256sum 两枚 ⇒ 8aa59c87a088247037141081cd308ad22d6fb4e1bc6c341f5a72b18eb69fc3c4  （两枚同一串）
```

**两处都在，同字节。** 新路径落在 `.scratch/wisp/probes/161/**` 且 `.scratch/wisp/issues/161-*.md` 真存在 ⇒ 归因规则的第一半＋第二半同时成立，乙形把它读成 `ticket-161`（下面 §2.3 的第 10 行）。

两枚都**不提交**，理由不是"临时件"而是**提交上去会把甲形点红**：那两枚字节是 gofumpt 拒绝的形状，一旦进了已跟踪集合，`gofmt (gofumpt)` 与新增那一步都会红（§3.2 实测的就是这条腿）。

### 2.2 派单那句"移动后乙形全部可归因、rc=0"在盘上不成立

| 读数（同一把尺，两份日志） | 命令 |
|---|---|
| 复制**前**：`(B) lines=9 files=7 tracked=0 attributed_tickets=[161] unattributable=1`、rc=1 | `sh .scratch/wisp/probes/161/r4/attrib.sh` → `logs/a0-r4-baseline.txt` |
| 复制**后**：`(B) lines=10 files=8 tracked=0 attributed_tickets=[161] unattributable=1`、rc=1 | `sh .scratch/wisp/probes/161/r5/attrib.sh` → `logs/c3-r5-both-forms-post-move.txt` |

多出来的那一行是新样本（可归因），**那一发归不出的仍然是 `probes/999/bad-sample.go` 原件**——它还在盘上，`gofumpt -l .` 就还看见它。"只建不删"＋"删除不在我权限内"⇒**乙形在这一棵树上仍默认红，直到 owner 把原件退役**。我没有为了让那句成立去动判据，也没有动那个文件（它在我写面之外）。

**因此本格的交件物是这两条，而不是一条：**
1. 会响这一发**已经不再需要那枚伤**（§2.3 的合成自测，两向，复跑即可）；
2. 那枚伤留在原地，是**一件待退役的事**，不是仪器的一部分。派单 §5 的"旧的由谁退役"——答案：**owner／编排者，实现程无权**。

### 2.3 合成自测：`classify`／`ticket_of`／`ticket_known` 喂文本，两向

r5 把 r4 只在乙形循环里内联做的判定抽成 `classify_line`（同一条规则、两个喂法：磁盘那一次和文本这一次），于是"喂一行文本 ⇒ 判 `UNATTRIBUTABLE` 并硬退出"成为可能。

`sh .scratch/wisp/probes/161/r5/attrib.sh --self-test` → `logs/c2-selftest-green.txt`，rc=0：

```
== attrib.sh --self-test: the classifier, fed as text (v0.12.0 (go1.27.1))
   PASS  want rc=1 verdict=UNATTRIBUTABLE    got rc=1  .scratch/wisp/probes/999/bad-sample.go (path claims ticket 999 but .scratch/wisp/issues/999-*.md does not exist, class=untracked, parse_err=0) <== REAL INJURY
   PASS  want rc=1 verdict=UNATTRIBUTABLE    got rc=1  .scratch/wisp/probes/999/bad-sample.go (... class=untracked, parse_err=1) <== REAL INJURY        <- 反斜杠＋:19:2: 诊断尾巴那一发
   PASS  want rc=1 verdict=UNATTRIBUTABLE    got rc=1  .scratch/wisp/probes/777/second-fake-ticket.go (... issues/777-*.md does not exist ...) <== REAL INJURY
   PASS  want rc=1 verdict=UNATTRIBUTABLE    got rc=1  internal/tools/not_a_bench_path.go (no probes/<NN>/ path segment, class=untracked, parse_err=0) <== REAL INJURY
   PASS  want rc=1 verdict=TRACKED<CI-RED>   got rc=1  tools/d22scan/main.go (parse_err=0) <== TRACKED BYTES ARE DIRTY
   PASS  want rc=0 verdict=ticket-161        got rc=0  .scratch/wisp/probes/161/r5/negative-control/bad-sample.go (untracked, parse_err=0)
   PASS  want rc=0 verdict=ticket-161        got rc=0  ...negative-control\bad-sample.go:19:2: expected declaration, found '\' (untracked, parse_err=1)
   PASS  want rc=0 verdict=ticket-169        got rc=0  .scratch/wisp/probes/169/a-bench-sample.go (untracked, parse_err=0)
attrib.sh: --self-test cases=8 failures=0
attrib.sh: SELF-TEST GREEN - the classifier rings on a line it cannot attribute and stays silent on a line a real ticket owns
```

**两向各 4／4：5 枚必须响（其中 1 枚是"已跟踪"那一支）、3 枚必须不响。** 不成对的那一半正是恒真检的形状：响的那五枚里放了 `777`（跟 `999` 不同的一枚假票号，证明"999 被特殊照顾"不成立）、`internal/tools/…`（形状就不对）、`tools/d22scan/main.go`（真的已跟踪文件，走的是另一支）；不响的那三枚里放了 `ticket-169`（**另一张真票**，证明归因不是"161 自己的目录永远过")。

单发（派单点名的那两发）：

```
$ sh .scratch/wisp/probes/161/r5/attrib.sh --classify-line '.scratch/wisp/probes/999/bad-sample.go'
UNATTRIBUTABLE	.scratch/wisp/probes/999/bad-sample.go (path claims ticket 999 but .scratch/wisp/issues/999-*.md does not exist, class=untracked, parse_err=0) <== REAL INJURY
$ echo $?   ->  1
$ sh .scratch/wisp/probes/161/r5/attrib.sh --classify-line '.scratch/wisp/probes/161/r5/negative-control/bad-sample.go'
ticket-161	.scratch/wisp/probes/161/r5/negative-control/bad-sample.go (untracked, parse_err=0)
$ echo $?   ->  0
```

**自测自己会响那一发（不是我说，是它红过一次）：** 第一版我把一枚期望写成 `ticket-158`（`probes/158/accept-r1/mut/guard.no1.go`），跑出来 **`--self-test` rc=1、`FAIL want rc=0 verdict=ticket-158 got rc=1 verdict=TRACKED<CI-RED>`**（全文 `logs/c1-selftest.txt`）。那枚文件在 `506cbae` 之后已被格式化并**已跟踪**，所以"它归 158"这句本来就错——**我的前提错，不是尺错**；我按 §5 的口径换掉那一发（换成 `169` 的文本一发）而没有动 `classify_line`。**同一把尺在两棵树上都不动判据。**

行为未漂移的机械证据：`diff logs/c4-r4-both-forms-post-move.txt logs/c3-r5-both-forms-post-move.txt` ⇒ **IDENTICAL**（r4 的原尺与 r5 的新尺，同一棵树、同一分钟，两形全部读数一字不差）。

---

## 3. 第②格：CI 只加甲形那一步

### 3.1 diff 形状（只插一步，没顶任何人）

`git diff --numstat -- .github/workflows/ci.yml` ⇒ `28	0	.github/workflows/ci.yml`（**加 28 行、删 0 行**）。既有步骤名册顺序（`python -c "yaml.safe_load(...)"` 现读，不是 grep 出来的）：

```
lint steps: 11
 - D22 scanner positive control (tools/d22scan tests, seeded red)
 - D22 scanner self-test (ticket 161 AC#2 - every ban, both directions)
 - D22 seven-ban + emoji scan (tools/d22scan)
 - gofmt (gofumpt)
 - gofmt (gofumpt) - the tracked set is the denominator (ticket 161 AC#7 form A)   <- 新增，就在上一行之后
 - go vet (module)
 - go vet (tools/d22scan module)
 - staticcheck
 - mockllm module vet
```
（前两枚 unnamed ＝ `uses: checkout`／`uses: setup-go`。）

新步骤解析出来的键 ⇒ `keys: ['name', 'run']`——**没有 `if:`、没有 `continue-on-error`、没有 `working-directory`**。`git diff -U0 | grep '^+' | grep -nE 'if:|continue-on-error'` 只命中 **1 行，且那是注释正文**（"# \`if:\`, or made continue-on-error. It sits directly BELOW the gofmt step"）。既有步骤一行未动：28 行全在 `gofmt (gofumpt)` 那一步的 `fi` 之后、`go vet (module)` 之前。

命令＝`run: sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only`。空心保护**复用 r4 那条**（我把它从主流程移进 `ruler_a()`，同一条 `[ "$TRACKED_GO" -gt 0 ] || … exit 2`，未新写第二份逻辑）：递给尺 0 枚 `.go` ⇒ **rc=2**，CI 读成红，永不读成绿（实测 §4.3 第 4 发）。

### 3.2 措辞按"今天不增加射程"写，不写成补了大洞

步骤注释里那两段的口径：现有 `gofumpt -l . tools/d22scan tools/mockllm` 在**检出树**上与甲形等价（检出树里只有已跟踪文件），所以这一步**今天不增加射程**；它买的是"分母＝已跟踪集合"这句话由机器执行，不再靠人的记忆；并且它调的是同一把仪器，两形读数不可能各自漂移。

### 3.3 那一步在本地跑出的颜色

```
$ sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only
== (A) tracked tree: git ls-files -z '*.go' | xargs -0 D:\work\base\gopath/bin/gofumpt.exe -l   (tracked .go files handed to it: 531)
attrib.sh: (A) tracked-tree  lines=0 files=0  rule: MUST be empty
attrib.sh: GREEN (tracked-only) - (A) empty over 531 tracked .go files
$ 上条命令的 rc -> 0
```
红的那一发不在这棵共享树上做（§4 就是它，仓外副本）。

---

## 4. 第③格：`TRACKED-DIRTY` 那一发，仓外隔离副本，四发读数

`sh .scratch/wisp/probes/161/r5/tracked-dirty-proof.sh` → `logs/d2-tracked-dirty-proof-full.txt`，脚本自身 **rc=0**（它自己判"三对一"是否齐全，缺一即 exit 3/1）。副本落在 `/tmp/wisp-161-r5-trackeddirty.Y82apT`（`mktemp -d`，在仓外；脚本对"副本竟然在仓内"那一支直接 exit 2）。清单＝`git ls-files '*.go' | sed -n '1,40p'`（40 枚），`cp --parents` 过去，`git init` ＋ `git add` 让它有自己的索引，`go.mod` 与 `.gitattributes` 一起带过去（不带 `.gitattributes` 时那 40 枚 `git add` 会刷出满屏 CRLF warning，见 `logs/d1-…txt` 的 40 行；带过去之后**只剩 1 行**，且那一行是 `go.mod`，与本判据无关）。

| # | 那一发 | 读数 |
|---|---|---|
| 1 | 原样副本 | `tracked .go files handed to it: 40`／`lines=0 files=0`／`GREEN (tracked-only)`／`TRACKED_PRISTINE_RC=0` |
| 2 | 往副本里的 `.scratch/wisp/probes/147/zz147probe_windows_test.go` **追加**两空格缩进的合成品（能解析、gofumpt 不收） | `TRACKED-DIRTY	.scratch/wisp/probes/147/zz147probe_windows_test.go	(last commit touching it: -)	<== CI IS RED ON THIS`／`lines=1 files=1`／`RED (tracked-only) - tracked-dirty=1 files_dirty=1`／`TRACKED_DIRTY_RC=1` |
| 3 | 从仓内同名件 `cp --parents` 复原，`cmp` 过（`byte-identical to …/probes/147/zz147probe_windows_test.go after restore`） | `lines=0 files=0`／`GREEN (tracked-only)`／`TRACKED_RESTORED_RC=0` |
| 4 | 第二条空副本（只有 `go.mod`，`git ls-files '*.go'` 数＝**0**）跑同一条命令 | `attrib.sh: (A) ruler saw 0 tracked .go files - refusing to report 'empty' from a ruler that was handed nothing`／`hollow-guard rc=2` |

`TRACKED-DIRTY` 这条腿**第一次被看见红**，红在没人共享的树上；`(last commit touching it: -)` 是副本里没有历史所致（`git log` 无 commit ⇒ 仪器那条 `|| echo '-'` 兜住），不是漏读。复原之后我在仓内复核过：`git status --porcelain -- .scratch/wisp/probes/147/zz147probe_windows_test.go` ⇒ **空**（受害者那棵共享树一字节未动）。

两枚禁止的形状都没用：`git worktree` 与仓内 `checkout` 一次未跑（`git status` 无新增 worktree，副本在 `/tmp`）；`git archive | tar -x` 一次未跑，副本的字节来自 `cp --parents`，且第 3 发用 `cmp` 钉住了"复原＝与仓内同字节"而不是"复原＝归档解出来的东西"。

---

## 4b. 台件全部入库之后再跑一次那三把尺（AC#7② 的"尺最后跑"，本程自己的顺序）

四处提交（cell 1 / cell 2 / cell 3 / 本件）落盘之后，`date` = `Sun Sep 27 08:5x CST 2026`，同一窗口重跑：

| 尺 | 读数 | 日志 |
|---|---|---|
| 甲形（CI 那一步的命令） | `lines=0 files=0`、递给尺 **531** 枚、rc=**1**（红只来自乙形那一支）→ 甲形自身 `GREEN (tracked-only) - (A) empty over 531 tracked .go files`、rc=0 | `logs/e1-attrib-post-commits.txt`（全文两形） |
| 乙形（工作树） | `lines=10 files=8 tracked=0 attributed_tickets=[161] unattributable=1`、rc=1 ⇒ **与提交前那一发（§2.2）完全同形**，本程落盘没多出一枚归不出的行 | 同上 |
| `sh scripts/d22scan.sh` | rc=**0**、`clean - no D22 ban violations` | `logs/e2-d22scan-post-commits.txt` |
| `go vet ./`（`tools/d22scan` 模块内） | rc=**0**、零行输出 | `logs/e3-vet-d22scan-module.txt` |
| `go vet ./...`（根模块） | rc=**0**、**0 行**输出（`wc -l` = 0） | `logs/e4-vet-root.txt` |

⚠ 上表第一行那句"rc=1"是**整把尺**的退码（两形一起跑那一发），不是甲形那一支的；甲形单独跑是 rc=0。派单 §6 说过"截断输出不是全表"，所以两发都留了全文日志。

---


## 5. 门禁（四数＋名册差集，同一窗口）

| 门禁 | 读数 | 命令 |
|---|---|---|
| `tools/d22scan` 全量 | `PASS=34 FAIL=0 SKIP=0`、`=== RUN=76`、`[no tests to run]`=0、rc=0 | `bash tools/d22scan/runtests.sh -C tools/d22scan ./...`（`logs/g2-runtests-post.log`；改前 `logs/g1-runtests-pre.log` 同四数） |
| 名册两向差集 | 各 **34** 行，`comm -3` 输出**空**（rc=0） | `grep -E '^(=== RUN\|--- (PASS\|FAIL\|SKIP)): ' <log> \| sed 's/ (.*//' \| sort` → `logs/g1.roster.txt`／`logs/g2.roster.txt` |
| 全仓扫描门 | rc=**0**，`d22scan: clean - no D22 ban violations`；`.scratch/**` 不在它的任何射程里（名册现读：bans #1-5 走 `internal/`＋`cmd/`、#6 走 `frontend/`、#7 走 `internal/tools/`、#8 走 `design/`＋`frontend/`＋`internal/`＋`cmd/`），所以 999 那枚在场与否都与这一发无关 | `sh scripts/d22scan.sh`（`logs/g3-d22scan-sh.txt`） |
| 甲形（CI 那一步） | `lines=0 files=0` over **531**，rc=0 | `sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only` |
| 乙形（工作树） | `lines=10 files=8 unattributable=1`，rc=1（§2.2） | `sh .scratch/wisp/probes/161/r5/attrib.sh` |
| 尺的版本 | `v0.12.0 (go1.27.1)`（两把：直接跑＋仪器自己印） | `"$(go env GOPATH)/bin/gofumpt.exe" --version`；`grep 'attrib.sh: gofumpt' logs/c3-…txt` |
| 分母两把尺 | `531` 与 `531` | `git ls-files '*.go' \| wc -l`；`git ls-files -z '*.go' \| tr -dc '\0' \| wc -c` |
| §0.3 那段 AC#7 的逐字性 | **1911 vs 1911 字节、`diff` 无输出**（把我加的行号列剥掉之后与票面 35–43 行 byte-identical，含那 8 行原有的两空格续行缩进） | `sed -n "/^35  - \\[ \\] \\*\\*AC#7/,/^43  /p" docs/evidence/s1/161-attrib-selftest-r5.md \| sed -E 's/^[0-9]{2}  //' > /tmp/q2.txt` 对 `sed -n '35,43p' $ISS > /tmp/q1.txt`，再 `diff` |

⚠ 派单里"`sh scripts/d22scan.sh` 预期 rc=0"——**符**。⚠ 派单说基线名册要与 34 对得上——**符**，且本程没加任何 Go 用例（不写 Go 就没有可加的），所以 34→34 是本程的预期而不是"没动"的借口：本程确实一行 Go 都没写。

---

## 6. 本程没做的事（免得下一位以为已经做了）

- 没退役 `probes/999/bad-sample.go`（§2.2），所以乙形在这棵树上仍 rc=1；
- 没接 `probes/161/r3/gate-order.sh`（编排者 00:0x 裁三：不编辑不删，`attrib.sh` 是唯一在用那把），本程只在 §0.4 的"从未红过"那一格里引用过它留下的日志名；
- 没测 CI runner（§1 第一条）；
- 没动 AC#2/AC#3/AC#4/AC#6 任何一格，没勾任何框。

## 7. 契约轴与写面（AC#5 的自查）

`git diff --name-only <本程各提交>` 只会出现 §5 派单写面四类：`probes/161/r5/**`、`docs/evidence/s1/161-attrib-selftest-r5.md`、`.github/workflows/ci.yml`、票面 Progress log 那一格。`tools/d22scan/**`、`allowlist.txt`、`docs/PLAN.md`、`docs/specs/**`、`internal/**`、`cmd/**`、`frontend/**`、`design/**`、`docs/reports/**`、`probes/161/r3/**`、`probes/161/r4/**`、别人的票面 ⇒ **零字节**（提交后以每枚提交的 `git show --stat` 为凭，§8）。

## 8. 被拒／没成功的调用（下面 1 是取数之前的探路，2–4 是取数之后；都不影响上面的读数）

1. `gofumpt --version`（裸名）⇒ `command not found`、rc=127 —— 取数**之前**探到的 PATH 事实，随后所有发都走 GOPATH 全路径（已写进抬头）。
2. 一次 `Edit`（`tracked-dirty-proof.sh`）⇒ `0 occurrences found`，因为我把注释里两个词记错了；改成小段编辑后成功。
3. `--self-test` 第一版 rc=1（§2.3 末：我的 `ticket-158` 期望错）⇒ 这是**读数不是故障**，日志留在 `logs/c1-selftest.txt`，并按它改了那一发而不是改尺。
4. `git status --porcelain -- …`（step 0）非空：两枚 161-r2 的未跟踪遗留（`r2/__pycache__/`、`r2/ctl/`），**未动**（§0.2）。

## 9. 有没有跑过删除命令

**没有。** 全程 `rm`／`rmdir`／`del`／`git rm`／`git clean`／`git checkout .`／`--amend`／`reset`／`rebase`／`stash` 零次；`git worktree` 零次。要区分开写的两件事：
- `tracked-dirty-proof.sh` **不清自己的临时目录**（派单 §7 允许它清，我选了不清，好让 §4 的四发可复跑），所以盘上多了 `/tmp/wisp-161-r5-trackeddirty.Y82apT` 与 `/tmp/wisp-161-r5-hollowguard.ncBBhW` 两枚仓外副本，都在 `mktemp` 的地界里，仓内无物；
- 第 3 发"复原"是 `cp --parents`（覆盖写），不是删除；`cmp` 的 byte-identical 那一行是它的凭据。

## 10. 伪授权两栏

**A 栏＝我说了、且每条旁边有命令的**：甲形 0 行／531 枚；乙形 10 行、1 枚归不出、rc=1；自测 8 发 5 响 3 不响 rc=0；`TRACKED-DIRTY` 在仓外副本 rc=1、复原 rc=0、空集合 rc=2；`ci.yml` 加 28 行删 0 行、新步骤只有 `name`＋`run` 两键；四数 34/0/0/76 与名册差集空；`sh scripts/d22scan.sh` rc=0。
**B 栏＝我没说、也别被读成我说了的**：没说乙形今天转绿（它没有，§2.2）；没说 CI 那一步"补上了一个大洞"（它今天不增加射程，§3.2）；没说 `attrib.sh` 能区分"故意坏"与"不小心坏"（§1）；没说 `gate-order.sh` 被修好了（一字未动，也没跑它）；没说票 161 可结案（AC#2/#3/#4/#6 与那枚非实现者表都还欠着）；没说 `probes/999` 那枚伤已由我处理（它等 owner）。

## 11. 凭据

零枚凭据、零枚明文密钥被读取或抄录。本程读过的文件全是尺、脚本与票面。

## 12. `next=`

1. **owner 退役 `.scratch/wisp/probes/999/bad-sample.go`**（一次删除，或改成 `probes/161/r5/negative-control/` 那份的形状）；退役后 `sh probes/161/r5/attrib.sh` 的乙形 ⇒ 期望 **`lines=9 files=7 unattributable=0` rc=0**，那是本程唯一没能自己交出的读数（§2.2 已给出现状两发）。
2. 把 `--self-test` 也接进 CI 需要一次授权：它今天只在本地跑（派单 §2 只准加甲形那一步）。若要加，判据应含"**自测自己红过的那一发**"——本程的证据是 `logs/c1-selftest.txt`，但那是我的期望错而非尺错，真正稳的做法是在 `tools/d22scan` 的 Go 测试里对 `attrib.sh --self-test` 的 rc 加一枚断言（**注意**：那是 `tools/d22scan/**`，本票 AC#2 允许的是"测试面与自检入口"，属可讨论范围，不属我这一程的写面）。
3. `classify_line` 不解析 `..` 那一格（§1 第 5 条）：如果哪天要它解析，记得 r4 那句"`filepath.Clean` 只能在 `risk.PathResolver` 之外…"是**代码禁令**、尺是 `d22scan`，与本脚本无关，但把 `..` 折进判据前请先想清楚"归因会不会因此变宽"。
4. 两枚 `/tmp` 副本（`trackeddirty.Y82apT`／`hollowguard.ncBBhW`）留在盘上等 owner；`tracked-dirty-proof.sh` 每复跑一次就新增两枚。
5. 派单 §0.2 那两枚 161-r2 未跟踪遗留（`r2/__pycache__/`、`r2/ctl/`）等编排者归口。
