# pool-validity-1 / 批三 —— 零勾开放票有效性普查（号段 100–149）

> 只读普查腿。取数时刻 **2026-10-05 22:3x +08**，起手 HEAD `f73558c0`（2026-10-05 22:28 +0800），分支 `dev`。
> 本腿**零 go 命令**（并行腿 `268-v2` 正在 cmd/wisp 跑定向用例与突变）；只用 `ls`/`grep`/`wc`/`find`/`git log`/Read。
> 票面零改动：一枚勾框没翻、一个票名没改、`Status:` 行一字没动、`docs/` 与产码没碰。
> 唯一写入路径＝`.scratch/wisp/probes/pool-validity/1/**`。⛔ 零读零引 `frontend/**`、`design/**`。
> 分桶参照（叠在其桶号之上，不替代）：`rate-census-4a`（1–49）／`-4b`（50–99）／`-3b`（100–149）／`-3a`（150–199）／`-2a`（≥200），皆为 **10-02/10-03 的读数**。

## §0 三档尺（本腿唯一的判法）

| 档 | 判据 | 复法 |
|---|---|---|
| **仍成立** | 票面点名的那枚缺陷今天在树里还在 | 把它 §现量／AC 里**最硬的一条静态断言**自己跑一遍：`grep -n <产码字面或函数名> <它指的文件>`，或读那枚行号现在写什么 |
| **已失效／已被别人做掉** | 缺陷不在了，或被别的票的 commit 顺手填了 | 给 commit 号＋`git log --name-only` 命中；缺陷字面在今树 0 命中 |
| **量不到** | 要靠运行期行为／真开窗／真机硬件／前端面才能判 | 归口写清缺哪一行读数；⛔ 不许由"grep 命中"外推成"这功能能用" |

★ `上次动过` 列＝`git log -1 --format='%h %ad' --date=format:'%m-%d' -- <票文件>`（编排者排"先打哪几张"的唯一依据）。
⚠ 诚实边界：三档里"仍成立"是**静态尺**复认，运行期的一律落第三档；本腿不翻任何 AC 框、不裁修法、不裁形。

## §1 名册（13 枚）

| 票号 | 票文件 | 上次动过 | 判定 | 复法（跑的那把尺） | 读数 |
|---|---|---|---|---|---|
| 100 | `100-appcontainer-landing-ticket-gated-by-trigger.md` | 6114e3df 09-21 | 未判 | — | — |
| 103 | `103-seal-seam-is-unguarded-and-removeunlinked-follows-junctions.md` | a64d06fb 09-21 | 未判 | — | — |
| 107 | `107-ticket102-allowlist-judgment-is-red-on-posix.md` | b878b30f 09-21 | 未判 | — | — |
| 108 | `108-103s-seal-guard-is-bypassable-three-ways.md` | 3f00217e 09-21 | 未判 | — | — |
| 109 | `109-models-have-a-cross-account-write-window-after-ensure.md` | 52ad9335 09-23 | 未判 | — | — |
| 111 | `111-ci-tests-20-of-33-packages.md` | c9907bb2 09-21 | 未判 | — | — |
| 112 | `112-winsec-step-goes-red-on-the-ci-runner-first-run.md` | e3fb933f 09-21 | 未判 | — | — |
| 120 | `120-posix-seal-is-check-then-act-154-of-4000-retries-hit-the-window.md` | 76fc5b07 09-21 | 未判 | — | — |
| 122 | `122-clear-the-staticcheck-backlog-after-85a.md` | 048a9e4c 09-23 | 未判 | — | — |
| 123 | `123-cmd-wisp-cli-tests-assume-a-human-approver-on-ci.md` | c9907bb2 09-21 | 未判 | — | — |
| 132 | `132-log-files-are-never-sealed-sealdirdir-has-zero-production-callers.md` | d5a59d66 09-29 | 未判 | — | — |
| 140 | `140-job-level-wisp-env-test-makes-the-ask-the-os-hardening-legs-never-get-consulted-on-the-windows-runner.md` | d171b427 09-25 | 未判 | — | — |
| 148 | `148-clonedecision-only-guarded-one-door-the-sibling-replay-outlet-still-hands-out-a-shallow-value-copy-and-what-the-panel-should-receive-is-unwritten.md` | d08999a4 09-25 | 未判 | — | — |

## §2 本批计数

仍成立 **未判** ／ 已失效 **未判** ／ 量不到 **未判** ／ 合计 13
