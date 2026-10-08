# verdict.md — 33-v4（非实现者验收腿，攻写腿 `33-r10`／票 33 `AC#13`）

被审三枚提交（我自己现量，⛔ 不是转述）：

```
44d178c946ca8c21ad35211ef5d93789d96edfe5 | 2026-10-08 10:08:42 +0800 | 33-r10 anchor: AC#13 three-cell re-measurement ...
70b00885683457956c253713641bfbb9723a97f0 | 2026-10-08 10:16:41 +0800 | test(33-r10 AC#13 item 2): give the cold-start page handover a window-free content ruler
db6ef1b1912f2b7d7b4a1518d6f2d98b472f44b9 | 2026-10-08 10:23:18 +0800 | docs(33-r10): file the AC#13 evidence set and the ticket addendum naming the expired premise
```

起手 HEAD `165b45c4fab20e59c91eff58968f41e7a68708c3`；交件时 HEAD 已被并行腿推到
`dcbb93d22ffc403d7a1f75b23676f25fa92897d8`（尺与逐枚哈希对拉在 `logs/99-hash-table.txt`：
我跑完之后，被审的四枚跟踪件工作树哈希 **仍等于** 起手时记下的哈希 ⇒ 本程一字未改跟踪件）。

一句话总判：**三格判据里 ① 与 ③ 已被实测钉住，② 只在带 bundle 的机器上有主体；产码那一改是干净的逐字搬；
它没新增红；框不该翻勾。** 下面逐格。

---

## Q1｜"逐字搬"这句话成立吗 —— **成立**（语句级；⚠ 锁那一问的判定在末段，它不是 33-r10 的账）

尺与读数：`logs/29-q1-statement-differ.txt`、`logs/04-before-statements.txt`、`logs/05-after-statements.txt`、
`logs/01-product-diff-raw.txt`、`logs/02-before-bringup.txt`、`logs/03-after-bringup-and-helper.txt`。

1. **RULER A（把搬动的那一段逐字节对拉）**

   ```
   diff <(git show 70b00885^:cmd/wisp/panel_host_windows.go | awk "NR>=425 && NR<=429") \
        <(git show 70b00885:cmd/wisp/panel_host_windows.go  | awk "NR>=448 && NR<=452")
   rcA=0
   ```

   ⇒ 搬动前后的那 5 行（`rtMs := m.firstRoundTripLocked(ctx, t0)` / 空行 /
   `if err := m.serveEntry(); err != nil {` / `m.serveNotBuiltNoticeLocked()` / `}`）**逐字节相同**。

2. **RULER B（语句序列作差：改前 `bringUp` 全函数 vs 改后 `bringUp`＋`coldStartPageHandover` 拼起来，
   剥掉注释行与空行）** → `rcB=1`，差异**只有三处**：

   ```
   32,35c32   改前四句  ->  改后一行 rtMs := m.coldStartPageHandover(ctx, t0)
   39a37,43   新增：} / func ... coldStartPageHandover(...) float64 { / 那四句原样 / return rtMs / }
   ```

   ⇒ **语句零增、零删、零换序**；唯一新增的可执行语句是 `return rtMs`。
   非注释的新增行全列在尺里（`git show 70b00885 -- <file> | grep -E '^\+' | grep -vE '^\+\s*//'`），
   逐枚就是上面那些。

3. **`return rtMs` 的值语义／赋值链**：`bringUp` 侧 `m.mu.Lock(); m.lastColdMs = rtMs; m.mu.Unlock()`
   三句**一行未动**（`:423-425`，`git blame` 仍归 `697b4faeb`，见 `logs/08-blame-lock-bracket.txt`）。
   `rtMs` 以值过函数边界，中间无任何改写；`firstRoundTripLocked` 的四条 `-1` 出口
   （`:755-757` 无控件／`:770-772` 绑定失败／`:787-791` ctx 取消／`:792-794` 5 s 超时）今天仍然**照原样**
   穿过新函数写进 `m.lastColdMs`——**包括"探测被拒的那次冷启动仍然继续供页"这一支**（改前改后都如此，
   我没找到任何一处把它改了）。⇒ 值语义未变。
