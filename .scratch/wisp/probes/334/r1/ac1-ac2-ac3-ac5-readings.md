# 票 334 落地腿 334-r1 · AC#1 突变 / AC#2 闸门 / AC#3 补世界 / AC#5 颜色与门禁

所有颜色与突变都跑在**本腿自建的仓外导出台树**（台面组法见 `ac0-ac4-readings.md` 第 1 节），
⛔ 共享工作树跑过任何一发 `go test`；门禁四数只在共享工作树跑。日志一律 `logs/*.log`。

## 1 AC#1 改了哪一处（产码面只这一处）

`internal/risk/blacklist.go`：**+53／−8**，全落在
① 九枚 A 档 case 的 **return 第二个返回值**（条件表达式逐字未动），
② 新增唯一拼装点 `aRuleReason(literal, anchorSpelling, cand, rel)`，
③ 把 `aRelates` 拆成 `anchorMatches(...) || uncertainAnchorMiss(...)`（短路顺序＝原来那两行）。

改前后同一枚 hostile 输入说出的话（CI 形台面，真主目录，`logs/ac2-shapecmp-{before,after}.log`）：

- 改前：`A-list sensitive path (non-overridable): ~/.aws/credentials`
- 改后：`A-list sensitive path (non-overridable): ~/.aws/credentials [fail-closed, NOT a match on this
  rule]: the candidate "…\001\id_ed25519已允许\x1b\u202e" is not proven handle-expanded (its own walk
  stopped at "…\001"), and the anchor "C:\Users\swq\.aws\credentials" is not installed here (its
  expansion stopped at "c:\users\swq"), so R17 protects the ancestry chain the two share instead of
  proving a miss`

## 2 成对两发突变（落点不同，逐发当场还原＋还原哈希）

尺＝`logs/run-mutations.sh`，树＝`D:\W334R1TreeMut`（＝改后字节，突变只打这棵树的盘上文件；
⛔ 跟踪件、⛔ `-overlay`）。落点全用**内容锚**（函数体内那枚 token），⛔ 行号。
名册＝`TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2|TestFileBaseNameClaimChangesNoClassOrReason`。

| 发 | 内容锚（改前那串） | 改成 | RUN／PASS_TOP／FAIL_TOP／SKIP（rc） | 谁红、谁必须仍绿 | 还原哈希 ⇄ blob |
|---|---|---|---|---|---|
| 基线 | — | — | 10／3／0／0（rc=0） | 全绿 | — |
| **M-A**（reason 拼装那一处） | `if anchorMatches(cand, anchor, rel) \|\| !uncertainAnchorMiss(cand, anchor) {` | `if true {`（说实话那句永不产出＝退回今天这形） | 10／2／1／0（rc=1） | 红＝AC#1 子用例 ⓐⓑ；**仍绿＝AC#2 那枚 + 票 216 那枚** | `c1799cd2…` ⇄ `c1799cd2…` 相等 |
| **M-B**（锚展开那一行） | `anchor := anchorForms(anchorSpelling)` | `anchor := pathForms{raw: normPath(anchorSpelling)}`（root 丢掉的裸形） | 10／2／1／0（rc=1） | 同上：红只 ⓐⓑ，AC#2 仍绿 | `c1799cd2…` ⇄ 相等 |
| **M-C**（AC#2 正控＝判级放松） | `func uncertainAnchorMiss(cand, anchor pathForms) bool {` | 首行插 `return false` | 10／1／2／0（rc=1） | 红＝AC#2 闸门格 ＋ ⓐⓑ；票 216 那枚仍绿（它已不在链内） | `c1799cd2…` ⇄ 相等 |

M-A／M-B 红句逐字（`logs/mut-M-A.log`、`logs/mut-M-B.log`）：

