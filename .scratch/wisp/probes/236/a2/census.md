# 236-a2 — 格式门 8 枚脏样逐枚普查（归属＋引用面）

- 腿：`236-a2`（只读普查）。起手锚点〔我现跑〕`git rev-parse --short HEAD` ＝ **`98df640a`**（branch `dev`）。
- **交回时锚点已漂**：〔我现跑〕终态 `git rev-parse --short HEAD` ＝ **`a207d129`**（其间推进 **1 枚**）。`git diff --name-only 98df640a..a207d129` ＝ **只 1 枚文件**：`docs/evidence/s1/235-pool-nail-and-bridge-preread-v1.md` ⇒ **8 枚脏样一枚都没被碰**，我在 `98df640a` 取的全部字节／行／列读数在 `a207d129` 逐枚复量**同值**（复量输出见下）。
- 零编译／零 `go test`／零 `go build`／零 `go vet`／零 commit／零 push。唯一被授权的静态尺＝`gofumpt -l`／`-d`（纯静态、不编译）。
- 骨架先落盘（八行表＋节标题），之后逐节追加。每条结论标〔我现跑〕或〔转述 236-a1，未复跑〕，两者不混。

---

## 0. 八行总表（终态；全部〔我现跑〕）

| # | 路径 | 字节 | 行数 | 脏在哪（具名到行） | 归属票 | 引用面结论 | 产品码快照？ |
|---|---|---|---|---|---|---|---|
| 1 | `.scratch/wisp/probes/163/a1/main.go` | 3,927 | 98 | 纯格式 1 处：文档注释缺 `//` 分隔行（`:5` 前插入） | 票 163（立而不派） | **只路径提及**（4 处），0 行列／0 字节／0 md5 | 否 |
| 2 | `.scratch/wisp/probes/174/c2/zz174c2_wiring_pair_windows_test.go` | 5,353 | 156 | 纯格式 1 处：`:141` map 字面量闭括号未独占行 | 票 174（ready-for-agent） | 只路径提及 ＋ **1 处可复跑 grep 尺**（格式中性），0 行列／字节／md5 | 否（overlay 夹具） |
| 3 | `.scratch/wisp/probes/183/r2/mut/task-boxset-off.go` | 15,275 | 376 | 纯格式 1 处：`:253` 少一个 tab（脱缩进） | 票 183（已 `-done`） | 只路径提及（含裁决表具名引用它作突变载体），0 行列／字节／md5 | 否（`task.go` 的 overlay 突变副本） |
| 4 | `.scratch/wisp/probes/185/r1/mut-m1/hostpath_185.go` | 4,980 | 89 | 纯格式 1 处：`:89` 整个函数体挤一行 | 票 185（已立按住不派） | **仓内 0 处字节/md5**；唯一 md5 线索在 `185-paged-reread-r1.md:93`，但那枚 md5 **"in the leg reply"**（聊天正文），不在任何仓内文件里 | 否（overlay 突变副本） |
| 5 | `.scratch/wisp/probes/197/r1c/pre/subagent_197.go` | 19,314 | 434 | 纯格式 1 处：`:305–308` `Result{…}` 字面量括号未独占行 | 票 197（在派） | 只路径提及 ＋ **具名声明它＝`d1c440ad` 的 blob 抽取件**（字节保真本身就是凭据） | **是**：产品码在树，27,205 字节，**快照比本体小 7,891** |
| 6 | `.scratch/wisp/probes/197/r1c/pre/subagent_197_test.go` | 25,429 | 696 | 纯格式 **6 处**（见 §1 枚 6） | 票 197（在派） | 同上（同一条 `:93` 声明覆盖两枚） | **是**：产品码在树，34,952 字节，**快照比本体小 9,523** |
| 7 | `.scratch/wisp/probes/222/v1/mutations/m11-budget-50ms-gated/subagent_222_test.go` | 20,076 | 532 | 纯格式 2 处：`:184`／`:291` 字段**对齐**多一空格 | 票 222（产码已交、六格未勾） | **有实锤行列引用**：`222-...-v1.md:130–132` 具名引**突变件行号** `:385`／`:400`／`:418`；⇒ 但本枚修法**行数零变化**（两 hunk 均 7→7），行号全部存活 | 否（overlay 突变副本） |
| 8 | `.scratch/wisp/probes/185/c1/mut/fs_broken.go` | 11,041 | 338 | **解析错误**（非格式）：`:4:1 imports must appear before other declarations`，尺 rc=**2** | 票 185 | **引用面最重但全是行列＋错误文案**：`fs_broken` 命中〔我现跑〕**20 处**；字节数 `11041`/`11,041` 命中 **0 处**；md5 **0 处** | 否（正控夹具） |

**八枚共同事实**〔我现跑〕：**全仓 `docs` ＋ `.scratch/wisp/issues` 里，没有一处引用这 8 枚里任何一枚的字节数或 md5/sha1**。我把 8 个字节数（含千分位写法）逐枚当字符串全扫了一遍，命中的 15 处全是假阳（提交号 `11f3927`／`0e95353`、JSON 里的 `10498048`），逐条抽查后**无一枚是真·字节引用**。

---

## 1. 逐枚详情

### 尺的现跑复核
- `gofumpt -l` 逐枚喂这 7 枚 ⇒ **7 枚全部仍列**〔我现跑〕，与 236-a1 名册**逐字相同**。
- `gofumpt -d` 逐枚 ⇒ 各枚 rc=1 并给出可定位 hunk；第 8 枚 ⇒ **rc=2 ＋ 单行错误**〔我现跑，把 a1 的转述升成实测〕。
- 落盘：`diff-1…diff-8`（8 份）。

