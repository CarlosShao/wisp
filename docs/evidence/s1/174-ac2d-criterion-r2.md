# 174-AC#2d 判据第一枚（写码腿 `174-r2`·纯测试面）

> 派单：`.scratch/wisp/dispatches/2026-09-28-151x-impl-174-r2-ac2d-first-criterion.md`
> 工单：`.scratch/wisp/issues/174-...md`（**本程只碰 AC#2d 这一格**）
> 上游普查：`docs/evidence/s1/174-wiring-cost-census-c2.md` ＋ 台账 `A372`
> 骨架落盘时刻 `2026-09-28 14:58 +0800`；本表终值读数时刻见 §5（终态取在最后一枚 commit 之后）
> ⛔ **AC 框一枚没勾**（勾由编排者翻）；`internal/tools/task.go`／`bridge.go`／`cmd/wisp/run.go`／`internal/risk/**` **零字节**。

## ① 起手锚 ＋ 写面闸门

- 起手 HEAD：**`40aee084`**（`Mon Sep 28 14:56:28 2026 +0800` · `docs(evidence): 33-a1 land the twelve-section skeleton for the panel inbound-hop survey (no product code touched)`）
- 派单声明的编排者锚点是 `ee0ef8ec` 之后的实际 HEAD ⇒ 共享树里别人在推进，**按实际 HEAD 做**；**偏差登记**：起手锚 `40aee084`（非 `ee0ef8ec` 直接后继），不改判、不改做。
- 起手写面自证：`git status --porcelain -- internal/ cmd/` ＝**空**。
- 骨架落盘在第 **6** 枚工具调用（派单 §4 要求 ≤5）：**偏差登记**——同批并发了三枚必读件（工单 174／`internal/tools/task.go:215-364`／`task_output_pointer_notice_test.go` 全文）＋本表 Write，多一枚；不影响任何读数。
- 本程写面只有两枚新文件：
  - `internal/tools/task_output_canonicalize_fail_174_test.go`（判据本体）
  - `.scratch/wisp/probes/174/r2/**`（变异载具＋overlay，**全部新建、零删除**）
- **本程没接线**：`cmd/wisp/run.go` 那枚 `TaskDeps{}` 补 `Paths` 一个字节的动作**没做**（派单 §1 禁，且同批有写腿在碰该文件）。

## ② 三形判据与名字

新文件两枚顶层用例：

**`TestCanonicalizeFailureFailsClosedInReply174`**（正向＋不越界）
夹具＝`ArtifactPath: "   "`（纯空白）＋**判定者已接线**（`NewPathCanonicalizer([]string{dir}, nil)`，`dir` 是真树）。这枚输入的形状：非空 ⇒ 过 `task.go:242` 的 `rec.ArtifactPath != ""`，进指针支 ⇒ `paths.go:120-122` 的 `Canonicalize` 在 `risk.Resolve` 之前就返回 `tools: empty path` ⇒ 走 `task.go:325-326` 那一支。**这是"规范化都没通过"那一形，不是"规范化成功但不在授权根内"那一形（后者已由 `TestPointerOutsideAuthorizedRootSpeaks` 钉住）。**
断言（逐条各守一块地）：
1. `!out.IsError && out.Truncated`——**不许空成功**：截断仍要答，不能变成拒绝也不能变成"什么都没发生"；
2. 含 `注意：` 且含 `读不到`——**按读不到处理**；
3. 含 **`连规范化都没通过`**——**这一支自己的说法**（关键：摘掉它不能躲在别的支的句子里）；
4. 含 `"全文见 " + 夹具串`——指针照旧给（`pointerRe` 在此夹具上不可用，见 §6 的现量说明）；
5. 禁字表：`读得回来`／`随时可读`／`可以读回`／`不在你被授权的目录范围内`／`不可找回`／`未接线`——**回执没把这条路径说成能读**，也没把别的支的判词冒充过来；
6. 回执**不含 `dir`**（授权根列表不得转述给模型）。

**`TestCanonicalizeFailureIsNotTheUnwiredArm174`**（正向第二枚＋两支持久不粘连）
同一夹具摆两臂：接线臂必须含 `连规范化都没通过` 且**不含** `未接线`；`Paths: nil` 臂必须含 `未接线`＋`按读不到处理` 且**不含** `连规范化都没通过`。
⇒ 挡掉两种"修法"：把规范化失败改口成"未接线"（或反向），以及两臂共用一句通用文案。

