# T17 对抗验收报告（编排者执行）

> 执行者：orchestrator（实现者 T17-impl 独立）。时间：2026-09-20T00:27:30Z（截止窗口内联验收）。
> 注：全包测试绿（0.28s）+ race 绿（1.32s）——T17/T18 两棒文件合流后联合验证。

| # | 项 | 裁决 | 证据 |
|---|---|---|---|
| 1 | R1–R9 全实现 | PASS | assessor.go（契约冻结块+融合 max+R9 recover）+ rules_gateway（R1–R4 接口注入）+ rules_network（R5 SSRF/白名单/URL>2048）+ rules_shell（R6 冻结元字符集）+ rules_scale（R7 ≥50）+ rules_irreversible（R8 未知类别 fail-closed） |
| 2 | 测试真实性 | PASS | 14 测试（正/反/边界、panic 注入×2、Deny 压 L2、send 压声明 L0、R4 会话覆盖阻断、10 个 golden 快照）——全包合流后全绿 |
| 3 | R9 fail-closed | PASS | panic 注入×2（插件判定器 + 已接线依赖）→ L2 |
| 4 | R4 语义 | PASS | SessionOverrideBlocked 字段（D45 会话授权不可覆盖 taint 升级） |
| 5 | 越界/D22 | PASS | 提交仅 assess/rules 文件（T18 的 pathresolver/blacklist 零触碰）；零 emoji |
| 6 | 票据对照 | PASS | 接线剩余清单已登记（18 With* 接线、19 taint、20 红队矩阵） |

**VERDICT: PASS**

## 更正（A30，2026-09-20，编排者授权）

本节**只追加**：上表第 5 行原文保留、一字未删，未勾/退任何框，未改票面 Status。
第 5 行是一行**三无 PASS**：既没写命令、也没给 `file:line`、也没给 commit SHA。
它其实捆了三个子主张，逐个补凭据后**结论不一样**，故分列如下。

### ① 被更正的原句（本文件第 12 行，原文照抄）

> | 5 | 越界/D22 | PASS | 提交仅 assess/rules 文件（T18 的 pathresolver/blacklist 零触碰）；零 emoji |

### ② 子主张 (a)「越界：提交仅 assess/rules 文件」——可复现，今天复现成立

票 17 的提交 SHA 在票面 Progress log（`.scratch/wisp/issues/17-risk-assessor-c19-done.md:46`），
报告本身没抄过来。补上后当场可跑：

```
$ git show --stat 67ffbd8
 ... 9 files changed, 1222 insertions(+), 1 deletion(-)
 .scratch/wisp/issues/17-risk-assessor-c19.md |   3 +-
 internal/risk/assessor.go                    | 353 +++++++++++++++++++++++++++
 internal/risk/assessor_test.go               | 182 ++++++++++++++
 internal/risk/rules_gateway.go               | 115 +++++++++
 internal/risk/rules_irreversible.go          |  48 ++++
 internal/risk/rules_network.go               | 128 ++++++++++
 internal/risk/rules_scale.go                 |  26 ++
 internal/risk/rules_shell.go                 | 104 ++++++++
 internal/risk/rules_test.go                  | 264 ++++++++++++++++++++
EXIT=0
```

⇒ 9 个文件全在 `internal/risk/`（assessor + rules_*）加票面自身；T18 的
`pathresolver*.go` / `blacklist*` 一路未触碰。**该子主张成立**（但当时未留凭据，属"事后补上的凭据"）。

### ③ 子主张 (b)「D22」——两条命令并排，错命令确实不扫任何东西

```
$ cd "D:\work\workspace\projects plans\Wisp" && go run ./tools/d22scan -root .
main module (github.com/CarlosShao/wisp) does not contain package github.com/CarlosShao/wisp/tools/d22scan
EXIT=1

$ cd "D:\work\workspace\projects plans\Wisp\tools\d22scan" && go run . -root ../..
d22scan: examined 194 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp
d22scan: clean - no D22 ban violations, no emoji in design/ or frontend/
EXIT=0
```

第一条（仓根调用，`tools/d22scan` 是独立 module 所以父模块里没有这个包）**不产生任何 findings 输出、
也不产生任何"扫了几文件"的自报**：用它得到的"clean"等于没测。第二条自报 `examined 194` + exit 0。
扫描器非空跑的阳性对照（`cd tools/d22scan && go test . -count=1 -v`，全 PASS、exit 0，
含 `TestScanDetectsAllSeededViolations` 与 `TestCheckRootRejectsBlindRoots/the_d22scan_module_itself`，
后者明确拒绝"扫描根看不见任何生产 Go 文件"的盲跑）。
**测量条件**：②③④跑于 2026-09-20T15:54Z，工作树含另一代理在途的 67 个 gofumpt 未提交 `.go` 改动；
本更正未执行任何格式化、未 stage 任何 `.go`；**未出现意外 findings**。

### ④ 子主张 (c)「零 emoji」——工具口径要先讲清，再给真命令

`d22scan` 的 emoji 禁令（ban #8）**只管 `design/` 与 `frontend/`**（`tools/d22scan/main.go:21、150-153`），
所以"票 17 提交零 emoji"从来不在 d22scan 的覆盖面里，不能拿③的 clean 给它背书。直接扫该提交：

```
$ git show 67ffbd8 | grep -cP "[\x{1F300}-\x{1FAFF}\x{2600}-\x{27BF}\x{FE0F}]"
0        （grep 无命中 ⇒ EXIT=1）
```

同一正则同一管线的阳性对照（换个确实含 emoji 的提交就会翻出来，证明这不是"管道断了也报 0"）：

```
$ git show 5aa3258 | grep -cP "[\x{1F300}-\x{1FAFF}\x{2600}-\x{27BF}\x{FE0F}]"
9        （EXIT=0）
```

⇒ **该子主张成立**（票 17 提交内零 emoji）。

### 重新判定

| 第 5 行的子主张 | 当时的证明力 | 今天的判定 |
|---|---|---|
| 越界：仅 assess/rules 文件 | **零**（无 commit SHA、无命令） | **成立**，凭据由本节 ② 补上（`git show --stat 67ffbd8`） |
| D22 | **零**（无命令；唯一可推断的仓根调用覆盖面为 0） | **无凭据**：今天 `examined 194 / exit 0` 只证明今日树干净，**不追认** 2026-09-20T00:27:30Z 那一刻，故不写 PASS |
| 零 emoji | **零**（无命令，且工具口径本就不覆盖 `internal/`） | **成立**，凭据由本节 ④ 补上（含同正则阳性对照） |

- **这一行按 README 规则 6 的口径不再是一行可采信的 PASS**：三个子主张里两个今天补上了真凭据，
  D22 那半**仍是空的**，且它当时也是空的。
- **对票 17 整体的影响范围**：只涉及第 5 行的 D22 半句。第 1–4、6 行（R1–R9 实现与规则文件、
  14 条测试、R9 panic 注入×2、`SessionOverrideBlocked` 字段、接线剩余清单登记）各有自己的
  `file:line` 或测试证据，**不因此节而失效**。票 17 的 done 状态本节不置可否——
  改票面（含把本节的 D22 半句登记为待重跑）归编排者。
- 若要给"D22 在 S1 期间确实干净"补一个**有据可查**的锚点，可引用的不是本报告，而是票 67 的
  正确调用记录（`examined 194`、exit 0）与本文 ③。
