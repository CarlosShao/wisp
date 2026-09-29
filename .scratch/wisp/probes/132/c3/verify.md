# 132-c3 — 对 `132-c2` 六处关键引用的独立取证复核（只读腿）

腿：`132-c3`（只读复核，零产码，零编译，零 `go test`，零 commit／push）。
被复核件：`.scratch/wisp/probes/132/c2/census.md`（35,495 字节／196 行，本腿现量字节数与行数逐字相同）。
票：`.scratch/wisp/issues/132-log-files-are-never-sealed-sealdirdir-has-zero-production-callers.md`（工单号由文件名与第 1 行标题 `# 132 — …` 双向核对，一致）。
本腿起手锚点：`git rev-parse --short HEAD` = **`98df640a`**；落盘期间 HEAD 被并行腿推进三次（`19909091`→`8f59a180`→`b490f8ae`），六条的行号在每一枚新锚点都重跑过，见 §9。⚠ c2 自述取数锚点 `795ed767`、交付锚点 `5759d7fb` **都早于本腿起手锚点**；本腿不沿用 c2 的任何行号，全部现读重取。
起手名册（`git status --porcelain`，本腿未做任何 git 写操作）：工作树本来就有别人的脏改动；本腿起手那把尺只取了前 50 行（缺陷，如实记在 §7 第 4 条），终态判据改按"限定到本腿那一格逐字相同＋全量名册里无一枚出自本腿"给，见 §9。
取数方式：只用 Read / grep / awk-printf 行号 / ls / diff；**没有跑过 `go build`／`go test`／`go vet`／任何编译或突变**。§6 的正控判定＝纯代码推理＋既有日志读数，未真跑突变。

## 0. 骨架（六行表：编号／腿的断言／我的尺／读数／判）

| # | c2 的断言（它写在哪儿） | 我现跑的尺 | 我的读数 | 判 |
|---|---|---|---|---|
| 1 | `open.go` 五处落点行号；自报票面 `:505`→`:506`、承诺注释 `:497-499`→`:497`–`:500` 两处漂 | 整段现读 480–542 行＋`grep -n` 逐句核＋票面现读 | 五处全中；两处漂移与它报的**逐字相同** | **成立** |
| 2 | 调用链在生产路径上（`:328`→`:191`→`run.go:384`／`providers.go:180`）；首跑必写两枚备份 | 逐跳 `grep -n`＋读上下文＋`planMigration` 体＋CI 日志枚举行 | 每一跳都真；`:303` 早退**不**短路"每安装／升级各一次" | **成立**（两处旁证计数漂 1） |
| 3 | 修法三支；丙＝现役 idiom 且零新增导出名；甲峰值＝单文件全量 | `parse.go:185–232` 现读＋`grep -n "func Seal"`＋导出面清点＋`winsec.go:75–145` 现读 | 三支成立；`SealFile` 已有 2 枚生产调用者⇒丙零新增导出名；甲行号全中 | **成立** |
| 4 | 两枚断言把备份目录钉成"恰好两枚、无侧车名"；孤儿态无 fixture；排他创建⇒`Open` 失败⇒rc=2 | 两文件现读（`:110`／`:350` 逐字）＋孤儿 grep＋`run.go:375–395` 现读 | 两行号逐字中；孤儿态**确实**无 fixture；rc=2 文案逐字找到 | **成立**（外加两格它没说的空白覆盖） |
| 5 | `internal/proc` 对 `winsec` 导入数 0⇒常驻腿先跑时日志树**没有可继承的封条** | `grep -rn "CarlosShao/wisp/internal/winsec" internal/proc/`＋`resident_windows.go:40–70` 现读＋MkdirAll 清点 | 导入 0 ✓；`:57`＝`installLogSink` ✓；`:54` 注释逐字 ✓ | **成立**（它的尺文产不出它报的读数，见 §5.2） |
| 6 | 今天零枚尺判"备份文件没加锁"；唯一近尺判结果不判程序；单删 `:503` 换裸 `MkdirAll` 它照绿 | 2 枚 memory 测试件＋`d22scan` 9 条正则＋`sweepPrivate`/`assertPrivateACL`/`aclSIDs` 逐行＋8 份日志比对 | 尺全数复核为中；照绿推理成立；**但它的"现状＝inherited-only"对 `bak-0-1` 为假** | **成立／一处被顶正**（见 §6.4，本腿唯一新增硬事实） |

## 1. 第 1 条：`internal/memory/open.go` 五处落点＋票面行号

尺：`Read internal/memory/open.go` 的 480–542 行段（整段现读，不看 c2 的转述）＋`grep -n` 逐句核＋`wc -l`。
`wc -l internal/memory/open.go` = **542**（c2 自称"542 行文件"✓）。

| 落点 | 本腿现读行号 | 逐字原文（本腿抄，非 c2 抄） | c2 报的 | 票面报的 | 判 |
|---|---|---|---|---|---|
| 目录封条 | `:503` | `	if err := winsec.PrivateDirAll(s.backupDir, 0o700); err != nil {` | `:503` | `:503`（增量段第 3 行） | 一致 |
| 基名 guard | `:506` | `	if _, err := os.Stat(dst); err == nil {` | `:506`（自称票面漂 1） | **`:505`**（硬约束 2） | c2 的"漂 +1"**成立**：票面写 `:505`，今天逐字在 `:506` |
| guard 的早退 | `:507`–`:509` | `s.logger.Warn("migration backup already exists, keeping it",` / `"component", "memory", "backup", dst)` / `return nil` | `:507`–`:509` | 票面未标 | 一致 |
| 侧车循环 | `:511` | `	for _, suffix := range []string{"", "-wal", "-shm"} {` | `:511` | 票面写成"循环里 `:518` 写 `dst+"-wal"`／`dst+"-shm"`" | 一致（措辞层级不同，c2 已自注） |
| 侧车调用点 | `:518` | `		if err := copyFile(src, dst+suffix); err != nil {` | `:518` | `:518` | 一致 |
| `copyFile` 体起 | `:527` | `func copyFile(src, dst string) error {` | `:527` | "体起 `:527`" | 一致 |
| 裸 `os.Create` | `:533` | `	out, err := os.Create(dst)` | `:533` | `:533`（增量段第 1、4 行两处） | 一致 |
| 流式拷贝／收尾 | `:537`／`:541` | `	if _, err := io.Copy(out, in); err != nil {`／`	return out.Close()` | `:537`／`:541` | "io.Copy 流式" | 一致；`:541` 之后**无 sync、无 SealFile**（`:539`–`:542` 只有 `out.Close(); return err` / `}` / `return out.Close()` / `}`） |
| 承诺注释（票面引的那句） | `:497`–`:500` | 见下框 | "实际跨 `:497`–`:500`，承诺句 498–500" | **`:497-499`**（硬约束 2） | c2 的"少一行"**成立** |

`:497`–`:500` 逐字（本腿现抄，四行全量）：

```
// backupDatabase snapshots wisp.db (plus -wal/-shm when present) into
// backup\wisp.db.bak-<from>-<to>. An existing backup is kept (the oldest
// pre-migration snapshot is the most valuable one; a retry after a failed
// migration must not destroy it).
```

票面硬约束 2 引的「An existing backup is kept (the oldest pre-migration snapshot is the most valuable one; a retry after a failed migration must not destroy it)」＝`498`–`500` 三行连读，**逐字为真**；票面把它标成 `:497-499` 是**范围下移一行**（`497` 是函数名那一行，`500` 才是承诺句末行）⇒ c2 的"漂 1 行"两处都复核为真，且**票面这两处至今未改**（票面 append-only，第 74–77 行原文照旧）。

另外两处本腿顺手钉住的：`backupPath` 在 `:493`、名字渲染在 `:494`（`dbFileName`＝`:39` 的 `	dbFileName       = "wisp.db" `）⇒ c2 §1 前两行✓。`backup` 这个目录名只有 `:40` `	backupDirName    = "backup"` 一处定义，`backupDir` 字段只在 `:91`/`:182`/`:494`/`:503` 出现 ⇒ **`:503` 是全仓唯一造出 `backup` 目录的站点**（这条是第 6 条"c2 §6.2 那句『直接删掉 :503 会红，但红因是目录没被造出来』"的前提，本腿确认为真）。

**尺本身复跑（票面增量段第一把尺，逐字照抄题面给的命令）**：

```
$ grep -rn 'os\.Create(\|os\.OpenFile(' --include='*.go' internal/memory internal/secret internal/agent | grep -v '_test.go'
internal/memory/open.go:533:	out, err := os.Create(dst)
$ 同上 | wc -l
1
```

⇒ **恰好 1 处、行号 `:533` 逐字相同**，票面第 66 行的读数在本锚点上仍为真。

**判定：第 1 条＝成立。** c2 报的两处漂移与票面原文互相对得上，本腿不采信它的因果，只确认"它引的那个文件那一行今天就是这么写的"。

## 2. 第 2 条：调用链真不真在生产路径上

尺：`grep -rn "backupDatabase" --include="*.go" .` 全仓现跑＋逐跳 `Read` 上下文＋`grep -rn "memory\.Open("`＋`open.go:340–386`（`planMigration` 体）现读。

**调用点＝恰好 1 枚**（本腿现跑，三行全量，逐字）：

