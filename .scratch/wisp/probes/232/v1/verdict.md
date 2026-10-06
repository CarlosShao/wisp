# 票 232 终裁件 · 腿 `232-v1`（非实现者）

裁决对象：`cmd/wisp/config_reload_223_test.go`（提交 `3d9b8374`，实现者＝`232-r2`，编排者代提）。
本件所有判语建立在**本腿自己造的突变读数**上；`232-r2` 的日志与自述**未被当作任何一格的凭据**。
本腿全部突变走 `go test -overlay`，**共享工作树零编辑**（终态 `git status --porcelain -- cmd internal`＝起手 0 行）。
本腿不勾任何 `- [ ]`、不写"完成"、不动 `docs/**` 与台账。

图例：`旧形`＝本腿用 git blob `ae9f346b`（＝ `a16d1ff7:cmd/wisp/config_reload_223_test.go` ＝ `3d9b8374` 的父态测试文件）经 overlay 重建的**232 之前那一把尺**（它是从 git 对象抽出的旧形，⚠ 不是盘上当前内容，也不是本腿编造）；`终态`＝盘上 HEAD 的那把尺；`HEAD 产码`＝`cmd/wisp/config_reload.go` 当前字节（`5ce441ca…`，全程未改）。

---

## §0 起手四把尺（原文读数）

**尺 1 — HEAD／近三笔**（本腿 18:39 取；派单时的 `1e8200b7` 之后已漂）
```
d6243757 Tue Oct 6 18:26:49 2026 +0800 probes(A642 落账)：机主桌面连弹 3 次"选择一个应用打开 projects"＝我造的不是外人，根因逐字取到
---
d6243757 probes(A642 落账)：机主桌面连弹 3 次"选择一个应用打开 projects"＝我造的不是外人，根因逐字取到
1e8200b7 probes(236-v1 交件): 票 236 三格终裁完成——§1 现状／§2 逐格判语／§3 名册 29 发全填实，三格一律判"成立"
55324446 probes(236-v1 §4 门禁＋§5 判不动): 票 236 终裁腿两节写满并 commit（28 发 overlay 全读数在案）
```

**尺 2 — 写面**：`git status --porcelain -- cmd internal` ⇒ **0 行**（起手 0，终态复量 0，见 §4 G5）。

**尺 3 — 产品文案对拉**
```
5ce441ca5e72b64d18a6c26f1c066882 *cmd/wisp/config_reload.go
5ce441ca5e72b64d18a6c26f1c066882 *-            (git show HEAD:cmd/wisp/config_reload.go | md5sum)
```
⇒ 相同；20 发突变跑完后复量仍为同值（§4 G5）。

**尺 4 — 基线三枚具名用例**。第一发（18:39，`tail -40` 把 RestartTier 那一行截在窗口外，只留可读到的部分）：
```
--- PASS: TestTicket223FailureSentencesAreDistinct (12.26s)
--- PASS: TestTicket223PanelInboundSaysHotReloadIsDisabled (0.01s)
ok  	github.com/CarlosShao/wisp/cmd/wisp	15.425s
```
为把三枚名字逐枚取全，本腿 19:1x 重跑并留盘（`logs/A0-baseline-3named.txt`）：
```
--- PASS: TestTicket223RestartTierSaysItWillNotApply (2.36s)
--- PASS: TestTicket223FailureSentencesAreDistinct (8.64s)
--- PASS: TestTicket223PanelInboundSaysHotReloadIsDisabled (0.01s)
ok  	github.com/CarlosShao/wisp/cmd/wisp	11.088s
```
`--- FAIL` 计数＝**0**。★本件判语只钉这三枚名字（加 AC#5 用的 `TestTicket223R2FailureSentenceRouting`）。

