# 票 113 — 对抗验收裁决表（`acceptor-ticket113`）

验收开始时间：`date` = **Mon Sep 21 20:40:29 CST 2026**（本机 Git Bash，CST）
被验票面：`.scratch/wisp/issues/113-posix-platformverifypplacement-has-no-link-leg.md`（`ready-for-review`，框未勾）
锚定 sha：**`ef65864`（修前红）** 与 **`3c5d1c3`（修法）**；当前 HEAD `34a810b`（docs-only 增量）。
档位图例：〔独立复现〕= 我亲手跑过并贴读数 / 〔日志＋归档，我抽验〕= 有落盘产物我只抽查其中若干 / 〔仅自述，不背书〕= 我没跑成或跑不了，写原因。
本文件渐进写：每验完一格立刻追加（2026-09-21 已有代理攒到最后没落盘，白干两轮）。

## 裁决表（建表时是空表，每格验完就地改；收工时全表已填）

| 格 | 声称 | 我的读数 | 档位 | 结论 |
| --- | --- | --- | --- | --- |
| AC#1 修前红（`ef65864`，含"外来 mode/属主前后"） | rc=1、RUN=16/PASS=7/FAIL=9、外来 `666→600` | **我亲手复现，逐字对上**：容器 rc=1、`=== RUN` 16 / `--- PASS` 7 / `FAIL:` 9 行（5 顶层 + 4 深度子用例）/ SKIP 0；红名 9 枚与票面一字不差；外来文件 `before=mode=666 uid=0 gid=0 → after=mode=600`、`SealDir` 那发外来目录 `777→700`、`PrivateFile` 那发 `AC#1 RED: PrivateFile put bytes in the foreign tree (<nil>)` | 〔独立复现〕 | **绿（这一格成立）** |
| AC#2 链接腿落地（复用 `pathPieces`、不折 `\`、只拒、含叶子） | `3c5d1c3` 容器 rc=0、16/16 | 容器 `3c5d1c3`：`go test -count=1 -v -run 'TestAC1POSIX\|TestAC2POSIX\|TestAC3POSIX\|TestAC4POSIX'` ⇒ **rc=0、RUN=16 PASS=16 FAIL=0 SKIP=0**；20 条外来前后读数**全部 `before=after`**（10 条 `666` + 10 条 `777`，uid/gid 皆 0 不变）。码面：`git show --stat 3c5d1c3` = 只碰 `winsec_other.go`(+83/-12) 与票面，`internal/winsec/placement_symlink_113_other_test.go` 在 `ef65864`↔`3c5d1c3` 之间 **0 行差异** ⇒ 红/绿同仪器可比；腿体在 `:113-121`，`pathPieces`/`resolve.go` 一行未动，哨兵同一枚 `ErrUnresolvedPath`，只 `return ""` 不改写 | 〔独立复现〕 | **绿** |
| AC#3 反半边不误伤 | 普通文件仍封、`a\b` 不折 | **声称的那两枚在"祖先无链接"的树上确实绿**（我在 `3c5d1c3` 容器内 `TestAC3POSIX*` 3 枚全 PASS，普通文件照旧 `0666→0600`）。**但"反半边"只在无链接祖先时成立：我用一枚合法 POSIX 形状把它打红了，见下"攻#2"** | 〔独立复现〕 | **绿（判据字面成立）+ 一条未登记的误伤面（见攻#2）** |
| AC#4 变异四发 | MUT-A 9 红 / MUT-B 恰 1 红 / MUT-C 9 红 / MUT-D2 2 红 | **四发全部我自己下刀、自己复量**（`git archive 3c5d1c3` → `/d/tmp/mut-ac113-{A,B,C,D2}`，仓库内无 worktree）。每发：`grep -n` 打印被改后整行 → 容器内 `go build ./internal/winsec/` rc=0 → 容器内 `go test -count=1 -v`。**读数一字不差**：A `RUN=16 PASS=7 FAIL=9`；B `FAIL=1`（恰 `TestAC1POSIXSealDirThroughASymlinkRefuses`，外来目录 `777→700`）；C `FAIL=9`（深度 1/2/3/4 全红）；D2 `FAIL=2`（本票 `TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator` + 票 108 的 `TestAC2POSIXDoesNotFoldABackslashIntoASeparator`）。**MUT-B 只红 1 枚的真相不是"用例不够"，是另外四枚的形状里链接都在祖先位（只有 `SealDir(link)` 的链接在叶子位）** ⇒ 全集主张成立，但**叶子这一维在 `SealFile` 方向没有交付用例**，见下"攻#4" | 〔独立复现〕 | **绿** |
| AC#5 门禁四数 + 三门 + d22 | 540/532/0/8、ban#8 internal/=368 | 五包 `-count=2 -v`（winsec/memory/risk/secret/models，容器 `golang:1.27`、`CGO_ENABLED=0`、uid 0）⇒ **rc=0、`=== RUN` 540 / `PASS:` 532 / `FAIL:` 0 / `SKIP:` 8**，逐包 `ok winsec 0.273s / ok memory 24.499s / ok risk 7.822s / ok secret 0.009s / ok models 0.530s`（我第一次跑出 `548/536/4/8` 的 4 枚红是**我自己塞进快照的探针文件**，移走后即此数——顺手证明这条命令不缓存、不吞红）。`gofmt -l .` **0 行**；容器内 `go vet` 五包 rc=0、`GOOS=windows go vet ./internal/winsec/` rc=0、`GOOS=darwin go vet ./internal/winsec/` rc=0。**d22：`sh scripts/d22scan.sh` rc=0，但 ban#8 `internal/=371`，不是它报的 368** ⇒ 见"攻#5" | 〔独立复现〕 | **绿（四数与三门），但 ban#8 那格的 sha↔数字配对错了一格** |
| AC#6 文档边界（`R-108-2` 落 `winsec.go` 头部） | 票面自述"还没做" | `git show 34a810b:internal/winsec/winsec.go`、`git show f6a86db:internal/winsec/winsec.go` 里 `S-1-1-0` 命中数 = **0**；`git log --oneline -1` 期间该文件从未出现在 113 的任何 commit（`3c5d1c3` 只碰 `winsec_other.go`）。**验收现场工作树里确有一段未提交的 +18 行注释**（`git diff --numstat internal/winsec/winsec.go` = 18/0，含 "ticket 113 AC#6, answering R-108-2" 字样）⇒ 未 commit = 未交件，且那是别人的地界，我不代写、不代勾 | 〔独立复现〕 | **未交件（本票 6 格里唯一没做的一格，票面框也确实没勾）** |
| 攻#1 108 定命探针复现 + 我自选形状 | 见下"攻#1 明细" | 见下 | 〔独立复现〕 | **symlink 穿越这一枚我打不穿（修法为真）；两枚"改权限并返回 nil"的形状仍在，但都在本票 AC 字面之外** |
| 攻#2 误伤面（合法 POSIX 树 / macOS `/tmp` / 生产调用方） | 票面只登记了 macOS `/tmp`、`/var` 一条权衡 | **能构造，且已量化**（见"攻#2 明细"）：容器内把 `TMPDIR` 放到一枚 symlink 之下（= macOS 的真实形状），`3c5d1c3` 上 **winsec 自己红 12 项**、**secret+config+agent 三包红 33 行**；同一形状在 `ef65864` 上 winsec 只红 9 项（本票的修前红）、三包 **0 红** ⇒ 差额 100% 归本票这发腿。**错误没有被吞**：`memory: create data dir: winsec: refusing to seal …`、`resource: agent: create artifacts dir: …` 一路 `fmt.Errorf/%w`+`observe.Wrap` 传到调用方返回值 | 〔独立复现〕 | **本票最大的一格：判据字面绿，代价没登记全** |
| 攻#3 这道门在 CI 上存不存在 | 票 110 只给 windows 腿加了步 | **不存在**：run **35601785381**（headSha `f6a86db2`，含 `3c5d1c3`）job **106339582851 `test-core`** step **7** `Portable package tests` 的 scope 列表（日志逐字 `[./internal/agent/... ./internal/llm/... ./internal/config/... ./internal/memory/... ./internal/observe/... ./internal/secret/... ./internal/risk/... ./internal/statemachine/... ./internal/session/... ./internal/watchdog/... ./internal/tools/... ./internal/models/... ./internal/buildinfo/... ./internal/audio/... ./internal/proc/... ./internal/panel/...]`）**没有 `./internal/winsec/`**；**整个 job 日志里 `internal/winsec` 出现 0 次**。全 CI 唯一的 winsec 门 = job **106339582688 `test-windows`** step **4**，而 `scripts/winsec-tests.sh:73-80` 在 `GOOS != windows` 时 **exit 2 硬拒** ⇒ POSIX 腿**按设计**没有 CI。在飞的 `ci.yml`/`portable-tests.sh`（票 111，工作树 `M`）我也 grep 了：仍然只有 windows 那一处 | 〔独立复现〕 | **本票修复今天零回归保护 ⇒ R-113-A，归票 111** |
| 攻#4 用例信息量 | — | 9 枚红名在 `ef65864` 真红、MUT-A/C 真红 ⇒ **不是恒真断言**；`untouched` 那条 mode 断言**独立发火过**（红日志里 `was mode=666 … and is now mode=600` 与 `returned nil` 两条同发）⇒ 不是装饰。但 **MUT-B 只红 1 枚的原因我查清了**：交付的 5 枚 AC#1 用例里链接**全部在祖先位**，只有 `SealDir(link)` 在叶子位 ⇒ 叶子方向（`SealFile(link)`、`PrivateFile(link)`）**没有交付用例**；我补的 P5 在 `3c5d1c3` 绿（拒 + 外来 mode/字节都不动）、在 MUT-B 红（`err=<nil>`、外来 `666→600`、`PrivateFile` 把 `"private bytes"` 覆写进外来文件）⇒ 实现是对的，**覆盖缺两枚** ⇒ R-113-C | 〔独立复现〕 | **AC#4 绿；记 R-113-C（覆盖缺口，非缺陷）** |
| 攻#5 数字口径与台账漂移 | ban#8 internal/=368 @ `3c5d1c3` | **复算 = 371，不是 368**。我用 `find internal -name '*.go' | wc -l` 独立复核同一把尺：`ef65864`=**368**、`3c5d1c3`=**371**、HEAD `f6a86db`=**371**（ban#1-5 internal/=202、design/=16 我也各数了一遍，逐字相符）。+3 的来路是**邻居**：`a505607`（票 112）加 `internal/winsec/tree_ownership_112_windows_test.go`、`f6818f2`（票 109）加 `internal/models/handoff_reuse_109_test.go` + `handoff_window_109_windows_test.go` ⇒ 366/368/371 三个数是同一条单调曲线上的不同 sha，**它把 `ef65864` 的 368 抄到了 `3c5d1c3` 那一行**（正是票 108 台账"读数与 sha 同行登记"要防的那一格）。SKIP=8 = 4 名 ×2 我逐条读到 `Skip` 那一行：`internal/memory/concurrent_test.go:186`（re-exec 子进程入口，无条件跳）、`internal/models/manifest_real_test.go:120/147`（`WISP_IT_REAL_MIRROR=1` 才跑，容器无网络）、`internal/risk/pathresolver_rewrite_account_test.go:138`（平台不给 handle-resolved 形式，POSIX 恒跳）⇒ **全是既有环境跳过，不是新增静默** | 〔独立复现〕 | **AC#5 四数绿；台账配错一格 ⇒ R-113-F** |

