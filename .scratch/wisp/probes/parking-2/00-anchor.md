# parking-2 起手锚＋检查笔记（2026-10-08 18:2x +0800；起手 HEAD `351ba36`）

本件＝文档腿 `parking-2` 的检查笔记（⛔ 不删、只追加）。派单要点：唯一任务＝**把停车点文档 `docs/reports/HANDOVER.md` 刷新到当前状态**，**只许在真末尾追加一节**（节头带 `date` 现量与 HEAD），⛔ 不回改任何既有行、⛔ 不删；⛔ 零 Go 命令；⛔ 不改票面/台账/不新建票；commit 带显式 pathspec；⛔ 零 push；除 HANDOVER 那一节外只许新建本目录 `.md`（临时件建仓库外、只建不删）。

## 1. 起手现量（2026-10-08 18:20–18:30 +0800，逐把现跑）

- `date`（本程第一发）＝`2026-10-08 18:20 +0800`。
- `git log -1`（起手）＝`351ba36238fc8b1438b0cd84a1ce3375f8a6ac43`，18:20，「253-v1 起手锚（HEAD 8dd239f1，靶= panel_dispatch_binding_roster_253r1_windows_test.go@65f4c968）」。
- ⚠ 落笔期间 HEAD 又前进四次（`git log` 现量）：`55488e3f`（18:27 `A726` 落账）／`14115e8d`（18:27 `comment-fix-prep-1` 起手锚）／`fdd2a72e`（18:28 `A727` 落账）／`6538b3c1`（18:30 `ledger-audit-1` 起手锚）。
- `docs/reports/HANDOVER.md` 追加前：`wc -l`＝**1882**；`git status --porcelain`＝**空**（rc=0）；哈希三面同值＝`0dd29f8357e4778cc871e8a52574d3e9d45dfffe`（尺＝`git rev-parse HEAD:docs/reports/HANDOVER.md`＋`git hash-object` 工作树＋`head -n 1882 docs/reports/HANDOVER.md | git hash-object --stdin`，三发同值；且 `git log --oneline 351ba36..HEAD -- docs/reports/HANDOVER.md` 回空 rc=0＝自起手起无人动过它）。
- 节头枚数（追加前）：`grep -cE '^## '`＝**38**；`grep -cE '^> \*\*[0-9]'`（本文件块引用日志节头形状）＝**57**。
- `git rev-list --count origin/dev..HEAD`＝**401**（18:30 现量；⚠ 共享树里会漂，引前自取）。

## 2. 真尾部末 3 行（追加前；机器抄录＝`tail -n 3 docs/reports/HANDOVER.md`，⛔ 未改一字节）

```text
> ⚠ **一句待答（不是否决、是别让我悄悄改契约）**：这句顶到两行既有定案——`PLAN.md:1544` REJECTED「多 Agent 协作…单 Agent，无编排」＋`:1447` S7 前不得开放多任务并发。**要的口令＝「子代理这条改判」**；不落它我只做到"能画的那部分"，实体层按住。另 `StreamLog` 今天**键上限 32、溢出"合并而不是丢"**（`pump.go:275/:287-290`）对子代理是错形状，197 要重裁＋配结构判据（造 N＋1 枚，不许任何一条被并进别人的）。
> **波次五枚的落点（账 `A392`–`A395`）**：`182-a1` 堆数 7 复跑对格（**口径＝原型 `design/doubao/demo` 注册表**，且实质订正了早上那份普查的成员）；`145-r2`／`188-r1` 两程**零产码停手上报**（`AC#6`"宁缺毋造"首次真生效；188 挖出**两套状态词表并存**⇒ 新立**票 196**）；`189-a1` 量清纯 Go 读 git 的边界（已跟踪那半 51.6 ms 与权威逐档一致；未跟踪那半朴素遍历把 46 枚算成 8,684 枚**且不报错**⇒ 算不出就报"读不到"，绝不报 0）。**我自己两处自纠**：`run.go:442` 那把锁是我自己划的禁区（已具名只放开两处并派 `145-r3`）；往跟踪文件写含反引号中文用了 `printf` 被 shell 执行掉两个符号（已补回，定式升级＝只走 `Edit`/`Write` 或 `cat >> <<'EOF'`）。
> **此刻（18:1x 现量 `date`）**：在飞 **2 枚**＝`145-r3`（写·`internal/panel/**`＋`cmd/wisp/run.go:442`/`:365` 两处）／`188-r2`（写·`internal/tools/task.go`＋`task_backfill.go`，状态值只准取 `states.go:11-31`）——写面不相交。**串行队列**：两枚交回 → `189` 落地腿（嵌在 `composer.git` 下开子格、不加顶层 key）→ `190` 落地腿 → `185-r2` → 票 194 堆1 名册补齐 → 票 191 终端（甲）→ 票 192 预览／联网 → 票 181 AC#7 → 票 186／187 → 票 178；**票 197 实体层按住等「子代理这条改判」**。**推送继续按住**。**待 owner：两句**＝「子代理这条改判」／`196` 甲或乙（默认甲＝只加映射表不动 schema）。台账到 `A397`。
```

## 3. 与派单转述的出入（具名，⛔ 不照抄）

1. 派单说"在飞…`253-v1`＋`167-a2`"：盘上现量（18:27–18:30）＝`167-a2` 17:58 已由编排者**代提死腿遗产**（`d33d492d`，标〔未验证半成品〕，死因逐字＝`You've reached your daily usage limit for Chat.`）；其续程 `167-a2b` **18:27 已交**（`A726`，`fa86e2d5`＋`9e8477ed` 两笔、三枚 `.md`）。⇒ 现量在飞＝**3 枚**＝`253-v1`（18:20 `351ba36`）＋`comment-fix-prep-1`（18:27 `14115e8d`）＋`ledger-audit-1`（18:30 `6538b3c1`，自登记句逐字"只读复核腿，复跑 `A720`-`A727` 十枚承重读数"）。
2. 派单说新建路径 `.scratch/wisp/probes/parking-2/`：盘上先例＝`probes/parking/1/`、`probes/parking/2b/`（**无** `parking-2` 目录、也无 `parking/2`）。本程按派单字面建 `probes/parking-2/`；与先例并存＝事实，⛔ 本程不合并、不删、不迁移。
3. HANDOVER 真尾部＝**2026-09-28 18:1x 的旧日志段**（块引用"再接"体，`A368`–`A397` 那两天）；§4 区"新会话从这一节起读"的最新节＝`## 4.0ad`（**第 261 行**，10-08 11:0x 版）。本程追加点＝文件真末尾（第 1883 行起），与 §4 区相隔约 1600 行 ⇒ reader 入口两处并存（事实登记；⛔ 本程未动 §4 区任何既有行——派单硬规矩）。
4. 三面现量时间线：18:28:33 曾见 **1 枚** ` M cmd/wisp/panel_dispatch_binding_roster_253r1_windows_test.go`（＝`253-v1` 的突变窗口）；18:30:39 复量**空**（rc=0）⇒ 该腿窗口已还原。共享树，引前重跑。

