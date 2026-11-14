# 票 441 — POSIX 封印窗口仪器在「弃权」时变红：普查、修、两侧读数

- 腿＝票 441 实现腿（委派）；2026-11-14 14:4x–15:3x +08；起手 HEAD `b062d52a7`（写腿期间共享树被别的腿推进到 `a46352fd3`，本件所有读数都在**自己的工作树**上现取，并逐条标了取数时的容器/命令）。
- 写面：`internal/winsec/posix_seal_window_120_other_test.go`（唯一一枚 `.go`，**测试文件**）＋
  `docs/evidence/s1/441-instrument-abstention.md`（本件）＋ `.scratch/wisp/probes/441/**`（新）＋ 票 441 尾一条 Progress bullet。
- ⛔ 未动：`internal/winsec` 的任何非测试文件（产码 0 字节，见 §6）、阈值／golden／契约、其它腿的未提交改动、`.github/**`、`scripts/**`、`frontend/**`。⛔ push。

## 一、AC#1 复现：在受控争用下让窗口「不开」，读出现有那条红句

### 1.1 CI 那一发（判据来源，现取不是抄）

`gh run view 43565131755 --log-failed` → `job=test-core`、`step 7 Portable package tests (core scope)`、
`head 5fb21317`（该 commit 只改 `docs/reports/pending-and-issues.md` 一行）。逐字：

```
=== RUN   TestPOSIX120SealWindowLeafSwapIsClosed
    posix_seal_window_120_other_test.go:365: AC#1 sample 1/4 shape=leaf-swap: attempts=500 hit-foreign=0 of-which-reported-nil=0 refused-at-check=468 of-which-refused-at-act=2 sealed-own-object=0 other=32
    posix_seal_window_120_other_test.go:365: AC#1 sample 2/4 shape=leaf-swap: attempts=500 hit-foreign=0 of-which-reported-nil=0 refused-at-check=466 of-which-refused-at-act=0 sealed-own-object=0 other=34
    posix_seal_window_120_other_test.go:365: AC#1 sample 3/4 shape=leaf-swap: attempts=500 hit-foreign=0 of-which-reported-nil=0 refused-at-check=471 of-which-refused-at-act=0 sealed-own-object=0 other=29
    posix_seal_window_120_other_test.go:365: AC#1 sample 4/4 shape=leaf-swap: attempts=500 hit-foreign=0 of-which-reported-nil=0 refused-at-check=481 of-which-refused-at-act=0 sealed-own-object=0 other=19
    posix_seal_window_120_other_test.go:370: AC#1 total shape=leaf-swap: attempts=2000 hit-foreign=0 of-which-reported-nil=0 refused-at-check=1886 of-which-refused-at-act=2 sealed-own-object=0 other=114
    posix_seal_window_120_other_test.go:370: AC#1 RED instrument is dead shape=leaf-swap: refused=1886, sealed-own=0 - a zero window reading here proves nothing
--- FAIL: TestPOSIX120SealWindowLeafSwapIsClosed (0.04s)
```

★同一发里**另两枚形状是绿的**，逐字 `ancestor-swap-relative-sibling-link: attempts=2000 ... refused-at-check=625 of-which-refused-at-act=50 sealed-own-object=57 other=1318` ⇒ 红的不是产码，是**同一枚仪器的一个计数臂**。

★关键读数：`refused=1886 > 0`（**窗口到达过**：密封 2000 次里有 1886 次被 floor/gate 拒在换过的状态上），
而 `sealed-own=0`（**静侧一次也没被 raced 循环看到**）。旧断言把这两个问句捆成一句
`totals.refused == 0 || totals.sealedOwn == 0`，于是「静侧被调度饿死」被报成「仪器死了」。

### 1.2 本机台面与平台射程（先说清）

- 该文件首行 `//go:build !windows`：本机 Windows 档**结构性不编译**（与票 346 同一条射程）。所有 `!windows` 读数都在
  `golang:1.27-bookworm` 容器里取（与 346 同镜像）。容器内 `/tmp` 用 tmpfs 是为了提速（见 1.3 的机理）。
- 台件：`.scratch/wisp/probes/441/repro-leafswap.sh`（只跑
  `go test ./internal/winsec -run '^TestPOSIX120SealWindowLeafSwapIsClosed$'`，`-count` / `-cpu` 是运行旋钮，**不改断言**）。

### 1.3 复现：同一臂（refused>0, sealed-own=0），N 与命中率

