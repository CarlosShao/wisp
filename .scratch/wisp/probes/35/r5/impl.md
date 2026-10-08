# 35-r5 · 票 35 `AC#8` 四格交付正文（写码腿，⛔ 零产码改动）

本腿范围＝**只改一枚测试夹具＋一枚真开窗台件的头注释**。工单框一枚没勾、`-done` 没加、
`.gitignore` 与 `design/**` 里别人在飞的改动一字未动、零 push。

---

## 0. 锚

| 项 | 现量 |
|---|---|
| 起手 HEAD | `daf1f5a028f274b113db545c0776f6e409796ab9`（短号 `daf1f5a0`，分支 `dev`） |
| 起手时刻 | `2026-10-08 09:28:58 +0800`（`logs/opening-rulers.txt:1`） |
| 起手工作树 | `git status --porcelain \| wc -l` = **756**；对本腿三枚目标件（夹具／winlive／产码）= **0**（`opening-rulers.txt:7-9`） |
| 票 35 框数 | 起手 **8 未勾／2 勾**（`opening-rulers.txt:11-12`）；交件前复量仍 **8／2**（见 §6 末） |
| `AC#8` 框身 | 票 35 **`:77`**（正文 `:78-82` 四支），**不在** `:286` 那节标题里；`:286/291/292/294/296/300` 是编排者收 `35-v4` 的登记节 |
| 本腿 commit | `633c4afa`（锚件）→ `3a343bc7`（(i)(ii)(iii)）→ `099d8fe6`（(iv)）→ 本文件与票面追加为第 4 笔 |

**框面原文与我派单转述的冲突：一处，具名如下。** 派单转述 (iii) 第一处＝「`:1284-1286`（"Chrome and WebView2 answer an invocation that lost its
receiver with 'TypeError: Illegal invocation'"）」，**盘上框身 `:80` 逐字相同**（(iii) 在 `:80`、
(iv) 在 `:81`；派单说「正文在它下面几行」时框身命中的是 `:286` 那节标题行，真正的框在 **`:77`**，
四支在 **`:78`/`:79`/`:80`/`:81`，覆盖面边界在 `:82`**）。
★**唯一一处冲突具名**：派单说 (ii) 只拆「`:1976` 那枚的 4 条红句与 `:2026` 那枚的 2 条」＝6 条，
框身 `:79` 也只点名这 6 条；但在夹具里同一两枚用例共用 `M-B ` 前缀的 token 实际是 **14 枚**
——同一两枚用例里还有 `M-B READING WRONG`（2＋3 枚）、`M-B CONTROL RED`（2 枚）、
`t.Logf` 的 `M-B face has teeth`（1 枚），合计 14 枚共用 `M-B ` 前缀。框文只点名了 `M-B RED`
那一族 6 枚。我按派单后半句「要成对、要各面唯一」把**整族 14 枚**一起拆了（理由与自证见 §3(ii)）；
如果验收腿认为越界，回退办法是把 `READING WRONG`／`CONTROL RED`／`face has teeth` 那 8 枚改回
`M-B ` 前缀——但那会留下一把我明知道会串色的尺，我不主动做。

---

## 1. 落了什么（逐件＋行号）＋ ⛔ 产码零改动的自证

| 件 | 动了什么 | 现量行号 |
|---|---|---|
| `cmd/wisp/panel_transport_35r2_test.go` | (i) 新用例一枚；(ii) 前缀拆面；(iii) 三处浏览器事实口径 | 2066 → **2142** 行（＋76）；新用例 `:2074-2142`，台件那行在 `:2116`，本面红句在 `:2123` |
| `cmd/wisp/panel_transport_live_35v2_windows_test.go` | (iv) 文件头注释（**纯注释**） | 261 → **282** 行（＋21），`//go:build windows && winlive` 行未动 |
| `cmd/wisp/panel_host_windows.go` | **⛔ 零改动**（它只是 overlay 的突变对象） | blob 起手 `26b5de83…`＝HEAD＝交件时同一枚（`logs/repo-untouched.txt:3`） |
| `.scratch/wisp/probes/35/r5/**` | 锚尺、反形脚本与 5 发读数、后尺、门、本文件 | 10 枚件，零枚 `.out` |

