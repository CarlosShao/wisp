# 35-a6 — 起手锚（只读普查腿，先于一切长跑命令）

腿代号 `35-a6`。任务＝票 35 面板桥 **6 枚未勾 AC 框**（`:42 :44 :45 :47 :49 :51`）的
**可开工性普查**：逐格答四问（有没有被测物／最少落点／要不要新增 C17 面／有没有仪器能钉）。
⛔ 本腿不实现、不改产码、不改测试、不翻任何一枚 AC 框、不 push、不开任何真窗。

## 1. 现量 HEAD（⛔ 不采用派单里转述的号）

- `git log -1 --format=%H` = `7a367a08292ef94908d86656608974e81c42f0d4`
- 分支 = `dev`（`git rev-parse --abbrev-ref HEAD`）
- 起手时刻：现量于本件写入前一把命令。共享工作树随时前进，下面所有引用都以**本锚时刻**的 HEAD 为准；
  跨 ref 的行号一律先 `git show <ref>:<file> | wc -l` 验存在性。

## 2. 现量脏面

- `git status --porcelain -- cmd internal docs .github scripts` = **0 行**（非空输出见下面 §5 的自证：本腿未跑任何写操作之前该尺就为 0）
- 含义：本腿的 go 读数窗口约束（`33-v4` 独占编译面/卫生仪器）不适用在「工作树脏」这一维，
  但仍适用「⛔ 零 `go build`/`go vet`/`go test`」——本腿全程零 `go` 命令，除 §4 具名的只读查询。

## 3. 禁区 / 写面对拉用的 blob 号（现量 `git rev-parse HEAD:<path>`）

| 用途 | 路径 | blob |
|---|---|---|
| 禁区（C17 名册＋入向三守卫，⛔ 动名册＝改契约＝人工批准） | `internal/panel/bridge.go` | `bebe8e702a85c640641551e93dd09936530950f4` |
| 禁区（传输钩子＋门＋守卫文案，本腿零写） | `cmd/wisp/panel_host_windows.go` | `1f9060dfff33cffab1317e5f653e1a95f9678e97` |
| 票面（⛔ AC 框一字不动，只在文末追加一节） | `.scratch/wisp/issues/35-panel-bridge-c17.md` | `e8644ab7f484966ed4db69708001a7d2d044c3cb` |
| 仪器#1（`:42`/`:44` 最近的尺） | `cmd/wisp/panel_transport_35r1_test.go` | `20de69803f818afb24c0f29c4eb8e124243f6494` |
| 仪器#2（行为尺，12 枚名册载体） | `cmd/wisp/panel_transport_35r2_test.go` | `a466b4ed49420561d65d22ff7e30d4dc0f6ae571` |
| 仪器#3（入向三支守卫红） | `cmd/wisp/panel_inbound_guards_35r3_test.go` | `a917e0dc1e2ef0801c189da66a6b726de3c45301` |
| 仪器#4（真窗，〔仅本机可量〕，⛔ 本腿不开窗） | `cmd/wisp/panel_transport_live_35v2_windows_test.go` | `e778be886e5590965ba5a9d5f01e2f09c3e3e786` |

⚠ 与票面历史读数的差异具名：票 `:202`/`:211` 记的 `panel_host_windows.go` blob 在 `35-v3` 时是
另一枚、`:205`/`:209` 段又记 `26b5de83`（`35-r4`/`35-r5` 时期）。本锚的 `1f9060df` ＝ **此刻 HEAD 的真身**，
后面所有行号引用按它量，不许沿用票面旧号。

## 4. 本腿的 go 命令口径（自报）

- 起手锚之前：**零条** `go` 命令。
- 全程只允许只读查询；实际跑了哪几条 `go env` / `go list` 会在终局具名登记（若一条没跑，也照实写「零条」）。
- ⛔ 未跑 `go build` / `go vet` / `go test` / `go run`（`33-v4` 独占编译面与 d22scan/gofmt/crlf 读数）。

## 5. 落点约定

- 写面只有 `.scratch/wisp/probes/35/a6/**`（全 `.md`）＋票 35 文末追加一节。
- ⛔ 不新建 `.sh`/`.ps1`/`.txt`/`.out`（`.scratch/**` 里的脚本与产物会挪 gofumpt/d22scan 分母名册，票 276 那枚病）。
- 需要较长 shell 的工作件放仓外 `D:/tmp/wisp35a6/`；⛔ 不在仓库目录内建 worktree／checkout。
- 每枚证据件自带 `rc=N` 行；0 字节件＝那一格没交。
- 追加票面后自数 `grep -cE '^[[:space:]]*- \[ \]'` 前后同数并写进件里。

## 6. 框数尺（追加任何内容之前）

- 尺＝`grep -nE '^[[:space:]]*- \[ \]' .scratch/wisp/issues/35-*.md`
- 未勾＝**6** 枚，行号＝`:42 :44 :45 :47 :49 :51`（与派单转述**同色**，本腿复跑确认，未顶回）
- 已勾＝**4** 枚（`:52` Transport agreement／`:63` inbound judge／`:75` behavioural yard／`:77` AC#8）
- 票面总行数＝**377**
