# 票 307 — `syncdirs` 那道"写在用户资料目录下 ⇒ 当同步可疑"的兜底**只折一边**：候选走 C26 解析器折成长形，`home` 只小写留短形 ⇒ 资料目录下的写读成"不落在任何同步根"（**fail-open**，方向与票 252 相反）

**立票**：2026-10-11 08:0x 编排者（来路＝只读归因腿 `risk-red-attrib-1` 的表②乙族＋**我自己的两发复跑**，件 `probes/risk-attrib/orch/logs/orch-aliased.txt`／`orch-control.txt`；台账 `A817` §5 那句"归口另待一枚普查腿定射程"由这一程结掉，`A831`）
**性质**★：**缺陷票**，方向＝**少报**（该被当外逃通道看的写被放行），⛔ 仪器票、⛔ 文档票。
**为什么要紧（用户可见形状）**：这条链是"宿主写文件之前那道风险判定"的一个输入。它现在只在**拼法恰好对得上**的时候才响——如果用户目录那一侧的拼写带着 8.3 短名（或任何解析器会折长、而 `normPath` 不折的形状），那一次写就**不会被标成 sync-suspect**，⇒ 后面那道"要不要问用户"的门少了一枚输入。⚠ **这不是"有人入侵了这台机器"**：现象出现在**测试与风险判定的字符串比较层**，本机没有第二条真账号、没有实测被绕过的生产写（见"射程与⛔ 射程"）。

## 现量（每条带尺；⛔ 引用前先重跑，行号一律当快照）

- ★**产码那处不对称**（编排者 `sed -n` 逐字复现，HEAD `e4740e35` 之后）：
  - `internal/risk/syncdirs.go:88`＝`s := &syncSet{home: normPath(orDefault(o.HomeDir, userHomeDir())), …}` ⇒ `home` **只过 `normPath`**；
  - `normPath`（`internal/risk/blacklist.go:133`）逐字＝`strings.ToLower(unifySeparators(p))` ＋ 剥尾分隔符 ⇒ **折小写与斜杠，⛔ 折 8.3／⛔ 走解析器**；
  - `internal/risk/syncdirs.go:197 resolveTarget()` ⇒ `res, err := Resolve(raw, s.excs)` → `canon, err := res.Actable()` → `return normPath(canon), nil` ⇒ 候选那一侧**先经 C26 折成长形再小写**；
  - `internal/risk/syncdirs.go:429`＝`if !s.complete && s.home != "" && isUnder(cand, s.home)`，而 `isUnder`（`blacklist.go:169`）逐字＝`path == dir || strings.HasPrefix(path, dir+sepStr)` ⇒ **纯前缀比较，一边长一边短就永远不等**。
  - 同一处旁边 `:425`＝`if isUnder(cand, e.canon)` 走的是 `roots` 那一支（`e.canon` 是否同形＝**本票要现量的一格，见 `AC#0`**）。
- ★**生产调用者＝有，一枚**（尺＝`grep -rn --include=*.go '\.IsSyncPath(' internal cmd | grep -v _test.go`＝**1 行**）＝`internal/risk/provenance.go:928`：`if p.IsSyncPath(t).Sync { return true, ChSyncWrite }` ⇒ `Sync:false` 直接让那次写**不被当成同步写**。
- ★**决定性的两发（编排者自己跑，⛔ 抄腿）**：同一份码、同一把尺 `go test ./internal/risk/ -count=1 -timeout 300s -v -run '<12 枚名册>'`，只差 `TEMP/TMP/TMPDIR` 一枚 env＝
  - 仓外一枚**真有 8.3 别名**的目录（`C:\Users\swq\tmp\riskattrib1\RunnerTempAreaForAttribution1` → 短形 `C:\Users\swq\tmp\RISKAT~1\RUNNER~1`，短形我自己用 `Scripting.FileSystemObject.GetFolder().ShortPath` 取的）⇒ **`rc_aliased=1`、`PASS=0`／`FAIL=12`／`SKIP=0`**；
  - 本机默认 TEMP（`realtemp_short` 与长形**逐字相同**＝这台机器 TEMP 没别名）⇒ **`rc_control=0`、`PASS=12`／`FAIL=0`／`SKIP=0`**。
  - 红句逐字（本票那族的一枚，件＝`orch-aliased.txt`）＝`syncdirs_test.go:402: under-profile path must be sync-suspect, got {Sync:false Root:{Provider: Path: Source:} Why:write target is not under any sync root}` ⇒ **"少报"这一支被当场造出来**，⛔ 靠读码推。
