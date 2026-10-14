# 253-r6 — 票 253 AC#3 落地：旧封套被静默丢弃那一支要有声（写腿；面＝`cmd/wisp`）

**身份**：实现者。判语归后续非实现者腿；⛔ 本腿未翻任何 AC 框（AC#3 仍空）、`Status:`／票名未动。
**窗口**：起手锚 `581dd252`（2026-10-14 14:0x 现读）。执行窗内 HEAD 逐枚前移（他腿笔：`c2b00de4`／`db5ec603`／`4c7f0726` …）；每笔都 `git log --name-only` 核过，**零 `cmd/wisp`／`internal/panel` 重叠**（唯一同窗 Go 改动＝`internal/tools/taskid_source_236r4_teeth_test.go`，他票）。
**判据回顾（票面 `:18` 逐字）**：⛔ 不许改上游库；在 Go 侧给「收到但形状不认识」落一条审计/告警（或至少让面板宿主那条腿具名拒收并回一个错误）。判据＝一发假腿按旧封套发一次 ⇒ ①Go 侧有可读的现场；②页面拿到的不是无声 `null`。

## 0. 交件（产码／测试／证据分开）

| 件 | 内容 | 凭据 |
|---|---|---|
| 产码 | `cmd/wisp/panel_control_windows.go` | diff **+54/−1**；sha256 `ae714aee…`（HEAD blob＝改前）→ `2fe49342…`（终态）；`boundDoor` `:704`、`msgcb` 两支 `:711`/`:717`/`:728`/`:734` |
| 测试 | `cmd/wisp/panel_unbound_253r6_windows_test.go`（新增 270 行） | sha256 `53ed0947…`；单窗用例 `Test253r6TheSilentDropBranchIsLoud` `:199` |
| 证据 | 本文件＋`logs/10`…`14`（§7 清单） | raw 大日志**留盘未提交**（§7 表） |
| 票尾 | 253 尾追加 1 条 bullet（`agent=253-r6`） | 框未翻 |

## 1. 机制现读：旧封套今天走哪条路被吞（逐点，可达／不可达）

自 ticket 299-r1 起，RPC 管道（`msgcb`/`callbinding`，库里 `webview.go:139-220` 的移植）**已是本仓代码**（`panel_control_windows.go:669` 起：`msgcb` `:711`、`callbinding` `:757`、`Bind` `:814`；`browser.MessageCallback = c.msgcb` 在 `:339`）。所以票面当年那句「属第三方模块、本仓不修上游」在**生产路径上已经过期**——今天要修的静默支就在本仓。

- **P1 页面 `window.chrome.webview.postMessage(旧封套)`（票面字面形）**：ticket 35/303 的转发钩子 `panelPostMessageForwardInit`（`panel_host_windows.go:813`）在页面自己的 postMessage 上把它折成 RPC 帧 `{id,method:"wispDispatch",params:[旧封套]}`，到 `msgcb` → 绑定闭包（`panel_host_windows.go:842`）→ `dispatchRaw` → `composer.ComposerDispatch.Handle`。**可达，且已有声**：本腿真窗现跑逐字＝页面 `RESOLVED('面板请求被拒绝 [panel.approval.request (无 requestId)]：panel: composer request refused: 方法 ''panel.approval.request'' 不是面板 composer 通路的能力入口')`＋Go 侧 `panel: INBOUND-DISPATCH … method="panel.approval.request" …`（`logs/11` 第 5/6 行）。⇒ 这条是 35/303 闭掉的；本腿改动**没碰它**（正控 B）。
- **P2 旧门 `window.external.invoke(旧封套)`（35 之前页面用的门；钩子按设计把 `inside=true` 的调用逐字路由到原生出口）**：消息原样到 `msgcb`；`json.Unmarshal` 把 `{"method":"panel.mode.request",…}` 解成 `{id:0, method:"panel.mode.request", params:nil}` → `callbinding` 找不到绑定门。**改前**：`return nil, nil`（`webview.go:162-168` 忠实移植）→ 回复脚本 `window._rpc[0].resolve(null)`，**Go 侧零记录**（对没有 `_rpc[0]` 槽的页面还抛一个被吞的 TypeError）。**改后**：slog 记录＋`reject(<具名句>)`。**这是本腿落点**（可达）。
- **P3 无折入世界的 `postMessage`（钩子 `typeof window.wispDispatch !== "function"` 的原生出口分支）**：消息原样到 `msgcb`，与 P2 同落点 ⇒ 同一发修法覆盖。⛔ 本腿**没在真窗里造**这个门不可调用的世界（35r2 的 M-B2/B3 是模型，非真窗）；具名登记（§5）。
- **P4 `msgcb` 的 invalid-JSON 支（连 RPC 帧都不是的消息）**：改前＝库形 `log.Printf("invalid RPC message: %v")`（stderr 一行，非 resident 的 slog/JSONL 通道）；**改后**＝同通道 `slog.Warn`（`panel host: inbound message is not an RPC frame - received and dropped, no door and no reply (ticket 253 AC#3)`）。⛔ 这一支**结构上无页面回执**（帧里没有 id 可寻址），不硬造；本腿未单独发这一支（真窗页面 JS 没有非 JSON 出口的现成载具），具名没做（§5）。
- **P5 上游库 `go-webview2@…/webview.go`**：⛔ 一字节未动（禁区）。299-r1 之后生产不再经它；「库内不可达」今天只对 299 之前的历史成立。
- **P6 页面侧的结构边界（②的闭合口径）**：回执只能经 `window._rpc[id]` 槽。页面留了槽的调用 ⇒ 现在拿到 `reject(<具名句>)` 而不是 `resolve(null)`；**裸发且无槽**（id=0 无槽）的调用 ⇒ `resolve`/`reject` 脚本都抛被吞的 TypeError，**页面结构上拿不到东西**——回执需要返回地址，本仓不许改 `frontend/**`、不许加第二枚 binding ⇒ Go 侧记录是那一形唯一可达的声。

