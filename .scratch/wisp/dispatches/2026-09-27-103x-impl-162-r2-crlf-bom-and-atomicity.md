# 派单 162-r2（写码位，**两格＋一句注释**）＝票 162 的 AC#3（行尾/BOM 正反两向）＋ AC#4（原子性反面实验），外加编排者裁给你的一句不完整清单文案

- 派单时刻：09-27 10:3x｜前一程 **162-r1 干净交件**（`0cb8864 a9d0576 7aaff7a 0e36753`，AC#1＋AC#2；我这边 `go test ./internal/tools/ -run TestFSEdit -count=1` ⇒ ok 0.183s）。
- 票面＝`.scratch/wisp/issues/162-…-rewriting-whole-files.md`。**必读最末两格**（09:4x 派单存档 ＋ 10:2x 我代落的 r1 交件与三裁）。
- ⚠ **预算硬顶 ≤60 次工具调用**：到点停手→写 `next=`→未提交增量用**显式 pathspec**提交。**每裁完一格 commit 一次。**
- ⚠ **AC#5（风险档路由）与 AC#6（契约轴复算）不是你的**——留给 r3。尤其**不许动 `internal/risk/**`**。

## 0. 起手四件
`date`｜`git rev-parse HEAD`｜`git status --porcelain -- internal/tools/ docs/evidence/s1/162-*.md`（**必须空**）｜`Read` 票面把 AC#3／AC#4 连行号抄进证据件。

## 1. 第①格＝AC#3（行尾与 BOM，**正反两向**）
起点（r1 现量，你要复跑）：`internal/tools/fs_edit.go:136` 是 `content := string(raw)`＝**纯字节匹配** ⇒ "LF 文件被改成 CRLF"今天**不可能发生**；但这两枚**还没有用例**：
- **带 BOM 的文件改完 BOM 还在**（正向）。
- **CRLF 文件里 `old` 写成 LF ⇒ 0 命中**（这一枚今天**会**响亮失败＝它响，但要看清它是"正确拒绝"还是"根本读不到"）；反过来 **CRLF 文件里 `old` 也写成 CRLF ⇒ 必须命中并落盘**。
⚠ 票面 `:27` 那条"行尾归一（内部统一成 LF、写完还原）"是**放宽匹配**——本格**只做取证、不做实现**：把"不实现它，哪些真写法会被 0 命中拒掉"量成读数并登记，**要不要做由我裁**。

## 2. 第②格＝AC#4（原子性反面实验）
r1 现读：`fs.edit` 复用的就是 `TestAtomicWriteKillsMidWrite` 那套缝（`FSDeps{Hooks:…}` 的 `Hooks.Kill`）。⇒ 把同一枚缝装进 `fsEditBridge`（**用原有的 helper，不许新造测试脚手架**），在 `write:` 边界杀一次，断言：
(a) 目标文件字节＝**改前原文**（不是半截）；(b) 目录里**不留** `.wisp-tmp-*` 残件**或**残件被点名（两种都行，但要说清今天实际是哪种）；(c) 答承重那一问：**摘掉 temp＋rename 那一味，是否存在一发损坏从此看不见？**
⚠ 判据不许写成"我看了代码觉得安全"。

## 3. 第③件（裁给你的那句注释，**只改文字**）
`internal/tools/fs.go:22-25` 那句 `The write half (fs.write / fs.trash / fs.move …)` 在 `fs.edit` 落地后**是不完全清单**（r1 按最窄读法故意没改并报回，我裁给 r2）。⇒ **允许你把这一句改准**（只改注释文字、**函数行为一字节不许动**），并答一句："删掉这句更正，哪一发会重新误导？"