| # | 台面 | 分母 | N | 红 | 读数 |
|---|---|---|---|---|---|
| A | overlay `/tmp`，GOMAXPROCS=1 | 2000 | 10 | **0/10** | sealed-own 588…1407（sealed-own 不是 0）→ `repro-cpu1-count10.log` |
| B | overlay `/tmp`，GOMAXPROCS=4 | 2000 | 10 | **0/10** | sealed-own 675…1021 |
| C | tmpfs `/tmp`，GOMAXPROCS=2 | 2000 | 10 | **0/10** | sealed-own **43…90**（最低 2.2%） |
| D | tmpfs `/tmp`，GOMAXPROCS=4 | 2000 | 10 | **0/10** | sealed-own **24…221**（最低 **24/2000 = 1.2%**） |
| E | tmpfs `/tmp`，GOMAXPROCS=8 / 12 | 2000 | 6+6 | **0/12** | sealed-own 62…80 / 47…135 |
| F | tmpfs `/tmp` + `WISP_120_ATTEMPTS=1`（该文件自带的探针旋钮） | **4** | **20** | **15/20 = 75%** | **15 发全为 `refused=3..4, sealed-own=0`**，红句逐字同 CI |

F 的 15 条红句（逐条同形，N 与命中率如上）：

```
posix_seal_window_120_other_test.go:370: AC#1 RED instrument is dead shape=leaf-swap: refused=3, sealed-own=0 - a zero window reading here proves nothing
...（refused=3 共 11 发，refused=4 共 4 发；15/20）
```

**为什么 F 能把全分母的 0 变成可复现的 75%，而 A–E 不能**（诚实交代，不粉饰）：
`sealed-own>0` 是「racd 循环里恰好采到叶子是**真文件**的那一刻」。这个概率随**密封 syscall 相对 racer 翻转周期的速度**单调下降：
overlay 2.2%–70%（A/B）→ tmpfs 1.2%–4.5%（C/D）→ CI runner 更快，读数落到 **0/2000**。
本台面比 CI runner 慢 ~2.5x（同一发：本机 0.20–0.42s/2000 次，CI 0.04s/2000 次），所以**全分母下本机到不了 0**（最低 1.2%），
而**把分母缩到 4 次抽样**（`WISP_120_ATTEMPTS=1`，文件头自己写着这是「cheaper probe run」的旋钮）就落回同一臂：75% 红。
⇒ 这不是「跑一百次没红」，是**同一句红句、同一个计数臂、分母 2000→4 的可复现读数**；全分母下 1.2% 的残余静侧观测也一并列出，读者可自己判边界。

## 二、AC#2 弃权 vs 真红：两个问句、两个载体

### 2.1 改了什么（一枚文件，测试侧）

旧：

```go
if totals.refused == 0 || totals.sealedOwn == 0 {
    t.Errorf("AC#1 RED instrument is dead shape=%s: refused=%d, sealed-own=%d - ...")
}
```

新（`judgeRun120` 拆成三条臂 + 一枚**不参与竞速**的对照）：

| 问句 | 载体 | 结果 |
|---|---|---|
| 窗口**到达过**吗？ | **raced** 行：`totals.refused > 0` | 0 ⇒ **真红**（逐字点名 attempts/samples/attempts-per-sample/sealed-own）。删 racer、floor 死掉仍会红——这臂是牙，不许松 |
| 密封**还封得动**自己的对象吗？ | **un-raced** 对照 `sealWindowQuietControl120`：停 racer → 用例自己重植真对象 → 调一次 `SealFile` → 要 `nil` 且证人未被碰 | 失败/命中 ⇒ **真红**（逐字点名 attempts/refused/refused-at-act/sealed-own/hit-foreign/first-other）。不再靠运气，也抓得住「密封死了」 |
| raced 循环**没看到静侧**呢？ | `totals.refused > 0 && totals.sealedOwn == 0` ⇒ 具名、带数字的一行 | `AC#1 ABSTAIN(441:<shape>:raced-quiet) attempts=… samples=… attempts-per-sample=… threshold=sealed-own>0 raced…`——**不硬红**，也**不是从缺席推绿**：静侧由对照实测（sealed-own=1） |

另保留：`attempts == 0` ⇒ 真红；`flipErrs > 0` ⇒ 真红（换态本身错＝真缺陷，不是调度）；
`refused>0 && refusedAct==0` 的「这一发没考到 identity gate」照旧同屏打读数（`t.Logf`，与 346 修订后的口径一致）。

### 2.2 修后读数（同一批台面，逐条对得上）

