# panel-l2-grant-nail-fix-r4-accept-r1 —— 非实现者对"门钉 r4"这一批的对抗验收

> 被验物：`internal/panel/l2_grant_boundary_test.go` @ 锚点 **`ee2a92d`**（blob `58f545144af54c0b16eaed013a0904fb1ffd562a`，2220 行）
> ＋自述件 `docs/evidence/s1/panel-l2-grant-nail-fix-r4.md`（449 行，同批落地）。
> 任务来源：`docs/evidence/s1/panel-l2-grant-nail-fix-r3-accept-r1.md` §9 点名的 **F-ACC-1／F-ACC-2／F-ACC-3**。
> 本程＝**验收程**：一字未改被验物，只 commit 本件；不勾票面／台账；不 push。
> 本程全部变异落在**仓外副本**（§1），读数原始件在 `D:\tmp\wisp-r4acc1\out\`（§12 列路径）。

---

## §0 锚点与落点：本程自己量的，不采信简报

```
git rev-parse ee2a92d:internal/panel/l2_grant_boundary_test.go -> 58f545144af54c0b16eaed013a0904fb1ffd562a   (2220 行)
git rev-parse 88eab34:internal/panel/l2_grant_boundary_test.go -> 9fd6defbd8dd113a7bd9c70c619ab92ac01f32d7   (2081 行)
git rev-parse 88eab34:internal/panel/bridge.go                 -> d2cd6362ecc6941a0cee58073a2bf8ab6669a590
git rev-parse ee2a92d:internal/panel/bridge.go                 -> d2cd6362ecc6941a0cee58073a2bf8ab6669a590   ⇒ 生产码零字节（同值）
git log --oneline 88eab34..ee2a92d -- internal/panel/          -> 只有 594a99e 一枚（本批那一枚）
git show --numstat 594a99e                                     -> 145  6  internal/panel/l2_grant_boundary_test.go（路径集合只有这一枚）
git status --porcelain -- internal/panel/                      -> 空（我进场与交件各复算一次，全程未写进仓内）
```

ⓑ 级细节一条：**简报那句"代码落点只有一枚"按"区间"字面读不成立、按"本批"成立**。
`git diff --numstat 88eab34..ee2a92d` 里有 **7 枚路径**，除本批那一枚 `_test.go` 之外还有
`frontend/src/styles/theme.css`（14/3，取自 `8b35f52`，别家色值尺那一程）与 5 枚文档/会话件。
实现件 §5.5 自己就把 theme.css 点名成"别家"，所以**这不是它的缺陷，是我简报的口径含糊**，按本程纪律写在这里。

---

## §1 台件自证：两条仪器坑我都独立复现了（不是照抄它的说法）

### §1.1 正路（本程实际用的那一支）

```
git -c core.autocrlf=false -c core.eol=lf archive ee2a92d | tar -x -C /d/tmp/wisp-r4acc1/anchor
逐枚验字节：git ls-tree -r ee2a92d（1172 枚）对比 git hash-object --no-filters <副本>
   -> same=1166  mismatch=6
   那 6 枚全部是 *.ps1：scripts/build.ps1 · scripts/dev/ball-cycle.ps1 · scripts/fetch-deps.ps1
                        scripts/sign-models.ps1 · scripts/slo-check.ps1 · scripts/spike/run.ps1
   成因＝本仓 .gitattributes 逐字写着 `*.ps1 text eol=crlf`（显式属性压过 core.eol=lf），
   是仓自己要求的 CRLF，不是污染；`grep -o '"[^"]*\.ps1"' internal/panel/*_test.go` = 空
   ⇒ 本包没有一枚测试读 .ps1，这 6 枚差值不进我的分母。
两枚被验物逐枚：
   l2_grant_boundary_test.go -> 58f545144af54c0b16eaed013a0904fb1ffd562a = 锚点（字节同值）
   bridge.go                 -> d2cd6362ecc6941a0cee58073a2bf8ab6669a590 = 锚点（字节同值）
```

### §1.2 坑①：只给 `core.autocrlf=false` 会凭空造出一枚红（我复现到了）

```
git -c core.autocrlf=false archive ee2a92d | tar -x -C /d/tmp/wisp-r4acc1/crlfcopy
frontend/fixtures/composer-states.html 现量：crlfcopy 6187 字节 / CR 计数 8 / `-->` 后是 \r\n
                                          anchor  6179 字节 / CR 计数 0 / `-->` 后是 \n
