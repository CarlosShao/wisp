# 72-r1 读数／门禁／反形（两台面 × 改前改后 × 两发 ＋ 三发变异 ＋ 三把门禁）

时刻＝`2026-10-11 11:33`（改前）／`11:43`（改后）／`11:46`（门禁＋变异）+0800。锚＝`9c83598d`。
本腿 commit＝`f9b38d5f`（锚件）／`f278c57f`（甲组 7 枚）／`79e4974d`（乙组 2 枚）。

**三态（每件读数都标）**
- **工作树**＝母仓 `D:\work\workspace\projects plans\Wisp`（只用于 `gofmt -l`／`git status`／门禁三把；⛔ 在母仓跑那 8 发颜色）。
- **仓外导出树**＝颜色读数的射程：`/c/Users/swq/tmp/w72r1/treeB`（改前，`git archive 9c83598d go.mod go.sum internal cmd`）与 `/c/Users/swq/tmp/w72r1/treeA`（改后，`git archive HEAD …`）。建树尺＝`archive_rc=0`／`tree_top=cmd go.mod go.sum internal`／`tree_go_files=648`（两棵都对得上，⛔ 少 `go.mod go.sum` 那一形）。⛔ 母仓⇄clone 那根轴、⛔ 在仓内建 worktree。
- **HEAD blob**＝格式名册第二档射程（`git archive HEAD internal/risk | tar -x` 后 `gofmt -l`）。
- 尺脚本＝仓内 `.scratch/wisp/probes/72/r1/run-two-surface.sh`（逐字可重跑）。原始件＝`.scratch/wisp/probes/72/r1/logs/tmp-{before,after}-{alias,machine}-{7,16}.txt`（工作树，已入库）＋同目录 `before-summary.txt`/`after-summary.txt`/`30-gates-mutations.txt`/`gate-{build,vet,d22scan}.txt`。三发变异的**原始** go test 输出在**仓外** `treeA/tmp-mut{A,B,C}.txt`、`treeA/tmp-restored.txt`（未入库，承重读数已逐行抄进本件 §3）。

---

## 1. 两台面 × 改前／改后 × 两发（尺＝`grep -c`，四数；含子测试行形）

命令一律 `go test ./internal/risk/ -count=1 -timeout 600s -v -run '<名册>'`，输出先落文件再取数。
名册甲＝7 枚指名（`TestClassifyAnchorSpellingIsNotVerdict|TestCanonicalInputGainsNoSecondForm|TestAListWinsWhereBothTablesHit|TestPathResolverShortNameAListDenied|TestPathResolverUNCAListDenied|TestPathResolverExtendedLengthPrefixAListDenied|TestBListDefaultDenyAndOverride`）；名册乙＝两文件全量 16 枚（`-run 'TestPathResolver|TestReparse|TestAList|TestBList|TestResolve|TestClassify|TestCanonical|TestOverride'`；**射程注**：`-run` 作用于整枚 `./internal/risk/` 包，该正则实际命中 **20 枚**顶层用例＝两文件 16 枚⛔ 全在里面＋包内另 4 枚同前缀（`TestOverrideKeyKeepsPosixBackslashesDistinct`／`TestResolveCanonicalIsPlatformShaped`／`TestResolveKeepsSeparatorsOutOfPosixNames`／`TestResolvePerCallBudget`），改前改后同一把尺⛔ 变，故作差仍然成对）。

| 台面 | 发次 | 改前（`9c83598d`） | 改后（HEAD＝`79e4974d`） |
|---|---|---|---|
| ①别名 `TEMP/TMP/TMPDIR=C:\Users\swq\tmp\RISKAT~1\RUNNER~1` | 甲（7 枚） | `rc=1` RUN=7 PASS=0 **FAIL=7** SKIP=0 | `rc=0` RUN=7 **PASS=7 FAIL=0** SKIP=0 |
| ①别名 | 乙（全量） | `rc=1` RUN=20 PASS=13 **FAIL=7** SKIP=0 | `rc=0` RUN=20 **PASS=20 FAIL=0** SKIP=0 |
| ②本机 `TEMP/TMP=C:\Users\swq\AppData\Local\Temp` | 甲（7 枚） | `rc=0` RUN=7 PASS=7 FAIL=0 SKIP=0 | `rc=0` RUN=7 PASS=7 FAIL=0 SKIP=0（**颜色⛔ 变**） |
| ②本机 | 乙（全量） | `rc=0` RUN=20 PASS=20 FAIL=0 SKIP=0 | `rc=0` RUN=20 PASS=20 FAIL=0 SKIP=0 |