| 台面 | N | 结果 | 读数 |
|---|---|---|---|
| tmpfs，全分母 2000 | 10 | **10/10 绿，0 弃权** | 每发都打出 `un-raced quiet control ... sealed-own-object=1`；`repro-postfix-quiet-count10.log` |
| tmpfs + `WISP_120_ATTEMPTS=1`（= §1.3 F 的同一批抽样） | 20 | **15 绿 / 5 红** | 15 绿里 **7 发**打出 `ABSTAIN(441:leaf-swap:raced-quiet)` 并带数字；5 红全是**另一臂**（`refused=0, sealed-own=4`：racer 一次也没把换态递给密封）→ `repro-postfix-att1-count20.log` |

5 条红的逐字（这臂**故意保留红**，理由见 2.4）：

```
AC#1 RED instrument is dead shape=leaf-swap: across 4 raced attempts (samples=4, attempts-per-sample=1) the floor/gate refused 0 and sealed-own=4, so the racer never presented the swapped state and a zero window reading here proves nothing about a window that was never reached
```

7 条弃权行之一（逐字）：

```
AC#1 ABSTAIN(441:leaf-swap:raced-quiet) attempts=4 samples=4 attempts-per-sample=1 refused-at-check=4 of-which-refused-at-act=0 sealed-own=0 other=0 - the raced loop reached the swapped state (refused>0) but the quiet state was starved by scheduling, so this run cannot attest the quiet side from the race; the threshold that carries it is the un-raced quiet control above (sealed-own=1), not a green from absence
```

### 2.3 为什么不是 `t.Skip`（AC#2 的「具名跳过」在这里的真实边界）

`t.Skip` 在本档是**结构性红**，不是修：CI 的这一步走 `scripts/portable-tests.sh` → `tools/d22scan/runtests.sh`，
后者 `:24` 逐字写着 `2. any \`--- SKIP\` -> SKIP is NOT a pass`，`:234` 就是失败分支。要在 core 档里让 SKIP 合法，只能：
①改 `scripts/portable-tests.sh` 的 ledger（**不在本票写面**，且那正是「把 skip 长期停放」的禁手）；
②或接受「弃权仍把 CI 拉红」——那等于没修票面那个 flake。
⇒ 本腿采用票面同一句话的后半句「也不许静默绿」的实现：**具名的弃权行 + 数字 + 可 grep 的 token**，
并把静侧改成**确定性实测**（un-raced 对照），所以它不是「弃权变绿」，是「静侧不再靠抽签」。
若编排者坚持要 Go 层的 `t.Skip`，那需要一张新票去改 `scripts/**` 的 ledger，并接受 SKIP=红。

### 2.4 剩下 5 发为什么必须红（不是没修完）

`refused==0` ＝ racer 在这 4 次抽样里**一次也没把换态递给密封** ⇒ 这一发对「窗口被测到过」零证据。
把它变绿／跳过，删掉 racer 的那一发就会**假绿**——正是票 322 AC#1 逐字禁的「拿 0 命中当凭据」。
这一臂在**提交分母**（2000 次）下几乎不可达（A–E 共 48 发、sealed-own 全 >0，refused 从未为 0）；
只有把分母缩到 4 次抽样时才会出现。⇒ 票面那个 flake（refused>0 + sealed-own=0）被修掉，这臂保留作牙。

## 三、AC#3 牙还在：仓外导出树里的两发突变

导出树 `D:\tmp\wisp\441\export`＝`git clone --shared`（`a46352fd3`）＋把**修后**的测试文件拷进去；
仓库里的产码字节一个没动（见 §6）。两发突变都在导出树里做、都在容器里跑。

| # | 突变（仅导出树） | 产码 SHA256 | 读数 |
|---|---|---|---|
| 1 | **拆护栏**：`sameObjectPinned` 函数体 → `return nil`（identity gate 整枚摘掉） | `0AAC3650B2A08D123D66541D84FCB544DCB721451543B26EB2C5BF91B00643A0` | **真红**：`TestPOSIX120SealWindowAncestorRelativeLinkIsClosed` FAIL，逐字 `AC#1 RED sample 1/4 shape=ancestor-swap-relative-sibling-link: the seal landed on the foreign tree 23 time(s), 23 of them while returning nil`（4 个样本 23/13/15/14 次落外）；同发里 leaf / absolute 两形仍 PASS |
| 2 | **把密封做死**：`applyDescriptorPOSIX` 顶部直接 `return placementRefusal(...)`（每次必拒） | `3520D222965BBBA40B8B57DE84763DAF34CC4FBA11D93B8FED46BA07EA77AE7C` | **真红**：`AC#1 RED instrument is dead shape=leaf-swap: the un-raced quiet control could not seal the tree's own object with the racer stopped (attempts=1 refused-at-check=1 of-which-refused-at-act=1 sealed-own=0 hit-foreign=0 ...)`；同屏还打出 `ABSTAIN(441:...)`（raced 侧 refused=2000/sealed-own=0）⇒ **同一发里「弃权行」与「真红」并存**，证明弃权不会把死密封洗绿 |
| 还原 | 拷回 `winsec_other.go` | **`19F3C211055011CC699A8D8027EF49FF543EDE0DA5C863F58AAAAD21A640D3E2`（与突变前同串）** | 三枚全 PASS：`ok github.com/CarlosShao/wisp/internal/winsec 1.913s` |

