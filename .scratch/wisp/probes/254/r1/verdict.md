# 票 254 — 交件（腿 `254-r1`）

> 本腿只做读数与落地，**不做验收判语**；票面三格 AC 的勾选框一枚未碰。

## 0. 起手锚

- 起手时刻：`2026-10-02T10:22:30+08:00`（`date -Iseconds` 逐字）
- 起手 `git log -1`：`5801b91f probes(180-a1): §6 我可能判错的八条 + §7 判不动的四处（甲四乙三）先写满`
- 起手 HEAD 全号：`5801b91f06ecc8e3caef81bccbd361a94742168c`
- 分支：`dev`
- 起手 `git status --porcelain scripts/`：**空**（逐字：零行，`| wc -l` = 0）⇒ `scripts/` 面上没有别的写腿在飞，本腿未并发写。

### 0.1 票面行号复认（逐条亲验，非抄票面）

| 票面断言 | 复认结果 | 我用来验的尺 |
|---|---|---|
| `scripts/winsec-tests.sh:94` 把包路径**显式**传给 `portable-tests.sh` | **成立**（逐字见 §1.0） | `Read` 该文件全量 |
| CI 侧由 `.github/workflows/ci.yml:439` 调 `winsec-tests.sh` | **成立**，且该行不带任何参数 ⇒ 子进程实收 1 枚参数 | `grep -n 'winsec-tests\|portable-tests' .github/workflows/ci.yml` |
| 「今天没有任何一步以 `--scope=winsec` 的名义去跑那档」 | **名字层成立；推论层不成立** ⇒ 见 §5 推翻的第一句 | 同上 + `portable-tests.sh:328-349` 的显式路径识别支 |
| `scripts/portable-tests.sh:211` 是 `core` 档里那行 `./internal/winsec/` | **成立**（该行逐字 `        ./internal/plugin/ ./cmd/llmrecord/ ./internal/winsec/`） | `Read` `portable-tests.sh:195-235` |
| 现量 #1/#2（载具 18 绿 / unknown 消息 5 名 rc=2） | 见 §4（本腿复跑过） | `bash scripts/portable-tests-selftest.sh` |

## 1. AC#1：两支的代价读数与选定

（本节在骨架之后填写。）

## 2. AC#1 的「改前必红」那一发

（本节在骨架之后填写。）

## 3. 逐格读数

（本节在骨架之后填写。）

## 4. 门禁与卫生读数（hunk 枚数与位置）

（本节在骨架之后填写。）

## 5. 我推翻编排者哪一句

（本节在骨架之后填写。）

## 6. 判不动／没测到的地方

（本节在骨架之后填写。）

## 7. 待裁（AC#3 那两形的代价）

（本节在骨架之后填写。）

---

# 10-02 接管腿 `254-r1b`（接 `254-r1`，前腿死于模型服务连接中断）

> 本节是**接管腿**写的，⛔ 上面 §0–§7 一行未删、一字未改（那七节的空句仍在原处，属前腿原话；
> 本腿只把同一批读数在自己的节里填满）。票面三格 AC 的勾选框一枚未碰。
> 本节编号 8.x，与上面 §0–§7 不同名，避免"看起来像前腿填的"。

## 8.0 起手锚（本腿自取，非抄前腿）

- 起手时刻：`2026-10-02T10:37:11+08:00`（`date -Iseconds` 逐字）
- 起手 `git log -3`：`580d6153`／`d253703a`／`ed459d09` ⇒ 前腿原话与代提增量都在库里 ✓（派单要求的锚）
- 分支 `dev`；起手 `git status --porcelain scripts/` = **0 行**（尺：`| wc -l` 形，逐字零行）⇒ `scripts/` 面无并发写腿
- 本腿全程 **零枚真 `go test`**（`248-v1b` 正整包跑 `cmd/wisp`）。用到的真 go 只有两枚不编译调用：
  `go list ./internal/winsec/...` 与 `go list ./internal/winsec/`（见 8.4 的 core_pin 那格）
- 编排者已量的两把尺（`bash -n` 各 rc=0；票 250/251 载具打在该增量上 18/0）**本腿未重跑**；
  本腿另量的第一把是 `bash -n` 打在**自己新写的字节**上（载具 rc=0，那是新对象不是复跑）

