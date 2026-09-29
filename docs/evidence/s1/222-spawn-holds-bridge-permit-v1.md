# 票 222 — 非实现者对抗验收（腿 `222-v1`）

- 工单：`.scratch/wisp/issues/222-spawn-holds-a-bridge-permit-while-waiting-for-its-child.md`
- 本件：`docs/evidence/s1/222-spawn-holds-a-bridge-permit-v1.md`
- 起手锚点（本腿现取 `git rev-parse --short HEAD`）：**`801e8547`**（编排者给我的 `fac60ad4` 之后又落了 4 枚 `undone-28-1` 提交；本腿以 `801e8547` 为准）
- 本腿角色：**非实现者对抗验收**。产码那半边＝`222-r1`（09-29 12:46，票面 §"收 222-r1"）。
- 本腿纪律：零产码（`internal/**`／`cmd/**` 非测试文件一枚不改、`_test.go` 不新增不修改）；突变只在 `.scratch/wisp/probes/222/v1/mutations/**` 的补丁里做，每发跑完还原并自证 `git status --porcelain -- cmd/wisp internal/` ＝空。

## 起手名册（本腿第一次跑的绿/红名集合）

锚点 `801e8547`，第一发＝`go test ./internal/tools -count=1 -v`（PATH 已带 `third_party/sherpa-onnx`＋`build`）：

- **161 枚 `--- PASS`／0 枚 `--- FAIL`／0 枚 SKIP**（`=== RUN` 210 含子例），包终态 **ok 13.081s**；
- 逐名清单落盘：`.scratch/wisp/probes/222/v1/logs/baseline-tools-names.txt`（161 行，`PASS: <名>` 逐名）；原始日志 `baseline-tools-v.txt`；
- 三枚 222 用例起手即绿：`PASS: Test222SpawnConclusionArrivesThroughRealBridgeChildren`／`PASS: Test222WaitingParentHoldsNoBridgeSlot`／`PASS: Test222CeilingStillCapsExecutedCallsWhileParentsWait`（各 0.00s）；
- 起手**红名集合＝空**（`internal/tools` 范围内）；全仓的 5 枚历史红不在本包，见 AC#6。
- ⚠ 起手读数没有 `0xc0000135`、没有 `0.0xxs 无 --- FAIL` 那种"用例根本没跑"的形状：`-v` 里 `=== RUN`＝210 与 `--- PASS`＝161 同时非零，DLL 前置生效。


---

## AC#1 生产形状第一次有尺（票面 `:35`）

**判语：成立（带两枚注，注不推翻判据本身）**

**现读依据（文件:行＝本腿 `801e8547` 上现跑的尺，非转述）**

| 事实 | 出处（现读） |
|---|---|
| 孩子的工具面走真桥＝**唯一**载荷行是 `ParentTools: h.bridge` | `internal/tools/subagent_222_test.go:278` |
| `Tools: h.bridge` 那行是**死行**：产码随后无条件覆写 `opt.Tools = newSubagentToolProvider(t.d.ParentTools)` | `internal/tools/subagent_197.go:284` ⟸ 覆写 `subagent_222_test.go:288` |
| 包装器把孩子的调用**逐字转交桥**，不另开路 | `internal/tools/subagent_197.go:567` `return p.inner.Execute(ctx, req)` |
| 桥的取许可点在 `run`，而 `Execute` 逐字转交它 | `bridge.go:245` `func (b *Bridge) Execute` → `:317` `return b.run(ctx, req, entry, dec)` → `:438-442`（`slot.take(b.sem)` ＋ `defer slot.giveBack()`） |
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

**判语：不成立（本格未闭；两句话今天还在，仍在对模型许诺。归票 221＝结案路径，不是翻勾理由）**

**现读依据**

