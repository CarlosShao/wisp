# 72-a2 — 甲族 7 枚逐枚七问（本腿唯一承重件）

尺与前件（凡"红在哪一行"的读数都出自这两枚件，⛔ 转述）：
- 别名台面：`logs/aliased-7named.txt`（`TMPDIR=TEMP=TMP='C:\Users\swq\tmp\RISKAT~1\RUNNER~1'`；
  命令逐字＝`go test ./internal/risk/ -count=1 -v -run '<7 枚名册>'`；**`rc=1`／`--- FAIL`=7／`--- PASS`=0／`--- SKIP`=0**）
- 默认台面：`logs/default-7named.txt`（`TMPDIR` 清空、`TEMP=TMP='C:\Users\swq\AppData\Local\Temp'`；同一把 `-run`；
  **`rc=0`／PASS=7／FAIL=0／SKIP=0**）
- 住址尺：`git grep -n -F "func <名>(" HEAD -- internal/risk`（逐枚各 1 命中，见下表"def"列；件 `00-anchor.md` §2 记名册）
- 源码一律从 HEAD 对象层读（`git show HEAD:<path>`），行号＝本腿 10:4x–10:5x 快照，内容锚逐字抄。

⇒ 两台面**只差一枚 env**、同一份码（工作树 `internal cmd` 干净，`git status --porcelain -- internal cmd`＝0 行，每发前各量一次）。
⇒ ★**默认台面的绿⛔ 是"守卫的绿"**：本机 `%TEMP%` 的短形与长形逐字相同（我用 `Scripting.FileSystemObject.GetFolder().ShortPath` 现量＝`C:\Users\swq\AppData\Local\Temp`，⛔ 别名）。

## 四种"红形"（先给结论，逐枚证据在下面）

| 形 | 这一枚把 `t.TempDir()` 的返回拼写用在哪一侧 | 枚 |
|---|---|---|
| **I 上界/`want` 字面量** | `if !strings.EqualFold(res.Canonical, long) { t.Fatalf }`，而那个叫 `long` 的串其实是短拼 | #4 #5 #6 |
| **II 前提侧 `t.Fatalf`** | `if !strings.EqualFold(res.Canonical, p) { t.Fatalf("test premise broken…") }` | #1 |
| **III `want` 计数面** | 对 `formsOf(...).spellings()` 断言 `== 1`，而输入名册里含一枚"尚不存在的叶子" | #2 |
| **IV 把短拼当"OS 担保拼写"用** | `Gate(p, map{strings.ToLower(p): true})`，入参与 override key 同源 | #3 #7 |

---

## #1 `TestClassifyAnchorSpellingIsNotVerdict`

**def**＝`internal/risk/pathresolver_anchor_spelling_windows_test.go:29`（内容锚 `func TestClassifyAnchorSpellingIsNotVerdict(t *testing.T) {`）

1. **红的真那一行**＝同文件 `:57`，内容锚逐字＝`t.Fatalf("Resolve(%q) canonical %q: test premise broken (candidate must be the real path)", p, res.Canonical)`。
   判据位置＝**前提侧**（`:56` 那枚 `if !strings.EqualFold(res.Canonical, p)` 的分支体内）。它是 `t.Fatalf` ⇒ `:59`（`Classify(res.Canonical) != ClassA`）与 `:63`（`Gate(...).Allow || .NeedL2`）**两枚真断言在别名台面上一行都没执行**。
