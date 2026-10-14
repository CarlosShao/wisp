# 236s — 票 236 **AC#6 甲案**（`fs_broken.go` 只排 import 顺序、保留"未定义符号"坏点）

- 腿：`236s`（实现者；⛔ 框未翻、`Status:`／票名未动、零 push、零 `add -A`）。时刻：2026-10-14 19:1x–19:3x +08。
- 面：`.scratch/wisp/probes/185/c1/mut/fs_broken.go` 一枚 ＋ `probes/236/s/`（本件＋logs 6 枚）。
- 窗口：起手锚 `81730cea`（19:1x）→ 编辑窗 HEAD `da389ca5`（19:25:13）→ 交件前复读 HEAD 仍 **`da389ca5`**。
  漂移 8 笔（逐笔）：`72db0f39`（编排者：159 翻四格＋**派 236s**）／`898ce5f3`／`b6e959b1`／`b6fc818d`／`41071cdc`／`d01f4b71`／`da389ca5` 等；**零笔与本枚夹具重叠**（`git log --oneline 81730cea..HEAD` 原文在 `logs/05-gates.txt` 首段）。

## §0 结论先行

AC#6 **甲案已落地**：`fs_broken.go` 的 import 块挪到其它声明之前，marker 行原样保留在 import 块之后；`undefinedSymbol185c1` 仍未定义（overlay 复现 rc=1 逐字 `undefined: undefinedSymbol185c1`）。
改动＝**纯顺序移动，−1／+2 行**（`git diff --numstat` ＝ `2 1`），两 hunk 全在文件头 import 区，唯一落笔文件。
改后 `gofumpt -l` **空 rc=0**；`sh scripts/d22scan.sh` rc=0（正控 35／0／0 先绿）；长度门 **GREEN 60=60**。
**硬闸判读＝不咬，继续甲**（依据全在 §1：机械钉 0、逐字节可复核副本 claim 0、185 名下 md5／字节／行引用 0；仅有 H 档 digest 历史读数两处，与已获接受的错误文案同级——⚠ 若编排者判这两处 digest 读数＝"硬引用"，本腿整个 commit 可退，退法见 §5）。

## §1 引用普查（步骤 1 硬闸；全部只读；原文落 `logs/01-census.txt`）

1. **机械钉＝0 枚**：`git grep -n "fs_broken\|undefinedSymbol185c1\|probe185c1BrokenMarker" -- "*_test.go" "scripts/*" ".github/*"` ⇒ **rc=1、零命中**。全仓没有任何 test／脚本／workflow 会 open／断言这枚文件或其符号（逐字在 `logs/01-census.txt`）。
2. **逐字节可复核副本 claim＝0 枚**：本票口径（a2 表头末列「产品码快照？」）该枚＝**否**（`a2/census.md:21`）；全仓搜 `逐字节|byte-identical|IDENTICAL|cmp |原字节` 与 `fs_broken` 同现＝零枚针对本文件的字节保真要求（"逐字节"命中全是 overlay 机制描述／别的台件）。
3. **185 名下 md5／字节数／行数引用＝0 处**：`grep -rn "fs_broken" .scratch/wisp/probes/185 docs/evidence/s1/185-reread-owner-census-c1.md` ＝ **5 命中**（`overlay-proof.json` 机器消费的路径 1／`overlay-proof-broken.txt` 错误文案 1／185 证据件 `:48`（错误文案）`:185`（rc=1 设计）`:191`（证明件保留）3）；185 名下 `[0-9a-f]{32,64}` 扫过＝只有门条款日志里的**别的对象**（提交号 `e9ef94d0…`）。
4. **全仓 digest 级引用＝2 处，皆为 H 档历史读数（非钉）**，具名：
   - `.scratch/wisp/probes/232/v1/verdict.md:152`：识别那枚 tracked 坏样本时录了**全文 md5 `f5c325e04a9a0d94421a665d9479e7b5`**（与错误文案同句；232 已归档）。
   - `.scratch/wisp/probes/269/a1/census.md:101` 与 `:159`：录了 **全文 blob `ee58385366d8fb6693c12fe0a4367225538f9dd4`**（"本腿现读"，用于论证其甲形"样本原字节"⇒引用零风险）。
   字节／行数读数：`236/a1/census.md:185`（`11,041 字节／338 行`）、`236/a2/census.md:21`（表，`11,041｜338`）、`275/a1/census.md:330`＋`census-2.md:180`（`ls -la … 11041`）。
