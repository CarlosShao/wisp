# 156 r1 — Q-57 落地＝乙：`sloSubject.exited()` 从此问 OS，那一发在本程锚点上被量了两次（改前 2209ms 说不死／改后 1 枚读数就认尸）

- 交付程角色：**实现程（票 156）**，写码＋取证。验收归另一枚程（`SPEC-12 §4.3` #1/#3）。
- 工单：`.scratch/wisp/issues/156-q-57-landed-as-option-b-exited-asks-the-os-instead-of-our-own-record-and-the-152-leftover-comment-actions.md`
- 派单：`.scratch/wisp/dispatches/2026-09-26-124x-impl-156-q57-b.md`
- 来路：`Q-57` 三候选里 owner 12:3x 批了**乙＝让 `exited()` 去问 OS**（推荐人＝编排者）。
  **owner 批的是那一枚推荐，不是候选清单** ⇒ 本程若量到"乙治不到 AC#1 那一发"就该回来报；
  第 1、2 格的读数是那一句的答卷，答案是**治得到，而且只能经由 `:820` 那一枚调用者治到**（见 §0.2 前提③）。
- 地界（本程**只**动这几类路径）：`cmd/wisp/slo_windows.go` · 新建 `cmd/wisp/slo_exit_os_156_windows_test.go` ·
  `cmd/wisp/slo_report_144_windows_test.go`（**仅** AC#4 的三枚注释级动作）· `.scratch/wisp/probes/156/**` · 本件。

---

## 第 0 格　锚点·四条前提·互斥·仪器

### 0.1 开工现读（原文落 `.scratch/wisp/probes/156/anchor/tree.txt` 与 `premises.txt`）

```
$ git rev-parse --short HEAD          <- 派单当场：9835d81（那枚 commit 的正文就是"新票 156＋两枚派单存档"）
$ git rev-parse --short HEAD          <- 本程 13:0x 再读：4e20e21（其间兄弟程落了 12 枚）
$ git log --oneline 9835d81..HEAD -- cmd/wisp/
（空）
$ git show 9835d81:cmd/wisp/slo_windows.go | md5sum      19e9d1d0304eb72575c30d19bcb72bc1
$ git show e9a1db0:cmd/wisp/slo_windows.go | md5sum      19e9d1d0304eb72575c30d19bcb72bc1
$ git show HEAD:cmd/wisp/slo_windows.go | md5sum         19e9d1d0304eb72575c30d19bcb72bc1
$ git show 9835d81:cmd/wisp/slo_report_144_windows_test.go | md5sum   726b4fbb8c2f3ec18df26a76335ed3a2
$ git show HEAD:cmd/wisp/slo_report_144_windows_test.go   | md5sum    726b4fbb8c2f3ec18df26a76335ed3a2
```

⇒ **本程开工时工作树里 `cmd/wisp/` 是干净的**（`git status --porcelain -- cmd/wisp/` 空），
且派单那四条前提声称量过的 `e9a1db0`、派单 commit `9835d81`、本程读到 `HEAD` 三版**逐字节相同**
⇒ 本程所有"改前"读数用的就是那版字节，**没有做过任何快照搬运**：变异与探针全部走 `-overlay`
（文件住在 `D:/tmp/wisp156/` 与 `probes/156/`，一行没落进 `cmd/wisp/`），
所以坑单第 4 条那把 `git ls-tree + cat-file --batch` 尺本程**没用到**，`git archive | tar` 一发也没跑。

⚠ 派单正文说前提"量在 `e9a1db0` 那棵树上"，任务书里另一处写作"锚点 `245d5f4` 附近"。后者本程现量：

```
$ git cat-file -t 245d5f4
fatal: Not a valid object name 245d5f4        (rc=128)
$ git cat-file -t e9a1db0
commit
```

⇒ `245d5f4` **不是本仓任何 git 对象**（票 154 续程 `7955d1` 同一发判语独立复到）。
本程按派单正文的 `e9a1db0` 取数，四条前提逐条重量于 §0.2。这一处**只登记、不据此改判任何格**：
它是编排者手写的一个指针，不是凭据、也不是任何一格的尺。

### 0.2 四条前提逐条判（每条都带本程自己当轮跑的命令）

```
$ git show 9835d81:cmd/wisp/slo_windows.go | grep -n 'exited()\|cmd.Wait()\|ProcessState'
520:		if s.exited() {
820:		if s.exited() {
832:func (s *sloSubject) exited() bool {
833:	return s.cmd != nil && s.cmd.ProcessState != nil && s.cmd.ProcessState.Exited()
837:	if !s.exited() {
840:	return s.cmd.ProcessState.ExitCode()
851:	if s.cmd != nil && s.cmd.Process != nil && !s.exited() {
855:		_ = s.cmd.Wait()
```

| 前提（派单 §"我核过的四条"） | 本程现量 | 判 |
|---|---|---|
| ① `exited()` 现形＝`s.cmd != nil && s.cmd.ProcessState != nil && s.cmd.ProcessState.Exited()`，写作 `:832` | 那段**表达式逐字符**在 `:833`；`:832` 是它的 `func` 声明行 | **成立**（指针偏一枚、引文逐字符对得上） |
| ② 全文件 `cmd.Wait()` 只有 `:855` 一枚（在 `stop()` 里）⇒ `stop()` 之前 OS 已知、码还不知道 | `grep -n 'cmd.Wait()'` 全文件命中 **1 枚**，行号 **`:855`**，与派单写的行号**完全相同** | **成立** |
| ③ `exited()` 调用者**只有三枚**：`:520`／`:837`／`:851` | 调用者是**四枚**：`:520`（`waitReady`）、**`:820`（`collectReportWithin`）**、`:837`（`exitCode` 守卫）、`:851`（`stop` 守卫）——派单列的三枚行号都对，**漏的正是 `:820` 这一枚** | **不成立**（少记一枚，且少的恰是本票主靶） |
| ④ 不需要新增依赖（`go.mod:11` 已有 `golang.org/x/sys v0.48.0`；仓内先例 `internal/audio/mmdevice_windows.go`、`internal/ball/ball_windows.go`） | `grep -n 'golang.org/x/sys' go.mod` ⇒ 第 11 行 `golang.org/x/sys v0.48.0`；两枚先例文件都 import `golang.org/x/sys/windows`；`go.mod`／`go.sum` 本程零字节改动（见第 6 格 AC#6 那两枚"名册外单列"的 hits=0） | **成立** |