2. **"写死 `t.TempDir()` 拼写"在这一枚上＝** `root := t.TempDir()`（`:30`）→ `profile := filepath.Join(root, "profile")`（`:31`）→ `p := filepath.Join(profile, rel)`（`:50`），然后在 `:56` 把这个 Join 产物与 `Resolve` 的返回**逐字符串比**。⛔ 是 `want` 字面量、⛔ 是 `err==nil` 那形——是"解析器返回值 vs 测试自己拼的串"那一形，且挂在前提位上。
3. **判级对不对**＝**对，且这一枚在别名台面⛔ 验证过判级**。`got`（红句里）＝`C:\Users\swq\tmp\riskattrib1\RunnerTempAreaForAttribution1\...\profile\.ssh\id_testkey`＝完全展开的长形＝解析器干对了。别名台面下"同一个 `profile` 路径该判 ClassA"这件事**从未被跑过**⇒ 本枚的处境是**仪器哑**、⛔ 判级被证伪。
4. **恒真／无信号**＝〔结构判断〕真信号在 `:59`（分类）与 `:63`（不可放行）。若只把 `:56-58` 删掉／降级成 `t.Logf`⇒ 本枚立刻绿、且它⛔ 再⛔ 是票 72 `AC#3①` 变异名册里那枚"退 `formsOf` 就红"的靶（票面 `:70-72`、log `:139-141`）。**这一问要仪器答**（要一发 `formsOf` 退回 raw 的变异），本腿⛔ 变异授权、⛔ 跑。
5. **需真窗／真账号／真别名机器**＝需**别名台面**才现红（本腿已在盘上真别名目录 `RISKAT~1\RUNNER~1` 上复现）。`mkJunction` 走 `cmd /c mklink /J`（`pathresolver_junction_windows_test.go:48-54`）⇒ ⛔ 需管理员、⛔ 需真窗、⛔ 需真账号。
6. **`t.Skip` 逃逸口**＝**无**。尺＝`git grep -n -E '^[^/]*t\.Skip\(' HEAD -- internal/risk`＝整包 4 命中，本文件 **0 命中**。
7. **改它会碰坏哪些既有钉**＝见 `40-nails-that-would-break.md` 第 4 条（票 72 `AC#3①` 的变异名册逐枚点名它 ⇒ 改期望侧后那两发变异必须重跑）；它住在 `anchor_spelling` 文件，**⛔ 在 AC#5 那枚逐字节钉的射程内**（`logs/ac5-byte-ruler.txt` 第二发：`git diff --quiet a04d3e2 HEAD -- <本文件>` **rc2=1**，即这枚文件本来就在 `a04d3e2` 之后由票 72 自己写的，⛔ 被钉）。

---

## #2 `TestCanonicalInputGainsNoSecondForm`

**def**＝`internal/risk/pathresolver_anchor_spelling_windows_test.go:125`

1. **红的真那一行**＝同文件 `:140`，内容锚逐字＝`t.Errorf("canonical %q produced %d comparison forms %q, want 1 (verdict must not move for a path with one spelling)", res.Canonical, got, f.spellings())`。
   判据位置＝**`want` 字面量侧**（`want 1`）。⚠ 红件里 `:143`（`!f.certain()`）**没报**——机制与代码自洽：`certain()` 定义在 `pathresolver.go:327` 逐字＝`return f.real != "" && !f.partial`，而缺失叶子使 `tailExistsBelow`（`:413`）＝false ⇒ `partial`＝false ⇒ 两形照样 `certain()`。
2. **"写死拼写"在这一枚上＝** 名册里那一枚 `filepath.Join(home, "not-created-yet", "deep.txt")`（`:132`），`home` 来自 `tmpHome(t)`（`pathresolver_junction_windows_test.go:22-35`，本体＝`home := t.TempDir()` ＋ 四枚 `t.Setenv` ＋ 种 5 个文件）。它把 `t.TempDir()` 的返回拼写**当成"只有一种拼法的那一种"**送进一条对**形状**的断言。
   ★本腿量到的一处**产码事实**（⛔ 缺陷判定，只是读数）：红句里那枚 `canonical` 自身**仍带 `RISKAT~1\RUNNER~1`** ⇒ 叶子不存在时 `Resolve` 交回的是**词汇形**、⛔ 句柄展开形；而 `formsOf` 会沿祖先把 `real` 折成**长形** ⇒ 同一条路径两种形＝**2 枚 `spellings()` 是 `formsOf` 的正确行为**（`pathresolver.go:329-331` 注释逐字＝"only ever ADD a deny, never remove one"）。名册里另外 4 枚**存在的路径**（`.ssh\id_testkey`／`.git-credentials`／`plain.txt`／`home`）⛔ 红。
