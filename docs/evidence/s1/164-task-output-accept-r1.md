# 票 164 · AC#2＋AC#3 独立验收表 v1（非实现者，只裁不改）

- 派单＝`.scratch/wisp/dispatches/2026-09-27-195x-accept-164-r1-v1-non-impl-ac2-ac3.md`（提交 `e4d8e047`）
- 票面＝`.scratch/wisp/issues/164-background-jobs-cannot-be-read-fix-the-roster-then-add-the-output-leg.md`
- 被告陈述＝`docs/evidence/s1/164-task-output-impl-r1-ac2-ac3.md`（224 行）——**本表每一个数都是自己重跑的，没有一个从它那里引**
- 本件性质＝裁决位；**只 commit、不 push**；**AC 框一枚不勾**（勾归编排者）
- 台件＝`.scratch/wisp/probes/164/v1/ac2-anchor-reading.sh` ＋ 它本次的 `.log`（**仓外副本只落在 `/tmp`，仓库目录内没建任何 worktree/checkout**）
- 分档：〔现跑〕＝本程同树跑过｜〔锚点〕＝在 `f206e9f` 解出的仓外副本上跑过｜〔读码〕＝读了产码/测试本体、没跑

---

## 0. 票面射程（判据的一部分，不是背景）——AC 两格＋五条"编排者定案"连行号

票面行号取自 `.scratch/wisp/issues/164-…-the-output-leg.md`（HEAD 版，本程 `sed -n` 现量）。

- **`:25` AC#2 未修码读数**：今天"起一个后台东西、再读它的输出"这一发**能不能做到**？量它响不响
  （做不到＝本票成立；做得到＝编排者判断错，直接报回）。
- **`:26` AC#3 输出太长那一格**：截断策略必须**带可找回的指针**（"全文在哪、怎么续读"），
  **不许只截不指**；照 `PLAN.md:431` 既有的"超 4000 token 落 artifacts、上下文留头尾＋路径"那条做，**不新造**。

**09-27 18:5x 五条编排者定案（`:31-36`），本表逐条裁定它落没落地：**

| 定案 | 行号 | 一句话内容 | 本程裁 |
|---|---|---|---|
| ① 续读形状＝**不新造游标** | `:32` | `task.output` 只给"头 500＋尾 200＋总长＋可续读路径"那一形，不发明第二套分页；`fs.read` 无偏移那一半归 `Q-59`，批不批都有归宿 | **落地**：参数只有 `{task_id}`（`task.go:166-168`），无 offset/cursor；形状照 `spill.go:141` 同句 |
| ② 持久性＝**v1 不做跨重启名册** | `:33` | 按 taskID 查不到必须**响亮**返回"查不到这个任务"，**绝对不许空返回** | **落地**：`task.go:204-208` 返回 `IsError`＋点名 id＋"这是『没有这条记录』，不是『任务没有输出』"；`TestUnknownTaskIDIsLoudNotEmpty` 钉住 |
| ③ 能力声明＝**不碰 C3** | `:34` | `PLAN.md:2564` 第 4 列是 `—`，"要第 12 枚能力"这一支不触发 | **落地**：`TaskOutputDecl()` 的 `Capabilities`/`Needs` 皆 `nil`（`task.go:253-262`） |
| ④ 落点硬约束＝**入口不许做成 `internal/agent` 上收 `taskID` 的导出方法** | `:35` | 那枚形状正被 `gate-clauses.sh` 的 **G3 腿**盯着；入口放 tools 侧；裸 `go func(` 禁；不许在 `risk.PathResolver` 之外用 `filepath.Clean\|Abs` 拼 artifacts 路径 | **落地**（但**理由里有半枚伪约束**，见 §7.2）：`internal/agent` 零字节；G3 pattern 命中 0 行 |
| ⑤ AC#4〔今天无法判定〕登记 | `:36` | 生产零后台写者 ⇒ **不许硬造一发永不响的判据当装饰**；实现那一格时先造对象、再答"未修码响不响" | **本程核：被告没有为 AC#4 造装饰判据**（新增 10 枚用例全在 AC#2/AC#3 射程，名册差集可查，见 §8） |

