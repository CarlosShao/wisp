# 72-r1 起手锚（写腿 72-r1，甲族 7 枚期望侧 ＋ 307 负向钉 2 枚）

时刻 `2026-10-11 11:30 +0800`。锚＝`9c83598d`，分支 `dev`。原始读数＝本件同目录 `logs/anchor-readings.txt`（一把尺：`date`／`git rev-parse --short HEAD`／`git status --porcelain -- <三枚授权测试文件>`／`ls -d` 两枚别名目录／`env | grep -iE '^(TEMP|TMP|TMPDIR)='`／`sed -n` 逐枚红位原文）。

## 0. 起手读数（全部本腿现跑，⛔ 抄腿／⛔ 抄派单）

- `git rev-parse --short HEAD` ＝ **`9c83598d`**；`git rev-parse --abbrev-ref HEAD` ＝ `dev`。
- 三枚授权测试文件在工作树**干净**（`git status --porcelain -- internal/risk/pathresolver_anchor_spelling_windows_test.go internal/risk/pathresolver_junction_windows_test.go internal/risk/syncdirs_test.go` ＝ **0 行**）⇒ 本腿的写面起点＝HEAD blob。
- 别名台面**盘上已在**：`ls -d C:/Users/swq/tmp/RISKAT~1/RUNNER~1` 与 `C:/Users/swq/tmp/riskattrib1/RunnerTempAreaForAttribution1` 两枚**都命中**（短形目录＋长形本体）。
- 本 shell 起手 env＝`TEMP=/tmp`、`TMP=/tmp`、无 `TMPDIR`（Git Bash 内部值，⛔ 是别名台面）⇒ 每发台面由该发命令行**显式设定／显式清空**，⛔ 依赖继承。
- 七枚红位逐枚 `sed -n` 读回＝派单快照**逐字对得上**（行号当 `9c83598d` 时刻快照，认锚用下面那七句内容锚）。

## 1. 名册尺（现跑，`grep -n 'func Test'`，射程＝两枚文件）

- `internal/risk/pathresolver_anchor_spelling_windows_test.go`：**6 枚** `func Test`（`:29` `TestClassifyAnchorSpellingIsNotVerdict`／`:74` `TestClassifySpellingInvariance`／`:125` `TestCanonicalInputGainsNoSecondForm`／`:157` `TestOverrideOnlyAcceptsTheResolvedForm`／`:204` `TestAListWinsWhereBothTablesHit`／`:250` `TestClassifyFailClosedWhenSpellingUnprovable`）
- `internal/risk/pathresolver_junction_windows_test.go`：**10 枚**（`:88` `TestPathResolverJunctionWindows`／`:123` ShortName／`:148` UNC／`:174` ExtendedLengthPrefix／`:191` CaseAndDotMix／`:214` `TestReparseExceptionExplicitPathOnly`／`:238` `TestAListDenyAndUnoverridable`／`:272` `TestBListDefaultDenyAndOverride`／`:303` `TestResolveIdempotent`／`:323` `TestResolveExpansion`）
- 合计 **16 枚**＝票 72 追加节 §1 那把 `grep -n 'func Test'` 尺的读数（6＋10＝16），**未变**。
- `internal/risk/syncdirs_test.go`：`TestSyncFallbackNotDisarmableByWeakRoot` 起 `:198`，`TestSyncUnverifiedRootKeepsFallback` 起 `:224`（与派单一致）。

## 2. 红位内容锚（七枚，逐字 `sed` 读回）

| # | 文件:行（快照） | 内容锚（逐字） | 期望侧坏形状（本腿判） |
|---|---|---|---|
| 1 | `anchor_spelling:57` | `t.Fatalf("Resolve(%q) canonical %q: test premise broken (candidate must be the real path)", p, res.Canonical)` | 前提句拿"测试自己 Join 出来的拼写"当逐字符串靶 |
| 2 | `anchor_spelling:140` | `... produced %d comparison forms %q, want 1 (verdict must not move for a path with one spelling)` | 计数硬编码（`want 1`） |
| 3 | `anchor_spelling:238` | `t.Fatalf("B-only control must stay overridable by design: %+v", d)` | override 键取短拼（＝把短拼当 OS 担保形） |
| 4 | `junction:135` | `t.Fatalf("short name did not expand to long path: got %q want %q", res.Canonical, long)` | 拼写字面量（`long` 来自 `tmpHome(t)`） |
| 5 | `junction:164` | `t.Fatalf("UNC spelling %q normalized to %q, want %q", spelling, res.Canonical, long)` | 同上 |
| 6 | `junction:182` | `t.Fatalf("\\\\?\\ spelling normalized to %q, want %q", res.Canonical, long)` | 同上 |
| 7 | `junction:292` | `t.Fatalf("B-list single-file override must allow: %+v", d)` | override 键取短拼（同 #3，形 IV） |

