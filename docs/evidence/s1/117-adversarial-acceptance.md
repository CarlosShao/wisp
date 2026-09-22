# 票 117 — 独立对抗验收（`acceptor-ticket117`）

被验：`.scratch/wisp/issues/117-security-warnings-have-no-listener-in-production.md`
本仓**第六次**"能力做出来了、装配根里没人调"族缺陷。票面六格 AC 由实现方**全部自勾 `[x]`**。

- 验收人：`acceptor-ticket117`（**只读**；唯一写件＝本文件；一行生产码/测试码都没改，见 §末自证）
- 开工 `date` 实测：**2026-09-22 10:53:38 CST**；收尾 `date` 见文末。每格读数各自带时间戳，不复用上一条。
- 被验 commit：`ce666ea`（5 枚代码文件）/ `23403ab`（票面入库）/ `e9d70ca`（接续者复算账）
- 实现方声称基线 `823d457`。**我实测 `ce666ea` 的父是 `2611558`**；`823d457..ce666ea` 的差集除票 117 五枚外只有
  `.../115-*.md` 与 `internal/winsec/notice_attribution_115_windows_test.go`（都不是 117 的改动）⇒ 基线口径成立。
- **我的控制组用 `2611558`（= `ce666ea` 的父）而不是 `823d457`** —— 这是最干净的同基座对照。
- 纯净快照（全部在仓外；仓库内**没有**建 worktree/checkout，A38④）：
  - `/tmp/ac117-gate` = `git archive e9d70ca | tar -x` + `cp third_party/sherpa-onnx/*.dll`
  - `/tmp/ac117-ctl` = `git archive 2611558 | tar -x`（控制组）
  - `/tmp/ac117-mut` = `git archive e9d70ca | tar -x` + dll（变异专用，四发都在这一份里做，做完各自还原）
  - 每把都 `ls -l go.mod` 证明文件真解出来了（`-rw-... 883 ...`），不是空目录假绿。

---

## 〇、开工前的仪器自检：两发**我自己的仪器**先假红了一次

派单点名 `git archive` 注入 CR、`cmd | grep x; echo $?` 量的是 `grep` 的 rc 这两类坑。**两发我都踩到了，都在读数之前纠正**：

**坑 1 —— "archive 注入 CR" 这一发我第一遍读成了假红，且方向相反。**

```
grep -c $'\r' cmd/wisp/logsink.go        →  193     ← 看着像"193 行全带 CR"
```

四把独立尺子重量，结论**相反**：

| 尺子 | 读数 |
|---|---|
| `head -c 40 cmd/wisp/logsink.go \| od -c` | `p a c k a g e   m a i n \n \n / /` ⇒ 行尾是**裸 `\n`** |
| `tr -cd '\r' < logsink.go \| wc -c` | **0**（文件 8770 字节） |
| `tr -cd '\r' < run.go \| wc -c` | **0**（28064 字节） |
| `grep -cP '\r'` | **0** |
| `git cat-file -p e9d70ca:cmd/wisp/logsink.go` 的 CR | **0**（blob 本身干净） |
| `git ls-files --eol` | `i/lf w/lf attr=text eol=lf` |

⇒ **本仓这个 commit 的 `git archive` 不注入 CR**：`.gitattributes` 把 `*.go` 钉成 `text eol=lf`，
per-path 属性盖住了 `text=auto` + `core.autocrlf=true`。第一发命中 193 的真实原因是
**Git Bash 把"只含一个 CR"的参数当空模式**（空模式匹配每一行）——**是我的仪器，不是产品码，也不是 archive**。
登记在此，为了下一位别再拿"archive 注入 CR"当本仓的读数；也说明我为什么后面所有按 `\n` 切块的用例都可信。

**坑 2 —— `GOOS=linux go vet ./cmd/wisp/ 2>&1 | head -4; echo $?` 量到的是 `head` 的 rc。**
我第一次读到 `rc=0`（假绿），当场改成接文件再量：真实 **rc=1**（3 行错误）。
⇒ 前任 §四 自述踩过同一发（`GOOS=linux go vet ./... | tail -5` 量成 0、真实 1）；**这条坑在本仓是活的**，我的读数全部走"先落文件再取 `$?`"。

**坑 3 —— dll 缺件。** `git archive` 里没有未入库的 `third_party/sherpa-onnx/*.dll`，
必须先 `mkdir -p` 再 `cp`（我第一次 `cp` 就撞了 "is not a directory"）。不补 dll 会红 4 行
（`TestSecretArgvCarriesNoSecret` 一族），那是**仪器缺件不是回归**。

**坑 4 —— MSYS 参数转换双向都咬人。** `icacls /grant` 要 `MSYS2_ARG_CONV_EXCL='*'`（否则报
`Invalid parameter "D:/work/soft/Git/grant"`）；**但同一条 `export` 会让 `taskkill //PID` / `//IM` 变成非法参数**
（`ERROR: Invalid argument/option - '//IM'`）。两种形状我都各踩一发，分开 shell 量才拿到正确读数。
另：`go build -o bin-a117/...` 的相对路径 + `APPDATA` 必须是 **Windows 形态**路径（给 MSYS 形态会静默落错地方）——
我全程用 `cygpath -w` 出来的形态。

---

## 一、裁决表（六格 × 三档证据）

| 格 | AC 一句话 | **判定** | 证据档位 |
|---|---|---|---|
| **AC#1** | 现状表：逐条点名"会出声的安全事件现在出到哪里"，全仓非测试口径 | **通过附条件** | 〔独立复现〕 |
| **AC#2** | `wisp run` 与 GUI 腿**两条路**都装上持久 sink，落点可 grep | **通过附条件（条件未了 ⇒ 本票不得翻 `-done`）** | 〔独立复现〕 |
| **AC#3** | 落点走票 95 私有目录口径；数据根解析不出来就拒绝，不许 fallback | **通过附条件** | 〔独立复现〕 |
| **AC#4** | 端到端真机/真进程证据，记录从盘上那个文件读出来 | **通过** | 〔独立复现〕 |
| **AC#5** | 噪声与体积量化读数 | **通过** | 〔独立复现〕（部分 〔日志＋归档，我抽验〕） |
| **AC#6** | 门禁四数 + gofmt/gofumpt/vet/d22scan | **通过附条件** | 〔独立复现〕 |

**六格里没有一格判"未复现"或"退回"**——理由与边界写在下面各格，其中 **AC#2 的附条件是阻断结案的**（§二 AC#2 的 M4）。

档位定义：〔独立复现〕我亲手跑过并拿到读数 / 〔日志＋归档，我抽验〕artifacts 在盘上、我抽查过关键条目 / 〔仅自述，不背书〕只有票面文字。

---

## 二、逐格正文

### AC#1 现状表 —— **通过附条件**〔独立复现〕

**我复算的口径**（不抄它的）：

```
grep -rn --include="*.go" -E "observe\.InitLog|InitLogWithRegistry|InstallAsDefault|slog\.SetDefault" . \
  | grep -v _test.go | grep -v third_party
```
非测试命中**恰好**是：`cmd/wisp/logsink.go:133/:137`、`cmd/wisp/slo_windows.go:237/:246`、
`internal/observe/logging.go:64/68/102/106/107`（定义与文档行）、`cmd/balldebug/main.go:118`（stderr TextHandler，**不是**持久 sink）
⇒ **与前任 §〇 的自述逐字一致**（除 `logsink.go` 是它自己新增的那四处）。

**事件名与五个字段名两侧逐字对得上**（首要攻击点，我自己 grep 的两边）：

- 生产者 `internal/winsec/winsec_windows.go:143`（`var noticeNarrowed` 体内，定义 `:135`）：
  `slog.Warn("winsec: seal cleared principals that stood on this object", "path", …, "kind", …, "cleared", …, "cleared_inherited", …, "policy", …)`
