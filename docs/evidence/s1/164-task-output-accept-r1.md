# 票 164 · AC#2＋AC#3 独立验收表 v1（非实现者，只裁不改）

- 派单＝`.scratch/wisp/dispatches/2026-09-27-195x-accept-164-r1-v1-non-impl-ac2-ac3.md`（提交 `e4d8e047`）
- 票面＝`.scratch/wisp/issues/164-background-jobs-cannot-be-read-fix-the-roster-then-add-the-output-leg.md`
- 被告陈述＝`docs/evidence/s1/164-task-output-impl-r1-ac2-ac3.md`（224 行）——**本表每一个数都是自己重跑的，没有一个从它那里引**
- 本件性质＝裁决位；**只 commit、不 push**；**AC 框一枚不勾**（勾归编排者）
- 台件＝`.scratch/wisp/probes/164/v1/ac2-anchor-reading.sh` ＋ 它本次的 `.log`（**仓外副本只落在 `/tmp`，仓库目录内没建任何 worktree/checkout**）
- 分档：〔现跑〕＝本程同树跑过｜〔锚点〕＝在 `f206e9f` 解出的仓外副本上跑过｜〔读码〕＝读了产码/测试本体、没跑

---

## 1. step-0 五件〔现跑〕

```
date                                          -> Sun Sep 27 19:50:55 CST 2026
git rev-parse --abbrev-ref HEAD               -> dev            （＝派单要求，未停手）
git rev-parse HEAD                            -> e4d8e047fce4829fab3ebfe0431ba5d5a1775249
git status --porcelain -- internal/ cmd/ docs/evidence/s1/
   -> " M docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md"   ← 【非空】
```

⚠ **派单 §0 那条"`git status --porcelain -- internal/ cmd/ docs/evidence/s1/` 必须空"起手即不符**。
不符的是 `docs/evidence/s1/` 里**别人的一枚证据件**，`internal/` 与 `cmd/` 两栏是空的。
按派单同一句的后半（"非空报回别动"）与本仓规矩处理：**不提交、不还原、不评论它**，
本程每一次提交都只带显式 pathspec，绝不吞它。⇒ 这一格**没有让本程停手**的理由：
它不污染被测码、不污染基线读数，且写面与本程无交集。

基线（逐包，禁全仓 `./...`）〔现跑〕：

```
go test -count=1 ./internal/tools/ ./internal/agent/
  ok  github.com/CarlosShao/wisp/internal/tools   19.131s
  ok  github.com/CarlosShao/wisp/internal/agent    2.551s      rc=0
```

名册口径（两数不同义，本表全程用**顶层**口径）〔现跑，锚点树 vs HEAD 树〕：

| 口径 | 命令 | 锚点 `f206e9f` | HEAD |
|---|---|---|---|
| 含子测试的 RUN 行数 | `grep -c '^=== RUN' <log>` | 239→**净 238**（见下注） | **248** |
| **顶层** PASS | `grep -c '^--- PASS' <log>` | **172** | **182** |
| 顶层 FAIL / SKIP | `grep -c '^--- (FAIL\|SKIP)'` | 0 / 0 | 0 / 0 |

注：锚点那一栏是**本程自己解出来的**，不是引用被告的号。解法＝`git archive f206e9f \| tar -x -C $(mktemp -d)`（仓外），
把 `task_output_ac2_before_test.go` 单枚文件塞进去跑第 5 节那一发，所以锚点树的 RUN/PASS 各多 1 枚；
名册求差时**显式 `grep -v` 掉那枚塞进去的用例**再比，得净名册 **172 枚**。
⇒ 被告自报的 238／172／248／182 **四数全部复算相符**。

---

## 2. 本程**没**测什么（照实记，别当查过）

1. **没跑** `probes/161/r6/flip-declaration.sh`（派单硬令：它脏 tracked 的 `logs/flip-*.txt`，那是别人的现场）。
   本程因此**没有**改动那族日志——它们在本程开工前就已是 ` M`。
2. **没跑**全仓 `go test ./...`（同树躺着别人的未提交件）。门禁只 scope 到 `./internal/tools/ ./internal/agent/` ⇒ **本表不声称全仓绿**。
3. **没复测** `go test ./cmd/wisp/` 的 `0xc0000135` 那一发（被告与台账 `pending-and-issues.md:2207` 称＝缺 `sherpa-onnx-c-api.dll` 的既有环境事实）。
   本程只独立量了 **`go build ./cmd/wisp/` = BUILD_OK**〔现跑〕，即"注册那一发在编译期过"这半句；
   "零枚用例也 0xc0000135"那半句**没重跑**，本表把它记成〔第二见证〕，不当本程读数。
