# `10-impl.md` — 票 310 `AC#1` 落地形＝甲-G（腿 `310-r1`，实现者）

锚＝`00-anchor.md`（起手 HEAD `7e6c1579`，分支 `dev`）。本件只写"改了哪几枚、调用点行号、专用根怎么命名"；判语归下一程非实现者腿。

## 1. 改了哪几枚 ＋ `go test` 调用点逐枚点名（本腿自己现跑，⛔ 照普查腿的名册抄）

尺＝`grep -n -E 'go test |bash "$portable"|sh "$strict"' scripts/*.sh scripts/*.ps1` ＋人工剔除注释行；下表行号＝**改后工作树**的现行号（改前＝HEAD `7e6c1579` 的号）。

| 枚 | `go test` 调用点（本腿现跑点名） | 本腿加的那一步 |
|---|---|---|
| `scripts/portable-tests.sh` | **真起点三处**：`:686`（`go test -list '.*'`，改前 `:617`）、`:721`（`go test -list "^<名>$" <包>`，改前 `:652`）、`:764`（`sh "$strict" "${scope[@]}" -count=1 -skip ...` → `tools/d22scan/runtests.sh:75` 那枚真 `go test -v -count=1`，改前 `:695`） | 新增块起于 **`:99`**（`export GOTMPDIR` 在 `:147`，两枚降级句在 `:152`/`:158`），位置＝`strict` 守卫之后、任何 `go` 调用之前；增量 **＋69 行／删 0** |
| `scripts/winsec-tests.sh` | **零枚直接 `go test`**；唯一出向＝`:115`（改前 `:94`）`bash "$portable" "${scope[@]}" 2>&1 \| tee "$capture"` | 新增块起于 **`:46`**（`export GOTMPDIR` 在 `:57`，降级句 `:61`）＝委派之前；增量 **＋21 行／删 0** |
| `scripts/wisp-cli-tests.sh` | **零枚直接 `go test`**；唯一出向＝`:133`（改前 `:113`）`bash "$portable" --scope=cli` | 新增块起于 **`:60`**（`export GOTMPDIR` 在 `:70`，降级句 `:74`）＝PATH 守卫之前；增量 **＋20 行／删 0** |
| `scripts/slo-check.ps1` | **零枚 `go test` 调用点**：`grep -c -i 'go test'`＝**1**，那一枚命中是 `:336` 的注释（"agents are mid `go test`"）；它的 go 侧只经 `:296-314` 起 `scripts/build.ps1`（`go build`），`:361` 里的 `go.exe` 是争用侦测的进程名清单、⛔ 调用 | ★**本腿没在这一枚里造任何一步**，具名顶回派单（`00-anchor.md` §5 第 1 条）。⛔ 在只有 `go build` 的入口里塞"跑测临时根"＝无的放矢的语义改动 |

三枚 delegating 入口的块与 `portable-tests.sh` 那块**同形且幂等**：`[ -z "${GOTMPDIR:-}" ]` ⇒ 一次运行只有一枚根、由最外面那枚入口命名，⛔ 第二枚根、⛔ 覆盖调用方已有的值。三处逻辑相同而**⛔ 抽成一枚共享件**＝派单写死"唯一落点＝那四枚"，新建第五枚 `scripts/**`  helper 属超授权扩面（记在这里给裁决腿审这一处重复的代价）。

★**顶回第二条（不在名册上的第五枚）**：`scripts/d22scan.sh` 经 `tools/d22scan/runtests.sh:75` 也起 `go test`（`tools/**` 是本票⛔ 落点，`scripts/d22scan.sh` 又⛔ 在那四枚名册上）。⇒ 本腿⛔ 动它，把它漏下的那面具名登记成欠账 `Q-310-r1-2`。

## 2. 专用根怎么命名（★这一枚是量出来的，⛔ 拍脑袋）

```
默认根＝${WISP_TEST_GOTMP_ROOT:-"${LOCALAPPDATA//\\//}/wisp-test-gotmp"}/<入口>-<yyyymmdd-HHMMSS>-<pid>
本机现量形＝C:/Users/swq/AppData/Local/wisp-test-gotmp/portable-20261012-115401-2013
```

四条命名约束与它们的**实测出处**（读数全在 `20-paired-readings.md`）：

