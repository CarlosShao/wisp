# 300-v4 · `AC#6` 三问（ⓐ 形状／ⓑ 权威／ⓒ 射程）＋ 落点裁语反攻

裁决者＝`300-v4`（非实现者）。锚＝`34e1962b`。台面＝**仓外导出树**（`git archive 34e1962bb166b4494670b77296263c79bd8671a3 | tar -x -C /tmp/wisp300v4/tree`，件 `logs/00-archive-rc.txt`）。
被裁的件＝`internal/audio/wave_format_float_300_windows_test.go`（271 行，`//go:build windows`，落点笔 `9a442923`）。
本腿**零产码／零测试文件改动**（`git status --porcelain internal/audio` = **0 行**，读数在 `90-hygiene.md`）。

---

## ⓐ verdict：**⛔ 踩——`AC#6` ③ 禁的那个形状它没踩；这枚钉子有牙，且两个方向都有牙**

`AC#6` ③ 原文（票面逐字）：`新判据⛔ 能写成"常量 == 3"的字面等值再包一层同义反复——判据＝**同一枚换形攻击跑在新尺上必须红**`。
⇒ 判据本身就是一把我自己重跑的尺，不是我给的印象尺。

### 结构面（词面尺，逐字命中行＝件 `logs/20-expected-side-hits.txt`）

```
70:	mmregWaveFormatPCM uint16 = 1
74:	mmregWaveFormatIEEEFloat uint16 = 0x0003
78:	mmregWaveFormatExtensible uint16 = 0xFFFE
159:			got:       waveFormatFloat,
160:			want:      mmregWaveFormatIEEEFloat,
169:			got:       waveFormatExt,
170:			want:      mmregWaveFormatExtensible,
178:			got:       waveFormatPCM,
179:			want:      mmregWaveFormatPCM,
```

- **期望侧**＝`want:` 三枚全是 `mmreg*` 那三枚**测试自带字面量**（声明在 `:70/:74/:78`，值形 `1`／`0x0003`／`0xFFFE` 与 `mmreg.h` 那一行同形，⛔ 被"顺手化简"成 `3`——那正是 17:5x 裁语 ② 点名的形式要求）。
- **观察侧**＝`got:` 三枚是那三枚**常量本身**，⛔ 解析结果。
- 我在整份文件里逐名核对：`waveFormatFloat`／`waveFormatExt`／`waveFormatPCM` 出现在期望侧的位置**只有注释与红句文案**（`:22`、`:133-134`、`:161-164`、`:181-183`），**⛔ 一处进 `want:` 表达式**。
- 两枚行为子测试（`:246`、`:262`）的**输入字节与期望值全走测试自带字面量**，⛔ 读那三枚常量 ⇒ 它们是纯行为面，靠 `parseWaveFormat` 的生产分支给结论。

### 突变面（`AC#6` 判据两发＋我这四发反攻，全部仓外，`cmp` 逐字节还原＝件 `logs/60-restore-cmp-all.txt` 六枚 `restore_cmp=IDENTICAL`）

| 发 | 改哪一行 | 改成 | rc | 红掉的子测试（尺＝`^    --- FAIL` 子测试尺） | 红句点名 |
|---|---|---|---|---|---|
| 票面判据 N2 | `wasapi_windows.go:98` | `3`→`1` | **1** | `ieee_float_low_word`＋`the_three_constants_map_one_to_one…` | `:191`／`:214`／`:229` **三句全点名 `waveFormatFloat`／`waveFormatPCM`** |
| 票面判据 N2 | 同上 | `3`→`0` | **1** | 同上 2 枚 | `:191`／`:235` 两句点名 `waveFormatFloat` |
| pristine | — | — | 0 | 0 红（顶层 1 PASS＋子测试 6 PASS，两把尺分开数） | — |
| **我的反攻 A1** | **期望侧** `test:74` | `0x0003`→`1` | **1** | 2 枚 | `:191 waveFormatFloat = 0x0003, want 0x0001`＋`:229 waveFormatPCM = 0x0001, but mmreg.h reserves 0x0001 for WAVE_FORMAT_IEEE_FLOAT` |
| **我的反攻 A2** | **期望侧** `test:74` | `0x0003`→`0xFFFE` | **1** | 2 枚 | `:191`＋`:235`＋`:229`（点名 `waveFormatExt`） |
| 反攻 A3 | `wasapi_windows.go:99` | `waveFormatExt` `0xFFFE`→`0xFFFD` | **1** | **3 枚**（含行为面 `parseWaveFormat_expands_…`） | `:191`＋`:235`＋`:250` |
| 反攻 A4 | `wasapi_windows.go:97` | `waveFormatPCM` `1`→`2` | **1** | 2 枚 | `:191`＋`:235`（`pcm_tag` 红） |

