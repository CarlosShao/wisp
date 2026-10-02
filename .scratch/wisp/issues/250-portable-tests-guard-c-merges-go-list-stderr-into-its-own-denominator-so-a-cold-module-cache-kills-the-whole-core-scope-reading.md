# 250 — `portable-tests.sh` 的 GUARD C 把 `go list` 的 **stderr 并进分母**：冷模块缓存一发就把 25 枚包数成 35 枚，于是**整段 core scope 的 1365 个 `=== RUN` 读数永久采不到**

- Status: **已立，待派**（10-02 08:5x，编排者立；来源＝只读取证腿 `ci-delta-1`，裁决件 `.scratch/wisp/probes/orchestrator/ci-delta-1/delta.md` §2 第 4 条＋§6 B 第 5 条；台账 `A512`）。
- ⚠ **这不是产品缺陷，是我方仪器的口径缺陷**，但它的代价比大多数产品缺陷大：**它一次吃掉整段 core scope 的全部测试读数**，并把 4 枚基线红伪装成"消失了/修好了"。
- ⚠ **形状上它与我 09 月裁过的那一族同源**：「一步红会吃掉后续步 ⇒ 那步的日志永久采不到」（见本仓既有账：`gofumpt` 那枚被跟踪的坏变体吃掉两枚 `go vet` 读数）。区别＝这次不是"后续步"，是**同一步内部在 `go test` 之前就死**，所以连失败用例名都没有。

## 现量（编排者 08:5x 自己跑，别信行号、自己复算）

| 事实 | 读数 | 尺 |
|---|---|---|
| **并流的就这一行** | `go list "${scope[@]}" >"$resolved" 2>&1` ＝ **`scripts/portable-tests.sh:261`** | 我 08:5x 现跑 `grep -n "" scripts/portable-tests.sh \| sed -n '252,275p'` 复认，逐字 |
| **stderr 真进了计数** | `:268` `grep -v '^$' "$resolved" \| sort -u` → `:270` `pkgcount=$(wc -l <"$resolved")` ⇒ **stdout+stderr 混在一个文件里被逐行计数** | 同上，`:268-270` 我逐行读到 |
| **GUARD C 拿这个数去比钉** | `:283` 起那几行打印 `resolved to a DIFFERENT package` / `Pinned: 25, resolved: 35.` | 我现跑 `grep -n "GUARD C" scripts/portable-tests.sh`＝命中 `:64`／`:249`／`:272-274`／`:283` |
| **CI 上的那发逐字** | `Pinned: 25, resolved: 35.`，后面 10 行逐字是 `> go: downloading github.com/dustin/go-humanize v1.0.1`、`… google/uuid v1.6.0`、`… pelletier/go-toml/v2 v2.2.4`、`… remyoudompheng/bigfft …`、`… golang.org/x/crypto v0.57.0`、`… x/sys v0.48.0`、`… modernc.org/libc v1.75.7`、`… mathutil v1.7.1`、`… memory v1.12.1`、`… sqlite v1.59.0` | 〔`ci-delta-1` 读数，run `36889094435` 作业 `test-core` 步骤 #7；日志切片在本仓 `.scratch/wisp/probes/orchestrator/ci-delta-1/`〕 |
| **包集合没变** | `git diff --stat 0589fd9c..8ae4c23e -- scripts/portable-tests.sh` **为空**，`-- internal/risk/` 也为空 | 〔同上腿 §3 第 5 行与 §6 B 第 5 条〕 |
| **代价** | 该步 3 秒内死守卫，`go test` **一行没执行** ⇒ 基线那发 1365 个 `=== RUN` 在新发整体缺失；4 枚基线红（`TestC21DesignTokensFourWayAgree`、`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`、`TestComposerContractTypesMatchFrontend`、`TestApprovalCardViewJSONKeysMatchFrontendTypes`）因此"看起来消失了" | 〔同上腿 §1.3 与 §2 第 4 条〕 |

**触发条件（为什么本机跑不出来）**：`go: downloading …` 只在**模块缓存是冷的**时打印。本机 `GOMODCACHE` 早已温热 ⇒ 这枚缺陷在本机**永远静默**；CI 的那台机器每次换 runner 实例就可能冷 ⇒ **同一枚 SHA 两次跑，可以一次绿一次把整段读数吃掉**。现成对照样本＝同 `headSha 8ae4c23e` 的 schedule 触发 run `36940536372`（`ci-delta-1` 按票面没并进差集，但它就是"同码二次采样"）。

## 要建什么（一块石头：分母只能由"看起来像 import path 的行"构成）

