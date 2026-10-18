# 票 341 落地腿 341-r1 report（写腿·本波唯一写 `.go` 的腿）

## 0. 起手读数（我第一笔现跑，原样）

```
git log -1 --format='%H %s'
98dbf20a35027a9290f52649d8ab2634646cb14e 第三波派单写档＝4 枚腿（先存档再发）｜放出前逐面现量…
git status --porcelain -- internal/config cmd/wisp
<零行>
```

★我的基线＝**`98dbf20a`（我自己现认的）**，⛔ 派单里编排者写的号。测量期间 HEAD 漂到 `062dbe44`：漂的 3 笔
全部现验**不落我写面**（`git diff --name-only 98dbf20a..HEAD -- internal/config cmd/wisp`＝零行），
且八枚相关件 blob 同一性 same=YES（`logs/01-start-gauges.txt` gauge 1）。
写面冲突尺：`find internal/config -name '*.go' -newermt '-45 minutes'` 只回我自己那一枚
（`logs/12-roster-strict.txt` 末段）⇒ **本波该面零别腿落笔**，⛔ 读／⛔ 提交／⛔ 采信任何对方字节（无对方）。
票面框尺＝未勾 **7**／已勾 **0**（起手⇄交件逐字相同）＝**本腿一枚框没翻**。

## 1. 提交名册（严格主题前缀尺，⛔ 子串尺）

| 提交 | 内容 | name-only 越界（对我允许面之外） |
|---|---|---|
| `6948eb08` | `internal/config/writeguard_unnailed_341_test.go`（numstat `389/0`） | 0 |
| `fd82f786` | `logs/` 33 枚证据件 | 0 |
| 本笔所在 | report＋`logs/16`＋16 枚台件脚本 | 0 |
| 后续 | 票 341 票尾一条 bullet | 0 |

★名册尺自指警告：本表写时"后续"两笔尚不存在 ⇒ 终尺＝`logs/12-roster-strict.txt` 头两行那把
（`git log --perl-regexp --grep='^票 341 落地腿 341-r1·' <anchor>..HEAD`）＋交件侧现跑。
⚠ **两处具名更正（本腿自己的缺陷，⛔ 当读数）**：
1. `fd82f786` 主题写"48 枚"，该笔实含 **33** 枚（余 15 枚台件在本笔）＝枚数写坏，⛔ 改历史，在此具名。
2. `logs/11-gates-final.txt` 的越界读数 16 枚＝我这把名册尺用了**子串** `--grep='票 341 落地腿 341-r1'`，
   把编排者的派单笔（主题含 `341-r1`）吃了＝正是派单"三键"警告的那形。改严格前缀后越界＝**0**（`logs/12`）。
3. `logs/00-start.txt`＝起手尺第一发的**缺陷件**（`$'\r'` 嵌双引号串破坏解析＋`cut -c` 截多字节 ⇒ 件非 UTF-8），
   原件留盘⛔ 删，重跑＝`logs/01`。同病灶第二犯＝`logs/10-…-rulerdefect.txt`、`logs/11-…-rulerdefect.txt`。
4. `logs/01` gauge 3 的 **P4 命中 0＝尺缺陷**（转义写坏）；改断言形锚重量＝载体 1 枚
   （`writeguard_handedit_227_test.go:92`，用例起 `:70`），P3/P6 的 0 同尺配正控后**才是读数**（`logs/09`）。

## 2. AC#0 起手现量（自加尺，⛔ 抄 227-v1 的表）

逐格问"有没有任何 `*_test.go` 真断这一形的**结果**"，射程＝锚 `98dbf20a` 的全部 `*_test.go`：

