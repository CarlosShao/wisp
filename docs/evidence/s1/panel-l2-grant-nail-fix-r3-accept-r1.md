# panel-l2-grant-nail-fix-r3-accept-r1 —— 非实现者对"门钉 r3"这一批的对抗验收

> 被验物：`docs/evidence/s1/panel-l2-grant-nail-fix-r3.md`（376 行）＋它交付的那枚常驻测试
> `internal/panel/l2_grant_boundary_test.go`（`81ad6fd` ＋ `121006d`）。任务来源：
> `docs/evidence/s1/panel-l2-grant-nail-accept-r2.md` §7.2（F-R2-3）与 §7.4 第 1 句（F-R2-1）。
> 本程＝**验收程**，一字未改代码；发现的问题只给"最小闭合集合"。

---

## §0 进场锚点、台件、以及简报前提的复算

### §0.1 锚点自量（不采信简报给的 `5ef1632`）

```
git rev-parse HEAD                     -> a5e1c8c81be138abd545add496f50c1f60f1ef9b   （2026-09-25 13:0x）
git rev-parse --abbrev-ref HEAD        -> dev
git status --porcelain -- internal/panel/   -> （空，rc=0）  ⇒ 那枚目录进场即干净，简报这句成立
git rev-parse HEAD:internal/panel/l2_grant_boundary_test.go -> 9fd6defbd8dd113a7bd9c70c619ab92ac01f32d7
git rev-parse HEAD:internal/panel/bridge.go                 -> d2cd6362ecc6941a0cee58073a2bf8ab6669a590
git log --oneline 121006d..HEAD -- internal/panel/          -> （空）
git diff --numstat a5e1c8c HEAD -- internal/panel/          -> （空）
```

**简报给的 `5ef1632` 不是进场锚点**：`git rev-list --count 5ef1632..a5e1c8c` = **4**
（`a0aaa16`、`f1cdafa`、`4037539`、`a5e1c8c`）。其中 **`f1cdafa` 直接关掉了本批红名分母里的两枚**
（前端按 `Q-50=甲` 删掉 `panel.view.request`）⇒ 简报第 4 句"另两枚红名是 renderer 那两枚"在我这一版已过期，见 §4。
（另：§8 那枚 `cmd/` 39→40 的差取自实现程 `cb60b94`→交付版之间，新增文件由 `9be3288` 引入，
那枚**早于** `5ef1632` 就已入库，所以它不影响简报锚点、只影响 §1.4 与 §8 两发的分母——§8 现量。）
⇒ 我下面凡是引用"几枚红／哪几枚红"，一律带**它自己那一发的 sha**，不写"当前 HEAD"。

共享树此刻仍在动：我这一程跑完后 HEAD 已漂到 `f50c037`（d22scan 那一发取于此后）与 `12ec7bb`（三门复算取于此后）；
上面第 6、7 行证明**漂动没有碰 `internal/panel/`**，所以被验物在我全程是同一枚 blob `9fd6defb…`。
另记一枚别家在飞的形状（非我地界、未碰）：`git status --porcelain` 里有一枚
`?? docs/evidences1142-non-quiescent-index-guard-accept-r1.md`——文件名把 `docs/evidence/s1/` 挤掉了路径分隔符，
落在仓库根；本程未动、只报备。

### §0.2 台件：仓外副本，以及简报那句"一律 git archive"在这台机器上会自己造出一枚红

```
mkdir -p /d/tmp/panel-l2-nail-fix-r3-accept-dwl/headcopy
git archive a5e1c8c | (cd …/headcopy && tar -xf -)
git hash-object …/headcopy/internal/panel/l2_grant_boundary_test.go -> 9fd6defb…   （与被验物同值）
git hash-object …/headcopy/internal/panel/bridge.go                 -> d2cd6362…   （同值）
```

**⚠ 一枚本程现量到的台件坑，必须先写**：`git archive` 在本仓对 `text=auto` 且没有 `eol=lf` 属性的文件
（`.html` 属于这一类；`.go`/`.md`/`.sh` 在 `.gitattributes` 里被钉成 lf，不受影响）**按 `core.autocrlf=true`
换成 CRLF**，于是 `internal/panel/composer_test.go` 里那枚按 `" -->\n"` 切块的
`TestComposerRenderFixtureTellsTheTruth` **在纯 archive 副本里必红**，而红因与门钉无关：

```
第一发（未归一化，out/A-baseline.txt，锚点 a5e1c8c）
  RUN=99 TOPPASS=52 SUBPASS=46 FAIL=1 SKIP=0  ^panic:=0
  唯一红名 = TestComposerRenderFixtureTellsTheTruth
  红句 = composer_test.go:699 the render fixture has no block for "no host attached" …
        :728 render fixture verified across 0 painted states          ← 切块全失败＝CRLF 指纹
凭据：git cat-file blob HEAD:frontend/fixtures/composer-states.html  -> `-->` 后是 \n
      副本里同一枚文件                 -> `-->` 后是 \r\n（od -c 现量）
      md5：实现者的副本 97c5b068043c1e293cd6bacad999b83b（LF）／我的副本 06b50261771ed270eb746c2222713fd1（CRLF）
第二发（只把那批 .html/.json 的 \r 去掉，out/A2-baseline-lf.txt）
  RUN=99 TOPPASS=53 SUBPASS=46 FAIL=0 SKIP=0  ^panic:=0  distinct=53
  ⇒ 包内零红
```

⇒ **本程之后所有"几枚红"的读数都取于归一化后的副本，并且我在正文里逐格标明是哪一发。**
这条不是实现者的缺陷，是**派单那句"变异一律在 `git archive <锚点>` 的副本里做"在本仓要配一枚 `.html` 归一化步骤**，
否则验收方会凭空多一枚红、还得给它编归因。

### §0.3 简报那五枚前提，逐枚复算（两枚对、三枚要改口径，改口径的三枚都动了我的读法）

| 简报的说法 | 本程现量 | 判 |
|---|---|---|
| 进场锚点 `5ef1632` | `git rev-parse HEAD` = **`a5e1c8c`**（漂 6 枚） | **不符**，见 §0.1，按自量值走 |
| `internal/panel/` 应当干净 | `git status --porcelain -- internal/panel/` 空，rc=0；全程 12 次复算皆空 | 对 |
| 词表 11 枚词、`grantRouteWords` 在 `:133` | 枚数 **11** 对；`:133` 是**改前那一版**的行号（`git show 81ad6fd^:… \| grep -n` = **133**），**交付版在 `:162`**（同尺现量，差 29 行＝本批往头段加的注释） | **枚数对、简报引的行号过期一格**：实现件 §0 那行量的是进场版、说得对，简报照抄时没带"取于哪一版"。⇒ 我下面一律用交付版行号，并标明取于 `121006d` |
| 另两枚红名＝renderer 那两枚，起因 `frontend/src/lib/panel.ts:252` | **在它交付那一版成立、在我锚点已被 `f1cdafa` 关掉**（§4 现量两枚皆 PASS），且新红一枚是台件假象（§0.2） | **不符（过期）**，见 §4 |
| `tools/d22scan` 是独立 module | `tools/d22scan/go.mod` 存在（`git ls-files` 现量）；`sh scripts/d22scan.sh` 内部 `runtests.sh -C tools/d22scan ./...` rc=0 | 对 |

`internal/panel/` 全程状态（每发变异之后现跑，未删任何临时件）：

```
（本程全部 12 发变异／复算结束时）git status --porcelain -- internal/panel/  -> 空
```

---

## §1 攻点 1：这枚新测试自己会不会静默空转 —— **部分会，且我造出了两形；判入账不判退回**

