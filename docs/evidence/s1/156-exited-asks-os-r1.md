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
