
## 编排者收 `307-a1`（2026-10-11 10:1x，只读普查腿，四笔 `08cc614e`→`4d4c10d4`→`205ce69f`→`293a4d34`，件 `probes/307/a1/**`；⛔ 上面各节原句一字不改，本节追加）

★**`AC#0` 那一格我翻勾**（凭据＝我这发**逐条自己重跑**，⛔ 抄腿；见第 2 节），并据此**裁 `AC#1` 的落点形状＝形 C 换比较函数那一支、且两处出口都换、⛔ 改任何存值**（第 3 节）。撤销本裁＝口令**「307 落点改形 B」**（那我把 `:160` 的存值一并改写，代价见第 3 节那一枚陷阱栏）。

### 1. 台面与越界（⛔ 抄腿自陈）

- 逐笔 `git show --name-only --format=` 我自己重跑＝四笔名册每行都在 `.scratch/wisp/probes/307/a1/**` 内，**写面外 0 格**；新建 `.go`＝**0**、`.out`＝**0**、0 字节件＝**0**（尺＝`find .scratch/wisp/probes/307/a1 -type f -size 0`，带 `-type f`）。
- ⚠ 腿自报两枚它自己的缺陷，我复核**处置正确**：第 1 笔里 `logs/00-anchor-readings.log.2` 与 `g-k-crossface.txt.2` 当时 0 字节（那两把尺的**空输出就是有效读数**＝`git status --porcelain -- internal cmd` 0 行、跨面 grep 0 命中）；它⛔ 删、⛔ `--amend`，在第 2 笔往件里补了"empty by design＋尺名＋rc"说明行 ⇒ **HEAD blob 上不再 0 字节，历史那一笔的 blob 确实 0 字节**。记为格式缺陷，⛔ 读数缺陷。
- 编译面：腿自陈 Go 面发数＝**0**（连 `go env`／`go list` 都没用），我复核它的件里**没有任何"我跑过 test"的声称**，凡后果一律标〔读码推到〕⇒ **本票的颜色读数欠着，由 `307-r1`＋验收腿两台面各取**（第 4 节）。

### 2. 承重读数我这一发重跑（每条带尺，⛔ 引用前先重跑；行号＝10:0x 快照）