- 消费者：`git log --oneline ce666ea..HEAD -- internal/winsec/winsec_windows.go` **空** ⇒ 被验的码到今天一字未动，
  且我在**真进程写出的那一行 JSONL 上**逐字拿到了这六个格（§AC#4 与 §五），不是拿票面比票面。
- 该文件里 `slog.` **仅此 1 处**（`grep -c`），与前任自述一致。

**其余点名行的坐标我也逐条重量**（左＝票面，右＝我量到的）：

| 票面点名 | 我量到 | 结论 |
|---|---|---|
| `run.go:463` `auditf` / `:466` 走 sink | **`:463`** / **`:466`**（`rt.spec.sink.logger().Info("audit: " + line)`） | **未漂** |
| `internal/tools/bridge.go:586`（`tools: d31 report`）/ `:822`（`tools: PATH-ACCOUNT`） | **`:586` / `:822`** | **未漂** |
| `internal/config/manager.go:244` / `:247` 两条 WARN | **`:244` / `:247`** | **未漂** |
| `internal/secret/migrate.go:198` WARN | **`:198`** | **未漂** |
| `cmd/wisp/secret.go:471` WARN | **`:471`** | **未漂** |
| `internal/proc/shutdown.go:145/:147/:172` | **`:145/:147/:172`** | **未漂** |
| `internal/observe/goroutine.go:168` | **`:168`** | **未漂** |
| `internal/winsec/resolve.go:150/:157/:165/:181` | **`:150/:157/:165/:181`** | **未漂** |
| `internal/risk/winsec_c26.go:21` 在 `init()` 里装 resolver | `func init()` 在 **`:20`**，`winsec.SetPathResolver(c26Pipeline{})` 在 **`:21`** | **成立** |
| `run.go:347`（`perm: MODE-READ`）/ `:421`（`MODE-READ-FAILED`） | **`:355` / `:429`**（各漂 **+8**） | **漂移**：正是它自己那 18 行 install 块把它们推下的 |
| `internal/winsec/winsec_windows.go:143` | **`:143`** | 未漂（它 §〇 已自报"票面 `:87` 漂到 `:143`"） |

**⇒ 附条件的理由（三条，都是"现状表"这一格自己的职责范围没做满，不是别格的锅）**：

1. **漏了两个在 `cmd/wisp` 依赖图里、改后会被它装的听众直接带走的生产 ERROR 族。**
   我用 `go list -deps ./cmd/wisp` 量的可达性：`winsec`、`risk`、`plugin`、`tools`、`statemachine` **全在图里**。而它的表里**没有**：
   - `internal/statemachine/machine.go:139` `slog.Error("state machine rejected transition", …)` ——
     **这一条最像"安全事件"而不是"顺手噪声"**：`IllegalTransitionError` 就是状态机被要求走一步不合法的路（takeover/confirm 那一族），
     AC#1 括号里那句"至少"没点到它，但它符合"会出声的安全事件"的判据。
   - `internal/plugin/disposal.go:208` 与 `:332` 两条 `slog.Error`（dispose 链失败 / Defer-after-Dispose）。
   - `internal/observe/goroutine.go:271` `slog.Warn("goroutine outside the D38 roster (leak symptom)")` —— 它只列了同文件的 `:168`。
   - `internal/observe/logging.go:204` 与 `:307` `slog.Warn("observe: log retention sweep failed")` ——
     **日志管道自己失败的告警也进了日志**（这条是"听众自己喊救命"，形状特殊，表里一个字没提）。
2. **"顺手量到的非安全族"两个计数偏小**：它写 `internal/audio/*`（**4** 处）、`internal/ball/*`（**14** 处）。
   我用 `grep -rn --include="*.go" -E "slog\.(Warn|Error)" internal/ cmd/ | grep -v _test.go | grep -v third_party` 量到
   **audio 5 处**（`audio.go:136/152`、`mmdevice_windows.go:60/65`、`wasapimic_windows.go:112`）、
   **ball 17 处**（`ball_windows.go` 7、`hotkey_reload.go` 3、`hotkey_windows.go` 6、`renderer_windows.go` 1）。
   它的"前瞻不是读数"那句本身是对的（今天的 GUI 腿确实一条都不发——我 §AC#4 腿 B 的读数只有 2 条可作旁证），但**分母报错了**。
3. **两处 `+8` 行漂移**（上表倒数第二行）。它 §〇 专门立了"验收方量的坐标我逐条重量"这一节，
   却把**自己改动之后**才成立的 `run.go` 行号写成了改动前的坐标。措辞级，不影响任何判据。

**没抓到**的部分（据实说）：它的 9 行表**逐行都成立**，包括两条最难的自证——
第 5 行"`secret.MigratePlaintext` 生产零调用方"与第 9 行"`init()` 期那族任何安装点都追不上"，
后者我**在真进程里独立复现了**（见 §AC#4 末与 R-117-2）。
第 6 行"`wisp secret` 这条腿仍然只到 stderr"我也复算为真：非测试里 `installLogSink` 只有 3 个调用者
（`run.go:164`、`resident_windows.go:57`、`cmd/wisp/models.go:276`），`secret` 不在其中。

> **补一条它没有、我也量了的新事实**：`installLogSink` 今天已有**第三个生产调用者**
> `cmd/wisp/models.go:276`（票 121 `7fe5e73` 接的 `wisp models` 腿，**不是 117 的改动**）。
> 这**加强**而不是削弱本票：装配根现在有三条腿装了听众。三条腿各自 `main()` 互斥（`case "run"` / `case "models"` / 无参数），
> 所以**同一进程里不会装两次**；我也确认管道开档是 `O_CREATE|O_WRONLY|O_APPEND`（`internal/observe/logging.go:244`），
> 即便将来同进程装两次也**不会 truncate 掉已有内容**。这条是我主动找的"双装会不会互相打烂"，**结论为无害**。

**判**：表本身逐行真、口径真、坐标基本真；缺的是"现状表"这一格的**穷尽性**（漏 2 个可达的安全/生命周期 ERROR 族 + 非安全族分母报错）。
⇒ **通过附条件**：条件＝把上列 4 处补进表（纯纸面，零代码），措辞级，不阻断结案。

---

### AC#2 两条腿都装上 + 落点可 grep —— **通过附条件，且这一条阻断结案**〔独立复现〕

**先给结论性的好消息：这一格不是"只有测试能到"。** 我自己 grep 的三条非测试调用点：

```
cmd/wisp/run.go:164               installLogSink(s.dataDir)             ← wisp run 腿
cmd/wisp/resident_windows.go:57   installLogSink(rt.Layout.DataDir)     ← 无参数 GUI 腿
cmd/wisp/models.go:276            installLogSink(io_.dataDir)           ← 票 121 的第三条腿（非本票）
```

`main()` 的分派我读了：无参数 ⇒ `attachParentConsole(); runResident()`（`:53-54`）；`case "run":` ⇒ `cmdRun` ⇒ `runTextTask`（`:58-59`、`:110`）。
⇒ **票面 AC#2 明写的两条腿，两条都在真进程的装配根上**。我在 §AC#4 用**两条真进程**分别把它坐实了，
不是靠用例绿。

**"没有手搓 sink"这一条我复算为真**：`grep -n "installLogSink\|InitLog\|SetDefault\|teeHandler\|runTextTask\|CountLogFiles"` 打在
`logsink_windows_test.go` / `logsink_test.go` 上——三条 AC#2 用例（`:137`/`:234`/`:308`）**都调生产入口 `runTextTask`**
（`runTextTask(runSpec{…})` 在 `:152`、`:319`），读的是生产路径自己装出来的那本文件，
`observe.CountLogFiles`（`:61`）在 grep 之前先数一遍。`logsink_test.go:110` 调的是**生产函数 `installLogSink` 本体**，不是自造 handler。
⇒ 这一族缺陷在本票**没有复发于用例层**。

**落点可 grep 的判据我逐字复跑为真**（§AC#4 那次真进程的盘上文件）：