3. **判级对不对**＝本枚**⛔ 产任何 allow/deny 判级**（它一句 `Classify`/`Gate` 都⛔ 调）。它是票 72 `AC#5` 的**爆炸半径钉**（票面 `:96-98` 原话："Resolve 交回的 canonical 只产**一种**比较形且 `certain()==true`"）。⇒ "判级被证伪"这一问在它身上是**空集**；被证伪的是那句"恒一种形"**⛔ 带台面限定**这件事。
4. **恒真／无信号**＝〔结构判断〕若 `want` 写成 `<= 2`，或名册换成"只存在的 4 枚"⇒ 两台面都绿 ⇒ 它**有**信号；但也正因为它是单向钉（只钉"⛔ 扩爆炸半径"），把它改绿最省事的两条路（删那枚缺失叶子／改阈值）**都⛔ 是放宽**。要确证它有没有牙需要一发"把 `Resolve` 改成对缺失叶子也展开祖先"的变异——本腿⛔ 跑、⛔ 授权。
5. **需真窗／真账号／真别名机器**＝需别名台面（已复现）；其余⛔ 需。
6. **`t.Skip` 逃逸口**＝**无**。
7. **会碰坏**＝`40-nails` 第 3 条（AC#5 的爆炸半径承诺以它为唯一仪器）＋第 4 条（票 72 log `:151-152` 明写它是"新增仪器第五条用例"）。⚠ 若为消红去改 `Resolve` 的缺失叶子回退 ⇒ 波及 `Actable()`（`pathresolver.go:93`）、`syncdirs.go`（票 307 刚落地的两支比较）、`internal/tools/paths.go:131`（`Classify` 的唯一生产 caller）——**⛔ 属票 72 射程，且属"⛔ 动存值"那枚裁语的反面**。

---

## #3 `TestAListWinsWhereBothTablesHit`

**def**＝`internal/risk/pathresolver_anchor_spelling_windows_test.go:204`

1. **红的真那一行**＝同文件 `:238`，内容锚逐字＝`t.Fatalf("B-only control must stay overridable by design: %+v", d)`。
   红句逐字（`logs/aliased-7named.txt` 第 9 行）＝`pathresolver_anchor_spelling_windows_test.go:238: B-only control must stay overridable by design: {Allow:false NeedL2:true Class:B Reason:B-list sensitive path; single-file override requires one L2 confirm: id_*}`。
   判据位置＝**正控那一支**（`:237` 的 `Gate(out, map…)  !d.Allow`）。⚠ 注意 `Class:B` **是对的**。
2. **"写死拼写"在这一枚上＝** `out := filepath.Join(home, "proj", "id_bothkey")`（`:233`）＋ `Gate(out, map[string]bool{strings.ToLower(out): true})`（`:237`）——**入参拼法与 override key 都取 `t.TempDir()` 那一形**。这是**第四形**，⛔ 是 `want` 字面量、⛔ 是 Join-vs-canonical：它把短拼**当成"OS 担保拼写"**去领一张放行票，而票 72 故意⛔ 允许（见下）。
3. **判级对不对**＝**对，而且这一枚的主判定体在别名台面全部通过**：`:211-229` 那个双重循环（`both` 三枚 × `{p 短拼, res.Canonical 长拼}`）跑完了⛔ 报错 ⇒ **A∩B 判 A、且两种拼法同判**＝票 72 的不变量在这一枚上被**正面验过**（★这条是本腿自己读红件读出来的，⛔ 抄票 307）。红的只有"放行侧只认已解析形"那条窄规则被测试当成"任意拼法都放行"用。
   窄规则本体逐字＝`blacklist.go:115-126 overrideApplies`，关键两支＝`if cand.real == "" { return cand.raw != "" && bOverrides[cand.raw] }`（`:119-121`）与 `if cand.raw != cand.real { return false }`（`:122-124`）。它的来⛔ 是臆造：票面 `:86-92`（编排者插单 10:08 的落地说明）＋票面 `:236-243` 逐条回执。
