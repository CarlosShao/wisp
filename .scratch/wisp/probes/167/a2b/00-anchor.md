# 167-a2b 起手锚＋射程（只读普查腿：只量「草稿」与「崩溃自救」两格）

- 腿名 `167-a2b`；工作语言中文；本程**零 Go 命令**（另有一条腿 `253-v1` 在 `cmd/wisp` 种突变跑测试，任何并发 Go 面都会互洗读数）；写面＝只新建 `.scratch/wisp/probes/167/a2b/*.md`；零翻框／零改既有文件／不 `-done`／不 push。
- 承接：前身 `167-a2` 被每日额度掐死在起手锚之后（只落 `.scratch/wisp/probes/167/a2/00-anchor.md`，已由编排者代提入库 `d33d492d`）；本程接它没干的活，**只量两格**：格一＝草稿（用户打了字还没发出去那段文本）；格二＝崩溃自救（程序崩过一次之后界面怎么知道自己该显示什么）。
- 不重做 `167-c2` 已量三格（占用条／任务序号／停止按钮），件在 `.scratch/wisp/probes/167/c2/census.md`，commit `9a2a3e00`。
- 前身 `a2/00-anchor.md` §2 曾记一枚"这两格不是零现量、曾由 167-a2 两轮量清（census.md 194 行／42,736 字节）"的前提冲突提醒：本程对**当前 HEAD 全部现量重跑**，不抄旧行号（`cmd/wisp` 行号会漂）。

## 1. 起手锚（全部现量，2026-10-08 18:2x +0800）

| 尺 | 读数 | 命令形状 |
|---|---|---|
| 当前时间 | `Thu Oct 8 18:20:06 CST 2026` | `date` |
| HEAD | `8dd239f14d8338752882b16a1fdfa94719ad0454`（`Thu Oct 8 18:01:41 2026 +0800`） | `git log -1 --format='%H %ad %s'` |
| 在飞总数 | **759** 条 | `git status --porcelain \| wc -l` |
| 在飞首屏（只登记不动） | `.gitignore`；`probes/152`；`probes/161/r6/logs/flip-*` 9 枚；`probes/242`；`probes/268`；`design/**`（多枚 M＋`design/assets/*` 4 枚 D） | `git status --porcelain \| head -20` |

> 写面归属（写面可达性三答的登记面）：`git status --porcelain -- cmd/wisp internal/panel frontend` 输出为空、rc=0 ⇒ 本程起手时这三条路径**工作树零在飞改动**。

## 2. 票 167 四把尺（现量，rc 全 0）

| 尺 | 读数 | 命令 |
|---|---|---|
| `wc -l` | **92** | `wc -l .scratch/wisp/issues/167-four-small-outputs-the-panel-needs-occupancy-queue-stop-draft-plus-crash-self-rescue.md` |
| 未勾 `- [ ]` | **6**（AC#2–AC#7） | `grep -cE '^[[:space:]]*- \[ \]' <票>` |
| 已勾 `- [x]` | **1**（AC#1） | `grep -cE '^[[:space:]]*- \[x\]' <票>` |
| 票面 porcelain | 空（票面一字未动） | `git status --porcelain -- <票>`（rc=0，无输出） |

## 3. 本程主尺（派单指定，两格各一把）

- 格一（草稿）：`git show HEAD:internal/panel/composer.go \| sed -n '180,260p'`（`ComposerState` 那族字段名册，实际文件 346 行、须再续读 260–346）＋ `git grep -niE 'draft' HEAD -- 'internal/panel/*.go' 'cmd/wisp/*.go' ':(exclude)*_test.go'`。须分开答「有没有读者」「有没有写者（入向有没有一条路能把页面那段未发的文本递给 Go）」；档位三选一：〔有字段有人读〕／〔有字段没人读〕／〔没字段〕。
- 格二（崩溃自救）：`git grep -niE 'recover|resume|restart|crash' HEAD -- 'internal/panel/*.go' 'cmd/wisp/*.go' ':(exclude)*_test.go'`，逐枚剥掉日志/错误处理里的同名用法（**按类型报枚数，别按词面**）；判据＝面板出向读面（`Snapshot`／`ApprovalCardView`／视图那族结构体）里有没有一枚字段承载"上次崩溃/上次异常退出"。不许把 `internal/observe` 的 `defaultPanicSink` 当"面板知道自己崩过"。
- 判"零调用者"锚调用形状（`x.Method(`）；函数值/字段赋值那族不带括号，另跑赋值形状尺。引"裸文件名:行号"先定包再量行数（本仓 `bridge.go` 有三枚：`internal/panel`／`internal/tools`／`internal/models`）。范围尺不接 `2>/dev/null`、跑完必检 rc；枚数尺写明射程文件。