⚠ 定案④里那句"不许在 `risk.PathResolver` 之外用 `filepath.Clean|Abs`"——本程现量
`internal/tools/task.go` 的 import 块（`:3-13`）**既无 `path/filepath` 也无 `os`** ⇒ 那一脚从形状上不存在；
`d22scan.sh` rc=0 独立见证（§8）。

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

---

## 7. 三枚点名的攻击

### 7.1 攻击①：同源拷贝 `takeHeadTokens`/`takeTailTokens` vs `takeTokens`/`takeTokensLast` —— **语义今天没分叉；归属没写；不可接受的是"没归属"这一半**

逐字节比（本程做法：把 `internal/agent/spill.go:259-287` 两枚函数**只做改名**
（`takeTokens→takeHeadTokens`、`takeTokensLast→takeTailTokens`）与**只把 `utf8Start` 换成 `utf8.RuneStart`**，
再与 `internal/tools/task.go` 的同两枚 `diff`）〔现跑〕：

```
$ diff /tmp/a.txt /tmp/b.txt
18,19c18,19
< 	n := budget * 4          |  > 	n := len(s) - budget*4
< 	if n >= len(s) {         |  > 	if n <= 0 {
22,24c22,23
< 	start := len(s) - n      |  > 	for n < len(s) && !utf8.RuneStart(s[n]) {
< 	for start < len(s) …     |  > 		n++
< 		start++                |  > 	return s[n:]
< 	}
< 	return s[start:]
diff-rc=1        ← 差异**全部落在 takeTailTokens 一枚里**；takeHeadTokens 归一化后逐字节相同
```

**语义裁：两枚函数等价，rune 边界处理同源。** 逐条：
- 谓词同源：`internal/agent/prompt.go:301` 写的是 `func utf8Start(b byte) bool { return b&0xC0 != 0x80 }`，
  这正是标准库 `utf8.RuneStart` 的定义体〔读码〕⇒ 换名不换行为，**不是**"看着像同源其实不同"那一形。
- 头：`budget<=0→""`、`n=budget*4`、`n>=len(s)→s`、回退到 rune 起始位、`s[:n]` —— **逐字节相同**。
- 尾：被告把"剩余长度"直接写进 `n`（`n := len(s)-budget*4`），agent 那枚写成"预算"再倒减（`start := len(s)-n`）。
  守卫条件互为改写：`budget*4 >= len(s)` ⟺ `len(s)-budget*4 <= 0`；起点同为 `len(s)-budget*4`；
  前扫到 rune 起始位的循环同体 ⇒ **数值上逐点相同，只是变量承载的东西换了一个人**。
- 唯一真实的语义**差异在调用侧**，不在拷贝里：agent 那族作用在已经 `capped` 的串上并顺带记 `KeptHead/KeptTail`
  （`spill.go:136-140`），`task.output` 作用在名册全文上、只把数字写进桩句。本程判这**不是分叉**，
  是同一把刀切两种东西。

**两枚拷贝各归谁、以后谁改谁？—— 票面、`AGENTS.md`、两份证据件里都没有这一句。**〔读码〕
`task.go:275-279` 只标了**来路**（"local copies of … `spill.go takeTokens/takeTokensLast`"），
没标**方向**。这正是票 141 返工过的那枚形状（同源拷贝逐枚没归属）。
⇒ 本程给的处置口径（**不由本程落**，本程不动 `internal/**`）：
名册里 `takeTokens`/`takeTokensLast` 是权威（它服务 D15 的 4000/500/200 那套预算，`PLAN.md:431` 认的是它），
`task.go` 那两枚是**追随者**；要么把两枚收进一处（`internal/agent` 导出自由函数，或落一枚 `internal/textcut` 之类的小包，
两边都引），要么在 `task.go` 头上写死"改 `spill.go` 必须同步改这里，反之亦然"＋补一枚**逐枚对照的等值测试**
（同输入两族输出必须相同）。本程倾向后者之外的那一枚：**收进一处**，因为 §7.2 刚把"为了不响 G3 所以必须拷贝"这唯一的技术理由拆掉。

### 7.2 攻击②：那半枚伪约束 —— **不可接受（比拷贝本身更重）**

现读的 G3 腿本体〔现跑，`.scratch/wisp/probes/154/gate-clauses.sh`〕：

```
:363  want G3 quiet
:364  want_n 0
:365  run "G3 Q-56 那一支落地：Loop 上出现收 taskID 的导出方法" \
:366      '^func \(l \*Loop\) [A-Z][A-Za-z0-9]*\(.*taskID' '' internal/agent
```