### 枚 1 — `probes/163/a1/main.go`
- 1 个 hunk：`@@ -2,6 +2,7 @@`。包文档注释里编号列表（`:2–4` 三行 `// 1./2./3.`）与 `:5` `// logdir 由旗标传入…` 之间缺 `//` 空注释分隔行。
- **不是 import 顺序、不是未定义符号**；文件本来就能编译（票 163 的正经只读核台件程序）。
- 整理代价：插 3 字节（`//\n`）⇒ 3,927→**3,930**、98→**99** 行。
- 引用面〔我现跑，尺＝`grep -rn "163/a1/main" docs .scratch/wisp/issues`〕：4 处 —— `docs/evidence/s1/197-subagent-entity-r1b.md:102`（表格：路径＋入库号 `e13196bd`）、`docs/evidence/s1/200-project-instructions-r2.md:152`（名册）、`docs/reports/pending-and-issues.md:8918` 与 `:9610`。**全部只是路径提及**（配提交号），没有行列、没有字节、没有 md5。

### 枚 2 — `probes/174/c2/zz174c2_wiring_pair_windows_test.go`
- 1 个 hunk：`@@ -138,7 +138,8 @@` 即 **`:141`** —— `c2Call(t, b, "fs.read", map[string]any{"path": filepath.ToSlash(shape.path)})` 跨行写时闭括号没独占行。
- 引用面〔我现跑〕7 处。要紧的**不是**路径提及，是这一条：
  - `docs/evidence/s1/174-wiring-cost-census-c2.md:101` 在这枚文件上挂了一把**可复跑的尺**：「`grep -c "filepath.Clean\|filepath.Abs" probes/174/c2/zz174c2_wiring_pair_windows_test.go`＝0」。⇒ 这是**内容属性引用**而非行列/字节引用；本枚修法只动一处换行，**不引入任何 `Clean`/`Abs`** ⇒ 该尺复跑读数不变 ⇒ **零凭据风险**。
  - 同文件 `:19`、`:54` ＝ 路径＋它靠 `overlay.json` 被 `-overlay` 编进 `internal/tools`（"跟踪件零字节"）。其余为名册路径提及（`197-…r1b.md:103`、`200-…r2.md:152`、`pending:8918`、`issues/174-*.md:60`）。
- 整理代价：字节＋行数都变（+1 行）。

### 枚 3 — `probes/183/r2/mut/task-boxset-off.go`
- 1 个 hunk：`@@ -250,7 +250,7 @@` 即 **`:253`** —— `_ = hostPathBoxFromCtx(ctx)` 少一层 tab。
- ⚠ 关键定性：这枚的"做坏点"是**语义**（把 `box.set` 摘掉、只留 `_ =` 丢弃调用），**语法面 gofumpt 只看到缩进**。修缩进＝语义零改动。
- 归属＝**票 183**〔我现跑核对票面：标题＝「宿主铸造的指针豁免在真机 CLI 上没护住宿主自己的指针…」，Status 带 `-done` 后缀〕。它的裁决表 `docs/evidence/s1/183-reread-permanent-teeth-r2.md:109` **具名把它称作"牙"的载体**：「摘掉 `internal/tools/task.go` 里那三行 `box.set`（…副本 `mut/task-boxset-off.go`、`overlay-boxset.json`、`logs/tools-183-on-boxset-off.txt`）⇒ 两枚新腿双双红」，同文件 `:226` 再引其日志。
- 我现跑 `overlay-boxset.json`：`internal/tools/task.go` ← 本枚文件。⇒ **它是可复跑夹具**，裁决表引的是「跑它得到那两枚 FAIL」这件事，**不是它的字节**。整理代价：`:253` 前 +1 tab ⇒ **15,275→15,276、行数不变**，md5 变；复跑红读照样成立。
- **0 处行列引用／0 处字节引用**〔我现跑，行号扫描 `task-boxset-off\.go:[0-9]+` 命中数＝**0**〕。

### 枚 4 — `probes/185/r1/mut-m1/hostpath_185.go`
- 1 个 hunk：`@@ -86,4 +86,10 @@` 即 **`:89`** —— `func tail24(s string) string { r := []rune(s); if len(r) > 24 { return … }; return s }` 全挤一行，gofumpt 要展开成 7 行。
- `:88` 注释自陈"tail24 is mutation M1 only: relax the whole-path compare to the last 24 runes" ⇒ 做坏点＝语义（后 24 rune 比较），语法只是挤一行。
- 归属＝**票 185**〔我现跑核对票面：标题＝「每一次成功的续读都会亲手造出下一次续读的阻断者…」，Status＝已立、按住不派〕。
- ⚠ **本枚是 8 枚里唯一沾 md5 的**，但引用**不在仓内**：`docs/evidence/s1/185-paged-reread-r1.md:93` 逐字＝「mutants + overlays: .scratch/wisp/probes/185/r1/mut-m1..m4, ov-m1..m4.json (**pre/post md5 in the leg reply**)」。⇒ md5 落在那条腿的**回报正文**里（仓外），**仓内没有任何文件存着这枚的 md5 或字节数**。后果要如实说：改字节会让**那条腿回报里的数**对不上，但那是仓外凭据、不属"事后复算会对不上的盘上凭据"。
- 同名撞车提醒：`hostpath_185.go` 也是产品码 `internal/risk/hostpath_185.go` 的名字。其余 5 处命中（`33-inbound-listener-r2.md:25/29`、`pending:8451/8466/8469`）逐条抽查**全部指产品码那枚**，与本枚无关。
- 语义核 `tail24` 在 `docs`＋`issues` 里〔我现跑〕**0 命中** ⇒ 没有任何裁决表逐字抄过这枚函数体。
- 整理代价：`:89` 一行→7 行 ⇒ 89→**95** 行、字节约 +60，md5 变。**行号会变**（本枚是唯一行数显著变的格式枚）——但仓内无人引它的行号。

