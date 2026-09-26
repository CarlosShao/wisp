# 161-r4 — AC#7② 归因仪器（两形读数）＋ AC#4 真摘一次

票面＝`.scratch/wisp/issues/161-gates-need-their-own-gates-self-tests-in-ci-plus-is-the-approval-door-ever-rung.md`
派单＝`.scratch/wisp/dispatches/2026-09-26-235x-impl-161-r4-gate-order-attribution.md`
锚点（step 0 现量，非抄号）：

```
$ git rev-parse HEAD
3e56e68051e57edd59a0522936a23057ffb7f92e        # dev
$ date
Sat Sep 26 23:38:57 CST 2026
$ git status --porcelain -- .scratch/wisp/issues/161-*.md .github/workflows/ci.yml tools/d22scan/
（空）
```

档位＝**实现者自证**（本程自己跑的每一条读数旁边都放了产生它的命令）。
本程**没有改动 `tools/d22scan/**` 的任何一字节**（见 §4 的 sha256 前后同值），也**没有动 `.github/workflows/ci.yml`**。

---

## 1. AC#7② — 两把各自能出读数的尺 ＋ 一枚归因仪器

### 1.1 判据为什么必须拆成两形（承重的那句，不是修辞）

原句「`gofumpt -l . tools/d22scan tools/mockllm` 必须为空」在**工作树**上是**结构上永远产不出**的读数：

```
$ git ls-files '.scratch/**/*.go' | wc -l
40        # gofumpt 走 .scratch，这 40 枚是已入库取证台件
```

取证机械的定义就是「往里放故意坏的样本」，所以工作树那一发非空**是常态而不是伤**。
CI 那一发能过，只是因为 CI 的检出树里只有已跟踪字节。⇒ 拆成：

| 形 | 命令 | 判据 |
|---|---|---|
| （甲）已跟踪树 | `git ls-files -z '*.go' \| xargs -0 gofumpt -l` | **必须为空**（＝CI 实际执行的那形） |
| （乙）工作树 | `gofumpt -l . tools/d22scan tools/mockllm` | 预期非空；**每一行必须归到一枚真实票号**，归不出＝真伤 |

### 1.2 仪器＝`.scratch/wisp/probes/161/r4/attrib.sh`

硬退码（这是它成为仪器而不是打印机的原因）：

```
rc=1  （甲）非空                                   -> 已入库字节是脏的，CI 红
rc=1  （乙）出现归不出的行                          -> 真伤，追它
rc=1  （乙）出现 TRACKED 行而（甲）为空              -> 两形互相打脸
rc=2  仪器自己跑不起来（没有 gofumpt / 不在 git 仓 / （甲）被递给 0 枚文件）
rc=0  （甲）空 且（乙）每一行都归得出
```

归因规则（**两条都要成立**，缺第二条这格就成恒真检）：

```
path == .scratch/wisp/probes/<NN>/**   AND   .scratch/wisp/issues/<NN>-*.md 存在
```

### 1.3 读数（本程现量，gofumpt `v0.12.0 (go1.27.1)`，版本由 `gofumpt --version` 现跑）

**（乙）8 行 / 6 枚不同文件，全部归得出＝票 161；（甲）空：**

