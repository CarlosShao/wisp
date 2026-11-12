# 346-a2 - AC#2 后半：去牙反形（变异体）＝假绿直证

变异只打仓外拷贝树（`mut-jia`/`mut-yi`）的**测试件物理字节**；⛔ `-overlay`、⛔ 动真实仓库、⛔ 放宽/删红句文字本身（红句连同参数表整段保留在文件里，只是它的守卫不再触发）。

## 变异形状（逐字 literal 替换，各一处/两处）

- 甲（`mut-jia`，两处，同一条规则）：`judgeRun120` 内两支 `if sealWindowAssertClosed120() {` ⇒ `if false { // MUTATION-346a2-jia: abstention guard removed (anti-form, exported tree only)`
  ——覆盖 `:225`（racer 翻面报错支）与 `:233`（"refused==0 || sealedOwn==0" 死仪器支）两句弃权红的**触发条件**；
  `else` 的 Logf 读数支原样保留（默认档反而变为无条件打印）＝窗口未到达的读数仍在屏上。
- 乙（`mut-yi`，一处）：`if liveness == 0 || flips == 0 {` ⇒ `if false && (liveness == 0 || flips == 0) { // MUTATION-346a2-yi: abstention guard removed (anti-form, exported tree only)`
  ——`:206` 弃权红的触发条件摘掉；`damaged>0` 与 `landed>0||narrowed>0` 两支真损伤红**一字未动**。
- 编译面：两棵树各 `go vet ./internal/winsec/`（容器）＝rc=0。

## 哈希（变异前＝pristine 值；变异后）

| 文件 | 变异前 | 变异后 |
|---|---|---|
| mut-jia `posix_seal_window_120_other_test.go` | `e8223d01f05629af…1311a` | `5fcebeec4601f1dde2ca65c82f057f470e51739d880aee28fc4fd61841795e58` |
| mut-yi `openleaf_322_other_test.go` | `0e6aba36ed64849d…4774` | `f0f3751689ea1e88ed5e5fd9abe72e8468e65ab75c9f8eeb06deacf6ee3d29b6` |

## 成对读数（同一条件、同一名册、`-count=1 -timeout 30m`）

**甲**（条件＝容器 `-cpu=1`＋`WISP_120_ATTEMPTS=1`）：

| 发 | 树 | 总量行 | 判定行 | rc |
|---|---|---|---|---|
| 1..3 | export | `attempts=4 … refused-at-check=0 … sealed-own-object=4 other=0` | `RED instrument is dead shape=leaf-swap: refused=0, sealed-own=4 - a zero window reading here proves nothing` | 1,1,1 |
| 1..2 | mut-jia | **逐字同上**（attempts=4, refused=0, sealed-own=4） | `AC#1 INSTRUMENT DEAD reading shape=leaf-swap: refused=0, sealed-own=4`（Logf，不响红） | 0,0 |

⇒ 同一读数、同一未到达窗口：有牙＝红、去牙＝绿。假绿直证 2/2。

**乙**（条件＝容器 `-cpu=1`）：

| 发 | 树 | paired reading（节选 flips / leafLegRefusals） | rc |
|---|---|---|---|
| 1,3,4 | export | flips=16/82/106, leafLegRefusals=**0** | 1（弃权红 :206 逐字在场） |
| 2,5 | export | flips=96/140, leafLegRefusals=1 | 0（窗口到达＝真绿） |
| 1..5 | mut-yi | flips=20/10/10/62/22, leafLegRefusals=1/**0**/**0**/**0**/**0** | **0×5** |
| 对照 | mut-yi | 安静档 flips=1076, leafLegRefusals=257 | 0（真绿） |

⇒ mut-yi 的 shot2/3/4/5 是**严格配对的假绿**：`leafLegRefusals=0`＝窗口确实没到、与 pristine 红发同一形状，但变异体不红＝PASS（4/4）。
shot1（llr=1）是真绿、对照发也是真绿＝变异只摘了弃权那一支、⛔ 碰真损伤红（`damagedWitnesses=0` 三支全零，全 5 发无任何 RED 行——
抽查尺对 5 发 raw 日志大小写敏感 grep `RED instrument|RED:`＝0 命中；01 件里那些 "paired/answered" 命中是尺大小写不敏感的自伤，已弃）。

## 这一笔点名了哪一维

- 乙维＝"叶腿一次答不上（liveness=0）时，仪器必须自陈红"——摘牙后该维对 `TestAC1POSIX…322r1` **失明**（一把两种形状都放行的尺）。
- 甲维＝"racer 打不出任何状态翻转（refused=0 或 sealed-own=0）时，仪器必须自陈红"——摘牙后同维失明。
- 真实仓库侧三句红句逐字仍在（判据自检，尺与读数在 04 件）：`openleaf_322_other_test.go:206`、`posix_seal_window_120_other_test.go:225`、`:233`。
