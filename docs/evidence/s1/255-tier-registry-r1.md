# 255-r1 证据件 —— AC#2-ⓑ 按键/按段档位登记表 ＋ AC#3 段外哑键的一把会响的尺

**代号**：`255-r1` · **性质**：落地写腿 · **写面**：`internal/config/**` ＋ 本件
**票面**：`.scratch/wisp/issues/255-config-says-these-sections-took-effect-immediately-panel-for-a-width-nobody-reads-while-the-panel-host-never-receives-config-and-tier-has-zero-readers.md`
**本腿只做两格**：AC#2-ⓑ（编排者裁定 `A527`）＋ AC#3。⛔ 不翻任何 AC 勾选框（勾归非实现者裁）。

---

## §0 起手锚（逐字读数）

| 尺 | 逐字读数 |
|---|---|
| `date '+%Y-%m-%d %H:%M:%S %z'`（进场第一发） | `2026-10-02 16:33:18 +0800` |
| `git rev-parse --short HEAD` | `3df8b82b` |
| `git branch --show-current` | `dev` |
| `git status --porcelain \| wc -l`（全仓，含他腿在飞的脏件） | `366` |
| `git status --porcelain internal/config`（本腿写面） | **空（0 行）** |
| 基线测试：`go test ./internal/config -count=1` | `ok github.com/CarlosShao/wisp/internal/config 1.124s`（16:35:38 +0800） |
| 基线绿名册（`-v` 抽 `--- PASS` 顶层） | `78` 枚，逐名录 `.scratch/wisp/probes/255/r1/baseline-green-roster.txt` |

**同机在飞声明**：`258-a1` 只读腿在飞（禁跑 Go、不互洗）——本腿跑 Go 尺时它不跑任何东西；本腿每一发尺前后
都复认 `git status --porcelain internal/config` 为空或只含本腿自己的新件。工作树里其余 366 行脏
（`.scratch`、`design`、根目录散件）全是既有痕迹，本腿一字未碰。

**排程预检**：票面排程节「按住」判据＝`248-v1` 退出＋`git status --porcelain cmd/wisp internal/config` 为空。
本腿实测 `git status --porcelain internal/config`＝0 行、`cmd/wisp` 不在本腿写面（AC#1/AC#5 归 `255-r2`），
写面冲突不存在，开工。

---

## §1 修法与判据（带 file:line）

（待填——实现进行中）

## §2 门禁四数

（待填）

## §3 变异自证表

（待填）

## §4 我可能写错的条目

（待填）

## §5 判不动的地方

（待填）

## §6 交件判语

（待填）
