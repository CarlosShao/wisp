# 306-r1 · 00 起手锚＋台面＋进程闸门自陈

腿＝`306-r1`（落地腿，只写测试侧）。票＝`.scratch/wisp/issues/306-inline-authoritative-format-tag-literals-in-wavinjector.go-have-no-ruler-at-all.md`
（实际文件名＝`306-inline-authoritative-format-tag-literals-in-wavinjector-go-have-no-ruler-at-all.md`）。

## 1. 起手锚（现跑，⛔ 抄派单）

| 尺 | 读数 |
|---|---|
| `git rev-parse HEAD`（第一发，本腿进场） | `a81c2980eb261940001e90190b90456139373225` ＝ 与派单锚号**一致** |
| 派单里给的锚号 | `a81c2980` |
| ★`git rev-parse HEAD`（第二发，落锚件时同一条命令里再量） | **`5480434f82f2a2344cb78701decc97b6acba2958`** ⇒ **锚在两次读数之间漂移了一枚** |
| ★**具名报回（派单第 3 节第 1 条）** | 漂移的那枚＝编排者自己的 `5480434f`（逐字标题起头＝`销 A826 那笔欠账：票 303 AC#3 的"不可满足判据"正式登记（A828，形照票 141）＋派两枚（306-r1 写腿 只读）`），它同时把 **`306-r1` 派了出去** ⇒ 本腿进场时它已经在 `dev` 上。**以盘上为准**：本腿的工作锚＝**`5480434f`**，"改前那发"的导出树也从这一枚取。 |
| ★漂移的影响面（尺＝`git diff --name-only a81c2980 HEAD`） | 名册 22 行里**⛔ 一枚 `internal/`**（尺＝`git diff --name-only a81c2980 HEAD -- internal` ＝ **0 行**，件内 `rc_diff_internal=0`／`diff_internal_lines=0`）⇒ 只动 `docs/reports/**` ＋ `.scratch/wisp/probes/**` ＋ 票 303 文件；`wavinjector.go` 的 `:192`／`:196`／`:218` 三枚内容锚在新旧两枚 HEAD blob 上**逐字相同**（件内 `rc_sed_anchors_head=0` 那三行）。⇒ 漂移⛔ 影响本票任何一行号或判据。 |
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git status --porcelain -- internal cmd docs scripts .github` | **0 行**，件内自落 `rc_porcelain_scoped=0` ⇒ 本腿写面（`internal/audio`）起手干净 |
| `git show HEAD:internal/audio/wavinjector.go \| cmp - internal/audio/wavinjector.go` | **`rc_blob_vs_worktree_cmp=0`**（漂移后复量一次仍＝0）⇒ 票面/复验件引用的 `:192`／`:196`／`:218` 三枚行号在工作树与 blob 上逐字相同（两把尺并排，读数⛔ 假设） |
| 车道在飞写者自查（尺＝`git log --oneline -3 --name-only -- internal/audio`） | 最近三枚碰 `internal/audio` 的全是**票 300 那一程**（`9a442923`／`64ceaacd`／`028529fc`），⛔ 与本票冲突的在飞写者；编排者这一笔自陈"编队两枚＝`306-r1` 写腿 ∥ `305-a1` 只读"⇒ `305-a1` 只写 `.scratch/**`，⛔ 抢 `internal/audio` 编译面 |

读数件＝`logs/00-anchor.txt`（每把尺自己落一行 `rc=N`；漂移那一节在文件后半段 `## drift readings`，同一枚件、⛔ 二次生成）。

## 2. 台面（整棵树的全貌，⛔ 本腿的活但⛔ 动、⛔ 暂存、⛔ 抱怨）

- `git status --porcelain`（无 pathspec 那把尺）＝**828 行**：别人的脏面本来就在，逐字可见 `.gitignore` 被改、`design/**` 大量 `D`/`M`、别家的 `.scratch/wisp/probes/**`、票 303 的文件被改。
- ⇒ 本腿**每一笔** commit 的显式 pathspec 只圈：`internal/audio/**` ＋ `.scratch/wisp/probes/306/r1/**` ＋ 票 306 文件（只在末尾追加一条 Progress log bullet）。
- 那 828 行本腿一枚⛔ 碰。逐笔越界自查见 `40-gates.md`。

## 3. 进程闸门自陈（派单第 3 节，本腿逐条怎么走的）

1. **commit-first**＝本件是起手锚跑完后落的第一笔，之前只跑了 `rev-parse`／`status`／`cmp`／`mkdir` 这四把短尺，⛔ 长跑命令。
2. **只 commit、⛔ push**（本腿全程零 push）。
3. commit 带**显式 pathspec**、写在 `$( … )` **之外**；含中文／反引号／`$` 的提交消息一律先 `Write` 成 `logs/msg-*.txt` 再 `git commit -F`；⛔ `git add -A`/`.`、⛔ `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`/`--no-verify`；⛔ 在仓库目录内建 worktree／checkout（"改前那发"台面＝`git archive 5480434f | tar -x` 到 `C:/Users/swq/tmp/` 下的仓外导出树，⛔ `a81c2980`——见 §1 漂移那一行，`internal/` 两枚锚点逐字相同，取盘上最新那枚）。
4. **退码一律纯重定向** `命令 > 件 2>&1; rc=$?`；⛔ `| tee 件; echo $?`（那取的是 tee 的退码）。每把门禁件自己落一行 `rc=N`。
5. 证据件⛔ 叫 `.out`（根 `.gitignore` 第 8 行是全仓 `*.out`）；本腿件名一律 `.txt`／`.md`，`cover` 档用 `.cover`／`.prof` 之类非 `.out` 后缀。
6. 临时件**只建不删**（含仓外导出树，⛔ 入库）。
7. 大输出先落文件再 grep（`go test -v` 整包日志、d22scan 输出都走文件）。
8. **⛔ 自勾任何 `AC` 框、⛔ 写"本格闭合"判语**；本腿只交读数＋凭据＋做到／⛔ 做到的名册；票面只在末尾追加一条 Progress log bullet（append-only）。
9. 目标调用数 ≤ 60 次。

## 4. 本腿要交的五格（照派单，⛔ 新判据）

- `AC#1`＝乙形夹具：同包 `_test.go` 里一份 40 字节 extensible（`tag=0xFFFE`／`cbSize=22`／`fmt` size=40）＋ float32（`bits=32`、`SubFormat.Data1` 低字＝`3` 在偏移 24）的 wav 面，断言**解析出来的值**（rate／chans／逐枚样本／长度），⛔ 断"某分支被调用过"。
- `AC#2`／`AC#2b`／`AC#1` 的牙＝**三对**成对突变凭据（`:192` `0xFFFE`→`0xFFFD`、`:196` `body+24`→`body+26`、`:218` `3`→`2`），每对＝改⇒目标用例当场红／`cmp` 逐字节还原⇒复绿，红句逐字引。
- `AC#3`＝`wave_format_float_300_windows_test.go:38-41` 那句注释改成实话（⛔ 动断言、⛔ 动期望值、⛔ 删 `waveFormatPCM`、⛔ 顺手改别的注释）。
- `AC#4`＝门禁＋越界（见 `40-gates.md`）。

## 5. 硬边界承认（本腿确认收到并逐条守住）

⛔ 把 `wavinjector.go` 改成引用 `waveFormatExt`/`waveFormatFloat`/`waveFormatPCM`（编排者仓外实测＝`GOOS=linux go build ./internal/audio/` rc=1，逐字两句 `undefined:`；本腿⛔ 重走这一形，只在 `40-gates.md` 复跑中立面的 rc=0）。
⛔ 动 `wavinjector.go` 任何产码语义／⛔ 动 `wasapi_windows.go:378`（R2＝票 300 形 C5，需 owner 另批）／⛔ 放宽票 300 `AC#6` 那枚钉的形状／⛔ 新二进制 `.wav` 夹具／⛔ 给新测试文件加 `//go:build windows`／⛔ 顺手补 injector 其它分支（含 `:194` 的 `size<40` 错误支与裸 `tag=3/bits=32` 非 extensible 面，本腿**刻意⛔ 装**，理由见 `10-fixture.md`）／⛔ SLO 阈值·golden·`thresholds.go`·`allowlist.txt`·D43·C1–C32·D1–D47·三枚冻结件·`PLAN.md`·`frontend/**`·`design/**`／⛔ 真窗真机／⛔ 跑 `cmd/wisp` 或整包 `go test ./...`。