- ⚠ **本机绿⛔ 是"守卫的绿"**：这台机器 TEMP 没有 8.3 别名，所以那 12 枚在本机全绿只证明"这台机器的拼法恰好对得上"。⇒ 本票**任何一格⛔ 用"本机全绿"当凭据**（腿表②把这形叫"缺世界"，我认）。
- **归因分族**（腿的表②，甲 7＋乙 5＝12＝名册行数；我复认了乙族的 5 枚与名册）：⚠ 甲族那 7 枚**⛔ 属本票**（它们的红是"用例期望侧写死了 `t.TempDir()` 的拼写"，判级本身没被证伪）⇒ **归口＝票 72**（那 7 枚用例就是它写的），见下面 `AC#5` 的对账格。

## 要建什么

- [ ] **`AC#0` 先把 `:425` 那一支现量清楚（⛔ 预设答案）**：`roots` 侧的 `e.canon` 与 `home` 侧的 `s.home` 是不是同一把折法？尺＝逐枚赋值点 `file:行`＋那条路径到底过不过 `Resolve`（⛔ 按注释读）。若 `e.canon` 也是短形 ⇒ 同一枚缺陷有**两**处出口，`AC#1` 的落点要同时覆盖；若只有 `home` 单边 ⇒ 落点只有一处。**这一格的答案决定落点形状，⛔ 由落地腿自己裁。**
- [ ] **`AC#1` 把兜底那一边折成同形**：判据＝在**真有 8.3 别名**的台面上，`syncdirs_test.go` 那 **5** 枚指名用例（`:162`/`:213`/`:233`/`:355`/`:402` 那族）**全绿**，而本机默认台面**颜色不变**（⛔ 把别的台面弄红）。**硬约束**＝① 折法**⛔ 新造一套**——只许复用现成的 `Resolve`/`Actable` 或既有 `normPath`＋一次具名展开，⛔ 在 `risk.PathResolver` 之外用 `filepath.Clean|Abs` 做判定（AGENTS §1.2，D22）；② 改完⛔ 让"该 fail-closed 的那几支"变松（`resolveTarget` 那两句 `errTargetUnverified`／`ErrReparseDenied` 的 fail-closed 语义一字不许动）；③ ⛔ 顺手把 `allowed_dirs` 做成硬边界（它是判级输入，先例已定）；④ ⛔ 动 `internal/panel/tokens_fourway_test.go`、`l2_grant_boundary_test.go`、`internal/perm/ticket90_persist_test.go` 三枚冻结件与 `thresholds.go`／`allowlist.txt`。
- [ ] **`AC#2` 反形敏感度（★本票最要紧的一格，⛔ 只交"改⇒红"）**：除了"改完变绿"，必须各交一发**把修复挪开**的凭据＝(a) 把 `home` 那一侧退回单边折法 ⇒ 那 5 枚当场红、**红句逐字引**；(b) `cmp` 逐字节还原 ⇒ 复绿；(c) ★**第三形**＝"把判据换成反形它也不响"那一问（本项目记忆第 107–155 条那一族）：例如只把测试的 `want` 改成恒真／只断言 `err == nil` 的那形，**指名用例必须仍然响**，否则那枚用例对这件事不敏感、本格⛔ 成立。⚠ 三缺一不算闭。
- [ ] **`AC#3` 台面判据要能进仪器（⛔ 造 `--- SKIP` 换绿）**：今天那 12 枚的颜色**只在带别名的台面上才存在**——CI 托管 runner 的 TEMP 就是别名形（`C:\Users\RUNNER~1`），本机不是。⇒ 本票要交＝**修复后两把台面各自的颜色**（别名台面 12 枚与本机台面 12 枚，逐名同集合），并**具名回答**："如果 CI 那台机器的 TEMP 哪天不再有别名，这 5 枚会不会变成静悄悄的绿？"——若会，就要一枚**自己造得出那个世界**的夹具（⛔ 依赖机器形状），并给那枚夹具的落点与它⛔ 需要真窗／真账号的说明。⛔ 为省事把这几枚搬进"需真 OneDrive／真账号"那一档。
- [ ] **`AC#4` 门禁与越界**：`GOFLAGS= go build ./...`／`go vet ./internal/risk/`／`sh scripts/d22scan.sh` 各 rc=0；格式名册**两把并排**（工作树／HEAD blob，各自写射程目录）；`internal/risk` 整包改前／改后**各两发、同一台面取交集**＝稳定新增红 0 枚，且**单发红过的那一枚⛔ 被这把尺藏住**——必具名报出并写清归因（先例＝票 303 `AC#3b` 第③件）；每把门禁件**自落一行 `rc=N`**（0 字节件＝那格没交；⚠ "0 字节"那把尺必带 `-type f`，本轮编排者自己踩过目录那一发）；`git show --name-only --format=` **逐笔**量越界（⛔ 区间尺）；⛔ 零 push；commit 显式 pathspec 写在 `$( … )` **之外**、只落自家 `probes/307/<腿名>/**`（票 301 `AC#4b`）。
- [ ] **`AC#5` 三本账的对账（⛔ 同一件事记两次）**：逐枚写清"⛔ 翻谁的格"——① 台账 `A817` §3/§5 那"12 枚〔待归因〕"⇒ **由本票＋票 72 合起来销**（甲 7 归票 72、乙 5 归本票），归因凭据＝`probes/risk-attrib/1/**` ＋ 编排者那两发 `probes/risk-attrib/orch/logs/orch-*.txt`；② **票 252**＝同一种"一条路径两种拼法分家"的**先例**（它的后果方向相反：根内写读成越界＝**多报**）⇒ 本票⛔ 扩它的面、⛔ 在它身上翻勾，只引它当形状先例；③ **票 72**＝那 7 枚用例的母票 ⇒ 甲族的修法⛔ 进本票（本票⛔ 动那 7 枚的期望侧）。先例＝`A817`／票 302 那处"同一物理缺陷链上只记一次"。