go test ./internal/panel/ -count=1 -v  ->  RUN=99 TOPPASS=52 SUBPASS=46 FAIL=1
   唯一红名 = TestComposerRenderFixtureTellsTheTruth
   红句（逐字，out/PITFALL-crlf-baseline.txt:152 起）：
     composer_test.go:699: the render fixture has no block for "unsupported attachment told to the user" …
     composer_test.go:728: render fixture verified across 0 painted states
⇒ 简报与实现件 §1 那一支的读数是**真的**，我这一发独立跑到同值（同一枚红名、同一句"0 painted states"）。
```

### §1.3 坑②：不带 `--no-filters` 的 `git hash-object` 撑不住"副本＝锚点"这句

```
同一枚被 CRLF 污染的 frontend/fixtures/composer-states.html：
  git hash-object            <副本>  -> ab39390cca966ccc3ec2df4db6d4ed2028a119aa
  git rev-parse ee2a92d:该路径        -> ab39390cca966ccc3ec2df4db6d4ed2028a119aa   ← **两值相同**
  git hash-object --no-filters <副本>  -> 680b908b12c66112abc3231a4725e554b461d24d   ← 与锚点不同
⇒ 那把"看起来对"的尺把污染读成同值。本程下面每一枚"副本字节＝锚点字节"的主张都只由
   `--no-filters`（或 §1.2 那种字节级 CR 计数）撑着，简报这一条我复算成立。
```

台件脚本（本程自己的，不复用实现者的）：`D:\tmp\wisp-r4acc1\scripts\{mut14.py,run14.sh,decoy_demo.py,decoy_run.py}`。
每发变异：锚点计数不等于 1 即 rc=1 拒绝落盘（本程真实拒过 1 次：`MDEC-reg` 我把注册表写在了 `bridge.go` 上，
而 `inboundTypeRegistry` 其实在 `l2_grant_boundary_test.go:1560`，仪器当场拒绝、半行未落，改对后重打）；
落盘后打印两枚被验物的 `--no-filters` 哈希，每发结束 `restore` 回锚点同值（末次复算＝`58f5451…`/`d2cd636…`）。

---

## §2 简报那枚前提我自己现量了（ⓖ）：`internal/panel/` 是**一枚**红，不是两枚

```
仓库工作树（HEAD 取于我量这一发时，`internal/panel/` 干净、blob 同锚点）：
go test ./internal/panel/ -count=1 -v   -> rc=1  RUN=99 TOPPASS=52 SUBPASS=46 FAIL=1 SKIP=0 ^panic:=0
唯一红名（只认 `--- FAIL:`，逐字）：TestC21DesignTokensFourWayAgree
红句：tokens_fourway_test.go:441: read design/assets/tokens.css: open …\design\assets\tokens.css:
      The system cannot find the path specified.
grep -rn "func TestC21" internal/panel/  -> 1 枚（tokens_fourway_test.go:439）
```

**裁：简报错了，实现件对了。** "两枚 C21 红"这句在 `internal/panel/` 这一枚包上不成立——只有 1 枚。
**并且我把"两枚"这个数的来路也量出来了**（简报没让我查，但只登记不解释等于没核）：

```
go test ./internal/ball/ -count=1   （同一枚脏工作树）
  -> --- FAIL: TestC21TableColourRowsMatchTokensCSS   tokens_table_test.go:1468: read design/assets/tokens.css …
