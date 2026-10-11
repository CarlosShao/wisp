# 票 308 — 拼法不变量那枚用例在 `filepath.VolumeName` 返回空串时**直接 panic 掉整枚测试二进制**：`vol[:1]` 零守卫 ⇒ 一次"临时目录没有卷标"的运行里，同包其余用例**零求值**（形状＝看着像红、实际是"这一发根本没跑完"）

**立票**：2026-10-11 11:1x 编排者（来路＝只读普查腿 `72-a2` 的自报第 4 问＋**我自己读了崩点与它的原始件**；台账 `A840` §2）
**性质**★：**仪器票**（测试面健壮性），⛔ 产码票。触发面窄（要一枚无卷标的 temp），但坏形状是"整包静默少跑"，与票 300／301／302 那一族"看着像读数、其实是⛔ 跑"同族。
**为什么要紧（用户可见形状）**：这枚用例的名字就是"同一份文件、换任何拼法都该给同一个分类"——它本来的职责是**守住拼法不变量**。它一 panic，那一发里同包**其余几十枚**（含票 307 刚落地的别名夹具）⛔ 出结果，而 CI／本地都只看到 `rc=1` 与两行名册 ⇒ 拿它当"这个包红了两枚"来读会**读错射程**。⚠ ⛔ 生产写被绕过：产码那两处同款切片我自己看过（下表现量）。

## 现量（每条带尺；行号一律当 11:1x 快照）

- ★**崩点＝测试面两枚同形**（尺＝`grep -rn '\[:1\]' internal/risk`）：
  - `internal/risk/pathresolver_anchor_spelling_windows_test.go:77 vol := filepath.VolumeName(long)` → `:78 share := strings.ToLower(vol[:1]) + "$"`；
  - `internal/risk/pathresolver_junction_windows_test.go:151 vol := filepath.VolumeName(long) // e.g. C:` → `:152 share := strings.ToLower(vol[:1]) + "$"`；
  - `tmpHome(t)`＝`pathresolver_junction_windows_test.go:22-25`，体内 `home := t.TempDir()`＋`t.Setenv("USERPROFILE", home)` ⇒ `long` 的卷标完全**由 temp env 决定**。
- ★**panic 逐字**（件＝`probes/72/a2/logs/t03-broken-surface-note.txt`，24 行，腿把被覆盖的那发原样保了一份）＝`panic: runtime error: slice bounds out of range [:1] with length 0 [recovered, repanicked]`；那一发四数＝**`rc=1`／名册只 2 枚／`PASS=0`／`FAIL=2`**（腿自己在件里写明"`PASS=0` ⛔ 代表'全红'，是测试二进制在**第 2 枚**就 panic 死了"）；⚠ 那一发的命令**没设任何 temp env**（工具跨命令⛔ 保留 shell 状态，Go 从 Git Bash 继承无卷标的 POSIX 形 `TMPDIR`）⇒ 触发条件＝**env 形状**，⛔ 是"某台机器坏了"。
- ★**产码⛔ 在射程内**（我这发自己读的，两条都带尺）：`internal/risk/pathresolver.go:262-263` 切 `drive[:1]` 之前有 `if len(drive) == 2 && drive[1] == '$' && …` 守卫；`internal/risk/pathresolver_windows.go:61 reparseComponents()` 取了 `vol` 但**⛔ 做 `[:1]` 切片** ⇒ **本票⛔ 动产码**，射程＝上面那两枚测试文件。
- ⚠ **一枚注释过期（登记、⛔ 顺手改对）**：`pathresolver_anchor_spelling_windows_test.go:70-73` 逐字写着"every other spelling is unconditional, so the test never quietly stops measuring"——在 `vol == ""` 那一形下这半句为假（它⛔ 是 quietly，但它确实 **stops measuring**，因为整枚二进制没了）。
- **枚数**：两枚用例会走那对切片（尺＝上面那条 `grep`，命中＝2 枚）；同文件其余用例（`TestClassifyAnchorSpellingIsNotVerdict` 等）⛔ 碰 `VolumeName`。

## 要建什么