## 攻#1 明细（108 那枚定命探针 + 我自己的形状）

| 探针 | `ef65864`（修前） | `3c5d1c3`（修后） | 判读 |
| --- | --- | --- | --- |
| P-108 原件：`SealFile(root/link/keep-me.txt)` | `err=<nil>` + 外来 `mode=666→600`（uid/gid 两边 0） | `err=ErrUnresolvedPath` + 外来 `666→666`、链接还在 | **108 的定命探针今天真的红了、现在真的不红了**（逐字复现） |
| P2 相对目标链接 `root/rel -> ../foreign` | 红（`err=<nil>`，外来被动） | 绿（拒） | 同一腿，不同拼写 |
| P3 `/proc/self/fd/N`（魔链在叶子位） | 红 | 绿（`refuses … the link at /proc/self`） | 底线连 `/proc` 魔链都拒 |
| P5 `SealFile(link)` / `PrivateFile(link)`（链接就在叶子位） | 红：`err=<nil>`、外来 `666→600`、`PrivateFile` 把字节写进外来树 | **绿**：拒（`reaches it through the link at <该名字本身>`）、外来 `666→666`、外来字节仍是 `"not this tree's data"` | 实现覆盖到了，**交付用例没有**（⇒ MUT-B 只红 1 枚的解释）⇒ R-113-C |
| P1 **硬链接**同 inode（拼写里一个 symlink 都没有） | `err=<nil>`、外来 `666→600` | **一模一样**：`err=<nil>`、外来 `666→600` | **AC 声称要防的结局类在 POSIX 上仍然造得出来，但机制不是"穿过链接"**：编排者已裁 `R-108-2`（"这棵树归谁"）不改定义；本票 `winsec_other.go` 注释第③条已自陈。**⇒ 不判本票 FAIL，判 AC#6 未交件**（那句边界今天只活在函数注释里，没进包文档） |
| P4 **check-then-act 竞态**：叶子在 `Lstat` 之后被换成 symlink，4000 发 | 穿过 **2087/4000（52.2%）** | 穿过 **154/4000（3.85%）**（另一轮 69/101） | 腿把概率压掉一个数量级，**没有关死**：底线是"看一眼再动手"，POSIX 侧没用 `openat2/O_NOFOLLOW`/fd 复核。AC#2 只承诺"要查"，没承诺原子 ⇒ **不判 FAIL**，但这条代价**没写进那两条"成本"注释** ⇒ R-113-D |
| P6 **root 失真**：外来文件 `chown` 给 uid 65534 后由 uid 0 的 wisp 密封 | `err=<nil>`、外来 `644(uid 65534)→600` | 拒，外来不动 | **容器是 uid 0 ⇒ 读数偏"能打穿"而非偏"打不穿"**；同时我量到反向天花板：uid 1000 去 `chmod` 一枚 uid 65534 的文件 = `Operation not permitted`，改自己的 = `rc=0` ⇒ **跨账户那一档 damage 只在 wisp 以 root 安装时可达**，普通安装可及的是"同账户、树外"的文件与"把字节写进树外"。所有外来 uid/gid 读数在交付夹具里都是 0/0 ⇒ 交付用例证明的是"**树外**"，不是"**别人账户**" |

