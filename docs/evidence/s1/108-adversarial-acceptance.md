# 票 108 独立对抗验收裁决表（acceptor-ticket108）

- 验收对象：票 `.scratch/wisp/issues/108-103s-seal-guard-is-bypassable-three-ways.md`，实现方 `agent-ticket108`，commits `87f6f6d`（修前红）+ `c02c609`（修法）。
- 验收时间：2026-09-21 19:34 起（`date` 取值 `Mon Sep 21 19:34:07 CST 2026`；收尾时 HEAD 已被邻居推进到 `21da8f3`）。
- 我是**新代理**，未读先前对话；前情只读票面 + `docs/evidence/s1/103-adversarial-acceptance.md`（P1b/P2/P3 的原始形状）。
- **总判：FAIL（不通过，退回）**。理由见"总判"一节：Windows 四格我都独立量到绿，但 **AC#3 声称要防的结局我在 POSIX 上真实造出来了**（`SealFile` 穿过 symlink 改掉外来文件的权限并返回 nil）。
- 档位图例：〔独立复现〕= 我在 `/tmp` 仓外快照或容器里自己跑出来的读数；〔日志＋归档，我抽验〕= 实现方日志 + 我抽验其中关键量；〔仅自述，不背书〕= 我无法复现，且不承担背书。

## 我的仪器（先交代，防"守卫拒一切"被当成绿）

- 快照全部在仓外：`/tmp/ac108`（探针）、`/tmp/gate108`（门禁，纯净 `git archive HEAD`）、`/tmp/mut108-ac108`（变异）。**仓库内未建 worktree、未 checkout、未改任何生产码**（A38④）。
- 我自己的探针文件两枚（**只存在于 `/tmp`，不进仓**）：
  `C:\Users\swq\AppData\Local\Temp\ac108\internal\winsec\zz_accept108_windows_test.go`、`...\zz_accept108_other_test.go`。
- 阳性对照同批跑：我自己的 `TestAC108NewSealShapesKeepForeignDaclAndStillNarrowOurs` 末腿证明**守卫没有打死正常密封**——
  我们自己那颗故意宽授权的文件 `icacls` 读数 `before=[S-1-1-0 S-1-5-32-544 S-1-5-18 S-1-5-21-...-1001]`
  → `after=[S-1-5-18 S-1-5-32-544 S-1-5-21-...-1001]`（`S-1-1-0` **由有到无**，即该收窄的真收窄了）。〔独立复现〕

## AC#1 缝不可解除（P1b：先 `SetPathResolver(nil)` 解除再装恒改写伪造）

