# 141 / Q-46c 对抗验收表 r1（裁决者 ≠ 实现者）

* **被验锚点**：`bb61dc5`（`2026-09-24 22:32 +0800`，`dev`）。裁决程 **只读快照、不改一字节生产码**。
* **纯净树取法（所有"整树读数"只在这里量）**：
  `mkdir -p /d/tmp/141-acc-base && git archive bb61dc5 | tar -x -C /d/tmp/141-acc-base`
* **变异树**（全部在 `/d/tmp` 下，**不在仓里、不删**）：`141-acc-parent`（`470e6c5`＝`ed2c077` 的父）、
  `141-acc-mut-strlen`、`141-acc-mut-lexical`、`141-acc-mut-noband`、`141-acc-mut-witharrow`、
  `141-acc-mut-backseed`、`141-acc-rev-schema`、`141-mut-noex`、`141-acc-probe`（③④样本）。
* **共享树纪律**：`internal/observe/sampler_settle_gate_136_test.go`（在飞程）与 `design/**`
  的 16 枚未提交删除＋未跟踪 `design/old/`、`design/doubao/` 全程**未读作事实、未 stage、未还原、未删**。
  本表每节只 `git add` 这一枚文件路径。

## 0. 被验的四枚 commit（`git show --name-only` 原文，逐枚复量"只带自己点名的路径"）

```
ed2c077 22:17  d22scan(ban #8): 面向用户的字符变严、注释面豁免、并同时补上原本漏掉的数学符号段 (票 141, Q-46 已批＝按推荐)
M	tools/d22scan/main.go
M	tools/d22scan/scan_test.go            （+469/−30：main.go 197、scan_test.go 302）
1218192 22:21  risk(票 141 具名解冻 ③): 审批卡 R7 文案那枚 `≥` 换成 ASCII，生产串与期望值同批 (Q-46 已批＝按推荐)
M	internal/risk/assessor_test.go
M	internal/risk/rules_scale.go
2970c79 22:21  memory(票 141 自决档): DDL raw 串里那两处"Go 看是字符串、SQL 看是注释"按严格读法清成 ASCII
M	internal/memory/schema.go
bb61dc5 22:32  evidence(141 Q-46c 交件表): 五条判据逐条读数＋判据(3)在盘上不可达那一维＋§10 对上 A197/Q-48
A	docs/evidence/s1/141-q46c-impl.md
```
⇒ **四枚都只带自己点名的路径，无越界、无删除（`--name-status` 全为 M/M/M/A）**。
`gofmt -l tools/d22scan internal/risk internal/memory` ＝ **空输出**；`go vet ./`（`tools/d22scan`）＝ **静默**。

⚠ 但 `git diff --stat 470e6c5..bb61dc5` 是 **9 枚文件／8 枚 commit**，不是 4 枚：
`054faae`、`f1b9c1c`、`e21d3bf`、`e848a93` 四枚 docs/evidence 夹在本批中间（`HANDOVER.md`、
`pending-and-issues.md`、`136-instr-fixes-r1.md`）。**不是越界**（不属于被验四枚），
只是"整批区间"与"本票四枚"不能混为一谈，读差异表时要点开。

## 1. 门禁基线四数（锚点快照现跑，`-v`，未动任何阈值）

```
$ sh tools/d22scan/runtests.sh -C tools/d22scan ./...      # 在 D:/tmp/141-acc-base
rc=1   PASS=22  FAIL=2  SKIP=0  === RUN=64  '^panic:'=0
--- FAIL: TestScannerSelfScanOfRealRepoIsGreen (1.53s)
--- FAIL: TestRealRepoLedgerIsHonest (0.91s)
```
* **名册差集**（防"一枚 panic 吞掉同包读数"）：`scan_test.go` 里 top-level `func Test*` **24 枚**，
  跑出 `--- PASS/FAIL` 的 top-level **24 枚**，`comm` 两向差集 **皆空** ⇒ 没有测试被吞。
* 整树真扫（同一快照，独立于测试）：`go run . -root /d/tmp/141-acc-base` ⇒ **rc=1／6 finding**，
  逐枚＝`frontend/fixtures/composer-states.html:2,5,8`、`frontend/src/components/ai-native/thinking.tsx:213`、
  `.../tool-chips.tsx:186`、`frontend/src/components/composer.tsx:170`；`design/`=0、`internal/`=0、`cmd/`=0。
  各 scope 计数：`design/ 16、frontend/ 40、internal/ 405、cmd/ 39`。
* 码点现量（我自己逐字节取，不是抄表）：前 4 枚里 `composer-states.html:2,5,8` 与 `composer.tsx:170`
  ＝ **U+2264**，`thinking.tsx:213`、`tool-chips.tsx:186` ＝ **U+2212**。**简报那句前提成立。**
* **旧正则不覆盖 ⇒ 新增红**：`git show ed2c077^:tools/d22scan/main.go | sed -n '111p'` 的字符类是
  `[1F000-1FAFF 2600-27BF 2B00-2BFF FE0F 1F1E6-1F1FF]`，与新类逐字符比对**只差 `\x{2200}-\x{22FF}` 一段**
  （箭頭 `2190–21FF`／带圈 `2460–24FF` 仍不在内）。反向实测见 §4(a)：把新段摘掉，
  `TestScannerSelfScanOfRealRepoIsGreen` 与 `TestRealRepoLedgerIsHonest` **立刻复绿**（6 枚 finding 归 0）
  ⇒ 这 6 行确实是本批点亮的，不是存量。

### 那两枚红＝"门能看见"，不是放水（简报点名要明写的一句）

* 两枚红是**同一个因**：`:269` 逐行点名那 6 枚（`t.Errorf`×6），`:1277` 是
  `t.Fatalf("HEAD must be green, rc=1")`。`:1273` 那句 `ban #6 examined 43` 是 `t.Logf`，
  **不是断言**，别读成"台账数不对"；且 `git show ed2c077^:...scan_test.go` 里它就已经是 `t.Logf`
  （`git log -L 1273,1274` 指回 `84e4161`，票 88），**不是本批降下来的**。
* 全仓**没有一处**为了变绿放宽断言：`ed2c077` 的测试 diff 里被删的断言行 **0 枚**、新增 **10 枚**（§5）。
* CI 侧后果（现量，不是推演）：`scripts/d22scan.sh` 是 `set -eu`，第一步就是这个正控 ⇒
  **正控红 ⇒ 第二步整树扫描根本不跑**，`.github/workflows/ci.yml:109` 调的正是它。
* **若 owner 在 `Q-48` 选"退回射程"（摘掉 `2200–22FF`）会不会自动复绿**：
  **这两枚会自动复绿**（`141-acc-mut-noband` 实测 PASS，见 §4(a)），
  **但另有 3 枚会红**：`TestBan8MathBandAndRemainingGaps`（U+2265/2264/2229/2212 四行）、
  `TestBan8CommentExemptionInGoSources/raw_string_that_reads_as_a_SQL_comment`、
  `TestBan8CommentExemptionInTextScopes/JSX_text_node_with_a_math_glyph`。
  ⇒ **"退回射程"不是一行 revert**：它要同时改回具名解冻②要求钉下的那三枚钉子，
  并让 `1218192`/`2970c79` 的清存量失去门禁支撑（生产串改成了 `>=`，却再没有任何仪器钉住"不许回到 `≥`"）。
  这一格只报形状，**不替 owner 选**。

## 2. 〔占位〕裁决格 1：豁免有没有吞掉本该算违规的东西

（本节在下面填，节末带本程 commit。）
