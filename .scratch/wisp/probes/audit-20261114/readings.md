# audit-20261114 原始读数（本审计跑的每一条尺的现取输出）

审计者 `agent-26`。本件是 `docs/evidence/s1/audit-20261114-closures.md` 的读数底稿；
所有数都是 2026-11-14 17:3x–17:5x +08 现跑，**没有一条是转述票面**。

## 0. 树锚

```
$ git rev-parse HEAD
3e9d9779fabab80c2c7c1f26559332d466374b53      # 起点；落笔时 HEAD 已到 81ba68874

$ git diff --name-status 3e9d9779f HEAD
A   .scratch/wisp/probes/orch/20261114/txt/survey-pool-dispatch.txt
M   docs/reports/pending-and-issues.md
# ⇒ 5 张票的 blob 在此区间未变

$ git status --porcelain | group
 5 M / 1725 ??      # 他人在飞件与 probes 未跟踪件，本审计未碰
```

## 1. HEAD 双计数（现树）

```
$ git ls-tree -r --name-only HEAD -- .scratch/wisp/issues | grep -E '/(441|442|274|443|444)-'
.scratch/wisp/issues/274-no-nail-shipped-exe-carry-page-build-step-2-comment-false-done.md
.scratch/wisp/issues/441-posix120-seal-window-instrument-reds-on-abstention-under-ci-contention-done.md
.scratch/wisp/issues/442-panel-profile-ruler-still-misses-two-interprocedural-shapes-done.md
.scratch/wisp/issues/443-test52-timeout-tree-instrument-reds-when-fixture-never-reports-a-pid-done.md
.scratch/wisp/issues/444-seven-more-race-window-instruments-can-red-on-abstention.md
```

每票恰 1 枚。磁盘同样 1 枚。

## 2. HEAD 双计数（历史窗口，逐 commit 扫 tree）

尺＝`dupescan.ps1`（同目录）；输出（只列 count>1 的 commit）：

```
e4f5baa56 2026-11-14 15:42:45 +0800  DOUBLE: 441=2
f1f3533ab 2026-11-14 15:43:35 +0800  DOUBLE: 441=2
36e812759 2026-11-14 15:46:18 +0800  DOUBLE: 441=2
db47f9c36 2026-11-14 15:48:34 +0800  DOUBLE: 441=2
5209a38e7 2026-11-14 15:54:18 +0800  DOUBLE: 441=2
936bbaa5b 2026-11-14 16:00:32 +0800  DOUBLE: 441=2
813bf8574 2026-11-14 16:26:12 +0800  DOUBLE: 274=2
```

修：`ae4e506ee` 16:01:48（441，`D` 旧路径）、`c0e48dbc4` 16:26:32（274，`D` 旧路径）。
442／443／444 从未出现在这张表里。

## 3. 441：引证与 ancestry

```
$ git log -1 --format='%h %p' 8e1367bfd
8e1367bfd eb9cf0e63

$ git cat-file -t eb9cf0e63 a46352fd3 8e1367bfd b412f50e4 952f4462b 2e8233aeb
commit ×6

$ git show -s --format='%h %ci %s' a46352fd3 eb9cf0e63
a46352fd3 2026-11-14 14:55:16 +0800  ticket 442 opens ...
eb9cf0e63 2026-11-14 15:22:51 +0800  ticket 442 pivots ...
（修复主体 8e1367bfd = 15:23:38）
```

CI（`gh run view <id> --json headSha,conclusion,jobs`；完成时刻取 `gh run list --json updatedAt`）：

| run | head | 创建 | **完成** | run 结论 | 逐 job |
|---|---|---|---|---|---|
| 43565131755 | 05fb21317 | 14:11 | — | failure | 只 test-core=failure，其余 5 job success |
| 43568643680 | a46352fd3 | 14:55 | — | success | 6/6（**修复前**） |
| 43570834384 | eb9cf0e63 | 15:22 | **15:41:19** | success | 6/6（**修复的父提交**） |
| 43570813207 | 650eff6cf | 15:22 | **15:42:09** | success | 6/6（**修复前**） |
| 43571235710 | 952f4462b | 15:28 | **15:47:00** | success | 6/6（**含修复；晚于结案 15:42:45**） |
| 43571511743 | 2e8233aeb | 15:31 | **15:47:39** | success | 6/6（含修复） |
| 43572043677 | b30cc7cf2 | 15:38 | **15:57:33** | failure | test-windows |