- [ ] **AC#1 分流**：`go list` 的 **stderr 不许进分母**。判据＝`:261` 那行的 stdout 与 stderr 落到**不同**去处，且 `pkgcount` 只数 stdout；⛔ **`go list` 真失败时那一步仍须 `exit 1` 并把两路都打出来**（今天 `:262-266` 这条失败支不许被削掉——那枚是本票的**反面**，不是要放松的东西）。
  **正控**＝造一发"stderr 有内容、stdout 完全正确"的场景（最省形＝`GOMODCACHE` 指到空目录再跑 `--scope=core`），要求 `resolved` 枚数**仍等于钉住的那枚数**、守卫**不响**。
- [ ] **AC#2 能力形自证（不是词表）**：分母里的每一行必须**过一道形状检查**——不是把 `go: downloading` 这个词组 grep 掉（词面型尺的漏计：Go 换一句进度文案就复活），而是**"不含空格、含 `/` 或以模块路径形状出现"**这一类能力判据；任一行不过 ⇒ **响亮拒绝并退出非 0**，⛔ 不许"丢掉不像包名的行然后继续"（那是把守卫改成少测东西，正是 `:272-274` 那段注释立起来要拦的形状）。
  **正控**＝往 stdout 里种一行带空格的伪包名 ⇒ 指名用例必须红。
- [ ] **AC#3 常驻判据（不能只在 CI 冷缓存时才有牙）**：本票的尺必须**在本机就能判红/判绿**——用可控的冷缓存载具（临时 `GOMODCACHE`／或把 `go list` 包一层能注入进度行的载具），⛔ 不许写成"只有 CI 上才测得出"然后挂〔仅本机可量〕。判据＝**改前该载具必红、改后必绿**，两形读数都要落文件。
- [ ] **AC#4 不许顺带动别的守卫**：GUARD A／GUARD B／空 scope 硬退出（`:254-256` 的 `exit 2`）三形一字不动，各自配一枚"种 X 必响"的正控读数。⚠ 本仓有过一次"改一处把另一处的硬退出拆软"的事故，这一格就是把那件事写成判据。

## 禁区

- ⛔ **阈值 / golden / `internal/observe/thresholds.go` 一字节不动**；`docs/PLAN.md`、`docs/specs/**`、`docs/BUILD.md`、`docs/SLO.md`、`tools/d22scan/allowlist.txt` 一字不动。
- ⛔ 三枚冻结测试件一字不动：`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`。本票射程＝`scripts/portable-tests.sh` ＋它自己的台件/用例，**不碰 `internal/panel`**。
- ⛔ `frontend/**`／`design/**` 零读零写零转述（owner 已把界面侧委托给外部 agent）。
- ⛔ **不许为了变绿放宽任何断言**；⛔ 不许把 `TestPanelHostRealWindowHopAndLifecycle` 那类"只在真桌面"的用例塞进 `-skip` 蒙混（那是把守卫改成少测东西）。
- ⛔ 不许动 `tools/d22scan/runtests.sh:98-102` 那句「SKIP is not a pass (ticket 71 AC#3)」——本票与它无关，且那枚是**既有承重墙**。

## 排程与串行（共享工作树）

- 写面＝`scripts/portable-tests.sh`（**此刻 `git status --porcelain scripts/`＝空**，编排者 10-02 08:5x 现跑）＋新增台件目录。与在飞的 `248-r1`（`internal/panel`／`internal/config`／`cmd/wisp`）、`33-r8`（`internal/ball`）**包级不重叠 ⇒ 可并行**。
- ⚠ **本腿禁跑整树 `go test ./...`**（毫秒级计时用例正被另两枚腿读，负载会洗它们的读数）。要跑只许：本票的脚本级载具＋最多单包 `-count=1`。
- 落地后腿交：改前/改后两形读数、四格正控的确切尺、以及**它推翻编排者题面哪一句**（我上面的行号全是待验断言）。

## 归口与它挡住的事

- 这枚修好之前，**"core scope 在 CI 上到底几枚红"这一格无人能答**（包括那 4 枚"消失"的前端契约红到底是变绿还是被跳过——本仓既有定式：**判"不再红"先分清变绿还是被跳过**）。
- `ci-delta-1` §6 B 另外四条（常驻面文案出口、`TestPanelHostRealWindowHopAndLifecycle` 的自述与 CI 分母不一致、`TestResolvePerCallBudget` 负载账、`TestComposedGateBlocksAWriteForTwoSeconds` 的 300s 账）**不归本票**，逐条已具名登记于台账 `A512`，各归票 246／票 33／票 119·128 那一族／票 224·C18 那一族。
