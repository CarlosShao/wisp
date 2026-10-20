# 351-a1 件02 = `AC#1` 定性（那一红属哪类）＋三支各自盘上凭据 ⇄ 到达数分子

★本腿 0 枚 `go test`、⛔ 造任何测量。凭据＝五份已入库归一日志（`346-a1` 留档）＋`git grep HEAD` 现读源码＋
`git show HEAD:` 的 blob。四数＝本腿脚下重跑（件00 第二节，与 346-a1 留档**逐数同**）。

## 一、定性＝**仪器形状（运行学时序）**；两支排除在下文第三、四节

红的那一格响的**不是**"卡片文案缺那枚 key"，而是"卡片正文一行都还没进 stdout"。机制读自 HEAD 源码（逐字行号）：

```
cmd/wisp/run.go:1421  func (u *consoleApprovalUI) Prompt(_ context.Context, p approval.Prompt) error {
cmd/wisp/run.go:1423      u.cards++                                  // windowCount 的分子在这一刻就已经 +1
cmd/wisp/run.go:1431      u.live.record(p)                           // ★ledger 落地＝awaitCard 轮询的那一面
cmd/wisp/run.go:1438      u.run.bookWaitingState("ui-prompt")        // 中间还隔一次 auditf
cmd/wisp/run.go:1441  fmt.Fprintf(u.out, "\n[确认 %s %s] %s\n", p.Level, p.Tool, p.Reason)  // ★stdout 才落地
```

`:1425-1429` 的注释自己写着「The grant is booked before anything is printed. The order is the point」
＝**先记账、后印字**是产码刻意为之的形状。而仪器这一侧：

```
cmd/wisp/config_reload_223_test.go:542  card := r.awaitCard(t, configReloadTool)   // 轮 :259 liveCards.pending()
cmd/wisp/config_reload_223_test.go:543  shown := r.h.out.String()                  // ★单次读、零同步边
cmd/wisp/config_reload_223_test.go:544-545  if !Contains(shown,"risk.permission_mode"){ t.Errorf(...) }
```

⇒ 存在一个**产码保证的窗口**：ledger 已可见（`:1431`）⇄ stdout 未写（`:1441` 之前）。
`awaitCard` 每 20ms 醒一次（`:256`），落在这个窗口里就返回；`:543` 于是读到"只有两条启动横幅"的 stdout。
★窗口⛔ 是我推的——它在五份日志里有直接读数（第二节的到达数），机制那一跳是源码现读、⛔ 猜测。

## 二、★"该形状到没到达"的分子（⛔ 拿总数当到达）

分子尺＝逐台面 grep 两条**具名红句**（`-SimpleMatch`、件＝归一件）：
`the console did not render the reload card`（＝`:403`，"卡片正文⛔ 在 stdout"那一格）与
`the card does not name`（＝`:406`/`:545`，"文案含 key"那一格）。

| 台面 | 该枚结果（尺＝`^--- (PASS\|FAIL): <名>`） | `:403` 型到达 | `:54x`/`:406` 型到达 |
|---|---|---|---|
| CI `39244425666`（`win.norm`，本枚红） | **FAIL** L1844 | 0（该枚⛔ 有在场格＝射程缺陷） | **1** L1841 |
| CI `39101060753`（`old39101`，本枚 **PASS** L6480） | PASS | **1** L6433（fs 那枚红） | **1** L6435（**与 `:403` 同段成对响**） |
| CI `39071959151` | PASS L3304 | 0 | 0 |
| CI `39124122121` | PASS L331 | 0 | 0 |
| 本机 `orch227-b/full.log` | PASS L352 | 0 | 0 |
| CI `39244425666/test-core`（`core.norm`） | ⛔ 现身＝⛔ 射程 | — | — |

⇒ **到达数＝2 台面 / 跑到 5 台面**；且这 2 次到达分属两枚用例（本票那枚＋它的 fs 姊妹枚），
**同一形状**。★正控（防"中文 needle 失灵"）＝同一把尺在同一件上命中已知在场句：
`win.norm` 的 `热加载`＝4、`答复监听`＝1（原始 ⇄ `-Encoding utf8` 两种读法**同数**）⇒ 尺灵。
★`core.norm` 的射程正控＝同件 `TestManager` 命中 24、`cmd/wisp` 命中 16 ⇒ ⛔ 失明；
`TestTicket223` 于 `core.norm`＝0 ⇒ 该包根本⛔ 在 test-core 射程（Linux 档），⛔ 读成"绿"。

**决定性一条（这就是⛔ 判"真缺陷"的凭据）**：两枚到达里，"文案含 key"那格**从未单独响过**——
`39101060753` 那次是 `:403`＋`:406` **一起**响（同一次单读、卡片整段缺位），
`39244425666` 那次的红句 dump（`win.norm` L1841-1842，逐字只两条横幅）同样是**卡片零字节**。
⇒ 红发的观测⛔ 是"卡片在场但没点名"，而是"卡片还没印出来"。

## 三、〔真缺陷〕那支＝**排除**（凭据在本枚红发自己身上）

