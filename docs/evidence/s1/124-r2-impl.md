# 票 124 落地腿 124-r2 —— internal/tools 一批：用例侧递根点接上 proc.SealableRoot

本件是**实现者**交件，NOT 裁决。裁决表由验收方另出（`docs/evidence/s1/124-*.md` 既有各份不属本腿）。

## 0. 一句话

把 `internal/tools` 里"该转还没转"的 **35 枚可转**用例（38 名 only-in-link 减票 123 那三枚 300 s 豁免腿）的
**递根点**接上票 119 那条纪律（OS 给的答案先过 `proc.SealableRoot` 再递密封底线）：软链 TMPDIR 形状下
本包 only-in-link 红名 **38 -> 3**（余 3 = 票 123 豁免），**新增红 0**（名册集合，NOT 枚数相等）。

## 1. 开工前两把尺（我自己现跑，NOT 抄派单）

(1) 写面脏态：`git status --porcelain -- internal cmd go.mod go.sum` = **空**（写面内无别人未提交改动；
全树 1250 条脏项全在 `.scratch/**` 等非写面）。开工首读 HEAD `c50218bf`，收尾复读 `e61861b4`，
`git diff --name-only c50218bf HEAD -- internal cmd` **零 hunk** => 量的是同一版码。

(2) **票面前提新鲜度（派单点名的腐坏尺）**：那把"`winsec.` 与 `TempDir` 同行"的静态尺今天**恒等于已修态**，
本腿**一枚都没拿它当分母**（分母全部走 sec 2 的终判据：纯净快照 + 两形各一枚新容器 + `-count=1 -v` + 名册 `comm`）。

(3) 分母独立复跑（NOT 抄 35/16）：见 sec 2 —— 我一手量到 **38 名 / 16 文件**，扣豁免三枚 => **35**，与派单逐数对上。

辅助判据上限一栏（防"四数全 0 被读成没红"）：五发 `-v` 的 `RUN` 行恒 **255**，`panic` 与 `test timed out` 各 **0**。

## 2. 分母复跑（两形名册 `comm` 差集，NOT 枚数相减）

台面：纯净快照 `git archive HEAD | tar -x`（17504 文件）、`golang:1.27` / `go1.27.1 linux/amd64`、
`GOPROXY=off`、**每形一枚新容器**、挂载自证 `go.mod` 1026 字节 + `md5 0c266f4d…` + `resolve.go md5 b6876a5e…`。
软链形 `ln -s /realpriv /varlink` + `TMPDIR=/varlink/w124tmp` + `exit 90/97/98`；
普通形另加 `exit 99`（`/varlink` 必须根本不存在）。**五发无一命中形状断言。**

```
$ comm -23 PRE-LINK-fail.txt PRE-PLAIN-fail.txt | wc -l    =>  38
$ comm -13 PRE-LINK-fail.txt PRE-PLAIN-fail.txt | wc -l    =>   0     (only-in-plain 空)
$ 子测试 only-in-link                                       =>   3
$ 38 名落到文件：distinct files                             =>  16     UNKNOWN=0
```

逐文件分布：`pointer_185_cli_seam` 6、`task_output_pointer_notice` 4、`grant` 4、`wiring` 3、
`ticket175r2_stamp_live` 3、`subagent_depth_fixture_287` 3、`tier_consistency_172_cases` 2、
`ticket176r1_start_port` 2、`task_output_leg` 2、`refusal_20r2` 2、`pointer_183_cli_seam` 2、
`ticket285_corr_rows_rulers` 1、`ticket283_corr_identity_rulers` 1、`ticket224_grantid_column` 1、
`task_cancel_221_legs` 1、`plan_present_165` 1。

**票 123 那族三枚 300 s 腿逐名在册**（本发一手量到时延，全部 `300.04s`）：
`TestL1WriteGoesThroughTheRealBlockWindow`、`TestLateVetoRendersTheApprovalLayersAppliedStepsReport`、
`TestFSReadOnlyNeverOpensACard` => 按账**豁免**、**不计入分母**、**未顺手修**。