⇒ 结案（15:42:45）那一刻**已跑完**的两发 = `eb9cf0e63` + `650eff6cf`，**都在修复之前**；含修复的第一发 15:47:00 才完成。

## 4. 441：票面字节

```
bytes=7127  NUL=3（偏移 213 / 726 / 854）  bare-CR=3  LF=67
偏移 213 : 出处行  "head \<NUL>5fb21317\"
偏移 726 : 普查行  "同一套 \<TAB>est-core\ 在 \<NUL>69f97cac\ ／..."
偏移 854 : 立案行  "触发这一发的 commit \<NUL>5fb21317\"
另有      : AC#5 那行  "\ash -n\"      # 原文 \bash -n\（b 被吃掉）
            /bin/bash 之外，行号不可引用（3 枚裸 CR 会让不同读法数出不同行数）
$ git show --stat 8e1367bfd | grep 441-posix
 ...ument-reds-on-abstention-under-ci-contention.md | Bin 3031 -> 4211 bytes
```

## 5. 441：census 复算（票面 AC#4 说 37 行 / 13 文件）

```
$ git grep -n -E 'RED instrument|instrument is dead|RED INSTRUMENT' 8e1367bfd -- internal
-> 37 行（全为 *_test.go）、13 个文件     # 与票面逐数一致
$ sh .scratch/wisp/probes/441/census.sh     # 当前 HEAD 树
rows: 43   files: 13   rc=0                 # 444 那波加行后变大，非矛盾
```

## 6. 442：字节 / 行号 / CI

```
$ (Get-Item docs/evidence/s1/442-runtime-criterion.md).Length              -> 14139
$ (Get-Item cmd/wisp/panel_data_path_distinct_442_windows_test.go).Length  -> 11536
$ (Get-Item cmd/wisp/panel_data_path_seam_windows.go).Length               -> 1851

$ grep -n 'panelDataPathSeam' cmd/wisp/panel_host_windows.go
335:    if panelDataPathSeam != nil {
336:        panelDataPathSeam(dataPath)

$ grep -n 'SAME data path' cmd/wisp/panel_data_path_distinct_442_windows_test.go
254: t.Errorf("construction point %s handed the host the SAME data path ...

# 票面写的另外几条守卫的真实行号
199: if want := rounds * len(points); len(reads) != want {
216: if len(order) != len(points) {
224: for _, point := range points {          # 票面说这里是「期望文件恰一个点」，实际守卫在 231
231: if len(found) != 1 {
235: }                                        # 票面说这里是「每点每轮都构造」，实际守卫在 240
240: if len(values) != rounds {
248: if value == "" {

$ git show --shortstat --format='%h' 6949484e1
18 files changed, 2369 insertions(+), 146 deletions(-)
$ git show --name-status --format='' 6949484e1 | grep -v 'probes/\|docs/\|issues/'
A  cmd/wisp/panel_data_path_distinct_442_windows_test.go
A  cmd/wisp/panel_data_path_seam_windows.go
M  cmd/wisp/panel_host_windows.go
M  cmd/wisp/panel_profile_leaf_shape_windows_test.go

$ gh run view 43575142588 (6949484e1) -> success, 6/6
$ gh run view 43575466311 (021f15661) -> success, 6/6
$ gh run view 43575529807 (312a58612) -> failure, test-windows / cmd-wisp CLI tests:
    --- FAIL: TestWireR1AnnounceEndToEnd (7.31s)     # 与 442 无关
```

## 7. 274：门禁时点

```
$ gh run view 43575830504 (head 813bf8574) --json jobs
  lint [failure] step: Tracked path-length budget (ticket 262)
  ... 日志逐字：
  check-path-length-budget.sh: denominator: tracked paths=20028  over-budget=70  covered by roster=69  not in roster=1
  check-path-length-budget.sh: RED - over budget and NOT in the roster: .scratch/wisp/issues/274-no-nail-requires-shipped-exe-to-carry-page-build-ps-step-2-comment-is-false-in-both-halves-done.md  (relative 123 chars, name 102 chars, ...)
  check-path-length-budget.sh: VERDICT RED
$ gh run view 43575854808 (head c0e48dbc4) -> lint failure 同一步
$ git show --name-status 6236a6110
  R100 .../274-no-nail-requires-shipped-exe-…-done.md -> .../274-no-nail-shipped-exe-…-done.md   # 102 -> 69 字符
$ gh run view 43577965575 (45ccaf186，改名后下一发) -> failure 仅 test-windows；lint success

$ gh run view 43575142588 / 43575466311 / 43578030753 / 43578228598 --json jobs
  slo-smoke: "Shipped wisp.exe carries the page bundle (panel-assets -check on the artifact)" => success ×4

$ git show --stat 45c68fa79
 .gitignore | 1 -
 .../274-…-is-false-in-both-halves.md | 24 +
 docs/specs/SPEC-01-architecture.md | 4 +---
 docs/specs/SPEC-11-build-deploy-containerization.md | 23 +++------------------
 4 files changed, 28 insertions(+), 24 deletions(-)
```

