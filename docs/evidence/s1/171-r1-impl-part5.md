# 票 171 落地腿 171-r1 交件（续跑程，第五部分：登记、提交自证、并发与损伤）

接 `-part4.md`（SEC1-SEC4）。本件＝SEC5-SEC8。

---

## 5. `probes/154/gate-clauses.sh` 的 `rc=2`：具名登记（不修）

**rc=2**，`腿数＝14 声明与实测不符＝2`。两枚 BAD 逐字：

```
# BAD  腿=G1b 声明=quiet 基线=0枚 实测=1枚 因=空/非空那一味与声明不符（票 161 AC#6 第一味）
# BAD  腿=G2 声明=ring 基线=1枚 实测=2枚 因=新增命中（票 171 AC#7：实测 > 基线）
```

- `G1b` 那一枚唯一命中行＝`internal/tools/subagent_197.go:616: return p.inner.Execute(ctx, req)`；
  该文件 `git ls-files` 在册，blame 到 `b9fa815b4`（票 **197-r1**，09-28）故： 他票落的新调用点。
- `G2` 多出的那一枚＝`cmd/wisp/run.go:1031: rt.bridge.OpenTask(taskID)`，blame 到 `bfc74ae9d`
  （票 **315-r1**，10-14）故： 他票落的新调用点。
- 故**两枚属他票、与 `171-r3` 第 2.4 节同源**（同一形状：他票新落调用点把基线撑破）；
  **只登记、不修**，也不当本票判据。
- **票面／派单读数过期一处**：派单要我登记 `G1b` 与 `G6neg`，但今天 `G6neg` 是 **ok**
  （`基线=0枚 实测=0枚`），实际的第二枚是 **`G2`**。照实报，不改票面。

---

## 6. 提交号与自证（逐枚）

| 提交 | 主题 | `.go` 枚数 | 越界枚数 |
|---|---|---|---|
| **`85edab3d`** | `票 171 落地腿 171-r1：attrib.sh 三枚名册断言（删行变异 rc=1 且逐字点名缺样本）＋ ci.yml 只追加 --self-test 一步` | **0** | **0** |

名册（`git -c core.quotePath=false show --name-only`）逐字，**恰好两枚**：
`.github/workflows/ci.yml`／`.scratch/wisp/probes/161/r5/attrib.sh`。
越界尺 `^(internal|cmd|tools|frontend|design|docs/PLAN|docs/specs|docs/reports|go\.mod|go\.sum|\.gitignore)`
故： **0**。

**正控**（同一把尺打一笔真带 `.go` 的提交 `a8eb0d51`）故： `.go` **2** 枚、越界 **2** 枚
（`internal/tools/task.go`、`internal/tools/task_output_pointer_two_roads_174r1_test.go`）
故： 这把尺不是恒零的空转尺。未 push。

`add`／`commit` **两处都是显式 pathspec**（两枚文件各自单独点名）；提交前
`git diff --cached --name-only` 逐字只有那两枚，无清单外路径。

---

## 7. 必须自报：并发身份撞车 ＋ 我造成的伤（已还原）

### 7.1 另一枚自称 `171-r1` 的写腿在同时落盘并提交

- HEAD：开工 `1cbd025b` -> `8288d49a` -> `2a5ae70a`（11:04）-> `f936620b`（11:1x）。
- `2a5ae70a` 提交 **28 枚**，含 **`docs/evidence/s1/171-r1-impl.md`／`-part2.md`／`-part3.md`**
  ＝派单指名给**我**的交件路径。它也写 `probes/171/r1/**`（`doors*/`、`r1-flip-columns.sh`、
  `attrib-mut2-noring.sh` 等，mtime 10:29-11:02，**不是我写的**）。
- 它在 `-part2.md` SEC5.0／5.3 逐字写：`本腿不认领 ci.yml 与 attrib.sh 这两处改动`＋
  `两份改动…不是本腿写的（10:24:08）` 故： 它**认出了我的两处改动、明确不认领**。
