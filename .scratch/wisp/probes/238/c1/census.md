# 238-c1 — 只读普查腿：首启数据根的"谁造的"两形＋常驻尺＋POSIX 半边

- 腿：`238-c1`（**只读普查**，不产码、不跑测试、不占桌面）
- 起手锚点：`git log -1` 现读＝`bf26dfc4d027f85067549018e0d15a89aa4900e2`（Wed Sep 30 11:37:56 2026 +0800）
- 落笔钟：`date` 现读＝`2026-09-30 11:44:47 +0800`（§0 复跑时刻）／末次更正见各行内联时间
- 工单面：`.scratch/wisp/issues/238-first-run-on-a-fresh-machine-creates-the-whole-data-root-with-no-seal-to-inherit.md`
- 本节任何"今天红/绿"的判断**一律来自读断言**（本轮硬红线禁止 `go test`），逐条已具名标注〔未跑包，此条来自读断言〕。

---

## 0. 现量口径（我跑了哪些尺）

⛔ 本轮硬红线：零 `go test`、零容器、零产码。下面每一把尺都是 `git grep` / `sed -n` / `Read`；
`rc` 一律取 **git grep 自己**的退码（不接 `| cat`，接了读到的是 `cat` 的 0——本轮一条这类坑已实测：
`git grep ... | cat; echo $?` 恒等于 0，会"证明"零命中那把尺跑成功）。

### 0.1 票面"现量"三行复跑

**行 1（票面 `:11`）——复跑＝成立，但其中半句按另一种读法是假的，具名更正。**

| 尺（逐字） | 本机读数 |
|---|---|
| `git grep -nE '"github.com/CarlosShao/wisp/internal/winsec"' -- internal/proc internal/buildinfo` | **零命中，rc=1** ⇒ 两包的 winsec **导入语句数＝0**：成立 |
| `git grep -n "internal/winsec" -- internal/proc` | 5 行命中，**全在注释里**（`envfork.go:113/116/148`＋`envfork_test.go:115/147`）⇒ "文字提到、代码不依赖" |
| `git grep -nE "winsec\|Mkdir" -- internal/buildinfo` | **零命中，rc=1** |
| `git grep -n "internal/buildinfo" -- internal/proc` | **5 行命中（rc=0）**，其中产码 2 枚：`boot_windows.go:17`、`envfork.go:10` |

⚠ **具名更正**：票面那一格写「`internal/proc` 对 `internal/winsec` 导入数＝0；**`internal/buildinfo` 同样**」。
这句有两种读法，**只有一种成立**：
- 读法甲＝"`internal/buildinfo` 对 winsec 的导入数也是 0"⇒ **为真**（上面第二／三把尺，连 `winsec` 这个词都在该包零命中）。
- 读法乙＝"`internal/proc` 对 `internal/buildinfo` 的导入数也是 0"⇒ **为假**：`internal/proc` 的两枚产码文件都导入 `internal/buildinfo`。
出处＝`132-c3` §5.3 的原尺是 `grep -rn "MkdirAll\|os.Mkdir" ... internal/proc/ internal/buildinfo/`（比的是**建目录动作**，不是导入数），
票面把它压缩成"同样"两个字时把比较对象丢了。**落地腿按甲读，别照乙写。**

**行 2（票面 `:12`）——复跑＝成立，名册逐枚。** `PrivateDirAll` 生产点＝**恰好 5 枚**（尺：逐文件 `grep -Hn "PrivateDirAll("` over `git ls-files internal cmd | grep -v _test.go`）：

| 落点 | 它封的是 |
|---|---|
| `internal/secret/store.go:49` | `<data>\secrets` |
| `internal/agent/spill.go:111` | spill 目录（`s.dir`） |
| `internal/memory/open.go:176` | **`<data>` 本身**（`abs`＝数据根，注释 `:173`–`:175` 逐字「The data root is sealed before anything is written into it, and sealed **with inheritance**, because wisp.db-wal / wisp.db-shm are created by SQLite and nothing in this repository gets to seal them afterwards」） |
| `internal/memory/open.go:188` | `<data>\artifacts` |
| `internal/memory/open.go:503` | `<data>\backup` |

⇒ **正控形状确实存在，但"数据根自己"只有 `open.go:176` 一枚**；`memory.Open` 的生产调用者**只有 2 枚**
（`cmd/wisp/run.go:384`、`cmd/wisp/providers.go:180`；尺：同法 `grep -Hn "memory\.Open("`）。
⇒ 本票的真正射程比票面标题更宽，见 §0.4。

**行 3（票面 `:13`）——复跑＝成立，行号漂一处，具名更正。**
`BuildDiagnosticsBundle` 生产调用者＝**0 枚**（尺：`git grep -n "BuildDiagnosticsBundle" -- '*.go'` ⇒ 命中只有
定义 `internal/observe/diagnostics.go:61`（注释）／`:62`（func）＋同包 `diagnostics_test.go:35/:125/:140` 三枚测试）。
⚠ **"把日志打进 zip"那段不在 `:73`**：现读在 `internal/observe/diagnostics.go:99`–`:129`
（`:100 entries, lerr := os.ReadDir(o.LogDir)`、`:122 writeZipEntry(zw, "logs/"+name, ...)`）；
`:73` 现读是 `os.Create(o.OutPath)` 那一行。

### 0.2 票面正文引用的两行——逐字命中，无漂

- `cmd/wisp/resident_windows.go:57`：`sink, sinkErr := installLogSink(rt.Layout.DataDir)`（`:53`–`:56` 的时序注释逐字含「proc.Boot is the only thing that ran before this, and it seals nothing - internal/proc has zero winsec imports」）。
- `internal/observe/logging.go:72`：`if err := os.MkdirAll(cfg.Dir, 0o755); err != nil {`。
- 落盘第二跳：`internal/observe/logging.go:244`：`f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)`（滚动日志文件自己**既不经 `MkdirAll` 也不经任何封条**）。

### 0.3 票面没写、但直接决定 AC#1 的一把尺：`SealDir` 的名册

尺：`git grep -n "SealDir" -- internal cmd tools ':!*/*_test.go'`
读数＝**6 行，全部在 `internal/winsec` 包内，且只有 `winsec.go:222` 是函数本身**：

```
internal/winsec/winsec.go:174  注释
internal/winsec/winsec.go:190  注释
internal/winsec/winsec.go:220  注释（doc）
internal/winsec/winsec.go:222  func SealDir(path string) error {     ← 定义
internal/winsec/winsec_windows.go:52  注释
internal/winsec/winsec_windows.go:313 注释
```

⇒ 票 132 标题那句「`SealDir` 有 16 处测试引用、**0 个生产调用者**」**到今天仍未变**（测试侧命中：`git grep -n "SealDir(" -- 'internal/winsec/*_test.go'` 逐枚可数，本腿数到 16 枚调用＋2 枚 POSIX 用例名注释）。
⇒ **`winsec` 包里"封一棵树"这个动作，产码侧唯一的入口是 `PrivateDirAll`，而它只封"自己创建的那些层＋请求点名的那一层"**（`winsec.go:152`–`:157` 逐字：「Ancestors above path are never touched」）。

### 0.4 本腿新量到的四条（票面没有，是 AC#1 的输入）

1. **造根的腿不止常驻腿一枚。** 尺：逐文件 `grep -Hn "installLogSink\|observe.InitLog\|os.MkdirAll"`（排 `_test.go`）。读数：
   - 经 `installLogSink`→`observe.InitLog`→`logging.go:72 MkdirAll` 建树的产码腿＝**4 枚**：
     `cmd/wisp/resident_windows.go:57`、`cmd/wisp/run.go:195`、`cmd/wisp/models.go:284`、`cmd/wisp/secret.go:226`；
   - **绕过 `installLogSink` 直接 `observe.InitLog`** 的腿＝1 枚：`cmd/wisp/slo_windows.go:272`（`slo_other.go:14` 说明这条腿不进 Boot）；
   - 不经任何 sink 就建根的腿＝1 枚：`cmd/wisp/doctor.go:301` `probeWritable` 的 `os.MkdirAll(dir, 0o755)`（调用点 `doctor.go:108`）。
   ⇒ **六枚腿里只有 `run` 那一枚会在后面某一步把根收窄**（`run.go:384 memory.Open` → `open.go:176`）。
   ⚠ 这条把票面标题的"先跑常驻腿"改写得更准确：**"没有一条腿把'根'当成自己的责任"**才是形状；常驻腿只是唯一一条**今天连 `memory.Open` 都不会路过**的腿。