## 8.1 对前腿那两枚形状的判断（先判断，再决定补什么）

| 前腿那枚形状 | 本腿判 | 凭据（全是本腿自己跑的） |
|---|---|---|
| `scripts/portable-tests.sh` 8/4：core 档那行目录形 → 与第五档同一个 `"$winsec_scope"` | **形状对，但它落盘时不被任何判据看见** | 载具 core 案 1–10 跑 `--scope=core` 时**不设** `FAKEGO_SCOPE_FILTER` ⇒ 假 go 照整份宇宙打印，目录形与 glob 形在那十枚里逐字同形；票 251 的 11–18 全跑 winsec/census，不看 core 行。⇒ 本腿 case 21 打在 `d253703a^` 字节上：`FAIL exit 0, expected 1`（core 在拆包种子上照绿）＋文本尺 `1 private spelling(s) / 0 shared` |
| `scripts/winsec-tests.sh` 37/0：新增 GUARD 3（只判"本已绿"那一发） | **形状对，两条条件都不冗余；它把 AC#1 ⓑ 从"我在显式路径上"变成可判的一格** | 三条：① 同一颗宇宙、只换调用参数 ⇒ widened 两支在 HEAD `rc=2`＋红句、在 `d253703a^` `rc=0`（3 条断言差＝这道守卫本身就是差额）；② 突变 A（子脚本仍宣告 `IS the winsec tier` 却 `pinned=''`）⇒ 本腿 case 19 红 4 条 ⇒ **case 19 判的是牙不是话**；③ 突变 C（子脚本仍钉但沿用调用方的目录形 scope）⇒ case 20 反控形 `rc=2` ⇒ **GUARD 3 的第二条（glob 形）不是装饰** |

⇒ 结论：**两枚都不必推翻、不必重写。本腿 `scripts/portable-tests.sh` 与 `scripts/winsec-tests.sh` 各 0 hunk**，
补的全在载具（判据层）。前腿 §1–§7 一个数都没填 ⇒ 它 §0.1 那句"名字层成立；推论层不成立"本腿按
〔仅自述〕处理，并在 8.2 把它变成有读数的格。

## 8.2 AC#1 收尾：ⓑ 那一格现在有了判据

新增两枚具名场景（`scripts/portable-tests-selftest.sh`），**驱动真链** `winsec-tests.sh -> portable-tests.sh`，
只假 `go`：

- `chain-explicit-split-audit-bites`（19）
  - ⚠ 两支都用**零参数**复现 `ci.yml:439` 真跑的那一发（ci 不带参数 ⇒ 父脚本自拼 `scope=("$target")`）。
    前腿的台件（`chain-run.sh`）与本腿早期的写法是**手传路径**，那在父脚本里走的是 `$# -gt 0` 的 else 支；
    本腿把它改到 ci 真走的那一支，两支不是同一条代码路。
  - 正控（同一颗种子拆包 ⇒ **必响**）：宇宙 = `winsec-pin-split`（`internal/winsec` ＋ `internal/winsec/acl`）
    ⇒ 链 `rc=1`，子脚本 GUARD C 点名 `> github.com/CarlosShao/wisp/internal/winsec/acl`、`Pinned: 1, resolved: 2`。
  - 反控（正常 ⇒ **不响**）：同链同假 go、宇宙 = `winsec-pin` ⇒ `rc=0`、`IS the winsec tier` 1 次、父脚本零条 GUARD 行。
  - 附带钉住："子脚本已经红了"那一发上 GUARD 3 **不叠第二句**（`hasnt 'winsec-tests\.sh: GUARD 3'`）。
- `chain-widened-scope-loses-audit`（20）
  - 一颗宇宙（`core_pin`，它认得 winsec 也认得 agent）、三种形，只换调用参数：
    widened 两路 ⇒ `rc=2`＋红句 `the delegated audit therefore did not happen`＋`IS the winsec tier` **0 次**
    ＋`runtests.sh: OK` **在**（子脚本是绿的＝缺口不报错的形状被钉死）；零参数 ⇒ `rc=0`；
    `--scope=winsec` 打给父脚本 ⇒ `rc=2` `GUARD 1 - the package list for this step does not name`。
  - 第三形的用处＝**ⓐ 不可从调用点偷换**（见 8.6，它同时是票面"ⓐ 最小"那句的反证）。

