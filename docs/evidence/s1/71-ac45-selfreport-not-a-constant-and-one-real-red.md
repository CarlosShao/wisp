# 71 — 自报的数「不是常量」（AC#4）＋ CI 那一步的一次真红（AC#5）

- **票**：`.scratch/wisp/issues/71-gates-must-self-report.md`（AC#4、AC#5）
- **本腿**：`71-ac45`（写面：`tools/d22scan/**`、`scripts/d22scan.sh`、本文件、本票 Progress log）
- **开工 HEAD**：`8d56a948d`（2026-11-14 18:0x，`git rev-parse --short HEAD`）
- **交件 commit**：`4ce5c86ff`（3 文件 +537 / -2：`tools/d22scan/selfreport_delta_test.go` 新、`tools/d22scan/export_red_proof.sh` 新、`scripts/d22scan.sh` 改）
- **本文件是**：实现者自账（`.scratch/wisp/issues/README.md` 规则 6 的 1:1 表）。**不是**对抗验收 —— AC#7 要求「验收者 ≠ 实现者」，归编排者派。

---

## 0. 先量现状：AC#4 的「自报」本体**早在盘上**，本腿没有从零重写

派单硬要求 3：若 AC#4 描述的形状与盘上现状不符，先报现状数字再改。现跑（`cd tools/d22scan && go run . -root <abs>`）：

| 作用域（票面 AC#4 要求的那一行 `scope … examined N files`） | 工作树 `8d56a948d` | 仓外导出树 `4ce5c86ff` |
|---|---|---|
| `bans #1-5 internal/` | 405 | 404 |
| `bans #1-5 cmd/` | 78 | 78 |
| `ban #6 frontend/` | 123 | 122 |
| `ban #7 internal/tools/` | 51 | 51 |
| `ban #10 internal/risk/` | 21 | 21 |
| `ban #8 design/` | 39 | 14 |
| `ban #8 frontend/` | 123 | 122 |
| `ban #8 internal/` | 1014 | 1012 |
| `ban #8 cmd/` | 283 | 283 |
| **合计作用域行数 / rc** | **9 行 / rc=0** | **9 行 / rc=0** |

（两列不同不是矛盾：导出树没有 `.git`，扫描器因此**不适用** `.gitignore` 而全量扫，工作树则跳过 `frontend/dist/` 这类未跟踪产物 —— 扫描器每次都会在 stderr 大声说明这件事，是 A207 的「响亮方向」。）

**结论**：AC#4 的「自报」本体（8→9 条作用域台账、`emptyLiveScope` / `undeclaredKeys` / `driftedAbsentScope` 三条守卫）由票 71 前一位实现者 2026-09-21（commit `b551fef`，见票面 Progress log 10:2x）落地，票 350 又加了 `ban #10`。⇒ **⛔ 本腿没有改自报本体、没有重写台账**，只补它当时缺的两件可判物：

1. **这个数不是常量**的证明（AC#4 的字面判据「像门面」的反面）；
2. **CI 那一步的一次真红**（AC#5 的字面判据）。

---

## 1. 落地物

| 文件 | 作用 |
|---|---|
| `tools/d22scan/selfreport_delta_test.go`（新） | **9 个子测试**，每个作用域一份：先量 `n0`，往该作用域塞**一个**坏文件，要求**同一个作用域那一行** `n0 → n0+1`、当次 `rc 0 → 1`、红句**点名**该路径与该 ban；另有一条**覆盖守卫**：`declaredScopes()` 多出一个 live 作用域而用例没跟上 ⇒ 自己红（新作用域不能带着未钉的数落地）。断言同时打在**内部计数器**和**打印出来的那一行**上（只钉内部映射的话，打印一个过期数字仍会全绿 —— 那正是本票「话术」病灶）。 |
| `tools/d22scan/export_red_proof.sh`（新） | AC#5 的可复跑仪器：把**钉住的 revision** 导出到**仓外**目录，跑**CI 那条命令** `sh scripts/d22scan.sh` 三次（基线 / 塞一个违规 / 删掉），并要求 ①基线 rc=0、②塞入那次 rc=1 且输出点名该文件、③删掉后整树 sha256 **与塞前逐字相同**且 rc=0。导出完整性用 `git ls-tree -r <rev>` 逐条 `test -e` 校验（不是数文件个数：Windows 大小写折叠会让个数对不上而树是完整的）。 |
| `scripts/d22scan.sh`（改，写面允许的「CI 接线那一处」） | 堵掉一个真洞，见 §3.2。**只严不宽**：仍无 `|| true`、无 continue-on-error、无可跳过步；退出码取**第一个非零**，control 红的码与改前逐字一致。 |

