# 归属收口腿 `close3-r1` —— 三格判性质（1×①／2×②）；① 那格＝票 175 AC#5 三态分账（零产码）

- 腿：`close3-r1`（写腿·实现者）｜锚 `385afead`｜分支 `dev`｜**只 commit、零 push**
- 交件：`docs/evidence/s1/close3-r1-impl.md`（本文件）｜**产码 diff＝0**（§5 自证）｜票面／台账／`docs/specs/**`／`docs/PLAN.md` 一字未动
- 一句话：**175→①（已做：写分账）／246→②（归口票 228）／316→②（归口票 272，票已 `withdrawn`）**；后两格**具名交回、⛔ 硬做**（上一批「剩 1 格」的教训）。

## §1 七项读数（速览）

| # | 读数 | 本腿 |
|---|---|---|
| 1 | 三格性质 | **①／②／②**（逐字依据＝§2） |
| 2 | 做了哪几格 | 仅 `175 AC#5`（§3）；246／316 零写面、零产码 |
| 3 | 反向判据 | 175 AC#5 的完备性机械尺：**变异红×3／还原绿／原样绿＋还原核**（§4） |
| 4 | 零回归 | `internal/tools` 改前 `ok 20.844s` → 改后 `ok 16.518s`；`internal/risk` 改前 2 枚红 → 改后 1 枚红，**两枚都是既存**（§5） |
| 5 | 四数 | 改动包 **0**（0 文件／0 行 Go）｜`cmd/wisp` **0 字节**｜`internal/speech` **0 字节**｜`internal/{session,config,secret,memory}` **0 字节** |
| 6 | 格式两尺 | **导出树（HEAD）＝甲/乙 `0/0`**；共享树＝甲/乙 `6/6`＝**既存 3**（worktree CRLF 假红）＋**`wire-r1` 正在写 3 枚**（逐名 §6） |
| 7 | `d22scan` | **导出树 HEAD＝`rc 0` clean**（正控 `PASS=35 FAIL=0`）｜**导出树＋本件＝`rc 0`**｜共享树 `rc 1`＝3 findings **全在 `wire-r1` 的未跟踪文件**（§6） |

## §2 三格性质判定（逐字依据＋理由）

### 2.1 票 175 最后一格 → ①「能做：写分账」

- **逐字**：`- [ ] **AC#5 与票 164／174 分账**：本票结题时必须逐条写明它跟 Q-59（偏移）、票 174（授权根）各是什么状态——三枚是同一根管子（"读回外部内容"）上的三处，批一枚不通另一枚。`
- **理由**：正文动词是「逐条写明」，产物是**状态说明**；无「归 X」／「X 落地之后再谈」／「归机主」字样，也不要求动产码 ⇒ **①**。可落面＝本腿写面允许的 `docs/evidence/s1/close3-r1-*.md` ⇒ 已落成 §3。
- ⚠ **一处交接**：台账 `A955`（10-14 20:14）已把本格记成「编排者桌面（随其族收口时做）」⇒ §3 是**可直接粘进票 175 Progress log** 的底稿；**框仍归编排者翻**（本腿 ⛔ 改票面）。

### 2.2 票 246 最后一格 → ②「归口别的票」

- **逐字**：`- [ ] **AC#8（新补）那枚 boot 时序 flake 要有归属**：修法在常驻启动路径（候选＝把 signal.Notify(bootExit) 提到 installLogSink 之前，消掉 [sink→Notify] 那段窗口），**地界属票 228**（用例名就是它的 TestAC228…）⇒ 本格只做"归口＋判据"，⛔ 不在本票改；判据＝整包序 -count=10 该用例零枚 0xc000013a＋红，且不得靠放宽断言或加 sleep 达成。`
- **逐字 2**（同票）：「`AC#8 保持未勾、保持归票 228`」｜「`本票状态：勾 8／未勾 1（只剩 AC#8，地界在票 228）`」
- **理由**：票面自己写「**地界属票 228**」「**⛔ 不在本票改**」⇒ **②**。另有两条独立硬边界同时挡住：① 判据落点 `TestAC228…` 在 `cmd/wisp/**`，本波该面归 `wire-r1` 独占（`cmd/` 命中必须 0）；② 判据要的「整包序 `-count=10`」今天撞票 314 的整包死亡区（台账 `A957`：9 趟里 2 趟整进程死亡＋1 趟 10m 超时）⇒ **不是不愿跑，是今天测不出干净读数**。
- **交回**：票 228（修法候选已由票面给出：`resident_windows.go` 的 `[sink→Notify]` 窗口）＋编排者（`cmd/wisp` 面／314）。