| 事实 | 出处（现读） |
|---|---|
| 许诺一：说明书正文写着"可以用 `task.cancel` 单独停它" | `internal/tools/subagent_197.go:196`（票面 `:189` 已漂 7 行，票面更正块 `:61` 说的 `:196` ＝ 本腿现读一致） |
| 许诺二：父侧放弃那条文本写着"可以单独停它：它的流键是 %s" | `internal/tools/subagent_197.go:388-389`（票面 `:362-363` → 现读 `:388-389`） |
| 而 `task.cancel` **没注册**：`BuiltinTaskEntries` 只返回 `task.output` 一枚 | `internal/tools/task.go:595-597` |
| 子代理族注册名册只有 `task.spawn` | `internal/tools/subagent_197.go:218`（`BuiltinSubagentEntries`） |
| `task.cancel` 在 Go 侧的唯一身份是**注释里的 DEFERRED 标记** | `internal/tools/task.go:23`（"DEFERRED with five fields, PLAN.md §7 :1531"）／`:576`／`:593` |
| 两枚许诺行**不是 `222-r1` 写的、也没被它碰过** | `git blame -L194,197`／`-L386,390` 逐枚＝`7ea14ce3e`（09-28 19:54，票 197 那批）；`git show bfc55b5b -- internal/tools/subagent_197.go` 的变更行里两枚都不在 |

**还在对谁许诺**：`Description()` 的读者是模型（它进工具目录随请求一起发出去）。票 221 文件 `:18` 逐字「模型会照说明书办事：它以为能停掉自己派出去的孩子，于是**不会去找别的退路**」——本腿核过这句的前提今天仍成立：全仓 Go 侧 `task.cancel` 的**非注释命中只有那一行说明书**（尺见下）。

**可复跑的尺（三把 grep，零产码；本腿逐把现跑，读数写在后面）**
```
# ① 说明书里的 task.* 词根          ⇒ 读数：task.cancel（唯一命中）
sed -n '/func (subagentSpawn) Description/,/^}/p' internal/tools/subagent_197.go \
  | grep -o 'task\.[a-z]*' | sort -u
# ② 注册名册（Tool.Name() 的返回值）⇒ 读数：task.output ＋ task.spawn（两枚）
grep -rh 'func (.*) Name() string { return "task\.' internal/tools/*.go \
  | grep -o 'task\.[a-z]*' | sort -u
# ③ 全仓 Go 侧 task.cancel 的非注释命中面 ⇒ 读数：只剩 subagent_197.go:196 那一行说明书
grep -rn "task\.cancel" --include=*.go internal/ cmd/ | grep -v "_test.go" \
  | grep -vE "^\S+:[0-9]+:\s*//"
```
①−② 的差集非空＝这把尺今天该响。票 221 的 AC#1 要的正是把它做成常驻能力尺，**改前必须响**——本腿这三发就是它的改前响度样本。
（⚠ 自我更正一次：本格子件首版把①写成了 `grep -n '"task\.'`，那把打的是 `Name()` 的返回值而不是说明书句子，现跑读数里 :196 并不出现；上面这版才是能复跑的那把。改判据文字＝零，只是把尺写对。）


**登记一枚本格射程外、但同族的谎（只具名，不改判语）**：`subagent_197.go:195` 说明书第一句"派生一枚子代理去独立完成一个子任务，**等它跑完并把结论带回本任务**"——本腿 **M11** 已把它量成假（孩子比父任务 per-tool 预算久时，父任务收到的是桥的超时文案、结论 0 次到达，见 AC#3 那格）。AC#4 的字面只点名两处"可以单独停它"，所以本格不因此加判；但**谁结案 AC#4，就得同一批面对这一句**，否则票 221 交完之后说明书仍有一句为假。⛔ 两枚票同撞 `internal/tools/subagent_197.go` ⇒ 串行、同批只一枚碰（票面 `:38`／`:46`）。

**为什么必须留着（具名）**：票 221 文件 `:4` 逐字「⇒ **本票现在是"可派写腿"状态**，但排程在后面：`226-v1` 交完 → 票 223 → `222-v1` → 本票（甲形＋乙形三处同一发）」⇒ 甲形（注册 `task.cancel`）已有 owner 批准记录＝台账 `A434`（13:10，撤销口令「撤 221 甲」），但**注册代码今天不在树上**（上面那把尺就是证据）。AC#4 只能在票 221 的写腿落地后随它结案。

**本腿零动作自证**：本格**一枚文件都没改**，也**没做任何突变**（突变只会让两句许诺更难看的红，量不出新东西）。`git status --porcelain -- cmd/wisp internal/` ＝空。


