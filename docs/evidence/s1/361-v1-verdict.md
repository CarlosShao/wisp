# 票 361 裁决表（AC#2〔乙〕）—— leg `361-v1`，非实现者，只读＋四读现跑

- HEAD 起手 ＝ `f2ce92a6e`（`361 AC#2乙:attrib.sh ruler_a读index(LF)取代worktree(CRLF假红);…;A1518`）
- 时刻 ＝ 2026-11-04 20:0x +08；本机 `core.autocrlf=true`；gofumpt `v0.12.0 (go1.27.1)`；`sh` = GNU bash 5.3.9 cygwin
- 票面 ＝ `.scratch/wisp/issues/361-attrib-tracked-only-is-red-locally-but-green-in-ci-the-third-cross-shape-case.md`
- 纪律 ＝ 0 个票框翻转、0 字节 `.go`、0 次 push；只新增本件一枚文件。`attrib.sh`／票面／台账 **一字未改**。
- 判 ＝ **PASS**

## 一、范围核对（SCOPE）

`git show --stat HEAD` 逐字：

```
 ...y-but-green-in-ci-the-third-cross-shape-case.md |  9 +++++++
 .scratch/wisp/probes/161/r5/attrib.sh              | 29 +++++++++++++++++++++-
 docs/reports/pending-and-issues.md                 |  4 +++
 3 files changed, 41 insertions(+), 1 deletion(-)
```

三枚全部落在应许面内（探针脚本 ＋ 361 票面 ＋ 台账 `A1518`）。**无意外文件**。

- `attrib.sh` diff 只有 **一枚 hunk**（`@@ -582,7 +582,34 @@`），只动 `ruler_a()` 的 `A_OUT` 取值源；`--self-test` 的 `rosterctl`、parse-error `A_RC>1` 分支、`A_FILES` 的 `:N:N:` 剥离**逐字未动**。
- 票面 diff ＝ 追加「AC#2 实现（编排者 161轮…）」小节（记录，不改任何判定框）。
- 台账 diff ＝ 追加 `## A1518`（记录，只追加）。
- 未碰 `PLAN.md`／`SPEC-*`／`thresholds.go`／`golden`／`.github/workflows/**`／任何 `.go`。
- 复跑前后 `git status --porcelain | grep -v '^??'` **均为空**（工作树与索引无 tracked 脏）。

## 二、四读现跑（逐字）

### (a) `sh .scratch/wisp/probes/161/r5/attrib.sh --self-test`

**rc=0**，关键行逐字：

```
   PASS  want rc=2+roster-names-the-sample  got rc=2 rows=1  throwaway tracked tree left at: /tmp/wisp-275-r2-rosterctl.6s29j0
   PASS  want rc=1 verdict=TRACKED<CI-RED>   got rc=1  tools/d22scan/main.go (parse_err=0) <== TRACKED BYTES ARE DIRTY
   PASS  want ring=5 silent=3 named=3  got ring=5 silent=3  all named must-ring samples present
attrib.sh: --self-test cases=10 failures=0
attrib.sh: SELF-TEST GREEN - the classifier rings on a line it cannot attribute and stays silent on a line a real ticket owns, the ticket 275 roster control saw an A-ROSTER row naming its own sample, and the roster nails hold (registered counts match, every named must-ring sample is present)
```

判读：**GREEN，正控仍红得起来**。`rosterctl` 正控（rc=2＋A-ROSTER 点名）与 `TRACKED<CI-RED>` 正控都在改后仍 PASS ⇒ 读 index 之后「脏/坏」这一侧没有被一起读没。self-test 的抛树走 `git add`（进 index、不进 commit），恰好落在新路径上而被验证，不是绕过。

### (b) 新版 `--tracked-only`（干净树）

**rc=0**，日志全文 8 行，关键行逐字：

```
   (A) denominator: handed to gofumpt = 1089; total tracked .go = 1455; subtracted by the one named exclusion :(exclude).scratch/** = 366
attrib.sh: (A) tracked-tree  lines=0 files=0  rule: MUST be empty
attrib.sh: GREEN (tracked-only) - (A) empty over 1089 tracked .go files (1455 tracked in total; the one named exclusion :(exclude).scratch/** subtracted 366)
```

判读：**GREEN 0 脏 / 1089 枚**，与票面目标一致；分母没缩水（1089 ＝ 独立复算值，见第三节）。

### (c) 牙：临时 stage 一枚真脏 `internal/agent/zz_v1_tooth.go`

```
printf 'package agent\n\nfunc  T ( ){x:=1;return}\n' > internal/agent/zz_v1_tooth.go
git add -- internal/agent/zz_v1_tooth.go      # ls-files 确认：internal/agent/zz_v1_tooth.go
sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only
```

**rc=1**，关键行逐字：

```
   (A) denominator: handed to gofumpt = 1090; total tracked .go = 1456; …
   TRACKED-DIRTY	internal/agent/zz_v1_tooth.go	(last commit touching it: )	<== CI IS RED ON THIS
attrib.sh: (A) tracked-tree  lines=1 files=1  rule: MUST be empty
attrib.sh: RED (tracked-only) - tracked-dirty=1 files_dirty=1, over 1090 tracked .go files (1456 tracked in total; …).
```

判读：**红，且逐名点名**，`files=1`、分母 1090（＝1089+1）⇒ index 里新 stage 的那枚被读到，说明「读 index」不是「读 HEAD、看不见暂存」。

清理与洁净确认（逐条现跑）：