4. **额外盯的锁**：调用点 `:421` 落在 `m.mu.Unlock()`（`:417`）与 `m.mu.Lock()`（`:423`）**之间** ⇒ 跨调用不持锁，
   新函数注释那句 "holds no lock across the call" **成立**。但被调的三枚里两枚带 `Locked` 后缀，
   而按本仓自己的形制（`setPriorFocusLocked`，`:535` 明写 "The caller must hold m.mu"，唯一调用点 `:514` 在
   `:511-515` 的持锁段里＝形制被守），`Locked` 后缀**意味着调用点应持锁**。`firstRoundTripLocked`（`:751`）与
   `serveNotBuiltNoticeLocked`（`:460`）**自己**取锁（`:752`／`:461`），`sync.Mutex` 不可重入 ⇒
   如果谁照名字在持锁状态下调它们，就是**自死锁**。
   **判定＝既有形制，不是这枚改动引入的**，凭据（`logs/06-blame-before-callsite.txt`）：
   改前同一对调用行的 blame 是 `13acad460`（`rtMs := m.firstRoundTripLocked(ctx, t0)` 与
   `m.serveNotBuiltNoticeLocked()`，2026-10-01 13:52:48）与 `697b4faeb`（`if err := m.serveEntry()`，
   2026-10-01 10:10:29）；锁括号 `:413-417`／`:423-425` 也全是 `697b4faeb`。
   ⇒ **具名写成"与本程无关的旧账"**：命名与真实锁纪律相反，载体是 10-01 那两刀，⛔ 不记在 33-r10 账上，
   我**也没顺手修**（改它＝动锁形状，超出验收腿射程）。
   这格我只记一笔**新增**的小账：新函数自己没有 `Locked` 后缀、也不取锁 ⇒ 它把"这两枚被调者不许持锁调用"
   这件事**原地搬了一遍**，没改善也没恶化。

**Q1 判语：成立。** 附一条不属于 33-r10 的旧账（`*Locked` 命名与锁纪律相反，`13acad46`/`697b4fa` 的）。

---

## Q2｜那三句"最终文档"断言到底有没有牙 —— **有牙，三问各自独立咬得住**；另有一枚具名盲区

⛔ 我没有拿写腿自报的那四发当凭据（只在 `logs/32-audited-leg-filed-readings.txt` 里存档它们的色，
供名册对拉）。下面六发**全是我自己造的、它没试过的坏法**，载具＝`go test -c -overlay` ＋仓外拷贝
`D:/tmp/wisp33v4/`（⛔ 未改任何跟踪件；`go test -c` 只写我指定的仓外 `-o` 路径，`go build` 中间产物进 GOCACHE）。

分母尺写法（⛔ 不是 `grep -c` 的行数）：每发的尺都是
`<exe> -test.run '^TestAC13ColdStartPageOverEndsOnEntryContentNotTheProbe$' -test.v -test.count=1`，
`-test.v` 下名册恒为 **1 枚用例**（`logs/v4-mut-summary.txt` 头部记着这把尺原文）。
本树世界读数（我**自己**跑出来的 pristine）：`probe document is 136 byte(s)`／
`embed resolves 1044 entry byte(s), 1 element id(s) [root]`／`--- PASS`，`runrc=0` ⇒ **本树 `entryOK=true`**。

