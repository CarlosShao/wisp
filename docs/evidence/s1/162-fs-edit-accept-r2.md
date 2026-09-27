# 票 162 · 162-v2 验收表（非实现者）＝裁 r3 四件：AC#3b 两形响亮 ／ 那半格反向尺 ／ AC#4b 真硬杀 ／ 判据顺序

- 验收程＝`162-v2`｜派单＝`.scratch/wisp/dispatches/2026-09-27-165x-accept-162-r3-v2-four-items.md`
- 被验对象＝`6823cbc6`（四枚先 `git cat-file -t` 复核＝`commit` 再用）／`b1cbdad4`／`3d98a8ab`／`87ca8d86`，件＝`docs/evidence/s1/162-fs-edit-r3.md`＋台件 `probes/162/r3/**`（被验件，本程一字未动）。
- **总判语：四格一律「成立」，零退回、零附条件**。三处"证据文字 vs 现量"的过期读数（行号差一、红枚数时序漂移、残件持续性措辞）入账不返工，逐条见各节与 §14。
- 本程票面一格未勾、`docs/reports/**` 一字未写、`internal/**` 零字节（变异全走 `-overlay`＋仓外副本 `probes/162/v2/mut/`）。

## 0. step-0（分支＋锚的 sha 逐发入表）

```
$ git rev-parse --abbrev-ref HEAD                     → dev          （起手现量）
$ git rev-parse HEAD                                  → 4193f3a4…    （本程锚＝被验四枚全部入库之后）
$ git cat-file -t 6823cbc6 b1cbdad4 3d98a8ab 87ca8d86 → commit×4
$ git diff --numstat 87ca8d86..HEAD -- internal/tools → （空＝HEAD 的 internal/tools 与被验收口逐字节同）
$ git status --porcelain -- internal/tools            → 0 行（安静窗口成立：本程全程无兄弟写手进场）
$ git diff --numstat ff000784..HEAD -- internal/tools → 132/6 fs_edit.go ＋ 54/56 ac34 ＋ 384/0 ac3b ＋ 274/0 ac4b
                                                       （与被验件 §2 逐字节对得上）
```

- 名册／判据性读数一律 `git grep <pat> <sha>` 钉死：基线名册锚 **`ff000784`**，现名册锚 **`4193f3a4`**，改前顺序文件锚 **`3d98a8ab^`（＝`b1cbdad4`）**。工作树只用于跑 `go test`（跑前验过与 HEAD 无差）。
- 别程 `171-v1` 在场：本程写面只有 `docs/evidence/s1/162-fs-edit-accept-r2.md` ＋ `.scratch/wisp/probes/162/v2/**`，名册筛路径全程用 `git -c core.quotePath=false` 语义的 sha 锚定命令，未取任何工作树名册读数。

## 1. 本程没测什么

1. **AC#5／AC#6 两格一律未判**（不在本派单四件之内）。
2. **非 Windows 那一支没跑**：本机即 Windows，Linux/CI 侧对 `fs_edit_ac4b_kill_windows_test.go` 是"根本不进编译"——被验件标"未经验证"这个措辞本程认（§5c 判语），但"认措辞"≠"测过那一支"，那一支依旧无人测过。
3. 没跑全仓 `go test ./...`（逐包规矩）；`-overlay` 与 `-cover*` 全程未合跑。
4. **残件的"持续在盘"行为只验了反方向**：两次 PASS 跑之后 `find` 整个 `%TEMP%`（只读）＝`.wisp-tmp-*` 0 枚、`TestFSEditRealTaskkill*` 目录 0 个——即 PASS 之后 Go 的 `t.TempDir()` 收尾把残件连目录一起收了。本程没有、也不许去"恢复"它（见 §11）。
5. 面板/UI 那一层没看：AC#3b 的"外部可见"本程只复算到**回执＋审计行**两处（票面 `:38` 的口径）；这两处再往外（panel 渲染）不在被验断言里，本程不替它加判。
6. `TestFSEditRealTaskkill` 的 30s 超时失败支与 `taskkill` 缺失硬失败支**没有被主动演练**（只在 M-3 overlay 下被动看见一次超时支真响——那反而证明它响得响亮，见 §4 红枚数 12 的归位）。
7. 四数名册 `comm` 之外的第三向（子测试名册）没做——本程口径＝顶层四数＋顶层名册，与派单一致。