**改前必红那一发（具名）**：打在 **`d253703a^`** 的两枚脚本字节上（`d253703a^` = `ed459d09`），
整份载具 ⇒ **`27 case(s) ran, 9 assertion(s) failed`、`rc=1`**，红的 9 条**全部**落在本腿新写的场景里，
票 251 的 18 次调用在该基上仍全绿 ⇒ 基选对了（只动本票那一枚变量）。逐字红句（行号是本腿日志的）：

```
137:== case chain-widened-scope-loses-audit: rc=0
138:   FAIL exit 0, expected 2
139:   FAIL /winsec-tests\.sh: GUARD 3/ appeared 0 time(s), wanted at least 1
140:   FAIL /the delegated audit therefore did not happen/ appeared 0 time(s), wanted at least 1
153:== case core-tier-split-goes-red: rc=0
154:   FAIL exit 0, expected 1
155:   FAIL /GUARD C - scope mode=core resolved to a DIFFERENT package/ appeared 0 time(s), wanted at least 1
156:   FAIL /> github\.com/CarlosShao/wisp/internal/winsec/acl/ appeared 0 time(s), wanted at least 1
162:   FAIL the two tiers disagree about what bit - core says [], the winsec tier says [github.com/CarlosShao/wisp/internal/winsec/acl].
164:   FAIL core branch carries 1 private spelling(s) of ./internal/winsec/ and 0 shared
174:   FAIL the mutant seed took 0 line(s), expected exactly 1 - the core tier's winsec row moved
```

日志：`.scratch/wisp/probes/254/r1b/logs/prefix-full-at-d253703a-parent.txt`（逐枚基线另有
`prefix-chain-*.txt`、`prefix-core-*.txt`）。

⚠ **为什么不是 `f5f9cc34`**：那枚是票 251 落地**之前**的字节（`git merge-base --is-ancestor f5f9cc34 d253703a^` = yes；
它的 `portable-tests.sh` 572 行 vs `d253703a^` 668 行）。同一发打上去 ⇒ **17 枚场景、39 条断言红**
（`prefix-full-at-f5f9cc34.txt`），红的里面既有票 251 自己的 11–18 也有本腿的 19（那基上连
`IS the winsec tier` 都不会打印）。两枚变量一起动 ⇒ **不可归因到本票**，所以它只能当"更老的基"记着，
不能当本票的改前必红。⛔ 本票的差额只在 `d253703a` 那一枚提交里，基就得在那儿取。

## 8.3 AC#2 收尾：两档一致性的尺 + `core_pin` 要不要同笔动

新增两枚具名场景：

- `core-and-winsec-see-the-same-split`（21）
  - 一颗宇宙 `core-pin-split`（`core_pin` 25 行 ＋ `internal/winsec/acl` = 26 行），两档各跑一次：
    `--scope=core` ⇒ `rc=1` 且 GUARD C 点名 acl；`--scope=winsec` ⇒ `rc=1` 且点名同一枚。
  - **一致性尺**（这把才是"两档必须给同一个答复"）：把两份红日志里 GUARD C 的 `>` 行各自
    `sed -n 's/^portable-tests\.sh:   > //p' | sort -u` 后**逐字相等**且非空；在 `d253703a^` 上这把尺
    给出 `core says []`（core 根本没红）vs `winsec tier says [.../acl]` ⇒ 一绿一红，当场判失败。
  - **文本尺**：`core)` 分支体里 `./internal/winsec/` 私拼 **0** 次、`"$winsec_scope"` 共享引用 **恰 1** 次
    ⇒ 版式漂走（两枚都不满足）时当场红，不静默；`shared` 那把同时兼作"锚还在不在"的自guard。
  - 控制形：未拆包宇宙 ⇒ `--scope=core` `rc=0`、GUARD B 分母 25、零条 GUARD 行。
