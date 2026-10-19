# 346-a1 件03 = 另两味成因的盘上取证（临时根拼法族 ⇄ CI 摆动那枚 ⇄ 争用）

尺＝件00 的归一日志＋顶格锚定；⛔ 本腿跑任何包。★⛔ 与件04 重复（blob 尺在那儿）。

## 一、味乙＝临时根**拼法**造的红族（9 枚 junction ＋ SID 那族）
1. 编排者那一发（长 D 卷 `GOTMPDIR`）＝`D:\tmp\wisp\20261019\orch227-b\full.log`（1544211B，四数 2317/1529/29/8）：
   - L7452（行长 309）逐字＝`bridge_junction_windows_test.go:268: 前置条件缺失：造不出真 NTFS junction
     （mklink /J D:\tmp\wisp\20261019\orch227-b\gotmp\TestBridgeRefusesARealJunctionOnTheReadRoute387489775\001\jn
     D:/tmp/wisp/20261019/orch227-b/gotmp\…\002: exit status 1: Invalid switch - "tmp".`
     ⇒ **link 侧＝反斜杠、target 侧＝正斜杠**（`D:/tmp/…`）＝cmd 把 `D:/tmp` 当开关吃掉。
   - 同族逐名行区间＝L7459／L7464／L7469／L7474／L7479／L7498（`bridge_junction_windows_test.go:394`／`:592`）、
     L7698（`fs_staging_windows_test.go:229`）、L8389（`refusal_20r2_windows_test.go:34`）
     ＋ `internal/risk` 两枚 L6068／L6071（`syncdirs_ancestor_reparse_ticket105_windows_test.go:59`／`:116`，
     同 `Invalid switch - "tmp"`）＝**9 枚**逐名对上编排者名册（量的）。
   - ★那发红册 29 枚里**两枚仪器 0 命中**＝它们⛔ 参与 Windows 编译（件04 第二节 tag 尺）
     ⇒ **临时根那味对本票两枚仪器⛔ 是成因**（⛔ 把它们归进这一族＝本腿最重要的一条排除）。
     ⚠ 该发射程＝本腿现量：最后一枚包结果行是 L10587 `FAIL … internal/winsec 37.588s`、其后⛔ 包行
     ⇒ **⛔ 整包终态名册**（`cmd/wisp` 那批⛔ 轮到）＝与编排者 `0ca0966a` 那句"未跑到 179 枚"独立同向；
     本排除⛔ 依赖"没跑到"那一维（`internal/winsec`⇄`internal/tools`⇄`internal/risk` 三包在该发**都跑到了**、
     两枚名在其包段内 0 命中＝有效读数）。
2. 裁腿 `227-v2` 成对两发（已入库件 `probes/227/v2/logs/03-unqualified-and-gates.md` 8159B／44 行，第 5-10 行）＝
   甲发（反斜杠短根 `C:\Users\…\227v2-pf-bs`）`2/2/0/0 rc=0` 两枚全绿 ⇄
   乙发（正斜杠短根 `C:/Users/…/227v2-pf-sl`）`2/0/2/0 rc=1` 两枚全红，
   乙发红句逐字＝`…mklink /J C:\…\jn C:/Users/…: exit status 1: Invalid switch - "Users".`（该行第 21 行）
   ⇒ **根因＝拼法，⛔ 长度、⛔ 卷**；该件 Ⓒ 段（第 26 行）已具名更正 `A1053`⑵ 那句"长临时根让 mklink 建不出来"。
   ★本腿独立复尺＝上面那 9 枚长根红句里 target 全为 `D:/tmp/…` 正斜杠形（逐枚量的）＝与甲⇄乙两发同形 ⇒ **认该更正**。
3. 同形早现（⛔ 本波新账）＝`probes/236/v2/raw/baseline-tools-v1-forward-slash-gotmpdir.txt:33,40,45,50,55,60`
   ＝`C:\tmp\wisp236v2gotmp\…\jn C:/tmp/wisp236v2gotmp\…` 混拼下同一句"前置条件缺失"（票 236 那波已现过一次）。
4. 造 junction 的码（读码，⛔ 运行）＝`internal/tools/bridge_junction_windows_test.go:50-58`
   `exec.Command("cmd", "/c", "mklink", "/J", link, target).CombinedOutput()` ⇒ 失败即 `t.Fatalf`；
   ★该文件 `:33-36` 自己写着 "If the machine cannot produce one, the case fails LOUDLY …; it never skips"
   ＝**弃权即红**的第三枚同族仪器（本票 AC#0 那把"全仓还有几枚"的尺的一格，逐名见结论件 (ii) 末）。