```
grep -F -c '"msg":"winsec: seal cleared principals that stood on this object"' <数据根>\logs\wisp-20260922-001.jsonl  →  2   (grep_rc=0)
```
命中行同级格 = `"level"`/`"path"`/`"kind"`/`"cleared"`/`"cleared_inherited"`/`"policy"` **六格全在**（我把 7 个键名逐个 `grep -oF` 量过，全 YES）。
`"msg":"audit: perm: MODE-READ` 族我也在同一本文件里逐字拿到了（§五 的 msg kinds 统计）。

**"反半边不被牺牲"我也在真进程里量到**：腿 B 的常驻进程控制台文件 `/tmp/ac117-legB-A.txt` 里同时有
`time=… level=INFO msg="wisp: persistent log sink installed" …`（**镜像，TextHandler 形态**）与
`time=… level=INFO msg="activation requested by second launch …"` ⇒ **装文件没有把控制台弄哑**，
这条它只有用例（`TestAC3InstallingTheFileSinkDoesNotSilenceTheConsole`），我有真进程读数。

**⇒ 附条件（阻断）的理由 = 变异 M4，这是本票最硬的一条读数：**

`grep -rn "runResident\|resident_windows" cmd/wisp/*_test.go` ⇒ **命中 0 个测试文件**。
**GUI 腿那 20 行 install 没有任何钉。** 我在纯净快照里把它整段删掉：

- 落地证明：`grep -c installLogSink cmd/wisp/resident_windows.go` = **0**（删前 1），文件从 105 行变 **76** 行；
- `go build ./...` **rc=0**、`go vet ./cmd/wisp/` **rc=0**；
- `go test -count=1 ./cmd/wisp/` ⇒ **`ok … 60.987s`，`--- FAIL` 命中 0，rc=0** ——
  **全套 146 条用例在 GUI 腿的听众被整段删光之后全绿。**

**这正是本票立票的那一族缺陷换了个位置**：票面 AC#2 要"两条路都装上"，今天两条都装着（我用真进程各验了一次），
但**只有 run 腿装了钉**。谁下一次重构 `runResident`（票 07 接球、票 35 接面板都会碰它）把这 20 行挪丢，
门禁一行都不会红，而 owner 双击图标那条腿**恰恰是唯一一条没有终端可看、因此唯一真正需要持久听众的腿**。
票 121 为同一族缺陷造的仪器就在隔壁：`internal/models/assembly_reachability_121_test.go`。

**判**：**通过附条件**。条件＝**在 `cmd/wisp/` 内补一条钉住"常驻腿装了持久 sink"的用例**（形状可以照
`TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite`：驱动 `runResident` 的可注入前半段，或用装配可达性/AST 断言
`runResident` 体内调用 `installLogSink`）。
按派单那条硬规矩复核："`告警从来没人看见` 是否还在？"——**不在 shipped 码里**（腿 B 我拿到了盘上读数），
所以**不判退回**；但**条件未了 ⇒ 本票不得翻 `-done`**（见 §六）。

---

### AC#3 落点 —— **通过附条件**〔独立复现〕

**判据按编排者给的那句话执行**："**有没有任何一条路能把 jsonl 写到数据根之外**"。**不许因为 `Q-31` 没答就判待拍板。**

**`Q-31` 我先自己查了**（不是抄它）：`grep -n "Q-31" docs/reports/pending-and-issues.md` ⇒ 命中 **10 行**
（`:22` 顶部 `[H7]`、`:990` Q 表、`:1242/:1269/:1312/:1340/:1374/:1397/:1404`），
**无一被划掉、无"已答"字样** ⇒ owner 仍未拍板；现行裁定＝**建议封**，实现按"封"的保守形状做 ⇒ 与我的现行裁定一致。
它**没有**在任何文档/注释/用例里把"日志属于私有数据"写成 owner 已认定（我读了 `logsink.go:39-45`
"…is the one question of ticket 95's that owner has NOT answered" 与 `logsink_test.go:6-12`
"sealing the log file is a decision this ticket must not make"）⇒ **这条纪律守住了**。

**两级 grep 我复算并自己找漏。** 它的两级：

- `cmd/wisp/logsink.go`：`installLogSink` 只认 `dataDir`；`logSinkDir = filepath.Join(dataDir, logDirName)`（`:71-73`，`logDirName="logs"` `:60`）；
  `dataDir==""` ⇒ **直接 error**（`:129-131`）。全文件 `TempDir|MkdirTemp|fallback|defaultDir|UserHomeDir` 命中 **0**。
- `internal/observe/logging.go`：`InitLogWithRegistry`（`:68`）里 `TempDir|MkdirTemp|fallback|defaultDir|UserHomeDir` 我量到
  **只剩 `logging.go:316` 注释里的 "the mtime is the fallback"**（一句关于**文件名解析**的注释，与落点无关），
  `cfg.Dir==""` ⇒ error、`MkdirAll` 失败 ⇒ error ⇒ **管道本体没有"给定目录打不开就换个地方写"**。
- 与 `wisp slo` 同形：`cmd/wisp/slo_windows.go:238` 也是 `Dir: filepath.Join(rt.Layout.DataDir,"logs")` + `Level:"info"` ⇒ **一个 env 一本日志、一套轮转**。

**我自己找漏的四发实跑**（都在仓外临时根，全程没动 owner 的真实数据目录）：

| 探法 | 命令形状 | 落点读数 | 跑出数据根？ |
|---|---|---|---|
| 正常 dev 根 | `APPDATA=<仓外> WISP_ENV=dev … run` | `<仓外>\wisp-dev\logs\wisp-20260922-001.jsonl` | **否** |
| **空 dataDir** | 单元层已被 `TestAC3EmptyDataRootIsARefusalNotAFallback` 钉死；生产里 `installLogSink` 拿不到 `""` | — | 否 |
| **`APPDATA` 未设（run 腿）** | `(unset APPDATA) … run` | **`dir=wisp-dev\logs`** ⇒ 文件落在 `<CWD>\wisp-dev\logs\…jsonl` | **否（跑出了"%APPDATA 之下"，但仍在它自己解析出的数据根之内）** |
| **`APPDATA` 未设（GUI 腿）** | `(unset APPDATA) …`（无参数） | **rc=1**，`wisp: boot failed: proc: user config dir: %AppData% is not defined` | **否：这条腿直接拒绝启动** |
| **`WISP_*` 环境覆盖** | `WISP_ENV=test …run` | `dir=C:\Users\swq\AppData\Local\Temp\wisp-test-10188\logs` | **否**（`test` env 的数据根本来就是 TEMP，SPEC-03 §5） |

⇒ **按 AC#3 那句判据（jsonl 会不会跑出数据根）：我造不出逃逸。这一格该红的没红。**

**附条件的两条（都是真读数，不是措辞）**：

