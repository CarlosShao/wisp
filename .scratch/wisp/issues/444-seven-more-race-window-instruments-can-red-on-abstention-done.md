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

- [x] **AC#1 逐枚定性**：7 枚逐枚给「这一发窗口到不了时它现在报什么」，用受控争用（WISP_120_ATTEMPTS=1 那种缩分母，或并发负载）**真的让窗口不开**，读出现有那条红句。⛔ 不许写「大概率会」。**票 322 那枚若证实无确定性载体，具名说明并给出「诚实射程」写法，算通过。**
- [x] **AC#2 逐枚改法**：照票 441 已确立的形（docs/evidence/s1/441-instrument-abstention.md）——**把「条件到达过」与「护栏还在」拆成两个载体**，弃权打**具名行＋数字＋可 grep token**（AC#1 ABSTAIN(<票号>:<shape>:<why>) attempts=... samples=...）。⛔ 不用 	.Skip（untests.sh 把 core 档顶层 SKIP 判致命）。
- [x] **AC#3 每枚都要牙**：每一枚各造一发「护栏被拆」的突变（仓外导出树）⇒ 修后**真红并逐字点名**；还原 SHA256 同串、再跑绿。⛔ 只验其中一枚就算。
- [x] **AC#4 ⛔ 不许钝化**：7 枚各自**原有的真红断言一条不许删/放宽**；改完把这 7 个测试的断言清单列成表（改前 vs 改后逐条对照）。
- [x] **AC#5 门禁四数**：ash -n（若碰）；sh scripts/d22scan.sh clean；path-length GREEN；go vet ./internal/winsec rc=0；**产码 0 字节**；CI 上 	est-core 与 	est-windows 各**两发连续绿**。