## 边界（⛔ 塞进本票）

- ⛔ **甲族那 7 枚用例的期望侧**（归票 72）；⛔ **`junction` 文件里那枚 `t.Skip` 逃逸口**与**"标题那件事在 runner 上从未被执行"**（腿的表③"仪器三形"）＝登记，落点见 `A831`，⛔ 本票顺手修。
- ⛔ **反向那枚 `TestC21TableColourRowsMatchTokensCSS`**（`internal/ball/tokens_table_test.go:1465`，判据读工作树资产字节）＝**⛔ 与这 12 枚同族**（腿独立判的，我复核其尺：`os.ReadFile` 工作树资产，成因属资产/放置那一族）⇒ 归票 77／`P1` 那条旧账，⛔ 本票。
- ⛔ **`internal/risk` 那 12 枚的搬运／档位**＝票 302 的射程（`A817` §5 原话：混一票必然顺手改错那半）。
- ⛔ 任何"要不要把它做成执行时硬边界"的扩张（先例：`allowed_dirs` 是判级输入，不是边界）。

## 规矩（本票全程）

- 子代理**只 commit、⛔ push**；⛔ `git add -A`／`.`；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`；⛔ 在仓内建 worktree 或 checkout；临时件**只建不删**（且从一开始就落自家 `logs/`，⛔ 落仓库根）；证据件⛔ 叫 `.out`（根 `.gitignore` 有全仓 `*.out`）。
- ⛔ 为变绿放宽任何断言；⛔ 造 `--- SKIP`；⛔ 动 SLO 阈值／golden／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt`／D43 转移表／C1–C32／D1–D47／三枚冻结件。
- ★**台面规则三根轴**（票 305 立的第ⓐⓑ轴＋本轮 `305-a2` 现量的第⏂轴，本票同样吃）：ⓐ母仓↔干净 clone 颜色⛔ 可比；ⓑ母仓里那枚未跟踪 `frontend/dist/index.html` 自己会漂；⏂**干净 clone／导出树⛔ 带 `third_party/sherpa-onnx` 与 `build` 两枚 DLL 目录**（尺＝`git ls-files third_party/sherpa-onnx`＝**0**、`git ls-files build`＝**0**，编排者 08:0x 现跑复认）⇒ 照抄 `PATH="$PWD/…"` 进 clone＝`0xc0000135` 且零 `--- FAIL`（用例根本没跑的假绿形）⇒ 跑 `cmd/wisp` 一律 PATH 写**母仓绝对** `/d/…` 路径，且**第一发必须先断言具名用例真发出 `--- PASS/FAIL/SKIP`**。本票⛔ 需要 `cmd/wisp`，但这条记在这儿防下一枚腿踩。
- 裁决者≠实现者（`SPEC-12 §4.3`）：`AC#1`/`AC#2`/`AC#3` 的判语⛔ 由落地腿自勾。
- 凡报枚数必写"哪把尺＋射程目录＋含不含注释行＋工作树还是 HEAD blob"；凡写"排除 N 枚"，N 必须等于那张名册的行数。