产码零改动的三把尺（全部现跑，`logs/zero-prod-change.txt`）：
`git diff --numstat daf1f5a0 -- cmd/wisp/panel_host_windows.go` = **空输出 rc=0**（＝那枚文件零行）；
`git hash-object` 现值 = `git rev-parse daf1f5a0:` 值 = `26b5de83b93a9a141f1546dc19a9b03ffa45ede0`；
`git status --porcelain -- cmd/wisp/panel_host_windows.go` = 空。
五发 overlay 跑完之后同一把尺复量：`prod-untouched-by-overlays=YES`、`fix-untouched-by-overlays=YES`。

---

## 2. (i) 的新用例与五发反形

### 新用例

`TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsNotCallable`（`:2097`，文档注释 `:2074` 起）。

世界＝**门的名字在，值是真值非函数**：`bf.window.set(panelDispatchBinding, "not-a-function")`（`:2116`）。
台件写在 Go 侧，与 `35-r4` 的 `markNonWritable` 同族手法；⛔ 没用 `Object.defineProperty`
（它在本解释器里对页面不可见，`35-r4` 已具名，框文 `:78` 也明写「⛔ 不需要 descriptor」）。
守卫读的是 **post 时刻**的 `window.<binding>` 而不是挂勾时刻，所以在 `openDocument()` 与
`pagePost()` 之间覆写就是这一面需要的整个世界。

三枚 `M-B3 CALIBRATION RED` 先把世界钉住（`:2111`/`:2114`/`:2118`）：名字必须在（否则这发
静默退化成 `:2026` 的"门不在"世界、同一枚突变会被两枚用例重复报销）；`Bind` 桩给的必须是
function（否则"把它变成不可调用"这句没发生过）；覆写后 `jsTypeOf` 既非 `undefined` 也非
`function`。断言口径照同族既有那两枚：不抛、帽未跳、原生出口恰一度且带 receiver、信封**逐字**
转发、门零轮次且信封死在 `msgcb` 的 unbound 支（`doorRounds=0 routerRuns=0 unboundSlots=1 evalThrew=1`）。

### 五发读数（`go test -c -overlay`，突变源全在仓外 `D:/tmp/wisp35r5/`；脚本＝`scripts/countershape-35r5.sh`）

尺的分母＝`-test.run 'TestForwarding|TestPagePostMessage|TestShapeA3|TestLegacySubShape|TestUnforwardedPageEnvelope'`
在这台夹具上点到的用例枚数（旧名册 **9**／本腿新名册 **10**）。⚠ `35-v4` 报的「12 全绿」用的是
更宽的 `-run`，两个分母不许对拉，读数各自带名册。

| 发 | overlay 内容 | 色 | 读数 |
|---|---|---|---|
| **Z-compile-probe** | 产码里塞一枚故意的语法错（`fmt(((syntax error`） | **buildrc=1**，跑不了 | ★这发不是反形，是**落地正控**：Go 去 open 替换件了 ⇒ overlay 的**键路径确实匹配**。我第一版脚本把替换件文件名写成含 `cmd/wisp/` 的形状，五发全 `buildrc=1` 且报 `open D:\tmp\…pristine-cmd\wisp\panel_host_windows.go: The system cannot find the path specified` ⇒ 正是这发让我看见键路径能匹配、只是值写歪（旧件与读数保留未删） |
| **A-pristine** | 产码换回**逐字节相同**的一份（`changed-line-count=0`） | **PASS=10 FAIL=0，buildrc=0 runrc=0** | 替换生效且不改颜色。⚠ 单看这发**证明不了**匹配（同色既可能是"换成相同原件"也可能是"没换"），所以 Z 那发必须存在 |
| **B-truthy-new** | `:680` 的 `typeof window.%[1]s !== "function"` → `!window.%[1]s`（`changed-line-count=2`，逐字见 `logs/mutant-landing.txt`） | **FAIL=1**，唯一红＝`TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsNotCallable`；红句 `panel_transport_35r2_test.go:2123: M-B3 RED (guard reduced to a truthiness test): … the page's post threw: TypeError: string is not a function` | ★框文 `:78` 要求的「这一枚必须红」**成立**，且只红这一枚 |
| **C-truthy-old** | 同一枚产码突变 ＋ 夹具换回 `daf1f5a0` 那版（＝没有本腿新用例） | **PASS=9 FAIL=0，buildrc=0 runrc=0** | ★本腿自己复现了 `35-v4` 的恒绿读数：**同一枚突变，旧尺全绿** ⇒ 新用例是那一面唯一的载体，不是给旧面再补一枪 |
| **D-typeofgone-new** | `:680` → `if (inside)`（35-v2 的原 M-B），当前夹具 | **FAIL=2**：`:2054 M-B2 RED … TypeError: undefined is not a function`（门不在那枚）＋ `:2123 M-B3 RED … TypeError: string is not a function`（新那枚） | 两面**分得开**：摘掉 `typeof` 半支会把两个世界一起打坏（合理），而换成真值判断**只**打坏"真值非函数"那一面（B 发 1 红）⇒ B 与 D 的红名册不同＝两台探测器不是同一台。⚠ 但红**名册**分不开"哪种 receiver 丢失"以外的东西，只有**文案**里的 `undefined` vs `string` 可分（跟着 `35-v4` 判 1 的粒度限制一起引） |

