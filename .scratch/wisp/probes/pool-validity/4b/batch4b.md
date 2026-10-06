# pool-validity-4b / 批4b —— 零勾开放票有效性普查（号段 ≥238，现量 8 枚）

> 只读普查腿 `pool-validity-4b`。任务＝把"238 以上号段里一格都没勾的票"逐枚判档，供编排者回答"还剩多少真活"。
> 起手 HEAD `3d9b8374`（2026-10-06 13:12:18 +0800）。⚠ 换 HEAD 要重量。
> ★交件时刻 HEAD＝`6f0a2d83`（10-06 14:01:55 +0800）——起手到交件之间漂进 8 枚 commit（本腿自落 `8aae805e`／`41d12ca7`＋同职腿 4c 四枚＋写腿 236-r2 一枚＋台账 `7aeca858`），
> 所以 §2 表里所有 `HEAD` 都是"我这发命令当时的 HEAD"。⇒ 本腿在落这枚 commit 之前**把八枚的判据尺在 `6f0a2d83` 上逐枚复跑了一遍**：
> 238＝11 枚 importer／proc·observe·buildinfo 各 0（rc=1）；242＝`:109` 逐字仍在＋capability 尺 6 枚在树；244＝`build.ps1:115` 与 `console_other.go:7` 两行原样；
> 247＝`b.SetAudioLevel` 产码命中仍只有 `cmd/balldebug/main.go:418/:477`；253＝`inbound_roster_253_test.go` 五枚尺仍在；262＝脚本＋`ci.yml:166` 仍在；
> 264＝agent/tools 非测试 `redact|Redact` 尺 rc_all＝1、rc_nontest＝1（**两把都零命中**，不是过滤把命中洗掉了）；266＝`'"scripts'`＝0 与钉 `:233/:291` 原样。
> **八枚档位一枚没因这次重量而改。**唯一受 HEAD 漂移影响的读数已具名写在 §4 第 1 处（跟踪路径枚数 6357→6358），八枚判定无一枚依赖那 1 枚差。
> 本腿**零 go 命令**（`go test`／`build`／`vet`／`run`／`list`／`env` 一次都没跑过，一次都不跑）：
> 此刻唯一被允许跑 go 的是写腿 `236-r2`（在 `internal/tools` 取突变读数）。
> 枚数／行数一律 `git ls-tree`／`git ls-files`／`wc`。
> 工具全集＝`git log`／`git grep HEAD`／`git show HEAD:<path>`／`git ls-tree`／`git ls-files`／`ls`／`grep`／`wc`／Read。
> 票面零改动：零翻勾、零改名、零 `Status:` 改写、零产码改动、零台账（`docs/**`）写入。
> 唯一写入路径＝`.scratch/wisp/probes/pool-validity/4b/**`（新建；起手 `ls` 过 `4b` 不存在）。
> ⛔ 零读零引 `frontend/**`、`design/**`、`.gitignore`、`.scratch/wisp/probes/236/**`、`.scratch/wisp/probes/pool-validity/4c/**`。
> ⛔ `cmd/wisp/**`、`internal/**` 工作树内容面不读（写腿地界），要看一律 `git show HEAD:<path>`。
> ⛔ `docs/evidence/s1/**` 只按**文件名**取、不读内容。
> 前人件 `pool-validity/1/**`、`pool-validity/2/batch3.md`·`batch4.md`（死腿骨架）一字不动；
> `pool-validity/2`（46 枚）、`3a`（13 枚，100–149）、`3b`（18 枚，150–199）、`4a`（11 枚，200–238）
> 只读只学**表形**，⛔ 不照抄其口径、不采信其判档。
> ⚠ 各批是**不同日期、不同对象**的量，本件不写成"上一版 vs 新版"，也不写"谁推翻了谁"。

## §0 四档尺（本腿判法，口径自立）

| 档 | 判据 | 复法（必须把命令原文＋读数留在表里） |
|---|---|---|
| **仍成立** | 票面点名的那枚缺陷今天还在树上 | 把票面 §现场／AC 里**最硬的一条**自己跑尺：`git --no-pager grep -n '<那句字面或符号>' HEAD -- internal/ cmd/ scripts/ .github/` 的命中数＋file:line |
| **已失效／已被别人做掉** | 缺陷不在了 | ⛔ 硬门：一枚 commit 号＋`git log --name-only` 命中，**或**今树 0 命中的复算尺；拿不出就不许写这一档 |
| **差翻勾** | 缺陷不在了，但凭据是"票内注记／台账 `A##`／非实现者裁决表" | 凭据种类写进复法列；⚠ **不等于本腿裁它可结案**，翻勾归编排者 |
| **量不到** | 判据落运行期行为／真开窗／真机双击／DPAPI／多账号／CI 侧读数／`frontend/**`·`design/**`（禁读） | 归口写清**缺哪一行读数**；⛔ 不许由 grep 命中外推成"这功能能用" |

★ 本段三枚已知**有特殊状态**，判之前先核：
- `244`：有一格（AC#5）是编排者裁过的"下游消费者名册"欠账。
- `264`：**待人拍板**的功能级票（owner 未答＝默认不做）⇒ 判"量不到／待人裁"而不是"没人做"。
- `253`：涉及 `internal/panel` 的**冻结件**射程。
  这三枚的判定各带一行"它卡在哪一层（欠读数／等人拍／契约邻接）"。

## §1 名册（8 枚，本腿现量尺）

分母尺（本腿 10-06 13:1x 实跑，逐字）：

```
for f in .scratch/wisp/issues/2*.md; do case "$f" in *-done.md) continue;; esac;
  b=$(grep -c '^- \[ \]' "$f"); x=$(grep -c '^- \[x\]' "$f");
  if [ "$b" -gt 0 ] && [ "$x" -eq 0 ]; then basename "$f"; fi; done
```

全量输出（`2*` 号段，未按 ≥238 过滤，含 <238 的 15–26 段与 200–236 段，那些**不归本腿**）共 31 行：
`200 21 22 220 225 227 229 23 230 232 233 234 236 237 238 24 242 244 247 249 25 253 26 262 264 266 269 27 28 29`。
其中 `249`、`269` 已撤销／并案（票面 `Status:` 行已就地标注），⛔ 本腿不判；
`232`、`236` 此刻有写腿在做，⛔ 本腿不判（判会互相洗读数）。
⇒ 本腿分母＝**号段 ≥238 且非 249/269 的零勾票**：