| 发次 | 我造的坏法（目标｜替换成对给，overlay JSON 原文落件） | 落地自证 | 色 | 响在哪一问 |
|---|---|---|---|---|
| `m1` | `serveEntry` 交给页面的字节换成**探测页那一份**（`w.SetHtml(string(data))` → `_ = data` ＋ 探测页字面量） | snippet 命中 1 次／替换件 sha16≠目标／exe sha16 `d0849d95c9000473` ≠ pristine `73d9306ac9499591` | **红** | 三问全响：`:261`＋`:264`＋`:268`＋`:277` |
| `m2` | 换向之后**只让明示页当最后一份**：`if err := m.serveEntry(); err != nil { … }` → `m.serveEntry()` ＋ 无条件 `m.serveNotBuiltNoticeLocked()` | 同上，exe `add890e798913219` | **红** | **只**响 `:268`＋`:277`（第 3 问）——第 1、2 问在这发里按设计是绿的 ⇒ 第 3 问**独立有牙** |
| `m3` | 把探测页 markup 改成**可看的**（补 `<p>loading the panel</p>`），**次序不动** | exe `053b81dad1b297a7` | **绿** | 内容锚自证：身份比对用的是运行时取的那一份，改字面量钉不死它；同时证明"第 2 问单独不够"（与 impl.md §5.5 口径一致） |
| `m4` | **打在内容判据本身**：入口的 `id="root"` 改名（交给页面一枚手写面板形状页）＋ 把台件第 3 问的"字节本体"那半枚**换成常量判据**（`if !strings.Contains(last, string(entryBytes)) {` → `if false {`） | 两枚 Replace 成对给（产码件＋台件）；`:268` 在输出里**零命中**＝那半枚确实被关掉了，不是空操作；exe `80fbca11c022320e` | **红** | 只响 `:277` ⇒ **`id` 那一问单独也有牙** |
| `m5` | **CI／新鲜检出的形状**：overlay 把 `frontend/dist/index.html` 映射成**空值＝该文件不存在**（⛔ 未动盘上那枚未跟踪件），再把明示页正文换成**只有 `<script>`** | 落地自证是台件自己的世界读数变了：`embed reports Built()=false (not-built branch)` ＋ `:282` 那句 "no subject here"；exe `3cb4fe47aed68abc` | **红** | 只响 `:264` ⇒ 在第 3 问**没有主体**的那半边世界，第 2 问仍单独咬得住 |
| `m6` | **具名盲区**：把探测从尾段里**删掉**（`rtMs := m.firstRoundTripLocked(ctx, t0)` → `rtMs := float64(-1)`） | exe sha16 不同、build rc=0 | **绿** | 谁都没响 |

`m6` 这发要说清：票面 `AC#13` 的完成判据里明写着"⛔ 不许用删掉探测来糊"，而这枚无窗尺**看不见探测被删**
（它只问最后一份文档的内容，`rtMs` 不参与任何判定——这本身是它的题中之义）。那一格的牙在**别处**：
`cmd/wisp/panel_host_windows_test.go:661-663`（`coldMs <= 0` → `t.Fatalf`）——⚠ 那枚用例带真窗
（在 `logs/25-window-dependent-roster.txt` 的名册里），我按约束**没跑**。
⇒ 判语：**不许把这枚无窗绿读成"探测还在"**；这条不是 33-r10 的缺陷（尺的范围就是最终文档），是**名册的分工**，
要登记就得登记成"这一格由真窗族 `:661` 守，本波没跑"。

⚠ 读 `logs/v4-mut-summary.txt` 的"响在哪一问"那一列时注意：表里我给的是**参与判定的 `t.Errorf` 行**
（`:261`/`:264`/`:268`/`:277`）。同一把 grep 也会带出 `:246`／`:279`／`:282` 三行，那是这枚用例自己的
`t.Logf` 读数（探测页字节数／世界读数／最后一份文档摘要），⛔ 不是断言，不许拿它们当"响了一问"。

另两枚我量过的"它没说谎"的形状检查（`logs/99-hash-table.txt` 末段）：
台件里 `t.Errorf/t.Fatalf` 参杂调用计数（`len(c.docs)` 等）＝**0 处**，参杂行号锚（`:\d{2,4}`／`line N`）＝**0 处**
⇒ "不问调用枚数、不用行号"两句**属实**。（`docs` 切片存在，但只作为"最后一份"的来源，不参与判定。）

