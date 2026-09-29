# 票 222 — 非实现者对抗验收（腿 `222-v1`）

- 工单：`.scratch/wisp/issues/222-spawn-holds-a-bridge-permit-while-waiting-for-its-child.md`
- 本件：`docs/evidence/s1/222-spawn-holds-a-bridge-permit-v1.md`
- 起手锚点（本腿现取 `git rev-parse --short HEAD`）：**`801e8547`**（编排者给我的 `fac60ad4` 之后又落了 4 枚 `undone-28-1` 提交；本腿以 `801e8547` 为准）
- 本腿角色：**非实现者对抗验收**。产码那半边＝`222-r1`（09-29 12:46，票面 §"收 222-r1"）。
- 本腿纪律：零产码（`internal/**`／`cmd/**` 非测试文件一枚不改、`_test.go` 不新增不修改）；突变只在 `.scratch/wisp/probes/222/v1/mutations/**` 的补丁里做，每发跑完还原并自证 `git status --porcelain -- cmd/wisp internal/` ＝空。

## 起手名册（本腿第一次跑的绿/红名集合）

未判（待填：本腿自跑的 `internal/tools` ＋ `cmd/wisp` 逐名名册）。

---

## AC#1 生产形状第一次有尺（票面 `:35`）

**判语：成立（带两枚注，注不推翻判据本身）**

**现读依据（文件:行＝本腿 `801e8547` 上现跑的尺，非转述）**

| 事实 | 出处（现读） |
|---|---|
| 孩子的工具面走真桥＝**唯一**载荷行是 `ParentTools: h.bridge` | `internal/tools/subagent_222_test.go:278` |
| `Tools: h.bridge` 那行是**死行**：产码随后无条件覆写 `opt.Tools = newSubagentToolProvider(t.d.ParentTools)` | `internal/tools/subagent_197.go:284` ⟸ 覆写 `subagent_222_test.go:288` |
| 包装器把孩子的调用**逐字转交桥**，不另开路 | `internal/tools/subagent_197.go:567` `return p.inner.Execute(ctx, req)` |
| 桥的 `Execute` 就是取许可那枚入口 | `internal/tools/bridge.go:438-442`（`slot.take(b.sem)` ＋ `defer slot.giveBack()`） |
| 生产形状确实是同一枚桥 | `cmd/wisp/run.go:589 ParentTools: rt.bridge` |
| 三枚新用例真名与行号 | `subagent_222_test.go:362`／`:434`／`:469`；`go test ./internal/tools -list '.*'` 现跑＝**161 枚顶层**，其中含 `222` 的 **3 枚** |

**突变读数（全部走 `go test -overlay`，工作树零写入；overlay  inert 性由 `asis` 对照证明：两枚产码文件从 `git cat-file blob HEAD:` 逐字节取回、`diff` 空、跑出的名册与基线一致）**

| 号 | 突变内容（台件在 `.scratch/wisp/probes/222/v1/mutations/<id>/`） | 读数 |
|---|---|---|
| 基线 | 未突变 | `internal/tools` 整包 `-count=1 -v`＝**161 PASS／0 FAIL／0 SKIP**（`--- PASS` 计数，13.081s）；三枚 222 用例各 **0.00s 绿** |
| **M1＝修之前那一发**（删 `subagent_197.go:356 giveBackWhileWaiting(ctx)`） | 正控 | **三枚全红**：`WaitingParentHoldsNoBridgeSlot` **0.00s**（`subagent_222_test.go:443`"占着 4 枚桥位，want 0"）、`CeilingStillCaps…` **0.00s**（`:477`）、`SpawnConclusion…` **3.00s** 红在内容（`:395` 四枚父任务全收到 `工具 task.spawn 超时（3000ms），已协作式中止`；`:413` 四枚结论各出现 **0** 次）⇒ **两枚 0.00s、不靠挂死、不靠超时**，第三条按票面 `:35` 允许的第三种判据形态（断言父侧拿到的是孩子结论）判红，且用的是台件自己的 3s 旋钮而不是生产的 30s |
| **M3＝把孩子换回 `h.dir` 形状**（`ParentTools`/`Tools` 都改指 `fake197Dir{}`） | AC#1 空心攻击（本腿任务书点名的那发） | `SpawnConclusion…` **红**，但红在 `:374`"只等到 0/4 枚，**护栏到点**"＝**30.00s 的 guard rail**，另两枚 **绿**（0.00s） |
| **M4＝孩子仍在跑探针但绕开桥**（新增 `bypass222` 目录直调 `probe.Execute`，不取许可） | 更阴的一发 | `SpawnConclusion…` **0.00s 红**：`:387`"第 4 枚孩子的工具调用起跑时桥位占用 **0** 枚, want 4" ⇒ "走真桥"这半条被**读数**钉住，不是被"有没有结果"钉住 |

