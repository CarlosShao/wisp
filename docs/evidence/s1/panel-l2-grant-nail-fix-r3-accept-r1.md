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