---

## AC#5 池帽与桥帽的关系要现读定性（票面 `:39`）

**判语：成立带注（论证成立；注＝守这枚论证的**盘上理由**已经过期，本格不许靠它复述）**

**盘上形状（现读，本腿未动任何帽值）**

| 事实 | 出处 |
|---|---|
| 桥的执行许可仍是字面 4 | `internal/tools/bridge.go:24 const MaxToolConcurrency = 4` |
| 池仍是字面 4（不是桥帽的别名，故意写成数字） | `internal/tools/subagent_197.go:85 MaxConcurrentSubagents = 4`（理由写在 `:60-84`） |
| 冻结文字逐字仍在原行 | `docs/PLAN.md:2849`「**工具并发上限 = 4**（防 LLM 一次发 50 个 tool call 打爆机器；也与 D45 的批量聚合配合 —— 4 路并发足以让批量场景快起来）」（`:2850` 是它的续行；D38 标题在 `:2818`，背压小标题在 `:2845`） |
| 同值三处旁证本腿逐枚现读 | `docs/specs/SPEC-01-architecture.md:134`、`SPEC-01:26`（「工具执行（工具并发 ≤4）」）、`SPEC-05-agent-core.md:70`（「与工具并发 ≤4」） |
| 契约面一字未动 | `git diff --stat beaeaeba..HEAD -- docs/PLAN.md docs/specs/ internal/risk/thresholds.go` ＝ **空**（本腿现跑） |

**本腿对两问的裁定**

1. **"修完之后池＝4 是否仍必要？"→ 机制上不再必要。** 丙形落地后，父等待不占许可（AC#2 的 M1 读数：占用 0），池与桥在机制上**解耦**：桥帽约束"同一瞬间在跑的工具调用数"，池帽约束"同时存在的孩子数"，两枚数字相等**不再是被桥逼出来的**。盘上仍写 4，是一枚**被保留的设计选择**。
2. **"池若写 8 违不违反 `PLAN.md:2849`？"→ 不违反。** 那句话的射程逐字是"工具并发上限"＋给出的理由是"一次发 50 个打爆机器"＝**正在执行的工具调用数**；轮转不抬高同时执行数。本腿的反控读数支持这句：`maxSeen`（同时执行数）＝4（`subagent_222_test.go:497`），且 M5 证明该断言不是恒真——把桥帽抬到 8，同发它就读到 6 并判红。**⇒ "8 枚孩子在 4 枚许可上轮转、单轮并发仍是 4"这条论证成立，不需要为它先改契约。**（票面 `:30`／台账 `A420:9141` 的推断，本腿复现。）

**突变读数（M6＝只在台件里把池写 8，桥仍 4；整包 `-v` 现跑）**：红 4 枚，逐名与逐句原文如下——

- `Test197SubagentPoolNeverExceedsBridgeCeiling`（`subagent_197_test.go:403`）**0.00s 红**，`:405` 原文：「池 8 大于桥的 D38d 天花板 4：多出来的 4 枚会被桥排成队，名册却说它们在跑；诚实数 = 天花板（**票 211 甲**），要更多并发**得先动契约**（票 211 乙），不是改这枚常量」；
- `Test197SubagentPoolCapsAtBridgeCeiling`（`:441`）红，`:444`；`Test197FullPoolRefusesNextSpawnWithReadableReason`（`:526`）红，`:529`（这两枚 r1 未点名 ⇒ **它的"两枚钉子"数少了两枚**，实测三枚 197 钉子＋一枚 222 用例）；
- `Test222SpawnConclusionArrivesThroughRealBridgeChildren` **3.00s 红**：8 枚父任务全收「工具 task.spawn 超时（3000ms），已协作式中止」、8 枚结论各 0 次、`:380` 记「孩子 c5/c6 起跑时已有 **8** 枚父任务返回」＝**池 8 时后 4 枚孩子只能在 gate 打开后才上桥**。

