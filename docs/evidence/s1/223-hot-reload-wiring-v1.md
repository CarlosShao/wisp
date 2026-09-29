# 票 223 非实现者对抗验收 `223-v1` —— CheckAndReload 热加载接线

- 工单：`.scratch/wisp/issues/223-checkandreload-has-zero-production-callers-hot-reload-never-runs.md`
- 被裁的实现件：`docs/evidence/s1/223-hot-reload-wiring-r1.md`（实现腿 `223-r1` 死于 150 轮上限；AC#5 由编排者代跑；五枚文件未提交改动由编排者代落＝commit `248095d1`）
- 落点普查前案：`.scratch/wisp/probes/223/c1/census.md`
- 跨票前案：`.scratch/wisp/issues/226-always-answer-rewrites-the-whole-config-snapshot-and-hides-hand-edits.md`／`docs/evidence/s1/226-config-write-no-clobber-v2.md`／`internal/config/writeguard.go`
- 本腿起手锚点：`git rev-parse --short HEAD` = **`577ae8c7`**，分支 `dev`，`date` = 见下方起手时刻表
- 台件目录：`.scratch/wisp/probes/223/v1/`（原始输出全量落盘，绝不接 `| head`／`| tail`）
- 本腿性质：**只裁、只跑台件、不产码、不 commit 任何源码改动**
- ⛔ 零读零写：`frontend/**`、`design/**`（本腿对这两目录不读、不引、不转述）
- ⛔ 不碰：`PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`、三枚冻结件

## 入库清单（本腿会 commit 的路径，逐枚具名）

| 路径 | 是什么 |
|---|---|
| `docs/evidence/s1/223-hot-reload-wiring-v1.md` | 本表 |
| `.scratch/wisp/probes/223/v1/` | 全部原始读数与台件（含 `full-v.txt`、突变日志、loader 对抗输入台件） |

---

## 起手读数（本腿现跑，不复用他人数字）

〔待填〕

## AC#1 生产里真有人在轮询

〔待填〕

## AC#2 三档生效级别各有读数

〔待填〕

## AC#3 放宽必带 L2 复确认（D33 正控）

〔待填〕

## AC#4 不生效与读不到是四句话

〔待填〕

## AC#5 整包终态读数（任务一：带 `-v`）

〔待填〕

## AC#6 裁决不在锁内（任务二：专属突变）

〔待填〕

## AC#7 "重启后生效"那一档要有出口

〔待填〕

## 编排者代落＋代跑这两件事的独立裁决

〔待填〕

## 跨票攻击：票 226 的"手改／不认领"在真轮询路径下是否仍成立

〔待填〕

## 推翻清单（实现腿＋编排者）

〔待填〕

## 票面写的 vs 我读到的（不一致单列）

〔待填〕

## 没做完／判不了（具名清单）

〔待填〕

## 门禁与尺：时刻表

〔待填〕