## 4. 追加后自证（四把尺）

（追加并 commit 后在此补读数。）

## 5. 本腿 commit 名册

（追加完成后补：c1＝起手锚（本件最初形）／c2＝HANDOVER 新节＋`01-section.md`／c3＝本件读数更新；笔号用 `git log --oneline -- .scratch/wisp/probes/parking-2 docs/reports/HANDOVER.md` 现取。）

## 4. 追加后自证（四把尺；本笔＝c3 补）

追加动作＝`cat .scratch/wisp/probes/parking-2/01-section.md >> docs/reports/HANDOVER.md`（守门三条追加前全过：`wc -l`＝1882、`porcelain` 空、前 1882 行 hash＝`0dd29f83…`；一次通过，18:33）。

- ① `wc -l`：**1882 → 1911**，增量＝**29**＝`01-section.md` 行数（`wc -l` 两处同值 29，逐字相等）。
- ② 节头枚数：`grep -cE '^## '` **38 → 38**（本节＝块引用形，未新增 `## `）；`grep -cE '^> \*\*[0-9]'`（本文件块引用日志节头形状）**57 → 58**（+1＝本节头，位于第 **1884** 行；第 1883 行＝空行分隔）。
- ③ `git diff --numstat 351ba36 HEAD -- docs/reports/HANDOVER.md`＝**`29	0	docs/reports/HANDOVER.md`**（只有这一枚文件、零删除）；同尺限本腿两处＝**`00-anchor.md 35	0`／`01-section.md 29	0`／HANDOVER `29	0`**（三枚逐枚纯增）。
- ④ `git status --porcelain -- docs/reports/HANDOVER.md`（commit 后）＝**空、rc=0**；`-- docs/reports/HANDOVER.md .scratch/wisp/probes/parking-2` 同尺亦空。
- 防吞锚（加厚）：`head -n 1882 docs/reports/HANDOVER.md | git hash-object --stdin`＝`0dd29f8357e4778cc871e8a52574d3e9d45dfffe`＝追加前同值 ⇒ 前 1882 行逐字节未动；追加段 `tail -n 29 … | cmp - 01-section.md`＝**identical**。
- commit 名册（逐笔现取）：**`4263812a`**（c1 起手锚＋本件 35 行）→ **`1ae79210`**（c2 HANDOVER 29 行＋`01-section.md` 29 行＝2 文件 58 增 0 删）→ **本笔（c3）**。⛔ 零 push、三笔都带显式 pathspec。
- ⚠ 后验（⛔ 不回改 HANDOVER 追加段）：追加段落笔后盘上又动一笔——`73b2438c`（18:32，「`253-v1` 对抗验收件：五组盘上突变读数…判'AC#1 成立但附条件'」，逐字见 `git log`）＝**`253-v1` 已产出验收件**；追加段内"在飞 3 枚"的读数时刻＝18:27–18:30，最新状态以本条为准。