```

⇒ 脏树（`design/**` 那 16 枚未提交删除，owner 自己的活）在**两枚不同包**里各打红一枚 C21：
`internal/panel` 的 `TestC21DesignTokensFourWayAgree` ＋ `internal/ball` 的 `TestC21TableColourRowsMatchTokensCSS`。
"两枚"是按 `./internal/...` 口径读出来的，"在 `internal/panel/`"是按包口径写的——**两句各自都对，拼起来才错**。
本程不替简报圆：简报这一格写错了范围。实现件 §0/§12 那句"现量只有 1 枚、未修未跳未放宽"复算成立。

---

## §3 三态四数＋名册差集（ⓗ）：`^panic:` 三发皆 0，红名差集只有那一枚 C21 的有无

| 台件 | 四数（`-count=1 -v`） | 红名 | `asked=` | 原始件 |
|---|---|---|---|---|
| 改前副本 `88eab34`（blob `9fd6defb`）| RUN=99 TOPPASS=53 SUBPASS=46 FAIL=0 SKIP=0 `^panic:`=0 | 无 | asked=528 answeredByRealGuard=0 hits=[] | `out/BASE-pre.txt` |
| 交付副本 `ee2a92d`（blob `58f54514`）| RUN=99 TOPPASS=53 SUBPASS=46 FAIL=0 SKIP=0 `^panic:`=0 | 无 | asked=528 answeredByRealGuard=0 hits=[] | `out/BASE-anchor.txt` |
| 仓库工作树（同 blob，脏 `design/`）| RUN=99 TOPPASS=52 SUBPASS=46 **FAIL=1** SKIP=0 `^panic:`=0 | `TestC21DesignTokensFourWayAgree` | asked=528 answeredByRealGuard=0 hits=[] | `out/BASE-repo-worktree.txt` |

```
名册差集（三态各取运行输出，未用 -run / -skip）：
  顶层 `^--- (PASS|FAIL|SKIP)` 归一化名：53 -> 53 -> 53    diff(pre,anchor)=空  diff(anchor,repo)=空   ⇒ +0/−0
  子测试 `^    --- …`          ：46 -> 46 -> 46    两两 diff 皆空                                     ⇒ +0/−0
  `func Test` 名册（git grep -h '^func Test' <rev> -- internal/panel/）：88eab34=53 枚 -> ee2a92d=53 枚，diff 空
  本文件内 `func Test(`：7 -> 7（§0 断言尺同值）
`^panic:` 三发皆 0 ⇒ 没有"一条用例 panic 吞掉同包其余读数"的形状，上面那 99 条 RUN 是真取到的。
```

**裁（ⓗ）**：实现件 §5.1／§11 那三发四数、名册 +0/−0、`asked=528` 未变、`^panic:=0` **逐格复算成立**。
一处口径差要写下来免得下游误读：它 §5.1/§11 那几发的 `TOPPASS=52`（因为工作树里 C21 那枚红），
我这三发里副本态是 `TOPPASS=53/FAIL=0`、工作树态才是 `52/1`——**同一枚事实的两种台件**，不是读数冲突；
它的 §11 那一发明写"取于工作树、唯一红名＝那枚 C21"，与我的第三行一字同。

---

## §4 三门与断言尺（本程自己跑，各带当场 HEAD）

```
HEAD(我跑这一批时)=2ddbe65（共享树在漂；`git status --porcelain -- internal/panel/` 全程空）
gofmt -l internal/panel/  -> 空，rc=0
go vet  ./internal/panel/ -> 空，rc=0
sh scripts/d22scan.sh     -> rc=0
  bans #1-5 internal/=203 cmd/=22 | ban #6 frontend/=46 | ban #7 internal/tools/=18
  ban #8 design/=32 frontend/=46 internal/=407 cmd/=40 | d22scan: clean - no D22 ban violations
```

⇒ 与实现件 §5.3／§11 那两发八枚分母**逐枚同值**（它取于 `fb10fc1`／`ffe9f42`，我取于 `2ddbe65`；
本批只动一枚 `_test.go` ⇒ 只进 ban #8 的 `internal/=407` 分母，这句我按 §0 的路径集合核过）。

断言尺（两版都 `git show` 取，`grep -oF | wc -l`，不数行）：

```
                     88eab34(9fd6defb)  ee2a92d(58f54514)  delta
t.Fatalf(                    33                35           +2
t.Fatal(                      5                 5            0
t.Errorf(                    40                45           +5
t.Error(                      1                 1            0
t.Logf(                      15                15            0
t.Skip                        0                 0            0   ⇒ 保持 0
func Test(                    7                 7            0
Fatal/Error 合计             79                86           +7   ⇒ 只增不减
行数                        2081              2220
```

⇒ 实现件 §5.4 那一格**逐格同值**。本程没有为了变绿动过任何断言、任何阈值、任何 golden（§0 路径集合为证）。

---

## §5 ⓒ F-ACC-1：乘积第三因子——**四发掏空我全部自己造、自己打，不抽它的数**

被验控（交付版逐字，`l2_grant_boundary_test.go:1325`）：
`if len(grantRouteWords) == 0 || len(grantRoutePrefixes) == 0 || len(grantRouteSuffixes) == 0 {`
＋乘积式（`:1353`）：`if want := len(grantRoutePrefixes) * len(grantRouteWords) * len(grantRouteSuffixes) * 2; asked != want {`。
下面每一发都是 `D:\tmp\wisp-r4acc1\anchor`（§1.1 那台件）现跑 `go test ./internal/panel/ -count=1 -v`。

| 发（我这边的编号） | 我只改哪一行 | 四数 | 红名（逐字） | 红因首句（逐字） | 原始件 |
|---|---|---|---|---|---|
| `M1-suffix-empty` **本轮真正的新账** | `grantRouteSuffixes` -> `[]string{}` | RUN=99 TOPPASS=52 SUBPASS=46 **FAIL=1** | `TestRealGuardRefusesEveryAssemblableApprovalRouteName` | `:1380 the sweep vocabulary is empty (grantRouteWords=11, grantRoutePrefixes=8, grantRouteSuffixes=0): a sweep that asks nothing answers clean forever` | `out/M1-suffix-empty.txt` |
| `M1-no-plural` | `for _, form := range []string{w, w + "s"}` -> `[]string{w}` | RUN=99 TOPPASS=52 SUBPASS=46 **FAIL=1** | 同上 | `:1380 the sweep asked 264 names but the three lists it reads multiply to 528 (8 prefixes x 11 words x 3 suffixes x 2 plural forms): a factor that stopped being ranged is a grid that silently shrank, and hits=[] …` | `out/M1-no-plural.txt` |
| `M1-words-empty` | `grantRouteWords` -> `[]string{}` | RUN=99 TOPPASS=50 SUBPASS=44 TOPFAIL=3 SUBFAIL=2（**5 枚红条目**） | `TestRealGuardRefuses…` ＋ `TestGrantVocabularyIsNot…` ＋ `TestPlantedGrantWiringGoesRedInASnapshot`（父＋`/B_Go_answering_panel.approval.request_goes_red`＋`/C_a_wired_grant_door_goes_red_on_both_halves`） | `:1377 the sweep vocabulary is empty (grantRouteWords=0, grantRoutePrefixes=8, grantRouteSuffixes=3)` | `out/M1-words-empty.txt` |
| `M1-prefix-empty` | `grantRoutePrefixes` -> `[]string{}`（8 枚全删） | RUN=99 TOPPASS=52 SUBPASS=46 **FAIL=1** | `TestRealGuardRefuses…` | `:1378 the sweep vocabulary is empty (grantRouteWords=11, grantRoutePrefixes=0, grantRouteSuffixes=3)` | `out/M1-prefix-empty.txt` |
| `M2-all-three-empty`（简报没要求，我加：ⓑ 那枚"0==0 空转角"） | 三枚清单一起清空 | RUN=99 TOPPASS=50 SUBPASS=44 TOPFAIL=3 SUBFAIL=2 | 同 `M1-words-empty` 那 5 枚 | `:1375 the sweep vocabulary is empty (grantRouteWords=0, grantRoutePrefixes=0, grantRouteSuffixes=0)` | `out/M2-words-empty+prefix-empty+suffix-empty.txt` |

**裁（ⓒ）**：
- 上游 r3-accept §1.1 的 **V1**（掏空第三因子 -> 本测试 PASS、`asked=0`、全包 0 枚红）**今天不成立**：我这一发
  现量 **红 1 枚**，红句逐字点名 `grantRouteSuffixes=0`。⇒ **实现件 §2 那格"复判原来记的『这一形 0 枚红』已死"复算成立**，
  而且它选的方案（ⓑ 乘积式 ＋ 控里补齐第三因子）确实是把 V1 与 V4 一并堵死的那一支。
- 上游的 **V4**（分母减半 `asked=264` 仍绿）同样**今天不成立**：现量红 1 枚，红因是乘积式而不是那句"空"。
- 它 §2 那句"ⓑ 不能单独存在：三枚清单一起清空时 `asked == 0 == want`，它会绿，所以 ⓐ 是前提"——我加发了
  `M2-all-three-empty` 验这一角：**这一角由控（ⓐ）拿住**（红在 `:1375`），不是由乘积式拿住。⇒ 它"两枚互为前提"
  的说法是**活的判断**，不是我抄来的修辞：摘掉任一枚都会掉一发（见 §7 的 ⓔ）。
- 行号口径：同一句控在 words 那发报 `:1377`、prefix 那发 `:1378`、suffix 与 no-plural 两发 `:1380`、
  三枚一起清空 `:1375`——**调用点跟着我自己那一刀的宽度漂**（我删的行数与它不同：它 §2 记 prefix-empty 一发
  `delta_lines=-3` 报 `:1377`，我这发删 2 行报 `:1378`）。这不是读数冲突，是两把台件切法不同；
  它那句"行号会跟着文件漂、读数活的不是抄的"我这一发复现到同一形状。

---

## §6 ⓓ F-ACC-2：11 枚证人我**全部逐枚删过**（不只抽样），含它自称 4 枚红的那枚

台件同上。脚本只动 `var grantRouteWords` 那一块内的字面量（块内计数不等于 1 即拒），证人册里那枚同名词**不碰**。

| 删掉的词 | 四数 | 红条目数 | 红名 | `asked=`（扫掠自己那一发） |
|---|---|---|---|---|
| `approve` | 99/52/46/1 | 1 | `TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes` | 480 |
| **`approval`**（必查那枚） | 99/51/44/2 | **4** | 同左 ＋ `TestPlantedGrantWiringGoesRedInASnapshot`（父＋`/B_Go_answering_panel.approval.request_goes_red`＋`/C_a_wired_grant_door_goes_red_on_both_halves`） | 480 |
| `grant` | 99/52/46/1 | 1 | `TestGrantVocabularyIsNot…` | 480 |
| `allow` | 99/52/46/1 | 1 | 同左 | 480 |
| `permit` | 99/52/46/1 | 1 | 同左 | 480 |
| `ratify` | 99/52/46/1 | 1 | 同左 | 480 |
| `authorize`（上游标"未逐枚打＝推定"，我实跑） | 99/52/46/1 | 1 | 同左 | 480 |
| `authorised` | 99/52/46/1 | 1 | 同左 | 480 |
| `decide` | 99/52/46/1 | 1 | 同左 | 480 |
| `decision`（上游同样标"推定"，我实跑） | 99/52/46/1 | 1 | 同左 | 480 |
| `verdict` | 99/52/46/1 | 1 | 同左 | 480 |

红因逐字（两枚代表，`out/W-approval.txt` / `out/W-verdict.txt`）：

```
:1512 witness route "panel.approval.request" is flagged by [] under grantRouteWords as it stands, and this file
      wrote down [approval] for it: the vocabulary and the evidence attached to it are no longer the same list … (F-ACC-2)
:1540 grantRouteWords no longer carries "approval" while "panel.approval.request" still stands here as its witness:
      the sweep stopped asking the 48 names that word used to build … this roster is the only thing in the package
      that noticed (F-ACC-2)
:1512 witness route "panel.verdict" is flagged by [] … wrote down [verdict] …
:1522 the route name assembled here from the vocabulary word "verdict" ("panel.verdict.request") does not read as
      an approval door - that word no longer carries the meaning the sweep assumes …
```

**裁（ⓓ）**：**11/11 有红，无一枚装饰**——实现件 §3 那一整表（含 `approval` 4 枚、其余各 1 枚）我逐枚复跑到同值，
连"`asked` 每发都从 528 掉到 480 而扫掠本身照旧绿"这句也复现到（`out/W-*.txt` 里那枚 `asked=480` 就是扫掠自己的
`t.Logf`）。上游 r3-accept §1.2 明标"未逐枚打＝推定"的两枚（`authorize`/`decision`）本程**从推定升级为现量**。
证人册的三条主张不是"表自证含它自己"：判决者确实是盘上的谓词与真守卫（`out/W-verdict.txt` 那发一次报出
主张 1/主张 2/双向对账三处，各是不同行号、不同判据）。

顺带把它 §3 末那枚"加分读数"也复算成现量：**M16（生产码运行期拼 `panel.review.ratify`）在交付态红 2 枚**
（`out/M2-M16.txt`：`TestRealGuardRefuses…` 的 `hits=[panel.review.ratify]` ＋ `TestGrantVocabularyIsNot…`
@`:1515` 那句 witness 主张），**M14 红恰好 1 枚**（`out/M2-M14.txt`，只有常驻扫掠响）。
⇒ 同一枚生产码形状现在有两名独立证人、M14 那一形仍只有一名——它 §6 那句"必须照实写的差别"是真的。

---

## §7 ⓑ 承重判定：摘掉任意一味，是否还存在一发变异从此打不红

本仓对"承重"的操作定义就这一句，我按枚判。交付的 5 处改动我按盘上形状归成 **4 味**
（证人册那枚清单变量与读它的断言块**拆不开**：`revert-roster` 单发我打了，包直接编译不过
（`out/V-revert-roster.txt`：`RUN=0`，`grantRouteWordWitnesses` 未定义）⇒ 它俩是一味的两半，不是两味）。

| 味 | 摘掉它之后**唯一由它拿住**的那一发 | 我那一发的读数 | 判 |
|---|---|---|---|
| ⓐ 控里补齐第三因子（`:1325` 那枚 `\|\| len(grantRouteSuffixes) == 0`） | `suffix-empty`（掏空后缀清单） | `revert-empty-ctrl,suffix-empty` -> **RUN=99 TOPFAIL=0 全绿**（`out/T-revert-empty-ctrl+suffix-empty.txt`） | **承重** |
| ⓑ 乘积式断言（`:1353`） | `no-plural`（砍掉复数那一圈） | `revert-product,no-plural` -> **TOPFAIL=0 全绿**（`out/T-revert-product+no-plural.txt`） | **承重** |
| ⓒ 证人册＋其三条主张＋双向对账（`:183-225` 与 `:1490-1548`，两半） | `dropword-ratify`（删一枚"历史上没人质"的词）／`addword-consent`（加一枚没配证的词）／`standing-out,M16`（常驻扫掠被摘、真守卫答了证人名字） | `revert-witness-mechanism,dropword-ratify` -> **TOPFAIL=0**（`out/X-revert-witness-mechanism+dropword-ratify.txt`，与上游 r3-accept §1.2 那发同形）；`revert-xcheck,addword-consent` -> **TOPFAIL=0**（`out/U-revert-xcheck+addword-consent.txt`）；`standing-out,revert-witness-rows,M16` -> **TOPFAIL=0**（`out/Y-standing-out+revert-witness-rows+M16.txt`） | **承重**，且三发各指向册内不同分支 |
| ⓓ 头段"两步路"措辞（`:119-132`） | 无一发 | `revert-header` -> 四数一字同（99/53/46/0，`out/H-revert-header.txt`）；`revert-header,MDEC-struct` -> 仍红 1 枚、同一句 `:1699 decode destination acProbe has no reflection twin in inboundTypeRegistry`（`out/H-revert-header+MDEC-struct.txt`） | **不承重、也不该承重**：它是文档级修法，盘上判据不读注释（本文件 10 处 `F-ACC` 字样在 `:102/:120/:183/:1316/:1354/:1441/:1490/:1512/:1532/:1540`，除两枚消息文本外全是注释，无一处是谓词） |

**两条反向核对（防"互相顶账"与防"把冗余读成装饰"）**：

```
① 上游 r3-accept §9 反向判据第 1 条我独立重打：摘常驻测试 ＋ M14 -> RUN=98 TOPFAIL=0（out/R-standing-out+M14.txt）
   摘常驻测试单发（不叠码）：RUN=98 TOPFAIL=0（out/R-standing-out.txt）⇒ 摘它本身不造红。
   ⇒ **新加的证人册没有替常驻测试兜这一发**（册里 `allow` 那行写的是 `panel.l2.allow`，M14 答的是
      `panel.review.allow`，两条主张各自只对各自的名字负责）⇒ 实现件 §6① 那句"新旧不互相顶账"复算成立。
   同一台件下 `standing-out,M16` -> **TOPFAIL=1**（红句 `:1502 the running guard answers "panel.review.ratify",
   the witness name for the vocabulary word "ratify"`，out/Y-standing-out+M16.txt；行号比交付态的 `:1515`
   低 13＝我那一刀摘掉的正是上面那 13 行常驻用例，调用点跟着漂，与 §5 那格同一条规律），
   而 r3-accept §4.3 那两行记的是 0 枚红 ⇒ 这一格**读数确实变了**，实现件 §6 不但没藏，还自己写清了归因
   （M16 那枚名字正是它给 `ratify` 挑的活证据名字）⇒ 加分，不是缺陷。
② 冗余不等于装饰：删词这一发被主张 1（`:1512`）与双向对账（`:1540`）**同时**报，两枚各自单独摘掉都还有红
   （`revert-xcheck,dropword-ratify` 红在 `:1512`＝主张 1 那行；`revert-witness-rows,dropword-ratify` 红在
   移位后的 `:1512`＝对账那行，原文与 `:1540` 一字同，out/X-revert-witness-rows+dropword-ratify.txt），
   但加词那一发只有对账报（`out/U-addword-consent.txt` 交付态红 1 枚 @`:1546`）、
   M16＋摘常驻那一发只有主张 3 报 ⇒ 三味分支各有独占形状，无一枚是花瓶。
```

按简报点名的那条反面，我**没有**用它做判据：任何"同一发变异里新加那支必须先响"的要求都不成立——
乘积式挂在 `sort.Strings(hits)` 之后、控的 `||` 第三项前提正是被抹掉的东西，这类"后响"不作数。

---

## §8 ⓔ 单点回退：只回退 5 处里的 1 处，答"哪条用例变不响"

| 回退哪一处 | 变不响的那条用例（逐字名） | 凭哪一发变异量出来 | 读数 |
|---|---|---|---|
| 第 1 处：控里补齐第三因子 | `TestRealGuardRefusesEveryAssemblableApprovalRouteName` | `suffix-empty` | 交付态红 1 枚 -> 回退后 **FAIL=0**（`out/T-revert-empty-ctrl+suffix-empty.txt`，该用例由红转绿，全包无第二枚响） |
| 第 2 处：乘积式断言 | 同上那一条 | `no-plural` | 交付态红 1 枚 -> 回退后 **FAIL=0**（`out/T-revert-product+no-plural.txt`） |
| 第 3＋4 处：证人册与其断言（必须同回，单回编译不过） | `TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes` | `dropword-ratify`（另 `approve`/`permit`/`authorize`/`authorised`/`decide`/`decision`/`verdict` 同形，§6 枚枚有红） | 交付态红 1 枚 -> 回退后 **FAIL=0**＝上游 §1.2 那形原样复活（`out/X-revert-witness-mechanism+dropword-ratify.txt`） |
| 第 5 处：头段措辞 | **没有用例变不响**（注释无人读） | `revert-header` 单发 ＋叠 `MDEC-struct` | 四数一字同（`out/H-revert-header.txt`）；`:1699` 那枚红句一字未变（`out/H-revert-header+MDEC-struct.txt`） |

⇒ **没有一枚是装饰**（第 5 处按它自己的修法性质豁免这一问：它自始是文档级，判据是"措辞与盘上一致"，
我在 §9 单独验它）。第 3、4 两处编译上互为一体，按"一处"计不构成对它的不利读法，因为它俩共同负责的形状
（`dropword-*`）我实测要**两处一起回**才打得穿。

---

## §9 F-ACC-3 那格：两步路的措辞与盘上一字同（本程独立复跑两发）

```
MDEC-struct     （bridge.go 末尾加 type acProbe + 一枚 json.Unmarshal 目的地，不注册）
  RUN=97 TOPPASS=52 SUBPASS=44 TOPFAIL=1 SKIP=0 PANIC=0   唯一红名 = TestJSONKeyDerivationAgreesWithEncodingJSON
  红句 l2_grant_boundary_test.go:1699: decode destination acProbe has no reflection twin in inboundTypeRegistry: …
  原始件 out/S-MDEC-struct.txt                       ⇒ 头段"只走第一步会红在这一枚、红句是这一句"逐字成立
MDEC-registered （同上一发，再把 acProbe 写进 inboundTypeRegistry——它在本文件 :1560，不在 bridge.go）
  RUN=99 TOPPASS=53 SUBPASS=46 TOPFAIL=0 SKIP=0 PANIC=0  包内零红   原始件 out/T-MDEC-registered.txt
                                                    ⇒ 第二步确实把包带回全绿
```

**裁**：实现件 §4 那两发的四数与红名**复算同值**（它记 `RUN=97/…/FAIL=1` 与 `RUN=99/…/FAIL=0`）。
头段那句"两步、且新目的地要与注册同一枚 commit 落地"不是一句宣言：它给的第 1 步后果我打出来了、
第 2 步的出口我也走通了。F-ACC-3 这笔**文档账结掉**。
⚠ 一处台件教训写在这里（与实现件无关，是我自己的）：`inboundTypeRegistry` 住在
`l2_grant_boundary_test.go:1560`，不在 `bridge.go`；我第一次把注册那一步打进生产码，被我自己那把
"锚点计数不等于 1 即拒"的尺挡下（`REFUSED anchor MDEC-reg: count=0`），改对后才有的上面那两行读数。


