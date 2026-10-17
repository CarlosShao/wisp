# 335-v1 判语＝票 335 七格终裁（非实现者腿 335-v1，2026-10-17 13:5x +08）

被裁腿＝`335-r1`（六枚提交 `9567b094`/`77d98354`/`a58b12dc`/`ce7224c2`/`64265f0e`/`23d36a81`）。
本腿全部读数出自**自己现跑**：台面＝仓外导出树 `D:\V335V1\T`（`git archive 23d36a81` 八件套＋仓根 `docs/`
＋手工三件 dll/dist；`.go`＝758 证非空；D: 卷短名目录；`GOTMPDIR` 先 mkdir）；
`internal/panel` 一律读**具名提交树**（判定侧 `77d98354`、终态侧 `23d36a81`），起手与收尾各跑一把
`git status --porcelain -- internal/panel`（两把皆空＝工作树对本人只读，`328-r1` 在飞字节零接触）。
下文【跑】＝真跑出来的；【推】＝读码推的，两层逐句标开。

## 逐格判语

**AC#0 起手复跑＝成立**【跑】旧钉子字节（`git show a06a73d8` 抽出）对范畴产码重跑＝`rc=1`、红句整 14 行＝`:448`×7＋`:451`×7、
七枚＝`U+2061 2062 2063 2064 061C 180E 206A`、`U+206B` 0 命中（`logs/start-red-rerun.txt`，还原哈希与前置同串）。
"7/7 在 `Cf` 范畴"的独立尺＝改形后用例臂 1 的前提断言（`:474` 逐枚 `unicode.Is(unicode.Cf,r)`）在本腿台面
绿跑（shot A0，`rc=0`）＝每一枚都被实测过一遍，⛔ 借票 333 的结论。三处过期文字按内容锚现认（diff 在 `logs/ac2-diff-77d98354.txt`）。
名册枚数在导出树尺上读 7（本腿）＝票面 5＋两枚台面红，机制见 AC#3，与腿的"具名更正 8⇄7"同形。

**AC#1 改形后的闸不是自指＝成立（本腿自造反形，未借 333-v1 的 M1/M3）**
正形读数＝臂 2 实际侧 `washSetOver(isDisplayWashCf)`（`:498`）⇄期望侧范畴闭包（`:499`）。反形＝把 `:498` 抄成 `:499` 那枚表达式（内容锚替换，2 行）。【跑】三发对拉：
- shot B＝正形测试＋产码退回 15 枚手抄名册：`rc=1`，红 `:480 :484 :505 :522`，臂 2 读数 `product-washed=15 missed=155`；
- shot D＝**反形测试＋同一漂移产码**：`rc=1`，但红只剩臂 1 的 `:480 :484`，臂 2 沉默且**谎报** `product-washed=170 missed=0`（产码其实只洗 15 枚）；
- shot C＝反形测试＋正常产码：`rc=0` 全绿（反形能伪装无害，才需要 D 定罪）。
⇒ 反形与正形在同一漂移世界里读数不同（D⊂B 的断言行集、`:505/:522` 只在正形侧红、`missed` 155⇄0），**正形臂 2 量的确实是产品**，此闸非恒真。
**改名不留痕（腿没问的那层）**【跑】`git grep` 于 `23d36a81`：代码面引用用例名的只有测试文件自己（`-l` 尺）；
⇒ 若有人改用例名，**没有任何会红的代码侧尺**，只有红名册作差这道**流程**（`grep '^--- FAIL'`⇄`comm` 对 pre/post 日志）会露出 lost=1。写进判语，不翻格。

