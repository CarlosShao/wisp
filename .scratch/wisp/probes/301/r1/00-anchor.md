# 301-r1 — 起手锚（⛔ 任何产码改动之前跑完并落本笔）

钟点＝`date` stdout 插值：2026-10-10 14:26:28 +0800

## 0. 本文件里每一个数的尺口径（三句逐字）

1. **HEAD 尺**＝`git rev-parse --short HEAD`，射程＝整仓，读的是**提交对象**不是工作树。
2. **porcelain 尺**＝`git status --porcelain | wc -l`（整仓，含未跟踪；`-uall` 未用），
   另有**写面专用尺**＝`git status --porcelain -- scripts internal .github | wc -l`（同一把 `git status`，射程收窄到三枚目录）。
3. **名册尺**＝对 **HEAD blob** 逐枚 `git show HEAD:<path>`（⛔ 对工作树量），
   tagged 判据＝`head -3 | grep -c 'go:build windows'`（逐字模式串 `go:build windows`），
   用例计数＝`grep -c '^func Test'`（逐字模式串 `^func Test`，**顶层 only／⛔ 含子测试／⛔ 含注释行**）。

## 1. 起手两枚（双锚规矩的第 1 锚）

- `git rev-parse --short HEAD` ⇒ **`64ceaacd`**（＝派单给的期望起手号，对上）
- `git branch --show-current` ⇒ `dev`
- `git status --porcelain | wc -l` ⇒ **810 行**
  - 按码分类：` M` 16 枚／` D` 16 枚／`??` 778 枚
  - 16 枚 ` M` 全在 `.gitignore`·`.scratch/wisp/probes/**`·`design/doubao/**`
  - 16 枚 ` D` 全在 `design/**`
  - 778 枚 `??` 里 724 枚在 `.scratch/**`，其余＝`design/**` 11 枚＋仓根散件（`part*-*.txt`、`logs`、`agent`、`4`、`397`、`2379`、`2026-09-24`、`-`、数枚中文名垃圾件）
- ★**写面专用尺**＝`git status --porcelain -- scripts internal .github` ⇒ **0 行**（输出为空，已现量）。
  ⇒ 派单开头那 810 行**与本票写面零交集**：`scripts/`／`internal/`／`.github/` 三枚目录在工作树上**一枚字节都不脏**。
  本票不需要那 810 行里的任何东西，也⛔ 会碰它们。

## 2. ⚠ 派单与盘上事实的三处具名差别（开工前报，⛔ 自行假设）

1. **起手 porcelain 枚数**：派单未给预期数；票面 `301-a2` 那节记的是「`86d89478`／porcelain 2 行」＋「交回 `41329475`／porcelain 0 行」。
   现量＝**810 行**。差别**不来自本票**（写面专用尺 0 行坐实），全部来自 `.scratch/**` 与 `design/**` 的在飞/遗留件。
   ⇒ 本票按「写面专用尺＝0」开工，⛔ 按整仓 810 行判定任何东西，也⛔ 去清理别人的件。
2. **探针落点**：派单逐字写 `probes/301/r1/logs/…`，票面 `AC#4` 逐字写名册「只含 `scripts/portable-tests.sh` ＋ `probes/301/**`」。
   盘上先例（`git log --name-only` 逐笔读到）＝同族腿全部落在 `.scratch/wisp/probes/<票号>/<腿>/`：
   `300-r1`＝`.scratch/wisp/probes/300/r1/logs/rosters.txt`（commit `41329475`）、`301-a1`＝`.scratch/wisp/probes/301/a1/**`、
   `301-a2`＝`.scratch/wisp/probes/301/a2/**`。仓根 `probes/` 确实存在但**只含未跟踪的 `280`，`git ls-files probes/`＝0 枚**（`git check-ignore` 退码 1＝⛔ 被忽略，纯粹没人跟踪）。
   ⇒ 本腿按盘上先例落 **`.scratch/wisp/probes/301/r1/**`**（`AC#4` 那句 `probes/301/**` 作后缀匹配成立），
   ⛔ 在仓根新建第二座探针树。这一处如果编排者要的是字面仓根，口令下达后本腿整批平移。
3. **`AC#1` 那句「⛔ 拆两笔＝第二笔会让 GUARD C 在中间态红」**：票面 `AC#1` 逐字仍这么写，而派单与票面 14:0x 裁定节都已更正
   ＝那两句「Pull the package / into a named scope and update that tier's pin in the SAME commit, or」长在 **GUARD D 的退码文本**里。
   本腿现量坐实（见 §4），两枚判红者都在：**GUARD D 咬 census 的「无档认领」，GUARD C 咬档清单↔pin 分家**。
   ⇒ 本腿照旧**同一笔**落两处（这条硬约束的成立理由由 C 提供，⛔ 由 D 的引文提供），⛔ 改票面一字。

## 3. `internal/audio` 名册（尺＝§0 第 3 句，射程＝HEAD blob `64ceaacd`）

