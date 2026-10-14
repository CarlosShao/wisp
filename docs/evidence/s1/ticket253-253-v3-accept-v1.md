# ticket 253 (253-v3): AC#3 非实现者验收裁决 — 成立（附 1 条具名更正・非阻塞）；AC#4 结线建议＝三格各自结线

**腿**：`253-v3`（验收腿，非实现者；只读＋自跑＋两窗 overlay 突变，已逐字节还原）。
**被裁**：落地腿 `253-r6`＝产码 `2298f7d2`（`cmd/wisp/panel_control_windows.go`，sha256 `2fe49342…`）＋测试 `cmd/wisp/panel_unbound_253r6_windows_test.go`（sha256 `53ed0947…`）＋交付 `89d29fce`。
**判语**：⛔ 本腿未翻框、未动 `Status:`／票名；下面每一条都是我自己取的读数（logs 在 `.scratch/wisp/probes/253/v3/logs/`，锚与机制链在 `probes/253/v3/00-anchor.md`）。

## 0. 结论先行

1. **AC#3＝成立**。判据①（假腿旧封套 ⇒ Go 侧可读现场）与判据②（页面拿到的不是无声 `null`）在**真窗**上逐条成立，且落点/正控/反控都由我独立复跑（§2/§3）。
2. **附一条具名更正（不影响成立、但必须收口时留档）**：`253-r6` 件 §1-P4 与实现自述把「invalid-JSON 支」改前描写成「stderr 一行、**非** resident 的 slog/JSONL 通道」——**不成立**。stdlib `slog.SetDefault` 把 `log` 包默认输出桥进当前 default handler（Go 源 `log/slog/logger.go:62-73`），而 resident 启动即 `installLogSink` 把 tee 设为 default（`resident_windows.go:101` → `logsink.go:153-161`）⇒ 改前那行 `log.Printf("invalid RPC message…")` **会落 resident JSONL**（INFO、无 attrs）。我的 (C) 改前读数当场证（capture 收到该行）。该支修法的**真实**收益＝WARN＋结构化 attrs（`err`/`bytes`）＋具名丢弃句＋不再依赖那条隐式桥；「补上了通道」这句要收回。
3. **AC#4 结线建议**（翻与不翻归编排者）：**可以按"三格各自结线"收口**，逐格留档句与撤销口令见 §5；⛔ 不许合并成"入向修好了"。
4. 本腿**没做**：`-tags winlive`、整包 `go test ./cmd/wisp`（314 死区，r6 同款具名）、`TestPanelHostRealWindowHopAndLifecycle`／AC13／AC14 的复跑（我的真窗预算给了被裁对象的判据本身；如需可再派）。真窗共 5 发，全 rc 如 §3。

## 1. 机制链自核（我现读，行号＝HEAD；⛔ 不引 r6 的件）

- 管道已是本仓代码：`panel_control_windows.go:339` `browser.MessageCallback = c.msgcb`；库侧 `edge.Chromium.MessageReceived → MessageCallback`（模块缓存 `pkg/edge/chromium.go:233-241`）⇒ 票面「静默面属第三方」在生产路径上**已过期**——成立。
- **改前静默 null 的必要条件（代码链）**：封套是合法 JSON ⇒ `json.Unmarshal` 成功；`method` 不在 `bindings` ⇒ `callbinding:761-763` 回 `(nil,nil)`、`err==nil` ⇒ 落 `json.Marshal(nil)="null"` 支 ⇒ `Eval("window._rpc[0].resolve(null); …")`；该支沿途无日志调用。⇒ 有槽页面得 `null`、Go 零记录。
- **三分界**：①旧门 `window.external.invoke(封套)`＝35/303 钩子（`panel_host_windows.go:813-832`）的 `ex.invoke` 包装把 `inside=true`，库 shim（`pkg/edge/chromium.go:112`）调用时读属性 ⇒ 走原生口 ⇒ 原封到 `msgcb`（本次落点）；②页面自己的 `chrome.webview.postMessage(封套)`（`inside=false` 且门可调）⇒ 折入 `wispDispatch`（35/303 已覆盖，有声）；③P3＝`wispDispatch` 不可调用 ⇒ 同走原生口、同落点。
- 生产绑定门全集＝`wispDispatch`（`panel_host_windows.go:842`）＋瞬态 `wispProbeRT`（`:892`）⇒ `panel.mode.request` 无门。