⛔ 反形只经 `-overlay`，跟踪文件全程未写（`logs/repo-untouched.txt`）。
⚠ overlay 的两个已知死法我都撞到或避开了：①**对形**——`mk_overlay` 强制 `"orig|replacement"` 成对，
Z 发那次是**值路径写歪**（不是少传对），症状同样是"颜色不对"，但表现成 buildrc=1 而非假绿；
②**裸单反斜杠**——键与值一律 `cygpath -m` 出正斜杠。⚠ 别把②读成"必须正斜杠"：本腿没有复现
「正确转义的反斜杠也能替换成功」那一发，我只证明了我这一族的写法可行，写法先例见
`35-v4` 与本腿 Z 发的报错路径（`D:\tmp\…` 那种反斜杠形式 Go 是能 open 的）。

---

## 3. (ii)(iii)(iv) 各前后枚数尺

### (ii) 红句前缀（`logs/after-ii.txt`）

- 尺 `'M-B RED'` 枚数：改前（`daf1f5a0` blob）**6** → 改后 **0**。
- 改后各面 token：`M-B1 `=**9**（原 `:1976` 那枚＝交件行号 `:1982` 起）、`M-B2 `=**5**（原 `:2026` 那枚＝交件行号 `:2032` 起）、`M-B3 `=**9**（本腿新用例）。
- 定义：`M-B1`＝"静默拒覆写"那一面，`M-B2`＝"门不在"那一面，`M-B3`＝"门在但不可调用"那一面。三面各唯一、成对给到框文点名的两枚。
- **只改前缀**的自证：把 `M-B1 `/`M-B2 `/`M-B3 ` 折回 `M-B ` 再与 `daf1f5a0` blob 对拉（`logs/scope-diff-ii.txt`），
  差异只剩 (i) 的新用例与 (iii) 的三处文案＝**88 行**，没有第三类改动；`t.Fatalf` 枚数 **47 → 56**，
  增量恰＝新用例的 9 枚 ⇒ 既有用例的断言强度、枚数、用例名一字未动（`func Test` 枚数 8 → 9）。
- ⚠ 框文只点名 `M-B RED` 6 枚；我按"各面唯一"把同族另外 8 枚一起拆了（见 §0 的冲突上报）。

### (iii) 三处浏览器事实（`logs/after-iii.txt`）

尺＝`grep -c 'Chrome and WebView2 answer an\|the way a browser.s host method does\|A real browser answers that with Illegal invocation'`：改前 **3** → 改后 **0**。
改后逐字（短）：
1. `:1284` 起（原 3 行变 7 行，`:1286` 是承重那句）：
   `// A browser's host method is not a detached function: Chrome and WebView2 are reported` /
   `// to answer an invocation that lost its receiver with "TypeError: Illegal invocation".` /
   `// THIS RULER CANNOT PRODUCE THAT CREDENTIAL (ticket 35 AC#8(iii)) … the behaviour recorded below is THIS FIXTURE'S OWN ANSWER, modelled after Chrome/WebView2 and never offered as a WebView2 behaviour record`，末句指到 winlive 那枚文件并写明 `〔仅本机可量〕`。
