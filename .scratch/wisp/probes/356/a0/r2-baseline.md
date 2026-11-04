# 356 AC#0 探针 · R2 正控成对读数（2026-11-04，编排者现跑）

harness：`internal/tools/task_output_pointer_two_roads_174r1_test.go` 的 `twoRoadsBridge174r1`（真 Bridge／真 C26／真 fs.read／真 gate seam）＋ `readThroughBridge`；审批计数能力＝`gateSpy.counts() (windows, approvals)`（`bridge_test.go:95`）。

## 命令
`go test ./internal/tools/ -run "TestPointerSecondRoadBackIsRealOnThisBridge174r1|TestPointerOutsideAllowlistNamesEveryRoadBack174r1" -v -count=1` ⇒ **PASS**。

## 读数（逐字）
1. **allowlist 臂 notice 原文**（out-of-allowlist，R2 形）：
   `IsError=false Truncated=true` Text=「…注意：这条路径现在读不到，它不在你被授权的目录范围内：fs.read 会要一张 L2 卡，没人批就是直接拒；回来的路有两条——要么用户把所属目录加进 [fs] allowed_dirs，要么批下那一张卡，全文见 <path>」
2. **road two，无人批（NoGate）**：`IsError=true level=L2 class="user_rejected"` text=「L2 审批通道尚未接入（票 21），已拒绝执行」。
3. **road two，卡片批准**：`IsError=false level=L2 bytes=20000`（整读回）。

⇒ R2-only 形：「只判到 L2，两条路都真」——NoGate 拒（L2/user_rejected）、批准即整读回 20000B。这是 AC#0 的**正控**那一半。

## 待做（AC#0 另一半：R3 A 档＝绝对 Deny）
R3 A 档「没有卡、批也无效、approval 被问 0 次」的**成对读数尚未取**：需按 `internal/tools/bridge_junction_windows_test.go:499-513` 的形（A-list 文件落在 allowed root 内、R2 满足、只 R3 抓得到 ⇒ `Decision.Level==risk.Deny && RulesHit∋R3`）构造 Windows junction 形台件，再折回 `task.output` 的 `pointerNotice`：预期 notice 仍吐「两路」文案（`!InAllowlist` 支先命中、无 Deny 支），而 `fs.read` 给 `permission_denied`、approval 计数 0。⛔ AC#0 未翻（半读数不够，须两形齐）。
