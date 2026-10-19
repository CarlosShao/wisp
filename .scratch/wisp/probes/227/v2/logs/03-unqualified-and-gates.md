# 227-v2 读数件③＝两枚未定性的定性⇄★三处派单／台账前提更正＋门禁四数

## 定性尺＝成对两发（`scripts/10-paired-gotmpdir.sh`，同树同 `-run`、只差 GOTMPDIR 拼法）

| 发 | GOTMPDIR 拼法 | 四数 RUN/PASS_TOP/FAIL_TOP/SKIP rc | 两枚名各自 |
|---|---|---|---|
| 甲 | **反斜杠** `C:\Users\swq\AppData\Local\wisp-test-gotmp\227v2-pf-bs` | `2/2/0/0 rc=0` | 两枚**全绿** |
| 乙 | **正斜杠** `C:/Users/swq/AppData/Local/wisp-test-gotmp/227v2-pf-sl` | `2/0/2/0 rc=1` | 两枚**全红** |

⇒ 两枚红的根因＝**同一枚：GOTMPDIR 的拼法**（⛔ 卷、⛔ 长短）。★两枚都在**同一棵树、同一个 C 卷短根**上翻转＝形状敏感被完整摘回来了。

### ① `TestDirOfBaseOfCutOnThePlatformSeparator`（`internal/tools/pathshape_portable_test.go:99`）＝**台损**，附一枚真缺陷形状

- 红句逐字（乙发）：`dirOf("C:\\Users\\swq\\AppData\\Local\\wisp-test-gotmp\\227v2-pf-sl\\TestDirOfBaseOfCutOnThePlatformSeparator2408872237\\001\\notes.txt") = "C:\\…\\001", want "C:/Users/swq/AppData/Local/wisp-test-gotmp/227v2-pf-sl\\TestDirOfBaseOfCutOnThePlatformSeparator2408872237\\001"`。
- 机制（〔读码推的〕＋上面那发〔真机量的〕）：`os.MkdirTemp` 把 GOTMPDIR 前缀**原样带进** `t.TempDir()`（乙发的 want 串＝`C:/…/227v2-pf-sl` 正斜杠段＋`\001` 反斜杠尾＝混拼），而 `dirOf`/`baseOf` 的文档自己写明入参是 **ALREADY canonical**（`fs_write.go:177-183`：切在平台分隔符、切前 `unifySeparators`）⇒ 归一化后的 got ⇄ 未归一化的 raw want 必然不同串。甲发（反斜杠）下 raw＝归一化形＝绿。
- ⇒ 判**台损**：本票 AC#5 的红名册里它⛔ 产品缺陷（票 75 那族消费者断言在出厂拼法下绿；`A1053`⑵ 把它记在"台损 17"内＝方向对，⚠ 根因那句要改，见下更正Ⓒ）。
- ★同时具名一枚**真缺陷形状**（⛔ 我写面、交编排者决定要不要立票）：该用例把 raw `t.TempDir()` 当 canonical 入参用＝断言成立与否取决于宿主环境变量拼法；同包 `helpers_test.go:83 tempCanonical` 已有"把夹具先归一化再比"的现成形（`bridge_junction_windows_test.go` 那族都用了它）。立票形状一句话＝「`internal/tools/pathshape_portable_test.go:95-118` 那枚用例的夹具改成先过 `tempCanonical`（或等价归一化），断言不动＝消费者尺⛔ 再随 `GOTMPDIR` 拼法改色」。

### ② `Test20r2JunctionRefusalIsNeverAskable` ＝**台损**，并更正派单／编排者两处前提