**最要紧的一枚（ⓐ 的正身）**：A1/A2 把**期望侧字面量**改掉 ⇒ 它**⛔ 变成恒真**，而是**照红**——红句逐字 `wave_format_float_300_windows_test.go:191: waveFormatFloat = 0x0003, want 0x0001; mmreg.h:2110 WAVE_FORMAT_IEEE_FLOAT 0x0003, copied out of the SDK header, not derived from this package.`
⇒ 期望侧与观察侧**真的解耦**：两侧各自独立可坏，⛔ 存在 `:152-153` 那种"改常量同时挪两侧"的通道。这一点 `300-v3` 的原始病灶（恒等式随常量一起移动）在新尺上**结构上不可能复发**。

**红句指向的判据**：`AC#6` 原文还要求"红句具名指向那枚常量而⛔ 指向 `tag` 的等值比较"——七发红句**逐枚点名常量名**（`waveFormatFloat`／`waveFormatExt`／`waveFormatPCM`），⛔ 一枚只写 `tag = …, want …`。唯一带 `tag` 的是**行为面** `:250`，而它点名的仍是 `waveFormatExt is the constant the parser compares that tag against` ⇒ 判**满足**。

### 两点必须写下来的**限定**（⛔ 藏）

1. `:214` 那一枚"两两互不相同"的断言，**两侧都读常量**，所以它是枚**只防碰撞**的半尺：`3`→`5` 这类不碰撞的移位它⛔ 看见（A3/A4 实测它绿而 `:235` 红）。承重的是 `:235`/`:229`（对权威字面量表求匹配）。⇒ 本格⛔ 单靠 `:214`，`300-r2` 在 `10-shape.md` 里那格"两枚独立断言在 3→1 与 3→0 都红"的自述与我的读数一致。
2. 期望侧三枚字面量是**抄录**而⛔ 求值（见 ⓑ）。若三枚一起抄歪，这枚尺会绿着错——这是"抄录权威"这一形的**固有**残余，⛔ 本件能消，只能靠盘上引文核对（我这发就是那枚核对）。

**结论**：`AC#6` ①②③ 三条硬约束逐条核过——① 未碰 SLO 阈值／golden／`thresholds.go`（本腿与落地腿的写面名册皆无）；② `parse_wave_format_300_windows_test.go` **一字未动**（`git show --name-only 9a442923` 名册只有那枚新件＋三枚自家探针件，逐字在下面 ⓕ）；③ 未写成同义反复。**判：`AC#6` 的判据成立。**

---

## ⓑ verdict：**够——权威这一层是硬的，但我这枚独立证人同时量到裁语自己错了一句**

### 我现读的盘上权威（件 `logs/80-authority-mmreg-v4.txt`，`rc_sed=0`）

- 完整路径＝`C:\Program Files (x86)\Windows Kits\10\Include\10.0.26100.0\shared\mmreg.h`；**SDK 版本＝`10.0.26100.0`**（`Include` 下只有这一个版本目录）。
- 全机唯一性：`find "/c/Program Files (x86)/Windows Kits" "/c/Program Files/Windows Kits" -name mmreg.h` ⇒ **恰好 1 枚命中**（我现跑）。
- 三行逐字（`sed -n '<N>p'`，行号＝件里引的那三枚）：

