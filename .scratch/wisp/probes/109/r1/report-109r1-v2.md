# 票 109 返工腿 109-r1（v2 重放）交件读数汇总

工作树＝`D:\work\workspace\projects plans\Wisp`，分支 `dev`。起手 HEAD＝`fd091cd1`（交件侧 HEAD 随在飞腿漂，
⛔ 当通则）。名册尺三把同号＝**6 笔**（`git log --grep='^票 109 返工腿 109-r1'` ⇄ `%s` 内存 `index($2)==1` ⇄
主题 grep），其中 1 笔是编排者的派单存档 `11804945`，本腿自己落 **5 笔**：

| 笔 | 内容 | 写面（逐枚 `git show --name-only`；四枚超限读数件在本腿最后一笔按具名 `-partN` 重落，故下表写的是**最终入库名**） |
|---|---|---|
| `0242dc07` | 〔甲〕AC#2 形① 产码＋测试 | `internal/models/downloader.go`、`internal/models/handoff_guard_109r1_test.go`、`probes/109/r1/{baseline-models-head-part1,cellA-pre-red,cellA-landing-green}.log` |
| `3951247f` | 〔乙〕AC#4 两洞用例 | `internal/models/handoff_ac4_109r1_test.go`、`internal/models/handoff_ac4_principal_109r1_windows_test.go`、`internal/models/handoff_window_109_windows_test.go`、`probes/109/r1/{cellB-pre-red,cellB-landing-green,landing-models-count1-part1,landing-models-count2-part1}.log` |
| `c2722dfb` | 变异两发＋POSIX＋台面复量 | `probes/109/r1/mutation-M-{A,B}{,-transcript}.log`、`posix-docker-golang1.27.log`、`testbed-*.log` |
| `a2234c00` | 票面登记（只追加） | `.scratch/wisp/issues/109-…-after-ensure.md`（numstat `73  0`，六格框 `- [ ]` 前后同为 6、`- [x]`=0） |
| 本笔 | ★读数件收口 ≤160 行帽（超限者按具名 `-partN` 拆分、拼接后与仓外完整日志逐字节相同）＋默认台根成对复量＋汇总件 | 只含 `.scratch/wisp/probes/109/r1/**`（四枚超限件按删除＋`-partN` 形入库＝⛔ 截断证据、⛔ `--amend`） |

越界尺＝逐枚 name-only 全集里 `internal/winsec`⇄`internal/risk`⇄`tools/d22scan`⇄`allowlist`⇄`go.mod`⇄`go.sum`⇄
`ci.yml`⇄`docs/PLAN`⇄`docs/specs`⇄`frontend/`⇄`design/`⇄`cmd/wisp` 命中 **0 枚**（`OVERLAP=NONE`）；
`git status --porcelain internal/models/` 现＝空（本腿写面全部入库）。

## 〔甲〕AC#2 ＝ 形①（守卫挪到 `Ensure` 的返回路径），⛔ 形②

裁决表 `docs/evidence/s1/109-adversarial-acceptance.md:43-44` 给的两条路。选①的凭据是**真在跑的装配链**，
不是口头承诺：

- `cmd/wisp/models.go:305` `bridge := models.WireDownloading(ms.mgr, machine)` →
  `internal/models/bridge.go:39 Run` → `:44 b.mgr.Ensure` → `:58 VerifyInstalled`；
  `cmd/wisp/models.go:249` 直调 `store.mgr.VerifyInstalled(id)`（`wisp models verify`）。
- 可达性仪器 `internal/models/assembly_reachability_121_test.go:126`
  `TestAC3CapabilityPackagesReachTheCompositionRoot`：`--- PASS`／GOOS=windows⇄linux⇄darwin 三支全 PASS，
  逐字读数在 `probes/109/r1/baseline-models-head-part1.log`（HEAD 树，⛔ 我的字节）与
  `shared-worktree-models-count2-part1.log`（含我的字节）两处都真跑到。
- 形②（"预置守卫待接线"＋"没人调"红钉）今天会**恒绿**＝它要挡的第 8 类缺陷已被票 121 拆掉，
  且形②治不了裁决表 `:297-301` 逐字点名的另一半（守卫"是调用点自带的，不是 choke point 强制的"）。