## 2. 件① AC#3b 两形响亮 —— 判语：**成立**

### 2a 硬拒的落点（现量，锚 `4193f3a4`）

```
$ git grep -n 'leadingBOMLen(content)\|lineEndingConvention(content)' 4193f3a4 -- internal/tools/fs_edit.go
internal/tools/fs_edit.go:192:  if n := leadingBOMLen(content); n > 0 && start < n {     ← 形①（定位循环内）
internal/tools/fs_edit.go:245:  if conv := lineEndingConvention(content); conv != "" && … ← 形②（updated 之后）
internal/tools/fs_edit.go:271:  res := t.d.stageAndRename(…                               ← 第一发写入
```

⇒ 两枚检查都真在**任何写入之前**，都走 `refusal(...)`（`IsError=true`），整批 all-or-nothing 性质继承成立。
⚠ 被验件 §3 写的 `:193`／`:244` 与现量 `:192`／`:245` 各差一枚——行号是读数、天然带版本（本仓既有账），**实质（先于 :271）不受影响**，入账不返工。

### 2b "响亮"两处外部读数——本程自己跑出来的原文（非引它的数）

`go test -count=1 -v ./internal/tools/`（rc=0；原始件 `probes/162/v2/logs/gates-gotest.txt`）内：

```
fs_edit_ac3b_test.go:134: AC#3b① 改后读数（回执）：IsError=true ErrorClass="tool" AppliedSteps=0
fs_edit_ac3b_test.go:137: AC#3b① 改后读数（审计行，外部可见）：tools: call kind=error … tool=fs.edit risk=L2 decision=allow outcome=error rules_hit=[R1] …
fs_edit_ac3b_test.go:251: AC#3b② 改后读数（回执）：IsError=true ErrorClass="tool" AppliedSteps=0
fs_edit_ac3b_test.go:254: AC#3b② 改后读数（审计行，外部可见）：… outcome=error …
fs_edit_ac3b_test.go:272: AC#3b② 合法对照一（同一枚编辑、new 按 CRLF 拼写）：23 字节 → 31 字节，CRLF 3→4 组、裸 CR/裸 LF 均 0
```

审计行由 `auditLineFor` **抓进断言**且钉"恰好命中一枚"（`fs_edit_ac3b_test.go:77-98`），对照一另钉 `outcome=success` 一枚（`:271`）。

**那个桥算不算"为了自证而造的旁路"——不算。** 判据三条（逐条现量）：
1. `fsEditBridgeAudited` 相对既有 `fsEditBridge` 只多一个 `Options.Logf` 捕获——`bridge.go:70-73` 是生产字段、`:174` 接进 `b.logf`、注释点名 `wired to agentRuntime.auditf by cmd/wisp`，即抓的是**生产审计汇**吐出的原文，不是测试自造字符串；
2. 桥的构造走的是生产构造器 `New(Options{...})`＋真 `BuiltinFSEntries`＋真 C26 路径解析，与被验行为同一代码通路；
3. 断言对象是**审计原文本身**（`Contains "tool=fs.edit"＋"outcome="＋wantSubs`），且命中数不恰为一即 `t.Fatalf`——"探针没接到也算绿"这一形被钉死。

### 2c 合法对照是不是装饰——**不是（本程自造变异亲测）**

自造 t1（`probes/162/v2/mut/t1/fs_edit.go`，锚 `4193f3a4` 的原件上把形①开火条件放宽成"无 BOM 文件命中第 0 字节也拒"；副本反扫 `TAMPER-V2-t1` present=1／原条件 absent=0，非静默 no-op）：

