
## 编排者收 `307-r1`（2026-10-11 10:3x，落地腿五笔 `c6691756`→`51a0a0f2`→`183032c1`→`fb6a171e`→`d8743be5`，⛔ push，件 `probes/307/r1/**`；⛔ 上面各节原句一字不改，本节追加；台账＝`A834`）

★**本票五格（`AC#1`…`AC#5`）我⛔ 翻**——判语归非实现者 `307-v1`（已在飞）；本节只交"我自己复跑的那四发颜色"与"两处我派单写歪"。

### 1. 我的独立复跑（⛔ 抄腿；尺与件全部具名）

尺＝`go test -C <目录> ./internal/risk/ -count=1 -v`，四数＝对落盘件本身 `grep -cE '^[[:space:]]*(=== RUN|--- PASS|--- FAIL|--- SKIP)'`（含子测试行；退码＝先 `> 件 2>&1` 再 `rc=$?`）；驱动逐字＝`probes/307/orch/logs/orch-307-rerun.sh`；台面按上面"规矩"段那三根轴搭（`aliased`＝`TMPDIR/TEMP/TMP` 全指 `C:\Users\swq\tmp\RISKAT~1\RUNNER~1`，进场前我现量短长两形都在盘上）。

| 状态 | 台面 | rc | RUN | PASS | FAIL | SKIP | 件 |
|---|---|---|---|---|---|---|---|
| 改前（仓外 clone detach `858ee7f4`＝本票修复的父台面） | default | 0 | 198 | 197 | 0 | 1 | `orch-pre-clone-default.txt` |
| 改前 | aliased | 1 | 198 | 185 | **12** | 1 | `orch-pre-clone-aliased.txt` |
| 改后（母仓工作树，跑时 HEAD `183032c1`） | default | 0 | 200 | 199 | 0 | 1 | `orch-post-head-default.txt` |
| 改后 | aliased | 1 | 200 | 192 | **7** | 1 | `orch-post-head-aliased.txt` |

⇒ 与腿 `10-ac1-colours.md` 那四发**逐枚同数**；两枚红名册我逐名抄回＝改前 12 枚（甲 7＋乙 5，名字全列在 `A834` §1）、改后 7 枚＝**甲族原样、一枚不少一枚不多**⇒ 新增红 **0**、且**本机默认台面颜色不变**（0 红两态都是）。⚠ **这⛔ 等于"本格成立"**：颜色一致只销"腿没造假读数"那一问，有没有牙归 `307-v1`。⇒ 我在派单里另外要了它一把⛔ 自己并排写过的尺：**甲族那 7 枚要连红句一起对拉**（红名册⛔ 区分坏法种类）。

### 2. 越界与卫生（逐笔尺我自己重跑，五笔各一把）

`git show --format= --name-only` 逐笔＝名册每行都在 `internal/risk/**` 或 `probes/307/r1/**` 内 ⇒ **越出授权写面 0 行／笔**；`.out` 0；禁区（`frontend/`／`design/`／`thresholds.go`／`allowlist.txt`／`PLAN.md`）命中 0；`51a0a0f2` 原始 numstat 我逐行看过（30 行）＝新建 `.go` **1 枚**（`internal/risk/syncdirs_alias_307_windows_test.go`，`141 0`）＋`syncdirs.go` **`17 3`**（那 3 行删＝要换掉的两枚 `isUnder` 比较句本体）；其余四笔删除行 0（`d8743be5` 那 1 行删在**它自己**的 `41-boundaries.md`，主题自陈"节号更正"）。
★**记我一枚测量陷阱（本轮新形）**：我第一版那五连尺把 `51a0a0f2` 念成"删除 0 行／新 `.go` 0 枚"——**同一枚 commit 我先前用 `--stat` 明明读到 3 删与那枚新文件**。根因＝尺写成 `grep -cE '\.go$'` 而**这台机器上 git 输出行尾带 `\r`**，`$` 锚全部落空 ⇒ 计数静默归 0。⇒ 定式＝**凡以行尾锚取名的 git／log 尺，先 `tr -d '\r'`；全称式"0 枚"落笔前必有一把把原始名册逐行打印的尺对拉**。

### 3. ★记我派单两处（⛔ 算腿的欠账，处置写明）