被验体（交付版 `121006d`，本程现读）：`grantRouteWords` `:162`（11 枚词）· `grantRoutePrefixes` `:1239`（8 枚）·
`grantRouteSuffixes` `:1247`（3 枚：`""`/`.request`/`.now`）· helper `knownComposerRefusesAssembledGrantNames`
`:1255`，体内两枚 `t.Fatalf`：`:1257` 词表空控、`:1264` 死守卫反空转控；测试本体 `:1307`。
下面每一发都是**我的仓外副本**（`git archive a5e1c8c` ＋ §0.2 的 `.html` 归一化，钉 = `9fd6defb…`）现跑
`go test ./internal/panel/ -count=1`，脚本 `scripts/mut.py`（锚点计数＝1 才落盘，落盘后打印 `delta_lines`）。

### §1.1 ⓐ "asked=0 还报绿"这一格：**前两枚清单有牙，第三枚没有**

| 发 | 我只改哪一行 | 本测试 | 全包红名（-count=1） | 原始件 |
|---|---|---|---|---|
| V2 | `grantRouteWords` → `[]string{}` | **红**，`:1305`（＝调用点，helper 的 `t.Helper()` 生效）| 3 枚：`TestRealGuardRefuses…` ＋ `TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes` ＋ `TestPlantedGrantWiringGoesRedInASnapshot` | `out/V2-words-empty.txt`、`out/FULL-V2-words-empty.txt` |
| V3 | `grantRoutePrefixes` → `[]string{}` | **红**，同一句 `the sweep vocabulary is empty (grantRouteWords=11, grantRoutePrefixes=0)` | — | `out/V3-prefix-empty.txt` |
| **V1** | **`grantRouteSuffixes` → `[]string{}`** | **绿**，日志 `asked=0 assembled route names, answeredByRealGuard=0, hits=[]` | **0 枚** | `out/V1-suffix-empty.txt`、`out/FULL-V1-suffix-empty.txt` |
| V4 | 复数那一圈 `[]string{w, w+"s"}` → `[]string{w}` | 绿，`asked=264` | 0 枚 | `out/V4-no-plural.txt` |
| V5 | 问守卫那一句改成 `if false && knownComposerMethod(name)` | 绿，`asked=528 hits=[]`——**与干净树同一枚读数** | 0 枚 | `out/V5-guard-stubbed.txt` |
| V8 | `hits = append(…)` 那一支短路 | 绿，`asked=528` | 0 枚 | `out/V8-assert-false.txt` |

⇒ **简报那句"asked=0 还报绿＝尺不存在"在这一格成立**：`:1257` 那枚反空转控逐字只写了
`len(grantRouteWords) == 0 || len(grantRoutePrefixes) == 0`，**乘积的第三枚因子没人问**，
所以把 `grantRouteSuffixes` 掏空 ⇒ 这枚测试照问 0 枚、照报绿，且**全包没有任何别的测试变红**
（现量：`grantRouteSuffixes` 在包内的代码引用只有 `:1273` 那一处 `range`，其余全是注释）。
V4/V5/V8 是同一族的更广义形状：`asked=` 那行只统计"循环走了多少趟"，没有任何断言把它钉在
`len(prefixes)×len(words)×len(suffixes)×2` 上，所以分母减半、甚至把守卫那一问短路掉，
报出来的都还是同一句"asked=528 / hits=\[\]"。**V5/V8 我不记入账**——那是"任何断言都能被 `if false` 短路"
这一类，任何 Go 测试都治不了；**V1 记入账**，因为它不是短路而是**它自己那枚防呆少写了一枚变量**，
而同一枚 helper 里另外两枚因子已经被防住了（形状就摆在自己三行之上）。

### §1.2 ⓐ 追加大获：词表**内容**只剩 4/11 枚有人质，删掉第 5 枚就放走 M16

`TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes` `:1400` 用的是"逐枚点名几枚路由名，它们必须仍然读成批准门"
这一招——这正是本仓对付"词表被人悄悄改窄"的既有武器。本程按枚复算它到底覆盖了几枚词
（`scripts/routevocab.py`，只重排 `var grantRouteWords` 那一块，一次删一枚）：

| 删掉的词 | 全包读数 | 判 |
|---|---|---|
| `allow` | **红 1 枚**：`TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes`（`:1402` 那句 "reads as an innocent name"）| 有人质 |
| `approval` / `decide` / `grant` | 同一枚证人测试点名这四枚路由（`panel.approval.request` / `panel.l2.allow` / `approval.decide` / `panel.grant`）| 有人质（`allow` 那发＝现量，其余三枚＝同一圈代码，未逐枚打） |
| **`ratify`** | **0 枚红**；再叠上**被验物自己那发 M16**（生产码运行期拼 `panel.review.ratify` 并被守卫回答）⇒ **全包 53 枚顶层全绿、`asked=480 hits=[]`** | **没人质** |
| `permit` / `verdict` / `approve` / `authorised` | 各 0 枚红 | 没人质（现量四枚）；`authorize` / `decision` 同形状**未逐枚打＝推定** |

原始件：`out/RV-ratify-M16.txt`（那一发既是"删词"又是"M16"，是全包绿的那一发）、`out/RV-allow.txt`、
`out/RV-permit.txt`、`out/RV-verdict.txt`、`out/RV-approve.txt`、`out/RV-authorised.txt`。

⇒ **实现件 §2.2 那行 T-D（"抓得住 M16"）在交付态成立，但它承重的不只是那三行代码，还承着 11 枚词一枚都不能少；
而"少一枚"这件事今天只有 4 枚词有人质。** 记入账（F-ACC-2），最小闭合在 §9。
⚠ 这不是"退回"：删词是**编辑这把尺自己**，被验物从头到尾没声称它防得住人把词表改窄（`:91` 那句
"widening them is a slice card rather than a quiet edit here" 说的正是这一族的处置权），
而**生产码形状**（M14/M16 那两发）在交付态确实被抓得住——那是 §3 的独立复算，不是我推的。

### §1.3 ⓑ 那两枚锚点／长度守卫是不是真在——我把它自己的摘除逻辑接过来打台件

我端口径复用它的 `plant.py nail standing`（同样两枚锚点、同样"切完函数名不许还在"、同样打印
`delta_lines`），然后**故意把台件弄坏**看它拒不拒（原始件 `out/`，读数在下面这段里逐字贴）：

```
P1 backup 里那枚起始锚点出现两次  -> rc=1  removal: anchor count=2 expected 1        ← 拒
P2 backup 里收尾锚点被改写（计数 0）-> rc=1  removal: anchor count=0 expected 1        ← 拒
P4 backup 里只少一枚闭括号          -> rc=1  removal: anchor count=0 expected 1        ← 拒
P3 backup 里函数被改名、锚点完好     -> rc=0  removal: delta_lines=-32, textual mentions left=5   ← 不拒
```

⇒ "零变化即拒绝报告"那句**是真的**（三发坏台件全被 `anchor count` 挡下，报的是退出码 1 而不是"-32"）。
**P3 那一形它挡不住，而它本来也不声称挡**：守卫验的是"我这刀落下去了"，不是"那枚禁令还活着"。
另——它的脚本是**从 `backup/` 整枚重写副本**的，所以 §2.4 那次 row U 的自曝（"plantf 塞回去的三行被
nail standing 又抹掉，那一发实测等同 T-E"）**在这台机器上是机械必然**，我复算了它的写文件路径
（`open(NAIL_OUT).write(from_backup.replace(…))`）才敢这样判；那一发自曝本身是加分不是减分。

### §1.4 ⓒ 词表被清空／被改名，有没有**另一枚**测试变红