```
$ go test -count=1 -overlay=.scratch/wisp/probes/162/v2/overlay-t1.json -v -run 'TestFSEdit(StillLands|RefusesAnOldThatSwallows)' ./internal/tools/   → rc=1
--- PASS: TestFSEditRefusesAnOldThatSwallowsTheFileBOM        （拒形仍响，不受干扰）
--- FAIL: TestFSEditStillLandsAnOldStartingAtByteZeroWithoutABOM
    "an offset-0 old in a BOM-less file must land: … IsError:true … 圈住了文件头那 0 字节 EF BB BF …"
--- PASS: TestFSEditStillLandsNewEndingsWhereThereIsNoConventionToBreak
```

⇒ 把合法那形也拒了，对照用例**红**且红因点名到位。五枚对照（BOM 之后抄起／无 BOM 第 0 字节／CRLF 拼写的 new／LF+LF／已混合不被告知）本遍全绿在册（`gates-gotest.txt`），抽测其一即证整类非装饰。

### 2d JSON 残缺 BOM 射程读数——独立复算＝**成立，且比证据件写得更强**

被验件只说"Go 编码器换成 U+FFFD"。本程 `probes/162/v2/jsonprobe/main.go` 现跑双向：

```
input valid-utf8=false bytes=bf 61 6c
MARSHAL out={"old":"<U+FFFD>al"} err=<nil>          ← 编码侧：残缺 BOM 字节进不了 JSON 字符串
UNmarshal rc-err=<nil>
decoded old bytes=ef bf bd 61 6c valid-utf8=true     ← 解码侧：手搓 raw JSON 塞 raw 0xBF，Go 解码器同样换成 U+FFFD
would decoded old match file bytes [2:5]=bf 61 6c ? false
```

⇒ 这一形**根本进不来**（编、解两向都回不来那三个字节），不是"修好了"；①唯一可达形＝"覆盖整枚 BOM"，而那形有响亮拒（本程 2b）。`TestFSEditCannotBeFedAPartialBOMThroughJSON` 两向读数与本程独立复算一致。

## 3. 件② 那半格反向尺 —— 判语：**成立**（最硬一问：那句话**成立**，v1 两枚绿**不是洞**）

**最硬的一问先答**：派单要求"要么找到一发非恒等的强制改行尾（＝假闭退回），要么独立复算支持它"。本程**独立复算支持**，过程如下。

1. **M-3 副本亲验**（`.scratch/wisp/probes/162/r3/mut/m3/fs_edit.go`，diff 对 `ff000784` 原件）：变异体＝
   `updated := strings.ReplaceAll(strings.ReplaceAll(b.String(), "\r\n", "\n"), "\n", "\r\n")`——先塌后涨。
   代数复算：**纯 CRLF**（无裸 `\r` 无裸 `\n`）塌成 LF 再涨回 CRLF＝逐字节恒等；**纯 CR**（既无 `\r\n` 也无 `\n`）两遍 ReplaceAll 都不可作用＝恒等；**含裸 LF 的混合**（裸 LF 被涨成 CRLF）＝非恒等；**纯 LF**＝非恒等。副本反扫 `present=1/absent=0`（`MUT-R3 same shape` 在、`updated := b.String()` 0 命中）＝变异非静默。**找不到任何"对纯 CRLF 或纯 CR 非恒等、又配得上'强制 CRLF'这个名字"的改法**——若把强制写成裸 `\n`→`\r\n` 不先塌（会造出 `\r\r\n`），那已经不是"盖约定"而是"损坏字节"，v1 正向那枚逐字节相等断言当场就能抓（本程 4 验证同一判断：真红的全在非恒等形上）。
