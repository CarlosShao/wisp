# D32 第一份真结论 — CI run `37249563077` 的 `slo-full` 摘要（编排者 10-05 09:0x 提取）

**为什么单独留这份摘要**：原件 645,665 字节不入库（它每发 run 可再生，不是不可再生物），但它是本仓**第一份**走到结论的 D32 读数（此前每一次都在 `$LASTEXITCODE` 那一句上死掉，`Upload SLO report` 全是 skipped）。原件的逐字节指纹记在下面，下一任若怀疑摘要，可用同一枚 artifact 重取比对。

## 出处（逐条可对）

| 项 | 读数 |
|---|---|
| run | `37249563077`（event＝push，workflow＝ci） |
| headSha | `21bec8a1c98c69a0de0221d87052918dd234ba06`（＝编排者 08:58 推上双远端的那枚 HEAD） |
| job | `slo-full`＝id `111574243797`，conclusion＝**success**；同发 `slo-smoke`＝id `111574243677`，也是 success |
| 步 | `SLO full gate (six states + settle + leak)`＝success；**`Upload SLO report`＝success**（不是 skipped） |
| 产物 | artifact `slo-full-report`，id `11320002666`，20,514 字节（zip） |
| 展开件 | `slo-report.json`＝**645,665 字节**，带 UTF-8 BOM（须以 `utf-8-sig` 读） |
| SHA-256 | `c94334cdada48f05efbc224ef163202d75f2264facd91a6625c5e601b7b36394` |
| 机器 | `DESKTOP-LVS7839`（self-hosted runner，＝这台机本机） |
| `subset` / `all_pass` / `generated_at` | `full` / **true** / `2026-10-05T01:01:47Z`（本地 09:01:47 +08） |

## 六档逐档读数（每档 24 枚样、`sample_errors=0`、`duration_sec≈6.01`、`posture="skeleton"`、`measurement_basis="out-of-tree"`，另有第二读者 `observer`）

| 档 | 内存中位 | 帽 | CPU 均值 | 帽 | 该档门数(gate=true) | 该档不响的判 |
|---|---|---|---|---|---|---|
| Sleeping | 4.5MB | <=25MB | 0.000% | <=0.5% | 7 | 无 |
| Armed | 4.5MB | <=90MB | 0.000% | <=2.0% | 5 | 无 |
| Warm | 4.6MB | <=350MB | 0.000% | <=1.0% | 4 | 无 |
| Conversation | 4.6MB | <=350MB | 0.000% | <=5.0% | 4 | 无 |
| PanelOpen | 4.5MB | <=600MB | 0.000% | <=10.0% | 4 | 无 |
| WorkPeak | 4.6MB | <=700MB | 0.000% | <=30.0% | 3 | 无 |

**settle**（`exit_code=0`、`pass=true`）：`target_state=Sleeping`、`cap_bytes=26,214,400`、`peak_bytes=131,846,144`、**`back_within_cap_ms=257`**、`final_bytes=9,117,696`、`free_os_memory_count=2`、`elapsed_ms=10007`、40 枚样、`sample_errors=0`。
**leak 自检**（门自己的阳性对照）：**`flipped_to_fail=true`**＝种 100MB 那发真把门打红。⚠ 票 263 §Progress log 具名过：这一族今天只剩"leak 要求 `exit==1`"这枚兜底，删钉那一形仍可被伪造 `exit=0`（⇒ 归票 266 的射程，⛔ 本票不修）。

## 这份读数买到什么、没买到什么（⛔ 下一任不许读成"D32 达标"）

- **买到**：票 263 那条因果链（"求值器第一次真跑到它就死"）的**反面**——求值器活着、六档全评、settle 与 leak 都真跑，且读数可逐字节复核。
- **没买到**：被测对象的完整度。六档 `posture` 全是 `skeleton` ⇒ 今天**没有真麦克风、没有面板宿主在泵、没有 LLM 在跑**。⇒ 这四个数是"**空骨架在预算内**"，不是"产品达标"；票 247（电平接线）／票 33（面板接进常驻）／票 197（子代理三层）落地之后，同一把尺的颜色才是那一格的答案。
- 另注：`cpu_mean_percent` 六档全 0.000% 与 `mem_median≈4.5MB` 一起看＝骨架确实安静；但**样本只有 6 秒/档**（CI 用 `-SecondsPerState 6`），⛔ 不许把它当"长时间驻留"的证据（D32 里那条驻留判定另有口径）。
