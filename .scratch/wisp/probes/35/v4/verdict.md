# 35-v4 对抗验收裁决｜票 35 `:75` (a) 支那枚判据仪器（手写 JS 子集解释器）

裁决者＝`35-v4`（非实现者）。实现者＝写腿 `35-r4`＋编排者代跑反形。⛔ 翻勾归编排者，本文只给判。
起跑 HEAD `168a91f0747a7f71c36344dd48f75edd98fe7729`；交件时 HEAD 已被别人推进到 `8e98b8f0…`（不是我写的）。
全程零 push；⛔ 零真窗（我每一条 go 命令都不带 `-tags winlive`；那枚 winlive 文件由 build tag 排除，尺＝`logs/roster-coverage.txt`）。
产码与夹具一字未动：`git hash-object` 终态 `26b5de83b93a9a141f1546dc19a9b03ffa45ede0`（产码）／
`31c96f37d32adcd1625a2c01748d5b30016227b2`（夹具）＝起手锚与 `HEAD:` blob 逐枚相同，`git status --porcelain -- cmd/wisp/`＝0 行（`logs/repo-untouched.txt`＋本腿末次现量）。

## 读数一览（我自己在 `-overlay` 仓外拷贝上跑的 16 发，名册＝12 枚，与腿 r4 同一名册可比）
尺＝`logs/mut-summary.txt`（发次行＝`PASS=/FAIL=`）＋`logs/names-*.txt`（逐名）＋`logs/mut-landing.txt`（`all-one-line=YES`，15 枚各 `changed-line-count=2`）。
baseline（未突变）＝PASS=12 FAIL=0。⛔ 本表只算 pass 2；pass 1 的十三发是空操作 overlay 造出来的假绿，见"我这把尺的两处缺陷"。

| 突变（恰一行） | 改了哪一行 | 红名册 | 判 |
|---|---|---|---|
| `v-ma-hook` | 产码 `:680` `native.call(cw, message)`→`native(message)` | 9 枚（名与腿的 `orch-mut-ma-hook.txt` **逐名相同**，尺＝`diff` 两份名册＝空） | 复现腿的读数 ✓ |
| `v-ma-hook+detector-gone` | ＋夹具 `:1288` `if recv != cw {`→`if false {` | 0 枚（PASS=12） | 复现 ✓ ⇒ 级联红的归因成立 |
| `v-native-window` ★新 | `:680` →`native.call(window, message)`（**非 nil 的错 receiver**） | 同样 9 枚 | 探测器是**按身份**判的，不是伪装成 nil 检查 |
| `v-ma-hook+no-panic` ★新 | 夹具 `:1290` 保留计数器、把 `panic(&jsPanic{…})` 换成 `_ = (&jsPanic{…})` | 2 枚，含 M-A | M-A 的红**第二条独立载体**（`:1953` 的计数钉子）存在 |
| `v-guard-truthy` ★新 | `:680` `typeof window.%[1]s !== "function"`→`!window.%[1]s` | **0 枚＝绿** | ★**第三面，成立**（下文 §3） |
| `v-return-gone` ★新 | `:682` 去掉门那一支的 `return` | 0 枚＝绿 | 绿，但**等价**，不算面（论证见 §3） |
| `v-idem-guard-gone` ★新 | `:675` 去掉 `cw.__wispForwardInstalled` 幂等守卫 | 1 枚＝`:1786` 那枚 | 我的预判（"双包还是一次送达"）**错了**，牙在（见 §3） |
| `v-finally-stuck` ★新 | `:682` `finally { inside = false; }`→`{ inside = true; }` | 1 枚＝`:1764` 那枚 | **三枚具名用例全绿**＝覆盖面边界（见 §3） |
| `v-detector-gone` ★新 | 只摘探测器，产码干净 | 0 枚 | 探测器在正确代码上不做任何工＝纯引线，符合预期 |
| `v-ma-world-swap` ★新 | 夹具 `:1947` 载具世界 `false`→`true` | 恰 1 枚＝M-A | 台件被真选择；但红句在 `:1961` expectDelivered，不在三枚 receiver 钉子上 |
| `v-mb-world-swap` ★新 | 夹具 `:1977` `true`→`false` | 恰 1 枚＝M-B | 同上，反向 |
| `v-init-keep-bind` ★新 | 夹具 `:2037` `initScripts[1:]`→`[0:1]` | 恰 1 枚＝`:2026` | "门不在"那格是真造出来的（calibration `:2042` 有牙） |
| `v-guard-gone`（复跑腿的 mb-typeof-gone） | `:680` 整支 typeof 条件删 | 恰 1 枚＝`:2026` | 复现 ✓ |
| `v-set-gone`（复跑腿的 mb-set-gone） | 夹具 `:164` 静默拒覆写改成 `&& false` | 恰 1 枚＝`:1976` | 复现 ✓ |
| `v-guard-gone+label-mb` ★新文案形 | 上一条 ＋ `:2048` 的红句开头换成 `:1980` 的"M-B RED (wrong death)" | 恰 1 枚＝`:2026` | 红**不串名**（见 §2） |
| `v-ma-hook+label-ma` ★新文案形 | `v-ma-hook` ＋ `:1954` 的红句开头换成 `:1984` 的"M-B RED (fixture blind to writability)" | 9 枚，仍含 M-A 具名 | 同上 |
| `v-cap-loose`／`v-cap-tight` ★新 | 夹具 `:1184` 帽 `8`→`2`／`8`→`1` | 8 枚／9 枚（正确产码被判红） | 帽数有**下界凭据**（合法深度现量＝3），无上传感；见 §3 末 |