## 3. 本腿读过的原文（⛔ 只按派单转述干活）

- `.scratch/wisp/issues/72-atble-classification-runner.md` 全文（294 行）：`AC#5`（`:79-97`）、三种禁止解法（`:199-203`）、不变式（`:215-222`）、编排者插单 10:08（`:224-247` 区）、**11:2x 追加节**（`:256-294`：§1 现量、§2 换形后的取证方式（两台面×两发＋红位集合⛔ 变小）、§3 四条禁线、§4 负向钉归口）。
- `.scratch/wisp/issues/307-...-done.md:218`（负向钉欠账**归票 72 下一程**，撤销口令「307 负向钉另立一票」）与 `:236`（`AC#5` 勾＋甲 7→票 72／乙 5→票 307 的分家裁语）。⚠ 台账行号快照里写 `syncdirs_test.go:209/:212/:216`、`:229/:232`，本腿 `sed` 现量红句实际在 `:213`/`:217` 与 `:233` ⇒ **按内容锚认**（`"under-profile fallback must stay armed"`、`"suspect fallback must stay armed"`、`"writes under the injected root must be sync"`），行号差异⛔ 影响归口。
- `.scratch/wisp/probes/72/a2/*.md`（只读普查腿结论）：`10-seven-cases.md` 逐枚七问、`30-answer-B-line.md` 形 IV（override 族）＋"两枚正控⛔ 能删也⛔ 能反转"那条恒真风险、`40-nails-that-would-break.md`（AC#3① 变异名册点名 `:57` 那枚前提）。`logs/*.txt` 原始件⛔ 整枚 cat（按派单：只 `grep -E '--- FAIL|\.go:[0-9]+:'` 取行）。
- `.scratch/wisp/issues/README.md:206-218`「Hard global constraints」＋规矩 1/2/8（commit 显式 pathspec／⛔ push／临时件只建不删）。

## 4. 冲突／交不了（开头先记，收尾再补）

- **无推翻性冲突**。一处**行号级**差异已在 §3 具名（票 307 台账 `:209/:212/:216` vs 盘上 `:213/:217`）＝按内容锚认，⛔ 改票面原文。
- 一处**派单转述 vs 票面原文**的口径差：派单 §4 说本机台面"⛔ 设这三枚 env"，票面追加节 `:283` 同义（"不设这三枚 env"）＝一致，⛔ 冲突。
- 本 shell 的 `TEMP=/tmp` 是 Git Bash 注入的**非 Windows 台面值** ⇒ 本机台面那一发必须**显式**把 `TEMP`/`TMP` 指回 `C:\Users\swq\AppData\Local\Temp`（⛔ 留 `/tmp`，否则那枚台面既⛔ 是别名也⛔ 是本机，读数无效）。这一点具名报回，属"派单转述没覆盖到的台面细节"。

## 5. 计划（本格逐枚，⛔ 动产码）

甲组 7 枚期望侧（`#1/#2/#3`＝anchor 文件、`#4/#5/#6/#7`＝junction 文件）＋乙组 2 枚负向钉（`syncdirs_test.go`）；`#2` 在 Ⓐ／Ⓑ 里选一支并写明代价（本腿**倾向 Ⓑ**：期望值由 OS 独立测量导出，`want 1` 那句"恰一形"的爆炸半径承诺⛔ 缩水，⛔ 读 env／⛔ 台面开关）。
凭据＝两台面 × 改前/改后 × 两发（7 枚指名／16 枚全量）＝8 发，＋反形 3 发（仓外树内做变异，逐发还原落 `sha256sum`），＋门禁三把 `rc=`，＋格式名册两把并排。
