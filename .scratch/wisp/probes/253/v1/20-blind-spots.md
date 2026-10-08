# 253-v1 · 我新找到的恒真面／裂缝（20-blind-spots）

只登记**本腿盘上种出或推理出、且不在 r1 自报清单（其 `30-blind-spots.md` 十条）里**的面。⛔ 不重复它已自报的 #1 JS 文本体不读／#3 名册外 Bind 不在射程／#6 闭包不读／#7 POSIX 分母为零／#8 读工作树不读 HEAD。

## N1 ★`Init` tied 判据的「出现≠生效」裂缝〔实测：V1 读1 绿 + V2 读2 红〕— **打破 r1 自报**

r1 在 `30-blind-spots.md` 末条自报：「`Init` 传进来的不是那枚 var 而是别的 var（尺会红在 `:398`）：⛔ 没种过」。**实测反例**：把 `w.Init` 的实参换成另一枚 var，只要那枚 var 的**任何**一枚 `fmt.Sprintf` 的**实参列表里出现过** `panelDispatchBinding` 这个 ident（哪怕该 ident 在格式串里从不生效），tied 就成立 ⇒ 盘上 rc=0 全绿，而脚本真转发的是 `window.wispNothing`（Go 从没绑过）。
机理（读尺 `absorbFile`/`censusTransport253r1` 的代码）：`varSprintf[var名]` 收集的是「该 var 名下所有 Sprintf 的 args 中**是 Ident 形**的那些」，tied 检查只问 `fed == doorBindingConst253r1` 出现过没有；**不**检查：①ident 是否被格式动词消费；②那枚 var 是不是传输真正交给 `Init` 的那段脚本；③脚本文本里有没有别的门名。V1/V2 对照（有 ident→绿 0.081s；无 ident→红 `:398`）证明开关就在「ident 是否出现过」。
**与 r1 自报 #1（M5）不是同一条**：M5 是「同一枚 var 的 JS 文本体被改」；N1 是「`Init` 实参被换成另一枚 var」——r1 恰在自报里断言换 var 会红，本腿给出了反例读数。共同底层＝ident 出现过≠ident 生效，但触发点不同（改文本 vs 换实参）。
扩展面（推理）：`varSprintf` 是**全包所有 var** 的收集表；将来任何人写下任何 var，只要某处 Sprintf 的 args 里恰巧有 `panelDispatchBinding` 一字，把 `Init` 指向它就自动 tied——攻击面随包内代码自然增长。

## N2 ★自夹具与全局名册耦合：名册键集合的任何变化都会被夹具咬住〔实测：V5a/V5b〕

`TestBindingRosterBitesItsOwnFixtures253r1` 的 4 枚夹具共用**全局** `doorRoster253r1`（`rosterRedsFor253r1` 直接读它），wantRed 按「名册＝1 名」写死。把名册加一枚 `"wispProbeRT"`（V5a）后：
- 盘上用例红 `:384`（reviewed door 无 Bind）；
- **4 枚夹具全部错位红**（每枚多出 `rostered but unbound door wispProbeRT`，want 数全不对）——即使 V5b 把产码补绑到 2 枚（盘上用例单枚 **PASS**），夹具仍 4 红。
⇒ 后果一（正面）：「名册加名」在整腿下**不是无痕动作**，夹具兜住了它。
⇒ 后果二（必须记住的读数规）：**判此尺只许整腿（两枚一起）取读数**——V5b 里盘上用例单枚是绿的；若验收只跑 `-run 'TestTransportDoorBindingMatchesRoster253r1'`，一个「名册加名＋补绑定」的提交就过了。单枚绿≠整腿绿。
⇒ 这一发同时把 r1 自报 #2 的「改名＋同批改名册=绿」（它标注推理）按整腿实测**改成不成立**：整腿读数 rc=1。

## N3 接收者不敏感：函数体内**任何** `x.Bind`/`x.Init` 都被当传输门〔推理，⛔ 没种〕

`censusTransport253r1` 对 `installPanelTransport` 体内的调用只匹配 `sel.Sel.Name == "Bind"`／`"Init"`，**不看接收者**。⇒ ①假红方向：将来在该函数体里对与传输无关的对象（如别的组件）调 `.Bind`，会被当作门名入账（可能红）；②假绿方向：若有人把真正对 `w` 的绑定写成别的调用形（如经包装函数 `bindThrough(w, …)`——selector 名不是 Bind），尺看不见。两个方向都**没种**（预算），⛔ 按「没跑不许写成立」只作推理登记。

## N4 「换成另一个常量」（题面字面）未种〔推理〕

把 `panelPostMessageForwardInit` 的 Sprintf 实参从 `panelDispatchBinding` 换成产码里另一枚**真实字符串常量**（如 `panelTitle = "Wisp panel"`，`panel_host_windows.go:96`）：`absorbFile` 收到的 ident 是 `panelTitle` ⇒ `fed == "panelDispatchBinding"` 不匹配 ⇒ tied="" ⇒ 预期红 `:398`（与 M3 同类，只是「换成字面量」vs「换成另一常量」）。本腿只种了 var 形（V1/V2），**常量形未种**，此处推理，⛔ 不当成立。

## N5 V6 的单调形（原处留空壳、无人复制）未单独种〔推理＋合成发读数限定〕

V6 合成发（复制＋空壳）的读数被 `len(sites) > 1` 支挡在前面（Fatal `:353`），「空壳独自在场」的支（空名册 Fatal `:364`，与 M2 同文案）**没拿到独立读数**。另：若攻击者「把方法整体搬到另一文件（同名同 receiver）＋原处删净」，sites=1、尺照绿——但那正是尺声明支持的「按名跨文件找」（其头注释 `:79-81` 明说 not hardcode the file），属设计，不算洞。

## 我试了但**没**打死的（⛔ 不许省）

- V3（拼接）：没打死——fail-closed Fatal `:364`。
- V5a（名册空堆名）：没打死——盘上红 `:384`＋夹具红。
- V5b（名册加名＋补绑定）：整腿**没打死**（夹具红、rc=1）；但盘上用例**单枚**被我洗绿（N2 的读数规）。
- V6（复制＋空壳）：没打死——Fatal `:353`。
- V1 读2（无 ident 的 var）：没打死——红 `:398`。