**注①（这一条是 AC#5 的真正交付物，别抄 r1 的理由）**：守池=4 的那枚钉子，**它写在盘上的理由已经被票 222 证伪**。`internal/tools/subagent_197_test.go:398-402` 逐字：「one in-flight spawn **holds one bridge slot for its child's whole life** (bridge.run keeps the semaphore held across entry.Tool.Execute)」——丙形之后这句**不再为真**（`subagent_197.go:356` 在等待前交还）。同理它那句"要更多并发得先动契约（票 211 乙）"正是票 222 §"为什么 211 的前提塌了"（票面 `:26-30`）作废的推断。⇒ **AC#5 的结论只能站在新理据上：抬池不是契约问题，是"名册把排队的孩子印成「在跑」"这枚诚实性问题＋三枚钉子的重写问题 ⇒ 归 owner 一句话。** `222-r1` 自己在 §⑧ 登记过"两段理由过期但断言照旧绿、本腿未碰那枚文件"，本腿复核为**真**且把它升级为具名缺陷：文件＝`internal/tools/subagent_197_test.go:398-402`，与票 221 同撞 `internal/tools`，**串行**。
**注②（不许把 M6 的红当成"池 8 有害"的证据）**：M6 里那枚 222 用例红，成因是台件的 `gate222` 要求 8 枚孩子**同时**inside 桥，而桥只有 4 枚许可——这是**测试构造出来的阻塞**，不是生产形状。它只证明"这枚用例的分母跟着 `MaxConcurrentSubagents` 走"。票面 `:39` 禁的是"靠测试绿了过"，本腿对称地补一句：**也不许靠测试红了就判池 8 有害**。
**注③（本腿边界）**：本格只裁论证与盘上形状，⛔ 没动 `PLAN.md`／`docs/specs/**` 一字、没把任何帽值改成 8（M5/M6 只存在于 `.scratch/wisp/probes/222/v1/mutations/*/` 的台件里，工作树 `git status --porcelain -- cmd/wisp internal/` ＝空）。

**可复跑的尺**
```
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
python .scratch/wisp/probes/222/v1/mut222.py m6-pool-to-8
go test ./internal/tools -count=1 -v \
  -overlay=.scratch/wisp/probes/222/v1/mutations/m6-pool-to-8/overlay.json | grep -E '^--- FAIL'
# 期望 4 枚红：Test197SubagentPoolNeverExceedsBridgeCeiling / …CapsAtBridgeCeiling /
#              Test197FullPoolRefusesNextSpawnWithReadableReason / Test222SpawnConclusion…(3.00s)
grep -n "工具并发上限 = 4" docs/PLAN.md            # 2849，仍在
git diff --stat beaeaeba..HEAD -- docs/PLAN.md docs/specs/   # 空
```


---

## AC#6 整包终态读数 + 门禁（票面 `:40`）

**判语：成立（带两枚口径注）**

**先立可比性**：`git diff --name-only fac60ad4..HEAD -- '*.go'` ＝ **空** ⇒ 编排者那两把尺跑的 Go 树与本腿锚点 `801e8547` **逐字节同码**（那 4 枚提交只动 `.scratch/**` 与台账）；`git diff --stat beaeaeba..HEAD -- internal/tools/bridge.go internal/tools/subagent_197.go internal/tools/subagent_222_test.go` ＝ **空** ⇒ 222 的三枚交付件自 `222-r1` 起没人再动过。

**引用编排者的两把尺＋口径核对（件在 `.scratch/wisp/probes/222/orch/`，本腿逐行现读）**

