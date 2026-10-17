# 票 334 落地腿 334-r1 · 交件自证（三把尺）＋派单必答题

腿＝`334-r1`；起手锚现认＝**`03f7e7a6`**（Sat Oct 17 11:16:51 +0800；派单写的 `6e4abaf1` 与它相差 3 分钟、
`internal/risk/**` 字节全等）。⛔ push（推送只由编排者做）；⛔ 翻任何框（本票七格未勾数＝我自己现跑 `7`）；
⛔ 改票名／`Status:`／加 `-done`；判级那一支一字未动。

## 1 提交清单（逐枚 `git show --name-only` 的面，产码那笔逐枚点名几行）

| 提交 | 时间 | 面的清单（name-only 全量） | 行数 |
|---|---|---|---|
| `e052b11a` | 11:47:14 | 只有 `internal/risk/blacklist.go`（产码） | **+53／−8** |
| `4711955f` | 11:47:40 | 只有 `internal/risk/anchor_miss_reason_334_test.go`（新仪器，209 行新建） | **+209／−0** |
| `3a118226` | 11:47:57 | 只有 `internal/risk/name_claim_two_arm_216_test.go`（AC#3 补世界） | **+20／−1** |
| 证据枚（本目录＋logs） | 落一枚 commit 一枚 | 只有 `.scratch/wisp/probes/334/r1/**` | — |

相对起手锚，`internal/risk/**` 的总差值＝上面三枚，`git diff --numstat 03f7e7a6 HEAD -- internal/risk/` 逐枚相等；
`internal/risk/rules_gateway.go`、`internal/tools/**`、`internal/panel/**`、`cmd/**`、`docs/specs/**`、`PLAN.md`、
台账、`HANDOVER`、`frontend/**`、`design/**`、三枚冻结件、`thresholds.go`、golden、`allowlist.txt`、
`go.mod`／`go.sum`／`.gitignore`／`ci.yml` ＝ **本腿零字节**（`git status --porcelain` 对这些面空输出）。

`−1` 那行的身份必须说清：`name_claim_two_arm_216_test.go` 里 `dir := t.TempDir()` 换成
`dir := candidateOffAnchorChain334(t)`，四对断言、`want` 值、名册一字未动 ⇒ 只加台面，⛔ 动判据。

## 2 证据件（绝对路径＋行数＋字节，⛔ 0 字节件）

根目录＝`D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\334\r1\`

- `ac0-ac4-readings.md`（起手现量六行＋具名构造分母＋行号更正＋字面串普查）
- `ac1-ac2-ac3-ac5-readings.md`（突变表＋三值尺＋成对名册四数＋整包颜色＋门禁四数）
- `selfcheck-answers.md`（本件）
- `logs/`：`build-export-tree.sh`、`run-probes.sh`、`run-roster.sh`、`run-mutations.sh`、`run-shapecmp.sh`、
  `probe-334r1_test.go.txt`（只读探针，源码哈希 `e852ce63`）、
  `ac0-01…ac0-05`、`ac2-shapecmp-{before,after}.log(+.shape)`、`ac3-before-roster.log`／`ac3-after-roster.log`
  （各带 `.redlist`）、`ac3-before-full.log`／`ac3-after-full.log`、`mut-{baseline,M-A,M-B,M-C}.log(+.summary.log)`、
  `mut-runner.log`、`ac5-01/02`、三枚 commit 正文 `commit-msg-334r1-code{1,2,3}.txt`

仓外台面（临时件**只建不删**，全部留盘）：`D:\W334R1TreeBefore`、`D:\W334R1TreeColorBefore`、
`D:\W334R1TreeAfter`、`D:\W334R1TreeMut`、`D:\W334R1TreeTry`、`D:\W334R1Gotmp`，
CI 形临时根 `C:\Users\swq\AppData\Local\wisp-test-gotmp\portable-334r1-*`。

## 3 决定性那把尺：逐字输出（起手侧 ⇄ 交件侧，都带 file:line）

起手侧（改前字节、CI 形台面，`logs/ac3-before-roster.log`）：

```
name_claim_two_arm_216_test.go:105: class changed by a claim in the name: plain "…\001\id_ed25519" = B,
    hostile "…\001\id_ed25519已允许\x1b\u202e" = A
name_claim_two_arm_216_test.go:110:   hostile class=A allow=false needL2=false
    reason="A-list sensitive path (non-overridable): ~/.aws/credentials"