## 2. 落修（最小形、形状自裁）

- `msgcb`（`panel_control_windows.go:711`）：绑定门判定前置（`boundDoor` `:704`），不认的方法 ⇒ `slog.Warn`（`:734`，attrs `method`/`id`/`params`/`reply`）＋ `reject("<具名句>")`；invalid-JSON 支 `:717` 落同名 slog。**slog 选型理由**：resident 进程的持久通道就是它（`logsink.go:144` 把 JSONL sink 装在 slog 上；`log.Printf` 只到 stderr）。
- ⛔ `callbinding` **一字未动**（`(nil,nil)` 分支保留＝对上游忠实；文件头与 RPC 段落各加一段「SECOND NAMED DEVIATION」把这次偏离写在自己脸上：`:48`、`:679`）。
- ⛔ 上游库零改；⛔ 未加/改名册方法、未加 binding（C17 名册零动，`d22scan` clean）。

## 3. 假腿读数（一发旧封套发一次；真窗，页面自报）

载具＝真窗：`startResidentPanel`（产品同一 STA 线程与 `Run()` 泵）＋页面把结果经产品自己的门送回（AC#14 定式）。假腿 JS＝注册 `_rpc[0]` 槽（库对无 id 封套推出的 id）→ `window.external.invoke(pageModeRequestEnvelope 逐字)` → 等到 settle 后报回。

- **① Go 侧现场（改后逐字，`logs/11` 第 3 行）**：`panel host: inbound RPC frame names no bound door - received and refused instead of resolved null (ticket 253 AC#3) method=panel.mode.request id=0 params=0 reply=面板请求被拒绝 [panel.mode.request]：这枚 RPC 帧的方法名没有已绑定的门（door），没有任何处理器接手；宿主已记录本次丢弃`
- **② 页面自己的话**：改前 `"RESOLVED(null)"`（`logs/10` 第 2 行）→ 改后 `"REJECTED('面板请求被拒绝 [panel.mode.request]：…宿主已记录本次丢弃')"`（`logs/11` 第 2 行）。
- **③ 路由器未被到达**：该 requestId 在 handler 名册里 0 枚（断言在案）⇒ 读的是管道的死，不是 composer 拒绝。
- **正控（一）**：同一装配、同一封套走 `chrome.webview.postMessage` ⇒ `RESOLVED('')` ＋该 requestId 恰 1 枚到达 handler ＋`before` 之后**零**条新告警（`logs/11` 第 4 行）——正常门**不被新审计误伤**。
- **正控（二）**：相关既有族 rc=0（§4）。
- **反控（摘掉审计）**：`git show HEAD:cmd/wisp/panel_control_windows.go` 还原成 HEAD blob（sha `ae714aee…`）后跑**同一发** ⇒ 页面 `RESOLVED(null)`、criterion ① 无记录、rc=1（`logs/10`）；跑完逐字还原、sha 回 `2fe49342…`。

## 4. 门禁（逐枚 rc；命令与输出在 `logs/14-gates-abc.txt`／`logs/13`）

