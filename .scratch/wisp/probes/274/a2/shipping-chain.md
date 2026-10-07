# 274-a2 出货通路普查 — 页面字节的落地机在哪台机器上

> 腿：`274-a2`（只读普查，Wisp 仓 `D:/work/workspace/projects plans/Wisp`，分支 `dev`）。
> 交付件本体＝本文件；每条读数指名它出自 `logs/` 哪枚文件。
> 所有读数均为本腿**现跑**；票面/台账旧数只作为「待复量的对象」出现，不当凭据。
> 唯一写面＝`.scratch/wisp/probes/274/a2/**` 与本腿的 commit。

## §0 起手锚

原始输出：`logs/anchor.txt`。

| 尺 | 读数 |
|---|---|
| `date` | `Wed Oct  7 11:02:00 CST 2026` |
| `git rev-parse --short HEAD` | `933a8341` |
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git status --porcelain -- cmd internal scripts tools .github docs frontend` | **0 行（干净）**，退出码 0 |
| `git ls-files frontend/src \| wc -l` | **64** |
| `git ls-files frontend \| wc -l` | **85** |

工作树现状说明（本腿不处理）：`design/**`、`.gitignore`、`.scratch/**` 里有别的腿/机主的**未提交改动**（含 `design/**` 的一批 `D` 删除记录与若干 `??` 未跟踪件），**不是本腿的**；本腿既不还原也不提交它们。终态自查一律只按上面那六族路径（`cmd internal scripts tools .github docs frontend`）判定。

复量对照（待复量对象，非凭据）：编排者给的背景与工单 `274` 里「dev 的 64 枚页面源码」这一枚，本腿现跑 `git ls-files frontend/src` **＝ 64，对上**。