## 8. 443

```
$ (Get-Item docs/evidence/s1/443-instrument-abstention.md).Length  -> 19537    # 票面自报 19537
$ grep -n 'ABSTAIN(443' internal/plugin/manifestslot/command_52_test.go
738: //     `AC#2 ABSTAIN(443:timeout-tree:raced-no-pid) ...
811: t.Logf("AC#2 ABSTAIN(443:timeout-tree:raced-no-pid) %s pid-seen=no ...
$ grep -n 'RED(443' (same file) -> 830 / 833 / 840 / 858（858 = tree-survived，牙）
$ git show --name-status --format='' 1f0d04714
 A .scratch/wisp/probes/314-ruler2/d22scan-postcommit.log
 A .scratch/wisp/probes/314-ruler2/path-budget-postcommit.log
 A .scratch/wisp/probes/314-ruler2/post-commit-green.log
 M docs/evidence/s1/314-shape-ruler.md            # 与产码零关系，成立

$ go test ./internal/plugin/manifestslot -run 'Test52|Test138' -count=1
ok  github.com/CarlosShao/wisp/internal/plugin/manifestslot  18.701s

# probes/443 里 abstain 计数
03-natural-25.txt               abstain=3
04-amplified-20-delay3000.txt   abstain=20
11-amplified-5-delay1500.txt    abstain=5
05-mut-singlekill-raced.txt     abstain=1 red=1 FAIL=1
13-natural-25-postedit.txt      abstain=0
→ 50 发 3 弃权 = 6%；1500ms 5/5；3000ms 20/20；突变真红

# CI（含 443 修复的 step）
43577965575 (45ccaf186): Portable windows tests (.../plugin/...) => success；该发只红 cmd-wisp 的 Test390r1
43578030753 (721a1d8a0): 6/6 success
```

## 9. 444

```
$ grep -n 'ABSTAIN(444' internal/winsec/*_test.go
openleaf_322_other_test.go:229
ancestoropen_342_unix_test.go:138
ancestoropen_342_windows_test.go:76
ancestorunlink_348_unix_test.go:172
ancestorunlink_348_windows_test.go:128
deleteseam_348r2_windows_test.go:170 / :192          # 共 7 条

# t.Errorf|t.Fatalf 计数（45ccaf186^ -> HEAD）
ancestoropen_342_unix          3 -> 3
ancestoropen_342_windows       2 -> 3
ancestorunlink_348_unix        7 -> 8
ancestorunlink_348_windows    15 -> 16
ancestorworld_342_348_unix    14 -> 23
ancestorworld_342_348_windows 12 -> 16
deleteseam_348r2_windows       4 -> 7
openleaf_322_other             8 -> 13
# 无一枚下降

$ git ls-tree -r --name-only HEAD -- .scratch/wisp/probes/444 | wc -l   -> 22（文件）
  其中 logs/ 下 20 枚、另有 commit-msg.txt 与 repro-444.sh

$ gh run view 43578030753 / 43578127657 / 43578228598 -> success 6/6 ×3（done 17:14:21 / 17:14:15 / 17:17:55）
```

## 10. 门禁现跑（HEAD 树）

```
$ bash scripts/check-path-length-budget.sh
check-path-length-budget.sh: bands: over the hat=69  of which in the 122..180 middle=65  past the old debt line(180)=4  in the wall interval=0  roster entries=69
check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree
rc=0
```

`bash` 在本机默认走 WSL（无 `/bin/bash`，rc=1 的**环境**假红）；真 bash 在
`D:\work\soft\Git\bin\bash.exe`。这一点与 274-w1 票面记的「尺的尺」同族，记下来免得下一个人把它当门禁红。
