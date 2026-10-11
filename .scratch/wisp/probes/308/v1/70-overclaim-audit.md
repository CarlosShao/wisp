# 308-v1 · 70-overclaim-audit（盘上句子读起来像"这一族⛔ 会再发生"吗）

尺＝对票 308／`probes/308/{a1,r1}/**`／三枚交付件逐句扫 `prevented|never again|杜绝|根除|再也不会|cannot panic|fully fixed|guarantee`（件：本腿扫描输出，命中＝**0**）＋逐枚读交付注释与两份 90-final。

## 结论先说

**没找到一句把"这一族"写成"被防住了"**。交付面用的动词全是**当输入喂进来时如何处理**那一形：
`driveShare` 文档＝"returns the admin-share token … **WITHOUT** the `vol[:1]` slice that panicked"；新用例文档＝"It feeds … a volume-less **LITERAL** as INPUT"；裁语 ⓑ 复述＝"the guard's teeth are on the INPUT layer … **NOT** on a machine that happens to hand out a volume-less temp"。
`308-r1/20-anti-form.md:19` 甚至正面写着"⛔ 挂在机器台面：本机三把 temp 台面⛔ 能现场触发无卷标那形"。⇒ ⛔ 退回（⛔ 触发"AC 声称要防的结局被造出来"那一族）。

## 但三句要具名（按危险度排）

| # | 位置（内容锚） | 逐字 | 本腿读数 | 处置 |
|---|---|---|---|---|
| 1 | `pathresolver_anchor_spelling_windows_test.go:100-101`（`TestClassifySpellingInvariance` 文档段末） | "every other spelling is **unconditional**, so the test **never quietly stops measuring**." | ★**此刻为假，且比票 §现量 登记的那版更假**：票登记的是"在 `vol==""` 那形下这半句为假"；本票落地后 `unc-admin`/`unc-extended` 变成**条件拼法**（`driveShare` 报 ¬ok 就去掉），于是"every other spelling is unconditional"这半句**在任何台面上都假**——被同一枚函数下面 8 行的代码正面否定。第二半"never quietly stops measuring"在**可达台面**上为真（分母 9⇄9 实测，见 10 件 §4），在**⛔ 可达台面**上只是"⛔ 再 panic"这一半有保证 | **要它自己那一刀**（票面早写了"登记、⛔ 顺手改对"，`Q-308-r1-C` 由此而来）。建议替换句（⛔ 本腿⛔ 改，见 90 件 §3） |
| 2 | `…anchor…:108-110`（ticket 308 注释块） | "the two UNC spellings cannot be built, so volumeInvariantSpellings drops them BY NAME and **the remaining seven volume-free spellings still run** -- the case ends with a verdict" | 量的部分＝helper 在⛔ 卷标**输入**下确实交回 7 枚名字（`logs/v45`、`logs/v-target-case.txt`）。⛔ 量的部分＝"still **run**"＋"**the case** ends with a verdict"：那两枚 windows-tagged 用例的那一支**到达数＝0**（哨兵命中 0，见 20 件 (c) 族），"still run"是从 helper 推到用例的**外推**，本机⛔ 能验 | ⛔ 判假句（推论合理、措辞超出证据）；要一句收窄＝"the helper keeps the seven volume-free spellings in the map; only the selection half is pinned by `TestVolumeInvariantSpellingsVolumeLess308`（the windows-tagged drop branch has no fixture on any reachable surface）" |
| 3 | `…volumeless_spelling_308_test.go:19` | "**This is the single place** the drive letter is taken." | 射程歧义：本腿普查（`grep -rn --include=*.go 'vol\[:1\]\|drive\[:1\]' internal cmd tools`）＝求值位**只剩 2 枚**：产码 `pathresolver.go:263`（有守卫，也在取驱动器字母）＋本 helper:27。若这句写的是"internal/risk 的测试面里唯一一处"⇒ **真**；若写成"仓里唯一一处"⇒ **假**（`internal/tools/volume_windows.go:28/:54` 也在读 `VolumeName`，⛔ 切片） | 措辞债（⛔ 影响任何判语）；建议补射程二字＝"the single place **in this test surface** the drive letter is taken" |

## 另两处⛔ 算越claim，但⛔ 值得记一句

- 票 §现量 的行号锚（`anchor:77/78`、`junction:151/152`、`grep 命中＝2 枚`）今天⛔ 再是那两行（内容漂到 `:113`／`:151`＋`:158-162`；`grep '\[:1\]' internal/risk` 现命中 7 行，其中 4 行是注释）。票自己在 §裁定 §2ⓓ 里已把两处缺陷改成**内容锚**，同一把尺⛔ 回头改 §现量 那几行 ⇒ ⛔ 本票账（⛔ 改原句），记入终报供 A846 续。
- `junction:157-159`"On any volume-qualified temp -- every reachable surface on this box **and on the windows CI gear** -- haveShare is true and the three spellings below are built exactly as before." 里"on the windows CI gear"＝**⛔ 读数**（本腿连本机 `fsutil 8dot3name` 都被 `Access denied` 挡了，CI 那台的卷形制⛔ 知；件 `logs/v46-fsutil.txt`）。⛔ 越claim 到"防住了"那一档，但它**引用了一枚没量过的台面** ⇒ 归进 `Q-308-v1-A`（CI 色那笔）一起收。

## 建议替换句（对象＝上面 #1；⛔ 本腿动手，⛔ 本腿⛔ 表态该⛔ 该改）

> The 8.3 spelling is only folded in when this volume generates short names (same probe the ticket-18 case uses). Since ticket 308 the two UNC spellings are conditional too: `driveShare` reports no drive letter on a volume-less path and `volumeInvariantSpellings` drops those two BY NAME, printing kept/dropped. So the denominator this case measures is 9 (10 with 8.3) on every volume-qualified temp -- the only kind Windows hands out here -- and 7 (8) on a volume-less one, which is pinned on the INPUT face by `TestVolumeInvariantSpellingsVolumeLess308`, never by an env shape. "Never quietly stops measuring" holds because a drop is NAMED and LOGGED, not because the spellings are unconditional.

（理由：这句把**枚数**、**哪一面被钉住**、**为什么⛔ 是"无条件"守住的**三件事同时写清，⛔ 留"every other spelling is unconditional"那枚现在⛔ 假的支柱。刀⛔ 归本腿＝票面"⛔ 顺手改对"＋本腿＝⛔ 实现者但同样⛔ 有权改别人写的码。）