### 枚 5 — `probes/197/r1c/pre/subagent_197.go`
- 1 个 hunk：`@@ -302,10 +302,12 @@` 即 **`:305–308`** —— `return Result{Text: fmt.Sprintf(…), IsError: true}, nil` 括号未独占行。中文回执文案在 hunk 里是**上下文行、原样保留** ⇒ 修格式不动那三段字。
- **是产品码快照**〔我现跑，逐尺〕：
  - 产品码路径＝`internal/tools/subagent_197.go`，`git ls-files` **在树里**（tracked、存在）。
  - 快照现量 **19,314 字节**；`git cat-file blob d1c440ad:internal/tools/subagent_197.go | wc -c` ＝ **19,314** ⇒ **逐字节相等**，证实裁决表声明的来源。
  - `git cat-file blob HEAD:internal/tools/subagent_197.go | wc -c` ＝ **27,205** ⇒ 快照比产品码本体**小 7,891 字节**（本体自 09-28 后又长了）。
  - 出处声明（我现读）＝`docs/evidence/s1/197-subagent-entity-r1b.md:93`：「改前（**本腿亲手量**：把 `d1c440ad` 的两枚 blob 抽到 `.scratch/wisp/probes/197/r1c/pre/` 再 `gofumpt -l`）｜点出 **2 枚**」。同文件 `:124`、`:272` 再引 `pre/`。
- ⚠ **这一枚给乙法新增一条 a1 没量的代价**：凭据不是"字节数被抄下"，而是**"它＝`d1c440ad` 的逐字节副本"这件事本身**＋那句"`gofumpt -l` 点出 2 枚"的**改前读数**。格式化它之后，任何人复算 197-r1b 那句改前读数会**跑空**（`-l` 不再点出这 2 枚）⇒ 与 185 同形、但**是第二处**，且它此前没被记进票面 AC#6 乙段的代价清单。
- ⚠ 甲法在这枚上的代价**是零**：裁决表 `200-project-instructions-r2.md:155-156` 已具名记「`internal/tools/subagent_197*.go` 那两枚不是我清的：它们在我起手后由 197-r1b 的 `7ea14ce3` 交件时格式化掉了」⇒ **产品码那两枚今天已是净的、照被扫**，`pre/` 免检**不藏任何产品码形状**（现量证据见 §2）。

### 枚 6 — `probes/197/r1c/pre/subagent_197_test.go`
- **6 个 hunk**，全部同一形状（跨行字面量／参数表括号未独占行）〔我现跑，hunk 头逐字〕：
  - `@@ -78,8 +78,10 @@` → **`:81–82`**（`_ = emit(llm.StreamEvent{Type: …EvToolCallArgsDelta,` ＋ `ArgsDelta:`）
  - `@@ -167,7 +169,8 @@` → **`:170–171`**（`func (h *sub197Harness) buildWith(…)` 参数表尾）
  - `@@ -201,8 +204,10 @@` → **`:207–208`**（`h.bridge = New(Options{…` ＋ `Logf:`）
  - `@@ -318,8 +323,10 @@` → **`:321–322`**（`prov := &fake197Provider{…`）
  - `@@ -365,8 +372,10 @@` → **`:368–369`**（同型 `&fake197Provider{…}`）
  - `@@ -431,8 +440,10 @@` → **`:434–435`** 附近（`child := &fake197Provider{…` ＋ `askSpawnFirstTurn:`）
- **是产品码快照**〔我现跑〕：产品码 `internal/tools/subagent_197_test.go` **在树里**；快照 **25,429** 字节 ＝ `d1c440ad` 同路径 blob **25,429**（逐字节相等）；`HEAD:` 现量 **34,952** ⇒ 快照比本体**小 9,523 字节**。
- 引用面同枚 5（`197-…r1b.md:93` 一句覆盖两枚）。行号扫描〔我现跑〕：`subagent_197_test.go:<数字>` 在 docs＋issues 有命中，但**逐条抽查全部指产品码**（`issues/197-*.md:48`、`issues/211-*.md:23`、`issues/212-*.md:6`、`issues/221-*.md:49`、`issues/222-*.md:100/:19`、`issues/235-*.md`）——**没有一处指 `pre/` 这枚**。
- ⚠ a1 提醒过的漂移已在本腿证实存在但**不影响本枚读数**：`a71f0be9`/`1d2ad737` 动的是**产品码**测试件（现量 34,952 字节／614 行那类），`pre/` 快照仍是 `d1c440ad` 的原样。