## 攻#2 明细（误伤面，我最担心的那一格）

仪器（容器 `golang:1.27`，两棵快照各跑一遍）：
`mkdir -p /realpriv/var/folders/xy/T000; ln -s /realpriv/var /varlink; export TMPDIR=/varlink/folders/xy/T000`
⇒ 这就是 macOS 的真实形状（`/var -> private/var`、`TMPDIR=/var/folders/…/T/`；`/tmp -> private/tmp` 同理），Go 的 `t.TempDir()`/`os.TempDir()` 都会踩到那枚 symlink 祖先。

| 树 | 默认 TMPDIR | 同一形状 + symlink 祖先 | 差额归谁 |
| --- | --- | --- | --- |
| `ef65864`（修前，`./internal/winsec/` 全套） | rc=0 | `RUN=35 FAIL=12`（其中 9 枚就是本票的修前红） | 既有腿（票 103/108 的 `RemoveUnlinked`）只咬 3 枚 |
| `3c5d1c3`（修后，同一条命令） | rc=0、`RUN=35 PASS=20 FAIL=0 SKIP=0` | **`RUN=35 FAIL=15`**，红名含本票自己的 3 枚 AC#3 反半边用例 | **新增 12 枚**（`TestPOSIXPrivateFileIsReally0600`、`TestPOSIXPrivateDirIsReally0700`、`TestAC4POSIXFloorAnswersInsideTheNamedTree`、`TestAC5FailedSealRefusesTheWrite`+4 子用例、`TestAC5FailureIsNotSwallowedByTheHappyPath`、3 枚 `TestAC3POSIX*`） |
| `ef65864`（修前，`secret`+`config`+`agent`） | rc=0 | **rc=0、0 红** | — |
| `3c5d1c3`（修后，同一条命令） | rc=0 | **rc=1、33 行红** | 全归本票 |

