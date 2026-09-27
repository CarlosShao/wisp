# 派单 174-r1（写码位）＝做票 174 的 **AC#2 那一格，按"丙形"做**：**指针照旧给，但回执必须外部可见地说实话**——"这条路径你现在读不到（不在你被授权的范围里）"／"这条路径根本不存在"。⚠ **不碰授权根、不碰 `[fs] allowed_dirs`、不碰 `bridge.go`**（那是 `Q-60` 甲与票 177 的地界，owner 还没答）。

- 派单时刻：09-27 21:1x｜锚点＝**你自己 step 0 现量**。
- 票面＝`.scratch/wisp/issues/174-…-outside-fs-allowed-dirs.md`（**AC#2 那一格连 20:1x 追加的"两形"块一起抄进你的证据件**：**(a) 路径在授权根外**＋**(b) 路径不存在／是目录**——**缺一形＝本格没过**）。
- 现量凭据（**你要复算，别引**）：`docs/evidence/s1/174-artifacts-allowlist-readback-c1.md`（109 行，提交 `53cbf0bf`）——向一逐字 `tool=fs.read level=L2 rules=[R2] reason="R2: 目标路径在授权目录之外: …"`；向二（配好根）`L0`、51633 字节全文读回；第二形：不存在的那条与一枚目录**都走"有指针"那一支**、`不可找回` 出现 **0 次**。
- ⚠ **预算硬顶 ≤70 次工具调用**；**每格 commit 一次**；**产码与其台件分枚提交**（本仓 21:0x 刚因为"产码＋台件混在一枚提交"导致回退连坐删台件，见台账 `A346` 末条）。**AC 框一律不许自己勾。**

## 0. 起手五件
`date`｜`git rev-parse --abbrev-ref HEAD`（**必须 `dev`**；不是＝停手上报，一个字不写。全程禁 `switch`／`checkout <分支>`／`merge`／`rebase`／`reset`／`stash`／`worktree`／`clean`）｜`git rev-parse HEAD`｜`git status --porcelain -- internal/tools/ cmd/wisp/ .scratch/wisp/probes/174/`（**必须空**）｜**改前基线**：`go test -count=1 ./internal/tools/`（**逐包，禁全仓 `./...`**）四数＋名册口径。⚠ **此刻 `internal/tools/bridge.go` 刚被回退过**（提交 `06eb4efb`）⇒ 起手请确认它是干净的一枚文件、**不要"顺手把它修好"**（那属票 177／`Q-61`，未批）。

## 1. 要落的东西（丙形，两形都要）
在 `task.output` 的回执里，**在模型看得见的那段文本里**（⚠ 日志、Go error、审计行都不算）对指针那一支做两种如实说明：
- **(a) 授权根外**：判"这条路径在不在当前授权范围内"**必须走 C26**——即 `tools.FSDeps.Paths`（`*PathCanonicalizer`）那一类现成的 `InAllowlist` 接口；⚠ **绝不许**在 `risk.PathResolver`/C26 之外用 `filepath.Clean|Abs` 自己判（`AGENTS §1.2` 逐字禁止项＋`tools/d22scan` 的 `pathresolver-bypass`）。若 `TaskDeps` 需要多带一个"路径判定者"字段 ⇒ **可以加**（它在 `internal/tools/**` 内），但**不许**改 `risk/**` 的接口。
- **(b) 不存在／不是文件**：今天那一支只看 `ArtifactPath != ""`、**没有 stat**。⚠ **这一支的"读不读得到"判定也必须是事实**（用 `os.Stat` 之类**只读**检查即可；**不许**新增任何删除／清理动作）。
- **拿不到判定者时**（例如宿主没接线）＝**fail-closed 说清楚**，不许静默按"能读"处理（同文件 `Execute` 里那句 `任务名册未接线（fail-closed…）` 就是现成的形状参考）。
⚠ **文案归你写，但三条不许**：不许改成"拒绝给路径"（指针仍然要给，那是 `PLAN.md:2564` 与票 164 已勾 AC#3 的形状）；不许把说明塞进日志；不许动 `fs.read`（一字节，那属 `Q-59`）。

