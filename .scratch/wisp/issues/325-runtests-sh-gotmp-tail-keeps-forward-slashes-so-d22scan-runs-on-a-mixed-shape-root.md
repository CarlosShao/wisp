# 325 — `tools/d22scan/runtests.sh` 的测试临时根＝**混形活坑**（基座反斜杠＋尾接正斜杠）：同一族分隔符病在票 317 射程外的那一处

- **Status**：已立（2026-10-16 09:2x，编排者），**派 `325-r1`**（写面＝`tools/d22scan/**`，起手时该面无写腿）
- **Type**：测量完整性缺陷（与票 317 同族、同一种写法），⛔ 不是功能缺陷、⛔ 不是为了变绿
- **来源**：票 317 的非实现者验收腿 `317-v1`（表＝`docs/evidence/s1/317-gotmp-separator.md`，笔 `ceaa50ac`）§4 具名上报；编排者 09:2x 现读复核并附自己的成对复跑（`.scratch/wisp/probes/orch/20261016/logs/10-nat.txt`／`10-psx.txt`）。

## 这是什么（人话）

跑测试时，Go 要把临时文件放在哪，是由一个环境变量（`GOTMPDIR`）说定的。我们这台 Windows 上，这个值**必须写成带反斜杠的样子**（`C:\Users\…`）；一旦里面混进正斜杠（`C:/Users/…`），有几条专门测"目录链接（junction）"的用例就会**假红**——不是代码坏了，是量它的尺自己把路径写坏了（`mklink` 会把 `C:/…` 当成命令行开关，报 `Invalid switch`）。

票 317 已经修好了三处这种写法（都在 `scripts/` 里）。**同一族写法还剩一处没修**，就在 `tools/d22scan/runtests.sh:92`：它的开头用了反斜杠、结尾接名字时用了正斜杠，所以生成的路径是"半对半错"的混形。今天没炸，只因为这条脚本的射程里没有那几条会炸的用例——**这是"坑还在、只是暂时没人踩"，不是"没坑"**。

## 现量（编排者 2026-10-16 09:2x，全部自己跑）

1. 字节形（`sed -n '92p' tools/d22scan/runtests.sh`）＝`gotmp_run="${WISP_TEST_GOTMP_ROOT:-"$LOCALAPPDATA\\wisp-test-gotmp"}/d22scan-$(date +%Y%m%d-%H%M%S)-$$"` ⇒ 基座 `\\`（原生）＋尾接 `/`（正斜杠）＝**混形**。
2. 已修的三处对照（同尺）＝`scripts/portable-tests.sh:155-156`／`scripts/winsec-tests.sh:69`／`scripts/wisp-cli-tests.sh:82` 三处**尾接都是 `\\`**＝原生形。
3. 病本身可复现（我这把尺，仓外导出树、同一个目录只差分隔符）：`GOTMPDIR='C:\tmp\wisp317orch\pair'` ⇒ `RUN=4 PASS=4 FAIL=0` rc=0 ⇄ `GOTMPDIR='C:/tmp/wisp317orch/pair'` ⇒ `RUN=4 PASS=0 FAIL=4` rc=1，红句逐字同行同号（`internal/risk/syncdirs_ancestor_actable_leg_116_test.go:182`／`:286`＋`…reparse_ticket105_windows_test.go:59`／`:116` 的 `Invalid switch - "tmp".`）。

## 要建什么