| 他的读数 | 本腿核到的口径 |
|---|---|
| 整包 `-count=1`：26 枚包＝23 ok／3 FAIL，`cmd/wisp` ok 135.067s | `clean-nov.txt`（首行 `start=17:53:11`）现数：**23 行 `ok`＋3 行 `FAIL <pkg>`＋5 行 `[no test files]`＝31 行包级读数**；"26 枚"＝去掉 5 枚无测试包的口径 ✓；`cmd/wisp 135.067s` ✓ 逐字 |
| 逐名红册＝恰 5 枚、零新增 | ⚠ **口径要说清**：`three-v.txt` 的 `-v` 只覆盖了**那 3 枚 FAIL 包**（ball／panel／risk，405 枚 `=== RUN`，包级 1 ok＋2 FAIL），**不是全量 `-v`**。全量非 `-v` 那发（`clean-nov.txt`／`full-v.txt`）里有 **6 枚 `--- FAIL` 名**＝5 枚＋`TestResolvePerCallBudget`。⇒ "恰 5 枚"是**扣除 risk 那枚争用型红之后**的集合，不是原始全量 `-v` 名册（票 223-v2 那种）。本格按"5 枚＋一枚已归因"记账，不按"原始 5 枚"记账 |
| risk 那枚红＝两发重叠的争用型假红 | **归因方向本腿支持、机制句无法从盘上证实**：`full-v.txt` `start=17:50:18`、`clean-nov.txt` `start=17:53:11`，两发文件里**都没有结束时刻**，相差 173s——是否重叠算不出来（A432 说的是撞在 `226-w1` 收尾的门上，那不发不在本目录）。本腿**安静复量 3 发**：`-run TestResolvePerCallBudget` **PASS 1.322s**、整包 `-count=3` **ok 8.174s**、整包 `-count=1 -v` **ok 2.339s**＝**3/3 绿**；`internal/risk/thresholds.go` 与预算常量 `git diff beaeaeba..HEAD` ＝ **空**（一字节未动）。⇒ 判**不计入新增红**成立，且不放宽任何断言 |
| 5 枚归属 | 本腿**未碰、也不许碰**：`internal/ball` 1 枚＋`internal/panel` 4 枚，逐字原因全指另一队地界（C21 令牌表与面板契约字段对表，含一枚缺失的样式源文件；那枚文件的内容本腿零读零转述） |

**本腿自己现跑的门禁（票面 `:40` 点名要自己跑的这几枚）**

| 门 | 命令 | 读数 |
|---|---|---|
| gofumpt（⚠ 裸名不在 PATH，会 127） | `"$(go env GOPATH)/bin/gofumpt.exe" -l` 三枚文件 | **空输出，exit 0**；`--version` ＝ **v0.12.0**（与票面 `:40` 同版） |
| 编译 | `go build ./...` | **exit 0** |
| vet | `go vet ./internal/tools` | **exit 0** |
| D22 扫描 | `bash scripts/d22scan.sh` | **clean — no D22 ban violations**；射程计数现读：bans#1-5 internal/=219、cmd/=29、ban#6 frontend/=85、**ban#7 internal/tools/=22**、ban#8 design/=39＋frontend/=85＋internal/=460＋cmd/=63 |
| 用例名册规模 | `go test ./internal/tools -list '.*'` | **161 枚顶层**，含 `222` 的 **3 枚**（与 r1 自述 158＋3 相符） |
| 本腿自跑的包终态 | `go test ./internal/tools -count=1 -v` | **161 `--- PASS`／0 FAIL／0 SKIP，13.081s** |
| 本腿自跑的 cmd/wisp | `go test ./cmd/wisp -count=1` | **ok 110.491s**（编排者 135.067s／r1 两发 77.622s／83.788s ⇒ 差异＝机器负载，**耗时不作判据**，只记三枚样本） |

**注①（时间敏感量法的复量）**：本腿所有会受负载影响的读数都带了复量或跨腿样本——并发帽三读互证（`bridge.go:24` 常量字面＝4／M5 台件抬到 8 后 `cap(bridge.sem)` 读数 8 判红／`maxSeen` 读数 6 判红），`-overlay` 的 asis 对照与基线同绿（0.042s vs 0.00s 级用例）证明突变架本身不改变结果。**未做**：全量 `-v` 终态名册（那是编排者的尺，他已声明由他自己复跑翻勾，本腿不越俎）。
**注②**：`gofumpt -l` 只对**三枚交付件**跑过（票 222 的写面就这三枚）；整仓 `-l` 没跑，那不是本票的判据。

**可复跑的尺**
```
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
"$(go env GOPATH)/bin/gofumpt.exe" -l internal/tools/bridge.go internal/tools/subagent_197.go internal/tools/subagent_222_test.go   # 期望：空
bash scripts/d22scan.sh | tail -2                                                                                                       # 期望：clean
go test ./internal/tools -count=1 -v | grep -cE '^--- PASS'                                                                             # 期望：161，且 grep -c '^--- FAIL' ＝ 0
go test ./internal/risk -count=1                                                                                        # 期望：ok（争用假红的复量样本）
```


