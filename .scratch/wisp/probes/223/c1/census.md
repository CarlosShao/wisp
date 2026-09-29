# 票 223 · C1 现读普查 —— CheckAndReload 的落点在哪

- 锚点（取证时 HEAD）：起手 `git rev-parse --short HEAD` = **`6233dedc`**（11:29）；⚠ **本轮 HEAD 中途前进了**——写腿在会话内提交了 `49eb440b`→`b2e24701`→**`5a251ece`**（`git log --name-only` 现读：这三枚都动过 `cmd/wisp/run.go`，`5a251ece` 另动 `cmd/wisp/approval_always.go`）。⇒ **本件的定锚锚点是 `5a251ece`**，我在 11:44 对它重跑了全部承重尺（见下"中途重锚"），**所有被引用的 run.go / approval_always.go 行号在新锚点下一字未变**，且 `git status --short -- cmd/wisp/` 现为**空**（写腿已把这两个文件提交，不再是在飞脏件）。
- 取证时间：起手 `2026-09-29 11:29 +0800`；重锚复跑 `2026-09-29 11:44 +0800`（各节写盘前另行现跑 `date`，节内已带各自时间戳）。
- 分支：`dev`
- ⛔ 本件所有行号均为本轮 `grep`/`Read` 现跑所得；未从 `.scratch/wisp/probes/**` 任何归档快照抄过一行号（归档件在本轮 grep 里被显式排除：`grep -v '/probes/'`）。

### 中途重锚（11:44 +0800，锚点 `5a251ece`，现跑复认的承重读数）

| 承重结论 | 重跑尺 | 结果 |
|---|---|---|
| `CheckAndReload` 非测试调用者只有 balldebug | `grep -rn CheckAndReload --include='*.go' . \| grep -v _test \| grep -v /probes/` | 仍＝`cmd/balldebug/main.go:243` 一枚，其余全是注释 |
| `cmd/wisp` 里没有 reload tick | `grep -rn -E 'CheckAndReload\|NewTicker\|ConfirmLocked' --include='*.go' cmd/wisp/ \| grep -v _test` | **0 命中**（三项全零） |
| mgr / gate 构造行 | `grep -n -E 'config\.NewManager\|approval\.New\('` | 仍 `run.go:341` / `run.go:444` |
| 装配其余引用行 | `grep -n` 逐枚 | `:355 :356 :390 :447 :448 :458 :471 :473 :474 :626 :631 :633` 全部原样成立；`:433` 那句 "an L2 card can only time out into a reject" 仍在 |
| 201 的"长期允许"写口 | `grep -n -E 'AddAllowedDir\|Default\.Spawn' cmd/wisp/approval_always.go` | 仍 `:134` / `:99` |

## 用了哪些尺（可原样重跑）

只读命令清单（在仓库根执行）：
```
git rev-parse --short HEAD ; date "+%Y-%m-%d %H:%M %z" ; git status --short
grep -rn "CheckAndReload" --include='*.go' . | grep -v '_test.go'
grep -rn -E 'ConfirmLocked|OnReload|OnRestartPending' --include='*.go' . | grep -v '_test.go'
grep -rn "config.NewManager\|NewManager(" --include='*.go' . | grep -v '_test.go' | grep -v 'func NewManager'
grep -rn -i 'watchdog' --include='*.go' . | grep -v '_test.go' | grep -v '/probes/'
grep -rn 'NewTicker' --include='*.go' cmd/wisp/                 # 0 命中
grep -rn 'NewTicker' --include='*.go' internal/ | grep -v '_test.go'
grep -rn -E 'func .*RunEventLoop' --include='*.go' .
grep -rn 'SLOSampleIntervalSec' --include='*.go' . | grep -v '_test.go' | grep -v '/probes/'
grep -rn -E 'AddAllowedDir|SetAllowedDirs' --include='*.go' . | grep -v '_test.go' | grep -v '/probes/'
grep -rn -E 'NewHotkeyReloader|bridge\.Refresh' --include='*.go' cmd/ | grep -v '_test.go'
sed -n '626,660p' cmd/wisp/run.go                               # confirmModeSwitch
# d22scan 尺：
grep -n -E 'wallclockRe|unixTimeRe|timeoutWordRe|GoStmt|_test\.go|bare-goroutine' tools/d22scan/main.go
sed -n '650,777p' tools/d22scan/main.go                          # walkGo + scanGoFile
```
读过的文件（Read）：`internal/config/manager.go`（全文）、`internal/config/allowdirs.go`（全文）、`internal/config/permmode.go`（全文）、`internal/config/unwired.go`（全文）、`internal/config/schema.go:28-135`、`internal/ball/hotkey_reload.go`（全文）、`cmd/wisp/run.go:330-529,626-660`、`cmd/wisp/resident_windows.go`（全文）、`internal/proc/boot_windows.go:118-188`、`internal/agent/approval/gate.go`（签名/Options）、`tools/d22scan/main.go:1-102,650-789`。

## 结论速览

- **A**：三档（Hot/Reload/Restart）在 `manager.go` 里有完整实现，但生产 `wisp run` 里没有任何 tick 调 `CheckAndReload` → 三档今天在生产里**一条都没走过"立即"那档**（引擎在、驱动器没接）。
- **B**：`CheckAndReload() (*Report, error)`，nil Report = "无变化"；自写抑制靠 `statOwnWrite`（allowdirs.go:125 / permmode.go:78）——它只吞"自己这次写"，用户随后手动改动仍会被下一次 stat 看到，不会误吞。
- **C**：SPEC-03 §4.3 那枚"watchdog 1s tick"**在代码里不存在**为配置轮询；`cmd/wisp` 全仓无 `NewTicker`；`HotkeyReloader.Refresh` 的唯一生产调用者在 `cmd/balldebug`，`wisp run` 里零。
- **D**：`mgr`（run.go:341）与 `gate`（run.go:444）都挂在同一个 `rt` 上，二者之间**没有任何 config→gate 的引用**；`ConfirmLocked` 全仓从未被赋值 = nil = fail-closed。最小那条线只需在 `cmd/wisp`（package main 已同时 import config 与 approval）加一句赋值，不新增包依赖边——但 `ConfirmLocked` 是在持锁态里被调的（见 B/D），有阻塞/死锁代价。
- **E**：`time.Ticker` 轮询本身不碰 ban #4（那条只扫 `.Sub(time.Now())` / `time.Now().Unix(`+timeout 词）；但若用裸 `go func()` 或 `go namedCall()` 起这圈 ticker，会被 **ban #1 bare-goroutine** 判（现读 AST `*ast.GoStmt`，closure 与具名调用都算，唯一合法姿势是 `observe.Registry.Spawn`）。
- **F**：给出 3 支落点候选，每支都具名回答"生产里谁调它"。

---

## A. D36 三档生效级别的落地形态

（写盘时间 `2026-09-29 11:34 +0800`）

### A.1 先数分母

`Config` 结构体（`internal/config/schema.go:104-134`）是 `config.toml` 的完整 section 树：除 `SchemaVersion`（`:108`）外共 **18 枚顶层 section** —— `app / ball / hotkey / session / voice / audio / llm / agent / risk / fs / net / privacy / memory / panel / cost / plugins / models / observe`（逐枚 `schema.go:110-133`）。其中 4 枚是 locked section（`risk`/`fs`/`net`/`plugins`，`schema.go:118-121,131` 的行内注释 "locked section: loosening needs L2 re-confirm"）。