4. **恒真／无信号**＝〔结构判断〕它**有**双向牙：`overrideApplies` 恒 false ⇒ `:238` 红（本腿别名台面正是这个读数）；A 规则失效 ⇒ `:219` 红。**⛔ 把它变绿的最省事形＝删掉/反转 `:237-239`**——那恰好把它自己的"防'什么都不放'"对照组拆了（`:192-197` 与 `:231-232` 两处注释都在说这件事）。要确证需一发"`overrideApplies` 直接 `return false`"的变异，本腿⛔ 跑。
5. **需真窗／真账号／真别名机器**＝需别名台面（已复现）；其余⛔ 需。
6. **`t.Skip` 逃逸口**＝**无**。
7. **会碰坏**＝★`TestOverrideOnlyAcceptsTheResolvedForm`（同文件 `:157`）＝这条窄规则的正钉，票面 `:90-92` 与 log `:171-173` 记它的变异"改回遍历全部形式"实测转红 ⇒ 任何"为了 #3 绿而去动 `overrideApplies`"的修法会**当场撞它**。另见 `40-nails` 第 5 条（台账/证据表逐枚点名它的名字）。

---

## #4 `TestPathResolverShortNameAListDenied`

**def**＝`internal/risk/pathresolver_junction_windows_test.go:123`

1. **红的真那一行**＝同文件 `:135`，内容锚逐字＝`t.Fatalf("short name did not expand to long path: got %q want %q", res.Canonical, long)`。
   红件逐字（`logs/aliased-7named.txt` 第 12 行）＝`got "C:\\Users\\swq\\tmp\\riskattrib1\\RunnerTempAreaForAttribution1\\TestPathResolverShortNameAListDenied4260811456\\001\\.git-credentials" want "C:\\Users\\swq\\tmp\\RISKAT~1\\RUNNER~1\\...\\.git-credentials"`。
   判据位置＝**上界/`want` 侧**。⇒ 真判级断言 `:137`（`Classify(res.Canonical) != ClassA`）与 `:140`（`gateFor(...).Allow`）**一行都没执行**。
2. **"写死拼写"在这一枚上＝** `home := tmpHome(t)`（`:124`）→ `long := filepath.Join(home, ".git-credentials")`（`:125`）——**变量名叫 `long`、实为短拼**（`t.TempDir()` 的返回形）。`got != want` 那形，且 `want` 是 Join 出来的字符串。★这一枚的 `got` **恰恰是完全展开的长形**＝解析器干对了 ⇒ **★本腿独立复核票 307 那句"判级⛔ 被证伪"＝成立**（件在上面，逐字红句可验，⛔ 抄）。
3. **判级对不对**＝对（同上）。⚠ 但这一枚⛔ 能"顺带就绿"：`Resolve(short)` 展开这一环是它的**唯一主题**（票 18 红队 Case 2），把 `:134` 比法改对之后仍要保住 `:137/:140` 被执行。
4. **恒真／无信号**＝〔结构判断，★且这一枚**必须靠仪器才答得了**〕两条形失效叠在一起：
   - `:134` 一红就 `t.Fatalf` ⇒ 分类断言从未执行（本腿量到的形）；
   - `:127-129` 那枚 **`t.Skip`** ⇒ 卷禁用 8.3 时整枚静默 SKIP（本腿这次 `--- SKIP`=0，件里可验，⛔ 代表它⛔ 生效过）。
   ⇒ 结论：**这一枚的"有没有牙"本腿⛔ 能判**，需两发仪器＝(i) `formsOf` 退回 raw 的变异、(ii) 一枚 8.3 关闭的台面（`fsutil 8dot3name query C:` 需提权，票 307 §追加 已记 `Error 5` 为具名欠读数）。⛔ 假装量过。