2. **`internal/observe` 对 winsec 的导入数＝0**（尺：`git grep -n "internal/winsec" -- internal/observe` ⇒ 2 行命中，全在 `logging.go:374/:386` 注释里）。
   而 `:386`–`:387` 禁的是**反方向**那枚边（逐字「internal/winsec **must not import this package**（the graph runs observe -> secret -> winsec, so the edge would close a cycle）」）⇒ **"observe 能不能导入 winsec"这句话票面、注释、仓里都没裁过**，见 §4 第 1 条。
3. **依赖图现读**（尺：逐文件 `grep -n "CarlosShao/wisp/internal"`，排 `_test.go`）：
   `internal/winsec` 产码侧**零内部导入**（是一枚叶子）；`internal/proc` 产码导入＝`buildinfo`＋`observe`（`boot_windows.go:17/18`、`envfork.go:10`、`externalsampler_windows.go:10`、`treemetrics_windows.go:12`）；
   `internal/observe/redact.go:10` 导入 `secret`；`internal/secret/store.go:10`＋`migrate.go:14` 导入 `winsec`。
   ⇒ **`internal/proc` 今天就已经通过 `observe -> secret -> winsec` 传递依赖 `winsec`**；给它加一条**直接**边不闭环、不新增可执行依赖，只是把一条已经在图里的边走到明面上（这一条票面 `:17` 点名要量，答案在 §1）。
4. **静态门里没有任何一格与封条有关**：`tools/d22scan/main.go:141`–`:164` 现读 9 条规则（`secretNameRe`/`wallclockRe`/`unixTimeRe`/`timeoutWordRe`/`mirrorHashRe`/`mirrorWordRe`/`approvalPanelRe`/`artifactToolRe`/`assignKeyShapeRe`）＋`emojiRe`，**没有一条认得 `os.MkdirAll`/`Seal`/`Create`**（复认 `132-c3` §6.1 的同一条尺）。⇒ 票面"今天没有任何尺会因此变红"这句**成立**。

### 0.5 已存在的真机 ACL 归属读数（本腿没跑 `icacls`，见 §5；引仓内既有读数）

`docs/evidence/s1/128-ac1-consequences.md:143`–`:153`（票 128 AC#1，本机 `swq`）摆的就是本票 AC#2 那把尺要用的**逐枚 `(I)` 归属**：

| 对象 | 读数（原文摘要） | 性质 |
|---|---|---|
| `…\wisp-dev\secrets` | 六条 ACE **全都没有 `(I)`** | 显式私有＝winsec 自己写的（`PrivateDirAll` 站点） |
| `…\wisp-dev`（**数据根自己**） | 三条**全 `(I)(OI)(CI)(F)`** | 纯继承 ⇒ 建它的不是封条 |
| `…\wisp-dev\logs` | 同上，三条 `(I)` | 纯继承 |
| `…\wisp-dev\config.toml`／`logs\wisp-20260923-001.jsonl` | 三条 `(I)(F)` | 纯继承 |

同文件 `:160`–`:165` 给了"继承源最坏能多宽"的**只读**读数（`icacls C:\Users\Public`：`NT AUTHORITY\INTERACTIVE:(RX,WD,AD)` 等）
⇒ 这是本仓现存**唯一**一组成对的（宽／窄）真机读数，**票 238 的尺能不能咬，就拿它对表**。
⚠ 口径：那棵树是**票 128 修前的"回落到 CWD"**形态（根在启动目录下），不是 `%APPDATA%` 上的首启树；本腿只借它的**归属判据形状**，不把它当"首启现量"。本机 `%APPDATA%\wisp` 现读＝**存在且为空**（`ls` 零条目），`%APPDATA%\wisp-dev` 同样在名册里 ⇒ 这台机器上**没有**一棵真实的宽树可测。

### 0.6 三行读法（本票碰隐私边界，措辞不许重于证据）

1. **现象在哪出现**＝**盘上那棵数据根树的创建时序**（哪条腿先跑、哪枚 `MkdirAll` 建了根），不在任何对外接口上。
2. **有没有本机被入侵的证据**＝**没有**。本腿全程 `git grep`／`sed -n`／`Read`，零写入、零执行被测程序、零 `icacls`；本机 `%APPDATA%\wisp` 至今是空目录。
3. **最坏后果是什么形状**＝**首启由非封条腿建的那棵树的描述符完全由父目录决定**：父目录窄则今天无人能读（本机就是这种），父目录宽（共享卷／漫游档案／`C:\Users\Public` 那种 `(RX,WD,AD)` 形状）则日志与根跟着宽，而 `docs/evidence/s1/128-ac1-consequences.md:149` 证明同一条腿上 `secrets` 是显式窄的 ⇒ **差别的存在已被量到，差别的"宽度"没被量到**。⛔ 定性＝**未证实的暴露面差异**，不写"被窃听/已越权"。

---

## 1. 两形代价表（AC#1）

先定一件不在争议里的事：**"封根"这一手今天只有三种下法**，因为 `winsec` 的产码面就这么多（§0.3）——
`SealDir(path)`（**Windows 侧＝`applyDescriptor` ＋ `propagatePrivate` 全树走查**，`winsec_windows.go:588`–`:620`；POSIX 侧＝单档 `chmod`，`winsec_other.go:170`），
`PrivateDirAll(path)`（`winsec.go:165`→`:177`–`:218`：只封"自己创建的那几层＋点名的那一层"，**祖先一律不碰**，`winsec.go:152`–`:157`），
`SealFile`／`PrivateFile*`（单档）。
⚠ **仓里没有"只封一层、不走查子树"的导出名**（`applyDescriptor` 是包内的，`winsec.go`/`winsec_windows.go` 都没外露）。
⇒ **任何一形如果不想在首启时把整棵已存在的数据根重写一遍描述符，就要新增一枚导出名＝按票 132 `:76` 乙支的先例属契约级，实现腿不许自己造**（这一条要编排者裁，见 §4 第 2 条）。

### 1.1 形①：常驻腿（`internal/proc`）自己调封条

