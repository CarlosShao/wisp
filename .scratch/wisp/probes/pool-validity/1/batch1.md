# pool-validity-1 / 批一 —— 零勾开放票有效性普查（号段 1–49）

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

## §1 名册（30 枚）

| 票号 | 票文件 | 上次动过 | 判定 | 复法（跑的那把尺） | 读数 |
|---|---|---|---|---|---|
| 15 | `15-speech-engines-cer-harness.md` | da81a6ac 09-19 | 未判 | — | — |
| 16 | `16-s2-acceptance.md` | 144151ae 09-19 | 未判 | — | — |
| 21 | `21-approval-gates-minimal.md` | 4a75482b 09-20 | 未判 | — | — |
| 22 | `22-web-tools-d30.md` | 144151ae 09-19 | 未判 | — | — |
| 23 | `23-system-window-input-tools.md` | 144151ae 09-19 | 未判 | — | — |
| 24 | `24-doc-search-tools.md` | 144151ae 09-19 | 未判 | — | — |
| 25 | `25-s3-acceptance.md` | 144151ae 09-19 | 未判 | — | — |
| 26 | `26-tts-output.md` | 6e2a68fa 09-20 | 未判 | — | — |
| 27 | `27-punctuation.md` | 144151ae 09-19 | 未判 | — | — |
| 28 | `28-session-scope-warm.md` | 130c0943 09-19 | 未判 | — | — |
| 29 | `29-memory-l1l2.md` | 144151ae 09-19 | 未判 | — | — |
| 30 | `30-result-routing-d10.md` | 144151ae 09-19 | 未判 | — | — |
| 31 | `31-reminders.md` | 144151ae 09-19 | 未判 | — | — |
| 32 | `32-s4-acceptance.md` | 130c0943 09-19 | 未判 | — | — |
| 34 | `34-frontend-scaffold.md` | 144151ae 09-19 | 未判 | — | — |
| 35 | `35-panel-bridge-c17.md` | 144151ae 09-19 | 未判 | — | — |
| 36 | `36-result-history-panel.md` | 15ff2bee 09-26 | 未判 | — | — |
| 37 | `37-approval-ui-l2.md` | 15ff2bee 09-26 | 未判 | — | — |
| 38 | `38-command-palette-tasks.md` | 15ff2bee 09-26 | 未判 | — | — |
| 39 | `39-config-editor-gui.md` | 15ff2bee 09-26 | 未判 | — | — |
| 40 | `40-security-privacy-cost-pages.md` | 2df182e1 09-28 | 未判 | — | — |
| 41 | `41-kws-wake-word.md` | 130c0943 09-19 | 未判 | — | — |
| 42 | `42-watchdog.md` | 144151ae 09-19 | 未判 | — | — |
| 43 | `43-power-lifecycle-events.md` | 144151ae 09-19 | 未判 | — | — |
| 44 | `44-costmeter-c23.md` | 5cba2d92 09-19 | 未判 | — | — |
| 45 | `45-diagnostics-guards.md` | 144151ae 09-19 | 未判 | — | — |
| 46 | `46-s6-acceptance.md` | 144151ae 09-19 | 未判 | — | — |
| 47 | `47-task-scheduler-pathlock.md` | 144151ae 09-19 | 未判 | — | — |
| 48 | `48-approval-queue-full.md` | 144151ae 09-19 | 未判 | — | — |
| 49 | `49-session-grants-d45-2.md` | 144151ae 09-19 | 未判 | — | — |

## §2 本批计数

仍成立 **未判** ／ 已失效 **未判** ／ 量不到 **未判** ／ 合计 30