| 格 | 命中枚数 | 载体 file:line | 与票面/227-v1 表的关系 |
|---|---|---|---|
| P0 改值 | 4 处 | `manager_test.go:60`/`:270`、`writeguard_226_test.go:37`、`writeguard_handedit_227_test.go:25` | 同 |
| P1 删非默认标量行 | 1 | `writeguard_handedit_227_test.go:31`（用例起 `:23`） | 同 |
| **P2 删值本就是默认的行** | 1（路过形） | `writeguard_handedit_227_test.go:184`（用例起 `:181`）＝只断认领形状，⛔ 断 P2 独有的"补回值＝内存值⇒看不出差别" | ★票面第 2 行把 P2 归〔只有台件〕；本腿现量与之**相容**（同 227-v1 的判定），本格判点确实⛔ 在场 |
| **P3 删整节** | 0 | 无（唯一近似命中 `cmd/wisp/config_reload_223_test.go:353`＝`r.plant` **写入**一节，⛔ 删一节；`logs/13` (a) 逐字读出） | 同 |
| P4 删 map 条目 | 1 | `writeguard_handedit_227_test.go:92`（用例起 `:70`） | 第一发尺转义坏⇒假 0，已更正（`logs/09`） |
| **P5 删 schema_version 走写支** | 0 | 名册里只有 `:110` 那一枚删该行，它走 nothing-written 支 | 同 |
| P5′ 走 nothing-written 支 | 1 | `writeguard_handedit_227_test.go:110`（用例起 `:106`） | 同 |
| **P6 删 provider base_url** | 0 | 无（同形正控 P5′ 回 1 ⇒ 这把尺活着） | 同 |

⇒ 票面"现量"四行**逐行复现、零处需更正**；未覆盖四格＝P2／P3／P5／P6，与本腿要做的四格一致。

## 3. AC#1 交付主体：四格逐格表（红名册逐名对撞用这张）

载体文件＝`internal/config/writeguard_unnailed_341_test.go`（389 行／5 枚 `^func Test`／CR 0／blob `c577f1cc`）。
"未突变读数"＝仓外导出树、`-count=1 -timeout 30m -v`、两把临时根台面（默认 C 卷 ⇄ `D:` 预建）逐枚同值。

| 格 | 用例名（载体行） | 未突变读数 | 去牙落点（⛔ 与他格相同） | 去牙红句（逐字，件内可 grep） |
|---|---|---|---|---|
| **P2** | `TestHandDeletedDefaultValuedLineIsRestoredInvisibly` `:58` | PASS（两面） | `writeguard.go:174-175` 摘掉 **INFO 支**那句删行人话（WARN 支同句保留＝存活锚计数 1 未变） | `writeguard_unnailed_341_test.go:111: the INFO branch must state the delete-shape even when the values hide it; log:` |
| **P3** | `TestHandDeletedWholeSectionComesBackFromDefaultsNotMemory` `:138` | PASS（两面） | `writeguard.go:141` `base = deepCopyConfig(disk)` ⇒ `(m.cur)`（合并基准退回内存快照＝票 226 的原始缺陷形） | `:162: the resurrected section carried memory's value instead of the schema default:` ／`:174: the size divergence was not reported; log:` ／`:180: a reload reads size=60 click_through=true, want 56 and true` |
| **P5** | `TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed` `:209` | PASS（两面） | `loader.go:257` `cp.SchemaVersion = SchemaVersionCurrent` ⇒ `= 1` | `:231: the write branch did not re-stamp the deleted schema_version line at the current version:` |
| **P6** | `TestHandDeletedProviderBaseUrlComesBackAsThePresetValue` `:289` | PASS（两面） | `loader.go:213` `p.BaseURL = pre.BaseURL` ⇒ `= ""`（preset 填充失效） | `:320: the deleted base_url line came back as something other than the built-in preset:` ／`:330: a reload reads "", want the preset endpoint the file now carries` |
| P6 正控 | `TestHandDeletedBaseUrlOfAnUnknownProviderStaysEmptyInTheFile` `:349` | PASS（两面） | 同 P6 落点 ⇒ **该发它必须仍绿**（实测：p6 那发红名册＝`TestLoadFilePresetInheritance`＋P6 两枚，本正控⛔ 在列＝绿） | —（正控无红句＝读数本身） |

四发落点＝两枚文件四行不同（`writeguard.go:174`／`:141`、`loader.go:257`／`:213`）；每发落刀前锚命中＝1、
落刀后＝0、新锚≥1、存活锚计数不变；逐发当场还原并回报还原哈希＝pristine：
`writeguard.go efe12661…47bd`（p2/p3 两发）、`loader.go 40cb3e12…7bdb`（p5/p6 两发），件＝`logs/05-mu-*.txt`。

