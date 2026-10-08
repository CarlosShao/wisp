# stale-claim-1 / 04 — 判不动的与看不见的

⛔ 本节的规矩：**不用"应该会／大概"填空**。每一枚只写两件事 —— *"判不动"* ＋ *"缺的是哪一发读数（命令原文）"*。
读数一律可用（本程未跑的都标了原因），ref 一律 `HEAD 6547fd30`。

---

## 1. 判不动的（15 枚，按表序）

| # | 出处 | 判不动的那一句 | 缺的那一发读数 | 为什么本程没跑 |
|---|---|---|---|---|
| B-01 | 表一 P16 `HEAD:cmd/wisp/models.go:36-39` | `There is still no reader of model bytes anywhere in `cmd/wisp`'s dependency graph … so the segment has no end point in the product today` | `git grep -nE "os\.Open\(|os\.ReadFile\(|ReadDir\(" HEAD -- ':(exclude).scratch' ':(exclude)*_test.go' 'internal/models/*.go' 'cmd/wisp/*.go'` ＋ `cmd/wisp` 的**依赖闭包**（"anywhere in the dependency graph" 这句的射程是闭包，不是 grep） | 闭包要 `go list -deps`，**派单明令禁跑**（`go doc`/`go list` 也不要跑） |
| B-02 | 表一 P17 `HEAD:internal/risk/provenance.go:64` | `no *plugin.DisposalScope reaches the plugin agent command in production today (160-c1 §2.2)` | `git grep -nE "DisposalScope" HEAD -- ':(exclude).scratch' ':(exclude)*_test.go' '*.go'`（调用形状那一步我没做）＋复跑 `160-c1 §2.2` 那把尺本身 | 我只读了句子，**没复跑尺**；引 `160-c1` 的读数不是我这一程的数 |
| B-03 | 表一 P18 `HEAD:internal/tools/bridge.go:872` | `Today no production code dispatches on the bridge except loop.go, so being` | 先把方法名抽出来：`git show HEAD:internal/tools/bridge.go \| grep -nE "^func \(b \*Bridge\)"`，再对该方法名跑带接收者的形状尺。**我先前那把 `.Dispatch(` 尺对本枚不可用**（命中全是同形异符号：`internal/models/bridge.go:40/46/59/64`、`internal/statemachine/machine.go:202`、`cmd/wisp/panel_host_windows.go:613`、`cmd/wisp/panel_resident_windows.go:360`、`cmd/balldebug/main.go:618`） | 同一次命令里没先做"方法名抽取"那一步；**这正是 §1.2 那条"别用裸符号名/裸方法名"的反例，我撞上了并具名** |
| B-04 | 表一 P19 `HEAD:internal/winsec/resolve.go:231` | `no caller's path passes through here; ResolvePath` | `git show HEAD:internal/winsec/resolve.go` 的 `:150-200` 与 `:330-345`（拿到 `:198`／`:340` 两枚调用点各自的**宿主函数名**），再对宿主函数跑入口可达尺 | 调用点数量我量到了（2 枚非测试），但这句说的不是数量、是"真发 OS 路径不流经此处"——**该命题的形状不是调用计数** |
| B-05 | 表一 P20 `HEAD:internal/ball/ball_windows.go:808-811` ＋ 表三 R10 | `no caller does that today / because the config poll and the card live on different legs` | 读 `git show HEAD:cmd/wisp/resident_windows.go \| sed -n '195,275p'` ＋ `git show HEAD:cmd/wisp/resident_ball_windows.go \| sed -n '200,330p'`，判 `hotReload258` 这条 rebind 链与 `:877 b.TakeEscForCancel()` 那条 card 链**是否同一条腿**（调用顺序／闭包同一性） | 两枚调用点各自量到了（`hotkey_reload.go:96`／`resident_approval_windows.go:877`），但"同腿/异腿"要读控制流，不是读 grep |
| B-06 | 表一 P06／表三 R02 | `TierOf` 的入口可达性第 3 层 | 从 `hotRowsFor`（`config_readers_255.go:234`）向上追：`git grep -nE "hotRowsFor" HEAD -- ':(exclude).scratch' '*.go'` 只回 `:234/:266/:292`（全在本文件），再要的是**本文件那两个调用者的宿主函数是谁、由 `run.go` 哪一行驱动** | 追到本文件内两枚调用点为止；第 3 层（入口）我没追完。**第 1/2 层结论不受影响** |
| B-07 | 表一 P05／表一·补 T03 | **产码与仪器直接矛盾**：`HEAD:internal/panel/composer_dispatch_test.go:440` 逐字 `no native host / WebView2 message channel is attached in this tree: no production source carries a …` ⇄ `HEAD:cmd/wisp/panel_host_windows.go:821 return m.disp.Handle(ctx, raw)` | 两把尺：① `git show HEAD:internal/panel/composer_dispatch_test.go \| sed -n '425,445p'` 读它到底断的是**哪一枚符号**（可能是"消息通道"而不是"`Handle` 有人调"）；② `git show HEAD:cmd/wisp/panel_host_windows.go \| sed -n '795,825p'` 读 `raw` 从哪来（真 WebView2 `postMessage` 回调 vs. 注入/测试路径） | ⛔ 我不能跑那枚测试去裁谁红谁绿（派单禁 `go test`，且此刻另有一条腿在 `cmd/wisp` 写码）。**这是本格最该有人去核的一枚**：它同时决定 P05（给用户看的那句文案）、P03、Q01、Q12 四枚的档位 |
| B-08 | 表二 Q05 | `它说未修码上 `rt.admitTask` 不存在、编译都过不去` | 需要"未修码那棵树"（某枚 parent commit）上的 `go build`，或 `git show <pre-fix-ref>:cmd/wisp/run.go \| grep -c admitTask` | 前者＝派单禁 `go build`；后者我只在 **HEAD** 上量了现状（`run.go:1008`、`internal/agent/loop.go:368-369`），**没有那棵树** |
| B-09 | 表三 R12 | `It is the first thing in this repository that calls panel.NewSnapshot / panel.NewComposerState from a running process rather than from a test` | `git grep -nE "NewSnapshot\(|NewComposerState\(" HEAD -- ':(exclude).scratch' '*.go'` 全谱，逐枚分"测试／非测试／侧程序(`cmd/balldebug`、`cmd/llmrecord`)" | 历史序（"first"）本程的尺给不出；**当前集**（"有没有别的非测试调用者"）我能给但没跑这一发 |
| B-10 | 表三 R13 | `runGitIndex is its only production caller.` —— **"its" 的先行词不在句内** | `git show HEAD:tools/d22scan/gitignore.go \| sed -n '286,300p'`（整块读，确认那段讲的是哪枚 helper） | 我只读了 `:287`、`:294`、`:297` 三行。⚠ **本枚是我这把尺的一次自打**：`git grep` 给的是单行，代词先行词在块里，**单行读数不足以定分档** |
| B-11 | 表一 P21 `HEAD:internal/ball/hotkey_windows.go:262` | `The other three slots are never standby: summon / mute / panel are modifier combinations by default` | 两把：① `HotkeyStandby` 的赋值点全谱尺；② `git show HEAD:internal/ball/*.go \| grep -n "DefaultHotkeys" ` 之后读那三枚默认绑定字面量 | 都不是"调用形状"尺，本程的定式覆盖不到（**它断的是数据默认值，不是符号可达性**） |
| B-12 | 表一 P29 `HEAD:internal/config/unwired.go:81` | `nothing populates the bOverrides map risk.Gate reads` | `git grep -nE "risk\.Gate\("` 已给 1 枚（`internal/tools/mode.go:98`）；缺的是 `overrides` 那枚实参的**生产生产者链**：读 `mode.go:80-100` ＋ `internal/tools/bridge.go:104/:145` 的 `Blacklist`/override 字段填充点 | 读到调用点为止，没往上追实参来源 |
| B-13 | 表一 P30 `HEAD:internal/config/validate.go:112` | `queue.go:88-89; ticket 267's census: no production caller passes one at all` | `git grep -nE "WarningLead:" HEAD -- ':(exclude).scratch' '*.go'` —— **组合字面量赋值**才是这枚的尺（"有没有生产者传"）。我这把只扫到字段/getter，没扫 `WarningLead:` 这个形状 | 同一次没补这一发。**注**：票 267 带 `-done`，所以这句是**归档裁决里的一条读数**，过期成本更高 |
| B-14 | 表二 Q10 | 票 255 `:89` 自陈的 cite 漂移：`config_readers_255.go:161` cite 了 `panel_host_windows.go:262/:392`、同族还 cite `:304` | `git show HEAD:cmd/wisp/panel_host_windows.go \| sed -n '150,175p;255,270p;300,310p;388,398p'` 逐枚比对被 cite 的内容是否还在那一行 | 本程判了"行号会漂"这条**规律**（P01 的 `:48->:50`、P14 的 `:261->:697`/`:205->:641`、Q06 的 `run.go:137->:181`），但没逐枚复量票 255 的那三处 cite |
| B-15 | 表二 Q14 | 台账比对：`docs/reports/pending-and-issues.md` 的 `A##`/`R##`/`Q##` 与 `docs/reports/HANDOVER.md` §4 | `git show HEAD:docs/reports/pending-and-issues.md \| grep -nE "零调用者|no production caller|生产调用者"` ＋ `git show HEAD:docs/reports/HANDOVER.md` 的最新时间戳节 | **本程主动排除**：本格只需产码与票面名册，且工作树里那两份可能被别的腿在飞改动，我只走 HEAD、且没把它们的读数引为依据。⇒ 这一格里**"哪一枚已登记过"我无从比对**，后续程不要以为我核过了 |

