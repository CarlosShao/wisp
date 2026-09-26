# 138 — 取消一轮任务时，那一轮**自己拉起来的子进程没有归属**：`Cancel()` 只掐 context，全仓零枚 `TerminateJobObject`，唯一的 `JobScope` 是**整机一枚**

**Status:** **`blocked`（09-26 00:5x 编排者按 138 验收程的第 4 节判；等 `Q-55` 那一次批准一起落）**
              ⇒ **不走 `-done`、也不作废**：验收程判"正解是池内规则 5＝`blocked` ＋ `Q##`"，
              并点名我票面 `:46` 那个 `ready-for-human` **在池子词表里不存在**（`README:9-15` 只有
              `ready-for-agent`／`in-progress`／`blocked`／`review`／`done`）——那一处是**我写的票面缺陷**，
              原句不抹、更正见文末 Progress log。
              ⚠ 同时登记：我票面 `:12`（冻 `docs/specs/*.md`）与 `:44-46`（命令去 `SPEC-12 §5` 登记）**是同票互斥**。
              ／ 上一状态 **`ready-for-agent`**（原文照抄在下一行，不覆盖它）

**Status:** ready-for-agent（2026-09-24 10:0x 编排者建；来源＝owner 批准的 `Q-43` 前半，
              原始缺口＝`docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md` **GAP-09**，
              裁决见 `docs/reports/2026-09-24-gap-analysis-audit-verdict.md` §2 表第 09 行。
              ⚠ **前提**：六个审核员里对这一条的结论是「28 条里最干净的一条、可直入排程」，
              但它同时判定**原文的"最坏后果"写歪了**（详见 AC#2）⇒ 本票按**核过的形状**做，不按原文做。）

**Packages:** `internal/proc/**`（`jobscope_windows.go`）· `internal/agent/loop.go`（取消路径）·
              `cmd/wisp/**` 若有接线点。**⚠ 冻结面照旧禁改**：`internal/risk/**` · `rules_gateway.go` ·
              `tools/d22scan/**` · `allowlist.txt` · D32 阈值（CPU≤0.5%／RSS≤25MB）· `thresholds.go` ·
              任何 golden · `frontend/**`（owner 已交外部 agent）· `docs/PLAN.md` · `docs/specs/*.md`。

## 0. 编排者现量的两条锚（不是审计代理转述的，是我自己刚走的）

| 事实 | 读数 |
|---|---|
| 取消动作到底是什么 | `internal/agent/loop.go:300` —— `func (t *RunningTask) Cancel() { t.root.Cancel() }`。**只有 context 取消，没有任何进程动作。** |
| Job Object 在仓里有几枚 | `git grep 'type JobScope'` ⇒ **仅 `internal/proc/jobscope_windows.go:62` 一枚声明**；`git grep 'TerminateJobObject'` ⇒ **全仓 0 处** |

⚠ 上面第二行**只到"声明"这一层**。"这一枚 JobScope 是不是全进程共用、有没有生产调用者、
有没有 `AssignProcessToJobObject`" 一律**留给 AC#1 现算**，本票不预先假定（那正是审计代理给出、
我未复算的那半截）。**若 AC#1 量出与 §0 或本票判据打反，停手报回，不许按票面硬改。**

## 1. 为什么现在立案

GAP-09 的原话是：「取消只有 `RunningTask.Cancel()`→root ctx；JobScope 全进程一枚、无
`TerminateJobObject` ⇒ 任务级级联确无归属」。裁决表接受其**结论**、否掉其**最坏后果**，并留下一条硬约束：
**不许顺手改 `cancelled` 文案**（那要动 D37 的 17 类错误表＝契约变更）。本票就是按这两条拆出来的形状。

- **要的东西**：一轮任务（一次用户请求）自己起的那些子进程，在用户按"停"之后，**跟着这一轮一起没了**。
  今天的形状是"context 取消"——它只让**我们自己的 Go 代码**不再等，
  **不保证那个已经被 exec 出去的子进程收到任何信号**。
- **不要的东西**：批量授权、审批文案、错误分类，一概不在本票范围内。

## 2. 结案判据

- [x] **AC#1**（把"无归属"从断言变成读数）**现算三件事并逐名给表**：
      ① 那一枚 `JobScope` 的**创建点数**与**调用者**（生产码里几枚、测试里几枚；是不是整机共用一枚）；
      ② 仓内有没有 `AssignProcessToJobObject`（及等价形状）；有 ⇒ 谁被分配进去过；
      ③ **今天哪一条生产路径真的会 exec 出子进程**：逐枚列出候选（`fs.*` 之外还有哪些），
      并注明 `shell.exec` / D46 命令插件今天**注册了没有**（`internal/tools/unwired.go` 那条
      "no shell.exec tool is registered" 的注释是不是仍然属实）。
      ⚠ **这一格允许得出的结论是"今天无入口踩得到"**——那也是有效交付，但要**同时**答：
      那它是"已实现功能里的缺口"还是"给未来功能预留的洞"？若是后者 ⇒ **必须**去 `SPEC-12 §5` 按五字段
      登记为 RESERVED／DEFERRED，并把本票降级为 `ready-for-human`，**不许**为了结案去造一条假的现实危害。
