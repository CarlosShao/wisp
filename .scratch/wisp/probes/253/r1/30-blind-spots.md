# 253-r1 · 这枚尺看不见什么（30-blind-spots，逐枚具名）

时刻 `2026-10-08 16:53 +0800`（`date` 现跑）。每条标 **〔实测〕**＝盘上种出来的读数，或 **〔推理〕**＝从尺的代码形状推、我没种反形。派单要求"说不出就写'我没造出反形'"——下面 #1/#3/#10 有实测，其余是〔推理〕，我逐枚标明。

## #1 ★JS 模板文本里 `%[1]s` 有没有真被用到＝看不见〔实测，M5〕

尺断的是"`w.Init` 交出去的那段脚本由 `panelDispatchBinding` 经 `fmt.Sprintf` 生成"（`fmt.Sprintf` 的**实参**层锁步），⛔ 它不读那段 JS 文本本身。
种法＝`cmd/wisp/panel_host_windows.go:790` 把 `window.%[1]s(message)` 改成 `window.wispNothing(message)`、实参一枚不动 ⇒ **尺 rc=0 全绿**，而那一刻页面转发到的门名 Go 从没绑过。
⇒ 那一格今天仍**只有能力形那半**守着（`cmd/wisp/panel_transport_35r1_test.go:124-131` 读 `Init` 串含不含 `window.wispDispatch` ＋ `:195→:203` 真调 `installPanelTransport`）。我没为它扩射程——扩了正撞 `A717` §5 那句"⛔ 重复劳动"。

## #2 名册＝写死的枚名，"改名＋同批改名册"是绿的〔推理，代码形状；M4 反证了更窄的那一式〕

`doorRoster253r1` 写死 `"wispDispatch"`。改常量值但⛔ 不改名册 → **红三面（M4 实测）**；改常量值**且**同批改名册 → 绿。后者是本仓名册尺的"已评审动作"定义（同形出处＝`cmd/wisp/panel_locked_naming_33r11_windows_test.go:66-80` 的名册：`Adding a name here is the reviewed act`）。
⇒ 后果具名：**这枚尺不裁决"这个名字对不对"，它只裁决"名字与绑定与页面脚本三者是不是同一枚、且是被评审过的那一枚"**。要问"前端真在叫这个名字吗"得去读 `frontend/**`——⛔ 那是派单禁区，本程一枚没读。

## #3 ★名册外的 `Bind`＝不在射程〔实测（干净树绿＋grep 同存）〕

尺只走 `installPanelTransport` 那一枚函数体内的 `Bind`。这一刻整包产码里有**两枚** `Bind`：
- `cmd/wisp/panel_host_windows.go:802`（传输门，在射程内）
- `cmd/wisp/panel_host_windows.go:852` `w.Bind("wispProbeRT", ...)`（`firstRoundTrip` 的一次性探针门，⛔ 我按形状具名排除）
尺在干净树上**全绿**（`10-gates.md` §4，rc=0）而 `:852` 那枚门名同时在盘上存在 ⇒ 这就是实证：**枚一枚不在传输函数里绑的门，这枚尺既不数它也不评审它**。
⇒ 若日后有人在 `bringUp` 或别的产码函数里再绑一枚页面可叫的门，本尺**不会红**。这条是刻意留在票 253 `AC#2` 那把名册尺（`internal/panel/inbound_roster_253_test.go`）之外的：一枚管"绑定处的锁步"，一枚管"路由名的封闭集"，⛔ 不许混成一枚（`probes/ruler-dedup-1/ruler-dedup-1.md:25` R6 判"两枚封闭集不在同一宇宙"，我采）。

## #4 只读本包目录〔推理〕

射程＝`os.ReadDir(".")` 拿到的 `cmd/wisp` 目录里非 `_test.go` 的 `.go`。别的包（`internal/panel`、`internal/agent/approval`）里若出现同名常量或同名门，本尺读不到也不判。

## #5 跨包路由名册不在射程＝票 253 AC#2 那格〔推理〕

