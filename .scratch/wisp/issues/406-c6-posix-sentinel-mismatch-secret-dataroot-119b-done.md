# 406 — C6 孤儿：POSIX 上 TestAC3POSIXSecretRoute...119b 的哨兵不一致（ErrReparseDenied vs winsec.ErrUnresolvedPath），票 119 已 done、无票主

**Status:** **done**（2026-11-05 15:5x 机主裁定 **③＝维持各认各的并登记、不改代码**，A1544。来路＝98 普查 C6，119 已 accepted-done、孤儿。裁定后本票结案；下方那枚 POSIX 红＝已登记 known-red，非本编队修）。
**Type:** 结后遗留红（done 票引入的测试/代码不一致），非新功能。
**Packages:** 候选三面均顶：【internal/risk pathresolver.go:36 ErrReparseDenied（冻结面）】/【internal/winsec POSIX 半边（票94 已done）】/【cmd/wisp secret_dataroot_119b_test.go:241-242（!windows build，整包 cmd/wisp 死区）。
**Blocked by:** 写面顶冻结面 ⇒ 须 owner 先裁：①解冻 internal/risk 统一哨兵 ②改 cmd/wisp 测试认 ErrReparseDenied ③维持各认各的并登记。

## 现场
docs/evidence/s1/98-attribution.md:65 现读：POSIX 上拒绝确实发生，但断言 errors.Is(err, winsec.ErrUnresolvedPath) 为假——实际哨兵是 internal/risk 的 ErrReparseDenied（pathresolver.go:36）⇒ 测试 :242 红。

## 建议（编排者按 119 的甲）
写面顶 frozen ⇒ 本票先挂台账待 owner 三选一（见 Blocked by）。所有者裁定前勿派写腿（避免第 22 处重派）。

## 机主裁定（2026-11-05 15:5x，A1544）＝③
- **③ 维持各认各的并登记、不改代码。** 理由：功能安全（POSIX 上拒绝确实发生），只是**测试期望的错误代号写错**（断言 winsec.ErrUnresolvedPath，实际 risk.ErrReparseDenied）。
- 代价：POSIX 腿保留 **1 枚 known-red**（`TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119`），登记在案、不算回归、不再重派。
- 不做 ①（解冻 internal/risk）/②（改 cmd/wisp 测试）—— 为一条测试期望代号去动冻结区/禁写区不值当。
- 撤销口令：机主改口「406 走②」即可重开改 cmd/wisp 测试断言。