1. 产码把 keys 逐字排进卡片文案（`cmd/wisp/config_reload.go:283-292`，`Reason: fmt.Sprintf("config.toml 里手改的这一条会放宽 [%s]：%v。…", section, keys)`），
   经 `internal/agent/approval/gate.go:606` `Reason: d.Reason` → `ui.go:35` → `run.go:1441` 原样打印＝**中间无人改写这枚串**。
2. ★**同一发红日志里就有那枚 keys 的到达读数**：`win.norm` L1843
   `config: locked loosening approved via L2 re-confirmation (D36 rule 1) section=risk keys=[risk.permission_mode]`。
   这句的 `keys` 与喂给卡片 `Reason` 的是**同一次调用的同一枚切片**（`internal/config/manager.go:180 confirm(c.section, c.keys)`
   ⇄ `:184-185` 打的还是 `c.keys`）⇒ 那一发 `keys` 确实携带 `risk.permission_mode`，
   于是"卡片一旦印出就⛔ 点名"这一支**⛔ 成立**。
3. 盘上既有**已渲染**卡片正文逐字（含 key 名）＝`.scratch/wisp/probes/221/v1/logs/ac5-full.txt` L142：
   `[确认 L2 config.reload] config.toml 里手改的这一条会放宽 [fs]：[fs.allowed_dirs]。…`
   ＝"这一跳真的会把 keys 印出来"的实测形状。★⚠ 诚实边界：盘上**没有** `[risk]：[risk.permission_mode]` 的已渲染件
   （尺＝`git grep -F -n '会放宽 [risk]' HEAD -- ':!.scratch/wisp/probes/351/a1'`＝**0 命中**；★那把不带排除的尺
   （`git grep -F '会放宽 [risk]' HEAD`）现在＝1 命中＝**吃进我自己这份件的这一句描述**＝`A1045`⇄`A1033` 那族
   "尺把自己的归档件吃进去"在本腿身上复发一发，具名记我＝**写档之后该尺必排除本腿写面**），
   所以第 3 条只证 `fs` 那支的渲染形状；
   `risk` 那支靠第 1＋2 条（同一条码路＋同一次调用的 keys 读数）。⛔ 把这条写满＝写"我⛔ 量到的"。
   ★同族另一枚读数＝尺 `git grep -F -n '会放宽 [fs]' HEAD`＝3 命中里有 **1 枚是另一句产码**
   （`cmd/wisp/approval_always.go:92`＝D45 长期规则那句、同词⛔ 同句）⇒ 凡"含某串"的尺⛔ 自动等于同一句，
   本件所有 needle 一律带文件＋行号定位。
4. `internal/config/manager_test.go:170-171` 已把"钩子收到的 keys ＝ `fs.allowed_dirs`"钉住（⛔ 本腿跑、引用既有件）。

## 四、〔CI 环境〕那支＝**只算放大器、⛔ 独立成因**；机主常驻 `wisp.exe` 抢键那形＝**不在射程**

- ★派单点名的反例本腿独立复现（分子尺＝逐台面 `REFUSED` 行总数，`-SimpleMatch`）：
  红发 `win`＝**122** ⇄ 绿发 `old39071`＝298／`old39101`＝254／`old39124`＝326／本机＝269
  ⇒ **红发总数最小**＝总通知数与红**反相关**。⇒ ⛔ 任何"通知很多＝形状到达"的读法在这一发上直接反号；
  只有第二节的逐台面具名 needle 才是到达数。（`result=deny` 于五件＝0 命中，正控见第二节尺灵。）
- runner 负载这一维：盘上只⛔ 够写"未知"——`gh run view --json jobs` 给 `startedAt/completedAt`
  （`39244425666`：test-windows 11:24:41Z→11:33:12Z＝同 run 里最长的一枚，⛔ 由此推负载），
  本腿⛔ 任何 runner CPU／负载读数 ⇒ 具名**未知**（与 `346-a1` 件03 第二节同判语）。
- ★"机主常驻 `wisp.exe` 抢键"那形＝**⛔ 射程**：CI runner 上的进程与机主常驻进程无关；且红句 dump 里没有任何热键／占用迹象
  （尺＝红发窗口 L1780..1860 grep `hotkey|occupied|RegisterHotKey|ball:`＝**0 命中**，正控＝同件 L1841 已知在场）。
- ★既有独立读数同向（⛔ 采信、只对回）＝`289-v1` 曾以 **21 发**复跑（隔离×3、`-count=5`、组串、`-cpu=1`×2、
  5 进程并发×5、原状×2）**零复现**，且并发把 2.0s 压成 2.8–3.0s 仍⛔ 红（`docs/reports/pending-and-issues.md:14857` 在册）
  ⇒ 单纯加负载⛔ 能造那枚窗口＝与本节判语一致：**成因是形状，负载只决定它这次有没有被撞上。**
  ⛔ 据此结案"flake 已定案"（票面禁区那条我照守：本件给的是**定性＋机制＋分子**，⛔ 是"已知摆动"那枚弃权章）。