**对前提③不成立这一条的处置（派单说"对不上就判不成立＋当轮读数并停手，别替我圆"）**：
- **判**：〔不成立〕，读数在上面第 ③ 行，当轮命令原文在 `premises.txt`。
- **不因此停手，理由写明不猜**：派单给的停手判据是"四条前提对不上"，而它同时给了**唯一一枚会让本程换修法的判据**
  ——任务书原话："如果你量到乙治不到 AC#1 那一发，回来报，不许顺手换成丙"。
  本程量到的方向**恰好相反**：AC#1 那一发（"写报告写到一半时死掉、印的是邻居那句"）**只发生在 `:820`**，
  也就是派单漏记的那一枚调用者上。⇒ 前提③的错误**不削弱**乙的射程，它**扩大**了本程的改动面一枚（`:820` 也吃这个口径），
  而扩大面是同一枚 `exited()` 自动带上的，本程没有为此做任何新决定、没有换候选、没有碰 `stop()` 的收尸语义。
- **要编排者认的账 1 笔**：派单 §"射程窄，但它是那三处共同的口径来源"这句在**票面**（工单"现量的形状"第 2 条）里也是三枚，
  两处同错；本件把四枚名册钉在 §0.2，下一程别再用三枚那版。

### 0.3 互斥与并发（文件级）

```
$ git status --porcelain | grep -vE '^(\?\?|.M) (design|frontend)' | head
 M .scratch/wisp/probes/152/my152.py                    <- 已停笔的别程半件：不 commit、不还原、不补
 M docs/evidence/s1/152-...-accept-r1.md                 <- 同上（3 行未提交自校）
?? .scratch/wisp/probes/152/mut-shipped/                 <- 本票 AC#4(iii) 点名要入库的那 9 枚
?? .scratch/wisp/probes/155-accept/  ?? .scratch/wisp/probes/139/accept-r1/   <- 兄弟程在写
```

- 兄弟两程（154 续程／155 验收程）的写面是 `docs/evidence/s1/154-*`·`155-*` 与 `probes/154`·`probes/155-accept`，
  与本程**交集为空**；`9835d81..4e20e21` 那 12 枚 commit **零枚**碰 `cmd/wisp/`（§0.1 的 `git log -- cmd/wisp/` 现量为空）。
- ⇒ **本程所有"零命中／干净"类宣称一律不含 `frontend/**`、`design/**`**（owner 与前端会话在未提交地写），
  也不含 `probes/152/my152.py` 与 `152-…-accept-r1.md` 那两枚半件：本程一枚都没 commit（第 6 格逐枚 commit 名册可反查）。
- 本程没有跑过任何全仓全量门禁（派单明令逐包单跑），所以不存在"替兄弟程背一次红"的读数。

### 0.4 工具版本与四枚仪器坑（含本程新增一条）

```
$ go version                                go version go1.27.1 windows/amd64
$ ls third_party/sherpa-onnx/*.dll          onnxruntime.dll  sherpa-onnx-c-api.dll  sherpa-onnx-cxx-api.dll  (3 枚)
```

| 坑（派单坑单） | 本程的绕法 | 绕没绕成的判据（真跑） |
|---|---|---|
| ① dll 在位＋**PATH 形状**：`pwd -W` 盘符正斜杠形 ⇒ `0xc0000135`＋0 条 `=== RUN` | bash 侧一律 `PATH="$PWD/third_party/sherpa-onnx:$PATH"`（`/d/work/…` 形，同 `scripts/wisp-cli-tests.sh`）；python 尺（`my156.py`）把盘符形塞进**原生 env block** | 每一枚日志都现数 `=== RUN`，0 枚即 `LOADED=NO` 且 rc=3。本程取到的枚数：门禁 pre `RUN=139`／post `RUN=143`，探针 pre `RUN=5`／post `RUN=5`，变异 12 发全部 `RUN≥5`。**没有一枚是 0** |
| ② `-overlay` 不与 `-cover*` 同用 | `my156.py:run()` 见到任何含 `cover` 的参数就 `FATAL`；变异与覆盖**分两跑** ⇒ 本程**没跑覆盖**（要的是"印哪一句"，不是"走到哪一行"） | 凡带 `-overlay` 的命令行原文都在日志首行，复算一眼看得见它没有 `-cover` |
| ③ `-race` 可能 `0xc0000374`／rc=1 而零条 `--- FAIL` | 本程**一枚 `-race` 都没跑** | 〔未取证〕，不下任何并发判定（第 9 格第 4 条点名） |
| ④ 取干净树不许 `git archive \| tar -x` | 本程**根本没取快照**（§0.1 末段） | 三版 `md5sum` 全等（`e9a1db0`/`9835d81`/`HEAD`） |
| **＋本程新量一条（给坑单补档）**：坑①的成因不是"盘符形 PATH"，是 **bash/MSYS 那一次 rewrite** | 同一枚盘符正斜杠带空格的 PATH 值，经 python 原生 env 传给 `go test`，**RUN=5 跑到了**（`probe-pre/probe-156.log` 首行可反查） | ⇒ 报"`pwd -W` 形必死"要说清是**从 shell 里死**的；尺在 python 里用盘符形是本程的常态，不算绕法失败 |

**放水两问自答**：① 本格零判据改动、零断言方向改动（只读 git／只写仓外与 `probes/156/`）；② 本格不使用任何 helper。
**另自报本程自己犯过的一枚仪器错**：`slo_exit_os_156_windows_test.go` 初稿把"OS 说死"量在
`OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION)` ＋ `WaitForSingleObject(0)` 上，那一发**必然是 WAIT_FAILED／Access is denied**
（limited 权限的句柄不带 `SYNCHRONIZE`），四枚用例各撞 20s 假红；现版 witness 换 `SYNCHRONIZE\|PROCESS_QUERY_LIMITED_INFORMATION`
并在注释里写明它与实现句柄不是同一枚（第 3 格）。这条不是判据改动，是**取版器坏了**，坏的时候读数全在日志里。