原始 log：`probes/441/logs/03-mutation-guard-removed.log`、`04-mutation-dead-seal.log`、`05-restored-green.log`。

## 四、AC#4 同类普查（先普查，它决定了本腿改几处）

尺（可复跑）：`sh .scratch/wisp/probes/441/census.sh`（`git grep -n -E 'RED instrument|instrument is dead|RED INSTRUMENT' -- internal`，
逐行取**上面那一行守卫**判臂）。现量＝**37 行 / 13 个文件**（本腿修完后的树；346 的台账是 26 行/11 文件，且行号会腐烂，⛔ 抄旧表）。

下面这张表里 **「到达数」臂**＝`provedX == 0 || liveness == 0` 这类「窗口没到达就红」的自检，是**唯一**会在 CI 争用下弃权变红的一族。
其余行是**世界造不出来**（`Fatalf`/`if err != nil`）或**注入式确定性 shot**，不属这一族。

| `file:line` | 守卫 | 档 | 跑在哪 | 会不会弃权变红 | 本腿 |
|---|---|---|---|---|---|
| `posix_seal_window_120_other_test.go:312,317,324` | `attempts==0` / `refused==0` / `control.sealedOwn!=1` | `!windows` | test-core | **票面这一枚**；`refused==0` 臂保留红，静侧改由 un-raced 对照承担 | **已修** |
| `openleaf_322_other_test.go:206` | `liveness == 0 \|\| flips == 0` | `!windows` | **test-core** | **会**（346 的 log 实录：`answered 0 and the racer landed 6 flips`） | ⛔ 未修，具名交回（见下） |
| `ancestoropen_342_unix_test.go:117` | `provedLinkCycles == 0 \|\| pinnedParentLeg == 0` | `unix` | **test-core** | **会**（racer 被饿死时两数可为 0） | ⛔ 未修 |
| `ancestorunlink_348_unix_test.go:153` | `provedLinkCycles == 0 \|\| liveness == 0` | `unix` | **test-core** | **会** | ⛔ 未修 |
| `ancestoropen_342_windows_test.go:56` | `provedReparse == 0 \|\| formRefusals == 0` | `windows` | test-windows | **会** | ⛔ 未修 |
| `ancestorunlink_348_windows_test.go:111` | `provedReparse == 0 \|\| formRefusals == 0` | `windows` | test-windows | **会** | ⛔ 未修 |
| `deleteseam_348r2_windows_test.go:131` | `landed == 0 && refused == 0` | `windows` | test-windows | **会**（翻转一次都没被问到） | ⛔ 未修 |
| `deleteseam_348r2_windows_test.go:141` | `landed != calls` | `windows` | test-windows | **会**（逐发都要求翻转站稳，调度可饿死它） | ⛔ 未修 |
| `deleteseam_348r2_windows_test.go:126` | `calls != int32(winAncAttempts342r1)` | `windows` | test-windows | **不会**（`calls` 是台面自己数的注入次数） | — |
| `leafrace_342r2_windows_test.go:389,395` | `flipped == 0` / `recheck == 0` | `windows && wisp342probe` | **⛔ 不进默认构建** | 形同族，但 CI 不编译 | — |
| `sealhandle_privateroute_120_other_test.go:212` | `flipErrs > 0` | `!windows` | test-core | **不会**（换态报错＝真缺陷，不是饿死） | — |
| `ancestorunlink_348_unix_test.go:105,108,111` | `plantErr != ""` / `!planted` / `err == nil` | `unix` | test-core | **不会**（世界造不出来，确定性） | — |
| `ancestorunlink_348_windows_test.go:79,144,178,181` | `err == nil` / pin 失败 / `toFlip` 失败 / `!reparseAt` | `windows` | test-windows | **不会**（台面构造失败，确定性） | — |
| `ancestorwindow_342r2_windows_test.go:366,567,662,699` | `!w.sealRun` / 正控没毁掉它自己的证人 / 基线 rename 失败 | `windows` | test-windows | **不会**（正控/基线构造，确定性） | — |
| `leftstanding_354_windows_test.go:84,88,127,131,201` | 68 字节台面 / `!w.sealRun` / `kSwap != "absent"` / 拒绝哨兵 | `windows` | test-windows | **不会**（注入式确定性 shot + 台面） | — |
| `nametopup_353_windows_test.go:67,124,128,187,191` | 同上族（`entered/sealRun/kSwap/kind`） | `windows` | test-windows | **不会**（注入式确定性 shot + 台面） | — |
| `internal/agent/compress_trace_test.go:108` | `len(tc.all()) != before+1` | 无 tag | test-core | **不会**（自己写一条哨兵、要求 +1，确定性探针） | — |

