# 270-a2 — 票 270「一票三形」配对读数（只读腿；⛔ 不选形、不翻框、不改产码一字）

射程＝**配对读数**：复跑 `270-a1` 的三把锚定尺并逐枚对数；判"一票三形装不装得住"；判 AC#1 那枚判据今天**会不会恒真**；把"零调用者／唯一产出者"那类句子的尺锚回调用形状并逐跳追到 main。⛔ 本件不重做普查（代价表在 `a1/census.md`），也不引用 `a1/census.md` 的**判语**当凭据——它的数我一枚枚复跑，它的 194 那枚我量到 153，差在口径，见 §1 第 4 行。

---

## §0 起手锚

| 项 | 命令（逐字） | 读数 |
|---|---|---|
| HEAD（起手） | `git rev-parse HEAD` ／ `git rev-parse --short HEAD` ／ `git branch --show-current` | `82a10da4d4398dc5de2b66c5f878a8ce888ff493` ／ `82a10da4` ／ `dev` |
| HEAD（交件前复验） | 同上 | **`7ec2f518de4c13c0fbf345675ccdad54f21020c4`（本腿期间 HEAD 动过两枚：`82a10da4`→`9aef5939`→`7ec2f518`）** ⇒ 下面所有行号在两个时刻都复验过，未漂 |
| 产码面工作树 | `git status --porcelain -- cmd internal \| wc -l` | **0 行**（起手＝交件前）⇒ 两枚在飞写腿（`272-r2`／`111-r5`）此刻在 `cmd`／`internal`／`.github` 上**都还是 0 行**（`git status --porcelain -- .github \| wc -l` ＝ 0），本腿读到的产码面＝干净面 |
| 三枚本体文件对 HEAD 的字节自证 | `md5sum < f` vs `git show HEAD:f \| md5sum` | `internal/tools/subagent_197.go`＝`06caf8f5465ff1c47c310da27fe3ffab` 工＝HEAD；`cmd/wisp/run.go`＝`1d95cfaf116d9092514ee0a905a82b23` 工＝HEAD；`internal/tools/failclosed_236_teeth_test.go`＝`c793966c2730c64896d9512de597ddf4` 工＝HEAD ⇒ 本腿读数落在与票面 §现量同一枚盘上 |
| 锚点复验（派单五枚＋`:293`） | `awk 'NR==135\|\|NR==256\|\|NR==258\|\|NR==278\|\|NR==293'` ＋ `sed -n '789p' cmd/wisp/run.go` | 逐字对上：`:135` `BaseOptions func() (agent.Options, bool)`・`:256` `if t.d.BaseOptions == nil \|\| t.d.ParentTools == nil {`・`:258` `}`・`:278` `base, wired := t.d.BaseOptions()`・`:293` `opt.Tools = newSubagentToolProvider(t.d.ParentTools)`・`run.go:789` `BaseOptions: func() (agent.Options, bool) { return rt.loopOpt, rt.loopOptSet },` ⇒ **一枚未漂**（与 `a1` §1.1 同判，但我自己量的） |
| 本腿跑过的 `go` 命令 | 全文＝**1 发**：`go test ./internal/tools/ -run 'Test236R2TaskSpawnRefusesWhenBaseOptionsIsUnwired\|Test236R2TaskSpawnRefusesWhenParentToolsIsUnwired' -count=1 -v` | rc=0／2 枚 `--- PASS`（`logs/narrow-baseline.txt`，7 行）⇒ 单包·具名 2 枚·`-count=1`，⛔ 未跑 `./cmd/wisp`、⛔ 未跑整包、⛔ 未跑任何突变 |
| 落盘日志 | `.scratch/wisp/probes/270/a2/logs/` | `r1-callshape.txt`・`r2-nilcheck-eq.txt`・`r3-nilcheck-ne.txt`・`r4-bare.txt`・`a1-goonly-paths.txt`・`a2-scratch-goonly.txt`・`a1-tracked-goonly.txt`・`a1-files.txt`・`a2-files.txt`・`k-literals.txt`・`mut2.txt`・`mut-readings.txt`・`narrow-baseline.txt`・`baseline.txt` |

