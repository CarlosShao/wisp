# 票 235 r1（产码腿 235-r1）逐发读数台件

起手锚点 `2555687973fdd1130998e7f5bc1679715184db71`（HEAD，09-29 20:42 那枚 `ledger(A452 …)`）。
本文件是**读数台账**：每一发的命令、期望、真实读数、还原自证，逐条落在下面各节。
⛔ 本腿不翻票面任何 AC 框（勾由编排者核过之后翻）。
⛔ 零行为变更：本腿只动 `internal/tools/subagent_197_test.go`（注释段）与 `internal/tools/subagent_222_test.go`（测试面）。
所有突变一律 `go test -overlay`，工作树零写入（每发后的 `git status --porcelain -- internal cmd` 读数见 §5）。
本目录的 `mut-*/` 副本与 `overlay-*.json` 就是那六发的现场，只建不删。

## 0. 现量复跑（起手四把尺，本腿自己跑的）

尺 A —— `sed -n '396,404p' internal/tools/subagent_197_test.go`（改前）：注释起于 **`:398`**
（"Why a bigger pool is a lie and not extra capacity…"），被证伪那句在 **`:399-401`**
（"one in-flight spawn holds one bridge slot for its child's whole life (bridge.run keeps the semaphore
held across entry.Tool.Execute)"），函数名行 **`:403`**，比较方向行 **`:404`**
（`if MaxConcurrentSubagents > MaxToolConcurrency {`）。⇒ 与编排者给的读数逐字一致。

尺 B —— `grep -rn "bridge slot for its child|whole life" --include=*.go internal cmd`（改前）：命中 4 行——
`internal/tools/subagent_197_test.go:399`（**就是那枚注释本身**）、`internal/agent/spill_acl_windows_test.go:33`、
`internal/audio/wasapimic_windows.go:147`、`internal/tools/bridge.go:637`（后三处与这段注释无关）。
⇒ ⚠ **没有任何测试断言这段注释的文字**（这条预告在 §3 的 AC#3② 被量成读数）。

尺 C —— `grep -n "runs" internal/tools/subagent_222_test.go`（改前）：
`:115` 字段声明、`:153` 自增、`:374`/`:375` 现有读法（`runs := await222Tokens(...)` **之后**才读）、`:513` 第二处读法。
⇒ AC#2 要用的计数器**已存在**，本腿没有新造 `runs`，只加了一枚"到过桥"的信号（见 §2）。

尺 D —— `internal/tools` 基线名册（`-count=1 -v` 现跑，`PATH` 带 Sherpa，台件 `baseline-tools-v.txt`）：
`=== RUN` **218**／`^[[:space:]]*--- (PASS|FAIL)` 全名 **218**（含缩进子测试）／只数顶层 `^--- PASS` **169**／
`^[[:space:]]*--- FAIL` **0**／包级 `^FAIL\t` **0**／`ok … 12.764s`。
⚠ 票面 AC#4 记的"起手 161 顶层声明／161 PASS"是 `222-v1` 当时的数；今天同一把只数顶层的尺量出 **169**，
含缩进的全名尺量出 **218** ⇒ 本腿拿**现跑的 218** 当"零新增红"的分母，没为对上 161 改任何计数口径。

## 1. AC#1 注释理由改写（`internal/tools/subagent_197_test.go`，commit `1d2ad737`）

改后落在 **`:398-410`**（旧 5 行 ⇒ 新 13 行），函数名行推到 **`:411`**、比较方向行 **`:412`**。新文字（逐字）：

