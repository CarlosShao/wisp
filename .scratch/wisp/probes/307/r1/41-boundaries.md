# 307-r1 · 41-boundaries —— 逐笔越界＋收口确认发

尺与全量输出＝`logs/commit-bounds.txt`（`date`／`git rev-parse HEAD`／卫生计数／逐笔名册全在件里，⛔ 区间尺）。

## 1. 逐笔越界（`git show --name-only --format=` 一笔一把尺）

| 笔 | 谁的 | 越出 `internal/risk/**` ＋ `.scratch/wisp/probes/307/r1/**` 的行数 |
|---|---|---|
| `c6691756` | 307-r1（起手锚） | **0** |
| `249753a4` | **编排者**（收 307-a1：票面追加节＋A833＋HANDOVER） | 7 ⇒ **别家的笔，⛔ 记在本腿账上** |
| `858ee7f4` | **编排者**（欠读数①②交件，`probes/307/orch/**`） | 2 ⇒ 同上 |
| `51a0a0f2` | 307-r1（修复＋AC#1 凭据） | **0** |
| `183032c1` | 307-r1（AC#2–AC#5 凭据） | **0** |
| 收口这一笔（本件＋`logs/*-final.txt`＋`commit-bounds.txt`） | 307-r1 | **本件⛔ 能自量**（尺跑在内容落定之前）⇒ 名册逐字＝该笔的 `git show --name-only --format=`，**留给编排者／验收腿复跑**；本腿自陈写面只有 `internal/risk/**` ＋ `.scratch/wisp/probes/307/r1/**`，`git status --porcelain -- internal cmd`＝0 行即其旁证 |

⇒ **本腿三笔（＋本笔）越界＝0 格**；其余两笔的"越界"恰好证明派单那句"⛔ 区间尺"——
`git rev-list cb7bb3c..HEAD` 一把尺会把编排者的两笔数进来（本腿第一发真踩过，具名见 `40-ac4-gates.md` §5）。

其它卫生尺（逐字在 `logs/commit-bounds.txt`）：

- `git status --porcelain -- internal cmd` = **0 行**；`-- .scratch/wisp/issues docs` = **0 行**（⇒ 本腿⛔ 动过票面与台账）。
- 三把进程尺（`wisp.exe`／`balldebug.exe`／`mockllm.exe`）收尾复跑 = **0／0／0**。
- `find .scratch/wisp/probes/307/r1 -type f -size 0` = **0 枚**（带 `-type f`）。
- 本腿各笔里 `.out` 件名 = **0**（`git show --name-only` 各笔合并后 `grep -c '\.out$'` = 0）。
- 新 `.go` 落在 `probes/**` = **0 枚**（本腿唯一的新 `.go` 是 `internal/risk/syncdirs_alias_307_windows_test.go`）。
- 仓外导出树 `C:\Users\swq\tmp\wisp-307-r1-tree-main\`（变异面，⛔ 是 git 仓、⛔ 在仓库目录内）与 pristine 副本
  `C:\Users\swq\tmp\wisp-307-r1-pristine-syncdirs.go` ⇒ 临时件**只建不删**；导出树里的 `zz307probe_test.go` 只活在那棵树里。
- 别名台面目录 `C:\Users\swq\tmp\riskattrib1\RunnerTempAreaForAttribution1` **⛔ 删、⛔ 重建**（只读复用）。

## 2. 收口确认发（最终 HEAD 上的第三发，⛔ 新码⛔ 新配置）

尺＝`logs/run-roster.ps1 -Tag final`（与改前／改后同名册同参数：整包 `go test ./internal/risk/ -count=1 -v`，两台面各两发）：

| 台面 | 发 | rc | RUN | PASS | FAIL | SKIP |
|---|---|---|---|---|---|---|
| default | full / full2 | 0 / 0 | 200 / 200 | 199 / 199 | **0 / 0** | 1 / 1 |
| aliased | full / full2 | 1 / 1 | 200 / 200 | 192 / 192 | **7 / 7**（＝甲族 7 枚，⛔ 本票射程） | 1 / 1 |
| default | `-run TestSync` | 0 | 38 | 37 | 0 | 1 |
| aliased | `-run TestSync` | 0 | 38 | 37 | **0** | 1 |

那 5 枚指名在 aliased 台面这一发逐字（`logs/full-aliased-final.txt`）：
`--- PASS: TestSyncFixtureFallbackAndMatch (0.01s)` / `TestSyncFallbackNotDisarmableByWeakRoot (0.02s)` /
`TestSyncUnverifiedRootKeepsFallback (0.01s)` / `TestSyncSuspectFallbackIsComponentBounded (0.01s)` / `TestSyncSuspectFallbackWhenUndetectable (0.01s)`。
⇒ 与 §1–§3 各件里"改后"那一发**逐字同集合同计数**（三发同台面交集⛔ 空）。

## 3. ★共享工作树的脏面（本腿⛔ 制造、⛔ 提交，具名为编排者的越界审计存档）

起手锚只量了 `git status --porcelain -- internal cmd`（＝0 行，本腿车道），收口时本腿另用**全工作树**尺看了一遍，
盘上现量＝工作树里躺着**别家的**未提交改动，本腿⛔ 写过一枚、⛔ 提交过一枚：

- ` M .gitignore`（新增 `.worktrees/` 一段）、` M go.mod`（`go-webview2` 从 `// indirect` 升成直接依赖那一族改写）、
  ` D design/**`（十余枚删除）、` M .scratch/wisp/probes/152|161|242|268/**`、以及根目录一堆 `??` 未跟踪件
  （`.scratch/commit-msg-*.txt`、`.scratch/ci-logs/**`、`-` 之类）。