---

## 推翻派单／票面里的哪句话（逐枚带尺；推翻不了的就写"未推翻"）

1. **推翻（票面 `:78` 给下一位的待办②，＝`222-r1` 的自述）**："把 `giveBack()` 的 `sync.Once` 拆掉 ⇒ '同一枚许可被扣两次'那形必须红"。**实测不成立**：M2（只拆 `bridge.go:672` 的 `once.Do`、留 `:675` 的 `s.sem = nil`）三枚用例 **0.040s 全绿**；M9（只删 `s.sem = nil`、留 once）也 **0.032s 全绿**。真正的双扣形要把**两枚护栏同时拆掉**（M8）才红，而且它红在 **30s 挂死护栏**上，不是 r1 自述（票面 `:68`）那句"终态 `len(sem)=0` 兼作双扣那形的尺"——M8 下那一行（`:524`）根本没被执行到。**这条要按"追加更正"进票面，不改写原文。**
2. **推翻（我派单里的口径句，不是结论句）**："逐名红册（`-v`，按用例名集合）＝恰 5 枚"。**盘上那发 `-v`（`orch/three-v.txt`）只覆盖了 3 枚 FAIL 包**（ball／panel／risk，405 枚 `=== RUN`），不是全量 `-v`；全量非 `-v` 那两发（`full-v.txt`／`clean-nov.txt`）逐名是 **6 枚**＝5 枚历史红＋`TestResolvePerCallBudget`。"零新增红"这个**结论**本腿认（risk 那枚 3 发安静复量全绿），但**"恰 5 枚"不能当成全量 `-v` 名册来引用**。
3. **推翻（我派单里的一句状态描述）**："票 221 还没派"。**现读**：票 221 的只读普查腿 `221-c1` **已交**（台账 `A433`，28,379 字节，`dc4f3870`），甲形（注册 `task.cancel`）**已获 owner 批准**（`A434`，13:10，撤销口令「撤 221 甲」）；**没派的是写腿**，且票 221 文件 `:4` 逐字把写腿排在本腿（`222-v1`）**之后**。⇒ AC#4 的"未闭"结论不变，但结案路径比我派单说的更靠前、更具体。
4. **推翻一半（票面 `:71`，r1 的 AC#5 自述）**："池抬到 8 会**立刻红** `Test197SubagentPoolNeverExceedsBridgeCeiling` 与 `Test197SubagentPoolCapsAtBridgeCeiling`"。**两枚点名的钉子确实红**（0.00s，`subagent_197_test.go:405`／`:444`），但**实测红 4 枚**：还多 `Test197FullPoolRefusesNextSpawnWithReadableReason`（`:529`）与本票自己的 `Test222SpawnConclusion…`（3.00s，见 AC#5 注②对它的**限定**）。⇒ 它没说错方向，说少了数目。
5. **未推翻**（逐枚列出，免得被当成"顺手都推翻了"）：
   - 待办①"删掉 `subagent_197.go:356` ⇒ 三枚用例必须红"＝**成立**（M1：三枚红，两枚 0.00s、一枚 3.00s 判在正文）。
   - 待办③"`MaxToolConcurrency`／`MaxConcurrentSubagents`／`PLAN.md:2849` 未被为了并发而说破"＝**成立**（M5 证明反控能看见帽被抬；常量与契约 diff 全净）。
   - 票面 `:52` 的三枚提交文件清单与消息逐枚相符＝**成立**（`git show --stat 92a4b5b7／bfc55b5b／beaeaeba` 现读逐枚相符，每枚都带 pathspec 落点）。
   - 票面更正块 `:60` 那句"父任务实际收到的是桥的超时文案而不是'不等了'"＝**成立**（M1/M11 读数原文：`工具 task.spawn 超时（Nms），已协作式中止`）。
   - AC#3 的"不许新造导出名"这一枚**没被踩**：新增六枚顶层声明逐枚小写开头，源名复用 `risk.SrcTaskOutput`，`SubagentDeps` 仍 5 枚字段。

## 票面 `:24` 留给验收腿的那枚归因，本腿的读数（"永挂"一支）

