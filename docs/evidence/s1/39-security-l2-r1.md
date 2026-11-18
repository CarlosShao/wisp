# 票 39 AC#2 — 安全放宽 ⇒ 原生 L2 确认 ⇒ 接受生效并留审计／拒绝保持旧值

写腿 agent-128 · 2026-11-18 · 派单 `.scratch/wisp/probes/orch/20261118/txt/39-security-l2-dispatch.txt`
（只做 AC#2：`.scratch/wisp/issues/39-config-editor-gui.md:38-39`；e2e with 票 21 的卡）

## 0. 一句话

程序化配置写路（`WriteSection`）遇到具名"安全放宽键"时，先弹**一张**原生 L2 卡（复用票 21 的
`gate.PendingApproval` 卡构造，没有第二套卡、没有面板侧 allow）；接受 ⇒ 真写 + 一条审计事件
（老值→新值 + 卡 id）；拒绝/超时 ⇒ 一个字节不写（真读回旧值断言）。

## 1. 具名清单（18 行 = 13 枚静态键 + 5 个 plugins 动态族）

清单＝`internal/config/loosenroster_39.go` 的 `SecurityLoosenRoster`。**⛔ 自己发明分类**：每一行
的来源就是 D36 方向审计函数本身（`manager.go` 的 `riskDirection`/`fsDirection`/`netDirection`/
`pluginsDirection`——`lockedDirectionOf` 与 reload 的 `planLocked` 都派发给它们），`Source` 列逐键
给 file:line，`Pin` 列给该行必须还在的字面量（动态族 pin 的是审计函数自己的拼接片段，钉的是代码
不是拷贝）：

| 键 | 段 | 来源（现值） |
|---|---|---|
| risk.shell_enabled | risk | manager.go:467 |
| risk.allow_shell_string | risk | manager.go:472 |
| risk.l1_window_sec | risk | manager.go:477 |
| risk.permission_mode | risk | manager.go:491 |
| risk.shell_allowlist | risk | manager.go:495 |
| risk.blacklist_overrides | risk | manager.go:496 |
| fs.allowed_dirs | fs | manager.go:502 |
| fs.reparse_point_exceptions | fs | manager.go:503 |
| fs.delete_enabled | fs | manager.go:505 |
| net.allowlist | net | manager.go:515 |
| net.block_private_ranges | net | manager.go:517 |
| net.proxy.mode | net | manager.go:525 |
| plugins.tier2_enabled | plugins | manager.go:536 |
| plugins.\<id\>（新条目） | plugins | manager.go:542 |
| plugins.\<id\>.enabled | plugins | manager.go:547 |
| plugins.\<id\>.capabilities | plugins | manager.go:551 |
| plugins.\<id\>.net_allowlist | plugins | manager.go:552 |
| plugins.\<id\>.host_api | plugins | manager.go:553 |

匹配器 `IsSecurityLoosenKey`：静态键精确匹配 + plugins 四叶后缀族 + 新条目族；族模板拼写本身与
未知路径都答 false。测试四枚（`loosenroster_39_test.go`）：
①`Test39LoosenRosterSourcesReReadTheAuditors`——测试现读 `manager.go`，逐行核对 Pin（清单烂了
立刻红）；②`Test39LoosenRosterCoversEveryAuditorEmittedLoosening`——驱动四个审计函数做全放宽
翻转，每个被发出的放宽键必须过匹配器（审计学了新键而清单没学 ⇒ 红）；③
`Test39LoosenRosterKeysAreAuditorEmittedInLockedSections`——每枚静态键必须真能被审计函数以放宽
方向发出（发明行＝红）、每行段必须在 `TierRegistry` 标 locked；④`Test39IsSecurityLoosenKeyWidth`
——匹配器宽度双向钉 + 清单常量与匹配器不许漂。

## 2. 写门（三件）

1. **清单钉进写路**（`sectionwrite.go` 的 DirLoosen 分支头部）：方向审计判为放宽的键若有任何一枚
   不在清单上 ⇒ `outcome=refused-roster`，fail-closed 拒写（清单与审计漂移时宁可拒，不给没有
   名字的键发卡）。
2. **排队即登记**（`sectionwrite.go` + `loosenresolve_39.go`）：`WriteSection` 排队前先铸造
   `LoosenRequest.RequestID`（`host:config-loosen#N`）并把它注册进 Manager 的 pending 表，再交给
   `LoosenQueue`——卡一存在就可解析，不存在"答得比登记快"的窗口；队列拒收则注销。
3. **接受/拒绝的落点**（`Manager.ResolveQueuedLoosen(id, approved)`，新导出面）：
   - 接受：对**现盘**重算（卡可能挂了很久，原 base 已过期）→ Check → Apply → 单键所有权复核 →
     validate → 方向复裁。仍是放宽 ⇒ 新键必须被卡上点名的键覆盖（父子/精确；插件族父子拼写都算
     覆盖），否则 `outcome=refused-drift` 拒写（卡没批过的键一个都不写）；已是盘上现值 ⇒
     `outcome=unchanged`（不谎报 written）；文件反向走成收紧 ⇒ 按收紧落盘（收紧不需要卡，D36
     rule 1）。落盘走**同一条** `commitSectionWrite`（内存先行 + 受守合并写 + 失败回滚）。成功后
     **恰一条**审计行：`SECTION-WRITE … outcome=written card=<id> changes=<key:old->new,…>`。
   - 拒绝/超时（gate 把无回答的卡自身判成 reject，两条臂合一）：`outcome=denied card=<id>`，
     零写入，旧值原样。
   - 一卡一写：pending 表在 resolve 时即取走（one-shot），同 id 二次 resolve 按名拒绝；未知 id 拒绝。