**Q2 判语：成立（有牙）**，五红一绿，绿的那一发是**我造的盲区**、不是它藏起来的假绿；
唯一必须写的限制是 `m5` 那类世界的**另一半**——没有 bundle 时第 3 问无主体（impl.md §5.4 自己具名欠了这格）。

---

## Q3｜去掉"②今天成立"那句反证 —— **字面不成立；带范围才成立**

负向尺（⛔ 不是裸 grep 符号名，按**调用／断言形状**搜；件＝`logs/30-q3-consumer-census.txt`）：

```
git grep -nE 'lastDoc|firstDoc|\.docs\[|getElementById|SetHtml\(string\(data\)\)' -- '*_test.go' cmd internal
git grep -nE 'func \([a-z] \*[A-Za-z0-9_]+\) SetHtml\(' -- '*_test.go'      # 谁在记录"交给控件的文档"
git grep -nE '^func Test.*(ColdStart|PageOver|Handover|EndsOn).*\(' -- '*_test.go'
```

读数：全仓**只有两枚**用例在问"最终那份文档是什么"：

1. `cmd/wisp/panel_pageover_33r10_windows_test.go:241`（本次新件，无窗）；
2. `cmd/wisp/panel_resident_windows_test.go:314`
   `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` —— 断言在 `:330`：
   "the live document contains NONE of the %d element ids the embedded entry declares"。
   它的 `id` 走的是**同一条接缝**（`:265-291 entryIDProbes` → `panel.BuiltinAssets()` ＋
   `Resolve(panel.EntryFile)` ＋ `entryIDRe`）。

⇒ **"没有仪器钉着最终文档含真内容"这句话在本机不成立**：`:314` 那枚今天就有主体（`entryIDProbes` 在
本树解析出 1 枚 `root`，我的 pristine 世界读数 `1 element id(s) [root]` 量的是同一枚入口的同一条接缝），
它跳不跳？不跳——它 `t.Skipf` 的条件是 `len(probes)==0`（`:316-318`）。我**没跑它**（它买真窗，约束 3 禁）。
唯一仍然成立的窄读是**范围限定在 CI／新鲜检出**：`frontend/dist/*` 被 `frontend/.gitignore:12`（根
`.gitignore:22-23` 同样口径）忽略，跟踪件只有 `.gitkeep`；而 CI 里唯一跑 `npm run build` 的 job 是
`lint-frontend`（`.github/workflows/ci.yml:953-954`），Go 测试的 `test-windows`（`:505`）是**另一个 job**、
没有 bundle ⇒ 在 CI 那半边，两枚尺的第 3 问**同时**没有主体。

票 35 那族（`Forwarding` 五枚，`panel_transport_35r2_test.go:1791/1952/1982/2032/2097`）拦的是**消息**
（`window.postMessage` 转发），不读最终文档；33-r9 的出口判据（`panel_resident_windows_test.go:646`）
买真窗、判的是线程退净。⇒ 都不是这一格的既有仪器。

**Q3 判语：不成立（按字面）／部分成立（按 CI 范围）。** 具名说：**②已被
`cmd/wisp/panel_resident_windows_test.go:314`（默认 `windows` 档，非 winlive）钉着**，
33-r10 补的是"不开窗也能跑"那一面——这句话它自己在票 33 追加节 ② 与 `impl.md` §5.4 里写的就是范围版，
写成了"新鲜检出／CI 里这一格零仪器"；**是转述把它扩成了全称句**（见末段顶回清单第 2 条）。

---

## Q4｜"前提已过期"这个判语对不对 —— **对，框里那句作为产码事实已过期**（⛔ 但翻勾不该跟着这个判语走）

我自己复跑，⛔ 没抄它的结论：

```
$ git log -1 --format="%H | %ad | %s" 13acad46
13acad460f4d163914cfc8593ae03e40336dee02 | 2026-10-01 13:52:48 +0800 | feat(panel-host): 33-r5 land the panel in the resident process (ticket 33 items 1-6)
```

