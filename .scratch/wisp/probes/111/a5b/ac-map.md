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
| AC#1 | 票面 `:22-23` 逐字：`- [ ] **AC#1** 复算并出一张**全仓对账表**：每个包 ×（有无测试文件 / 在不在 CI 某一步 / 那一步真给过结论的 run id + step 号）。` ⏎ `      表格必须能自证完整（\`go list ./...\` 的 33 行都在，不许只列零覆盖那几个）。` | 表本体＝`.scratch/wisp/probes/111/r4/census-table.tsv`（36 行＝表头 1＋数据 35）＋`r4/evidence.md:31-99`（§1 表体、`:33` 自证尺、`:37-47` 四列口径、`:87-95` 调用点图例）；commit `11b8f288`（骨架＋12 枚台件）／`418d5b58`（表补满，2026-10-06 09:50:38）；第④列凭据＝`r4/evidence.md:182-207`（§5 方法＋读数表）；票面"33"的出处＝`:23` 与 `:16-18` 现场段 | r4 读数：表体 **35 行 ＝ `go list ./...` 现量 35**（本机 GOOS=windows，`r4:20`/`:103` 时刻 10-06 09:07:54），`comm -13`／`comm -23` **双向空**（`r4:35`）；第④列 **9 行写"待推送后回填"**＋census 整层待回填（`r4:287-288`）。★本腿自己的零 go 复算尺（10-06 10:1x，`78072249`）：`git ls-files '*.go'` 剥三枚独立模块（`tools/d22scan`·`tools/mockllm`·`scripts/spike`）与 `testdata/` 后，`cmd/`＋`internal/`＋`frontend/`＋`tools/signmodels` 含 `.go` 的目录数＝**35**，与 `r4/golist-win.txt`（35 行）`comm -13`＝**0 行**、`comm -23`＝**0 行** ⇒ 35 这一格不靠任何腿的记忆、也没跑 `go list` | ①第④列那 9 行＋census 那一层要**census 步（`ci.yml:596`，`6c0e3e31`）推到远程并跑出一枚 run 的步级结论**；②票面那句"33 行都在"要裁决者认"现量 35 ⇒ 判据按 35 行核"（本腿不改票面一字） | **推送**（第④列的 census 层＋`internal/session`/`internal/projctx` 两行认领读数，`1bb654e3` 晚于 tip `c6cf66e6`，本腿 `merge-base --is-ancestor` 判 NOT on origin/dev）＋**前提已翻**（"33 行"这半；复算尺见④栏末：35 目录 vs 35 行双向 `comm` 空） | **不翻**。三列里第③列（在不在 CI 某一步）与表的自证完整已交齐，但票面把"那一步真给过结论的 run id + step 号"写进同一格，那一半在推送前不可能有读数——拿本地表冒充整格闭合就是票面 `:57` 禁的那件事。⛔ 本腿不是裁决者 |
| AC#2 | 票面 `:24-26` 逐字：`- [ ] **AC#2** 逐包接入，**先易后难**，每包一次可核对的步级读数。` ⏎ `      ⚠ **加严可以直接做**；**不许**为了让某包变绿而放宽它的断言、调它的阈值、或给它加 \`//go:build\`/\`t.Skip\` 挡掉。` ⏎ `      接入第一天就红 ⇒ **那是发现**：红名逐条登记进本票面并**开票**，不许撤步骤。` | 名册真身（本腿 10-06 10:1x 现量，HEAD `78072249`／blob `2ff02dd7`）：`scripts/portable-tests.sh` core_pin `:164`（**27 行**）·win_pin `:193`（**10 行**）·cli_pin `:205`（1 行 `cmd/wisp` `:206`）·winsec_pin `:208`（1 行 `:209`）；票面点名的五枚 pin 行＝`cmd/llmrecord :165`/`:194` · `internal/ball :169`/`:195` · `internal/perm :182`/`:197` · `internal/plugin :183`/`:198` · `cmd/wisp :206`；scope 数组 `:242-243`（core）／`:251`（win）／`:260`（cli）；CI 调用点 `ci.yml:373` core／`:471` winsec-tests／`:507` wisp-cli／`:596` census／`:632` windows。腿证据件＝`r1b/evidence.md` §0/§2、`r2/evidence.md` §2 表（7 枚逐枚）、`r4/evidence.md` §1。产码凭据＝`8fe5c7ce`（2026-09-21 21:05:37，标题即含 AC#2）／`7699ec3f`（21:22:15）／`1bb654e3`（session+projctx 进 scope/pin，+106/−11） | ★本腿今天自己取的推送读数（run `37396530365`／head `c6cf66e6`／created 2026-10-06T00:55:36Z／run conclusion `failure`；取数时刻 10:34:29→10:42:24，台件 `/tmp/a5b-logs/core-37396530365.log` 750,298 字节＋`/tmp/a5b-logs/win-37396530365.log` 784,023 字节）：core `step 7 Portable package tests`＝**failure**，own-line 表 **25 行**（尺＝`grep -cE` 数 `portable-tests.sh:` 开头＋`ok/FAIL (own line)` 的行，结果＝25，10:40:19），四数逐字 `=== RUN=1518 --- PASS=1028 --- FAIL=4 --- SKIP=1`（日志 `:5396`），4 枚红逐名＝`TestApprovalCardViewJSONKeysMatchFrontendTypes`/`TestC21DesignTokensFourWayAgree`/`TestComposerContractTypesMatchFrontend`/`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`（全 `internal/panel`）；windows `step 8 Portable windows tests`＝**failure**，该步 own-line 表 **8 行**（按步分组数：日志 `:3840` 起属 step8 的行＝7 枚 `ok`（`cmd/llmrecord`·`internal/ball`·`internal/config`·`internal/perm`·`internal/plugin`·`internal/proc`·`internal/secret`）＋`FAIL github.com/CarlosShao/wisp/internal/risk`；全日志 own-line 合计 10 行＝step4 winsec 门 1 行＋step7 CLI 表 1 行 `FAIL cmd/wisp`＋step8 这 8 行），四数 `=== RUN=536 --- PASS=362 --- FAIL=12 --- SKIP=1`；cli `step 7`＝**failure**，四数 `=== RUN=323 --- PASS=220 --- FAIL=6 --- SKIP=1`。⚠ 已推送配置的 core 表是 **25 行**，本地 pin 是 **27 枚**——差的正是 `session`/`projctx`（`r4:203-204` 同判） | ①`internal/session`＋`internal/projctx` 各一枚**可核对的步级读数**（要 push）；②"红名逐条登记进本票面并开票"这一半：panel 那 4 枚、risk 那 12 枚、`cmd/wisp` 那 6 枚今天都在 CI 上红，票面要求登记＋开票，而票面正文本腿一字不许改 | **推送**（2 枚新包的读数）＋**别的票**（红名归口：panel＝票 92/115 地界 `r4:250-253`；risk windows 12 枚＋`TestSyncRedTeamRealOneDrive` 的 SKIP＝票 123 AC#5 `r4:254-259`；`cmd/wisp` POSIX 那族＝票 98） | **不翻**。"逐包接入"这一半在名册上已闭合（五枚早已在册＝`r1b` §0 推翻票面前提，本腿复认行号），但判据是"每包一次可核对的步级读数"——`session`/`projctx` 两枚今天仍零读数，且红名登记这一半没人做过。⛔ 本腿不是裁决者 |
| AC#3 | 未填 | 未填 | 未填 | 未填 | 未填 | 未填 |
| AC#4 | 未填 | 未填 | 未填 | 未填 | 未填 | 未填 |
| AC#5 | 未填 | 未填 | 未填 | 未填 | 未填 | 未填 |
| AC#6 | 未填 | 未填 | 未填 | 未填 | 未填 | 未填 |
| AC#7 | 未填 | 未填 | 未填 | 未填 | 未填 | 未填 |
| AC#8 | 未填 | 未填 | 未填 | 未填 | 未填 | 未填 |
| AC#9 | 未填 | 未填 | 未填 | 未填 | 未填 | 未填 |
| AC#10 | 未填 | 未填 | 未填 | 未填 | 未填 | 未填 |
