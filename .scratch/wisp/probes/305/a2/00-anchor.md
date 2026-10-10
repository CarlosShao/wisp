# 305-a2 · 00 起手闸门重锚（只读普查腿）

腿：`305-a2`（只读普查，零产码）。仓库：`D:/work/workspace/projects plans/Wisp`，分支 `dev`。
取数时刻：`2026-10-11 07:50:27 +0800`（`logs/00-date.txt`，rc=0）

## 四件闸门（逐件附尺与 rc）

| 尺 | 命令 | rc | 结果 |
|---|---|---|---|
| 1 | `date "+%Y-%m-%d %H:%M:%S %z"` | 0 | `2026-10-11 07:50:27 +0800` |
| 2 | `git log -1 --format='%H%n%ad%n%cd%n%s' --date=iso-strict` | 0 | `e4740e35ff62d9dc8a9732a20cfedbbf4d501a17` / `2026-10-11T07:47:35+08:00` / `306-r1 交件名册逐字（终局）` |
| 3 | `git status --porcelain -- .scratch/wisp/probes/305`（scoped，只数自家路径） | 0 | 1 行：`?? .scratch/wisp/probes/305/a2/`（＝本腿自己刚建的目录，见下） |
| 4 | `git rev-parse HEAD` | 0 | `e4740e35ff62d9dc8a9732a20cfedbbf4d501a17` |
| 附 | `git rev-parse --abbrev-ref HEAD` | 0 | `dev` |

raw 件：`logs/00-date.txt` · `logs/00-git-log-1.txt` · `logs/00-scoped-porcelain.txt` · `logs/00-head.txt` · `logs/00-branch.txt`

## 与我这一单（任务书）不符处 —— 具名报回

- **HEAD 号腐烂**：任务书写着「`git rev-parse HEAD` 重锚，⛔ 用我写的 `6284a489`」。盘上现量 HEAD＝
  **`e4740e35ff62d9dc8a9732a20cfedbbf4d501a17`**（subject `306-r1 交件名册逐字（终局）`，提交时刻 `2026-10-11T07:47:35+08:00`，
  即本腿起手前约 3 分钟）。⇒ 本件一律以 `e4740e35` 为锚，`6284a489` 不作为任何锚点引用。
- scoped porcelain 那 1 行 `?? .scratch/wisp/probes/305/a2/` 是本腿 `mkdir` 的产物，非别家脏面。
  任务书点名的别家脏面（`.gitignore` 被改、`design/**` 大量 `D`、别家 `probes/**`、票 303 文件）**不在**我这把 scoped 尺的射程内，
  本腿对它们**零接触**（既不碰也不抱怨）。

## 本腿车道（自我声明，供裁决者核）

- ⛔ 任何 Go 编译面：`go build`／`go vet`／`go test`（含 `cmd/wisp`）／`scripts/d22scan.sh`／`gofumpt` —— 全程未跑、不会跑。
  `go env`／`go list` 若跑必逐条自报命令与 rc。**本次交件截至此节：`go env`／`go list` 亦未跑过（0 条）。**
- ⛔ 真窗／真机；⛔ `git checkout`／`stash`／`reset`／`rebase`／`clean`／`--amend`；读旧笔一律 `git show <ref>:<path>` 对象层。
- ⛔ 读 `%APPDATA%\wisp-dev\secrets\`；⛔ 读机主 config 的值；⛔ 写 `frontend/src/**`／`design/**`。
- 临时件只建不删，全部落本目录 `logs/`；证据件不叫 `.out`。
- 每把尺自落一行 `rc=N`，取退码一律 `> 件 2>&1; rc=$?`。