★前提也被我复量了（谁的自述都不引）：把**产码突变**配到 `2fc5f5c9^` 那份 1824 行旧夹具上（blob `7185ab56d75952247226feb404add664…`），
`ma-hook`＝9 PASS/0 FAIL、`guard-gone`＝9 PASS/0 FAIL、干净＝9 PASS/0 FAIL ⇒ **35-v2 那两面恒真今天由我第二次量到**，
而同一枚突变在今天的夹具上会红 ⇒ r4 这轮**确实改了颜色**。尺＝`logs/prem-replication-table-final.txt`（原始段 `logs/prem-replication.txt`）。
旧夹具里 `writable|configurable|defineProperty|unwritable` 现量 0 处、今天 22 处 ⇒ `:1966` 那句自述成立。

---

## 判 1｜M-A 的牙算不算数：**成立**（在模型内），并带一条粒度限制
1. 因果对＝我自己复跑的：`v-ma-hook` 9 红（`logs/run-v-ma-hook.txt`，名册与腿逐名相同，尺＝上面那行 `diff` 为空）
   ＋`v-ma-hook+detector-gone` 12 全绿（`logs/run-v-ma-hook+detector-gone.txt`）。级联红的**唯一变量**就是那枚 receiver 检查 ⇒ 归因支撑得住。
2. 我另外顶了一发它没试过的形状：`native.call(window, message)`（`v-native-window`）。它同样 9 红 ⇒ 探测器读的是**身份**（`:1288` 的 `recv != cw`），
   不是"有没有传第二个参数"的伪装的 nil 检查。**第三种"骗过它"的形状我在模型内没找到**：模型里 `cw` 每个文档只造一次（`:1249`），
   没有任何表达式能产出一个"与 cw 不同却等价"的对象；而"不捕获、调用时再读属性"那一支由 `:1250-1276` 的 hop 计数在 8 层处引线（`:1812` 那枚旧形钉子量到 re-entered 9 levels）。
3. 第二载体我也量了：摘掉 panic 只留计数器（`v-ma-hook+no-panic`），M-A 仍红 ⇒ 红的来源不只是"异常逃到 `thrown`"，`:1953` 那枚计数钉子也会落。
   反过来推：panic 在位时红句走的是**第一支**（`:1949`），与 `:1942-1944` 的 TOOTH 自述一致（这条我是**从两发读数推的**，没有单独把 M-A 的红句行号提出来量，⛔ 写成"直接量过"是假话）。
4. ⛔ 边界（不许被读成别的东西）：这枚探测器**证明的是夹具里"丢了 receiver 会红"**。"WebView2 真绑 receiver"不在这里面，那格在票 `:52`／`A684`／`A686` 的 winlive 凭据里，今天零次开窗。
5. 粒度限制（新发现，具名）：nil receiver 与错对象 receiver 两发给出**同一组 9 枚红名**，只在 panic 文案里可分（`called with undefined` vs `called with object`，尺＝`logs/mut-summary.txt` 的 red-reasons＋两份 run 文件）。
   ⇒ 红名册本身**不区分坏法种类**；以后谁按名册读"是哪种 receiver 丢失"必须连文案一起读。