`git ls-tree --name-only HEAD internal/audio/` ⇒ 20 个文件。
其中 `_test.go` 逐枚问 `head -3 | grep -c 'go:build windows'`：

| 文件 | tagged | `grep -c '^func Test'`（顶层） |
|---|---|---|
| `internal/audio/capturelevel_windows_test.go` | 1 | 1 |
| `internal/audio/hotplug_test.go` | 1 | 8 |
| `internal/audio/parse_wave_format_300_windows_test.go` | 1 | 1 |

⇒ **3 枚 tagged 文件／10 枚顶层用例**。分母按 **10**（票面 14:0x 裁定节末条已裁「`300-r1` 那枚算第 10 枚」，`A809`）。
★`internal/audio/parse_wave_format_300_windows_test.go` **此刻在 HEAD blob 上**（票 300 落地腿已交回，`64ceaacd` 那笔带它），
⇒ 派单「⛔ 动 `internal/audio/**`」这条在本票写面上天然满足：本腿对 `internal/audio` 零字节。

## 4. 三句逐字禁区的盘上现场（本腿⛔ 动，逐枚读到那一行）

- **`:593` ledger 行**逐字起头 `"TestLiveWasapiSmoke|./internal/audio/|windows|fixture|live WASAPI capture needs WISP_LIVE_MIC=1 …`，
  理由串里带行号锚 `hotplug_test.go:527`。⇒ **⛔ 动**（编排者 13:3x 裁定：那行第二条腿活着，删它会摘掉一枚管「改名/删用例」的钉，
  且删后第 10 枚会真跑出 `--- SKIP` 把 `test-windows` 判红，`tools/d22scan/runtests.sh:98`→`:102`）。
- **GUARD D 块＝`:407`–`:423`**（`:409` 起头逐字 `echo "portable-tests.sh: GUARD D - $guardd package(s) compile a test file for GOOS=$goos and"`）。
  那两句在 **`:417` 尾部 `… zero coverage for four days after their tests landed. Pull the package`** ／
  **`:418` 开头 `echo "portable-tests.sh: into a named scope and update that tier's pin in the SAME commit, or"`**，
  两行各以 `echo "portable-tests.sh: ` 起头 ⇒ 票面 `sed -n '417,418p'` 那把尺对上，引用时⛔ 写成一行「逐字」。
- **GUARD C 块从 `:542` 起**（本腿 `sed -n '553,560p'` 现量），其同义句在 **`:556`–`:558`**：
  `GUARD C - scope mode=$mode resolved to a DIFFERENT package` / `set than the one pinned next to it. Pinned: …, resolved: …` /
  `a deleted or newly-uncovered package must fail the step, not` + `shorten it. Either restore the scope entry or update the pin` + `in the SAME commit and say why in the CI log.`
  ⇒ ⛔ 在 GUARD C 里补任何一句（重复立法）。

## 5. 改动前的两处靶（＝本票要修的那件事，现量）

- `windows)` 档清单（`:248-253`）**10 枚路径**，逐字＝
  `./internal/proc/ ./internal/secret/ ./internal/config/ ./internal/risk/ ./internal/ball/ ./internal/perm/ ./internal/plugin/ ./cmd/llmrecord/ ./internal/session/ ./internal/projctx/`
  ⇒ `grep -c audio` 于这一段＝**0**（audio ⛔ 在）。`pinned=$win_pin` 在 `:254`。
- `win_pin` 正文（`:193` 起）**10 枚导入路径**＝`cmd/llmrecord` `internal/ball` `internal/config` `internal/perm` `internal/plugin`
  `internal/proc` `internal/projctx` `internal/risk` `internal/secret` `internal/session` ⇒ ⛔ `internal/audio`。
- `core` 档侧（尺＝`grep -n 'internal/audio' scripts/portable-tests.sh`，现量**三处命中**）＝
  `:168` `github.com/CarlosShao/wisp/internal/audio`（**core_pin 正文**）／
  `:241` `./internal/buildinfo/... ./internal/audio/... ./internal/proc/...`（**core 档清单**）／
  `:593` ledger 那行。⇒ audio 已被 `core` 认领，而 `core` 跑在 `ubuntu-latest`
  （尺＝`sed -n '402,403p' .github/workflows/ci.yml` ⇒ `  test-core:` ＋ `    runs-on: ubuntu-latest`）
  ⇒ 本票要修的形状成立。`:542` 现量逐字起头 `# GUARD C: the named scopes pin their own resolved set. Deleting an entry from the`
  （＝派单说的「C 从 `:542` 起」对上）。

## 6. 进程前置尺（跑长跑件之前现量）

`tasklist //FI "IMAGENAME eq wisp.exe"` ⇒ 匹配计数 **0**；`tasklist //FI "IMAGENAME eq balldebug.exe"` ⇒ 匹配计数 **0**
（后者 stdout 逐字 `No matches found` 且 `find` 退码非 0）。⇒ ⛔ 真窗在跑，长跑件放行。

## 7. 格式卫生

`bash -n scripts/portable-tests.sh` ⇒ **rc=0**（改动前基线）。