**AC#2 三处散文改完说实话＝成立（附一处措辞保留）**
零改动尺复跑【跑】（腿的 `zc335r1.sh` 原样，输入取自具名树）：`a06a73d8`⇄`77d98354`⇄`23d36a81` 剥整行注释后
同串 `5274713dead6…8b305`、`diff-rc=0` 两发；**正控**（副本里非注释那枚 `Cf`→`Cc`，token 现数命中 1 处）⇒ 同尺
`same_sha=no diff-rc=1`、`bdb58d7b…`；**负控**（只改一枚整行注释）⇒ `rc=0`（本尺按构造对整行注释失明）。
逐字读三句新文字【推】：`:90` "format rune above U+00FF (any member of the Cf category) -> \uXXXX"——不再指名册；
"any member" 措辞偏松，但该行左边已限定 above U+00FF，`:93` 同段自己举 `U+00AD -> \xad` 例外 ⇒ **不实未犯、留一处措辞松**。
`:92` "every invisible Cf-category rune lands in"——与产码事实相符：写手分支 `case isDisplayWashCf`（`:124`）收全部 `Cf`，
转义拼写分 `≤0xff/\xNN` 与 `>0xff/\uXXXX` 两支且同段并举 ⇒ 真。**`:193` 那句 `Cc` 支**【跑＋推】"C0 U+0000-U+001F and C1 U+007F-U+009F"
＝`unicode.IsControl` 的射程，与产码调用点（`:121/:204`）对得上；`case 0x` 全文件命中 0＝产码不枚举个位。

**AC#3 红名册归位＝成立，fixture 字节差机制归到底（本腿复现，不再"未归到底"）**
【跑】终态导出树红名册 7 枚逐名（`logs/ac4-pair-summary.txt`）与票起手 5 枚作差 gained 0、lost 0；多的两枚具名台面红：
`TestReadGitOnThisRepositoryIsSelfConsistent`（归档树无 `.git`）与 `TestComposerRenderFixtureTellsTheTruth`。
后者红因本腿量穿【跑】：`frontend/fixtures/composer-states.html` 三源字节＝**blob `23d36a81`=16604（CR 0/LF 45）⇄
`git archive` 导出件=16649（CR 45＝每行 CRLF）⇄ 工作树件=16641（CR 37/LF 45）**；8 字节差＝恰好 8 行裸 LF（行号 1、14、15、16、30、31、32、45，
perl 尺逐行现数，`logs/fixture-lf-lines-perl.txt`）；机制＝`core.autocrlf=true`＋`.gitattributes` `* text=auto` 令 **archive 导出时按 smudge 全量换 CRLF**，
工作树那份 8 行只剩 LF；两边剥 CR 后与 blob `cmp` 逐字节同⇒`git status` 对它是干净（实测）。
用例在导出树读到的字节形状多了 45 枚 CR 故红、在共享树读混合 LF 故绿＝台面字节差、非内容、⛔ 计入 335 判据。
**那 8 行 CR 是哪个进程剥的未归到底**【推】（`frontend/**` 归机主外部 agent，本编队零写面、无从取证）——如实具名。
起手 5 枚逐名复认同名集合【跑】；逐枚归因（票 35 两枚有意留红／两枚前端契约键／`tokens.css` 一枚）＝腿与编排者的账，本腿只裁"不在 335 射程"。

**AC#4 两把临时根台面成对＝成立**【跑】同一棵树同一条命令（`go test -v ./internal/panel/ -count=1 -timeout 30m`）：
c 侧（`TMP/TEMP=/tmp`（C:）、`GOTMPDIR` 不指）＝`RUN=257 PASS_TOP=160 FAIL_TOP=7 SKIP=0 rc=1`；
d 侧（`TMP/TEMP/GOTMPDIR=/d/V335V1/gotmp` 先建）＝`RUN=257 PASS_TOP=160 FAIL_TOP=7 SKIP=0 rc=1`；两侧红名逐字同（diff 空）。
四数各尺独立 `grep -c` 现量、⛔ 混用⛔ 相减；零枚红句含 `walk|build failed|no test files|setup failed|0xc0000135`（现数=0）。⛔ 跑 `./cmd/wisp` 整包、⛔ `go test ./...`、PID 21892 零接触。