## 3. 装配根接线（cmd/wisp，票 195 证据 §7 点名的欠账）

`cmd/wisp/config_loosen_gate_39.go`：
- `configLoosenQueue` 实现 `config.LoosenQueue`：**只入队不等人**（裁定 N9），卡等待协程经
  `observe.Default.Spawn("approval-waiter", <RequestID>, loosenRoot, …)` 起（D38b 名册名，D22 禁令
  #1 的唯一合法协程形状），`close()` 时 cancel + join（run.go 的 close 增 loosenRoot/loosenHandle
  两枚字段与收口，废弃卡按 deny 落地——关机永远不是放宽溜进去的时刻）。
- `resolveLoosenCard`：`gate.AdmitTextTask(configLoosenTaskID)`（D47 合成任务，退出即撤销）→
  `gate.PendingApproval`（risk.L2，`Tool=config.loosen`，`CorrelationID=RequestID`——卡 id、
  pending 键、审计 join 三者同一值）→ allow ⇒ `ResolveQueuedLoosen(id,true)`；其余一切出口（拒绝、
  超时、gate 缺席）⇒ `ResolveQueuedLoosen(id,false)`。
- `startConfigLoosenGate()` 在 `assembleRuntime` 的答复面之后调用（与 `startConfigReload` 同序
  规则）；gate 缺席的宿主不接线，保持 fail-closed 拒绝。
- **面板侧不动**：面板来源放宽仍被 `ErrPanelLoosenRefused` 结构性拒绝（R20/票 195 冻结形状）；
  本腿没有给面板加任何 allow 面（禁区②/禁令 #6 原样）。
- 无新 C17 方法（C17=PanelBridge；本腿全部改动在 config.Manager / cmd/wisp 装配层，理由：卡与
  pending 表的生命周期都归 Manager，装配根只做适配——与票 195 §7 的预同名）。

## 4. e2e（载具实跑，票 21 的卡）

`cmd/wisp/config_loosen_gate_39_test.go` 四枚（真 gate + 真 native Allow/Reject + 真 config.toml，
断言全部读盘回）：
- `Test39NativeLooseningCardAcceptedAppliesAndAudits`：排队不写 → 卡到 native 面（id=卡命名空间、
  tool=config.loosen、L2）→ `Gate.Native().Allow`（grant 单次）→ 盘上 zed.enabled=true；
  恰一条 `SECTION-WRITE outcome=written card=<id> changes=plugins.zed.enabled:false->true`；
  `LOOSEN-GATE state=answered … result=allow` 与卡同 id join。
- `Test39NativeLooseningCardRejectedKeepsOldValue`：native Reject ⇒ 读盘回 zed.enabled 仍 false；
  恰一条 denied 行点名卡 id。
- `Test39NativeLooseningCardTimeoutKeepsOldValue`：无人应答，gate 自己的 C18 `ApprovalTimeout`
  （300ms 测试值）到期 ⇒ 读盘回旧值；恰一条 denied 行。
- `Test39LooseningWithoutGateFailsClosed`：不接线的运行时 ⇒ `ErrLoosenQueueUnwired`，盘零变化。

## 5. 三枚牙（变异只做仓外副本：`%TEMP%\wisp39-mut`，验完已删）

| 牙 | 变异（仓外副本内） | 载具读数 |
|---|---|---|
| ① 摘 L2 门（直接写） | `sectionwrite.go`：`if dir == DirLoosen` → `… && false`，放宽直落盘 | `Test39NativeLooseningCardAcceptedAppliesAndAudits` **FAIL**（无卡即写入，"queueing must return no receipt" 红） |
| ② 接受分支不写审计 | `loosenresolve_39.go`：删掉接受落盘后的 `m.auditf(...)` 行 | 同一枚 **FAIL**（"accepted SECTION-WRITE lines = 0" 红） |
| ③ 拒绝分支仍写入 | `loosenresolve_39.go`：`if !approved {` → `if false && !approved {` | `Test39NativeLooseningCardRejectedKeepsOldValue` **FAIL**（读盘回 zed.enabled=true 红，日志留 `outcome=written card=…`） |

原树（无变异）下同一批测试全绿。

## 6. 门禁读数（本机 windows，2026-11-18）

- `sh scripts/d22scan.sh`：`d22scan.sh: OK - control rc=0, scan rc=0`（末行），0 findings。
- 载具：`scripts/wisp-cli-staged-tests.sh`（A1888 串行纪律，跑前 `Get-Process wisp.test` 空）。
- `go build ./...` rc=0；`go vet ./internal/config ./cmd/wisp` rc=0。
- `go test ./internal/config/ -count=1` 全绿（含票 195/215 存量牙，无一放宽）。
- ⛔ 未 push；thresholds/golden/allowlist/PLAN/specs/D43/design 零触碰；面板侧 allow 零新增。

## 7. 边界与诚实声明

- **接受路径的"宽松面"仍是收紧的**：面板来源的放宽被结构性拒绝（本腿未动）；原生来源要卡；卡上
  没点名的键在 resolve 时拒写（drift fail-closed）。
- 本腿交付的是 **Go 侧门**。AC#2 的勾格仍归对抗验收（裁决者≠实现者）；面板页把锁行画出来并接到
  `config.set` 的半格在票 39 的前端腿射程（`internal/panel` 的 `lockedFieldFamilies` 拒绝文案今
  天已把"走原生二次确认那条路"指到本腿接的这条路上）。
- `SetPermissionMode`/`AddAllowedDir`/`SetAllowedDirs`（票 90/201 的自有卡路径）不经过
  `WriteSection`，本腿不碰——它们各自的头注已声明"caller 已弹卡"，混并会双卡。
