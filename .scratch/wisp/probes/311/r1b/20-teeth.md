# 20-teeth.md — 票 311 落地腿 `311-r1b` 格二：AC#2 那几发牙（ⓐⓑⓒⓓ）

台面（每发件头逐字自带 `workdir=`／`env_assignment=`／`command_line=`／尾行 `rc=`）：
`T_post`＝`git archive 0f888a4f` 的**仓外**导出树（＝改后那侧，四枚块都在）；`T_pre`＝`git archive 2567e3a9`（改前，只有 winsec 一枚块）；
`T_strip`＝`T_post` 的副本、按 ⓑ 摘块用，摘⇄复原每步带 `sha256`（件 `21-strip-proof.txt`）；`WT`＝共享工作树（**只**用于"拒绝在委派之前"那一形，见 §6 第 1 条）。
逐枚计数器与分母并排＝件 `41-shot-table.txt`（66 枚件、`bytes` 与 `rc=` 同行）；名册两把尺并排＝件 `42-winsec-roster.txt`。

## 1. ⓐ 换卷＝起测试**之前**就拒绝（每枚入口一发＋两枚变体）

| 入口（件） | 命令行（逐字在件头） | rc | 出声的那枚前缀 | `=== RUN` | 分母／同发证据 |
|---|---|---|---|---|---|
| winsec `20-A1` | `bash scripts/winsec-tests.sh ./internal/winsec/` | **2** | `winsec-tests.sh:` | **0** | 同发那句 `ticket 310 … GOTMPDIR=D:/tmp/…/winsec-20261012-153120-830`＝根**真建在 D:**、GUARD 行 1 枚、件 1350 字节 |
| portable `20-A2` | `bash scripts/portable-tests.sh ./internal/winsec/` | **2** | `portable-tests.sh:` | **0** | 件 4343 字节（ledger 那几步行都在，`sh "$strict"` 那委派没走到） |
| wisp-cli `20-A4`（台面＝WT） | `bash scripts/wisp-cli-tests.sh` | **2** | `wisp-cli-tests.sh:` | **0** | 件 1353 字节；★该档自己的 `scope=./cmd/wisp/ dll_dir=…` 行**没出现**＝拒绝在委派之前，`cmd/wisp` 一枚测试也没起 |
| runtests `20-A3` | `sh tools/d22scan/runtests.sh -C tools/d22scan ./...` | **2** | `runtests.sh:` | **0** | 件 1340 字节 |
| runtests，offender＝调用方给的 `GOTMPDIR` `20-C3` | `GOTMPDIR=D:/tmp/…/gotmpdir-only sh … -C tools/d22scan ./...` | **2** | `runtests.sh:` | **0** | 同发 `ticket 310` 命中 **0**（分母＝那块整支没进：`GOTMPDIR` 非空即跳过建根）⇒ 守卫读的是**调用方**的根 |
| wisp-cli，同上 `20-C4`（WT） | `GOTMPDIR=… bash scripts/wisp-cli-tests.sh` | **2** | `wisp-cli-tests.sh:` | **0** | 件 1119 字节 |

退码形状＝裁形乙（命中换卷＝stderr 具名 GUARD 行＋`exit 2`），四枚同串，`allowlist/thresholds/golden` 一律⛔ 动。

## 2. ⓑ 敏感性正控：把守卫摘掉，颜色**必须**回来（★没这一发＝那格没交）

摘除只在 `T_strip` 里做（`txt/mk-block.sh strip`，按内容锚删块）；每步 `strip_removed_lines=`**71**、复原后
`sha_tstrip_now` 与备份逐字相同、`restore_equals_pristine=`**yes**（`21-strip-proof.txt`，两枚摘除序列都在）。