| 判据 | 结论 | 证据与档位 |
| --- | --- | --- |
| `SetPathResolver(nil)` 在已装/已闩后一律拒（单向语义） | **通过** | 我读 `resolve.go:127-183`：`nil` 分支只在 `resolver == nil && !seamLatched` 时放行，其余走 ERROR 审计 `refusing to release the sealing path resolver` 并**原样保留在位者**；一次性由 `seamLatched` 决定而不是"槽位当前是不是 nil"。〔独立复现（读码 + 下面两条跑）〕 |
| 探针 P1b 修后必须拒，且判据是"外来 DACL 没被动过" | **通过** | 基线（HEAD 纯净快照 + `-v -run 'TestAC1\|TestAC2\|TestAC3\|TestAC4'`）：`rc=0 RUN=56 PASS=56 FAIL=0`，含 `TestAC1SeamCannotBeFreedThenGivenATreeMovingFake` **PASS**、`TestAC1SeamIsOneWayUseIsTheOnlyDirection` **PASS**。外来 victim `icacls` SID 级 `before=[S-1-1-0 S-1-5-32-544 S-1-5-18 S-1-5-21-...-1001]`，同形状守卫生效后 `after=` **同一个列表**（`S-1-1-0` **由有到仍有**）。〔独立复现〕 |
| 闩锁有没有 race（一次性/单向在并发下是否还成立） | **通过** | 我自己写的 `TestAC108SeamLatchIsRaceFree`：48 个 goroutine 混打 `SetPathResolver(nil)` / `SetPathResolver(伪造 latchFake108{C:\Windows\win.ini})` / 读 `PathResolverInstalled()`，`go test -count=2 -race` ⇒ `rc=0`、`=== RUN 2`、`--- PASS 2`、`--- FAIL 0`、`WARNING: DATA RACE 0`，伪造**一次都没进缝**、槽位类型未变手。注：`resolve.go:128-135` 把探针放在锁外跑，但**改槽位的动作全在 `resolverMu.Lock()` 内**，所以"锁外探测"不构成解除窗口。〔独立复现〕 |
| R-103-4（恒等重装把账抹成第 0 份） | **通过（附一条如实登记）** | `name()` 用的是 `%T` ⇒ "同一类型的不同实例"会走"identically installed"分支；该分支**不换槽位**（保留第一份），方向是安全的一侧，但"identity"这个词在类型级别而非指针级别，实现方自己钉的 `countingResolver` 指针身份用例我在基线 56 里算到 PASS。登记为 `R-108-4`（措辞/记账精度，非 fail-open）。〔独立复现（读码）〕 |
| 生产码有无解除路径 | **通过** | `internal/winsec/` 非测试文件里只有 `SetPathResolver` 一处写 `resolver`/`seamLatched`；`SetSeamForTest` 在 `export_test.go`（`go build ./...` 产物不含）。安装点唯一：`internal/risk/winsec_c26.go`。〔独立复现（grep + build）〕 |

## AC#2 祖先链对所有分隔符形状生效（含 POSIX 不折 `\` 的反向用例）

**A. 三枚旧探针的形状（重跑，不是背书）**：基线 56 全绿 ⇒ `all-forward-slash` / `mixed-separators` / `doubled-separator` / `volume-only-relative-tail` 都拒了。〔独立复现〕

**B. 我自己新造的 15 枚形状**（Windows 真机 junction，`cmd /c mklink /J`，全在 `t.TempDir()` 下，测完随 tmp 清）：

