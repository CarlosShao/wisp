# 票 294 × 验收腿 `294-v1` —— 非实现者裁决表（`AC#1`..`AC#5` 逐格；翻框归编排者）

- 裁决腿＝`294-v1`（**非实现者**：⛔ 写过票 294 任何产码/测试件）；起手 HEAD＝`3d937469`（`Tue Oct 13 21:20:37 2026 +0800`，原件逐字在 `00-head.txt`）；运行期 HEAD 漂到 `8cfd7b1b`（别腿/编排者提交在飞；`cmd/wisp` 全程零动，收手复量）。
- **本表⛔ 翻任何框、⛔ 动 `Status:`、⛔ push。**
- 仪器两枚（代提 `e480e15b` 起零改动）：`cmd/wisp/resident_mute_294_test.go`=`7b20b7f8…cfde49`、`…_windows_test.go`=`ede84f7a…6969229`；blob 锚 `e480e15b:<件>`==`HEAD:<件>`==`9704ff08…`。
- 读数件＝`.scratch/wisp/probes/294/v1/logs/{00-head,10-ac1-astpin-v1,20-ac2-window-v1,30-ac3-executed-v1,40-ac4-ac5-v1}.txt`（每枚内嵌 tee 原件；变异一律 `-overlay` 或仓外备份＋`sha256` 还原证）。
- ⚠ 落点更正：起手原件曾落仓库根 `probes/294/v1/logs/`，该目录（含同波 229-v1 的同形目录）运行期间被人从盘上清掉（⛔ 本腿删、⛔ 本腿跑过 clean/checkout 族）；全部读数另留仓外档，本表与 5 枚件按其逐字重落至 229-v1 同波落点。

大白话一句：三枚判据我都亲手掰坏过——把走查的 `SelectorExpr` 层去掉、把生产实参/接收者换掉、把镜像常量改假、把 `executed` 的真相翻掉——每一发都红在它该红的行上，还原后全绿；AC#4/AC#5 是逐笔扫和逐名对，0 命中／逐名相同。

| 格 | 判语 | 本腿自跑的尺（一行；件名在尾） |
|---|---|---|
| **AC#1** 装配根形钉 | **成立（有牙；方向=尺侧 1 发＋产码侧 2 发）** | 尺侧：去掉 `SelectorExpr` 层（overlay）⇒ `-run TestAC294AssemblyRootAttachesTheMuteGateByShape` rc=1，红在 (a)「contains **0** rb.attachMuteGate call(s), want exactly 1」＝该层承重、`==1` 不是空集通过形。产码侧（仓外备份＋还原）：实参 `raudio.toggleMute`→`nil` ⇒ 红在 `:133`（现读 `""`）；接收者 `rb`→别名 ⇒ 红在 `:127`（现读 `"rbalias"`）。还原：`S0=S1=72fcee7a…`、`git status --porcelain -- cmd/wisp` 零行〔10-…〕 |
| **AC#2** 真窗 mute 注册 | **成立（有牙；非 skip）** | 绿读：`AC#2 reading: id=2 name="mute" binding="Ctrl+Alt+X" status="live" vk=0x58`（起手/收手两发，`--- SKIP` 零命中）。牙：镜像 `hkMute294` 2→3（overlay）⇒ 同源尺红在 `:82`「mirror drift: cmd/wisp hkMute294 = 3 but internal/ball hkMute = 2…」＋真窗用例红在 `:110`（`IsLive(3)`＝cancel 槽 idle 未绑；红句表里 slot2=mute 仍 live＝红在槽号不红在环境）［20-…〕 |
| **AC#3** `executed==true`＋M2 反形 | **成立（有牙；正反两向）** | E1：翻 `case gate.Muted():` 支 `, true`→`, false`（r1 翻的是 Open 支，互不重叠）⇒ rc=1 红在 `:152`「answered executed=false while the gate DID turn (muted=true open=false)」、句面一字未动。E3a：句面删字而 bool 不动 ⇒ 被验尺**绿**（＝没读句面）。句子读器对照：同字节上绿（被骗）、句变才红（E3b 红在变体 `:152`）⇒ 判据锚＝gate 自读＋`executed`，⛔ 句面〔30-…〕 |
| **AC#4** 越界 | **成立（0 命中）** | `e480e15b^..HEAD` 逐笔 `git show --name-only`（扫描时 HEAD=`8cfd7b1b`，17 笔）：十枚禁列逐笔 0；点名五笔 `e480e15b`=44 件／`ca777e8f`=58／`ff384b3d`=1／`b3959665`=1／`95708485`=1 全 0（子串形复扫同 0）；票面 `ca777e8f` numstat=2/0、后三笔票面零动；仪器 blob 未变〔40-…〕 |
| **AC#5** 门禁四数 | **成立（附一处降格＋一处欠账，见下）** | 交集：核 r1b 原件（stable 122=122 `cmp` rc=0、`diff` 0 行；target 13⇄17 唯一差＝4 枚 `TestAC294*`）＋本腿自抽 `-run TestAC29` rc=0／17 PASS／名册与 r1b `stable-target-post` **diff 空**；`GOFLAGS= go build ./...` rc=0；`gofumpt -l cmd/wisp` 仅两枚既有脏（⛔ 本票账）；`sh scripts/d22scan.sh` worktree rc=0 clean。⚠ 降格：prescribed 全包 5/5 死（157.188–174.374s；死点 AC14×1／AC13×2／PanelThread×1＋companion ProvidersProbe×1），`resident_*` 区零覆盖⇒定向补尺代偿（⛔ 互替）〔40-…〕 |

## 哪里没牙／哪里没跑（具名，⛔ 掩盖）

1. **AC#1 子句 (c)（时序）没有独立红形**：可编译的"倒序"变异不存在（`raudio` 必先绑定才可被引用，Go 作用域使然）⇒ 该子句的牙属结构性；若编排者要逐子句红形须另派（本腿⛔ 拿编译不过的假红充数）。
2. **AC#2 的 SKIP 退路**：headless 机上的 skip 不算读数（票面口径）；本机有桌面、两发绿读均非 skip——其余机器形本腿未覆盖。
3. **AC#3**：本腿与 r1 各翻一支（Muted／Open）合起来两向齐；第三段 `muteGesture` 只判"到没到门"（件内自陈，⛔ 句面）、未另造反形。
4. **AC#5 prescribed 全包覆盖未取得**（`resident_*` 区在票 314 死区之后）：本腿⛔ 撞整包死区；定向补尺只覆盖 `TestAC29*` 17 枚。完整形＝314 定/修后同形 2+2 重跑（编排者队列）。
5. 本腿未跑：整包套件（⛔ 故意）；`.scratch` 外仓根 `probes/` 目录已被外部清掉（读数有仓外档；重落件在 `.scratch/wisp/probes/294/v1/`）。
6. 本腿踩坑具名（防下一位）：`go test` 重复 `-overlay` **只取最后一枚**——首轮 E2/E3b 因此把产码侧变异整枚覆盖、读数作废；已用合并 JSON 重取（作废原件留在 30 件里并标注）。

## 翻框建议（翻框归编排者；本腿零翻）

- **建议 `AC#1`..`AC#5` 全翻**（`AC#0` 已翻）。附三条件：① `AC#5` 的翻随收档带上第 4 条欠账（314 定/修后重跑同形 2+2）；② `AC#1` 的翻知晓子句 (c) 无独立红形（结构性牙）；③ 三红（`TestTicket255RosterStillMatchesTheActualReadSites`＝台账 A878/313 账、票 314 死区、真窗族漂色）⛔ 算 294 的账、⛔ 本腿修。
