# 307-r1 · 10-ac1-colours —— 两台面名册颜色（改前／改后，各两发，同台面取交集）

尺＝`powershell -NoProfile -ExecutionPolicy Bypass -File .scratch/wisp/probes/307/r1/logs/run-roster.ps1 -Tag pre|post`
（脚本逐字在 `logs/run-roster.ps1`；每发＝`go test ./internal/risk/ -count=1 -v`，射程目录＝`./internal/risk/`，工作树）
四数口径＝`^\s*--- PASS|FAIL|SKIP` 与 `^\s*=== RUN` 的 `Select-String` 计数（**含子测试行**），逐发落在 `logs/counts-summary.txt`。

台面：
- `default` = 本机默认（`TMPDIR` 清空、`TEMP`/`TMP` 回原值＝`C:\Users\swq\AppData\Local\Temp`，`40-machine-alias.md` §1 已证全链 SHORT==LONG）
- `aliased` = `$env:TMPDIR/TEMP/TMP` = `C:\Users\swq\tmp\RISKAT~1\RUNNER~1`（真 8.3 别名短形，长形 `...\riskattrib1\RunnerTempAreaForAttribution1`，⛔ 删、⛔ 重建，仍在盘上）
- 每发的 env 实读回落在 `logs/env-<name>-<table>-<tag>.txt`

## 1. 整包名册（`-v`，全 198/200 枚）

| 发 | 台面 | rc | RUN | PASS | FAIL | SKIP | 秒 | 件 |
|---|---|---|---|---|---|---|---|---|
| full-pre | default | 0 | 198 | 197 | 0 | 1 | 6.95 | `logs/full-default-pre.txt` |
| full2-pre | default | 0 | 198 | 197 | 0 | 1 | 3.20 | `logs/full2-default-pre.txt` |
| full-pre | aliased | 1 | 198 | 185 | **12** | 1 | 3.37 | `logs/full-aliased-pre.txt` |
| full2-pre | aliased | 1 | 198 | 185 | **12** | 1 | 3.50 | `logs/full2-aliased-pre.txt` |
| full-post | default | 0 | 200 | 199 | **0** | 1 | 5.42 | `logs/full-default-post.txt` |
| full2-post | default | 0 | 200 | 199 | **0** | 1 | 3.65 | `logs/full2-default-post.txt` |
| full-post | aliased | 1 | 200 | 192 | **7** | 1 | 3.40 | `logs/full-aliased-post.txt` |
| full2-post | aliased | 1 | 200 | 192 | **7** | 1 | 3.50 | `logs/full2-aliased-post.txt` |

RUN 198→200 = 本票新增的 **2** 枚 Windows 夹具（`TestSync307AliasHomeUnderProfileStaysSuspect`、`TestSync307AliasShortRootMembership`），⛔ 搬动既有任何一枚。

### 交集与"单发红过的那一枚"

- aliased 台面改前名册（12 枚，两发**逐字同集合**）＝
  `TestClassifyAnchorSpellingIsNotVerdict`、`TestCanonicalInputGainsNoSecondForm`、`TestAListWinsWhereBothTablesHit`、
  `TestPathResolverShortNameAListDenied`、`TestPathResolverUNCAListDenied`、`TestPathResolverExtendedLengthPrefixAListDenied`、
  `TestBListDefaultDenyAndOverride`（＝**甲 7，归票 72**）
  ＋ `TestSyncFixtureFallbackAndMatch`、`TestSyncFallbackNotDisarmableByWeakRoot`、`TestSyncUnverifiedRootKeepsFallback`、
  `TestSyncSuspectFallbackIsComponentBounded`、`TestSyncSuspectFallbackWhenUndetectable`（＝**乙 5，本票**）
- aliased 台面改后名册（7 枚，两发同集合）＝甲 7，⛔ 一枚不少⛔ 一枚不多。
- **改后 − 改前 的新增红＝0 枚**（两台面、两发都是）。尺＝`comm -13 <post 集> <pre 集>`（`/c/Users/swq/tmp/pre7.txt` ⇄ `post7.txt`）。
- **单发红过、被交集藏住的那一枚＝0 枚**（发与发之间⛔ 任何集合差）。
  ⚠ 但有一枚**发间读数差**必须具名、⛔ 被"集合相同"糊过去：`TestSyncFallbackNotDisarmableByWeakRoot` 用的是 `t.Errorf`（不中断循环），
  两发的 **FAIL 集合相同**，而它体内 `:213` 那行的**打印次数**在改前 aliased 两发里都是 4 次（`"default"`/`"fixture"`/`"options"`/`""` 四个 grade 全红）＝两发一致；
  改后 aliased 两发该枚**全绿**（`:212`/`:209`/`:216` 三处断言⛔ 一条都不再打印）。⇒ 无"只红一发"的枚，这一格交的是"0 枚"这个读数本身。

## 2. ★AC#1 判据逐字：5 枚指名用例在**真别名台面**上全绿

尺＝`grep -E '^\s*--- (PASS|FAIL): TestSync…'`，逐字取自四份件：