```
./internal/memory/open.go:328:		if err := s.backupDatabase(step.from, step.to); err != nil {
./internal/memory/open.go:497:// backupDatabase snapshots wisp.db (plus -wal/-shm when present) into
./internal/memory/open.go:501:func (s *Store) backupDatabase(from, to int) error {
```

⇒ c2 §2 那句"全仓 3 行命中（`:328` 调用、`:497` 注释、`:501` 定义）⇒ 调用点＝恰好 1 枚"**逐字成立**。

**逐跳（每一跳本腿都读了上下文，不是只 grep 到行号）**：

| 跳 | 现读原文（逐字） | 上下文判定 |
|---|---|---|
| 1 `open.go:328` | `if err := s.backupDatabase(step.from, step.to); err != nil {` → `:329 return err` | 在 `:324 for _, step := range steps {` 循环内；**紧接的上一句是 `:325 if err := boot.Close(); err != nil {`**（c2 §4 用来论证"正常路径下侧车 Stat 落空"的那一跳，本腿确认 `:325` 在 `:328` 之前，中间没有别的语句） |
| 2 `open.go:191` | `if err := s.migrate(); err != nil {` → `:192 return nil, err` | 在 `Open()` 内、`:188` 的 `PrivateDirAll(s.artifactsDir, 0o700)` 之后、`:199 s.reader, err = sql.Open(...)` 之前；**`migrate()` 是无条件调用的**，没有任何 feature flag／分支包住它 |
| 3 `cmd/wisp/run.go:384` | `mem, err := memory.Open(s.dataDir)` | 在 `assembleRuntime`（定义体 `run.go:329 func assembleRuntime(s runSpec) (*agentRuntime, int) {`）内；该函数由 `runTextTask`（定义 `run.go:154 func runTextTask(s runSpec) int {`）在 `run.go:207 rt, code := assembleRuntime(s)` 调；入口 `main.go:145 func cmdRun(args []string) int {` → `main.go:153 return runTextTask(runSpec{` |
| 4 `cmd/wisp/providers.go:180` | `store, err := memory.Open(io_.dataDir)` | 同一函数里 `:170 eps, err := res.ResolveChain()` 与 `:175 prov, err := llm.BuildEndpointProvider(...)` **都成功之后**才打开（c2 写"端点就绪之后才打开"✓）；`:181–:184` 失败 ⇒ `wisp providers: 存储不可用：%v` + `return 2` |

⇒ c2 的四跳**逐跳为真**。生产调用者计数：`grep -rn "memory\.Open(" --include="*.go" internal cmd tools` 排除 `_test.go` ⇒ **只有 `providers.go:180` 与 `run.go:384` 两枚**（c2 §2 末行"非测试调用者＝2 枚"✓）。

**"会不会真跑起来"（这条决定本票是安全洞还是未接线）——本腿独立推的结论与 c2 同向**：

- `internal/memory/schema.go:129` `const SchemaVersionTarget = 2`（逐字）；`schema.go:168`–`:171`：
  ```
  var migrationChain = []migrationStep{
  	{from: 0, to: 1, apply: applyV1},
  	{from: 1, to: 2, apply: applyV2},
  }
  ```
- 生产 target 就是它：`open.go:121` `		target:      SchemaVersionTarget,`（`withSchemaTarget` 在 `:139` 且注释明写"Unexported: the production target is SchemaVersionTarget; only the migration-chain test"用它）。
- 全新库：`open.go:430`–`:432` `	if nTables == 0 {` / `		return 0, true, nil` / `	}` ⇒ `(version=0, fresh=true)`；`open.go:351` `planMigration` 的 `:352 if fresh {` + `:353 version = 0` + `:364 for v < target {` ⇒ **两步**。
- 首跑写两枚备份＝**代码推论成立**，并且**有今天的实跑读数兜着**（不是本腿跑的，是本锚点工作树里既有腿的日志）：`​.scratch/wisp/probes/gate-rerun-1/ac1-baseline-v.txt:615`/`:617` 那两行 `level=INFO msg="pre-migration backup written" ... backup=...\root\data\backup\wisp.db.bak-0-1 from=0 to=1` 与 `...bak-1-2 from=1 to=2`（同一读数的第二枚：`m5a-red.txt:613`/`:615`；CI 侧另有直接证据——`.scratch/wisp/probes/ci-red/logs/run-36559498617-testwindows-job.txt:881`/`:887` 打出这两枚文件的 `icacls` 原文），出自 `TestAC3ProductionDataRootIsPrivateEndToEnd`（真实 `memory.Open(data)`，`production_windows_test.go:25`）⇒ **两枚备份在真实 `Open()` 路径上确实被造出来过**，而且不是本机偶然：GitHub 的 `test-windows` job 里也能翻到同一对。

**题面点名要判的那句："每安装／升级各一次"有没有被 `:303` 的早退短路掉？——没有，而且这正是"各一次"的含义本身。**
现读 `open.go:303`–`:306`：

```
	if version == s.target {
		_ = boot.Close()
		return nil
	}
```

早退条件＝`version == s.target`（2==2），只在**已经迁完的库**上成立；全新库走 `:430` 拿到 `version=0`、v1 库拿到 `version=1`，两者都 `!= 2` ⇒ **不进早退**，直落 `:315 planMigration` ⇒ `:324` 循环 ⇒ `:328`。所以"每安装各一次（fresh：2 枚）、每级升级各一次（v1→v2：1 枚）、稳态零触发"这三句**互不矛盾、同时为真**：c2 §2 把这三句都写了，本腿逐句复核为真。⚠ 唯一能被 `:303` 短路掉的是"每次 `wisp run` 都写"这种读法，而**票面与 c2 都没有这么写**（票面增量段说的是"`backupDatabase` 有生产调用者，落在真实落盘路径上"，没承诺频率）。

**旁证计数（本腿与 c2 差 1 的两处，都不承重）**：
- c2 写"另有 5 行命中在 `.scratch/wisp/probes/**` 的快照件"⇒ 本腿现跑＝**4 行**（`151/run.head.go:313`、`151/run.mut-noclose.go:313`、`151-accept/run.shipped.go:313`、`197/r3b/pre/run.go:347`）。这类计数本来就随并行腿每分钟在动，判**漂 1、不影响结论**。
- c2 写"测试侧引用它的文件＝17 枚"⇒ 本腿：`grep -rln --include="*_test.go"` 全仓＝17，但其中一枚是快照件 `.scratch/wisp/probes/151-accept/accept151_e2e_test.go`；**在树真测试文件＝16 枚**。判**口径差（含/不含 `.scratch`），不是错**。

## 3. 第 3 条：修法是不是三支

尺：`Read internal/config/parse.go` 的 185–232 行＋`Read internal/winsec/winsec.go` 的 75–145 行＋`grep -n "func Seal" internal/winsec/*.go`＋导出面清点（照 c2 给的那把尺逐字复跑）＋`grep -rn "SealFile(" --include="*.go" internal cmd tools | grep -v _test.go`。

### 3.1 丙支（c2 说"票面只写了两支、实际有三支"）——**成立，且是本腿认为最要紧的一条**

票面硬约束 1（票面 `:76`）原文含：**「落点只有两支，代价必须具名：甲＝整档进内存…乙＝要一枚"先封、后流式写"的新入口＝新增导出名＝契约级」**⇒ 票面确实**只写了甲／乙两支**。

仓里那第三支的逐字证据（`internal/config/parse.go`，本腿现读）：

```
196	func atomicWrite(path string, data []byte) error {
201		tmp, err := os.CreateTemp(parent, ".wisp-config-*.tmp")
207		// Seal before a single content byte, not after: os.Chmod(0o600) at the end
215		if err := winsec.SealFile(tmpName); err != nil {
219		if _, err := tmp.Write(data); err != nil {
226		if err := os.Rename(tmpName, path); err != nil {
```

**题面点名要逐字抄的那句注释**＝`parse.go:207` 行首：`	// Seal before a single content byte, not after: os.Chmod(0o600) at the end`
⇒ c2 引的「Seal before a single content byte, not after」**逐字为真**（它是 `:207` 那行到冒号为止的前缀；整条注释块跨 `:207`–`:214`）。c2 给的三个行号 `:201`/`:215`/`:219` **全中**。

⚠ **本腿要补的一处形状差**（c2 没写）：`parse.go` 这枚先例是 **CreateTemp → Seal → Write → Rename**，靠 `:206 defer os.Remove(tmpName)` 兜失败删档、靠 rename 把描述符带过去。把它照搬到 `open.go:533` 的 `copyFile` 时，落地的是"直接 `os.Create` 到最终名 → SealFile → io.Copy"，**rename 那一步没有对应物**，于是"封失败要删档"必须由写腿自己补（见 3.3）。⇒ 三支的说法仍成立，但"丙＝现役 idiom"的**等价度**比 c2 那句"零新增导出名"要低一档：现役 idiom 的失败清理是 `defer os.Remove`，不是 winsec 内部那一条。

### 3.2 "丙真的零新增导出名吗"——**成立**

```
$ grep -n "func Seal" internal/winsec/*.go
internal/winsec/winsec.go:136:func SealFile(path string) error {
internal/winsec/winsec.go:222:func SealDir(path string) error {
（另两枚是 _test.go 里的 SealableTempDirForTest124，不算导出面）

$ grep -rn "SealFile(" --include="*.go" internal cmd tools | grep -v _test.go
internal/config/parse.go:215:	if err := winsec.SealFile(tmpName); err != nil {
internal/secret/migrate.go:174:	} else if serr := winsec.SealFile(backupPath); serr != nil {
internal/winsec/winsec.go:136:func SealFile(path string) error {
```