next＝`307-a1`（只读，答 `AC#0` 那格：`:425` 与 `:429` 两支的折法名册＋最小诚实落点＋会碰坏哪些既有钉）→ 我裁落点形状 → `307-r1`（落地腿，`AC#1`＋`AC#2` 三形）→ `307-v1`（非实现者终裁）。

## 编排者收 `risk-red-attrib-1`（2026-10-11 08:0x，四笔 `63bc6335`→`f4fd71bc`→`1f9e781f`→`8936b100`，件 `probes/risk-attrib/1/**`；⛔ 我自己的复跑＝这一节正文）⇒ **`A817` 那 12 枚〔待归因〕结掉**：甲 7＝**归本票**、乙 5＝**立成新票 307**

### 0. 台面与越界（⛔ 抄腿自陈）

- 起手锚＝`e4740e35`（腿自己重锚的；我派单写的 `6284a489` 又腐烂了——**同族第三枚**，定式照旧＝派单号只当"派发时刻快照"，腿一律自己 `rev-parse`）。
- 逐笔越界尺我自己在四笔上重跑（`git show --name-only --format=`）＝名册每行都在 `.scratch/wisp/probes/risk-attrib/1/**` 下；新建 `.go`＝**0** 枚；`.out`＝**0**；0 字节件＝**0**（尺＝`find … -type f -size 0`，⚠ 我上一程踩过不带 `-type f` 那发）。
- 编译面：它自报了跑过哪些 `go test`，我复核＝只落 `./internal/risk/`，⛔ `cmd/wisp`、⛔ `./...`、⛔ `d22scan`、⛔ 格式名册（我给的那条特批边界守住了）；⚠ 它自陈一枚自己的尺缺陷＝`T02` 写了 `echo "rc=$rc"` 却漏了取退码那一拍 ⇒ 那格退码空（颜色由件尾与四数承担，**我认它不是假读数、但记为格式缺陷**）。

### 1. 承重读数我自己复跑（⛔ 抄腿；每条带尺）

| 腿的读数 | 我这发的独立复跑 |
|---|---|
| 分母＝12 枚**全在 `internal/risk`**（3 个文件） | 尺＝逐枚 `grep -rn --include=*_test.go "func <名>(" internal/risk/`：**12 枚各 1 命中**（我自己列的名册行数＝12）✓；它件里另出现过第 **13** 个名字 `TestAListWinsBothTablesHit`，我量＝**defs 0 枚**（⚠ 那是**它引我 `A817` 时脱了一个 "Where" 造出来的不存在名字**，盘上真名 `TestAListWinsWhereBothTablesHit` 在台账里 4 处命中；⛔ 因此追认任何格，也⛔ 让它替我 `A817` 背书——**注释／引用里的测试名一律当待验断言**，这条我早写过，本轮又多一枚实例） |
| ★那 2 枚"同文件同族"**⛔ 同文件** | `TestAListWinsWhereBothTablesHit` 在 `pathresolver_anchor_spelling_windows_test.go:204`、`TestBListDefaultDenyAndOverride` 在 `pathresolver_junction_windows_test.go:272`——**我现跑对上** ⇒ `A817` §3 那句"同文件"确实过期（同包、不同文件、同成因），**更正记在我账上** |
| 决定性两发＝别名台面 12 枚全红／本机台面 12 枚全绿 | ★**我自己重跑了这两发**（同一份码、同一把 `-run` 名册，只差 `TEMP/TMP/TMPDIR`）：短形我用 `Scripting.FileSystemObject.GetFolder().ShortPath` 自己取＝`C:\Users\swq\tmp\RISKAT~1\RUNNER~1`（长形目录确实存在，腿⛔ 删、合规）；读数＝**`rc_aliased=1`／`PASS=0`／`FAIL=12`／`SKIP=0`** 与 **`rc_control=0`／`PASS=12`／`FAIL=0`／`SKIP=0`**，件＝`probes/risk-attrib/orch/logs/orch-aliased.txt`（内含 `rc_aliased=1` 行）／`orch-control.txt`（`rc_control=0`）⇒ **成立，且这是我本轮唯一一次"腿的凭据我自己原样复现"** |
| 乙族那 5 枚＝**fail-open**（该可疑的写读成不可疑） | 我自己 `sed` 读了四处产码：`syncdirs.go:88` `home: normPath(…)`；`blacklist.go:133` `normPath`＝`ToLower(unifySeparators(…))`（⛔ 折 8.3）；`syncdirs.go:197 resolveTarget()`＝`Resolve`→`Actable()`→`normPath(canon)`（**折长**）；`syncdirs.go:429 isUnder(cand, s.home)`＋`blacklist.go:169 isUnder`＝纯前缀 ⇒ 一边长一边短＝永不等；我自己的红句逐字＝`syncdirs_test.go:402: under-profile path must be sync-suspect, got {Sync:false … write target is not under any sync root}` ⇒ **机制成**，⛔ 读码推 |
| 生产调用者＝1 枚 | `grep -rn --include=*.go '\.IsSyncPath(' internal cmd \| grep -v _test.go`＝**1 行**＝`internal/risk/provenance.go:928`（我现跑）⇒ **不是"建了没接"**，这条判级真在卡里被读 |
| 本机 TEMP ⛔ 别名 | 我这发量＝`realtemp_short` 与长形 `C:\Users\swq\AppData\Local\Temp` **逐字相同** ⇒ 本机绿⛔ 是守卫的绿（腿那格"缺世界"我认） |