| 形状 | 现量 | 判 |
|---|---|---|
| `grantRouteWords` 清空 | 全包 **3 枚红**（本测试 ＋ §1.1 那两枚）| **不是"存活由它自己守"**，这一格它说得过去 |
| `grantRouteWords` 枚数 11→10（删 `ratify`） | 全包 **0 枚红**（§1.2） | **只有它自己守，而它自己没守** |
| `grantRouteWords` 整体改名（18 处一起改） | `ok github.com/CarlosShao/wisp/internal/panel`（**rc=0，包全绿**）| 改名＝无害重构，不构成空转；**真正可乘的是删内容** |
| `grantRouteSuffixes` 清空 | 全包 **0 枚红**（§1.1 V1） | **只有它自己守，而它自己没守** |

⇒ 简报那句"点名'这把尺的存活由它自己守'是不是够"——我的裁法：**枚数级不够、清单级不够、清空级（两枚因子）够**。
补齐只需要两行，见 §9 的 F-ACC-1／F-ACC-2。

---

## §2 攻点 2：§2.4 选"删"那一支的反噬 —— **判"挪"不判"删"；独占性的三种暴露里只有一种真裸**

### §2.1 ⓐ 删掉那三行之后，植物 F 那枚 subtest 还剩什么在主张

按枚数（脚本 `git show <rev>:…` 取块，块＝`t.Run("F a route named outside the guard` 到其 `}`）：

```
A 版（81ad6fd^，进场版）：块内 24 行、断言 5 枚  {t.Fatal, t.Fatalf, t.Error, t.Errorf}
B 版（121006d，交付版）：块内 21 行、断言 4 枚  {t.Fatal, t.Fatalf, t.Errorf}
少的那一枚＝被搬走的那句 t.Error（与 §6 全局尺 t.Error 2 -> 1 同一枚，不是第二枚）
```

剩下四枚主张的仍是"pool 那半边"，逐枚点名（交付版行号）：`:1884` 副本 parse 不掉就 `t.Fatalf`、
`:1886` `requireReadableInstrument`（读不动树就响亮失败）、`:1889` **pool 没在别的函数里看见这枚名字就 `t.Fatal`**、
`:1892` pool 必须报出在哪个文件找到它（`t.Errorf`）、`:1895` 快照里守卫自己的 case list **不许**把这枚植进去的链
解析成已答（`t.Errorf`）。⇒ **"pool 那半边"照旧被钉住，没有因为删三行而少问一个形状。**

**并且没有任何形状从此没人问**：被删那句问的是"真守卫不许答 `panel.review.allow`"，而这枚名字
**落在新测试那张网里**（`panel.review` × `allow` × `""` × 单数）——不是我推的，是我在**自己的副本**里
把 M14 打进去现量到的（§3 的 T-B 那发 `hits=[panel.review.allow]`）。
反向也成立：交付版在 M14 树下，植物 F 那枚 subtest **照旧 PASS**
（`out/T-B-M14.txt:247` `--- PASS: TestPlantedGrantWiringGoesRedInASnapshot/F_a_route_named_outside_the_guard_still_reaches_the_pool`，
同发 A/B/C/D/E/G/H 七枚 subtest 也全 PASS）⇒ 搬走那半句没有把快照那台仪器一起拖走。

### §2.2 ⓑ 禁令独占在新测试里 ⇒ "改名 / t.Skip / 整段摘掉"三种情况各有没有人响

| 情况 | 我怎么打 | 现量 | 谁响 |
|---|---|---|---|
| **整段摘掉** | `scripts/mut.py standing-out`（helper 故意留在原地），再叠 M14／M16 | `out/T-C2-out-M14.txt`、`out/T-E-standingout-M16.txt`：`FAIL=0`、**全包包内一声不响** | **只有被摘那枚**。⇒ 承重成立（这是本格的正解），但见下两行的不对称 |
| **被改名** | `func TestRealGuard…` → `func notARealGuardRefuses`（锚点全留，编译照过） | `out/V6-rename-clean.txt`：RUN 99→98、红 0；**再叠 M14**：`out/V6-rename-M14.txt` 仍 0 枚红 | **没人响**。包内没有第二枚测试钉"这枚用例存在"；`grantRouteSuffixes` 之外，我把 `scripts/`、`tools/`、`.github/workflows/` 全 grep 过一遍，**没有任何一道门把顶层用例名册钉成常量**（`portable-tests.sh` 的 census 只数每包 PASS/FAIL 枚数，−1 那枚混在里面看不见） |
| **被 `t.Skip`** | 测试第一行插 `t.Skip(…)` | `out/V7-skip-clean.txt`：SKIP=1，而 **`go test ./internal/panel/` 退出码 = 0**（Go 把 SKIP 记成 ok）；换成 CI 那把严尺 `bash tools/d22scan/runtests.sh ./internal/panel/ -run TestRealGuard…` ⇒ **rc=1**，原文 `runtests.sh: 1 test(s) SKIPPED and SKIP is not a pass (ticket 71 AC#3)`；**同一道命令在干净钉上 rc=0**（正控，`out/T-A-runtests.txt`） | **CI 那一步响**（`scripts/portable-tests.sh:80` 把 verdict 交给 `runtests.sh`，`--scope=core` 的名单里逐字写着 `./internal/panel/...`，同文件 `:179`）。⇒ 这一形不是裸的 |

⇒ 三种里只有一种真裸，而那一种是**全仓通用形状**（任何一枚 Go 测试被改名都不会有人响；本仓至今是靠
"名册差集 +1/−0"这一枚人工动作在防，实现件 §1.2 自己就做了那一步）。它**不是本批独有的洞**，
所以我不拿它当退回理由；但**它使 §2.4 那句"删掉才量得出来"变成一条必须长期保留的账**：
这枚新测试一旦被改名，`out/V6-rename-M14.txt` 那种"守卫真答了 panel.review.allow 而全包绿"的读数会**静默回来**。
另记一枚附带的自指残留：交付版里那枚名字被 6 处文字引用（`:47`/`:95`/`:1250`/`:1288`/`:1865`/`:1898`，
最后一处在 `t.Logf` 的**字符串**里），改名之后这 6 处全部指向一枚不存在的测试而无任何人响——
形状与验收件 r2 §7.3 记过的那枚"`bridge.go:33` 指向不存在的用例"同类，属读数/文档级。

---

## §3 攻点 3：零误伤与抓得住，**独立重跑并对账**（我自己的副本目录，不复用它的）

台件：`D:\tmp\panel-l2-nail-fix-r3-accept-dwl\headcopy`（`git archive a5e1c8c` ＋ §0.2 归一化），
钉 `9fd6defbd8dd113a7bd9c70c619ab92ac01f32d7` ＝ **它交付的那枚 blob**（`git rev-parse HEAD:…` 现量，§0.1），
`bridge.go = d2cd6362…`（同它的十二次读数）。全部 `-count=1 -v`。