1. **`APPDATA` 未设时数据根整棵搬到 CWD 相对路径**——`cmd/wisp/doctor.go` 的 `resolveDataDir`（run 腿与 models 腿用它）里有一句
   `base, err := os.UserConfigDir(); if err != nil { base = "." }`。
   我 `git log -S'base = "."' --oneline -- cmd/wisp/doctor.go` 量到它来自 **`bcc892c`（S0 建票 12/doctor 那批）**，
   **不是 117 引入的**。而且**私有数据整族跟着搬**：`find /tmp/ac117-cwdtest -type d` 同时出现 `wisp-dev\logs` **和 `wisp-dev\secrets`**，
   事后 `icacls …\ac117-cwdtest\wisp-dev\secrets` = **只有 SYSTEM/Administrators/swq 三条 `(F)`(+3 条 `(OI)(CI)(IO)`)** ⇒
   `secrets` 在新位置**照样被封过**。所以 117 没引入不对称：**它只是让安全日志也多长了一份在那个位置、且那份没有密封**。
   真实代价：`logs\` 落在 CWD（双进程的 CWD 可能是任何地方）且**未密封**。⇒ **另开一张票**（R-117-B，地界在 `resolveDataDir`，
   比日志大得多：它同时挪动 `config.toml`、DPAPI 目录与 `memory.db`）。
2. **它钉的那条不变式，与它自述"不 fallback"那句的覆盖面不一致。** 我把 `installLogSink` 的 `""` 拒绝改成
   `dataDir = filepath.Join(os.TempDir(), "wisp-log-fallback")`（变异 M3，落地证明：`grep -c TempDir logsink.go` = **1**、`go build` **rc=0**）：
   `TestAC3EmptyDataRootIsARefusalNotAFallback` **红**（`logsink_test.go:72` "installLogSink("") returned no error: an unresolved data root must be a refusal"）——
   **钉有牙**；但同一次里 **`TestAC3LogSinkLandsInsideTheEnvDataRoot` 全绿**，即
   **"落在数据根之内"那条用例看不见"落点被换成 TempDir"**（它自己传的是非空 dataDir）。
   结合 (1)：**真实可达的那条 fallback 在 `installLogSink` 的上游一层，所以 `cmd/wisp` 里任何用例都看不见它。**

**顺手把"封不封"那一面的代价复量了**（我这一发的，不是抄它的）：
我自己的 hermetic 根同样建在 `%TEMP%` 下（票 95 的规矩），同一本 jsonl 读到的是

```
…\ac117-appdata\wisp-dev\logs\wisp-20260922-001.jsonl
   DESKTOP-LVS7839\CodexSandboxUsers:(I)(M,DC)
   S-1-5-21-3623186960-731165060-4091685855-1717338598:(I)(M,DC)
   NT AUTHORITY\SYSTEM:(I)(F)  BUILTIN\Administrators:(I)(F)  DESKTOP-LVS7839\swq:(I)(F)
```
而同根下被 `PrivateDirAll` 封过的 `secrets\` 只剩 **SYSTEM/Administrators/swq 三条 `(F)`**
⇒ **"这类主机上别的账户能读、能改、能删这本安全日志"这句我复现成立**（它 §三 第 2 条不是空口）。
而 owner 真实根 `%APPDATA%\wisp` 与 `%APPDATA%\wisp-dev` 我**只读**重量：两处都恰好
`(I)(OI)(CI)(F)` × SYSTEM/Administrators/swq **三条**，与它 §三 第 1 条逐字一致；
并且**两处都没有 `logs` 目录**（`ls -d …\logs` ⇒ No such file or directory）
⇒ **"没人往 owner 的真实数据目录写过东西"这一条我拿到的是"缺失即证据"级别的读数，不是自述**。

**判**：**通过附条件**（条件＝R-117-A/R-117-B 立票；纸面与跨票，不在本票码内地界）。
**没有把 `Q-31` 当挡箭牌**：这一格要的是落点，落点我四发都钉在数据根之内；它也没因为等不到答复就改成"不封"。

---

### AC#4 真机端到端 —— **通过**〔独立复现〕

**我换了一个 APPDATA 根自己重造了一次带外授权清除**（不是复用它 22:53 那一发）：

仪器：`/tmp/ac117-gate` 里 `go build -o bin-a117/wisp-a117.exe ./cmd/wisp`（**没碰仓库的 `wisp.exe`**；
事后 `find . -name "wisp-a117*" -not -path "./.git/*"` 命中 **0**，`git status` 里仓库没多出 exe）；
`APPDATA=C:\Users\swq\AppData\Local\Temp\ac117-appdata`（**仓外**）+ `WISP_ENV=dev` ⇒ 数据根 `<appdata>\wisp-dev`；
**全程没打开任何 GUI 窗口**（自证见下面腿 B）。

**腿 A = `wisp run`，三件全做**：

1. **run #1**：`runA1_rc=`**`2`**（Unconfigured，无 `config.toml`）——它自己就把 `%TEMP%` 父目录带进来的**两条外来继承 ACE** 清了，
   `kind":"inherited"` 那条 WARN 当场落盘（盘上 3 条 / **1315 字节**）。
2. **制造带外授权**：`icacls <根>\secrets /grant *S-1-1-0:(OI)(CI)(RX)` ⇒ `grant_rc=0`。
   BEFORE = 只有 SYSTEM/Administrators/swq 各 `(F)` + `(OI)(CI)(IO)(F)`；
   **AFTER-GRANT 明确多出 `Everyone:(OI)(CI)(RX)`**（我贴的 icacls 原文，不是推断）。
3. **run #2**：`runA2_rc=`**`2`**；`grep -F -c '"msg":"winsec: seal cleared principals that stood on this object"'` ⇒ **2**（`grep_rc=0`）。
   **我这一发的盘上原文（逐字）**：
   ```
   {"time":"2026-09-22T10:54:27.7202401+08:00","level":"WARN","msg":"winsec: seal cleared principals that stood on this object","path":"C:\\Users\\swq\\AppData\\Local\\Temp\\ac117-appdata\\wisp-dev\\secrets","kind":"explicit","cleared":"S-1-1-0(A;OICI;0x1200a9;;;WD)","cleared_inherited":"","policy":"winsec owns the grants on this tree; out-of-band ACEs are removed at the next seal"}
   ```
   **六格逐个点名**（`grep -oF` 单键量，全 YES）：`"level":"WARN"` / `"msg"` / `"path"`=**我那个仓外根的 secrets** /
   `"kind":"explicit"`（我这次只塞了**一条显式** ACE，所以不是前任那发的 `explicit+inherited`——**不是数据有出入**）/
   `"cleared":"S-1-1-0(A;OICI;0x1200a9;;;WD)"`（**就是我 `icacls` 塞进去的那条 Everyone**）/ `"cleared_inherited":""` /
   `"policy":"winsec owns the grants on this tree; out-of-band ACEs are removed at the next seal"`。
4. **确认清除真发生**：run #2 之后 `icacls <根>\secrets` = SYSTEM/Administrators/swq 三条 `(F)` + 三条 `(OI)(CI)(IO)(F)`；
   `icacls … | grep -c Everyone` = **0** ⇒ **`Everyone` 真的从 ACL 消失了**。
   ⇒ **这条记录描述的是真发生了的清除，不是"打算清除"**——AC#4 这一格要的形状，我拿到了。

**腿 B = 无参数 GUI 腿（一次真进程，并且确实没开窗）**：
常驻 A（PID **36856**）→ `Get-Process wisp-a117 | Select Id,MainWindowTitle,MainWindowHandle` ⇒
**`MainWindowTitle` 为空、`MainWindowHandle` = `0`**（不是"我绕过了不许开窗"，是这条腿今天只有 Job Object + 单实例 + 空事件循环）。
第二实例 ⇒ 打印 `wisp: another instance is running in this session; activated it; exiting`，`second_rc=`**`0`**；
A 的事件循环因此出声，而 A **没有可看的屏幕**，盘上多出的两条（A 还活着、没走 Close 时就已在盘上）：
```
{"time":"2026-09-22T10:56:02.33821+08:00","level":"INFO","msg":"wisp: persistent log sink installed","dir":"C:\\Users\\swq\\AppData\\Local\\Temp\\ac117-appdata\\wisp-dev\\logs","min_level":"info"}
{"time":"2026-09-22T10:56:33.0981657+08:00","level":"INFO","msg":"activation requested by second launch (ball bring-to-front lands with ticket 07)"}
```
`Get-Process -Id 36856 | Measure-Object` = **1**（量这两条时 A 仍在跑）⇒ **不是退出路径替它兜出来的**。
测毕 `taskkill /IM wisp-a117.exe /F` ⇒ 复查计数 **0**，没留常驻进程。