4. **没有**在真机（Windows GUI 进程）上跑过一次 `wisp run` 去问模型"后台那东西吐了什么"——
   第 6 节那一问是用**产码 composition**（名册零写者＋`RunAsync` 零调用点）答的，不是用一次真机对话答的。
5. **没测** `TaskRoster` 的并发形状（`sync.RWMutex` 读写、`Record` 覆盖语义）——票面 AC 里没有这一格，本程不扩测。
6. **没测** linux 那一支（被告自报容器真跑 10/10）；本程**没重跑容器**，该项记〔第二见证〕。
7. **没核**票面 AC#4／AC#5 两格（不在本派单射程内），也**没动**票 174／175 的任何账。

---

## 3. 格① AC#2 —— **先判"判据本身"**：它算不算"未修码读数"的判据？

### 3.1 裁定：**不算**。那枚 `TestTaskOutputAC2BeforeLegIsUnreachable` 对 AC#2 这一问是**恒真形状**。

**决定性凭据（不是推理，是两棵树各跑一遍）**〔锚点＋现跑〕：同一枚文件、同一枚用例，
塞进 `f206e9f`（改前树）的仓外副本跑一次，在 HEAD（改后树）跑一次——

```
-- 改前树（f206e9f）--
--- PASS: TestTaskOutputAC2BeforeLegIsUnreachable
    AC#2 BEFORE verbatim: IsError=true ErrorClass="tool" Truncated=false Text="未知工具 task.output，可用工具见 list_tools"
-- 改后树（HEAD）--
--- PASS: TestTaskOutputAC2BeforeLegIsUnreachable
    AC#2 BEFORE verbatim: IsError=true ErrorClass="tool" Truncated=false Text="未知工具 task.output，可用工具见 list_tools"
```

**两行逐字相同。** 一棵"已经能把那发请求路由进真工具"的树，和一棵"什么都还没有"的树，
给出同一份读数 ⇒ 它量的不是树，是**它自己搭的那间屋**。复算命令＝`sh .scratch/wisp/probes/164/v1/ac2-anchor-reading.sh`（第 [5] 段）。

形状上的成因（读码）：它走的是 `helpers_test.go:32 fsBridgeWith`，而那枚 helper 里只
`for _, e := range BuiltinFSEntries(...)`（`helpers_test.go:45`）——**手抄**了一份"今天生产注册了什么"。
产码怎么改都进不到这间屋：`internal/tools/task.go` 注册不注册、`cmd/wisp/run.go` 注不注册 `BuiltinTaskEntries`，
都不改变 `fsBridgeWith` 的注册表。文件里第二发（`:69-79` 遍历 `Tools()` 见到 `task.` 前缀就 Fatal）**同病**——
它遍历的还是那间屋。

### 3.2 它的**证伪条件**到底是什么（老实答，分两支）

派单问"产码怎么改会让它红"。答得出，但**答出来的那一支不是 AC#2 那一问**：

| 能把它改红的动作 | 这是不是 AC#2 问的事 |
|---|---|
| `BuiltinFSEntries` 从此返回一枚 `task.*` 工具 | 不是——那是 fs 家族长错了牙 |
| `fsBridgeWith` 这枚 helper 开始也注册 task 家族 | 不是——那是测试脚手架自己的形状 |
| `bridge.go:251` 那句 `"未知工具 "+req.Name+"，可用工具见 list_tools"` 换字，或 `error_class` 不再是 `tool`（D37） | **也不是**——它量的是**拒绝的形状**，不是"这一发落不落在工具上" |
| `task.output` 变成可达（＝AC#2 真正问的那件事） | **它不红。** 这一行就是本格的判决 |

⇒ 精确表述：**它不是"永远不可能红"的死件**（上面前三行都能让它红，它对"拒绝文案＋D37 分类"是一枚活的守卫，
那两枚读数本身没被任何人钉过，留着不亏），**但它对 AC#2 的命题不可证伪**——
AC#2 要的是"未修码读数"，它给的是"我这间屋里读不到"。**本格判：恒真形状＝不合格作 AC#2 的凭据。**

