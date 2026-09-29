# 236 — **六枚只在"读数层"暴露的仪器缺口**（`221-v1` 与 `ci-red-1` 两枚非实现者腿交回）：两支 fail-closed 无尺／那把 DEFERRED 尺的 (b) 支扫词面无牙／"任务 id 必须宿主铸造"零钉／`ci.yml` 那句"今天会红"已过期／**在册红名册不是稳定集合**／**一枚刻意坏的夹具被入库之后污染了格式门的 tracked 分母、并吃掉两道 `go vet`**

- **Status**：**待派**（编排者 09-29 20:2x 立，起手锚点 `c5d88a7f`＝本票现读 `git rev-parse --short HEAD`）。
  来源＝`docs/evidence/s1/221-task-cancel-v1.md`（**37,097 字节**，我本人 `wc -c`；三发 `98007cc0`→`f6186158`→`c5d88a7f`）与 `.scratch/wisp/probes/ci-red/ci-red-1.md`（**35,452 字节**，我本人 `wc -c`，未提交）。台账 `A451`。
- **本票的射程不是"功能没做"，是"没人能证明它没坏"**。⛔ **零枚 AC 允许放宽任何现有断言**；凡"改门／改分母"那一支一律标〔契约邻接，要另批〕。

## 现量（起手逐条复算，别信这里的行号）

| # | 事实 | 谁的读数 | 我复跑了吗 |
|---|---|---|---|
| 1 | `taskCancel.Execute` 里 `caller == ""` 与 `t.d.Roster == nil` 两支 fail-closed，摘掉之后**各 0 枚红、rc=0** | `221-v1` 的 m4／m5 | ⛔ 未复跑（腿自述） |
| 2 | 票 221 新写那把 DEFERRED 尺：**(a) 支有牙、(b) 支无牙**——`teeth-m13` 摘掉 `task.list` 的标记行仍 PASS | `221-v1` | ✅ **我现读机制**：`internal/tools/task.go:279` 那行注释**同时含** `task.list` 与 `DEFERRED` ⇒ 词面计数被叙述句撑住（HEAD 命中 2 行：`:23` 与 `:279`） |
| 3 | 全仓**没有任何尺钉"任务 id 必须由宿主铸造"**；只有 `internal/agent/compress_trace_test.go:524` 一枚**形状**检查（36 字符＋4 枚连字符），射程是 trace 那一族 | `221-v1` | ✅ 我 `grep -rn newTaskID --include=*_test.go`＝5 命中，逐条读：`:524` 形状、`:557` 两枚不同、`pointer_183_cli_seam_test.go:59` 是指针注释不是断言 ⇒ **无"谁铸造"这一维** |
| 4 | `.scratch/wisp/probes/185/c1/mut/fs_broken.go` **已被入库**（`4813567e`，09-28 14:42）⇒ 进了 `lint::gofmt (gofumpt) - the tracked set is the denominator` 那一步的分母 | `ci-red-1` | ✅ **两半都我自己跑**：`git ls-files --error-unmatch` **成功**＝tracked；`gofumpt.exe -l` 该文件 **rc=2**、原文 4:1 "imports must appear before other declarations" |
| 5 | 今晚两发 push CI 红名册**逐字相同＝顶层 23 枚（含缩进子测试 25 枚）**，⛔ **不是我 A450 里写的 11 枚** | `ci-red-1` | ✅ 我认：我那 11 是用 `gh run view --log-failed \| grep -oE` 抽的，**任何一条过滤器都复现不出 23**（见 AC#5 那条尺） |
| 6 | `ci.yml:519-524` 那句"这一步今天在 windows-latest **会失败**（`RUNNER~1` 8.3 短名让 A 档判成 B）" | `ci-red-1` 称今晚实读 `--- PASS: TestPathResolverJunctionWindows` | ⛔ 未复跑（腿自述；且需要托管 runner 的环境才知道过期是否长期成立） |

## 为什么必须先修（不是文案洁癖）

- 第 1、3 条是**同一族**：两枚"拒绝"今天靠**代码里的那一行**存在，而**没有任何仪器**能在有人删掉它时报警。票 221 的整条安全结论（子代理不许停兄弟／不许停自己）就建在这两行之上。
- 第 4 条已经在**吃读数**：`lint` 那一步一红，后面的 `go vet (module)` 与 `go vet (tools/d22scan module)` **今晚两发零读数**（项目已知坑"一步红会吃掉后续步，skipped 不产日志＝读数永久采不到"，这次不是推测、是实测到了）。
- 第 5 条影响**所有"逐名比红名集合"的判据**：我这一整天用它做过若干次"零新增红"的结论，而名册本身在两发之间会换位 ⇒ 那条尺**缺了一半**（缺隔离复量那一半）。

## 落点（每格都要现跑读数；⛔ 不许用 mock 代替真的）