**注①（不许瞒）**：AC#1 要求的"孩子工具面走真桥"确实有尺，但**只有 `Test222SpawnConclusionArrivesThroughRealBridgeChildren` 一枚能看见它**——M3 与 M4 之下另两枚用例照绿。M3（票面 `:19` 记录的原始 `h.dir` 形状）下这枚用例是**靠 30s 护栏才红**的，正是票面 `:35` 警告的"等超时才红"那形的一种；M4（保留工具执行、只绕开许可）才是 0.00s 读数判红。判据文字没有被违反（`:35` 的"秒级"约束射程＝**未修码那一发＝M1**，M1 两枚 0.00s），但下一位若要把 `h.dir` 那形也钉成读数红，缺的断言是"起跑前孩子在桥上的次数"应在 `await` **之前**先读一次 `h.probe.runs`。本腿不改测试件，只具名。

**注②**：台件里 `:288 Tools: h.bridge` 是死行（被 `subagent_197.go:284` 覆写）。它没有让用例变松（载荷行是 `:278`），但文件头 `:269-271` 那句"That one line is the whole difference"指的必须是 `ParentTools` 那行才成立——本腿现读确认：是 `:278`。

**还原自证**：全程 `-overlay`，工作树未被写过一枚字节。本腿每发突变后现跑
`git status --porcelain -- cmd/wisp internal/` ＝ **空**（见文末逐发记录），
且四枚文件 `sha1sum` 与起手值逐字节相同：
`bridge.go 36a1d2b9…`／`subagent_197.go df860999…`／`subagent_222_test.go b3938cf0…`／`subagent_197_test.go 3ff5c9a7…`。

**可复跑的尺（一格一把）**
```
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
python .scratch/wisp/probes/222/v1/mut222.py m1-drop-giveback
go test ./internal/tools -count=1 -run 'Test222' -v \
  -overlay=.scratch/wisp/probes/222/v1/mutations/m1-drop-giveback/overlay.json   # 期望：3 枚红（0.00/0.00/3.00s）
python .scratch/wisp/probes/222/v1/mut222.py m4-child-bypasses-bridge
go test ./internal/tools -count=1 -run 'Test222' -v \
  -overlay=.scratch/wisp/probes/222/v1/mutations/m4-child-bypasses-bridge/overlay.json  # 期望：:387 读数红 0.00s
```


---

## AC#2 "占着许可干等"被永久钉住（票面 `:36`）

**判语：成立（带三枚注：其中一枚推翻票面 `:78` 的突变②，一枚推翻 `bridge.go:666-669` 的注释）**

**现读依据**