**★关键不对称（本腿自己读产品源文 `cmd/wisp/config_reload.go:317-326`，另有一处独立实测见 §3 `P-final`）**
`重启进程后生效` 这五字**只在 `rt.auditf` 那一枚审计 detail（`:319`）**；操作员那句 `fmt.Fprintf(rt.stdout, …)`（`:321-326`）里
有的是 `本次运行不会生效`／`需要重启进程`／`原因：`／`涉及：%s`／`没有被丢掉`，**没有那五字**。

---

## §1 现状：`232-r2` 改了什么（自读 `git show 3d9b8374`，不读它的自述件）

改动**只有一枚文件**：`cmd/wisp/config_reload_223_test.go`（`git show --stat` 里 `cmd/` 只有这一枚；产品码零字节，尺 3 已拉平）。

1. **新增四枚文件内小写 helper**：`windowSince232`（起点快照；`HasPrefix` 不成立就 `t.Fatalf`，不静默换窗）·`stdoutSince232`·`awaitStdoutSince232`·`awaitAuditSince232`（现盘 `:174/:185/:200/:223`）。零新增导出名。
2. **起点标记**：`TestTicket223RestartTierSaysItWillNotApply` 在 `plant` **之前**取两枚标记（stdout＋stderr，`:566-567`），随后三枚 await 全换成 `…Since232`（`:569`/`:573`/`:592`）。
3. **needle 的落点变了（本票核心）**：旧尺 `for needle in {app.autostart, 开机自启, 重启进程后生效} { Contains(why+out, needle) }`＝**拼接串**，任一路出现就绿；
   新尺把六枚内容 needle（`app.autostart`／`开机自启`／`需要重启进程`／`原因：`／`交给平台层`／`没有被丢掉`）钉在 **stdout 那一路的种后窗口**（`:605-609`），
   并把 `重启进程后生效` **单独钉到具名审计流窗口**（`:610-613`）。⇒ 全文件 `why+out` 计数：旧形 1、终态 **0**（本腿 grep）。
4. **删除列落点**：diff 里 9 枚删除全部落在它自己替换掉的 `why+out` 那一块与三枚 await 调用行。本腿另用**断言计数**独立复核同一枚用例体：旧形 `t.Errorf`5／`t.Fatal`2／`strings.Contains`4 ⇒ 终态 **6／2／5**（只增不减，§3 `X-count`）。
5. **负向断言读得更宽不是更窄**：`这些段已立即生效：[app]` 那一枚旧尺读 `out`，新尺读 `r.h.out.String()`（整流）。
6. **未动**：`reloadCaseBudget`=40s、Timer+20ms Ticker 形状、无 `Skip`、三枚冻结件、`thresholds.go`、golden、产品文案（§4 逐把复量）。

一句话：**它把"钉句子存在"改成了"钉句子在哪一路、且在种下改动之后才出现"**；唯一让步＝产品真源没有那五字，那一枚只能钉审计流（§5①）。

---

## §2 逐格判语表（六格）