| 腿的读数 | 我这发的尺与结果 |
|---|---|
| 票面 `:425`/`:429` 已腐烂，真值 **422/428** | 尺＝`grep -n "isUnder(cand" internal/risk/syncdirs.go` ⇒ **422**（`isUnder(cand, e.canon)`）与 **428**（`!s.complete && s.home != "" && isUnder(cand, s.home)`）⇒ **成立**。同段 `:88`/`:197`/`blacklist.go:133`/`:169`/`provenance.go:928`/`:345` 与测试侧 5 枚行号我另用 `sed` 逐枚对，**全部同号未腐烂** ⇒ **入账方式＝内容锚**（同 `A831` 对 `gate.go` 的定式） |
| `home` 赋值点＝2 枚、**0 枚过 `Resolve`** | 尺＝`grep -rn --include=*.go "home: normPath" internal` ⇒ 逐字两行 `provenance.go:345`＋`syncdirs.go:88`，**两边正文完全同形** ⇒ 成立；`normPath` 本体我 `sed -n '126,145p'` 复看＝只 `ToLower`＋`unifySeparators`＋剥尾分隔符 ⇒ **⛔ 折 8.3** |
| `e.canon` 赋值点＝3 枚，其中 **`:160` 连 `lexCanonical` 都没过** | 尺＝`grep -n "c\.canon" internal/risk/syncdirs.go` ⇒ `:152`/`:156`/`:160`（＋`:163` 那枚比较）⇒ 成立；`:160` 正文＝`normPath(expandInput(r.Path))` |
| **生产面 `roots`＝空 ⇒ `:422` 那个循环零次迭代** | 我逐环重跑：`grep -rn --include=*.go "NewProvenance(" internal cmd \| grep -v _test.go` ⇒ 非测试调用者＝**2 枚**，且 `cmd/wisp/run.go:743`／`cmd/wisp/panel_assets.go:239` **都传 `ProvOptions{NoProbe: true}`**；`grep "SyncRoots:"` 与 `grep "SyncFixturePath:"` 非测试命中各 **0 行**；`NoProbe` 分支只 `add` 这两路 ⇒ `finalize()`（我 `sed -n '118,138p'` 复看：`complete` 只在"canonical＋gradeConfirmed"的根上置真）在 0 枚根上 ⇒ `complete` 恒 false ⇒ **生产唯一阳性出口＝`:428`** |
| 那 5 枚用例**全部走 `provenance.go:345`**、无一枚走 `:88` | 尺＝逐枚 `sed` 取函数体里的 `NewProvenance(ProvOptions{…})` ⇒ 五枚**全带 `NoProbe: true`** ⇒ 成立。唯一吃 `:88` 的 fold 断言＝`syncdirs_test.go:501`（同发现量 `p := NewProvenance(ProvOptions{}) // live env`＋前置 `if p.SyncDetectionComplete()`） |
| 形 B 那枚空串陷阱 | `sed -n '167,172p' internal/risk/blacklist.go` ⇒ `isUnder(path, dir) = path == dir \|\| HasPrefix(path, dir+sepStr)` ⇒ `dir==""` 时＝`HasPrefix(path, sepStr)` ⇒ **POSIX 上每条绝对路径全命中** ⇒ 陷阱**成真** |
| ★**我自己读到的一条、腿没写**：`aRelates` 会跳过空形 | `sed -n '270,287p' internal/risk/blacklist.go` ⇒ 双层循环里 `if c == "" { continue }` 与 `if a == "" { continue }` ⇒ **走"换比较函数"那一支时，枚空锚不会变成全局前缀**；这条把形 B 的最大风险直接消掉，是我裁形 C 的主要理由之一 |
| ★**我自己读到的一条、腿没写**：`:422` 那个根循环**⛔ 受 `!s.complete` 保护** | `sed -n '421,429p'` ⇒ `for _, e := range s.roots { if isUnder(cand, e.canon) … }` 在 `if !s.complete && …`（`:428`）**之前且无该条件** ⇒ **改 `:422` 的比较函数会波及"已检出确认根"那一批台面**，而 `uncertainAnchorMiss`（`blacklist.go:314-322`，需两侧 `root` 非空且互为祖先）有可能把某枚"⛔ 是 sync"翻成"是 sync" ⇒ 这是 `AC#1`"⛔ 把别的台面弄红"最可能现形的地方，`307-r1` **必两台面各一发**（我⛔ 在此替它判颜色） |
| `:162` 那一枚在形 A 下会**为错误的原因绿** | 我这发把两环接起来：fixture 根 `home\Cloud\NutstoreSync` **从未创建**（`sandbox()` 只建 `profile`/`local`/`roaming`，`syncdirs_test.go:15-27`；`:155` 现量），`Resolve` 对缺失叶子走 `pathresolver.go:127-140` 的词汇回退（`Resolved=false`、且 `Rewritten＝len(kinds)>0` 在这里＝false）⇒ 落 `:160` ⇒ 根成员那一支在别名台面上⛔ 命中；而候选确实在 `home` 之下 ⇒ 只折 `home` 会让它从 **`:428`（suspect-fallback）** 拿到 `Sync:true`，**用例标题那句"under fixture root must be sync"就名不副实了**（能区分两支的断言只有 `syncdirs_other_test.go:99` 那枚 `Root.Source`，而它在 `//go:build !windows` 档）⇒ **单折 `home` 我不接受**，落点必须同时把 `:422` 的比较换形 |

### 3. ★裁语：`AC#1` 落点形状＝形 C 换比较函数、两处出口同换、⛔ 动存值

**形状（写给 `307-r1`，⛔ 由它自勾框）**：

1. `:428` 把 `isUnder(cand, s.home)` 换成包内既有的 `aUnder(<候选那一侧的 pathForms>, anchorForms(s.home))`；
2. `:422` 把 `isUnder(cand, e.canon)` 换成同一族（`aUnder(…, anchorForms(e.canon))`）；
3. **存值一字节⛔ 改**：`home` 那两枚赋值点、`c.canon` 那三枚赋值点、`finalize()`、`resolveTarget` 那两句 `errTargetUnverified`／`ErrReparseDenied` 的 fail-closed 语义**全⛔ 碰**（腿在 `20-landing-shapes.md` 给形 B 算的那枚语义扩张——"把根不可验证从 `canonical=false` 变成报错"——我**⛔ 采纳**，理由见下面第 4 点）；
4. `s.home != ""` 那枚守卫原地留着（`blacklist.go:199` 既有写法 `homeS != "" && aUnder(…)` 就是先例）；
5. `cand` 那一侧怎么变成 `pathForms`＝**本票唯一一处新写法**（它已是 `resolveTarget` 折过的字符串，再进 `formsOf` 会多一次祖先走步）。⇒ 判据＝**落地腿必须自己现量这一枚的代价**（同一台面 `go test ./internal/risk/ -run 'TestSync' -count=1 -v` 改前／改后各自总耗时＋一枚具名用例的 `--- PASS (0.xx s)`），并把数写进件里；⛔ 由我在此预填数。

**为什么是这一形（四条，逐条可驳）**：