2. **红绿名册本程亲跑**（`-overlay=r3/overlay-m3.json -run '^TestFSEdit'`，锚 HEAD；原始件 `probes/162/v2/logs/m3-on-r3.log`）：
   - `TestFSEditLeavesNeitherTheCRLFNorTheBareLFOfAMixedFileRewritten` **红**（那半格要的正是这一条）；
   - `TestFSEditOnACRLFFileAppliesACRLFSpelledOldAndStaysCRLF` **绿**、`TestFSEditLineEndingForensicsRefusesRealSpellings` **绿**——v1 点名的两枚复算仍绿；
   - 新尺在 base overlay（未修码）与修码后都绿（§4 恒真栏）＝它不是"喂什么都会响"。
3. **两枚绿的"为什么绿"各有结构性理由（读断言本体，不是听转述）**：CRLF 正向那枚夹具＝纯 CRLF、old/new 全 CRLF 拼写，M-3 对它的落盘字节是恒等（1 的复算）⇒ **该夹具形状上不可能长出对这一族变异的判别力**；五形取证那枚每一腿都在**写入之前** 0 命中拒掉（`fs_edit_ac34_test.go:328-329` 断"拒＋文件未动"）⇒ 落盘期的任何变异**物理上到不了它**。判别器唯一能长的形状＝"含 CRLF、又不是只有 CRLF"——正是新尺的 22 字节混合夹具。**结论：那半格不是假闭。**
4. **一处时序漂移入账（不返工）**：本程 HEAD 上同一命令红 **12** 枚 ≠ r3 存档 logs 的 **11** 枚。差的那一枚＝`TestFSEditRealTaskkill…`（30.45s），红因**合法且正当**：M-3 把它的 LF 夹具目标盖成 CRLF，落盘前撞上②号守卫、子进程硬拒不写盘，父进程"30s 没等到 write:8 边界"响亮失败——恰是 r3 §4 那句"这一味现在有人会挡"的读数。r3 那一遍量于件②用例入库**之前**，两本账各对自己那一刻为真；登记给下一程，别让读者以为谁数错了。
5. **残留一枚射程外盲点（只登记，不折进本格）**："改完还是 **CR**"方向今天没有任何落盘用例（纯 CR 只存在于取证那枚的**拒**腿里）——对"裸 CR→CRLF"这一族规范化，整包当前无红。它与 r3 自报的"混合行尾不被告知"同族、越出票面 `:37` 那两形 ⇒ 归 owner 新判据，不归本程加判。

## 4. 件③ AC#4b 真硬杀 —— 判语：**成立**

### 4a SKIP=0 复跑（本程口径＝编排者 16:52 同口径）

```
$ go test -count=1 -v ./internal/tools/    → rc=0
RUN=150  顶层 PASS=101  FAIL=0  SKIP=0  panic=0
--- PASS: TestA18RealTaskkillLeavesTheTargetWholeAndTheNextWriteReclaimsTheStagingFile (0.98s)
--- PASS: TestFSEditRealTaskkillLeavesTheTargetWholeAndNamesTheStagingResidue (0.68s)
```

与派单转述的"101/0/0"对上；本格与 A18 同遍 PASS＝"本机有分母"的凭据独立成立（非推论）。
`t.Skip` 复核：`grep -n "t\.Skip"` 该文件 1 命中、位置＝`:36` **注释句**（"There is no t.Skip in this file"）——真调用 0 枚；`taskkill` 缺失走 `t.Fatalf`（`:159`）＝硬失败不静退。**"非 Windows 那一支未经验证"这个措辞：认。** `//go:build windows` 使它在 Linux 上根本不进编译，既无分母也无读数，诚实标法只有这一种。

### 4b 三条判据逐条复算