---

## §1 三把尺复跑对照表（⛔ 不抄 a1 的判语，只与它的数对）

| # | 尺 | 命令（逐字，我在仓根跑） | `270-a1` 报的 | 我量到的 | 对上？ |
|---|---|---|---|---|---|
| 1 | 调用形尺 | `grep -rn "\.BaseOptions(" --include=*.go .` | 19 命中＝产码真调用 1（`subagent_197.go:278`）＋ `.scratch` 副本 17 ＋ 测试注释 1（`failclosed_236_teeth_test.go:203`） | **19 行**＝`./.scratch` 17 ＋ tracked 2，且 tracked 那 2 枚逐枚＝`internal/tools/failclosed_236_teeth_test.go:203`・`internal/tools/subagent_197.go:278`；17 枚副本的第一枚仍是嵌套重复目录 `.scratch/.scratch/wisp/probes/222/v1/mutations/asis/subagent_197.go:269` | **对上**（19／17／tracked 2 枚的行号一字不差） |
| 2 | 判空尺 | `grep -rn "BaseOptions == nil" --include=*.go .` | 19 命中＝产码判空 1（`:256`）＋ 副本 17 ＋ 注释 1（`failclosed…:12`） | **19 行**＝副本 17 ＋ tracked 2＝`failclosed…:12`・`subagent_197.go:256` | **对上** |
| 3 | 正向判空尺 | `grep -rn "BaseOptions != nil" --include=*.go .` | 0 命中（无第二支正向判空） | **0 行** | **对上** |
| 4 | 裸符号尺（对照用，不是判据尺） | `grep -rn "BaseOptions" --include=*.go .` | **194 命中，其中 172 落在 `.scratch`** | **153 行＝`.scratch` 131 ＋ tracked 22**（tracked 逐文件＝`subagent_197.go` 6・`run.go` 2・`failclosed…_test.go` 10・`subagent_197_test.go` 1・`subagent_222_test.go` 1・`task_cancel_221_legs_test.go` 1・`cmd/wisp/subagent_selfapproval_197_test.go` 1 ＝ 8 产码＋14 测试） | **对不上：它报了 194／172，我量到 153／131，差 41 枚全在 `.scratch` 那侧。** 差在口径＝它那把尺**没带 `--include=*.go`**：按路径后缀清它落盘的 `logs/baseoptions-all.txt`（194 行）＝`.go` 路径 153 行＋非 `.go`（工单与 probes／docs 下的 `.md`，含票 270 与票 236 自己）41 行；`.go` 那 153 行与我**逐枚对上**（副本 131／tracked 22、tracked 逐文件枚数与它 §1.2 的 8＋14 一致）。⇒ 它写错了命令串，不是量错了盘。⚠ 引用 a1 的裸符号数时按 **153／131** 用 |
| 5 | 副本集合增删 | `diff <(a1 logs 里 .scratch 下 .go 路径) <(我这把尺的 .scratch 下 .go 路径)` | （a1 未给集合级读数） | **diff 无输出＝0 增 0 删** ⇒ 131 枚不可编译副本自 `a1` 起没动，第 4 行的差不是"盘变了"造成的 | 新读（口径差的排除证明） |
| 6 | 名册尺（配套计数） | `grep -rn "SubagentDeps{" --include=*.go cmd internal`／`grep -rn "BuiltinSubagentEntries" --include=*.go cmd internal` | 7 处／8 命中 | **7 行／8 行**，分布与 a1 一致（测试里两枚不派发：`task_cancel_221_test.go:51` 只枚名、`subagent_selfapproval_197_test.go:551` 只做 `reflect.TypeOf` 类型遍历） | **对上** |
| 7 | 产出者尺（调用形＋字面量形两把） | `grep -rn "BaseOptions:" --include=*.go cmd internal`（去 `_test.go` 后） | 「生产侧唯一产出者＝`run.go:789`」 | **4 处字面量＝产码 1（`cmd/wisp/run.go:789`）＋ 测试 3**（`subagent_197_test.go`・`subagent_222_test.go`・`task_cancel_221_legs_test.go`）⇒ 产码侧恰 1 枚 | **对上** |
| 8 | 文案等值钉尺 | `grep -rn "lit236SpawnNoAssembly" --include=*.go internal cmd` | 乙顶 3 枚（`:95-96` 共用常量＋`:195`＋`:218`） | **5 行＝`:95` 声明・`:195` 等值・`:197` 红句・`:218` 等值・`:221` 红句**。⇒ "顶红"那 3 枚（声明＋两枚等值）成立；另 2 行是红句文案本体（改串时也要同批改） | **对上，并补 2 枚**（`:197`／`:221` 两句红话里写着"subagent_197.go:256 那一支"，乙拆支后这两句也会说谎） |
| 9 | 我没复跑的那把 | — | 「`internal/tools` 25 处＋`cmd` 1 处无 recover 的 spawn 派发点」（a1 §1.3／§2.3） | **未复跑**：a1 件内只给了位置∶枚数，没逐字给那把尺的命令，本腿不猜它的口径（`h.spawn(`／`.spawn(t,`／`Execute(` 三种口径数不同） | ⇒ 记 §4 第 3 条 |
| 10 | "读盘尺存不存在"这把（来自转述，不是来自 a1 件内） | `grep -rln "go/parser\|go/ast" --include=*.go internal cmd tools scripts` ＋ `sed -n '193,212p' internal/tools/paths_twocontainments_252_r2_test.go` ＋ `sed -n '224p' internal/tools/task_cancel_221_legs_test.go` ＋ `grep -n "读盘\|embed\|ReadFile\|go/ast\|parser" .scratch/wisp/probes/270/a1/census.md` | 转述称"没有任何读盘／embed 的尺去读 `subagent_197.go`"；a1 件内**无此句**（那把关键词尺只中 `:121`・`:203` 两行，别有所指） | **22 枚文件用 `go/ast`／`go/parser`**；`paths_twocontainments_252_r2_test.go:193-212` 用 `os.ReadDir(".")`＋`parser.ParseFile(..., ParseComments)` 逐枚解析本包全部非 `_test.go` 源文件 ⇒ **`subagent_197.go` 今天已被打开**；`task_cancel_221_legs_test.go:224`＝`os.ReadFile("task.go")`＝同包"读产码源文件钉常驻判据"的现成先例 | **对不上（错在转述侧，不在 a1 侧）**；修正与它带来的约束见 §5 第 1 条 |

