# 306-v1 · 10 反形敏感度／AC#2b 的形／AC#1↔AC#2 是否同一发（必答 1·3·4）

台面＝仓外导出树 `C:/Users/swq/tmp/306v1-mut-20261011-075658`（尺＝`git archive e4740e35… | tar -x -C <树>`，
`rc_archive_tar=0`；件＝`logs/TREE-PATH.txt`）。⛔ 共享工作树做任何突变。
驱动＝`logs/sens-driver.sh`（逐枚突变，读回 `logs/sens/ZZ-summary.txt`）＋`logs/cov-roster-driver.sh`
（覆盖块图＋整包名册，`logs/sens/ZZ-cov-summary.txt`）＋`logs/controls-driver.sh`（无夹具对照，`logs/ZZ-controls.txt`）。
定向跑的每一发都带 `-vet=off`（因为变体形里带 `if false && …` 这类刻意恒真的形状，⛔ 让 vet 替我读）；
正控：真件⛔ 变体时 `GOFLAGS= go vet ./internal/audio/` 工作树 `rc=0`（`logs/g1-vet-worktree.txt`）、
同树默认 vet 的定向跑 `rc=0`（`logs/sens-smoke.txt`，两枚指名用例逐字 `--- PASS`）。

## 0. 三枚权威值在盘上的形状（本腿自己取，⛔ 抄腿）

尺＝`sed -n '192p;196p;218p' logs/00-blob-wavinjector.go`（`logs/00-blobcmp.txt`），
且 `git show HEAD:internal/audio/wavinjector.go | cmp - internal/audio/wavinjector.go`＝**`rc_blob_vs_worktree_cmp=0`**
⇒ blob＝工作树，下面三枚行号在两把尺上逐字相同：

```
:192  if fmtTag == 0xFFFE { // WAVE_FORMAT_EXTENSIBLE: real tag is SubFormat[0:2]
:196  fmtTag = binary.LittleEndian.Uint16(data[body+24 : body+26])
:218  case fmtTag == 3 && bits == 32:
```

## 1. ★必答 1：反形敏感度（三形并排，⛔ 一票结论）

新面＝`internal/audio/wavinjector_extensible_float_306_test.go`，两枚指名用例
`TestParseWavExtensibleFloat32306`（同包直调 `parseWav`）＋`TestWavInjectorExtensibleFloat32306`（真 C8 入口 `NewWavInjector`）。

### 1.1 形一＝把产码取坏值（改⇒红），本腿自己跑的 5 发

尺＝`sens-driver.sh` 的 `run()`（`go test ./internal/audio/ -count=1 -run ExtensibleFloat32306`），件＝`logs/sens/P*.txt`：

| 发 | 产码那一行的盘上形状（每发都逐字打回 `ZZ-summary.txt`） | rc | `--- FAIL` | 红句（逐字，本腿取的） |
|---|---|---|---|---|
| `P0-pristine` | ⛔ 突变 | **0** | 0 | `ok … internal/audio 0.047s`（绿＝基线） |
| `P1-mut192` | `if fmtTag == 0xFFFD {` | **1** | **2**（两枚指名用例） | `:177: parseWav error = unsupported wav format: tag=65534 bits=32 … want nil` |
| `P2-mut218` | `case fmtTag == 2 && bits == 32:` | **1** | **2** | 同形状 `tag=3 bits=32` |
| `P3-mut196-boundsshift` | `data[body+26 : body+28]` | **1** | **2** | 同形状 `tag=0 bits=32` |
| `P4-mut196-literalwording` | `data[body+26 : body+26]` | **1** | **1**＋`panic:` | ⛔ 任何断言句；只有 `panic: runtime error: index out of range [1] with length 0 [recovered, repanicked]` |

三枚红句**互异**（`65534`／`3`／`0`）⇒ ⛔ 一枚读数解释三格（见 §3）。

### 1.2 形二＝把断言的期望侧换成恒真（VA），看它还响不响