⚠ **重要更正（相对任务描述）**：代码里的三档是 **Hot / Reload / Restart**（`schema.go:33-44` 的 `TierHot / TierReload / TierRestart`），任务里写的第三档叫"下次会话生效"——**代码里没有"下次会话"这一枚 Tier 常量**，`TierReload` 的自陈是"applies immediately AND emits a reload event"（`schema.go:37-39`），即"立即生效＋发一个重载事件"，语义上更接近"立即"而非"下次会话"。这条命名对不上，已一并登记到末节「没查清」。

### A.2 每枚键属于哪一档（判定处＝代码，不是文档）

三档的分配**不是靠一张 tier 表**，而是硬编码在 `Manager.apply()` 及其三个子过程里。逐枚：

- **Locked 4 枚 section**（`risk/fs/net/plugins`）：走 `applyLocked`（`manager.go:158-192, 228-263`）。其内：
  - **收紧（tighten）/ 中性（neutral）→ 立即热生效**（`manager.go:250-260`，各打一条 Info 日志）；
  - **放宽（loosen）→ 必须过 `ConfirmLocked` 钩子拿 L2 复确认，nil 钩子即拒（fail-closed）**（`manager.go:238-249`，`approved := m.ConfirmLocked != nil && m.ConfirmLocked(section, loosen)`）。
  - 具体被点名的 loosen/tighten 键（方向引擎）：`risk.shell_enabled / risk.allow_shell_string / risk.l1_window_sec / risk.permission_mode / risk.shell_allowlist / risk.blacklist_overrides`（`manager.go:340-372`）；`fs.allowed_dirs / fs.reparse_point_exceptions / fs.delete_enabled`（`manager.go:376-385`）；`net.allowlist / net.block_private_ranges / net.proxy.mode`（`manager.go:389-405`）；`plugins.tier2_enabled / plugins.<id> / plugins.<id>.enabled / .capabilities / .net_allowlist / .host_api`（`manager.go:409-436`）。
- **[app] 拆分**（`applyApp`，`manager.go:267-284`）：`app.theme` = 立即（hot，`:276-279`）；`app.language / app.autostart / app.single_instance` = **重启后生效**（保留旧值，只往 `rep.Restart` 记名，`:273-283`；`schema.go:141-143` 亦注 autostart/single_instance 为 restart-tier）。
- **[voice] 拆分**（`applyVoice`，`manager.go:288-336`）：**reload 档**＝重塑音频/模型管线的键：`voice.enabled, voice.wake_word.enabled, voice.wake_word.keywords, voice.asr, voice.tts.provider, voice.tts.voice, voice.conversation_mode, voice.aec, voice.realtime, voice.cloud_asr_chain, voice.cloud_tts_chain`（`manager.go:297-306`，判定＝置 `reload=true`→`rep.Reload`+`OnReload` 事件）；**hot 档**＝调优键：`voice.wake_word.thresholds, voice.wake_word.veto_words, voice.tts.speed, voice.punctuation`（`manager.go:320-324`）。
- **其余 12 枚 section 全为 hot（立即）**：`ball, hotkey, session, audio, llm, agent, privacy, memory, panel, cost, models, observe`（`manager.go:198-221` 的 `rest` 切片逐枚列名，"Everything else is hot-tier: apply wholesale on change"，`:197`）。`schema.go` 行内注释与此一致：`BallSection` "entirely hot-tier"（`:162`）、`HotkeySection` "hot-tier"（`:176`）、`SessionSection` "hot-tier"（`:186`）。

### A.3 tier 常量本身有没有被代码用？—— 没有

`grep -rn TierHot|TierReload|TierRestart`（非测试、排除归档）只命中它们的**声明处**（`schema.go:36/39/43`），**零引用者**。也就是说 `apply()` 靠的是上面那些 `rep.Hot/Reload/Restart` 切片与方向引擎，`Tier` 枚举目前**只是文档级类型**。〔建了但没接〕（类型建了，没有消费者）

### A.4 三档里今天有没有任何一条真的走"立即"那档？

**在生产 `wisp run` 里：一条都没有。** 理由＝"立即/重载/重启"三档全部只在 `Manager.apply()` 里被计算，而 `apply()` 的唯一入口是 `CheckAndReload()`（`manager.go:117-150`），`CheckAndReload` 全仓**唯一非测试调用者是 `cmd/balldebug/main.go:243`**（调试旁支程序），`cmd/wisp/**` 生产装配路径里 grep 不到任何一次调用（`grep -rn CheckAndReload --include='*.go' cmd/wisp/` = 0 命中）。所以：**引擎（三档逻辑）在，驱动器（轮询 tick）没接**。〔建了但没接〕

- 唯一的两个"运行时真会改配置"的生产入口**不走**这套 tier：
  - `Manager.SetPermissionMode`（`permmode.go:64-82`）——由 perm.Store 在用户显式切档时调（工牌 90/101/114），直接改内存+落盘+`statOwnWrite` 抑制自读，**绕过** `CheckAndReload`/`apply`。
  - `Manager.AddAllowedDir`（`allowdirs.go:56-77`）——生产调用者 `cmd/wisp/approval_always.go:134`（工牌 201 的"长期允许"回答把目录持久化），同样直接改内存+落盘+抑制自读，**绕过** tier 引擎。
  这两条是"程序自己写的、由显式命令触发"的路径，不是"用户手改文件→热加载"的路径；票 223 问的"热加载"三档在生产里没有任何一条真的跑起来。〔已证（现读）〕

---

## B. CheckAndReload 自身的形状

（写盘时间 `2026-09-29 11:35 +0800`）

### B.1 签名与实现

- 签名：`func (m *Manager) CheckAndReload() (*Report, error)`（`internal/config/manager.go:117`）。无参、无 ctx；轮询节奏由调用者决定（注释假定"watchdog 1s tick"，见 C 节：该 tick 生产不存在）。
- 幂等语义：先 `os.Stat(m.path)`（`:118`），在持锁态比对 `mtime+size`，未变则立即 `return nil, nil`（`:125-128`）——这就是"幂等轮询"的实现，靠 mtime+size 双指纹做 no-op 判定（`:113-116` 注释自陈）。
- 变化时才重跑整条加载管线 `LoadFile(m.path, m.res)`（`:129`），失败则**仍采纳新的 mtime+size** 再返回错误（`:130-137`），注释解释动机＝"避免对一个不会自愈的文件每 tick 都重报同一个错"。代价见 B.4。
- 成功路径：`rep := m.apply(fresh)`（`:138`，**在持锁态里**）→ 更新 resolved 与 stat → 解锁（`:139-141`）→ 解锁后才 fire `OnReload`/`OnRestartPending`（`:143-148`）。

### B.2 Report 结构体逐字段（`manager.go:64-74`）