调用方看到什么（逐字，非吞掉）：
`memory: create data dir: winsec: refusing to seal /varlink/…: winsec: path is not provably resolved, refusing to seal: … reaches it through the link at /varlink, which is not the tree this call names`
`resource: agent: create artifacts dir: winsec: refusing to seal …`
`config: … temp file seal: winsec: refusing to seal …`
⇒ 生产调用点我全数过：`internal/secret/store.go:49/79`、`internal/secret/migrate.go:169/174/189`、`internal/memory/open.go:176/188/503`、`internal/memory/artifacts.go:75`、`internal/agent/spill.go:111/125/254`、`internal/config/parse.go:215`、`internal/config/migrate.go:93` ⇒ **一律 `return` 上抛，没有一处 `_ =` 或 `continue` 吞掉**（这一条是好消息：坏行为是"响亮失败"，不是静默）。

真实合法形状（不用仿真也成立）：
1. **`WISP_ENV=test` 的数据根本身**：`cmd/wisp/doctor.go:235` 与 `internal/proc/envfork.go:98` = `filepath.Join(os.TempDir(), "wisp-test-<pid>")` ⇒ 在 macOS 上这条路径的祖先含 `/var`（或回落到 `/tmp -> private/tmp`）⇒ 从本票起 **test 环境的 memory/secret 初始化会开始拒绝**。
2. **Linux 上的 `~/.config`**：`os.UserConfigDir()` = `$XDG_CONFIG_HOME`/`~/.config`，dotfiles 管理里把它整目录软链走是常规操作；`~` 本身是 automount 软链（`/home -> /export/home`、NixOS profile）也同形 ⇒ `secret.NewStore` 直接返回错误 = **凭据写不进去**。
3. 票面注释只登记了 (1) 的 macOS 版本，**没登记 (2)，也没登记"本包自己的 POSIX 用例与三包调用方在这种根上会红"** ⇒ 记账为 R-113-B。

