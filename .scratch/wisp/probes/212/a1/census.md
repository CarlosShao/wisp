# 212-a1 — 注释引用仓内路径的全仓普查（票 212 AC#1 的分母/分子名册）

> 本件是**只读普查腿 `212-a1`** 的交件，对应票 `.scratch/wisp/issues/212-comments-cite-evidence-files-that-do-not-exist.md` 的 **AC#1**。
> 本腿**零产码、零 Go 命令**（闸门见下 §0.3），所有结论只来自「读文件 + grep + `test -f`」这一类盘上尺。
> 每一格都带**尺的原文命令 + 真实读数**；只给结论不给尺的一律不算答案。

---

## §0 锚点、脏名册、本腿闸门

### 0.1 锚点（本腿自取）

| 项 | 值 | 尺 |
|---|---|---|
| 起手工号时间 | `2026-10-03T10:51:40+08:00` | `date -Iseconds`〔量 S0-1〕 |
| HEAD | `1a14f203b8ae42cb61178feae9969b48c6677649` | `git rev-parse HEAD`〔量 S0-2〕 |
| 分支 | `dev` | `git rev-parse --abbrev-ref HEAD`〔量 S0-3〕 |
| 工作树脏条目总数（起手） | **482** | `git status --porcelain \| grep -c "^"`〔量 S0-4〕 |

> ⚠ 本件里**每一枚基线数都必须现跑并带 HEAD**。上面这枚 HEAD 是本腿全部读数的唯一锚点；
> 后续任何一格若换了锚点，会在该格具名标注。

### 0.2 起手脏名册（凡落在这批文件上的读数一律标「未提交中间态」，⛔ 不当终态）

尺：`git status --porcelain`〔量 S0-4，读数 482 条〕

**2.1 已跟踪且被改动的产码（M/D）——本腿关心的子集**

| 文件 | 状态 | 对本腿的意义 |
|---|---|---|
| `cmd/wisp/panel_host_windows.go` | ` M` | **写腿 `255-r3` 在飞**；此件上的注释引用＝中间态 |
| `cmd/wisp/panel_resident_windows.go` | ` M` | 同上 |
| `.scratch/wisp/probes/152/my152.py` | ` M` | 探针件，非产码 |
| `.scratch/wisp/probes/154/gate-clauses.sh` | ` M` | 门禁件（票 212 禁区点名不许改，本腿只读） |
| `.scratch/wisp/probes/161/r6/logs/flip-*.txt`（8 枚） | ` M` | 日志 |
| `.gitignore` | ` M` | 影响 d22scan 的忽略过滤口径（§D 相关） |
| `design/**`（22 枚 ` D`/` M`） | — | **本腿不读不引**（两层禁令） |

**2.2 未跟踪的产码形状（`??`）——本腿关心的子集**

| 文件 | 意义 |
|---|---|
| `internal/tools/paths_twocontainments_252_r2_test.go` | **写腿 `252-r2` 在飞**，测试件未落树；此件里的注释引用不在 `git ls-files` 名册里 |
| `internal/tools/paths_twocontainments_252_r2_windows_test.go` | 同上 |
| `internal/tools/`（` M` 无、但整目录可能有别的未跟踪件） | 见 §7 量不到的格子 |
| `.scratch/wisp/probes/171/r3/**` | **写腿 `171-r3` 在飞** |

**2.3 仓根垃圾件（与本票同族、但不是本票射程）**

起手 `ls -la` 复读：仓根存在 **30+ 枚零字节、文件名是被误执行的中文/emoji 片段** 的产物
（`****（）。`、`⚠`、`「干净`、`📐`、`§15`、`1`、`2026-09-24`、`2379`、`-`、`落位：PersonalAudioAssistant\`` 等），
时间戳集中在 `Oct 3 10:20–10:21`。
尺：`ls -la`〔量 S0-5〕

⇒ 判读材料：这正是本票那一族病的**孪生形状**——**该加引号的 heredoc 没加引号，反引号里的路径被 shell 真的执行了**，
于是命令替换的输出被当成文件名落成了零字节文件。本腿**不处置、不删除**（禁区：临时件只建不删），只登记。

### 0.3 本腿闸门（自陈遵守情况）

| 闸门 | 遵守 |
|---|---|
| ⛔ 禁跑任何 Go 命令（含 `tools/d22scan`） | ✅ 全程 `grep`/`find`/`sed`/`test -f`/`wc -l`/`git`，**未执行一次 `go`** |
| ⛔ 产码一字不动，写面只有本件 | ✅ 单文件 `.scratch/wisp/probes/212/a1/census.md` |
| ⛔ `frontend/**`、`design/**` 不读不引 | ✅ 全部尺的搜索根显式写作 `cmd internal tools docs scripts .scratch`；`design/**` 仅在 §0.2 以**文件名**出现在脏名册（来自 `git status`，非读取内容） |
| ⛔ 不碰 `^- [ ]` 勾选框、不改 `docs/**` | ✅ |
| ⛔ 不 push / 不 `add -A` / 不 `--amend` / 不 `reset` / 不 `stash` / 不 `checkout .` / 不 `clean` | ✅ 只 commit，显式 pathspec |

---

## §1 A — 分母：注释里引用仓内路径的**形状穷举**

（口径与逐形尺读数以 §1.2 表为准；本节先立骨架）

## §2 B — 分子：引用对象盘上不存在的

## §3 C — 那 1 枚「今天的尺读不到」的形状

## §4 D — 装牙的代价表（给编排者裁，本腿不裁）

## §5 E — 不加 ban 之外的路

## §6 F — 推翻清单（票面与本腿派单每一句都是待验断言）

## §7 量不到的格子 + 复量法
