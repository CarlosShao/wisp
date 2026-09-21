# 93 — portable 测试步是裸 `go test` ⇒ **SKIP 记成 ok**：`TestSyncRegistryProbeLive` 两侧都跳，CI 却全绿

**Status:** open（2026-09-21 17:1x 编排者建，来源=票 82 交件时**自己点名**的残留；我没让它顺手修，因为那超出票 82 的界）
**Type:** 门禁完整性（票 71 家族："门没跑"和"跑了没问题"在 CI 输出上长得一模一样）
**Blocks:** nothing（但它在**削弱** `test-core` 全绿这句话的含义）· **Blocked by:** nothing
**Packages:** `.github/workflows/ci.yml`（`test-core` 的 "Portable package tests" 步）、
              `scripts/` 下一个可复用 runner（**先读 `tools/d22scan/runtests.sh`，它已经会拒 SKIP**）。
              **禁改**：任何测试的**断言**、任何阈值、`docs/PLAN.md`、`docs/specs/*.md`、
              `tools/d22scan/**` 与 `allowlist.txt`（那是编排者的门禁面）。

## 现场（可直接复现）

`internal/risk` 的 `TestSyncRegistryProbeLive` 在 **Windows 与 ubuntu 两侧都 `--- SKIP`**
（Windows：这台机器 HKCU 无 `UserFolder`；POSIX：压根没有注册表）。CI 的 portable 步是
`go test <16 个包> -count=1`（`ci.yml` 的 "Portable package tests"），**裸 go test 对 SKIP 的记账是 `ok`**
⇒ 这条用例在 CI 上**从来没有产出过一个结论**，而输出看起来和"它通过了"完全相同。
票 82 的代理在 `-count=2` 下量到 **2 行 SKIP，同名 1 个**。

⚠ 同一仓里已经有**两套仪器对同一件事一严一松**：`tools/d22scan/runtests.sh` **拒 SKIP**，
portable 步不拒。这就是本票的缺陷本体——**不是那条用例该不该跳，是"跳了没人记账"**。

## AC（1:1，裁决表 `docs/evidence/s1/93-*.md` 由验收方出）

- [ ] **AC#1** 先做**全量点名**：把 CI portable 清单里 16 个包在 **Windows 与 ubuntu 两侧**当前所有
      `--- SKIP` 逐条列出（包 / 测试名 / file:line / 跳过原因）。**判据是"清单是完整的"**：
      你的取数命令必须能证明它抓到了全部输出（贴 `grep -c -- "--- SKIP"` 与总行数的对照），
      不许只贴 `internal/risk` 那一条。
- [ ] **AC#2** 落一个**真会红**的机制：portable 步遇到任何未被显式记账的 SKIP ⇒ 该步 rc≠0。
      ⚠ **不许**用"往 allowlist 里加一行"把它糊过去——allowlist 只会把问题从"红"变成"没人看的白名单"。
      若某条 SKIP 是**合法的**（平台 API 就是不存在），修法必须是把那条判据**搬到平台层**
      （`//go:build`，票 82/81 已经做过两次，照做），让它在不适用的平台上**不参评**，
      而不是在适用的平台上"跳过并算过"。
- [ ] **AC#3** 双向变异：(i) 把机制改成 no-op ⇒ **必须有既有用例红**（不是本票新写的）；
      (ii) 往快照里**种一条 `t.Skip`** ⇒ portable 步必须红并点名它。锚点=承载行为那一行，
      同链 grep 证落地，**编译失败不算变异**，还原后 `git diff --quiet` 证干净。
- [ ] **AC#4** CI 上有**真实 run id + 步级结论**（`gh run view --job` 读**step**，不读 job status；
      并区分 failure / cancelled（`total_count=0` 不是样本）/ 未跑完）。⚠ 代理**不 push** ⇒
      这条若要 push 才能拿到，就在票面明写"欠编排者 push 后复跑"，**不许用"本地跑过了"替代**。
- [ ] **AC#5** 门禁（只跑自己碰的范围）：`gofmt -l` 空、改动脚本 `bash -n` / `go vet` rc=0、
      本地按新 runner 跑一遍并贴出**四数**（`=== RUN` / PASS / FAIL / **SKIP 必须为 0 或被逐条点名**）。

## Rules（本仓固定）

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc）；提交前 `git diff --cached --name-only`；
禁 `git add -A`/`.`、`--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34）；**不 push**；
不在仓内建 worktree（A38④，快照 `git archive <sha> | tar -x -C /tmp/<带会话后缀>`）；
不跑整仓门禁（共树：`internal/winsec/**` 票 89、`frontend/`+`internal/panel/` 票 77 有人在写，别碰）；
票面 append-only，改行前先读，在标题前插段落要把标题重抄进 `new_string` 且 `git diff --numstat` 删除列为 0；
四种假绿逐条点名；数字不达标就 FAIL 附数字，不许调阈值、不许挑运气那次、多样本全报。
15 次工具调用内交回第一枚 checkpoint；每次提交同步 Status + 勾框 + `next=`；接近轮数上限主动收尾留断点。

## Progress log（append-only）

- 2026-09-21 17:1x（编排者）：建票。来源是票 82 报告的两条残留之一："`TestSyncRegistryProbeLive` 两侧都 SKIP，
  而 CI 的 portable 步是裸 `go test` ⇒ SKIP 记成 ok（只有 `tools/d22scan/runtests.sh` 拒 SKIP）"。
  我把它单独开票而**不塞进票 82**，因为票 82 的界是那 8 条红；顺手改它=扩界（它自己也是这么判的）。
  ⚠ 判据里我最在意的一条是 AC#2 那句"**不许用 allowlist 糊过去**"：这票的修法如果让白名单变长，
  那它就是把"红"改成了"没人读的清单"，等于没修。
  next= 派单（它不依赖任何在飞的票，可以插队；但**先等 `internal/panel`/`frontend/` 那两路收敛**，
  因为它要动 `ci.yml`，而票 77 也在等 CI 结论 ⇒ 同文件风险）。
