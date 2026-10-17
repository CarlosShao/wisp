# 票 334 裁腿 334-v1 · 读数件（⛔ 判语；判语在 `docs/evidence/s1/334-v1-verdict.md`）

腿＝`334-v1`（非实现者终裁）。首枚现量 `12:35:27` HEAD＝`648eda04`；三件入库那一刻（`logs/00-intake.log`）HEAD＝`0a030cea`、
跟踪件面当时三枚脏（票 328／334／335，⛔ 本腿字节，编排者 13:26 `ed6fa14e` 归位）；起手锚现认＝`03f7e7a6`（11:16:51）。
两枚号同份这把尺我复跑了：`git diff --stat 6e4abaf1 03f7e7a6 -- internal/risk/` ＝ **空输出 rc=0**，
两枚号全差只含票 334 正文＋派单三件（`+123/-0`）⇒ 腿的"两号同份"成立。

## 0 台面（⛔ 共享工作树跑过任何一发 go test／go build）

尺＝`logs/build-tree.sh`（八件套＋仓根 `docs`＋手工三件），三棵树全在 `D:` 卷且逐枚读到 8.3 别名：

| 树 | 锚 | `.go` 枚数 | blacklist.go blob | 8.3 别名 |
|---|---|---|---|---|
| `D:\W334V1TreeBefore` | `03f7e7a6` | 742 | `c2bcd2c8…` | `W334V1~2` |
| `D:\W334V1TreeAfter` | `3a118226` | 746 | `c1799cd2…` | `W334V1~3` |
| `D:\W334V1TreeMut` | `3a118226`（突变只打这棵） | 746 | `c1799cd2…` | `W334V1~4` |

`docs_plan=present`／`dlls_build=3 dlls_sherpa=3 dist=present` 逐枚读到（⛔ 0＝台面故障）。
CI 形台面＝`GOTMPDIR=$LOCALAPPDATA\wisp-test-gotmp\portable-334v1-…`（逐字照 `portable-tests.sh:152-159`），
`TMP/TEMP` 一字不动，`mkdir_gotmp_rc=0`。

## 1 AC#2 三值尺：我自己那把（42 行）在两棵树上对撞

只读探针 `logs/probe-334v1_test.go.txt`（哈希 `69a70fa4…`，两棵树装的字节相同）跑 `-run TestProbe334V1`：

- `logs/72-probe-v2-before.log` ⇄ `logs/73-probe-v2-after.log`，四数两发同为 `RUN=22 PASS_TOP=4 FAIL_TOP=0 SKIP=0 rc=0`
- `diff logs/v2-before.shape logs/v2-after.shape` ＝ **0 行**（`THREE_VALUE_DIFF_lines=0`），42 枚具名输入
  的 `class／allow／needL2／reason 点名的规则` 逐字相等 ⇒ 腿的 `THREE_VALUE_IDENTICAL=YES lines=14` 我这把更宽的尺上同样成立。
- 我这把尺比腿的宽在哪：腿的 14 行只经过 **`.git-credentials`／`wisp config.toml`／`.aws`** 三枚锚（真机＋两枚假主目录），
  我按七枚走 `aEq/aUnder` 的锚逐枚造世界（`chain:case1/3/4/5/6/8/9` ＋ 两枚 solo `case4b/8b`），
  每枚再带一发 override-key 行与一发"干净文件"行与一发"链外"行。基线读数（`logs/60-v2-baseline-chain.shape` 18 行）：
  `chain:*` 全部 `class=A allow=false needL2=false`，`chain-clean:*` 全部 `class=none`，`offchain:*` 全部 `class=B needL2=true`。

## 2 突变八发（全打在 `D:\W334V1TreeMut` 的盘上字节；逐发「从具名提交重取字节覆盖＋断言哈希」）

⚠ 自纠：台件 v1 的 perl 还原锚写坏（列 0 的 `}` 被写成制表符）⇒ `M-V1` 残留在突变树上，
v1 的 `M-V2／M-V3` 两发是在脏态上取的、**作废**；v2 全部重跑，末态三枚哈希＝交件 blob
（`c1799cd2…`／`8f683b04…`／`cbc0a610…`），逐字见 `logs/60-v2-runner.log` 末节。

