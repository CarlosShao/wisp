# 165-v1 读数件 00：起手 · 台件 · AC#1 · AC#2 · AC#5 面尺

裁决腿＝165-v1（非实现者·只读）。实现者＝`165-r1`（`internal/tools/plan_present.go`＋`plan_present_165_test.go`）
与 `165-r2`（死腿，遗产由编排者代提：`cmd/wisp/run.go` 注册那一处＋两枚 cmd/wisp 用例）。
本件里每一枚数都是**我自己现跑的**；腿报与编排者复跑记录一律标〔未验证自陈〕，只作对照不作凭据。

## 0. 起手三把尺（现跑）

- `git rev-parse HEAD`＝`71002d001af60feaf2bbd9b2a07863a6cb738a23`，分支 `dev`；时刻 `2026-10-16 10:17:11 +0800`。
- 票框尺：`grep -c '^[[:space:]]*- \[ \]' .scratch/wisp/issues/165-*.md`＝**5**；`...'- \[x\]'`＝**0**（与派单预期一致）。
- 起树锚之后的链（尺＝`git log --oneline b74b87d0..HEAD`，分母＝**93 枚提交**）：动 `internal/tools` 的只有 **3 枚**——
  `4a49480da`（165-r1 交件）、`03da4941b`（**270-r1**，task.spawn 判空）、`e1973191b`（**287-r1**，深度夹具）。
  ⚠ 与派单第 5 行那句"本波只有一枚写腿 216-r1、⛔ 动 tools"不符：`216-r1` 在这一段里**零提交**（尺＝
  `git log --all --grep='216-r1'` 只命中编排者的放腿记录），而 `internal/tools` 里真动了刀的是 270/287 两枚**更早波次**的腿。
  归属现量（尺＝`git log -1 -- <那四枚 165 件>`）：`plan_present.go`／`plan_present_165_test.go` 的最后 toucher＝
  `4a49480da`；`cmd/wisp/plan_*_165r2_test.go`／`run.go` 那一段的最后 toucher＝`17ca21ee`。⇒ 270/287 没碰过 165 的任何一枚件。
- 165 的字节面（我自己取的号，尺＝`git log --all --grep='165'`＋逐枚 `git show --stat`，⛔ 抄派单）：
  `4a49480da`（r1 交件 9 路径）＋`a754a390`（r1 票尾 bullet 1 路径）＋`17ca21ee`（r2 遗产代提 15 路径）＋`988ce413`（补 r2 漏交的证据件 1 路径）＝**4 枚、26 路径**。
  另有一枚**契约冻结笔**`824bb9c`（2026-09-27，动 `docs/PLAN.md`＋`docs/specs/SPEC-07-tools-and-plugins.md`）——那是 D34 那一行与 SPEC-07 镜像的**来源**，
  ⛔ 记进 165-r1/r2 的账（它早于本波起手锚 18 天，且不属这两腿）。

## 1. 台件（仓外导出树，⛔ 共享工作树、⛔ worktree、⛔ `-overlay`）

| 树 | 取法 | `.go` 枚数（先数后信） | 用途 |
|---|---|---|---|
| `D:\tmp\165v1\t0` | `git archive b74b87d0 \| tar -x` | **1062** | AC#1 未修码那发 |
| `D:\tmp\165v1\t1` | `git archive HEAD \| tar -x` | **1091** | 基线＋全部变异重放 |
| `D:\tmp\165v1\t2` | `git archive HEAD \| tar -x` | **1094** | AC#2 反形（Gate 加第三条回调），**留在突变态不复用** |

- `third_party/` 不被 git 跟踪 ⇒ 三枚 DLL（`onnxruntime.dll`／`sherpa-onnx-c-api.dll`／`sherpa-onnx-cxx-api.dll`）我从仓里**复制**进三棵树的 `third_party/sherpa-onnx/`，
  再 `PATH="$TREE/third_party/sherpa-onnx:$PATH"`——这正是编排者在 A998 具名过的 `0xc0000135` 那一族，本腿三棵树零发踩到。
