
## 编排者收 `307-v1`（2026-10-11 11:0x，非实现者终裁腿五笔 `c04ced26`→`ed1d2bad`→`b2fb50e9`→`aebceb5e`→`94c2808c`，件 `probes/307/v1/**`＋裁决表 `docs/evidence/s1/ticket307-307-r1-accept-v1.md`；⛔ 上面各节原句一字不改，本节追加；台账＝`A839`）⇒ **五格逐格翻勾（本票 6 勾／0 未勾）**，凭据＝下面我这发**逐把自己重跑**的尺

### 0. 台面与越界（⛔ 抄腿自陈）

- 取读数时 HEAD＝`a6267fbb`（其间别家长了两笔：`94c2808c`＝本腿更正笔、`a6267fbb`＝`72-a2`）；产码面 `git diff --name-only f6f4eca8..HEAD -- internal cmd`＝**0 行** ⇒ 上面所有行号与本节行号同一份码。
- 五笔越界我自己逐笔重跑（尺＝`git show --name-only --format=`，⛔ 区间尺）：`c04ced26` 1 枚／落自家外 **0**、`ed1d2bad` 27／**0**、`b2fb50e9` 6／**1**、`aebceb5e` 2／**1**、`94c2808c` 4／**1**。⚠ 那三枚"外"分别是 `docs/evidence/s1/ticket307-307-r1-accept-v1.md`（**AGENTS §1.5「裁决表：`docs/evidence/s1/`」逐字规定的落点**）与本票末尾那条 Progress bullet（`1 加／0 删`）⇒ **⛔ 算越界**；⚠ **记我派单一处**：我给验收腿写的"只落 `probes/307/v1/**`"这把尺**比规矩窄**，同一枚我在 `A838` §4 刚记过第二遍 ⇒ 模板那句要改成"自家 `probes/**` ＋裁决表 `docs/evidence/s1/**` ＋本票末尾 bullet"。
- 腿的写面⛔ 碰产码：`git status --porcelain -- internal cmd`＝**0 行**（四形突变全在仓外树做、每发还原；我这发复跑也验了还原，见 §1 第 3 条）。

### 1. ★承重读数我这发逐把自己复跑（每条带尺；⛔ 引用本节前先重跑）

1. **落地形状＝我裁的形 C，逐字对得上**（尺＝`sed -n '415,450p' internal/risk/syncdirs.go`，HEAD）：`:433` `candForms := formsOf(cand)`（★只算一次，⛔ 每枚根重走祖先）、`:434` `if aUnder(candForms, anchorForms(e.canon))`（原 `isUnder(cand, e.canon)`）、`:442` `if !s.complete && s.home != "" && aUnder(candForms, anchorForms(s.home))`（原 `isUnder(cand, s.home)`）⇒ **两处出口同换、`s.home != ""` 那枚守卫原地留着**（＝追加节 §3 第 4 点）。**存值⛔ 动**：`grep -n 'home: normPath' internal/risk/*.go`＝**2 行**逐字未变（`provenance.go:345`＋`syncdirs.go:88`）；`Root.Path` 仍报 `s.home` 那枚比较形 ⇒ `A834` 登记①那条**照原样欠着**，⛔ 因翻勾而被当成已修。
2. **两把台面的颜色（我自己四发，⛔ 抄腿）**：尺＝同一份 HEAD 码**导出到仓外树**、`go test ./internal/risk/ -count=1 -timeout 600s -v`，只差 `TEMP/TMP/TMPDIR` 一枚 env；建树用**改好的那把尺** `git archive HEAD go.mod go.sum internal cmd | tar -x`（⇒ `archive_rc=0`、`tree_go_files=648`、测试真跑；腿复现我原尺那次＝`rc=1`＋`go: go.mod file not found`，**我认这枚记在我账上**）。件＝`probes/orch/20261011/logs/orch-307-twosurface.txt`（＋两枚 `*.tmp-*` 原始件）。
   - 别名台面（短形 `C:\Users\swq\tmp\RISKAT~1\RUNNER~1`，我自己早先那发里用 COM `ShortPath` 取的）＝**`rc=1`／`PASS=192`／`FAIL=7`／`SKIP=1`**，七枚红**逐名**＝`TestClassifyAnchorSpellingIsNotVerdict`／`TestCanonicalInputGainsNoSecondForm`／`TestAListWinsWhereBothTablesHit`／`TestPathResolverShortNameAListDenied`／`TestPathResolverUNCAListDenied`／`TestPathResolverExtendedLengthPrefixAListDenied`／`TestBListDefaultDenyAndOverride` ⇒ **全是票 72 甲族、乙族那 5 枚已由 FAIL→PASS**；
   - 本机默认台面＝**`rc=0`／`PASS=199`／`FAIL=0`／`SKIP=1`** ⇒ 本机颜色未被动；
   - ⚠ 口径差具名：我这把尺只数 `--- PASS/FAIL/SKIP` 行（⇒ 192／199），腿的"四数"第一枚是 `RUN`（200＝192＋7＋1）⇒ **同一次运行两种打印**，引数必带尺名；`SKIP=1` 恒为 `TestSyncRegistryProbeLive`（腿报的枚名与我这发一致）。