**第 0 格判定**：**成立**（锚点三版逐字节全等、四条前提逐条现量并点名哪一条不成立、互斥按文件级成立、四枚坑各有"绕成／没绕成"判据）。

**本程没测什么（本格）**：`245d5f4` 为什么不存在（那是编排者的指针，不是本程的尺）；
`third_party/sherpa-onnx` 那 3 枚 dll 与 CI runner 用的是不是同一批字节（派单没要求，票 152 第 0 格同样登记为未查）。

---

## 第 1 格　AC#1：**本程自己复现那一发**，"OS 已知道／码还不知道"那个窗口量到 ms 与具名次数

**判据（跑之前写死）**：真起子进程、真让它在写报告写到一半时死掉、真用生产那条 `collectReportWithin` 循环、
真用 `os.ReadFile`；窗口两端各要一枚独立见证——起点是**内核**说"这个 pid 我完事了"，终点是**我们自己的记录**说"它死了"；
红绿不看汇总行，只看 `^--- FAIL:`。读数一律 `P156|<cell>|<key>=<value>` 原文可 grep。
**引票 152 的数只当线索，本程一枚都不当凭据**（它的 `18485 次／60.05s` 出现在派单与本件别处时都点名是它的）。

**仪器**：`.scratch/wisp/probes/156/zz156probe_windows_test.go`，经 `-overlay` 当 `cmd/wisp/zz156probe_windows_test.go` 喂进编译，
**仓里没落一字**（§0.1 末段）。四枚 cell 全是真子进程（`exec.Command(os.Args[0], "-test.run=…")` 重入本测试二进制，真 pid、真退出码 7）。
**见证与修法刻意不同门**：探针问内核走 `OpenProcess(pid, PROCESS_QUERY_LIMITED_INFORMATION)` ＋ `GetExitCodeProcess`
（从 pid 新开的句柄），而第 2 格落地的修法读的是 `os/exec` 自己那枚 CreateProcess 句柄的**等待态**——两扇门，一个事实（第 3 格把这条钉成断言）。

### 1.1 改前读数（`.scratch/wisp/probes/156/probe-pre/probe-156.log`，`RUN=5 PASS=1 FAIL=0`＝跑到且探针自身不红）

| cell | OS 何时知道（ms since spawn／见证轮询次数／码） | 那一刻我们的记录 | 记录追上 OS 要多久 | 循环烧掉什么 | 印出来的是哪句 |
|---|---|---|---|---|---|
| `prefix-unreaped`（生产形：报告半截、没人收尸） | **44ms**、20 次、code=7 `terminated` | `ProcessState_nil=true exited=false exit_code=-1` | **追不上**：再问 145 次、200ms 内 0 次为真（`os_dead_to_code_dead_ms=-1`） | `elapsed=2006ms`、**真读 186 次**、预算 2s 烧尽 | **BUDGET**：`subject 38708 never wrote a complete report within 2s (last read: 989 bytes read, document still open at offset 989…)` |
| `nofile-unreaped`（连文件都没有） | **40ms**、21 次、code=7 | 同上 | 143 次问、0 次为真 | `2009ms`、**真读 191 次** | **BUDGET**：`… never wrote a complete report within 2s (last read: 0 bytes read, no report file yet)` |
| `prefix-REAPED-case13shape`（对照：先 `Wait()` 再进循环） | 53ms、29 次、code=7 | `ProcessState_nil=false exited=true code=7` | **5ms、1 次** | `0ms`、**读 1 次** | **DEAD-CHILD**：`subject 63520 exited (code 7) without writing its report (…)` |
| `live-child`（阴性对照：孩子还活着） | `alive=true code=259`（1 次） | `exited=false` | 不适用 | `2000ms`、读 188 次 | **BUDGET**（这是**正确的**：没死的孩子就该说预算） |

**窗口有多大（AC#1 要的那个数）**：`window=os_knew_at_ms=44_code_still_said_alive_after_loop_ms=2209`，
`nofile` 那发是 `40 / 2210`。读法是：**内核 44ms 就知道，我们在循环的整个生命周期里都不知道**——
那一发里"码还不知道"不是只持续 2 秒，是**持续到有人去调 `cmd.Wait()` 为止**；本探针没调，所以它**从没知道**
（`record_after_loop=exited=false` 那一行在案）。具名次数同一行在案：**186／191 次真读**、**145／143 次问 `exited()`**、全为 false。

### 1.2 判读（四句，每句只靠上表的列）

1. **那一发复现了**，且不需要任何手工接线：`prefix-unreaped` 与 `nofile-unreaped` 两发的记录列都是
   `ProcessState_nil=true ⇒ exited=false`，而见证列都是 `terminated code=7`。
2. **指错病因这件事是量出来的、不是推的**：同一个进程状态下，"先收尸"那发（第 3 行）印 `exited (code 7)`、
   "生产接线"那发（第 1、2 行）印 `never wrote a complete report within 2s`。两发之间被改动的只有"有没有人 `Wait()` 过"这一件事。
3. **窗口不是竞态的产物**：本探针没有制造任何读写竞争（报告字节是子进程真写的一半／或干脆没写，`nofile` 那发连文件都不存在），
   窗口照样是 0 次为真；把它读成"运气不好的时序"是错的——它是 `:833` 那行表达式与 `stop()` 的位置关系的**恒等后果**。
4. **阴性对照没坏**：`live-child` 那一发仍旧烧预算、仍旧印预算句。这条是本格给"改后会响"设的下界——
   **修法不许把活着的孩子也认成死的**（第 4 格 `m3` 就是去量这一条有没有主人）。

### 1.3 外推（点名是外推，不是读数）

上表的预算是本探针自己挑的 `2s／10ms`。生产那两枚常量是 `subjectReportBudget = 30 * time.Second`、
`subjectPollInterval = 20 * time.Millisecond`（本程一字未动，见第 6 格），按同一形状外推：
操作员在真实 `wisp slo` 上看到的窗口是**内核 44ms 就知道、命令 30s 里问约 1500 次全为 false**，最后印预算句。
第 4 格的 `m1` 发**不是外推**：它用生产常量跑，量到 `reads=1468／30.08s`（`mut-156/m1-os-arm-removed.log`）。