### 2. ★裁：乙族 5 枚⛔ 挂本票，**立成新票 307**（我不照腿的建议走，具名理由）

腿的表③建议"乙 5 枚挂票 252（同一物理缺陷的第二枚调用点，⛔ 另立一票记两次）"。我**不采纳**，理由三条，逐条可驳：
1. **后果方向相反**：票 252 那一族是"根内的写读成越界"＝**多报**（用户体验＝被多问一次）；本票这 5 枚是"资料目录下的写读成不可疑"＝**少报**（安全输入丢失）。两族修法在"折哪一边、折完谁变严"上是**相反的一支**——这正是 `A817` §5 自己警告过的那种"混一票必然顺手改错那半"。
2. **票 252 自己还有 2 格未闭**（现量＝4 勾／2 未勾，task #298：`AC#1`/`AC#6` 的非实现者读数与 `AC#4`/`AC#5` 按 `cmd/wisp` 空），且⛔ 改 `-done`。把一枚**新落点**并进去＝让它的收口去等一枚本不属于它的修复，且会把"谁翻哪格"糊成一团。
3. **"同一物理缺陷链上只记一次"⛔ 被读成"同形只开一票"**：这条规矩管的是**同一处事实记在两本账**（先例＝`A817`／票 302）；本票与 252 是**不同代码、不同判定消费者**（`PathResolver`／R2 升 L2 ⇄ `IsSyncPath`／`provenance.go:928`），形状相同而已。⇒ 处置＝**立票 307**，并在 307 的 `AC#5` 里**双向写死对账**（⛔ 在 252 身上翻勾、只引它当形状先例）。撤销本裁语＝口令**「307 并回 252」**（那我把 307 五格转写进 252 的追加格，⛔ 删 307 文件、改名留档）。

- **甲族 7 枚＝归票 72**（本票⛔ 收）：那 7 枚用例是票 72 写的，且它票面自己写着"AC#4 的 runner 半边不在代理手里"；腿现量到机制＝**期望侧写死了 `t.TempDir()` 的返回拼写**（红在前提／上界／`want`／正控那一行，**判级本身⛔ 被证伪**——#4/#5/#6 的 `got` 恰恰是完全展开的长形＝解析器干对了）。⇒ 这一格**归本票射程⛔**，我只在此登记"归口已定＝票 72"，⛔ 动票 72 一枚框（票 72 现态＝4 勾／2 未勾、`Status: review`，由它的下一程自己读这一节）。
- **腿的表③"仪器三形"（缺世界／`junction` 那枚 `t.Skip` 逃逸口／#11 那件事在 runner 上从未被执行）＝登记**，落点＝票 307 的"边界"节（⛔ 本票顺手修）＋票 72 的下一程射程。
- ⚠ **`A817` §3 那格"2 枚机制未证"⇒ 已填**；`A817` 欠账⑥（12 枚〔待归因〕）**由本节结掉**：分族＝甲 7（→票 72）＋乙 5（→票 307），⛔ 丙、⛔ 丁（腿响亮报了"丙那把尺我⛔ 打，⛔ 变异授权"，我认＝**丙未测是本票与票 72 的欠读数**，⛔ 算腿的欠账——派单里我就⛔ 给它变异授权）。

