# 票 301 — 编排者独立基线 ＋ 一枚写在腿交件**之前**的 Δ 预测（14:3x）

**谁写的**：编排者本人（`301-r1` 在飞期间补位的那一发；写面＝本件一枚，⛔ 碰 `scripts/**`、⛔ 碰 `internal/**`、⛔ 任何 go 编译面）。
**为什么值得单独落一枚件**：`301-r1` 交回来的名册作差我不能按它的转述去核。这一枚件先把我自己量到的基线和**我会拿哪三行去对它的读数**写死，交件时它说什么都不影响判据已经长什么样。

---

## 0. 锚与射程（⛔ 引用本件任何一枚数之前先看这一节）

- HEAD＝`64ceaacd`；`scripts/portable-tests.sh` 的 blob＝`2ff02dd7476c8f623a4489932d0d0ff95480e667`，767 行，md5 `328eb3ead0545736a33e2c131687d60d`。
- ★同一枚 blob 在 `41329475` 和 `86d89478` 上也是同一枚（尺＝`git rev-parse <ref>:scripts/portable-tests.sh`，三行输出逐字相同）⇒ **`301-a2` 那三张表和我这一枚量的是同一份字节**，枚数分歧⛔ 来自文件，只能来自尺。
- 射程＝**blob**（`git show HEAD:<path>`、`git ls-tree -r --name-only HEAD`）。⛔ 工作树——`301-r1` 此刻正在写 `scripts/portable-tests.sh`，盘上那份随时是中间态（先例＝票 294 那一笔的定式："树会呼吸，写腿在飞时盘上现量自带歧义，锚点一律引 HEAD"）。
- `.github/workflows/ci.yml` 同样只读 HEAD blob。
- 另有一枚**仓外**射程＝`.scratch/ci-logs/`（两枚归档 CI 运行日志：run `37158259050`＝2026-10-03、run `37166458550`＝2026-10-04）。这一射程里的字节是**托管 windows runner 真跑出来的输出**，⛔ 是我推的。

---

## 1. `AC#1` 要动的两处落点（逐字基线）

- **`win_pin` 正文在 `:193-204`**：`:193` 是 `win_pin='`，`:194-203` 十枚导入路径（字典序），`:204` 是收尾那个 `'`。十枚逐枚读过＝`cmd/llmrecord`／`internal/ball`／`internal/config`／`internal/perm`／`internal/plugin`／`internal/proc`／`internal/projctx`／`internal/risk`／`internal/secret`／`internal/session`。⛔ audio（尺＝`grep -c audio` 对 `:193,204p` ⇒ 0）。
- **`windows)` 档清单在 `:248-254`**：`:248` = `windows)`、`:249` = `scope=(`、`:250-252` = 十枚路径三行、`:253` = `)`、`:254` = `pinned=$win_pin`。⛔ `./internal/audio/`（同一把尺 ⇒ 0）。
- 该插的位置（⛔ 判据要求的形，只是本文件按现有规矩给出的观察）：`win_pin` 里 `internal/audio` 按字典序落在第 2 枚（`cmd/llmrecord` 之后、`internal/ball` 之前）；档清单里 `./internal/audio/` 落点自由（GUARD C 比的是集合，`:380` 那把尺是 `grep -qxF`，⛔ 顺序）。
- 相关调用点（沿用 `301-a2` 已交的读数，我这一枚⛔ 重跑 yml）：`--scope=windows` 唯一调用点 `ci.yml:780`，job＝`test-windows`、`runs-on: windows-latest`（**托管**）。

## 2. ★ 新现量（任何一份腿报告里都⛔ 有这一条）：被 `-skip` 过滤掉的用例**⛔ 产 `--- SKIP` 行**

这一条决定 `AC#1` 落地那一发到底会不会因为 ledger 而红，所以它必须在翻勾之前有人量过。

