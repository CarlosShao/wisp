# 301 orch — CI 色回填（票 301 `AC#2` 欠的那半句，编排者本人跑）

时刻＝`2026-10-10 16:0x +0800`　锚＝HEAD `a0e55a51`（`git log -1 --format='%h %ad' --date=iso-strict` 见本件末「起手／交件双锚」）
射程＝**只新建本件＋`probes/301/orch/logs/**`**；⛔ 产码、⛔ 票面/台账/HANDOVER（落账归编排者另一笔）、⛔ push 之外的写面

## 0. 这格过去欠着什么

`301-v1` 的判语（`probes/301/v1/50-verdict.md` §3-②/§2-E）翻过 `AC#2` 时留了一句：**「欠一枚具名 CI 色＝`test-windows` 在托管 runner 上的首跑，归编排者名下」**。
我 15:5x 想跑，三次抓 GitHub API 遇 `EOF` / `TLS handshake timeout`，那半句就没销。本件是那枚读数的正文。

## 1. 载具（两发同机对照，尺＝`gh api repos/CarlosShao/wisp/actions/jobs/<id>/logs`）

| | run | job | sha | 那笔代码日期 |
|---|---|---|---|---|
| 改前基线 | `38004428359` | `114069831344`（`test-windows`） | `cc315261` | 2026-10-06（schedule 于 10-09 23:26 跑） |
| 改后 | `38034689386` | `114162576251`（`test-windows`） | `bcd0a543` | 2026-10-10，含 `AC#1` 那笔 `0a0f62ef` |

⚠ **基线为什么挑这一发**：origin/dev 的 tip 在 10-06 停在 `cc315261`，今天第一次推送一次带了 **742 枚 commit**（尺＝`git rev-list --count cc315261..bcd0a543`＝742）。⇒ 盘上**不存在**「`AC#1` 前一笔」的 CI 读数；能拿的最接近对照就是同一台 runner 上的 `cc315261`。这一枚⛔ 是「`AC#1` 前一笔」，它是「本批之前」。两件事在本件里全程分开写。
`git merge-base --is-ancestor 0a0f62ef bcd0a543` ⇒ YES（改后那发确实带着 `AC#1`）。

## 2. `--scope=windows`（`ci.yml` 里那步 `Portable windows tests (…)`）

four numbers（尺＝该步日志里 `portable-tests.sh: four numbers` 那一行，逐字引）：

- 基线 `RUN=577  --- PASS=386  --- FAIL=12  --- SKIP=1`
- 改后 `RUN=624  --- PASS=427  --- FAIL=12  --- SKIP=1`

派生量（本件自己算，⛔ 引腿的数）：A-evaluated＝PASS+FAIL+SKIP ⇒ 基线 **399** ↔ 改后 **440** ＝ **+41**。
这与我在腿交件前（14:3x，`probes/301/orch/2026-10-10-baseline-and-delta-prediction.md`）写死的 Δ 预测**同数**，且这次是在 CI 那台机器上，不是本机。

名集合尺（⛔ 只比枚数）：

- 红名册逐名 diff ＝**空**（`diff ci-baseline-fail-names.txt ci-windows-fail-names.txt` rc=0，12↔12 枚枚同名）。
- 改后红名 ∩ `internal/audio` 顶层名册（`probes/301/r1/rosters/audio-toplevel-names.txt`，尺＝`comm -12`）＝**0 枚**。
- 本机（`301-r1` 的 `windows-after`）↔ runner 的逐名色差＝**13 枚**：12 枚本机绿／runner 红 ＋ `TestSyncRedTeamRealOneDrive`（本机跑、runner `--- SKIP`）。反方向 1 枚：`TestC21TableColourRowsMatchTokensCSS` 在 runner 绿、本机不绿（本机那枚的成因是机主工作树里 `design/assets/tokens.css` 脏，早具名过）。
- `RUN` 624 与本机 `windows-after.B-evaluated.txt`（624 行）**逐枚同数**。

⇒ **`AC#1` 的 CI 色＝该档确实开口了**（+47 枚进编译面、+41 枚被求值），**新增红 0 枚**，门的颜色没被这次改动挪过。

## 3. 那 12 枚红是什么（⛔ 算进本票成绩，也⛔ 算本票造成的）

红句逐字引在 `logs/ci-windows-fail-reasons.txt`（尺＝每枚取 `=== RUN` ↔ `--- FAIL` 之间那段，⛔ 取 FAIL 之后——Go 的 `-v` 把明细写在 `--- FAIL:` **之前**，我第一把尺取反了，重跑过）。