| 发 | 落点（内容锚，⇄ 交件 blob 的行差逐字在日志里） | 腿的三枚仪器名册四数 | 我这把 42 行尺看到什么 |
|---|---|---|---|
| 基线 | — | `10／3／0／0 rc=0` | 18 行全 A／none／B 三形分明 |
| **M-V1** 只动判级支：`aRelates` 尾行 `return uncertainAnchorMiss(cand, anchor)` → `return false` | `blacklist.go:281` | `10／1／2／0 rc=1`，红＝闸门格＋ⓐⓑ | 七枚 `chain:*` 全 `A→B needL2=true` |
| **M-V2a** 只摘 `%LOCALAPPDATA%\Microsoft\Credentials` 那一枚 case 的回落支（`aUnder`→`anchorMatches`） | `blacklist.go:237` | **`10／3／0／0 rc=0` 全绿** | `chain:case9-local-creds class=A → class=B needL2=false → true` |
| **M-V2b** 只摘 `~/.ssh/**` 那一枚 case 的回落支 | `blacklist.go:227` | **`10／3／0／0 rc=0` 全绿** | `chain:case4b-ssh-solo A→B`；`chain:case4-ssh` 点名的锚换成 `.aws`（class 仍 A） |
| **M-V2c** 只摘 `~/.kube/config` 那一枚 case 的回落支 | `blacklist.go:231` | `10／2／1／0 rc=1` 红＝ⓑ | `chain:case6-kube A→B` ⇒ `.kube` 是被 ⓑ 钉住的 |
| **M-V3** 反形（变严）：回落支去掉祖先关系界 `return true` | `blacklist.go:366` | `10／1／2／0 rc=1`，**闸门格不在红名册里**，红＝ⓐⓑ＋票 216 那枚 | 九枚 `chain:*` 与九枚 `offchain:*` 全部 `B→A`＝整台机器都算 A |
| **M-V4**＝腿的 M-A 形：`if anchorMatches(… ) \|\| !uncertainAnchorMiss(… )` → `if true` | `blacklist.go:326` | `10／2／1／0 rc=1` 红只 ⓐⓑ，闸门仍绿 | —（reason 面） |
| **M-V5**＝腿的 M-B 形：`anchor := anchorForms(anchorSpelling)` → `pathForms{raw: normPath(…)}` | `blacklist.go:325` | 同上 | — |
| **M-V6** 只把那句"实话"换成**另一种假话**（仍带 `fail-closed` 字样、不提候选与缺位锚） | `blacklist.go:329-330` | `10／2／1／0 rc=1` 红只 ⓐⓑ | —（⇒ reason 尺不是"看见 fail-closed 就放行"） |

M-V1 的红句逐字（`logs/61-v2-MV1-leg.log`，行号＝仪器自己的行）：

```
anchor_miss_reason_334_test.go:193: class: got B, want A (the fail-closed tier is unchanged by AC#1)
anchor_miss_reason_334_test.go:197: Gate: got class=B allow=false needL2=true reason="B-list sensitive path; single-file override requires one L2 confirm: id_*", want class=A needL2=false allow=false (non-overridable)
anchor_miss_reason_334_test.go:201: an A deny from the fail-closed clause was unlocked by an override key: class=B allow=true needL2=false reason="B-list single-file override (one L2 confirm + logged): id_*"
anchor_miss_reason_334_test.go:208: SHAPE-READING334 inside-profile=class=B … || off-profile=class=B …
```

⇒ 那三句是**同一次降级的三面镜子**（`Classify`／`Gate`／豁免键），不是两件事；第 208 行那发读数
在 M-V1 下变成两形同 B，正是"回落支整支被摘"的形状。腿的 `M-C` 打的是同一枚落点（函数首行插
`return false` ⇒ 也是回落支整支），⛔ 打到某一枚 A 档 case 的判定链。

## 3 AC#3 那枚仪器补世界的两问