| 新形状 | 结果 | 读数 |
| --- | --- | --- |
| N2 `\\?\C:\...\link\sub\keep-me.txt` | **守住** | 拒 `ErrIsReparsePoint`，点名链接祖先 `\\?\C:\...\data\link` |
| N3 `\\.\C:\...`（设备前缀） | **守住** | 同上（祖先被点名） |
| N4 `\\localhost\Users\...`（UNC 拼写） | **守住** | 同上 |
| N6 小写盘符 + 全 `/`（`c:/.../link/sub/keep-me.txt`） | **守住** | 拒并点名 `c:/.../data/link` |
| N7 祖先尾点 `link.\sub\...` | **守住** | 拒（Win32 剥尾点后仍命中链接） |
| N8 祖先尾空格 `link \sub\...` | **守住** | 拒 |
| N10 混合 + 双分隔符 `data\link//sub\keep-me.txt` | **守住** | 拒 |
| N11 全大写整条路径 | **守住** | 拒 |
| N12 祖先里插 `..\`（`data\link\..\link\sub\...`） | **守住** | 拒（Lstat 走 OS 解析，仍命中链接） |
| N13/N14 **叶子本身是链接 + 尾分隔符**（`link\`、`link/`） | **守住（我原以为这里是洞）** | 见下条专项 |
| N15 8.3 短名祖先（`TESTAC~1\002\data\LINK\sub\...`） | 不成立为绕过 | 该子用例跑在 N13/N14 之后，链接已被前面合法解掉 ⇒ 路径指向不存在的名字，`removeUnlinked` 对 `ErrNotExist` 返回 nil（既有契约"已经不在了算成功"）。外来文件/目录/`icacls` 三条断言**全部未红**。 |
| N5 `\\Users\swq\...`（我拼错的准 UNC） | 不是绕过 | 返回 nil，但三条外来断言（文件在、目录在、DACL 逐 SID 相同）**全绿** ⇒ 什么都没动；这是"路径不存在 ⇒ 已 Gone"的既有语义，不是祖先漏查。 |
| N9 `%2F` URL 风格分隔符（`link%2Fsub%2Fkeep-me.txt`） | 不是绕过 | OS 不解码 ⇒ 这是一个不存在的名字 ⇒ nil（同上，外来三条断言未红）。 |

〔独立复现〕`=== RUN` 计数与逐子用例 PASS/FAIL 我逐条看过；上面"不是绕过"三枚我**没有**当成守卫的功劳，也没当成缺陷。

**C. 专项：尾分隔符会不会让守卫"顺着链接删别人的目录"？**（这是"修法的新形状"里我认为最可能中的一项，两侧都量）

- Windows：`TestAC108TrailingSeparatorMustNotFollowTheLinkIntoAForeignTree`（链接指向一个**外来空目录**，`link\` 与 `link/` 两种尾分隔符）⇒ `--- PASS` 两个子用例：外来目录 `os.Stat` 仍在，链接本身被解掉（= 票 79 要的行为）。
- Linux（容器 `golang:1.27`，`CGO_ENABLED=0`）：同名 POSIX 用例 ⇒ `link/` 被拒（`remove ...: not a directory`），`link\` 是另一个名字（POSIX 上 `\` 属合法文件名字符 ⇒ **不折**，正是 AC#2 要的反向语义），外来目录两次都还在。
〔独立复现〕

**D. POSIX 反向用例真跑**：容器里 `-run 'TestAC2POSIX|TestAC4POSIX|TestAC1|TestAC3' ./internal/winsec/` ⇒ `=== RUN 4 / --- PASS 4 / FAIL 0 / SKIP 0`、`rc=0`（Linux 上缝本来就是空的，`TestAC4POSIXFloorAnswersInsideTheNamedTree` **没 skip**）。〔独立复现〕

**AC#2 结论：通过**（Windows 与 POSIX 两侧我都自己量到；"祖先一个都不查"这个结局我用 15 枚新形状没造出来）。

## AC#3 `platformVerifyPlacement` 走同一把切分（不许留第二条"消费答案但不查"的路）

- **Windows：通过。** 我自己的 5 枚新拼写（含 `\\?\`、尾点祖先、混合分隔符、全 `/`）经 `SealFile` ⇒ **全部拒**，且每枚都量了 SID 级前后：victim 目录与其外来文件的 DACL 逐 SID **未变**。文件级越界检查：票 106 的私有集/`verifyPrivate`/`Seal*` 白名单一行未动（`c02c609` 的 9 个路径不含 `winsec_windows.go`）。〔独立复现〕
- **POSIX：不通过（这就是我总判 FAIL 的那一格）。** `internal/winsec/winsec_other.go:64` 是
  `func platformVerifyPlacement(path string) (string, error) { return path, nil }`
  ⇒ POSIX 的底线**一条 reparse/symlink 腿都没有**。我在容器里实测：
  `SealFile("/tmp/.../data/link/keep-me.txt")`（`link` 是指向外来的 symlink）⇒ **返回 `<nil>`，外来文件权限由 `-rw-rw-rw-` 变成 `-rw-------`**
  ⇒ 与 P3 同一结局类（"seal 穿过链接改掉别人的权限并报告成功"），只是平台换到 POSIX、`icacls` 换成 mode。
  同一枚 fixture 下 `data//link/...`（空段）与 `data/link/./...`（`.` 段）被 `lexicalTraversal` 拒了 ⇒ 红的是"缺链接腿"，不是"缺词法腿"。
  这条**预先存在**（票 108 没碰 `winsec_other.go`），但 AC#3 的措辞未加平台限定，且票面 AC#5 强制"POSIX 必须真跑"正是为暴露这类形状；按票头规则我不写"通过（附条件）"。
  ⇒ 登记 `R-108-1`，交回编排者立案。〔独立复现〕