| 票号 | 票文件（全路径 `.scratch/wisp/issues/` 下） | 票面标题（原文，已剥加粗标记） | 未勾格数 | 上次动过（`git log -1 --format='%h %ad' --date=format:'%m-%d' -- <票文件>`） |
|---|---|---|---|---|
| 238 | `238-first-run-on-a-fresh-machine-creates-the-whole-data-root-with-no-seal-to-inherit.md` | 238 — 新机器上先跑常驻（GUI）那条腿时，整个数据根是"没有任何封条可继承"造出来的（`internal/proc` 对 `winsec` 导入数＝0，连 `MkdirAll` 的父档都不是窄的） | 4 | `cf03fa31 09-30` |
| 242 | `242-the-grant-binding-layer-and-the-panel-facing-read-surface-both-have-zero-rulers.md` | 242 — 一次性令牌的绑定那一层全仓零尺，而"面板那面不带令牌"这句话只由三行注释守着：两形我都有读数，都不是 AC#5 的洞，但下一枚腿照它们写就会签错字 | 3 | `745173ad 10-03` |
| 244 | `244-spec-11-wants-the-gui-subsystem-build.md` | 244 — SPEC-11 要求"无参数＝GUI 子系统"，`docs/BUILD.md` 把这一切换推迟给票 07、票 07 结案时没把它带走：今天 `build/wisp.exe` 的 PE 子系统我读到的是 CUI(3)＝双击连一个黑控制台窗口一起起，而两行注释反过来声称"已经链接成 GUI 子系统" | 7 | `378e9ca6 10-05` |
| 247 | `247-the-capture-stack-has-zero-importers-so-no-real-microphone-level-ever-reaches-the-ball.md` | 247 — 采集栈写完了但全仓零 importer：真麦克风的电平到今天没有任何一处会去读，球上那个数是命令行手填的 | 10 | `8d17e8cf 09-30` |
| 253 | `253-panel-inbound-three-ruler-holes.md` | 票 253 — 面板入向那三枚尺洞：可达性钉只认 `postMessage(`、`config.get`/`config.set` 不带 `panel.` 前缀掉在两把名册尺之外、旧封套被宿主库静默 `resolve(null)` 且 Go 侧零审计 | 4 | `46079fcc 10-04` |
| 262 | `262-tracked-path-length-gate-for-ci-checkout.md` | 票 262 — 全仓没有任何一把尺盯"跟踪路径长度"，而过长路径会让本机 runner 的 `actions/checkout` 整步失败，`slo-full`（D32 唯一求值路径）已被两枚长名工单打死过一整天 | 8 | `da929f2b 10-04` |
| 264 | `264-error-text-carries-absolute-paths-into-model-visible-text.md` | 票 264 — `fs.*` 与 `fs_write` 的失败句今天把你电脑上的绝对路径逐字发给云端服务商，而这一侧一枚仪器都没有：要不要在发出去之前抹掉路径（⛔ 取舍归机主，本格不派写腿） | 4 | `275864f1 10-04` |
| 266 | `266-no-external-tooth-over-the-slo-check-self-locking-nails.md` | 票 266 — `scripts/*.ps1` 里的门自锁钉可以被就地删掉而本机零枚外牙会发现：唯一兜得住的是 `slo-freshness.sh` 的 P3 aging，而那枚牙按"天"钝（票 263 的两枚钉实测删掉即 rc=0 全绿） | 5 | `a32f30fb 10-04` |

**与编排者名册的对账**：编排者 10-06 13:1x 现量给出 `238／242／244／247／253／262／264／266` 共 8 枚；
本腿 10-06 13:17 于 HEAD `3d9b8374` 复量的 ≥238 零勾票是**同一 8 枚，一枚不差**（排除 249/269 后）。无差异要具名说明。

## §2 判档表（8 枚 × 4 列）