尺＝`mkVA()`：六枚期望侧逐枚前缀 `if false &&`（`rate != 48000`／`chans != 2`／`len(samples) != len(w306WantI16)`／
`samples[i] != want`／`len(inj.samples) != …`／`inj.samples[i] != want`，`VA_inert_count=6`），**只留两枚 `err != nil` 门**。
件＝`logs/sens/VA-*.txt`：

| 发 | rc | FAIL |  fired 的那一枚断言 |
|---|---|---|---|
| `VA-pristine` | **0** | 0 | 恒真形⛔ 自己红（＝变体合法、⛔ 假红） |
| `VA-mut192` / `VA-mut218` / `VA-mut196boundsshift` | **1** | 各 2 | `:177`／`:213` 的 `err != nil`（`tag=65534`／`tag=3`／`tag=0`） |
| `VA-mut196literal` | 1 | 1 | `panic`（与 P4 同形） |

⇒ **值断言全部恒真之后，那三枚突变照红**＝`err` 门是第二道独立的牙。

### 1.3 形三＝反过来问："它是不是只因为'返回了 err'才红？"（VR：把 err 门拔掉）

尺＝`mkVR()`：只把 `TestParseWavExtensibleFloat32306` 的 `:176 if err != nil {` 前缀 `if false &&`（`VR_lines=176:`），
值断言一字⛔ 动；`-run TestParseWavExtensibleFloat32306`（只那一枚，⛔ 让第二枚的 err 门替它响）。
件＝`logs/sens/VR-*.txt`：

| 发 | rc | 红句（逐字，本腿现量） |
|---|---|---|
| `VR-pristine` | **0** | —（拔掉 err 门⛔ 制造假红） |
| `VR-mut192` | **1** | `:182: rate = 0, want 48000` ＋ `:185: chans = 0, want 2` ＋ `:188: len(samples) = 0, want 8: 8 float32 samples are 32 data bytes …` |
| `VR-mut218` | **1** | 同上三句（`tag=3` 那一发） |
| `VR-mut196boundsshift` | **1** | 同上三句（`tag=0` 那一发） |
| `VR-mut196literal` | 1 | `panic`（崩在断言之前，⛔ 任何值判断） |

⇒ **把"确实返回了 err"那一形拔掉之后，三枚取坏值仍然让指名用例红，且红句是解析出来的值本身**。
这一条正是票面 `AC#1` 那句"断言的是**解析出来的那个值**，⛔ 断'被调用过'"的正面凭据。

### 1.4 形四＝把判据换成反形（逐枚翻转，⛔ 一枚恒真装饰）

尺＝`FLIP` 组：逐枚把一行里的 `!=` 翻成 `==`（产码 pristine），期望＝**每枚都该红**（红＝那个判断真的在Binding）。
件＝`logs/sens/FLIP-line*.txt`，`logs/sens/ZZ-summary.txt`：

| 翻转的那一行 | 行号 | 结果 |  fired 句（逐字） |
|---|---|---|---|
| `if len(raw) != 12+8+40+8+len(w306Floats)*4` | :171 | rc=1 FAIL 1 | 夹具自检（⛔ 产码的牙，只证它活着） |
| `if err != nil`（用例 1） | :176 | rc=1 FAIL 1 | — |
| `if rate != 48000` | :181 | rc=1 FAIL 1 | — |
| `if chans != 2` | :184 | rc=1 FAIL 1 | — |
| `if len(samples) != len(w306WantI16)` | :187 | rc=1 FAIL 1 | `:188: len(samples) = 8, want 8 …` |
| `if samples[i] != want` | :192 | rc=1 FAIL 1 | `:193: sample 0 = 32767, want 32767 …`（逐枚，⛔ 死循环） |
| `if err != nil`（用例 2） | :212 | rc=1 FAIL 1 | — |
| `if len(inj.samples) != …` | :217 | rc=1 FAIL 1 | `:218: injector holds 8 samples, want 8` |
| `if inj.samples[i] != want` | :221 | rc=1 FAIL 1 | — |
| 169 行以下全翻（除 :207 环境检查） | — | rc=1 FAIL **2** | 两枚指名用例各自红 |

⇒ **9/9 逐枚翻转都红**＝文件里⛔ 一枚恒真装饰、⛔ 空循环（`sample 0 = 32767, want 32767` 这类句子本身就是"值确实在流"的凭证）。