### 3.3 替代形状：两形选一形，并说为什么

**选形 A＝归档锚 `f206e9f` 上的真树读数（`git grep`/`git show` ＋ 一份仓外副本），已落地成台件。**
不选"往 `probes/164/v1/**` 塞一份完整仓外副本"那一形，理由三条（写进了台件头部注释）：
① 那等于把同一棵树提交两遍进历史，共享工作树里是别人的现场风险；
② AC#2 那一问在锚点上由三件事实完全决定，三件都能就地量，不需要跑测试；
③ 真需要跑测试的那一步（第 3.1 节的两树对照），用 `mktemp -d` 建在 `/tmp` 就够了，**结论一样硬**。

形 A 的量到的读数（台件第 [1]–[4] 段，全〔锚点现跑〕）：

```
git grep -n "task\.output" f206e9f -- internal cmd      -> 0 命中（锚点上 Go 源码里没这个名字）
git ls-tree f206e9f internal/tools/ | grep -i task      -> 无 task*.go
git show f206e9f:cmd/wisp/run.go | grep Builtin.*Entries -> 341: 只有 tools.BuiltinFSEntries 这一发
git show f206e9f:internal/tools/bridge.go | grep 未知工具  -> 251: "未知工具 "+req.Name+"，可用工具见 list_tools"
```

⇒ **"改前那一发是未知工具"这句话，本程独立核过为真**——只是它的凭据是**锚点树的生产 composition**，不是那枚自屋用例。
被告的**结论对了，交的尺子不对**。

### 3.4 本格判语

**AC#2＝〔成立·带条件〕。**
成立的是**命题**（改前做不到＝本票成立；派单那句判断未被推翻，被告没看错产品），
条件是：**AC#2 记的凭据必须换成 §3.3 的锚点读数**；`task_output_ac2_before_test.go` 可以留，
但它的位置要改写成"未知工具拒绝形状＋D37 分类的守卫"，**不许**署名为"未修码读数"。
编排者若按现状勾 AC#2 的框，勾的是一枚自证文件。

---

## 4. 格② AC#3 正向那一半——"全文见某路径"逐字节读得回来吗？

**判语（只裁"正向读回"这一半）：〔成立〕，带一条本程补上的形状缺口（§4.3）。** 整枚 AC#3 的判语在 §6。
复算方式＝逐行读
`internal/tools/task_output_leg_test.go:165-239` 的**断言本体**（不是它的注释），再看 `internal/tools/task.go:226-240` 的产码。

### 4.1 那四段断言是真的，不是"提了个路径就算过"〔读码〕

| # | 断言（行号） | 真/假 | 本程核法 |
|---|---|---|---|
| 1 | `os.ReadFile(pointer)` 的字节 == 名册里的 `full`（`:201-208`） | 真 | `pointer` 是从**桩文本**里正则抠出来的（`:155` `全文见 (\S+)…`），不是测试自己手里的变量 ⇒ 模型看到的那句话就是被走的那句话 |
| 2 | 同一枚 pointer 再过一次**真 `fs.read`**（经同一个 `Bridge.Execute`），字节 == `full`（`:212-224`） | 真 | 注入面＝既有接缝（C1 工具＋真桥），无 mock 顶真桥（AGENTS §1.3 合规） |
| 3 | 宣布的 `总长 <len(full)> 字节` 必须出现，且 `总长 <len(桩自己)> 字节` 必须**不**出现（`:227-232`） | 真 | 第二发是负向断言，正是"拿桩长冒充全文长"那枚作弊形状的克星 |
| 4 | `省略 N` 与实际保留的头尾字节对得上账：`len(full)-N == len(headOf)+len(tailOf)`（`:234-238`） | 真 | `headOf/tailOf/omittedRe`（`:463-492`）都是从桩文本反解，产码常量没被测试引用 |

桩的读回**不依赖 fixture 之外的东西**：pointer 是从模型会看到的那句桩里抠出来的，
所以"提了个路径就算过"这一形在本用例里过不去（变异 1 就是把它摘掉，红了，见 §5）。

### 4.2 桩文本本体（复算用的那一行）〔现跑，HEAD 未变异状态〕

```
$ grep -n "全文见 %s" internal/tools/task.go
229:			"%s\n[…输出已落文件：省略 %d 字符，总长 %d 字节 / 约 %d token，全文见 %s…]\n%s",
```