（派单 §2.2 说 pattern 在 `:366-368` —— 实际是 `:363-366`，`want` 两行在它上面；本程以量到的为准。）

拿三枚候选形状**喂进那枚 pattern**（仓外 `/tmp/g3probe.txt`，不改主树）〔现跑〕：

```
func (l *Loop) TakeTokens(s string, budget int) string {     -> 不响
func TakeTokens(s string, budget int) string {               -> 不响
func (l *Loop) ReadTaskOutput(taskID string) string {        -> 响（第 3 行，命中 1 行）
```

⇒ 结论三条，一条比一条硬：
1. **导出自由函数 `agent.TakeTokens(s, budget)` 撞不到 G3**（pattern 要求 `^func (l *Loop) ` 这个接收者）。
2. 连"**导出成 `Loop` 的方法**"这一支也撞不到——pattern 还要求同一行出现字面量 `taskID`，
   而这两枚纯字符串切片函数**签名里没有、也不该有 `taskID`**。⇒ **那枚仪器对本程这一对拷贝，一个形状都不响。**
3. 所以"为了不响 G3"不能当"必须拷贝"的理由。被告注释里**每一枚事实都真**（G3 存在、`want_n 0`、
   射程是 `Loop` 收 `taskID` 那一形、定案④确实要求入口在 tools 侧）——
   但它被放在**"deliberately local copies … rather than new exported methods on agent.Loop"** 这句的**理由位**上，
   而那枚理由能否掉的形状**从来就不是唯一的备选**（真备选是自由函数/共享小包，G3 看不见）。

**裁：不可接受。** 这按派单 §2.2 的说法就是"**给一个设计选择编了一个仪器理由**"，
本程同意那一句的严重性判据：它会误导下一程——下一程读到"因为 G3 所以我们不能导出"会以为
**G3 是一堵覆盖'任何 agent 侧新出口'的墙**，从而继续在 tools 侧造第二份、第三份拷贝。
票 141 那起返工烧掉的正是"逐枚没归属的同源拷贝"，而这次拷贝还附带一枚**站不住的仪器理由**，比那次更值得纠。
**该删的是那半句仪器引用，留下定案④那半句真约束**（"入口在 tools 侧"是 owner 裁的，不是尺子裁的）。
⚠ 本程不动 `internal/**`，这一条是**给编排者的改文案账**，不是本程的产出。

### 7.3 攻击③：C25 污染源标记 —— **复算为真，归票 175，不算本票 AC#3 的账**

派单点名的两处读数，本程**复算＋纠行号**〔现跑／读码〕：

| 派单说法 | 本程量到的 | 真伪 |
|---|---|---|
| `bridge.go:551-553`：`mark` 在 `!risk.IsSensitiveSource(dec.Tool)` 时直接 return | **逐字相符**：`:551 func (b *Bridge) mark(dec Decision, res Result) {`／`:552 if b.prov == nil \|\| dec.TaskID == "" \|\| !risk.IsSensitiveSource(dec.Tool) {`／`:553 return` | 真 |
| `internal/risk/provenance.go:94-99`＝8 个名字的名册，`task.output` 不在内 | **名册确为 8 枚，但行号偏**：8 枚是**常量块 `:84-91`**（`SrcFSRead`/`SrcSearchContent`/`SrcClipboardRead`/`SrcSystemGet`/`SrcWebFetch`/`SrcDocRead`/`SrcScreenCapture`/`SrcTranscript`），`:96-98` 是同一批名字组成的 `sensitiveSourceTools` 切片；`IsSensitiveSource` 是它的线性查表 | 实质真、行号需纠 |
| 派单问："有没有看漏 `Origin`／`res.Origin` 是否产码自己填／还有别的盖戳路径吗" | 三处本程**都查了**：① `task.go` 全文**从不设 `Origin`**——生产里只有 `fs.go:148/158/170/212/218/238`、`fs_edit.go:138` 设它〔现跑 grep〕；② `PathParams` 为 `nil` ⇒ `dec.Paths` 空 ⇒ 就算过了名册那道闸，`origin` 也会是 `""`；③ **全仓另一枚 `Mark` 调用点只有 `cmd/wisp/panel_assets.go:245`**，那是 `taintSourceFlag` 的**显式声明源**清单（面板侧仿真用，`len(f.sources)==0` 就返回 nil），**不是**一条通用结果盖戳路径；④ `internal/agent/` 里 `Provenance`／`.Mark(` **零命中**〔现跑〕⇒ 环路不盖戳，桥是唯一生产盖戳点 | 无看漏；派单的读法成立 |