- [ ] **AC#2**（修法形状，**先裁再写**）如果 AC#1 判定确有可走路径 ⇒ 修法只许落在
      **"每轮任务一枚自己的 Job Object、取消时级联"** 这一族，且必须：
      ① 不动 D4／C19 的任何门控判定（`risk` 侧一个字不改）；
      ② 不新增重量级依赖（D22 闸门②依赖白名单）；
      ③ Windows 专属 API 必须留在 `internal/proc/**` 的 `_windows.go` 半边，
         **非 Windows 侧要么有等价实现、要么明确"此平台无此能力"并留可见告警**
         （不许静默 no-op——那正是本仓抓过的"一步存在却从不产出结论"）；
      ④ 若判定需要改 `cancelled` 相关文案或错误分类 ⇒ **停手上报**（那半截属 D37 契约变更，不在本票）。
- [ ] **AC#3**（牙）**必须有一枚用例在"级联"被摘掉时转红**，且**逐名列出**它红在哪一行；
      先证变异落地（`grep -n` 到你改那一行的原文 ＋ `go build ./...` rc=0）再读数。
      ⚠ 不许用 `t.Skip` 或调高阈值换绿；不许改断言去迁就实现。
- [x] **AC#4**（门禁）`gofmt -l`／`gofumpt`（⚠ **宿主已装 v0.12.0，直接跑现成 binary，
      绝对不要执行 `go install mvdan.cc/gofumpt@latest`**——那会改写宿主工具）／
      `go vet` 宿主原生 rc=0 ／ `sh scripts/d22scan.sh` rc=0 且各 scope 不降 ／ `-count=2 -v` 两形四数
      ＋**名册差集**（四数之外还要比 `=== RUN` 与 `--- FAIL` 逐名名单，一条 panic 会吞掉同包其余读数）。
      ⚠ Linux 交叉 vet 对 `cmd/wisp` 与 `cmd/balldebug` **本来就 rc=1**（sherpa build constraint），
      与被审对象无关 ⇒ 按包作用域跑，别把它当真伤报回，也别拿它当"整树假象"放过别处。

## 3. 明确不在本票范围

- `cancelled` 文案／错误分类（D37 17 类表，契约级）—— AC#2④ 已划出去。
- 副作用清单与 per-call `AppliedSteps`：归**票 21 段 2**。
- "撤销本次任务的全部改动"（回退/checkpoint）：那是 **GAP-02**，裁决表判为**需 owner 拍板**，
  且与本票**不是一件事**（本票管"进程还活着吗"，那票管"文件改动能不能退回去"）。

## 4. 进度日志

- [2026-09-24 10:0x +08] agent=orchestrator did=建票。
  来源＝owner 在对话里批准 `Q-43`（原话「都按你的推荐来」）；撤销口令＝**「138 撤」**（改回未立案、票文件删除）。
  自己现量了 §0 那两枚锚（`loop.go:300`、`type JobScope` 唯一声明、`TerminateJobObject` 零命中）；
  **审计代理报的"全进程一枚／无调用者"那半截我没复算，已写进 AC#1 要求现算**。
  next＝排在今日队列之后（137 AC#3／AC#4、135 原六格、136 AC#10／AC#11、#40 全树终判据复算）。

- [2026-09-26 00:5x +08] 编排者按 138 验收程（`docs/evidence/s1/138-cancel-kills-children-r1-accept-r1.md`，6 格／6 枚 commit、零代码改动）处置本票。
  **`AC#1`／`AC#4` 我已翻勾；`AC#2`／`AC#3` 原句不动、也不勾**——它俩的"前件被判假⇒未触发"是真的（今天没有可走路径），
  本票**既不作废也不 `-done`**，按验收程第 4 节改判 **`blocked` ＋ 等 `Q-55` 那一次批准一起落**。
  > **更正一处、认错两处（原句不抹，逐条在下面追加）：**
  > ① **`:46` 那句 `ready-for-human` 是我写的票面缺陷**：`README:9-15` 的状态词表里**根本没有这个值**
  >    （只有 `ready-for-agent`/`in-progress`/`blocked`/`review`/`done`）。正确写法就是本行上面那个 `blocked`。
  > ② **`:12`（冻 `docs/specs/*.md`）与 `:44-46`（命令去 `SPEC-12 §5` 登记）在同一次批准之前互斥**——
  >    这是我自己开票时没量到的自相矛盾；验收程点名了，处置是"两行一起交给 `Q-55`"，不由实现方硬挑一边。
  > ③ **`A267③` 我那格"它引用了我没说过的话"降级为未定性**：验收程现量 **票面整份零枚 `CI|workflow|runner|gh run` 字样**（比我说得强，排除"改写票面"那一支），
  >    但**我的派单正文在盘上没有凭据**（`.scratch/wisp/dispatches/` 的规矩当时只归档"授予范围／解冻／例外"那类，普通实现派单不入盘）
  >    ⇒ 所以"我派单里没这句"是**我对上下文的回忆、不是一件可核的物证**。甲／乙／丙／丁四支**一支都不能定罪**。
  >    ⇒ 制度已改：**从本行起，每一枚派单正文落盘进 `dispatches/`**（见该目录 README 末段追加的那条），今后这类争议有物证可查。