## 判 2｜两枚"恰红 1 枚"的名字唯一性：**成立**，但文案层的混淆面是真的（不算塌）
1. 复现：`v-set-gone`→恰 `TestForwardingHookIsSilentlyUnarmedByANonWritableNativeExit`；`v-guard-gone`→恰 `TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsAbsent`（`logs/names-v-set-gone.txt`／`logs/names-v-guard-gone.txt`）。
2. 我自造的两发**合并文案**形（腿与编排者都没做过）：把 `:2048` 的红句开头换成 `:1980` 那句原话、把 `:1954` 的红句开头换成 `:1984` 那句原话，各自再配它对应的坏法（`v-guard-gone+label-mb`／`v-ma-hook+label-ma`）。
   读数：红**仍按用例名落**，没有串到另一枚（`logs/names-…label….txt` 逐名）。⇒ "恰红 1 枚"这个归因**不靠文案互异**，靠的是用例名，站得住。
3. ⚠ 但我顺手量出现仓里今天就有这层混淆：`:1976` 那枚用例的 4 条红句与 `:2026` 那枚用例的 2 条红句**共用同一个前缀"M-B RED"**（`:1980`/`:1984`/`:1992`/`:1996` vs `:2048`/`:2052`）。
   谁按红句关键词（不是按 `--- FAIL` 名）读日志，就会把"静默拒覆写那一面"和"门不在那一面"读成同一面——这正是票 `:63` 那三支被 35-v3 抓过的形。⛔ 我不改夹具（裁决权归我、翻勾与改码归编排者），**建议**：那两条各换成本面专属前缀（例如 `M-B1`／`M-B2`），并把这条记进 `:75` 的下一轮。
4. 台件不是装饰（三发新形）：`:1947`/`:1977` 的世界开关各自单独翻一次，各红恰 1 枚；`:2037` 的"丢脚本 0"换成"留脚本 0"后由 `:2042` 的 calibration 红 ⇒ 三枚用例的世界是**被选择**的，不是顺手全绿。
   ⚠ 一处必须写细：`v-ma-world-swap` 的红落在 `:1961`（expectDelivered），而它里面三枚 receiver 钉子（`:1949`/`:1953`/`:1957`）在**不可写世界里全过**（现量 `receiverOK=1 receiverLost=0`，`logs/reasons-for-verdict.txt`）
   ⇒ M-A 的"专指 receiver 那一面"是靠**世界选择＋送达钉子**一起成立的，那三枚钉子本身不排他。编排者转述里"M-A 是级联红所以弱一档"这句我按此**改读**：级联归因有对（判 1），但 M-A 单发读数在别的世界上与 M-B 不可分，这一半我原话补上。

## 判 3｜第三面：**成立一枚真面**（具名如下），另有两形作为覆盖面边界
**★ 真面（恒绿，产码坏掉而仪器不红）：`typeof window.wispDispatch !== "function"` 这一支只被"门不存在"这一种世界验过。**
- 坏法＝一行：`:680` 的守卫换成 `!window.%[1]s`（真值判断）。读数＝**PASS=12 FAIL=0**（`logs/run-v-guard-truthy.txt`＋`logs/names-v-guard-truthy.txt`）。
- 为什么这是坏法不是等价改写：两支只在"门存在但不可调用"时才分岔——`typeof x !== "function"` 会把一个真值非函数（对象／字符串）打回原生出口，`!x` 会放它过去然后 `x(message)` 抛 TypeError。
- 为什么仪器看不见：这台尺**没有那种世界**。门在页面里唯一的构造点是 `Bind` 注册的桩脚本（`:1229-1232` → `libBindStubScript35r2` `:97`，注入的是一个 function）；
  全仓没有任何一处把 `window[panelDispatchBinding]` 设成非函数（尺＝本腿 grep，命中的 `wispDispatch` 全是"调用它"或字符串比较，见 `logs/roster-coverage.txt` 同批现量）。
  `:2026` 造的世界是**门缺失**＝`undefined`，两支在那一点上取值相同。
- 与作者自己那份"未建模清单"的关系：`impl.md §4` 列的是 `defineProperty`／`getOwnPropertyDescriptor`／`configurable`／`delete`／accessor——**这一枚不在那张表上**，也不在 `:75` 框的 (a) 支措辞里（框只要求"drop-the-receiver 与 silently-ignored-override 两面各自红"）。
  而 `:2019-2021` 那句注释已经写了这枚用例要干的活是"the hook must fall back to the native exit and forward byte-for-byte, **not call a non-function**"——**今天它只证明了 `undefined` 那一格**。⇒ 具名报为 `:75` 的第三面：*"守卫的非函数分支缺世界"*。
  修法建议（不是本腿的活）：再加一枚用例，世界＝门存在但是真值非函数（Go 侧台件把 binding 值换成 `win.set(panelDispatchBinding, "not-a-function")` 即可，⛔ 不需要 descriptor），断言"必须走原生出口、逐字转发、不抛"。