### 枚 7 — `probes/222/v1/mutations/m11-budget-50ms-gated/subagent_222_test.go`
- 2 个 hunk，**纯对齐**〔我现跑〕：`@@ -181,7 +181,7 @@` → **`:184`** `late      time.Duration`；`@@ -288,7 +288,7 @@` → **`:291`** `late:     200 * time.Millisecond,`。两处各多一个空格。
- ⚠ 本枚的"做坏点"（50ms 预算门控）**不在这两个 hunk 里**（那在两处对齐之外，是语义旋钮）。
- 归属＝**票 222**〔我现跑核对票面：标题＝「`task.spawn` 一边占着桥的执行许可、一边等孩子跑完…」，Status＝产码已交并静态核过，⛔ 六枚格一枚没勾〕。
- **引用面＝本腿挖出的最实质一条**（裁决表 `docs/evidence/s1/222-spawn-holds-bridge-permit-v1.md`）：
  - `:130–132` **具名引"突变件"行号**，逐字：「`:395`（**突变件里漂到 `:400`**）四枚父任务全收到…」「`:413`→`:418` 四枚孩子的结论各出现 0 次」「`:380`→`:385` 因果读数」。我现读**证实这三枚数指的就是本枚探针**：探针 `:385`／`:400`／`:418` 逐字是 `t.Errorf("孩子 %s 的工具调用起跑时已有 %d 枚父任务返回…`／`t.Errorf("父任务 %d 收到的是错误而不是结论：%q"…`／`t.Errorf("结论 %q 出现在 %d 枚父任务回复里, want 1（%v）"…`；而**产品码现行**同处号位内容**不是**这些。
  - ⇒ **决定性结论**：本枚的修法（两处各去一个空格）**行数为零变化**（两个 hunk 都 7→7）⇒ 那三枚被引的突变件行号**逐枚存活**，`sha1` 变但**仓内无人引它**（我算探针 sha1＝`98cbc06f432179912d01f0dbf1d571cb3302bc8f`，全仓 `docs`＋`issues` **0 命中**）。
  - ⚠ 澄清一处**容易误判成凭据**的：`:127`「`subagent_222_test.go:70 h222PreFixBudget`」我逐枚验号后**指产品码当时的版本**，不指探针——`beaeaeba`/`0a22b1a7`/`bfc55b5b`/`f23586fa` 四枚提交下产品测试件都 527 行、`:70` 逐字都是 `h222PreFixBudget = 3 * time.Second`；探针 `:70` 是 `= 50 * time.Millisecond`。⇒ 这条不是"探针被引行号"。
  - ⚠ 同理 `:141`「测试件 `sha1sum` 仍 `b3938cf0…`」＝**产品码**在 `beaeaeba` 时的 sha1（我逐枚核：`beaeaeba` blob sha1 前缀＝`b3938cf07ee8`；产品码现 sha1＝`d558fcfced32`；探针＝`98cbc06f`／`3c5d0d38` 都不等）。⇒ **与本枚探针无关**，改探针字节动不到它。
- `overlay.json` 现读＝`internal/tools/subagent_222_test.go` ← 本枚 ⇒ 可复跑夹具。
- 整理代价：字节 −1／−1（两处各删一个空格）⇒ **20,076→20,074、532 行不变**。

### 枚 8 — `probes/185/c1/mut/fs_broken.go`（唯一非格式枚、唯一让尺退出 2 的枚）
- 现读前几行核对形状〔我现跑〕：`package tools` → `var probe185c1BrokenMarker = undefinedSymbol185c1` → 才 `import (`。**双重坏点**＝① 语法位置（`:4:1`，让 gofumpt **rc=2**）② 未定义符号（让 `go build` rc=1）。
- `gofumpt -d` 现跑＝**rc=2、stdout 空、stderr 单行** `.scratch\wisp\probes\185\c1\mut\fs_broken.go:4:1: imports must appear before other declarations`〔我现跑，把 a1 的转述升成实测〕。
- 归属＝**票 185**（c1 只读普查腿的正控夹具；`overlay-proof.json` 现读＝`internal/tools/fs.go` ← 本枚）。
- **引用面〔我现跑，尺＝`grep -rn --exclude-dir=specs "fs_broken" docs .scratch/wisp/issues`〕＝20 命中**（⚠ **与 a1 报的 13 命中不一致**；`docs/specs` 排除与否都是 20，差在 a1 那次之后台账／票面又追加了文字。以本腿 20 为准）。**具名引行列或错误文案**的逐字如下：
  - `docs/evidence/s1/185-reread-owner-census-c1.md:48`（**防恒绿自证那一格**）：「同一张 overlay 名册换成故意写坏的副本 `mut/fs_broken.go` → `go build -overlay=... ./internal/tools/` **rc=1**（`logs/overlay-proof-broken.txt`：`syntax error: imports must appear before other declarations`，报的就是我这枚副本），真实变异副本 rc=0 ⇒ 映射确实生效，"零枚红"是**读数不是没跑**」
  - 同文件 `:185`（"唯一的非零 rc 是设计如此"）、`:191`（"临时件只建不删：`mut/fs_broken.go`（证明件）… 全部保留"）
  - `docs/evidence/s1/197-subagent-entity-r1b.md:105`（路径＋`4813567e` 09-28 14:42＋那句 `4:1`）、`:288`
  - `docs/evidence/s1/200-project-instructions-r2.md:153`（枚 gofumpt 自己的 parse 报错行）
  - `docs/reports/HANDOVER.md:360`、`:379`、`:409`；`docs/reports/pending-and-issues.md:8983`、`:9566`、`:9610`、`:9614`
  - `.scratch/wisp/issues/236-*.md:14`、`:31`、`:40`、`:54`、`:100`、`:104`
  - **另有 `.scratch/wisp/probes/185/c1/logs/overlay-proof-broken.txt:2`** ＝ `.\.scratch\wisp\probes\185\c1\mut\fs_broken.go:4:1: syntax error: imports must appear before other declarations`
