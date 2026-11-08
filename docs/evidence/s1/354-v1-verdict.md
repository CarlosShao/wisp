# 票 354 `354-v1` 裁决腿判语：AC#2 确定性仪器 ＋ AC#3 成对去牙（非实现者裁）

- 日期：2026-11-08 23:2x +08。腿＝`354-v1`（裁决腿；实现者＝`354-r1`，两格判语全部出自本腿**盘上独立复跑**，⛔ 信 `354-r1` 任何转述）。
- 锚：主树 HEAD **`b45ff3339`**；全部读数跑在 `git clone --shared` 导出树 `/d/tmp/wisp/20261108/v354/tree-v1`（HEAD 同 `b45ff3339`，含 `354-r1` 两笔 `393430dc9`/`921263112`）。
- 判据出处（三核现读，⛔ 信转述）：票面 AC#2＝`.scratch/wisp/issues/354-identity-recheck-refuses-the-delete-so-the-callers-name-keeps-an-attackers-object.md:33-34`；票面 AC#3＝`:36`；判据原文（票面 `:18` 具名指向）＝`.scratch/wisp/probes/353/v1/02-ticket-draft-secC.md:29-30`。
- 本腿写面＝本件＋票 354 Progress log 一行。**零产码、零去牙入库**：两枚牙只落在导出树、每发当场还原，主树 `internal/winsec/**` 字节零触碰。

## 一、起手三核（全过）

1. 派单点名的串/行**全部现读在**：仪器件 `internal/winsec/leftstanding_354_windows_test.go`（214 行）；去牙落点 `internal/winsec/pinnedleaf_342_windows.go:277-283`（复验块 `Lstat`⇄`SameFile`⇄`return`）＋`:284`（`_ = p.root.Remove(base)`）。
2. 引文出处核过：票面 AC#2/AC#3 原文、草案 SEC-C `:29-30` 判据原文、证据件 `docs/evidence/s1/354-r1-deterministic-teeth.md` 全文，均现读。
3. 红/绿取自导出树（`git clone --shared`），非主树在飞读数；主树 HEAD 现量 `b45ff3339`。

## 二、AC#2 确定性仪器 —— **成立**

判据＝「同输入必得同输出」，三条逐一在**码里**核（⛔ 听声称）：

| # | 判据 | 码内证据 | 判 |
|---|---|---|---|
| ① | 有无 `t.Parallel()`（有＝不确定源） | 仪器件 grep `t.Parallel` **0 命中** | 过 |
| ② | 有无靠时序/睡眠（`time.Sleep`／`Eventually`） | 仪器件 **0 命中**；全套件 grep 证实 sleep 只在 `ancestorworld_342_348`／`leafrace_342r2` 两族 racer 件里，本仪器复用的世界（`newWinDetSwapWorld342r2`／`installWindow`，宿主 `ancestorwindow_342r2_windows_test.go`）零 sleep、同步推进 | 过 |
| ③ | 断言是否逐字节 | 全在：`:87-90` 固定 64 字节长度断言（改常量即红）；`:100-101` `wantStanding` 由两个固定常量拼出；`:115` 计数门三条件（`present`＋逐字节 `wantStanding`＋own name `absent`）；`:134-137` swap-in 自名缺席断言；`:159-163` 幸存名逐字节比对 | 过 |

配套三件也现核在码里：

- **轮内序号**（票面 ：35「尺必须带轮内序号」）：`:104-105` 两枚夹具名各带发号 `shot<N>`，读数行 ：119-122 带同一发号——我手跑出的日志逐行可归发（shot1/2/3 各自文件名可辨）。过。
- **分子不借既有**（票面 ：34）：判据分子＝新计数器 `leftStandingForeignName`；`entered` 只作活性配对不作判据。过。
- **零新接缝**：世界 `newWinDetSwapWorld342r2(t)`（:97）＋`installWindow`（:108）逐字复用，注入走包内已声明接缝 `applyDescriptor`（`winsec_windows.go:35-38`）。过。

**本腿独立复跑（裁决的全部，非抄读数）**：导出树 `go test ./internal/winsec -run '<两枚 354 测试精确名>' -count=3 -v` ⇒ **3 轮 × 3 发全 PASS**，逐发读数逐字节同：

```
354r1 AC#2 over 3 shots: entered=3/3, leftStandingForeignName=3 (the caller's name standing after the refusal, carrying the foreign object as its ONLY name)
```

每发 `the caller's name=present (size 64 body [an object the caller never named, standing in the same directory])` 且 `the swapped-in object's own name=absent`；负控 `name=absent`。与证据件 `04` 号日志**逐字对上**。同输入同输出成立（日志 `/d/tmp/wisp/20261108/v354/01-positive-count3.log`）。

**两点如实记**（不影响判）：

1. 派单给的正控模式 `-run 'Test354|LeftStanding'` 按字面**匹配零测试**（现跑留证：`testing: warning: no tests to run`——测试名里无 `Test354` 连串、无 `LeftStanding` 串，后者只在件名）。本腿改用两枚测试精确名实跑。
2. 「`-count=3` 九行全同」只作旁证（本腿同样量到全同）；**判据是逐字节断言在码里**（上表③），那才是同输入必得同输出的构造性保证。

**判：AC#2 成立。**

## 三、AC#3 成对去牙 —— **成立**（附一处措辞层澄清，非偏差）