⚠ **没有**任何断言钉精确文案字节：`:299`／`:349` 两枚模板冻结钉一字未动（复量见 §5）。

## ③ 变异现量（摘哪处会红，逐字）

载具＝`-overlay`（**跟踪件零改动**；`internal/tools/task.go` 从未被写）：
`D:/work/workspace/projects plans/Wisp/.scratch/wisp/probes/174/r2/{mut/task-M1.go,mut/task-M2.go,overlay-m1.json,overlay-m2.json}`

| 变异 | 摘的是哪一处 | 落地证明（逐字 `grep -n`） | 红了的用例 | 顶层 |
|---|---|---|---|---|
| **M1** | `task.go:326` 那一支的文案整体摘掉（＝174-v1 的 `MUT-V5` 同形） | `mut/task-M1.go:326: // MUTATION 174R2-M1 (probe only, never committed into internal/): this` ＋ `sed -n '323,330p'` 显示 `case err != nil:` 之后**只剩注释**、`notes = append` 那行不在 | `--- FAIL: TestCanonicalizeFailureFailsClosedInReply174 (0.00s)`<br>`--- FAIL: TestCanonicalizeFailureIsNotTheUnwiredArm174 (0.00s)` | `FAIL github.com/CarlosShao/wisp/internal/tools 15.608s` |
| **M2** | 同一支换成 `return ""`（＝"当作可读"） | `mut/task-M2.go:329: return ""`（紧跟 `case err != nil:` 与两行 MUTATION 注释；`:340`／`:356` 的 `return ""` 是既有码） | 同上两枚，逐字同名 | `FAIL github.com/CarlosShao/wisp/internal/tools 15.365s` |

**未变异基线**（同一发 `-run TestCanonicalizeFailure`）：
`--- PASS: TestCanonicalizeFailureFailsClosedInReply174 (0.01s)`／`--- PASS: TestCanonicalizeFailureIsNotTheUnwiredArm174 (0.00s)`／`ok github.com/CarlosShao/wisp/internal/tools 0.065s`。

**名册之外零枚红**：两次变异跑的都是**整包** `-v`，`grep '^--- FAIL'` 只出这两枚 ⇒ 判据不粘连别人，别的判据也不替它守这块地。
**历史对照（本票的正面结论）**：`174-v1` 的 `MUT-V5` 摘同一支 ⇒ `rc=0`、顶层 **PASS 124 全绿**；同一形状今天 ⇒ **2 枚红**。AC#2d 的"零覆盖"在测试面这一侧被合上。

正向读数逐字（`-v` 现量，本程自己跑的）：
```
canonicalize-error verbatim: IsError=false Truncated=true announce="[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：这条路径现在读不到，C26 连规范化都没通过（tools: empty path）；注意：这条路径现在读不到，宿主登记的那份副本文件并不存在，全文见    …]"
unwired control verbatim: announce="[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：路径授权判定者未接线（fail-closed：C26 没接进来，这条路径是否还读得回来无法核实，按读不到处理），全文见    …]"
```
⚠ 本机现量：`os.Stat("   ")` 报错 ⇒ 存在性那记说明**也**响了（`注意：` 出现 **2** 次）。这是 OS 事实、不是判据射程，**故意没钉枚数**。

## ④ 本程没测什么 / 没钉什么

1. **没接线**：`cmd/wisp/run.go` 的 `TaskDeps` 今天仍不带 `Paths` ⇒ **真机走的是 `:319` 那一支**，本表钉的 `:325-326` 只在接缝上可达（AC#2b 名下，欠下一枚腿，见 §7）。
2. **`risk.ErrReparseDenied` 那一形本程未造**：它要在盘上有一枚真 junction（`mkRealJunction` 只在 windows-only 的 `bridge_junction_windows_test.go` 里），本程为可移植性选了"纯空白路径 ⇒ `tools: empty path`"这一形。⇒ 同一支的第二枚入口今天仍无判据。
3. **AC#2d 票面后半句"拼进去的 `err.Error()` 里不许带路径原文／根列表／C26 内部状态"——本程只能钉住前两样**：
   - 根列表：已钉（断言 6，回执不含 `dir`）；
   - 路径原文：这一形里夹具串就是空白、无从断言（`strings.Contains(text, "   ")` 恒真）；
   - **C26 内部状态：今天必带**——`task.go:326` 就是 `+err.Error()+` 原样拼接，现量回执里那句 `（tools: empty path）` 正是包内错误文本。⇒ **钉"不许带内部文本"会让判据在未修码上红＝要改产码**，本程按派单**停手上报**（登记事实，不写断言、不改产码一个字）。