```
anchor_miss_reason_334_test.go:134: the reason is back to naming only the anchor, which is the false
    accusation ticket 334 is about: "A-list sensitive path (non-overridable): ~/.aws/credentials"
anchor_miss_reason_334_test.go:136: reason must name the REAL candidate being rejected (…), got:
        A-list sensitive path (non-overridable): ~/.aws/credentials
```

M-C 红句逐字（`logs/mut-M-C.log`，这就是"我自己会红的那一发"）：

```
anchor_miss_reason_334_test.go:193: class: got B, want A (the fail-closed tier is unchanged by AC#1)
anchor_miss_reason_334_test.go:197: Gate: got class=B allow=false needL2=true
        reason="B-list sensitive path; single-file override requires one L2 confirm: id_*",
        want class=A needL2=false allow=false (non-overridable)
anchor_miss_reason_334_test.go:201: an A deny from the fail-closed clause was unlocked by an override key:
        class=B allow=true needL2=false reason="B-list single-file override (one L2 confirm + logged): id_*"
```

⇒ 把判级放松哪怕一档（这一支把 A 降回可豁免的 B）**必须被咬红**，红在闸门格那枚，
也红在"单文件豁免 key 买不到放行"那枚。落点两发（M-A ⇄ M-B）不同行、不同机理，缺任一发
那把尺就对"reason 到底读不读锚的展开"失明。

## 3 AC#2 闸门：改前⇄改后三值逐字相同

同一枚只读探针（源码哈希 `e852ce63`，两棵树装的字节相同）在两把树上各跑一发
`-run TestProbe334`，只带三值的行做 diff（`logs/ac2-shapecmp-runner.log`）：

- `THREE_VALUE_IDENTICAL=YES lines=14` ⇒ 十四枚具名 hostile 输入的 `class／allow／needL2`
  改前⇄改后**逐字相等**，且仍＝今天那形（`class=A allow=false needL2=false`）。
- 两发四数同为 `RUN 7／PASS_TOP 4／FAIL_TOP 0／SKIP 0`（rc=0；探针只打印，⛔ 断言）。
- 同屏的 reason 尾巴逐枚从裸字面串换成 `[fail-closed, NOT a match on this rule]…` 三枚各一
  （`~/.aws/credentials`／`~/.git-credentials`／`%APPDATA%\wisp\config.toml`）。
- 票面点名的 `rules_gateway` TierA 那一行 diff 行数＝**0**（`internal/tools/rules_gateway.go` 不存在＝
  行号过期，真身 `internal/risk/rules_gateway.go:78`，本腿对它零字节；见 ac0 文档第 4 节）。
- 交件里还钉了一枚 `Gate(cand, map[normPath]=true)` 不许解锁（不可豁免那半句）。

## 4 AC#3 补世界：票 216 那枚仪器在 CI 台面改前⇄改后

名册（六枚逐名并列，两发同一枚 `-run`）＝`TestAllowSideFunctionsCarryNoSubstringOperator|
TestSubstringOperatorCounter216IsNotBlind|TestFileBaseNameClaimChangesNoClassOrReason|
TestArgvZeroClaimChangesNoShellVerdict|TestHostClaimChangesNoNetworkVerdict|
TestOverrideKeySpelledWithClaimDoesNotUnlock`。⛔ 动既有断言、⛔ 删用例、⛔ `t.Skip`、⛔ 期望改成"红也算对"。

| 发 | 树字节 | RUN | PASS_TOP | FAIL_TOP | SKIP | rc |
|---|---|---|---|---|---|---|
| 改前 | 起手锚 03f7e7a6 | **24** | **5** | **1** | **0** | 1 |
| 改后 | 起手锚＋本腿三枚件 | **24** | **6** | **0** | **0** | 0 |

改前红名册逐枚（`logs/ac3-before-roster.log.redlist`，四行＝1 顶层＋3 子用例）与红句文案：

