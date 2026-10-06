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
| 100 | `100-appcontainer-landing-ticket-gated-by-trigger.md` | 6114e3df 09-21 | 待填 | 待填 | 待填 |
| 103 | `103-seal-seam-is-unguarded-and-removeunlinked-follows-junctions.md` | a64d06fb 09-21 | 待填 | 待填 | 待填 |
| 107 | `107-ticket102-allowlist-judgment-is-red-on-posix.md` | b878b30f 09-21 | 待填 | 待填 | 待填 |
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
| 仍成立 | 待填 | — |
| 已失效／已被别人做掉 | 待填 | — |
| 差翻勾 | 待填 | — |
| 量不到 | 待填 | — |
| 合计 | 13 | ✓ 与 §1 名册对上 |