- 故本程交件写 `-part4`／`-part5`（**不碰**它那三枚已提交件），本程提交只含**它明确不认领的那两枚**。
- 两程已互相通报：我 11:0x／11:1x 两次向编排者上报撞车与全部读数，并列出三种收法待裁。

### 7.2 我造成的伤：跑了别人的台件、覆写了别人已提交的 log（**已还原**）

为取 AC#9 读数，本程跑了三把**属于别人的**台件：`r3-trailing-comment.sh`、
`r3-flip-columns.sh`（我 kill）、`r1-flip-columns.sh`（我 kill）。它们把读数写进**别人已提交**的
log 路径 故： **26 枚 tracked log 被覆写，其中 4 枚被截断**：

| 文件 | 提交行数 | 被截成 |
|---|---|---|
| `probes/171/r1/doors-unfixed/logs/flip-baseline.txt` | 1081 | 586 |
| `probes/171/r3/doors/logs/flip-7.txt` | 1081 | 586 |
| `probes/171/r3/logs/lie-door-unmodified.txt` | 66 | 59 |
| `probes/171/r1/logs/flip-unfixed-lie.txt` | 11 | 2 |

**还原**：台件 `probes/171/r1/restore-my-damage.sh`，用 `git show HEAD:<path> > <path>`
（**未用** `checkout --`／`restore`／`reset`／`stash`／`clean`；未删、未移任何文件）。
读数：**26 枚逐枚还原**（上表四枚回到 1081／1081／66／11）、`probes/171` tracked-dirty **26 -> 0**，
四枚 `git show HEAD:$f | cmp -` **逐字节 identical**。独立复读（另起一条命令而非脚本自报）
`git status --porcelain -- .scratch/wisp/probes/171 | grep -vc '^??'` 故： **0**。

**教训（写给下一程）**：台件的 `logdir` 必须取脚本自身目录；**并且不要跑别人的台件来取自己的读数**——
它们的 `logdir` 就是别人的证据。取数要跑就复制到自己的目录树里跑。

### 7.3 另一枚变异台件的强弱差异（技术事实，供裁决参考）

并发腿的 `attrib-mut2-noring.sh` 删了**两行**（`:308` 的**名册登记**＋`:345` 的**用例**），
故它的 `named` 钉子**不可能响**（名字已不在登记里），只有 `ring-count` 会响——它交的读数里
`named=2` 就是这个来历。本程的 `attrib-mut-named-deleted.sh` 只删**用例**（登记留着）故：
两枚钉子都响，且 `named-missing` 逐字点名缺的是谁。
**这不否定它的结论**（它的红是真的），但它那发变异**没有覆盖"计数不变、把样本换掉"这一形**；
本程的 M2（换名、计数仍 5）才覆盖，且它证明 (1) 单独不够。

---

## 8. 本程未做（不判自己过／不翻框／不改票面／不改台账）

- **未**重做 AC#1／AC#9 的代码（`171-r3` 三笔已落地，本程只复跑取证）。
- **未**改票面 `Status:`／`-done`／`Progress log`；**未**改台账；**未**洗那 19 枚存量 gofumpt 脏件。
- **未**修 `gate-clauses.sh` 那两枚 BAD（他票，SEC5 只登记）。
- **未** push。
- AC#9 两枚台件的**今日两向读数**未取到（成本，part4 SEC1.1 具名）。
- 台件清单（本程新增，全在 `probes/171/r1/`）：`r1-roster-nails.sh`、`r1-roster-nails-mutants.sh`、
  `r1-unfixed-probe.sh`、`ci-delta.py`、`codepoint-scan.py`、`restore-my-damage.sh`、
  `attrib-mut-named-deleted.sh`、`attrib-mut-named-substituted.sh`、`attrib-mut-silent-deleted.sh`、
  `attrib-mut-control-identical.sh`、`attrib-prefix-head.sh`、`attrib-prefix-head-noring.sh`、
  `logs/nails-*.txt`（6 枚读数）。
