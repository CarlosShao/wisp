# 票 124 落地腿 124-r3 ---- cmd/wisp 一批：用例侧递根点接上 proc.SealableRoot

本件是**实现者**交件，NOT 裁决。裁决表由验收方另出。

## 0. 一句话

把 `cmd/wisp` 里"该转还没转"的 **37 枚 only-in-link** 用例的**递根点**接上票 119 那条纪律
（OS 给的答案先过 `proc.SealableRoot` 再递密封底线）：软链 TMPDIR 形状下本包
only-in-link 红名 **37 -> 0**、**新增红 0**（名册集合 `comm`，NOT 枚数相等）。

> 本件因 160 行帽拆两份：本件 = sec 0..4；`124-r3-impl-part2.md` = sec 5..9
> （三读数与 sha256、名册出口判据、两把分母、门禁 rc、红名册、提交与正控）。

## 1. 开工前三闸

(1) **写面脏态**：`git status --porcelain -- cmd/wisp` = **空**（开工首读 HEAD `9697a9bd`）。
全树 1250+ 条脏项全在 `.scratch/**` 与未跟踪件，`cmd/wisp/**` 内零枚别人未提交改动 => 可开工。

(2) **票面前提新鲜度**：票面那条"`winsec.` 与 `TempDir` 同行"的静态尺我**一枚都没当分母**
（裁腿 `124-v1` sec 一.5 已证它恒等于已修态、恒 20）。分母全部走终判据：纯净快照 + 两形各一枚新容器
+ `-count=1 -v` + 名册 `comm`。

(3) **两把分母独立复跑**：见 sec 3。我一手量到 **37**（票 124 族）与 **68+9**（DPAPI 族），
两把都与裁腿 `124-v1` 的名册**逐名 set-equal**（双向差集各 0 行），NOT 抄任何数。

## 2. 台面

纯净快照 `git archive HEAD | tar -x`（17503 文件）；容器 `golang:1.27` / `go1.27.1 linux/amd64` /
`GOPROXY=off` / `CGO_ENABLED=1` / **每形一枚新容器** / `--rm` / 命名卷 `ac119-gomodcache`+`ac119-gocache`。

七发 `-v` 每发进容器第一件事打挂载自证，七发**逐字相同**：
`-rwxrwxrwx 1 root root 1026 /src/go.mod` + `md5 go.mod 0c266f4d52895a65100bee94dc16bdcb`
+ `md5 internal/winsec/resolve.go b6876a5efe759f6e17434d1b50a129c3`（与票 124 前各批逐字同字）。

软链形 `ln -s /realpriv /varlink` + `TMPDIR=/varlink/w124tmp` + `exit 90/97/98`；
普通形另加 `exit 99`（`/varlink` 必须根本不存在）＋ `/plainroot`。**七发无一命中形状断言**
（每发 `SHAPE=... TMPDIR=...` 行逐字在 sec 5 的日志目录里）。

## 3. 两把分母分开报（NOT 抄任何一个数）

### (1) only-in-link（票 124 那族）= 本腿的清零目标

```
$ comm -23 PRE-LINK.fail.txt PRE-PLAIN.fail.txt   =>  37  名
$ comm -13 PRE-LINK.fail.txt PRE-PLAIN.fail.txt   =>   0  名   (only-in-plain 空)
$ 子测试 only-in-link                              =>  12  名
```

四数（每数分列，NOT 相减）：

| 发 | 顶层 RUN | 顶层 PASS | 顶层 FAIL | 顶层 SKIP | 子 RUN | 子 PASS | 子 FAIL | npr | rts | dpapi |
|---|---|---|---|---|---|---|---|---|---|---|
| PRE-LINK  | 171 | 66  | 105 | 0 | 63 | 42 | 21 | 115 | 115 | 0  |
| PRE-PLAIN | 171 | 103 | 68  | 0 | 63 | 54 | 9  | 1   | 1   | 70 |

（`npr` = `not provably resolved` 行数、`rts` = `refusing to seal` 行数、`dpapi` = DPAPI 串行数；
三者**只是机制旁证，NOT 判据**。`panic` 与 `test timed out` 七发**各 0**。）

**对账裁腿 `124-v1`**：它记 37 名 / 13 文件、两形都红顶层 68 + 子 9。
`comm -23 mine verifier` **0 行**、`comm -13 mine verifier` **0 行**（37 那侧两向都 0）；
68 那侧同样两向 0 行 => **分母与裁腿逐名 set-equal**。

### (2) 两形都红（DPAPI 那族）= 登记，NOT 修、NOT 计入出口

**顶层 68 + 子 9**。逐名机制归类（`harness/classify.py`，块内命中 DPAPI 串即归该族）：

```
DPAPI=64  NON-DPAPI=4  TOTAL=68
```

14 枚非 DPAPI 的 4 枚逐名与逐字首红行：

| 名 | 机制 | 首红行（节选） |
|---|---|---|
| `TestAC1AC2DispatchHopGate133` | 票 133 门 | `leg_dispatch_gate_133_test.go:251: AC#1/#2 RED: leg "slo" (main.go:110) ...` |
| `TestAC4EveryLegIsNailedOrRuled` | 票 131 门 | `leg_sink_gate_131_test.go:275: AC#4 RED: leg "models" (main.go:102) ...` |
| `TestCleanCheckoutBuilds_AC11` | 真实 git 环境 | `panel_host_gate_test.go:404: git rev-parse --show-toplevel: exit status 128` |
| `TestSecretOverwriteIsAnnounced` | DPAPI 下游 | `secret_test.go:1196: first set: 1` |