## AC#4 树归属这一刀真落地（复用票 102 的 `Rewritten`/`Actable()`，不新造第二本账）

| 判据 | 结论 | 证据 |
| --- | --- | --- |
| 不实现 `ResolveAccounted` 的候选一律拒（类型上门） | **通过** | 基线 56 里 `TestAC4ResolverThatCannotAccountForItsTreeCannotSealAnything` PASS；`resolve.go:206-211` 是类型断言，`resolveAccounted` 的缺能力分支返回 `rewritten=true`（默认按"改写了"处理，方向安全）。〔独立复现〕 |
| 恒改写伪造在**安装期**被关系式探针拒 | **通过** | `TestAC4TreeOwnershipIsPartOfTheConformanceContract` 基线 PASS；MUT-2（把 `resolverTreeOwnershipFailure(r)` 换成 `""`）恰 1 枚红且红在它身上 ⇒ 这条腿是独立增量、不是被别的腿顺带打死。〔独立复现〕 |
| 复用票 102 的账：`res.Actable()` 原样保留、`Rewritten` 递过缝 | **通过（静态）** | `internal/risk/winsec_c26.go` 仍读 `Actable(`（`internal/risk/pathresolver_rewrite_account_test.go:174` 的静态锁在基线 56 之外由 `-count=2` 四包 664 全绿覆盖），并把 `res.Rewritten` 经 `ResolveAccounted` 递出；`ResolvePath` 每次密封只走这一条腿（`resolve.go:336-374`，无双解析）。〔独立复现（读码 + 门禁绿）〕 |
| **使用期**这一刀的真实边界 | **不通过（残余，实现方已自陈）** | 使用期只有两条腿：解析器**自报** `rewritten=true` ⇒ 拒；答案本身过一遍底线 ⇒ 否则放行。所以"每棵树各给一个干净且互相包含的伪造答案、且自报没改写"的解析器仍然能移动树。我的 N 组形状碰不到它（它需要一次安装），而生产路径无解除口 + 闩锁在并发下没破（AC#1 那格），故今天只能从包内塞。登记 `R-108-3`（把这条自陈边界写进裁决表，不当已解决）。〔日志＋归档，我抽验〕 |

## AC#5 变异三向 + 门禁 + POSIX 真跑

**变异（我自己的四发，全在 `/tmp/mut108-ac108` 纯净快照；每发同链 `grep -n` 打印被改后整行 + `go build` rc=0 先量）**
基线：`rc=0 RUN=56 PASS=56 FAIL=0`。