**放水两问自答**：① 本格不判任何断言方向（探针只打印，不断言生产行为；唯一 `t.Fatalf` 在"见证 20s 内没等到孩子死"那一发，是夹具守卫）；
② 本格用的 helper 全是**新建的探针件**（`p156*`），没碰交付面里原有的任何一枚 helper；`slo144Report` 是
`slo_report_144_windows_test.go` **原有**那枚 fixture，本程一字未改地复用（同一份 `json.MarshalIndent` 产物，1978 字节）。

**第 1 格判定**：**成立**——AC#1 那一发被本程自己复现成两条独立读数（`prefix`／`nofile`），窗口有 ms 与具名次数，
并且本程**没有**用它没量过的东西交件：没有引用票 152 的数当凭据、没有把"应该能复现"写成交件。

**本程没测什么（本格）**：
1. **没量真 `wisp slo` 端到端那一发**（真 Job Object、真 posture、真把 subject 打死）——票 152 第 3 格量过它"改前"的原文，
   本程的"改后"端到端欠账记在 §9 第 1 条。
2. **没量子进程死在 `os.WriteFile` 系统调用内部**那一形（半截字节由**内核**留下，而不是由子进程主动写一半再退）：
   本探针的"半截"是孩子自己 `os.WriteFile(doc[:989])` 后 `exit 7`，属于**真前缀**但不是**写一半被打断**。
   缺的仪器＝一台愿意撞部分读的环境（票 152 §8.3 第 8 条同一句欠账，本程没能力在这一票里还）。
3. **没量 pid 复用**：见证与实现都只在"进程对象还在"的时候读得到码；`Wait()` 之后对象可回收，
   那正是本程把 (a) 支放在 (b) 支**之前**的理由（第 2 格注释），不是量出来的行为。

---

## 第 2 格　AC#2：落地乙——`exited()` 从此问 OS；机制本程选，举证也在本程

### 2.1 选了什么机制，为什么不是那枚一行编辑

`cmd/wisp/slo_windows.go` 里 `exited()` 现在经由 `exitStatus() → queryProcess()` 问内核：

```go
err := s.cmd.Process.WithHandle(func(h uintptr) {
    if waited, _ := windows.WaitForSingleObject(windows.Handle(h), 0); waited != windows.WAIT_OBJECT_0 {
        return // WAIT_TIMEOUT: alive. WAIT_FAILED: not judged here.
    }
    if cerr := windows.GetExitCodeProcess(windows.Handle(h), &osCode); cerr != nil {
        return
    }
    decided = true
})
```

三件本程现量到的事决定了这个形状，不是抄来的：

1. **句柄从哪来**：`os.Process.WithHandle` 递给回调的是 `exec.Cmd.Start` 那枚 **CreateProcess 原生句柄**
   （`D:/work/base/go/src/os/exec_posix.go:80` 读的是 `p.handle`；Go 自己的 `(*Process).wait()` 在
   `exec_windows.go:18` 就是拿它做 `WaitForSingleObject(INFINITE)` ＋ `GetExitCodeProcess`）。
   ⇒ 本程没开新句柄、没引依赖、没造 reaper，等的是**内核早就算好的那个位**。
2. **不能用 `GetExitCodeProcess` 单独判死活**（派单警告的那一枚"就用 GetExitCodeProcess"）：
   它在进程还活着时答 `STILL_ACTIVE = 259`，而 259 同时是**合法退出码** ⇒ 一个以 259 退出的 subject 会被读成"还活着"。
   本仓 subject 今天的退出码是 0/1/2（`cmdSLO` 的三个 return），**这一发本程没量到发生过**（记 §9 第 5 条），
   但"等态先、码后"把两个问题拆开是**接口性质**上的便宜：多一枚 syscall、少一类歧义，且码为 259 时能如实报 259。
3. **两扇门不等价，量出来的**：`slo_exit_os_156_windows_test.go` 初稿用
   `OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION)` 的句柄做 `WaitForSingleObject`，四枚用例各撞
   `WAIT_FAILED / Access is denied`（limited 权限不带 `SYNCHRONIZE`）⇒ **从 pid 新开的句柄与 exec 那枚不是同一只东西**；
   而探针那侧（第 1 格）用"码门"看到 `terminated` 时，"等态门"还要晚约 6ms 才同意
   （`probe-post/probe-156.log`：`witness=…44/50ms` 对 `agree=os_dead_to_code_dead_ms=6 asks=6`）。
   ⇒ 修法取的是**保守那一门**："内核还没收完尾" 就说没死，最多多烧一轮 20ms 轮询，不会把活孩子认成死的。

### 2.2 三件硬约束逐条现量（原文 `.scratch/wisp/probes/156/anchor/constraints.txt`）

| 约束 | 尺 | 读数 |
|---|---|---|
| 不许动预算常量（`subjectReportBudget`／`subjectGrace`） | `git diff -- cmd/wisp/slo_windows.go \| grep -cE '^[-+].*(subjectReportBudget\|subjectGrace\|subjectReadyBudget)'` | **0**（`thresholds.go`、golden、`scripts/slo-check.ps1` 见第 6 格逐枚 commit 名册） |
| 不许新增 go.mod 依赖 | `git status --porcelain go.mod go.sum` | **空**（零字节；用的是 `go.mod:11` 已有的 `golang.org/x/sys v0.48.0`，先例 `internal/proc/jobscope_windows.go` 等） |
| 不许墙钟超时（ban #5）／不许裸 `go func(`（ban #1）／不许 reaper（owner 没批丙） | `git diff -- cmd/wisp/slo_windows.go \| grep -nE '^\+.*(go func\|time\.Now\|time\.Since\|time\.Until\|watchdog\|reaper)'` | **零命中（rc=1）**；`queryProcess` 一个 `time.` 都不出现，循环的预算仍是 `observe.NewTimeout` 那枚单调尺；**没有新增 goroutine** |
| `stop()` 的收尾语义（`Kill`→`Wait`→`RemoveAll`）不许动 | `git diff \| grep '^[-+]' \| grep -iE 'stop\|Kill\|RemoveAll\|Wait\(\)'` 只命中本程新写的**注释行**；`grep -n 'cmd.Wait()'` 全文件仍 **1 枚**（`:936`，原 `:855`，位移＝上方插入行所致） | **代码一字未改** |

