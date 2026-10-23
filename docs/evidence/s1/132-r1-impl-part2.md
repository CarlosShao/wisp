# 票 132 落地腿 `132-r1` —— 门禁、红名册、提交、诚实边界（`132-r1-impl.md` 的续件）

> 前半 §1–§5（改前三个读数／落地逐行＋为什么那一跳产品会走／产品面逐字 DACL／新用例逐名／反向判据三读数）在
> `132-r1-impl.md`。本件＝§6–§9。拆件只因为交付帽 160 行（整件 189 行），⛔ 内容删减。

## §6 门禁 rc ＋ 红名册（逐名，取主检出；新增红 **0**）

| 门 | rc | 读数 |
|---|---|---|
| `go build ./...` | **0** | — |
| `go vet ./internal/observe/ ./internal/winsec/` | **0** | — |
| `GOOS=linux go vet ./internal/observe/` | **0** | POSIX 半边（`logprivate_132_other_test.go`）真编译 |
| `go test -count=2 -v ./internal/observe/` | **0** | **154 PASS / 0 SKIP / 0 FAIL**（77×2；六枚新例各跑两遍皆 PASS） |
| `sh scripts/d22scan.sh` | **0** | 正向控制跑过；`examined 290 production Go files`（与 `find` 现数的 290 **对上**） |
| `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`，最长 185、runner 229 |
| `go test ./cmd/wisp/` | **1** | 红 2 枚（下） |
| `go test ./internal/winsec/` | **1** | 红 1–2 枚（下） |

**已知红（逐名）＋ 为什么不是本腿引入**：

1. `cmd/wisp`：`TestAC2RealProcessRefusesOnEveryLegWithoutAppData128/resident`、`TestAC14GoSideEvalPushReachesThePage`。
   **正控**：`git archive HEAD` 导出**纯净 HEAD 树**（无本腿任何改动）到仓外，同两枚**同样红、同样名** ⇒ **新增红 0**。
   机理（第一枚）：`internal/proc/envfork.go:243` 的文案是中文「用户配置目录不可得」，而用例 marker 是英文 `"user config dir"`。
2. `internal/winsec`：`TestAC1WindowsPrivateFile…342r1`／`TestAC2WindowsLeafIdentityReCheck…342r1`（各 400 次尝试的**竞态仪器**，本会话两次跑红名单不同 ⇒ 本就抖）。本腿 `git show --name-only` 里**无任何 winsec 文件**。

---

## §7 格式两尺

- **甲尺（判据）**：`git archive b22225513` 导出到仓外再跑 `/d/work/base/gopath/bin/gofumpt.exe -l .` ⇒ 20 行，
  **20/20 全在 `.scratch/wisp/probes/**`（历史突变残留）**，`internal|cmd` 命中 **0**；我 5 枚文件 **0 命中**。
  正控：在导出树里放一枚故意不格式化的 `internal/observe/zz132_control.go` ⇒ 同尺**逐字点名** `internal\observe\zz132_control.go`
  （★反斜杠形状 ⇒ 判归属必须 `[\\/]`）。
- **乙尺（诊断，工作树）**：`gofumpt -l internal/observe/` ⇒ 空；`gofmt -l` ⇒ 空；5 枚文件 `tr -cd "\r" | wc -c` ⇒ **全 0**。

---

## §8 提交

| 提交 | `git show --name-only` | 越界 |
|---|---|---|
| `b22225513` | `internal/observe/logging.go` ＋ `logprivate_132_test.go`／`logprivate_132_windows_test.go`／`logprivate_132_other_test.go`／`sealcalls_132_test.go` | **0**（全在 `internal/observe/**`；无 `cmd/**`、无 docs、无 scripts） |

**正控（同尺打一笔含 `.go` 的真提交，否则"越界 0"可能是尺瞎）**：同一条 `git show --name-only` 在
`6e19f4a43`（票 196 落地腿，**别人的**一笔）上列出 **7** 枚 `.go`、在 `5ffa553e3` 上列出 **2** 枚 ⇒ 尺看得见 `.go`。
本笔 `b22225513` 同尺列出 **5** 枚 `.go`，逐枚都在 `internal/observe/`，**越界 0**。
（本腿第一版把 `96f1af9c9` 当正控，现量它 `.go` **0 枚** ⇒ 那是个空正控，已换掉并留痕于此。）

---

## §9 诚实边界（未证／未做，逐条具名）

1. **`SealDir`（导出函数）生产调用者仍是 0**。落地用的是 `winsec.PrivateDirAll`——它在**创建每一层之前**就封
   （`MkdirAll`+`SealDir` 做不到，正是 AC#2「创建时就」的语义），且 `privateDirAll` 内部调的就是同一枚 `sealDir` 原语。
   ⇒ 票面标题那把尺（`SealDir(` 计数）今天仍读 0，**能封的那个能力**在生产里第一次有调用者。这一处请裁决者按"能力"而非"函数名"读。
2. **`cmd/wisp/logsink.go:45-49` 的 `INTERIM (ticket 117 AC#3)` 注释现已过期**（它写"nothing here seals or un-seals the log file"）。
   该文件在 `cmd/**`＝本腿**禁写**，⛔ 未动；请派 `cmd/wisp` 那一枚写腿收口。
3. 票面 §地界 写的是 `cmd/wisp/logsink.go`＋`internal/winsec/winsec.go`＋`internal/memory/**`，而票面自己 `:88` 又判
   `internal/observe/logging.go:72` **不在本票地界**（归票 238）。本腿按派单写面落在 `internal/observe/**` ⇒ **两处地界文字需编排者对齐**。
4. **POSIX 侧**：文件 mode 由创建那一跳就 0600，"封文件"与"建文件"同 mode ⇒ 该平台**没有**"文件自带显式 ACE"这枚逐枚判别子；
   POSIX 上靠 §4 第 6 枚（调用枚举门）兜住"摘调用"，本件不假装 mode 能区分。
5. **数据根只在"它自己由本调用创建"时才被封**（§3 形状 A）；根**已存在且宽**时本腿**不碰它**（形状 B：根仍 `(I)` 宽，
   `PrivateDirAll` 明文契约"绝不碰 `path` 以上的祖先"）。⇒ 老机器上"宽的数据根"这一族**仍在**，按票面 `:88` 归票 238；
   ⛔ 本件不把它写成已修。
6. 反向对照是**组 SID 口径**（`S-1-5-11`／`S-1-5-32-545` 两枚 `findsid` 由 Found 变 Not found），
   ⛔ 不是"以第二枚真实账户登录试读"——本机无第二账户，不假称做过。
7. `wisp doctor` 在那棵封过的根上**整包 rc=1**：FAIL 项（sherpa-onnx／onnxruntime DLL 共置与 pin）**与改前逐行相同**，
   因为本腿把二进制放在仓外 `bin\` 而不是 `build\`。⇒ 它**不是**封条引起的，但本件不用它当"零影响"的正控；
   零影响的正控是 `slo` 连发 rc=0 与 `[PASS] data dir writable (test)`。