- `core-dir-form-hides-the-split`（22，反证，与票 251 的 case 18 同形）
  - 突变体只把 core 那一行换回目录形（种子尺＝`diff` 恰 1 行；取不上就报"种子取不上"并判失败）。
  - 同一颗宇宙上：突变体 `--scope=core` ⇒ **`rc=0`、分母仍 25**（漏计新包那一枚）；
    而**同一份突变体** `--scope=winsec` ⇒ `rc=1` 点名 acl ⇒ 一档绿一档红＝票面 ⛔ 禁的形状，现在被钉在载具里。

**`core_pin` 要不要同笔动？今天不要。**凭据是真 `go list`（不编译）：

- `go list ./internal/winsec/...` = 1 行 = `github.com/CarlosShao/wisp/internal/winsec`；
  `go list ./internal/winsec/` 亦 1 行 ⇒ 两形**今天**解析同集；
- 该行与 `winsec_pin` 抽出来的集**逐字相等**（尺：`YES`），与 `core_pin` 里 winsec 那 1 行**逐字相等**（尺：`SAME`）。

⇒ "同笔"不是欠着一笔没做，而是**被 GUARD C 钉成"发生拆包就必须同笔"**：拆包那天 core 与 winsec
两档同时红、且红在同一枚包名上（8.3 第一枚场景），谁只改一档谁吃红。票面那句"要改 glob 就同笔动
`core_pin`"在今天的宇宙下不成立为动作、成立为约束。⛔ 本腿没往 `core` 档塞包名（禁区），也⛔ 没动
`core_pin` 的任何一行。

## 8.4 门禁读数（逐格终值）

- `bash -n`：`scripts/portable-tests-selftest.sh`（新字节）⇒ `rc=0`。
- **载具终值（HEAD）**：`27 case(s) ran, 0 assertion(s) failed`、`rc=0`
  （`carrier-all-at-HEAD.txt`；18 → **27** 次脚本调用＝新增 4 枚具名场景／9 次调用。
  ⚠ "case" 这把尺一直数的是**调用次数**，不是场景枚数——本腿把 21 的两枚文本尺与集合尺的
  `ran` 增量去掉了，免得字节检查冒充一次执行）。
- 编排者那句提醒成立：**载具仍全绿不等于前腿那两枚被验过**——18 枚里没有一枚看得见 GUARD 3，
  也没有一枚钉"两档对同一枚包解析必须一致"。现在各有场景（19/20 与 21/22）。
- 每枚新场景单跑（HEAD）：19 ⇒ `2 case(s) ran, 0 failed`；20 ⇒ `3/0`；21 ⇒ `3/0`；22 ⇒ `1/0`。
- 突变台件读数（输出先落文件，⛔ 未改任何跟踪文件；副本在 `/tmp`，判据写进本文件）：
  - A `pinned=''`（子脚本仍宣告仍走 glob 但不钉）⇒ case 19 **红 4 条**、链 `rc=0`＋`runtests.sh: OK` 在
    ⇒ "只有宣告没有对账"这一形被 case 19 抓住，不是被文案抓住。
  - C `scope=("$@")`（子脚本仍钉仍宣告但沿用调用方目录形）⇒ case 20 **红 2 条**、反控形 `rc=2` GUARD 3 响
    ⇒ GUARD 3 第二条是承重的。
  - D 父脚本 `target` 丢尾斜杠（＝票 252 那一族拼法漂）⇒ case 19 **红 9 条**、链 `rc=2`、`IS the winsec tier` 0 次
    ⇒ 拼法漂**不可能静默变绿**；⚠ 残余：红句里 `scope=[${target}...]` 是拼出来的，`target` 一漂那句括号跟着漂成
    `scope=[./internal/winsec...]` ⇒ **读数看退码与短语，不看那句括号**（产码本腿未动，要不要改成引用子脚本原文归编排者）。
- ⛔ 票面 3 枚 AC 勾选框一枚未碰（`scripts docs tools cmd internal .scratch` 内 AC 框本腿零改动）。

## 8.5 卫生（hunk 枚数与位置，逐枚进表）