| 字段 | 类型 | 语义 | 出处 |
|---|---|---|---|
| `Hot` | `[]string` | 已被热生效（立即档）的 section 名 | `:65-66` |
| `Reload` | `[]string` | reload 档 section（已生效＋已发事件） | `:67-68` |
| `Restart` | `[]string` | restart 档 section，新值**顺延到下次进程启动** | `:69-71` |
| `Locked` | `[]LockedDecision` | 每枚发生变化的 locked section 一条裁决 | `:72-73` |

`LockedDecision`（`manager.go:49-60`）四字段：`Section`（`risk|fs|net|plugins`）、`Direction`（`DirLoosen`/`DirTighten`/`DirNeutral`，枚举在 `schema.go:51-55`；取"最强姿态"，放宽压过收紧压过中性，`:52-53`）、`Approved`（放宽时＝钩子裁决；收紧/中性恒 `true`，`:55-57`）、`Keys`（变化键路径，放宽＋收紧合并，`:58-59`；构造见 `applyLocked` 的 `keys := append(slices.Clone(loosen), tighten...)`，`:235`）。

### B.3 nil Report 的语义

`manager.go:62-64` 自陈：**"A nil Report from CheckAndReload means 'nothing changed'."** 因此 nil 有两种来路，含义不同、调用者今天无处区分：
1. `return nil, nil`＝无变化（`:125-128`）；
2. `return nil, err`＝stat 失败（`:119-122`）或加载失败（`:130-137`）。
注意第二类的"采纳 stat"分支里也是 `return nil, err`——所以一次**加载失败**会同时留下"旧配置继续生效"＋"这批新字节再也不被重看"两个后果。票 223 AC#4 要求的"不生效与读不到是两句话（缺失／语法错／权限不够／热加载被禁用各一句）"，在 `Report`/返回值这一层**只有两态（nil vs error）**，四态并不存在。〔仅文档写了、代码没有〕（AC#4 属判据要求，不是现状）

### B.4 "Manager 自己的写不能被自己当成外部改动"——自写抑制怎么实现的

两句被点名的注释：`internal/config/allowdirs.go:121`、`internal/config/permmode.go:59`（措辞均为"CheckAndReload 不得把程序自己的写读成一次未经确认的手改"）。

机制＝**写完立刻把新指纹吞掉**，两处同形四行：
- `Manager.statOwnWrite()`（`allowdirs.go:125-131`）：`os.Stat` → 把 `seenMtime/seenSize/seenValid` 刷成新值；被 `writeAllowedDirs` 在 `SaveFile` 成功后调用（`allowdirs.go:108-118`，`:116`）。stat 失败就**什么都不做**（`:127-129`），注释说这样"下一次 poll 会看见这次改动"，属保守方向。
- `SetPermissionMode` 里内联同样四行（`permmode.go:78-80`），注释明说"the same four lines SetPermissionMode runs, shared here rather than copied a third time"（`allowdirs.go:122-123`）。

**它会不会把"用户真的手动改了 allowed_dirs"一起吞掉？** 分三种情形，只有第一种安全：

1. **手改发生在程序自写之后**（正常串行）：文件 mtime/size 再次变化 → 下一次 `CheckAndReload` 的指纹比对不再成立（`manager.go:125-128`）→ 走 `LoadFile`+`applyLocked`→ 放宽会去问 `ConfirmLocked`。**不会被吞**。〔已证（现读）〕
2. **手改与程序自写落在同一个指纹窗口里**：自写抑制吞掉的是"写完那一刻的 mtime+size"，因此在此之前已经混进文件里的用户改动**一并被认成自己写的**，复确认被跳过。这一支的风险被 `SaveFile` 的全文件回写放大了：`writeAllowedDirs` 与 `SetPermissionMode` 都是 `SaveFile(m.path, m.cur)`（`allowdirs.go:111`、`permmode.go:73`），而 `schema.go:101-103` 声明写出的是**完整 struct、无 omitempty**。⇒ 程序只知道自己内存里的值，任何"本进程最后一次 load 之后"的手改（可能属于任意别的 section）都会在回写时被**静默覆盖丢失**，且覆盖动作被 `statOwnWrite` 认成自写、永不复确认。再叠加 A.4 的事实（生产里根本没人调 `CheckAndReload`，所以"最后一次 load"永远是 boot 时那一次，`run.go:341`），这一支在生产里是**默认情形而非边角**。〔已证（现读）＝代码形状；后果推断，未做动态复现〕
3. **手改保住了 mtime 与 size**（同长度改值、或 `touch -r`/复制工具保留 mtime）：mtime+size 双指纹是唯一信号（`manager.go:125`），无内容哈希、无 fsnotify（`manager.go:15-17` 注释自陈"no fsnotify"）→ 这类改动**永远不会被发现**，不是被抑制、是根本不存在检测路径。〔建了但没接——不，更准确：〔已证（现读）＝检测手段只有 mtime+size〕〕

⇒ 对票 219"长期允许"的含义：先有 `CheckAndReload` 的调用者只是必要条件之一；上面第 2、3 两支说明"接上 tick"之后仍需要（a）写前重读或内容指纹、(b) 与手改的写序约定，否则"放宽必复确认"仍可能被绕过。这两条不在本票判据里，属我量到的额外代价。

---

## C. 生产里能挂 tick 的现成载体

（写盘时间 `2026-09-29 11:37 +0800`）

### C.1 `cmd/wisp` 今天的并发/循环面（逐枚具名）

- **`cmd/wisp/` 里 `time.NewTicker` 命中数 = 0**（`grep -rn 'NewTicker' --include='*.go' cmd/wisp/`，含测试文件也是 0）。⇒ 生产 CLI 侧今天没有任何一枚自建 ticker。〔已证（现读）〕
- **协程只有两枚，且都走 `observe.Default.Spawn`（不是裸 `go`）**：
  - `cmd/wisp/approval_always.go:99` → `observe.Default.Spawn("approval-waiter", "approval", root, ...)`
  - `cmd/wisp/approval_reply.go:421` → `rt.replyHandle = observe.Default.Spawn("approval-waiter", ...)`，其体内是 `bufio.NewScanner(in)` + `for sc.Scan()`（`approval_reply.go:434-436`）＝**阻塞读 stdin 的应答监听**，不是定时轮询。
  ⇒ "approval-waiter" 是 `cmd/wisp` 唯一的既有常驻协程名。〔已证（现读）〕
- **`wisp run` 的主面是 stdin 行循环**（`cmd/wisp/panel_inbound.go:149,153` 的 `bufio.NewScanner` + `for lineNo := 0; sc.Scan(); ...`），阻塞在无输入时不会自己醒来 → 这条面**不能**充当 tick 载体。〔已证（现读）〕
- `cmd/wisp/slo_windows.go:518` 与 `:801` 各有一枚 `for {}`，但它们是 `wisp slo` 子命令里等 subject ready / 等报告的忙等循环（用 `observe.NewTimeout` 做预算，`:517`、`:798`），跑完即退，不是常驻周期任务。〔已证（现读）〕

### C.2 SPEC-03 §4.3 的"watchdog 1s tick"到底存不存在？—— **不存在**

三处代码都在引用这枚不存在的 tick：
- `internal/config/manager.go:15-17` 头注释："owns the in-memory config, polls mtime+size (no fsnotify, no extra goroutine - **the watchdog tick calls CheckAndReload**)"；
- `internal/config/manager.go:113`：`CheckAndReload is the idempotent poll (SPEC-03 sec 4.3, **watchdog 1s tick**)`；
- `internal/config/schema.go:605-606`：`SLOSampleIntervalSec ... the watchdog SLO sampling cadence (D32)`。