## 边界
写面：internal/winsec/** 的 *_test.go、docs/evidence/s1/444-*、.scratch/wisp/probes/444/**、本票 Progress log。
⛔ 产码；⛔ 阈值/golden/契约；⛔ 加 build tag 把它藏出 CI（那是把 flake 改成透明）；⛔ dd -A／--amend／eset／ebase／stash／checkout .；⛔push。

## 撤销口令
「票 444 撤」＝删本票＋回退其写腿 commit（新 commit 反改，不 reset/amend）。

## Progress log

- **2026-11-14 17:1x +08 — 七枚全部改成「到达过」可弃权、「护栏还在」必真红，真人实测。** 形照 441：每枚新增一条**无竞态控制**（钉父 → 翻祖先 → 单次 act，读回外来见证与自己那片叶子），控制违例＝真红；racer 侧零仍红（护栏被拆的牙），route 侧零打 `ABSTAIN(444:<shape>:<why>) attempts=… samples=…` 并带该发数字。实测读数：① 基线红句 7 枚逐字取到（`WISP_444_ATTEMPTS` 缩分母 1/4/0，Windows 侧 6/6、5/5、2/2，unix 侧 3/3；日志 `.scratch/wisp/probes/444/logs/0[1-6]-*`）；② 另发现票面没列的第三条：`ancestoropen_342_unix_test.go:96` 的 `!landed` 分支原来 `t.Errorf` 后 **`return`**，把下面的 damage 断言整段跳过——该 `return` 已删，damage 断言反而**多跑**（严格更强）；③ 突变 7 发（仓外导出树 `D:\tmp\wisp\444\export`，每发还原后 SHA256 同串再绿）：`oNoFollowPlatform=0`、`Openat→Open(path)`、`Unlinkat→os.Remove`、`p.Root.OpenFile→os.OpenFile`、`p.Root.Remove→os.Remove`、seam 同 `Remove` 突变、`winsec.go` 删掉 refusal-path removal 调用 ⇒ 各自真红且逐字点名，其中删除 removal 那发**同时**打出 `AC#2 RED instrument: the removal seam was reached 0 times in 400 attempts` 与 `ABSTAIN(444:348r2-seam:removal-never-reached)`，证明弃权**没有**吞掉红；④ 受控争用（4 个 CPU 燃烧器 + `GOMAXPROCS=1`、committed 分母、count=3）九发全绿，且 322 的那一发 `leafLegRefusals=0 flips=130` **走的是弃权行**——正是过去会假红的那一发；⑤ 门禁：`go vet` 双平台 rc=0、`gofmt -l` 空、`bash -n` OK、`d22scan clean`、产码 0 字节；**path-length 是 RED 但不是本腿的**（唯一出册路径 `.scratch/wisp/issues/274-…-done.md`，`git log -1` 指向前一条腿），本腿新路径 36/37/45/56 字符全在册内。⑥ 诚实射程：缩分母到 20 时 racer 侧仍会红 3/6（racer 自己一轮都没翻完），到 1 时 6/6 红——那是台面假象、句子点名的是 racer，不是 CI 的状态；刻意**不**把 racer 侧零也改成弃权，否则删掉 racer 也永远绿。证据表：`docs/evidence/s1/444-instrument-abstention.md`。
---

## 编排者收件（2026-11-14 17:0x，非实现者，只追加）：AC#1–#4 翻，**AC#5 留着等 CI 两发连续绿**

**尺1 提交真在**：`d857ae57d`（`cat-file` rc=0）。**产码 0 字节我自己核**：`git show --name-only` 里非 `_test.go`、非 `probes/`/`docs/`/`scripts/` 的文件 = **0 枚**。

**尺2 件真在（都是现读）**
- `docs/evidence/s1/444-instrument-abstention.md` **23277 B**；可复跑脚本 `.scratch/wisp/probes/444/repro-444.sh` 在。
- **7 条具名弃权 token 逐条在盘上**：`ancestoropen_342_unix_test.go:138`、`ancestoropen_342_windows_test.go:76`、`ancestorunlink_348_unix_test.go:172`、`ancestorunlink_348_windows_test.go:128`、`deleteseam_348r2_windows_test.go:170`、`deleteseam_348r2_windows_test.go:192`、`openleaf_322_other_test.go:229` ⇒ **token 总数 = 7**，与票面点名表一一对上。
- 基线／突变日志 **22 枚入库**（`logs/01`–`05` 基线、`13`/`14` 起为突变），可逐条打开核。

**逐格**
- **AC#1** ✓ 7 枚逐枚给了「窗口不到时报什么」：Windows 侧 6/6、5/5、2/2 红；unix 侧 3/3×3 红（日志 `01`–`05`，我已抽读 `01`/`03` 头几行确认它们真是对应的测试）。票 322 那枚**按票面要求**给了「诚实射程」写法（不是硬改成绿）。
  ⚠ 我**没有**逐发复现全部 7 枚的争用（那是 9 发 × minutes），我复跑的是**修后**那一发（见 AC#5 下的尺3）。
- **AC#2** ✓ 照 441 的形：到达性可弃权（具名行＋attempts/samples/flips 等数字），护栏本身仍真红；腿另报一条**关键的**：删掉 removal 那发**同时**打 `AC#2 RED …` 与 `ABSTAIN(…)` ⇒ **弃权没吞红**，这正是我要的证据。
- **AC#3** ✓ 七发突变（仓外 `D:\tmp\wisp\444\export`）：`oNoFollowPlatform=0`、`Openat→Open(path)`、`Unlinkat→os.Remove`、`Root.OpenFile→os.OpenFile`、`Root.Remove→os.Remove`×2、`winsec.go` 删 refusal-path removal 调用 ⇒ 每发真红并逐字点名，还原后 SHA256 同串（`19F3C211…`/`22D5D118…`/`1C15A7DD…`/`F0BBF9FF…`）再绿。
- **AC#4** ✓ 断言语义未被钝化：**342-unix 那枚的 `return` 被删掉后，下面的 damage 断言反而多跑**（不是少跑）；`calls!=attempts` 那条 wiring 红一字未动；措辞 token 只增不改。
- **AC#5** ⏳ **一半已绿、一半等 CI**：`gofmt -l` 空、`go vet` 双平台 rc=0、`bash -n` OK、`d22scan` clean（482 产码文件）、**产码 0 字节**、7 条定向测试绿色（腿跑的）；**`path-length` 那条它报的红是我自己造的**（274 结票文件名 102 字符），我已改短并现跑 `VERDICT GREEN`；AC#5 还要求 **CI 上 `test-core` 与 `test-windows` 各两发连续绿** ⇒ **必须等 CI 跑完才敢勾**，⛔ 不拿本地绿冒充 CI 绿。

**并行带出的两个不在它写面的红（不吸收，只记）**：`scripts/check-path-length-budget.sh` 的 274 文件名（我修的）；`internal/models/acl_sid_121_test.go` 的 gofumpt（CRLF→LF，历史件）。
**2026-11-14 18:0x 编排者补（AC#5 的 CI 那一半，现读）**：三发 **6/6 job 全 success** —— `43578030753`（head `721a1d8a0`）、`43578127657`（head `5b56792c6`）、`43578228598`（head `c64a4fa22`）⇒ 「`test-core` 与 `test-windows` **两发连续绿**」满足（实际有三发）。⚠ 我 17:1x 写「等 CI 两发连续绿」时它们还没跑完，是**独立审计腿 `dedb4089d` 查出已经够了**才翻的 —— 不是我事后放宽标准。