### 3. 现态与 next

- 现态＝**票 307 新立（0 勾／6 格）**；票 305 **1 勾／5 未勾**（下一节收 `305-a2`）；票 306 **1 勾／5 未勾**（`306-v1` 在飞）；票 303 **6 勾／1 未勾**＋规则 6 那道闸（`A830` §5）；票 72 **4 勾／2 未勾**、⛔ 动框。
- next＝`307-a1`（只读答 `AC#0`）→ 我裁落点 → `307-r1` → `307-v1`；同批补位＝`303-v2`（补那张与 AC 编号 1:1 的七行裁决表）；我自己在 `cmd/wisp` 车道那四发照 `A830` §7 排队。
- Progress log：
  - [2026-10-11 08:0x +0800] agent=编排者 did=收 `risk-red-attrib-1`（四笔，逐笔越界我自己重跑＝越出 `probes/risk-attrib/1/**` **0 格**、新 `.go` 0、`.out` 0、0 字节件 0）＋**承重读数我自己复跑**：12 枚名册逐枚 `grep "func <名>("` 各 1 命中；`TestAListWinsBothTablesHit` defs＝**0**（腿引我 `A817` 时脱字造出的不存在名，⛔ 追认、记我"引用里的测试名一律待验"再中一枚）；那 2 枚"同文件"⛔ 成立（盘上 `anchor_spelling:204` ⇄ `junction:272`，`A817` §3 那句过期、更正记我）；★**决定性两发我原样复现**＝别名台面（短形我自己用 COM `ShortPath` 取：`C:\Users\swq\tmp\RISKAT~1\RUNNER~1`）`rc=1`／`PASS=0`／`FAIL=12`，本机台面（TEMP 长短形逐字相同＝⛔ 别名）`rc=0`／`PASS=12`／`FAIL=0`，件 `probes/risk-attrib/orch/logs/orch-aliased.txt`／`orch-control.txt` 各带 `rc=` 行（⚠ 我第一发⛔ 把 rc 落进件、只 echo 到 stdout，补了 note 行，踩的是我自己那条规矩）；产码四处我自己 `sed` 复现（`syncdirs.go:88` 只 `normPath`、`blacklist.go:133` ⛔ 折 8.3、`resolveTarget:197` 走 `Resolve`→`Actable`→折长、`isUnder:429` 纯前缀）＋红句逐字 `under-profile path must be sync-suspect, got {Sync:false …}` ⇒ **fail-open 机制成立**；`IsSyncPath(` 非测试面调用点＝**1**（`provenance.go:928`）⇒ 不是"建了没接"。⇒ **裁**：乙 5 枚⛔ 挂票 252、**立票 307**（三条理由：后果方向相反／252 自己两格未闭⛔ 再等一枚新落点／"只记一次"管的是同一事实记两本账⛔ 是"同形只开一票"，撤销口令「307 并回 252」）；甲 7 枚归票 72（机制＝期望侧写死 `t.TempDir()` 拼写，⛔ 动票 72 一枚框，现态 4 勾／2 未勾）；`A817` 欠账⑥**结掉**；丙那把尺⛔ 打＝⛔ 算腿欠账、算 307／72 的欠读数；`TestC21TableColourRowsMatchTokensCSS` 独立判⛔ 同族（读工作树资产字节，归票 77/`P1` 旧账）。⚠ 台面⛔ 三根轴：`git ls-files third_party/sherpa-onnx`＝0、`git ls-files build`＝0 ⇒ 干净 clone／导出树⛔ 带那两枚 DLL 目录，照抄 `$PWD` PATH 进 clone＝`0xc0000135` 且零 `--- FAIL`。现态＝票 307 **0 勾／6 未勾**；在飞＝`306-v1`；next=`307-a1`（只读答 `AC#0`）→我裁→`307-r1`→`307-v1`，同批 `303-v2` 补 1:1 表；⛔ push（`A813` 未过期、我仍⛔ 用＝推送自启 `slo-full` 抢机主 CPU，且票 300／303／111 那批 CI 色要同批取）