- 尺＝`for c in c6691756 51a0a0f2 183032c1 fb6a171e; do git show --name-only --format= $c; done | grep -E '^(go\.mod|go\.sum|\.gitignore)$' | wc -l` ⇒ **0**
  ⇒ 本腿四笔**没把任何一枚别家脏件带进 commit**（每笔都带显式 pathspec，且 pathspec 写在 `$( … )` 之外）。
- ⚠ 一枚⛔ 得⛔ 解释的向：本腿跑过 `GOFLAGS= go build ./...` 与 `go vet ./internal/risk/`，
  Go ≥1.16 默认 `-mod=readonly`（改写 go.mod 会**报错**⛔ 静默改），且那枚 `go.mod` diff 的形状＝`go mod tidy` 的产物⛔ 是 build 的产物
  ⇒ 判读＝**先于本腿存在**；本腿交的是"⛔ 在我这手里变过"这一枚**推理**，⛔ 是"实测过它在我进场时已脏"（起手⛔ 量全工作树，那一格是本腿的尺缺陷，具名认）。

## 4. 本票欠着的读数（具名，⛔ 掩）（＝`msg-04` 里我写的那节"§3"，本节落定后renumber 成 §4，具名报回）

1. `AC#2(d)` 那两枚⛔ 响的指名用例（`TestSyncFallbackNotDisarmableByWeakRoot`、`TestSyncUnverifiedRootKeepsFallback`）＝枚级无负向钉，本票⛔ 动其期望侧。
2. `alias307Pair` 的 `t.Fatalf`（无别名即 Fatal）分支在本机从未被执行（需关卷的 8.3 生成＝改机器配置，⛔ 授权）。
3. 托管 CI runner 上的两台面颜色⛔ 本腿取（⛔ push）。
4. `:431` 把 `normPath` 比较形当 `Root.Path` 报出——本票按登记**原地未动**，
   但换形后这一枚报出的拼法在 `AC#2(d)` 的红句里现量过一次：
   `Path:c:\users\swq\tmp\riskat~1\runner~1\testsy~1\001\thirty~1`（小写＋短形链），
   与 `blacklist.go:128-132` 那句 ticket 75 的注释相抵 ⇒ **仍归编排者那本欠账**（追加节 §3 登记①）。
5. 枚级耗时代价⛔ 可测（`-v` 钟的 0.01 s 粒度吃掉了 `formsOf` 那一次祖先走步），只交到"包级 ±0.06 s 内"这一层（`10-ac1-colours.md` §4）。
