# 351-a1 件04 = 决定性窗口的逐字入库（把仓外盘上件的关键段搬进仓，让交件可在盘上复查）

来源＝`346-a1` 已入库归一件所对应的仓外件（⛔ 本腿复跑任何包、⛔ 重造那把尺）：
`D:\tmp\wisp\20261019\346a1\win.norm.log`（＝run `39244425666`／job `118032776659` test-windows，归一件 683457B）⇄
`D:\tmp\wisp\20261019\346a1\old39101.norm.log`（＝run `39101060753` 全量，1132935B）。
★四数本腿复尺与 346-a1 留档逐数同＝件00 第二节；本件只搬**窗口**、⛔ 搬整件（帽 160 行）。

## 一、`win.norm.log` L1829..1844＝本票那枚红的完整段（★go test 正文印在 `--- FAIL` 之前）

```
1829: === RUN   TestTicket223ModeLooseningChangesTheRunningModeAfterAllow
1841:     config_reload_223_test.go:545: the card does not name risk.permission_mode: wisp run: 答复监听已接入（卡片上给了编号）。原生侧：yes <编号> / no <编号> [理由] / veto <编号>；面板路线：panel-no <编号> [理由]、panel-yes <编号>（这一条只会被服务端 API 拒绝）、head / view <编号> 看卡片；quit 退出监听。
1842:         wisp run: 配置热加载已接管（每 1s 检查一次 config.toml）。手改会按 D36 三档处理：可热加载段立即生效；[risk]/[fs]/[net]/[plugins] 的放宽要先答一张 L2 卡，不答按拒绝保留旧值；重启档的改动本次不生效，会另有一句告诉你为什么不生效。
1843: time=2026-10-19T11:29:30.864Z level=WARN msg="config: locked loosening approved via L2 re-confirmation (D36 rule 1)" section=risk keys=[risk.permission_mode]
1844: --- FAIL: TestTicket223ModeLooseningChangesTheRunningModeAfterAllow (2.22s)
```
（L1830..1840＝sink 安装／mockllm 起听／schema 三步迁移，逐字在仓外件；本件省略那 11 行并具名省略。）

★读法＝红句 dump 里只有 L1841＋L1842 两条**操作员横幅**、卡片那一行（`[确认 L2 config.reload] …`）**⛔ 在场**；
而 L1843＝同一次调用的 keys 读数**带上了那枚 key**。⇒ 件02 的定性（仪器形状、⛔ 真缺陷）就在这四行上。

## 二、`old39101.norm.log` L6433..6438＝同一形状的第二枚到达（★成对响，决定性）

★逐字性声明（免得本件被当全文抄）＝上面两块里带 `…` 的行（L6434／L6436、以及第一节说明的 L1830..1840）＝
**具名省略**；⛔ 省略号的两行（L1841／L1843、L6433／L6435／L6437／L6438）＝**逐字**，
尺＝仓外件行号（`D:\tmp\wisp\20261019\346a1\`），要全文按件00 的归一尺在那两件上取那几行。

```
6433:     config_reload_223_test.go:403: the console did not render the reload card: wisp run: 答复监听已接入（卡片上给了编号）。原生侧：yes <编号> / no <编号> [理由] / veto <编号>；面板路线：panel-no <编号> [理由]、panel-yes <编号>（这一条只会被服务端 API 拒绝）、head / view <编号> 看卡片；quit 退出监听。
6434:         wisp run: 配置热加载已接管（每 1s 检查一次 config.toml）。手改会按 D36 三档处理：…
6435:     config_reload_223_test.go:406: the card does not name the key it is asking about: wisp run: 答复监听已接入（卡片上给了编号）。原生侧：yes <编号> / no <编号> [理由] / veto <编号>；面板路线：panel-no <编号> [理由]、panel-yes <编号>（这一条只会被服务端 API 拒绝）、head / view <编号> 看卡片；quit 退出监听。
6436:         wisp run: 配置热加载已接管（每 1s 检查一次 config.toml）。…
6437: time=2026-10-18T05:38:19.935Z level=WARN msg="config: locked loosening approved via L2 re-confirmation (D36 rule 1)" section=fs keys=[fs.allowed_dirs]
6438: --- FAIL: TestTicket223HandEditedFsLooseningCostsAnL2Card (1.84s)
```
⇒ ★"没渲染"(`:403`) 与"没点名"(`:406`) **一起**响＝同一次单读、卡片整段缺位。
这条是本票最要紧的一味：**"文案含 key"那格⛔ 曾单独响过**＝盘上从⛔ 出现"卡片在场但没点名"那一形。
同发里本票那枚＝`old39101` L6480 `--- PASS: …(1.84s)`＝绿。

## 三、五台面的该枚结果行（尺＝归一件顶格 `^--- (PASS|FAIL): <名>`，逐枚带行号）

| 件 | 该枚 | 结果行号⇄耗时 |
|---|---|---|
| `win.norm` | 本票 | **FAIL** L1844 (2.22s) |
| `old39101` | 本票 | PASS L6480 (1.84s)；★姊妹枚 `HandEdited…` FAIL L6438 |
| `old39071` | 本票 | PASS L3304 (2.13s) |
| `old39124` | 本票 | PASS L331 (2.23s) |
| 本机 `orch227-b/full.log` | 本票 | PASS L352 (2.20s) |
| `core.norm`（Linux） | 本票 | ⛔ 现身＝⛔ 射程（正控＝同件 `TestManager` 24／`cmd/wisp` 16） |

★同族其余七枚在该包的五台读数＝件01／件02 的表（⛔ 抄 `346-a1`，本腿自己顶格重数）；
`TestTicket223PermissionDeniedSitsInItsOwnSentence`＝四发 CI 全红（稳定红在册、⛔ 与摆动那枚混记）＝
尺＝L1927／L3387／L6563／L414，其红句逐字＝`config_reload_perm_223_windows_test.go:95: attempt 1: icacls denied nothing,
so this case cannot show the permission sentence`＝**弃权即红**那形，⛔ 本票射程、⛔ 混批。
