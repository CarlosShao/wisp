# 票 111 — 111-a5b 交件：AC#1…AC#10 十格映射表（唯一交付）

代号 `111-a5b`（只读普查腿，接零足迹死掉的前腿 `111-a5-acmap`——它的 `.scratch/wisp/probes/111/a5/` 从未建立，本腿不找它的半成品）。
写面＝`.scratch/wisp/probes/111/a5b/**` ＋票 111 Progress log 行。票面正文一字未改、十枚 `- [ ]` 框一枚未翻。

## §0 起手锚（本腿自己现读，时刻 `2026-10-06 10:10:37 → 10:2x +0800`）

| 项 | 读数 |
|---|---|
| `date` | 2026-10-06 10:10:37 +0800（复量 10:21:17） |
| `git rev-parse HEAD` | `78072249780203d57ff641f58f5796d2946c4a8d`（短 `78072249`，分支 `dev`） |
| `git rev-parse --short origin/dev` | **`c6cf66e6`**（`git log -1 --date=iso`＝2026-10-05 10:58:27 +0800） |
| `git rev-list --count origin/dev..HEAD` | **167**（本票四格 `1bb654e3`／`6c0e3e31`／`98f62fde` 逐枚 `merge-base --is-ancestor` 判 **NOT on origin/dev**） |
| ★推送那道硬闸（本腿复量） | `git show c6cf66e6:.github/workflows/ci.yml \| grep -c 'scope=census'` = **0**；同文件 `grep -c census` = **0**；`git show c6cf66e6:scripts/portable-tests.sh \| grep -c 'GUARD D'` = **0** |
| 同一闸的边界（别读成"推送前什么都还没接"） | `c6cf66e6` 的 `portable-tests.sh`：GUARD A=6 / B=6 / C=11 / `escape_re`=2 / `internal/winsec`=14 / `unaccounted`=3 / `TestSyncRegistryProbeLive`=2 枚命中，`TestWorkspaceSwitchRefusesAJunctionToOutside`（AC#10 那枚 skip）**已在册 1 次**；ci.yml 调用点行号 `:373`core / `:471`winsec-tests / `:507`wisp-cli / `:543`windows（`scope=windows`） |
| HEAD 上的四枚 blob（blob 同一性＝前腿行号可复认） | `ci.yml`＝`014861149bde`（＝r3/r4 记的那枚，未漂）·`portable-tests.sh`＝`2ff02dd7`（＝r2/r3/r4 记的那枚，未漂）·`portable-tests-selftest.sh`＝**`b8222356`**（≠ r3 §0 的 `6c5a5a36`，`98f62fde` 之后漂过）·`wisp-cli-tests.sh`＝`5fd918ce`（与票面换源复核同值） |
| 本腿零 go 命令的替代尺（枚数一律 `git ls-files`／`git ls-tree` 族） | `git ls-files '*.go'` 剥 `tools/d22scan`·`tools/mockllm`·`scripts/spike` 三枚独立模块与 `testdata/` 后，`cmd/`+`internal/`+`frontend/`+`tools/signmodels` 含 `.go` 的目录数＝**35**，与 r4 留盘 `golist-win.txt`（35 行）**双向 `comm` 空**（`comm -13`／`comm -23` 各 0 行）⇒ AC#1 的 35 这一格不靠任何腿的记忆，也不靠本腿跑 `go list` |
| 票面格数尺（我自己现量） | `grep -c '^- \[ \] \*\*AC#'`＝**10**（逐条 AC#1…AC#10；★不是四枚腿一直在说的"五格"） |

## §1 十行映射表

★AC#6 的生死判据（票面 `:291` 与 `:255` 两处逐字，本腿原样抄，⛔ 不替谁裁）：
> `是本票唯一的生死判据：**只要还有一步是 skipped，就判 AC#6 FAIL**，不接受"通过附条件"。`

| ① AC | ② 票面那一格原文逐字 | ③ 现有凭据（file:line ＋ commit ＋ 腿证据件节号） | ④ 读数（rc／枚数／红句，逐字抄台件＋取数时刻） | ⑤ 还缺什么才能翻勾 | ⑥ 缺的那半被什么挡 | ⑦ 按今天的盘该翻还是不该翻（建议位） |
|---|---|---|---|---|---|---|
| AC#1 | 未填 | 未填 | 未填 | 未填 | 未填 | 未填 |
| AC#2 | 未填 | 未填 | 未填 | 未填 | 未填 | 未填 |
| AC#3 | 未填 | 未填 | 未填 | 未填 | 未填 | 未填 |
| AC#4 | 未填 | 未填 | 未填 | 未填 | 未填 | 未填 |
| AC#5 | 未填 | 未填 | 未填 | 未填 | 未填 | 未填 |
| AC#6 | 未填 | 未填 | 未填 | 未填 | 未填 | 未填 |
| AC#7 | 未填 | 未填 | 未填 | 未填 | 未填 | 未填 |
| AC#8 | 未填 | 未填 | 未填 | 未填 | 未填 | 未填 |
| AC#9 | 未填 | 未填 | 未填 | 未填 | 未填 | 未填 |
| AC#10 | 未填 | 未填 | 未填 | 未填 | 未填 | 未填 |