票面 `:24` 逐字："⚠ **'永挂'这一支我不下结论**（腿报'无限挂住'，但父侧有 `ectx` 超时兜着 ⇒ 真形状是**周期性松绑＋父侧一律报错**，不是死锁）。**归因未做，留给验收腿现跑。**"
**本腿裁定：编排者那句推断成立。** M1（未修码）读数：4 枚父任务在 **3.00s** 全部拿到错误正文并**释放许可**，孩子在它们返回**之后**照样挤上桥跑完（`:380` 逐枚记"已有 4 枚父任务返回"），终态池位回收由 `:417` 读到"没回收：2"＝父走了、孩子还在跑。⇒ **不是死锁，是"父侧一律报错＋孩子继续跑完并落名册"的分裂**。**真会挂死的是另一枚形状**：许可被扣两次（M8），那形三枚用例各 30.00s 撞护栏、整包 90.037s，且盘上用"终态 `len(sem)=0`＋nil 护栏＋once"三处把它按住（前两者的单独贡献见 AC#2 的 M2/M9/M8 读数）。

## 票 222 现在该勾哪几格 / 哪几格必须留着（具名原因）

| AC | 本腿判语 | 该不该勾 | 留着的原因（具名） |
|---|---|---|---|
| AC#1 | 成立（两注） | **可勾** | 注①②是缺陷登记，不改变"生产形状第一次有尺＋修前红修后绿"这条判据被满足 |
| AC#2 | 成立（三注） | **可勾**，但**必须同时把票面 `:78` 待办②就地更正**（追加、不抹原文） | 反控非恒真（M5 读数在案），"永久钉住"成立；被推翻的是 r1 描述突变的那句**话**，不是那枚钉子 |
| AC#3 | **不成立（半格）** | ⛔ **必须留着** | 甲形（`task.spawn` 立即返回句柄）**至今没有 `A##`**（台账 `:9144`／`:9273`／`:9330` 逐字；`A434` 批的是票 221 的 `task.cancel`，不含本格）；乙形撞 `bridge.go:30-31` 的 C22 契约面；M11 三发 3/3 红证明通道仍未真 ⇒ **要 owner 一句话，agent 不自裁** |
| AC#4 | **不成立（未闭）** | ⛔ **必须留着** | `subagent_197.go:196`／`:388-389` 两句"可以单独停它"仍在，`task.cancel` 仍不在注册名册（`task.go:595-597`）；结案＝票 221 的写腿，排程在本腿之后，且同撞 `subagent_197.go` ⇒ 串行、同批只一枚碰 |
| AC#5 | 成立带注 | **可勾**（它只要求"现读定性"） | 勾的是"定性已交付"，**不是**"池可以写 8"。抬不抬池＝owner 独立问题；注①那枚过期理由（`subagent_197_test.go:398-402`）要另开一格收，别让它继续当论据 |
| AC#6 | 成立 | **可勾** | 编排者本人的那发整包已在（`A432`＋`orch/*`），本腿核了口径、补了 risk 三发复量与全套门禁（gofumpt v0.12.0/build/vet/d22scan/-list/cmd-wisp 110.491s） |

⇒ **本腿建议：勾 AC#1／AC#2／AC#5／AC#6 四格，留 AC#3／AC#4 两格。** 票面 Status 那句"这枚票今天只做成半格的那一半要说给 owner"（`:3`）**照旧成立且更硬**：本腿把"剩的那半"量成了读数（M11），不是文字推断。

## 没做完 / 判不了（具名清单）