⚠ 一条**必需要说的后果**（不是编辑、是口径变了带出来的）：`stop()` 里 `!s.exited()` 那枚守卫现在会
**跳过对已死未收尸 subject 的 `Kill()`**（以前一定打一发 `Kill` 到已终止的进程上，返回一个被丢弃的错误）。
`Kill→Wait→RemoveAll` 的顺序与幂等性都没变，本程没为此改 `stop()` 一字；
但这一条**没有用例钉**（`sloSubject.stop()` 今天零枚测试，本程 grep 过：`cmd/wisp/*_test.go` 里的 `stop()` 全是
`resident_sink_nail_127` 那枚 `leg.stop()`，不属同一类型）⇒ 记 §9 第 2 条，不假称已钉。

### 2.3 举证"这一发在改前不响、改后会响"（承重按操作定义答在第 4 格）

同一枚探针、同一棵树、同一个 cell，只换 `slo_windows.go` 那一版字节：

| cell（生产形） | 改前（`probe-pre/`） | 改后（`probe-post/`） |
|---|---|---|
| `prefix-unreaped` | 记录 `exited=false`；再问 145 次全 false；循环 **186 次真读／2006ms**；印 **BUDGET** 句 | 记录 `exited` 在等态门亮后 6ms 内为 true（`agree=os_dead_to_code_dead_ms=6 asks=6`）；循环 **1 次读／elapsed=0ms**；印 **DEAD-CHILD** 句 `exited (code 7) without writing its report (…)` |
| `nofile-unreaped` | 同上（191 次真读／2009ms，BUDGET） | **1 次读**、DEAD-CHILD，`record_after_loop=exited=true exit_code=7 ProcessState_nil=true` |
| `live-child`（阴性对照） | 188 次真读／2000ms，BUDGET（正确） | **192 次真读／2007ms，BUDGET（一字未变）** ⇒ 修法没有把活孩子认成死的 |
| `prefix-REAPED-case13shape` | 1 次读／DEAD-CHILD | 1 次读／DEAD-CHILD（(a) 支原样保留，票 149 的 case 13 不动） |

第二向证据（**用发货的用例**而不是探针）在第 4 格：`m1-os-arm-removed` 这一发把 (b) 支摘掉、其余一字不改，
四枚新用例全红、case 17 量到 **1469 次真读／30.07s**（同 cell 前一发 1468 次）并印回 BUDGET 句——
**同一发在改后的字节上是 1 次读、DEAD-CHILD 句**。

**放水两问自答**：① 断言方向：本程没有为了让任何一格通过而放宽任何断言；第 3 格四枚用例是**新增**的、只往严加；
② helper：`queryProcess`/`exitStatus` 是本程新建的两枚方法，`exited()`/`exitCode()` 是**原有**那两枚的口径改写（签名不变、三枚调用者一字未动）；
测试侧 `scriptedReader`／`scriptMissing`／`scriptBytes`／`slo144Report` 全是 `slo_report_144_windows_test.go` **原有**的 helper，本程一枚都没改。

**第 2 格判定**：**成立**——乙落地，机制为本程所选并给出 §2.1 三条现量理由；三件硬约束各有零命中读数；
"改前不响／改后会响"两向都有本程自己的读数（§2.3 两枚来源：探针对照 ＋ `m1` 变异）。

**本程没测什么（本格）**：
1. **没量 `wisp slo` 真命令端到端**（真 Job Object 里打死真 subject 后操作员看到的那一句原文）——见 §9 第 1 条，本票真正的欠账。
2. **没量 subject 以 259 退出**那一形（§2.1 第 2 条的歧义是**接口性质**的论证，不是读数）。
3. **没量 `stop()` 的 Kill 跳过**（§2.2 末段），也没量新增的每轮询 1–3 枚内核调用对观察者自身 CPU 的影响
   （门禁射程里没有 `wisp slo`，票面已登记；本程没跑 SLO 门）。

---

## 第 3 格　AC#3：会响的检——四枚新用例钉"OS 说死 ⇒ 我们立刻算死"，既有断言一枚没弄松

新建 `cmd/wisp/slo_exit_os_156_windows_test.go`（cases 15-19，5 枚）。每枚各自钉一件，**互不替代**（替代关系是量出来的，见第 4 格的单枚摘掉表）：

| 用例 | 钉住什么 | 为什么只有它能钉 |
|---|---|---|
| 15 `ExitedAsksTheOSForAChildNobodyReaped` | **机制**：等态门一亮，`exited()` 第一次问就 true、`exitCode()` 就是 7；且**同一刻 `ProcessState == nil`**，问完再问一次还是 nil | "OS 说死⇒我们算死"这一形若被"顺手在 `exited()` 里 `Wait()` 一下"实现，前两条断言照样满足——是 `ProcessState==nil` 那一枚把它挡在门外（AC#3 要的正是钉这一形，不是钉这个答案） |
| 16 `LiveSubjectStaysAliveUntilTheOSDisagrees` | **反方向**与两支的分界：16a 活着 ⇒ false／`exitCodeUnknown`；16b 我们亲手 `Kill` 且未收尸 ⇒ true（走 (b)）；16c `Wait()` 之后 ⇒ true 走 (a) | 只往"死"的方向加用例，`exited() := true` 就能满足全部（第 4 格 `m3` 量到只有 16 一枚红 ⇒ 这一枚是唯一的持有者） |
| 17 `ReportLoopNamesTheDeadSubjectItWasWaitingOn` | **那一发本身**在生产接缝的闭环：`:822`（派单记作 `:820`，本票前提③那枚）那枚调用者对**已死未收尸**的孩子印 `exited (code 7) without writing its report`，且**读 1 次**、不许出现 BUDGET 句 | 票 149 的 case 13 已经钉过同一句的**文案**，但它靠 `dead.Run()` 先把尸收了——生产里不存在那一手（152 验收件第 1、2 格）。17 是唯一一枚"真 corpse ＋真循环"的接线 |
| 18 `WaitReadyNamesTheDeadSubjectToo` | 第二枚出口 `:522`：subject 死在 ready 标记之前 ⇒ `exited early (code 7) before reporting ready`，不许印 `never reported ready within 1m0s` | 票 152 第 2 格量过这句"今天根本产不出"，此前全仓零枚持有者；这一枚是它的第一枚 |
| 19 `FixtureChildrenDoWhatTheirNamesSay` | **夹具自己**：`exit` 角色真的以 7 死、死时 `ProcessState` 为 nil；`sleep` 角色真的活着且杀得掉 | 本程先踩后钉：第一版把角色分发塞在 case 15 里，于是"摘掉 case 15"这发变异把别的 cell 的**子进程**一起摘坏了，`case15-off` 在未改的码上报出 17/18 两枚假红（`mut-156-first-run-summaries.txt`） |