```
$ sh .scratch/wisp/probes/161/r4/attrib.sh          # 逐字取自 logs/attrib-run3-final-clean.txt
== (A) tracked tree: git ls-files -z '*.go' | xargs -0 D:\work\base\gopath/bin/gofumpt.exe -l   (tracked .go files handed to it: 529)
== (B) working tree: D:\work\base\gopath/bin/gofumpt.exe -l . tools/d22scan tools/mockllm
   ticket-161	.scratch/wisp/probes/161/r1/runs/b8-probe-bands/internal/a/p_b8.go (ignored, parse_err=0)
   ticket-161	.scratch/wisp/probes/161/r1/runs/unp-bad/internal/a/p_unp.go (ignored, parse_err=1)
   ticket-161	.scratch/wisp/probes/161/r1/runs/unp-bad/internal/a/p_unp.go (ignored, parse_err=1)
   ticket-161	.scratch/wisp/probes/161/r1/runs/unp-hides-bans/internal/a/p_unp.go (ignored, parse_err=1)
   ticket-161	.scratch/wisp/probes/161/r1/runs/unp-hides-bans/internal/a/p_unp.go (ignored, parse_err=1)
   ticket-161	.scratch/wisp/probes/161/r3/pre-format/guard.no1.go (ignored, parse_err=0)
   ticket-161	.scratch/wisp/probes/161/r3/pre-format/guard.no2.go (ignored, parse_err=0)
   ticket-161	.scratch/wisp/probes/161/r3/pre-format/guard.no3.go (ignored, parse_err=0)

attrib.sh: gofumpt v0.12.0 (go1.27.1)
attrib.sh: (A) tracked-tree  lines=0 files=0  rule: MUST be empty
attrib.sh: (B) working-tree  lines=8 files=6  tracked=0 attributed_tickets=[161] unattributable=0
attrib.sh: GREEN - (A) empty and all 8 working-tree lines attributed to a named ticket   # rc=0
```

全文＝`logs/attrib-run3-final-clean.txt`。三条独立核对：

```
$ git ls-files '.scratch/wisp/probes/161/r1/runs/' | wc -l
0                  # 那 5 行全在 .gitignore 之下（r1/.gitignore:11 runs/，r3/.gitignore:1 pre-format/*.go）
$ git check-ignore -v .scratch/wisp/probes/161/r1/runs/unp-bad/internal/a/p_unp.go
.scratch/wisp/probes/161/r1/.gitignore:11:runs/   ...p_unp.go
$ grep -c 'parse_err=1' logs/attrib-run3-final-clean.txt ; grep -c 'parse_err=0' ...
4  /  4            # 8 行里 4 行是"解析不了"（2 枚 unp 文件各 2 行），4 行是"要格式化"
```

⇒ 派单里「今天 8 行、全部在 probes 之下、全归得出」＝**成立**；r3 脚本注释里「5 of today's 8 lines are parse errors」＝**差一枚**（实测 4 行诊断 / 4 行纯文件名）。

### 1.4 负控：这把仪器不是恒绿（**这一发是判据要求的那一发**）

在 `.scratch/wisp/probes/999/bad-sample.go` 放一枚故意不格式化的 `.go`（**只建不删，未提交，也不加 .gitignore**——它偏在 `probes/<NN>` 约定里，但票 999 不存在）：

```
$ sh .scratch/wisp/probes/161/r4/attrib.sh ; echo RC=$?
   UNATTRIBUTABLE	.scratch/wisp/probes/999/bad-sample.go (path claims ticket 999 but .scratch/wisp/issues/999-*.md does not exist, class=untracked, parse_err=0) <== REAL INJURY
attrib.sh: (A) tracked-tree  lines=0 files=0  rule: MUST be empty
attrib.sh: (B) working-tree  lines=9 files=7  tracked=0 attributed_tickets=[161] unattributable=1
attrib.sh: RED - tracked-dirty=0 unattributable=1 forms-disagree=0
RC=1
```

全文＝`logs/attrib-run4-negative-control-red.txt`。
⚠ **由此产生的现场**：这枚 999 样本按"只建不删"会长留，所以**今天之后 `attrib.sh` 默认就是 rc=1**。
这不是坏掉的仪器，是它按要求认下的一枚已知伤：**CI 侧的硬门是（甲）那一形**（999 未跟踪 ⇒ 甲不受影响，实测 `lines=0`）。
要把（乙）当门用，唯一诚实的动作是**给 999 开一张真票**或在下一程把这枚样本换掉，**不是**给它加豁免。

### 1.5 与 161-r3 留下的 `gate-order.sh` 的关系：接不上，自己写了一枚

派单说「r3 留了一枚**未提交**的 `gate-order.sh`」——**盘上不符**：它在 HEAD 里。

```
$ git status --porcelain -- .scratch/wisp/probes/161/          -> 只有 ?? r2/__pycache__/ 与 ?? r2/ctl/
$ git log --diff-filter=A --format='%h %s' -- .scratch/wisp/probes/161/r3/gate-order.sh
3e56e68 票 161 r3 撞顶代提（AC#2 自检入口落地＋AC#7 判据改两形读）＋ledger(A318)
```

