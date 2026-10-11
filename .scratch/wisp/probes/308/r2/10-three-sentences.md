# 308-r2 · 10-three-sentences（三句各自：旧句 → 新句 → 我为它跑的尺）

对象＝票 308 §4（编排者 14:2x）点名的三处注释假句／缺射程句。三句都只动**注释行**（尺见 `20` 件：剥掉整行注释后与派发锚 blob `cmp` 同一）。
句子的**建议文本**来自非实现者腿 `308-v1` 的 `probes/308/v1/70-overclaim-audit.md`（#1 那段编排者逐字读过才采纳），本腿**逐枚复核过里面的每个数**，复核不过的那一枚（`(8)`）已当场改数并在下面具名顶回。

---

## 句 1 · `internal/risk/pathresolver_anchor_spelling_windows_test.go`（`TestClassifySpellingInvariance` 文档段末）

旧句（逐字，HEAD blob `52c8fd44`）：
> The 8.3 spelling is only folded in when this volume generates short names (same probe the ticket-18 case uses); **every other spelling is unconditional, so the test never quietly stops measuring.**

新句（要点）：8.3 那句保留；把"其余拼法无条件"换成"自票 308 起两枚 UNC 拼法也是条件的（`driveShare` 报 ¬ok ⇒ `volumeInvariantSpellings` 按名去掉，本用例届时印出保留数与去掉的名）"，然后给出**量过的分母**：卷标齐＝**9**、8.3 命中＝**10**；无卷标输入＝**7**，并写明 7 是被**输入面**那枚用例钉住的、⛔ 是被 env 形状钉的；末句保留"never quietly stops measuring"但把理由换成"a drop is NAMED and LOGGED, not because the spellings are unconditional"。

复核尺与读数（全部 2026-10-11 现跑，件在 `logs/`）：

| 该句里的断言 | 尺 | 读数 | 判 |
|---|---|---|---|
| 卷标齐时保留 9 枚、去 0 枚 | 仓外导出树探针（只打印不断言）`logs/probe-denominator.txt`，默认 temp | `keep=9 [dot-and-dotdot extended long lower not-yet-there slashes unc-admin unc-extended upper] dropped=[]` | ✅ 与 `308-v1/logs/v45` 同值 |
| 同一枚输入在 8.3 别名台面上也＝9 | 同探针第二发 `logs/probe-denominator-alias.txt`（`TMP=RISKAT~1\RUNNER~1`，cmd.exe 绕开 MSYS 改写；`os.TempDir()` 当场证用上） | `keep=9`，名册逐字同 | ✅ |
| "10 with 8.3"＝8.3 槽位会命中 | 探针里 `getShortPath(t, long)`（真路径） | 两把台面都**非空**（`…\TESTR2~1\001\SSH~1\ID_TES~1`）⇒ 用例那行 `spellings["8.3-short"] = s` 命中 ⇒ 9＋1＝**10** | ✅ |
| "7 on a volume-less one" | 探针喂无卷标字面量 `\tmp\wisp-308-no-volume\.ssh\id_testkey` | `VolumeName=""`、`keep=7`、`dropped=[unc-admin unc-extended]` ⇒ 输入面用例自己那行 `t.Logf` 也在整包名册里印着同一句 | ✅ |
| **"（8）"**（无卷标＋8.3＝8） | 探针 `getShortPath(t, 无卷标字面量)` | **`""`** ⇒ 无卷标那形**停在 7，不是 8** | ❌ **顶回**（见末节） |
| "a drop is NAMED and LOGGED" | 读码＋整包名册 | 用例的 `t.Logf` 串逐字含 `dropped the two volume-derived spellings by name %v; %d volume-free spellings still evaluated` ⇒ 去的是名、留的是数 | ✅ |
| "the only kind Windows hands out here" | 我只量了两把台面（默认＋别名），`os.TempDir()` 各＝`C:\Users\swq\AppData\Local\Temp`／`C:\Users\swq\tmp\RISKAT~1\RUNNER~1`，`VolumeName` 都＝`C:` | 编排者采纳的原句是"every volume-qualified temp -- the only kind Windows hands out here"＝**全类断言**；本腿⛔ 复跑 v1 的 env 网格（`v1/logs/v14`/`v50`），⇒ 新句改写成"measured 2026-10-11 on the two temp surfaces reachable here"，⛔ 把两发读数写成通则 | ⚠ 收窄（⛔ 造新承诺） |

## 句 2 · 同文件 ticket 308 注释块

旧句（逐字）：
> ... volumeInvariantSpellings drops them BY NAME and **the remaining seven volume-free spellings still run** -- **the case ends with a verdict**, it does NOT panic the binary and does NOT t.Skip ...

新句（＝采纳 `70-overclaim-audit.md` #2 的收窄，并按我自己的尺具名）：
- **钉住的那一半**：helper 喂无卷标字面量时**返回⛔ panic**、按名去掉两枚 UNC、把七枚无卷标拼法留在 map 里 —— 判据＝`TestVolumeInvariantSpellingsVolumeLess308`（选择那一半）。
- **⛔ 任何夹具钉不住的那一半**：本枚 `//go:build windows` 用例自己的 `len(dropped) > 0` 分支。⛔ 再写"still run / the case ends with a verdict"。
- 保留旧块里**已被量过**的两件事实：⛔ panic 的形状（旧 inline `strings.ToLower(vol[:1])` 逐字红句＋"took every other case in the package's evaluation with it"＝票 §现量 与 `308-v1/20-ac2-teeth.md` §a 的 201⇄35 那两发）；"this file holds zero `t.Skip` calls"。

复核尺与读数：