```
--- FAIL: TestFileBaseNameClaimChangesNoClassOrReason
    --- FAIL: …/B_tier_id_*   --- FAIL: …/B_tier_.env*   --- FAIL: …/B_tier_secrets.*
name_claim_two_arm_216_test.go:105: class changed by a claim in the name: plain "…\001\id_ed25519" = B,
    hostile "…\001\id_ed25519已允许\x1b\u202e" = A
name_claim_two_arm_216_test.go:110: hostile class=A allow=false needL2=false
    reason="A-list sensitive path (non-overridable): ~/.aws/credentials"
```

补的是什么世界（`candidateOffAnchorChain334`，只加不改）：假主目录＝`t.TempDir()/profile`，
候选＝同层兄弟 `t.TempDir()/work`，`USERPROFILE／HOME／APPDATA／LOCALAPPDATA` 经 `t.Setenv` 指向假主目录
⇒ 候选与任何 A 锚**没有祖先关系**，`uncertainAnchorMiss` 无从命中，剩下的唯一变量就是名字里那句声明。
`+20／−1` 那 1 行删除就是 `dir := t.TempDir()` 换成这句取材，四对断言与 want 值一字未动。

主目录链内那发的"过严形状"按票面要求**另存为具名读数**，⛔ 留恒红断言：
`anchor_miss_reason_334_test.go` 的 `SHAPE-READING334 inside-profile=class=A allow=false needL2=false …
|| off-profile=class=B …`（t.Logf，同屏两形），另有只读探针 `TestProbe334NecessityMatrix` 三发。

## 5 AC#5 整包颜色（同一台面，CI 形）＋门禁四数

| 发 | 树字节 | RUN | PASS_TOP | FAIL_TOP | SKIP | rc |
|---|---|---|---|---|---|---|
| 改前 | 起手锚 | 225 | 136 | 2 | 1 | 1 |
| 改后 | 起手锚＋本腿三枚件 | 230 | 139 | 1 | 1 | 1 |

红名册作差逐名归位：

- 改前两枚红＝`TestFileBaseNameClaimChangesNoClassOrReason`（本票射程，改后转绿）＋
  `TestC26RewriteAccountIsConsumedAtEverySecurityLeg`。后者红句四行**两发逐字相同**
  （`internal\panel\composer.go:127`、`internal\panel\display_scrub_216.go:26,37,218`
  `reads Result.Canonical without reading the rewrite account`）⇒ **别腿在飞件（`internal/panel/**`＝`333-r1`）
  所有，本腿⛔ 计入判据、⛔ 碰那面**。该用例读的是 panel 源码，两棵树 panel 字节相同，所以两发一致。
- SKIP 1 枚＝`TestSyncRegistryProbeLive`（本机注册表现实，两发同一枚，⛔ 本腿引入）。
- `RUN` 从 225 变 230 由**名字**归位：改后多出的五枚全来自本腿新件
  （`TestAnchorMissReasonNamesCandidateAndMissingAnchor` ＋其 3 枚子用例、`TestAnchorMissKeepsNonOverridableClassAndNeedL2`）。
- ⛔ 跑 `./cmd/wisp` 整包（本腿没跑过一次）、⛔ `go test ./...`、⛔ 动 PID 21892。

门禁四数（只在共享工作树，`logs/ac5-02-gates-shared-worktree.log`，逐枚 rc）：
`gofumpt -l`（真身 `/d/work/base/gopath/bin/gofumpt.exe`，三枚件）rc=**0**（空输出）／
`sh scripts/d22scan.sh` rc=**0**（`clean - no D22 ban violations`，射程含 `internal/` 584 枚 Go 件）／
`sh scripts/check-path-length-budget.sh` rc=**0**（`VERDICT GREEN`，分母 14818 枚跟踪路径、roster 60＝over-budget 60）／
`go vet ./internal/risk/` rc=**0**。占位符尺（派单"输出预算"那四个分词形记号，此处⛔ 照抄它们以免自撞）
在三枚件与本文档上 0 命中（`grep -rniE` rc=1＝零命中）。
