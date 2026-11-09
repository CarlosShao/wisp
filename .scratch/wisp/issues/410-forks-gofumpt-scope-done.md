# 410 — `forks/go-webview2` 8 枚 gofumpt 红：射程裁定与落地（甲＝格式化，非乙＝exclude）

**Status:** 落地待收（本腿只 commit 不 push；编排者收件后追 `-done`）
**Claimed by:** Wisp 写腿（2026-11-08 23:1x 领单）
**Last update:** 2026-11-08 23:4x

## 背景与立案理由

276 腿现跑 `bash .scratch/wisp/probes/161/r5/attrib.sh --tracked-only` ⇒ **RED rc=1**，红点＝`forks/go-webview2` 8 枚
（`2bb639e4e` 11-07 票65 CarlosShao，已推 origin/dev）。A1556 记录的「门今天 green」**已再变**。
派单问「归票 70 还是新立案」：**新立案**——票70 的 AC#1 格式化 sweep 已于 09-21 闭（`6ce43c7`，当时 fork 尚不存在——fork 是 11-07 才进的树），
票70 剩余项全是 blocked-on-owner 的 CI run 证据（ubuntu 归属／staticcheck pin／push 后五 job 全绿），与本裁定不同轴。

## 判决（甲，推翻派单的乙倾向；逐条理由见证据件）

**乙（exclude `forks/**`）不成立**：其判据「本仓没有作者权」事实为假——票65 已写 `Close()/Release()` 包装、
票408 r3（`da5d559b7`）已在本仓作者权下改 `pkg/edge/chromium.go`（nil-env 守卫，gofumpt 干净）。
exclude 等于把本仓持续写码的目录摘出格式门。**甲的代价（上游同步冲突）是伪命题**：`go.mod:24` replace 硬叉＋pseudo-version 钉死，
票65 证据件逐字「上游更新＝手动同步一次」，语义层已分叉，8 枚格式 hunk 只是噪音。先例＝票70 AC#1「放宽门不是修门」。

## Progress log

- 2026-11-08 23:4x（写腿 r1）：起手三核全现跑（红名单恰 8 枚／exclude 段=:560-600 仅 `.scratch/**` 一枚具名排除／`git log` 证 fork 非纯直搬）；拍甲落地：`gofumpt -w forks/go-webview2` 恰 8 件 +16/−12 纯格式，commit `504858fd1`，`attrib.sh` 一字未动；导出树（`git clone --shared`）验收 **attrib.sh --tracked-only GREEN rc=0，四数＝lines=0/files=0/handed=1286/total=1657(excluded 371)**；gofumpt 改动件 0 行、vet rc0、build rc0、d22scan clean、8 件全 LF。证据件 `docs/evidence/s1/forks-gofumpt-scope.md`。只 commit 不 push。