**其余新形（都不算面，但都是覆盖面边界，具名）**
- `v-return-gone`（去掉门那一支的 `return`）＝12 全绿。我判它**不是面**，理由现量两条：页面侧那两处发送把 `postMessage` 当语句用、函数签名是 `void`（`frontend/src/lib/panel.ts:170-190` 我读了原行）；
  门答复的回程在模型里走 `msgcb → Eval("window._rpc[id].resolve(…)")`（夹具 `:1388-1438`），钉子是 `evalThrew`/`lastReply`，不走包裹器的返回值。
  ⇒ 只有当页面将来 `await` 那次 post 的返回值、或答复改成同步返回给调用方时，这一格才变成面。今天写"恒真面"是越格。
- `v-finally-stuck`（`finally { inside = true; }`）＝**三枚具名用例全绿**，只红 `:1764`（红句现量 `doorRounds=1 want 2`）。
  ⇒ 台账第 127 条问的"`finally` 到底会不会执行"：模型执行**是对的**（解释器 `:655-685` 在 return 支与 throw 支都跑 finally，throw 时丢弃 finally 的返回值），但**"释放 inside"这件事只有第二次 post 那枚兄弟用例在量**。
  票 `:75` 要的三枚具名用例各只 post 一次 ⇒ 这三枚**单独不覆盖再入标志的释放**。写成"这轮把 finally 也验了"是超售。
- `v-idem-guard-gone`＝我**押错**的一发：我以为双包仍会送达一次从而全绿，实测红 `:1786`，红句是"`:1802` 到达 RPC 帧里的字符串不是页面信封逐字"——双包把帧**二次封装**了。⇒ 幂等守卫有牙，此形未破。记我：又一枚"我先验的直觉不作数"（第 122 条同族）。
- 帽数（`postMessageHopCap = 8`）：往下改（8→2、8→1）会把**正确产码**判红（8/9 枚），红句现量"re-entered 3 levels (cap 2)"⇒ 合法路径的真实嵌套深度＝3，帽数有下界凭据；往上改（更松）我**没测**——那只会让未界循环更晚引线，不改变颜色，本腿不据此写"帽数被钉死"。

## 判 4｜越格借光：**不成立**（没有任何一句把"解释器里成立"写成"WebView2 里成立"的结论），但三处**前提式**浏览器事实要具名收紧
- 结论层：新写的整块 `:1897-1914` 逐字写着它**不买**什么（`:1908-1909` "It does NOT buy 'WebView2 really binds the receiver / really lets the page replace postMessage'"），
  并把 `defineProperty` 非页面可见、descriptor 读、`configurable`、`delete` 具名放在外面（`:1910-1914`）；`impl.md §3` 那三句口径同向。⇒ "两格互相作证"这件事**没发生**。
- 前提层（本轮新写、且是浏览器事实，⛔ 本尺今天给不出凭据）：
  1. `cmd/wisp/panel_transport_35r2_test.go:1284-1286` "A browser's host method is not a detached function: **Chrome and WebView2 answer** an invocation that lost its receiver with 'TypeError: Illegal invocation'"；
  2. 同文件 `:1939` "the way **a browser's host method does**"；
  3. 同文件 `:1954`——这条在 `t.Fatalf` 的**运行时文案**里："**A real browser answers that with Illegal invocation** and the page's letter never leaves the page"。只读 `--- FAIL` 上下几行的人看不到 `:1908` 那句免责声明。
  4. 模型自己造的串 `:1290-1291` 与 `:1269`（后者是 r2 旧文，不算本轮新写）都在**模仿浏览器错误文案**，`cat` 到日志里像实测。
- 我的判：1–3 应改成"this fixture answers with…"或加"modelled after Chrome/WebView2"的限定词，⛔ 不该由解释器的绿来充当 WebView2 的行为记录。这是可读性/口径缺陷，**不是恒真面**，也不构成本格退回理由（`:75` 的 (a) 支要的两枚具名红都在，(c) 支本来就写了"⛔ 本框不许用解释器绿来满足"，所以那 7 枚未勾里这一格**不该被翻**——见下面"给编排者的一句话"）。