| 项 | 现量 |
|---|---|
| 要动的文件:函数 | `internal/proc/boot_windows.go`：`Boot()`（`:57`）在 `:77`–`:84` 拿到 `layout` 之后插一枚封条；新增导入 `internal/winsec`；或走 `bootConfig`（`:40`–`:43`）＋`BootOption`（`:37`–`:54`，现有两枚 `WithLayout`／`WithRegistry`） |
| 覆盖哪些腿 | **只 2 枚**：`cmd/wisp/resident_windows.go:33`、`cmd/wisp/slo_windows.go:261`（全仓 `proc.Boot` 产码调用者＝这 2 枚，尺见 §0.4）。⚠ `wisp secret` 明写不走 Boot（`cmd/wisp/secret.go:31` 逐字「This command never calls proc.Boot」），`run`／`models`／`providers`／`doctor` 也不走 |
| 漏掉的腿 | §0.4 那 6 枚造根腿里**有 4 枚这条边管不到**（`run.go:195`、`models.go:284`、`secret.go:226`、`doctor.go:301`）⇒ **同一台机器换一条腿首启，根照样宽**，而这枚票的判据要的正是"哪条腿先跑" |
| 跨包依赖边 | **不新开可执行依赖**：`internal/proc` 今天就经 `observe`(`boot_windows.go:18`)→`secret`(`redact.go:10`)→`winsec`(`store.go:10`) 传递依赖 `winsec`；`winsec` 产码侧零内部导入（叶子，§0.4 第 3 条）⇒ 直接边**不可能闭环**。**"给 proc 新增对 winsec 的依赖"这一条本仓允不允许＝没有任何成文禁止**：`AGENTS.md` 禁止清单里没有依赖方向条；仓里唯一写死的同类规矩是反方向的（`logging.go:386`「winsec must not import observe」、`logsink.go:132`、`risk/winsec_c26.go:7`–`:11`「winsec cannot import risk」）。⇒ **此条要编排者裁**（§4 第 1 条） |
| 有没有现成注入模式可照抄 | **有两枚**：① `internal/risk/winsec_c26.go:20` 的 `init()` 把 resolver **装进** winsec（"能力在下游、接线在上游"这一族的样板）；② `proc.BootOption`（`WithLayout`／`WithRegistry`）已是"装配方注入函数／对象"的现成 seam ⇒ 形①′＝由 `cmd/wisp` 传 `WithRootSealer(func(string) error)` 进 Boot，**proc 零新增导入**，代价是多一枚导出 option＋每个 Boot 调用点都得传（漏传＝静默不封，正好是这族票抓了九次的形状） |
| 失败会不会堵死启动 | **会，而且是最直接的一枚**：`resident_windows.go:42`–`:45` 现读是 `if err != nil { "wisp: boot failed"; os.Exit(1) }` ⇒ 封条一旦从 `Boot` 返回错误，**双击图标那条腿直接起不来**。本仓已定案这是要躲的开法：票 128 的"拒绝启动"之所以没堵死联调，靠的是台账 `docs/reports/pending-and-issues.md:3689`–`:3696`（`A111④`）现读的两条前置路径——**① `portable.txt` 挨着 exe 直接返回 `<exeDir>\data` 并结束，② `WISP_ENV=test` 走 `proc.TestDataDir()`，③ 只有前两条都不成立才去问 `%APPDATA%`**（代码＝`cmd/wisp/doctor.go:247`–`:267`）。⚠ **这两条前置路径救不了封条失败**：它们在"**选哪棵树**"这一步之前，不在"**这棵树封不封得上**"这一步之前——便携目录插在一块 FAT/exFAT 的 U 盘上、test 根落在一台只读卷的 CI 上，两条路都会走到"封不上"。⇒ **要么封条失败降级为响亮告警＋照起（照 `resident_windows.go:58`–`:65`／`run.go:196`–`:203` 那两枚既有开法的原话「a log directory that will not open must not become a way to keep Wisp from starting」），要么这一枚要先给 owner 摆代价再裁**（§4 第 3 条） |
| 最坏后果形状 | (a) **半修**：常驻／`slo` 两条时序收窄，另外 4 条照旧 ⇒ AC#2 那把尺若不限定腿，会读成"时好时坏"，很容易被下一位当成仪器坏了；(b) 封条在 `Boot` 里＝**每次启动都走一遍 `propagatePrivate` 全树**，在 Windows 上这意味着 `wisp` 常驻进程每次起都会**把整棵数据根的描述符重写一次**（含 `<data>\models`，见 1.3 的票 95 撞钉），常驻腿是天天起的进程，这把锤子的**重复开销与不可逆面**都要读数（本腿没量，见 §5） |

### 1.2 形②：由装配根在起任何腿之前把根封一次

| 项 | 现量 |
|---|---|
| 装配根在哪 | **今天不存在单一装配根**。6 枚造根腿各有入口：`installLogSink`（4 枚调用点，§0.4）→ `observe.InitLog`（`logsink.go:149`）；`slo` 绕过 `installLogSink` 直连 `observe.InitLog`（`slo_windows.go:272`）；`doctor` 连 sink 都不装，直接 `probeWritable` 的 `os.MkdirAll`（`doctor.go:301`）。⇒ 形② 的"单一点"只有两枚候选：**(甲) `internal/observe/logging.go:72` 那一行**（覆盖 6 枚里的 5 枚，`doctor` 那枚天生不在管子里）；**(乙) `cmd/wisp` 新加一枚 `sealDataRoot(dir)` 并在 6 枚腿各自调用**（覆盖面全，但"单一点"退化成人肉六处，漏一处就是本票的形状） |
| 要动的文件:函数 | 甲形：`internal/observe`（`logging.go:66`–`:75` `InitLogWithRegistry` 的 `os.MkdirAll` 换成／加上 `winsec.PrivateDirAll`），新增 `internal/observe -> internal/winsec` 直接导入；乙形：`cmd/wisp`（新 helper＋`resident_windows.go`／`run.go`／`models.go`／`secret.go`／`slo_windows.go`／`doctor.go` 六处调用点） |
| 跨包依赖边 | 甲形＝**`observe` 新增对 `winsec` 的直接依赖**：`winsec` 是叶子 ⇒ 不闭环；且 `observe` 今天已经经 `secret` 依赖 `winsec`。⚠ 但 `logging.go:386`–`:387` 那段注释把"这条管子不许成环"写成了本包的一处理由，**它禁的是 `winsec -> observe`，没有说 `observe -> winsec` 可以**——仓里**没裁过** ⇒ 与形① 同一枚问题，见 §4 第 1 条。乙形＝**给 `cmd/wisp` 开一枚新直接导入**（现读：产码侧 `cmd/wisp` 里 `internal/winsec` 命中 5 行**全在注释**——`doctor.go:232`、`logsink.go:7`／`:132`、`resident_windows.go:48`、`secret.go:139`；`grep '"github.com/CarlosShao/wisp/internal/winsec"'` 在 `cmd/` 只命中测试件 `secret_dataroot_119b_test.go:49` ⇒ **产码导入数＝0**），这是**本票第一枚真正的新边**，且开在产码里 |
| 失败会不会堵死启动 | 甲形：**四条腿"响亮但不拦"、一条腿直接退 2**——`installLogSink` 失败在 `resident_windows.go:58`–`:65` 与 `run.go:196`–`:203` 都是"打印＋照跑"（原话逐字「a log directory that will not open must not become a way to keep Wisp from starting」／「it does not stop the run」），`models.go:284`／`secret.go:226` 同形；⚠ 但 **`slo` 那枚腿不是**：`cmd/wisp/slo_windows.go:272`–`:278` 现读是 `observe.InitLog` 一失败就 `wisp slo: log pipeline: %v` ＋ `return 2` ⇒ 把封条塞进 `InitLog` 会让**体检／SLO 这条自救腿**新增一种"起不来"，而 owner 恰恰要用它去诊断（这一条与 §3.4 第 3 点的 POSIX 拒绝面叠在一起看）。⇒ **同一枚改动在三条腿上给出两种失败方向，这一格要编排者裁**（§4 第 3 条）。乙形：6 处调用点各裁一次，最容易六处不一致 |
| 最坏后果形状 | 见 1.3（甲乙同命）：**首启那一次 `PrivateDirAll(数据根)`／`SealDir(数据根)` 在 Windows 上＝对已存在整棵子树做 `propagatePrivate`**，把每一枚子档的描述符**换成显式窄集并切断继承**（`winsec_windows.go:595`–`:620`），这不是"给新建物一枚好父亲"，是**对既有树的批量 ACL 改写** |

### 1.3 两形共有的一枚硬撞钉（本腿认为 AC#1 必须先裁它，不然落地腿会当场撞上）

`internal/models/no_seal_ruling_windows_test.go:27` `TestAC3ExtractionIsDeliberatelyNotSealed` 现读逐字（`:13`–`:22`）：

> 「Ticket 95 AC#3's other half: this class is left inherit-wide **ON PURPOSE**, and a test that measures it is what separates that from an unwired site. …
> If somebody later routes archive.go's extraction through winsec, this test goes red. **That is intended**: sealing here is a product decision
> （it is what would stop a second instance running under another account, e.g. a service account, from reusing a GB-sized cache），
> and **decisions like that do not get made by an unrelated hardening sweep**.」