```
:2110#define  WAVE_FORMAT_IEEE_FLOAT                 0x0003 /* Microsoft Corporation */
:2376#define  WAVE_FORMAT_EXTENSIBLE                 0xFFFE /* Microsoft */
:2418#define WAVE_FORMAT_PCM         1
```

⇒ **三处对上一处不差**，连`#define` 后**两枚空格**、`0x0003` 的写法、行尾注释都逐字对上；与编排者 `probes/300/orch/2026-10-10-mmreg-authority.txt`（17:58:13 那发）三行**逐字同形**。件里 `:69/:73/:77` 的注释引文＝**抄对了**。

### 权威强度的裁：`⛔` 运行时读盘这一取舍**⛔ 成立**，我用第三把尺独立核过它的前提

`scripts/portable-tests.sh` 我现读（只读，未碰）：
- `:18` 逐字 `2. any top-level --- SKIP in what actually ran is fatal;`
- `:704` 逐字 `if [ "$skipped" -ne 0 ]; then` ⇒ 未登记 SKIP **＝红**；
- `:587` 逐字 `class: fixture = the machine or OS cannot supply the subject at all`；`:594` 是 `./internal/audio/` 名下**唯一**一枚 `fixture` 登记（`TestLiveWasapiSmoke`），`:593` 证明登记可以是 `windows` 档的。
⇒ 裁语那条结构性代价（要嘛红、要嘛改 `scripts/**`＝撞 `AC#4` 名册）**⛔ 是虚的**，我信了它，因为我自己量到了同一句。
另：新件 `t.Skip` **0 处**，我的两对整包名册里 `--- SKIP` 两侧**同为 1 枚且同名 `TestLiveWasapiSmoke`**（件 `logs/70-*.txt`）⇒ 本格⛔ 往 CI 分母里塞了一枚新 SKIP。

### ★但我这枚证人同时量到：**17:5x 那句"术语更正"自己错了半句，新件把它抄成了断言级注释**

- 我这发逐字：`KSDATAFORMAT_SUBTYPE_IEEE_FLOAT` 在 **`mmreg.h` 里就有**——`:2480 #if !defined( STATIC_KSDATAFORMAT_SUBTYPE_IEEE_FLOAT )`、`:2483 DEFINE_GUIDSTRUCT("00000003-0000-0010-8000-00aa00389b71", KSDATAFORMAT_SUBTYPE_IEEE_FLOAT);`、`:2484`；同一枚 GUID 在 **`ksmedia.h:851/854/855`** 逐字同串。
- ⇒ 正确说法＝**两份头文件各带一份同 GUID 的替代定义**（两枚 `#if !defined(STATIC_…)` 守卫），而⛔「`ksmedia.h` 有、`mmreg.h` ⛔ 有」。
- 新件 `:38-41` 的注释把裁语那句照抄成 `the GUID KSDATAFORMAT_SUBTYPE_IEEE_FLOAT … lives in ksmedia.h, not in mmreg.h` ⇒ **该行文字与盘上权威冲突**；而**同包兄弟件** `parse_wave_format_300_windows_test.go:101` 引的 `KSDATAFORMAT_SUBTYPE_IEEE_FLOAT (mmreg.h:2483, Data1=0x00000003)` **是对的**（我 `:2483` 现读到同一行）⇒ 两枚件现在**互相打脸**，兄弟件对、新件的注释错。
- 后果分级：**⛔ 触及任何断言**（该件的断言只涉及三枚 WORD 值，三枚的引文我逐字对上）⇒ **⛔ 动 `AC#6` 的成立**。但它是**一件被写进测试文件的假事实**，按本仓"注释也是证据面"的口径必须销。
- 处置建议（⛔ 我动手，本腿⛔ 改测试文件）：另开一枚**只改注释**的小笔（限 `:38-41` 那四行文字，写成"同一枚 GUID 在 `mmreg.h:2480-2484` 与 `ksmedia.h:851-855` 各有一份替代定义；WORD 值那三行⛔ 是 GUID"），⛔ 碰断言、⛔ 碰 `:152-153`；并**由编排者把 17:5x 那句更正本身追加更正**（只追加⛔ 抹）。