**小结：三把锚定尺（#1・#2・#3）3/3 逐枚对上；配套计数（#6・#7・#8）也逐枚对上。对不上 2 枚：#4 裸符号尺差在 a1 少写 `--include=*.go`（口径，不是盘），#10 那枚"没有读盘尺"是**转述造出来的**、a1 件内没有这句话，且它在今天的盘上是假的。**

---

## §2 配对结论：先钉哪枚＋它今天会不会恒真

**落点共用情况（现量）**：三形的产码落点是**同一枚文件的同一枚函数体**——`internal/tools/subagent_197.go`（甲在 `:278` 前插判空；乙把 `:256-258` 拆两支＋同批装 `:278`；丙在 `:256`／`:278` 处留注释或断言性拒收）。只有乙额外必须动第二枚文件 `internal/tools/failclosed_236_teeth_test.go`（`:95` 常量拆分＋`:195`／`:218` want＋`:197`／`:221` 红句）。

★**票面 AC#1 那条完成判据今天就是恒真的，不许当凭据。**
- 判据原文＝"删除 `:256` 整支后的台件读数：`internal/tools` 整包仍出全量名册（PASS+FAIL ≈ 206，⛔ 不许出现 25／0 那种失明形）"。
- **判据换成反形还绿吗？反形正是它自己要读的那一发**＝"产码一字不装、只摘 `:256`"。引用件读数（⛔ 不是我跑的，本腿直接开在 `236-v1` 原文上核过，不转抄 a1 的转述）：`.scratch/wisp/probes/236/v1/verdict.md:153` 记 `A-4-spawn-assembly-false`（`:256` → `if false {`＝两半一起没、**保留** recover）＝**rc=1／PASS=204／FAIL=2**，红名册逐枚＝`…WhenParentToolsIsUnwired`＋`…WhenBaseOptionsIsUnwired` ⇒ **合计 206＝名册全**；同件 `:164` 记 `A-4-without-recover`（同一产码突变＋overlay 摘掉测试里那 5 行 recover）＝**25／2＋`^panic` 2 行＋其后约 181 枚零读数**；`.scratch/wisp/probes/236/r2/evidence.md:112` 记 m-1d＝**rc=1／PASS=202／FAIL=2**、`:213` 记"摘整支＝2 枚红"。⇒ 名册之所以还全，是 `spawnGuarded`（`failclosed_236_teeth_test.go:230-238`）那层 recover 替我们把"崩"降级成了 1 枚红。**"零读数消失"这件事今天已经成立，与有没有装甲无关 ⇒ 该尺对 `:278` 装没装甲不敏感。**
- 我现量的旁证：`grep -rc "^func Test" internal/tools/*_test.go` 求和＝**206**（静态顶层枚数，不能当台件读数）⇒ 票面那个"≈206"这个数今天仍是真的，但它真的是名册大小，不是形状。
- ★**有牙的第一枚**＝**同一发突变里的具名用例翻色读数**：`Test236R2TaskSpawnRefusesWhenBaseOptionsIsUnwired` 在"摘掉那一支 guard"的那一发必须从**红**翻**绿**（等值钉在 `failclosed_236_teeth_test.go:218`）。今天摘支后它是**红**（引用件 `evidence.md:117`／`:139` 的 got 逐字以 `panic（guard 摘掉后不是拒绝而是崩）: runtime error: invalid memory address or nil pointer dereference` 开头），装上 `:278` 判空且拒收串与 `lit236SpawnNoAssembly` 同串才绿 ⇒ **敏感、非恒真 ⇒ 先钉这枚**。
  - ⚠ 这枚钉的前提要写死：**`:278` 的 nil 支必须复用 `:257` 那句串**。若甲／乙选一句新串 ⇒ 摘支后 `:218` 仍红 ⇒ 那就不许拿"旧钉仍红"读成"形状破"，得新造一枚具名用例（这一格属编排者裁，a1 §3 第 4 条已在问）。
  - ⚠ 摘哪一支按所选形定：甲形＝摘 `:256` 整支（两支都不留）；乙形＝摘新拆出来的 `BaseOptions` 那一支（`ParentTools` 那支留着，否则 `:293` 那一维会混进读数）。
