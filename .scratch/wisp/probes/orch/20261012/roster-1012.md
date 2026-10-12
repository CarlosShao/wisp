# `roster-1012.md` — 编排者 13:2x 现量：59 枚残渣的构成＋专用根里那 12 枚"进不去"的目录

台面＝工作树 HEAD `8e3cb370`（在飞两枚腿：`310-r2` 落地／`311-a1` 只读）。尺件＝`logs/orch-roster-temp-1012.txt`（`ps_rc=0`，stderr **0 字节**，正控同行）＋`logs/orch-sealed-dirs-1012.txt`（`roster_rc=0`）。
★本节目光所及＝**只数、只名册、⛔ 删任何一枚**（票 310 那条铁规矩；机主没回「清」之前⛔ 动）。

## 1. 顶层名册（一把尺：`Get-ChildItem $env:TEMP -Directory | ? Name -like 'Test*'`）

- **枚数＝59**（与 13:0x 那发逐字同数＝我这批编排者的活⛔ 新增残渣；正控同行 `temp_all_dirs_positive_control=`**2834**，⛔ 拿 59 当"整个临时目录只有 59 个"）。
- 按前缀分桶（`sed` 剥掉数字后缀后 `uniq -c`，尺＝本件 §1 那把名册）：**mockllm 一族 25 枚**（`TestMockllmGoldenByteIdenticalEvents` 17 ＋ `TestMockllmSharedFaultSuite` 8）／`TestRoundTripLoadMarshalLoad` **9 枚**（票 310 本族）／`TestRollingWriterRetentionSweep` 4／`TestGateJudgesThePrivateSetByResolvedSID` 3（票 311 那枚红用例的族）／`TestFailoverAfterPrimaryExhaustion` 3／其余 15 个前缀各 1 枚。
- ⚠ 归因口径提醒（正典已写）：**桶名只到"哪个测试留下的"，⛔ 到"哪一发"**；逐枚 `CreationTime` 在件里，要做"哪一程"的归属还得配 `comm -13` 作差那把尺。
- ⚠⚠ **同一格我第一次读到 58 枚＝尺的假数**（★**我第一版把成因写成"`grep` 打印 Binary file ⇒ 少一枚"，那句⛔ 成立**；下面是重跑后的现量）。真形＝件里有**一枚带中文的目录名**：`TestTicket223R2FailureSentenceRouting…_\xbb\xb5\xb1\xed\xcd\xb7_A2_3165385141`，PowerShell 经 bash 重定向落盘＝**GBK 字节流**（件内 NUL＝0，但 `utf-8` 解码在 offset 1192 失败＝`invalid continuation byte 0xc9`）⇒ UTF-8 locale 下 `grep` ①对整件打印过 `Binary file … matches` ②**严格正则那一把整行漏掉该枚＝58**。三把尺并排：PS `.Count`＝**59** ⇄ `grep -a -E '^Test'`＝**59** ⇄ `LC_ALL=C grep -cE '^Test[^|]*\|2026'`＝**59**。⇒ 定式：**名册件含非 ASCII 时，数行的尺一律 `LC_ALL=C`（或 `grep -a`），⛔ 拿一把 locale 敏感的尺的"少一枚"当事实**；★而且**⛔ 用"哪一把尺报了少一枚"去倒推成因**——我这一次就是从成因写歪了一步（同族错＝`A858` §4 那句根据）。

## 2. 两枚专用根（票 310 `AC#1` 的落点）的条目数

| 根（一把尺＝`Get-ChildItem -Directory`，⛔ 与 `ls -1` 混用） | 运行目录枚数 | 根内 `Test*` 残渣条目 | 根里的普通文件 |
|---|---|---|---|
| `C:/Users/swq/AppData/Local/wisp-test-gotmp`（出厂默认） | **11** | **2** | 1 枚（`portable-20261012-115401-2013/TestRoundTripLoadMarshalLoad1761995879/001/config.toml`＝某次测试写进自己临时根里的文件，⛔ 当残渣、⛔ 当泄漏） |
| `D:/tmp/wisp-test-gotmp`（`310-r1` 换卷实验那一族，票 311 的台面） | **6** | **10** | ★**2 枚读数件落错了地方**：`shapec-115000-c2_redirected.run.txt`（**0 字节**）＋同名 `.err`（35458 字节）——正典"读数件⛔ 写进任何专用根里"就是被这一形教出来的；件⛔ 删（票 310 铁规矩＋⛔ 动别家的件），只登记 |

⇒ **专用根⛔ 是"收干净"的意思**：它把残渣从 `%TEMP%` 顶层搬进一枚具名、可整族数、可整族交别人复核的地方，收尾那一半仍⛔ 发生（票 310 `AC#1` 框里那句成立范围已写明）。

## 3. ★新现量（给 `311-a1` 的 `AC#0` ③ 递料，⛔ 我据此定罪）

`find` 递归那两枚 `D:` 根的运行目录 ⇒ **`perm_denied_dirs_unreachable=`12**（尺＝`find /d/tmp/wisp-test-gotmp -type f` 的 `Permission denied` 行数；件 `logs/orch-sealed-dirs-1012.txt`；正控同行＝同把尺对一枚⛔ 被封的运行目录 `winsec-20261012-122026-451` 递归得到 **1** 枚目录、**0** 行被拒）。

- 被拒的那 12 枚全是个**具名形状**：`TestGateJudgesThePrivateSetByResolvedSID…/001/{la,named,numeric,runner}` 与 `TestSealNarrowsAndNamesThePrincipalItRemovedBySID…/002/003/la`。
- ⚠ **目录名本身还数得到**（`Get-ChildItem -Directory` 能列出＝上面那 10 枚），**descend 才被拒**；⛔ 拿"数得到"当"进得去"、⛔ 拿"进不去"当"不存在"。
- ⛔ 我说这是谁干的：**〔真机量的〕只有"创建者这一侧被拒"这一句**（我就是那枚创建者账号）；"是哪一枚 SID 被留在 ACE 里、`WRITE_DAC` 还在不在"＝**〔待 `311-a1` 量〕**，归票 311 `AC#0` ①②③三问，⛔ 由我这发读数代裁。

**next**＝等 `310-r2`／`311-a1` 交回，我按老规矩逐把自己要看的尺重跑才翻框；`310-r2` 核过 ⇒ 票 310 改 `-done`；`311-a1` 交回 ⇒ 裁票 311 `AC#1` 三形（甲／乙／丙，默认仍押在"读数没到手前⛔ 动 `internal/winsec`"）。⛔ push、⛔ 删目录。
