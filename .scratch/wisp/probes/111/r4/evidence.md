# 票 111 — 111-r4 证据件（交票面 AC#1 那张全仓对账表）

代号 `111-r4`（只读普查腿）。本腿唯一交付物＝票面 **AC#1 的全仓对账表**，另附 (a) `go list ./...` 枚数现量＋35−33 差集、
(b) "空分母"那一族的现状、(c) winlive 半边与 cmd/wisp ubuntu 半边的现状登记。
前三格（闸门写进脚本 `1bb654e3`／闸门接进 CI `6c0e3e31`／闸门配自检 `98f62fde`）**本腿不重做**，只引其行号与 blob 凭据。

⛔ 本腿零改动面：`.github/**`、`scripts/**`、`cmd/**`、`internal/**`、票面五枚 `- [ ]` 框、台账、`docs/reports/**`、`docs/evidence/s1/**` 一字未动。
⛔ 本腿零 `go test`／`go build`／`go vet`；`go` 命令只用了 `go list`（含 `go list -f`，不编译测试二进制），给并行腿 `231-r1` 的计时窗让路。
写面＝`.scratch/wisp/probes/111/r4/**` ＋票 111 Progress log 一行。

---

## §0 起手锚（写满）

| 项 | 读数 | 时刻（+0800） |
|---|---|---|
| `git rev-parse HEAD` | `a16d1ff7291fa9ab59f328cbf5f2429184a881f4` | 09:14:54 |
| `git status --porcelain -- .github scripts cmd internal` | **0 行** | 09:14:54 |
| `git status --porcelain` 全仓 | **734 行**（共树在飞；本腿未 add 他人任何路径） | 09:14:54 |
| `go list ./...` 枚数现量（本机 GOOS=windows） | **35**（`| wc -l`＝35，stderr 0 字节） | 09:07:54 |
| `git rev-parse HEAD:scripts/portable-tests.sh` | `2ff02dd7476c8f623a4489932d0d0ff95480e667` | 09:14:54 |
| `git rev-parse HEAD:.github/workflows/ci.yml` | `014861149bde08ee8e4bb6e988f8a605796a856f` | 09:14:54 |
| `bash -n scripts/portable-tests.sh` | rc=**0**（语法自证，只读） | 09:14:54 |
| `bash -n scripts/portable-tests-selftest.sh` | rc=**0** | 09:14:54 |

★**锚点稳定性**（免得被读成"我在 734 行的脏树上量的门"）：上面两枚 blob 在 `1bb654e3`／`15d8b60e`／`6c0e3e31`／`b3d29b0d`／起手 HEAD
**五枚 HEAD 上逐字节同一**（`git rev-parse <ref>:<path>` 逐枚复量＝`2ff02dd7` 与 `014861149bde`）。
⇒ 本件里所有 `scripts/portable-tests.sh` 与 `ci.yml` 的行号对这五枚同时成立，不需要随 HEAD 漂移重量。
`go list ./...` 的**名册**（35 枚包路径）在 `15d8b60e..起手 HEAD` 之间同样未变（该区间 `-- cmd internal frontend tools` 的 `.go` 改动 4 枚，全在 `cmd/wisp` 包内，无新增/删除包目录）。

## §1 AC#1 全仓对账表 —— IN PROGRESS（下一笔补满）

## §2 附件 (a) 枚数现量与 35−33 差集 —— IN PROGRESS

## §3 附件 (b) "空分母"那一族的现状（GUARD A / GUARD D 各咬什么）—— IN PROGRESS

## §4 附件 (c) winlive 半边 ＋ cmd/wisp ubuntu 半边现状登记 —— IN PROGRESS

## §5 CI 读数列（第④列）：语义＝"已推送配置给过的结论" —— IN PROGRESS

## §6 量不到的格子（具名归口）—— IN PROGRESS

## §7 本腿 commit 链 —— IN PROGRESS