- **判据甲**＝别名 7 枚 FAIL→PASS ✓；本机两发颜色⛔ 变 ✓（同一份码只差那三枚 env）。
- **判据乙**＝红名册作差尺（`comm`，改前⇄改后**同一枚台面**各取一次）：
  `alias-7: before=7 after=0 only_in_before=7 only_in_after=0`／`alias-16: before=7 after=0 only_in_before=7 only_in_after=0`／`machine-7: 0/0/0/0`／`machine-16: 0/0/0/0` ⇒ **⛔ 一枚新红**，16 枚全量⛔ 出现"从没有→有"。
- **台面在场尺**（每发之前）：①`--- SKIP`＝0（别名 8 数里全 0；`junction:126` 那枚 `getShortPath` 空返回的票 18 既有 Skip 从未触发＝短名真在）；②`t.TempDir()` 返回值含 `~1`＝仓外独立探针直接问同一支 API：`os.MkdirTemp("", "W72R1PROBE")` 在别名台面＝`C:\Users\swq\tmp\RISKAT~1\RUNNER~1\W72R1PROBE638957454 hasTilde=true`、本机台面＝`C:\Users\swq\AppData\Local\Temp\W72R1PROBE3605685462 hasTilde=false`（件＝`after-summary.txt`/`before-summary.txt` 的 `tempdir_probe` 两行）。
- 附带（本腿⛔ 需要但顺手取了）：整包 `go test ./internal/risk/`（别名台面、工作树）**`ok 3.956s`**＝⛔ 只有那 16 枚，包内其余用例也⛔ 被我的改动打红。

## 2. 门禁三把（各自一行 `rc=`，件⛔ 是 0 字节的那两把由退码承担）

```
build_rc=0 outbytes=0      # GOFLAGS= go build ./...      （干净时本就零输出）
vet_rc=0   outbytes=0      # go vet ./internal/risk/
d22scan_rc=0 outbytes=23088  # sh scripts/d22scan.sh       （件已入库：logs/gate-d22scan.txt）
```
- 格式名册两把并排（⛔ 把 CRLF 假枚读成"有人改了格式"）：
  工作树射程 `gofmt -l internal/risk` ＝ **1 枚**＝`internal\risk\provenance.go`；HEAD blob 射程 ＝ **0 枚**。⇒ 那一枚是本仓已知 CRLF 假枚形状（尺＝`logs/eol-and-gofmt-facet.txt`：`git ls-files --eol` 逐枚＝`provenance.go` 是 `i/lf w/crlf`（工作树 CRLF、索引 LF），而我三枚写面文件全是 `i/lf w/lf`）；`gofmt -l` 点名这四枚的射程里**只命中 `provenance.go`** ⇒ **本腿三枚写面文件在两档射程里都⛔ 出现**。

## 3. 反形三发（全在仓外树 `treeA`；每发之后 `git archive HEAD <file> | tar -x` 逐字节还原并落 sha）

pristine sha（前 16 位）：`syncdirs.go=aa4779bb52704f43`、`blacklist.go=b8f316794704b04b`、`pathresolver_junction_windows_test.go=ab07caa3f3083f40`。