```
// Why the pool must still not exceed the ceiling. The reason this nail was
// filed with has since been falsified: ticket 222 AC#2 makes a spawn that is
// waiting for its child hand the bridge slot back (giveBackWhileWaiting in
// subagent_197.go), so a waiting parent holds none of the four permits and the
// ceiling caps no longer how many children may EXIST at once - they may, what
// the ceiling caps is how many of them may be executing a tool in the same
// instant. The conclusion stands on the other half of the same fact: a pool
// above the ceiling still admits children the machine cannot run at once, they
// queue for the permits while the roster prints every admitted row as 「在跑」 -
// which is exactly the 7/8-red shape ticket 197 leg A reported. Raise
// MaxConcurrentSubagents and this leg still goes red; ticket 222's M6 reading
// (pool 8 against a bridge of 4, four legs red) is recorded in
// docs/evidence/s1/222-spawn-holds-bridge-permit-v1.md.
```

禁法三条的遵守（本腿自证，`git show 1d2ad737` 的 diff **只有那段注释**，函数体零字符变更）：
- 没拆钉子：`Test197SubagentPoolNeverExceedsBridgeCeiling` 仍在（`:411` 起），两枚 `t.Errorf` 一字未动。
- 没放宽比较方向：`MaxConcurrentSubagents > MaxToolConcurrency` 原样，两个常量数值也没碰。
- 历史句原样留着：`which is exactly the 7/8-red shape ticket 197 leg A reported` 逐字在位，
  同段把它标成**出处**（并点了今日机制 `giveBackWhileWaiting` 与 M6 的出处文件）。

新注释最后一句"抬池 ⇒ 这枚钉子仍红"本腿没有照抄，现验过（M6 复跑，产码一枚字节未动、只 overlay）：

```
PATH=… go test -overlay .scratch/wisp/probes/235/r1/overlay-m6-pool8.json ./internal/tools/ -count=1 -v
```

读数（`m6-pool8-rerun.txt`）：**4 枚红**——
`Test197SubagentPoolNeverExceedsBridgeCeiling (0.00s)`、`Test197SubagentPoolCapsAtBridgeCeiling (0.00s)`、
`Test197FullPoolRefusesNextSpawnWithReadableReason (0.00s)`、
`Test222SpawnConclusionArrivesThroughRealBridgeChildren (3.00s)`。
⇒ 与票面现量 #3 的"M6 抬池 ⇒ 4 枚红"同形。⚠ 那枚 3.00s 的红来自 `subagent_222_test.go:488` 的结论正文计数，
按 `222-v1` 注② 的读法＝**测试构造的阻塞**（gate 要 8 枚同时进桥、桥只有 4 枚许可），不是生产形状，
本腿不拿它当"池 8 有害"的证据，也只用它支持注释里那句"这枚钉子仍会红"。

## 2. AC#2 三枚用例的前置读数（`internal/tools/subagent_222_test.go`，commit `a71f0be9`）

补法形状＝一次"**桥上的到达 vs 父任务的返回**"的因果读数，落在各枚用例第一枚 `await` **之前**；
先到"到达"⇒ 通过（不判任何东西、微秒级返回），先到"父任务返回"⇒ 当场 `t.Fatalf`。

新增件（全在 `_test.go` 内，零产码）：
- `probe222.arrived chan struct{}`（`:134`，缓冲 16）＋ `Execute` 里 `p.runs.Add(1)` 之后的非阻塞投签（`:172-176`）。
  这枚信号**不承载任何判据**，只报"这次执行真上过桥"；缓冲 16 大于本文件最多要投的 10 枚（反控那发的 6 枚宿主调用＋4 枚孩子），所以投签永远不会变成谁在等谁。
- `h.awaitChildOnBridge222(t)`（`:378-407`）：`select` 三头——到达／`h.parents`／`h222SafetyBound` 护栏；
  中间那头把取走的那枚结论**原样交回**（`h.parents` 缓冲 `2*MaxConcurrentSubagents`，同时最多 n 枚 ⇒ 回投不阻塞），
  所以后面收 n 枚结论的分母没被改动。
- `h.drainBridgeArrivals()`（`:409-421`）：只在反控那发用——宿主调用先投过签，`wg.Wait()` 之后清空，
  之后读到的每一枚才只可能是孩子的。