| 枚 | 出处文件:行 | 机制那一半 |
|---|---|---|
| `TestClassifyAnchorSpellingIsNotVerdict` | `pathresolver_anchor_spelling_windows_test.go:57` | 红句逐字含 `C:\Users\RUNNER~1` |
| `TestCanonicalInputGainsNoSecondForm` | 同文件 `:140` | 同上（canonical 输入多出第二形） |
| `TestPathResolverShortNameAListDenied` | `pathresolver_junction_windows_test.go:135` | 同上 |
| `TestPathResolverUNCAListDenied` | 同文件 `:164` | 同上 |
| `TestPathResolverExtendedLengthPrefixAListDenied` | 同文件 `:182` | 同上 |
| `TestAListWinsWhereBothTablesHit` | `pathresolver_anchor_spelling_windows_test.go:238` | 红句**不含** `RUNNER~1`，同文件同族，机制**未证** |
| `TestBListDefaultDenyAndOverride` | `pathresolver_junction_windows_test.go:292` | 同上，机制**未证** |
| `TestSyncFixtureFallbackAndMatch` | `syncdirs_test.go:162` | under-profile fallback 家族，机制**未证** |
| `TestSyncFallbackNotDisarmableByWeakRoot` | `syncdirs_test.go:213` | 同上 |
| `TestSyncUnverifiedRootKeepsFallback` | `syncdirs_test.go:233` | 同上 |
| `TestSyncSuspectFallbackIsComponentBounded` | `syncdirs_test.go:355` | 同上 |
| `TestSyncSuspectFallbackWhenUndetectable` | `syncdirs_test.go:402` | 同上 |

5/12 逐字点名 runner TEMP 的 8.3 别名（本机 TEMP 没有别名——这是「本机绿／runner 红」的方向差，与 09-24 那族 `RUNNER~1` 归因同形）。
**这三件文件都远早于本票**：`pathresolver_anchor_spelling_windows_test.go` 建于 `f1033e1e`（09-21），`pathresolver_junction_windows_test.go` 建于 `f088ce3f`（09-20），`syncdirs_test.go` 同族；三件都在 `cc315261` 上就存在（尺＝`git cat-file -e cc315261:<path>` 全 YES）。
⚠ 但「10-06 那批代码在同一台 runner 上也红这 12 枚」**本件没证**（基线那发就是 `cc315261`，它确实红这 12 枚、名集合逐字相同 —— 这一句是证的起；⛔ 证的是「这 12 枚⛔ 是 742 枚里任何一笔带来的」，那需要逐笔 bisect，不在本件射程）。
本件只证到：**它们⛔ 属于票 301**（交集 0 枚 + 名集合改前改后逐字相同 + `AC#1` 那笔只动 `scripts/portable-tests.sh` 两行，逐笔名册尺见 `301-v1` §3-⑥）。

## 4. 顺手量到的一枚**与本票无关**的新红（必须具名入账）

同一 job 里 `cmd/wisp` 那步（`--scope=cli`）也红，且比基线多 2 枚：

- 基线 `RUN=335  --- PASS=231  --- FAIL=6  --- SKIP=1` → 改后 `RUN=398  --- PASS=278  --- FAIL=8  --- SKIP=2`
- 新增红 3 枚（尺＝`comm -13` 两发红名册）：`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`、`TestAC14AwaitedBindingReplyReachesThePage`、`TestAC14GoSideEvalPushReachesThePage`
- 消红 1 枚：`TestPanelHostLatencyPercentilesAC2`（基线红 → 改后 `--- SKIP`）
- 三枚红句同一形状（逐字）：`no report "ac13-probe"/"ac14r-0"/"ac14-push" from the page within 15s (what DID arrive at the door: nothing at all)`；同发上下文里窗体侧是**建成**的（`embed resolves 1068 entry byte(s)`、`entry bytes carried=true`、`probe document is 136 byte(s)`），缺的是**页面自己回话**那一跳。
- 出处文件＝`cmd/wisp/panel_resident_windows_test.go:325/:821/:866`、`cmd/wisp/panel_pageover_33r10_windows_test.go`（票 33 那族 AC#13/AC#14 判据）。
⇒ 判读：**票 33 的真窗判据今天躺在 `--scope=cli` 档里被托管 runner 求值**，而那台机器给不了页面回执 ⇒ CI 恒红。这与 10-08 对票 35 `:75(c)` 的裁定（真窗⛔ 搬进 CI）是同一根，落点归**新票 302**（本件⛔ 动任何产码、⛔ 改任何档）。
⚠ 归因边界：这 3 枚⛔ 是本票造成的（本票两笔产码面＝`scripts/portable-tests.sh` 与 `ci.yml` 注释，一枚没碰 `cmd/wisp`）；它们出现在本批只是因为 742 枚一次进 CI。

## 5. job 色与门色（六枚逐枚对拉）

