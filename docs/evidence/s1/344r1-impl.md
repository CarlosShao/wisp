# 票 344 落地腿 344-r1 —— C21 真相源改指第四代 ＋ 三枚死玻璃色码删除（甲形）

`s1/` 证据 · 写腿 `344r1` · 2026-11-02 11:2x +08 · 锚点 HEAD `c9ef940c8`
机主明示＝甲（逐字「都按推荐」，台账 `A1395`：`344＝甲（C21真相源改指v3＋删三枚死玻璃色死码）`）。
派单存档＝`.scratch/wisp/probes/orch/20261102/txt/344r1-dispatch.txt`。
写面＝`internal/ball/**` ＋ `docs/evidence/` ＋必要 `tokens_table` 测试；⛔ `frontend/**`／`design/**`／`cmd/wisp/**`／thr/golden／票面／台账／push（全部未碰）。

## 一、AC#0 起手现跑（八行逐枚重跑，2026-11-02 10:5x–11:2x）

| 票面读数 | 本腿现跑 | 一致? |
|---|---|---|
| 对照表 20 行引 `tokens.css` | `Select-String` 现数 21 行（表体 20 ＋ 本文件生成注记 1）；`docs/evidence/s1/c21-native-tokens.md:5` 逐字 `权威真相源：design/assets/tokens.css（C21 契约参考实现）` | ✅ |
| v3 无玻璃族 | `theme.css` 现数 glass＝**0**；`tokens.generated.css`＝**0** | ✅ |
| 三枚死码位置 | `tokens.go:77/78/79`（声明）＋ `:132-134`（暗赋）＋ `:183-185`（亮赋） | ✅ |
| `tokens.go` 外零读者 | `git grep -l -e ".GlassRing" -e ".GlassHi" -e ".GlassHiSoft" -- "*.go" ":!*_test.go" ":!internal/ball/tokens.go"` ＝ **0 文件**；含测试仅 `tokens_test.go:45` 一行 golden | ✅ |
| ⚠ 活玻璃族 | `LiquidLook.Glass`（bool）与 `GlassCaustic`/`GlassRimPx`/`GlassLipPx` 仍被 `renderer_windows.go`/`liquid.go` 读 ⇒ **未碰** | ✅ |
| 检查 5 硬要 CSS 在场 | `const c21TokensCSSPath` 已＝c21-r1 提交入库的第一代**留档副本**（14,392 B，逐字等于 `72b65c614^` 原件），`go test` 现跑 PASS ⇒ 票面那枚在册红**早已被改指副本消解**，本腿在此副本上继续钉 | ✅（比票面新） |

**★票面过期点（具名更正）**：票 344（2026-10-18 立）记 `TestC21TableColourRowsMatchTokensCSS` 为在册红；
现 HEAD `c9ef940c8` 上该检查读的是留档副本、**跑绿**（起手基线：`internal/ball` 整包 `ok`，FAIL=0，红名册为空）。
故本腿的「改钉」＝在既有副本钉之上补真相源文字与缺席反向钉，不是修红。

**AC#0 加扫（⛔ 只交三枚）**：`Palette` 全部 40 字段逐枚问「本包外有无读者」——除这三枚外，其余零读者字段
（`FgPrimary`/`FgTertiary`/`InfoLine`/`OrbShadow`/`Tint*` 等）均已在对照表「未接线豁免」册上有名有因，
由 `TestC21TokenConsumerReport` 双向钉着；三枚玻璃正是本票授权删除的增量，名册其余不动。

## 二、删前读数（三枚，逐值）

| 字段 | 暗（`:root`） | 亮（`[data-theme=light]`） | CSS 变量 |
|---|---|---|---|
| `GlassRing` | `rgba(255, 255, 255, 0.14)` | `rgba(255, 255, 255, 0.66)` | `--glass-ring` |
| `GlassHi` | `rgba(255, 255, 255, 0.18)` | `rgba(255, 255, 255, 0.90)` | `--glass-hi` |
| `GlassHiSoft` | `rgba(255, 255, 255, 0.08)` | `rgba(255, 255, 255, 0.60)` | `--glass-hi-soft` |

