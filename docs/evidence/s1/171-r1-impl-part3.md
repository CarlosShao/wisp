# 票 171 落地腿 171-r1 实现记录（第三部分：门禁、写面、码位、边界）

接 `171-r1-impl.md` 与 `171-r1-impl-part2.md`。本件＝SEC7-SEC11。

---

## 7. 门禁（本腿自己现跑，逐把贴末行）

| 命令 | rc | 末行／要点 |
|---|---|---|
| `GOFLAGS= go build ./...` | **0** | 无输出 |
| `sh scripts/d22scan.sh` | **0** | 全通过 |
| `sh scripts/check-path-length-budget.sh` | **0** | 全通过 |
| `sh probes/154/gate-clauses.sh`（默认锚 HEAD） | **2** | `# 腿数＝14 声明与实测不符＝2`／`# 基线过期枚数＝0`／`# 聚合退码＝2` |
| `sh probes/161/r5/attrib.sh --tracked-only` | **1** | `attrib.sh: RED (tracked-only) - tracked-dirty=35 files_dirty=35.` |
| `sh probes/171/r1/r1-two-shots.sh both` | **0** | `r1-two-shots.sh: GREEN - mode=both, every assertion above held.` |
| `sh probes/171/r1/r1-flip-columns.sh both` | **0** | `r1-flip-columns.sh: GREEN - mode=both, every assertion above held.` |
| `sh probes/171/r3/r3-trailing-comment.sh post` | **0** | `GREEN - mode=post` |
| `sh probes/171/r3/r3-sample-deleted.sh` | **0** | `GREEN.` |

`gate-clauses.sh` 那两枚 `BAD`（`G1b`、`G6neg`）**不是本票的读数**，是别程落的新调用点
（与 `probes/171/r3/impl.md` 的 2.4 节同源）。本腿**没有**为它们改任何 `want`／`want_n`。

---

## 8. 写面自报（越界＝本件作废）

本腿**新建**：

```
.scratch/wisp/probes/171/r1/r1-two-shots.sh
.scratch/wisp/probes/171/r1/r1-flip-columns.sh
.scratch/wisp/probes/171/r1/fixtures/t1_only_trailing_comment.go.txt
.scratch/wisp/probes/171/r1/fixtures/t2_comment_plus_real_violation.go.txt
.scratch/wisp/probes/171/r1/fixtures/t3_real_violation_no_mention.go.txt
.scratch/wisp/probes/171/r1/fixtures/gate-lie-columns-r1.sh
.scratch/wisp/probes/171/r1/gate-unfixed-2a5489d3-parent.sh
.scratch/wisp/probes/171/r1/flip-unfixed-2a5489d3-parent.sh
.scratch/wisp/probes/171/r1/attrib-mut-noring.sh
.scratch/wisp/probes/171/r1/attrib-mut2-noring.sh
.scratch/wisp/probes/171/r1/doors-unfixed/  probes/171/r1/doors-fixed/  probes/171/r1/logs/
docs/evidence/s1/171-r1-impl.md
docs/evidence/s1/171-r1-impl-part2.md
```

**本腿一字节未碰**：`internal/**`、`cmd/**`、`tools/d22scan/**`、`docs/PLAN.md`、`docs/specs/**`、
`go.mod`／`go.sum`／`.gitignore`、`frontend/**`、`design/**`、任何阈值／golden、`.gitattributes`、
`probes/154/**`（只读）、`probes/161/r6/**`（只读）。
`.github/workflows/ci.yml`：开工后本腿**零字节**改动（改动者为并发腿，见 SEC5.0）。

---

## 9. 禁用码位自检

本两枚证据件按派单要求做码位自检：`U+00A7`（节号）、`U+2713`（勾）、`U+2264`（小于等于）、
`U+2460`-`U+24FF`（带圈数字）、`U+1F000`-`U+1FAFF`／`U+2600`-`U+27BF`／`U+2B00`-`U+2BFF`／`U+FE0F`（emoji）
——本件正文**零命中**；需要写节号处一律写成"SEC"的名称形（`U+00A7`），需要写序数处一律用阿拉伯数字。

---

## 10. 本腿没做／判不动（具名）

- **AC#3 的两处落地**不在本腿名下（SEC5.0、SEC5.3）；本腿只给独立复验读数。
- **票面 AC#1／AC#9 与票面文字不一致**（第一部分 SEC0.3 过期甲）：票面 `Progress log` 缺 `171-r3`／`171s`
  两行，故票面读起来像"未闭"。本腿**没有**替票面补那两行（禁改台账／禁改票面）。
- **`probes/171/r3/**` 22 枚 ` M`**：是本腿跑 `r3-*` 台件时那些台件**自己覆写**的日志
  （`logdir` 取脚本自身目录的既有设计）。**未还原、未删除**，归属待编排者裁。
- **`ci.yml` 那步没有 CI 现场读数**：本腿不能 push，也没触发真实 run；SEC5.2 三行是**本机**现跑。
  要在 CI 上看见它绿，得等编排者推送后的那一 run。
- **`--self-test` 的 `cases=9` 会随名册变化**：并发腿加了名册钉子后它变成 `cases=10`。
  这不是过期，是同一枚尺的枚数随登记增长；本件两处读数都带时刻与来源（改前 552 行／改后 662 行）。
