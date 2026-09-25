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
