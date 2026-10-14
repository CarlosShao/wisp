# 174-r9 —— 票 174 AC#2b 现读复核＋两向三态（结论：接线早已在盘上；零产码改动）

- 腿：174-r9（实现者；⛔ 框未翻、⛔ `Status:`/票名未动、只 commit 未 push）。面＝`cmd/wisp/`（含其测试文件）＋`probes/174/r9/**`＋票 174 尾 1 bullet。
- 起手锚（21:2x 首读）＝`dd2e4fe3`；**漂移逐报**：读数期间 HEAD 走到 `ccdceff5`（`9905bb50`＋`ccdceff5` 两枚编排者 bookkeeping：只动 `issues/174-*.md`／`docs/evidence/s1/model-page-v3-brief.md`／`docs/reports/pending-and-issues.md`，码面零命中）。
- 提交：代提 `14b42887`（独立一枚，仅 r8 遗产件）｜交付＝（本枚：probes/174/r9/**＋票尾 bullet）。
- 窗口自证：21:22:04／21:24:13–34／21:26:25–49／21:28:07–28／21:29:18／21:29:57 皆安静（tasklist 锚定形）；⚠ 21:38:23 一发撞见别人在飞笔 ⇒ 那一发不作主证，已错峰重跑（§5）。

## 0. 结论先行

1. **AC#2b 的接线不是"待做"，而是 2026-09-28 18:15 就落了盘**：提交 `4db5f3f6`（腿 `145-r3`，提交标题逐字含"`145-r3(格②/票174 AC#2b)`"）。尺：`git log -S "Paths: rt.paths" -- cmd/wisp/run.go` **单命中**。
2. **票面 `:362` 与现 `:547` 是同一跳**：同一语句 `tools.BuiltinTaskEntries(tools.TaskDeps{…})`，行号只被后来的注释/插入推后 185 行（`4db5f3f6^` 那一行现读在 `:365`；它的 diff hunk 头写 `@@ -362,7 +362,7 @@`，票面 `:362` 即由此而来）。
3. **判据两向＋安静正例三态全齐，且是端到端的**：走 `assembleRuntime` 交出的真桥、真 `TaskDeps`、真 C26 实例、真落盘副本（不是 174-c1 的接缝级测法）。
4. ⇒ 本腿**零产码改动**（照编排者 20:5x 预警的①支执行：⛔ 不重复造接线、⛔ 不改已接好的行）。"改动 diff"＝**零**；本腿唯一的 diff 是 overlay 里的**反形**（§3），用于读"接前"那一向。
5. ⚠ 一枚新情报：r8 **并非"什么都没留下"**——盘上有一枚**未跟踪**测试件 `cmd/wisp/task_output_pointer_notice_174r8_test.go`（7531 B，21:09:36），它就是 §3 三态的 harness，也是 AC#2b 目前**唯一**的端到端判据。按编排者 21:3x 裁定**代提**进 git（§4，**原件一字节未改、未改名**）。

## 1. 接线现读复核（逐条可重跑）

- 那一行（`cmd/wisp/run.go:547`，现读）：`for _, e := range tools.BuiltinTaskEntries(tools.TaskDeps{Roster: rt.tasks, Paths: rt.paths}) {`
- **对象同一性**（AC#2b 的"路径判定者"就是 C26 本尊）：
  - `rt.paths` 的类型＝`*tools.PathCanonicalizer`（`cmd/wisp/run.go:275`）；
  - 它由 `cfg.FS.AllowedDirs` 在 `cmd/wisp/run.go:500-504` 建成（`allowed := …cfg.FS.AllowedDirs…`；`tools.NewPathCanonicalizer(allowed, cfg.FS.ReparsePointExceptions)`）；
  - **同一实例**另投两处：`FSDeps{Paths: rt.paths}`（`:519`，`fs.read` 用它）与 `tools.New(tools.Options{Paths: rt.paths})`（`:750`，桥用它）；
  - 消费侧：`internal/tools/task.go:49` `Paths *PathCanonicalizer`，只在 `:836 d.Paths.Canonicalize(raw)` 与 `:847 d.Paths.InAllowlist(canon)` 两处被问。
- **AC#2b 的 ⛔ 未被违反**（"不许让 task.go 自己去 filepath 拼路径或读配置绕过 C26"）：`grep -n "filepath\.\|os.Getenv\|config\." internal/tools/task.go` 唯一命中在 `:47` 的**注释**里（"…filepath.Clean/Abs of one's own…"那句自我约束），零实际调用、零配置读。
- `internal/tools/task.go`（判定者语义）与 `cmd/wisp/run.go` 在本腿**零字节改动**：`git status --porcelain -- internal/ cmd/` 在代提前只剩那枚未跟踪件，代提后为空。

## 2. 三态读数（逐字在 `logs/01`、`logs/02`；原始在 `logs/raw/`）

| 态 | 怎么造 | 回执里那一段（模型可见文本） | 判决 |
|---|---|---|---|
| **接前**（未接线） | overlay 把 `run.go:547` 那行**字节级**种回接线前形态 | `；注意：路径授权判定者未接线（fail-closed：C26 没接进来，这条路径是否还读得回来无法核实，按读不到处理）` | **T1、T2 双红** rc=1 |
| **接后**（现树／已提交树） | 不动 | `；注意：这条路径现在读不到，它不在你被授权的目录范围内，fs.read 会被拒；要用户先把所属目录加进 [fs] allowed_dirs 才读得回来` ＋ `全文见 <artifact>` | **双绿** rc=0 |
| **安静正例** | 根内＋真实普通文件（`allowed_dirs=[数据根]`） | 同段里 `注意：` **0 次**、`读不到` **0 次**，指针仍在 | **绿** rc=0 |
| 还原 | 撤 overlay | 复绿；`run.go` md5 与 HEAD blob **同值** `07541a1f…` | rc=0 |

- 三态共用同一发形状：`task.output` 读一枚 48000 字节（约 12000 token）的登记行 ⇒ 走 D15(3) 截断＋指针那一支；两态的唯一差别只有"授权根里有没有这枚副本"。
- 精确文案与票面 AC#2b 引的逐字对得上（"不在你被授权的目录范围内"＋"加进 `[fs] allowed_dirs`"两条都在）。
- 判据件（§4 那枚）的两向牙在同一次 overlay 里同时给出：**接前**那一发不但"精确文案"缺，**连安静正例都被污染**（`注意：` 1 次）⇒ 接线这一行的因果是可被两枚用例从两个方向咬住的，不是单向恒真。

## 3. 反形（本腿唯一的 diff）与字节同一性

- 造法：`sed` 把 `tools.TaskDeps{Roster: rt.tasks, Paths: rt.paths}` 换成 `tools.TaskDeps{Roster: rt.tasks}`，落到 `mut/run-unwired.go`；`diff cmd/wisp/run.go mut/run-unwired.go` ＝ **恰 1 行**（`547c547`，全文＝`raw/06-mutant-diff.txt`）。
- **字节同一性证明**（"字节级种回"不是形容词）：
  - `git show 4db5f3f6^:cmd/wisp/run.go | grep 'BuiltinTaskEntries(tools.TaskDeps' | sha256sum` → `e3b19cba121f41de8a029ec974201cfeba2d32bd5a91837fadecb52e13e2fcbf`
  - `grep 'BuiltinTaskEntries(tools.TaskDeps' mut/run-unwired.go | sha256sum` → **同值**。
- overlay 生效是**行为级**证明的：若不生效，回执里不可能出现"未接线"那句（`d.Paths == nil` 才走那支，`internal/tools/task.go:833`）。

## 4. 代提（与 §2/§3 的读数**分开**记：谁跑的门、跑了哪几门）

- **裁决**：编排者 21:3x 一句话＝"收编进 git（甲），按死腿遗产代提办"。
- **这一笔是谁做的**：`174-r9` 代 `174-r8`。原件的产码者仍是 r8（文件内注释自述"写腿 174-r8"），本腿**只做了"跑门＋提交"两件事**，未改内容一字节、未改名。
- **本腿为它跑的门（只此几道）**：`gofumpt -l <该件>` **空** rc=0（⛔ 未触发"停手上报"）／`go build ./...` rc=0／`scripts/d22scan.sh` **clean** rc=0／`check-path-length-budget.sh` **VERDICT GREEN** rc=0／`go test ./cmd/wisp/ -run 'TestTaskOutputPointerNotice174r8' -count=1` `ok 12.305s` rc=0（这一发编译到本件）。
- **本笔未跑（具名）**：`-race`、整包 `cmd/wisp` 四数、端到端 `wisp run`。
- 提交：`14b42887`（独立一枚，显式 pathspec `cmd/wisp/task_output_pointer_notice_174r8_test.go`，1 file / +184）。
- 冻结副本 `mut/harness-174r8-frozen-copy.go.txt`（md5 与原文件同值 `4bf7eadceb7435c41513b53815cf3e9a`）留盘**不删**；`probes/174/r8/**` 与其它临时件本腿只建不删、未碰。

## 5. 窗口与协调

- 主证窗口全安静（§0 列出的六枚时刻，tasklist 锚定形自证；`logs/03` R3）。
- ⚠ **自报一处仪器缺陷并更正**：本腿中途的窗口尺曾写成 `grep -cE 'go\.exe|wisp\.test\.exe'`，而 `cargo.exe` 含子串 `go.exe` ⇒ 假报"2 枚在飞"。改为行首/空白锚定后归零（`logs/03` R3）。
- ⚠ **21:38:23 一发撞见别人在飞笔**（`go.exe 36364`＋`wisp.test.exe 38860`，非本腿所起；21:39:32／21:41:16／21:43:29 复查仍在，约 7 分钟，形状像 314 死区那种整趟长跑）⇒ 派单"重包同一时刻只跑一个"触发：**那一发不作主证**（只作旁证）。**已错峰重跑**：21:45:30 现读归零 → 21:45:45–21:46:03 安静窗复跑**已提交树**，两枚 `--- PASS`、`ok 11.434s` rc=0（`logs/05`）。主证仍是 §2 的三个安静窗（21:24:13／21:26:25／21:28:07）。
- 面归属：`internal/tools/**` 归 `174-r7`（本腿零字节）；`cmd/wisp` 面本波另有 `133-v` 在跑重包（上面撞见的在飞笔疑似它，未读、未归因）。

## 6. 门禁现量（本腿面）

- `go build ./...` rc=0；`gofumpt -l cmd/wisp/run.go mut/run-unwired.go` 空 rc=0。
- `go test ./internal/tools/ -run '174|PointerNotice|Pointer' -count=1` → `ok 3.276s` rc=0（21:29:57）。
- `scripts/d22scan.sh` → **clean** rc=0（`internal/=557`、`cmd/=129`）。
- `scripts/check-path-length-budget.sh` → **VERDICT GREEN** rc=0；代提后复量 `tracked paths=13000 over-budget=60 covered by roster=60 not in roster=0`（+1 恰是代提那一件，仍在预算内）。

## 7. 具名没做（本腿射程外／未跑）

- 未改任何产码（`cmd/wisp/run.go`、`internal/tools/**`、`internal/risk/**`、`[fs]` 默认值、`docs/PLAN.md`、`docs/specs/**`、`frontend/**`、`design/**`、`thresholds.go`、golden、`allowlist.txt` 全零字节）。
- 未跑：`-race`、整包 `cmd/wisp` 四数（314 死区，派单明禁当主证）、端到端 `wisp run` 真进程、`scripts/wisp-cli-tests.sh` 全量、`portable-tests.sh --scope=cli`、`gate-clauses.sh`。
- ⛔ 未翻 AC 框、⛔ 未动 `Status:`/票名、⛔ 未 push。
- 未替 AC#2b 出**判语**（那属非实现者；本腿只交读数与复核）。

## 8. 复现原料（留盘未提交，规则 8 只建不删）

- `logs/raw/` 9 枚：`01-now-r8harness`（首跑＝缺 DLL 的 `0xc0000135` 复现，见 174-c1 旧账）／`01-now-r8harness-wired`／`02-unwired-overlay`／`03-restore-green`／`04-gates-a`／`05-gates-b`／`06-mutant-diff`／`07-adopt-gates`／`08-post-adopt-gates`／`09-offset-rerun`（＋`09-window-and-rerun`）。
- `mut/` 4 枚：`run-unwired.go`＋`overlay-unwired.json`＋`harness-174r8-frozen-copy.go.txt`＋`commit-msg-adopt.txt`。
- 提交面的 `logs/01..05` 为上述 raw 的紧凑摘录（决定性逐字＋rc）。
- 复跑一条命令（安静窗内）：`export PATH=<root>/third_party/sherpa-onnx:$PATH; GOTMPDIR='C:\tmp\wisp174r9gotmp' go test ./cmd/wisp/ -run 'TestTaskOutputPointerNotice174r8' -count=1 -v`；反形加 `-overlay .scratch/wisp/probes/174/r9/mut/overlay-unwired.json`。