## 3. 十六枚文件逐枚：原来交了什么、现在过哪一跳

（`=> ` 在本件写作 `=> `；"递根点"指把 `t.TempDir()`/`tempCanonical` 的答案交给密封底线或 C26 的那一处表达式。）

| 文件 | 原来交的根 | 现在过哪一跳 | 覆盖红 |
|---|---|---|---|
| `plan_present_165_test.go` | `dir := t.TempDir()`（:209，清单外那一发 fs.write 的根） | `sealableTempDir124(t)` | 165 一枚 |
| `tier_consistency_172_machine_test.go` | `in, out := t.TempDir(), t.TempDir()`（:36，机器侧两枚根） | 两处 `sealableTempDir124(t)` | 172 两枚 |
| `refusal_20r2_test.go` | `cfgPath := filepath.Join(t.TempDir(),"config.toml")`（:182/:276）＋ `outside := tempCanonical(t)`（:192/:285） | 四处 `sealableTempDir124` / `sealableTempCanonical124` | 20r2 两枚（含 1 子测试） |
| `task_cancel_221_legs_test.go` | `x.mem = openStore(t, t.TempDir())`（:82，真 SQLite 的根） | `openStore(t, sealableTempDir124(t))` | 221 一枚 |
| `ticket283_corr_identity_rulers_test.go` | `ArtifactsDir: filepath.Join(t.TempDir(),"artifacts")`（:160） | `sealableTempDir124(t)` | 283 一枚 |
| `ticket285_corr_rows_rulers_test.go` | 同上（:42） | `sealableTempDir124(t)` | 285 一枚 |
| `subagent_depth_fixture_287_test.go` | 同上（:84） | `sealableTempDir124(t)` | 287 三枚 |
| `ticket176r1_start_port_test.go` | `dir := t.TempDir()`（:59，Spiller 的 artifacts 根） | `sealableTempDir124(t)` | 176r1 两枚 |
| `ticket175r2_stamp_live_test.go` | `dir := tempCanonical(t)`（:184/:208/:281） | `sealableTempCanonical124(t)` | 175r2 三枚 |
| `task_output_pointer_notice_test.go` | `dir := tempCanonical(t)`（:148/:205/:351）＋ `rootA, rootB := tempCanonical(t), tempCanonical(t)`（:498） | 四处 `sealableTempCanonical124(t)` | 174 族四枚 |
| `task_output_leg_test.go` | `dir := tempCanonical(t)`（:166/:347） | `sealableTempCanonical124(t)` | 164 族两枚 |
| `pointer_183_cli_seam_test.go` | `dir := t.TempDir()`（:94，`p185Seam`/`r183Seam` 同形） | `sealableTempDir124(t)` | 183 两枚 |
| `pointer_185_cli_seam_test.go` | `dir := t.TempDir()`（:93） | `sealableTempDir124(t)` | 185 六枚 |
| `grant_test.go` | `dir := t.TempDir()`（:195/:246/:388/:448） | 四处 `sealableTempDir124(t)` | 224 四枚（含 2 子测试） |
| `ticket224_grantid_column_test.go` | `memory.Open(filepath.Join(t.TempDir(),"data"))`（:41）＋ `dir := t.TempDir()`（:101） | 两处 `sealableTempDir124(t)` | 224 一枚 |

**零新增走法**：只用票 124 前几批已入树的两枚单行委托
（`tempdir_resolved_124_test.go` 的 `sealableTempDir124`/`sealableTempCanonical124`，
`tempdir_resolved_124x_test.go` 的外包孪生一枚），函数体逐字是 `return proc.SealableRoot(t.TempDir())`
（sha256 `3d949a47…` / `9c4b7381…`，与 HEAD blob 同字，**本腿未碰这两枚文件**）
=> 票 125 `R-125-2` 那本三枚副本账**不增行**。