## 与我转述冲突的原文（硬要求那条，逐条）
1. **路径**：派单说"台件脚本 `scripts/orchestrator-countershape.sh`（腿原件 `scripts/leg-countershape.sh`）"。现量：`scripts/` 里**没有**这两枚（`ls scripts/` 无 countershape；`find` 全仓命中＝`.scratch/wisp/probes/35/r4/scripts/` 两枚）。你那句"它在仓里、可复现"成立，但**目录不是 `scripts/`**；以后派单别再写 `scripts/`（第 108 条镜像形：在 `scripts/` 找不到尺≠仓里没有）。
2. **行号**：派单说"载具 `runShippedHookInOneWorld35r4(t, nonWritable)`（`:1926`）"。现量：函数签名在 **`:1921`**，`:1926` 是函数体里的 `bf := newFakeDoc35r2()`；`:1927` 才把世界开关交给夹具、`:1928` 才是 `installPanelTransport`。其余你给的行号我逐枚复量**全对**（`:137`/`:149-160`(注释 149-153、体 154-161)/`:163-166`(拒覆写在 164-169)/`:1285-1291`/`:1946`/`:1976`/`:2026`；产码 `:673`/`:693`/`:700`）。
3. **四发反形读数**：与仓内件一致（`.scratch/wisp/probes/35/r4/logs/orch/orch-mut-summary.txt`：ma-hook 9 FAIL/3 PASS、ma-detector-gone 12/0、mb-set-gone 恰红 1、mb-typeof-gone 恰红 1），我全部自己复跑对上，⛔ 没有一条需要你改。
4. **"名册＝341 PASS/5 FAIL/1 SKIP"**：那是**整包**读数（`.scratch/wisp/probes/35/r4/logs/orch/full-package-roster.txt`）。我整包只跑到 150/346 枚就用掉 ~14 分钟（`logs/pass1-noopoverlay-run-baseline.txt` 现量）⇒ 我改跑 12 枚名册并具名理由（判不动第 2 条），⛔ 我没做整包复跑，别把我这份当整包染色凭据。
5. 你那句"M-A 是级联红……不许读成红 9 枚所以更凶"我照办，并且**补了反方向的读数**：级联归因是有凭据的（那对摘除），但 M-A 的三枚钉子本身在不可写世界上全过 ⇒ 单看 M-A 的红句不足以指认 receiver（判 2 第 4 条）。

## 我这把尺的两处缺陷（先记账，pass 1 不算任何凭据）
1. **pass 1 的十三发＝空操作 overlay 的假绿**：我把 `mk_overlay` 的"仓路径|突变路径"对形写成了只传突变路径 ⇒ JSON 变成 `{"Replace":{突变:突变}}`，两侧还都是 MSYS `/d/tmp/…`。症状是 buildrc=0、12 全绿、**连腿自己的 ma-hook 都绿**。
   抓它的是正控（那条我明知该红）。全文＋修法＝`logs/self-catch-noop-overlay.txt`；旧件保留改名 `logs/pass1-noopoverlay-*`，⛔ 未删。
2. **两把提取尺第一版是坏的**：`logs/prem-replication-table.txt` 的 awk 分段没匹配 ⇒ 三行都报 PASS=0 FAIL=0（不是读数，凭据只有 `-final` 那份）；
   第一版 final 我又用了猜的字节范围、跨界含进邻段表头，第二次按现量行号（3/31/59）重算才是本文引的数。⛔ 同一坑（第 128 条）我在这程里连踩两次。

## 没跑完的（具名，不含糊）
- 整包逐突变复跑：没做（时间不可行，见判不动第 2 条）。我做的是**引用面**尺：全仓只有 3 枚 test 文件＋1 枚 winlive 文件读 `installPanelTransport`/`panelPostMessageForwardInit`（`logs/roster-coverage.txt`）。
- "更松的帽数（8→80）仍全绿"这一发没跑，所以帽数上界我只写"无上传感"，不写"未验"之外的结论。
- 判 1 第 3 条那句"M-A 在 ma-hook 下走第一支"是从两发读数推的，没单独提行号量。
- `:63` 那三支（票的另一格）本腿未碰；我只验了 `:75` 的 (a) 支，(b) 支（M-C1/M-C2 那颗牙）我以 `v-idem-guard-gone`/`v-cap-*` 从旁边顶过，⛔ 没重跑那三支的裁决。

## 给编排者的一句话（⛔ 不是要求，不代你翻勾）
`:75` 那框的 (a) 支我判**成立**（两面各有具名红，且我复跑＋新形都没破它，只破出**新的一枚面**）；(c) 支按框里那句"⛔ 本框不许用解释器绿来满足"今天**没落地**，所以这一格该继续留在未勾。第三面（守卫的非函数分支缺世界）＋红句前缀共用"M-B RED"这两样，是下一枚写腿的活，本腿一行代码都没动。