- ⚠ **本腿对乙段代价的一处顶正（新事实，a1 未给）**：那句被引用的错误文本**存在一枚独立的日志文件里**（`overlay-proof-broken.txt`），**不是** `fs_broken.go` 自己的字节。⇒ 格式化 `fs_broken.go` **不会改动那枚日志一个字节**，"已归档凭据"盘上仍然逐字对得上。**真正的风险只在复跑**：重打 `go build -overlay=overlay-proof.json` 会给出新理由（`undefined: undefinedSymbol185c1`）而 rc=1 不变。⇒ 乙段该说的不是"改掉被引用的错误文本"，而是"**改掉那条判据的可复现性**"。这一条对裁定方向有影响（乙的代价被高估了一档）。
- 字节数 `11041`／`11,041` 与 md5 在 `docs`＋`issues` **各 0 命中**〔我现跑〕。

---

## 2. 甲法代价的两个实测数

### 数一：排除 `.scratch/**` 之后被扫的枚数〔我现跑〕
| 尺 | 读数 |
|---|---|
| `git ls-files '*.go' \| wc -l` | **731** |
| `git ls-files '*.go' \| grep -c '^\.scratch/'` | **176** |
| `git ls-files '*.go' \| grep -vc '^\.scratch/'` | **555** |
| `git status --porcelain --untracked-files=all -- '*.go'` 里 `^??` 且不在 `.scratch/` | **0** |

- ⚠ **与 a1／票面的差**：分母是 **731／176／555**，不是 a1 的 **726／171**。5 枚增量是 09-29 之后新入库的 `.go`（同时进了 `.scratch/**` 与总集），所以**非-`.scratch` 那 555 枚恰好仍是票面 `Q-66` 甲段预言的 555**〔我现跑证实该数现在成立，但它是"撞对的"：726−171＝555 与 731−176＝555 两条路都到 555〕。
- 工作树与 tracked 集在 `.scratch` **之外今天完全重合**（非-`.scratch` 未跟踪 `.go`＝0）⇒ 甲法把 step 7 的分母换成"工作树减 `.scratch`"或"tracked 减 `.scratch`"，**今天读数相同**。

