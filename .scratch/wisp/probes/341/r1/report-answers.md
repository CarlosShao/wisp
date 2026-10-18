# 票 341 落地腿 341-r1 report 分件：五条必答题（⛔ 本腿自己盖章算，逐条给盘上尺）

取数时刻 HEAD＝`062dbe44`（起手锚＝`98dbf20a`，两者之间 `writeguard.go`／`loader.go`／`migrate.go`
blob 同一性 same=YES＝sec 4 表 `logs/12-roster-strict.txt` 末段现量）。台面＝`D:/tmp/wisp/20261018/341r1/`
的导出树（base⇄deliv⇄mu_*），⛔ 在共享工作树跑 `go test`。

## (i) 四格逐格：去牙那一发红在哪一行、红句逐字；两发落点分别是什么

载体文件＝`internal/config/writeguard_unnailed_341_test.go`（blob `c577f1cc`，389 行，CR 0）。
"红在哪一行"＝**用例里那枚判点断言的行**；"落点"＝**产码里被摘的那一行**。四格落点互不相同
（两枚文件、四行：`writeguard.go:174`／`writeguard.go:141`／`loader.go:257`／`loader.go:213`）：

| 格 | 产码落点（摘什么） | 用例判点行 → 红句逐字 | 该发红名册（全量） |
|---|---|---|---|
| P2 | `writeguard.go:174-175` 摘 INFO 支删行人话（WARN 支同句＝存活锚，计数 1 未变） | `:111:` `the INFO branch must state the delete-shape even when the values hide it; log:` | **1 枚**＝只红本枚（`FAIL_TOP=1`） |
| P3 | `writeguard.go:141` `base = deepCopyConfig(disk)` ⇒ `(m.cur)` | `:162:` `the resurrected section carried memory's value instead of the schema default:`（另 `:174`／`:180`） | 12 枚（守卫被摘 ⇒ 226/227 既有＋本腿 4 枚一起红） |
| P5 | `loader.go:257` `cp.SchemaVersion = SchemaVersionCurrent` ⇒ `= 1` | `:231:` `the write branch did not re-stamp the deleted schema_version line at the current version:` | 2 枚（＋既有 `TestTicket257R1AC1CleanMachineRefusesAllSevenFields`） |
| P6 | `loader.go:213` `p.BaseURL = pre.BaseURL` ⇒ `= ""` | `:320:` `the deleted base_url line came back as something other than the built-in preset:`（另 `:330`） | 2 枚（＋既有 `TestLoadFilePresetInheritance`） |
| P6 正控 | 同 P6 落点 | 无红句＝**该发它必须仍绿**（实测：p6 红名册里没有本正控） | — |

逐字件＝`logs/05-mu-p2.txt`／`05-mu-p3.txt`／`05-mu-p5.txt`／`05-mu-p6.txt`（＋同前缀 `-full.txt` 原始输出）。
每发都带：落刀前锚命中＝1（＝期望）、落刀后＝0、新锚≥1、存活锚计数不变、**当场还原哈希＝pristine**
（`writeguard.go efe12661e71b940f…d47bd`／`loader.go 40cb3e12b9bb425…7bdb`，件内逐字 `restored==pristine ? YES`）。

## (ii) P2 与 P1 在写侧到底走不走同一支（⛔ 读码推断当读数）

**同一枚函数，两条不同的支**，用两发突变把"同不同"量出来（⛔ 靠读码断言）：

- 同一函数＝`writeguard.go:111 func (m *Manager) mergeWrite`；`AddAllowedDir` 的唯一写路径＝
  `allowdirs.go:142`（`permmode.go:79` 是另一枚 caller）。两枚用例都从 `m.AddAllowedDir(testDir)` 进去，
  这一点是**同一支**（尺＝`logs/15-carrier-table.txt` 末段的 `grep -n` 行号，锚⇄HEAD blob 同一）。
- **差别在出口**（P1 走分叉侧 Warn 支 `:184-190`；P2 走 `len(diverged)==0` 的 Info 支 `:171-176`），
  机读证据＝两发交叉：
  - **p2 那发**（只摘 INFO 支那句）⇒ 红名册 `FAIL_TOP=1`＝**只红 P2 一枚**，227 的 P1 枚
    （`TestHandDeletedScalarLineComesBackAtSchemaDefault`）**仍绿**＝P1 读的⛔ 是那句（它读 WARN 支）。
  - **p3 那发**（摘"合并基准＝文件"）⇒ P1 红、P3 红、P6 红、226 的 AC#1/AC#2/AC#3 红，
    而**P2 一枚⛔ 在红名册里**（尺＝`logs/15` 末段 `grep -c '^--- FAIL: TestHandDeletedDefaultValuedLine…'` ＝ **0**）。
    原因也是量出来的：P2 的内存值与文件值同为一枚默认值 ⇒ 基准取谁都得到同一份合并结果，
    这枚守卫对它**不可见**＝正是"看不出差别"那一判据。