- **基线我量了**：`go test ./internal/tools/ -run '<那两枚具名用例>' -count=1 -v` ＝ rc=0／两枚 PASS ⇒ "今天两枚等值钉都绿"是本腿实跑；"摘支后剩什么"全是引用件，本腿禁跑突变 ⇒ §4 第 1 条。
- **另两枚钉的敏感性现量（结构性，不是跑的）**：`:195` 与 `:218` **共用 `:95` 那一枚常量**（`lit236SpawnNoAssembly`，串文里同时写着 `BaseOptions/ParentTools`）⇒ 两枚钉对"是哪半边拦的"结构性失明 ⇒ 甲不使任一枚翻色、乙若两支复用同串也不使任一枚翻色 ⇒ 它们只能当**文案等值凭据**，不能当"nil 防护已就位"凭据。

**一句话结论**：本票的第一枚钉应落在"定向突变＋具名用例翻色（红→绿）＋`^panic` 行数＝0"那一发上，⛔ 不落在"整包名册数"上；后者今天零改动也绿，按其自身判据即被淘汰。

---

## §3 一票三形装不装得住

**判定：装得住"崩溃那一维"，装不住"三形同尺"这句——票面标题的一票两形／a1 补的一票三形里混了三种不同射程的东西，其中第三形按现量不是崩溃形。**

理由四条（都带现量或指名的引用件）：