- **复用而非新造**＝`AC#1`① 的字面要求：`formsOf`／`anchorForms`／`aUnder`／`uncertainAnchorMiss` 都在同一个包里，且 `blacklist.go:199`–`:221` 那 8 处**早就在拿它比 `home` 锚**，票 72 的不变量原话（`blacklist.go:189-193`）逐字＝"the verdict must not depend on the spelling a path arrived in"。折法一套⛔ 出现第二套。
- **单调朝 fail-closed**＝`pathresolver.go:330-337` 那句注释逐字"only ever turn a miss into a deny, never remove one" ⇒ 本票的病是**少报**，这一族天生只朝"多报"那一侧动，方向对；而形 A／形 B 是靠"把某一侧折得更长"来加命中，方向相同但⛔ 带工具性保证。
- **⛔ 空串陷阱**（第 2 节我那枚 `aRelates` 读数）＝形 B 需要额外写守卫，形 C 天生不需要。
- **⛔ 碰 `finalize()`／`canonical`**＝腿在 D 组算清楚的那件事：`complete` 对包外不可观测（非测试调用者＝**0**，我这发复核 `grep … | grep -v _test.go` 只剩定义行），但它**是测试可见的契约**（13 枚调用点里 8 枚是前置闸）⇒ 换比较函数⛔ 动分级，红就不会先出现在那批闸上。**这一条是我⛔ 采形 B 的主因**：形 B 改的是"哪枚根算已确认"，那等于顺手改了 P12 安全网能不能被卸掉这件事（票 307 的边界节写着"⛔ 顺手把 `allowed_dirs` 做成硬边界"，同一族道理）。

**登记、本轮⛔ 做（三枚，具名）**：① `:431` 仍把 `normPath` 的比较形当 `Root.Path` 报出去，与 `blacklist.go:128-132` 那句 ticket 75 的注释相抵——**换形会改变它报出的拼法**，本格⛔ 修、⛔ 顺手"改对"，只记；② `SyncRoots()` 的消费面（`syncdirs.go:181-182` 注释声称面板 Security page 吃这枚判定）腿只按注释读、⛔ 复核 ⇒ 归我车道一发只读尺；③ 逐卷短名状态（`fsutil 8dot3name query C:` 需提权、`Error 5`）＝**具名欠读数**，⛔ 算腿的欠账。

### 4. `AC#2`／`AC#3` 的判据按 `307-a1` 的读数收窄成三句可执行的话

- `AC#2` 第三形（★本票最要紧那一格）现在有了具体靶子：**必须交一发"只折 `home`、⛔ 换 `:422`"的中间形**，并在件里具名回答"`:162` 在这一形下从哪一支拿到 `Sync:true`"。⇒ 若它从 `:428` 拿到（我这发的读码推到＝会），那枚用例**对"根成员判定"这件事⛔ 敏感**，`AC#2` 那格⛔ 成立，除非另加一枚断言 `Root.Provider/Source` 的用例把它钉住（判据形状＝`syncdirs_other_test.go:99` 那枚的 Windows 同形）。
- `AC#3` 的"自己造得出那个世界"＝腿在 `40-machine-alias.md` §2 已证成：**这台机器新建长名目录取得到短名**（`wisp307a1AliasWorldProbe\LongSegmentNameBeyondEight` → `WI0253~1\LONGSE~1`），且包内已有两枚 OS 级取短名的手（`pathresolver_junction_windows_test.go:58 getShortPath`、`syncdirs_redteam_windows_test.go:29 tryShortPath`）——⚠ **这两枚都带"无别名就 `t.Skip`"的出口**，而票 115 的定式是"种不出即 `t.Fatal`、⛔ Skip"。⇒ `AC#3` 的交付＝**新夹具⛔ 依赖 `%TEMP%` 的机器形状**（`t.TempDir()` 下自己建一枚 >8 字符的目录、向 OS 要它的短形、取不到即 `t.Fatal`），⛔ 是把那 5 枚搬进"需真 OneDrive／真账号"那一档。
  ⚠ **台面边界一条（腿的 `90-final.md` §2 与我这边的既有事实）**：8.3 是 Windows 独有的，而 `syncdirs_test.go` **无 build tag、POSIX 也求值**（`A817` 那条编译面差）⇒ 新夹具只能落在**带 `//go:build windows` 的新文件**里，⛔ 是给那 5 枚加 tag（那会把它们从 core 档分母里摘掉，撞票 301 的包级 guard D——票面"边界"节已写死）。
- **★必答那一问我照腿的读数认：会。** 换一台 profile 无短名别名的机器，这 5 枚**静悄悄全绿**（⛔ SKIP、⛔ 失败消息、⛔ 任何自陈），因为触发条件在断言之前就不存在 ⇒ **本机绿⛔ 能当凭据**（腿的第 1 环我这一发复跑同形：`C:\Users\swq` 与 `%TEMP%` 全链逐段 SHORT==LONG，而 `riskattrib1` 那一支 SHORT!=LONG）。

