# 300-v4 — anchor gate (non-implementer adjudication leg)

现量时刻：本腿起手。

## 锚
- `git rev-parse HEAD` = `34e1962bb166b4494670b77296263c79bd8671a3` (短 `34e1962b`) — 与派单声称一致。
- 工作树**脏**（`design/**` 删除、`.scratch/**` 修改等，皆他人正在飞的活）⇒ 本腿**只读**，⛔ 任何 checkout/reset/stash/clean。
- 本腿允许写面：`.scratch/wisp/probes/300/v4/**`（新建）。⛔ 起名 `.out`（根 `.gitignore:8` 是全仓 `*.out`）。

## 进程闸门现量（rc=0）
- `tasklist /FI "IMAGENAME eq wisp.exe"` → `INFO: No tasks are running which match the specified criteria.` ⇒ **0**
- `tasklist /FI "IMAGENAME eq balldebug.exe"` → 同上 ⇒ **0**
- `tasklist /FI "IMAGENAME eq mockllm.exe"` → 同上 ⇒ **0**
- `tasklist /FI "IMAGENAME eq msedgewebview2.exe"` 计数 = **24** ⇒ 按台账先例 `A821` **⛔ 判残留**（挂在机主自己的应用树下），⛔ 杀任何进程。

⇒ 起手闸门**放行**，零归属不明。