现读结果：
- `internal/watchdog/` 目录**只有 `doc.go` 一枚文件**（`ls -1 internal/watchdog/` → 仅 `doc.go`），全文 20 行，末段 `internal/watchdog/doc.go:18-19` 写着 `DEFERRED(watchdog loop/thresholds): implemented by ticket 42. This ticket only freezes the package boundary.`；而它的职责清单 `doc.go:3-4` 恰恰包含 "plus **config-file mtime polling reusing its loop (SPEC-03 §4.3)**"。⇒ 设计上由 watchdog 的循环顺带驱 config 轮询，**该循环今天没有任何实现文件**。〔仅文档写了、代码没有〕
- 附带现读：`grep -rn SLOSampleIntervalSec --include='*.go' .`（排除测试与归档）只命中它的**声明处** `internal/config/schema.go:605-606`，**零读者** → 那枚"5 秒采样间隔"配置键今天也没有消费者（与 C.2 的"tick 不存在"互为印证）。〔已证（现读）〕
- 全仓唯一的周期 tick 在**常驻（无参数）路径**，不是 `wisp run`：`cmd/wisp/resident_windows.go:83` → `rt.RunEventLoop()` → `internal/proc/boot_windows.go:120-138`，循环体是 `windows.WaitForSingleObject(..., 50)` 激活探测（`:126`）＋ `case <-time.After(50 * time.Millisecond)` 注释为 "loop tick: activation pump + signal check"（`:134-135`）。两个限制值得落地腿注意：(a) 它是 **50ms**、SPEC-03 说的是 1s；(b) 它跑在 `proc.Runtime` 上，而 `config.Manager` 只在 `cmd/wisp/run.go:341` 的 `wisp run` 装配里诞生（`resident_windows.go` 全文不 import `internal/config`）⇒ **这枚 tick 今天手里没有 mgr**。〔已证（现读）〕
- `cmd/wisp/run.go` 装配出的 `rt`（agentRuntime，含 `rt.mgr`）之后交给谁、有没有循环，属写腿正在改的区域（`cmd/wisp/run.go` 状态为 ` M`）；我只读到装配段 `:330-529`，未见任何 `CheckAndReload` 调用（`grep -rn CheckAndReload cmd/wisp/` = 0 命中）。

### C.3 `internal/ball/hotkey_reload.go` 的 `Refresh` 面——生产调用者现跑点数

- 字段：`HotkeyReloader.Refresh func() error`（`internal/ball/hotkey_reload.go:61`），每次 `Check()` 开头调一次（`:82-86`）；`Check()` 的三个来路是 `Run` 的 ticker 分支（`:130-138`，`time.NewTicker(every)` 在 `:130`）、`Sections()` 事件分支（`:110-117`）、以及宿主直接调。
- **非测试赋值/调用点数（现跑）＝ 2，同一枚程序里的两行**：
  - `cmd/balldebug/main.go:243` `bridge.Refresh = func() error { _, err := mgr.CheckAndReload(); return err }`
  - `cmd/balldebug/main.go:252` `bridge.Run(ctx, time.Second)`（由它驱动 `Check`→`Refresh`）
  - 构造点 `cmd/balldebug/main.go:237` `ball.NewHotkeyReloader(...)`。
- `cmd/wisp` 侧 **`NewHotkeyReloader` / `Refresh` / `bridge.Run` 命中数 = 0**；`grep -rn 'internal/ball' cmd/wisp/` 非测试只剩两句**注释**（`cmd/wisp/approval_always.go:163`、`cmd/wisp/notify_windows.go:13`，后者明说 "nothing here touches it"）⇒ `wisp run` 连 ball 包都没接（悬浮球是工牌 07 的活，`resident_windows.go:81` 亦自陈 "the floating ball arrives in ticket 07"）。〔已证（现读）〕
- ⇒ **结论**：这张"配置改了→热键立刻换绑"的桥是**建了但没接**〔建了但没接〕——它自带一份 1s ticker（`hotkey_reload.go:130`，`Run` 的注释 `:122-125` 明说这是 "the poll loop for hosts without a watchdog tick"），也就是说仓库里其实**已经存在一枚可用的 1s 轮询实现**，只是它的宿主是调试旁支 `balldebug`，不在生产装配里。这枚现成件对 F 节第 1、3 支是直接可借的形状。

---

## D. 触发 L2 复确认要碰到谁

（写盘时间 `2026-09-29 11:38 +0800`）

### D.1 两枚构造点与它们在装配顺序里的位置

- `config.Manager`：`cmd/wisp/run.go:341` `mgr, err := config.NewManager(cfgPath, nil)` → 存进 `rt.mgr`（`run.go:357`，字段声明 `run.go:227 mgr *config.Manager`）。同时刻取一份快照 `cfg := mgr.Config()`（`run.go:355`）给后面的 gate 用。
- 审批 `Gate`：`cmd/wisp/run.go:444` `rt.gate = approval.New(approval.Options{...})`（字段声明 `run.go:250 gate *approval.Gate`；`approval.New` 定义在 `internal/agent/approval/gate.go:93`，Options 形状 `gate.go:17-43`）。注意它的 `Window`/`ApprovalTimeout` 是从**上面那份 boot 快照** `cfg.Risk.*` 换算的（`run.go:447-448`）⇒ 即便 reload 接上，这两个值今天也是 boot 冻结的，不会随热加载变。〔已证（现读）〕
- 顺序上 **mgr(:341) 早于 gate(:444)**，所以任何"把 gate 交给 config"的接线只能写在 `:444` 之后。

### D.2 今天有没有任何一条从 config 侧指向 Gate 的引用路径？

- **热加载那条：没有。** `Manager.ConfirmLocked`（`internal/config/manager.go:29`，签名 `func(section string, keys []string) bool`）在非测试代码里**只有声明处 `:29` 和读取处 `:240`，全仓零赋值点**（现跑 `grep -rn -E 'ConfirmLocked' --include='*.go' . | grep -v _test.go` 的其余命中全是注释：`allowdirs.go:8,28`、`permmode.go:54`）。⇒ 生产里它恒为 `nil`，`manager.go:240` 的 `approved := m.ConfirmLocked != nil && m.ConfirmLocked(section, loosen)` 直接短路成 **deny（fail-closed）**。也就是说：D33"放宽必须触发 L2 复确认"今天不是"接了但不响"，而是**双重断开**——没人调 `CheckAndReload`（C 节）＋即使调了也没有钩子会弹卡（这一条）。〔建了但没接〕
- **但同类接线已经有一枚现成模板**，就在同一文件里：`perm.New(perm.Options{Manager: mgr, Confirm: confirm, Logf: rt.auditf})`（`run.go:473-477`），其中 `confirm` 默认取 `rt.confirmModeSwitch`（`run.go:469-472`）。`confirmModeSwitch` 的体（`sed -n '626,660p' cmd/wisp/run.go`，函数起于 `run.go:626`）做的事正是"拿着 config 侧的上下文去 gate 上签一张 L2 卡"：先 `rt.gate.AdmitTextTask(taskID)`（`run.go:631`，定义 `gate.go:160`），再 `ans, why := rt.gate.PendingApproval(ctx, tools.Decision{TaskID: "host:mode-switch", Tool: modeSwitchToolName, Level: risk.L2, Reason: ..., Mode: sw.From})`（`run.go:633-639`，定义 `gate.go:468`），答案非 `tools.AnswerAllow` 就返回错误（`run.go:641-646`）。
  ⇒ 结论：**"config 侧持有人 ↔ gate"这条边今天存在，但只服务权限档位的显式切换**（由用户命令触发，`internal/perm/store.go` 全程不碰 `ConfirmLocked`，现跑 `grep -rn -E 'ConfirmLocked|approval|Gate' internal/perm/` 非测试只有 `store.go:85-86` 两句注释指向 C18 路由）。把它复制到 reload 触发面上，形状就是 D.4 那条最小线。〔已证（现读）〕