### 2.3 票 316 最后一格 → ②「归口票 272；票已 `withdrawn`」

- **逐字**：`- [ ] **AC#1 收口（三选一，落地腿现读后裁并具名理由）**：` ＋ 三候选 `(a)` 窗口内恰一次／`(b)` 标记取在 `state=armed` 之后／`(c)` 句偏移晚于 `HOT-RELOAD state=applied` 行。
- **逐字 2（Status）**：`**withdrawn（撤票；口令「316 撤」已用）**…**与票 272 重复**`
- **逐字 3（撤票记录）**：`"正例 G 仍绿"与判据①**互斥**`｜`**裁定**：arm1 保留`｜`去向：本票读数并入票 272 的收尾（272 现只剩 AC#5／AC#7 未勾）`
- **理由**：**②** —— 票面写「与票 272 重复」且去向具名 272；更重的是「现读后裁」**已由非实现者 `316-r1` 做完**（三支逐支现读：`(a)` 已在盘 `:668`／`(b)` 空操作／`(c)` 结构性误伤在正常一发 FAIL 42.20s）⇒ 再做＝**重审一枚已撤的票**，且与 272 的 arm1 裁定（M 今天已红、FAIL 4.38s 红在 `:596`）冲突。撤票口令已用 ⇒ **⛔ 重开**。
- **交回**：票 272（`-done`，读数已并入其 AC#5／AC#7）＋编排者（⛔ 重开票）。

## §3 落地格：票 175 AC#5 三态分账（现读，锚 `385afead`；可直接粘进票面 Progress log）

> **行号口径**：本 §3 与 §4 的每一条 `path:line` 都是 **`git show HEAD:<path>` 现读**（锚 `385afead`）。⚠ 共享树里同几枚文件正被 `wire-r1` 改（`cmd/wisp/run.go` 的 `RunAsync` 在 worktree 已从 `:1374` 漂到 `:1380`）⇒ **本件按 HEAD 记锚**，避免读数随别人的写而腐烂；§6 的共享树读数都注明了时刻。

### 3.1 `Q-59`（偏移参数）

- **状态**：owner 09-27 21:3x 批**乙（暂不加偏移）**；五字段 DEFERRED 行在盘 `docs/PLAN.md:1538`（行首 `| **DEFERRED** | **超长工具结果的"尾段读不回"**`；解冻条件逐字「需 owner 另批一次（`Q-59` 甲）」）⇒ **乙已兑现、甲未批**。
- **现读**：`fs.read` 参数仍只 `{path, max_bytes}`（`internal/tools/fs.go:113-116`）；`defaultMaxReadBytes = 256 * 1024` 落 `internal/tools/fs.go:55`；`if a.MaxBytes > 0 && a.MaxBytes < limit { limit = a.MaxBytes }`（`fs.go:153`）⇒ **只能改小、不能放大**；`truncated := len(buf) > limit`（`fs.go:160`）。
- **报警器**：`internal/tools/task_output_leg_test.go:346 TestPointerPast256KiBIsNotFullyReadable`（PLAN 那一行逐字：「那条用例红就是这一行到期未做的报警器」）⇒ 今天**绿**（§5）。
- **不通别的**：175 的盖戳修好**不通**这一枚——超 256 KiB 的后半截今天仍读不回。

### 3.2 票 174（授权根）