- **(a) 字节**：本程跑内读数 `fs_edit_ac4b_kill_windows_test.go:220`＝目标 31 字节＝改前原文、74 字节 `"ALPHA-KILLED…"` 一字未现；`:214` 的 `HasPrefix(got,intended)`／`Contains(got,new)` 反向断言在位。尺非死开关：本程 M-4 overlay 亲跑该枚（`logs/m4-on-ac4b.log`）⇒ **红于 :210 `got 4 字节 41 4c 50 48`**，与 r3 读数同向；`(b)` 被 `(a)` 挡＝编排者 16:5x 已裁"故意不改是对的"，本程复核无新证据翻它。
- **(b) 残件点名**：本程遍内 `:253` 点名 `.wisp-tmp-5c06588b-36180-1572977897`（4 字节、内容 `"ALPH"`＝新内容前 4 字节，`len!=1` 即 `t.Fatalf`）——形状与 r3 那枚逐字节同款，**可重跑性由本程这一遍亲证**。⚠ **一处持续性更正**：本程跑后 `find %TEMP%`（只读）＝r3 具名残件与本程残件都已不在盘上——用例 PASS 后 Go 的 `t.TempDir()` 收尾连目录一起收（§5b 自注"目录级回收是 Go 收尾做的"与此一致），所以"残件留在盘上等回收"**只在执行窗内为真**，可复核凭据是那行点名日志而不是那个文件本身。判语不受影响（断言在窗内、先量残件后跑 fs.read 的顺序在码上）；措辞入账。本程全程**没清、没还原、也无法清**（不存在可清之物）。
- **清扫器那枚小账（编排者"顺手判一下"，判语如下）**：**不必在本格钉用例，机制已有尺。** 现量：清扫器＝`reclaimStaging`（定义 `fs_staging.go:192`），调用点在**共用的** `stageAndRename` 内（`fs_write.go:288`，先扫后 `CreateTemp`），`fs.edit` 恰走这一枚函数（`fs_edit.go:271`）；"下一次写盘带走残件"这一行为已被 A18 的第③子例钉住（`bridge_a18_kill_windows_test.go:216 the_next_write_reclaims_the_staging_residue`）＋ `fs_staging_windows_test.go` 两腿。缺的只是"`fs.edit` 杀后 → 下一次写盘"这一特定链条的点名回归——属共享层之上的归因糖，**建议归 AC#6 那程一并处理**（r3 §12 第 5 条的最小形状可用），不因它退回、不因它加派。

## 5. 件④ 判据顺序那一笔 —— 判语：**成立**

```
$ git show 3d98a8ab^:internal/tools/fs_edit_ac34_test.go | sed -n '448p'
    t.Fatalf("expected the write to reach the staging boundaries, got %q", killed)   ← 改前：边界探针真在 :448
$ go test -count=1 -overlay=.scratch/wisp/probes/162/r3/overlay-m4.json -v -run '^TestFSEditKilledMidWriteLeavesTheTargetByteIdentical$'   → rc=1
    fs_edit_ac34_test.go:456: AC#4 violated: after a mid-write kill the target holds 8 bytes "ALPHA-LO", want the original …
```

⇒ 字节断言（`:455-458`）现量**先于**三关边界探针（`:460/:465-471`）执行：M-4 下红落在字节上（本程亲跑，
`probes/162/v2/logs/order-after-m4.log`），不再是探针吞读数。探针一枚未删一枚未弱；`Skip` 0 枚；
该枚 diff 只落测试文件（`ff000784..HEAD` numstat 见 §0）。r3 没照抄编排者的 `:458`、自己现量成 `:448/:456`——
本程对这两个号也各自现量复核（`git show 3d98a8ab^` 那一发钉的是**改前树**，不是工作树）。

## 6. 恒真栏·逐件两向读数（本程亲跑，非引它的数）

