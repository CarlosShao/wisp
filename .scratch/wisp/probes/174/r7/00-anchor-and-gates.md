# 174-r7 —— 票 174 两格交件（AC#2c 复核＋三态｜AC#2d 第二半：判据钉住 err-text 转述＋形乙修）

- 腿：174-r7（实现者；⛔ 框未翻、⛔ `Status:`/票名未动、只 commit 未 push、零 push）。面：`internal/tools/task.go`（修）＋`internal/tools/task_output_pointer_notice_errtext_174r7_test.go`（新判据）＋`probes/174/r7/**`＋票 174 尾 1 bullet。
- 提交：产码 `67a7ab41`（task.go＋新判据，2 文件 +122/−1）｜交付＝本枚＋logs＋票尾 bullet（同一枚 commit）。
- 起手锚（20:28 首读）＝`e9810346`；**漂移逐报**：读数期间 HEAD 走到 `eb141e98`(20:33)→`ccab14fb`→`46dee9ef`→`51e2f5ea`→`b75089cc`→`f79c2d74`→`705c4655`(20:54，本腿 67a7ab41 的父)。逐枚 `git show --stat`＝编排者/133-r1/136 的 bookkeeping（probes/**＋issues/**＋台账＋证据表＋一枚 `scripts/slo-freshness.sh` 修复），`internal/tools/**`＋`internal/agent/**` 零命中。
- 窗口自证：门禁那一发（20:56:03）`tasklist` 对 go.exe/wisp.test.exe 皆"No tasks…"＝安静窗；133-r1 已交件（46dee9ef）。
- ⚠ 一处自报：产码 commit 正文里 `` `case err != nil` `` 那半句被 bash 反引号吃掉（提交信息现读为"pointerNotice  支剥掉"）——内容无损、句有缺口，按本仓"不改写历史"以本行更正，不改提交。

## 0. 结论先行

1. **AC#2c：盘上【已落地】**（`5b14afe5`，10-03 14:51，`merge-base --is-ancestor` 本程现验为真；台账 `A575` 逐字"2c＝174-r2 落地（三 commit＋编排者补记在案）**待非实现者验收翻勾**"）。派单前提"**今天仍然静默**"**与盘上不符**——本腿的活＝现读复核＋补齐两向三态读数，**未重做、未再动 `spill.go` 一字节**。剩下的动作是"非实现者翻勾"，不是再落产码。
2. **AC#2d：第二半落齐**。派单引的"全仓零覆盖/MUT-V5"是 09-27 口径（r2/r3/r6 之后，①"说了没有"与②的路径/根列表两半已有牙）；**真仍无牙的恰是②里 `err.Error()` 内文的那一半**。本次：新判据（全平台、零 skip、禁词独立推导）＋**形乙 1 行修**，两向牙＋三态齐（§3 表）。
3. 死线判定（派单"两支取舍"）：**不触发停手**——177-census 判"盖章×续读今天不互斥"，故"显式裁掉 spill 桩这一支"没有更正确的理由；落地形本就是 177-safe（判据#5：注意句零路径字符、"全文见"恰一次＝豁免正则纹丝不动）。

## 1. AC#2c —— 复核与三态（读数逐字在 logs/04、logs/05）

- **落点形状**（`internal/agent/spill.go`，本腿零字节）：`PointerJudge{Canonicalize; InAllowlist}` 消费侧窄接口（零新依赖边；判"在不在根里"走与 fs.read 同一枚 C26 实例）＋三支说话：未接线＝fail-closed 明说／`Canonicalize` 出错＝只说"连规范化都没通过"、**不拼 err.Error()**／根外＝说清"要一张 L2 卡"并给出两条回来的路；健康＝零说明。
- **未修不响**（log 05 读数10）：overlay 把 `pointerNotice` 种回恒返 `""`（＝5b14afe5 之前的静默形；字节级等价：该 `%s` 槽为空后即修前模板）⇒ 174 族 **7 枚红**、桩文本只剩"全文见 <路径>"；仍绿恰是两枚不该被静默影响的（安静正例／指针唯一）。
- **修后响**（log 04 读数9）：HEAD 上 174 族 **9 枚逐枚 PASS**；安静正例单跑 PASS（判据自带"零注意＋逐字节等于修前模板"）。
- **还原绿**（log 05 读数11）：撤 overlay ⇒ `ok 2.537s`，`internal/agent/**` porcelain 零行（md5 与 HEAD blob 同）。
- **具名没做（AC#2c 内）**：① "不存在／非普通文件"那两支**未补**——spill 腿的存在性由构造保证（文件是同一函数几行前刚写的，三条失败路都在拼句前 return），补 `os.Stat` 只会得到一枚恒真尺（spill.go:203-209 的既有论证＋判据#4/#5 钉着）；若编排者要字面三支对齐，那是一处**设计变更**、不是本格的缺牙。② "未修码真身"不可单独编（判据引 `WithPointerJudge`，与新码同 commit；静默形 overlay 是 r2 既定替代形）。③ 177-census 的两包取数本腿未复跑（属它腿读数，我只复核了它的结论与 `spill.go` 判据#5 的对应关系）。
- **具名环境事实**：生产 `cmd/**` 未接 `Config.PointerJudge`（本程 `grep -rn PointerJudge cmd/` 零命中）⇒ 真机 spill 桩今天走 fail-closed"未接线"句——**明说而非撒谎**，接线属 AC#2b 面（r8 在飞）。

## 2. AC#2d —— 泄漏两半、钉、形乙修、两向牙（读数逐字在 logs/01-03）

- **零覆盖尺**（log 01）：钉"连规范化都没通过"的 6 枚断言全在问"说了没有"；`grep -rn "err.Error()" internal/tools/*_test.go` 无一枚问本支拼了什么。
- **泄漏两半（未修码上"今天不响"）**：同一发跑既有两条泄漏扫 ⇒ **双 PASS**，而 notice 区逐字含上游文本——(a) junction 臂：`（risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions）`（含配置键名＝C26 内部状态）；(b) 便携臂：`（tools: empty path）`。
- **装钉**：`TestCanonicalizeErrorNoticeRelaysNoUpstreamErrorText174r7`（114 行；全平台、零 skip，故**不需要** `scripts/portable-tests.sh` 新 ledger 行）——① 正向：注意句在场、指针仍给；② 禁词三类同问：filed 字面量（含 `%q` 形）／根目录列表（含 fold 形）／**C26 自身错误文本＋`tools:`、`risk:` 包前缀**；禁词**独立推导**（同一构造器、同一输入现算，禁"抄一枚字面量"）。未修码上即红（log 01 读数2）。
- **形乙修**（1 行＋注释）：`task.go` 的 `case err != nil` 剥掉 `（"+err.Error()+"）`；目标句与 `internal/agent/spill.go:229` **逐字同形**。依据链：`A375`（该半格要钉必须动产码）→ `A587`（撤下 owner 清单：1 行/零判据要改/可逆/不需批准；目标串非新造）＋`174-a3 §4`（契约轴不挡：D15(3)/PLAN.md:2564 冻的是"头/尾/总长/路径"四样在场，`pointerNotice` 是塞进 `%s` 槽的附加句）。
- **两向牙**（log 03）：把根列表拼回句里 ⇒ 恰两枚根禁词红；把修前整枚种回 ⇒ 新判据重红；撤 overlay ⇒ 绿。既有判据**零枚要改**（r1 八枚、r2 两枚、r3 两枚、r6 一枚、r4 三枚修后全绿）。
- **具名没做（AC#2d 内）**：① 两枚我自造的首发缺陷照实留档（漏 import 一发；一枚过严的"恰一条注意"断言已摘）；② 本腿未改 `spill.go` 文案做"两腿对齐"以外的任何事（它本就已是目标形）。

## 3. 三态总表（两格）

| 格 | 未修 | 已修 | 还原 |
|---|---|---|---|
| AC#2c | log 05 读数10：静默形 ⇒ 族 7 红、桩零说明（不响） | log 04 读数9：HEAD ⇒ 族 9 绿、桩说话（响） | log 05 读数11：撤 overlay ⇒ 绿＋agent 面零字节 |
| AC#2d | log 01：泄漏两半逐字＋既有扫双绿（不响）；新判据红 | log 02：新判据绿＋两条既有扫 notice 变干净＋族 9 绿 | log 03：种回修前句重红→撤 overlay 绿＋porcelain 只含本腿两件 |

## 4. 门禁现量（20:56 安静窗；逐字＝logs/raw/raw-gates.txt）

- G1 `go test -count=1 ./internal/tools/` → `ok 21.402s` **rc=0**
- G2 `go test -count=1 ./internal/agent/` → `ok 4.477s` **rc=0**
- G3 `sh scripts/d22scan.sh` → `clean - no D22 ban violations` **rc=0**（ban #8 internal/=557）
- G4 `sh scripts/check-path-length-budget.sh` → `VERDICT GREEN … roster equals the tree` **rc=0**
- G5 `gofumpt -l`（`D:\work\base\gopath\bin\gofumpt`，我动过的两枚 .go）→ 空 **rc=0**
- G6 `go vet ./internal/tools/ ./internal/agent/` → **rc=0**

## 5. 具名没做（本腿射程外）

- AC#2b（`cmd/wisp` 接线）：非本腿面；**现读旁证** `cmd/wisp/run.go:547` 已带 `tools.TaskDeps{Roster…, Paths: rt.paths}`（145-r3 落过）⇒ 该格今天缺的是端到端三臂复核，不是接线行本身。
- AC#3／AC#4／AC#5：非本腿（AC#3 已由 174-v3 判、AC#4 owner 面、AC#5 结题动作）。
- 端到端 `wisp run`／`cmd/wisp` 整包／`-race`／`portable-tests.sh` 全量：未跑（前二者非本腿面＋本机缺 DLL 旧账 174-c1；全量属 CI 面）。
- `spill.go` 零字节；`internal/agent/**` 零字节（只读＋overlay）；`[fs]` 默认值／C26 语义／`docs/PLAN.md`／`docs/specs/**`／`risk/**`／`panel/**`／`thresholds.go`／golden／`allowlist.txt`／`frontend/**`／`design/**`／`cmd/**` 全零字节。
- 未动 `scripts/**`（新判据零 skip ⇒ 不需 ledger 行）；未碰别人脏件。

## 6. 复现原料（留盘未提交；规则 8 只建不删）

- `logs/raw/` 6 枚：01-ac2d-before／02-ac2d-after／03-ac2d-restore-and-teeth／04-ac2c-head／05-ac2c-silent-and-restore／raw-gates（全量原始捕获）。
- `mut/` 6 枚：`task-splice.go`（修前整枚，md5 `ded82068…`）／`task-roots.go`（`eb44b292…`）／`spill-silent.go`（`a0fc92bc…`）＋三枚 overlay json。
- 提交面的 `logs/01..05` 为上述 raw 的紧凑摘录（决定性逐字行＋rc）。

## 7. 与在飞/邻格的协调

- 133-r1 已交（46dee9ef）；`174-r8`（AC#2b，cmd/wisp，21:0x 派）与 `133-v` 在飞——与本腿**不同文件**；本腿两次 commit 全走显式 pathspec，只含自己的件。
- 邻格读数我未越权：`scripts/slo-freshness.sh` 的工作树脏件（136 的修复）与 `probes/133/**`、`probes/136/**` 一律未读未碰未归因。