---

## 2. AC#4 判据：塞一个坏文件 ⇒ 那一行 N 加一，且当次红

命令（CI 现有的阳性对照步会跑同一条，见 §3.1，所以**不需要新接线**）：

```
bash tools/d22scan/runtests.sh -C tools/d22scan ./... -run TestEveryDeclaredScopeCountsTheFileItIsGiven
```

实测（本 commit，`-v` 子测试日志逐行）：

| 作用域 | 塞入的坏文件 | n0 → n1 | rc | 打印行（节选） |
|---|---|---|---|---|
| `bans #1-5 internal/` | `internal/zz71delta/goroutine.go`（`go worker()`） | 3 → 4 | 0 → 1 | `scope bans #1-5 internal/      examined   4 production Go files` |
| `bans #1-5 cmd/` | `cmd/zz71delta/goroutine.go` | 1 → 2 | 0 → 1 | `scope bans #1-5 cmd/           examined   2 production Go files` |
| `ban #6 frontend/` | `frontend/src/zz71delta/decision.tsx`（`approval.decide`） | 1 → 2 | 0 → 1 | `scope ban #6 frontend/         examined   2 text files` |
| `ban #7 internal/tools/` | `internal/tools/zz71delta.go`（artifact-tool） | 1 → 2 | 0 → 1 | `scope ban #7 internal/tools/   examined   2 production Go files` |
| `ban #10 internal/risk/` | `internal/risk/zz71delta.go`（引 `aZz71Delta`，包内无声明） | 1 → 2 | 0 → 1 | `scope ban #10 internal/risk/   examined   2 production Go files, comment symbol citations` |
| `ban #8 design/` | `design/zz71delta/ball.md`（U+1F600） | 1 → 2 | 0 → 1 | `scope ban #8 design/           examined   2 text files` |
| `ban #8 frontend/` | `frontend/src/zz71delta/banner.txt`（U+1F600） | 1 → 2 | 0 → 1 | `scope ban #8 frontend/         examined   2 text files` |
| `ban #8 internal/` | `internal/zz71delta/emoji.go`（字符串里 U+1F600） | 3 → 4 | 0 → 1 | `scope ban #8 internal/         examined   4 Go files, comments and _test.go included` |
| `ban #8 cmd/` | `cmd/zz71delta/emoji.go` | 1 → 2 | 0 → 1 | `scope ban #8 cmd/              examined   2 Go files, comments and _test.go included` |

**整包读数**（本 commit，工作树）：`runtests.sh -C tools/d22scan ./...` ⇒ **rc=0**、`PASS=46 FAIL=0 SKIP=0`、`=== RUN=103`（改前 `PASS=45 / RUN=93`，新增正好这一条 + 9 个子测试）。`gofumpt -l .` 在模块目录为空。

**为什么这不是既有测试的重复**（不然就是白干）：

| 既有测试 | 它钉什么 | 它钉不到什么 |
|---|---|---|
| `TestLedgerCountsMatchAnIndependentWalk` | 真树上每个数与**独立走树**相等 | 一个恰好等于今天树值的**常量**照样过（等值 ≠ 会动） |
| `TestRealRepoLedgerIsHonest` | 真树无非空作用域、无非登记计数 | 同上：数不动它也绿 |
| `TestBuiltBinaryGoesRedEndToEnd` | 进程退出码在种违规时 1 / 空作用域时 2 | 不钉「哪个作用域的数动了」—— 它连数都不读 |
| **本条（新）** | **同一作用域：一个文件 ⇒ +1，且当次红、且点名** | — |

---

## 3. AC#5 判据：仓外导出树里造成真红、点名、还原后同 hash 再绿

### 3.1 CI 接线现状（先量，未新接线）

`.github/workflows/ci.yml`（HEAD `8d56a948d`）里 `lint` 作业的次序：