---

## ⓒ verdict：**`waveFormatExt`／`waveFormatPCM`＝本格授权内的附带覆盖面；`F2` 那枚"声明后零使用"该另立一格，且本格让它变贵了**

- **来路清楚**：两枚⛔ 在 `AC#6` 原文里，出处＝17:5x 裁语 ③ 第 4 条（"顺带把 `waveFormatExt`、`waveFormatPCM` 一起钉……⛔ 顺手删 `waveFormatPCM`"）。落地腿照派单办＝**⛔ 越界**：派单⛔ 改契约（`C1–C32`／`D1–D47`／`R1–R9`／D43 表一字未动），它只是把同一枚判据的覆盖面顺着同一份权威文本摊开。
- **裁语**：`AC#6` 那一格的**闭合判据⛔ 依赖它们**（原文判据只点名 `waveFormatFloat`；我 N2 两发独立坐实），但**本格交付物里它们是本格的一部分**——同件、同笔、同授权链。⇒ 翻勾时请把这一层写进判语（`AC#6` 原文＝float 一枚；裁语 ③＝附带两枚；**⛔ 把两枚读成 scope creep，也⛔ 读成"原文要求"**），否则下一次缺口审计按原文逐字对时会出现"多两枚"的假差。
- **⛔ 删任何未使用声明**（本腿⛔ 删，也⛔ 建议现在删）：`F2`＝词面尺读数（`git grep` 只命中声明行 `:97`），它问的是**包面该⛔ 该留一枚没人用的常量**，那是一枚**独立射程**（改产码声明面＝又一次 AC#2 级授权）。
- ★**本格改写了 `F2` 的读数**，必须具名：`waveFormatPCM` 今天**⛔ 再也不是"零使用"**——它被本件的 `:178` 与 `:207` 读着。由此
  1. 今后任何腿重跑 `F2` 那把尺会拿到**不同答案**，`F2` 的原句只在 `9a442923` 之前的世代成立 ⇒ 台账里请标**世代锚**（与本仓"编号锚／行号锚会腐烂"同一条定式）。
  2. 那一枚未来的"删 `waveFormatPCM`"格**成本上升**：删了它，本件直接**编译不过**（`undefined: waveFormatPCM`）⇒ 那格从此必须同时动测试面。**这就是裁语 ③"顺带钉"买到的代价**，写下来比藏着值钱。
  3. 另具名一枚：本件的 `r2Authority` 是**切片而⛔ 以常量作 map 键**（`:120-123` 注释自己说了原因）⇒ 突变跑的是值、⛔ 撞编译错误。这一形我认可：用编译错误冒充断言红＝本格⛔ 要的牙。

---

## ★反攻落点裁语本身（C1／C3／C4／C5）——结论：**四枚"⛔"我全维持；⛔ 一枚推翻，但三枚的**理由**我各改一笔**

- **C1（往 `parse_wave_format_300_windows_test.go` 里加一段）＝维持 ⛔，但裁语给的理由有半句⛔ 立**。
  裁语说 C1"钉的是浮点判定⇔解析值的一致性、⛔ 值本身"——**这句只对那两枚既有行（`:152-153`）成立**；在同一个文件里加一段"常量⇔权威字面量"的断言**做得到钉值**，技术上 C1 与 C2 等价。真正⛔ 允许 C1 的是**另一条**：那两枚行是 `AC#1`/`AC#2` 的**逐字被引凭据**（票面两处直接引 `:152-153`），动同一枚文件会把两格的 prose 锚与那对成对基线一起推歪——**共享证据文件＝⛔ 动**，这才是承重那条。⇒ 结论⛔ 变，理由换。**代价对比**：C2 多一枚件（`internal/audio` 顶层用例 41→42，尺＝`^--- PASS` 顶层，我两侧各跑两发实测）；C1 省一枚件但欠两格锚。**我维持 C2。**