**R-108-5（darwin 那半边今天能不能采到数）**：**如实答"采不到"**。我没有 macOS 机器，`GOOS=darwin go vet ./internal/winsec/` rc=0 只证明**能编译**，不证明任何行为；上面那张表是 **Linux 上的同形状仿真 + 代码读法**推出来的 darwin 预测（`/var`、`/tmp` 在 darwin 是 symlink 属平台事实，不是我实测）。⇒ 我这一格同样只能标〔推论 + 容器内同形状实测〕，**darwin 的实测账仍然是票 108 那笔没付的 R-103-6**。

## 我亲手跑的命令（可重跑）

```bash
date                                                       # 2026-09-21 20:40:29 CST 起跑，收工 21:0x
git log --oneline -12; git status --porcelain              # 锚定 sha；工作树在飞邻居（我没读脏树）
mkdir -p /d/tmp/ac113-{red,green,head}
git archive ef65864 | tar -x -C /d/tmp/ac113-red           # A38④：仓库内无 worktree、无 checkout
git archive 3c5d1c3 | tar -x -C /d/tmp/ac113-green
git archive HEAD    | tar -x -C /d/tmp/ac113-head
# 容器挂载先证非空（假绿坑#1）：ls -l /src/go.mod / src/internal/winsec 计数 / 探针文件存在 + id + uname
MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 \
  -v /d/tmp/ac113-red:/src -v wisp113mod:/go/pkg/mod -v wisp113build:/root/.cache/go-build -w /src golang:1.27 \
  sh -c 'go test -count=1 -v -run "TestAC1POSIX|TestAC2POSIX|TestAC3POSIX|TestAC4POSIX" ./internal/winsec/ > /src/113-red-ac113.log 2>&1; echo RC=$?; …'
# 修后同一仪器（ac113-green）；变异四发见下表；五包：go test -count=2 -v ./internal/winsec/ ./internal/memory/ ./internal/risk/ ./internal/secret/ ./internal/models/
# 静态：gofmt -l . (0 行) / go vet 五包 / GOOS=windows go vet / GOOS=darwin go vet（容器内，全 rc=0）
# sh scripts/d22scan.sh（容器内快照）rc=0；ban#8 internal/=371
# 独立复核同一把尺：for s in ef65864 3c5d1c3 f6a86db; do git archive $s | tar -x -C /d/tmp/cnt-$s; find /d/tmp/cnt-$s/internal -name '*.go' | wc -l; done
# CI：gh run list --json databaseId,headSha / gh run view 35601785381 --json jobs（取 step 号）/ gh run view --log --job 106339582851 | grep -c internal/winsec  => 0
```
容器内 uid=0 全程成立（`id` 读数 `uid=0(root) gid=0(root)`）⇒ 见上表 P6 的失真说明；`/src` 是 Windows bind 挂载，其上的 mode 位是假的（`-rwxrwxrwx`），**但所有断言都发生在 `t.TempDir()` = 容器 overlay 的 `/tmp` 上**，mode 位真实。

