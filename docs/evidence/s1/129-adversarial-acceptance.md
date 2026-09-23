# 票 129 — 独立对抗验收（AC#1..AC#5）

裁决方：`acceptor-ticket129-r1`。我不是实现者，下面每一个数都是我自己打的。

## 0. 锚定、快照与被验版本自证

| 项 | 读数 |
| --- | --- |
| 开工时 `git rev-parse HEAD` | `3029415284f80739b3e53e28126379ae5ecb4df0`（＝简报所说的 `3029415`，成立） |
| 被验版本（锚定 sha） | **`ea05cf5`** |
| 我实际用的快照目录 | **`/tmp/wisp129-acc-r1`** ＝ `C:\Users\swq\AppData\Local\Temp\wisp129-acc-r1`（`git archive ea05cf5 \| tar -x`，仓库内零 worktree / 零 checkout） |
| 快照保真自证 | `cmp` 快照 `internal/winsec/resolve.go` 与 `git show ea05cf5:internal/winsec/resolve.go` ⇒ **字节全等**；`.gitattributes` 在归档内（`*.go text eol=lf`） |
| `go` 版本 | `go1.27.1 windows/amd64` |
| 在飞 CI | 全程有 self-hosted run 在跑（`35819355656` → `35822181825`）⇒ 本机读数与 CI 同机争 CPU。本轮**所有**四数读数均与预期同形、无 `0xc000013a`/`0xc0000135` 早死，故不取任何时序结论。 |

### 0.1 简报里三条断言的复核（有一条不成立）

1. ⚠ **不成立**：简报说「`ea05cf5..HEAD` 差异全在 `cmd/wisp/` 与 `internal/observe/`」。
   实测 `git diff --name-only ea05cf5..HEAD` **只有六枚文档**：
   `.scratch/wisp/issues/{130,131,133}-*.md`、`docs/evidence/s1/130-ac3-a-with-mirror-implementation.md`、
   `docs/reports/HANDOVER.md`、`docs/reports/pending-and-issues.md`。
   `git diff --stat ea05cf5..HEAD -- '*.go'` ⇒ **空输出**；`-- cmd/wisp internal/observe` ⇒ **空输出**。
   ⇒ 结论方向反而更有利：**`ea05cf5` 与 `HEAD` 的全部 Go 源码逐字相同**，我在 `ea05cf5` 上读到的东西
   对 `HEAD` 同样成立。简报那句话是错的，但不影响被验版本的有效性（按"不成立就报回来、不硬改"处理）。
2. 成立：CI run `35817761098` 的 `go vet (module)` 真实报错是
   `vet: cmd/wisp/leg_sink_nail_131_test.go:127:12: undefined: sinkInstallRecord`
   （我从 `gh run view --job 107042866455 --log` 逐字读出，同发 `staticcheck` 另点 5 处 undefined）。
   ⇒ 那是我们自己仓库里的真类型错误，**不是** cgo 交叉编译假象。详见第 5 格。
3. 成立（但只覆盖一部分，见 0.2）：`git diff --stat a71b2d8..HEAD -- internal/winsec` 确为**空输出**。

### 0.2 票 129 的完整改动面（我自己算的，不接受任何人给的范围）

`git log --oneline -- internal/winsec` 里属于 129 的是两枚；129 全部六枚 commit 是我按
`git log --oneline ea05cf5` ＋ 逐枚 `git show --stat` 认出来的：

| commit | 内容 | 碰了什么 |
| --- | --- | --- |
| `a45b2e9` | AC#1＋AC#2 | **`internal/winsec/resolve.go`（生产码）** ＋ 两枚新 `_windows_test.go` ＋ 票面 |
| `64f4811` | AC#4 STEP 0 | 票面 only |
| `db9fafc` | AC#3 | `volume_attribution_126_windows_test.go` ＋ `absoluteness_attribution_129_windows_test.go` ＋ 票面 |
| `e10ca09` | AC#4 | 票面 ＋ `docs/evidence/s1/129-ac4-ac5-mutation-and-gates.md` |
| `4785ae7` | AC#5 | 票面 ＋ 同一枚 evidence |
| `ea05cf5` | 注入登记 | 票面 ＋ 同一枚 evidence |

⇒ **完整 winsec 改动面** = `git diff --stat 4824bb8..db9fafc -- internal/winsec` =
`absoluteness_attribution_129_windows_test.go +327`、`absoluteness_seam_landing_129_windows_test.go +451`、
**`resolve.go +53/-…`**、`volume_attribution_126_windows_test.go +43/-4`。

⇒ 简报那句提醒是对的：实现者"`git diff a71b2d8..HEAD -- internal/winsec` 为空"**只**证明
AC#4/AC#5 那两批没动生产码，**不证明整票没动生产码**。129 确实动了 `resolve.go`：
新增 `sameAbsoluteness`（`resolve.go:519`），并在 `sameTree`（`:456`）与 `answerInsideTree`（`:552`）
各加一枚 `|| !sameAbsoluteness(…)` 提前返回。`pathComponents` 一字未动（边界③守住，我核过 hunks）。

