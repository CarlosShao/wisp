# 票 223 非实现者对抗验收 `223-v2` —— AC#5 终态那一发 ＋ `223-r2` 修法的牙

- 工单：`.scratch/wisp/issues/223-checkandreload-has-zero-production-callers-hot-reload-never-runs.md`
- 被裁的实现件：`docs/evidence/s1/223-hot-reload-wiring-r2.md`（写腿 `223-r2`，收码 commit `4d925c2d` 一系）
- 上一枚验收腿：`docs/evidence/s1/223-hot-reload-wiring-v1.md`（判 AC#5 不成立、AC#4 两句归反）
- 更早的实现件：`docs/evidence/s1/223-hot-reload-wiring-r1.md`
- 本腿起手锚点：`git rev-parse --short HEAD` = **〔待填〕**，分支 `dev`
- 本腿性质：**只裁、只跑台件、不产码、不 commit 任何源码改动**（突变＝一进一出，逐字还原＋SHA256 自证）
- 台件目录：`.scratch/wisp/probes/223/v2/`（原始输出全量落盘，绝不接 `| head`／`| tail`）
- ⛔ 零读零写零转述：`frontend/**`、`design/**`
- ⛔ 不碰：`PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`、三枚冻结件、`docs/reports/**`、任何工单票面
- ⛔ 不碰别人的脏改动：`.gitignore`、`design/**`、`.scratch/wisp/probes/152/**`、`.scratch/wisp/probes/161/r6/logs/**`、`docs/evidence/s1/152-*.md`

## 入库清单（本腿会 commit 的路径，逐枚具名）

| 路径 | 是什么 |
|---|---|
| `docs/evidence/s1/223-hot-reload-wiring-v2.md` | 本表 |
| `.scratch/wisp/probes/223/v2/` | 全部原始读数与台件（含 `full-v.txt`、两发突变日志、J1 三形台件） |

---

## 起手读数

〔待填：起手 `date`／锚点／在飞检查／索引干净／别人地界的脏改动清单〕

## ① AC#5 缺的那一发：带 `-v` 的终态整包读数

〔待填：尺逐字、起手与终态 `date`、包级读数、逐名绿册、逐名红册、本票那枚用例在不在红册、`TestResolvePerCallBudget` 复量两发〕

## ② 攻 `223-r2` 那发正控与新轮询助手的牙

〔待填：正控独立重做（原始红句＋前后 SHA256）、`awaitStdout` 四问逐答带 `file:line`、AC#4 拿掉实现那发的红形名册 vs 实现腿自述〕

## ③ J1 那一句裁定的实测（三形）

〔待填：(a) 未来版＋语法坏／(b) 未来版＋解析得开／(c) 未来版＋未知键 的实际句子原文，裁定坐实或推翻〕

## ④ AC#1..AC#7 逐格重裁

〔待填：七格各一句＋依据行；重点核三笔＝`config_reload.go` 11/6 纯注释归真是否零行为、新增 124 行用例文件是否只测归句、AC#2 reload 档半格归属〕

## 没做完／判不了（具名清单）

〔待填〕

## 门禁与尺：时刻表

〔待填〕

## 交回编排者的六节（大白话）

〔待填〕