### 5. 与派单不符处 4 枚——**全部我认，且有两枚是派单写错**

1. 票面 `:425`/`:429` **腐烂**（真值 422/428）⇒ 入账方式改内容锚（第 2 节），⛔ 改票面原句。
2. 我在派单里把 `provenance.go:345` 说成"同形复制件"＝**说轻了**：它是**生产面唯一那枚 `home` 赋值点**（两个产码调用者都 `NoProbe:true`，我这发现量），那 5 枚用例也全走它 ⇒ **"要不要一起算进落点"⛔ 是覆盖问题，是"⛔ 覆盖它就没修到生产"的问题**。记我派单一处。
3. 我在派单里把第三形说成"住在 `pathresolver.go`"＝**说错了地方**：`formsOf`/`anchorForms` 定义在 `pathresolver.go:345`/`:391`，但**做 `home` 锚点比较那一族住在 `blacklist.go`**（8 处 `anchorForms(...)` ＋ `aUnder`/`aEq`/`aRelates`/`uncertainAnchorMiss`）⇒ 本程第 3 节的裁语照此写，改动点在 `syncdirs.go:422`/`:428` 调 `blacklist.go` 的关系函数。
4. 腿报"票面把 5 枚都算作 `home` 单边折法的出口，而 `:162` 同时吃 `:422` 的 `:160` 支"＝**票面归因不完整，我认**，并按它改写 `AC#2` 的第三形判据（第 4 节第 1 句）。⛔ 改票面 `AC#0`–`AC#5` 任何一句原句，本节就是把那半句补齐的地方。

### 6. 现态与 next

- 现态＝**票 307 `AC#0` 翻勾 ⇒ 1 勾／5 未勾**（尺＝`grep -cE '^[[:space:]]*- \[[ x]\]' <本票>` 改前 0／6、改后 1／5）；票 306 6 勾／0 未勾⛔ `-done`（欠 core 档 CI 色）；票 303 6 勾／1 未勾⛔ `-done`（`AC#3`＝永久未勾，见 `A832`）；票 305 1 勾／5 未勾；票 300 6 勾／1 未勾；票 298 0 勾／5 未勾；票 72 4 勾／2 未勾。
- 在飞＝**`307-r1`**（落地腿，写面＝`internal/risk/**` ＋ `probes/307/r1/**`）。⚠ **车道约束具名**：它动 `internal/risk`＝`cmd/wisp` 的导入图，⇒ **它落笔期间我⛔ 在母仓跑票 305 那六发**（那六发要跨 9 枚 commit 作差，混进别家源码改动就整批发废），改在它交完之后跑。
- next＝① 收 `307-r1` → 我自己重跑两台面名册 → 派 `307-v1`（非实现者终裁，⛔ 由实现腿自勾 `AC#1`/`AC#2`/`AC#3`）；② 取 38103274181 的 CI 色 → 票 306 `-done`＋索引；③ 我跑票 305 那六发（`probes/305/a2/30-deciding-run-recipe.md`）；④ 补票 300 `AC#4`；⑤ 队列 `298-v1`→`296-r2`→`295-r1`→`294-r1`→`294-v1`；⑥ 我那两枚欠读数＝`SyncRoots()` 消费面＋`:431` 报出的拼法。

- [2026-10-11 10:07:54 +08] agent=编排者 did=收 307-a1（四笔，逐笔越界我重跑＝写面外 0／新 .go 0／.out 0／0 字节件 0；承重读数逐条我自己重跑：422/428 真值、home 两枚 0 过 Resolve、e.canon 三枚其中 :160 不过、生产两枚调用者都 NoProbe:true 且 SyncRoots:/SyncFixturePath: 非测试 0 命中 ⇒ :422 生产零次迭代、5 枚全走 provenance.go:345）＋我自己补了两枚腿没写的读数（aRelates 跳过空形 ⇒ 形 B 的空串陷阱在换比较那一支不存在；:422 根循环不受 !s.complete 保护 ⇒ 改它的波及面要在两台面现量）⇒ AC#0 翻勾、裁 AC#1 落点＝形 C 换比较函数两处同换⛔ 动存值⛔ 碰 finalize()（撤销口令「307 落点改形 B」）；记我派单两处错（把 provenance.go:345 说成"同形复制件"＝其实是生产唯一赋值点；把第三形说成住 pathresolver.go＝真身住 blacklist.go）；:162 在只折 home 的形下会从 :428 拿到 Sync:true ⇒ AC#2 第三形判据改写＋必交那枚中间形 next=307-r1 在飞（写面 internal/risk/** ＝ cmd/wisp 导入图，故我⛔ 此时跑票 305 那六发）；⛔ 翻 AC#1/AC#2/AC#3、⛔ push
