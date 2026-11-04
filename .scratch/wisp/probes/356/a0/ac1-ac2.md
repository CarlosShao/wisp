# 356 AC#1 + AC#2 落地证据（2026-11-04，编排者）

## AC#1：回执区分三种形状（i 只判 L2 两路真 / ii 越权+Deny / iii 授权内+Deny）

改法（判据用**分类**不用 RiskLevel 串）：
- `internal/tools/paths.go` 新增 `PathCanonicalizer.DenyA(canonical) bool`＝`risk.Classify(canonical)==risk.ClassA`（A 档绝对 Deny）。只读分类、不扩授权（已在 `refusal_20r2_test.go` 的 read/judge 面名册登记，守卫测绿）。
- `internal/tools/task.go` 的 `pointerNotice`：switch 增 `case d.Paths.DenyA(canon):`（置于 `!InAllowlist` 之前），吐「绝对禁止…不会有卡，加 [fs] allowed_dirs 或批卡都无效」。
- `internal/agent/spill.go` 的 `pointerNotice` 同形加 `if s.judge.DenyA(canon)`；`PointerJudge` 接口加 `DenyA(canonical string) bool`；两个测试 fake（`argAware174r4Judge`/`fake174Judge`）补 `DenyA`（返回 false 保旧行为）。

**逐字同形**：两腿 deny 句逐字相同（现跑 `[ "$A" = "$B" ]` → SAME，len **184**）。
**探针三形**（`zz_356_ac0_probe_test.go` tag `probe356`，现跑 PASS）：
- i（R2-only）：notice 两路、approval=1、整读回 20000B。
- ii（越权+Deny）：notice 现吐「绝对禁止…不会有卡…批卡无效」（原为两路）、fs.read permission_denied、approval=0。
- iii（授权内+Deny）：notice 现吐同句（原为无告警/沉默）、fs.read permission_denied、approval=0。

## AC#2：补 MUT-5 缺口——只属第二路的标记

- 新测 `internal/tools/task_output_pointer_second_road_pin_356_test.go` 断言 R2 回执含**第二路整句** `要么批下那一张卡`（旧标记集 {不在你被授权的目录范围内, [fs] allowed_dirs, L2 卡} 全在句首、删第二路仍绿＝MUT-5 缺口）。
- **变异自证**：删 `task.go` 的「，要么批下那一张卡」→ `go build` OK、pin 测 **RED**（报 `task_output_pointer_second_road_pin_356_test.go:29: ...second-road clause ... missing`）；从备份还原后 `git diff` 空、count 回 1。旧 174-r1 标记在变异下发仍绿（实证 MUT-5），新 pin 独红。

## 门禁（AC#3 相关）
`go build ./cmd/wisp` OK；`go vet ./cmd/wisp` 净；`gofumpt -l` 被碰文件零输出；`d22scan` rc=0 clean；`go test ./internal/tools/ ./internal/agent/ -count=1` 全 **ok**（tools 22.1s／agent 6.7s）。⚠ `go vet ./internal/tools/` 仅余 `system_backend_windows.go` 两处 unsafe.Pointer **基线**告警（该文件未碰）。