而模型缓存**就在数据根里**（`cmd/wisp/models.go:198` 逐字 `return filepath.Join(dataDir, "models")`）。
⇒ 形①／形② 任何一枚用 `SealDir(根)` 或 `PrivateDirAll(根)`（Windows 都带 `propagatePrivate` 走查），**首启那一次就会把 `<data>\models\**` 的 `BUILTIN\Users` 读权摘掉**——那正是票 95 拍板"故意留着"的那一枚产品决定。
⚠ 具名口径：那枚测试**本身不会因此变红**（它自建 `t.TempDir()` 里的宽父档再解压，不跑生产腿；〔未跑包，此条来自读断言〕），
**变红的是决定本身**：盘上真实数据根的 models 缓存被无关的一轮加固走查收窄，而仓里唯一提醒这件事的文字就在那枚测试的注释里。
⇒ 这一条我判不准该怎么裁，**明写：要编排者裁**（§4 第 2 条），三支候选代价摆在那里：
①封根但**不走查**（⇒ 要新增导出名＝契约级）；②走查但**排除**某些子树（⇒ 要在 winsec 里种一张例外表，本仓对"例外表"的历史包袱见票 106／`private_set_sid`，且这是第二真相源）；③只封 `logs` 那一棵、根不动（⇒ 回到"根是谁造的"仍未答，本票的 AC#2 就无从判）。

### 1.4 还有两枚不能顺手带上的口子

1. ⛔ **别把"封根"写成"根边界"**：本仓已裁过 `allowed_dirs` 是**判级输入、不是执行时硬边界**（批准的 L2 卡能读根外，见 `AGENTS.md` 引的 D4／`PLAN.md` 那条与台账 `Q-60` 的处置）。ACL 收窄描述符**不改变任何一条判级路径**，所以"顺手让它成为硬边界"这件事既不会发生，也不该被写进方案当免费午餐。
2. ⚠ **"路径已解析"这句话说不满**：`PrivateDirAll`／`SealDir` 的参数会过 C26（`winsec.go:158`–`:164`＋`risk/winsec_c26.go:39`–`:45` 的 `Actable()`，票 102 的改写会报**拒绝**而不是报成功）；
   但 `logging.go:72` 那枚 `os.MkdirAll` **什么也不解析**——所以"同一个拼写在两条腿上指的是同一棵树"这句话**今天不成立**，
   首启那棵树是谁造的，可能取决于拼写而不是取决于腿。另：`internal/memory/open.go:170` 在封根前用的是 `filepath.Abs(dir)`（相对拼写会拼到 CWD 上）⇒
   落地腿若新增第二处绝对化，直接撞 `AGENTS.md` §1.2 那条 ban #2（`risk.PathResolver` 之外不许用 `filepath.Clean|Abs` 做文件系统决策）。

### 1.5 两形并排（一行一栏，给编排者裁的时候对着看）

| 判据 | 形①（常驻腿自封） | 形②‑甲（`logging.go:72` 处封） | 形②‑乙（`cmd/wisp` 每腿封一次） |
|---|---|---|---|
| 覆盖造根腿 | **2／6**（resident、slo） | **5／6**（doctor 不在管子里） | 6／6（靠人肉六处） |
| 新增包依赖边 | proc→winsec（传递已有，不闭环；**未裁**） | observe→winsec（传递已有，不闭环；**未裁**） | cmd/wisp→winsec（**产码侧第一枚**） |
| 照抄现成注入模式 | `BootOption`（`WithLayout` 同形）／`risk/winsec_c26.go` 的 `init()` 安装 | 无（改的是包内第一行创建） | `installLogSink` 那一枚"每腿各调一次"的现成形状 |
| 封条失败是否堵启动 | **默认会堵**（`resident_windows.go:42` `os.Exit(1)`），需另裁 | **默认不堵**（两条调用腿已定"响亮但不拦"） | 取决于六处各怎么写＝**最易六处不一致** |
| POSIX 半边 | **碰不到**（`Boot` 是 `//go:build windows`，非 Windows 没有常驻腿：`resident_other.go:23`–`:24` 逐字拒绝启动并 `os.Exit(2)`） | 碰得到（`logging.go:72` 是 untagged；但 `chmod` 单档不递归，见 §3） | 碰得到，同样受 §3 限制 |
| 全树 ACL 改写风险 | **有**（每次启动走查一遍） | **有**（首次建根那一次走查） | **有**（同左） |
| 本票 AC#2 那把尺能不能判它 | 只能判 2 条时序，另 4 条仍红 ⇒ 会被读成仪器坏 | 能判 5 条时序 | 能判 6 条，但漏任一条就整枚失效 |

---

## 2. 一把不恒真的常驻尺（AC#2）

### 2.1 这把尺要判的那一句，物理上只有一种问法

`(I)` 不是"权限宽不宽"的记号，是**这条 ACE 存在哪儿**的记号。仓里已把它说死（`internal/winsec/acl_windows_test.go:383`–`:387` 逐字）：

> 「every child above holds its foreign grant by **inheritance only**, and an inherited ACE is **not stored in the child** - it is recomputed
> from the parent the moment the parent's DACL changes. … A test written against an arbitrary stray file would **pass for the wrong reason and pin nothing**.」

⇒ 三个直接后果，写在这里是为了让落地腿不必再撞一遍：

1. **"根是谁造的"这件事只有在根自己那枚对象上读得出来**：`<data>` 带非 `(I)` 的私有 ACE ⇔ 有封条碰过它；纯 `(I)` ⇔ 从建好到今天没有任何一枚封条碰过它（**也包括"曾经过宽、后来父被收窄"那种**——OS 重算之后它看起来和"生来窄"一模一样）。
2. **子档上问不出创建者**：`wisp.db-wal`/`-shm` 这类"封好的父里后来生的孩子"（`internal/memory/open.go:173`–`:175` 逐字选择**用继承**做机制）与"生来宽的树里没人管的孩子"，文本读数是同一枚 `pure (I)` ⇒ **把它们逐枚判成"必须带非 (I)"＝把本仓选定的机制判成缺陷**。
3. 所以本票 AC#2 那句"判据必须逐枚"，本腿把它落成：**对本次时序里"由这个进程创建的那串目录"逐枚判**（首启常驻腿＝`<data>` 与 `<data>\logs` 两枚；`wisp run` 那形再多一枚 `<data>\secrets` 之类），
   ⛔ **不是**"整棵树每一枚文件都必须带非 `(I)`"。要不要把范围扩到每枚 jsonl 文件＝**这一条要编排者裁**（§4 第 4 条：扩了会把票 132 的 AC#2「每个文件首写前 `SealFile`」变成前置条件，两票就从"串行"变成"互相承重"）。

### 2.2 恒真形状在本票有两种犯法，都点名

- 形状 A（票面 `:18` 点名的那枚）：「至少一枚带非 `(I)` 的 ACE 就算过」——同族事故是 `132-c3` §6.4：`bak-0-1` 今天就在 8 份互不相关的日志里都带非 `(I)`（机制＝`:503` 第二次调用时 `sealDir`→`propagatePrivate` **事后补封**），"至少一枚"那种尺**今天就绿＝测空气**。本票若写"这棵树里至少一枚对象带非 `(I)`"，**同一条顶正原样适用**：`run` 腿跑过一次 `memory.Open` 之后 `<data>` 就是显式窄的，而常驻腿那条时序仍然全树纯 `(I)`——"至少一枚"会**在修好的机器上绿、在有洞的机器上也绿**。
- 形状 B（本腿新点的一枚）：**沿用现有助手就等于沿用它的盲区**。`aclSIDs`（`acl_windows_test.go:86`）取 `line[:i]`（`i := strings.LastIndex(line, ":(")`），**flags 被丢掉**；`privateACLError`（`:119`）只看 SID 集合；这正是 `132-c3` §6.2 说的"判结果不判程序"三环。⇒ 新尺**必须自带一枚逐 ACE 保 flag 的解析器**（`internal/models/acl_sid_121_test.go:50 parseACELines` 是同族里唯一保留整条 flag 文本的解析器，可作起点；⚠ 它在 `internal/models`，跨包复用要么抽公共件、要么在测试内复制并注明出处——本腿不裁这一枚，见 §4 第 5 条）。

### 2.3 判据草案（写成一枚能常驻的行为尺，不写常量名）