| 发 | 我这边的读数 | 它的读数（实现件 §2.2） | 对账 |
|---|---|---|---|
| **T-A** 交付态＋干净守卫 | `FAIL=0`，`asked=528 assembled route names, answeredByRealGuard=0, hits=[]`（`out/A2-baseline-lf.txt:178`、`out/T-A-clean.txt`）| `asked=528 / 0 / hits=[]` → PASS | **逐字同值** |
| **T-B** ＋M14（运行期拼 `panel.review.allow`） | 红 **1 枚**＝`TestRealGuardRefusesEveryAssemblableApprovalRouteName`；红句两行：`:1310 … answers "panel.review.allow" …` 与 `:1315 … plant F plants is now answered by the guard …`（`out/T-B-M14.txt:178-179`）| 同测试红，`hits=[panel.review.allow]` | **同值**（那枚"共红诊断"我也量到，同发同一处） |
| **T-D** ＋M16（运行期拼 `panel.review.ratify`） | 红 **1 枚**＝同一枚；`asked=528 … answeredByRealGuard=1, hits=[panel.review.ratify]`（`out/T-D-M16.txt:178-179`）| 同值 | **同值** |
| **T-G** ＋M17（死守卫控） | 红 **4 枚**：`TestComposerEnvelopeAcceptsItsFourRequests`、`TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests`、`TestAnsweredPanelRoutesCarryNoApprovalDecision`、`TestRealGuardRefuses…`；本测试的红因**不是** hits，而是那句 `the running guard refuses its own declared route "panel.mode.request" …`，**报在 `:1308`**（＝调用点，`t.Helper()` 生效）（`out/T-G-M17.txt:192`）| 它记"红因是 `:1266` 那枚反空转 `t.Fatalf`，读数报在调用点 `:1308`"，同发还红 `TestAnsweredPanelRoutes…`/`TestComposerEnvelope…` | **同值**，含"`t.Helper()` 把行号搬到调用点"这一枚细节：我把词表掏空那一发（少 3 行）里同一句控报在 **`:1305`**（`out/V2-words-empty.txt`）⇒ 调用点跟着文件漂，它这句描述是活的不是抄的 |
| 加分行 前缀 8→14 | `asked=924 … hits=[]` → PASS（`out/W-widen14.txt`）| 924 / 0 | **同值**，分母可复算＝`14×11×3×2` |

**分母对账（不许"差不多"）**：`grantRoutePrefixes` 现量 **8 枚**（`:1239-1242`）、`grantRouteWords` **11 枚**
（`:162-165`）、`grantRouteSuffixes` **3 枚**（`:1247`）、单复数 ×2 ⇒ `8×11×3×2 = 528`。
⇒ 我这枚 528 与它那枚 528 是**同一版词表的同一张网**（钉 blob 同值，前缀/词/后缀三枚清单逐枚同）。
**没有"前缀集不同"或"词表版本不同"这类差**；唯一两处不同都 attributable 到台件而非尺，且我在别处归了因：
① 先存在的红枚数（我这里 0 枚、它那里 3 枚）→ §4；② 我那发 M14/M16 的**植入文本**与它不同（我 7 行、它记 `+352B/9 行`）
——都是"三段字面量 `strings.Join` 拼出一枚 map，守卫末尾返回 `map[m]`"，行为同形，文本不必同。

**"asked=528"这枚数今天仍然只是 `t.Logf` 里的一个读数，没有任何断言钉它**（§1.1 的 V4/V1 就是这么打穿的）。
实现件 §2.1 那句"乘积写死在 `t.Logf` 里，人事后可复算"——**可复算这句对**（清单就在 `:1239/:162/:1247`，
我按枚数过），但请把"写死"两字读成"打印"，别读成"断言"。

---

## §4 攻点 4：那句"验收件说只剩 C21、今天要说只剩三枚" —— **判成立（连归因一起复算），且形状结论我这边更强**

### §4.1 三枚红名的归因，逐枚独立复算

**ⓐ 那两枚 renderer 红：我在它交付那一版（`121006d`）的纯净副本上重新造出来。**

```
git archive 121006d -> /d/tmp/…/headcopy-r3anchor（同样只归一化 .html）
go test ./internal/panel/ -count=1 -v   -> FAIL=2：
    TestTheRendererHoldsExactlyOneDoorToTheHost
    TestPlantedRendererDoorShapesGoRed
红句原文（out/S-r3anchor-121006d.txt:152/156）：
    src/lib/panel.ts:252:   sendRequest("panel.view.request", …)  names "panel.view.request"
```

⇒ 实现件 §1.3 那三行归因（同一枚红因行 `:252`、`git log -1 -S'panel.view.request'` 指到 `d61281c`、
两枚用例都定义在 `internal/panel/composer_test.go`、"不在本程地界、未修未跳未放宽"）**逐条对得上**；
它引的台账出处我也核了：`docs/reports/pending-and-issues.md:6124`（`A237②`，那一句正是"那道门被前端
会话那枚第六个方法名打响"）与 `:1084`（`Q-50`）都是真条目，不是它编的指针。

**ⓑ C21 那枚：它的"今天的机制是脏树，不是内容缺陷"这句，我用台件差集证明，不是听来的。**

```
同一枚钉、同一次归档，只差 design/ 在不在：
  我的 git archive 副本（design/assets/tokens.css 在，HEAD 里跟踪着）
      -> TestC21DesignTokensFourWayAgree **PASS**（out/A2-baseline-lf.txt / out/S-r3anchor-121006d.txt 两版皆 PASS）
  仓库工作树（design/ 那 16 枚未提交删除，别家在飞）
      -> 同一枚用例 **FAIL**（out/R-repo-head.txt，锚点 a5e1c8c，红句 read design/assets/tokens.css: … cannot find the path）
```

⇒ 它 §1.3 末段那句"编排者若按'C21 两枚红'核账会数不上，特此报备"**报备得对**，而且它把机制
（脏树）与内容缺陷分开写、没有顺手把 C21 记成自己那批发现的、也没修它（不在地界）——这是加分。

### §4.2 简报第 4 句（"另两枚应是 renderer 那两枚"）**在我这一版已过期**，与它无关

```
git show HEAD:frontend/src/lib/panel.ts | grep -n 'view\.request'
  -> 只剩 :228 那行注释（"There is deliberately NO fifth outbound route here…"），sendRequest 的实参只剩四枚
git log --oneline 5ef1632..a5e1c8c  -> f1cdafa fix(前端·Q-50=甲): 删掉第五枚出站路由 panel.view.request
```

⇒ **`f1cdafa` 已经落了 owner 那句"删掉那个请求"**，于是在我的锚点 `a5e1c8c` 上，归一化后的纯净副本
**`FAIL=0`、包内一枚红都没有**（§0.2 第二发）。⇒ 简报那句前提在 `121006d`/`5ef1632` 上成立、在 `a5e1c8c` 上作废。
**编排者若现在去复核"三门"，四数与实现件 §1.2/§7 那两发必然不同**（我这里 `count=1` 是
`RUN=99 TOPPASS=53 SUBPASS=46 FAIL=0 SKIP=0`，它那里 `count=2` 是 `RUN=198 TOPPASS=100 FAIL=6`），
差别**全部**由这两枚别家 commit（`f1cdafa` 关红、`design/` 脏树决定 C21 那枚的红/绿）解释，
与被验物无关——这句话必须连着锚点 `a5e1c8c` 一起引用，别只引数字。

### §4.3 "形状结论不变"那句是否诚实 —— **诚实，而且我把它推到更强的那一头**

它写的是："⚠ 验收件 §7.2 那句'只剩先存在的 C21 红'在今天要读作'只剩先存在的**三枚**红'，**形状结论不变**"。
我这边的复算不必借用它的分母：在我的纯净副本上，**先存在的红是 0 枚**，于是

```
摘掉常驻测试 ＋ M14  -> 包内 0 枚红（out/T-C2-out-M14.txt）
摘掉常驻测试 ＋ M16  -> 包内 0 枚红（out/T-E-standingout-M16.txt）
交付态（不摘） ＋ M14/M16 -> 各自恰好 1 枚红＝那枚常驻测试（§3 T-B/T-D）
```