3. **反形我自己三发（★AC#2 那格的凭据⛔ 是腿的一面之词）**，全在仓外树、每发之后还原＋`sha256sum` 自证：件＝`logs/orch-307-reverse-shapes.txt`。
   - 形 A（**两支同退**回 `isUnder`，＝把本票修复整个挪开）⇒ `rc=1`／`PASS=185`／**`FAIL=14`**／`SKIP=1`，红名册＝甲 7 ＋ 乙族指名 5（`TestSyncFixtureFallbackAndMatch`、`TestSyncFallbackNotDisarmableByWeakRoot`、`TestSyncUnverifiedRootKeepsFallback`、`TestSyncSuspectFallbackIsComponentBounded`、`TestSyncSuspectFallbackWhenUndetectable`）＋**新夹具 2 枚**（`TestSync307AliasHomeUnderProfileStaysSuspect`、`TestSync307AliasShortRootMembership`）⇒ 修复**是承重的**，⛔ 它只是让别的红变绿；
   - 形 C（**只退 `:434` 那一支**＝腿报的 (c)）⇒ `rc=1`／`PASS=191`／**`FAIL=8`**，其中 sync 族**唯一**那一枚红＝**`TestSync307AliasShortRootMembership`**，而乙族那 5 枚在这一形下**全绿** ⇒ ★**根成员判定今天由这一枚新夹具单独钉住**，与腿报的"`FAIL=1`＝`syncdirs_alias_307_windows_test.go:120`"同结论（⚠ 射程不同：腿打的是名册档、我打的是整包档，8 枚里 7 枚是票 72 的）；
   - 正控（pristine 树、别名台面）⇒ `FAIL=7`＝甲族那 7 枚，⛔ 多一枚 ⇒ 我这两发的红⛔ 是噪声；
   - ★还原互证：我这发的 `sha256(syncdirs.go)`＝`aa4779bb5270…`，**与腿在裁决表里报的前缀逐字同** ⇒ 两枚独立的手摸到的是同一份件（腿那把尺是 `Get-FileHash`、我的是 `sha256sum`）。
4. **夹具⛔ 依赖机器形状（★AC#3 那格我自己量的那把词面尺）**：`internal/risk/syncdirs_alias_307_windows_test.go`＝141 行、`:1` 逐字 `//go:build windows`、`t.Skip`＝**0**、`t.Fatal*`＝**15**、`tryShortPath`＝**2 命中**、`shortPathOf`＝**0**、`getShortPath`＝**0** ⇒ 票 115 那枚"种不出即 `t.Fatal`、⛔ Skip"的定式真落地，带逃逸口的两枚旧手⛔ 被走；我另读了它的正文（`:100-130`）：它**显式**断言 `st.Root.Source == "suspect-fallback"` 即 `t.Fatalf("verdict must come from ROOT MEMBERSHIP, not from the fallback…")` ⇒ 这一枚正是形 C 唯一响的那枚，机制自洽。
5. **门禁我自己三把**：`GOFLAGS= go build ./...`＝**0**、`go vet ./internal/risk/`＝**0**、`sh scripts/d22scan.sh`＝**0**（件＝`logs/orch-307-gates.txt`，三行 `rc=` 自落，⛔ 0 字节件）。
6. ⚠ **格式名册那一把我纠我自己也纠腿的尺**：腿报"工作树 1 枚（`provenance.go`）⇄ HEAD blob 0 枚"，射程＝`internal/risk`；我把同一把尺放宽到 `gofmt -l internal cmd`＝**5 枚**（`internal/agent/approval/pending_read.go`、`internal/agent/tools.go`、`internal/risk/provenance.go`、`internal/tools/bridge.go`、`cmd/wisp/models.go`），逐枚 `git ls-files --eol`＝**全是 `i/lf w/crlf attr/text eol=lf`**，且 `git status --porcelain -- internal cmd`＝**0 行** ⇒ **五枚全＝CRLF 假枚、零真编辑**；blob 侧我自己复跑（`git show HEAD:internal/risk/provenance.go > f && gofmt -l f`）＝**空** ⇒ 腿那格"blob 0 枚"成立。⇒ 这是"引数必带射程"那族口径的**第四枚实例**，⛔ 谁把 1 与 5 当成两次打架的读数。
7. **生产面（腿的必答第 1 问，我自己量的同形）**：`IsSyncPath` 非测试调用者＝**1**（`provenance.go:928`）；非测试 `NewProvenance`＝**2**（`cmd/wisp/run.go:743`、`cmd/wisp/panel_assets.go:239`）**都只传 `NoProbe: true`**；`.SyncRoots()` 非测试＝**0**。⇒ 腿那句"这台机器今天在生产上观测不到颜色差别"我认；⛔ 但⛔ 读成"没接"——我这发把判定往门外追了四跳（`bridge.go:931 WithTaintDetector(b.prov.Detector(taskID))` → `assessor.go:168` 的 C19 缝 → `rules_gateway.go:105 ctx.taint.TaintHit(...)` ⇒ `R4` → `provenance.go:725 Inspect` → `:755/:908 writeGate` → `:928 IsSyncPath`），全表在 **`A838` §1** ⇒ **本票翻勾的凭据性质＝测试台面凭据**（⛔ 生产台面凭据、⛔ 拿 CI runner 的 `TEMP` 别名替生产背书）。

