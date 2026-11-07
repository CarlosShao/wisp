# 138 r2 裁定书 — internal/tools 面 AC#1 审计：取消级联在本面"无入口踩得到"，AC#2 零改动裁定

- 工单：`.scratch/wisp/issues/138-cancelling-a-task-kills-nothing-the-task-itself-spawned-single-jobscope-no-terminatejobobject.md`
- 本件：wave75 写腿（票 138，写面 `internal/tools/**`）的审计与裁定落盘；被验基线 `HEAD = 938918176`（A1618 派单枚）。
- 时间：2026-11-07 20:33–20:4x +08（每条读数现跑，不背旧数）。
- LF：本件 LF；`git add`/`commit` 两处显式 pathspec；只 commit 不 push。

---

## §0 起手三核（现跑读数，2026-11-07 20:3x）

### 0.1 ① 核票面点名的串/行还在不在

| 票面锚点 | 现跑命令 | 现读数 | 判定 |
|---|---|---|---|
| 全仓 0 处 `TerminateJobObject` | `grep -rn "TerminateJobObject" internal/` | 生产 1 处：`internal/plugin/manifestslot/command_job_windows.go:146`（票 52 落地） | 票面锚失效（r1 裁定书已记，本腿复核同读数） |
| `internal/tools/unwired.go` "no shell.exec tool is registered" 注释 | `grep -rln unwired internal/tools/` | **该文件已不在盘** | 票面此锚整段过期（09-24 世界） |
| 票 138 修法族前提"唯一的 JobScope 是整机一枚" | `grep -rn "type JobScope" internal/` | 仍唯一：`internal/proc/jobscope_windows.go:62`；但 manifestslot 另有每子进程一枚 `childJob`（r1 裁定书 §1.3） | 与 r1 裁定书一致 |

### 0.2 ② 核引文出处（r1 裁定书收件核——收件三把尺）

- **commit 真在**：`git cat-file -t` 逐枚过：`8429d0e60`（步1）／`b3e990300`（步2）／`bf224be19`（步3）／`b576c10a3`（步4）／`9287d1179`（步5）——五枚全在盘。
- **件真在**：`docs/evidence/s1/138-cascade-ruling.md` ＋ `internal/plugin/manifestslot/command_138_cancel_test.go` 均 tracked；票面 Progress log 16:3x 行已在库（工作树 tracked 零脏）。
- **自己复跑**：`go test -count=1 -run Test138TaskCancelKillsTheProcessTree ./internal/plugin/manifestslot` ⇒ **PASS 1.97s，rc=0**（2026-11-07 20:40:58 现跑）。

### 0.3 ③ 红/绿是不是别人的在飞件造成的

- `git status --porcelain` 在本腿三处写面（`internal/tools/`、`docs/evidence/s1/`、票 138 文件）⇒ **空**；wave75 另两腿（145-snapshot=`internal/panel`、68-parity=`internal/ball`）的面本腿零触碰。
- 未跟踪件里与本票写面零交集（`cmd/wisp/*323*test.go`、`docs/evidence/s1/{195,351,378}*` 等），一枚不碰、不 commit。

---

## §1 步1：internal/tools 面的进程 spawn 审计（AC#1，现跑读数）

### 1.1 本包全部进程创建点（`grep -rn "exec.Command|os.StartProcess|SysProcAttr|CreateProcess|ShellExecute" --include=*.go internal/tools/ | grep -v _test.go`）

| 创建点 | 形状 | 生产可达？ |
|---|---|---|
| `launcher_windows.go:113` | app.launch path 臂 → `w.deps.Jobs.StartInJob(cmd)`（C30 Job） | **不可达**（见 §1.2） |
| `launcher_windows.go:116` | app.launch AUMID 臂 → explorer 载体进 Job（UWP 本体由系统激活宿主承载，件头如实登记） | 不可达 |
| `launcher_windows.go:54` | web.open/file.open → `ShellExecuteW "open"`（shell 移交，fire-and-forget） | 不可达 |
| `launcher_other.go:27/:30/:71` | 非 Windows 半边（xdg-open/open/StartInJob） | 不可达（D7 只 ship Windows） |

### 1.2 可达性裁定：web 家族今天未接进生产装配

- 装配根现跑（`cmd/wisp/run.go`，只读）：注册的是 `BuiltinFSEntries`(`:651`)／`BuiltinTaskEntries`(`:689`)／`BuiltinSystemEntries`(`:720`)／`BuiltinSubagentEntries`(`:1124`)／`BuiltinPlanEntries`(`:1163`)——**五族，无 `BuiltinWebEntries`**。
- 全树精确 grep：`grep -rn "BuiltinWebEntries" --include=*.go cmd/ internal/ | grep -v _test.go` ⇒ **零命中**（rc=1）。
- ⇒ 本包唯一的子进程 spawn 面（launcher，含 app.launch）今天在生产里**没有入口**。票 138 的缺陷（"取消杀不掉自己 spawn 的子进程"）在 internal/tools 面**无入口踩得到**。

### 1.3 两个"看起来像缺口"的点的定性（按票面 AC#1③ 的两问逐条答）