1. **尺写歪**：我给腿的建树尺逐字是 `git archive HEAD internal cmd | tar -x` ⇒ 那棵树**没有 `go.mod`/`go.sum`**，`go test ./internal/risk/` 直接 `cannot find main module`。腿按盘上事实把 `go.mod go.sum` 一并 archive 才跑得起来，并**具名报回了这处不符**⇒ 处置正确；"改了尺之后那批变异读数还有没有效"交 `307-v1` 裁。（我在这枚票上第三处写歪，前两处见追加节 §5 第 2、3 条。）
2. **写面漏给票面文件**：我给腿的授权写面＝`internal/risk/**` ＋ `probes/307/r1/**`，**里面⛔ 有本票这枚文件** ⇒ 它无法履行"每腿在票尾追加一条 Progress log"这条全局规矩（它自己具名报了这枚相撞，并因此连 `AC#4`/`AC#5` 的框都⛔ 碰——那是正确的：撞上⛔ 禁区时选窄的那一支，⛔ 猜）。⇒ **处置＝那两条 bullet 由我代提**（见下面 Progress log，正文标注"由编排者代写、读数出自腿件"）。⚠ 定式＝给落地腿划写面时，票面文件要么一起给，要么当场说清"Progress log 由我代提"。

### 4. 现态与 next

- 现态＝**本票 `1 勾／5 未勾`**（尺＝`grep -cE '^[[:space:]]*- \[[ x]\]'`＝6、`- \[ \]`＝5，改前改后同数⇒本笔⛔ 动框）；票 306 **同期收口 `-done`**（`A835`）；票 303 6 勾／1 未勾⛔ `-done`；票 305 1 勾／5 未勾（前提动摇，见 `A834` §3）；票 72 4 勾／2 未勾⛔ 动框。
- 在飞＝**`307-v1`**（本轮编队 1 枚）。我要它额外答两枚⛔ 在任何格判据里的必答：① 生产面这台机器上这次修复⛔ 观测得到颜色差别（本机 `C:\Users\swq` 逐段 SHORT==LONG、且生产 `roots` 空 ⇒ 唯一阳性出口只有 `:428`）；② `AC#1` 那格凭据属**测试台面**还是**生产台面**。
- next＝收 `307-v1` → 我逐格翻或退回 → 推送⛔ 排（`A813` 未过期仍⛔ 用；`51a0a0f2` 起 `internal/risk` 那批 CI 色要和票 300／303／111 欠的 CI 色同批作差）→ 我跑票 305 那六发 → 队列 `298-v1`→`296-r2`→`295-r1`→`294-r1`→`294-v1`。

