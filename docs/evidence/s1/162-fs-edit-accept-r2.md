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

<!-- 后续格逐格追加 -->
