# `10-ac1-verdicts.md` — `310-v1` 对票 310 `AC#1`（落地形＝甲-G）三个必答题的判语

**判语一句**：`AC#1` **成立**（按已裁射程落干净），但"四枚入口件"这枚数**在盘上就是⛔ 全量**（真数＝五枚仓内 shell 入口），漏的那一枚**⛔ 影响本格成立、⛔ 归本腿落地**（归口见 §3 末）。

本件里**每一行读数都是我自己现跑的**（台件＝`gate-310v1.sh`／`ac3-gates-310v1.sh`／`block-bytes-310v1.sh`，读数件全在 `logs/`）。⛔ 抄 `310-r1` 的 logs、⛔ 抄编排者的转述。

---

## ① 这一步真的跑过吗（判据＝`A858` §6 新立那句："要么出现 `GOTMPDIR=` 那行、要么出现 `DEGRADED` 那行，两行都缺＝那一发⛔ 算这一步的凭据"）

**尺**＝对 `.scratch/wisp/probes/310/r1/logs/` 全部 66 枚件逐枚 `grep -c 'GOTMPDIR='` ＋ `grep -c 'DEGRADED'`（件＝`logs/10-r1-log-presence.txt`；"哪枚两行都缺"＝`logs/10-r1-entry-candidates.txt`，射程＝内容里出现三枚入口件名的 13 枚件）。

**逐字读数（带 `GOTMPDIR=` 那五枚，腿用来支撑"这一步跑过"的全部凭据）**＝`logs/10-gotmpdir-verbatim.txt`：

```
winsec-tests.sh: ticket 310 AC#1: GOTMPDIR=D:/tmp/wisp-test-gotmp/winsec-20261012-115046-1584 - this run's t.TempDir() family lands here; nothing deletes it.
winsec-tests.sh: ticket 310 AC#1: GOTMPDIR=C:/Users/swq/AppData/Local/wisp-test-gotmp/winsec-20261012-115521-883 - ...
wisp-cli-tests.sh: ticket 310 AC#1: GOTMPDIR=C:/Users/swq/AppData/Local/wisp-test-gotmp/cli-20261012-115545-1029 - ...
portable-tests.sh: ticket 310 AC#1: GOTMPDIR=D:/tmp/wisp-test-gotmp/portable-20261012-114637-408 - ...
portable-tests.sh: ticket 310 AC#1: GOTMPDIR=C:/Users/swq/AppData/Local/wisp-test-gotmp/portable-20261012-115401-2013 - ...
```

★**那五枚报出来的根名，逐枚在盘上存在**（我用 PS 的 `Get-ChildItem -Directory` 现量，件＝`logs/00-gate-start.txt` 第 71～84 行：`winsec-…115046-1584` entries=4、`winsec-…115521-883` entries=0、`cli-…115545-1029` entries=0、`portable-…114637-408` entries=0、`portable-…115401-2013` entries=1）⇒ **`GOTMPDIR=` 那行⛔ 是打字打出来的，它命名的东西在盘上**（这一条正是 10-12 那两枚腿栽过的"名字只认盘上来路"）。

**★两行都缺的那几枚，逐枚归位（这是本格的主活）**：