1. **⛔ 建在 `%TEMP%` 顶层**（普查腿具名报过的假绿坑：非递归名册尺会读出"0 新增"、而残渣只是出了尺的射程）。⇒ 根是 `%TEMP%` 的**兄弟**（`%LOCALAPPDATA%\wisp-test-gotmp`），⛔ 其子目录；起手 `mkdir -p` 建，收尾只把路径打进读数件。
2. **⛔ 换卷**（派单⛔ 写这一条，是本腿撞出来的）。第一稿把根建在**仓所在卷**（`D:/tmp/wisp-test-gotmp/...`），当场把 `--scope=winsec` 那档从 **58 PASS / 0 FAIL** 打成 **56 PASS / 2 FAIL**：`TestGateJudgesThePrivateSetByResolvedSID`（两枚子例）＋`TestSealNarrowsAndNamesThePrincipalItRemovedBySID`——它们判的是"封进去的那套 ACL 继承自哪儿"，换一枚卷的继承集合⛔ 同一套答案。同一档在 `%LOCALAPPDATA%` 根上复跑＝**58 PASS / 0 FAIL**，与⛔ 重定向（HEAD 导出树，件 `logs/11-winsec-headtree.txt`）逐字同数。⇒ **根落哪一枚卷⛔ 是自由变量**；本腿取"与 `%TEMP%` 同卷同 ACL 族"那一枚，因为它是唯一一把⛔ 改颜色的位置。这条同时把派单那句"甲-G ⛔ 换掉那 17 枚逐字读 `os.TempDir()` 的测试件的答案"**收紧成一句实话**：那 17 枚的**词面答案**确实没换（Windows 上 `os.TempDir()` 走 `GetTempPath`、⛔ 读 `GOTMPDIR`），但**物理位置一换**就够让三枚封 ACL 的用例改颜色 ⇒ "⛔ 打红既有钉"这一条⛔ 靠机制推理成立，只能靠复量。
3. **起手建、收尾只写路径**：⛔ 删整根（票面那条铁规矩＝先逐枚名册、先经非实现者复核）。根路径由每枚入口件在自己的输出里具名打印（`ticket 310 AC#1: GOTMPDIR=...`），⛔ 静默。
4. **建不出来＝响在嘴上、⛔ 响在退码上**：`mkdir` 失败或 `LOCALAPPDATA` 为空 ⇒ 往 stderr 打一行具名 `DEGRADED`（"GOTMPDIR 保持不设，残渣还会落回 OS 临时根"），本步⛔ 因此把闸门翻红。★这是**本腿自加的判断，写给裁决腿审**：⛔ 我能测 CI 的 windows runner 上 `%TEMP%` 的兄弟目录可不可写，硬失败会把一枚我⛔ 看得见的闸门改成红；换成响亮的降级句，宁可让那发读数留在 stderr。若裁决者认为这里该 fail-closed，这是第一个要改的点。

## 3. 语义面（⛔ 动了什么）

- 三枚 `.sh` 各**只加**一段独立块（`git diff --numstat -- scripts`＝`69 0`／`21 0`／`20 0`）；`git diff --name-only HEAD -- '*.go'`＝**0 枚**。⛔ 产码、⛔ `*_test.go`、⛔ CI／`.github/**`、⛔ 票 304 的 Job Object 面（`internal/proc/jobscope_*`、`internal/llm/adaptertest/childjob_*`、两枚 mockllm spawn 点）、⛔ golden／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt`／D43 转移表／C1–C32／D1–D47、⛔ 任何断言、⛔ 三枚冻结件。
- 机制依据（本腿现读 GOROOT，⛔ 引用别人的一句话）：`D:/work/base/go/src/testing/testing.go:1613` ＝ `c.tempDir, c.tempDirErr = os.MkdirTemp(os.Getenv("GOTMPDIR"), pattern)` ⇒ `GOTMPDIR` **只**搬 `t.TempDir()` 一族；`D:/work/base/go/src/cmd/go/internal/work/action.go:298` ＝ 同一枚变量还被 go 工具自己拿去放 `go-build*` 工作目录（本票⛔ 射程，只是同根顺带受益）。
- 平台门＝`go env GOOS` 为 `windows` 才落这一步；非 Windows 打一行 `left as-is` 且⛔ 设值（ubuntu 那一腿一个字节都不改，理由与欠账见 `90-final.md` 的 `Q-310-r1-3`）。
- ⛔ 任何一处删目录：本腿建的专用根（`C:/Users/swq/AppData/Local/wisp-test-gotmp/*` 与实验期那两枚 `D:/tmp/wisp-test-gotmp/*`）与落在里面的残渣**全部留在盘上**，逐枚名册见 `20-paired-readings.md` §4。