- `GOTMPDIR`＝`D:\tmp\165v1\gotmp`（原生反斜杠形，先 `mkdir -p`）。原始日志一律落 `D:\tmp\165v1\logs\*.log`（仓外、只建不删），本件只抄逐字关键句。
- 还原对账尺（跑完全部变异之后现跑，尺＝仓库里 `git cat-file -p HEAD:<f> \| sha256sum` 对 `sha256sum <树>/<f>`）：
  `internal/tools/tool.go`＝`602d632b…` MATCH、`plan_present.go`＝`3fc87895…` MATCH、`bridge.go`＝`a0c0736c…` MATCH、`cmd/wisp/run.go`＝`7f28db4f…` MATCH。
- 门禁两把（⚠ 归属具名）：`go vet ./internal/tools`（HEAD 导出树与共享工作树两处）**rc=0**；`go vet ./cmd/wisp` 在 **HEAD 导出树 rc=0**，
  在**共享工作树 rc=1**，红句逐字＝`cmd\wisp\panel_l1_blocked_220_test.go:236:84: (*testing.common).Errorf call has arguments but no formatting directives`。
  该文件 `git status --porcelain` 读作 `?? cmd/wisp/panel_l1_blocked_220_test.go`（**未跟踪**）＝本波唯一写手 `220-r2` 在飞的件 ⇒ **⛔ 记成 165 的账**。
  `gofumpt` 本机仍不存在（`which gofumpt` 空，与 A998 那句一致）；`gofmt -l` 对五枚 165 件＝**列出 0 枚**（⛔ 拿这把 0 冒充 gofumpt 的 0）。

## 2. AC#1 未修码读数（凭据＝三把尺＋两形）

1. **代码名册尺（静态）**：`grep -rn "plan\.present" --include=*.go t0/internal t0/cmd`＝**0 命中**；
   `grep -rn "BuiltinPlanEntries\|planPresent"`＝**0 命中**。对照同一把尺在 t1：`BuiltinPlanEntries` 命中 4 枚文件
   （`plan_present.go`／`plan_present_165_test.go`／`plan_wiring_165r2_test.go`／`run.go`）⇒ 尺不是瞎的。
2. **D34 名册尺（契约面）**：`plan.present` 那一行**在冻结文本里、不在码里**——`docs/PLAN.md:2567`（＋镜像 `docs/specs/SPEC-07-tools-and-plugins.md:73`）
   两处现读都在（`824bb9c` 落的）；`internal/tools` 在 t0 里零枚 plan 形 ⇒ 「契约立了、码里没有」这一形坐实。
3. **未修码那发（跑）**：`t0`＋我自己写的一枚同名形状尺（`v1ac1_roster_test.go`，只读 fs＋task＋subagent 三个既有 builder，另加 registry 面）：
   rc=0、`--- PASS`，逐字红位没出现，Logf 逐字＝
   `V1-AC1 未修码名册 10 枚、零枚 plan/propose 形、registry.Lookup(plan.present)=false：[fs.read fs.list fs.write fs.edit fs.trash fs.move fs.delete task.output task.cancel task.spawn]`
   ⇒ 名册 10 枚、plan 形 0 枚、**连 registry 面都查不到**（这一味比 r1 那枚只数名册的尺多一层）。
4. **反形（加一味即响）**：同一把尺在 t1 上把 `BuiltinPlanEntries()` append 进去 ⇒ `V1-AC1 反形成立：名册 11 枚、plan 形 1 枚＝[plan.present]`，rc=0（用例要求它就是响一次）。
5. **生产可达那一发（未修码必红）**：把 r2 交付的那枚调用点尺（`plan_wiring_165r2_test.go`）原样盖进 `t0`，靶向 `-run Test165r2PlanWiringCallSiteCensus`：
   **rc=1**，阳性对照那行仍绿（`cmd/wisp/run.go:783 range=true register=true`＝先例在 t0 里也已接），红句逐字＝
   `接线 RED：tools.BuiltinPlanEntries 在 internal/ 与 cmd/ 的产码里调用点 0 枚，要 1 枚（cmd/wisp/run.go 的装配那一处）；命中清单＝[]（0 枚＝写了没接线那一族，plan.present 在生产装配里根本够不到）`
   ⇒ **今天（未修码）"只读地出一份方案"做不到**，票 165 前提成立；这一发同时是 AC#1 的"做不到"与 r2 那枚尺的去牙正控。