DPAPI 那 64 枚失败串几乎全是 `secret: DPAPI is only available on Windows; use an env: ref on this platform`
（`internal/secret/protect_other.go:19`，该件带 `//go:build !windows` => **容器 Linux 档才有、宿主 Windows 档不红**）。
这 64 枚在 PRE-PLAIN / POST-PLAIN / MUT-PLAIN / RESTORED 里**逐名不动**（名册 `diff` 0 行，sec 5）。

## 4. 十三枚文件逐枚：原来交了什么、现在过哪一跳

（"递根点"＝把 `t.TempDir()` 的答案交给密封底线那一个表达式。下表 13 枚是按用例名归集的定义文件；
**实际改的是 9 枚**——另 4 枚（`panel_config_320` / `panel_delete_320` / `panel_models_320` / `config_reload_223`）
的根全部来自这 9 枚里的共享 helper，故零 hunk。）

| 文件 | 原来交的根 | 现在过哪一跳 | 覆盖名 |
|---|---|---|---|
| `panel_config_248_test.go` | `dir := t.TempDir()`（:88，`settingsDir248` 体内） | `sealableTempDir124(t)` | 5（并供 320 族 9 处调用） |
| `panel_config_313_test.go` | `dir := t.TempDir()`（:52 `empty313Dir` 体内；:296 用例体内） | 两处 `sealableTempDir124(t)` | 2 |
| `panel_inbound_33_test.go` | `dir := t.TempDir()`（:63，`newInboundDataDir33` 体内） | `sealableTempDir124(t)` | 4（并供 223 那枚） |
| `firstrun_198_test.go` | `dir := t.TempDir()`（:75/:158/:193/:238，四枚用例各自的根） | 四处 `sealableTempDir124(t)` | 4 |
| `firstrun_198r2_test.go` | `dir := t.TempDir()`（:222/:364）＋ `runDir := t.TempDir()`（:462） | 三处 `sealableTempDir124(t)` | 3 |
| `firstrun_257_test.go` | `dir = t.TempDir()`（:74，`t257CleanMachine` 体内） | `sealableTempDir124(t)` | 6 |
| `firstrun_257_nonpreset_test.go` | `dir = t.TempDir()`（:50，`t257npCleanMachine` 体内） | `sealableTempDir124(t)` | 2 |
| `config_sentences_223r2_test.go` | `filepath.Join(t.TempDir(), "config.toml")`（:105，子测试体内） | `filepath.Join(sealableTempDir124(t), "config.toml")` | 1（含 9 子测试） |
| `restart_tier_keys_255r2_test.go` | `filepath.Join(t.TempDir(), configFileName)`（:95，`reloadOneAppKey255` 体内） | `filepath.Join(sealableTempDir124(t), configFileName)` | 1（含 3 子测试） |
| `panel_config_320_test.go` | 经 `settingsDir248`（无自有递根点） | 由上一枚 helper 带上 | 3 |
| `panel_delete_320_test.go` | 同上 | 同上 | 2 |
| `panel_models_320_test.go` | 同上 | 同上 | 3 |
| `config_reload_223_test.go` | 经 `newInboundDataDir33`（无自有递根点） | 由上一枚 helper 带上 | 1 |

**零新增走法**：只用票 124 批 4 已入树的那一枚单行委托（`tempdir_resolved_124_test.go:35-37`，
函数体逐字 `return proc.SealableRoot(t.TempDir())`）。**本腿未碰该文件**（`git diff` 空、`sha256` 见 sec 5）
=> 票 125 `R-125-2` 那本副本账**不增行**。`cmd/wisp` 侧调用点由 3 处（`providers_test.go:40` / `secret_test.go:667`
/ 定义处）增到 **18** 处，全走同一枚。

**按名登记"有意不改"的 3 处裸 `t.TempDir()`**（同族内，**实测两形都 PASS**，摘解析无可观测差别）：

| 站点 | 为什么不动 | 两形实测 |
|---|---|---|
| `panel_inbound_33_test.go:415` `empty := t.TempDir()` | 该用例把 `empty` 交给 `t.Chdir`，根**不进密封底线**；它的断言是"拒绝腿一个字节都不写进 CWD" | 两形各 PASS |
| `firstrun_198r2_test.go:300` `filepath.Join(t.TempDir(),"data")` | AC#5 用例的**被测对象就是"建不出来"**：它在根里种一个挡路的 `secrets` 文件，交的是一枚**故意不可建**的根；接解析会把这条用例的语义改掉 | 两形各 PASS（PRE-LINK `--- PASS ... 0.03s`） |
| `firstrun_198r2_test.go:434` `rootDir := t.TempDir()` | J1(a) 半边的断言**正是"装配根什么都不建"**；摘掉解析会让它读成"建不出来"的假绿。它已由 (b) 半边 `runDir`（本腿已转）承担真实建文件那一跳 | 两形各 PASS |

**这 3 处不是漏掉、也不是判断**：我另打一台 `PROBE3`（`post` ＋ 把这 3 处也换成 `sealableTempDir124`）
在软链形跑一遍 —— 顶层 `171/103/68/0` + 子 `63/54/9/0`，红名册与 POST-LINK **双向 `comm` 差集各 0 行**、
子测试名册 `diff` 0 行 => **转了也是 0 个名**。设备在 `logs/PROBE3-LINK.log`。