1. **同落点、互不打红**＝甲 0 枚红／丙 0 枚红／乙 3 枚红全属它自己要改的那批 `236` 钉（`:95`／`:195`／`:218`），不是别的形替它背的。⇒ 三形塞同一枚票不会因为落点冲突而互相洗读数。
2. **第三形（`:293` 无条件使用 `ParentTools`）今天不是崩溃形**＝静态读 `subagent_197.go:546-568`：`newSubagentToolProvider(nil)` 只造 `&subagentToolProvider{inner: nil}`，**不当场崩**；`Tools()` 在 `:553` 有 `if p.inner == nil { return nil, nil }` 守卫；真正会崩的是 `:578` 的 `p.inner.Execute(ctx, req)`，而那一步只有子环路真去调一枚工具才走得到。⇒ 摘掉 `ParentTools` 那半边的形状＝"真派生＋空工具目录"，引用件印证：`v-m1d1`／m-1d1 的 got 逐字是"子代理 …（父工具面没接线）已结束，状态 Settling"（`evidence.md:115`／`:130`，rc=1／PASS=203／FAIL=1＝名册仍全）。⇒ **a1 §4 第 6 条那格（"会不会自己崩，需实跑定性"）本腿能按静态读给一半答复：`:293` 不是当场崩形。**
3. ⇒ 连带结论：**票面 §26 那句判据（"拆掉任一半 guard 之后整包不许出现零读数"）只对 `BaseOptions` 半边有意义**；`ParentTools` 半边今天产不出零读数（现量静态＋引用件读数两路同向）。若把这句按字面当三形共用的验收尺，会出现"乙拆支后 `ParentTools` 那半边照出全量名册 ⇒ 判据自动满足"的空转。⇒ **这句该收窄成"摘 `BaseOptions` 半边／整支 ⇒ 不许 `^panic`、不许零读数"**，且与 §2 的具名翻色钉配成一对用。
4. **真正需要同批的只有乙**（票面禁区第 3 条"不许把 `:278` 的修法与拆语句分批"）：本腿复认其成立——拆完的那一形正是 `:278` 裸奔的形；甲／丙不新增这条约束。

**拆票建议（⛔ 本腿不建票、不动票面）**：建议拆成 **2 枚**＋1 枚待议：
- 票 270-A「`:278` 装甲＋AC#1 的敏感读数」＝甲，判据＝§2 那枚具名翻色；AC#5 那句注释可同批（丙的另一半）。改动面 1 文件、0 枚红。
- ★丙**不是**"没法钉"那一形：本腿现量同包已有源形状尺先例（`paths_twocontainments_252_r2_test.go:193-212` 已把 `subagent_197.go` 解析在内、`task_cancel_221_legs_test.go:224` 已用 `os.ReadFile` 钉常驻判据）⇒ 丙的"具名常驻判据"零新仪具成本，写成一枚源形状尺即可钉住"后人别把装甲拆回去"；⛔ 但它对 `-overlay` 突变结构性失明（`failclosed_236_teeth_test.go:68-75` 记着这条仪器事实），**不能顶替 AC#1 的突变读数**，两枚尺各管一事。
- 票 270-B「`:256` 拆两支＋`236` 那批钉同批改串」＝乙，判据＝AC#2（半边摘各恰 1 枚红）＋AC#3（三句互为替身的防护）＋新增的"红句文案也得改"（§1 第 8 行补的 `:197`／`:221`）；这一枚要过"换串不是放宽"的人工批准口径（a1 §3 第 5 条在问的就是它）。
- 待议第 3 枚：`ParentTools` 那半边"摘了会真派生一枚空目录子代理"算不算诚实形状＝**行为题不是崩溃题**，别塞进本票判据句。
- ⇒ **一票三形装得住文件、装不住判据句**；不拆也能做，代价＝AC#1 与 AC#2 分属两枚互不相干的读数口径，验收时容易把名册读数（§2 已证恒真）误当装甲凭据。

---

## §4 判不动的地方（缺什么料写清楚，⛔ 不推测填空）