5. **需真窗／真账号／真别名机器**＝需**8.3 短名生成开启的卷**（本机 C 盘满足；`getShortPath` 见 `:58-78`，返回 "" 即无别名）。⛔ 需真窗／真账号。
6. **`t.Skip` 逃逸口**＝**有**，`:128`，内容锚逐字＝`t.Skip("volume has 8.3 short-name generation disabled (GetShortPathName returned the long name); legal workaround per ticket 18: covered by junction/case cases")`。
   ⚠ 与票 307 §追加 `:132` 引的**票 115 定式**（"种不出即 `t.Fatal`、⛔ Skip"）正相反；同包既有两枚取短名的手：`getShortPath`（本文件 `:58`，失败→`""`→上层 `t.Skip`）与 `syncdirs_redteam_windows_test.go:29 tryShortPath`（`(string, bool)`），而**票 307 刚落地的第三枚手**＝`syncdirs_alias_307_windows_test.go:38 alias307Pair`，取不到短形**即 `t.Fatal`**（`:46`，内容锚＝`t.Fatalf("AC#3 premise unbuildable: the volume generated no 8.3 alias for %q (SHORT == LONG), so this instrument measured nothing", long)`）——⛔ Skip。⇒ 这一枚逃逸口＝票 307"边界"节点名登记给票 72 的那件事（票 307 `:34`），**⛔ 属本腿射程⛔ 属修它的人顺手**：它就在 #4 体内。
7. **会碰坏**＝★**AC#5 的逐字节自证尺**。票面 `:82` 逐字＝`git diff --quiet a04d3e2 HEAD -- internal/risk/pathresolver_junction_windows_test.go` **rc=0**。
   本腿现跑＝**rc=0 成立**（件 `logs/ac5-byte-ruler.txt`）。⇒ **改 #4/#5/#6/#7 这四枚＝直接把这枚已勾 `AC#5` 的书面证据弄过期**。⛔ 写腿自裁，见 `40-nails` 第 1 条。另见 `40-nails` 第 7 条（⛔ 靠加 build tag／搬文件来消红，票 302 射程＋票 301 包级 guard D 先例）。

---

## #5 `TestPathResolverUNCAListDenied`

**def**＝`internal/risk/pathresolver_junction_windows_test.go:148`

1. **红的真那一行**＝`:164`，内容锚逐字＝`t.Fatalf("UNC spelling %q normalized to %q, want %q", spelling, res.Canonical, long)`（`spelling` 那枚是 `\\?\UNC\localhost\c$\...` 循环第一支）。
   判据位置＝**上界/`want` 侧**。⇒ `:166-168`（`gateFor(...).Allow` ⇒ `t.Fatalf`）**从未执行**。