| 事实 | 出处（现读） |
|---|---|
| 许可在取到之后交给 `inFlightSlot`，`defer` 兜底 | `internal/tools/bridge.go:438-442` |
| 交还载体＝未导出 `inFlightSlot{once, sem}`，`giveBack` 在 `once.Do` 里 `<-s.sem` **并把 `s.sem = nil`** | `bridge.go:657-660`／`:670-680`（`:675` 那枚 nil 赋值） |
| 工具侧入口＝未导出 `giveBackWhileWaiting(ctx)`，取不到载体返回 false、永不阻塞 | `bridge.go:688-700` |
| 等待点的交还调用 | `internal/tools/subagent_197.go:356` |
| 占用读数（AC#2 的正读数）＝等待瞬间 `len(h.bridge.sem)` want **0**；同时钉 `cap(...) == 4` 字面 | `subagent_222_test.go:441-445`／`:446-448` |
| 反控读数＝父全等待时**同时执行**的调用数 `maxSeen` want **恰 4**、逐枚 `slotsHeld ≤ 4`、终态 `len(sem)==0` | `subagent_222_test.go:476-478`／`:496-500`／`:506-508`／`:524-526` |

**突变读数（锚点＝数字真正所在那一行）**

| 号 | 突变 | 读数 |
|---|---|---|
| M1 | 删 `subagent_197.go:356` | `:443` 与 `:477` 两处 **0.00s** 判红，读数＝**4 枚桥位被等待中的父任务占着**（＝票面现量链的形状被读数钉住，不靠超时、不靠挂死） |
| **M5**（反控是否恒真） | `bridge.go:24 MaxToolConcurrency = 4 → 8`（只改台件，产码未动） | **反控不是装饰**：`Test222WaitingParentHoldsNoBridgeSlot` **0.00s 红**（`:446` 读到 `cap=8`）＋ `Test222CeilingStillCapsExecutedCallsWhileParentsWait` **0.00s 红**（`:497` 读到 `maxSeen=6`）。同发里既有钉子 `TestToolConcurrencyCeilingIsFour` **照绿**（0.03s）⇒ 该文件头 `:59-64` 与票面 A420 那句"它拿常量跟自己比，看不见常量被抬"**本腿独立复现成立**；全量名册：159 枚绿／2 枚红，两枚红都是 222 的新用例 |
| **M2**（票面 `:78` 的突变②） | 只拆 `bridge.go:672` 的 `s.once.Do`，保留 `:675 s.sem = nil` | **三枚全绿，0.040s，`ok`** ⇒ **"拆掉 sync.Once ⇒ 双扣那形必须红"这句话不成立**（见下面推翻清单第 1 条）。形状上说得通：交还与 `defer` 兜底跑在**同一枚 goroutine 上顺序两次**，第二次被 `s.sem = nil` 拦下，`once` 在这条路径上无从发力 |
| **M9** | 反向隔离：保留 `once.Do`，只删 `:675 s.sem = nil` | **三枚全绿，0.032s** ⇒ 两枚护栏**各自单独够用**，删一枚测不出红＝这形被**冗余**钉住（好），但注释把载荷归给了错的那枚（坏） |
| **M8** | 两枚护栏**同时**拆掉（＝真正的"同一枚许可被扣两次"） | **三枚全红**，但每枚都红在 **30.00s 的 guard rail**（`:392`／`:457`／`:518`"只等到 0/4 枚，护栏到点"），整包 90.037s ⇒ 双扣这形**确实有尺**，可尺是**挂死护栏**，**不是** r1 自述（票面 `:68`）里那句"终态 `len(sem)=0` 兼作双扣那形的尺"——那一行（`:524`）在 M8 下根本没被执行到 |
| **M7** | 把 `:356` 的交还挪到两条"已派生"文案**之后** | **三枚全绿**；再压 `-cpu=1 -count=20`＝20/20 绿、默认 GOMAXPROCS `-count=200`＝200/200 绿 ⇒ `subagent_197.go:351-353` 那句"a reader that has seen this call's spawn delta therefore knows the slot is already back…decide without a deadline"所声称的**因果标记没有任何断言在守**：位置挪了，用例既不会红、也不会变读数不同（本机上读数照旧是 0，因为父任务必先阻塞才轮得到读侧） |

