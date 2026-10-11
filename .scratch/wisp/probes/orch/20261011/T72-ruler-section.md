
---

## [2026-10-11 11:2x] 编排者追加一节：放写腿进去**之前**先把 `AC#5` 那枚取证尺换形（⛔ 翻任何一枚框、⛔ 改上面任何一句原话、⛔ 删上面任何一段）

**这一节是什么**：本票 `AC#5` 已勾，它交出的凭据里有一把尺是**逐字节钉住一枚测试文件**。现量发现：甲族那 7 枚红里**有 4 枚就住在那枚文件里** ⇒ 任何去修期望侧的腿，第一笔 commit 就会把自己票上已勾的凭据打红。写腿⛔ 有权裁自己票上的格（`SPEC-12 §4.3` 裁决者≠实现者），所以这枚冲突**由我在落笔前解掉**，⛔ 留给腿去"顺手"处理。

### 1. 现量（三把尺，都我这发自己跑，⛔ 抄腿）

| 尺（逐字可重跑） | 读数 |
|---|---|
| `git diff --quiet a04d3e2 HEAD -- internal/risk/pathresolver_junction_windows_test.go` | **rc=0** ⇒ `AC#5` 那句"逐字节相同"今天**仍成立**（改动尚未开始，这一节⛔ 作废任何已有读数） |
| `git diff --quiet a04d3e2 HEAD -- internal/risk/pathresolver_anchor_spelling_windows_test.go` | **rc=1** ⇒ 另一枚文件⛔ 在这把尺的射程里（它当初就没被逐字节钉住） |
| `grep -n 'func Test' internal/risk/pathresolver_junction_windows_test.go` | **10** 枚；`anchor_spelling` 那枚 **6** 枚；两文件合计 **16** |

甲族 7 枚红的**红位**（腿件 `probes/72/a2/logs/aliased-7named.txt`，我 `grep -E '^\s*--- FAIL|\.go:[0-9]+:'` 自读过，逐枚按行号落进函数体对得上）：

- `anchor_spelling.go:57` → `TestClassifyAnchorSpellingIsNotVerdict`（`:29` 起）／`:140` → `TestCanonicalInputGainsNoSecondForm`（`:125`）／`:238` → `TestAListWinsWhereBothTablesHit`（`:204`）＝**3 枚住⛔ 被钉的那枚文件**
- `junction.go:135` → `TestPathResolverShortNameAListDenied`（`:123`）／`:164` → `TestPathResolverUNCAListDenied`（`:148`）／`:182` → `TestPathResolverExtendedLengthPrefixAListDenied`（`:174`）／`:292` → `TestBListDefaultDenyAndOverride`（`:272`）＝**4 枚住⛔ 被逐字节钉住的那枚文件** ⇒ 本节存在的因

### 2. ★裁语：从这一节起，`AC#5` 的取证方式换成"两台面 × 两发 ＋ 红位集合⛔ 变小"

**⛔ 再引上面那把 `git diff --quiet … -- pathresolver_junction_windows_test.go` 当"⛔ 放宽断言"的凭据。** 它防的那件事（有人为了买绿去动仪器）改由下面两发放牙：

- **台面两枚**＝①**别名台面** `TEMP/TMP/TMPDIR` 指到 `…\RISKAT~1\RUNNER~1`；②**本机台面** 不设这三枚 env。每枚台面各建一份**仓外导出树**（`git archive HEAD go.mod go.sum internal cmd | tar -x`），⛔ 拿母仓⇄clone 当那根轴（先例＝`A839` §1 那八发）。
- **每台面两发**＝ⓐ那 **7 枚**指名的名册；ⓑ两文件**全量 16 枚**。命令一律 `go test ./internal/risk/ -count=1 -v -run '<名>'`，四数（RUN/PASS/FAIL/SKIP）逐发落件。
- **判据甲（该绿的绿）**＝别名台面 7 枚 **FAIL→PASS**；**本机台面颜色⛔ 变**（同一份码只差那三枚 env ⇒ 若本机台面也变了色，就是拿环境换绿）。
- **判据乙（⛔ 放宽，替掉"逐字节"那半句）**＝**红位集合⛔ 变小**：改后任何一发里新出现的红必须逐枚指名；且改后那 16 枚的全量名册⛔ 允许出现"从没有→有"的新红。这条尺咬的是"把某枚断言删掉/降级/`Skip` 掉"那一族假绿。
- **`AC#5` 里仍然有效、⛔ 由本节改动的那半句**＝"`ErrReparseDenied` 与 `reparse_point_exceptions` 的『只按精确路径』语义零改动"＋"A/B 表内容与规则字符串逐字保留"。这两句⛔ 依赖那枚文件逐字节不变，改由**逐笔 hunk 尺**取证：`git show --numstat --format=` 每一笔写腿 commit，产码侧（非 `*_test.go`）命中必须 **0 行**，且测试侧那几笔的 hunk 逐枚读回、⛔ 出现表字符串或豁免规则的改动。

### 3. 四条禁线（下一程写腿的派单里逐字重复一遍；这里只落一处，⛔ 另开仪器票）

1. ⛔ 动 `overrideApplies` 让它遍历全部拼写形。票面 `:95` 逐字写过「插单前我在 AC#5 写的那段"豁免查表改成全部拼写形"**作废**，那正是被指出来的放宽」，正钉＝`TestOverrideOnlyAcceptsTheResolvedForm`（`anchor_spelling.go:157`，今天**绿**，⛔ 在 7 枚名册里）。⚠ 而且它今天**连着生产分档**（`internal/tools/mode.go:91` `overrides := b.confirmations()` → `:93` `Canonicalize` → `:98` `risk.Gate(c, overrides)`），动它⛔ 是"只动测试"；票面 `:245-247` 那句"今天还没有生产 caller 传 override 表"**已过期**，我在台账 `A840` §1 落这句，⛔ 改票面原文。
2. ⛔ 把 `t.TempDir()` 换成手写常量＝票面 `:201` 逐字禁止的"在测试里复刻 bug"。期望侧要按**四形各写一次**（长形／8.3 短形／`\\?\`／UNC），⛔ 钉死其中一形。
3. ⛔ 把 #1 那枚 `t.Fatalf` 降级（`:57` 那句 `test premise broken (candidate must be the real path)` 搬走了 `AC#3①` 的变异靶）。
4. ⛔ 把 #2 的 `want 1` 改成 `want<=2`、或删掉那枚尚不存在的叶子（那会悄悄缩小 `AC#5` 的爆炸半径承诺）。

### 4. 归口（一句，⛔ 在本节里做）

票 307 收档时追加给我的**负向钉**欠账＝`TestSyncFallbackNotDisarmableByWeakRoot`（`syncdirs_test.go:209/:212/:216`）与 `TestSyncUnverifiedRootKeepsFallback`（`:229/:232`）体内只有正向断言（恒真形下它们绿）。本票写腿**同批**带上（同一族"期望侧⛔ 负向钉"，且甲 7 的期望侧已归本票）＝**裁语在台账 `A839`/`A840`，撤销口令「307 负向钉另立一票」**。

**撤销本节**：口令**「72 取证尺换回逐字节」**——回到 `AC#5` 原文那把 `git diff --quiet … -- pathresolver_junction_windows_test.go` 尺，代价＝甲族 4 枚红⛔ 有合法的修路（任何修它们的 commit 都会把已勾的 `AC#5` 打红）。