**修前红逐字**（HEAD 产码不动、只加测试件＝`probes/109/r1/cellA-pre-red.log`，另在默认台根下复量一次＝
`cellAB-pair-pre-default-temproot.log`）：`--- FAIL: TestAC21FreshInstallIsGuardedOnEnsuresOwnReturnPath`
＋ `--- FAIL: TestAC41NailByteFace…` 等 5 枚红；AC#2 那枚的断言原文：

    AC#2/R-109-4 (shape 1): Ensure returned a model directory as available although 1 byte of model.onnx was
    swapped after the last check that ran on its return path. dir=… swaps=1 - a caller that uses Manager.Ensure
    directly, without DownloadingBridge, gets the hand-off with no guard on it at all, which is exactly what
    R-109-4 says the bridge-only placement left open.

换字节＝构造在 `Ensure` 自己的 `Progress(PhaseDone)` 回调里＝该分支最后一次落点侧检查之后、返回之前，
⛔ 假设。修后绿三枚全 PASS，拒绝句带 `sha256 mismatch (pins 785b0751fc2c..., got dd7bafd7c0c4...)`。
**反半边**：`server hits: install=1, second-instance reuse=0` 修前⇄修后同读数；`dir1==dir2`；
`TestAC2NormalInstallAndSecondInstanceReuseAreNotBroken`⇄`TestAC3WideTemp…` 照旧绿。

## 〔乙〕AC#4 ＝ 补 `R-109-2` ⇄ `R-109-3` 两条界线的洞

`R-109-2`（未点名文件⇄extra 目录进哈希面）＝`handoff_ac4_109r1_test.go` 四枚，红→绿见
`cellB-pre-red.log` ⇄ `cellB-landing-green.log`。字节面装进那枚钉子**自己走的动词**（`Ensure` 返回路径）。

`R-109-3`（⛔ 按账户名⇄⛔ 按单一主体过滤 ACL，让显式外来 ACE 逃过去）＝新增单一判据 `ac4BreakLineTwo`
（对象上**所有**非 `(I)` 的 ACE），并把**已入库那枚钉子** `TestAC4InstallDirWidthIsDeliberateAndHasABreakLine`
改为调用同一判据 ⇒ 仓内不留两套口径。修前红逐字（`cellB-pre-red.log`）：

    AC#4/R-109-3: break line 2 reports nothing while this object carries an explicit non-inherited write grant
    for Everyone (Everyone:(OI)(CI)(M)). … the object has 1 ACE(s) that did not arrive by inheritance: []

播种＝`icacls /grant:r *S-1-1-0:(OI)(CI)(M)`，落没落由 `allTrusteesOn` **独立于判据本身**复核（防止"⛔ 落地"
冒充"界有洞"）。反半边＝`TestAC4R4NailBreakLineTwoStaysQuietOnInheritedWidth`（实测新装目录一条非 `(I)` 都没有
⇒ "对象上有非继承 ACE"⇄"有人手工加宽这一枚对象"是同一件事、票 95 的继承宽仍读"允许"）。
口径边界照裁决表 `:321` 逐字执行＝判据落 `internal/models/` 测试、**`internal/winsec/**` 零字节**。

注记一格诚实话：`R-109-3` 的修复本体**就是测试侧判据**，所以它的"修前红"两形＝① 该判据的 pre-fix 形落盘后跑红
（`cellB-pre-red.log`，逐字如上）＋② 变异 M-B 把它换回旧口径后跑红（`mutation-M-B-transcript.log`）；
⛔ 像 〔甲〕那样存在一枚"HEAD 产码不动就红"的独立形态——HEAD 上根本没有这个判据。

## 变异两发（落点⛔ 同、每发落刀前后各验锚命中数、当场还原报 pristine 哈希）

