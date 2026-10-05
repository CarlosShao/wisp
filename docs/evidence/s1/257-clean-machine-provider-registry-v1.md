# 257-v1 终裁件 — 票 257（干净机器上设置写不进去／形 ⓒ 首启回执）的对抗验收

> 腿：`257-v1`（非实现者终裁，`SPEC-12 §4.3` #1／`AGENTS.md` §0.3）。
> 票：`.scratch/wisp/issues/257-clean-machine-provider-registry-nil-blocks-writes.md`（AC 框一枚未碰）。
> 被裁的产码：`cmd/wisp/firstrun.go`（＋61 行，只进 stderr）；被裁的测试：
> `cmd/wisp/firstrun_257_test.go`（六枚）与 `cmd/wisp/firstrun_257_nonpreset_test.go`（两枚）。
> 本件只出判语，翻勾与删不删都归编排者。

## 0. 起手锚

| 项 | 读数（本腿现跑，逐字） |
|---|---|
| 进场时刻 | `2026-10-05 12:20:31 +0800`（`date` 自取） |
| 进场 HEAD | `5b638498`（派单给的锚是 `21bec8a1`→HEAD；这 15 分钟内链头已被别的腿推走，见 §6 用绝对 range） |
| 分支 | `dev` |
| `git status --porcelain -- cmd internal` | **0 行**＝`cmd/**`＋`internal/**` 此刻无别人的活在树上（12:20:51 现量） |
| `go env GOFLAGS` | 空串（两次现量 12:20:31／12:20:51 相同） |
| `git status --porcelain -- internal/agent/approval` | **0 行**＝`259-r2` 此刻在该包无未提交活（12:20:51） |
| 并发腿 | 派单告知 `259-r2` 在 `internal/agent/approval` 做注释级改动，会被 `cmd/wisp` 编译带上 ⇒ 本腿两发整包前后各取一次上面两行（§9、§10） |
| 本腿写面（授权） | `docs/evidence/s1/257-clean-machine-provider-registry-v1.md`＋`.scratch/wisp/probes/257/v1/**`；突变只许种进 `cmd/wisp/**`（⛔ 不种 `internal/agent/approval`／`internal/config`／`internal/panel`） |

## 1. 编队容量（起计时敏感腿前的三发，⛔ 任一 ≥70 % 就不抢机器）

`powershell -NoProfile -ExecutionPolicy Bypass -File .scratch/fleet-load.ps1` 连取三发（逐字）：

| 发 | 时刻 | 读数 |
|---|---|---|
| 1 | 12:20:31 | `CPU=58% MEM=56.4%` |
| 2 | 12:20:49 | `CPU=67% MEM=56.5%` |
| 3 | 12:20:51 | `CPU=60% MEM=56.6%` |

⇒ 三发均 <70 %，本腿放行去跑计时敏感那两枚；CPU 67 % 那一枚贴边，故 §9 那两发整包**各自起腿前再取一发**，任一 ≥70 % 就停下登记。

## 2. ★裁定一：两测试套是不是同一枚判据做了两遍、哪一套算交件

判语：未判。凭据：未判。

## 3. AC#1 干净机真跑（无 `config.toml` 起步＋ ⓒ 三样指引兑现）

判语：未判。凭据：未判。

## 4. AC#2 三种拒因三句话（含定向突变）

判语：未判。凭据：未判。

## 5. AC#3 凭据面（不新增回显、不落明文）

判语：未判。凭据：未判。

## 6. AC#4 越界（`21bec8a1..HEAD` 禁区文件名命中数）

判语：未判。凭据：未判。

## 7. ★裁定三：那 61 行有没有把指引漏进生成的 `config.toml`（实弹突变）

判语：未判。凭据：未判。

## 8. ★裁定二：正向断言（"提示真到达 stderr/stdout"）是不是恒真

判语：未判。凭据：未判。

## 9. 那枚新增红＋那枚命名跳的复跑两形与归因

判语：未判。凭据：未判。

## 10. 门禁四门（本腿自跑，不引用腿的读数）

判语：未判。凭据：未判。

## 11. 突变名册（每发起止时刻到秒，编排者拿它作差）

未判。

## 12. 七项对抗检查（`SPEC-10 §8`）

未判。

## 13. 我攻不动的地方／判不动的格子（具名归口，⛔ 一枚不写成"应该没问题"）

未判。

## 14. 交件判语与 commit 链

未判。