- 三个调用点：leg1 **`:444`**（在 `runs := await222Tokens(… h.probe.in, n)` 之前）、leg2 **`:536`**、leg3 **`:602` drain ＋ `:604`**。
- 文件头多了一段"票 235 AC#2 加的是哪把尺、为什么它微秒级就能判"（`:48-58`）。

为什么这一发不算放宽也不算假绿：修好的形状下"到达必然早于任何父任务返回"是**因果**
（父任务的结论来自孩子，孩子的结论又来自它的工具调用从桥上回来之后），所以不需要墙钟；
M3 形状下桥上一次都没到、父任务却照样各自拿到结论 ⇒ 当场红。那枚 30s 护栏退回成护栏，不再是判据。

**未修码（改前）的 M3 读数**——判据要"某发必须红"，所以先跑改前再动码：

```
PATH=… go test -overlay .scratch/wisp/probes/235/r1/overlay-m3-original.json ./internal/tools/ -count=1 -v
```

读数（`m3-before-fix.txt`）：`=== RUN` 218／全名 218／`--- FAIL` **1**——
`Test222SpawnConclusionArrivesThroughRealBridgeChildren (30.00s)`，红句逐字
`subagent_222_test.go:374: 孩子在桥上跑工具调用：只等到 0/4 枚，护栏到点（这一发真挂住了，不是读数不同）`；
另两枚 `Test222WaitingParentHoldsNoBridgeSlot (0.00s)`、`Test222CeilingStillCapsExecutedCallsWhileParentsWait (0.00s)` **绿**。
⇒ 与票面现量 #5／`222-v1` 注① 完全同形（**只有 1/3 看得见，且靠 30s 挂死护栏才红**）。

**改后（终态码）的同一枚 M3**：

```
PATH=… go test -overlay .scratch/wisp/probes/235/r1/overlay-m3-afterfix.json ./internal/tools/ -count=1 -v
```

读数（`m3-after-fix-rerun.txt`）：**3/3 全红，全 0.00s**——

```
--- FAIL: Test222SpawnConclusionArrivesThroughRealBridgeChildren (0.00s)   ← :444
--- FAIL: Test222WaitingParentHoldsNoBridgeSlot (0.00s)                    ← :536
--- FAIL: Test222CeilingStillCapsExecutedCallsWhileParentsWait (0.00s)     ← :604
```

红句（三枚同形，逐字片段）：`桥上的到达信号一枚都没有，而第一枚父任务已经拿到了结论（真桥上累计执行 0 次／0 次／6 次，
反控那一发里这些全是宿主调用，孩子一次都没上桥）：… 孩子的工具面绕开了真桥＝A420 的第 7 环，票 235 AC#2 这一发就是为它加的`。
⛔ 零枚来自 30s 护栏。⚠ 反控那发的"6 次"就是它自己的 6 枚宿主调用，措辞已按这个事实写清，没假装是孩子在跑。

M3 形状＝`ParentTools` 与 `Tools` **两行同时**改指 `&fake197Dir{}`（票面现量 #4 说过 `:288` 单独改会被
`subagent_197.go:284` 的 `opt.Tools = newSubagentToolChain(… from base)` 覆写 ⇒ 两行一起改，与 `222-v1` 的 M3 同形）。

顺带补的指向（票面 AC#2 要求）：票 222 `:35`（AC#1，已勾）下面追加了一行 `- 〔09-29 由票 235 追加，原句一字未动〕…`，
`git diff`＝**1 insertion(＋)**，那枚 `- [x]` 的原句逐字未动。

## 3. AC#3 两枚突变自证

**① 把 AC#2 新加的前置读数注释掉 ⇒ 期望 M3 退回"只有 1/3 看得见"——成立**
副本 `mut-ac3prereadoff-m3/subagent_222_test.go`＝在**终态码**上把 `h.awaitChildOnBridge222(t)` ×3 与
`h.drainBridgeArrivals()` ×1 注释掉，再叠加同一枚 M3。