- **状态**：**回执诚实化＋授权根判定接生产已落**；**契约轴（AC#4）未批**；**AC#5 分账未写**（本 §3 可直接喂它——两票 AC#5 逐字同源）。票面 Status 仍 `ready-for-agent`（未 `-done`）。
- **现读（洞没修的那一层）**：`cmd/wisp/run.go:534-535` 授权根**只**来自 `cfg.FS.AllowedDirs`，无第二来源；artifacts 目录 `:1271`／`:1384` 仍**不在**根里 ⇒ **默认配置下 spill 文件今天仍读不到**。174 修的是「回执不许诺一条读不回来的路」，**不是「给根」**。
- **现读（判定者接生产）**：`cmd/wisp/run.go:553`／`:581`／`:849` 三处 `Paths: rt.paths`；`pointerNotice`（`internal/tools/task.go:812`）据此只说真话。
- **不通别的**：174 修好回执**不通** `Q-59`（偏移），也**不通** 175（读回来之后打不打戳）。

### 3.3 票 175（本票：盖戳）

- **状态**：AC#1–#4 已勾，AC#5＝本格。名册 **9 枚**（含 `SrcTaskOutput`）＋桥级五枚判据在盘。
- **现读**：`internal/risk/provenance.go:94` `SrcTaskOutput = "task.output"`／`:99-102` 在 `sensitiveSourceTools`／`:106 IsSensitiveSource`／`:573` 是唯一闸门；用例 `internal/tools/ticket175r2_stamp_live_test.go:142/183/207/280/349`。
- ⚠ **票面一处残言已过期（本格现读抓到）**：`175-r2`（09-28）写「端到端那一发仍站着不动——`TaskRoster.Record` 生产零写者、`cmd/wisp` 只建空名册」。**今天为假**：起跑口 `cmd/wisp/run.go:1374` `bg := loop.RunAsync(ctx, task)` → `:1382-1386` `tools.TaskBackfill{…}.Backfill(...)` → `internal/tools/task_backfill.go:141` `b.Roster.Record(taskID, rec)`（落地提交 `41e8873bd`＋`06dae8316`，票 176-r1）。读侧也不再被「风险未分级」拒：`internal/agent/loop.go:828-834` 逐字「fs.read / fs.list / task.output arrive here as L0 (ticket 179…)」＋`decideRisk` 的 `case RiskL0`（`:857-860`）直接放行。⇒ **洗戳路今天已通**，175 的盖戳是承重件，不是「等接线」。
- **不通别的**：175 修好盖戳**不通** `Q-59`／174。

### 3.4 一句话分账（三态＝两处半）

> **AC#5 分账（`close3-r1` 现读，锚 `385afead`）**：三枚是同一根管子上的三处，**今天＝两处半**。① `Q-59`：owner 批乙、五字段行 `docs/PLAN.md:1538`、**甲未批**；`fs.read` 仍 `{path,max_bytes}`、256 KiB 上限且只能改小 ⇒ **后半截读不回**，报警器 `TestPointerPast256KiBIsNotFullyReadable` 在岗且绿。② 票 174：回执诚实化＋判定者接生产**已落**，授权根仍只来自 `[fs] allowed_dirs`、artifacts 未入根 ⇒ **默认仍读不到**；契约轴 AC#4 未批、AC#5 未写。③ 本票 175：名册 9 枚＋桥级五枚判据已落；**`175-r2`「生产零写者」已过期**（`run.go:1374` → `task_backfill.go:141`；读侧 L0 由 179 放行）⇒ **洗戳路已通、盖戳承重**。**三处要分别批：批 175 不通 Q-59／174；批 `Q-59` 甲也不通 175／174。**

## §4 反向判据（只对 ① 那一格：完备性机械尺＋三读数）