### 1.5 ★第四形："那道守卫分岔所需的世界，夹具里到底造出来没有"

尺＝覆盖块图，`-mode` 默认＝set（读数只有 0/1，本腿⛔ 把它当命中次数读）。三发并排，件＝`logs/sens/C*.txt`＋名册 `logs/sens/ZZ-cov-summary.txt`：

| 块（`wavinjector.go`） | 语义 | `C2`＝夹具 `mv` 走、整包 | `C1`＝夹具在、整包 | `C3`＝**只**跑那两枚新用例 |
|---|---|---|---|---|
| `188.4,192.24` | fmt 支到 `if fmtTag == 0xFFFE` 头 | **1** | 1 | 1 |
| `193.5,193.18` | `if size < 40` 那一判 | **0** | 1 | **1** |
| `194.6,195.1` | `size<40` 的错误 return | **0** | **0** | **0** |
| `196.5,196.65` | `data[body+24 : body+26]`（第三枚权威值） | **0** | 1 | **1** |
| `212.3,214.26` | PCM16 支体 | 1 | 1 | **0** |
| `219.3,221.26` | float32 支体（`:218` 命中之后） | **0** | 1 | **1** |
| `222.4,223.1` | float32 解码循环 | **0** | 1 | **1** |
| `224.3,224.43` | `return FloatToPCM16(f) …` | **0** | 1 | **1** |
| `226.3,226.111` | default（unsupported）支 | 1 | 1 | **0** |

⇒ 判语：**世界真的造出来了，而且是同一支。** `C3` 那一列是决定性的：只跑那两枚新用例时，
`:192` 的判定块、`:196` 的 SubFormat 读、`:218` 之后的 float32 三支（体／循环／return）**同时**从 0 变 1，
而 PCM16 支（`212.3,214.26`）与 default 支（`226.3,226.111`）保持 0
＝这一枚面⛔ "走到其中一支、其余两支仍零执行者"，也⛔ 靠 default 错误支混进"被执行"的名册。
唯一仍零执行者的是 `194.6,195.1`（`size<40` 的错误 return）——文件第 61-62 行自己写死了⛔ 装它，
票面 §边界也写着⛔ 顺手补别的分支 ⇒ ⛔ 本格缺陷，具名残余见 `30-judgments.md` §AC#1。

## 2. ★必答 3：`AC#2b` 的形——本腿选**甲**

两发都本腿自己重跑过（`P3`＝甲、`P4`＝乙），另外补了两发名册级现量：

- 甲（`data[body+24 : body+26]` → `data[body+26 : body+28]`，上下界同移）：`rc=1`、**两枚**指名用例逐枚 `--- FAIL`、
  红句 `tag=0 bits=32`（＝一个关于那枚权威值的**判断**）；整包尺 `D6`：`rc=1`、FAIL 恰好那 2 枚、顶层 PASS 仍 42（⛔ 连带红）。
- 乙（派单字面＝只动起点，`data[body+26 : body+26]`）：`rc=1`、FAIL **1 枚**、`panic:` 1 行；
  本腿再跑一发整包 `-v`（件＝`logs/N4-fixture-mut196literal-fullpkg-v.txt`）：
  `topFAIL=1 topPASS=36 topSKIP=1 RUN306=1 panics=1`
  ⇒ 名册里只有 `=== RUN TestParseWavExtensibleFloat32306` 那一枚 306 用例**开始过**，
  第二枚指名用例（`TestWavInjectorExtensibleFloat32306`）**根本没跑到**，且顶层 PASS 从 44 掉到 36＝**整个二进制被截断**。
- 两枚的正控都在盘上：夹具 `mv` 走后，甲（`D4`）与乙（`D5`）都是 `rc=0`、顶层 PASS 42、`panic` 0 枚
  ⇒ 那一枚 `panic` **确实**只在"新夹具走到 `:196`"时发生（＝它证明的是**可达**）。

**本腿的判语（写清那两问）：**