判语见证据件 §AC#1。⚠ 本格**没有**写成"应该做不到"：第 3／5 两发是真跑出来的。

## 3. AC#2 复用形状（判据＝"没新增状态"，⛔ 读成"功能做全了"）

- **面尺（`--name-only`，分母＝26 路径）**：四枚 165 提交里 `internal/risk/**`＝**0 路径**、`internal/statemachine/**`＝**0 路径**、
  `internal/agent/approval/**`＝**0 路径**、`internal/panel/**`＝**0 路径**。⇒「没动状态机文件」这一句在字节面上成立。
- **`Gate` 回调枚数（两把尺同读）**：文本尺＝`sed -n '/type Gate interface/,/^}/p' internal/tools/gate.go | tr -d '\r' | grep -cE '^[[:space:]]+Pending[A-Za-z]*\(ctx context.Context'`＝**2**
  （`PendingWindow`＝L1 支、`PendingApproval`＝L2 支，定义在 `internal/tools/gate.go:110`）；
  反射尺＝`reflect.TypeOf((*Gate)(nil)).Elem().NumMethod()`，在 `Test165PlanPresentAsksThroughTheExistingL2Card` 里判 `!= 2` 即 `t.Fatalf`（用例件 `:152`）⇒ 基线 rc=0 全绿。
- **`GrantSource` 方法面**：同一把文本尺＝**1**（`Covering`，定义在 `internal/tools/grant.go:39`）；反射尺在 `Test165PlanPresentAddsNoAllowExit`（`:314`）判 `NumMethod()!=1 || 第一枚不叫 Covering` ⇒ 基线绿。
- **本票自己有没有造状态/答案形状**：`grep -cE 'Answer[A-Z]|Mode[A-Z]|State[A-Z]' internal/tools/plan_present.go`＝**0 命中**（产码里零枚新常量引用）；
  `planPresent` 的方法名册只有 4 枚＝`Name`／`Description`／`Parameters`／`Execute`（`grep -n '^func (planPresent'`），即 `Tool` 接口那一面，零枚额外出口。
- **D43 那四行的引用点（现读，⛔ 抄 r1 的号）**：冻结表里 `#17`＝`docs/PLAN.md:3077`（`Acting` 需确认 → `Confirming`(L1)／`AwaitingApproval`(L2)）、
  `#23`＝`:3083`（队头被**原生侧**允许 → `Acting`）、`#24`＝`:3084`（拒绝／**超时 300s** → `Acting` 取消该调用）、`#25`＝`:3085`（队列空且任务完成）。
  码侧对应＝`internal/statemachine/table.go` 里 `D43: 17`（`:118`、`:123` 两条 to）、`D43: 23`（`:158`、`:163`）、`D43: 24`（`:168`、`:172` 含 `EvApprovalTimeout`）、`D43: 25`（`:176`、`:180`）＝**8 条转移条目、4 个 D43 号**，
  全部是 165 之前就存在的条目（面尺已证 165 零字节）。⇒ r1 报的"复用 #17/#23/#24/#25"这四枚号**名实相符**：L2 那一发走的正是 `Acting → AwaitingApproval → (native allow｜reject/timeout) → Acting`，
  没有第五个号、也没有新态。