- Progress log：
  - [2026-10-11 10:36:54 +08] agent=307-r1（★正文由编排者代写，读数出自它自家件 `probes/307/r1/**`，逐枚点名如左）did=票 307 `AC#1` 落地＋`AC#2`–`AC#5` 凭据（五笔：`c6691756` 起手锚→`51a0a0f2` 产码＋`AC#1` 两台面名册→`183032c1` 四形变异→`fb6a171e` 收口→`d8743be5` 补件）＝**产码**：`syncdirs.go` `match()` 两支比较都换成包内既有形比较族（`candForms := formsOf(cand)` ＋ `aUnder(candForms, anchorForms(e.canon))` / `…anchorForms(s.home)`），**存值／`finalize()`／`resolveTarget()` 那两句 `errTargetUnverified`／`ErrReparseDenied` 一字节⛔ 动**（我用对象层尺复认：numstat `17 3`，那 3 行删就是被换掉的两句本体），非测试面新增 `filepath.Clean|Abs`＝0；**新建夹具 1 枚**＝`internal/risk/syncdirs_alias_307_windows_test.go`（`//go:build windows`，141 行，2 枚用例，取不到短形即 `t.Fatal`、`t.Skip` 0 枚＝票 115 定式）。**颜色**＝`go test ./internal/risk/ -count=1 -v` 两台面×改前改后各两发：default `198/197/0/1`→`200/199/0/1`、aliased `198/185/12/1`→`200/192/7/1`（余 7＝甲族归票 72，⛔ 少⛔ 多；RUN 多的 2 枚＝它自己新增那两枚夹具）；**交集新增红 0、单发红那枚 0**（件 `logs/full*-{pre,post}.txt`＋`counts-summary.txt`，编排者我自己四发复跑逐枚同数）。**`AC#2` 四形**＝(a) 只撤 `:428` ⇒ 5 枚指名全红（★它具名报回：`TestSyncFixtureFallbackAndMatch` 这一形下红在 **`:175`** 而⛔ 票面写的 `:162`，因 `:162` 改由已修好的 `:422` 拿到 `Sync:true`）；(b) `sha256`/`cmp` 三处还原逐字节相同⇒复绿；(c) ★中间形（只换 `:428`、`:422` 退回）用一枚只活在导出树里的打印用例具名回答"`:162` 从哪一支拿到阳性"＝**`:428`**（`Root.Provider=sync-suspect`／`Source=suspect-fallback`／`Why` 是 `:428` 那句）⇒ 编排者那枚读码推到成立，据此它交了钉住 `:422` 的那枚新用例 `TestSync307AliasShortRootMembership`；(d) 恒真形＝**5 枚里只有 3 枚响**（`TestSyncFallbackNotDisarmableByWeakRoot`／`TestSyncUnverifiedRootKeepsFallback` 体内⛔ 负向钉＝天生无信号），它⛔ 写"本格闭合"、只写"包级仪器层成立（同发另有 8 枚既有钉一起红）"＋枚级欠账具名。**门禁**＝`go build ./...`／`go vet ./internal/risk/`／`d22scan` 各 `rc=0`；格式名册两把并排＝工作树 `gofmt -l internal/risk` 1 枚（`provenance.go`＝`i/lf w/crlf` **假枚**，⛔ 计债⛔ 动）、HEAD blob 导出树 **0 枚**；0 字节件 0、`.out` 0、probes 下新 `.go` 0；逐笔越界尺（它自陈**第一版误用区间尺**把编排者的两笔数了进来、已改逐笔）。**⛔ 做到的**＝⛔ 翻 `AC#1`–`AC#5` 任何框、⛔ push（0 笔）、⛔ 在共享工作树做变异（全在 `C:/Users/swq/tmp/wisp-307-r1-tree-main`）；**具名欠读数**＝CI runner 两枚台面色（⛔ 本机可量）、逐卷 `fsutil 8dot3name query C:`（需提权 `Error 5`）、`alias307Pair` 的 Fatal 分支本机零执行、枚级耗时在 0.01 s 粒度下量不出来（它⛔ 主张"代价为零"）。⚠ **它报回两枚与派单不符**＝① `git archive HEAD internal cmd` 那把尺缺 `go.mod`/`go.sum` 跑不起来（它补了才建树，记编排者派单）；② 授权写面⛔ 含票面文件 ⇒ 它无法履行"票尾 Progress log"（⇒ 本条由编排者代提）。
  - [2026-10-11 10:36:54 +08] agent=编排者 did=收 `307-r1`（五笔逐笔越界我自己重跑＝写面外 0／`.out` 0／禁区命中 0；`51a0a0f2` 原始 numstat 30 行逐行看过＝新 `.go` 1 枚＋`syncdirs.go` `17 3`）＋**两台面四把我自己跑**（件 `probes/307/orch/logs/orch-*.txt`，与腿账**逐枚同数**：改前 aliased `198/185/12/1`、改后 `200/192/7/1`、default 两态 0 红；红名册逐名抄回＝改后剩的 7 枚正是甲族、乙 5 枚全翻绿、新增红 0）＋★**自抓一枚测量陷阱并改正**：`grep -cE '\.go$'` 在这台机器上因 git 输出行尾带 `\r` 而**静默归 0**（我第一版因此把 `51a0a0f2` 念成"删除 0 行／新 `.go` 0 枚"，与我自己先前 `--stat` 那发矛盾）⇒ 定式＝行尾锚的 git 尺先 `tr -d '\r'`，"0 枚"这种全称读数必配一把逐行打印原始名册的尺对拉；记我派单两处（建树尺缺 `go.mod`／写面漏给票面文件 ⇒ 本票 Progress log 由我代提）。⇒ **本票五格⛔ 动，等 `307-v1`**（在飞；我另要它答两枚⛔ 在任何格判据里的必答：生产面这台机器可观测性、`AC#1` 凭据属测试台面还是生产台面）。现态＝票 307 1 勾／5 未勾；票 306 本波收口 `-done`（`A835`）；next=收 307-v1→逐格翻或退回→推送⛔ 排（A813 仍⛔ 用，`51a0a0f2` 起那批 CI 色同批作差）→我跑票 305 那六发→队列 298-v1→296-r2→295-r1→294-r1→294-v1