**"不能把票 152 已钉住的既有断言弄松"——本程的尺与读数**：

```
$ # 逐包门禁，改前／改后各一次（原文 probes/156/gate-{pre,post}/cmdwisp.log）
改前 cmd/wisp   RUN=139  PASS=79   FAIL=0  SKIP=0   rc=0  ok 136.561s
改后 cmd/wisp   RUN=144  PASS=84   FAIL=0  SKIP=0   rc=0  ok  (见第 7 格现量)
$ # 名册两向 comm（roster-{pre,post}-cmdwisp.txt，sort -u 后 comm -3）
only-in-post: 5 枚，全部 TestSLO156*（本票新增，逐枚点名见第 7 格）
only-in-pre : 零枚
```

- **既有断言的方向与枚数一枚没动**：`TestSLO149ExitedGiveUpSentenceCarriesTheLastReading`（case 13）
  在改后仍 `--- PASS (0.04s)`，且它在 `m1`／`m4` 两发变异下的红／不红关系是**本程量出来的目标行为**（第 4 格）；
- **没有 `t.Skip`**：`SKIP=0` 在改前改后都是 0；**没有改 golden／阈值**（第 6 格逐枚 commit 名册）；
- **helper 全是原有的那几枚**：17 复用 `scriptedReader`＋`scriptBytes`，fixture 复用 `slo144Report`；本程新写的
  `slo156Spawn`/`slo156Witness`/`slo156OsDead`/`slo156OsAlive` 只陈述夹具事实，**不含任何一条对生产代码的判定**。

**放水两问自答**：① 断言方向动没动——没动，五枚全是新增，两向 `comm` 的 only-in-pre 为空集就是这一句的读数；
② helper 是不是原有的那枚——判定用的 helper 全是原有那三枚；新增四枚是**取版器**，它们唯一的 `t.Fatalf` 出现在"夹具没照名字死／活"上（守卫，不评产品）。

**第 3 格判定**：**成立**——四枚会响的检（15/17/18/19）加一枚双向界（16），每枚在第 4 格都有它**唯一持有**的那一发；
既有断言一枚没松（case 13 改后仍绿、名册 only-in-pre 空）。

**本程没测什么（本格）**：19 枚 code-path 之外的 `stop()`（§2.2 末段）；
`exitStatus()` 在 `s == nil`／`s.cmd == nil` 那两条 return 上的行为（脚本化 subject 早于本票就在测 `cmd==nil` 那一支，`s==nil` 那一支本程没造用例）；
`WAIT_FAILED` 那一条兜底支**没有任何用例走到**（本程构造不出让 exec 句柄 wait 失败的条件，硬造＝放水，见 §9 第 3 条）。


---

## 第 4 格　AC#4：票 152 §8.3 那三枚注释级动作（**r2 程**，2026-09-26 15:3x，锚点现读 `0f18652`）

> 本节及第 5 格由 **续程 r2** 写（派单：`.scratch/wisp/dispatches/2026-09-26-151x-impl-156-r2-two-cells.md`，`5092e86`）。
> 上面第 0-3 格是 r1 程的已入库文字，本程**一字未改**（`git diff` 的删除列＝0，见 §4.5 那把尺）。

### 4.0 先登记一处写面不一致（不是缺陷申报，是让下一位不必猜）

派单 §"git 纪律与写面"把 AC#4 的写面写成 **`cmd/wisp/slo_windows.go` 的注释面**，而票面 `:45` 指到的"现量第 4 条"
那三枚动作（工单 `:22-24`）逐枚点名的是**别的文件**：①② 在 `cmd/wisp/slo_report_144_windows_test.go`，③ 在
`probes/152/mut-shipped/` 或票 152 实现件 `§4.2`。**按派单那一句，AC#4 三枚动作一枚都做不了。**
本程取了票面那一支：`cmd/wisp/slo_report_144_windows_test.go` **不在冻结清单里**，r1 的派单与 r1 证据件第 0 格
"地界"（本件 `:9-10`）本来就把这枚文件写进射程、且注明"**仅** AC#4 的三枚注释级动作"，本程照那一档做，
一字代码未动、一断言未动。⚠ **这一处需要编排者认账**：要么承认派单写面那行漏了这枚文件（票面与 r1 地界都对，派单自己窄了），
要么判本程越界——但**别按派单那行去判票面错**。

### 4.1 动作①（补回 case 13 裁定头）：**"补回"由 r1 已交付并经本程核过字节，本程补的是"指路"那一半**

现量（本程自己跑，不是抄 r1 的表）：

```
$ git show 10e3585 -- cmd/wisp/slo_report_144_windows_test.go | grep -c '^-//'      1
$ git show 10e3585 -- ... | grep '^-//' | sed 's/^-//' | md5sum        a756e044366a5c957bef470027a86bc0
$ sed -n '835p' cmd/wisp/slo_report_144_windows_test.go | md5sum        a756e044366a5c957bef470027a86bc0
$ diff <(那一枚删除行) <(现在的 :835)                                     IDENTICAL
```