| 发 | 锚点 | 前 | 后 | build/vet | 红 | 绿（⇒ 没顺手砸别人） | 还原哈希 |
|---|---|---|---|---|---|---|---|
| M-A | `downloader.go` `Ensure` 三条返回路径的 `m.handOffCheck(...)` 调用行（`:195/:205/:241`），`VerifyInstalled` 那枚不动 | 1/2/1 | 0/0＋死守卫行 3 | `go build` rc=0 | `TestAC21`、`TestAC22`、`TestAC41/42/43Nail` | `TestAC1HandoffRefuses…`（bridge 守卫仍在）、`TestAC2NormalInstall…`、`TestAC3WideTemp…`、`TestAC23`、`TestAC44` | `079430b4d12b…8a3738a8c` == pristine |
| M-B | `handoff_ac4_principal_109r1_windows_test.go` `ac4BreakLineTwo` 的 `range allTrusteesOn(t, path)` 换回 `range aceLinesOwnedBy(t, path, 'BUILTIN\Users')` | 1/0 | 0/1 | `go vet` rc=0 | `TestAC4R3…` | `TestAC4R4…`、`TestAC4InstallDirWidth…`、`TestAC1TheWriteWindow…` | `d3de244d1cbb…4c9b7e56` == pristine |

## POSIX 那半真跑（容器逐字读数）

`golang:1.27`＋`CGO_ENABLED=0`，仓库字节走 tar on stdin，**容器内先 `ls -l` 证明看得见文件**
（⛔ Git Bash `-v "C:\…"` 静默挂空＝假绿；本机第一发踩到，改用具名卷＋清内容后重跑才有下文）：

    -rw-r--r-- 1 197609 197609 1026 Oct 20 13:17 /w/go.mod
    -rw-r--r-- 1 197609 197609 29447 Oct 20 13:19 internal/models/downloader.go
    -rw-r--r-- 1 197609 197609  9096 Oct 20 13:19 internal/models/handoff_ac4_109r1_test.go
    -rw-r--r-- 1 197609 197609 11082 Oct 20 13:19 internal/models/handoff_guard_109r1_test.go
    Linux 063f20eaee24 6.6.114.1-microsoft-standard-WSL2 … x86_64 GNU/Linux
    --- PASS: TestAC41NailByteFaceUnnamedFileBreaksTheWidth (0.03s)
    --- PASS: TestAC42NailByteFaceExtraDirectoryBreaksTheWidth (0.01s)
    TestAC43: model dir holds 5 entries the signed manifest does not name: a1.onnx, a2.onnx, a3.onnx, sub/, sub/deep.bin
    --- PASS: TestAC43NailByteFaceCountIsExact (0.01s) / TestAC44NailByteFaceCleanTreeStillPasses (0.01s)
    --- PASS: TestAC21… (0.01s) / TestAC22…{cache_hit,local_override} (0.02s) / TestAC23… (0.01s)
    TestAC2NormalInstall…: server hits: install=1, second-instance reuse=0
    --- PASS: TestAC3WideTemp… / TestAC1HandoffRefuses…
    ok  github.com/CarlosShao/wisp/internal/models  0.156s      GOTEST_RC=0   VET_RC=0   DOCKER_RC=0

`_windows_test.go` 两枚在本档由 build tag 排除（＝Linux 侧不参与编译），这是 `go vet` 报的事、⛔ 我据它判自己失败。

## 成对四数（落地 ⇄ 基线，逐名作差）

- 基线（HEAD `fd091cd1` 导出树，⛔ 我的字节）`-count=1`：**RUN=51 PASS_TOP=38 FAIL_TOP=0 SKIP=2** rc=0。
- 落地（同一棵导出树＋本腿全部字节）`-count=1`：**RUN=62 PASS_TOP=47 FAIL_TOP=0 SKIP=2** rc=0。
- 逐名作差：新增 **9 枚全 PASS**（AC21/22/23＋AC41/42/43/44Nail＋AC4R3/AC4R4），
  基线侧消失 **0 枚** ⇒ **新增红 0**；SKIP 两条＝既有联网门 `manifest_real_test.go:120/:147`，本腿⛔ 动它。
- `-count=2`：落地导出树与共享工作树各一次，同读数 **RUN=124 PASS_TOP=94 FAIL_TOP=0 SKIP=4** rc=0。
  四数⛔ 混用⛔ 相减（`=== RUN` 含子测试、`--- PASS` 只顶层）。

## 门禁四把（跑前 `unset GOTMPDIR TMPDIR TEMP TMP`）