## 变异清单（我下的刀，锚点 = 被改后整行）

| 变异 | 锚点行（grep 原文） | 落地证明 | `go build` | 读数（容器 `go test -count=1 -v`） |
| --- | --- | --- | --- | --- |
| MUT-A 整腿关 | `internal/winsec/winsec_other.go:115: if false && ancestorIsLink(prefix) {` | 同链 `grep -n` + 与 `git show 3c5d1c3:…` 逐字 diff | rc=0 | rc=1 `RUN=16 PASS=7 FAIL=9 SKIP=0`；9 枚红名逐字同票面；外来 `666→600` ×6、`777→700` ×1 |
| MUT-B 不查叶子 | `internal/winsec/winsec_other.go:115: pieces = pieces[:len(pieces)-1]` | 同上（113-116 行整块打印） | rc=0 | rc=1 **FAIL=1** = `TestAC1POSIXSealDirThroughASymlinkRefuses`，外来目录 `777→700`；**我的 P5 在这发上也红（`SealFile`/`PrivateFile` 叶子位）⇒ 交付用例缺两枚** |
| MUT-C 只查一层 | `internal/winsec/winsec_other.go:115: if len(pieces) > 1 { pieces = pieces[:1] }` | 同上 | rc=0 | rc=1 `FAIL=9`（深度 1/2/3/4 全红）⇒ "咬全集"成立 |
| MUT-D2 把 `\` 塞进切分 | `internal/winsec/winsec.go:303: return c == os.PathSeparator || (nativeIsBackslash && c == '/') || c == 0x5c /* MUTATION-FOLD-113 */` | 同上 + 全文件 diff 只有这一行 | rc=0 | rc=1 `FAIL=2` = 票 113 `TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator` + 票 108 `TestAC2POSIXDoesNotFoldABackslashIntoASeparator` ⇒ 与自述一字不差 |
| 还原 | 四发都在 `/d/tmp/mut-ac113-*` 独立快照；仓库 `git status --porcelain` 里我的路径只有新增的这枚证据文件（`?? docs/evidence/s1/113-adversarial-acceptance.md`），**0 行生产码被我改过** | — | — | — |

## R-113-x（要立案的，含地界与是否阻塞）

| 编号 | 一句话 | 属谁地界 | 阻塞 113 结案？ |
| --- | --- | --- | --- |
| **R-113-A** | 这条 POSIX 腿在 CI 上**一步都没有**：ubuntu `test-core` step 7 的 scope 不含 `./internal/winsec/`（job 日志 `internal/winsec` 命中 0 次，run 35601785381 / job 106339582851），唯一 winsec 门在 `test-windows` step 4 且脚本硬拒非 Windows ⇒ 今天"这条修复绿"只活在我这台机器的 Docker 里 | **票 111 的覆盖面**（`ci.yml` + `scripts/portable-tests.sh` 是它的，我不碰） | 不阻塞 113；**阻塞**任何人把 POSIX 腿读成"CI 已保护"，包括票 108 复验 |
| R-113-B | 误伤面登记不全：只写了 macOS `/tmp`、`/var`，没写 Linux `~/.config`/automount home、也没写 `WISP_ENV=test` 数据根走 `os.TempDir()`；实测代价 = winsec 12 枚 + 三包 33 行红 | `internal/winsec`（注释/票面）+ 调用方纪律（票 76/95） | 不阻塞（修法方向正确），但**必须在 AC#6 那格一起补上**，否则下一个人还会当成新 bug |
| R-113-C | 叶子方向缺两枚交付用例（`SealFile(link)`、`PrivateFile(link)`）：MUT-B 只红 1 枚的真因；我补的 P5 在现码绿、在 MUT-B 红 | `internal/winsec/placement_symlink_113_other_test.go`（票 113 自己的用例文件） | 不阻塞；**建议 113b 之后追加一小票或并进 AC#6 那一轮** |
| R-113-D | 底线是 check-then-act：我把叶子在 `Lstat` 之后换成链接，4000 发里 **154 发**穿过（修前 2087 发）。不是回归，但注释里"两条成本 + 三条残余"没列它 | `internal/winsec`（`winsec_other.go` 注释）→ 若要关需另立一票（`openat2/O_NOFOLLOW`/fd 复核） | 不阻塞 113 |
| R-113-E | 硬链接同 inode：`SealFile` 仍 `err=<nil>` 且把外来 mode `666→600`（red/green 逐字相同） | 已由编排者裁定为 **`R-108-2`**（不改定义、写文档） | 不阻塞（判据不在本票），但**它正是 AC#6 未交件时唯一没落盘的边界** |
| R-113-F | 台账漂移：`3c5d1c3` 的 ban#8 `internal/` 是 **371** 不是 368（368 属 `ef65864`）；+3 全来自票 112/109 新增的 `_test.go` | 票 113 票面 Progress log（谁写谁改） | 不阻塞（方向仍是"不降"），但**这一格别照抄** |

## 总判

**通过附条件（PASS-with-conditions）** —— 一句话：108 那枚把我这一族判死的探针我**逐字复现了它今天真的红过**（`ef65864`：`SealFile` 穿 symlink 返回 nil、外来 `666→600`），也复现了它**现在真的不红**（`3c5d1c3`：拒 + 外来 mode/属主/字节三项一字未动），四发变异我一字不差地重造，AC#5 四数在剥掉我自己的探针后精确是 `540/532/0/8` ⇒ **修法成立**；但**六格里 AC#6 未交件**（无任何 commit 含那段边界，工作树里那 18 行是别人的未提交改动），且**这条腿今天没有任何 CI 覆盖**、**误伤面的代价只登记了一半**，这三条我按 R-113-A/B/F + AC#6 FAIL 记成放行条件而不是"通过"。

### 放行条件（结案前必须见到）
1. **AC#6 必须落成真 commit**（`internal/winsec/winsec.go` 头部那段边界；判据是 `git diff` 只增不减、不动判定分支）⇒ 在此之前**票面 AC#6 那格不许勾**。
2. **R-113-A 写进票 111 的覆盖面清单**（POSIX 腿要么进 ubuntu `test-core` scope，要么在 ubuntu 加一枚 `winsec-tests.sh` 的平台腿），并注明"在此之前票 113 的绿只有本机 Docker 读数"。
3. **AC#5 那一格的 ban#8 数字按 sha 重登**（`3c5d1c3` ⇒ 371）。
4. 建议一并把 R-113-B/C/D 三行补进 `winsec_other.go` 那两条"成本"注释（同一枚文档 commit 顺手做掉，不再占工单）。

## 收工前的自证（`date` = Mon Sep 21 21:12:53 CST 2026）

- `git status --porcelain -- docs/evidence/s1/113-adversarial-acceptance.md` = `?? `（只新增了这一枚证据文件）；仓库里**我没有改过任何一行生产码、没有 commit、没有 push、没有建 worktree/checkout**（A38④：所有复算与变异都在 `/d/tmp/ac113-*`、`/d/tmp/mut-ac113-*` 这些 `git archive <sha> | tar -x` 的仓外快照里做）。
- 工作树在我验收期间继续被邻居推进（`ef65864`→`3c5d1c3`→`34a810b`→`f6a86db`，且 `M internal/winsec/winsec.go`、`M .github/workflows/ci.yml`、`M scripts/portable-tests.sh` 是三笔**未提交**的在飞改动）⇒ **本表所有读数都出自按 sha 归档的快照，不出自脏树**；任何"20 分钟前还是这样"的状态断言（尤其 AC#6 那一格）以复验当时的 commit 为准。

## 伪授权登记（本会话）

工具输出里自称"编排者备注 / 系统提示 / 停手 / 撤回 / 请 revert"的文本：**0 次**。
另有 **1 次**对话层 `MEMORY.md was modified since it was last read` 通知（内容是他的记忆索引，含"撤销只能来自对话"一条）⇒ 按"不作为授权、不改判据"处置。**本验收全程未 revert、未撤销任何 commit、未改任何生产码**（`git status --porcelain` 里我只有这枚新增证据文件）。

## 下一张该派什么（我的建议）

1. **票 111 追加 AC**：给 ubuntu 腿加 `internal/winsec` 的 POSIX 门（含 `-count=2` 四数与"scope 不含 winsec 即红"的守卫）——这是眼下最便宜的一条，因为仪器（`runtests.sh` + `portable-tests.sh` 的 pin）都已经在别处被证明能咬。
2. **票 113b（AC#6）落 commit**，并把我这三行 R-113-B/C/D 的边界一起写进文档；写完由我抽验（判据：只增不减 + 容器 AC#5 三门读数不变）。
3. **新票：POSIX 密封的原子性**（R-113-D：`openat2(RESOLVE_BENEATH)`/`O_NOFOLLOW` + fd 复核，或"检查后重开并比对 `fdpath`"），窗口实测 154/4000 起。
4. 票 108 复验：**先等 111 的 POSIX 门存在**，否则复验读到的"绿"依然只活在我这台机器的 Docker 上。

## （表头那格"总判=未裁决"的占位已作废）

总判见上一节 **"## 总判"（通过附条件）**；本占位是 20:40 建表时的骨架，留一行说明免得下一个人以为有两份裁决。