1. **AC#1 与 AC#2 的全部台件读数**——本腿禁跑突变（本波只有 `272-r2` 允许 Go 突变）。缺：overlay 摘 `:256`（甲形整支／乙形单支）后那一发的 `^--- (PASS|FAIL)` 名册、`^panic` 行数、以及"具名用例红→绿"的翻色证明。⇒ 必须另派一枚能跑突变的读数腿。
2. **名册今日全量基线**（rc=0 时 `PASS=?`）——本腿被禁跑整包 `-count=1`。我只给了静态顶层 `^func Test` 求和＝206 当**数口径的旁证**；引用件两发（v1 的 206／r2 的 204）分属两个时点 ⇒ "≈206"这个阈值该按今天盘重定一次，并由**能跑整包的腿**先声明口径（顶层 vs 含缩进子测试；`236-v1` 自己声明过 `grep -c '^--- PASS'` 行首锚定看不见子测试）。
3. **a1 的"25＋1 处无 recover 的 spawn 派发点"那把尺**——它没在件内逐字给命令，口径不唯一（`h.spawn(`／`.spawn(t,`／`Execute(` 三种数法），本腿不猜 ⇒ 未复跑。要配对得先让 a1 那把尺把命令写进件，或现采一种口径重采。
4. **"子环路带着空工具目录会不会真去调 `Execute`（→ `:578` 延迟崩）"**——要么读 `internal/agent` 环路对"目录外工具名"的处理，要么实跑。本腿射程外（实跑属突变腿；环路读数不在票 270 面），⇒ 只登记，不定性。
5. **`ErrorClass` 那一维**（票面禁区第 4 条〔契约邻接＝D37〕需人工批准）——不是只读腿能判的。
6. **⚠ 登记两条今天已过期的语句（⛔ 本腿未改票面一字，也未改盘上注释一字）**：
   - 票面 §现量 第 2 行"…**无条件调用**，就在 guard 之后**两行**"（与人话段"两行之后"同源）。现量＝`:258`（guard 收尾）→ `:278` 相隔 **20 行**；`:256` → `:278` 相隔 **22 行**；中间还夹着 `:259-263` 父任务 id 支、深度上限支、池满支。⇒ "两行之后"在今天的盘上是假的，且它是立案语义（"顺带挡住下一行"）里最直觉的那一句，失真代价高。
   - 同一句过期话还写进了盘上注释：`internal/tools/failclosed_236_teeth_test.go:201-203`（"…**two lines later**"）。⇒ 若最终选丙（在注释里留常驻判据），这两处过期文案要一起处置；谁处置、算不算"动 236 的凭据"，属编排者裁（a1 §3 第 6 条已在问）。
7. **引用 a1 时的两处口径修正**（不算它的错，属本腿配对产物）：裸符号数按 **153／131** 读（§1 第 4 行）；"乙顶 3 枚"要读成"3 枚红＋2 句红话要同批改"（§1 第 8 行）。

---

## §5 与我收到的转述不符之处（具名）

1. 转述说"`270-a1` 的 `census.md` 里写甲形顶 0 枚钉、乙形顶 3 枚（`:95-96` 共用常量、`:195`、`:218`）、丙形 0 枚但没有尺能钉住它"——**枚数我全部复对上**（0／3／0），但**"没有任何读盘／embed 的尺去读 `subagent_197.go`"这半句是转述造出来的，`a1/census.md` 里没有这句话**（我按"这句话在说什么"扫过它全文，`读盘|embed|ReadFile|go/ast|parser` 只命中两处：`:121` 与 `:203`，两处都不是那个意思）。而且**这半句在今天的盘上是假的**，本腿现量三枚反例：
   - `internal/tools/paths_twocontainments_252_r2_test.go:193-212`＝`os.ReadDir(".")` ＋ 对每一枚非 `_test.go` 的 `.go` 跑 `parser.ParseFile(fset, name, nil, parser.ParseComments)` ⇒ **`subagent_197.go` 每天都在被这枚尺打开**（只是它断言的是路径包含关系，不是 `:278`）。
   - `internal/tools/task_cancel_221_legs_test.go:224`＝`src, err := os.ReadFile("task.go")`＝**同包已有"读产码源文件来钉一条常驻判据"的先例**（它钉的是 `DEFERRED` 标记）。
   - `internal/tools/task_state_188_test.go:253-261`＝`filepath.WalkDir` ＋ `parser.ParseFile`，且 `:247-249` 写着分母纪律（"a scan whose denominator is zero reports 'clean' for a tree it never read"）；`cmd/wisp/leg_dispatch_gate_133_test.go:375`／`leg_sink_gate_131_test.go:581`／`panel_geometry_255_test.go:210` 同族 ⇒ 全仓 22 枚文件用 `go/ast`|`go/parser`。
   ⇒ 修正后的结论：**丙那条"具名常驻判据"零新仪具成本**（同包已有源形状尺技术＋分母纪律先例），代价只是新写一枚断言；⛔ 但它**不能兼任 AC#1 的突变读数**——盘上注释 `failclosed_236_teeth_test.go:68-75` 记着一条仪器事实（本腿亲验原文，非转抄）：**读盘型尺对 `go test -overlay` 换进去的字节结构性失明**（编译器看的是 overlay 字节，`os.ReadFile` 看的是盘上原字节）。⇒ 两枚尺各管一事：源形状尺钉"形状不许被后人拆回去"，突变台件读数钉"装甲真的在拦崩"。
   另：`scripts/` 只有 `d22scan.sh`／`slo-*.sh`／`portable-tests*.sh`／`wisp-cli-tests.sh`／`winsec-tests.sh`／`check-path-length-budget.sh`／`build.ps1` 等，`tools/` 只有 `d22scan`／`mockllm`／`signmodels` ⇒ **仓里没有常驻突变执行器**这一条我复认（"把突变读数搬进 CI"那一支要新建仪具＋人工点头）。
