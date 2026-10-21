# 票 171 落地腿 171-r1 交件（续跑程，第四部分：AC#1／AC#9 的落地账 ＋ AC#3 两条腿）

接 `171-r1-impl.md`／`-part2.md`／`-part3.md`（那三枚由并发腿 `2a5ae70a`／`f936620b` 提交，**本程一字未改**）。
本件＝续跑程（同身份 171-r1，第二班），取数 `2026-10-21 10:02–11:2x +08`。
凡"响／不响"都是当次现跑；量不到的逐条具名。SEC5 之后见 `-part5.md`。

- 本程提交：**`85edab3d`**（两枚文件，`git show --name-only` 逐字见 part5 SEC6）
- 锚：开工 `1cbd025b`；取数期间 HEAD 三度前移（`8288d49a` -> `2a5ae70a` -> `f936620b`），见 part5 SEC7

---

## 1. 过期一：AC#1／AC#9 的落地账（供翻框，未重做）

编排者裁：AC#1／AC#9 的代码早已落地，按"已有落地＋今天复跑 GREEN"结案。**逐笔现量**：

| 提交 | 主题（逐字截取） | 动了什么 |
|---|---|---|
| `f4d64fcc` | `171-r3 step 1: unmodified-code readings for AC#1 (both leg families) landed before the ruler is touched` | `probes/171/r3/gate-prefix-f78d3cb.sh`／`impl.md`／`logs/*`（改前尺快照与读数） |
| `2a5489d3` | `171-r3 step 2: AC#1 numerator covers both leg families (trailing comments stop counting) and AC#9-1 makes the flip door read the stale row's own columns` | **`probes/154/gate-clauses.sh`** ＋ **`probes/161/r6/flip-declaration.sh`** ＋ fixtures ＋ 台件 |
| `7cb9d570` | `171-r3 step 3: evidence for the post-commit re-readings, and the named limit of the confirmatory rerun` | `probes/171/r3/impl.md` ＋ `logs/*` |

故：`probes/154/gate-clauses.sh` 与 `probes/161/r6/flip-declaration.sh` 的改动**真在盘**（编排者结论成立）。
尾随注释那枚洞在 `2a5489d3` 之后由 `pair()` 与 `run()` **两族腿同时**收窄，不是只修一侧。

### 1.1 我今天那枚台件的 rc 与读数（现跑，`mode=post`）

台件 `probes/171/r3/r3-trailing-comment.sh`（AC#1 两族腿，四形）故： **rc=0 GREEN**：

| 形 | 样本 | pair 未成对（未修 -> 修后） | run 命中行数（未修 -> 修后） | 判 |
|---|---|---|---|---|
| S1 | 词根只在**行尾注释**里 | `1` -> `0` | `1` -> `0` | 该剔的剔掉了 |
| S2 | 真调用点（**正控**） | `1` -> `1` | `1` -> `1` | 该留的留住了（不是把尺子弄瞎） |
| S3 | 真开、关方只在行尾注释 | `0` -> `1` | `2` -> `1` | 未修**静默**（漏），修后**点名** |
| S4 | 词根只在整行注释 | `0` -> `0` | `1` -> `0` | run 那一侧的同族洞一并收掉 |

故 AC#1 的判据**两向都成立**：该不计的不计（S1／S4）、该计的照计（S2）、
危险方向（S3：一条注释买来一个"关"）从静默变成点名。

台件 `r3-flip-columns.sh` 与 `r1-flip-columns.sh`（AC#9-1 两向对照）：**本程未取到今日读数**。
成本具名：整门 11 轮 × 单轮约 7m20s 约 80 分钟，两扇门对照约 4 小时；本程两次启动都只跑到中途
（kill 记录见 part5 SEC7.2）。故 AC#9 的**代码落地**有 `2a5489d3` 为证
（`flip-declaration.sh` 里 `the row prints its own numbers` 那一味由 `Select-String` 数得 3 处、
改前 `2a5489d3^` 为 0 处），但**今天的两向读数我没跑出来**，不拿别腿读数冒充。

---

## 2. 过期二：AC#3〔甲〕两条腿都做完（本程主活）

### 2.1 腿一：`.scratch/wisp/probes/161/r5/attrib.sh` 加名册断言（`+114/-4`，552 -> 662 行）

三枚断言，**分开钉**（合成一个总数就允许拿一枚 ring 换一枚 silent 而读数不变）：

```
ST_RING_EXPECT=5          # (1) 必须响的用例登记数
ST_SILENT_EXPECT=3        # (2) 必须静默的用例登记数
ST_RING_NAMED=".scratch/wisp/probes/999/bad-sample.go
.scratch/wisp/probes/777/second-fake-ticket.go
internal/tools/not_a_bench_path.go"   # (3) 具名必须在场
```

`st_case()` 里按 `want` 分别累加 `ST_RING_SEEN`／`ST_SILENT_SEEN`，并把喂进去的原文行攒进
`ST_RING_LINES`；收尾处三枚**各自独立判定**（不是 `elif` 链——链会让"计数"先命中就吞掉"点名"，
正是 275-v1 `:313` 记的"红而无名"那一形）。缺一不可的理由：
(1) 删一枚 ring 会响，但只说"4 不是 5"说不出**缺谁**；(2) 是它的镜像（只钉总数会让两向互换）；
(3) 抓 (1) 抓不到的那一形——**计数不变、把被删的样本换成另一枚**。

### 2.2 删行变异：三读数（本格的牙，全部现跑）