- **形状**：仓外 `D:\tmp\wisp\20261026\close3r1\doc\check-close3r1-175ac5.ps1`（判「结题说明写全没有」）。五件要求：① 三枚具名节（`3.1 Q-59`／`3.2 票 174`／`3.3 票 175`）；② 每节至少一条**行级现读**锚；③ **无未完成标记**（词表在尺内，本件全缺）；④ 三态词表命中（`未批`／`已批`／`已落`／`未写`）；⑤ **交叉复算五条硬锚**：doc 里写的 `path:line` 必须在 **HEAD blob** 的同一行命中对应模式——`internal/tools/fs.go:55` ↔ `defaultMaxReadBytes = 256 * 1024`、`docs/PLAN.md:1538` ↔ `尾段读不回`、`internal/risk/provenance.go:94` ↔ `SrcTaskOutput = "task.output"`、`cmd/wisp/run.go:1374` ↔ `loop.RunAsync(ctx, task)`、`internal/tools/task_backfill.go:141` ↔ `b.Roster.Record(taskID, rec)`。任一不符 ⇒ 红。
- **尺本身被实测纠过一次（照实记）**：第一发按**工作树**读锚 ⇒ `as-is rc 1`，红句逐字 `line drift: cmd/wisp/run.go doc=1374 disk=1380`——**不是 doc 错，是 `wire-r1` 正在改那个文件**（worktree 漂到 `:1380`／`:1388`）。改正＝**尺改读 `git show HEAD:<path>`**（doc 的锚本来就该钉 HEAD），**没有为了让绿去改数**（doc 仍写 1374，HEAD 现读仍是 1374）。两发读数都在日志里。
- **三读数**（第二发，日志 `logs/06-doc-check.txt`，全跑在仓外副本上）：**原样绿**＝真件 `rc 0 GREEN`；**变异红×3**＝`M0` 指向不存在路径（摘掉整件）`rc 1`／`M1` 删掉 `票 174` 整节 `rc 1`（红句＝`missing section`，点名 3.2 节）／`M2` 把 `docs/PLAN.md:1538` 改成 `…:999` `rc 1`（`doc missing anchor`）；**还原绿**＝副本还原后 `rc 0 GREEN`。
- **还原核**（时机＝**最后一发变异之后**，逐字在日志 `logs/06-doc-check.txt` 的 `[t0]` 与 `restore proof` 两段）：真件 sha256「开跑前」与「全部变异跑完后」**是同一个值**（真件从未进突变窗口，两值逐字相同）；`sha256(restored)`＝`sha256(pristine)` ⇒ `real unchanged = True`／`restored == pristine = True`。⚠ 本件不写自己最终版的哈希（自指必漂）——它就在日志最后那一发的 `[t0]` 行；`logs/07-commit.txt` 里再核一次「提交后的 blob 哈希＝该值」。
- ⚠ **边界如实说**：这是**完备性**尺，把「结题说明写全没有」从散文变机械；**不是**安全门，替代不了 §5 的桥级五枚用例。

## §5 零回归与四数

- **改前**（`logs/02-pkg-pre.txt`，锚 `385afead`，起手 `internal/tools`＋`internal/risk` 均无脏件）：`go test -count=1 ./internal/tools/ ./internal/risk/` ⇒ `ok …/internal/tools 20.844s`；`internal/risk` **FAIL**（2 枚）。
- **改后**（`logs/02-pkg-post.txt`）：同一命令 ⇒ `ok …/internal/tools 16.518s`；`internal/risk` **FAIL**，**只剩 1 枚**＝`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（另一枚时序 flake 本发绿）。
- **既存红两枚**（HEAD `385afead` 复现，⛔ 本腿责任）：`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（点名 `internal/panel/display_scrub_216.go:218`；台账 `22104`／`22694`／`26002`／`27800` 四处均记「HEAD 复跑同样红」，本腿孤立跑 `rc=1` 0.112s）＋`TestResolvePerCallBudget`（时序；本腿孤立跑 `ok 1.219s` ⇒ 争用，非回归）。
- **175 家族现跑**（`logs/05-175-family-v.txt`）：`go test -count=1 -v -run '175r2|PointerPast256KiB|LongOutputPointerRecoversEveryByte' ./internal/tools/` ⇒ **`RUN=7 PASS=7 FAIL=0`**、`ok 0.204s`；七枚逐名 PASS＝`TestLongOutputPointerRecoversEveryByte`／`TestPointerPast256KiBIsNotFullyReadable`＋五枚 `…175r2`。
- **四数**：改动包＝**0**（0 文件／0 行 Go）｜`git diff --numstat -- cmd/` 与 `-- internal/speech/` 均**空**＝各 **0 字节**｜`internal/{session,config,secret,memory}` 空；提交后 `git show --name-only` 里 `cmd/` 命中＝**0**。