---

## 2. 看不见的（结构性盲区，具名，不填空）

1. **接口派发看不见调用者类型。** `internal/agent/loop.go:737 return l.opt.Tools.Execute(ctx, req)` 是接口形状；"谁实现了 `Tools`"要靠 `internal/tools/bridge.go` 的方法集与 `cmd/wisp/run.go` 的装配，**grep 调用点回答不了"到不到得了"**。表一 P18 卡在这上面。
2. **方法值 / 函数字段看不见。** `Manager.ConfirmLocked`、`OnRestartPending`、`bridge.Refresh` 是**函数字段**：赋值点＝"调用者存在"，尺是 `x.Field = ` 而不是 `x.Field(`。本程对 P08/P09 用了赋值形状尺（量到了），但对**把函数字段作为参数传下去**的那一族（如 `startResidentBall(rt.Registry, ra.vetoByEsc, hotCfg258, hotReload258, …)`，`resident_windows.go:217`）**没有尺**。
3. **`//go:build` 的 GOOS 分岔看不见。** `_windows.go` 与 `_other_test.go` 的调用点都算进"生产调用者"，但对某一台机器可能根本不编译。表一 P25（`d2d_windows.go`）与 P20/P21（ball 族）都受这一条限制；表二 Q03 的 `RebindHotkeys` 落点 `internal/ball/hotkey_reload.go` **没有 `_windows` 后缀**——这枚我是**按名字判的，没读它的 build tag**。
4. **前端那一侧我一行没扫。** 名册 pathspec 是 `'*.go'` ＋ `docs/*.md`。⇒ `frontend/**`（含 `frontend/dist/index.html`）里"那一跳"的另一半**完全在本程射程外**。⚠ 这直接影响 P05／B-07／Q01／Q12 四枚：`网页事件 → Go 那一跳` 的**网页侧**我没有任何读数。
   （顺带具名一枚**我在 commit message 里读到、但没有读数的**事实：`frontend/dist/index.html` mtime＝`09-27 10:59`，而 HEAD 只跟踪 `dist/.gitkeep` —— `go:embed` 解析的那 1044 字节来自哪版源码**未验**。⛔ 我不据此判任何东西，只登记"这一格是空的"。）
