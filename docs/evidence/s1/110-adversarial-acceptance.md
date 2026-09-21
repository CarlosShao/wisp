# 票 110 对抗验收裁决表（acceptor-110-97）

被验收对象：`.scratch/wisp/issues/110-no-ci-step-runs-internal-winsec.md` · commit `9e9a2f5`（已 push）
验收会话后缀：`ac110` · 全部快照在 `/tmp`（仓内零 worktree）
本表由独立对抗验收方产出，未参考先前对话；自述读数一律标档位。

档位图例：〔独立复现〕= 验收方自己跑出同形读数 · 〔日志＋归档，我抽验〕= 依赖票内日志/归档，验收方抽样核对 · 〔仅自述，不背书〕= 只有交件方一句话。

## 裁决表（与票面 AC 1:1）

| AC | 票面要求（摘要） | 裁决 | 证据档位 | 独立读数 |
| --- | --- | --- | --- | --- |
| AC#1 | 全 CI 每一步实跑哪些包 × `go list ./...` 对账，点名 winsec 0 次与其余零覆盖包 | TBD | TBD | TBD |
| AC#2 | 落一步真跑 winsec 测试 + 一次**步级**成功读数（run/job/step） | TBD | TBD | run `35595651898` / job `106319703680` / step 4 |
| AC#3 | 新步自己会红 + 不是空仪器（两发变异，同链 grep + `go build` rc=0） | TBD | TBD | TBD |
| AC#4 | `R-93-4` 收口：`TestSyncRegistryProbeLive` 的步级证据或"不该纳入"的论证 | TBD | TBD | TBD |
| AC#5 | `bash -n` rc=0；`sh scripts/d22scan.sh` 纯净快照 rc=0 且台账不降；新步不得排在常红步之后 | TBD | TBD | TBD |

## 步骤位置判定（AC#5 关键，验收方自己读 `ci.yml` + 真实 job steps）

`test-windows` 实排（来自 run `35595651898` / job `106319703680` 的 `.steps[].name`，与 `.github/workflows/ci.yml` 一致）：

1. Set up job · 2. actions/checkout@v4 · 3. actions/setup-go@v5 · **4. Windows ACL sealing gate (internal/winsec's own tests, ticket 110)** · 5. Cache third_party · 6. cgo build smoke · 7. Portable windows tests (proc/secret/config/risk) · 8. PathResolver junction placeholder

⇒ 新步之前只有 checkout / setup-go，无第三个可常红步骤（`if: always()` 的必要性判定：TBD）。

## 同 run 的其它步级读数（区分"旧红"与"新步的账"）

TBD（lint / test-core 红名归因）。

## R-110-x 登记

TBD

## 未验完 / 断点

本文件为第一枚 checkpoint 骨架：AC#2 步级读数待 run 跑完回填，AC#1 复算与 AC#3 变异在进行中。
