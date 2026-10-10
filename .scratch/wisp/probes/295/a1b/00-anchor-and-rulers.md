# 票 295 · 补位只读普查腿 `295-a1b` · 第 1 笔：射程声明 + 锚 + 尺

**本件性质**：对死腿 `295-a1` 已交件 `.scratch/wisp/probes/295/a1/10-ac0-field-roster.md`
（提交号 `3f9bedf7fa309da74a2c0320c32937b423c13f6b`，`git log -1 --format=%H -- <该路径>` 现跑所得，非转述）
里 ⓐⓑ 两张表的**每一个锚点**，对着 **HEAD 对象层**重跑一遍，并具名报出对不上的地方。
⛔ 本腿不是落地腿、⛔ 不是验收腿、⛔ 零源码改动、⛔ 零翻框。

## 三枚锚（全部现跑，⛔ 不引编排者转述）

| 锚 | 值（现跑） | 尺（逐字命令） |
|---|---|---|
| 本腿 HEAD | `29081a1367fa4024716378f1b18486fc97b6edbf`（短号 `29081a13`，branch `dev`） | `git log -1 --oneline` / `git rev-parse HEAD` / `git rev-parse --abbrev-ref HEAD` |
| 死腿交件的提交 | `3f9bedf7fa309da74a2c0320c32937b423c13f6b`，件正文 `20365` 字节 | `git log -1 --format=%H -- .scratch/wisp/probes/295/a1/10-ac0-field-roster.md` ／ `git show <该号>:…10-ac0-field-roster.md \| wc -c` |
| 死腿自称锚点 | `448f5a57`（其件第 4 行逐字：`**锚点** \`448f5a57\`，branch \`dev\``） | 本腿 Read 该件所得 |

⚠ **锚差即"树会呼吸"的量化面**：死腿量的是 `448f5a57`，本腿复跑的是 `29081a13` ⇒ 两锚之间 `cmd/wisp/**` 若被写腿动过，
行号漂移是**预期**而不是异常；本件的任务就是把漂移逐枚钉出来，并把每个锚改写成**内容锚**（逐字抄那一行原文）。

## 射程声明（三件事逐枚写清：尺／目录／含不含 `_test.go` 与注释行）

- **锚一律取 HEAD 对象层**：`git show HEAD:<path>`、`git grep -n '<内容片段>' HEAD -- <路径>`。⛔ 不拿工作树当锚点来源。
  凡本腿也看了工作树（为对照漂移），同一把尺**并排跑两棵**，两把读数都具名写出。
- 射程目录：`cmd/wisp`（生产面＋测试面同目录）；对照面在 `internal/ball`、`internal/audio` 只做线程判定用。
- 尺命令逐字：
  1. `git grep -n -o "\.started\b" HEAD -- cmd/wisp` — **含** `_test.go`、**含**注释行（`-o` 只打印命中片段，不过滤注释/测试）。
  2. `git grep -n "atomic.Bool" HEAD -- internal cmd` — **含** `_test.go`、**含**注释行。
  3. `git grep -n -E "voiceEnabled|mutedAtBoot" HEAD -- cmd/wisp` — **含** `_test.go`、**含**注释行。
  4. 定义窗整窗读：`git show HEAD:cmd/wisp/resident_audio_windows.go | cat -n`（ball／resident_windows 同形）。
- **go 命令**：⛔ 零编译面（`go build`/`vet`/`test`/`run`/`gofumpt`/`d22scan` 一枚不跑）。`go env`/`go list` 若跑过在本件末节自报，没跑就写零枚。
- ⚠ 本腿⛔ 不起任何进程（编排者正在同机跑出厂构建链）；⛔ 不动 `build/**`、`frontend/dist/**`。
- ⛔ 不改死腿那份件；本腿结论**另写新件**，凡更正一律"死腿原话 ↔ HEAD 现量 ↔ 差在哪"三列式。