- **C3（测试运行时读盘上的 SDK 头）＝维持 ⛔，并补一枚裁语没列的变体，量过之后它更⛔ 该否**。
  裁语列的是 `t.Skip` 那一支（未登记 SKIP＝红，我已在 ⓑ 独立核到 `:18`/`:704` 逐字）。**还剩一形：⛔ Skip、头文件找不到就 `t.Fatalf`**。我否它，因为它只是把红的位置挪了：托管 runner 无 SDK ⇒ `./internal/audio/` 整档常红，而 `portable-tests.sh` 的红名交集尺（`AC#4` 那把）会把这一枚**记成本票新增红**——比 SKIP 更糟（SKIP 尚有一枚 `:594` 的登记先例可走，常红没有）。⇒ **C3 无论哪一支都⛔**，而"值⇔权威"的证据归探针件：本件 `logs/80-authority-mmreg-v4.txt` 就是那一格，**独立第二人、不同时刻（20:17:48）、同一份头文件、逐字同形**。
- **C4（拿仓内另一枚实现当裁判）＝维持 ⛔，而且我这发让它**比裁语更⛔ 成立**（★本腿对裁语最硬的一枚新读数）**。
  裁语的理由＝"两枚实现可以一起错"。我量到一件更直接的事：仓内那枚"另一实现"**自己就是零尺**——`internal/audio/wavinjector.go` 把同源的两枚权威值写成**行内字面量**：`:192 if fmtTag == 0xFFFE` ／ `:218 case fmtTag == 3 && bits == 32`，⛔ 引用 `waveFormatExt`／`waveFormatFloat`（`wavinjector.go` 是无 tag 的平台中立件，而那两枚常量住在 `wasapi_windows.go`＝`//go:build windows`，Linux 档⛔ 得引用）。我在仓外树各跑一发突变（件 `logs/70-50-*.txt`、`logs/70-51-*.txt`，`cp`→`sed`→跑→`cmp`＝`IDENTICAL`）：

  | 发 | 改哪一行 | rc | `^--- FAIL`＋`^    --- FAIL` |
  |---|---|---|---|
  | `:218` | `fmtTag == 3` → `fmtTag == 2` | **0** | **0 行＝整包照绿** |
  | `:192` | `fmtTag == 0xFFFE` → `0xFFFD` | **0** | **0 行＝整包照绿** |
  | pristine | — | 0 | 0 行 |

  ⇒ 若按 C4，本格的"权威"就是一枚**今天改 3→2 全仓无人喊红**的行内字面量。**用零牙的尺当牙**＝本格要治的病复发。C4 的 ⛔ 从"两枚可以一起错"升级成"**那一枚此刻正是一颗零尺的常量**"。
- **C5（把 `:378` 抽成包内谓词）＝维持 ⛔**（本格射程），细节与"这一步要另一次批准"的写法在 **ⓓ**（件 `20-residuals-and-ac4.md`）。
- 附带一枚⛔ 推翻、只算补充：**"顺带钉"值不值得**——A3 实测 `waveFormatExt` 一改就红**三处**、其中 `:250` 走的是 `parseWaveFormat` 的**生产分支**（`wasapi_windows.go:210`），⇒ 同一份件里 `waveFormatExt` 连**生产消费者**都观察得到，只有 `waveFormatFloat` 观察不到（ⓓ 的 R2）。裁语 ③ 那句"边际成本≈0"我这发看＝**低估了**，它买到的比我以为的多。
