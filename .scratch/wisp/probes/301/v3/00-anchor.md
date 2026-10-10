# 301-v3 — 起手锚（非实现者验收腿，票 301 `AC#5` + `AC#4` 残留必答）

尺与读数＝逐字来自 `logs/anchor-start.txt`（先落件再量，双锚先例＝本票 `301-a1`/`301-a2`、`301-v1` §0）。

| 格 | 尺 | 读数（2026-10-10 16:07:23 +0800） |
|---|---|---|
| 起手时刻 | `date '+%F %T %z'` | `2026-10-10 16:07:23 +0800` |
| 起手 HEAD | `git log -1 --format='%h %ad %s' --date=iso-strict` | `a0e55a51` ＋ `2026-10-10T15:51:05+08:00`（subject 截 80 字节存件） |
| 起手 porcelain 本票射程 | `git status --porcelain -- scripts internal .github docs frontend` | **0 行**（`porcelain_lines=0`） |
| 起手 tasklist `wisp.exe` | `tasklist //FI "IMAGENAME eq wisp.exe"` | **0 枚**（stdout 逐字＝`INFO: No tasks are running which match the specified criteria.`） |
| 起手 tasklist `balldebug.exe` | `tasklist //FI "IMAGENAME eq balldebug.exe"` | **0 枚**（同形逐字） |
| 件退码 | `echo "anchor_rc=$?"` | `anchor_rc=0` |

射程声明＝本腿写面只有 `.scratch/wisp/probes/301/v3/**`，且全部为新建 `.md`／`.txt`；
⛔ 新建或修改 `.sh`／`.ps1`／`.go`；⛔ 碰票面、台账、`HANDOVER`、`scripts/portable-tests.sh`、
`.github/workflows/ci.yml`、`internal/**`、`frontend/**`、`design/**`；⛔ push／`--amend`／`reset`／
`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`／`git add -A`／`git add .`；⛔ 改 git 配置。

环境＝本机 `GOOS=windows / GOARCH=amd64`（与 `301-v1` §0 同机同模块；本腿⛔ 跑任何计时敏感档，
`--scope=windows`／`--scope=core` 两档的实跑颜色一律**采 CI 字节**或引既有件并标〔⛔ 我复跑〕）。
