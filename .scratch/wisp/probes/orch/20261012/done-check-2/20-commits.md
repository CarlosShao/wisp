# `20-commits.md` — B 节：逐笔 commit 的越界名册（⛔ 区间尺，一笔一量）

尺（逐字，本腿现跑；输出＝本腿 `logs/dc2-R6-commit-roster.txt`／`logs/dc2-R6-legshas.txt`／`logs/dc2-R7-percommit.txt`／`logs/dc2-R8-suspicious.txt`，均 rc 自落）：
```
git log --format='%h %s' | grep -E '(^|[ (])30[49]-[a-z]'            # 名册＝自己数，⛔ 用派单给的枚数
git -c core.quotePath=false show --name-only --format= <每一笔>      # 逐笔，⛔ 区间尺
```

## 1. 分母（本腿自己数出来的，⇄ 编排者给的数）

- 名册总命中＝**42 笔**；剔除四笔纯编排者笔（`49153fa8`／`db1694bd`／`688a1e76`／`5f6323c0`）⇒ **腿笔 37 枚**。
  - 票 304＝**14 笔**：`304-a1` 2（`b5954339`、`93de4b4a`）／`304-r1` 5（`b9f37892`、`e875f00f`、`bd123485`、`27f7b197`、`04f2ecbc`）／`304-v1` 1（`c27928d6`）／`304-v1b` 6（`83b1aea9`、`245ab5d3`、`afd954ed`、`9fe10b73`、`7a9d00d7`、`b6080266`）。
  - 票 309＝**23 笔**：`a1` 3／`r1` 5／`r2` 4（含代提 `f6067e5d`）／`r3` 5／`v1` 1／`v2` 1／`v3` 1／`v3b` 3。
- ⚠ A856 §1／A855 §1 里那句"六笔／三笔逐笔越界"＝**只指最后一枚验收腿**，⛔ 是这两票的 commit 全体。本腿按全体量的。

## 2. 越界名册（全体 37 笔的归类结果）

合法落点集合（派单 §1.B 的口径）＝该腿**自家** `probes/<票>/<腿>/**` ＋ 该票票尾 ＋ 该票自认的实现件（304＝`internal/llm/**`，309＝`internal/tools/**`）＋裁决表 `docs/evidence/s1/**`（AGENTS §1.5 要求的落点）。

| 判 | 笔 | 路径 | 具名 |
|---|---|---|---|
| ★**越界一枚（跨腿目录）** | `f6067e5d`（＝编排者代提 `309-r2` 的最后九枚件） | `.scratch/wisp/probes/309/**v1**/logs/pathlen2.log`（2539 B） | 这枚件是 **`309-r2` 那一程的 pathlen 门禁读数**，却落在**另一枚腿 `309-v1` 的自家目录**里。`309-r2` 自己另有四枚 pathlen 件（`r2/logs/{final-pathlen,gate-pathlen-final,gate-pathlen,pathlen-start}.log`），所以这枚⛔ 是"v1 自己写的"（v1 笔 `47d3f6da` 早于它）。⇒ **A847/A849 那句"越界 0"⛔ 成立**：按"每笔只带自家路径"这把尺，这里恰有 1 枚。归因＝代提时把件放错了腿的目录（本腿⛔ 动，只报状态） |
| ⛔ 假阳性（记本腿自己的尺） | `47d3f6da` | `".scratch/wisp/probes/309/v1/logs/spotcheck-现量.log"` | 第一把尺没带 `-c core.quotePath=false` ⇒ 中文件名被转成八进制引号串、被我的分类器当成 FOREIGN。**同一条尺加了 quotePath 后命中＝0**（`logs/dc2-R8-suspicious.txt`，rc=0）。⇒ 名册尺⛔ 忘 quotePath |
| 其余 35 笔 | — | 只落自家 `probes/30[49]/<腿>/**`＋票尾＋自认实现件＋裁决表 | 逐笔原文在本腿 `logs/dc2-R7-percommit.txt`（37 笔全量，非抽样） |

实现件逐笔（本腿现取，`logs/dc2-R15-codecommits.txt`）：
- `e875f00f`＝4 枚 `internal/llm/**`（两枚 spawn 点同批＋两枚新 helper）⇒ 在 `AC#1` 授权面内（票面 `:26` 硬约束③逐枚点名了这两枚）。
- `16e86bae`＝`internal/tools/{assessorfor_taint_instrument_309_test.go, bridge.go, toolrequest_nail_309_test.go}`；`d75db9ee`＝`{bridge.go, toolrequest_nail_309_test.go}`；`df8d517d`＝`{toolrequest_nail_309_test.go}`；`a1f84d39`＝`{bridge.go}`。
- ★**具名判"那两笔动了 `bridge.go`（注释面）与枚数钉，是否在本票 `AC` 授权范围内"＝在范围内**，凭据逐枚：`AC#1` 票面 `:32` 自写"写面＝`internal/tools/*_test.go`"；枚数钉落 `AC#2`（`:33`）；`bridge.go` 的注释面由 `AC#0` 裁语落地节 `:70` 逐字授权（"Ⓑ 的注释（`bridge.go:920-923` 追加＋`:925` 上方具名）"），文案面由 `AC#3`（`:34`）授权；`AC#4`（`:35`）那句"写面⛔ 出 `internal/tools/**`＋`probes/309/**`＋票尾一条 bullet"是**上限**而非⛔ 碰 `bridge.go`。**⇒ 本腿判这两笔⛔ 越界；真正越界的只有上表那一枚跨腿目录件。**