| 发 | 变异（`perl -0pi -e`，落地尺＝命中数） | 台面／名册 | 读数 | 还原 sha |
|---|---|---|---|---|
| ⓐ 甲组某枚新期望侧**退回旧形** | `if !namesSameFile(res.Canonical, long) {` → `if !strings.EqualFold(res.Canonical, long) {`（3 处；`landed=3`，剩 1 处＝#1 的前提句⛔ 在射程内） | 别名／7 枚指名 | `rc=1 FAIL=3 SKIP=0`＝**恰是那三枚**：`TestPathResolverShortNameAListDenied`／`TestPathResolverUNCAListDenied`／`TestPathResolverExtendedLengthPrefixAListDenied` ⇒ 甲族 #4/#5/#6 的新期望侧是承重的 | `restored_sha_junction=ab07caa3f3083f40` ✓ |
| ⓑ 乙组负向钉所要防的那支：**恒真形** | `syncdirs.go:442` `if !s.complete && s.home != "" && aUnder(candForms, anchorForms(s.home)) {` → `if true {`（`landed=1`） | 别名／`-run 'TestSync'` | `rc=1 FAIL=13 SKIP=1 PASS=16`。SKIP 那枚＝`TestSyncRegistryProbeLive`（票面 `AC#4` 早具名的既有 live-probe，⛔ 本腿引入）。**响的文案只有本腿新增那两句**：`syncdirs_test.go:230 … must NOT be sync`（四档 src 各一次）与 `:254 a write outside the profile and outside the unverified root must not be flagged` ⇒ 三枚**正向**断言（`must stay armed`／`writes under the injected root must be sync`）在恒真形下**全绿**＝旧仪器确实哑，新负向钉是唯一的牙。同发另 11 枚既有 sync 用例转红（恒真形把"界内有界外"整族抹平）＝预期形状，⛔ 计入本腿产物 | `restored_sha_syncdirs=aa4779bb52704f43` ✓ |
| ⓒ 拔掉产码那枚守卫分岔（对照，⛔ 改母仓文件） | `blacklist.go:122` `if cand.raw != cand.real {` → `if false {`（＝票面 `:95` 逐字作废过的那条放宽，只作对照） | 别名／`-run 'TestOverride\|TestAList\|TestBList\|TestClassify'` | `rc=1 FAIL=1`＝**恰 `TestOverrideOnlyAcceptsTheResolvedForm`**（窄规则正钉），而本腿改过的 `TestAListWinsWhereBothTablesHit`／`TestBListDefaultDenyAndOverride` 在该放宽下**仍绿** ⇒ 证明甲族 #3/#7 的修法是"把键换成已解析形"、⛔ 是靠放宽换绿（若我把 `overrideApplies` 改了，这两枚会⛔ 响，正钉会响） | `restored_sha_blacklist=b8f316794704b04b` ✓ |
| 正控 | 全部还原后同一棵 `treeA` | 别名／7 枚＋`TestSync` | `restored_rc=0` ✓（⛔ 还原＝空操作假绿） | — |

## 4. ⛔ 做到的事（具名，⛔"应该没问题"含过去时）

1. **CI／runner 台面⛔ 跑**：本腿只在那两枚**本机可造的台面**（别名／本机 `%LOCALAPPDATA%\Temp`）上取数。票面 `AC#4` 那句"runner 侧必须来自一次真 run"仍⛔ 闭（要 push 才有一次真 run，本票禁 push）⇒ 那一格的颜色⛔ 是我交的，我⛔ 判。
2. **POSIX 档零求值**：两枚写面文件都带 `//go:build windows`，`syncdirs_test.go` 无 tag 但我的两句负向钉走的判定与平台无关；**darwin/linux 档我⛔ 跑**（本机⛔ 交叉测试台、且票 307 §3 已把"新夹具在 POSIX 档⛔ 求值"记成既有欠读数）。⇒ 我⛔ 能交"⛔ 新红在 POSIX 上"这一色。
3. **#2 的"恰 1 形"那半句承诺现在零仪器**（选型与代价逐字在 `10-class-a-expected-side.md` §#2）：我⛔ 能补，补它要动 `Resolve` 的缺失叶子分支＝产码、⛔ 本票射程。
4. **三发改异的原始 go test 输出躺在仓外树** `treeA/tmp-mut{A,B,C}.txt`（临时件只建不删，⛔ 入库；承重读数已逐行抄进 §3）。下一枚验收腿要重跑变异，尺＝本件 §3 那张表的三行 `perl` 表达式 ＋ `run-two-surface.sh`。
5. **`gofmt -l internal/risk` 工作树那 1 枚（`provenance.go`）⛔ 是我改的**（我⛔ 碰过它），是已知 CRLF 形状，我⛔ 顺手"改绿"格式名册。
6. **`-run 'TestSync'` 那一发的 `SKIP=1`**＝既有 live-probe，⛔ 是本腿引入也⛔ 是本腿能消。
7. 共享索引里躺着**⛔ 我名的 9 枚 tracked 改动**（票 309、`probes/orch/20261011/*`、`docs/reports/HANDOVER.md`、`docs/reports/pending-and-issues.md`）：我的三笔都用显式 pathspec，`git show --numstat --format=` 逐笔核对＝**越出授权面 0 行**（尺读数＝每笔 `non_test_internal_files=0`、`forbidden_files=0`）；我⛔ `reset`／⛔ 取消别人的暂存，那些文件我⛔ 读过内容。