| 格 | 判语 | 本腿自造读数（凭据） | 缺的那一行读数／残余形状 |
|---|---|---|---|
| **AC#1** 复现瘦句不响 | **成立** | 前提 `oldform-only`（旧尺＋HEAD 产码）PASS 2.88s；**B**oldform PASS **3.16s**；**C** PASS **3.08s**；**E** PASS **3.07s**；**F** PASS **3.67s**；正控 **K** FAIL 43.17s（旧尺独有红句＝旧尺确已落地） | 旧形是从 git blob 重建的（已注明）；每形只跑 1×，未跑 `-count` |
| **AC#2** 内容 needle 只在 stdout＋窗口化 | **成立（带具名让步：那五字只能钉审计流）** | 流那一支：`B/C/E/F/N` 在终态**全红**，红句里的 stdout 窗口**逐字带着被删薄的那一句**⇒ 读的确实是操作员那一路；`Q`（五字搬到 stdout、审计删掉）终态 FAIL 3.63s 而旧形同发 PASS 4.21s⇒ 拼接串跨流抵账已被切断。窗口那一支：`SA`（两枚审计句挪到启动、重启时不再打印）旧形 PASS 2.29s → 终态 FAIL 41.45s⇒ 窗口有牙 | 票面字面要的"把 `重启进程后生效` 钉在 stdout"**做不到**（产品那句没这五字，钉上去＝改文案＝禁区）。算不算半格＝编排者裁（§5①） |
| **AC#3** 反向正控要真响 | **成立** | 终态：`B` FAIL 3.39s（7 枚 omits）·`C` FAIL 2.78s（3）·`E` FAIL 2.85s（1）·`F` FAIL 3.79s（1）·`N`（"原因"换成一句无关话＝票面点名的新突变）FAIL 3.62s（2）·`K` FAIL 43.34s（窗口空行）·`P`（审计尾部删那五字）FAIL 2.61s·`Q` FAIL 3.63s·`SA` FAIL 41.45s·`MK` FAIL 42.71s；`MIRROR`（逐字节相同副本）PASS 2.97s＝overlay 落地正控 | 无缺口；"两发都仍绿＝装饰"那一支**未被触发** |
| **AC#4** 横幅假绿那一枚也要钉 | **部分成立**（stdout 那一路没被拦住，本腿量到了） | **拦住**：`MK`（横幅先带全字样＋操作员重启句整体删除）旧形 PASS 4.38s → 终态 **FAIL 42.71s**，机制＝窗口内确实没有该句（红句 dump 里 `上一版横幅…` 逐字可见）；`SA`（审计侧同形）旧形 PASS 2.29s → 终态 FAIL 41.45s。**没拦住**：`M`（横幅先带同样字样、**重启句照打**）终态仍 **PASS 2.62s**（旧形同发 PASS 3.26s＝两代尺都绿） | 票面"不许 **0ms** 通过"在本形状上从不出现 0ms（最快 2.29s、最慢 42.71s）⇒ 措辞怎么读归你裁（§5②）；`M` 那一形是**真缺口**，分开答＝"被拦住了"只覆盖 `SA/MK`，"根本没测到"不成立（`M` 测到了、仍绿），收口形状候选见 §5③ |
| **AC#5** 零放宽 | **成立** | 断言计数 6/2/5 ≥ 旧 5/2/4；删除列只落在被替换的 `why+out` 块；无"或满足任一即可"；`K`/`SA`/`MK` 都在 40s 期限**报红不 skip**（helper 出口是 `t.Fatalf`）；本腿具名复跑：基线三枚 0 红＋**10 枚 `TestTicket223*` 全 PASS**（`logs/AC5-all-ticket223.txt`：PASS=10／FAIL=0，含 `TestTicket223R2FailureSentenceRouting` 0.62s、`TestTicket223PermissionDeniedSitsInItsOwnSentence` 3.10s、`ok … 35.561s`） | 本腿**未跑** `-count=5` 稳定性发（预算所限）；间歇红那一枚属票 223 地界，本腿只证"没被放宽" |
| **AC#6** 整包终态 | **成立（零新增持久红；一枚新面孔待归账）** | `PATH=… go test ./cmd/wisp ./internal/... -count=1` ⇒ rc=1／8 枚红逐名（§4 G4）：`TestAC1AlwaysBranchDoesNotRevertAHandEditedKey` 安静单跑 PASS 1.19s；`TestTenOpsInOneToolCallGetOneConfirm` 安静单跑 PASS 0.07s；`TestResolvePerCallBudget` `-count=3` ⇒ `ok … 4.063s`（三遍全绿，不记账）；ball 1＋panel 4＝派单在册别人地界；`Test258*` 今日未红 | `TestTenOpsInOneToolCallGetOneConfirm` **不在派单给的在册名册里**（带载型假红）⇒ 要不要进在册归你记（§5④） |