⇒ `SealFile` **早已导出、且已有 2 枚生产调用者**（config 的"造新档"腿＋secret 的"修旧档"腿，正是 c2 §3 说的那两枚）。丙支用的名字全在现有导出面上 ⇒ **零新增导出名＝真**，票面硬约束 1 里"要造就先回编排者"那句**只在乙支触发，丙支不触发**。

导出面清点（照 c2 的尺逐字复跑，`grep -rn "^func [A-Z]\|^type [A-Z]\|^var [A-Z]" internal/winsec/ --include="*.go" | grep -v _test.go`）＝**15 行**，拆法与 c2 报的完全一致：函数 9 枚（`SetPathResolver`/`PathResolverInstalled`/`ResolvePath`/`PrivateFile`/`PrivateFileExclusive`/`SealFile`/`PrivateDirAll`/`SealDir`/`RemoveUnlinked`）＋sentinel error 3 枚（`ErrUnresolvedPath`/`ErrNotSealable`/`ErrIsReparsePoint`）＋类型 3 枚（`C26Resolver`/`RewriteAccounted`/`ResolvedPath`）。**乙支"新增导出名"这笔账的基数本腿复核为真。**

乙支的另一半——机制已存在但没导出——也**为真**：
`internal/winsec/winsec_windows.go:328` `func sealHandle(f *os.File) error { return applyDescriptor(f.Name(), false) }`、`winsec_other.go:45` 同名同体；`winsec_windows.go:322` 注释逐字「// sealHandle seals a file that is open but still empty. It goes through the」⇒ c2 引的那句**逐字为真**。

### 3.3 甲支的行号与"峰值＝单文件全量"这个后果——**成立**