| 摘了什么 | 件 | 结果（逐字） | 判语 |
|---|---|---|---|
| winsec 档，**四枚块全摘** | `22-B1-winsec-xvol-guard-removed` | `=== RUN=101 --- PASS=56 --- FAIL=2 --- SKIP=0`、`rc=1`、`TempDir RemoveAll cleanup` **4 行**、`Access is denied` **4 行**、红名＝`TestGateJudgesThePrivateSetByResolvedSID`＋`TestSealNarrowsAndNamesThePrincipalItRemovedBySID` | ★这段守卫⛔ 是装饰：摘掉之后票 311 现量那对颜色**一枚不差地回来**；买到的＝"从这发红里把结论救回来" |
| runtests 档，**只摘它自己那一枚** | `22-B3-runtests-xvol-guard-removed` | 同 `PASS=56 FAIL=2`、`rc=1`、cleanup **4 行** | 单变量＝那一枚块；最深那档入口同样"摘了就红回来" |
| portable 档，**只摘它自己那一枚** | `22-B2-portable-xvol-own-guard-removed` | `rc=2`、`=== RUN=`**0**、GUARD 行还在但**换了一枚名字**＝`runtests.sh: ticket 311 AC#1 shape B: GUARD`（offender 里逐字带着 `GOTMPDIR=D:/tmp/…/portable-20261012-154030-1620(volume D)`） | "这一枚入口自己会不会说话"是被测出来的：摘掉 `portable` 那块，本入口那句就⛔ 再不属于本入口——拒绝下沉一层、由**子级**具名。⇒ 分层防住，⛔ 一次防 |
| cli 档 | —— | **⛔ 跑**（§6 第 1 条） | 只交了在场那发（`20-A4`／`20-C4`）＋静态尺 |

摘除侧⇄在场侧的**名册/四数并排**（同一形状两发）：
`26-B1x`（摘除，默认根）`=== RUN=101 PASS=58 FAIL=0 --- SKIP=0 rc=0` ⇄ `30-G4`/`30-G5`（在场，默认根）逐字同串；
`26-B1y`（摘除，portable 默认）`PASS=58 FAIL=0 rc=0` ⇄ 在场侧＝委派路那行 `portable-tests.sh: four numbers … PASS=58 FAIL=0`（在 `30-G4`/`30-G5` 件里，⚠ 射程具名：portable **自己**默认形状⛔ 独立起过一发，见 §6 第 3 条）；
`26-B3x`（摘除，runtests）`top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77 rc=0` ⇄ `30-G3`（在场）逐字同串。

## 3. ⓒ 三形不许破：在场 ⇄ 摘除，两发**同色**（守卫本来就不该在这三形说话）

