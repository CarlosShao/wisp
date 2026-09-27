# 161-v1 非实现者验收表：票 161 的 AC#1／AC#2／AC#3 三格（对抗验收，逐格独立读数）

- 派单＝`.scratch/wisp/dispatches/2026-09-27-092x-accept-161-v1-selftests-ac1-ac2-ac3.md`（09-27 09:2x）
- 本程性质＝**非实现者裁决**。三格的实现者分别是 161-r1（AC#1）／161-r3 写码＋编排者代提复跑（AC#2）／161-r2（AC#3）。
- 锚＝**`8efece378a38e74bd95f6cf7aabd971fb3dfc3eb`**（`git rev-parse HEAD` 现量，`git cat-file -t` = `commit`，2026-09-27 09:20）。
- ⚠ 并发事实：另程此刻在写 `probes/154|161/**`；本程**全部被验读数取自仓外副本**——
  `git ls-files` 全清单（2385 枚）`tar -T` 到 `/tmp/wisp161v1/tree`，再把 22 枚工作树相对锚有出入的已跟踪路径
  逐一 `git show <锚>:<path>` 复写回锚字节，副本内 `git init`＋全量 add/commit 建索引。
  副本内 `git ls-files '*.go'` = **531**（与 r5 的甲形分母同数）。锚到 HEAD 之间 `tools/d22scan`、`.github/workflows/ci.yml`、
  `docs/evidence/s1`、`internal`、`cmd` 的 `git diff --stat` 为**空**。
- 本机仪器：`go1.27.1 windows/amd64`、`gofumpt v0.12.0 (go1.27.1)`（两枚都现跑 `--version`，日志 `probes/161/v1/logs/go-version.txt`）。
- 档位约定：本件默认〔我本轮现跑过〕；凡引用实现者的话都显式降档。

---

## 0. 本程**没**测什么／哪些读数受并发影响（先写这节）

1. **乙形（工作树归因）没在真树复跑**：真树 `probes/**` 正被另一程写，任何"此刻跑 `gofumpt -l .`"的读数都不属于任何锚。
   本程只在仓外锚副本上跑了**甲形**（§2.4）。⇒ 编排者 09:0x 那句"乙形 9 行全归票 161、rc=0"的档位在本件里是
   〔仅自述，我未复算〕——要钉它需要一次安静窗口的真树现跑。
2. **CI runner 一次没跑**：GitHub 上任何一步的真实 conclusion 没量过；§2.5 的可达性是**按 ci.yml 锚字节的作业/步骤结构**现读的。
3. **B3 盲区（bans #1–#5 不看 `_test.go`）没独立复算**：派单只要求 ≥2 枚盲区复算，我做了 B1/B4/B5/B6 四枚＋B2 的扩张取样（§1.3），
   B3 属"未定义即停"族，留给 owner 的话。
4. **r2 的尺 B/C/D 本体（AST 尺、链接器尺）没重跑**：AC#3 我换的是**自己的 grep 尺**（独立于它的四把），它的尺只抽验了 file:line 对得上。
5. 门铃的**运行期端到端**（真起进程看卡）没测——那本来就是 r2 §2 N1 自报未测，本格判据不含它。

---

## 1. AC#1 裁决：〔成立〕——"8 枚全点得到＋六枚盲区登记"经独立复算成立；本件另登记两枚 r1 未取样的相邻形状

### 1.1 名册现读（锚字节，不抄话）

`tools/d22scan/main.go:5–49` 的 `// Bans` 块：编号禁令 **8 枚**（bare-goroutine / pathresolver-bypass / plaintext-key /
wallclock-timeout / mirror-hash / panel-approval / internal-artifact-tool / emoji），`unparseable` 是发现类型不算分母。
⇒ 与 r1 §0.5、票面 09-26 21:3x 更正一致。**"九条"作废**这一条成立。

### 1.2 我自己的新样本（27 发之外、r1 名册没出现过的形状）