本尺**不判**"入向路由名的封闭集"（`panel.mode.request`／`config.get`／`config.set` 那批）。那格已由 `253-r5` 交件（票面 §9 `:47`，`internal/panel/inbound_roster_253_test.go` 834 行），⛔ 我不重复也不"顺手扩"。

## #6 回调里到底干了什么＝看不见〔推理〕

尺只看 `Bind` 的**第一个实参**，不读第二个实参（那个闭包）。所以"名字绑对了但闭包喂错了对象"（比如不接 `dispatchRaw` 而接了别的）本尺不红。那一格仍归能力形那半（`panel_transport_35r1_test.go:176` 断的正是帧进了哪枚 door）。

## #7 `//go:build windows` ⇒ POSIX 分母为零〔实测：文件头 `:1`；后果属推理〕

本尺与票 33/35/255 那族同源带 `//go:build windows`。Linux 那几步编不到它＝它给的绿**只覆盖 windows 腿**（`probes/ruler-dedup-1/ruler-dedup-1.md:24` 给那族记过同一枚面，我这条复现）。改档位＝动别的腿的门禁面，本程不做。

## #8 尺读的是工作树，不是 HEAD〔推理，但本程已按规矩压住〕

`os.ReadDir`/`os.ReadFile` 吃的是盘上字节 ⇒ 如果起跑时 `cmd/wisp` 里有别人的未提交改动，我会把"别人的半成品"读成"现状"（本仓记过这形：票 154-r1 的派单里就写着"复跑 cmd/wisp 必须先取锚点的树，否则读到别人的半成品并记成没判据"）。
我这发的处置＝起手 `git status --porcelain -- cmd/wisp` **空**（`00-anchor.md`），每发还原后**再量一次空**（`20-mutations.md` 逐条），落笔前再量一次空 ⇒ 本程读到的树与 HEAD 逐字一致（哈希对拉见 `10-gates.md` §1）。⚠ 这是**时刻读数**，⛔ 不保证后来者：任何人复跑前该自己再量一遍。

## #9 包里有文件解析不了＝尺红，但红因与门无关〔推理〕

尺对"读不到／解析不了"走 `t.Fatalf`（不许静默跳过，形状抄 `panel_geometry_255r6_range_windows_test.go:84`）。后果＝如果别人在同一棵工作树里塞进一枚语法坏掉的 `.go`，本尺会红在"could not parse"上，**看着像门的红其实不是**。后来者读到这枚尺的 `FAIL` 时，先认红句里是 `census could not parse` 还是 `AC#1 RED`。

## #10 空名册那支的文案不区分坏法〔实测（M2 与 M6 落在同一句）〕

M2（`w.Bind("", ...)`＝名字空）与 M6（`w.Bind((panelDispatchBinding), ...)`＝名字变成解析不了的形）**打的是同一句红**：
`panel_dispatch_binding_roster_253r1_windows_test.go:364: installPanelTransport (...) issued 1 Bind call(s) but bound zero readable door names`
⇒ **fail-closed 是对的**（两种坏法都不许绿），但**文案把"被摘掉"和"被间接化"压成同一种输出**。这条正是 `A716` §4 定式ⓑ（"尺不许把'读不到'和'没有'压成同一个输出"）的一个实例——我这枚是"两种读不出"压成一枚，轻一格，但仍具名认下。要拆成两句得改尺自己的文案，本程没做（⛔ 不留恒红、也不改判据凑绿）。

## 我没造出反形的，直说

- **`installPanelTransport` 出现两枚**那一支（尺的 `len(sites) > 1` ⇒ `t.Fatalf`）：⛔ 我没在盘上种过（种它要在第二枚产码文件里造一枚同名方法，会连累别的腿的编译），只按代码形状陈述。
- **常量整枚不存在**那一支（尺的 `!constFound` ⇒ `t.Fatalf`）：⛔ 没种，同上（删常量会打断 `:792`/`:802` 两处引用＝编译门先红，尺跑不到，那是 M6 首发踩过的同类坑）。
- **`Init` 传进来的不是那枚 var 而是别的 var**（尺会红在 `:398`）：⛔ 没种过（与 M3 相邻，M3 用的是"换成字面量"那形）。