我没改它一个字（派单也这么要求），但把它**照 committed 路径跑了一遍**，结果就是我不接它的理由：

```
$ sh .scratch/wisp/probes/161/r3/gate-order.sh      # 全文 logs/gate-order-r3-probe.txt
gate-order.sh: repo root = /d/work/workspace/projects plans/Wisp/.scratch   # <-- 根算短了一级（它写 ../../../..）
   tracked=0 untracked=11                            # 11 行里含两条 "GetFileAttributesEx tools/... not found" 错误串
== ruler 2/3: sh scripts/d22scan.sh
   0 finding lines; script rc=127                    # <-- 那一腿根本没跑到，却被记成 0 finding
== ruler 3/3: go vet ...
   root: found packages tools (...) and main (...) in ...Wisp\.scratch\wisp\probes\151
gate-order.sh: TRACKED offenders across the three rulers = 0
gate-order.sh: clean on the committed set            # <-- rc=0
```

⇒ 它在**自己的路径**上把根算到 `.scratch`：三把尺没有一把在正确的树上跑，而它照样打「clean」并退 0。
这正是 AC#7② 要防的形状（"门点红了没人能在推送前知道是谁的台件"），也解释了为什么它归不出 999 那一行。**它归因只做 tracked/untracked 两类，没有票号那一档，也没有"归不出⇒红"这条硬退码。**
本程没去修它（那是别人的已提交件），只把它的读数留在上面这条命令里。

---

## 2. AC#4 — 真摘一次：把一枚「必须响」样本从名册里拿走

### 2.1 怎么摘的（**没删任何字节，没改任何已跟踪文件**）

样本表是 `tools/d22scan/selftestsamples.go` 里的 Go 字面量（不是磁盘上的 fixture 文件），所以"挪走"用
`go build/run -overlay` 把**一张挪开的副本**递给编译器，原件一字节未动：

```
$ sed '190,195d' tools/d22scan/selftestsamples.go > .scratch/wisp/probes/161/r4/removed/selftestsamples.no-mirror-ring.go
$ cat .scratch/wisp/probes/161/r4/removed/overlay.json      # 只把 selftestsamples.go 指向上面那枚副本
$ sha256sum tools/d22scan/selftestsamples.go                # 摘之前
7ea3f1d8f664ac61a0e3287b5fc2fc4c6a0bc94871a0fcc9005f57790d85daf5
$ sha256sum tools/d22scan/selftestsamples.go                # 全部读数跑完之后
7ea3f1d8f664ac61a0e3287b5fc2fc4c6a0bc94871a0fcc9005f57790d85daf5      # 同值 = 原件从未被改
```

摘掉的那一枚＝**ban #5 `mirror-hash` 的 `wantRing` 样本**（`internal/probe/mirror.go`，`selftestsamples.go:190-195`）。
选它的原因（现读名册，`grep -n 'tag:' selftestsamples.go`）：它是 `mirror-hash` 这一 tag **唯一**的 ring 样本，
而这台名册里有 **7 枚 tag 各带 2–4 枚 ring 样本**——所以我还摘了第二枚做对照（§2.4）。

### 2.2 读数

```
$ cd tools/d22scan && go run . -self-test                     # 控制组：名册完整
rc=0
d22scan -self-test: clean - all 34 direction checks passed (19 expect-ring, 15 expect-silent)

$ go run -overlay='.../removed/overlay.json' . -self-test      # 摘掉之后
rc=1
d22scan -self-test: FATAL the table does not cover the tool (1 hole(s)); an unrun self-test is not a green self-test
d22scan -self-test:   HOLE tag "mirror-hash" has only an expect-silent sample - a ban whose violating sample was deleted cannot be distinguished from a ban that was never implemented

$ go build -overlay='.../removed/overlay.json' -o .../d22scan.mirror-removed.exe . && ./d22scan.mirror-removed.exe -self-test
REMOVED_BINARY_rc=2      # 程序自己的退码；`go run` 把它塌成 1（CI 那一步用的就是 go run，见 ci.yml:105）
$ ./.../d22scan.intact.exe -self-test
INTACT_BINARY_rc=0
```