5. 同一发里的另两味（⛔ 混进拼法族）＝SID 两枚 L9249／L9256＝**卷形在册**（票 311，⛔ 本票新账）；
   `TestCrossVolumeMove…` L7779＝红句逐字 `open C:\.wisp-tmp-…: Access is denied`＝C 卷根（`A1053`⑷ 在册）。

## 二、味丙＝调度争用（★盘上只钉得住一条，其余写未知）
- 三发 `test-core`（39101060753 `a42b35af` ⇄ 39124122121 `a5080612` ⇄ 39244425666 `6f808a83`）＝
  **同一作业、同一批包**，两枚仪器**同 run 同色**（双绿⇄双红⇄双绿，件00 第四节速查＋件01/件02 表①）
  ⇒ 它们抢的是**同一件事**的两个读数，⛔ 两枚各自独立的运气（这句是本腿对票面"CI 随机红"散文的形状收窄）。
- ★"runner 当时负载⇄兄弟作业数"这一维＝**未知**（尺＝`gh run view --json jobs` 只给 startedAt/completedAt；
  本腿⛔ 任何 runner 侧 CPU／负载读数，⛔ 现跑）⇒ 调度争用这味**⛔ 可证**，本腿据实写未知、⛔ 圆场。

## 三、第三枚摆动（票面喂本腿的"新料②"，本腿只对回、⛔ 采信为成因）
`TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`（`cmd/wisp/config_reload_223_test.go`）：
| 台面 | 颜色 | 行区间（归一件） |
|---|---|---|
| CI 39071959151 `test-windows` | PASS (2.13s) | `old39071.norm.log` L3291／L3304 |
| CI 39101060753 `test-windows` | PASS (1.84s) | `old39101.norm.log` L6467／L6480 |
| CI 39124122121 `test-windows` | PASS (2.23s) | `old39124.norm.log` L318／L331 |
| **CI 39244425666 `test-windows`** | **FAIL (2.22s)** | `win.norm.log` L1829..L1844 ＝ **B321082..B324162** |
| 本机整包（编排者那一发） | PASS (2.20s) | `orch227-b/full.log` L337／L352（量的＝读盘上件，⛔ 本腿复跑） |
| CI `test-core`（Linux） | ⛔ 现身＝⛔ 射程（0 命中，正控＝同件 7 枚 FAIL） | — |
- 红句逐字（段内、`--- FAIL` 前一行）＝L1841 `config_reload_223_test.go:545: the card does not name
  risk.permission_mode: wisp run: 答复监听已接入（卡片上给了编号）。…`；★同段 L1843 仍打出
  `config: locked loosening approved via L2 re-confirmation (D36 rule 1) section=risk keys=[risk.permission_mode]`
  ＝**那一发批准发生了、卡片文案没点到那枚 key** ⇒ 形状＝1s tick／等待预算的到达竞态
  （读码＝`config_reload_223_test.go:31` `reloadCaseBudget`、`:98`/`:108` `time.Sleep(20 * time.Millisecond)`）。
- 字节同一性（本腿现跑 `git rev-parse`）＝该测试文件 `a5080612` ⇄ `6f808a83` **同串 `f1cac524…`**；
  `git log --oneline a5080612..6f808a83 -- cmd/wisp` ＝ **0 枚** ⇒ 同码两绿一红＝时序／环境敏感
  （与 `A1054`② 独立同向，⛔ 抄它）。
- ★与票 346 那两枚**⛔ 同判据族**：两枚靠"抢到并发窗口"作数（到达数＝`refused-at-act`/`sealed-own`/`flips`），
  这枚靠"1s tick 内到达"作数；两枚有弃权红句、这枚红句是**真断言失败**（⛔ "零读数⛔ 是证据"那句）。
- 同包另一枚 `TestTicket223PermissionDeniedSitsInItsOwnSentence` ＝**四发 CI 全红**
  （`old39071` L3387／`old39101` L6563／`old39124` L414／`win.norm` L1927）＝**稳定红在册**
  （票 346 现量表第 5 行已逐名归位）⇒ ⛔ 与摆动那枚混记、⛔ 新账。