**三枚硬点里本腿真正裁不动的一句**：AC#4 的 `M` 形——窗口下界取在 plant **之前**，"启动横幅先含同样字样"天然落在窗口内，
而用例里没有任何一条要求"这句话在窗口内只出现一次"或"必须晚于 applied 行"（本腿现读终态用例体确认，非推断）。
⇒ 本腿判 AC#4 **部分成立**，不判成立也不判失败。

---

## §3 突变名册（本腿自造，全部 `-overlay`，逐发附落地证明）

产物：`D:/tmp/wisp232v1/mut/prod/p*.go`（9 枚产品副本）·`mut/test/tOLDFORM.go`（重建旧尺）·`overlay/*.json`（24 枚）·
`logs/*.txt`（24 枚，已入库 `.scratch/wisp/probes/232/v1/logs/`）。生成器 `D:/tmp/wisp232v1/gen_mut{,2,3}.py`（只读 base，⛔ 不写工作树）。

| 发 | 内容 | 尺 | 读数 | 落地证明 |
|---|---|---|---|---|
| `oldform-only` | HEAD 产码，仅换测试文件 | 旧形 | **PASS 2.88s** | 同通道另一发（`K-oldform`）红句带出旧形独有措辞 |
| `B` | 操作员句只留 `"wisp run: 这些段的改动本次运行不会生效。\n"`，三段（原因／涉及／两件事都没发生）删光，needle 保留 | 旧形 | **PASS 3.16s**＝**瘦句不响复现** | `diff base pB`＝`:322-326` 五行→一行 |
| `B` | 同上 | 终态 | **FAIL 3.39s**：`the operator's restart sentence omits "app.autostart"`×共 7 枚，窗口逐字 `wisp run: 这些段的改动本次运行不会生效。` | 红句带出新字节 |
| `C` | 只把"原因：…交给平台层（开机自启注册、单实例锁、界面语言）"换成`这几枚键在进程启动时读一次。` | 旧形 | **PASS 3.08s** | 终态红句窗口逐字带出该替换句 |
| `C` | 同上 | 终态 | **FAIL 2.78s**（omits 开机自启／原因：／交给平台层） | 同上 |
| `E` | 只删 `涉及：%s` 一支（并去掉第二个 Fprintf 实参） | 旧形 | **PASS 3.07s** | 终态红句窗口里无 `app.autostart` |
| `E` | 同上 | 终态 | **FAIL 2.85s**（omits `app.autostart`） | 同上 |
| `F` | 只把"注意两件事都没发生…"换成`文件与内存都按 D36 处理。` | 旧形 | **PASS 3.67s** | 终态红句逐字带出新句 |
| `F` | 同上 | 终态 | **FAIL 3.79s**（omits `没有被丢掉`） | 同上 |
| `N` | **票面 AC#3 点名的"换成一句无关话"**：`原因：…` → `原因：这条提示与配置无关，只是占位的一句话。` | 终态 | **FAIL 3.62s**（omits 开机自启／交给平台层） | 红句窗口逐字带出占位句 |
| `K` | 操作员整句 `fmt.Fprintf(rt.stdout, …)` 删除（审计两句不动） | 旧形 | **FAIL 43.17s** `stdout never carried "本次运行不会生效" within 40s` | 红句＝旧 helper 措辞（无 `AFTER the plant`） |
| `K` | 同上 | 终态 | **FAIL 43.34s** `stdout never carried "本次运行不会生效" **AFTER the plant** within 40s; window since the mark:`（窗口为空行） | 窗口空＝真没打印，非拼接串兜住 |
| `P` | 只删**审计 detail** 尾部 `，重启进程后生效`（操作员句一字不动） | 终态 | **FAIL 2.61s** at `:611`；红句里的 **stdout 窗口逐字没有那五字** | ★不对称由这一发独立证实 |
| `P` | 同上 | 旧形 | FAIL 2.71s at `:493` `the restart sentence omits "重启进程后生效"; audit:` | 旧形独有文案＝旧尺已落地 |
| `Q` | 把那五字**从审计句搬到操作员句**（跨流换位，审计删） | 旧形 | **PASS 4.21s**＝**拼接串跨流抵账复现** | `diff`＝两处一行级 |
| `Q` | 同上 | 终态 | **FAIL 3.63s** at `:611` | ★同一句字面搬到另一路，终态仍红＝具名流真的分了家 |
| `SA` | 两枚重启**审计**句挪到启动横幅处（plant 之前打印），`reportRestartPending` 内删除；操作员句不动 | 旧形 | **PASS 2.29s**＝**前向假绿复现（审计侧）** | 终态红句 `full stderr` 里逐字可见这两行 |
| `SA` | 同上 | 终态 | **FAIL 41.45s** at `:569` `the audit trail never carried "…restart-pending…" AFTER the plant`；窗口只含 `state=applied … restart=[app]` | ★窗口化的牙量在审计流 |
| `M` | 启动横幅**追加**一句同样字样的操作员句（重启句照打） | 旧形 | **PASS 3.26s**（非 0ms） | 红/绿 dump 均带横幅文本 |
| `M` | 同上 | 终态 | **PASS 2.62s**＝**AC#4 残余缺口**（本腿量到的"没拦住"） | 同上；本腿没说成"没测到"——它就是同一发在修好的尺上仍绿 |
| `MK` | `M` 的横幅 ＋ 操作员重启句整体删除 | 旧形 | **PASS 4.38s**＝票面 现量那格的本腿复算 | `AC4b/AC4d` 日志 full stdout 里 `上一版横幅` 逐字在 |
| `MK` | 同上 | 终态 | **FAIL 42.71s** at `:592` | ★复合形被拦住（机制＝窗口内无该句＋期限红，非变慢） |
| `MIRROR` | **逐字节相同副本**（`diff` 空、md5 同 `5ce441ca…`） | 终态 | **PASS 2.97s** | ★把"走了 overlay 这个动作"本身洗清的正控 |
| `X-count` | 读盘计数（非 go test） | — | 旧 5/2/4 → 终 6/2/5 | AC#5"只增不减"独立尺 |
| `R`（**未跑**） | 计划：needle 换成无关句量 helper 自身期限 | — | **未跑**：与 `K` 同族同出口，`K` 已量到 43.34s 期限红 | 具名声明：本腿以 `K` 代替，不当第二枚凭据 |