与 `internal/agent/spill.go:141` 是**同一句**〔读码〕⇒ 没给模型造第二种方言（定案①合规）。
本格的**反向**凭据（两枚变异）在 §5，不在这里。

### 4.3 本程补的一条形状缺口（不改判语，但记账）

产码 `task.go:227` 的分支条件是 **`rec.ArtifactPath != ""`**，**没有任何 stat**〔读码〕。
⇒ 一枚"宿主填了路径、文件却不在／读不了"的记录，会走**有指针那一支**，
而"后半段不可找回"那一记响**不会响**——桩会指着空气。被告自报的 `TestTruncationShapeIsTheD15Triple`
正好就用了 `ArtifactPath: dir`（一枚目录，`:248`），这枚用例自己就坐在"指针不保证可读"的形状上。
本程判它**不在 AC#3 判据射程内**（判据问"截断带不带能找回全文的指针"，带了；填指针那手是宿主的活，
`TaskOutput` 注释 `:89-92` 也写明由宿主填），**且票 174 已立在"指针默认走不通"那一账上**，本表不重复开票、不动它。
⇒ 只登记一句：**"响"只覆盖 `ArtifactPath==""` 一种，覆盖不了"路径是假的"**。

---

## 5. 格③ AC#3 反向那一半——判定：〔成立〕

派单要求：两枚里**至少自己重跑一枚**、变异**先证落地**、还原来一发空 status。本程**两枚都重跑了**。
还原只用了派单指定的那一形 `git cat-file blob HEAD:internal/tools/task.go > internal/tools/task.go`；
**没跑** `checkout`／`reset`／`stash`。

### 5.1 变异 1——摘掉 `，全文见 %s`（含其参数 `rec.ArtifactPath`）

落地证明〔现跑〕：

```
$ grep -n "全文见\|totalTokens, tail" internal/tools/task.go
230:			head, totalBytes-len(head)-len(tail), totalBytes, totalTokens, tail)
235:			head, totalBytes-len(head)-len(tail), totalBytes, totalTokens, tail)
（"全文见" 整枚文件里已不存在 ⇒ 摘的就是那一味，不是没生效）
```

响不响两向〔现跑，**整包** `-v`，不是只挑目标用例〕：

```
$ go test -count=1 -v ./internal/tools/          # 变异树上
rc=1  RUN=165 顶层PASS=114 顶层FAIL=2 SKIP=0
--- FAIL: TestLongOutputPointerRecoversEveryByte (0.01s)
    task_output_leg_test.go:196: the stub must carry a recoverable pointer, got: "…省略 17200 字符，总长 20000 字节 / 约 5000 token…]"
--- FAIL: TestPointerPast256KiBIsNotFullyReadable (0.01s)
    task_output_leg_test.go:369: a >256 KiB artifact must still be pointed at, got: "…省略 304400 字符，总长 307200 字节 / 约 76800 token…]"
--- PASS: TestTruncationShapeIsTheD15Triple
--- PASS: TestNoCopyFileAnnouncesLoss
```

⇒ 被告自报"读回用例红、顺带 256KiB 那枚也红、D15 形状那枚仍绿"**三句全部复算相符**。
"未修码上响"这一问的答案＝**响**：把指针那一味摘掉（＝回到"只截不指"的那一形），承接者立刻红。

### 5.2 变异 2——`refTaskSpillHeadTokens` 500 → 400

```
$ grep -n "refTaskSpillHeadTokens = " internal/tools/task.go
52:	refTaskSpillHeadTokens = 400                          ← 落地证明
$ go test -count=1 -v ./internal/tools/
rc=1  RUN=165 顶层PASS=115 顶层FAIL=1 SKIP=0
--- FAIL: TestTruncationShapeIsTheD15Triple
    task_output_leg_test.go:264: head = 1600 bytes, want 2000 (D15 head 500 tokens x 4)
    task_output_leg_test.go:273: head/tail are not the leading and trailing windows of the full text
```

⇒ 被告"全库只它红"**复算相符**（顶层 FAIL 恰 1 枚，`TestLongOutputPointerRecoversEveryByte` 仍绿）。
期望值是字面量 `500*4`／`200*4`（`:263,266`）而不是引用产码常量——被告 §6.3 自陈"一开始引了常量、抓不住、后改成字面量"，
这一枚变异就是那句话的凭据：**它现在抓得住**。