**载体**：`cmd/wisp` 里一枚新钉子（与 `resident_sink_nail_127_windows_test.go`／`leg_sink_nail_131_windows_test.go` 同族同形），**驱动真实腿、跑真实进程或真实函数**，⛔ 不是"用 mock 代替真的"。

**判据**（每一步都在本机新建的、**父档未被人为加宽**的临时根上跑；这正是"非恒真"的一半，见 §2.4 的对照）：

1. 取一个空的新根 `R`（`t.TempDir()` 拼接一个尚不存在的 `wisp` 目录名；⛔ 不许碰 `%APPDATA%`，票 132 AC#1 那句"一律用仓外临时根、owner 真实数据目录一字不许多写"同样管本票）。
2. 只跑**首启时序**：常驻腿＝`installLogSink(R)`（或 `runResident` 的真进程形，同 127 那枚钉子的 `buildWispForTest`＋`bootResidentLeg`）；⛔ 中途**不许**调用 `memory.Open`／`secret.NewStore`——一旦调用，后来的封条会把根修成显式窄，那把"曾经谁造的"这一问题**擦掉了**（本票 1.5 表里那枚"事后走查"就是这件事，`132-c3` §6.4 已实测过一次）。
3. 逐枚取"这一步新建出来的目录串"：`R`、`R\logs`。对每一枚：
 - 读 `icacls`，**保留 flag**，筛出"主体 SID ∈ {当前用户, SYSTEM, Administrators}"的那几行；
 - **要求其中至少一条不带 `(I)`**（＝这条私有授权是这个对象**自己的**，是某枚封条写上去的）；
 - 同时要求**不带任何外来主体的授权**（沿用 `privateACLError` 那句"foreign SID"判语，或它的等价），否则"被谁封的"没意义。
 任一枚不满足 ⇒ 红，红名点到**那一枚对象的路径**，不点整棵树（逐枚的红名才读得出是谁造的）。
4. **对照枚（同一次跑里必须同时在场，否则这枚尺自己就是恒真的）**：对 `R\secrets`（由 `secret.NewStore`→`PrivateDirAll` 造，**已知带非 `(I)`**）跑**同一条判据**，要求**绿**。
 既有真机读数对表：`docs/evidence/s1/128-ac1-consequences.md:149`（`secrets` 六条**全不带 `(I)`**）vs `:150`–`:153`（`wisp-dev`、`logs`、`config.toml`、`jsonl` **全带 `(I)`**）⇒ 判据的两侧各有实测锚，不是推理出来的。

### 2.4 自带正控（票面 `:18` 要求的"把新封条注释掉必红"）

| 变异 | 必须红的是哪一枚 | 为什么红得起来 |
|---|---|---|
| **X1＝把新加的封条那一行注释掉** | §2.3 那枚新钉子的**第 3 步**（`R`、`R\logs` 逐枚） | 注释掉之后建根只剩 `logging.go:72` 那枚 `os.MkdirAll`，`MkdirAll` 在 Windows 上不写任何描述符（本仓原话：`acl_windows_test.go:158`–`:159`「os.MkdirAll's 0o700 … on Windows it does nothing to access control, which is the point」），也没有 `PROTECTED_DACL` 那一手（`winsec_windows.go:377` 的 `inherit = OBJECT_INHERIT_ACE|CONTAINER_INHERIT_ACE` 才是决定子档权限的那一手）⇒ 根**只剩 `(I)`** ⇒ 第 3 步必红 |
| **X2＝把封条挪到"装完 listener 之后"** | 同一枚钉子的第 3 步**不会红**（根最终是显式窄） | 这一发专门用来证明 §2.3 判的是**存在性**不是**时序**；所以钉子必须另配一发**时序半**：在第 2 步之后、封条之前插一次读数（同形做法仓里有现成参照——`logsink_windows_test.go:308 TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite` 就是把"先后"当独立一发来判的）。**这一半要不要写进判据＝要编排者裁**（§4 第 4 条） |
| **X3＝把 §2.3 的第 4 步对照枚删掉** | 新钉子自己 | 删掉对照枚之后，任何"整棵树的 `(I)` 状态读错了"的走法都不会被发现；本仓对负向尺的规矩就是"必配种 X 必响的正控" |

⚠ **不许写成"至少一枚"**（§2.2 形状 A）；⛔ **不许把判据降级成 grep 新符号名**——第 64 条纪律：行为型钉只写产物不写常量名，`d22scan` 也没有任何一条认得封条（§0.4 第 4 条），所以**唯一能常驻的是行为读数尺**。

### 2.5 今天这枚尺是红的还是绿的——具名口径

- **本腿不跑测试**（硬红线）。按读断言＋读代码给结论：§2.3 在第 3 步上**今天必红**，因为产码侧没有任何一条腿对 `<data>` 或 `<data>\logs` 下过封条（§0.3：`SealDir` 产码调用者 0；§0.4 第 1 条：6 枚造根腿里没有一枚先封根；唯一会收窄根的 `open.go:176` 只有 `run.go:384`/`providers.go:180` 两枚调用者，且都在 sink 之后）。
- **"算未修码"的锚点**（本仓规矩：派单写"这枚判据今天应该红"必须自带哪枚 commit 算未修码＋自证尺）：**起手锚点 `bf26dfc`（本腿现读 `git log -1`）＝未修码**；自证尺两条＝§0.3 的 `git grep -n "SealDir" -- internal cmd tools ':!*/*_test.go'`（6 行全在 winsec 包内）＋§0.4 第 1 条的造根名册。⇒ 落地腿必须**钉子与封条同枚 commit 落地**，否则门禁当天就多一条红。

### 2.6 名册：今天绿着、会被这把新尺／这枚新封条打红的用例（逐枚读断言得来）

⛔ 全表判语一律标〔未跑包，此条来自读断言〕（本轮禁 `go test`）。分两类：**R 类＝被"尺"打红**（只有在新尺改写**共享助手**时才发生，草案 §2.3 特意自带解析器，就是为了让这一类不发生）；**F 类＝被"封条落地"打红**（无论选哪一形都要正面处理，属派单里"要重写就写清被删断言的存在理由"的那一档）。

