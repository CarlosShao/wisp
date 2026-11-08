# 354-a1 现读原始读数（AC#0 只读腿；零产码字节；⛔ 翻框）

- 日期：2026-11-08 21:1x +08。腿＝`354-a1`（票 354 排程第一节＝只读现读＋归因）。
- 读数锚：两棵独立导出树（起手三核③）＝`/d/tmp/wisp/20261108/leg78/tree-leg78` 与 `tree-leg78b`，
  两棵 HEAD 同为 **`ca790c46e`**（`git clone --shared` 于 21:10/21:11 取得；共享树起手 HEAD `75ff9c4eb`，
  `internal/winsec` 起手 `git status --porcelain` 空＝非他腿半保存态）。
- 跑前环境：`unset GOTMPDIR TMPDIR TEMP TMP`；⛔ CGO_ENABLED=0 未设（默认 cgo 开，winsec 仪器需要真 Windows 调用）。
- 原始日志落盘：本目录 `01-ac2windows-instruments.log`（定向 `TestAC2Windows` 6 发，rc0）＋
  `02-winsec-full-package.log`（`./internal/winsec` 整包，`ok ... 20.148s`，rc0）。

## 一、两族仪器现跑读数（Windows 档）

### 族 A＝`ancestorunlink_348_windows_test.go`（删除分子那条族）

- `TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1` PASS（0.88s）：
  `form=symlink-dir attempts=400 cycles=67 provedReparse=67 formBLegRefusals=12 truncated=0 deleted=0 returnedNil=76 flipErrs=11`；
  回答腿分布：C26 稳态 237／pin-walk 拒 51／act-time reparse 拒 12／seal 拒 10／C26 拼写拒 14／returnedNil 76。
- `TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1` PASS（0.76s）：
  `form=leaf-symlink attempts=400 cycles=81 provedReparse=83 formBLegRefusals=6 truncated=0 deleted=0 returnedNil=66 flipErrs=3`；
  **残留原句（`:169` 逐字）**：`AC#2' RESIDUAL, deterministic on this spelling: the followed open left the bystander as "" (read err <nil>) - the bytes were refused afterwards and the NAME was never removed, which is the half this platform may claim`；
  叶腿拆分（REPORTED not gated）：`preOpenLeafGuard=0, openPinnedLeaf re-check=0, provedReparse=83`。
- `TestAC2WindowsPinRefusalOnAnAbsentAncestorIsNotFiledAsASealRefusal348r1` PASS（0.01s）：absent 祖先归 pin-walk 行、不记 seal 行。

### 族 B＝世界台件消费者（`ancestorworld_342_348_windows_test.go` 提供的世界；本机命中三枚）

- `TestAC2WindowsSteadyStateAncestorJunctionReportsDamageAndRefusalAsTwoNumbers342r2` PASS（0.17s）：
  三发 deterministic，两列均 0 damage、拒绝来自**任何 act 之前**的腿（C26 稳态，verbatim
  `risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions`）；
  用例自陈该形复现 348-v1 的"缺世界"，且 **NOT ticket 348's pre-landing numerator**。
- `TestAC2WindowsNameTopUpBeforeTheReCheckLeavesTheTopUpObjectNamed353r1` PASS（0.03s）——**残留①的确定性仪器面**：
  三发汇总逐字 `entered=3/3, worlds the platform froze=0, unnamed objects left with no name (DAMAGE)=0, refusals that left the caller's name standing=3`；
  每发 DAMAGE 行逐字：`the caller's name=present (size 64 body [an object the caller never named, standing in the same directory]) | the swapped-in object's own name=absent ... | swapped-in object left with no name=false`；
  GREEN SCOPE 声明逐字在场（recheck 在 removal 之前 ⇒ 拒掉所有已落地的 top-up 并**留下调用者的名字**；
  recheck 应答之后、`Deleteat` 的 NtOpenFile 绑定之前落地的 top-up ⛔ 覆盖、平台内部无钩）。
- `TestAC2WindowsRemovalStillRunsWhenNobodyTopsUpTheName353r1` PASS（0.01s）＝**负控**：
  无 top-up 时移除照跑，`name=absent`——证明"standing"计数不是恒 3 的钝尺。

## 二、正控（尺能命中非零，先跑）

1. **分子尺非零命中**＝上面 353r1 的 shot 级 DAMAGE 行：调用者名字下站着**64 字节攻击者对象本体**（body 逐字可读），
   计数器 `refusals that left the caller's name standing=3/3`＝非零且确定性。
2. **负控**＝RemovalStillRuns 那发 `name=absent`——同一把尺在"没人 top-up"世界读 0，尺有两向。
3. **敏感度（引用，未重跑）**：`353-v1` 判语第 3 件记录的三发去牙（摘删除 ⇒ 正控转红 `name=present (size 0 body [])`）；
   本腿不重跑去牙——那是 `354-r1` `AC#3` 的活，本腿只交现读。

## 三、本机跑不了的那半边（如实登记）

`ancestorunlink_348_unix_test.go`／`ancestorworld_342_348_unix_test.go` 带 unix 构建标签，本机（Windows）⛔ 可跑。
两族读数本腿只交 **Windows 档**；unix 半边登记〔仅另机/容器可量，未跑〕，不冒领。

## 四、门禁（只读腿的分档）

- 本腿 `.go` 字节＝0（`git status --porcelain -- internal/ cmd/` 起手⇄落笔同空）⇒ `gofumpt -l` 射程＝空集（非"跳过"）。
- `./internal/winsec` 整包色＝上面 02 号日志（rc0）——两族仪器所在包的现势色由本腿亲跑。
- `sh scripts/d22scan.sh`：另发现跑（见交件提交正文）。