1. **app.launch 子进程在取消后继续活着**——不是缺口，是交付物语义：用户要的就是"应用被打开"。取消任务杀掉它＝毁掉任务的交付。工作型子进程（agent 代表任务跑的命令）的级联杀归属 **D46 命令插件**（manifestslot per-call Job + `TerminateJobObject`），r1 裁定书 §1.2 八环＋取消半边牙（`Test138TaskCancelKillsTheProcessTree`，本腿复跑 PASS）已闭环。C30 的 `KILL_ON_JOB_CLOSE` 是坠机网不是取消机制，其牙已在盘：`Test22AppLaunchNotepadUnderRealJob`（`app_launch_windows_test.go:90-107`，关 Job ⇒ pid 死）。
2. **子代理父取消不级联**——票 197 §0 的冻结设计，不是本票缺陷：`subagent_197.go:353-374`（`context.WithoutCancel`）＋诚实答复（`:430-440` "没有被级联取消…task.cancel 单独停"）＋载荷钉两处（`subagent_197_test.go:729` AC#5、`:806` Description "停掉父任务不会级联"）。子代理的命令子进程仍走 manifestslot per-call Job，级联杀对它们成立。

### 1.4 AC#1③ 两问的答案

- 是"已实现功能里的缺口"吗？**不是**——web 家族未接线，无功能在跑，无现实危害。
- 是"给未来功能预留的洞"吗？**不是**——接线归属已有 owner（`web.go` 件头注明 wiring ticket＝票 12，per R12），且接线后形状 by construction 就是每子进程 `StartInJob` 进 C30 Job（`startInJobOrBare`，`launcher_windows.go:129-141`；Jobs 未接线时显式打弱姿态告警，不是静默）。
- ⇒ **无需去 `SPEC-12 §5` 新登记，无需降级本票**。r1 裁定书对修法的归属（manifestslot）与本腿对 internal/tools 面的"无入口"裁定合起来，票 138 的"取消级联"在今天的仓里**已闭环**。

---

## §2 步2：AC#2 裁定——internal/tools 面零代码改动（先裁再写）

- **不落"取消任务杀掉 app.launch 子进程"**：理由 §1.3（交付物语义；主流 harness 的分法＝工作进程级联杀、交付物存续——Claude Code 类 harness 对"打开应用/网页"类工具同样不在 cancel 时回收目标应用）。
- **不叠 per-task Job Object**：r1 裁定书 §2.1-§2.2 已裁（票 52 后 per-call Job 占位，层级复杂买不到真价值），本腿不改判。
- **不动子代理取消语义**：票 197 §0 冻结＋双钉在盘；本腿改它＝改别人的冻结契约。
- **不新增牙**：本面无可钉的"取消杀子进程"实现（无入口）；已有的三颗牙覆盖了全部真实形状——manifestslot 取消半边（r1 落）、C30 kill-on-close（票 22 落）、子代理不级联（票 197 落）。突变牙的前提是"有 kill 那一跳可摘"，本面那一跳不存在于可达路径，硬造一枚=为结案造假危害，正是票面 AC#1 警句禁止的形状。

---

## §3 步4：门禁读数（2026-11-07 20:43–20:4x +08，全部现跑）

| 门禁 | 命令（要点） | 读数 |
|---|---|---|
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l internal/tools/`（绝对路径） | **空表**，rc=0 |
| go vet | `go vet ./internal/tools` | **rc=0** |
| build | `go build ./internal/tools` | **rc=0** |
| test | `go test ./internal/tools -count=1 -timeout 30m`（跑前 `unset GOTMPDIR TMPDIR TEMP TMP`＋`export PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`） | **ok 36.674s，rc=0** |
| d22scan | `sh scripts/d22scan.sh` | **clean rc=0**；448 生产 Go 文件（internal 385＋cmd 63）、ban#6 frontend=116、ban#7 tools=48、ban#10 risk=21、ban#8 internal=956/cmd=236/design=39——各 scope 无降 |

---

## §4 交付清单与边界

- 本腿 commit（一枚，只 commit 不 push）：证据件＋票面 Progress log 一行。
- 写面：`docs/evidence/s1/138-cancel-kills-cascade.md`（新增）＋票 138 Progress log 一行。**零 .go 改动**。
- 未做／边界（如实）：
  - internal/tools 零代码改动——AC#1 判"无入口踩得到"，AC#2 判本面零改动（§2 三条理由）；派单步2/步3 的"落地＋新牙"在本面**没有可落对象**，这是裁定的结果不是省略；
  - AC 勾与 Status 零触碰（裁决者≠实现者；本腿连 Progress log 也只追加一行，未翻任何勾）；
  - `internal/plugin/manifestslot/**`、`internal/proc/**`、`internal/agent/**`、`cmd/wisp/**`、`internal/panel/**`、`internal/ball/**`、`internal/config/**`、`frontend/**`、`design/**`、thresholds/golden：零触碰（manifestslot 的牙只跑读数，未改一字）；
  - web 家族接线（票 12/R12）未做、不归本票。