## 4. 写面（**逐枚列名**——上一程我写"只许改那一句"没料到注册链的连带，这是我的缺陷，这次列全）
`internal/tools/fs_edit.go`｜`internal/tools/fs_edit_test.go`｜`internal/tools/fswrite_silentloss_ac1_test.go`｜新测试文件 `internal/tools/fs_edit_ac34_test.go`（建议单独一枚，别挤进 r1 的文件）｜`internal/tools/fs.go`（**只许 `:22-25` 那一句注释**）｜`docs/evidence/s1/162-fs-edit-r2.md`｜`.scratch/wisp/probes/162/**`。
**禁改**：`internal/risk/**`、`internal/config/**`、`internal/agent/**`、`internal/panel/**`、`cmd/**`、`docs/PLAN.md`、`docs/specs/**`、`tools/d22scan/**`（含禁令射程）、`allowlist.txt`、thresholds.go／golden、`.github/workflows/**`、`frontend/**`、`design/**`、`docs/reports/**`、票面（进度由我代落）、别人的票面。
⚠ **若两格任一必须落在禁改名单里 ⇒ 停手报回**（哪个文件、哪一行、非动不可的理由、有没有不动的形状）。**owner 这次给 162 的范围里没有 `internal/risk`、`internal/config`。**
**现场理由**：`design/**` 有 owner 未提交的删除；`probes/152/my152.py`（` M`）与 `docs/evidence/s1/152-…-accept-r1.md`（3 行未提交）是停下来的半件——不提交、不还原、不补完 ⇒ `git add -A`／`git add .`／`git commit -a` 一律禁止。

## 5. 门禁（**两格全部落盘之后**再跑，顺序别倒）
`go test ./internal/tools/ -count=1`（⚠ **四数之外必须报名册差集**：r1 的基线是"顶层 run=86 pass=86 fail=0 skip=0"，你新增的用例逐枚点名，**"基线有而现在无"必须 0 枚**；一条 panic 会吞同包其余读数）｜`go vet ./internal/tools/`｜`bash tools/d22scan/runtests.sh -C tools/d22scan ./...`（基线 `PASS=34 FAIL=0 SKIP=0`、`RUN=76`）｜`sh scripts/d22scan.sh` 预期 rc=0｜`gofumpt --version` 现跑｜甲形（已跟踪集合）必须 0 行、乙形每一行可归因。

## 6. 仪器坑（实测过）
`grep -c` 命中 0 ⇒ rc=1 吃掉 `&&`｜`git ls-files -Z` 不存在（非法开关仍退 0、产出恒真的"空"）｜**`git grep` 的行是 4 段 `<commit>:<path>:<line>:<text>`**｜Windows 反斜杠路径｜截断输出不是全表｜`-overlay` 与 `-cover*` 不许合跑｜**本机 `staticcheck` 读不了 go1.27 的 export data**（会产"干净的绿"，别拿它当证据）｜`cmd/wisp` 的用例要先有 `third_party\sherpa-onnx` 在 PATH 才不会 `0xc0000135` 加载期崩（既有账，见 `HANDOVER` §4.0s）。

## 7. Git 与噪声
只 commit 不 push；显式 pathspec ＋ **加引号的 heredoc**；提交前 `git diff --cached --name-only`（别人路径＝共享索引正常噪声），**提交后 `git show --name-only <自己号>` 出现别人路径＝停手报回**；禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；**一律不许跑删除命令**（每次 `rm` 都会在真人机器上弹授权窗）；台件**只建不删**、落盘前 gofumpt 干净。

## 8. 交件报告（顺序固定）
1. step-0 四件；2. **本程没测什么**；3. AC#3 两枚各正反两向＋"不实现归一会拒掉哪些真写法"的读数；4. AC#4 三问 (a)(b)(c)；5. 第③件：那句注释的 diff ＋ 删掉它哪一发会重新误导；6. 门禁四数＋**名册差集**；7. **有没有动过禁改名单里的文件**（应为"无"）；8. 被拒／没成功的调用（取数前还是后）；9. **有没有跑过删除命令**；10. 伪授权两栏计数（各带出处＝工具名＋命令前 40 字）；11. 凭据值零抄录；12. 档位（四档）；13. `next=`（给 r3 的 AC#5/AC#6 各一句起点）。
⚠ 我这段话里每条前提（含"`fs_edit.go:136` 是纯字节匹配""基线 86 枚""`Hooks.Kill` 那套缝可复用"）**都是未验证断言**：不符就报回并继续做做得动的部分，**不许为迁就我去改判据、改测试或放宽断言**。