(a) **成对反控（CI 形台面，我自己真跑）**：

| 发 | 字节 | RUN／PASS_TOP／FAIL_TOP／SKIP（rc） | 红名册 |
|---|---|---|---|
| 改前 | `03f7e7a6`（夹具＝`t.TempDir()`） | `24／5／1／0 rc=1` | `TestFileBaseNameClaimChangesNoClassOrReason` ＋子用例 `B_tier_id_*`／`B_tier_.env*`／`B_tier_secrets.*` |
| 改后 | `3a118226`（夹具＝`candidateOffAnchorChain334(t)`） | `24／6／0／0 rc=0` | `redlist=0` |
| **反控** | `3a118226` ＋ 只把夹具改回 `dir := t.TempDir()` | **`24／5／1／0 rc=1`**，同枚名册同枚用例同三发子用例 | 落点行差逐字 `89c89` 在 `logs/70-ac3-revert-runner.log` |
| **正控** | 同一棵树夹具还原后 | `24／6／0／0 rc=0` | — |

夹具还原哈希 `cbc0a610…`＝交件 blob（逐字同串）。⚠ 我第一发反控把 `GOTMPDIR` 放在 `D:` 卷上跑＝
**台面缺陷**（那一发红本身活在"临时根落在真主目录链内"这件事上，换根＝换判据），已在 `logs/70-*` 用 CI 形重跑。

(b) **那枚新夹具把世界造对了吗**——尺＝`V1PREFIX` 行（候选 root 与每枚锚展开 root 的公共前缀段数）：
候选放在假主目录**的兄弟**（`base/work`）时，七枚锚逐枚 `shared_segments=9`（只到 temp 根，`…\001`）、
`clause=false`（`uncertainAnchorMiss` 用**严格祖先**关系，不是公共前缀）；`offchain:*` 九行全 `class=B needL2=true`
⇒ 锚这一维真的出局，剩下的唯一变量是名字里那句声明。**两臂不同哑**：链外臂在基线里读 B（票 216 要的形），
链内臂（`chain:*`）读 A（票 334 的过严形），333-v1 抓到过的那一形在这里没有复现。

## 4 幻影凭据普查（`logs/50-phantom-census.log`）

尺①＝`git grep -h -oE '^func (Test[A-Za-z0-9_]+)' 3a118226 -- '*_test.go'` → 名册 **1944** 枚；
尺②＝`git ls-tree -r --name-only 3a118226` → **14816** 枚跟踪路径。正控：两枚真名 `在册=1`、`blacklist.go 在树=1`。
待验对象＝三枚产码件的注释＋四枚交付 `.md`（`scan_bytes=66502`，抽出 38 枚引用）逐枚验：

- 指名的 14 枚测试名里 **11 枚在册**；`TestProbe334`／`TestProbe334NecessityMatrix`／`TestProbe334RealMachineAnchors`
  **在册=0** ⇒ 但它们不是"假凭据"：字节真在 `.scratch/wisp/probes/334/r1/logs/probe-334r1_test.go.txt`
  （⛔ 跟踪的 `_test.go`，是腿自报的仓外只读探针），读数在 `logs/*.log`。
- 指名的 `.go` 件逐枚命中；**唯一 NONE＝`internal/tools/rules_gateway.go`**——那是**票面与派单**写的路径
  （真身 `internal/risk/rules_gateway.go:78 case TierA:`），腿已在 ac0 文档第 4 节具名更正过。
- 编排者点名的那枚中途报告引用：`internal/risk/blacklist_test.go` 在 `3a118226` **不存在**、
  `TestApprovalSensitivePathIsNotOverridable` 在该枚树 **0 命中**、**四枚交付件里 0 命中**（只在未入库的中途报告）。
- `non-overridable` 一词在 `3a118226` 的 `internal/**` 命中四枚：三枚腿自报的那三件 ＋ `internal/tools/mode.go`
  （那里是 `Absolute names the A-tier paths: non-overridable, no mode, no grant.` 这句散文，⛔ 用例名）。

## 5 AC#0 六行＋AC#4＋AC#5 我这边的读数

