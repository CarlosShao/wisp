# 55 — S8: macOS port (second-class platform) — DEFERRED until open-source prep

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 54-s7-acceptance + self-use week + explicit user go-ahead (S8 starts only
pre-open-source, D23)
**Parallel slots:** ≤2 sub-agents
**Spec refs:** D7, D23, DEFERRED registry macOS row, SPEC-12 §5

## What to build
The macOS platform layer: borderless NSWindow + `level = .floating` ball, Keychain-backed
SecretStore (S8 row of C28), symlink-realpath path rules (18 stub), and CI cross-builds with
prebuilt artifacts.

## Key constraints
- D7 stance unchanged: macOS ships as "second-class, unverified on real hardware" unless a
  maintainer with hardware joins; README must say so.
- Ball parity: 20 states via NSWindow; layer-shell Linux still RESERVED (do not start).
- CI: darwin/amd64 + arm64 artifacts (cgo strategy per BUILD.md extension); notarization
  DEFERRED (Gatekeeper docs `xattr` workaround, D17).
- All DEFERRED rows' completion criteria from SPEC-12 §5 must be met before S8 closes.

## Acceptance criteria
- [ ] Both ball implementations pass interaction tests; CI produces signed-optional macOS
      artifacts; CER-equivalent English models out of scope (i18n is 58).
- [ ] Keychain SecretStore passes 06's test matrix ported.
- [ ] PathResolver macOS branch (realpath + lstat) passes red-team port.
- [ ] README platform matrix updated honestly.
- [ ] P12/macOS sync-client probing (registered by ticket 19's DEFERRED(P12-macos),
      adversarial report N-6): `~/Library/CloudStorage/<Provider>` plus per-client config
      probing must yield env/registry/config-GRADE evidence — a bare `os.Stat` of a guessed
      default directory does NOT disarm the under-profile sync-suspect fallback
      (`gradeConfirmed` in internal/risk/syncdirs.go). Acceptance: a relocated-Dropbox and a
      relocated-OneDrive fixture both resolve to a confirmed root, and
      `internal/risk/syncdirs_redteam_*_test.go` spellings equivalents pass on macOS.

## Progress log (append-only, newest last)

- 2026-09-21 17:1x（编排者，票 82 结案时追加的**验收锚点**，不是新工作）：票 82 已经把
  `internal/risk` 的 sync-root 家族在 POSIX 上**分层**了——Windows 专属判据进 `//go:build windows`，
  POSIX 侧留了三条"断言当前真实行为"的用例，并把它们做成**票 55 落地即红**的 tripwire：
  **`TestSyncNoGradeIsConfirmedOnPosix`**、**`TestSyncMembershipDecidesOffProfilePosix`**、
  **`TestSyncRegistryProbeIsAStubHere`**（`internal/risk/syncdirs_other_test.go`）。
  ⇒ **macOS 移植交件时必须逐条点名这三条**：它们变红是**预期且有意义的**（POSIX 侧"还没实现"这个事实被实现
  打破了），修法是**把断言升级为真实判据**，**绝不是**放宽或删掉。⚠ 前两条断言的是"没有确认等级 / 成员关系
  决定嫌疑"，第三条断言的是"registry 探针在此平台是桩"——macOS 若有自己的等价物（iCloud Drive 同步根），
  第三条要改成 macOS 层的真实探测用例而不是直接删。


- 2026-09-21 16:2x（编排者，**票 82 验收带回的第四格锚点，macOS 交件时一并点名**）：
  `acceptor-ticket82` 在 ubuntu 上做变异 **MUT-c**：把 `gradeConfirmed` 加宽去认 `"default"` 这一档
  ⇒ **ubuntu rc=0 全绿，而 Windows rc=1**。含义：票 82 定性表里"弱等级永不拆兜底网——两侧都有对象"那句
  在 **Linux 侧其实是空仪器**（它对"弱等级"这个概念没有可断言的对象，所以加宽它不被发现）。
  ⇒ **票 55 落地时这一格必须变成有牙的**：macOS 有自己的同步根来源（iCloud Drive 等），
  届时"某个弱来源被误升为 confirmed ⇒ 兜底网被拆"要有一条**在 macOS 上会红**的用例，
  而不是继续靠 Windows 一侧的 rc=1 独扛。判据形状照票 82 的 MUT-c 反向做：**加宽 ⇒ 用例红**。
  与上面那三条 POSIX tripwire 同批点名。