**本程额外量到的一枚（派单没问，但它决定票 175 怎么修）**：风险层**并没有**把非名册工具挡在外面——
`provenance.go:489-491` 写的是 `if !IsSensitiveSource(tool) { logf("… recorded fail-closed as sensitive anyway …") }`，
即 **`Mark()` 对任何工具名都会照记**，只是打一行日志。⇒ **闸门只有 `bridge.go:552` 那一枚调用侧的早退**，
不是契约层的硬墙。后果：票 175 的最小修法可以只在 tools 侧动（把 `task.output` 加进名册，或让 `mark` 的闸门看
`Decl` 上的一枚"结果可能含外部内容"属性），**不需要**改 `risk` 的语义；但那样一来 `origin` 会是空串，
桩文本里那枚 artifacts 路径才是真正有信息量的 origin —— 这一句本程**只登记、不设计**，射程外。

**是不是真破口＋最坏后果什么形状**：是真破口，**今天潜伏**——`task.output` 的**全部工作**就是把一段
"别的东西吐出来的字节"（后台命令的 stdout，其上游可能是网络/文件/子进程）搬回模型上下文，
而 D30 那族间接提示注入防护靠的就是 C25 污染标记；名字不在名册 ⇒ 这段内容**以原生工具结果的身份**进上下文，
下游看不出它不是模型自己写的。
⚠ **但本程拒绝用"线没接"给它降级**（派单 §1 末bullet 引的 `Q-49` 前例就是这一形）：
名册今天零写者（§6 现量）⇒ 这枚破口在真机上**打不通**；而**票 164 AC#4／票 163 要做的正是那根线**。
⇒ 排序风险（给编排者的一句话）：**票 175 必须与"后台写者"同片落，不能排在它后面**，
否则本票交完的那一程就是 `Q-49` 的复读。本程**没动** `internal/risk/**`，**没把它算进 AC#3 的账**。

---

## 8. 门禁（本程自己重跑，没有一个数从被告那里引）

