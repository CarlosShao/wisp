# 票 444 — 还有 7 枚竞态窗口仪器会在「弃权」时变红（票 441 普查交回，⛔ 不在 441 里假装已修）

**状态**：待派
**出处**：票 441 的 AC#4a 普查（docs/evidence/s1/441-instrument-abstention.md §4；尺＝sh .scratch/wisp/probes/441/census.sh，**37 行 / 13 文件**）。票 441 已结，只关它自己那一枚。

## 点名表（普查现取，⛔ 不是猜）

| 文件:行 | 判据 | 备注 |
|---|---|---|
| internal/winsec/openleaf_322_other_test.go:206 | 会弃权变红 | 文件头自陈「种一个叶子 link 会被 placement floor 先拒掉」⇒ **无确定性载体**；改成弃权＝没证人时变绿（⛔ 不做），正解只能是**换一种能确定性造窗口的载体**或诚实标注射程 |
| internal/winsec/ancestoropen_342_unix_test.go:117 | 会弃权变红 | 需单独载体 |
| internal/winsec/ancestorunlink_348_unix_test.go:153 | 会弃权变红 | 需单独载体 |
| internal/winsec/ancestoropen_342_windows_test.go:56 | 会弃权变红 | 需单独载体 |
| internal/winsec/ancestorunlink_348_windows_test.go:111 | 会弃权变红 | 需单独载体 |
| internal/winsec/deleteseam_348r2_test.go:131 | 会弃权变红 | 需单独载体 |
| internal/winsec/deleteseam_348r2_test.go:141 | 会弃权变红 | 需单独载体 |

其余普查行是台面／注入式**确定性**自检（不会弃权变红），另有两枚藏在 wisp342probe build tag 里、CI 根本不编译。

## 为什么值得修
每一枚都是一个**独立的 CI 假红源**：CI 负载下窗口不开 ⇒ 仪器报「护栏没了」⇒ 整条 	est-core/	est-windows 变红，而产品没毛病。票 441 那枚实测 **10 发里 3 红**，可见这不理论。假红一多，真的红就没人信了。

## 完成判据（AC）

- [ ] **AC#1 逐枚定性**：7 枚逐枚给「这一发窗口到不了时它现在报什么」，用受控争用（WISP_120_ATTEMPTS=1 那种缩分母，或并发负载）**真的让窗口不开**，读出现有那条红句。⛔ 不许写「大概率会」。**票 322 那枚若证实无确定性载体，具名说明并给出「诚实射程」写法，算通过。**
- [ ] **AC#2 逐枚改法**：照票 441 已确立的形（docs/evidence/s1/441-instrument-abstention.md）——**把「条件到达过」与「护栏还在」拆成两个载体**，弃权打**具名行＋数字＋可 grep token**（AC#1 ABSTAIN(<票号>:<shape>:<why>) attempts=... samples=...）。⛔ 不用 	.Skip（untests.sh 把 core 档顶层 SKIP 判致命）。
- [ ] **AC#3 每枚都要牙**：每一枚各造一发「护栏被拆」的突变（仓外导出树）⇒ 修后**真红并逐字点名**；还原 SHA256 同串、再跑绿。⛔ 只验其中一枚就算。
- [ ] **AC#4 ⛔ 不许钝化**：7 枚各自**原有的真红断言一条不许删/放宽**；改完把这 7 个测试的断言清单列成表（改前 vs 改后逐条对照）。
- [ ] **AC#5 门禁四数**：ash -n（若碰）；sh scripts/d22scan.sh clean；path-length GREEN；go vet ./internal/winsec rc=0；**产码 0 字节**；CI 上 	est-core 与 	est-windows 各**两发连续绿**。

## 边界
写面：internal/winsec/** 的 *_test.go、docs/evidence/s1/444-*、.scratch/wisp/probes/444/**、本票 Progress log。
⛔ 产码；⛔ 阈值/golden/契约；⛔ 加 build tag 把它藏出 CI（那是把 flake 改成透明）；⛔ dd -A／--amend／eset／ebase／stash／checkout .；⛔push。

## 撤销口令
「票 444 撤」＝删本票＋回退其写腿 commit（新 commit 反改，不 reset/amend）。