- **换向前**（`13acad46^`，件 `logs/14-order-before-13acad46.txt`）：`bringUp` 尾段＝
  `:221 if err := m.serveEntry()` → `:227 rtMs := m.firstRoundTripLocked(ctx, t0)`
  ⇒ **探测在最后**，与票面 `:39` 那句行号（`:221`→`:227`）逐枚对得上，"最终显示的是探测页"当时为真。
- **换向后**（`13acad46`，件 `logs/15-order-after-13acad46.txt`＋`logs/16-blame-flip.txt`）：
  `:268 rtMs := m.firstRoundTripLocked(ctx, t0)` → `:270 if err := m.serveEntry()`，
  且 `:268`/`:271` 两行的 blame 就是 `13acad460` 自己。
- **今天**（HEAD，件 `logs/03-after-bringup-and-helper.txt`）：`:448` 探测 → `:450` 供页 → 失败支 `:451`
  明示页（`:467-468`，宿主自己写的 `panel assets unavailable` 页）。
  ⇒ **成功支最终文档＝入口页字节；失败支最终文档＝明示页；没有任何一条路径停在探测页上。**
- 顺带把 `git log -S` 那把尺也交了（`logs/13-log-S-probe-call.txt`）：
  `rtMs := m.firstRoundTripLocked(ctx, t0)` 这枚语句只在 `697b4faeb` 计入——**这恰是 `-S` 会漏的**
  （同一串文本在两个位置各算一次，净数不变），所以我用**逐版本语句 dump ＋ blame** 定案，⛔ 不靠 `-S`。

**判语（翻勾归编排者，我一格框没动）**：票 33 `AC#13` 这框**现在不该翻**。理由不是"前提没过期的那半边"，
而是它自己的三条完成判据里 **② 尚未在 CI 可复现地闭合**：本波补的无窗尺在带 bundle 的树上有牙（我的 `m1/m2/m4`），
在新鲜检出里只剩第 1、2 问（`m5` 就是那一半世界），而票面 ② 要的是"最终文档里含 embed 入口的真内容"。
⇒ **该翻的是"框里那句事实描述"（已由 `db6ef1b1` 具名追加更正，尺：框枚数前后 `13 未勾／1 已勾`，
`git diff 44d178c9^..db6ef1b1 -- <ticket>` 里 `- [ ]`/`- [x]` 行零增删，见 `logs/31-*`）；
不该翻的是勾。** 要闭合 ② 的 CI 面，出路＝把 `frontend/dist` 的入口由构建闸门喂给 Go 测试
（票 274 那条链），⛔ 不是本腿或本程能改 `ci.yml` 拧的。

---

## Q5｜无窗绿与真窗绿之间不许互相冒充 —— **两枚都不许替对方交货**，口径写死如下

现量的两枚尺各自问的是什么：

- **无窗尺**（本波新件）：`m.coldStartPageHandover(...)` 跑完后，读**记录型 sink** 的最后一份文档
  （`:254-255` → `:261`／`:264`／`:268`／`:277`）。它对"控件有没有真的渲染、页面自己的脚本有没有把 `#root`
  填出来、资源过滤器能不能解析 `./assets/*`"**一个字都没问**；而且 sink 自己**建模**了"新文档一显示，
  页面就当场调已绑定的无参回调"（`:72-85`）——这是建模，⛔ 不是测量。
- **真窗尺**（票 35，编排者本人 10-08 在 `2dcef1a6` 跑过，件 `.scratch/wisp/probes/35/v6/live-wave.md`）：
  它断的是 `:244` "chrome.webview.postMessage from the page did not reach the Go door"、`:249` 信封逐字段、
  `:272` 泵返回、`:279` 进程树零残留——**问的是消息到没到 Go 的门，⛔ 从没问最终那份文档是什么**。
  同一族里唯一问最终文档的那枚是 `panel_resident_windows_test.go:314`，本波**谁都没跑**（写腿 `impl.md` §5.2
  明写没跑，我这腿按约束 3 也不跑）。

**口径（写死，两句都不许省）**：