| 号 | 档 | 复法（命令原文） | 读数 |
|---|---|---|---|
| 238 | **仍成立**（缺陷在树上；⚠ 状态＝owner 已否决目的，非"可派活"） | ① 模块名现取：`git show HEAD:go.mod \| head -3`；② 票面现量 1 的尺复跑：`git --no-pager grep -l 'wisp/internal/winsec' HEAD -- internal/ cmd/ \| grep -v 'internal/winsec/'` | ①＝`module github.com/CarlosShao/wisp`。② ⇒ 全仓 winsec 非本包 importer 文件＝**11 枚**（产码 8 枚：`internal/agent/spill.go:15`、`internal/config/migrate.go:10`、`internal/config/parse.go:15`、`internal/memory/artifacts.go:14`、`internal/memory/open.go:18`、`internal/risk/winsec_c26.go:3`、`internal/secret/migrate.go:14`、`internal/secret/store.go:10`；`_test.go` 3 枚：`internal/risk/pathresolver_expansion_test.go:9`、`cmd/wisp/secret_dataroot_119b_test.go:49`、`internal/config/c26_seam_posix_125_test.go:47`）——**`internal/proc`、`internal/observe`、`internal/buildinfo` 三包各 0 枚**（逐包尺：`git --no-pager grep -c 'wisp/internal/winsec' HEAD -- internal/observe/` 与 `-- internal/buildinfo/` 均 rc=1 零命中）⇒ 票面现量 1"导入数＝0"今复跑仍成立（含它点名的 `buildinfo` 那半句）；`internal/proc` 包内出现的 `winsec` 字样全部是注释（`envfork.go:112/113/116/138/145/148/156/241`、`envfork_test.go:115/147`）。⚠ 状态凭据（不是档位凭据）＝票面 `:36-43`「owner 裁定（09-30 16:5x，台账 `A476`）＝本条不做」逐字原话＋台账 `docs/reports/pending-and-issues.md:9904` `A476`。本腿不据"owner 否决"把档改成"已失效"——缺陷本体没从树上消失，消失的是修它的授权。 |
| 242 | **差翻勾**（凭据种类＝票内注记＋台账 `A559`＋非实现者裁决表**文件名**） | ① `git ls-tree -r --name-only HEAD -- internal/agent/approval/ \| grep -iE '242\|259'`；② `git --no-pager grep -n 'func Test' HEAD -- internal/agent/approval/ticket242_panelface_test.go internal/agent/approval/ticket259_panel_capability_rulers_test.go`；③ 载具前置复跑：`git --no-pager grep -n 'CorrelationID: taskID' HEAD -- cmd/wisp/`；④ `ls docs/evidence/s1/ \| grep -E '^242'`（只取件名，⛔ 未读内容） | ①⇒ 4 件在树：`ticket242_binding_test.go`／`ticket242_panelface_test.go`／`ticket259_denial_rulers_test.go`／`ticket259_panel_capability_rulers_test.go`。②⇒ 出向读面**能力尺 8 枚**在树（`ticket242_panelface_test.go:30/:53` 两枚＝commit `2048d6b6` 10-02「242-r2 甲形落地」；`ticket259_panel_capability_rulers_test.go:66/:93/:169/:184/:216/:283` 六枚＝commit `b6b1d6a4` 10-05「259-r1 AC#2 拒因＋AC#3 三枚能力尺」）⇒ 票面现量 3 那句"仪器＝0 枚"今复跑**不再成立**（缺陷的"零仪器"半边已被做掉）。③⇒ `HEAD:cmd/wisp/subagent_selfapproval_197_test.go:109` 逐字 `TaskID: taskID, CorrelationID: taskID,` **仍在**＝票面 AC#1 的载具前置在盘上未成立。④⇒ `242-grant-binding-v1.md` 存在。⚠ 本腿**不裁它可结案**：票面 `:63` 逐字「三格保持 `[ ]`、⛔ 不加 `-done`、⛔ 不撤票」，且绑定层"恒等式"本体已被编排者就地拆出另立**票 259**（`ls` 现量：`.scratch/wisp/issues/259-grant-binding-identity-ruler-holes.md` 在池、`AC#0` 一格已勾、AC#1-#5 未勾），翻勾与"242 还欠什么"归编排者。 |
| 244 | **量不到**（卡在哪一层＝**欠读数**） | ① 构建链半边：`git --no-pager grep -n 'H=windowsgui' HEAD -- scripts/ .github/`；② 产物读数：`objdump -p build/wisp.exe \| grep -i subsystem`＋`ls -l build/wisp.exe`；③ subsystem 那把尺在场性：`git --no-pager grep -n 'debug/pe' HEAD -- scripts/ tools/ internal/ cmd/` 与 `git --no-pager grep -n 'objdump' HEAD -- scripts/ .github/ tools/`；④ 注释两行真身：`git --no-pager grep -n 'windowsgui\|GUI-subsystem' HEAD -- cmd/wisp/console_other.go cmd/wisp/console_windows.go` | ①⇒ 命中 `scripts/build.ps1:103`（why 注释）＋`:115` 逐字 `"-H=windowsgui"`（＝commit `cc6eaa65` 10-02 19:39「244-r1: build.ps1 追加 -H=windowsgui」，`--name-only` 命中 `scripts/build.ps1`＋`docs/evidence/s1/244-black-console-r1.md`），另有 `scripts/slo-check.ps1:58/:158` 两处引用。⇒ **"构建链里根本没有这一维"那半边缺陷已被写腿做掉**，本腿不判它。②⇒ 盘上产物今读＝`Subsystem 00000003 (Windows CUI)`，mtime `Sep 30 23:43` **早于** `cc6eaa65`＝票面 §9:92 那句"'旗标已落'≠'台上那枚 exe 双击没黑窗'"今天照原样成立，⛔ 本腿不许拿这枚 CUI 读成"构建链缺陷还在"、也不许拿它读成"双击还有黑窗"（产物不可再生面＝要真跑一次 `scripts/build.ps1` 才兑现，那是 go build，本腿零 go）。③⇒ 两把尺**各 0 命中**＝票 §9:93"全仓零枚测试/脚本读 PE subsystem"的"会响的尺真空"半边今仍在树。④⇒ `console_other.go:7`「links as a GUI-subsystem binary (ticket 07)」＋`console_windows.go:26`「the windowsgui subsystem (the final GUI build, ticket 07)」两行**原文未改**（AC#3 未落地；注：构建链现已带 `-H`，那句对 build.ps1 产物已是真话、对台件二进制仍是假话——差别在 `244-c2` 名册，本腿不复量）。**缺的读数逐行具名**：(a) AC#1"改后必绿"那发 GUI(2) objdump 读数；(b) AC#2 三发真机读数（父控制台里 `wisp run`／`doctor > out.txt` 重定向／**真双击·explorer 拉起无黑框**）；(c) AC#5 ⓐ"GUI 进程任务源＋停机源从哪来"的答句与 ⓑ"常驻腿零任务生产者必响"那枚尺；(d) **AC#6＝编排者裁过的"下游消费者名册"欠账**（票面 `:31-32`，10-04 追加，与 AC#5 同档、排 AC#1 之前；本腿现量：`git --no-pager grep -rn 'windowsgui' HEAD -- scripts/` 只 4 命中、无一名册件 ⇒ 名册在盘上仍未见）。 |
| 247 | **仍成立**（票面三把现量尺今复跑全部照原样） | ① 现量 1（非测试 importer＝0）：`git --no-pager grep -l 'wisp/internal/audio' HEAD -- internal/ cmd/ tools/ \| grep -v 'HEAD:internal/audio/'`；② 现量 2（`SetAudioLevel` 非测试生产者）：`git --no-pager grep -n 'SetAudioLevel' HEAD -- internal/ cmd/ \| grep -v '_test'`；③ 现量 3：`git ls-tree -r --name-only HEAD -- internal/speech/`；④ 那枚接缝的调用者：`git --no-pager grep -n 'SpawnCapture' HEAD -- internal/ cmd/` | ①⇒ **零命中（rc=1）**＝`internal/audio` 的非测试 importer 今天仍 0 枚。★这一把"零命中"我配了两发正控（防 §4 第 2 处那类坏尺）：同一把尺的形状换目标包 `git --no-pager grep -l 'wisp/internal/observe' HEAD -- internal/ cmd/ \| grep -v 'HEAD:internal/observe/'` ⇒ **120 枚文件**（尺是活的），而 `'wisp/internal/audio'` 全仓跟踪文件里出现在 **86 枚**（票面／probes／CI 日志），但 `git --no-pager grep -l 'wisp/internal/audio' HEAD -- '*.go'` ⇒ **0 枚 `.go`**＝连测试文件都没 import 它，所以"零 importer"是真读数不是拼写错。②⇒ 产码命中只有 `cmd/balldebug/main.go:418`（`b.SetAudioLevel(v)`，`:396` 逐字 "pushes a syllabified synthetic envelope"）与 `:477`，另加注释面 `internal/ball/liquid_windows.go:42` 的定义与 5 行说明 ⇒ 球的数仍只来自那枚调试 cmd。③⇒ `internal/speech/` 今树**只有 `doc.go` 一枚**。④⇒ `SpawnCapture` 只有定义（`internal/audio/audio.go:176` 注／`:180` func），**调用者 0 枚**。⛔ 本腿不据 ④ 外推"接上就能用"：AC#2（真机两形电平读数）／AC#6（拔麦降级照跑）落真设备，属"量不到"那一档的**格子**，但本票**判据本体**（"零 importer、没人读真电平"）是静态可读的形状，我已直接量到 ⇒ 整票判〔仍成立〕。AC#0 那格要求的"四问代价表"已由 `247-a1` 交件（票面 `:35-50` 编排者八问裁定在册），排程那句"`247-r1` 排在票 33 之后"是**在飞/待派**状态，不是"已被别人做掉"。 |
| 253 | **仍成立**（三格中两格的缺陷在树；一格已被做掉但编排者明令不翻） | ① AC#1 的钉形状：`git --no-pager grep -n 'postMessage(' HEAD -- internal/panel/composer_test.go`；② 有没有任何尺断言产码里 `wispDispatch` 那枚绑定在：`git --no-pager grep -rn 'panelDispatchBinding' HEAD -- cmd/wisp/*_test.go`（Git Bash 展开）；③ AC#2 那把尺在场性：`git --no-pager grep -n 'func Test' HEAD -- internal/panel/inbound_roster_253_test.go`；④ 冻结件射程只按名取：`git ls-tree -r --name-only HEAD -- internal/panel/ \| grep -E 'l2_grant_boundary\|tokens_fourway'` | ①⇒ 词面钉本体仍在：`composer_test.go:380` 逐字注释 "is not counted as a call site"＋`:381` 那行正则（匹配"点号＋postMessage＋左括号"这一形状），`routeLiteralRe` 在 `:394` 仍只认 `panel.` 前缀的字面量 ⇒ 票面现量 1"可达性钉看不见生产那扇门"今成立。②⇒ **零命中**＝没有任何尺钉 `w.Bind("wispDispatch")` 的存在，票 §8:42 裁的"真洞只在窄义"那格**仍未落**（写面＝`cmd/wisp`）。③⇒ 五枚测试在树：`inbound_roster_253_test.go:419/:597/:696/:769/:809`＝`87bc6aca`（10-03 09:50「253-r5 票 253 AC#2 落地（形ⓐ）」，`--name-only` 命中 `internal/panel/inbound_roster_253_test.go`）⇒ **AC#2 那格的缺陷（"新增方法名不进名册也绿"）已被别人做掉**，正控 `TestPlantedUnregisteredInboundMethodNameGoesRed` 亦在树。④⇒ `internal/panel/l2_grant_boundary_test.go`＋`tokens_fourway_test.go` 均在树（本腿一字未读内容）＝**契约邻接面**未动，形ⓑ 已被裁"不解冻"。**它卡在哪一层＝契约邻接＋欠读数**：AC#2 卡"冻结件射程已裁ⓐ、编排者票面 `:52` 逐字「⛔ AC#2 这格我不翻（实现者自述不算凭据）：待非实现者验收腿 `253-v1`」"，而本腿现量 `ls docs/evidence/s1/ \| grep -E '^25[0-9]'`＝`252-allowlist-sameform-v1.md`／`255-tier-registry-r1.md`／`255-tier-registry-v1.md`／`257-clean-machine-provider-registry-v1.md`——**无 `253-*` 件** ⇒ 验收腿未交；AC#3 卡欠读数（"页面拿到的不是无声 `null`"必须真开 WebView2 窗才读得到；Go 侧那半边我读到 `composer_dispatch.go:155-158` 的 `Handle` 在 parse 失败时 `d.record(req, err)`＋`RefusedEnvelopeForUser` ⇒ 进得了派发的形状**已有声**，但票面现量 3 指的是**没走 `wispDispatch`、直接 `postMessage` 的旧封套**那一支，落点在第三方库 `webview.go:162-168`，本仓射程外 ⇒ 判不到）。 |
| 262 | **差翻勾**（凭据种类＝今树复算尺＋`git log --name-only` 命中一枚 commit＋票内 Progress log `262-r1` 注记＋台账 `A585`；⚠ **不等于可结案**，两格仍欠读数） | ① `git ls-tree -r --name-only HEAD -- scripts/ \| grep -i path-length`；② `git --no-pager grep -n 'check-path-length-budget' HEAD -- .github/`；③ `git --no-pager log -1 --format='%h %ad %s' --name-only 519e1ade`；④ 本腿自己现跑那把门（只读，零写树）：`sh scripts/check-path-length-budget.sh` 与 `sh scripts/check-path-length-budget.sh --with-self-test`；⑤ 分母独立复算：`git ls-files \| awk 'length($0)>121' \| wc -l` | ①⇒ `scripts/check-path-length-budget.sh` 在树，`git --no-pager grep -c '' HEAD -- scripts/check-path-length-budget.sh`＝**577 行**。②⇒ 调用点 `.github/workflows/ci.yml:166` 逐字 `run: sh scripts/check-path-length-budget.sh --with-self-test`（＋`:153` 注释具名同源）。③⇒ `519e1ade 10-04 09:42 262-r1 落地：跟踪路径长度门禁（丙＋丁同批，阈值钉在帽＝相对 >121）`，`--name-only` 命中 `.github/workflows/ci.yml`＋`scripts/check-path-length-budget.sh`；续枚 `6a735939`（10-04 09:57，终态读数＋脚本 `:178-213` 的 AC#7「旧名可追」11 枚对）。④⇒ 本腿现跑两发：只读那发 `VERDICT GREEN … roster equals the tree`、rc=0，读数含 `tracked paths=6357 / over-budget=57 / covered by roster=57 / not in roster=0 / longest=180 … 252-…-r2-l2.md / worst full path=224`；带 self-test 那发首行 `positive control PASSED`（＝AC#3 那枚"种相对 >121 必红"的自证今天在门里跑），末行仍 GREEN，rc=0。⑤⇒ **57**，与脚本名册枚数一致。⇒ 票面标题那句"全仓没有任何一把尺盯跟踪路径长度"**已失效**（满足"已失效"的硬门：commit 号＋`--name-only` 命中＋今尺在场且有牙）。⚠ 停在〔差翻勾〕不升〔已失效〕的具名理由＝AC#5（"该步在 self-hosted 那一侧跑过且有日志行可指"属 **CI 侧读数**，`262-r1` 票面 `:67` 自己逐字标注"本腿给不了（无推送权）"）与 AC#6（卫生四门要跑 `gofumpt`/`go vet`/`d22scan`，本腿**零 go**）两格无凭据；且台账 `docs/reports/pending-and-issues.md:11537` 逐字「`262-r1` 已落地 ⇒ **`262-v1` 待派**（含 `A593` §3 那枚"墙口径"头一格）」＋本腿现量 `ls docs/evidence/s1/ \| grep -E '^26[2468]'`＝**0 件** ⇒ 非实现者裁决表尚未存在。翻勾与"能不能只按 AC#1-#4 结"归编排者。 |
| 264 | **量不到**（卡在哪一层＝**等人拍**，不是"没人做"） | ① 判据本体（AC#0 那一格要求的读数）：`grep -n '264 的三栏\|264-a1' docs/reports/pending-and-issues.md` 取出台账两行原文；② 票面三形代价表的在场性：`grep -n 'AC#0' .scratch/wisp/issues/264-*.md`；③ 我**只**对票面 §现量 做静态复跑（当辅助读数，不当档位凭据）：`git --no-pager grep -nE 'edact' HEAD -- internal/agent/ internal/tools/ \| grep -v '_test.go'`、`git --no-pager grep -n '路径无法解析（按 fail-closed 拒绝）\|打开失败：\|原子重命名失败' HEAD -- internal/tools/fs.go internal/tools/fs_write.go`、`git --no-pager grep -n 'log.Text = out.Text\|func toolResultMessage' HEAD -- internal/agent/loop.go` | ①⇒ 台账两处：`docs/reports/pending-and-issues.md:11537` 逐字「`264-a1` 在飞（回来后把票 264 的三栏代价补实**再摆给机主**）」、`:11829` 把「票 264 的三栏（含"不做"那一栏）」列进待人拍板那一组，同句逐字「**默认动作全部是不做／不动**」；★本腿另按"台账只追加不删"把同一把尺推到**最新一节**（`grep -n '^## A[0-9]' docs/reports/pending-and-issues.md \| tail -3`＝`A633`／`A634`／`A635`）：`:12489` 逐字「**待人项**：原八枚在册未摆（…／**票 264**／…）…机主两次回"继续／恢复"＝未选，不重复追问」⇒ 机主那一句话**到 10-06 仍未取回**。补充读数（把 `:11537` 那句"在飞"纠正过来）：`A597` §4＝台账 `:11604` 已收 `264-a1`（件 `.scratch/wisp/probes/264/a1/census.md`，本腿 `wc -l -c` 现量＝**353 行／58,686 字节**），编排者在该节逐字改口"⛔ 那句推荐我改口"并"重写后的三栏与我在同一轮里发给他"⇒ **普查腿已交、"摆给机主"这一步还没走完**；本格要的读数不在盘上任何一棵树里，静态尺读不出（⛔ 不许由"没找到那句话"外推成"这功能没人做"）。②⇒ 票面 `:26` 逐字「**AC#0（第一格，⛔ 不许动产业代码，⛔ 本轮不派写腿）**：把上表三栏摆给机主并**取回他的一句话**（甲／乙／不做）。取回之前本票任何格都不许开工」；`:40` 「**此刻不派**（AC#0 未取回）」⇒ AC#1/AC#2/AC#3 三格**按票面规则今天就不许开工**，判"差翻勾"或"已失效"都无凭据。③辅助读数（只说明"外发那半边今天仍是真"，⛔ 不用它定档）：agent/tools 非测试 `edact` ⇒ **零命中**（`-l` 尺回的 files＝空，脱敏那半边今天仍不在这两包里）；源头装配点今在树＝`internal/tools/fs.go:144`／`:148`／`:208`、`internal/tools/fs_write.go:260`／`:338`／`:343`／`:389`／`:460`／`:464`／`:611`；投递链＝`internal/agent/loop.go:698`(`log.Text = out.Text`)→`:722`(`l.append(toolResultMessage(c.ID, log.Text, …))`)→`:1080`(`func toolResultMessage`) 三跳逐字仍在。⚠ 一枚必须跟着这行读的限定（台账 `A597` §4＝`:11605` 逐字）：「`fs.go:144` 那一支今天**走不到**（`Actable()` 在 `internal/tools` 零调用者）」⇒ 我在表里数到的是**装配点在树上的字面**，不是"那一条今天真会发出去"；同节还纠正两枚坐标（`cmd/wisp/task_output_pointer_notice_test.go` 应为 `internal/tools/…`；`internal/risk/paths.go:122-126` 不存在，真身 `internal/tools/paths.go:119`），编排者在那一节写死「在我重写票面之前 ⛔ 任何腿不许照这三个坐标下刀」——本腿只读不判，照这一句把坐标原样带着。另复量到一枚**与票面"不做"那一栏同一维**的事实：`internal/config/schema.go:519` 逐字 `RedactPaths bool` ＋ toml 标签 `redact_paths` ＋ `default:"false"`（`:518` 注释逐字 "masks filesystem paths in user-visible text"），而 `cmd/wisp/config_readers_255.go:134` 自述 "the mirror field exists and is filled by callers, **never from cfg.Privacy**" ⇒ 开关默认关且模型侧那条路上没有读者。 |
| 266 | **仍成立**（票面现量 1 的尺今复跑照原样；⚠ 唯一能关掉它的形落在人工批准面上） | ① 现量 1（d22scan 的 scope 里没有 `scripts/`）：`git show HEAD:tools/d22scan/main.go \| grep -c '"scripts'` ＋ walk 起点逐枚 `git show HEAD:tools/d22scan/main.go \| grep -n 'filepath.Join(root,'`；② 有没有一枚外部件读那两枚钉的名字：`git --no-pager grep -rn 'Test-ScriptShapeNail' HEAD -- . \| grep -v 'scripts/slo-check.ps1'`；③ 唯一外尺的射程：`git show HEAD:scripts/slo-freshness.sh \| grep -c 'Test-ScriptShapeNail\|Test-ExitCodeInstrument'` ＋ `git --no-pager grep -n 'sample_max_age_days' HEAD -- scripts/slo-freshness.sh`；④ 两枚钉今天真身（不照抄票面行号）：`git --no-pager grep -n 'Fail-InstrumentBroken' HEAD -- scripts/slo-check.ps1 \| grep -i 'shape nail'` | ①⇒ **0**；walk 起点只有 `:264 internal/`、`:267 cmd/`、`:271 frontend/`、`:274 internal/tools/`（另 `:409/:413/:417/:433` 是同一四名的 label 表）⇒ `scripts/**` 今天仍不在任何扫描器射程内。②⇒ 命中**全在 `.scratch/wisp/probes/**`**（`263/r1/*`、`263/v1/run-mutations-v1.ps1:77`、`263/v1/run_mutations_v1.py:29`、`266/a1/census.md:47`）＝只有**一次性突变台件**与普查件在引用那两枚钉的名字，**入库的门／CI／扫描器里没有一枚外部件数它**（⛔ 台件不算外牙：它只在那几发突变跑的当刻被调用，不在任何一条常跑门禁路径上，谁删钉都不会让它响）。③⇒ **0**＝`slo-freshness.sh` 只查 ci.yml 的 `slo-full` job 是否还跑 `-Subset full`（`:160-171` 那段，判的是**触发器在场**），不读那两枚钉；aging 常量在 `:97` `sample_max_age_days=${SLO_FULL_SAMPLE_MAX_AGE_DAYS:-3}`、判定在 `:368` `if [ "$sample_age_days" -gt "$sample_max_age_days" ]` ⇒ 票面现量 3 那句"唯一外牙按天钝、且看产物不看钉的存在"今成立。④⇒ 词面钉 `scripts/slo-check.ps1:249`（PSCommandPath 空）／`:256`（命中即 `Fail-InstrumentBroken`）仍在，函数体 `:233 Test-ScriptShapeNail`、调用点 `:291`——**两枚钉都还是只由被检文件自己守着**。★ 一格如实附注（不改档位）：票面 `:43-46` 已裁「形先不裁，因为唯一闭合那一形（③）落在人工批准面上；已落 `Q-78` 待机主一句话」，台账 `:11785` `A604` 与 `:11810` 逐字「这一格我登记为**未闭合**，不写成"已缓解"」，`Q-78` 在 `:11829` 那组待拍板项里在案 ⇒ 判〔仍成立〕不等于"催派单"，它等的是 owner 那一句。 |