**腿 B 采不到的两半，我复核"属实且没被掩盖"**：
- 我**自己复现了那堵仪器墙**：`taskkill /PID 36856`（不带 `/F`）⇒ **rc=1**，错误原文
  `ERROR: The process with PID 36856 could not be terminated. Reason: This process can only be terminated forcefully (with /F option).`
  期间 `lines_before=8 → lines_after=8`（关停记录确实一条都没进文件）。
  ⇒ **D38(e) 十步的 `shutdown step …` 记录仍未采到**，缺的动作＝**owner 签收窗口里按一次 Ctrl+C**
  （常驻进程自己的 stdout 文件里那句 `Ctrl+C exits cleanly` 就是它给的入口）。**这是纸面缺口不是假绿。**
- **GUI 腿今天出不了 winsec WARN**（这条腿没有任何密封点：不读 config、不开 store）。
  我能证的与它一样，只有"同一个 `slog.SetDefault` 扇出在这条腿上也装上了、控制台外确实多了一本文件"——**这一句我有真进程读数（上面那两条 + A 的控制台镜像文件）**。

**第 3 件 = 剥光 PATH**：`PATH="/usr/bin:/bin" ./bin-a117/wisp-a117.exe run …` ⇒ **`stripped_rc=127`**，原文
`C:/Users/swq/AppData/Local/Temp/ac117-gate/bin-a117/wisp-a117.exe: error while loading shared libraries: sherpa-onnx-c-api.dll: cannot open shared object file: No such file or directory`；
这一发之后 `wc -l` 那本 jsonl **仍是 6，一条没多** ⇒ **加载期起不来的进程连 sink 都没机会装**。
⇒ **接续者这一件我从头到尾自己跑了一遍**，前任那发 `rc=127`"说过没留证据"的形状现在对我也不成立了——**我有自己的 artifacts**。

**R-117-2 我独立复现了**（不靠它的 `run-leg-3.txt`）：我腿 A 控制台文件的**第一行**是
`2026-09-22 10:53:38 INFO winsec: sealing path resolver installed resolver=risk.c26Pipeline probes_passed=2`
（**stock 默认 logger 的 `YYYY-MM-DD HH:MM:SS` 形态**，与 install 之后那两条 `time=…level=…` 的 TextHandler 形态**肉眼可分**），
而同一次 run 的 jsonl **3 条里没有它** ⇒ **`init()` 期出声的记录任何安装点都追不上，登记属实**。

**判**：**通过**〔独立复现〕。这一格原本最薄，且前任确实被抓过一次"说过没留证据"；
我这把是**全新一个 APPDATA 根、全新一条 Everyone ACE、全新一次逐字比对**，六格全中、清除可验、rc=127 与零增长可验。
两半采不到的地方它照实标了，我也照实复核了那两半是真的采不到（**包括我自己撞了同一堵 `taskkill` 墙**）。

---

### AC#5 噪声与体积 —— **通过**〔独立复现〕（13 条那一发为〔日志＋归档，我抽验〕）

**逐条量，不用平均抹尾巴；我没动过任何阈值/轮转/级别。** 单位：**字符数 = `awk length($0)+1`**；`wc -c` 单列（audit 那一格有中文，UTF-8 一字三字节，所以 `wc -c` 更大——它这句我自己复现了）。

**样本 1（腿 A run #1，3 条 / `wc -c`=1315 字节）**
`199`（install INFO）/ **`591`**（seal WARN，`cleared_inherited` 里两条外来 SID —— **这就是我这一把的最长单条**）/ `455`（`audit: perm: MODE-READ-FAILED` INFO）

**样本 2（腿 A run #2，3 条 / 增量 1107 字节）**
`199` / **`383`**（seal WARN，只有一条显式 `S-1-1-0`）/ `455`

**样本 3（腿 B 常驻启动，2 条 / 346 字符）**
`197`（install INFO）/ `149`（activation INFO —— **全表最短**）

**样本 4（`WISP_ENV=test` 一次 run，3 条 / `wc -c`=1282 字节）**
`191` / `584` / `437`

**样本 5（并发探针：常驻活着 + 连着 5 次 `wisp run`，同一本文件）——这一发是我自己加的**
**13 条 / 4508 字节**，逐条：`196 149 196 589 449 196 449 196 449 195 449 196 449`；
`msg` 种类统计：`wisp: persistent log sink installed` × **6**、`audit: perm: MODE-READ-FAILED` × **5**、
`winsec: seal cleared principals …` × **1**、`activation requested by second launch` × **1**。
⇒ 每进程 **1 条 install**；5 次重复 run **各 1 条 audit**；**seal WARN 只有 1 条**——
**它是自限的**：只有"真有东西被清"才出声，第二次启动发现树已经窄了就不再喊。
⇒ **体积预算**：`4508 字节 / 6 次启动 = 751 字节/次`（含一次常驻与 5 次 CLI）。

**JSONL 完整性（我加的那发攻击）**：两个进程**同时**往同一本文件追加，这在装听众之前不可能发生
（改前只有 `wisp slo` 会写）。我用 `python json.loads` **逐行验**：
`checked=13 bad=0`、`checked=8 bad=0` ⇒ **没有被撕开/交错的半行**。
机制我读了：`internal/observe/logging.go:244` 是 `os.O_CREATE|os.O_WRONLY|os.O_APPEND`，
`flushInterval = 500 * time.Millisecond`。⚠ **边界照实说**：我这一发只测了同一本日内 6 个进程；
**跨 10MB 换档那一瞬两个 writer 各自 roll 的竞态我没测**（既有形状，`wisp slo` 时代就有，不是 117 引入的新面积）。

**与 `wisp slo` 既有轮转/上限策略的关系——我自己量的常数**：
`LogConfig` 只传 `Dir`/`Level`（`logsink.go:133`），`RollSizeMB`/`RollDays` **留 0**；
`internal/observe/logging.go:46-49` 明写 `0 = 10MB, schema default` / `0 = 7 days, schema default`；
`flushInterval = 500ms`（`:29`）；单格长度上界 `observe.MaxLoggedString = 512`（`redact.go:36-37`，`truncate(s, MaxLoggedString)` 在 `:141`）
⇒ **没另开一套、没把级别或轮转调松**（它 §五 那句"同一套、一字未改"我复核为真）。
`Level:"info"`（`logsink.go:66`）= schema 默认 = `wisp slo` 今天的传法；winsec 那条是 **WARN**，`info` 收得住，
**级别往严里调（warn/error）反而会留得住它、但会丢掉 audit 那族 INFO**——它没调，我也没调。

**我自己重算的预算**（不抄它的"3,400 次"）：
- 按我的 Unconfigured 启动 **1315 字节/次** ⇒ **10MB ÷ 1315 ≈ 7,970 次**才滚一次；
- 按样本 5 的 **751 字节/次** ⇒ ≈ **13,900 次**；
- 它自述的完整成功一次 run = **13 条 / 3131 与 3134 字节** ⇒ ≈ **3,330 次**（这一发**我没重跑**：要一份能过 `config.toml` 校验的真配置 + mockllm + 一次 host `fs.read`。
  属〔日志＋归档，我抽验〕：我在快照里以 `-count=1` 跑过装它的 `TestAC2AuditTrailLandsInTheRunLegLogFile`，**它绿**，
  且 M1 把它打成红时它明写"文件里只剩 5 条"这种数量断言 ⇒ **该用例确实在数条数**，只是我没把那个 13 抄成我的读数）。
- **最坏一发的形状**：4 个密封点（`internal/secret/store.go:49` secrets、`internal/memory/open.go:176/188/503` db/artifacts/backups，
  另 `internal/secret/migrate.go:174` 每个备份一枚 `SealFile`）若同时都需要收窄 ⇒ 每次启动**最多约 5 条 seal WARN**，
  按我最长单条 **591 字符**上界估 ≈ **3.0KB 一次启动** ⇒ 10MB 仍容 ≈ 3,400 次。**这是算术不是测量**，别当读数。

