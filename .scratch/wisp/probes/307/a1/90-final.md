# 307-a1 · 90-final —— 现态／欠账／冲突处具名报回

腿＝`307-a1`（只读普查）。任务＝票 307 的 `AC#0` 那一格的读数 ＋ 派单加的四路射程。⛔ 碰任何 `AC` 框、⛔ 写"本格闭合"（裁决者≠实现者）。

## 0. 现态

- 进场锚＝HEAD `1eaa91f66502fdf56c33f220d2c5911d6182bbb3`、分支 `dev`（`00-anchor.md`；派单未给快照号 ⇒ ⛔ 冲突）。
- 台面读数：`git status --porcelain -- internal cmd`＝**0 行**；`frontend/dist/index.html`＝**1044 字节**、mtime `2026-10-10 08:51:10 +0800`（`stat` 到纳秒＝`...08:51:10.506316400 +0800`）。⇒ 〔现跑〕**`internal`/`cmd` 工作树＝HEAD blob**，所以本腿报的每一个枚数在"工作树"与"HEAD blob"两把尺下同值（这一句只在 `internal`/`cmd` 射程内成立，`.scratch/**` 不保证）。
- 笔数／写面／越界：见本件 §3。

## 1. 交了什么（四路射程 → 件）

| 件 | 内容 |
|---|---|
| `10-fold-roster.md` | `AC#0` 的折法名册：`home` 侧 2 枚赋值点、`e.canon` 侧 3 枚赋值点、`cand` 侧折法、`roots` 上游链、票面行号腐烂对照 |
| `20-landing-shapes.md` | 三形候选（A 只折 home／B 两侧都过 `Resolve`／C 复用 `formsOf`·`anchorForms`）＋一枚登记用的反形，逐形写动哪几行、变严变松、撞⛔ 撞 `AC#1` 四条 |
| `30-test-nails.md` | 会碰坏的既有钉名册（A 组＝票面点名 5 枚、B 组＝同族未点名、C 组＝13 枚负向钉、D 组＝`SyncDetectionComplete(` 调用者＝0 的三义）＋需要编排者车道跑的具名发数 |
| `40-machine-alias.md` | ★必答的台面现量与具名回答 |
| `logs/*` | 每一发的原始 stdout，逐把尺自落 `rc_*=`；两枚 `.ps1` 脚本本体 |

### `AC#0` 那一问的读数（一句话版，细节与尺在 `10-fold-roster.md`）

**两枚折法不是同一把**：`e.canon` 的 3 枚赋值点里 2 枚过 `Resolve`（`syncdirs.go:152`、`:156`，其中 `:156` 必为 handle 长形），`s.home` 的 2 枚赋值点（`syncdirs.go:88`、`provenance.go:345`）**0 枚过 `Resolve`**（只有 `normPath`）。
但 `e.canon` 另有 1 枚（`:160` `default:` 支，`normPath(expandInput(r.Path))`）**连 `Resolve` 都没过、也没过 `lexCanonical`** ⇒〔读码推到〕"一边长一边短"这枚机制在 `:422`（票面写作 `:425`）也成立，触发条件＝该根盘上不存在（POSIX 上恒成立，因 `pathresolver_other.go:10` 的 `DEFERRED` stub）。
⇒ 出口枚数读数＝**2 处**（`:422` 的 `:160` 支 ＋ `:428`），其中**生产可达＝1 处**（`roots` 在生产为空：两个产码 `NewProvenance` 都是 `ProvOptions{NoProbe: true}`、⛔ `SyncRoots`/`SyncFixturePath` ⇒ `finalize()` 零循环 ⇒ `complete` 恒 false ⇒ `:422` 循环零次迭代）。
⇒ `home` 那枚折法的两枚同形赋值点里，**生产走的是 `provenance.go:345`**，而票面点名的 5 枚用例**全部**也走 `:345`（无一枚走 `syncdirs.go:88`；唯一吃 `:88` 的 fold 断言是 `syncdirs_test.go:501`）。落点形状请据此裁——**本腿⛔ 选形**。

### ★必答（`40-machine-alias.md`）