### 数二：那 555 枚里 `gofumpt -l` 有没有脏样〔我现跑〕——**没有，一枚都没有**
- 尺＝被授权的静态尺：`git ls-files -z '*.go' \| xargs -0 "$(go env GOPATH)/bin/gofumpt.exe" -l`（gofumpt ＝ `D:\work\base\gopath/bin/gofumpt.exe`，5,028,352 字节，09-23 装）。
- **stdout＝7 行，逐字全部以 `.scratch\wisp\probes\` 开头**（名册见 §3）；stderr＝1 行 `fs_broken.go:4:1`。⇒ **`internal/`、`cmd/`、`tools/`、`design` 外的产品码 555 枚里脏样＝0**。
- ⇒ **对裁定的直接后果**：**甲法确实能让门变绿**（前提是把 step 7 的**正文** `gofumpt -l . tools/d22scan tools/mockllm` 的分母真的换成排除式；只改 `probes/161/r5/attrib.sh:341` 改不到今晚红的那一步——此点为〔转述 236-a1，未复跑〕，本腿未复跑 `gh api`）。乙法同样能让门变绿（8 枚全净后 stdout 空、尺 rc=0）。
- 附带一条会让"绿"被误解的事实：**同一 job 的 step 11 `staticcheck` 今晚也是 `failure`**〔转述 236-a1 的 `gh api` 步级读数，本腿未复跑，且本腿禁碰 CI 读数〕⇒ 修好格式门**不会让整个 lint job 变绿**，只会把 step 9／10 两道 `go vet` 从 `skipped` 放出来。这正是票 236 AC#6 要的读数。

---

## 3. 8 枚之外还有没有下一枚会让门再红 —— **有，7 枚，今天就躺在工作树里（未跟踪）**

### 3.1 tracked 那把尺：名册与 a1 **逐字相同**，无第 9 枚〔我现跑〕
| 比对 | 读数 |
|---|---|
| a1 的 7 枚 | 见 §0 表枚 1–7 |
| 本腿现跑 7 枚 | **逐字相同**（同序、同路径、分隔符同为反斜杠） |
| **多出的枚／消失的枚** | **0／0** |
| 第 8 枚（尺 rc=2 的错误行） | **同一枚**：`fs_broken.go:4:1` |
| xargs rc | **123**（与 a1 同） |

### 3.2 ⛑ **但那只回答了"已入库"的那一半。本腿多跑了一把未跟踪尺，答案变了**〔我现跑〕
- 尺＝`git ls-files --others --exclude-standard -z -- '*.go' \| xargs -0 gofumpt.exe -l`
- 未跟踪 `.go` 总数＝**38 枚**（全部在 `.scratch/**` 下；`.scratch` 之外＝**0**）；gofumpt rc=**0**、stderr 空；**其中脏样＝7 枚**：

| 未跟踪脏样 | 字节 | 行数 | 归属票 | 落盘后台号 |
|---|---|---|---|---|
| `.scratch/wisp/probes/161/r5/negative-control/bad-sample.go` | 1,046 | 24 | 161 | **负控件**（见下） |
| `.scratch/wisp/probes/162/r4/fs_write_m2.go` | 28,618 | 783 | 162 | 突变样本 |
| `.scratch/wisp/probes/226/v1/probe2/main.go` | 5,267 | 185 | 226 | 探针程序（可编译形状） |
| `.scratch/wisp/probes/226/v1/probe3/main.go` | 3,750 | 130 | 226 | 同上 |
| `.scratch/wisp/probes/226/v1/probe4/main.go` | 3,574 | 126 | 226 | 同上 |
| `.scratch/wisp/probes/226/v1/probe7/main.go` | 4,277 | 168 | 226 | 同上 |
| `.scratch/wisp/probes/235/v1/mut/prereadoff-m3_222_test.go` | 24,485 | 614 | 235 | 突变样本 |

- **为什么这 7 枚现在不红、将来一定红**：今晚红的是 `ci.yml:136`＝`gofumpt -l . tools/d22scan tools/mockllm`，它走的是**工作树遍历**、不是 tracked 名册。⇒ 这 7 枚**在本机bench上今天就已经被那把尺看到**（所以 `attrib.sh` 的 form-(B) 形状在本机与在 CI 上不同），只是 **CI 是一次 fresh checkout，未跟踪件不存在** ⇒ 它们在 CI 上隐身。**一旦各自的主人把 `probes/161/r5`、`probes/162/r4`、`probes/226/v1`、`probes/235/v1` 提交，门的名册立刻从 8 涨到 15**。
- ⇒ **对裁定的直接影响（这条会改变裁定）**：**乙法（逐枚整理）今天做完 ≠ 门变稳**，它只清得掉已经入库的 8 枚，**追不上正在飞的 7 枚**（那 4 张票的腿还在跑，`235-v1` 与 `226-v1` 都是"在飞"状态）。**甲法（分母排除 `.scratch/**`）天然把这 7 枚和以后任何一枚台件一起挡在门外**，且 §2 数二已实测"产品码 555 枚脏样＝0"⇒ **甲一次到位、乙要反复追**。
- ⚠ 这正好是**本票 AC#5「在册红名册不是稳定集合」**同一根管子换了一枚读数：**名册不稳定的原因不只是"别人新入库"，还包括"未跟踪件在 CI 隐身、在 bench 可见"**。

### 3.3 负控件一枚（不要顺手清它）〔我现读，尺＝读正文＋读 `attrib.sh`〕
- `.scratch/wisp/probes/161/r5/negative-control/bad-sample.go` 自己的文件头注释逐字写着：「This file is the AC#7 step-2 **NEGATIVE CONTROL** for `.scratch/wisp/probes/161/r4/attrib.sh`。It exists so that one reading of that instrument is **RED**：it is an **unformatted** .go file…」，并列了三个刻意做坏的形状（两空格缩进、`x:=1` 不加空格、闭括号前多一空行）。
- ⚠ **但本腿核到一处能让这枚的风险降级**：`attrib.sh:256-262` 的注释逐字＝「The r4 program proved "the ruler can ring" by leaving a real wound at `.scratch/wisp/probes/999/bad-sample.go`… **The ring now comes from text**, so it is reproducible on a fresh clone, and the physical sample it used to require **has been copied to r5/negative-control/, a path that IS attributable**」。⇒ `attrib.sh --self-test` 的 `st_case`（`:286-298`）喂的是**文本行**，**不读盘上那枚文件的字节** ⇒ **格式化 r5 这枚盘上副本不会让 self-test 失牙**。代价只在"文件头那句 self-declared purpose 会变成假话"这一层。
- ⚠ 顺带一枚反向读数（值得写进裁定的注脚）：**当年那枚真伤口 `.scratch/wisp/probes/999/bad-sample.go` 今天已经在盘上、未跟踪（`git ls-files` ＝ 0），而且 `gofumpt -l` 对它 rc=0、不列＝它已经被净了**（1,467 字节；与 r5 副本自第 1,026 字节／第 20 行起不同）。⇒ 「刻意留一道伤口让尺有牙」这件事**已经被后来的某次整理洗掉过一次了**，这正是乙法形状的实际先例。

### 3.4 台账口径过期再记一次
`docs/reports/pending-and-issues.md:8983` 与 `HANDOVER.md:397` 那句「`-l` 名单确实只有 **4** 枚」今天过期 **4 档**（现量 7＋1）。⚠ 谁按 4 派活都会再低估一次；按台账规矩**只追加更正、不删原句**。

---

## 4. 台件名册（`.scratch/wisp/probes/236/a2/`，只建不删）〔我现跑〕

| 文件 | 作用 |
|---|---|
| `gitstatus-start.txt` | 起手名册（152 行；`?? .scratch/wisp/probes/236/a2/` **已在起手第 101 行**，因为重定向先建了目录） |
| `gitstatus-end.txt` | 终态名册（154 行） |
| `tracked-go-zlist.bin` | `git ls-files -z '*.go'` 原始名册（731 枚） |
| `gofumpt-tracked-stdout.txt` | tracked 尺 `-l` 清单（7 行） |
| `gofumpt-tracked-stderr.txt` | 解析错误行（1 行） |
| `gofumpt-UNTRACKED-stdout.txt` | **未跟踪尺 `-l` 清单（7 行，§3.2 那批）**〔本腿新增，非 a1 台件〕 |
| `gofumpt-UNTRACKED-stderr.txt` | 未跟踪尺 stderr（**空**，rc=0） |
| `diff-1-163-main.txt` … `diff-8-185-fs-broken.txt` | 逐枚 `gofumpt -d` 原文（8 份，脏点的唯一凭据） |
| `census.md` | 本文件 |

**名册漂移归属（不是本腿造成的）**〔我现跑，终态比对尺＝`diff gitstatus-start.txt gitstatus-final.txt`〕：

起手 **152 行** → 交回 **157 行**（HEAD 也从 `98df640a` 走到 `a207d129`）。`diff` ＝ **只 5 行新增，逐条都不是本腿**：

| 新增行 | 归属 |
|---|---|
| ` M .scratch/wisp/issues/62-liquid-glass-ball-visuals.md` | **别的腿**（本腿起手名册里此行不存在、且我全程零触碰 `.scratch/wisp/issues/**`；我只 `head -3` 过 163/174/183/185/197/222/236 七张票的标题行，枚 62 那张我**没读**） |
| `?? .scratch/wisp/probes/132/c3/` | 别的腿 |
| `?? .scratch/wisp/probes/235/v1/` | 别的腿（正在写，见 §3.2） |
| `?? .scratch/wisp/probes/62/` | 别的腿 |
| `?? .scratch/wisp/probes/nail/` | 别的腿 |

- **本腿建的 13 枚件没有让名册相对起手多出一条**：`?? .scratch/wisp/probes/236/a2/` 在**起手第 101 行就已存在**（折叠目录态，因为重定向先建了目录），所以我在它下面建多少枚都不新增行 ⇒ 与 236-a1 遇到的同一条现象，**不是巧合，是判据要按"起手已含"来读**。
- **跟踪文件零改动**〔我现跑〕：`git diff --name-only` 列出的 32 枚（`.gitignore`、`design/**` 一批、`probes/152/my152.py`、`probes/161/r6/logs/**` 九份、`docs/evidence/s1/152-*.md`）**逐枚在起手名册里就已是 `M`/`D`**（我逐枚 grep 过 6 个代表，start 计数全＝final 计数）。**唯一一枚 start=0／final=1 的是 `issues/62-*.md`，不是我。**
- `git diff --cached --name-only` ＝ **空**（本腿从未 `git add`）。零 commit、零 push、零 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`、零删除命令。
- ⚠ 起手名册里本来就有的别人脏改动（含 `design/**` 那一批删除）**一枚没动、一枚没读**（禁区）。

---

## 5. 一句话给裁定的账（本腿实测汇总）

| 裁定向问题 | 本腿实测答案 |
|---|---|
| 乙法（逐枚整理）会不会洗掉别人凭据里的**字节数或 md5**？ | **不会**：8 个字节数＋md5/sha1 全扫 `docs`＋`issues`＝**0 处真引用**〔我现跑〕 |
| 乙法会不会改掉别人凭据里的**行列**？ | **只有枚 7 一处**（`222-…-v1.md:130–132` 引突变件 `:385`／`:400`／`:418`），而**它的修法行数零变化 ⇒ 那三枚行号全部存活**〔我现跑〕；枚 8 被引的 `4:1`＋错误文案存在**另一份日志文件**里，**改 `fs_broken.go` 动不到那枚日志**〔我现跑〕 |
| 乙法真正的凭据代价在哪？ | **枚 5／枚 6**：它们的凭据是"＝`d1c440ad` 的逐字节副本"（`197-…r1b.md:93` 具名声明），我实测两枚与那两枚 blob **字节完全相等**（19,314／25,429）⇒ 格式化即销毁这条可复核性，**并让同句"改前 `-l` 点出 2 枚"的读数复算跑空**。**这一条 a1 没量到** |
| 甲法能不能让门绿？ | **能**：排除 `.scratch/**` 后剩 **555 枚**（731−176），这 555 枚里 `gofumpt -l` 脏样＝**0**〔我现跑〕。⚠ 但必须改 `ci.yml:137–144` 的**正文分母**（红在那一步），只改 `attrib.sh:341` 碰不到——后半句〔转述 236-a1〕 |
| 乙法能不能让门绿？ | 今天能（8 枚全净后 stdout 空、rc=0），**但追不上正在飞的 7 枚未跟踪脏样**（§3.2），且**枚 4 的 md5 在腿回报正文里、枚 5/6 的字节保真会被破** |
| 门绿了 lint job 就绿了吗？ | **不是**：step 11 `staticcheck` 今晚同为 `failure`〔转述 236-a1，本腿未复跑〕。修格式门的**收益是 step 9／10 两道 `go vet` 从 `skipped` 放出来**，这正是 AC#6 要的读数 |

- ⚠ **交回后继续漂**（本腿最后一发现跑，写档时刻）：HEAD 又从 `a207d129` 走到 **`83915ab8`**（其间 3 枚文件、**新增 `.go` ＝ 0**，`git diff --name-only a207d129..83915ab8` 里**没有一枚属于这 8 枚**）⇒ 本 census 的 8 枚字节／行／列读数**到 `83915ab8` 仍然成立**。
- ⚠ 唯一一处会随漂失效的是 §3.2 那 **7 枚未跟踪脏样**的名册（它们在 `235-v1`／`226-v1`／`132-c3`／`62` 那几枚正在飞的腿手里）⇒ **下一腿请重跑 `git ls-files --others --exclude-standard -z -- '*.go' | xargs -0 gofumpt -l`，别引我这一发。**

---

## 6. 没做完／留给编排者（⛔ 非空）

1. **我没有复跑 `gh api` 的步级读数**。§2「step 11 staticcheck 也 failure」、§2「只改 `attrib.sh:341` 碰不到今晚红的那一步」两条都是**〔转述 236-a1，未复跑〕**。⇒ 若要按"改哪个文件能翻颜色"派写腿，**那两条必须重取一发 `gh api` 读数**（本腿禁区是编译/突变，不是 CI 读数，但我 40 枚调用预算里把它让给了逐枚引用面）。
2. **8 枚的脏点我只做到"行号＋形状"，没做"整理后的实测读数"**。我在 §0/§1 给的**新字节数／新行数全部是算术推演**（例如枚 7 ＝ 20,076−2、枚 5 ＝ 434→436 行），**不是我把文件写干净后再量一遍**——那要改文件，本腿禁改。⇒ 写腿若走乙，交回时必须逐枚 `wc -c`／`wc -l` 复量，别引用我这几枚推演值当实测。
3. **枚 5／枚 6 的"乙法新增代价"我只量了字节与来源声明，没量"197-r1b 的改前读数是否还有别的载具"**。`197-subagent-entity-r1b.md:93` 那句「改前＝`gofumpt -l` 点出 2 枚」的原始 stdout 落盘在 `probes/197/r1c/`（该文件 `:124` 具名 `gofumpt-tracked-final.txt` 等）。**我没有逐枚打开 `probes/197/r1c/` 里那些 `base-*/mut*/gate-*` 日志去确认有没有第二处存着这 2 枚的名册**。⇒ 乙法若动这 2 枚，请先把那目录里"含 `pre/` 两枚名册"的日志逐枚 grep 一遍（`grep -rln "r1c/pre" .scratch/wisp/probes/197/`），我没跑这一发。
4. **枚 4 的 md5 在"腿回报正文"里，而腿回报正文不在这仓里**。`185-paged-reread-r1.md:93` 明写 "in the leg reply"。⇒ 编排者若走乙，这条凭据**在盘上无法核也无法破**，只能靠台账追加一句声明它作废；要不要作废属裁定，不属普查。
5. **命中数 20 vs a1 的 13 这处不一致我只归了一半因**。我确认 `--exclude-dir=specs` 加与不加都是 20（所以不是射程口径差），也确认 `98df640a..a207d129` 只动了一枚 235 的裁决表（不是它带来的 7 条增量）。⇒ **剩下那半因（a1 是 22:51 交的、台账在那之后被编排者追加过 `A455` 一带的文字）我只是形状推理，没逐条比对两次名册差集**。要钉死就再跑一次 `grep -rn "fs_broken" docs .scratch/wisp/issues` 与 a1 census 里那 13 条对拉。
6. ~~**"下一枚"我只排了 tracked 两条**~~ ⇒ **这一条本腿已经补成实测**（见 §3.2：未跟踪 38 枚里脏样＝**7 枚**，逐枚有名有号）。**仍未做的是另一半**：我没核对这 7 枚**各归属哪张票的哪一腿、那腿现在是不是还在飞**（我只从路径号推出 161/162/226/235 四张票，没逐张打开工单确认状态行）。⇒ 编排者若按"甲一次到位"裁，这句不影响；若要派腿去清那 7 枚，**必须先确认对应腿交件了没有**（`probes/235/v1/`、`probes/nail/` 在我起手名册里是 `??` 折叠态、终态才展开，说明**这两腿正在写**）。
7. **`docs/specs/**` 我全程排除没读**。⚠ 这里更正我先前写的一句：我现跑三把口径——`grep -rn "fs_broken" docs .scratch/wisp/issues`＝**20**、叠 `--exclude-dir=specs`＝**20**、单独 `grep -rn "fs_broken" docs/specs`＝**0**。⇒ **specs 层对这 8 枚的命中是实测零**，我那个 20 **就是**含 specs 的全数，不必重跑（先前稿子里我把这写成"可能要重跑"，现已改掉）。**但那只验了 `fs_broken` 一枚**；其余 7 枚文件名在 `docs/specs/**` 里的命中数**我没逐枚测**，那 7 个口径仍需重跑才能声称"全仓零引用"。
8. **枚 2 那把 `grep -c "filepath.Clean|filepath.Abs"`＝0 的尺我没复跑**（我只跑了它的结论所依赖的行号与形状）。⇒ 我说"格式中性、复跑不变"是**读 hunk 推的**（本腿的 hunk 只把 `})` 拆行、不新增任何 `filepath` 调用），不是我把那把尺重打了一遍。要钉死请现跑一次。
9. **我没有量"甲法改 `ci.yml:137-144` 正文"具体要改成什么形状**（是 `gofumpt -l $(git ls-files …  \| grep -v '^\.scratch/')`、还是 `-whl` 排除、还是把 `.` 换成显式目录列表）。这属门的形状＝〔契约邻接〕，且 `ci.yml:146-151` 有"不许移动/编辑/加 `if:` 既有步"的先例 ⇒ **本腿只交"改哪一步、改了能不能绿"的实测数，不交改法**。甲法落地前那一形状选型请另派。
10. **本腿全程没读 CI 读数**（`gh api` 一枚没打）。所有"今晚哪一步红/skipped"都是〔转述 236-a1〕。⛔ 这也意味着：**如果 09-29 22:51 之后又有推送过，a1 那套步级结论可能已经不是最新色**，写腿动手前请自己取一发。