## 3. ★注释面那三笔会不会挪别人的尺（本腿自己那一把，⛔ 用别人报过的数）

尺（逐字）＝
```
git show <h>^:internal/tools/bridge.go | grep -vE '^[[:space:]]*//'   # 剥**整行**注释（⛔ 剥成空行）
git show <h>:internal/tools/bridge.go  | grep -vE '^[[:space:]]*//'
cmp 两边；两边行数并排先证非空
```
读数（`logs/dc2-R10-bridge-comment-only.txt`，rc=0）：

| 笔 | 父版非注释行数 | 子版非注释行数 | `cmp` | 该笔 numstat |
|---|---|---|---|---|
| `16e86bae` | **854** | **854** | IDENTICAL | 73/0 ＋ 32/3 ＋ 134/0 |
| `d75db9ee` | **854** | **854** | IDENTICAL | 44/9 |
| `a1f84d39` | **854** | **854** | IDENTICAL | 14/4 |

⇒ 两边**都非空**（各 854 行）⇒ 比⛔ 是空⇄空的假清洁；三笔的注释面**没挪任何人的尺**（`bridge.go` 代码字节逐笔同一）。⚠ `16e86bae` 那 `32/3` 是**测试件**（`assessorfor…_test.go`）的增删，⛔ 属注释面，本表只量 `bridge.go`。

## 4. ★另抓一枚与"尺的面"有关的实际越界后果（票 304）

`bd123485`（`304-r1` 的 `AC#2`/`AC#3` 交件笔）**新建了 4 枚 tracked `.go`**：
`probes/304/r1/logs/landed/{childjob_other,childjob_windows,mockllm,mockllm_integ_test}.go`（尺＝`git log --diff-filter=A`，`logs/dc2-R15-codecommits.txt`）。
- 这与 `309-r2` 那笔**同名不同命**：`309-r2` 自己把 `probes/309/r2/logs/bridge-head.go` 改名成 `.go.txt`，具名理由是"会挪别家 gofumpt/d22scan 的 tracked 分母"，并交了 `find … -name '*.go'`＝0 那发。
- 本腿现量 tracked `.go` 落在 `.scratch/wisp/probes/**` 下的总数＝**318 枚**（尺＝`git ls-files '*.go' | grep -cE '^\.scratch/wisp/probes/'`）⇄ `309-v2` 记该笔欠账时的 314 枚 ⇒ **+4，正是这一笔长的**。
- ⇒ 后果＝任何"全仓 `--include=*.go`"的词面尺（含票 309 那枚钉的词面对照、票 304 `AC#0` 的 spawn 名册）多 4 枚污染面；`309-v2-6` 那笔欠账的**账数已过期**，且 304 侧**没记这一笔**。（四道门禁今天对它零反应＝与 `Q-309-v2-6` 原判一致，⛔ 因此当作零后果。）

## 5. Git 纪律面（本腿能核的）

- 禁区面（尺＝`grep -cE '(\.gitignore|go\.mod|^design/|internal/risk/|\.github/)' `打在本腿 37 笔逐笔 name-only 汇总件上，件＝`logs/dc2-R17-forbidden-paths.txt`）＝**命中 0 枚**（empty above＝0 枚；无 stderr）。`^docs/` 命中 **4 枚＝全是裁决表**（`ticket304-304-v1-accept-v1.md`／`309-v1-verdict.md`／`309-v2-verdicts.md`／`309-v3b-verdict.md`）⇒ 落点＝AGENTS §1.5 要求的 `docs/evidence/s1/`，⛔ 改别处 `docs/**`。
- 票尾 bullet（本腿现量 `git show --numstat --format=`，件＝`logs/dc2-R16-bullets.txt`）：`7a9d00d7`（304）＝**1/0** 一枚文件＝票面；`9404ef3b`（309）＝票面那一行 **1/0** ＋同笔另带 7 枚自家 probes 件（`44/0 56/0 3/0 5/0 252/0 20/0 14/0`）⇒ 两票的票尾都是 **append-only 纯新增**，✓ 与 A855/A856 那句"1 加／0 删"对得上。
- ⛔ push 与否⛔ 本腿能判（那是编排者的权，且要跑 `git log origin/dev..dev` 之外还要网络面读，本腿未做＝记 `Q-dc2-*`，见 `90-final.md`）。

rc=0
