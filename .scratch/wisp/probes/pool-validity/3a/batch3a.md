# pool-validity-3a / 批3a —— 零勾开放票有效性普查（号段 100–149，13 枚）

> 只读普查腿 `pool-validity-3a`，接 `pool-validity-1`／`pool-validity-2`（两枚同职腿都死在"扫全池"）。
> 起手 HEAD `4b15869d`（2026-10-06 10:12 +0800，`## A632`）。⚠ 换 HEAD 要重量。
> 本腿**零 go 命令**（并行腿 `231-v1` 此刻真在 `cmd/wisp` 取整包终态）。
> 工具全集＝`git log`／`git grep HEAD`／`git show HEAD:<path>`／`git ls-tree`／`ls`／`grep`／`wc`／Read。
> 票面零改动：零翻勾、零改名、零 `Status:` 改写、零产码改动、零台账写入。
> 唯一写入路径＝`.scratch/wisp/probes/pool-validity/3a/**`（新建）。
> ⛔ `pool-validity/1/**` 与 `pool-validity/2/**` 只读只引、不采信其结论；
> ⛔ `pool-validity/2/batch3.md`（前一枚死腿的空骨架，13 枚全"未判"）**本腿一字不动**，本件另立，防"谁判的"被洗。
> ⛔ 零读零引 `frontend/**`、`design/**`、`.gitignore`、`.scratch/wisp/probes/{231,269,111}/**`。
> ⚠ 写面纪律：`cmd/wisp/**`、`scripts/**`、`.github/**` 只按**文件名／行数**取，⛔ 不读内容（并行腿地界）。

## §0 四档尺（本腿判法，与 pv2 同构但口径自立）

| 档 | 判据 | 复法（必须留命令原文＋读数） |
|---|---|---|
| **仍成立** | 票面点名的那枚缺陷今天还在树上 | 把它 §现场／AC 里**最硬的一条**自己跑尺：`git --no-pager grep -n '<那句字面或符号>' HEAD -- internal/ cmd/` 的命中数与 file:line |
| **已失效／已被别人做掉** | 缺陷不在了 | ⛔ 硬门：一枚 commit 号＋`git log --name-only` 命中，**或**今树 0 命中的复算尺；拿不出就不许写这一档 |
| **差翻勾** | 缺陷不在了，但凭据是"票内注记／台账 `A##`／一枚非实现者裁决表" | 凭据种类写进复法列；**不等于**本腿裁它可结案 |
| **量不到** | 判据落运行期行为／真开窗／真 seal／DPAPI／多账号登录／`frontend/**`·`design/**`（本腿禁读） | 归口写清缺哪一行读数；⛔ 不许由 grep 命中外推成"这功能能用" |

★ `上次动过` 列＝`git log -1 --format='%h %ad' --date=format:'%m-%d' -- <票文件>`（本腿逐枚实跑）。
★ 本族（印章／封根／门禁）预登记的**一处陷阱**：真 seal、DPAPI、多账号登录的运行期行为静态判不到 ⇒
  归〔量不到〕；"代码里那个函数还在"只是〔仍成立〕里最弱的一种凭据，**备注列一律标弱**。

## §1 名册与判定（13 枚）