1. "无窗绿" 只能支撑这一句：** shipped 尾段交给控件的最后一份文档，字节上等于宿主从 embed 解析出的入口，
   且剥掉 `<script>` 后仍有可看内容。** 它⛔ 不能写成"真浏览器里用户最终看到的是入口页"。
2. "票 35 真窗绿" 只能支撑"页面发的信封确实到了 Go 的门"。它⛔ 不能替 AC#13 交货，
   因为它跑的是另一枚命题，而且它 `t.Logf` 那行冷启动时延本波已知是**哨兵 `-1` 未替换**的坏读数
   （`live-wave.md:44`），⛔ 不许当凭据引用。
3. 要把 AC#13 的"用户看得见面板"写进任何句子，**唯一合法凭据是 `panel_resident_windows_test.go:314` 那一发真窗**
   （带 bundle 的树），本波**没有**这一发的色。

**读错了最坏会放行什么形状的真 bug**（这正是我造 `m6`、也是本格存在的理由）：

- 产码把**最后一份交给控件的文档**换成一枚"含入口字节的壳"（我的 `m1`/`m4` 形状），或把探测步骤整个删掉
  （`m6`，绿）——无窗尺的后两问仍可能绿；只有真窗那枚会响。
- 更狠的一类**两枚无窗仪器永远看不见**：字节交对了，但**控件没渲染**——CSP `default-src 'none';
  script-src 'self'`（`frontend/dist/index.html` 里就是这一条）＋ `SetHtml` 喂进去的文档没有可解析的
  `./assets/*` 源 ⇒ `#root` 恒空。这时无窗尺全绿（第 3 问只要求"含入口字节与 `id`"，而 `id="root"` 确实在），
  票 35 真窗尺也全绿（它测的是 postMessage 通道）。**⇒ 最坏形状＝"面板开出来了，但里面是空的"，
  而两枚绿会把它一起放行。** 这一格今天的唯一牙就是那枚没跑的真窗 AC#13 用例。

---

## Q6｜卫生与不新增红 —— **两枚被审件干净；包级三枚既有未格式化件归别的腿；名册零收缩**

件：`logs/20-vet.txt`、`logs/21-d22scan.txt`、`logs/22-gofmt-list.txt`、`logs/23-crlf-flagged.txt`、
`logs/24-gofmt-headblobs.txt`、`logs/26-hygiene-details.txt`、`logs/25-window-dependent-roster.txt`、
`logs/19f-roster-prechange.txt`、`logs/19g-roster-HEAD.txt`、`logs/19d/19e/27/28*.colours.txt`、
`logs/34-prior-windowed-reds.txt`。

| 尺 | 命令原文 | rc | 读数 |
|---|---|---|---|
| vet | `go vet ./cmd/wisp/` | 0 | 零输出 |
| vet（winlive 档，防我新增件把 CI 那步编译弄坏） | `go vet -tags winlive ./cmd/wisp/` | 0 | 零输出 |
| d22scan | `sh scripts/d22scan.sh` | 0 | 末行 `d22scan: clean - no D22 ban violations`；分母现量 `bans #1-5 internal/=228 cmd/=38`、`ban #6 frontend/=85`、`ban #7 internal/tools/=23`、`ban #8 design/=39 frontend/=85 internal/=514 cmd/=109` |
| gofmt（被审两枚） | `gofmt -l cmd/wisp/panel_host_windows.go cmd/wisp/panel_pageover_33r10_windows_test.go` | 0 | **零行**＝两枚都格式化 |
| gofmt（整包） | `gofmt -l cmd/wisp/` | 0 | 三枚：`models.go`／`panel_inbound_guards_35r3_test.go`／`panel_transport_35r2_test.go` |
| CRLF 幻影尺 | 对拉 `tr -cd '\r' \| wc -c`（工作树 vs `git show HEAD:<file>`） | 0 | `models.go worktreeCR=334 headCR=0` ⇒ **幻影**（CR 剥掉后 `gofmt` 不响）；另两枚 `worktreeCR=0 headCR=0` ⇒ **真未格式化**；被审两枚 `0/0` |
| 真未格式化两枚的归属 | `gofmt -l` 跑在 CR 剥掉的 HEAD blob 上仍列出；`git log -1 -- <file>` | 0 | `panel_inbound_guards_35r3_test.go` 最后一刀 `7f9d6e40`(10-07)、`panel_transport_35r2_test.go` `3a343bc7`(10-08 票 35) ⇒ ⛔ 不记 33-r10 账，我**一枚未改** |
| 名册不收缩 | 两支 exe 各 `-test.list='.*'`，**先剥掉汇总／计数行再作差**（`FAIL: <名>` 与计数行同形那条坑） | 0 | 改前 251 枚、改后 252 枚；差集**恰好一枚**＝`TestAC13ColdStartPageOverEndsOnEntryContentNotTheProbe`；⛔ 零枚消失 |
| 不新增红（无窗面） | 12 枚点名集（新件＋写腿自证的 5 枚邻居＋票 35 族 5 枚＋255 一枚），两支 exe 各跑一遍 | 0 | 改前 `PASS=11 FAIL=0 SKIP=0`，改后 `PASS=12 FAIL=0 SKIP=0`；色名册作差只多**一行 PASS** ⇒ 33-r10 在无窗面**零新增红** |