## 2. 判据①②逐条凭据（真窗，页面自报＝唯一读法）

改后（我自跑，shot4 rc=0；raw `logs/08-…`、摘取 `logs/03-realwindow-postfix.txt`）：

- **(A) 我自种的形：显式 id=977、method=`legacy.echo`（r6 没种过）、params=["x"]** ⇒ 页面：`REJECTED('面板请求被拒绝 [legacy.echo]：这枚 RPC 帧的方法名没有已绑定的门（door），没有任何处理器接手；宿主已记录本次丢弃')`；slog 逐字：`panel host: inbound RPC frame names no bound door - received and refused instead of resolved null (ticket 253 AC#3) method=legacy.echo id=977 params=1 reply=…`⇒ 驳回句/记录都带**帧自己的**名字与数字（不是 id=0 的特写）；该 requestId 在 handler 名册 **0 枚**。
- **(B) 我自种的形：裸发且无 `_rpc` 槽**（页面只装 `_rpc["0"]` 的 getter 陷阱）⇒ `getterHit=true`（宿主回执脚本确曾寻址该槽）＋页面**拿不到任何东西**（getter 返 undefined、`reject` 抛被吞 TypeError）＋Go 側 slog 记录在案。⇒ **r6 对"裸发无槽结构不可达、闭合只在 Go 记录"的读法与我一致**（页面侧结构不可达；能变的只有 Go 那一条记录）。
- **(C) 我自种的形：非 JSON 消息**（r6 未真窗发过）⇒ slog 逐字：`panel host: inbound message is not an RPC frame - received and dropped, no door and no reply (ticket 253 AC#3) err=invalid character 'v' after top-level value bytes=34`。
- **(E) 我自种的形：P3（门瞬时不可调用→postMessage 走原生口）** ⇒ 与 (A) 同形 `REJECTED('…[panel.mode.request]…')`＋恰多 1 条 slog 记录（0→1 计数在断言里）。
- r6 自己的用例我复跑（shot1 rc=0）：fake leg `REJECTED('…宿主已记录本次丢弃')`＋slog 行逐字＋folded 路 `RESOLVED('…不是面板 composer 通路的能力入口')`＋`INBOUND-DISPATCH` 审计行——与它留档逐字相同。

## 3. 正控与反控（自跑逐字）

- **正控不误伤**（我自跑）：(D) 直接 `window.wispDispatch` 正门调用 ⇒ `RESOLVED('')`＋该 requestId 在 handler **恰 1 枚**＋控制窗口内**零**条新"no bound door"告警；r6 的反方向正控（同封套走 postMessage）我复跑亦绿（shot1）。(A)/(E) 的 handler 计数 **0**（死的是管道、不是 composer 拒绝）。
- **反控（改前 blob overlay，sha 我现跑＝`ae714aee…` 与 r6 报逐字相同）**：同一测试件 rc=1——(A) `"RESOLVED(null)"`、(E) `"RESOLVED(null)"`、(B) `"getterHit=true"` 且无记录、(C) 只见桥过来的 `invalid RPC message: …`（见 §0.2 更正）、(D) 正控仍 `"RESOLVED('')"`。⇒ 改后两个判据的红都是**被裁码造成的**，不是环境。
- 还原：产码 sha256 回 `2fe49342…`、`git hash-object` == HEAD blob `84ff38c9…`、`git status --porcelain cmd/wisp internal/panel` 空；overlay 测试件移出仓外（master 留 `probes/253/v3/`）。

## 4. 门禁（我现跑，各 rc；`logs/01-gates.txt`／`02-families.txt`）