| 件 | 未修码向（base＝`ff000784` 的 `fs_edit.go`） | 修码向（HEAD） | 答复 |
|---|---|---|---|
| ① 断言拒绝四腿（2 枚 Refuses＋SilentShape 两子例） | **全红**（顶层红 3 枚＝`base-on-unfixed.log`，枚枚点名；PASS=17/FAIL=3，其 20 枚＝r3 存 19 枚＋件③用例后入） | 全绿 | 非恒真，红因＝回执给成功 |
| ① 五枚合法对照 | 全绿 | 全绿 | 两向绿＝钉"没被放宽"，不参与恒真（它红的那一形见 t1 变异，§2c） |
| ① `…CannotBeFedAPartialBOM…` | 绿 | 绿 | 钉传输射程，两向绿**有意义**（两版都吞不进残缺 BOM，§2d 独立复算） |
| ② 混合反向尺 | 绿 | 绿 | 红只在 M-3 下（§3.2）＝反向尺本体，非前提尺 |
| ③ AC#4b 用例 | （base 含 stageAndRename ⇒ 编译期内不变异无从测"未修"；其"未修"凭据＝票面 AC#4 历史：Hooks.Kill 一支从未验过真杀） | 绿；M-4 下红于字节 | 两向＝PASS 遍＋M-4 遍（§4b） |
| ④ 顺序那一笔 | 改前树 M-4 红在 `:448` 探针（r3 存档＋本程对改前树行号复量） | 改后树 M-4 红在 `:456` 字节（本程亲跑） | 顺序差被同一变异两向摊开 |

## 7. 单点回退（逐枚："撤掉这处，哪条用例变得不响？"）

- 撤 `fs_edit.go:192`（BOM 拒）→ `TestFSEditRefusesAnOldThatSwallowsTheFileBOM` 与 `…/an_old_that_carries_the_bom_is_refused_r3` **变不响**（回执翻回成功、审计行断言失败连带红）；撤 `:245`（行尾拒）→ `TestFSEditRefusesALFSpelledNewInAPureCRLFFile` 与 `…/an_lf_spelled_new…_r3` 不响。**答得出。**
- 撤件②新尺用例 → 全套件在 M-3 下 CRLF 方向**失去唯一红**（§3.2 名册），"盖约定"从此不可见。**答得出。**
- 撤件③整文件 → M-4 下再无人报"目标停在新内容前缀"（本程那发红于 `:210` 消失）；撤其 `(b)` 段 → 残件计数不响。**答得出。**
- 撤件④换回旧序 → 红从 `:456` 退回 `:448`，字节读数再次被探针吞掉（改前存档 `order-before-m4.log` 就是这一形的历史样本）。**答得出。**
- 无一处需要"附条件"：没有哪枚是被别枚顺带照亮的装饰。

## 8. "最小可见单位"逐件答复

- **件①**：单位＝**一枚工具调用**（回执与审计行都是调用粒度——派单里"审计行是文件粒度"这一转述不成立，`bridge.book` 每调用一行、带 `corr=`）。两半在**同一次跑、同一枚调用**上各自断言（回执 `:134`、审计 `:137`），不存在"分开各自不响"：模型侧看得见回执、运维侧看得见审计行，两个受众各拿整份证据。
- **件②**：单位＝**一次落盘后的整文件字节**（逐字节相等＋行尾计数双钉）。
- **件③**：单位＝**目标文件字节**＋**目录内残件计数**——同一枚调用的两个不同对象，不是把一份证据拆两半。
- **件④**：单位＝**单条用例内的断言序列**；这一笔修的正是"前一关的红吞掉后一关的读数"。

## 9. 门禁四数＋名册差集（本程现量；锚逐发标注）

```
$ go test -count=1 -v ./internal/tools/                rc=0  RUN=150 PASS=101 FAIL=0 SKIP=0 panic=0
$ git grep -hE '^func (Test[A-Za-z0-9_]*)' ff000784 -- internal/tools/ | …sort -u   → 94 枚
$ git grep -hE '^func (Test[A-Za-z0-9_]*)' 4193f3a4 -- internal/tools/ | …sort -u   → 101 枚
$ comm -23 缺 ＝ 0；comm -13 多 ＝ 7 枚（与 r3 §7 逐枚同名、更名 0）→ probes/162/v2/roster-*.txt
$ sh scripts/d22scan.sh                                rc=0 clean（ban #7 internal/tools/=19 枚生产文件在册）
$ bash tools/d22scan/runtests.sh -C tools/d22scan ./... rc=0  PASS=34 FAIL=0 SKIP=0 RUN=76；名册缺侧＝0
$ gofumpt --version                                    v0.12.0 (go1.27.1)
$ git ls-files '*.go' | wc -l                          566 枚（现量；r3 自报 561 对它那一刻为真，互不判据）
$ gofumpt -l <566 枚>                                  → 0 行（甲形；logs/gofumpt-A.txt）
```