`gh run view <id> --json jobs` 两发逐枚：`lint`／`test-core`／`test-windows`＝**红**，`lint-frontend`／`slo-smoke`／`slo-full`＝**绿**，基线与改后**同色同枚**。census totals 两发逐字相同＝`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`。
`continue-on-error` 全仓 0 枚真用（尺＝`grep -n "continue-on-error" .github/workflows/ci.yml`，9 处命中枚枚在注释里）⇒ 「步红＝job 红＝run 红」这条读法在本仓是硬的。

最新一笔 `a0e55a51`（含 `AC#5` 的注释笔 `b761b584`）＝run `38035842314`，16:0x 现量：`slo-smoke`／`slo-full`／`lint-frontend` success、`test-core`／`lint` failure、`test-windows` **in_progress**（那半句等它落终态再补，`AC#5` 是纯注释面，预期不变色；**预期⛔ 等于读数**）。

## 6. 这格买到什么（一句话）

`AC#2` 的「CI 现在真会跑它」从今天起**有 CI 自己的字节**：该档被求值面 399→440、新增红 0 枚、门的颜色没动。

## 7. 本发自己的缺口（具名，⛔ 我给自己盖章）

- 我没跑 `git bisect` 去证「这 12 枚与 742 枚里哪一笔相关」——那⛔ 是本票射程，且会在共享工作树上动 checkout（⛔ 禁）。
- `--scope=cli` 那 3 枚新增红的**机制**我只量到「页面不回话」，⛔ 量到「runner 上有没有 WebView2 Runtime」——那枚读数归票 33 的 `33-n1`（在飞）与归口新票 302。
- `a0e55a51` 那发的 `test-windows` 终态本件写时未落；补法＝下一个人重跑第 5 节那条 `gh run view`。
- ⛔ 翻任何框、⛔ 改任何产码／档／门禁、⛔ 动 `:594` 那行 ledger 文案。

## 8. 双锚＋进程台面（现量，⛔ 事后补）

```
[锚] 2026-10-10T16:04:54+08:00
     git log -1  = a0e55a51 2026-10-10T15:51:05+08:00
     git status --porcelain -- scripts internal .github docs frontend  = 0 行
     tasklist IMAGENAME eq wisp.exe      = No tasks are running which match the specified criteria (0 行)
     tasklist IMAGENAME eq balldebug.exe = 同上 (0 行)
```

⚠ 口径写清：本发的第 2 枚锚与交件锚**同一枚时间窗**（16:0x），因为本件全程只读 `gh api` ＋读盘上日志，中间没跑任何长跑、⛔ 动过工作树 ⇒ 两枚读数之间的漂移窗口今天⛔ 存在，但这条**是我判的、⛔ 是量出来的**，下一个人若要引用请按「一发读数」读。
`gh` 的失败重试记录：起手 `gh run view 38034689386 --json jobs` 第一次回 `failed to get run: ... EOF`，第二次成功（原始字节都在 `logs/`）。

## 9. 欠账①已销：`AC#5`（注释笔 `b761b584`）落地后的 CI 读数（16:1x 现跑，⛔ 预测）

本件 §5 结尾那句"`a0e55a51` 那发的 `test-windows` 终态本件写时未落"到 16:1x 落了。尺＝`gh api repos/CarlosShao/wisp/actions/jobs/114165960653/logs`（run `38035842314`，sha `a0e55a51`＝`b761b584` ＋我的落账笔），原始字节＝`logs/job-114165960653-ac5.log`（6,163 行／915,516 B）。

- `--scope=windows` four numbers＝`RUN=624  --- PASS=427  --- FAIL=12  --- SKIP=1` ⇒ 与 §2 那发（`bcd0a543`）**逐字相同**。
- 该档 `--- FAIL` 名册 12 枚，`diff` 对 `ci-windows-fail-names.txt` ＝**空**（`logs/ci-ac5-fail-names.txt`）。
- `--scope=cli`＝`RUN=398 PASS=278 FAIL=8 SKIP=2` ⇒ 同样逐字相同（票 302 那三枚仍在）。
- census totals＝`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0` ⇒ 逐字相同（GUARD D 没被顺手改动，这条正是票 301 `AC#5` 形状⑥要的那把尺）。
- 六枚 job 色＝`slo-smoke`/`slo-full`/`lint-frontend` 绿、`test-core`/`test-windows`/`lint` 红 ⇒ 与基线（§1、§5）逐枚相同。

⇒ **判语（这格仍归 `301-v3`，本节只交读数）**：`AC#5` 那笔**零变色**，与"纯注释面"的判据形状一致；`A817` §6 当时写的"预期不变色"从今天起是**量过的**，⛔ 是预期的。