5. **判读**：以上 digest／字节／行数引用全在 census／verdict **散文件**里，是"记的是当时读数"的 **H 档历史快照**（269 §6 自己的三档定义：M＝机器真会去 open＝仅 overlay 路径一枚／R＝路径级规范引用／H＝历史读数）；无一处是"机器会照着比"的**钉**。与**已获编排者接受的错误文案级联**同类（票 236 `:104`／a1 `:191` 已裁"改后错误文本换句+台账更正即可"）⇒ 不满足"停手上报"条件，继续。
   **具名自认**：本腿是凭这条判读走下去的；若编排者判 digest 读数也算断闸，请改判（退法 §5 第 1 条）。

## §2 改动（唯一落笔文件；`git diff` 全文在 `logs/03-diff.txt`）

```diff
@@ -1,5 +1,4 @@
 package tools
-var probe185c1BrokenMarker = undefinedSymbol185c1
 
 import (
@@ -15,6 +14,8 @@ import (
 	"github.com/CarlosShao/wisp/internal/risk"
 )
 
+var probe185c1BrokenMarker = undefinedSymbol185c1
+
 // The D34 fs family lands slice by slice. This file holds the L0 pair:
```

- `--numstat` ＝ `2	1`；两 hunk 均在文件头 import 区；marker 行文与注释**一字未改**（只是位置：:2 → **:17**；文件内既有注释全部逐字节保留，含 `:169` 的 `185-c1 MUTANT …` 段——"注释仍说清它是 185 的突变夹具"＝**未动**）。
- 仓外演练（`C:\tmp\wisp236sgotmp\work\cand{1,2,3,4}.go`，只建不删）：cand1（`)` 后直接 var）／cand3（package 后无空行）／cand4（var 与后随注释间无空行）**都被 gofumpt 判要格式**；**唯一干净形＝cand2**（package 空行＋import 块＋空行＋var＋空行＋既有注释），即落盘形；`cmp` 与落盘文件 **IDENTICAL**。
- ⚠ **顶一条 a1 预测**：`a1/census.md:185` 预测甲"字节总数与行数不变"——**现测不成立**：gofumpt 干净形必然在 marker 与后随注释间多一个空行 ⇒ **338→339 行、11041→11042 字节**。语义仍是"只动 import 顺序"（那一空行是格式器的分隔要求，不是内容增删）。
- 指纹（改前→改后）：md5 `f5c325e0…`→`c43e4b1e…`；sha1 `f9bc1dbd…`→`f4fd2d47…`；sha256 `6db34456…`→`a14dbcc9…`；blob `ee583853…`→`20f5ae54…`；字节 `11041`→`11042`；行 `338`→`339`。rc=1 的位置随行号变（见 §3）。

## §3 判据（逐发原文；`logs/02-pre.txt`／`03-diff.txt`／`04-post.txt`）

- **改前 `gofumpt -l <file>`**（本机 `v0.12.0 (go1.27.1)`，`D:\work\base\gopath\bin\gofumpt.exe`）：
  `.scratch\wisp\probes\185\c1\mut\fs_broken.go:4:1: imports must appear before other declarations`；**rc=2**（`-d` 同句同 rc）。＝**非空** ✓
- **改后 `gofumpt -l <file>`＝空、rc=0**；`gofumpt -d` 空、rc=0。＝**空** ✓
- **"符号坏"仍成立·证法＝按票 185 原有用法重放一发（overlay 构建）；选择与结果具名**：
  - 改前（同仪器）：`go build -overlay=.scratch/wisp/probes/185/c1/overlay-proof.json ./internal/tools/` ⇒ rc=1、逐字 `.\.scratch\wisp\probes\185\c1\mut\fs_broken.go:4:1: syntax error: imports must appear before other declarations`。
  - **改后**：同命令 ⇒ **rc=1、逐字 `.\.scratch\wisp\probes\185\c1\mut\fs_broken.go:17:30: undefined: undefinedSymbol185c1`** ⇒ 换句不换坏点（rc=1 不变、理由＝未定义符号）。
  - 正控（同批）：`-overlay=…/overlay-muta.json`（干净突变副本 fs.go）改前改后 **各 rc=0** ⇒ 映射生效、"符号坏"不来自构建动作。
  - （未选"临时包名拷进仓外树"那形：该文件 import `internal/risk`，仓外模块过不了 internal 边界；overlay 复现＝票 185 原仪器，证据更强。）
- **`sh scripts/d22scan.sh`**＝rc=0；先绿正控 `top-level: PASS=35 FAIL=0 SKIP=0`；`clean - no D22 ban violations`；各 scope：bans #1-5 internal/=237・cmd/=40・#6 frontend/=115・#7 internal/tools/=24・**#8 internal/=556**・#8 cmd/=128（原文 `logs/05-gates.txt`）。
- **长度门** `sh scripts/check-path-length-budget.sh`＝rc=0，`denominator: tracked paths=12954 over-budget=60 covered by roster=60 not in roster=0` ⇒ **GREEN 60=60**（原文 `logs/06-gates.txt`）。
- **附加读数**：tracked 集 gofumpt 尺里 `fs_broken` 命中数 **0**（改前它在；其余 14+ 枚 A-ROSTER 夹具原样在册、一枚未动）；scoped porcelain：`internal cmd`＝**0 行**、`scripts .github`＝**0 行**；本腿面只有 ` M fs_broken.go` ＋ `?? probes/236/s/`。