| 尺 | 读数 |
|---|---|
| `git diff --numstat d253703a..HEAD -- scripts/`（截至提交 #1） | `231 5 scripts/portable-tests-selftest.sh`；**`portable-tests.sh`、`winsec-tests.sh` 不在表里**（本腿各 0/0） |
| `git diff d253703a -- scripts/portable-tests-selftest.sh \| grep -cE '^@@'` | **6** 枚（默认 `-U3`）；同尺取 `-U0` 为 **7** 枚（无上下文时 `make_shadow_root` 与 chain 助手两处被拆开）⇒ 两把都记，差一枚就是上下文差 |
| 6 枚的位置 | `@@ -8,8 +8,10`（抬头"eighteen cases"改口）／`@@ -30,6 +32,24`（254 抬头段＋改前必红命令形）／`@@ -57,6 +77,10`（Usage 里的 `CARRIER_WINSEC_TESTS`）／`@@ -188,14 +212,24`（`make_shadow_root` 收第二枚参数＝父脚本）／`@@ -237,6 +271,51`（`run_chain` ＋ 254 的拆包宇宙）／`@@ -514,6 +593,153`（四枚新场景） |
| `git diff f5f9cc34..HEAD -- scripts/portable-tests.sh \| grep -cE '^@@'` | **8** 枚，⛔ **零枚属本腿**。逐枚归因（同一把尺按提交取）：`709d2589`=7、`4cca2c0e`=2、`d253703a`=2 ⇒ 逐枚相加 11 > 区间尺 8（跨提交的相邻改动在区间尺里并成一枚），两把尺都摆出来，不许拿其中一个说"只有 N 枚" |
| `wc -l -c`（实测） | `scripts/portable-tests-selftest.sh` = 755 行 / 38467 字节（改前 529/24878 ⇒ 本腿 +226/+13589）；本文件（verdict.md）＝49 + 本节；`portable-tests.sh` 672/38940 与 `winsec-tests.sh` 174/9142 **未变** |
| 证据件体量 | `.scratch/wisp/probes/254/r1b/logs/` = 17 个文件（含两枚基的字节副本 ×2 对、三次突变读数、四枚改前单跑、两枚改前整跑、HEAD 整跑） |
| 残留与空节自查 | 一把词族字符类尺（⛔ 词面只在命令行里跑、不落被扫文件；本行第一版抄了它的词面、当场被它抓到，见下方"自纠"），对 `scripts/portable-tests-selftest.sh` 命中 **0**；一把空节尺：全文件命中 **7**（⛔ 全是前腿 §1–§7 那七行原句，本腿一行未删），限定在 `^# 10-02 接管腿` 之后 = **0**。**⚠ 这条尺复现了本仓 10-02 那一枚复发：尺的字面文本一旦写进被扫文件，计数就不再是 0——本腿写完 8.5 那版自查句自己红了一次，改掉后重跑才是 0**（见 §8.9 自纠第 2 条） |

## 8.6 我推翻／更正编排者哪一句

1. **推翻（票面 AC#1 ⓐ 那句"写面只在 `scripts/`，最小"）**：ⓐ **不最小**。现量＝把 `--scope=winsec`
   直接打给 `winsec-tests.sh` ⇒ `rc=2`，红在 `GUARD 1 - the package list for this step does not name`
   （case 20 第三形）。ⓐ 至少要**同笔**改同一枚文件里的 GUARD 1（`:56-66` 认的是包路径）与 GUARD 3
   （`:156-168` 认的是子脚本 `IS the winsec tier` 那句宣告）；GUARD 2 依赖的 `pkgpath` 也得跟着换来源。
   ⇒ ⓐ/ⓑ 之差不是"写面大小"，是"要不要连守卫的名册一起换"。本腿仍按前腿的 ⓑ 交判据，**不改判**。
2. **更正（票面现量 #1"载具 18 枚"）**：现值 `27 case(s) ran`（本腿加的）。属过期读数非判错——
   那一行写的是 10:2x 的量，票面自己标了"待验断言"。
3. **证实（票面现量 #3 那条〔仅自述〕腿报）**：`winsec-tests.sh:94` 显式传路径 ✓（本文件 Read 到该行）；
   `ci.yml:439` = `run: bash scripts/winsec-tests.sh` **不带参数** ✓（尺：`grep -n 'winsec-tests' .github/workflows/ci.yml`）
   ⇒ "今天没有任何一步以 `--scope=winsec` 的名义去跑那档" **成立**。