## §2bis 本批计数（8 枚全判完，⛔ 未判 0 枚）

| 档 | 枚数 | 票号 |
|---|---|---|
| 仍成立 | **4** | 238 247 253 266 |
| 已失效／已被别人做掉 | **0** | —（无一枚满足"硬门＋缺陷整体消失"两样同时成立；262／242 的缺陷半边已消失但都被下面两档接走） |
| 差翻勾 | **2** | 242 262 |
| 量不到 | **2** | 244（欠读数）264（等人拍） |
| 合计 | **8** | ✓ 与 §1 名册逐枚对上 |

★ 三枚有特殊状态的那一行（按 §0 的要求逐枚给"卡在哪一层"）：
`244`＝**欠读数**（AC#1 改后产物 objdump 那一发要真跑构建、AC#2 三发真机含真双击、AC#5 两截、AC#6 名册）；
`264`＝**等人拍**（票面 `:26`／`:40` 自己规定"AC#0 未取回之前任何格不许开工"，台账 `:11537` 与 `:11829` 把它列在待人拍板那一组，默认动作＝不做）；
`253`＝**契约邻接**（AC#2 已被 `87bc6aca` 的形ⓐ 做掉，编排者明令不冻结结件；AC#1／AC#3 写面撞三枚冻结件同包与 `cmd/wisp`，且 `253-v1` 未交）；
另附一枚同层：`238`＝**owner 已否决目的**（票面 `:36-43`／台账 `A476`），缺陷在树但授权不在 ⇒ 不计入"真活"。