**前瞻那一格它自己标了"未测，别当读数"**（球/音频接上之后 `internal/ball` 逐帧可能重复的 `slow drawFrame`），我复核同意，
并且**它数错了分母**（17 不是 14、5 不是 4，见 AC#1 第 2 条）——**这反而让那句前瞻更该被认真对待**，
所以我把"票 07/34 接线时补一次最坏每秒几条"记进 R-117-D。

**判**：**通过**。多样本全报、尾部（591）没被平均抹掉、没调阈值、轮转关系用常数量了、并发追加的完整性我自己攻过一轮。

---

### AC#6 门禁 —— **通过附条件**〔独立复现〕

**基座：`/tmp/ac117-gate`（= `git archive e9d70ca`，即票 117 三枚 commit 全在里的树）+ dll。控制组：`/tmp/ac117-ctl`（= `git archive 2611558` = `ce666ea` 的父）。**

`go build ./...` ⇒ **rc=0**（快照内，含我的 dll 复制）。

**四数（`-count=2` 不走缓存；四数全部用 `-v` 量的）**：

| 包 | rc | `=== RUN` 全部 | `--- PASS` 顶层 | `    --- PASS` 子用例 | `--- FAIL` | `--- SKIP` |
|---|---|---|---|---|---|---|
| `./cmd/wisp/` | **0** | **146** | **78** | **68** | **0** | **0** |
| `./internal/observe/` | **0** | **94** | **94** | 0 | **0** | **0** |

`78+68=146` 两口径自洽；`ok … 110.595s` / `ok … 2.583s`。
⇒ **与前任自述、与接续者复算的 146/146/0/0 与 94/94/0/0 逐字一致**，这次是我量的第三把。

**四发假绿我逐个排（不抄它的说法）**：

1. **`-count=2` 才不缓存** —— 我用的就是 `-count=2`（另一发 `-count=1` 只在下面"非 `-v`"那发与变异里用，各自标明）。
2. **非 `-v` 既不印 PASS 也不印 SKIP** —— 我实跑 `go test -count=1 ./cmd/wisp/`（无 `-v`）：
   输出只有 `ok github.com/CarlosShao/wisp/cmd/wisp 45.061s`，
   `grep -c 'PASS'` = **0**、`grep -c 'SKIP'` = **0**，rc=0
   ⇒ **"ok" 不等于任何一条用例真的跑过**；我表里那四个数全部来自 `-v`，这一点被我自己坐实。
3. **`GOOS=linux go vet` 只编译不执行** —— 承认，并且**我抓到了它这一格的一处over-claim**（见下面"附条件"）。
4. **Git Bash 下 `docker run -v "C:\…"` 静默挂空且 rc=0** —— 本票六格里**没有任何一条读数走 docker**，
   所以我没用容器复算（`docker version` 本机 **29.6.2** 可用，我确认过；用它只会引入这一族新坑）。
   受影响的是票 118 那一族 Linux 容器读数，不归本票。

**格式与静态**：
`gofmt -l cmd/wisp internal/observe` ⇒ **rc=0，输出空**；
`"$(go env GOPATH)/bin/gofumpt.exe" -l cmd/wisp internal/observe` ⇒ **rc=0，输出空**，
本机该二进制**存在**（`-rwxr-xr-x … D:\work\base\gopath/bin/gofumpt.exe`），`--version` 原文 **`v0.7.0 (go1.27.1)`**
⇒ **不是"未跑"**；派单提醒过 `logsink_windows_test.go` 曾被 `gofmt` 点名，**我这把快照里也不再被点名**，
与接续者"我没为格式改过一行"那句一致。
`go vet ./cmd/wisp/ ./internal/observe/` ⇒ **rc=0**。

**`GOOS=linux go vet ./...` ⇒ 我的 rc=1**，错误原文（快照 `e9d70ca`）：
```
package github.com/CarlosShao/wisp/cmd/wisp
	imports github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx
	imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in D:\work\base\gopath\pkg\mod\github.com\k2-fsa\sherpa-onnx-go-linux@v1.13.8
```
**控制组 `2611558`（不含票 117 任何一枚）跑同一条命令 ⇒ rc 同为 1，`diff` 两份输出 ⇒ 完全一致（`IDENTICAL=yes`）**。
⇒ **不是 117 引入的**。**这一格正是接续者 §六 自己点名的"唯一没重量的一格"——我重量了，它成立。**

**收尾 `sh scripts/d22scan.sh`：纯净快照 rc=0 / clean**，台账**逐 scope 对照我的控制组**：

| scope | 控制组 `2611558` | 被验 `e9d70ca` | 判 |
|---|---|---|---|
| bans #1-5 `internal/` | 202 | 202 | 平 |
| bans #1-5 `cmd/` | 20 | **21** | **+1 = `logsink.go`（生产）** |
| ban #6 `frontend/` | 40 | 40 | 平 |
| ban #7 `internal/tools/` | 18 | 18 | 平 |
| ban #8 `design/` | 16 | 16 | 平 |
| ban #8 `frontend/` | 40 | 40 | 平 |
| ban #8 `internal/` | 373 | **374** | **+1 = 票 119 的 `dataroot_symlink_119_other_test.go`（`189cb1e`，不是我看的这枚 commit 的，也不是降）** |
| ban #8 `cmd/` | 26 | **29** | **+3 = `logsink.go` + 两枚 `_test.go`** |
⇒ **没有任何 scope 下降**；两处上升都能逐枚归到文件，且 `cmd/` 那两笔与接续者自述逐字相同。
`tools/d22scan/**`、`allowlist.txt`、任何阈值/golden：**我只读，一字未碰**（本文件末尾 `git status` 自证）。

**PATH 依赖如实登记**：本票所有"能跑"都带 `PATH=$PWD/third_party/sherpa-onnx:$PATH` 这个前提（票 98 的洞还在），
**没有一句被我说成"不依赖私有 PATH"**；剥光 PATH 那一发的 rc=127 与"日志一条不增"都在 §AC#4。

**⇒ 附条件（我抓到的唯一一处真 over-claim）**：票面 §六 写着
> `GOOS=linux go vet ./internal/observe/` **rc=0**（**我的新生产文件是 untagged 的，这一格证明它在 linux 也编得过**）

**这句的括号不成立**，两条独立读数：
- `GOOS=linux go list -deps ./internal/observe/ | grep -c cmd/wisp` = **0** ⇒ **`logsink.go` 根本不在这条命令的依赖闭包里**，
  它 rc=0 只证明 `internal/observe` 在 linux 编得过，**一个字都证明不了 `logsink.go`**。
- 唯一能覆盖它的命令 `GOOS=linux go vet ./cmd/wisp/` ⇒ **rc=1**（上面那发 sherpa build-constraints，控制组同形）。
  （我第一次也用 `| head -4; echo $?` 量成了 rc=0，改成接文件才拿到 1 —— 见 §〇 坑 2。）

顺带一个形状事实：`logsink_test.go` 那 3 条 AC#3 用例**没有 build tag**，看着像"跨平台"，
但它的包在 linux **根本编不出来**，而 CI 里 `cmd/wisp` 的用例只在 **`windows-latest`**（`.github/workflows/ci.yml`：`--scope=cli` 由
`scripts/wisp-cli-tests.sh` 跑，step `:389`，该 job `runs-on: windows-latest` `:302`）
⇒ **这 3 条钉 AC#3 落点的用例在任何 CI 平台上都只在 Windows 执行**，untagged 给不出 linux 覆盖。
**后果评估（据实，别夸）**：`logSinkDir` 只是 `filepath.Join(dataDir, "logs")`，平台中立，真红风险低；
**该修的是那句"证明它在 linux 也编得过"**（措辞级），以及别让"untagged"被下一位当"跨平台已测"。

**判**：**通过附条件**（条件＝把 §六 那句 linux 措辞改成"`logsink.go` 的 linux 可编译性**未证**，`GOOS=linux go vet ./cmd/wisp/` rc=1 且与父 commit 逐字同形"）。
四数、格式、vet、台账、控制组、假绿四发——**全部我自己在纯净快照里量的，且比它多做了一把 `2611558` 控制组**。