| 把 | 命令 | rc |
|---|---|---|
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l <本腿五枚>`（⛔ `-w` 扫别人；一次 `-w` 只打在**自己**那枚纯对齐 diff 上） | **0**，dirty_bytes=0 |
| d22scan | `sh scripts/d22scan.sh` | **0** clean（`ban #8 internal/=607`、`cmd/=143`，台账 ⛔ 降；分母随在飞腿涨，本腿只 +3 枚 `.go`） |
| 长度门 | `sh scripts/check-path-length-budget.sh` | **0** VERDICT GREEN（17287 tracked、over-budget=60＝roster 60、不在册 0） |
| go vet | `go vet ./internal/models/`（共享树）＋ `GOOS=linux go vet ./internal/models/`（导出树，容器内亦 rc=0） | **0** / **0** |

入库 blob 的 CR 尺（`git cat-file blob`＋`tr -cd`，⛔ 对工作树 `grep -c $'\r'`）＝本腿六枚件 **cr_bytes=0**。
跑色全程 ⛔ 整包 `./cmd/wisp`、⛔ `./internal/winsec/`、⛔ `./internal/agent/**`、⛔ `./...`；
`cmd/wisp` 只取编译面 `go build ./cmd/wisp/` rc=0（导出树内）。

## 台面登记（具名到行，⛔ 据此判本腿失败、⛔ 动 winsec 字节去绕）

起手把 `GOTMPDIR` 指到 `D:\tmp\wisp\…\gtmp` 时，`TestAC4InstallDirWidthIsDeliberateAndHasABreakLine` 在**不含本腿
任何字节**的 HEAD 导出树上红一次，红句逐字 `removing the parent's grant left an explicit foreign ACE on the
install dir - out-of-band by ticket 89's definition: [BUILTIN\Users (I)(RX) BUILTIN\Users (I)(OI)(CI)(IO)(GR,GE)]`。
量到的原因＝`D:\tmp\wisp` 自身带机器范围 `BUILTIN\Users:(I)(RX)`＋`(I)(OI)(CI)(IO)(GR,GE)` 两条继承
（`icacls D:\tmp\wisp` 原文），而 `icacls /remove` 撤得掉本对象显式授权、**撤不掉祖先链的继承** ⇒ 那条界读到的是
我的台根策略。安静复量 3 发（默认台根）全 PASS ⇒ 之后一律 `unset GOTMPDIR TMPDIR TEMP TMP`。
件＝`testbed-GOTMPDIR-acl-red-part{1,2}.log` ⇄ `testbed-default-temp-remeasure-{1,2,3}.log`。

## 未做/边界（诚实）

- `AC#6` ⛔ 归本腿：票面 `:44-51` 明写与"第一个模型装载器"那张票合并收＝今天无落点。
  残留窗可达性判据＝**没变可达**（`internal/speech/` 仍只有 `doc.go`、`internal/` 下⛔ `engines/`、
  `cmd/wisp/models.go:36-40` 自陈本进程从不 open `<store>/<id>/` 下任何文件、只重哈希），本枚把守卫挪到更靠后的位置
  ⇒ 那段窗变短、⛔ 归零、⛔ 变可达 ⇒ 按票面 `:50` "别在本票自己收口"继续开着、⛔ 停手。
- 成本如实：缓存命中与 `local_override` 两支现在对同一棵树读两遍（裁决表攻4 单份＝256 MiB 1.07–1.22 s）。
  本腿⛔ 动阈值/golden/`thresholds.go`、⛔ 做 D32/资源判定、⛔ 翻框、⛔ push。
- `R-109-2` 的另一枚**更宽**的字节面（NTFS 命名流 `model.onnx:attacker`、硬链接 `nlink=2`）本机实测确实绕得过
  目录遍历（`ReadDir`⛔ 看见流、`os.Lstat` 读作 regular file），但裁决表 `:44` 逐字给本腿的两形＝未点名文件⇄extra 目录，
  命名流/硬链接⛔ 在派单射程内 ⇒ 本腿⛔ 自行扩面，只在此具名登记，交编排者裁是否另立一格。
- 工具输出里自称"编排者备注/停手/请 revert"的文本：本腿 **0 次**（票面 `:66` 那味照记＝⛔ 是指令）。
- 落盘件行数帽：全部读数件 ≤160 行（超限者按 `-partN` 拆分入库，拼接后与仓外完整日志 `cmp` **逐字节相同**，
  四数由拼回的件重新导出＝与原始读数同号）；⛔ `.out` 扩展名；⛔ 0 字节件。