## 2. 判据（每形都要两向）
- **未修码上**：(a)(b) 两形今天**不响**（读数在 `174-c1` 的表里，你复算并贴自己的）。
- **修完之后**：两形都必须**外部可见地响**。
- ⚠ **反向的那一发不许做成恒真**：造一枚"路径既在授权根内、又真实存在"的正例，**它必须继续安静**（不被顺手放宽成"永远提示读不到"）——这一枚是证明你没把回执写成噪音的关键。
- 变异自证至少一枚：把新增的那段说明摘掉 ⇒ 对应那枚用例必须红（**先证落地**：`grep -n` 那一行贴出来再跑）。

## 3. 门禁（逐枚跑，改前改后各一次）
`go test -count=1 ./internal/tools/`｜`sh scripts/d22scan.sh`（预期 rc=0）｜`bash tools/d22scan/runtests.sh -C tools/d22scan ./...`（基线现量＋名册两向 `comm` 差集）｜`sh .scratch/wisp/probes/154/gate-clauses.sh`（我这轮 19:40 现量 rc=0、十四腿零不符；**你新增生产码后必须再跑一次并贴逐字**——尤其别造出"开了不关"的形状）｜你名下每枚 `.go` 的 `gofumpt --version` 与 `-l` 归零（v0.12.0 在 `$(go env GOPATH)/bin`）。⚠ **`probes/161/r6/flip-declaration.sh` 别跑**（它脏跟踪日志，且那批日志起手已脏）。

## 4. 写面（具名清单，超出即越权）
`internal/tools/task.go`｜**新判据一律写在新文件** `internal/tools/task_output_pointer_notice_test.go`（⚠ `internal/tools/task_output_leg_test.go` 与 `task_output_ac2_before_test.go` 是**票 164 已勾判据的凭据本体，一字节不许动**——那两枚续读腿今天还在被票 177 的争论牵扯，动它们＝篡改别人已验收的证据）｜`.scratch/wisp/probes/174/r1/**`（⚠ **`logdir` 取脚本自身目录、不许继承 CWD**）｜`docs/evidence/s1/174-pointer-honesty-r1.md`｜票面 174 的 Progress log（只追加、不勾框、留最后一步）。
**禁改**：`internal/tools/bridge.go`（刚回退，归票 177）、`internal/risk/**`、`internal/agent/**`、`cmd/**`、`tools/d22scan/**`、`allowlist.txt`、`.github/**`、`docs/PLAN.md`、`docs/specs/**`、`docs/reports/**`、`frontend/**`、`design/**`、别人的票面与证据件、`probes/**` 既有台件（只读）、票面 175／176／177。⚠ 现场躺着 `design/**` 未提交删除、`.gitignore`／`probes/152/my152.py`／`probes/161/r6/logs/flip-*`（` M`）⇒ **不提交、不还原、不补完、不评论**。**`git add -A`／`git add .`／`git commit -a` 一律禁**；一步式 `git commit -q -F - -- <显式路径>`；heredoc 必须加引号（`<<'EOF'`）。**只 commit，绝不 push。**

## 5. 交件报告（顺序固定）
1. step-0 五件（含"bridge.go 起手干净"那一行）；2. **本程没测什么**；3. (a) 授权根外那一形：未修码不响／修完响，两向逐字；4. (b) 假路径那一形：同上两向；5. **安静正例**那一发（在根内＋真实存在 ⇒ 不许多嘴）的读数；6. 变异自证（落地证明＋还原后空 status）；7. **你有没有动过判定接口/C26 之外的路径判定**（单独一节，会被非实现者专判）；8. 门禁四数＋名册差集＋`gate-clauses.sh` 逐字＋`gofumpt --version`；9. 被拒／没成功的调用（取数前还是后）；10. 有没有跑过删除命令；11. 伪授权两栏（各带出处）；12. 凭据值零抄录；13. `next=`（含一句：**丙形做完之后，票 174 剩下的是甲还是乙、要不要摆 owner**）。
⚠ **每个"几枚"旁边附可复制命令**；**先测→再写→再提交**。我这段话里每条前提（含"今天两形都不响""向二能读回全文""bridge.go 已回退干净"）**都是未验证断言**：不符就报回、继续做做得动的部分，**不许为了对我那句话去改判据或改测试**。