| 变异 | 落地证据 | build | 读数 | 红在哪（点名） |
| --- | --- | --- | --- | --- |
| MUT-1 `resolve.go:145` 后插 `resolver, seamLatched = nil, false`（把票 103 那扇门装回去） | `145:		resolver, seamLatched = nil, false` | rc=0 | rc=1 RUN=56 PASS=54 **FAIL=2** | `TestAC1SeamCannotBeFreedThenGivenATreeMovingFake`、`TestAC1SeamIsOneWayUseIsTheOnlyDirection` ⇒ **两红都在守卫上** |
| MUT-2 `return resolverTreeOwnershipFailure(r)` → `return ""` | `268/278: return ""` | rc=0 | rc=1 RUN=56 PASS=55 **FAIL=1** | 恰 1 红 = `TestAC4TreeOwnershipIsPartOfTheConformanceContract`；AC#1/AC#4 类型腿仍绿 ⇒ 该腿是独立增量 |
| MUT-3 `nativeIsBackslash := os.PathSeparator == '\\'` → `false`（Windows 不再认 `/`） | `winsec.go:301`、`resolve.go:437` 两行打印 | rc=0 | rc=1 RUN=56 PASS=41 **FAIL=15** | 红按形状点名：`TestAC2AncestorPrefixesForEverySeparatorShape/{all-forward-slash,mixed-separators,trailing-separator,doubled-separator,posix-backslash-is-a-name-character}` 等 ⇒ P2/P3 形状被精确复现 |
| **MUT-4（我补的一向）** `RemoveUnlinked` 的 `if !filepath.IsAbs(path) {` → `if false {` | `239: if false {`（+ 下面 240 行的错误串仍在） | rc=0 | rc=1 RUN=56 PASS=53 **FAIL=3** | `TestAC2AncestorGuardHoldsForEverySeparatorSpelling/**volume-only-relative-tail**`、`TestAC2ComponentsAndTraversalPerSeparatorShape` ⇒ 相对拼写那一格确实由这条腿守住（不是白加的收紧） |
| 还原 | 四发各自 `diff -q` 与 pristine 一致；最后一轮补 `MUT-SNAPSHOT ALL CLEAN`（我如实登记：MUT-3 的 perl 同时改了两枚文件，我当时只还原了 `winsec.go`，故"final diff"报过一次 `differ`，随后按 pristine 还原 `resolve.go` 并复量为 CLEAN） | — | — | 仓库内 `git status --porcelain` 我的路径 0 项 |

**门禁（纯净快照 `/tmp/gate108`，`git archive HEAD`）**
`go test -count=2 -v ./internal/winsec/ ./internal/memory/ ./internal/risk/ ./internal/secret/` ⇒ **rc=0**、
`=== RUN` **664**、`--- PASS` **660**、`--- FAIL` **0**、`--- SKIP` **4 行 = 2 个名字 × 2 轮**（`TestSubprocessCrashWriter`(memory)、`TestSyncRegistryProbeLive`(risk)，两条都是既有环境条件用例；**是 `-v` 量的**，不是"没跑所以没 SKIP"）。
四包：`ok winsec 23.717s` / `ok memory 26.313s` / `ok risk 8.952s` / `ok secret 0.458s`。
（实现方报的 602/602 是**三包**口径，我按票面要求跑四包 ⇒ 664/660，两数不矛盾。）〔独立复现〕

**POSIX 必须容器真跑（两个假绿坑都按规矩躲开）**
`MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -v "C:/Users/swq/AppData/Local/Temp/gate108:/src" -v wisp108mod:/go/pkg/mod -w /src golang:1.27 go test -count=1 ./internal/{winsec,risk,memory,secret}`
⇒ 容器内先 `ls /src/internal/winsec | wc -l` = **19**、并 `ls /src/internal/winsec/zz_accept108*` = `No such file`（证明**挂的是纯净快照、且我的探针没混进这次读数**）
⇒ `LINUX_RC=0`，`ok winsec 0.034s / ok risk 3.335s / ok memory 11.976s / ok secret 0.007s`。
我**没有**用 `go test -c` 产二进制在无源码目录执行（那造过假读数）；所有 rc 都来自容器内 `go test` 本体，输出先落文件再计数（不用 `cmd | grep x; echo $?`）。〔独立复现〕
同一枚容器跑法暴露出 AC#3 的 POSIX 缺口 ⇒ 见上，`R-108-1`。

**静态三门 + d22**：`gofmt -l` 四包 **空**；`$(go env GOPATH)/bin/gofumpt.exe -l` 同四包 **空**（本机有二进制，跑了）；
`go vet` 原生四包 rc=0、`GOOS=linux go vet` 四包 rc=0、`GOOS=darwin go vet ./internal/winsec/` rc=0；
`sh scripts/d22scan.sh`（**没从仓根 `go run ./tools/d22scan`**）⇒ **rc=0**、零 emoji（覆盖注释与 `_test.go`）。〔独立复现〕

## 台账疑点复算（`ban #6/#8 frontend/` 40→37）