⇒ 形状结论不但"不变"，在干净台件上是它的最强形：**没有那枚新测试，Go 真答一条面板侧批准路由时，
整个 `internal/panel` 包会一声不响地全绿**。它那三枚先存在的红反而**稀释**了它自己的读数（读起来像"还有别人兜"），
我这里连那点稀释都没有。它没有把这件事说重，也没有借我的口径改自己的数——按本仓"量出与上游不同就照自己读数写"
那条，这一格判**成立**。

---

## §5 攻点 5：F-R2-1 那一格（文件头 `:102` 起那段）—— **自纠后的"四枚"判对；那句出路要补一步**

### §5.1 我自己造的那发 MDEC：红名枚数与逐枚名册

```
scripts/mut.py MDEC   -> 副本 bridge.go 里加一枚合法的 var acProbeMap map[string]any + json.Unmarshal
go test ./internal/panel/ -count=1 -v   （out/T-V-MDEC.txt）
FAIL=4：TestAnsweredPanelRoutesCarryNoApprovalDecision   （facet 1）
        TestNoInboundEnvelopeCanBindAnApprovalVerdict    （facet 2）
        TestJSONKeyDerivationAgreesWithEncodingJSON      （三向对照）
        TestPlantedGrantWiringGoesRedInASnapshot         （facet 4 的 pristine-snapshot 前置）
四处同一句，报在 :1162 / :1351 / :1554 / :1740 —— 与实现件 §3 引的四枚行号**逐字同值**
红句：a JSON decode destination in 1 place(s) cannot be enumerated: bridge.go:110:
      decodes into &acProbeMap (type "map[string]any"), which is not a same-package struct
```

⇒ **`121006d` 那次"三枚 → 四枚"自纠判对**；验收件 r2 §7.4 第 1 句原文的"三枚 ban 测试 `t.Fatalf`"确实是错的
（错在把 facet 4 那枚快照前置漏计了，它走的是同一道 `requireReadableInstrument`）。
它 §3 末段那句"本程那枚常驻测试在此发仍绿（它不读 decode）"我也复现：`out/T-V-MDEC.txt:179` `--- PASS: TestRealGuard…`。
⚠ 读数口径提醒：这一发 `=== RUN` 从 99 掉到 **89**（四枚父用例 `t.Fatalf` 后它们体内 10 枚子测试不再跑），
`^panic:` = 0 —— **MDEC 那一行的四数不能和干净态那行的四数直接比**，别下游谁把它当"少了 10 枚用例"来记账。

**另两发我把它声称的两种目的地都打了**（头段原话是 "a decode into a `map[string]any` **or** a `json.RawMessage`"）：

```
MDEC2（目的地换成 json.RawMessage）-> 同一批 4 枚红、同四处行号、红句 … decodes into &acRaw (type "json.RawMessage")
                                      （out/T-V2-MDEC-raw.txt）  ⇒ 头段那句"两种合法形状都算"是真的
```

### §5.2 那条约束写进文件头的措辞**可证伪**吗 —— 可，且我按它的句子造了反例方向

头段那句（交付版 `:104-111`）的主语是**一形代码**（"a SECOND decode destination in `internal/panel`"）、
谓语是一枚**可数的后果**（"makes **four** of this file's tests abort with `<逐字引出的消息>`"）——
两样我都在副本上量到了，所以它不是宣言。**它没有写成"必须怎样才对"那种无从否证的句子。**

**但它的出路那半句少写了一步**（记入账 F-ACC-3）：头段说 `the message carries its own way out, which is to
route the bytes through a same-package struct **or** to register the type in inboundTypeRegistry …`。现量：

```
只走第一路（新增同包 struct 作目的地）      -> FAIL=1：TestJSONKeyDerivationAgreesWithEncodingJSON
      红句 :1560 decode destination acProbe has no reflection twin in inboundTypeRegistry
                                              （out/T-V3-MDEC-struct.txt，RUN=97）
两路一起走（同包 struct + 注册进 inboundTypeRegistry） -> FAIL=0，且那枚类型真被两把尺对着判：
      :1576 acProbe: 1 keys agreed by both instruments [text]（out/T-V4-MDEC-registered.txt，RUN=99）
```

⇒ 出路**存在**（我按它给的材料走通了，零红），但它是一条**两步**的路，头段那个 "or" 会让下一个撞上来的人
以为"换成 struct 就完事"，然后在第二道引信上多吃一枚红。修法是一行措辞：把 "or" 改成
"and then"（或补一句"新目的地进了包就得进 `inboundTypeRegistry`，`:1560` 那枚引信会替你说"）。
这一条**与"三枚/四枚"那处自纠同级**，属文档级、不改判据、不放宽断言。

---

## §6 攻点 6：门禁／名册／契约轴／断言尺 —— **四门全过，名册 +1/−0，契约轴零命中且正控活**

### §6.1 三门（都在仓库工作树上跑，`internal/panel/` 全程干净）

```
gofmt -l internal/panel/     -> （空）rc=0
go vet  ./internal/panel/    -> （空）rc=0
sh scripts/d22scan.sh        -> rc=0（原文 out/d22scan-head.txt，取于 HEAD=f50c037）
  d22scan: examined 225 production Go files under internal/ and cmd/
  bans #1-5 internal/=203  cmd/=22 | ban #6 frontend/=46 | ban #7 internal/tools/=18
  ban #8 design/=32  frontend/=46  internal/=407  cmd/=40
  d22scan: clean - no D22 ban violations
  runtests.sh: OK - packages=[./...] top-level: PASS=30 FAIL=0 SKIP=0, === RUN=70
```

⚠ 简报那句"`tools/d22scan` 是独立 module、别在根目录 `go vet ./tools/d22scan/`"**复算对**：
`tools/d22scan/go.mod` 在 `git ls-files` 里；本程没有对它下过任何一条根目录尺。
⚠ 简报那句"bans #1-5 不含 `_test.go`，只有 ban #8 含"也**复算对，我用文件枚数证的**：
`git ls-files internal/ | grep '\.go$' | grep -vc '_test.go$'` = **203**（＝bans #1-5 的 internal/ 分母），
`git ls-files internal/ | grep -c '\.go$'` = **407**（＝ban #8 的 internal/ 分母）。
⇒ 本批改的恰好是一枚 `_test.go`，**这次改动只进 ban #8 那一枚分母**，实现件 §1.4 那句口径正确。
另：`scripts/d22scan.sh:14-16` 逐字写明它的第一步＝`runtests.sh -C tools/d22scan ./...` 是**植入违规的正控**
（"proves the gate CAN go red"），所以我这一发的 rc=0 不是"门是瞎的"那种绿。

### §6.2 第四门：包体四数 ＋ 名册差集（两种口径都取运行输出，不取 `t.Skip`）

| 台件 | 四数（`-count=1 -v`） | 红名 |
|---|---|---|
| 仓库工作树，锚点 `a5e1c8c`（`out/R-repo-head.txt`）| `RUN=99 TOPPASS=52 SUBPASS=46 FAIL=1 SKIP=0 ^panic:=0 distinct=99` | `TestC21DesignTokensFourWayAgree`（脏树机制，§4.1ⓑ） |
| 我的 archive 副本（同钉、归一化，`out/A2-baseline-lf.txt`）| `RUN=99 TOPPASS=53 SUBPASS=46 FAIL=0 SKIP=0 ^panic:=0` | 无 |
| 它交付那一版 `121006d` 的 archive 副本（`out/S-r3anchor-121006d.txt`）| 同上口径 | 2 枚 renderer（§4.1ⓐ）|