两条日志全文：`logs/selftest-control-intact.txt`、`logs/selftest-mirror-ring-removed.txt`。
⇒ 派单里「名册空心 ⇒ exit 2、拒绝出判语」这条**是真的会走**：它一行 case 都没跑（输出里没有任何 per-case 行），
只出 FATAL＋HOLE，退 2。**这不是引用代码，是跑出来的。**

### 2.3 测试面与 CI 侧的可见读数

```
$ bash tools/d22scan/runtests.sh -C tools/d22scan -overlay='.../removed/overlay.json' ./...
rc=1   packages=[-overlay=... ./...] top-level: PASS=32 FAIL=2 SKIP=0, === RUN=76, '[no tests to run]'=0
--- FAIL: TestSelfTestEntryPassesEveryCase
--- FAIL: TestSelfTestRosterAuditRejectsAHollowTable
        selftest_test.go:75: the shipped table must audit clean, got 1 hole(s): ...
```

名册差集（把时间戳尾巴剥掉再 `comm`）：

```
$ comm -3 r4-pre-runtests.log.roster.txt logs/runtests-mirror-ring-removed.log.roster.txt
	--- FAIL: TestSelfTestEntryPassesEveryCase
	--- FAIL: TestSelfTestRosterAuditRejectsAHollowTable
--- PASS: TestSelfTestEntryPassesEveryCase
--- PASS: TestSelfTestRosterAuditRejectsAHollowTable
```

### 2.4 承重两问（各一句）

- **① 摘掉之后有没有哪一发从此打不红＝那枚样本承重吗？** 承重，且是**tag 级**承重不是 case 级：
  摘掉 `mirror-hash` 唯一的 ring 样本后，**没有一发变哑**——`-self-test`（退 2／`go run` 退 1）、`go test`（2 枚红）、
  CI 的两步（`D22 scanner self-test` 与 `D22 scanner positive control`）**同时**红。
  对照组证明它不是"摘任何一枚都红"：摘 `bare-goroutine` 三枚 ring 样本里的一枚（`removed/selftestsamples.no-baregoroutine-ring1.go`＋`overlay2.json`）⇒
  **rc=0、`clean - all 33 direction checks passed (18 expect-ring, 15 expect-silent)`**（`logs/selftest-baregoroutine-ring1-removed.txt`）——
  那枚冗余样本单独摘掉**不承重**，仪器也不因此恒红。
- **② 摘掉之后有没有任何外部可见读数变过？** 三处：退码 `0→1`（`go run`，CI 看到的那形）与 `0→2`（二进制自身）；
  判语从 `clean - all 34 direction checks passed` 变成 `FATAL the table does not cover the tool (1 hole(s))`；
  `go test` 的 `PASS=34 FAIL=0` 变成 `PASS=32 FAIL=2`（CI 那两步同色变红）。

### 2.5 这一格**没有**证明的事

它证明的是"摘掉名册里的一枚唯一 ring 样本 ⇒ 门拒绝出绿色判语"。它**没有**证明：
把 `main.go` 里某条禁令的**正则**改瞎之后自检会红（那是另一发，属于禁令射程面，本票 AC#5 冻结，本程一字节未动）。

---

## 3. 门禁四数（同一安静窗口，改前／改后）

```
$ bash tools/d22scan/runtests.sh -C tools/d22scan ./...       # 改前  logs: r4-pre-runtests.log
runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0   rc=0
$ bash tools/d22scan/runtests.sh -C tools/d22scan ./...       # 改后  logs: r4/logs/runtests-post-restore.log
runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0   rc=0
$ comm -3 <(两向，剥时间戳) -> 空（34 vs 34 枚同名）
$ sh scripts/d22scan.sh                                        # 预期 rc=0
d22scan: clean - no D22 ban violations; ...                    rc=0   （999 那枚坏样本在场时仍是 0）
```