## 4. AC#4 颜色与门禁

四数（⛔ 混用⛔ 相减；尺名＝`RUN=^=== RUN`／`PASS_TOP=^--- PASS`／`FAIL_TOP=^--- FAIL`／`SKIP=^--- SKIP`）：

| 那把尺（`./internal/config/`，`-count=1 -timeout 30m -v`） | RUN | PASS_TOP | FAIL_TOP | SKIP |
|---|---|---|---|---|
| 基线树（锚 `98dbf20a` 原样导出）·台面 A 默认 C 卷 | 190 | 109 | 0 | 0 |
| 基线树·台面 B `D:` 预建 | 190 | 109 | 0 | 0 |
| 交件树（基线＋本腿 5 枚用例）·台面 A | 195 | 114 | 0 | 0 |
| 交件树·台面 B | 195 | 114 | 0 | 0 |
| 去牙 p2（INFO 支那句） | 195 | 113 | **1** | 0 |
| 去牙 p3（合并基准＝文件） | 195 | 102 | **12** | 0 |
| 去牙 p5（SaveFile 钉版本） | 195 | 112 | **2** | 0 |
| 去牙 p6（preset 填 base_url） | 195 | 112 | **2** | 0 |

- 红名册**逐名**作差（基线⇄交件）＝**零枚新增红**；交件 195−190＝本腿 5 枚、114−109＝同 5 枚（⛔ 相减当读数，此处只作对照）。
- ★第一发交件树是 **2 枚红**，红句逐字在 `logs/06`／`logs/07`：本腿用例把 canonical 序列化的引号形状写成双引号，
  而本仓输出**单引号字面串**（`b2` 台件 `base_url = ''`、既有 `settings_313_test.go:270` 同形）＝**用例缺陷**，
  改后重跑＝上表；⛔ 为此放宽任何断言、⛔ 改产码。
- p3 那发 12 枚红＝摘掉"合并基准来自文件"这枚守卫后，226／227 的既有用例与本腿 4 枚一起红＝
  **守卫被摘**的合法读数（⛔ 计入"我的字节把包打红"，那发在仓外树）。p5/p6 各 2 枚：另一枚分别是
  `TestTicket257R1AC1CleanMachineRefusesAllSevenFields`、`TestLoadFilePresetInheritance`＝既有常驻用例，
  同属"守卫摘掉该红"，⛔ 本腿新造。
- 门禁四数（只在共享工作树，跑前 `unset GOTMPDIR TMPDIR`，逐枚 `rc=`）＝`logs/11-gates-final.txt`：
  gofumpt `-l` 我的件＝空输出 **rc=0**；`go vet ./internal/config/` **rc=0**；`sh scripts/d22scan.sh`
  末行逐字 `d22scan: clean - no D22 ban violations` **rc=0**；`sh scripts/check-path-length-budget.sh`
  **rc=0**（`VERDICT GREEN`，分母 16056 tracked／over-budget 60 全在册／not in roster 0）。
  ★`gofumpt -l internal/config` 本轮**空输出**＝派单说的"6 枚 CRLF 存量脏"此刻不在该面包现（读数，非我改动）。
- EOL 自证（派单"我的新字节必须 LF"）＝`logs/16-eol-proof.txt`：工作树 CR **0**、committed blob CR **0**、
  `git archive` 重现的检出件 CR **0**、三处 gofumpt rc=0；对照 227 配方件同样 CR 0。
  ⚠ 本机 `core.autocrlf=true` ⇒ git 会打"LF will be replaced by CRLF"警告；blob 里存的是 LF（现量）。
- ⛔ `./cmd/wisp` 整包（本波该面零写腿＋单车道，持跑者＝`338-v1`）；⛔ `go test ./...`＝本腿零跑零判。

## 5. AC#2（⛔ 把"没覆盖"说成"覆盖"）