- [ ] **AC#1 改形**：把 `tools/d22scan/runtests.sh:92` 的尾接 `/d22scan-` 改成 `\\d22scan-`（与 `scripts/**` 那三处逐字同形），旁边留一条说明"为什么这里不许折成正斜杠"的注释（⛔ 只写"风格"、要写**会炸的那条用例**）。判据＝改前/改后各跑一发 `sh scripts/d22scan.sh`（干净台面：`GOTMPDIR`／`WISP_TEST_GOTMP_ROOT` 抽干），两发都 rc=0 clean，并把那一行 `GOTMPDIR=` 的逐字形状落在读数件里。⚠ **⛔ 顺手把 `wisp_311_ref_path` 那四处对照位（`portable-tests.sh:818`／`winsec-tests.sh:174`／`wisp-cli-tests.sh:186`／`runtests.sh:155`）一起"改一致"**——票 317 已逐处判过它们是对照位显示值、⛔ 流向 `GOTMPDIR`；动了＝越界，要具名上报。
- [ ] **AC#2 成对反控（这族尺的牙）**：在**仓外导出树**上把这一行退回混形（物理改字节，⛔ `-overlay`——这族尺运行时读盘），跑一次能让 `./internal/risk/` 那四枚用例红的构造 ⇒ 证明"混形会炸"今天仍然真；然后还原、复跑绿、记录还原哈希。⛔ 只做"改前必红"就交，还要写清"射程里今天有没有那四枚"——如果 `d22scan` 那档根本不调 `./internal/risk/`，就具名写成**缺世界**（坑活着但今天看不见），⛔ 假装测到了。
- [ ] **AC#3 顺带交出票 317 AC#3 那一发**（同一台面、同一档，省一枚腿）：修完后跑一发 `bash scripts/portable-tests.sh --scope=windows`，三样齐＝① `GOTMPDIR=` 行逐字原生形 ② 八包四数**逐包**落盘 ③ 红名册与 10-14 那两发 CI 的 `Invalid switch` 名册作差＝新增 0。
- [ ] **AC#4 门禁**：`bash -n tools/d22scan/runtests.sh`／`sh scripts/d22scan.sh`／`sh scripts/check-path-length-budget.sh --with-self-test`／`go build ./...` 逐枚自落 `rc=`。⚠ 两枚台面雷：d22scan 会把 `GOTMPDIR` 指到非 `%LOCALAPPDATA%` 卷时回 **rc=2 的票 311 GUARD 射程拒绝**（⛔ 当红）；`cmd/wisp` 测试缺 `third_party/sherpa-onnx` 那三枚 DLL 会回 `0xc0000135`＋`RUN=0`（＝台面故障，⛔ 读数）。
- [ ] **AC#5 账**：票 317 面上那句"同形、今天未爆"要**重述**成"混形，射程遮蔽中的活坑"——由编排者在 317 的 Progress log 追加，⛔ 本腿改 317 原文。

## 禁区

⛔ 动 `scripts/**`（那三处已由票 317 修毕并验过）／`wisp_311_ref_path` 四处对照位／`.github/**`／`ci.yml`／`thresholds.go`／golden／`tools/d22scan/allowlist.txt`／三枚冻结件／`frontend/**`／`design/**`／`PLAN.md`／`specs/**`；⛔ 为变绿放宽任何断言或删用例或加 `t.Skip`；⛔ push（只 commit，显式 pathspec）；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删。

## 派单前置条件

无（`tools/**` 面空）。写面三处＝`.scratch/wisp/probes/325/r1/`＋`docs/evidence/s1/325-runtests-gotmp.md`＋票 325 末尾 1 条 bullet；⛔ 勾 AC 框（翻框归编排者）。

## Progress log (append-only, newest last)