- **AC#2 的反形（我自己造的世界，`t2`）**：给 `Gate` 加第三条回调 `PendingPlan(ctx, Decision) (Answer, string)`，
  并为让它编译而在导出树里补 4 枚桩（`v1m4_stubs_test.go` 里的 `GateFuncs`／`gateSpy`／`plan165Gate`／`t90Gate`／`NoGate`，＋一枚新文件 `internal/agent/approval/v1m4_stub.go`——
  ⚠ 具名：这枚桩落在 **AC#5 禁区面**里，但**只在仓外树**、⛔ 入仓、⛔ 算 165 的越界账；t1 的越界尺与 HEAD 对账都不含它）。
  结果 **rc=1**，唯一红＝`Test165PlanPresentAsksThroughTheExistingL2Card`，逐字＝
  `plan_present_165_test.go:153: Gate 有 3 个方法：本票复用现成的 PendingWindow/PendingApproval 两个回调，多出来的第三个就是一张新审批通道`
  另四枚（含 AC#1/AC#3/AC#4 那三枚）仍 PASS ⇒ 这枚"两回调"尺**有牙**、且红在它该红的那一格上。
  突变态哈希留存：`internal/tools/gate.go`＝`e31a43136ee091c152192c97eb1edb06f5cf4e04564ed61760d689cc6875e7bd`、
  `internal/tools/v1m4_stubs_test.go`＝`15aba439ad3e9203b115dfb6e59fda29006bc84487afb1ddb337eb367b35d176`、
  `internal/agent/approval/v1m4_stub.go`＝`469f0bce70c83d8d222d2ef06c46d6c180bb1d06e3fa587c612a4b587ff309bb`（`t2` 整棵留在突变态，⛔ 再用作任何基线）。
- ⚠ 具名一处"缺世界"：AC#2 没有一枚尺去判「`AwaitingApproval` 遇系统挂起」那一格——那本就不是本票的判据，且它是 `AGENTS.md §2` 与 `docs/PLAN.md:1530` 的**具名停手项**（DEFERRED，归 S7），
  ⛔ 由我裁、⛔ 由 165 落。这一格我只登记"没造过这个世界"，不写判语。

## 4. AC#5 面尺原始读数（判语与分母见证据件 §AC#5）

- 分母＝**26 枚路径**（`4a49480da` 9 ＋ `a754a390` 1 ＋ `17ca21ee` 15 ＋ `988ce413` 1），尺＝`git show --name-only --format='' <四枚>`。
- 禁区正则名册（`docs/PLAN.md`、`docs/specs/`、`internal/risk/`、`internal/panel/`、`internal/agent/approval/`、`internal/observe/thresholds.go`、`*golden*`、
  `tools/d22scan/allowlist.txt`、`frontend/`、`design/`、`go.mod`、`go.sum`、`.gitignore`、`.github/`、`*ci.yml`、`internal/perm/ticket90_persist_test.go`）
  对这 26 枚的命中＝**0 枚**（`grep -c` 回 0、退出码 1＝真零命中，stderr 空）。
- 同一把尺的**正控**：喂 `824bb9c`（D34 冻结笔）的 2 枚路径 ⇒ 命中 **2 枚** ⇒ 尺有牙。
- 26 枚路径的**实际名册**（逐枚）：`internal/tools/plan_present.go`、`internal/tools/plan_present_165_test.go`、
  `cmd/wisp/plan_assembled_165r2_test.go`、`cmd/wisp/plan_wiring_165r2_test.go`、`cmd/wisp/run.go`、`docs/evidence/s1/165-wiring.md`、
  `.scratch/wisp/issues/165-plan-first-then-act-a-brake-the-voice-entry-needs.md`、`.scratch/wisp/probes/165/00..06`（7 枚）、`.scratch/wisp/probes/165/r2/00..09,11`（12 枚）。
  ⇒ 全落在本票写面与三处交付面内。
- 区间归属（⛔ 把别腿的路径算成 165 的账）：`b74b87d0..HEAD` 这 93 枚里动过上面禁区的共 **8 枚**（尺＝`git log --oneline b74b87d0..HEAD -- <禁区>`），
  逐枚＝`7c80ac032`、`f82b909e9`、`9b820f8e8`、`c7752e536`、`73eda0294`、`f3d447346`、`272405f53`（320-r2）、`4f2967e90`——**零枚**署名 165-r1/r2，且这 8 枚全在 `4a49480da` 之后或属前端/编排者亲改面。
