# 票 193 · `193-a1` 只读普查表 —— "卡前端的门禁"逐枚点名（A1 腿）

> 性质：**只读普查**。本表不勾任何 AC 框、不动 `internal/**`／`cmd/**`／`go.mod`／`docs/PLAN.md`／`docs/specs/**` 一字。
> 一切数字均为本程现跑；抄自别家表的读数一律标〔转述，未复量〕。

## ① 起手锚＋写面闸门

- 派单锚点：`398e3344`｜本程起手 HEAD 现量：**`72078779`** ⇒ **已漂**，按派单 §3 首行处置（照实际 HEAD 做、登记、不停下等）。
- `git status --porcelain -- internal/ cmd/ go.mod` 起手现量（**两行，全是别人的在飞件，本程不动不提交不评论**）：
  - ` M internal/panel/git.go`
  - ` M internal/panel/git_test.go`
- `git status --porcelain -- design/` 起手现量：**16 删除／4 修改／11 未跟踪**（与派单 §2-Q1 预警同形；本程**不读 `design/**` 内容**，只报计数）。
- 全仓脏件行数现量：74 行。
- 写面本程只碰：本文件、`.scratch/wisp/probes/193/a1/**`、票 193 的 Progress log（追加）。

## ② Q1 门禁射程逐枚

### ②.1 `scripts/d22scan.sh` 起手现量（名册逐行照抄）

```
d22scan: examined 232 production Go files under internal/ and cmd/
d22scan: scope bans #1-5 internal/      examined 209 production Go files
d22scan: scope bans #1-5 cmd/           examined  23 production Go files
d22scan: scope ban #6 frontend/         examined  85 text files
d22scan: scope ban #7 internal/tools/   examined  21 production Go files
d22scan: scope ban #8 design/           examined  39 text files
d22scan: scope ban #8 frontend/         examined  85 text files
d22scan: scope ban #8 internal/         examined 436 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined  45 Go files, comments and _test.go included
verdict: clean - no D22 ban violations
```

- **`ban #8 internal/` 起手＝436，不是基线 433** ⇒ 按派单 §3 具名登记为**别人的在飞件/已落地新件抬上去的计数**，不当漂移、不去动它（本程未碰任何 `.go`）。
  旁证（现跑）：`git ls-files --others --exclude-standard -- internal/ cmd/` = **0 行**（无未跟踪 `.go`），`find internal cmd -name '*.go' | wc -l` = **481** = 436+45 ⇒ 436 是这棵树今天的真实文件数，不是脏件造成的读数污染。
- `frontend/` 分母：ban #6 与 ban #8 各 **85**；`design/` 分母：ban #8 **39**（本程**只报计数，未打开其中任何一枚文件**）。
- 扫描器自身那枚正向对照（`runtests.sh -C tools/d22scan ./...`）现量：`PASS=34 FAIL=0 SKIP=0`，退码 0。
- `scripts/d22scan.sh` 整体退码：**0**（clean）。
- 影子件：本次 stdout 抄件落在 `.scratch/wisp/probes/193/a1/`（见 §⑨）。

### ②.2 逐枚判定（谁真的会点前端红）

（本节在 §⑥ 之前逐枚补完：`tools/d22scan/main.go:249`（ban #6 walk）、`:557-564`（ban #8 `emojiScopes()` 四棵树的字面清单）、`internal/panel/frontend_hygiene_test.go`、`internal/panel/composer_test.go`、`internal/panel/bridge_test.go:135`、`internal/panel/approval_test.go:109`、`internal/panel/tokens_fourway_test.go`、`internal/ball/tokens_table_test.go:1364`、`tools/d22scan/scan_test.go` 的实树自扫。）

## ③ Q2 规格原句两列

（待补：`PLAN.md:64`／`:980`／`:2219`／`SPEC-08*` 的逐字摘句，两列＝"能约束前端视觉/动画的"／"本来就管不到前端的"。）

## ④ Q3 字段读者现量＋生效档推荐

（待补：`internal/config` 里 `[panel]`／视觉字段逐枚读者现跑；上游 18 枚零读者**复核**；D36 档位建议。）

## ⑤ Q4 关闭档可证性

（待补：Go 侧／测试侧能凭什么说动画真停了；没有机读证据就具名说没有。）

## ⑥ Q5 逐枚放开／保留表

（待补：一行一枚，含代价列与 `needs-approval` 标记。）

## ⑦ 我判错/没测的

（待补。）

## ⑧ 门禁终态

（待补：最后一枚 commit 之后重取。）

## ⑨ 被拒调用＋零删除自证＋工具调用终值

- 零删除自证：本程未执行 `rm`／`rmdir`／`git clean`／`git restore`／`git checkout .`／`worktree`；未跑 `probes/161/r6/flip-declaration.sh`；未改 `scripts/d22scan.sh` 与 `probes/154/gate-clauses.sh`。
- （待补：被拒调用、工具调用终值。）

## ⑩ next

（待补：哪几枚要 owner 先批、哪些可直接派写腿。）
