# 300-a2 · 90 — 终态锚＋名册差集＋欠读数（交回件）

- 起手锚＝4e00357f 2026-10-10T17:34:47+08:00（00-anchor.md）
- 终态锚与逐笔名册、差集＝**由尺现跑追加在本行下面**（同形＝git log -1 / git status --porcelain 计数 / tasklist 三枚计数；名册＝对 起手锚..HEAD 逐笔 git show 型 --name-only）
- 授权名册＝probes/300/a2/** 之内；⛔ 越界一枚、⛔ 动产码、⛔ push
- ⚠ 自我引用一处：本件的**最后一笔 commit**（追加终态读数那一笔）⛔ 能把⛔ 自己的号写进⛔ 自己 ⇒ 那一枚号在腿的回报正文里给，⛔ 编排者按 git log 现跑核对。

## 三件交付物（各带尺，见文件）

| 件 | 是什么 | 尺面 |
|---|---|---|
| 10-candidate-landing-points.md | `AC#6` 候选落点 **5 枚**（C1 既有 tagged 用例加段／C2 新枚同包 tagged 文件／C3 权威表读盘／C4 跨实现一致性／C5 生产侧抽谓词），逐枚写"动哪枚文件／判据形状／3→1 与 3→0 红⛔ 红的**依据**／买到⛔ 买到" | 五枚硬事实 F1–F5，全 blob 层 |
| 20-near-miss-nails.md | 擦边钉 **7 枚**（N1–N7）逐枚"断的是什么／撞⛔ 撞"，含"会撞的只有 N4（动容差/token 值）与 N6（新增未登记 SKIP 要改 scripts）两形" | 7 把名册级 grep，跨 internal/** ＋ cmd/wisp/** ＋ scripts ＋ tools/d22scan |
| 30-ci-visibility.md | 答复＝**落在 CI 上**，档＝test-windows 的 --scope=windows 那一发，OS＝windows-latest（托管），step 级 if 逐字为 !cancelled() 故**有求值**；ubuntu 档永远⛔ 求值它（GUARD A 注释自己写着） | portable-tests.sh 名册 × ci.yml job/runs-on/调用点，调用点一律内容锚 |

## 欠读数（具名，⛔ 静默；本腿⛔ 可取）

1. **新判据的成对两发颜色**（改⇒红／逐字节还原⇒绿）＝**归落地腿**（本腿⛔ 编译面授权：`go build/test/vet/run/list` 全部 0 发；`go env` 亦 **0 发**，此处自报）。
2. **CI 那一发色**（test-windows 档跑起来是 PASS/FAIL/SKIP）＝**归编排者**，票面 next 已钉"跟票 303 修复后那批推送一起取、⛔ 为它单推"。
3. **托管 windows 镜像上 SDK 头文件在不在**（候选 C3 的生死判）＝**归那一发 2**；仓内名册级尺对 Windows Kits/mmreg.h 在 scripts/tools/.github 的命中＝**0 枚**（本腿现跑）。
4. **internal/audio 在 windows 档的当前基线色**（票 301 AC#2 那发前后名册作差）＝**归编排者／301 的腿**；probes/301 在我授权名册外，本腿只按票面引、⛔ 查。
5. **TestLiveWasapiSmoke 在托管 runner 上到底是登记过的 SKIP 还是红**＝**归那一发 2**（ledger 那行的理由句自称托管无音频端点，⛔ 实测凭据在盘上）。
6. **cmd/wisp 那族 255 名册仪器会不会去看 audio 的字节面**＝**归 cmd/wisp 台面持有腿（303-r1）或编排者**；本腿只量到"⛔ 人按行号引 wasapi_windows.go／parse_wave_format_300*"（0 命中），⛔ 证明那族仪器⛔ 看字节面。
7. **顶回一处枚数口径**：派单/票面那句"整包 **41** 枚顶层"与本腿源码尺（blob 名册逐枚统计 func Test 行）数到的 **42 枚**差 1，而该包 func TestMain 命中 **0** ⇒ 两把尺⛔ 同形（PASS 行 vs 源码行），⛔ 我判谁对，落笔请带尺名。
8. **顶回一处派单措辞**：票面 AC#6 原文还有一半边派单⛔ 写——"且红句具名指向那枚常量而⛔ 指向 tag 的等值比较"；本腿按**原文**办（原文优先于转述），差异具名报回。
9. ⛔ 派单一处**⛔ 我拍**：C4（跨实现一致性）撞票面 AC#0 那句"⛔ 拿本仓另一枚实现当裁判"的**字面射程**——那一格钉偏移、本候选钉**值**，算⛔ 同一条禁令＝编排者裁；C3 撞既有注释声明＋要改 scripts（越出 AC#4 名册）＝编排者裁；C5 改产码语义面＝⛔ 本票授权范围，需编排者明示。

## 追加：终态读数（尺现跑，取于本件入库那一发之前）

anchor-start:
4e00357f 2026-10-10T17:34:47+08:00
porcelain-mine:
1
porcelain-full:
823
tasklist:
wisp.exe=0 balldebug.exe=0 msedgewebview2.exe=24 
tip:
62688150 2026-10-10T17:51:14+08:00
rc-last=$?

## 追加二：逐笔名册与差集（尺现跑，--no-renames；取于 tip=62688150 之后）

roster-of-my-commits (4e00357f..HEAD, --no-renames):
COMMIT 62688150 2026-10-10T17:51:14+08:00 probes/300/a2: final piece 90 plus remaining log files (landing survey complete)

.scratch/wisp/probes/300/a2/90-final.md
COMMIT ae8ebfbd 2026-10-10T17:50:19+08:00 303-r1: AC#2 机制三形分开答（活动形=丙，有硬读数）

.scratch/wisp/probes/303/r1/10-culprit-prod-diff.txt
.scratch/wisp/probes/303/r1/14-ac2-shape3-reentry-read.txt
.scratch/wisp/probes/303/r1/16-ac2-mechanism-three-shapes.md
.scratch/wisp/probes/303/r1/20-targeted-before-red.txt
.scratch/wisp/probes/303/r1/msg-ac2.txt
COMMIT 3d0f8ccc 2026-10-10T17:49:55+08:00 probes/300/a2: AC#6 landing-point survey (candidates, near-miss nails, CI visibility)

.scratch/wisp/probes/300/a2/10-candidate-landing-points.md
.scratch/wisp/probes/300/a2/20-near-miss-nails.md
.scratch/wisp/probes/300/a2/30-ci-visibility.md
.scratch/wisp/probes/300/a2/logs/anchor-rot.txt
.scratch/wisp/probes/300/a2/logs/final-anchor-raw.txt
.scratch/wisp/probes/300/a2/logs/near-miss.txt
.scratch/wisp/probes/300/a2/logs/wavinjector-and-scopes.txt
COMMIT e4d9a223 2026-10-10T17:48:22+08:00 A821 取证否证我自己那笔"24 枚孤儿 webview"的欠账（那是机主四个应用的进程树）＋★同一次抓到真残留立票 304（测试二进制被外力杀 ⇒ 子进程不随父退，2 枚 mockllm.exe 已挂 3420 分钟）

.scratch/wisp/issues/304-test-binaries-killed-from-outside-orphan-their-children-two-mockllm-exe-still-listening-on-the-owners-machine.md
.scratch/wisp/probes/orch/webview-orphans/forensics-1.txt
.scratch/wisp/probes/orch/webview-orphans/forensics-2.txt
.scratch/wisp/probes/orch/webview-orphans/forensics-3.txt
.scratch/wisp/probes/orch/webview-orphans/forensics.ps1
.scratch/wisp/probes/orch/webview-orphans/forensics2.ps1
.scratch/wisp/probes/orch/webview-orphans/forensics3.ps1
.scratch/wisp/probes/orch/webview-orphans/kill-mine-1.txt
.scratch/wisp/probes/orch/webview-orphans/kill-mine.ps1
.scratch/wisp/probes/orch/webview-orphans/msgs/A821.txt
.scratch/wisp/probes/orch/webview-orphans/msgs/cm-1.txt
.scratch/wisp/probes/orch/webview-orphans/msgs/handover-40bs.txt
docs/reports/HANDOVER.md
docs/reports/pending-and-issues.md
COMMIT ae5ee86d 2026-10-10T17:41:30+08:00 probes/300/a2: anchor deposit 1 (HEAD-object-layer greps only, no readings verdict yet)

.scratch/wisp/probes/300/a2/00-anchor.md
.scratch/wisp/probes/300/a2/logs/ci-guards.txt
.scratch/wisp/probes/300/a2/logs/ci-jobs.txt
.scratch/wisp/probes/300/a2/logs/ci-visibility.txt
.scratch/wisp/probes/300/a2/logs/grep-constants.txt
.scratch/wisp/probes/300/a2/logs/roster-audio.txt
.scratch/wisp/probes/300/a2/logs/sec-B-and-prod.txt
.scratch/wisp/probes/300/a2/logs/sec-B2-prod.txt
.scratch/wisp/probes/300/a2/snap/ledger.md
.scratch/wisp/probes/300/a2/snap/parse_wave_format_300_windows_test.go.txt
.scratch/wisp/probes/300/a2/snap/ticket-300.md
.scratch/wisp/probes/300/a2/snap/v3-20-teeth.md
COMMIT 03fa5840 2026-10-10T17:36:43+08:00 303-r1: 起手锚 (HEAD 4e00357f, 台面 0 脏, 进程 wisp/balldebug=0)

.scratch/wisp/probes/303/r1/00-anchor.md
diff-set (mine minus authorized roster probes/300/a2), count:
20
authorized-roster file count:
20
my-file-sizes:
2463 .scratch/wisp/probes/300/a2/00-anchor.md
11603 .scratch/wisp/probes/300/a2/10-candidate-landing-points.md
7240 .scratch/wisp/probes/300/a2/20-near-miss-nails.md
7289 .scratch/wisp/probes/300/a2/30-ci-visibility.md
7575 .scratch/wisp/probes/300/a2/90-final.md
56919 .scratch/wisp/probes/300/a2/logs/anchor-rot.txt
18936 .scratch/wisp/probes/300/a2/logs/ci-guards.txt
1309 .scratch/wisp/probes/300/a2/logs/ci-jobs.txt
3114 .scratch/wisp/probes/300/a2/logs/ci-visibility.txt
166 .scratch/wisp/probes/300/a2/logs/final-anchor-raw.txt
10676 .scratch/wisp/probes/300/a2/logs/grep-constants.txt
15407 .scratch/wisp/probes/300/a2/logs/near-miss.txt
6052 .scratch/wisp/probes/300/a2/logs/roster-audio.txt
10889 .scratch/wisp/probes/300/a2/logs/sec-B-and-prod.txt
6694 .scratch/wisp/probes/300/a2/logs/sec-B2-prod.txt
30135 .scratch/wisp/probes/300/a2/logs/wavinjector-and-scopes.txt
4250525 .scratch/wisp/probes/300/a2/snap/ledger.md
7492 .scratch/wisp/probes/300/a2/snap/parse_wave_format_300_windows_test.go.txt
72382 .scratch/wisp/probes/300/a2/snap/ticket-300.md
13433 .scratch/wisp/probes/300/a2/snap/v3-20-teeth.md
go-compile-face self-report: go build/test/vet/run/list = 0 calls; go env = 0 calls
rc-line: rc=0

## 追加三：更正追加二那把尺的射程（⛔ 抹原句，只追加；出处＝本腿自查）

原句diff-set count: 20＝**射程打错了**：追加二那条 git log 的范围 4e00357f..HEAD ⛔ 只含本腿四笔，它把中间**别腿的两笔**一起算进来了（e4d9a223 A821 台账笔／ae8ebfbd 303-r1 产码笔，正合在飞写腿一 commit，HEAD 就会动那句派单警告）。⇒ 20 枚＝那两笔的文件，⛔ 本腿写的。

正确尺＝对本腿**逐笔** show 名册（--no-renames）再作差：
本腿名册总枚数(去重):
24
越界枚数(名册减去授权 probes/300/a2/**):
4
逐笔名册（去重后）:
-- commit 3d0f8ccc:
-- commit 62688150:
-- commit ae5ee86d:
-- commit bb314f3b:
.scratch/wisp/probes/300/a2/00-anchor.md
.scratch/wisp/probes/300/a2/10-candidate-landing-points.md
.scratch/wisp/probes/300/a2/20-near-miss-nails.md
.scratch/wisp/probes/300/a2/30-ci-visibility.md
.scratch/wisp/probes/300/a2/90-final.md
.scratch/wisp/probes/300/a2/logs/anchor-rot.txt
.scratch/wisp/probes/300/a2/logs/ci-guards.txt
.scratch/wisp/probes/300/a2/logs/ci-jobs.txt
.scratch/wisp/probes/300/a2/logs/ci-visibility.txt
.scratch/wisp/probes/300/a2/logs/final-anchor-raw.txt
.scratch/wisp/probes/300/a2/logs/grep-constants.txt
.scratch/wisp/probes/300/a2/logs/near-miss.txt
.scratch/wisp/probes/300/a2/logs/roster-audio.txt
.scratch/wisp/probes/300/a2/logs/sec-B-and-prod.txt
.scratch/wisp/probes/300/a2/logs/sec-B2-prod.txt
.scratch/wisp/probes/300/a2/logs/wavinjector-and-scopes.txt
.scratch/wisp/probes/300/a2/snap/ledger.md
.scratch/wisp/probes/300/a2/snap/parse_wave_format_300_windows_test.go.txt
.scratch/wisp/probes/300/a2/snap/ticket-300.md
.scratch/wisp/probes/300/a2/snap/v3-20-teeth.md
rc=0