5. **测试枚数类断言一律量不了。** 例：表一 P34 的 `had nine green tests`、表二 Q02 的 `有 16 处测试引用`。前者要 `go test`（禁），后者我可以量但没跑（`git grep -c SealDir HEAD -- ':(exclude).scratch' '*_test.go'`）。⇒ 这类枚数**引的时候必须带 ref**，它们随 HEAD 动。
6. **`.scratch/wisp/probes/**` 里的复制件会让"枚数"虚高。** 同一枚产码注释在 `probes/174/r3/mut/`、`probes/183/a1/mut/`、`probes/221/r1/mutations/`、`probes/197/r3b/pre/` 等处被复制 **3–9 份**（起手那一发不加排除就输出 159 KB）。⇒ 本表所有计数**都是在 `:(exclude).scratch` 之下量的**；换一把不加排除的尺，枚数会假涨。
7. **工作树我一次未读。** 派单具名的三类在飞改动（`design/**` 的 16 删 4 改、`cmd/wisp` 的在飞、被重写的 `probes/**`）里，**`cmd/wisp` 那一批有可能已经在改 P01/P05/P39 这几枚**。⇒ 本表所有"已过期"都是 **HEAD 上的过期**，不是"工作树上的过期"。若另一条腿正在这几枚文件上接线，交件之后**行号与档位都可能再动**。
8. **"到得了 `main`"只有两条入口被我确证。** `main.go:91 case "run" -> cmdRun -> runTextTask`（`main.go:159`）与 `main.go:115 case "panel-inbound" -> cmdPanelInbound`（`main.go:120`），常驻腿经 `resident_windows.go:260 startResidentTaskSource`。**`main.go` 一共 172 行，其余子命令（`providers`、`panel-host`、`slo`、`balldebug`、`llmrecord` 等）我没有逐枚建入口尺** —— 凡本表写"入口＝某条腿"的，只有那两条是我确证过的。