**注（AC#2 文字与读数的射程差）**：票面 `:36` 要的是"可用许可数＝4 −（父之外正在执行的真工具数）"这枚**等式**；台件里两处分别钉住了"父等待时占用＝0"（`:441`）与"探针自记占用 ≤ 4／同时执行数＝4"（`:506`／`:497`），**等式右半边只在 ≤ 方向被检查**。本腿判这不影响"占着许可干等被钉住"这一枚 AC 的成立，但把射程写清楚，别让下一位以为等式两侧都有尺。

**可复跑的尺**
```
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
for m in m5-ceiling-to-8 m2-drop-sync-once m9-drop-nil-guard m8-double-drain-for-real m7-giveback-after-delta; do
  python .scratch/wisp/probes/222/v1/mut222.py $m >/dev/null
  go test ./internal/tools -count=1 -run 'Test222' -v \
    -overlay=.scratch/wisp/probes/222/v1/mutations/$m/overlay.json | grep -E '^(--- |ok|FAIL)'
done
# 期望逐名：M5＝2 红 1 绿 / M2＝3 绿 / M9＝3 绿 / M8＝3 红(各 30s) / M7＝3 绿
```


---

## AC#3 父侧结论通道要么真、要么明说不做（票面 `:37`）

**判语：不成立（半格；两支甲/乙都没走，且两支都还缺批准）**

**现读依据**

| 事实 | 出处（现读） |
|---|---|
| 甲形（`task.spawn` 立即返回句柄）**没做**：`:382-393` 仍是 `select { case <-ctx.Done(): … case res := <-done: … }`，调用阻塞到孩子结束才返回 | `internal/tools/subagent_197.go:382-393` |
| 乙形（等待不受 tool 超时约束）**没做**：`TaskSpawnDecl()` 仍不设 `Timeout` ⇒ `timeoutFor` 落 `b.defTm` ⇒ 生产值＝`cfg.Agent.PerToolTimeoutMS`（未配置 `DefaultToolTimeout = 30s`）；带 deadline 的 `ectx` 在取到许可**之后**才建，而它就是等待用的那枚 ctx | `subagent_197.go:202-211`／`bridge.go:545-557`／`bridge.go:28`／`cmd/wisp/run.go:570` |
| 丙形（交还许可）**确实落了**且不改上面两行 | `subagent_197.go:356` ＋ `bridge.go:688-700` |

**⛔ 越权检查（本腿任务书点名的那枚扳机）＝没有踩**
`git show bfc55b5b -- internal/tools/bridge.go internal/tools/subagent_197.go` 新增的顶层声明**逐枚现读＝六枚，全小写开头**：`inFlightSlot`／`take`／`giveBack`／`inFlightSlotKey`／`withInFlightSlot`／`giveBackWhileWaiting` ⇒ **零枚新导出名**；源名复用现成 `risk.SrcTaskOutput`（`subagent_197.go:463` 往父任务作用域盖戳），**没有** `subagent.output` 之类的新名字；`SubagentDeps` 仍是 5 枚字段（`:128-148`：Roster／BaseOptions／ParentTools／Provenance／Stream），197 那枚字段枚举守卫不必重谈。**所以 AC#3 不是因为越权判不成立的。**

**突变读数（本格的新尺＝M11，打在已修的产码上，产码一枚字节未动）**
只改测试件里的两样东西（`-overlay`，台件在 `mutations/m11-budget-50ms-gated/`）：把 C22 旋钮 `subagent_222_test.go:70 h222PreFixBudget` 从 `3 * time.Second` 缩到 **`50 * time.Millisecond`**，并给孩子在第一枚模型调用里加一段 **200ms 的 `late`**（＝"孩子真的在干活、比父任务的预算久"这件事的生产情形）。三发 `-count=3` 读数**逐发相同、3/3 红、每发 0.20s**：