**AC#5 洗集范围一字未动＝成立**【跑】blob 对撞：`a06a73d8`=git`230d3443…` ⇄ `77d98354`=git`6ab8ad85…`＝`23d36a81`=git`6ab8ad85…`
（终态与判定侧同一 blob）；两树之差＝3 枚整行注释（剥注释 sha 同串即 AC#2 那把尺）；导出树物理件全突变跑完后 `sha256=5c68a930…`
＝腿报的还原哈希同串，且 `git hash-object` 抽出字节＝`6ab8ad85` 三枚点名：注释行差异在 `:90/:92/:193`，判定层谓词 `:164` 逐字 `return unicode.Is(unicode.Cf, r)`。
六枚提交 `--name-only` 并集⛔ 含 `internal/risk/**`、`cmd/**`、三枚冻结件、`thresholds.go`、golden、`allowlist.txt`、`go.mod`/`go.sum`/`.gitignore`/`ci.yml`【跑】⇒ 逐枚零字节。

**AC#6 纪律格＝成立**【跑】严格前缀尺（行首 `票 335 ` 带空格）6/6 中；子串尺会多吞两枚非本腿提交（`03b827d8`、`e3333c13`＝票 333 面上提了"票 335"）——编排者 10-17 自踩的坑本腿复现具名。
唯一路径⊆｛`internal/panel/**` 两枚文件、`.scratch/wisp/probes/335/r1/**`、票尾一条｝。`t.Skip` 新增＝0（两枚代码提交各自尺）；`-func Test` 删除＝0；用例名在终态树逐字在册。
⛔ 在共享工作树跑 `go test` 那半格属**不可盘核自报**：腿的共享树终态一发（`257/162/5/0`）本腿只能当自报读、⛔ 计入凭据（本腿凭据全部出自仓外导出树）。

## 派单 ② 那一问（ⓐ/ⓑ 对称差）的裁语

**裁：尺已够敏感，第四发"逐枚加进一枚 `Cf` 成员"⛔ 必要——本腿实测它是空操作。**【跑】shot E＝在范畴形状的产码里显式加 `if r == 0x2061 { return true }`：
`rc=0`、臂 2 读数与 A0 基线逐字同（`170/170/0/0`）——该形**造不出红**，加它零增量凭据。凭据形状已钉满的论证【跑＋推】：
"名册被悄悄加长"唯一可响的形＝加**范畴外**成员，发 ⓒ（本腿 shot G，`U+00A0` 纳入）承担＝红 `:526 :540 :543`、`extra=1`，与 ⓐⓑ 断言行集**交集为空**；
"名册被摘减"由 missed 臂承担，最小敏感＝发 ⓑ（本腿 shot F，摘 `U+2061`）`missed=1` 红 `:480 :484 :522`；
ⓐ 独红 `:505`（地板绊线，本腿 shot B 同位复现）把"回到手抄"与"单枚特例"分开。子集关系（F 行集⊂B 行集）是构造使然，
腿的"级联⛔ 读成更强"纪律本腿照守：分离凭据逐发是**作差**读出的，不是枚数。三发还原哈希复认＝三方同串 `5c68a930…`（本腿七发 PRE⇄RESTORE 全对同串，见各 `shot-*.txt`）。

## 两枚必答

**(i) 若裁定再把洗集收窄，这台仪器会红还是会静？它买的是"产码漂移"还是"裁定变更"？**
会红，两种收窄形状都实测过【跑】：收窄回手抄名册＝shot B 红 `:505`（`len(product)<=15` 绊线）＋`:522`（missed=155）；
收窄成 `Cf` 的任一真子集（＝特例摘除族）＝shot F 红 `:522`（missed≥1 点名到枚）。它买的是**产码漂移**——谓词与范畴的相等关系由测试的期望侧
锚在 `unicode.Cf` 上，产码单方面动就红；**裁定变更**它买不起也不该买：改裁定唯一受控的过法＝连测试期望侧一起改的那一笔成对编辑，
而那笔编辑本身就是"此路须经裁定"的留痕。若将来有人只改产码期望它静——已证不静（B/F）；若有人把测试也改成名册形——`:505` 绊线红（B 实测）。
最低限度的定向突变形状（供后人复打，⛔ 真打跟踪件）：把 `isDisplayWashCf` 体内 `return unicode.Is(unicode.Cf, r)` 整行换成 15 枚手抄 `switch`（本腿 `mut335v1.sh` 发 B 原文可重放）。