| 断言 | 现读原文（逐字） | 判 |
|---|---|---|
| 写入口只吃 `[]byte` | `winsec.go:86` `func PrivateFile(path string, data []byte, perm fs.FileMode) error {` | ✓ 票面 `:76` 与 c2 都指这一行，一致 |
| 一次全量 Write | `winsec.go:122` `	if _, err := f.Write(data); err != nil {`（在 `privateFile` 体内，`:108` 定义） | ✓ c2 报 `:122`，一致 |
| flags 与 `os.Create` 同 | `winsec.go:87` `	return privateFile(path, data, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)` | ✓ ⇒ c2 §3"甲不改用户可见行为"的前提为真 |
| 封失败＝不写字节并删档 | `winsec.go:117`–`:121` `	if err := sealHandle(f); err != nil {` / `_ = f.Close()` / `_ = os.Remove(resolved)` / `	return sealError(resolved, err)` | ✓ c2 报的 `:117`–`:121` 一致；`:100`–`:103` doc 逐字写「A seal that fails leaves no content behind: the entry is removed before the error is returned」 |
| 命名二次解析由 `SealFile` 自己做 | `winsec.go:137` `	resolved, err := resolveString(path)`；`winsec_windows.go:300`–`:302`「// SetNamedSecurityInfo takes the object by name, so it also validates that / // the path still resolves to the object we meant (a rename between create / // and here fails instead of sealing the wrong file).」 | ✓ c2 引的句子逐字为真（它标的 `:300`–`:302` 正中） |
| 峰值内存＝单枚文件全量、`wisp.db` 无上限 | `retention.go:38` `	ArtifactsQuotaBytes = 500 * 1024 * 1024`（只管 `artifacts\`）；`:31`–`:36` 只有行级 TTL；包内**没有任何一处给 `wisp.db` 设字节上限** | ✓ 本腿没找到反例；c2 的"结构性无上限"表述**不越证据** |

⇒ **判定：第 3 条＝成立。** 三支的形状、丙支的零新增导出名、甲支的行号与代价，本腿逐字复核为中。
唯一要写的减分：丙支那句"零新增导出名"容易被读成"零代价"，而 c2 自己在同一行末尾确实写了代价（把 winsec 的两条内部保证搬到调用方），**但"config 那枚 idiom 靠 rename＋defer Remove，搬过来要换成最终名直造"这一步它没摊开**——编排者若选丙，这一格要写进切片卡。

## 4. 第 4 条：既有断言、孤儿态 fixture、rc=2 因果链

尺：`Read internal/memory/migrate_v1seed_test.go` 的 95–125 行＋`Read internal/memory/schema_test.go` 的 324–360／382–440 行＋孤儿侧车 grep（照 c2 §7.4 那把尺逐字复跑）＋`Read cmd/wisp/run.go` 的 375–395 行＋`grep -rn "wisp\.db\.bak-0-1 NT AUTHORITY" .scratch/`（拿既有读数验"侧车今天会不会出现"）。

### 4.1 两枚断言的行号与射程——**成立，行号逐字中**

`migrate_v1seed_test.go`（测试名 `:18 func TestMigrateSeededV1DatabaseToV2(t *testing.T)`）：

```
100		// 5. The pre-migration backup for the 1->2 step exists (SPEC-02 §7);
101		// the dir also holds the 0->1 backup from CREATING the v1 database, and
102		// nothing else.
103		if _, err := os.Stat(filepath.Join(dir, "backup", "wisp.db.bak-1-2")); err != nil {
106		entries, err := os.ReadDir(filepath.Join(dir, "backup"))
110		wantBackups := []string{"wisp.db.bak-0-1", "wisp.db.bak-1-2"}
111		if len(entries) != len(wantBackups) {
112			t.Errorf("backup dir = %v, want %v", namesOf(entries), wantBackups)
115			if i < len(entries) && entries[i].Name() != w {
```

`schema_test.go`（`:324 func TestMigrateFreshDatabaseCreatesCurrentVersion(t *testing.T)`）：

```
344		// Reopen must not have written more backups: the two migration backups
345		// (0->1 and 1->2) exist exactly once from the first open.
350		wantBackups := []string{"wisp.db.bak-0-1", "wisp.db.bak-1-2"}
351		if len(entries) != len(wantBackups) {
355		if i < len(entries) && entries[i].Name() != w {
```

⇒ c2 报的 `:110` 与 `:350` **两枚行号逐字命中**，"恰好两枚＋名单就是这两个基名（无侧车名）"这个描述**与代码一致**（`:111`/`:351` 判长度、`:114`–`:118`/`:354`–`:358` 判逐位名字；两枚的 `t.Errorf` 分别在 `:112`/`:352` 与 `:116`/`:356`）。

**"恰好两枚、无侧车名"还有既有实跑读数兜着**（不是本腿跑的）：`.scratch/wisp/probes/gate-rerun-1/ac1-baseline-v.txt` 里 `TestAC3ProductionDataRootIsPrivateEndToEnd` 那枚 sweep 一共枚举了 **8** 项，逐字可数（同件 `:700` `sweep checked 8 entries under ...\root\data, all private`）：
`data`(根本身) + `\wisp.db` + `\secrets` + `\secrets\endtoend` + `\backup` + `\backup\wisp.db.bak-0-1` + `\backup\wisp.db.bak-1-2` + `\artifacts` = 8 ⇒ **`backup\` 里今天确实只有那两枚基名件，`-wal`／`-shm` 侧车一枚都没落**。这从"结果侧"独立证实了 c2 §4 那句"`:325 boot.Close()` 紧接 `:328`，正常路径下侧车 Stat 落空、那两圈直接 continue"。

### 4.2 孤儿态（基名缺而侧车在）今天到底有没有 fixture——**没有，c2 为真；而且比它说的更空**

```
$ grep -rn "bak-0-1-wal\|bak-1-2-wal\|\.bak-.*-wal" --include="*.go" .
（零命中）
$ grep -rln "bak-.*-wal" .scratch/wisp/issues docs
.scratch/wisp/issues/132-…md   .scratch/wisp/issues/89-…md   docs/reports/pending-and-issues.md
```

⇒ 文本层（票面／台账）有，**代码层 0 命中**＝c2 §7.4 为真。
**本腿再加两格它没说的空白**（都直接改 AC#4 的用例设计）：
1. **A 态（基名已存在 ⇒ `:506`–`:509` 早退"keeping it"）也没有 fixture**：`grep -rn "WriteFile\|Create(" --include="*_test.go" internal/memory/ | grep -i "backup\|bak"`＝**零命中**⇒ 全仓没有任何测试往 `backup\` 里预先放过基名件，那枚"保留最旧快照"的承诺**今天一行仪器都没有**。票面硬约束 2 说"这句今天只对基数名成立"——本腿读数更严：**它对基数名成立，但没人测过；对侧车不成立，也没人测过。**
2. **侧车那两圈（`:511` 的 `-wal`/`-shm`）零覆盖**：上面 4.1 的 8 项枚举＋这里的零 fixture ⇒ `copyFile` 被调用的**只有基名一圈**，`dst+"-wal"`／`dst+"-shm"` 这两个名字在任何测试里都没落到过盘。写腿若改这个循环，**咬不到任何现有断言**（c2 说的"最可能被误伤的既有尺"是 `len(entries)!=2`，那枚只在"多出第三枚"时咬，测不出"少拷了侧车"）。

### 4.3 "改成排他创建会把静默覆盖变成 `Open` 失败→rc=2"——**因果链在生产代码里逐环成立**

`winsec.go:90`–`:95` doc 逐字（含 c2 引的那半句）：

```
// PrivateFileExclusive creates path with data and fails with fs.ErrExist if any
// entry already occupies the name. The exclusive create is the point: artifact
// names are injective in the logical id (ticket 79), so a collision means "this
// id's own earlier artifact", and silently truncating it would let one tool call
// overwrite another's evidence. Deciding what a same-id retry means stays the
// caller's job - see agent.Spiller.Prepare's documented last-writer-wins swap.
```

链条（每一环本腿都现读到原文）：
1. `open.go:518` `if err := copyFile(src, dst+suffix); err != nil {` → `:519 return fmt.Errorf("memory: backup %s: %w", filepath.Base(src), err)` ⇒ 排他返回的 `fs.ErrExist` 会被包进 `backupDatabase` 的错误。
2. `open.go:328`–`:329` `if err := s.backupDatabase(...); err != nil { return err }` ⇒ 上抛。
3. `open.go:191`–`:192` `if err := s.migrate(); err != nil { return nil, err }` ⇒ `Open()` 返错。
4. **题面点名要找的那句文案**：`cmd/wisp/run.go:384`–`:387`
   ```
   	mem, err := memory.Open(s.dataDir)
   	if err != nil {
   		fmt.Fprintf(s.stderr, "wisp run: 存储不可用：%v\n", err)
   		return rt, 2
   	}
   ```
   ⇒ `rc=2` **逐字在生产代码里**；另一条腿 `providers.go:180`–`:184` 同形（`wisp providers: 存储不可用：%v` + `return 2`）。

⇒ c2 §4 那句"不是纯洁癖，是一枚低频、但拿可用性换隐私性的行为变更"**本腿不裁因果，只确认它的每一环引的文件那一行今天都这么写**。B 态无 fixture＝c2 自己已标〔纸面风险〕，本腿加 4.2 的两格之后，这条"风险"更准确的写法是：**代码允许、仪器没见过、方向是拒绝启动**。

**顺带把 c2 §4 表里 A 行的一处说法钉住**（它写"早退在写之前，排他根本不会被问到 ⇒ 无变更"）：✓ 成立，`:506` 的 `os.Stat(dst)` 在 `:511` 循环之前，`:509 return nil`。
但**"基名在 ⇒ 侧车也不会被重写"这句要加一个条件**：只有当**基名与侧车同名同代**时才成立；若基名是上一代留下的、这一代 `dst` 不同名（`<from>-<to>` 变了），guard 就不拦。今天 `from/to` 由链条决定，同一对 `(from,to)` 的第二次跑才会撞 guard ⇒ 表述**在本票射程内没错，只是别抄成"只要有一个基名在就全保"**。

## 5. 第 5 条：`internal/proc` 零 winsec 导入 ⇒ 日志树无封条可继承

⚠ **题面要求这一条写成独立一段**，因为它是"要不要单开一张票"的依据。本腿的结论：**成立**，而且比 c2 写的更宽一档。

### 5.1 本腿现跑的尺与读数

```
$ grep -rn "CarlosShao/wisp/internal/winsec" --include="*.go" internal/proc/ | grep -v _test.go
（零命中，rc=1）                      ⇒ internal/proc 对 winsec 的【导入数＝0】：成立

$ grep -rn -A12 "^import (" internal/proc/*.go | grep -i winsec
（零命中，rc=1）                      ⇒ 第二把尺（直接扫 import 块）同读数

$ grep -rn "winsec" internal/proc/
10 行命中，全部在注释里：envfork.go:112/113/116/138/145/148/156/241
        + envfork_test.go:115/147    ⇒ "文字提到、代码不依赖"，与 c2 的"导入 0"不矛盾
```

`internal/proc/` 全目录 17 枚 `.go`（`boot_windows.go`／`doc.go`／`envfork.go`／`shutdown.go`／`singleinstance_windows.go`／`jobscope_windows.go`／`externalsampler_windows.go`／`systemprocs_windows.go`／`treemetrics_windows.go` ＋测试件）；`boot_windows.go` 里 **`MkdirAll`／`PrivateDirAll`／`os.Create`／`WriteFile`／`OpenFile` 五种造文件的动作一个都没有**（`grep -n` rc=1）。

**题面点名要读的那一行**：`cmd/wisp/resident_windows.go:57`
```
57		sink, sinkErr := installLogSink(rt.Layout.DataDir)
```
上下文逐字（`:53`–`:56`，同文件）：
```
53		// Ordering: proc.Boot is the only thing that ran before this, and it seals
54		// nothing - internal/proc has zero winsec imports - so no event can have
55		// been missed. The close defer is registered BEFORE the shutdown defer so
56		// LIFO runs the D38(e) sequence first and its own log lines still land.
```
⇒ c2 §5 末段引的「internal/proc has zero winsec imports」**逐字为真**（在 `:54`），文件 imports 只有 `errors`/`fmt`/`os`/`internal/buildinfo`/`internal/proc`（`:5`–`:12`）——**常驻进程体本身一枚 winsec 名字都不出现**（`grep -n "winsec" cmd/wisp/resident_windows.go` 只命中 `:48`、`:54` 两行注释）。

**落盘那两跳（本腿现读）**：`installLogSink`（`cmd/wisp/logsink.go:144`）→ `:149 p, err := observe.InitLog(observe.LogConfig{Dir: dir, Level: logSinkLevel})` → `internal/observe/logging.go:72` `	if err := os.MkdirAll(cfg.Dir, 0o755); err != nil {` → 滚动文件 `logging.go:244` `	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)`。
⇒ c2 §5 表第 2 行给的两行号（`:72`／`:244`）**逐字命中**；`internal/observe/` 对 winsec 的**导入数也是 0**（`grep -rn "CarlosShao/wisp/internal/winsec" --include="*.go" internal/observe/` rc=1；只有 `logging.go:374/386/387` 三行注释提它，其中 `:386`–`:387` 还写明"不许成环"）。

**`installLogSink` 的生产调用者＝恰好 4 枚**（照 c2 的尺逐字复跑 `grep -rn "installLogSink(" --include="*.go" cmd/ | grep -v _test.go`）：
```
cmd/wisp/logsink.go:144:func installLogSink(dataDir string) (*logSink, error) {   ← 定义
cmd/wisp/models.go:284:	if sink, sinkErr := installLogSink(io_.dataDir); sinkErr != nil {
cmd/wisp/resident_windows.go:57:	sink, sinkErr := installLogSink(rt.Layout.DataDir)
cmd/wisp/run.go:195:	sink, sinkErr := installLogSink(s.dataDir)
cmd/wisp/secret.go:226:	if sink, sinkErr := installLogSink(layout.DataDir); sinkErr != nil {
```
⇒ c2 报的 4 枚＋四行号**全中**（定义行另算，它写的"共 4 枚"是排除定义的调用者数，口径一致）。

### 5.2 一处"尺文与读数不配对"（不翻结论，但按纪律必须点名）

c2 §5 末段写：「`grep -n "memory\.Open\|secret\.NewStore\|PrivateDirAll\|winsec\." cmd/wisp/resident_windows.go` ＝ **只命中两行注释**（`:48`、`:54`…）」。
本腿逐字复跑这把尺：**零命中（rc=1）**——因为该正则要求 `winsec` 后面跟一个点，而 `:48` 是「internal/winsec, a locked config」、`:54` 是「winsec imports」，都没有点。
⇒ **"命中两行"这个读数不是这条命令产出的**（`:48`/`:54` 两行注释的存在本腿用 `grep -n "winsec"` 单独确认了）。
判定：c2 的**结论为真**（resident 腿一枚封条都不下），但**它给出的尺跑不出它给出的读数**——按本仓"别信行号信尺"的规矩，这处必须记名，编排者若照抄那条尺会拿到空结果并误以为"这文件根本不提 winsec"。

### 5.3 本腿把这一条加宽的地方（c2 只说了 `logs`，没说到 `<data>` 本身）

`grep -rn "MkdirAll\|os.Mkdir" --include="*.go" internal/proc/ internal/buildinfo/ | grep -v _test.go` ＝ **零命中**。
⇒ 常驻腿从 `buildinfo.ResolveEnv`＋`proc.Boot` 到 `installLogSink` 之间**没有任何一层建目录**；`<data>` 与 `<data>\logs` 两枚**都由 `logging.go:72` 那一次 `os.MkdirAll(dir, 0o755)` 连根一起造出来**（MkdirAll 建父级时用的是同一个 perm，Windows 上 perm 是装饰、POSIX 上是 0755）。
⇒ 所以 c2 那句"没有可继承的封条"**成立且应当写成更硬的一版**：在"新机器上第一枚命令是常驻腿"这条时序里，**数据根自己就是宽的**（继承的是 `%USERPROFILE%` 的 `BUILTIN\Users` 那一族），`logs` 与 `*.jsonl` 只是跟着宽——不是"数据根封了、只有 logs 没封"。
这条直接把票 132 第 23 行那格（「**数据根目录本身：没锁**」）从 c2 §7.7 判的"按字面读为假"里**救回来一半**：c2 说"数据根已经在 `open.go:176` 用 `PrivateDirAll` 封过"——那只在**有腿调用 `memory.Open`/`secret.NewStore` 之后**为真；常驻腿不满足这个前提。⇒ **两枚说法都对，但都必须带时序条件**，这正是 AC#1 那张表"区分已存在但没锁／新建时才锁"两种形状要落的格。

### 5.4 要不要单开一张票——本腿的意见（不裁，只把依据摊平）

**支持单开的读数（三条都是本腿现跑）：**
1. **射程不同**：票 132 的 AC#2 是"日志目录在创建时 `SealDir`、滚动文件首写前 `SealFile`"，**默认了有人先把 `<data>` 封掉**；上面 5.3 说的是"**建根的那枚 MkdirAll 在哪个进程里跑**"，属**进程拓扑／启动时序**问题，改的是 `internal/observe/logging.go` 与 `cmd/wisp/*.go` 的先后，而不是 `logsink.go` 一处。票面 `:6` 地界写的是 `cmd/wisp/logsink.go`＋`internal/winsec/winsec.go`＋`internal/memory/**`——**`internal/observe/` 不在地界内**。
2. **既有文字已经把这件事当成一条独立规矩在讲**：`cmd/wisp/logsink.go:27`–`:32` 的 invariant #1（逐字含「on the resident leg it means immediately after proc.Boot, which performs no sealing at all (internal/proc has zero winsec imports)」）＋`resident_windows.go:53`–`:56`。两处注释都在**解释为什么这么排**，都没**修**"谁都不封"这件事本身。
3. **它能被独立判红**：`internal/observe` 零 winsec 导入是可 grep 的事实；POSIX 半边 c2 §5 末段已给出（`winsec_other.go:170` `func sealDir(path string) error { return applyDescriptor(path, true) }`，`applyDescriptorPOSIX` `:24`–`:41` 只 `os.Chmod` 单档、**不递归**；`winsec.go:181`–`:217` 的 `missing` 表只覆盖 PrivateDirAll **自己创建的那些层**）⇒ "常驻腿先跑 ⇒ 树宽 ⇒ 后来的 `PrivateDirAll(数据根)` 靠 `sealDir` 的传播走查（`winsec_windows.go:588`/`:592` `propagatePrivate`/`:620`）在 Windows 上**会**把它收窄"这一串，Windows 侧有机制、POSIX 侧没有 ⇒ **两平台的修法不同**，是一张票的体量，不是 132 的顺手一格。

**反对单开（并账）的读数：**
- 票面 AC#2 那句「日志目录在**创建时**就 `SealDir`」本身就要求在"造目录那一句"里改，`observe` 与 `logsink` 谁改都是这一格；`internal/proc has zero winsec imports` 在票面是**同一族的第九次**（票面 `:4` "能力已实现但生产里没人调"），不是新形状。

⇒ **本腿建议（只当依据呈上）**：**单开一张"常驻腿／数据根建立时序"的票，票 132 的 AC#1 现状表里给它一枚具名行并把票面 `:23` 那格改写成带时序条件的句子**。理由＝地界（`internal/observe` 不在 132 的地界里）＋修法跨两平台且不对称＋132 的 AC#2 判据（"创建时就 SealDir"）在常驻腿这条时序下**即使全绿也仍然有宽树**。⛔ 本腿不替编排者拍这一枚；这条**必须进 `docs/reports/pending-and-issues.md` 的 `A##` 或 `Q##`**，因为它是"要不要新增一张票"的决定，不是普查结论。

## 6. 第 6 条：今天没有任何尺判得出备份件没加锁＋正控形状

### 6.1 "直接尺＝0 枚"——**成立**（三把尺全部复跑）

```
$ grep -rln "icacls\|assertPrivate\|PrivateFile\|SealFile\|winsec\." --include="*_test.go" internal/memory/
internal/memory/artifacts_junction_tripwire_windows_test.go
internal/memory/tempdir_resolved_124_test.go            ⇒ 恰好 2 枚：与 c2 一致
$ grep -in "backup" internal/memory/artifacts_junction_tripwire_windows_test.go internal/memory/tempdir_resolved_124_test.go
（零命中，rc=1)                                          ⇒ 这 2 枚【都不判备份件】：成立

$ grep -n "Create\|Seal\|seal" tools/d22scan/main.go
（零命中，rc=1)                                          ⇒ 静态门里连 "seal" 这个词都不出现
```
`tools/d22scan/main.go` 的规则清点（现读 `:141`–`:164`）＝ `var (…)` 块里 **9 条**（`:143` `secretNameRe`／`:144` `wallclockRe`／`:145` `unixTimeRe`／`:146` `timeoutWordRe`／`:147` `mirrorHashRe`／`:148` `mirrorWordRe`／`:149` `approvalPanelRe`／`:150` `artifactToolRe`／`:151` `assignKeyShapeRe`）＋ `:164` 的 `emojiRe` ⇒ c2 报的"`:143`–`:151` 九条 + `:164`"**逐字命中**，无一条与 `os.Create`／封条有关。
⇒ **"没有任何一格会因备份文件没加锁变红"＝成立**。

同族尺在别的包确实存在（c2 §6.3 的三枚，本腿逐枚验行号）：`internal/config/private_acl_windows_test.go:57`（`	backup := filepath.Join(dir, "config.toml.bak-1")`）＋`:85`/`:89`（`	assertPrivate(t, backup)`）、`internal/winsec/migrate_windows_test.go:72`（`func TestAC2MigrationBackupIsPrivate`）与 `:113`（`func TestAC2PreExistingMigrationBackupIsRepaired`）⇒ "memory 包 0 枚、config／winsec／agent 各有"**成立**。

### 6.2 唯一近似尺 `sweepPrivate`：判结果不判程序——**成立**

行号：`internal/winsec/production_windows_test.go:18` `func TestAC3ProductionDataRootIsPrivateEndToEnd(t *testing.T) {` ✓、`:65` `func sweepPrivate(t *testing.T, root string) {` ✓（调用点在 `:59`）、`:29` `	t.Logf("data root contains: %v", mustReadDir(t, data))` ✓。
它**确实递归扫到备份件**：`:92` `				assertPrivateACL(t, full)` 对每个非目录条目执行，`:97` `	t.Logf("sweep checked %d entries under %s, all private", checked, root)`；既有读数（§4.1 那 8 项枚举）里就有 `\backup\wisp.db.bak-0-1` 与 `\bak-1-2` 两行。

**为什么它判不到程序**——本腿逐行读判据，三环：
1. `acl_windows_test.go:86` `aclSIDs` 取的是 `line[:i]`（`i := strings.LastIndex(line, ":(")`）⇒ **只留主体名，`(I)`／`(OI)(CI)` 这些 flags 被丢掉**（c2 特意警告"判据不能沿用它"✓ 逐字为真）。
2. `:119 privateACLError` 判的是**SID 集合**：`s == me` ⇒ `grantsMe`；`systemSID`/`adminsSID`/完整性标签放过；其余进 `bad`。**没有"这条 ACE 是谁写的"这一维**。
3. `:142 assertPrivateACL` 只是 1+2 的包装。
⇒ **只要盘上的有效主体集合是窄集，不管窄集是显式给的还是继承来的，它都绿。**

### 6.3 c2 给的正控形状（只读推理＋引代码，**未真跑突变**）——**结论成立**

题面问：把 `open.go:503` 的目录封条换成裸 `os.MkdirAll` 会不会让 `sweepPrivate` 变红？
**答：不会，照绿。** 推理链（每一环都引代码）：
- `os.MkdirAll(backupDir, 0o700)` 的 mode 参数在 Windows 上不产生任何访问控制效果——同仓把这句话说死在 `acl_windows_test.go:158`–`:159`：「os.MkdirAll's 0o700 is what internal/memory and internal/agent do today; on Windows it does nothing to access control, which is the point」，而 `winsec_windows.go:377` 那条 `inherit = OBJECT_INHERIT_ACE | CONTAINER_INHERIT_ACE` 才是决定子档权限的那一手。
- `backup` 的父级 `data` 在同一次 `Open()` 的 `:176` 被 `PrivateDirAll` 封过，而 winsec 给目录写的是**带继承位**的 ACE：`winsec_windows.go:377` `		inherit = windows.OBJECT_INHERIT_ACE | windows.CONTAINER_INHERIT_ACE` ⇒ 新建的 `backup` 与 `backup\*` 都继承到同一窄集（读数印证：`m5a-red.txt:684`–`:686` 那枚 `bak-1-2` 的三条全带 `(I)(F)`，主体只有 SYSTEM／Administrators／当前用户）。
- 于是 `aclSIDs`/`privateACLError` 拿到的集合与今日**逐字相同** ⇒ 绿。
- 反过来"直接删掉 `:503`"会红，但红因与 ACL 无关：`backup` 目录**只有 `:503` 这一枚创建者**（`grep -rn "backupDir" --include="*.go" internal/ cmd/ | grep -v _test.go` ⇒ `:40`/`:91`/`:182`/`:494`/`:503`，其中只有 `:503` 建目录）⇒ `os.Create` 在未存在的目录里失败 ⇒ `memory.Open` 报错 ⇒ `production_windows_test.go:27` `t.Fatalf("memory.Open: %v", err)`。
⇒ c2 §6.2 那一整格（包括它给的两个反直觉警告）**本腿复核为真**。它引的"要让这枚尺红，得连数据根那层一起拆"有既有读数兜着：`docs/evidence/s1/89-adversarial-acceptance.md:152`–`:160`，`PROTECTED_DACL_SECURITY_INFORMATION` 被摘掉那枚变异里，红名清单**逐字包含** `FAIL TestAC3ProductionDataRootIsPrivateEndToEnd`（`:160`），计数写法是「**8 条顶层 FAIL + 3 条子测试 FAIL = 11 红，2 顶层 PASS**」（`:153`）⇒ c2 说的"11 红里点名 TestAC3"**为真**。

### 6.4 【本腿唯一新增硬事实，方向＝顶正 c2 §6.4】`bak-0-1` 今天**不是** inherited-only

c2 §6.4 的判据设计建在这句上：「**不种 X 时它今天本来就红（现状即 inherited-only）**，所以这一枚同时是负向尺与"改动生效"探针」——它提议的正控＝对 `backup\wisp.db.bak-*` 要求"至少存在一条非 `(I)` 的 ACE"。

**本腿现读数（跨 8 份互不相关的日志，含 GitHub CI）显示：两枚备份件的 flag 状态是分裂的。**

| 来源（都不是本腿跑的） | `bak-0-1` | `bak-1-2` |
|---|---|---|
| `gate-rerun-1/ac1-baseline-v.txt:680` / `:686` | `NT AUTHORITY\SYSTEM:(F)`（**无 `(I)`**，三条皆显式） | `NT AUTHORITY\SYSTEM:(I)(F)`（**全 `(I)`**） |
| `gate-rerun-1/m5a-red.txt:678` / `:684` | 同上 `(F)` | 同上 `(I)(F)` |
| `probes/211/r1/logs/baseline-v-full.txt:5998` / `:6004` | `(F)` | `(I)(F)` |
| `probes/223/r1/ac5-full-final.txt:6325` / `:6331`（＋ `v1/full-v.txt`、`v2/full-v.txt`） | `(F)` | `(I)(F)` |
| `probes/235/orch-rerun/full.txt:6393` / `:6399` | `(F)` | `(I)(F)` |
| `probes/ci-red/logs/run-36559498617-testwindows-job.txt:881` / `:887`（GitHub `test-windows`） | `(F)` | `(I)(F)` |
| `probes/ci-red/logs/run-36500353477-log-failed.txt:881` / `:887`、`run-36556778053-log-failed.txt:901` / `:907` | `(F)` | `(I)(F)` |

⇒ **"非 `(I)` 的 ACE 今天已在场"，就在 `bak-0-1` 上。** 若 AC#4 的负向判据写成"备份目录里**至少一枚**文件带非 `(I)` ACE"，它**今天就绿＝测空气**，注释掉新加的封条也不会红，正控失效；写成"**每一枚**都必须带非 `(I)`"，它今天是红的（`bak-1-2` 不满足），这才是一枚有效的"改动生效"探针。

**机制（本腿代码推论，⛔ 未跑码验证，但与该读数的预测完全吻合）**：`:503` 那枚 `PrivateDirAll` **每次 `backupDatabase` 调用都会跑一遍**，而 Windows 侧的 `sealDir` 不是只封目录自己：
```
winsec_windows.go:588	func sealDir(path string) error {
589		if err := applyDescriptor(path, true); err != nil {
592		return propagatePrivate(path)
595	func propagatePrivate(dir string) error {
620			if err := applyDescriptor(full, false); err != nil {     ← 对目录里【已存在的每一枚文件】下显式描述符
```
⇒ 时序：第 1 步（0→1）造并封 `backup`（当时是空的）→ `os.Create(bak-0-1)` ⇒ 出生时 `(I)`；第 2 步（1→2）**再封一次** `backup` ⇒ 传播走查给 `bak-0-1` **补上显式窄描述符** → 随后 `os.Create(bak-1-2)` ⇒ 出生时 `(I)`。
**同一件事在仓里有第三处文字钉着**：`migrate_windows_test.go:113` 的 `TestAC2PreExistingMigrationBackupIsRepaired` 就是"已存在的旧备份被下一次封条走查收窄"这一族的具名尺，其 `:110`–`:112` 注释逐字写「it has to carry a grant of its own for that leg to be load-bearing (an inherited-only grant would be recomputed by the OS when the parent changes, and the test would pass with the fix deleted)」。

**这条为什么承重（三句，给编排者裁 AC#1／AC#4 用）**：
1. **AC#1 那张现状表不能给备份族写一个统一状态**：同一目录里两枚同类件，一枚"显式（由**下一次**封条的传播走查在事后补的）"、一枚"仅继承"。按票面 `:32` 那句"明写已锁/未锁/继承自父目录"，这里必须**逐件分行**。
2. **c2 §6.4 关于"`icacls` 不带 `(I)`＝被 winsec 直接封过；全带 `(I)`＝只靠继承"的现量对照方法是对的，但用法要反过来说**：不带 `(I)` 不必然意味着"写的时候封了"，也可能是**事后被传播走查补的**——这一格 c2 没区分，照它写的判据会把 `bak-0-1` 读成"这枚站点已经修好了"。
3. **AC#4 的正控必须按"逐枚"落**，否则"注释掉新封条"这一枚变异**不会红**（`bak-0-1` 的显式 ACE 由 `:503` 的传播走查提供，与新加的封条无关）；而**如果**写腿选了甲／丙并把封条挪到 `copyFile` 里、同时**删掉** `:503`，`bak-0-1` 会退回 `(I)` ⇒ 判据能咬。这两句是可以直接写进切片卡的验收形状。
4. **⚠ 它顶正的不仅是 c2，还包括票面自己的一句**（这句本腿认为必须记名）：票面 `:72`「照三行读」的第 ③ 行写「Windows 侧副本自身**无显式描述符**、只靠 `:503` 那枚目录 DACL 的**继承**」。按 6.4 的读数＋机制，这句**只在单步时序（v1→v2 升级，`backupDatabase` 只跑一次）或"最后一次写出的那枚"上为真**；两步首跑时较早的那枚（`bak-0-1`）会被第二次 `:503` 的传播走查**事后补上显式描述符**。⇒ 票面 `:60` 那句"唯一一处写私有数据而不经封条的裸创建"**仍然为真**（落盘那一句确实不经封条），但"权限面比原文件宽"这件事**在结果侧已被仪器间接否掉一半**：票面 `:72` 末句"不许在 AC#1 的表里把它写成已证，只能写〔已证：不经封条〕＋〔待量：宽在哪一格〕"——本腿给的就是那一格的**部分答案**：两枚里的一枚结果并不宽、另一枚只靠继承。

### 6.5 判定

**第 6 条＝成立（c2 报的三把尺、零枚直接尺、`sweepPrivate` 的盲区、单删 `:503` 照绿的正控推理，本腿全部复核为中）；一处被顶正**：它的"现状＝inherited-only"只对 `bak-1-2`、不对 `bak-0-1`，因此它的正控判据要按 6.4 改写为逐枚形式。**该顶正不翻任何一条 c2 的因果结论**，但它翻的是 AC#4 那一格会不会白写。

## 7. 我写错的条目（本节不许为空）

**先给结论**：六条里 **5 条判"成立"、第 6 条判"成立＋一处顶正"**——所以本节不是"全成立"那种情况，它本来就是非空的。但本节的主要职责是记**本腿自己**的错，以下三处都是本腿在本文件写完之前自己发现并当场改掉的，逐条附"我怎么发现的"，编排者可据此判本文件有没有被偷偷修过：

1. **§2 的日志行号抄串了文件**：本腿第一版把"pre-migration backup written"两行写成 `ac1-baseline-v.txt:613`/`:615`——那对行号属于 `m5a-red.txt`（本腿是从 m5a 的读数里记下来的）。改法＝对 `ac1-baseline-v.txt` 重跑一次 `grep -n`，真值是 **`:615`/`:617`**，于是文中同时给出两枚来源并标清各自行号。**这条如果被留在文里，就是本腿犯了 c2 在 §6.2 犯的那枚错（见下面第 4 条）。**
2. **§4.1 的代码块里本腿手抖把变量名抄成 `len(wiskBackups…)` 并配了一句自我说明**：重读现文原文（`	if len(entries) != len(wantBackups) {`）后把该行换成逐字，并删掉那句多余的自我说明。同时把 `:112`/`:352`、`:116`/`:356` 这四枚 `t.Errorf` 行号补做二次核（一次跑 `awk` 打印 95–125／335–360 行，另一次只用 `grep -n` 独立取，两法一致）。
3. **§6.3 有一句把不相关的行号硬接进了因果**：第一版写"`os.MkdirAll` 的 mode 在 Windows 上是装饰——`winsec_windows.go:377` 之外没有任何地方读它"。`winsec_windows.go:377` 讲的是 winsec 自己写继承位，跟"谁读 MkdirAll 的 mode"是两件事，这句话**产不出它暗示的结论**。改后＝mode 无效这件事只引 `acl_windows_test.go:158`–`:159` 那处逐字先例，`:377` 单列为"决定子档权限的另一手"。
4. **本腿把"起手名册"取浅了（一处程序性缺陷，如实报）**：起手那把 `git status --porcelain` 本腿加了 `head -50`，因此**没有拿到完整名册**（落盘前最后一次全量量＝**152 行，其中受版本管理的 M/D＝31 行**；会话中途另一枚量法是 155/34——**这组数在共享树里本来就在动**，并行腿一边 commit 一边把 M/D 变成干净行，见 §9）。补救＝本腿对**自己那一格**跑了限定路径的 `git status --porcelain .scratch/wisp/probes/132/` ⇒ 唯一条目 `?? .scratch/wisp/probes/132/c3/`，再对全部受版本管理的改动逐条核对：**没有一枚是被本腿改的**（名单落在 `.gitignore`、`issues/236`、`probes/152`、`probes/161/r6/logs/flip-*`、`design/**`、`docs/evidence/s1/152-*`，都在本腿零接触范围内；`design/**` 与 `frontend/**` 本腿**一枚都没读**，名册里的名字是 git 元数据、不是文件内容）。⚠ 名册在本腿会话内**还在长**：起手的 50 行里还没有 `.scratch/wisp/issues/236-*.md` 这枚 M，终态有了——它是并行腿在 `98df640a`（＝本腿起手锚点，那枚 ledger commit 正文写的就是 A455 收 132-c2）之后动的工作树，与本腿无关，但**足以说明本腿无法给"与起手名册逐字相同"这一判据**：只能给"限定到本腿那一格逐字相同＋全量名册里无一枚出自本腿"。这一格差是本腿造成的（`head -50`），写在这里不藏。
5. **锚点在会话中被推进了两次，本腿重新取了一次全部行号**：起手 `98df640a`，中途 `git rev-parse --short HEAD` 读到 **`19909091`**，落盘前最后一次读到 **`8f59a180`**（并行腿在 commit）。本腿在两处各跑一次 `git diff --name-only 98df640a..HEAD`：动的 `.go` 全在 `.scratch/wisp/probes/235/v1/mut/` 的突变快照里，**与本腿引用的源文件零交集**；六条用到的关键行号在新锚点全部重跑，**全部原位**（详见 §9 第一行）。

## 8. 没做完／留给编排者（两问必答）

### 8.1 问①：票 132 面上现在一共几格？哪几格的内容已被这轮普查顶正？（给票面行号）

**格数（本腿现读票面全文，逐行定位）**：**5 枚 AC 勾框**——`AC#1`＝`:31`–`:33`、`AC#2`＝`:34`–`:37`、`AC#3`＝`:38`–`:39`、`AC#4`＝`:40`–`:42`、`AC#5`＝`:43`–`:46`；**五枚全是空框**（`- [ ]`），`Progress log`（`:79`–`:81`）只有 2026-09-23 建票那一行。⛔ `docs/evidence/s1/132-*.md` 今天仍是 **0 枚**（`ls docs/evidence/s1/ | grep -c '^132'` ＝ 0），与 c2 §7.8 同读数 ⇒ **本票按"未开工"处理仍为真**。
除 5 枚 AC 外，票面还有 4 个**会被抄进切片卡**的非 AC 单元：原现状表 `:18`–`:24`（5 行）、增量段 `:60`、增量现量表 `:64`–`:70`（5 行）、"照三行读" `:72`、两条硬约束 `:76` 与 `:77`。

**已被这轮普查（c2＋本腿）顶正／需要改写的格，逐格给行号与顶正内容：**

| 票面行 | 原文（节选） | 顶正成什么 | 谁量出来的 |
|---|---|---|---|
| `:22` | 「`cmd/wisp/logsink.go:71` `const logDirName = "logs"`」 | 行号**漂 +5**：今天逐字在 **`logsink.go:76`**（全仓唯一一处 `logDirName =` 定义） | **本腿新发现**（c2 未查这格） |
| `:23` | 「**数据根目录本身：没锁**」 | 必须**带时序条件**：`open.go:176` 的 `PrivateDirAll` 对已存在路径也下封条（`winsec.go:155`–`:156` 逐字），所以**跑过 `memory.Open`/`secret.NewStore` 之后为已锁**；但**常驻腿先跑的新机器上，`<data>` 连根都是 `observe` 裸 `MkdirAll(0o755)` 造的**（`logging.go:72`；`internal/proc`＋`internal/buildinfo` 零 `MkdirAll`）⇒ 这格要写成"依哪条腿先跑"，不能写成恒假也不能写成恒真 | c2 §7.7 提出问题；**本腿 §5.3 把它加宽到数据根本身并给出两平台不对称** |
| `:60`＋`:72` 的③ | 「副本自身**无显式描述符**、只靠 `:503` 那枚目录 DACL 的**继承**」 | **两步首跑时序下为假（半格）**：`bak-0-1` 带**非 `(I)` 的显式窄描述符**，由第二次 `:503` 的传播走查（`winsec_windows.go:588`/`:592`/`:620`）事后补上；单步升级时序／最后写出的那枚仍为仅继承。读数来自 8 份互不相关日志含 GitHub CI（本文件 §6.4 的表） | **本腿新发现**（c2 §6.4 的"现状＝inherited-only"是同方向的半个错） |
| `:76`（硬约束 1） | 「落点**只有两支**」 | 改写为**三支**：甲＝整档进内存（`winsec.go:86`/`:122`，峰值＝单文件全量）／乙＝新导出流式入口（契约级）／**丙＝现役 idiom `CreateTemp`→`SealFile`→才写（`parse.go:201`/`:215`/`:219`，注释 `:207` 逐字），零新增导出名**。⚠ 本腿给丙补一条代价：现役 idiom 的失败清理靠 `defer os.Remove`＋rename，搬到 `open.go` 的最终名直造时**没有 rename 那一步**，"封失败删档"要写腿自己补 | c2 §3／**本腿 §3 逐字复核为中** |
| `:77`（硬约束 2） | 「`open.go:497-499` 的注释…`:505` 的 `os.Stat(dst)`」 | 两处行号各 **+1／+1 行**漂：承诺注释实跨 **`:497`–`:500`**（句体 498–500），guard 实为 **`:506`**；并补一句：`"must not destroy"` 在侧车上不成立**之外**，A 态（基名在 ⇒ 早退）也**零 fixture**、侧车两圈（`-wal`/`-shm`）**零覆盖** | c2 §1/§4；**本腿 §1/§4.1–4.2 复核并加两格** |
| AC#1 `:31`–`:33` | 「每一样私有落点…明写已锁/未锁/继承自父目录」 | 这轮给它加两样：① 增量段要求的 `backup\wisp.db.bak-*` 具名行（`:60` 已写明）；② **同目录同类件要分行**——`bak-0-1`（显式，事后补）与 `bak-1-2`（仅继承）不是一种状态 | c2 §5／**本腿 §6.4** |
| AC#4 `:40`–`:42` | 「把新加的 seal 调用注释掉 ⇒ 必须红」 | 判据形状必须按**逐枚**落（"目录里至少一枚带非 `(I)`"今天已绿＝测空气）；有效形状＝"每一枚 `backup\*` 都带至少一条非 `(I)` ACE"，帮手是 `icaclsRaw`（`acl_windows_test.go:50`）不是 `aclSIDs`（`:86` 丢 flags） | c2 §6.4 提方法；**本腿 §6.4 顶正其适用面** |

**没被这轮动的格**：`AC#3`（`:38`–`:39` 反向对照）与 `AC#5`（`:43`–`:46` 门禁四数）——本腿零跑码，对这两格**无读数**；c2 §7.3 也已声明 POSIX 半边与 AC#5 的 Docker 真跑在它射程外。⇒ **AC#3/AC#5 今天仍是白纸**。

### 8.2 问②：第 5 条那枚"日志树无封条可继承"该不该单开票？

**本腿的意见：该单开，并同时在票 132 的 AC#1 表里给它一枚具名行。**完整依据已经写在 §5.4（正反两栏都在），这里只留三句最硬的：
1. **地界不合**：132 的 `:6` 地界是 `cmd/wisp/logsink.go`＋`internal/winsec/winsec.go`＋`internal/memory/**`；这件事要改的是 **`internal/observe/logging.go:72`／`cmd/wisp/resident_windows.go:57` 的先后与谁来建根**，两枚都不在地界里。
2. **判据不合**：132 的 AC#2（`:34`）问的是"日志目录在创建时封没封"；这件事问的是"**创建它的是哪条腿、那条腿之前有没有任何一层下过封条**"。⇒ 即使 AC#2 全绿，"新机器第一枚命令是常驻腿"这条时序下**仍然留着一棵宽树**（因为该时序里根本没有一次 `PrivateDirAll(数据根)`）。一票两问会把 AC#2 的验收拖成两件事。
3. **修法跨平台不对称**：Windows 侧后来的 `PrivateDirAll` 会靠传播走查**事后收窄**（`winsec_windows.go:588`/`:592`/`:620`，读数见 §6.4）；POSIX 侧 `sealDir`＝单档 `chmod`、不递归（`winsec_other.go:170` → `:24`–`:41`），`privateDirAll` 的下降只覆盖它自己创建的层（`winsec.go:181`–`:217` 的 `missing` 表）⇒ 同一枚 bug 在两平台上**修法不同**，是一枚独立票的体量。
⛔ 本腿不裁这一枚：单开／并账都需要编排者落 `A##` 或摆 `Q##`（票面 `:56` 的规矩）。本腿只把上面三条写成"要不要单开"的可核依据。

### 8.3 本腿没做完的（诚实清单）

1. **零跑码**：本文件所有"会红／会绿／照绿"的判定都是**代码推理＋既有日志读数**，没有一枚是本腿跑的突变或 `go test`（题面硬门）。§6.3/§6.4 的正控形状**必须由写腿或裁决腿用真变异自证一次**才算数。
2. **§6.4 的机制归属是推论**：`(F)`/`(I)(F)` 分裂的**读数**是硬的（8 份日志、含 CI），"**为什么**分裂"＝`:503` 每次调用都跑 `sealDir`→`propagatePrivate` 这条本腿是读码推出来的。反证只需要一步：任何**单步**迁移（`backupDatabase` 只跑一次）里，那唯一一枚应当全是 `(I)`——本腿没有这种读数，写腿或裁决腿跑一次 `wisp run`（fresh 是两步，所以要手工造一个 v1 库）即可钉死。
3. **票面 AC#1 表里 `wisp.db` 那一行（`:24` "未量"）本腿没动**：c2 §5 已把 `wisp.db`/`-wal`/`-shm` 判成"规则内（`:173`–`:175` 注释＋`acl_windows_test.go:332` 那把尺）"，本腿复核该注释逐字为真，但**没量 `icacls`**，那格仍要现量。
4. **`models`／`ball`／`diagnostics` 那三枚同形站点本腿没复核**（不在六条题面内，c2 §5 已登记）。本腿只顺手跑了 `grep -rn "CarlosShao/wisp/internal/winsec" --include="*.go" internal/models/ | grep -v _test.go`＝**零命中**（⇒ c2 那句"`internal/models` 零 winsec 导入"为真），其余（`archive.go:81`、`ball/position.go:73`/`:77`、`diagnostics.go:73` 的调用者计数）本腿**未复跑**，沿用 c2 的〔仅自述〕。
5. **`-wal` 上界与本机 DB 体积无读数**（承 c2 §7.1/§7.2）：`%APPDATA%\wisp` 与 `%APPDATA%\wisp-dev` 本腿没再验，`walPageSizeQuery` 那两枚死常量本腿没再验 ⇒ 甲支代价的"数量级"这一格**两腿都取不到数**。

## 9. 终态自证（git 纪律）

- 起手锚点 `98df640a`；**本文件落盘期间 HEAD 被并行腿推进多次**：`19909091` → `8f59a180` → `b490f8ae` → 交付时 `c28858fc`。各跳之间动的 `.go` 文件本腿逐条看过（`git diff --name-only 98df640a..HEAD | grep '\.go$'`）＝**全在 `.scratch/wisp/probes/235/v1/mut/` 的突变快照里，受版本管理的 `.go` 源文件改动数＝0**（最后一跳现量），**本腿引用的每一枚源文件都不在其中**；六条用到的关键行号（`open.go:176/191/503/506/511/518/527/533`、`schema.go:129`、`parse.go:201/215/219`、`winsec.go:86/136/222`、`run.go:384`、`resident_windows.go:57`）在**每一枚新锚点各重跑一遍，全部原位**（含交付时的 `c28858fc`）。⚠ 起手锚点 `98df640a` 那枚 commit 的正文＝"ledger(A455 收三枚腿 236-a1/**132-c2**/235-r1…)"，也就是编排者写下 A455 的那一枚——本腿正是被派去复跑 A455 所依据的 c2 读数；交付时它仍在 `HEAD` 的祖先里，未被改写。
- **零编译／零 `go test`／零 `go vet`／零突变／零 commit／零 push／零 `git add`**：本腿跑过的命令全集＝`git rev-parse` / `git status --porcelain` / `git diff --name-only` / `git log --oneline -1 -- <path>` / `Read` / `grep` / `awk`（只为打行号）/ `ls` / `wc` / `sed -n`（只为打原文）/ `diff`（只为把 c2 留在树里的**既有**突变快照 `gate-rerun-1/m5a-mutant.go` 与 `winsec_windows.go` 比一下，用来判定 §6.2 那批日志出自哪一枚变异——**没有将任何突变放进工作树，也没有覆盖过任何源文件**）。
- 本腿在树里的**唯一足迹**＝新建 `.scratch/wisp/probes/132/c3/`（本文件一枚），限定路径名册逐字：`git status --porcelain .scratch/wisp/probes/132/` ⇒ **`?? .scratch/wisp/probes/132/c3/`**（c2 那一格 `?? .../c2/` 未被本腿触碰）。
- 未修改／未删除任何既有件；全量受版本管理名册（落盘前最后一次量＝152 行、其中 31 行 M/D；会话中途另一次量＝155/34，这组数在共享树里本来就在动，见 §7 第 4 条）**里没有一枚是被本腿改的**；`frontend/**`、`design/**` **零读零引零转述**（`design/**` 的名字只出现在 git 名册里，本腿没有打开过其中任何一枚）；`PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`、`internal/panel/tokens_fourway_test.go`、`internal/panel/l2_grant_boundary_test.go`、`internal/perm/ticket90_persist_test.go` 本腿**一枚都没读**（`internal/perm/ticket90_persist_test.go` 只在 §2 的 `grep -rln` 输出里以文件名出现，正文未读、未引）。凭据值：本文件零出现。
## 10. 六条各一行（给回报正文用的同文收口）＋本腿自量

| # | 一行判定 |
|---|---|
| 1 | **成立**：五处落点（`:503`/`:506`/`:511`/`:518`/`:527`＋`:533`）本腿逐字命中；票面那两处漂移（`:505`→`:506`、`:497-499`→`:497`–`:500`）c2 报得**不多不少**；那把 grep 现跑仍＝恰好 1 处。 |
| 2 | **成立**：四跳全真（`:328`→`:191`→`run.go:384`／`providers.go:180`），生产调用者 2 枚；`:303` 的早退**只短路稳态**，"每安装／升级各一次"不被短路，两枚备份有 CI 日志兜底；旁证计数漂 1（scratch 快照 4 枚非 5、测试文件 16 枚在树＋1 枚快照）。 |
| 3 | **成立**：**三支**为真，丙支逐字坐实（`parse.go:201`→`:215`→`:219`，注释在 `:207`「Seal before a single content byte, not after」），`SealFile` 已有 2 枚生产调用者 ⇒ **零新增导出名＝真**；甲支行号（`winsec.go:86`/`:87`/`:122`）全中、"峰值＝单文件全量、`wisp.db` 无上限"本腿找不到反例。补一条 c2 没摊开的代价：config 那枚 idiom 靠 rename＋`defer os.Remove`，照搬到最终名直造时"封失败删档"要写腿自补。 |
| 4 | **成立**：`:110` 与 `:350` 两枚行号逐字命中、"恰好两枚／无侧车名"另有既有 sweep 枚举（8 项里 `backup\` 只含那两枚基名件）从结果侧坐实；孤儿态（基名缺／侧车在）**代码层零 fixture**；rc=2 因果链在生产代码里逐环成立，`run.go:386` 那句「wisp run: 存储不可用：%v」＋`return rt, 2` 本腿现读到。**外加两格它没说的空白**：A 态早退也零 fixture、侧车两圈零覆盖。 |
| 5 | **成立**（尺要换）：`internal/proc` 对 winsec **导入数 0**（两把尺同读数），`resident_windows.go:57`＝`installLogSink`、`:54` 注释逐字「internal/proc has zero winsec imports」；**但 c2 写的那把尺（`winsec\.`）跑零命中，产不出它报的"两行注释"**。本腿把它加宽：`internal/proc`＋`internal/buildinfo` 连 `MkdirAll` 都没有 ⇒ 常驻腿先跑时**`<data>` 整棵都是宽的**，不只是 `logs`。→ 该不该单开票：**本腿倾向单开**（地界不合／判据不合／跨平台修法不对称），依据写在 §5.4，裁归编排者。 |
| 6 | **成立一处顶正**：三把尺复跑（memory 包 2 枚 ACL 件都不提 backup、`d22scan` 里连 "seal" 这个词都搜不到、`aclSIDs` 丢 flags ⇒ `sweepPrivate` 判结果不判程序），单删 `:503` 换裸 `MkdirAll` **照绿**的推理本腿逐环确认为真；**但 c2 说的"现状＝inherited-only"只对 `bak-1-2`**——`bak-0-1` 在 8 份互不相关日志（含 GitHub CI）里都带**非 `(I)` 的显式窄描述符**，机制＝`:503` 第二次调用时 `sealDir`→`propagatePrivate`（`winsec_windows.go:588`/`:592`/`:620`）事后补封。⇒ AC#4 的正控必须按**逐枚**写，否则"注释掉新封条"这枚变异不会红；这一处同时顶正票面 `:72` 的③。 |

**本腿自量（最后一次现跑，这句即本文件最后一次改动）**：`wc -l` 本文件＝**473 行**；字节数**故意不写进本文件**——一写进去它就一直变，是自指死结（c2 那枚"35,495 字节／196 行"同样是被别人量出来的，本腿同理：以编排者手里那一把 `wc -c` 为准）。六条中本腿**现读原文**的源文件 18 枚（`open.go`／`schema.go`／`retention.go`／`winsec.go`／`winsec_windows.go`／`winsec_other.go`／`resolve.go` 的导出面／`parse.go`／`secret/migrate.go` 调用行／`logging.go`／`logsink.go`／`run.go`／`providers.go`／`resident_windows.go`／`migrate_v1seed_test.go`／`schema_test.go`／`acl_windows_test.go`／`production_windows_test.go`／`migrate_windows_test.go`／`private_acl_windows_test.go`／`tools/d22scan/main.go`／`internal/proc` 全目录清点），既有读数引用 8 份日志＋1 枚 `docs/evidence/s1/89-adversarial-acceptance.md`；票面 132 全文现读，工单号由文件名＋`:1` 标题双向核对。`docs/evidence/s1/132-*` 本腿现量＝**0 枚**（该目录今天共 320 份）。