- 红句逐字（乙发）：`refusal_20r2_windows_test.go:34: 前置条件缺失：造不出真 NTFS junction（mklink /J C:\Users\…\227v2-pf-sl\Test20r2JunctionRefusalIsNeverAskable2509583917\001\jn C:/Users/swq/AppData/Local/wisp-test-gotmp/227v2-pf-sl\Test20r2JunctionRefusalIsNeverAskable2509583917\002: exit status 1: Invalid switch - "Users".）`。
- ★本腿自己补打的四形 `mklink /J` 矩阵（`09-targeted.sh` A 段，短 C 卷根，⛔ D 卷、⛔ 长路径）：link⇄target 都反斜杠＝`Invalid switch - "bs-bs-link"`（本腿这条⛔ 带引号、cmd 把整串切开＝我这发的写法缺陷，⛔ 当读数用）；★**唯一与本腿整包同形的两发**＝link 反斜杠⇄target 正斜杠＝`Invalid switch - "Users"`，与编排者长 D 卷红句、与 `A1053`⑴ 那两句逐字同串 ⇒ `mklink` 在 `C:/…` 那形上被 cmd 的开关解析吃掉＝**拼法依赖坐实**（长 D 卷那族同理：那发 GOTMPDIR＝`D:/tmp/…` 正斜杠形，见编排者红句里 target 全是 `D:/tmp/…`）。
- ⇒ 判**台损**（凭据＝上面配对两发＋四形矩阵）；本票红名册里它⛔ 产品缺陷，且★甲发（反斜杠短根）下该枚 `--- PASS (0.08s)`＝它在正确台面形下是绿的。
- Ⓐ **派单前提更正（具名，⛔ 我改派单正文）**：派单第 14 行写「★`internal/panel` 同名那枚短路径下已绿＝两枚同名不同包、行为不同」⇄ 编排者 `A1053`③ 写「panel 同名列已绿＝两枚同名不同包」。本腿现量尺＝`git grep -n '^func Test20r2JunctionRefusalIsNeverAskable' HEAD` **只命中 1 枚**＝`internal/tools/refusal_20r2_windows_test.go:33`；`git grep -c 'Junction' HEAD -- internal/panel`＝**0**、`-- cmd/wisp`＝**0**。⇒ **全仓只一枚该名的用例、在 `internal/tools`；`cmd/wisp` 侧⇄`internal/panel` 侧那两枚⛔ 存在**——"两枚同名不同包"那形＝把同一枚（`internal/tools`）的长路径发⇄短路径发读成了两枚。旁证＝编排者自己那份 `attrib/final-attr.log:14` 紧跟那枚红打的包行是 `FAIL github.com/CarlosShao/wisp/internal/tools`。⛔ 据此退任何人：两处描述方向（"这枚要单独定性"）与结论（台损）都成立。
- Ⓑ `A1053`③ 那句"我红句没读全"（①那枚）本腿已读全并归口（上面①）。
- Ⓒ **根因更正**：`A1053`⑵ 那句"我的长临时根让 `mklink` 建不出来"——★`mklink` 那族的失败⛔ 系于"长"或"D 卷"，系于**拼法**（正斜杠 GOTMPDIR）；短根＋正斜杠照样红（本腿乙发）。同一条更正适用于 junction 那 9 枚⇄SID 那族（甲发＝反斜杠短根下逐名全绿）。**编排者归的类（台损）⛔ 变，只是可复现的修法要改写**：定式＝`GOTMPDIR` 一律 `cygpath -w`／反斜杠形，与 PATH 那条同批守。
- Ⓓ 一处⛔ 归口＝`A1053`⑦ 那枚 `TestAC2SealNoticeLandsInTheRunLegLogFile`（"长路径形状敏感"）：本腿甲发＝`--- PASS (0.10s)`＝与Ⓒ同源（它那发的 `attrib` 里 GOTMPDIR 是正斜杠形）；⛔ 推翻它"⛔ 落地引入"那句，只补一句：根因与Ⓒ同。

## 门禁四数（★共享工作树、先 `unset GOTMPDIR TMPDIR`＝A1027⑤ 定式）

| 门 | 命令 | rc | 判读 |
|---|---|---|---|
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l <我的件>` | — | ★**⛔ 对象**：本腿写面＝`.md`＋`.sh`＋`.ps1`，全史 `.go` 命中 **0**（逐枚 `git show --name-only` 尺见件④）⇒ 那把尺今天无物可量，具名登记、⛔ 编数 |
| d22scan | `sh scripts/d22scan.sh` | **0** | 末行逐字 `clean - no D22 ban violations …`（`gate-d22scan.log`） |
| 长度门 | `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN - every over-budget tracked path is rostered …`（`gate-len.log`） |
| go vet | `go vet ./cmd/wisp/` | **0** | 本腿⛔ 改产码，这发＝证明我读的包在共享工作树上仍干净（`gate-vet.log` 空输出） |

★门禁跑前现量＝`unset GOTMPDIR TMPDIR` 后 d22scan rc=0（⛔ 带 env 那形 rc=2 射程拒绝＝票 311 GUARD，⛔ 读成红）。

## 本腿自报三枚（⛔ 圆场）

1. **★第一发整包（full-A）＝本腿自造台损**：导出树命名 `tree` ⇒ 整条路径每段≤8 字符、无需别名 ⇒ `GetShortPathName252` 那族 8 枚前置红（多算 8 枚、且把 `TestTicket252R1Alignment…` 的父段问题盖住）。改树名重打＝full-B。那发四数（2537/1706/27/13）入库⛔ 当 AC#5 读数用、只作"台损形状"的凭据。
2. **`grep -c $'\r'` 那把尺本腿打坏一次**：`$'\r'` 在 `$(…)` 里被吞成空模式＝命中数恰好等于行数（一度让本腿以为"入库字节带 CR、违反公共块 LF 条"）。改用 `tr -cd "\r" | wc -c` 重打＝**6 枚已入库 blob CR 字节全 0** ⇄ 工作树同 0（`core.autocrlf=true`＋`.gitattributes` 里 `text: set / eol: lf`）⇒ 新字节确实 LF；那枚"本腿缺陷"⛔ 成立、⛔ 撤，具名留下坏尺本身。
3. **`09-targeted.sh` A 段那条四形矩阵**：本腿写 mklink 探针时⛔ 带引号，甲形（两侧都反斜杠）也被 cmd 切开报 `Invalid switch`＝那一行的读数⛔ 能用来判"反斜杠形是否可用"（本腿的整包⇄配对台跑的是 Go 侧 `exec.Command("cmd","/c","mklink","/J",link,target)` 传参、⛔ 手工拼串，那才是与用例同形的凭据）。Ⓐ 段结论只取"target 为正斜杠 ⇒ `Invalid switch - "Users"`"那一格（与编排者红句逐字同串），甲行作废并就地具名。