⚠ **对派单/票面基线的一处更正**：票面 21:3x 那格写「`tools/d22scan` 那一包本来就有两枚红
（`TestScannerSelfScanOfRealRepoIsGreen`、`TestRealRepoLedgerIsHonest`，病因＝票 169）」——
**在 `3e56e68` 上这两枚今天都是 PASS**（`FAIL=0`）。本程没有为凑这个说法去动任何断言；
只是提醒下一位：别按"改前有两枚红"去核对名册差集，会以为别人修了什么。

`gofumpt -l` 两形：见 §1.3（甲空 rc=0；乙 8 行全归得出）与 §1.4（负控 rc=1）。
`--version`：`v0.12.0 (go1.27.1)`（`attrib.sh` 每次自己现跑并打这一行）。

**提交之后再跑最后一遍**（AC#7② 的那条顺序规矩：三把尺在台件全部入库之后跑最后一次）——
本程新提交的两枚 `removed/*.go` 是已跟踪 .go，所以它们**进了（甲）的分母**：

```
$ git rev-parse --short HEAD ; sh .scratch/wisp/probes/161/r4/attrib.sh
77e9ed2
attrib.sh: (A) tracked-tree  lines=0 files=0  rule: MUST be empty        # 仍空 = CI 不会因我的台件变红
attrib.sh: (B) working-tree  lines=9 files=7  tracked=0 attributed_tickets=[161] unattributable=1
（rc=1，唯一归不出的那行还是 §1.4 的 999 负控）   # 全文 logs/attrib-run6-post-commit.txt
$ git show --name-only --format= 77e9ed2 | grep -v -E '^\.scratch/wisp/probes/161/r4|161-gate-order-r4\.md|161-gates-need'
（空）   # 30 枚文件里没有别人的路径
```

---

## 4. 本程写面 / 没写面

写：`.scratch/wisp/probes/161/r4/**`（含 `logs/`、`removed/`）、`.scratch/wisp/probes/999/bad-sample.go`（负控，**不提交**）、
`docs/evidence/s1/161-gate-order-r4.md`（本件）、票 161 Progress log 一格。
**没写**：`tools/d22scan/**`（sha256 前后同值，见 §2.1）、`allowlist.txt`、`.github/workflows/ci.yml`、
`docs/PLAN.md`、`docs/specs/**`、`internal/**`、`cmd/**`、`thresholds.go`、golden、`scripts/slo-check.ps1`、`frontend/**`、`design/**`、别人的票面。
**没跑任何删除命令**：`rm`/`rmdir`/`del` 一发未动（唯一被执行过的 `rm -f` 是 `runtests.sh` 内部删它自己 mktemp 的那一行，
由派单 §4 指定的那条命令带出来，不是我下的手，也没落在别人的字节上）。

## 5. next=（给编排者裁，本程不自行执行）

1. **（乙）那一形要不要进门**：今天它 rc=1 是设计如此（999 负控）。若要当门用，先定 999 那一行归谁——
   要么真开一张票 999（`.scratch/wisp/issues/999-*.md`），要么把它登记成常驻负控并明写"（乙）不进 CI"。
   ⚠ 别用"加豁免"解决，那是本票 AC#7 明令排除的形状。
2. **CI 若加一步，只加（甲）**：`git ls-files -z '*.go' | xargs -0 gofumpt -l` 空判非空——这一形与
   `ci.yml:139` 现在那发（`gofumpt -l . tools/d22scan tools/mockllm`）在检出树里等价，但**在未跟踪台件在场的本机不等价**；
   要不要换、要不要并存，由你裁（本程按派单未动 `ci.yml` 一字）。
3. **`gate-order.sh` 的根算短一级**（§1.5）：它是 r3 已提交件，本程无权修。建议起一张小票（或并进 161 的下一程）
   把 `../../../..` 改成"向上找 `.git`＋校验 `go.mod`"，并给 ruler 2 的 `rc=127` 加硬退码——
   **它现在会在三把尺都没跑对的情况下打 clean 并退 0**，这比没有仪器更坏。
4. **（甲）的 TRACKED-DIRTY 那一发从未实测红过**（§5' of my report / 本文件 §4）：要证它得把一枚已跟踪 `.go` 弄脏，
   共享工作树里不该由实现程顺手做。可复算做法＝另开一个 worktree/`git archive` 检出，在里面弄脏一枚再跑 `attrib.sh`。