这台机器 TEMP 全链逐段 SHORT==LONG（`C:\Users`/`swq`/`AppData`/`Local`/`Temp` 五段无一枚别名），`USERPROFILE=C:\Users\swq` 亦 SHORT==LONG；而 `C:\Users\swq\tmp\riskattrib1\RunnerTempAreaForAttribution1` SHORT!=LONG（`RISKAT~1\RUNNER~1`，与票面逐字同）；`%TEMP%` 下新建长名目录**取得到**别名（`WI0253~1`、`LONGSE~1`）⇒ 短名生成对这枚卷是开着的，本机绿只因"账户名恰好 ≤8 字符"。
**具名回答：会**——换一台 profile 无短名别名的机器，那 5 枚**静悄悄全绿**（⛔ SKIP、⛔ 失败消息、⛔ 任何自陈），因为触发条件在断言之前就不存在。**⇒ 本机绿⛔ 能当凭据。**

## 2. 我⛔ 做到的（逐条具名，⛔ 写"没跑／没做"就当没这回事）

1. ⛔ 跑任何 Go 编译面：**`go build`／`go vet`／`go test`／`gofmt -l`／`d22scan` 一发都⛔ 跑**；连 `go env`／`go list` 也⛔ 跑过 ⇒ **本腿的 Go 面发数＝0**（派单允许 `go env`/`go list`，我选择一枚都⛔ 用，因为本格的读数全在静态名册＋台面 PS 读数里）。需要颜色的两发已具名写在 `30-test-nails.md` §E，交回编排者车道。
2. ⛔ 选定落点形状、⛔ 判某形可行／⛔ 撞约束（`20-landing-shapes.md` 每形只给"动哪几行＋后果＋对四条硬约束的边缘位"）。
3. ⛔ 造/提任何修复；⛔ 动一字节产码或 `_test.go`。
4. ⛔ 回答 `AC#3` 的修法（夹具落点、该不该进仪器）——只交到"这台机器给不给得出短名世界 ＋ 包内已有两枚 OS 级取短名手（且都带 `t.Skip` 出口）"这一读数层。
5. ⛔ 碰 `AC#0`…`AC#5` 任何一枚框、⛔ 写"本格已闭"。
6. 欠读数（具名）：① `fsutil 8dot3name query C:` 取不到（`Error 5: Access is denied`，`fsutil_volume_rc=1`，本腿⛔ 提权）⇒ **逐卷短名状态没量**；② `SyncRoots()` 的消费面⛔ 复核（`syncdirs.go:181-182` 注释声称面板 Security page 消费本判定，我只量了 `SyncDetectionComplete(`/`IsSyncPath(` 两枚 API＝0 跨面调用者）；③ B 组/C 组那批用例的**颜色**一枚都⛔ 跑过（全靠读码推到，已逐条标〔读码推到〕）。
7. ⛔ 删任何临时件：两枚探针目录（`%TEMP%\wisp307a1AliasWorldProbe\...`）与票面那条链上的 `riskattrib1` 两枚目录**原样留盘**；本腿的临时件从一开始就全在自家 `logs/`，仓库根 0 枚新增。

## 3. 越界自证尺（逐笔，⛔ 区间尺）

见本件末段"逐笔读数"。尺＝对每一枚自己的 commit 跑 `git show --name-only --format=%H <sha>`，名册每一行须落在 `.scratch/wisp/probes/307/a1/**` 内；新建 `.go`＝0 枚、`.out`＝0 枚、0 字节件＝0 枚（尺＝`find .scratch/wisp/probes/307/a1 -type f -size 0`，**带 `-type f`**）。⛔ push。

## 4. 与派单／票面**冲突处**具名报回（4 枚）