| 项 | 我的读数 | 尺（射程） | rc |
|---|---|---|---|
| ① `pathExists` 刻意偏向存在 | 注释在 `pathresolver.go:427-429`、函数 `:430`（具名树 `03f7e7a6`）；`git log -L 427,430` ＝ `f1033e1e7 2026-09-21` ⇒ 既有设计 | `git show \| awk NR` ＋ `git log -L`（读码，非跑） | 0 |
| ② 那一跳的错串 | `V1LSTAT i=0/1/2 err_is_FS_ErrNotExist=false`，`i=3`（干净缺失件）＝`true` 正控 | `logs/72-probe-v2-before.log`（真跑） | 0 |
| ③ 锚侧缺位 | 真机 `chain` 之外那发：`real-temp-i0..2 class=A rule=A:~-.aws-credentials` | 同上（真跑） | 0 |
| ③′ 补造 `.aws` 换 `.kube` | 我在 `3a118226` 树里用**自己的** solo 世界复现：`chain:case6-kube rule=A:~-.kube-config` ⇄ `chain:case5-aws rule=A:~-.aws-credentials`（⛔ 动机主 profile） | `logs/60-v2-baseline-chain.log` | 0 |
| ④ 必要条件 | 链内 `A` ⇄ 链外 `B`（九对，`offchain:*` 全 B）；` Lad 挪走仍 A` 那发我没重造，腿的读数我不背书 | 我的 `offchain` 行 | 0 |
| ⑤ 今天真实会撞 | `real-appdata-wisp-logs class=A allow=false needL2=false rule=A:%APPDATA%-wisp-config.toml` ＋ `V1OVERRIDE … allow=false`（豁免键买不到放行） | `logs/72/73` | 0 |
| ⑥ CI 会红 | `24／5／1／0 rc=1`（第 3 节那发） | `logs/40-ac3-before-roster.log` | 1 |
| AC#0 加的那枚分母 | 我这把尺具名构造并实测：链内读 A **9 枚** ＋ 真机读 A **4 枚** ＝ **13 枚**；同层干净文件 10 枚读 `none`（9 枚 `chain-clean:*` ＋ 1 枚 `real-temp-i3`）；链外 9 枚读 `B`（⛔ 全称句；尺＝`grep -c 'V1SHAPE …class=' `，射程＝`logs/72-probe-v2-before.log` 全部 32 枚 `V1SHAPE` 行） | `logs/60-v2-baseline-chain.shape` ＋ `logs/72-probe-v2-before.log` | 0 |
| AC#4 消费者普查 | 正控＝`A-list sensitive path` 在 `*.go` 命中 2 枚文件（生产者 `blacklist.go:84` ＋ 仪器自己）；`fail-closed, NOT a match` **只命中生产者 1 枚**；各锚字面串在测试里只做**子用例名／注释**（`pathshape_portable_test.go:134` 那张表比的是 `Classify()==ClassA`，⛔ 比 reason）；`design/doubao/demo/screens/security.js`＝机主外部 mock；`docs/reports/pending-and-issues.md` 里的 `B-list sensitive path`＝台账散文 ⇒ **判定型消费者 0 枚** | `git grep -nF/-lF`，射程＝`3a118226` 全仓跟踪件（剔 `.scratch/**`） | 0 |
| AC#5 门禁四数（只在共享工作树） | `gofumpt.exe -l` 三枚件 rc=**0**（空输出）／`sh scripts/d22scan.sh` rc=**0** `clean - no D22 ban violations`（射程 internal=586 cmd=143 design=39 frontend=116）／`sh scripts/check-path-length-budget.sh` rc=**0** `VERDICT GREEN`（分母现量 **15035**、over-budget 60＝roster 60；腿报的 14818 是它自己那刻的值，不是差异）／`go vet ./internal/risk/` rc=**0** | `logs/81-gates-shared-worktree.log` | 0 |
| 腿对判级面的零字节 | `git diff 03f7e7a6 3a118226 --numstat -- internal/risk/rules_gateway.go internal/tools/` ＝ **空输出** | 现跑 | 0 |