### 2. ★裁：五格逐格翻勾，三枚 caveat 跟着勾一起走（裁决者＝`307-v1`，⛔ 实现者；落笔＝编排者）

- `AC#1` ⇒ **勾**，但票面上这格从此带一句：凭据性质＝**测试台面凭据**（§1 第 7 条），⛔ 只写"修复生效"。
- `AC#2` ⇒ **勾**；那枚"乙族 2 枚⛔ 负向钉"的欠账（`syncdirs_test.go:209`/`:212`/`:216`、`:229`/`:232` 体内只有正向断言，恒真形下它们绿）**归票 72 下一程**（同文件、同族成因、甲族期望侧本来就归它）＝裁语，撤销口令**「307 负向钉另立一票」**；⛔ 塞回本票（票面"边界"第 1 条禁的就是动那 7 枚的期望侧）、⛔ 归票 302（它只管档位／搬运）。
- `AC#3` ⇒ **勾**，附两枚**具名欠读数**（⛔ 因翻勾而消失）：`alias307Pair` 的 `t.Fatalf` 分支本机从未执行（要一台⛔ 8.3 别名的机器才走得到）；新夹具在 POSIX 档⛔ 求值（`//go:build windows`＝票面授权的形状，票 301 包级 guard D）。
- `AC#4` ⇒ **勾**；那枚⛔ 可证的格认（`git add` 的命令行本身⛔ 入库，只能靠"五笔名册⛔ 一枚落在授权外"旁证）；另按 §1 第 6 条把格式名册那格**改写成带射程的两句**，⛔ 让下一位拿 1 与 5 打架。
- `AC#5` ⇒ **勾**；腿那把"红名册文案层对拉"的尺（甲 7 枚改前⇄改后只剩时长与 `t.TempDir()` 随机后缀）与我这发独立取的**名册逐名**结果同向 ⇒ **甲 7→票 72／乙 5→票 307 的分家成立**，`A817` 欠账⑥**由本票结掉**（正文落 `A839`）。
- ⚠ 腿自抓两枚数位（`94c2808c`：把目测点数当尺的读数；把"结论"写成"0 命中"）＝**处置正确**（⛔ 擦历史、另笔更正），我复核两枚⛔ 改判语；那条定式（凡写进判语的枚数必须直接来自那一发打印的数字，且"0 命中"必配一把把命中的那枚打印出来的尺）**进我车道记忆**。

### 3. 残余与永久欠读数（具名，全部落在别处，⛔ 说"零残余"）