| 形 | 入口 | 在场（件／色／守卫出声枚数） | 摘除或改前那侧（件／色） | 同色？ |
|---|---|---|---|---|
| ① 出厂默认（两枚变量都不设） | winsec | `30-G4`+`30-G5` `58/0 rc=0` GUARD=INCONC=**0** | `26-B1x` `58/0 rc=0`；另改前 `10-winsec-pre-run1/2` `58/0 rc=0` | ✅ |
| ① | portable | 委派路 `30-G4` 内 `portable-tests.sh: four numbers … 58/0`＋GUARD 0 | `26-B1y` `58/0 rc=0` | ✅（射程见 §6.3） |
| ① | runtests | `30-G3` `PASS=35 RUN=77 rc=0` GUARD 0 | `26-B3x` 同串 | ✅ |
| ① | wisp-cli | ⛔ 跑（§6.1） | ⛔ 跑 | —— |
| ② 同卷不同目录（`WISP_TEST_GOTMP_ROOT=C:/Users/swq/AppData/Local/311r1b-samevol`） | winsec | `26-D1` `58/0 rc=0` GUARD 0，同发 GOTMPDIR 逐字指那枚 C: 根 | `26-C1` `58/0 rc=0` | ✅ |
| ② | portable | `26-D2` `58/0 rc=0` GUARD 0 | `29-P7`（`T_pre`，那棵树 portable⛔ 有块＝等价摘除）同串 | ✅ |
| ② | runtests | `26-D3` `PASS=35 RUN=77 rc=0` GUARD 0 | `29-P5`（`T_pre`）同串 | ✅ |
| ③a `LOCALAPPDATA` 为空 | winsec | `27-E1` `58/0 rc=0`，GUARD/INCONC=**0**，`ticket 310 AC#1` 命中 **3**（委派链上三枚入口各打一行 DEGRADED，逐字同句首） | `29-P2`（`T_pre`，winsec 有块／下两枚⛔）命中 **3**、同色 | ✅ 不多一行不少一行 |
| ③a | portable | `27-E3` `58/0`，310 命中 **2**（portable＋runtests），GUARD 0 | `29-P6`（`T_pre`）310 命中 **2**、同色 | ✅ |
| ③a | runtests | `27-E2` `PASS=35 RUN=77 rc=0`，310 命中 **1**，GUARD 0 | `29-P1`（`T_pre`）310 命中 **1**、同色 | ✅ |
| ③b `go env GOOS`⛔ windows | winsec | `28-F1` `rc=2`（它自己那句 windows-leg GUARD），310 命中 0、311 命中 **0** | `29-P4`（`T_pre`）`rc=2`、311 命中 0 | ✅ |
| ③b | wisp-cli | `28-F2`（`T_post`）`rc=2`、311 命中 **0**（分母＝同发那句 `GUARD - this is the windows leg of the cmd/wisp gate`） | 摘除侧⛔ 另跑：那一支在 GOOS 守卫就 `exit 2`，两枚台面行为同一把尺可读（块前后都⛔ 起测试） | ✅ |
| ③b | runtests | `28-F3` `rc=1`（Go 跑⛔ 了非本机目标），310 命中 0、311 命中 0 | `29-P3`（`T_pre`）`rc=1`、同样 0 命中 | ✅ |

★`INCONCLUSIVE` 支那发＝`25-C2`：`runtests.sh` 打了 **2 行** INCONCLUSIVE（`GOTMPDIR` 与 `WISP_TEST_GOTMP_ROOT` 各一枚，都是裸相对名），
块**⛔ 改退码**——那发 `rc=1` 来自 Go 自己（逐字 `go: creating work dir: GetFileAttributesEx 311r1b-relbare/…: The system cannot find the path specified.`），
⛔ 是守卫响。这一形是已登记欠账 `Q-311-r1-1`，本腿⛔ 扩（派单第 3 节硬约束）。拼法正控＝`25-C1`（`D:\tmp\…` 反斜杠形）GUARD 读到 volume D、`rc=2`。

## 4. ⓓ 顺序：`3b857e53` 那段注释指名的两枚件，今天真在盘上

- `23-order-degraded-then-guard.txt`（`T_post`，`WISP_TEST_GOTMP_ROOT=/d/tmp/311r1b-notadir`＝**一枚真文件**站在那儿）：
  行序逐字＝`winsec-tests.sh: ticket 310 AC#1: DEGRADED - cannot create /d/tmp/311r1b-notadir/winsec-…` → 三行 `ticket 311 AC#1 shape B: GUARD …` → **没有一行业务测试**（`=== RUN` 0，件 1332 字节）→ `rc=2`。
- `24-order-guard-after-guard1.txt`（有效换卷根＋非法 scope）：`ticket 310 … GOTMPDIR=D:/tmp/…` → `GUARD 1 - the package list …` → **块⛔ 出声**（311 命中 0）→ `rc=2`
  ⇒ 311 那块真在 GUARD 1 **之后**（这正是"三形不许破"里 DEGRADED 仍能先出声的前提）。
- `24b-order-degraded-before-guard1.txt`（`LOCALAPPDATA=`＋非法 scope，`T_post`）：`ticket 310 AC#1: DEGRADED - LOCALAPPDATA is empty …` → `GUARD 1 …` → `rc=2`
  ＝台账 `A862` ⑤ 那一发在本腿台面上的复现，也是 10 号件 §4 那句"跑在 GUARD 之前"的真机凭据。
