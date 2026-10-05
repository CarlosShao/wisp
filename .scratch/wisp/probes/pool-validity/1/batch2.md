# pool-validity-1 / 批二 —— 零勾开放票有效性普查（号段 50–99）

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

## §1 名册（16 枚）

| 票号 | 票文件 | 上次动过 | 判定 | 复法（跑的那把尺） | 读数 |
|---|---|---|---|---|---|
| 50 | `50-tier1-manifest-plugins.md` | 144151ae 09-19 | 未判 | — | — |
| 51 | `51-tier2-goja.md` | 6160bc8f 09-19 | 未判 | — | — |
| 52 | `52-d46-command-plugins.md` | 144151ae 09-19 | 未判 | — | — |
| 53 | `53-larkcli-plugin-e2e.md` | 144151ae 09-19 | 未判 | — | — |
| 54 | `54-s7-acceptance.md` | 144151ae 09-19 | 未判 | — | — |
| 55 | `55-s8-macos-port.md` | e208efb1 09-21 | 未判 | — | — |
| 56 | `56-s8-signing-distribution-naming.md` | 144151ae 09-19 | 未判 | — | — |
| 57 | `57-s8-plugin-sdk-registry.md` | 144151ae 09-19 | 未判 | — | — |
| 58 | `58-s8-i18n-english.md` | 144151ae 09-19 | 未判 | — | — |
| 59 | `59-p15-aec-spike.md` | 130c0943 09-19 | 未判 | — | — |
| 60 | `60-c32-realtime-engine.md` | 5cba2d92 09-19 | 未判 | — | — |
| 61 | `61-cloud-voice-providers-c9.md` | 5cba2d92 09-19 | 未判 | — | — |
| 65 | `65-ball-glass-quality-rework.md` | ed1e7c62 09-20 | 未判 | — | — |
| 86 | `86-resolvepercallbudget-wallclock-fragility.md` | 3d716a4b 09-21 | 未判 | — | — |
| 91 | `91-os-isolation-spike-restricted-token-vs-appcontainer.md` | 6114e3df 09-21 | 未判 | — | — |
| 98 | `98-cmd-wisp-tests-never-run-on-this-host.md` | 92818af9 09-28 | 未判 | — | — |

## §2 本批计数

仍成立 **未判** ／ 已失效 **未判** ／ 量不到 **未判** ／ 合计 16