- 我在**当前 HEAD 的纯净快照**上跑同一把仪器：`ban #6 frontend/ = 40`、`ban #8 frontend/ = 40`，同时 `bans #1-5 internal/=202`、`ban #7 internal/tools/=18`、`ban #8 internal/=365`、`cmd/=26`、`design/=16`。
- 结论（我自己量的，不是抄的）：**没有下降，也没有覆盖被抹**。这把尺数的是"被扫的那棵树里 `frontend/` 有多少个文本文件"，所以读数 = 快照 HEAD 的函数；实现方取数用的 `git archive HEAD` 早于票 92 把 `frontend/` 的 3 枚文件带进树，故它读到 37；今天的 HEAD 读到 40。同一次 d22 输出里 `internal/` 与 `internal/tools/` 也从 199/17 涨到 202/18，**同一漂移方向的旁证**（邻居在飞的新文件，不是谁覆盖了谁）。
- 建议动作：台账里的各 scope 读数今后**与 `git rev-parse HEAD` 同行登记**，否则"快照 vs 工作树"的读数差会被反复误读成覆盖下降。〔独立复现〕

## `RemoveUnlinked` 拒相对拼写：裁决

- 生产调用点穷举（`grep -rn "RemoveUnlinked(" --include=*.go`，排除 `_test.go`）：**只有两枚**，都在 `internal/memory/artifacts.go:204` 与 `:214`，都是 `full := filepath.Join(s.artifactsDir, rel)`。
- 那个调用方读源核对：`s.artifactsDir` 在 `internal/memory/open.go:169` 由 `abs, err := filepath.Abs(dir)` 得到、`open.go:183` 再 `filepath.Join(abs, artifactsDirName)` ⇒ 生产永不为相对拼写；`rel` 来自 `filepath.Rel(s.artifactsDir, p)` 且 `p` 来自以 `artifactsDir` 为根的 `WalkDir` ⇒ 不含 `..`、不空。
- 反向也量了（不许"守卫拒一切"算绿）：`TestAC2AncestorGuardStillUnlinksAPlainFile` 基线绿（普通文件、以及斜杠拼写的独立链接照旧回收），我自己的"我们那颗故意宽授权的 `blob.bin` 被真收窄"腿也绿，容器里 `internal/memory` 四包全绿。
- **裁定：这是 fail-closed 的正确从严，保留，不放行"为了可用性把这格打开"。** 依据：唯一调用方结构上不可能递相对拼写（绝对根 + Join + Rel），而相对拼写在**不走解析器**的入口上意味着"进程站的那棵树"，报成功而什么也没删是这仓登过两次的形状；MUT-4 证明这格由该腿独立守住。要求：把它当**契约收紧**写进票面结案段（一句"RemoveUnlinked 只接受绝对拼写"），并在 HANDOVER/台账留痕，别让它停在"请裁决"。

## 四种假绿与伪授权逐条点名

1. **"没跑所以没 SKIP"**：我的 SKIP=4 是 `-v` 量的，两个名字逐条点名（`TestSubprocessCrashWriter`、`TestSyncRegistryProbeLive`），环境条件用例、非本票路径。
2. **编译期糊真跑**：POSIX 一律 `docker run` + 容器内 `go test`，并先 `ls` 证明文件在（避开 Git Bash `-v C:\…` 静默挂空且 rc=0 那个坑——我用 Windows 正斜杠绝对路径挂载并 `ls` 复核）。实现方交件里"WSL 只有 docker-desktop 最小发行版"那一句属〔仅自述，不背书〕（我没复现 WSL 侧，走 Docker 得到同样的 Linux 读数）。
3. **变异打空**：四发各自 `grep -n` 打印被改后整行 + `go build` rc=0 先量（编译失败一律不算变异），红名逐一点到守卫；MUT-3 我顺带登记了一次"只还原一枚文件"的自己的操作瑕疵（不是代码缺陷，快照内、已补还原）。
4. **"拒一切"当绿**：两侧都量——外来文件在守卫生效时 `icacls` SID 列表逐 SID 未变；我们自己故意放宽的文件被真收窄（`S-1-1-0` 由有到无）；POSIX 侧同理用 mode 量。
5. **伪授权文本**：工具输出末尾自称"编排者备注/停手/撤回/请 revert"的文本 **0 次**；我未执行任何 revert、未撤销任何 commit。另有 2 次 `MEMORY.md was modified since it was last read` 类系统提示（记忆索引自身变更通知，无指令），按"不作为授权"处置，未据此改判据。