2. **"写死拼写"＝** `long := filepath.Join(home, ".git-credentials")`（`:150`）同 #4。**★本腿额外量到次生写死**：`:151-153` 用 `filepath.VolumeName(long)` 与 `strings.TrimPrefix(long, vol)` **现场把短形段拼进 UNC**，所以即使把 `:163` 的比较改对，那三种 UNC 拼法（`:154-158`）带的仍是"人类临时目录那一形"——**夹具构造侧是这一枚的第二处要动的地方**（⛔ 只动比较那一行⛔ 够）。
3. **判级对不对**＝对。红件逐字（`logs/aliased-7named.txt` 第 15 行）里 `normalized to "C:\\Users\\swq\\tmp\\riskattrib1\\RunnerTempAreaForAttribution1\\..."` ＝长形；且红发生在 `:164` 而⛔ 在 `:161`（`t.Fatalf("UNC spelling %q: unexpected reparse error: %v", spelling, err)`）⇒ **`Resolve` 对 `\\localhost\c$` 与 `\\?\UNC\` 两种 UNC 都成功返回**＝本机可用、⛔ 需真账号／真共享。⇒ ★本腿独立复核票 307 那句"解析器干对了"在这一枚上＝**成立**。
4. **恒真／无信号**＝〔结构判断〕它的唯一安全断言是 `:166` 的"⛔ 放行"，而**⛔ 有 `Classify==ClassA` 那一支**（与 #4/#6 相比它少一枚正判级钉）⇒ 它的信号⛔ 完全依赖 `:163` 那句形状比较先绿；把 `:163` 改成"⛔ 比"＝这一枚退化成"UNC 不报错"。要确证需变异，本腿⛔ 跑。
5. **需真窗／真账号／真别名机器**＝需别名台面；UNC 走 `localhost`/`c$` **本腿量到⛔ 需真账号**（红位在 `:164`⛔ 在 `:161`）。⛔ 需真窗。
6. **`t.Skip` 逃逸口**＝**无**（整包 4 命中⛔ 落在这一枚）。
7. **会碰坏**＝同 #4 的 **AC#5 逐字节钉**（同一枚文件）；＋`40-nails` 第 5 条（台账/证据表点名它）。

---

## #6 `TestPathResolverExtendedLengthPrefixAListDenied`

**def**＝`internal/risk/pathresolver_junction_windows_test.go:174`

1. **红的真那一行**＝`:182`，内容锚逐字＝`t.Fatalf("\\\\?\\ spelling normalized to %q, want %q", res.Canonical, long)`。判据位置＝**上界/`want` 侧**。⇒ `:184-186`（`gateFor(...).Allow`）**从未执行**。
2. **"写死拼写"＝** `long := filepath.Join(home, ".git-credentials")`（`:176`）＋ `Resolve(`\\?\`+long, nil)`（`:177`）。Join 产物 vs 解析器返回的字符串直比，⛔ 是 `want` 字面量⛔ 是 `err==nil` 那形。
3. **判级对不对**＝对。红件逐字（`logs/aliased-7named.txt` 第 18 行）`normalized to "C:\\Users\\swq\\tmp\\riskattrib1\\RunnerTempAreaForAttribution1\\..."`＝完全展开长形；⛔ 在 `:179`（`t.Fatalf("\\\\?\\ prefix is not a reparse point: %v", err)`）报 ⇒ 剥前缀这一步**解析器干对了**。★独立复核成立。
4. **恒真／无信号**＝〔结构判断〕同 #5：唯一安全断言在 `:184`，把 `:181` 改"⛔ 比"就只剩"⛔ 报错"。本枚**⛔ 有** `Classify` 钉（比 #4 弱一档）。要确证需变异，本腿⛔ 跑。
5. **需真窗／真账号／真别名机器**＝需别名台面；`\\?\` 前缀本身⛔ 需特权。⛔ 真窗／⛔ 真账号。
6. **`t.Skip` 逃逸口**＝**无**。
7. **会碰坏**＝同 #4 的 AC#5 逐字节钉（同文件）；`40-nails` 第 5 条。

---

## #7 `TestBListDefaultDenyAndOverride`

**def**＝`internal/risk/pathresolver_junction_windows_test.go:272`

1. **红的真那一行**＝`:292`，内容锚逐字＝`t.Fatalf("B-list single-file override must allow: %+v", d)`。
   红件逐字（`logs/aliased-7named.txt` 第 21 行）＝`...:292: B-list single-file override must allow: {Allow:false NeedL2:true Class:B Reason:B-list sensitive path; single-file override requires one L2 confirm: .env*}`。
   判据位置＝**正控（override 必须放）那一支**。⚠ `Class:B` 是对的。
2. **"写死拼写"＝** `home := tmpHome(t)`（`:273`）→ `p := filepath.Join(home, "proj", n)`（`:279`）→ `Gate(p, nil)`（`:283`）与 `Gate(p, map[string]bool{strings.ToLower(p): true})`（`:291`）。**入参与 override key 同源、都取短拼**＝与 #3 **完全同形**的第四形（IV）。
3. **判级对不对**＝**对，且这一枚的判级主体在别名台面全部通过**：`:280`（8 个文件名逐一 `Classify == ClassB`）与 `:283`（默认 `!Allow && NeedL2`）跑完⛔ 报错（否则红会落在那两行）。⇒ **"分类⛔ 依赖拼法"这一枚是正面验过的**（★本腿自己从红位读出来的，⛔ 抄）。红的只有窄规则的用法。
   ⚠ 一处**次生**风险具名：`:294`（`if len(logs) == 0 { t.Fatal("B-list override must emit a log entry") }`）**从未到达**——它靠 `:288` 临时接管 `Logf`、`:289` `defer func() { Logf = nil }()` 还原。本腿现量（尺＝`git grep -n 'Logf =' HEAD -- internal/risk`）＝整包 **6 处命中／3 个还原形**：本枚还原成 `nil`（`:289`），而 `provenance_test.go:131-132` 与 `:825-826` 两处还原成 `old`。⇒ 改这一枚时如果只把 `:291` 弄绿、让 `:294` 第一次真的被执行，那枚"还原成 `nil`"的写法就**第一次被跑到**（`logf` 在 `blacklist.go:52` 要判空）——这是一个⛔ 在派单里、也⛔ 在票面里的落点侧风险，写腿要么保住 `:288-289` 原样、要么显式改用例内局部 logger。
4. **恒真／无信号**＝〔结构判断〕删掉 `:291-293` ⇒ 本枚变成"B 永远要 L2"＝与"什么都不放"同义 ⇒ **那一支⛔ 能删**，只能改拼法。它另有 `:297`（`main.go` 必须 `ClassNone`）当反向对照 ⇒ 整体⛔ 是恒真。要确证需变异（把 `overrideApplies` 恒 true ⇒ 本枚会⛔ 红、而 `TestOverrideOnlyAcceptsTheResolvedForm` 会红＝它把信号让给了同族另一枚），本腿⛔ 跑。
5. **需真窗／真账号／真别名机器**＝需别名台面；其余⛔ 需。
6. **`t.Skip` 逃逸口**＝**无**。
7. **会碰坏**＝★同 #4 的 **AC#5 逐字节钉**（票面 `:82`，本腿现跑 rc=0，件 `logs/ac5-byte-ruler.txt`）＋★`TestOverrideOnlyAcceptsTheResolvedForm`（`anchor_spelling:157`，窄规则正钉）。⇒ **本枚是必答 B 的最要紧靶**，详见 `30-answer-B-line.md` 形 IV。

---

## 本件用到的全部尺（可重跑）

```bash
# 名册住址（逐枚各 1 命中）
git grep -n -F "func <名>(" HEAD -- internal/risk
# 两台面颜色（每发前先 git status --porcelain -- internal cmd，应为 0 行）
R='TestClassifyAnchorSpellingIsNotVerdict|TestCanonicalInputGainsNoSecondForm|TestAListWinsWhereBothTablesHit|\
TestPathResolverShortNameAListDenied|TestPathResolverUNCAListDenied|TestPathResolverExtendedLengthPrefixAListDenied|\
TestBListDefaultDenyAndOverride'
export TMPDIR='C:\Users\swq\tmp\RISKAT~1\RUNNER~1' TEMP="$TMPDIR" TMP="$TMPDIR"; go test ./internal/risk/ -count=1 -v -run "$R"
unset TMPDIR; export TEMP='C:\Users\swq\AppData\Local\Temp' TMP="$TEMP"; go test ./internal/risk/ -count=1 -v -run "$R"
# 四数尺（对落盘件本身，含子测试行）
grep -cE '^[[:space:]]*--- (FAIL|PASS|SKIP)' <件>
# Skip 逃逸口名册（整包 4 命中，含注释豁免那一形）
git grep -n -E '^[^/]*t\.Skip\(' HEAD -- internal/risk
# AC#5 逐字节钉
git diff --quiet a04d3e2 HEAD -- internal/risk/pathresolver_junction_windows_test.go; echo rc=$?
```