| 门 | 命令〔现跑〕 | 本程读数 | 与被告自报 |
|---|---|---|---|
| 逐包测试 | `go test -count=1 ./internal/tools/ ./internal/agent/` | **rc=0**；`ok tools 19.131s`／`ok agent 2.551s` | 一致 |
| 同上 `-v` 四数 | `go test -count=1 -v …`＋`grep -c` | RUN **248** ／顶层 PASS **182** ／FAIL **0** ／SKIP **0** | 一致（被告 248/182/0/0） |
| 名册两向差集 | `comm -3 <(锚点净名册) <(HEAD 名册)` | **左栏 0 枚**（无被吞读数）／**右栏 10 枚**，全名见台件第 [7] 段；锚点净名册 **172** | 一致（被告 0/10）；**锚点那 172 是本程自己解出来的，不是引的** |
| D22 静态扫 | `sh scripts/d22scan.sh` | **rc=0**，`clean - no D22 ban violations`；与本票相关两行：`ban #7 internal/tools/=20 production Go files`、`ban #8 internal/=425`、`ban #8 cmd/=45` | 一致 |
| d22scan 自测 | `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | **rc=0**，`top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0` | 一致 |
| 门_clause 尺 | `sh .scratch/wisp/probes/154/gate-clauses.sh` | **rc=0**；`腿数＝14 声明与实测不符＝0`；`腿数断言：名册=14 声明=14 记账=14 缺腿=0 空头声明=0`；`聚合退码＝0`；G3/G5/G6/G7 四腿**都在场且言行相符** | **与派单 §3 那句"19:40:17 现量 rc=0／14 腿／不符 0"完全一致 ⇒ 编排者给实现程的那句"没造出未成对收尾"没被推翻** |
| 格式 | `export PATH="$PATH:$(go env GOPATH)/bin"; gofumpt --version` | **`v0.12.0 (go1.27.1)`** | 一致 |
| 格式（`-l`） | `gofumpt -l internal/tools/task.go internal/tools/task_output_leg_test.go internal/tools/task_output_ac2_before_test.go cmd/wisp/run.go` | **空**（unformatted＝0 枚） | 一致 |
| 契约轴 | `git -c core.quotePath=false diff --name-only f206e9f HEAD` ∩ `^(docs/specs/\|docs/PLAN\.md\|internal/risk/\|internal/panel/\|internal/agent/\|…/thresholds\.go\|tools/d22scan/\|allowlist\|\.github\|golden)` | **0 命中**（11 枚变更全在 `cmd/wisp/run.go`／`internal/tools/` 三枚／票面／台账／派单／证据件／新立的 174、175 票面） | 一致 |
| 生产可编译 | `go build ./cmd/wisp/` | **BUILD_OK** | 一致 |
| `flip-declaration.sh` | **没跑**（派单硬令：它会脏 tracked 的 `logs/flip-*.txt`） | 本程结束后那 6＋1 枚 `flip-*.txt` 仍是开工前就有的 ` M` 状态，**本程没改它们** | — |

⚠ **gofumpt 对"本程名下 `.go` 文件"这一格**：本程名下**没有** `.go`（台件是一枚 `.sh` ＋一份 `.log`），
所以对**被测的**四枚 `.go` 现量 `-l`＝空。`probes/164/v1/**` 里**没有 SLO 阈值、没有 golden、没有 `thresholds.go` 的抄录**（§12）。

**被告的一处口径本程复算为不符（不影响判语）**：证据件 §4 说 `RunAsync` "非测试命中只有两行注释"——
HEAD 上实际是 **4 行**（定义本体 `loop.go:321` ＋它自己的文档注释 `:320` ＋`bridge.go:663` ＋
**它自己这次新加的 `cmd/wisp/run.go:357`**），被告写那行时 `run.go:357` 还不存在 ⇒ **过期一句**，不是假话。
关键结论（**生产调用点 0 枚**）复算为**真**。

**派单一处行号本程纠**：票面定案①（`:32`）写 `fs.read` 上限"同文件 `:36`"，本程现量
`defaultMaxReadBytes = 256 * 1024` 在 **`internal/tools/fs.go:55`**；被告证据件写的是 `:55`（对）。
**而本程读到 `task_output_leg_test.go:341` 那句注释里把 `:36` 抄进了产码树** ⇒ 那行注释是**过期行号**，
建议与 §7.2 那半句一起清账（本程不动 `internal/**`）。

---

## 9. 被拒／没成功的调用（发生在取数之前还是之后）

**被权限系统拒绝：0 次。** 没成功的调用 2 次，**全部发生在取数之后**（交付段），没有一枚影响读数：

1. `git commit -q -F - -- <三枚显式路径>`（第 1 格那次）**rc=1**：
   `error: pathspec '…' did not match any file(s) known to git` ×3 ——
   原因是那三枚是本程**新建的未跟踪文件**，`git commit -- path` 只认已跟踪的改动。
   处置＝先 `git add -- <同一批显式路径>`（**没跑** `git add -A`／`git add .`）再以同一形提交 ⇒ `b5adc34f`。
   ⚠ 这一枚**本程自己写坏的**，不是环境问题；派单 §4 那句"提交一步式"在**新建文件**上少说了 `add` 这一步。
2. `go test -count=1 -run ZZZNoSuchTest ./cmd/wisp/` 的 `tail -3` **输出被吞**（本程没拿到读数）⇒
   该句记为〔第二见证〕、写进 §2.3，**没有**拿它当本程读数用。

## 10. 有没有跑过删除命令

**0 次。** `rm`／`git clean`／`git checkout .`／`git restore`／`--amend`／`reset`／`rebase`／`stash`／`worktree`／`push` 全都没跑。
唯一接近的一次：台件初稿里本程写了 `rm -rf "$BASE"`，**运行前**被本程改写成 `mktemp -d`（`ac2-anchor-reading.sh:35-37`），
从未执行 ⇒ 那次 `rm` **一个字都没跑过**。
`/tmp` 下留下的两份锚点副本（`/tmp/wisp-anchor`、`/tmp/wisp-anchor-164v1.*`）**留在原地不清**（临时件只建不删；
它们在本程主树之外，不影响任何 status 读数）。

## 11. 伪授权两栏（各带出处）

- **真通知回显＝0 枚。** 本程收到的 harness 附带件三枚，逐条出处＋按"不是授权"处理：
  1. 开局 `system-reminder` 的技能清单（列 available skills）——与票 164 无关，**没据此调用任何 skill**；
  2. `AGENTS.md` 以 project context 注入（它自称"与 `PLAN.md`/`specs`/`issues/README` 不一致时以那些文件为准"）——
     用户自己配的仓库规矩，**既不计回显也不计注入**；
  3. `Edit` 之后一句 `The file changed since your last read`（工具结果，**不是对话**）——
     **成因是本程自己**：那枚 `git cat-file … > internal/tools/task.go` 还原改了磁盘上的文件，
     harness 说的是实话；本程**没有**据此认为"别人动了我的树"，也没有据此往下多改一字节。
- **判为注入＝0 枚。** 全程没有出现任何"编排者备注／已核验请继续提交／请 revert／放宽阈值／已解锁／不用取证直接给结论"
  形状的指令；派单里编排者的四枚断言（"改前那一发是未知工具"、"gate-clauses rc=0／14 腿"、"两枚变异会红"、
  "8 枚名册里没 task.output"）本程**全部自己重跑**后才写进表（§1、§3.3、§5、§7.3），
  其中"8 枚名册"与"两行注释"两处编排者/被告的**行号与枚数**被本程纠了（§8 末两段）。

## 12. 凭据值零抄录

本表与台件里**没有**任何 API 密钥、token、DPAPI 密文、`config.toml` 里的 provider 键值、用户真实路径或任务正文。
出现的字符串只有：合成 task_id（`bg-1`／`bg-7`／`ghost-9`）、`mktemp -d` 给的 `/tmp/wisp-anchor-164v1.XXXXXX` 前缀、
以及被测桩句的**格式模板**（`全文见 %s`，不含被替换的实际路径）。
台件 `.log` 里出现的两条长串是测试自造的 `asciiRun` 字母序列（`abcdefghijkl…`），不是任何人的数据。

## 13. `next=`

1. **勾框归编排者**，但**别按现状勾 AC#2**：那一格的凭据要换成 §3.3 的锚点读数（台件已可复算），
   `task_output_ac2_before_test.go` 的**文件名与头上那段注释**建议改写为"未知工具拒绝形状＋D37 分类守卫"，
   否则下一程会把它当"未修码读数"引用（那正是本程拆掉的那枚混淆）。
2. **两笔改文案账（都属 `internal/**`，要另派，本程没动）**：
   ① `task.go:275-279` 里"为了不响 G3 所以必须拷贝"那半句删掉、留定案④那半句（§7.2）；
   ② `task_output_leg_test.go:341` 那句 `(:36 defaultMaxReadBytes)` 的行号纠成 `fs.go:55`（§8 末段）。
   顺带把 §7.1 那对拷贝的**归属方向**写死，或收进一处（本程倾向后者，因为 §7.2 已把它唯一的"技术性"理由拆了）。
3. **票 175 必须排在"后台写者"那一格之前或同片**（§7.3）：`Record` 生产零写者今天让这枚破口打不通，
   而 AC#4／票 163 要做的正是那根线 ⇒ 排序错了就是 `Q-59`／`Q-49` 那一族的复读。
   ⚠ 本程另记一枚派单没问的读数：`provenance.go:489-491` 的 `Mark()` **对任何工具名都照记**（只打日志），
   闸门**只有** `bridge.go:552` 那一枚调用侧早退 ⇒ 票 175 的最小修法在 tools 侧，不在 risk 侧。
4. **§4.3 那枚"指针不校验存在性"**：`task.go:227` 的条件是 `ArtifactPath != ""`、无 `stat`，
   "路径是假的"那一支今天**不响**。它不在 AC#3 的字面射程内、且与票 174 相邻——
   建议由编排者判它归 174（"指针走不通"）还是单立一枚，本程**没有**为它写判据（那需要动 `internal/**`）。
5. **AC#4 仍然〔今天无法判定〕**（定案⑤原样成立）：本程现量 `RunAsync` 生产调用点 **0 枚**、
   `TaskRoster.Record` 生产写者 **0 枚** ⇒ 那一格要先造出"还在写的后台尾巴"再谈判据。