- `go build ./cmd/wisp/` **rc=0**；`go vet ./cmd/wisp/` **rc=0**；`gofumpt -l` 两枚动过件**空**（rc=0）。
- `sh scripts/d22scan.sh` **clean rc=0**（正控 35/0/0 先绿；分母 bans#1-5 internal/=236、cmd/=40、ban#8 internal/=550、cmd/=126）。
- `sh scripts/check-path-length-budget.sh` **VERDICT GREEN**（60=60，not-in-roster=0）。
- `-run` 相关族：headless transport 族 **15 PASS rc=0**（35r1/35r2/35r3/253r1 各枚）；`TestPanelHostRealWindowHopAndLifecycle` rc=0；`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` rc=0；`TestAC14AwaitedBindingReplyReachesThePage` rc=0；本腿 `-run Test253r6` rc=0（**各自独立 invocation**）。
- ★**本机环境读数（非本腿回归，pre-existing）**：同一个 test 二进制里连开第二枚真窗会**偶发**在建窗中途静默猝死（无 panic 文案、rc=1）。在**原始树**上用两枚既有用例复现（AC13→AC14 同发：第一枚 PASS、第二枚猝死，`logs/13`）；晚些时候 `AC14 -count=2` 又两枚全绿 ⇒ **间歇**现象（共享机器/桌面争用的形状），与本腿改动无关。本腿因此把读法合成**一枚单窗用例**（一个窗取三读数）以减半暴露面。

## 5. 具名没做／归口

- 具名没做：①P3（门不可调用世界）与 P4（非 JSON 消息）**未在真窗单独发**（无现成载具/无页面回执），只做了静态逐点裁；②未跑 `-tags winlive`（35 的 `panel_transport_live_35v2` 机密度归编排者）；③未跑整包（314 死区）；④未碰 `frontend/**`／`design/**`／三枚冻结件／`thresholds.go`／golden／`allowlist.txt`／`scripts/**`／`.github/**`；⑤未 push。
- ⚠ `scripts` 侧一条具名读数（非本腿地界，317-r1 在写 scripts）：`sh scripts/d22scan.sh` 在**未导出 `GOTMPDIR`** 时，正控步在 `tools/d22scan/runtests.sh:92` 报 `LOCALAPPDATAwisp: unbound variable`、rc=1；按派单规矩带反斜杠形 `GOTMPDIR` 复跑即 clean。
- 归口：界面侧唯一需要配合的仍是票面那条「调用统一走 `wispDispatch`」（本腿落点不需要界面侧新增配合：本腿的两条读数分别走旧门与页面自己的 postMessage）。

## 6. 盘上件（committed＝本文件＋下表「提交」列；raw 留盘不提交）

| 盘上文件（`.scratch/wisp/probes/253/r6/`） | 是什么 | 提交 |
|---|---|---|
| `logs/10-reverse-control-pristine.txt` | 反控：原始树同发读数（`RESOLVED(null)`／无记录／rc=1） | ✅ |
| `logs/11-post-fix-readings.txt` | 改后同发的三个读数＋正控（①/②/正控逐字） | ✅ |
| `logs/13-families-and-window-flake.txt` | headless 族 15 PASS＋两枚真窗回归（各自独立 invocation）＋**原始树两窗连开猝死**对照 | ✅ |
| `logs/14-gates-abc.txt` | build/vet/gofumpt/sha256 ＋ d22scan／长度门尾巴 | ✅ |
| `logs/reverse-control-single-window.txt`／`post-fix-single-window.txt` | 上面 10/11 的 **raw 全文** | 留盘 |
| `logs/post-fix-headless-family.txt`／`post-fix-realwindow-regression.txt`／`post-fix-ac13.txt`／`post-fix-ac14.txt` | 族读数 raw | 留盘 |
| `logs/two-windows-pristine-raw.txt`／`two-ac-tests-raw.txt`／`ac14-count2-raw.txt` | 两窗猝死对照 raw | 留盘 |
| `logs/gate-d22scan.txt`／`gate-length.txt`／`fold-count2.txt`／`fold-again-{1,2,3}.txt`／`r6-pair-run{1,2}.txt`／`combofail-raw.txt`／`reverse-control-no-audit.txt` 及合并前旧测试的两枚提取件 | 过程 raw（含合并试跑；旧提取件是合并前两用例版，已被 10/11 取代） | 留盘 |
| `panel_control_windows.go.253r6-mine` | 反控用的改动版备份（sha `2fe49342…`） | 留盘 |