```
名册差集（简报要的 ^func Test 逐名对比，两枚版本都取 git show）：
  单文件：81ad6fd^ 6 枚 -> 121006d 7 枚      diff = 7a8 > TestRealGuardRefusesEveryAssemblableApprovalRouteName
  整包  ：git grep -h '^func Test' <rev> -- internal/panel/  -> 52 枚 -> 53 枚，diff 同上**一行**、删除侧**空**
```

**panic 那一格按简报要求明写**：两发都 `^panic:` = **0**，所以**没有**"一条用例 panic 吞掉同包其余几十条"的
情况发生，上面那 99 条 `=== RUN` 是真取到的；我也**没有**用 `-run`/`-skip` 排除任何用例（`-run` 只在
§1.3 那三发单点复算与 §2.2 的 `runtests.sh` 那两发里出现，且那两发的完整包体读数在同表另有整跑）。
⇒ 实现件 §1.2 那句"名册 +1/−0、FAIL 改前后同值"**成立**（它那里 6→6 是 count=2、三枚红 ×2；我这里 1→1 是 count=1，
换算同值）。

### §6.3 契约轴（先证存在，再说零命中，同尺带正控）

```
本批自己的落点（逐枚 commit 现量，git show --numstat）：
  81ad6fd  130  14  internal/panel/l2_grant_boundary_test.go        （show --name-only 只有这一枚路径）
  121006d    3   2  internal/panel/l2_grant_boundary_test.go        （同上）
  合并 81ad6fd^..121006d，把保护面一次性写进 pathspec：
    git diff --numstat 81ad6fd^ 121006d -- internal/ frontend/ design/ cmd/ docs/PLAN.md docs/specs/ \
        internal/risk/ internal/observe/thresholds.go internal/risk/rules_gateway.go \
        tools/d22scan/allowlist.txt '*.sse' '*.go'
      -> 只有一行： 131 14  internal/panel/l2_grant_boundary_test.go
      ⇒ **正控就是这一行**（同一把尺、同一区间、指自己那枚文件，非空）
先证存在（不是"找不到就当零命中"）：
  git ls-files | grep -icE 'thresholds\.go$|allowlist\.txt$|rules_gateway\.go$'  -> 3（三枚逐名都在）
  git ls-files | grep -c '\.sse$'                                                -> 52 枚 golden
  git ls-files frontend | wc -l -> 46   design -> 30   internal/risk -> 37   cmd/wisp -> 35   docs/specs -> 14
```

⇒ `cmd/wisp/**`（143 那程在写）**存在且本批零命中**，`git status --porcelain -- cmd/wisp/` 亦空（它的活已入库）。
**生产码零字节**：`internal/panel/bridge.go` 在我副本里每发变异之后复算，
`git hash-object` 全程 `d2cd6362…`（副本被 `restore` 之后再次现量仍为该值）。

### §6.4 断言尺（我自己的正则，两版都从 `git show` 取，不取工作树）

```
                A(81ad6fd^)  B(121006d)  delta
t.Fatalf(            31          33        +2
t.Fatal(              5           5        +0
t.Errorf(            38          40        +2
t.Error(              2           1        -1     <- 唯一那枚减少
t.Logf(              14          15        +1
t.Skip                0           0        +0
func Test(            6           7        +1
assertions(总)       76          79        +3
行数               1964        2081
```

⇒ 与实现件 §1.5 **逐格同值**（我原本用一枚写坏的正则数过一次 `t.Fatalf`＝0，发现是自己的式子错了，
换成 `t\.Fatalf\(` 才是上面这行——**这条自纠写在这里免得下游把这格当成对它的复算不符**）。
**那 3 行是"挪"不是"删"**，凭三条现量：① 全局减少的恰好只有 `t.Error` 一枚，枚数与 §2.1 植物 F 块内
`5 → 4` 对得上；② `81ad6fd` 的 14 枚删除行我逐枚看过（10 枚注释 ＋ 1 枚 `t.Logf` 改写 ＋ 3 枚那句
`if knownComposerMethod("panel.review.allow") / t.Error / }`），**它 §1.5 那句"11＋3＝14、没有第 15 枚删除"
逐字对**；③ 那句谓词在交付态仍然会红，且红在新测试 `:1310`——我在自己的副本上把 M14 打进去量到了
（§3 T-B）。**判：挪判据，可以**。

---

## §7 攻点 7：§2.5"一处没做的事"与 §4"没测什么"是否诚实 —— **诚实；那处"没做"不构成退回，但 §4 少了一条我该说的**

### §7.1 可核的那几条，我逐条核到盘上

| 它的自陈 | 我怎么核 | 判 |
|---|---|---|
| "只写两枚路径"／"未碰 `internal/panel/` 下别的文件、`tools/d22scan/**`、`frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`internal/risk/**`、`rules_gateway.go`、`allowlist.txt`、`thresholds.go`、任何 golden、台账、票面、`Q-49` 的勾 | 三枚 commit 的 pathset 现量：`81ad6fd` → 只有 `internal/panel/l2_grant_boundary_test.go`；`121006d` → 同一枚；`062d389` → 只有那枚证据件。区间内碰台账／票面的 commit（`7d6a019`/`6103120`/`97de647`/…）**作者不是我核的这一程** | **对**（路径级可证） |
| "临时件只建不删（未 `rm`/`rmdir`）" | 它的 `D:\tmp\panel-l2-nail-fix-r3\` 里 `backup/` 2 枚、`scripts/` 8 枚、`out/` **55** 枚全在 | **对**，但枚数见下表 §7.3 |
| "没跑 `go test ./...`" | 它 55 枚原始件里出现的包只有两枚：`internal/panel` 与 `tools/d22scan`（后者是 `scripts/d22scan.sh` 自带的正控测试，不是它跑的包体扫描） | **与盘上无矛盾**（这条只能证"没留证据说跑了"，不是硬证） |
| "没给自己判成立" | §0 第 8 行与 §4 末条都写着"承重结论要另一程复算" | **对**，本程就是那一程 |
| §2.4 那发台件失败（row U 等同 T-E） | 我核了它的脚本：`nail standing` 是**从 `backup/` 整枚重写**副本 ⇒ 上一步塞回植物 F 的三行必被抹掉；它 `out/U-…txt` 里 `grep -c 'no longer a clean tree'` = **0**，`out/G-…txt` = **1**（现量同值） | **自曝得对，且方向是往不利于自己的方向自曝**（本程加分项） |

### §7.2 那处"没做"（§2.5：`widen.py` 8→14 前缀未交付）按本仓分界怎么裁

分界摆在这里：**验收方自己造得出"被防的结局" ⇒ 退回**，否则才是"附条件"。我照这句打：

```
被验物声称防的结局 = "守卫回答了一枚运行期拼得出的批准路由而没人红"
在交付网格（8×11×3×2）内：M14 / M16 各恰好红 1 枚 = 那枚常驻测试（§3，我自己的副本现量）
网格之外：wisp.review.ok / composer.allow 之类没问，且文件头 :104-100 逐字写着
          "the hole that stays open is the grid's outside: another namespace, another vocabulary,
           or a route that only exists once config is read"