```
PATH=… go test -overlay .scratch/wisp/probes/235/r1/overlay-ac3prereadoff-m3.json ./internal/tools/ -count=1 -v
```

读数（`ac3-1-prereadoff-m3-final.txt`）：`=== RUN` 218／全名 218／`--- FAIL` **1**——
`Test222SpawnConclusionArrivesThroughRealBridgeChildren (30.00s)`，另两枚 **0.00s 绿**。
⇒ 检测力确实是 AC#2 新加的那枚前置读数买的：去掉就退回 §2 那份改前读数（1/3＋护栏形）。

**② 把 AC#1 改过的注释再改回旧那句 ⇒ 期望"必须有测试红"——这一发本腿跑不出来**
副本 `mut-ac3-commentback/subagent_197_test.go`＝终态件里把 `:398-410` 换回旧那 5 行，其余一字不动。

```
PATH=… go test -overlay .scratch/wisp/probes/235/r1/overlay-ac3-commentback.json ./internal/tools/ -count=1 -v
```

读数（`ac3-2-commentback-final.txt`）：`goexit=0`／`=== RUN` 218／全名 218／**`--- FAIL` = 0**——
`Test197SubagentPoolNeverExceedsBridgeCeiling (0.00s)` **PASS**、`Test197SubagentPoolCapsAtBridgeCeiling (0.00s)` **PASS**、
`ok github.com/CarlosShao/wisp/internal/tools 13.308s`。**红名册为空。**

⇒ ⚠ **具名承认：这枚注释没有任何仪器保护。** 把已经被证伪的假理由原样写回去，整包 218 枚用例一枚都不红——
那枚钉子只比 `MaxConcurrentSubagents > MaxToolConcurrency`，不看注释写了什么（§0 尺 B 的 grep 已经预告了这件事）。
本腿走票面给的第二条出路：**写明"只能靠人读"**；⛔ **没有**为了让它变红而新写一枚"断言注释文字"的测试——
那正是本票要修的毛病（拿假检测力冒充有尺）。
该有什么尺（**留给编排者定案，本腿不自裁**）：注释里被证伪的那句其实是**关于实现的断言**
（"父等待占不占桥位"），而那件事**已经有行为尺**——`Test222WaitingParentHoldsNoBridgeSlot` 的
`len(bridge.sem)=0`（`222-v1` 的 M1 之下 0.00s 红）。缺的不是第二把尺，而是"**注释 ↔ 尺**"的对应关系本身，
任何仪器都读不到那层（除非动契约面或写一枚 grep 文字的测试，两样本票都禁）。
⇒ 结论：**注释文字＝人读级；行为面＝有尺级**。这一句已同时写进本台件与票面 Progress log。

## 4. AC#4 门禁与逐名（终态，工作树干净）

```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./internal/tools/ -count=1 -v   → head-tools-v.txt
```

- `=== RUN` **218**／全名（含缩进） **218**／只数顶层 `^--- PASS` **169**（缩进子测试 49 枚 ⇒ ⛔ 别拿顶层尺当分母）／
  `^[[:space:]]*--- FAIL` **0**／包级 `^FAIL\t` **0**／`ok … 13.206s`。
- 逐名对拉：`base-names.txt` ↔ `final-names.txt`（`sort -u` 后 `comm -23`／`comm -13`）⇒ **双向差集为空**（218↔218）。
  只按用例名集合比，没按枚数比。
- `"$GOPATH/bin/gofumpt.exe" -l internal/tools/` ⇒ **零输出**。
- `go vet ./internal/tools/` ⇒ **clean**（无输出）。
- `scripts/d22scan.sh` ⇒ `d22scan: clean - no D22 ban violations`（`internal/` 462 枚 Go 件含注释与 `_test.go` 在 ban #8 射程内；
  `runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76`）。