★**为什么另外 7 枚「会」的不在本腿一并改**（这是普查给出的判断，不是省事）：

1. **票 322 那枚没有可用的确定性载体**，它自己的文件头写死了理由：`a leaf link planted before the call is refused by the placement floor before the open runs, so planting cannot reach this leg`
   ⇒ 那枚仪器**只能用 racer 到达**那个 leg，静侧／到达侧不存在「不竞速也能测」的替代。把它改成具名弃权
   ＝「没有任何独立证人时变绿」，正是本仓禁的假绿。要修它得先有新的注入面（=动产码，须另票＋人工批准）。
2. 其余 6 枚（`ancestoropen/unlink` 的 unix/windows、`deleteseam`）各自有自己的世界构造与 `pinned`/`reparse` 探针，
   每一枚的确定性载体都要单独设计并单独证明；本票 AC 只要「普查表＋这一枚」，把它们硬塞进这一腿就是**未经验证的改动**。
   ⇒ 具名交回：建议下一票按本件 §2.1 的三问句逐枚落，优先 `ancestoropen_342_unix` 与 `ancestorunlink_348_unix`（都在 test-core，都会红 CI）。
3. ★这一族**不在** ticket 346 的 `-done` 结论里：346 选的是「丙＝弃权照旧红但同屏给四数」，本票是它的后续（弃权不再硬红）。

## 五、AC#5 门禁四数

| 门 | 读数 | 命令 |
|---|---|---|
| `bash -n`（本腿碰了 shell：`probes/441/{census,repro-leafswap}.sh`） | **rc=0 / rc=0** | 容器内 `bash -n .scratch/wisp/probes/441/census.sh`、`.../repro-leafswap.sh` |
| `sh scripts/d22scan.sh`（clean） | **rc=0**：`d22scan: clean - no D22 ban violations`；step1 正控 `go test -count=1 ./...` rc=0 | 本机无 POSIX sh（WSL 无发行版）；**容器里跑该 wrapper 时 `tools/d22scan` 自检命中 `go test` 默认 10m 超时**（`FAIL .../tools/d22scan 600.331s`），故按 wrapper 的两条命令在**本机原生**逐条跑（`cd tools/d22scan && go test -count=1 ./...` → rc=0；`go run . -root <repo>` → rc=0）。口径同 wrapper，扫描集：internal/ 1010 Go 文件、cmd/ 280 |
| path-length budget | **VERDICT GREEN**（tracked paths=19952，over-budget=69 全在册，not in roster=0） | `sh scripts/check-path-length-budget.sh`（提交后复跑，见票尾） |
| 改动只碰写面 | `internal/winsec` 仅 1 枚 `*_test.go`；产码 0 字节 | `git status --porcelain internal/winsec`、`git diff --stat -- internal/winsec/*.go ':!*_test.go'`（空） |

## 六、边界自证（可核）

- **产码 0 字节**：`git diff --name-only -- internal/winsec ':!*_test.go'` ＝ 空；`winsec_other.go` 在**仓库里**的 SHA256 仍是
  `19F3C211055011CC699A8D8027EF49FF543EDE0DA5C863F58AAAAD21A640D3E2`（＝突变/还原那一串的pristine值）。
- **未碰**阈值／golden／`thresholds.go`／契约／`go.mod`／`go.sum`／`.github/**`／`scripts/**`／`frontend/**`／票面正文（只追加票尾 bullet）。
- **本件的射程**：证「这一枚仪器的静侧不再靠抽签＋牙仍在（两发突变一红一绿）」；⛔ 证「CI 从此不会再红」
  （还有 §4 的 7 枚在册，且 CI 的其它红各有其主）。⛔ 把「本机全分母没红」读成「CI 绿」。
- ⛔ `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`；commit 用显式 pathspec；⛔ push。