4. **更正一处行号**：票面 AC#2 说 core 那行在 `scripts/portable-tests.sh:211`〔行号待验〕⇒ **现值是 215**，
   且那一行已从目录形换成 `"$winsec_scope"`（`d253703a` 干的），票面引的"它是目录形"这句对 `d253703a^` 成立、
   对 HEAD 不成立。⛔ 后续腿别按 211 找。
5. **补一句前腿没说的**：它的 GUARD 3 注释引 `c0aab35b`（＝前腿自己的骨架提交，10:26）与
   `before-widened.txt`（在库里 ✓）——本腿把这两个引用当真验过了：`git log -1 c0aab35b` 存在、
   日志存在、且该基上 widened 那发确实 `rc=0`。注释里的测试名/票号这次不是空头引用。

## 8.7 判不动／没测到的地方

- **载具不在 CI 里跑**（尺：`grep -n 'selftest' .github/workflows/ci.yml` ⇒ **0 命中**）。
  ⇒ 本腿新加的 4 枚场景是**笔记本尺**；CI 上有牙的仍然只是 GUARD 3 本身。要把载具接进 CI＝
  **动 `.github/workflows/ci.yml` 的触发档＝契约级 ⇒ 本腿停手上报，一个字未动**（票面禁区第 1 条）。
  ⛔ 别把这条读成"case 白写"：CI 的守卫与笔记本的载具是互相同一面的两枚，任一枚被删另一枚会漂；
  但"漂了没人知道"的窗口确实存在，这就是本腿交出来的缺口，不是妥协。
- **ⓐ/ⓑ 的最终选择归编排者翻勾**：本腿只交两形代价（8.6 第 1 条）＋ⓑ 的判据。
- **没测到的**：
  - `internal/winsec` **真被拆包那天**的同笔提交体验（载具只能造宇宙，不能替人写 pin）；
  - winsec 那族的**真机 `GOOS=windows` 腿**（禁区：本腿不许跑真 `go test`，`248-v1b` 在整包跑 `cmd/wisp`）；
  - 冷模块缓存下的 stderr 形状（票 250 已结案，本腿不复跑，也⛔ 不据它改口径）；
  - `--scope=census` 在**真 go** 下的名册（本腿只在假 go 下跑 case 15；真跑要 `go list ./...` 全树，
    属可读但本腿没排到，写〔待验〕⇒ ⛔ 不许据此排除任何候选形）。
- **一处文案未改的判断**：GUARD 3 的红句把 `scope=[${target}...]` 拼在自己嘴里（突变 D 实测会跟着漂）。
  要不要改成引用子脚本原文＝**动 `scripts/winsec-tests.sh`**，而本腿对该文件的判断是"形状对、不必动"，
  且没有第二条腿在场读它 ⇒ 本腿不自作改，报这里等裁。

## 8.8 待裁：AC#3 unknown 消息里那份 known 列表（4 名 → 5 名）算不算射程

⛔ 本腿不选。现量三处＋两形代价如下。

**现量**：

- 唯一生成点＝`scripts/portable-tests.sh:246`
  `echo "portable-tests.sh: unknown --scope=$mode (known: ${tiers// /, })" >&2`
  ⇒ 那份名册**今天已经**是 `tiers=`（`:192`）的派生物；"扩到 5 名"不是有人改文案，是多一枚档的必然后果。
- 名册与分支的一致性另有一把尺在守：`--scope=census` 读自己文件的 `case $mode in` 分支（`:277-295`），
  分支集 ≠ 名册集 ⇒ **拒答 rc=1**（票 251 AC#4 的形状）。
- 文本尺：载具 case 17 `unknown-scope-still-hard` 逐字断言那串 5 名（`rc=2` 不许软＝票面禁区第 2 条）。
- 射程尺（显式根、`| wc -l` 形）：串 `core, windows, cli, winsec, census` 命中文件枚数 =
  `scripts` **1**（case 17 的断言）· `docs` **1**（`docs/reports/pending-and-issues.md:10658`＝编排者自己的读数）
  · `tools` **0** · `cmd` **0** · `internal` **0** · `.scratch` **9**（工单与证据件）。
  产码里读这句文案的消费者：**0 处**（尺：`grep -rn 'unknown --scope' scripts docs tools cmd internal .scratch` 里
  只有生成点 1 枚＋断言 1 枚＋台账/工单若干）。