1. **票面 `AC#0`/现量节的产码行号腐烂**：票面写 `syncdirs.go:425` 与 `:429`，盘上现量＝**422 / 428**（内容锚见 `10-fold-roster.md` §0）。同段的 `:88`/`:197`/`blacklist.go:133`/`:169`/`provenance.go:928`/`:345` 以及测试侧 `:162/:213/:233/:355/:402` **全部同号未腐烂**。⇒ 派单转述的"`:425` 那一支"我按**内容**（`isUnder(cand, e.canon)`）理解，⛔ 按号理解。
2. **派单说"同形复制件 `internal/risk/provenance.go:345` 附近那一族"**——现量结果比"附近那一族"更硬：那一枚⛔ 只是复制件，它是**生产面唯一那枚 `home` 赋值点**（两个产码 `NewProvenance` 都传 `NoProbe: true`），而票面点名的 5 枚用例也全走它。⇒ 这一枚要不要算进落点不是"顺手覆盖"的问题，是"⛔ 覆盖它就没有修到生产"的问题。读数交回，⛔ 由我裁。
3. **派单转述的第三形名字对不上**：派单写"复用 `pathresolver.go` 里既有的 `formsOf`·`anchorForms` 那一族"。现跑：`formsOf` 定义在 `pathresolver.go:345`、`anchorForms` 在 `:391`（⛔ 腐烂，与产码那两枚无关），但**实际消费它们做 `home` 锚点比较的那一族住在 `blacklist.go`**（`:199`、`:214`、`:221` 等 8 枚 `anchorForms(...)` 调用 ＋ `aUnder`/`aEq`/`aRelates`/`uncertainAnchorMiss` 四枚关系函数，`blacklist.go:262-322`）。⇒ 落点若选这形，**改动点在 `syncdirs.go:428` 调 `blacklist.go` 的关系函数**，⛔ 落在 `pathresolver.go` 里。
4. **票面 §1"乙族那 5 枚"里 `:162` 那一枚的归因**：票面把 5 枚都算作"home 单边折法"的出口（`AC#0` 自己也没预设）。现跑读数：`:162` 所在的 fixture 根 `home\Cloud\NutstoreSync` 从未被创建（`sandbox()` 只建 `home`/`local`/`roaming`，`syncdirs_test.go:15-27/152-163`）⇒ 它同时吃 `:422` 的 `:160` 支。若落点只折 `home`（形 A），〔读码推到〕`:162` 能不能复绿取决于 fallback 是否接住它——**这一枚是本腿唯一"读数可能与票面归因不同"的地方**，具名报回，⛔ 由我改判任何归因。

## 5. 逐笔越界尺（现跑，⛔ 区间尺）

尺＝逐笔 `git show --name-only --format= <sha>`（第 2 笔那批循环输出的原始留档＝`logs/52-commit2.txt`）。

| 笔 | SHA | 名册行数 | 越出 `.scratch/wisp/probes/307/a1/**` 的格 |
|---|---|---|---|
| 1 | `08cc614ec1915c76e385ab9e070645512f5861e0` | 4 | **0** |
| 2 | `4d4c10d4f03cfd29def0713b42dc5bc9078e8445` | 25 | **0** |
| 3 | 本笔（`90-final.md` ＋ `logs/52-commit2.txt`） | 2 | **0** |

- 新建 `.go`＝**0** 枚、`.out`＝**0** 枚、0 字节件（尺＝`find .scratch/wisp/probes/307/a1 -type f -size 0`，**带 `-type f`**）＝**0** 枚（`logs/51-zero-fix.txt` 现跑：`count_zero=0`、`count_files=27`）。
- ⚠ **自报一枚自己的缺陷**：第 1 笔带了一枚 0 字节件 `logs/00-anchor-readings.log.2`（那是 `git status --porcelain -- internal cmd` 的**空输出**，读数本身有效＝0 行）。第 2 笔已往该件补了"empty output by design ＋尺名＋rc"说明行 ⇒ **HEAD blob 不再 0 字节**，但**历史那一笔的 blob 确实 0 字节**。⛔ 删（临时件只建不删）、⛔ 改写历史（⛔ `--amend`/`rebase`），要更正就追加——已按"追加"更正。同形的另一枚＝`logs/g-k-crossface.txt.2`（第 2 笔一并补好，故第 2 笔起 count_zero=0）。
- ⚠ **共享工作树的并发读数（⛔ 本腿的越界）**：`git rev-list --reverse 1eaa91f..HEAD` 里**夹着别家的笔**（`probes/303/v2/**`、`probes/306/orch/**`、`issues/305`/`306`/`307`/`72` 四枚票面、`docs/reports/HANDOVER.md`、`docs/reports/pending-and-issues.md`、`docs/evidence/s1/306-*.md`、`docs/evidence/s1/303-ac3b-ac5-v2.md`）。⇒ 这正是派单"终局自证尺⛔ 用区间尺"的原因：**拿 `1eaa91f..HEAD` 那把区间尺会数出上百行别家路径，看着像本腿越界**。⇒ 逐笔尺才是判据。
- 台面另有一批别家未跟踪/删除态（`design/**` 16 行 ` D`、`?? -`、`?? .scratch/.scratch/`、`.scratch/commit-msg-*.txt` 等，`logs/50-selfcheck.txt` ruler5 段，尺＝`git status --porcelain | grep -v "^ M" | head -30`）⇒ 本腿**一枚都⛔ 碰、⛔ 删、⛔ 收进自家 commit**（`git add` 只带自家目录那一枚 pathspec，⛔ `-A`/`.`）。
- push＝**0**（本腿全程只 commit 三笔）。