## §3 判得心虚的枚数与具名理由（6 枚，逐枚点名）

1. **238（仍成立／但它不是活）**：四档尺第一档的定义是"票面点名的缺陷今天还在树上"，我量到的确实是"在树"（`internal/proc`、`internal/observe` 对 `winsec` 的 import 各 0 枚，尺见 §2 ②）。心虚的**不是档位**，是它容易被读成"还剩一件真活"：owner 09-30 已把**目的本身**否掉（票面逐字「锁个屁…」＝台账 `A476`），四枚框保持未勾是**裁定的结果**而不是欠账。⇒ 我给编排者的换算句应是"仍成立 4 枚里 1 枚已被 owner 关死"，而不是"还剩 4 件"。
2. **242（差翻勾／但盘上读数更像已失效）**：我能拿 commit 号（`2048d6b6`、`b6b1d6a4`）＋今树尺在场（8 枚能力尺）——按第二档的硬门它够格写"已被别人做掉"。我落在第三档的**唯一理由**是凭据种类：票面 `:50-63` 那节是编排者收件＋台账 `A559`＋裁决表 `242-grant-binding-v1.md`（我只按件名取），且票面 `:63` 逐字"三格保持 `[ ]`、⛔ 不加 `-done`、⛔ 不撤票"。**若编排者按"盘上有没有那把尺"来换算，这一枚应当读成第二档**；两形读数都在 §2 的复法列里，我没有换措辞。另有一处我没量：绑定层"恒等式"本体已移交给**票 259**（在池、`AC#0` 已勾、AC#1-#5 未勾），所以"242 的缺陷还剩多少"这句话单靠 242 票面读不出来。
3. **244（量不到／它其实是一枚混合票）**：按 AC 逐枚读，档位并不一致——AC#1 的构建链半边已落（`cc6eaa65`＋`scripts/build.ps1:115` 逐字 `"-H=windowsgui"`＝硬门满足）、AC#3 那两行注释我量到**原文未改**（`console_other.go:7`／`console_windows.go:26`＝按第一档定义它是"仍成立"）、AC#2／AC#5 ⓐ／AC#6 只能是运行期或人给的读数（＝第四档）。我按**结案凭据那一枚**（AC#1 的"改后必绿"objdump 读数只能由真跑 `build.ps1` 兑现）归了〔量不到〕。★这一枚如果编排者要的是"还剩多少真活"，请按**混合票**处理，别在"0 件"或"1 件"里选一边。
4. **264（量不到／但外发那半边静态量得到）**：我判"等人拍"，可票面现量的三跳投递链与"agent/tools 非测试零脱敏"我今天复跑仍全部成立——也就是说"缺陷还在"这件事**是有静态凭据的**。我没有写成〔仍成立〕是因为票面自己把 AC#0 定成"取回一句话之前任何格不许开工"（`:26`），写"仍成立"会让人以为可以派写腿。⚠ 换算请写"仍成立＋等 owner 一句话"，两种写法都对盘，只是档位尺只许选一边。
5. **262（差翻勾／正控那半边是我自己跑的，不是别人交的）**：AC#3 的"种一枚相对 >121 必红"我拿的是脚本**自带的** `--with-self-test` 读数（`positive control PASSED`）。本仓有一条铁律是"实现者自述不算凭据"（票 242 §9 与票 253 票面 `:52` 都逐字写过），而 `262-r1` 的读数件属实现者自跑。我把"缺陷已被做掉"的半边凭据压在这把**自证**尺上 ⇒ 偏松，如实标出；真正**外**于实现者的凭据只有 `ls docs/evidence/s1/ | grep -E '^26[2468]'`＝**0 件**这一条（＝还没有非实现者裁过），这也是我不写〔已失效〕的第二个理由。
6. **253（仍成立／但三格方向不同）**：AC#2 已被 `87bc6aca` 做掉（今树五枚尺在场，含一枚正控名 `TestPlantedUnregisteredInboundMethodNameGoesRed`），AC#1 我量到"钉只认词面"仍在＋没有任何尺钉 `wispDispatch` 那枚绑定（`grep -rn 'panelDispatchBinding' HEAD -- cmd/wisp/*_test.go` 零命中），AC#3 落点在依赖库里、静态尺结构上读不到。整票写〔仍成立〕靠的是 AC#1 那一格；`docs/evidence/s1/` 里 `253-*`＝0 件 ⇒ 编排者那句"待非实现者验收腿 253-v1"仍未兑现，翻不翻勾我一律不碰。