4. **`spill.go` 那枚桩（AC#2c）一字未看未动**；票 177 的豁免落点本程不评。
5. 端到端 `wisp run` 未跑（本机缺 DLL 是 174-c1 的现量事实）。
6. 全仓其它包的门禁本程未跑（派单只列四把尺）。

## ⑤ 门禁终态

（本节的数取在**本程最后一枚 commit 之后**；若终表与提交顺序有先后，以 `git log` 为凭）

| 尺 | 读数 |
|---|---|
| `go test -count=1 ./internal/tools/` | 待填 |
| `sh scripts/d22scan.sh` | 待填（`ban #8 internal/` 枚数解释见下行） |
| `bash .scratch/wisp/probes/154/gate-clauses.sh` | 待填（**比红腿名册不比退码**） |
| `gofumpt -l internal/tools/task_output_canonicalize_fail_174_test.go` | 空（第 14 枚调用现量：`gofumpt-rc=0`、零行输出） |
| 两枚模板冻结钉 md5 复量 | 待填 |
| 产码零改动自证 `git diff --numstat 40aee084..HEAD -- internal/ cmd/` | 待填 |

## ⑥ 被拒调用 ＋ 零删除 ＋ 终值

- 被拒／失败调用：**零枚被权限系统拒绝**。一处**仪器自伤**要如实登记：第一次并门禁那发里 `grep -cE '^--- FAIL...'` 命中数＝**0** ⇒ 退出码 1 ⇒ 同条 `&&` 链在它之后的 `sh scripts/d22scan.sh` **没跑**（我据此重跑了一发）。这不是被拒，是链式写法把"零枚红"读成了失败——**读数只作废那一发里 d22scan 的缺位**，`ok 15.523s`／`=== RUN 186`／`FAIL+SKIP 0` 三行是真读数。
- **零删除命令**：全程只用 `mkdir -p`／`cp`／`printf > 新文件`／`Write`／`Edit`，没有 `rm`／`rmdir`／`git clean`／`git restore`／`git checkout .`，也没有动 `probes/161/r6/flip-declaration.sh`（未跑）。
- 还原自证：变异只在 `.scratch/wisp/probes/174/r2/mut/` 的副本上；`internal/tools/task.go` 从未成为写入目标 ⇒ `git status --porcelain -- internal/ cmd/` 只应出现**新建的 `_test.go`**（别人在飞的 `internal/panel/git.go` 不是本程的，见下）。
- 共享树噪声登记（**不评论、不动**）：终态写面出现 `?? internal/panel/git.go`＝同批 `181-r1` 的在飞件，非本程产物。
- 工具调用终值：见 §7 末行。

## ⑦ next

- **接线那一行还欠谁**：`cmd/wisp/run.go` 里构造 `tools.TaskDeps{}` 那一行补 `Paths`——**AC#2b 名下、下一枚写腿的活**（174-c2 §1 给的现量号位是 `:365`，行号请它自己复算）。派单 §1 明写本程不许碰，且同批有写腿正在改 `cmd/wisp/run.go`。
- 接线之前，本表这两枚判据是**接缝级**：真机仍只会说"未接线"。
- 接线之后要做而本程没做的两样：①`ErrReparseDenied` 那一形（要真 junction，windows-only 台件）；②`err.Error()` 文本泄漏面（现量在 §4.3，**要改产码** ⇒ 是契约/措辞轴，须先有人裁一句 `PLAN.md:2564` 与 `task.go:326` 的关系）。
- 本程**没碰** AC#2c／票 177／`allowlist.txt`／`thresholds.go`／golden／审批超时常量，`docs/PLAN.md` 与 `docs/specs/**` 零字节。