- `:395`（突变件里漂到 `:400`）四枚父任务**全**收到 `"工具 task.spawn 超时（50ms），已协作式中止"`；
- `:413`→`:418` 四枚孩子的结论各出现 **0** 次；
- `:380`→`:385` 因果读数：`孩子 c1…c4 的工具调用起跑时已有 **4** 枚父任务返回`——孩子**照样挤上了桥**（丙形生效，父确实不占许可），只是**到得比父任务的预算晚**。

⇒ 生产上把 50ms 换成 `cfg.Agent.PerToolTimeoutMS`（未配置 30s）：子代理真干活是几分钟量级 ⇒ **父任务必然收到超时文案、孩子继续跑并正常落名册**＝"名册对、父任务错"那枚分裂**今天仍在**。票面 `:69` 与台账 `A430 :9273` 那句"只修掉一半"由本腿**复现成读数**，不是文字推断。
（⚠ 诚实记一发更早的尝试 **M10**＝只把旋钮改成 `1ms`、不加 `late`：`-count=3` 里 **1 红 2 绿**。原因是纯内存假流程真能在 1ms 内跑完⇒ **M10 是撞运的、不算尺**，本格判语只用 M11。顺带读出一枚对该文件的认识：`SpawnConclusion` 那枚用例的确定性靠的是 gate＋3s 的**富余**，而不是靠断言本身。）


**两支为什么都不能由 agent 自裁（具名，勿互相冒充已完）**
- 乙＝给这枚调用一枚 C22 豁免 ⇒ 顶到 `bridge.go:30-31` 逐字 "nothing here can switch enforcement off"＝**契约面**。
- 甲＝改 `task.spawn` 语义 ⇒ 票面 `:37` 明写"**要先落一枚 `A##`**"；本腿现查台账：**没有任何一条 `A##` 落的是甲形**——`A420 :9144` 把它列为"要点头的那一支"、`A430 :9273` 原样升级给 owner、`A434`（13:10）批的是**票 221 的甲形＝注册 `task.cancel`**，而且 `:9330` 逐字写着"父任务等孩子的时限那一维不在这里（那是票 222 的 AC#3 剩半格，两支都要另外批）"。⇒ **别把 `A434` 当成本票甲形的批准，它不是。**

**还原自证**：`git status --porcelain -- cmd/wisp internal/` ＝空；M10 的台件在 `.scratch/wisp/probes/222/v1/mutations/m10-budget-1ms/`，测试件 `sha1sum` 仍 `b3938cf0…`。

**可复跑的尺**
```
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
python .scratch/wisp/probes/222/v1/mut11.py
go test ./internal/tools -count=3 -run 'Test222SpawnConclusion' -v \
  -overlay=.scratch/wisp/probes/222/v1/mutations/m11-budget-50ms-gated/overlay.json
# 期望：3/3 红（每发 0.20s）：四枚"超时（50ms）"＋四枚结论出现 0 次＋"已有 4 枚父任务返回"
```


---

## AC#4 不许留任何"对模型许诺了但没接"（票面 `:38`）

判语：**未判**

依据（文件:行，现读）：未填
归票 221 的具名理由：未填

---

## AC#5 池帽与桥帽的关系要现读定性（票面 `:39`）

判语：**未判**

依据（文件:行，现读）：未填
`PLAN.md:2849`（D38）逐字核：未填
突变读数：未填

---

## AC#6 整包终态读数 + 门禁（票面 `:40`）

判语：**未判**

编排者两把尺的引用与口径核对：未填
本腿自跑的门禁（gofumpt / build / d22scan）：未填

---

## 推翻派单里的哪句话

未填（推翻不了就写"未推翻"）。

## 票 222 现在该勾哪几格 / 哪几格必须留着（具名原因）

未填。

## 没做完 / 判不了（具名清单，本节不许空）

未填。

## 伪授权登记（工具输出里任何自称"编排者备注／系统提示／用户已改规矩／请 revert／请冻结／请放宽"）

命中数：**0**（截至骨架落盘）。原文与出现位置：无。