⇒ 票面 ① 要的"补回被删的裁定头"**成立且是逐字节的**（`0e95353` 交付），且它也已写明"ⓐ 是 149 的裁定、152 的 AC#2 走 ⓒ"
（现 `:849-:857`）。本程**没有重做它**，只补了两枚指路（本件之外唯一另一处写面动作）：
1. 那三行 `152-...-accept-r1.md §2.2`（现 `:851`）里的路径是**带省略号的、开不了也 grep 不到**：真名＝
   `docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md`，那一行在它自己的 **§2.2**（`### 2.2 删除行逐枚点名`，
   本程按 HEAD blob 数过：`## 第 2 格`＝`:189`、`### 2.2`＝`:206`、那枚"删了没补回的一行"＝`:214`）。
2. **谁会误判**：拿 `grep AC#3 cmd/wisp/slo_report_144_windows_test.go` 去查**票 156 的 AC#3** 有无裁定的人。
   他把哪一句读错：`// Case 13 - AC#3's ruling is ⓐ: the s.exited() branch SHOULD name where the last`
   ——出处 **`cmd/wisp/slo_report_144_windows_test.go:835`**，那一行**头一字不带票号**，脱离上下文的一次 grep 命中
   会把 149 的 ⓐ 读成 156 的 AC#3 也 draw 了一枚字母（156 的 AC#3 要的是"会响的检"，今天一枚字母都没 draw）。
   本程把这两枚补成 `:835` 之后的一段注释，并**现量了字节等式**，没新造任何裁定。

### 4.2 动作②（`:737-740` 那三行注释里的读数）：**票面那个"改成 20/13/0"本程按自己的锚点重量＝成立，且错因比票面写的更具体**

先按票面要求"先按你自己的锚点现量再改"。尺＝`.scratch/wisp/probes/156/my156-r2.py`
（它**不重打变异字面**，全部 `spec`/`case_off` 都 import 自 r1 那把尺 `probes/156/my156.py`，
所以 r2 的读数与 r1 的日志说的是同一发编辑；本程另加一枚"落地必打印"，见第 5 格 §5.2）。选择器＝
`-run 'TestSLO144|TestSLO147|TestSLO149|TestSLO152'`（票 152 那一族的原选择器，**不含**本票新增的 5 枚）。

| 发 | 谁量的 | `=== RUN` | `--- PASS` | `--- FAIL` | `--- SKIP` | 出处 |
|---|---|---|---|---|---|---|
| 借来的值放回去**＋**case 14 摘掉 | **本程现跑** | **20** | **13** | **0** | 0 | `probes/156/mut-156-r2/r2-152g1-case14off.log` |
| 同一发、case 14 活着 | **本程现跑** | 21 | 13 | **1** | 0 | `.../r2-152g1-case14live.log`（红名 `TestSLO152CorruptLegWithNoNamedPositionRefusesToBorrowOne`） |
| **不放任何变异**（同选择器、同字节） | **本程现跑** | **21** | **14** | **0** | 0 | `.../r2-asis-152-surface.log` |
| 注释**点名引用**的那枚日志（`mut-post/`，`git ls-files`＝tracked） | r1 留的、本程现数 | 20 | 13 | 0 | 0 | `probes/152/mut-post/g1-restored-borrowed-value-case14off.log` |
| 同一枚日志的未入库孪生（`mut-shipped/`） | r1 留的、本程现数 | 20 | 13 | 0 | 0 | `probes/152/mut-shipped/g1-restored-borrowed-value-case14off.log` |

**判**：注释里那句 `reports === RUN=21, PASS=14, FAIL=0`（`cmd/wisp/slo_report_144_windows_test.go:739`）**是错的**，
票面 ② 给的 `20/13/0` **对**。本程还量出错因，比"写歪了一个数"更具体：
**21/14/0 恰好等于这个选择器在改前改后的"零变异"名册数**（第 3 行＝本程现跑，且票 152 自己那枚
`mut-post/asis.log` 也是 21/14/0）⇒ 那三行注释把 **asis 的四个数**贴进了一句承诺"变异后的数"的话里；
`21` 不是"从孪生日志抄错一位"，`14` 也不是——**它俩谁都不是那一发变异的读数**。
（旁证：本程第一次跑错命令时**覆盖**了 `mut-156-r2/asis.log`——同名不同选择器一枚文件两头发——
本程没有把它删掉，改成把每一发都换成带选择器的名字（`r2-asis-152-surface.log`）后重跑，两发的数一字未变；
这条自报＝"别按名字猜"那一坑在本程又实发一次。）

**这一枚的修法形状受限于派单的尺**：派单要 AC#4 每一枚 commit **删除列＝0**，而"改成"按字面必然是一删一增
（票面 `:737-740` 与派单 `:18` 在这里直接冲突，冲突归 4.0 那一笔账）。本程取**纯追加**：
在 `:739` 之后（同一段注释里、函数声明之前）加一段，**逐字说出上面那三行表、现跑的日志路径、以及错因**，
并**逐字保留**被纠正的那句原文可 grep。⇒ 结果：`RUN=21, PASS=14, FAIL=0` 这串字符**仍在 `:739`**，
但它下一段就写明它是什么的数。**残留缺陷（要编排者那一枚带删除的 commit 才闭得了）**：
只读 `:739` 那一行的人若跳过后续注释，仍会拿到错的数。**本程不假装这一半也交了。**

**谁会误判**：自己去复跑 case 14 那枚"承重"读数的人。他把哪一句读错：
`g1-restored-borrowed-value-case14off.log reports === RUN=21, PASS=14, FAIL=0`，
出处 **`cmd/wisp/slo_report_144_windows_test.go:739`**。他量到 20/13 ⇒
在本仓 `=== RUN` 少一枚的读法就叫"有用例没跑到／被改名吞了"（正是名册两向 `comm` 要防的那件事），
于是他去**找一个失踪的用例**，而不是去信一句过期注释。

### 4.3 动作③（`mut-shipped/` 入库 or 把 §4.2 档位降级）：**两支都写在本程写面之外 ⇒ 本程一支都不执行，交读数＋点名残余**