```
git rm --cached -q -- internal/agent/zz_v1_tooth.go ; rm -f -- internal/agent/zz_v1_tooth.go
git ls-files -- internal/agent/zz_v1_tooth.go      → 空
git status --porcelain | grep zz_v1_tooth          → 空
ls internal/agent/zz_v1_tooth.go                   → No such file or directory
git status --porcelain | grep -v '^??'             → 空
```

**未留脏。**

### (d) 反向牙：原版（HEAD~1）脚本

```
git show HEAD~1:.scratch/wisp/probes/161/r5/attrib.sh > /tmp/orig361.sh   # sha256=da44d40651fad481022bac555df72782d2c294516dee74409b3eeab18269ba26
ATTRIB_ROOT="$PWD" sh /tmp/orig361.sh --tracked-only
```

**rc=1**，关键行逐字：

```
   TRACKED-DIRTY	internal/agent/approval/pending_read.go	(last commit touching it: d03d166fe 10-04 10:37 CarlosShao)	<== CI IS RED ON THIS
   TRACKED-DIRTY	internal/agent/tools.go	(last commit touching it: 5413f46d4 10-04 09:19 CarlosShao)	<== CI IS RED ON THIS
attrib.sh: (A) tracked-tree  lines=2 files=2  rule: MUST be empty
attrib.sh: RED (tracked-only) - tracked-dirty=2 files_dirty=2, over 1089 tracked .go files (1455 tracked in total; …).
```

判读：**同一棵树、同一分母 1089、只换脚本**，旧读红 2 枚、新读绿 0 枚 ⇒ 差异**确实归因于本次改动**，而非树变了。两枚正是票面点名的 `pending_read.go` / `tools.go`。

## 三、机制独立复核（不只信脚本自报）

**前提成立性**（逐字节量，不靠叙述）：

| 文件 | worktree 字节／CR 字节 | index 字节／CR 字节 |
|---|---|---|
| `internal/agent/approval/pending_read.go` | 6366 / **131** | 6235 / **0** |
| `internal/agent/tools.go` | 7581 / **225** | 7356 / **0** |

⇒ worktree 确是 CRLF、index 确是 LF；`gofumpt -l` 对 worktree 点名这两枚，对 index 字节不点名。「假红」这一形是真的，不是托词。

**独立复算（我自己搭两套源，不跑票里的脚本）**：

```
git -c core.autocrlf=false checkout-index -a --prefix=$T1/   rc=0
  → 提取 .go = 1089 ; 期望(git ls-files 减 .scratch) = 1089 ; gofumpt -l 输出 = 空
git archive HEAD | tar -x -C $T2                             rc=0
  → .go = 1089 ; gofumpt -l 输出 = 空
```

判读：绿**不是空集假绿**（两法都 1089，分母没掉），也**不是 worktree 残留**（出树后 gofumpt 全空）；且换一条完全独立的 blob 源（`git archive`）结论一致。最长的 tracked `.go` 路径 90 字符、临时前缀 19 字符，远未触 Windows 长路径。

`.gitattributes` 逐字有 `*.go text eol=lf` ⇒ index 字节恒为 LF，这条不变式不依赖开发者的 `core.autocrlf` 取值。

## 四、判定

**PASS。** 理由三条，缺一不成立：

1. **目标达成**：改后本机 `--tracked-only` 从「红 2 枚纯 CRLF 假红」变为「rc=0 / 1089 枚 / 0 脏」，与 CI 读同形（LF）。
2. **不是靠读没**：三枚正控各自独立响——self-test 10/10（含 rc=2 roster 点名与 TRACKED<CI-RED>）、真脏 stage 后 rc=1 且**逐名**点名、原版脚本对同一棵树仍红 2 枚。绿与红都还在，只是把「假红」这一格拿掉了。
3. **范围干净**：只动 `ruler_a` 的 `A_OUT` 源一枚 hunk，票面/台账为记录追加，无契约、无门禁阈值、无 `.go`。

## 五、残留关切（不阻断本判，留给后续票）

1. **仪器的自述行已陈旧**：`attrib.sh:582` 那行仍打印 `== (A) tracked tree: git ls-files -z '*.go' ':(exclude).scratch/**' | xargs -0 …gofumpt.exe -l`，而 0 脏那条路径实际已改为「checkout-index → 临时树 → gofumpt」。日志自述与实际动作不再一致（本判的四读日志里可见该行）。属可读性/取证诚实度问题，非功能缺陷。
2. **守卫只防「全空」，不防「部分缺失」**：`A_TMP_N -gt 0` 只在提取出 0 枚时回退；若 `checkout-index` 部分失败（如 Windows 造不出的路径），分母会静默变小而仍可能报绿。当前实测 1089/1089、rc=0，无现实风险，但这是一条潜在假绿通道。
3. **度量面从 worktree 移到 index**：未 `git add` 的工作树改动从此对 (A) 不可见（旧版会红）。就「与 CI 同形」的目标而言这是对的（CI 读入库字节，self-test 也依赖 index 语义），但「本机 pre-commit 自查」的射程确实变窄了一格，值得票面写明。
4. **`find | xargs`（无 `-0`）**：临时树路径用 `find | sed | xargs`，原路径是 `-z | xargs -0`；含空格的文件名会被拆词。现仓 tracked `.go` 无一枚含空白字符（已扫），故是潜在项而非现行缺陷。
5. **本判未复跑** `scripts/d22scan.sh` / `check-path-length-budget.sh`：本 AC 只动 `.scratch` 下的探针脚本，不在两把尺的扫描面；且票面把门禁归 AC#3，不属 v1 射程。