## §4 具名没做（一律零触碰）

- 其余 14+ 枚 A-ROSTER 夹具（163/a1・174/c2・183/r2・185/r1・197/r1c×2・212/v1・220/r1・222/v1・224/v2・231/v1・235/v1・241/v1・253/v2・257/r2・259/r1・263/v1・33/p1…）**一枚未动**（那是乙案/别的账）。
- `.github/**`、`tools/d22scan/**`、`scripts/attrib.sh` 零触碰；**未改门的分母、未加 `if:`、未动 `continue-on-error`**。
- 185 名下其它件（三枚 overlay json、`logs/**`、`mut/fs.go`、`mut/taintmatch.go`、证据件）**一字未动**；未删件（`issues/README` 规则 8）。
- 台账 `docs/reports/pending-and-issues.md`、`docs/evidence/**`、停车点 **未写**（归编排者收件时裁）。
- 未跑 CI 复读、未 push、未碰他腿脏件（r5/r6/r7/v2 探针、`cmd/wisp`、`scripts/slo-*`、`probes/154` 系）。

## §5 给后程（收件与后算的影响面）

1. **建议的台账更正（一条，覆盖三件事，⛔ 本腿未写＝归编排者）**：自 `<本腿新锚>` 起 `.scratch/wisp/probes/185/c1/mut/fs_broken.go`：(a) overlay 生效证明错误行由 `:4:1 imports must appear before other declarations` 变为 **`:17:30 undefined: undefinedSymbol185c1`**，rc=1 不变；(b) md5 `f5c325e0…`→`c43e4b1e…`、blob `ee583853…`→`20f5ae54…`（232 `verdict.md:152` 与 269 `census.md:101/:159` 记的是**当时**读数，复算时按本条对表）；(c) 字节 `11041`→`11042`、行 `338`→`339`（干净形多出的分隔空行）。若编排者判 (b) 属"硬引用"（见 §1 第 5 条自认），退法＝对本 commit 追加反向 commit（恢复 `ee583853…` 那版字节），并将裁定改到乙／丙。
2. **一枚观察（未动，具名报文）**：改后文件 `:169` 的注释仍写 `185-c1 MUTANT (probe copy only, never a tracked file)`——"never a tracked file"与事实相反（它自 `4813567e` 起就是 tracked 件，正是 AC#6 的病灶本身）。它属于 185 突变体的原生注释、**不在"只动 import"射程**，未碰；是否改写归编排者／185 的账。
3. 改后该夹具只剩**符号级**坏点：任何"顺手补个 `var undefinedSymbol185c1 any`"式的修复都会把 185 的防恒绿自证弄绿 —— ⛔ 后程勿修。
4. overlay 复跑入口仍只有人工（脚本层 0 枚，269 §6 已记）；本腿的复现命令＝§3 那两行，可直接重放。

## §6 件与 logs

| 件 | 内容 |
|---|---|
| `logs/01-census.txt`（48 行） | 步骤 1 普查原文：机械钉尺、digest 检索、185 名下扫、锚与漂移 |
| `logs/02-pre.txt`（31 行） | 改前指纹／`gofumpt -l/-d` rc=2／overlay 前发 rc=1＋muta 正控 rc=0 |
| `logs/03-diff.txt`（33 行） | `git diff` 全文＋`--numstat 2 1`＋`cmp` cand2＝IDENTICAL＋改后指纹 |
| `logs/04-post.txt`（20 行） | 改后 `gofumpt -l/-d` 空 rc=0／overlay 复现 rc=1 `undefined`／muta 正控 rc=0 |
| `logs/05-gates.txt`（39 行） | 窗口与漂移 8 笔＋d22scan 尾段 rc=0 |
| `logs/06-gates.txt`（81 行） | 长度门 GREEN 60=60＋scoped porcelain＋tracked 尺里 `fs_broken`＝0 |
| 仓外临时件 | `C:\tmp\wisp236sgotmp\work\cand{1,2,3,4}.go`（只建不删） |

**提交**（只 commit、⛔ 零 push；显式 pathspec，逐笔 `git show --stat` 已在手）：夹具 `4961c3b5`（1 枚、+2/−1）／交付 `e915413a`（7 枚、340 行）／票 236 bullet `4abdbf7d`（1 行）／本行补记（第 4 笔）。三笔皆落在编辑窗 `da389ca5` 之上（起手锚 `81730cea` → 收件前 HEAD 逐笔现读为准）。