```

⇒ **我造不出"声称要防而没防住"的结局**（在声称范围内它防住了，两向都防住了）。
把 924 那发留在副本不外发，是"两程账对得上"这一枚取舍，代价是网格窄 6 枚命名空间——
**这属它已经写明的"残余"，不属"做过了却没做"**。⇒ **§2.5 不构成退回理由**。

### §7.3 §4 那节少写的一条（我替他补上，免得下一程以为"没测"就等于"没人发现"）

§4 列了"没测词表外的结论键／没测运行期从 config 来的名字／没测网格之外的名字"，**方向都对**；
但它没列 **"没测这把尺自己会不会少问"**——也就是本程 §1 那两发（`grantRouteSuffixes` 掏空 ⇒ `asked=0` 仍绿；
词表 11 枚里删任意一枚无人质 ⇒ 该词覆盖的 48 枚候选静默消失）。实现件 §2.1 客观描述了
"内含两枚 `t.Fatalf`：词表为空、死守卫控"，所以**没有说谎**；但"两枚"读起来像"自空转已被堵完"，
实际是"三枚乘数里堵了两枚"。⇒ 这一条记在我这里，别记成它隐瞒了什么。

另记一枚读数级漂移（同 r2 记它上游的 F-R2-2 同一族，**不新造账**）：

```
§4 写 "out/ 50 枚读数与名册"  ->  现量 ls D:\tmp\panel-l2-nail-fix-r3\out | wc -l = 55
差在哪：按 mtime 落在 12:46 之后的 7 枚里，committed 态那批发（committed-count2.txt(.roster)、
        d22scan-committed.txt、postcommit-count2.txt(.roster)）＝ 5 枚，正是 §7/§8 两条追加所产
判：写"50"的当时是对的，§7/§8 之后没回来改 —— 形状与它自己在 §6 拒给 commit 3 预写 sha 是同一条规律的两侧
```

---

## §8 攻点 8：`ban #8 cmd/ 39→40` 那处分母时效 —— **差值与归因判对；"带 sha"这条它自己没做到**

### §8.1 我这一版复跑到同值（两次，各自带当场 HEAD）

```
HEAD=f50c037  sh scripts/d22scan.sh -> rc=0   （out/d22scan-head.txt）
HEAD=0e17914  sh scripts/d22scan.sh -> rc=0   （out/d22scan-0e17914.txt，即本件 §7 那枚 commit 入树之后）
两发八枚数一字同：bans #1-5 internal/=203 cmd/=22 | ban #6 frontend/=46 | ban #7 internal/tools/=18
                    ban #8 design/=32 frontend/=46 internal/=407 cmd/=40
⇒ 与实现件 §8（交付版钉 9fd6defb／2081 行）那一发**逐枚同值**；§1.4 那一发唯一的差别就是 cmd/ 39。
```

### §8.2 差值归因（它只说"别家在 `cmd/` 里多提了一枚文件"，我把它没点的名点出来）

```
git ls-tree -r --name-only <rev> cmd/ | grep '\.go$' | wc -l
   cb60b94 -> 39      （§1.4 那一发的锚点）
   1785a77 -> 39      （实现程进场）
   9be3288 -> 40      ← 分母在这一枚掉头
   a5e1c8c -> 40      （我的锚点）    HEAD(0e17914) -> 40
git log --format="%h %s" --name-status --diff-filter=A cb60b94..a5e1c8c -- cmd/
   9be3288 feat(票143 AC#1+AC#2): panel-assets 那条 CLI 路接上 C25 引擎…
   A	cmd/wisp/panel_assets_143_test.go
```

⇒ **归因对**：多出来那枚是 `cmd/wisp/panel_assets_143_test.go`，由票 143 那一程在 `9be3288` 提入；
它是一枚 `_test.go`，所以只进 ban #8 的分母（`cmd/=40`），**不进** bans #1-5 那枚 `cmd/=22`——
这正好与 §6.1 我用文件枚数复算出的口径咬合（22 那枚数在 `9be3288` 前后都是 22，因为我这一发的 `cmd/=22`
与它 §1.4 那发 `cmd/=22` 同值，而 `cmd/=39→40`）。**它这句"不是本程造的"成立。**

### §8.3 "分母随锚点变"这条它有没有写成可复算 —— **半格**

它 §8 给的是**两枚钉的 blob**（`9c06f3f6`／2080 行 ↔ `9fd6defb`／2081 行）加一句规矩
（"引本件任何一枚 d22scan 分母时连锚点一起引"）。问题是：**扫描分母取决于整棵树，不取决于那枚钉**——
把钉的 blob 当锚点引，下一程照样对不上（我今天这两发就是钉同值 `9fd6defb`、而 `cmd/` 与 §1.4 差一枚）。
⇒ 缺的是**那一发自己的 `git rev-parse HEAD`**。它 §1.4 那发写了 `锚点 cb60b94`（合格），
§8 这发只写了 blob（不合格）。补法是一行、append-only：在 §8 那段末尾追加"取于 HEAD `<sha>`（当时现量）＋
差值文件 `cmd/wisp/panel_assets_143_test.go`（`9be3288`）"。这一格记入账 **F-ACC-4**，
属读数/口径级，**不推翻它任何一条结论**（我今天两发独立跑到同值就是证据）。

---

## §9 总裁：**成立（附条件入账）**——条件是五笔账一起带，缺前两笔就变成下一次假绿

八格里没有一格推翻被验物：**F-R2-3 那三行确实从"借来的负控"变成了常驻承重测试**
（承重形最硬：摘掉它 ＋ M14/M16 ⇒ 我的干净副本**包内 0 枚红**，§4.3），
**F-R2-1 那条约束确实写进了文件头且措辞可证伪**（§5.2 我按它的句子把两种合法目的地都打了、四处红名逐枚对上），
三门＋名册＋契约轴＋断言尺全部自量通过，且它六句"要复算的话"我逐句复算到同值（§3、§6）。
判**不退回**的第二条理由：本程最强的那两发打穿（掏空第三枚清单、删一枚词）**都要编辑这把尺自己**，
而被验物从没声称它防得住人改它（`:91` 那句把这类改动的处置权明确推给了切片卡）。

入账五笔，前两笔是**这把尺的自空转面**，必须在它下一次被当作"边界干净"的证人之前关掉：

| 号 | 落在哪一枚文件与行为 | 最小闭合（到行为级；本程不替它改） | 我量到的现状 |
|---|---|---|---|
| **F-ACC-1** | `internal/panel/l2_grant_boundary_test.go` `:1257` 那枚反空转 `t.Fatalf` | 那枚 `if` 只点名 `grantRouteWords`／`grantRoutePrefixes` 两枚清单，**乘积的第三枚因子 `grantRouteSuffixes` 没人问**。闭合取一：把 `len(grantRouteSuffixes) == 0` 并进去；**或**在 `:1308` 之后钉一枚 `asked != len(prefixes)*len(words)*len(suffixes)*2 ⇒ t.Fatalf`（后者一并治掉 §1.1 的 V4 那形——砍掉复数那一圈 ⇒ `asked=264` 今天也照样绿） | `grantRouteSuffixes=[]` ⇒ 本测试 **PASS**、日志 `asked=0`、**全包 0 枚红**（`out/V1-suffix-empty.txt`、`out/FULL-V1-suffix-empty.txt`） |
| **F-ACC-2** | 同文件 `:1400` 那圈"路由证人" | 词表**内容**只有 4/11 枚有人质（`allow`/`approval`/`decide`/`grant`）。闭合：给 `approve`/`permit`/`ratify`/`authorize`/`authorised`/`verdict`/`decision` 各补至少一枚"必须仍读成批准门"的路由名（`ratify` 那枚现成就是 `panel.review.ratify`＝M16）。修完的判据：**从 `grantRouteWords` 删掉任意一枚，全包必须有红** | 删 `ratify` ⇒ 0 枚红；**再叠它自己那发 M16 ⇒ 53 枚顶层全绿、`asked=480`**（`out/RV-ratify-M16.txt`）。删 `allow` ⇒ 红（证人那圈有效，`out/RV-allow.txt`）；`permit`/`verdict`/`approve`/`authorised` 各 0 枚红（现量），`authorize`/`decision` **未逐枚打＝推定同形** |
| **F-ACC-3** | 同文件头段 `:112-116`（F-R2-1 那段） | 出路那句的 "or" 要写成**两步**：换成同包 struct **还不够**，那枚类型还得进 `inboundTypeRegistry`（`:1560` 那枚引信会替你喊）。只改措辞，不动判据 | 只换 struct ⇒ `FAIL=1`（`:1560`，`out/T-V3-MDEC-struct.txt`）；struct＋注册 ⇒ `FAIL=0` 且该类型真被两把尺对着判（`:1576 acProbe: 1 keys agreed`，`out/T-V4-MDEC-registered.txt`） |
| **F-ACC-4** | 被验件 `docs/evidence/s1/panel-l2-grant-nail-fix-r3.md` §8 | 那一发补它自己的 `git rev-parse HEAD`，并点名差值来源 `cmd/wisp/panel_assets_143_test.go`（`9be3288`）。append-only 追加，不改已提交行 | 我两发（`f50c037`／`0e17914`）跑到同值可证它是真读数，只是没带锚点 |
| **F-ACC-5** | 同件 §4 | `out/ 50 枚` 改（追加）成现量 55，并把"没测这把尺自己会不会少问"补进那一节（本程已在 §7.3 代补一句） | 按 mtime 现量：§7/§8 追加那批发出 5 枚 |