### D.3 一个签名/运行环境的硬约束（会影响写法，先量出来）

1. **签名阻抗**：`ConfirmLocked` 是同步的 `func(section string, keys []string) bool`、**不接 ctx**（`manager.go:29`）；而 gate 的取卡 API `PendingApproval` **要 ctx** 且要先有被 `AdmitTextTask` 登记过的 taskID（`run.go:631-633`）。⇒ 桥接闭包必须自带一个 ctx（从轮询处捕获）并自造一个合成任务身份（模板即 `run.go:630` 的 `"host:mode-switch"` 常量）。
2. **`ConfirmLocked` 是在持锁态里被调的**：`CheckAndReload` 在 `manager.go:124` 加锁 → `:138` 在锁内调 `m.apply(fresh)` → `apply`→`applyLocked`→`:240` 调 `ConfirmLocked` → 直到 `:141` 才 `Unlock`。而 `manager.go:19-22` 的头注释写的是"Callbacks run OUTSIDE the lock so they may call Config()"——那句对 `OnReload`/`OnRestartPending` 成立（它们在 `:143-148`、即 `:141` 解锁之后才 fire），对 `ConfirmLocked` **不成立**。两条实际后果：
   - 卡片在最长 C18 300s 内不答（`gate.go:30-32` 注释 default 300s；值由 `run.go:448` 传入），这 300s 里 `Manager` 的 `mu` 一直被占 ⇒ 并发的 `Config()`（`:92-96`）/`Resolved()`（`:100-111`）/`AddAllowedDir`（`allowdirs.go:68`）/`SetPermissionMode`（`permmode.go:68`）全部阻塞在同一把锁上。
   - 若闭包内部再回头读 Manager（例如把当前值渲进卡面文案）→ `sync.Mutex` 不可重入，**自死锁**。
   〔已证（现读）＝锁序；后果是推演，未动态复现〕
3. **console 侧今天能不能真接到答案**：`run.go:458` 是 `rt.liveCards.bind(rt.gate, "", "", "")`（无 reply source / 无 veto 通道），`run.go:431-441` 的注释自陈"an L2 card can only time out into a reject"，工牌 201 只改了"能答"这一半、且应答面是在函数尾部由 `attachReplyListener` 挂上（`cmd/wisp/approval_reply.go:421` 的 `observe.Default.Spawn("approval-waiter", ...)`）。⇒ 落点若起在 stdin 未挂监听之前，D33 复确认会表现为"阻塞 300s 然后拒绝"，行为上等价于静默不生效——这正是票 223 AC#2 禁止拿"静默不生效"充当第二档的那个形状。

### D.4 最小的一条线建在哪个包 / 会不会新增依赖边（两条边的方向与理由，不替你选）

- **边①（宿主侧闭包）**：`cmd/wisp` 已经同时 import 两者（`run.go:47` → `internal/agent/approval`，`run.go:49` → `internal/config`），所以把 `mgr.ConfirmLocked = func(section string, keys []string) bool { ... 复用 confirmModeSwitch 的形状签卡 ... }` 写在 `run.go:444` 之后，**不新增任何包级依赖边**：新增的只是一条数据边 `main → (config 的函数字段)`，方向仍是"装配根指两个叶子"。代价：轮询 tick 也得住在 `cmd/wisp`（或由它托管），且 D.3 的持锁/300s 后果全落在这个闭包上。
- **边②（把逻辑放进 `internal/config`）**：会新增 `internal/config → internal/agent/approval` 这条包依赖边。**现读证明这条边今天编译不过**（不是风格问题）：`internal/agent/approval` 的非测试文件 import `internal/tools`（`approval/gate.go` 等 5 枚），`internal/tools/bridge.go:13` import `internal/agent`，而 `internal/agent/cost.go:6` 与 `internal/agent/loop.go:15` import `internal/config` ⇒ `config → approval → tools → agent → config` 是一个 4 节点闭环。要走走这条边必须先拆 `agent → config` 或 `tools → agent` 其中一条既有边（那就是 D22 意义上的契约/结构变更，须人工批准）。〔已证（现读）〕
- 顺带一条事实供选边参考：`internal/config` 今天的非测试依赖只有 `observe / risk / secret / winsec`（现读 `grep -rE '"github.com/CarlosShao/wisp/' internal/config/*.go`），是个近乎叶子的包；反过来 `grep -rl 'wisp/internal/config' internal/ cmd/`（非测试）给出 10 枚引用者（`internal/agent/cost.go`、`internal/agent/loop.go`、`internal/llm/{discover,probe_health,resolver}.go`、`internal/perm/store.go`、`cmd/balldebug/main.go`、`cmd/wisp/{models,panel_inbound,providers,run}.go`）——它被 agent 与 tools 链路共用，所以任何"让 config 反向认识审批"的动作波及面比看上去大。

---

## E. 写法会不会被门禁点名

（写盘时间 `2026-09-29 11:40 +0800`）

### E.1 ban #4 "墙钟时间差实现超时"实际扫的拼写

现读 `tools/d22scan/main.go`，三枚正则（`:144-146`）：