台件＝`.scratch/wisp/probes/161/v1/{samples,logs}/**`；跑法＝副本内 `go build` 出的 `d22scan-v1.exe -root <假根>`，
四枚假根在 `/tmp/wisp161v1/fakeroots/{rootA..rootD}`（不进生产树）。**假根下限读数先撞了一发好门**：
rootA 骨架（3 枚 .go）被 main 的"production .go < 10 ⇒ rc=2 拒答"守卫挡下 ⇒ 我的"干净基线"对照由
**同跑同 tag 的必须响样本**承担（见每格控制列），这一发顺带证明 -root 侧也不存在"递给尺 3 枚文件报 clean"的恒真形状。

| 发 | 形状（新造的） | 预期 | 实测（rootB rc=1，finding 逐条在 `logs/v1-rootB.txt`） | 控制（同跑同 tag） |
|---|---|---|---|---|
| s1 | 方法值裸协程 `go s.worker()` | 响 | **响**（`bare \`go s.worker(...)\` is banned ... R16: named calls count too`） | — |
| s6 | 裸 `filepath.Abs/Clean` | 响 | **响**（两臂各一条） | — |
| s7 | `const probeAPIKey = "PROBE161V1..."` | 响 | **响** | — |
| s8 | 同行 mirror+sha256 | 响 | **响**（2 条） | — |
| s9 | `now.Sub(time.Now())` | 响 | **响** | — |
| s5 | 三行三字形：`→`／`①`／`✓` 各占一行 | 只 `✓` 响 | **只 `s5.go:7`（✓ 行）响**，3/5 行静默 | ✓ 行响 |
| s2 | `import fp "path/filepath"` + `fp.Clean` | （盲区）不响 | **不响** | s6 同跑响 |
| s3 | `time.Now().Sub(deadline)`（反向操作数） | 我认为它漏 | **不响** | s9 同跑响 |
| s4 | map 字面量 `"apiSecret": "PROBE161V1..."` ＋ `cache["apiKey"] = "..."` 索引赋值 | 我认为它漏 | **不响** | s7 同跑响 |
| unp | rootD：故意不解析的 .go 内含 `go func(){` | 只响 unparseable | **该文件只有 `[unparseable]`**，全树 `bare-goroutine` 计数=1（是 s1.go 的，不是它的） | s1 同跑响 |
| al | rootC：与 rootB 同文件集＋allowlist 条目 `pathresolver-bypass→internal/` | 被压掉且零打印 | **pathresolver-bypass 0 条**、bare-goroutine 照响、全输出 `grep -ci suppress`=**0** | rootB 同文件响过 |

**r1 表格的复算结论**：八枚编号禁令里我在假根上亲手打响 6 枚（#1 #2 #3 #4 #5 #8）；#6/#7 两枚经 §2 的
`-self-test`（AC#2 表，本程现跑 rc=0，#6 两响一静、#7 两响一静）在我本轮过了一遍 ⇒ **"8 枚今天都点得到"经独立路径成立**。
登记成盲区的六枚里我独立复算 **B1（s2）／B4（s5）／B5（rootD）／B6（rootC）四枚，全部与其读数一致**；
B2 被我的 s4 扩张（map 字面量与 IndexExpr 左值赋值同族两形今天都不响——r1 §2 N3 自己标了"没取样"，这一格从"没量"变成"量过、静默"）。

### 1.3 本件新增登记（r1 表格里没有的两枚，不算退回、算续表）

- **V1-N1** ban #4 反向操作数：`time.Now().Sub(deadline)` 不响，而禁令文字是"timeout/deadline 逻辑里的墙钟差"。
  会响条件（若哪天扩它）：上表 s3 翻响。今天它是**登记外**的相邻形状，r1 的"射程比文字窄"结论因它**更强**而不是被推翻。
- **V1-N2** ban #3 的 map 索引赋值：`cache["apiKey"] = "<24 字符>"` 不响（左值是 IndexExpr，AST 两支都不看）。
  同族同因（只认 ValueSpec 与 Ident 左值），会响条件同 B2 那条。
- 顺手一枚反向哨兵（不是洞）：s1 的 `go s.worker()` 证明 r1 §2 N3 列的"没量的相邻形状"之一今天**在射程内**。

**恒真/空心自查**：表内每一枚"不响"都有同跑同 tag 的"响"作控制（见控制列）；rootA 的骨架下限由 main 的 ≥10 文件守卫硬拒。
⇒ 本程没有一枚"永远不响的检"，也没有一枚"递给尺 0 枚文件报干净"。