## §4 一处尺有歧义／判不动的地方（6 处，全部留两把读数原文）

1. **HEAD 在我这双腿之间漂了（同一枚尺第二遍拿不到同一读数）**
   起手：`git --no-pager log -1 --format='%h %ad' --date=iso` ＝ `3d9b8374 2026-10-06 13:12:18 +0800`。
   13:38 我自己落 `41d12ca7`（骨架＋前三枚判档），13:43 同职腿 4c 落 `9621e15c`。
   复跑 262 那把门两发：只读那发＝`denominator read: 6357 tracked paths`；`--with-self-test` 那发＝`6358 tracked paths`，两发都 `VERDICT GREEN`、rc=0。
   同一段时间里 `git ls-files | wc -l` 也从 6357 变 6358（差 1 枚＝我自己那枚 `msg-b1.txt` 进了跟踪树）。
   ⇒ 歧义在**口径**不在结论："跟踪路径枚数"这类读数在共享工作树里会随**本腿自己的 commit** 漂；任何引用它的人都得连 HEAD 短号＋取数时刻一起写（262 票面 AC#2 正是这一形）。我没有因此改判 262。

2. **238 那把尺的第一遍是坏尺（短词面把注释算成 import）**
   第一遍原文：`git --no-pager grep -c 'internal/winsec' HEAD -- internal/proc/` ⇒ 回 `HEAD:internal/proc/envfork.go:3`、`envfork_test.go:2` ＝**看起来有 import**。
   第二遍原文：`git --no-pager grep -n 'wisp/internal/winsec' HEAD -- internal/ cmd/`（先 `git show HEAD:go.mod` 现取模块名＝`github.com/CarlosShao/wisp`，import 串必带 `wisp/`）⇒ `internal/proc` **零命中**，包内 10 处 `winsec` 字样逐枚读回全是注释行（`envfork.go:112/113/116/138/145/148/156/241`、`envfork_test.go:115/147`）。
   ⇒ 具名教训：**照票面标题那种短写法（`对 winsec 导入数＝0`）直接 grep 会得假红**。244 票面现量表 `:14` 已经自陈过一次同类学费（"把中文分支和 `windowsgui` 拼进同一个模式串，回了零命中＝坏尺"，台账 `A473`）——两枚是同一种病的两个方向（一个假红、一个假绿）。