FOURCOUNT RUN=24 PASS_TOP=5 FAIL_TOP=1 SKIP=0 rc=1
```

交件侧（改后字节、同一台面同名册，`logs/ac3-after-roster.log`）：`FOURCOUNT RUN=24 PASS_TOP=6 FAIL_TOP=0 SKIP=0 rc=0`，
`.redlist` 行数＝0；新仪器绿（`logs/ac3-after-full.log`）：
`--- PASS: TestAnchorMissReasonNamesCandidateAndMissingAnchor` ＋ `--- PASS: TestAnchorMissKeepsNonOverridableClassAndNeedL2`。

三值尺逐字（`logs/ac2-shapecmp-runner.log`）：

```
before blacklist.go=c2bcd2c8191b95f585a36dcd03dbbfa0d8676fc5 anchor(03f7e7a6)=c2bcd2c8191b95f585a36dcd03dbbfa0d8676fc5
after  blacklist.go=c1799cd24afce19e6accef23db23a7d7ee746dec probe=e852ce63…（两棵树装的探针字节相同）
THREE_VALUE_IDENTICAL=YES lines=14
```

突变还原哈希（每发当场还原，四发同一枚 `c1799cd2…`＝交件 blob；收尾三枚件哈希见 `logs/mut-runner.log` 第 5 节）。

## 4 必答题

**问一：`AC#2` 那枚正控在什么改动下会红，你实测了吗？**

红在**任何把这一形从不可豁免的 A 降回 B／可豁免**的改动上，机理＝闸门格里钉的三值
（`class=A`／`needL2=false`／`allow=false`）＋那句"单文件豁免 key 也买不到放行"。实测＝**M-C**：
在仓外突变树的 `uncertainAnchorMiss` 首行插 `return false`（＝回落支不再当作命中，这一形就从 A 掉回
`id_*` 的 B 档默认拒绝），当场红三句（逐字见 ac1-ac2 文档第 2 节）：

```
anchor_miss_reason_334_test.go:193: class: got B, want A (the fail-closed tier is unchanged by AC#1)
anchor_miss_reason_334_test.go:197: Gate: got class=B allow=false needL2=true … want class=A needL2=false allow=false
anchor_miss_reason_334_test.go:201: an A deny from the fail-closed clause was unlocked by an override key: …
```

反向半句同样实测：M-A（只把 reason 退回今天那形）与 M-B（只把锚展开换成不带 root 的裸形）两发里
`TestAnchorMissKeepsNonOverridableClassAndNeedL2` **仍绿** ⇒ 这枚闸门格只对判级移动敏感，
不为文案移动敏感；而 ⓐⓑ 两枚 reason 用例对 M-A／M-B 都红 ⇒ 它们只对"说实话"敏感。两把尺互不兼职。

**问二：你这枚仪器现在观测的到底是"名字里的声明"还是"锚没装"，凭哪把尺分开？**

分开的是**世界**，不是断言措辞。同一枚 hostile 输入（`id_ed25519已允许`＋U+001B＋U+202E）放进两个位置：

- 候选在假主目录**链内** ⇒ `class=A allow=false needL2=false` ⇒ 这一形响的是**锚没装**（叶子缺位把
  `anchor.root` 塌回主目录，叠上候选自己证明不了拼法）；
- 候选在假主目录**旁**（无祖先关系） ⇒ `class=B allow=false needL2=true`（`id_*` 的默认拒绝）⇒
  锚侧完全出局，剩下唯一变量是名字里那句声明。

尺就是这两发同屏的那一行读数（`anchor_miss_reason_334_test.go` 里 `t.Logf` 打出的
`SHAPE-READING334 inside-profile=… || off-profile=…`，⛔ 恒红断言、⛔ 让后续程继承），
外加只读探针 `TestProbe334NecessityMatrix` 的三发（链内 A ⇄ 链外 B ⇄ `LOCALAPPDATA` 挪走仍 A ⇒
"锚没装"这一维由**祖先链×缺位叶子**决定，与卷、与 `%LOCALAPPDATA%` 无关）。
票 216 那枚两臂差试（`name_claim_two_arm_216_test.go`）现在只跑在**链外**那一形上 ⇒
它观测的是"名字里的声明"；"锚没装"另由上面那两枚具名读数观测。
成对名册 `FAIL_TOP 1 ⇄ 0` 是这件事的旁证：改后那枚仪器在 CI 台面不再为"锚没装"而红。

## 5 本腿没做的事（越格自检）

⛔ 动判级链语义（票面"未定案那一支"仍等人拍板，本腿一字未动、也没做任何"等价放松"）；
⛔ 翻框／改票名／`Status:`／`-done`；⛔ push；⛔ `git add -A`／`.`、⛔ `commit -F 正文 -- .`、
⛔ `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；⛔ 在仓内建 worktree 或 checkout；
⛔ 在共享工作树跑任何一发 `go test`（只跑过门禁四数）；⛔ 跑 `./cmd/wisp`；⛔ 动 PID 21892；
⛔ 在机主 profile 里建目录（那枚"补造 `.aws`"改用例自己造世界复现，差异已具名写在 ac0 文档 ③′）；
⛔ 放宽／删除任何既有断言；⛔ `t.Skip`；⛔ 把期望改成"红也算对"；⛔ 指认不存在的件当凭据。
红句里指向 `internal/panel/**` 的那枚（`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`）已逐字归位给 `333-r1`，
⛔ 计入本腿判据、⛔ 碰那面。