**既有窗口依赖族（带着交件并具名，⛔ 我按约束 3 一枚都没跑，也不判它们今天的色）**：
尺＝逐枚扫函数体里的开窗调用（`startPanelForTest|webview2.New|bringUp|.Show(|HotShow(|BringUp(`），
件 `logs/25-window-dependent-roster.txt`，**14 枚**：

```
windows 档（12）：TestAC13BringUpRefusesAThreadWithAQueuedClose  TestAC13BringUpSurvivesAReusedThreadQuit
  TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe  TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever
  TestAC14AwaitedBindingReplyReachesThePage  TestAC14GoSideEvalPushReachesThePage
  TestAC4PriorFocusSurvivesARefusedPanelSample  TestBallPanelGesturesReachThePanelThread
  TestPanelThreadIsSTAAndExitsCleanly  TestTicket255PanelHostBuildsItsWindowOptions
  TestPanelHostRealWindowHopAndLifecycle (panel_host_windows_test.go:614)
  TestAC4FocusReturnToPriorWindowGap33r5 (panel_host_windows_test.go)
winlive 档（2）：TestLive35v2PagePostMessageEnvelopeReachesTheGoDoor  TestTicket255RealWindowWidthFollowsTheConfig
```

它们里已记录过的红（⛔ 我**没有重跑**，只把出处名册带回来）：`33-r9` 的整包件
`.scratch/wisp/probes/33/r9/09-fullpack-1.txt` 里 `--- FAIL: TestAC4FocusReturnToPriorWindowGap33r5`；
`33-r6` 的 `.scratch/wisp/probes/33/r6/fullpack2.log` 里 `--- FAIL: TestBallPanelGesturesReachThePanelThread`。
⇒ **这些不记在 33-r10 账上**（它一枚没动），⛔ 我也没拿"放宽任何断言"去换绿。

同一批 10-01 的整包件里还有**一枚与开窗无关**的既有红，单独具名：
`--- FAIL: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`（`cmd/wisp/config_reload_223_test.go:535`，
无 build tag；红因是 `:445` 那句"the card does not name risk.permission_mode"＝票 223 的 L2 卡面，⛔ 不是窗口族）。
我按尺把它从窗口名册里摘出来：它在我那把 `webview2.New|startPanelForTest|bringUp|.Show(|HotShow(` 的体扫里
**命中 0 次**，所以不属上面那 14 枚；它的色是 10-01 的读数，本程同样**没重跑**。
⇒ 本腿**没有**交整包色：那要开真窗，与本程硬约束 3 冲突——这一格是**具名的射程限制**，⛔ 不是"整包全绿"。

---

## 我具名顶回的转述（约束 6：以盘上原文为准）

