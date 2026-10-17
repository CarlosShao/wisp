# 332-a1 件1＝AC#0 起手四把尺 ＋ AC#1 名册（`Cf` 全集 170 枚里未被洗的那 155 枚）

腿＝`332-a1`（只读普查，⛔ 产码 ⛔ 测试件 ⛔ `go test`/`go build`/`go vet` ⛔ 判语）。
锚＝具名提交 `7535d83b`（dev，2026-10-17 10:07 +0800）；读数一律 `git show 7535d83b:<路径>`，
⛔ 读工作树在飞字节（`internal/panel` 有 `327-r2` 在写、`cmd/wisp` 有 `228-r3` 在写）。
逐枚原始转储＝本目录 `logs/`（枚枚已入库，见件3 §自证①）。

## AC#0 起手四把尺（★"量的还是算的"逐枚写明）

**尺① 显示边界名单逐枚＝15 枚（量的）**
尺＝`git show 7535d83b:internal/panel/display_scrub_216.go` 取 `isDisplayWashCf` 的 `switch` 支，
数其中码面枚数＝**15**（`:166-170` 一条跨行 `case`；⚠ `grep -c 'case 0x'` 回 **1** 不是 15，
因为 15 枚写在同一条 `case` 里——这一枚先记下来免得下一位照那把尺读成"只有 1 枚"）。
名单：`0x00AD, 0x200B, 0x200C, 0x200D, 0x2060, 0xFEFF, 0x202A, 0x202B, 0x202C, 0x202D, 0x202E, 0x2066, 0x2067, 0x2068, 0x2069`。
同一名单被写入侧 `scrubForDisplay`（`:123`）与判据侧 `needsDisplayWash`（`:214`）各调一次（票 327 AC#2 的"一枚名单两处调用"）。
转储＝`logs/roster-display-scrub-216-at-7535d83b.log`（239 行全文）。

**尺② 全标量 `unicode.Is(unicode.Cf, r)` 枚数＝170（量的，独立复跑，⛔ 抄 `327-v1`）**
尺＝仓外单文件纯 stdlib 探针 `D:\tmp\332a1\cfprobe\main.go`（`go run` rc=0，stderr 空），
遍历 0..`utf8.MaxRune`、跳过代理区 `D800-DFFF`、逐枚 `unicode.Is(unicode.Cf, r)` 计数。读数逐字：
`GO_UNICODE_VERSION=17.0.0` ／ `CF_TOTAL=170` ／ `CF_IN_ROSTER=15` ／ `CF_OUTSIDE_ROSTER=155` ／
`ROSTER_SIZE=15 ROSTER_NOT_CF=0` ／ `MISMATCH_COUNT=0`
（★最后一枚是本探针的自证：170 枚逐枚送进**逐字重打**的 `scrubForDisplay`，比"名单判定"与"实际是否被换掉"，
两轴 0 枚不一致 ⇒ 下面那张表的"洗／穿"这一列是量出来的、不是按名单推出来的）。
正控（报任何 0 命中之前跑过）：`U+00AD` 那一行 `out="\xad" WASHED`、`U+0600` 那一行 `out="\u0600" RAW`，两形都在盘上。
对合尺（另一枚独立来源，算的）：unicode.org 官方 `UnicodeData.txt`（17.0.0，2243593 字节）里 `;Cf;` 行 **170** 枚，
与 Go 那 170 枚逐枚比＝`IN_GO_NOT_UD=[]`、`IN_UD_NOT_GO=[]` ⇒ 两个来源在同一版本上**同一枚集合**。
转储＝`logs/ruler2-cf-170-census-go-probe.log`（178 行）＋`logs/ac1-cf-170-roster-with-unicode-names.log`（171 行，逐枚码位／官方名／Category／名单内外／行为）。