```go
wallclockRe   = regexp.MustCompile(`\.Sub\(time\.Now\(\)\)`)                    // :144
unixTimeRe    = regexp.MustCompile(`time\.Now\(\)\.(Unix|UnixNano|UnixMilli)\(`) // :145
timeoutWordRe = regexp.MustCompile(`(?i)timeout|deadline|expire|\bttl\b|budget|until`) // :146
```

判定处是 `scanGoFile` 的逐行循环（`:759-770`），对**原始源码行**跑正则，只跳过空行与以 `//` 开头的纯注释行（`:760-763`）⇒ 注释豁免、代码与字符串不豁免。两支触发：
- `wallclockRe.MatchString(code)` 单独命中即报（`:764-766`，文案 "wall-clock delta (Sub(time.Now())) in what must be monotonic logic"）；
- `unixTimeRe` **且** `timeoutWordRe` 同行命中才报（`:767-769`）。

⚠ **射程比看上去窄**（按上述正则文本现读所得，非动态复现）：`:144` 要求 `.Sub(` 的**实参**就是字面量 `time.Now()`，即只抓 `prev.Sub(time.Now())` 这一种写法。常见两种同义写法**不命中**：`time.Now().Sub(prev)`（`.Sub(` 后跟的是 `prev`）、`time.Since(prev)`（源码文本里根本没有 `.Sub(time.Now())`）。⇒ "别写墙钟差"这条禁令在仪器层面是有空隙的，落地腿就算侥幸不被判，也不等于合规；仓里既成的替代形状是 `observe.NewTimeout(budget)`（现读用例 `cmd/wisp/slo_windows.go:517`、`:798`）。

### E.2 裸 `go func(` 要配什么才算有 owner/recover —— 现读：注解不算，必须换 API

- 判定是 **AST 级**、不是正则：`scanGoFile` 里 `case *ast.GoStmt:`（`tools/d22scan/main.go:697`）。
  - 闭包字面量 → 报 "bare `go func(` is banned (D22/D38b): use observe.Registry.Spawn (named, owner, recover boundary)"（`:705-708`）；
  - **具名调用也算**（R16#1 加宽）→ 报 "bare `go f()` is banned ... named calls count too"（`:710-711`，配合 `goStmtText` 于 `:781`）。
- ⇒ 现读结论：**今天没有"加个 recover／写个 owner 注释就放过"这条路**。头注释 `:8-14` 把话写死了：唯一合法出生方式是 `observe.Registry.Spawn`，唯一的文件级豁免是**按文件路径**豁免的 `internal/observe/goroutine.go`（它自己实现 Registry），见 `tools/d22scan/allowlist.txt:7`（"Exempted BY FILE PATH only (never by call shape), so every other `go <anything>` in the tree is a finding"）。
- 仓里现成的合规样板（现读两处）：
  - `cmd/balldebug/main.go:250-254`：`observe.NewRoot(goroutineHotkeyBridge)` + `observe.Default.Spawn(goroutineHotkeyBridge, "balldebug", bridgeRoot, func(ctx){ bridge.Run(ctx, time.Second) })`，注释 `:245-249` 具名理由（R16#2：Spawn 给共享 recover 边界 D37b ＋ roster 可见的生命周期，使 main 的 join 被泄漏检测器核，而不只靠 reviewer 的眼）。
  - `cmd/wisp` 自己的两枚协程都走同一 API：`cmd/wisp/approval_always.go:99`、`cmd/wisp/approval_reply.go:421`（`observe.Default.Spawn("approval-waiter", "approval", root, ...)`）。

### E.3 bans #1-5 与 #8 各自是否包含 `_test.go` —— 两组射程确实不同

- **bans #1-5：不含 `_test.go`，也不含 `testdata/`**。头注释 `:5-6` 定调 "production scope internal/ + cmd/ (non-test, non-testdata) unless stated"；实现处 `walkGo`：目录级跳过 `testdata`/`.git`（`:660-662`），文件级 `if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") { return nil }`（`:665-667`）。作用目录由 `declaredScopes` 声明为 `"bans #1-5 internal/"` 与 `"bans #1-5 cmd/"`、kind `"production Go files"`（`:386-393`）。（ban #7 `internal/tools/` 同走 `walkText` 的 `goOnly` 分支，`:858` 同样跳 `_test.go`。）
- **ban #8（emoji）：含 `_test.go`，且含注释**。头注释 `:40-49` 明说它是**唯一**脱离"production, non-test"默认射程的一条（"This is the one ban whose scope is NOT the 'production, non-test' default above, and it is deliberately so"）；`ban8Scopes` 给的 kind 字符串是 `"Go files, comments and _test.go included"`（`:423-428`，另一处 `:574`），`walkEmoji`（`:919`）的取舍理由在 `:900-902`："comments exempt ... `_test.go` still counts: a test source is a string-literal carrier"。
- 直接后果（给判据设计用）：**同一枚 `go func` 写在 `_test.go` 里不会被判**（#1 不扫测试文件），所以"测试里绿"永远证明不了生产形状干净；反过来 emoji 写在测试里照样红。

### E.4 正面回答："落地腿用 `time.Ticker` 起一圈轮询，会不会被扫描器判违规？被哪一条？"

- **`time.Ticker` 本身不被判。** 它既不产生 `.Sub(time.Now())`，也不产生 `time.Now().Unix(`＋timeout 词，因此碰不到 ban #4 的任何一支（`:764`、`:767`）。仓里已有一枚同样形状、且位于 #1-5 射程**之内**的 1s ticker 轮询是干净的现读先例：`internal/ball/hotkey_reload.go:126-140`（`t := time.NewTicker(every)` 于 `:130`，`select { case <-ctx.Done(): case <-t.C: r.checkSafe() }` 于 `:133-138`）——它在 `internal/` 下、属非测试文件，会被扫，但没有任何 finding 形状。
- **会被判的是"这圈轮询怎么出生"：ban #1 `bare-goroutine`。** 只要落地腿写 `go pollConfigReload(ctx)` 或 `go func(){ ... }()`，`*ast.GoStmt` 分支立即命中（`tools/d22scan/main.go:697` → `:705-708` 闭包 / `:710-711` 具名）。唯一正确姿势＝`observe.Default.Spawn("config-reload", <owner>, root, func(ctx){ ...ticker loop... })`（样板：`cmd/balldebug/main.go:250-254`；recover 边界由 Registry 提供，不必自己再写一层）。
- **另外两条会顺手撞上的**：(a) 若为了"改了没生效超过 N 秒就告警"而手写 `prev.Sub(time.Now())` → ban #4 `:764` 直接命中（改用 `observe.NewTimeout`，`cmd/wisp/slo_windows.go:517` 的用法）；(b) 若在 `cmd/wisp` 里为了比对路径写 `filepath.Clean/Abs` → ban #2 `pathresolver-bypass`（`:713-721`，`AGENTS.md` §1.2 同条，唯一豁免是 `internal/risk/pathresolver.go` 之外的 `allowlist.txt:3,4,6` 三处具名文件）。本票正常不需要动路径，列出来是因为"reload 后要重新构造 C26 canonicalizer"（`cmd/wisp/run.go:390`）很容易被顺手写成那样。

---

## F. 落点候选（3 支）

（写盘时间 `2026-09-29 11:41 +0800`）

### 先决条件（三支共用，现读）

任何一支都绕不开两枚**今天全仓零赋值点**的钩子：
- `Manager.ConfirmLocked`（声明 `internal/config/manager.go:29`，读取 `:240`）——不接它，放宽永远是 `nil → deny`（D.2）。
- `Manager.OnRestartPending`（声明 `:36`，读取 `:146-147`）——非测试赋值点也是 0（现跑 `grep -rn -E 'ConfirmLocked|OnReload|OnRestartPending'` 只有 `cmd/balldebug/main.go:244` 接了 `OnReload`）。⇒ 票 223 AC#2 要的"第二档：改了→不重启→**明确告知未生效**"今天**连出口都没接**，`Report.Restart`（`manager.go:69-71`）里那份名单目前没人读。

---

### F1 在 `wisp run` 的装配里起一枚 reload 协程 + 赋 `ConfirmLocked`

- **要动的文件（file:line 级）**：
  - `cmd/wisp/run.go`：在 gate 构造与 bind 之后（`:444` 建 gate、`:458` bind）加一句 `mgr.ConfirmLocked = ...`；在装配尾部加 `observe.Default.Spawn(<roster名>, "config", root, func(ctx){ ticker 1s → mgr.CheckAndReload() })`。
  - 签卡闭包复用现成形状：`rt.gate.AdmitTextTask` + `rt.gate.PendingApproval(ctx, tools.Decision{Level: risk.L2, ...})`，即 `run.go:631-639`（`confirmModeSwitch`，函数起于 `run.go:626`）。
  - 隐藏必动项：`[fs]` 放宽生效后要真正影响路径判定，必须重造 C26 canonicalizer——它今天**只在 boot 造一次**（`run.go:390` `tools.NewPathCanonicalizer(allowed, cfg.FS.ReparsePointExceptions)`，`:383-385` 注释说明同一实例被 bridge 与 fs 工具共用），且 `rt.cfg` 是 boot 快照（`run.go:355-356`）、gate 的 `Window/ApprovalTimeout` 也从该快照换算（`run.go:447-448`）。这三处不处理的话，"接上 tick"只会让内存里的 `Config` 变了而行为不变。
- **新增导出面？要不要先落批准记录**：`ConfirmLocked` 是既有导出字段（`manager.go:29`）⇒ **不需要**新 A##。但若为上面的隐藏必动项新加导出符号（如 `Manager.OnApplied`、`rt.rebuildPaths()` 之类）就触发票 223 禁区末条"不新增导出名（要新增先落 `A##`）"。
- **协程名字不是免费的**：`internal/observe/goroutine.go:41` `ResidentBaseline = 6` 与 `:44-46` `ResidentNames` 是 D38b 冻结表，`:270` 处 `CategoryUnknown` 名会被记成**泄漏症状**（`Registry.Spawn` 于 `:262`）。现跑 `grep -rn 'Default.Spawn('` 显示 `"watchdog"` 与 `"hotkey-listener"` 两枚在册槽**今天无人 spawn** ⇒ 复用 `"watchdog"` 不动冻结表；新造名字则要改 `goroutine.go`，那是契约级变更。
- **生产里谁会调它**：`cmd/wisp/run.go` 里这枚 Spawn 出来的协程就是 `CheckAndReload` 的**第一个生产调用者**——AC#1 的读数正是"改前 `cmd/wisp` 命中 0、改后 ≥1"（现跑 `grep -rn CheckAndReload --include='*.go' cmd/wisp/` = 0）。
- **判据怎么写才在生产真观测到**：只能经接缝（`AGENTS.md` §1.3：CLI `wisp run` 是允许注入面）。起 `wisp run` → 手改 `config.toml` 的 `[fs] allowed_dirs` 加一枚目录 → **不重启** → 断言 (a) 出现一张 L2 卡（audit 行可 grep），(b) 拒 ⇒ 该目录仍被 C26 拒，(c) 允 ⇒ 之后对该目录的写通过（这一步才证明 canonicalizer 真被换过）。变异：`ConfirmLocked = nil` ⇒ 判据必须红（票 223 AC#3）；收紧方向（删目录）不得弹卡。
- **代价**：D.3 三条全落——锁内最长 300s 阻塞其他 `Manager` 调用、闭包内回头读 `Manager` 即自死锁、console 无 reply source 时表现为"阻塞后拒绝"（`run.go:431-441`＋`:458` 的 `bind(rt.gate, "", "", "")`）。外加 canonicalizer/快照失效这块未标价的工作量。⚠ 与编排者写腿同撞 `cmd/wisp/run.go`（该文件本轮 `git status` 为 ` M`），票 223 禁区末条要求与 201/222 **串行**。

### F2 挂到"设计里本该存在的那枚 tick"上：实现 `internal/watchdog` 的循环，由它驱 `CheckAndReload`

- **设计原意（三处逐字互指）**：`internal/watchdog/doc.go:3-4` "plus **config-file mtime polling reusing its loop** (SPEC-03 §4.3)"；`internal/config/manager.go:15-17` "the watchdog tick calls CheckAndReload"；`manager.go:113` "(SPEC-03 sec 4.3, watchdog 1s tick)"。
- **现状**：`internal/watchdog/` 目录**只有 `doc.go`**（现跑 `ls -1 internal/watchdog/`），`doc.go:18-19` = `DEFERRED(watchdog loop/thresholds): implemented by ticket 42. This ticket only freezes the package boundary.`；roster 里 `"watchdog"` 槽已在（`internal/observe/goroutine.go:44`）且无人 spawn；`SLOSampleIntervalSec`（`internal/config/schema.go:605-606`，注释"the watchdog SLO sampling cadence (D32)"）零读者。
- **要动的文件**：新建 `internal/watchdog/*.go`（今天不存在此文件）；`cmd/wisp/run.go` 构造并 Spawn 它、把 `mgr.CheckAndReload` 注入为其 config 面（保持 `internal/config` 不认识 approval，见 D.4 的环）。
- **新增导出面？要**：`watchdog.New/Loop` 等新导出符号 ⇒ 须先落 `A##`；并且 `DEFERRED(watchdog loop/thresholds)` 属 `SPEC-12 §5` 登记表条目，`AGENTS.md` §1.1 要求标记与登记表 **1:1 双向**对得上，把它从"推迟"改成"已落地"要动那张表＝人工批准。
- **生产里谁会调它**：watchdog 协程（由 `cmd/wisp/run.go` Spawn）调 `CheckAndReload` ⇒ AC#1 达成，且 roster 层可观测（resident 槽＋D38b 泄漏检测核它）。
- **判据**：同 F1 的 CLI 面；再加一条只有 F2 能给的：tick 存在性可由 roster 读数证明（不是靠 grep 源码）。
- **代价**：**最大且越界**。等于替票 42 把整包实现（采样、阈值、D32 表、`WatchdogAlert` 态），且要把 `SLOSampleIntervalSec` 一并接上；正撞 `AGENTS.md` §0.1"未定义即停"——这不再是"量落点"，是替别人定方案。本普查不推荐、也不在此替编排者选。

### F3 不轮询：改成显式命令入口（一次性 `CheckAndReload`）

- **要动的文件**：`cmd/wisp/main.go:82` 的 `switch args[0]`（现有 verb：`run/providers/doctor/secret/models/slo/panel-assets/panel-inbound/version/help`，`:83-119`）加一枚 verb；或在 `wisp run` 的行循环里加一条命令（`cmd/wisp/panel_inbound.go:149,153` 的 `bufio.NewScanner` + `for sc.Scan()`）。⚠ 加 verb 必须同步 `const usage`（`cmd/wisp/main.go:21-23`）：`:73-74` 具名说 usage 块由 `censusVsUsage133` 与 dispatch 做 census 核对，不同步就红。落点数据侧 `rt.mgr` 已就位（`run.go:357`）。
- **新增导出面**：不加 Go 导出符号 ⇒ 不触发 A##。**但"新增一枚用户可见 CLI 命令"要不要算契约面，我没在 `AGENTS.md`/票 223 禁区里找到明文**，已列入末节待问。
- **生产里谁会调它**：`wisp` 可执行文件自身的命令分发（`main.go:82`）⇒ AC#1 的"`cmd/wisp` 里 ≥1"达成，且**零新协程**（不碰 roster、完全不碰 ban #1）。
- **判据**：跑该 verb → 断言输出含 `Report` 四字段（`manager.go:64-74`：`Hot/Reload/Restart/Locked`）；种 `[fs]` 放宽 → 必出 L2 卡（AC#3，去掉复确认必红）；restart 档 → 必须打出"restart required"（依赖上面那枚 `OnRestartPending` 被接，或直接把 `Report.Restart` 打出来）。AC#4 的"四种各一句"在此支里最自然：`CheckAndReload` 的 `error`（`manager.go:119-122, 130-137`）与 `nil Report`（`:125-128`）在命令输出上可以直接分列——但注意 B.3 说今天代码只有两态，四态需要新写。
- **代价**：**语义不是 D33 说的"热加载"**。票 223 后果#3 抱怨的是"改了不重启不生效、而且没人告诉他"，F3 治"没人告诉他"、不治"不重启"。而且票 219 的"长期允许"**已经有**显式路径了（`cmd/wisp/approval_always.go:134` → `mgr.AddAllowedDir`，其调用者是 reply 协程 `approval_always.go:99` 的 `observe.Default.Spawn("approval-waiter", ...)`），所以 F3 的净增量只是把 reload 面也做成显式。对 D33 那条逐字要求属**打折实现**。
- 与 F1/F2 不互斥：`applyLocked` 不看谁来调（`manager.go:240`），所以 F3 同样必须接 `ConfirmLocked` 才会弹卡。

---

## 我没查清 / 查不动的

（写盘时间 `2026-09-29 11:43 +0800`）

1. **`wisp run` 装配之后的执行面我没读全。** `cmd/wisp/run.go` 全文 1173 行，我只读了 `:330-529` 与 `:626-660`。装配根把 `rt` 交给谁、之后有没有一枚循环、`attachReplyListener`（`cmd/wisp/approval_reply.go:421`）在 run 路径上究竟何时挂——都没读。⇒ **C.2 的"生产无配置 tick"只覆盖 `cmd/wisp` 里 `NewTicker`/`CheckAndReload`/`ConfirmLocked` 三项 grep 为零这一事实**（这三把尺我在 11:44 对定锚 `5a251ece` 又复跑过一次，仍是 0），不等于我逐行看过 run.go 全部 1173 行。
   ⚠ 本轮内 HEAD 从 `6233dedc` 前进到 `5a251ece`（写腿提交了 `49eb440b`/`5a251ece`，两枚都改过 `cmd/wisp/run.go`，后者还改了 `cmd/wisp/approval_always.go`）。我已在文件头的"中途重锚"表里逐枚复认被引行号**在新锚点未漂移**，`git status --short -- cmd/wisp/` 现读为空（不再是脏件）。但**若后续再有提交落进 `cmd/wisp/`，本件的 run.go 行号需重跑一遍才能引用**——尺已写在文件头，可直接原样重跑。
2. **三档命名对不上，我没去文档里定谁对。** 任务描述写"立即／重启后／下次会话"，代码只有 `TierHot/TierReload/TierRestart`（`internal/config/schema.go:33-44`），`TierReload` 自陈是"立即生效＋发重载事件"。我没有去 `docs/PLAN.md` 的 D36 正文核对原始三档措辞（只核了票 223 引用的 `PLAN.md:1644-1645` 那句 D33），所以**"下次会话生效"到底是 D36 原词、还是任务描述转述时的偏差，我没查清**。这影响 A 节的分类口径，需要人拍一下。
3. **所有"阻塞/死锁/覆盖丢失"都是读码推演，没有一次动态复现。** 硬约束禁止任何编译与运行，所以我没有跑过 `wisp run`、没有种过手改、没有观测过 300s 阻塞，也没有实测 B.4 第 2 支的"手改被程序回写覆盖"。三处后果我都已在正文标了"未动态复现"。**B.4 第 2 支是本件最值得怀疑的一条**：它成立的前提是 `SaveFile` 写整个文件且写前不重读磁盘——前者我只由 `internal/config/schema.go:101-103` 的"every field is emitted by Write (no omitempty)"＋调用形 `SaveFile(m.path, m.cur)`（`allowdirs.go:111`、`permmode.go:73`）推得，**我没有打开 `loader.go` 去核 `SaveFile`/`atomicWrite` 的实现**（只看到 `permmode.go:50` 提到 atomicWrite）。若 `SaveFile` 内部其实有"读盘合并"或写前 re-stat，那一支结论要降级。后续腿请优先复核这里。
4. **`DEFERRED(watchdog loop/thresholds)` 的账我没能两向核。** 我现读了标记本身（`internal/watchdog/doc.go:18-19`）与"该目录只有 doc.go"，但**没有**去 `SPEC-12 §5` 登记表逐条核它是否在册、也没有核票 42 的闭环状态。`AGENTS.md` §1.1 要求标记与登记表 1:1 双向——这一条我只做到单向。若 F2 被选中，这是必须先补的读数。
5. **`ResidentBaseline = 6` 与实际 spawn 数的矛盾我没查清。** 现跑 `grep -rn 'Default.Spawn('`（非测试）里没有 `"watchdog"`、也没有 `"hotkey-listener"`，而 `internal/observe/goroutine.go:41` 把 resident 预算冻在 6、`:44-46` 列了 6 个名字。⇒ 今天到底是"只有 4 枚常驻真在跑、另外两枚待接"，还是有我 grep 漏掉的 spawn 形状（例如 `reg.Spawn`/`q.reg.Spawn`/`opts.Registry.Spawn` 这几枚我看到的都是别的名字），我没有把 6 枚逐一对上号。F1 的"复用 `"watchdog"` 槽"这步棋依赖这里，请先复核。
6. **CLI verb 是否算"新增面"没有明文依据。** F3 说"不加导出符号就不触发 A##"，但"新增一枚用户可见命令"要不要落批准记录，我在 `AGENTS.md` §1.1 与票 223 禁区里都没找到针对 CLI 表面的条款（只找到 `cmd/wisp/main.go:73-74` 那枚 usage↔dispatch census 机制）。这一条留给编排者定，我没有替你判断。
7. **`internal/panel` 的既有红名与三枚冻结件我没碰。** 票 223 AC#5 说历史在册的 `internal/panel` 红不算本票新增，禁区又点名三枚冻结测试文件一字不动。我既没跑测试（禁止），也没读那三枚文件，所以**没有**任何关于"reload 接上后 `internal/panel` 会怎样"的读数。
8. **`frontend/**` 与 `design/**` 我一行没读、也没转述**（两层都禁）。唯一相关的一句是扫描器自己的射程描述（`tools/d22scan/main.go:24-37`、`:395-407` 讲 ban #6 扫 `frontend/` 的文件数），那是**门禁代码里对扫描范围的自述**，不是我读取前端内容的结果——我不据此对前端下任何结论。
9. **归档件我只当作"尺"，没当作"读数"。** 本轮 grep 里凡是可能命中 `.scratch/wisp/probes/**` 的尺都加了 `grep -v '/probes/'`；`internal/watchdog` 只有 doc.go 这一条是我的 `ls` 现读，不是引自 `docs/evidence/s1/121-*` 或 `164-*`（我确实在一次未过滤的 grep 里瞥见过那两份文件提到同一事实，但没有采用它们的行号或结论）。
