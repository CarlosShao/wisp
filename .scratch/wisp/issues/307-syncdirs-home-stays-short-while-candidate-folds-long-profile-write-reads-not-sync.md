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

- [x] **`AC#0` 先把 `:425` 那一支现量清楚（⛔ 预设答案）**：`roots` 侧的 `e.canon` 与 `home` 侧的 `s.home` 是不是同一把折法？尺＝逐枚赋值点 `file:行`＋那条路径到底过不过 `Resolve`（⛔ 按注释读）。若 `e.canon` 也是短形 ⇒ 同一枚缺陷有**两**处出口，`AC#1` 的落点要同时覆盖；若只有 `home` 单边 ⇒ 落点只有一处。**这一格的答案决定落点形状，⛔ 由落地腿自己裁。**
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