- [2026-10-16 09:2x +08] 编排者立票（来路＝`317-v1` 具名上报＋我自己现读三把尺：字节形／已修三处对照／成对复跑 4 红⇄4 绿）。★出身笔**我已在立票时反查、不必等腿**：`git log -1 581dd252`＝10-14 13:58「编排者：微修 runtests.sh 第四处折法（镜像317-r1形）」＝尾接 `/d22scan-` 的来路；`git log -1 76b9e214`＝10-14 14:49「编排者：253-r6 收件＋自伤修复（runtests.sh 单反斜杠→双…）」＝基座那半的来路（`git log -L 92,92:tools/d22scan/runtests.sh` 现读命中它）。⇒ **两枚都是编排者自己的手**，本票的账记在我头上、⛔ 记在任何腿上。
- [2026-10-16 09:5x +08] agent=`325-r1`（落地腿·测量完整性；⛔ 翻框 ⛔ 动 `Status:`/票名/`-done` ⛔ 改票 317 一字 ⛔ 动 `scripts/**` ⛔ 动四处对照位 ⛔ push）did=★**产码一枚**＝`tools/d22scan/runtests.sh:92` 尾接 `/d22scan-` → `\\d22scan-`（笔 `428ef21e`，与 `portable-tests.sh:156`／`winsec-tests.sh:69`／`wisp-cli-tests.sh:82` 现量同形：`}\\d22scan-$(d` ⇄ `}\\winsec-$(da` ⇄ `}\\cli-$(date`），旁边注释 23 行点名**会炸的那族用例**（⛔ 只写风格）。★**AC#1 判据两发**（台面＝仓外导出树，pre＝锚 `46007856`／post＝`428ef21e`，六件套逐枚验在＋`find -name '*.go'|wc -l`＝1090 非空；共享树有 `322-r1b` 在写 winsec ⇒ 整树读数一律没在仓里跑）：改前 rc=0 clean＋发射根逐字 `…\wisp-test-gotmp/d22scan-20261016-093651-329`、改后 rc=0 clean＋`…\wisp-test-gotmp\d22scan-20261016-093745-842`，**两发正控制四数逐字同形** `PASS=35 FAIL=0 SKIP=0 === RUN=77` ⇒ **这一格在未修码上"不响"＝门绿着、坑没人看见**。★**AC#2 真红（不是构造的正斜杠，是那一行自己生成的根）**：四枚名先 `-list` 逐枚 1/1 自证，pre 树 `sh tools/d22scan/runtests.sh ./internal/risk/ -run '^(四名)$'` ⇒ rc=**1**、`RUN=4 PASS=0 FAIL=4 SKIP=0`，四行＝`syncdirs_ancestor_actable_leg_116_test.go:182`／`:286`（前置条件文案里同一条根两侧分隔符分家）＋`syncdirs_ancestor_reparse_ticket105_windows_test.go:59`／`:116`（真 `mklink /J` 被 cmd 拒）；⚠ **措辞要说准**＝混形（尾接那枚 `/`）报 `Parameter format not correct - "d22scan-…"`，⛔ 与编排者/`317-v1` 那对**纯正斜杠根**的 `Invalid switch - "tmp"/"Users"` 混为一句，两形同因（开关误读）。★**反控链**＝在导出树副本上**物理**退回混形（⛔ `-overlay`）：突变自证 `mixed-shape-hits=1／native-shape-hits=0` ⇒ 四红 rc=1 ⇒ 用 `git show 428ef21e:` 的字节还原 ⇒ **sha256 回到提交态 `ad6f6a8d72cb…4bf6`**、`mixed=0／native=1` ⇒ 复跑 rc=0 `PASS=4 FAIL=0`。★**射程答＝缺世界（具名，⛔ 假装测到）**：本入口唯一的门 `scripts/d22scan.sh:51` 只交 `tools/d22scan ./...`，`grep -rn mklink tools/`＝**0**（全仓 mklink 只在 `internal/` 七枚）⇒ `./internal/risk/` 今天不在这档分母里；CI 唯一直呼本 runner 打 risk 的那步（`ci.yml:831`）`-run` 只盖 `TestPathResolverJunctionWindows`，本腿在混形码上亲测那发＝rc=**0** `PASS=1`（它的 link/target 走 `filepath.Join` 归一化，吃得下混形）⇒ 它绿≠坑没；⛔ 为凑红把档位改成跑 risk（那是改门、要人工批准）。★**AC#3（顺带交票 317 那格）**＝post 树 `bash scripts/portable-tests.sh --scope=windows` rc=**0**：① 发射行逐字原生 `GOTMPDIR=C:\Users\swq\AppData\Local\wisp-test-gotmp\portable-20261016-094356-1108`；② **逐包四数**（口径 `RUN`＝`^=== RUN`／`PASS`＝`^--- PASS` 只数顶层，本发两把差 200）＝proc 48/40、secret 29/21、config 185/104、risk 200/132、ball 75/69、perm 14/14、plugin 7/7、llmrecord 2/2、session 36/19、projctx 10/10、audio 56/44，总计 `=== RUN=662 --- PASS=462 --- FAIL=0 --- SKIP=0`（与脚本自打印行一致；risk 那行与 `317-v1` 同形＝跨腿复认）；③ 红名册＝**空** ⇒ 与 10-14 两发 CI（`38508311376` vs `38484004521`，账＝A922）的 `Invalid switch` 名册（5 增里那四枚 Sync/116 族；第五枚 `TestSessionListByWorkspace` 文案不同、另账）作差＝**新增 0**；⛔ `--scope=cli`／⛔ `./cmd/wisp` 整包（票 314 死亡区，派单禁）。★**AC#4 门禁四枚全 rc=0**：`bash -n`（post 树）／`sh scripts/d22scan.sh`（post 树）／长度门 `--with-self-test` `VERDICT GREEN`（★该门要 git 索引、`git archive` 不带 ⇒ **唯一在仓内只读跑的一发**，台面值已具名）／`GOFLAGS= go build ./...`（0 字节输出＝无诊断）；`gofumpt`＝`which` rc=**1** 且四处候选目录（`~/go/bin` 空／`~/go/pkg/mod`／`$LOCALAPPDATA/go/bin`／`/d/work/base/go/bin`）零命中真身，本腿写面 **0 枚 `.go`** ⇒ 射程 0，据实写"不适用"、⛔ 拿 `gofmt` 的 0 冒充；票 311 的 rc=2 射程拒绝未触发（根始终在 `%LOCALAPPDATA%` 同卷 C:），五类台面孔（`walk`/`build failed`/`no tests to run`/`0xc0000135`/`RUN=0`）一次未现。★**两处更正（记我=编排者的派单，不记腿）**：① 票面/派单那句"八包"＝**过期计数**，现量 `windows` scope＝**11 包**（`portable-tests.sh:329-336`）；② 四处对照位里 `runtests.sh` 那一处因本笔注释**行号位移 `:155`→`:178`**（字节未变），那把行号尺请按此重锚。件＝证据 `docs/evidence/s1/325-runtests-gotmp.md`（179 行，笔 `c7190b27`）＋读数 `probes/325/r1/`（22 枚，笔 `73131d29`）。**没做完（具名）**＝CI `test-windows` 换色那一发（推送权在编排者，本腿只有本机导出树读数）／`--scope=cli` 那档（禁跑，非遗漏）／纯正斜杠基座那发本腿未重跑（只有编排者 `10-nat/10-psx` 与 `317-v1` M1 两份）。next＝非实现者验收腿逐格判语（★重点攻：AC#2 的"缺世界"是否被两发读数钉死、AC#3 的 11 包口径与作差基准、AC#1 那 23 行注释有没有越界解释别票）→ 编排者重述票 317 那句＋翻框。
- [2026-10-16 09:5x +08] 编排者★更正派单与本票那句"八包"＝现量 **11 包**（尺＝取 `scripts/portable-tests.sh` 第 329-336 行那段 `windows` scope 数组逐枚数：proc／secret／config／risk／ball／perm／plugin／cmd 下 llmrecord／session／projctx／audio）。腿 `325-r1` 顶回、我复跑同意 ⇒ 原句⛔ 抹平、只加这条更正；AC#3 判据按现量读成"十一包逐包四数"。★同时记下我自己刚才造的一次事故：我用双引号串把这段更正正文直接喂给 shell，串里的反引号**真的被执行了**（bash 报了 `windows: command not found` 之类三条）——这条规矩是我自己立的（第 109 条），第三次踩；后果核实＝只往 `/dev/null` 打、票 325 那次 append 没落笔、票 317 那条走的是先 Write 成文件再 `cat`，字节完好（反引号 34 枚成对、"八包"在、无 `command not found` 残迹）。