- 机制位置：`portable-tests.sh:694` 把 `go test` 交给 `$strict`（`:92` = `tools/d22scan/runtests.sh`），而 `runtests.sh:98-104` 逐字是 `if [ "$skipped" -ne 0 ]; then ... SKIP is not a pass (ticket 71 AC#3) ... exit 1`。⇒ **只要有任何一枚 `^--- SKIP`，那一档必红。**
- 那么 `-skip` 会不会造出一枚 `^--- SKIP`？两发独立 CI 字节回答**会⛔**：
  - run `37158259050`（10-03）`test-windows` 的 windows 档那一步：`=== RUN=488  --- PASS=335  --- FAIL=12  --- SKIP=1`，而它的 `-skip` 模式串里**含 7 枚名字**，其中 `TestHelperProcess`（`./internal/proc/`，在那一档里）和 `TestSyncRegistryProbeLive`（`./internal/risk/`，也在那一档里）**都在档内**。如果 `-skip` 会产 `--- SKIP`，这一发至少该是 `SKIP=3`。实测 `SKIP=1`。
  - 那一枚唯一的 SKIP 的**名字与出处**也在日志里：`--- SKIP: TestSyncRedTeamRealOneDrive` ＋ `syncdirs_redteam_windows_test.go:220: no live sync root on this machine (detected roots: [])`——它⛔ 在 ledger 里，所以 `portable-tests.sh` 把它打进"unaccounted SKIP lines"那一节。**它是自己 `t.Skip` 出来的**，与 `-skip` 无关。
  - run `37166458550`（10-04）同一档：`=== RUN=492  --- PASS=339  --- FAIL=12  --- SKIP=1`——第二发，同一枚 1，同一个形状。
- ⇒ **结论（判据级）**：`-skip` 是**过滤器**，被过滤的用例连 `=== RUN` 都⛔ 打；只有用例自己 `t.Skip` 才产 `--- SKIP`、才把那一档打红。
- ⇒ **对 `AC#1` 的直接后果**：audio 进档后，ledger `:593` 那行（`TestLiveWasapiSmoke|./internal/audio/|windows|fixture|`）会让那一枚用例被**静默过滤**——⛔ 新增红、⛔ 新增 SKIP、⛔ 新增 `=== RUN`。这一枚"落地会不会把 `--scope=windows` 打红"的担心，到此为止，⛔ 需要再跑一发来排除。
- ⚠ 但同一枚结论⛔ 能推出"落地必绿"：它只消掉了 ledger 那一支，⛔ 消掉第 3 节那 32 枚。

## 3. ★ 更正票面那半句"⛔ '改一枚清单'不是零成本"（**记我**——那句话是我自己写的）

- 票面原句（`现量` 那一节末条）写的是：把 `./internal/audio/` 拉进 `windows)` 档，"会同时把该包**其余 windows-tagged 用例**（上面那 2 枚，以及票 300 正要交进来的那一枚）一起拉进 `test-windows` 那个 job 的分母"。
- **这句把代价写小了整整一枚粒度**：分档的单位是**包**（`scope=( ./internal/audio/ )` → `go test` 收的是包），进档之后进分母的⛔ 是"其余 tagged"，**是该包全部顶层用例 42 枚**。尺＝HEAD blob 逐枚 `grep -c '^func Test'`（rc=0）：

| 文件 | 顶层用例 | windows tag |
|---|---|---|
| `capturelevel_windows_test.go` | 1 | 有 |
| `hotplug_test.go` | 8 | 有 |
| `parse_wave_format_300_windows_test.go` | 1 | 有 |
| `captureopt_test.go` | 4 | 无 |
| `gate_test.go` | 5 | 无 |
| `level_test.go` | 10 | 无 |
| `resample_test.go` | 7 | 无 |
| `wavinjector_test.go` | 6 | 无 |
| **合计** | **42**（tagged 10／无 tag 32） | 3 枚文件带 tag |