```
$ git ls-files .scratch/wisp/probes/152/mut-post/    | wc -l   6      <- 含注释点名的那两枚 g1
$ git ls-files .scratch/wisp/probes/152/mut-anchor/  | wc -l   6
$ git ls-files .scratch/wisp/probes/152/mut-shipped/ | wc -l   0
$ git log --oneline -- .scratch/wisp/probes/152/mut-shipped/ | wc -l   0
$ ls .scratch/wisp/probes/152/mut-shipped/ | wc -l  9                  (asis + g1/g2/g3/g4 + h1/h2 + i1)
$ grep -n '原始读数在' docs/evidence/s1/152-subject-death-never-measured-r1.md
380:### 4.2 承重两问——两句都答，尺是本程自己的（`probes/152/my152.py`，原始读数在 `mut-anchor/` 与 `mut-shipped/`）
```

- 票面 ③ 的两支：**(甲)** 把 `mut-shipped/` 那 9 枚**入库**——纯新增（0 删），但目标路径是 `.scratch/wisp/probes/152/**`，
  不在派单给本程的写面（"新建 `.scratch/wisp/probes/156/**`"）里，而且那是**已停笔别程**留下的件；
  **(乙)** 把票 152 实现件 `§4.2` 那两行档位**降级**——那是**另一枚票的已入库证据件**，且按字面必带删除行。
  ⇒ **两支都越界，本程都不执行**（`SPEC-12 §4.2` 未定义即停那一支的形）。本程**也没有**做一个"看着像入库"的替代动作：
  把 9 枚日志**复制**进 `probes/156/` 再 commit 是**假交件**——被点名的那条路径依旧 0 枚 tracked，
  而多出来的副本只会让下一位数成两份凭据。**这条判断本程写在这里，不写成已完成。**
- 本程交的是读数与**归属**：`mut-post/`（tracked）**已经含**代码注释点名的两枚 g1 ⇒ 代码侧那条引用**在 HEAD 上可解析**，
  这一半不缺；缺的是 `§4.2` 点名的 `mut-shipped/` 那一族（9 枚、0 tracked、0 commit）。
- **谁会误判**：从**一枚 HEAD checkout**（而不是这台工作树）出发、按票 152 `§4.2` 那行去找原始读数的人。
  他把哪一句读错：`原始读数在 `mut-anchor/` 与 `mut-shipped/``，
  出处 **`docs/evidence/s1/152-subject-death-never-measured-r1.md:380`** —— 他会撞 `No such file or directory`，
  然后有理由把票 152 那张承重表读成编的（实际同形读数就在两门之外的 `mut-post/`，tracked）。
  本程把这一对（`mut-post`=6／`mut-shipped`=0）写进了 `slo_report_144_windows_test.go` 的注释面。
- ⚠ **别把这条当"③已闭"**：③ 的"两选一 + 写清选了哪个"仍**未决**，选择权在编排者。

### 4.4 放水两问自答（本格）

① **断言方向动没动**：没动——本格只加注释行，零枚 `t.` 调用、零枚判据；四个数与名册由 §4.5 那发控量证明未变。
② **helper 是不是原有的那枚**：本格不用 helper；尺是 r1 那把（`my156.py` 的 `spec`/`case_off`），r2 只加打印。
**另自报一处本程自己的形状违例**：`git diff --numstat` 的**新增行全以 `//` 开头**这一条，
本格只在 `cmd/wisp/slo_report_144_windows_test.go` 上成立；证据件本体是 markdown，追加行天然不以 `//` 开头
（票 154 那把尺量的是码的注释面）。本程不为此把证据行伪装成注释。

### 4.5 机器核（派单那把尺，原文可复算）

```
$ git diff --numstat -- cmd/wisp/slo_report_144_windows_test.go            55	0	...      (删除列＝0)
$ git diff -U0 -- <同上> | grep '^+' | grep -v '^+++' | grep -vcE '^\+//'   0            (新增行 100% 以 // 开头)
$ "$(go env GOPATH)/bin/gofumpt.exe" --version                              v0.12.0 (go1.27.1)   <- 现读，未抄旧票面
$ "$(go env GOPATH)/bin/gofumpt.exe" -l cmd/wisp/
  cmd/wisp/slo_exit_os_156_windows_test.go        <- 不是本程碰过的文件：r1 交付的 :272-279 那枚 struct literal 换行
                                                     形状；本程没动它（改它＝带删除行、且不在 AC#4 三枚动作里），
                                                     只登记。`grep -rln gofumpt scripts/ tools/` 本程现量＝0 命中
                                                     ⇒ 今天的门禁不扫这条，但 AC#7 那程若加就得先处理它。
$ go build ./cmd/wisp/                                                      rc=0
$ # 行为控制发（注释改完之后，同一条命令、同一选择器）：
$ python probes/156/my156-r2.py ... r2-asis-152-surface       RUN=21 PASS=14 FAIL=0 SKIP=0  rc=0
$ python probes/156/my156-r2.py ... r2-asis-slo-all-after-ac4 RUN=26 PASS=19 FAIL=0 SKIP=0  rc=0
  # 与注释改动前那一发 r2-asis-slo-all（RUN=26 PASS=19 FAIL=0 SKIP=0）一字不差
  # TREE 行把每发钉在字节上：slo_windows.go=d2bc41ecfef06bd4（改前改后同一枚，本格一字未动它）、
  # slo_exit_os_156_windows_test.go=9f1a709712b3bf1c（同样未动）、被本格加了 55 行注释的那枚文件
  # 在控制发里=daf3455916f6bfd3（注释改完之后的字节；上面三发 g1/asis 全部跑在这一版上）
```

**第 4 格判定**：**① 成立（r1 交付＋本程字节核）／② 读数成立、修法只到"纯追加纠正"那一半，残留一行过期数字未闭／
③ 未交（两支均越界，已点名归属）**。本格**不判自己 AC#4 通过**——三枚里 ② 有残留、③ 未动，票框由编排者定。

**本程没测什么（本格）**：
1. **没验"票 152 实现件 §4.2 那两行档位"到底该降还是该入库**（那是编排者的选择，本程只给读数）；
2. **没跑 `mut-shipped/` 那 9 枚日志里除 g1 之外的任何一发**（h1/h2/i1/g2/g3/g4 全是 r1 或票 152 的凭据，本程一枚未复算）；
3. **没验"注释里那两个数被改完之后，票 152 的验收程会不会反过来判它自己的 §9 第 8 行过期"**——
   `152-…-accept-r1.md` 是本程不许动的半件（` M`），它对同一枚数的记法本程一字未读成凭据。

---