| 件（逐字名） | 两行都缺？ | 它是什么发 | 本格判它算什么 |
|---|---|---|---|
| `11-winsec-headtree.txt` | 缺 | HEAD 导出树（⛔ 本腿步子那台）的 winsec 档 | **⛔ 算这一步的凭据**（它算"改前"那一侧）✓ 腿也⛔ 拿它当"跑过" |
| `12-cli-headtree.txt` | 缺 | HEAD 导出树的 cli 档，**被自家 GUARD 打回**（末行逐字 `wisp-cli-tests.sh: GUARD - /d/tmp/310r1-tree2.dieduj/third_party/sherpa-onnx does not exist.`） | ⛔ 算凭据；⚠ 腿的交付件**没点名这枚件**（见下"缺口"） |
| `12-cli-noredirect-cf.txt` | 缺 | 反控：同台面把根**指回 OS 临时根**（`grep -i gotmpdir`＝0 命中＝那发本来⛔ 设值） | ⛔ 算"跑过"的凭据、**算反控**，腿正是这么用的 ✓ |
| `30-census-prepatch.txt`／`-totals.txt` | 缺 | HEAD 导出树的 census 跑前读数 | ⛔ 算凭据（改前基线）✓ |
| `30-census-postpatch-totals.txt`／`postpatch2-totals.txt`／`prepatch-totals.txt` | 缺 | **单行 totals 摘录件**（`wc -l`＝1），⛔ 一发射击 | ⛔ 适用这条尺（它⛔ 声称自己是一发） |
| `11-winsec-lappdata.txt` | 缺 | 1426 枚行的 winsec 档读数，**腿的 `.md` 五枚交付件里一枚都没点名它**（`grep -c -F` 命中＝0，件＝`logs/10-r1-uncited.txt`） | 判＝**交付面缺口**（见下），⛔ 判成"造假" |
| `12-cli-headtree2.txt` | （同上未点名族） | 未点名的第二发 cli 导出树 | 同上 |

⇒ **判语①**：拿这条尺**回头量 `310-r1` 每一件**的结果＝**没有一枚"声称这一步跑过"的发两行都缺**；那五枚带 `GOTMPDIR=` 的发⛔ 止带行、根名还在盘上。**两行都缺的 8 枚全部是改前发／反控发／单行摘录件，腿⛔ 拿它们当过"跑过"的凭据** ⇒ ① **成立**。
⚠ **一条⛔ 成立的边**：`DEGRADED` 那一支在腿的 66 枚件里**零命中**（件＝`logs/10-r1-log-presence.txt`，`degraded_hits=[1-9]` 的枚数＝**0**）⇒ 降级分支⛔ 被任何一发 exercised 过，腿自己在 `Q-310-r1-3` 认了这一点。**本腿替它补了一发**（⛔ 改 `scripts/**` 一个字节，只设环境变量）：
`WISP_TEST_GOTMP_ROOT="D:/work/workspace/projects plans/Wisp/go.mod/nope" bash scripts/portable-tests.sh --scope=census` ⇒ 逐字（件＝`logs/10-degraded-verbatim.txt`，那发 `probe_rc=0`、stdout `GOTMPDIR` 行＝0、stderr 293 字节）：
```
portable-tests.sh: ticket 310 AC#1: DEGRADED - cannot create D:/work/.../go.mod/nope/portable-20261012-125524-657, so GOTMPDIR stays unset and a t…
```
⇒ **降级句真会响、且⛔ 把退码翻红（rc 仍 0）**＝编排者 `A858` §6 那句"响在嘴上、⛔ 响在退码上"**被我独立证实**；census totals 那行在这一发里逐字仍是 `packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`。

## ② 跨平台那一支长什么样（判据＝现读三段块字节＋`go env GOOS` 一把尺）

尺＝`logs/10-block-bytes.txt`（awk 逐块数行/字节，`go env GOOS` 同行落）：

```
go_env_GOOS=windows
block=scripts/portable-tests.sh lines=68 bytes=4381 export_GOTMPDIR_lines=1 DEGRADED_lines=2 left_as_is_lines=1
block=scripts/winsec-tests.sh lines=21 bytes=1349 export_GOTMPDIR_lines=1 DEGRADED_lines=1 left_as_is_lines=0
block=scripts/wisp-cli-tests.sh lines=20 bytes=1298 export_GOTMPDIR_lines=1 DEGRADED_lines=1 left_as_is_lines=0
```

逐字读数（`scripts/portable-tests.sh:162-166`，我自己 `awk` 取的）：

```
    else
        echo "portable-tests.sh: ticket 310 AC#1: GOTMPDIR left as-is on GOOS=$(go env GOOS)" \
            "(the ticket's measurement surface is the Windows host's temp root)."
    fi
```