**两发反向判据**（升级或复核时用，别再靠话）：

1. **拆掉它必须红**：`python scripts/mut.py standing-out && python scripts/mut.py M14` → `go test ./internal/panel/ -count=1`。
   今天读数 **0 枚红**（`out/T-C2-out-M14.txt`）＝承重。任何"再补一枚常驻断言"的修法做完这一发**必须仍 0 枚红**；
   若做完变红，说明两枚判据在互相顶账，实现件 §2.4 那一支（"留下＝互相抵账"）就重演了。
2. **修完 F-ACC-1/2 之后必须红的两发**：`grantRouteSuffixes` 掏空 ⇒ 那枚测试必须 `t.Fatalf`；
   `grantRouteWords` 逐枚各删一次（11 发）⇒ **每一发**全包必须有红。删一枚而全包绿＝那一枚词今天仍然没人质。

**不并入退回的四条形状**（都记在这里免得下一程当新洞）：把新测试**改名**全仓无人响（§2.2，Go 通用形状，
本仓唯一 backstop 是人工名册差集）；改名后那 6 处文字引用悬空（§2.2 末，同 r2 记过的 `bridge.go:33` 同类）；
`asked=` 只是打印不是断言（§3 末，已由 F-ACC-1 的乘积方案一并覆盖）；V5/V8 那种 `if false` 短路任何 Go 测试都治不了
（记为极限，不记为洞）。

**给编排者的一句（我不勾票、不写台账）**：`Q-49` Go 侧那半的措辞在 r2 §7.5 给过，本批**没有把它推得更远**——
现在可以记的还是那句"限静态写下的形状＋网格内运行期拼名已被行为扫收住（528 问零误伤、双向红）"，
网格之外与词表之外的两形仍按 r2 §7.4 第 4 句（F-R2-4）开着。前两笔账（F-ACC-1/2）不改变这句的射程，
但它们决定这枚新尺**明年还作不作数**。

---

## §10 本程**没测**什么（不靠沉默读成通过）

- **没跑 `go test ./...`、没跑 CI、没跑前端四道门**（eslint/vitest/build/lint）。本件所有"全包／包内"字样一律只指
  **`internal/panel` 那一枚包**。
- **零枚计时判据**：这台机器上另有两程在跑（票 142 已交、票 143 在飞、`slo-full` 偶发），本程只取颜色／枚数／名册；
  所有命令的耗时我一秒都没记，正文里也没有任何"多少秒"的结论。
- **没在 Linux 容器里量任何东西**——本程复核过这个豁免站不站得住：`git grep -c 'go:build' HEAD -- internal/panel/`
  = **0 枚文件**，即本包没有"只能在 linux 量"的腿，所以这一格没测**不**留下平台盲区（只对本包而言）。
- **没测 F-R2-4 那两形**：词表外的结论键、以及名字从 config／网络来的路由。本程只复算了"被验件没声称闭上它们"。
- **没逐枚打完 11 枚词**：`authorize`／`decision` 两枚的删词后果是**推定**（同一圈代码），未实跑。
- **没测网格之外**的名字（`wisp.review.ok`、`composer.allow` 之类）会不会被包里别的测试兜住——没造那发。
- **没复算"今天什么都没接线"那枚自陈**（`ParseComposerRequest` 生产零调用者、`grep -c webview go.mod`=0）：
  它在文件头 `:65-73`，r2 §7.3 已独立验过一轮，本程**没有重做**，因此不替它担保。
- **没核**别的文件头是否还有 F-ACC-3 那种"两路其实是一路两步"的措辞（只查了 §5 那一处）。
- **没读** `design/**` 脏树内容作凭据（只用"它有 16 枚未提交删除"这一枚事实去解释 C21 的红/绿差）。
- **没写**台账、票面、`Q-49` 的勾；**未 push、未 `add -A`/`.`、未 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`、未 `rm`**。
- 临时件只建不删：`D:\tmp\panel-l2-nail-fix-r3-accept-dwl\` 下 `backup/` 3 枚、`out/` 50 枚读数、
  `scripts/` 4 枚（`mut.py`／`routelex.py`／`routevocab.py`／`worddrop.py`，前两枚是同族两版）、
  两枚副本目录（`headcopy`＝`a5e1c8c`、`headcopy-r3anchor`＝`121006d`）；**仓库树内未建 worktree、未 checkout**。
- 每发变异之后现跑的 `git status --porcelain -- internal/panel/` 全程为**空**（最后一次：`d59d99a` 之后复算＝空，§8 末）。

---

## §11 给 owner 的一段人话（不用术语）

我们让机器去问一个新问题："如果后台程序**自己拼**出一个'同意'的暗号，那扇门会不会开？"
它把能拼出来的五百多种拼法全部试了一遍，今天**一种都不开**——这活儿干得对，我另外找台机器重算也是同一个结果。
更关键的是一件事：**如果把这枚新加的问法整段删掉，那扇本来会被拼开的门就再也没人报了。**
也就是说这次补的不是装饰，是真在顶东西的。它前面那一版被退回，是因为那个"不许答应"的规矩**寄在别人的段落里**，
谁改那段话它就没了；这次它站成了独立的一条，这个洞补上了。

但我找到两处**它自己会偷偷不干活**的地方，都要改，都很小：
第一，这条问法由三样东西相乘凑出来（地点清单 × 动作清单 × 尾巴清单）。前两样清单**空了会报错**，
第三样**空了它什么都不问，还照样报"一切正常"**。第二，动作清单里有十一个词，
其中七个词今天**没有任何人守着**——只要有人把其中一个删掉（比如"批准"那个词），
那么"后台真答应了一次批准"这件事就会**从头到尾全绿**，而这正是我们花钱要它拦的那一幕。
两处的补法都是在文件里加一两行，不动产品代码。

还有三笔是文档账，不是产品账，一并列在上面：它给后来人指的那条"出路"其实要走两步，
它写成了"或"；它引用扫描门的那串数字时，漏了自己那一次是在哪个版本上量的；
它自报的临时文件数目在交件后又长了五枚。这些都不改变结论，但如果现在不记，
明年就没人知道那枚"全绿"值多少。

顺带替它说句公道话：这一程有两次**主动报了自己不利的读数**（跑错顺序把一发读数废掉、
"三枚"改成"四枚"），这两处我都复算过、都是真的，而且都往不利于自己的方向改。这类行为在这份工作上比结论值钱。