| 行 | 步 | 读法 |
|---|---|---|
| `:119` | `run: bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | 阳性对照：`-count=1`、`-v`、**SKIP 不算过** ⇒ 本腿 §2 的测试由它自动跑到 |
| `:144` | `run: go run . -self-test` | 扫描器自检（每个 ban 正反两向） |
| `:174` | `run: sh scripts/d22scan.sh` | **本文件 §3.3 证明会红的那一步**（含 control + 扫描） |

三步都带 `if: ${{ !cancelled() }}`（前一步红了也照样跑），没有 `continue-on-error`。⇒ **接线本来就在**，本腿补的是「它真的会红、红句能点名」的读数。

### 3.2 写这一发红时撞到的真洞（已堵，同一个 commit）

在导出树上塞违规后跑 `sh scripts/d22scan.sh`，**第一次读数**是：退出码 1，但**扫描步一行都没跑** —— 因为旧包装器是 `set -eu`，阳性对照（`runtests.sh`）一红就**结束脚本**，于是「树上有违规」这种最需要定位行的跑里，**定位行根本不存在**。这正是本票主题（「第一步红了，后面的步就被跳过」，D22 run-away mode 6）**下沉一层**：包装器把自己的门的判决藏在了自己的守卫后面。

处置（`scripts/d22scan.sh`）：两步都跑、都报；control 红时打一行 `WARNING - … read it as untrusted rather than as absent`；退出码取**第一个非零**。⇒ 红不会变绿、control 的退出码与改前逐字一致、多出来的只是「扫描的判决+定位行」。实测（导出树，塞入态）：`d22scan.sh: FAIL - control rc=1, scan rc=1`，选择加严而非"绕过"。

### 3.3 三跑读数（`sh tools/d22scan/export_red_proof.sh -r 4ce5c86ff`）

| 跑 | 命令 | 读数 |
|---|---|---|
| 导出 | `git archive <rev>` → 仓外目录（zip 路线；msys GNU tar 在 Windows 丢非 ASCII 名） | **20124 条 tracked 路径 / 0 缺失**；`内部`/`cmd`/`frontend`/`design` 四棵树齐 |
| ① 基线 | `sh scripts/d22scan.sh` | **rc=0**；9 条作用域行；`d22scan: clean`；`d22scan.sh: OK - control rc=0, scan rc=0`；整树 sha256 = `14d07cc44adbb715ddff54cbf0039b9c9bcc47dce5ca73fe188a288a32cb5975` |
| ② 塞入 | 加 `internal/zz71export/leak.go`（`go worker()`）后同一条命令 | **rc=1**；`bans #1-5 internal/` **404 → 405**；红句：`internal/zz71export/leak.go:5: [bare-goroutine] bare ``go worker(...)`` is banned (D22/D38b, R16: named calls count too): use observe.Registry.Spawn …`；control 侧同一点名（`PASS=44 FAIL=2`，两枚失败正是 `TestScannerSelfScanOfRealRepoIsGreen` / `TestRealRepoLedgerIsHonest`）；`d22scan.sh: FAIL - control rc=1, scan rc=1` |
| ③ 还原 | `rm internal/zz71export/leak.go` 后同一条命令 | **rc=0**；整树 sha256 = `14d07cc4…cb5975` **与①逐字相同**；数回 404；`d22scan.sh: OK - control rc=0, scan rc=0` |

⇒ AC#5 字面要求的三件（**红了 / 红句点名那个文件 / 还原后同 SHA256 再绿**）都有读数，且**是仓外导出树、跑的是 CI 那一条命令**。

---

## 4. ⛔ 不声称的东西（宁缺毋造）

1. **没有 GitHub Actions 上的一次红**：本腿被明令不得 push（派单与票面），那一半**物理不可证**，本文件不作任何「CI 已红」的暗示。§3.3 的措辞只到「CI 那一条命令，在仓外导出树上」。
2. **不是对抗验收**：AC#7 要求验收者 ≠ 实现者，本表只是实现者自账。
3. **导出树的计数 ≠ 工作树计数**（见 §0 注）：导出树无 `.git` ⇒ 不适用 gitignore ⇒ 全量扫。这不是缺陷，是扫描器每次都会大声说的事。
4. **一处观察到但没动的口径问题**：`bans #1-5 internal/` 这一行**也承载 ban #9（phantom-citation）的计数**（`selftest.go` 的 `selfCoverForTag` 把 `phantom-citation` 映射到 `internalGo`），即这一行的**名字比它实际承载的 ban 少说了 #9**。这属「话术与覆盖面脱钩」，改的是文案而非仪器；本腿没动（也不该由本腿单人改一个全队门禁的措辞），登记在此交编排者。
5. **本腿不勾任何 AC 框**：AC#4/#5 的框留给编排者按其判据判（票面规则：框由编排者/验收者勾）。

---

## 5. 复跑命令（三把尺）

```
# AC#4：9 个作用域各「塞一个坏文件 ⇒ 那一行 +1 且当次红」
bash tools/d22scan/runtests.sh -C tools/d22scan ./... -run TestEveryDeclaredScopeCountsTheFileItIsGiven

# AC#5：仓外导出树里的一次真红（默认导出 HEAD，-r 钉 revision，-k 留证据）
sh tools/d22scan/export_red_proof.sh -r 4ce5c86ff -C "$TMPDIR/wisp71-proof" -k

# 整包（CI 阳性对照跑的那条）
bash tools/d22scan/runtests.sh -C tools/d22scan ./...
```