- [ ] **AC#1 两支 fail-closed 要有常驻尺**：给 `taskCancel.Execute` 的 `caller == ""` 与 `Roster == nil` 两支各立一枚判据，**形状＝摘掉那一支必然有具名用例变红**（照 `221-v1` 的 m2／m3 那两发的做法：各恰 1 枚红才算有牙）。完成判据＝两发突变读数（红→还原绿）＋工作树 `git status --porcelain -- internal cmd` 为空。残缺表现＝今天有人把任一支删掉，全仓零枚测试会响。
- [ ] **AC#2 那把 DEFERRED 尺的 (b) 支改扫能力**：从"数词面"改成"看 `BuiltinTaskEntries` 的注册名册里有没有 `task.list` 这一行"，并**先证它今天无牙**（`teeth-m13` 那形：摘标记行仍 PASS＝未修码读数）。⚠ 与票 225 分开算账：225 管"标记与 `SPEC-12 §5` 双向对账"，本格只管**这把尺本身有没有牙**。⚠ 附一条仪器事实要写进判据注释：**读盘型尺（`os.ReadFile`）对 `-overlay` 结构性不可见**，测它的牙必须 overlay 替换测试文件里的读路径并配正控。
- [ ] **AC#3 "任务 id 由谁铸造"要有一枚钉**：先量清今天有几枚尺会在"task id 换成模型给的值"时变红（我现读＝**0 枚**，只有形状检查），再定钉的落点。⚠ 背景两条硬事实：`internal/tools/task_backfill.go:38` 明写**普通 spill 路径的 call id 是 model-supplied**、且"裸 task id 是完全可以被够到的 call id"；`agent.newTaskID()` 熵源失败那一支**退化成时钟派生**。⛔ 本票**不许**顺手把 `allowed_dirs` 之类做成硬边界、也不许新造工具面（D34 无对应行）。
- [ ] **AC#4 `ci.yml` 那句过期注释要不要动**：只改注释、零行为改动；改前先**现读两发**确认 junction 那步真的 PASS（`ci-red-1` 给的是今晚两发，别拿它当长期事实——托管 runner 的环境会变）。⚠ 低利害、可逆，**不该拿去找 owner 拍**，由编排者或下一位读 CI 的人顺手定。
- [ ] **AC#5 "逐名比红名集合"这条尺要自带隔离复量那一半**：把规矩写成两发——**整包名册红 ＋ 单包隔离复量**（`221-v1` 实测：`cmd/wisp/TestTicket223RefusedLooseningKeepsOldValues` 整包 5.15s 红、隔离 **3/3 绿**（5.055／5.019／5.034s）＝计时红，不记账）；并把"在册常红名册今晚与实现腿那发不一致（两枚计时红换位）"重钉一次，落点写进哪份仪器说明由执行者定。⚠ 顺带纠正一条**我自己的坏尺**：`gh run view --log-failed` 抽名册会**漏计**（我抽到 11、真数 23），且**步名归属在部分 run 上全是 `UNKNOWN STEP`** ⇒ 数名册必须与 `portable-tests.sh` 自报的 four numbers **互咬**（今晚 4／7／12 三发各自相加＝23 对得上），并写清"名级 23 ≠ 包级 3"两个口径。
- [ ] **AC#6 tracked 分母被入库的刻意夹具污染**：`.scratch/wisp/probes/185/c1/mut/fs_broken.go` 是**故意编译不过**的正控夹具（`undefinedSymbol185c1`），它入库之后合法地进了格式门的 tracked 分母 ⇒ `lint::gofmt (tracked set)` 常红，并**吃掉后面两道 `go vet`**。三支摆清楚（⛔ 本票不预选，因为涉及门的分母＝〔契约邻接〕）：
  **甲**＝只把该文件的 **import 顺序排好**、保留"未定义符号"那一枚坏点（它的用途是符号坏、不是格式坏）——最小、零门改动；⚠ 代价＝它是票 185 已归档读数的对象，**改字节会洗掉别人引用的读数**，动之前要先核 `185` 名下有无 md5／字节数引用它。
  **乙**＝门的分母排除 `mut/` 这类刻意夹具目录——**动 CI 形状**，照票 134 的 C+B 先例要人工批准。
  **丙**＝不动，只把"这一步今天为真红、且它会吃掉两道 go vet"写进 CI 拓扑文档与停车点，让后续程别再当"CI 没测"是环境噪声。
  完成判据＝任选一支之后**必须真拿到一次 `go vet` 读数**（不是"应该能拿到"）。

## 禁区

- ⛔ 不放宽／不删除任何现有断言；不许为了变绿把 `lint` 那一步改成 `continue-on-error`。
- ⛔ **临时件只建不删**（`issues/README` 规则 8）：不许用"删掉 `fs_broken.go`"当作甲的省事做法。
- ⛔ 不碰 `frontend/**`／`design/**`、`PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`、三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）。
- ⛔ 工作树里别人的 46 枚脏改动（含 `docs/evidence/s1/152-*.md`、`probes/161/r6/logs/**`）一枚不许顺手提交。
- 派单必写：**产码腿的 AC 框一枚都不许碰**（票 221 刚为这一漏撞过一次，`A450`）；凡写进派单的 grep 尺先在本机跑一遍并把真实读数抄进派单。

## Progress log (append-only, newest last)