1. **panic 证不证明"`:196` 的偏移有人钉"？——证明"有人走到"，⛔ 证明"有人钉"。**
   钉＝一个会判断读出来的值的断言。乙那一发的失败来源是 `encoding/binary` 的运行时越界检查，
   件里⛔ 一句 `wavinjector_extensible_float_306_test.go:<行号>:` 的断言输出（本腿逐字 grep，见 §1.1 那一行的红句列），
   它说的是"这枚偏移让一次读越界了"，⛔ 是"这枚偏移读出来的值⛔ 等于 3"。
   **炸二进制与断言失败在"有没有牙"这件事上⛔ 是同一回事**：断言失败是**每枚指名用例各自**给一句关于值的判词；
   崩溃是把判词⛔ 产生、并把同批其余 35-36 枚顶层用例一起带走（`topPASS` 44→36、第二枚指名用例⛔ 跑到）⇒
   票面 `AC#2b` 要求的"成对两发"与"同一枚指名用例红"在乙那一发上只满足了一半（1 枚 FAIL、0 句判词）。
   还有一层**敏感度差**：乙只对"挪到越界"这一枚取值有反应；挪到仍越界的其它偏移（如甲那枚 `+26/+28`）走的是判断那侧，
   而挪到界内（例如 `+22`）乙形⛔ 响——一枚只会因越界而响的凭据，钉⛔ 住那枚偏移的值。
2. **盘上还有一枚同族先例支持甲**：姊妹件 `internal/audio/parse_wave_format_300_windows_test.go:116-119`
   逐字＝`name: "S4_data1_shifted_to_26", face: face300{… data1At: 26}, wantTag: 0, why: "positive control for offset sensitivity …"`
   ⇒ 在本仓自己的钉法里，"Data1 晚两字节"这一形是**一枚有期望值（`wantTag: 0`）的子测试**（本腿现量它 PASS，件＝`logs/Y3-ticket300-pins.txt`），⛔ 一次崩溃。甲与这枚先例同形。

**选形＝甲**（本格凭据＝`P3`/`D6` 那一对：改⇒两枚指名用例各自断言红、还原⇒绿）。
乙⛔ 丢、⛔ 升格：本腿把 `P4`/`N4`/`D5` 三发完整留在件里，读作**"新夹具是 `:196` 唯一执行者"的可达性凭据**（它⛔ 是 `AC#2b` 的牙）。

**要不要在票面追加一条更正？——要，且只追加（append-only，⛔ 删原句）。** 建议措辞（由编排者落笔，本腿⛔ 动框）：
`AC#2b` 那句"只改那一处偏移、语句其余字不动"的字面一发在盘上产出的是 `data[body+26 : body+26]`＝合法空切片
⇒ `panic: runtime error: index out of range [1] with length 0`、整个测试二进制中止、第二枚指名用例⛔ 跑到；
⇒ 该格的凭据取"同一枚切片上下界同移"那一发（判词逐枚、名册可归因），字面那一发降为可达性正控。

## 3. 必答 4：`AC#2` 与 `AC#1` ⛔ 同一发的两个说法

| 格 | 凭据的**尺名**（各自独立，本腿各跑各的） | 读数 |
|---|---|---|
| `AC#1` | ①覆盖块图 before/after（`C2` vs `C1`/`C3`，§1.5）②值断言拔门对照（`VR` 组，§1.3）③`:218` `3`→`2` 那一发（`P2`） | `C2` 五块 0→`C1`/`C3` 1；`VR` 三发红句是值；`P2` rc=1／FAIL 2／`tag=3` |
| `AC#2` | `:192` `0xFFFE`→`0xFFFD` 那一发（`P1`）＋同树无夹具正控（`D2`） | `P1` rc=1／FAIL 2／`tag=65534`；`D2` rc=0／PASS 42 |

⇒ 两格的**主尺⛔ 同**（覆盖块图 vs 突变红句），连突变枚都⛔ 同（`:218` vs `:192`），红句取值也⛔ 同（`3` vs `65534`）。
本腿⛔ 用任何一发读数追认两格。⚠ 唯一**共享**的东西是那枚夹具本身（票面 §3 射程那句"一枚夹具同时给三处装牙"是设计意图，
"三枚各要自己那一发"这条本腿逐枚量到了）。