---

## 3. 一句话事实（本程给的全部结论，⛔ 不含修法、⛔ 不新建判据）

- 复跑之后的**枚数与分档**（各表自洽，与 `00`/`01`/`03` 一致）：
  - 表一 产码注释 **39 枚** -> **已过期 15 ／ 判不动 9 ／ 仍成立 15**
  - 表一·补 `_test.go` **4 枚** -> **已过期 1（T01）／判不动 3（T02/T03/T04，其中 T03 与产码直接矛盾＝B-07）**
  - 表二 票面／派单 **14 枚** -> **已过期 7（Q01/Q03/Q04/Q06/Q07/Q08/Q12）／判不动或未取数 4（Q05/Q10/Q11/Q14）／仍成立 2（Q02/Q09）／规矩层读数 1（Q13）**
  - 表三 正向断言 **14 枚** -> **仍成立 10（含 R11"仍成立但脆"）／已过期 1（R10）／判不动 3（R12/R13/R14）**
  - 剔除未计：101 行命中里 **62 行**属运行时拒绝文案／fail-closed 行为分支／标识符自指，其中三组已在 `03` R14 具名剔除
- 台账里点名的那 5 枚，逐枚落点：
  - `resolve.go` -> 表一 **P19**，判不动（缺 B-04）
  - `SealDir` 生产零调用者 -> 表一·附 ＋ 表二 **Q02**，**仍成立**（票 132 未 `-done`）
  - `askConfirmation` 生产调用者 0 -> 表一 **P14**，**已过期**（`resident_approval_windows.go:700`）；但它住在 `AskOnTaskRoot` 里、而后者生产调用点仍是 **0**（`resident_windows.go:248` 那句 `This call is the caller` **已过期**，P15）
  - 票 114 那句"网页事件→Go 那一跳不存在" -> 表二 **Q01** ＋ 表一 **P39**，**已过期**（`internal/panel/composer_dispatch.go:155`）
  - `composer_dispatch.go:48` 那句 "NO production caller yet" -> 表一 **P01**：**位置在 `:50`（`:48` 已漂 2 行），且已过期**（三处生产入口，见 P01）
- **本程未改任何一句话，未跑任何 Go 构建/测试/文档命令，未 add/checkout/restore/clean/stash 任何文件，未 push。** 唯一写入＝本目录 5 枚 `.md`。