- 四格各一枚以上真断言，⛔ 放宽／⛔ 删用例／⛔ `t.Skip`（文件里 `t.Skip` 命中 0，占位符尺 0 命中＝`logs/15`）。
- ★P2 与 P1 **合并＝造谎**这条，本腿给出机读差分（不是嘴上说不同）：
  - 摘掉"合并基准＝文件"（P3 那发）⇒ P1/P3/P4/P6 与 226 用例全红，**P2 用例仍绿**（`logs/15` 末段计数 0）；
  - 摘掉 INFO 支那句（P2 那发）⇒ **只红 P2 一枚**（FAIL_TOP＝1）；
  - ⇒ 两形被测的是**不同的支**（P2＝`len(diverged)==0` 的 Info 支 `writeguard.go:171-176`；P1＝Warn 支 `:184-190`）。
- **缺世界＝零枚**。四格今天的产码下都造得出可断的世界，⛔ 假绿填空。唯一具名的"造不出"是**形状本身没有**：
  P2 那一形在报告侧**没有任何可指认的差别**（这正是它的判据），所以本腿把它钉成"那句仍要说清丢了什么"＋
  "⛔ 冒出 divergence 清单"的反向断言，⛔ 假装它可被值侧观测。

## 6. AC#3 那句说过头的句子（只指认、⛔ 落笔）

现量尺＝`grep -c`／`grep -n`，取数时刻 HEAD＝`062dbe44`（★⛔ 动 `docs/evidence/s1/226-*.md` 一字）：

- **哪份件第几行**：`docs/evidence/s1/226-config-write-no-clobber-r1.md` **第 111 行**（sec 5「逐格判据」的 AC#5 那一格）。
- **原句逐字**：`附带一枚：文件里已经有这枚键要的值时**一个字节都不写**（不白重排、不白丢注释、不动 mtime）。`
- **按盘上现量该收成什么形**（建议编排者落笔的更正文字，本腿只给形）：
  "一个字节都不写"只在**该文件 schema_version 已是当前版**时成立；手删 `schema_version` 行或塞入旧版件时，
  `readConfigFile` 内部的迁移会先把文件重排成 canonical（bytes 变、mtime 变、操作者注释丢、落 `config.toml.bak-1`），
  即便 `mergeWrite` 走 `nothing written` 早返回支也一样＝产码 `writeguard.go:152-156` 那句自己承认
  （`…had already been re-laid in canonical form…`，锚 `98dbf20a` 处行号 `:154`／`:187`）。
  钉住它的常驻用例＝票 227 的 `writeguard_handedit_227_test.go:106` ⇄ 正控 `:150`，本腿再加 P5 那一枚
  （`writeguard_unnailed_341_test.go:209`＝写支同形）。
- ★**票面本身与盘不符之处（具名更正，＝派单起手第 2 步要的东西）**：票 341 现量第 3 行与本票 AC#3 说
  "sec 1 与 sec 5 那句不白重排不白丢注释"；盘上现量＝`不白重排` 与 `不白丢注释` 在 r1.md **各只 1 处、都在 :111（sec 5）**，
  sec 1（`:14-27`）里 `:22` 那句写的是"**值不丢、格式丢**"＋"注释与手排键序本来就不保留"＝**那句说实话、⛔ 说过头**。
  ⇒ 该更正只一处（:111），⛔ 按票面去追 sec 1 那句。另：v2 件 `:117`／`:162` 已把同一问题记成"要据此收窄／不成立"，
  与盘上产码现状一致（收窄文字已由票 227 落到 `writeguard.go:152-156`）。

## 7. 五条必答题 → 见 `report-answers.md`（同一提交内，⛔ 本件写不下；先例＝322-r2 的 `report-sec9.md`）

## 8. 本腿⛔ 做／交回的（具名）

1. **AC#3 的落笔**＝归编排者（写面⛔ `docs/evidence/**`），本腿只交sec 6 那条形。
2. **七格翻框**＝归非实现者（`341-v1` 之类）；本腿⛔ 动 `Status:`／票名／`-done`。
3. **P5 暴露的"迁移-on-读真写盘"行为变更＝乙形**＝票 345 射程，本票⛔ 顺手做（AC#5 明文）。
4. `cmd/wisp` 整包颜色＝本腿零跑零判（车道归 `338-v1`）。
5. `internal/config` 包级 `gofumpt` 存量脏此刻读空（sec 4）＝与派单"6 枚"不符，**读数具名**，⛔ 我改。
