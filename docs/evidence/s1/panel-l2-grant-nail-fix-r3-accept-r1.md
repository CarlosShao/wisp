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