1. **"起手 HEAD 不早于 `2dcef1a6`"** —— 现量起手 HEAD＝`165b45c4…`（`180-c1` 那刀，2026-10-08 10:26:13），
   我在锚件里按**现量＋不早于**记，⛔ 不写成等于。交件时 HEAD 又漂到 `dcbb93d2…`；被审件哈希逐枚对拉未变
   （`logs/99-hash-table.txt`），这正是不写"应为 <sha>"的理由。
2. **"Q3 里那句『含 winlive 那枚 `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`』"** —— 盘上原文：
   那枚用例在 `cmd/wisp/panel_resident_windows_test.go`，该文件**只有 `//go:build windows`**（第 1 行），
   ⛔ 不带 `winlive`；winlive 档的是 `panel_host_windows_live_test.go` / `*_live_*` 那几枚。
   范围我按原文改判：它是**默认 windows 档的真窗尺**（见 Q3）。
3. **"它自报的反控与正控：`swap` 红／`pristine` 绿／`shift` 绿／`widen_swap` 红"** —— 与盘上
   `.scratch/wisp/probes/33/r10/mut/mut-summary.log` 逐枚对得上（存档在 `logs/32-*.txt`）；
   但 `impl.md` §3 那张表把 `shift` 标成"行号全位移"而**没**说清 `pristine` 那发的 `diff=0`，
   我只把**我自己的**六发当凭据。
4. **"②今天成立"（无仪器钉着最终文档含真内容）** —— 我按全称句**判它不成立**（`panel_resident_windows_test.go:314`
   在本树有主体、不跳）；写腿盘上原文本来就是范围句（"新鲜检出／CI 里这一格零仪器"），
   **是转述把它去掉了范围**。两处都记，⛔ 我按原文判。
5. **"三笔 commit… 证据＋票面追加"** —— `git show --stat db6ef1b1` 现量：票 33 `+39 行／-0 行`，
   另 14 枚证据件全新增；框枚数尺 `13 未勾／1 已勾` 在 `44d178c9^`／`db6ef1b1`／HEAD **三处同数**，
   且 `- [ ]`/`- [x]` 行在整段区间作差里零增删 ⇒ 追加属实，⛔ 没人翻框。
6. **产码改动"只动 `panel_host_windows.go`，41 行（＋30／－11）"** —— `git show --numstat` 现量
   `30 11 cmd/wisp/panel_host_windows.go`（41＝两侧和），⛔ 不是"41 行新增"。形状＝逐字搬，见 Q1。

---

## 本程硬约束自证

- ⛔ 未改任何跟踪件：起手记下的三枚被审件哈希（`00-anchor.md`）与交件时逐枚 `git hash-object`
  对拉**全等**，票 33 与 HEAD blob 亦全等（`logs/99-hash-table.txt`）。写点只有 `.scratch/wisp/probes/33/v4/**`
  ＋仓外 `D:/tmp/wisp33v4/**`。
- ⛔ 未 push、未 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；两笔 commit 都用**显式 pathspec**，
  未用 `git add -A`/`.`。在飞件（`design/`、`.gitignore`、`.scratch/wisp/probes/161/**`、`?? -`）一字未动。
- ⛔ 未开真窗（约束 3）；突变全走 `-overlay` ＋仓外拷贝，overlay 每发都是**成对** `"目标": "替换件"`，
  `m5` 那一发的"空值＝文件不存在"只作用于编译器视图，⛔ 盘上那枚未跟踪的 `frontend/dist/index.html`
  交件时仍是 `21e42f24…`（`logs/99-hash-table.txt` 里带着它）。
- 证据件零 `.out`、零 0 字节（`find -size 0`＝空，`find -name '*.out'`＝空），只用 `.md`/`.txt`；
  每把自己落一行 `rc=N`。
- ⛔ 一格框未动，⛔ 一行台账未动，⛔ 一枚跟踪件未顺手修（包括 Q1 那枚旧账与 Q6 那三枚未格式化件）。