## §6 门禁

- `d22scan`（**三读数，逐条具名**）：
  - 仓外导出树（`git clone --shared` → `D:\tmp\…\export\repo`，**HEAD `385afead` 原样**）：`sh scripts/d22scan.sh` ⇒ **`rc 0` clean**，正控逐字 `PASS=35 FAIL=0 SKIP=0`（`logs/40-d22scan-export-base.txt`）。
  - 导出树 **＋本件**（唯一差异＝这一枚 `.md`）：scan-only `go run . -root <export>` ⇒ **`rc 0` clean**（`logs/43-d22scan-export-withdoc.txt`）。
  - 共享树现况（19:4x）：**`rc 1`**、`3 finding(s)` **全部**点名 `cmd/wisp/resident_voice_engines_windows.go`（`git status`＝`??`，**`wire-r1` 正在写**）：`:736`／`:859` `[emoji] ban #8` ＋ `:1` `[unparseable] …:323:10: rune literal not terminated` ⇒ **tracked 文件零命中、本件零命中**（`logs/44-d22scan-live-now.txt`）。⛔ 本腿不碰 `cmd/**`，未代改。
- 格式两尺（甲＝`gofmt -l`，乙＝`/d/work/base/gopath/bin/gofumpt.exe -l`，范围 `./internal ./cmd ./tools`）：
  - **导出树 HEAD＝甲 `0` 行／乙 `0` 行**（`logs/41-…`／`logs/42-…`）⇒ HEAD 本身是干净的。
  - 共享树＝甲 `6` 行／乙 `6` 行，逐名归属（按 `[\\/]` 判，本腿零 `.go` 改动 ⇒ 无一枚属本腿）：
    - **既存 3 枚**（你派单里那 3 条假红，成因本腿现量＝**worktree CRLF vs blob LF**：`internal\agent\tools.go` CR=225／B=7581、`cmd\wisp\models.go` CR=389／B=17445、`internal\agent\approval\pending_read.go` CR=131／B=6366；导出树同三枚 CR=0 且两尺皆空）＝`internal\agent\approval\pending_read.go`／`internal\agent\tools.go`／`cmd\wisp\models.go`。
    - **`wire-r1` 在写 3 枚**＝`cmd\wisp\resident_voice_engines_windows.go`／`cmd\wisp\resident_voice_windows.go`（`M`）／`cmd\wisp\resident_voice_wire_r1_windows_test.go`（`??`）。
- `gate-clauses.sh`：**未用作本腿判据**（如实说：它在**导出树 HEAD 也 `rc 5`**，BAD＝`G1b`／`G2`／`G5neg`，属票 171／178 那族基线漂移；且其输入含 `.scratch` 状态 ⇒ 不是本腿可归因的尺）。
- 提交：**本文件所在一笔**（主题以 `归属收口腿 close3-r1` 开头）；`git show --name-only` 逐名＝`docs/evidence/s1/close3-r1-impl.md`，`cmd/` 命中 **0**（`logs/07-commit.txt`）。

## §7 顺带交回（⛔ 本腿不做）

- **票 177 AC#5** 是同一张分账的更宽版（逐字要求 `175`／`174`／`Q-59`／`Q-60` 四枚状态）⇒ 本 §3 可直接喂它，缺 `Q-60` 一栏（账 `A349` 已批丙）。
- **票 177 AC#6**（「票 176 不得早于本票合入」）**已过期**：176-r1 于 09-28 已落（`41e8873bd`／`06dae8316`）⇒ 该排程约束的时点已过，归属请编排者裁。
- **票 174 AC#5** 与本票 AC#5 逐字同源（同一根管子两处）⇒ 一 §3 可喂两票。