票面 ：36 原文：「摘复验那步 ⇒ 回到 `damage>0`；摘删除那步 ⇒ damage 变体各自红；落点⛔ 同（⛔ 单发当判据）」。判据原文（草案 SEC-C `:30`）：「摘复验那步 ⇒ 回到 `damage>0`；**摘删除那步 ⇒ 正控转红**」。四问四答，全部本腿手跑（牙只在导出树，日志 `02-toothA.log`／`03-toothB.log`）：

1. **落点是否真不同**：是。牙 A＝`:277-283`（整块复验：`now, err := p.root.Lstat(base)` / `if err != nil || !os.SameFile(pinned, now)` / `return`）；牙 B＝`:284`（`_ = p.root.Remove(base)` 单行）。同函数内两个不同语句块，摘除后函数体我手 sed 现读核实。**落点⛔同成立**。
2. **牙 A 复跑（摘复验）**：分子 **3/3 → 0/3，断言红**。红句逐字：
   `354r1 AC#2 RED: leftStandingForeignName=0 of 3 entered worlds - the ratified residual moved. Form 乙's price is that the caller's name STAYS under the refused write and stays carrying the object the caller never named; a landing that changed either half changed what stands under a caller's name without saying so`
   且 353r1 主案**回到 `damage>0`**：`unnamed objects left with no name (DAMAGE)=3`，红句逐字：
   `353r1 AC#2 RED: 3 of 3 entered worlds left the swapped-in object with NO name, so the refusal-path removal ran on a name that carried an object the caller never named - ticket 353's damage, on the platform where the route's own handle makes that name vacate-able (bench row S4)`
   两枚控制照旧 PASS。**「摘复验⇒回到 damage>0」命中，与证据件逐字对上。**
3. **牙 B 复跑（摘 `:284`）**：**两枚无 top-up 控制齐红**（354r1 控制红句逐字：
   `354r1 AC#2 RED: with NOTHING topped up the refused write still left the name standing (present: size 0 body []), so the refusal-path removal no longer runs - the numerator above would then be counting a removal that cannot happen, and AC#5's refused write now leaves a half-made file`
   ；353r1 控制同形红），**主分子照绿**（`leftStandingForeignName=3`）、353r1 主案照绿（DAMAGE=0）。
4. **「主分子照绿」票面允许吗——现读判**：允许。★理由在判据文本本身：草案 SEC-C `:30` 对牙 B 的要求逐字是「**正控转红**」，而草案 AC#2 `:29` 定义正控＝「没人顶名的世界里名字仍必须消失」的那条读数——正是这两枚转红的控制。票面压缩语「damage 变体各自红」＝「量『删除那步』的各条读数各自红」，两枚控制即其全部在场的成员。主分子量的是「复验拒后名字下站着谁」：摘掉 `:284` 不改变复验留下的形（复验在 `:283` 已 return，根本走不到 Remove），所以主分子照绿是**测量对象决定的必然**——若牙 B 下主分子也红，反而说明它断错了腿。仪器靠成对控制补上这一腿＝**成对判据完整**。⇒ **非偏差**。
   ★一处措辞层澄清（如实记，不改判）：证据件把这两枚控制称「负控」，草案/仪器件内注释称「正控」——同一对象两套称呼，读数与红绿无差；本判语按草案口径记「正控（no-top-up 控制）」。
5. **⛔ 单发当判据**：每发去牙跑的是仪器内置 3 发循环（354r1 三发＋353r1 三发），红的是**三发全体**＋汇总行，非单发。过。
6. **还原**：两牙各跑后从 pristine 拷贝还原，`git status --porcelain` 空、`git diff --stat` 空、整包复绿。落 commit 字节零去牙（主树从未载入任何去牙形，本腿现量确认）。

**判：AC#3 成立。**

## 四、门禁（导出树 `b45ff3339`，跑前 `unset GOTMPDIR TMPDIR TEMP TMP`＋`PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`，CGO_ENABLED 未设）

| 门禁 | 本腿读数 |
|---|---|
| gofumpt（`/d/work/base/gopath/bin/gofumpt.exe -l internal/winsec/`） | 空，rc0 |
| `go vet ./internal/winsec` | rc0 |
| `go build ./...` | rc0 |
| `go test ./internal/winsec -count=1` | `ok github.com/CarlosShao/wisp/internal/winsec 35.199s` rc0 |
| d22scan（导出树现编 `tools/d22scan/d22scan.exe -root .`） | `clean - no D22 ban violations` rc0 |
| 两牙还原后 `git status --porcelain` | 空 |

## 五、总裁决

- **AC#2 成立**；**AC#3 成立**。两格均非实现者自证：确定性判在码内三判据＋本腿手跑读数；成对判在两落点现摘现红＋判据原文（草案 SEC-C）对「摘删除⇒正控转红」的逐字命中。
- ⇒ **编排者可以翻 AC#2／AC#3 两勾**（本腿⛔ 翻框）。
- 无「实现者自证」项需要退回；唯一具名的差异是上文第三节第 4 条的「负控/正控」称呼不一致（措辞层，无实质）。
- 后续：AC#4（颜色与门禁其余半格：两平台成对＋逐名红册）、AC#5、AC#6 归收口腿。
- 撤销口令：主树 revert 本腿这笔 commit 即回原状；导出树与 `/d/tmp/wisp/20261108/v354/` 全为临时件、不入库。
