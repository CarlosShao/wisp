# 票 223 写腿 `223-r2` —— AC#5 间歇红改确定性 ＋ AC#4 两句归正

- 工单：`.scratch/wisp/issues/223-checkandreload-has-zero-production-callers-hot-reload-never-runs.md`
- 前案裁决表（必读，本腿的射程由它划出）：`docs/evidence/s1/223-hot-reload-wiring-v1.md`（56,522 字节／265 行）
- 本腿起手锚点：`git rev-parse --short HEAD` = **`c774f8da`**，分支 `dev`
- 本腿性质：**写腿**，只做派单的两件事（AC#5 轮询化／AC#4 句子归正＋常驻用例），⛔ 不扩射程
- 台件目录：`.scratch/wisp/probes/223/r2/`（原始输出全量落盘，不接 `| head`／`| tail`）
- 入库 pathspec（只此四枚）：`cmd/wisp/**`、`internal/config/**`、本表、`.scratch/wisp/probes/223/r2/`、票 223 工单面（只追加）
- ⛔ 零读零写零转述：`frontend/**`、`design/**`；⛔ 不碰：`PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`、三枚冻结件
- ⛔ 票 223 不翻 `-done`、不勾任何 AC 框（翻勾由编排者按非实现者表＋自己复跑裁）

## 起手读数（本腿现跑）

- 起手 `date`：〔待填〕
- `git diff --cached --name-only` 起手：〔待填〕
- `ps -W` 在飞门检查：〔待填〕
- 别人地界的脏改动（起手 `git status --porcelain` 现读、本腿一律不碰）：`.gitignore`、`design/**`、`.scratch/wisp/probes/152/**`、`.scratch/wisp/probes/161/r6/logs/**`、`docs/evidence/s1/152-*.md`

## 任务一：`TestTicket223RestartTierSaysItWillNotApply` 间歇红 → 有界轮询

### 1.1 复现（先确认红的是读时序，不是产品句子缺失）

〔待填：红发读数、红句原文逐字〕

### 1.2 修法

〔待填：改动位置、轮询形状、期限怎么写（monotonic timer，非墙钟减法）、为什么不动产品文案〕

### 1.3 二十连跑判据（同一把尺）

〔待填：逐发 exit 码，原始件 `.scratch/wisp/probes/223/r2/flake-20.txt`〕

### 1.4 正控（拿掉那句产品文案必须红）

〔待填：突变形状、红句原文、还原前后 certutil SHA256 两个哈希逐字〕

## 任务二：AC#4 两句会说反的话归位

### 2.1 现读归因（loader.go 的＋72/−6 与句子分支）

〔待填：读 `internal/config/loader.go` 的 peekSchemaVersion/declaredSchemaVersion、`cmd/wisp/config_reload.go` 的 describeReloadFailure、v1 腿 overlay 17 发读数 `ac4-probe.txt`〕

### 2.2 修法

〔待填：改哪一处、三条形（读得出当前版＋正文语法坏→语法错；读不出版本→config.toml parse；版本低于当前→迁移管线）、migrate_test.go:123 不受影响的论证〕

### 2.3 新增常驻用例

〔待填：用例名、形状、四形各归其句的断言〕

### 2.4 牙（把修法拿掉必须红）

〔待填：红句原文逐字、还原读数〕

## 三把门（原始结论落盘）

| 门 | 读数 | 原始件 |
|---|---|---|
| `go build ./...` | 〔待填〕 | 〔待填〕 |
| `gofumpt -l <本腿动过的文件>` | 〔待填〕 | 〔待填〕 |
| `sh scripts/d22scan.sh` | 〔待填〕 | 〔待填〕 |

## 定性（涉及"会说反的话"字样，按派单要求自带三行）

〔待填：①现象出现在哪 ②有没有本机被入侵的证据 ③最坏后果是什么形状〕

## 没做完／判不了（具名）

〔待填；含：若 20 连跑里冒出 `always_write_no_clobber_226_test.go`（票 227 射程）的红——不修、具名登记〕

## 时刻表（本腿现跑）

| 时刻(`date` 现读) | 命令 | 落盘 | 结果 |
|---|---|---|---|
| 〔待填〕 | | | |

## 交回编排者的六节（大白话）

**① 两件事各自动了哪几枚文件、行数**〔待填：`git diff --numstat` 原文〕
**② 任务一 20 发逐发 exit 码与正控红句**〔待填〕
**③ 任务二改前／改后两句原文并排＋新增用例有牙读数**〔待填〕
**④ 三把门原始结论**〔待填〕
**⑤ 票 223 该不该翻勾**〔待填——本腿是写腿，不自翻勾；只报"我这边两格缺口堵没堵上"〕
**⑥ 没做完／判不了**〔待填〕