| 用例（文件:行:名） | 类 | 今天为什么绿 | 会不会红／为什么 | 本腿判 |
|---|---|---|---|---|
| `internal/winsec/acl_windows_test.go:332 TestAC2SealedDirCoversFilesItNeverTouched`（`:355`/`:371`/`:372`/`:378` 四处 `assertPrivateACL`） | R | 它测的就是"**winsec 从没碰过的孩子靠继承也该窄**"（`:349`–`:350` 逐字「Files winsec never touched, created inside the sealed root, must inherit privacy from it - that is the whole coverage argument」） | **若把 `(I)` 判据塞进共享助手 `assertPrivateACL`／`privateACLError` ⇒ 必红**，而且红得没道理（判的是本仓选定机制） | ⛔ 禁止复用那两枚助手；新尺自带解析器 |
| `internal/winsec/production_windows_test.go:18 TestAC3ProductionDataRootIsPrivateEndToEnd`（`:42`–`:52` 判 SQLite 写的 `wisp.db`/`-wal`/`-shm`；`:59 sweepPrivate`） | R | 同上：`-wal`/`-shm` **只可能靠继承**（`open.go:173`–`:175` 的原话） | 同上一条，塞进共享助手 ⇒ 必红 | ⛔ 同上 |
| `internal/winsec/migrate_windows_test.go:72`、`:113`（`:104 sweepPrivate(secrets)`；`:42` 注释逐字「The parent is not winsec's to seal」） | R/F | 它同时钉住"**父目录不归 winsec 封**" | 塞共享助手 ⇒ 红（父档纯 `(I)` 是**故意的**）；形② 若封到根以上一层 ⇒ 直接违背 `:42` 那句话 | 具名：`winsec.go:152`–`:157`「Ancestors above path are never touched」是既有裁决，新尺射程不许越过它 |
| `cmd/wisp/logsink_windows_test.go:137 TestAC2SealNoticeLandsInTheRunLegLogFile` | F | 它数的是 `hits`（`msg == sealNoticeMsg` 的记录）**恰好 1 条**（`:174`），且 `:181` 要求那条的路径等于 `secrets`（`assertNamesTree117` 比叶子名） | 新封条若在 `secret.NewStore` **之前**封了根，而 `run` 腿这枚用例已把 `R\secrets` 用 `icacls … /grant *S-1-1-0:(OI)(CI)(RX)` 加宽（`:146`）⇒ 封根那次走查会**先把 Everyone 摘掉**（`winsec_windows.go:595`–`:620`），到 `PrivateDirAll(secrets)` 时无东西可清 ⇒ `len(hits)!=1` 直接 `t.Fatalf`；即便仍有 1 条，那条点名的是**根**、叶子名不等 ⇒ `assertNamesTree117` 红 | **必被红**（形②‑甲/乙 在 `installLogSink` 之前封根这一形）；处置＝要么封根放在"只封不 propagation"的形态（需新导出名，§1 开头），要么这枚用例的判语要重写并写清存在理由 |
| `cmd/wisp/logsink_windows_test.go:308 TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite` | F | 它钉 `first == 1`（只有 install 记录可以在 seal 通知之前） | 封条早于 listener ⇒ 通知进票 130 的早缓冲、被 `FlushEarlyLogRecords`（`logsink.go:173`）排在**记录 0** ⇒ `first` 变 0；封条插在 `InitLog` 之后、booking 之前 ⇒ 通知落文件仍在 booking 之前，同样 `first!=1` | **必被红**，除非新封条放在 install booking **之后**（那是 §2.4 的 X2，把窗口留在生产里） |
| `cmd/wisp/resident_sink_nail_127_windows_test.go:393`（断言在 `:423`–`:436`：记录 0 必须＝早缓冲回放的裁决、`:434` 「install booking is record %d, want 1: the replayed early records come first **and nothing else**」）＋同文件 `:543`（`:587`–`:588` 再钉一次记录 0） | F | 票 130 把"记录 0 是谁"钉成了两枚钉子 | 同上：常驻腿的封条若早于 `installLogSink`，回放名册多一条 ⇒ 两枚都红；⚠ **且红不红取决于这台机器的临时根父档给不给 `BUILTIN\Users`**（`132-c3` §6.3 记过「ticket 95 measured hosts where %TEMP%'s parent grants BUILTIN\Users」）⇒ 是一枚**主机相关**的红，最难查的那种 | 派单里必须点名这两枚；本腿判：**要么显式规定"封条不许发通知到早缓冲"，要么改这两枚的索引判语并重写其存在理由**（不许放宽成 `>=1`） |
| `internal/models/no_seal_ruling_windows_test.go:27 TestAC3ExtractionIsDeliberatelyNotSealed`（注释 `:13`–`:22`） | F（决定层） | 它自建宽父档＋本地解压，不走生产腿 | 用例本身不会红〔未跑包，此条来自读断言〕；**红的是票 95 的决定**：封根的全树走查会在真机把 `<data>\models\**` 的 `BUILTIN\Users` 读权摘掉（`models.go:198`） | ⚠ 见 §1.3／§4 第 2 条 |
| `cmd/wisp/resident_sink_nail_127_windows_test.go:498 TestAC1ResidentLegOutlivesItsOwnLogFailure` | F | 它要求"日志目录开不出来时**照起**"，并要求那枚挡路的普通文件**还在**（`:527`–`:529`） | 只要封条失败仍走 `resident_windows.go:58`–`:65` 那种"响亮但不拦"就**不红**；若形① 把封条塞进 `Boot` 并 `return err` ⇒ 这枚钉子（和 `:393`／`:543` 两枚）会一起红，红因＝启动被堵死 | 具名：这一行就是"失败方向"的仪器代价，编排者裁 §4 第 3 条时对着它看 |
| `internal/winsec/reparse_windows_test.go:103/:170/:209/:327`（`TestAC4…IsNotRecursed/IsNotWalked/SealedWalkSkipsLinks`）、`internal/winsec/private_other_test.go:60 TestPOSIXSymlinkAtArtifactPositionIsNotRecursed` | F | 钉的是"走查**永不穿链接**"与 POSIX 侧**故意不递归** | 若为补 POSIX 半边而给 `sealDir` 加递归（§3 的丙支），**POSIX 这枚不会红**（它测的是 `RemoveUnlinked`），但 `winsec_other.go:165`–`:169` 那段**存在理由**会被推翻，而那句写的是「It deliberately does not walk the existing subtree … would be a chmod over files whose permissions somebody else set on purpose」 | ⚠ POSIX 递归＝要重写这段注释＋给新钉子的存在理由，属裁量项（§4 第 6 条） |

---

## 3. POSIX 半边具名（AC#3）

### 3.1 先纠正一个会被顺嘴说错的对称句

票面 `:4` 那句「POSIX 只 `chmod` 单档、不递归」是**真的**，但它推不出"POSIX 那半边和 Windows 那半边是同一件事的两套写法"。现读两条：

- `internal/winsec/winsec_other.go:165`–`:170` 逐字：
 「sealDir narrows a directory. It **deliberately does not walk** the existing subtree the way the Windows implementation does:
 on POSIX a child **never inherited its parent's mode in the first place**, so 'seal the tree' here would be a chmod over files
 whose permissions somebody else set on purpose, and **the hole being closed does not exist on this platform**.」
- 同文件 `:24`–`:41 applyDescriptorPOSIX`：`os.Chmod` 之后**回读 `info.Mode().Perm()` 并要求逐字等于 0600/0700**，读不回来就 `ErrNotSealable`（FAT/exFAT/`all_squash` 的 NFS 会在这里变红，不是跳过）。

⇒ **本票的"根"这一格在 POSIX 上其实只要单档 `chmod` 就完备**（根与 `logs` 两枚目录各自 0700），
**"递归"这一支在 POSIX 上不属于本票的必需项**；真正在 POSIX 上漏的是**每枚 `*.jsonl` 的 mode**——
那是 `logging.go:244` 用 `os.OpenFile(name, O_CREATE|O_WRONLY|O_APPEND, 0o644)` 建的，
⚠ **POSIX 上没有继承这回事，所以"封了根"不会让日志文件变窄**，只有"每枚文件首写前 `SealFile`"才会＝**票 132 的 AC#2，不是本票的射程**。
⇒ 这条边界要在判据注释里逐字写出来（§3.3 给了文案草案），否则下一位会把"根封了"读成"日志不泄露了"。

### 3.2 两支走法的代价（票面 `:19` 要的两选一）

| 走法 | 内容 | 代价（具名） |
|---|---|---|
| **甲：不递归，只封"本时序里由这个进程创建的那串目录"**（本腿推荐它作为 AC#3 的答案，⛔ 但"要不要按甲写"这一条仍要编排者裁，§4 第 6 条） | `PrivateDirAll(R)`／`SealDir(R)` 各一枚单档动作，Windows 侧靠 `sealDir` 的走查顺带收窄既有子树（副作用＝§1.3 的 models 撞钉），POSIX 侧就是 `chmod 0700` | ① POSIX 的"已有子树照旧宽"＝**语义不对称是刻意的、有注释在先**（`winsec_other.go:165`–`:169`），不写进判据注释就会被读成漏了；② 每枚 jsonl 的 mode **仍不受本票管**（§3.1） |
| **乙：给 POSIX 补一枚递归走查**（把 Windows 的 `propagatePrivate` 在 `!windows` 上也做一遍） | 新增平台半边：`winsec_other.go` 里写一棵 walk，或把 `sealDir` 改成两平台都走 | ① **正面推翻 `winsec_other.go:165`–`:169` 的存在理由**（那句写的是"这儿的洞不存在"＋"走查＝覆盖别人有意设的 mode"）⇒ 要改注释＋给新钉子的存在理由，属裁量项；② 会把票 95 那张"故意留着宽"的决定在 POSIX 上一并改写（同 §1.3 的 models 撞钉，只是换了平台）；③ **本机验不了**：这半边的任何读数只能在一台 linux 容器里产生（§3.4），而 `cmd/wisp` 连 Linux 构建都不存在（`cmd/wisp/logsink.go:52`–`:61` 逐字：「package main has no Linux build, because `GOOS=linux go vet ./cmd/wisp/` is rc=1 on the sherpa-onnx import main.go carries … nothing in this file has ever been *compiled* for Linux」）⇒ **乙形若把钉子放在 `cmd/wisp`，POSIX 那半边今天连"编译过"都证明不了，更别谈跑** |