| 断言 | 尺 | 读数 |
|---|---|---|
| 那一支到达数＝0（默认台面） | 整包 `go test ./internal/risk/ -count=1 -v`：派发锚 `202ed4a5` **仓外导出树**一发（`logs/pre-head-full.txt`）⇄ 改后**工作树**一发（`logs/post-worktree-final-full.txt`），哨兵串 `dropped the two volume-derived` | 两发都＝**0** |
| 同一枚尺在别名台面 | 别名台面整包一发（`logs/pre-head-full-alias.txt`，`RISKAT~1` 字样命中 10 枚＝台面真被用上） | **0** |
| 另一枚哨兵（junction 的 `cannot be constructed`） | 同三发 | **0／0／0**（⚠ 那一支在 junction 文件里，**本腿⛔ 改它**＝票 §4 只点名 anchor 两处＋新件一处） |
| "linux 永不求值这枚文件" | `GOOS=linux GOFLAGS= go list -f '{{range .TestGoFiles}}…'` ⇄ `GOOS=windows` | windows 名册**含** `pathresolver_anchor_spelling_windows_test.go`；linux 名册**⛔ 含**（22 ⇄ 17 枚）⇒ 结构性尺，⛔ 靠跑色推断 |
| "helper 那一支真不 panic" | 输入面用例在整包名册里的颜色 | `--- PASS: TestVolumeInvariantSpellingsVolumeLess308`、`panic: runtime error` 命中 **0**、rc=0 |
| "本文件零 t.Skip" | `grep -cE '^[[:space:]]+t\.Skip\('` | anchor **0**／新件 **0**／junction **1**（＝票 18 那枚既有，⛔ 本腿的、⛔ 计入我的账） |

## 句 3 · `internal/risk/pathresolver_volumeless_spelling_308_test.go`（`driveShare` 文档段）

旧句（逐字）：
> **This is the single place** the drive letter is taken.

新句：`Within this test surface this is the single place the drive letter is taken.` ＋把另外两处具名：产码 `internal/risk/pathresolver.go`（`normalizeLocalUNC` 在 `drive[:1]` 之前逐字有 `len(drive) == 2 && drive[1] == '$'` 守卫）、`internal/tools/volume_windows.go`（`volumeID`／`volumeRootOf` 读 `filepath.VolumeName` 但**完全⛔ 切片**），并给出普查尺＋射程＋"评估位＝2"，⛔ 写会随注释自身腐烂的枚数（那枚"另外 5 行是注释"的 5 在我落笔后变 6 ⇒ 改成"pattern 命中的其余行都是注释"，尺留在句里可重跑）。

复核尺（各一发，`2026-10-11`）：

| 断言 | 尺 | 读数 |
|---|---|---|
| 求值位＝2 | `grep -rn --include=*.go 'vol\[:1\]\|drive\[:1\]' internal cmd tools`（⛔ 注释另筛）⇄ `git grep -e 'vol[:1]' -e 'drive[:1]' -- '*.go'`（整棵跟踪树） | 求值位**只剩 2**：`pathresolver.go:263` ＋ 本 helper；其余命中全在注释行 ⇒ 两把尺同值 |
| 产码确有守卫且仍在切片之前 | 读 `internal/risk/pathresolver.go:262-263` 逐字 | `262` 守卫行、`263` `drive[:1]` ⇒ 编排者 §1 那句复认 |
| `volume_windows.go` 读 `VolumeName` 且⛔ 切片 | `grep -n 'filepath.VolumeName' internal/tools/volume_windows.go` | `:28`（走 `len(vol) != 2 \|\| vol[1] != ':'` 一支）＋ `:54`；文件 65 行 ⇒ `308-v1` 那句 `:28/:54` **是对的** |
| "测试面唯一一处" | 同普查：`internal/risk` 的 `_test.go` 里驱动器字母切片位＝本 helper 一枚（junction 那枚已改经 `driveShare`） | ✅ |

---

## ★顶回 `308-v1` 建议文本的地方（交回时同样写明，编排者复跑再判）

1. **`(8)` 那一枚分母⛔ 成立**：`70-overclaim-audit.md` 建议句写 "and 7 **(8)** on a volume-less one"。本腿用同一枚 `getShortPath` 探针喂那枚无卷标字面量，两把台面都交回 `""`（件 `logs/probe-denominator.txt` 第 7 行／`logs/probe-denominator-alias.txt` 第 7 行，同句 `R2PROBE getShortPath(volume-less literal)="" -> 8.3 slot on volume-less input = 7`）⇒ 无卷标那形**到不了 8**：8.3 槽位与卷标⛔ 独立，它是"真路径＋该卷开 8.3"的函数。⇒ 新句写成 "so it is 7, not 8"，并把量它的尺具名留在句里。
2. **"the only kind Windows hands out here"＝全类断言，本腿只复现了两把台面**：⛔ 复跑 v1 的 env 网格（那是别腿的台面账），所以新句把射程收到 "the two temp surfaces reachable here"。⚠ 这不是说 v1 那句错——它另有两发（`v1/logs/v14`/`v50`）——只是⛔ 让我把⛔ 量过的范围抄进交付件。
3. **"printing kept/dropped" 措辞偏松**：用例印的是"保留的**枚数**＋去掉的**名**"，且只在真去掉时印（`if len(dropped) > 0`）。新句按代码逐字改写，免得下一个人以为名册每次都印。
4. ⚠ **⛔ 顶回但记一笔**：`70-overclaim-audit.md` §另两处 那句 `junction:157-159` "on the windows CI gear"＝⛔ 读数——它落在 `pathresolver_junction_windows_test.go`，一枚**票 §4 ⛔ 点名、本腿⛔ 碰**的文件（那笔账在 `Q-308-v1-A` 的 CI 色里）。本腿⛔ 改它，只在此具名说"句子仍在盘上"。