删除面：`tokens.go` 结构体声明 ×3 ＋ `DarkPalette()` 赋值 ×3 ＋ `LightPalette()` 覆写 ×3；
`tokens_test.go` golden 表随删 `light glass-ring` 一行（被删字段的值钉无处可钉，其职责移交缺席反向钉，注释写明去处）。
对照表 78 行（明 40 ＋ 亮 38）**一行未删**（票 344 AC#2：删行＝把债洗成不存在）。

## 三、改钉形状（⛔ 删断言换绿 —— 断言只增不减、只反向不清零）

1. **真相源**（AC#1 甲，逐字对照 `PLAN.md:3316`/`PLAN.md:3612` 契约文字）：对照表头「权威真相源」改指
   `frontend/src/styles/tokens.generated.css`（第四代生成件，`:root`＝light、`[data-theme="dark"]`＝dark）；
   并写明**逐值相等只对同代同文件断言**——第四代腿由「生成器表↔生成件」守（`internal/panel`），
   原生世代腿由本表 ＋ `TestC21TableColourRowsMatch*` 守（读第一代留档副本）。`tokens.go` 值一字未动（⛔ 改值）。
2. **反向钉**（`c21V2RetiredGlass`，`internal/ball/tokens_table_test.go`）：玻璃 6 行仍走 leg1（行＝留档副本逐字一致），
   leg2 从「三方逐值」改为「**字段必须缺席**」——`Palette` 若被偷偷加回字段，check1 与 check5 当场红。
3. **豁免册随删**：「未接线豁免」表里玻璃三行删除——登记对象（字段）已不存在，双向机器检查（多行⇒红）本就要求删；
   这是债的**清偿**非抹除，行/读数/裁定全部留在本表、表注记与本件。

## 四、两平台成对 ＋ 零回归（红名册逐名作差）

- `go test ./internal/ball/ -count=1`（windows/amd64，go1.27.1，跑前 `Remove-Item env:GOTMPDIR/env:TMPDIR`）＝ **ok**，FAIL=0。
- `-run TestC21|TestTokenGoldenValues -v` 关键行：
  - `verified 108/114 colour rows ... plus 6 retired-glass rows pinned ABSENT (ticket 344 shape A) (37 palette fields ...)`
  - `three-way colour check: 78/78 cited rows agree ... of which 6 retired-glass rows verified as row-cites-reference AND field-absent (:root=129 decls, light=65)`
  - `未接线豁免: 21 row(s) attribute all 24 zero-consumer token(s); 89/113 checked tokens are live`
- **改前基线红名册＝空（整包 ok）；改后红名册＝空 ⇒ 作差为空**；票面在册红 `TestC21TableColourRowsMatchTokensCSS`
  改前已因留档副本改指而绿（见上「过期点」），改后**仍绿且钉得更狠**——逐名点名，⛔ 拿别的红抵数。
- `GOOS=linux GOARCH=amd64 go vet ./internal/ball/` rc=0；`go build ./internal/panel/` rc=0（ball 的下游编译面）。
- **反证一发（mutation probe）**：临时把 `GlassRing` 加回 `Palette` ⇒ check1＋check5 双红（暗行 `:58`、亮行 `:103` 各点名，
  句逐字 `re-adding a retired colour is the owner's call, not a test edit`）；还原后复绿。证明缺席钉不是空转。
- `gofumpt -l internal/ball/`（`D:\work\base\gopath\bin\gofumpt.exe`，绝对路径）改后 ＝ 空；`gofmt -l` ＝ 空。

## 五、纪律格（AC#3/AC#6）

- `git diff --stat` 与题面写面一致：仅 `internal/ball/tokens.go`、`internal/ball/tokens_test.go`、
  `internal/ball/tokens_table_test.go`、`docs/evidence/s1/c21-native-tokens.md` ＋本新文件；票名/`Status:`/框一律未动；台账零字节；只 commit 不 push。
- AC#5 落档形状：`A##` 由编排者落（本腿⛔ 自行落档、⛔ 先动 `PLAN.md`——`PLAN.md` 的契约文字早在票 380 已落地，本腿零改动）。
- 本件与表注记互为指涉：表「亮色系」下注记 ＋「审计」补录条目 → 本件；本件 → 表与测试的具体行。