3. **264 的"零命中"我用两把分母尺证伪过，但它仍可能被下一位读错**
   `git --no-pager grep -nE 'redact|Redact' HEAD -- internal/agent/ internal/tools/ \| grep -v '_test.go'` ⇒ rc=1、零命中。
   证伪（射程非空）：`git ls-tree -r --name-only HEAD -- internal/agent/ \| grep -c '\.go$'`＝**62**、`internal/tools/`＝**72**。
   交件前在 `6f0a2d83` 又分两发重量：不带过滤 `git --no-pager grep -nE 'redact\|Redact' HEAD -- internal/agent/ internal/tools/` ⇒ **rc_all＝1（零命中）**，加 `\| grep -v '_test\.go'` ⇒ **rc_nontest＝1**（两把都零＝那条 `grep -v` 没洗掉任何命中，"零"不是过滤造出来的）。
   ★但**全仓 `Redact` 不是零命中**：`internal/config/schema.go:519` `RedactPaths bool toml:"redact_paths" default:"false"`、`internal/observe/logging.go:51`、`internal/observe/diagnostics.go:120/:138` 都在树。⇒ 那句"这一侧一枚仪器都没有"的射程是**模型侧那一跳**（agent/tools），不是"仓里没有脱敏码"。
   歧义点：`cmd/wisp/config_readers_255.go:134` 逐字 "the mirror field exists and is filled by callers, **never from cfg.Privacy**"——即那枚开关连日志侧都吃不到配置。这句话属票面"不做"那一栏现量要求"再写"的那一格（`:22`），本腿把它复认了一行，⛔ 不替它翻勾。