**尺③ 显示边界的生产调用点枚数＝1 枚（量的）**
尺＝`git grep -n 'scrubComposerForDisplay\|scrubForDisplay' 7535d83b -- internal cmd | grep -v _test`（射程＝两目录全量 `.go`，滤测试）。
生产调用者**只有 1 枚**＝`internal/panel/composer.go:145`（`Composer: scrubComposerForDisplay(composer)`，在 `NewSnapshot` 体内）。
其余 7 行命中全在同一件内（`:106` 定义、`:227` 定义、`:228`/`:234` 两处赋值、`:70`/`:201`/`:222` 注释）。
⇒ 显示边界洗的字段**只有两枚**：`display_scrub_216.go:228`（`Workspace.Canonical`）与 `:234`（`Attachments[].Name`）。
⚠ 这两行行号相对票 327 普查时点（`327-a1` 记作 `:169`/`:175`）**已漂移**，`327-r1` 落地 `isDisplayWashCf` 把文件加长了。

**尺④ 判定层那张相反名单＝6 枚显名 ＋ 整个 `unicode.Mn` 族（量的；★⛔ 动它一字）**
尺＝`git show 7535d83b:internal/risk/taintmatch.go` 读 `:39-45`，显名码位 6 枚＝
`0x00AD, 0x200B, 0x200C, 0x200D, 0x2060, 0xFEFF`，随后 `return unicode.Is(unicode.Mn, r)`。
Mn 族枚数＝**2059**（尺＝仓外纯 stdlib 探针 `D:\tmp\332a1\mncount\main.go`，`go run` rc=0：`MN_TOTAL=2059`）
⇒ 那张表的真实射程＝`6+2059=2065` 枚，但**落在 `Cf` 轴上的只有那 6 枚**（同一探针 `MN_INTERSECT_CF=0`，
另有件2 §筛检尺的第 3 列在 155 枚上逐枚复算＝0 枚响）。
两表交集＝那 6 枚（`TAINT_NAMED_ALSO_IN_DISPLAY=6`，算的）⇒ 显示名单 15 枚 ＝ 判定名单 6 枚 ＋ 方向族 9 枚（`U+202A–U+202E`、`U+2066–U+2069`），
方向族**不在**判定表里。
**为什么取向相反**（照它自己的注释逐字读，`:32-38`）：那张表服务的是**判定轴**——
"U+00AD/U+200B..U+200D/U+2060/U+FEFF and the combining marks (Mn) render as nothing (or as a diacritic on the previous rune)
but are NOT unicode.IsSpace, so leaving them in place lets an exfiltrator cut a tainted fragment into sub-8-rune pieces
with zero-width characters"，出处＝对抗报告 `16.9#1`/`M-6`。即判定要**把看不见字符从归一化串里丢掉**
（`normalizeTaint` `:59` 用它做删除），好让 C25 的 ≥8 字符片段匹配不因零宽切口失效；
显示边界要的是**把看不见字符显出来**（转义、⛔ 删）。一删一显，同一枚码位在两轴上的正确动作相反，故两份名单各写各的、
`display_scrub_216.go:142-147` 也自陈"NOT internal/risk's isIgnorableTaintRune"。
转储＝`logs/roster-taintmatch-at-7535d83b.log`（283 行全文）。

## AC#1 名册：整族 170 枚，未被洗的 155 枚逐枚在下表（⛔ 抽样）

逐枚底表（170 行，码位／Unicode 17.0.0 官方名／Category／`IN`|`OUT`／`WASHED`|`RAW`）＝
`logs/ac1-cf-170-roster-with-unicode-names.log`；下面是同一张表按 Unicode 名分族求和，
★分族尺＝码位区间＋官方名族名（算的，脚本在件3 §自证③逐字给出），**求和＝155、LEFTOVER=[]、GHOST=[]**（既不漏也不重）。
名单内 15 枚另列（件1 开头尺①那 15 枚），不进这张 155 表。