- **非 Windows 那支＝明确不动**：`portable-tests.sh:163` 打一行 `left as-is`，块内那一枚 `export GOTMPDIR` **只在 Windows 支里**（`left_as_is_lines=1` 那一支⛔ 设值）；下面两枚件的块里同样⛔ 有一枚 `export`，也都关在 `GOOS = windows` 那个合取条件里。✓
- `winsec-tests.sh:54`／`wisp-cli-tests.sh:67` 的块是**合取条件**（逐字：`if [ -z "${GOTMPDIR:-}" ] && [ "$(go env GOOS)" = windows ] && [ -n "${LOCALAPPDATA:-}" ]; then`）⇒ 非 Windows 时那两枚件**自己一行都⛔ 打**（`left_as_is_lines=0`）。这算⛔ 算"明确不动"？**算**，两条理由＋一条具名缺口：① 两枚件都**委派**给 `portable-tests.sh`（`winsec-tests.sh:115`／`wisp-cli-tests.sh:133`，腿件 `10-impl.md` §1 的行号），而被委派那枚件会打 `left as-is`；② `wisp-cli-tests.sh:80-87` 的 `if [ "$(go env GOOS)" != windows ]` GUARD **改前的 blob 里就在**（我的尺＝`git show a21dcd5b~1:scripts/wisp-cli-tests.sh | grep -c 'this is the windows leg of the cmd/wisp gate'` ＝ **1**）⇒ 本腿⛔ 动过那一支，它在非 Windows 上本来就是 `exit 2`；③ ⚠ 缺口：`winsec-tests.sh` 在 `LOCALAPPDATA` 为空的 Windows 上是**静默**（⛔ `DEGRADED`、⛔ `left as-is`），要靠委派的 `portable-tests.sh:158` 才响——**记一笔、⛔ 翻格**（`Q-310-v1-5`）。
- **ubuntu 那一支今天仍会漏＝盘上仍然可核**：尺＝`grep -n 'ubuntu leg is left exactly as it is' scripts/portable-tests.sh` ⇒ 逐字命中 `:138`：`# are about this host's %TEMP%. The ubuntu leg is left exactly as it is because`（后两行写明"这一腿⛔ 能跑 CI 去量重定向的临时根会把这些步子改成什么，已登记成欠账、⛔ 默默当它安全"）⇄ 腿件 `90-final.md` 的 `Q-310-r1-3` 第 ① 条具名"那一族残渣在 `/tmp` 继续长、无人管"。⇒ **⛔ 被读成"两边都管住了"** ✓
- **改动面我自己复跑**（⛔ 信转述）：`git log --oneline a21dcd5b~1..3122306e`＝**10 笔**；`git diff --numstat a21dcd5b~1 3122306e -- scripts` 逐字＝`69 0 scripts/portable-tests.sh`／`21 0 scripts/winsec-tests.sh`／`20 0 scripts/wisp-cli-tests.sh`（件＝`logs/10-numstat-scripts.txt`）；`git diff --name-only a21dcd5b~1 3122306e -- '*.go'` 枚数＝**0**（件＝`logs/10-r1-commits.txt` 同批）。⇒ 与票面/派单转述的 69/21/20 加、0 删、`*.go` 零枚**逐字一致**。

## ③ ★名册是不是全量（两把尺都我自己现跑，⛔ 抄腿的话）

**尺 A**（具名射程＝`scripts/` ＋ `tools/`，件＝`logs/10-roster-r1.txt`／stderr＝`logs/10-roster-r1.err` 0 字节）：
```
grep -rn --include=*.sh --include=*.ps1 'go test' scripts/ tools/   →  rc=0  hits=30  files=7  err_bytes=0
```
7 枚命中文件逐枚归位（"真起点"＝那一行是命令而⛔ 是注释）：