**按名登记"有意不改"的 6 处裸 `t.TempDir()`**（都在同一族内，但**两形都 PASS**，摘解析无可观测差别，
故 NOT 计入 35 枚、也 NOT 改）：`plan_present_165_test.go:123/:148`、`task_output_pointer_notice_test.go:579`、
`grant_test.go:302/:303`（L2/Deny 拒绝腿，根不进密封底线）、`grant_test.go:349`（`HOME` 重定向）、
`grant_test.go:422`（无路径调用）。逐枚实测见 sec 6 的"两形恒 PASS"名册。

## 4. 判定分支未动（票面 AC#3）

```
$ git diff --numstat -- internal/tools
4  4  grant_test.go                    1  1  plan_present_165_test.go
1  1  pointer_183_cli_seam_test.go     1  1  pointer_185_cli_seam_test.go
4  4  refusal_20r2_test.go             1  1  subagent_depth_fixture_287_test.go
1  1  task_cancel_221_legs_test.go     2  2  task_output_leg_test.go
4  4  task_output_pointer_notice_test.go
3  3  ticket175r2_stamp_live_test.go   1  1  ticket176r1_start_port_test.go
2  2  ticket224_grantid_column_test.go
1  1  ticket283_corr_identity_rulers_test.go
1  1  ticket285_corr_rows_rulers_test.go
1  1  tier_consistency_172_machine_test.go
```

`sorted(added) == sorted(deleted)` 逐枚对称（**28 增 / 28 删**）。**删除侧全文只有 28 行**（`--unified=0` 现量，
`grep -c` 与逐行枚举两条尺同数），逐行都是"根路径怎么来"那一个表达式：

```
  - dir := t.TempDir()                                            x9
  - dir := tempCanonical(t)                                       x8
  - cfgPath := filepath.Join(t.TempDir(), "config.toml")           x2
  - outside := tempCanonical(t)                                    x2
  - ArtifactsDir: filepath.Join(t.TempDir(), "artifacts")          x3
  - rootA, rootB := tempCanonical(t), tempCanonical(t)             x1
  - s, err := memory.Open(filepath.Join(t.TempDir(), "data"))      x1
  - x.mem = openStore(t, t.TempDir())                              x1
  - in, out := t.TempDir(), t.TempDir()                            x1
                                                  total            28
```

**断言／阈值／golden／`t.Skip`／判定分支 0 hunk。** 用 `Fatal|Errorf|want|expect|assert|Skip|threshold|golden`
扫全部 `+/-` 行，命中 **0**。禁改列（`internal/{winsec,proc,risk}/**`、`internal/agent/**`、`cmd/wisp/**`、
`internal/config/**`、`internal/tools` **生产码**、`tools/d22scan/**`、`allowlist.txt`、`thresholds.go`、
任何 golden、`go.mod`/`go.sum`、`frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`docs/reports/**`）
**一律零 hunk**（`git diff --name-only` 仅 sec 3 那 15 枚 `*_test.go`）。

**票 113/119 那两族拒绝用例未变成通过方式**：本包内 `assertRefused113`/`leafLinkTo118`/`resolveProbeRoot`
0 命中（它们在 `internal/winsec/`）；票 125 seam 探针 `c26_seam_posix_125_test.go` 不在本包。
决定性一条与批 1 同形：sec 5 的 35 枚 `went-green` 在**改前普通形逐名 PASS**（`PRE-PLAIN` 那侧）=> 断言本就成立，
这次摘掉的只是 harness 形状造成的那次拒，**放宽换绿在这个形状下无从发生**。

> 本件因 160 行帽拆四份：本件 = sec 0..4；`124-r2-impl-part2.md` = sec 5..7（三读数与 sha256 还原、
> 名册集合出口判据、豁免逐名）；`124-r2-impl-part3.md` = sec 8..10（门禁 rc、红名册逐名、未验证项）；
> `124-r2-impl-part4.md` = sec 11（提交号、逐枚 `.go`、越界与正控）。