- [ ] **`AC#0` 先把"哪一发真能崩"现量清楚（只读）**：在**无卷标 temp**（如 `TMPDIR=/tmp` 那形）与**有卷标**两把台面上各跑一次那两枚用例所在的名册，逐枚交 `rc`＋`RUN/PASS/FAIL/SKIP`＋panic 行逐字；⛔ 预设"只有 `anchor_spelling` 会崩"（两枚同形，另一枚的实测颜色⛔ 知）。
- [ ] **`AC#1` 加一枚守卫**：判据＝无卷标那一形下，这枚用例**以可解释的颜色结束**（⛔ panic、⛔ 整枚二进制没了）；★**守卫的落点⛔ 是"把整枚用例 `t.Skip` 掉"**——那正是它守的不变量还剩下长/大写/小写/扩展长度几种拼法照跑，只有 `unc-admin` 那一枚拼法在⛔ 卷标的 temp 上无法构造（可具名去掉**那一项**，⛔ 去掉整枚）。硬约束＝⛔ 改任何判级逻辑、⛔ 动 `pathresolver.go`／`pathresolver_windows.go`、⛔ 把 `t.TempDir()` 换成手写常量（票 72 `:201` 那条禁线同形）。
- [ ] **`AC#2` 反形敏感（⛔ 只交"改⇒绿"）**：必须各交一发凭据＝(a) 守卫**拔掉** ⇒ 无卷标台面那一发**回到 panic**（红句逐字）；(b) `cmp`／`sha256` 逐字节还原 ⇒ 复绿；(c) ★**别名台面颜色⛔ 因这枚守卫而变**（改前⇄改后逐名同集合）——否则"守卫"其实是把用例关掉了。
- [ ] **`AC#3` 分母与档位**：两枚件都带 `//go:build windows` ⇒ 具名回答"CI 哪一档真跑它们、`windows` tag 在 `test-windows` 里求不求值"（先例＝票 301 包级 guard D、票 111 那 12 枚的"从未求值"）；⛔ 用"本机跑过"顶替 CI 色；拿不到 CI 色就具名登记欠读数（⛔ 由落地腿自己盖章）。
- [ ] **`AC#4` 门禁与越界**：`GOFLAGS= go build ./...`／`go vet ./internal/risk/`／`sh scripts/d22scan.sh` 各 `rc=0` 且**每把件自落一行 `rc=N`**；格式名册**两把并排且各写射程**（⚠ 本轮刚又实证一次：射程 `internal/risk` 工作树＝1 枚 CRLF 假枚、放宽到 `internal cmd`＝5 枚，见 `A839` §1 末行）；`git show --name-only --format=` **逐笔**量越界；⛔ push；commit 显式 pathspec；写面＝自家 `probes/308/**` ＋裁决表 `docs/evidence/s1/**` ＋本票末尾 bullet（⛔ 再拿"只落 `probes/**`"那把窄尺误报，见 `A838` §4／`A839` §0）。
- [ ] **`AC#5` 对账**：⛔ 翻票 72 与票 307 任何一格；本票与票 72 的关系＝**同一枚文件、不同件事**（票 72 管"期望侧写死了 temp 拼法"，本票管"崩点让整包零求值"）；⛔ 扩票 307（它已 `-done`，残余见 `A839` §3）。

## 边界（⛔ 塞进本票）

- ⛔ 甲族那 7 枚的期望侧（归票 72，`A840` §2 已换好取证尺并派腿）。
- ⛔ `overrideApplies` 那条放宽（票 72 `:95` 逐字作废，且今天它有真生产 caller `internal/tools/mode.go:91/:93/:98`，见 `A840` §1）。
- ⛔ 把票 300／301／302 那三档"windows 档在 CI 恒零求值"的搬运活并进本票（本票只交自己那两枚的档位读数）。

## 规矩（本票全程）

- 子代理**只 commit、⛔ push**；⛔ `git add -A`／`.`；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`；临时件**只建不删**且从一开始就落自家 `logs/`（⚠ 仓库根现存 **48** 枚深度 1 未跟踪件，`A840` §1 已登记，⛔ 再往根上丢）；证据件⛔ 叫 `.out`。
- ⛔ 为变绿放宽任何断言；⛔ 造 `--- SKIP` 换绿（本票 `AC#1` 那枚"具名去掉一项拼法"要写清它⛔ 是 skip 整枚用例）；⛔ 动 SLO 阈值／golden／`thresholds.go`／`allowlist.txt`／D43 转移表／C1–C32／D1–D47／三枚冻结件。
- ★票名⛔ 超 100 字符（长度门＝票 262 那把尺 `sh scripts/check-path-length-budget.sh --with-self-test`，本轮 11:0x 读数 **GREEN**、`worst=228`、`over the hat=57 ⇄ roster=57`）⇒ 改名或新建前先跑它。
- 裁决者≠实现者（`SPEC-12 §4.3`）：`AC#1`/`AC#2`/`AC#3` 的判语⛔ 由落地腿自勾。

next＝`308-a1`（只读答 `AC#0`：两把台面各自的颜色＋panic 逐字＋CI 档位那问）→ 我裁守卫落点 → `308-r1`（写腿）→ `308-v1`（非实现者终裁）。⚠ 排在票 72 写腿**之后**（同抢 `internal/risk` 的 `*_test.go` 写面）。