| 族 | 枚数 | 码位范围 | 官方名族（逐枚见底表） | 三形归类（本腿⛔ 测像素） |
|---|---|---|---|---|
| A | 13 | `U+0600–U+0605`,`U+06DD`,`U+070F`,`U+0890`,`U+0891`,`U+08E2`,`U+110BD`,`U+110CD` | ARABIC NUMBER SIGN / SANAH / FOOTNOTE MARKER / SAFHA / SAMVAT / NUMBER MARK ABOVE / END OF AYAH、SYRIAC ABBREVIATION MARK、ARABIC POUND／PIASTRE MARK ABOVE、ARABIC DISPUTED END OF AYAH、KAITHI NUMBER SIGN(ABOVE) | 名含 MARK/SIGN＝**前缀标记族**（官方名不属 INVISIBLE/JOINER 那一支） |
| B | 1 | `U+061C` | ARABIC LETTER MARK | **方向族**（bidi 控制 ALM） |
| C | 1 | `U+180E` | MONGOLIAN VOWEL SEPARATOR | 票面点名的"蒙古元音分隔符"；Unicode 17 里 Category＝`Cf`（⛔ `Zs`，见件2 筛检尺：155 枚里 `unicode.IsSpace` 0 枚响） |
| D | 2 | `U+200E`,`U+200F` | LEFT-TO-RIGHT MARK、RIGHT-TO-LEFT MARK | **方向族**（bidi 标记，与已洗的 `U+202A–U+202E` 同一支） |
| E | 4 | `U+2061–U+2064` | FUNCTION APPLICATION、INVISIBLE TIMES、INVISIBLE SEPARATOR、INVISIBLE PLUS | **隐形／零宽族**；＝票 327 `AC#1` 里"具名裁要不要"那一族（3 枚官方名直写 INVISIBLE） |
| F | 6 | `U+206A–U+206F` | INHIBIT／ACTIVATE SYMMETRIC SWAPPING、**INHIBIT／ACTIVATE ARABIC FORM SHAPING**、NATIONAL／NOMINAL DIGIT SHAPES | 票面点名的两族都在此：`U+206A`/`U+206B`＝激活／去激活对称换形符，`U+206C`/`U+206D`＝阿拉伯字母形式控制；Unicode 17 里整族已 DEPRECATED |
| G | 3 | `U+FFF9–U+FFFB` | INTERLINEAR ANNOTATION ANCHOR／SEPARATOR／TERMINATOR | **隐形／分隔族**（行间注释控制） |
| H | 16 | `U+13430–U+1343F` | EGYPTIAN HIEROGLYPH VERTICAL／HORIZONTAL JOINER、INSERT AT …、BEGIN／END SEGMENT／ENCLOSURE／WALLED ENCLOSURE、OVERLAY MIDDLE | JOINER／边界控制族；★`U+13439–U+1343F` 是 16.0 才进 `Cf` 的（Python 14.0 表把它们读成非 `Cf`，本腿以 Go 17.0.0＋UnicodeData 17.0.0 为准） |
| I | 4 | `U+1BCA0–U+1BCA3` | SHORTHAND FORMAT LETTER OVERLAP／CONTINUING OVERLAP／DOWN STEP／UP STEP | 速记格式控制族 |
| J | 8 | `U+1D173–U+1D17A` | MUSICAL SYMBOL BEGIN／END BEAM／TIE／SLUR／PHRASE | 音乐符号配对控制族 |
| K | 97 | `U+E0001`,`U+E0020–U+E007F` | LANGUAGE TAG、TAG SPACE＋TAG <ASCII 逐枚同名>（`U+E0021`–`U+E007E` 一枚一个 ASCII 字符）、CANCEL TAG | **隐形族**；★这一族的官方名是"TAG ＋ 一枚 ASCII 字符名"，即 97 枚里 95 枚按名各对应一个 ASCII 字符（算的，尺＝名表逐枚前缀） |
| **合计** | **155** | | | 求和 13+1+1+2+4+6+3+16+4+8+97＝**155** |

★票面要求的三族点名对照（哪一族在这里）：数学作用符＝E（4 枚）；阿拉伯字母形式控制＝F 里的 `U+206C–U+206D`（2 枚）；
蒙古元音分隔符＝C（1 枚）；`U+206A`/`U+206B` 激活／去激活对称换形符＝F 里的 2 枚；
零宽／隐形一支＝E＋G＋K（4+3+97＝104 枚，算的）；方向一支＝B＋D（2 枚）＋F 的换形／数字形状 6 枚（同一支 Unicode "format controls for bidi and shaping" 名族，本腿只按官方名归族）。
⛔ 本件不给"该不该把这 155 枚（或其中任何一族）纳进名单"的判语——那是编排者与裁决腿的活。