### 3.3 具名注释草案（若 AC#3 选"把 POSIX 侧今天不验写进判据"）

给落地腿当起点（措辞要自带"为什么不是跳过"，本仓对跳过的态度见 `internal/winsec/private_other_test.go:15`–`:16`「A skip would be the wrong shape here」）：

```
// POSIX 半边，具名不验（票 238 AC#3）。
// 本钉子断言的是"根与 logs 这两枚目录自己带一枚非继承的窄描述符"，
// 这是 Windows 的说法：icacls 的 (I) 记号在 POSIX 上没有对应物，
// chmod 之后 mode 就是那个 mode，没有"存哪儿"这一维（winsec_other.go:24-:41 只回读、不回读 flag）。
// 因此本钉子在 !windows 下【不跑】，而【不是跳过】：
//   - 它在 POSIX 上能判的那一半已经有人判了：internal/winsec/private_other_test.go:37
//     TestPOSIXPrivateDirIsReally0700（PrivateDirAll 的父层与自身都必须读回 0700）；
//   - POSIX 上真正没被管住的是每枚 wisp-*.jsonl 的 mode（logging.go:244 的 0o644），
//     那一格属票 132 AC#2 的"每枚文件首写前 SealFile"，不属本票；
//   - ⛔ 不许用"容器里跑过了"代替本行：这一半的实测读数必须由 §3.4 那一发产生并落在 docs/evidence/s1/238-*.md。
```

### 3.4 要跑哪一发、在哪跑（⛔ 本腿没跑，硬红线禁容器；下面每条都是给落地腿的派单手递）

1. **在哪**：本机 Docker，镜像 `golang:1.27`（本地已有）。⚠ 本仓已登记的挂载坑（票 132 Rules 逐字）：
   「POSIX 半边要 Docker **真跑**（Git Bash 下 `docker -v C:\…` 会静默挂空且 rc=0）」⇒ 仓库路径含空格（`projects plans`），
   挂载源必须整体加引号，且**起手要先证"容器里看得见 go.mod"**再谈读数（`ls /src/go.mod` 一发，空即仪器坏了）。
2. **跑哪一发**（本票 POSIX 半边最小集）：
   `go test -count=2 -v ./internal/winsec/ ./internal/observe/`
   ⛔ **不要加 `./cmd/wisp/`**：该包在 Linux 无构建（`logsink.go:52`–`:61` 逐字＋`internal/proc/crossvet_test.go:22`–`:26` 解释 cgo 与 `CGO_ENABLED=0` 那一段），
   加上去只会拿到一条"包加载失败"，会被误读成本票的红。
3. **正控（必带，否则这发只是重跑别人的绿）**：同一容器内先跑一次**默认 TMPDIR**，再跑一次 **`TMPDIR` 指到一枚 symlink 底下**的形式。
   第二发的预期结果在仓里有既数：`internal/winsec/winsec_other.go:97`–`:99` 逐字
   「R-113-B, measured in a Linux container: **81 failing lines** across this package and internal/config, internal/agent, internal/memory with **TMPDIR behind a symlink**」
   ⇒ **这一发是本票最要紧的 POSIX 读数**：任何把封条放进 `internal/observe/logging.go`（形②‑甲）的走法，
   都会让 `logging_test.go:243`/`:286` 两枚直接调 `InitLogWithRegistry(LogConfig{Dir: t.TempDir()…})` 的用例
   **走进 `platformVerifyPlacement` 的拒绝面**（`winsec_other.go:155`–`:163`：任一前缀是 symlink 就 `ErrUnresolvedPath`），
   而 `internal/observe` **不在 `proc.SealableRoot` 的下游**（那条 resolve 只发生在根解析点：`cmd/wisp/doctor.go:263`、`internal/proc/envfork.go:125`/`:245`）
   ⇒ **后果形状**：在 macOS（`/var` 是 symlink，`winsec_other.go:100`–`:102` 已点名）与"TMPDIR 挂在链接底下"的容器里，**日志管道会开始拒绝启动**，
   而 `wisp slo` 那枚腿的失败方向是 `return 2`（`cmd/wisp/slo_windows.go:272`–`:278` 现读：`InitLog` 一失败就打印 `wisp slo: log pipeline: %v` 并退 2），
   **不是**常驻/`run` 那两条"响亮但不拦"⇒ 同一枚改动在三条腿上给出两种失败方向，这条要进 AC#1 的裁据（§4 第 3 条）。
4. **另一枚 POSIX 专属的拒绝面（具名）**：`ResolvePath` 的内置闸门 `builtinVerifier.Resolve`（`internal/winsec/resolve.go:664`–`:675`）
   第一腿就是 `filepath.IsAbs` ⇒ 一旦 `logging.go:72` 换成 `PrivateDirAll`，**一个相对拼写的 `LogConfig.Dir` 会从今天"照建"变成"拒绝并报错"**。
   现读全仓 `LogConfig{` 只有 4 处（`logsink.go:149`、`slo_windows.go:272`、`logging_test.go:243`/`:286`），四处传的都是绝对拼写（`t.TempDir()`／`filepath.Join(DataDir,…)`）
   ⇒ **今天不会因此多红一枚用例**〔未跑包，此条来自读断言＋名册穷举〕，但它是一条**新长出来的行为边界**，判据注释要写。

### 3.5 本票 POSIX 半边的三行读法

① 现象在哪出现＝盘上那枚目录的 mode（POSIX）／描述符（Windows），不在任何对外接口；
② 有没有本机被入侵的证据＝**没有**（本腿零写入，且这台机器是 Windows，POSIX 那半边今天连现场都没有）；
③ 最坏后果形状＝POSIX 上 `logs` 的 mode 由 `logging.go:72` 的 `0o755`（去 umask）定、每枚 jsonl 由 `0o644` 定，
**同机其它账户可读**这一句在本腿这儿**只是代码读数、不是实测读数**（要实测就得 §3.4 那一发）⇒ 写进证据表时只能标〔读码推定〕，不许标〔实测〕。

---

## 4. 要编排者裁的（不许为空）

⛔ 本节每一条都是**票面没覆盖、本腿不许自己假设方向**的。每条给：争议／三支候选／本腿能量到的代价／不裁的代价。