### 0.3 基线重量（我自己的，不引用实现者的数）

`/tmp/wisp129-acc-r1`（纯净、未变异）里 `go test -count=2 -v ./internal/winsec/` ⇒ **rc=0**：

| 读数（按 `scripts/winsec-tests.sh:97-100` 那四条 grep 的逐字形状） | 值 |
| --- | --- |
| `^=== RUN` | 202 |
| `^--- PASS` | 116 |
| `^--- FAIL` | 0 |
| `^--- SKIP` | 0 |
| `grep -c '(cached)'` | 0 |
| 顶层 `--- PASS` 去重名字数 | **58** |

⇒ 与实现者 STEP 0 自报的 202/116/0/0 与 58 枚**同形**，但我是在 `ea05cf5` 上自己数的（`-v` 输出、
`^` 锚定，注释与缩进的子用例不进分子）。

---

## 3. AC#3 —— `R-126-3` 那枚被问侧 guard（首要攻击点）

**裁决：PASS** ｜ 标签：〔独立复现〕＋〔独立新造〕

### 3.1 先回答简报的那两问

**问①：这枚 guard 在生产码里还是只在 `_test.go` 里？**

只在测试里。逐字定位：

- 定义：`internal/winsec/volume_attribution_126_windows_test.go:82` `func vouchedSpelling129(t *testing.T, spelling string)`
- 用在 `:173`（`…AnotherVolumeSpelling`，`noticesAboutTree(*got, planted)` 之前）
- 用在 `:249`（`…ASecondRealVolume/two_real_volumes`，`noticesAboutTree(*got, b)` 之前）
- 第二半：逐枚通知那一圈从裸 `noticeNamesTree` 换成 `answerNamesTree115`（`volume_attribution_126_windows_test.go:192`）
- 全仓 `grep -rn vouchedSpelling129 --include="*.go"` ⇒ **4 处命中，全在同一枚 `_windows_test.go` 内**（1 枚定义注释 + 定义 + 2 处使用）

**它算不算数？算。** 三条依据：

1. **没有别的地方可放。** 它的动作是 `t.Fatalf`（`:85`），`*testing.T` 的方法。生产码里不存在能承载
   这一枚判决的位置。所以"只在测试里"不是选了个弱位置，是**只有这一个位置**。
2. **它拒的那个东西只在测试里存在。** `ResolvePath` 答"不认这枚拼写"在**生产**里是被**故意**设计成
   `false` 方向的，原文在 `winsec_windows.go:98-102`：
   *"Failure direction is a false alarm rather than a false all-clear: a spelling ResolvePath will not
   vouch for attributes to no notice at all, so a caller asking 'was my tree reported on?' hears 'no'
   and goes looking"* —— 生产里 `false` 是安全侧；只有**测试**会把"没人应答"读成"两棵树不同"。
   被问侧那个"种植出来的拼写"（`Z:\…`）根本不是生产对象，是测试自己造的。
3. **按票 126 AC#1 那把尺量，结论也是同一侧。** 我独立复算 `noticeNamesTree`/`noticesAboutTree` 的
   非测试消费者：`grep -rn "noticeNamesTree\|noticesAboutTree" --include="*.go" internal/winsec \| grep -v _test.go`
   ⇒ 命中只有它自己的定义（`winsec_windows.go:106`/`:120`）、`:123` 的内部调用、以及两处**注释**。
   **零枚非测试消费者**（与票 126 验收方 1.2 同向）。⇒ 归属侧连事故面的生产通路都没开通，
   那一侧按尺子本来就"按事故面记"，而 AC#3 要修的比事故面还窄一格：它要的是**仪器自己在别的机器上
   别撒谎**。所以这枚 guard 打的是"我们自己别写错"，而票面 AC#3 的原话要求的**正是**这一格
   （"自证它挡得住'换台机器就什么都没比较而报绿'"），不是打攻击面。**它答的是被问的那一题。**

**问②：亲手造一次"换台机器什么都没比较而报绿"，门必须红。** ⇒ 造出来了，见 3.2。

### 3.2 三态读数（`MUT-ASKED-REFUSES`，我自己造的形，不是我抄它的）

变异内容（打在 `builtinVerifier.Resolve` 的 `IsAbs` 腿之后，两枚树上逐字同一发）：
"这台机器的 `ResolvePath` 不肯为一枚**未挂载的卷字母**作证" —— 这正是 R-126-3 点名的那台机器。
落地证明：两枚快照里 `grep -c MUT-ASKED-REFUSES` 各 1、`go build ./internal/winsec/` rc=0。

| 状态 | 树 | `go test -count=2 -v ./internal/winsec/` | 判决 |
| --- | --- | --- | --- |
| **改前（缺 guard）＋ 变异** | `f5bbccd`（票 126 交付原样，`grep -c vouchedSpelling129` = **0**、`grep -c sameAbsoluteness` = **0**） | **rc=0**，`RUN=182 PASS=100 FAIL=0 SKIP=0` | **假绿量成读数** |
| **改后（有 guard）＋ 同一发变异** | `ea05cf5` | **rc=1**，`RUN=202 PASS=114 FAIL=2 SKIP=0`，顶层红名去重 **1 枚** | **门红了** |
| **还原复绿** | `ea05cf5` ＋ `resolve.go` 还原（`cmp` 证 byte-identical、`grep -c MUT-ASKED-REFUSES` 回 **0**） | **rc=0**，`RUN=202 PASS=116 FAIL=0 SKIP=0` | 复绿 |

