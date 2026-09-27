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