- ⇒ 本票真正在赌的是那 **32 枚无 tag 的用例**：它们在 CI 上到今天为止**只在 ubuntu 的 `core` 档里被跑过**，⛔ 有任何一枚在 windows 平台上被求值过。票面那两句（以及 `301-a1` 的 AC#0 三张表）讨论的前置条件全部集中在 tagged 那 10 枚上，**射程窄了 3 倍**。
- 载入风险⛔ 是风险（尺＝8 枚测试文件的 import 块逐枚读，rc=0）：全纯 Go，⛔ cgo、⛔ sherpa、⛔ 设备类依赖，唯一的包外依赖＝`internal/observe`。⇒ `0xc0000135` 那形（`cmd/wisp` 缺 DLL 的旧例）在这里⛔ 适用。
- **计时**才是那一族的风险面，具名两枚先记下来：`TestPinnedThreadStable10s`（`hotplug_test.go:443`）——它只在 `-short` 下跳，而 `:694` 那条命令**从⛔ 传 `-short`**（尺＝`grep -n -- '-short' scripts/portable-tests.sh` ⇒ 0 命中），⇒ **落地一次＝那一档多 10 秒墙钟**，且它判的是 `LockOSThread` 独占性（托管 runner 上是共享机器）；`capturelevel_windows_test.go` 的 `waitForLevels(t, rec, frames, 5*time.Second)` 同族。⚠ 这两枚⛔ 是"会红"的判定，是"红的时候先看哪两枚"的路线图。

## 4. 现在写死、腿交件时我拿来对的三行 Δ 预测

判据射程＝本机 `bash scripts/portable-tests.sh --scope=windows` 改前／改后各一发（`AC#2` 那形）。⚠ 下面这三行是**预测**，⛔ 是读数；它们的作用是：腿的读数如果落在这里面，说明它量的是同一件事；落在外面，**先怀疑尺**。

1. **Δ `=== RUN` ＝ ＋41**（42 枚顶层，减掉被 `-skip` 静默过滤的 `TestLiveWasapiSmoke`）。⇒ 谁报 **+9 或 +10**，它量的就是"tagged 那一半"，⛔ 是这一档的真实分母；谁报 **＋42**，说明那次跑里 `-skip` 没生效（ledger 没被读进去＝GUARD C/D 之外的另一枚缺陷，具名上报）。
2. **Δ `--- SKIP` ＝ 0**（第 2 节两发 CI 字节）。⇒ 谁报 Δ SKIP ≥ 1 **且那枚名字来自 `internal/audio`**，那就是我这节被推翻了：那一枚红是 `runtests.sh:99-104` 的形状（"SKIP 是红"），⛔ 是 audio 的判据坏——归因⛔ 许写成"用例红"。
3. **Δ `--- PASS` ＝ 41 −（具名新增红枚数）**。⇒ 新增红必须**逐名＋连红句文案**一起报（先例口径：红名册⛔ 区分坏法种类）。⛔ 拿"枚数一致"当"没问题"。

## 5. 我这一枚⛔ 量到的（具名，⛔ 当已知引用）

- `--scope=windows` 在**本机**今天的四数（跑它＝go 编译面，此刻归 `301-r1`）。
- `--scope=census` 的真读数（`unclaimed-with-tests=` 到底几枚、totals 行原文）——这条欠账从 `301-a2` 起就挂在"归编排者"那一节，⛔ 腿跑。
- `--scope=core` 改前改后的颜色（同一枚原因）。
- 那 32 枚无 tag 用例在 windows 上的真颜色——只有第 4 节那一发跑完才知道，本件只给出"该往哪两枚看"。
- ⚠ `301-a2` 表①里 `internal/audio` 那行的 **total `_test.go`＝8**，我这一枚独立复认（同一把 `git ls-tree` 尺，8 枚文件、42 枚用例），两把尺对得上；它表里 tagged＝3 也对得上。**这一枚算我复跑掉了，⛔ 再标〔仅腿量〕。**

## Progress log

- [2026-10-10 14:3x +0800] agent=orch did=blob 层基线（HEAD 64ceaacd＝同一枚 blob 三号一致）＋新现量"`-skip` ⛔ 产 `--- SKIP`"（两发归档 CI 字节：RUN=488/492、SKIP=1、那枚 SKIP 的名字与 file:line 具名）＋更正票面"其余 tagged 进分母"＝真实分母 42 枚（tagged 10／无 tag 32，逐枚文件表）＋写死 Δ 预测三行（RUN +41／SKIP 0／PASS 41−具名红）next=核 `301-r1` 的改前改后名册；⛔ 零 push；⛔ 任何翻勾（翻勾归非实现者）