**甲形（"消息文本属射程、不许扩"）的代价**：

1. 要让 5 枚档只报 4 名，就得把 `:246` 的 `${tiers// /, }` 换成硬写四名 ⇒ **重新造出"两份名册可以各说各话"
   那一枚形状**，正是票 251 AC#4 当场拆掉的病（`portable-tests.sh:184-192` 的注释具名记着它当时的读数：
   census 报 `CLAIMED BY ... winsec` 而没有分支实现它、还 rc=0）。
2. 或者保留生成式但把第五档从名册里摘掉 ⇒ 名册与分支不再相等 ⇒ `--scope=census` **从此拒答**（`:283-295`），
   等于拿一条 CI 步骤换一句文案不动。
3. 副作用：`unknown --scope=winsec`（票 251 的立票那一发，台账 `:10505` 记的就是 4 名旧串）从此**不点名第五档**
   ⇒ 响亮失败退化成"你不知道有这一档"。这与票面"rc=2 不许软"是同向还是反向，归编排者判。

**乙形（"列表随档自动生成、扩名不算改动"）的代价**：

1. 那句文案**不再是定长串**：case 17 的逐字断言必须随档改（今天已随票 251 改过一次，台账里留着 4 名旧串）
   ⇒ 扩一名会把载具弄红一次——那是"要人看一眼"的代价，不是静默。
2. 台账/证据件里引过的旧串成为过期读数（现量 `docs` 1 枚、`.scratch` 9 枚）⇒ 只能**追加**更正、⛔ 不改写。
3. 零产码代价：现物**已经是**乙的形状。

⇒ 本腿落哪一形都不必改产码：甲要动 `:246` 一句（并连带把 census 的形状改掉或接受它拒答），乙动 0 句。
⛔ 选哪一形由编排者裁；本腿没改 `:246`、没改 case 17 的断言、没动 `rc=2`。

## 8.9 提交清单（逐枚：号＋pathspec＋numstat）

| 号 | pathspec（提交时显式写的） | numstat |
|---|---|---|
| `6bcb934a` | `scripts/portable-tests-selftest.sh`、`.scratch/wisp/probes/254/r1b` | 15 文件／+2427／−5；其中 `scripts/` 只有那一行 `231 5 scripts/portable-tests-selftest.sh` |
| 本文件所在那一枚 | `.scratch/wisp/probes/254/r1/verdict.md`、`.scratch/wisp/probes/254/r1b/logs/prefix-full-at-d253703a-parent.txt`、`.scratch/wisp/probes/254/r1b/logs/carrier-all-at-HEAD.txt` | 由 `git show --numstat` 在下一程取；⛔ 本腿不给自己编号（编过的号在台账里出过事） |

⚠ **两处必须记着的自纠**：

1. `prefix-full-at-d253703a-parent.txt` 的第一版是**半途夭折**的读数——本腿在它跑到一半时改了载具
   （把 case 19 改成零参数那一支），bash 逐行读脚本 ⇒ 那一版不可用。它已随提交 `6bcb934a` 进了库
   （156 行那一版），本腿**不删文件、不改历史**，只在同一路径上以 `>` 重跑覆盖
   （覆盖后：`27 case(s) ran, 9 assertion(s) failed`、`rc=1`＝本节 8.2 引的那一发）。
   ⇒ 引这一发只认覆盖后的那份；`git show 6bcb934a -- <该路径>` 与本腿后一枚提交的差异即两版状态。
2. 8.5 那把残留尺的**词面**被我写进了被扫文件 ⇒ 第一次跑它对本文件命中 1（命中的就是我自己写下的那串词），
   而**不是**文件里真有未填的空节。改掉词面、只留计数后才是 0。⇒ 这一枚与本仓 10-02 已知的"把尺面抄进
   被扫件"同形，本腿复踩了一次，写在表里而不是抹掉。