合计 **23 发本腿自造的 `go test` 读数**（终态尺 12 发＝2 绿 10 红；旧形尺 11 发＝9 绿 2 红），
红因**逐字全部抄在 `logs/*.txt`**；**零枚"改了读路径却没测到"**（每一发都配了落地证明那一列）。
`go test` 之外的门禁发（基线、10 枚 `TestTicket223*`、整包、三枚安静复量）另计，见 §4。

**读盘型断言检查（另一枚腿的教训那一支）**：本票三枚用例的内容 needle **全部读内存流**（`r.h.out`／`r.h.err`）；
`os.ReadFile` 在本文件只出现在 `:78`／`:482`，读的是**测试自己 tempdir 里的 config.toml**（验种植是否落地），
`restart_tier_keys_255r2_test.go:99` 读的是它自己的 `t.TempDir()` ——**没有一枚"读产品源文件找字"型断言**，
故"`-overlay` 对读盘断言结构性不可见"那一支在本票**不适用**（本腿 grep 复量，非推断）。

---

## §4 门禁（本腿亲手跑，逐字读数）

**G1 · `sh scripts/d22scan.sh`**（留盘 `logs/G1-d22scan.txt`）
```
runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0
d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=228, bans #1-5 cmd/=38, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=514, ban #8 cmd/=104; ban #8 emoji coverage: … cmd/ 104 Go files, comments and _test.go included
rc=0
```
⇒ **rc=0 且正控先绿**（`PASS=35 FAIL=0 SKIP=0`，不是"仪器是瞎的"）。