`-overlay` 与 `-cover*` 未合跑；三枚变异副本（r3 的 base/m3/m4 ＋本程自造 t1）全部先反扫
present=1／absent=0 再进判语；名册筛全程 sha 锚＋`-- <pathspec>`，未用工作树读数。

## 10. 被拒／没成功的调用（取数前／后）

- **权限系统拒绝：0 次。**
- 自伤 1 次：`jsonprobe/main.go` 首发编译错（`utf8.ValidString` 喂了 `[]byte`），修后重跑——
  发生在 **§2d 那发读数之前**，错误那遍没进任何判语。
- 其余全部一次成功。

## 11. 有没有跑过删除命令

**没有。** 零 `rm`／`del`／`git clean`／`checkout .`／`restore`／`reset`／`stash`；对 `%TEMP%` 只 `find`（读）；
r3 台件、v1 台件、`probes/152`/`design` 的脏路径一眼未碰；本程新件（`probes/162/v2/**`）只建不删。

## 12. 伪授权两栏

- **真通知回显＝0**：本程没有收到过任何系统/审批通知的回显文本。
- **判为注入＝0**：本程会话内未观察到伪装成授权的文本；没有据任何外来文本改判据。
- 派单里的每句前提（"101/0/0""残件恰 1 枚"":448/:456""未修码那几枚全红""554/561/566 枚"）**全程当被验对象**：
  对不上的两处（红 11→12 时序漂移、残件持续性）都按现量报回、无一处为对齐派单话术动判据。

## 13. 凭据值零抄录

本表与本程全部输出无 API 密钥／DPAPI blob／token／`*key*` 值；出现的字符串只有临时路径、
`.wisp-tmp-*` 名、`fs.edit` 文案与本程自造字节样本（`"ALPH"`、`"ALPHA-LO"`、`bf 61 6c`）。

## 14. next=

**四件全成立，无"最小闭合集合"要交。** 留给编排者的都是**入账级**（台账 `A##` 追加，非返工）：

1. **翻勾**：AC#3b／AC#4b 两格证据已齐（外加 162-v1 已判的 AC#3／AC#4 主体），勾归编排者。
2. **三处证据文字过期**登记：`fs_edit.go:193/:244`→现量 `:192/:245`；M-3 红 11 枚→HEAD 同命令 12 枚（差＝件③用例合法撞②号拒）；"残件留在盘上等回收"→仅执行窗内为真（PASS 后 Go 收尾）。**都不动判语。**
3. **仍开着的射程边界**（归 owner 新判据，别在 r4 偷加）：混合行尾无提示（16:5x 已裁"不加"）；"改完还是 CR"方向整包无落盘尺（§3.5）。
4. **AC#5／AC#6** 照票面：AC#5 若需动 `internal/risk/**`＝停手上报；清扫器那枚特定链条（fs.edit 杀后 → 下一次写盘带走）建议并入 AC#6 程，最小形状＝r3 §12 第 5 条。
5. 本程台件：`probes/162/v2/{logs/*.txt,mut/t1/,overlay-t1.json,jsonprobe/,roster-*.txt,tracked-go.txt}`，全部只建不删。

> 本表与派单 §5 固定顺序的对应：step-0＝§0；没测＝§1；四件判语＝§2–§5（"那半格最硬一问"在 §3、
> "SKIP=0＋清扫器"在 §4）；恒真栏＝§6；单点回退＝§7；最小可见单位＝§8；门禁＝§9；其余 §10–§14 原位。