---

## 三、变异清单（四发，全在 `/tmp/ac117-mut` = `git archive e9d70ca` 纯净快照里；**仓库内一字未改 ⇒ 无需 revert**）

每发都走"改一处 → **同一条 `&&` 链里** grep 那几字节 → `go build` rc=0 → 才读红名"。红名都点到用例自己。

| 编号 | 改什么 | 落地证明（同链 grep） | build | 红名（用例本名 + 失败原文要点） | 结论 |
|---|---|---|---|---|---|
| **M1** 时机 | 把 `installLogSink` 整段挪到 `assembleRuntime` **之后**（`grep -n` ⇒ assemble 在 `:158`、install 在 `:165`） | 行号对调可见 | **rc=0** | `--- FAIL: TestAC2SealNoticeLandsInTheRunLegLogFile`<br>`--- FAIL: TestAC2AuditTrailLandsInTheRunLegLogFile`<br>`--- FAIL: TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite`<br>（`-run 'TestAC2|TestAC3'` 六条里 **3 红 3 绿**，绿的正是三条 AC#3）；日志里能看见 seal WARN **排在** `wisp: persistent log sink installed` **之前** | **有牙**。它自述的 M1/M2 我复跑为真 |
| **M2** 反半边 | 把 `teeHandler.Handle` 里 mirror 那三行摘掉（只留 `primary`） | `sed -n '177,180p'` ⇒ 函数体只剩 2 行 | **rc=0** | `--- FAIL: TestAC3InstallingTheFileSinkDoesNotSilenceTheConsole`<br>`logsink_test.go:134: the console saw nothing of the WARN`<br>`logsink_test.go:137: the console line lost the level the stock default logger used to print` | **有牙**。**这一发前任明写"我没有实跑"** ⇒ 现在跑过了，它猜的红名是真的 |
| **M3** 拒绝 vs fallback | `dataDir==""` 从 `errors.New(…)` 改成 `dataDir = filepath.Join(os.TempDir(),"wisp-log-fallback")` | `grep -c TempDir logsink.go` = **1** | **rc=0** | `--- FAIL: TestAC3EmptyDataRootIsARefusalNotAFallback`<br>`logsink_test.go:72: installLogSink("") returned no error: an unresolved data root must be a refusal`<br>**但 `TestAC3LogSinkLandsInsideTheEnvDataRoot` 与 `…DoesNotSilenceTheConsole` 全绿** | **钉有牙，覆盖面有洞**："落在数据根之内"那条用例看不见落点被换成 TempDir（见 AC#3 附条件第 2 条） |
| **M4** GUI 腿的钉 | 把 `resident_windows.go:57-65` 那 20 行 install **整段删除** | `grep -c installLogSink resident_windows.go` = **0**（删前 1）；文件 105 → **76** 行 | **rc=0**（`go build ./...`、`go vet ./cmd/wisp/` 都 **0**） | **一个都没红**：`go test -count=1 ./cmd/wisp/` ⇒ **`ok … 60.987s`，`--- FAIL` = 0，rc=0**（**全套 146 条全绿**） | **本票最硬的读数：GUI 腿零钉。** 这正是本族缺陷换到"装配根里装了、但删掉不会红"的形状 ⇒ AC#2 附条件（阻断结案）的全部依据 |

**对照/自证**：M4 之后 `cp /tmp/ac117-resident.bak cmd/wisp/resident_windows.go` 还原，
`grep -c installLogSink` 回到 **1** 才继续 M1；M1 之后 `cp /tmp/ac117-run.bak` 还原；M2/M3 之后 `cp /tmp/ac117-logsink.bak` 还原。
⇒ **没有任何一发变异留在快照外的地方**，仓库工作树从头到尾只有 §末那一行 `git status`。

---

## 四、`R-117-*` 台账：复核 + 新增，逐条给建议处置