**变异落地证明**（不是只说不跑）：删的是 `st_case ... 777/second-fake-ticket.go` 那一行，
`diff` 逐字 `345d344`／`<     st_case 1 UNATTRIBUTABLE '.scratch/wisp/probes/777/second-fake-ticket.go'`
（**1 removed / 0 added**；名字登记一字未动——登记还在、样本没了，这才是要抓的形）；
`sh -n` 故： **rc=0**（不是因为语法错而红）。台件 `r1-roster-nails.sh` ＋ 生成器 `r1-roster-nails-mutants.sh`。

| 读数 | 命令 | rc | 读数逐字 |
|---|---|---|---|
| **红** | 变异体 `--self-test` | **1** | `FAIL  roster assertion  want ring=5 silent=3 named=3  got ring=4 silent=3`；`ring-count: registered must-RING cases=4, this file declares 5`；`named-missing: a NAMED must-ring sample is not among the cases that ran`；`missing must-ring sample(s), verbatim:` 后逐行 `.scratch/wisp/probes/777/second-fake-ticket.go`；`cases=9 failures=1` |
| **还原绿** | 原尺再跑 | **0** | `cases=10 failures=0` |
| **原样绿** | 同一发 | **0** | `PASS  want ring=5 silent=3 named=3  got ring=5 silent=3  all named must-ring samples present` |

**抗装饰读数（这一发必须今天不响）**：HEAD 那版（552 行、无名册断言）跑**同一发删行** 故：
**rc=0、`cases=8 failures=0`、SELF-TEST GREEN**。故洞今天真在；本程不是给一把已经会响的尺子补记录。

**另两枚变异各咬各的**（同一台件现跑 **rc=0 GREEN**）：
换名但计数仍 5 故： `cases=10 failures=1`，**只有** `named-missing` 响（证明 (1) 单独不够）；
删 silent 那枚 故： `cases=9 failures=1`，**只有** `silent-count` 响（2 登记 / 3 声明）。

### 2.3 腿二：`ci.yml` **只追加**一步（`+48/-0`，零删除行）

```yaml
      # ADDED STEP ONLY (ticket 171 AC#3, leg 171-r1): the attribution instrument's
      # OWN SELF-TEST, run in CI. No existing step was moved, edited, deleted, given
      # an `if:` it did not have, or made continue-on-error.
      # ...（46 行理由：这一步本来缺什么、为什么名册必须先落地、为什么它会真红）
      - name: "Attribution instrument self-test with its roster nails (ticket 171 AC#3 - the mutation goes red)"
        if: ${{ !cancelled() }}
        run: sh .scratch/wisp/probes/161/r5/attrib.sh --self-test
```

**为什么它不会让既有 job 变松**（逐条，读字节／解析 YAML 得出，不是 grep 注释）：
1. `git diff --numstat` 故： `48  0` 故： **纯新增**；删除行数 `grep -c "^-[^-]"` = **0**。
2. **无** `continue-on-error`（`yaml.safe_load` 后查键，不是查注释里有没有这个词）。
3. `if: ${{ !cancelled() }}` 是该 job 里 `staticcheck`／`mockllm module vet`／
   `portable-tests-selftest` 已在用的同一枚 guard——它改的是"这一步报不报"，**不是**"能不能红"。
4. 位置＝lint job **末尾下一个空位**（step 14／14），没有越过任何既有步骤。

### 2.4 步骤名名册作差（改完后**重新快照**再作差）

改前 `40` 枚 -> 改后 `41` 枚；**消失 0 枚**、**新增 1 枚**（逐名两向 `comm`）：
`"Attribution instrument self-test with its roster nails (ticket 171 AC#3 - the mutation goes red)"`。
六个 job 步数：`lint 13->14`；`lint-frontend 11->11`／`slo-full 5->5`／`slo-smoke 6->6`／
`test-core 7->7`／`test-windows 10->10` 全 SAME；**job 名册与触发表一字未动**（含 `schedule` cron）。
台件 `probes/171/r1/ci-delta.py` 故： **rc=0**，末行
`GREEN - one step added, nothing removed, triggers untouched`。

---

## 3. AC#9 的补齐：表到底被核到几分（只报我量到的）

| 尺／点 | 我今天的读数 | 出处 |
|---|---|---|
| `flip-declaration.sh` 的"列核"味 | `the row prints its own numbers` **3 处**（改前 `2a5489d3^` **0 处**） | `Select-String` 现跑 |
| 那味的判据形 | `want_row="基线=${new_n}枚 实测=${m0}枚 差=$((new_n - m0))枚"`，与实到行**全文**比而不只比存在 | 读码 `:202`-`:223` |
| 两向对照（造假行／诚实行） | **本程未跑完**（成本见 SEC1.1） | — |
| `r3-sample-deleted.sh`（删样本不删腿） | **本程未取到**（同上） | — |

故 AC#9 我只能给"代码在哪、那味有几处、它的判据形是什么"，**给不出今天的两向读数**；缺口具名，不拿别腿顶替。

---

## 4. 门禁（本提交所在这棵树现跑，逐把 rc ＋末行逐字）

| 门 | rc | 末行逐字 |
|---|---|---|
| `GOFLAGS= go build ./...` | **0** | （零输出＝clean） |
| `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=244, bans #1-5 cmd/=42, ...` |
| `sh scripts/check-path-length-budget.sh` | **0** | `... VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree` |
| `attrib.sh --tracked-only` | **1** | `RED (tracked-only) - tracked-dirty=35 files_dirty=35. ...`（照实报，未洗） |
| `attrib.sh --self-test` | **0** | `--self-test cases=10 failures=0`（原 **9**，多出的那一枚即名册断言） |