| # | 要裁的那一刀 | 候选与代价（本腿现读撑着的） | 不裁会怎样 |
|---|---|---|---|
| **1** | **两枚"传递已有、直接未裁"的依赖边能不能开到明面上**：`internal/proc -> internal/winsec`（形①）、`internal/observe -> internal/winsec`（形②‑甲） | 仓里**没有成文禁止**：`AGENTS.md` §1.2 那串禁止形状里没有依赖方向；成文的是**反方向禁**（`resolve.go:35`–`:40` 逐字「internal/winsec cannot import internal/risk, because the graph already runs risk -> observe -> secret -> winsec … so the reverse edge closes a cycle」、`logging.go:386`–`:387`、`logsink.go:132`–`:134`）。两枚直接边都不会闭环（`winsec` 是叶子，§0.4） | 落地腿自己判＝agent 单方面改架构形状；下一位会拿"注释里那句话"当规矩用，而那句话禁的是另一枚边 |
| **2** | **封根这一手要不要带 propagation（全树走查）** | 甲＝新增"只封一层不走查"的导出名（**契约级**，照票 132 `:76` 乙支的先例：实现腿不许自己造，要回编排者落 `A##` 并摆给 owner）；乙＝走查但加子树例外表（⇒ 第二真相源，且本仓对"例外表"已有包袱：`[fs] reparse_point_exceptions` 只在判级侧、`risk/winsec_c26.go:14`–`:19` 明写"placement 不许带例外"）；丙＝只封 `<data>\logs`、根不动（⇒ **本票标题那句"根是谁造的"仍然无人答**，AC#2 判据落不了地） | 直接撞 §1.3：首启那一次会把 `<data>\models\**`（`models.go:198`）的 `BUILTIN\Users` 读权摘掉，而那正是票 95 拍板"故意留着"的决定，`no_seal_ruling_windows_test.go:21`–`:22` 逐字「decisions like that do not get made by an unrelated hardening sweep」 |
| **3** | **封条失败的方向**（拒启动／降级照跑／分腿不同处理） | 形① 在 `Boot` 里失败＝`resident_windows.go:42`–`:45` `os.Exit(1)`＝**双击图标起不来**；形②‑甲 失败＝4 条腿照跑（`resident_windows.go:58`–`:65`、`run.go:196`–`:203`、`models.go:284`–`:286`、`secret.go:226`–`:229` 现读全是"打印＋继续"）＋**`slo` 那条直接 `return 2`**（`slo_windows.go:272`–`:278`）。⚠ 票 128 的"拒绝启动"之所以没堵死联调，靠的是台账 `A111④`（`docs/reports/pending-and-issues.md:3689`–`:3696`）那两条**在选树之前**的前置路径（`portable.txt`／`WISP_ENV=test`，代码＝`doctor.go:247`–`:267`）；**"封不上"这一步没有任何前置路径**，且全仓**没有** `--data-dir`／`WISP_DATA_DIR` 这类对外口子（同一条台账 `:3695`–`:3696` 已记为缺口） | 无条件拒启动＝堵死 owner 联调（本仓定案过的教训）；无条件降级＝"根宽"只剩一条日志，且 `slo` 这条**自救用的**腿可能先起不来 |
| **4** | **AC#2 那把尺的射程到哪一层** | 甲＝只判"本进程创建的那串目录"（`R`、`R\logs`）：能判、正控成立、不撞 §2.6 的 R 类三枚；乙＝扩到每枚 `*.jsonl`：⇒ **票 132 的 AC#2（每枚文件首写前 `SealFile`）从"另一张票"变成本票的前置条件**，两票关系从"串行"升级成"互相承重"，要重排；丙＝再加"时序半"（X2 那一发，见 §2.4）⇒ 需要一枚能在中途读 ACL 的载体，`cmd/wisp` 里现成形状是 `logsink_windows_test.go:308`（日志文件里读记录序号），真机中途读数则要有个能停下来的钩子 | 射程写歪的后果是两种：写窄了＝AC#2 只证了"有封条"没证"首启那棵树"；写宽了＝本票悄悄把票 132 的格占了，132 的裁决表还没有一枚（`docs/evidence/s1/132-*`＝0 枚，票 132 `:88` 现读） |
| **5** | **新尺的保 flag 解析器放哪儿** | 仓里唯一保留整条 flag 文本的解析器是 `internal/models/acl_sid_121_test.go:50 parseACELines`（＋`aceLinesForSID :100`），它在 `internal/models`；`internal/winsec/acl_windows_test.go:86 aclSIDs` 明确丢 flag。⇒ 要么抽公共测试件（跨包，动别人地界），要么在新钉子内复制一份并注明出处（多一份 parser＝多一份漂移） | 复制＝将来两边判定不一致；抽公共＝碰 `internal/models` 那枚包的地界（本腿不裁） |
| **6** | **POSIX 半边选甲还是乙**（§3.2） | 甲＝不递归＋判据注释写明"POSIX 上封根不保护任何东西，每枚 jsonl 属票 132"；乙＝补递归＝推翻 `winsec_other.go:165`–`:169` 那段存在理由（那句现在还逐字写着"the hole being closed does not exist on this platform"），要改注释＋给新钉子的存在理由 | 不选＝落地腿会照 Windows 的直觉写一支"两边都对"的修法，而 §3.1 说这两边**不是**同一件事 |
| **7** | **本票没有"地界"这一节** | 票面只有 `:21`–`:23`「禁区」，**没有**票 132 `:6` 那种"地界"行。而本票要动的文件横跨 `internal/observe/logging.go`、`internal/proc/boot_windows.go`、`cmd/wisp/{logsink,resident_windows,run,models,secret,slo_windows,doctor}.go`、`internal/winsec`（若走甲形新导出名）——⚠ 且票面 `:23` 自己写着"与票 132／174／175 同撞 `internal/winsec`／`internal/observe` ⇒ 串行" | 落地腿改到哪算合规无法判；跨包改动在共享工作树里容易撞别人在飞的写腿（此刻 `internal/tools`／`cmd/wisp` 有写腿） |
| **8** | **AC#1 完成判据里那发"真机首启读数"由谁在哪跑** | 票面 `:17` 要求"一发起跑时序的真机读数（新目录、先跑常驻腿那形）"。本腿按红线**没跑**（零执行被测程序、零 `icacls`）；参照做法在 `docs/evidence/s1/128-ac1-consequences.md:9`–`:10`（仓外临时根＋两个 CWD）与 `:217`（对 `%USERPROFILE%` 之外的树**只读**） | 没这发读数，AC#1 只有"两形＋代价"，判据那一格仍不成立 |

## 5. 没做完的（不许为空）

1. **没跑任何测试／容器／`icacls`／被编译出来的 exe**（硬红线）。因此下面这些全是**空缺**，不许被当成已成立：
   - §2.6 名册的"今天绿"一律标〔未跑包，此条来自读断言〕；落地腿要跑 `-count=2 -v ./cmd/wisp/ ./internal/winsec/ ./internal/observe/ ./internal/models/` 对表；
   - **AC#1 那发真机首启读数**（新目录、先跑常驻腿）未取；
   - **POSIX 半边未量**（§3.4 给了在哪跑、跑哪一发、两发 TMPDIR 对照与"容器起手先证 `/src/go.mod` 在不在"的仪器自检）；
   - "宽多少、谁能读到"仍**没有读数**（票面 `:5` 那句维持原样，本腿没升级它）。
2. **没跑 `go build ./...`**——虽然它在允许清单里，本腿主动不跑：另一枚验收腿正在**独占桌面采样 CPU 与句柄**，一发全模块编译（含 cgo/sherpa）会洗掉它的读数。
   自证本腿没写坏东西的替代尺＝本腿三枚 commit 的文件清单（`git show --stat` 只看那三枚），**只含 `.scratch/wisp/probes/238/c1/census.md` 一枚文件、零 `.go`**。⚠ 此条是"我用编译换来了另一枚腿的读数干净"，属**口径偏离**，编排者若坚持要编译证据，请在验收轮补跑。
3. **没量 propagation 走查的开销形状**：形① 会让常驻腿**每次启动**重写整棵数据根的描述符（`winsec_windows.go:588`–`:620`），这与 D32 的启动预算／空闲上限有没有冲突，本腿没算（要算得有读数，而那要跑真进程）。
4. **`memory.Open` 在 `open.go:170` 用的是 `filepath.Abs`** 这一条我只登记、没有追查它与 D22 ban #2／票 102 那族"改写"的关系（不属本票射程，但会污染"哪个根被谁封"的名册）；要不要单开一枚，请裁。
5. **没查 `frontend/**`／`design/**`**（红线：零读零写零转述）；没动台账 `docs/reports/*`（按派单：我不写）；没翻票面任何 `- [ ]`，没改票面原文；只在票面 `## Progress log` 末尾**追加**了一条指针。
6. **票 132 面上那 5 枚勾框一枚未碰**（票面 `:23` 明令"不许顺手把票 132 的 AC 一起翻"）。
7. **AC#2 的"逐枚"落成本腿的一种读法**（只判目录串，见 §2.1 第 3 点）——⚠ 若编排者认为票面那句"判据必须逐枚"指的是**整棵树每枚文件**，则本腿这一版**没做到**，需按 §4 第 4 条重派；我没有自行扩到那一层，因为扩了就要动票 132 的射程。