### 5.3 还原后那行空 status〔现跑，派单指定必须贴〕

```
$ git cat-file blob HEAD:internal/tools/task.go > internal/tools/task.go
$ grep -n "refTaskSpillHeadTokens = \|全文见 %s" internal/tools/task.go
52:	refTaskSpillHeadTokens = 500
229:			"%s\n[…输出已落文件：省略 %d 字符，总长 %d 字节 / 约 %d token，全文见 %s…]\n%s",
$ git status --porcelain -- internal/tools/
（空）
```

### 5.4 为什么判〔成立〕

三枚反向判据**各管一段、没糊成一团**：变异 1 红两枚（读回＋256KiB canary）而形状枚绿，
变异 2 红一枚（形状）而读回枚绿，`TestNoCopyFileAnnouncesLoss` 在两枚变异上都保持绿——
它断的是"没落文件时必须响 `后半段不可找回`、且不许出现 `全文见`"，与前两味不重叠。
⇒ 票面 AC#3 那句"不许只截不指"现在是**有克星**的一条规矩，不是散文。
本格**没往未修码上打**的只有 §4.3 那枚"路径是假的"形状——它不在 AC#3 的判据字面里，本程登记为缺口、不判红。

## 6. 格④ 能力类 AC 必问生产调用者：**真机上今天读得到吗？**

**判语：AC#3＝〔成立·带条件〕；"读一个后台任务的输出"这件事在真机上今天做不到。**

本程现量（复算命令＝`sh .scratch/wisp/probes/164/v1/ac2-anchor-reading.sh` 第 [6] 段）：

```
# 名册写者：TaskRoster.Record 在生产码里的调用点（排 *_test.go）
$ grep -rn "\.Record(" --include=*.go internal cmd | grep -v _test.go
  -> 0 枚写者
# RunAsync 的非测试命中（4 行，其中 0 枚是调用点）
internal/agent/loop.go:320  // 文档注释
internal/agent/loop.go:321  func (l *Loop) RunAsync(...)   ← 定义本体，不是调用
internal/tools/bridge.go:663  // 一句注释
cmd/wisp/run.go:357           // 一句注释（本程 AC#3 自己写的）
```

- `task.output` **确实注册在生产里**：`cmd/wisp/run.go:361-367`〔读码＋现跑 build〕，
  `go build ./cmd/wisp/` = **BUILD_OK**〔现跑〕。
- 但那枚名册**永远零条记录**：没有任何生产代码往 `rt.tasks` 里 `Record` ⇒ 每一次真机调用
  都落在 `task.go:204-208` 的 `查不到这个任务` 那一支。**响亮地做不到，仍然是做不到。**
- `cmd/wisp/run.go` 的注释（`:354-360`）自己就把这一句写明了（"WHO FILLS IT IS NOT THIS TICKET"）⇒ 被告**没有藏这一问**，
  它自报的"生产里那句话今天还是不成立"与本程读数一致。

为什么这一格还是〔成立〕而不是〔不成立〕：票面 AC#3 的原句是**"输出太长那一格：截断策略必须带可找回的指针"**，
是一枚**策略／形状**判据，不是一枚端到端能力判据；策略已落地且**被两枚可复算的变异钉住**（§4.2）。
"没有后台写者"是 AC#4 与票 163 的地界，票面定案⑤也已把 AC#4 判成"今天连被测对象都没有"。

为什么**不许记成〔成立〕而不带条件**（派单 §1 末bullet 点名的前例＝`Q-49`：安全性只因"线没接"才成立，后被独立验收打成完整绕过）：
本程如果只报"判据已过"，读者会得到"后台任务输出可读"这一枚**产品级印象**，而真实状态是
**在册＋有实现＋生产无人喂**——正是票 164 AC#1 那笔硬账的**新形状**（从"零实现"前进到"零写者"，`Q-58` 那两枚前例的同族）。
⇒ 条件三条，勾框前必须一起写进票面：
① AC#3 的〔成立〕**只覆盖截断与指针形状**，不覆盖"真机能读到"；
② `task.output` 目前对真模型**不可达**（名册零写者），下一格（后台起跑口）没做之前**不许**把它算进产品能力；
③ 指针能不能被模型**自己**走通，取决于票 174（artifacts 目录默认不在 `[fs] allowed_dirs`）——那枚不归本票，本程没测。