2. `:1944`：`// now demands cw itself - the rule THIS FIXTURE answers with, modelled after what Chrome` ＋ `// and WebView2 are reported to do (this yard never observes a browser; AC#8(iii)) - so the`。
3. `:1960`（**在 `t.Fatalf` 的运行时文案里**）：`… their receiver. THIS FIXTURE answers that with Illegal invocation - modelled after what Chrome and WebView2 are reported to do, NOT a WebView2 behaviour record (AC#8(iii); the real-window reading is the 〔仅本机可量〕 winlive rig) - and on that model the page's letter never leaves the page, so this hook must not be certified as a delivery.`
   ⇒ 红句仍能指认用例：`M-A RED` 前缀仍在（2 枚命中），`%d`／`%s`／`%s` 参数序未动，D 发那发的红名册里这枚没被我的文案改动惊动。

### (iv) `〔仅本机可量〕` 上自己的脸（`logs/after-iv.txt`）

尺＝`grep -n '仅本机\|machine-local\|CI 永看不见' cmd/wisp/panel_transport_live_35v2_windows_test.go`：
改前（`daf1f5a0` blob）**rc=1／枚数 0**（＝票面 `AC#8(iv)` 让派单前自己复量的那一发，复现成立）→
改后 **rc=0／枚数 1**（第 3 行）。四要素逐枚命中各 1（`① WHAT IS MACHINE-LOCAL`／`② WHO OWNS IT`＝**ORCHESTRATOR**／
`③ RE-RUN CADENCE`＝每波真窗**并同波一发控制组**、件名＋`rc` 逐波进 `.scratch/wisp/probes/35/**`／
`④ THE PRICE PAID`＝CI 永不拦下"这条传输在真浏览器里坏掉"），另 `Revocation` 撤销口令 1 枚＝
「给 `:52` 的决定性读数一枚 CI 可达替身，然后这段改写、不许删」。写法形先例＝`docs/evidence/s1/33-panel-host-c27-r5.md:78`。
文件 261 → 282 行**全是注释**：`//go:build` 行、任何断言、任何 `t.Fatalf` 都没碰。

---

## 4. 门三件各 `rc`（`logs/gates.txt`）

| 门 | 载具 | `rc` |
|---|---|---|
| `go vet ./`（CWD＝`cmd/wisp`） | 件体 0 字节＝**vet 成功时本就零输出**，不是没跑；我在同一格又复跑一次确认 | **0**（两次） |
| `d22scan`（`cd tools/d22scan && go run . -root ../../`；独立 module） | 末行 `clean - no D22 ban violations`，`cmd/` 108 枚 Go 文件含注释与 `_test.go` 全扫 | **0** |
| `gofmt -l` 两枚被改件 | 列出 1 枚＝`panel_transport_35r2_test.go`，**既有漂移**：同一把尺对 `daf1f5a0` blob 也列 1 枚，hunk 是 `jsParser peek/next` 那 4 行对齐。⛔ 没顺手改无关行 | **0**（rc 与"是否列出文件"无关，这里 rc＝0） |

跑测试的机器载具（照派单 §2，⛔ 不照做就是"用例根本没跑"的假绿）：CWD＝`cmd/wisp`，
`PATH` 前置仓内 `third_party/sherpa-onnx`（现量 3 枚 dll：`onnxruntime.dll`／`sherpa-onnx-c-api.dll`／
`sherpa-onnx-cxx-api.dll`）。⛔ 本腿**没有**整包跑 `go test ./cmd/wisp/`（346 枚、~14 分钟、5 枚窗口依赖既有红），
也⛔ 没为 rc 好看放宽任何断言。

**交件名册读数（`logs/delivered-baseline.txt`，在 `099d8fe6` 之上现build现跑）**：
`-test.run 'TestForwarding|TestPagePostMessage|TestShapeA3|TestLegacySubShape|TestUnforwardedPageEnvelope|TestInbound'`
＝ **PASS=13 / FAIL=0 / SKIP=0，buildrc=0 runrc=0**。

---

## 5. 这格不许被读成什么（照 `35-r4` 的写法）

- ⛔ **不买「WebView2 真绑 receiver」**。夹具要求 `recv == cw` 是本尺**自陈的一条规则**（(iii) 已把它降级成
  `this fixture answers with…`），不是 WebView2 的行为记录。真窗那一格在票 35 `:52` 的 `-tags winlive` 凭据里。
- ⛔ **不买「WebView2 真让页面换掉 `postMessage`」**，也不买「真浏览器里 `chrome.webview.postMessage` 真不可写」。
  `markNonWritable` 是 Go 侧台件，`Object.defineProperty` 对页面仍不可见；descriptor 读、`configurable`、
  `delete`、accessor 全在 `35-r4` 具名的未建模清单上，本腿一格都没往里走。