注释已由 `e4e787fb` 改指这三枚件，并具名写了"r1 那两枚从没落盘"（死腿遗产那一格销掉）。

## 5. 每一发都自带"坏世界真被造出来"那一行

ⓐ 六发的坏世界＝同发里 310 那块**真把 GOTMPDIR 定在 D: 上**那一行＋盘上真目录（逐枚根名＋`CreationTime`＝`51-residue.txt` §1，其中两枚含锁死叶）；
ⓑ 两发的坏世界＝4 行 `TempDir RemoveAll cleanup: openfdat D:/tmp/… : Access is denied` 逐字在件里；
ⓒ 那些"守卫不出声"的发⛔ 当坏世界证据，它们的凭据是**颜色等于基线**＋出声枚数 0 旁边站着同发四数行。
唯一一枚按定式判**作废**的是 `logs/51-residue-raw.txt` 第一版 `unreadable_leaf_count=0`——假阴性（只探三层，锁死叶在第四层），
作废句写在 `51-residue.txt` §7，替换尺＝`ps1/residue2.ps1`（`52-locked-leaves-raw.txt`：本腿 12 枚，正控＝同一把尺在别腿的根上读到 7/6/6/6/1 枚）。

## 6. 本腿⛔ 做完的，具名（⛔ 把没跑的读数写成跑过）

1. **wisp-cli 的形①／形②／摘除侧没跑**。两枚墙：形①②会 `bash portable --scope=cli` ⇒ 真起 `cmd/wisp` 档（本波 `cmd/wisp/**` 归 `295-r1`，在飞）；
   `T_post` 导出树里 `third_party/sherpa-onnx` ⛔ 在（尺＝`git ls-files third_party`＝**0** 枚 tracked）⇒ 该入口会先死在 DLL 守卫上、够不着我这枚块。
   交出去的＝在场两发（`20-A4`／`20-C4`，都证明"拒绝在委派之前"）＋形③b 一发（`28-F2`）＋静态尺（块内只有 4 行 echo，全在 cross/inconclusive 分支里，`grep -c 'ticket 311 AC#1 shape B'`＝**4**×四枚，件 `43` §3）。
2. **cli 档与 tools/d22scan 档的换卷颜色本机⛔ 量过**（票面那对颜色只属于 winsec 档）⇒ 块里那句 gate 措辞我保持 winsec 原文，⛔ 给别档编一对颜色（理由与尺在 `10-replicate.md` §2）。
3. **portable 档"默认形状"我⛔ 独立发过一发**：`--scope=census` 在 `:758` 之前就返回（够不着块），独立发默认形状＝整包 suite。
   我用的是委派路（`30-G4`/`30-G5` 件里 portable 自己那行四数）＋同卷那一发（`26-D2`）＋摘除侧（`26-B1y`/`29-P7`）三枚合起来覆盖，射程已具名。
4. `runtests` 那枚形③b 的 `rc=1` 与 `25-C2` 的 `rc=1` 都**⛔ 是档位判定**（前者 Go 跑不了非本机目标、后者相对根解析不到）；我把它们当"形状没被守卫改坏"的读数用，⛔ 当"门绿"。
5. 我自己的两枚仪器缺陷（相对件名在 `cd` 后写／块锚当动态正则）＝10 号件 §5；两枚都在**落件之前**被抓，`grep -L '^rc='` 与 0 字节自查＝件 `55-audit.txt`（66 枚件，`rc=` 缺 0 枚，0 字节 0 枚）。
6. 票面那条**未复核**读数（`A860`/`A862` 记的"BA 在令牌里是 deny-only"）本腿⛔ 碰过、也⛔ 复核得了——它⛔ 属本票三格，留在编排者名下。