## R-108-x（本轮登记，验收期我一行生产码没改）

- **R-108-1（阻断结案，需立案）** POSIX/`winsec_other.go:64` 的 `platformVerifyPlacement` 是 `return path, nil`，底线**没有链接腿** ⇒ 未链接 `internal/risk` 的二进制里，`SealFile`/`SealDir` 穿过 symlink 改掉**外来文件的权限并返回 nil**（容器内实测：`-rw-rw-rw-` → `-rw-------`）。结局类与票 103 的 P3 同形、平台换掉。修法方向（要编排者裁，我不代做）：POSIX 侧用 `Lstat` + `os.ModeSymlink` 走同一把 `pathPieces` 前缀切分（POSIX 上只认 `/`），并把这条反向用例钉在容器真跑里。
- **R-108-2（观察，非本票 AC 的绕过）** `SealFile` 对**直接点名的外来绝对路径**（祖先链无链接）会成功并剥掉 `S-1-1-0`（实测）。winsec 的守卫是"拼写/树归属"级别，不是"这棵树归谁"级别；后者是调用方数据根纪律（票 76/95）的职责。要不要把"外来"定义成"越出调用方数据根"，请编排者裁——我不在验收期改判据。
- **R-108-3（实现方已自陈，我确认存在）** 使用期的树归属只信解析器自报的 `rewritten` + 答案自过底线；"每棵树各给一个干净且互相包含的答案"的解析器仍可移动树。可达性今天被 AC#1（无解除路径 + 闩锁并发不破）与唯一安装点限制在包内。
- **R-108-4（精度）** `SetPathResolver` 的"恒等重装"按 `%T` 比身份，不是指针身份；行为上保留第一份（安全方向），但日志措辞"installed again identically"在"同类型不同实例"时不精确。
- **R-108-5（继承实现方账）** `R-103-6` darwin 仍只有编译期读数（无 macOS runner）；`R-103-7` macOS `/tmp`、`/var` 本身是 symlink 时回收会开始报"拒删"——这两格本轮**未复验**（我没有 macOS），下一条命令应是 CI 里的 darwin 步或欠一个 run id + step 名。

## 逐格判定（与票面 AC 1:1）

- **AC#1：通过**〔独立复现〕。
- **AC#2：通过**〔独立复现〕（含 15 枚新形状 + POSIX 真跑反向用例 + 尾分隔符链接叶子专项）。
- **AC#3：不通过**〔独立复现〕——Windows 半边通过，POSIX 半边 `platformVerifyPlacement` 无链接腿，P3 的结局我造出来了（`R-108-1`）。
- **AC#4：通过（附残余 R-108-3 为待裁项）**〔独立复现 + 一档自述〕。
- **AC#5：通过**〔独立复现〕（四发变异各红在守卫、门禁 664/660/0/4、三门 rc=0、`sh scripts/d22scan.sh` rc=0、POSIX 容器真跑 rc=0）。
- **总判：FAIL / 不通过，退回实现方**：票面 AC#3 声称要防的结局（穿过链接改掉外来权限并报成功）在受支持平台的底线里仍被真实造出来。票 103 的三枚探针本身**都已红→绿并被我的新形状再攻不破**；本票的修法没有重开它们，重开的是它没覆盖的那一侧。