1. **没做（有意不做，越权面）**：抬池帽到 8／改 `MaxToolConcurrency`／改 per-tool 超时／给 `task.spawn` 开 C22 豁免——全归 owner（AC#3/AC#5）。
2. **没做（尺的射程不够，具名）**：AC#1 注①——M3（`h.dir` 原形）下只有 1/3 用例能看见真桥、且靠 30s 护栏才红。补法本腿已写清（在 `await` 之前先读一次 `h.probe.runs`），**但本腿不改 `_test.go`**。
3. **没做（不是本票的尺）**：全量 `-v` 终态名册（编排者的尺，他已声明由他自己复跑翻勾）；本腿只跑了 `internal/tools`（全量 `-v`）＋`cmd/wisp`＋`internal/risk` 三包。
4. **没做（外部对标）**：票面 `:31` 那些"别家分设两枚帽"的路径（Step-Code／DSH／minimax／openchamber）本腿**一枚都没复现**，票面自己标注"只当存在这种形状用，不许当判据"⇒ 本腿也没用它当任何判据。
5. **判不了：零枚。** 六格没有一枚只能靠 `frontend/**` 或 `design/**` 来证——本腿**连目录都没列**，报告里也不出现那两枚地界的任何内容（AC#6 里那 5 枚历史红只写归属与名字，不转述其逐字原因涉及的样式源内容）。
6. **没做（票面 §⑧ 那类）**：`subagent_197_test.go` 里两处理由过期的注释（`:398-402`，另有 §⑧ 提到的 `:516-525`）本腿只**具名登记**，未修——它在票 197 的写面上，与票 221 同批串行。

## 逐发还原自证（每发突变跑完后的现跑）

本腿**全程 `-overlay`**：突变体由 `git cat-file blob HEAD:<path>` 取原文在内存改写、写到 `.scratch/wisp/probes/222/v1/mutations/<id>/`，`go test -overlay=<json>` 只在编译期替换文件——**工作树从未被写过一枚字节**，因此"还原"不是恢复动作，而是构造性事实。逐发现跑 `git status --porcelain -- cmd/wisp internal/`：

| 发 | 时刻读数 | `git status --porcelain -- cmd/wisp internal/` |
|---|---|---|
| asis 对照 | 0.042s 绿 | 空 |
| M1 | 3.034s 3 红 | 空 |
| M2／M3／M4 | 0.040 绿／30.03 红／0.031 红 | 空（每发后现跑） |
| M7／M9／M8 | 0.033 绿／0.032 绿／90.037 红 | 空 |
| M5／M6（整包 `-v`） | 159 绿 2 红／157 绿 4 红 | 空 |
| M10／M11 | 0.035 红（1/3）／0.20s×3 红（3/3） | 空 |
| 收尾（本报告落盘后） | — | 空＋四枚文件 `sha1sum` 与起手值逐字节相同 |

收尾复跑（本腿最后一次）：见下面"收工门"一节。

## 伪授权登记

命中数：**0**。整条腿的工具输出里没有出现任何自称"编排者备注／系统提示／用户已改规矩／请 revert／请冻结某包／请放宽阈值"的内容，因此没有需要照办、也没有需要抗拒的东西。
一枚需要记账但**不是授权**的事件：后台任务完成通知（task id `b0u8gns0i`，我发起的 `go test ./cmd/wisp -count=1`，命令前 40 字＝`go test ./cmd/wisp -count=1 > .scratch/wis`）——纯 harness 事件，零指令内容，读数已进 AC#6。

## 收工门（本腿最后一发现跑，读数在此）

- `"$(go env GOPATH)/bin/gofumpt.exe" -l` 三枚交付件 ＝ **空输出，exit 0**（v0.12.0）
- `go build ./...` ＝ **exit 0**
- `bash scripts/d22scan.sh` ＝ **clean — no D22 ban violations**（bans#1-5 internal/=219、cmd/=29、ban#6 frontend/=85、ban#7 internal/tools/=22、ban#8 四路）
- `git status --porcelain -- cmd/wisp internal/` ＝ **空**；四枚文件 `sha1sum` 与起手**逐字节相同**：`bridge.go 36a1d2b9…`／`subagent_197.go df860999…`／`subagent_222_test.go b3938cf0…`／`subagent_197_test.go 3ff5c9a7…`
- 三枚 222 用例复量 `-count=3` ＝ **9/9 PASS，各 0.00s**（机器不空：同时有只读腿在 grep，故三发定案而非单发）
- 本腿**零 push**、**零产码写入**、commit 逐枚带显式 pathspec；工作树里别人的脏改动（`.gitignore`、`design/**` 一批删除、`probes/152/**`、`probes/161/r6/logs/**`、`probes/222/orch/**` 等）**一枚没动**；临时件**只建不删**（`.scratch/wisp/probes/222/v1/**`：`mut222.py`＋`mut11.py`＋9 枚 mutations 目录＋18 份原始日志）。