**G2 · 格式尺（只喂 `.go`，⛔ 一枚 `.md` 都没喂）**
- 本票地界（`cmd/wisp/config_reload_223_test.go`＋`cmd/wisp/config_reload.go`）：`gofmt -l` **空**、`gofumpt.exe -l` **空**（rc=0）⇒ **本票两枚文件干净**；tracked 尺输出里也**不含**这两枚（本腿 grep，rc=1）。
- `gofmt -l cmd internal` 与 `gofumpt -l cmd internal` 另报同一 5 枚：`cmd/wisp/models.go`、`internal/agent/approval/pending_read.go`、`internal/agent/tools.go`、`internal/risk/provenance.go`、`internal/tools/bridge.go`。
  本腿量的性质：`git ls-files --eol` 逐枚为 `i/lf w/crlf attr/text eol=lf`；python 逐字节计数例 `provenance.go CRLF=1115 LF=1115`；`git diff --quiet HEAD --` 对 5 枚逐枚通过（＝与本腿无关）；`.gitattributes` 写 `*.go text eol=lf`。
  ⇒ 判定＝**本机 checkout 行尾噪声，非本票缺陷、非产品缺陷**；本腿**未触碰这五枚**。
- CI 那一把的形 `sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only` ⇒ **rc=2**，原文：
```
== (A) tracked tree: git ls-files -z '*.go' | xargs -0 …/gofumpt.exe -l   (tracked .go files handed to it: 935)
attrib.sh: (A) gofumpt exited 123 on the tracked set (>=2 means a file would not parse) - the tracked tree is not readable by this ruler
```
  本腿按文件拆开复现，唯一解析失败＝`.scratch/wisp/probes/185/c1/mut/fs_broken.go:4:1: imports must appear before other declarations`
  （一枚 **tracked 的故意坏样本**：`git ls-tree HEAD` 在案、md5 `f5c325e04a9a0d94421a665d9479e7b5`）；分批（`xargs -n40`）重跑后报的就是上面那 5 枚 CRLF。
  ⇒ 这条 rc=2 是**仪器在 bench 树上的既有形状**（早于本票、与本票无关），要不要修／记账归你（§5⑤）。

**G3 · `go vet ./cmd/wisp/`** ⇒ **rc=0**（零输出）。

**G4 · 终态整包**（`logs/G5-fullpackage.txt`，932 行）`PATH=… go test ./cmd/wisp ./internal/... -count=1` ⇒ **rc=1**，逐名红名册（8 枚）＋安静复量：
```
--- FAIL: TestAC1AlwaysBranchDoesNotRevertAHandEditedKey (17.71s)      cmd/wisp   安静单跑 PASS (1.19s)  ⇒ 带载型
--- FAIL: TestTenOpsInOneToolCallGetOneConfirm (5.24s)                 agent/approval 安静单跑 PASS (0.07s) ⇒ 带载型，★派单名册里没有这一枚
--- FAIL: TestC21TableColourRowsMatchTokensCSS (0.00s)                 ball        ⇒ 派单在册（别人地界）
--- FAIL: TestApprovalCardViewJSONKeysMatchFrontendTypes (0.02s)        panel       ⇒ 派单在册
--- FAIL: TestComposerContractTypesMatchFrontend (0.01s)                panel       ⇒ 派单在册
--- FAIL: TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme (0.22s)     panel       ⇒ 派单在册
--- FAIL: TestC21DesignTokensFourWayAgree (0.00s)                       panel       ⇒ 派单在册
--- FAIL: TestResolvePerCallBudget (3.52s)                              risk        ⇒ 安静复量 `go test ./internal/risk/ -count=3 -run TestResolvePerCallBudget` 原文 `ok github.com/CarlosShao/wisp/internal/risk 4.063s`（三遍全绿 ⇒ 不记账，两遍读数都在案）
```
`Test258Occupied*`／`Test258V1ProbeSummon*`（`cmd/wisp/resident_hotkey_258_windows_test.go`）**今日整包未红**（本腿的包内 `TestTicket223*` 整跑亦未见其红）⇒ 无需复现"热键被占"那枚归因。
⇒ **零新增持久红**；唯一待归账＝`TestTenOpsInOneToolCallGetOneConfirm`（§5④）。