| 用例 | aliased 改前 | aliased 改后 | default 改前 | default 改后 |
|---|---|---|---|---|
| `TestSyncFixtureFallbackAndMatch` | FAIL 0.01s | **PASS 0.01s** | PASS 0.01s | PASS 0.01s |
| `TestSyncFallbackNotDisarmableByWeakRoot` | FAIL 0.03s / 0.01s（两发） | **PASS 0.02s / 0.01s** | PASS 0.01s / 0.02s | PASS 0.01s / 0.01s |
| `TestSyncUnverifiedRootKeepsFallback` | FAIL 0.01s / 0.00s | **PASS 0.01s / 0.00s** | PASS 0.01s / 0.00s | PASS 0.00s / 0.00s |
| `TestSyncSuspectFallbackIsComponentBounded` | FAIL 0.01s / 0.00s | **PASS 0.01s / 0.00s** | PASS 0.01s / 0.00s | PASS 0.00s / 0.00s |
| `TestSyncSuspectFallbackWhenUndetectable` | FAIL 0.00s | **PASS 0.01s / 0.00s** | PASS 0.01s / 0.00s | PASS 0.01s / 0.00s |
| 本票新夹具 ×2（仅 Windows 档） | 未存在 | **PASS 0.01s / 0.01s**（两台面四发全绿） | 未存在 | PASS 0.01s / 0.01s |

改前 aliased 的 5 枚红句逐字（件＝`logs/full-aliased-pre.txt`、`logs/full2-aliased-pre.txt`，两发逐字相同）：

```
syncdirs_test.go:162: write under fixture root must be sync
syncdirs_test.go:213: source "default": under-profile fallback must stay armed
syncdirs_test.go:213: source "fixture": under-profile fallback must stay armed
syncdirs_test.go:213: source "options": under-profile fallback must stay armed
syncdirs_test.go:213: source "": under-profile fallback must stay armed
syncdirs_test.go:233: suspect fallback must stay armed
syncdirs_test.go:355: under-profile must be suspect
syncdirs_test.go:402: under-profile path must be sync-suspect, got {Sync:false Root:{Provider: Path: Source:} Why:write target is not under any sync root}
```

⇒ 票面 §1 那枚编排者自己的红句（`:402`）**逐字复现**。

## 3. 本机默认台面颜色不变

改前/改后 default 都是 `FAIL=0`、`SKIP=1`；`comm -13` 尺给新增红＝**0 枚**；`RUN` 只多本票新增那 2 枚。
⇒ `AC#1` 的"⛔ 把别的台面弄红"这一条，在**这台机器给得出的两枚台面**上成立。
⚠ 该读的边界具名：default 台面上那 5 枚**改前就绿**（`40-machine-alias.md` §3 的"缺世界"），所以 default 的颜色**⛔ 能当修复的凭据**，它只证明"⛔ 新增红"。
⚠ 另有一枚**⛔ 属本票射程**的台面形状被现量到但⛔ 动它：`full2-aliased-pre.txt` 里 `syncdirs_test.go:491`/`:509` 打印
`P12 evidence: detected OneDrive root via env: C:\Users\swq\OneDrive` ⇒ 这台机器 `TestSyncDetectionOnThisMachine`（live env，唯一一枚吃 `syncdirs.go:88` 的 `:501` 断言）
的 `SyncDetectionComplete()` 为 **true** ⇒ `:500` 那枚闸在本机恒不执行。具名欠读数，⛔ 本票改（票 72／票 302 射程）。

## 4. `AC#4`／编排者点名要的**唯一新写法的代价**（`cand` → `pathForms` 多一次祖先走步）

尺＝`go test ./internal/risk/ -count=1 -v -run 'TestSync'`（脚本里 `syncrun` 那一发；同一台面各一发，`-run 'TestSync'` 射程＝36→38 枚）：

| 台面 | 发 | 改前总耗时 | 改后总耗时 | Δ |
|---|---|---|---|---|
| default | syncrun | **1.09 s** | **1.15 s** | +0.06 s |
| aliased | syncrun | **1.08 s** | **1.05 s** | −0.03 s |

具名用例逐枚 `--- PASS (0.xx s)` 并排（同一台面 aliased，四发里取两发）：

| 用例 | 改前 aliased | 改后 aliased | 改前 default | 改后 default |
|---|---|---|---|---|
| `TestSyncFixtureFallbackAndMatch` | FAIL 0.01s | PASS 0.01s | PASS 0.01s | PASS 0.01s |
| `TestSyncFallbackNotDisarmableByWeakRoot` | FAIL 0.03s | PASS 0.02s | PASS 0.01s | PASS 0.03s（full）/0.01s（full2） |
| `TestSyncUnverifiedRootKeepsFallback` | FAIL 0.01s | PASS 0.01s | PASS 0.01s | PASS 0.01s |
| `TestSyncSuspectFallbackIsComponentBounded` | FAIL 0.01s | PASS 0.01s | PASS 0.01s | PASS 0.01s |
| `TestSyncSuspectFallbackWhenUndetectable` | FAIL 0.00s | PASS 0.01s | PASS 0.01s | PASS 0.01s |

读数＝**枚级耗时 0.00–0.03 s 带⛔ 可分辨地变化**（`-v` 的钟是 0.01 s 粒度，本票那枚走步被粒度吃掉了）；
包级总耗时两台面各自 ±0.06 s 内、且 aliased 那一发改后**更快**，⛔ 支持"代价可测"这一说⇒ **本票⛔ 主张代价为零，也只主张"在 0.01 s 粒度下量不出来"**。
⛔ 引用编排者或 `307-a1` 的任何数（那两枚都⛔ 跑过编译面）。

## 5. 本文件⛔ 做到的

- ⛔ 勾 `AC#1`/`AC#2`/`AC#3`（裁决者≠实现者）。
- ⛔ 覆盖甲族那 7 枚（票 72 射程，本票⛔ 动其期望侧）。
- ⛔ 声称"CI runner 台面色"：那需要非实现腿或编排者在托管 runner 上取。