| 票号 | 票文件 | 上次动过 | 判定 | 复法（命令原文） | 读数 |
|---|---|---|---|---|---|
| 100 | `100-appcontainer-landing-ticket-gated-by-trigger.md` | 6114e3df 09-21 | **仍成立**（门仍未开） | 票面 AC#0 自己那把尺逐字跑：`git --no-pager grep -n "exec\.Command" HEAD -- internal/tools/ internal/plugin/ \| grep -v '_test\.go'` → **rc=1／0 命中**；加宽一把：`git --no-pager grep -n "exec\.CommandContext\|exec\.Cmd\|\"os/exec\"" HEAD -- internal/tools/ internal/plugin/ internal/agent/ internal/proc/ \| grep -v '_test\.go'` → **1 命中**＝`internal/proc/jobscope_windows.go:110 StartInJob(cmd *exec.Cmd)`（收一个别人的 `*exec.Cmd`，自己不建命令）；`git --no-pager grep -n "shell\.exec" HEAD -- internal/tools/ \| grep -v '_test\.go'` → **0 命中**；载体尺：`git --no-pager grep -iln "appcontainer\|CreateAppContainerProfile\|SECURITY_CAPABILITY" HEAD -- internal/ cmd/ tools/` → **0 枚文件** | 门（AC#0）**没开**：生产面 `exec.Command` 非测试命中 **0**，`shell.exec` 工具**未注册**（票 52 同读），AppContainer 全仓**零字符**（连注释都没有）⇒ AC#1–AC#6 六格**无载体**、票面 `next= 不派单` 今天仍然成立。⚠ 票面 AC#0 那句"今天的全仓非测试 `exec.Command` 只有 `mockllm`/SLO 自测/`doctor`"已**过期**：今树非测试命中只有 `internal/llm/adaptertest/mockllm.go` 与 `tools/d22scan/gitignore.go` 两处（`git --no-pager grep -ln "exec\.Command" HEAD -- internal/ tools/ \| grep -v _test\.go` → 2），`doctor`/SLO 那两处不在 `internal/`＋`tools/` 的读得范围内（在禁读的 `cmd/wisp`／只列名的 `scripts`）⇒ 不影响判定（门仍关），但票面那句别当现价读。AC#1–#6 的判据（真起 AC 子进程、`icacls` SID 前后、D32 采样）**全落运行期＋真机** ⇒ 若门开，本腿只能量到"零载体"这一半 |
| 103 | `103-seal-seam-is-unguarded-and-removeunlinked-follows-junctions.md` | a64d06fb 09-21 | **差翻勾**（缺陷已不在，凭据＝票内注记＋裁决表＋它自己的 commit） | 票面两条读数的现树复算：①`git --no-pager grep -n "resolverProbeShapes\|refusing to install a path resolver\|the seam is single-use and one-way" HEAD -- internal/winsec/ \| grep -v '_test\.go'` → **6 命中**（`resolve.go:111/158/169/183/196/292`）；②`git --no-pager grep -n "firstLinkAncestor" HEAD -- internal/winsec/ \| grep -v '_test\.go'` → **3 命中**（`winsec.go:261/268/279`）；③用例载体：`git ls-tree -r --name-only HEAD -- internal/winsec internal/memory \| grep -E 'seam_guard\|junction_tripwire'` → **2 枚都在**；④commit 尺：`git --no-pager show --name-only --format='%h %ad %s' --date=format:'%m-%d' 0717bf29`（标题逐字 `fix(103,AC#1#2): 密封缝从此只能装一次…RemoveUnlinked 动手前先查祖先链是不是链接`）→ `--name-only`＝`internal/winsec/{resolve.go,winsec.go,winsec_other.go,winsec_windows.go}`；裁决表存在尺＝`git ls-tree -r --name-only HEAD -- docs/evidence/s1 \| grep -E '/103-'` → `docs/evidence/s1/103-adversarial-acceptance.md` | R-c 那条"没有任何守卫"的句子今树**读不出真**：`SetPathResolver`（`resolve.go:128`）现在有三道守卫＋**单向闩**（`:145-155` 拒 `nil` 解除，理由逐字点名 "which is ticket 103's probe P1b"），`RemoveUnlinked`（`winsec.go:255-266`）先 `IsAbs` 后 `firstLinkAncestor`。⇒ **两条缺陷都不是"今天还在"**。但**这活是 103 自己做的**（`0717bf29` 标题就叫 fix(103)），所以不写〔已失效／已被别人做掉〕，写〔差翻勾〕。**⛔ 不等于可结案**：①它被 `rejected-needs-fix` 卡在"守卫可被绕过"，续跑单＝票 108（本批判 108 的 P1b/P2/P3 也已闭合）；②AC#3 三向变异与 AC#4 的 `-count=2`/`go vet`/`d22scan` 四数是运行期 ⇒ 本腿量不到；③裁决表内容本腿禁读，只核到文件名与 commit 标题 |
| 107 | `107-ticket102-allowlist-judgment-is-red-on-posix.md` | b878b30f 09-21 | **已失效／已被别人做掉**（硬门过得去：commit 号＋`--name-only` 命中＋今树复算） | 硬门第一把：`git --no-pager log --format='%h %ad %s' --date=format:'%m-%d' -S "treeResolvedAsNamed" -- internal/tools/paths.go` → **`8e100950 09-21`**，标题逐字 `feat(92,AC#1-#7): 工作区选择 + composer 封套 + 前端输入框 v2…`；`git --no-pager show --name-only --format= 8e100950 \| grep internal/tools` → `internal/tools/paths.go`／`paths_workspace.go`／`paths_workspace_test.go`。硬门第二把（今树复算，票面点名的错法本身）：`git --no-pager grep -n "treeOnDisk\|treeResolvedAsNamed\|resolvedForm" HEAD -- internal/tools/paths.go` → **12 命中**，授权腿在 `:99 if !res.Resolved && !treeResolvedAsNamed(res.Canonical) {`、放行腿在 `:191`／`:220 rf, ok := resolvedForm(canonical)`；`git ls-tree -r --name-only HEAD -- internal/tools \| grep -i ticket107` → `paths_ticket107_portable_test.go` ＋ `paths_ticket107b_probes_test.go`（三枚探针**已进仓**，正是验收 `next=` 要求的"收进可跑用例"） | 票面点名的两件事今树**都不在**：①"`InAllowlist` 只看存在性"——`:99` 现在只认 `treeResolvedAsNamed`（`paths.go:294` 定义，注释逐字 "answers ticket 102's confirmation question"），`treeOnDisk`（`:275`）降级为"WEAKER question"且**授权腿零调用**（非测试调用点只剩 `paths_workspace.go:100`，票 92 的存在性腿，方向是收紧）；②"探针只在 `/tmp`、进不了 CI"——两枚 107 用例文件在 HEAD 名册里。⇒ 〔已失效〕。**⚠ 一处措辞要说死**：装进的是 `8e100950`＝**票 92 的 commit**，代码却是 `agent-ticket107b` 自己写的（票面 19:5x 编排者注记逐字"那枚 commit **少署了 107b**"，台账 `A80` 记修法被攻破那一轮）⇒ "被别人做掉"在这里＝**被别人的 commit 装进去的**，不是"别的票替它把活干完"。**⛔ 不等于零活**：AC#1 的 POSIX 复现/run-id 空位、AC#3"两侧都被真正执行过"、AC#5 四数全是 go 读数 ⇒ 本腿量不到，欠的是**以修复进树后的 HEAD 复验**（编排者 19:5x 那句 `next= 派复验`），不是产码 |
| 108 | `108-103s-seal-guard-is-bypassable-three-ways.md` | 3f00217e 09-21 | 待填 | 待填 | 待填 |
| 109 | `109-models-have-a-cross-account-write-window-after-ensure.md` | 52ad9335 09-23 | 待填 | 待填 | 待填 |
| 111 | `111-ci-tests-20-of-33-packages.md` | c9907bb2 09-21 | 待填 | 待填 | 待填 |
| 112 | `112-winsec-step-goes-red-on-the-ci-runner-first-run.md` | e3fb933f 09-21 | 待填 | 待填 | 待填 |
| 120 | `120-posix-seal-is-check-then-act-154-of-4000-retries-hit-the-window.md` | 76fc5b07 09-21 | 待填 | 待填 | 待填 |
| 122 | `122-clear-the-staticcheck-backlog-after-85a.md` | 048a9e4c 09-23 | 待填 | 待填 | 待填 |
| 123 | `123-cmd-wisp-cli-tests-assume-a-human-approver-on-ci.md` | c9907bb2 09-21 | 待填 | 待填 | 待填 |
| 132 | `132-log-files-are-never-sealed-sealdirdir-has-zero-production-callers.md` | d5a59d66 09-29 | 待填 | 待填 | 待填 |
| 140 | `140-job-level-wisp-env-test-makes-the-ask-the-os-hardening-legs-never-get-consulted-on-the-windows-runner.md` | d171b427 09-25 | 待填 | 待填 | 待填 |
| 148 | `148-clonedecision-only-guarded-one-door-the-sibling-replay-outlet-still-hands-out-a-shallow-value-copy-and-what-the-panel-should-receive-is-unwritten.md` | d08999a4 09-25 | 待填 | 待填 | 待填 |

## §2 本批计数

| 档 | 枚数 | 票号 |
|---|---|---|
| 仍成立 | 1 | 100 |
| 已失效／已被别人做掉 | 1 | 107 |
| 差翻勾 | 1 | 103 |
| 量不到 | 0 | — |
| 未判（余） | 10 | 108 109 111 112 120 122 123 132 140 148 |
| 合计 | 13 | ✓ 与 §1 名册对上 |