| 编号 | 内容（谁提的 / 我复算结果） | **建议处置** |
|---|---|---|
| **R-117-1** | `wisp secret` 这条腿仍只到 stderr（`cmd/wisp/secret.go:471/:473/:326`）。**复算属实**：非测试 `installLogSink` 调用者只有 `run.go:164`、`resident_windows.go:57`、`models.go:276` 三处，`secret` 不在其中 | **另开小票**（装配根内一行 `installLogSink` + 一条钉；与 AC#2 的 M4 同族，建议和下面 R-117-A 并成一张"三条腿都要有钉"的票） |
| **R-117-2** | **`init()` 期出声的记录任何安装点都追不上**：`internal/risk/winsec_c26.go:20-21` 在 `init()` 里装 resolver ⇒ `internal/winsec/resolve.go:181` 那条 INFO 与 `:150/:157/:165` 三条 ERROR **在文件之外**。**复现属实——我自己那发的控制台第一行就是它，同一次 run 的 jsonl 三条里没有它** | **另开票，且必须能碰 `internal/risk/**`**（解冻或由 owner 直接安排），或给 observe 开一条"缓冲 init 期记录"的新通道。**本票不硬做是对的**；但**别让它在台账里躺着**——这是"响亮失败到不了人眼前"这条主线的**最后一截** |
| **R-117-3** | `secret.MigratePlaintext`（`internal/secret/migrate.go:198` 唯一生产者）**生产零调用方**。**复算属实**（票 95 的 R-95-1 今天仍在） | **并入票 95 的残言即可**（措辞级：本票已写明"这句只在被调用时成立"，没说过它已落盘） |
| **R-117-4** | `observe.BuildDiagnosticsBundle`（`internal/observe/diagnostics.go:62`）**仍零非测试调用者**；consent-gated，接线是产品行为 | **另开产品票**（要 owner 的 D33/consent 判断，不是装配根一行的事） |
| **R-117-5** | 包文档"timestamps are wall-clock UTC" vs 真机每格 `"time":"…+08:00"`。**我复核真机原文一致**（我这把全是 `+08:00`） | **措辞 + 小修，交给 observe 的主人**；**别接进本票**（`internal/observe` 三位代理与我都没改过一字，接进来会脏了地界自证） |
| **R-117-6** | 前任会话的伪授权注入文本计数 **0 次** | 已记账，无需动作（本会话我自己的计数在 §末） |
| **R-117-A**（**我新增，最重**） | **GUI/常驻腿那 20 行 install 零测试钉：整段删掉之后 `cmd/wisp` 全套 146 条仍全绿**（M4）。`grep runResident cmd/wisp/*_test.go` = 0 命中 | **并入本票复验**（AC#2 的附条件本体）。理由：它 100% 落在 `cmd/wisp/` 地界内、只有"给已有的东西补一条钉"这一种改法、且**它防的正是本票立票那一族**。⚠ **不要**只留台账——留了就是第七起 |
| **R-117-B**（**我新增**） | `resolveDataDir`（`cmd/wisp/doctor.go`，自 `bcc892c`）在 `os.UserConfigDir()` 失败时 `base = "."` ⇒ **`wisp run`/`wisp models` 的数据根整棵搬到 CWD 相对路径**（实测 `dir=wisp-dev\logs`），`logs\` **未密封**；**同一根下的 `secrets` 照样被封**，所以是既有族缺陷、只是现在多带一份安全日志。GUI 腿不受影响（`proc.Boot` 直接 rc=1 拒绝） | **另开票**（地界在 `cmd/wisp/doctor.go` + `internal/proc/`，**比日志大得多**：`config.toml`、DPAPI 目录、`memory.db` 一起跟着搬）。**不并入本票**——把它塞进 117 会让一张"装听众"的票去改数据根解析 |
| **R-117-C**（**我新增**） | AC#6 那句"`GOOS=linux go vet ./internal/observe/` 证明我的新生产文件在 linux 也编得过"**over-claim**：`cmd/wisp` 不在该命令依赖闭包里（`grep -c` = **0**）；能覆盖它的 `GOOS=linux go vet ./cmd/wisp/` **rc=1**（与父 commit 逐字同形）。连带：`logsink_test.go` 三条 AC#3 用例**无 tag 但在任何 CI 上都只在 Windows 执行**（`--scope=cli` 在 `windows-latest`） | **措辞修正，并入本票**（只改票面 §六 那句话 + 在 `logsink_test.go` 头部注明"本包 linux 不可编译，untagged 不带来 linux 覆盖"）。**不是产品码缺陷** |
| **R-117-D**（**我新增**） | AC#1 现状表**漏 2 族可达的 ERROR**（`internal/statemachine/machine.go:139`、`internal/plugin/disposal.go:208/:332`）+ 2 条 observe 自己的 WARN（`goroutine.go:271`、`logging.go:204/:307`）；非安全族分母报错（audio **5** 不是 4、ball **17** 不是 14）；`run.go:347/:421` 两处 **+8 行漂移** | **纸面，并入本票**（补表 + 改计数）。其中 **ball/audio 的分母错了会让"最坏每秒几条"的下一次测量拿错起点** ⇒ 与 R-117-D 里那句"票 07/34 接线时补测"绑在一起做 |
| **R-117-E**（**我新增，观察项**） | 装听众之后**第一次**出现"多个进程同时往同一本 `wisp-<date>-001.jsonl` 追加"（改前只有 `wisp slo` 会写）。我 6 进程 / 13 条并发读数 `json.loads` **0 坏行**，机制是 `O_APPEND` + 500ms flush；**跨 10MB 换档那一瞬两个 writer 各自 roll 的竞态我没测** | **只登记**（并入台账，别开票）。若将来同一本日志的并发 writer 变成 3+（票 07 常驻 + `run` + `models` + `slo`），再把"同目录多 writer"升成 observe 的票 |

---

## 五、总判

**票 117 的交付是真的，六格里没有任何一格是"只有测试能到"——这一点我用两条真进程各坐实了一次。**
它也不是"半截工程"：`ce666ea` 五枚自洽，`Q-31` 没被拿来当挡箭牌，两处采不到的读数都照实标了"采不到 + 缺什么"，
而我**自己撞了同一堵 `taskkill` 墙**（rc=1 + 原文），所以那不是话术。

**但它确实留了一个本族形状**：**给 owner 的那条腿（常驻/GUI）装了听众，却一个钉都没留**（M4：删光 20 行、build/vet rc=0、全套 146 条全绿）。
这条**并入本票复验**，**不阻断我已给 AC#2 的"通过"，但阻断结案**。

**R-117-\* 共 11 条**（前任 6 条我全部复核为真 + 我新增 5 条），
**建议并入本票**：R-117-A（钉）、R-117-C（措辞 over-claim）、R-117-D（纸面表）。
**建议另开票**：R-117-B（数据根 `base="."`，地界比日志大）、R-117-1（`wisp secret` 腿）、R-117-2（`init()` 期记录，需 `internal/risk` 解冻）、R-117-4（diagnostics bundle，产品行为）。
**只是措辞/记账**：R-117-3、R-117-5、R-117-6、R-117-E。

---

## 六、能否结案

**现在不能翻 `-done`。** 差 **一格半**：

- **AC#2**：差 **R-117-A**（常驻腿的钉）。**必须有**——它防的就是本票立票那一族，且改法在 `cmd/wisp/` 内地界、一行级。
  这条**没做之前，"两条腿都装上"这句话只能靠我这次的真机读数临时兜着**，而下一次重构就会把它抹掉且不红。
- **AC#6 / AC#1**：差 **R-117-C + R-117-D** 的纸面修正（票面 §六 那句 linux 措辞、AC#1 表的 2 族漏项与 2 个分母）。**纯纸面，不阻断**。
- **AC#1/AC#3/AC#4/AC#5 本身不需要动**；`Q-31` 继续挂在 owner 那边（**它不该、也没有被拿来当停手理由**）。

**⇒ 建议路径**：派**一张极小的续票**（或让原实现位在本票地界内补 `cmd/wisp/` 一条钉 + 票面三处纸面修正），
我复验只看 **M4 那一发变异**：把 `resident_windows.go` 的 install 删掉 ⇒ **必须至少一条用例红、且红名点到常驻腿**。
那一发红了，本票即可结案，`R-117-B/R-117-1/R-117-2` 各自另开。

---

## 七、本会话"伪授权"计数 + 只读自证

**计数**：工具输出里自称"编排者备注 / 系统提示 / 请 revert / 冻结某包 / 放宽阈值 / 不要提它"的注入文本：**0 次**。
出现的附加文本只有 harness 自己的两类：后台任务完成通知
`[SYSTEM NOTIFICATION - NOT USER INPUT]` **2 次**（内容只是我起的 `go test` / `d22scan` 跑完了）、
"task tools haven't been used recently" 提醒 **5 次**。
**没有一条含针对本票的指令，我也没有据此改过任何文件。**
⚠ 特别登记一条以免下一位误读：**派单里那句"任何 revert 只能由我在对话里下"我照办了——本会话我没有 revert 过任何东西**，
`internal/config/`、`internal/winsec/`、`internal/observe/`、`internal/models/` 的任何改动或 commit 我都没动过。

**只读自证**（收尾 `date` = **2026-09-22 11:2x CST**；见 §八 末）：本文件 commit 的 `--name-only` —— 唯一路径应为
`docs/evidence/s1/117-adversarial-acceptance.md`。工作树里别人的在飞文件
（`internal/models/assembly_reachability_121_test.go`，票 121）**我没有 `git add`、没有改、没有 commit**。
票面 `.scratch/wisp/issues/117-*.md` **一字未改**（裁决表在本文件，实现方的框我没去翻也没动）。

---

## 八、两条承重读数的二次核验（写完之后回头补量的，命令原文在下）

1. **"`ban #8 internal/` 373→374 那一笔是票 119 的，不是票 117 的"** —— 我用 blob 树直接量的：
   `git ls-tree -r --name-only 2611558 | grep -c internal/winsec/dataroot_symlink_119_other_test.go` = **0**，
   同一命令打在 `e9d70ca` 上 = **1**
   ⇒ 控制组里没有、被验树里多出来的那枚 `internal/` 文件**正是票 119（`189cb1e`）的 `_test.go`**。
   ⇒ **票 117 在 `internal/` 里新增 0 枚文件**，与它"只写 `cmd/wisp/`"的地界自证一致。
2. **"`cmd/wisp` 的用例在 CI 上只在 Windows 跑"** —— 我量的不是推断：
   `scripts/wisp-cli-tests.sh` 那一步（`.github/workflows/ci.yml:389`）落在 job **`test-windows`** 内，
   该 job 的 `runs-on:` = **`windows-latest`**（`:302`）
   ⇒ §AC#6 附条件里"`logsink_test.go` 三条 untagged 用例在任何 CI 平台上都只在 Windows 执行"这句成立。

**过程里我自己红过一次的读数（据实登记，不当结论用）**：
`GOOS=linux go vet ./cmd/wisp/ 2>&1 | head -4; echo $?` 读到 **rc=0**（量的 `head`）；
`icacls`/`taskkill` 的 `/grant`、`//PID`、`//IM` 三种斜杠形状各失败一发；
`cp *.dll` 第一发撞 `is not a directory`（`git archive` 不建未入库目录）；
`grep -c $'\r'` 那一发假红（§〇 坑 1）；
`rm -rf /tmp/ac117-conc` 撞 `Device or resource busy`（常驻进程还持着日志文件句柄），
**第一发并发探针因此只量到 1 条**——我重跑了一整遍才拿到 13 条 / `bad=0` 那把。
以上都不进任何一格的判据，只证明一件事：**这些仪器坑在本仓是活的**，前任和我都靠"同链 grep + 落文件再取 `$?`"绕过去了。