- `go build ./cmd/wisp/` rc=0；`go vet ./cmd/wisp/` rc=0；`gofumpt -l`（GOPATH/bin 形）两枚动过件**空** rc=0；`gofmt -l` 空。
- `-run 'Test253r6TheSilentDropBranchIsLoud'` **rc=0**（单窗）。
- headless transport 族 15 枚（一发 invocation）**15 PASS rc=0**（含 `TestUnforwardedPageEnvelopeDiesAtTheUnboundSlot`——⚠ 它模型里仍写死库形 `resolve(null)` 回执，是 35 院的自校准、未随本仓落点更新；⛔ 非本票地界，仅具名）。
- `sh scripts/d22scan.sh` **clean rc=0**（正控 35/0/0；bans#8 internal/=550、cmd/=126）。
- `sh scripts/check-path-length-budget.sh` **VERDICT GREEN**（roster entries=60=60；byte 视图与字符视图同选 60 枚）。

## 5. AC#4 结线建议（三格各自结线；留档句措辞）

**建议：可以收口（AC#4 可翻），但收口件必须分三格写、且每格带自己的读数与条件。** 措辞建议（逐格，⛔ 不许合并）：

- ①AC#1（钉的射程）＝**成立附条件**：凭据＝非实现者腿 `253-v1`（`probes/253/v1/10-attacks.md`／`20-blind-spots.md`／`99-verdict.md`）＋编排者 §10 复现的 N1（`Init` tied 支看不见"转发目标被换"）与 N2（只许整腿取读数）；撤销口令「撤 253 AC#1」。
- ②AC#2（名册的射程）＝**成立附条件**：凭据＝非实现者腿 `253-v2` 五分母逐项＋各造正控＋反形不敏感对照（`docs/evidence/s1/ticket253-253-v2-accept-v1.md`）；条件＝"更严结案口径"下 ④ `default` 改名/内联误红形未种（保留项）；半量欠账裁定：`config.reload` **不算入向名**；撤销口令「撤 253 AC#2」。
- ③AC#3（产品行为）＝**成立**（本裁决）：判据①＝旧封套假腿改前零记录→改后 slog 具名句（我 (A)/(B)/(E) 三形自种）；判据②＝有槽形 `reject(具名句)`、无槽形结构不可达且闭合在 Go 记录；正控/反控见 §3；**附具名更正 1 条**（§0.2，(C) 支的通道 lore）；残余＝`-tags winlive`／整包未跑（r6 与本人同款具名）。撤销口令「撤 253 AC#3」。
- ④**"不许合并"这一条自身**：AC#4 的判据是"结案时逐格给读数"——上面三格各自的读数是**三个不同维度**（钉/名册/产品行为），三枚非实现者已成（`253-v1`／`253-v2`／`253-v3`）。建议结尾句逐字写"三格各自结线：钉的射程／名册的射程／产品行为，⛔ 不以'入向修好了'一言收口"。

## 6. 具名没做（本腿）

- 未跑 `-tags winlive`（r6 同款归口）；未跑整包 `go test ./cmd/wisp`（314 死区在册）；未复跑 AC13/AC14/RealWindowHop（真窗预算＝5 发全给了判据与反控；如需回归可再派）。
- 未读 `frontend/**`／`design/**`（本票禁令）；未碰三枚冻结件／`thresholds.go`／golden／`allowlist.txt`／`PLAN.md`／`specs/**`／`scripts/**`。
- 未把 overlay 测试件提交进 `cmd/wisp`（master 留在 `probes/253/v3/`，可重放）；未 push。

## 7. 足迹

- commit：只 1 笔（本文件＋`probes/253/v3/**`＋票 253 尾 1 bullet），显式 pathspec；⛔ 未 push。
- 真树突变 2 窗（overlay 测试＋改前产码），两窗都在跑完后逐字节还原（§3）；他腿窗口核对：执行窗内 HEAD 3 笔全为编排者件，`git log --name-only` 零 `cmd/wisp` 重叠。