**"假绿"那一格的证据不是推理，是它自己的日志行**（`f5bbccd` ＋ 变异）：

```
=== RUN   TestNoticeFromOneVolumeIsNotAttributedToAnotherVolumeSpelling
--- PASS: TestNoticeFromOneVolumeIsNotAttributedToAnotherVolumeSpelling (0.04s)
AC#4 shapes: sealed=C:\Users\swq\AppData\Local\Temp\TestNotice...2334399730\001\store-44440\artifact.txt
             planted=Z:\Users\swq\AppData\Local\Temp\TestNotice...2334399730\001\store-44440\artifact.txt
```

⇒ 腿**照跑**（不是 SKIP）、`planted=Z:\…` 说明它种植的正是那枚这台机器"看不见"的拼写，
而它一次都没比较过任何东西就报绿 —— 因为 `noticeNamesTree` 的 `false` 全部来自"没人应答"那一问。

**"门红了"那一格的逐字红名**（`ea05cf5` ＋ 同一发变异）：

```
--- FAIL: TestNoticeFromOneVolumeIsNotAttributedToAnotherVolumeSpelling (0.07s)
    volume_attribution_126_windows_test.go:173: this leg cannot be asked at all: ResolvePath("Z:\Users\swq\...\store-44440\artifact.txt")
      refused to vouch for the spelling it was handed: winsec: refusing to seal Z:\Users\swq\...\artifact.txt:
      the installed risk.c26Pipeline answered "Z:\...", a spelling the floor itself refuses:
      winsec: path is not provably resolved, refusing to seal: Z:\...\artifact.txt names volume Z:, which this machine does not have
```

⇒ 红名**就是 guard 自己那句原文**（`volume_attribution_126_windows_test.go:85`），
且它顺带把"拒点在 `ResolvePath` 对答案重跑底线的那一腿"钉成了读数。

**反向对照（变异是定点的、没把整包打死）**：同一发变异下
`TestNoticeFromOneVolumeIsNotAttributedToASecondRealVolume/two_real_volumes` 在两枚树上都 `--- PASS`
（我这台机器有 `C:`/`D:` 两枚真卷，日志原文 `AC#8 shapes: A=C:\wisp126-xvol-43908\… B=D:\wisp126-xvol-43908\…`）。

### 3.3 我另外去攻、但**没**攻开的两处（如实写）

- **`:249` 那一处（真卷那枚）我造不出假绿。** 试着让 `ResolvePath` 拒绝 B 树 —— B 的尾段与 A 相同，
  所以 `SealFile(a)` 会先 `t.Fatalf`（`:241`），门红而不是假绿。也试过让
  `writableVolumeRoots126` 交出一枚不存在的卷根 —— `os.MkdirAll` 先失败（`:229`）。
  ⇒ `:249` 这枚 guard 今天**被两枚更早的前提间接保护**；它自己的注释写的就是这个
  （"this precondition costs nothing here and is the only thing that keeps it honest somewhere else"）。
  这枚 guard 不承重，但它是**唯一**一处直接读数，删了它不会立刻红、留着才会在需要时红 —— 判它成立。
- **同族残留：别的文件里还有两枚相同极性的裸调用**（`len(...) != 0` / `if ... { Errorf }` 期望 `false`）：
  `inherited_narrow_notice_104_windows_test.go:293`、`narrow_notice_windows_test.go:90`。
  我试着用同一发变异打红它们 ⇒ **打不出来**：它们问的 `k`/`inherited` 都是**同一枚测试刚创建并解析过的真文件**，
  不是种植出来的字母，这台机器与"看不见 Z: 的机器"两种形态下它们都真在比较。
  ⇒ 不构成本格的缺陷；作为**同族残账**登记为 `R-129-4`（低严重度、非本票 AC#3 的要求面）。

### 3.4 AC#3 附赠：本格还量出一枚 AC#2 的归属侧读数

`TestAttributionFaceNeverSeesAMixedAbsolutenessPair`（`absoluteness_attribution_129_windows_test.go:121`）
不是我读的结论，是我复跑绿的：真 seal 一发、逐枚 `filepath.IsAbs(n.Path)` 与 `ResolvePath(child)` 的答案
都绝对 ⇒ 新 leg 在归属面上**无输入可达**。与 3.1 问① 第 3 条（零枚非测试消费者）同向。

**⇒ AC#3 通过。** 票面那句"补上并自证它挡得住'换台机器就什么都没比较而报绿'"两半都成立：
补上了（4 处命中、2 个使用点），且我**独立地**把那一枚假绿在缺 guard 的树上造出来、在有 guard 的树上
打成红、还原复绿。