- ⇒ 两形⛔ 同一支持此结论＝上面两发红名册的**交集为空**，⛔ 本腿嘴上说。

## (iii) 有没有任何一格判"缺世界"

**零枚**。四格都在今天的产码下造得出可断的世界，五枚用例在未突变树上逐枚 PASS
（`logs/05-color-{Adflt,Bdtemp}-deliv.txt` 逐名）。
具名写一处**形状侧**的实话（⛔ 缺世界，属判据内容）：P2 那一形在**报告侧没有任何可指认的差别**
（文件⇄内存同值、走 Info 支、⛔ divergence 清单），所以本格能钉的只有两件事——
①那一行确实回到默认值（`readText` 在场），②那句文案仍必须把"删的行回到默认"说清（`:111`），
外加反向断言"⛔ 冒出 kept 清单"（`:101`，尺见 `logs/17-banned-char-ruler.txt`）。⛔ 因为它难钉就把 P2 并进 P1 交差（AC#2 明文），
也⛔ 用假绿填空。

## (iv) 新造的五枚里，哪一枚在未修码（起手锚 `98dbf20a`）上会红

**没有一枚会红**——本票是补仪器、⛔ 修行为，守卫今天是对的，所以正确答案就是这句。
盘上凭据（⛔ 糊过去）：

- 交件树＝基线树字节＋本腿那一枚新件；三枚产码 blob 在锚⇄HEAD 同一性 `same=YES`
  （`logs/12-roster-strict.txt`），交件四数＝`RUN=195 PASS_TOP=114 FAIL_TOP=0 SKIP=0`（两把临时根同值），
  基线＝`190/109/0/0` ⇒ 新增的 5 枚**逐名 PASS**（`logs/15` 逐名列出，两面各 5 行）。
- ⇒ "这枚钉真有牙"的凭据**只剩去牙那一发**＝四发逐枚点名红（(i) 那张表），每发另有
  锚命中数⇄存活锚⇄还原哈希三把尺守着，防"空操作伪装成成功复跑"。
- 反证一把：⛔ 恒红也排除了——同一批用例若 setup 坏就会红；实测第一发交件确实红过 2 枚，
  红句逐字在 `logs/06`／`logs/07`，那是本腿引号形状写坏（本仓 canonical 输出单引号字面串），
  改正后为绿＝证这批用例**会红、也⛔ 恒红**，两向都有读数。

## (v) `git diff --numstat` 逐笔：产码 `.go` 新增/删除＝0/0

| 我的提交（严格主题前缀尺 `^票 341 落地腿 341-r1·`） | numstat 全部行 | 产码 `.go` 路径数 |
|---|---|---|
| `6948eb08` | `389  0  internal/config/writeguard_unnailed_341_test.go` | **0** |
| `fd82f786` | 33 行，全部 `.scratch/wisp/probes/341/r1/logs/*.txt` | **0** |
| 本笔所在 | report＋report-answers＋`logs/16`/`17`＋台件脚本 | **0** |
| 票尾 bullet 那笔 | `.scratch/wisp/issues/341-*.md`（只追加 1 条） | **0** |

- 注释也没动：产码文件在锚⇄HEAD 之间 **blob 同一性 same=YES**（同一性比"0/0"更强——一字未动，含注释）。
  ⇒ "只许注释动，注释动了也要逐枚点名"这半句＝本腿**零枚注释改动**，无需点名。
- 越界尺（`logs/12`）：本腿全部路径 ⊆ {`internal/config/**_test.go`, `.scratch/wisp/probes/341/r1/**`,
  票 341 票尾一条} ＝ 之外 **0** 枚。⚠ 第一发这把尺用子串吃了编排者的派单笔（越界读数 16＝尺缺陷），
  改严格前缀后＝0，两处都留盘具名。
- 票面框尺＝7 未勾／0 已勾（起手⇄交件逐字同）＝⛔ 翻框。