**G5 · 终态写面与还原对拉**
```
git status --porcelain -- cmd internal   ⇒ 0 行（与起手相同）
md5sum cmd/wisp/config_reload.go         ⇒ 5ce441ca5e72b64d18a6c26f1c066882（＝起手＝HEAD）
```

---

## §5 判不动的地方／要编排者裁（本腿一律未自行假设）

1. **AC#2 的字面那一半**：现量（§0 末＋§3 `P-final` 两处独立读数）＝`重启进程后生效` 只在审计句，操作员句真源没这五字；钉去 stdout 必须改产品文案＝本票禁区。
   `r2` 的处置（钉具名审计流＋种后窗口，另在 stdout 窗口钉 `需要重启进程`）本腿复核**有牙**（`P` 2.61s 红、`Q` 3.63s 红而同发旧形 4.21s 绿）。
   **要裁**：算"AC#2 成立（带具名让步）"还是"半格"。若你要"操作员那一路也必须有那五字"，那是**产品文案改动**，得另立票。
2. **AC#4 票面"不许 0ms 通过"怎么读**：本腿所有读数里**没有一枚 0ms**（最快 2.29s／最慢 42.71s）。旧形在 `M/SA/MK` 上是**真绿**，终态在 `SA/MK` 上是红。
   按"不许假绿"读＝两形满足；按字面"不许 0ms"读＝票面现量与盘上不符。**要裁**（本件按"不许假绿"写判语并逐字给秒数）。
3. **★本腿量到的真缺口（建议立新票，属判据改动故本腿不实现）**：`M` 形＝启动横幅先带同样字样、之后种改动、重启句照打——**终态仍 PASS 2.62s**。
   根因：窗口下界取在 plant 之前，"先含同样字样"的横幅天然落在窗口内；用例里没有任何一条要求"窗口内该句出现次数＝1"或"必须晚于 `state=applied` 行"。
   现网产码里该句唯一来源仍是 `config_reload.go`（本腿 grep 全仓：无其它产码来源）⇒ **今天不红，是前向形状**。
   收口候选：(a) `strings.Count(win, 该句)==1`；(b) 标记取在 `state=armed` 横幅**之后**；(c) 要求该句偏移晚于 `HOT-RELOAD state=applied` 那一行。
4. **`TestTenOpsInOneToolCallGetOneConfirm`（`internal/agent/approval`）**：整包红 5.24s、安静单跑 PASS 0.07s，**不在派单给的历史在册名册里**（在册只有 ball 1＋panel 4＋risk 争用＋r2 具名的 `TestAC1ResidentLeg…`）。
   本腿判＝带载型假红（与 `TestAC1*` 同族）；**A##／在册归你写**，本腿 `docs/**` 一字未动。
5. **格式仪器在 bench 树上的形状**：tracked 尺 rc=2 由一枚 tracked 故意坏样本（`probes/185/c1/mut/fs_broken.go`）触发 `gofumpt` 退出 123；
   `attrib.sh` 的退出规则表写了"非空要可归因"，**没写"解析失败"那一支**。另有 5 枚工作树 CRLF 命中（本机 checkout 现象）。
   ⇒ 要不要立票收口归你；本腿既没动那五枚也没动 `probes/185/**`。
6. **本腿没跑的**：`-count=5` 稳定性发（AC#5 的加强版，`r2` 自述跑过、本腿不采信亦不复用）；`-tags winlive`（禁区，未碰）；`go test ./...` 全仓（派单只要求 `./cmd/wisp ./internal/...`）。