2. 转述把 `:256` 那一枚描述为"同时是 `:278` 唯一的 nil 保护"——**成立**（本腿现量：`grep -n "t\.d\.[A-Z][A-Za-z]*(" internal/tools/subagent_197.go` ＝**只 1 命中＝`:278`**，⇒ 全文件只有这一枚"函数值字段被无条件调用"的崩溃点；`BaseOptions != nil` 尺＝0 ⇒ 无第二支正向判空）。
3. 转述给的"生产者只有一处 `cmd/wisp/run.go:789`"——**成立**，但要把"生产者"和"调用者"分清（本轮射程第 3 条要求的正是这个）：调用形尺在 tracked 侧给的是 `subagent_197.go:278`（唯一产码调用）＋`failclosed_236_teeth_test.go:203`（**注释里的一次提及，不是调用**）；`BaseOptions:` 字面量形另有测试侧 3 处（都接非 nil 闭包）。追到 main 的完整链路（逐跳，本腿现量；⚠ 比 a1 多写了 `main.go:159` 与 `run.go:249` 两跳，a1 件内跳过了它们）：`cmd/wisp/main.go:58 func main` → `main.go:91 os.Exit(cmdRun(args[1:]))` → `main.go:151 func cmdRun` → `main.go:159 return runTextTask(runSpec{…})` → `run.go:181 func runTextTask` → `run.go:249 rt, code := assembleRuntime(s)` → `run.go:382 func assembleRuntime` → `run.go:780 tools.BuiltinSubagentEntries(tools.SubagentDeps{` → `run.go:789 BaseOptions: func() (agent.Options, bool) { return rt.loopOpt, rt.loopOptSet },` ⇒ **闭合；产码产出者恰 1 枚。**
4. 转述背景说"一票三形"——本腿按现量给的是：**`:256` 今天实际一票两形＋一枚非崩溃形**（①`:278` 的 nil 崩溃防护＝唯一的；②`ParentTools` 未接线的诚实拒绝文案；③`:293` 的 nil 入参＝今天不崩，是"真派生＋空目录"）。⇒ "三形"这个词在票面/a1／本件里指的不是同一件事，转述与 a1 的"三形"（＝票面甲／乙／丙三个**修法**）也撞名。建议编排者下次派单把"三形"改写成"甲／乙／丙三形"与"`:256` 承担的两件半件事"分开写，⛔ 本腿不改票面。
5. 派单红线"不许跑 `./cmd/wisp`／不许整包"我全守；本腿唯一一次 `go` 调用是 §0 那 1 发窄测。另：仓里 753 行工作树脏（`git status --porcelain | wc -l`）但 `cmd`／`internal`／`.github` 三面均 0 行，与本腿读数无冲突。
