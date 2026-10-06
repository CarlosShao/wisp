# pool-validity-4b / 批4b —— 零勾开放票有效性普查（号段 ≥238，现量 8 枚）

> 只读普查腿 `pool-validity-4b`。任务＝把"238 以上号段里一格都没勾的票"逐枚判档，供编排者回答"还剩多少真活"。
> 起手 HEAD `3d9b8374`（2026-10-06 13:12:18 +0800）。⚠ 换 HEAD 要重量。
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
| 238 | **仍成立**（缺陷在树上；⚠ 状态＝owner 已否决目的，非"可派活"） | ① 模块名现取：`git show HEAD:go.mod \| head -3`；② 票面现量 1 的尺复跑：`git --no-pager grep -l 'wisp/internal/winsec' HEAD -- internal/ cmd/ \| grep -v 'internal/winsec/'` | ①＝`module github.com/CarlosShao/wisp`。② ⇒ 全仓 winsec 非本包 importer 文件＝**11 枚**（产码 8 枚：`internal/agent/spill.go:15`、`internal/config/migrate.go:10`、`internal/config/parse.go:15`、`internal/memory/artifacts.go:14`、`internal/memory/open.go:18`、`internal/risk/winsec_c26.go:3`、`internal/secret/migrate.go:14`、`internal/secret/store.go:10`；`_test.go` 3 枚：`internal/risk/pathresolver_expansion_test.go:9`、`cmd/wisp/secret_dataroot_119b_test.go:49`、`internal/config/c26_seam_posix_125_test.go:47`）——**其中 `internal/proc` 0 枚、`internal/observe` 0 枚** ⇒ 票面现量 1"导入数＝0"今复跑仍成立；`internal/proc` 包内出现的 `winsec` 字样全部是注释（`envfork.go:112/113/116/138/145/148/156/241`、`envfork_test.go:115/147`）。⚠ 状态凭据（不是档位凭据）＝票面 `:36-43`「owner 裁定（09-30 16:5x，台账 `A476`）＝本条不做」逐字原话＋台账 `docs/reports/pending-and-issues.md:9904` `A476`。本腿不据"owner 否决"把档改成"已失效"——缺陷本体没从树上消失，消失的是修它的授权。 |
| 242 | **差翻勾**（凭据种类＝票内注记＋台账 `A559`＋非实现者裁决表**文件名**） | ① `git ls-tree -r --name-only HEAD -- internal/agent/approval/ \| grep -iE '242\|259'`；② `git --no-pager grep -n 'func Test' HEAD -- internal/agent/approval/ticket242_panelface_test.go internal/agent/approval/ticket259_panel_capability_rulers_test.go`；③ 载具前置复跑：`git --no-pager grep -n 'CorrelationID: taskID' HEAD -- cmd/wisp/`；④ `ls docs/evidence/s1/ \| grep -E '^242'`（只取件名，⛔ 未读内容） | ①⇒ 4 件在树：`ticket242_binding_test.go`／`ticket242_panelface_test.go`／`ticket259_denial_rulers_test.go`／`ticket259_panel_capability_rulers_test.go`。②⇒ 出向读面**能力尺 8 枚**在树（`ticket242_panelface_test.go:30/:53` 两枚＝commit `2048d6b6` 10-02「242-r2 甲形落地」；`ticket259_panel_capability_rulers_test.go:66/:93/:169/:184/:216/:283` 六枚＝commit `b6b1d6a4` 10-05「259-r1 AC#2 拒因＋AC#3 三枚能力尺」）⇒ 票面现量 3 那句"仪器＝0 枚"今复跑**不再成立**（缺陷的"零仪器"半边已被做掉）。③⇒ `HEAD:cmd/wisp/subagent_selfapproval_197_test.go:109` 逐字 `TaskID: taskID, CorrelationID: taskID,` **仍在**＝票面 AC#1 的载具前置在盘上未成立。④⇒ `242-grant-binding-v1.md` 存在。⚠ 本腿**不裁它可结案**：票面 `:63` 逐字「三格保持 `[ ]`、⛔ 不加 `-done`、⛔ 不撤票」，且绑定层"恒等式"本体已被编排者就地拆出另立**票 259**（`ls` 现量：`.scratch/wisp/issues/259-grant-binding-identity-ruler-holes.md` 在池、`AC#0` 一格已勾、AC#1-#5 未勾），翻勾与"242 还欠什么"归编排者。 |
| 244 | **量不到**（卡在哪一层＝**欠读数**） | ① 构建链半边：`git --no-pager grep -n 'H=windowsgui' HEAD -- scripts/ .github/`；② 产物读数：`objdump -p build/wisp.exe \| grep -i subsystem`＋`ls -l build/wisp.exe`；③ subsystem 那把尺在场性：`git --no-pager grep -n 'debug/pe' HEAD -- scripts/ tools/ internal/ cmd/` 与 `git --no-pager grep -n 'objdump' HEAD -- scripts/ .github/ tools/`；④ 注释两行真身：`git --no-pager grep -n 'windowsgui\|GUI-subsystem' HEAD -- cmd/wisp/console_other.go cmd/wisp/console_windows.go` | ①⇒ 命中 `scripts/build.ps1:103`（why 注释）＋`:115` 逐字 `"-H=windowsgui"`（＝commit `cc6eaa65` 10-02 19:39「244-r1: build.ps1 追加 -H=windowsgui」，`--name-only` 命中 `scripts/build.ps1`＋`docs/evidence/s1/244-black-console-r1.md`），另有 `scripts/slo-check.ps1:58/:158` 两处引用。⇒ **"构建链里根本没有这一维"那半边缺陷已被写腿做掉**，本腿不判它。②⇒ 盘上产物今读＝`Subsystem 00000003 (Windows CUI)`，mtime `Sep 30 23:43` **早于** `cc6eaa65`＝票面 §9:92 那句"'旗标已落'≠'台上那枚 exe 双击没黑窗'"今天照原样成立，⛔ 本腿不许拿这枚 CUI 读成"构建链缺陷还在"、也不许拿它读成"双击还有黑窗"（产物不可再生面＝要真跑一次 `scripts/build.ps1` 才兑现，那是 go build，本腿零 go）。③⇒ 两把尺**各 0 命中**＝票 §9:93"全仓零枚测试/脚本读 PE subsystem"的"会响的尺真空"半边今仍在树。④⇒ `console_other.go:7`「links as a GUI-subsystem binary (ticket 07)」＋`console_windows.go:26`「the windowsgui subsystem (the final GUI build, ticket 07)」两行**原文未改**（AC#3 未落地；注：构建链现已带 `-H`，那句对 build.ps1 产物已是真话、对台件二进制仍是假话——差别在 `244-c2` 名册，本腿不复量）。**缺的读数逐行具名**：(a) AC#1"改后必绿"那发 GUI(2) objdump 读数；(b) AC#2 三发真机读数（父控制台里 `wisp run`／`doctor > out.txt` 重定向／**真双击·explorer 拉起无黑框**）；(c) AC#5 ⓐ"GUI 进程任务源＋停机源从哪来"的答句与 ⓑ"常驻腿零任务生产者必响"那枚尺；(d) **AC#6＝编排者裁过的"下游消费者名册"欠账**（票面 `:31-32`，10-04 追加，与 AC#5 同档、排 AC#1 之前；本腿现量：`git --no-pager grep -rn 'windowsgui' HEAD -- scripts/` 只 4 命中、无一名册件 ⇒ 名册在盘上仍未见）。 |
| 247 | 尚未判，本腿下一步判它 | — | — |
| 253 | 尚未判，本腿下一步判它 | — | — |
| 262 | 尚未判，本腿下一步判它 | — | — |
| 264 | 尚未判，本腿下一步判它 | — | — |
| 266 | 尚未判，本腿下一步判它 | — | — |

## §3 判得心虚的枚数与具名理由

本骨架阶段＝0 枚已判，所以此节当前只写"还没开始判"这个事实；判档过程中每出现一枚心虚的，本腿会当场把
票号＋心虚的具体一处（哪把尺读不出"缺陷在／不在"）追加进来，不留空节。

## §4 一处尺有歧义／判不动的地方

同上：骨架阶段尚无实跑歧义记录；一旦某把尺第二遍拿不到期望读数，本腿停手并把两把读数原文写进这里。