- ★**永久欠读数一枚**＝`formsOf(cand)` 那次祖先走步的代价：腿与落地腿两把尺（含 `-run 'TestSync'` 前后总耗时＋一枚具名用例的 `--- PASS (0.xx s)`）都量不出来，只剩结构性上界"每查询 O(1)＋锚侧有缓存" ⇒ 裁＝**⛔ 为它立仪器票**，只记成永久欠读数，撤销口令**「307 补 bench 尺」**（⛔ 因为量不出来就说代价为零）。
- `:431` 报 `normPath` 形为 `Root.Path`（`A834` 登记①，本轮⛔ 修、⛔ 顺手"改对"）。
- `SyncRoots()`／`SyncDetectionComplete()` 生产调用者＝0，两处注释按一枚尚未建的 Security page 写现在时 ⇒ **已归票 40**（`A838` §2/§3）。
- `bridge.go:925` 那三枚前提里 `taskID == ""` 在真跑时会不会发生（发生 ⇒ R4 探测器整枚不挂）＝**`A838` 新抓的具名欠读数**，归我车道下一发只读尺。
- 逐卷 `fsutil 8dot3name query C:` 需提权（`Error 5`）＝⛔ 我取得到。
- 别名台面那 7 枚（票 72 甲族）仍是红 ⇒ 本票结掉⛔ 让 `internal/risk` 全绿；CI 色与推送照旧在我车道（`@{u}`＝`cnb/dev` 落后 1032，`A813` 未过期、我仍⛔ 用）。

### 4. 收口与 next

- **收口动作**（README 规则 4）：六格全勾 ⇒ 本节写 `Status: done` ⇒ 同笔 `git mv` 加 `-done` 后缀；`-done` ⛔ 是"零残余"证明（README `L201` 原话），残余全在 §3 那六枚里、各有归宿。
- next＝收 `72-a2`（在飞，已交两笔 `493f2cb8`／`a6267fbb`）⇒ 我据它的七问＋本票 §2 那枚负向钉欠账**派 72 的写腿**（甲族期望侧＋两枚负向钉同批，⛔ 两批都动 `syncdirs_test.go` 之外的文件）→ 我自己那发只读尺答 `taskID==""` 那问 → 票 305 那六发（要编队静的窗口、⛔ bisect）→ 补票 300 `AC#4` → 推送与 CI 色同批。队列 `298-v1`→`296-r2`→`295-r1`→`294-r1`→`294-v1`。
- Progress log：
  - [2026-10-11 11:0x +0800] agent=编排者 did=收 `307-v1`（五笔逐笔越界我自己重跑＝落在授权面外 3 枚全是裁决表／本票 bullet；仓外树建树尺改用带 `go.mod go.sum` 那把 ⇒ `archive_rc=0`／648 枚 `.go`／测试真跑）＋**承重读数逐把自己复跑**：两台面四发（别名 `rc=1`/`PASS=192`/`FAIL=7`/`SKIP=1`，七枚红逐名全＝票 72 甲族；本机 `rc=0`/`199`/`0`/`1`）、反形三发（两支同退 `FAIL=14`；只退 `:434` ⇒ sync 族唯一红＝`TestSync307AliasShortRootMembership`；pristine 正控 `FAIL=7`；还原 `sha256 aa4779bb…` 与腿报同前缀）、夹具词面尺（`//go:build windows`／`t.Skip`0／`t.Fatal*`15／`tryShortPath`2／`shortPathOf`与`getShortPath`各 0）、门禁三把各 `rc=0`、格式名册两档射程（`internal/risk` 工作树 1 ⇄ `internal cmd` 工作树 **5**，`git ls-files --eol` 逐枚＝`i/lf w/crlf` ⇒ 全 CRLF 假枚；HEAD blob 0）、生产面（`IsSyncPath` 非测试 1、`NewProvenance` 非测试 2 全 `NoProbe:true`、`.SyncRoots()` 非测试 0）＋往门外追四跳（`bridge.go:931`→`assessor.go:168`→`rules_gateway.go:105`⇒`R4`→`provenance.go:725/:755/:908/:928`）⇒ **翻 `AC#1`–`AC#5` 五格（本票 6 勾／0 未勾）并三枚 caveat 跟勾走**：AC#1 的凭据性质＝**测试台面**、AC#2 那枚负向钉欠账**归票 72 下一程**（口令「307 负向钉另立一票」）、AC#3 两枚欠读数照记；`formsOf` 那枚代价值裁＝**⛔ 立仪器票、记永久欠读数**（口令「307 补 bench 尺」）；`A817` 欠账⑥**结掉**（甲 7→票 72／乙 5→本票，文案层对拉与我逐名复跑同向）。⇒ 收口 `Status: done`＋`git mv` 加 `-done`（⛔ 等于零残余：§3 六枚残余各有归宿）。⛔ push（落后 `cnb/dev` 1032，`A813` 我仍⛔ 用）。

**Status:** **done**（README 规则 4 收口：六格全勾 → 本节写 `Status: done` → 同笔 `git mv` 加 `-done` 后缀。收口凭据＝上面 §1 那七条我自己重跑的尺；残余＝§3 那六枚，逐枚已落到票 72／票 40／我车道，⛔ 说"零残余"。）