4. **253 AC#3 判不动的那一层是尺的物理边界，不是我没跑**
   票面现量 3 的落点＝依赖模块 `webview.go:162-168`（"旧封套被宿主库静默 `resolve(null)`"）。本腿现查：`git ls-files \| grep -i webview` 只得 `docs/evidence/s0/json/06-webview2-run2.json`、`06-webview2.json`、`scripts/spike/webview2-latency/main.go` 三枚 ⇒ 仓里没有那枚文件；`git show HEAD:go.mod` 的 require 行是 `github.com/jchv/go-webview2`（`:19`）。
   ⇒ `git grep HEAD` 这类**只读跟踪树**的尺**结构上不可能**读到它。本仓可读的那半边（`internal/panel/composer_dispatch.go:155-158`：parse 失败走 `d.record(req, err)`＋`RefusedEnvelopeForUser`）**与票面现量 3 不是同一支**（那一支要求页面绕开 `wispDispatch`、直接 `postMessage`）。我把这一格写成"量不到"，⛔ 不许下一位据此说"Go 侧已经有声了＝AC#3 齐了"。

5. **244 的两把尺方向相反，引用时必须连"读的是链还是台"一起写**
   `git --no-pager grep -n 'H=windowsgui' HEAD -- scripts/ .github/` ⇒ `scripts/build.ps1:115` 逐字 `"-H=windowsgui"`（＝"链上已落"）。
   `objdump -p build/wisp.exe \| grep -i subsystem` ⇒ `Subsystem 00000003 (Windows CUI)`，`ls -l build/wisp.exe` mtime＝`Sep 30 23:43`，而 `cc6eaa65` 落盘 10-02 19:39（＝"台上还没落"，产物比 commit 旧）。
   ⇒ 两把都对，但**同一条票上互不相容**：只引第一把会说"切完了"，只引第二把会说"没切"。票面 §9:92 已自陈这一形（"旗标已落 ≠ 台上那枚 exe 双击没黑窗"）。本腿两把都留在表里，没挑一把。

6. **244 顺带复跑到两枚前置仍未落（非档位凭据，给下一位省一把尺）**
   `git --no-pager grep -n '允许一次\|AllowOnce' HEAD -- internal/ball/ cmd/wisp/` ⇒ 零命中；`cmd/wisp/resident_ball_windows.go:455` 今树仍逐字写着 "Giving this item an executor needs either a stop-request hook in…"。
   ⇒ 票面 AC#5 ⓐ（GUI 进程的任务源／停机源）与 §6 G2／J10 归口票 228 AC#2／AC#11 那两枚执行者**今天都还不存在**，"244 能不能切 `-H`"仍卡在它们后面。★这一条的零命中我按"中文多分支 grep 得零命中先证伪"那条规矩配了两发正控：`git --no-pager grep -c '允许' HEAD -- cmd/wisp/ internal/ball/` 非空（`cmd/wisp/approval_reply.go`＝17、`approval_always.go`＝3 等），且同一枚托盘文件的四名名册 `git --no-pager grep -n '打开面板\|静音\|暂停唤醒\|退出' HEAD -- internal/ball/tray_windows.go` 命中 `:86/:88/:89/:91` 四枚 ⇒ 射程非空、尺是好的，"第五枚『允许一次』没落"是真读数不是坏尺；`:455` 那句注释命中在树同理（同文件尺）。⛔ 两把零命中里我只信配了正控的那一把。

## §5 交件自证（本腿到底做了什么、没做什么）

- 零 go 命令：本腿跑过的全部命令只有 `git`（`log`／`grep HEAD`／`show HEAD:<path>`／`ls-tree`／`ls-files`／`rev-list`／`add`／`commit`／`diff`／`status`）、`grep`、`wc`、`awk`、`ls`、`sed -n`、`objdump -p`（读一枚已存在二进制的 PE 头，不构建、不执行它）、`sh scripts/check-path-length-budget.sh`（票 262 那把门本体，纯 shell＋`git ls-files`，只读，无写树）＋一枚本目录下的 `check-table.py`（数 markdown 表格列数，只为核 §2 那张表没被中文管道符打断）。⛔ `go test`／`go build`／`go vet`／`go run`／`go list`／`go env` **一次都没出现**。
- ★交件前重量：八枚的判据尺在 `6f0a2d83`（14:01:55）逐枚复跑一遍（逐枚读数写在开头那块引言里），档位零变化；两把"零命中"尺（247 的 importer、264 的 `redact`）各配了正控或双发对照，防 §4 第 2 处那类坏尺。
- 零写已有文件：写入只有 `.scratch/wisp/probes/pool-validity/4b/batch4b.md` 与本目录下几枚临时件（`msg-skeleton.txt`、`msg-b1.txt`、`msg-final.txt` 三枚 commit-msg 件＋`check-table.py` 一把数表格列的尺；按"临时件只建不删"留着）。commit 前 `git --no-pager diff --cached --numstat` 逐枚核过删除列只落在我自己那枚新文件上（`3 3`＝§2 表里那三行"尚未判"被我自己的判档替换，全在同一枚我的新文件内）。
- 票面零改动：八枚票的 `- [ ]`／`- [x]` 枚数、`Status:` 行、文件名一枚没动；台账与 `docs/**` 一字不碰（只按行号**读**了两处 `A476`／`Q-78` 相关的原文行与 `docs/evidence/s1/` 的**件名**）。
- 未读：`frontend/**`、`design/**`、`.gitignore`、`docs/evidence/s1/**` 内容、`.scratch/wisp/probes/236/**`、`.scratch/wisp/probes/pool-validity/4c/**`；`cmd/wisp/**` 与 `internal/**` 一律经 `git show HEAD:<path>` 取，没读工作树内容面。