- 稳定性：`-count=3 -run 'Test222|Test197SubagentPool|Test197SpawnDescription'` ⇒ 18 名全 PASS／**0 红**（`stability-count3.txt`）。
- 数据竞争：`-race -count=2 -run 'Test222'` ⇒ **`DATA RACE` 0**、6 名全 PASS、`ok 1.173s`（`race-222.txt`）。
- 历史在册红：本腿射程内**零枚新增红**；`internal/panel` 4 枚＋`internal/ball` 1 枚属别人地界，本腿没跑那两个包的终态（见 §7）。

## 5. 每发之后的还原自证（`git status --porcelain -- internal cmd`）

| 发 | 现读 |
|---|---|
| 改前 M3（overlay，还没动码） | 空 |
| AC#1 落码后（未 commit） | `M internal/tools/subagent_197_test.go` ⇒ 已进 `1d2ad737` |
| 改后 M3（终态 overlay） | `M internal/tools/subagent_222_test.go`（本腿自己的在途件）⇒ 已进 `a71f0be9` |
| AC#2 commit 之后 | **空** |
| AC#3①（终态复跑） | **空** |
| AC#3②（终态复跑） | **空** |
| M6 复跑（overlay 改产码副本，工作树未动） | **空** |
| 终态门禁 | **空** |

⛔ 六发突变全走 `-overlay`，没有一发往工作树写过突变字节；副本只存在于本目录 `mut-*/`。

## 6. 本腿落下的 commit（各带显式 pathspec，零 push）

- `795ed767` `chore(235)` — 台件骨架＋票面 Progress log 起手一节。
- `1d2ad737` `test(235 AC#1)` — `internal/tools/subagent_197_test.go`（注释段）。
- `a71f0be9` `test(235 AC#2)` — `internal/tools/subagent_222_test.go` ＋ 票 222 `:35` 下的追加指向行。
- 第四枚＝本台件全文＋票 235 的 Progress log 结案一节（只带这两枚 pathspec）。

## 7. 没做完／留给编排者

1. **AC#3② 的另一支"把它转成一条常驻断言"本腿没做**：票面写的是"或写明只能靠人读"，本腿走了后一支（§3）。
   若编排者要留一枚常驻尺盯"注释↔机制"，只有两条路：动契约面（要 `A##`），或新写一枚断言注释文字的测试——
   后者正是本票禁的形状。**这一句要人拍，本腿不自裁。**
2. **票面四枚 AC 框本腿一枚没翻**（`- [ ]` 原样，也没新增框）。可核的最硬读数：
   AC#1 diff 只含注释段＋M6 复跑 4 枚红；AC#2 的 M3 由 1/3（30.00s 护栏）变 3/3（0.00s）；
   AC#3① 去掉前置读数即退回 1/3；AC#3② 零枚红（已具名承认）；AC#4 门禁全绿、逐名双向零差集。
3. **台账登记没做**：`docs/reports/pending-and-issues.md` 的 `A##` 追加属编排者的笔，本腿没碰台账。
   要落两枚：① 票 222 AC#1 的覆盖面补强已交（指向票 235），② §3 那句"注释无人保护＝只能靠人读"的具名承认。
4. **票面 AC#4 的名册基线数字已过期**（写的 161，本腿现跑 169 顶层／218 全名）。
   本腿按现跑的当分母；⛔ 没为对上 161 改任何计数口径。**改不改票面那句＝编排者的笔。**
5. **整包 `./cmd/wisp ./internal/...` 的终态本腿没跑**：本票 AC#4 的射程只点到 `internal/tools`＋三件门禁，
   本腿按票面跑足；`internal/panel` 4 枚＋`internal/ball` 1 枚历史在册红因此**未经本腿复核**。
   要整包终态（票 222 AC#6 那形）请另派发，且与本腿同撞 `internal/tools` ⇒ 串行。