- ⛔ **不买「`window.wispDispatch` 在真 WebView2 里可能是个非函数值」**。本腿只说：如果那种世界出现，
  守卫必须退回原生出口——这是**守卫的形状**，不是页面/宿主的事实。谁把它读成后者就是借光。
- ★ **两格⛔ 不许互相借光**：`:52`（真窗、〔仅本机可量〕）与本腿的 `AC#8(i)`（解释器里的第三个世界）
  是两枚不同的凭据。(i) 的绿**不**减少 `:52` 的欠账，(iv) 把那句话写到台件脸上也**不**等于那一发跑过——
  ⛔ 本腿没开过窗，因此本腿不声称取过任何 winlive 读数；票 111 那条 CI 步也不替它复跑（winlive 零 CI 岗位）。
- ⛔ 不许读成「`AC#8` 这一框可以勾了」：本腿只交**夹具与注释**，(i) 的"新那枚必须红"由 **B/C 两发**给到，
  但框体还挂着 `:82` 那两条覆盖面边界（再入标志的**释放**不在这三枚射程；去 `return` 是等价形；
  帽数 8 只有下界凭据，上界那发 `35-v4` 没跑）。勾不勾由 `35-v5` 裁。
- ⛔ 本腿没动 SLO 阈值／`thresholds.go`／golden；没新增 `DEFERRED(D-xx)`；`DEFERRED` 登记 1:1 未受影响。

---

## 6. 留给验收腿的问题（⛔ 我不代答）

1. **名册分母**：`35-v4` 报「12 全绿」、本腿报 C 发「9 全绿」。两把尺的 `-run` 不同形。
   你复现时**用哪份名册作准**？（我的件里名册＝`a-pristine`/`b-truthy-new` 的 10 枚，逐名可从
   `logs/run-*.txt` 里 `--- PASS` 抽。）
2. `M-B3` 的三枚 `CALIBRATION RED` 算不算**新增断言强度**？框文 `:80` 说 (ii)「⛔ 只许改前缀」，
   但那三枚在**新用例**里（属 (i)）。你认为 (i) 的载体应当只保留与同族同形的 6 枚（对照＝`:2032` 那枚的 8 枚里去掉 2 枚 calibration）、
   还是允许世界自证？（现量：新用例 9 枚 `t.Fatalf`＝1 install ＋ 3 CALIBRATION ＋ 5 面断言）
3. 我把 (ii) 从框文点名的 6 枚扩到同族 14 枚。**这是补尺还是越格**？越格的话回退点在哪一行？
4. B 发只红 `:2123`、D 发红 2 枚。⚠ **有没有第四种突变**能把"门在但不可调用"那一面变绿而
   仍留着 `typeof` 半支？（我只试了 `!window.…` 与 `if (inside)` 两支。）比如把守卫换成
   `typeof window.wispDispatch !== "function" || inside` 的**重排形**、或给 `native.call` 前加真值判断。
5. `M-B3` 的红句在 `--- FAIL` **上面一行**（现量＝文案在 `logs/run-b-truthy-new.txt:30`、`--- FAIL` 在 `:31`；Go 把 `t.Fatalf` 印在 FAIL 行之前），
   按关键词读日志的人会不会把 `M-B2`/`M-B3` 的红**错认成同一枚用例**？（两条都含
   `the forwarding hook must fall back`／`must still fall back`，用例名不同。）
6. (iii) 我把三处都降级成夹具自陈。**有没有哪一处其实是产码注释、必须按产码规矩另开一票**？
   （我判断三处都在 `_test.go` 里，尺＝`grep -n` 的行号全在夹具，见 `logs/after-iii.txt`。）
7. (iv) 的「同波一发控制组」写成了口径、没写成脚本。⛔ 我**没有**为 winlive 复跑造台件
   （那要开窗，越本腿预算）。这一格缺不缺"可执行载体"？缺的话由谁在什么预算内补？

票 35 框数复量（本文件与票面追加**前后各一次**）：追加前 **8 未勾／2 勾**（本文件 §0 那行）；
追加后见 `logs/box-census-after-append.txt`＝仍 **8 未勾／2 勾**，⛔ 一枚框没勾、没加 `-done`。