**(ii) 票 335 有没有哪一格其实没修完？**
逐枚点名，⛔ "七格全绿"含糊：
1. **绊线不覆盖 `Cc` 臂**【推，构造即证】：臂 2 两枚被收集者都是 `Cf` 侧（`:498/:499`），臂 3 样本皆非控制符——`Cc` 那支若被换成手抄清单，此钉静。
   覆盖在邻钉 `internal/panel/display_scrub_333_333r1_test.go:44`（全家族扫描含 `Cc`，`:118-135`，本腿台面绿跑【跑】）⇒ **不算 335 欠修**，
   但若裁定要**同一枚钉**双族守卫，那是另立票的活，本判语具名留此。
2. **改名不留痕**【跑】：代码侧无一枚尺按名引用本用例（见 AC#1），改名的痕迹只活在流程作差里＝仪器射程外的程序性洞，写进判语。
3. `:90` "any member of the Cf category" 措辞松一行（同段 `:93` 自证例外，判不实未犯）＝**措辞保留、不返工**。
4. 台账 `A1017` 的旧枚数（起手 5⇄交件 6）与本腿台面读数（5+2 台面⇄7）之差已在本件归因（`.git`／CRLF 两枚台面红）；**登记与否归编排者，本腿⛔ 动台账一字**。
除这四条具名外，七格判**成立**。票面七框**本腿一枚没翻**（终裁凭据＝本件＋`logs/`，翻框归编排者）。

## 自证（三把尺）

- 本腿逐枚提交：起手状态＋AC#2 diff＝`e0353344`；脚本＝`fb967707`；七发突变读数＝`939f5f93`；fixture＋起手红复跑＝`2ff5b17a`；zc/AC4/AC5/AC6 审计＝`ca60e624`；本判语＝下一笔；票尾 bullet＝末笔。
  各笔 `git show --name-only` 全集⊆｛`.scratch/wisp/probes/335/v1/**`、`docs/evidence/s1/335-v1-verdict.md`、票 335 文件｝；`.go` 命中＝**0**（尺＝对各笔 name-only 全集 `grep -c '\.go$'`，见收尾 log）。
- 证据件绝对路径：`D:\work\workspace\projects plans\Wisp\docs\evidence\s1\335-v1-verdict.md`（本件，`wc` 与行数见收尾 commit bullet 内）＋
  `...\Wisp\.scratch\wisp\probes\335\v1\`（`mut335v1.sh` ＋ `logs/` 25 枚，逐枚 `wc -c`>0，枚数尺＝`ls logs | wc -l` 现数）。空槽尺（分词形，逐枚精确命令在 r1 `logs/placeholders.log` 第 1 行）对本件与脚本逐枚现跑＝0 命中。
- 决定性尺逐字：起手侧红＝`logs/start-red-rerun.txt`（14 句 `:448/:451`×7，rc=1）⇄ 交件侧绿＝`logs/mut-A0.log`（`:519 arm 2: product-washed=170, unicode.Cf=170, missed=0, extra=0, cfRoster=15`，rc=0）；
  成对突变红句与还原哈希＝`shot-{A0,B,C,D,E,F,G}-summary.txt` 逐字（七发 PRE⇄RESTORE 逐对同串；产码基线 `5c68a930…`、测试基线 `999672af…`）。
- 台面卫生：仓外树 `D:\V335V1\{T,gotmp,mut,zc}` 只建不删；`internal/**`/`cmd/**` 本腿零字节（工作树两把 status 尺：`logs/start-panel-status.txt`＋收尾同名件）。