| 文件 | 命中行 | 是不是真起点 |
|---|---|---|
| `scripts/portable-tests.sh` | 3 行（`:686`/`:721` 的 `go test -list`＋`:764` 委派 `runtests.sh`） | **是**（腿落了步子） |
| `scripts/winsec-tests.sh`／`scripts/wisp-cli-tests.sh` | 各 1 行委派 | **间接是**（腿落了步子） |
| `tools/d22scan/runtests.sh` | `:75` 逐字 `go test -v -count=1 "$@" >"$out" 2>&1` | **★是＝真·第五枚** |
| `scripts/d22scan.sh` | 3 行**全是注释**（`:19`/`:22`/`:35`，逐字如 `# Why step 1 goes through tools/d22scan/runtests.sh and not bare `go test``） | 它⛔ 直接起、它**调用**上面那枚 runner ⇒ 入口经它 |
| `scripts/slo-check.ps1` | 1 行＝`:336` 注释（逐字 `# agents are mid `go test`. The resulting D32 reading can then be falsely red`） | **⛔ 是** ⇒ 腿"第四枚名册件零枚调用点"**被我独立复认** |
| `scripts/portable-tests-selftest.sh` | 1 行＝`:8` 注释（`# "Pinned: 25, resolved: 36" and killed the core step BEFORE `go test` ran - while`） | **⛔ 是**（★这把尺我自己也差点把它当第六枚，逐字读完才归位） |

**尺 B**（把 CI 与产码面并进来，只为把"全量"这枚数说清，件＝`logs/10-roster-scope-extra.txt`／stderr 0 字节）：
```
grep -rn --include=*.yml --include=*.yaml --include=*.ps1 --include=*.sh --include=*.bat 'go test' .github cmd internal frontend
  →  rc=0  hits=10  files=1（全部在 .github/workflows/ci.yml）
```
其中 `ci.yml:82` 逐字＝`run: bash tools/d22scan/runtests.sh -C tools/d22scan ./...` ⇒ **CI 直接点那枚 runner**（连 `scripts/d22scan.sh` 那层都跳过）；其余 9 行是注释。CI 那一层＝票 310 裁形**显式排除**的射程（`⛔ 动 CI`），属于"被排除"⛔ 属于"被漏掉"。

**★本格要的那一句裁决**：**"四枚"⛔ 是全量**（盘上真数＝**五枚**仓内 shell 入口：`portable-tests.sh`／`winsec-tests.sh`／`wisp-cli-tests.sh`／＋经 `d22scan.sh` 的 `tools/d22scan/runtests.sh:75`；第四枚名册件 `slo-check.ps1` 零调用点）。
**⇒ `AC#1` 仍判"成立"**，理由是这一格裁的是"按已裁射程落干净"：裁形按**路径**点名四枚，那三枚真起点全落、那一枚（`slo-check.ps1`）经我两把尺复认**本就没有可落的调用点**，⛔ 一枚漏落、⛔ 一枚多落、⛔ 一处越界（越界尺见 `30-ac3-gates.md`）。漏的那一枚**在名册之外**＝**名册错**（编排者 `A858` §6 已自认"我的名册漏"），⛔ 是落地腿跑歪；把它读成 `AC#1` 不成立＝惩罚腿⛔ 惩罚名册。
⛔ 本腿动第五枚（落地活归下一程）。★**它的漏⛔ 是纸面推测，本腿现量到了**：我跑门禁那一发 `sh scripts/d22scan.sh`（rc=0）的读数件里 `GOTMPDIR` 命中＝**0**（件＝`logs/30-ac3-run.txt` 的 `d22scan_rc=0 ... gotmpdir_hits=0` 行）⇒ **走那一条路死掉的测试二进制今天照旧在 `%TEMP%` 顶层留目录**。要不要把 `tools/d22scan/runtests.sh`（或 `scripts/d22scan.sh`）纳入射程＝编排者一句话的事，`Q-310-v1-4` 已挂。
